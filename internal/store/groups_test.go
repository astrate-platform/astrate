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

func testGroupsBatch(t *testing.T, s *Store) {
	ctx := context.Background()
	realm := mustCreateRealm(t, s)
	d1 := mustRegisterDevice(t, s, realm.ID)
	d2 := mustRegisterDevice(t, s, realm.ID)
	d3 := mustRegisterDevice(t, s, realm.ID) // in no group

	ga, err := s.CreateGroup(ctx, realm.ID, "batch-a")
	if err != nil {
		t.Fatal(err)
	}
	gb, err := s.CreateGroup(ctx, realm.ID, "batch-b")
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range []struct {
		g *Group
		d deviceid.ID
	}{{ga, d1}, {gb, d1}, {ga, d2}} {
		if err := s.AddGroupDevice(ctx, m.g.ID, realm.ID, m.d); err != nil {
			t.Fatal(err)
		}
	}

	got, err := s.ListDeviceGroupsBatch(ctx, realm.ID, []deviceid.ID{d1, d2, d3})
	if err != nil {
		t.Fatalf("ListDeviceGroupsBatch: %v", err)
	}
	if len(got[d1]) != 2 || got[d1][0] != "batch-a" || got[d1][1] != "batch-b" {
		t.Errorf("d1 groups = %v, want [batch-a batch-b] sorted", got[d1])
	}
	if len(got[d2]) != 1 || got[d2][0] != "batch-a" {
		t.Errorf("d2 groups = %v", got[d2])
	}
	if _, ok := got[d3]; ok {
		t.Errorf("groupless device present in map: %v", got[d3])
	}

	empty, err := s.ListDeviceGroupsBatch(ctx, realm.ID, nil)
	if err != nil || len(empty) != 0 {
		t.Errorf("empty batch = %v, err %v", empty, err)
	}
}
