package channels

import (
	"encoding/json"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/engine/triggers"
)

type fakeBus struct {
	ch          chan stream.Event
	cancelCount atomic.Int64
}

func (f *fakeBus) Subscribe(_ string, _ stream.Filter, _ int) (<-chan stream.Event, func()) {
	return f.ch, func() { f.cancelCount.Add(1) }
}

func (f *fakeBus) cancelCalled() bool {
	return f.cancelCount.Load() > 0
}

const (
	validDeviceIDA = "f0VMRgIBAQAAAAAAAAAAAA"
	validDeviceIDB = "f0VMRgIBAQAAAAAAAAAAAQ"
)

func TestDashboardFourWatches(t *testing.T) {
	bus := &fakeBus{ch: make(chan stream.Event, 16)}
	reg := NewRegistry(bus)
	rm := reg.Join("test", "rooms:test:devices")

	connReq := WatchRequest{
		Name:          "connectiontrigger-D",
		DeviceID:      validDeviceIDA,
		SimpleTrigger: json.RawMessage(`{"type":"device_trigger","on":"device_connected","device_id":"` + validDeviceIDA + `"}`),
	}
	discReq := WatchRequest{
		Name:          "disconnectiontrigger-D",
		DeviceID:      validDeviceIDA,
		SimpleTrigger: json.RawMessage(`{"type":"device_trigger","on":"device_disconnected","device_id":"` + validDeviceIDA + `"}`),
	}
	errReq := WatchRequest{
		Name:          "errortrigger-D",
		DeviceID:      validDeviceIDA,
		SimpleTrigger: json.RawMessage(`{"type":"device_trigger","on":"device_error","device_id":"` + validDeviceIDA + `"}`),
	}
	dataReq := WatchRequest{
		Name:          "datatrigger-D",
		DeviceID:      validDeviceIDA,
		SimpleTrigger: json.RawMessage(`{"type":"data_trigger","on":"incoming_data","interface_name":"*","value_match_operator":"*","match_path":"/*"}`),
	}

	for _, req := range []WatchRequest{connReq, discReq, errReq, dataReq} {
		if err := rm.Watch(req); err != nil {
			t.Fatalf("Watch(%s): %v", req.Name, err)
		}
	}

	if rm.Watches() != 4 {
		t.Fatalf("Watches() = %d, want 4", rm.Watches())
	}

	m := rm.AddMember(0)
	defer m.Leave()

	now := time.Now().UTC().Truncate(time.Millisecond)

	bus.ch <- stream.Event{Kind: stream.KindDeviceConnected, DeviceID: validDeviceIDA, IP: "1.2.3.4", Timestamp: now}
	bus.ch <- stream.Event{Kind: stream.KindDeviceDisconnected, DeviceID: validDeviceIDA, Timestamp: now}
	bus.ch <- stream.Event{Kind: stream.KindDeviceError, DeviceID: validDeviceIDA, ErrorName: "unknown_format", ErrorMetadata: map[string]string{"detail": "slow"}, Timestamp: now}
	bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: validDeviceIDA, Interface: "org.example.V1", Path: "/sensor/temp", Value: 42.0, InterfaceMajor: 1, Timestamp: now}

	type expected struct {
		TriggerName string
		EventType   string
	}
	want := []expected{
		{"connectiontrigger-D", "device_connected"},
		{"disconnectiontrigger-D", "device_disconnected"},
		{"errortrigger-D", "device_error"},
		{"datatrigger-D", "incoming_data"},
	}

	for _, w := range want {
		select {
		case ev := <-m.Events():
			if ev.TriggerName != w.TriggerName {
				t.Errorf("got TriggerName %q, want %q", ev.TriggerName, w.TriggerName)
			}
			switch body := ev.Event.(type) {
			case triggers.DeviceConnectedEvent:
				if w.EventType != "device_connected" {
					t.Errorf("got DeviceConnectedEvent, want %s", w.EventType)
				}
				if body.DeviceIPAddress != "1.2.3.4" {
					t.Errorf("DeviceIPAddress = %q, want %q", body.DeviceIPAddress, "1.2.3.4")
				}
			case triggers.DeviceDisconnectedEvent:
				if w.EventType != "device_disconnected" {
					t.Errorf("got DeviceDisconnectedEvent, want %s", w.EventType)
				}
			case triggers.DeviceErrorEvent:
				if w.EventType != "device_error" {
					t.Errorf("got DeviceErrorEvent, want %s", w.EventType)
				}
				// The bus carries Astrate's own reject reason; the event
				// body must carry the upstream enum value the client
				// validates against (triggers.UpstreamErrorName).
				if body.ErrorName != "undecodable_bson_payload" {
					t.Errorf("ErrorName = %q, want %q", body.ErrorName, "undecodable_bson_payload")
				}
				if body.Metadata["astrate_reason"] != "unknown_format" {
					t.Errorf("astrate_reason = %q, want %q", body.Metadata["astrate_reason"], "unknown_format")
				}
			case triggers.IncomingDataEvent:
				if w.EventType != "incoming_data" {
					t.Errorf("got IncomingDataEvent, want %s", w.EventType)
				}
				if body.Interface != "org.example.V1" {
					t.Errorf("Interface = %q, want %q", body.Interface, "org.example.V1")
				}
				if body.Path != "/sensor/temp" {
					t.Errorf("Path = %q, want %q", body.Path, "/sensor/temp")
				}
			default:
				t.Errorf("unexpected event type %T", ev.Event)
			}
		case <-time.After(time.Second):
			t.Fatalf("timeout waiting for %s", w.TriggerName)
		}
	}
}

