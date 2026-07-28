package flow

import (
	"encoding/json"
	"fmt"
)

// Pipeline is an acyclic graph (DAG) of named blocks with typed connections.
// It is a serialisable description; calling Manager.StartFlow instantiates it
// into a running Flow.
type Pipeline struct {
	// ID is a unique identifier for this pipeline.
	ID string `json:"id"`
	// Name is a human-readable label.
	Name string `json:"name"`
	// Blocks is the set of nodes in the graph.
	Blocks []PipelineNode `json:"blocks"`
	// Connections is the set of edges linking block output ports to input ports.
	Connections []Connection `json:"connections"`
}

// PipelineNode describes one block within a pipeline.
type PipelineNode struct {
	// Name is a unique identifier for this node within the pipeline.
	Name string `json:"name"`
	// BlockType identifies which block implementation to use.
	BlockType string `json:"block_type"`
	// Config holds block-specific parameters.
	Config map[string]any `json:"config,omitempty"`
}

// Connection describes a typed edge between two blocks in a pipeline.
type Connection struct {
	// From is the name of the source block.
	From string `json:"from"`
	// FromPort is the output port on the source block (empty means default).
	FromPort string `json:"from_port,omitempty"`
	// To is the name of the target block.
	To string `json:"to"`
	// ToPort is the input port on the target block (empty means default).
	ToPort string `json:"to_port,omitempty"`
}

// Validate checks that the pipeline is a valid acyclic directed graph: all
// block names are unique, all connections reference existing blocks, the graph
// contains no cycles, and there is at least one source and one sink.
func (p *Pipeline) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("flow: pipeline ID is empty")
	}
	if len(p.Blocks) == 0 {
		return fmt.Errorf("flow: pipeline %q has no blocks", p.ID)
	}

	names := make(map[string]bool, len(p.Blocks))
	for _, b := range p.Blocks {
		if b.Name == "" {
			return fmt.Errorf("flow: pipeline %q has a block with empty name", p.ID)
		}
		if names[b.Name] {
			return fmt.Errorf("flow: pipeline %q has duplicate block name %q", p.ID, b.Name)
		}
		names[b.Name] = true
	}

	for _, c := range p.Connections {
		if !names[c.From] {
			return fmt.Errorf("flow: pipeline %q connection references unknown block %q", p.ID, c.From)
		}
		if !names[c.To] {
			return fmt.Errorf("flow: pipeline %q connection references unknown block %q", p.ID, c.To)
		}
		if c.From == c.To {
			return fmt.Errorf("flow: pipeline %q has self-loop on block %q", p.ID, c.From)
		}
	}

	// Topological sort to detect cycles and verify there is at least one
	// source and one sink.
	inDeg := make(map[string]int, len(p.Blocks))
	outDeg := make(map[string]int, len(p.Blocks))
	adj := make(map[string][]string, len(p.Blocks))
	for _, b := range p.Blocks {
		inDeg[b.Name] = 0
		outDeg[b.Name] = 0
	}
	for _, c := range p.Connections {
		outDeg[c.From]++
		inDeg[c.To]++
		adj[c.From] = append(adj[c.From], c.To)
	}

	var queue []string
	for _, b := range p.Blocks {
		if inDeg[b.Name] == 0 {
			queue = append(queue, b.Name)
		}
	}

	visited := 0
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		visited++
		for _, nb := range adj[n] {
			inDeg[nb]--
			if inDeg[nb] == 0 {
				queue = append(queue, nb)
			}
		}
	}

	if visited != len(p.Blocks) {
		return fmt.Errorf("flow: pipeline %q contains a cycle", p.ID)
	}

	hasSource := false
	hasSink := false
	for _, b := range p.Blocks {
		if inDeg[b.Name] == 0 || outDeg[b.Name] == 0 {
			// Re-check using original degrees (topo sort modified inDeg).
		}
	}
	for _, b := range p.Blocks {
		if outDeg[b.Name] == 0 {
			hasSink = true
		}
		if inDeg[b.Name] == 0 {
			hasSource = true
		}
	}
	// Re-compute source/sink from original edges since topo sort zeroed inDeg.
	inDeg2 := make(map[string]int, len(p.Blocks))
	outDeg2 := make(map[string]int, len(p.Blocks))
	for _, b := range p.Blocks {
		inDeg2[b.Name] = 0
		outDeg2[b.Name] = 0
	}
	for _, c := range p.Connections {
		outDeg2[c.From]++
		inDeg2[c.To]++
	}
	hasSource = false
	hasSink = false
	for _, b := range p.Blocks {
		if inDeg2[b.Name] == 0 {
			hasSource = true
		}
		if outDeg2[b.Name] == 0 {
			hasSink = true
		}
	}
	if !hasSource {
		return fmt.Errorf("flow: pipeline %q has no source block (every block has incoming edges)", p.ID)
	}
	if !hasSink {
		return fmt.Errorf("flow: pipeline %q has no sink block (every block has outgoing edges)", p.ID)
	}

	return nil
}

// MarshalJSON serialises the Pipeline. It runs validation before encoding.
func (p *Pipeline) MarshalJSON() ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	type Alias Pipeline
	return json.Marshal((*Alias)(p))
}
