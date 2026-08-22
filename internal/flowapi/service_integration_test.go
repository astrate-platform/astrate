//go:build integration

package flowapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
)

var realmSeq int64

// uniqueTestRealm returns a schema-valid realm name unique per test run.
func uniqueTestRealm() string {
	realmSeq++
	return "t" + strconv.FormatInt(time.Now().UnixNano(), 36) + strconv.Itoa(int(realmSeq))
}

// newTestService wires a real store + manager + default block registry
// against the shared test TimescaleDB.
func newTestService(t *testing.T) (*Service, *store.Store, string) {
	t.Helper()
	pool := testutil.StartTimescale(t)
	ctx := context.Background()
	st, err := store.New(ctx, pool.Config().ConnString())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(st.Close)

	mgr := flow.NewManager()
	reg := blocks.DefaultRegistry()
	bus := stream.New(nil)
	svc := NewService(st, mgr, reg, bus, slog.New(slog.DiscardHandler))

	sealed := make([]byte, 48)
	for i := range sealed {
		sealed[i] = byte(i)
	}
	realm, err := st.CreateRealm(ctx, store.NewRealm{
		Name:               uniqueTestRealm(),
		JWTPublicKeysPEM:   []string{"-----BEGIN PUBLIC KEY-----\nplaceholder\n-----END PUBLIC KEY-----\n"},
		CACertificatePEM:   "-----BEGIN CERTIFICATE-----\nplaceholder\n-----END CERTIFICATE-----\n",
		CAPrivateKeySealed: sealed,
	})
	if err != nil {
		t.Fatalf("CreateRealm: %v", err)
	}
	return svc, st, realm.Name
}

func mustStartFlow(t *testing.T, svc *Service, realm, name, pipeline string, config json.RawMessage) *FlowView {
	t.Helper()
	view, err := svc.CreateAndStartFlow(context.Background(), realm, CreateFlowRequest{
		Name: name, Pipeline: pipeline, Config: config, AutoRestart: true,
	})
	if err != nil {
		t.Fatalf("CreateAndStartFlow: %v", err)
	}
	t.Cleanup(func() { _ = svc.DeleteFlow(context.Background(), realm, name) })
	return view
}

