slug: docs-native-socket-sse-exact-accept
verdict: done
at:  1322214
ran: 2026-10-06T20:11:12Z on DietPi in 410s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astrate_native_api.yaml [offset=260, limit=120]
→ Read internal/appengine/stream/ws.go [offset=100, limit=90]
→ Read internal/appengine/stream/ws.go [limit=99]

$ rg -n "426|event-stream|transport" docs/api/astrate_native_api.yaml
291:        `Accept: text/event-stream` or `?transport=sse`.
328:        - name: transport
334:          description: Force SSE transport instead of WebSocket.
339:          description: "SSE stream started (when `Accept: text/event-stream`)."

$ rg -rn "426" internal/appengine/stream/ docs/ | head -30
(no output)

$ rg -rn "astrate_native_api" --glob '!docs/api/*' | head -20; ls docs
docs/swagger-ui/index.html:      <option value="../api/n.yaml">Astrate Native API — Health, readiness, metrics, WebSocket</option>
docs/handoff/phase-2-handoff.md:- Fix YAML: n.yaml ha 2 valori quotati (backtick+colon)
docs/handoff/phase-2-memory.md:Due fix di indentazione YAML in `n.yaml` (problemi di parser Redocly con valori non quotati contenenti `:` dopo backtick):
docs/handoff/phase-2-memory.md:| `docs/api/n.yaml` | **Modificato** — 2 fix quotatura valori |
docs/handoff/phase-1-memory.md:| `n.yaml` | 8 | 336 | Health/readiness/metrics (no auth), 3 compat health, 2 WebSocket (native stream + Phoenix V2) |
docs/handoff/phase-1-memory.md:└── n.yaml          (8 op, 336 righe)
internal/swagger/swagger_test.go:	b, err := docs.APIYAML.ReadFile("api/n.yaml")
internal/swagger/swagger_test.go:		t.Fatalf("reading n.yaml: %v", err)
internal/swagger/swagger_test.go:	b, err := docs.APIYAML.ReadFile("api/n.yaml")
internal/swagger/swagger_test.go:		t.Fatalf("reading n.yaml: %v", err)
internal/swagger/swagger_test.go:	b, err := docs.APIYAML.ReadFile("api/n.yaml")
internal/swagger/swagger_test.go:		t.Fatalf("reading n.yaml: %v", err)
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

$ rg -n "openapi:|headers:|text/event-stream|content:" docs/api/astrate_native_api.yaml | head -40
1:openapi: "3.0.3"
37:          content:
56:          content:
67:          content:
99:          content:
125:          content:
134:          content:
153:          content:
162:          content:
181:          content:
190:          content:
213:          content:
233:          content:
253:          content:
275:          content:
291:        `Accept: text/event-stream` or `?transport=sse`.
339:          description: "SSE stream started (when `Accept: text/event-stream`)."
483:      content:
493:      content:
503:      content:

$ rg -n "\"200\":" -A 8 docs/api/astrate_native_api.yaml | sed -n '1,60p'
35:        "200":
36-          description: Service is alive.
37-          content:
38-            application/json:
39-              schema:
40-                $ref: "#/components/schemas/HealthStatus"
41-              example:
42-                status: ok
43-
--
54:        "200":
55-          description: All dependencies are ready.
56-          content:
57-            application/json:
58-              schema:
59-                $ref: "#/components/schemas/ReadinessStatus"
60-              example:
61-                status: ok
62-                checks:
--
97:        "200":
98-          description: Prometheus metrics.
99-          content:
100-            text/plain; version=0.0.4; charset=utf-8; escaping=underscores:
101-              schema:
102-                type: string
103-              example: |
104-                # HELP astrate_broker_sessions Live authenticated MQTT device sessions.
105-                # TYPE astrate_broker_sessions gauge
--
123:        "200":
124-          description: Service is healthy.
125-          content:
126-            application/json:
127-              schema:
128-                $ref: "#/components/schemas/DataEnvelopeHealthStatus"
129-              example:
130-                data:
131-                  status: ok
--
151:        "200":
152-          description: Service is healthy.
153-          content:
154-            application/json:
155-              schema:
156-                $ref: "#/components/schemas/DataEnvelopeHealthStatus"
157-              example:
158-                data:
159-                  status: ok
--
179:        "200":
180-          description: Service is healthy.
181-          content:
182-            application/json:
183-              schema:
184-                $ref: "#/components/schemas/DataEnvelopeHealthStatus"
185-              example:
186-                data:
187-                  status: ok
--

