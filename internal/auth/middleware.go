package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

// KeySource supplies per-realm key material. *store.Store satisfies it; the
// middleware only reads JWTPublicKeysPEM from the returned realm.
type KeySource interface {
	GetRealmByName(ctx context.Context, name string) (*store.Realm, error)
}

// tokenContextKey carries the verified *Token in the request context.
type tokenContextKey struct{}

// TokenFromContext returns the verified token stored by the middleware, if
// the request passed through Require*.
func TokenFromContext(ctx context.Context) (*Token, bool) {
	tok, ok := ctx.Value(tokenContextKey{}).(*Token)
	return tok, ok
}

// Middleware authenticates and authorizes REST requests with realm JWTs
// (docs/DESIGN.md §4.2). Status mapping is upstream parity: missing or
// unverifiable token → 401, verified token whose claims do not authorize the
// request → 403, both with the canonical envelopes.
type Middleware struct {
	keys  KeySource
	cache *Cache
}

// NewMiddleware builds a Middleware over the given key source with a
// DefaultCacheSize token cache.
func NewMiddleware(keys KeySource) *Middleware {
	return &Middleware{keys: keys, cache: NewCache(DefaultCacheSize)}
}

// RequireRealm guards a realm-scoped route (path pattern must carry a
// {realm} segment): it resolves the realm's JWT public keys, verifies the
// bearer token, and matches the claim's authorization strings against the
// method and the path relative to the realm base.
func (m *Middleware) RequireRealm(claim Claim) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			realm := r.PathValue("realm")
			if realm == "" {
				_ = astarteapi.WriteUnauthorized(w)
				return
			}

			row, err := m.keys.GetRealmByName(r.Context(), realm)
			switch {
			case errors.Is(err, store.ErrNotFound):
				// An unknown realm has no keys: unauthenticated, not 404
				// (no existence oracle on auth failures).
				_ = astarteapi.WriteUnauthorized(w)
				return
			case err != nil:
				_ = astarteapi.WriteInternalServerError(w)
				return
			}

			m.authorize(w, r, next, claim, row.JWTPublicKeysPEM, realm)
		})
	}
}

// RequireStatic guards an instance-level route (Housekeeping) with a fixed
// key set instead of per-realm keys. The authorization path is the request
// path relative to the service base (the segment after "v1").
func (m *Middleware) RequireStatic(claim Claim, keysPEM []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.authorize(w, r, next, claim, keysPEM, "v1")
		})
	}
}

// authorize runs the shared bearer-extract → verify → claim-match pipeline.
// base is the path segment after which the authorization path starts (the
// realm name, or "v1" for instance-level routes).
func (m *Middleware) authorize(w http.ResponseWriter, r *http.Request, next http.Handler, claim Claim, keysPEM []string, base string) {
	tokenString, ok := bearerToken(r)
	if !ok {
		_ = astarteapi.WriteUnauthorized(w)
		return
	}

	tok, err := m.cache.Verify(tokenString, keysPEM)
	if err != nil {
		_ = astarteapi.WriteUnauthorized(w)
		return
	}

	authPath, ok := RelativePath(r.URL.Path, base)
	if !ok {
		// Upstream parity: a path the authorizer cannot anchor is an
		// authorization failure (403), not an authentication one.
		_ = astarteapi.WriteForbidden(w)
		return
	}
	if !tok.Authorizes(claim, r.Method, authPath) {
		_ = astarteapi.WriteForbidden(w)
		return
	}

	ctx := context.WithValue(r.Context(), tokenContextKey{}, tok)
	next.ServeHTTP(w, r.WithContext(ctx))
}

// RelativePath computes the authorization path with upstream parity (Astarte's
// GuardianAuthorizePath plug): split the URL path into segments, drop
// everything up to and including the first segment equal to base, and join
// the rest with "/". ok is false when base does not appear in the path.
//
// Example: RelativePath("/pairing/v1/test/agent/devices", "test") returns
// ("agent/devices", true).
func RelativePath(urlPath, base string) (string, bool) {
	segments := strings.Split(strings.Trim(urlPath, "/"), "/")
	for i, s := range segments {
		if s == base {
			return strings.Join(segments[i+1:], "/"), true
		}
	}
	return "", false
}

// bearerToken extracts the credential from the Authorization header.
// Upstream parity (~r/bearer\:?\s+(.*)$/i): the scheme is case-insensitive
// "Bearer", optionally followed by a colon, then whitespace and the token.
func bearerToken(r *http.Request) (string, bool) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", false
	}
	scheme, rest, found := strings.Cut(header, " ")
	if !found {
		return "", false
	}
	scheme = strings.TrimSuffix(scheme, ":")
	if !strings.EqualFold(scheme, "bearer") {
		return "", false
	}
	token := strings.TrimSpace(rest)
	if token == "" {
		return "", false
	}
	return token, true
}
