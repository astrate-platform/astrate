slug: recipe-docs-sync
verdict: proposed
at:  94b0d28
ran: 2026-10-02T18:17:12Z on DietPi in 530s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls -la /root/astrate-mule/.mule/ /root/astrate-mule/docs/ /root/astrate-mule/docs/api/ 2>&1
/root/astrate-mule/.mule/:
total 10200
drwxr-xr-x  7 root root     4096 Oct  2 19:50 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       13 Oct  2 20:08 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  2 19:39 .heartbeat
-rw-r--r--  1 root root        0 Oct  2 20:08 .last-output
-rw-r--r--  1 root root        0 Oct  2 20:09 .lock
-rw-r--r--  1 root root       10 Oct  2 20:08 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 10129148 Oct  2 20:09 cron.log
drwxr-xr-x  2 root root     4096 Oct  2 19:50 failed
-rw-r--r--  1 root root    97814 Oct  2 13:44 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    12236 Oct  2 19:50 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  2 13:43 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   117207 Oct  2 19:50 todo.md

/root/astrate-mule/docs/:
total 184
drwxr-xr-x  6 root root  4096 Sep  4 20:46 .
drwxr-xr-x 17 root root  4096 Sep 25 19:41 ..
-rw-r--r--  1 root root 23369 Sep 25 22:05 COMPATIBILITY.md
-rw-r--r--  1 root root 50906 Sep  4 20:46 DESIGN.md
-rw-r--r--  1 root root  6699 Jul 27 05:10 JSON-PAYLOAD-PROFILE.md
-rw-r--r--  1 root root   426 Sep  4 20:46 Makefile
-rw-r--r--  1 root root  6196 Jul 27 05:40 OPERATIONS.md
-rw-r--r--  1 root root 45928 Sep  4 20:46 ROADMAP.md
-rw-r--r--  1 root root  1336 Sep  4 20:46 UPSTREAM-EXPERIMENTAL.md
drwxr-xr-x  2 root root  4096 Sep  4 18:55 api
-rw-r--r--  1 root root   435 Jul 27 22:18 embed.go
drwxr-xr-x  2 root root  4096 Sep  4 20:46 handoff
-rw-r--r--  1 root root  1763 Sep  4 20:46 mkdocs.yml
-rw-r--r--  1 root root    40 Jul 27 05:10 requirements.txt
drwxr-xr-x  4 root root  4096 Sep 28 19:57 site
drwxr-xr-x  2 root root  4096 Jul 27 05:40 swagger-ui

/root/astrate-mule/docs/api/:
total 172
drwxr-xr-x 2 root root  4096 Sep  4 18:55 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 59526 Sep 27 13:31 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 15547 Sep 25 21:48 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 21724 Sep 28 19:59 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 45387 Sep 25 22:07 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 14977 Sep 20 21:25 astrate_native_api.yaml

$ ls /root/astrate-mule/internal/ 2>&1
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

$ cat /root/astrate-mule/.mule/.rotation; echo "---"; tail -60 /root/astrate-mule/.mule/todo.md
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_0fdce5bbd001YPI0t8nJwbGp7V

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
- [x] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` path parameter (docs/api/astarte_pairing_api.yaml:360) and `RegisterRequest.hw_id` (yaml:373) — but `deviceid.Parse` accepts only the 22-character **unpadded base64url** form, in strict mode, and rejects the standard alphabet's `+` and `/` and any `=` padding with an error wrapping `ErrInvalid` (pkg/deviceid/deviceid.go:26, 37, 39-60; upstream parity cited at 36). Measured: with the same `encoding` the parser uses, `dT6hS2W9TT6LEnP25ks+lg` and `.../lg` fail at byte 19 and the padded 23-character form fails at byte 21, so a client generated from the current spec can send a well-formed standard-base64 ID and get 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}` (internal/pairing/http.go:307-309). State the exact encoding and the 22-character length in both descriptions and add a `pattern`/`minLength`/`maxLength` to `hw_id`. The existing `hw_id: dT6hS2W9TT6LEnP25ks_lg` example (yaml:123) is valid and stays. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in docs/api/astarte_pairing_api.yaml (yaml:374-376), an unconstrained string described only as "Optional Astrate extension specifying the device's initial payload format" — so a generated client can send any string, while `Service.Register` accepts only `"bson"` or `"json"` and rejects everything else with 422 `{"errors":{"initial_payload_format":["is invalid"]}}` (internal/pairing/service.go:179-181 -> internal/pairing/http.go:310-312; the allowed set is stated at service.go:173). Add `enum: [bson, json]` plus a note that absent/empty means unset, and add the three missing `ValidationError` examples while you are in the component: `initial_payload_format: ["is invalid"]` (http.go:310-312), `csr: ["can't be blank"]` (http.go:173-177) and `client_crt: ["can't be blank"]` (http.go:253-257) — the component currently shows only the four hw_id/error_name rows (yaml:594-618) and the `csr`/`client_crt` blank-field cases are what a device actually hits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` description in docs/api/astarte_pairing_api.yaml (yaml:153), "Removes a device from the realm. Returns 204 on success.", which reads as a device-and-data deletion. The row and every byte of its stored data survive: `store.UnregisterDevice` only clears the credential trail (`credentials_secret_hash = ''`, `first_credentials_request = NULL`, `cert_serial = NULL`, `cert_aki = NULL`) and flips `status` back to `'registered'` (internal/store/devices.go:93-110), and the service docstring says so outright ("its data is kept", internal/pairing/service.go:227-229). State the real semantics: the device becomes registrable again, its interfaces, datastream data and group memberships are retained, and a second DELETE is a 204 no-op (the row is still there, so RowsAffected is 1, devices.go:106-108) — not the documented 404. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `base64.RawURLEncoding.Strict()` (pkg/deviceid/deviceid.go:37), so of the 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with `ErrInvalid` while `..._lA` is accepted, and the pairing spec's `hw_id` pattern `^[A-Za-z0-9_-]{22}$` (added by docs-sync-pairing-deviceid-base64url) deliberately does not encode the rule because a pattern is a client-side contract and upstream's leniency is unmeasured. The docstring justifies strictness as byte-identical round-tripping, which is a reason, not a parity claim, and docs/COMPATIBILITY.md records no entry for it. Measure upstream (Elixir `Base.url_decode64!(padding: false)`, which does not check those bits) on the Legion Go, then either record the divergence in docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.
- [x] fdo-rc6-scope-delta-for-giulio [auto]: upstream v1.4.0-rc.6 (2026-09-30, 27 commits, 154 files, almost all FDO) grows the onboarding surface issue #78 is scoped against, and none of it is in #78's body (written against upstream 1.3 / rc.5): new `POST /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}/to0` re-runs TO0 and returns the refreshed rendezvous expiry (`{"data":{"expiry":…}}`, ownership_voucher.ex:63-73; 409 `device_already_onboarded` via fallback_controller.ex:120-125); voucher deletion, which revokes the rendezvous registration first and *keeps* the voucher if that fails (`:rendezvous_revocation_failed` → 500) and refuses another realm's voucher through `ensure_voucher_in_realm` (ownership_voucher.ex:40-60 — now expressible because the voucher row carries its realm); `device_id` accepted on voucher upload; the device registered immediately on upload rather than mid-TO2 (`add_unconfirmed_credentials`, device.ex:63-85); single-entry x5chain certs accepted; invalid vouchers no longer stored; vouchers and to2-sessions moved from the per-realm keyspace to the global `astarte` one with `fdo_guid` added to the device row (housekeeping/realms/queries.ex:784, 1129-1133); and RM device deletion now cascading to the voucher (device_removal/core.ex). Write the delta into `.mule/for-giulio.md` as a scope update for #78 and say which of these a client could actually observe — do not edit or comment on the issue (GitHub state is read-only for the mule) and do not implement anything.
- [x] compat-note-v1.4-rc.6 [auto]: record in `.mule/for-giulio.md` that v1.4.0-rc.6 is wire-inert for Astrate, as a delta on the open v1.3.2/v1.3.4 wording proposal rather than a re-derivation. Measured, the whole non-FDO delta is inert: DUP RPC-server reliability is OTP supervisor/GC tuning; the Dashboard's new `unknown_status` is computed client-side from `last_connection && !last_disconnection && !isConnected` (models/Device/index.ts:213-215) and both timestamps are already on Astrate's device shape (internal/appengine/service.go:118), so a device whose disconnection was never recorded reads as "Unknown status" rather than "Disconnected"; the `on_device_registered` dispatch moved from `Agent.register_device` to `Engine.register_device` (pairing/engine.ex:157-165, 186-193) so it also covers the FDO path, but the Agent path is byte-for-byte the same rule and Astrate already fires it on every successful registration (internal/pairing/service.go:221-223); the `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"); and `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma. So the proposal's "v1.4.0 is still rc.5-only" reference becomes rc.6, and its "FDO authentication (pairing, disabled by default)" line understates the 1.4 line, which ships a full ownership-voucher surface Astrate does not emulate (#78, parked). Do not edit docs/COMPATIBILITY.md and do not restate the standing proposal.

