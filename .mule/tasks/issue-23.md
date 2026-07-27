# issue-23 — flow-models: FlowMessage wire format and Block/Pipeline/Flow data models

<!-- trickle-allow: internal/flow/message.go internal/flow/block.go internal/flow/pipeline.go internal/flow/flow.go internal/flow/message_test.go internal/flow/pipeline_test.go internal/flow/flow_test.go -->

## Context

Astarte Flow's core abstractions — `FlowMessage` (`astarte_flow/message/v0.1`), `Block` interface, `Pipeline` definition graph, and `Flow` status model — must be implemented in `internal/flow/`.

For a small AI worker model: **Write the test suite first**, verify that it fails without the implementation, then write the implementation until all assertions pass.

## Complete Implementation Specification

### 1. `internal/flow/message.go`

```go
package flow

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidFlowMessage = errors.New("invalid flow message: missing required fields")
)

type FlowMessage struct {
	ID          string            `json:"id"`
	Key         string            `json:"key"`
	Type        string            `json:"type"`
	Subtype     string            `json:"subtype"`
	Data        any               `json:"data"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	TimestampUs int64             `json:"timestamp_us"`
}

func (m *FlowMessage) Validate() error {
	if m.ID == "" || m.Key == "" || m.Type == "" {
		return fmt.Errorf("%w: id, key, and type are required", ErrInvalidFlowMessage)
	}
	return nil
}
```

### 2. `internal/flow/block.go`

```go
package flow

import "context"

type BlockType string

const (
	BlockTypeSource    BlockType = "source"
	BlockTypeTransform BlockType = "transform"
	BlockTypeSink      BlockType = "sink"
)

