slug: docs-sync-rm-delete-device-async-operation-param
verdict: done
at:  f965f39
ran: 2026-09-25T20:19:10Z on DietPi in 903s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md

# Todos
[•] Read task-named files: yaml deleteDevice block, internal/realm/http.go, COMPATIBILITY.md deviation 17, async_operation_test.go
[ ] Add async_operation param to deleteDevice in docs/api/astarte_realm_management_api.yaml
[ ] Update deviation 17 in docs/COMPATIBILITY.md to include the device leg
[ ] Add device-delete case to TestRealmManagementAsyncOperationParam
[ ] Run make -C docs build and confirm YAML loads
[ ] Run gate: go test ./... (or -race) and gofmt -l .

→ Read internal/realm/async_operation_test.go

$ cd /root/astrate-mule && git status --short && echo "---BRANCH---" && git rev-parse --abbrev-ref HEAD && echo "---GREP async_operation yaml---" && rg -n "async_operation" docs/api/astarte_realm_management_api.yaml
---BRANCH---
mule/queue
---GREP async_operation yaml---
903:      name: async_operation
910:        `async_operation=false`; Astrate always works synchronously and answers

→ Read docs/api/astarte_realm_management_api.yaml [offset=860, limit=80]
→ Read docs/api/astarte_realm_management_api.yaml [offset=640, limit=120]

$ cd /root/astrate-mule && rg -n "AsyncOperation" docs/api/astarte_realm_management_api.yaml && echo "---CONTEXT---" && rg -n -B4 -A2 '\$ref: "#/components/parameters/AsyncOperation"' docs/api/astarte_realm_management_api.yaml
106:        - $ref: "#/components/parameters/AsyncOperation"
249:        - $ref: "#/components/parameters/AsyncOperation"
322:        - $ref: "#/components/parameters/AsyncOperation"
827:        - $ref: "#/components/parameters/AsyncOperation"
902:    AsyncOperation:
---CONTEXT---
102-      security:
103-        - a_rma: []
104-      parameters:
105-        - $ref: "#/components/parameters/RealmName"
106:        - $ref: "#/components/parameters/AsyncOperation"
107-      requestBody:
108-        required: true
--
245-      parameters:
246-        - $ref: "#/components/parameters/RealmName"
247-        - $ref: "#/components/parameters/InterfaceName"
248-        - $ref: "#/components/parameters/MajorVersion"
249:        - $ref: "#/components/parameters/AsyncOperation"
250-      requestBody:
251-        required: true
--
318-      parameters:
319-        - $ref: "#/components/parameters/RealmName"
320-        - $ref: "#/components/parameters/InterfaceName"
321-        - $ref: "#/components/parameters/MajorVersion"
322:        - $ref: "#/components/parameters/AsyncOperation"
323-      responses:
324-        "204":
--
823-        - a_rma: []
824-      parameters:
825-        - $ref: "#/components/parameters/RealmName"
826-        - $ref: "#/components/parameters/PolicyName"
827:        - $ref: "#/components/parameters/AsyncOperation"
828-      responses:
829-        "204":

→ Read docs/COMPATIBILITY.md [offset=320, limit=50]
→ Read internal/realm/http.go [limit=130, offset=60]

$ cd /root/astrate-mule && echo "=== yaml-reading tests ===" && rg -ln "astrate_realm_management_api|docs/api" --glob '*.go' | head -30 && echo "=== UPSTREAM-EXPERIMENTAL async/device mentions ===" && rg -n "async_operation|delete.*device|device.*delet" docs/UPSTREAM-EXPERIMENTAL.md | head -40
=== yaml-reading tests ===
=== UPSTREAM-EXPERIMENTAL async/device mentions ===
25:| `async_operation=false` on mutating RM endpoints | #68 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |

$ cd /root/astrate-mule && ls docs && echo "=== Makefile in docs ===" && cat docs/Makefile 2>/dev/null | head -60
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
=== Makefile in docs ===
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

