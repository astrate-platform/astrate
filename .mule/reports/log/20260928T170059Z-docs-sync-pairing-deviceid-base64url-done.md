slug: docs-sync-pairing-deviceid-base64url
verdict: done
at:  32693cf
ran: 2026-09-28T17:00:59Z on DietPi in 480s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ rg -n "base64" docs/api/astarte_pairing_api.yaml
360:      description: The device hardware ID (base64-encoded 128-bit).
373:          description: Device hardware ID (base64-encoded 128-bit).
402:          description: One-time device credentials secret (base64-encoded).
618:                    - "is not a valid base64 encoded 128 bits id"

$ ls docs/ && cat docs/Makefile 2>/dev/null | head -40
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

→ Read docs/api/astarte_pairing_api.yaml [offset=100, limit=50]
→ Read docs/api/astarte_pairing_api.yaml [offset=340, limit=60]

$ rg -n "deviceID|hw_id" docs/api/astarte_pairing_api.yaml
123:                hw_id: dT6hS2W9TT6LEnP25ks_lg
149:  /pairing/v1/{realm}/agent/devices/{deviceID}:
172:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
224:  /pairing/v1/{realm}/devices/{deviceID}:
258:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:
357:      name: deviceID
369:      required: [hw_id]
371:        hw_id:
596:              summary: Blank hw_id
599:                  hw_id:
614:              summary: Invalid hw_id format
617:                  hw_id:

→ Read pkg/deviceid/deviceid.go

