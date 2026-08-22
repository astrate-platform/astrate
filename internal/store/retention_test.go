//go:build integration

package store

import (
	"context"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/testutil"
)

// TestEnforceRealmRetentionCeilings covers the realm-level retention sweep
// (#72): rows older than a capped realm's ceiling vanish from both hypertables
// even when their interface declares no_ttl, fresh rows survive, and realms
// without a ceiling are untouched.
func TestEnforceRealmRetentionCeilings(t *testing.T) {
	pool := testutil.StartTimescale(t)
	dsn := pool.Config().ConnString()
	ctx := context.Background()

	s, err := New(ctx, dsn)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(s.Close)

	capped := mustCreateRealm(t, s)
	max := int64(3600)
	if err := s.UpdateRealm(ctx, capped.Name, RealmPatch{PatchRetention: true, SetRetention: max}); err != nil {
		t.Fatalf("UpdateRealm(capped): %v", err)
	}
	free := mustCreateRealm(t, s)

	device := mustRegisterDevice(t, s, capped.ID)
	freeDev := mustRegisterDevice(t, s, free.ID)

	// no_ttl mappings everywhere: only the realm ceiling can age these out.
	individual := mustInstallInterface(t, s, capped.ID, `{
		"interface_name": "com.astrate.test.RetentionValues",
		"version_major": 1,
		"version_minor": 0,
		"type": "datastream",
		"ownership": "device",
		"mappings": [{"endpoint": "/v", "type": "double"}]
	}`)
	object := mustInstallInterface(t, s, capped.ID, `{
		"interface_name": "com.astrate.test.RetentionObject",
		"version_major": 1,
		"version_minor": 0,
		"type": "datastream",
		"ownership": "device",
		"aggregation": "object",
		"mappings": [
			{"endpoint": "/%{id}/a", "type": "double"},
			{"endpoint": "/%{id}/b", "type": "double"}
		]
	}`)
	freeIface := mustInstallInterface(t, s, free.ID, `{
		"interface_name": "com.astrate.test.RetentionFree",
		"version_major": 1,
		"version_minor": 0,
		"type": "datastream",
		"ownership": "device",
		"mappings": [{"endpoint": "/v", "type": "double"}]
	}`)

	now := time.Now().UTC()
	old := now.Add(-2 * time.Hour) // beyond the 3600 s ceiling
	fresh := now.Add(-10 * time.Minute)
	v := 1.0
	batch := DatastreamBatch{
		Individual: []IndividualRow{
			{RealmID: capped.ID, DeviceID: device, InterfaceID: individual.ID,
				EndpointID: individual.Endpoints["/v"], Path: "/v",
				TS: old, ReceptionTS: now, ValueDouble: &v},
			{RealmID: capped.ID, DeviceID: device, InterfaceID: individual.ID,
				EndpointID: individual.Endpoints["/v"], Path: "/v",
				TS: fresh, ReceptionTS: now, ValueDouble: &v},
			{RealmID: free.ID, DeviceID: freeDev, InterfaceID: freeIface.ID,
				EndpointID: freeIface.Endpoints["/v"], Path: "/v",
				TS: old, ReceptionTS: now, ValueDouble: &v},
		},
		Objects: []ObjectRow{
			{RealmID: capped.ID, DeviceID: device, InterfaceID: object.ID, Path: "/1",
				TS: old, ReceptionTS: now, Value: []byte(`{"a":1}`)},
		},
	}
	if err := s.AppendDatastreams(ctx, batch); err != nil {
		t.Fatalf("AppendDatastreams: %v", err)
	}

	if err := s.EnforceRealmRetentionCeilings(ctx); err != nil {
		t.Fatalf("EnforceRealmRetentionCeilings: %v", err)
	}

	count := func(table string, realmID int16) int {
		var n int
		if err := s.pool.QueryRow(ctx,
			`SELECT count(*) FROM `+table+` WHERE realm_id = $1`, realmID).Scan(&n); err != nil {
			t.Fatalf("counting %s of realm %d: %v", table, realmID, err)
		}
		return n
	}

	if n := count("individual_datastreams", capped.ID); n != 1 {
		t.Errorf("capped realm individual rows after sweep = %d, want 1 (only the fresh row)", n)
	}
	if n := count("object_datastreams", capped.ID); n != 0 {
		t.Errorf("capped realm object rows after sweep = %d, want 0 (aged row deleted)", n)
	}
	if n := count("individual_datastreams", free.ID); n != 1 {
		t.Errorf("uncapped realm individual rows after sweep = %d, want 1 (old row untouched)", n)
	}

	// Idempotent re-run: nothing left to delete, no error.
	if err := s.EnforceRealmRetentionCeilings(ctx); err != nil {
		t.Fatalf("EnforceRealmRetentionCeilings(re-run): %v", err)
	}
}
