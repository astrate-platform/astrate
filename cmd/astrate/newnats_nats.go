//go:build nats

package main

import (
	"github.com/astrate-platform/astrate/internal/config"
	"github.com/astrate-platform/astrate/internal/engine/forward"
	"github.com/astrate-platform/astrate/internal/engine/triggers"
)

// natsBuildTagEnabled lets untagged test files (forward_test.go) tell which
// half of the build-tag split they are running against.
const natsBuildTagEnabled = true

// newNATSForwarder builds the NATS bus forwarder. Only compiled into a binary
// built with -tags nats; see newnats_default.go for the untagged build.
func newNATSForwarder(f config.ForwardConfig) (triggers.Forwarder, error) {
	return forward.NewNATS(forward.NATSConfig{URL: f.URL, Subject: f.Subject})
}
