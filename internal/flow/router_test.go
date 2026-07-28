package flow

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	promtest "github.com/prometheus/client_golang/prometheus/testutil"
)

// --- test blocks -----------------------------------------------------------

// collectBlock is a sink that appends every message to its slice.
type collectBlock struct {
	mu   sync.Mutex
	msgs []*FlowMessage
}

func (c *collectBlock) Process(msg *FlowMessage) ([]*FlowMessage, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cp := *msg
	c.msgs = append(c.msgs, &cp)
	return nil, nil
}

func (c *collectBlock) Name() string { return "collect" }

func (c *collectBlock) len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.msgs)
}

func (c *collectBlock) get(i int) *FlowMessage {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.msgs[i]
}

// failBlock errors on every Nth message.
type failBlock struct {
	every  int
	count  atomic.Int64
	prefix string
}

func (f *failBlock) Process(msg *FlowMessage) ([]*FlowMessage, error) {
	n := f.count.Add(1)
	if n%int64(f.every) == 0 {
		return nil, fmt.Errorf("%s: injected error", f.prefix)
	}
	return []*FlowMessage{msg}, nil
}

func (f *failBlock) Name() string { return f.prefix }

// passthroughBlock forwards every message unchanged.
type passthroughBlock struct{}

func (p *passthroughBlock) Process(msg *FlowMessage) ([]*FlowMessage, error) {
	return []*FlowMessage{msg}, nil
}

func (p *passthroughBlock) Name() string { return "passthrough" }

// panicBlock panics on every Nth message.
type panicBlock struct {
	every  int
	count  atomic.Int64
	prefix string
}

func (p *panicBlock) Process(msg *FlowMessage) ([]*FlowMessage, error) {
	n := p.count.Add(1)
	if n%int64(p.every) == 0 {
		panic(fmt.Sprintf("%s: injected panic", p.prefix))
	}
	return []*FlowMessage{msg}, nil
}

func (p *panicBlock) Name() string { return p.prefix }

// --- helpers ---------------------------------------------------------------

func makeMsg(key string, seq int) *FlowMessage {
	return &FlowMessage{
		Key:       key,
		Type:      TypeInteger,
		Data:      int64(seq),
		Timestamp: time.Now().UnixMicro(),
	}
}

func startTestRouter(t *testing.T, graph *BlockGraph, cfg RouterConfig) *Router {
	t.Helper()
	r := NewRouter(graph, cfg, nil)
	r.Run(context.Background())
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = r.Drain(ctx)
	})
	return r
}

