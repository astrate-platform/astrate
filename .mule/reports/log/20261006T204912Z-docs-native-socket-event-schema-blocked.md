slug: docs-native-socket-event-schema
verdict: blocked
at:  bc3dc9a
ran: 2026-10-06T20:49:12Z on DietPi in 662s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astrate_native_api.yaml
→ Read internal/appengine/stream/ws.go [limit=120]
→ Read internal/engine/stream/bus.go [limit=80]
→ Read internal/appengine/channels/ws.go [offset=395, limit=35]

$ ls docs && cat docs/Makefile 2>/dev/null
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

$ ls docs/swagger-ui docs/api && rg -n "astrate_native_api|swagger" docs --glob '!site/**' -g '!swagger-ui/**' | head -50
docs/api:
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

docs/swagger-ui:
index.html
docs/swagger-ui/index.html:7:  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5.32.11/favicon-32x32.png">
docs/swagger-ui/index.html:8:  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui.css">
docs/swagger-ui/index.html:32:    #swagger-ui { max-width: 100%; }
docs/swagger-ui/index.html:48:      <option value="../api/astrate_native_api.yaml">Astrate Native API — Health, readiness, metrics, WebSocket</option>
docs/swagger-ui/index.html:53:  <div id="swagger-ui"></div>
docs/swagger-ui/index.html:55:  <script src="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui-bundle.js"></script>
docs/swagger-ui/index.html:56:  <script src="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui-standalone-preset.js"></script>
docs/swagger-ui/index.html:60:    const uiContainer = document.getElementById('swagger-ui');
docs/swagger-ui/index.html:66:        dom_id: '#swagger-ui',
docs/embed.go:9://go:embed swagger-ui/*
docs/handoff/phase-5-handoff.md:16:- Phase 2 (Swagger UI) is complete: ~/astrate/docs/swagger-ui/index.html.
docs/handoff/phase-7-memory.md:20:- `docs/swagger-ui/index.html` — bundled Swagger UI
docs/handoff/phase-4-handoff.md:16:- Phase 2 (Swagger UI) is complete: ~/astrate/docs/swagger-ui/index.html.
docs/handoff/phase-4-handoff.md:20:- DESIGN.md, ROADMAP.md, api/, swagger-ui/ are copied into docs/site/ for MkDocs. Run make sync to refresh.
docs/handoff/phase-2-memory.md:7:#### Step 2.1: `docs/swagger-ui/index.html`
docs/handoff/phase-2-memory.md:16:#### Step 2.2: Embed nel binary + route `/swagger`
docs/handoff/phase-2-memory.md:21:   - `SwaggerUI` (via `//go:embed swagger-ui/*`)
docs/handoff/phase-2-memory.md:26:2. **`internal/swagger/swagger.go`** — package `swagger`, funzione `Mount(mux *http.ServeMux)`:
docs/handoff/phase-2-memory.md:27:   - `GET /swagger` → redirect 302 a `/swagger/index.html`
docs/handoff/phase-2-memory.md:28:   - `GET /swagger/*` → `http.FileServer` su `fs.Sub(SwaggerUI, "swagger-ui")` con `StripPrefix`
docs/handoff/phase-2-memory.md:31:   Il path relativo `../api/*.yaml` nell'HTML risolve correttamente: da `/swagger/index.html` → `/api/*.yaml`.
docs/handoff/phase-2-memory.md:33:3. **`cmd/astrate/main.go`** — aggiunto import `swagger` + chiamata `swagger.Mount(mux)` in `mountAPIs()`, prima della riga `handler := httpx.NotFound(mux)`.
docs/handoff/phase-2-memory.md:36:- `go vet ./cmd/astrate/ ./internal/swagger/ ./docs/` → clean
docs/handoff/phase-2-memory.md:42:Due fix di indentazione YAML in `astrate_native_api.yaml` (problemi di parser Redocly con valori non quotati contenenti `:` dopo backtick):
docs/handoff/phase-2-memory.md:55:| `docs/swagger-ui/index.html` | **Nuovo** — Swagger UI statica |
docs/handoff/phase-2-memory.md:56:| `docs/embed.go` | **Nuovo** — go:embed per swagger-ui + api YAML |
docs/handoff/phase-2-memory.md:57:| `internal/swagger/swagger.go` | **Nuovo** — handler Mount + Specs() |
docs/handoff/phase-2-memory.md:58:| `cmd/astrate/main.go` | **Modificato** — import swagger + Mount(mux) |
docs/handoff/phase-2-memory.md:59:| `docs/api/astrate_native_api.yaml` | **Modificato** — 2 fix quotatura valori |
docs/handoff/phase-2-memory.md:64:- [x] `go vet ./cmd/astrate/ ./internal/swagger/ ./docs/` — clean
docs/handoff/phase-2-memory.md:65:- [x] Route testing concettuale: `/swagger` → redirect, `/swagger/index.html` → HTML, `/api/*.yaml` → YAML
docs/site/quickstart.md:68:- **API Explorer:** [Swagger UI](swagger.md) — interactive API docs (no manual server needed)
docs/site/index.md:44:- [API Explorer](swagger.md) -- interactive API docs
docs/handoff/phase-6-handoff.md:16:- Phase 2 (Swagger UI): ~/astrate/docs/swagger-ui/index.html.
docs/handoff/README.md:73:- Serve from `~/astrate/docs`, then open `http://localhost:9090/swagger-ui/`.
docs/handoff/README.md:74:- Do not serve from `~/astrate/docs/swagger-ui`; the `../api/` YAML paths will not resolve.
docs/handoff/README.md:80:- Run `make sync` before build to refresh DESIGN.md, ROADMAP.md, api/, swagger-ui/ copies.
docs/handoff/phase-1-memory.md:20:| `astrate_native_api.yaml` | 8 | 336 | Health/readiness/metrics (no auth), 3 compat health, 2 WebSocket (native stream + Phoenix V2) |
docs/handoff/phase-1-memory.md:56:└── astrate_native_api.yaml          (8 op, 336 righe)
docs/handoff/phase-1-memory.md:62:2. **Fase 2 — Swagger UI**: creare `docs/swagger-ui/index.html` con CDN refs e configurazione per le 5 YAML.
docs/handoff/phase-1-memory.md:63:3. **Opzionale: Fase 2.2** — embed `go:embed` + route `/swagger` nel binary.
docs/handoff/phase-3-handoff.md:17:- Phase 2 (Swagger UI) is complete: ~/astrate/docs/swagger-ui/index.html.
docs/handoff/phase-3-handoff.md:19:- To view Swagger UI locally, serve from ~/astrate/docs, then open http://localhost:9090/swagger-ui/.
docs/handoff/phase-3-handoff.md:22:- Phase 2.2 (go:embed + /swagger route) was not done.
docs/handoff/phase-4-memory.md:7:- Copied DESIGN.md, ROADMAP.md, api/, swagger-ui/ into `docs/site/` so MkDocs can serve them.
docs/handoff/phase-4-memory.md:9:- Verified clean build: `mkdocs build` succeeds, all 15 pages + DESIGN + ROADMAP + api + swagger-ui present in output.
docs/handoff/phase-4-memory.md:23:- `docs/site/swagger-ui/` — copied from `docs/swagger-ui/`
docs/handoff/phase-4-memory.md:28:- All nav links, DESIGN, ROADMAP, api/, swagger-ui/ links resolve in built output
docs/handoff/phase-4-memory.md:29:- `site-dist/` contains: index.html, 15 page dirs, DESIGN/, ROADMAP/, api/, swagger-ui/, assets/, search/
docs/handoff/phase-6-memory.md:6:- Verified full `docs/` directory structure: 15 narrative pages in `site/`, 5 OpenAPI YAML specs in `api/`, Swagger UI in `swagger-ui/`, Makefile, mkdocs.yml, requirements.txt.

