//go:build integration

package store

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

func testGroups(t *testing.T, s *Store) {
	ctx := context.Background()
	realm := mustCreateRealm(t, s)

	group, err := s.CreateGroup(ctx, realm.ID, "fleet1")
	if err != nil {
		t.Fatalf("CreateGroup: %v", err)
	}
	if group.ID == 0 || group.Name != "fleet1" || group.RealmID != realm.ID {
		t.Errorf("created group: %+v", group)
	}
	if _, err := s.CreateGroup(ctx, realm.ID, "fleet1"); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("duplicate group: got %v, want ErrAlreadyExists", err)
	}

	got, err := s.GetGroupByName(ctx, realm.ID, "fleet1")
	if err != nil {
		t.Fatalf("GetGroupByName: %v", err)
	}
	if got.ID != group.ID {
		t.Errorf("get round-trip: %+v", got)
	}

	groups, err := s.ListGroups(ctx, realm.ID)
	if err != nil {
		t.Fatalf("ListGroups: %v", err)
	}
	if len(groups) != 1 || groups[0].Name != "fleet1" {
		t.Errorf("ListGroups: %+v", groups)
	}

	// Membership: composite FK semantics.
	d1 := mustRegisterDevice(t, s, realm.ID)
	d2 := mustRegisterDevice(t, s, realm.ID)
	if err := s.AddGroupDevice(ctx, group.ID, realm.ID, d1); err != nil {
		t.Fatalf("AddGroupDevice d1: %v", err)
	}
	if err := s.AddGroupDevice(ctx, group.ID, realm.ID, d2); err != nil {
		t.Fatalf("AddGroupDevice d2: %v", err)
	}
	if err := s.AddGroupDevice(ctx, group.ID, realm.ID, d1); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("duplicate membership: got %v, want ErrAlreadyExists", err)
	}

	ghost, _ := deviceid.Random()
	if err := s.AddGroupDevice(ctx, group.ID, realm.ID, ghost); !errors.Is(err, ErrNotFound) {
		t.Errorf("membership of unknown device: got %v, want ErrNotFound", err)
	}
	if err := s.AddGroupDevice(ctx, group.ID+9999, realm.ID, d1); !errors.Is(err, ErrNotFound) {
		t.Errorf("membership in unknown group: got %v, want ErrNotFound", err)
	}
	// A group of another realm must not accept this realm's devices.
	otherRealm := mustCreateRealm(t, s)
	if err := s.AddGroupDevice(ctx, group.ID, otherRealm.ID, d1); !errors.Is(err, ErrNotFound) {
		t.Errorf("cross-realm membership: got %v, want ErrNotFound", err)
	}

	members, err := s.ListGroupDevices(ctx, group.ID)
	if err != nil {
		t.Fatalf("ListGroupDevices: %v", err)
	}
	if len(members) != 2 {
		t.Errorf("members: %v", members)
	}

	deviceGroups, err := s.ListDeviceGroups(ctx, realm.ID, d1)
	if err != nil {
		t.Fatalf("ListDeviceGroups: %v", err)
	}
	if len(deviceGroups) != 1 || deviceGroups[0] != "fleet1" {
		t.Errorf("groups of d1: %v", deviceGroups)
	}

	// Deleting a device cascades it out of the group (composite FK).
	if _, err := s.pool.Exec(ctx, `DELETE FROM devices WHERE realm_id = $1 AND id = $2`,
		realm.ID, pgtype.UUID{Bytes: d2, Valid: true}); err != nil {
		t.Fatalf("deleting device row: %v", err)
	}
	members, err = s.ListGroupDevices(ctx, group.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(members) != 1 || members[0] != d1 {
		t.Errorf("members after device delete: %v", members)
	}

	if err := s.RemoveGroupDevice(ctx, group.ID, realm.ID, d1); err != nil {
		t.Fatalf("RemoveGroupDevice: %v", err)
	}
	if err := s.RemoveGroupDevice(ctx, group.ID, realm.ID, d1); !errors.Is(err, ErrNotFound) {
		t.Errorf("second removal: got %v, want ErrNotFound", err)
	}

	if err := s.DeleteGroup(ctx, realm.ID, "fleet1"); err != nil {
		t.Fatalf("DeleteGroup: %v", err)
	}
	if _, err := s.GetGroupByName(ctx, realm.ID, "fleet1"); !errors.Is(err, ErrNotFound) {
		t.Errorf("group survived delete: %v", err)
	}
	if err := s.DeleteGroup(ctx, realm.ID, "fleet1"); !errors.Is(err, ErrNotFound) {
		t.Errorf("second group delete: got %v, want ErrNotFound", err)
	}
}