$ rg -n "base64|deviceid|device ID" docs/JSON-PAYLOAD-PROFILE.md | head -30; echo "--- todo:"; cat .mule/todo.md 2>/dev/null | tail -30
84:| `binaryblob`   | base64 string, **standard alphabet, padded**                    |
130:{ "data": { "hw_id": "<22-char base64url id>", "initial_payload_format": "json" } }
--- todo:
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
- [x] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion Go — the device-delete leg added by docs-sync-rm-delete-device-async-operation-param (register a device, `DELETE /devices/{id}?async_operation=…` answers 204 and the row is gone on both values) compiles clean under `go vet -tags integration ./internal/realm/` but has never been executed, because the Pi has no database; until it runs, TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go is the only assertion of that leg that has actually been verified.
- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] forward-static-header-validation [auto]: `New` (internal/engine/forward/http.go:38-71) validates the URL and the method at boot precisely because of the rule its own comment states at 42-45 — "an unusable endpoint must fail here rather than surface once per delivery" — but `cfg.StaticHeaders` is stored unvalidated (http.go:68) and applied per delivery with `req.Header.Set` (http.go:110-112), and net/http checks header names and values only at write time. Measured on the real net/http path with a throwaway program: a name of `X Bad Name` or `X:Foo` fails per delivery with `net/http: invalid header field name`, and a value containing a newline fails with `net/http: invalid header field value for "X-Foo"` (a plain space in a value is legal). So one typo in `triggers.forward.static_headers` makes the process boot clean, accept triggers, and then count `astrate_engine_trigger_deliveries_total{outcome="failed"}` and log "custom trigger action failed" (internal/engine/triggers/actions.go:499) on every custom action forever. Validate in `New`: the RFC 7230 field-name grammar is exactly the token grammar `validMethod` (http.go:77-95) already implements, so reuse that predicate per name (rename it if you like) and add a CR/LF rejection per value. Add a table case to `TestNewRejectsBadConfig` (internal/engine/forward/http_test.go:271) asserting `New` returns an error for a bad name and a bad value and that no request is attempted; this test fails today because `New` returns nil error for both.
- [x] forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
- [x] forward-status-error-body [auto]: a non-2xx forward returns `forward: status %d` and nothing else (internal/engine/forward/http.go:121-123) after copying the **entire** response body into `io.Discard` with no limit (http.go:117-120) — the one thing an operator needs to debug a 500 from the bus is destroyed, and the sibling webhook request in the same codebase does bound its drain with `io.LimitReader(resp.Body, 1<<20)` (internal/engine/triggers/actions.go:626), so the two near-identical request paths have drifted on exactly that line. Read a bounded prefix (e.g. 512 bytes) of the body into the returned error and bound the discard the way actions.go:626 does. Add a case beside `TestStatusTable` (internal/engine/forward/http_test.go:164) where the handler answers 500 with body `{"error":"boom"}` and assert the error mentions both `500` and `boom`; it fails today because the error carries only the status code.
- [x] forward-envelope-bytes-test [auto]: `marshalEnvelope` (internal/engine/forward/envelope.go:19-29) has no direct test — `rg marshalEnvelope` finds exactly two callers (http.go:99, nats.go:45) and no test file — while envelope.go:6-8 makes the strong claim that "every Forwarder implementation in this package must produce byte-identical envelopes for the same inputs". Both suites check the shape only by unmarshalling into `bodyShape` (http_test.go:24-29), which cannot see key order or the exact bytes and would pass unchanged if a field were renamed on both sides. Add internal/engine/forward/envelope_test.go with a table pinning the exact output for (nil, nil) → `{"realm":"r","trigger":"t","action":null,"event":null}`, for the empty-non-nil pair the rule at 16-18 is about, and for a valued pair; container-free, no HTTP server needed.
- [x] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.
- [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml — the string `"403"` appears nowhere in the file and `components` defines no `Forbidden` response, yet every route is wrapped by `a.require = mw.RequireRealm(auth.ClaimAppEngine)` (internal/appengine/http.go:35), which answers `403 Forbidden` for a verified a_aea JWT whose grants do not authorize the method+path (`astarteapi.WriteForbidden`, internal/auth/middleware.go:90; `DetailForbidden`, pkg/astarteapi/envelope.go:41; pinned today by internal/appengine/http_test.go:184-185). Add a `Forbidden` response component copied from the realm-management one (docs/api/astarte_realm_management_api.yaml:1277-1285 — `ErrorDetail` schema, example `errors.detail: Forbidden`) and `$ref` it from all 32 operations, which is exactly the parity docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go:123-170) pins; add the same docs-half test for the AppEngine spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-patch-by-alias-409 [auto]: add the missing `409` response to `patchDeviceByAlias` in docs/api/astarte_appengine_api.yaml (yaml:241-291, currently 200/400/401/404/422/500) — the handler resolves the alias then calls the same `applyPatch` as the device- and group-scoped PATCH, and a rename to an alias already in use returns `ErrAliasAlreadyInUse` (internal/appengine/service.go:337) which `writeError` maps to `409 ConflictAliasInUse` (internal/appengine/http.go:650-651); the other two PATCH operations already document that response thanks to docs-sync-appengine-query-params-status and docs-sync-appengine-group-patch-status, so this is the last PATCH of the three with a hole. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml (yaml:939-980, currently 201/400/401/404/409/500) — a well-formed body whose `device_id` does not parse is rejected as a field error (`FieldErrors{"device_id": {"is not a valid device id"}}`, internal/appengine/service.go:631-634) and answered 422 with the `FieldErrorsDetail` body (`astarteapi.WriteFieldErrors`, internal/appengine/http.go:647), so `$ref` the existing `ValidationErrors` response (yaml:1707-1716) rather than `BadRequest` — the documented 400 only covers a body `DecodeData` cannot parse (internal/appengine/http.go:488-491). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data DELETE operations in docs/api/astarte_appengine_api.yaml (`deleteDataByAlias` yaml:519-547, `deleteData` yaml:756-784, `deleteDataInGroup` yaml:1285-1313 — each currently 204/401/404/405/500) — unsetting a path that matches no endpoint mapping returns `engine.ErrPathNotFound` (`UnsetServerProperty`, internal/engine/serverdata.go:241-244) which `writeError` answers as `400 "Endpoint not found"` (internal/appengine/http.go:679-680), the same 400 the PUT/POST twins already document via `BadRequest`; internal/appengine/writeerror_test.go pins the mapping. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 200/401/500) — a `from_token` that does not parse as a device id is refused with `ErrValidation: invalid cursor` (internal/appengine/service.go:150-155) which `writeError` maps to 422 with the flat `{"errors":{"detail":"invalid cursor"}}` body, i.e. the `ErrorDetail` shape and not `FieldErrorsDetail`, so `$ref` the same response the group twin uses for the identical parameter (`ValidationErrors`, yaml:934-935) only after fixing its example, or add a dedicated one. Extend the `from_token` parameter description (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:556) and `InternalServerError` shows `detail: Internal Server Error` (yaml:638), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37, golden fixture pkg/astarteapi/testdata/error_bad_request.json) and `"Internal server error"` (envelope.go:58, testdata/error_internal.json), both measured on upstream 1.2.0 (envelope.go:33-36, 55-57). The other five detail examples on this file (`Not Found`, `Device not found`, `Unauthorized`, `Forbidden`, `Too Many Requests`) are already correct — change only those two. The identical wrong pair is in the other four specs (appengine 1732/1792, housekeeping 362/385/477, realm_management 1265/1348, native 466); leave those for a run of their surface. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
- [x] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pairing_api.yaml, which shows `version: "0.1.0-test"` (yaml:252) — a value no build can produce. The field is filled from `pairing.Config.Version` (internal/pairing/service.go:302), wired from `main.version` (cmd/astrate/main.go:151), whose default is `"0.1.0-dev"` (main.go:53) and which is otherwise set only via `-ldflags "-X main.version=..."`. Use `"0.1.0-dev"`, matching the four compat-version examples docs-sync-native-version-value already corrected to that (docs/api/astrate_native_api.yaml:203, 223, 243, 263). Do NOT document the service's own fallback `DefaultVersion = "0.1.0-astrate"` (service.go:86-88): it is unreachable from `main`, which always passes a non-empty value. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` path parameter (docs/api/astarte_pairing_api.yaml:360) and `RegisterRequest.hw_id` (yaml:373) — but `deviceid.Parse` accepts only the 22-character **unpadded base64url** form, in strict mode, and rejects the standard alphabet's `+` and `/` and any `=` padding with an error wrapping `ErrInvalid` (pkg/deviceid/deviceid.go:26, 37, 39-60; upstream parity cited at 36). Measured: with the same `encoding` the parser uses, `dT6hS2W9TT6LEnP25ks+lg` and `.../lg` fail at byte 19 and the padded 23-character form fails at byte 21, so a client generated from the current spec can send a well-formed standard-base64 ID and get 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}` (internal/pairing/http.go:307-309). State the exact encoding and the 22-character length in both descriptions and add a `pattern`/`minLength`/`maxLength` to `hw_id`. The existing `hw_id: dT6hS2W9TT6LEnP25ks_lg` example (yaml:123) is valid and stays. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in docs/api/astarte_pairing_api.yaml (yaml:374-376), an unconstrained string described only as "Optional Astrate extension specifying the device's initial payload format" — so a generated client can send any string, while `Service.Register` accepts only `"bson"` or `"json"` and rejects everything else with 422 `{"errors":{"initial_payload_format":["is invalid"]}}` (internal/pairing/service.go:179-181 -> internal/pairing/http.go:310-312; the allowed set is stated at service.go:173). Add `enum: [bson, json]` plus a note that absent/empty means unset, and add the three missing `ValidationError` examples while you are in the component: `initial_payload_format: ["is invalid"]` (http.go:310-312), `csr: ["can't be blank"]` (http.go:173-177) and `client_crt: ["can't be blank"]` (http.go:253-257) — the component currently shows only the four hw_id/error_name rows (yaml:594-618) and the `csr`/`client_crt` blank-field cases are what a device actually hits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` description in docs/api/astarte_pairing_api.yaml (yaml:153), "Removes a device from the realm. Returns 204 on success.", which reads as a device-and-data deletion. The row and every byte of its stored data survive: `store.UnregisterDevice` only clears the credential trail (`credentials_secret_hash = ''`, `first_credentials_request = NULL`, `cert_serial = NULL`, `cert_aki = NULL`) and flips `status` back to `'registered'` (internal/store/devices.go:93-110), and the service docstring says so outright ("its data is kept", internal/pairing/service.go:227-229). State the real semantics: the device becomes registrable again, its interfaces, datastream data and group memberships are retained, and a second DELETE is a 204 no-op (the row is still there, so RowsAffected is 1, devices.go:106-108) — not the documented 404. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ rg -rn "docs/api" --glob '!docs/site' -l . | head -20; echo ---; rg -n "docs" Makefile | head -20
./docs/handoff/phase-1-handoff.md
./docs/handoff/phase-5-handoff.md
./docs/handoff/phase-7-memory.md
./docs/handoff/phase-2-handoff.md
./docs/handoff/phase-4-memory.md
./docs/handoff/phase-3-handoff.md
./docs/handoff/phase-4-handoff.md
./docs/handoff/phase-3-memory.md
./docs/handoff/phase-2-memory.md
./docs/handoff/phase-6-handoff.md
./docs/handoff/phase-1-memory.md
---
1:# Astrate build & verification entry points (docs/ROADMAP.md §1.2, verification tiers §0.2).
53:		echo "conformance harness lands in M9 (docs/ROADMAP.md §10); nothing to run yet"; \

