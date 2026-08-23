package blocks

import "encoding/json"

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
	// ConfigSchema is a Draft-07-style JSON Schema (object form) describing the
	// block's config keys, served through the blocks-discovery API so clients can
	// render forms without hardcoding block knowledge.
	ConfigSchema json.RawMessage `json:"config_schema,omitempty"`
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
	TypeToJSON: {
		Type:    TypeToJSON,
		Role:    RoleTransform,
		Summary: "Encode the payload as JSON bytes (binary, application/json)",
		Config:  "(none)",
	},
	TypeUpdateMetadata: {
		Type:    TypeUpdateMetadata,
		Role:    RoleTransform,
		Summary: "Merge set_metadata into Metadata, then delete delete_metadata keys",
		Config:  "set_metadata (object string→string), delete_metadata (string array); at least one required",
	},
	TypeSplitMap: {
		Type:    TypeSplitMap,
		Role:    RoleTransform,
		Summary: "Explode a map payload into one message per field",
		Config:  `key_template (default "{key}/{field}"; placeholders {key}, {field})`,
	},
	TypeRandomSource: {
		Type:    TypeRandomSource,
		Role:    RoleSource,
		Summary: "Emit a random integer, real, or boolean value every interval",
		Config:  `type (required: integer|real|boolean), interval_ms (default 1000), min, max, key (default "random")`,
	},
	TypeSort: {
		Type:    TypeSort,
		Role:    RoleTransform,
		Summary: "Buffer messages and release them in ascending timestamp order behind a window",
		Config:  "window_ms (default 1000, must be ≥ 0), dedup (bool, default false)",
	},
	TypeJSONPathMap: {
		Type:    TypeJSONPathMap,
		Role:    RoleTransform,
		Summary: "Reshape JSON messages through a template with $MESSAGE / $METADATA placeholders",
		Config:  "template (required; placeholders $MESSAGE, $MESSAGE.<path>, $METADATA.<name>; unresolvable paths fail the message)",
	},
	TypeHTTPSource: {
		Type:    TypeHTTPSource,
		Role:    RoleSource,
		Summary: "Poll GET URLs round-robin and emit each response body as a binary message",
		Config:  "urls (required string array, at least one), interval_ms (default 1000), timeout_ms (default 5000)",
	},
	TypeHTTPSink: {
		Type:    TypeHTTPSink,
		Role:    RoleSink,
		Summary: "POST each message payload to a URL (binary as-is, strings as text/plain, others JSON)",
		Config:  `url (required), method (default "POST"), timeout_ms (default 5000), headers (object string→string)`,
	},
	TypeMQTTSource: {
		Type:    TypeMQTTSource,
		Role:    RoleSource,
		Summary: "Subscribe to MQTT topics and emit each delivery as a binary message",
		Config:  `url (required broker URL, e.g. tcp://127.0.0.1:1883), topics (required string array, at least one), qos (0|1|2, default 0), client_id, username, password`,
	},
	TypeMQTTSink: {
		Type:    TypeMQTTSink,
		Role:    RoleSink,
		Summary: "Publish each message payload to an MQTT topic (binary/strings raw, others JSON)",
		Config:  `url (required), topic (required), qos (0|1|2, default 0), retained (bool, default false), client_id, username, password`,
	},
	TypeContainer: {
		Type:    TypeContainer,
		Role:    RoleTransform,
		Summary: "Run a Docker image as a transform (HTTP POST /v1/message); local Docker only (PoC)",
		Config:  "image (required), config (object→ASTRATE_FLOW_CONFIG), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",
	},
	TypeVirtualDevicePool: {
		Type:    TypeVirtualDevicePool,
		Role:    RoleSink,
		Summary: "Publish each message as a registered virtual device; key is <device_id>/<interface></path>",
		Config:  "devices (required string array of registered device_ids)",
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
	if ok {
		if s := builtinSchemas[blockType]; s != "" {
			info.ConfigSchema = json.RawMessage(s)
		}
	}
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
