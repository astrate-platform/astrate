package flow

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"github.com/astrate-platform/astrate/internal/engine/stream"
)

// ErrUnknownBlockType is returned when Instantiate sees a block_type with no
// registered constructor.
var ErrUnknownBlockType = errors.New("flow: unknown block type")

// Deps holds process-level dependencies block constructors may need.
// Zero fields are valid when a pipeline uses only blocks that do not need them.
type Deps struct {
	// Bus is the live event bus (required by astarte_source).
	Bus *stream.Bus
	// Realm is the tenant the pipeline runs for; sources may default to it
	// when their config omits realm.
	Realm string
	// FlowName is the durable/named flow instance name (optional). Used by
	// container labels and similar operator metadata.
	FlowName string
}

// Constructor builds one Block from a pipeline node. name is the pipeline
// node name and should be returned by Block.Name() for metrics and logging.
type Constructor func(name string, config map[string]any, deps Deps) (Block, error)

// Registry maps pipeline block_type strings to constructors.
type Registry struct {
	ctors map[string]Constructor
}

// NewRegistry returns an empty block-type registry.
func NewRegistry() *Registry {
	return &Registry{ctors: make(map[string]Constructor)}
}

// Register associates blockType with ctor. Later Register calls for the same
// type replace the previous constructor.
func (r *Registry) Register(blockType string, ctor Constructor) {
	if r.ctors == nil {
		r.ctors = make(map[string]Constructor)
	}
	r.ctors[blockType] = ctor
}

// Has reports whether blockType has a registered constructor.
func (r *Registry) Has(blockType string) bool {
	_, ok := r.ctors[blockType]
	return ok
}

// Types returns the registered block_type names in sorted order.
func (r *Registry) Types() []string {
	out := make([]string, 0, len(r.ctors))
	for t := range r.ctors {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// Instantiate turns a validated Pipeline description into an ordered block
// list suitable for Manager.StartFlow. Blocks are returned in topological
// order (sources first, sinks last). Linear graphs are the supported
// production shape; DAGs are flattened in topo order (fan-in/fan-out is not
// yet modelled in BlockGraph).
func (r *Registry) Instantiate(p *Pipeline, deps Deps) ([]Block, error) {
	if p == nil {
		return nil, fmt.Errorf("flow: pipeline is nil")
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	order, err := p.topoOrder()
	if err != nil {
		return nil, err
	}
	byName := make(map[string]PipelineNode, len(p.Blocks))
	for _, n := range p.Blocks {
		byName[n.Name] = n
	}
	out := make([]Block, 0, len(order))
	for _, name := range order {
		node := byName[name]
		if node.BlockType == "" {
			return nil, fmt.Errorf("flow: pipeline %q block %q has empty block_type", p.ID, name)
		}
		ctor, ok := r.ctors[node.BlockType]
		if !ok {
			return nil, fmt.Errorf("%w: %q (block %q)", ErrUnknownBlockType, node.BlockType, name)
		}
		b, err := ctor(name, node.Config, deps)
		if err != nil {
			stopAll(out)
			return nil, fmt.Errorf("flow: construct %q (%s): %w", name, node.BlockType, err)
		}
		if b == nil {
			stopAll(out)
			return nil, fmt.Errorf("flow: construct %q (%s) returned nil block", name, node.BlockType)
		}
		out = append(out, b)
	}
	return out, nil
}

// stopAll calls Stop on every block in blocks that implements Stopper, so a
// partially-built pipeline does not leak resources already acquired by
// earlier blocks (e.g. a container block's running Docker container).
func stopAll(blocks []Block) {
	for _, b := range blocks {
		if s, ok := b.(Stopper); ok {
			s.Stop()
		}
	}
}

// topoOrder returns block names in topological order (Kahn). Validate must
// already have accepted the graph (acyclic, non-empty).
func (p *Pipeline) topoOrder() ([]string, error) {
	inDeg := make(map[string]int, len(p.Blocks))
	adj := make(map[string][]string, len(p.Blocks))
	for _, b := range p.Blocks {
		inDeg[b.Name] = 0
	}
	for _, c := range p.Connections {
		inDeg[c.To]++
		adj[c.From] = append(adj[c.From], c.To)
	}
	var queue []string
	for _, b := range p.Blocks {
		if inDeg[b.Name] == 0 {
			queue = append(queue, b.Name)
		}
	}
	// Stable: process sources in declaration order among zero in-degree nodes.
	sort.Slice(queue, func(i, j int) bool {
		return blockIndex(p, queue[i]) < blockIndex(p, queue[j])
	})

	order := make([]string, 0, len(p.Blocks))
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		order = append(order, n)
		// Keep queue stable when multiple successors become ready.
		var ready []string
		for _, nb := range adj[n] {
			inDeg[nb]--
			if inDeg[nb] == 0 {
				ready = append(ready, nb)
			}
		}
		sort.Slice(ready, func(i, j int) bool {
			return blockIndex(p, ready[i]) < blockIndex(p, ready[j])
		})
		queue = append(queue, ready...)
	}
	if len(order) != len(p.Blocks) {
		return nil, fmt.Errorf("flow: pipeline %q contains a cycle", p.ID)
	}
	return order, nil
}

func blockIndex(p *Pipeline, name string) int {
	for i, b := range p.Blocks {
		if b.Name == name {
			return i
		}
	}
	return len(p.Blocks)
}

// ParseDefinition unmarshals a stored pipeline definition (blocks +
// connections JSON) into a Pipeline, setting ID and Name from the caller
// (the store keeps those outside the definition blob).
func ParseDefinition(id, name string, definition []byte) (*Pipeline, error) {
	if len(definition) == 0 {
		return nil, fmt.Errorf("flow: empty pipeline definition")
	}
	p := &Pipeline{ID: id, Name: name}
	if err := json.Unmarshal(definition, p); err != nil {
		return nil, fmt.Errorf("flow: pipeline definition does not parse: %w", err)
	}
	// Unmarshal may overwrite id/name if present in JSON; prefer caller values
	// when non-empty so store name wins over a stale embedded field.
	if id != "" {
		p.ID = id
	}
	if name != "" {
		p.Name = name
	}
	if p.ID == "" {
		p.ID = p.Name
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

// InstanceID builds the Manager map key for a named flow instance
// (realm + "/" + flowName). Different flows may share one pipeline name.
func InstanceID(realm, flowName string) string {
	return realm + "/" + flowName
}

// PipelineID is a legacy alias for InstanceID. Prefer InstanceID;
// the map key is the flow instance name, not the pipeline name.
func PipelineID(realm, name string) string {
	return InstanceID(realm, name)
}
