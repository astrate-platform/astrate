//go:build integration

package store

import (
	"context"
	"errors"
	"testing"
)

func testProperties(t *testing.T, s *Store) {
	ctx := context.Background()

	realm := mustCreateRealm(t, s)
	device := mustRegisterDevice(t, s, realm.ID)

	deviceProps := mustInstallInterface(t, s, realm.ID, `{
		"interface_name": "com.astrate.test.DeviceProps",
		"version_major": 1,
		"version_minor": 0,
		"type": "properties",
		"ownership": "device",
		"mappings": [
			{"endpoint": "/p1", "type": "integer", "allow_unset": true},
			{"endpoint": "/%{id}/setting", "type": "string"}
		]
	}`)
	serverProps := mustInstallInterface(t, s, realm.ID, `{
		"interface_name": "com.astrate.test.ServerProps",
		"version_major": 1,
		"version_minor": 0,
		"type": "properties",
		"ownership": "server",
		"mappings": [{"endpoint": "/s1", "type": "string"}]
	}`)

	set := func(t *testing.T, iface *StoredInterface, endpoint, path, valueJSON, valueType string) {
		t.Helper()
		if err := s.UpsertProperty(ctx, Property{
			RealmID: realm.ID, DeviceID: device, InterfaceID: iface.ID,
			EndpointID: iface.Endpoints[endpoint], Path: path,
			Value: []byte(valueJSON), ValueType: mustValueType(t, valueType),
		}); err != nil {
			t.Fatalf("UpsertProperty %s: %v", path, err)
		}
	}

	t.Run("UpsertIdempotent", func(t *testing.T) {
		set(t, deviceProps, "/p1", "/p1", "1", "integer")
		set(t, deviceProps, "/p1", "/p1", "2", "integer")

		p, err := s.GetProperty(ctx, realm.ID, device, deviceProps.ID, "/p1")
		if err != nil {
			t.Fatalf("GetProperty: %v", err)
		}
		if string(p.Value) != "2" {
			t.Errorf("value after double upsert: %s", p.Value)
		}
		if p.ValueType.String() != "integer" || p.EndpointID != deviceProps.Endpoints["/p1"] {
			t.Errorf("metadata: type=%v endpoint=%d", p.ValueType, p.EndpointID)
		}

		list, err := s.ListProperties(ctx, realm.ID, device, deviceProps.ID)
		if err != nil {
			t.Fatalf("ListProperties: %v", err)
		}
		if len(list) != 1 {
			t.Errorf("double upsert produced %d rows", len(list))
		}
	})

	t.Run("Unset", func(t *testing.T) {
		set(t, deviceProps, "/p1", "/p1", "3", "integer")
		found, err := s.UnsetProperty(ctx, realm.ID, device, deviceProps.ID, "/p1")
		if err != nil {
			t.Fatalf("UnsetProperty: %v", err)
		}
		if !found {
			t.Error("unset of existing property reported not-found")
		}
		if _, err := s.GetProperty(ctx, realm.ID, device, deviceProps.ID, "/p1"); !errors.Is(err, ErrNotFound) {
			t.Errorf("property survived unset: %v", err)
		}
		found, err = s.UnsetProperty(ctx, realm.ID, device, deviceProps.ID, "/p1")
		if err != nil {
			t.Fatal(err)
		}
		if found {
			t.Error("second unset reported a deletion")
		}
	})

	t.Run("PurgeComplement", func(t *testing.T) {
		set(t, deviceProps, "/p1", "/p1", "10", "integer")
		set(t, deviceProps, "/%{id}/setting", "/3/setting", `"a"`, "string")
		set(t, deviceProps, "/%{id}/setting", "/4/setting", `"b"`, "string")
		set(t, serverProps, "/s1", "/s1", `"server"`, "string")

		// The device reports it still holds /p1 and /4/setting: exactly the
		// complement (/3/setting) goes, server-owned rows are untouched.
		deleted, err := s.PurgeDeviceOwnedExcept(ctx, realm.ID, device, []PropertyRef{
			{InterfaceID: deviceProps.ID, Path: "/p1"},
			{InterfaceID: deviceProps.ID, Path: "/4/setting"},
		})
		if err != nil {
			t.Fatalf("PurgeDeviceOwnedExcept: %v", err)
		}
		if deleted != 1 {
			t.Errorf("purge deleted %d rows, want 1", deleted)
		}

		remaining, err := s.ListProperties(ctx, realm.ID, device, deviceProps.ID)
		if err != nil {
			t.Fatal(err)
		}
		paths := map[string]bool{}
		for _, p := range remaining {
			paths[p.Path] = true
		}
		if len(paths) != 2 || !paths["/p1"] || !paths["/4/setting"] {
			t.Errorf("device-owned after purge: %v", paths)
		}

		serverOwned, err := s.ListServerOwnedProperties(ctx, realm.ID, device)
		if err != nil {
			t.Fatalf("ListServerOwnedProperties: %v", err)
		}
		if len(serverOwned) != 1 || serverOwned[0].Path != "/s1" || string(serverOwned[0].Value) != `"server"` {
			t.Errorf("server-owned after purge: %+v", serverOwned)
		}

		// An empty keep list purges every device-owned property.
		deleted, err = s.PurgeDeviceOwnedExcept(ctx, realm.ID, device, nil)
		if err != nil {
			t.Fatal(err)
		}
		if deleted != 2 {
			t.Errorf("full purge deleted %d rows, want 2", deleted)
		}
		serverOwned, err = s.ListServerOwnedProperties(ctx, realm.ID, device)
		if err != nil {
			t.Fatal(err)
		}
		if len(serverOwned) != 1 {
			t.Errorf("server-owned property purged by device resync: %+v", serverOwned)
		}
	})
}
