slug: examples-echo-container-contract-test
verdict: done
at:  f35af21
ran: 2026-09-23T17:07:27Z on DietPi in 251s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read examples/flow-container-echo/main.go
✱ Glob "examples/flow-container-echo/**" 3 matches
→ Read examples/flow-container-echo/README.md

$ ls /root/astrate-mule && head -5 /root/astrate-mule/go.mod
AGENTS.md
Dockerfile
LICENSE
Makefile
README.md
bench
cmd
deploy
docker-compose.yml
docs
examples
go.mod
go.sum
internal
migrations
opencode.json
pkg
site-dist
test
tools
module github.com/astrate-platform/astrate

go 1.26.1

// Pinned dependency decisions — docs/ROADMAP.md §1.1 (frozen in docs/DESIGN.md).

$ ls /root/astrate-mule/examples && ls /root/astrate-mule/examples/*/
flow-container-echo
pokemon-agent
/root/astrate-mule/examples/flow-container-echo/:
Dockerfile
README.md
main.go

/root/astrate-mule/examples/pokemon-agent/:
Makefile
README.md
astrate-interfaces
docs
emulator-agent
llm-orchestrator

✱ Grep "1016 / 1<<20|1<<20|1 MiB|MiB" in . · 7 matches

$ ls /root/astrate-mule/docs/handoff/ | head -30
README.md
flow-design-a-named-durable-flows-2026-07-29.md
flow-design-b-container-block-2026-07-29.md
flow-lifecycle-followups-design-2026-08-22.md
flow-parity-audit-2026-07-29.md
flow-v2-decisions-2026-07-29.md
phase-1-handoff.md
phase-1-memory.md
phase-10-handoff.md
phase-10-memory.md
phase-11-handoff.md
phase-2-handoff.md
phase-2-memory.md
phase-3-handoff.md
phase-3-memory.md
phase-4-handoff.md
phase-4-memory.md
phase-5-handoff.md
phase-5-memory.md
phase-6-handoff.md
phase-6-memory.md
phase-7-memory.md
phase-8-handoff.md
phase-8-memory.md
phase-9-handoff.md
phase-9-memory.md
pokemon-agent-handoff.md
pokemon-agent-memory.md
session-2026-07-29-flow-factory-handoff.md
session-2026-07-29-handoff.md

→ Read internal/flow/blocks/container/httpbridge.go
→ Read docs/handoff/flow-design-b-container-block-2026-07-29.md

$ go version && cd /root/astrate-mule && go build ./examples/flow-container-echo/ && echo BUILD_OK
go version go1.26.7 linux/arm64
BUILD_OK

← Write examples/flow-container-echo/main_test.go
Wrote file successfully.

$ gofmt -l . && go vet ./examples/flow-container-echo/ && go test ./examples/flow-container-echo/ -v 2>&1 | tail -40
examples/flow-container-echo/main_test.go
=== RUN   TestHandleMessageContract
=== RUN   TestHandleMessageContract/echoes_valid_JSON_object_verbatim
=== RUN   TestHandleMessageContract/echoes_pretty-printed_JSON_byte-for-byte
=== RUN   TestHandleMessageContract/400_on_invalid_JSON
=== RUN   TestHandleMessageContract/400_on_non-object_JSON
=== RUN   TestHandleMessageContract/204_on_empty_body
=== RUN   TestHandleMessageContract/204_when_metadata_echo_drop_is_1
=== RUN   TestHandleMessageContract/echoes_when_echo_drop_is_not_the_string_1
=== RUN   TestHandleMessageContract/echoes_when_echo_drop_is_a_non-string_value
--- PASS: TestHandleMessageContract (0.00s)
    --- PASS: TestHandleMessageContract/echoes_valid_JSON_object_verbatim (0.00s)
    --- PASS: TestHandleMessageContract/echoes_pretty-printed_JSON_byte-for-byte (0.00s)
    --- PASS: TestHandleMessageContract/400_on_invalid_JSON (0.00s)
    --- PASS: TestHandleMessageContract/400_on_non-object_JSON (0.00s)
    --- PASS: TestHandleMessageContract/204_on_empty_body (0.00s)
    --- PASS: TestHandleMessageContract/204_when_metadata_echo_drop_is_1 (0.00s)
    --- PASS: TestHandleMessageContract/echoes_when_echo_drop_is_not_the_string_1 (0.00s)
    --- PASS: TestHandleMessageContract/echoes_when_echo_drop_is_a_non-string_value (0.00s)
=== RUN   TestHandleMessageReadCap
--- PASS: TestHandleMessageReadCap (0.07s)
=== RUN   TestSanitizeLog
    main_test.go:129: sanitizeLog("ä\x1a¸") = "ä¸", want "龸"
--- FAIL: TestSanitizeLog (0.00s)
FAIL
FAIL	github.com/astrate-platform/astrate/examples/flow-container-echo	0.092s
FAIL

$ gofmt -d examples/flow-container-echo/main_test.go | head -30
diff examples/flow-container-echo/main_test.go.orig examples/flow-container-echo/main_test.go
--- examples/flow-container-echo/main_test.go.orig
+++ examples/flow-container-echo/main_test.go
@@ -26,8 +26,8 @@
 			wantBody: []byte(`{"schema":"astarte_flow/message/v0.1","key":"demo","type":"string","data":"hello","timestamp_us":0}`),
 		},
 		{
-			name: "echoes pretty-printed JSON byte-for-byte",
-			body: []byte("{\n  \"schema\": \"astarte_flow/message/v0.1\",\n  \"key\": \"demo\",\n  \"data\": 1\n}\n"),
+			name:     "echoes pretty-printed JSON byte-for-byte",
+			body:     []byte("{\n  \"schema\": \"astarte_flow/message/v0.1\",\n  \"key\": \"demo\",\n  \"data\": 1\n}\n"),
 			wantCode: http.StatusOK,
 			wantBody: []byte("{\n  \"schema\": \"astarte_flow/message/v0.1\",\n  \"key\": \"demo\",\n  \"data\": 1\n}\n"),
 		},
@@ -47,8 +47,8 @@
 			wantCode: http.StatusNoContent,
 		},
 		{
-			name: "204 when metadata echo_drop is 1",
-			body: []byte(`{"schema":"astarte_flow/message/v0.1","key":"demo","data":"x","metadata":{"echo_drop":"1"}}`),
+			name:     "204 when metadata echo_drop is 1",
+			body:     []byte(`{"schema":"astarte_flow/message/v0.1","key":"demo","data":"x","metadata":{"echo_drop":"1"}}`),
 			wantCode: http.StatusNoContent,
 		},
 		{
@@ -129,4 +129,4 @@
 			t.Errorf("sanitizeLog(%q) = %q, want %q", tc.in, got, tc.want)
 		}
 	}
-}

← Edit examples/flow-container-echo/main_test.go
Index: /root/astrate-mule/examples/flow-container-echo/main_test.go
===================================================================
--- /root/astrate-mule/examples/flow-container-echo/main_test.go
+++ /root/astrate-mule/examples/flow-container-echo/main_test.go
@@ -121,9 +121,9 @@
 	{"a\x00b", "ab"},
 	{"a\x1fb", "ab"},
 	{"a\x7fb", "ab"},
 	{"a\r\nb", "ab"},
-	{string([]rune{0xe4, 0x1a, 0xb8}), "龸"}, // printable kept, C0 stripped
+	{string([]rune{0xe4, 0x1a, 0xb8}), "ä¸"}, // printable kept, C0 stripped
 }
 for _, tc := range cases {
 	if got := sanitizeLog(tc.in); got != tc.want {
 		t.Errorf("sanitizeLog(%q) = %q, want %q", tc.in, got, tc.want)

$ gofmt -w examples/flow-container-echo/main_test.go && gofmt -l . && go vet ./examples/flow-container-echo/ && go test ./examples/flow-container-echo/ -count=1 && echo ALL_PASS
ok  	github.com/astrate-platform/astrate/examples/flow-container-echo	0.098s
ALL_PASS

All green. Now running the full mule gate (`go vet ./... && go test ./...`):
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

Done: added a container-free httptest suite for examples/flow-container-echo pinning the full message contract (verbatim echo, 400/204 paths, 1 MiB read cap, echo_drop filter) plus sanitizeLog.
Files: examples/flow-container-echo/main_test.go
Verified: gofmt -l . -> clean; go vet ./... -> pass; go test ./... -> pass (plain go test, not -race: race detector unavailable on this machine per AGENTS.md)
Unsure: nothing. No behavior change — this pins existing contract; the negative echo_drop cases ("0", non-string) and 400-on-non-object assert exact-code details beyond the task's literal list, in line with the documented "one JSON object" contract.
Follow-ups: none
