slug: flow-msg-json-integer-precision
verdict: done
at:  4df3f6b
ran: 2026-09-23T19:07:10Z on DietPi in 612s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
✱ Glob "**/message.go" 1 match
✱ Glob "**/message_test.go" 1 match
→ Read internal/flow/message.go
→ Read internal/flow/message_test.go
✱ Grep "payload-longinteger-fraction-quantize" 34 matches
✱ Grep "UseNumber" 1 match
Let me look at the reference fix and its report.
→ Read pkg/payload/json.go [offset=230, limit=120]
→ Read internal/appengine/data.go [offset=330, limit=60]
✱ Grep "flow\.Message|flow\.UnmarshalMessage|Data\.\(map\[string\]any\)|Data\.\(float64\)" 100 matches
Now let me check who consumes these messages, to see the impact of switching the decoder to `UseNumber`.
✱ Grep "TypeReal|TypeInteger" in internal · 55 matches
→ Read internal/flow/blocks/mqtt.go [offset=270, limit=40]
→ Read internal/flow/blocks/splitmap.go
→ Read internal/flow/blocks/httpblocks.go [limit=70, offset=173]
→ Read internal/flow/blocks/jsonpathmap.go
✱ Grep "setDataFromWire|wireMessage\{|json\.Unmarshal\(.*&w" 10 matches
→ Read internal/flow/pump_test.go [limit=130]
✱ Grep "json\.Unmarshal|\.Marshal\(" in internal/flow · 29 matches
Let me check for any other consumers of the wire format across the repo, then look at the payload test rows for reference style.
✱ Grep "astarte_flow/message|WireSchema|TypeMap|flow\.Message\{" 100 matches
→ Read examples/flow-container-echo/main_test.go
→ Read internal/flow/blocks/container/httpbridge.go [limit=90, offset=100]
→ Read internal/flow/blocks/container/block_test.go [limit=190, offset=140]
Now I have a full picture. Let me implement the fix in message.go.

← Edit internal/flow/message.go
Index: /root/astrate-mule/internal/flow/message.go
===================================================================
--- /root/astrate-mule/internal/flow/message.go
+++ /root/astrate-mule/internal/flow/message.go
@@ -4,11 +4,13 @@
 // may interleave across lanes.
 package flow
 
 import (
+	"bytes"
 	"encoding/base64"
 	"encoding/json"
 	"fmt"
+	"math"
 	"time"
 )
 
 // WireSchema is the JSON schema identifier used on the wire.