- [x] container-stop-deadline [auto]: make `cliInstance.Stop` (internal/flow/blocks/container/docker.go:167-175) derive its 15s bound whenever the passed context carries **no deadline**, not only when it is nil — both cleanup callers pass `context.Background()` (docker.go:122, the "Best-effort cleanup so we do not leave orphans on mapping failure" path, and block.go:121, the not-ready cleanup in `New`), so a wedged docker daemon leaves an unbounded `docker rm -f` inside `exec.CommandContext` holding flow instantiation, while `Block.Stop` (block.go:310-311) does the same job correctly with an explicit 15s timeout. Measured with the package's own `Run` hook: the `rm` invocation arrives with `ctx.Deadline()` unset on both paths. Add a case asserting `inst.Stop(context.Background())` reaches the injected `Run` with a deadline set (it does not today) and that the hostPort-failure path still issues `rm -f <id>` (docker.go:120-124, never executed by any test today). Container-free.
- [!] container-timeout-bounds [auto]: bound `timeout_ms` and `ready_timeout_ms` in `parseConfig` (internal/flow/blocks/container/block.go:205-219) — both are only checked `> 0` and then multiplied unchecked (`time.Duration(n) * time.Millisecond`), which wraps: measured `ready_timeout_ms: 18446744073709` gives `-551.616µs`, so `New` fails with "container: not ready after wait" without ever issuing a probe (the wrapped negative deadline fails the loop guard at httpbridge.go:64), while `ready_timeout_ms: 9223372036854` gives 2562047h47m16s and `New` never returns (still polling `/healthz` after 3s in the probe). The two consumers of a non-positive duration also disagree — `Bridge.timeout()` (httpbridge.go:39-44) treats it as "use the 5s default" while net/http treats a non-positive `Client.Timeout` as no deadline at all (measured: a client with `Timeout: -551µs` returned err=nil after a 3s handler). Reject a value that would overflow or exceed a documented maximum, naming the key and the value in the error, and add table cases asserting `New` refuses both (it accepts them today). Nothing bounds these on the wire either: internal/flow/blocks/schema.go:19 advertises both as bare `"type":"integer"`. Container-free. — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [x] container-parseconfig-rules-test [auto]: pin the four `parseConfig` rules that no test touches (internal/flow/blocks/container/block.go:188-219 — all `count=0` in `go test -coverprofile`, which is 71.2% for the package): `config` must be a JSON object (192-193), `port` an integer in 1–65535 including the 0 / 65536 / non-integer rejections (197-202), and `timeout_ms` and `ready_timeout_ms` positive (205-210, 215-217) — beside the existing `TestConstructor_RequiresImage` (block_test.go:99), which is the only config case the suite has. The port rule is the one with teeth: it is what keeps `127.0.0.1::99999` out of the `docker run` argv (docker.go:93). Container-free, pure.
- [!] container-response-cap-test [auto]: pin `Bridge.RoundTrip`'s 1 MiB response cap (internal/flow/blocks/container/httpbridge.go:133-140 — `io.LimitReader(resp.Body, maxBody+1)` plus the explicit rejection; lines 136, 138.16 and 138.40 are all uncovered), the only bound between a tenant-supplied container and process memory on the message path: assert a handler answering exactly `maxBody+1` bytes is rejected with the "exceeds" error and that exactly `maxBody` bytes still round-trip, driving a `Bridge{MaxBodyBytes: N}` so the fixture stays small. Also worth one line in the same change: the cap is absent from the operator-facing catalog entry (internal/flow/blocks/info.go:116) and from the schema, so an author whose container emits 2 MiB learns about it only from a runtime error. Container-free. — BLOCKED: lint failed: internal/flow/blocks/container/block_test.go:500:2: redefines-builtin-id: redefinition of the built-in function cap (rev
- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.

