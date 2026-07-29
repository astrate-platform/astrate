//go:build integration

package store

import (
	"context"
	"errors"
	"testing"
)

func testPipelines(t *testing.T, s *Store) {
	ctx := context.Background()
	realm := mustCreateRealm(t, s)

	def := []byte(`{"blocks":[{"name":"src"},{"name":"sink"}],"connections":[{"from":"src","to":"sink"}]}`)
	p, err := s.CreatePipeline(ctx, realm.ID, "p1", def)
	if err != nil {
		t.Fatalf("CreatePipeline: %v", err)
	}
	if p.ID == 0 || p.Name != "p1" {
		t.Errorf("created pipeline = %+v", p)
	}
	if _, err := s.CreatePipeline(ctx, realm.ID, "p1", def); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("duplicate create = %v, want ErrAlreadyExists", err)
	}

	cyclic := []byte(`{"blocks":[{"name":"a"},{"name":"b"}],"connections":[{"from":"a","to":"b"},{"from":"b","to":"a"}]}`)
	if _, err := s.CreatePipeline(ctx, realm.ID, "cyclic", cyclic); !errors.Is(err, ErrPipelineCyclic) {
		t.Errorf("cyclic create = %v, want ErrPipelineCyclic", err)
	}

	if _, err := s.GetPipeline(ctx, realm.ID, "ghost"); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown get = %v, want ErrNotFound", err)
	}
	got, err := s.GetPipeline(ctx, realm.ID, "p1")
	if err != nil {
		t.Fatalf("GetPipeline: %v", err)
	}
	if got.RealmID != realm.ID {
		t.Errorf("got.RealmID = %v, want %v", got.RealmID, realm.ID)
	}

	def2 := []byte(`{"blocks":[{"name":"src"},{"name":"mid"},{"name":"sink"}],"connections":[{"from":"src","to":"mid"},{"from":"mid","to":"sink"}]}`)
	updated, err := s.UpdatePipeline(ctx, realm.ID, "p1", def2)
	if err != nil {
		t.Fatalf("UpdatePipeline: %v", err)
	}
	if !updated.UpdatedAt.After(p.CreatedAt) && !updated.UpdatedAt.Equal(p.CreatedAt) {
		t.Errorf("UpdatedAt not advanced: created=%v updated=%v", p.CreatedAt, updated.UpdatedAt)
	}
	if _, err := s.UpdatePipeline(ctx, realm.ID, "ghost", def2); !errors.Is(err, ErrNotFound) {
		t.Errorf("update unknown = %v, want ErrNotFound", err)
	}

	other := mustCreateRealm(t, s)
	if _, err := s.CreatePipeline(ctx, other.ID, "elsewhere", def); err != nil {
		t.Fatal(err)
	}
	if _, err := s.CreatePipeline(ctx, realm.ID, "p2", def); err != nil {
		t.Fatal(err)
	}
	list, err := s.ListPipelines(ctx, realm.ID)
	if err != nil {
		t.Fatalf("ListPipelines: %v", err)
	}
	if len(list) != 2 || list[0].Name != "p1" || list[1].Name != "p2" {
		t.Errorf("list = %+v, want [p1 p2]", list)
	}

	if err := s.DeletePipeline(ctx, realm.ID, "p1"); err != nil {
		t.Fatalf("DeletePipeline: %v", err)
	}
	if err := s.DeletePipeline(ctx, realm.ID, "p1"); !errors.Is(err, ErrNotFound) {
		t.Errorf("double delete = %v, want ErrNotFound", err)
	}
}
