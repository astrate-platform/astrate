slug: recipe-docs-sync
verdict: proposed
at:  ee37dc7
ran: 2026-09-28T10:59:15Z on DietPi in 355s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ tail -40 .mule/todo.md; echo "=== for-giulio tail ==="; tail -30 .mule/for-giulio.md 2>&1
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
=== for-giulio tail ===

#78 is no longer parked: zero-touch onboarding is strategic for commercial
viability. Scope frozen on the issue: owner-side TO1/TO2 in our Pairing
service only (last mile, like upstream), reuse fdo-rs for
manufacturing/rendezvous, acceptance = official `astarte-device-fdo-rust`
SDK completes onboarding against Astrate, docs as a first-class deliverable.
When v3.0 is marked DONE, the v4.0 section of `.mule/milestones.md` should be
drafted with this investigation as its first item (issue #78 has the full
verified context).
- **The mule has been idle 17h.** Filed by the dead-man's switch; see journalctl on the Pi.

## 2026-09-25 — the `-tags nats` half of the tree is compiled by no gate

`cmd/astrate/newnats_nats.go`, `internal/engine/forward/nats.go` and
`internal/engine/forward/nats_test.go` are all behind `//go:build nats`, and nothing
builds or tests that tag: not `make build`/`make test`/`make test-integration`/
`make test-e2e` (Makefile:23, 38, 42, 46), not the mule gate, not the lint target.
So the whole NATS trigger-forwarding path is invisible to CI — while
`config.validate` accepts `triggers.forward.kind = "nats"` (internal/config/config.go:339-343)
and the untagged binary then fails at boot with the build-tag message
(cmd/astrate/newnats_default.go). I measured it today: `go build -tags nats ./cmd/...`
and `go vet -tags nats ./cmd/...` both exit 0, so there is nothing broken — just nothing
keeping it unbroken. Worth a `make build-nats` (and the NATS test tier, which needs a
container, so it would be a `[legion]` line) in the gate; not a queue line because the
only test I can imagine for it shells out to `go build`, and the mule's own check strips
the implementation and the guard together, so it would never fail.
Not queued for the same reason: `IdleTimeout` is missing on the HTTP server
(cmd/astrate/main.go:221) and `-healthcheck` ignores the config file
(cmd/astrate/main.go:531) — both real, both yours to call, neither has a test worth
writing. See `.mule/reviews/cmd-astrate-2026-09-25.md`.

$ ls docs/api/ docs/site/ 2>&1; echo ---; ls internal/ 2>&1
docs/api/:
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

docs/site/:
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
housekeeping-api.md
index.md
interface-schema.md
json-payload-profile.md
migration-from-astarte.md
mqtt-protocol.md
observability.md
operations.md
pairing-and-security.md
payload-formats.md
pokemon-agent.md
quickstart.md
realm-management-api.md
swagger-ui
swagger.md
triggers.md
troubleshooting.md
---
appengine
auth
broker
config
engine
flow
flowapi
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil

