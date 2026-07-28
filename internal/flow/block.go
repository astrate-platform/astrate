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

// SourceFunc is a function that produces messages from external events. It
// receives nil and returns zero or more messages.
type SourceFunc func() ([]*FlowMessage, error)

// TransformFunc consumes one message and returns zero or more transformed
// messages.
type TransformFunc func(msg *FlowMessage) ([]*FlowMessage, error)

// SinkFunc consumes a message for external output. Return value is ignored
// by the pipeline.
type SinkFunc func(msg *FlowMessage) error

// sourceBlock adapts a SourceFunc to the Block interface.
type sourceBlock struct {
	fn   SourceFunc
	name string
}

func (s *sourceBlock) Process(_ *FlowMessage) ([]*FlowMessage, error) {
	return s.fn()
}

func (s *sourceBlock) Name() string { return s.name }

// NewSourceBlock wraps fn as a Block with the given name.
func NewSourceBlock(name string, fn SourceFunc) Block {
	return &sourceBlock{fn: fn, name: name}
}

// transformBlock adapts a TransformFunc to the Block interface.
type transformBlock struct {
	fn   TransformFunc
	name string
}

func (t *transformBlock) Process(msg *FlowMessage) ([]*FlowMessage, error) {
	return t.fn(msg)
}

func (t *transformBlock) Name() string { return t.name }

// NewTransformBlock wraps fn as a Block with the given name.
func NewTransformBlock(name string, fn TransformFunc) Block {
	return &transformBlock{fn: fn, name: name}
}

// sinkBlock adapts a SinkFunc to the Block interface.
type sinkBlock struct {
	fn   SinkFunc
	name string
}

func (s *sinkBlock) Process(msg *FlowMessage) ([]*FlowMessage, error) {
	return nil, s.fn(msg)
}

func (s *sinkBlock) Name() string { return s.name }

// NewSinkBlock wraps fn as a Block with the given name.
func NewSinkBlock(name string, fn SinkFunc) Block {
	return &sinkBlock{fn: fn, name: name}
}

// errBlock is a helper that wraps an error with the block name for diagnostics.
func errBlock(b Block, err error) error {
	return fmt.Errorf("block %s: %w", b.Name(), err)
}