$ cd /root/astrate-mule && rg -n 'docs-api|docs-sync|docs/api' .mule/todo.md | tail -40
157:- [x] docs-sync-native-compat-version-endpoints [auto]: add the four undocumented unauthenticated compat version endpoints — `GET /appengine/version`, `GET /realmmanagement/version`, `GET /pairing/version`, `GET /housekeeping/version` — to docs/api/astrate_native_api.yaml, answering 200 `{"data":"<APICompatVersion>"}`. They are registered for all four services via `observability.MountVersionCompat` (cmd/astrate/main.go:393-395, compat.go:49-52) but no spec documents them; the native spec already covers their sibling compat-health paths (`/{service}/health`). The realm-scoped `/realmmanagement/v1/{realm}/version` is covered in the RM spec (yaml:532); the appengine/pairing realm-scoped twins (main.go:396-398) are out of this spec's scope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
158:- [x] docs-sync-rm-delete-interface-status [auto]: fix the `deleteInterface` response in docs/api/astarte_realm_management_api.yaml — the spec documents `422` with examples "Interface major version is not 0, can't be deleted" / "Cannot delete an interface that is used by a device introspection" (yaml:255-272), but the handler returns `403` with "Interface can't be deleted" / "Interface can't be deleted since it's currently used" (internal/realm/http.go:387-391, `writeInterfaceError`; both are 403 `StatusForbidden`). The spec's detail strings were lifted from the `writeError` branch (http.go:414-419) which DELETE /interfaces does not use. Fix status to `403` and update the two examples to match the actual wire messages. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
159:- [x] docs-sync-rm-mapping-required-encrypted [auto]: add `required` (boolean, datastream-only) and `encrypted` (boolean, datastream-only) to the `InterfaceMapping` schema in docs/api/astarte_realm_management_api.yaml — the parser accepts and persists both on datastream mappings (internal/interfaceschema/parse.go:410-414), rejects them on properties mappings (parse.go:395-398), the detailed listing emits them (`detailed.go:69`), and the compatibility check compares them (interfaceschema/compat.go:121-123); the schema (yaml:820-862) currently has `additionalProperties: false` and neither field, so the OpenAPI would mark a valid datastream interface install body as invalid. Also add both fields to the detailed-listing example (yaml:70-85). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
160:- [x] docs-sync-rm-put-auth-422 [auto]: add the missing `422` (ValidationError) response to `PUT /realmmanagement/v1/{realm}/config/auth` in docs/api/astarte_realm_management_api.yaml — `SetAuthKey` returns `ErrValidation` for a blank `jwt_public_key_pem` (internal/realm/service.go:632-640), which `writeError` maps to 422; the spec (yaml:431-464) documents only 204/400/401/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
161:- [!] docs-sync-rm-version-example [auto]: fix the `getVersion` example in docs/api/astarte_realm_management_api.yaml from `data: "1.1.0"` (yaml:557) to `data: "1.2.2"` — `APICompatVersion` is `"1.2.2"` (internal/realm/service.go:588); the native spec's four compat-version examples already use "1.2.2". Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1
167:- [x] docs-sync-ae-write-405 [auto]: add the missing `405` response to the nine data-write ops in docs/api/astarte_appengine_api.yaml — PUT/POST/DELETE `/devices/{device}/interfaces/{interface}/{path}`, the three `/devices-by-alias/{alias}` twins, and the three `/groups/{group}/devices/{device}` twins. `writeError` emits `405 Method Not Allowed` with "Cannot write to device owned resource" (`ErrNotServerOwned`) and "Cannot write to read-only resource" (`ErrNotAProperty`) at internal/appengine/http.go:673-676, reachable from putData/putDataByAlias/putDataInGroup (http.go:292-357) and deleteData/deleteDataByAlias/deleteDataInGroup (http.go:311-377); `"405"` currently appears nowhere in the spec (those ops document only 200/204, 400, 401, 404, 500). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
168:- [x] docs-sync-ae-write-value-422 [auto]: add the missing `422` response to the six PUT/POST data-write ops in docs/api/astarte_appengine_api.yaml (device-scoped, by-alias, and in-group at `/interfaces/{interface}/{path}`). A value larger than `maxValueBytes` (64 KiB, internal/appengine/http.go:24) is answered `422` "Value size exceeds size limits" by putData/putDataByAlias/putDataInGroup (http.go:298-301, 327-329, 346-348), but the specs for PUT/POST (yaml:686-695, 726-736, 453-463, 499-509, 1176-1186, 1218-1228) list only 200/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
169:- [x] docs-sync-ae-post-groups-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups` in docs/api/astarte_appengine_api.yaml — `Service.CreateGroup` returns `ErrGroupAlreadyExists` (internal/appengine/service.go:495, mapped at http.go:660-661) which the handler answers as `409` "Group already exists"; the spec (yaml:821-839) documents only 201/400/401/422/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
172:- [x] docs-sync-ae-post-group-devices-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml — re-adding a member the group already holds answers `409` "Device already in group" (`Service.AddGroupDevice` returns `ErrDeviceAlreadyInGroup`, internal/appengine/service.go:636-638, mapped at http.go:662-663), but the spec (yaml:967-977) lists only 201/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
173:- [x] docs-sync-ae-patch-merge-patch-media-type [auto]: fix the three PATCH request bodies in docs/api/astarte_appengine_api.yaml — PATCH /devices/{device} (yaml:176), PATCH /devices-by-alias/{alias} (yaml:260), PATCH /groups/{group}/devices/{device} (yaml:1023) declare `application/json; charset=utf-8`, but all three handlers enforce `Content-Type: application/merge-patch+json` exactly and reject every other value — 500 on the device and group paths (internal/appengine/http.go:209-211, 521-523), 400 on the alias path (http.go:230-232). Change the requestBody media type to `application/merge-patch+json` on all three. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
178:- [x] docs-sync-hk-delete-gating-responses [auto]: add the missing `405` and `422` responses to `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — `writeError` answers 405 "Realm deletion disabled" for `ErrDeletionDisabled` and 422 `{"errors":{"error_name":["connected_devices_present"]}}` for `ErrConnectedDevicesPresent` (internal/housekeeping/http.go:233-249), both reachable from `Service.DeleteRealm` (service.go:248-266), but the spec lists only 204/401/403/404/500 (yaml:198-208). Its sole `405` reference sits on PATCH (yaml:182-183), which can never return it (`UpdateRealm`, service.go:196-216, has no deletion-gating path) — re-point 405 to DELETE (reusing the MethodNotAllowed component) and add 422 with the FieldErrors shape and a connected_devices_present example. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
179:- [x] docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_housekeeping_api.yaml (yaml:417) — it shows `detail: "housekeeping: validation failed: realm name must not be empty"`, but `validationDetail` strips the `housekeeping: validation failed: ` prefix on the wire (internal/housekeeping/http.go:251-258) and the real blank-name message is `"realm_name can't be blank"` (service.go:138). Change the example to `detail: "realm_name can't be blank"`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
185:- [x] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication is "via Authorization header or `?token=` query parameter" (yaml:293-294), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
188:- [x] docs-sync-pairing-deviceendpoints-dead-404-403 [auto]: remove the unreachable `404` DeviceNotFound response from the three device-facing ops in docs/api/astarte_pairing_api.yaml — `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials` (yaml:226-227), `GET /pairing/v1/{realm}/devices/{deviceID}` (yaml:268-269), `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify` (yaml:339-340) — and the `403` Forbidden from the info GET (yaml:266-267). `Service.authenticateDevice` folds every failure (malformed/unknown deviceID, unknown realm, unregistered device, wrong secret) into uniform `ErrUnauthorized` (internal/pairing/service.go:380-411), so `store.ErrNotFound` never reaches `writeServiceError` (http.go:326-327) from those handlers; and `Service.Info` never returns `ErrInhibited` (inhibited devices still read info with 200 + status "inhibited", service.go:286-307, comment at 282-285). Reachable on these three are only 200/201/400/401/422/429/500 plus 403 on credentials and verify (both genuinely raise `ErrInhibited`, service.go:252-254 and 329-331). Keep 404 on the agent endpoints (register yaml:149-150, unregister yaml:176-177) where `GetRealmByName`/`deviceid.Parse` really do return `store.ErrNotFound`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
189:- [x] docs-sync-pairing-version-404-unreachable [auto]: remove the `404` "Unknown realm" response from `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml (yaml:96-104) — the route is served by `observability.VersionHandler` (cmd/astrate/main.go:398, internal/observability/compat.go:38-47), which answers static `200 {"data":"<version>"}` with no realm lookup, so a 404 can never be emitted there. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
194:- [x] docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
196:- [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
201:- [x] docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate_devices_total{realm="test"} 42`, but no such series is ever registered — the registry carries `astrate_broker_sessions` (internal/observability/metrics.go:54-62), `astrate_db_pool_{acquired,idle,total,max}_conns` (metrics.go:65-79), plus the engine/flow/trigger families (`astrate_engine_*` engine/router.go:411-444, `astrate_engine_stream_dropped_total` engine/stream/bus.go:119, `astrate_engine_trigger_*` engine/triggers/actions.go:296-300, `astrate_flow_router_*` flow/router.go:260-276) and the standard `go_*`/`process_*` runtime collectors — so the example advertises a scrape result the wire can never return. Swap it for a real series (e.g. `astrate_broker_sessions 0`), keeping the HELP/TYPE comment style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
202:- [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
208:- [!] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine_api.yaml — `retrieve_metadata` and `downsample_key` are parsed for every data read by `parseQueryOpts` (internal/appengine/http.go:624 `retrieve_metadata`, 633 `downsample_key`) yet appear on none of the read endpoints: device-scoped GET `{interface}` (yaml:583-606) and `{interface}/{path}` (yaml:623-649), by-alias (yaml:324-368, 369-417), and in-group (yaml:1108-1132, 1149-1176) — the parameter block on each lists only since/since_after/to/limit/downsample_to/sort/format/allow_bigintegers/allow_safe_bigintegers (components DataSince..DataAllowSafeBigIntegers, yaml:1360-1439). Add the two `$ref`s (new DataRetrieveMetadata and DataDownsampleKey components, or inline — mirror the existing style) to all six, and update the DataSort-independent `sort` default note only if verified. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
216:- [x] docs-sync-rm-put-interface-409 [auto]: add the missing `409` Conflict response to `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` in docs/api/astarte_realm_management_api.yaml (yaml:222-234 lists only 204/400/401/404/422/500) — `updateInterface` answers 409 for every url/body disagreement and minor-upgrade incompatibility: `ErrNameMismatch` (service.go:204 → "Interface name doesn't match the one in the interface json", http.go:367-369), `ErrMajorMismatch` (service.go:207 → http.go:370-372), and the `CheckMinorUpgrade` sentinels `ErrMinorNotIncreased`/`ErrDowngradeNotAllowed`/`ErrMissingEndpoints`/`ErrIncompatibleEndpointChange` (service.go:223-226 → http.go:373-384) plus `store.ErrAlreadyExists` (http.go:385-386) — the exact set `TestRealmManagementErrorCodes` pins as StatusConflict (http_test.go:417-437). Reuse the `Conflict` response component (already used on POST, yaml:131-132) and add name-mismatch/major-mismatch examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
217:- [x] docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` (yaml:231-232) in docs/api/astarte_realm_management_api.yaml as a oneOf — both handlers go through `writeInterfaceError`, which answers three distinct 422 bodies: the flat ErrorDetail (`ErrValidation` via `validationDetail`, http.go:392-393), the nested violations changeset envelope (`writeViolations`, http.go:358-359 + 434-533, e.g. `{"errors":{"description":["should be at most 1000 character(s)"]}}` and the aligned full-length `mappings` array, http_test.go:480-501), and the named FieldErrors envelope for `ErrMaximumDatabaseRetentionExceeded` (`{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, http.go:355-357, http_test.go:282-320); today only the flat ValidationError is referenced. Mirror the createTrigger oneOf pattern (yaml:344-368). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
218:- [x] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "realm: validation failed: interface definition is invalid"`, but `validationDetail` strips the `realm: validation failed: ` prefix on the wire (http.go:562-568) and the parser's real messages are `invalid interface: ...` (pkg/interfaceschema/parse.go:38, wrapped at service.go:163); the example should show the stripped, real message. Same class as the already-fixed docs-sync-hk-validation-example (housekeeping http.go:251-258). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
219:- [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
226:- [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
227:- [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
228:- [x] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
229:- [x] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion Go — the device-delete leg added by docs-sync-rm-delete-device-async-operation-param (register a device, `DELETE /devices/{id}?async_operation=…` answers 204 and the row is gone on both values) compiles clean under `go vet -tags integration ./internal/realm/` but has never been executed, because the Pi has no database; until it runs, TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go is the only assertion of that leg that has actually been verified.
230:- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
236:- [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml — the string `"403"` appears nowhere in the file and `components` defines no `Forbidden` response, yet every route is wrapped by `a.require = mw.RequireRealm(auth.ClaimAppEngine)` (internal/appengine/http.go:35), which answers `403 Forbidden` for a verified a_aea JWT whose grants do not authorize the method+path (`astarteapi.WriteForbidden`, internal/auth/middleware.go:90; `DetailForbidden`, pkg/astarteapi/envelope.go:41; pinned today by internal/appengine/http_test.go:184-185). Add a `Forbidden` response component copied from the realm-management one (docs/api/astarte_realm_management_api.yaml:1277-1285 — `ErrorDetail` schema, example `errors.detail: Forbidden`) and `$ref` it from all 32 operations, which is exactly the parity docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go:123-170) pins; add the same docs-half test for the AppEngine spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
237:- [x] docs-sync-ae-patch-by-alias-409 [auto]: add the missing `409` response to `patchDeviceByAlias` in docs/api/astarte_appengine_api.yaml (yaml:241-291, currently 200/400/401/404/422/500) — the handler resolves the alias then calls the same `applyPatch` as the device- and group-scoped PATCH, and a rename to an alias already in use returns `ErrAliasAlreadyInUse` (internal/appengine/service.go:337) which `writeError` maps to `409 ConflictAliasInUse` (internal/appengine/http.go:650-651); the other two PATCH operations already document that response thanks to docs-sync-appengine-query-params-status and docs-sync-appengine-group-patch-status, so this is the last PATCH of the three with a hole. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
238:- [x] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml (yaml:939-980, currently 201/400/401/404/409/500) — a well-formed body whose `device_id` does not parse is rejected as a field error (`FieldErrors{"device_id": {"is not a valid device id"}}`, internal/appengine/service.go:631-634) and answered 422 with the `FieldErrorsDetail` body (`astarteapi.WriteFieldErrors`, internal/appengine/http.go:647), so `$ref` the existing `ValidationErrors` response (yaml:1707-1716) rather than `BadRequest` — the documented 400 only covers a body `DecodeData` cannot parse (internal/appengine/http.go:488-491). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
239:- [x] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data DELETE operations in docs/api/astarte_appengine_api.yaml (`deleteDataByAlias` yaml:519-547, `deleteData` yaml:756-784, `deleteDataInGroup` yaml:1285-1313 — each currently 204/401/404/405/500) — unsetting a path that matches no endpoint mapping returns `engine.ErrPathNotFound` (`UnsetServerProperty`, internal/engine/serverdata.go:241-244) which `writeError` answers as `400 "Endpoint not found"` (internal/appengine/http.go:679-680), the same 400 the PUT/POST twins already document via `BadRequest`; internal/appengine/writeerror_test.go pins the mapping. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
240:- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 200/401/500) — a `from_token` that does not parse as a device id is refused with `ErrValidation: invalid cursor` (internal/appengine/service.go:150-155) which `writeError` maps to 422 with the flat `{"errors":{"detail":"invalid cursor"}}` body, i.e. the `ErrorDetail` shape and not `FieldErrorsDetail`, so `$ref` the same response the group twin uses for the identical parameter (`ValidationErrors`, yaml:934-935) only after fixing its example, or add a dedicated one. Extend the `from_token` parameter description (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
241:- [x] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:556) and `InternalServerError` shows `detail: Internal Server Error` (yaml:638), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37, golden fixture pkg/astarteapi/testdata/error_bad_request.json) and `"Internal server error"` (envelope.go:58, testdata/error_internal.json), both measured on upstream 1.2.0 (envelope.go:33-36, 55-57). The other five detail examples on this file (`Not Found`, `Device not found`, `Unauthorized`, `Forbidden`, `Too Many Requests`) are already correct — change only those two. The identical wrong pair is in the other four specs (appengine 1732/1792, housekeeping 362/385/477, realm_management 1265/1348, native 466); leave those for a run of their surface. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
242:- [x] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pairing_api.yaml, which shows `version: "0.1.0-test"` (yaml:252) — a value no build can produce. The field is filled from `pairing.Config.Version` (internal/pairing/service.go:302), wired from `main.version` (cmd/astrate/main.go:151), whose default is `"0.1.0-dev"` (main.go:53) and which is otherwise set only via `-ldflags "-X main.version=..."`. Use `"0.1.0-dev"`, matching the four compat-version examples docs-sync-native-version-value already corrected to that (docs/api/astrate_native_api.yaml:203, 223, 243, 263). Do NOT document the service's own fallback `DefaultVersion = "0.1.0-astrate"` (service.go:86-88): it is unreachable from `main`, which always passes a non-empty value. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
243:- [x] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` path parameter (docs/api/astarte_pairing_api.yaml:360) and `RegisterRequest.hw_id` (yaml:373) — but `deviceid.Parse` accepts only the 22-character **unpadded base64url** form, in strict mode, and rejects the standard alphabet's `+` and `/` and any `=` padding with an error wrapping `ErrInvalid` (pkg/deviceid/deviceid.go:26, 37, 39-60; upstream parity cited at 36). Measured: with the same `encoding` the parser uses, `dT6hS2W9TT6LEnP25ks+lg` and `.../lg` fail at byte 19 and the padded 23-character form fails at byte 21, so a client generated from the current spec can send a well-formed standard-base64 ID and get 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}` (internal/pairing/http.go:307-309). State the exact encoding and the 22-character length in both descriptions and add a `pattern`/`minLength`/`maxLength` to `hw_id`. The existing `hw_id: dT6hS2W9TT6LEnP25ks_lg` example (yaml:123) is valid and stays. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
244:- [x] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in docs/api/astarte_pairing_api.yaml (yaml:374-376), an unconstrained string described only as "Optional Astrate extension specifying the device's initial payload format" — so a generated client can send any string, while `Service.Register` accepts only `"bson"` or `"json"` and rejects everything else with 422 `{"errors":{"initial_payload_format":["is invalid"]}}` (internal/pairing/service.go:179-181 -> internal/pairing/http.go:310-312; the allowed set is stated at service.go:173). Add `enum: [bson, json]` plus a note that absent/empty means unset, and add the three missing `ValidationError` examples while you are in the component: `initial_payload_format: ["is invalid"]` (http.go:310-312), `csr: ["can't be blank"]` (http.go:173-177) and `client_crt: ["can't be blank"]` (http.go:253-257) — the component currently shows only the four hw_id/error_name rows (yaml:594-618) and the `csr`/`client_crt` blank-field cases are what a device actually hits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
245:- [x] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` description in docs/api/astarte_pairing_api.yaml (yaml:153), "Removes a device from the realm. Returns 204 on success.", which reads as a device-and-data deletion. The row and every byte of its stored data survive: `store.UnregisterDevice` only clears the credential trail (`credentials_secret_hash = ''`, `first_credentials_request = NULL`, `cert_serial = NULL`, `cert_aki = NULL`) and flips `status` back to `'registered'` (internal/store/devices.go:93-110), and the service docstring says so outright ("its data is kept", internal/pairing/service.go:227-229). State the real semantics: the device becomes registrable again, its interfaces, datastream data and group memberships are retained, and a second DELETE is a 204 no-op (the row is still there, so RowsAffected is 1, devices.go:106-108) — not the documented 404. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
246:- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `base64.RawURLEncoding.Strict()` (pkg/deviceid/deviceid.go:37), so of the 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with `ErrInvalid` while `..._lA` is accepted, and the pairing spec's `hw_id` pattern `^[A-Za-z0-9_-]{22}$` (added by docs-sync-pairing-deviceid-base64url) deliberately does not encode the rule because a pattern is a client-side contract and upstream's leniency is unmeasured. The docstring justifies strictness as byte-identical round-tripping, which is a reason, not a parity claim, and docs/COMPATIBILITY.md records no entry for it. Measure upstream (Elixir `Base.url_decode64!(padding: false)`, which does not check those bits) on the Legion Go, then either record the divergence in docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.

$ cd /root/astrate-mule && rg -n 'docs-sync-hk-|docs-sync-housekeeping' .mule/todo.md | head -30
149:- [x] docs-sync-hk-patch-endpoint [auto]: add the undocumented `PATCH /housekeeping/v1/realms/{realm}` route to docs/api/astarte_housekeeping_api.yaml — it exists in code (internal/housekeeping/http.go:40, handler patchRealm at http.go:163-229, returns 200 on success) but the spec jumps straight from GET to DELETE. Document the request body (`RealmPatch` — same three optional fields as `RealmCreate`: `jwt_public_key_pem`, `device_registration_limit`, `datastream_maximum_storage_retention`, all nullable), the success response, and the error responses: 400 (bad envelope), 401, 403, 404, 405 (deletion disabled), 422 (invalid update parameters for unknown fields, blank jwt_public_key_pem, negative limit/retention, connected devices present), 500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
150:- [x] docs-sync-hk-retention-field [auto]: add the `datastream_maximum_storage_retention` field to the `Realm` and `RealmCreate` schemas in docs/api/astarte_housekeeping_api.yaml — the field is on the wire (realmBody struct, internal/housekeeping/http.go:50, returned by viewBody at http.go:99, accepted in create at http.go:68 and patch at http.go:216) but absent from both schemas (yaml:200-231). Type: integer, format int64, nullable, description matching the device_registration_limit style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
178:- [x] docs-sync-hk-delete-gating-responses [auto]: add the missing `405` and `422` responses to `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — `writeError` answers 405 "Realm deletion disabled" for `ErrDeletionDisabled` and 422 `{"errors":{"error_name":["connected_devices_present"]}}` for `ErrConnectedDevicesPresent` (internal/housekeeping/http.go:233-249), both reachable from `Service.DeleteRealm` (service.go:248-266), but the spec lists only 204/401/403/404/500 (yaml:198-208). Its sole `405` reference sits on PATCH (yaml:182-183), which can never return it (`UpdateRealm`, service.go:196-216, has no deletion-gating path) — re-point 405 to DELETE (reusing the MethodNotAllowed component) and add 422 with the FieldErrors shape and a connected_devices_present example. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
179:- [x] docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_housekeeping_api.yaml (yaml:417) — it shows `detail: "housekeeping: validation failed: realm name must not be empty"`, but `validationDetail` strips the `housekeeping: validation failed: ` prefix on the wire (internal/housekeeping/http.go:251-258) and the real blank-name message is `"realm_name can't be blank"` (service.go:138). Change the example to `detail: "realm_name can't be blank"`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
218:- [x] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "realm: validation failed: interface definition is invalid"`, but `validationDetail` strips the `realm: validation failed: ` prefix on the wire (http.go:562-568) and the parser's real messages are `invalid interface: ...` (pkg/interfaceschema/parse.go:38, wrapped at service.go:163); the example should show the stripped, real message. Same class as the already-fixed docs-sync-hk-validation-example (housekeeping http.go:251-258). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
226:- [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
227:- [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

→ Read .mule/.rotation

$ cd /root/astrate-mule && rg -n '^- \[' .mule/todo.md | tail -12 | cut -c1-120
242:- [x] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pa
243:- [x] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-
244:- [x] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in do
245:- [x] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` description in docs/api/astarte_pa
246:- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailin
247:- [x] fdo-rc6-scope-delta-for-giulio [auto]: upstream v1.4.0-rc.6 (2026-09-30, 27 commits, 154 files, almost all FDO
248:- [x] compat-note-v1.4-rc.6 [auto]: record in `.mule/for-giulio.md` that v1.4.0-rc.6 is wire-inert for Astrate, as a
250:- [x] container-stop-deadline [auto]: make `cliInstance.Stop` (internal/flow/blocks/container/docker.go:167-175) der
251:- [!] container-timeout-bounds [auto]: bound `timeout_ms` and `ready_timeout_ms` in `parseConfig` (internal/flow/blo
252:- [x] container-parseconfig-rules-test [auto]: pin the four `parseConfig` rules that no test touches (internal/flow/
253:- [!] container-response-cap-test [auto]: pin `Bridge.RoundTrip`'s 1 MiB response cap (internal/flow/blocks/containe
254:- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-con

$ cd /root/astrate-mule && rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*realm*'
(no output)

Picking the realm-management surface for this run (pairing was last, appengine the run before).
$ cd /root/astrate-mule && rg -n '^  /' docs/api/astrate_realm_management_api.yaml
rg: docs/api/astrate_realm_management_api.yaml: IO error for operation on docs/api/astrate_realm_management_api.yaml: No such file or directory (os error 2)

$ ls -la /root/astrate-mule/docs/api/; pwd
total 172
drwxr-xr-x 2 root root  4096 Sep  4 18:55 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 59526 Sep 27 13:31 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 15547 Sep 25 21:48 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 21724 Sep 28 19:59 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 45387 Sep 25 22:07 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 14977 Sep 20 21:25 astrate_native_api.yaml
/root/astrate-mule

$ rg -n '^  /' docs/api/astrate_realm_management_api.yaml
rg: docs/api/astrate_realm_management_api.yaml: IO error for operation on docs/api/astrate_realm_management_api.yaml: No such file or directory (os error 2)

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/realm/
(no output)

$ ls -lab /root/astrate-mule/docs/api/; echo "---"; ls -la /root/astrate-mule/internal/realm/ 2>&1 | head -20; echo "---"; stat /root/astrate-mule/docs/api/astrate_realm_management_api.yaml 2>&1 | head
total 172
drwxr-xr-x 2 root root  4096 Sep  4 18:55 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 59526 Sep 27 13:31 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 15547 Sep 25 21:48 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 21724 Sep 28 19:59 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 45387 Sep 25 22:07 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 14977 Sep 20 21:25 astrate_native_api.yaml
---
total 132
drwxr-xr-x  2 root root  4096 Sep  9 21:46 .
drwxr-xr-x 17 root root  4096 Sep  4 20:46 ..
-rw-r--r--  1 root root  4344 Sep 25 22:06 async_operation_test.go
-rw-r--r--  1 root root  8885 Sep  4 20:46 dashboard_compat_test.go
-rw-r--r--  1 root root  2674 Sep  8 18:56 detailed.go
-rw-r--r--  1 root root  1315 Sep  8 18:54 detailed_required_encrypted_test.go
-rw-r--r--  1 root root  2666 Sep  4 20:46 detailed_test.go
-rw-r--r--  1 root root  3765 Sep  9 21:46 helpers_test.go
-rw-r--r--  1 root root 19956 Sep  9 21:49 http.go
-rw-r--r--  1 root root 35006 Sep  4 20:46 http_test.go
-rw-r--r--  1 root root   427 Jul 27 05:10 policy.go
-rw-r--r--  1 root root  3593 Jul 27 05:10 policy_test.go
-rw-r--r--  1 root root 24167 Sep  9 21:49 service.go
---
stat: cannot statx '/root/astrate-mule/docs/api/astrate_realm_management_api.yaml': No such file or directory

$ cd /root/astrate-mule/docs/api && wc -c astarte_realm_management_api.yaml; echo "rc=$?"; python3 -c "import os; print(len(os.listdir('.')))"
45387 astarte_realm_management_api.yaml
rc=0
5

$ cd /root/astrate-mule && ls -la docs/ | head; echo "== readlink"; readlink -f docs/api; readlink -f /root; echo "== direct"; ls docs/api/astrate_realm_management_api.yaml; echo "rc=$?"; ls -la /root/astrate-mule/docs/api/astrate_realm_management_api.yaml; echo "rc2=$?"
total 184
drwxr-xr-x  6 root root  4096 Sep  4 20:46 .
drwxr-xr-x 17 root root  4096 Sep 25 19:41 ..
-rw-r--r--  1 root root 23369 Sep 25 22:05 COMPATIBILITY.md
-rw-r--r--  1 root root 50906 Sep  4 20:46 DESIGN.md
-rw-r--r--  1 root root  6699 Jul 27 05:10 JSON-PAYLOAD-PROFILE.md
-rw-r--r--  1 root root   426 Sep  4 20:46 Makefile
-rw-r--r--  1 root root  6196 Jul 27 05:40 OPERATIONS.md
-rw-r--r--  1 root root 45928 Sep  4 20:46 ROADMAP.md
-rw-r--r--  1 root root  1336 Sep  4 20:46 UPSTREAM-EXPERIMENTAL.md
== readlink
/root/astrate-mule/docs/api
/root
== direct
rc=2
rc2=2
ls: cannot access 'docs/api/astrate_realm_management_api.yaml': No such file or directory
ls: cannot access '/root/astrate-mule/docs/api/astrate_realm_management_api.yaml': No such file or directory

$ cd /root/astrate-mule/docs/api && for f in *; do printf '%s: ' "$f"; if [ -f "$f" ]; then echo "FILE $(wc -c < "$f")"; else echo "NOT-FILE"; fi; done; echo "=== ls -f"; ls -f; echo "=== find"; find . -maxdepth 1
astarte_appengine_api.yaml: FILE 59526
astarte_housekeeping_api.yaml: FILE 15547
astarte_pairing_api.yaml: FILE 21724
astarte_realm_management_api.yaml: FILE 45387
astrate_native_api.yaml: FILE 14977
=== ls -f
..
astarte_realm_management_api.yaml
.
astrate_native_api.yaml
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
=== find
.
./astarte_realm_management_api.yaml
./astrate_native_api.yaml
./astarte_appengine_api.yaml
./astarte_housekeeping_api.yaml
./astarte_pairing_api.yaml

$ rg -n '^  /' astarte_realm_management_api.yaml
27:  /realmmanagement/v1/{realm}/interfaces:
172:  /realmmanagement/v1/{realm}/interfaces/{name}:
203:  /realmmanagement/v1/{realm}/interfaces/{name}/{major}:
352:  /realmmanagement/v1/{realm}/triggers:
449:  /realmmanagement/v1/{realm}/triggers/{name}:
505:  /realmmanagement/v1/{realm}/config/auth:
578:  /realmmanagement/v1/{realm}/config/device_registration_limit:
611:  /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention:
648:  /realmmanagement/v1/{realm}/version:
683:  /realmmanagement/v1/{realm}/devices/{device}:
712:  /realmmanagement/v1/{realm}/policies:
786:  /realmmanagement/v1/{realm}/policies/{name}:

$ rg -n 'Handle|\.(Get|Post|Put|Delete|Patch)\(' http.go | head -60
25:	require func(http.Handler) http.Handler
37:	h := func(f http.HandlerFunc) http.Handler { return a.require(f) }
38:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces", h(a.listInterfaces))
39:	mux.Handle("POST /realmmanagement/v1/{realm}/interfaces", h(a.installInterface))
40:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}", h(a.listInterfaceMajors))
41:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.getInterface))
42:	mux.Handle("PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.updateInterface))
43:	mux.Handle("DELETE /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.deleteInterface))
44:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers", h(a.listTriggers))
45:	mux.Handle("POST /realmmanagement/v1/{realm}/triggers", h(a.createTrigger))
46:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers/{name}", h(a.getTrigger))
47:	mux.Handle("DELETE /realmmanagement/v1/{realm}/triggers/{name}", h(a.deleteTrigger))
48:	mux.Handle("GET /realmmanagement/v1/{realm}/config/auth", h(a.getAuth))
49:	mux.Handle("PUT /realmmanagement/v1/{realm}/config/auth", h(a.putAuth))
50:	mux.Handle("GET /realmmanagement/v1/{realm}/config/device_registration_limit", h(a.getRegistrationLimit))
51:	mux.Handle("GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention", h(a.getDatastreamMaximumStorageRetention))
52:	mux.Handle("GET /realmmanagement/v1/{realm}/version", h(a.getVersion))
53:	mux.Handle("DELETE /realmmanagement/v1/{realm}/devices/{device}", h(a.deleteDevice))
54:	mux.Handle("GET /realmmanagement/v1/{realm}/policies", h(a.listPolicies))
55:	mux.Handle("POST /realmmanagement/v1/{realm}/policies", h(a.createPolicy))
56:	mux.Handle("GET /realmmanagement/v1/{realm}/policies/{name}", h(a.getPolicy))
57:	mux.Handle("DELETE /realmmanagement/v1/{realm}/policies/{name}", h(a.deletePolicy))
150:	if r.URL.Query().Get("detailed") == "true" {

→ Read internal/realm/http.go

$ rg -n '^  [a-zA-Z/{}._-]+:|^      operationId:|^        [0-9]{3}:|^      parameters:|^    [A-Za-z]+:' astarte_realm_management_api.yaml
3:  title: Astarte Realm Management API
4:  description: |
10:  version: "0.0.1"
11:  license:
12:    name: "Apache-2.0 with SECO Exclusion"
13:    url: https://www.apache.org/licenses/LICENSE-2.0
17:    variables:
28:    get:
29:      operationId: listInterfaces
41:      parameters:
95:    post:
96:      operationId: installInterface
104:      parameters:
173:    get:
174:      operationId: listInterfaceMajors
180:      parameters:
204:    get:
205:      operationId: getInterface
211:      parameters:
236:    put:
237:      operationId: updateInterface
245:      parameters:
308:    delete:
309:      operationId: deleteInterface
318:      parameters:
353:    get:
354:      operationId: listTriggers
360:      parameters:
379:    post:
380:      operationId: createTrigger
388:      parameters:
450:    get:
451:      operationId: getTrigger
457:      parameters:
481:    delete:
482:      operationId: deleteTrigger
488:      parameters:
506:    get:
507:      operationId: getAuthConfig
513:      parameters:
539:    put:
540:      operationId: putAuthConfig
546:      parameters:
579:    get:
580:      operationId: getDeviceRegistrationLimit
586:      parameters:
612:    get:
613:      operationId: getDatastreamMaximumStorageRetention
622:      parameters:
649:    get:
650:      operationId: getVersion
658:      parameters:
684:    delete:
685:      operationId: deleteDevice
694:      parameters:
713:    get:
714:      operationId: listPolicies
720:      parameters:
739:    post:
740:      operationId: createPolicy
748:      parameters:
787:    get:
788:      operationId: getPolicy
794:      parameters:
818:    delete:
819:      operationId: deletePolicy
825:      parameters:
844:  securitySchemes:
853:  parameters:
854:    RealmName:
862:    InterfaceName:
870:    MajorVersion:
879:    TriggerName:
887:    DeviceID:
895:    PolicyName:
903:    AsyncOperation:
919:  schemas:
922:    AuthConfig:
932:    InterfaceDefinition:
968:    InterfaceMapping:
1020:    TriggerDefinition:
1038:    SimpleTrigger:
1069:    TriggerAction:
1098:    TriggerPolicyDefinition:
1125:    TriggerPolicyErrorHandler:
1155:    DataEnvelopeStringArray:
1164:    DataEnvelopeIntArray:
1173:    DataEnvelopeInterfaceDetailedArray:
1185:    ErrorDetail:
1196:    TriggerErrorsDetail:
1221:    InterfaceViolationsDetail:
1243:    InterfaceFieldErrorsDetail:
1256:  responses:
1257:    BadRequest:
1267:    Unauthorized:
1277:    Forbidden:
1287:    NotFound:
1297:    DeviceNotFound:
1307:    Conflict:
1330:    ValidationError:
1340:    InternalServerError:

$ python3 - <<'EOF'
import re
lines=open('astarte_realm_management_api.yaml').read().split('\n')
op=None
cur=None
for i,l in enumerate(lines,1):
    m=re.match(r'^      operationId: (\S+)',l)
    if m: op=m.group(1); print(f'\n== {op} (line {i})'); continue
    m=re.match(r'^  (/[^:]*):',l)
    if m: print(f'\n--- PATH {m.group(1)}')
    m=re.match(r'^        (get|post|put|delete|patch|parameters|requestBody|responses):',l)
    if m: cur=m.group(1)
    m=re.match(r'^        (?:description:)?',l)
    m2=re.match(r'^          (\d{3}):',l)
    if m2 and cur=='responses': print(f'   {i}: {m2.group(1)}')
    m3=re.match(r'^      (get|post|put|delete|patch):',l)
    if m3: cur=m3.group(1)
EOF
--- PATH /realmmanagement/v1/{realm}/interfaces

== listInterfaces (line 29)

== installInterface (line 96)

--- PATH /realmmanagement/v1/{realm}/interfaces/{name}

== listInterfaceMajors (line 174)

--- PATH /realmmanagement/v1/{realm}/interfaces/{name}/{major}

== getInterface (line 205)

== updateInterface (line 237)

== deleteInterface (line 309)

--- PATH /realmmanagement/v1/{realm}/triggers

== listTriggers (line 354)

== createTrigger (line 380)

--- PATH /realmmanagement/v1/{realm}/triggers/{name}

== getTrigger (line 451)

== deleteTrigger (line 482)

--- PATH /realmmanagement/v1/{realm}/config/auth

== getAuthConfig (line 507)

== putAuthConfig (line 540)

--- PATH /realmmanagement/v1/{realm}/config/device_registration_limit

== getDeviceRegistrationLimit (line 580)

--- PATH /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention

== getDatastreamMaximumStorageRetention (line 613)

--- PATH /realmmanagement/v1/{realm}/version

== getVersion (line 650)

--- PATH /realmmanagement/v1/{realm}/devices/{device}

== deleteDevice (line 685)

--- PATH /realmmanagement/v1/{realm}/policies

== listPolicies (line 714)

== createPolicy (line 740)

--- PATH /realmmanagement/v1/{realm}/policies/{name}

== getPolicy (line 788)

== deletePolicy (line 819)

→ Read docs/api/astarte_realm_management_api.yaml [offset=27, limit=180]

$ rg -n 'func \(s \*Service\)' service.go
100:func (s *Service) WithDisconnecter(d Disconnecter) *Service {
113:func (s *Service) DeleteDevice(ctx context.Context, realm, deviceID string) error {
139:func (s *Service) realmID(ctx context.Context, realm string) (int16, error) {
156:func (s *Service) InstallInterface(ctx context.Context, realm string, def []byte) (*store.StoredInterface, error) {
191:func (s *Service) UpdateInterface(ctx context.Context, realm, urlName string, urlMajor int, def []byte) (*store.StoredInterface, error) {
238:func (s *Service) checkRetentionCeiling(r *store.Realm, iface *interfaceschema.Interface) error {
255:func (s *Service) DeleteInterface(ctx context.Context, realm, name string, major int) error {
274:func (s *Service) checkNameCollision(ctx context.Context, rid int16, name string) error {
309:func (s *Service) ListInterfaces(ctx context.Context, realm string) ([]string, error) {
335:func (s *Service) ListInterfacesDetailed(ctx context.Context, realm string) ([]json.RawMessage, error) {
364:func (s *Service) ListInterfaceMajors(ctx context.Context, realm, name string) ([]int, error) {
388:func (s *Service) GetInterface(ctx context.Context, realm, name string, major int) (json.RawMessage, error) {
419:func (s *Service) CreateTrigger(ctx context.Context, realm string, def []byte) (*store.Trigger, error) {
460:func (s *Service) GetTrigger(ctx context.Context, realm, name string) (json.RawMessage, error) {
473:func (s *Service) DeleteTrigger(ctx context.Context, realm, name string) error {
486:func (s *Service) ListTriggers(ctx context.Context, realm string) ([]string, error) {
515:func (s *Service) CreatePolicy(ctx context.Context, realm string, def []byte) (*store.TriggerPolicy, error) {
528:func (s *Service) GetPolicy(ctx context.Context, realm, name string) (json.RawMessage, error) {
541:func (s *Service) ListPolicies(ctx context.Context, realm string) ([]string, error) {
559:func (s *Service) DeletePolicy(ctx context.Context, realm, name string) error {
592:func (s *Service) GetDeviceRegistrationLimit(ctx context.Context, realm string) (*int32, error) {
605:func (s *Service) GetDatastreamMaximumStorageRetention(ctx context.Context, realm string) (int64, error) {
622:func (s *Service) GetAuthKey(ctx context.Context, realm string) (string, error) {
632:func (s *Service) SetAuthKey(ctx context.Context, realm, keyPEM string) error {
664:func (s *Service) interfacesChanged(ctx context.Context, rid int16, realm string) {
677:func (s *Service) triggersChanged(ctx context.Context, rid int16, realm string) {

→ Read internal/realm/service.go [offset=305, limit=130]
→ Read internal/auth/middleware.go [limit=70, offset=45]
→ Read internal/realm/service.go [limit=30, offset=135]
→ Read docs/api/astarte_realm_management_api.yaml [limit=150, offset=203]
→ Read docs/api/astarte_realm_management_api.yaml [limit=160, offset=353]
→ Read internal/realm/service.go [limit=130, offset=434]
→ Read internal/realm/service.go [limit=125, offset=562]
→ Read docs/api/astarte_realm_management_api.yaml [limit=345, offset=505]
→ Read docs/api/astarte_realm_management_api.yaml [limit=140, offset=850]
→ Read docs/api/astarte_realm_management_api.yaml [limit=360, offset=989]
→ Read internal/realm/policy.go
→ Read internal/realm/service.go [limit=30, offset=105]

$ rg -n 'Detail[A-Za-z]+ *=' envelope.go
37:	DetailBadRequest = "Bad request"
39:	DetailUnauthorized = "Unauthorized"
41:	DetailForbidden = "Forbidden"
43:	DetailNotFound = "Not Found"
46:	DetailDeviceNotFound = "Device not found"
51:	DetailRouteNotFound = "Not found"
54:	DetailPageNotFound = "Page not found"
58:	DetailInternalServerError = "Internal server error"

$ rg -n 'quality|aggregate|alias|"path"' parse.go | head -30
107:// The omitempty tags only matter on re-encode: after a legacy-alias
116:	// Quality is the upstream legacy alias of ownership (issue #61): accepted
119:	Quality *Ownership `json:"quality,omitempty"`
120:	// Aggregate is the upstream legacy alias of aggregation (issue #61).
121:	Aggregate   *bool         `json:"aggregate,omitempty"`
127:	// aliases records whether any legacy alias field was consumed, so
130:	aliases bool `json:"-"`
136:	// Path is the upstream legacy alias of endpoint (issue #61): accepted in
166:// the canonical re-encoding of the document when any upstream legacy alias
167:// field (interface quality/aggregate, mapping path — issue #61) was consumed
190:	if !raw.aliases {
212:	// Legacy aliases (issue #61): normalise before any rule runs. quality is
214:	// aggregate only without aggregation (both → violation on aggregation).
216:		raw.aliases = true
218:			vc.add("ownership", "ownership and quality are mutually exclusive")
225:		raw.aliases = true
227:			vc.add("aggregation", "aggregation and aggregate are mutually exclusive")
294:	// Interface-level length bounds (#61), after alias normalisation.
307:			// Mapping path alias (issue #61): endpoint wins when both are
310:			raw.aliases = true
605:			return fmt.Errorf("object-aggregated mappings must all have the same depth: %q has %d levels, %q has %d",
610:				return fmt.Errorf("object-aggregated mappings must share the same prefix: %q diverges from %q",
615:			return fmt.Errorf("object-aggregated mapping %q must end in a literal segment", m.Endpoint)
625:// an object-aggregated interface (description/doc may differ).
643:		return fmt.Errorf("object-aggregated mappings must all share the same %s: %q differs from %q",

$ cd /root/astrate-mule && rg -n 'Bad request|Internal server error|capitalis|DeviceID|deviceid|legacy alias|quality|ValidationError.*example|updateInterface|additionalProperties' .mule/todo.md | cut -c1-180
159:- [x] docs-sync-rm-mapping-required-encrypted [auto]: add `required` (boolean, datastream-only) and `encrypted` (boolean, datastream-only) to the `InterfaceMapping` schema in d
179:- [x] docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_housekeeping_api.yaml (yaml:417) — it shows `detail: "housekeeping: va
188:- [x] docs-sync-pairing-deviceendpoints-dead-404-403 [auto]: remove the unreachable `404` DeviceNotFound response from the three device-facing ops in docs/api/astarte_pairing_a
200:- [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a locally-run `elixir -e` with `Base.url_decode64!`) with the dev
216:- [x] docs-sync-rm-put-interface-409 [auto]: add the missing `409` Conflict response to `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` in docs/api/astarte_realm_ma
218:- [x] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "rea
240:- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 
241:- [x] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadRequest` component shows `
243:- [x] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` path parameter (docs/a
244:- [x] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in docs/api/astarte_pairing_api.yaml (yaml:374-376), an unconstra
246:- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `bas

→ Read internal/realm/service.go [limit=50, offset=191]
→ Read pkg/interfaceschema/parse.go [limit=70, offset=100]
→ Read internal/realm/service.go [limit=60, offset=30]

$ rg -n 'func CompilePolicy' -A 40 policy.go 2>/dev/null | head -60; ls
97:func CompilePolicy(def []byte) (*Policy, error) {
98-	var doc struct {
99-		Name          string `json:"name"`
100-		ErrorHandlers []struct {
101-			On       json.RawMessage `json:"on"`
102-			Strategy string          `json:"strategy"`
103-		} `json:"error_handlers"`
104-		MaximumCapacity int  `json:"maximum_capacity"`
105-		PrefetchCount   *int `json:"prefetch_count"`
106-		RetryTimes      *int `json:"retry_times"`
107-		EventTTL        *int `json:"event_ttl"`
108-	}
109-	if err := json.Unmarshal(def, &doc); err != nil {
110-		return nil, fmt.Errorf("triggers: policy does not parse as JSON: %w", err)
111-	}
112-	if !policyNameRE.MatchString(doc.Name) {
113-		return nil, fmt.Errorf("triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]")
114-	}
115-	if len(doc.ErrorHandlers) == 0 {
116-		return nil, fmt.Errorf("triggers: policy must have at least one error handler")
117-	}
118-
119-	handlers := make([]compiledHandler, 0, len(doc.ErrorHandlers))
120-	hasRetry := false
121-	for i, h := range doc.ErrorHandlers {
122-		var strat Strategy
123-		switch h.Strategy {
124-		case "discard":
125-			strat = StrategyDiscard
126-		case "retry":
127-			strat = StrategyRetry
128-			hasRetry = true
129-		default:
130-			return nil, fmt.Errorf("triggers: error handler strategy must be discard or retry")
131-		}
132-		ch, err := compileHandlerOn(h.On)
133-		if err != nil {
134-			return nil, fmt.Errorf("triggers: error handler %d: %w", i+1, err)
135-		}
136-		ch.strategy = strat
137-		handlers = append(handlers, ch)
actions.go
actions_test.go
change_test.go
errorname.go
errorname_test.go
events.go
events_test.go
match.go
match_test.go
policy.go
policy_test.go
testdata
transient.go
transient_test.go
validation.go

$ rg -n 'ErrValidation\s*=' /root/astrate-mule/internal/realm/
/root/astrate-mule/internal/realm/service.go:29:var ErrValidation = errors.New("realm: validation failed")

→ Read pkg/interfaceschema/parse.go [limit=30, offset=25]

$ rg -rn 'invalid interface|missing "interface_name"' /root/astrate-mule/pkg/interfaceschema/ /root/astrate-mule/internal/realm/ | head -20
/root/astrate-mule/pkg/interfaceschema/violations_test.go:		if !strings.HasPrefix(ve.Error(), "n: ") {
/root/astrate-mule/pkg/interfaceschema/parse.go:var ErrInvalid = errors.New("n")
/root/astrate-mule/pkg/interfaceschema/parse.go:		return nil, errors.New(`n`)
/root/astrate-mule/pkg/interfaceschema/parse.go:		return nil, fmt.Errorf("n name %q", name)
/root/astrate-mule/pkg/interfaceschema/violations.go://	n: description: should be at most 1000 character(s); mappings[1].doc: ...
/root/astrate-mule/pkg/interfaceschema/violations.go:	b.WriteString("n: ")
/root/astrate-mule/pkg/interfaceschema/testdata/invalid/manifest.json:  "01-bad-name.json": "n name",

$ cd /root/astrate-mule && rg -n 'invalid interface|missing "interface_name"' pkg/interfaceschema/ internal/realm/ | head -20
pkg/interfaceschema/violations_test.go:228:		if !strings.HasPrefix(ve.Error(), "invalid interface: ") {
pkg/interfaceschema/parse.go:38:var ErrInvalid = errors.New("invalid interface")
pkg/interfaceschema/parse.go:240:		return nil, errors.New(`missing "interface_name"`)
pkg/interfaceschema/parse.go:256:		return nil, fmt.Errorf("invalid interface name %q", name)
pkg/interfaceschema/violations.go:65://	invalid interface: description: should be at most 1000 character(s); mappings[1].doc: ...
pkg/interfaceschema/violations.go:68:	b.WriteString("invalid interface: ")
pkg/interfaceschema/testdata/invalid/manifest.json:2:  "01-bad-name.json": "invalid interface name",

→ Read pkg/interfaceschema/parse.go [limit=30, offset=230]

$ rg -n 'responses/ValidationError' astarte_realm_management_api.yaml; echo "--- Conflict refs"; rg -n 'responses/Conflict' astarte_realm_management_api.yaml; echo "--- NotFound refs"; rg -n 'responses/NotFound' astarte_realm_management_api.yaml
574:          $ref: "#/components/responses/ValidationError"
782:          $ref: "#/components/responses/ValidationError"
839:          $ref: "#/components/responses/ValidationError"
--- Conflict refs
137:          $ref: "#/components/responses/Conflict"
273:          $ref: "#/components/responses/Conflict"
420:          $ref: "#/components/responses/Conflict"
780:          $ref: "#/components/responses/Conflict"
--- NotFound refs
199:          $ref: "#/components/responses/NotFound"
232:          $ref: "#/components/responses/NotFound"
271:          $ref: "#/components/responses/NotFound"
329:          $ref: "#/components/responses/NotFound"
477:          $ref: "#/components/responses/NotFound"
499:          $ref: "#/components/responses/NotFound"
814:          $ref: "#/components/responses/NotFound"
837:          $ref: "#/components/responses/NotFound"

$ cd /root/astrate-mule && rg -n 'async' docs/COMPATIBILITY.md | head -12; echo "=== for-giulio config-key entries"; rg -n 'configuration-reference' .mule/for-giulio.md | head -5; echo "=== todo config tasks"; rg -c 'configuration-reference' .mule/todo.md
247:   an asynchronous multi-service wipe and exposes a transient
335:17. **Always synchronous where upstream 1.4 defaults to asynchronous** —
340:    `?async_operation=false`. Astrate performs every one of these
=== for-giulio config-key entries
79:- **`docs/site/configuration-reference.md` `[housekeeping]` lists 2 of the 4 keys it has, and
152:- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`ASTRATE_FLOW_CONFIG`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
222:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
305:  `docs/site/configuration-reference.md`** (site prose is yours, never edited by the mule).
308:  (lines 270-271 and 288); neither appears in `docs/site/configuration-reference.md`.
=== todo config tasks

→ Read docs/COMPATIBILITY.md [limit=16, offset=335]
→ Read pkg/deviceid/deviceid.go [limit=45, offset=20]

$ cd /root/astrate-mule && ls pkg/astarteapi/testdata/ | head
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

$ cd /root/astrate-mule && head -3 docs/api/astrate_realm_management_api.yaml 2>/dev/null || head -3 /root/astrate-mule/docs/api/astarte_realm_management_api.yaml; ls pkg/astarteapi/testdata/ | rg 'error_'
openapi: "3.0.3"
info:
  title: Astarte Realm Management API
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

$ cd /root/astrate-mule && tail -c 400 .mule/todo.md | cat -A | tail -5 | cut -c1-100
ker ps` M-bM-^FM-^R `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown

→ Read .mule/todo.md [offset=253, limit=2]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -250,5 +250,10 @@
 - [x] container-stop-deadline [auto]: make `cliInstance.Stop` (internal/flow/blocks/container/docker.go:167-175) derive its 15s bound whenever the passed context carries **no deadline**, not only when it is nil — both cleanup callers pass `context.Background()` (docker.go:122, the "Best-effort cleanup so we do not leave orphans on mapping failure" path, and block.go:121, the not-ready cleanup in `New`), so a wedged docker daemon leaves an unbounded `docker rm -f` inside `exec.CommandContext` holding flow instantiation, while `Block.Stop` (block.go:310-311) does the same job correctly with an explicit 15s timeout. Measured with the package's own `Run` hook: the `rm` invocation arrives with `ctx.Deadline()` unset on both paths. Add a case asserting `inst.Stop(context.Background())` reaches the injected `Run` with a deadline set (it does not today) and that the hostPort-failure path still issues `rm -f <id>` (docker.go:120-124, never executed by any test today). Container-free.
 - [!] container-timeout-bounds [auto]: bound `timeout_ms` and `ready_timeout_ms` in `parseConfig` (internal/flow/blocks/container/block.go:205-219) — both are only checked `> 0` and then multiplied unchecked (`time.Duration(n) * time.Millisecond`), which wraps: measured `ready_timeout_ms: 18446744073709` gives `-551.616µs`, so `New` fails with "container: not ready after wait" without ever issuing a probe (the wrapped negative deadline fails the loop guard at httpbridge.go:64), while `ready_timeout_ms: 9223372036854` gives 2562047h47m16s and `New` never returns (still polling `/healthz` after 3s in the probe). The two consumers of a non-positive duration also disagree — `Bridge.timeout()` (httpbridge.go:39-44) treats it as "use the 5s default" while net/http treats a non-positive `Client.Timeout` as no deadline at all (measured: a client with `Timeout: -551µs` returned err=nil after a 3s handler). Reject a value that would overflow or exceed a documented maximum, naming the key and the value in the error, and add table cases asserting `New` refuses both (it accepts them today). Nothing bounds these on the wire either: internal/flow/blocks/schema.go:19 advertises both as bare `"type":"integer"`. Container-free. — BLOCKED: TIMEOUT after 1201s — task too big, split it
 - [x] container-parseconfig-rules-test [auto]: pin the four `parseConfig` rules that no test touches (internal/flow/blocks/container/block.go:188-219 — all `count=0` in `go test -coverprofile`, which is 71.2% for the package): `config` must be a JSON object (192-193), `port` an integer in 1–65535 including the 0 / 65536 / non-integer rejections (197-202), and `timeout_ms` and `ready_timeout_ms` positive (205-210, 215-217) — beside the existing `TestConstructor_RequiresImage` (block_test.go:99), which is the only config case the suite has. The port rule is the one with teeth: it is what keeps `127.0.0.1::99999` out of the `docker run` argv (docker.go:93). Container-free, pure.
 - [!] container-response-cap-test [auto]: pin `Bridge.RoundTrip`'s 1 MiB response cap (internal/flow/blocks/container/httpbridge.go:133-140 — `io.LimitReader(resp.Body, maxBody+1)` plus the explicit rejection; lines 136, 138.16 and 138.40 are all uncovered), the only bound between a tenant-supplied container and process memory on the message path: assert a handler answering exactly `maxBody+1` bytes is rejected with the "exceeds" error and that exactly `maxBody` bytes still round-trip, driving a `Bridge{MaxBodyBytes: N}` so the fixture stays small. Also worth one line in the same change: the cap is absent from the operator-facing catalog entry (internal/flow/blocks/info.go:116) and from the schema, so an author whose container emits 2 MiB learns about it only from a runtime error. Container-free. — BLOCKED: lint failed: internal/flow/blocks/container/block_test.go:500:2: redefines-builtin-id: redefinition of the built-in function cap (rev
+- [ ] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, housekeeping and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
+- [ ] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy alias fields to `InterfaceDefinition` and `InterfaceMapping` in docs/api/astarte_realm_management_api.yaml — both declare `additionalProperties: false` (yaml:935, 972) yet omit interface-level `quality` (alias of `ownership`) and `aggregate` (alias of `aggregation`) and mapping-level `path` (alias of `endpoint`), all three of which the parser accepts and canonicalises: decoded into dedicated fields (pkg/interfaceschema/parse.go:116-122, 136-138), normalised before any rule runs (parse.go:212-236), and re-encoded canonically by `ParseInterfaceCanonical` so the store never sees them (parse.go:165-170, 307-311). A spec-conformant client — or any generated client with this schema baked in — therefore rejects a valid upstream-shaped install/update with a spurious 422 while the service would have accepted and stored it as `ownership`/`aggregation`/`endpoint`. Add the three as documented deprecated aliases (enums mirroring ownership/aggregation, `path` a plain string) or relax `additionalProperties`, and record the three rules the code enforces: `ownership` + `quality` is a violation (parse.go:218) and `aggregation` + `aggregate` is a violation (parse.go:227), while `endpoint` + `path` silently keeps `endpoint` (parse.go:307-311). None of the three may appear in a `required` list. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstrained string described only as "The device hardware ID." — `Service.DeleteDevice` runs the path value through `deviceid.Parse` (internal/realm/service.go:118-121), which accepts exactly 22 characters of unpadded base64url in strict mode and rejects padding, the standard alphabet's `+` and `/`, and non-canonical trailing bits (pkg/deviceid/deviceid.go:24-26, 33-37, 42-60); every rejection is folded into `store.ErrNotFound` (service.go:119-121) and answered `404 {"errors":{"detail":"Device not found"}}` (internal/realm/http.go:107-109), so a client that sends a well-formed standard-base64 or 23-character padded id is told the device does not exist rather than that its id is malformed. State the encoding and the 22-character length and add `pattern`/`minLength`/`maxLength`, exactly as docs-sync-pairing-deviceid-base64url did for the pairing twin. Do NOT encode the strict trailing-bits rule in the pattern: whether it is a deviation worth recording is still open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_realm_management_api.yaml — `data` is declared as a bare `type: object` with no properties (yaml:250-261) although the handler parses it with the same `interfaceschema.ParseInterfaceCanonical` as `installInterface` (internal/realm/service.go:196 vs 161) and then checks the body's `interface_name`/`version_major` against the URL's `{name}`/`{major}` (service.go:203-208), so the update body — the one carrying every interface rule — is the only interface body in the spec with nothing to validate against. `$ref` `#/components/schemas/InterfaceDefinition` as `installInterface` does (yaml:115-117), and say in the operation description that the body must name the same interface and major as the URL (disagreement is 409, already documented) and that `version_minor` must increase (409). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.

$ cd /root/astrate-mule && rg -n 'docs-sync-rm-' .mule/todo.md | tail -5 | cut -c1-90; echo "== total lines:"; wc -l < .mule/todo.md
254:- [ ] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail
255:- [ ] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy 
256:- [ ] docs-sync-rm-validationerror-example [auto]: fix or split the shared `Validation
257:- [ ] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in d
258:- [ ] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request bo
== total lines:
259

$ cd /root/astrate-mule && find . -newermt '-30 minutes' -type f -not -path './.git/*' -not -path './.mule/cron.log' 2>/dev/null; echo "== last line ok:"; sed -n '259p' .mule/todo.md | cut -c1-70
./internal/flow/blocks/container/httpbridge.go
./internal/flow/blocks/container/block_test.go
./internal/flow/blocks/schema.go
./internal/flow/blocks/info.go
./.mule/failed/container-response-cap-test.diff
./.mule/failed/container-response-cap-test.gate.log
./.mule/failed/gate.log
./.mule/log.md
./.mule/.budget
./.mule/.last-output
./.mule/.rotation
./.mule/reports/log/20261002T175003Z-container-response-cap-test-blocked.md
./.mule/.lock
./.mule/todo.md
== last line ok:
- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the def

Done: docs-sync recipe on the realm-management surface — compared all 20 mounted routes against the 12 documented paths (every route documented, no dead paths), then hand-checked status codes and field schemas against the handlers; five `docs/api/` corrections appended to `.mule/todo.md`, nothing implemented, no GitHub or git action taken.
Files: `.mule/todo.md` (lines 254-258 appended; the only file written — no code, no `docs/api/*.yaml`, no `.mule/for-giulio.md` entry was warranted)
Verified: `rg`/`read` cross-check of every route in `internal/realm/http.go:38-57` vs `docs/api/astarte_realm_management_api.yaml`, plus the two error writers (`http.go:352-425`), the service paths (`service.go:113-134, 309-359, 419-457, 515-578, 592-640`), `pkg/astarteapi/envelope.go:37,58` and `pkg/interfaceschema/parse.go:116-138,212-236,238-240,307-311` -> pass (5 findings, each traced to a source line; no build/gate run, this run changes no code)
Unsure: two things I deliberately did **not** propose. (1) `getVersion`'s example is still `data: "1.1.0"` (yaml:673) against `APICompatVersion = "1.2.2"` (service.go:588) — that is already task line 161 (`docs-sync-rm-version-example`), left `[!] BLOCKED: opencode exited 1`, so re-proposing it would duplicate; someone needs to clear that marker or it stays stuck. (2) `docs/site/configuration-reference.md` missing `ASTRATE_HOUSEKEEPING_*` keys is already escalated in `for-giulio.md` at least three times (lines 79, 152, 305-308) and is not converging — a fourth restatement would be noise. Also skipped the realm-not-found 404s on every operation: `realmID`'s `store.ErrNotFound` is unreachable behind the guard, which answers 401 for an unknown realm (auth/middleware.go:66-72), so those gaps are deliberate, not drift.
Follow-ups: five lines appended — `docs-sync-rm-error-example-capitalisation`, `docs-sync-rm-legacy-alias-fields`, `docs-sync-rm-validationerror-example`, `docs-sync-rm-deviceid-param`, `docs-sync-rm-update-interface-body`
