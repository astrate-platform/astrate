slug: docs-sync-native-error-example-capitalisation
verdict: done
at:  7cc237f
ran: 2026-10-03T16:52:09Z on DietPi in 284s


> build · big-pickle

I'll read the required files first.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ ls internal/swagger/ docs/ && ls docs/api/
docs/:
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

internal/swagger/:
swagger.go
swagger_test.go
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

$ rg -n "ErrorDetailExamples" internal/swagger/swagger_test.go
460:// TestPairingErrorDetailExamples pins the error-detail examples in the pairing
470:func TestPairingErrorDetailExamples(t *testing.T) {
509:// TestRealmManagementErrorDetailExamples pins the error-detail examples in the
512:// from TestPairingErrorDetailExamples. The BadRequest and InternalServerError
521:func TestRealmManagementErrorDetailExamples(t *testing.T) {
561:// TestAppEngineErrorDetailExamples pins the error-detail examples in the app
564:// management specs get from TestPairingErrorDetailExamples and
565:// TestRealmManagementErrorDetailExamples. The BadRequest and InternalServerError
573:func TestAppEngineErrorDetailExamples(t *testing.T) {

$ rg -n "detail: Internal" docs/api/astrate_native_api.yaml; echo "---"; rg -n "detail:" docs/api/astrate_native_api.yaml
466:              detail: Internal Server Error
---
434:            detail:
446:              detail: Unauthorized
456:              detail: Forbidden
466:              detail: Internal Server Error

→ Read internal/swagger/swagger_test.go [offset=400, limit=260]

$ rg -n "Detail[A-Za-z]+ +=|DetailInternalServerError|DetailUnauthorized|DetailForbidden" pkg/astarteapi/envelope.go | head -30
37:	DetailBadRequest = "Bad request"
38:	// DetailUnauthorized is the canonical 401 detail.
39:	DetailUnauthorized = "Unauthorized"
40:	// DetailForbidden is the canonical 403 detail.
41:	DetailForbidden = "Forbidden"
43:	DetailNotFound = "Not Found"
46:	DetailDeviceNotFound = "Device not found"
51:	DetailRouteNotFound = "Not found"
54:	DetailPageNotFound = "Page not found"
55:	// DetailInternalServerError is the canonical 500 detail. Measured
58:	DetailInternalServerError = "Internal server error"
110:		body, _ = marshal(detailEnvelope{Errors: detailBody{Detail: DetailInternalServerError}})
197:	return WriteError(w, http.StatusUnauthorized, DetailUnauthorized)
202:	return WriteError(w, http.StatusForbidden, DetailForbidden)
217:	return WriteError(w, http.StatusInternalServerError, DetailInternalServerError)

→ Read docs/api/astrate_native_api.yaml [offset=415, limit=60]

$ rg -n "^func (componentBlock|containsLine|operationBlock)" -A 30 internal/swagger/swagger_test.go | head -80
909:func operationBlock(t *testing.T, lines []string, operationID string) []string {
910-	t.Helper()
911-	marker := "      operationId: " + operationID
912-	start := -1
913-	for i, l := range lines {
914-		if l == marker {
915-			start = i
916-			break
917-		}
918-	}
919-	if start < 0 {
920-		t.Fatalf("spec declares no operation %q", operationID)
921-	}
922-	endRe := regexp.MustCompile(`^      operationId: |^  \S|^components:`)
923-	for i := start + 1; i < len(lines); i++ {
924-		if endRe.MatchString(lines[i]) {
925-			return lines[start+1 : i]
926-		}
927-	}
928-	return lines[start+1:]
929-}
930-
931-// componentBlock returns the lines of the components entry starting with the
932-// given key line, up to the next entry at the same or shallower indentation.
933:func componentBlock(t *testing.T, lines []string, key string) []string {
934-	t.Helper()
935-	start := -1
936-	for i, l := range lines {
937-		if l == key {
938-			start = i
939-			break
940-		}
941-	}
942-	if start < 0 {
943-		t.Fatalf("spec declares no components entry %q", key)
944-	}
945-	for i := start + 1; i < len(lines); i++ {
946-		if strings.TrimSpace(lines[i]) == "" {
947-			continue
948-		}
949-		if !strings.HasPrefix(lines[i], "      ") {
950-			return lines[start+1 : i]
951-		}
952-	}
953-	return lines[start+1:]
954-}
955-
956:func containsLine(lines []string, want string) bool {
957-	for _, l := range lines {
958-		if l == want {
959-			return true
960-		}
961-	}
962-	return false
963-}
964-
965-func embeddedYAMLFilenames() ([]string, error) {
966-	var names []string
967-	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
968-		if err != nil {
969-			return err
970-		}
971-		if d.IsDir() {
972-			return nil
973-		}
974-		if strings.HasSuffix(path, ".yaml") {
975-			names = append(names, strings.TrimPrefix(path, "api/"))
976-		}
977-		return nil
978-	})
979-	sort.Strings(names)
980-	return names, err
981-}
982-
983-func get(t *testing.T, url string) string {
984-	t.Helper()
985-	resp, err := http.Get(url)
986-	if err != nil {

$ rg -n "detail: (Bad Request|Internal Server Error|Not Found|Device Not Found|Page Not Found)" docs/api/ ; echo "=== exit $?"
docs/api/astrate_native_api.yaml:466:              detail: Internal Server Error
docs/api/astarte_pairing_api.yaml:59:                  detail: Not Found
docs/api/astarte_housekeeping_api.yaml:362:            detail: Bad Request
docs/api/astarte_housekeeping_api.yaml:385:              detail: Bad Request
docs/api/astarte_housekeeping_api.yaml:415:              detail: Not Found
docs/api/astarte_housekeeping_api.yaml:477:              detail: Internal Server Error
docs/api/astarte_appengine_api.yaml:1762:              detail: Not Found
docs/api/astarte_realm_management_api.yaml:1347:              detail: Not Found
=== exit 0

$ rg -n "housekeeping" internal/swagger/swagger_test.go | head -20; echo "=== todo grep"; rg -rn "housekeeping" .mule/todo.md | head -20
127:// upstream-parity the housekeeping and pairing specs already document.
196:// TestHousekeepingAsyncOperationParamDocumented guards that the housekeeping
202:// TestHousekeepingAsyncOperationParam in internal/housekeeping.
204:	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
206:		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
309:// TestHousekeepingRetentionZeroFoldDocumented guards that the housekeeping spec
312:// `null || val == 0` to ClearRetention (internal/housekeeping/http.go) and
314:// (internal/housekeeping/service.go) — upstream parity measured on v1.2.0. A
320:// in internal/housekeeping.
322:	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
324:		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
=== todo grep
    ssh legion 'cd ~/astrate/bench && ./scripts/run-tier.sh small astrate -base-url ... -n-key ...'
- [ ] race-check-store: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/store/... ./internal/n/... ./migrations/...`. Report any failure to .mule/for-giulio.md with the full race report. Split out of the former single `race-check` line, which timed out running the whole tree at once. [legion] [readonly]
- [x] n-tests [auto]: add a container-free unit test file for `internal/n` (service.go, 262 lines, zero test coverage). Cover two rules: (1) `CreateRealm` validation — blank name, blank JWT key, negative regLimit, negative retention, and default retention injection when `WithDefaultDatastreamMaximumStorageRetention` is set and the caller passes nil; (2) `DeleteRealm` gating — `ErrDeletionDisabled` when `WithRealmDeletionDisabled(true)`, `ErrConnectedDevicesPresent` when `DeviceStats` reports connected > 0. Use httptest-style fakes for the store and sealer (no Docker).
- [x] docs-sync-hk-patch-endpoint [auto]: add the undocumented `PATCH /n/v1/realms/{realm}` route to docs/api/astarte_n_api.yaml — it exists in code (internal/n/http.go:40, handler patchRealm at http.go:163-229, returns 200 on success) but the spec jumps straight from GET to DELETE. Document the request body (`RealmPatch` — same three optional fields as `RealmCreate`: `jwt_public_key_pem`, `device_registration_limit`, `datastream_maximum_storage_retention`, all nullable), the success response, and the error responses: 400 (bad envelope), 401, 403, 404, 405 (deletion disabled), 422 (invalid update parameters for unknown fields, blank jwt_public_key_pem, negative limit/retention, connected devices present), 500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-retention-field [auto]: add the `datastream_maximum_storage_retention` field to the `Realm` and `RealmCreate` schemas in docs/api/astarte_n_api.yaml — the field is on the wire (realmBody struct, internal/n/http.go:50, returned by viewBody at http.go:99, accepted in create at http.go:68 and patch at http.go:216) but absent from both schemas (yaml:200-231). Type: integer, format int64, nullable, description matching the device_registration_limit style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /n/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-delete-gating-responses [auto]: add the missing `405` and `422` responses to `DELETE /n/v1/realms/{realm}` in docs/api/astarte_n_api.yaml — `writeError` answers 405 "Realm deletion disabled" for `ErrDeletionDisabled` and 422 `{"errors":{"error_name":["connected_devices_present"]}}` for `ErrConnectedDevicesPresent` (internal/n/http.go:233-249), both reachable from `Service.DeleteRealm` (service.go:248-266), but the spec lists only 204/401/403/404/500 (yaml:198-208). Its sole `405` reference sits on PATCH (yaml:182-183), which can never return it (`UpdateRealm`, service.go:196-216, has no deletion-gating path) — re-point 405 to DELETE (reusing the MethodNotAllowed component) and add 422 with the FieldErrors shape and a connected_devices_present example. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_n_api.yaml (yaml:417) — it shows `detail: "n: validation failed: realm name must not be empty"`, but `validationDetail` strips the `n: validation failed: ` prefix on the wire (internal/n/http.go:251-258) and the real blank-name message is `"realm_name can't be blank"` (service.go:138). Change the example to `detail: "realm_name can't be blank"`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] hk-zero-retention-asymmetry [auto]: `PATCH /n/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/n/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto] — already resolved and re-queued by the refill: commit 015a706 chose the upstream-parity behaviour (0 ≡ unset, not invalid, not a literal 0) and pinned it — create folds 0 → nil with the container-free TestCreateRealmZeroRetentionUnsets (internal/n/service_test.go:186) and PATCH 0 → ClearRetention in the integration suite (internal/n/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).
- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/n/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "realm: validation failed: interface definition is invalid"`, but `validationDetail` strips the `realm: validation failed: ` prefix on the wire (http.go:562-568) and the parser's real messages are `invalid interface: ...` (pkg/interfaceschema/parse.go:38, wrapped at service.go:163); the example should show the stripped, real message. Same class as the already-fixed docs-sync-hk-validation-example (n http.go:251-258). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the n and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /n/v1/realms` and `DELETE /n/v1/realms/{realm}` in docs/api/astarte_n_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/n/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/n/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_n_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/n/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/n/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the n spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:556) and `InternalServerError` shows `detail: Internal Server Error` (yaml:638), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37, golden fixture pkg/astarteapi/testdata/error_bad_request.json) and `"Internal server error"` (envelope.go:58, testdata/error_internal.json), both measured on upstream 1.2.0 (envelope.go:33-36, 55-57). The other five detail examples on this file (`Not Found`, `Device not found`, `Unauthorized`, `Forbidden`, `Too Many Requests`) are already correct — change only those two. The identical wrong pair is in the other four specs (appengine 1732/1792, n 362/385/477, realm_management 1265/1348, native 466); leave those for a run of their surface. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
- [x] fdo-rc6-scope-delta-for-giulio [auto]: upstream v1.4.0-rc.6 (2026-09-30, 27 commits, 154 files, almost all FDO) grows the onboarding surface issue #78 is scoped against, and none of it is in #78's body (written against upstream 1.3 / rc.5): new `POST /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}/to0` re-runs TO0 and returns the refreshed rendezvous expiry (`{"data":{"expiry":…}}`, ownership_voucher.ex:63-73; 409 `device_already_onboarded` via fallback_controller.ex:120-125); voucher deletion, which revokes the rendezvous registration first and *keeps* the voucher if that fails (`:rendezvous_revocation_failed` → 500) and refuses another realm's voucher through `ensure_voucher_in_realm` (ownership_voucher.ex:40-60 — now expressible because the voucher row carries its realm); `device_id` accepted on voucher upload; the device registered immediately on upload rather than mid-TO2 (`add_unconfirmed_credentials`, device.ex:63-85); single-entry x5chain certs accepted; invalid vouchers no longer stored; vouchers and to2-sessions moved from the per-realm keyspace to the global `astarte` one with `fdo_guid` added to the device row (n/realms/queries.ex:784, 1129-1133); and RM device deletion now cascading to the voucher (device_removal/core.ex). Write the delta into `.mule/for-giulio.md` as a scope update for #78 and say which of these a client could actually observe — do not edit or comment on the issue (GitHub state is read-only for the mule) and do not implement anything.
- [x] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, n and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_n_api.yaml — the `ErrorDetail` schema-level example (yaml:361-362) and the `BadRequest` response (yaml:385) show `detail: Bad Request`, and `InternalServerError` (yaml:477) shows `detail: Internal Server Error`, but the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58, via WriteBadRequest/WriteInternalServerError, pinned by pkg/astarteapi/testdata/error_bad_request.json and error_internal.json). While you are in `ErrorDetail`: its `example` sits at the schema level and shows a bare `{detail: ...}` where the envelope the wire sends is `{"errors":{"detail":...}}` — move it under `errors` (or drop it, as the realm-management and pairing `ErrorDetail` schemas carry no example) so a generated client is not shown the wrong shape. Add `TestHousekeepingErrorDetailExamples` beside the pairing and realm-management halves in internal/swagger/swagger_test.go. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it

