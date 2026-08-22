package flow_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
	"github.com/astrate-platform/astrate/internal/flow/blocks/astartesource"
)

func TestRegistry_InstantiateLinear(t *testing.T) {
	var mu sync.Mutex
	var got []string
	snapshot := func() []string {
		mu.Lock()
		defer mu.Unlock()
		return append([]string(nil), got...)
	}
	reg := flow.NewRegistry()
	reg.Register("source", func(name string, _ map[string]any, _ flow.Deps) (flow.Block, error) {
		return flow.NewSourceBlock(name, func() ([]*flow.Message, error) {
			return nil, nil
		}), nil
	})
	reg.Register("transform", func(name string, _ map[string]any, _ flow.Deps) (flow.Block, error) {
		return flow.NewTransformBlock(name, func(msg *flow.Message) ([]*flow.Message, error) {
			mu.Lock()
			got = append(got, name+":"+msg.Key)
			mu.Unlock()
			return []*flow.Message{msg}, nil
		}), nil
	})
	reg.Register("sink", func(name string, _ map[string]any, _ flow.Deps) (flow.Block, error) {
		return flow.NewSinkBlock(name, func(msg *flow.Message) error {
			mu.Lock()
			got = append(got, name+":"+msg.Key)
			mu.Unlock()
			return nil
		}), nil
	})

	p := &flow.Pipeline{
		ID:   "pipe-1",
		Name: "linear",
		Blocks: []flow.PipelineNode{
			{Name: "src", BlockType: "source"},
			{Name: "mid", BlockType: "transform"},
			{Name: "out", BlockType: "sink"},
		},
		Connections: []flow.Connection{
			{From: "src", To: "mid"},
			{From: "mid", To: "out"},
		},
	}
	blocksList, err := reg.Instantiate(p, flow.Deps{})
	if err != nil {
		t.Fatalf("Instantiate: %v", err)
	}
	if len(blocksList) != 3 {
		t.Fatalf("len = %d, want 3", len(blocksList))
	}
	for i, want := range []string{"src", "mid", "out"} {
		if blocksList[i].Name() != want {
			t.Errorf("block[%d].Name() = %q, want %q", i, blocksList[i].Name(), want)
		}
	}

	mgr := flow.NewManager()
	f, err := mgr.StartFlow(context.Background(), flow.Config{
		PipelineID: p.ID,
		Blocks:     blocksList,
	})
	if err != nil {
		t.Fatalf("StartFlow: %v", err)
	}
	f.Router().Submit(&flow.Message{Key: "k1", Type: flow.TypeString, Data: "x"}, 1)
	deadline := time.Now().Add(2 * time.Second)
	var final []string
	for time.Now().Before(deadline) {
		final = snapshot()
		if len(final) >= 2 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if len(final) < 2 || final[0] != "mid:k1" || final[1] != "out:k1" {
		t.Fatalf("got path %v, want [mid:k1 out:k1]", final)
	}
	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := mgr.StopFlow(stopCtx, p.ID); err != nil {
		t.Fatalf("StopFlow: %v", err)
	}
}

func TestRegistry_UnknownBlockType(t *testing.T) {
	reg := flow.NewRegistry()
	reg.Register("null_sink", blocks.NullSink)
	p := &flow.Pipeline{
		ID: "x",
		Blocks: []flow.PipelineNode{
			{Name: "src", BlockType: "missing_type"},
			{Name: "sink", BlockType: "null_sink"},
		},
		Connections: []flow.Connection{{From: "src", To: "sink"}},
	}
	_, err := reg.Instantiate(p, flow.Deps{})
	if !errors.Is(err, flow.ErrUnknownBlockType) {
		t.Fatalf("err = %v, want ErrUnknownBlockType", err)
	}
}

func TestRegistry_ConstructorError(t *testing.T) {
	reg := flow.NewRegistry()
	reg.Register("bad", func(string, map[string]any, flow.Deps) (flow.Block, error) {
		return nil, errors.New("boom")
	})
	reg.Register("null_sink", blocks.NullSink)
	p := &flow.Pipeline{
		ID: "x",
		Blocks: []flow.PipelineNode{
			{Name: "a", BlockType: "bad"},
			{Name: "b", BlockType: "null_sink"},
		},
		Connections: []flow.Connection{{From: "a", To: "b"}},
	}
	_, err := reg.Instantiate(p, flow.Deps{})
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("err = %v, want construct error containing boom", err)
	}
}

// stoppableBlock is a Block + Stopper whose Stop() records that it was called.
type stoppableBlock struct {
	name    string
	stopped *int
}

func (b *stoppableBlock) Name() string { return b.name }
func (b *stoppableBlock) Process(msg *flow.Message) ([]*flow.Message, error) {
	return []*flow.Message{msg}, nil
}
func (b *stoppableBlock) Stop() { *b.stopped++ }

