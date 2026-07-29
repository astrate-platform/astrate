package flow_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks/astartesource"
)

// TestSourcePump_AstarteSourceToSink verifies issue #37: the manager pumps
// an AstarteSource and delivers converted messages to a downstream sink.
func TestSourcePump_AstarteSourceToSink(t *testing.T) {
	bus := stream.New(nil)
	src := astartesource.New(bus, astartesource.Config{Realm: "test"})
	sink := &countingSink{}

	mgr := flow.NewManager()
	f, err := mgr.StartFlow(context.Background(), flow.FlowConfig{
		PipelineID: "pump-1",
		Blocks:     []flow.Block{src, sink},
		RouterCfg: flow.RouterConfig{
			Lanes:        2,
			LaneCapacity: 64,
		},
	})
	if err != nil {
		t.Fatalf("StartFlow: %v", err)
	}

	ts := time.Date(2026, 7, 29, 12, 0, 0, 0, time.UTC)
	bus.Publish(stream.Event{
		Kind:      stream.KindIncomingData,
		Realm:     "test",
		DeviceID:  "dev1",
		Interface: "com.ex.Sensors",
		Path:      "/temp",
		Value:     21.5,
		Timestamp: ts,
	})

	deadline := time.Now().Add(2 * time.Second)
	for sink.load() == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if got := sink.load(); got != 1 {
		t.Fatalf("sink received %d messages, want 1", got)
	}
	got := sink.lastMsg()
	if got == nil {
		t.Fatal("sink last message is nil")
	}
	if got.Key != "test/dev1" {
		t.Errorf("Key = %q, want test/dev1", got.Key)
	}
	if got.Type != flow.TypeReal || got.Data.(float64) != 21.5 {
		t.Errorf("Type/Data = %v/%v, want TypeReal/21.5", got.Type, got.Data)
	}

	drainCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := mgr.StopFlow(drainCtx, "pump-1"); err != nil {
		t.Fatalf("StopFlow: %v", err)
	}
	if f.Status() != flow.FlowStatusStopped {
		t.Errorf("status = %v, want stopped", f.Status())
	}
}

// TestSourcePump_StopReleasesSubscription verifies that StopFlow calls Stop
// on AstarteSource and drops the bus subscription.
func TestSourcePump_StopReleasesSubscription(t *testing.T) {
	bus := stream.New(nil)
	src := astartesource.New(bus, astartesource.Config{Realm: "test"})
	if bus.Subscribers() != 1 {
		t.Fatalf("Subscribers before start = %d, want 1", bus.Subscribers())
	}

	mgr := flow.NewManager()
	_, err := mgr.StartFlow(context.Background(), flow.FlowConfig{
		PipelineID: "pump-stop",
		Blocks:     []flow.Block{src, &countingSink{}},
	})
	if err != nil {
		t.Fatalf("StartFlow: %v", err)
	}
	if bus.Subscribers() != 1 {
		t.Fatalf("Subscribers while running = %d, want 1", bus.Subscribers())
	}

	drainCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := mgr.StopFlow(drainCtx, "pump-stop"); err != nil {
		t.Fatalf("StopFlow: %v", err)
	}
	if bus.Subscribers() != 0 {
		t.Fatalf("Subscribers after StopFlow = %d, want 0", bus.Subscribers())
	}
}

// TestSourcePump_FuncSourceBlock verifies NewSourceBlock is pumped the same way.
func TestSourcePump_FuncSourceBlock(t *testing.T) {
	var emitted atomic.Bool
	src := flow.NewSourceBlock("tick", func() ([]*flow.FlowMessage, error) {
		if emitted.Swap(true) {
			// Subsequent polls: nothing more (non-blocking SourceFunc).
			time.Sleep(50 * time.Millisecond)
			return nil, nil
		}
		return []*flow.FlowMessage{{
			Key:       "k",
			Type:      flow.TypeInteger,
			Data:      int64(7),
			Timestamp: time.Now().UnixMicro(),
		}}, nil
	})
	sink := &countingSink{}

	mgr := flow.NewManager()
	_, err := mgr.StartFlow(context.Background(), flow.FlowConfig{
		PipelineID: "pump-fn",
		Blocks:     []flow.Block{src, sink},
	})
	if err != nil {
		t.Fatalf("StartFlow: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = mgr.StopFlow(ctx, "pump-fn")
	})

	deadline := time.Now().Add(2 * time.Second)
	for sink.load() == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if got := sink.load(); got != 1 {
		t.Fatalf("sink received %d messages, want 1", got)
	}
}

// TestGraph_RunSkipsSource ensures submitted messages do not re-enter Source
// Process (which would discard the payload for AstarteSource).
func TestGraph_RunSkipsSource(t *testing.T) {
	var sourceCalls atomic.Int64
	src := flow.NewSourceBlock("src", func() ([]*flow.FlowMessage, error) {
		sourceCalls.Add(1)
		return []*flow.FlowMessage{{Key: "from-source", Type: flow.TypeString, Data: "x"}}, nil
	})
	sink := &countingSink{}

	g, err := flow.NewBlockGraph(src, sink)
	if err != nil {
		t.Fatal(err)
	}
	msg := &flow.FlowMessage{Key: "submitted", Type: flow.TypeInteger, Data: int64(1)}
	if _, err := g.Run(msg); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if sourceCalls.Load() != 0 {
		t.Fatalf("Source.Process called %d times during Run, want 0", sourceCalls.Load())
	}
	if sink.load() != 1 {
		t.Fatalf("sink got %d, want 1", sink.load())
	}
	if sink.lastMsg().Key != "submitted" {
		t.Errorf("sink key = %q, want submitted", sink.lastMsg().Key)
	}
}

// countingSink is a flow.Block sink used by pump tests.
type countingSink struct {
	n    atomic.Int64
	last atomic.Pointer[flow.FlowMessage]
}

func (c *countingSink) Process(msg *flow.FlowMessage) ([]*flow.FlowMessage, error) {
	if msg != nil {
		cp := *msg
		c.last.Store(&cp)
		c.n.Add(1)
	}
	return nil, nil
}

func (c *countingSink) Name() string { return "counting-sink" }

func (c *countingSink) load() int64 { return c.n.Load() }

func (c *countingSink) lastMsg() *flow.FlowMessage { return c.last.Load() }
