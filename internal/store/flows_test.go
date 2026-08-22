//go:build integration

package store

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"
)

// jsonEqual compares two JSON blobs semantically (jsonb normalizes spacing).
func jsonEqual(a, b []byte) bool {
	var va, vb any
	if err := json.Unmarshal(a, &va); err != nil {
		return false
	}
	if err := json.Unmarshal(b, &vb); err != nil {
		return false
	}
	return reflect.DeepEqual(va, vb)
}

func testFlows(t *testing.T, s *Store) {
	ctx := context.Background()
	realm := mustCreateRealm(t, s)
	other := mustCreateRealm(t, s)

	cfg := []byte(`{"webhook_url":"https://example.com/hook"}`)
	f, err := s.CreateFlow(ctx, realm.ID, "prod-webhooks", "device-to-http", cfg, true)
	if err != nil {
		t.Fatalf("CreateFlow: %v", err)
	}
	if f.ID == 0 || f.Name != "prod-webhooks" || f.PipelineName != "device-to-http" {
		t.Errorf("created = %+v", f)
	}
	if !f.AutoRestart || f.Status != "stopped" {
		t.Errorf("defaults: auto_restart=%v status=%q", f.AutoRestart, f.Status)
	}
	if !jsonEqual(f.Config, cfg) {
		t.Errorf("config = %s, want %s", f.Config, cfg)
	}

	if _, err := s.CreateFlow(ctx, realm.ID, "prod-webhooks", "other", nil, false); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("duplicate create = %v, want ErrAlreadyExists", err)
	}

	// Empty config becomes {}.
	f2, err := s.CreateFlow(ctx, realm.ID, "batch", "device-to-http", nil, false)
	if err != nil {
		t.Fatalf("CreateFlow empty config: %v", err)
	}
	if string(f2.Config) != "{}" {
		t.Errorf("empty config stored as %q, want {}", f2.Config)
	}
	if f2.AutoRestart {
		t.Error("auto_restart should be false")
	}

	if _, err := s.GetFlow(ctx, realm.ID, "ghost"); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown get = %v, want ErrNotFound", err)
	}
	got, err := s.GetFlow(ctx, realm.ID, "prod-webhooks")
	if err != nil {
		t.Fatalf("GetFlow: %v", err)
	}
	if got.RealmID != realm.ID || got.PipelineName != "device-to-http" {
		t.Errorf("got = %+v", got)
	}

	now := time.Now().UTC().Truncate(time.Microsecond)
	errMsg := "pipeline not found"
	if err := s.UpdateFlowRuntime(ctx, realm.ID, "prod-webhooks", "failed", &errMsg, nil, nil); err != nil {
		t.Fatalf("UpdateFlowRuntime failed: %v", err)
	}
	got, err = s.GetFlow(ctx, realm.ID, "prod-webhooks")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "failed" || got.ErrorMessage == nil || *got.ErrorMessage != errMsg {
		t.Errorf("after fail: status=%q err=%v", got.Status, got.ErrorMessage)
	}

	if err := s.UpdateFlowRuntime(ctx, realm.ID, "prod-webhooks", "running", nil, &now, nil); err != nil {
		t.Fatalf("UpdateFlowRuntime running: %v", err)
	}
	got, err = s.GetFlow(ctx, realm.ID, "prod-webhooks")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != "running" || got.ErrorMessage != nil {
		t.Errorf("after run: status=%q err=%v", got.Status, got.ErrorMessage)
	}
	if got.StartedAt == nil {
		t.Error("started_at not set")
	}

	// Other realm isolation + list.
	if _, err := s.CreateFlow(ctx, other.ID, "prod-webhooks", "x", nil, true); err != nil {
		t.Fatal(err)
	}
	list, err := s.ListFlows(ctx, realm.ID)
	if err != nil {
		t.Fatalf("ListFlows: %v", err)
	}
	if len(list) != 2 || list[0].Name != "batch" || list[1].Name != "prod-webhooks" {
		t.Errorf("list = %+v, want [batch prod-webhooks]", list)
	}

	// Auto-restart list includes only auto_restart=true across realms.
	rehydrate, err := s.ListAutoRestartFlows(ctx)
	if err != nil {
		t.Fatalf("ListAutoRestartFlows: %v", err)
	}
	foundProd, foundBatch := false, false
	for _, fr := range rehydrate {
		if fr.RealmID == realm.ID && fr.Name == "prod-webhooks" {
			foundProd = true
			if fr.RealmName == "" {
				t.Error("RealmName empty on rehydrate row")
			}
		}
		if fr.RealmID == realm.ID && fr.Name == "batch" {
			foundBatch = true
		}
	}
	if !foundProd {
		t.Error("auto_restart=true flow missing from ListAutoRestartFlows")
	}
	if foundBatch {
		t.Error("auto_restart=false flow should not appear in ListAutoRestartFlows")
	}

	if err := s.DeleteFlow(ctx, realm.ID, "batch"); err != nil {
		t.Fatalf("DeleteFlow: %v", err)
	}
	if err := s.DeleteFlow(ctx, realm.ID, "batch"); !errors.Is(err, ErrNotFound) {
		t.Errorf("double delete = %v, want ErrNotFound", err)
	}
	if err := s.UpdateFlowRuntime(ctx, realm.ID, "ghost", "stopped", nil, nil, nil); !errors.Is(err, ErrNotFound) {
		t.Errorf("update unknown = %v, want ErrNotFound", err)
	}

	// UpdateFlowConfig replaces the snapshot and touches nothing else.
	newCfg := []byte(`{"webhook_url":"https://example.com/v2"}`)
	upd, err := s.UpdateFlowConfig(ctx, realm.ID, "prod-webhooks", newCfg)
	if err != nil {
		t.Fatalf("UpdateFlowConfig: %v", err)
	}
	if !jsonEqual(upd.Config, newCfg) {
		t.Errorf("updated config = %s, want %s", upd.Config, newCfg)
	}
	if upd.Status != "running" {
		t.Errorf("UpdateFlowConfig changed status to %q, want running (untouched)", upd.Status)
	}
	if _, err := s.UpdateFlowConfig(ctx, realm.ID, "ghost", newCfg); !errors.Is(err, ErrNotFound) {
		t.Errorf("UpdateFlowConfig unknown = %v, want ErrNotFound", err)
	}
}