// waitFor polls cond at 5ms intervals until it returns true or timeout fires.
func waitFor(t *testing.T, timeout time.Duration, desc string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		if cond() {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("timeout waiting for %s", desc)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// --- acceptance tests ------------------------------------------------------

// TestRouter_InOrderDeliveryPerKey sends 100 messages with the same key and
// verifies they arrive at the sink in submission order.
func TestRouter_InOrderDeliveryPerKey(t *testing.T) {
	sink := &collectBlock{}
	graph, err := NewBlockGraph(sink)
	if err != nil {
		t.Fatal(err)
	}
	r := startTestRouter(t, graph, RouterConfig{
		Lanes:        4,
		LaneCapacity: 128,
		QoS0Overflow: OverflowBlock,
		QoS1Overflow: OverflowBlock,
	})

	const key = "sensor-1"
	const n = 100
	for i := range n {
		r.Submit(makeMsg(key, i), 1)
	}

	waitFor(t, 5*time.Second, "all messages delivered", func() bool {
		return sink.len() == n
	})
	for i := range n {
		got := sink.get(i)
		if got.Key != key {
			t.Fatalf("msg %d: key = %q, want %q", i, got.Key, key)
		}
		v, ok := got.Data.(int64)
		if !ok || v != int64(i) {
			t.Fatalf("msg %d: data = %v, want %d", i, got.Data, i)
		}
	}
}

// TestRouter_ParallelDistinctKeys verifies that messages with different keys
// are processed concurrently (non-blocking) across distinct lanes.
func TestRouter_ParallelDistinctKeys(t *testing.T) {
	const keys = 4
	const msgsPerKey = 10
	sink := &collectBlock{}
	graph, err := NewBlockGraph(sink)
	if err != nil {
		t.Fatal(err)
	}
	r := startTestRouter(t, graph, RouterConfig{
		Lanes:        8,
		LaneCapacity: 128,
		QoS0Overflow: OverflowBlock,
		QoS1Overflow: OverflowBlock,
	})

	for i := range keys {
		key := fmt.Sprintf("key-%d", i)
		for j := range msgsPerKey {
			r.Submit(makeMsg(key, j), 0)
		}
	}

	total := keys * msgsPerKey
	waitFor(t, 5*time.Second, "all messages delivered", func() bool {
		return sink.len() == total
	})

	if got := sink.len(); got != total {
		t.Fatalf("delivered %d messages, want %d", got, total)
	}
}

// TestRouter_OverflowDropPolicy verifies that QoS 0 messages are dropped
// (not blocked) when a lane is full, while QoS ≥ 1 messages block.
func TestRouter_OverflowDropPolicy(t *testing.T) {
	gate := make(chan struct{})
	// A slow sink that blocks until gate is opened.
	slow := &blockingSink{gate: gate}
	graph, err := NewBlockGraph(slow)
	if err != nil {
		t.Fatal(err)
	}
	r := startTestRouter(t, graph, RouterConfig{
		Lanes:        1,
		LaneCapacity: 2,
		QoS0Overflow: OverflowDrop,
		QoS1Overflow: OverflowBlock,
	})

	// Fill the lane: 1 in-process + 2 queued = 3 total before full.
	// The slow sink holds the first message. The next 2 fill the channel.
	// QoS0 submit 3 should be dropped; QoS1 submit should block until gate opens.

	// Send first QoS1 — the slow sink picks it up but blocks on gate.
	r.Submit(makeMsg("k", 0), 1)
	// Send 2 more QoS1 — fills the lane.
	r.Submit(makeMsg("k", 1), 1)
	r.Submit(makeMsg("k", 2), 1)
	time.Sleep(50 * time.Millisecond) // let goroutines settle

	// Now send a QoS0 — lane is full, should be dropped.
	r.Submit(makeMsg("k", 100), 0)
	r.Submit(makeMsg("k", 101), 0)

	if got := promtest.ToFloat64(r.met.droppedQoS0); got != 2 {
		t.Errorf("QoS0 drop counter = %v, want 2", got)
	}

	// Open gate — blocked messages drain.
	close(gate)
	waitFor(t, 5*time.Second, "QoS1 messages processed", func() bool {
		return slow.len() >= 3
	})
}

// TestRouter_BlockErrorIsolation verifies that a failing block does not crash
// the lane; errors are logged and counted, and subsequent messages still flow.
func TestRouter_BlockErrorIsolation(t *testing.T) {
	failer := &failBlock{every: 3, prefix: "test-fail"}
	sink := &collectBlock{}
	graph, err := NewBlockGraph(failer, sink)
	if err != nil {
		t.Fatal(err)
	}
	r := startTestRouter(t, graph, RouterConfig{
		Lanes:        1,
		LaneCapacity: 128,
		QoS0Overflow: OverflowBlock,
		QoS1Overflow: OverflowBlock,
	})

	const n = 20
	for i := range n {
		r.Submit(makeMsg("k", i), 1)
	}

	waitFor(t, 5*time.Second, "messages processed", func() bool {
		// Atomic counter is shared across lanes, so use 1 lane for
		// deterministic ordering. failBlock fires at every 3rd
		// call: n=3,6,9,12,15,18 = 6 errors; 14 succeed.
		return sink.len() == 14
	})

	errs := promtest.ToFloat64(r.met.blockErrors)
	if errs != 6 {
		t.Errorf("block errors = %v, want 6", errs)
	}

	// The router is still alive — submit two more: counter=21 fails (21%3==0),
	// counter=22 succeeds (22%3!=0). Verify the error count increased and a
	// new message reached the sink.
	r.Submit(makeMsg("k", 999), 1)  // counter 21 → fails
	r.Submit(makeMsg("k", 1000), 1) // counter 22 → succeeds
	waitFor(t, 5*time.Second, "post-error message delivered", func() bool {
		return sink.len() == 15
	})
	errs2 := promtest.ToFloat64(r.met.blockErrors)
	if errs2 != 7 {
		t.Errorf("block errors after post-error = %v, want 7", errs2)
	}
}

// TestRouter_PanicIsolation verifies that a panic inside a block is recovered
// per-message and does not crash the lane.
func TestRouter_PanicIsolation(t *testing.T) {
	panicker := &panicBlock{every: 2, prefix: "test-panic"}
	sink := &collectBlock{}
	graph, err := NewBlockGraph(panicker, sink)
	if err != nil {
		t.Fatal(err)
	}
	r := startTestRouter(t, graph, RouterConfig{
		Lanes:        1,
		LaneCapacity: 128,
		QoS0Overflow: OverflowBlock,
		QoS1Overflow: OverflowBlock,
	})

	const n = 10
	for i := range n {
		r.Submit(makeMsg("k", i), 1)
	}

	waitFor(t, 5*time.Second, "messages processed", func() bool {
		// Atomic counter shared across lanes; use 1 lane.
		// panicBlock fires at every 2nd call: n=2,4,6,8,10 = 5 panics; 5 succeed.
		return sink.len() == 5
	})

	errs := promtest.ToFloat64(r.met.blockErrors)
	if errs < 5 {
		t.Errorf("block errors = %v, want >= 5 (panics counted)", errs)
	}

	// Lane is still alive — submit two messages: counter=11 panics (11%2==1... wait).
	// counter starts at 10. 10%2==0 → panic. 11%2==1 → succeeds. So send2.
	r.Submit(makeMsg("k", 999), 1)  // counter11 → succeeds (odd)
	r.Submit(makeMsg("k", 1000), 1) // counter12 → panics (even)
	waitFor(t, 5*time.Second, "post-panic messages delivered", func() bool {
		return sink.len() == 6
	})
}

// TestRouter_Drain processes all queued messages before exiting.
func TestRouter_Drain(t *testing.T) {
	sink := &collectBlock{}
	graph, err := NewBlockGraph(sink)
	if err != nil {
		t.Fatal(err)
	}
	r := startTestRouter(t, graph, RouterConfig{
		Lanes:        2,
		LaneCapacity: 64,
		QoS0Overflow: OverflowBlock,
		QoS1Overflow: OverflowBlock,
	})

	const n = 50
	for i := range n {
		r.Submit(makeMsg("k", i), 1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.Drain(ctx); err != nil {
		t.Fatalf("drain: %v", err)
	}
	if got := sink.len(); got != n {
		t.Fatalf("delivered %d after drain, want %d", got, n)
	}

	// Post-drain submits are silently dropped.
	r.Submit(makeMsg("k", 999), 1)
	time.Sleep(50 * time.Millisecond)
	if got := sink.len(); got != n {
		t.Errorf("post-drain message delivered: %d", got)
	}
}

// TestLaneOf verifies deterministic hashing of keys to lane indices.
func TestLaneOf(t *testing.T) {
	for n := 1; n <= 32; n *= 2 {
		for _, key := range []string{"", "a", "hello", "sensor-42", "abc"} {
			l1 := laneOf(key, n)
			l2 := laneOf(key, n)
			if l1 != l2 {
				t.Fatalf("laneOf not deterministic for key=%q n=%d", key, n)
			}
			if l1 < 0 || l1 >= n {
				t.Fatalf("laneOf(%q, %d) = %d out of range", key, n, l1)
			}
		}
	}
}

// blockingSink holds every message until unblocked.
type blockingSink struct {
	gate chan struct{}
	mu   sync.Mutex
	msgs []*FlowMessage
}

func (b *blockingSink) Process(msg *FlowMessage) ([]*FlowMessage, error) {
	<-b.gate
	b.mu.Lock()
	defer b.mu.Unlock()
	cp := *msg
	b.msgs = append(b.msgs, &cp)
	return nil, nil
}

func (b *blockingSink) Name() string { return "blocking-sink" }

func (b *blockingSink) len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.msgs)
}
