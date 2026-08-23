package flow

import (
	"encoding/json"
	"fmt"
)

// UserBlock is one realm-level composite definition (#85).
type UserBlock struct {
	Name         string
	BlockType    string // producer | consumer | producer_consumer; stored, not enforced here
	Source       []byte // a pipeline definition body: {"blocks":[…],"connections":[…]}
	ConfigSchema json.RawMessage
}

// ExpandComposites inlines every user-defined block found in def (issue #85):
// a node whose block_type resolves through resolve is replaced by its stored
// sub-chain, spliced in place — incoming edges attach to the sub-chain's
// first node and outgoing edges leave from its last node (stable topological
// order). Nested composites are expanded recursively and their block names
// nest as "outer.inner.block". resolve returns (nil, nil) for built-in types,
// which pass through unchanged. Definitions are handled as raw node/connection
// slices; Pipeline.Validate is never run on inner bodies.
func ExpandComposites(def []byte, resolve func(name string) (*UserBlock, error)) ([]byte, error) {
	if resolve == nil {
		return nil, fmt.Errorf("flow: nil composite resolver")
	}
	if len(def) == 0 {
		return nil, fmt.Errorf("flow: empty pipeline definition")
	}
	var root rawDefinition
	if err := json.Unmarshal(def, &root); err != nil {
		return nil, fmt.Errorf("flow: pipeline definition does not parse: %w", err)
	}
	blocks, conns, err := expandComposites(resolve, root.Blocks, root.Connections, map[string]bool{}, "")
	if err != nil {
		return nil, err
	}
	return json.Marshal(rawDefinition{Blocks: blocks, Connections: conns})
}

// rawDefinition mirrors Pipeline's JSON shape without ID/Name, so mid-chain
// composite bodies (no global source/sink) can be decoded and re-encoded
// without triggering Validate via MarshalJSON.
type rawDefinition struct {
	Blocks      []PipelineNode `json:"blocks"`
	Connections []Connection   `json:"connections"`
}

// maxCompositeDepth caps composite nesting depth (#85).
const maxCompositeDepth = 32

// expandComposites is the recursive core: it returns nodes/connections with
// every resolvable composite under prefix spliced in place.
func expandComposites(resolve func(name string) (*UserBlock, error), nodes []PipelineNode, conns []Connection, stack map[string]bool, prefix string) ([]PipelineNode, []Connection, error) {
	outConns := append([]Connection(nil), conns...)
	var outNodes []PipelineNode
	var splicedConns []Connection
	for _, node := range nodes {
		ub, err := resolve(node.BlockType)
		if err != nil {
			return nil, nil, err
		}
		if ub == nil {
			outNodes = append(outNodes, node)
			continue
		}
		if stack[ub.Name] {
			return nil, nil, fmt.Errorf("flow: composite cycle detected at %q", ub.Name)
		}
		if len(stack) >= maxCompositeDepth {
			return nil, nil, fmt.Errorf("flow: composite nesting too deep at %q", ub.Name)
		}
		src, err := SubstituteConfig(ub.Source, node.Config)
		if err != nil {
			return nil, nil, err
		}
		var inner rawDefinition
		if err := json.Unmarshal(src, &inner); err != nil {
			return nil, nil, fmt.Errorf("flow: composite %q definition does not parse: %w", ub.Name, err)
		}
		if len(inner.Blocks) == 0 {
			return nil, nil, fmt.Errorf("flow: composite %q has an empty definition", ub.Name)
		}
		if err := checkInnerGraph(ub.Name, inner.Blocks, inner.Connections); err != nil {
			return nil, nil, err
		}
		if _, err := (&Pipeline{Blocks: inner.Blocks, Connections: inner.Connections}).topoOrder(); err != nil {
			return nil, nil, fmt.Errorf("flow: composite %q contains a cycle", ub.Name)
		}

		subStack := make(map[string]bool, len(stack)+1)
		for k := range stack {
			subStack[k] = true
		}
		subStack[ub.Name] = true
		subPrefix := prefix + node.Name + "."
		// Relative names come up: this level's rename loops below are the
		// only place prefixes are applied.
		expNodes, expConns, err := expandComposites(resolve, inner.Blocks, inner.Connections, subStack, "")
		if err != nil {
			return nil, nil, err
		}
		for i := range expNodes {
			expNodes[i].Name = subPrefix + expNodes[i].Name
		}
		for i := range expConns {
			expConns[i].From = subPrefix + expConns[i].From
			expConns[i].To = subPrefix + expConns[i].To
		}
		// Boundary endpoints must be the expanded chain's real first/last
		// nodes: nested composites may have replaced the immediate ones.
		expOrder, err := (&Pipeline{Blocks: expNodes, Connections: expConns}).topoOrder()
		if err != nil {
			return nil, nil, fmt.Errorf("flow: composite %q contains a cycle", ub.Name)
		}
		first, last := expOrder[0], expOrder[len(expOrder)-1]
		for i := range outConns {
			if outConns[i].To == node.Name {
				outConns[i].To = first
			}
			if outConns[i].From == node.Name {
				outConns[i].From = last
			}
		}
		outNodes = append(outNodes, expNodes...)
		splicedConns = append(splicedConns, expConns...)
	}
	return outNodes, append(outConns, splicedConns...), nil
}

// checkInnerGraph performs the structural sanity checks allowed for composite
// bodies (full Validate demands a global source/sink they need not have):
// unique non-empty names and connection endpoints that exist. Cycles are
// rejected separately via topoOrder above.
func checkInnerGraph(composite string, blocks []PipelineNode, conns []Connection) error {
	names := make(map[string]bool, len(blocks))
	for _, b := range blocks {
		if b.Name == "" {
			return fmt.Errorf("flow: composite %q has a block with empty name", composite)
		}
		if names[b.Name] {
			return fmt.Errorf("flow: composite %q has duplicate block name %q", composite, b.Name)
		}
		names[b.Name] = true
	}
	for _, c := range conns {
		if !names[c.From] {
			return fmt.Errorf("flow: composite %q connection references unknown block %q", composite, c.From)
		}
		if !names[c.To] {
			return fmt.Errorf("flow: composite %q connection references unknown block %q", composite, c.To)
		}
	}
	return nil
}