$ ls internal/swagger/ && wc -l internal/swagger/*.go && rg -n "^func Test" internal/swagger/swagger_test.go
swagger.go
swagger_test.go
   69 internal/swagger/swagger.go
  597 internal/swagger/swagger_test.go
  666 total
19:func TestMount(t *testing.T) {
78:func TestMountSubPanicsOnBrokenFS(t *testing.T) {
91:func TestSpecs(t *testing.T) {
124:func TestRealmManagement403(t *testing.T) {
135:func TestAppEngine403(t *testing.T) {
199:func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
242:func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
317:func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {
377:func TestAppEngineAddGroupDevice422Documented(t *testing.T) {
421:func TestAppEngineDataDelete400Documented(t *testing.T) {
466:func TestPairingErrorDetailExamples(t *testing.T) {

→ Read internal/swagger/swagger_test.go [limit=130]
→ Read internal/swagger/swagger_test.go [offset=460, limit=140]
→ Read internal/swagger/swagger_test.go [offset=300, limit=130]

$ rg -n "pattern:|minLength|maxLength|enum:" docs/api/*.yaml | head -40
docs/api/astrate_native_api.yaml:315:            enum: [websocket, sse]
docs/api/astrate_native_api.yaml:389:          enum: [ok]
docs/api/astrate_native_api.yaml:398:          enum: [ok, unavailable]
docs/api/astrate_native_api.yaml:416:              enum: [ok, unhealthy]
docs/api/astarte_appengine_api.yaml:1480:        enum: [ascending, descending]
docs/api/astarte_appengine_api.yaml:1489:        enum: [structured, table, disjoint_tables]
docs/api/astarte_pairing_api.yaml:431:          enum: [confirmed, pending, inhibited]
docs/api/astarte_pairing_api.yaml:449:          enum: [true]
docs/api/astarte_pairing_api.yaml:457:          enum: [EXPIRED, INVALID, REVOKED]
docs/api/astarte_pairing_api.yaml:468:          enum: [false]
docs/api/astarte_pairing_api.yaml:499:          enum: [ok]
docs/api/astarte_realm_management_api.yaml:939:          maxLength: 128
docs/api/astarte_realm_management_api.yaml:940:          pattern: '^([a-zA-Z][a-zA-Z0-9]*\.([a-zA-Z0-9][a-zA-Z0-9-]*\.)*)?[a-zA-Z][a-zA-Z0-9]*$'
docs/api/astarte_realm_management_api.yaml:949:          enum: [datastream, properties]
docs/api/astarte_realm_management_api.yaml:952:          enum: [device, server]
docs/api/astarte_realm_management_api.yaml:955:          enum: [individual, object]
docs/api/astarte_realm_management_api.yaml:979:          enum: [double, integer, boolean, longinteger, string, binaryblob, datetime, doublearray, integerarray, booleanarray, longintegerarray, stringarray, binaryblobarray, datetimearray]
docs/api/astarte_realm_management_api.yaml:982:          enum: [unreliable, guaranteed, unique]
docs/api/astarte_realm_management_api.yaml:986:          enum: [discard, volatile, stored]
docs/api/astarte_realm_management_api.yaml:995:          enum: [no_ttl, use_ttl]
docs/api/astarte_realm_management_api.yaml:1044:          enum: [data_trigger, device_trigger]
docs/api/astarte_realm_management_api.yaml:1047:          enum: [incoming_data, value_change, value_change_applied, path_created, path_removed, value_stored, device_registered, device_connected, device_disconnected, device_empty_cache_received, device_error, incoming_introspection, interface_added, interface_removed, interface_minor_updated]
docs/api/astarte_realm_management_api.yaml:1065:          enum: ["*", "==", "!=", ">", ">=", "<", "<=", contains, not_contains]
docs/api/astarte_realm_management_api.yaml:1078:          enum: [delete, get, head, options, patch, post, put]
docs/api/astarte_realm_management_api.yaml:1104:          minLength: 1
docs/api/astarte_realm_management_api.yaml:1105:          maxLength: 128
docs/api/astarte_realm_management_api.yaml:1106:          pattern: '^[a-zA-Z0-9_.~-]{1,128}$'
docs/api/astarte_realm_management_api.yaml:1132:              enum: [any_error, client_error, server_error]
docs/api/astarte_realm_management_api.yaml:1151:          enum: [discard, retry]

→ Read docs/api/astarte_pairing_api.yaml [offset=585, limit=45]

$ mkdir -p /tmp/opencode/probe && cd /tmp/opencode/probe && printf '%s\n' 'module probe' '' 'go 1.24' '' 'require github.com/astrate-platform/astrate v0.0.0' '' 'replace github.com/astrate-platform/astrate => /root/astrate-mule' > go.mod && cat go.mod
module probe

go 1.24

require github.com/astrate-platform/astrate v0.0.0

replace github.com/astrate-platform/astrate => /root/astrate-mule

← Write /tmp/opencode/probe/main.go
Wrote file successfully.

$ GOFLAGS=-mod=mod go run . 2>&1 | tail -20
"dT6hS2W9TT6LEnP25ks_lg" parse_err=false pattern_match=true
"dT6hS2W9TT6LEnP25ks+lg" parse_err=true  pattern_match=false
"dT6hS2W9TT6LEnP25ks/lg" parse_err=true  pattern_match=false
"dT6hS2W9TT6LEnP25ks_lg=" parse_err=true  pattern_match=false
"dT6hS2W9TT6LEnP25ks_lx" parse_err=true  pattern_match=true
"dT6hS2W9TT6LEnP25ks_lA" parse_err=false pattern_match=true
EncodedLen = 22

$ rg -n -i "device.?id|canonical" docs/COMPATIBILITY.md | head -30
40:- Certificate `Subject CN = <realm>/<device_id>`; 128-bit serials; clientAuth.
94:     against `<device_id>/<interface_name><match_path>` — the match path is
96:     the bare `<device_id>`. Astrate built `<device_id>/<interface_name>` for
105:     device_trigger must carry `device_id: "*"` inside `simple_trigger`
106:     (refusal reason: `device_id must be * for group triggers`, the mirror
111:   - **A device trigger's `device_id` must sit inside `simple_trigger` and
112:     equal the request's own.** Upstream refuses a watch whose `device_id` is
114:     and refuses a wildcard `device_id: "*"`, both with the misleading reason
115:     `unauthorized`. Astrate used to fall back to the top-level `device_id` and
306:15. **Canonical detail casing follows measurement, not memory** — upstream's
320:    contract — key grammar `<device_id>/<interface></path…>` (upstream's

$ rg -rn "Strict\(\)|canonical" pkg/deviceid/*_test.go internal/pairing/*.go | head -20
pkg/deviceid/deviceid_test.go:	// Output is always lowercase n.
pkg/deviceid/deviceid_test.go:		t.Errorf("UUID() = %q, want lowercase n", got)
pkg/deviceid/deviceid_test.go:		{"non-n trailing bits", "AAAAAAAAAAAAAAAAAAAAAB"},

→ Read internal/swagger/swagger_test.go [limit=75, offset=228]
Now the two description edits and the `hw_id` constraints:

← Edit docs/api/astarte_pairing_api.yaml
Index: /root/astrate-mule/docs/api/astarte_pairing_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_pairing_api.yaml
+++ /root/astrate-mule/docs/api/astarte_pairing_api.yaml
@@ -356,9 +356,13 @@
   DeviceID:
     name: deviceID
     in: path
     required: true
-    description: The device hardware ID (base64-encoded 128-bit).
+    description: |
+      The device hardware ID: a 128-bit value in the 22-character unpadded
+      base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and `_`,
+      with no `=` padding. The standard-alphabet characters `+` and `/`, and
+      a padded 23-character spelling of the same value, are not accepted.
     schema:
       type: string
 
 schemas:

← Edit docs/api/astarte_pairing_api.yaml
Index: /root/astrate-mule/docs/api/astarte_pairing_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_pairing_api.yaml
+++ /root/astrate-mule/docs/api/astarte_pairing_api.yaml
@@ -373,9 +373,19 @@
 required: [hw_id]
 properties:
   hw_id:
     type: string
-    description: Device hardware ID (base64-encoded 128-bit).
+    minLength: 22
+    maxLength: 22
+    pattern: '^[A-Za-z0-9_-]{22}$'
+    description: |
+      Device hardware ID: a 128-bit value in the 22-character unpadded
+      base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and
+      `_`, with no `=` padding. The standard-alphabet characters `+` and
+      `/`, and a padded 23-character spelling of the same value, are
+      rejected with 422 and the
+      `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}`
+      body.
   initial_payload_format:
     type: string
     description: Optional Astrate extension specifying the device's initial payload format.
 

→ Read docs/api/astarte_pairing_api.yaml [limit=36, offset=352]

$ head -5 go.mod; echo ---; rg -n "^(test|lint|vet|check|fmt)" Makefile | head -20
module github.com/astrate-platform/astrate

go 1.26.1

// Pinned dependency decisions — docs/ROADMAP.md §1.1 (frozen in docs/DESIGN.md).
---
31:lint:
37:test:
41:test-integration:
45:test-e2e:
49:test-conformance:

→ Read Makefile [limit=30, offset=25]
Now the docs-half test:

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -501,8 +501,86 @@
 		}
 	}
 }
 
+// TestPairingDeviceIDEncodingDocumented guards that the pairing spec states the
+// exact wire form of a device identifier — the 22-character unpadded base64url
+// string deviceid.Parse accepts — instead of the "base64-encoded 128-bit"
+// wording, which read as standard base64 and let a generated client send a
+// well-formed `+`/`/` identifier or a padded 23-character spelling and collect
+// a 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}`
+// (internal/pairing/http.go). The shipped pattern is compiled here and held to
+// the same spellings the parser is, so the two cannot drift apart: the bounds
+// and the description wording are derived from deviceid.EncodedLen rather than
+// written out, so a change to the accepted length fails this test too.
+func TestPairingDeviceIDEncodingDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	wantLen := strconv.Itoa(deviceid.EncodedLen)
+	wantDesc := fmt.Sprintf("%s-character unpadded", wantLen)
+
+	param := strings.Join(componentBlock(t, lines, "    DeviceID:"), "\n")
+	for _, want := range []string{wantDesc, "base64url"} {
+		if !strings.Contains(param, want) {
+			t.Errorf("DeviceID parameter description does not say %q", want)
+		}
+	}
+	if strings.Contains(param, "(base64-encoded 128-bit)") {
+		t.Error("DeviceID parameter still describes the ID as plain base64")
+	}
+
+	block := propertyBlock(t, lines, "hw_id")
+	hwID := strings.Join(block, "\n")
+	for _, want := range []string{wantDesc, "base64url"} {
+		if !strings.Contains(hwID, want) {
+			t.Errorf("RegisterRequest.hw_id description does not say %q", want)
+		}
+	}
+	for _, want := range []string{
+		"          minLength: " + wantLen,
+		"          maxLength: " + wantLen,
+	} {
+		if !containsLine(block, want) {
+			t.Errorf("RegisterRequest.hw_id is missing line %q", want)
+		}
+	}
+
+	pattern := ""
+	for _, l := range block {
+		if v, ok := strings.CutPrefix(l, "          pattern: "); ok {
+			pattern = strings.Trim(v, `'"`)
+			break
+		}
+	}
+	if pattern == "" {
+		t.Fatal("RegisterRequest.hw_id carries no pattern")
+	}
+	re, err := regexp.Compile(pattern)
+	if err != nil {
+		t.Fatalf("hw_id pattern %q does not compile: %v", pattern, err)
+	}
+
+	if !re.MatchString("dT6hS2W9TT6LEnP25ks_lg") {
+		t.Errorf("hw_id pattern %q rejects the spec's own example dT6hS2W9TT6LEnP25ks_lg", pattern)
+	}
+	for _, spelling := range []string{
+		"dT6hS2W9TT6LEnP25ks+lg", // standard-alphabet '+'
+		"dT6hS2W9TT6LEnP25ks/lg", // standard-alphabet '/'
+		"dT6hS2W9TT6LEnP25ks_lg=", // padded
+	} {
+		if _, err := deviceid.Parse(spelling); err == nil {
+			t.Errorf("deviceid.Parse accepts %q; this test's premise no longer holds", spelling)
+		}
+		if re.MatchString(spelling) {
+			t.Errorf("hw_id pattern %q accepts %q, which deviceid.Parse rejects", pattern, spelling)
+		}
+	}
+}
+
 // operationBlock returns the lines of the operation with the given operationId,
 // up to the next operation, path, or the components section.
 func operationBlock(t *testing.T, lines []string, operationID string) []string {
 	t.Helper()

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -1,20 +1,23 @@
 package swagger
 
 import (
 	"embed"
+	"fmt"
 	"io"
 	"io/fs"
 	"net/http"
 	"net/http/httptest"
 	"reflect"
 	"regexp"
 	"sort"
+	"strconv"
 	"strings"
 	"testing"
 
 	docs "github.com/astrate-platform/astrate/docs"
 	"github.com/astrate-platform/astrate/pkg/astarteapi"
+	"github.com/astrate-platform/astrate/pkg/deviceid"
 )
 
 func TestMount(t *testing.T) {
 	mux := http.NewServeMux()

$ gofmt -l internal/swagger/ && go vet ./internal/swagger/ && go test ./internal/swagger/ 2>&1 | tail -20
internal/swagger/swagger_test.go
ok  	github.com/astrate-platform/astrate/internal/swagger	0.077s

$ gofmt -w internal/swagger/swagger_test.go && gofmt -l . ; go test ./internal/swagger/ -run TestPairingDeviceIDEncodingDocumented -v 2>&1 | tail -10
=== RUN   TestPairingDeviceIDEncodingDocumented
--- PASS: TestPairingDeviceIDEncodingDocumented (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.018s

Now verify the test actually fails without the YAML change (MULE.md step 4):
$ cp docs/api/astarte_pairing_api.yaml /tmp/opencode/pairing.bak && git stash list >/dev/null; git diff --stat docs/api/astarte_pairing_api.yaml
docs/api/astarte_pairing_api.yaml | 18 ++++++++++++++++--
 1 file changed, 16 insertions(+), 2 deletions(-)

$ git show HEAD:docs/api/astarte_pairing_api.yaml > docs/api/astarte_pairing_api.yaml && go test ./internal/swagger/ -run TestPairingDeviceIDEncodingDocumented 2>&1 | tail -15; cp /tmp/opencode/pairing.bak docs/api/astarte_pairing_api.yaml && git diff --stat docs/api/astarte_pairing_api.yaml
--- FAIL: TestPairingDeviceIDEncodingDocumented (0.00s)
    swagger_test.go:531: DeviceID parameter description does not say "22-character unpadded"
    swagger_test.go:531: DeviceID parameter description does not say "base64url"
    swagger_test.go:535: DeviceID parameter still describes the ID as plain base64
    swagger_test.go:542: RegisterRequest.hw_id description does not say "22-character unpadded"
    swagger_test.go:542: RegisterRequest.hw_id description does not say "base64url"
    swagger_test.go:550: RegisterRequest.hw_id is missing line "          minLength: 22"
    swagger_test.go:550: RegisterRequest.hw_id is missing line "          maxLength: 22"
    swagger_test.go:562: RegisterRequest.hw_id carries no pattern
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.020s
FAIL
 docs/api/astarte_pairing_api.yaml | 18 ++++++++++++++++--
 1 file changed, 16 insertions(+), 2 deletions(-)

$ cat .gitignore 2>/dev/null; echo "--- tracked site files:"; git ls-files docs/site | head -5; git ls-files docs/site | wc -l; echo "--- mkdocs:"; command -v mkdocs || echo "no mkdocs"
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
--- tracked site files:
docs/site/appengine-api.md
docs/site/architecture.md
docs/site/compatibility.md
docs/site/configuration-reference.md
docs/site/contributing.md
23
--- mkdocs:
no mkdocs

$ cat docs/requirements.txt; command -v pip3 pip python3; timeout 60 pip3 install -r docs/requirements.txt 2>&1 | tail -5
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10
timeout: failed to run command ‘pip3’: No such file or directory

$ command -v python3 python pipx uv 2>&1; ls /usr/bin | grep -i -E "^(python|pip)" | head
/usr/bin/python3
pipewire
pipewire-aes67
pipewire-avb
python3
python3-config
python3.13
python3.13-config

$ python3 -c "import yaml, sys; print('pyyaml', yaml.__version__)" 2>&1; python3 -m pip --version 2>&1 | tail -2; rg -n "yaml" go.mod
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml, sys; print('pyyaml', yaml.__version__)
    ^^^^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'
/usr/bin/python3: No module named pip
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

$ timeout 90 make -C docs build 2>&1 | tail -5
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'

$ rg -c "kin-openapi|libopenapi" go.sum 2>/dev/null || echo "no openapi validator dep"; git status --porcelain
no openapi validator dep
 M docs/api/astarte_pairing_api.yaml
 M internal/swagger/swagger_test.go

← Write /tmp/opencode/probe/yamlcheck/main.go
Wrote file successfully.

$ GOFLAGS=-mod=mod go run ./yamlcheck 2>&1 | tail -20
astarte_appengine_api.yaml                           openapi=3.0.3 paths=17 ok
astarte_housekeeping_api.yaml                        openapi=3.0.3 paths=2 ok
astarte_pairing_api.yaml                             openapi=3.0.3 paths=7 ok
astarte_realm_management_api.yaml                    openapi=3.0.3 paths=12 ok
astrate_native_api.yaml                              openapi=3.0.3 paths=12 ok
hw_id minLength=22 maxLength=22 pattern="^[A-Za-z0-9_-]{22}$"
hw_id description first line: "Device hardware ID: a 128-bit value in the 22-character unpa"

$ command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./internal/swagger/ || echo "golangci-lint not installed; skipped"
golangci-lint not installed; skipped

$ go vet ./... 2>&1 | tail -5
(no output)

$ go test -race ./... 2>&1 | grep -v "^ok\|no test files" | tail -20; echo "exit=$?"
FAIL	github.com/astrate-platform/astrate/internal/store	0.137s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.050s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/testutil	0.087s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/astarteapi	0.063s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/deviceid	0.062s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/interfaceschema	0.066s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/pkg/payload	0.046s
FAIL
exit=0

$ go test ./... 2>&1 | grep -v "^ok\|no test files" | tail -20; echo "gate-exit=${PIPESTATUS[0]}"
/bin/sh: 1: Bad substitution

$ go test ./... > /tmp/opencode/gate.txt 2>&1; echo "exit=$?"; grep -v "^ok\|no test files" /tmp/opencode/gate.txt | tail -20
exit=0

$ gofmt -l . | grep -v '^docs/site/' ; echo "gofmt-clean"; git status --porcelain; git diff --stat
gofmt-clean
 M docs/api/astarte_pairing_api.yaml
 M internal/swagger/swagger_test.go
 docs/api/astarte_pairing_api.yaml | 18 ++++++++-
 internal/swagger/swagger_test.go  | 81 +++++++++++++++++++++++++++++++++++++++
 2 files changed, 97 insertions(+), 2 deletions(-)

$ tail -c 400 .mule/todo.md | tail -2 | cut -c1-200
("its data is kept", internal/pairing/service.go:227-229). State the real semantics: the device becomes registrable again, its interfaces, datastream data and group memberships are retained, and a se

$ cat >> .mule/todo.md <<'EOF'
- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `base64.RawURLEncoding.Strict()` (pkg/deviceid/deviceid.go:37), so of the 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with `ErrInvalid` while `..._lA` is accepted, and the pairing spec's `hw_id` pattern `^[A-Za-z0-9_-]{22}$` (added by docs-sync-pairing-deviceid-base64url) deliberately does not encode the rule because a pattern is a client-side contract and upstream's leniency is unmeasured. The docstring justifies strictness as byte-identical round-tripping, which is a reason, not a parity claim, and docs/COMPATIBILITY.md records no entry for it. Measure upstream (Elixir `Base.url_decode64!(padding: false)`, which does not check those bits) on the Legion Go, then either record the divergence in docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.
EOF
git status --porcelain
M .mule/todo.md
 M docs/api/astarte_pairing_api.yaml
 M internal/swagger/swagger_test.go

$ git diff docs/api/astarte_pairing_api.yaml internal/swagger/swagger_test.go
diff --git a/docs/api/astarte_pairing_api.yaml b/docs/api/astarte_pairing_api.yaml
index 8eb5867..b07c004 100644
--- a/docs/api/astarte_pairing_api.yaml
+++ b/docs/api/astarte_pairing_api.yaml
@@ -357,7 +357,11 @@ components:
       name: deviceID
       in: path
       required: true
-      description: The device hardware ID (base64-encoded 128-bit).
+      description: |
+        The device hardware ID: a 128-bit value in the 22-character unpadded
+        base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and `_`,
+        with no `=` padding. The standard-alphabet characters `+` and `/`, and
+        a padded 23-character spelling of the same value, are not accepted.
       schema:
         type: string
 
@@ -370,7 +374,17 @@ components:
       properties:
         hw_id:
           type: string
-          description: Device hardware ID (base64-encoded 128-bit).
+          minLength: 22
+          maxLength: 22
+          pattern: '^[A-Za-z0-9_-]{22}$'
+          description: |
+            Device hardware ID: a 128-bit value in the 22-character unpadded
+            base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and
+            `_`, with no `=` padding. The standard-alphabet characters `+` and
+            `/`, and a padded 23-character spelling of the same value, are
+            rejected with 422 and the
+            `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}`
+            body.
         initial_payload_format:
           type: string
           description: Optional Astrate extension specifying the device's initial payload format.
diff --git a/internal/swagger/swagger_test.go b/internal/swagger/swagger_test.go
index 9aebf6f..4da5545 100644
--- a/internal/swagger/swagger_test.go
+++ b/internal/swagger/swagger_test.go
@@ -2,6 +2,7 @@ package swagger
 
 import (
 	"embed"
+	"fmt"
 	"io"
 	"io/fs"
 	"net/http"
@@ -9,11 +10,13 @@ import (
 	"reflect"
 	"regexp"
 	"sort"
+	"strconv"
 	"strings"
 	"testing"
 
 	docs "github.com/astrate-platform/astrate/docs"
 	"github.com/astrate-platform/astrate/pkg/astarteapi"
+	"github.com/astrate-platform/astrate/pkg/deviceid"
 )
 
 func TestMount(t *testing.T) {
@@ -502,6 +505,84 @@ func TestPairingErrorDetailExamples(t *testing.T) {
 	}
 }
 
+// TestPairingDeviceIDEncodingDocumented guards that the pairing spec states the
+// exact wire form of a device identifier — the 22-character unpadded base64url
+// string deviceid.Parse accepts — instead of the "base64-encoded 128-bit"
+// wording, which read as standard base64 and let a generated client send a
+// well-formed `+`/`/` identifier or a padded 23-character spelling and collect
+// a 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}`
+// (internal/pairing/http.go). The shipped pattern is compiled here and held to
+// the same spellings the parser is, so the two cannot drift apart: the bounds
+// and the description wording are derived from deviceid.EncodedLen rather than
+// written out, so a change to the accepted length fails this test too.
+func TestPairingDeviceIDEncodingDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	wantLen := strconv.Itoa(deviceid.EncodedLen)
+	wantDesc := fmt.Sprintf("%s-character unpadded", wantLen)
+
+	param := strings.Join(componentBlock(t, lines, "    DeviceID:"), "\n")
+	for _, want := range []string{wantDesc, "base64url"} {
+		if !strings.Contains(param, want) {
+			t.Errorf("DeviceID parameter description does not say %q", want)
+		}
+	}
+	if strings.Contains(param, "(base64-encoded 128-bit)") {
+		t.Error("DeviceID parameter still describes the ID as plain base64")
+	}
+
+	block := propertyBlock(t, lines, "hw_id")
+	hwID := strings.Join(block, "\n")
+	for _, want := range []string{wantDesc, "base64url"} {
+		if !strings.Contains(hwID, want) {
+			t.Errorf("RegisterRequest.hw_id description does not say %q", want)
+		}
+	}
+	for _, want := range []string{
+		"          minLength: " + wantLen,
+		"          maxLength: " + wantLen,
+	} {
+		if !containsLine(block, want) {
+			t.Errorf("RegisterRequest.hw_id is missing line %q", want)
+		}
+	}
+
+	pattern := ""
+	for _, l := range block {
+		if v, ok := strings.CutPrefix(l, "          pattern: "); ok {
+			pattern = strings.Trim(v, `'"`)
+			break
+		}
+	}
+	if pattern == "" {
+		t.Fatal("RegisterRequest.hw_id carries no pattern")
+	}
+	re, err := regexp.Compile(pattern)
+	if err != nil {
+		t.Fatalf("hw_id pattern %q does not compile: %v", pattern, err)
+	}
+
+	if !re.MatchString("dT6hS2W9TT6LEnP25ks_lg") {
+		t.Errorf("hw_id pattern %q rejects the spec's own example dT6hS2W9TT6LEnP25ks_lg", pattern)
+	}
+	for _, spelling := range []string{
+		"dT6hS2W9TT6LEnP25ks+lg",  // standard-alphabet '+'
+		"dT6hS2W9TT6LEnP25ks/lg",  // standard-alphabet '/'
+		"dT6hS2W9TT6LEnP25ks_lg=", // padded
+	} {
+		if _, err := deviceid.Parse(spelling); err == nil {
+			t.Errorf("deviceid.Parse accepts %q; this test's premise no longer holds", spelling)
+		}
+		if re.MatchString(spelling) {
+			t.Errorf("hw_id pattern %q accepts %q, which deviceid.Parse rejects", pattern, spelling)
+		}
+	}
+}
+
 // operationBlock returns the lines of the operation with the given operationId,
 // up to the next operation, path, or the components section.
 func operationBlock(t *testing.T, lines []string, operationID string) []string {

```
Done: the pairing spec now states the device identifier's exact wire form (22-character unpadded base64url, alphabet `A-Za-z0-9_-`, no `=` padding) in both the `DeviceID` path parameter and `RegisterRequest.hw_id`, with `minLength: 22` / `maxLength: 22` / `pattern: '^[A-Za-z0-9_-]{22}$'` on `hw_id`, plus a docs-half test that compiles the shipped pattern and holds it to `deviceid.Parse`.
Files: docs/api/astarte_pairing_api.yaml, internal/swagger/swagger_test.go, .mule/todo.md (one appended line)
Verified: go vet ./... -> pass; go test ./... -> pass (exit 0, whole suite; `go test -race` cannot run on this host — ThreadSanitizer "unsupported VMA range: Found 39 - Supported 48", the known Pi limitation, so the gate is the non-race one and no race coverage was obtained for this change); gofmt -l . -> clean; golangci-lint not installed, not run; `make -C docs build` -> fails with `mkdocs: No such file or directory` (no pip on this host, so `pip install -r docs/requirements.txt` is impossible) — fallback run instead: a /tmp Go probe parsing all five docs/api/*.yaml with gopkg.in/yaml.v3, all five parse and decode `hw_id minLength=22 maxLength=22 pattern="^[A-Za-z0-9_-]{22}$"`; TestPairingDeviceIDEncodingDocumented fails on the pre-change spec (8 failures, reverted-and-rerun check)
Unsure: two judgement calls. (1) The pattern encodes alphabet + length only, not `deviceid.Parse`'s strict trailing-bits check — measured: of the 22-character unpadded-base64url strings only those whose last character is in `A,Q,g,w` parse (`..._lx` refused, `..._lA` accepted), so the shipped pattern is deliberately a superset of the parser; putting `[AQgw]` in a pattern would assert an upstream leniency I have not measured. (2) I added the 422 body verbatim into the `hw_id` description as a consequence of the encoding rule, which is slightly beyond "state the encoding and length"; it duplicates the `ValidationError` example at yaml:613-618.
Follow-ups: `.mule/todo.md` + `deviceid-strict-trailing-bits [legion]`: measure whether the strict trailing-bits rejection diverges from upstream and either record it in docs/COMPATIBILITY.md or relax the encoding
```
