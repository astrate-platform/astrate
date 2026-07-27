# issue-27 — flow-astarte-source: Astarte device events source block

<!-- trickle-allow: internal/flow/blocks/astartesource/source.go internal/flow/blocks/astartesource/source_test.go internal/engine/stream/bus.go -->

## Context

The `AstarteSource` block is a built-in `flow.SourceBlock` that ingests device events directly from Astrate's live engine stream bus (`internal/engine/stream`) and converts them into `FlowMessage`s.

This specification provides exact struct definitions, filter matching against `stream.Filter`, conversion rules, and a **loud-failing test suite**.

## Exact Code Specifications

### 1. `internal/flow/blocks/astartesource/source.go`

```go
package astartesource

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/flow"
)

type Config struct {
	Realm           string `json:"realm"`
	InterfaceFilter string `json:"interface_filter"`
	PathFilter      string `json:"path_filter"`
}

type AstarteSourceBlock struct {
	id     string
	cfg    Config
	bus    *stream.Bus
	out    chan flow.FlowMessage
	cancel func()
	ctx    context.Context
	wg     sync.WaitGroup
}

func New(id string, bus *stream.Bus) *AstarteSourceBlock {
	return &AstarteSourceBlock{
		id:  id,
		bus: bus,
		out: make(chan flow.FlowMessage, 100),
	}
}

func (b *AstarteSourceBlock) ID() string { return b.id }
func (b *AstarteSourceBlock) Type() flow.BlockType { return flow.BlockTypeSource }
func (b *AstarteSourceBlock) Out() <-chan flow.FlowMessage { return b.out }

func (b *AstarteSourceBlock) Init(ctx context.Context, config map[string]any) error {
	b.cfg.Realm, _ = config["realm"].(string)
	b.cfg.InterfaceFilter, _ = config["interface_filter"].(string)
	b.cfg.PathFilter, _ = config["path_filter"].(string)
	return nil
}

func (b *AstarteSourceBlock) Start(ctx context.Context) error {
	filter := stream.Filter{
		Interface: b.cfg.InterfaceFilter,
	}

	eventChan, cancel := b.bus.Subscribe(b.cfg.Realm, filter, 64)
	b.cancel = cancel

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for ev := range eventChan {
			b.processEvent(ev)
		}
	}()

	return nil
}

func (b *AstarteSourceBlock) processEvent(ev stream.Event) {
	if b.cfg.PathFilter != "" && b.cfg.PathFilter != "*" && !strings.Contains(ev.Path, b.cfg.PathFilter) {
		return
	}

	ts := ev.Timestamp.UnixNano() / 1000
	if ts == 0 {
		ts = 1
	}

	msg := flow.FlowMessage{
		ID:          fmt.Sprintf("%s-%s-%d", ev.DeviceID, ev.Interface, ts),
		Key:         fmt.Sprintf("%s/%s", ev.Realm, ev.DeviceID),
		Type:        ev.Kind,
		Subtype:     ev.Interface,
		Data:        ev.Value,
		Metadata: map[string]string{
			"realm":     ev.Realm,
			"device_id": ev.DeviceID,
			"interface": ev.Interface,
			"path":      ev.Path,
		},
		TimestampUs: ts,
	}

	select {
	case b.out <- msg:
	default:
		// Non-blocking drop if output channel buffer full (§1.4 philosophy)
	}
}

func (b *AstarteSourceBlock) Stop(ctx context.Context) error {
	if b.cancel != nil {
		b.cancel()
	}
	b.wg.Wait()
	close(b.out)
	return nil
}
```

---

## Test Suite Specifications (Loud-Failing Assertions)

Create `internal/flow/blocks/astartesource/source_test.go`:

```go
package astartesource_test

import (
	"context"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/flow/blocks/astartesource"
)

func TestAstarteSource_IngestAndConvert(t *testing.T) {
	bus := stream.New(nil)
	block := astartesource.New("source-1", bus)
	ctx := context.Background()

	_ = block.Init(ctx, map[string]any{
		"realm":            "testrealm",
		"interface_filter": "com.example.Sensors",
	})

	if err := block.Start(ctx); err != nil {
		t.Fatalf("FAIL: Start failed: %v", err)
	}

	now := time.Now()
	event := stream.Event{
		Kind:      stream.KindIncomingData,
		Realm:     "testrealm",
		DeviceID:  "device-123",
		Interface: "com.example.Sensors",
		Path:      "/temp",
		Value:     36.6,
		Timestamp: now,
	}

	bus.Publish(event)

	select {
	case msg := <-block.Out():
		if msg.Key != "testrealm/device-123" {
			t.Fatalf("FAIL: expected Key 'testrealm/device-123', got %q", msg.Key)
		}
		if msg.Type != stream.KindIncomingData {
			t.Fatalf("FAIL: expected Type 'incoming_data', got %q", msg.Type)
		}
		if msg.Subtype != "com.example.Sensors" {
			t.Fatalf("FAIL: expected Subtype 'com.example.Sensors', got %q", msg.Subtype)
		}
		if msg.Data != 36.6 {
			t.Fatalf("FAIL: expected Data 36.6, got %v", msg.Data)
		}
		if msg.Metadata["path"] != "/temp" {
			t.Fatalf("FAIL: expected Metadata['path'] '/temp', got %q", msg.Metadata["path"])
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("FAIL TIMEOUT: AstarteSource did not emit FlowMessage for published event")
	}

	if err := block.Stop(ctx); err != nil {
		t.Fatalf("FAIL: Stop failed: %v", err)
	}

	if bus.Subscribers() != 0 {
		t.Fatalf("FAIL: subscription not cleaned up on Stop; subscribers count = %d", bus.Subscribers())
	}
}
```

---

## Negative Constraints & TDD Workflow

1. ❌ **Do NOT block the stream bus publisher**.
2. ❌ **Do NOT re-invent event structures**. Import `github.com/astrate-platform/astrate/internal/engine/stream`.
3. **Execution order**:
   - Write `internal/flow/blocks/astartesource/source_test.go`.
   - Run `go test ./internal/flow/blocks/astartesource/...` (must fail).
   - Write `internal/flow/blocks/astartesource/source.go`.
   - Run `go test -v ./internal/flow/blocks/astartesource/...` (must pass clean).

## Acceptance criteria

Run from repo root:

```sh
go build ./internal/flow/blocks/astartesource/...
go test -v ./internal/flow/blocks/astartesource/...
```
