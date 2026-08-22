package flow

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// TestManager_StartFlowTransitionToRunning verifies that StartFlow creates a
// flow with status running and that the flow is retrievable.
func TestManager_StartFlowTransitionToRunning(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	sink := &collectBlock{}
	f, err := mgr.StartFlow(ctx, Config{
		PipelineID: "pipe-1",
		Blocks:     []Block{sink},
	})
	if err != nil {
		t.Fatalf("StartFlow: %v", err)
	}
	t.Cleanup(func() { _ = mgr.StopFlow(context.Background(), "pipe-1") })

	if f == nil {
		t.Fatal("StartFlow returned nil flow")
	}
	if f.ID() == "" {
		t.Fatal("flow ID is empty")
	}
	if got := f.Status(); got != FlowStatusRunning {
		t.Fatalf("status = %v, want running", got)
	}
	if got := f.PipelineID(); got != "pipe-1" {
		t.Fatalf("pipeline ID = %q, want %q", got, "pipe-1")
	}
	if f.StartedAt().IsZero() {
		t.Fatal("started-at not set")
	}

	got, err := mgr.GetFlowStatus("pipe-1")
	if err != nil {
		t.Fatalf("GetFlowStatus: %v", err)
	}
	if got != FlowStatusRunning {
		t.Fatalf("GetFlowStatus = %v, want running", got)
	}
}

// TestManager_StopFlowDrainsAndReleases verifies that StopFlow drains all
// in-flight messages and transitions the flow to stopped.
func TestManager_StopFlowDrainsAndReleases(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	sink := &collectBlock{}
	f, err := mgr.StartFlow(ctx, Config{
		PipelineID: "pipe-2",
		Blocks:     []Block{sink},
		RouterCfg: RouterConfig{
			Lanes:        2,
			LaneCapacity: 128,
			QoS0Overflow: OverflowBlock,
			QoS1Overflow: OverflowBlock,
		},
	})
	if err != nil {
		t.Fatalf("StartFlow: %v", err)
	}

	const n = 50
	for i := range n {
		f.Router().Submit(makeMsg("key", i), 1)
	}

	drainCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := mgr.StopFlow(drainCtx, "pipe-2"); err != nil {
		t.Fatalf("StopFlow: %v", err)
	}

	got := sink.len()
	if got != n {
		t.Fatalf("delivered %d messages after StopFlow, want %d", got, n)
	}

	if got := f.Status(); got != FlowStatusStopped {
		t.Fatalf("status = %v, want stopped", got)
	}
	if f.StoppedAt().IsZero() {
		t.Fatal("stopped-at not set")
	}

	// Post-stop submits are silently dropped.
	f.Router().Submit(makeMsg("key", 999), 1)
	time.Sleep(50 * time.Millisecond)
	if got := sink.len(); got != n {
		t.Errorf("post-stop message delivered: %d", got)
	}
}

// TestManager_FailedInitLeavesNoMapEntry verifies that a graph construction
// failure (nil block) returns an error and does not register the instance.
func TestManager_FailedInitLeavesNoMapEntry(t *testing.T) {
	mgr := NewManager()

	f, err := mgr.StartFlow(context.Background(), Config{
		PipelineID: "pipe-3",
		Blocks:     []Block{nil},
	})
	if err == nil {
		t.Fatal("StartFlow with nil block should return error")
	}
	if f != nil {
		t.Fatal("StartFlow should not return a flow on construction failure")
	}

	if _, err := mgr.GetFlowStatus("pipe-3"); !errors.Is(err, ErrFlowNotFound) {
		t.Fatalf("GetFlowStatus = %v, want ErrFlowNotFound", err)
	}
	// Retry should not hit ErrFlowExists.
	sink := NewSinkBlock("s", func(*Message) error { return nil })
	f2, err := mgr.StartFlow(context.Background(), Config{
		PipelineID: "pipe-3",
		Blocks:     []Block{sink},
	})
	if err != nil {
		t.Fatalf("retry StartFlow: %v", err)
	}
	t.Cleanup(func() { _ = mgr.StopFlow(context.Background(), "pipe-3") })
	if f2.Status() != FlowStatusRunning {
		t.Fatalf("status = %v", f2.Status())
	}
}

// TestManager_ListFlows verifies that ListFlows returns all registered flows.
func TestManager_ListFlows(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	blk := &passthroughBlock{}
	for _, id := range []string{"a", "b", "c"} {
		_, err := mgr.StartFlow(ctx, Config{
			PipelineID: id,
			Blocks:     []Block{blk},
		})
		if err != nil {
			t.Fatalf("StartFlow %s: %v", id, err)
		}
	}
	t.Cleanup(func() { _ = mgr.StopFlow(context.Background(), "a") })
	t.Cleanup(func() { _ = mgr.StopFlow(context.Background(), "b") })
	t.Cleanup(func() { _ = mgr.StopFlow(context.Background(), "c") })

	flows := mgr.ListFlows()
	if got := len(flows); got != 3 {
		t.Fatalf("ListFlows returned %d flows, want 3", got)
	}
}

