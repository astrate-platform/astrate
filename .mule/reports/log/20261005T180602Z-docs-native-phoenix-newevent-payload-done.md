slug: docs-native-phoenix-newevent-payload
verdict: done
at:  a16b5fe
ran: 2026-10-05T18:06:02Z on DietPi in 267s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ ls /root/astrate-mule/docs; ls /root/astrate-mule/docs/Makefile 2>/dev/null && echo "--- makefile exists"
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
UPSTREAM-EXPERIMENTAL.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui
/root/astrate-mule/docs/Makefile
--- makefile exists

→ Read docs/api/astrate_native_api.yaml [offset=260, limit=120]

$ cat /root/astrate-mule/docs/Makefile
.PHONY: serve build clean install sync

install:
	pip install -r requirements.txt

sync:
	cp DESIGN.md site/DESIGN.md
	cp ROADMAP.md site/ROADMAP.md
	cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
	rm -rf site/api && cp -r api site/api
	rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui

serve: sync
	mkdocs serve -f mkdocs.yml

build: sync
	mkdocs build -f mkdocs.yml

clean:
	rm -rf ../site-dist

→ Read internal/engine/stream/bus.go [limit=80, offset=25]
→ Read internal/appengine/channels/ws.go [limit=40, offset=390]
→ Read internal/appengine/stream/ws.go [limit=25, offset=45]

$ sed -n '1,60p' /root/astrate-mule/internal/appengine/channels/protocol.go
// Package channels implements the upstream Phoenix Channels V2 wire protocol
// served at /appengine/v1/socket/websocket. It is the compatibility counterpart
// to the Astrate-native socket in internal/appengine/stream.
package channels

import (
	"encoding/json"
	"fmt"
)

// Event-name constants used by the Phoenix V2 protocol.
const (
	EventPhxJoin   = "phx_join"
	EventPhxReply  = "phx_reply"
	EventPhxLeave  = "phx_leave"
	EventPhxClose  = "phx_close"
	EventPhxError  = "phx_error"
	EventHeartbeat = "heartbeat"

	// TopicHeartbeat is the reserved topic for heartbeat messages.
	TopicHeartbeat = "phoenix"
)

// Frame is one Phoenix V2 message on the wire, encoded as a five-element JSON
// array: [join_ref, ref, topic, event, payload]. A nil JoinRef or Ref is
// marshalled as JSON null; the two are distinguished from the empty string.
type Frame struct {
	JoinRef *string
	Ref     *string
	Topic   string
	Event   string
	Payload json.RawMessage
}

// MarshalJSON renders the frame as a five-element JSON array. A nil JoinRef or
// Ref becomes JSON null. A nil or zero-length Payload becomes {}.
func (f Frame) MarshalJSON() ([]byte, error) {
	payload := json.RawMessage(`{}`)
	if len(f.Payload) > 0 {
		payload = f.Payload
	}
	return json.Marshal([5]interface{}{f.JoinRef, f.Ref, f.Topic, f.Event, payload})
}

// UnmarshalJSON parses a five-element JSON array back into the frame. It
// rejects anything that is not a JSON array of exactly five elements.
func (f *Frame) UnmarshalJSON(b []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return fmt.Errorf("channels.Frame: invalid JSON: %w", err)
	}
	if len(raw) != 5 {
		return fmt.Errorf("channels.Frame: expected a 5-element JSON array, got %d elements", len(raw))
	}
	// Decode join_ref and ref as nullable strings.
	jr, err := unmarshalNullString(raw[0])
	if err != nil {
		return fmt.Errorf("channels.Frame: join_ref: %w", err)
	}
	f.JoinRef = jr

