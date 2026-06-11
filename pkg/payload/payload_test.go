package payload

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// captureTime mirrors tools/bsoncapture: the `t` carried by every *_t.hex
// golden vector.
var captureTime = time.Date(2026, 6, 10, 12, 34, 56, 789000000, time.UTC)

// readVector loads one golden vector captured from the official Go SDK
// encoder (testdata/bson/README.md).
func readVector(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "bson", name+".hex"))
	if err != nil {
		t.Fatalf("reading golden vector: %v", err)
	}
	raw, err := hex.DecodeString(strings.TrimSpace(string(b)))
	if err != nil {
		t.Fatalf("golden vector %s is not valid hex: %v", name, err)
	}
	return raw
}

// valuesEqual compares decoded Values semantically (time.Time via Equal,
// byte slices via bytes.Equal, objects recursively).
func valuesEqual(a, b Value) bool {
	switch av := a.(type) {
	case time.Time:
		bv, ok := b.(time.Time)
		return ok && av.Equal(bv)
	case []byte:
		bv, ok := b.([]byte)
		return ok && bytes.Equal(av, bv)
	case [][]byte:
		bv, ok := b.([][]byte)
		if !ok || len(av) != len(bv) {
			return false
		}
		for i := range av {
			if !bytes.Equal(av[i], bv[i]) {
				return false
			}
		}
		return true
	case []time.Time:
		bv, ok := b.([]time.Time)
		if !ok || len(av) != len(bv) {
			return false
		}
		for i := range av {
			if !av[i].Equal(bv[i]) {
				return false
			}
		}
		return true
	case map[string]Value:
		bv, ok := b.(map[string]Value)
		if !ok || len(av) != len(bv) {
			return false
		}
		for k, v := range av {
			ev, ok := bv[k]
			if !ok || !valuesEqual(v, ev) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(a, b)
	}
}

// mapping builds a minimal CompiledMapping for decode tests.
func mapping(vt interfaceschema.ValueType, explicit bool) *interfaceschema.CompiledMapping {
	return &interfaceschema.CompiledMapping{ValueType: vt, ExplicitTimestamp: explicit}
}

// objectLeaves matches the golden `object` vector and the JSON object tests.
func objectLeaves(explicit bool) map[string]*interfaceschema.CompiledMapping {
	return map[string]*interfaceschema.CompiledMapping{
		"lat":     mapping(interfaceschema.Double, explicit),
		"lon":     mapping(interfaceschema.Double, explicit),
		"samples": mapping(interfaceschema.Integer, explicit),
		"ok":      mapping(interfaceschema.Boolean, explicit),
	}
}

// goldenScalarVectors maps each golden vector base name to its declared type
// and expected decoded value. Must stay in sync with tools/bsoncapture.
var goldenScalarVectors = []struct {
	name string
	vt   interfaceschema.ValueType
	want Value
}{
	{"double", interfaceschema.Double, 22.5},
	{"integer", interfaceschema.Integer, int32(42)},
	{"boolean", interfaceschema.Boolean, true},
	{"longinteger", interfaceschema.LongInteger, int64(9007199254740993)},
	{"string", interfaceschema.String, "héllo, Astarte ✓"},
	{"binaryblob", interfaceschema.BinaryBlob, []byte{0x00, 0x01, 0x02, 0xFE, 0xFF}},
	{"datetime", interfaceschema.DateTime, time.Date(2025, 12, 31, 23, 59, 59, 999000000, time.UTC)},
	{"doublearray", interfaceschema.DoubleArray, []float64{1.5, -2.25, 0}},
	{"integerarray", interfaceschema.IntegerArray, []int32{-1, 0, 2147483647}},
	{"booleanarray", interfaceschema.BooleanArray, []bool{true, false, true}},
	{"longintegerarray", interfaceschema.LongIntegerArray, []int64{1, -9007199254740993, 9223372036854775807}},
	{"stringarray", interfaceschema.StringArray, []string{"a", "β", ""}},
	{"binaryblobarray", interfaceschema.BinaryBlobArray, [][]byte{{0xDE, 0xAD}, {0xBE, 0xEF}, {}}},
	{"datetimearray", interfaceschema.DateTimeArray, []time.Time{
		time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 6, 10, 12, 0, 0, 500000000, time.UTC),
	}},
	{"doublearray_empty", interfaceschema.DoubleArray, []float64{}},
}

