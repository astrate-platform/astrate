slug: container-response-cap-test
verdict: blocked
at:  e5753dd
ran: 2026-10-02T17:50:03Z on DietPi in 271s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/task.md
→ Read .mule/MULE.md
→ Read internal/flow/blocks/container/httpbridge.go

$ ls /root/astrate-mule/internal/flow/blocks/container/
block.go
block_test.go
docker.go
httpbridge.go

→ Read internal/flow/blocks/container/block_test.go
→ Read internal/flow/blocks/info.go [offset=80, limit=80]

$ grep -rn "max_body\|maxBody\|MaxBodyBytes\|exceeds" /root/astrate-mule --include=*.go --include=*.json --include=*.yaml --include=*.yml --include=*.md -l | head -30
/root/astrate-mule/internal/realm/service.go
/root/astrate-mule/internal/realm/http_test.go
/root/astrate-mule/internal/realm/http.go
/root/astrate-mule/internal/flowapi/http.go
/root/astrate-mule/internal/pairing/ratelimit_test.go
/root/astrate-mule/internal/pairing/ca/ca_test.go
/root/astrate-mule/internal/pairing/http.go
/root/astrate-mule/internal/housekeeping/http.go
/root/astrate-mule/internal/appengine/http.go
/root/astrate-mule/internal/engine/introspection.go
/root/astrate-mule/internal/flow/blocks/container/httpbridge.go
/root/astrate-mule/internal/flow/blocks/sort.go
/root/astrate-mule/internal/flow/blocks/schema.go
/root/astrate-mule/test/conformance/upstream/verify-server-writes.json
/root/astrate-mule/test/conformance/upstream/recordchannels/main.go
/root/astrate-mule/test/conformance/upstream/README.md
/root/astrate-mule/.mule/for-giulio.md
/root/astrate-mule/.mule/reviews/container-2026-10-02.md
/root/astrate-mule/.mule/reports/log/20261002T172128Z-container-timeout-bounds-blocked.md
/root/astrate-mule/.mule/reports/log/20260925T201910Z-docs-sync-rm-delete-device-async-operation-param-done.md
/root/astrate-mule/.mule/reports/log/20260915T173739Z-hk-zero-retention-asymmetry-done.md
/root/astrate-mule/.mule/reports/log/20260910T184706Z-recipe-docs-sync-proposed.md
/root/astrate-mule/.mule/reports/log/20260912T171256Z-recipe-milestones-proposed.md
/root/astrate-mule/.mule/reports/log/20260926T194649Z-recipe-docs-sync-proposed.md
/root/astrate-mule/.mule/reports/log/20260917T190323Z-recipe-github-issues-proposed.md
/root/astrate-mule/.mule/reports/log/20260927T111624Z-docs-sync-ae-add-group-device-422-done.md
/root/astrate-mule/.mule/reports/log/20260909T195057Z-realm-pure-helper-tests-done.md
/root/astrate-mule/.mule/reports/log/20260912T113032Z-docs-sync-ae-write-value-422-done.md
/root/astrate-mule/.mule/reports/log/20261002T114427Z-recipe-code-review-proposed.md
/root/astrate-mule/.mule/reports/log/20260910T200120Z-recipe-github-issues-proposed.md

✱ Grep "container" in internal/flow/blocks/schema.go · 100 matches

$ cd /root/astrate-mule && rg -n "ready_timeout_ms" --glob '!.mule/reports/**' --glob '!.mule/reviews/**' -l
internal/flow/blocks/info.go
internal/flow/blocks/schema.go
internal/flow/blocks/container/block.go
internal/flow/blocks/container/block_test.go

→ Read internal/flow/blocks/schema.go [limit=40]

