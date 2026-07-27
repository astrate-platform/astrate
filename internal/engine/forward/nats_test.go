//go:build nats

package forward

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	tcnats "github.com/testcontainers/testcontainers-go/modules/nats"
)

// Assert at compile time that *NATS satisfies the triggers.Forwarder contract
// without importing the triggers package (which would create an import cycle).
// forwarder and bodyShape are defined in http_test.go.
var _ forwarder = (*NATS)(nil)

// startContainer starts a NATS container and returns a cleanup function and
// the connection string.
func startContainer(t *testing.T) string {
	t.Helper()
	ctx := context.Background()
	ctr, err := tcnats.Run(ctx, "nats:2.10-alpine")
	if err != nil {
		t.Fatalf("failed to start NATS container: %v", err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(ctx) })
	addr, err := ctr.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}
	return addr
}

func TestNATSForwardHappyPath(t *testing.T) {
	addr := startContainer(t)

	fwd, err := NewNATS(NATSConfig{URL: addr, Subject: "test.subject"})
	if err != nil {
		t.Fatalf("NewNATS: %v", err)
	}

	ctx := context.Background()
	nc, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("nats.Connect: %v", err)
	}
	defer nc.Close()
	sub, err := nc.SubscribeSync("test.subject")
	if err != nil {
		t.Fatalf("SubscribeSync: %v", err)
	}
	if err := nc.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	action := json.RawMessage(`{"amqp_exchange":"x"}`)
	event := []byte(`{"device_id":"abc"}`)
	if err := fwd.Forward(ctx, "realm1", "my.trigger", action, event); err != nil {
		t.Fatalf("Forward: %v", err)
	}

	msg, err := sub.NextMsg(5 * time.Second)
	if err != nil {
		t.Fatalf("NextMsg: %v", err)
	}

	var got bodyShape
	if err := json.Unmarshal(msg.Data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got.Realm != "realm1" {
		t.Errorf("realm = %q, want %q", got.Realm, "realm1")
	}
	if got.Trigger != "my.trigger" {
		t.Errorf("trigger = %q, want %q", got.Trigger, "my.trigger")
	}
	if string(got.Action) != string(action) {
		t.Errorf("action = %s, want %s", got.Action, action)
	}
	if string(got.Event) != string(event) {
		t.Errorf("event = %s, want %s", got.Event, event)
	}
}

func TestNATSForwardNilActionAndEvent(t *testing.T) {
	addr := startContainer(t)

	fwd, err := NewNATS(NATSConfig{URL: addr, Subject: "nil.subject"})
	if err != nil {
		t.Fatalf("NewNATS: %v", err)
	}

	ctx := context.Background()
	nc, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("nats.Connect: %v", err)
	}
	defer nc.Close()
	sub, err := nc.SubscribeSync("nil.subject")
	if err != nil {
		t.Fatalf("SubscribeSync: %v", err)
	}
	if err := nc.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	if err := fwd.Forward(ctx, "r", "t", nil, nil); err != nil {
		t.Fatalf("Forward: %v", err)
	}

	msg, err := sub.NextMsg(5 * time.Second)
	if err != nil {
		t.Fatalf("NextMsg: %v", err)
	}

	var got bodyShape
	if err := json.Unmarshal(msg.Data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if string(got.Action) != "null" {
		t.Errorf("action = %s, want null", got.Action)
	}
	if string(got.Event) != "null" {
		t.Errorf("event = %s, want null", got.Event)
	}
}

func TestNATSForwardEmptyNonNilActionAndEvent(t *testing.T) {
	addr := startContainer(t)

	fwd, err := NewNATS(NATSConfig{URL: addr, Subject: "empty.subject"})
	if err != nil {
		t.Fatalf("NewNATS: %v", err)
	}

	ctx := context.Background()
	nc, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("nats.Connect: %v", err)
	}
	defer nc.Close()
	sub, err := nc.SubscribeSync("empty.subject")
	if err != nil {
		t.Fatalf("SubscribeSync: %v", err)
	}
	if err := nc.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	if err := fwd.Forward(ctx, "r", "t", json.RawMessage{}, []byte{}); err != nil {
		t.Fatalf("Forward: %v", err)
	}

	msg, err := sub.NextMsg(5 * time.Second)
	if err != nil {
		t.Fatalf("NextMsg: %v", err)
	}

	var got bodyShape
	if err := json.Unmarshal(msg.Data, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if string(got.Action) != "null" {
		t.Errorf("action = %s, want null", got.Action)
	}
	if string(got.Event) != "null" {
		t.Errorf("event = %s, want null", got.Event)
	}
}

func TestNewNATSRejectsBadConfig(t *testing.T) {
	addr := startContainer(t)

	// Empty URL
	if _, err := NewNATS(NATSConfig{Subject: "s"}); err == nil {
		t.Error("empty URL: got nil error, want non-nil")
	}

	// Empty Subject
	if _, err := NewNATS(NATSConfig{URL: addr}); err == nil {
		t.Error("empty Subject: got nil error, want non-nil")
	}

	// Unreachable URL
	if _, err := NewNATS(NATSConfig{URL: "nats://127.0.0.1:1", Subject: "s"}); err == nil {
		t.Error("unreachable URL: got nil error, want non-nil")
	}

	// Valid config succeeds
	if _, err := NewNATS(NATSConfig{URL: addr, Subject: "s"}); err != nil {
		t.Errorf("valid config: got error %v, want nil", err)
	}
}
