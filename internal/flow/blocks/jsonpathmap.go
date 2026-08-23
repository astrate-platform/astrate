package blocks

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"maps"
	"strings"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
)

// JSONPathMap reshapes JSON messages through a template: placeholders in the
// template are substituted with the original payload or metadata values and
// the result is parsed as a JSON object payload.
//
// Config keys:
//   - template (string, required): a JSON document containing placeholders.
//
// Placeholder grammar:
//   - $MESSAGE — the entire original payload, rendered as raw JSON
//   - $MESSAGE.<seg>(.<seg>)* — dot-separated path into the payload; each seg
//     is one or more of [A-Za-z0-9_-]. Unresolvable paths fail Process.
//   - $METADATA.<name> — metadata lookup; missing names fail Process.
//     $MESSAGE followed by anything not matching .<seg> is the whole-payload
//     token ($MESSAGES renders $MESSAGE then a literal "S").
//
// Whole-message and path tokens insert raw (unquoted) JSON; time.Time becomes
// an RFC3339Nano string and []byte a base64 string before marshaling.
// Metadata tokens substitute the value JSON-escaped, so a template occurrence
// like "$METADATA.user" lands in the rendered document as the JSON string
// "alice". The output message is
// TypeMap with Key/Timestamp inherited, Metadata cloned, and per-field types
// inferred: float64→real, bool→boolean, string→string; nested objects, arrays
// and null become strings carrying their compact JSON rendering. The input
// message is never mutated.
func JSONPathMap(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
	template := stringConfig(config, "template", "")
	if template == "" {
		return nil, fmt.Errorf("json_path_map: template is required")
	}
	return flow.NewTransformBlock(name, func(msg *flow.Message) ([]*flow.Message, error) {
		if msg == nil {
			return nil, nil
		}
		rendered, err := expandJSONTemplate(template, msg)
		if err != nil {
			return nil, err
		}
		var parsed map[string]any
		if err := json.Unmarshal([]byte(rendered), &parsed); err != nil {
			return nil, fmt.Errorf("json_path_map: template did not produce a JSON object: %v", err)
		}
		fieldTypes := make(map[string]flow.DataType, len(parsed))
		for k, v := range parsed {
			switch v.(type) {
			case float64:
				fieldTypes[k] = flow.TypeReal
			case bool:
				fieldTypes[k] = flow.TypeBoolean
			case string:
				fieldTypes[k] = flow.TypeString
			default:
				raw, err := json.Marshal(v)
				if err != nil {
					return nil, fmt.Errorf("json_path_map: %v", err)
				}
				parsed[k] = string(raw)
				fieldTypes[k] = flow.TypeString
			}
		}
		out := &flow.Message{
			Key:        msg.Key,
			Timestamp:  msg.Timestamp,
			Type:       flow.TypeMap,
			Data:       parsed,
			FieldTypes: fieldTypes,
		}
		if msg.Metadata != nil {
			out.Metadata = maps.Clone(msg.Metadata)
		}
		return []*flow.Message{out}, nil
	}), nil
}

const (
	msgPlaceholder      = "$MESSAGE"
	metadataPlaceholder = "$METADATA"
)

// expandJSONTemplate substitutes every placeholder in tmpl against msg,
// scanning left-to-right with longest match.
func expandJSONTemplate(tmpl string, msg *flow.Message) (string, error) {
	var b strings.Builder
	i := 0
	for i < len(tmpl) {
		if tmpl[i] != '$' {
			b.WriteByte(tmpl[i])
			i++
			continue
		}
		switch {
		case strings.HasPrefix(tmpl[i:], msgPlaceholder):
			rest := tmpl[i+len(msgPlaceholder):]
			if segs, n, ok := matchDottedPath(rest); ok {
				v, err := resolveMessagePath(msg, segs)
				if err != nil {
					return "", err
				}
				raw, err := marshalRawJSON(v)
				if err != nil {
					return "", fmt.Errorf("json_path_map: %v", err)
				}
				b.Write(raw)
				i += len(msgPlaceholder) + n
				continue
			}
			raw, err := toJSONBytes(msg)
			if err != nil {
				return "", fmt.Errorf("json_path_map: %v", err)
			}
			b.Write(raw)
			i += len(msgPlaceholder)
		case strings.HasPrefix(tmpl[i:], metadataPlaceholder):
			rest := tmpl[i+len(metadataPlaceholder):]
			name, n, ok := matchDotSeg(rest)
			if !ok {
				b.WriteByte('$')
				i++
				continue
			}
			v, found := msg.Metadata[name]
			if !found {
				return "", fmt.Errorf("json_path_map: cannot resolve $METADATA.%s", name)
			}
			// Substitute the JSON-escaped body; the template's surrounding
			// quotes make it a JSON string (e.g. "$METADATA.user" → "alice").
			q, err := json.Marshal(v)
			if err != nil {
				return "", fmt.Errorf("json_path_map: %v", err)
			}
			b.Write(q[1 : len(q)-1])
			i += len(metadataPlaceholder) + n
		default:
			b.WriteByte('$')
			i++
		}
	}
	return b.String(), nil
}

// isSegChar reports whether c belongs to a placeholder segment
// ([A-Za-z0-9_-]).
func isSegChar(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '-'
}

// matchDottedPath parses ".seg(.seg)*" at the start of s, greedily consuming
// as many segments as possible. It returns the segments and the number of
// bytes consumed, or ok=false when s does not start with a dot followed by at
// least one segment character.
func matchDottedPath(s string) ([]string, int, bool) {
	if !strings.HasPrefix(s, ".") {
		return nil, 0, false
	}
	var segs []string
	i := 1
	for {
		j := i
		for j < len(s) && isSegChar(s[j]) {
			j++
		}
		if j == i {
			break
		}
		segs = append(segs, s[i:j])
		i = j
		if i >= len(s) || s[i] != '.' || i+1 >= len(s) || !isSegChar(s[i+1]) {
			break
		}
		i++
	}
	if len(segs) == 0 {
		return nil, 0, false
	}
	return segs, i, true
}

// matchDotSeg parses a single ".seg" at the start of s, returning the segment
// and bytes consumed.
func matchDotSeg(s string) (string, int, bool) {
	if !strings.HasPrefix(s, ".") {
		return "", 0, false
	}
	j := 1
	for j < len(s) && isSegChar(s[j]) {
		j++
	}
	if j == 1 {
		return "", 0, false
	}
	return s[1:j], j, true
}

// resolveMessagePath walks segs through map levels starting at msg.Data; the
// first segment names a top-level field. A non-map level or a missing key is
// an error naming the full path.
func resolveMessagePath(msg *flow.Message, segs []string) (any, error) {
	full := strings.Join(segs, ".")
	cur := msg.Data
	for _, seg := range segs {
		level, ok := cur.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("json_path_map: cannot resolve $MESSAGE.%s", full)
		}
		v, ok := level[seg]
		if !ok {
			return nil, fmt.Errorf("json_path_map: cannot resolve $MESSAGE.%s", full)
		}
		cur = v
	}
	return cur, nil
}

// marshalRawJSON renders v as unquoted JSON, converting time.Time to its
// RFC3339Nano string and []byte to its base64 string first (same convention
// as tojson.go).
func marshalRawJSON(v any) ([]byte, error) {
	switch t := v.(type) {
	case time.Time:
		return json.Marshal(t.Format(time.RFC3339Nano))
	case []byte:
		return json.Marshal(base64.StdEncoding.EncodeToString(t))
	default:
		return json.Marshal(v)
	}
}
