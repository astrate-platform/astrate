# issue-25 — flow-lifecycle: Flow start/stop lifecycle and runtime management

<!-- trickle-allow: internal/flow/manager.go internal/flow/manager_test.go internal/store/flows.go internal/store/flows_test.go migrations/000009_flows.up.sql migrations/000009_flows.down.sql -->

## Context

Flows in `astarte_flow` are running instances of pipelines — they allocate resources, consume CPU/memory, process message streams, and track operational state.

Astrate currently has no runtime manager for initializing, executing, monitoring, or gracefully stopping flows.

## What to do

1. **`migrations/000009_flows.up.sql` and `migrations/000009_flows.down.sql`**:
   - Create table `flows`:
     - `id` (UUID, PRIMARY KEY)
     - `realm_id` (TEXT, NOT NULL, REFERENCES realms(id) ON DELETE CASCADE)
     - `pipeline_id` (UUID, NOT NULL, REFERENCES pipelines(id) ON DELETE CASCADE)
     - `status` (TEXT, NOT NULL) -- 'creating', 'running', 'stopped', 'failed'
     - `config` (JSONB, NOT NULL DEFAULT '{}')
     - `error_message` (TEXT)
     - `created_at` (TIMESTAMPTZ, NOT NULL, DEFAULT clock_timestamp())
     - `updated_at` (TIMESTAMPTZ, NOT NULL, DEFAULT clock_timestamp())
   - Add indices on `(realm_id)` and `(pipeline_id)`.

2. **`internal/store/flows.go`**:
   - Store methods for flow persistence: `CreateFlow`, `GetFlow`, `UpdateFlowStatus`, `ListFlows`, `DeleteFlow`.

3. **`internal/flow/manager.go`**:
   - Implement `Manager` struct holding active running flows in memory:
     - `StartFlow(ctx context.Context, realmID string, pipelineID string, config map[string]any) (*Flow, error)`:
       - Loads pipeline definition from store.
       - Validates pipeline graph.
       - Spawns block execution contexts and channels.
       - Sets flow status to `FlowStatusRunning` in store & memory.
     - `StopFlow(ctx context.Context, flowID string) error`:
       - Signals cancellation context to running blocks.
       - Drains in-flight messages within timeout.
       - Releases block resources and goroutines.
       - Updates status to `FlowStatusStopped`.
     - `GetFlowStatus(ctx context.Context, flowID string) (*FlowStatusReport, error)`:
       - Returns current state, uptime, and message counters.

## Constraints

- **Goroutine Leak Safety**: Every spawned goroutine must monitor context cancellation (`ctx.Done()`) or channel closure.
- **Graceful Shutdown**: `StopFlow` must enforce a configurable drain timeout (default e.g. 5 seconds) after which in-flight routines are forcibly terminated.
- If executed on Raspberry Pi (no race detector), log loud warning in report if concurrency primitives are added.

## Acceptance criteria

Run from repo root:

```sh
go build ./internal/flow/... ./internal/store/...
go test ./internal/flow/... -run TestFlowManager
```

Tests to add:

- `internal/flow/manager_test.go`:
  - `TestFlowManager_StartStopLifecycle`: starts a mock flow, asserts status is `running`, calls `StopFlow`, asserts status transitions to `stopped` and all goroutines terminate cleanly.
  - `TestFlowManager_StopDrainsInFlight`: sends messages to a mock flow, invokes `StopFlow`, verifies buffered messages are processed before shutdown completes.
  - `TestFlowManager_FailedInitialization`: attempts to start a flow with an invalid pipeline ID, asserts status transitions to `failed` with recorded error.
