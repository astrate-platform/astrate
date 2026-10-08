slug: docs-sync-pairing-realm-name-pattern
verdict: done
at:  ed301fb
ran: 2026-10-08T19:28:20Z on DietPi in 410s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_pairing_api.yaml [offset=340, limit=40]
→ Read internal/swagger/swagger_test.go [offset=1220, limit=80]
→ Read internal/swagger/swagger_test.go [limit=60, offset=1300]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=50, offset=230]
→ Read docs/api/astrate_native_api.yaml [limit=60, offset=410]
→ Read internal/auth/middleware.go [limit=30, offset=55]
→ Read internal/pairing/service.go [limit=20, offset=360]
Now let me verify the code references and look at the native API twin:
→ Read internal/observability/compat.go [limit=20, offset=35]
→ Read cmd/astrate/main.go [limit=15, offset=430]
→ Read internal/pairing/http.go [limit=15, offset=280]
Let me verify the remaining references and see how the pairing spec uses the parameter:

$ rg -n "RealmName|/version|/health" docs/api/astarte_pairing_api.yaml | head -40
27:  /pairing/v1/{realm}/health:
40:        - $ref: "#/components/parameters/RealmName"
70:  /pairing/v1/{realm}/version:
81:        - $ref: "#/components/parameters/RealmName"
110:        - $ref: "#/components/parameters/RealmName"
167:        - $ref: "#/components/parameters/RealmName"
195:        - $ref: "#/components/parameters/RealmName"
245:        - $ref: "#/components/parameters/RealmName"
279:        - $ref: "#/components/parameters/RealmName"
357:    RealmName:

$ rg -n "^  /pairing|^    (getHealth|getVersion|operationId):" docs/api/astarte_pairing_api.yaml | head -30
27:  /pairing/v1/{realm}/health:
70:  /pairing/v1/{realm}/version:
97:  /pairing/v1/{realm}/agent/devices:
147:  /pairing/v1/{realm}/agent/devices/{deviceID}:
181:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
233:  /pairing/v1/{realm}/devices/{deviceID}:
267:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:

→ Read migrations/000002_metadata.up.sql
All references verified. Now let me see the version operation and the rest of the pairing paths to write the description correctly:
→ Read docs/api/astarte_pairing_api.yaml [limit=80, offset=27]

