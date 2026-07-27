//go:build nats

package main

import (
	"context"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/config"
	"github.com/nats-io/nats.go"
	tcnats "github.com/testcontainers/testcontainers-go/modules/nats"
)

// TestForwarderNATSBuilt proves the config→newForwarder→forward.NewNATS wire
// end to end against a real broker, the -tags nats counterpart of
// TestForwarderNATSWithoutBuildTag's untagged-build error path.
func TestForwarderNATSBuilt(t *testing.T) {
	ctx := context.Background()
	ctr, err := tcnats.Run(ctx, "nats:2.10-alpine")
	if err != nil {
		t.Fatalf("failed to start NATS container: %v", err)
	}
	t.Cleanup(func() { _ = ctr.Terminate(ctx) })
	addr, err := ctr.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("ConnectionString: %v", err)
	}

	var cfg config.Config
	cfg.Triggers.Forward = config.ForwardConfig{Kind: "nats", URL: addr, Subject: "astrate.triggers"}
	f, err := newForwarder(cfg, quietLogger())
	if err != nil {
		t.Fatalf("newForwarder: %v", err)
	}
	if f == nil {
		t.Fatal("forwarder = nil, want a NATS forwarder")
	}

	nc, err := nats.Connect(addr)
	if err != nil {
		t.Fatalf("nats.Connect: %v", err)
	}
	defer nc.Close()
	sub, err := nc.SubscribeSync("astrate.triggers")
	if err != nil {
		t.Fatalf("SubscribeSync: %v", err)
	}
	if err := nc.Flush(); err != nil {
		t.Fatalf("Flush: %v", err)
	}

	if err := f.Forward(ctx, "realm1", "my.trigger", nil, nil); err != nil {
		t.Fatalf("Forward: %v", err)
	}
	if _, err := sub.NextMsg(5 * time.Second); err != nil {
		t.Fatalf("NextMsg: %v", err)
	}
}
