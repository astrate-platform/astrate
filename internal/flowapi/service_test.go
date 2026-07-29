package flowapi

import (
	"testing"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
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

func TestListFlowsFiltersByRealm(t *testing.T) {
	mgr := flow.NewManager()
	// Register two flows via StartFlow with distinct pipeline IDs.
	// Empty transform→sink graphs still need a non-source block only? Need ≥1 block.
	// Use a single sink as a degenerate graph for list filtering.
	sink := flow.NewSinkBlock("s", func(*flow.FlowMessage) error { return nil })
	// BlockGraph requires at least one block; a lone sink is OK for list tests.
	if _, err := mgr.StartFlow(t.Context(), flow.FlowConfig{
		PipelineID: "acme/p1",
		Blocks:     []flow.Block{sink},
	}); err != nil {
		t.Fatal(err)
	}
	sink2 := flow.NewSinkBlock("s2", func(*flow.FlowMessage) error { return nil })
	if _, err := mgr.StartFlow(t.Context(), flow.FlowConfig{
		PipelineID: "other/p1",
		Blocks:     []flow.Block{sink2},
	}); err != nil {
		t.Fatal(err)
	}

	s := &Service{mgr: mgr, reg: blocks.DefaultRegistry()}
	list := s.ListFlows("acme")
	if len(list) != 1 || list[0].Pipeline != "p1" || list[0].Realm != "acme" {
		t.Fatalf("ListFlows acme = %+v", list)
	}
	if err := mgr.Shutdown(t.Context()); err != nil {
		t.Fatal(err)
	}
}
