package payload

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// This file implements the Astrate JSON payload profile (docs/DESIGN.md
// §3.5.3) — the strict, documented contract for constrained clients (AtomVM
// and friends):
//
//	{ "v": <value>, "t": "2026-06-10T12:34:56.789Z" }
//
//   - The envelope is mandatory: a bare value (`22.5`) is rejected, keeping
//     symmetry with BSON and `t` unambiguous.
//   - `t` is optional: RFC 3339 string or integer milliseconds since epoch.
//   - Scalars map by the declared interface type: double/integer ← number;
//     longinteger ← number or decimal string (JS 2^53 safety); boolean ←
//     bool; string ← string; binaryblob ← base64 (standard alphabet, padded)
//     string; datetime ← RFC 3339 string or epoch-ms number; arrays ← JSON
//     arrays; object aggregation ← `v` is an object of last-level keys.
//   - Extra envelope members are tolerated and ignored, matching the BSON
//     decoder's lookup semantics. JSON null is never a value: property unset
//     is an empty payload.

// maxSafeJSONInt is 2^53 — the largest integer magnitude that both float64
// and a JavaScript number represent exactly. Beyond it, longintegers must
// travel as decimal strings (§3.5.3).
const maxSafeJSONInt = int64(1) << 53

// rfc3339Milli is the outbound timestamp layout: RFC 3339, UTC, fixed
// millisecond precision (the Astarte datetime resolution).
const rfc3339Milli = "2006-01-02T15:04:05.000Z07:00"

// jsonEnvelope captures the `{v, t}` members as raw tokens; per-type
// decoding happens against the mapping's declared ValueType.
type jsonEnvelope struct {
	V json.RawMessage `json:"v"`
	T json.RawMessage `json:"t"`
}

// decodeJSONEnvelope parses p as a JSON-profile envelope and extracts the
// raw `v` member plus the optional decoded `t` timestamp.
func decodeJSONEnvelope(p []byte) (json.RawMessage, *time.Time, error) {
	var env jsonEnvelope
	if err := json.Unmarshal(p, &env); err != nil {
		return nil, nil, rejectf(ReasonMalformed, "invalid JSON document: %v", err)
	}
	if env.V == nil {
		return nil, nil, rejectf(ReasonNoValue, `JSON document has no "v" member`)
	}
	if env.T == nil {
		return env.V, nil, nil // t absent: reception time applies (decided by the facade).
	}
	ts, err := parseJSONTimestamp(env.T)
	if err != nil {
		return nil, nil, err
	}
	return env.V, &ts, nil
}

// parseJSONTimestamp decodes the envelope `t` member, re-tagging every
// failure as ReasonBadTimestamp (a broken `t` is a timestamp problem, not a
// value-type mismatch).
func parseJSONTimestamp(raw json.RawMessage) (time.Time, error) {
	ts, err := decodeJSONDateTime(raw)
	if err != nil {
		var re *RejectError
		if errors.As(err, &re) && re.Reason != ReasonBadTimestamp {
			return time.Time{}, &RejectError{Reason: ReasonBadTimestamp, Detail: `"t": ` + re.Detail}
		}
		return time.Time{}, err
	}
	return ts, nil
}

// decodeJSONValue coerces one raw JSON token to the declared mapping type
// under the docs/DESIGN.md §2.6 step 5 rules and the §3.5.3 profile.
func decodeJSONValue(raw json.RawMessage, vt interfaceschema.ValueType) (Value, error) {
	if isJSONNull(raw) {
		return nil, rejectf(ReasonTypeMismatch, "JSON null is not a value (property unset is an empty payload)")
	}
	if vt.IsArray() {
		if len(raw) == 0 || raw[0] != '[' {
			return nil, rejectf(ReasonTypeMismatch, "JSON value %s does not coerce to %s", clip(raw), vt)
		}
		var elems []json.RawMessage
		if err := json.Unmarshal(raw, &elems); err != nil {
			return nil, rejectf(ReasonMalformed, "invalid JSON array: %v", err)
		}
		return decodeArray(vt, len(elems), func(i int) (Value, error) {
			return decodeJSONValue(elems[i], vt.Elem())
		})
	}

	switch vt {
	case interfaceschema.Double:
		num, err := jsonNumber(raw)
		if err != nil {
			return nil, err
		}
		f, err := strconv.ParseFloat(string(num), 64)
		if err != nil {
			return nil, rejectf(ReasonTypeMismatch, "number %s does not fit double", num)
		}
		return asValue(checkDouble(f))
	case interfaceschema.Integer:
		num, err := jsonNumber(raw)
		if err != nil {
			return nil, err
		}
		i, err := jsonInt64(num)
		if err != nil {
			return nil, err
		}
		return asValue(int32FromInt64(i))
	case interfaceschema.LongInteger:
		// Number, or decimal string for values beyond 2^53 (§3.5.3).
		if len(raw) > 0 && raw[0] == '"' {
			var s string
			if err := json.Unmarshal(raw, &s); err != nil {
				return nil, rejectf(ReasonMalformed, "invalid JSON string: %v", err)
			}
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, rejectf(ReasonTypeMismatch, "string %q is not a decimal longinteger", s)
			}
			return i, nil
		}
		num, err := jsonNumber(raw)
		if err != nil {
			return nil, err
		}
		return asValue(jsonInt64(num))
	case interfaceschema.Boolean:
		var b bool
		if err := json.Unmarshal(raw, &b); err != nil {
			return nil, rejectf(ReasonTypeMismatch, "JSON value %s is not a boolean", clip(raw))
		}
		return b, nil
	case interfaceschema.String:
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return nil, rejectf(ReasonTypeMismatch, "JSON value %s is not a string", clip(raw))
		}
		return asValue(checkString(s))
	case interfaceschema.BinaryBlob:
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return nil, rejectf(ReasonTypeMismatch, "binaryblob must be a base64 JSON string, got %s", clip(raw))
		}
		data, err := base64.StdEncoding.Strict().DecodeString(s)
		if err != nil {
			return nil, rejectf(ReasonTypeMismatch, "invalid base64 (standard alphabet, padded required): %v", err)
		}
		return data, nil
	case interfaceschema.DateTime:
		return asValue(decodeJSONDateTime(raw))
	}
	return nil, rejectf(ReasonTypeMismatch, "JSON value %s does not coerce to %s", clip(raw), vt)
}