// TestReloadAndConfigUpdate covers #44 (explicit reload picks up an edited
// pipeline; update alone does not touch running flows) and #46 (config-only
// update rebuilds a live flow).
func TestReloadAndConfigUpdate(t *testing.T) {
	svc, st, realm := newTestService(t)
	ctx := context.Background()

	def := []byte(`{"blocks":[
		{"name":"src","block_type":"astarte_source","config":{"interface":"${config.iface}"}},
		{"name":"sink","block_type":"null_sink"}
	],"connections":[{"from":"src","to":"sink"}]}`)
	if _, err := svc.CreatePipeline(ctx, realm, "p1", def); err != nil {
		t.Fatalf("CreatePipeline: %v", err)
	}

	started := mustStartFlow(t, svc, realm, "f1", "p1", json.RawMessage(`{"iface":"dev.v1"}`))
	if started.Status != "running" {
		t.Fatalf("start status = %q", started.Status)
	}
	firstRuntime := started.RuntimeID

	// Editing the pipeline is a no-op for the running instance...
	def2 := []byte(`{"blocks":[
		{"name":"src","block_type":"astarte_source","config":{"interface":"${config.iface}"}},
		{"name":"mid","block_type":"filter","config":{"key_prefix":"a"}},
		{"name":"sink","block_type":"log_sink"}
	],"connections":[{"from":"src","to":"mid"},{"from":"mid","to":"sink"}]}`)
	upd, err := svc.UpdatePipeline(ctx, realm, "p1", def2)
	if err != nil {
		t.Fatalf("UpdatePipeline: %v", err)
	}
	if len(upd.ReferencingFlows) != 1 || upd.ReferencingFlows[0] != "f1" {
		t.Errorf("referencing_flows = %v, want [f1]", upd.ReferencingFlows)
	}

	// ...and reload picks the new definition up (new runtime id).
	reloaded, err := svc.ReloadFlow(ctx, realm, "f1")
	if err != nil {
		t.Fatalf("ReloadFlow: %v", err)
	}
	if reloaded.Status != "running" {
		t.Fatalf("reload status = %q err=%v", reloaded.Status, reloaded.ErrorMessage)
	}
	if reloaded.RuntimeID == "" || reloaded.RuntimeID == firstRuntime {
		t.Errorf("reload runtime_id = %q, want a fresh instance", reloaded.RuntimeID)
	}

	// Config update on a live flow validates against the pipeline first:
	// missing ${config.iface} key fails loudly and leaves config untouched.
	if _, err := svc.UpdateFlowConfig(ctx, realm, "f1", json.RawMessage(`{"other":"x"}`)); !errors.Is(err, ErrValidation) {
		t.Errorf("missing key = %v, want ErrValidation", err)
	}
	got, _ := svc.GetFlow(ctx, realm, "f1")
	if !jsonEqual(got.Config, []byte(`{"iface":"dev.v1"}`)) {
		t.Errorf("config mutated by failed update: %s", got.Config)
	}

	// Valid config update persists and rebuilds the live graph.
	updView, err := svc.UpdateFlowConfig(ctx, realm, "f1", json.RawMessage(`{"iface":"dev.v2"}`))
	if err != nil {
		t.Fatalf("UpdateFlowConfig: %v", err)
	}
	if updView.Status != "running" || updView.RuntimeID == reloaded.RuntimeID {
		t.Errorf("after config update: status=%q runtime=%q (prev %q)", updView.Status, updView.RuntimeID, reloaded.RuntimeID)
	}
	if !jsonEqual(updView.Config, []byte(`{"iface":"dev.v2"}`)) {
		t.Errorf("persisted config = %s", updView.Config)
	}

	// Config update on a durably-stopped flow just persists (no rebuild).
	// Simulate a clean stop: tear the live graph down and flip the row.
	if err := svc.Manager().StopFlow(ctx, flow.InstanceID(realm, "f1")); err != nil {
		t.Fatal(err)
	}
	svc.Manager().UnregisterFlow(flow.InstanceID(realm, "f1"))
	r, err := st.GetRealmByName(ctx, realm)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	if err := st.UpdateFlowRuntime(ctx, r.ID, "f1", "stopped", nil, nil, &now); err != nil {
		t.Fatal(err)
	}
	stopped, err := svc.UpdateFlowConfig(ctx, realm, "f1", json.RawMessage(`{"iface":"dev.v3"}`))
	if err != nil {
		t.Fatalf("UpdateFlowConfig stopped flow: %v", err)
	}
	if stopped.Status != "stopped" {
		t.Errorf("stopped flow status after config update = %q", stopped.Status)
	}
	if !jsonEqual(stopped.Config, []byte(`{"iface":"dev.v3"}`)) {
		t.Errorf("stopped flow config = %s, want {}", stopped.Config)
	}
}

// TestReloadFailedFlow covers the manual-restart escape hatch: a failed flow
// can be brought back via reload once the pipeline exists again.
func TestReloadFailedFlow(t *testing.T) {
	svc, st, realm := newTestService(t)
	ctx := context.Background()

	// Flow referencing a pipeline that doesn't exist yet -> start fails.
	if _, err := svc.CreateAndStartFlow(ctx, realm, CreateFlowRequest{
		Name: "f2", Pipeline: "ghost-pipeline", AutoRestart: false,
	}); err == nil {
		t.Fatal("expected start failure for unknown pipeline")
	}
	row, err := st.GetFlow(ctx, func() int16 {
		r, rerr := st.GetRealmByName(ctx, realm)
		if rerr != nil {
			t.Fatal(rerr)
		}
		return r.ID
	}(), "f2")
	if err != nil {
		t.Fatal(err)
	}
	if row.Status != "failed" {
		t.Fatalf("status = %q, want failed", row.Status)
	}

	// Store the pipeline under the same name, then reload succeeds.
	def := []byte(`{"blocks":[
		{"name":"src","block_type":"astarte_source"},
		{"name":"sink","block_type":"null_sink"}
	],"connections":[{"from":"src","to":"sink"}]}`)
	if _, err := svc.CreatePipeline(ctx, realm, "ghost-pipeline", def); err != nil {
		t.Fatal(err)
	}
	view, err := svc.ReloadFlow(ctx, realm, "f2")
	if err != nil {
		t.Fatalf("ReloadFlow after pipeline created: %v", err)
	}
	if view.Status != "running" {
		t.Fatalf("reloaded failed flow status = %q", view.Status)
	}
	t.Cleanup(func() { _ = svc.DeleteFlow(ctx, realm, "f2") })
}

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
