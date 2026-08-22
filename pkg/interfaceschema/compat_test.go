package interfaceschema_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// baseV10 is the v1.0 document the upgrade table mutates.
const baseV10 = `{
	"interface_name": "com.astrate.test.Upgrade",
	"version_major": 1,
	"version_minor": 0,
	"type": "datastream",
	"ownership": "device",
	"description": "before",
	"mappings": [
		{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
		 "description": "before"},
		{"endpoint": "/%{sensor_id}/status", "type": "string"}
	]
}`

func mustParse(t *testing.T, doc string) *interfaceschema.Interface {
	t.Helper()
	iface, err := interfaceschema.ParseInterface([]byte(doc))
	if err != nil {
		t.Fatalf("ParseInterface: %v", err)
	}
	return iface
}

func TestCheckMinorUpgradeTable(t *testing.T) {
	cases := []struct {
		name      string
		oldDoc    string // "" = baseV10
		next      string
		wantSub   string // "" = accepted
		wantClass error  // classification sentinel the rejection must wrap
	}{
		{
			name: "additive mapping with minor bump",
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed"},
					{"endpoint": "/%{sensor_id}/status", "type": "string"},
					{"endpoint": "/%{sensor_id}/extra", "type": "boolean"}
				]
			}`,
		},
		{
			name: "doc and description changes are non-semantic",
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 2,
				"type": "datastream", "ownership": "device",
				"description": "after",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
					 "description": "after", "doc": "new docs"},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
		},
		{
			name: "placeholder rename is non-semantic",
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{device}/value", "type": "double", "reliability": "guaranteed"},
					{"endpoint": "/%{device}/status", "type": "string"}
				]
			}`,
		},
		{
			name: "minor not increased",
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 0,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed"},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			wantSub:   "minor version was not increased",
			wantClass: interfaceschema.ErrMinorNotIncreased,
		},
		{
			name: "minor downgraded",
			oldDoc: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed"},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 0,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed"},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			wantSub:   "downgrade not allowed",
			wantClass: interfaceschema.ErrDowngradeNotAllowed,
		},
		{
			name: "major changed",
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 2, "version_minor": 0,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed"},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			wantSub: "major version changed",
		},
		{
			name: "name changed",
			next: `{
				"interface_name": "com.astrate.test.Renamed",
				"version_major": 1, "version_minor": 1,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed"},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			wantSub: "name changed",
		},
		{
			name: "type changed",
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "properties", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double"},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			wantSub:   "type changed",
			wantClass: interfaceschema.ErrIncompatibleEndpointChange,
		},
		{
			name: "ownership changed",
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "datastream", "ownership": "server",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed"},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			wantSub:   "ownership changed",
			wantClass: interfaceschema.ErrIncompatibleEndpointChange,
		},
		{
			name: "mapping removed",
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed"}
				]
			}`,
			wantSub:   "removed",
			wantClass: interfaceschema.ErrMissingEndpoints,
		},
		{
			name: "value type mutated",
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "integer", "reliability": "guaranteed"},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			wantSub:   "changed type",
			wantClass: interfaceschema.ErrIncompatibleEndpointChange,
		},
		{
			name: "reliability mutated",
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "unique"},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			wantSub:   "changed reliability",
			wantClass: interfaceschema.ErrIncompatibleEndpointChange,
		},
		{
			name: "explicit_timestamp mutated",
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
					 "explicit_timestamp": true},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			wantSub:   "changed explicit_timestamp",
			wantClass: interfaceschema.ErrIncompatibleEndpointChange,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			oldDoc := tc.oldDoc
			if oldDoc == "" {
				oldDoc = baseV10
			}
			prev := mustParse(t, oldDoc)
			next := mustParse(t, tc.next)
			err := interfaceschema.CheckMinorUpgrade(prev, next)
			if tc.wantSub == "" {
				if err != nil {
					t.Fatalf("CheckMinorUpgrade rejected valid upgrade: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("CheckMinorUpgrade accepted, want error containing %q", tc.wantSub)
			}
			if !errors.Is(err, interfaceschema.ErrIncompatibleUpgrade) {
				t.Errorf("error %v does not wrap ErrIncompatibleUpgrade", err)
			}
			if !strings.Contains(err.Error(), tc.wantSub) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.wantSub)
			}
			if tc.wantClass != nil && !errors.Is(err, tc.wantClass) {
				t.Errorf("error %v does not classify as %v", err, tc.wantClass)
			}
			if tc.wantClass == nil {
				for _, other := range []error{
					interfaceschema.ErrMinorNotIncreased,
					interfaceschema.ErrDowngradeNotAllowed,
					interfaceschema.ErrMissingEndpoints,
					interfaceschema.ErrIncompatibleEndpointChange,
				} {
					if errors.Is(err, other) {
						t.Errorf("unclassified error %v unexpectedly wraps %v", err, other)
					}
				}
			}
		})
	}
}

func TestCheckMinorUpgradeAggregationChange(t *testing.T) {
	old := mustParse(t, `{
		"interface_name": "com.astrate.test.Agg",
		"version_major": 1, "version_minor": 0,
		"type": "datastream", "ownership": "device",
		"mappings": [{"endpoint": "/lat", "type": "double"}, {"endpoint": "/lng", "type": "double"}]
	}`)
	next := mustParse(t, `{
		"interface_name": "com.astrate.test.Agg",
		"version_major": 1, "version_minor": 1,
		"type": "datastream", "ownership": "device", "aggregation": "object",
		"mappings": [{"endpoint": "/lat", "type": "double"}, {"endpoint": "/lng", "type": "double"}]
	}`)
	err := interfaceschema.CheckMinorUpgrade(old, next)
	if err == nil || !strings.Contains(err.Error(), "aggregation changed") {
		t.Errorf("err = %v, want aggregation change rejection", err)
	}
	if !errors.Is(err, interfaceschema.ErrIncompatibleEndpointChange) {
		t.Errorf("error %v does not classify as ErrIncompatibleEndpointChange", err)
	}
}

func TestCheckMinorUpgradeNil(t *testing.T) {
	old := mustParse(t, baseV10)
	if err := interfaceschema.CheckMinorUpgrade(nil, old); err == nil {
		t.Error("CheckMinorUpgrade(nil, x) succeeded")
	}
	if err := interfaceschema.CheckMinorUpgrade(old, nil); err == nil {
		t.Error("CheckMinorUpgrade(x, nil) succeeded")
	}
}
