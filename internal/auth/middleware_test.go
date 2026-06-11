package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
)

// fakeKeySource is an in-memory KeySource.
type fakeKeySource map[string]*store.Realm

func (f fakeKeySource) GetRealmByName(_ context.Context, name string) (*store.Realm, error) {
	r, ok := f[name]
	if !ok {
		return nil, fmt.Errorf("%w: realm %q", store.ErrNotFound, name)
	}
	return r, nil
}

// newTestMux wires a realm-guarded pairing route, a realm-guarded appengine
// route, and a static-key housekeeping route, mirroring how the real
// surfaces mount the middleware.
func newTestMux(t *testing.T, m *Middleware, hkKeys []string) *http.ServeMux {
	t.Helper()
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := TokenFromContext(r.Context()); !ok {
			t.Error("verified token missing from request context")
		}
		w.WriteHeader(http.StatusOK)
	})

	mux := http.NewServeMux()
	mux.Handle("POST /pairing/v1/{realm}/agent/devices",
		m.RequireRealm(ClaimPairing)(okHandler))
	mux.Handle("GET /appengine/v1/{realm}/devices/{deviceID}",
		m.RequireRealm(ClaimAppEngine)(okHandler))
	mux.Handle("GET /housekeeping/v1/realms",
		m.RequireStatic(ClaimHousekeeping, hkKeys)(okHandler))
	return mux
}

// do performs a request against the mux and returns the recorder.
func do(mux *http.ServeMux, method, path, token string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, nil)
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w
}

func TestMiddleware(t *testing.T) {
	tk := keys(t)
	pubPEM := publicPEM(t, &tk.rsaKey.PublicKey)
	hkPEM := publicPEM(t, &tk.ecKey.PublicKey)

	ks := fakeKeySource{
		"test": &store.Realm{ID: 1, Name: "test", JWTPublicKeysPEM: []string{pubPEM}},
	}
	m := NewMiddleware(ks)
	mux := newTestMux(t, m, []string{hkPEM})

	pairingToken := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{
		"a_pa": []string{".*::.*"},
		"exp":  time.Now().Add(time.Hour).Unix(),
	})
	narrowToken := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{
		"a_aea": []string{"^GET$::^devices/.*$"},
	})
	hkToken := signToken(t, jwt.SigningMethodES256, tk.ecKey, jwt.MapClaims{
		"a_ha": []string{".*::.*"},
	})

	t.Run("MissingToken401", func(t *testing.T) {
		w := do(mux, "POST", "/pairing/v1/test/agent/devices", "")
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status: got %d, want 401", w.Code)
		}
		testutil.Golden(t, "envelope_401.json", w.Body.Bytes())
	})

	t.Run("GarbageToken401", func(t *testing.T) {
		w := do(mux, "POST", "/pairing/v1/test/agent/devices", "not.a.token")
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status: got %d, want 401", w.Code)
		}
		testutil.Golden(t, "envelope_401.json", w.Body.Bytes())
	})

	t.Run("UnknownRealm401", func(t *testing.T) {
		w := do(mux, "POST", "/pairing/v1/nosuchrealm/agent/devices", pairingToken)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status: got %d, want 401", w.Code)
		}
		testutil.Golden(t, "envelope_401.json", w.Body.Bytes())
	})

	t.Run("WrongClaim403", func(t *testing.T) {
		// narrowToken has a_aea but no a_pa: authenticated, not authorized.
		w := do(mux, "POST", "/pairing/v1/test/agent/devices", narrowToken)
		if w.Code != http.StatusForbidden {
			t.Fatalf("status: got %d, want 403", w.Code)
		}
		testutil.Golden(t, "envelope_403.json", w.Body.Bytes())
	})

	t.Run("ClaimPathMismatch403", func(t *testing.T) {
		// a_aea grants ^GET$::^devices/.*$ — a bare 'devices/<id>' path
		// matches, but only with method GET; here we check that the *path*
		// regex is matched against the realm-relative path by using a
		// non-matching deeper route.
		mux2 := http.NewServeMux()
		mux2.Handle("GET /appengine/v1/{realm}/groups",
			m.RequireRealm(ClaimAppEngine)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})))
		w := do(mux2, "GET", "/appengine/v1/test/groups", narrowToken)
		if w.Code != http.StatusForbidden {
			t.Fatalf("status: got %d, want 403", w.Code)
		}
		testutil.Golden(t, "envelope_403.json", w.Body.Bytes())
	})

	t.Run("Authorized200", func(t *testing.T) {
		w := do(mux, "POST", "/pairing/v1/test/agent/devices", pairingToken)
		if w.Code != http.StatusOK {
			t.Fatalf("status: got %d, want 200 (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("NarrowClaimAuthorized200", func(t *testing.T) {
		w := do(mux, "GET", "/appengine/v1/test/devices/h4-Dx_RYTU-RbpDOTabhRg", narrowToken)
		if w.Code != http.StatusOK {
			t.Fatalf("status: got %d, want 200 (body %s)", w.Code, w.Body.String())
		}
	})

	t.Run("HousekeepingStaticKeys", func(t *testing.T) {
		w := do(mux, "GET", "/housekeeping/v1/realms", hkToken)
		if w.Code != http.StatusOK {
			t.Fatalf("status with hk token: got %d, want 200 (body %s)", w.Code, w.Body.String())
		}
		// The realm token is signed with a key outside the housekeeping
		// set: 401.
		w = do(mux, "GET", "/housekeeping/v1/realms", pairingToken)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("status with realm token: got %d, want 401", w.Code)
		}
		testutil.Golden(t, "envelope_401.json", w.Body.Bytes())
	})

	t.Run("BearerSchemeVariants", func(t *testing.T) {
		for _, header := range []string{
			"Bearer " + pairingToken,
			"bearer " + pairingToken,
			"Bearer: " + pairingToken,
		} {
			r := httptest.NewRequest("POST", "/pairing/v1/test/agent/devices", nil)
			r.Header.Set("Authorization", header)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			if w.Code != http.StatusOK {
				t.Errorf("header %q: got %d, want 200", header, w.Code)
			}
		}
		r := httptest.NewRequest("POST", "/pairing/v1/test/agent/devices", nil)
		r.Header.Set("Authorization", "Basic "+pairingToken)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Basic scheme: got %d, want 401", w.Code)
		}
	})
}

func TestRelativePath(t *testing.T) {
	cases := []struct {
		urlPath string
		base    string
		want    string
		ok      bool
	}{
		{"/pairing/v1/test/agent/devices", "test", "agent/devices", true},
		{"/appengine/v1/test/devices/x/interfaces/i", "test", "devices/x/interfaces/i", true},
		{"/realmmanagement/v1/test", "test", "", true},
		{"/housekeeping/v1/realms", "v1", "realms", true},
		{"/housekeeping/v1/realms/test", "v1", "realms/test", true},
		{"/pairing/v1/other/agent/devices", "test", "", false},
		// First occurrence of the base segment wins (upstream drop_while).
		{"/pairing/v1/test/agent/test/x", "test", "agent/test/x", true},
	}
	for _, tc := range cases {
		got, ok := RelativePath(tc.urlPath, tc.base)
		if got != tc.want || ok != tc.ok {
			t.Errorf("RelativePath(%q, %q) = (%q, %v), want (%q, %v)",
				tc.urlPath, tc.base, got, ok, tc.want, tc.ok)
		}
	}
}