// decodeJSONDateTime decodes a datetime token: RFC 3339 string or integer
// epoch-milliseconds number, bounds-checked to the accepted window.
func decodeJSONDateTime(raw json.RawMessage) (time.Time, error) {
	if isJSONNull(raw) || len(raw) == 0 {
		return time.Time{}, rejectf(ReasonTypeMismatch, "JSON null is not a datetime")
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return time.Time{}, rejectf(ReasonMalformed, "invalid JSON string: %v", err)
		}
		ts, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return time.Time{}, rejectf(ReasonTypeMismatch, "%q is not an RFC 3339 timestamp", s)
		}
		return checkDateTime(ts)
	}
	num, err := jsonNumber(raw)
	if err != nil {
		return time.Time{}, rejectf(ReasonTypeMismatch,
			"JSON value %s is not an RFC 3339 string or epoch-ms number", clip(raw))
	}
	ms, err := jsonInt64(num)
	if err != nil {
		return time.Time{}, err
	}
	return dateTimeFromMillis(ms)
}

// decodeJSONObject decodes an object-aggregation `v`: a JSON object of
// last-level endpoint names, each of which must resolve in leaves
// (docs/DESIGN.md §2.6 step 6).
func decodeJSONObject(raw json.RawMessage, leaves map[string]*interfaceschema.CompiledMapping) (Value, error) {
	if isJSONNull(raw) || len(raw) == 0 || raw[0] != '{' {
		return nil, rejectf(ReasonBadObject, "JSON value %s where an object-aggregation document was expected", clip(raw))
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, rejectf(ReasonMalformed, "invalid JSON object: %v", err)
	}
	if len(obj) == 0 {
		return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
	}
	out := make(map[string]Value, len(obj))
	for key, eraw := range obj {
		leaf, ok := leaves[key]
		if !ok || leaf == nil {
			return nil, rejectf(ReasonBadObject, "key %q matches no declared object leaf", key)
		}
		val, err := decodeJSONValue(eraw, leaf.ValueType)
		if err != nil {
			return nil, annotate(err, "key %q", key)
		}
		out[key] = val
	}
	return out, nil
}

// ---------------------------------------------------------------------------
// JSON token helpers
// ---------------------------------------------------------------------------

// isJSONNull reports whether raw is the JSON null literal (json.RawMessage
// holds the exact value token, so a direct comparison suffices).
func isJSONNull(raw json.RawMessage) bool { return string(raw) == "null" }

// jsonNumber extracts a JSON number token; strings (including quoted
// numbers, which encoding/json would otherwise accept into a json.Number),
// booleans, and structured values are rejected.
func jsonNumber(raw json.RawMessage) (json.Number, error) {
	if len(raw) == 0 || (raw[0] != '-' && (raw[0] < '0' || raw[0] > '9')) {
		return "", rejectf(ReasonTypeMismatch, "JSON value %s is not a number", clip(raw))
	}
	var num json.Number
	if err := json.Unmarshal(raw, &num); err != nil || num == "" {
		return "", rejectf(ReasonTypeMismatch, "JSON value %s is not a number", clip(raw))
	}
	return num, nil
}

