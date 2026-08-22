package blocks

import (
	"fmt"
	"maps"
	"strings"

	"github.com/astrate-platform/astrate/internal/flow"
)

// Well-known transform block_type strings.
const (
	TypeFilter = "filter"
	TypeMap    = "map"
)

// Filter keeps or drops messages. All non-empty conditions must match (AND).
// Config keys (all optional; at least one required at construct time):
//   - key_prefix (string): message Key must have this prefix
//   - key_contains (string): message Key must contain this substring
//   - type (string): wire type name — integer|real|boolean|datetime|binary|string|map
//   - metadata (object string→string): every key must equal the given value
//
// Matching messages pass through unchanged; non-matching messages are dropped
// (zero outputs, no error).
func Filter(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
	cfg, err := parseFilterConfig(config)
	if err != nil {
		return nil, fmt.Errorf("filter: %w", err)
	}
	return flow.NewTransformBlock(name, func(msg *flow.Message) ([]*flow.Message, error) {
		if msg == nil || !cfg.match(msg) {
			return nil, nil
		}
		return []*flow.Message{msg}, nil
	}), nil
}

type filterConfig struct {
	keyPrefix   string
	keyContains string
	dataType    *flow.DataType // nil = no type constraint
	metadata    map[string]string
}

func parseFilterConfig(config map[string]any) (filterConfig, error) {
	cfg := filterConfig{
		keyPrefix:   stringConfig(config, "key_prefix", ""),
		keyContains: stringConfig(config, "key_contains", ""),
		metadata:    stringMapConfig(config, "metadata"),
	}
	if t := stringConfig(config, "type", ""); t != "" {
		dt, err := parseTypeName(t)
		if err != nil {
			return cfg, err
		}
		cfg.dataType = &dt
	}
	if cfg.keyPrefix == "" && cfg.keyContains == "" && cfg.dataType == nil && len(cfg.metadata) == 0 {
		return cfg, fmt.Errorf("at least one of key_prefix, key_contains, type, metadata is required")
	}
	return cfg, nil
}

func (c filterConfig) match(msg *flow.Message) bool {
	if c.keyPrefix != "" && !strings.HasPrefix(msg.Key, c.keyPrefix) {
		return false
	}
	if c.keyContains != "" && !strings.Contains(msg.Key, c.keyContains) {
		return false
	}
	if c.dataType != nil && msg.Type != *c.dataType {
		return false
	}
	for k, want := range c.metadata {
		if msg.Metadata == nil || msg.Metadata[k] != want {
			return false
		}
	}
	return true
}

// Map rewrites key and/or metadata. The payload (Type/Data) is never changed.
// Config keys (all optional; at least one required):
//   - key (string): new stream key; supports placeholders {key} and
//     {metadata.<name>} (missing metadata expands to empty)
//   - set_metadata (object string→string): merge into Metadata (overwrites)
//   - delete_metadata ([]string): remove these Metadata keys after set
//
// The message is shallow-copied so concurrent lanes never share Metadata maps.
func Map(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
	cfg, err := parseMapConfig(config)
	if err != nil {
		return nil, fmt.Errorf("map: %w", err)
	}
	return flow.NewTransformBlock(name, func(msg *flow.Message) ([]*flow.Message, error) {
		if msg == nil {
			return nil, nil
		}
		out := cloneMessage(msg)
		if cfg.keyTemplate != "" {
			out.Key = expandTemplate(cfg.keyTemplate, msg)
		}
		if len(cfg.setMetadata) > 0 {
			if out.Metadata == nil {
				out.Metadata = make(map[string]string, len(cfg.setMetadata))
			}
			for k, v := range cfg.setMetadata {
				out.Metadata[k] = v
			}
		}
		for _, k := range cfg.deleteMetadata {
			delete(out.Metadata, k)
		}
		return []*flow.Message{out}, nil
	}), nil
}

type mapConfig struct {
	keyTemplate    string
	setMetadata    map[string]string
	deleteMetadata []string
}

func parseMapConfig(config map[string]any) (mapConfig, error) {
	cfg := mapConfig{
		keyTemplate:    stringConfig(config, "key", ""),
		setMetadata:    stringMapConfig(config, "set_metadata"),
		deleteMetadata: stringSliceConfig(config, "delete_metadata"),
	}
	if cfg.keyTemplate == "" && len(cfg.setMetadata) == 0 && len(cfg.deleteMetadata) == 0 {
		return cfg, fmt.Errorf("at least one of key, set_metadata, delete_metadata is required")
	}
	return cfg, nil
}

func cloneMessage(msg *flow.Message) *flow.Message {
	out := *msg
	if msg.Metadata != nil {
		out.Metadata = maps.Clone(msg.Metadata)
	}
	// FieldTypes / FieldSubtypes / Data are not mutated by Map; share is fine.
	return &out
}

// expandTemplate replaces {key} and {metadata.<name>} in tmpl.
func expandTemplate(tmpl string, msg *flow.Message) string {
	s := strings.ReplaceAll(tmpl, "{key}", msg.Key)
	// Replace {metadata.X} for each metadata key present; then any remaining
	// known-shape placeholders for missing keys become empty.
	const metaPrefix = "{metadata."
	for {
		i := strings.Index(s, metaPrefix)
		if i < 0 {
			break
		}
		rest := s[i+len(metaPrefix):]
		j := strings.IndexByte(rest, '}')
		if j < 0 {
			break
		}
		name := rest[:j]
		val := ""
		if msg.Metadata != nil {
			val = msg.Metadata[name]
		}
		s = s[:i] + val + rest[j+1:]
	}
	return s
}

func parseTypeName(s string) (flow.DataType, error) {
	switch s {
	case "integer":
		return flow.TypeInteger, nil
	case "real":
		return flow.TypeReal, nil
	case "boolean":
		return flow.TypeBoolean, nil
	case "datetime":
		return flow.TypeDatetime, nil
	case "binary":
		return flow.TypeBinary, nil
	case "string":
		return flow.TypeString, nil
	case "map":
		return flow.TypeMap, nil
	default:
		return 0, fmt.Errorf("unknown type %q (want integer|real|boolean|datetime|binary|string|map)", s)
	}
}

func stringMapConfig(config map[string]any, key string) map[string]string {
	if config == nil {
		return nil
	}
	raw, ok := config[key]
	if !ok || raw == nil {
		return nil
	}
	// JSON objects decode as map[string]any; tests may pass map[string]string.
	switch m := raw.(type) {
	case map[string]string:
		if len(m) == 0 {
			return nil
		}
		out := make(map[string]string, len(m))
		for k, v := range m {
			out[k] = v
		}
		return out
	case map[string]any:
		if len(m) == 0 {
			return nil
		}
		out := make(map[string]string, len(m))
		for k, v := range m {
			if s, ok := v.(string); ok {
				out[k] = s
			} else if v != nil {
				out[k] = fmt.Sprint(v)
			}
		}
		if len(out) == 0 {
			return nil
		}
		return out
	default:
		return nil
	}
}

func stringSliceConfig(config map[string]any, key string) []string {
	if config == nil {
		return nil
	}
	raw, ok := config[key]
	if !ok || raw == nil {
		return nil
	}
	switch s := raw.(type) {
	case []string:
		return append([]string(nil), s...)
	case []any:
		out := make([]string, 0, len(s))
		for _, v := range s {
			if str, ok := v.(string); ok {
				out = append(out, str)
			}
		}
		return out
	default:
		return nil
	}
}
