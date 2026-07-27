# issue-24 — flow-pipeline-store: Pipeline CRUD storage and graph validation

<!-- trickle-allow: migrations/000008_pipelines.up.sql migrations/000008_pipelines.down.sql internal/store/pipelines.go internal/store/pipelines_test.go internal/store/store.go -->

## Context

Pipelines in `astarte_flow` are stored as JSON documents and validated for acyclicity and block connections upon creation/update.

This specification provides the exact SQL migration, Store methods wrapping `store.ErrNotFound` / `store.ErrAlreadyExists`, and a **loud-failing test suite** for graph validation and CRUD.

## Exact Code & Schema Specifications

### 1. `migrations/000008_pipelines.up.sql`

```sql
CREATE TABLE IF NOT EXISTS pipelines (
    id UUID PRIMARY KEY,
    realm_id SMALLINT NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    definition JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT clock_timestamp(),
    UNIQUE (realm_id, name)
);

CREATE INDEX IF NOT EXISTS idx_pipelines_realm_id ON pipelines(realm_id);
```

### 2. `migrations/000008_pipelines.down.sql`

```sql
DROP TABLE IF EXISTS pipelines;
```

### 3. `internal/store/pipelines.go`

```go
package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/astrate-platform/astrate/internal/flow"
)

var (
	ErrPipelineCyclic   = errors.New("pipeline graph contains cycles")
	ErrInvalidBlockRef  = errors.New("connection references nonexistent block")
	ErrDuplicateBlockID = errors.New("duplicate block id in pipeline definition")
)

type PipelineRecord struct {
	ID         string        `db:"id"`
	RealmID    int16         `db:"realm_id"`
	Name       string        `db:"name"`
	Definition flow.Pipeline `db:"definition"`
	CreatedAt  time.Time     `db:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at"`
}

// ValidatePipelineGraph performs in-memory graph validation without DB access.
func ValidatePipelineGraph(p *flow.Pipeline) error {
	blockMap := make(map[string]bool)
	for _, b := range p.Blocks {
		if b.ID == "" {
			return fmt.Errorf("%w: block missing id", ErrInvalidBlockRef)
		}
		if blockMap[b.ID] {
			return fmt.Errorf("%w: %s", ErrDuplicateBlockID, b.ID)
		}
		blockMap[b.ID] = true
	}

	adj := make(map[string][]string)
	inDegree := make(map[string]int)
	for bID := range blockMap {
		inDegree[bID] = 0
	}

	for _, conn := range p.Connections {
		if !blockMap[conn.FromBlock] {
			return fmt.Errorf("%w: from_block %s", ErrInvalidBlockRef, conn.FromBlock)
		}
		if !blockMap[conn.ToBlock] {
			return fmt.Errorf("%w: to_block %s", ErrInvalidBlockRef, conn.ToBlock)
		}
		adj[conn.FromBlock] = append(adj[conn.FromBlock], conn.ToBlock)
		inDegree[conn.ToBlock]++
	}

	// Kahn's algorithm for cycle detection
	queue := make([]string, 0)
	for bID, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, bID)
		}
	}

	visitedCount := 0
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		visitedCount++

		for _, neighbor := range adj[curr] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if visitedCount != len(blockMap) {
		return ErrPipelineCyclic
	}

	return nil
}
```

---

## Test Suite Specifications (Loud-Failing Assertions)

Create `internal/store/pipelines_test.go` with exact loud-failing test cases:

```go
package store_test

import (
	"errors"
	"testing"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/store"
)

func TestValidatePipelineGraph_Acyclic(t *testing.T) {
	// A -> B -> C (valid DAG)
	p := &flow.Pipeline{
		Blocks: []flow.BlockConfig{
			{ID: "A"}, {ID: "B"}, {ID: "C"},
		},
		Connections: []flow.ConnectionConfig{
			{FromBlock: "A", ToBlock: "B"},
			{FromBlock: "B", ToBlock: "C"},
		},
	}

	if err := store.ValidatePipelineGraph(p); err != nil {
		t.Fatalf("FAIL: expected valid DAG to pass graph validation, got error: %v", err)
	}
}

func TestValidatePipelineGraph_DirectCycle(t *testing.T) {
	// A -> B -> A (cycle)
	p := &flow.Pipeline{
		Blocks: []flow.BlockConfig{
			{ID: "A"}, {ID: "B"},
		},
		Connections: []flow.ConnectionConfig{
			{FromBlock: "A", ToBlock: "B"},
			{FromBlock: "B", ToBlock: "A"},
		},
	}

	err := store.ValidatePipelineGraph(p)
	if !errors.Is(err, store.ErrPipelineCyclic) {
		t.Fatalf("FAIL: expected ErrPipelineCyclic for A->B->A cycle, got: %v", err)
	}
}

func TestValidatePipelineGraph_IndirectCycle(t *testing.T) {
	// A -> B -> C -> A (cycle)
	p := &flow.Pipeline{
		Blocks: []flow.BlockConfig{
			{ID: "A"}, {ID: "B"}, {ID: "C"},
		},
		Connections: []flow.ConnectionConfig{
			{FromBlock: "A", ToBlock: "B"},
			{FromBlock: "B", ToBlock: "C"},
			{FromBlock: "C", ToBlock: "A"},
		},
	}

	err := store.ValidatePipelineGraph(p)
	if !errors.Is(err, store.ErrPipelineCyclic) {
		t.Fatalf("FAIL: expected ErrPipelineCyclic for A->B->C->A cycle, got: %v", err)
	}
}

func TestValidatePipelineGraph_MissingBlockReference(t *testing.T) {
	// A -> NONEXISTENT
	p := &flow.Pipeline{
		Blocks: []flow.BlockConfig{{ID: "A"}},
		Connections: []flow.ConnectionConfig{
			{FromBlock: "A", ToBlock: "NONEXISTENT"},
		},
	}

	err := store.ValidatePipelineGraph(p)
	if !errors.Is(err, store.ErrInvalidBlockRef) {
		t.Fatalf("FAIL: expected ErrInvalidBlockRef for missing block target, got: %v", err)
	}
}

func TestValidatePipelineGraph_DuplicateBlockID(t *testing.T) {
	p := &flow.Pipeline{
		Blocks: []flow.BlockConfig{
			{ID: "dup"}, {ID: "dup"},
		},
	}

	err := store.ValidatePipelineGraph(p)
	if !errors.Is(err, store.ErrDuplicateBlockID) {
		t.Fatalf("FAIL: expected ErrDuplicateBlockID for duplicate block ID, got: %v", err)
	}
}
```

---

## Negative Constraints & TDD Workflow

1. ❌ **Do NOT use external graph libraries**. Use Kahn's algorithm provided above.
2. ❌ **Do NOT use string for `realm_id` in SQL**. Astrate realms use `smallint` (`int16`).
3. **Execution order**:
   - Write `internal/store/pipelines_test.go`.
   - Run `go test ./internal/store/...` (must fail until `ValidatePipelineGraph` is implemented).
   - Write `migrations/000008_pipelines.up.sql` and `internal/store/pipelines.go`.
   - Run `go test -v ./internal/store/...` (must pass clean).

## Acceptance criteria

Run from repo root:

```sh
go build ./internal/store/...
go test -v ./internal/store/... -run TestValidatePipelineGraph
```