// TestGoldenVectorsDecode pins the ROADMAP §2.4 gate: every golden vector
// captured from the official SDK decodes to the expected typed Value, and
// re-encoding produces a semantically equal document with deterministic
// v-then-t order. For scalar/array vectors the re-encoded `v` and `t`
// elements must be byte-identical to the SDK's (same BSON element types);
// the object vector is compared semantically because the SDK encodes map
// keys in randomized order.
func TestGoldenVectorsDecode(t *testing.T) {
	for _, tc := range goldenScalarVectors {
		for _, withT := range []bool{false, true} {
			name := tc.name
			var wantTS *time.Time
			if withT {
				name += "_t"
				wantTS = &captureTime
			}
			t.Run(name, func(t *testing.T) {
				raw := readVector(t, name)
				dec, err := Decode(raw, mapping(tc.vt, withT))
				if err != nil {
					t.Fatalf("Decode: %v", err)
				}
				if dec.Format != FormatBSON {
					t.Errorf("Format = %v; want bson", dec.Format)
				}
				if !valuesEqual(dec.Value, tc.want) {
					t.Errorf("Value = %#v; want %#v", dec.Value, tc.want)
				}
				checkTimestamp(t, dec.Timestamp, wantTS)
				reencodeAndCompare(t, raw, dec, mapping(tc.vt, withT), false)
			})
		}
	}

	wantObject := map[string]Value{
		"lat": 45.4642, "lon": 9.19, "samples": int32(3), "ok": true,
	}
	for _, withT := range []bool{false, true} {
		name := "object"
		var wantTS *time.Time
		if withT {
			name += "_t"
			wantTS = &captureTime
		}
		t.Run(name, func(t *testing.T) {
			raw := readVector(t, name)
			dec, err := DecodeObject(raw, objectLeaves(withT))
			if err != nil {
				t.Fatalf("DecodeObject: %v", err)
			}
			if !valuesEqual(dec.Value, wantObject) {
				t.Errorf("Value = %#v; want %#v", dec.Value, wantObject)
			}
			checkTimestamp(t, dec.Timestamp, wantTS)
			reencodeAndCompare(t, raw, dec, nil, true)
		})
	}
}

// checkTimestamp asserts the decoded explicit timestamp matches want.
func checkTimestamp(t *testing.T, got, want *time.Time) {
	t.Helper()
	switch {
	case want == nil && got != nil:
		t.Errorf("Timestamp = %v; want nil", got)
	case want != nil && got == nil:
		t.Errorf("Timestamp = nil; want %v", want)
	case want != nil && !got.Equal(*want):
		t.Errorf("Timestamp = %v; want %v", got, want)
	}
}

// reencodeAndCompare re-encodes a decoded golden payload and asserts:
// deterministic v(,t) element order; per-element byte equality with the SDK
// capture (skipped for objects, whose SDK key order is randomized); and a
// lossless second decode.
func reencodeAndCompare(t *testing.T, golden []byte, dec DecodedPayload, m *interfaceschema.CompiledMapping, object bool) {
	t.Helper()
	enc, err := encodeBSON(dec.Value, dec.Timestamp)
	if err != nil {
		t.Fatalf("re-encode: %v", err)
	}

	els, err := bson.Raw(enc).Elements()
	if err != nil {
		t.Fatalf("re-encoded document invalid: %v", err)
	}
	wantElems := 1
	if dec.Timestamp != nil {
		wantElems = 2
	}
	if len(els) != wantElems || els[0].Key() != "v" {
		t.Fatalf("re-encoded element layout %v; want [v]/[v t]", els)
	}
	if dec.Timestamp != nil && els[1].Key() != "t" {
		t.Fatalf("re-encoded second element %q; want t", els[1].Key())
	}

	if !object {
		gv, err := bson.Raw(golden).LookupErr("v")
		if err != nil {
			t.Fatalf("golden v lookup: %v", err)
		}
		ev, err := bson.Raw(enc).LookupErr("v")
		if err != nil {
			t.Fatalf("re-encoded v lookup: %v", err)
		}
		if !ev.Equal(gv) {
			t.Errorf("re-encoded v element %s differs from SDK capture %s", ev, gv)
		}
		if dec.Timestamp != nil {
			gt, _ := bson.Raw(golden).LookupErr("t")
			et, _ := bson.Raw(enc).LookupErr("t")
			if !et.Equal(gt) {
				t.Errorf("re-encoded t element %s differs from SDK capture %s", et, gt)
			}
		}
	}

	// Second decode must be lossless.
	var redec DecodedPayload
	if object {
		redec, err = DecodeObject(enc, objectLeaves(dec.Timestamp != nil))
	} else {
		redec, err = Decode(enc, m)
	}
	if err != nil {
		t.Fatalf("decode of re-encoded payload: %v", err)
	}
	if !valuesEqual(redec.Value, dec.Value) {
		t.Errorf("re-encode round trip: %#v != %#v", redec.Value, dec.Value)
	}
	checkTimestamp(t, redec.Timestamp, dec.Timestamp)
}