$ cd /root/astrate-mule && rg -n "Frame\b|writeFrame" internal/appengine/channels/protocol.go | head -30
24:// Frame is one Phoenix V2 message on the wire, encoded as a five-element JSON
27:type Frame struct {
37:func (f Frame) MarshalJSON() ([]byte, error) {
47:func (f *Frame) UnmarshalJSON(b []byte) error {
50:		return fmt.Errorf("channels.Frame: invalid JSON: %w", err)
53:		return fmt.Errorf("channels.Frame: expected a 5-element JSON array, got %d elements", len(raw))
58:		return fmt.Errorf("channels.Frame: join_ref: %w", err)
63:		return fmt.Errorf("channels.Frame: ref: %w", err)
67:		return fmt.Errorf("channels.Frame: topic: %w", err)
70:		return fmt.Errorf("channels.Frame: event: %w", err)
93:func Reply(in Frame, status string, response any) (Frame, error) {
98:			return Frame{}, fmt.Errorf("channels.Reply: marshal response: %w", err)
108:		return Frame{}, fmt.Errorf("channels.Reply: marshal payload: %w", err)
110:	return Frame{
120:func OK(in Frame, response any) (Frame, error) {
126:func Err(in Frame, reason string) (Frame, error) {

$ cd /root/astrate-mule && rg -rn "event_data" --glob '!.git' . | head -20; echo "--- exit $?"
./docs/api/astrate_native_api.yaml:        [null, null, "rooms:test:all", "new_event", {"n": "..."}]
--- exit 0

$ cd /root/astrate-mule && head -3 go.mod && go version && rg -n "channels" internal/appengine/channels/ws.go | head -5
module github.com/astrate-platform/astrate

go 1.26.1
go version go1.26.7 linux/arm64
1:package channels
129:	// the nested shape before authorization is ever consulted (channels.json,
139:	// with this exact reason (channels.json, 2026-08-22 recording) — the
154:// test/conformance/upstream/channels.json — device shapes on 2026-07-26, the
345:		// a reason string — recorded verbatim in channels.json.

$ cd /root/astrate-mule && rg -n "event_data" --glob '!.git' . ; echo "--- done"
./docs/api/astrate_native_api.yaml:365:        [null, null, "rooms:test:all", "new_event", {"event_data": "..."}]
--- done

$ mkdir -p /tmp/opencode/newframe && cd /tmp/opencode/newframe && cat > go.mod <<'EOF'
module probe

go 1.26.1

require github.com/astrate-platform/astrate v0.0.0

replace github.com/astrate-platform/astrate => /root/astrate-mule
EOF
echo ok
ok

← Write /tmp/opencode/newframe/main.go
Wrote file successfully.

$ cd /tmp/opencode/newframe && GOFLAGS=-mod=mod go mod tidy >/dev/null 2>&1; go run . 2>&1 | tail -20
package probe
	main.go:8:2: use of internal package github.com/astrate-platform/astrate/internal/appengine/channels not allowed

← Write /tmp/opencode/newframe/main.go
Wrote file successfully.

$ cd /tmp/opencode/newframe && rm -f go.mod && cat > go.mod <<'EOF'
module probe

go 1.26.1
EOF
go run . 2>&1
payload:    {"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}
full frame: [null,null,"rooms:test:all","new_event",{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}]
event_data client: err=<nil> value=map[]
snake_case client: err=<nil> event="" device_id="" value=23.5 ts=false

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -361,10 +361,18 @@
   - Server pushes `new_event` frames for matching events
 
   Frame format:
   ```json
-  [null, null, "rooms:test:all", "new_event", {"event_data": "..."}]
+  [null, null, "rooms:test:all", "new_event", {"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}]
   ```
+
+  The payload is the bus event marshalled straight into the frame, so
+  its keys are the capitalised Go field names — unlike
+  `/astrate/v1/{realm}/socket` above, whose example is snake_case
+  because that socket serialises through its own `wireEvent` struct.
+  A client generated from any other shape (snake_case, `event_data`)
+  reads nothing at all here: Go's unmarshal leaves unknown fields zero
+  rather than failing, so every field arrives empty with no error.
 tags: [WebSocket]
 parameters:
   - name: realm
     in: query

→ Read docs/api/astrate_native_api.yaml [limit=32, offset=347]

$ cd /root/astrate-mule && python3 - <<'EOF'
import json, yaml
d = yaml.safe_load(open('docs/api/astrate_native_api.yaml'))
op = d['paths']['/appengine/v1/socket/websocket']['get']
desc = op['description']
print("YAML OK; operationId =", op['operationId'])
# extract the json fenced block and validate it parses
block = desc.split('```json\n')[1].split('\n```')[0].strip()
frame = json.loads(block)
print("frame elements:", len(frame))
print("payload keys:", list(frame[4].keys()))
EOF
Traceback (most recent call last):
  File "<stdin>", line 1, in <module>
ModuleNotFoundError: No module named 'yaml'

$ cd /root/astrate-mule && gofmt -l . && echo "gofmt clean" && go vet ./... 2>&1 | tail -5 && echo "vet done"
gofmt clean
vet done

$ cd /root/astrate-mule && (python3 -c "import yaml" 2>/dev/null && echo have-yaml) ; pip3 install pyyaml 2>&1 | tail -2; which mkdocs; mkdocs --version 2>&1 | tail -1
/bin/sh: 1: pip3: not found
/bin/sh: 1: mkdocs: not found

$ cd /root/astrate-mule && python3 -m pip --version 2>&1 | tail -1; ls docs/site 2>/dev/null | head; ls -d .venv venv docs/.venv 2>/dev/null
/usr/bin/python3: No module named pip
DESIGN.md
ROADMAP.md
api
appengine-api.md
architecture.md
compatibility.md
configuration-reference.md
contributing.md
data-modeling.md
deployment.md

$ cd /root/astrate-mule && make -C docs build 2>&1 | tail -15; echo "=== exit ${PIPESTATUS[0]}"
make: Entering directory '/root/astrate-mule/docs'
cp DESIGN.md site/DESIGN.md
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'
/bin/sh: 1: Bad substitution

$ cd /root/astrate-mule && git status --short && echo "--- yaml deps:" && rg -n "yaml" go.mod
M docs/api/astrate_native_api.yaml
--- yaml deps:
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

$ cd /root/astrate-mule && ls *.go 2>/dev/null; ls tools/ 2>/dev/null | head
bsoncapture
mule-plan-day.sh
mule-setup-pi.sh
mule-survey.sh
mule.sh
reconcile.sh

$ mkdir -p /root/astrate-mule/.mule/tmpprobe && cat > /root/astrate-mule/.mule/tmpprobe/main.go <<'EOF'
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	b, err := os.ReadFile("docs/api/astrate_native_api.yaml")
	if err != nil {
		panic(err)
	}
	var doc map[string]any
	if err := yaml.Unmarshal(b, &doc); err != nil {
		fmt.Println("YAML PARSE FAIL:", err)
		os.Exit(1)
	}
	fmt.Println("YAML parses OK; top-level keys:", len(doc))
	paths := doc["paths"].(map[string]any)
	p := paths["/appengine/v1/socket/websocket"].(map[string]any)["get"].(map[string]any)
	desc := p["description"].(string)
	block := strings.Split(desc, "```json\n")[1]
	block = strings.Split(block, "\n```")[0]
	block = strings.TrimSpace(block)
	frame := []any{}
	if err := json.Unmarshal([]byte(block), &frame); err != nil {
		fmt.Println("JSON EXAMPLE FAIL:", err)
		os.Exit(1)
	}
	fmt.Println("json example parses; frame len =", len(frame))
	payload := frame[4].(map[string]any)
	for _, k := range []string{"Kind", "Realm", "DeviceID", "Interface", "Path", "InterfaceMajor", "IP", "ErrorName", "ErrorMetadata", "Value", "Timestamp"} {
		if _, ok := payload[k]; !ok {
			fmt.Println("MISSING payload key:", k)
			os.Exit(1)
		}
	}
	fmt.Println("payload keys =", len(payload), "(all expected present)")
	// the native socket example must still be snake_case and unchanged
	np := paths["/astrate/v1/{realm}/socket"].(map[string]any)["get"].(map[string]any)["description"].(string)
	nb := strings.Split(strings.Split(np, "```json\n")[1], "\n```")[0]
	for _, want := range []string{`"event":`, `"device_id":`, `"timestamp":`} {
		if !strings.Contains(nb, want) {
			fmt.Println("native example missing", want)
			os.Exit(1)
		}
	}
	fmt.Println("native socket example still snake_case: OK")
}
EOF
cd /root/astrate-mule && go run ./.mule/tmpprobe 2>&1
YAML parses OK; top-level keys: 6
json example parses; frame len = 5
payload keys = 11 (all expected present)
native socket example still snake_case: OK

$ rm -rf /root/astrate-mule/.mule/tmpprobe && cd /root/astrate-mule && git status --short
M docs/api/astrate_native_api.yaml

$ cd /root/astrate-mule && go test ./... 2>&1 | rg -v "^ok|no test files" | head -20; echo "=== go test done"
=== go test done

$ cd /root/astrate-mule && go test -race ./internal/appengine/channels/... 2>&1 | tail -5
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/appengine/channels	0.023s
FAIL

$ cd /root/astrate-mule && git diff
diff --git a/docs/api/astrate_native_api.yaml b/docs/api/astrate_native_api.yaml
index 853e0c6..58f62e8 100644
--- a/docs/api/astrate_native_api.yaml
+++ b/docs/api/astrate_native_api.yaml
@@ -362,8 +362,16 @@ paths:
 
         Frame format:
         ```json
-        [null, null, "rooms:test:all", "new_event", {"event_data": "..."}]
+        [null, null, "rooms:test:all", "new_event", {"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}]
         ```
+
+        The payload is the bus event marshalled straight into the frame, so
+        its keys are the capitalised Go field names — unlike
+        `/astrate/v1/{realm}/socket` above, whose example is snake_case
+        because that socket serialises through its own `wireEvent` struct.
+        A client generated from any other shape (snake_case, `event_data`)
+        reads nothing at all here: Go's unmarshal leaves unknown fields zero
+        rather than failing, so every field arrives empty with no error.
       tags: [WebSocket]
       parameters:
         - name: realm

Done: Replaced the `new_event` payload example in `docs/api/astrate_native_api.yaml` with the measured `stream.Event` shape (capitalised Go field names) and documented why it differs from the snake_case native socket; no behaviour change, so no test written.

Files: docs/api/astrate_native_api.yaml (frame example at :363-366, new prose :368-374)

Verified: probe on a verbatim copy of `stream.Event` + `channels.Frame` (`/tmp/opencode/newframe`) -> payload and full 5-element frame byte-identical to the task's stated measurement; `event_data` client unmarshals with `err=<nil>` and empty map; snake_case client gets zero values; `yaml.v3` parse of the whole spec -> OK, the JSON block parses as a 5-element frame with all 11 capitalised keys, and the `/astrate/v1/{realm}/socket` example still snake_case; `rg event_data` -> this YAML was the only occurrence in the repo; `go vet ./...` -> pass; `go test ./...` -> pass; `gofmt -l .` -> clean; `make -C docs build` -> **the `sync` half passed, `mkdocs build` failed with `mkdocs: No such file or directory`** (no mkdocs and no pip on this box — environment gap, not a content failure); `go test -race` -> unusable here as documented, `ThreadSanitizer: unsupported VMA range (39-bit)`.

Unsure: two wording points in the task line. It says "the four-element frame shape" but the example array has five elements (join_ref, ref, topic, event, payload) — I left the array structure untouched as instructed and pinned nothing new. And "swagger-referenced YAML still loads" was verified with a yaml.v3 parse plus a JSON-parse of the example, not with mkdocs/swagger itself, since neither is installed. I also called the serialiser "the bus event marshalled straight into the frame" rather than naming `stream.Event`/`wireEvent` files in the prose, to keep the public API doc free of internal paths — say the word if you want the citations in.

Follow-ups: none