$ rg -n 'docs-sync' .mule/todo.md | sed 's/\(.\{160\}\).*/\1.../' 
65:`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
80:- [x] docs-sync-rm-policies-delete-422 [auto]: add the missing `422` response to `DELETE /realmmanagement/v1/{realm}/policies/{name}` in docs/api/astarte_rea...
81:- [x] docs-sync-rm-triggers-422-nested-envelope [auto]: document the nested changeset envelope on `POST /realmmanagement/v1/{realm}/triggers` in docs/api/ast...
116:- [x] docs-sync-pairing-health-path [auto]: add the undocumented `GET /pairing/v1/{realm}/health` route to docs/api/astarte_pairing_api.yaml — it exists in ...
117:- [x] docs-sync-pairing-register-404 [auto]: document the missing `404` (DeviceNotFound) on `POST /pairing/v1/{realm}/agent/devices` in docs/api/astarte_pai...
123:- [x] docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` ...
124:- [x] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/as...
127:- [x] docs-sync-appengine-by-alias-endpoints [auto]: add the 7 undocumented by-alias mirror routes to docs/api/astarte_appengine_api.yaml — PATCH /devices-b...
128:- [x] docs-sync-appengine-group-endpoints [auto]: add the 8 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (ret...
129:- [x] docs-sync-appengine-get-group-device [auto]: add the undocumented `GET /groups/{group}/devices/{device}` route to docs/api/astarte_appengine_api.yaml ...
130:- [x] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups...
132:- [x] docs-sync-appengine-group-patch-status [auto]: PATCH /groups/{group}/devices/{device} (internal/appengine/http.go:518-537, patchGroupDevice) shares ap...
133:- [x] docs-sync-appengine-data-422-interface-level [auto]: the interface-level data GETs — GET /devices/{device}/interfaces/{interface}, GET /devices-by-ali...
136:- [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accept...
143:- [x] docs-sync-appengine-device-status-schema [auto]: fix the `DeviceStatus` schema in docs/api/astrate_appengine_api.yaml — the wire emits `id` (internal/...
144:- [!] docs-sync-appengine-data-set-422 [auto]: add the missing `422` response to the six PUT/POST data-set operations in docs/api/astrate_appengine_api.yaml...
145:- [x] docs-sync-appengine-downsample-min [auto]: fix the `downsample_to` constraint in docs/api/astrate_appengine_api.yaml — the DataDownsample description ...
149:- [x] docs-sync-hk-patch-endpoint [auto]: add the undocumented `PATCH /housekeeping/v1/realms/{realm}` route to docs/api/astarte_housekeeping_api.yaml — it ...
150:- [x] docs-sync-hk-retention-field [auto]: add the `datastream_maximum_storage_retention` field to the `Realm` and `RealmCreate` schemas in docs/api/astarte...
154:- [x] docs-sync-native-compat-health-503 [auto]: add the missing `503` response to the three compat-health endpoints — `GET /appengine/health`, `GET /realmm...
155:- [x] docs-sync-pairing-status-enum [auto]: fix the `PairingInfo.status` enum in docs/api/astarte_pairing_api.yaml from `[confirmed, pending, denied, expire...
156:- [x] docs-sync-pairing-version-endpoint [auto]: add the undocumented `GET /pairing/v1/{realm}/version` route to docs/api/astarte_pairing_api.yaml — it exis...
157:- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET ...
158:- [x] docs-sync-rm-delete-interface-status [auto]: fix the `deleteInterface` response in docs/api/astarte_realm_management_api.yaml — the spec documents `42...
159:- [x] docs-sync-rm-mapping-required-encrypted [auto]: add `required` (boolean, datastream-only) and `encrypted` (boolean, datastream-only) to the `Interface...
160:- [x] docs-sync-rm-put-auth-422 [auto]: add the missing `422` (ValidationError) response to `PUT /realmmanagement/v1/{realm}/config/auth` in docs/api/astart...
161:- [!] docs-sync-rm-version-example [auto]: fix the `getVersion` example in docs/api/astarte_realm_management_api.yaml from `data: "1.1.0"` (yaml:557) to `da...
167:- [x] docs-sync-ae-write-405 [auto]: add the missing `405` response to the nine data-write ops in docs/api/astarte_appengine_api.yaml — PUT/POST/DELETE `/de...
168:- [x] docs-sync-ae-write-value-422 [auto]: add the missing `422` response to the six PUT/POST data-write ops in docs/api/astarte_appengine_api.yaml (device-...
169:- [x] docs-sync-ae-post-groups-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups` in docs/api/astarte_appengine_api.yaml — `S...
172:- [x] docs-sync-ae-post-group-devices-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_...
173:- [x] docs-sync-ae-patch-merge-patch-media-type [auto]: fix the three PATCH request bodies in docs/api/astarte_appengine_api.yaml — PATCH /devices/{device} ...
178:- [x] docs-sync-hk-delete-gating-responses [auto]: add the missing `405` and `422` responses to `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte...
179:- [x] docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_housekeeping_api.yaml (yaml:417) — it shows `detail:...
185:- [x] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication ...
188:- [x] docs-sync-pairing-deviceendpoints-dead-404-403 [auto]: remove the unreachable `404` DeviceNotFound response from the three device-facing ops in docs/a...
189:- [x] docs-sync-pairing-version-404-unreachable [auto]: remove the `404` "Unknown realm" response from `GET /pairing/v1/{realm}/version` in docs/api/astarte...
194:- [x] docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Return...
196:- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/re...
201:- [x] docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate...
202:- [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yam...
208:- [!] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine...
216:- [x] docs-sync-rm-put-interface-409 [auto]: add the missing `409` Conflict response to `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` in docs/...
217:- [x] docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1...
218:- [x] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it s...
219:- [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `Requi...
226:- [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /hous...
227:- [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say ...
228:- [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1....
229:- [x] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion G...
230:- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no ...
236:- [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml —...
237:- [x] docs-sync-ae-patch-by-alias-409 [auto]: add the missing `409` response to `patchDeviceByAlias` in docs/api/astarte_appengine_api.yaml (yaml:241-291, c...
238:- [x] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_ap...
239:- [x] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data DELETE operations in docs/api/astarte_appengine_api.yaml (`delet...
240:- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (ya...

$ ls .mule/reviews/ | tail -20; echo ---; rg -ln 'docs-sync' .mule/reviews/ | tail -5
channels-2026-09-18.md
cmd-astrate-2026-09-25.md
config-2026-09-13.md
deviceid-2026-09-18.md
engine-2026-09-06.md
engine-triggers-2026-09-20.md
flow-2026-09-05.md
flow-2026-09-23.md
flowapi-2026-09-10.md
forward-2026-09-26.md
housekeeping-2026-09-15.md
httpx-2026-09-17.md
interfaceschema-2026-09-11.md
observability-2026-09-16.md
pairing-2026-09-09.md
payload-2026-09-12.md
realm-2026-09-09.md
store-2026-09-08.md
store-2026-09-21.md
swagger-2026-09-24.md
---
.mule/reviews/swagger-2026-09-24.md

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*pairing*'
(no output)

$ ls internal/pairing/; echo ---; rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/pairing/
bearersecret_test.go
ca
http.go
http_test.go
ratelimit.go
ratelimit_test.go
remoteip_test.go
service.go
service_test.go
testdata
---
internal/pairing/http.go:81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/pairing/http.go:86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
internal/pairing/http.go:88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
internal/pairing/http.go:90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
internal/pairing/http.go:344:	header := r.Header.Get("Authorization")

$ rg -n '^\s{2}/' docs/api/astarte_pairing_api.yaml
27:  /pairing/v1/{realm}/health:
70:  /pairing/v1/{realm}/version:
97:  /pairing/v1/{realm}/agent/devices:
149:  /pairing/v1/{realm}/agent/devices/{deviceID}:
172:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
224:  /pairing/v1/{realm}/devices/{deviceID}:
258:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:

→ Read internal/pairing/http.go [limit=140]
→ Read internal/pairing/http.go [offset=140, limit=230]
→ Read docs/api/astarte_pairing_api.yaml

$ rg -n 'func \(s \*Service\)|Err[A-Za-z]* =|ErrInhibited|ErrUnauthorized|Version' internal/pairing/service.go | head -60
46:	ErrInvalidHWID = errors.New("pairing: invalid hw_id")
49:	ErrInvalidPayloadFormat = errors.New("pairing: invalid initial_payload_format")
53:	ErrAlreadyRegistered = errors.New("pairing: device already registered")
57:	ErrRegistrationLimitReached = errors.New("pairing: device registration limit reached")
58:	// ErrUnauthorized is the uniform device-authentication failure: unknown
61:	ErrUnauthorized = errors.New("pairing: unauthorized")
62:	// ErrInhibited reports a device blocked by credentials_inhibited (403).
63:	ErrInhibited = errors.New("pairing: credentials request inhibited")
65:	ErrInvalidCSR = errors.New("pairing: invalid CSR")
86:	// DefaultVersion is reported by the info endpoint when Config.Version
88:	DefaultVersion = "0.1.0-astrate"
114:	// Version is reported by the info endpoint; empty selects
115:	// DefaultVersion.
116:	Version string
141:	if cfg.Version == "" {
142:		cfg.Version = DefaultVersion
174:func (s *Service) Register(ctx context.Context, realmName, hwID, initialFormat string) (string, error) {
230:func (s *Service) Unregister(ctx context.Context, realmName, deviceIDStr string) error {
243:// device by credentials secret (uniform ErrUnauthorized on any mismatch),
247:func (s *Service) Credentials(ctx context.Context, realmName, deviceIDStr, secret, csrPEM string, ip netip.Addr) (string, error) {
253:		return "", ErrInhibited
275:	Version   string
286:func (s *Service) Info(ctx context.Context, realmName, deviceIDStr, secret string) (*Info, error) {
302:		Version:   s.cfg.Version,
324:func (s *Service) VerifyCredentials(ctx context.Context, realmName, deviceIDStr, secret, clientCrtPEM string) (*VerifyResult, error) {
330:		return nil, ErrInhibited
369:func (s *Service) Health(ctx context.Context, realm string) error {
378:// device, wrong secret) return the same ErrUnauthorized after a comparable
380:func (s *Service) authenticateDevice(ctx context.Context, realmName, deviceIDStr, secret string) (*store.Realm, *store.Device, error) {
384:		return nil, nil, ErrUnauthorized
389:		return nil, nil, ErrUnauthorized
397:		return nil, nil, ErrUnauthorized
405:		return nil, nil, ErrUnauthorized
408:		return nil, nil, ErrUnauthorized
414:func (s *Service) loadCA(realm *store.Realm) (*ca.CA, error) {

$ rg -n 'func ' internal/pairing/service.go
140:func New(st Store, sealer *store.KeySealer, cfg Config) *Service {
156:func ProvisionCA(realmName string, lifetime time.Duration, sealer *store.KeySealer) (certPEM string, sealedKey []byte, err error) {
174:func (s *Service) Register(ctx context.Context, realmName, hwID, initialFormat string) (string, error) {
230:func (s *Service) Unregister(ctx context.Context, realmName, deviceIDStr string) error {
247:func (s *Service) Credentials(ctx context.Context, realmName, deviceIDStr, secret, csrPEM string, ip netip.Addr) (string, error) {
286:func (s *Service) Info(ctx context.Context, realmName, deviceIDStr, secret string) (*Info, error) {
324:func (s *Service) VerifyCredentials(ctx context.Context, realmName, deviceIDStr, secret, clientCrtPEM string) (*VerifyResult, error) {
369:func (s *Service) Health(ctx context.Context, realm string) error {
380:func (s *Service) authenticateDevice(ctx context.Context, realmName, deviceIDStr, secret string) (*store.Realm, *store.Device, error) {
414:func (s *Service) loadCA(realm *store.Realm) (*ca.CA, error) {
428:func burnBcrypt(secret string) {
434:func generateSecret() (string, error) {

→ Read internal/pairing/service.go [offset=170, limit=160]

$ rg -n 'pairing' docs/site/configuration-reference.md | head -40; echo "=== ENV ==="; rg -o '\bASTRATE_[A-Z_]+' -N internal/ cmd/ | sort -u | head -80
36:| `addr` | string | `":8080"` | `ASTRATE_HTTP_ADDR` | No | Bind address for the single REST listener (pairing, realm management, housekeeping, appengine, native endpoints). |
51:| `advertised_url` | string | `""` | `ASTRATE_MQTT_ADVERTISED_URL` | No | Broker URL handed to devices by the pairing info endpoint. Empty derives `"mqtts://<addr>"`; set when devices reach the broker by another host. |
66:## `[pairing]` — Credential issuance & rate limits
=== ENV ===
cmd/astrate/main.go:ASTRATE_HTTP_ADDR
cmd/astrate/main.go:ASTRATE_MASTER_KEY
cmd/astrate/main.go:ASTRATE_MASTER_KEY_FILE
cmd/astrate/main_test.go:ASTRATE_MASTER_KEY
cmd/astrate/main_test.go:ASTRATE_TEST_DSN
cmd/astrate/sealer_test.go:ASTRATE_MASTER_KEY
internal/config/config.example.toml:ASTRATE_DATABASE_DSN
internal/config/config.example.toml:ASTRATE_ENGINE_SHARDS
internal/config/config.example.toml:ASTRATE_HTTP_ADDR
internal/config/config.example.toml:ASTRATE_LOG_LEVEL
internal/config/config.example.toml:ASTRATE_MASTER_KEY
internal/config/config.example.toml:ASTRATE_MASTER_KEY_FILE
internal/config/config.example.toml:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/config/config.go:ASTRATE_DATABASE_DSN
internal/config/config.go:ASTRATE_ENGINE_SHARDS
internal/config/config.go:ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
internal/config/config.go:ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
internal/config/config.go:ASTRATE_HTTP_ADDR
internal/config/config.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
internal/config/config.go:ASTRATE_HTTP_TLS_CERT_FILE
internal/config/config.go:ASTRATE_HTTP_TLS_KEY_FILE
internal/config/config.go:ASTRATE_LOG_FORMAT
internal/config/config.go:ASTRATE_LOG_LEVEL
internal/config/config.go:ASTRATE_MASTER_KEY
internal/config/config.go:ASTRATE_MASTER_KEY_FILE
internal/config/config.go:ASTRATE_MQTT_ADDR
internal/config/config.go:ASTRATE_MQTT_ADVERTISED_URL
internal/config/config.go:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/config/config.go:ASTRATE_MQTT_SESSION_STORE_PATH
internal/config/config.go:ASTRATE_MQTT_TLS_CERT_FILE
internal/config/config.go:ASTRATE_MQTT_TLS_KEY_FILE
internal/config/config.go:ASTRATE_REALM_JWT_PUBLIC_KEY
internal/config/config.go:ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
internal/config/config.go:ASTRATE_REALM_NAME
internal/config/config.go:ASTRATE_SECURITY_MASTER_KEY_FILE
internal/config/config_test.go:ASTRATE_DATABASE_DSN
internal/config/config_test.go:ASTRATE_ENGINE_SHARDS
internal/config/config_test.go:ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
internal/config/config_test.go:ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
internal/config/config_test.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
internal/config/config_test.go:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/engine/engine_integration_test.go:ASTRATE_TEST_DSN
internal/flow/blocks/container/block.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/container/docker.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/info.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/schema.go:ASTRATE_FLOW_CONFIG
internal/store/crypto.go:ASTRATE_MASTER_KEY
internal/store/crypto.go:ASTRATE_MASTER_KEY_FILE
internal/store/store_test.go:ASTRATE_TEST_DSN
internal/testutil/pg.go:ASTRATE_TEST_DSN

$ rg -n 'APICompatVersion|1\.2\.2|Version' internal/config/config.go | head -30; echo "=== version handler ==="; rg -rn 'version' --glob '!*_test.go' internal/astrate* internal/engine/compat* 2>/dev/null | head; rg -ln 'compat' internal/ | head -20
=== version handler ===
internal/flow/blocks/sort.go
internal/pairing/service.go
internal/flowapi/service.go
internal/realm/http.go
internal/realm/http_test.go
internal/realm/service.go
internal/store/interfaces.go
internal/engine/capabilities.go
internal/auth/middleware.go
internal/appengine/channels/protocol.go

$ rg -n 'pairing|Pairing' internal/config/config.go | head -40
48:	Pairing      PairingConfig      `toml:"pairing"`
77:	// AdvertisedURL is the broker URL handed to devices by the pairing info
98:// PairingConfig tunes credential issuance and its rate limits (§4.3, §4.5).
99:type PairingConfig struct {
192:		Pairing: PairingConfig{

→ Read docs/site/configuration-reference.md [offset=60, limit=40]
→ Read internal/config/config.go [limit=30, offset=96]
→ Read internal/config/config.go [limit=25, offset=188]

$ rg -n 'APICompatVersion|CompatVersion|"1\.2\.' --glob '!*_test.go' internal/ cmd/ | head -20; echo "=== main pairing wiring ==="; rg -n 'pairing\.|Pairing|Version' cmd/astrate/main.go | head -30
cmd/astrate/main.go:437:		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(realm.APICompatVersion)))
cmd/astrate/main.go:438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
internal/realm/http.go:119:// feature UI on it — see APICompatVersion).
internal/realm/http.go:121:	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
internal/realm/service.go:580:// APICompatVersion is the upstream Realm Management API level Astrate
internal/realm/service.go:588:const APICompatVersion = "1.2.2"
=== main pairing wiring ===
63:	showVersion := flag.Bool("version", false, "print the version and exit")
67:	if *showVersion {
141:	// Pairing service is built here (not in mountAPIs) so the flow runtime
147:	pairer := pairing.New(st, sealer, pairing.Config{
149:		CertTTL:           cfg.Pairing.CertTTL.Std(),
150:		EnforceLatestCert: cfg.Pairing.EnforceLatestCert,
151:		Version:           version,
152:		BcryptCost:        cfg.Pairing.BcryptCost,
173:			if errors.Is(err, pairing.ErrAlreadyRegistered) || errors.Is(err, store.ErrDeviceAlreadyConfirmed) {
294:		EnforceLatestCert: cfg.Pairing.EnforceLatestCert,
353:func mountAPIs(cfg config.Config, st *store.Store, e *engine.Engine, b *broker.Broker, sealer *store.KeySealer, metrics *observability.Metrics, flowSvc *flowapi.Service, pairer *pairing.Service, log *slog.Logger) (http.Handler, *housekeeping.Service, error) {
357:	pairing.NewAPI(pairer, mw, pairing.APIConfig{
358:		RegisterRate:     cfg.Pairing.RegisterRate,
359:		RegisterBurst:    cfg.Pairing.RegisterBurst,
360:		CredentialsRate:  cfg.Pairing.CredentialsRate,
361:		CredentialsBurst: cfg.Pairing.CredentialsBurst,
393:	// Pairing's is public — both measured on upstream 1.2.0
397:		observability.MountVersionCompat(mux, svc, version)
399:	mountRealmVersion(mux, mw)
429:// mountRealmVersion registers the realm-scoped version endpoints observed on
435:func mountRealmVersion(mux *http.ServeMux, mw *auth.Middleware) {
437:		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(realm.APICompatVersion)))
438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))

$ rg -n '^var version|version =|version string|ldflags' cmd/astrate/*.go Makefile 2>/dev/null | head -20
Makefile:23:	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' ./...
Makefile:26:		CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(DIST)/astrate ./cmd/astrate; \
cmd/astrate/main.go:52:// -ldflags "-X main.version=vX.Y.Z").
cmd/astrate/main.go:53:var version = "0.1.0-dev"

$ rg -n '0\.1\.0-test|0\.1\.0-dev|0\.1\.0-astrate' docs/ --glob '!site/DESIGN.md' --glob '!site/ROADMAP.md' | head -20; echo "=== native spec version examples ==="; rg -n -A6 'data: "1\.2\.2"' docs/api/astrate_native_api.yaml | head -30
docs/site/api/astarte_pairing_api.yaml:252:                  version: "0.1.0-test"
docs/site/api/astrate_native_api.yaml:203:                data: "0.1.0-dev"
docs/site/api/astrate_native_api.yaml:223:                data: "0.1.0-dev"
docs/site/api/astrate_native_api.yaml:243:                data: "0.1.0-dev"
docs/site/api/astrate_native_api.yaml:263:                data: "0.1.0-dev"
docs/api/astarte_pairing_api.yaml:252:                  version: "0.1.0-test"
docs/api/astrate_native_api.yaml:203:                data: "0.1.0-dev"
docs/api/astrate_native_api.yaml:223:                data: "0.1.0-dev"
docs/api/astrate_native_api.yaml:243:                data: "0.1.0-dev"
docs/api/astrate_native_api.yaml:263:                data: "0.1.0-dev"
=== native spec version examples ===

$ ls -la docs/site/api/ docs/site/swagger-ui/ 2>&1 | head -20; echo "=== docs/Makefile ==="; cat docs/Makefile 2>/dev/null | head -40
docs/site/api/:
total 168
drwxr-xr-x 2 root root  4096 Sep 27 13:32 .
drwxr-xr-x 4 root root  4096 Sep 27 13:32 ..
-rw-r--r-- 1 root root 59526 Sep 27 13:32 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 15547 Sep 27 13:32 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 19506 Sep 27 13:32 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 45387 Sep 27 13:32 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 14977 Sep 27 13:32 astrate_native_api.yaml

docs/site/swagger-ui/:
total 16
drwxr-xr-x 2 root root 4096 Sep 27 13:32 .
drwxr-xr-x 4 root root 4096 Sep 27 13:32 ..
-rw-r--r-- 1 root root 4185 Sep 27 13:32 index.html
=== docs/Makefile ===
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

→ Read internal/pairing/service.go [limit=115, offset=329]

$ rg -n 'Cause[A-Za-z]* *=' internal/pairing/*.go; echo "=== swagger tests for pairing ==="; rg -n 'func Test' internal/swagger/swagger_test.go
internal/pairing/service.go:72:	CauseExpired = "EXPIRED"
internal/pairing/service.go:75:	CauseInvalid = "INVALID"
internal/pairing/service.go:78:	CauseRevoked = "REVOKED"
internal/pairing/service.go:342:		res.Cause = CauseExpired
internal/pairing/service.go:345:		res.Cause = CauseInvalid
internal/pairing/service.go:352:			res.Cause = CauseInvalid
internal/pairing/service.go:356:			res.Cause = CauseRevoked
=== swagger tests for pairing ===
18:func TestMount(t *testing.T) {
77:func TestMountSubPanicsOnBrokenFS(t *testing.T) {
90:func TestSpecs(t *testing.T) {
123:func TestRealmManagement403(t *testing.T) {
134:func TestAppEngine403(t *testing.T) {
198:func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
241:func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
316:func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {
376:func TestAppEngineAddGroupDevice422Documented(t *testing.T) {
420:func TestAppEngineDataDelete400Documented(t *testing.T) {

→ Read pkg/astarteapi/envelope.go

$ rg -n 'detail: (Bad Request|Internal Server Error|Bad request|Internal server error)' docs/api/*.yaml
docs/api/astrate_native_api.yaml:466:              detail: Internal Server Error
docs/api/astarte_housekeeping_api.yaml:362:            detail: Bad Request
docs/api/astarte_housekeeping_api.yaml:385:              detail: Bad Request
docs/api/astarte_housekeeping_api.yaml:477:              detail: Internal Server Error
docs/api/astarte_appengine_api.yaml:1732:              detail: Bad Request
docs/api/astarte_appengine_api.yaml:1792:              detail: Internal Server Error
docs/api/astarte_pairing_api.yaml:556:              detail: Bad Request
docs/api/astarte_pairing_api.yaml:638:              detail: Internal Server Error
docs/api/astarte_realm_management_api.yaml:1265:              detail: Bad Request
docs/api/astarte_realm_management_api.yaml:1348:              detail: Internal Server Error

$ ls pkg/astarteapi/testdata/ 2>/dev/null; rg -n 'Bad request|Internal server error' pkg/astarteapi/*_test.go | head -20
data_array.json
data_links.json
data_links_last.json
data_metadata.json
data_nil_metadata.json
data_null.json
data_object.json
data_pem.json
error_bad_request.json
error_custom_detail.json
error_device_not_found.json
error_fields.json
error_fields_multikey.json
error_forbidden.json
error_internal.json
error_not_found.json
error_raw_multikey.json
error_raw_nested.json
error_unauthorized.json

$ for f in error_bad_request.json error_internal.json error_not_found.json error_device_not_found.json error_forbidden.json error_unauthorized.json; do printf '%s: ' "$f"; cat pkg/astarteapi/testdata/$f; echo; done
error_bad_request.json: {"errors":{"detail":"Bad request"}}
error_internal.json: {"errors":{"detail":"Internal server error"}}
error_not_found.json: {"errors":{"detail":"Not Found"}}
error_device_not_found.json: {"errors":{"detail":"Device not found"}}
error_forbidden.json: {"errors":{"detail":"Forbidden"}}
error_unauthorized.json: {"errors":{"detail":"Unauthorized"}}

$ rg -n 'PayloadFormat|payload_format' --glob '!*_test.go' internal/store/*.go internal/pairing/*.go | head -20; echo "=== bson/json consts ==="; rg -rn '"bson"|"json"' --glob '!*_test.go' internal/payload/ 2>/dev/null | head -10
internal/store/devices.go:58:	PayloadFormatHint        string
internal/store/devices.go:64:	total_received_msgs, total_received_bytes, payload_format_hint`
internal/store/devices.go:406:// SetPayloadFormatHint flips the device's preferred outbound payload format
internal/store/devices.go:408:func (s *Store) SetPayloadFormatHint(ctx context.Context, realmID int16, id deviceid.ID, hint string) error {
internal/store/devices.go:413:		UPDATE devices SET payload_format_hint = $3 WHERE realm_id = $1 AND id = $2`,
internal/store/devices.go:435:		&d.PayloadFormatHint); err != nil {
internal/pairing/service.go:34:	SetPayloadFormatHint(ctx context.Context, realmID int16, id deviceid.ID, hint string) error
internal/pairing/service.go:47:	// ErrInvalidPayloadFormat reports an initial_payload_format outside
internal/pairing/service.go:49:	ErrInvalidPayloadFormat = errors.New("pairing: invalid initial_payload_format")
internal/pairing/service.go:173:// the Astrate initial_payload_format extension ("", "bson" or "json").
internal/pairing/service.go:180:		return "", fmt.Errorf("%w: %q", ErrInvalidPayloadFormat, initialFormat)
internal/pairing/service.go:217:		if err := s.st.SetPayloadFormatHint(ctx, realm.ID, id, initialFormat); err != nil {
internal/pairing/http.go:97:// initial_payload_format extension (docs/DESIGN.md §3.5.4; upstream-shaped
internal/pairing/http.go:101:	InitialPayloadFormat string `json:"initial_payload_format"`
internal/pairing/http.go:126:	secret, err := a.svc.Register(r.Context(), r.PathValue("realm"), req.HwID, req.InitialPayloadFormat)
internal/pairing/http.go:310:	case errors.Is(err, ErrInvalidPayloadFormat):
internal/pairing/http.go:312:			map[string][]string{"initial_payload_format": {"is invalid"}})
internal/store/devices_test.go:31:		if d.Status != DeviceStatusRegistered || d.Connected || d.PayloadFormatHint != "bson" {
internal/store/devices_test.go:32:			t.Errorf("defaults: status=%q connected=%v hint=%q", d.Status, d.Connected, d.PayloadFormatHint)
internal/store/devices_test.go:124:		if err := s.SetPayloadFormatHint(ctx, realm.ID, id, "json"); err != nil {
=== bson/json consts ===

$ rg -n -B8 -A25 'func \(s \*Store\) UnregisterDevice' internal/store/devices.go
88-		return fmt.Errorf("%w: device %s", ErrDeviceAlreadyConfirmed, id)
89-	}
90-	return nil
91-}
92-
93-// UnregisterDevice makes a device registrable again without losing its data
94-// (DELETE /agent/devices parity, docs/DESIGN.md §4.4): the credentials
95-// secret and certificate trail are cleared, the row and all stored data stay.
96:func (s *Store) UnregisterDevice(ctx context.Context, realmID int16, id deviceid.ID) error {
97-	tag, err := s.pool.Exec(ctx, `
98-		UPDATE devices
99-		SET credentials_secret_hash = '', status = 'registered',
100-		    first_credentials_request = NULL, cert_serial = NULL, cert_aki = NULL
101-		WHERE realm_id = $1 AND id = $2`,
102-		realmID, uuidParam(id))
103-	if err != nil {
104-		return fmt.Errorf("store: unregistering device %s: %w", id, err)
105-	}
106-	if tag.RowsAffected() == 0 {
107-		return fmt.Errorf("%w: device %s", ErrNotFound, id)
108-	}
109-	return nil
110-}
111-
112-// GetDevice fetches one device.
113-func (s *Store) GetDevice(ctx context.Context, realmID int16, id deviceid.ID) (*Device, error) {
114-	row := s.pool.QueryRow(ctx,
115-		`SELECT `+deviceColumns+` FROM devices WHERE realm_id = $1 AND id = $2`,
116-		realmID, uuidParam(id))
117-	d, err := scanDevice(row)
118-	if errors.Is(err, pgx.ErrNoRows) {
119-		return nil, fmt.Errorf("%w: device %s", ErrNotFound, id)
120-	}
121-	if err != nil {

$ rg -n 'rate|burst|Rate|Burst|per IP|429' docs/site/pairing-and-security.md | head -20; echo "=== unregister/removal claims ==="; rg -n -i 'unregister|removes|deleted|retain' docs/site/pairing-and-security.md docs/site/operations.md | head -20
3:Astrate reproduces Astarte's three credential planes exactly, replacing the CFSSL sidecar with an embedded per-realm CA.
74:Agent                Astrate(pairing)          Device                 Astrate(broker)
95:## Rate limiting
=== unregister/removal claims ===

→ Read docs/site/pairing-and-security.md [limit=40, offset=90]

$ command -v mkdocs >/dev/null 2>&1 && echo "mkdocs present" || echo "mkdocs MISSING"; rg -n 'api/|swagger' docs/mkdocs.yml | head -20
mkdocs MISSING
63:  - API Explorer: swagger.md

$ cat docs/requirements.txt 2>/dev/null
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10

→ Read internal/swagger/swagger_test.go [limit=35, offset=90]

$ rg -n 'initial_payload_format' docs/site/*.md docs/api/*.yaml | head -20
docs/site/json-payload-profile.md:125:  `initial_payload_format`, so a JSON-only device receives JSON from its very
docs/site/json-payload-profile.md:131:{ "data": { "hw_id": "<22-char base64url id>", "initial_payload_format": "json" } }
docs/site/json-payload-profile.md:139:1. Register (optionally with `initial_payload_format: "json"`), obtain the
docs/site/DESIGN.md:589:  "initial_payload_format": "json"}}` — additive, ignored-by-upstream-shaped).
docs/site/DESIGN.md:605:2. **JSON payloads and `initial_payload_format`** are Astrate extensions — pure supersets;
docs/site/DESIGN.md:704:- Optional Astrate extension field `initial_payload_format` (§3.5.4).
docs/site/pairing-and-security.md:44:- Optional extension: `"initial_payload_format": "json"` for JSON-profile devices.
docs/site/ROADMAP.md:270:| 4.3 | `internal/pairing/service.go` | Flows A–C business logic over `store`: register (secret gen 32 B → 44-char base64, bcrypt store, show-once; re-register rotates iff no credentials request yet else 422; `initial_payload_format` extension §3.5.4), unregister (re-registrable, data kept), credentials (bcrypt compare, uniform error, status flip `registered→confirmed`, stamp first-request + IP, record serial/AKI), info (broker_url + realm `ca_crt`), verify (valid/until vs cause `EXPIRED|INVALID|REVOKED`) | 260 |
docs/site/ROADMAP.md:423:| 9.4 | `test/conformance/atomvm/main.go` | **AtomVM-profile JSON simulator** (§3.5, §6): MQTT 3.1.1 client constrained to the documented JSON profile + zlib control payloads — registration with `initial_payload_format: "json"`, JSON datastreams/properties, asserts server-owned data arrives JSON-encoded |
docs/api/astarte_pairing_api.yaml:374:        initial_payload_format:
docs/site/payload-formats.md:82:{ "data": { "hw_id": "...", "initial_payload_format": "json" } }

$ rg -n -A40 'func Parse' internal/deviceid/*.go | head -60
rg: internal/deviceid/*.go: IO error for operation on internal/deviceid/*.go: No such file or directory (os error 2)

$ rg -rn 'package deviceid' --files-with-matches . 2>/dev/null | head; echo ---; find . -type d -name deviceid
./pkg/deviceid/deviceid_test.go
./pkg/deviceid/deviceid.go
---
./pkg/deviceid

→ Read pkg/deviceid/deviceid.go

$ mkdir -p /tmp/opencode/idprobe && cat > /tmp/opencode/idprobe/main.go <<'EOF'
package main

import (
	"encoding/base64"
	"fmt"
	"strings"
)

var enc = base64.RawURLEncoding.Strict()

func main() {
	for _, s := range []string{"dT6hS2W9TT6LEnP25ks_lg", "dT6hS2W9TT6LEnP25ks+lg", "dT6hS2W9TT6LEnP25ks/lg", "dT6hS2W9TT6LEnP25ksLg=="} {
		b := make([]byte, 16)
		n, err := enc.Decode(b, []byte(s))
		fmt.Printf("%-26q len=%2d n=%d err=%v\n", s, len(s), n, err)
	}
	fmt.Println("hyphen/underscore accepted:", strings.ContainsAny("a_b-c", "_-"))
}
EOF
cd /tmp/opencode/idprobe && (test -f go.mod || go mod init idprobe >/dev/null 2>&1) && go run . 
"dT6hS2W9TT6LEnP25ks_lg"   len=22 n=16 err=<nil>
"dT6hS2W9TT6LEnP25ks+lg"   len=22 n=12 err=illegal base64 data at input byte 19
"dT6hS2W9TT6LEnP25ks/lg"   len=22 n=12 err=illegal base64 data at input byte 19
"dT6hS2W9TT6LEnP25ksLg=="  len=23 n=15 err=illegal base64 data at input byte 21
hyphen/underscore accepted: true

$ rg -n 'base64url|base64url-encoded|unpadded' docs/api/*.yaml | head -20
(no output)

$ rg -n -B3 -A3 'base64-encoded 128-bit|128-bit' docs/api/*.yaml | head -60
docs/api/astarte_pairing_api.yaml-357-      name: deviceID
docs/api/astarte_pairing_api.yaml-358-      in: path
docs/api/astarte_pairing_api.yaml-359-      required: true
docs/api/astarte_pairing_api.yaml:360:      description: The device hardware ID (base64-encoded 128-bit).
docs/api/astarte_pairing_api.yaml-361-      schema:
docs/api/astarte_pairing_api.yaml-362-        type: string
docs/api/astarte_pairing_api.yaml-363-
--
docs/api/astarte_pairing_api.yaml-370-      properties:
docs/api/astarte_pairing_api.yaml-371-        hw_id:
docs/api/astarte_pairing_api.yaml-372-          type: string
docs/api/astarte_pairing_api.yaml:373:          description: Device hardware ID (base64-encoded 128-bit).
docs/api/astarte_pairing_api.yaml-374-        initial_payload_format:
docs/api/astarte_pairing_api.yaml-375-          type: string
docs/api/astarte_pairing_api.yaml-376-          description: Optional Astrate extension specifying the device's initial payload format.

$ rg -n -A6 '^\s{4}Device:|^    DeviceID:|^    Device:' docs/api/astarte_appengine_api.yaml docs/api/astarte_realm_management_api.yaml | head -30
docs/api/astarte_appengine_api.yaml:1413:    DeviceID:
docs/api/astarte_appengine_api.yaml-1414-      name: device
docs/api/astarte_appengine_api.yaml-1415-      in: path
docs/api/astarte_appengine_api.yaml-1416-      required: true
docs/api/astarte_appengine_api.yaml-1417-      description: The device hardware ID.
docs/api/astarte_appengine_api.yaml-1418-      schema:
docs/api/astarte_appengine_api.yaml-1419-        type: string
--
docs/api/astarte_realm_management_api.yaml:887:    DeviceID:
docs/api/astarte_realm_management_api.yaml-888-      name: device
docs/api/astarte_realm_management_api.yaml-889-      in: path
docs/api/astarte_realm_management_api.yaml-890-      required: true
docs/api/astarte_realm_management_api.yaml-891-      description: The device hardware ID.
docs/api/astarte_realm_management_api.yaml-892-      schema:
docs/api/astarte_realm_management_api.yaml-893-        type: string

$ tail -c 300 .mule/todo.md | od -c | tail -5
0000360 342 200 224       B   L   O   C   K   E   D   :       T   I   M
0000400   E   O   U   T       a   f   t   e   r       1   2   0   0   s
0000420     342 200 224       t   a   s   k       t   o   o       b   i
0000440   g   ,       s   p   l   i   t       i   t  \n
0000454

$ cat >> .mule/todo.md <<'EOF'
- [ ] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:556) and `InternalServerError` shows `detail: Internal Server Error` (yaml:638), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37, golden fixture pkg/astarteapi/testdata/error_bad_request.json) and `"Internal server error"` (envelope.go:58, testdata/error_internal.json), both measured on upstream 1.2.0 (envelope.go:33-36, 55-57). The other five detail examples on this file (`Not Found`, `Device not found`, `Unauthorized`, `Forbidden`, `Too Many Requests`) are already correct — change only those two. The identical wrong pair is in the other four specs (appengine 1732/1792, housekeeping 362/385/477, realm_management 1265/1348, native 466); leave those for a run of their surface. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
- [ ] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pairing_api.yaml, which shows `version: "0.1.0-test"` (yaml:252) — a value no build can produce. The field is filled from `pairing.Config.Version` (internal/pairing/service.go:302), wired from `main.version` (cmd/astrate/main.go:151), whose default is `"0.1.0-dev"` (main.go:53) and which is otherwise set only via `-ldflags "-X main.version=..."`. Use `"0.1.0-dev"`, matching the four compat-version examples docs-sync-native-version-value already corrected to that (docs/api/astrate_native_api.yaml:203, 223, 243, 263). Do NOT document the service's own fallback `DefaultVersion = "0.1.0-astrate"` (service.go:86-88): it is unreachable from `main`, which always passes a non-empty value. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` path parameter (docs/api/astarte_pairing_api.yaml:360) and `RegisterRequest.hw_id` (yaml:373) — but `deviceid.Parse` accepts only the 22-character **unpadded base64url** form, in strict mode, and rejects the standard alphabet's `+` and `/` and any `=` padding with an error wrapping `ErrInvalid` (pkg/deviceid/deviceid.go:26, 37, 39-60; upstream parity cited at 36). Measured: with the same `encoding` the parser uses, `dT6hS2W9TT6LEnP25ks+lg` and `.../lg` fail at byte 19 and the padded 23-character form fails at byte 21, so a client generated from the current spec can send a well-formed standard-base64 ID and get 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}` (internal/pairing/http.go:307-309). State the exact encoding and the 22-character length in both descriptions and add a `pattern`/`minLength`/`maxLength` to `hw_id`. The existing `hw_id: dT6hS2W9TT6LEnP25ks_lg` example (yaml:123) is valid and stays. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in docs/api/astarte_pairing_api.yaml (yaml:374-376), an unconstrained string described only as "Optional Astrate extension specifying the device's initial payload format" — so a generated client can send any string, while `Service.Register` accepts only `"bson"` or `"json"` and rejects everything else with 422 `{"errors":{"initial_payload_format":["is invalid"]}}` (internal/pairing/service.go:179-181 -> internal/pairing/http.go:310-312; the allowed set is stated at service.go:173). Add `enum: [bson, json]` plus a note that absent/empty means unset, and add the three missing `ValidationError` examples while you are in the component: `initial_payload_format: ["is invalid"]` (http.go:310-312), `csr: ["can't be blank"]` (http.go:173-177) and `client_crt: ["can't be blank"]` (http.go:253-257) — the component currently shows only the four hw_id/error_name rows (yaml:594-618) and the `csr`/`client_crt` blank-field cases are what a device actually hits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` description in docs/api/astarte_pairing_api.yaml (yaml:153), "Removes a device from the realm. Returns 204 on success.", which reads as a device-and-data deletion. The row and every byte of its stored data survive: `store.UnregisterDevice` only clears the credential trail (`credentials_secret_hash = ''`, `first_credentials_request = NULL`, `cert_serial = NULL`, `cert_aki = NULL`) and flips `status` back to `'registered'` (internal/store/devices.go:93-110), and the service docstring says so outright ("its data is kept", internal/pairing/service.go:227-229). State the real semantics: the device becomes registrable again, its interfaces, datastream data and group memberships are retained, and a second DELETE is a 204 no-op (the row is still there, so RowsAffected is 1, devices.go:106-108) — not the documented 404. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
EOF
tail -6 .mule/todo.md | cut -c1-90
- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appeng
- [ ] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detai
- [ ] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response exam
- [ ] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device ident
- [ ] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.ini
- [ ] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` descript

$ git status --porcelain; echo "--- diff stat ---"; git diff --stat
M .mule/todo.md
--- diff stat ---
 .mule/todo.md | 5 +++++
 1 file changed, 5 insertions(+)

Surface run: **pairing** (the least-recently-audited of the five). All six registered routes, all seven documented paths, and every operation's status set and response fields were checked against `internal/pairing/http.go` and `internal/pairing/service.go` — the route and status-code coverage is already accurate, so the five proposals are all field/description/example-level.

Verified while working, not guessed:
- `pkg/astarteapi/testdata/error_bad_request.json` and `error_internal.json` pin the exact wire bytes (`"Bad request"`, `"Internal server error"`) against the spec's `"Bad Request"` / `"Internal Server Error"`.
- The `getDeviceInfo` version example traced `service.go:302` → `main.go:151` → `main.go:53`.
- The device-ID encoding claim was measured with a throwaway program using the same `base64.RawURLEncoding.Strict()` the parser uses: `+`, `/` and `=` padding all fail; the spec's existing example is valid.
- All seven `[pairing]` config keys in `docs/site/configuration-reference.md` match `internal/config/config.go:99-107, 192-197` exactly (the divergent-looking constants in `internal/pairing/http.go:23-27` are only zero-value fallbacks that `main.go:357-361` never triggers), and no `docs/site/` page contradicts the code — so nothing went to `.mule/for-giulio.md`.

The wrong `Bad Request` / `Internal Server Error` example pair is verbatim in all five specs; each proposal names the other four file:line pairs but scopes itself to pairing, per one-surface-per-run.

```
Done: Ran the docs-sync recipe on the pairing surface and appended 5 verified docs/api/*.yaml proposals to .mule/todo.md
Files: .mule/todo.md (5 lines appended)
Verified: git diff --stat -> 1 file changed, 5 insertions; no code touched, no git mutation, no gh run (the recipe does not ask for one)
Unsure: whether "0.1.0-dev" or a neutral placeholder is the better example value for the info endpoint — a release build shows the ldflags-injected version instead; I chose the default because that is what the native spec was already corrected to
Follow-ups: the four non-pairing specs still carry the same wrong error-detail example pair (appengine 1732/1792, housekeeping 362/385/477, realm_management 1265/1348, native 466) — deliberately not queued, one surface per run
```
