// Package payload implements the dual-format Astarte data-payload codec
// (docs/DESIGN.md §3.5): standard BSON `{v, t}` documents as produced by the
// official Astarte device SDKs, and the strict Astrate JSON profile
// (§3.5.3) for constrained clients (AtomVM and friends) on the same topics
// with the same semantics.
//
// The package is pure: it sniffs the wire format (§3.5.2), decodes and
// validates a payload against a compiled interface mapping with the
// upstream coercion rules (docs/DESIGN.md §2.6 step 5), and encodes
// outbound `{v, t}` documents in either format (§3.5.4). Every rejection
// carries a typed RejectReason that internal/engine feeds into per-reason
// metrics and device_error triggers.
package payload

import (
	"errors"
	"fmt"
	"math"
	"time"
	"unicode/utf8"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// Format is the detected wire format of an inbound data payload
// (docs/DESIGN.md §3.5.2).
type Format uint8

// Format values returned by DetectFormat.
const (
	// FormatInvalid marks a payload that is neither empty, BSON, nor JSON.
	FormatInvalid Format = iota
	// FormatEmpty is the zero-length payload (property-unset semantics).
	FormatEmpty
	// FormatBSON is a BSON `{v, t}` document (official SDKs).
	FormatBSON
	// FormatJSON is an Astrate JSON-profile `{"v": ..., "t": ...}` document.
	FormatJSON
)

// String returns a stable lowercase label (also used as a metrics label).
func (f Format) String() string {
	switch f {
	case FormatInvalid:
		return "invalid"
	case FormatEmpty:
		return "empty"
	case FormatBSON:
		return "bson"
	case FormatJSON:
		return "json"
	default:
		return fmt.Sprintf("Format(%d)", uint8(f))
	}
}

// Value is a decoded payload value. It is always one of a closed set of Go
// types, determined by the mapping's ValueType:
//
//	double       float64        doublearray       []float64
//	integer      int32          integerarray      []int32
//	boolean      bool           booleanarray      []bool
//	longinteger  int64          longintegerarray  []int64
//	string       string         stringarray       []string
//	binaryblob   []byte         binaryblobarray   [][]byte
//	datetime     time.Time      datetimearray     []time.Time
//
// Object-aggregated payloads decode to map[string]Value keyed by the
// last-level endpoint name, each entry holding one of the types above.
// Encode accepts exactly the same set.
type Value = any

// DecodedPayload is the result of decoding one data payload.
type DecodedPayload struct {
	// Value is the decoded value (see Value); nil for an unset payload.
	Value Value
	// Timestamp is the explicit `t` timestamp, if the mapping declares
	// explicit_timestamp and the payload carried one; nil otherwise
	// (reception time applies).
	Timestamp *time.Time
	// Format is the wire format the payload arrived in. internal/engine
	// uses it to maintain the device's payload_format_hint
	// (docs/DESIGN.md §3.5.4).
	Format Format
}

// IsUnset reports whether the payload was the empty property-unset payload.
func (d DecodedPayload) IsUnset() bool { return d.Value == nil }

// Size and cardinality limits (docs/DESIGN.md §2.6 step 5, §3.5.3).
const (
	// DefaultMaxSize is the default accepted payload size for both
	// formats: 64 KiB (configurable per call site).
	DefaultMaxSize = 64 << 10
	// MaxStringLen is the maximum byte length of a decoded string value.
	MaxStringLen = 64 << 10
	// MaxArrayLen is the maximum number of elements in an array value.
	MaxArrayLen = 1024
)

// MinDateTime and MaxDateTime bound accepted datetime values (both payload
// values and explicit `t` timestamps). The window — years 0001 through 9999
// — is the intersection of what RFC 3339 (4-digit years), BSON UTC datetime,
// and PostgreSQL timestamptz all represent without surprises.
var (
	// MinDateTime is the earliest accepted instant: 0001-01-01T00:00:00Z.
	MinDateTime = time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	// MaxDateTime is the latest accepted instant: 9999-12-31T23:59:59.999Z.
	MaxDateTime = time.Date(9999, time.December, 31, 23, 59, 59, 999_000_000, time.UTC)
)

var (
	minDateTimeMillis = MinDateTime.UnixMilli()
	maxDateTimeMillis = MaxDateTime.UnixMilli()
)

// RejectReason classifies why a payload was rejected. internal/engine keys
// its per-reason rejection counters and device_error trigger events on it
// (docs/DESIGN.md §2.6 "failures are never silent").
type RejectReason uint8

// RejectReason values. The zero value means "not a rejection".
const (
	// ReasonNone is the zero value; never carried by an error.
	ReasonNone RejectReason = iota
	// ReasonTooLarge: payload exceeds the configured size cap.
	ReasonTooLarge
	// ReasonUnknownFormat: the payload is neither empty, BSON, nor JSON.
	ReasonUnknownFormat
	// ReasonMalformed: detected format, but not a decodable document.
	ReasonMalformed
	// ReasonNoValue: the `{v, t}` envelope carries no `v` field.
	ReasonNoValue
	// ReasonBadTimestamp: `t` missing where required, undecodable, or out
	// of the [MinDateTime, MaxDateTime] window.
	ReasonBadTimestamp
	// ReasonTypeMismatch: the value cannot coerce to the declared
	// ValueType under the docs/DESIGN.md §2.6 step 5 rules.
	ReasonTypeMismatch
	// ReasonValueTooLarge: a string exceeds MaxStringLen or an array
	// exceeds MaxArrayLen.
	ReasonValueTooLarge
	// ReasonBadObject: object-aggregation shape violation (not a
	// document, empty, or a key that resolves to no declared leaf).
	ReasonBadObject
	// ReasonUnsetNotAllowed: empty payload on a mapping without
	// allow_unset (datastreams never allow it).
	ReasonUnsetNotAllowed
)

// rejectReasonLabels maps RejectReason (offset by 1) to its metrics label.
var rejectReasonLabels = [...]string{
	"too_large", "unknown_format", "malformed", "no_value", "bad_timestamp",
	"type_mismatch", "value_too_large", "bad_object", "unset_not_allowed",
}

// String returns the stable snake_case label used for metrics and logs.
func (r RejectReason) String() string {
	if r >= ReasonTooLarge && r <= ReasonUnsetNotAllowed {
		return rejectReasonLabels[r-1]
	}
	return fmt.Sprintf("RejectReason(%d)", uint8(r))
}

// RejectReasons returns every reject reason, in order, so metric consumers
// can pre-register one counter per label.
func RejectReasons() []RejectReason {
	out := make([]RejectReason, 0, len(rejectReasonLabels))
	for i := range rejectReasonLabels {
		out = append(out, RejectReason(i+1)) //nolint:gosec // i < 9, no overflow
	}
	return out
}

// RejectError is the typed error returned for every payload rejection.
type RejectError struct {
	// Reason is the rejection class (metrics label).
	Reason RejectReason
	// Detail is a human-readable explanation for logs and trigger events.
	Detail string
}

// Error implements the error interface.
func (e *RejectError) Error() string {
	return fmt.Sprintf("payload rejected (%s): %s", e.Reason, e.Detail)
}

// rejectf builds a *RejectError with a formatted detail message.
func rejectf(reason RejectReason, format string, args ...any) error {
	return &RejectError{Reason: reason, Detail: fmt.Sprintf(format, args...)}
}

// ReasonOf extracts the RejectReason from err, or ReasonNone if err is nil
// or not a *RejectError.
func ReasonOf(err error) RejectReason {
	var re *RejectError
	if errors.As(err, &re) {
		return re.Reason
	}
	return ReasonNone
}

// ---------------------------------------------------------------------------
// Coercion table (docs/DESIGN.md §2.6 step 5). Each helper returns a
// *RejectError on failure so decoders can pass errors through unchanged.
// ---------------------------------------------------------------------------

// maxInt64Float is 2^63 as a float64 (exactly representable); float values
// in [-2^63, 2^63) convert to int64 without implementation-defined results.
const maxInt64Float = 9223372036854775808.0

// checkDouble rejects NaN and ±Inf, which neither Astarte's BSON profile
// nor JSON can represent.
func checkDouble(f float64) (float64, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, rejectf(ReasonTypeMismatch, "double value is %v", f)
	}
	return f, nil
}

