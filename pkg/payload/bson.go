package payload

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// This file implements the BSON side of the codec (docs/DESIGN.md §3.5.5):
// decoding uses the mongo-driver raw-document API exclusively — bson.Raw
// element lookups, no reflection, no intermediate maps — so the hot path
// stays near zero allocations. Encoding is a hand-rolled `{v, t}` document
// writer with a deterministic field order (v first, then t) and the exact
// element types the official Go SDK produces (double, int32, int64, boolean,
// UTF-8 string, generic binary, UTC datetime, arrays, embedded documents).

// decodeBSONEnvelope validates p as a single BSON document and extracts the
// raw `v` element plus the optional `t` UTC-datetime timestamp. The caller
// must have classified p as FormatBSON; the leading guard repeats the cheap
// structural checks so a misrouted payload degrades into ReasonMalformed
// instead of undefined driver behaviour.
func decodeBSONEnvelope(p []byte) (bson.RawValue, *time.Time, error) {
	if len(p) < 5 || int64(binary.LittleEndian.Uint32(p[0:4])) != int64(len(p)) {
		return bson.RawValue{}, nil, rejectf(ReasonMalformed, "BSON length prefix does not match the %d byte payload", len(p))
	}
	raw := bson.Raw(p)
	if err := raw.Validate(); err != nil {
		return bson.RawValue{}, nil, rejectf(ReasonMalformed, "invalid BSON document: %v", err)
	}
	v, err := raw.LookupErr("v")
	if err != nil {
		return bson.RawValue{}, nil, rejectf(ReasonNoValue, `BSON document has no "v" element`)
	}
	tRaw, err := raw.LookupErr("t")
	if err != nil {
		return v, nil, nil // t absent: reception time applies (decided by the facade).
	}
	ms, ok := tRaw.DateTimeOK()
	if !ok {
		return bson.RawValue{}, nil, rejectf(ReasonBadTimestamp, `BSON "t" element is %s, want UTC datetime`, tRaw.Type)
	}
	ts, err := dateTimeFromMillis(ms)
	if err != nil {
		return bson.RawValue{}, nil, err
	}
	return v, &ts, nil
}

// decodeBSONValue coerces one raw BSON value to the declared mapping type
// under the docs/DESIGN.md §2.6 step 5 rules.
func decodeBSONValue(rv bson.RawValue, vt interfaceschema.ValueType) (Value, error) {
	if vt.IsArray() {
		arr, ok := rv.ArrayOK()
		if !ok {
			return nil, rejectf(ReasonTypeMismatch, "BSON %s does not coerce to %s", rv.Type, vt)
		}
		vals, err := arr.Values()
		if err != nil {
			return nil, rejectf(ReasonMalformed, "invalid BSON array: %v", err)
		}
		return decodeArray(vt, len(vals), func(i int) (Value, error) {
			return decodeBSONValue(vals[i], vt.Elem())
		})
	}

	switch vt {
	case interfaceschema.Double:
		// double ← BSON double; int32/int64 widen only when lossless.
		switch rv.Type {
		case bson.TypeDouble:
			f, _ := rv.DoubleOK()
			return asValue(checkDouble(f))
		case bson.TypeInt32:
			i, _ := rv.Int32OK()
			return float64(i), nil // int32 always widens losslessly
		case bson.TypeInt64:
			i, _ := rv.Int64OK()
			return asValue(doubleFromInt64(i))
		}
	case interfaceschema.Integer:
		// integer ← int32, or int64/double that fits int32 exactly.
		switch rv.Type {
		case bson.TypeInt32:
			i, _ := rv.Int32OK()
			return i, nil
		case bson.TypeInt64:
			i, _ := rv.Int64OK()
			return asValue(int32FromInt64(i))
		case bson.TypeDouble:
			f, _ := rv.DoubleOK()
			return asValue(int32FromFloat64(f))
		}
	case interfaceschema.LongInteger:
		// longinteger ← int64/int32 only (no double form in BSON, §2.6).
		switch rv.Type {
		case bson.TypeInt64:
			i, _ := rv.Int64OK()
			return i, nil
		case bson.TypeInt32:
			i, _ := rv.Int32OK()
			return int64(i), nil
		}
	case interfaceschema.Boolean:
		if b, ok := rv.BooleanOK(); ok {
			return b, nil
		}
	case interfaceschema.String:
		if s, ok := rv.StringValueOK(); ok {
			return asValue(checkString(s))
		}
	case interfaceschema.BinaryBlob:
		// Any binary subtype is accepted as raw bytes; the data sub-slice is
		// cloned because the broker may reuse the payload buffer.
		if _, data, ok := rv.BinaryOK(); ok {
			return bytes.Clone(data), nil
		}
	case interfaceschema.DateTime:
		if ms, ok := rv.DateTimeOK(); ok {
			return asValue(dateTimeFromMillis(ms))
		}
	}
	return nil, rejectf(ReasonTypeMismatch, "BSON %s does not coerce to %s", rv.Type, vt)
}

