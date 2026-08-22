package interfaceschema_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// legacyProps is the minimal properties interface in pure legacy-alias form,
// mirroring the shape upstream 1.2.0 accepted during the #61 probes.
const legacyProps = `{"interface_name":"com.astrate.test.Legacy","version_major":0,"version_minor":1,` +
	`"type":"properties","quality":"device","aggregate":false,` +
	`"mappings":[{"path":"/value","type":"string"}]}`

func mustParseViolationDoc(t *testing.T, data string) *interfaceschema.Interface {
	t.Helper()
	iface, err := interfaceschema.ParseInterface([]byte(data))
	if err != nil {
		t.Fatalf("ParseInterface: %v", err)
	}
	return iface
}

// violationsOf fails the test unless err carries a *ViolationsError.
func violationsOf(t *testing.T, err error) *interfaceschema.ViolationsError {
	t.Helper()
	if err == nil {
		t.Fatal("ParseInterface succeeded, want rejection")
	}
	var ve *interfaceschema.ViolationsError
	if !errors.As(err, &ve) {
		t.Fatalf("error %v does not carry *ViolationsError", err)
	}
	return ve
}

// findViolation returns the violation for (field, index), failing otherwise.
func findViolation(t *testing.T, ve *interfaceschema.ViolationsError, field string, index int) *interfaceschema.Violation {
	t.Helper()
	for i := range ve.Violations {
		v := &ve.Violations[i]
		if v.Field == field && v.MappingIndex == index {
			return v
		}
	}
	t.Fatalf("no violation for field %q at mapping %d in %+v", field, index, ve.Violations)
	return nil
}

func TestLegacyAliasesAccepted(t *testing.T) {
	t.Run("quality_aggregate_path_alone", func(t *testing.T) {
		iface := mustParseViolationDoc(t, legacyProps)
		if iface.Ownership != interfaceschema.OwnershipDevice {
			t.Errorf("Ownership = %v, want device", iface.Ownership)
		}
		if iface.Aggregation != interfaceschema.AggregationIndividual {
			t.Errorf("Aggregation = %v, want individual", iface.Aggregation)
		}
		if len(iface.Mappings) != 1 || iface.Mappings[0].Endpoint != "/value" {
			t.Errorf("Mappings = %+v, want endpoint /value", iface.Mappings)
		}
	})

	t.Run("aggregate_true_object", func(t *testing.T) {
		iface := mustParseViolationDoc(t, `{"interface_name":"com.astrate.test.LegacyObj","version_major":1,`+
			`"version_minor":0,"type":"datastream","ownership":"device","aggregate":true,`+
			`"mappings":[{"endpoint":"/base/k1","type":"double"},{"endpoint":"/base/k2","type":"double"}]}`)
		if iface.Aggregation != interfaceschema.AggregationObject {
			t.Errorf("Aggregation = %v, want object", iface.Aggregation)
		}
	})

	t.Run("endpoint_wins_over_path", func(t *testing.T) {
		iface := mustParseViolationDoc(t, `{"interface_name":"com.astrate.test.BothPaths","version_major":1,`+
			`"version_minor":0,"type":"datastream","ownership":"device",`+
			`"mappings":[{"endpoint":"/a","path":"/b","type":"double"}]}`)
		if iface.Mappings[0].Endpoint != "/a" {
			t.Errorf("Endpoint = %q, want /a (endpoint wins)", iface.Mappings[0].Endpoint)
		}
	})
}