// doubleFromInt64 widens an integer to double only when lossless
// (round-trips through float64 exactly; always true for |i| <= 2^53).
func doubleFromInt64(i int64) (float64, error) {
	f := float64(i)
	if f < -maxInt64Float || f >= maxInt64Float || int64(f) != i {
		return 0, rejectf(ReasonTypeMismatch, "integer %d does not widen losslessly to double", i)
	}
	return f, nil
}

// int32FromInt64 narrows an int64 to int32 only when it fits exactly.
func int32FromInt64(i int64) (int32, error) {
	if i < math.MinInt32 || i > math.MaxInt32 {
		return 0, rejectf(ReasonTypeMismatch, "integer %d overflows int32", i)
	}
	return int32(i), nil
}

// int32FromFloat64 converts a double to int32 only when it is integral and
// in range (exact fit; 5.0 coerces, 5.5 does not).
func int32FromFloat64(f float64) (int32, error) {
	if f != math.Trunc(f) || f < math.MinInt32 || f > math.MaxInt32 {
		return 0, rejectf(ReasonTypeMismatch, "number %v does not fit int32 exactly", f)
	}
	return int32(f), nil
}

// int64FromFloat64 converts a double to int64 only when it is integral and
// in range (exact fit — integral float64 values in range are exact).
func int64FromFloat64(f float64) (int64, error) {
	if f != math.Trunc(f) || f < -maxInt64Float || f >= maxInt64Float {
		return 0, rejectf(ReasonTypeMismatch, "number %v does not fit int64 exactly", f)
	}
	return int64(f), nil
}

