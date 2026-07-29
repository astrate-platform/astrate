package astartesource

import (
	"context"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/flow"
)

func TestSourceIngestsAndConverts(t *testing.T) {
	bus := stream.New(nil)
	src := New(bus, Config{Realm: "test"})
	defer src.Stop()

	ts := time.Date(2026, 7, 28, 12, 0, 0, 0, time.UTC)
	bus.Publish(stream.Event{
		Kind: stream.KindIncomingData, Realm: "test", DeviceID: "dev1",
		Interface: "com.ex.Sensors", Path: "/v", Value: 1.5, Timestamp: ts,
	})

	var got []*flow.FlowMessage
	deadline := time.Now().Add(2 * time.Second)
	for len(got) == 0 && time.Now().Before(deadline) {
		msgs, err := src.Process(nil)
		if err != nil {
			t.Fatalf("Process: %v", err)
		}
		got = append(got, msgs...)
		if len(got) == 0 {
			time.Sleep(time.Millisecond)
		}
	}
	if len(got) != 1 {
		t.Fatalf("got %d messages, want 1", len(got))
	}
	m := got[0]
	if m.Key != "test/dev1" {
		t.Errorf("Key = %q, want test/dev1", m.Key)
	}
	if m.Type != flow.TypeReal || m.Data.(float64) != 1.5 {
		t.Errorf("Type/Data = %v/%v, want TypeReal/1.5", m.Type, m.Data)
	}
	if m.Timestamp != ts.UnixMicro() {
		t.Errorf("Timestamp = %d, want %d", m.Timestamp, ts.UnixMicro())
	}
}

func TestSourceRealmAndInterfaceFilter(t *testing.T) {
	bus := stream.New(nil)
	src := New(bus, Config{Realm: "test", Interface: "com.ex.Sensors"})
	defer src.Stop()

	bus.Publish(stream.Event{Kind: stream.KindIncomingData, Realm: "other", DeviceID: "dev1", Interface: "com.ex.Sensors", Value: 1.0})
	bus.Publish(stream.Event{Kind: stream.KindIncomingData, Realm: "test", DeviceID: "dev1", Interface: "com.ex.Other", Value: 2.0})
	bus.Publish(stream.Event{Kind: stream.KindIncomingData, Realm: "test", DeviceID: "dev1", Interface: "com.ex.Sensors", Value: 3.0})

	time.Sleep(20 * time.Millisecond)
	msgs, err := src.Process(nil)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(msgs) != 1 || msgs[0].Data.(float64) != 3.0 {
		t.Fatalf("msgs = %+v, want one message with Data 3.0", msgs)
	}
}

func TestSourcePathFilter(t *testing.T) {
	bus := stream.New(nil)
	src := New(bus, Config{Realm: "test", Path: "/keep"})
	defer src.Stop()

	bus.Publish(stream.Event{Kind: stream.KindIncomingData, Realm: "test", DeviceID: "d", Path: "/drop/x", Value: 1.0})
	bus.Publish(stream.Event{Kind: stream.KindIncomingData, Realm: "test", DeviceID: "d", Path: "/keep/x", Value: 2.0})

	time.Sleep(20 * time.Millisecond)
	msgs, err := src.Process(nil)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if len(msgs) != 1 || msgs[0].Data.(float64) != 2.0 {
		t.Fatalf("msgs = %+v, want one message with Data 2.0", msgs)
	}
}

func TestSourceStopUnsubscribesCleanly(t *testing.T) {
	bus := stream.New(nil)
	src := New(bus, Config{Realm: "test"})
	if bus.Subscribers() != 1 {
		t.Fatalf("Subscribers = %d, want 1", bus.Subscribers())
	}
	src.Stop()
	if bus.Subscribers() != 0 {
		t.Fatalf("Subscribers after Stop = %d, want 0", bus.Subscribers())
	}
	// Process on a stopped source must return cleanly, not block or panic.
	msgs, err := src.Process(nil)
	if err != nil || msgs != nil {
		t.Errorf("Process after Stop = %v, %v, want nil, nil", msgs, err)
	}
}

func TestSourceEmitBlocksUntilEvent(t *testing.T) {
	bus := stream.New(nil)
	src := New(bus, Config{Realm: "test"})
	defer src.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan []*flow.FlowMessage, 1)
	errc := make(chan error, 1)
	go func() {
		msgs, err := src.Emit(ctx)
		if err != nil {
			errc <- err
			return
		}
		done <- msgs
	}()

	// Give Emit a moment to block before publishing.
	time.Sleep(20 * time.Millisecond)
	bus.Publish(stream.Event{
		Kind: stream.KindIncomingData, Realm: "test", DeviceID: "d",
		Interface: "com.ex.S", Path: "/v", Value: true, Timestamp: time.Now(),
	})

	select {
	case err := <-errc:
		t.Fatalf("Emit: %v", err)
	case msgs := <-done:
		if len(msgs) != 1 || msgs[0].Data != true {
			t.Fatalf("Emit = %+v, want one boolean true", msgs)
		}
	case <-ctx.Done():
		t.Fatal("Emit timed out waiting for published event")
	}
}

func TestSourceEmitCancelled(t *testing.T) {
	bus := stream.New(nil)
	src := New(bus, Config{Realm: "test"})
	defer src.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := src.Emit(ctx)
	if err == nil {
		t.Fatal("Emit with cancelled context returned nil error")
	}
}
