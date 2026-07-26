package main

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// The realm-creation paths are the only part of the harness whose behaviour
// differs between the two platforms it drives, and both differences were found
// the hard way against a live upstream (2026-07-26). They are pinned here
// because neither can be reached without a running deployment otherwise.

func testKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func TestCreateRealmRekeysExistingRealm(t *testing.T) {
	// Upstream Astarte 1.2.0 reports an existing realm as 422 with
	// error_name "existing_realm"; Astrate reports 409. Both must lead to
	// the same PATCH, or provision is not re-runnable on one of them.
	for _, tc := range []struct {
		name   string
		status int
		body   string
	}{
		{"upstream 422", http.StatusUnprocessableEntity, `{"errors":{"error_name":["existing_realm"]}}`},
		{"astrate 409", http.StatusConflict, `{"errors":{"detail":"Realm already exists"}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var patched string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch {
				case r.Method == http.MethodPost && r.URL.Path == "/housekeeping/v1/realms":
					w.WriteHeader(tc.status)
					_, _ = w.Write([]byte(tc.body))
				case r.Method == http.MethodPatch && r.URL.Path == "/housekeeping/v1/realms/bench":
					patched = r.URL.Path
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"data":{}}`))
				default:
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusTeapot)
				}
			}))
			defer srv.Close()

			ep := deriveEndpoints(srv.URL)
			if err := newClient(5*time.Second).createRealm(ep, testKey(t), "bench", "PUB"); err != nil {
				t.Fatalf("createRealm: %v", err)
			}
			if patched == "" {
				t.Fatal("existing realm was not re-keyed: no PATCH reached the server")
			}
		})
	}
}

func TestCreateRealmWithoutRekeyRouteExplainsItself(t *testing.T) {
	// A deployment with no realm-update route (Astrate) cannot be handed the
	// signing key this run generated. The failure must name that, not surface
	// later as an unexplained 403 on the first realm-management call.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"errors":{"detail":"Realm already exists"}}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"errors":{"detail":"Not found"}}`))
	}))
	defer srv.Close()

	err := newClient(5*time.Second).createRealm(deriveEndpoints(srv.URL), testKey(t), "bench", "PUB")
	if err == nil {
		t.Fatal("expected an error when the realm exists and cannot be re-keyed")
	}
	if !strings.Contains(err.Error(), "cannot be re-keyed") {
		t.Fatalf("error does not explain the cause: %v", err)
	}
}

func TestAwaitRealmKeyToleratesPropagationDelay(t *testing.T) {
	// Realm creation upstream is asynchronous: realm-management answers 403
	// for a couple of seconds after housekeeping returns 201. A 403 must be
	// retried, and anything else must fail immediately.
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"errors":{"detail":"Forbidden"}}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer srv.Close()

	c := newClient(5 * time.Second)
	if err := c.awaitRealmKey(deriveEndpoints(srv.URL), testKey(t), "bench", time.Minute); err != nil {
		t.Fatalf("awaitRealmKey: %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected to retry through two 403s, got %d calls", calls)
	}
}

func TestAwaitRealmKeyDoesNotRetryOtherFailures(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"errors":{"detail":"Unauthorized"}}`))
	}))
	defer srv.Close()

	err := newClient(5*time.Second).awaitRealmKey(deriveEndpoints(srv.URL), testKey(t), "bench", time.Minute)
	if err == nil {
		t.Fatal("expected a 401 to fail immediately")
	}
	if calls != 1 {
		t.Fatalf("a non-403 must not be retried, got %d calls", calls)
	}
}
