//go:build integration

package store

import (
	"context"
	"errors"
	"testing"
)

func testUserBlocks(t *testing.T, s *Store) {
	ctx := context.Background()
	realm := mustCreateRealm(t, s)

	src := []byte("def handle_message(msg): pass")
	b, err := s.CreateUserBlock(ctx, realm.ID, &UserBlock{
		Name:         "b1",
		BlockType:    "producer",
		Source:       src,
		ConfigSchema: []byte(`{"type":"object"}`),
	})
	if err != nil {
		t.Fatalf("CreateUserBlock: %v", err)
	}
	if b.ID == 0 || b.Name != "b1" || b.BlockType != "producer" {
		t.Errorf("created user block = %+v", b)
	}
	if _, err := s.CreateUserBlock(ctx, realm.ID, &UserBlock{Name: "b1", BlockType: "consumer", Source: src}); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("duplicate create = %v, want ErrAlreadyExists", err)
	}

	if _, err := s.CreateUserBlock(ctx, realm.ID, &UserBlock{Name: "badtype", BlockType: "neither", Source: src}); err == nil {
		t.Error("invalid block_type create = nil error, want CHECK rejection")
	}

	nilSchema := mustCreateRealm(t, s)
	nb, err := s.CreateUserBlock(ctx, nilSchema.ID, &UserBlock{Name: "noschema", BlockType: "consumer", Source: src})
	if err != nil {
		t.Fatalf("CreateUserBlock (NULL schema): %v", err)
	}
	if nb.ConfigSchema != nil {
		t.Errorf("created ConfigSchema = %q, want nil", nb.ConfigSchema)
	}
	gotNil, err := s.GetUserBlock(ctx, nilSchema.ID, "noschema")
	if err != nil {
		t.Fatalf("GetUserBlock (NULL schema): %v", err)
	}
	if gotNil.ConfigSchema != nil {
		t.Errorf("round-tripped NULL ConfigSchema = %q, want nil", gotNil.ConfigSchema)
	}

	if _, err := s.GetUserBlock(ctx, realm.ID, "ghost"); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown get = %v, want ErrNotFound", err)
	}
	got, err := s.GetUserBlock(ctx, realm.ID, "b1")
	if err != nil {
		t.Fatalf("GetUserBlock: %v", err)
	}
	if got.RealmID != realm.ID || string(got.Source) != string(src) || string(got.ConfigSchema) != `{"type":"object"}` {
		t.Errorf("round-trip = %+v", got)
	}

	newSrc := []byte("def handle_message(msg): return msg")
	updated, err := s.UpdateUserBlock(ctx, realm.ID, "b1", "producer_consumer", newSrc, nil)
	if err != nil {
		t.Fatalf("UpdateUserBlock: %v", err)
	}
	if updated.BlockType != "producer_consumer" || string(updated.Source) != string(newSrc) || updated.ConfigSchema != nil {
		t.Errorf("updated block = %+v", updated)
	}
	if updated.ID != b.ID || updated.RealmID != realm.ID || updated.Name != "b1" {
		t.Errorf("update touched identity: was %+v, now %+v", b, updated)
	}
	if !updated.CreatedAt.Equal(b.CreatedAt) {
		t.Errorf("CreatedAt bumped: created=%v updated=%v", b.CreatedAt, updated.CreatedAt)
	}
	if updated.UpdatedAt.Before(updated.CreatedAt) {
		t.Errorf("UpdatedAt not advanced: created=%v updated=%v", updated.CreatedAt, updated.UpdatedAt)
	}
	if _, err := s.UpdateUserBlock(ctx, realm.ID, "ghost", "consumer", src, nil); !errors.Is(err, ErrNotFound) {
		t.Errorf("update unknown = %v, want ErrNotFound", err)
	}

	other := mustCreateRealm(t, s)
	if _, err := s.CreateUserBlock(ctx, other.ID, &UserBlock{Name: "elsewhere", BlockType: "producer", Source: src}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CreateUserBlock(ctx, realm.ID, &UserBlock{Name: "b2", BlockType: "consumer", Source: src}); err != nil {
		t.Fatal(err)
	}
	list, err := s.ListUserBlocks(ctx, realm.ID)
	if err != nil {
		t.Fatalf("ListUserBlocks: %v", err)
	}
	if len(list) != 2 || list[0].Name != "b1" || list[1].Name != "b2" {
		t.Errorf("list = %+v, want [b1 b2]", list)
	}
	otherList, err := s.ListUserBlocks(ctx, other.ID)
	if err != nil {
		t.Fatalf("ListUserBlocks (other realm): %v", err)
	}
	if len(otherList) != 1 || otherList[0].Name != "elsewhere" {
		t.Errorf("other-realm list = %+v, want [elsewhere]", otherList)
	}

	gotAfterUpdate, err := s.GetUserBlock(ctx, realm.ID, "b1")
	if err != nil {
		t.Fatal(err)
	}
	if gotAfterUpdate.BlockType != "producer_consumer" {
		t.Errorf("post-update get = %+v, want producer_consumer", gotAfterUpdate)
	}

	if err := s.DeleteUserBlock(ctx, realm.ID, "b2"); err != nil {
		t.Fatalf("DeleteUserBlock: %v", err)
	}
	if _, err := s.GetUserBlock(ctx, realm.ID, "b2"); !errors.Is(err, ErrNotFound) {
		t.Errorf("get after delete = %v, want ErrNotFound", err)
	}
	if err := s.DeleteUserBlock(ctx, realm.ID, "b2"); !errors.Is(err, ErrNotFound) {
		t.Errorf("delete of missing = %v, want ErrNotFound", err)
	}
}
