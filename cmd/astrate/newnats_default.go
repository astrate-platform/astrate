//go:build !nats

package main

import (
	"fmt"

	"github.com/astrate-platform/astrate/internal/config"
	"github.com/astrate-platform/astrate/internal/engine/triggers"
)

// natsBuildTagEnabled lets untagged test files (forward_test.go) tell which
// half of the build-tag split they are running against.
const natsBuildTagEnabled = false

// newNATSForwarder reports that this binary was not built with -tags nats, so
// triggers.forward.kind = "nats" fails at boot rather than silently
// forwarding nothing. See newnats_nats.go for the tagged build.
func newNATSForwarder(_ config.ForwardConfig) (triggers.Forwarder, error) {
	return nil, fmt.Errorf("triggers.forward.kind \"nats\" requires building astrate with -tags nats")
}
