package flow

import (
	"context"
	"fmt"
)

// Block is a computation unit in a flow graph. Implementations must be safe
// for concurrent use by a single lane goroutine (one goroutine calls Process
// sequentially per message); external concurrency safety is the Router's job.
//
// Three roles exist:
//   - Source: produces messages from external events (see the Source interface).
//   - Transform: consumes one message and emits zero or more transformed messages.
//   - Sink: consumes messages for external output (return value is ignored).
type Block interface {
	// Process handles one message. A source receives msg == nil and may return
	// zero or more messages. A transform receives exactly one non-nil message
	// and may return zero or more. A sink returns nil.
	Process(msg *Message) ([]*Message, error)
	// Name returns a human-readable label for metrics and logging.
	Name() string
}

// Source is a Block that produces messages from an external system without an
// input message. The flow source pump calls Emit on every Source in the graph
// and submits the results into the router; BlockGraph.Run skips Source stages
// so a submitted message is not re-consumed by the producer.
type Source interface {
	Block
	// Emit returns newly available messages. Implementations may block until
	// at least one message is ready or ctx is cancelled.
	Emit(ctx context.Context) ([]*Message, error)
}

// Stopper is optionally implemented by Blocks that own resources (bus
// subscriptions, goroutines, file handles). Manager.StopFlow calls Stop on
// every Stopper after the source pump exits and the router drains.
type Stopper interface {
	Stop()
}

// SourceFunc is a function that produces messages from external events. It
// receives nil and returns zero or more messages.
type SourceFunc func() ([]*Message, error)

// TransformFunc consumes one message and returns zero or more transformed
// messages.
type TransformFunc func(msg *Message) ([]*Message, error)

// SinkFunc consumes a message for external output. Return value is ignored
// by the pipeline.
type SinkFunc func(msg *Message) error

// sourceBlock adapts a SourceFunc to the Block and Source interfaces.
type sourceBlock struct {
	fn   SourceFunc
	name string
}

func (s *sourceBlock) Process(_ *Message) ([]*Message, error) {
	return s.fn()
}

func (s *sourceBlock) Emit(_ context.Context) ([]*Message, error) {
	return s.fn()
}

func (s *sourceBlock) Name() string { return s.name }

// NewSourceBlock wraps fn as a Source Block with the given name.
func NewSourceBlock(name string, fn SourceFunc) Block {
	return &sourceBlock{fn: fn, name: name}
}

// transformBlock adapts a TransformFunc to the Block interface.
type transformBlock struct {
	fn   TransformFunc
	name string
}

func (t *transformBlock) Process(msg *Message) ([]*Message, error) {
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

func (s *sinkBlock) Process(msg *Message) ([]*Message, error) {
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