// checkString enforces valid UTF-8 and the MaxStringLen cap.
func checkString(s string) (string, error) {
	if len(s) > MaxStringLen {
		return "", rejectf(ReasonValueTooLarge, "string of %d bytes exceeds %d byte cap", len(s), MaxStringLen)
	}
	if !utf8.ValidString(s) {
		return "", rejectf(ReasonTypeMismatch, "string is not valid UTF-8")
	}
	return s, nil
}

// dateTimeFromMillis converts an epoch-milliseconds timestamp, enforcing
// the [MinDateTime, MaxDateTime] window.
func dateTimeFromMillis(ms int64) (time.Time, error) {
	if ms < minDateTimeMillis || ms > maxDateTimeMillis {
		return time.Time{}, rejectf(ReasonBadTimestamp, "epoch-ms timestamp %d outside [%s, %s]",
			ms, MinDateTime.Format(time.RFC3339), MaxDateTime.Format(time.RFC3339))
	}
	return time.UnixMilli(ms).UTC(), nil
}

// checkDateTime bounds-checks an already-parsed instant and truncates it to
// millisecond precision (the Astarte datetime resolution; BSON UTC datetime
// carries no finer grain, and JSON inputs are normalised to match).
func checkDateTime(ts time.Time) (time.Time, error) {
	if ts.Before(MinDateTime) || ts.After(MaxDateTime) {
		return time.Time{}, rejectf(ReasonBadTimestamp, "timestamp %s outside [%s, %s]",
			ts.Format(time.RFC3339Nano), MinDateTime.Format(time.RFC3339), MaxDateTime.Format(time.RFC3339))
	}
	return ts.UTC().Truncate(time.Millisecond), nil
}

// collectArray builds a homogeneous []T of n elements, fetching each via
// elem (which performs the per-element coercion). Shared by the BSON and
// JSON array decoders.
func collectArray[T any](n int, elem func(i int) (Value, error)) ([]T, error) {
	if n > MaxArrayLen {
		return nil, rejectf(ReasonValueTooLarge, "array of %d elements exceeds %d element cap", n, MaxArrayLen)
	}
	out := make([]T, n)
	for i := range n {
		v, err := elem(i)
		if err != nil {
			return nil, annotate(err, "element %d", i)
		}
		out[i] = v.(T)
	}
	return out, nil
}

// annotate prefixes the Detail of a *RejectError with positional context
// ("element 3", `key "lat"`), preserving the reason. Non-RejectError errors
// pass through untouched.
func annotate(err error, format string, args ...any) error {
	var re *RejectError
	if errors.As(err, &re) {
		return &RejectError{Reason: re.Reason, Detail: fmt.Sprintf(format, args...) + ": " + re.Detail}
	}
	return err
}

// asValue adapts a typed coercion result (T, error) to the (Value, error)
// shape shared by the per-format decoders.
func asValue[T any](v T, err error) (Value, error) {
	if err != nil {
		return nil, err
	}
	return v, nil
}

// decodeArray dispatches an n-element array decode to the typed collector
// for the array value type vt, fetching each element via elem. The element
// callback must return the Go scalar type matching vt.Elem().
func decodeArray(vt interfaceschema.ValueType, n int, elem func(i int) (Value, error)) (Value, error) {
	switch vt {
	case interfaceschema.DoubleArray:
		return asValue(collectArray[float64](n, elem))
	case interfaceschema.IntegerArray:
		return asValue(collectArray[int32](n, elem))
	case interfaceschema.BooleanArray:
		return asValue(collectArray[bool](n, elem))
	case interfaceschema.LongIntegerArray:
		return asValue(collectArray[int64](n, elem))
	case interfaceschema.StringArray:
		return asValue(collectArray[string](n, elem))
	case interfaceschema.BinaryBlobArray:
		return asValue(collectArray[[]byte](n, elem))
	case interfaceschema.DateTimeArray:
		return asValue(collectArray[time.Time](n, elem))
	default:
		return nil, fmt.Errorf("payload: decodeArray: %s is not an array type", vt)
	}
}
