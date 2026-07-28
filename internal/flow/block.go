package flow

import (
	"fmt"
)

// Block is a computation unit in a flow graph. Implementations must be safe
// for concurrent use by a single lane goroutine (one goroutine calls Process
// sequentially per message); external concurrency safety is the Router's job.
//
// Three roles exist:
//   - Source: produces messages from external events (Process receives nil, may
//     return multiple messages).
//   - Transform: consumes one message and emits zero or more transformed messages.
//   - Sink: consumes messages for external output (return value is ignored).
type Block interface {
	// Process handles one message. A source receives msg == nil and may return
	// zero or more messages. A transform receives exactly one non-nil message
	// and may return zero or more. A sink returns nil.
	Process(msg *FlowMessage) ([]*FlowMessage, error)
	// Name returns a human-readable label for metrics and logging.
	Name() string
}

// errBlock is a helper that wraps an error with the block name for diagnostics.
func errBlock(b Block, err error) error {
	return fmt.Errorf("block %s: %w", b.Name(), err)
}
