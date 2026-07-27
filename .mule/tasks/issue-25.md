# issue-25 — flow-lifecycle: Flow start/stop lifecycle and runtime management

<!-- trickle-allow: internal/flow/manager.go internal/flow/manager_test.go internal/store/flows.go internal/store/flows_test.go migrations/000009_flows.up.sql migrations/000009_flows.down.sql -->

## Context

Flows in `astarte_flow` are running instances of pipelines. This task provides the `flows` DB migration, store procedures, and the runtime `flow.Manager` responsible for launching, monitoring, and gracefully stopping flows.

## Exact Code & Schema Specifications

### 1. `migrations/000009_flows.up.sql`

```sql
CREATE TABLE IF NOT EXISTS flows (
    id UUID PRIMARY KEY,
    realm_id SMALLINT NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    pipeline_id UUID NOT NULL REFERENCES pipelines(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp()
);

CREATE INDEX IF NOT EXISTS idx_flows_realm_id ON flows(realm_id);
CREATE INDEX IF NOT EXISTS idx_flows_pipeline_id ON flows(pipeline_id);
```

### 2. `migrations/000009_flows.down.sql`

```sql
DROP TABLE IF EXISTS flows;
```

### 3. `internal/flow/manager.go`

```go
package flow

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrFlowNotFound       = errors.New("flow not found")
	ErrFlowAlreadyRunning = errors.New("flow is already running")
	ErrFlowNotRunning     = errors.New("flow is not running")
)

type runningFlow struct {
	flow       *Flow
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
	doneChan   chan struct{}
}

type Manager struct {
	mu           sync.RWMutex
	running      map[string]*runningFlow
	drainTimeout time.Duration
}

func NewManager(drainTimeout time.Duration) *Manager {
	if drainTimeout <= 0 {
		drainTimeout = 5 * time.Second
	}
	return &Manager{
		running:      make(map[string]*runningFlow),
		drainTimeout: drainTimeout,
	}
}

func (m *Manager) StartFlow(ctx context.Context, f *Flow) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.running[f.ID]; exists {
		return ErrFlowAlreadyRunning
	}

	_, cancel := context.WithCancel(context.Background())
	rf := &runningFlow{
		flow:       f,
		cancelFunc: cancel,
		doneChan:   make(chan struct{}),
	}

	f.Status = FlowStatusRunning
	f.UpdatedAt = time.Now()
	m.running[f.ID] = rf
	return nil
}

func (m *Manager) StopFlow(ctx context.Context, flowID string) error {
	m.mu.Lock()
	rf, exists := m.running[flowID]
	if !exists {
		m.mu.Unlock()
		return ErrFlowNotFound
	}
	delete(m.running, flowID)
	m.mu.Unlock()

	// Signal cancellation to flow goroutines
	rf.cancelFunc()

	// Wait for graceful shutdown with timeout
	done := make(chan struct{})
	go func() {
		rf.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		rf.flow.Status = FlowStatusStopped
	case <-time.After(m.drainTimeout):
		rf.flow.Status = FlowStatusStopped
		rf.flow.ErrorMessage = "stopped forcibly after drain timeout"
	}

	rf.flow.UpdatedAt = time.Now()
	return nil
}

func (m *Manager) GetFlowStatus(flowID string) (*Flow, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rf, exists := m.running[flowID]
	if !exists {
		return nil, ErrFlowNotFound
	}
	return rf.flow, nil
}
```

---

## Test Suite Specifications (Loud-Failing Assertions)

Create `internal/flow/manager_test.go` with exact loud-failing test cases:

```go
package flow_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
)

func TestFlowManager_StartAndStop(t *testing.T) {
	mgr := flow.NewManager(1 * time.Second)
	ctx := context.Background()

	f := &flow.Flow{
		ID:         "flow-uuid-1",
		Realm:      "testrealm",
		PipelineID: "pipe-uuid-1",
	}

	// 1. Start Flow
	if err := mgr.StartFlow(ctx, f); err != nil {
		t.Fatalf("FAIL: StartFlow failed: %v", err)
	}

	st, err := mgr.GetFlowStatus("flow-uuid-1")
	if err != nil {
		t.Fatalf("FAIL: GetFlowStatus failed: %v", err)
	}
	if st.Status != flow.FlowStatusRunning {
		t.Fatalf("FAIL: expected status FlowStatusRunning, got %q", st.Status)
	}

	// 2. Duplicate Start must fail
	if err := mgr.StartFlow(ctx, f); !errors.Is(err, flow.ErrFlowAlreadyRunning) {
		t.Fatalf("FAIL: expected ErrFlowAlreadyRunning on duplicate start, got %v", err)
	}

	// 3. Stop Flow
	if err := mgr.StopFlow(ctx, "flow-uuid-1"); err != nil {
		t.Fatalf("FAIL: StopFlow failed: %v", err)
	}

	// Status after stop must be not found in active running map
	if _, err := mgr.GetFlowStatus("flow-uuid-1"); !errors.Is(err, flow.ErrFlowNotFound) {
		t.Fatalf("FAIL: expected ErrFlowNotFound after StopFlow, got %v", err)
	}
}

func TestFlowManager_StopNonExistent(t *testing.T) {
	mgr := flow.NewManager(1 * time.Second)
	if err := mgr.StopFlow(context.Background(), "unknown"); !errors.Is(err, flow.ErrFlowNotFound) {
		t.Fatalf("FAIL: expected ErrFlowNotFound stopping non-existent flow, got %v", err)
	}
}
```

---

## Negative Constraints & TDD Workflow

1. ❌ **Do NOT leak goroutines**. Always pair `WithCancel` with a deferred `cancel()` or explicit call during `StopFlow`.
2. ❌ **Report unverified races loud and clear in summary if running on Raspberry Pi** (`go test ./...` without `-race`).
3. **Execution order**:
   - Write `internal/flow/manager_test.go`.
   - Run `go test ./internal/flow/...` (must fail).
   - Write `internal/flow/manager.go` and migration files.
   - Run `go test -v ./internal/flow/...` (must pass clean).

## Acceptance criteria

Run from repo root:

```sh
go build ./internal/flow/... ./internal/store/...
go test -v ./internal/flow/... -run TestFlowManager
```
