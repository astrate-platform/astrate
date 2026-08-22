// Package flow implements stream-based message routing through a block graph,
// mirroring Astarte Flow's core concepts: messages belong to streams (identified
// by key), streams are processed in order within a lane, and different streams
// may interleave across lanes.
package flow

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// WireSchema is the JSON schema identifier used on the wire.
const WireSchema = "astarte_flow/message/v0.1"

// DataType enumerates the wire-level value types Astarte Flow supports.
type DataType uint8

// DataType values, in wire-encoding order. TypeMap is the only aggregate;
// its per-field types live in FieldTypes and FieldSubtypes.
const (
	TypeInteger DataType = iota
	TypeReal
	TypeBoolean
	TypeDatetime
	TypeBinary
	TypeString
	// TypeMap indicates the message carries a map payload. Per-field types
	// and subtypes are in FieldTypes and FieldSubtypes.
	TypeMap
)

// dataTypeString converts a DataType to its wire-format string.
func dataTypeString(dt DataType) string {
	switch dt {
	case TypeInteger:
		return "integer"
	case TypeReal:
		return "real"
	case TypeBoolean:
		return "boolean"
	case TypeDatetime:
		return "datetime"
	case TypeBinary:
		return "binary"
	case TypeString:
		return "string"
	case TypeMap:
		return "map"
	default:
		return "unknown"
	}
}

// parseDataType converts a wire-format type string to DataType.
func parseDataType(s string) (DataType, error) {
	switch s {
	case "integer":
		return TypeInteger, nil
	case "real":
		return TypeReal, nil
	case "boolean":
		return TypeBoolean, nil
	case "datetime":
		return TypeDatetime, nil
	case "binary":
		return TypeBinary, nil
	case "string":
		return TypeString, nil
	default:
		return 0, fmt.Errorf("flow: unknown data type %q", s)
	}
}

// Message is one unit of data flowing through a block graph. Every message
// carries a Key that identifies its stream; messages with the same key are
// processed in submission order by the same lane (consistent hashing).
type Message struct {
	// Key identifies the stream this message belongs to. It must be non-empty.
	Key string
	// Metadata is an optional string→string map carried alongside the payload.
	Metadata map[string]string
	// Type is the base data type of the payload.
	Type DataType
	// Subtype is an optional MIME hint (meaningful when Type is TypeBinary).
	Subtype string
	// Timestamp is the event-time in microseconds since Unix epoch.
	Timestamp int64
	// Data is the payload; its concrete type must match Type (int64 for
	// TypeInteger, float64 for TypeReal, bool for TypeBoolean, time.Time for
	// TypeDatetime, []byte for TypeBinary, string for TypeString, map[string]any
	// for TypeMap).
	Data any
	// FieldTypes holds per-field types when Type is TypeMap.
	FieldTypes map[string]DataType
	// FieldSubtypes holds per-field subtypes when Type is TypeMap.
	FieldSubtypes map[string]string
}

