package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/realm"
)

// TestRealmScopedVersionReportsAPICompatLevel proves the realm-scoped version
// endpoints serve the emulated upstream API level the Dashboard feature-gates
// its UI on (COMPATIBILITY.md deviation 10, issue #77), not Astrate's build
// version. The pairing route is public, so it is reachable without a token.
func TestRealmScopedVersionReportsAPICompatLevel(t *testing.T) {
	mux := http.NewServeMux()
	mountRealmVersion(mux, auth.NewMiddleware(nil))

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/pairing/v1/somerealm/version", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("pairing version status = %d, want 200", rec.Code)
	}
	want := `{"data":"` + realm.APICompatVersion + `"}`
	if got := rec.Body.String(); got != want {
		t.Errorf("pairing version body = %s, want %s", got, want)
	}
}
