//go:build integration

package housekeeping

import (
	"net/http"
	"testing"
)

// TestHousekeepingAsyncOperationParam pins deviation 17
// (docs/COMPATIBILITY.md): upstream 1.4 runs realm create and delete
// asynchronously unless the caller opts out with `?async_operation=false`.
// Astrate is always synchronous, so the parameter is accepted and ignored on
// both values — a client that sends it must never see a 4xx for it, and the
// response must be identical to the parameterless call.
func TestHousekeepingAsyncOperationParam(t *testing.T) {
	r := newHKRig(t)

	for _, q := range []string{"?async_operation=false", "?async_operation=true"} {
		realmName := "as" + randSuffix(t)
		createBody := `{"realm_name":` + jsonStr(realmName) + `,"jwt_public_key_pem":` + jsonStr(r.realmKey) + `}`

		if rec := r.req(t, http.MethodPost, "/realms"+q, createBody, r.haToken); rec.Code != http.StatusCreated {
			t.Fatalf("create realm %s: got %d, want 201 (%s)", q, rec.Code, rec.Body)
		}
		// Synchronous means done on return: the realm reads back immediately.
		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusOK {
			t.Errorf("get realm after create %s: got %d, want 200", q, rec.Code)
		}
		if rec := r.req(t, http.MethodDelete, "/realms/"+realmName+q, "", r.haToken); rec.Code != http.StatusNoContent {
			t.Fatalf("delete realm %s: got %d, want 204 (%s)", q, rec.Code, rec.Body)
		}
		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusNotFound {
			t.Errorf("get realm after delete %s: got %d, want 404", q, rec.Code)
		}
	}
}
