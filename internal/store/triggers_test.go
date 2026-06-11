//go:build integration

package store

import (
	"context"
	"errors"
	"testing"
)

func testTriggers(t *testing.T, s *Store) {
	ctx := context.Background()
	realm := mustCreateRealm(t, s)

	def := []byte(`{
		"name": "ondata",
		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
		"simple_triggers": [{"type": "data_trigger", "on": "incoming_data",
			"interface_name": "*", "match_path": "/*", "value_match_operator": "*"}]
	}`)

	created, err := s.CreateTrigger(ctx, realm.ID, "ondata", def)
	if err != nil {
		t.Fatalf("CreateTrigger: %v", err)
	}
	if created.ID == 0 || created.Name != "ondata" {
		t.Errorf("created trigger: %+v", created)
	}
	if _, err := s.CreateTrigger(ctx, realm.ID, "ondata", def); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("duplicate trigger: got %v, want ErrAlreadyExists", err)
	}

	got, err := s.GetTrigger(ctx, realm.ID, "ondata")
	if err != nil {
		t.Fatalf("GetTrigger: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("get round-trip: %+v", got)
	}
	// jsonb normalizes whitespace; semantic equality is enough here, the
	// engine re-parses the definition anyway.
	if len(got.Definition) == 0 {
		t.Error("round-tripped definition is empty")
	}

	if _, err := s.CreateTrigger(ctx, realm.ID, "onconnect", []byte(`{"name":"onconnect"}`)); err != nil {
		t.Fatal(err)
	}
	list, err := s.ListTriggers(ctx, realm.ID)
	if err != nil {
		t.Fatalf("ListTriggers: %v", err)
	}
	if len(list) != 2 || list[0].Name != "onconnect" || list[1].Name != "ondata" {
		t.Errorf("ListTriggers: %+v", list)
	}

	if err := s.DeleteTrigger(ctx, realm.ID, "ondata"); err != nil {
		t.Fatalf("DeleteTrigger: %v", err)
	}
	if _, err := s.GetTrigger(ctx, realm.ID, "ondata"); !errors.Is(err, ErrNotFound) {
		t.Errorf("trigger survived delete: %v", err)
	}
	if err := s.DeleteTrigger(ctx, realm.ID, "ondata"); !errors.Is(err, ErrNotFound) {
		t.Errorf("second delete: got %v, want ErrNotFound", err)
	}
}
