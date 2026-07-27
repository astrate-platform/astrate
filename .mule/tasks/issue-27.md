# issue-27 — flow-astarte-source: Astarte device events source block

<!-- trickle-allow: internal/flow/blocks/astartesource/source.go internal/flow/blocks/astartesource/source_test.go internal/engine/stream/bus.go -->

## Context

The most critical built-in block for Astrate's flow feature parity: a Source block that ingests Astarte device events from Astrate's engine stream bus (`internal/engine/stream`) and converts them into `FlowMessage`s for downstream flow processing.

Astrate currently has no flow source blocks connecting the ingestion engine to the flow framework.

## What to do

1. **`internal/flow/blocks/astartesource/source.go`**:
   - Implement `AstarteSourceBlock` implementing `internal/flow.Block`:
     - `Init(ctx context.Context, config map[string]any) error`:
       - Parse configuration parameters:
         - `realm` (string, required or wildcard `*`)
         - `interface_filter` (string pattern or regex)
         - `path_filter` (string pattern or regex)
     - `Start(ctx context.Context) error`:
       - Subscribe to `internal/engine/stream.Bus` for device events.
       - Convert incoming `stream.DeviceEvent` into `FlowMessage`:
         - `Key`: `<realm>/<device_id>`
         - `Type`: `"datastream"`, `"properties"`, or `"control"`
         - `Subtype`: Interface name (e.g. `com.example.Sensors`)
         - `Data`: Decoded payload value
         - `Metadata`: `map[string]string{"realm": realm, "device_id": deviceID, "path": path}`
         - `TimestampUs`: Event timestamp in microseconds
       - Filter out events not matching realm/interface/path filters.
       - Emit converted `FlowMessage` into output channel / router.
     - `Stop(ctx context.Context) error`:
       - Unsubscribe from engine stream bus, close output channels cleanly.

## Constraints

- Connector must be non-blocking with respect to the engine stream bus: if output queue is full, handle per flow backpressure configuration.
- Do not mutate original `stream.DeviceEvent` payloads.

## Acceptance criteria

Run from repo root:

```sh
go build ./internal/flow/blocks/astartesource/...
go test ./internal/flow/blocks/astartesource/... -run TestAstarteSource
```

Tests to add:

- `internal/flow/blocks/astartesource/source_test.go`:
  - `TestAstarteSource_IngestEvent`: publishes a mock `stream.DeviceEvent` to bus, asserts `AstarteSourceBlock` converts and emits a valid `FlowMessage` with key `<realm>/<device_id>`.
  - `TestAstarteSource_Filtering`: configures interface filter `com.example.Sensors`; publishes matching and non-matching events, verifies only matching events are emitted as `FlowMessage`s.
  - `TestAstarteSource_CleanStop`: verifies calling `Stop()` unsubscribes from stream bus without leaking subscriptions.