$ cd /root/astrate-mule && rg -n "type Message struct" -A 30 internal/flow/*.go
internal/flow/message.go:81:type Message struct {
internal/flow/message.go-82-	// Key identifies the stream this message belongs to. It must be non-empty.
internal/flow/message.go-83-	Key string
internal/flow/message.go-84-	// Metadata is an optional string→string map carried alongside the payload.
internal/flow/message.go-85-	Metadata map[string]string
internal/flow/message.go-86-	// Type is the base data type of the payload.
internal/flow/message.go-87-	Type DataType
internal/flow/message.go-88-	// Subtype is an optional MIME hint (meaningful when Type is TypeBinary).
internal/flow/message.go-89-	Subtype string
internal/flow/message.go-90-	// Timestamp is the event-time in microseconds since Unix epoch.
internal/flow/message.go-91-	Timestamp int64
internal/flow/message.go-92-	// Data is the payload; its concrete type must match Type (int64 for
internal/flow/message.go-93-	// TypeInteger, float64 for TypeReal, bool for TypeBoolean, time.Time for
internal/flow/message.go-94-	// TypeDatetime, []byte for TypeBinary, string for TypeString, map[string]any
internal/flow/message.go-95-	// for TypeMap).
internal/flow/message.go-96-	Data any
internal/flow/message.go-97-	// FieldTypes holds per-field types when Type is TypeMap.
internal/flow/message.go-98-	FieldTypes map[string]DataType
internal/flow/message.go-99-	// FieldSubtypes holds per-field subtypes when Type is TypeMap.
internal/flow/message.go-100-	FieldSubtypes map[string]string
internal/flow/message.go-101-}
internal/flow/message.go-102-
internal/flow/message.go-103-// wireMessage is the on-the-wire JSON representation matching the upstream
internal/flow/message.go-104-// astarte_flow/message/v0.1 schema.
internal/flow/message.go-105-type wireMessage struct {
internal/flow/message.go-106-	Schema      string            `json:"schema"`
internal/flow/message.go-107-	Key         string            `json:"key"`
internal/flow/message.go-108-	Type        any               `json:"type"`
internal/flow/message.go-109-	Subtype     any               `json:"subtype,omitempty"`
internal/flow/message.go-110-	Data        any               `json:"data"`
internal/flow/message.go-111-	Metadata    map[string]string `json:"metadata,omitempty"`

→ Read internal/flow/message.go [offset=100, limit=80]
→ Read internal/flow/blocks/schema_test.go

← Edit internal/flow/blocks/container/block_test.go
Index: /root/astrate-mule/internal/flow/blocks/container/block_test.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/container/block_test.go
+++ /root/astrate-mule/internal/flow/blocks/container/block_test.go
@@ -466,8 +466,82 @@
 		t.Fatal("want not ready error")
 	}
 }
 
+// sizedMessage marshals a valid string message whose JSON body is exactly n
+// bytes, returning the body and the payload it carries. Padding the payload
+// rather than hand-writing wire JSON keeps the length exact regardless of what
+// the marshaller puts around it.
+func sizedMessage(t *testing.T, n int) (body []byte, data string) {
+	t.Helper()
+	pad := 0
+	for {
+		data = strings.Repeat("x", pad)
+		body, err := json.Marshal(&flow.Message{Key: "k", Type: flow.TypeString, Data: data})
+		if err != nil {
+			t.Fatalf("marshal padded message: %v", err)
+		}
+		switch {
+		case len(body) == n:
+			return body, data
+		case len(body) > n:
+			t.Fatalf("cannot size a message to %d bytes: smallest is %d", n, len(body))
+		default:
+			pad += n - len(body)
+		}
+	}
+}
+
+// TestBridge_RoundTripResponseCap pins the cap on the container's answer. The
+// container is tenant-supplied and this is the only bound between what it
+// chooses to emit and process memory on the message path: a body one byte past
+// the cap must fail rather than be parsed, and one byte exactly at it must
+// arrive whole.
+func TestBridge_RoundTripResponseCap(t *testing.T) {
+	const cap = 4096 // stands in for the 1 MiB default, so the fixture stays small
+	for _, tc := range []struct {
+		name      string
+		bodyBytes int
+		wantErr   string // "" means the message must round-trip
+	}{
+		{name: "exactly at the cap", bodyBytes: cap},
+		{name: "one byte over the cap", bodyBytes: cap + 1, wantErr: "exceeds"},
+	} {
+		t.Run(tc.name, func(t *testing.T) {
+			body, data := sizedMessage(t, tc.bodyBytes)
+			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
+				w.Header().Set("Content-Type", "application/json")
+				_, _ = w.Write(body)
+			}))
+			t.Cleanup(srv.Close)
+
+			bridge := &container.Bridge{BaseURL: srv.URL, MaxBodyBytes: cap}
+			outs, err := bridge.RoundTrip(context.Background(),
+				&flow.Message{Key: "k", Type: flow.TypeString, Data: "in"})
+
+			if tc.wantErr != "" {
+				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
+					t.Fatalf("RoundTrip with a %d-byte answer: err = %v, want the %q error", len(body), err, tc.wantErr)
+				}
+				if outs != nil {
+					t.Errorf("outs = %v, want none alongside the error", outs)
+				}
+				return
+			}
+			if err != nil {
+				t.Fatalf("RoundTrip with a %d-byte answer: %v, want success", len(body), err)
+			}
+			if len(outs) != 1 {
+				t.Fatalf("outs len = %d, want 1", len(outs))
+			}
+			// Byte-exact: a cap that truncated the last byte would lose a pad char.
+			if got, _ := outs[0].Data.(string); got != data {
+				t.Errorf("data = %d bytes, want the %d bytes sent", len(got), len(data))
+			}
+		})
+	}
+}
+
 func TestCLIRunner_ParsesDockerPort(t *testing.T) {
 	var calls [][]string
 	r := &container.CLIRunner{
 		Run: func(_ context.Context, name string, args ...string) (string, string, error) {

Now the two documentation lines:
← Edit internal/flow/blocks/info.go
Index: /root/astrate-mule/internal/flow/blocks/info.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/info.go
+++ /root/astrate-mule/internal/flow/blocks/info.go
@@ -112,9 +112,9 @@
 TypeContainer: {
 	Type:    TypeContainer,
 	Role:    RoleTransform,
 	Summary: "Run a Docker image as a transform (HTTP POST /v1/message); local Docker only (PoC)",
-	Config:  "image (required), config (object→ASTRATE_FLOW_CONFIG), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",
+	Config:  "image (required), config (object→ASTRATE_FLOW_CONFIG), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000); responses over 1 MiB are rejected",
 },
 TypeVirtualDevicePool: {
 	Type:    TypeVirtualDevicePool,
 	Role:    RoleSink,

← Edit internal/flow/blocks/schema.go
Index: /root/astrate-mule/internal/flow/blocks/schema.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/schema.go
+++ /root/astrate-mule/internal/flow/blocks/schema.go
@@ -15,9 +15,9 @@
 	TypeHTTPSource:        `{"type":"object","properties":{"urls":{"type":"array","description":"URLs to poll GET round-robin (at least one)"},"interval_ms":{"type":"integer","description":"Milliseconds between polls (default 1000)"},"timeout_ms":{"type":"integer","description":"Per-request timeout in milliseconds (default 5000)"}},"required":["urls"]}`,
 	TypeHTTPSink:          `{"type":"object","properties":{"url":{"type":"string","description":"Base URL to POST payloads to"},"method":{"type":"string","description":"HTTP method (default \"POST\")"},"timeout_ms":{"type":"integer","description":"Per-request timeout in milliseconds (default 5000)"},"headers":{"type":"object","description":"Extra request headers string→string"}},"required":["url"]}`,
 	TypeMQTTSource:        `{"type":"object","properties":{"url":{"type":"string","description":"Broker URL (e.g. tcp://127.0.0.1:1883)"},"topics":{"type":"array","description":"Topics to subscribe to (at least one, non-empty strings)"},"qos":{"type":"integer","description":"Subscription QoS level 0|1|2 (default 0)"},"client_id":{"type":"string","description":"Optional MQTT client ID (random when unset)"},"username":{"type":"string","description":"Optional username"},"password":{"type":"string","description":"Optional password"}},"required":["url","topics"]}`,
 	TypeMQTTSink:          `{"type":"object","properties":{"url":{"type":"string","description":"Broker URL (e.g. tcp://127.0.0.1:1883)"},"topic":{"type":"string","description":"Topic to publish to"},"qos":{"type":"integer","description":"Publish QoS level 0|1|2 (default 0)"},"retained":{"type":"boolean","description":"Set the retain flag on published messages (default false)"},"client_id":{"type":"string","description":"Optional MQTT client ID (random when unset)"},"username":{"type":"string","description":"Optional username"},"password":{"type":"string","description":"Optional password"}},"required":["url","topic"]}`,
-	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as ASTRATE_FLOW_CONFIG"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
+	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as ASTRATE_FLOW_CONFIG"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"description":"Container answers over 1 MiB are rejected","required":["image"]}`,
 	TypeVirtualDevicePool: `{"type":"object","properties":{"devices":{"type":"array","items":{"type":"string"},"minItems":1,"description":"Registered device_ids this pool may publish as"}},"required":["devices"]}`,
 	TypeNullSink:          `{"type":"object","properties":{}}`,
 	TypeLogSink:           `{"type":"object","properties":{}}`,
 }

$ cd /root/astrate-mule && gofmt -l . && go vet ./internal/flow/blocks/... && go test ./internal/flow/blocks/... -run 'Cap|RoundTrip|Schema|Info' -v 2>&1 | tail -40
=== RUN   TestBuiltinSchemas_Complete
--- PASS: TestBuiltinSchemas_Complete (0.00s)
=== RUN   TestBuiltinSchemas_WellFormed
--- PASS: TestBuiltinSchemas_WellFormed (0.00s)
=== RUN   TestBuiltinSchemas_SpotRows
--- PASS: TestBuiltinSchemas_SpotRows (0.00s)
=== RUN   TestLookupInfo_ConfigSchema
--- PASS: TestLookupInfo_ConfigSchema (0.00s)
=== RUN   TestHTTPBlocks_RegistrationAndInfo
--- PASS: TestHTTPBlocks_RegistrationAndInfo (0.00s)
=== RUN   TestMQTTRegistration_InfoAndSchema
--- PASS: TestMQTTRegistration_InfoAndSchema (0.00s)
=== RUN   TestSort_CapWithSameTimestamps
--- PASS: TestSort_CapWithSameTimestamps (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	0.066s
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	0.046s [no tests to run]
=== RUN   TestBlock_RoundTripEcho
--- PASS: TestBlock_RoundTripEcho (0.01s)
=== RUN   TestBridge_RoundTripResponseCap
=== RUN   TestBridge_RoundTripResponseCap/exactly_at_the_cap
=== RUN   TestBridge_RoundTripResponseCap/one_byte_over_the_cap
--- PASS: TestBridge_RoundTripResponseCap (0.01s)
    --- PASS: TestBridge_RoundTripResponseCap/exactly_at_the_cap (0.01s)
    --- PASS: TestBridge_RoundTripResponseCap/one_byte_over_the_cap (0.01s)
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.067s
testing: warning: no tests to run
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	0.038s [no tests to run]

$ cd /root/astrate-mule && cp internal/flow/blocks/container/httpbridge.go /tmp/opencode/httpbridge.go.bak && python3 - <<'EOF'
import re
p='internal/flow/blocks/container/httpbridge.go'
s=open(p).read()
s=s.replace("""	limited := io.LimitReader(resp.Body, b.maxBody()+1)
	respBody, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("container: read response: %w", err)
	}
	if int64(len(respBody)) > b.maxBody() {
		return nil, fmt.Errorf("container: response body exceeds %d bytes", b.maxBody())
	}
""","""	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("container: read response: %w", err)
	}
""")
open(p,'w').write(s)
EOF
go test ./internal/flow/blocks/container/ -run TestBridge_RoundTripResponseCap 2>&1 | tail -12; cp /tmp/opencode/httpbridge.go.bak internal/flow/blocks/container/httpbridge.go && git diff --stat
--- FAIL: TestBridge_RoundTripResponseCap (0.01s)
    --- FAIL: TestBridge_RoundTripResponseCap/one_byte_over_the_cap (0.00s)
        block_test.go:523: RoundTrip with a 4097-byte answer: err = <nil>, want the "exceeds" error
FAIL
FAIL	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.026s
FAIL
 internal/flow/blocks/container/block_test.go | 74 ++++++++++++++++++++++++++++
 internal/flow/blocks/info.go                 |  2 +-
 internal/flow/blocks/schema.go               |  2 +-
 3 files changed, 76 insertions(+), 2 deletions(-)

$ cd /root/astrate-mule && gofmt -l . ; go vet ./... && go test ./... 2>&1 | grep -v "^ok\|no test files" | head -20; echo "EXIT_OK"
EXIT_OK

```
Done: Pinned the container response cap — a 4096-byte answer round-trips byte-exact through `Bridge{MaxBodyBytes: 4096}`, 4097 bytes is rejected with "exceeds"; the1 MiB limit is now named in the container catalog entry and schema.
Files: internal/flow/blocks/container/block_test.go, internal/flow/blocks/info.go:116, internal/flow/blocks/schema.go:19
Verified: gofmt -l . (clean), go vet ./..., go test ./... -> pass; mutation check (LimitReader + exceeds check deleted) makes TestBridge_RoundTripResponseCap/one_byte_over_the_cap fail -> pass
Unsure: nothing. The cap is a hardcoded default (defaultMaxBodyBytes, httpbridge.go:20) with no config key, so the doc lines state the fixed 1 MiB rather than offering a knob — MaxBodyBytes is settable only by code, not by a pipeline author.
Follow-ups: none
```