// jsonInt64 converts a JSON number token to int64. Plain decimal-integer
// tokens convert exactly across the whole int64 range. Fraction or exponent
// forms are accepted only when integral and strictly inside the float64-exact
// window (|x| < 2^53): beyond it the decimal→binary parse itself may round,
// silently corrupting the value — such values must use the decimal-string
// form (§3.5.3).
func jsonInt64(num json.Number) (int64, error) {
	i, err := strconv.ParseInt(string(num), 10, 64)
	if err == nil {
		return i, nil
	}
	if errors.Is(err, strconv.ErrRange) {
		return 0, rejectf(ReasonTypeMismatch, "number %s overflows int64", num)
	}
	f, err := strconv.ParseFloat(string(num), 64)
	if err != nil {
		return 0, rejectf(ReasonTypeMismatch, "number %s is not decodable", num)
	}
	if f != math.Trunc(f) || f <= -float64(maxSafeJSONInt) || f >= float64(maxSafeJSONInt) {
		return 0, rejectf(ReasonTypeMismatch, "number %s does not coerce to an integer exactly", num)
	}
	return int64(f), nil
}

// clip truncates a raw token for use in human-readable reject details.
func clip(raw []byte) string {
	const maxDetail = 64
	if len(raw) > maxDetail {
		return string(raw[:maxDetail]) + "…"
	}
	return string(raw)
}

// ---------------------------------------------------------------------------
// Outbound encoder
// ---------------------------------------------------------------------------

// encodeJSON builds the outbound JSON-profile `{v, t}` document. Member
// order is always v first, then t; timestamps are RFC 3339 UTC with fixed
// millisecond precision; longintegers beyond ±2^53 are emitted as decimal
// strings so JavaScript-grade parsers never round them.
func encodeJSON(v Value, ts *time.Time) ([]byte, error) {
	jv, err := jsonifyValue(v)
	if err != nil {
		return nil, err
	}
	env := struct {
		V any    `json:"v"`
		T string `json:"t,omitempty"`
	}{V: jv}
	if ts != nil {
		t, err := checkDateTime(*ts)
		if err != nil {
			return nil, err
		}
		env.T = t.Format(rfc3339Milli)
	}
	return json.Marshal(env)
}

// jsonifyValue converts a closed-set Value into its JSON-profile
// representation, applying the same caps as decoding.
func jsonifyValue(v Value) (any, error) {
	switch x := v.(type) {
	case float64:
		if _, err := checkDouble(x); err != nil {
			return nil, err
		}
		return x, nil
	case int32:
		return x, nil
	case int64:
		return jsonifyInt64(x), nil
	case bool:
		return x, nil
	case string:
		if _, err := checkString(x); err != nil {
			return nil, err
		}
		return x, nil
	case []byte:
		return base64.StdEncoding.EncodeToString(x), nil
	case time.Time:
		t, err := checkDateTime(x)
		if err != nil {
			return nil, err
		}
		return t.Format(rfc3339Milli), nil
	case []float64:
		return jsonifySlice(x)
	case []int32:
		return jsonifySlice(x)
	case []int64:
		return jsonifySlice(x)
	case []bool:
		return jsonifySlice(x)
	case []string:
		return jsonifySlice(x)
	case [][]byte:
		return jsonifySlice(x)
	case []time.Time:
		return jsonifySlice(x)
	case map[string]Value:
		if len(x) == 0 {
			return nil, rejectf(ReasonBadObject, "object-aggregation document is empty")
		}
		out := make(map[string]any, len(x))
		for k, e := range x {
			if _, nested := e.(map[string]Value); nested {
				return nil, rejectf(ReasonBadObject, "key %q: nested objects are not encodable", k)
			}
			jv, err := jsonifyValue(e)
			if err != nil {
				return nil, annotate(err, "key %q", k)
			}
			out[k] = jv
		}
		return out, nil // json.Marshal sorts map keys: deterministic output
	default:
		return nil, rejectf(ReasonTypeMismatch, "Go type %T is not an encodable payload value", v)
	}
}

// jsonifySlice converts a homogeneous array, enforcing the MaxArrayLen cap.
func jsonifySlice[T any](xs []T) (any, error) {
	if len(xs) > MaxArrayLen {
		return nil, rejectf(ReasonValueTooLarge, "array of %d elements exceeds %d element cap", len(xs), MaxArrayLen)
	}
	out := make([]any, len(xs))
	for i, e := range xs {
		jv, err := jsonifyValue(e)
		if err != nil {
			return nil, annotate(err, "element %d", i)
		}
		out[i] = jv
	}
	return out, nil
}

// jsonifyInt64 emits a longinteger as a JSON number when exactly
// representable by JavaScript-grade parsers (|x| ≤ 2^53), and as a decimal
// string beyond that (§3.5.3).
func jsonifyInt64(x int64) any {
	if x > maxSafeJSONInt || x < -maxSafeJSONInt {
		return strconv.FormatInt(x, 10)
	}
	return json.Number(strconv.FormatInt(x, 10))
}
