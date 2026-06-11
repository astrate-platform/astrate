//go:build integration

package store

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"testing"
	"time"
)

func testRealms(t *testing.T, s *Store) {
	ctx := context.Background()

	t.Run("CreateGetRoundTrip", func(t *testing.T) {
		masterKey := make([]byte, MasterKeySize)
		if _, err := rand.Read(masterKey); err != nil {
			t.Fatal(err)
		}
		ks, err := NewKeySealer(masterKey)
		if err != nil {
			t.Fatalf("NewKeySealer: %v", err)
		}
		caKeyPEM := []byte("-----BEGIN EC PRIVATE KEY-----\nrealmtest\n-----END EC PRIVATE KEY-----\n")
		sealed, err := ks.Seal(caKeyPEM)
		if err != nil {
			t.Fatalf("Seal: %v", err)
		}

		limit := int32(100)
		created, err := s.CreateRealm(ctx, NewRealm{
			Name:                    uniqueRealmName(t),
			JWTPublicKeysPEM:        []string{"pem-key-1", "pem-key-2"},
			CACertificatePEM:        "ca-cert-pem",
			CAPrivateKeySealed:      sealed,
			DeviceRegistrationLimit: &limit,
		})
		if err != nil {
			t.Fatalf("CreateRealm: %v", err)
		}
		if created.ID == 0 {
			t.Fatal("created realm has zero ID")
		}
		if created.CreatedAt.IsZero() || time.Since(created.CreatedAt) > time.Minute {
			t.Errorf("created_at implausible: %v", created.CreatedAt)
		}

		got, err := s.GetRealmByName(ctx, created.Name)
		if err != nil {
			t.Fatalf("GetRealmByName: %v", err)
		}
		if got.ID != created.ID || got.Name != created.Name {
			t.Errorf("round-trip identity mismatch: %+v vs %+v", got, created)
		}
		if len(got.JWTPublicKeysPEM) != 2 || got.JWTPublicKeysPEM[0] != "pem-key-1" || got.JWTPublicKeysPEM[1] != "pem-key-2" {
			t.Errorf("JWT keys round-trip: got %v", got.JWTPublicKeysPEM)
		}
		if got.CACertificatePEM != "ca-cert-pem" {
			t.Errorf("CA cert round-trip: got %q", got.CACertificatePEM)
		}
		if got.DeviceRegistrationLimit == nil || *got.DeviceRegistrationLimit != 100 {
			t.Errorf("registration limit round-trip: got %v", got.DeviceRegistrationLimit)
		}

		// CA key seal → store → load → open round-trip (M2 gate).
		opened, err := ks.Open(got.CAPrivateKeySealed)
		if err != nil {
			t.Fatalf("Open of stored CA key: %v", err)
		}
		if !bytes.Equal(opened, caKeyPEM) {
			t.Error("stored CA key does not open to the original plaintext")
		}

		byID, err := s.GetRealm(ctx, created.ID)
		if err != nil {
			t.Fatalf("GetRealm: %v", err)
		}
		if byID.Name != created.Name {
			t.Errorf("GetRealm by ID: got %q", byID.Name)
		}

		realms, err := s.ListRealms(ctx)
		if err != nil {
			t.Fatalf("ListRealms: %v", err)
		}
		found := false
		for _, r := range realms {
			if r.Name == created.Name {
				found = true
			}
		}
		if !found {
			t.Error("ListRealms does not contain the created realm")
		}
	})

	t.Run("DuplicateName", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		_, err := s.CreateRealm(ctx, NewRealm{
			Name:               realm.Name,
			CACertificatePEM:   "x",
			CAPrivateKeySealed: []byte("x"),
		})
		if !errors.Is(err, ErrAlreadyExists) {
			t.Fatalf("duplicate realm: got %v, want ErrAlreadyExists", err)
		}
	})

	t.Run("InvalidName", func(t *testing.T) {
		for _, name := range []string{"", "Realm", "1realm", "re-alm", "re_alm"} {
			_, err := s.CreateRealm(ctx, NewRealm{
				Name:               name,
				CACertificatePEM:   "x",
				CAPrivateKeySealed: []byte("x"),
			})
			if !errors.Is(err, ErrInvalidRealmName) {
				t.Errorf("name %q: got %v, want ErrInvalidRealmName", name, err)
			}
		}
	})

	t.Run("Updates", func(t *testing.T) {
		realm := mustCreateRealm(t, s)

		if err := s.SetRealmJWTPublicKeys(ctx, realm.Name, []string{"new-key"}); err != nil {
			t.Fatalf("SetRealmJWTPublicKeys: %v", err)
		}
		if err := s.SetRealmCA(ctx, realm.Name, "new-ca", []byte("new-sealed")); err != nil {
			t.Fatalf("SetRealmCA: %v", err)
		}
		got, err := s.GetRealmByName(ctx, realm.Name)
		if err != nil {
			t.Fatalf("GetRealmByName: %v", err)
		}
		if len(got.JWTPublicKeysPEM) != 1 || got.JWTPublicKeysPEM[0] != "new-key" {
			t.Errorf("JWT keys after update: %v", got.JWTPublicKeysPEM)
		}
		if got.CACertificatePEM != "new-ca" || !bytes.Equal(got.CAPrivateKeySealed, []byte("new-sealed")) {
			t.Error("CA material after update mismatch")
		}

		if err := s.SetRealmJWTPublicKeys(ctx, "nosuchrealm", nil); !errors.Is(err, ErrNotFound) {
			t.Errorf("update of unknown realm: got %v, want ErrNotFound", err)
		}
	})

	t.Run("CascadeDelete", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		device := mustRegisterDevice(t, s, realm.ID)
		si := mustInstallInterface(t, s, realm.ID, `{
			"interface_name": "com.astrate.test.CascadeValues",
			"version_major": 1,
			"version_minor": 0,
			"type": "datastream",
			"ownership": "device",
			"mappings": [{"endpoint": "/%{sensor_id}/value", "type": "double"}]
		}`)
		props := mustInstallInterface(t, s, realm.ID, `{
			"interface_name": "com.astrate.test.CascadeProps",
			"version_major": 1,
			"version_minor": 0,
			"type": "properties",
			"ownership": "device",
			"mappings": [{"endpoint": "/enabled", "type": "boolean"}]
		}`)

		if err := s.UpsertProperty(ctx, Property{
			RealmID: realm.ID, DeviceID: device, InterfaceID: props.ID,
			EndpointID: props.Endpoints["/enabled"], Path: "/enabled",
			Value: []byte("true"), ValueType: mustValueType(t, "boolean"),
		}); err != nil {
			t.Fatalf("UpsertProperty: %v", err)
		}

		now := time.Now().UTC()
		v := 1.5
		if err := s.AppendDatastreams(ctx, DatastreamBatch{
			Individual: []IndividualRow{{
				RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID,
				EndpointID: si.Endpoints["/%{sensor_id}/value"], Path: "/1/value",
				TS: now, ReceptionTS: now, ValueDouble: &v,
			}},
			Objects: []ObjectRow{{
				RealmID: realm.ID, DeviceID: device, InterfaceID: si.ID,
				Path: "/1", TS: now, ReceptionTS: now, Value: []byte(`{"a":1}`),
			}},
		}); err != nil {
			t.Fatalf("AppendDatastreams: %v", err)
		}

		group, err := s.CreateGroup(ctx, realm.ID, "cascadegroup")
		if err != nil {
			t.Fatalf("CreateGroup: %v", err)
		}
		if err := s.AddGroupDevice(ctx, group.ID, realm.ID, device); err != nil {
			t.Fatalf("AddGroupDevice: %v", err)
		}
		if _, err := s.CreateTrigger(ctx, realm.ID, "cascadetrigger", []byte(`{"name":"cascadetrigger"}`)); err != nil {
			t.Fatalf("CreateTrigger: %v", err)
		}

		if err := s.DeleteRealm(ctx, realm.Name); err != nil {
			t.Fatalf("DeleteRealm: %v", err)
		}

		if _, err := s.GetRealmByName(ctx, realm.Name); !errors.Is(err, ErrNotFound) {
			t.Fatalf("realm still present after delete: %v", err)
		}
		for _, q := range []struct {
			table string
			sql   string
		}{
			{"interfaces", `SELECT count(*) FROM interfaces WHERE realm_id = $1`},
			{"endpoints", `SELECT count(*) FROM endpoints e JOIN interfaces i ON i.id = e.interface_id WHERE i.realm_id = $1`},
			{"devices", `SELECT count(*) FROM devices WHERE realm_id = $1`},
			{"properties", `SELECT count(*) FROM properties WHERE realm_id = $1`},
			{"individual_datastreams", `SELECT count(*) FROM individual_datastreams WHERE realm_id = $1`},
			{"object_datastreams", `SELECT count(*) FROM object_datastreams WHERE realm_id = $1`},
			{"groups", `SELECT count(*) FROM groups WHERE realm_id = $1`},
			{"group_devices", `SELECT count(*) FROM group_devices WHERE realm_id = $1`},
			{"triggers", `SELECT count(*) FROM triggers WHERE realm_id = $1`},
		} {
			var n int
			if err := s.pool.QueryRow(ctx, q.sql, realm.ID).Scan(&n); err != nil {
				t.Fatalf("counting %s: %v", q.table, err)
			}
			if n != 0 {
				t.Errorf("%s: %d rows survived the realm cascade", q.table, n)
			}
		}

		if err := s.DeleteRealm(ctx, realm.Name); !errors.Is(err, ErrNotFound) {
			t.Errorf("second delete: got %v, want ErrNotFound", err)
		}
	})
}