$ cd /root/astrate-mule && echo "=== go tests mentioning openapi/swagger/yaml spec ===" && rg -ln "openapi|swagger|astarte_realm_management" --glob '*.go' | head -20 && echo "=== git log for async_operation docs work ===" && git log --oneline -15 && echo "=== files touching docs/api/astrate_realm_management_api.yaml ===" && git log --oneline -8 -- docs/api/astarte_realm_management_api.yaml
=== go tests mentioning openapi/swagger/yaml spec ===
cmd/astrate/main.go
docs/embed.go
internal/swagger/swagger.go
internal/swagger/swagger_test.go
internal/realm/http.go
internal/realm/service.go
=== git log for async_operation docs work ===
7d78bd3 mule: log docs-sync-rm-async-operation-param
3f0eca8 mule: docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
17bfa74 mule: log docs-sync-hk-retention-zero-is-unset
c2ad442 mule: docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
57501a7 mule: log docs-sync-hk-async-operation-param
ed1cd72 mule: docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
19f4788 mule: refill the queue
ee48763 mule: recipe docs-sync ran (335s)
3b06c46 mule: log cmd-loadsealer-masterkeyfile-test
a6711f9 mule: cmd-loadsealer-masterkeyfile-test [auto]: add a container-free test for `loadSealer` (cmd/astrate/main.go:515-526) — the `security.master_key_file` to env hand-off (`os.Setenv(store.EnvMasterKeyFile, ...)`, 517) is the only wiring that makes a configured key file reach the store loader, and dropping or misspelling that branch leaves the file silently ignored and boot dying with the generic "master key" error (523) even though the config named a valid file. Assert a 32-byte key file under t.TempDir() loads, the inline ASTRATE_MASTER_KEY path still wins, and a missing key returns the error naming all three env vars. t.Setenv both ASTRATE_MASTER_KEY and ASTRATE_MASTER_KEY_FILE first so the os.Setenv is restored at cleanup.
b84abe5 mule: blocked cmd-healthcheck-contract-test
93da5b6 mule: log cmd-devcert-fields-test
f1434eb mule: cmd-devcert-fields-test [auto]: add a container-free test for `selfSignedDevCert` (cmd/astrate/main.go:320-345), the throwaway mTLS identity `insecure_dev_mode` boots on, which no test touches — pin NotBefore = now-1h (331; a plain now would reject a device whose clock is behind), NotAfter = now+365d (332), DNSNames localhost and IPAddresses 127.0.0.1/::1 (336-337, dropping which breaks hostname verification for a device dialling mqtts://127.0.0.1:8883), ExtKeyUsage serverAuth (334), and that the leaf is self-signed and parses (x509.ParseCertificate on Certificate[0]).
5cf927e mule: log drain-per-stage-budget
84a050e mule: drain-per-stage-budget [auto]: `shutdown` spends a single 30s sctx (cmd/astrate/main.go:479) across three sequential stages — srv.Shutdown (484), b.Close (489), flowSvc.Manager().Shutdown (494) and MarkRunningFlowsStopped (497) — while `drainEngine` alone opens a second one (504), so the constant's own comment "bounds the whole graceful drain" (main.go:55-56) is wrong in both directions, and one slow in-flight HTTP request (the server sets no per-request deadline, 221) burns the whole first budget and leaves the flow stage on an expired context: StopFlow then returns at `f.router.Drain(ctx)` before step 3 releases block resources (internal/flow/flow.go:265-273) and MarkRunningFlowsStopped skips every flow (flowapi/service.go:895-897), both surfacing as one log.Warn. Give each stage its own budget through an extracted helper (shutdownTimeout has to become an injectable var for the test) and pin it container-free: a second stage still receives a live context after a first stage exhausts its own.
=== files touching docs/api/astrate_realm_management_api.yaml ===
3f0eca8 mule: docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
9c05370 mule: docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
02e269b mule: docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "realm: validation failed: interface definition is invalid"`, but `validationDetail` strips the `realm: validation failed: ` prefix on the wire (http.go:562-568) and the parser's real messages are `invalid interface: ...` (pkg/interfaceschema/parse.go:38, wrapped at service.go:163); the example should show the stripped, real message. Same class as the already-fixed docs-sync-hk-validation-example (housekeeping http.go:251-258). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
93e784e mule: docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` (yaml:231-232) in docs/api/astarte_realm_management_api.yaml as a oneOf — both handlers go through `writeInterfaceError`, which answers three distinct 422 bodies: the flat ErrorDetail (`ErrValidation` via `validationDetail`, http.go:392-393), the nested violations changeset envelope (`writeViolations`, http.go:358-359 + 434-533, e.g. `{"errors":{"description":["should be at most 1000 character(s)"]}}` and the aligned full-length `mappings` array, http_test.go:480-501), and the named FieldErrors envelope for `ErrMaximumDatabaseRetentionExceeded` (`{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, http.go:355-357, http_test.go:282-320); today only the flat ValidationError is referenced. Mirror the createTrigger oneOf pattern (yaml:344-368). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
7256af9 mule: docs-sync-rm-put-interface-409 [auto]: add the missing `409` Conflict response to `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` in docs/api/astarte_realm_management_api.yaml (yaml:222-234 lists only 204/400/401/404/422/500) — `updateInterface` answers 409 for every url/body disagreement and minor-upgrade incompatibility: `ErrNameMismatch` (service.go:204 → "Interface name doesn't match the one in the interface json", http.go:367-369), `ErrMajorMismatch` (service.go:207 → http.go:370-372), and the `CheckMinorUpgrade` sentinels `ErrMinorNotIncreased`/`ErrDowngradeNotAllowed`/`ErrMissingEndpoints`/`ErrIncompatibleEndpointChange` (service.go:223-226 → http.go:373-384) plus `store.ErrAlreadyExists` (http.go:385-386) — the exact set `TestRealmManagementErrorCodes` pins as StatusConflict (http_test.go:417-437). Reuse the `Conflict` response component (already used on POST, yaml:131-132) and add name-mismatch/major-mismatch examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
8bfedd6 mule: docs-sync-rm-triggers-422-nested-envelope [auto]: document the nested changeset envelope on `POST /realmmanagement/v1/{realm}/triggers` in docs/api/astarte_realm_management_api.yaml — when trigger compile fails, `createTrigger` answers 422 with the upstream nested body `{"errors":{"action":{...},"simple_triggers":[...]}}` (internal/realm/http.go:252-258, triggerErrorBody at http.go:268-285, issues #63/#70), but the spec's createTrigger 422 (yaml:344-345) refs only the flat `ValidationError` component (`{"errors":{"detail":...}}`) — which the handler also emits for plain ErrValidation (http.go:259-260, e.g. "trigger requires a name", service.go:429). Model the response as oneOf across the flat ErrorDetail and the nested envelope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
c0521bc mule: docs-sync-rm-policies-delete-422 [auto]: add the missing `422` response to `DELETE /realmmanagement/v1/{realm}/policies/{name}` in docs/api/astarte_realm_management_api.yaml — `Service.DeletePolicy` refuses to delete a policy still referenced by a trigger (`"policy %q is still used by trigger %q"`, internal/realm/service.go:573-574, wrapped as ErrValidation) and `writeError` answers 422 for it (internal/realm/http.go:410-411, detail stripped by validationDetail), but the spec lists only 204/401/404/500 for the delete op (yaml:695-713). Add the 422 ValidationError response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
8554720 mule: docs-sync-rm-put-auth-422 [auto]: add the missing `422` (ValidationError) response to `PUT /realmmanagement/v1/{realm}/config/auth` in docs/api/astarte_realm_management_api.yaml — `SetAuthKey` returns `ErrValidation` for a blank `jwt_public_key_pem` (internal/realm/service.go:632-640), which `writeError` maps to 422; the spec (yaml:431-464) documents only 204/400/401/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cd /root/astrate-mule && wc -l internal/swagger/swagger_test.go internal/swagger/swagger.go && rg -n "func Test" internal/swagger/swagger_test.go
435 internal/swagger/swagger_test.go
   69 internal/swagger/swagger.go
  504 total
18:func TestMount(t *testing.T) {
77:func TestMountSubPanicsOnBrokenFS(t *testing.T) {
90:func TestSpecs(t *testing.T) {
123:func TestRealmManagement403(t *testing.T) {
179:func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
222:func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
294:func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {

→ Read internal/swagger/swagger_test.go [limit=125, offset=170]

$ cd /root/astrate-mule && rg -n "func newRig|type rig|func \(r \*rig\)" internal/realm/*_test.go | head -20 && echo "=== rig file ===" && rg -ln "func newRig" internal/realm/
internal/realm/http_test.go:44:type rig struct {
internal/realm/http_test.go:56:func newRig(t *testing.T) *rig {
internal/realm/http_test.go:102:func (r *rig) req(t *testing.T, method, path, rawBody, token string) *httptest.ResponseRecorder {
=== rig file ===
internal/realm/http_test.go

→ Read internal/realm/http_test.go [limit=120, offset=1]

$ cd /root/astrate-mule && rg -n "devices" internal/realm/http_test.go | head -30 && echo "=== store device create helpers ===" && rg -n "func (.*Store) (UpsertDevice|CreateDevice|GetDevice|SetDevice)" internal/store/*.go | head
=== store device create helpers ===

$ cd /root/astrate-mule && rg -rn "devices" internal/realm/ | head -20 && echo "=== deleteDevice callers/tests ===" && rg -n "DeleteDevice|deleteDevice" --glob '*.go' | head -20
internal/realm/dashboard_compat_test.go:		if rec := r.req(t, http.MethodDelete, "/n/"+dev.String(), "", r.rmaToken); rec.Code != http.StatusNoContent {
internal/realm/dashboard_compat_test.go:		for _, path := range []string{"/n/" + dev.String(), "/n/not-an-id"} {
internal/realm/http.go:	mux.Handle("DELETE /realmmanagement/v1/{realm}/n/{device}", h(a.deleteDevice))
=== deleteDevice callers/tests ===
internal/store/devices.go:453:// DeleteDevice removes the device row and every stored datum in one
internal/store/devices.go:459:func (s *Store) DeleteDevice(ctx context.Context, realmID int16, id deviceid.ID) error {
internal/store/devices_test.go:358:		if err := s.DeleteDevice(ctx, realm.ID, d1); err != nil {
internal/store/devices_test.go:359:			t.Fatalf("DeleteDevice: %v", err)
internal/store/devices_test.go:372:		if err := s.DeleteDevice(ctx, realm.ID, d1); !errors.Is(err, ErrNotFound) {
internal/realm/http.go:53:	mux.Handle("DELETE /realmmanagement/v1/{realm}/devices/{device}", h(a.deleteDevice))
internal/realm/http.go:102:// deleteDevice synchronously removes a device and its data (the dashboard's
internal/realm/http.go:105:func (a *API) deleteDevice(w http.ResponseWriter, r *http.Request) {
internal/realm/http.go:106:	err := a.svc.DeleteDevice(r.Context(), r.PathValue("realm"), r.PathValue("device"))
internal/realm/service.go:83:	// OnDeletionStart / OnDeletionFinish bookend DeleteDevice (nil-safe).
internal/realm/service.go:105:// DeleteDevice synchronously removes a device and all its data (upstream
internal/realm/service.go:113:func (s *Service) DeleteDevice(ctx context.Context, realm, deviceID string) error {
internal/realm/service.go:129:	err = s.st.DeleteDevice(ctx, rid, id)

→ Read internal/realm/dashboard_compat_test.go
→ Read internal/realm/http.go [limit=59]
→ Read docs/UPSTREAM-EXPERIMENTAL.md [limit=40, offset=10]

$ cd /root/astrate-mule && ssh -o ConnectTimeout=5 -o BatchMode=yes legion 'echo ok' 2>&1 | tail -2
ssh: connect to host legion.tail4db093.ts.net port 22: Connection timed out

→ Read internal/swagger/swagger_test.go [limit=70, offset=100]

$ cd /root/astrate-mule && rg -n -A24 "^func operationBlock" internal/swagger/swagger_test.go && rg -n -A12 "^func componentBlock|^func containsLine" internal/swagger/swagger_test.go
345:func operationBlock(t *testing.T, lines []string, operationID string) []string {
346-	t.Helper()
347-	marker := "      operationId: " + operationID
348-	start := -1
349-	for i, l := range lines {
350-		if l == marker {
351-			start = i
352-			break
353-		}
354-	}
355-	if start < 0 {
356-		t.Fatalf("spec declares no operation %q", operationID)
357-	}
358-	endRe := regexp.MustCompile(`^      operationId: |^  \S|^components:`)
359-	for i := start + 1; i < len(lines); i++ {
360-		if endRe.MatchString(lines[i]) {
361-			return lines[start+1 : i]
362-		}
363-	}
364-	return lines[start+1:]
365-}
366-
367-// componentBlock returns the lines of the components entry starting with the
368-// given key line, up to the next entry at the same or shallower indentation.
369-func componentBlock(t *testing.T, lines []string, key string) []string {
369:func componentBlock(t *testing.T, lines []string, key string) []string {
370-	t.Helper()
371-	start := -1
372-	for i, l := range lines {
373-		if l == key {
374-			start = i
375-			break
376-		}
377-	}
378-	if start < 0 {
379-		t.Fatalf("spec declares no components entry %q", key)
380-	}
381-	for i := start + 1; i < len(lines); i++ {
--
392:func containsLine(lines []string, want string) bool {
393-	for _, l := range lines {
394-		if l == want {
395-			return true
396-		}
397-	}
398-	return false
399-}
400-
401-func embeddedYAMLFilenames() ([]string, error) {
402-	var names []string
403-	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
404-		if err != nil {

Now the edits. YAML first:
← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -693,8 +693,9 @@
   - a_rma: []
 parameters:
   - $ref: "#/components/parameters/RealmName"
   - $ref: "#/components/parameters/DeviceID"
+  - $ref: "#/components/parameters/AsyncOperation"
 responses:
   "204":
     description: Device deleted. No body returned.
   "401":

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -905,13 +905,14 @@
 in: query
 required: false
 description: |
   Accepted and ignored on either value. Upstream 1.4 runs interface
-  install/update/delete and trigger-delivery-policy delete in the
-  background and lets the caller opt into synchronous execution with
-  `async_operation=false`; Astrate always works synchronously and answers
-  only once the work is done, so the value changes nothing and the
-  response is identical either way (deviation 17 in COMPATIBILITY.md).
+  install/update/delete, device deletion and trigger-delivery-policy
+  delete in the background and lets the caller opt into synchronous
+  execution with `async_operation=false`; Astrate always works
+  synchronously and answers only once the work is done, so the value
+  changes nothing and the response is identical either way
+  (deviation 17 in COMPATIBILITY.md).
 schema:
   type: boolean
   default: false
 

← Edit docs/COMPATIBILITY.md
Index: /root/astrate-mule/docs/COMPATIBILITY.md
===================================================================
--- /root/astrate-mule/docs/COMPATIBILITY.md
+++ /root/astrate-mule/docs/COMPATIBILITY.md
@@ -333,16 +333,17 @@
     log-and-drop).
 
 17. **Always synchronous where upstream 1.4 defaults to asynchronous** —
     upstream 1.4 runs realm create/delete (Housekeeping), interface
-    install/update/delete (Realm Management) and trigger-delivery-policy
-    delete in the background, answering immediately and letting the caller opt
-    into synchronous execution with `?async_operation=false`. Astrate performs
-    every one of these synchronously and answers only once the work is done —
-    a strictly stronger guarantee, so a client written against upstream never
-    observes a regression. The parameter is accepted and ignored on either
-    value, on every one of those endpoints, so upstream clients that send it
-    keep working unchanged; pinned by
+    install/update/delete and device deletion (Realm Management) and
+    trigger-delivery-policy delete in the background, answering immediately
+    and letting the caller opt into synchronous execution with
+    `?async_operation=false`. Astrate performs every one of these
+    synchronously and answers only once the work is done — a strictly
+    stronger guarantee, so a client written against upstream never observes a
+    regression. The parameter is accepted and ignored on either value, on
+    every one of those endpoints, so upstream clients that send it keep
+    working unchanged; pinned by
     `TestRealmManagementAsyncOperationParam` and
     `TestHousekeepingAsyncOperationParam`. Tracked in
     `docs/UPSTREAM-EXPERIMENTAL.md` (#68) until upstream 1.4 is final.
 

Now the docs-half test:
✗ Edit internal/swagger/swagger_test.go failed
Error: Could not find oldString in the file. It must match exactly, including whitespace, indentation, and line endings.
✗ Edit internal/swagger/swagger_test.go failed
Error: Could not find oldString in the file. It must match exactly, including whitespace, indentation, and line endings.

$ cd /root/astrate-mule && sed -n '212,236p' internal/swagger/swagger_test.go | cat -A | sed -n '1,25p' | cut -c1-120
// TestRealmManagementAsyncOperationParamDocumented guards that the realm$
// management spec tells clients the `?async_operation` parameter exists on the$
// four operations upstream 1.4 runs in the background: interface$
// install/update/delete and trigger-delivery-policy delete. Astrate accepts and$
// ignores it there (deviation 17, docs/COMPATIBILITY.md) M-bM-^@M-^T none of the four$
// handlers reads the query string, which is the whole of the acceptance M-bM-^@M-^T so an$
// upstream client keeps working, but a client generated from the spec has to$
// learn the parameter from the spec. This is the documentation half of the$
// behaviour pinned by TestRealmManagementAsyncOperationParam in$
// internal/realm.$
func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {$
^Ib, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")$
^Iif err != nil {$
^I^It.Fatalf("reading astarte_realm_management_api.yaml: %v", err)$
^I}$
^Ilines := strings.Split(string(b), "\n")$
$
^Iconst ref = `        - $ref: "#/components/parameters/AsyncOperation"`$
^Ifor _, op := range []string{"installInterface", "updateInterface", "deleteInterface", "deletePolicy"} {$
^I^Iblock := operationBlock(t, lines, op)$
^I^Iif !containsLine(block, ref) {$
^I^I^It.Errorf("operation %s does not $ref the AsyncOperation parameter", op)$
^I^I}$
^I}$
$

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -210,15 +210,15 @@
 }
 
 // TestRealmManagementAsyncOperationParamDocumented guards that the realm
 // management spec tells clients the `?async_operation` parameter exists on the
-// four operations upstream 1.4 runs in the background: interface
-// install/update/delete and trigger-delivery-policy delete. Astrate accepts and
-// ignores it there (deviation 17, docs/COMPATIBILITY.md) — none of the four
-// handlers reads the query string, which is the whole of the acceptance — so an
-// upstream client keeps working, but a client generated from the spec has to
-// learn the parameter from the spec. This is the documentation half of the
-// behaviour pinned by TestRealmManagementAsyncOperationParam in
+// five operations upstream 1.4 runs in the background: interface
+// install/update/delete, device deletion and trigger-delivery-policy delete.
+// Astrate accepts and ignores it there (deviation 17, docs/COMPATIBILITY.md) —
+// none of the five handlers reads the query string, which is the whole of the
+// acceptance — so an upstream client keeps working, but a client generated from
+// the spec has to learn the parameter from the spec. This is the documentation
+// half of the behaviour pinned by TestRealmManagementAsyncOperationParam in
 // internal/realm.
 func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
 	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
 	if err != nil {

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -226,9 +226,12 @@
 }
 lines := strings.Split(string(b), "\n")
 
 const ref = `        - $ref: "#/components/parameters/AsyncOperation"`
-for _, op := range []string{"installInterface", "updateInterface", "deleteInterface", "deletePolicy"} {
+for _, op := range []string{
+	"installInterface", "updateInterface", "deleteInterface",
+	"deleteDevice", "deletePolicy",
+} {
 	block := operationBlock(t, lines, op)
 	if !containsLine(block, ref) {
 		t.Errorf("operation %s does not $ref the AsyncOperation parameter", op)
 	}

Now the behaviour half — the device-delete case in the integration test:
← Edit internal/realm/async_operation_test.go
Index: /root/astrate-mule/internal/realm/async_operation_test.go
===================================================================
--- /root/astrate-mule/internal/realm/async_operation_test.go
+++ /root/astrate-mule/internal/realm/async_operation_test.go
@@ -2,21 +2,25 @@
 
 package realm
 
 import (
+	"context"
 	"fmt"
 	"net/http"
 	"testing"
+
+	"github.com/astrate-platform/astrate/pkg/deviceid"
 )
 
 // TestRealmManagementAsyncOperationParam pins deviation 17
-// (docs/COMPATIBILITY.md): upstream 1.4 runs interface install/update/delete
-// and policy delete asynchronously unless the caller opts out with
-// `?async_operation=false`. Astrate is always synchronous, so the parameter is
-// accepted and ignored on either value — sending it must never draw a 4xx, and
-// the effect must be visible on the very next read.
+// (docs/COMPATIBILITY.md): upstream 1.4 runs interface install/update/delete,
+// device deletion and policy delete asynchronously unless the caller opts out
+// with `?async_operation=false`. Astrate is always synchronous, so the
+// parameter is accepted and ignored on either value — sending it must never
+// draw a 4xx, and the effect must be visible on the very next read.
 func TestRealmManagementAsyncOperationParam(t *testing.T) {
 	r := newRig(t)
+	ctx := context.Background()
 
 	for i, q := range []string{"?async_operation=false", "?async_operation=true"} {
 		name := fmt.Sprintf("com.ex.M7a.Async%d", i)
 		iface := func(major, minor int, extra string) string {

← Edit internal/realm/async_operation_test.go
Index: /root/astrate-mule/internal/realm/async_operation_test.go
===================================================================
--- /root/astrate-mule/internal/realm/async_operation_test.go
+++ /root/astrate-mule/internal/realm/async_operation_test.go
@@ -67,6 +67,23 @@
 		}
 		if rec := r.req(t, http.MethodGet, "/policies/"+policyName, "", r.rmaToken); rec.Code != http.StatusNotFound {
 			t.Errorf("get policy after delete %s: got %d, want 404", q, rec.Code)
 		}
+
+		// Device delete carries the parameter upstream too: the dashboard's
+		// Delete button, which upstream answers before the data is gone.
+		dev, err := deviceid.Random()
+		if err != nil {
+			t.Fatalf("random device id: %v", err)
+		}
+		if err := r.st.RegisterDevice(ctx, r.realmID, dev, "h"); err != nil {
+			t.Fatalf("register device %s: %v", dev, err)
+		}
+		if rec := r.req(t, http.MethodDelete, "/devices/"+dev.String()+q, "", r.rmaToken); rec.Code != http.StatusNoContent {
+			t.Fatalf("delete device %s: got %d, want 204 (%s)", q, rec.Code, rec.Body)
+		}
+		// Synchronous means gone by the time the 204 is written, on both values.
+		if _, err := r.st.GetDevice(ctx, r.realmID, dev); err == nil {
+			t.Errorf("device row survived the delete %s", q)
+		}
 	}
 }

Now verify: the docs-half test must fail without the yaml change, and the whole gate must pass.
$ cd /root/astrate-mule && gofmt -l . && echo "--- gofmt clean above (no output = clean) ---" && go vet ./... && echo "VET OK" && go test ./internal/swagger/ -run 'AsyncOperationParamDocumented' -v 2>&1 | tail -20
--- gofmt clean above (no output = clean) ---
VET OK
=== RUN   TestHousekeepingAsyncOperationParamDocumented
--- PASS: TestHousekeepingAsyncOperationParamDocumented (0.00s)
=== RUN   TestRealmManagementAsyncOperationParamDocumented
--- PASS: TestRealmManagementAsyncOperationParamDocumented (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.016s

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -693,9 +693,8 @@
   - a_rma: []
 parameters:
   - $ref: "#/components/parameters/RealmName"
   - $ref: "#/components/parameters/DeviceID"
-  - $ref: "#/components/parameters/AsyncOperation"
 responses:
   "204":
     description: Device deleted. No body returned.
   "401":

$ cd /root/astrate-mule && go test ./internal/swagger/ -run 'TestRealmManagementAsyncOperationParamDocumented' 2>&1 | tail -8
--- FAIL: TestRealmManagementAsyncOperationParamDocumented (0.01s)
    swagger_test.go:236: operation deleteDevice does not $ref the AsyncOperation parameter
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.027s
FAIL

✗ Edit docs/api/astrate_realm_management_api.yaml failed
Error: File /root/astrate-mule/docs/api/astrate_realm_management_api.yaml not found

$ cd /root/astrate-mule && ls -la docs/api/ && rg -n "AsyncOperation" docs/api/astrate_realm_management_api.yaml
total 164
drwxr-xr-x 2 root root  4096 Sep  4 18:55 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 56583 Sep 12 22:44 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 15547 Sep 25 21:48 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 19506 Sep 17 13:03 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 45330 Sep 25 22:07 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 14977 Sep 20 21:25 astrate_native_api.yaml
rg: docs/api/astrate_realm_management_api.yaml: IO error for operation on docs/api/astrate_realm_management_api.yaml: No such file or directory (os error 2)

$ cd /root/astrate-mule && git status --short && echo "---" && rg -n "AsyncOperation" docs/api/astrate_realm_management_api.yaml
M docs/COMPATIBILITY.md
 M docs/api/astarte_realm_management_api.yaml
 M internal/realm/async_operation_test.go
 M internal/swagger/swagger_test.go
---
rg: docs/api/astrate_realm_management_api.yaml: IO error for operation on docs/api/astrate_realm_management_api.yaml: No such file or directory (os error 2)

The real filename is `astarte_…`; my restore path had a typo. Fixing now:
← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -693,8 +693,9 @@
   - a_rma: []
 parameters:
   - $ref: "#/components/parameters/RealmName"
   - $ref: "#/components/parameters/DeviceID"
+  - $ref: "#/components/parameters/AsyncOperation"
 responses:
   "204":
     description: Device deleted. No body returned.
   "401":

$ cd /root/astrate-mule && git diff --stat && echo "=== yaml diff ===" && git diff docs/api/astarte_realm_management_api.yaml
docs/COMPATIBILITY.md                      | 17 ++++++++--------
 docs/api/astarte_realm_management_api.yaml | 12 +++++++-----
 internal/realm/async_operation_test.go     | 31 +++++++++++++++++++++++++-----
 internal/swagger/swagger_test.go           | 19 ++++++++++--------
 4 files changed, 53 insertions(+), 26 deletions(-)
=== yaml diff ===
diff --git a/docs/api/astarte_realm_management_api.yaml b/docs/api/astarte_realm_management_api.yaml
index 5a229f6..d6886f4 100644
--- a/docs/api/astarte_realm_management_api.yaml
+++ b/docs/api/astarte_realm_management_api.yaml
@@ -694,6 +694,7 @@ paths:
       parameters:
         - $ref: "#/components/parameters/RealmName"
         - $ref: "#/components/parameters/DeviceID"
+        - $ref: "#/components/parameters/AsyncOperation"
       responses:
         "204":
           description: Device deleted. No body returned.
@@ -905,11 +906,12 @@ components:
       required: false
       description: |
         Accepted and ignored on either value. Upstream 1.4 runs interface
-        install/update/delete and trigger-delivery-policy delete in the
-        background and lets the caller opt into synchronous execution with
-        `async_operation=false`; Astrate always works synchronously and answers
-        only once the work is done, so the value changes nothing and the
-        response is identical either way (deviation 17 in COMPATIBILITY.md).
+        install/update/delete, device deletion and trigger-delivery-policy
+        delete in the background and lets the caller opt into synchronous
+        execution with `async_operation=false`; Astrate always works
+        synchronously and answers only once the work is done, so the value
+        changes nothing and the response is identical either way
+        (deviation 17 in COMPATIBILITY.md).
       schema:
         type: boolean
         default: false

Now the full gate plus a compile check of the integration-tagged test:
$ cd /root/astrate-mule && go vet -tags integration ./internal/realm/ && echo "INTEGRATION COMPILE OK" && gofmt -l . && echo "GOFMT CLEAN" && go vet ./... && echo "VET OK"
INTEGRATION COMPILE OK
GOFMT CLEAN
VET OK

$ cd /root/astrate-mule && go test ./... 2>&1 | tail -45
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
ok  	github.com/astrate-platform/astrate/internal/swagger	0.072s
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

$ cd /root/astrate-mule && (command -v mkdocs && mkdocs --version) 2>&1 | tail -3
(no output)

$ cd /root/astrate-mule && timeout 300 make -C docs build 2>&1 | tail -15
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

$ cd /root/astrate-mule && git status --short && echo "=== site tracked? ===" && git check-ignore -v docs/site/DESIGN.md 2>&1 | tail -2; python3 -c "import yaml; print('pyyaml', yaml.__version__)" 2>&1 | tail -2; rg -n "yaml" go.mod
M docs/COMPATIBILITY.md
 M docs/api/astarte_realm_management_api.yaml
 M internal/realm/async_operation_test.go
 M internal/swagger/swagger_test.go
=== site tracked? ===
.gitignore:32:/docs/site/DESIGN.md	docs/site/DESIGN.md
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

$ cd /root/astrate-mule && (command -v ruby && ruby -ryaml -e 'puts "ruby-yaml ok"') 2>&1 | tail -2; echo "=== modcache ==="; ls "$(go env GOMODCACHE)/gopkg.in" 2>/dev/null | head; ls -d "$(go env GOMODCACHE)/gopkg.in/yaml.v3@v3.0.1" 2>/dev/null
=== modcache ===
check.v1@v0.0.0-20161208181325-20d25e280405
check.v1@v1.0.0-20201130134442-10cb98267c6c
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1
/root/go/pkg/mod/gopkg.in/yaml.v3@v3.0.1

$ cd /root/astrate-mule && rg -n "realm_management|url|yaml" docs/swagger-ui/* 2>/dev/null | head -20 && echo "=== swagger.go ===" && rg -n "Specs|APIYAML|yaml" internal/swagger/swagger.go | head -20
44:      <option value="../api/astarte_housekeeping_api.yaml">Housekeeping API — Realm management, a_ha auth</option>
45:      <option value="../api/astarte_pairing_api.yaml">Pairing API — Device pairing flows A/B/C, a_pa auth</option>
46:      <option value="../api/astarte_realm_management_api.yaml">Realm Management API — Interfaces, triggers, policies, a_rma auth</option>
47:      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
48:      <option value="../api/astrate_native_api.yaml">Astrate Native API — Health, readiness, metrics, WebSocket</option>
51:    <input type="text" id="server-url" value="http://localhost:8080" placeholder="e.g. http://localhost:8080">
59:    const serverUrl = document.getElementById('server-url');
65:        url: specUrl,
83:      const url = serverUrl.value.replace(/\/+$/, '');
89:        nativeSetter.call(input, url);
=== swagger.go ===
19:	apiRoot := mustSub(docs.APIYAML, "api")
35:	// Serve YAML files at /api/ so the relative ../api/*.yaml in index.html
56:// Specs returns the list of available YAML spec filenames (without path prefix).
57:func Specs() []string {
59:	_ = fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
63:		if strings.HasSuffix(path, ".yaml") {

$ mkdir -p /tmp/opencode/yamlprobe && cd /tmp/opencode/yamlprobe && printf 'module yamlprobe\n\ngo 1.24\n\nrequire gopkg.in/yaml.v3 v3.0.1\n\nreplace gopkg.in/yaml.v3 => /root/go/pkg/mod/gopkg.in/yaml.v3@v3.0.1\n' > go.mod && cat go.mod
module yamlprobe

go 1.24

require gopkg.in/yaml.v3 v3.0.1

replace gopkg.in/yaml.v3 => /root/go/pkg/mod/gopkg.in/yaml.v3@v3.0.1

← Write /tmp/opencode/yamlprobe/main.go
Wrote file successfully.

$ cd /tmp/opencode/yamlprobe && GOFLAGS=-mod=mod GOPROXY=off go run . 2>&1 | tail -20
YAML PARSE OK
paths: 12 operations: 20
AsyncOperation: name=async_operation in=query required=false type=boolean default=false
description mentions device deletion: true
  deleteDevice param: #/components/parameters/RealmName
  deleteDevice param: #/components/parameters/DeviceID
  deleteDevice param: #/components/parameters/AsyncOperation
RESULT: deleteDevice $refs AsyncOperation and the component resolves

$ cd /root/astrate-mule/docs && rg -n "api|nav|docs_dir" mkdocs.yml | head -25 && echo "=== requirements ===" && cat requirements.txt
5:docs_dir: site
24:    - navigation.instant
25:    - navigation.sections
26:    - navigation.top
42:nav:
50:  - AppEngine API: appengine-api.md
51:  - Realm Management API: realm-management-api.md
52:  - Housekeeping API: housekeeping-api.md
=== requirements ===
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10

$ cd /root/astrate-mule && timeout 90 pip install -q -r docs/requirements.txt 2>&1 | tail -5; command -v mkdocs || echo "mkdocs still absent"
timeout: failed to run command ‘pip’: No such file or directory
mkdocs still absent

$ cd /root/astrate-mule && go test -count=1 ./internal/swagger/ ./internal/realm/ 2>&1 | tail -5 && echo "=== FULL DIFF ===" && git diff
ok  	github.com/astrate-platform/astrate/internal/swagger	0.088s
ok  	github.com/astrate-platform/astrate/internal/realm	0.042s
=== FULL DIFF ===
diff --git a/docs/COMPATIBILITY.md b/docs/COMPATIBILITY.md
index e34dd45..879ee9b 100644
--- a/docs/COMPATIBILITY.md
+++ b/docs/COMPATIBILITY.md
@@ -334,14 +334,15 @@ All additive or strictly-safer; none affect unmodified device SDKs.
 
 17. **Always synchronous where upstream 1.4 defaults to asynchronous** —
     upstream 1.4 runs realm create/delete (Housekeeping), interface
-    install/update/delete (Realm Management) and trigger-delivery-policy
-    delete in the background, answering immediately and letting the caller opt
-    into synchronous execution with `?async_operation=false`. Astrate performs
-    every one of these synchronously and answers only once the work is done —
-    a strictly stronger guarantee, so a client written against upstream never
-    observes a regression. The parameter is accepted and ignored on either
-    value, on every one of those endpoints, so upstream clients that send it
-    keep working unchanged; pinned by
+    install/update/delete and device deletion (Realm Management) and
+    trigger-delivery-policy delete in the background, answering immediately
+    and letting the caller opt into synchronous execution with
+    `?async_operation=false`. Astrate performs every one of these
+    synchronously and answers only once the work is done — a strictly
+    stronger guarantee, so a client written against upstream never observes a
+    regression. The parameter is accepted and ignored on either value, on
+    every one of those endpoints, so upstream clients that send it keep
+    working unchanged; pinned by
     `TestRealmManagementAsyncOperationParam` and
     `TestHousekeepingAsyncOperationParam`. Tracked in
     `docs/UPSTREAM-EXPERIMENTAL.md` (#68) until upstream 1.4 is final.
diff --git a/docs/api/astarte_realm_management_api.yaml b/docs/api/astarte_realm_management_api.yaml
index 5a229f6..d6886f4 100644
--- a/docs/api/astarte_realm_management_api.yaml
+++ b/docs/api/astarte_realm_management_api.yaml
@@ -694,6 +694,7 @@ paths:
       parameters:
         - $ref: "#/components/parameters/RealmName"
         - $ref: "#/components/parameters/DeviceID"
+        - $ref: "#/components/parameters/AsyncOperation"
       responses:
         "204":
           description: Device deleted. No body returned.
@@ -905,11 +906,12 @@ components:
       required: false
       description: |
         Accepted and ignored on either value. Upstream 1.4 runs interface
-        install/update/delete and trigger-delivery-policy delete in the
-        background and lets the caller opt into synchronous execution with
-        `async_operation=false`; Astrate always works synchronously and answers
-        only once the work is done, so the value changes nothing and the
-        response is identical either way (deviation 17 in COMPATIBILITY.md).
+        install/update/delete, device deletion and trigger-delivery-policy
+        delete in the background and lets the caller opt into synchronous
+        execution with `async_operation=false`; Astrate always works
+        synchronously and answers only once the work is done, so the value
+        changes nothing and the response is identical either way
+        (deviation 17 in COMPATIBILITY.md).
       schema:
         type: boolean
         default: false
diff --git a/internal/realm/async_operation_test.go b/internal/realm/async_operation_test.go
index 1560203..a957bc6 100644
--- a/internal/realm/async_operation_test.go
+++ b/internal/realm/async_operation_test.go
@@ -3,19 +3,23 @@
 package realm
 
 import (
+	"context"
 	"fmt"
 	"net/http"
 	"testing"
+
+	"github.com/astrate-platform/astrate/pkg/deviceid"
 )
 
 // TestRealmManagementAsyncOperationParam pins deviation 17
-// (docs/COMPATIBILITY.md): upstream 1.4 runs interface install/update/delete
-// and policy delete asynchronously unless the caller opts out with
-// `?async_operation=false`. Astrate is always synchronous, so the parameter is
-// accepted and ignored on either value — sending it must never draw a 4xx, and
-// the effect must be visible on the very next read.
+// (docs/COMPATIBILITY.md): upstream 1.4 runs interface install/update/delete,
+// device deletion and policy delete asynchronously unless the caller opts out
+// with `?async_operation=false`. Astrate is always synchronous, so the
+// parameter is accepted and ignored on either value — sending it must never
+// draw a 4xx, and the effect must be visible on the very next read.
 func TestRealmManagementAsyncOperationParam(t *testing.T) {
 	r := newRig(t)
+	ctx := context.Background()
 
 	for i, q := range []string{"?async_operation=false", "?async_operation=true"} {
 		name := fmt.Sprintf("com.ex.M7a.Async%d", i)
@@ -64,5 +68,22 @@ func TestRealmManagementAsyncOperationParam(t *testing.T) {
 		if rec := r.req(t, http.MethodGet, "/policies/"+policyName, "", r.rmaToken); rec.Code != http.StatusNotFound {
 			t.Errorf("get policy after delete %s: got %d, want 404", q, rec.Code)
 		}
+
+		// Device delete carries the parameter upstream too: the dashboard's
+		// Delete button, which upstream answers before the data is gone.
+		dev, err := deviceid.Random()
+		if err != nil {
+			t.Fatalf("random device id: %v", err)
+		}
+		if err := r.st.RegisterDevice(ctx, r.realmID, dev, "h"); err != nil {
+			t.Fatalf("register device %s: %v", dev, err)
+		}
+		if rec := r.req(t, http.MethodDelete, "/devices/"+dev.String()+q, "", r.rmaToken); rec.Code != http.StatusNoContent {
+			t.Fatalf("delete device %s: got %d, want 204 (%s)", q, rec.Code, rec.Body)
+		}
+		// Synchronous means gone by the time the 204 is written, on both values.
+		if _, err := r.st.GetDevice(ctx, r.realmID, dev); err == nil {
+			t.Errorf("device row survived the delete %s", q)
+		}
 	}
 }
diff --git a/internal/swagger/swagger_test.go b/internal/swagger/swagger_test.go
index a2d357a..8f1e81b 100644
--- a/internal/swagger/swagger_test.go
+++ b/internal/swagger/swagger_test.go
@@ -211,13 +211,13 @@ func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
 
 // TestRealmManagementAsyncOperationParamDocumented guards that the realm
 // management spec tells clients the `?async_operation` parameter exists on the
-// four operations upstream 1.4 runs in the background: interface
-// install/update/delete and trigger-delivery-policy delete. Astrate accepts and
-// ignores it there (deviation 17, docs/COMPATIBILITY.md) — none of the four
-// handlers reads the query string, which is the whole of the acceptance — so an
-// upstream client keeps working, but a client generated from the spec has to
-// learn the parameter from the spec. This is the documentation half of the
-// behaviour pinned by TestRealmManagementAsyncOperationParam in
+// five operations upstream 1.4 runs in the background: interface
+// install/update/delete, device deletion and trigger-delivery-policy delete.
+// Astrate accepts and ignores it there (deviation 17, docs/COMPATIBILITY.md) —
+// none of the five handlers reads the query string, which is the whole of the
+// acceptance — so an upstream client keeps working, but a client generated from
+// the spec has to learn the parameter from the spec. This is the documentation
+// half of the behaviour pinned by TestRealmManagementAsyncOperationParam in
 // internal/realm.
 func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
 	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
@@ -227,7 +227,10 @@ func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
 	lines := strings.Split(string(b), "\n")
 
 	const ref = `        - $ref: "#/components/parameters/AsyncOperation"`
-	for _, op := range []string{"installInterface", "updateInterface", "deleteInterface", "deletePolicy"} {
+	for _, op := range []string{
+		"installInterface", "updateInterface", "deleteInterface",
+		"deleteDevice", "deletePolicy",
+	} {
 		block := operationBlock(t, lines, op)
 		if !containsLine(block, ref) {
 			t.Errorf("operation %s does not $ref the AsyncOperation parameter", op)

$ cd /root/astrate-mule && rg -n "async_operation|AsyncOperation" --glob '!docs/site/**' --glob '!*.sum' -l && echo "=== other count phrasings ===" && rg -n "four operations|four mutating|mutating RM" --glob '!docs/site/**' | head
docs/api/astarte_housekeeping_api.yaml
docs/api/astarte_realm_management_api.yaml
docs/COMPATIBILITY.md
docs/UPSTREAM-EXPERIMENTAL.md
internal/swagger/swagger_test.go
internal/housekeeping/async_operation_test.go
internal/realm/async_operation_test.go
=== other count phrasings ===
docs/UPSTREAM-EXPERIMENTAL.md:25:| `async_operation=false` on mutating RM endpoints | #68 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |

$ cd /root/astrate-mule && cat .mule/todo.md
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_0da309a0f001J8Z1xg6Zh5nqrk

- [x] docs-sync-rm-delete-interface-status [auto]: fix the `deleteInterface` response in docs/api/astarte_realm_management_api.yaml — the spec documents `422` with examples "Interface major version is not 0, can't be deleted" / "Cannot delete an interface that is used by a device introspection" (yaml:255-272), but the handler returns `403` with "Interface can't be deleted" / "Interface can't be deleted since it's currently used" (internal/realm/http.go:387-391, `writeInterfaceError`; both are 403 `StatusForbidden`). The spec's detail strings were lifted from the `writeError` branch (http.go:414-419) which DELETE /interfaces does not use. Fix status to `403` and update the two examples to match the actual wire messages. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-mapping-required-encrypted [auto]: add `required` (boolean, datastream-only) and `encrypted` (boolean, datastream-only) to the `InterfaceMapping` schema in docs/api/astarte_realm_management_api.yaml — the parser accepts and persists both on datastream mappings (internal/interfaceschema/parse.go:410-414), rejects them on properties mappings (parse.go:395-398), the detailed listing emits them (`detailed.go:69`), and the compatibility check compares them (interfaceschema/compat.go:121-123); the schema (yaml:820-862) currently has `additionalProperties: false` and neither field, so the OpenAPI would mark a valid datastream interface install body as invalid. Also add both fields to the detailed-listing example (yaml:70-85). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-put-auth-422 [auto]: add the missing `422` (ValidationError) response to `PUT /realmmanagement/v1/{realm}/config/auth` in docs/api/astarte_realm_management_api.yaml — `SetAuthKey` returns `ErrValidation` for a blank `jwt_public_key_pem` (internal/realm/service.go:632-640), which `writeError` maps to 422; the spec (yaml:431-464) documents only 204/400/401/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-rm-version-example [auto]: fix the `getVersion` example in docs/api/astarte_realm_management_api.yaml from `data: "1.1.0"` (yaml:557) to `data: "1.2.2"` — `APICompatVersion` is `"1.2.2"` (internal/realm/service.go:588); the native spec's four compat-version examples already use "1.2.2". Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1
- [x] interfaceschema-compat-attrs-flip-test [auto]: extend TestCheckMinorUpgradeTable (pkg/interfaceschema/compat_test.go) to assert the four currently-unpinned immutable attributes — retention, expiry, database_retention_policy/ttl (datastream old/next pair) and allow_unset (properties old/next pair) — each flipped on upgrade and rejected as ErrIncompatibleEndpointChange with an at-immutability acceptance twin; only type/reliability/explicit_timestamp/required/encrypted are tested today (sameMappingAttributes, compat.go:102, sits at 66.7% stmt coverage), so a regression relaxing expiry or DB-TTL immutability passes the suite.
- [x] interfaceschema-object-attrs-uniformity [auto]: pin the object-aggregation attribute-uniformity rule for the four unasserted attributes — retention, expiry, explicit_timestamp, database_retention_policy/ttl — via ParseInterface rejections ("must all share the same <attr>", pkg/interfaceschema/parse.go:626-647) each paired with an at-uniformity acceptance twin; only reliability is pinned today (testdata/invalid/17-object-mixed-attrs.json, sameObjectAttributes at 54.5% stmt coverage).
- [x] interfaceschema-malformed-placeholder-fixtures [auto]: add invalid fixtures + manifest entries for parseSegment's four untested param branches (pkg/interfaceschema/parse.go:509-517) — "/%{", "/%ax", "/%{x", "/%{1}" — each rejected with the malformed-placeholder / invalid-placeholder-name message; the only placeholder-syntax fixture today is the mid-segment 12-endpoint-partial-placeholder.json (ContainsAny branch), so a regression accepting a dangling "/%{" passes every test.
- [x] interfaceschema-enum-roundtrip-coverage [auto]: extend TestEnumRoundTrips (pkg/interfaceschema/parse_test.go:244) to cover Ownership, Aggregation, Retention, DatabaseRetentionPolicy — the four enums it omits (Retention.String/MarshalJSON and DatabaseRetentionPolicy.String at 0% stmt coverage in the package's own tests) — so a typo'd wire string like "retention":"descard" or "aggregation":"obilver" fails the suite instead of passing.
- [x] interfaceschema-trie-interior-literal-param [auto]: pin the rule that an interior-only literal trie child does not shadow its parametric sibling at the matching depth (pkg/interfaceschema/trie.go:124-151) — purpose-build a trie with "/a/b/c" and "/a/%{p}", assert Match("/a/b") returns the /a/%{p} mapping and Match("/a/b/c") the /a/b/c one; statement coverage cannot see this branch (matchNode reports 100%), so a "fix" that misses on shallow paths through an interior literal node sails through the suite.
- [x] docs-sync-ae-write-405 [auto]: add the missing `405` response to the nine data-write ops in docs/api/astarte_appengine_api.yaml — PUT/POST/DELETE `/devices/{device}/interfaces/{interface}/{path}`, the three `/devices-by-alias/{alias}` twins, and the three `/groups/{group}/devices/{device}` twins. `writeError` emits `405 Method Not Allowed` with "Cannot write to device owned resource" (`ErrNotServerOwned`) and "Cannot write to read-only resource" (`ErrNotAProperty`) at internal/appengine/http.go:673-676, reachable from putData/putDataByAlias/putDataInGroup (http.go:292-357) and deleteData/deleteDataByAlias/deleteDataInGroup (http.go:311-377); `"405"` currently appears nowhere in the spec (those ops document only 200/204, 400, 401, 404, 500). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-write-value-422 [auto]: add the missing `422` response to the six PUT/POST data-write ops in docs/api/astarte_appengine_api.yaml (device-scoped, by-alias, and in-group at `/interfaces/{interface}/{path}`). A value larger than `maxValueBytes` (64 KiB, internal/appengine/http.go:24) is answered `422` "Value size exceeds size limits" by putData/putDataByAlias/putDataInGroup (http.go:298-301, 327-329, 346-348), but the specs for PUT/POST (yaml:686-695, 726-736, 453-463, 499-509, 1176-1186, 1218-1228) list only 200/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-post-groups-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups` in docs/api/astarte_appengine_api.yaml — `Service.CreateGroup` returns `ErrGroupAlreadyExists` (internal/appengine/service.go:495, mapped at http.go:660-661) which the handler answers as `409` "Group already exists"; the spec (yaml:821-839) documents only 201/400/401/422/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] payload-longinteger-fraction-quantize [auto]: `jsonInt64` (pkg/payload/json.go:255-271) accepts fractional longinteger tokens that the ParseFloat fallback silently rounds onto an integer inside its own documented-safety window |x| < 2^53 (json.go:250-254 claims values there never round, "silently corrupting the value") — `{"v":4503599627370496.5}` decodes to 4503599627370496, `{"v":6000000000000000.9}` to 6000000000000001, and `{"v":5e-400}` underflows to 0, all without error, while the contract says fraction forms are "accepted only when integral". Verify the token is an exact decimal integer (`math/big.Rat` SetString + IsInt) on the float fallback — keeps `5.0`, `1e3`, `0.5e1` accepted, rejects the rounding cases as ReasonTypeMismatch; plain-integer ParseInt fast path and the > 2^53 window check unchanged (BSON unaffected: longinteger accepts int64/int32 only). Add failing rows to TestJSONProfile.
- [x] payload-json-malformed-t-not-tolerated [auto]: pin in pkg/payload/payload_test.go that a malformed `t` is rejected with ReasonBadTimestamp even on a non-explicit mapping — `decodeJSONEnvelope` (json.go:50-66) parses `t` before the facade learns the explicit-timestamp policy (payload.go:52-54), so only a *valid* `t` is "tolerated and ignored" while `{"v":1,"t":"yesterday"}` with explicit=false is still a rejection; only the valid-t-ignored row and the BSON "t wrong type" twin are asserted today, so a regression implementing the literal leniency reading (ignore a bad t on non-explicit mappings) passes every suite. Add one table row to TestJSONProfile.
- [x] docs-sync-ae-post-group-devices-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml — re-adding a member the group already holds answers `409` "Device already in group" (`Service.AddGroupDevice` returns `ErrDeviceAlreadyInGroup`, internal/appengine/service.go:636-638, mapped at http.go:662-663), but the spec (yaml:967-977) lists only 201/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-patch-merge-patch-media-type [auto]: fix the three PATCH request bodies in docs/api/astarte_appengine_api.yaml — PATCH /devices/{device} (yaml:176), PATCH /devices-by-alias/{alias} (yaml:260), PATCH /groups/{group}/devices/{device} (yaml:1023) declare `application/json; charset=utf-8`, but all three handlers enforce `Content-Type: application/merge-patch+json` exactly and reject every other value — 500 on the device and group paths (internal/appengine/http.go:209-211, 521-523), 400 on the alias path (http.go:230-232). Change the requestBody media type to `application/merge-patch+json` on all three. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
- [x] config-fail-loud-dev-mode [auto]: `applyEnv` silently discards a malformed `ASTRATE_MQTT_INSECURE_DEV_MODE` value (internal/config/config.go:255-257) — `strconv.ParseBool` error is discarded, the field becomes false, and TLS requirements are silently re-enabled (or a misleading "tls_cert_file required" error surfaces via validate). Return the parse error on bad input; add a failing test that asserts `ASTRATE_MQTT_INSECURE_DEV_MODE=yes` causes Load to error.
- [x] config-key-resolvers-test [auto]: add container-free unit tests for `HousekeepingKeys()` and `RealmJWTPublicKey()` (internal/config/config.go:383-411) — both file-reading helpers have zero test coverage. Cover: happy path (inline PEM + file PEM → concatenated/trimmed), file-not-found error (assert error message includes the path), inline-over-file priority (`RealmJWTPublicKey` returns the inline block when both are set), and TrimSpace behaviour (leading/trailing newlines stripped). `t.TempDir()` + `os.WriteFile` for fixtures; no Docker.
- [x] compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.
- [x] docs-sync-hk-delete-gating-responses [auto]: add the missing `405` and `422` responses to `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — `writeError` answers 405 "Realm deletion disabled" for `ErrDeletionDisabled` and 422 `{"errors":{"error_name":["connected_devices_present"]}}` for `ErrConnectedDevicesPresent` (internal/housekeeping/http.go:233-249), both reachable from `Service.DeleteRealm` (service.go:248-266), but the spec lists only 204/401/403/404/500 (yaml:198-208). Its sole `405` reference sits on PATCH (yaml:182-183), which can never return it (`UpdateRealm`, service.go:196-216, has no deletion-gating path) — re-point 405 to DELETE (reusing the MethodNotAllowed component) and add 422 with the FieldErrors shape and a connected_devices_present example. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_housekeeping_api.yaml (yaml:417) — it shows `detail: "housekeeping: validation failed: realm name must not be empty"`, but `validationDetail` strips the `housekeeping: validation failed: ` prefix on the wire (internal/housekeeping/http.go:251-258) and the real blank-name message is `"realm_name can't be blank"` (service.go:138). Change the example to `detail: "realm_name can't be blank"`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

- [x] auth-iat-not-required-test [auto]: pin in internal/auth/jwt_test.go that a token with a *future* `iat` still verifies — the "`iat` is not required — upstream parity" rule (jwt.go:104-105, no `jwt.WithIssuedAt` in the parser, jwt.go:115-118) has zero coverage, so a regression that starts validating `iat` or `WithIssuedAt` would silently reject tokens from a skewed-clock issuer and pass every existing suite. Sign RS256 with `iat` = now+1h and no `exp`/`nbf`, assert `Verify` succeeds.
- [x] auth-empty-claims-403-test [auto]: pin in internal/auth/middleware_test.go that a validly signed token carrying NO Astarte claims is `403` on a guarded route, not `401` — the authenticated-but-grant-less distinction (middleware.go:149-151, `Authorizes` on a nil grants map → false → `WriteForbidden`) is currently untested: every 403 fixture has a claim on the wrong surface, so a regression mapping grant-less tokens to 401 (or skipping the forbidden path) passes the suite. Sign RS256 with empty `jwt.MapClaims`, hit the pairing route, expect 403 + golden `envelope_403.json`.
- [x] auth-cache-default-size-test [auto]: pin in internal/auth/cache_test.go that `NewCache(0)` and `NewCache(-1)` fall back to `DefaultCacheSize` (cache.go:29-31) — the `size < 1` branch exists purely so `lru.New` cannot panic and has zero coverage; a regression dropping the fallback (panic on `lru.New` with a non-tick size, or a silently tiny cache) would pass every test. Construct with both sizes, assert no panic and a working `Verify` with `lru.Len() == 0`.
- [x] hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto] — already resolved and re-queued by the refill: commit 015a706 chose the upstream-parity behaviour (0 ≡ unset, not invalid, not a literal 0) and pinned it — create folds 0 → nil with the container-free TestCreateRealmZeroRetentionUnsets (internal/housekeeping/service_test.go:186) and PATCH 0 → ClearRetention in the integration suite (internal/housekeeping/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).
- [x] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication is "via Authorization header or `?token=` query parameter" (yaml:293-294), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] astarteapi-metadata-envelope-test [auto]: add a container-free golden test for `WriteDataWithMetadata` (pkg/astarteapi/envelope.go:154-158) pinning the two wire branches — nil metadata omits the `metadata` key (`{"data": v}` only, via the `omitempty` at envelope.go:151), populated metadata emits `{"data": v, "metadata": {...}}` — with 200 + the package `Content-Type`, matching the TestGoldenEnvelopes style. Today the non-nil branch is only covered by the integration-tagged formats_parity_test.go (runs on the Legion Go, not the Pi gate) and the nil-omission branch is covered nowhere (internal/appengine/data.go:262-268 always sets metadata), so a regression dropping the omitempty passes every suite on this box.
- [x] astarteapi-multikey-sorted-golden [auto]: extend TestGoldenEnvelopes (pkg/astarteapi/envelope_test.go:101-121) with multi-key `WriteFieldErrors` and `WriteRawErrors` maps asserting the sorted key order on the wire — the docs claim "keys are emitted in sorted order" (envelope.go:170, 184) but every fixture uses a single key, so a change to the marshalling path that stopped sorting (e.g. switching encoders) passes the suite.
- [x] docs-sync-pairing-deviceendpoints-dead-404-403 [auto]: remove the unreachable `404` DeviceNotFound response from the three device-facing ops in docs/api/astarte_pairing_api.yaml — `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials` (yaml:226-227), `GET /pairing/v1/{realm}/devices/{deviceID}` (yaml:268-269), `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify` (yaml:339-340) — and the `403` Forbidden from the info GET (yaml:266-267). `Service.authenticateDevice` folds every failure (malformed/unknown deviceID, unknown realm, unregistered device, wrong secret) into uniform `ErrUnauthorized` (internal/pairing/service.go:380-411), so `store.ErrNotFound` never reaches `writeServiceError` (http.go:326-327) from those handlers; and `Service.Info` never returns `ErrInhibited` (inhibited devices still read info with 200 + status "inhibited", service.go:286-307, comment at 282-285). Reachable on these three are only 200/201/400/401/422/429/500 plus 403 on credentials and verify (both genuinely raise `ErrInhibited`, service.go:252-254 and 329-331). Keep 404 on the agent endpoints (register yaml:149-150, unregister yaml:176-177) where `GetRealmByName`/`deviceid.Parse` really do return `store.ErrNotFound`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-version-404-unreachable [auto]: remove the `404` "Unknown realm" response from `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml (yaml:96-104) — the route is served by `observability.VersionHandler` (cmd/astrate/main.go:398, internal/observability/compat.go:38-47), which answers static `200 {"data":"<version>"}` with no realm lookup, so a 404 can never be emitted there. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] observability-readiness-wedge-test [auto]: make the readiness budget injectable and pinned — add a `timeout time.Duration` field on `Health` (internal/observability/health.go) defaulting to `readinessTimeout` (health.go:12, no test currently exercises the wedged-dependency path), and add a test in observability_test.go asserting a ctx-honoring wedged check makes `/astrate/v1/readiness` return 503 reporting it failing within a short injected budget, instead of hanging the endpoint.
- [x] observability-dbpool-gauge-coverage [auto]: extend `TestMetricsExposesGauges` (internal/observability/observability_test.go:25-34) to assert all four `astrate_db_pool_*` gauges — it currently pins only `acquired_conns` and `max_conns`, so a regression deleting or mis-wiring `idle_conns` or `total_conns` (metrics.go:77-78) passes the suite; expect idle 5, total 7 from the supplied `DBPoolStats{2,5,7,10}`.
- [x] observability-health-content-type-test [auto]: assert `Content-Type: application/json` on `/astrate/v1/health` and `/astrate/v1/readiness` (set by `writeJSON`, health.go:81) and on both the 200 and 503 paths of the per-service `/health` (set by hand in compat.go:26,31) — only the version endpoint pins the header today (compat_version_test.go:31-33), so a regression dropping the header on the /astrate/v1 surface passes the suite.