→ Read docs/api/astarte_housekeeping_api.yaml [limit=40, offset=350]

$ rg -n "docs-sync-native-error-example-capitalisation|^# |^- \[ \]" .mule/todo.md | head -30
5:# Mule queue
86:- [ ] broker-external-bus-intake [legion]: implement the second Intake implementation named in the TODO at internal/broker/intake.go:56-62 — an external-bus-backed (e.g. NATS JetStream) Intake for multi-instance deployment / restart survival, reproducing the embedded broker's per-device ordering, deferred-ack backpressure and QoS 0 drop semantics; an implementation of this interface needs containerised integration tests, so it runs on the Legion Go.
88:- [ ] hygiene-govulncheck [legion]: on the Legion Go, run govulncheck ./... on a fresh ~/astrate clone and file a task line for any REACHABLE advisory it reports (name the CVE and the call path, per the hygiene recipe); this box cannot build it — go install golang.org/x/vuln/cmd/govulncheck@latest is OOM-killed here (3.7GB, 4 cores), so the recipe's highest-priority check never runs unless done there.
89:- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.go:560) with no stop signal, so a block-death racing process shutdown (cmd/astrate/main.go:482-488 drains the Manager but never cancels the goroutine) can rebuild the flow via `mgr.StartFlow` on a background pump (internal/flow/flow.go:173-175) after `Manager.Shutdown` and flip the durable status back to running. Thread a stop channel/context through the Service, checked in the loop's sleep, cancelled on shutdown; verification needs a `[legion]` integration or timing-based test. [auto]
101:- [ ] race-check-store: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/store/... ./internal/housekeeping/... ./migrations/...`. Report any failure to .mule/for-giulio.md with the full race report. Split out of the former single `race-check` line, which timed out running the whole tree at once. [legion] [readonly]
102:- [ ] race-check-engine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/engine/... ./internal/broker/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
103:- [ ] race-check-flow: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/flow/... ./internal/realm/... ./internal/pairing/... ./internal/auth/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
104:- [ ] race-check-appengine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/appengine/... ./internal/observability/... ./internal/httpx/... ./internal/config/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
105:- [ ] race-check-pkg: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./pkg/... ./cmd/... ./internal/testutil/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
200:- [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a locally-run `elixir -e` with `Base.url_decode64!`) with the device ID `"AAAAAAAAAAAAAAAAAAAAAB"` — nonzero unused bits in the 22nd base64url char — and report whether upstream accepts it: `pkg/deviceid/deviceid.go:34-37` uses `base64.RawURLEncoding.Strict()` and deviceid_test.go:176 asserts rejection while the comment claims parity with `Elixir Base.url_decode64!(padding: false)`, which decodes and discards those bits rather than erroring. Astrate may therefore be stricter than upstream on the same wire form (a device registered with such an id on upstream would be rejected here). If upstream rejects it identically, the strict encoding is verified parity and the task is done; if upstream accepts it, escalate the relax-vs-strict call to `.mule/for-giulio.md` with the measurement. Probe first, no code change either way.
203:- [ ] store-devices-inhibit-re-register [legion] [auto]: `RegisterDevice` (internal/store/devices.go:75-91) silently clears the inhibit flag — its `ON CONFLICT ... DO UPDATE SET status = 'registered'` fires for any device with `first_credentials_request IS NULL`, including one the admin inhibited via `SetDeviceInhibited` (devices.go:251-268), which sets `status='inhibited'` on unconfirmed devices too; §5.3 says an inhibited device blocks new credentials and connections, so the re-registration re-opens it. Preserve `'inhibited'` in the SET (e.g. `status = CASE WHEN devices.status = 'inhibited' THEN 'inhibited' ELSE 'registered' END`) and add a Lifecycle case in internal/store/devices_test.go: inhibit an unconfirmed device, re-register, assert status stays inhibited and the secret still rotates. Verify the assert against upstream's register-not-touching-inhibit on the Legion while the integration suite runs.
205:- [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (internal/store/devices.go:130-131, "if several devices share an alias the lowest device ID wins", `ORDER BY id LIMIT 1`) — no test drives it today; add a case in internal/store/devices_test.go that gives two devices the same alias and asserts the lookup resolves the lower ID. Needs the integration DB.
221:- [ ] flow-boot-resumes-stopped-durable-flows [legion] [auto]: RehydrateAutoRestart (cmd/astrate/main.go:216) restarts every flow with auto_restart=true whatever its last status — `ListAutoRestartFlows` filters on `f.auto_restart = true` alone (internal/store/flows.go:125) and nothing ever clears the flag, which is written once at create (flowapi/http.go:141-146 -> service.go:318) and never on stop — so a durable flow stopped through the API is running again after the next boot, while the shutdown mark that says otherwise (MarkRunningFlowsStopped, main.go:497 -> internal/flowapi/service.go:885-901, which writes status="stopped") is never read by that query. The two halves disagree: service.go:325 documents "every durable flow", the shutdown call assumes "only what was running". Say which is the intent, then make them agree — either add `AND f.status = 'running'` to flows.go:125 (making the shutdown mark load-bearing) or delete the mark as dead code — and add a case in internal/store/flows_test.go pinning the chosen rule (a stopped auto_restart row must, or must not, come back). Needs the DB.
246:- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `base64.RawURLEncoding.Strict()` (pkg/deviceid/deviceid.go:37), so of the 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with `ErrInvalid` while `..._lA` is accepted, and the pairing spec's `hw_id` pattern `^[A-Za-z0-9_-]{22}$` (added by docs-sync-pairing-deviceid-base64url) deliberately does not encode the rule because a pattern is a client-side contract and upstream's leniency is unmeasured. The docstring justifies strictness as byte-identical round-tripping, which is a reason, not a parity claim, and docs/COMPATIBILITY.md records no entry for it. Measure upstream (Elixir `Base.url_decode64!(padding: false)`, which does not check those bits) on the Legion Go, then either record the divergence in docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.
259:- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.
262:- [ ] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `"Internal server error"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ head -30 .mule/todo.md; echo "=== total lines:"; wc -l .mule/todo.md; echo "=== unchecked:"; rg -c "^- \[ \]" .mule/todo.md
- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
- [x] triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
- [x] compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).

