# issue-26 — flow-router: Stream-based message routing through the block graph

<!-- trickle-allow: internal/flow/router.go internal/flow/router_test.go -->

## Context

When a Flow is running, messages emitted by source blocks must be routed through the pipeline graph to transform and sink blocks.

Upstream Astarte Flow guarantees **in-order delivery within a stream key** (messages sharing the same `Key` string) while allowing parallel processing across distinct stream keys. Astrate currently has no message router for flows.

## What to do

1. **`internal/flow/router.go`**:
   - Implement `Router` struct managing message dispatch across blocks:
     - `NewRouter(graph *PipelineGraph, options RouterOptions) *Router`
     - `Push(ctx context.Context, msg FlowMessage) error`: accepts a message from a source block and enqueues it for routing.
   - Stream lane assignment:
     - Map `msg.Key` to a worker lane using consistent hashing (`fnv` or `crc32` hash modulo number of lanes).
     - Each lane runs a worker goroutine reading from a bounded message queue, guaranteeing strict sequential execution per stream key.
   - Backpressure and Overflow handling:
     - `RouterOptions.QueueSize` (bounded channel capacity per lane, e.g. 1000 messages).
     - `OverflowPolicy` enum (`OverflowPolicyBlock`, `OverflowPolicyDrop`).
   - Per-block routing & error isolation:
     - Route output of block $A$ to target input of block $B$ per connection definition.
     - Recover/catch per-message panics or block errors without stopping the lane or crashing the entire flow. Increment error counter on the flow metrics context.

## Constraints

- Strictly preserve message order per stream key.
- Avoid global mutex locks during message push; use worker channel dispatch.
- If executed on Raspberry Pi (no race detector), log loud warning in report if concurrency primitives are added.

## Acceptance criteria

Run from repo root:

```sh
go build ./internal/flow/...
go test ./internal/flow/... -run TestFlowRouter
```

Tests to add:

- `internal/flow/router_test.go`:
  - `TestRouter_InOrderDeliveryPerKey`: pushes 100 messages with key `"device-A"` and sequence numbers $1 \dots 100$. Verifies sink receives all 100 messages strictly in order $1 \dots 100$.
  - `TestRouter_ParallelInterleavingAcrossKeys`: pushes interleaved messages for `"device-A"` and `"device-B"`. Verifies processing completes without blocking or deadlocks.
  - `TestRouter_OverflowPolicyDrop`: saturates a bounded lane under `OverflowPolicyDrop` and verifies queue overflow drops excess messages cleanly without blocking caller.
  - `TestRouter_ErrorIsolation`: injects a block that returns an error for specific message IDs; asserts other messages proceed normally and error counts increase.
