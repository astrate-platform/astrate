package auth

import (
	"encoding/json"
	"testing"
)

// tokenWithClaims builds a Token straight from claim JSON (no signature):
// claim matching is independent of signature verification.
func tokenWithClaims(t *testing.T, claimsJSON string) *Token {
	t.Helper()
	c := &astarteClaims{}
	if err := json.Unmarshal([]byte(claimsJSON), c); err != nil {
		t.Fatalf("unmarshalling claims: %v", err)
	}
	return newToken(c)
}

// TestClaimMatchingParity is the upstream-parity table (docs/ROADMAP.md §4):
// the semantics of Astarte's GuardianAuthorizePath plug, including the cases
// astartectl emits, implicit anchoring, OR across strings, and the
// literal-concatenation anchoring quirks.
func TestClaimMatchingParity(t *testing.T) {
	cases := []struct {
		name   string
		grants []string // a_aea authorization strings
		verb   string
		path   string
		want   bool
	}{
		// astartectl default: catch-all.
		{"catch-all matches GET", []string{".*::.*"}, "GET", "devices/h4-Dx_RYTU-RbpDOTabhRg", true},
		{"catch-all matches POST deep path", []string{".*::.*"}, "POST", "devices/x/interfaces/com.example.Iface/path", true},
		{"catch-all matches empty path", []string{".*::.*"}, "GET", "", true},

		// Verb restriction.
		{"POST devices/.* allows POST", []string{"^POST$::^devices/.*$"}, "POST", "devices/x/interfaces", true},
		{"POST devices/.* denies GET", []string{"^POST$::^devices/.*$"}, "GET", "devices/x/interfaces", false},
		{"POST devices/.* denies DELETE", []string{"^POST$::^devices/.*$"}, "DELETE", "devices/x", false},

		// Implicit anchoring: bare `devices` must NOT match `devices/abc`.
		{"anchored literal exact match", []string{".*::devices"}, "GET", "devices", true},
		{"anchored literal no subpath", []string{".*::devices"}, "GET", "devices/abc", false},
		{"anchored literal no prefix", []string{".*::devices"}, "GET", "xdevices", false},
		{"suffixed wildcard matches subpath", []string{".*::devices/.*"}, "GET", "devices/abc", true},
		{"suffixed wildcard requires slash", []string{".*::devices/.*"}, "GET", "devices", false},
		{"verb anchoring", []string{"GET::.*"}, "GETX", "anything", false},

		// Multiple authorization strings are OR-ed.
		{"OR first matches", []string{"^GET$::^a$", "^POST$::^b$"}, "GET", "a", true},
		{"OR second matches", []string{"^GET$::^a$", "^POST$::^b$"}, "POST", "b", true},
		{"OR cross combination denied", []string{"^GET$::^a$", "^POST$::^b$"}, "POST", "a", false},

		// Upstream anchors by literal concatenation ("^" <> s <> "$"), so a
		// top-level alternation escapes the anchors: "^a|b$" is
		// (left-anchored a) OR (right-anchored b). Parity, not a bug.
		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},
		{"alternation anchoring quirk middle", []string{".*::devices|other"}, "GET", "XotherX", false},

		// Upstream splits on ":" into exactly three parts; the middle
		// (opts) field is ignored.
		{"opts field ignored", []string{"^GET$:ignored:^a$"}, "GET", "a", true},

		// Malformed strings never match but never poison the token.
		{"missing separator never matches", []string{"justonestring"}, "GET", "justonestring", false},
		{"single colon never matches", []string{"GET:devices"}, "GET", "devices", false},
		{"bad regex never matches", []string{"^GET$::^devices[$"}, "GET", "devices[", false},
		{"bad regex does not poison siblings", []string{"^GET$::^devices[$", ".*::.*"}, "GET", "anything", true},

		// Astarte Channels verbs are JOIN/WATCH, not HTTP methods.
		{"channels JOIN verb", []string{"JOIN::.*"}, "JOIN", "rooms/test", true},
		{"channels JOIN verb denies WATCH", []string{"JOIN::.*"}, "WATCH", "rooms/test", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			grants, err := json.Marshal(tc.grants)
			if err != nil {
				t.Fatal(err)
			}
			tok := tokenWithClaims(t, `{"a_aea": `+string(grants)+`}`)
			claim := ClaimAppEngine
			if tc.name[:8] == "channels" {
				tok = tokenWithClaims(t, `{"a_ch": `+string(grants)+`}`)
				claim = ClaimChannels
			}
			if got := tok.Authorizes(claim, tc.verb, tc.path); got != tc.want {
				t.Errorf("Authorizes(%v, %q, %q) with %v = %v, want %v",
					claim, tc.verb, tc.path, tc.grants, got, tc.want)
			}
		})
	}
}

func TestClaimsAreScoped(t *testing.T) {
	tok := tokenWithClaims(t, `{"a_pa": [".*::.*"]}`)

	if !tok.Authorizes(ClaimPairing, "POST", "agent/devices") {
		t.Error("a_pa catch-all should authorize the pairing surface")
	}
	for _, other := range []Claim{ClaimAppEngine, ClaimRealmManagement, ClaimHousekeeping, ClaimChannels} {
		if tok.Authorizes(other, "GET", "anything") {
			t.Errorf("a_pa-only token must not authorize %s", other)
		}
	}
}

