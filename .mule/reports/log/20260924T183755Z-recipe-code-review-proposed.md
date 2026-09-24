slug: recipe-code-review
verdict: proposed
at:  b3d4fd9
ran: 2026-09-24T18:37:55Z on DietPi in 164s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md

$ ls internal/ pkg/
internal/:
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

pkg/:
astarteapi
deviceid
interfaceschema
payload

$ ls .mule/reviews/ 2>/dev/null; echo "---LOG---"; tail -40 .mule/log.md 2>/dev/null
appengine-2026-09-07.md
astarteapi-2026-09-15.md
auth-2026-09-14.md
broker-2026-09-04.md
broker-2026-09-21.md
channels-2026-09-18.md
config-2026-09-13.md
deviceid-2026-09-18.md
engine-2026-09-06.md
engine-triggers-2026-09-20.md
flow-2026-09-05.md
flow-2026-09-23.md
flowapi-2026-09-10.md
housekeeping-2026-09-15.md
httpx-2026-09-17.md
interfaceschema-2026-09-11.md
observability-2026-09-16.md
pairing-2026-09-09.md
payload-2026-09-12.md
realm-2026-09-09.md
store-2026-09-08.md
store-2026-09-21.md
---LOG---
| 2026-09-14 | auth-cache-default-size-test | done | 157s | 874797b |
| 2026-09-15 | hk-zero-retention-asymmetry | done | 904s | 015a706 |
| 2026-09-15 | hk-zero-reglimit-create-asymmetry | blocked | 37s | opencode exited 1 |
| 2026-09-15 | hk-zero-retention-asymmetry | done | 567s | 203af1c |
| 2026-09-15 | docs-sync-native-socket-query-token-auth | done | 98s | 8148aae |
| 2026-09-15 | astarteapi-metadata-envelope-test | done | 238s | 1365cb8 |
| 2026-09-16 | astarteapi-multikey-sorted-golden | done | 215s | bfffbed |
| 2026-09-16 | docs-sync-pairing-deviceendpoints-dead-404-403 | done | 379s | f0d6a4d |
| 2026-09-16 | docs-sync-pairing-version-404-unreachable | done | 200s | bc73f65 |
| 2026-09-16 | observability-readiness-wedge-test | done | 178s | 4ef2e23 |
| 2026-09-16 | observability-dbpool-gauge-coverage | done | 62s | 19cd4f1 |
| 2026-09-16 | observability-health-content-type-test | done | 250s | f6c27f5 |
| 2026-09-17 | docs-sync-pairing-version-value | done | 1059s | c69a164 |
| 2026-09-17 | docs-sync-native-version-value | done | 503s | 08048b3 |
| 2026-09-17 | httpx-cors-vary-origin-passthrough | done | 143s | 986c100 |
| 2026-09-17 | docs-sync-rm-policies-delete-422 | done | 117s | c0521bc |
| 2026-09-17 | docs-sync-rm-triggers-422-nested-envelope | done | 161s | 8bfedd6 |
| 2026-09-18 | compat-note-v134 | done | 154s | bcd4383 |
| 2026-09-18 | channels-group-watch-membership-scope | done | 677s | aeaeb47 |
| 2026-09-18 | channels-rejoin-joinref-tagging-test | done | 208s | d48e17f |
| 2026-09-19 | channels-rejoin-authz-mismatch | blocked | 179s | wrote nothing |
| 2026-09-20 | errorname-missing-required-test | done | 135s | 739df67 |
| 2026-09-20 | triggers-custom-action-policy-nodecide | done | 459s | 3c8b89a |
| 2026-09-20 | compat-note-custom-action-policy-boundary | done | 97s | 8cd17a8 |
| 2026-09-20 | docs-sync-native-metrics-example-fake-series | done | 93s | 43c4a8b |
| 2026-09-20 | docs-sync-native-socket-missing-403-500 | done | 116s | f9d8b27 |
| 2026-09-21 | store-pipelines-empty-name-zero-blocks-test | done | 91s | 82fc0a4 |
| 2026-09-21 | broker-acl-coldstart-fallback-flood | done | 218s | 0e5078f |
| 2026-09-21 | broker-offlineacl-entry-eviction | done | 406s | 72b2686 |
| 2026-09-22 | docs-sync-ae-read-query-params | blocked | 159s | wrote nothing |
| 2026-09-23 | examples-echo-container-contract-test | done | 251s | f35af21 |
| 2026-09-23 | flow-mqtt-source-reconnect-recovery | blocked | 84s | wrote nothing |
| 2026-09-23 | flow-msg-json-integer-precision | done | 612s | 4df3f6b |
| 2026-09-23 | flow-sort-bounded-buffer | done | 634s | 9bf0ab7 |
| 2026-09-23 | flow-randomsource-span-overflow | blocked | 288s | lint failed: internal/flow/blocks/randomsource.go:113:40: G115: integer overflow conversion int64 -> uint64 (gosec) |
| 2026-09-23 | flow-filter-key-contains-test | done | 107s | 3875ba6 |
| 2026-09-23 | docs-sync-rm-put-interface-409 | done | 356s | 7256af9 |
| 2026-09-24 | docs-sync-rm-interface-422-shapes | done | 490s | 93e784e |
| 2026-09-24 | docs-sync-rm-validation-example-prefix | done | 238s | 02e269b |
| 2026-09-24 | docs-sync-rm-auth-403 | done | 668s | 9c05370 |