// TestManager_DuplicatePipelineReturnsError verifies that starting two
// flows with the same pipeline ID returns ErrFlowExists.
func TestManager_DuplicatePipelineReturnsError(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	blk := &passthroughBlock{}
	_, err := mgr.StartFlow(ctx, Config{
		PipelineID: "dup",
		Blocks:     []Block{blk},
	})
	if err != nil {
		t.Fatalf("StartFlow: %v", err)
	}
	t.Cleanup(func() { _ = mgr.StopFlow(context.Background(), "dup") })

	_, err = mgr.StartFlow(ctx, Config{
		PipelineID: "dup",
		Blocks:     []Block{blk},
	})
	if !errors.Is(err, ErrFlowExists) {
		t.Fatalf("duplicate StartFlow error = %v, want ErrFlowExists", err)
	}
}

// TestManager_GetFlowNotFound verifies that querying a non-existent flow
// returns ErrFlowNotFound.
func TestManager_GetFlowNotFound(t *testing.T) {
	mgr := NewManager()
	_, err := mgr.GetFlowStatus("nope")
	if !errors.Is(err, ErrFlowNotFound) {
		t.Fatalf("GetFlowStatus error = %v, want ErrFlowNotFound", err)
	}
}

// TestManager_StopFlowNotFound verifies that stopping a non-existent flow
// returns ErrFlowNotFound.
func TestManager_StopFlowNotFound(t *testing.T) {
	mgr := NewManager()
	err := mgr.StopFlow(context.Background(), "nope")
	if !errors.Is(err, ErrFlowNotFound) {
		t.Fatalf("StopFlow error = %v, want ErrFlowNotFound", err)
	}
}

// TestManager_StopFlowIsGraceful submits messages, stops the flow, and
// verifies all in-flight messages are processed before shutdown completes.
func TestManager_StopFlowIsGraceful(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	var mu sync.Mutex
	var order []int
	sink := &orderedCollectBlock{mu: &mu, order: &order}

	f, err := mgr.StartFlow(ctx, Config{
		PipelineID: "graceful",
		Blocks:     []Block{sink},
		RouterCfg: RouterConfig{
			Lanes:        1,
			LaneCapacity: 64,
			QoS0Overflow: OverflowBlock,
			QoS1Overflow: OverflowBlock,
		},
	})
	if err != nil {
		t.Fatalf("StartFlow: %v", err)
	}

	const n = 100
	for i := range n {
		f.Router().Submit(makeMsg("key", i), 1)
	}

	drainCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := mgr.StopFlow(drainCtx, "graceful"); err != nil {
		t.Fatalf("StopFlow: %v", err)
	}

	mu.Lock()
	got := len(order)
	mu.Unlock()
	if got != n {
		t.Fatalf("processed %d messages, want %d", got, n)
	}

	mu.Lock()
	for i, v := range order {
		if v != i {
			t.Fatalf("message %d: got seq %d", i, v)
		}
	}
	mu.Unlock()
}

// orderedCollectBlock records the sequence numbers of received messages in
// order, for verifying that drain completes all in-flight work.
type orderedCollectBlock struct {
	mu    *sync.Mutex
	order *[]int
}

func (o *orderedCollectBlock) Process(msg *Message) ([]*Message, error) {
	v, _ := msg.Data.(int64)
	o.mu.Lock()
	defer o.mu.Unlock()
	*o.order = append(*o.order, int(v))
	return nil, nil
}

func (o *orderedCollectBlock) Name() string { return "ordered-collect" }

// TestManager_ShutdownDrainsAll verifies that Shutdown drains every
// running flow.
func TestManager_ShutdownDrainsAll(t *testing.T) {
	mgr := NewManager()
	ctx := context.Background()

	sinks := make(map[string]*collectBlock)
	for _, id := range []string{"x", "y"} {
		sink := &collectBlock{}
		sinks[id] = sink
		_, err := mgr.StartFlow(ctx, Config{
			PipelineID: id,
			Blocks:     []Block{sink},
		})
		if err != nil {
			t.Fatalf("StartFlow %s: %v", id, err)
		}
		sink2 := sinks[id] // nolint:govet
		f := findFlow(mgr, id)
		for i := range 10 {
			f.Router().Submit(makeMsg(id, i), 1)
		}
		_ = sink2 // ensure it's referenced
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := mgr.Shutdown(shutdownCtx); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}

	for id, sink := range sinks {
		if got := sink.len(); got != 10 {
			t.Errorf("flow %s: delivered %d, want 10", id, got)
		}
	}
}

func findFlow(mgr *Manager, pipelineID string) *Flow {
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()
	return mgr.flows[pipelineID]
}
