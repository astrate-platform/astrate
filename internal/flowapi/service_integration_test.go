//go:build integration

package flowapi

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
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
	svc := NewService(st, mgr, reg, bus, nil, slog.New(slog.DiscardHandler))

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
	if err := st.UpdateFlowRuntime(ctx, r.ID, "f1", "stopped", nil, nil, nil, &now); err != nil {
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

// --- #45 phase 1: block-death detection + auto-restart ----------------------

// boomOnce tracks, per flow name, whether the doomed block already died.
var boomOnce sync.Map

// boomConstructor builds a single-node "source+sink" pipeline block whose
// container-equivalent dies shortly after start — exactly once per flow.
func boomConstructor(name string, _ map[string]any, deps flow.Deps) (flow.Block, error) {
	v, _ := boomOnce.LoadOrStore(deps.FlowName, new(atomic.Bool))
	fired := v.(*atomic.Bool)
	if deps.NotifyFatal != nil {
		go func() {
			time.Sleep(100 * time.Millisecond)
			if fired.CompareAndSwap(false, true) {
				deps.NotifyFatal(name, errors.New("boom: container exited unexpectedly"))
			}
		}()
	}
	return flow.NewSinkBlock(name, func(*flow.Message) error { return nil }), nil
}

const boomPipeline = `{"blocks":[{"name":"n1","block_type":"boom"}],"connections":[]}`

// TestBlockDeathFailsFlow covers the non-auto-restart path: a dead block
// fails the whole flow with failed_block recorded and no live instance left.
func TestBlockDeathFailsFlow(t *testing.T) {
	svc, st, realm := newTestService(t)
	ctx := context.Background()

	reg := blocks.DefaultRegistry()
	reg.Register("boom", boomConstructor)
	svc2 := NewService(st, svc.Manager(), reg, stream.New(nil), nil, slog.New(slog.DiscardHandler))

	if _, err := svc2.CreatePipeline(ctx, realm, "boom-pipeline", []byte(boomPipeline)); err != nil {
		t.Fatal(err)
	}
	r, err := st.GetRealmByName(ctx, realm)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc2.CreateAndStartFlow(ctx, realm, CreateFlowRequest{
		Name: "bfatal", Pipeline: "boom-pipeline", AutoRestart: false,
	}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc2.DeleteFlow(ctx, realm, "bfatal") })

	deadline := time.Now().Add(5 * time.Second)
	for {
		view, err := svc2.GetFlow(ctx, realm, "bfatal")
		if err != nil {
			t.Fatal(err)
		}
		if view.Status == "failed" {
			if view.FailedBlock == nil || *view.FailedBlock != "n1" {
				t.Fatalf("failed_block = %v, want n1", view.FailedBlock)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("flow never failed; status=%q", view.Status)
		}
		time.Sleep(50 * time.Millisecond)
	}

	// No live graph may survive the death.
	if _, err := svc.Manager().GetFlowStatus(flow.InstanceID(realm, "bfatal")); err == nil {
		t.Fatal("live instance still registered after fatal")
	}
	_ = r
}

// TestBlockDeathAutoRestart covers the auto_restart path: the doomed block
// dies once, the service rebuilds with backoff, and the flow ends up running.
func TestBlockDeathAutoRestart(t *testing.T) {
	svc, st, realm := newTestService(t)
	ctx := context.Background()

	reg := blocks.DefaultRegistry()
	reg.Register("boom", boomConstructor)
	svc2 := NewService(st, svc.Manager(), reg, stream.New(nil), nil, slog.New(slog.DiscardHandler))

	if _, err := svc2.CreatePipeline(ctx, realm, "boom-pipeline", []byte(boomPipeline)); err != nil {
		t.Fatal(err)
	}
	if _, err := svc2.CreateAndStartFlow(ctx, realm, CreateFlowRequest{
		Name: "bretry", Pipeline: "boom-pipeline", AutoRestart: true,
	}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc2.DeleteFlow(ctx, realm, "bretry") })

	deadline := time.Now().Add(15 * time.Second)
	for {
		view, err := svc2.GetFlow(ctx, realm, "bretry")
		if err != nil {
			t.Fatal(err)
		}
		if view.Status == "running" && view.FailedBlock == nil {
			return // died once, came back clean
		}
		if time.Now().After(deadline) {
			t.Fatalf("flow never recovered: status=%q failed_block=%v", view.Status, view.FailedBlock)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