// TestJSONProfile is the §3.5.3 profile table: every value type with valid
// and invalid inputs, both `t` forms, envelope strictness, and the bare-value
// shorthand rejection (ROADMAP §2.4 gate).
func TestJSONProfile(t *testing.T) {
	ts := time.Date(2026, 6, 10, 12, 34, 56, 789000000, time.UTC)
	cases := []struct {
		name       string
		in         string
		vt         interfaceschema.ValueType
		explicit   bool
		want       Value
		wantTS     *time.Time
		wantReason RejectReason
	}{
		// double
		{name: "double", in: `{"v":22.5}`, vt: interfaceschema.Double, want: 22.5},
		{name: "double from int token", in: `{"v":3}`, vt: interfaceschema.Double, want: 3.0},
		{name: "double leading ws", in: "  \t\n" + `{"v":-0.25}`, vt: interfaceschema.Double, want: -0.25},
		{name: "double extra member ignored", in: `{"v":1.5,"x":"ignored"}`, vt: interfaceschema.Double, want: 1.5},
		{name: "double from string", in: `{"v":"22.5"}`, vt: interfaceschema.Double, wantReason: ReasonTypeMismatch},
		{name: "double overflow", in: `{"v":1e400}`, vt: interfaceschema.Double, wantReason: ReasonTypeMismatch},
		{name: "double from bool", in: `{"v":true}`, vt: interfaceschema.Double, wantReason: ReasonTypeMismatch},
		{name: "double from array", in: `{"v":[22.5]}`, vt: interfaceschema.Double, wantReason: ReasonTypeMismatch},
		// integer
		{name: "integer", in: `{"v":42}`, vt: interfaceschema.Integer, want: int32(42)},
		{name: "integer exact-fit float", in: `{"v":5.0}`, vt: interfaceschema.Integer, want: int32(5)},
		{name: "integer min", in: `{"v":-2147483648}`, vt: interfaceschema.Integer, want: int32(math.MinInt32)},
		{name: "integer fractional", in: `{"v":5.5}`, vt: interfaceschema.Integer, wantReason: ReasonTypeMismatch},
		{name: "integer overflow", in: `{"v":2147483648}`, vt: interfaceschema.Integer, wantReason: ReasonTypeMismatch},
		{name: "integer from bool", in: `{"v":true}`, vt: interfaceschema.Integer, wantReason: ReasonTypeMismatch},
		// boolean
		{name: "boolean", in: `{"v":true}`, vt: interfaceschema.Boolean, want: true},
		{name: "boolean false", in: `{"v":false}`, vt: interfaceschema.Boolean, want: false},
		{name: "boolean from number", in: `{"v":1}`, vt: interfaceschema.Boolean, wantReason: ReasonTypeMismatch},
		{name: "boolean from string", in: `{"v":"true"}`, vt: interfaceschema.Boolean, wantReason: ReasonTypeMismatch},
		// longinteger
		{name: "longinteger number beyond 2^53", in: `{"v":9007199254740993}`, vt: interfaceschema.LongInteger, want: int64(9007199254740993)},
		{name: "longinteger decimal string", in: `{"v":"9223372036854775807"}`, vt: interfaceschema.LongInteger, want: int64(math.MaxInt64)},
		{name: "longinteger negative string", in: `{"v":"-9223372036854775808"}`, vt: interfaceschema.LongInteger, want: int64(math.MinInt64)},
		{name: "longinteger exact-fit float", in: `{"v":1e3}`, vt: interfaceschema.LongInteger, want: int64(1000)},
		{name: "longinteger fractional", in: `{"v":1.25}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
		{name: "longinteger bad string", in: `{"v":"12abc"}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
		{name: "longinteger number overflow", in: `{"v":9223372036854775808}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
		{name: "longinteger string overflow", in: `{"v":"9223372036854775808"}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
		{name: "longinteger float form beyond 2^53", in: `{"v":9007199254740993.0}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
		// string
		{name: "string", in: `{"v":"héllo"}`, vt: interfaceschema.String, want: "héllo"},
		{name: "string empty", in: `{"v":""}`, vt: interfaceschema.String, want: ""},
		{name: "string from number", in: `{"v":5}`, vt: interfaceschema.String, wantReason: ReasonTypeMismatch},
		// binaryblob
		{name: "binaryblob", in: `{"v":"3q2+7w=="}`, vt: interfaceschema.BinaryBlob, want: []byte{0xDE, 0xAD, 0xBE, 0xEF}},
		{name: "binaryblob unpadded", in: `{"v":"3q2+7w"}`, vt: interfaceschema.BinaryBlob, wantReason: ReasonTypeMismatch},
		{name: "binaryblob url alphabet", in: `{"v":"3q2-7w=="}`, vt: interfaceschema.BinaryBlob, wantReason: ReasonTypeMismatch},
		{name: "binaryblob from number", in: `{"v":5}`, vt: interfaceschema.BinaryBlob, wantReason: ReasonTypeMismatch},
		// datetime
		{name: "datetime rfc3339", in: `{"v":"2026-06-10T12:34:56.789Z"}`, vt: interfaceschema.DateTime, want: ts},
		{name: "datetime rfc3339 no fraction", in: `{"v":"2026-06-10T12:34:56Z"}`, vt: interfaceschema.DateTime, want: ts.Truncate(time.Second)},
		{name: "datetime rfc3339 offset", in: `{"v":"2026-06-10T14:34:56.789+02:00"}`, vt: interfaceschema.DateTime, want: ts},
		{name: "datetime epoch ms", in: `{"v":1781094896789}`, vt: interfaceschema.DateTime, want: time.UnixMilli(1781094896789).UTC()},
		{name: "datetime bad string", in: `{"v":"10/06/2026"}`, vt: interfaceschema.DateTime, wantReason: ReasonTypeMismatch},
		{name: "datetime from bool", in: `{"v":true}`, vt: interfaceschema.DateTime, wantReason: ReasonTypeMismatch},
		{name: "datetime epoch below window", in: `{"v":-99999999999999999}`, vt: interfaceschema.DateTime, wantReason: ReasonBadTimestamp},
		{name: "datetime epoch above window", in: `{"v":253402300800001}`, vt: interfaceschema.DateTime, wantReason: ReasonBadTimestamp},
		// arrays
		{name: "doublearray", in: `{"v":[1.5,2,3.25]}`, vt: interfaceschema.DoubleArray, want: []float64{1.5, 2, 3.25}},
		{name: "doublearray empty", in: `{"v":[]}`, vt: interfaceschema.DoubleArray, want: []float64{}},
		{name: "doublearray heterogeneous", in: `{"v":[1.5,"x"]}`, vt: interfaceschema.DoubleArray, wantReason: ReasonTypeMismatch},
		{name: "doublearray from scalar", in: `{"v":5}`, vt: interfaceschema.DoubleArray, wantReason: ReasonTypeMismatch},
		{name: "integerarray", in: `{"v":[1,-2]}`, vt: interfaceschema.IntegerArray, want: []int32{1, -2}},
		{name: "integerarray fractional element", in: `{"v":[1.5]}`, vt: interfaceschema.IntegerArray, wantReason: ReasonTypeMismatch},
		{name: "booleanarray", in: `{"v":[true,false]}`, vt: interfaceschema.BooleanArray, want: []bool{true, false}},
		{name: "booleanarray null element", in: `{"v":[true,null]}`, vt: interfaceschema.BooleanArray, wantReason: ReasonTypeMismatch},
		{name: "longintegerarray mixed forms", in: `{"v":[1,"9007199254740993"]}`, vt: interfaceschema.LongIntegerArray, want: []int64{1, 9007199254740993}},
		{name: "longintegerarray fractional", in: `{"v":[1.5]}`, vt: interfaceschema.LongIntegerArray, wantReason: ReasonTypeMismatch},
		{name: "stringarray", in: `{"v":["a",""]}`, vt: interfaceschema.StringArray, want: []string{"a", ""}},
		{name: "stringarray bad element", in: `{"v":[5]}`, vt: interfaceschema.StringArray, wantReason: ReasonTypeMismatch},
		{name: "binaryblobarray", in: `{"v":["3q0=","vu8="]}`, vt: interfaceschema.BinaryBlobArray, want: [][]byte{{0xDE, 0xAD}, {0xBE, 0xEF}}},
		{name: "binaryblobarray unpadded element", in: `{"v":["3q2+7w"]}`, vt: interfaceschema.BinaryBlobArray, wantReason: ReasonTypeMismatch},
		{name: "datetimearray mixed forms", in: `{"v":["2026-06-10T12:34:56.789Z",1781094896789]}`, vt: interfaceschema.DateTimeArray, want: []time.Time{ts, time.UnixMilli(1781094896789).UTC()}},
		{name: "datetimearray bad element", in: `{"v":["nope"]}`, vt: interfaceschema.DateTimeArray, wantReason: ReasonTypeMismatch},
		// timestamps
		{name: "t rfc3339", in: `{"v":22.5,"t":"2026-06-10T12:34:56.789Z"}`, vt: interfaceschema.Double, explicit: true, want: 22.5, wantTS: &ts},
		{name: "t epoch ms", in: `{"v":22.5,"t":1781094896789}`, vt: interfaceschema.Double, explicit: true, want: 22.5, wantTS: &ts},
		{name: "t ignored without explicit", in: `{"v":22.5,"t":"2026-06-10T12:34:56.789Z"}`, vt: interfaceschema.Double, want: 22.5},
		{name: "t required but missing", in: `{"v":22.5}`, vt: interfaceschema.Double, explicit: true, wantReason: ReasonBadTimestamp},
		{name: "t bad string", in: `{"v":1,"t":"yesterday"}`, vt: interfaceschema.Double, explicit: true, wantReason: ReasonBadTimestamp},
		{name: "t bool", in: `{"v":1,"t":true}`, vt: interfaceschema.Double, explicit: true, wantReason: ReasonBadTimestamp},
		{name: "t null", in: `{"v":1,"t":null}`, vt: interfaceschema.Double, explicit: true, wantReason: ReasonBadTimestamp},
		{name: "t fractional epoch", in: `{"v":1,"t":1781094896789.5}`, vt: interfaceschema.Double, explicit: true, wantReason: ReasonBadTimestamp},
		// envelope strictness
		{name: "bare number shorthand", in: `22.5`, vt: interfaceschema.Double, wantReason: ReasonUnknownFormat},
		{name: "bare string", in: `"hi"`, vt: interfaceschema.String, wantReason: ReasonUnknownFormat},
		{name: "bare array", in: `[22.5]`, vt: interfaceschema.DoubleArray, wantReason: ReasonUnknownFormat},
		{name: "truncated document", in: `{"v":22.5`, vt: interfaceschema.Double, wantReason: ReasonMalformed},
		{name: "trailing garbage", in: `{"v":1} x`, vt: interfaceschema.Double, wantReason: ReasonMalformed},
		{name: "no v member", in: `{}`, vt: interfaceschema.Double, wantReason: ReasonNoValue},
		{name: "only t member", in: `{"t":"2026-06-10T12:34:56Z"}`, vt: interfaceschema.Double, wantReason: ReasonNoValue},
		{name: "null v", in: `{"v":null}`, vt: interfaceschema.Double, wantReason: ReasonTypeMismatch},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dec, err := Decode([]byte(tc.in), mapping(tc.vt, tc.explicit))
			if tc.wantReason != ReasonNone {
				if ReasonOf(err) != tc.wantReason {
					t.Fatalf("Decode(%s) err = %v (reason %v); want reason %v", tc.in, err, ReasonOf(err), tc.wantReason)
				}
				return
			}
			if err != nil {
				t.Fatalf("Decode(%s): %v", tc.in, err)
			}
			if dec.Format != FormatJSON {
				t.Errorf("Format = %v; want json", dec.Format)
			}
			if !valuesEqual(dec.Value, tc.want) {
				t.Errorf("Value = %#v; want %#v", dec.Value, tc.want)
			}
			checkTimestamp(t, dec.Timestamp, tc.wantTS)
		})
	}
}

// TestJSONObjectAggregation covers the object-aggregation JSON path.
func TestJSONObjectAggregation(t *testing.T) {
	leaves := objectLeaves(false)

	dec, err := DecodeObject([]byte(`{"v":{"lat":45.0,"lon":9.0,"samples":3,"ok":true}}`), leaves)
	if err != nil {
		t.Fatalf("DecodeObject: %v", err)
	}
	want := map[string]Value{"lat": 45.0, "lon": 9.0, "samples": int32(3), "ok": true}
	if !valuesEqual(dec.Value, want) {
		t.Errorf("Value = %#v; want %#v", dec.Value, want)
	}

	// Partial documents are accepted: only declared-leaf resolution is
	// enforced (docs/DESIGN.md §2.6 step 6).
	if _, err := DecodeObject([]byte(`{"v":{"lat":45.0}}`), leaves); err != nil {
		t.Errorf("partial object rejected: %v", err)
	}

	for in, wantReason := range map[string]RejectReason{
		`{"v":{"lat":1.0,"nope":2.0}}`: ReasonBadObject,    // undeclared key
		`{"v":{}}`:                     ReasonBadObject,    // empty document
		`{"v":5}`:                      ReasonBadObject,    // not a document
		`{"v":null}`:                   ReasonBadObject,    // null
		`{"v":[1]}`:                    ReasonBadObject,    // array
		`{"v":{"lat":"x"}}`:            ReasonTypeMismatch, // leaf type violation
	} {
		if _, err := DecodeObject([]byte(in), leaves); ReasonOf(err) != wantReason {
			t.Errorf("DecodeObject(%s) reason = %v; want %v", in, ReasonOf(err), wantReason)
		}
	}

	// Explicit-timestamp policy comes from the leaves.
	if _, err := DecodeObject([]byte(`{"v":{"lat":1.0}}`), objectLeaves(true)); ReasonOf(err) != ReasonBadTimestamp {
		t.Errorf("explicit object without t: reason = %v; want bad_timestamp", ReasonOf(err))
	}
}

// rawBSONDoc hand-crafts a BSON document from pre-encoded element bodies so
// tests can produce documents the package encoder refuses to emit.
func rawBSONDoc(t *testing.T, elems ...[]byte) []byte {
	t.Helper()
	doc := make([]byte, 4)
	for _, e := range elems {
		doc = append(doc, e...)
	}
	doc = append(doc, 0x00)
	binary.LittleEndian.PutUint32(doc[0:4], uint32(len(doc)))
	return doc
}

// rawBSONElem hand-crafts one element: type byte, key, NUL, body.
func rawBSONElem(typ byte, key string, body []byte) []byte {
	e := append([]byte{typ}, key...)
	e = append(e, 0x00)
	return append(e, body...)
}

// TestBSONCoercion covers the §2.6 step 5 BSON coercion edges that golden
// vectors cannot (cross-type widening/narrowing, NaN, malformed documents).
func TestBSONCoercion(t *testing.T) {
	enc := func(v Value) []byte {
		t.Helper()
		b, err := encodeBSON(v, nil)
		if err != nil {
			t.Fatalf("encodeBSON(%v): %v", v, err)
		}
		return b
	}
	float64Body := func(f float64) []byte {
		return binary.LittleEndian.AppendUint64(nil, math.Float64bits(f))
	}

	cases := []struct {
		name       string
		in         []byte
		vt         interfaceschema.ValueType
		want       Value
		wantReason RejectReason
	}{
		{name: "int64 widens to double", in: enc(int64(1 << 53)), vt: interfaceschema.Double, want: float64(1 << 53)},
		{name: "int32 widens to double", in: enc(int32(-7)), vt: interfaceschema.Double, want: -7.0},
		{name: "lossy int64 to double", in: enc(int64(1<<53 + 1)), vt: interfaceschema.Double, wantReason: ReasonTypeMismatch},
		{name: "double 5.0 to integer", in: enc(5.0), vt: interfaceschema.Integer, want: int32(5)},
		{name: "double 5.5 to integer", in: enc(5.5), vt: interfaceschema.Integer, wantReason: ReasonTypeMismatch},
		{name: "int64 fits integer", in: enc(int64(123)), vt: interfaceschema.Integer, want: int32(123)},
		{name: "int64 overflows integer", in: enc(int64(math.MaxInt32 + 1)), vt: interfaceschema.Integer, wantReason: ReasonTypeMismatch},
		{name: "int32 widens to longinteger", in: enc(int32(9)), vt: interfaceschema.LongInteger, want: int64(9)},
		{name: "double never longinteger", in: enc(5.0), vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
		{name: "string is not integer", in: enc("5"), vt: interfaceschema.Integer, wantReason: ReasonTypeMismatch},
		{name: "boolean is not integer", in: enc(true), vt: interfaceschema.Integer, wantReason: ReasonTypeMismatch},
		{
			name:       "NaN double rejected",
			in:         rawBSONDoc(t, rawBSONElem(byte(bson.TypeDouble), "v", float64Body(math.NaN()))),
			vt:         interfaceschema.Double,
			wantReason: ReasonTypeMismatch,
		},
		{
			name:       "+Inf double rejected",
			in:         rawBSONDoc(t, rawBSONElem(byte(bson.TypeDouble), "v", float64Body(math.Inf(1)))),
			vt:         interfaceschema.Double,
			wantReason: ReasonTypeMismatch,
		},
		{
			name: "datetime below window",
			in: rawBSONDoc(t, rawBSONElem(byte(bson.TypeDateTime), "v",
				binary.LittleEndian.AppendUint64(nil, uint64(minDateTimeMillis-1)))),
			vt:         interfaceschema.DateTime,
			wantReason: ReasonBadTimestamp,
		},
		{
			name:       "null v rejected",
			in:         rawBSONDoc(t, rawBSONElem(byte(bson.TypeNull), "v", nil)),
			vt:         interfaceschema.Double,
			wantReason: ReasonTypeMismatch,
		},
		{
			name:       "missing v",
			in:         rawBSONDoc(t, rawBSONElem(byte(bson.TypeBoolean), "x", []byte{1})),
			vt:         interfaceschema.Boolean,
			wantReason: ReasonNoValue,
		},
		{
			name: "t wrong type",
			in: rawBSONDoc(t,
				rawBSONElem(byte(bson.TypeBoolean), "v", []byte{1}),
				rawBSONElem(byte(bson.TypeString), "t", append(binary.LittleEndian.AppendUint32(nil, 4), 'n', 'o', 'w', 0x00))),
			vt:         interfaceschema.Boolean,
			wantReason: ReasonBadTimestamp,
		},
		{
			name:       "corrupt interior",
			in:         append(binary.LittleEndian.AppendUint32(nil, 9), 0xFF, 0xFF, 0xFF, 0xFF, 0x00),
			vt:         interfaceschema.Double,
			wantReason: ReasonMalformed,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dec, err := Decode(tc.in, mapping(tc.vt, false))
			if tc.wantReason != ReasonNone {
				if ReasonOf(err) != tc.wantReason {
					t.Fatalf("reason = %v (err %v); want %v", ReasonOf(err), err, tc.wantReason)
				}
				return
			}
			if err != nil {
				t.Fatalf("Decode: %v", err)
			}
			if !valuesEqual(dec.Value, tc.want) {
				t.Errorf("Value = %#v; want %#v", dec.Value, tc.want)
			}
		})
	}

	// Duplicate object keys are rejected.
	dup := rawBSONDoc(t, rawBSONElem(byte(bson.TypeEmbeddedDocument), "v", rawBSONDoc(t,
		rawBSONElem(byte(bson.TypeBoolean), "ok", []byte{1}),
		rawBSONElem(byte(bson.TypeBoolean), "ok", []byte{0}),
	)))
	if _, err := DecodeObject(dup, objectLeaves(false)); ReasonOf(err) != ReasonBadObject {
		t.Errorf("duplicate object key reason = %v; want bad_object", ReasonOf(err))
	}
}

// TestSizeAndCardinalityCaps pins the 64 KiB payload cap and the 1024-element
// array cap for both formats (ROADMAP §2.4 gate).
func TestSizeAndCardinalityCaps(t *testing.T) {
	// Default 64 KiB cap, BSON: a >64 KiB binary payload.
	bigBlob, err := encodeBSON(make([]byte, DefaultMaxSize), nil)
	if err != nil {
		t.Fatalf("encodeBSON: %v", err)
	}
	if _, err := Decode(bigBlob, mapping(interfaceschema.BinaryBlob, false)); ReasonOf(err) != ReasonTooLarge {
		t.Errorf("BSON over default cap: reason = %v; want too_large", ReasonOf(err))
	}
	// Default 64 KiB cap, JSON.
	bigJSON := []byte(`{"v":"` + strings.Repeat("a", DefaultMaxSize) + `"}`)
	if _, err := Decode(bigJSON, mapping(interfaceschema.String, false)); ReasonOf(err) != ReasonTooLarge {
		t.Errorf("JSON over default cap: reason = %v; want too_large", ReasonOf(err))
	}
	// Custom tighter cap applies to both formats.
	small := Decoder{MaxSize: 8}
	if _, err := small.Individual([]byte(`{"v":true}`), mapping(interfaceschema.Boolean, false)); ReasonOf(err) != ReasonTooLarge {
		t.Errorf("custom cap JSON: want too_large")
	}
	okDoc, _ := encodeBSON(true, nil)
	if _, err := small.Individual(okDoc, mapping(interfaceschema.Boolean, false)); ReasonOf(err) != ReasonTooLarge {
		t.Errorf("custom cap BSON: want too_large")
	}
	// String cap distinct from the payload cap (raised payload limit).
	wide := Decoder{MaxSize: 1 << 20}
	overString := []byte(`{"v":"` + strings.Repeat("a", MaxStringLen+1) + `"}`)
	if _, err := wide.Individual(overString, mapping(interfaceschema.String, false)); ReasonOf(err) != ReasonValueTooLarge {
		t.Errorf("over-cap string: reason = %v; want value_too_large", ReasonOf(err))
	}

	// 1024-element array cap, JSON.
	overJSON := []byte(`{"v":[` + strings.Repeat("0,", MaxArrayLen) + `0]}`)
	if _, err := Decode(overJSON, mapping(interfaceschema.IntegerArray, false)); ReasonOf(err) != ReasonValueTooLarge {
		t.Errorf("JSON 1025-element array: reason = %v; want value_too_large", ReasonOf(err))
	}
	atJSON := []byte(`{"v":[` + strings.Repeat("0,", MaxArrayLen-1) + `0]}`)
	if _, err := Decode(atJSON, mapping(interfaceschema.IntegerArray, false)); err != nil {
		t.Errorf("JSON 1024-element array rejected: %v", err)
	}
	// 1024-element array cap, BSON (hand-crafted: the encoder enforces the
	// cap itself, so the oversize document must be built manually).
	var elems []byte
	for i := range MaxArrayLen + 1 {
		elems = append(elems, rawBSONElem(byte(bson.TypeInt32), strconv.Itoa(i), []byte{0, 0, 0, 0})...)
	}
	arrDoc := make([]byte, 4)
	arrDoc = append(arrDoc, elems...)
	arrDoc = append(arrDoc, 0x00)
	binary.LittleEndian.PutUint32(arrDoc[0:4], uint32(len(arrDoc)))
	overBSON := rawBSONDoc(t, rawBSONElem(byte(bson.TypeArray), "v", arrDoc))
	if _, err := Decode(overBSON, mapping(interfaceschema.IntegerArray, false)); ReasonOf(err) != ReasonValueTooLarge {
		t.Errorf("BSON 1025-element array: reason = %v; want value_too_large", ReasonOf(err))
	}
}

// TestFacadeUnsetAndErrors covers property-unset handling and caller-bug
// (non-Reject) errors.
func TestFacadeUnsetAndErrors(t *testing.T) {
	unsettable := &interfaceschema.CompiledMapping{ValueType: interfaceschema.Double, AllowUnset: true}
	dec, err := Decode(nil, unsettable)
	if err != nil {
		t.Fatalf("unset decode: %v", err)
	}
	if !dec.IsUnset() || dec.Format != FormatEmpty || dec.Timestamp != nil {
		t.Errorf("unset payload decoded to %#v", dec)
	}

	if _, err := Decode([]byte{}, mapping(interfaceschema.Double, false)); ReasonOf(err) != ReasonUnsetNotAllowed {
		t.Errorf("unset without allow_unset: reason = %v; want unset_not_allowed", ReasonOf(err))
	}
	if _, err := DecodeObject(nil, objectLeaves(false)); ReasonOf(err) != ReasonUnsetNotAllowed {
		t.Errorf("unset object: reason = %v; want unset_not_allowed", ReasonOf(err))
	}
	if _, err := Decode([]byte{0xDE, 0xAD, 0xBE, 0xEF, 0x99}, mapping(interfaceschema.Double, false)); ReasonOf(err) != ReasonUnknownFormat {
		t.Errorf("garbage payload: want unknown_format")
	}

	// Caller bugs are plain errors, not payload rejections.
	if _, err := Decode([]byte(`{"v":1}`), nil); err == nil || ReasonOf(err) != ReasonNone {
		t.Errorf("nil mapping: want plain error, got %v", err)
	}
	if _, err := DecodeObject([]byte(`{"v":{}}`), nil); err == nil || ReasonOf(err) != ReasonNone {
		t.Errorf("nil leaves: want plain error, got %v", err)
	}
}

// TestEncodeRoundTrip drives Encode through every type in both formats and
// decodes the result back.
func TestEncodeRoundTrip(t *testing.T) {
	ts := time.Date(2026, 6, 10, 12, 34, 56, 789000000, time.UTC)
	cases := []struct {
		vt interfaceschema.ValueType
		v  Value
	}{
		{interfaceschema.Double, 22.5},
		{interfaceschema.Integer, int32(-42)},
		{interfaceschema.Boolean, true},
		{interfaceschema.LongInteger, int64(9007199254740993)},
		{interfaceschema.LongInteger, int64(-12)},
		{interfaceschema.String, "héllo ✓"},
		{interfaceschema.BinaryBlob, []byte{0, 1, 0xFE}},
		{interfaceschema.DateTime, time.Date(1969, 7, 20, 20, 17, 40, 0, time.UTC)},
		{interfaceschema.DoubleArray, []float64{1.5, -2.5}},
		{interfaceschema.IntegerArray, []int32{1, 2}},
		{interfaceschema.BooleanArray, []bool{false, true}},
		{interfaceschema.LongIntegerArray, []int64{math.MinInt64, math.MaxInt64}},
		{interfaceschema.StringArray, []string{"", "x"}},
		{interfaceschema.BinaryBlobArray, [][]byte{{0xAA}, {}}},
		{interfaceschema.DateTimeArray, []time.Time{ts, ts.Add(time.Hour)}},
	}
	for _, f := range []Format{FormatBSON, FormatJSON} {
		for _, tc := range cases {
			for _, withT := range []bool{false, true} {
				var expTS *time.Time
				if withT {
					expTS = &ts
				}
				enc, err := Encode(tc.v, expTS, f)
				if err != nil {
					t.Fatalf("Encode(%v, %v, %v): %v", tc.v, expTS, f, err)
				}
				dec, err := Decode(enc, mapping(tc.vt, withT))
				if err != nil {
					t.Fatalf("Decode(Encode(%#v), %v, withT=%v): %v", tc.v, f, withT, err)
				}
				if dec.Format != f {
					t.Errorf("round-trip format = %v; want %v", dec.Format, f)
				}
				if !valuesEqual(dec.Value, tc.v) {
					t.Errorf("%v round trip: %#v != %#v", f, dec.Value, tc.v)
				}
				checkTimestamp(t, dec.Timestamp, expTS)
			}
		}
		// Object aggregation round trip.
		obj := map[string]Value{"lat": 1.25, "lon": -2.5, "samples": int32(7), "ok": false}
		enc, err := Encode(obj, &ts, f)
		if err != nil {
			t.Fatalf("Encode(object, %v): %v", f, err)
		}
		dec, err := DecodeObject(enc, objectLeaves(true))
		if err != nil {
			t.Fatalf("DecodeObject(%v): %v", f, err)
		}
		if !valuesEqual(dec.Value, obj) {
			t.Errorf("%v object round trip: %#v != %#v", f, dec.Value, obj)
		}
	}
}

// TestEncodeJSONShape freezes the exact JSON profile bytes for canonical
// samples: member order v,t; RFC 3339 ms timestamps; longinteger string form
// beyond 2^53.
func TestEncodeJSONShape(t *testing.T) {
	ts := time.Date(2026, 6, 10, 12, 34, 56, 789000000, time.UTC)
	cases := []struct {
		v    Value
		ts   *time.Time
		want string
	}{
		{22.5, &ts, `{"v":22.5,"t":"2026-06-10T12:34:56.789Z"}`},
		{22.5, nil, `{"v":22.5}`},
		{int64(12), nil, `{"v":12}`},
		{int64(9007199254740993), nil, `{"v":"9007199254740993"}`},
		{int64(-9007199254740993), nil, `{"v":"-9007199254740993"}`},
		{[]byte{0xDE, 0xAD, 0xBE, 0xEF}, nil, `{"v":"3q2+7w=="}`},
		{ts, nil, `{"v":"2026-06-10T12:34:56.789Z"}`},
		{[]int32{1, 2}, nil, `{"v":[1,2]}`},
	}
	for _, tc := range cases {
		got, err := Encode(tc.v, tc.ts, FormatJSON)
		if err != nil {
			t.Fatalf("Encode(%#v): %v", tc.v, err)
		}
		if string(got) != tc.want {
			t.Errorf("Encode(%#v) = %s; want %s", tc.v, got, tc.want)
		}
	}
}

// TestEncodeRejects pins the encoder-side validation and format handling.
func TestEncodeRejects(t *testing.T) {
	ts := time.Now()

	// FormatEmpty: the property-unset payload.
	empty, err := Encode(nil, nil, FormatEmpty)
	if err != nil || len(empty) != 0 {
		t.Errorf("Encode(nil, nil, empty) = %v, %v", empty, err)
	}
	if _, err := Encode(22.5, nil, FormatEmpty); err == nil {
		t.Error("FormatEmpty with a value accepted")
	}
	if _, err := Encode(nil, &ts, FormatEmpty); err == nil {
		t.Error("FormatEmpty with a timestamp accepted")
	}
	if _, err := Encode(nil, nil, FormatBSON); err == nil {
		t.Error("nil value accepted for BSON")
	}
	if _, err := Encode(22.5, nil, FormatInvalid); err == nil {
		t.Error("FormatInvalid accepted")
	}

	for _, f := range []Format{FormatBSON, FormatJSON} {
		if _, err := Encode(math.NaN(), nil, f); ReasonOf(err) != ReasonTypeMismatch {
			t.Errorf("%v: NaN reason = %v; want type_mismatch", f, ReasonOf(err))
		}
		if _, err := Encode(math.Inf(-1), nil, f); ReasonOf(err) != ReasonTypeMismatch {
			t.Errorf("%v: -Inf reason = %v; want type_mismatch", f, ReasonOf(err))
		}
		if _, err := Encode(uint(5), nil, f); ReasonOf(err) != ReasonTypeMismatch {
			t.Errorf("%v: unsupported Go type reason = %v; want type_mismatch", f, ReasonOf(err))
		}
		if _, err := Encode(strings.Repeat("a", MaxStringLen+1), nil, f); ReasonOf(err) != ReasonValueTooLarge {
			t.Errorf("%v: over-cap string reason = %v; want value_too_large", f, ReasonOf(err))
		}
		if _, err := Encode(make([]bool, MaxArrayLen+1), nil, f); ReasonOf(err) != ReasonValueTooLarge {
			t.Errorf("%v: over-cap array reason = %v; want value_too_large", f, ReasonOf(err))
		}
		if _, err := Encode(time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC), nil, f); ReasonOf(err) != ReasonBadTimestamp {
			t.Errorf("%v: out-of-window datetime reason = %v; want bad_timestamp", f, ReasonOf(err))
		}
		badTS := time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)
		if _, err := Encode(22.5, &badTS, f); ReasonOf(err) != ReasonBadTimestamp {
			t.Errorf("%v: out-of-window t reason = %v; want bad_timestamp", f, ReasonOf(err))
		}
		if _, err := Encode(map[string]Value{}, nil, f); ReasonOf(err) != ReasonBadObject {
			t.Errorf("%v: empty object reason = %v; want bad_object", f, ReasonOf(err))
		}
		if _, err := Encode(map[string]Value{"a": map[string]Value{"b": 1.0}}, nil, f); ReasonOf(err) != ReasonBadObject {
			t.Errorf("%v: nested object reason = %v; want bad_object", f, ReasonOf(err))
		}
	}
}
