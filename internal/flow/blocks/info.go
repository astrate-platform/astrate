package blocks

// Role is the operator-facing position of a block in a pipeline.
type Role string

// The three Role values a block can occupy in a pipeline.
const (
	RoleSource    Role = "source"
	RoleTransform Role = "transform"
	RoleSink      Role = "sink"
)

// Info is static operator documentation for a built-in block type.
type Info struct {
	// Type is the block_type string stored in pipeline definitions.
	Type string `json:"type"`
	// Role is source | transform | sink.
	Role Role `json:"role"`
	// Summary is a one-line description for operators.
	Summary string `json:"summary"`
	// Config notes key config fields (not a full JSON Schema).
	Config string `json:"config,omitempty"`
}

// builtinInfo is the discovery catalog for DefaultRegistry types.
// Keep in sync when registering new built-ins in catalog.go / transform.go.
var builtinInfo = map[string]Info{
	TypeAstarteSource: {
		Type:    TypeAstarteSource,
		Role:    RoleSource,
		Summary: "Subscribe to the live stream bus and emit FlowMessages",
		Config:  "realm (string, default deps.Realm), interface (string), path (string prefix)",
	},
	TypeFilter: {
		Type:    TypeFilter,
		Role:    RoleTransform,
		Summary: "Drop messages that fail AND of configured conditions",
		Config:  "key_prefix, key_contains, type (integer|real|boolean|datetime|binary|string|map), metadata (object string→string); at least one required",
	},
	TypeMap: {
		Type:    TypeMap,
		Role:    RoleTransform,
		Summary: "Rewrite key and/or metadata; payload unchanged",
		Config:  "key (template with {key} and {metadata.<name>}), set_metadata (object), delete_metadata (string array)",
	},
	TypeContainer: {
		Type:    TypeContainer,
		Role:    RoleTransform,
		Summary: "Run a Docker image as a transform (HTTP POST /v1/message); local Docker only (PoC)",
		Config:  "image (required), config (object→ASTRATE_FLOW_CONFIG), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",
	},
	TypeNullSink: {
		Type:    TypeNullSink,
		Role:    RoleSink,
		Summary: "Discard every message",
		Config:  "(none)",
	},
	TypeLogSink: {
		Type:    TypeLogSink,
		Role:    RoleSink,
		Summary: "Log each message at debug level",
		Config:  "(none)",
	},
}

// LookupInfo returns static operator docs for a known built-in type.
// Unknown types return false (registry may still have a custom constructor).
func LookupInfo(blockType string) (Info, bool) {
	info, ok := builtinInfo[blockType]
	return info, ok
}

// InfoForTypes returns Info for each type, filling a minimal stub when the
// type is registered but has no static docs (custom constructors).
func InfoForTypes(types []string) []Info {
	out := make([]Info, 0, len(types))
	for _, t := range types {
		if info, ok := LookupInfo(t); ok {
			out = append(out, info)
			continue
		}
		out = append(out, Info{Type: t, Role: RoleTransform, Summary: "registered block (no built-in docs)"})
	}
	return out
}