// decodeBSONObject decodes an object-aggregation `v`: an embedded document
// of last-level endpoint names, each of which must resolve in leaves
// (docs/DESIGN.md §2.6 step 6).
func decodeBSONObject(rv bson.RawValue, leaves map[string]*interfaceschema.CompiledMapping) (Value, error) {
	doc, ok := rv.DocumentOK()
	if !ok {
		return nil, rejectf(ReasonBadObject, "BSON %s where an object-aggregation document was expected", rv.Type)
	}
	elems, err := doc.Elements()
	if err != nil {
		return nil, rejectf(ReasonMalformed, "invalid BSON object document: %v", err)
	}
	if len(elems) == 0 {
		return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
	}
	out := make(map[string]Value, len(elems))
	for _, el := range elems {
		key, err := el.KeyErr()
		if err != nil {
			return nil, rejectf(ReasonMalformed, "invalid BSON object element: %v", err)
		}
		leaf, ok := leaves[key]
		if !ok || leaf == nil {
			return nil, rejectf(ReasonBadObject, "key %q matches no declared object leaf", key)
		}
		if _, dup := out[key]; dup {
			return nil, rejectf(ReasonBadObject, "duplicate object key %q", key)
		}
		ev, err := el.ValueErr()
		if err != nil {
			return nil, rejectf(ReasonMalformed, "invalid BSON object element %q: %v", key, err)
		}
		val, err := decodeBSONValue(ev, leaf.ValueType)
		if err != nil {
			return nil, annotate(err, "key %q", key)
		}
		out[key] = val
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// Outbound encoder
// ---------------------------------------------------------------------------

// encodeBSON builds the outbound BSON `{v, t}` document. Element order is
// always v first, then t — the same envelope shape the official SDKs emit.
// Values are validated with the same caps as decoding (string/array limits,
// NaN/±Inf, datetime window) so a malformed publish is rejected here rather
// than on the device.
func encodeBSON(v Value, ts *time.Time) ([]byte, error) {
	doc := make([]byte, 4, 64)
	var err error
	doc, err = appendBSONElement(doc, "v", v)
	if err != nil {
		return nil, err
	}
	if ts != nil {
		t, err := checkDateTime(*ts)
		if err != nil {
			return nil, err
		}
		doc = appendBSONHeader(doc, byte(bson.TypeDateTime), "t")
		doc = appendBSONMillis(doc, t.UnixMilli())
	}
	return finishBSONDoc(doc, 0)
}

// appendBSONHeader appends one element prelude: type byte plus NUL-terminated
// key. Keys are "v"/"t", decimal array indexes, or object keys pre-checked by
// appendBSONObject, so they never contain NUL bytes.
func appendBSONHeader(dst []byte, typ byte, key string) []byte {
	dst = append(dst, typ)
	dst = append(dst, key...)
	return append(dst, 0x00)
}

// appendBSONMillis appends a BSON UTC-datetime body (int64 little-endian
// milliseconds since the Unix epoch).
func appendBSONMillis(dst []byte, ms int64) []byte {
	return binary.LittleEndian.AppendUint64(dst, uint64(ms)) //nolint:gosec // two's-complement int64 is the BSON datetime wire encoding (negative = pre-1970)
}

// finishBSONDoc appends the document terminator and patches the int32 length
// prefix located at offset start.
func finishBSONDoc(doc []byte, start int) ([]byte, error) {
	doc = append(doc, 0x00)
	size := len(doc) - start
	if size > math.MaxInt32 {
		return nil, fmt.Errorf("payload: BSON document of %d bytes overflows the int32 length prefix", size)
	}
	binary.LittleEndian.PutUint32(doc[start:start+4], uint32(size)) //nolint:gosec // bounds-checked above
	return doc, nil
}

// appendBSONElement appends one fully-typed element for any value of the
// closed Value set (see Value). Unknown Go types are rejected.
func appendBSONElement(dst []byte, key string, v Value) ([]byte, error) {
	switch x := v.(type) {
	case float64:
		if _, err := checkDouble(x); err != nil {
			return nil, err
		}
		dst = appendBSONHeader(dst, byte(bson.TypeDouble), key)
		return binary.LittleEndian.AppendUint64(dst, math.Float64bits(x)), nil
	case int32:
		dst = appendBSONHeader(dst, byte(bson.TypeInt32), key)
		return binary.LittleEndian.AppendUint32(dst, uint32(x)), nil //nolint:gosec // two's-complement int32 is the BSON wire encoding
	case int64:
		dst = appendBSONHeader(dst, byte(bson.TypeInt64), key)
		return binary.LittleEndian.AppendUint64(dst, uint64(x)), nil //nolint:gosec // two's-complement int64 is the BSON wire encoding
	case bool:
		dst = appendBSONHeader(dst, byte(bson.TypeBoolean), key)
		if x {
			return append(dst, 0x01), nil
		}
		return append(dst, 0x00), nil
	case string:
		if _, err := checkString(x); err != nil {
			return nil, err
		}
		dst = appendBSONHeader(dst, byte(bson.TypeString), key)
		dst = binary.LittleEndian.AppendUint32(dst, uint32(len(x)+1)) //nolint:gosec // bounded by MaxStringLen via checkString
		dst = append(dst, x...)
		return append(dst, 0x00), nil
	case []byte:
		if len(x) > math.MaxInt32 {
			return nil, rejectf(ReasonValueTooLarge, "binary blob of %d bytes overflows the int32 length prefix", len(x))
		}
		dst = appendBSONHeader(dst, byte(bson.TypeBinary), key)
		dst = binary.LittleEndian.AppendUint32(dst, uint32(len(x))) //nolint:gosec // bounds-checked above
		dst = append(dst, bson.TypeBinaryGeneric)
		return append(dst, x...), nil
	case time.Time:
		t, err := checkDateTime(x)
		if err != nil {
			return nil, err
		}
		dst = appendBSONHeader(dst, byte(bson.TypeDateTime), key)
		return appendBSONMillis(dst, t.UnixMilli()), nil
	case []float64:
		return appendBSONSlice(dst, key, x)
	case []int32:
		return appendBSONSlice(dst, key, x)
	case []int64:
		return appendBSONSlice(dst, key, x)
	case []bool:
		return appendBSONSlice(dst, key, x)
	case []string:
		return appendBSONSlice(dst, key, x)
	case [][]byte:
		return appendBSONSlice(dst, key, x)
	case []time.Time:
		return appendBSONSlice(dst, key, x)
	case map[string]Value:
		return appendBSONObject(dst, key, x)
	default:
		return nil, rejectf(ReasonTypeMismatch, "Go type %T is not an encodable payload value", v)
	}
}

// appendBSONSlice appends a homogeneous array element ("0", "1", … keys, the
// BSON array convention), enforcing the MaxArrayLen cap.
func appendBSONSlice[T any](dst []byte, key string, xs []T) ([]byte, error) {
	if len(xs) > MaxArrayLen {
		return nil, rejectf(ReasonValueTooLarge, "array of %d elements exceeds %d element cap", len(xs), MaxArrayLen)
	}
	dst = appendBSONHeader(dst, byte(bson.TypeArray), key)
	start := len(dst)
	dst = append(dst, 0, 0, 0, 0)
	for i, e := range xs {
		var err error
		dst, err = appendBSONElement(dst, strconv.Itoa(i), e)
		if err != nil {
			return nil, annotate(err, "element %d", i)
		}
	}
	return finishBSONDoc(dst, start)
}

// appendBSONObject appends an object-aggregation embedded document. Keys are
// emitted in sorted order so encoding is deterministic; nested objects are
// rejected (object aggregation is one level of last-level names, §2.6).
func appendBSONObject(dst []byte, key string, obj map[string]Value) ([]byte, error) {
	if len(obj) == 0 {
		return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
	}
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	dst = appendBSONHeader(dst, byte(bson.TypeEmbeddedDocument), key)
	start := len(dst)
	dst = append(dst, 0, 0, 0, 0)
	for _, k := range keys {
		if k == "" || strings.IndexByte(k, 0x00) >= 0 {
			return nil, rejectf(ReasonBadObject, "invalid object key %q", k)
		}
		if _, nested := obj[k].(map[string]Value); nested {
			return nil, rejectf(ReasonBadObject, "key %q: nested objects are not encodable", k)
		}
		var err error
		dst, err = appendBSONElement(dst, k, obj[k])
		if err != nil {
			return nil, annotate(err, "key %q", k)
		}
	}
	return finishBSONDoc(dst, start)
}
