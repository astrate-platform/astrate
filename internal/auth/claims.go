package auth

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claim names an Astarte authorization claim: the JWT key under which a list
// of authorization strings is carried (docs/DESIGN.md §4.2).
type Claim string

// The Astarte claim set. Each claim authorizes one API surface; the values
// are the exact JWT keys upstream tooling (astartectl, astarte-go) emits.
const (
	// ClaimAppEngine authorizes the AppEngine API (a_aea).
	ClaimAppEngine Claim = "a_aea"
	// ClaimChannels authorizes Astarte Channels; Astrate honours it on the
	// live stream socket (a_ch).
	ClaimChannels Claim = "a_ch"
	// ClaimHousekeeping authorizes the Housekeeping API (a_ha).
	ClaimHousekeeping Claim = "a_ha"
	// ClaimPairing authorizes the Pairing agent API (a_pa).
	ClaimPairing Claim = "a_pa"
	// ClaimRealmManagement authorizes the Realm Management API (a_rma).
	ClaimRealmManagement Claim = "a_rma"
)

// stringList unmarshals an authorization claim value. Upstream always emits
// JSON arrays of strings; a bare string is tolerated as a single-element
// list (a pure superset — upstream Guardian would reject it, accepting it
// breaks no valid token).
type stringList []string

// UnmarshalJSON implements json.Unmarshaler.
func (s *stringList) UnmarshalJSON(b []byte) error {
	trimmed := strings.TrimSpace(string(b))
	if trimmed == "null" {
		*s = nil
		return nil
	}
	if strings.HasPrefix(trimmed, "[") {
		var list []string
		if err := json.Unmarshal(b, &list); err != nil {
			return err
		}
		*s = list
		return nil
	}
	var single string
	if err := json.Unmarshal(b, &single); err != nil {
		return fmt.Errorf("authorization claim must be a string or a list of strings: %w", err)
	}
	*s = []string{single}
	return nil
}

// astarteClaims is the JWT claims document: registered claims plus the five
// Astarte authorization claims.
type astarteClaims struct {
	jwt.RegisteredClaims

	AppEngine       stringList `json:"a_aea,omitempty"`
	Channels        stringList `json:"a_ch,omitempty"`
	Housekeeping    stringList `json:"a_ha,omitempty"`
	Pairing         stringList `json:"a_pa,omitempty"`
	RealmManagement stringList `json:"a_rma,omitempty"`
}

// grant is one compiled authorization string: a verb regex and a path regex,
// both implicitly anchored.
type grant struct {
	verb *regexp.Regexp
	path *regexp.Regexp
}

// compileGrant compiles one authorization string with upstream parity
// (Astarte's GuardianAuthorizePath plug):
//
//   - the string splits on ":" into exactly three parts
//     "<verb-regex>:<opts>:<path-regex>" — the canonical form
//     "<verb>::<path>" is the empty-opts case, and the opts field is
//     ignored;
//   - each part is anchored by literal concatenation, "^" + part + "$"
//     (no grouping — "a|b" becomes "^a|b$", left-anchored a OR
//     right-anchored b, exactly as upstream's Regex.compile does);
//   - any malformed string (wrong arity, regex that does not compile)
//     simply never matches; it does not invalidate the token.
//
// ok reports whether the grant is usable.
func compileGrant(s string) (grant, bool) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) != 3 {
		return grant{}, false
	}
	verb, err := regexp.Compile("^" + parts[0] + "$")
	if err != nil {
		return grant{}, false
	}
	path, err := regexp.Compile("^" + parts[2] + "$")
	if err != nil {
		return grant{}, false
	}
	return grant{verb: verb, path: path}, true
}

// Token is a verified JWT: its expiry and its compiled authorization grants.
// Tokens are immutable and safe for concurrent use, which is what allows the
// LRU cache to hand the same *Token to many requests.
type Token struct {
	expiresAt *time.Time
	grants    map[Claim][]grant
}

// newToken compiles verified claims into an immutable Token.
func newToken(c *astarteClaims) *Token {
	t := &Token{grants: make(map[Claim][]grant, 5)}
	if c.ExpiresAt != nil {
		exp := c.ExpiresAt.Time
		t.expiresAt = &exp
	}
	for claim, raw := range map[Claim]stringList{
		ClaimAppEngine:       c.AppEngine,
		ClaimChannels:        c.Channels,
		ClaimHousekeeping:    c.Housekeeping,
		ClaimPairing:         c.Pairing,
		ClaimRealmManagement: c.RealmManagement,
	} {
		for _, s := range raw {
			if g, ok := compileGrant(s); ok {
				t.grants[claim] = append(t.grants[claim], g)
			}
		}
	}
	return t
}

// ExpiresAt returns the token expiry and whether one is set (`exp` is
// optional, upstream parity).
func (t *Token) ExpiresAt() (time.Time, bool) {
	if t.expiresAt == nil {
		return time.Time{}, false
	}
	return *t.expiresAt, true
}

// Authorizes reports whether the token grants `claim` for the given verb
// (HTTP method, or JOIN/WATCH on the stream socket) and authorization path
// (the request path relative to the realm base, e.g. "agent/devices").
// Multiple authorization strings within a claim are OR-ed; a token without
// the claim authorizes nothing on that surface.
func (t *Token) Authorizes(claim Claim, verb, authPath string) bool {
	for _, g := range t.grants[claim] {
		if g.verb.MatchString(verb) && g.path.MatchString(authPath) {
			return true
		}
	}
	return false
}
