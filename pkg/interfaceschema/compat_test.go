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

// pinnedV10 pins the three datastream-only immutable attribute pairs at
// non-zero values, so the flip cases below exercise a real mutation of one
// attribute instead of a default-to-value transition.
const pinnedV10 = `{
	"interface_name": "com.astrate.test.Upgrade",
	"version_major": 1,
	"version_minor": 0,
	"type": "datastream",
	"ownership": "device",
	"mappings": [
		{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
		 "retention": "stored", "expiry": 3600,
		 "database_retention_policy": "use_ttl", "database_retention_ttl": 3600},
		{"endpoint": "/%{sensor_id}/status", "type": "string"}
	]
}`

// pinnedV11 is pinnedV10 with the minor bumped: every pinned immutable
// attribute stays put, which an upgrade must accept. It is the
// at-immutability acceptance twin of the retention/expiry/database-retention
// flip cases below.
const pinnedV11 = `{
	"interface_name": "com.astrate.test.Upgrade",
	"version_major": 1,
	"version_minor": 1,
	"type": "datastream",
	"ownership": "device",
	"mappings": [
		{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
		 "retention": "stored", "expiry": 3600,
		 "database_retention_policy": "use_ttl", "database_retention_ttl": 3600},
		{"endpoint": "/%{sensor_id}/status", "type": "string"}
	]
}`

// pinnedPropsV10 is a properties interface carrying allow_unset (a
// properties-only immutable attribute) at its non-zero value.
const pinnedPropsV10 = `{
	"interface_name": "com.astrate.test.Upgrade",
	"version_major": 1,
	"version_minor": 0,
	"type": "properties",
	"ownership": "device",
	"mappings": [
		{"endpoint": "/%{sensor_id}/enabled", "type": "boolean", "allow_unset": true},
		{"endpoint": "/%{sensor_id}/mode", "type": "string"}
	]
}`

// pinnedPropsV11 is the at-immutability acceptance twin of the allow_unset
// flip case: allow_unset stays true across the minor bump.
const pinnedPropsV11 = `{
	"interface_name": "com.astrate.test.Upgrade",
	"version_major": 1,
	"version_minor": 1,
	"type": "properties",
	"ownership": "device",
	"mappings": [
		{"endpoint": "/%{sensor_id}/enabled", "type": "boolean", "allow_unset": true},
		{"endpoint": "/%{sensor_id}/mode", "type": "string"}
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
		{
			name:   "retention at immutability accepted",
			oldDoc: pinnedV10,
			next:   pinnedV11,
		},
		{
			name:   "retention mutated",
			oldDoc: pinnedV10,
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
					 "retention": "volatile", "expiry": 3600,
					 "database_retention_policy": "use_ttl", "database_retention_ttl": 3600},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			wantSub:   "changed retention",
			wantClass: interfaceschema.ErrIncompatibleEndpointChange,
		},
		{
			name:   "expiry at immutability accepted",
			oldDoc: pinnedV10,
			next:   pinnedV11,
		},
		{
			name:   "expiry mutated",
			oldDoc: pinnedV10,
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
					 "retention": "stored", "expiry": 7200,
					 "database_retention_policy": "use_ttl", "database_retention_ttl": 3600},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			wantSub:   "changed expiry",
			wantClass: interfaceschema.ErrIncompatibleEndpointChange,
		},
		{
			name:   "database retention at immutability accepted",
			oldDoc: pinnedV10,
			next:   pinnedV11,
		},
		{
			name:   "database_retention_policy mutated",
			oldDoc: pinnedV10,
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
					 "retention": "stored", "expiry": 3600,
					 "database_retention_policy": "no_ttl"},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			wantSub:   "changed database_retention_policy",
			wantClass: interfaceschema.ErrIncompatibleEndpointChange,
		},
		{
			name:   "database_retention_ttl mutated",
			oldDoc: pinnedV10,
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "datastream", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
					 "retention": "stored", "expiry": 3600,
					 "database_retention_policy": "use_ttl", "database_retention_ttl": 7200},
					{"endpoint": "/%{sensor_id}/status", "type": "string"}
				]
			}`,
			wantSub:   "changed database_retention_ttl",
			wantClass: interfaceschema.ErrIncompatibleEndpointChange,
		},
		{
			name:   "allow_unset at immutability accepted",
			oldDoc: pinnedPropsV10,
			next:   pinnedPropsV11,
		},
		{
			name:   "allow_unset mutated",
			oldDoc: pinnedPropsV10,
			next: `{
				"interface_name": "com.astrate.test.Upgrade",
				"version_major": 1, "version_minor": 1,
				"type": "properties", "ownership": "device",
				"mappings": [
					{"endpoint": "/%{sensor_id}/enabled", "type": "boolean", "allow_unset": false},
					{"endpoint": "/%{sensor_id}/mode", "type": "string"}
				]
			}`,
			wantSub:   "changed allow_unset",
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