func TestLegacyAliasConflictsRejected(t *testing.T) {
	cases := []struct {
		name  string
		doc   string
		field string
	}{
		{
			name: "ownership_and_quality",
			doc: `{"interface_name":"com.astrate.test.Conflict","version_major":1,"version_minor":0,` +
				`"type":"datastream","ownership":"device","quality":"server",` +
				`"mappings":[{"endpoint":"/v","type":"double"}]}`,
			field: "ownership",
		},
		{
			name: "aggregation_and_aggregate",
			doc: `{"interface_name":"com.astrate.test.Conflict","version_major":1,"version_minor":0,` +
				`"type":"datastream","ownership":"device","aggregation":"individual","aggregate":true,` +
				`"mappings":[{"endpoint":"/v","type":"double"}]}`,
			field: "aggregation",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := interfaceschema.ParseInterface([]byte(tc.doc))
			ve := violationsOf(t, err)
			v := findViolation(t, ve, tc.field, -1)
			if len(v.Messages) == 0 || !strings.Contains(v.Messages[0], "mutually exclusive") {
				t.Errorf("Messages = %v, want a mutual-exclusion message", v.Messages)
			}
			if !errors.Is(err, interfaceschema.ErrInvalid) {
				t.Errorf("error %v does not wrap ErrInvalid", err)
			}
		})
	}

	t.Run("unknown_field_still_rejected_not_violation", func(t *testing.T) {
		_, err := interfaceschema.ParseInterface([]byte(`{"interface_name":"com.astrate.test.Bogus",` +
			`"version_major":1,"version_minor":0,"type":"datastream","ownership":"device","bogus_field":1,` +
			`"mappings":[{"endpoint":"/v","type":"double"}]}`))
		if err == nil {
			t.Fatal("made-up field accepted, want rejection")
		}
		if !errors.Is(err, interfaceschema.ErrInvalid) {
			t.Errorf("error %v does not wrap ErrInvalid", err)
		}
		var ve *interfaceschema.ViolationsError
		if errors.As(err, &ve) {
			t.Errorf("unknown-field decode produced %+v, want a plain error", ve)
		}
	})
}

// TestCanonicalBytes pins the ParseInterfaceCanonical contract that makes
// stored documents render canonically after an alias install (#61).
func TestCanonicalBytes(t *testing.T) {
	iface, canon, err := interfaceschema.ParseInterfaceCanonical([]byte(legacyProps))
	if err != nil {
		t.Fatalf("ParseInterfaceCanonical: %v", err)
	}
	if canon == nil {
		t.Fatal("canonical bytes = nil for an alias document, want re-encoding")
	}
	if strings.Contains(string(canon), `"quality"`) ||
		strings.Contains(string(canon), `"aggregate"`) ||
		strings.Contains(string(canon), `"path"`) {
		t.Errorf("canonical bytes still contain aliases: %s", canon)
	}
	want := `{"interface_name":"com.astrate.test.Legacy","version_major":0,"version_minor":1,` +
		`"type":"properties","ownership":"device","aggregation":"individual",` +
		`"mappings":[{"endpoint":"/value","type":"string"}]}`
	if string(canon) != want {
		t.Errorf("canonical =\n%s\nwant\n%s", canon, want)
	}
	rt, canon2, err := interfaceschema.ParseInterfaceCanonical(canon)
	if err != nil {
		t.Fatalf("re-parsing canonical bytes: %v", err)
	}
	if rt.Ownership != iface.Ownership || canon2 != nil {
		t.Errorf("canonical bytes did not round-trip cleanly: canon2 = %s", canon2)
	}

	canonicalDoc := `{"interface_name":"com.astrate.test.Canonical","version_major":1,"version_minor":0,` +
		`"type":"datastream","ownership":"device",` +
		`"mappings":[{"endpoint":"/value","type":"double"}]}`
	if _, canon3, err := interfaceschema.ParseInterfaceCanonical([]byte(canonicalDoc)); err != nil || canon3 != nil {
		t.Errorf("alias-free document: canon = %s, err = %v, want nil/nil", canon3, err)
	}
}