func TestDeviceIDFilterIsAdditive(t *testing.T) {
	bus := &fakeBus{ch: make(chan stream.Event, 16)}
	reg := NewRegistry(bus)
	rm := reg.Join("test", "rooms:test:devices")

	dataReq := WatchRequest{
		Name:          "datatrigger-D",
		DeviceID:      validDeviceIDA,
		SimpleTrigger: json.RawMessage(`{"type":"data_trigger","on":"incoming_data","interface_name":"*","value_match_operator":"*","match_path":"/*"}`),
	}
	if err := rm.Watch(dataReq); err != nil {
		t.Fatalf("Watch: %v", err)
	}

	m := rm.AddMember(0)
	defer m.Leave()

	now := time.Now().UTC().Truncate(time.Millisecond)
	bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: validDeviceIDB, Interface: "org.example.V1", Path: "/sensor/temp", Value: 42.0, InterfaceMajor: 1, Timestamp: now}

	select {
	case ev := <-m.Events():
		t.Fatalf("unexpected event for wrong device: %+v", ev)
	case <-time.After(200 * time.Millisecond):
	}
}

func TestFullMailboxDropsNeverBlocks(t *testing.T) {
	bus := &fakeBus{ch: make(chan stream.Event, 16)}
	reg := NewRegistry(bus)
	rm := reg.Join("test", "rooms:test:devices")

	dataReq := WatchRequest{
		Name:          "datatrigger",
		SimpleTrigger: json.RawMessage(`{"type":"data_trigger","on":"incoming_data","interface_name":"*","value_match_operator":"*","match_path":"/*"}`),
	}
	if err := rm.Watch(dataReq); err != nil {
		t.Fatalf("Watch: %v", err)
	}

	slow := rm.AddMember(1)
	fast := rm.AddMember(16)

	now := time.Now().UTC().Truncate(time.Millisecond)
	for i := 0; i < 3; i++ {
		bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: "X", Interface: "org.example.V1", Path: "/a", Value: i, InterfaceMajor: 1, Timestamp: now}
	}

	// Wait for dispatch to finish by polling Dropped().
	deadline := time.After(2 * time.Second)
	for slow.Dropped() < 2 {
		select {
		case <-deadline:
			t.Fatalf("timeout waiting for drops: Dropped() = %d", slow.Dropped())
		default:
			time.Sleep(time.Millisecond)
		}
	}

	// Drain slow (1 readable).
	select {
	case <-slow.Events():
	case <-time.After(time.Second):
		t.Fatal("slow member: timeout")
	}

	// Fast member got all 3.
	var fastCount int