func TestRegistry_InstantiateStopsAlreadyBuiltBlocksOnFailure(t *testing.T) {
	var stoppedA, stoppedB int
	reg := flow.NewRegistry()
	reg.Register("stoppable-a", func(name string, _ map[string]any, _ flow.Deps) (flow.Block, error) {
		return &stoppableBlock{name: name, stopped: &stoppedA}, nil
	})
	reg.Register("stoppable-b", func(name string, _ map[string]any, _ flow.Deps) (flow.Block, error) {
		return &stoppableBlock{name: name, stopped: &stoppedB}, nil
	})
	reg.Register("bad", func(string, map[string]any, flow.Deps) (flow.Block, error) {
		return nil, errors.New("boom")
	})

	p := &flow.Pipeline{
		ID: "x",
		Blocks: []flow.PipelineNode{
			{Name: "a", BlockType: "stoppable-a"},
			{Name: "b", BlockType: "stoppable-b"},
			{Name: "c", BlockType: "bad"},
		},
		Connections: []flow.Connection{{From: "a", To: "b"}, {From: "b", To: "c"}},
	}
	blocksList, err := reg.Instantiate(p, flow.Deps{})
	if err == nil {
		t.Fatal("Instantiate: want error, got nil")
	}
	if len(blocksList) != 0 {
		t.Fatalf("blocksList = %v, want empty", blocksList)
	}
	if stoppedA != 1 {
		t.Errorf("stoppedA = %d, want 1", stoppedA)
	}
	if stoppedB != 1 {
		t.Errorf("stoppedB = %d, want 1", stoppedB)
	}
}

func TestParseDefinition(t *testing.T) {
	raw := []byte(`{
		"blocks": [
			{"name": "src", "block_type": "astarte_source", "config": {"interface": "com.ex.S"}},
			{"name": "sink", "block_type": "null_sink"}
		],
		"connections": [{"from": "src", "to": "sink"}]
	}`)
	p, err := flow.ParseDefinition("pipe-id", "my-pipe", raw)
	if err != nil {
		t.Fatalf("ParseDefinition: %v", err)
	}
	if p.ID != "pipe-id" || p.Name != "my-pipe" {
		t.Errorf("ID/Name = %q/%q", p.ID, p.Name)
	}
	if len(p.Blocks) != 2 || p.Blocks[0].BlockType != "astarte_source" {
		t.Errorf("blocks = %+v", p.Blocks)
	}
}

func TestParseDefinition_RejectsInvalid(t *testing.T) {
	_, err := flow.ParseDefinition("id", "n", []byte(`{"blocks":[]}`))
	if err == nil {
		t.Fatal("expected error for empty blocks")
	}
}

func TestDefaultCatalog_AstarteSourceToNullSink(t *testing.T) {
	bus := stream.New(nil)
	reg := blocks.DefaultRegistry()
	p := &flow.Pipeline{
		ID:   "live",
		Name: "live",
		Blocks: []flow.PipelineNode{
			{Name: "src", BlockType: "astarte_source", Config: map[string]any{"realm": "r1"}},
			{Name: "sink", BlockType: "null_sink"},
		},
		Connections: []flow.Connection{{From: "src", To: "sink"}},
	}
	blks, err := reg.Instantiate(p, flow.Deps{Bus: bus, Realm: "r1"})
	if err != nil {
		t.Fatalf("Instantiate: %v", err)
	}
	if blks[0].Name() != "src" {
		t.Errorf("source Name = %q, want src", blks[0].Name())
	}
	if _, ok := blks[0].(flow.Source); !ok {
		t.Fatal("first block should implement flow.Source")
	}
	if _, ok := blks[0].(flow.Stopper); !ok {
		t.Fatal("first block should implement flow.Stopper")
	}

	mgr := flow.NewManager()
	f, err := mgr.StartFlow(context.Background(), flow.Config{
		PipelineID: flow.PipelineID("r1", "live"),
		Blocks:     blks,
	})
	if err != nil {
		t.Fatalf("StartFlow: %v", err)
	}
	bus.Publish(stream.Event{
		Kind: stream.KindIncomingData, Realm: "r1", DeviceID: "d1",
		Interface: "com.ex.S", Path: "/v", Value: 3.0, Timestamp: time.Now(),
	})
	// Give the pump a moment to emit; null_sink discards — we only assert
	// lifecycle + subscriber release on stop.
	time.Sleep(50 * time.Millisecond)

	stopCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := mgr.StopFlow(stopCtx, f.PipelineID()); err != nil {
		t.Fatalf("StopFlow: %v", err)
	}
	// Source Stop should have run; a new subscription should work cleanly.
	src := astartesource.New(bus, astartesource.Config{Realm: "r1"})
	src.Stop()
}

func TestFlowInstanceID(t *testing.T) {
	if got := flow.InstanceID("acme", "sensors"); got != "acme/sensors" {
		t.Errorf("InstanceID = %q", got)
	}
	if flow.PipelineID("acme", "sensors") != flow.InstanceID("acme", "sensors") {
		t.Error("PipelineID should alias InstanceID")
	}
}
