slug: payload-longinteger-fraction-quantize
verdict: done
at:  d7babc6
ran: 2026-09-12T18:41:28Z on DietPi in 356s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read pkg/payload/json.go [offset=230, limit=60]
✱ Grep "TestJSONProfile" in pkg/payload · 2 matches
→ Read pkg/payload/payload_test.go [limit=120, offset=255]
✱ Grep "jsonInt64|maxSafeJSONInt" in pkg/payload · 9 matches
→ Read pkg/payload/json.go
✱ Grep "math\." in pkg/payload · 37 matches
✱ Grep "func checkDouble|func int32FromInt64" in pkg/payload · 2 matches
→ Read pkg/payload/value.go [limit=60, offset=210]
Now let me verify `big.Rat.SetString` behavior on the relevant tokens with a throwaway probe:

← Write /tmp/opencode/ratprobe/main.go
Wrote file successfully.

$ cd /tmp/opencode/ratprobe && go mod init ratprobe >/dev/null 2>&1; go run main.go
5.0                    -> IsInt=true  rat=5  float=5.0
1e3                    -> IsInt=true  rat=1000  float=1e3
0.5e1                  -> IsInt=true  rat=5  float=0.5e1
42                     -> IsInt=true  rat=42  float=42
4503599627370496.5     -> IsInt=false  rat=9007199254740993/2  float=4503599627370496.5
6000000000000000.9     -> IsInt=false  rat=60000000000000009/10  float=6000000000000000.9
5e-400                 -> IsInt=false  rat=1/2000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000  float=5e-400
1.25                   -> IsInt=false  rat=5/4  float=1.25
1e03                   -> IsInt=true  rat=1000  float=1e03
-0.25e-400             -> IsInt=false  rat=-1/40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000  float=-0.25e-400

Probe confirms the plan. Quick check of ParseFloat error semantics:
$ cd /tmp/opencode/ratprobe && cat > main.go <<'EOF'
package main

import (
	"fmt"
	"strconv"
)

func main() {
	for _, t := range []string{"5e-400", "1e400", "4503599627370496.5", "6000000000000000.9", "1e1000000000"} {
		f, err := strconv.ParseFloat(t, 64)
		fmt.Printf("%-22s -> f=%v err=%v\n", t, f, err)
	}
}
EOF
go run main.go
5e-400                 -> f=0 err=<nil>
1e400                  -> f=+Inf err=strconv.ParseFloat: parsing "1e400": value out of range
4503599627370496.5     -> f=4.503599627370496e+15 err=<nil>
6000000000000000.9     -> f=6.000000000000001e+15 err=<nil>
1e1000000000           -> f=+Inf err=strconv.ParseFloat: parsing "1e1000000000": value out of range

← Edit pkg/payload/json.go
Index: /root/astrate-mule/pkg/payload/json.go
===================================================================
--- /root/astrate-mule/pkg/payload/json.go
+++ /root/astrate-mule/pkg/payload/json.go
@@ -3,9 +3,9 @@
 import (
 	"encoding/base64"
 	"encoding/json"
 	"errors"
-	"math"
+	"math/big"
 	"strconv"
 	"time"
 
 	"github.com/astrate-platform/astrate/pkg/interfaceschema"

Probe confirms both the bug and the fix semantics. Now the edit:
← Edit pkg/payload/json.go
Index: /root/astrate-mule/pkg/payload/json.go
===================================================================
--- /root/astrate-mule/pkg/payload/json.go
+++ /root/astrate-mule/pkg/payload/json.go
@@ -247,12 +247,12 @@
 }
 
 // jsonInt64 converts a JSON number token to int64. Plain decimal-integer
 // tokens convert exactly across the whole int64 range. Fraction or exponent
