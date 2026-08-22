package flow

import "fmt"

// BlockGraph is a linear chain of blocks: source → transform₁ → … → sink.
// The graph is immutable after construction; calling Run feeds one message
// through every non-Source stage. Source blocks are driven by the flow
// source pump (see Manager.StartFlow), not by Run.
type BlockGraph struct {
	blocks []Block
}

// NewBlockGraph validates that the graph has at least one block and that the
// last block is a sink (it may return nil messages). The first block is
// typically the source. Returns an error if the chain is empty or nil blocks
// are present.
func NewBlockGraph(blocks ...Block) (*BlockGraph, error) {
	if len(blocks) == 0 {
		return nil, fmt.Errorf("flow: block graph must have at least one block")
	}
	for i, b := range blocks {
		if b == nil {
			return nil, fmt.Errorf("flow: block at index %d is nil", i)
		}
	}
	return &BlockGraph{blocks: blocks}, nil
}

// Blocks returns the graph's blocks in order. The slice must not be mutated.
func (g *BlockGraph) Blocks() []Block { return g.blocks }

// Sources returns every block that implements Source, in graph order.
func (g *BlockGraph) Sources() []Source {
	var out []Source
	for _, b := range g.blocks {
		if s, ok := b.(Source); ok {
			out = append(out, s)
		}
	}
	return out
}

// Run feeds one message through every non-Source stage. Source stages are
// skipped: they are polled by the source pump and their outputs are
// submitted into the router, which calls Run. Returns the messages produced
// by the final stage — typically nil for a sink.
func (g *BlockGraph) Run(msg *Message) ([]*Message, error) {
	cur := []*Message{msg}
	for _, b := range g.blocks {
		if _, isSource := b.(Source); isSource {
			continue
		}
		var next []*Message
		for _, m := range cur {
			out, err := b.Process(m)
			if err != nil {
				return nil, errBlock(b, err)
			}
			next = append(next, out...)
		}
		cur = next
	}
	return cur, nil
}
