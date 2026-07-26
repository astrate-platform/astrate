//go:build integration

package store

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func testPolicies(t *testing.T, s *Store) {
	ctx := context.Background()
	realm := mustCreateRealm(t, s)
	def := []byte(`{"name":"retry5xx","error_handlers":[{"on":"server_error","strategy":"retry"}],"maximum_capacity":100,"retry_times":3}`)

	p, err := s.CreateTriggerPolicy(ctx, realm.ID, "retry5xx", def)
	if err != nil {
		t.Fatalf("CreateTriggerPolicy: %v", err)
	}
	if p.ID == 0 || p.Name != "retry5xx" {
		t.Errorf("created policy = %+v", p)
	}
	if _, err := s.CreateTriggerPolicy(ctx, realm.ID, "retry5xx", def); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("duplicate create = %v, want ErrAlreadyExists", err)
	}

	// jsonb round-trips semantically, not byte-identically.
	got, err := s.GetTriggerPolicy(ctx, realm.ID, "retry5xx")
	if err != nil {
		t.Fatalf("GetTriggerPolicy: %v", err)
	}
	var wantDoc, gotDoc map[string]any
	if err := json.Unmarshal(def, &wantDoc); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(got.Definition, &gotDoc); err != nil {
		t.Fatalf("stored definition is not JSON: %v", err)
	}
	if !reflect.DeepEqual(gotDoc, wantDoc) {
		t.Errorf("stored definition = %v, want %v", gotDoc, wantDoc)
	}
	if _, err := s.GetTriggerPolicy(ctx, realm.ID, "ghost"); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown get = %v, want ErrNotFound", err)
	}

	// A second policy; list returns both, name-ordered, realm-scoped.
	if _, err := s.CreateTriggerPolicy(ctx, realm.ID, "discard", []byte(`{"name":"discard","error_handlers":[{"on":"any_error","strategy":"discard"}],"maximum_capacity":10}`)); err != nil {
		t.Fatal(err)
	}
	other := mustCreateRealm(t, s)
	if _, err := s.CreateTriggerPolicy(ctx, other.ID, "elsewhere", def); err != nil {
		t.Fatal(err)
	}
	list, err := s.ListTriggerPolicies(ctx, realm.ID)
	if err != nil {
		t.Fatalf("ListTriggerPolicies: %v", err)
	}
	if len(list) != 2 || list[0].Name != "discard" || list[1].Name != "retry5xx" {
		t.Errorf("list = %+v, want [discard retry5xx]", list)
	}

	if err := s.DeleteTriggerPolicy(ctx, realm.ID, "retry5xx"); err != nil {
		t.Fatalf("DeleteTriggerPolicy: %v", err)
	}
	if err := s.DeleteTriggerPolicy(ctx, realm.ID, "retry5xx"); !errors.Is(err, ErrNotFound) {
		t.Errorf("double delete = %v, want ErrNotFound", err)
	}
}
