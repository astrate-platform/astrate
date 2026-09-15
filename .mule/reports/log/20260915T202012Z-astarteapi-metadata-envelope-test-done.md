slug: astarteapi-metadata-envelope-test
verdict: done
at:  1365cb8
ran: 2026-09-15T20:20:12Z on DietPi in 238s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read pkg/astarteapi/envelope.go

$ ls pkg/astarteapi/
envelope.go
envelope_test.go
testdata

✱ Grep "TestGoldenEnvelopes" in pkg/astarteapi · 2 matches
→ Read pkg/astarteapi/envelope_test.go

$ ls pkg/astarteapi/testdata/
data_array.json
data_links.json
data_links_last.json
data_null.json
data_object.json
data_pem.json
error_bad_request.json
error_custom_detail.json
error_device_not_found.json
error_fields.json
error_forbidden.json
error_internal.json
error_not_found.json
error_raw_nested.json
error_unauthorized.json

$ cat pkg/astarteapi/testdata/data_links.json pkg/astarteapi/testdata/data_array.json
{"data":["dT6hS2W9TT6LEnP25ks_lg"],"links":{"self":"/v1/test/devices?details=true&limit=1","next":"/v1/test/devices?details=true&from_token=dT6hS2W9TT6LEnP25ks_lg&limit=1"}}{"data":["dT6hS2W9TT6LEnP25ks_lg"]}

→ Read internal/testutil
→ Read internal/testutil/golden.go
→ Read internal/appengine/data.go [offset=240, limit=40]
Now I understand the golden test pattern. Adding two cases for `WriteDataWithMetadata` to `TestGoldenEnvelopes`.

← Edit pkg/astarteapi/envelope_test.go
Index: /root/astrate-mule/pkg/astarteapi/envelope_test.go
===================================================================
--- /root/astrate-mule/pkg/astarteapi/envelope_test.go
+++ /root/astrate-mule/pkg/astarteapi/envelope_test.go
@@ -67,8 +67,26 @@
 			astarteapi.Links{Self: "/v1/test/devices?limit=1"})
 	},
 },
 {
+	name: "data with nil metadata", golden: "data_nil_metadata.json", wantStatus: 200,
+	write: func(w http.ResponseWriter) error {
+		// nil metadata must omit the key entirely, not emit null.
+		return astarteapi.WriteDataWithMetadata(w, 200, []string{"dT6hS2W9TT6LEnP25ks_lg"}, nil)
+	},
+},
+{
+	name: "data with metadata", golden: "data_metadata.json", wantStatus: 200,
+	write: func(w http.ResponseWriter) error {
+		// Realistic format=table shape from renderIndividual.
+		return astarteapi.WriteDataWithMetadata(w, 200, []string{"dT6hS2W9TT6LEnP25ks_lg"},
+			map[string]any{
+				"columns":      map[string]int{"timestamp": 0, "temperature": 1},
+				"table_header": []string{"timestamp", "temperature"},
+			})
+	},
+},
+{
 	name: "bad request", golden: "error_bad_request.json", wantStatus: 400,
 	write: astarteapi.WriteBadRequest,
 },
 {

Now generate the golden files:
$ go test ./pkg/astarteapi/ -run TestGoldenEnvelopes -update 2>&1
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	0.165s

$ cat pkg/astarteapi/testdata/data_nil_metadata.json; echo; cat pkg/astarteapi/testdata/data_metadata.json; echo
{"data":["dT6hS2W9TT6LEnP25ks_lg"]}
{"data":["dT6hS2W9TT6LEnP25ks_lg"],"metadata":{"columns":{"temperature":1,"timestamp":0},"table_header":["timestamp","temperature"]}}

$ go test ./pkg/astarteapi/ 2>&1 && gofmt -l pkg/astarteapi/ 2>&1
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	0.144s

Let me verify the nil-branch test actually pins `omitempty` by probing what removing it would produce:
$ cat > /tmp/opencode/probe_omitempty.go <<'EOF'
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type withOmit struct {
	Data     any            `json:"data"`
	Metadata map[string]any `json:"metadata,omitempty"`
}
type withoutOmit struct {
	Data     any            `json:"data"`
	Metadata map[string]any `json:"metadata"`
}

func marshal(v any) []byte {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
	return buf.Bytes()[:buf.Len()-1]
}

func main() {
	fmt.Printf("with omitempty,    nil: %s\n", marshal(withOmit{Data: []string{"x"}, Metadata: nil}))
	fmt.Printf("without omitempty, nil: %s\n", marshal(withoutOmit{Data: []string{"x"}, Metadata: nil}))
}
EOF
go run /tmp/opencode/probe_omitempty.go
with omitempty,    nil: {"data":["x"]}
without omitempty, nil: {"data":["x"],"metadata":null}

Confirmed the test genuinely breaks without `omitempty`. Now the gate:
$ go vet ./pkg/astarteapi/ && go test -race ./pkg/astarteapi/ 2>&1
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/astarteapi	0.045s
FAIL

Race detector unsupported here (39-bit VMA, the Pi), as MULE.md predicts. Running the real gate:
$ go vet ./... && go test ./... 2>&1
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
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
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

$ gofmt -l . 2>&1
(no output)

Done: added a container-free golden test pinning both `WriteDataWithMetadata` wire branches in `TestGoldenEnvelopes` — nil metadata omits `metadata` (`{"data":[...]}` only), populated metadata emits `{"data":[...],"metadata":{...}}`, both at 200 with the package `Content-Type`.
Files: pkg/astarteapi/envelope_test.go, pkg/astarteapi/testdata/data_nil_metadata.json, pkg/astarteapi/testdata/data_metadata.json
Verified: `go test ./pkg/astarteapi/` (pass), `go vet ./...` (pass), `go test ./...` (pass), `gofmt -l .` (clean), omitempty-removal probe in /tmp (shows the nil branch would emit `"metadata":null` — test provably breaks without the pinned behaviour). `go test -race` is unsupported on this box (39-bit VMA, the Pi) — flagged as expected.
Unsure: nothing
Follow-ups: none
