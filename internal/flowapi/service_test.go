package flowapi

import (
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
	"github.com/astrate-platform/astrate/internal/store"
)

func TestCheckBlockTypes(t *testing.T) {
	s := &Service{reg: blocks.DefaultRegistry()}
	ok := []byte(`{"blocks":[{"name":"src","block_type":"astarte_source"},{"name":"sink","block_type":"null_sink"}],"connections":[{"from":"src","to":"sink"}]}`)
	if err := s.checkBlockTypes(ok); err != nil {
		t.Fatalf("ok definition: %v", err)
	}
	bad := []byte(`{"blocks":[{"name":"src","block_type":"nope"},{"name":"sink","block_type":"null_sink"}]}`)
	if err := s.checkBlockTypes(bad); err == nil {
		t.Fatal("expected unknown block_type error")
	}
}

func TestListBlocks_DefaultCatalog(t *testing.T) {
	s := &Service{reg: blocks.DefaultRegistry()}
	list := s.ListBlocks("any-realm")
	want := map[string]bool{
		blocks.TypeAstarteSource: false,
		blocks.TypeFilter:        false,
		blocks.TypeMap:           false,
		blocks.TypeNullSink:      false,
		blocks.TypeLogSink:       false,
	}
	for _, info := range list {
		if _, ok := want[info.Type]; ok {
			want[info.Type] = true
		}
		if info.Summary == "" {
			t.Errorf("%q: empty summary", info.Type)
		}
	}
	for typ, seen := range want {
		if !seen {
			t.Errorf("missing type %q in ListBlocks", typ)
		}
	}
	got, err := s.GetBlock("any-realm", blocks.TypeFilter)
	if err != nil {
		t.Fatal(err)
	}
	if got.Role != blocks.RoleTransform || got.Config == "" {
		t.Fatalf("GetBlock filter = %+v", got)
	}
	if _, err := s.GetBlock("any-realm", "nope"); err == nil {
		t.Fatal("expected not found for unknown type")
	}
}

func TestMergeFlowView_LiveOverridesStatus(t *testing.T) {
	now := time.Now().UTC()
	row := &store.Flow{
		Name:         "prod",
		PipelineName: "device-to-http",
		Config:       []byte(`{"webhook_url":"https://x"}`),
		AutoRestart:  true,
		Status:       "stopped",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	// No live: durable status wins.
	v := mergeFlowView("acme", row, nil)
	if v.Name != "prod" || v.Pipeline != "device-to-http" || v.Realm != "acme" {
		t.Fatalf("view = %+v", v)
	}
	if v.Status != "stopped" || v.RuntimeID != "" {
		t.Fatalf("status/runtime = %q / %q", v.Status, v.RuntimeID)
	}
	if string(v.Config) != `{"webhook_url":"https://x"}` || !v.AutoRestart {
		t.Fatalf("config/auto = %s / %v", v.Config, v.AutoRestart)
	}

	// Live running overlays status + runtime_id.
	mgr := flow.NewManager()
	sink := flow.NewSinkBlock("s", func(*flow.FlowMessage) error { return nil })
	live, err := mgr.StartFlow(t.Context(), flow.FlowConfig{
		PipelineID: flow.FlowInstanceID("acme", "prod"),
		Blocks:     []flow.Block{sink},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = mgr.Shutdown(t.Context()) })

	v2 := mergeFlowView("acme", row, live)
	if v2.Status != "running" || v2.RuntimeID == "" {
		t.Fatalf("live merge = %+v", v2)
	}
	if v2.StartedAt == nil {
		t.Fatal("expected started_at from live")
	}
}

func TestSplitInstanceID(t *testing.T) {
	realm, name, ok := splitInstanceID("acme/prod-webhooks")
	if !ok || realm != "acme" || name != "prod-webhooks" {
		t.Fatalf("got %q %q %v", realm, name, ok)
	}
	if _, _, ok := splitInstanceID("noneslash"); ok {
		t.Fatal("expected fail")
	}
}

func TestFlowInstanceIDAlias(t *testing.T) {
	if flow.FlowInstanceID("r", "n") != flow.FlowPipelineID("r", "n") {
		t.Fatal("FlowPipelineID should alias FlowInstanceID")
	}
}