$ rg -c '' --glob 'internal/swagger/**/*.go' | sort -t: -k2 -rn | head; echo "---"; ls -la internal/swagger/
internal/swagger/swagger_test.go:192
internal/swagger/swagger.go:43
---
total 20
drwxr-xr-x  2 root root 4096 Sep  5 19:42 .
drwxr-xr-x 17 root root 4096 Sep  4 20:46 ..
-rw-r--r--  1 root root 1414 Jul 27 05:40 swagger.go
-rw-r--r--  1 root root 4865 Sep 24 18:54 swagger_test.go

→ Read internal/swagger/swagger.go
→ Read internal/swagger/swagger_test.go

$ rg -n "swagger.Mount|swagger\.Specs|Mount\(" --glob '*.go' | rg -v 'internal/swagger' ; echo "=== specs usage ==="; rg -n "Specs\(\)" --glob '*.go' | rg -v 'internal/swagger/swagger_test'
cmd/astrate/main.go:359:	}).Mount(mux)
cmd/astrate/main.go:368:	housekeeping.NewAPI(hkSvc, mw, hkKeys).Mount(mux)
cmd/astrate/main.go:372:	realm.NewAPI(realmSvc, mw).Mount(mux)
cmd/astrate/main.go:373:	appengine.NewAPI(appengine.NewService(st, e, log), mw).Mount(mux)
cmd/astrate/main.go:374:	apstream.NewAPI(e.Bus(), mw).Mount(mux)
cmd/astrate/main.go:377:	channels.NewAPI(e.Bus(), st, st).Mount(mux)
cmd/astrate/main.go:379:	flowapi.NewAPI(flowSvc, mw).Mount(mux)
cmd/astrate/main.go:411:	health.Mount(mux)
cmd/astrate/main.go:414:	swagger.Mount(mux)
test/conformance/instance/instance.go:183:	pairing.NewAPI(pairer, auth.NewMiddleware(st), pairing.APIConfig{}).Mount(mux)
test/conformance/cpb/cpb_test.go:222:	api.Mount(mux)
test/conformance/cpc/cpc_test.go:224:	pairing.NewAPI(pairer, mw, pairing.APIConfig{}).Mount(mux)
test/conformance/cpc/cpc_test.go:225:	housekeeping.NewAPI(housekeeping.NewService(st, sealer, b, nil), mw, []string{instPubPEM}).Mount(mux)
test/conformance/cpc/cpc_test.go:226:	realm.NewAPI(realm.NewService(st, e, nil), mw).Mount(mux)
test/conformance/cpc/cpc_test.go:227:	appengine.NewAPI(appengine.NewService(st, e, nil), mw).Mount(mux)
test/conformance/cpc/cpc_test.go:228:	apstream.NewAPI(e.Bus(), mw).Mount(mux)
test/conformance/cpa/cpa_test.go:193:	api.Mount(mux)
test/security/security_test.go:144:	pairing.NewAPI(pairer, mw, pairing.APIConfig{RegisterRate: 0.01, RegisterBurst: 2}).Mount(mux)
test/security/security_test.go:145:	housekeeping.NewAPI(housekeeping.NewService(st, sealer, b, discard()), mw, []string{instPub}).Mount(mux)
test/security/security_test.go:146:	realm.NewAPI(realm.NewService(st, e, discard()), mw).Mount(mux)
test/security/security_test.go:147:	appengine.NewAPI(appengine.NewService(st, e, discard()), mw).Mount(mux)
test/security/security_test.go:148:	apstream.NewAPI(e.Bus(), mw).Mount(mux)
internal/appengine/stream/ws_test.go:60:	NewAPI(bus, auth.NewMiddleware(fakeKeys{realm})).Mount(mux)
internal/appengine/stream/ws.go:45:func (a *API) Mount(mux *http.ServeMux) {
internal/appengine/http.go:40:func (a *API) Mount(mux *http.ServeMux) {
internal/flowapi/http.go:31:func (a *API) Mount(mux *http.ServeMux) {
internal/realm/http.go:36:func (a *API) Mount(mux *http.ServeMux) {
internal/realm/http_test.go:90:	NewAPI(svc, auth.NewMiddleware(st)).Mount(mux)
internal/observability/observability_test.go:48:	h.Mount(mux)
internal/observability/observability_test.go:89:	h.Mount(mux)
internal/observability/observability_test.go:112:	h.Mount(mux)
internal/observability/health.go:43:func (h *Health) Mount(mux *http.ServeMux) {
internal/pairing/http.go:77:func (a *API) Mount(mux *http.ServeMux) {
internal/appengine/channels/ws_test.go:61:	api.Mount(mux)
internal/appengine/channels/ws_test.go:74:	api.Mount(mux)
internal/appengine/channels/ws_test.go:87:	api.Mount(mux)
internal/appengine/channels/ws_test.go:99:	api.Mount(mux)
internal/appengine/channels/ws_test.go:112:	api.Mount(mux)
internal/appengine/channels/ws_test.go:141:	api.Mount(mux)
internal/appengine/channels/ws_test.go:194:	api.Mount(mux)
internal/appengine/channels/ws_test.go:234:	api.Mount(mux)
internal/appengine/channels/ws_test.go:290:	api.Mount(mux)
internal/appengine/channels/ws_test.go:339:	api.Mount(mux)
internal/appengine/channels/ws_test.go:410:	api.Mount(mux)
internal/appengine/channels/ws_test.go:444:	api.Mount(mux)
internal/appengine/channels/ws_test.go:546:	api.Mount(mux)
internal/appengine/channels/ws_test.go:637:	api.Mount(mux)
internal/appengine/channels/ws_test.go:710:	api.Mount(mux)
internal/appengine/channels/ws_test.go:802:			api.Mount(mux)
internal/appengine/channels/ws_test.go:975:			api.Mount(mux)
internal/appengine/channels/ws.go:47:func (a *API) Mount(mux *http.ServeMux) {
internal/appengine/e2e_test.go:236:	channels.NewAPI(env.engine.Bus(), env.st, env.st).Mount(mux)
internal/housekeeping/http.go:35:func (a *API) Mount(mux *http.ServeMux) {
internal/appengine/http_test.go:149:	NewAPI(NewService(st, sd, discardLogger()), auth.NewMiddleware(st)).Mount(mux)
internal/housekeeping/http_test.go:91:		auth.NewMiddleware(st), []string{instancePub}).Mount(mux)
internal/pairing/http_test.go:102:	api.Mount(mux)
=== specs usage ===
internal/swagger/swagger.go:31:func Specs() []string {

$ rg -n "SwaggerUI|APIYAML|go:embed" docs/*.go
7:// SwaggerUI holds the static Swagger UI files (index.html, CSS, JS refs).
9://go:embed swagger-ui/*
10:var SwaggerUI embed.FS
12:// APIYAML holds the OpenAPI 3.0 YAML specifications for all five API surfaces.
14://go:embed api/*.yaml
15:var APIYAML embed.FS

Confirmed the fs layout facts. The one real finding is the ignored `fs.Sub` errors in `Mount`. Writing the review file, then appending the proposal to todo.md.
← Write .mule/reviews/swagger-2026-09-24.md
Wrote file successfully.

$ tail -15 .mule/todo.md
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

$ printf '%s\n' '- [ ] swagger-sub-failfast [auto]: make `Mount` fail fast (panic) when `fs.Sub` over the `docs.SwaggerUI`/`docs.APIYAML` embeds errors — `uiRoot, _ := fs.Sub(...)` and `apiRoot, _ := fs.Sub(...)` (internal/swagger/swagger.go:17-18) silently serve an empty tree → 404s at /swagger/ and /api/ with no log the day the docs embed layout changes; add an fs-injecting internal variant (`MountWithFS(mux, uiRoot, apiRoot fs.FS)` called by `Mount` after the subs) so a broken-FS unit test in swagger_test.go can prove the fail-fast, keeping the public `Mount(mux)` signature and its single main.go:414 caller untouched. Container-free.' >> .mule/todo.md
(no output)

$ tail -2 .mule/todo.md
- [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] swagger-sub-failfast [auto]: make `Mount` fail fast (panic) when `fs.Sub` over the `docs.SwaggerUI`/`docs.APIYAML` embeds errors — `uiRoot, _ := fs.Sub(...)` and `apiRoot, _ := fs.Sub(...)` (internal/swagger/swagger.go:17-18) silently serve an empty tree → 404s at /swagger/ and /api/ with no log the day the docs embed layout changes; add an fs-injecting internal variant (`MountWithFS(mux, uiRoot, apiRoot fs.FS)` called by `Mount` after the subs) so a broken-FS unit test in swagger_test.go can prove the fail-fast, keeping the public `Mount(mux)` signature and its single main.go:414 caller untouched. Container-free.

```
Done: reviewed internal/swagger (small, never-reviewed package) — one proposal appended to todo.md
Files: .mule/reviews/swagger-2026-09-24.md, .mule/todo.md
Verified: read swagger.go + swagger_test.go (all), grepped all Mount/Specs callers, checked docs/*.go embeds -> pass
Unsure: nothing
Follow-ups: swagger-sub-failfast (fs.Sub errors silently ignored at swagger.go:17-18)
```