drain:
	for {
		select {
		case <-fast.Events():
			fastCount++
		default:
			break drain
		}
	}
	if fastCount != 3 {
		t.Fatalf("fast received %d, want 3", fastCount)
	}

	slow.Leave()
	fast.Leave()
}

func TestLastMemberTearsDown(t *testing.T) {
	bus := &fakeBus{ch: make(chan stream.Event, 16)}
	reg := NewRegistry(bus)
	rm := reg.Join("test", "rooms:test:devices")

	m1 := rm.AddMember(0)
	m2 := rm.AddMember(0)

	m1.Leave()
	if reg.Rooms() != 1 {
		t.Fatalf("after first Leave, Rooms() = %d, want 1", reg.Rooms())
	}
	if bus.cancelCalled() {
		t.Fatal("cancel should not have run yet")
	}

	m2.Leave()
	if reg.Rooms() != 0 {
		t.Fatalf("after second Leave, Rooms() = %d, want 0", reg.Rooms())
	}
	if !bus.cancelCalled() {
		t.Fatal("cancel should have run")
	}
}

func TestDoubleLeaveNoPanic(_ *testing.T) {
	bus := &fakeBus{ch: make(chan stream.Event, 16)}
	reg := NewRegistry(bus)
	rm := reg.Join("test", "rooms:test:devices")

	m := rm.AddMember(0)
	m.Leave()
	m.Leave()
}

func TestRejectedWatches(t *testing.T) {
	bus := &fakeBus{ch: make(chan stream.Event, 16)}
	reg := NewRegistry(bus)
	rm := reg.Join("test", "rooms:test:devices")

	tests := []struct {
		name    string
		req     WatchRequest
		wantErr string
	}{
		{
			name:    "empty name",
			req:     WatchRequest{Name: "", DeviceID: validDeviceIDA, SimpleTrigger: json.RawMessage(`{"type":"device_trigger","on":"device_connected","device_id":"` + validDeviceIDA + `"}`)},
			wantErr: "name",
		},
		{
			name:    "missing simple_trigger",
			req:     WatchRequest{Name: "t", DeviceID: validDeviceIDA},
			wantErr: "simple_trigger",
		},
		{
			name:    "invalid on value",
			req:     WatchRequest{Name: "t", DeviceID: validDeviceIDA, SimpleTrigger: json.RawMessage(`{"type":"device_trigger","on":"not_a_real_event","device_id":"` + validDeviceIDA + `"}`)},
			wantErr: "CompileCondition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := rm.Watches()
			err := rm.Watch(tt.req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if tt.wantErr != "" && !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error %q does not mention %q", err.Error(), tt.wantErr)
			}
			if rm.Watches() != before {
				t.Fatalf("Watches() grew to %d after rejected watch", rm.Watches())
			}
		})
	}
}

func TestBusCloseClosesMailboxes(t *testing.T) {
	bus := &fakeBus{ch: make(chan stream.Event)}
	reg := NewRegistry(bus)
	rm := reg.Join("test", "rooms:test:devices")

	m := rm.AddMember(0)
	close(bus.ch)

	select {
	case _, ok := <-m.Events():
		if ok {
			t.Fatal("expected closed channel")
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for mailbox close")
	}

	// The room retires itself: a topic whose bus subscription has died must not
	// stay in the registry, or a later Join would hand out a room with a dead
	// event channel and no dispatcher behind it.
	deadline := time.After(time.Second)
	for reg.Rooms() != 0 || !bus.cancelCalled() {
		select {
		case <-deadline:
			t.Fatalf("room not retired: Rooms() = %d, cancelCalled = %v", reg.Rooms(), bus.cancelCalled())
		default:
			time.Sleep(time.Millisecond)
		}
	}
}
