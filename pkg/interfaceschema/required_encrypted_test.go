package interfaceschema_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

func TestParseRequiredEncryptedDatastream(t *testing.T) {
	t.Run("required true parses", func(t *testing.T) {
		iface := mustParse(t, `{
			"interface_name": "com.astrate.test.RequiredField",
			"version_major": 0, "version_minor": 1,
			"type": "datastream", "ownership": "device",
			"mappings": [{"endpoint": "/v", "type": "integer", "required": true}]
		}`)
		if !iface.Mappings[0].Required {
			t.Error("Required = false, want true")
		}
		if iface.Mappings[0].Encrypted {
			t.Error("Encrypted = true, want default false")
		}
	})

	t.Run("encrypted true parses", func(t *testing.T) {
		iface := mustParse(t, `{
			"interface_name": "com.astrate.test.EncryptedField",
			"version_major": 0, "version_minor": 1,
			"type": "datastream", "ownership": "device",
			"mappings": [{"endpoint": "/v", "type": "integer", "encrypted": true}]
		}`)
		if !iface.Mappings[0].Encrypted {
			t.Error("Encrypted = false, want true")
		}
		if iface.Mappings[0].Required {
			t.Error("Required = true, want default false")
		}
	})

	t.Run("absent fields default false", func(t *testing.T) {
		iface := mustParse(t, `{
			"interface_name": "com.astrate.test.PlainField",
			"version_major": 0, "version_minor": 1,
			"type": "datastream", "ownership": "device",
			"mappings": [{"endpoint": "/v", "type": "integer"}]
		}`)
		if iface.Mappings[0].Required || iface.Mappings[0].Encrypted {
			t.Errorf("Required = %v, Encrypted = %v, want both false",
				iface.Mappings[0].Required, iface.Mappings[0].Encrypted)
		}
	})
}

func TestParseRequiredEncryptedRejectedOnProperties(t *testing.T) {
	t.Run("required rejected", func(t *testing.T) {
		_, err := interfaceschema.ParseInterface([]byte(`{
			"interface_name": "com.astrate.test.PropRequired",
			"version_major": 0, "version_minor": 1,
			"type": "properties", "ownership": "device",
			"mappings": [{"endpoint": "/v", "type": "integer", "required": true}]
		}`))
		if err == nil {
			t.Fatal(`ParseInterface accepted "required" on properties`)
		}
		if !strings.Contains(err.Error(), `"required" is not allowed on properties`) {
			t.Errorf("error %q does not mention the properties rejection", err.Error())
		}
	})

	t.Run("encrypted rejected", func(t *testing.T) {
		_, err := interfaceschema.ParseInterface([]byte(`{
			"interface_name": "com.astrate.test.PropEncrypted",
			"version_major": 0, "version_minor": 1,
			"type": "properties", "ownership": "device",
			"mappings": [{"endpoint": "/v", "type": "integer", "encrypted": true}]
		}`))
		if err == nil {
			t.Fatal(`ParseInterface accepted "encrypted" on properties`)
		}
		if !strings.Contains(err.Error(), `"encrypted" is not allowed on properties`) {
			t.Errorf("error %q does not mention the properties rejection", err.Error())
		}
	})
}

func TestObjectAggregationPerKeyRequired(t *testing.T) {
	iface := mustParse(t, `{
		"interface_name": "com.astrate.test.ObjectRequired",
		"version_major": 0, "version_minor": 1,
		"type": "datastream", "ownership": "device", "aggregation": "object",
		"mappings": [
			{"endpoint": "/lat", "type": "double", "required": true},
			{"endpoint": "/lng", "type": "double"}
		]
	}`)
	if len(iface.Mappings) != 2 {
		t.Fatalf("len(Mappings) = %d, want 2", len(iface.Mappings))
	}
	if !iface.Mappings[0].Required || iface.Mappings[1].Required {
		t.Errorf("Required = %v/%v, want true/false (no uniformity rule)",
			iface.Mappings[0].Required, iface.Mappings[1].Required)
	}
}

