# issue-26 — flow-router: Stream-based message routing through the block graph

<!-- trickle-allow: internal/flow/router.go internal/flow/router_test.go -->

## Context

When a Flow is running, messages must be routed through the pipeline graph with strict in-order delivery per stream key (`Key`), while allowing concurrent non-blocking processing across different keys.

This specification provides the exact `Router` struct, worker lane hashing, overflow policies, and a **loud-failing test suite**.

## Exact Code Specifications

### 1. `internal/flow/router.go`

```go
package flow

import (
	"context"
	"errors"
	"hash/fnv"
	"sync"
	"sync/atomic"
)

var ErrQueueFull = errors.New("lane queue full")

type OverflowPolicy string

const (
	OverflowPolicyBlock OverflowPolicy = "block"
	OverflowPolicyDrop  OverflowPolicy = "drop"
)

type RouterOptions struct {
	NumLanes       int
	QueueSize      int
	OverflowPolicy OverflowPolicy
}

type lane struct {
	ch chan FlowMessage
}

type Router struct {
	opts         RouterOptions
	lanes        []lane
	sinkFunc     func(ctx context.Context, msg FlowMessage) error
	errorCount   uint64
	successCount uint64
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewRouter(opts RouterOptions, sinkFunc func(ctx context.Context, msg FlowMessage) error) *Router {
	if opts.NumLanes <= 0 {
		opts.NumLanes = 4
	}
	if opts.QueueSize <= 0 {
		opts.QueueSize = 100
	}
	if opts.OverflowPolicy == "" {
		opts.OverflowPolicy = OverflowPolicyBlock
	}

	ctx, cancel := context.WithCancel(context.Background())
	r := &Router{
		opts:     opts,
		lanes:    make([]lane, opts.NumLanes),
		sinkFunc: sinkFunc,
		ctx:      ctx,
		cancel:   cancel,
	}

	for i := 0; i < opts.NumLanes; i++ {
		r.lanes[i] = lane{ch: make(chan FlowMessage, opts.QueueSize)}
		r.wg.Add(1)
		go r.worker(i, r.lanes[i].ch)
	}

	return r
}

func (r *Router) getLaneIndex(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) % r.opts.NumLanes
}

func (r *Router) Push(ctx context.Context, msg FlowMessage) error {
	idx := r.getLaneIndex(msg.Key)
	l := r.lanes[idx]

	if r.opts.OverflowPolicy == OverflowPolicyDrop {
		select {
		case l.ch <- msg:
			return nil
		default:
			atomic.AddUint64(&r.errorCount, 1)
			return ErrQueueFull
		}
	} else {
		select {
		case l.ch <- msg:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *Router) worker(laneID int, ch <-chan FlowMessage) {
	defer r.wg.Done()
	for {
		select {
		case <-r.ctx.Done():
			for msg := range ch {
				r.processMessage(msg)
			}
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			r.processMessage(msg)
		}
	}
}

func (r *Router) processMessage(msg FlowMessage) {
	if r.sinkFunc != nil {
		if err := r.sinkFunc(r.ctx, msg); err != nil {
			atomic.AddUint64(&r.errorCount, 1)
			return
		}
	}
	atomic.AddUint64(&r.successCount, 1)
}

func (r *Router) Stop() {
	r.cancel()
	for _, l := range r.lanes {
		close(l.ch)
	}
	r.wg.Wait()
}

func (r *Router) Stats() (success uint64, errors uint64) {
	return atomic.LoadUint64(&r.successCount), atomic.LoadUint64(&r.errorCount)
}
```

---

## Test Suite Specifications (Loud-Failing Assertions)

Create `internal/flow/router_test.go`:

```go
package flow_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
)

func TestRouter_StrictInOrderDeliveryPerKey(t *testing.T) {
	var mu sync.Mutex
	received := make(map[string][]int64)

	sink := func(ctx context.Context, msg flow.FlowMessage) error {
		mu.Lock()
		defer mu.Unlock()
		seq := msg.Data.(int64)
		received[msg.Key] = append(received[msg.Key], seq)
		return nil
	}

	opts := flow.RouterOptions{NumLanes: 4, QueueSize: 1000, OverflowPolicy: flow.OverflowPolicyBlock}
	r := flow.NewRouter(opts, sink)
	ctx := context.Background()

	// Push 100 messages for device-A and 100 for device-B
	const totalMsgs = 100
	for i := int64(1); i <= totalMsgs; i++ {
		_ = r.Push(ctx, flow.FlowMessage{ID: "m", Key: "realm/device-A", Type: "data", Data: i})
		_ = r.Push(ctx, flow.FlowMessage{ID: "m", Key: "realm/device-B", Type: "data", Data: i})
	}

	r.Stop()

	// Loud assertions on in-order delivery
	mu.Lock()
	defer mu.Unlock()

	for _, key := range []string{"realm/device-A", "realm/device-B"} {
		seqs := received[key]
		if len(seqs) != totalMsgs {
			t.Fatalf("FAIL: expected %d messages for key %s, got %d", totalMsgs, key, len(seqs))
		}
		for i := 0; i < totalMsgs; i++ {
			expectedSeq := int64(i + 1)
			if seqs[i] != expectedSeq {
				t.Fatalf("FAIL ORDER VIOLATION for key %s: position %d expected seq %d, got %d", key, i, expectedSeq, seqs[i])
			}
		}
	}
}

func TestRouter_OverflowPolicyDrop(t *testing.T) {
	blockSink := func(ctx context.Context, msg flow.FlowMessage) error {
		time.Sleep(500 * time.Millisecond) // artificially delay processing
		return nil
	}

	opts := flow.RouterOptions{NumLanes: 1, QueueSize: 2, OverflowPolicy: flow.OverflowPolicyDrop}
	r := flow.NewRouter(opts, blockSink)
	ctx := context.Background()

	// Push 10 messages into lane with queue size 2
	var droppedCount int
	for i := 0; i < 10; i++ {
		err := r.Push(ctx, flow.FlowMessage{ID: "m", Key: "same-key", Type: "t"})
		if errors.Is(err, flow.ErrQueueFull) {
			droppedCount++
		}
	}

	r.Stop()

	if droppedCount == 0 {
		t.Fatalf("FAIL: expected messages to be dropped with ErrQueueFull under OverflowPolicyDrop")
	}
}
```

---

## Negative Constraints & TDD Workflow

1. ❌ **Do NOT use global mutexes during `Push`**. Lane assignment must use lock-free `fnv.New32a()` hashing.
2. ❌ **Do NOT crash workers on error**.
3. **Execution order**:
   - Write `internal/flow/router_test.go`.
   - Run `go test ./internal/flow/...` (must fail).
   - Write `internal/flow/router.go`.
   - Run `go test -v ./internal/flow/...` (must pass clean).

## Acceptance criteria

Run from repo root:

```sh
go build ./internal/flow/...
go test -v ./internal/flow/... -run TestRouter
```