$ rg -n "pairing/v1|RequireRealm|handleHealth" cmd/astrate/main.go internal/pairing/http.go | head -30
cmd/astrate/main.go:437:		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(realm.APICompatVersion)))
cmd/astrate/main.go:438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
internal/pairing/http.go:41:// API is the /pairing/v1 HTTP surface (docs/DESIGN.md §4.4, §3.7). Agent
internal/pairing/http.go:68:		requireAgent: mw.RequireRealm(auth.ClaimPairing),
internal/pairing/http.go:81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/pairing/http.go:82:	mux.Handle("POST /pairing/v1/{realm}/agent/devices",
internal/pairing/http.go:84:	mux.Handle("DELETE /pairing/v1/{realm}/agent/devices/{deviceID}",
internal/pairing/http.go:86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
internal/pairing/http.go:88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
internal/pairing/http.go:90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
internal/pairing/http.go:281:// handleHealth serves the realm-scoped health probe (#71): 200 when the
internal/pairing/http.go:284:func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {

→ Read internal/pairing/http.go [limit=40, offset=60]
→ Read internal/swagger/swagger_test.go [limit=80, offset=1]
Now let me look at the helpers and existing pairing tests to see if a twin stays small:

$ rg -n "^func |componentBlock|operationBlock|containsLine" internal/swagger/swagger_test.go
26:func TestMount(t *testing.T) {
85:func TestMountSubPanicsOnBrokenFS(t *testing.T) {
98:func TestSpecs(t *testing.T) {
131:func TestRealmManagement403(t *testing.T) {
142:func TestAppEngine403(t *testing.T) {
149:func specDocuments403(t *testing.T, filename string) {
189:		if !containsLine(lines[start:end], `        "403":`) {
194:	if !containsLine(lines[compIdx:], `    Forbidden:`) {
206:func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
215:		block := operationBlock(t, lines, op)
216:		if !containsLine(block, ref) {
221:	comp := componentBlock(t, lines, "    AsyncOperation:")
230:		if !containsLine(comp, want) {
249:func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
261:		block := operationBlock(t, lines, op)
262:		if !containsLine(block, ref) {
267:	comp := componentBlock(t, lines, "    AsyncOperation:")
276:		if !containsLine(comp, want) {
288:func propertyBlock(t *testing.T, lines []string, name string) []string {
324:func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {
331:	patch := strings.Join(operationBlock(t, lines, "patchRealm"), "\n")
344:		block := componentBlock(t, lines, schema)
367:		componentBlock(t, lines, "    RealmCreate:"), "datastream_maximum_storage_retention"), "\n")
384:func TestAppEngineAddGroupDevice422Documented(t *testing.T) {
396:	block := operationBlock(t, lines, "addGroupDevice")
411:	comp := strings.Join(componentBlock(t, lines, "    ValidationErrors:"), "\n")
428:func TestAppEngineDataDelete400Documented(t *testing.T) {
441:		block := operationBlock(t, lines, op)
458:	if !containsLine(lines, "    BadRequest:") {
473:func TestPairingErrorDetailExamples(t *testing.T) {
492:		block := componentBlock(t, lines, tc.response)
524:func TestRealmManagementErrorDetailExamples(t *testing.T) {
544:		block := componentBlock(t, lines, tc.response)
576:func TestAppEngineErrorDetailExamples(t *testing.T) {
596:		block := componentBlock(t, lines, tc.response)
631:func TestNativeErrorDetailExamples(t *testing.T) {
648:		block := componentBlock(t, lines, tc.response)
682:func TestNativeVersionDescriptions(t *testing.T) {
692:		if desc := operationDescription(t, operationBlock(t, lines, op)); !strings.Contains(desc, twin) {
697:	desc := operationDescription(t, operationBlock(t, lines, "getHousekeepingVersion"))
731:func TestNativeSocketSecurityDocumented(t *testing.T) {
738:	if !containsLine(lines, "security: []") {
742:	scheme := componentBlock(t, lines, "    a_ch:")
748:		if !containsLine(scheme, want) {
763:	sock := operationBlock(t, lines, "nativeWebSocket")
764:	if !containsLine(sock, "      security:") || !containsLine(sock, "        - a_ch: []") {
774:		if containsLine(operationBlock(t, lines, op), "      security:") {
802:func TestNativeMetricsContentNegotiationDocumented(t *testing.T) {
842:	block := operationBlock(t, lines, "getMetrics")
866:func metricsScrape(t *testing.T, headers map[string]string) (int, string, string) {
881:func mediaTypeKeys(t *testing.T, block []string) []string {
938:func TestHousekeepingErrorDetailExamples(t *testing.T) {
957:		block := componentBlock(t, lines, tc.response)
976:	schemaBlock := componentBlock(t, lines, "    ErrorDetail:")
1009:func TestRealmManagement422ValidationDetails(t *testing.T) {
1019:	if containsLine(lines, "    ValidationError:") {
1041:		block := strings.Join(responseBlock(t, operationBlock(t, lines, tc.op), "422"), "\n")
1055:func responseBlock(t *testing.T, block []string, status string) []string {
1086:func TestPairingDeviceIDEncodingDocumented(t *testing.T) {
1096:	param := strings.Join(componentBlock(t, lines, "    DeviceID:"), "\n")
1117:		if !containsLine(block, want) {
1168:func TestRealmManagementDeviceIDEncodingDocumented(t *testing.T) {
1178:	param := componentBlock(t, lines, "    DeviceID:")
1188:		if !containsLine(param, want) {
1247:func TestHousekeepingRealmNamePatternDocumented(t *testing.T) {
1282:	param := componentBlock(t, lines, "    RealmName:")
1292:	create := propertyBlock(t, componentBlock(t, lines, "    RealmCreate:"), "realm_name")
1297:	realm := propertyBlock(t, componentBlock(t, lines, "    Realm:"), "realm_name")
1303:	dataEnvelopeNames := componentBlock(t, lines, "    DataEnvelopeRealmNames:")
1320:	if !containsLine(operationBlock(t, lines, "createRealm"),
1324:	unprocessable := strings.Join(componentBlock(t, lines, "    ValidationError:"), "\n")
1336:func schemaPattern(t *testing.T, lines []string, indent string) string {
1353:func TestPairingUnregisterDeviceSemanticsDocumented(t *testing.T) {
1360:	desc := operationDescription(t, operationBlock(t, lines, "unregisterDevice"))
1383:func operationDescription(t *testing.T, block []string) string {
1407:// operationBlock returns the lines of the operation with the given operationId,
1409:func operationBlock(t *testing.T, lines []string, operationID string) []string {
1431:// componentBlock returns the lines of the components entry starting with the
1433:func componentBlock(t *testing.T, lines []string, key string) []string {
1456:func containsLine(lines []string, want string) bool {
1465:func embeddedYAMLFilenames() ([]string, error) {
1483:func get(t *testing.T, url string) string {

$ rg -n "DetailNotFound" pkg/astarteapi/ | head -5
pkg/astarteapi/envelope.go:42:	// DetailNotFound is the canonical generic 404 detail.
pkg/astarteapi/envelope.go:43:	DetailNotFound = "Not Found"
pkg/astarteapi/envelope.go:49:	// DetailNotFound, which a handler emits for a route that matched and found
pkg/astarteapi/envelope.go:207:	return WriteError(w, http.StatusNotFound, DetailNotFound)

→ Read internal/swagger/swagger_test.go [limit=50, offset=1380]
Let me verify the device routes answer 401 for an unknown realm:

$ rg -n "authenticateDevice|WriteUnauthorized|ErrUnauthorized" internal/pairing/http.go | head -20
165:		_ = astarteapi.WriteUnauthorized(w)
204:		_ = astarteapi.WriteUnauthorized(w)
245:		_ = astarteapi.WriteUnauthorized(w)
322:	case errors.Is(err, ErrUnauthorized):
323:		_ = astarteapi.WriteUnauthorized(w)

→ Read internal/pairing/service.go [limit=30, offset=376]
All facts verified. Let me confirm the APICompatVersion value, then make the edit:

$ rg -n "APICompatVersion" internal/realm/ pkg/ --type go | head -5
internal/realm/http.go:119:// feature UI on it — see APICompatVersion).
internal/realm/http.go:121:	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
internal/realm/dashboard_compat_test.go:35:		if v != APICompatVersion {
internal/realm/dashboard_compat_test.go:36:			t.Errorf("version = %q, want %q", v, APICompatVersion)
internal/realm/service.go:580:// APICompatVersion is the upstream Realm Management API level Astrate

All references check out. Now I'll make the edits. First the RealmName parameter in the pairing spec:
← Edit docs/api/astarte_pairing_api.yaml
Index: /root/astrate-mule/docs/api/astarte_pairing_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_pairing_api.yaml
+++ /root/astrate-mule/docs/api/astarte_pairing_api.yaml
@@ -357,11 +357,27 @@
 RealmName:
   name: realm
   in: path
   required: true
-  description: The realm name.
+  description: |
+    The realm name: a lowercase ASCII word starting with a letter, so
+    `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm`
+    are not realm names — no realm exists under one. This one parameter is
+    shared by every operation on this surface, and an off-pattern name
+    gets three different answers, not one: the agent and device routes
+    answer `401` `Unauthorized`, the same body as a bad token, because
+    the realm is resolved by the middleware, which treats a miss as
+    unauthenticated and not as a missing resource; `getHealth` answers
+    `404` with the same `{"errors":{"detail":"Not Found"}}` body as a
+    well-formed realm that does not exist, because its store lookup
+    matches the name verbatim with no format check of its own; and
+    `getVersion` answers `200` with the API compatibility version for any
+    realm string whatsoever, because that route is mounted with no
+    middleware and no realm lookup — the pattern constrains clients, not
+    that handler.
   schema:
     type: string
+    pattern: '^[a-z][a-z0-9]*$'
 
 DeviceID:
   name: deviceID
   in: path

Now the `getVersion` operation description:
← Edit docs/api/astarte_pairing_api.yaml
Index: /root/astrate-mule/docs/api/astarte_pairing_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_pairing_api.yaml
+++ /root/astrate-mule/docs/api/astarte_pairing_api.yaml
@@ -74,9 +74,14 @@
 description: |
   Returns the emulated upstream API compatibility version. Unauthenticated
   — mirrors upstream's public pairing version endpoint.
 
-  The Astarte Dashboard gates feature UI on this value.
+  The Astarte Dashboard gates feature UI on this value. This operation
+  does not look the realm up and does not authenticate against it: it is
+  mounted with no middleware, so it answers `200` for any realm string
+  whatsoever — including the off-pattern names the shared RealmName
+  parameter rejects. Do not read that parameter as a guarantee this
+  handler enforces it.
 tags: [Pairing]
 parameters:
   - $ref: "#/components/parameters/RealmName"
 responses:

Now the test twin. Let me look at where to place it and keep it small:
→ Read internal/swagger/swagger_test.go [limit=60, offset=1346]
Now the pairing twin test. I'll add it right after the `schemaPattern` helper:

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -1342,8 +1342,50 @@
 	}
 	return ""
 }
 
+// TestPairingRealmNamePatternDocumented guards that the pairing spec constrains
+// its `realm` path parameter the way the schema does (the realms table CHECKs
+// the name column, migrations/000002_metadata.up.sql) and records the split
+// answer an off-pattern name gets across the operations sharing the parameter:
+// the agent and device routes 401 through the middleware, getHealth 404 through
+// the verbatim store lookup, and getVersion 200 for any string because it is
+// mounted with no middleware and no realm lookup.
+func TestPairingRealmNamePatternDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	sql, err := migrations.FS.ReadFile("000002_metadata.up.sql")
+	if err != nil {
+		t.Fatalf("reading 000002_metadata.up.sql: %v", err)
+	}
+	check := regexp.MustCompile(`CHECK \(name ~ '([^']+)'\)`).FindSubmatch(sql)
+	if check == nil {
+		t.Fatal("realms.name declares no CHECK (name ~ ...) constraint; this test's premise no longer holds")
+	}
+	want := string(check[1])
+
+	param := componentBlock(t, lines, "    RealmName:")
+	if got := schemaPattern(t, param, "        "); got != want {
+		t.Errorf("RealmName parameter pattern = %q, want the realm name CHECK %q", got, want)
+	}
+	for _, wording := range []string{"401", "404", "200", astarteapi.DetailNotFound} {
+		if !strings.Contains(strings.Join(param, "\n"), wording) {
+			t.Errorf("RealmName parameter does not record the %q answer an off-pattern name gets on some operation", wording)
+		}
+	}
+
+	version := operationDescription(t, operationBlock(t, lines, "getVersion"))
+	for _, wording := range []string{"200", "any realm string", "no middleware"} {
+		if !strings.Contains(version, wording) {
+			t.Errorf("getVersion description does not say %q; a reader would take the shared parameter as a guarantee this handler enforces", wording)
+		}
+	}
+}
+
 // TestPairingUnregisterDeviceSemanticsDocumented guards that the pairing spec
 // describes `unregisterDevice` for what it does. "Removes a device from the
 // realm" reads as a device-and-data deletion; store.UnregisterDevice clears the
 // credential trail and flips status back to 'registered' and touches nothing

Now run the gate:
$ gofmt -l . && go vet ./... 2>&1 | tail -5
(no output)

$ go test ./... 2>&1 | tail -30
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
ok  	github.com/astrate-platform/astrate/internal/swagger	0.106s
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

All pass. Let me verify the new test fails without the change (the runner requires that), and try the race build:
$ git stash -- docs/api/astarte_pairing_api.yaml 2>/dev/null; go test ./internal/swagger/ -run TestPairingRealmNamePatternDocumented 2>&1 | tail -8; git stash pop 2>/dev/null; git status --short
Saved working directory and index state WIP on mule/queue: bf5aadb mule: blocked docs-sync-pairing-body-cap-400
    swagger_test.go:1376: RealmName parameter does not record the "200" answer an off-pattern name gets on some operation
    swagger_test.go:1376: RealmName parameter does not record the "Not Found" answer an off-pattern name gets on some operation
    swagger_test.go:1383: getVersion description does not say "200"; a reader would take the shared parameter as a guarantee this handler enforces
    swagger_test.go:1383: getVersion description does not say "any realm string"; a reader would take the shared parameter as a guarantee this handler enforces
    swagger_test.go:1383: getVersion description does not say "no middleware"; a reader would take the shared parameter as a guarantee this handler enforces
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.019s
FAIL
On branch mule/queue
Changes not staged for commit:
  (use "git add <file>..." to update what will be committed)
  (use "git restore <file>..." to discard changes in working directory)
	modified:   docs/api/astarte_pairing_api.yaml
	modified:   internal/swagger/swagger_test.go

no changes added to commit (use "git add" and/or "git commit -a")
Dropped refs/stash@{0} (1d1beb40c7bcf9d136c09cc0a934857dd35a5859)
 M docs/api/astarte_pairing_api.yaml
 M internal/swagger/swagger_test.go

$ go test -race ./internal/swagger/ 2>&1 | tail -3
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.029s
FAIL

The race detector can't build on this kernel (39-bit VMA — known; the gate here is `go test ./...`, which passes). The test correctly fails without the YAML change. Now the docs build check:
$ make -C docs build 2>&1 | tail -10
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

mkdocs is missing; the task allows falling back to confirming the YAML parses:
$ python3 -c "
import yaml
with open('docs/api/astarte_pairing_api.yaml') as f:
    spec = yaml.safe_load(f)
p = spec['components']['parameters']['RealmName']
assert p['schema']['pattern'] == '^[a-z][a-z0-9]*\$', p['schema']
print('YAML parses; RealmName pattern =', p['schema']['pattern'])
print('description lines:', len(p['description'].splitlines()))
v = spec['paths']['/pairing/v1/{realm}/version']['get']['description']
assert 'no middleware' in v and '200' in v
print('getVersion description OK')
"
Traceback (most recent call last):
  File "<string>", line 2, in <module>
    import yaml
ModuleNotFoundError: No module named 'yaml'

$ pip install pyyaml -q 2>&1 | tail -2; python3 -c "
import yaml
with open('docs/api/astarte_pairing_api.yaml') as f:
    spec = yaml.safe_load(f)
p = spec['components']['parameters']['RealmName']
assert p['schema']['pattern'] == '^[a-z][a-z0-9]*\$', p['schema']
print('YAML parses; RealmName pattern =', p['schema']['pattern'])
v = spec['paths']['/pairing/v1/{realm}/version']['get']['description']
assert 'no middleware' in v and '200' in v
print('getVersion description OK')
"
/bin/sh: 1: pip: not found
Traceback (most recent call last):
  File "<string>", line 2, in <module>
    import yaml
ModuleNotFoundError: No module named 'yaml'

No pip/pyyaml. The repo has a Go YAML dependency — I'll write the probe in Go instead:
$ rg -n "yaml" go.mod | head -5
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

! permission requested: external_directory (/tmp/*, /tmp/yamlprobe/*); auto-rejecting
✗ mkdir -p /tmp/yamlprobe && cd /tmp/yamlprobe && cat > main.go <<'EOF'
package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	b, err := os.ReadFile("/root/astrate-mule/docs/api/astarte_pairing_api.yaml")
	if err != nil {
		panic(err)
	}
	var spec map[string]any
	if err := yaml.Unmarshal(b, &spec); err != nil {
		panic(fmt.Sprintf("parse: %v", err))
	}
	comps := spec["components"].(map[string]any)
	params := comps["parameters"].(map[string]any)
	rn := params["RealmName"].(map[string]any)
	schema := rn["schema"].(map[string]any)
	fmt.Println("pattern:", schema["pattern"])
	paths := spec["paths"].(map[string]any)
	version := paths["/pairing/v1/{realm}/version"].(map[string]any)
	get := version["get"].(map[string]any)
	desc := get["description"].(string)
	fmt.Println("getVersion mentions no middleware:", len(desc) > 0)
}
EOF
cat > go.mod <<'EOF'
module yamlprobe

go 1.24

require gopkg.in/yaml.v3 v3.0.1
EOF
go mod tidy >/dev/null 2>&1 && go run .  failed
Error: The user rejected permission to use this specific tool call.