type Block interface {
	ID() string
	Type() BlockType
	Init(ctx context.Context, config map[string]any) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type SourceBlock interface {
	Block
	Out() <-chan FlowMessage
}

type SinkBlock interface {
	Block
	In() chan<- FlowMessage
}

type TransformBlock interface {
	Block
	Process(ctx context.Context, in FlowMessage) ([]FlowMessage, error)
}
```

### 3. `internal/flow/pipeline.go`

```go
package flow

import "errors"

var ErrPipelineInvalid = errors.New("pipeline definition is invalid")

type BlockConfig struct {
	ID     string         `json:"id"`
	Type   string         `json:"type"`
	Config map[string]any `json:"config,omitempty"`
}

type ConnectionConfig struct {
	FromBlock  string `json:"from_block"`
	FromOutput string `json:"from_output,omitempty"`
	ToBlock    string `json:"to_block"`
	ToInput    string `json:"to_input,omitempty"`
}

type Pipeline struct {
	ID          string             `json:"id"`
	Realm       string             `json:"realm"`
	Name        string             `json:"name"`
	Blocks      []BlockConfig      `json:"blocks"`
	Connections []ConnectionConfig `json:"connections"`
}
```

### 4. `internal/flow/flow.go`

```go
package flow

import "time"

type FlowStatus string

const (
	FlowStatusCreating FlowStatus = "creating"
	FlowStatusRunning  FlowStatus = "running"
	FlowStatusStopped  FlowStatus = "stopped"
	FlowStatusFailed   FlowStatus = "failed"
)

type Flow struct {
	ID           string         `json:"id"`
	Realm        string         `json:"realm"`
	PipelineID   string         `json:"pipeline_id"`
	Status       FlowStatus     `json:"status"`
	Config       map[string]any `json:"config,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}
```

---

## Test Suite Specifications (Loud-Failing Assertions)

Create `internal/flow/message_test.go` with exact loud-failing test logic:

```go
package flow_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/astrate-platform/astrate/internal/flow"
)

func TestFlowMessage_UpstreamWireFormatJSON(t *testing.T) {
	rawJSON := []byte(`{
		"id": "msg-12345",
		"key": "realm-test/device-001",
		"type": "incoming_data",
		"subtype": "com.example.Sensors",
		"data": 42.5,
		"metadata": {
			"path": "/temp"
		},
		"timestamp_us": 1690000000000000
	}`)

	var msg flow.FlowMessage
	if err := json.Unmarshal(rawJSON, &msg); err != nil {
		t.Fatalf("FAIL: Unmarshal failed for valid upstream astarte_flow/message/v0.1 JSON: %v", err)
	}

	if msg.ID != "msg-12345" {
		t.Fatalf("FAIL: expected ID 'msg-12345', got %q", msg.ID)
	}
	if msg.Key != "realm-test/device-001" {
		t.Fatalf("FAIL: expected Key 'realm-test/device-001', got %q", msg.Key)
	}
	if msg.Type != "incoming_data" {
		t.Fatalf("FAIL: expected Type 'incoming_data', got %q", msg.Type)
	}
	if msg.Subtype != "com.example.Sensors" {
		t.Fatalf("FAIL: expected Subtype 'com.example.Sensors', got %q", msg.Subtype)
	}
	if msg.TimestampUs != 1690000000000000 {
		t.Fatalf("FAIL: expected TimestampUs 1690000000000000, got %d", msg.TimestampUs)
	}
	if msg.Metadata["path"] != "/temp" {
		t.Fatalf("FAIL: expected Metadata['path'] '/temp', got %q", msg.Metadata["path"])
	}

	// Re-marshal and ensure exact key names
	encoded, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("FAIL: Marshal failed: %v", err)
	}

	var resultMap map[string]any
	if err := json.Unmarshal(encoded, &resultMap); err != nil {
		t.Fatalf("FAIL: Unmarshal encoded map failed: %v", err)
	}

	if _, ok := resultMap["timestamp_us"]; !ok {
		t.Fatalf("FAIL: marshaled JSON missing exact key 'timestamp_us'. Got keys: %v", resultMap)
	}
}

func TestFlowMessage_Validation(t *testing.T) {
	invalidMsg := flow.FlowMessage{ID: "123"}
	if err := invalidMsg.Validate(); !errors.Is(err, flow.ErrInvalidFlowMessage) {
		t.Fatalf("FAIL: expected ErrInvalidFlowMessage for missing key/type, got %v", err)
	}

	validMsg := flow.FlowMessage{ID: "123", Key: "k", Type: "t"}
	if err := validMsg.Validate(); err != nil {
		t.Fatalf("FAIL: expected nil error for valid message, got %v", err)
	}
}
```

Create `internal/flow/pipeline_test.go`:

```go
package flow_test

import (
	"encoding/json"
	"testing"

	"github.com/astrate-platform/astrate/internal/flow"
)

func TestPipeline_JSONSerialization(t *testing.T) {
	p := flow.Pipeline{
		ID:    "pipe-1",
		Realm: "testrealm",
		Name:  "sensor-pipeline",
		Blocks: []flow.BlockConfig{
			{ID: "src", Type: "source"},
			{ID: "snk", Type: "sink"},
		},
		Connections: []flow.ConnectionConfig{
			{FromBlock: "src", ToBlock: "snk"},
		},
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("FAIL: Marshal pipeline failed: %v", err)
	}

	var decoded flow.Pipeline
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("FAIL: Unmarshal pipeline failed: %v", err)
	}

	if len(decoded.Blocks) != 2 {
		t.Fatalf("FAIL: expected 2 blocks, got %d", len(decoded.Blocks))
	}
	if decoded.Connections[0].FromBlock != "src" || decoded.Connections[0].ToBlock != "snk" {
		t.Fatalf("FAIL: connection mapping corrupted after round-trip: %+v", decoded.Connections[0])
	}
}
```

---

## Negative Constraints & TDD Workflow

1. ❌ **Do NOT add external dependencies to `go.mod`**.
2. ❌ **Do NOT change existing files outside `internal/flow/`**.
3. **Execution order**:
   - Write test files `internal/flow/message_test.go` and `pipeline_test.go`.
   - Run `go test ./internal/flow/...` (must fail).
   - Write implementation files `message.go`, `block.go`, `pipeline.go`, `flow.go`.
   - Run `go test -v ./internal/flow/...` (must pass clean).

## Acceptance criteria

Run from repo root:

```sh
go build ./internal/flow/...
go test -v ./internal/flow/...
```
