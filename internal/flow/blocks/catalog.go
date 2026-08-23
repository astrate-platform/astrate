// Package blocks registers the built-in Flow block catalog used to instantiate
// stored pipelines (milestone v2.0: block factory).
package blocks

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks/astartesource"
	"github.com/astrate-platform/astrate/internal/flow/blocks/container"
)

// Well-known block_type strings stored in pipeline definitions.
const (
	TypeAstarteSource = "astarte_source"
	TypeNullSink      = "null_sink"
	TypeLogSink       = "log_sink"
	// TypeContainer is the Docker-backed custom processing block (Design B / #43).
	TypeContainer = container.Type
	// TypeFilter and TypeMap are declared in transform.go.

	// Pure transforms and the random source (#80, astarte_flow parity).
	TypeToJSON         = "to_json"
	TypeUpdateMetadata = "update_metadata"
	TypeSplitMap       = "split_map"
	TypeRandomSource   = "random_source"
	TypeSort           = "sort"

	// TypeJSONPathMap reshapes JSON messages through a template (#81,
	// astarte_flow parity); declared in jsonpathmap.go's neighborhood.
	TypeJSONPathMap = "json_path_map"

	// TypeHTTPSource / TypeHTTPSink poll/POST over HTTP (#82, astarte_flow
	// parity); declared in httpblocks.go.
)

// DefaultRegistry returns a registry with the minimum useful built-in set:
// AstarteSource (bus → Message), filter/map transforms, container (custom
// image via local Docker HTTP bridge), and null/log sinks so operators can
// compose a complete source→transform→sink pipeline.
func DefaultRegistry() *flow.Registry {
	r := flow.NewRegistry()
	r.Register(TypeAstarteSource, AstarteSource)
	r.Register(TypeFilter, Filter)
	r.Register(TypeMap, Map)
	r.Register(TypeToJSON, ToJSON)
	r.Register(TypeUpdateMetadata, UpdateMetadata)
	r.Register(TypeSplitMap, SplitMap)
	r.Register(TypeRandomSource, RandomSource)
	r.Register(TypeSort, Sort)
	r.Register(TypeJSONPathMap, JSONPathMap)
	r.Register(TypeHTTPSource, HTTPSource)
	r.Register(TypeHTTPSink, HTTPSink)
	r.Register(TypeContainer, container.Constructor)
	r.Register(TypeNullSink, NullSink)
	r.Register(TypeLogSink, LogSink)
	return r
}

// AstarteSource constructs an astartesource.Source. Config keys:
//   - realm (string): tenant; defaults to deps.Realm
//   - interface (string): optional interface filter
//   - path (string): optional path prefix filter
func AstarteSource(name string, config map[string]any, deps flow.Deps) (flow.Block, error) {
	if deps.Bus == nil {
		return nil, fmt.Errorf("astarte_source requires a stream bus")
	}
	cfg := astartesource.Config{
		Realm:     stringConfig(config, "realm", deps.Realm),
		Interface: stringConfig(config, "interface", ""),
		Path:      stringConfig(config, "path", ""),
	}
	if cfg.Realm == "" {
		return nil, fmt.Errorf("astarte_source: realm is required (config or deps)")
	}
	src := astartesource.New(deps.Bus, cfg)
	return namedSourceStopper(name, src), nil
}

// NullSink discards every message. Useful as a placeholder sink in tests and
// for pipelines that only need side effects from transforms.
func NullSink(name string, _ map[string]any, _ flow.Deps) (flow.Block, error) {
	return flow.NewSinkBlock(name, func(*flow.Message) error { return nil }), nil
}

// LogSink logs each message at debug level via slog.Default().
func LogSink(name string, _ map[string]any, _ flow.Deps) (flow.Block, error) {
	return flow.NewSinkBlock(name, func(msg *flow.Message) error {
		if msg == nil {
			return nil
		}
		slog.Default().Debug("flow log_sink",
			"block", name,
			"key", msg.Key,
			"type", msg.Type,
			"data", msg.Data,
		)
		return nil
	}), nil
}

func stringConfig(config map[string]any, key, def string) string {
	if config == nil {
		return def
	}
	v, ok := config[key]
	if !ok || v == nil {
		return def
	}
	s, ok := v.(string)
	if !ok {
		return def
	}
	return s
}

// namedSS wraps a Source+Stopper so Block.Name returns the pipeline node name
// while preserving Emit and Stop for the manager pump/teardown.
type namedSS struct {
	inner flow.Block
	name  string
}

func namedSourceStopper(name string, inner flow.Block) flow.Block {
	return &namedSS{inner: inner, name: name}
}

func (n *namedSS) Process(msg *flow.Message) ([]*flow.Message, error) {
	return n.inner.Process(msg)
}

func (n *namedSS) Name() string { return n.name }

func (n *namedSS) Emit(ctx context.Context) ([]*flow.Message, error) {
	src, ok := n.inner.(flow.Source)
	if !ok {
		return nil, fmt.Errorf("flow: block %q is not a Source", n.name)
	}
	return src.Emit(ctx)
}

func (n *namedSS) Stop() {
	if s, ok := n.inner.(flow.Stopper); ok {
		s.Stop()
	}
}

// Ensure namedSS implements the interfaces the pump and StopFlow use.
var (
	_ flow.Block   = (*namedSS)(nil)
	_ flow.Source  = (*namedSS)(nil)
	_ flow.Stopper = (*namedSS)(nil)
)