-// forms are accepted only when integral and strictly inside the float64-exact
-// window (|x| < 2^53): beyond it the decimal→binary parse itself may round,
-// silently corrupting the value — such values must use the decimal-string
-// form (§3.5.3).
+// forms are accepted only when the token is an exact decimal integer (5.0,
+// 1e3, 0.5e1) and strictly inside the float64-exact window (|x| < 2^53):
+// beyond it the decimal→binary parse itself may round, silently corrupting
+// the value — such values must use the decimal-string form (§3.5.3).
 func jsonInt64(num json.Number) (int64, error) {
 	i, err := strconv.ParseInt(string(num), 10, 64)
 	if err == nil {
 		return i, nil
@@ -263,14 +263,26 @@
 	f, err := strconv.ParseFloat(string(num), 64)
 	if err != nil {
 		return 0, rejectf(ReasonTypeMismatch, "number %s is not decodable", num)
 	}
-	if f != math.Trunc(f) || f <= -float64(maxSafeJSONInt) || f >= float64(maxSafeJSONInt) {
+	// Fallback path: ParseFloat rounds, so judge the token, not the parsed
+	// value. 5e-400 parses to an integral 0 and 4503599627370496.5 rounds
+	// onto an integer, yet neither is an exact decimal integer.
+	if !jsonTokenIsExactInt(string(num)) || f <= -float64(maxSafeJSONInt) || f >= float64(maxSafeJSONInt) {
 		return 0, rejectf(ReasonTypeMismatch, "number %s does not coerce to an integer exactly", num)
 	}
 	return int64(f), nil
 }
 
+// jsonTokenIsExactInt reports whether the decimal token denotes an exact
+// integer (42, 5.0, 1e3, 0.5e1), as opposed to a fraction (1.25) or an
+// underflowing magnitude (5e-400) that a float64 parse would place on or
+// round toward an integer.
+func jsonTokenIsExactInt(token string) bool {
+	r, ok := new(big.Rat).SetString(token)
+	return ok && r.IsInt()
+}
+
 // clip truncates a raw token for use in human-readable reject details.
 func clip(raw []byte) string {
 	const maxDetail = 64
 	if len(raw) > maxDetail {

Now the test rows:
← Edit pkg/payload/payload_test.go
Index: /root/astrate-mule/pkg/payload/payload_test.go
===================================================================
--- /root/astrate-mule/pkg/payload/payload_test.go
+++ /root/astrate-mule/pkg/payload/payload_test.go
@@ -302,9 +302,13 @@
 {name: "longinteger number beyond 2^53", in: `{"v":9007199254740993}`, vt: interfaceschema.LongInteger, want: int64(9007199254740993)},
 {name: "longinteger decimal string", in: `{"v":"9223372036854775807"}`, vt: interfaceschema.LongInteger, want: int64(math.MaxInt64)},
 {name: "longinteger negative string", in: `{"v":"-9223372036854775808"}`, vt: interfaceschema.LongInteger, want: int64(math.MinInt64)},
 {name: "longinteger exact-fit float", in: `{"v":1e3}`, vt: interfaceschema.LongInteger, want: int64(1000)},
+{name: "longinteger exponent-integral 0.5e1", in: `{"v":0.5e1}`, vt: interfaceschema.LongInteger, want: int64(5)},
 {name: "longinteger fractional", in: `{"v":1.25}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
+{name: "longinteger fractional rounds", in: `{"v":4503599627370496.5}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
+{name: "longinteger fractional rounds up", in: `{"v":6000000000000000.9}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
+{name: "longinteger underflow to zero", in: `{"v":5e-400}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
 {name: "longinteger bad string", in: `{"v":"12abc"}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
 {name: "longinteger number overflow", in: `{"v":9223372036854775808}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
 {name: "longinteger string overflow", in: `{"v":"9223372036854775808"}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
 {name: "longinteger float form beyond 2^53", in: `{"v":9007199254740993.0}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},

Now run the gate:
$ go vet ./... && go test ./... && gofmt -l .
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	5.092s
ok  	github.com/astrate-platform/astrate/internal/broker	2.997s
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	3.460s
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	3.187s
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
ok  	github.com/astrate-platform/astrate/internal/testutil	0.397s
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	0.186s
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	0.108s

All green, `gofmt -l .` clean. Final diff review:
$ git diff -- pkg/payload/json.go pkg/payload/payload_test.go
diff --git a/pkg/payload/json.go b/pkg/payload/json.go
index b0f4967..a18a896 100644
--- a/pkg/payload/json.go
+++ b/pkg/payload/json.go
@@ -4,7 +4,7 @@ import (
 	"encoding/base64"
 	"encoding/json"
 	"errors"
-	"math"
+	"math/big"
 	"strconv"
 	"time"
 
@@ -248,10 +248,10 @@ func jsonNumber(raw json.RawMessage) (json.Number, error) {
 
 // jsonInt64 converts a JSON number token to int64. Plain decimal-integer
 // tokens convert exactly across the whole int64 range. Fraction or exponent
-// forms are accepted only when integral and strictly inside the float64-exact
-// window (|x| < 2^53): beyond it the decimal→binary parse itself may round,
-// silently corrupting the value — such values must use the decimal-string
-// form (§3.5.3).
+// forms are accepted only when the token is an exact decimal integer (5.0,
+// 1e3, 0.5e1) and strictly inside the float64-exact window (|x| < 2^53):
+// beyond it the decimal→binary parse itself may round, silently corrupting
+// the value — such values must use the decimal-string form (§3.5.3).
 func jsonInt64(num json.Number) (int64, error) {
 	i, err := strconv.ParseInt(string(num), 10, 64)
 	if err == nil {
@@ -264,12 +264,24 @@ func jsonInt64(num json.Number) (int64, error) {
 	if err != nil {
 		return 0, rejectf(ReasonTypeMismatch, "number %s is not decodable", num)
 	}
-	if f != math.Trunc(f) || f <= -float64(maxSafeJSONInt) || f >= float64(maxSafeJSONInt) {
+	// Fallback path: ParseFloat rounds, so judge the token, not the parsed
+	// value. 5e-400 parses to an integral 0 and 4503599627370496.5 rounds
+	// onto an integer, yet neither is an exact decimal integer.
+	if !jsonTokenIsExactInt(string(num)) || f <= -float64(maxSafeJSONInt) || f >= float64(maxSafeJSONInt) {
 		return 0, rejectf(ReasonTypeMismatch, "number %s does not coerce to an integer exactly", num)
 	}
 	return int64(f), nil
 }
 
+// jsonTokenIsExactInt reports whether the decimal token denotes an exact
+// integer (42, 5.0, 1e3, 0.5e1), as opposed to a fraction (1.25) or an
+// underflowing magnitude (5e-400) that a float64 parse would place on or
+// round toward an integer.
+func jsonTokenIsExactInt(token string) bool {
+	r, ok := new(big.Rat).SetString(token)
+	return ok && r.IsInt()
+}
+
 // clip truncates a raw token for use in human-readable reject details.
 func clip(raw []byte) string {
 	const maxDetail = 64
diff --git a/pkg/payload/payload_test.go b/pkg/payload/payload_test.go
index 60c000c..8d512ee 100644
--- a/pkg/payload/payload_test.go
+++ b/pkg/payload/payload_test.go
@@ -303,7 +303,11 @@ func TestJSONProfile(t *testing.T) {
 		{name: "longinteger decimal string", in: `{"v":"9223372036854775807"}`, vt: interfaceschema.LongInteger, want: int64(math.MaxInt64)},
 		{name: "longinteger negative string", in: `{"v":"-9223372036854775808"}`, vt: interfaceschema.LongInteger, want: int64(math.MinInt64)},
 		{name: "longinteger exact-fit float", in: `{"v":1e3}`, vt: interfaceschema.LongInteger, want: int64(1000)},
+		{name: "longinteger exponent-integral 0.5e1", in: `{"v":0.5e1}`, vt: interfaceschema.LongInteger, want: int64(5)},
 		{name: "longinteger fractional", in: `{"v":1.25}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
+		{name: "longinteger fractional rounds", in: `{"v":4503599627370496.5}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
+		{name: "longinteger fractional rounds up", in: `{"v":6000000000000000.9}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
+		{name: "longinteger underflow to zero", in: `{"v":5e-400}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
 		{name: "longinteger bad string", in: `{"v":"12abc"}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
 		{name: "longinteger number overflow", in: `{"v":9223372036854775808}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},
 		{name: "longinteger string overflow", in: `{"v":"9223372036854775808"}`, vt: interfaceschema.LongInteger, wantReason: ReasonTypeMismatch},

`jsonInt64` now verifies the token is an exact decimal integer (`math/big.Rat` SetString + IsInt) on the float fallback before the untouched `|x| < 2^53` window check. The three silent-rounding cases (`4503599627370496.5`, `6000000000000000.9`, `5e-400`) now reject as ReasonTypeMismatch; `5.0`/`1e3`/`0.5e1` stay accepted, and the plain-integer ParseInt fast path is unchanged. No goroutines touched, so no race concern. Ran `go test ./...` (without `-race`; this box is the Pi) — gate passed.

```
Done: payload-longinteger-fraction-quantize: jsonInt64 float fallback now requires the token to be an exact decimal integer (big.Rat IsInt), rejecting silent ParseFloat rounding/underflow as ReasonTypeMismatch
Files: pkg/payload/json.go, pkg/payload/payload_test.go
Verified: go vet ./... && go test ./... -> pass; gofmt -l . -> clean (no -race: Pi host, no goroutines touched)
Unsure: nothing
Follow-ups: none
```