→ Read internal/swagger/swagger_test.go
→ Read docs/Makefile

$ rg -n "description: [>|]" docs/api/astrate_native_api.yaml; rg -n "description: >-" docs/api/*.yaml | head -10
4:  description: |
30:      description: |
48:      description: |
81:      description: |
118:      description: |
146:      description: |
174:      description: |
204:      description: |
224:      description: |
244:      description: |
264:      description: |
288:      description: |
353:      description: |
380:          description: |
412:      description: |

→ Read docs/api/astrate_native_api.yaml [limit=30, offset=376]

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -286,10 +286,16 @@
 operationId: nativeWebSocket
 summary: Native live event stream (WebSocket / SSE)
 description: |
   Astrate-native real-time event stream. Upgrades to WebSocket by
-  default, or streams Server-Sent Events when the client sends
-  `Accept: text/event-stream` or `?transport=sse`.
+  default, or streams Server-Sent Events when the request carries an
+  `Accept` header whose value is exactly `text/event-stream`, or when
+  `?transport=sse`. The header is matched for exact string equality: a
+  list or a parameterised value such as `text/event-stream, text/plain`
+  or `text/event-stream;q=1.0` does not match, and neither does a
+  missing header. Any other `Accept` falls through to the WebSocket
+  upgrade attempt, which answers a client that is not upgrading with
+  `426 Upgrade Required`.
 
   Subscribes to the realm's committed data and device lifecycle events.
   Optional `device_id` and `interface` query parameters narrow the
   event filter.

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -336,9 +336,17 @@
     schema:
       type: string
       enum: [websocket, sse]
       default: websocket
-    description: Force SSE transport instead of WebSocket.
+    description: |
+      Requested transport. Server-Sent Events are selected only by the
+      literal value `sse`; `websocket`, the empty value and every other
+      string mean WebSocket, i.e. they fall through to the upgrade
+      attempt (an exactly matching `Accept: text/event-stream` header
+      selects SSE on its own and overrides this parameter). The `enum`
+      is a client-side contract: the handler never rejects a value
+      outside it — a non-`sse` value simply does not match and is not
+      reported as an error.
 responses:
   "101":
     description: WebSocket upgrade successful.
   "200":

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -349,9 +349,22 @@
 responses:
   "101":
     description: WebSocket upgrade successful.
   "200":
-    description: "SSE stream started (when `Accept: text/event-stream`)."
+    description: |
+      SSE stream started, when the `Accept` header is exactly
+      `text/event-stream` or `transport=sse`. The handler sets
+      `Content-Type: text/event-stream`, `Cache-Control: no-cache` and
+      `Connection: keep-alive`, then writes one `data: {...}` frame per
+      event followed by a blank line, flushing after each frame.
+    content:
+      text/event-stream:
+        schema:
+          type: string
+        example: |-
+          data: {"event":"incoming_data","realm":"test","device_id":"f0VMRgIBAQAAAAAAAAAAAA","interface":"org.example.Sensors","path":"/sensors/temperature/value","value":23.5,"timestamp":"2024-05-30T13:49:57.045Z"}
+
+          data: {"event":"device_lifecycle","realm":"test","device_id":"f0VMRgIBAQAAAAAAAAAAAA","timestamp":"2024-05-30T13:49:58.010Z"}
   "401":
     $ref: "#/components/responses/Unauthorized"
   "403":
     $ref: "#/components/responses/Forbidden"

$ rg -n "Kind:|Kind =" internal/engine/ | head -20
internal/engine/data.go:74:	OpIndividual OpKind = iota + 1
internal/engine/data.go:324:		op.Kind = OpPropertyUnset
internal/engine/data.go:326:		op.Kind = OpPropertySet
internal/engine/data.go:328:		op.Kind = OpObject
internal/engine/data.go:330:		op.Kind = OpIndividual
internal/engine/data.go:362:	if op.Kind == OpObject {
internal/engine/data.go:455:	if op.Kind == OpObject {
internal/engine/topics.go:22:	kindIntrospection topicKind = iota + 1
internal/engine/stream/bus_test.go:15:		Kind: KindIncomingData, Realm: realm, DeviceID: device,
internal/engine/stream/bus_test.go:92:		b.Publish(Event{Kind: KindDeviceConnected, Realm: "r", DeviceID: "d"})
internal/engine/stream/bus_test.go:101:		b.Publish(Event{Kind: KindDeviceConnected, Realm: "r", DeviceID: "d"})
internal/engine/engine.go:182:		if op.Kind == OpIndividual {
internal/engine/engine.go:208:		Kind:           stream.KindIncomingData,
internal/engine/engine.go:310:		Kind:          stream.KindDeviceError,
internal/engine/engine.go:351:			Kind: stream.KindDeviceConnected, Realm: rs.name, DeviceID: deviceID, Timestamp: ev.At, IP: ip,
internal/engine/engine.go:369:			Kind: stream.KindDeviceDisconnected, Realm: rs.name, DeviceID: deviceID, Timestamp: ev.At,
internal/engine/engine.go:390:		Kind: stream.KindDeviceRegistered, Realm: rs.name, DeviceID: hwID, Timestamp: at,
internal/engine/engine.go:424:		Kind: kind, Realm: rs.name, DeviceID: hwID, Timestamp: at,
internal/engine/batch_test.go:93:			Kind: OpIndividual, RealmID: 1, DeviceID: devAlpha, Interface: ci,
internal/engine/batch_test.go:186:		Kind: OpObject, RealmID: 1, DeviceID: devAlpha, Interface: flat, Path: "",

$ rg -n "Kind\w+ |= \"" internal/engine/stream/*.go | head -20
internal/engine/stream/bus.go:19:	// KindIncomingData is a committed device data operation (set or unset).
internal/engine/stream/bus.go:20:	KindIncomingData = "incoming_data"
internal/engine/stream/bus.go:21:	// KindDeviceRegistered is a device registration via the Pairing API.
internal/engine/stream/bus.go:22:	KindDeviceRegistered = "device_registered"
internal/engine/stream/bus.go:23:	// KindDeviceConnected is a device connection.
internal/engine/stream/bus.go:24:	KindDeviceConnected = "device_connected"
internal/engine/stream/bus.go:25:	// KindDeviceDisconnected is a device disconnection.
internal/engine/stream/bus.go:26:	KindDeviceDisconnected = "device_disconnected"
internal/engine/stream/bus.go:27:	// KindDeviceDeletionStarted is emitted immediately before a synchronous
internal/engine/stream/bus.go:29:	KindDeviceDeletionStarted = "device_deletion_started"
internal/engine/stream/bus.go:30:	// KindDeviceDeletionFinished is emitted immediately after a synchronous
internal/engine/stream/bus.go:32:	KindDeviceDeletionFinished = "device_deletion_finished"
internal/engine/stream/bus.go:33:	// KindDeviceError is a rejected device message (docs/DESIGN.md §2.6
internal/engine/stream/bus.go:36:	KindDeviceError = "device_error"
internal/engine/stream/bus.go:78:	if f.DeviceID != "" && f.DeviceID != ev.DeviceID {
internal/engine/stream/bus.go:81:	if f.Interface != "" && f.Interface != ev.Interface {
internal/engine/stream/bus_test.go:52:	if ev := recv(t, ch); ev.Realm != "realmA" || ev.DeviceID != "dev1" {
internal/engine/stream/bus_test.go:68:		if ev := recv(t, ch); ev.DeviceID != "keep" {
internal/engine/stream/bus_test.go:81:		if ev := recv(t, ch); ev.Interface != "com.ex.Keep" {
internal/engine/stream/bus_test.go:102:		if ev := recv(t, ch); ev.Kind != KindDeviceConnected {

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -362,9 +362,9 @@
         type: string
       example: |-
         data: {"event":"incoming_data","realm":"test","device_id":"f0VMRgIBAQAAAAAAAAAAAA","interface":"org.example.Sensors","path":"/sensors/temperature/value","value":23.5,"timestamp":"2024-05-30T13:49:57.045Z"}
 
-        data: {"event":"device_lifecycle","realm":"test","device_id":"f0VMRgIBAQAAAAAAAAAAAA","timestamp":"2024-05-30T13:49:58.010Z"}
+        data: {"event":"device_connected","realm":"test","device_id":"f0VMRgIBAQAAAAAAAAAAAA","timestamp":"2024-05-30T13:49:58.010Z"}
 "401":
   $ref: "#/components/responses/Unauthorized"
 "403":
   $ref: "#/components/responses/Forbidden"

→ Read docs/api/astrate_native_api.yaml [limit=90, offset=284]

$ make -C docs build 2>&1 | tail -20
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

$ which python3 pip3 mkdocs; python3 -c "import mkdocs, sys; print(mkdocs.__version__)" 2>&1 | tail -2; ls docs/requirements.txt && cat docs/requirements.txt
/usr/bin/python3
    ^^^^^^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'mkdocs'
docs/requirements.txt
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10

$ python3 -m pip --version 2>&1 | tail -3; python3 -c "import sys; print(sys.version)"
/usr/bin/python3: No module named pip
3.13.5 (main, Jul 15 2026, 20:25:40) [GCC 14.2.0]

$ python3 -m ensurepip --version 2>&1 | tail -2; ls /usr/lib/python3/dist-packages 2>/dev/null | head; ls docs/site | head
pip 25.1.1
Cryptodome
packaging
packaging-25.0.dist-info
pycryptodomex-3.20.0.egg-info
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

$ python3 -m venv /tmp/mkvenv && /tmp/mkvenv/bin/pip install -q -r docs/requirements.txt 2>&1 | tail -5; /tmp/mkvenv/bin/mkdocs --version
mkdocs, version 1.6.1 from /tmp/mkvenv/lib/python3.13/site-packages/mkdocs (Python 3.13)

$ PATH=/tmp/mkvenv/bin:$PATH make -C docs build 2>&1 | tail -25
mkdocs build -f mkdocs.yml

 │  ⚠  Warning from the Material for MkDocs team
 │
 │  MkDocs 2.0, the underlying framework of Material for MkDocs,
 │  will introduce backward-incompatible changes, including:
 │
 │  × All plugins will stop working – the plugin system has been removed
 │  × All theme overrides will break – the theming system has been rewritten
 │  × No migration path exists – existing projects cannot be upgraded
 │  × Closed contribution model – community members can't report bugs
 │  × Currently unlicensed – unsuitable for production use
 │
 │  Our full analysis:
 │
 │  https://squidfunk.github.io/mkdocs-material/blog/2026/02/18/mkdocs-2.0/

INFO    -  Cleaning site directory
INFO    -  Building documentation to directory: /root/astrate-mule/site-dist
INFO    -  The following pages exist in the docs directory, but are not included in the "nav" configuration:
  - DESIGN.md
  - ROADMAP.md
INFO    -  Doc file 'index.md' contains an unrecognized relative link 'api/', it was left as is.
INFO    -  Documentation built in 3.95 seconds
make: Leaving directory '/root/astrate-mule/docs'

$ ls internal/swagger/; rg -ln "yaml" internal/swagger/ | head; rg -ln "yaml.Unmarshal|openapi3|kin-openapi" --type go | head -20
swagger.go
swagger_test.go
internal/swagger/swagger.go
internal/swagger/swagger_test.go

$ /tmp/mkvenv/bin/python - <<'PY'
import yaml
d = yaml.safe_load(open('docs/api/astrate_native_api.yaml'))
op = d['paths']['/astrate/v1/{realm}/socket']['get']
print('openapi', d['openapi'], '| paths', len(d['paths']))
print('200 content keys:', list(op['responses']['200']['content'].keys()))
print('200 desc first line:', op['responses']['200']['description'].splitlines()[0])
print('transport desc has enum note:', 'client-side contract' in op['parameters'][3]['description'])
print('desc mentions exact:', 'exactly `text/event-stream`' in op['description'])
print('desc mentions 426:', '426' in op['description'])
PY
openapi 3.0.3 | paths 12
200 content keys: ['text/event-stream']
200 desc first line: SSE stream started, when the `Accept` header is exactly
transport desc has enum note: True
desc mentions exact: True
desc mentions 426: True

$ gofmt -l . ; echo "gofmt exit=$?"
gofmt exit=0

$ go vet ./... 2>&1 | tail -15; echo "vet=$?"
vet=0

$ go vet ./... > /tmp/vet.log 2>&1; echo "vet=$?"; go test -race ./... > /tmp/test.log 2>&1; echo "test=$?"; tail -25 /tmp/test.log
vet=0
test=1
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/realm	0.049s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/store	0.042s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.037s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/testutil	0.028s
?   	github.com/astrate-platform/astrate/migrations	[no test files]
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/astarteapi	0.028s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/deviceid	0.024s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/interfaceschema	0.023s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/payload	0.028s
FAIL

$ go test ./... > /tmp/test2.log 2>&1; echo "test=$?"; grep -E "^(ok|FAIL|---)" /tmp/test2.log | tail -40
test=0
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
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
ok  	github.com/astrate-platform/astrate/internal/swagger	0.130s
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

$ git status --short | head -30; echo "---"; git check-ignore -v site-dist docs/site 2>&1 | head
M docs/api/astrate_native_api.yaml
---
.gitignore:29:/site-dist/	site-dist

$ git ls-files docs/site | head -3; git status --short --untracked-files=all docs/site | head -5; ls -la docs/site/api/ 2>/dev/null | head; git check-ignore -v docs/site/api/astrate_native_api.yaml
docs/site/appengine-api.md
docs/site/architecture.md
docs/site/compatibility.md
total 180
drwxr-xr-x 2 root root  4096 Oct  6 22:08 .
drwxr-xr-x 4 root root  4096 Oct  6 22:08 ..
-rw-r--r-- 1 root root 61063 Oct  6 22:08 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 17437 Oct  6 22:08 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 21724 Oct  6 22:08 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 47623 Oct  6 22:08 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 19875 Oct  6 22:08 astrate_native_api.yaml
.gitignore:34:/docs/site/api/	docs/site/api/astrate_native_api.yaml

Done: rewrote the SSE negotiation prose on `/astrate/v1/{realm}/socket` — exact-match `Accept`, fall-through to the upgrade attempt (426 for non-upgrade clients), `transport` non-enforcing enum, and the `200` now carries `text/event-stream` with its real headers and frame format
Files: docs/api/astrate_native_api.yaml
Verified: `make -C docs build` -> pass; PyYAML parse of the spec (12 paths, `200.content` keys `['text/event-stream']`) -> pass; `go vet ./...` -> pass; `go test ./...` -> pass (incl. internal/swagger reading the embedded YAML); `gofmt -l .` -> clean. `go test -race ./...` did **not** run: ThreadSanitizer aborts with "unsupported VMA range" (39-bit VMA kernel) — this change is docs-only, no goroutines/locks touched.
Unsure: mkdocs was not installed system-wide here; I installed `docs/requirements.txt` into a throwaway venv at `/tmp/mkvenv` and put it on PATH for the build — whether the timer host has mkdocs is unknown.
Follow-ups: none
