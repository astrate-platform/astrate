//go:build integration

package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

const valuesV10 = `{
	"interface_name": "com.astrate.test.StoreValues",
	"version_major": 1,
	"version_minor": 0,
	"type": "datastream",
	"ownership": "device",
	"mappings": [
		{"endpoint": "/%{sensor_id}/value", "type": "double", "explicit_timestamp": true, "reliability": "guaranteed"}
	]
}`

const valuesV11 = `{
	"interface_name": "com.astrate.test.StoreValues",
	"version_major": 1,
	"version_minor": 1,
	"type": "datastream",
	"ownership": "device",
	"mappings": [
		{"endpoint": "/%{sensor_id}/value", "type": "double", "explicit_timestamp": true, "reliability": "guaranteed"},
		{"endpoint": "/%{sensor_id}/unit", "type": "string", "explicit_timestamp": true}
	]
}`

func testInterfaces(t *testing.T, s *Store) {
	ctx := context.Background()

	t.Run("InstallGeneratedColumns", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		si := mustInstallInterface(t, s, realm.ID, valuesV10)

		if si.Name != "com.astrate.test.StoreValues" || si.Major != 1 || si.Minor != 0 {
			t.Errorf("identity: %+v", si)
		}
		if si.Type != interfaceschema.Datastream || si.Ownership != interfaceschema.OwnershipDevice ||
			si.Aggregation != interfaceschema.AggregationIndividual {
			t.Errorf("enums: type=%v ownership=%v aggregation=%v", si.Type, si.Ownership, si.Aggregation)
		}

		// The generated columns must derive the same values from the raw JSON.
		var (
			name, typ, ownership, aggregation string
			major, minor                      int
		)
		err := s.pool.QueryRow(ctx, `
			SELECT name, major_version, minor_version, type, ownership, aggregation
			FROM interfaces WHERE id = $1`, si.ID).
			Scan(&name, &major, &minor, &typ, &ownership, &aggregation)
		if err != nil {
			t.Fatalf("reading generated columns: %v", err)
		}
		if name != si.Name || major != 1 || minor != 0 || typ != "datastream" ||
			ownership != "device" || aggregation != "individual" {
			t.Errorf("generated columns: name=%q major=%d minor=%d type=%q ownership=%q aggregation=%q",
				name, major, minor, typ, ownership, aggregation)
		}

		// Endpoint attribute columns mirror the parsed mapping.
		var (
			valueType, reliability, retention, policy string
			explicitTS, allowUnset                    bool
		)
		err = s.pool.QueryRow(ctx, `
			SELECT value_type, reliability, retention, database_retention_policy, explicit_timestamp, allow_unset
			FROM endpoints WHERE id = $1`, si.Endpoints["/%{sensor_id}/value"]).
			Scan(&valueType, &reliability, &retention, &policy, &explicitTS, &allowUnset)
		if err != nil {
			t.Fatalf("reading endpoint row: %v", err)
		}
		if valueType != "double" || reliability != "guaranteed" || retention != "discard" ||
			policy != "no_ttl" || !explicitTS || allowUnset {
			t.Errorf("endpoint attributes: %q %q %q %q %v %v",
				valueType, reliability, retention, policy, explicitTS, allowUnset)
		}

		if _, err := s.InstallInterface(ctx, realm.ID, []byte(valuesV10)); !errors.Is(err, ErrAlreadyExists) {
			t.Errorf("duplicate install: got %v, want ErrAlreadyExists", err)
		}
	})

	t.Run("EndpointIDStability", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		installed := mustInstallInterface(t, s, realm.ID, valuesV10)
		valueEpID := installed.Endpoints["/%{sensor_id}/value"]
		if valueEpID == 0 {
			t.Fatal("install returned zero endpoint ID")
		}

		got, err := s.GetInterface(ctx, realm.ID, installed.Name, installed.Major)
		if err != nil {
			t.Fatalf("GetInterface: %v", err)
		}
		if got.Endpoints["/%{sensor_id}/value"] != valueEpID {
			t.Errorf("reload changed endpoint ID: %d vs %d", got.Endpoints["/%{sensor_id}/value"], valueEpID)
		}

		updated, err := s.UpdateInterface(ctx, realm.ID, []byte(valuesV11))
		if err != nil {
			t.Fatalf("UpdateInterface: %v", err)
		}
		if updated.ID != installed.ID {
			t.Errorf("minor update changed interface ID: %d vs %d", updated.ID, installed.ID)
		}
		if updated.Minor != 1 {
			t.Errorf("minor after update: %d", updated.Minor)
		}
		if updated.Endpoints["/%{sensor_id}/value"] != valueEpID {
			t.Errorf("minor update changed existing endpoint ID: %d vs %d",
				updated.Endpoints["/%{sensor_id}/value"], valueEpID)
		}
		if updated.Endpoints["/%{sensor_id}/unit"] == 0 {
			t.Error("minor update did not insert the new endpoint")
		}

		loaded, err := s.LoadRealmInterfaces(ctx, realm.ID)
		if err != nil {
			t.Fatalf("LoadRealmInterfaces: %v", err)
		}
		if len(loaded) != 1 {
			t.Fatalf("LoadRealmInterfaces: got %d interfaces, want 1", len(loaded))
		}
		if loaded[0].Endpoints["/%{sensor_id}/value"] != valueEpID {
			t.Errorf("LoadRealmInterfaces changed endpoint ID: %d vs %d",
				loaded[0].Endpoints["/%{sensor_id}/value"], valueEpID)
		}

		// Dropping an endpoint in an update must be refused.
		droppedDef := `{
			"interface_name": "com.astrate.test.StoreValues",
			"version_major": 1,
			"version_minor": 2,
			"type": "datastream",
			"ownership": "device",
			"mappings": [
				{"endpoint": "/%{sensor_id}/value", "type": "double", "explicit_timestamp": true, "reliability": "guaranteed"}
			]
		}`
		if _, err := s.UpdateInterface(ctx, realm.ID, []byte(droppedDef)); err == nil {
			t.Error("update dropping an endpoint succeeded")
		}

		// Compile through the StoredInterface resolver: the stamped IDs must
		// be the stored ones.
		iface, err := interfaceschema.ParseInterface([]byte(valuesV11))
		if err != nil {
			t.Fatalf("ParseInterface: %v", err)
		}
		ci, err := interfaceschema.Compile(iface, updated)
		if err != nil {
			t.Fatalf("Compile with store resolver: %v", err)
		}
		if ci.ID != installed.ID {
			t.Errorf("compiled interface ID: %d vs %d", ci.ID, installed.ID)
		}
		m, ok := ci.Trie.Match("/4/value")
		if !ok {
			t.Fatal("trie match /4/value failed")
		}
		if m.EndpointID != valueEpID {
			t.Errorf("compiled endpoint ID: %d vs %d", m.EndpointID, valueEpID)
		}
	})

	t.Run("UpdateUnknown", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		if _, err := s.UpdateInterface(ctx, realm.ID, []byte(valuesV11)); !errors.Is(err, ErrNotFound) {
			t.Errorf("update of unknown interface: got %v, want ErrNotFound", err)
		}
	})

	t.Run("DeleteRules", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		mustInstallInterface(t, s, realm.ID, valuesV10)

		// Major != 0 is never deletable.
		err := s.DeleteInterface(ctx, realm.ID, "com.astrate.test.StoreValues", 1)
		if !errors.Is(err, ErrInterfaceMajorNotZero) {
			t.Fatalf("delete of major 1: got %v, want ErrInterfaceMajorNotZero", err)
		}

		draft := mustInstallInterface(t, s, realm.ID, `{
			"interface_name": "com.astrate.test.Draft",
			"version_major": 0,
			"version_minor": 1,
			"type": "datastream",
			"ownership": "device",
			"mappings": [{"endpoint": "/v", "type": "integer"}]
		}`)

		// A device declaring the interface blocks deletion.
		device := mustRegisterDevice(t, s, realm.ID)
		if _, err := s.UpdateIntrospection(ctx, realm.ID, device, map[string]InterfaceVersion{
			"com.astrate.test.Draft": {Major: 0, Minor: 1},
		}); err != nil {
			t.Fatalf("UpdateIntrospection: %v", err)
		}
		err = s.DeleteInterface(ctx, realm.ID, "com.astrate.test.Draft", 0)
		if !errors.Is(err, ErrInterfaceInUse) {
			t.Fatalf("delete of declared interface: got %v, want ErrInterfaceInUse", err)
		}

		// Stale datastream rows must be swept by the delete.
		now := time.Now().UTC()
		v := int32(7)
		if err := s.AppendDatastreams(ctx, DatastreamBatch{Individual: []IndividualRow{{
			RealmID: realm.ID, DeviceID: device, InterfaceID: draft.ID,
			EndpointID: draft.Endpoints["/v"], Path: "/v",
			TS: now, ReceptionTS: now, ValueInteger: &v,
		}}}); err != nil {
			t.Fatalf("AppendDatastreams: %v", err)
		}

		if _, err := s.UpdateIntrospection(ctx, realm.ID, device, nil); err != nil {
			t.Fatalf("clearing introspection: %v", err)
		}
		if err := s.DeleteInterface(ctx, realm.ID, "com.astrate.test.Draft", 0); err != nil {
			t.Fatalf("delete of drained draft interface: %v", err)
		}

		if _, err := s.GetInterface(ctx, realm.ID, "com.astrate.test.Draft", 0); !errors.Is(err, ErrNotFound) {
			t.Errorf("interface still present after delete: %v", err)
		}
		var n int
		if err := s.pool.QueryRow(ctx,
			`SELECT count(*) FROM individual_datastreams WHERE realm_id = $1 AND interface_id = $2`,
			realm.ID, draft.ID).Scan(&n); err != nil {
			t.Fatal(err)
		}
		if n != 0 {
			t.Errorf("%d datastream rows survived interface delete", n)
		}

		if err := s.DeleteInterface(ctx, realm.ID, "com.astrate.test.Draft", 0); !errors.Is(err, ErrNotFound) {
			t.Errorf("second delete: got %v, want ErrNotFound", err)
		}
	})
}
