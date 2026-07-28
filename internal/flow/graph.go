package flow

import "fmt"

// BlockGraph is a linear chain of blocks: source → transform₁ → … → sink.
// The graph is immutable after construction; calling Run feeds one message
// through every stage.
type BlockGraph struct {
	blocks []Block
}

// NewBlockGraph validates that the graph has at least one block and that the
// last block is a sink (it may return nil messages). The first block is the
// source. Returns an error if the chain is empty or nil blocks are present.
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

// Run feeds one source message through every stage. Returns the messages
// produced by the final (sink) stage — typically nil.
func (g *BlockGraph) Run(msg *FlowMessage) ([]*FlowMessage, error) {
	cur := []*FlowMessage{msg}
	for _, b := range g.blocks {
		var next []*FlowMessage
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
