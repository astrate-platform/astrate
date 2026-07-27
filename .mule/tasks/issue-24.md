# issue-24 — flow-pipeline-store: Pipeline CRUD storage and graph validation

<!-- trickle-allow: migrations/000008_pipelines.up.sql migrations/000008_pipelines.down.sql internal/store/pipelines.go internal/store/pipelines_test.go internal/store/store.go -->

## Context

Pipelines in `astarte_flow` are stored as JSON documents and validated upon creation and update. Astrate currently has no database schema or store methods for pipelines.

This task adds PostgreSQL migration `000008_pipelines` and implements pipeline CRUD methods and graph validation (acyclicity, node existence, connection validity) in `internal/store/`.

## What to do

1. **`migrations/000008_pipelines.up.sql` and `migrations/000008_pipelines.down.sql`**:
   - Create table `pipelines`:
     - `id` (UUID, PRIMARY KEY)
     - `realm_id` (TEXT, NOT NULL, REFERENCES realms(id) ON DELETE CASCADE)
     - `name` (TEXT, NOT NULL)
     - `definition` (JSONB, NOT NULL)
     - `created_at` (TIMESTAMPTZ, NOT NULL, DEFAULT clock_timestamp())
     - `updated_at` (TIMESTAMPTZ, NOT NULL, DEFAULT clock_timestamp())
     - `UNIQUE (realm_id, name)`
   - Add index on `(realm_id)`.

2. **`internal/store/pipelines.go`**:
   - Define `PipelineRecord` struct representing database row.
   - Implement pipeline CRUD methods:
     - `CreatePipeline(ctx context.Context, realmID string, p *PipelineRecord) error`
     - `GetPipeline(ctx context.Context, realmID string, idOrName string) (*PipelineRecord, error)`
     - `UpdatePipeline(ctx context.Context, realmID string, p *PipelineRecord) error`
     - `DeletePipeline(ctx context.Context, realmID string, idOrName string) error`
     - `ListPipelines(ctx context.Context, realmID string) ([]*PipelineRecord, error)`
   - Implement `ValidatePipelineGraph(def *PipelineDefinition) error`:
     - Verify all block IDs are unique within the pipeline.
     - Verify all connections reference valid `FromBlock` and `ToBlock` IDs.
     - Check for cycles using DFS or Kahn's algorithm (return an explicit `ErrPipelineCyclic` error if a cycle is detected).

## Constraints

- Store queries must match existing store error wrapping and SQL patterns (`fmt.Errorf("...: %w", err)`).
- Follow existing migration file naming conventions (`000008_pipelines.up.sql` / `000008_pipelines.down.sql`).
- Graph validation must be runnable in memory without requiring a database connection.

## Acceptance criteria

Run from repo root:

```sh
go build ./internal/store/...
go test ./internal/store/... -run TestPipelineStore
```

Tests to add:

- `internal/store/pipelines_test.go`:
  - `TestValidatePipelineGraph_Acyclic`: asserts a valid DAG passes graph validation.
  - `TestValidatePipelineGraph_Cyclic`: asserts a pipeline containing a loop fails validation with `ErrPipelineCyclic`.
  - `TestValidatePipelineGraph_InvalidBlockRef`: asserts a connection referencing a missing block ID fails validation.
  - `TestPipelineStore_CRUD` (integration test with test store): verifies Create, Get, Update, List, and Delete operations on PostgreSQL.