// wireMessage is the on-the-wire JSON representation matching the upstream
// astarte_flow/message/v0.1 schema.
type wireMessage struct {
	Schema      string            `json:"schema"`
	Key         string            `json:"key"`
	Type        any               `json:"type"`
	Subtype     any               `json:"subtype,omitempty"`
	Data        any               `json:"data"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	TimestampUs int64             `json:"timestamp_us"`
}

// MarshalJSON serialises the Message to the upstream JSON wire format.
func (m *Message) MarshalJSON() ([]byte, error) {
	w := wireMessage{
		Schema:      WireSchema,
		Key:         m.Key,
		Metadata:    m.Metadata,
		TimestampUs: m.Timestamp,
	}

	switch m.Type {
	case TypeMap:
		w.Type = m.fieldTypesWire()
		if st := m.fieldSubtypesWire(); st != nil {
			w.Subtype = st
		}
		w.Data = m.dataWireMap()
	default:
		w.Type = dataTypeString(m.Type)
		if m.Subtype != "" {
			w.Subtype = m.Subtype
		}
		w.Data = m.dataWireScalar()
	}

	return json.Marshal(w)
}

// UnmarshalJSON deserialises a Message from the upstream JSON wire format.
func (m *Message) UnmarshalJSON(b []byte) error {
	var w wireMessage
	if err := json.Unmarshal(b, &w); err != nil {
		return fmt.Errorf("flow: unmarshal message: %w", err)
	}

	if w.Schema != WireSchema {
		return fmt.Errorf("flow: unsupported schema %q", w.Schema)
	}

	m.Key = w.Key
	m.Metadata = w.Metadata
	m.Timestamp = w.TimestampUs

	// Parse type field — either a string or a map of field types.
	switch t := w.Type.(type) {
	case string:
		dt, err := parseDataType(t)
		if err != nil {
			return err
		}
		m.Type = dt
		if s, ok := w.Subtype.(string); ok {
			m.Subtype = s
		}
		return m.setDataFromWire(dt, w.Data)

	case map[string]any:
		m.Type = TypeMap
		m.FieldTypes = make(map[string]DataType, len(t))
		for k, v := range t {
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf("flow: field type for %q is not a string", k)
			}
			dt, err := parseDataType(s)
			if err != nil {
				return err
			}
			m.FieldTypes[k] = dt
		}
		if sm, ok := w.Subtype.(map[string]any); ok {
			m.FieldSubtypes = make(map[string]string, len(sm))
			for k, v := range sm {
				s, ok := v.(string)
				if !ok {
					return fmt.Errorf("flow: field subtype for %q is not a string", k)
				}
				m.FieldSubtypes[k] = s
			}
		}
		return m.setDataFromWireMap(w.Data)

	default:
		return fmt.Errorf("flow: type field is neither string nor map")
	}
}

// dataWireScalar returns the wire-format data value for non-map messages.
func (m *Message) dataWireScalar() any {
	switch m.Type {
	case TypeBinary:
		if bs, ok := m.Data.([]byte); ok {
			return base64.StdEncoding.EncodeToString(bs)
		}
	case TypeDatetime:
		if t, ok := m.Data.(time.Time); ok {
			return t.Format(time.RFC3339Nano)
		}
	}
	return m.Data
}

// dataWireMap returns the wire-format data value for map messages.
func (m *Message) dataWireMap() any {
	raw, ok := m.Data.(map[string]any)
	if !ok {
		return m.Data
	}
	out := make(map[string]any, len(raw))
	for k, v := range raw {
		dt := m.FieldTypes[k]
		switch dt {
		case TypeBinary:
			if bs, ok := v.([]byte); ok {
				out[k] = base64.StdEncoding.EncodeToString(bs)
				continue
			}
		case TypeDatetime:
			if t, ok := v.(time.Time); ok {
				out[k] = t.Format(time.RFC3339Nano)
				continue
			}
		}
		out[k] = v
	}
	return out
}

// fieldTypesWire returns the wire-format type field for map messages.
func (m *Message) fieldTypesWire() map[string]string {
	out := make(map[string]string, len(m.FieldTypes))
	for k, dt := range m.FieldTypes {
		out[k] = dataTypeString(dt)
	}
	return out
}

// fieldSubtypesWire returns the wire-format subtype field for map messages.
func (m *Message) fieldSubtypesWire() map[string]string {
	if len(m.FieldSubtypes) == 0 {
		return nil
	}
	out := make(map[string]string, len(m.FieldSubtypes))
	for k, st := range m.FieldSubtypes {
		out[k] = st
	}
	return out
}

// setDataFromWire decodes a scalar wire data value into m.Data.
func (m *Message) setDataFromWire(dt DataType, raw any) error {
	switch dt {
	case TypeInteger:
		switch v := raw.(type) {
		case float64:
			m.Data = int64(v)
		case json.Number:
			n, err := v.Int64()
			if err != nil {
				return fmt.Errorf("flow: integer data: %w", err)
			}
			m.Data = n
		default:
			return fmt.Errorf("flow: integer data: expected number, got %T", raw)
		}
	case TypeReal:
		switch v := raw.(type) {
		case float64:
			m.Data = v
		case json.Number:
			f, err := v.Float64()
			if err != nil {
				return fmt.Errorf("flow: real data: %w", err)
			}
			m.Data = f
		default:
			return fmt.Errorf("flow: real data: expected number, got %T", raw)
		}
	case TypeBoolean:
		b, ok := raw.(bool)
		if !ok {
			return fmt.Errorf("flow: boolean data: expected bool, got %T", raw)
		}
		m.Data = b
	case TypeDatetime:
		s, ok := raw.(string)
		if !ok {
			return fmt.Errorf("flow: datetime data: expected string, got %T", raw)
		}
		t, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return fmt.Errorf("flow: datetime data: %w", err)
		}
		m.Data = t
	case TypeBinary:
		s, ok := raw.(string)
		if !ok {
			return fmt.Errorf("flow: binary data: expected string, got %T", raw)
		}
		bs, err := base64.StdEncoding.DecodeString(s)
		if err != nil {
			return fmt.Errorf("flow: binary data: %w", err)
		}
		m.Data = bs
	case TypeString:
		s, ok := raw.(string)
		if !ok {
			return fmt.Errorf("flow: string data: expected string, got %T", raw)
		}
		m.Data = s
	default:
		return fmt.Errorf("flow: unsupported data type %d", dt)
	}
	return nil
}

// setDataFromWireMap decodes a map wire data value into m.Data.
func (m *Message) setDataFromWireMap(raw any) error {
	rm, ok := raw.(map[string]any)
	if !ok {
		return fmt.Errorf("flow: map data: expected object, got %T", raw)
	}
	out := make(map[string]any, len(rm))
	for k, v := range rm {
		dt := m.FieldTypes[k]
		switch dt {
		case TypeBinary:
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf("flow: map field %q: expected base64 string for binary", k)
			}
			bs, err := base64.StdEncoding.DecodeString(s)
			if err != nil {
				return fmt.Errorf("flow: map field %q: %w", k, err)
			}
			out[k] = bs
		case TypeDatetime:
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf("flow: map field %q: expected string for datetime", k)
			}
			t, err := time.Parse(time.RFC3339Nano, s)
			if err != nil {
				return fmt.Errorf("flow: map field %q: %w", k, err)
			}
			out[k] = t
		default:
			out[k] = v
		}
	}
	m.Data = out
	return nil
}
