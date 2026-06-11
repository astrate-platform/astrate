package payload

import (
	"errors"
	"fmt"
	"time"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// This file is the package facade (docs/ROADMAP.md §2.4 file 1.19):
// internal/engine and the AppEngine API talk to the codec exclusively
// through Decoder.Individual / Decoder.Object (or the package-level
// shorthands) and Encode. Every rejection is a *RejectError whose
// RejectReason keys the engine's per-reason metrics and device_error
// triggers; plain errors signal caller bugs (nil mapping, bad format
// argument), never payload content.

// Decoder decodes inbound data payloads. The zero value is ready to use
// with the default limits.
type Decoder struct {
	// MaxSize caps the accepted payload size in bytes for both formats
	// (docs/DESIGN.md §3.5.3); 0 means DefaultMaxSize.
	MaxSize int
}

// limit returns the effective payload size cap.
func (d Decoder) limit() int {
	if d.MaxSize > 0 {
		return d.MaxSize
	}
	return DefaultMaxSize
}

// gate applies the size cap and format sniff shared by both decode paths.
func (d Decoder) gate(p []byte) (Format, error) {
	if len(p) > d.limit() {
		return FormatInvalid, rejectf(ReasonTooLarge, "%d byte payload exceeds %d byte cap", len(p), d.limit())
	}
	f := DetectFormat(p)
	if f == FormatInvalid {
		return f, rejectf(ReasonUnknownFormat, "payload is neither empty, BSON, nor JSON")
	}
	return f, nil
}

// Individual decodes p against an individual-aggregation mapping: BSON or
// JSON `{v, t}` envelope, value coerced to m.ValueType, empty payload =
// property unset (allowed only with allow_unset). The explicit-timestamp
// policy follows docs/DESIGN.md §2.6 step 5: a mapping with
// explicit_timestamp requires `t`; otherwise a present `t` is tolerated and
// ignored (upstream leniency).
func (d Decoder) Individual(p []byte, m *interfaceschema.CompiledMapping) (DecodedPayload, error) {
	if m == nil {
		return DecodedPayload{}, errors.New("payload: nil mapping")
	}
	f, err := d.gate(p)
	if err != nil {
		return DecodedPayload{}, err
	}

	var (
		val Value
		ts  *time.Time
	)
	switch f {
	case FormatEmpty:
		if !m.AllowUnset {
			return DecodedPayload{}, rejectf(ReasonUnsetNotAllowed, "empty payload on a mapping without allow_unset")
		}
		return DecodedPayload{Format: FormatEmpty}, nil
	case FormatBSON:
		rv, t, err := decodeBSONEnvelope(p)
		if err != nil {
			return DecodedPayload{}, err
		}
		if val, err = decodeBSONValue(rv, m.ValueType); err != nil {
			return DecodedPayload{}, err
		}
		ts = t
	case FormatJSON:
		raw, t, err := decodeJSONEnvelope(p)
		if err != nil {
			return DecodedPayload{}, err
		}
		if val, err = decodeJSONValue(raw, m.ValueType); err != nil {
			return DecodedPayload{}, err
		}
		ts = t
	}
	return finishDecode(val, ts, f, m.ExplicitTimestamp)
}

// Object decodes p against an object-aggregated interface: `v` must be a
// document of last-level endpoint names, each resolving in leaves
// (CompiledInterface.ObjectLeaves). Object aggregation exists only on
// datastreams, so the empty (property-unset) payload is always rejected.
// The explicit-timestamp policy is taken from the leaves (uniform across an
// object-aggregated interface by construction).
func (d Decoder) Object(p []byte, leaves map[string]*interfaceschema.CompiledMapping) (DecodedPayload, error) {
	if len(leaves) == 0 {
		return DecodedPayload{}, errors.New("payload: no object leaves")
	}
	f, err := d.gate(p)
	if err != nil {
		return DecodedPayload{}, err
	}

	var (
		val Value
		ts  *time.Time
	)
	switch f {
	case FormatEmpty:
		return DecodedPayload{}, rejectf(ReasonUnsetNotAllowed, "empty payload on an object-aggregated datastream")
	case FormatBSON:
		rv, t, err := decodeBSONEnvelope(p)
		if err != nil {
			return DecodedPayload{}, err
		}
		if val, err = decodeBSONObject(rv, leaves); err != nil {
			return DecodedPayload{}, err
		}
		ts = t
	case FormatJSON:
		raw, t, err := decodeJSONEnvelope(p)
		if err != nil {
			return DecodedPayload{}, err
		}
		if val, err = decodeJSONObject(raw, leaves); err != nil {
			return DecodedPayload{}, err
		}
		ts = t
	}
	return finishDecode(val, ts, f, objectExplicitTimestamp(leaves))
}

// finishDecode applies the explicit-timestamp policy and assembles the
// result.
func finishDecode(val Value, ts *time.Time, f Format, explicit bool) (DecodedPayload, error) {
	if explicit && ts == nil {
		return DecodedPayload{}, rejectf(ReasonBadTimestamp, "mapping declares explicit_timestamp but payload carries no t")
	}
	if !explicit {
		ts = nil // tolerated and ignored (docs/DESIGN.md §2.6 step 5)
	}
	return DecodedPayload{Value: val, Timestamp: ts, Format: f}, nil
}

// objectExplicitTimestamp reads the explicit-timestamp attribute off any
// leaf; pkg/interfaceschema guarantees it is uniform across an
// object-aggregated interface.
func objectExplicitTimestamp(leaves map[string]*interfaceschema.CompiledMapping) bool {
	for _, m := range leaves {
		if m != nil {
			return m.ExplicitTimestamp
		}
	}
	return false
}

// Decode decodes an individual-aggregation payload with the default limits.
func Decode(p []byte, m *interfaceschema.CompiledMapping) (DecodedPayload, error) {
	return Decoder{}.Individual(p, m)
}

// DecodeObject decodes an object-aggregation payload with the default
// limits.
func DecodeObject(p []byte, leaves map[string]*interfaceschema.CompiledMapping) (DecodedPayload, error) {
	return Decoder{}.Object(p, leaves)
}

// Encode builds an outbound `{v, t}` document in the requested format
// (docs/DESIGN.md §3.5.4: BSON by default, JSON for devices hinted to the
// JSON profile). v must be one of the closed Value set — scalars, arrays,
// or map[string]Value for object aggregation; ts is the optional explicit
// timestamp. FormatEmpty encodes the property-unset payload: v and ts must
// be nil and the result is the empty payload.
func Encode(v Value, ts *time.Time, f Format) ([]byte, error) {
	switch f {
	case FormatEmpty:
		if v != nil || ts != nil {
			return nil, errors.New("payload: FormatEmpty encodes only a nil value (property unset)")
		}
		return []byte{}, nil
	case FormatBSON:
		if v == nil {
			return nil, errors.New("payload: nil value (use FormatEmpty for property unset)")
		}
		return encodeBSON(v, ts)
	case FormatJSON:
		if v == nil {
			return nil, errors.New("payload: nil value (use FormatEmpty for property unset)")
		}
		return encodeJSON(v, ts)
	default:
		return nil, fmt.Errorf("payload: cannot encode format %s", f)
	}
}
