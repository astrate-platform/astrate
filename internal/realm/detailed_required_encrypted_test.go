//go:build !integration

package realm

import "testing"

// TestRenderDetailedInterfaceRequiredEncrypted pins that the detailed
// listing round-trips the required/encrypted mapping flags for datastreams
// (issue: the renderer used to drop them even when the stored definition
// declares one). Properties reject both flags at parse, so the datastream
// branch is the only one that can carry them.
func TestRenderDetailedInterfaceRequiredEncrypted(t *testing.T) {
	def := []byte(`{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":0,` +
		`"type":"datastream","ownership":"device","aggregation":"individual",` +
		`"mappings":[{"endpoint":"/value","type":"double","required":true,"encrypted":true}]}`)
	got, err := renderDetailedInterface(def)
	if err != nil {
		t.Fatalf("renderDetailedInterface: %v", err)
	}
	want := `{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":0,` +
		`"type":"datastream","ownership":"device","aggregation":"individual",` +
		`"mappings":[{"endpoint":"/value","type":"double","reliability":"unreliable",` +
		`"retention":"discard","expiry":0,"explicit_timestamp":false,` +
		`"required":true,"encrypted":true,` +
		`"database_retention_policy":"no_ttl"}]}`
	if string(got) != want {
		t.Errorf("rendered =\n%s\nwant\n%s", got, want)
	}
}
