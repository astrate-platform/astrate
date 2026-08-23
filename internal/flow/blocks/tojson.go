package blocks

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
)

// ToJSON converts a message payload into its JSON representation, emitted as
// a TypeBinary message with subtype "application/json".
//
// Config keys: none.
//
// Map payloads are converted field-by-field (TypeBinary fields become base64
// strings, TypeDatetime fields become RFC3339Nano strings); scalar datetime
// and binary payloads get the same treatment. The input message is never
// mutated.
func ToJSON(name string, _ map[string]any, _ flow.Deps) (flow.Block, error) {
	return flow.NewTransformBlock(name, func(msg *flow.Message) ([]*flow.Message, error) {
		if msg == nil {
			return nil, nil
		}
		data, err := toJSONBytes(msg)
		if err != nil {
			return nil, fmt.Errorf("to_json: %w", err)
		}
		out := cloneMessage(msg)
		out.Subtype = "application/json"
		out.Type = flow.TypeBinary
		out.Data = data
		out.FieldTypes = nil
		out.FieldSubtypes = nil
		return []*flow.Message{out}, nil
	}), nil
}

// toJSONBytes renders msg.Data as JSON bytes. Map fields typed as binary or
// datetime are rendered as their wire string forms; other values pass
// through to encoding/json unchanged.
func toJSONBytes(msg *flow.Message) ([]byte, error) {
	switch msg.Type {
	case flow.TypeMap:
		raw, ok := msg.Data.(map[string]any)
		if !ok {
			break
		}
		out := make(map[string]any, len(raw))
		for k, v := range raw {
			switch msg.FieldTypes[k] {
			case flow.TypeBinary:
				if bs, ok := v.([]byte); ok {
					out[k] = base64.StdEncoding.EncodeToString(bs)
					continue
				}
			case flow.TypeDatetime:
				if t, ok := v.(time.Time); ok {
					out[k] = t.Format(time.RFC3339Nano)
					continue
				}
			}
			out[k] = v
		}
		return json.Marshal(out)
	case flow.TypeDatetime:
		if t, ok := msg.Data.(time.Time); ok {
			return json.Marshal(t.Format(time.RFC3339Nano))
		}
	case flow.TypeBinary:
		if bs, ok := msg.Data.([]byte); ok {
			return json.Marshal(base64.StdEncoding.EncodeToString(bs))
		}
	}
	return json.Marshal(msg.Data)
}