func TestParseCanonicalRequired(t *testing.T) {
	t.Run("alias plus required re-encodes canon", func(t *testing.T) {
		iface, canon, err := interfaceschema.ParseInterfaceCanonical([]byte(`{
			"interface_name": "com.astrate.test.CanonRequired",
			"version_major": 0, "version_minor": 1,
			"type": "datastream", "ownership": "device",
			"mappings": [{"path": "/v", "type": "integer", "required": true}]
		}`))
		if err != nil {
			t.Fatalf("ParseInterfaceCanonical: %v", err)
		}
		if iface == nil || canon == nil {
			t.Fatalf("iface = %v, canon = %v, want both non-nil", iface, canon)
		}
		if !strings.Contains(string(canon), `"required":true`) {
			t.Errorf("canonical form %s does not contain %q", canon, `"required":true`)
		}
	})

	t.Run("required alone does not trip re-encode", func(t *testing.T) {
		_, canon, err := interfaceschema.ParseInterfaceCanonical([]byte(`{
			"interface_name": "com.astrate.test.CanonRequired",
			"version_major": 0, "version_minor": 1,
			"type": "datastream", "ownership": "device",
			"mappings": [{"endpoint": "/v", "type": "integer", "required": true}]
		}`))
		if err != nil {
			t.Fatalf("ParseInterfaceCanonical: %v", err)
		}
		if canon != nil {
			t.Errorf("canon = %s, want nil (no legacy aliases consumed)", canon)
		}
	})
}

func TestCheckMinorUpgradeRequiredEncrypted(t *testing.T) {
	const oldDoc = `{
		"interface_name": "com.astrate.test.FlagUpgrade",
		"version_major": 1, "version_minor": 0,
		"type": "datastream", "ownership": "device",
		"mappings": [{"endpoint": "/v", "type": "integer"}]
	}`
	cases := []struct {
		name    string
		field   string
		wantSub string
	}{
		{name: "required flip rejected", field: `"required": true`, wantSub: "required"},
		{name: "encrypted flip rejected", field: `"encrypted": true`, wantSub: "encrypted"},
	}
	for _, tc := range cases {
		nextDoc := `{
			"interface_name": "com.astrate.test.FlagUpgrade",
			"version_major": 1, "version_minor": 1,
			"type": "datastream", "ownership": "device",
			"mappings": [{"endpoint": "/v", "type": "integer", ` + tc.field + `}]
		}`
		sameDoc := `{
			"interface_name": "com.astrate.test.FlagUpgrade",
			"version_major": 1, "version_minor": 1,
			"type": "datastream", "ownership": "device",
			"mappings": [{"endpoint": "/v", "type": "integer"}]
		}`

		t.Run(tc.name, func(t *testing.T) {
			prev := mustParse(t, oldDoc)
			next := mustParse(t, nextDoc)
			err := interfaceschema.CheckMinorUpgrade(prev, next)
			if err == nil {
				t.Fatalf("CheckMinorUpgrade accepted flipping %s on an existing mapping", tc.field)
			}
			if !strings.Contains(err.Error(), tc.wantSub) {
				t.Errorf("error %q does not mention %q", err.Error(), tc.wantSub)
			}
			if !errors.Is(err, interfaceschema.ErrIncompatibleEndpointChange) {
				t.Errorf("error %v does not classify as ErrIncompatibleEndpointChange", err)
			}
		})

		t.Run(tc.name+" acceptance twin", func(t *testing.T) {
			prev := mustParse(t, oldDoc)
			next := mustParse(t, sameDoc)
			if err := interfaceschema.CheckMinorUpgrade(prev, next); err != nil {
				t.Fatalf("CheckMinorUpgrade rejected unchanged-field upgrade: %v", err)
			}
		})
	}
}
