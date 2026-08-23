package blocks

import (
	"fmt"
	"maps"
	"strings"

	"github.com/astrate-platform/astrate/internal/flow"
)

// SplitMap explodes a TypeMap payload into one message per field. Non-map
// messages are dropped silently (zero outputs, no error).
//
// Config keys (all optional):
//   - key_template (string): output Key; placeholders {key} (original key) and
//     {field} (field name); default "{key}/{field}"
//
// Each output message carries a single field value as its payload: the type
// comes from FieldTypes when present, else is inferred from the Go value
// (bool→boolean, int64→integer, float64→real, anything else → string via
// fmt.Sprint). FieldSubtypes carry over per field; Metadata is cloned;
// Timestamp is inherited.
func SplitMap(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
	keyTemplate := stringConfig(config, "key_template", "{key}/{field}")
	return flow.NewTransformBlock(name, func(msg *flow.Message) ([]*flow.Message, error) {
		if msg == nil || msg.Type != flow.TypeMap {
			return nil, nil
		}
		fields, ok := msg.Data.(map[string]any)
		if !ok {
			return nil, nil
		}
		out := make([]*flow.Message, 0, len(fields))
		for k, v := range fields {
			m := &flow.Message{
				Key:       expandSplitKey(keyTemplate, msg.Key, k),
				Timestamp: msg.Timestamp,
			}
			if msg.Metadata != nil {
				m.Metadata = maps.Clone(msg.Metadata)
			}
			if dt, ok := msg.FieldTypes[k]; ok {
				m.Type = dt
				m.Data = v
			} else {
				switch val := v.(type) {
				case bool:
					m.Type, m.Data = flow.TypeBoolean, val
				case int64:
					m.Type, m.Data = flow.TypeInteger, val
				case float64:
					m.Type, m.Data = flow.TypeReal, val
				default:
					m.Type, m.Data = flow.TypeString, fmt.Sprint(v)
				}
			}
			m.Subtype = msg.FieldSubtypes[k]
			out = append(out, m)
		}
		return out, nil
	}), nil
}

// expandSplitKey replaces {key} and {field} in tmpl.
func expandSplitKey(tmpl, key, field string) string {
	s := strings.ReplaceAll(tmpl, "{key}", key)
	return strings.ReplaceAll(s, "{field}", field)
}
