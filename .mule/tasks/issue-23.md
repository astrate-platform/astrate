# issue-23 — flow-models: FlowMessage wire format and Block/Pipeline/Flow data models

<!-- trickle-allow: internal/flow/message.go internal/flow/block.go internal/flow/pipeline.go internal/flow/flow.go internal/flow/message_test.go internal/flow/pipeline_test.go internal/flow/flow_test.go -->

## Context

Astarte Flow's core abstractions — the `FlowMessage` wire format (`astarte_flow/message/v0.1`), `Block` interface (source/sink/transform), `Pipeline` graph, and `Flow` runtime instance — have no equivalent in Astrate today.

Astrate needs a clean `internal/flow` package establishing these fundamental data models and interface definitions in idiomatic Go before pipeline storage, routing, and blocks can be built.

## What to do

1. **`internal/flow/message.go`**:
   - Define `FlowMessage` struct matching `astarte_flow/message/v0.1` schema:
     - `ID` (string, UUID or unique string identifier)
     - `Key` (string, routing key for stream partition/ordering)
     - `Type` (string or `MessageType` enum, e.g. `datastream`, `properties`, `control`)
     - `Subtype` (string)
     - `Data` (any / `interface{}` representing payload value)
     - `Metadata` (`map[string]string` for arbitrary key-value headers)
     - `TimestampUs` (int64, timestamp in microseconds)
   - Implement JSON custom marshaling/unmarshaling ensuring compliance with upstream wire format.

2. **`internal/flow/block.go`**:
   - Define core `Block` interface:
     ```go
     type Block interface {
         ID() string
         Type() string
         Init(ctx context.Context, config map[string]any) error
         Start(ctx context.Context) error
         Stop(ctx context.Context) error
     }
     ```
   - Define `SourceBlock`, `TransformBlock`, `SinkBlock` specialized channels or signatures for message processing.

3. **`internal/flow/pipeline.go`**:
   - Define `Pipeline` struct:
     - `ID` (string)
     - `RealmID` (string)
     - `Name` (string)
     - `Blocks` ([]BlockConfig)
     - `Connections` ([]ConnectionConfig)
   - Define `ConnectionConfig`: `FromBlock`, `FromOutput`, `ToBlock`, `ToInput`.
   - Add JSON tags for pipeline graph representation.

4. **`internal/flow/flow.go`**:
   - Define `FlowStatus` enum (`FlowStatusCreating`, `FlowStatusRunning`, `FlowStatusStopped`, `FlowStatusFailed`).
   - Define `Flow` metadata struct representing a running pipeline instance.

## Constraints

- Restate capabilities in idiomatic Go. Do not attempt to replicate Elixir GenStage / OTP supervision architecture.
- Keep dependencies minimal. Use stdlib types or standard repo utilities.
- Do not touch existing `internal/engine` or `internal/store` files in this task (storage & ingestion wiring are covered in #24 and #27).

## Acceptance criteria

Run from repo root:

```sh
go build ./internal/flow/...
go test ./internal/flow/... -run TestFlowModels
go test ./internal/flow/...
```

Tests to add:

- `internal/flow/message_test.go`: `TestFlowMessage_JSONRoundTrip` asserting exact match with upstream `v0.1` JSON format (including `timestamp_us`, `key`, `data`, `metadata`).
- `internal/flow/pipeline_test.go`: `TestPipeline_Serialization` verifying pipeline definition JSON encoding and decoding.
- `internal/flow/flow_test.go`: `TestFlowStatus_String` verifying status enum string representations and validation.
