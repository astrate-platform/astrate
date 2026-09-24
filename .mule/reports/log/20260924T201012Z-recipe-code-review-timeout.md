slug: recipe-code-review
verdict: timeout
at:  f6b35a5
ran: 2026-09-24T20:10:12Z on DietPi in 1200s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md

$ ls internal/ pkg/ 2>/dev/null
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

$ ls .mule/reviews/ 2>/dev/null; echo "---"; ls .mule/ 2>/dev/null
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
swagger-2026-09-24.md
---
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
recipes
reports
reviews
tasks
todo.md

$ rg -n 'channels' -l internal/ pkg/ --glob '*.go' | head; echo "---"; ls internal/testutil/ 2>/dev/null
internal/flow/router_test.go
internal/store/groups.go
internal/flow/router.go
internal/engine/stream/bus.go
internal/appengine/channels/protocol_test.go
internal/appengine/channels/room_test.go
internal/auth/claims_test.go
internal/appengine/channels/room.go
internal/appengine/channels/ws_test.go
internal/appengine/channels/protocol.go
---
astartedevice.go
golden.go
golden_test.go
mqttclient.go
pg.go
pg_test.go

$ tail -40 .mule/log.md
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
| 2026-09-24 | swagger-sub-failfast | done | 371s | 1d321ea |

$ for p in internal/observability internal/testutil internal/broker internal/store; do echo "== $p"; rg -c '' --glob '*.go' $p 2>/dev/null | sort -t: -k2 -rn | head -5; done
== internal/observability
internal/observability/observability_test.go:180
internal/observability/health.go:84
internal/observability/metrics.go:80
internal/observability/compat_version_test.go:55
internal/observability/compat.go:52
== internal/testutil
internal/testutil/astartedevice.go:216
internal/testutil/mqttclient.go:163
internal/testutil/golden_test.go:117
internal/testutil/pg.go:106
internal/testutil/golden.go:59
== internal/broker
internal/broker/broker_test.go:640
internal/broker/broker_e2e_test.go:618
internal/broker/sessionstore.go:533
internal/broker/authhook.go:496
internal/broker/broker.go:343
== internal/store
internal/store/datastreams_test.go:823
internal/store/devices.go:485
internal/store/datastreams.go:434
internal/store/devices_test.go:427
internal/store/interfaces.go:398