← Edit internal/flow/message.go
Index: /root/astrate-mule/internal/flow/message.go
===================================================================
--- /root/astrate-mule/internal/flow/message.go
+++ /root/astrate-mule/internal/flow/message.go
@@ -141,9 +141,14 @@
 
 // UnmarshalJSON deserialises a Message from the upstream JSON wire format.
 func (m *Message) UnmarshalJSON(b []byte) error {
 	var w wireMessage
-	if err := json.Unmarshal(b, &w); err != nil {
+	// Decode with UseNumber so numeric tokens survive as json.Number instead
+	// of being pre-rounded to float64 (which would lose integer precision
+	// beyond 2^53 before setDataFromWire ever sees them).
+	dec := json.NewDecoder(bytes.NewReader(b))
+	dec.UseNumber()
+	if err := dec.Decode(&w); err != nil {
 		return fmt.Errorf("flow: unmarshal message: %w", err)
 	}
 
 	if w.Schema != WireSchema {

← Edit internal/flow/message.go
Index: /root/astrate-mule/internal/flow/message.go
===================================================================
--- /root/astrate-mule/internal/flow/message.go
+++ /root/astrate-mule/internal/flow/message.go
@@ -268,33 +268,19 @@
 // setDataFromWire decodes a scalar wire data value into m.Data.
 func (m *Message) setDataFromWire(dt DataType, raw any) error {
 	switch dt {
 	case TypeInteger:
-		switch v := raw.(type) {
-		case float64:
-			m.Data = int64(v)
-		case json.Number:
-			n, err := v.Int64()
-			if err != nil {
-				return fmt.Errorf("flow: integer data: %w", err)
-			}
-			m.Data = n
-		default:
-			return fmt.Errorf("flow: integer data: expected number, got %T", raw)
+		n, err := integerWireValue(raw)
+		if err != nil {
+			return fmt.Errorf("flow: integer data: %w", err)
 		}
+		m.Data = n
 	case TypeReal:
-		switch v := raw.(type) {
-		case float64:
-			m.Data = v
-		case json.Number:
-			f, err := v.Float64()
-			if err != nil {
-				return fmt.Errorf("flow: real data: %w", err)
-			}
-			m.Data = f
-		default:
-			return fmt.Errorf("flow: real data: expected number, got %T", raw)
+		f, err := realWireValue(raw)
+		if err != nil {
+			return fmt.Errorf("flow: real data: %w", err)
 		}
+		m.Data = f
 	case TypeBoolean:
 		b, ok := raw.(bool)
 		if !ok {
 			return fmt.Errorf("flow: boolean data: expected bool, got %T", raw)

Now update the map decoder to coerce numeric fields per their declared type, and add the helpers.
← Edit internal/flow/message.go
Index: /root/astrate-mule/internal/flow/message.go
===================================================================
--- /root/astrate-mule/internal/flow/message.go
+++ /root/astrate-mule/internal/flow/message.go
@@ -348,10 +348,57 @@
 				return fmt.Errorf("flow: map field %q: %w", k, err)
 			}
 			out[k] = t
 		default:
-			out[k] = v
+			switch dt {
+			case TypeInteger:
+				n, err := integerWireValue(v)
+				if err != nil {
+					return fmt.Errorf("flow: map field %q: %w", k, err)
+				}
+				out[k] = n
+			case TypeReal:
+				f, err := realWireValue(v)
+				if err != nil {
+					return fmt.Errorf("flow: map field %q: %w", k, err)
+				}
+				out[k] = f
+			default:
+				out[k] = v
+			}
 		}
 	}
 	m.Data = out
 	return nil
 }
+
+// integerWireValue coerces a wire value for an integer field to int64,
+// rejecting fraction and out-of-range forms that a float64 round-trip would
+// silently truncate or wrap (3.7 → 3, 1e300 → garbage). Wire decoding
+// delivers json.Number (the decoder runs with UseNumber), which is judged
+// exactly against the token; a bare float64 (programmatic data) is accepted
+// only when it is an exact in-range integer.
+func integerWireValue(v any) (int64, error) {
+	switch n := v.(type) {
+	case json.Number:
+		return n.Int64()
+	case float64:
+		if n != math.Trunc(n) || n < float64(math.MinInt64) || n >= float64(math.MaxInt64) {
+			return 0, fmt.Errorf("%v does not coerce to an integer exactly", n)
+		}
+		return int64(n), nil
+	default:
+		return 0, fmt.Errorf("expected number, got %T", v)
+	}
+}
+
+// realWireValue coerces a wire value for a real field to float64.
+func realWireValue(v any) (float64, error) {
+	switch n := v.(type) {
+	case json.Number:
+		return n.Float64()
+	case float64:
+		return n, nil
+	default:
+		return 0, fmt.Errorf("expected number, got %T", v)
+	}
+}

Now the tests. First add the `math` import and the `math.MaxInt64` round-trip row:
← Edit internal/flow/message_test.go
Index: /root/astrate-mule/internal/flow/message_test.go
===================================================================
--- /root/astrate-mule/internal/flow/message_test.go
+++ /root/astrate-mule/internal/flow/message_test.go
@@ -1,8 +1,9 @@
 package flow
 
 import (
 	"encoding/json"
+	"math"
 	"testing"
 	"time"
 )
 

← Edit internal/flow/message_test.go
Index: /root/astrate-mule/internal/flow/message_test.go
===================================================================
--- /root/astrate-mule/internal/flow/message_test.go
+++ /root/astrate-mule/internal/flow/message_test.go
@@ -13,8 +13,17 @@
 	msg     Message
 	wantErr bool
 }{
 	{
+		name: "integer max int64",
+		msg: Message{
+			Key:       "sensor/big",
+			Type:      TypeInteger,
+			Data:      int64(math.MaxInt64),
+			Timestamp: 1551884045074181,
+		},
+	},
+	{
 		name: "integer",
 		msg: Message{
 			Key:       "sensor/temp",
 			Type:      TypeInteger,

Now the rejection tests, added after `TestFlowMessage_UnmarshalRejectsUnknownSchema`:
← Edit internal/flow/message_test.go
Index: /root/astrate-mule/internal/flow/message_test.go
===================================================================
--- /root/astrate-mule/internal/flow/message_test.go
+++ /root/astrate-mule/internal/flow/message_test.go
@@ -239,8 +239,29 @@
 		t.Fatal("expected error for unknown schema")
 	}
 }
 
+func TestFlowMessage_IntegerWireRejectsFractionAndOverflow(t *testing.T) {
+	for _, data := range []string{"3.7", "1e300", "9223372036854775808"} {
+		t.Run(data, func(t *testing.T) {
+			raw := `{"schema":"` + WireSchema + `","key":"k","type":"integer","timestamp_us":0,"data":` + data + `}`
+			var msg Message
+			if err := json.Unmarshal([]byte(raw), &msg); err == nil {
+				t.Fatalf("unmarshal integer data %s: expected error, got %v", data, msg.Data)
+			}
+		})
+	}
+}
+
+func TestSetDataFromWire_IntegerRejectsFractionalFloat64(t *testing.T) {
+	for _, v := range []float64{3.7, 1e300} {
+		var msg Message
+		if err := msg.setDataFromWire(TypeInteger, v); err == nil {
+			t.Fatalf("setDataFromWire(%v): expected error, got %v", v, msg.Data)
+		}
+	}
+}
+
 func TestFlowMessage_MapSubtypesOmittedWhenEmpty(t *testing.T) {
 	msg := Message{
 		Key:       "m",
 		Type:      TypeMap,

Now run the gate. Start with vet, gofmt, and the changed package's tests, then the full suite.
$ gofmt -l internal/flow/message.go internal/flow/message_test.go && go vet ./internal/flow/ && go test ./internal/flow/
ok  	github.com/astrate-platform/astrate/internal/flow	0.510s

Gates pass with the change. Now let me verify the tests genuinely break against the old behavior — I'll replicate the pre-fix decode path in a `/tmp` probe (no git, no repo changes).
← Write /tmp/opencode/olddecode/main.go
Wrote file successfully.

$ cd /tmp/opencode/olddecode && go mod init probe >/dev/null 2>&1; go run .
1e300                          old decode -> 9223372036854775807
2^63                           old decode -> 9223372036854775807
maxint64 9223372036854775807   old decode -> 9223372036854775807
3.7                            old decode -> 3

The probe reveals the old float→int64 conversion **saturates to MaxInt64 on this machine**, not wrap to MinInt64 as the task assumed — so a plain `math.MaxInt64` round-trip row would pass even with the old code here. Let me check the platform and rework the test so it genuinely breaks pre-fix.
$ uname -m; go version; go env GOHOSTARCH GOARCH
aarch64
go version go1.26.7 linux/arm64
arm64
arm64

This is arm64 (saturating float→int conversion), so `math.MaxInt64` round-trips "correctly" by accident under the old code. Let me confirm a precision-loss value behaves as expected in the probe, then adjust the tests.
← Edit /tmp/opencode/olddecode/main.go
Index: /tmp/opencode/olddecode/main.go
===================================================================
--- /tmp/opencode/olddecode/main.go
+++ /tmp/opencode/olddecode/main.go
@@ -14,8 +14,9 @@
 }
 
 func main() {
 	cases := map[string]string{
+		"2^53+1 9007199254740993":      `{"schema":"astarte_flow/message/v0.1","key":"k","type":"integer","timestamp_us":0,"data":9007199254740993}`,
 		"maxint64 9223372036854775807": `{"schema":"astarte_flow/message/v0.1","key":"k","type":"integer","timestamp_us":0,"data":9223372036854775807}`,
 		"3.7":                          `{"schema":"astarte_flow/message/v0.1","key":"k","type":"integer","timestamp_us":0,"data":3.7}`,
 		"1e300":                        `{"schema":"astarte_flow/message/v0.1","key":"k","type":"integer","timestamp_us":0,"data":1e300}`,
 		"2^63":                         `{"schema":"astarte_flow/message/v0.1","key":"k","type":"integer","timestamp_us":0,"data":9223372036854775808}`,

$ cd /tmp/opencode/olddecode && go run .
maxint64 9223372036854775807   old decode -> 9223372036854775807
3.7                            old decode -> 3
1e300                          old decode -> 9223372036854775807
2^63                           old decode -> 9223372036854775807
2^53+1 9007199254740993        old decode -> 9007199254740992

Confirmed: `2^53+1` loses precision under the old path on this machine. I'll add that row too, so the test suite provably breaks without the fix on every architecture (while keeping the task-named `math.MaxInt64` row, which wraps→breaks on amd64 but saturates on arm64).
← Edit internal/flow/message_test.go
Index: /root/astrate-mule/internal/flow/message_test.go
===================================================================
--- /root/astrate-mule/internal/flow/message_test.go
+++ /root/astrate-mule/internal/flow/message_test.go
@@ -22,8 +22,17 @@
 		Timestamp: 1551884045074181,
 	},
 },
 {
+	name: "integer beyond 2^53 precision",
+	msg: Message{
+		Key:       "sensor/precise",
+		Type:      TypeInteger,
+		Data:      int64(1)<<53 + 1,
+		Timestamp: 1551884045074181,
+	},
+},
+{
 	name: "integer",
 	msg: Message{
 		Key:       "sensor/temp",
 		Type:      TypeInteger,

$ gofmt -l internal/flow/message.go internal/flow/message_test.go; go test ./internal/flow/ -run 'FlowMessage|SetDataFromWire' -v 2>&1 | tail -40
=== RUN   TestFlowMessage_RoundTrip/string
=== RUN   TestFlowMessage_RoundTrip/binary
=== RUN   TestFlowMessage_RoundTrip/datetime
=== RUN   TestFlowMessage_RoundTrip/map_with_mixed_types
=== RUN   TestFlowMessage_RoundTrip/empty_metadata
--- PASS: TestFlowMessage_RoundTrip (0.00s)
    --- PASS: TestFlowMessage_RoundTrip/integer_max_int64 (0.00s)
    --- PASS: TestFlowMessage_RoundTrip/integer_beyond_2^53_precision (0.00s)
    --- PASS: TestFlowMessage_RoundTrip/integer (0.00s)
    --- PASS: TestFlowMessage_RoundTrip/real (0.00s)
    --- PASS: TestFlowMessage_RoundTrip/boolean (0.00s)
    --- PASS: TestFlowMessage_RoundTrip/string (0.00s)
    --- PASS: TestFlowMessage_RoundTrip/binary (0.00s)
    --- PASS: TestFlowMessage_RoundTrip/datetime (0.00s)
    --- PASS: TestFlowMessage_RoundTrip/map_with_mixed_types (0.00s)
    --- PASS: TestFlowMessage_RoundTrip/empty_metadata (0.00s)
=== RUN   TestFlowMessage_WireFormatSchema
--- PASS: TestFlowMessage_WireFormatSchema (0.00s)
=== RUN   TestFlowMessage_UnmarshalRejectsUnknownSchema
--- PASS: TestFlowMessage_UnmarshalRejectsUnknownSchema (0.00s)
=== RUN   TestFlowMessage_IntegerWireRejectsFractionAndOverflow
=== RUN   TestFlowMessage_IntegerWireRejectsFractionAndOverflow/3.7
=== RUN   TestFlowMessage_IntegerWireRejectsFractionAndOverflow/1e300
=== RUN   TestFlowMessage_IntegerWireRejectsFractionAndOverflow/9223372036854775808
--- PASS: TestFlowMessage_IntegerWireRejectsFractionAndOverflow (0.00s)
    --- PASS: TestFlowMessage_IntegerWireRejectsFractionAndOverflow/3.7 (0.00s)
    --- PASS: TestFlowMessage_IntegerWireRejectsFractionAndOverflow/1e300 (0.00s)
    --- PASS: TestFlowMessage_IntegerWireRejectsFractionAndOverflow/9223372036854775808 (0.00s)
=== RUN   TestSetDataFromWire_IntegerRejectsFractionalFloat64
--- PASS: TestSetDataFromWire_IntegerRejectsFractionalFloat64 (0.00s)
=== RUN   TestFlowMessage_MapSubtypesOmittedWhenEmpty
--- PASS: TestFlowMessage_MapSubtypesOmittedWhenEmpty (0.00s)
=== RUN   TestFlowMessage_ScalarSubtypeOmittedWhenEmpty
--- PASS: TestFlowMessage_ScalarSubtypeOmittedWhenEmpty (0.00s)
=== RUN   TestFlowMessage_BinaryBase64Encoding
--- PASS: TestFlowMessage_BinaryBase64Encoding (0.00s)
=== RUN   TestFlowMessage_MapBinaryBase64Encoding
--- PASS: TestFlowMessage_MapBinaryBase64Encoding (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/flow	0.039s

All pass. Now the full gate: `go vet ./...` and `gofmt -l .`, then the full test suite (container-free).
$ go vet ./... && gofmt -l .
(no output)

$ go test ./...
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
ok  	github.com/astrate-platform/astrate/examples/flow-container-echo	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow	0.784s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	1.403s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	0.099s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.578s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	0.038s
ok  	github.com/astrate-platform/astrate/internal/flowapi	0.096s
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

Let me also try the race detector on the changed package (MULE gate preference).
$ go test -race ./internal/flow/
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/flow	0.048s
FAIL

The race detector is unavailable here (39-bit VMA kernel), as AGENTS.md anticipates — I ran `go test ./...` without `-race`.

Done: flow-msg-json-integer-precision: decode wire messages with `json.Decoder`+`UseNumber` and parse integer data via `json.Number.Int64()`, so >2^53 values survive exactly and non-integral/out-of-int64 forms (`3.7`, `1e300`) are rejected instead of silently truncated/wrapped
Files: internal/flow/message.go, internal/flow/message_test.go
Verified: `go vet ./...` -> pass; `gofmt -l .` -> clean; `go test ./...` -> pass (all packages); `go test -race ./internal/flow/` -> fail, ThreadSanitizer unsupported VMA range (39-bit kernel, expected here); scoped `go test ./internal/flow/` -> pass
Unsure:
- Two things the task's writeup assumed that this machine disproves, measured with a /tmp probe: on arm64 the old `int64(float64)` path *saturates* to MaxInt64 (it does not wrap to MinInt64), so a `math.MaxInt64` round-trip row alone would not break the old code here. I kept that row (it does break on wrapping architectures) and added a `2^53+1` row and a `9223372036854775808` wire row that provably fail the old path on every arch.
- Switching the decoder to UseNumber makes *all* wire numbers `json.Number`, so I also made `setDataFromWireMap` coerce integer/real fields per `FieldTypes` (int64/float64). That is a small fix beyond the named scalar bug — previously a wired TypeMap integer field silently stayed float64, and without the coercion it would now leak as `json.Number` (Data-contract violation). Flagging for review.
Follow-ups: none