func TestStringListLenientDecoding(t *testing.T) {
	// A bare string is tolerated as a single-element list.
	tok := tokenWithClaims(t, `{"a_rma": "^GET$::^interfaces$"}`)
	if !tok.Authorizes(ClaimRealmManagement, "GET", "interfaces") {
		t.Error("bare-string claim should behave as a one-element list")
	}

	// null claim: no grants.
	tok = tokenWithClaims(t, `{"a_rma": null}`)
	if tok.Authorizes(ClaimRealmManagement, "GET", "interfaces") {
		t.Error("null claim must not authorize anything")
	}

	// Non-string members are a decode error (token rejected upstream too).
	c := &astarteClaims{}
	if err := json.Unmarshal([]byte(`{"a_rma": [42]}`), c); err == nil {
		t.Error("numeric claim member should fail to decode")
	}
	if err := json.Unmarshal([]byte(`{"a_rma": 42}`), c); err == nil {
		t.Error("numeric claim should fail to decode")
	}
}

// TestAuthorizesChannelParity is the a_ch table recorded from upstream Astarte
// v1.2.0 in test/conformance/upstream/channels.json. Channels reads a_ch by a
// different rule than the REST claims: the verb field is matched by equality,
// not compiled as a regex, so the blanket grant that authorizes every REST
// surface authorizes nothing here.
//
// Each refusal is paired with an acceptance that differs only in the rule under
// test, so a matcher that refused everything could not pass the table.
func TestAuthorizesChannelParity(t *testing.T) {
	const room = "probe"

	cases := []struct {
		name   string
		grants []string // a_ch authorization strings
		verb   string
		path   string
		want   bool
		rule   string
	}{
		{
			name: "blanket .*::.* authorizes no join", grants: []string{".*::.*"},
			verb: VerbJoin, path: room, want: false,
			rule: `".*" is not the literal string "JOIN"`,
		},
		{
			name: "blanket .* authorizes no join", grants: []string{".*"},
			verb: VerbJoin, path: room, want: false,
			rule: "an entry without the three-field form is discarded",
		},
		{
			name: "JOIN::.* authorizes the join", grants: []string{"JOIN::.*"},
			verb: VerbJoin, path: room, want: true,
			rule: "the paired acceptance: an explicit JOIN verb is what upstream requires",
		},
		{
			name: "JOIN::<this room> authorizes the join", grants: []string{"JOIN::probe"},
			verb: VerbJoin, path: room, want: true,
			rule: "an exact-room grant is honoured",
		},
		{
			name: "JOIN::<other room> does not", grants: []string{"JOIN::other"},
			verb: VerbJoin, path: room, want: false,
			rule: "the room name is really matched, not merely the verb",
		},
		{
			name: "WATCH::.* authorizes a watch", grants: []string{"JOIN::.*", "WATCH::.*"},
			verb: VerbWatch, path: "device/iface", want: true,
			rule: "the watch bucket is selected by its own verb",
		},
		{
			name: "a JOIN grant does not authorize a watch", grants: []string{"JOIN::.*"},
			verb: VerbWatch, path: "device/iface", want: false,
			rule: "the two verbs are separate buckets, not a shared path list",
		},
		{
			name: "a WATCH grant does not authorize a join", grants: []string{"WATCH::.*"},
			verb: VerbJoin, path: room, want: false,
			rule: "the two verbs are separate buckets, not a shared path list",
		},
		{
			name: "an unknown verb is discarded", grants: []string{"SUBSCRIBE::.*"},
			verb: VerbJoin, path: room, want: false,
			rule: "upstream keeps only JOIN and WATCH entries",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.grants)
			if err != nil {
				t.Fatalf("marshalling grants: %v", err)
			}
			tok := tokenWithClaims(t, `{"a_ch": `+string(b)+`}`)
			if got := tok.AuthorizesChannel(tc.verb, tc.path); got != tc.want {
				t.Errorf("AuthorizesChannel(%q, %q) with a_ch=%v = %v, want %v — %s",
					tc.verb, tc.path, tc.grants, got, tc.want, tc.rule)
			}
		})
	}
}

// TestRESTVerbStaysARegex is the invariant the change above must not break.
// Upstream's REST plug compiles the verb field into a regex and matches it
// against the HTTP method, so ".*::.*" grants every REST surface and "GET|POST"
// selects two methods. Only Channels departs from this.
func TestRESTVerbStaysARegex(t *testing.T) {
	tok := tokenWithClaims(t, `{"a_aea": [".*::.*"]}`)
	if !tok.Authorizes(ClaimAppEngine, "GET", "devices") {
		t.Error(`".*::.*" must still authorize a REST GET — the verb field is a regex there`)
	}
	if tok.AuthorizesChannel(VerbJoin, "devices") {
		t.Error("an a_aea grant must not authorize the Channels surface at all")
	}

	alt := tokenWithClaims(t, `{"a_aea": ["GET|POST::devices"]}`)
	if !alt.Authorizes(ClaimAppEngine, "POST", "devices") {
		t.Error(`"GET|POST" must still match POST — the verb field is a regex there`)
	}
}