# Mule queue

One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).
The mule runs the topmost open line. Reorder freely; the order is the priority.

If a line needs more detail than fits, write `.mule/tasks/<slug>.md` and keep the line short.

    tools/mule.sh add "<slug>: <outcome>"
    tools/mule.sh loop

A line tagged `[legion]` needs the Legion Go and is skipped automatically while it is asleep;
the queue moves on to the next runnable line rather than stalling behind it.

A line tagged `[readonly]` is a **standing check**, not a piece of work. It verifies something
and is expected to produce no diff at all — so an empty diff is success, not the `wrote
nothing` failure it would be for any other line. It never gets ticked off, because a gate is
never done. Its output lands in `.mule/reports/<slug>.md` with the sha it ran against, and it
is skipped while that sha is still HEAD: re-checking code that has not moved buys nothing and
spends a call on a free provider.

**`mule/queue` is authoritative for this file.** The mule ticks on the Pi and writes its
progress there, so editing the queue on `main` resurrects completed tasks the moment the two
are merged — which has happened once. Add and reorder tasks on `mule/queue`; let them reach
`main` only when that branch is merged.

A line marked `- [~]` is **parked**: real work, but not tick-sized. The mule only ever
=== total lines:
262 .mule/todo.md
=== unchecked:
15

$ rg -n "BadRequest|InternalServerError" docs/api/astrate_native_api.yaml; echo "=== internal 500 writers:"; rg -n "WriteInternalServerError" internal/native/*.go internal/observability/*.go 2>/dev/null | head; echo "=== native pkg files:"; ls internal/native 2>/dev/null || rg -ln "native_api|astrate_native" --glob '*.go' | head
328:          $ref: "#/components/responses/InternalServerError"
370:          $ref: "#/components/responses/InternalServerError"
458:    InternalServerError:
=== internal 500 writers:
=== native pkg files:

$ rg -n "^  /|operationId:|      \"500\"|      \"503\"" docs/api/astrate_native_api.yaml | head -30
26:  /astrate/v1/health:
28:      operationId: getHealth
44:  /astrate/v1/readiness:
46:      operationId: getReadiness
65:        "503":
77:  /astrate/v1/metrics:
79:      operationId: getMetrics
99:  /appengine/health:
101:      operationId: getAppEngineHealth
117:        "503":
127:  /realmmanagement/health:
129:      operationId: getRealmManagementHealth
145:        "503":
155:  /pairing/health:
157:      operationId: getPairingHealth
173:        "503":
185:  /appengine/version:
187:      operationId: getAppEngineVersion
205:  /realmmanagement/version:
207:      operationId: getRealmManagementVersion
225:  /pairing/version:
227:      operationId: getPairingVersion
245:  /housekeeping/version:
247:      operationId: getHousekeepingVersion
267:  /astrate/v1/{realm}/socket:
269:      operationId: nativeWebSocket
327:        "500":
332:  /appengine/v1/socket/websocket:
334:      operationId: phoenixWebSocket
369:        "500":

$ rg -rn "WriteInternalServerError" --glob '*.go' internal/ | rg -i "native|observability|websocket|socket|phoenix" | head -20
(no output)

$ rg -rn "websocket" --glob '*.go' -l internal/ | head -20; echo "=== 500 writers near websocket/socket:"; rg -rn "StatusInternalServerError" --glob '*.go' internal/ | head -20
internal/engine/triggers/validation.go
internal/appengine/stream/ws.go
internal/appengine/channels/ws_test.go
internal/appengine/channels/protocol.go
internal/appengine/channels/ws.go
internal/appengine/e2e_test.go
=== 500 writers near websocket/socket:
internal/flow/blocks/container/block_test.go:					http.Error(w, "boom", http.n)
internal/flow/blocks/httpblocks_test.go:	srv := staticServer(t, http.n, "text/plain", "boom")
internal/flow/blocks/httpblocks_test.go:			w.WriteHeader(http.n)
internal/engine/forward/http_test.go:		w.WriteHeader(http.n)
internal/engine/forward/http_test.go:		w.WriteHeader(http.n)
internal/appengine/stream/ws.go:		http.Error(w, "streaming unsupported", http.n)
internal/appengine/writeerror_test.go:			http.n, "Internal server error"},
internal/appengine/mirror_test.go:	if rec.Code != http.n {
internal/appengine/patch_test.go:		if rec.Code != http.n {
internal/engine/triggers/actions_test.go:			w.WriteHeader(http.n)
internal/engine/triggers/actions_test.go:		w.WriteHeader(http.n)
internal/engine/triggers/actions_test.go:			w.WriteHeader(http.n) // second: 500 → discard
internal/flowapi/http_test.go:		{errors.New("other"), http.n, "Internal"},

✱ Grep "astrateapi\.WriteInternalServerError" in . · 0 matches

$ rg -rn "WriteInternalServerError" --glob '*.go' . | head -20
./pkg/astarteapi/envelope.go:// n writes the canonical 500 envelope.
./pkg/astarteapi/envelope.go:func n(w http.ResponseWriter) error {
./pkg/astarteapi/envelope_test.go:			write: astarteapi.n,
./internal/pairing/http.go:		_ = astarteapi.n(w)
./internal/appengine/http.go:		_ = astarteapi.n(w)
./internal/appengine/http.go:		_ = astarteapi.n(w)
./internal/appengine/http.go:		_ = astarteapi.n(w)
./internal/flowapi/http.go:		_ = astarteapi.n(w)
./internal/appengine/channels/ws.go:		_ = astarteapi.n(w)
./internal/realm/http.go:		_ = astarteapi.n(w)
./internal/realm/http.go:		_ = astarteapi.n(w)
./internal/housekeeping/http.go:		_ = astarteapi.n(w)
./internal/swagger/swagger_test.go:// astarteapi.WriteBadRequest and every 500 through n
./internal/swagger/swagger_test.go:// astarteapi.WriteBadRequest and every 500 through n
./internal/auth/middleware.go:				_ = astarteapi.n(w)

✱ Grep "astrate/v1/\{realm\}/socket|astrate/v1/health" 5 matches
→ Read internal/appengine/stream/ws.go [limit=70]
→ Read internal/appengine/stream/ws.go [offset=70, limit=40]

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -462,5 +462,5 @@
 schema:
   $ref: "#/components/schemas/ErrorDetail"
 example:
   errors:
-    detail: Internal Server Error
+    detail: Internal server error

$ rg -n "^    [A-Za-z]+:$" docs/api/astrate_native_api.yaml | sed -n '1,60p'
16:    variables:
27:    get:
45:    get:
78:    get:
100:    get:
128:    get:
156:    get:
186:    get:
206:    get:
226:    get:
246:    get:
268:    get:
333:    get:
374:    RealmName:
383:    HealthStatus:
392:    ReadinessStatus:
406:    DataEnvelopeHealthStatus:
418:    DataEnvelopeVersion:
426:    ErrorDetail:
438:    Unauthorized:
448:    Forbidden:
458:    InternalServerError:

Three response components: `Unauthorized`, `Forbidden`, `InternalServerError`. Adding the test beside the others.
← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -609,8 +609,60 @@
 		}
 	}
 }
 
+// TestNativeErrorDetailExamples pins the error-detail examples in the native
+// spec's components.responses to the frozen canonical strings the wire emits
+// (astarteapi/envelope.go), the same guard the pairing, realm management and
+// app engine specs get from TestPairingErrorDetailExamples,
+// TestRealmManagementErrorDetailExamples and TestAppEngineErrorDetailExamples.
+// Only the 500 was wrong: it showed "Internal Server Error", the capital-S
+// spelling reconstructed rather than measured, while every 500 on the two
+// websocket operations the component is $ref'd from answers through
+// astarteapi.WriteInternalServerError (auth.Middleware's upgrade failure,
+// internal/auth/middleware.go, and the channels socket,
+// internal/appengine/channels/ws.go) with the canonical "Internal server error".
+// The 401 and 403 are already the canonical spellings. This spec has no
+// BadRequest component at all: the native surface is the compat health and
+// version endpoints plus the two sockets, and none of them answers 400 through
+// WriteBadRequest, so there is nothing here to pin.
+func TestNativeErrorDetailExamples(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astrate_native_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astrate_native_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	const detailPrefix = "              detail: "
+
+	for _, tc := range []struct {
+		response string
+		want     string
+	}{
+		{"    Unauthorized:", astarteapi.DetailUnauthorized},
+		{"    Forbidden:", astarteapi.DetailForbidden},
+		{"    InternalServerError:", astarteapi.DetailInternalServerError},
+	} {
+		block := componentBlock(t, lines, tc.response)
+
+		got := ""
+		for _, l := range block {
+			if strings.HasPrefix(l, detailPrefix) {
+				got = strings.TrimPrefix(l, detailPrefix)
+				break
+			}
+		}
+		if got == "" {
+			t.Errorf("response %q carries no %q example", tc.response, "detail")
+			continue
+		}
+		if got != tc.want {
+			t.Errorf("response %q example detail = %q, want the canonical %q",
+				strings.TrimSpace(tc.response), got, tc.want)
+		}
+	}
+}
+
 // TestRealmManagement422ValidationDetails pins the 422 examples the three
 // non-interface realm-management operations carry to the details the wire
 // actually emits, so a generated client that copies one does not learn a
 // message no handler can produce.

$ go test ./internal/swagger/ -run TestNativeErrorDetailExamples -v 2>&1 | tail -15
=== RUN   TestNativeErrorDetailExamples
--- PASS: TestNativeErrorDetailExamples (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.019s

Now confirming the test actually fails without the one-line fix.
← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -462,5 +462,5 @@
 schema:
   $ref: "#/components/schemas/ErrorDetail"
 example:
   errors:
-    detail: Internal server error
+    detail: Internal Server Error

$ go test ./internal/swagger/ -run TestNativeErrorDetailExamples 2>&1 | tail -10
--- FAIL: TestNativeErrorDetailExamples (0.00s)
    swagger_test.go:659: response "InternalServerError:" example detail = "Internal Server Error", want the canonical "Internal server error"
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.014s
FAIL

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -462,5 +462,5 @@
 schema:
   $ref: "#/components/schemas/ErrorDetail"
 example:
   errors:
-    detail: Internal Server Error
+    detail: Internal server error

Now the docs build and the YAML-load check.
$ cat docs/Makefile; echo "=== mkdocs?"; command -v mkdocs || echo "no mkdocs"; python3 -c "import yaml; print('pyyaml ok')" 2>&1 | tail -1
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
=== mkdocs?
no mkdocs
ModuleNotFoundError: No module named 'yaml'

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

$ git status --porcelain | head -20; echo "=== users of embeddedYAMLFilenames:"; rg -n "embeddedYAMLFilenames|func Test.*YAML|LoadSpec|swagger.Load" internal/swagger/*.go | head -20
M docs/api/astrate_native_api.yaml
 M internal/swagger/swagger_test.go
=== users of embeddedYAMLFilenames:
internal/swagger/swagger_test.go:114:	want, err := embeddedYAMLFilenames()
internal/swagger/swagger_test.go:1017:func embeddedYAMLFilenames() ([]string, error) {

→ Read internal/swagger/swagger_test.go [limit=70, offset=100]

$ rg -n "yaml\.|Unmarshal|Parse" internal/swagger/swagger.go | head -20; echo "=== tests referencing parse:"; rg -n "func Test" internal/swagger/swagger_test.go | head -40
=== tests referencing parse:
23:func TestMount(t *testing.T) {
82:func TestMountSubPanicsOnBrokenFS(t *testing.T) {
95:func TestSpecs(t *testing.T) {
128:func TestRealmManagement403(t *testing.T) {
139:func TestAppEngine403(t *testing.T) {
203:func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
246:func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
321:func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {
381:func TestAppEngineAddGroupDevice422Documented(t *testing.T) {
425:func TestAppEngineDataDelete400Documented(t *testing.T) {
470:func TestPairingErrorDetailExamples(t *testing.T) {
521:func TestRealmManagementErrorDetailExamples(t *testing.T) {
573:func TestAppEngineErrorDetailExamples(t *testing.T) {
628:func TestNativeErrorDetailExamples(t *testing.T) {
682:func TestRealmManagement422ValidationDetails(t *testing.T) {
759:func TestPairingDeviceIDEncodingDocumented(t *testing.T) {
841:func TestRealmManagementDeviceIDEncodingDocumented(t *testing.T) {
905:func TestPairingUnregisterDeviceSemanticsDocumented(t *testing.T) {

$ rg -n "yaml" go.mod
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

← Write /tmp/opencode/yamlcheck/main.go
Wrote file successfully.

$ go run /tmp/opencode/yamlcheck/main.go docs/api/*.yaml 2>&1 | tail -10
docs/api/astarte_appengine_api.yaml: parses OK
docs/api/astarte_housekeeping_api.yaml: parses OK
docs/api/astarte_pairing_api.yaml: parses OK
docs/api/astarte_realm_management_api.yaml: parses OK
docs/api/astrate_native_api.yaml: parses OK

$ gofmt -l . | head; echo "=== gofmt done"; go vet ./... 2>&1 | tail -5
=== gofmt done

$ go test ./... 2>&1 | rg -v "^ok|no test files" | head -30; echo "=== exit ${PIPESTATUS[0]}"
/bin/sh: 1: Bad substitution

$ go test ./... > /tmp/opencode/gate.txt 2>&1; echo "exit=$?"; rg -v "^ok|no test files" /tmp/opencode/gate.txt | head -30
exit=0

$ go test -race ./internal/swagger/... 2>&1 | tail -5; echo "=== tail of todo:"; sed -n '255,262p' .mule/todo.md | cut -c1-120
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.025s
FAIL
=== tail of todo:
- [!] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy alias fields to `InterfaceDefiniti
- [x] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/
- [x] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_a
- [!] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_re
- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-contain
- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appe
- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_ho
- [ ] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_n

→ Read .mule/todo.md [offset=261]
✗ Edit .mule/todo.md failed
Error: Could not find oldString in the file. It must match exactly, including whitespace, indentation, and line endings.

$ tail -c 120 .mule/todo.md | od -c | tail -4
0000120   h   e       S   w   a   g   g   e   r   -   r   e   f   e   r
0000140   e   n   c   e   d       Y   A   M   L       s   t   i   l   l
0000160       l   o   a   d   s   .  \n
0000170

$ cat >> .mule/todo.md <<'EOF'
- [ ] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
EOF
tail -1 .mule/todo.md | cut -c1-80
- [ ] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of

$ git status --porcelain; echo "=== diff stat:"; git diff --stat
M .mule/todo.md
 M docs/api/astrate_native_api.yaml
 M internal/swagger/swagger_test.go
=== diff stat:
 .mule/todo.md                    |  1 +
 docs/api/astrate_native_api.yaml |  2 +-
 internal/swagger/swagger_test.go | 52 ++++++++++++++++++++++++++++++++++++++++
 3 files changed, 54 insertions(+), 1 deletion(-)

Done: corrected `InternalServerError`'s detail example in `docs/api/astrate_native_api.yaml` to the canonical `"Internal server error"` and added `TestNativeErrorDetailExamples` pinning all three of the spec's response examples to the `astarteapi` constants (verified it fails on the old string, passes on the new).

Files: docs/api/astrate_native_api.yaml (yaml:466), internal/swagger/swagger_test.go (new test at :628), .mule/todo.md (one appended line)

Verified: `go vet ./...` clean; `go test ./...` exit 0; `gofmt -l .` empty; `go test ./internal/swagger/ -run TestNativeErrorDetailExamples` pass, and FAIL (`example detail = "Internal Server Error", want the canonical "Internal server error"`) with the yaml change reverted; all five `docs/api/*.yaml` parse (throwaway `gopkg.in/yaml.v3` program in /tmp). `make -C docs build` fails with `mkdocs: No such file or directory` — mkdocs is not installed on this box, unrelated to the change; its `sync` step ran and `git status` shows no stray files. `-race` cannot run here at all: `ThreadSanitizer: unsupported VMA range, Found 39 - Supported 48`, i.e. this box has the Pi's 39-bit VMA kernel, so the gate is the non-race one.

Unsure: the task line claims "no spec carries a reconstructed capitalisation any more" — that is not true. `docs/api/astarte_housekeeping_api.yaml` still carries `detail: Bad Request` (:362 schema-level, :385 BadRequest) and `detail: Internal Server Error` (:477); pairing/appengine/realm-management are clean (`Not Found` there is canonical). I did not fix it: it is already queued as the blocked line `docs-sync-hk-error-example-capitalisation`, so I only appended the tick-sized first half of it as a new line rather than duplicating or silently expanding scope. Also, the native spec's 500 is documented on the two websocket operations; both answer via `astarteapi.WriteInternalServerError` (internal/auth/middleware.go upgrade failure, internal/appengine/channels/ws.go), which is what the test comment cites — I did not run a live probe of those two paths.

Follow-ups: one line appended to .mule/todo.md — `docs-sync-hk-error-detail-examples-split`.