- [x] docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] httpx-cors-vary-origin-passthrough [auto]: in internal/httpx/cors.go non-wildcard mode, the disallowed-origin pass-through (cors.go:26-29) returns without `Vary: Origin`, so responses whose CORS outcome depends on the Origin header reach shared caches without the RFC 9110 §12.5.5 variation marker (a cache could serve the CORS-approving variant to a disallowed origin, or the headerless variant to an allowed one); the stamp lives only in the allowed branch (cors.go:35-36) and only the preflight pins it (cors_test.go:42). Add `h.Add("Vary", "Origin")` on the pass-through branch too, and assert `Vary: Origin` in TestCORSActualRequestStamped (actual non-preflight, pinned only for Allow-Origin today) and TestCORSDisallowedOriginPassesThroughUnstamped. Additive header change; container-free.
- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] channels-group-watch-membership-scope: group watches authorize against `WATCH::groups/<name>/…` (ws.go:178-205, conformance channels.json) but `Room.dispatch` ever only filters by `watchEntry.deviceID` (empty for a group watch — room.go:36-40, 196-200) and the trigger compiles with no group/device scope (triggers/match.go:374-396, because upstream carries group_name at the payload top level), so a group-scoped token receives every realm event matching interface/path/on, members or not. Plumb a group-membership resolver into internal/appengine/channels (it knows only `Bus` today), resolve membership at Watch time and store it in `watchEntry`, filter deliveries by it. Add a container-free test with a fake resolver: watch group `probe`, publish an event from a member and a non-member, assert only the member's arrives. Verify the member-filter rule against upstream on the Legion if disputed; the boundary leak stands regardless.
- [x] channels-rejoin-joinref-tagging-test: pin in internal/appengine/channels/ws_test.go that a second `phx_join` on the same topic updates the session's `joinRef` (ws.go:284-291) and subsequent `new_event` pushes carry the fresh ref (read per event at ws.go:403) — no test drives a rejoin today (all `TestWireSession_*` join once), so the comment's "server pushes must be tagged with the current one" rots unasserted. Wire: join with ref "1", rejoin with "2", push an event, assert element 0 is `"2"`.
- [!] channels-rejoin-authz-mismatch: `handleJoin`'s comment says upstream authorizes every join including a rejoin (ws.go:276-277), but the rejoin branch (ws.go:284-291) replies OK without calling `s.tok.AuthorizesChannel(auth.VerbJoin, name)` — behaviourally identical today (a token is immutable for the session) but the comment lies about what the code does. Either move the authz check to cover the rejoin branch too, or fix the comment; mechanical, no behaviour change, choose one. — BLOCKED: wrote nothing
- [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a locally-run `elixir -e` with `Base.url_decode64!`) with the device ID `"AAAAAAAAAAAAAAAAAAAAAB"` — nonzero unused bits in the 22nd base64url char — and report whether upstream accepts it: `pkg/deviceid/deviceid.go:34-37` uses `base64.RawURLEncoding.Strict()` and deviceid_test.go:176 asserts rejection while the comment claims parity with `Elixir Base.url_decode64!(padding: false)`, which decodes and discards those bits rather than erroring. Astrate may therefore be stricter than upstream on the same wire form (a device registered with such an id on upstream would be rejected here). If upstream rejects it identically, the strict encoding is verified parity and the task is done; if upstream accepts it, escalate the relax-vs-strict call to `.mule/for-giulio.md` with the measurement. Probe first, no code change either way.
- [x] docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate_devices_total{realm="test"} 42`, but no such series is ever registered — the registry carries `astrate_broker_sessions` (internal/observability/metrics.go:54-62), `astrate_db_pool_{acquired,idle,total,max}_conns` (metrics.go:65-79), plus the engine/flow/trigger families (`astrate_engine_*` engine/router.go:411-444, `astrate_engine_stream_dropped_total` engine/stream/bus.go:119, `astrate_engine_trigger_*` engine/triggers/actions.go:296-300, `astrate_flow_router_*` flow/router.go:260-276) and the standard `go_*`/`process_*` runtime collectors — so the example advertises a scrape result the wire can never return. Swap it for a real series (e.g. `astrate_broker_sessions 0`), keeping the HELP/TYPE comment style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] store-devices-inhibit-re-register [legion] [auto]: `RegisterDevice` (internal/store/devices.go:75-91) silently clears the inhibit flag — its `ON CONFLICT ... DO UPDATE SET status = 'registered'` fires for any device with `first_credentials_request IS NULL`, including one the admin inhibited via `SetDeviceInhibited` (devices.go:251-268), which sets `status='inhibited'` on unconfirmed devices too; §5.3 says an inhibited device blocks new credentials and connections, so the re-registration re-opens it. Preserve `'inhibited'` in the SET (e.g. `status = CASE WHEN devices.status = 'inhibited' THEN 'inhibited' ELSE 'registered' END`) and add a Lifecycle case in internal/store/devices_test.go: inhibit an unconfirmed device, re-register, assert status stays inhibited and the secret still rotates. Verify the assert against upstream's register-not-touching-inhibit on the Legion while the integration suite runs.
- [x] store-pipelines-empty-name-zero-blocks-test [auto]: extend internal/store/pipelines_validate_test.go with the two untested error branches of `validatePipelineGraph` (internal/store/pipelines.go:51-64) — a definition with zero blocks and one whose block has an empty name must both be rejected. Pure helper, no DB; the existing suite covers only acyclic/cyclic/unknown-ref/duplicate-name.
- [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (internal/store/devices.go:130-131, "if several devices share an alias the lowest device ID wins", `ORDER BY id LIMIT 1`) — no test drives it today; add a case in internal/store/devices_test.go that gives two devices the same alias and asserts the lookup resolves the lower ID. Needs the integration DB.
- [x] broker-acl-coldstart-fallback-flood [auto]: `syncOwnershipOf` (internal/broker/authhook.go:212-237, reached from the ACL miss path aclhook.go:196-206) does a full synchronous `GetDevice`+`GetInterface` store read for every distinct interface name a device publishes to, with no rate limit — `introspectionReloadDebounce` (authhook.go:42-45, "so an adversarial topic flood cannot hammer the database") throttles only `refreshIfStale`'s full reload, and the fallback's per-name cache write means N distinct bogus names in one window cause N synchronous full-device reads per second, an open anti-flood hole the pre-fix denied-miss path did not have. Gate the fallback to one sync read per session per debounce window (or fold it into a full reload that stamps `lastIntroLoad`), keep `TestBrokerACLColdStartIntrospectionMiss` (one name) green, and add a T1 test in broker_test.go with a GetDevice-counting fake: connect cold, publish to K>=2 distinct unknown-interface names, assert no more than one fallback read was made.
- [x] broker-offlineacl-entry-eviction [auto]: `offlineACL.entries` (internal/broker/aclhook.go:100-162) is append-only — `ownershipOf` upserts one entry per CN ever ACL-checked offline and nothing ever deletes, so a long-lived broker grows the map unboundedly across device churn even after the 10s `offlineACLCacheTTL` passes (the entry is kept and re-stamped, not reaped). Evict stale entries lazily on access (drop an entry whose `loadedAt` predates some multiple of the TTL before inserting another, or cap+LRU), add an injectable `now func() time.Time` in the same style as `lifecycleHook.now`, and cover in a container-free unit test: seed N entries, advance the clock, access one, assert the map stays bounded.
- [!] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine_api.yaml — `retrieve_metadata` and `downsample_key` are parsed for every data read by `parseQueryOpts` (internal/appengine/http.go:624 `retrieve_metadata`, 633 `downsample_key`) yet appear on none of the read endpoints: device-scoped GET `{interface}` (yaml:583-606) and `{interface}/{path}` (yaml:623-649), by-alias (yaml:324-368, 369-417), and in-group (yaml:1108-1132, 1149-1176) — the parameter block on each lists only since/since_after/to/limit/downsample_to/sort/format/allow_bigintegers/allow_safe_bigintegers (components DataSince..DataAllowSafeBigIntegers, yaml:1360-1439). Add the two `$ref`s (new DataRetrieveMetadata and DataDownsampleKey components, or inline — mirror the existing style) to all six, and update the DataSort-independent `sort` default note only if verified. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] examples-echo-container-contract-test [auto]: add a container-free httptest suite for examples/flow-container-echo — the only no-test package in the repo that contains logic (`handleMessage`, main.go:49-75; `sanitizeLog`, main.go:40-47) — pinning (1) POST `/v1/message` echoes a valid JSON-object body verbatim with 200, answers 400 on invalid JSON and 204 on an empty body (the 1 MiB read cap), and (2) metadata `echo_drop` == "1" forces the 204 filter/drop path; the example is the canonical reference container contract (README.md, docs/handoff/flow-design-b-container-block-2026-07-29.md), so a regression silently misleads every container author who copies it.
- [!] flow-mqtt-source-reconnect-recovery [auto]: `mqttSource` latches `lost` on paho's ConnectionLost handler (internal/flow/blocks/mqtt.go:185-187, gates in Emit/Process at 234-236 and 252-254) but nothing ever clears it, even though `newMQTTClient` enables `SetAutoReconnect(true)` (mqtt.go:54) and the docstring promises background auto-reconnect — so any transient broker blip permanently kills the source for the flow's lifetime ("connection lost" forever) while the flow stays `running`; `TestMQTTSource_ConnectionLost` (mqtt_test.go:190) only asserts the error surfaces, never recovery. Clear the latch on reconnect (a paho OnConnect handler, or gate on `client.IsConnected()` instead of a sticky flag), and add a container-free test that stops and restarts the embedded mochi broker and asserts Emit resumes. — BLOCKED: wrote nothing
- [x] flow-msg-json-integer-precision [auto]: `setDataFromWire`'s TypeInteger float64 branch does `int64(v)` on a value that plain `json.Unmarshal` (message.go:141-145) already rounded to float64, so integer wire values > 2^53 silently lose precision and wrap (`float64(MaxInt64)` → 2^63 → MinInt64 on decode), a fractional `3.7` silently truncates to 3, and `1e300` becomes garbage — the `json.Number` branch (message.go:268-273) that would reject these is unreachable in this path. Mirror the accepted `payload-longinteger-fraction-quantize` rule: decode with `json.Decoder`+`UseNumber` (parse via `json.Number.Int64`) or reject non-integral/out-of-int64 float64; add round-trip rows for `math.MaxInt64`, `1e300`, `3.7` to message_test.go. Container-free.
- [x] flow-sort-bounded-buffer [auto]: `Sort` (internal/flow/blocks/sort.go:32-65) appends every message to an uncapped `buf` and only flushes while `len(buf) > 1 && buf[0].ts <= newest-windowUs` (sort.go:60) — a stream whose timestamps stay within `window_ms` of the newest (clustered bursts, or messages sharing one timestamp) never satisfies the condition, so the block emits nothing while the slice grows without bound, and the documented "newest is never emitted" mean the final message is held indefinitely; secondarily `window_ms * 1000` (sort.go:31) overflows for large configs and silently inverts the flush edge. Add a configurable max buffered count with a documented overflow policy (drop-oldest / error / force-flush) and a cap test feeding N same-timestamp messages; guard the multiply. Container-free.
- [!] flow-randomsource-span-overflow [auto]: `randomSource.next()` computes `rand.Int64N(s.maxInt-s.minInt+1)` (internal/flow/blocks/randomsource.go:171) whose argument overflows int64 to a non-positive value when min/max span more than the int64 range — config like `min:-9e18, max:9e18` passes the single `min <= max` check (randomsource.go:106-108) but makes `rand.Int64N` panic (it panics on n <= 0) inside the source pump goroutine (flow.go:205-239, no recover), crashing the whole process. Reject a span wider than `math.MaxInt64` at construction alongside the existing check, and add a construct/Emit test with the wide span. Container-free. — BLOCKED: lint failed: internal/flow/blocks/randomsource.go:113:40: G115: integer overflow conversion int64 -> uint64 (gosec)
- [x] flow-filter-key-contains-test [auto]: `Filter`'s `key_contains` rule (`strings.Contains`, internal/flow/blocks/transform.go:69) has zero test coverage — `transform_test.go` only exercises `key_prefix`, `type`, unknown-type and the at-least-one construct rule, so a regression silently dropping the Contains check would pass the whole suite. Add substring-match, substring-absent, and prefix+contains combined rows to the filter tests. Container-free.

- [x] docs-sync-rm-put-interface-409 [auto]: add the missing `409` Conflict response to `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` in docs/api/astarte_realm_management_api.yaml (yaml:222-234 lists only 204/400/401/404/422/500) — `updateInterface` answers 409 for every url/body disagreement and minor-upgrade incompatibility: `ErrNameMismatch` (service.go:204 → "Interface name doesn't match the one in the interface json", http.go:367-369), `ErrMajorMismatch` (service.go:207 → http.go:370-372), and the `CheckMinorUpgrade` sentinels `ErrMinorNotIncreased`/`ErrDowngradeNotAllowed`/`ErrMissingEndpoints`/`ErrIncompatibleEndpointChange` (service.go:223-226 → http.go:373-384) plus `store.ErrAlreadyExists` (http.go:385-386) — the exact set `TestRealmManagementErrorCodes` pins as StatusConflict (http_test.go:417-437). Reuse the `Conflict` response component (already used on POST, yaml:131-132) and add name-mismatch/major-mismatch examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` (yaml:231-232) in docs/api/astarte_realm_management_api.yaml as a oneOf — both handlers go through `writeInterfaceError`, which answers three distinct 422 bodies: the flat ErrorDetail (`ErrValidation` via `validationDetail`, http.go:392-393), the nested violations changeset envelope (`writeViolations`, http.go:358-359 + 434-533, e.g. `{"errors":{"description":["should be at most 1000 character(s)"]}}` and the aligned full-length `mappings` array, http_test.go:480-501), and the named FieldErrors envelope for `ErrMaximumDatabaseRetentionExceeded` (`{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, http.go:355-357, http_test.go:282-320); today only the flat ValidationError is referenced. Mirror the createTrigger oneOf pattern (yaml:344-368). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "realm: validation failed: interface definition is invalid"`, but `validationDetail` strips the `realm: validation failed: ` prefix on the wire (http.go:562-568) and the parser's real messages are `invalid interface: ...` (pkg/interfaceschema/parse.go:38, wrapped at service.go:163); the example should show the stripped, real message. Same class as the already-fixed docs-sync-hk-validation-example (housekeeping http.go:251-258). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] swagger-sub-failfast [auto]: make `Mount` fail fast (panic) when `fs.Sub` over the `docs.SwaggerUI`/`docs.APIYAML` embeds errors — `uiRoot, _ := fs.Sub(...)` and `apiRoot, _ := fs.Sub(...)` (internal/swagger/swagger.go:17-18) silently serve an empty tree → 404s at /swagger/ and /api/ with no log the day the docs embed layout changes; add an fs-injecting internal variant (`MountWithFS(mux, uiRoot, apiRoot fs.FS)` called by `Mount` after the subs) so a broken-FS unit test in swagger_test.go can prove the fail-fast, keeping the public `Mount(mux)` signature and its single main.go:414 caller untouched. Container-free.
- [ ] flow-boot-resumes-stopped-durable-flows [legion] [auto]: RehydrateAutoRestart (cmd/astrate/main.go:216) restarts every flow with auto_restart=true whatever its last status — `ListAutoRestartFlows` filters on `f.auto_restart = true` alone (internal/store/flows.go:125) and nothing ever clears the flag, which is written once at create (flowapi/http.go:141-146 -> service.go:318) and never on stop — so a durable flow stopped through the API is running again after the next boot, while the shutdown mark that says otherwise (MarkRunningFlowsStopped, main.go:497 -> internal/flowapi/service.go:885-901, which writes status="stopped") is never read by that query. The two halves disagree: service.go:325 documents "every durable flow", the shutdown call assumes "only what was running". Say which is the intent, then make them agree — either add `AND f.status = 'running'` to flows.go:125 (making the shutdown mark load-bearing) or delete the mark as dead code — and add a case in internal/store/flows_test.go pinning the chosen rule (a stopped auto_restart row must, or must not, come back). Needs the DB.
- [x] drain-per-stage-budget [auto]: `shutdown` spends a single 30s sctx (cmd/astrate/main.go:479) across three sequential stages — srv.Shutdown (484), b.Close (489), flowSvc.Manager().Shutdown (494) and MarkRunningFlowsStopped (497) — while `drainEngine` alone opens a second one (504), so the constant's own comment "bounds the whole graceful drain" (main.go:55-56) is wrong in both directions, and one slow in-flight HTTP request (the server sets no per-request deadline, 221) burns the whole first budget and leaves the flow stage on an expired context: StopFlow then returns at `f.router.Drain(ctx)` before step 3 releases block resources (internal/flow/flow.go:265-273) and MarkRunningFlowsStopped skips every flow (flowapi/service.go:895-897), both surfacing as one log.Warn. Give each stage its own budget through an extracted helper (shutdownTimeout has to become an injectable var for the test) and pin it container-free: a second stage still receives a live context after a first stage exhausts its own.
- [x] cmd-devcert-fields-test [auto]: add a container-free test for `selfSignedDevCert` (cmd/astrate/main.go:320-345), the throwaway mTLS identity `insecure_dev_mode` boots on, which no test touches — pin NotBefore = now-1h (331; a plain now would reject a device whose clock is behind), NotAfter = now+365d (332), DNSNames localhost and IPAddresses 127.0.0.1/::1 (336-337, dropping which breaks hostname verification for a device dialling mqtts://127.0.0.1:8883), ExtKeyUsage serverAuth (334), and that the leaf is self-signed and parses (x509.ParseCertificate on Certificate[0]).
- [!] cmd-healthcheck-contract-test [auto]: add a container-free test for `runHealthcheck` (cmd/astrate/main.go:531-558) — the `HEALTHCHECK CMD ["/astrate", "-healthcheck"]` contract (Dockerfile:35-36) has no coverage. Assert it probes /astrate/v1/readiness and not /health (543), returns 0 only on 200 so a 503 readiness exits 1 (554), and maps a bare ":PORT" in ASTRATE_HTTP_ADDR to 127.0.0.1 (536-538). httptest.NewServer plus t.Setenv("ASTRATE_HTTP_ADDR", srv.Listener.Addr().String()). — BLOCKED: tests failed: --- FAIL: TestRunHealthcheckProbesReadiness (0.01s)
- [x] cmd-loadsealer-masterkeyfile-test [auto]: add a container-free test for `loadSealer` (cmd/astrate/main.go:515-526) — the `security.master_key_file` to env hand-off (`os.Setenv(store.EnvMasterKeyFile, ...)`, 517) is the only wiring that makes a configured key file reach the store loader, and dropping or misspelling that branch leaves the file silently ignored and boot dying with the generic "master key" error (523) even though the config named a valid file. Assert a 32-byte key file under t.TempDir() loads, the inline ASTRATE_MASTER_KEY path still wins, and a missing key returns the error naming all three env vars. t.Setenv both ASTRATE_MASTER_KEY and ASTRATE_MASTER_KEY_FILE first so the os.Setenv is restored at cleanup.
- [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cd /root/astrate-mule && wc -l .mule/todo.md
229 .mule/todo.md

→ Read .mule/todo.md [offset=227, limit=3]
✗ Edit .mule/todo.md failed
Error: Could not find oldString in the file. It must match exactly, including whitespace, indentation, and line endings.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -225,5 +225,6 @@
 - [x] cmd-loadsealer-masterkeyfile-test [auto]: add a container-free test for `loadSealer` (cmd/astrate/main.go:515-526) — the `security.master_key_file` to env hand-off (`os.Setenv(store.EnvMasterKeyFile, ...)`, 517) is the only wiring that makes a configured key file reach the store loader, and dropping or misspelling that branch leaves the file silently ignored and boot dying with the generic "master key" error (523) even though the config named a valid file. Assert a 32-byte key file under t.TempDir() loads, the inline ASTRATE_MASTER_KEY path still wins, and a missing key returns the error naming all three env vars. t.Setenv both ASTRATE_MASTER_KEY and ASTRATE_MASTER_KEY_FILE first so the os.Setenv is restored at cleanup.
 - [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion Go — the device-delete leg added by docs-sync-rm-delete-device-async-operation-param (register a device, `DELETE /devices/{id}?async_operation=…` answers 204 and the row is gone on both values) compiles clean under `go vet -tags integration ./internal/realm/` but has never been executed, because the Pi has no database; until it runs, TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go is the only assertion of that leg that has actually been verified.
 - [ ] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cd /root/astrate-mule && git status --short && echo "=== final gate ===" && gofmt -l . && go vet ./... && go test -count=1 ./... 2>&1 | rg -v "^ok|no test files" ; echo "EXIT=$?"
M .mule/todo.md
 M docs/COMPATIBILITY.md
 M docs/api/astarte_realm_management_api.yaml
 M internal/realm/async_operation_test.go
 M internal/swagger/swagger_test.go
=== final gate ===
--- FAIL: TestChangeTriggersLifecycle (0.03s)
panic: runtime error: index out of range [-1] [recovered, repanicked]

goroutine 474 [running]:
testing.tRunner.func1.2({0x82dd80, 0x7f53e82708})
	/usr/lib/go-1.26/src/testing/testing.go:1974 +0x1a0
testing.tRunner.func1()
	/usr/lib/go-1.26/src/testing/testing.go:1977 +0x318
panic({0x82dd80?, 0x7f53e82708?})
	/usr/lib/go-1.26/src/runtime/panic.go:860 +0x12c
github.com/astrate-platform/astrate/internal/engine.lastOf(...)
	/root/astrate-mule/internal/engine/triggers_change_test.go:257
github.com/astrate-platform/astrate/internal/engine.TestChangeTriggersLifecycle(0x7f30013b08)
	/root/astrate-mule/internal/engine/triggers_change_test.go:136 +0x1774
testing.tRunner(0x7f30013b08, 0x8c8580)
	/usr/lib/go-1.26/src/testing/testing.go:2036 +0xc4
created by testing.(*T).Run in goroutine 1
	/usr/lib/go-1.26/src/testing/testing.go:2101 +0x3a8
FAIL	github.com/astrate-platform/astrate/internal/engine	2.877s
FAIL
EXIT=0

$ cd /root/astrate-mule && for i in 1 2 3 4 5; do go test -count=1 ./internal/engine/ -run TestChangeTriggersLifecycle 2>&1 | tail -1; done
ok  	github.com/astrate-platform/astrate/internal/engine	0.075s
ok  	github.com/astrate-platform/astrate/internal/engine	0.071s
ok  	github.com/astrate-platform/astrate/internal/engine	0.068s
ok  	github.com/astrate-platform/astrate/internal/engine	0.078s
ok  	github.com/astrate-platform/astrate/internal/engine	0.074s

$ cd /root/astrate-mule && go test -count=1 ./... 2>&1 | rg -v "^ok|no test files"; echo "---second run---"; go test -count=1 ./... 2>&1 | rg -v "^ok|no test files"; echo "---third run---"; go test -count=1 ./... 2>&1 | rg -v "^ok|no test files"
---second run---
---third run---
time=2026-09-25T22:15:53.488+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-25T22:15:53.489+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:38633
time=2026-09-25T22:15:53.489+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-25T22:15:53.500+02:00 level=INFO msg="mochi mqtt server started"
time=2026-09-25T22:15:53.510+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-25T22:15:53.516+02:00 level=WARN msg="" listener=test-tcp error="read tcp 127.0.0.1:38633->127.0.0.1:55396: use of closed network connection"
time=2026-09-25T22:15:53.516+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-25T22:15:53.517+02:00 level=INFO msg="mochi mqtt server stopped"
time=2026-09-25T22:15:53.517+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-25T22:15:53.517+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:40319
time=2026-09-25T22:15:53.518+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-25T22:15:53.522+02:00 level=INFO msg="mochi mqtt server started"
time=2026-09-25T22:15:53.543+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-25T22:15:53.543+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-25T22:15:53.543+02:00 level=INFO msg="mochi mqtt server stopped"
time=2026-09-25T22:15:53.544+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-25T22:15:53.544+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:45549
time=2026-09-25T22:15:53.545+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-25T22:15:53.552+02:00 level=INFO msg="mochi mqtt server started"
time=2026-09-25T22:15:53.556+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-25T22:15:53.560+02:00 level=WARN msg="" listener=test-tcp error="read tcp 127.0.0.1:45549->127.0.0.1:60016: use of closed network connection"
time=2026-09-25T22:15:53.560+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-25T22:15:53.560+02:00 level=INFO msg="mochi mqtt server stopped"
time=2026-09-25T22:15:53.814+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-25T22:15:53.815+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:46535
time=2026-09-25T22:15:53.816+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-25T22:15:53.820+02:00 level=INFO msg="mochi mqtt server started"
time=2026-09-25T22:16:03.833+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-25T22:16:03.835+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-25T22:16:03.835+02:00 level=INFO msg="mochi mqtt server stopped"
--- FAIL: TestMQTTSink_HappyPath (10.02s)
    mqtt_test.go:246: timed out waiting for published message
time=2026-09-25T22:16:03.837+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-25T22:16:03.838+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:34985
time=2026-09-25T22:16:03.839+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-25T22:16:03.840+02:00 level=INFO msg="mochi mqtt server started"
time=2026-09-25T22:16:03.857+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-25T22:16:03.858+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-25T22:16:03.859+02:00 level=INFO msg="mochi mqtt server stopped"
time=2026-09-25T22:16:03.860+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-25T22:16:03.860+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:45471
time=2026-09-25T22:16:03.861+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-25T22:16:03.866+02:00 level=INFO msg="mochi mqtt server started"
time=2026-09-25T22:16:03.881+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-25T22:16:03.882+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-25T22:16:03.882+02:00 level=INFO msg="mochi mqtt server stopped"
time=2026-09-25T22:16:03.888+02:00 level=INFO msg="added hook" hook=allow-all-auth
time=2026-09-25T22:16:03.890+02:00 level=INFO msg="attached listener" id=test-tcp protocol=tcp address=127.0.0.1:45663
time=2026-09-25T22:16:03.891+02:00 level=INFO msg="mochi mqtt starting" version=2.7.9
time=2026-09-25T22:16:03.903+02:00 level=INFO msg="gracefully stopping server"
time=2026-09-25T22:16:03.904+02:00 level=INFO msg="stopping hook" hook=allow-all-auth
time=2026-09-25T22:16:03.904+02:00 level=INFO msg="mochi mqtt server stopped"
time=2026-09-25T22:16:03.908+02:00 level=INFO msg="mochi mqtt server started"
FAIL
FAIL	github.com/astrate-platform/astrate/internal/flow/blocks	11.308s
FAIL

$ cd /root/astrate-mule && go test -count=1 ./... 2>&1 | rg -v "^ok|no test files" | rg "FAIL|panic|timed out" ; echo "--- run 5 done (no output above = green) ---"; go test -count=1 ./internal/swagger/ ./internal/realm/ 2>&1 | tail -3
--- run 5 done (no output above = green) ---
ok  	github.com/astrate-platform/astrate/internal/swagger	0.112s
ok  	github.com/astrate-platform/astrate/internal/realm	0.057s

```
Done: `deleteDevice` now `$ref`s the existing `AsyncOperation` component parameter in docs/api/astarte_realm_management_api.yaml (no new component), the shared parameter description and deviation 17 in docs/COMPATIBILITY.md now name the device-deletion leg, and both guards cover it — `deleteDevice` added to the op list in TestRealmManagementAsyncOperationParamDocumented (container-free, fails without the yaml change: verified by removing the $ref) plus a device-delete case in TestRealmManagementAsyncOperationParam.
Files: docs/api/astarte_realm_management_api.yaml, docs/COMPATIBILITY.md, internal/swagger/swagger_test.go, internal/realm/async_operation_test.go, .mule/todo.md
Verified: `gofmt -l .` -> clean; `go vet ./...` -> pass; `go test -count=1 ./...` -> pass in 3 of 5 runs, the other 2 hit pre-existing flakes in packages I did not touch (internal/engine TestChangeTriggersLifecycle panicking in its own `lastOf` helper, internal/flow/blocks TestMQTTSink_HappyPath 10s broker timeout); `go vet -tags integration ./internal/realm/` -> compiles; `make -C docs build` -> `sync` target ran, then failed with `mkdocs: No such file or directory` (no pip/mkdocs on this box, unrelated to the change), so instead confirmed with an offline yaml.v3 probe in /tmp: spec parses (12 paths, 20 operations) and `#/components/parameters/AsyncOperation` resolves on deleteDevice (in: query, boolean, default false)
Unsure: the new device-delete case in TestRealmManagementAsyncOperationParam compiles but has NEVER BEEN EXECUTED — it is `//go:build integration` and needs a live DB, which this box has none of; `ssh legion` timed out (machine off), so the behaviour half is unverified and the docs-half assertion is the only one that has actually run. I also did not probe upstream for the "device deletion is async" claim behind deviation 17 (needs the live upstream stack); I extended prose the spec already asserted at yaml:687-689, inventing nothing new. Minor: deviation 17's paragraph was re-wrapped rather than minimally patched, so the diff is wider than the one added phrase.
Follow-ups: `verify-rm-device-delete-async-leg [legion]` appended to .mule/todo.md — run the integration suite on the Legion Go so the device-delete leg is proven.
```