$ rg -ln "astrate_native_api|yaml" internal/swagger docs/*.go 2>/dev/null; rg -rn "NativeEvent" . --glob '!docs/site/**' | head; rg -n "schemas/" docs/api/astarte_*.yaml | head -5
internal/swagger/swagger.go
internal/swagger/swagger_test.go
docs/embed.go
docs/api/astarte_housekeeping_api.yaml:40:                $ref: "#/components/schemas/DataEnvelopeRealmNames"
docs/api/astarte_housekeeping_api.yaml:71:                  $ref: "#/components/schemas/RealmCreate"
docs/api/astarte_housekeeping_api.yaml:86:                $ref: "#/components/schemas/DataEnvelopeRealm"
docs/api/astarte_housekeeping_api.yaml:124:                $ref: "#/components/schemas/DataEnvelopeRealm"
docs/api/astarte_housekeeping_api.yaml:169:                  $ref: "#/components/schemas/RealmPatch"

→ Read internal/swagger/swagger_test.go
→ Read internal/swagger/swagger.go

$ rg -ln "redocly|openapi-generator|swagger-cli|spectral" --glob '!docs/site/**' . | head; rg -n "incoming_data|device_connected" internal/appengine/stream/*_test.go | head -20; ls internal/appengine/stream/
./docs/handoff/phase-1-handoff.md
./docs/handoff/phase-2-handoff.md
./docs/handoff/phase-3-handoff.md
./docs/handoff/phase-2-memory.md
./docs/handoff/phase-1-memory.md
ws.go
ws_test.go

$ which mkdocs python3 gofmt go; go version; python3 -c "import yaml; print('pyyaml ok')"
/usr/bin/python3
/usr/bin/gofmt
/usr/bin/go
go version go1.26.7 linux/arm64
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; print('pyyaml ok')
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ cat docs/requirements.txt; python3 -m pip --version 2>&1 | head -2; rg -n "event|json" internal/appengine/stream/ws_test.go | head -30
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10
/usr/bin/python3: No module named pip
8:	"encoding/json"
23:// closedBus yields one pre-loaded event then closes — the handler writes the
24:// event and returns, making the SSE assertion deterministic without sleeps.
67:		req.Header.Set("Accept", "text/event-stream")
77:		req.Header.Set("Accept", "text/event-stream")
81:		if ct := rec.Header().Get("Content-Type"); ct != "text/event-stream" {
89:		if err := json.Unmarshal([]byte(strings.TrimSpace(strings.TrimPrefix(body, "data: "))), &ev); err != nil {
93:			t.Errorf("event = %+v", ev)

$ rg -n "yaml" go.mod; rg -rln "yaml.Unmarshal|yaml.v3|Decode.*yaml" --glob '*.go' . | head; ls docs/*.go; sed -n 1,30p docs/embed.go
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect
docs/embed.go
// Package docs embeds the Swagger UI static files and OpenAPI YAML specs so
// they are served directly from the astrate binary without external files.
package docs

import "embed"

// SwaggerUI holds the static Swagger UI files (index.html, CSS, JS refs).
//
//go:embed swagger-ui/*
var SwaggerUI embed.FS

// APIYAML holds the OpenAPI 3.0 YAML specifications for all five API surfaces.
//
//go:embed api/*.yaml
var APIYAML embed.FS

$ command -v mkdocs; ls ~/.local/bin 2>/dev/null; ls /usr/local/bin | head -20; timeout 8 curl -sSI https://pypi.org/simple/ | head -1
claude
heif-convert
heif-dec
heif-enc
heif-info
heif-thumbnailer
hermes
net-watchdog.sh
node
npm
npx
opencode-telegram
pihole
HTTP/2 200

$ rg -n "^func (operationBlock|componentBlock|operationDescription|containsLine|schemaPattern)" -A 30 internal/swagger/swagger_test.go | sed -n 1,120p
1336:func schemaPattern(t *testing.T, lines []string, indent string) string {
1337-	t.Helper()
1338-	for _, l := range lines {
1339-		if v, ok := strings.CutPrefix(l, indent+"pattern: "); ok {
1340-			return strings.Trim(v, `'"`)
1341-		}
1342-	}
1343-	return ""
1344-}
1345-
1346-// TestPairingUnregisterDeviceSemanticsDocumented guards that the pairing spec
1347-// describes `unregisterDevice` for what it does. "Removes a device from the
1348-// realm" reads as a device-and-data deletion; store.UnregisterDevice clears the
1349-// credential trail and flips status back to 'registered' and touches nothing
1350-// else (internal/store/devices.go), so the row, its interfaces, its datastream
1351-// data and its group memberships all survive, and a second DELETE still matches
1352-// the row — a 204 no-op, not the 404 a deletion would have produced.
1353-func TestPairingUnregisterDeviceSemanticsDocumented(t *testing.T) {
1354-	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
1355-	if err != nil {
1356-		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
1357-	}
1358-	lines := strings.Split(string(b), "\n")
1359-
1360-	desc := operationDescription(t, operationBlock(t, lines, "unregisterDevice"))
1361-	for _, want := range []string{
1362-		"not a deletion",
1363-		"registrable again",
1364-		"interfaces",
1365-		"datastream data",
1366-		"group memberships",
--
1383:func operationDescription(t *testing.T, block []string) string {
1384-	t.Helper()
1385-	var text []string
1386-	for i, l := range block {
1387-		if l == "      description: |" {
1388-			for _, b := range block[i+1:] {
1389-				if strings.TrimSpace(b) == "" {
1390-					continue
1391-				}
1392-				if !strings.HasPrefix(b, "        ") {
1393-					break
1394-				}
1395-				text = append(text, strings.TrimSpace(b))
1396-			}
1397-			return strings.Join(text, " ")
1398-		}
1399-		if v, ok := strings.CutPrefix(l, "      description: "); ok {
1400-			return strings.Trim(strings.TrimSpace(v), `"'`)
1401-		}
1402-	}
1403-	t.Fatal("operation declares no description")
1404-	return ""
1405-}
1406-
1407-// operationBlock returns the lines of the operation with the given operationId,
1408-// up to the next operation, path, or the components section.
1409:func operationBlock(t *testing.T, lines []string, operationID string) []string {
1410-	t.Helper()
1411-	marker := "      operationId: " + operationID
1412-	start := -1
1413-	for i, l := range lines {
1414-		if l == marker {
1415-			start = i
1416-			break
1417-		}
1418-	}
1419-	if start < 0 {
1420-		t.Fatalf("spec declares no operation %q", operationID)
1421-	}
1422-	endRe := regexp.MustCompile(`^      operationId: |^  \S|^components:`)
1423-	for i := start + 1; i < len(lines); i++ {
1424-		if endRe.MatchString(lines[i]) {
1425-			return lines[start+1 : i]
1426-		}
1427-	}
1428-	return lines[start+1:]
1429-}
1430-
1431-// componentBlock returns the lines of the components entry starting with the
1432-// given key line, up to the next entry at the same or shallower indentation.
1433:func componentBlock(t *testing.T, lines []string, key string) []string {
1434-	t.Helper()
1435-	start := -1
1436-	for i, l := range lines {
1437-		if l == key {
1438-			start = i
1439-			break
1440-		}
1441-	}
1442-	if start < 0 {
1443-		t.Fatalf("spec declares no components entry %q", key)
1444-	}
1445-	for i := start + 1; i < len(lines); i++ {
1446-		if strings.TrimSpace(lines[i]) == "" {
1447-			continue
1448-		}
1449-		if !strings.HasPrefix(lines[i], "      ") {
1450-			return lines[start+1 : i]
1451-		}
1452-	}
1453-	return lines[start+1:]
1454-}
1455-
1456:func containsLine(lines []string, want string) bool {
1457-	for _, l := range lines {
1458-		if l == want {
1459-			return true
1460-		}
1461-	}
1462-	return false
1463-}
1464-
1465-func embeddedYAMLFilenames() ([]string, error) {
1466-	var names []string
1467-	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
1468-		if err != nil {
1469-			return err
1470-		}

$ mkdir -p /tmp/opencode/wireprobe && cat > /tmp/opencode/wireprobe/main.go <<'EOF'
package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// verbatim copy of internal/appengine/stream/wireEvent + toWire inputs
type wireEvent struct {
	Event     string    `json:"event"`
	Realm     string    `json:"realm"`
	DeviceID  string    `json:"device_id"`
	Interface string    `json:"interface,omitempty"`
	Path      string    `json:"path,omitempty"`
	Value     any       `json:"value,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	ts, _ := time.Parse(time.RFC3339Nano, "2024-05-30T13:49:57.045Z")
	cases := []wireEvent{
		{Event: "device_connected", Realm: "test", DeviceID: "f0VMRgIBAQAAAAAAAAAAAA", Timestamp: ts},
		{Event: "incoming_data", Realm: "test", DeviceID: "f0VMRgIBAQAAAAAAAAAAAA", Interface: "org.example.Sensors", Path: "/sensors/temperature/value", Value: false, Timestamp: ts},
		{Event: "incoming_data", Realm: "test", DeviceID: "f0VMRgIBAQAAAAAAAAAAAA", Interface: "org.example.Sensors", Path: "/sensors/temperature/value", Value: 0, Timestamp: ts},
		{Event: "incoming_data", Realm: "test", DeviceID: "f0VMRgIBAQAAAAAAAAAAAA", Interface: "org.example.Sensors", Path: "/sensors/temperature/value", Value: "", Timestamp: ts},
		{Event: "incoming_data", Realm: "test", DeviceID: "f0VMRgIBAQAAAAAAAAAAAA", Interface: "org.example.Sensors", Path: "/sensors/temperature/value", Value: nil, Timestamp: ts},
		{Event: "device_error", Realm: "test", DeviceID: "f0VMRgIBAQAAAAAAAAAAAA", Timestamp: ts},
	}
	for _, c := range cases {
		b, err := json.Marshal(c)
		fmt.Println(string(b), err)
	}
}
EOF
cd /tmp/opencode/wireprobe && go mod init probe >/dev/null 2>&1; go run .
{"event":"device_connected","realm":"test","device_id":"f0VMRgIBAQAAAAAAAAAAAA","timestamp":"2024-05-30T13:49:57.045Z"} <nil>
{"event":"incoming_data","realm":"test","device_id":"f0VMRgIBAQAAAAAAAAAAAA","interface":"org.example.Sensors","path":"/sensors/temperature/value","value":false,"timestamp":"2024-05-30T13:49:57.045Z"} <nil>
{"event":"incoming_data","realm":"test","device_id":"f0VMRgIBAQAAAAAAAAAAAA","interface":"org.example.Sensors","path":"/sensors/temperature/value","value":0,"timestamp":"2024-05-30T13:49:57.045Z"} <nil>
{"event":"incoming_data","realm":"test","device_id":"f0VMRgIBAQAAAAAAAAAAAA","interface":"org.example.Sensors","path":"/sensors/temperature/value","value":"","timestamp":"2024-05-30T13:49:57.045Z"} <nil>
{"event":"incoming_data","realm":"test","device_id":"f0VMRgIBAQAAAAAAAAAAAA","interface":"org.example.Sensors","path":"/sensors/temperature/value","timestamp":"2024-05-30T13:49:57.045Z"} <nil>
{"event":"device_error","realm":"test","device_id":"f0VMRgIBAQAAAAAAAAAAAA","timestamp":"2024-05-30T13:49:57.045Z"} <nil>

Measured the marshal shape (probe in `/tmp/opencode/wireprobe`): lifecycle events emit only the four always-present keys; `value: false`/`0`/`""` are emitted, `nil` is dropped. Now the spec edits.
← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -300,9 +300,10 @@
   Optional `device_id` and `interface` query parameters narrow the
   event filter.
 
   **WebSocket protocol:** server-push only — client frames are drained
-  and ignored. Events are JSON objects:
+  and ignored. Every server text frame is one JSON event object shaped
+  as `NativeEvent`, referenced from the 101 response below:
 
   ```json
   {
     "event": "incoming_data",
@@ -314,10 +315,35 @@
     "timestamp": "2024-05-30T13:49:57.045Z"
   }
   ```
 
-  **SSE protocol:** each event is a `data:` frame followed by a blank line.
+  That is the data-event shape, not the general one: `interface`, `path`
+  and `value` are dropped when empty, so a lifecycle event is only its
+  four always-present keys — measured on the wire:
 
+  ```json
+  {
+    "event": "device_connected",
+    "realm": "test",
+    "device_id": "f0VMRgIBAQAAAAAAAAAAAA",
+    "timestamp": "2024-05-30T13:49:57.045Z"
+  }
+  ```
+
+  **SSE protocol:** each event is a `data:` frame followed by a blank
+  line, carrying the same `NativeEvent` object as its `data:` body.
+
+  This socket serialises through its own `wireEvent` struct
+  (internal/appengine/stream/ws.go) rather than the bus event, so it
+  drops `InterfaceMajor`, `IP`, `ErrorName` and `ErrorMetadata`, all of
+  which `stream.Event` carries (internal/engine/stream/bus.go): a
+  `device_error` reaches the client with no reason and no metadata, and
+  a `device_connected` reaches it with no `IP`, unlike the Phoenix twin
+  below, which marshals the whole struct into the frame. That is a code
+  gap and not a spec one — the four fields are deliberately absent from
+  `NativeEvent`, and putting them on this wire would be a code change,
+  not a documentation change.
+
   Authentication: `a_ch` realm JWT via Authorization header.
 tags: [WebSocket]
 security:
   - a_ch: []

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -375,9 +375,17 @@
       outside it — a non-`sse` value simply does not match and is not
       reported as an error.
 responses:
   "101":
-    description: WebSocket upgrade successful.
+    description: |
+      WebSocket upgrade successful. OpenAPI 3.0 has no message field for
+      the WebSocket protocol, so the shape of every server text frame
+      after the upgrade is recorded here: one `NativeEvent` object per
+      event, pushed as JSON.
+    content:
+      application/json:
+        schema:
+          $ref: "#/components/schemas/NativeEvent"
   "200":
     description: |
       SSE stream started, when the `Accept` header is exactly
       `text/event-stream` or `transport=sse`. The handler sets

$ rg -n '"101"|NativeEvent|WebSocket upgrade' docs/api/astrate_native_api.yaml
305:        as `NativeEvent`, referenced from the 101 response below:
333:        line, carrying the same `NativeEvent` object as its `data:` body.
343:        `NativeEvent`, and putting them on this wire would be a code change,
378:        "101":
380:            WebSocket upgrade successful. OpenAPI 3.0 has no message field for
382:            after the upgrade is recorded here: one `NativeEvent` object per
387:                $ref: "#/components/schemas/NativeEvent"
462:        "101":
463:          description: WebSocket upgrade successful.

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -554,8 +554,61 @@
         properties:
           detail:
             type: string
 
+  NativeEvent:
+    type: object
+    description: |
+      One event pushed by /astrate/v1/{realm}/socket on either transport: a
+      WebSocket server text frame, or the body of an SSE `data:` line.
+      Serialised by the wireEvent struct in
+      internal/appengine/stream/ws.go from the engine's bus event.
+    required: [event, realm, device_id, timestamp]
+    properties:
+      event:
+        type: string
+        description: |
+          Event kind: one of the `Kind*` constants of
+          internal/engine/stream (bus.go), which mirror the Astarte trigger
+          event names.
+        enum:
+          - incoming_data
+          - device_registered
+          - device_connected
+          - device_disconnected
+          - device_deletion_started
+          - device_deletion_finished
+          - device_error
+      realm:
+        type: string
+        description: The tenant the event belongs to.
+      device_id:
+        type: string
+        description: The encoded device ID the event is about.
+      interface:
+        type: string
+        description: |
+          Interface of a data event. Empty for lifecycle events, and an
+          empty value is omitted from the JSON entirely, so this key is not
+          guaranteed to be present.
+      path:
+        type: string
+        description: |
+          Path inside the interface, for a data event. Omitted when empty,
+          like `interface`.
+      value:
+        description: |
+          JSON rendering of a data event's value: any JSON value, which is
+          why no type is declared here. Omitted only when the value is nil
+          (property unset, lifecycle events) — `omitempty` on an `any` field
+          drops nil alone, so `false`, `0` and `""` are all emitted.
+      timestamp:
+        type: string
+        format: date-time
+        description: |
+          Event instant, RFC 3339 (the effective sample timestamp for data
+          events).
+
 responses:
   Unauthorized:
     description: Unauthorized — missing or invalid JWT/token.
     content:

$ wc -l internal/swagger/swagger_test.go; sed -n '1465,$p' internal/swagger/swagger_test.go
1499 internal/swagger/swagger_test.go
func embeddedYAMLFilenames() ([]string, error) {
	var names []string
	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".yaml") {
			names = append(names, strings.TrimPrefix(path, "api/"))
		}
		return nil
	})
	sort.Strings(names)
	return names, err
}

func get(t *testing.T, url string) string {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s: status = %d, want %d", url, resp.StatusCode, http.StatusOK)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading %s: %v", url, err)
	}
	return string(b)
}

$ sed -n '560,625p' docs/api/astrate_native_api.yaml | cat -A | sed -n '1,70p' | cut -c1-120
description: |$
        One event pushed by /astrate/v1/{realm}/socket on either transport: a$
        WebSocket server text frame, or the body of an SSE `data:` line.$
        Serialised by the wireEvent struct in$
        internal/appengine/stream/ws.go from the engine's bus event.$
      required: [event, realm, device_id, timestamp]$
      properties:$
        event:$
          type: string$
          description: |$
            Event kind: one of the `Kind*` constants of$
            internal/engine/stream (bus.go), which mirror the Astarte trigger$
            event names.$
          enum:$
            - incoming_data$
            - device_registered$
            - device_connected$
            - device_disconnected$
            - device_deletion_started$
            - device_deletion_finished$
            - device_error$
        realm:$
          type: string$
          description: The tenant the event belongs to.$
        device_id:$
          type: string$
          description: The encoded device ID the event is about.$
        interface:$
          type: string$
          description: |$
            Interface of a data event. Empty for lifecycle events, and an$
            empty value is omitted from the JSON entirely, so this key is not$
            guaranteed to be present.$
        path:$
          type: string$
          description: |$
            Path inside the interface, for a data event. Omitted when empty,$
            like `interface`.$
        value:$
          description: |$
            JSON rendering of a data event's value: any JSON value, which is$
            why no type is declared here. Omitted only when the value is nil$
            (property unset, lifecycle events) M-bM-^@M-^T `omitempty` on an `any` field$
            drops nil alone, so `false`, `0` and `""` are all emitted.$
        timestamp:$
          type: string$
          format: date-time$
          description: |$
            Event instant, RFC 3339 (the effective sample timestamp for data$
            events).$
$
  responses:$
    Unauthorized:$
      description: Unauthorized M-bM-^@M-^T missing or invalid JWT/token.$
      content:$
        application/json; charset=utf-8:$
          schema:$
            $ref: "#/components/schemas/ErrorDetail"$
          example:$
            errors:$
              detail: Unauthorized$
$
    Forbidden:$
      description: Forbidden M-bM-^@M-^T the JWT is valid but does not authorize this path.$
      content:$
        application/json; charset=utf-8:$

$ tail -c 400 internal/swagger/swagger_test.go
func get(t *testing.T, url string) string {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s: status = %d, want %d", url, resp.StatusCode, http.StatusOK)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading %s: %v", url, err)
	}
	return string(b)
}

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -1496,4 +1496,91 @@
 		t.Fatalf("reading %s: %v", url, err)
 	}
 	return string(b)
 }
+
+// TestNativeEventSchemaDocumented guards that the native spec documents the
+// event payload of /astrate/v1/{realm}/socket as a components.schemas entry
+// that the operation $refs, instead of the fenced example that used to be the
+// only description of it. The example is not the general shape: wireEvent
+// marks interface, path and value omitempty (internal/appengine/stream/ws.go),
+// so a lifecycle event marshals to its four always-present keys alone, while
+// omitempty on an `any` drops nil only — value false, 0 and "" are emitted.
+// So the schema requires exactly those four keys and describes the other three
+// as conditional.
+//
+// The enum is read out of the Kind* constants in
+// internal/engine/stream/bus.go rather than written into the test, so a kind
+// the bus can publish is a kind the spec must list, in both directions.
+//
+// InterfaceMajor, IP, ErrorName and ErrorMetadata are asserted absent from the
+// schema: stream.Event carries all four (internal/engine/stream/bus.go), the
+// socket's wireEvent does not, so a device_error reaches the client with no
+// reason. Documenting them here would promise a wire this socket does not
+// send, so the operation description is where that code gap is recorded — and
+// it has to keep saying it, or the next reader "fixes" the schema instead of
+// the code.
+func TestNativeEventSchemaDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astrate_native_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astrate_native_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	block := componentBlock(t, lines, "    NativeEvent:")
+	if !containsLine(block, "      required: [event, realm, device_id, timestamp]") {
+		t.Error("NativeEvent does not require exactly event, realm, device_id and timestamp")
+	}
+	for _, prop := range []string{"event", "realm", "device_id", "interface", "path", "value", "timestamp"} {
+		propertyBlock(t, block, prop) // fails the test when the property is absent
+	}
+	for _, name := range []string{"interface_major", "ip", "error_name", "error_metadata"} {
+		if containsLine(block, "        "+name+":") {
+			t.Errorf("NativeEvent documents %q, which this socket's wireEvent does not send", name)
+		}
+	}
+
+	src, err := os.ReadFile("../../internal/engine/stream/bus.go")
+	if err != nil {
+		t.Fatalf("reading internal/engine/stream/bus.go: %v", err)
+	}
+	kinds := regexp.MustCompile(`(?m)^\s*Kind\w+\s*=\s*"([^"]+)"`).FindAllStringSubmatch(string(src), -1)
+	if len(kinds) == 0 {
+		t.Fatal("bus.go declares no Kind* constants; this test's premise no longer holds")
+	}
+	want := make(map[string]bool, len(kinds))
+	for _, k := range kinds {
+		want[k[1]] = true
+	}
+	got := map[string]bool{}
+	for _, l := range propertyBlock(t, block, "event") {
+		if v, ok := strings.CutPrefix(l, "            - "); ok {
+			got[v] = true
+		}
+	}
+	for k := range want {
+		if !got[k] {
+			t.Errorf("NativeEvent.event enum is missing %q, which bus.go publishes", k)
+		}
+	}
+	for k := range got {
+		if !want[k] {
+			t.Errorf("NativeEvent.event enum lists %q, which no Kind* constant in bus.go declares", k)
+		}
+	}
+
+	sock := operationBlock(t, lines, "nativeWebSocket")
+	if !containsLine(sock, `                $ref: "#/components/schemas/NativeEvent"`) {
+		t.Error("nativeWebSocket $refs no NativeEvent schema; a generated client learns the payload from prose alone")
+	}
+	desc := operationDescription(t, sock)
+	for _, want := range []string{
+		"NativeEvent",
+		"InterfaceMajor", "`IP`", "ErrorName", "ErrorMetadata",
+		"stream.Event", "Phoenix twin", "code gap",
+		"lifecycle event",
+	} {
+		if !strings.Contains(desc, want) {
+			t.Errorf("nativeWebSocket description does not say %q", want)
+		}
+	}
+}

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -6,8 +6,9 @@
 "io"
 "io/fs"
 "net/http"
 "net/http/httptest"
+"os"
 "reflect"
 "regexp"
 "slices"
 "sort"

$ gofmt -l . ; go vet ./internal/swagger/ && go test ./internal/swagger/ 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/swagger	0.088s

$ cd /root/astrate-mule && cp docs/api/astrate_native_api.yaml /tmp/opencode/native.orig.yaml && sha256sum docs/api/astrate_native_api.yaml && awk '
  /^    NativeEvent:$/ {skip=1}
  skip && /^  responses:$/ {skip=0}
  skip==0 {print}
' /tmp/opencode/native.orig.yaml | grep -v '"#/components/schemas/NativeEvent"' > /tmp/opencode/native.stripped.yaml && cp /tmp/opencode/native.stripped.yaml docs/api/astrate_native_api.yaml && go test ./internal/swagger/ -run TestNativeEventSchemaDocumented 2>&1 | tail -15; cp /tmp/opencode/native.orig.yaml docs/api/astrate_native_api.yaml; sha256sum docs/api/astrate_native_api.yaml
3c1db6a070a01bc249c7d673f11d4b655cbd07ad49b5cfc5faca0c6c5e9da1e9  docs/api/astrate_native_api.yaml
--- FAIL: TestNativeEventSchemaDocumented (0.00s)
    swagger_test.go:1530: spec declares no components entry "    NativeEvent:"
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.014s
FAIL
3c1db6a070a01bc249c7d673f11d4b655cbd07ad49b5cfc5faca0c6c5e9da1e9  docs/api/astrate_native_api.yaml

$ cd /root/astrate-mule && go vet ./... 2>&1 | tail -5; go test -race ./... 2>&1 | tail -40
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/housekeeping	0.030s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/httpx	0.058s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/observability	0.057s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/pairing	0.033s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/pairing/ca	0.025s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/realm	0.030s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/store	0.030s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.039s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/testutil	0.039s
?   	github.com/astrate-platform/astrate/migrations	[no test files]
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/astarteapi	0.039s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/deviceid	0.041s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/interfaceschema	0.029s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/payload	0.029s
FAIL

$ cd /root/astrate-mule && go test ./... 2>&1 | grep -v "^ok\|no test files" | head -30; echo "EXIT=$?"
EXIT=0

$ cd /root/astrate-mule && go test ./... > /tmp/opencode/gotest.log 2>&1; echo "go test exit=$?"; grep -c "^ok" /tmp/opencode/gotest.log; grep -v "^ok\|no test files" /tmp/opencode/gotest.log | head
go test exit=0
31

$ cd /root/astrate-mule && cat .gitignore 2>/dev/null; grep -n "site_dir\|docs_dir" docs/mkdocs.yml; python3 -c "import venv, ensurepip; print('venv ok')" 2>&1 | tail -1
# Build artifacts
dist/
/astrate
*.exe
*.test
*.out

# Coverage
coverage.*
*.cover

# bbolt session stores (broker persistence, §3.1)
*.db

# Local env / editor
.env
.envrc
.idea/
.vscode/
*.swp
.DS_Store

# Local AI documentation workflow. The handoff files under docs/handoff/ ARE tracked:
# they are the project's phase history, and the agents that need it most are the cheap ones
# that cannot re-derive it.
/docs/AI-DOCUMENTATION-WORKFLOW.md

# MkDocs build output
/site-dist/

# Docs site copies (refreshed by docs/make sync)
/docs/site/DESIGN.md
/docs/site/ROADMAP.md
/docs/site/api/
/docs/site/swagger-ui/

# Local-only compose overrides
docker-compose.override.yml
__pycache__/

# bench harness artifacts (bench/README.md). Note bench/results/ is deliberately NOT here:
# a benchmark number without its host and its history is worthless, so results are evidence
# and evidence is versioned.
/bench/.astarte/
/bench/keys/
/bench/astrate/conf/
/bench/*.json
/bench/bench

# The mule host's tick transcript: local to whichever machine ran it, grows without
# bound, and got committed by accident on 2026-09-04.
.mule/cron.log
5:docs_dir: site
6:site_dir: ../site-dist
venv ok

$ python3 -m venv /tmp/opencode/mkdocs-venv && /tmp/opencode/mkdocs-venv/bin/pip install -q -r /root/astrate-mule/docs/requirements.txt 2>&1 | tail -3; /tmp/opencode/mkdocs-venv/bin/mkdocs --version
mkdocs, version 1.6.1 from /tmp/opencode/mkdocs-venv/lib/python3.13/site-packages/mkdocs (Python 3.13)

$ cd /root/astrate-mule && PATH=/tmp/opencode/mkdocs-venv/bin:$PATH make -C docs build 2>&1 | tail -20
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
INFO    -  Documentation built in 3.97 seconds
make: Leaving directory '/root/astrate-mule/docs'

$ mkdir -p /tmp/opencode/specload && cd /tmp/opencode/specload && go mod init specload >/dev/null 2>&1; GOFLAGS=-mod=mod GOPROXY=off go get gopkg.in/yaml.v3@v3.0.1 2>&1 | tail -3; cat > main.go <<'EOF'
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	files, _ := filepath.Glob("/root/astrate-mule/docs/api/*.yaml")
	if len(files) == 0 {
		fmt.Println("no specs found")
		os.Exit(1)
	}
	bad := false
	for _, f := range files {
		raw, err := os.ReadFile(f)
		if err != nil {
			fmt.Println("read:", err)
			bad = true
			continue
		}
		var doc map[string]any
		if err := yaml.Unmarshal(raw, &doc); err != nil {
			fmt.Printf("FAIL %s: %v\n", filepath.Base(f), err)
			bad = true
			continue
		}
		var refs []string
		var walk func(any)
		walk = func(v any) {
			switch t := v.(type) {
			case map[string]any:
				for k, val := range t {
					if k == "$ref" {
						if s, ok := val.(string); ok {
							refs = append(refs, s)
						}
					}
					walk(val)
				}
			case []any:
				for _, val := range t {
					walk(val)
				}
			}
		}
		walk(doc)
		missing := 0
		for _, r := range refs {
			if !strings.HasPrefix(r, "#/") {
				continue
			}
			cur := doc
			for _, part := range strings.Split(strings.TrimPrefix(r, "#/"), "/") {
				next, ok := cur[part].(map[string]any)
				if !ok {
					fmt.Printf("FAIL %s: unresolvable %s\n", filepath.Base(f), r)
					missing++
					bad = true
					break
				}
				cur = next
			}
		}
		fmt.Printf("ok   %-32s refs=%d unresolved=%d\n", filepath.Base(f), len(refs), missing)
	}
	if bad {
		os.Exit(1)
	}
}
EOF
GOPROXY=off go run .
go: added gopkg.in/yaml.v3 v3.0.1
ok   astarte_appengine_api.yaml       refs=350 unresolved=0
ok   astarte_housekeeping_api.yaml    refs=46 unresolved=0
ok   astarte_pairing_api.yaml         refs=59 unresolved=0
ok   astarte_realm_management_api.yaml refs=154 unresolved=0
ok   astrate_native_api.yaml          refs=23 unresolved=0

$ cd /root/astrate-mule && git status --porcelain && git diff --stat
M docs/api/astrate_native_api.yaml
 M internal/swagger/swagger_test.go
 docs/api/astrate_native_api.yaml | 93 ++++++++++++++++++++++++++++++++++++++--
 internal/swagger/swagger_test.go | 88 +++++++++++++++++++++++++++++++++++++
 2 files changed, 178 insertions(+), 3 deletions(-)

→ Read docs/api/astrate_native_api.yaml [offset=284, limit=115]

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -1539,8 +1539,15 @@
 		t.Errorf("NativeEvent documents %q, which this socket's wireEvent does not send", name)
 	}
 }
 
+value := strings.Join(propertyBlock(t, block, "value"), "\n")
+for _, want := range []string{"nil", "`false`", "`0`"} {
+	if !strings.Contains(value, want) {
+		t.Errorf("NativeEvent.value description does not say %q is emitted rather than dropped", want)
+	}
+}
+
 src, err := os.ReadFile("../../internal/engine/stream/bus.go")
 if err != nil {
 	t.Fatalf("reading internal/engine/stream/bus.go: %v", err)
 }

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -1542,9 +1542,9 @@
 
 value := strings.Join(propertyBlock(t, block, "value"), "\n")
 for _, want := range []string{"nil", "`false`", "`0`"} {
 	if !strings.Contains(value, want) {
-		t.Errorf("NativeEvent.value description does not say %q is emitted rather than dropped", want)
+		t.Errorf("NativeEvent.value description does not record %q", want)
 	}
 }
 
 src, err := os.ReadFile("../../internal/engine/stream/bus.go")

Error: Error from provider (Console): Rate limit exceeded. Please try again later.
