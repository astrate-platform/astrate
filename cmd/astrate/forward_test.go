package main

import (
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/config"
)

// quietLogger keeps newForwarder's info line out of the test output.
func quietLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// TestForwarderDisabledIsUntypedNil is the one rule in newForwarder that a
// compiler cannot catch and a careless refactor would break: the executor
// decides whether to forward by comparing its Forwarder against nil, so
// returning a typed nil *forward.HTTP would pass that comparison and then
// panic on the first custom trigger action. The interface value itself must
// be nil, which is why the disabled branch returns a bare nil rather than a
// nil pointer.
func TestForwarderDisabledIsUntypedNil(t *testing.T) {
	f, err := newForwarder(config.Config{}, quietLogger())
	if err != nil {
		t.Fatalf("newForwarder: %v", err)
	}
	if f != nil {
		t.Fatalf("forwarder = %#v, want an untyped nil interface", f)
	}
}

// TestForwarderHTTPBuilt checks the configured kind produces a usable
// forwarder, so the disabled case above is not satisfied by a function that
// returns nil for everything.
func TestForwarderHTTPBuilt(t *testing.T) {
	var cfg config.Config
	cfg.Triggers.Forward = config.ForwardConfig{
		Kind: "http",
		URL:  "https://bus.example/trigger",
	}
	f, err := newForwarder(cfg, quietLogger())
	if err != nil {
		t.Fatalf("newForwarder: %v", err)
	}
	if f == nil {
		t.Fatal("forwarder = nil, want an HTTP forwarder")
	}
}

// TestForwarderRejectsBadURL: config.validate normally catches this, but
// newForwarder must not silently ignore a construction failure if a config
// ever reaches it unvalidated.
func TestForwarderRejectsBadURL(t *testing.T) {
	var cfg config.Config
	cfg.Triggers.Forward = config.ForwardConfig{Kind: "http", URL: "nope"}
	if _, err := newForwarder(cfg, quietLogger()); err == nil {
		t.Fatal("expected an error for a relative URL")
	}
}

// TestForwarderNATSWithoutBuildTag: this file (and newnats_default.go) build
// without -tags nats, so kind = "nats" must fail at boot with a message
// naming the missing build tag, not silently forward nothing.
func TestForwarderNATSWithoutBuildTag(t *testing.T) {
	if natsBuildTagEnabled {
		t.Skip("built with -tags nats; see forward_nats_test.go")
	}
	var cfg config.Config
	cfg.Triggers.Forward = config.ForwardConfig{Kind: "nats", URL: "nats://bus.example:4222", Subject: "astrate.triggers"}
	_, err := newForwarder(cfg, quietLogger())
	if err == nil {
		t.Fatal("expected an error building a NATS forwarder without -tags nats")
	}
	if !strings.Contains(err.Error(), "-tags nats") {
		t.Errorf("error = %q, want it to name the missing build tag", err)
	}
}