$ rg -c '' --glob '*.go' internal/engine/*.go 2>/dev/null | sort -t: -k2 -rn | head; echo "---subdirs---"; ls internal/engine/
internal/engine/fakestore_test.go:818
internal/engine/data.go:552
internal/engine/cache.go:502
internal/engine/router.go:490
internal/engine/data_test.go:461
internal/engine/batch.go:459
internal/engine/engine.go:456
internal/engine/batch_test.go:446
internal/engine/engine_e2e_test.go:440
internal/engine/engine_integration_test.go:435
---subdirs---
batch.go
batch_test.go
bus_events_test.go
cache.go
cache_test.go
capabilities.go
capabilities_test.go
control.go
control_test.go
data.go
data_test.go
devicedata_test.go
engine.go
engine_e2e_test.go
engine_integration_test.go
fakestore_test.go
forward
introspection.go
introspection_test.go
mqttv1_parity_test.go
router.go
router_test.go
serverdata.go
serverdata_test.go
stream
testdata
topics.go
topics_test.go
triggers
triggers_change_integration_test.go
triggers_change_test.go

→ Read .mule/reviews/engine-2026-09-06.md
The engine review is old (2026-09-06) but appengine's main package likewise (2026-09-07, only its `channels` subpackage re-reviewed since). Let me pick the right candidate.

$ tail -25 .mule/todo.md
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

$ rg -c '' --glob '*.go' internal/appengine/*.go 2>/dev/null | sort -t: -k2 -rn | head -15; echo "---total files---"; ls internal/appengine/*.go | wc -l
internal/appengine/service.go:846
internal/appengine/http.go:696
internal/appengine/data.go:514
internal/appengine/http_test.go:512
internal/appengine/formats_parity_test.go:467
internal/appengine/e2e_test.go:351
internal/appengine/groups_parity_test.go:309
internal/appengine/mirror_test.go:243
internal/appengine/query_opts_test.go:177
internal/appengine/dashboard_compat_test.go:160
internal/appengine/patch_test.go:113
internal/appengine/downsample_test.go:93
internal/appengine/links_test.go:65
internal/appengine/writeerror_test.go:55
internal/appengine/downsample.go:29
---total files---
15

→ Read .mule/reviews/appengine-2026-09-07.md
The oldest main-package reviews are `engine` (2026-09-06) and `appengine` (2026-09-07), both ~2.5 weeks old; `testutil` has never been reviewed. Let me check whether the prior engine review's top proposal was ever done, and look at what changed since.

$ rg -n 'publishAsOwner|devicedata|PublishDeviceValue|UnsetServerProperty' .mule/todo.md .mule/log.md 2>/dev/null | head -20
.mule/todo.md:131:- [x] server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.

$ git log --oneline -30 -- internal/engine/ | head -40
3c8b89a mule: triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
739df67 mule: errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
e46ce11 mule: server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.
fd562f7 payload: pre-wire missing_required reject reason label fixture and error-name translation (fase 4b prep)
5e1aa1d feat(engine): PublishDeviceValue — device-owned ingest with no MQTT delivery (issue #84, phase 84a)
ed0342d feat(triggers): upstream wildcard semantics + introspection interface rules (#70)
edf4e1b feat(triggers): upstream-measured HTTP action validation limits (#63)
029e0b6 feat(triggers): pairwise-disjoint policy handlers, prefetch_count bounds, AMQP action rejection (issues #65, #64)
06cb17a engine: tolerate unknown control subpaths silently (issue #51)
094d5c1 test(engine): MQTT v1 parity conformance for empty introspection (#48) and binaryblob BSON subtype-0 resend (#49)
6bd14a7 feat(triggers): evaluate change-derived data conditions (value_change, path_created/removed, value_stored)
4c881d0 fix(lint): resolve golangci-lint v2.12.2 failures across flow container work
aba2d18 fix: land blocked mule issues #14 #21 #22 #24 #27
b20a302 mule: issue-12: device-registered-trigger: emit DeviceRegisteredEvent from internal/pairing/service.go on registration
8bd9df9 mule: issue-15: prefetch-count-delivery-policies: add PrefetchCount to the trigger delivery policy config
4e6db62 mule: issue-16: capabilities-topic: handle the <realm>/<device_id>/capabilities BSON topic
3e23d8a M12-08c/08d: NATS bus Forwarder, wired end to end behind -tags nats (#11)
e53e57b Match the recorded upstream: housekeeping 404 envelope, two error_name rows
3ec2753 Select the trigger forwarder from [triggers.forward] config
9c6b22b Add an HTTP bus Forwarder that posts a custom trigger action verbatim
699cf11 Initial commit under updated license terms

$ git log --oneline --since=2026-09-06 -- internal/engine/serverdata.go internal/engine/devicedata_test.go internal/engine/batch.go internal/engine/engine.go internal/engine/data.go internal/engine/router.go internal/engine/cache.go internal/engine/control.go 2>/dev/null
e46ce11 mule: server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.

The engine main package was last reviewed 18 days ago (2026-09-06) and has since gained `PublishDeviceValue`, the server-data trigger/bus work, and parity conformance. Its top proposal was landed (server-data-trigger-bus). I'll re-review `internal/engine` main package (excluding triggers, which was reviewed 09-20).

Let me look at what changed since the last engine review.
$ wc -l internal/engine/serverdata.go internal/engine/devicedata_test.go internal/engine/fakestore_test.go internal/engine/introspection.go internal/engine/topics.go internal/engine/capabilities.go internal/engine/control.go internal/engine/mqttv1_parity_test.go
349 internal/engine/serverdata.go
  223 internal/engine/devicedata_test.go
  818 internal/engine/fakestore_test.go
  168 internal/engine/introspection.go
   97 internal/engine/topics.go
  102 internal/engine/capabilities.go
  316 internal/engine/control.go
  151 internal/engine/mqttv1_parity_test.go
 2224 total

→ Read internal/engine/serverdata.go
→ Read internal/engine/engine.go
Now let me examine the pipeline core files that the server path shares, plus `capturePrevious`.
→ Read internal/engine/data.go
Let me look at the remaining core files: `batch.go`, `router.go`, and check whether the shared decoder is stateless.
✗ Read pkg/payload/decoder.go failed
Error: File not found: /root/astrate-mule/pkg/payload/decoder.go
→ Read internal/engine/batch.go

$ ls pkg/payload/
bson.go
fuzz_test.go
json.go
payload.go
payload_test.go
required_test.go
sniff.go
sniff_test.go
testdata
value.go
value_test.go

$ rg -n 'e.dec|dec =' internal/ | head
internal/flow/blocks/catalog.go:23:	// TypeFilter and TypeMap are declared in transform.go.
internal/flow/composite.go:43:// composite bodies (no global source/sink) can be decoded and re-encoded
internal/engine/data.go:119:	// Value is the decoded payload value (payload.Value closed set;
internal/engine/data.go:258:	// Step 3 — ownership (checked before decoding: cheaper, and the payload
internal/engine/data.go:278:		dp, err = e.dec.Object(m.Payload, ci.ObjectLeaves)
internal/engine/data.go:286:		dp, err = e.dec.Individual(m.Payload, mapping)
internal/flow/message.go:268:// setDataFromWire decodes a scalar wire data value into m.Data.
internal/flow/message.go:376:// silently truncate or wrap (3.7 → 3, 1e300 → garbage). Wire decoding
internal/flow/message.go:377:// delivers json.Number (the decoder runs with UseNumber), which is judged
internal/engine/control.go:28:// at the declared size and an absolute ceiling"). 1 MiB comfortably holds

$ rg -n 'Dec()|dec ' internal/engine/router.go internal/engine/cache.go internal/engine/engine.go | head
internal/engine/router.go:73:	dec     payload.Decoder
internal/engine/router.go:175:		dec:     payload.Decoder{MaxSize: cfg.MaxPayloadBytes},

→ Read pkg/payload/payload.go
Let me look at the cache and control files, which the previous review covered but may have changed, plus check the new test.
→ Read internal/engine/cache.go
→ Read internal/engine/control.go
Let me read the remaining core files and the relevant tests.
→ Read internal/engine/devicedata_test.go
→ Read internal/engine/introspection.go
→ Read internal/engine/capabilities.go
→ Read internal/engine/topics.go
→ Read internal/engine/router.go

$ ls internal/engine/stream/ internal/engine/triggers/ internal/engine/forward/
internal/engine/forward/:
envelope.go
http.go
http_test.go
nats.go
nats_test.go

internal/engine/stream/:
bus.go
bus_test.go

internal/engine/triggers/:
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

→ Read internal/engine/mqttv1_parity_test.go
→ Read internal/engine/stream/bus.go
→ Read internal/engine/serverdata_test.go

$ rg -n 'DefaultMaxSize|MaxSize' pkg/payload/*.go internal/appengine/*.go | head -30
pkg/payload/value.go:92:	// DefaultMaxSize is the default accepted payload size for both
pkg/payload/value.go:94:	DefaultMaxSize = 64 << 10
pkg/payload/payload_test.go:562:	bigBlob, err := encodeBSON(make([]byte, DefaultMaxSize), nil)
pkg/payload/payload_test.go:570:	bigJSON := []byte(`{"v":"` + strings.Repeat("a", DefaultMaxSize) + `"}`)
pkg/payload/payload_test.go:575:	small := Decoder{MaxSize: 8}
pkg/payload/payload_test.go:584:	wide := Decoder{MaxSize: 1 << 20}
pkg/payload/payload.go:24:	// MaxSize caps the accepted payload size in bytes for both formats
pkg/payload/payload.go:25:	// (docs/DESIGN.md §3.5.3); 0 means DefaultMaxSize.
pkg/payload/payload.go:26:	MaxSize int
pkg/payload/payload.go:31:	if d.MaxSize > 0 {
pkg/payload/payload.go:32:		return d.MaxSize
pkg/payload/payload.go:34:	return DefaultMaxSize

✱ Grep "PublishServerValue|PublishDeviceValue|UnsetServerProperty" 69 matches

$ rg -n 'func \(x \*Executor\) Enqueue|type Executor|queue chan|sync\.' internal/engine/triggers/actions.go | head -30
210:type ExecutorConfig struct {
236:type Executor struct {
240:	wg       sync.WaitGroup
242:	stopOnce sync.Once
250:	policyMu       sync.Mutex
254:	sendMu sync.RWMutex
325:func (x *Executor) Enqueue(d Delivery) {

→ Read internal/appengine/data.go [offset=430, limit=90]
Let me check the fake broker port and the executor's Enqueue for concurrency safety (the server path now fires triggers/bus from AppEngine goroutines, separate from shards).

$ rg -n 'type fake|func .*published|func newWiredRig|brokerPort|Publish func|errPort|failPublish' internal/engine/fakestore_test.go | head -30
102:type fakeStore struct {
731:type fakePort struct {
759:func (p *fakePort) published() []pubCall {
768:func (p *fakePort) publishedTo(topic string) []pubCall {
788:func newWiredRig(t testing.TB, cfg Config) (*pipelineRig, *fakeStore, *fakePort) {

→ Read internal/engine/triggers/actions.go [offset=210, limit=130]

$ rg -n 'ObjectLeaves|aggregation' pkg/interfaceschema/*.go | head -20
pkg/interfaceschema/parse.go:94:	// enforced at runtime only for object aggregation (upstream 1.4
pkg/interfaceschema/parse.go:120:	// Aggregate is the upstream legacy alias of aggregation (issue #61).
pkg/interfaceschema/parse.go:122:	Aggregation *Aggregation  `json:"aggregation,omitempty"`
pkg/interfaceschema/parse.go:156:// aggregation and per-type field constraints) are rejected with an error
pkg/interfaceschema/parse.go:214:	// aggregate only without aggregation (both → violation on aggregation).
pkg/interfaceschema/parse.go:227:			vc.add("aggregation", "aggregation and aggregate are mutually exclusive")
pkg/interfaceschema/parse.go:284:		return nil, errors.New("properties interfaces cannot use object aggregation")
pkg/interfaceschema/parse.go:595:// checkObjectAggregation enforces the object-aggregation shape: every
pkg/interfaceschema/types.go:117:// ParseAggregation parses the wire form of an aggregation.
pkg/interfaceschema/types.go:125:		return 0, fmt.Errorf("unknown aggregation %q", s)
pkg/interfaceschema/compat.go:26:// ownership, and aggregation; strictly increase the minor version; keep
pkg/interfaceschema/compat.go:60:		return fmt.Errorf("%w: %w: aggregation changed from %s to %s",
pkg/interfaceschema/parse_test.go:77:	t.Run("Geolocation object aggregation", func(t *testing.T) {
pkg/interfaceschema/required_encrypted_test.go:92:		"type": "datastream", "ownership": "device", "aggregation": "object",
pkg/interfaceschema/compat_test.go:450:		"type": "datastream", "ownership": "device", "aggregation": "object",
pkg/interfaceschema/compat_test.go:454:	if err == nil || !strings.Contains(err.Error(), "aggregation changed") {
pkg/interfaceschema/compat_test.go:455:		t.Errorf("err = %v, want aggregation change rejection", err)
pkg/interfaceschema/compile.go:22:// object aggregation, the last-level key → mapping table.
pkg/interfaceschema/compile.go:38:	// ObjectLeaves maps, for object aggregation, each last-level name to its
pkg/interfaceschema/compile.go:39:	// mapping; nil for individual aggregation.

Let me verify a couple of suspicions: whether `objectPathOK`'s single-leaf check is sound, and how the store fake port is built (for a propose-able ordering test).
→ Read internal/engine/fakestore_test.go [limit=120, offset=731]
→ Read pkg/interfaceschema/parse.go [limit=40, offset=590]
The engine main package (last reviewed 2026-09-06, since gained PublishDeviceValue, server-data trigger/bus wiring, capabilities). Let me verify a couple of specific suspicions before writing up.

$ rg -n 'objectPathOK|ObjectLeaves' internal/engine/*.go pkg/interfaceschema/*.go | head
internal/engine/cache_test.go:49:	if obj := alpha.iface("com.astrate.test.ObjectFlat", 1); obj == nil || obj.ObjectLeaves["latitude"] == nil {
pkg/interfaceschema/compile_test.go:72:	if ci.ObjectLeaves != nil {
pkg/interfaceschema/compile_test.go:73:		t.Errorf("ObjectLeaves = %v, want nil for individual aggregation", ci.ObjectLeaves)
pkg/interfaceschema/compile_test.go:110:	if len(ci.ObjectLeaves) != len(wantLeaves) {
pkg/interfaceschema/compile_test.go:111:		t.Fatalf("len(ObjectLeaves) = %d, want %d", len(ci.ObjectLeaves), len(wantLeaves))
pkg/interfaceschema/compile_test.go:114:		m, ok := ci.ObjectLeaves[leaf]
pkg/interfaceschema/compile_test.go:116:			t.Errorf("ObjectLeaves missing %q", leaf)
pkg/interfaceschema/compile_test.go:128:	if m != ci.ObjectLeaves["latitude"] {
pkg/interfaceschema/compile_test.go:129:		t.Error("trie leaf and ObjectLeaves entry are not the same mapping")
internal/engine/data.go:273:		if !objectPathOK(ci, path) {

$ rg -n 'pubErrs' internal/engine/ -l; echo "---"; rg -n 'pubErrs' internal/engine/fakestore_test.go
internal/engine/fakestore_test.go
---
735:	// pubErrs holds errors returned by the next Publish calls.
736:	pubErrs []error
742:	if len(p.pubErrs) > 0 {
743:		err := p.pubErrs[0]
744:		p.pubErrs = p.pubErrs[1:]

## Objective
- Execute the code-review recipe on one area of the Astrate codebase: read `.mule/MULE.md` + `.mule/recipes/code-review.md` (done), pick an area not recently reviewed, go deep, then write a review evidence file under `.mule/reviews/` and append up to 5 task-proposal lines to `.mule/todo.md`.
- No code changes, no git operations. The chosen area (after rotation analysis) is the `internal/engine` main package, last fully reviewed 2026-09-06.

## Important Details
- Recipe is proposal-only: deliverable is the review file + new `.mule/todo.md` lines; do NOT implement proposed tasks, do NOT touch git (no commit/branch/checkout/stash).
- This recipe does NOT ask for `gh issue create` → do not run any `gh` commands at all.
- Rotation rule: check `.mule/log.md` and `.mule/reviews/`; picked `internal/engine` because its main-package review (2026-09-06) is the oldest not re-reviewed since; the `triggers` sub-area was re-reviewed 2026-09-20 (`engine-triggers-2026-09-20.md`), so the main package is in scope and triggers are not.
- Proposal value ordering from recipe: wrong/unguarded behavior > missing test on existing rule > perf w/ named workload > mechanical clarity > features. Renames/reformatting: never propose. Finding little in a well-written package is a correct result — say so.
- MULE.md standing rules: English only; issue #1 never proposed/raised; if a task grows, do the first useful part and add the rest as todo lines; read only what the task needs.
- Prior engine review's top finding (`server-data-trigger-bus`) was already landed (todo.md:131, commit e46ce11) — do not re-propose it.
- Review file naming convention: `.mule/reviews/<pkg>-<date>.md` (e.g. `appengine-2026-09-07.md`); most recent review is `swagger-2026-09-24.md`.
- `.mule/todo.md` task-line style observed: `- [ ] <slug>: <problem w/ file:line> … <test requirement, container-free when relevant>`.

## Work State
### Completed
- Read `.mule/MULE.md` (truncated view — standing instructions/rule #1 present) and `.mule/recipes/code-review.md` (full intent captured above).
- Inventoried areas: `internal/{appengine auth broker config engine flow flowapi housekeeping httpx observability pairing realm store swagger testutil}` and `pkg/{astarteapi deviceid interfaceschema payload}`.
- Mapped review history from `.mule/reviews/`: swagger 09-24, flow 09-23, broker+store 09-21, engine-triggers 09-20, channels+deviceid 09-18, httpx 09-17, observability 09-16, astarteapi+housekeeping 09-15, auth 09-14, config 09-13, payload 09-12, interfaceschema 09-11, flowapi 09-10, pairing+realm 09-09, store 09-08, engine 09-06, flow 09-05, broker 09-04, appengine 09-07; `testutil` never reviewed.
- Read prior `reviews/engine-2026-09-06.md` (background; its proposals mostly landed) and tail of `.mule/log.md` and `.mule/todo.md`.
- Read engine files largest-first + tests: `serverdata.go`, `engine.go`, `data.go` (partial), `batch.go`, `cache.go`, `control.go`, `introspection.go`, `capabilities.go`, `topics.go`, `router.go` (partial), `stream/bus.go`, `serverdata_test.go`, `devicedata_test.go`, `mqttv1_parity_test.go`.
- Verified: `payload.Decoder` is stateless (only `MaxSize`), so shared `e.dec` is safe; `triggers.Executor.Enqueue` is mutex-safe (sendMu) for calls from AppEngine goroutines; AppEngine `PublishData` calls `s.sd.PublishServerValue` (internal/appengine/data.go:466-475; 487 for UnsetProperty).

### Active
- Finalizing the finding list and writing the review evidence file; no files written yet.

### Blocked
- None.

## Next Move
1. Write `.mule/reviews/engine-<today's-date>.md` recording: area choice rationale (engine main package, stale since 09-06; triggers excluded due to 09-20 review), what was read, and the findings below with file:line.
2. Candidate findings to propose as todo lines (pick the strongest, ≤5):
   - **Failure-path ordering in server-owned writes**: `PublishServerValue`/`UnsetServerProperty` (internal/engine/serverdata.go, `publishAsOwner` ~98-214, `UnsetServerProperty` ~222-279, `PublishServerValue` ~up to 62) persist the op and fire `fireData` (triggers + live bus) *before* `broker.Publish`; a broker-publish failure returns an error to the caller although the value is committed and triggered — no test pins this. Propose a container-free test that forces a publish failure and asserts the resulting persisted/triggered state, deciding intent.
   - **`handleEmptyCache` arms hint reset before resync succeeds** (internal/engine/control.go): `dev.armHintReset()` runs before `retryStore`; on store failure the flag is armed anyway and the next BSON flips the hint as if emptyCache had been honored. Pin or fix.
   - **`decodeCapabilities` silently drops wrong-typed known capability values** (internal/engine/capabilities.go, `decodeCapabilities`/`validateCapability`): a non-string `purge_properties_compression_format` is dropped, never validated, silently defaulting to zlib. Decide reject-vs-accept and add a test.
   - **No-op unset emits events**: unsetting an absent property (both `UnsetServerProperty` and broker-ingress unset) still fires `value_change` null→null / bus `incoming_data`; behavior is consistent across paths but undocumented and untested.
3. Append the accepted proposals to `.mule/todo.md` in the existing `- [ ] slug:` format, each with file:line and a machine-checkable test requirement (container-free where possible).
4. Do NOT run git commands or `gh`; do NOT start implementing.

## Relevant Files
- `.mule/recipes/code-review.md`: the recipe being executed (proposal-only, rotation, output format).
- `.mule/MULE.md`: project safety rules (English only, never-raise issue #1, split large work).
- `internal/engine/*.go`: the area under review — `serverdata.go` (server/device-owned publish path, recently changed by e46ce11), `engine.go` (`fireData`/`fireCommitted`), `control.go`, `capabilities.go`, `data.go`, `batch.go`, `router.go`, `cache.go`, `introspection.go`, `topics.go`.
- `internal/engine/stream/bus.go`: live bus event kinds (`KindIncomingData`, etc.) used by `fireData`.
- `internal/engine/serverdata_test.go`, `devicedata_test.go`, `mqttv1_parity_test.go`, `fakestore_test.go` (`fakePort`, `newWiredRig`): existing test coverage and test rig for the proposed tests.
- `.mule/reviews/engine-2026-09-06.md` and `engine-triggers-2026-09-20.md`: prior reviews of this area (avoid duplicate proposals; triggers out of scope).
- `.mule/todo.md` + `.mule/log.md`: task append target and rotation evidence.
- `internal/appengine/data.go`: HTTP-layer caller of `PublishServerValue`/`UnsetServerProperty` (why the error-ordering matters to API semantics).
Continue if you have next steps, or stop and ask for clarification if you are unsure how to proceed.
Let me re-read the two exact spots I want to cite precisely before writing up.
✱ Grep "func objectPathOK" in internal/engine · 1 match
→ Read internal/engine/data.go [limit=30, offset=488]
Let me re-read the exact flow in the two spots I want to propose findings on, to be precise.
→ Read internal/engine/serverdata.go [limit=120, offset=40]
→ Read internal/engine/serverdata.go [limit=100, offset=160]
→ Read internal/engine/serverdata.go [limit=90, offset=260]
✱ Grep "func (e \*Engine) fireData|func (e \*Engine) fireDataChanges|func fireDataChanges" in internal/engine · 0 matches
→ Read internal/engine/control.go [limit=60, offset=40]

$ rg -n 'fireData|fireDataChanges|path_removed|created|removed|changed' internal/engine/engine.go | head -40
147:		e.fireData(rs, op)
151:// fireData evaluates data triggers for one committed op and publishes it on
153:func (e *Engine) fireData(rs *realmSchema, op *PersistOp) {
203:			e.fireDataChanges(rs, op, value)
219:// fireDataChanges emits the previous-value-derived events for one committed
224://   - path_created when the path had none and now receives a value;
225://   - path_removed when an existing property is unset;
227:func (e *Engine) fireDataChanges(rs *realmSchema, op *PersistOp, value any) {
236:	changed := !op.PrevExists || !jsonEqual(newJSON, op.Prev)
237:	created := !op.PrevExists && op.Value != nil
238:	removed := op.PrevExists && op.Value == nil
266:	if changed {
270:	if created {
274:	if removed {
278:	if changed {

→ Read internal/engine/engine.go [limit=135, offset=147]
→ Read internal/engine/data.go [limit=35, offset=517]
✱ Grep "func \(d \*deviceState\) armHintReset" in internal/engine · 1 match
→ Read internal/engine/cache.go [limit=30, offset=390]
Let me verify the `capturePrevious` behavior for the absent-property unset case, and confirm `pubErrs` is truly unused.

$ rg -n 'func capturePrevious|PrevExists|Tracked' internal/engine/serverdata.go internal/engine/data.go | head; echo ---; rg -n 'capturePrevious' internal/engine/*.go
internal/engine/data.go:133:	// previous value, PrevExists reports whether there was one, Tracked
internal/engine/data.go:136:	Tracked    bool
internal/engine/data.go:137:	PrevExists bool
internal/engine/data.go:406:	op.Tracked = true
internal/engine/data.go:407:	op.PrevExists = prevFound
internal/engine/data.go:479:	op.Tracked = true
internal/engine/data.go:480:	op.PrevExists = prevFound
---
internal/engine/serverdata.go:184:	e.capturePrevious(ctx, rs, &op)
internal/engine/serverdata.go:263:	e.capturePrevious(ctx, rs, &op)
internal/engine/data.go:446:// capturePrevious fills the previous-value snapshot of an op that persists
internal/engine/data.go:454:func (e *Engine) capturePrevious(ctx context.Context, rs *realmSchema, op *PersistOp) {

→ Read internal/engine/data.go [limit=35, offset=446]
I have enough for the change-triggers and failure-path findings. Let me verify two last things: whether `dev.declares` is mutex-guarded (deviceView race), and whether any test actually exercises `pubErrs`.

$ rg -n 'func \(d \*deviceState\) declares|func \(d \*deviceState\) declared|declares ' internal/engine/cache.go | head; echo ---; rg -rn 'pubErrs\s*=' internal/engine/ | head; echo ---; rg -n 'func \(d \*deviceState\) hint\b' internal/engine/cache.go
370:// declares reports whether the device's introspection contains name, and at
372:func (d *deviceState) declares(name string) (store.InterfaceVersion, bool) {
---
internal/engine/fakestore_test.go:		p.n p.pubErrs[1:]
---
388:func (d *deviceState) hint() string {

→ Read internal/engine/cache.go [limit=20, offset=370]
→ Read internal/engine/data.go [limit=80, offset=360]

$ rg -n 'func Test' internal/engine/serverdata_test.go internal/engine/devicedata_test.go
internal/engine/devicedata_test.go:66:func TestPublishDeviceDatastream(t *testing.T) {
internal/engine/devicedata_test.go:93:func TestPublishDeviceProperty(t *testing.T) {
internal/engine/devicedata_test.go:121:func TestPublishDeviceValueErrors(t *testing.T) {
internal/engine/devicedata_test.go:192:func TestPublishDeviceValueFiresBusAndTriggers(t *testing.T) {
internal/engine/serverdata_test.go:41:func TestPublishServerProperty(t *testing.T) {
internal/engine/serverdata_test.go:81:func TestPublishServerPropertyJSONHint(t *testing.T) {
internal/engine/serverdata_test.go:109:func TestPublishServerDatastream(t *testing.T) {
internal/engine/serverdata_test.go:154:func TestPublishServerValueErrors(t *testing.T) {
internal/engine/serverdata_test.go:213:func TestUnsetServerProperty(t *testing.T) {
internal/engine/serverdata_test.go:270:func TestServerObjectAggregate(t *testing.T) {
internal/engine/serverdata_test.go:336:func TestServerPublishFiresBusAndTriggers(t *testing.T) {

→ Read internal/engine/serverdata_test.go [limit=120, offset=154]
