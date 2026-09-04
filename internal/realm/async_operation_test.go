//go:build integration

package realm

import (
	"fmt"
	"net/http"
	"testing"
)

// TestRealmManagementAsyncOperationParam pins deviation 17
// (docs/COMPATIBILITY.md): upstream 1.4 runs interface install/update/delete
// and policy delete asynchronously unless the caller opts out with
// `?async_operation=false`. Astrate is always synchronous, so the parameter is
// accepted and ignored on either value — sending it must never draw a 4xx, and
// the effect must be visible on the very next read.
func TestRealmManagementAsyncOperationParam(t *testing.T) {
	r := newRig(t)

	for i, q := range []string{"?async_operation=false", "?async_operation=true"} {
		name := fmt.Sprintf("com.ex.M7a.Async%d", i)
		iface := func(major, minor int, extra string) string {
			return fmt.Sprintf(`{"interface_name":%q,"version_major":%d,"version_minor":%d,`+
				`"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}%s]}`,
				name, major, minor, extra)
		}

		// Install, then an additive minor upgrade, both with the parameter.
		if rec := r.req(t, http.MethodPost, "/interfaces"+q, iface(1, 0, ""), r.rmaToken); rec.Code != http.StatusCreated {
			t.Fatalf("install interface %s: got %d, want 201 (%s)", q, rec.Code, rec.Body)
		}
		if rec := r.req(t, http.MethodGet, "/interfaces/"+name+"/1", "", r.rmaToken); rec.Code != http.StatusOK {
			t.Errorf("get interface after install %s: got %d, want 200", q, rec.Code)
		}
		upgrade := iface(1, 1, `,{"endpoint":"/count","type":"integer"}`)
		if rec := r.req(t, http.MethodPut, "/interfaces/"+name+"/1"+q, upgrade, r.rmaToken); rec.Code != http.StatusNoContent {
			t.Fatalf("update interface %s: got %d, want 204 (%s)", q, rec.Code, rec.Body)
		}

		// Delete only accepts a draft (major 0), so that leg gets its own.
		draftName := fmt.Sprintf("com.ex.M7a.AsyncDraft%d", i)
		draft := fmt.Sprintf(`{"interface_name":%q,"version_major":0,"version_minor":1,`+
			`"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`, draftName)
		if rec := r.req(t, http.MethodPost, "/interfaces", draft, r.rmaToken); rec.Code != http.StatusCreated {
			t.Fatalf("install draft %s: got %d, want 201 (%s)", q, rec.Code, rec.Body)
		}
		if rec := r.req(t, http.MethodDelete, "/interfaces/"+draftName+"/0"+q, "", r.rmaToken); rec.Code != http.StatusNoContent {
			t.Fatalf("delete draft %s: got %d, want 204 (%s)", q, rec.Code, rec.Body)
		}
		if rec := r.req(t, http.MethodGet, "/interfaces/"+draftName+"/0", "", r.rmaToken); rec.Code != http.StatusNotFound {
			t.Errorf("get draft after delete %s: got %d, want 404", q, rec.Code)
		}

		// Policy delete carries the same parameter upstream.
		policyName := fmt.Sprintf("async%d", i)
		policy := fmt.Sprintf(`{"name":%q,"error_handlers":[{"on":"server_error","strategy":"retry"}],`+
			`"maximum_capacity":100,"retry_times":3}`, policyName)
		if rec := r.req(t, http.MethodPost, "/policies", policy, r.rmaToken); rec.Code != http.StatusCreated {
			t.Fatalf("create policy %s: got %d, want 201 (%s)", q, rec.Code, rec.Body)
		}
		if rec := r.req(t, http.MethodDelete, "/policies/"+policyName+q, "", r.rmaToken); rec.Code != http.StatusNoContent {
			t.Fatalf("delete policy %s: got %d, want 204 (%s)", q, rec.Code, rec.Body)
		}
		if rec := r.req(t, http.MethodGet, "/policies/"+policyName, "", r.rmaToken); rec.Code != http.StatusNotFound {
			t.Errorf("get policy after delete %s: got %d, want 404", q, rec.Code)
		}
	}
}
