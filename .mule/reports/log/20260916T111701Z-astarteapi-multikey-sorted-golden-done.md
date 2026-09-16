slug: astarteapi-multikey-sorted-golden
verdict: done
at:  bfffbed
ran: 2026-09-16T11:17:01Z on DietPi in 215s


> build · big-pickle

→ Read .mule/task.md
→ Read .mule/MULE.md
→ Read pkg/astarteapi/envelope_test.go [offset=1, limit=150]
→ Read pkg/astarteapi/envelope.go [offset=160, limit=40]
→ Read pkg/astarteapi/envelope.go [limit=40, offset=130]
✗ Read pkg/astarteapi/golden/error_fields.json failed
Error: File not found: /root/astrate-mule/pkg/astarteapi/golden/error_fields.json
✗ Read pkg/astarteapi/golden/error_raw_nested.json failed
Error: File not found: /root/astrate-mule/pkg/astarteapi/golden/error_raw_nested.json
✱ Glob "**/golden/**/*.json" 0 matches
→ Read pkg/astarteapi/envelope.go [limit=40, offset=1]
✱ Glob "pkg/astarteapi/testdata/*" 17 matches
→ Read pkg/astarteapi/testdata/error_fields.json
→ Read pkg/astarteapi/testdata/error_raw_nested.json
→ Read pkg/astarteapi/envelope_test.go [offset=150, limit=110]

← Edit pkg/astarteapi/envelope_test.go
Index: /root/astrate-mule/pkg/astarteapi/envelope_test.go
===================================================================
--- /root/astrate-mule/pkg/astarteapi/envelope_test.go
+++ /root/astrate-mule/pkg/astarteapi/envelope_test.go
@@ -127,17 +127,39 @@
 		name: "raw nested errors", golden: "error_raw_nested.json", wantStatus: 422,
 		write: func(w http.ResponseWriter) error {
 			// Upstream nests changeset errors inside the failing part of a
 			// trigger definition; keys must come out sorted with no HTML
-			// escaping and no trailing newline.
+			// escaping and without a trailing newline.
 			return astarteapi.WriteRawErrors(w, 422, map[string]any{
 				"action": map[string][]string{
 					"http_url":    {"should be at least 8 character(s)", "must be a valid http(s) URL"},
 					"http_method": {"is invalid"},
 				},
 			})
 		},
 	},
+	{
+		name: "field errors multi-key", golden: "error_fields_multikey.json", wantStatus: 422,
+		write: func(w http.ResponseWriter) error {
+			// Multi-key field errors: top-level map keys must be sorted
+			// on the wire so a sort regression cannot pass the suite.
+			return astarteapi.WriteFieldErrors(w, 422, map[string][]string{
+				"name": {"must not be empty"},
+				"email": {"is invalid"},
+			})
+		},
+	},
+	{
+		name: "raw errors multi-key", golden: "error_raw_multikey.json", wantStatus: 422,
+		write: func(w http.ResponseWriter) error {
+			// Multi-key raw errors: top-level map keys must be sorted
+			// on the wire so a sort regression cannot pass the suite.
+			return astarteapi.WriteRawErrors(w, 422, map[string][]string{
+				"http_url":    {"should be at least 8 character(s)"},
+				"http_method": {"is invalid"},
+			})
+		},
+	},
 }
 for _, tc := range cases {
 	t.Run(tc.name, func(t *testing.T) {
 		rec := httptest.NewRecorder()

Now I understand the task. The existing `WriteFieldErrors` and `WriteRawErrors` golden fixtures use single-key maps at the top level, so a sort regression would pass. I need to add multi-key test cases and their golden files.
Let me verify the exact JSON output before creating the golden files:
$ cat > /tmp/verify.go << 'EOF'
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	f1 := map[string][]string{
		"name":  {"must not be empty"},
		"email": {"is invalid"},
	}
	w1, _ := json.Marshal(struct {
		Errors map[string][]string `json:"errors"`
	}{Errors: f1})
	fmt.Println(string(w1))

	f2 := map[string][]string{
		"http_url":    {"should be at least 8 character(s)"},
		"http_method": {"is invalid"},
	}
	w2, _ := json.Marshal(struct {
		Errors map[string][]string `json:"errors"`
	}{Errors: f2})
	fmt.Println(string(w2))
}
EOF
go run /tmp/verify.go
{"errors":{"email":["is invalid"],"name":["must not be empty"]}}
{"errors":{"http_method":["is invalid"],"http_url":["should be at least 8 character(s)"]}}

$ cat > /tmp/verify2.go << 'EOF'
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// Simulate WriteRawErrors with map[string]any (same structure as existing test)
	f := map[string]any{
		"http_url":    []string{"should be at least 8 character(s)"},
		"http_method": []string{"is invalid"},
	}
	w, _ := json.Marshal(struct {
		Errors any `json:"errors"`
	}{Errors: f})
	fmt.Println(string(w))
}
EOF
go run /tmp/verify2.go
{"errors":{"http_method":["is invalid"],"http_url":["should be at least 8 character(s)"]}}

← Write pkg/astarteapi/testdata/error_fields_multikey.json
Wrote file successfully.

Now create the golden files and run the gate:
← Write pkg/astarteapi/testdata/error_raw_multikey.json
Wrote file successfully.

$ gofmt -l pkg/astarteapi/
pkg/astarteapi/envelope_test.go

$ go test ./pkg/astarteapi/...
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	0.169s

Tests pass, but `gofmt` flagged the test file. Let me fix formatting:
$ gofmt -w pkg/astarteapi/envelope_test.go
(no output)

→ Read pkg/astarteapi/envelope_test.go [limit=40, offset=117]

$ gofmt -l pkg/astarteapi/ && go test ./pkg/astarteapi/...
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)

Now run the full gate as the task requires:
$ go vet ./... && go test ./... 2>&1 | tail -20
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

All gates pass. Report:

```
Done: Extended TestGoldenEnvelopes (pkg/astarteapi/envelope_test.go) with multi-key WriteFieldErrors and WriteRawErrors cases plus golden fixtures asserting top-level map keys are emitted in sorted order on the wire
Files: pkg/astarteapi/envelope_test.go, pkg/astarteapi/testdata/error_fields_multikey.json, pkg/astarteapi/testdata/error_raw_multikey.json
Verified: go vet ./... && go test ./... -> pass; gofmt -l pkg/astarteapi/ -> clean
Unsure: nothing
Follow-ups: none
```

The two new cases freeze sorted multi-key output (`error_fields_multikey` sorts `email` before `name`; `error_raw_multikey` sorts `http_method` before `http_url` at the top level — the previous raw fixture only exercised sorted order one level down). A change that stopped sorting would now produce different bytes and fail the golden comparison.
