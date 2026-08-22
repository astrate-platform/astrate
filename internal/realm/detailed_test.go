//go:build integration

package realm

import (
	"net/http"
	"testing"
)

// TestRealmManagementDetailedListing pins issue #66 on the wire:
// ?detailed=true serves one fully materialised 1.4-style document per
// interface (every datastream default filled in), while the plain listing —
// and any other parameter value — stays byte-identical to upstream 1.2's
// names-only response.
func TestRealmManagementDetailedListing(t *testing.T) {
	r := newRig(t)

	if rec := r.req(t, http.MethodPost, "/interfaces", ifaceV1, r.rmaToken); rec.Code != http.StatusCreated {
		t.Fatalf("install minimal datastream interface: got %d, want 201 (%s)", rec.Code, rec.Body)
	}

	want := `{"data":[{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":0,` +
		`"type":"datastream","ownership":"device","aggregation":"individual",` +
		`"mappings":[{"endpoint":"/value","type":"double","reliability":"unreliable",` +
		`"retention":"discard","expiry":0,"explicit_timestamp":false,` +
		`"database_retention_policy":"no_ttl"}]}]}`
	rec := r.req(t, http.MethodGet, "/interfaces?detailed=true", "", r.rmaToken)
	if rec.Code != http.StatusOK {
		t.Fatalf("detailed listing: got %d, want 200 (%s)", rec.Code, rec.Body)
	}
	if got := rec.Body.String(); got != want {
		t.Errorf("detailed body =\n%s\nwant\n%s", got, want)
	}

	var names []string
	decodeData(t, r.req(t, http.MethodGet, "/interfaces", "", r.rmaToken), &names)
	if len(names) != 1 || names[0] != rmIface {
		t.Errorf("plain listing = %v, want bare names [%s]", names, rmIface)
	}
	decodeData(t, r.req(t, http.MethodGet, "/interfaces?detailed=false", "", r.rmaToken), &names)
	if len(names) != 1 || names[0] != rmIface {
		t.Errorf("detailed=false listing = %v, want bare names [%s]", names, rmIface)
	}
}

// TestRenderDetailedInterfaceProperties covers the properties branch of the
// renderer directly: no datastream defaults leak in, allow_unset is
// materialised.
func TestRenderDetailedInterfaceProperties(t *testing.T) {
	def := []byte(`{"interface_name":"com.ex.M7a.Props","version_major":0,"version_minor":2,` +
		`"type":"properties","ownership":"server","description":"props",` +
		`"mappings":[{"endpoint":"/state","type":"boolean"}]}`)
	got, err := renderDetailedInterface(def)
	if err != nil {
		t.Fatalf("renderDetailedInterface: %v", err)
	}
	want := `{"interface_name":"com.ex.M7a.Props","version_major":0,"version_minor":2,` +
		`"type":"properties","ownership":"server","aggregation":"individual","description":"props",` +
		`"mappings":[{"endpoint":"/state","type":"boolean","allow_unset":false}]}`
	if string(got) != want {
		t.Errorf("rendered =\n%s\nwant\n%s", got, want)
	}
}