// TestDescriptionDocLengthBounds walks the probe-verified bounds (issue #61)
// with each rejection paired with its at-the-bound acceptance twin.
func TestDescriptionDocLengthBounds(t *testing.T) {
	dsIface := func(desc, doc string) string {
		body := `{"interface_name":"com.astrate.test.Bounds","version_major":1,"version_minor":0,` +
			`"type":"datastream","ownership":"device",`
		if desc != "" {
			body += `"description":"` + desc + `",`
		}
		if doc != "" {
			body += `"doc":"` + doc + `",`
		}
		return body + `"mappings":[{"endpoint":"/v","type":"double"}]}`
	}

	cases := []struct {
		name  string
		desc  string
		doc   string
		field string // "" = accepted twin
	}{
		{"description_1000_ok", strings.Repeat("x", 1000), "", ""},
		{"description_1001_rejected", strings.Repeat("x", 1001), "", "description"},
		{"doc_100000_ok", "", strings.Repeat("x", 100000), ""},
		{"doc_100001_rejected", "", strings.Repeat("x", 100001), "doc"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := interfaceschema.ParseInterface([]byte(dsIface(tc.desc, tc.doc)))
			if tc.field == "" {
				if err != nil {
					t.Fatalf("accepted twin rejected: %v", err)
				}
				return
			}
			findViolation(t, violationsOf(t, err), tc.field, -1)
		})
	}

	t.Run("mapping_level_carries_index_and_count", func(t *testing.T) {
		long := strings.Repeat("x", 1001)
		doc := `{"interface_name":"com.astrate.test.MapBounds","version_major":1,"version_minor":0,` +
			`"type":"datastream","ownership":"device","mappings":[` +
			`{"endpoint":"/a","type":"double"},` +
			`{"endpoint":"/b","type":"double","description":"` + long + `"}]}`
		_, err := interfaceschema.ParseInterface([]byte(doc))
		ve := violationsOf(t, err)
		if ve.MappingCount != 2 {
			t.Errorf("MappingCount = %d, want 2 (declared mappings)", ve.MappingCount)
		}
		v := findViolation(t, ve, "description", 1)
		if v.Messages[0] != "should be at most 1000 character(s)" {
			t.Errorf("message = %q", v.Messages[0])
		}
		if !strings.HasPrefix(ve.Error(), "invalid interface: ") {
			t.Errorf("Error() = %q, want the invalid-interface prefix", ve.Error())
		}
	})

	t.Run("mapping_doc_over_bound", func(t *testing.T) {
		long := strings.Repeat("x", 100001)
		doc := `{"interface_name":"com.astrate.test.MapDoc","version_major":1,"version_minor":0,` +
			`"type":"datastream","ownership":"device","mappings":[` +
			`{"endpoint":"/a","type":"double","doc":"` + long + `"}]}`
		_, err := interfaceschema.ParseInterface([]byte(doc))
		v := findViolation(t, violationsOf(t, err), "doc", 0)
		if v.Messages[0] != "should be at most 100000 character(s)" {
			t.Errorf("message = %q", v.Messages[0])
		}
	})

	t.Run("properties_mapping_bounds_enforced_too", func(t *testing.T) {
		long := strings.Repeat("x", 1001)
		doc := `{"interface_name":"com.astrate.test.PropBounds","version_major":0,"version_minor":1,` +
			`"type":"properties","ownership":"device","mappings":[` +
			`{"endpoint":"/a","type":"string","description":"` + long + `"}]}`
		_, err := interfaceschema.ParseInterface([]byte(doc))
		findViolation(t, violationsOf(t, err), "description", 0)
	})
}

// TestDatabaseRetentionTTLBand pins the probe-verified [60, 630720000) band
// (issue #61) that replaces the old ">= 1" rule.
func TestDatabaseRetentionTTLBand(t *testing.T) {
	ttlIface := func(ttl string) string {
		body := `{"interface_name":"com.astrate.test.TTL","version_major":1,"version_minor":0,` +
			`"type":"datastream","ownership":"device","mappings":[` +
			`{"endpoint":"/v","type":"double","database_retention_policy":"use_ttl"`
		if ttl != "" {
			body += `,"database_retention_ttl":` + ttl
		}
		return body + `}]}`
	}

	cases := []struct {
		name    string
		ttl     string
		wantOK  bool
		message string // exact message when rejected
	}{
		{"ttl_59_rejected", "59", false, "must be greater than or equal to 60"},
		{"ttl_60_ok", "60", true, ""},
		{"ttl_630719999_ok", "630719999", true, ""},
		{"ttl_630720000_rejected", "630720000", false, "must be less than 630720000"},
		{"ttl_630720001_rejected", "630720001", false, "must be less than 630720000"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := interfaceschema.ParseInterface([]byte(ttlIface(tc.ttl)))
			if tc.wantOK {
				if err != nil {
					t.Fatalf("accepted twin ttl=%s rejected: %v", tc.ttl, err)
				}
				return
			}
			v := findViolation(t, violationsOf(t, err), "database_retention_ttl", 0)
			if v.Messages[0] != tc.message {
				t.Errorf("message = %q, want %q", v.Messages[0], tc.message)
			}
		})
	}

	t.Run("use_ttl_without_ttl_still_plain_error", func(t *testing.T) {
		_, err := interfaceschema.ParseInterface([]byte(ttlIface("")))
		if err == nil {
			t.Fatal("use_ttl without ttl accepted, want rejection")
		}
		var ve *interfaceschema.ViolationsError
		if errors.As(err, &ve) {
			t.Errorf("missing-ttl rejection produced %+v, want a plain error", ve)
		}
		if !strings.Contains(err.Error(), "requires database_retention_ttl") {
			t.Errorf("error %q lost the historical missing-ttl message", err)
		}
	})
}
