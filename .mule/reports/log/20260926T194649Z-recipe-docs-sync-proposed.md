slug: recipe-docs-sync
verdict: proposed
at:  6fc2103
ran: 2026-09-26T19:46:49Z on DietPi in 1112s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls /root/astrate-mule/internal/ 2>/dev/null; echo "---API---"; ls /root/astrate-mule/docs/api/ 2>/dev/null; echo "---SITE---"; ls /root/astrate-mule/docs/site/ 2>/dev/null
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
---API---
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
---SITE---
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

$ wc -l /root/astrate-mule/.mule/todo.md; echo "---"; tail -60 /root/astrate-mule/.mule/todo.md
235 /root/astrate-mule/.mule/todo.md
---
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
- [x] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion Go — the device-delete leg added by docs-sync-rm-delete-device-async-operation-param (register a device, `DELETE /devices/{id}?async_operation=…` answers 204 and the row is gone on both values) compiles clean under `go vet -tags integration ./internal/realm/` but has never been executed, because the Pi has no database; until it runs, TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go is the only assertion of that leg that has actually been verified.
- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] forward-static-header-validation [auto]: `New` (internal/engine/forward/http.go:38-71) validates the URL and the method at boot precisely because of the rule its own comment states at 42-45 — "an unusable endpoint must fail here rather than surface once per delivery" — but `cfg.StaticHeaders` is stored unvalidated (http.go:68) and applied per delivery with `req.Header.Set` (http.go:110-112), and net/http checks header names and values only at write time. Measured on the real net/http path with a throwaway program: a name of `X Bad Name` or `X:Foo` fails per delivery with `net/http: invalid header field name`, and a value containing a newline fails with `net/http: invalid header field value for "X-Foo"` (a plain space in a value is legal). So one typo in `triggers.forward.static_headers` makes the process boot clean, accept triggers, and then count `astrate_engine_trigger_deliveries_total{outcome="failed"}` and log "custom trigger action failed" (internal/engine/triggers/actions.go:499) on every custom action forever. Validate in `New`: the RFC 7230 field-name grammar is exactly the token grammar `validMethod` (http.go:77-95) already implements, so reuse that predicate per name (rename it if you like) and add a CR/LF rejection per value. Add a table case to `TestNewRejectsBadConfig` (internal/engine/forward/http_test.go:271) asserting `New` returns an error for a bad name and a bad value and that no request is attempted; this test fails today because `New` returns nil error for both.
- [x] forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
- [x] forward-status-error-body [auto]: a non-2xx forward returns `forward: status %d` and nothing else (internal/engine/forward/http.go:121-123) after copying the **entire** response body into `io.Discard` with no limit (http.go:117-120) — the one thing an operator needs to debug a 500 from the bus is destroyed, and the sibling webhook request in the same codebase does bound its drain with `io.LimitReader(resp.Body, 1<<20)` (internal/engine/triggers/actions.go:626), so the two near-identical request paths have drifted on exactly that line. Read a bounded prefix (e.g. 512 bytes) of the body into the returned error and bound the discard the way actions.go:626 does. Add a case beside `TestStatusTable` (internal/engine/forward/http_test.go:164) where the handler answers 500 with body `{"error":"boom"}` and assert the error mentions both `500` and `boom`; it fails today because the error carries only the status code.
- [x] forward-envelope-bytes-test [auto]: `marshalEnvelope` (internal/engine/forward/envelope.go:19-29) has no direct test — `rg marshalEnvelope` finds exactly two callers (http.go:99, nats.go:45) and no test file — while envelope.go:6-8 makes the strong claim that "every Forwarder implementation in this package must produce byte-identical envelopes for the same inputs". Both suites check the shape only by unmarshalling into `bodyShape` (http_test.go:24-29), which cannot see key order or the exact bytes and would pass unchanged if a field were renamed on both sides. Add internal/engine/forward/envelope_test.go with a table pinning the exact output for (nil, nil) → `{"realm":"r","trigger":"t","action":null,"event":null}`, for the empty-non-nil pair the rule at 16-18 is about, and for a valued pair; container-free, no HTTP server needed.
- [x] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.

$ rg -n 'docs-sync' /root/astrate-mule/.mule/todo.md -o | wc -l; echo "--- slugs ---"; rg -o 'docs-sync-[a-z0-9-]+' /root/astrate-mule/.mule/todo.md | sort | uniq -c | sort -rn
55
--- slugs ---
      2 docs-sync-rm-delete-device-async-operation-param
      2 docs-sync-rm-async-operation-param
      2 docs-sync-native-compat-version-endpoints
      2 docs-sync-hk-validation-example
      2 docs-sync-appengine-data-output-params
      1 docs-sync-rm-version-example
      1 docs-sync-rm-validation-example-prefix
      1 docs-sync-rm-triggers-422-nested-envelope
      1 docs-sync-rm-put-interface-409
      1 docs-sync-rm-put-auth-422
      1 docs-sync-rm-policies-delete-422
      1 docs-sync-rm-mapping-required-encrypted
      1 docs-sync-rm-interfaces-detailed-param
      1 docs-sync-rm-interface-422-shapes
      1 docs-sync-rm-delete-interface-status
      1 docs-sync-rm-datastream-retention-endpoint
      1 docs-sync-rm-auth-403
      1 docs-sync-pairing-version-value
      1 docs-sync-pairing-version-endpoint
      1 docs-sync-pairing-version-404-unreachable
      1 docs-sync-pairing-status-enum
      1 docs-sync-pairing-register-404
      1 docs-sync-pairing-health-path
      1 docs-sync-pairing-deviceendpoints-dead-404-403
      1 docs-sync-native-version-value
      1 docs-sync-native-socket-query-token-auth
      1 docs-sync-native-socket-missing-403-500
      1 docs-sync-native-metrics-example-fake-series
      1 docs-sync-native-compat-health-503
      1 docs-sync-hk-retention-zero-is-unset
      1 docs-sync-hk-retention-field
      1 docs-sync-hk-patch-endpoint
      1 docs-sync-hk-delete-gating-responses
      1 docs-sync-hk-async-operation-param
      1 docs-sync-appengine-query-params-status
      1 docs-sync-appengine-group-patch-status
      1 docs-sync-appengine-group-endpoints
      1 docs-sync-appengine-get-group-device
      1 docs-sync-appengine-downsample-min
      1 docs-sync-appengine-device-status-schema
      1 docs-sync-appengine-data-set-422
      1 docs-sync-appengine-data-422-interface-level
      1 docs-sync-appengine-by-alias-endpoints
      1 docs-sync-ae-write-value-422
      1 docs-sync-ae-write-405
      1 docs-sync-ae-read-query-params
      1 docs-sync-ae-post-groups-409
      1 docs-sync-ae-post-group-devices-409
      1 docs-sync-ae-patch-merge-patch-media-type

$ cd /root/astrate-mule && git log --oneline -15
6fc2103 mule: log webhook-static-headers-override-test
f9db0a6 mule: webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.
fe2ce8c mule: log forward-envelope-bytes-test
ed869b6 mule: forward-envelope-bytes-test [auto]: `marshalEnvelope` (internal/engine/forward/envelope.go:19-29) has no direct test — `rg marshalEnvelope` finds exactly two callers (http.go:99, nats.go:45) and no test file — while envelope.go:6-8 makes the strong claim that "every Forwarder implementation in this package must produce byte-identical envelopes for the same inputs". Both suites check the shape only by unmarshalling into `bodyShape` (http_test.go:24-29), which cannot see key order or the exact bytes and would pass unchanged if a field were renamed on both sides. Add internal/engine/forward/envelope_test.go with a table pinning the exact output for (nil, nil) → `{"realm":"r","trigger":"t","action":null,"event":null}`, for the empty-non-nil pair the rule at 16-18 is about, and for a valued pair; container-free, no HTTP server needed.
0b64cfe mule: log forward-status-error-body
915c260 mule: forward-status-error-body [auto]: a non-2xx forward returns `forward: status %d` and nothing else (internal/engine/forward/http.go:121-123) after copying the **entire** response body into `io.Discard` with no limit (http.go:117-120) — the one thing an operator needs to debug a 500 from the bus is destroyed, and the sibling webhook request in the same codebase does bound its drain with `io.LimitReader(resp.Body, 1<<20)` (internal/engine/triggers/actions.go:626), so the two near-identical request paths have drifted on exactly that line. Read a bounded prefix (e.g. 512 bytes) of the body into the returned error and bound the discard the way actions.go:626 does. Add a case beside `TestStatusTable` (internal/engine/forward/http_test.go:164) where the handler answers 500 with body `{"error":"boom"}` and assert the error mentions both `500` and `boom`; it fails today because the error carries only the status code.
0640a14 mule: log forward-static-headers-override-test
3a0f153 mule: forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
13c6094 mule: log forward-static-header-validation
1cbebd5 mule: forward-static-header-validation [auto]: `New` (internal/engine/forward/http.go:38-71) validates the URL and the method at boot precisely because of the rule its own comment states at 42-45 — "an unusable endpoint must fail here rather than surface once per delivery" — but `cfg.StaticHeaders` is stored unvalidated (http.go:68) and applied per delivery with `req.Header.Set` (http.go:110-112), and net/http checks header names and values only at write time. Measured on the real net/http path with a throwaway program: a name of `X Bad Name` or `X:Foo` fails per delivery with `net/http: invalid header field name`, and a value containing a newline fails with `net/http: invalid header field value for "X-Foo"` (a plain space in a value is legal). So one typo in `triggers.forward.static_headers` makes the process boot clean, accept triggers, and then count `astrate_engine_trigger_deliveries_total{outcome="failed"}` and log "custom trigger action failed" (internal/engine/triggers/actions.go:499) on every custom action forever. Validate in `New`: the RFC 7230 field-name grammar is exactly the token grammar `validMethod` (http.go:77-95) already implements, so reuse that predicate per name (rename it if you like) and add a CR/LF rejection per value. Add a table case to `TestNewRejectsBadConfig` (internal/engine/forward/http_test.go:271) asserting `New` returns an error for a bad name and a bad value and that no request is attempted; this test fails today because `New` returns nil error for both.
0a1b425 mule: refill the queue
722d308 mule: recipe code-review ran (342s)
6764463 mule: recipe astarte-upstream timed out
f0ac79c mule: refill the queue
f26d44c mule: recipe github-issues ran (247s)

$ cd /root/astrate-mule && ls .mule/ && ls .mule/recipes/
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
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

$ cd /root/astrate-mule && git log --oneline -60 | rg -o 'docs-sync-[a-z0-9-]+|^\w+' | rg 'docs-sync' | head -30
docs-sync-rm-delete-device-async-operation-param
docs-sync-rm-delete-device-async-operation-param
docs-sync-rm-delete-device-async-operation-param
docs-sync-rm-async-operation-param
docs-sync-rm-async-operation-param
docs-sync-rm-async-operation-param
docs-sync-hk-retention-zero-is-unset
docs-sync-hk-retention-zero-is-unset
docs-sync-hk-async-operation-param
docs-sync-hk-async-operation-param

$ cd /root/astrate-mule && rg -n '^  /' docs/api/astrate_appengine_api.yaml
rg: docs/api/astrate_appengine_api.yaml: IO error for operation on docs/api/astrate_appengine_api.yaml: No such file or directory (os error 2)

$ cd /root/astrate-mule && rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/appengine/ --glob '*.go' | rg -v '_test' | head -80
internal/appengine/http.go:209:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:230:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:521:	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
internal/appengine/stream/ws.go:152:	return r.Header.Get("Accept") == "text/event-stream" || r.URL.Query().Get("transport") == "sse"

$ cd /root/astrate-mule && rg -n 'HandleFunc|mux\.|Mount' internal/appengine/http.go | head -60
38:// Mount registers the routes on mux (paths wire-frozen to upstream
40:func (a *API) Mount(mux *http.ServeMux) {
43:	mux.Handle("GET "+base+"/devices", h(a.listDevices))
44:	mux.Handle("GET "+base+"/stats/devices", h(a.devicesStats))
45:	mux.Handle("GET "+base+"/devices/{device}", h(a.getDevice))
46:	mux.Handle("PATCH "+base+"/devices/{device}", h(a.patchDevice))
47:	mux.Handle("GET "+base+"/devices-by-alias/{alias}", h(a.getDeviceByAlias))
48:	mux.Handle("GET "+base+"/devices/{device}/interfaces", h(a.listDeviceInterfaces))
50:	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}", h(a.getData))
51:	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.getData))
52:	mux.Handle("PUT "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
53:	mux.Handle("POST "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
54:	mux.Handle("DELETE "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteData))
59:	mux.Handle("PATCH "+base+"/devices-by-alias/{alias}", h(a.patchDeviceByAlias))
60:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces", h(a.listInterfacesByAlias))
61:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces/{interface}", h(a.getDataByAlias))
62:	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.getDataByAlias))
63:	mux.Handle("PUT "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.putDataByAlias))
64:	mux.Handle("POST "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.putDataByAlias))
65:	mux.Handle("DELETE "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.deleteDataByAlias))
67:	mux.Handle("GET "+base+"/groups", h(a.listGroups))
68:	mux.Handle("POST "+base+"/groups", h(a.createGroup))
69:	mux.Handle("GET "+base+"/groups/{group}", h(a.getGroup))
70:	mux.Handle("GET "+base+"/groups/{group}/devices", h(a.listGroupDevices))
71:	mux.Handle("POST "+base+"/groups/{group}/devices", h(a.addGroupDevice))
72:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}", h(a.getGroupDevice))
73:	mux.Handle("PATCH "+base+"/groups/{group}/devices/{device}", h(a.patchGroupDevice))
74:	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}", h(a.removeGroupDevice))
75:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces", h(a.listInterfacesInGroup))
76:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}", h(a.getDataInGroup))
77:	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.getDataInGroup))
78:	mux.Handle("PUT "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.putDataInGroup))
79:	mux.Handle("POST "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.putDataInGroup))
80:	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteDataInGroup))

$ cd /root/astrate-mule && rg -n '^  /' docs/api/astarte_appengine_api.yaml
26:  /appengine/v1/{realm}/devices:
98:  /appengine/v1/{realm}/stats/devices:
136:  /appengine/v1/{realm}/devices/{device}:
210:  /appengine/v1/{realm}/devices-by-alias/{alias}:
292:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces:
324:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
369:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
555:  /appengine/v1/{realm}/devices/{device}/interfaces:
583:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
623:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
789:  /appengine/v1/{realm}/groups:
863:  /appengine/v1/{realm}/groups/{group}:
891:  /appengine/v1/{realm}/groups/{group}/devices:
981:  /appengine/v1/{realm}/groups/{group}/devices/{device}:
1080:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces:
1108:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}:
1149:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}:

✗ Read docs/api/astrate_appengine_api.yaml failed [offset=1, limit=330]
Error: File not found: /root/astrate-mule/docs/api/astrate_appengine_api.yaml

$ cd /root/astrate-mule && pwd && realpath docs/api/astrate_appengine_api.yaml && ls -la docs/api/
/root/astrate-mule
/root/astrate-mule/docs/api/astrate_appengine_api.yaml
total 164
drwxr-xr-x 2 root root  4096 Sep  4 18:55 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 56583 Sep 12 22:44 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 15547 Sep 25 21:48 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 19506 Sep 17 13:03 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 45387 Sep 25 22:07 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 14977 Sep 20 21:25 astrate_native_api.yaml

→ Read docs/api/astarte_appengine_api.yaml [limit=140, offset=1]
→ Read internal/appengine/http.go [limit=120, offset=30]

$ cd /root/astrate-mule && rg -n 'operationId:|^  /appengine' docs/api/astrate_appengine_api.yaml
rg: docs/api/astrate_appengine_api.yaml: IO error for operation on docs/api/astrate_appengine_api.yaml: No such file or directory (os error 2)

$ cd /root/astrate-mule && rg -n 'operationId:|^  /appengine' ./docs/api/astrate_appengine_api.yaml
rg: ./docs/api/astrate_appengine_api.yaml: IO error for operation on ./docs/api/astrate_appengine_api.yaml: No such file or directory (os error 2)

$ cd /root/astrate-mule && ls docs/api/ && rg -n 'operationId:' docs/api/astrate_appengine_api.yaml | head -40
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
rg: docs/api/astrate_appengine_api.yaml: IO error for operation on docs/api/astrate_appengine_api.yaml: No such file or directory (os error 2)

$ cd /root/astrate-mule/docs/api && ls | cat -A | head
astarte_appengine_api.yaml$
astarte_housekeeping_api.yaml$
astarte_pairing_api.yaml$
astarte_realm_management_api.yaml$
astrate_native_api.yaml$

$ cd /root/astrate-mule && rg -n 'operationId:' docs/api/astrate_appengine_api.yaml | head -50
rg: docs/api/astrate_appengine_api.yaml: IO error for operation on docs/api/astrate_appengine_api.yaml: No such file or directory (os error 2)

$ cd /root/astrate-mule && stat docs/api/astrate_appengine_api.yaml; rg --version | head -1
stat: cannot statx 'docs/api/astrate_appengine_api.yaml': No such file or directory
ripgrep 14.1.1

$ cd /root/astrate-mule && ls -la docs/ && python3 -c "
import os
p='docs/api/astarte_appengine_api.yaml'
print(os.path.exists(p), os.path.isfile(p))
print(os.listdir('docs/api'))
"
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
drwxr-xr-x  4 root root  4096 Sep 25 22:44 site
drwxr-xr-x  2 root root  4096 Jul 27 05:40 swagger-ui
True True
['astarte_realm_management_api.yaml', 'astrate_native_api.yaml', 'astarte_appengine_api.yaml', 'astarte_housekeeping_api.yaml', 'astarte_pairing_api.yaml']

$ cd /root/astrate-mule && stat docs/api/astrate_appengine_api.yaml && rg -c 'operationId:' docs/api/astarte_appengine_api.yaml
stat: cannot statx 'docs/api/astrate_appengine_api.yaml': No such file or directory

$ cd /root/astrate-mule && python3 - <<'EOF'
import os
p='docs/api/astrate_appengine_api.yaml'
print(os.stat(p))
print(open(p).readline())
EOF
Traceback (most recent call last):
  File "<stdin>", line 3, in <module>
FileNotFoundError: [Errno 2] No such file or directory: 'docs/api/astrate_appengine_api.yaml'

$ cd /root/astrate-mule && for i in 1 2 3 4 5; do python3 -c "import os;print(os.path.exists('docs/api/astarte_appengine_api.yaml'))"; done; mount | rg -i 'astrate|overlay|fuse' | head
True
True
True
True
True
fusectl on /sys/fs/fuse/connections type fusectl (rw,nosuid,nodev,noexec,relatime)

$ cd /root/astrate-mule && python3 -c "
import os
for f in os.listdir('docs/api'): print(repr(f))
"
'astarte_realm_management_api.yaml'
'astrate_native_api.yaml'
'astarte_appengine_api.yaml'
'astarte_housekeeping_api.yaml'
'astarte_pairing_api.yaml'

$ cd /root/astrate-mule && python3 -c "
import os
p='docs/api/astarte_appengine_api.yaml'
for i in range(5):
    try:
        print(i, os.stat(p).st_size)
    except Exception as e:
        print(i, 'ERR', e)
"
0 56583
1 56583
2 56583
3 56583
4 56583

$ cd /root/astrate-mule && rg -n 'operationId:' docs/api/astrate_appengine_api.yaml
rg: docs/api/astrate_appengine_api.yaml: IO error for operation on docs/api/astrate_appengine_api.yaml: No such file or directory (os error 2)

$ cd /root/astrate-mule && cp docs/api/astrate_appengine_api.yaml /tmp/opencode/ae.yaml && wc -l /tmp/opencode/ae.yaml
cp: cannot stat 'docs/api/astrate_appengine_api.yaml': No such file or directory

$ cd /root/astrate-mule && python3 -c "
import os
d=os.listdir('docs/api')
for f in d:
    print(f, [hex(ord(c)) for c in f][:5], len(f))
" ; echo "--- try cat ---"; cat docs/api/astrate_appengine_api.yaml | head -2
astarte_realm_management_api.yaml ['0x61', '0x73', '0x74', '0x61', '0x72'] 33
astrate_native_api.yaml ['0x61', '0x73', '0x74', '0x72', '0x61'] 23
astarte_appengine_api.yaml ['0x61', '0x73', '0x74', '0x61', '0x72'] 26
astarte_housekeeping_api.yaml ['0x61', '0x73', '0x74', '0x61', '0x72'] 29
astarte_pairing_api.yaml ['0x61', '0x73', '0x74', '0x61', '0x72'] 24
--- try cat ---
cat: docs/api/astrate_appengine_api.yaml: No such file or directory

$ cd /root/astrate-mule && for f in docs/api/*.yaml; do echo -n "$f: "; head -c 20 "$f" 2>&1 | head -1; echo; done
docs/api/astarte_appengine_api.yaml: openapi: "3.0.3"

docs/api/astarte_housekeeping_api.yaml: openapi: "3.0.3"

docs/api/astarte_pairing_api.yaml: openapi: "3.0.3"

docs/api/astarte_realm_management_api.yaml: openapi: "3.0.3"

docs/api/astrate_native_api.yaml: openapi: "3.0.3"

$ cd /root/astrate-mule && rg -n 'operationId:' docs/api/astrate_appengine_api.yaml
rg: docs/api/astrate_appengine_api.yaml: IO error for operation on docs/api/astrate_appengine_api.yaml: No such file or directory (os error 2)

✱ Grep "operationId:" in docs/api/astrate_appengine_api.yaml · 76 matches
✗ Read docs/api/astrate_appengine_api.yaml failed [offset=136, limit=160]
Error: File not found: /root/astrate-mule/docs/api/astrate_appengine_api.yaml
✗ Read docs/api/astrate_appengine_api.yaml failed [limit=160, offset=136]
Error: File not found: /root/astrate-mule/docs/api/astrate_appengine_api.yaml
→ Read docs/api/astarte_appengine_api.yaml [limit=80, offset=136]

$ cd /root/astrate-mule && rg -n 'appengine' .mule/todo.md | rg -i '403' 
136:- [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
198:- [x] channels-rejoin-joinref-tagging-test: pin in internal/appengine/channels/ws_test.go that a second `phx_join` on the same topic updates the session's `joinRef` (ws.go:284-291) and subsequent `new_event` pushes carry the fresh ref (read per event at ws.go:403) — no test drives a rejoin today (all `TestWireSession_*` join once), so the comment's "server pushes must be tagged with the current one" rots unasserted. Wire: join with ref "1", rejoin with "2", push an event, assert element 0 is `"2"`.
202:- [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cd /root/astrate-mule && rg -n 'Forbidden' docs/api/astrate_appengine_api.yaml
rg: docs/api/astrate_appengine_api.yaml: IO error for operation on docs/api/astrate_appengine_api.yaml: No such file or directory (os error 2)

✱ Grep "Forbidden|"403"" in docs/api/astarte_appengine_api.yaml · 72 matches
→ Read internal/auth/middleware.go

$ cd /root/astrate-mule && python3 - <<'PY'
import re
src=open('docs/api/astarte_appengine_api.yaml').read()
i=src.index('\ncomponents:')
comp=src[i:]
# list response component names
for m in re.finditer(r'^    ([A-Za-z0-9_]+):\n', comp, re.M):
    print(m.group(1), end='  ')
print()
PY
a_aea  RealmName  DeviceID  InterfaceName  GroupName  DataSince  DataSinceAfter  DataTo  DataLimit  DataDownsample  DataSort  DataFormat  DataAllowBigIntegers  DataAllowSafeBigIntegers  DeviceStatus  Introspection  DevicePatch  GroupCreate  GroupInfo  DataEnvelopeGroupInfo  DataEnvelopeDeviceStatus  DataEnvelopeDeviceIDList  DataEnvelopeDeviceStatusList  DataEnvelopeObject  DataEnvelopeStringArray  Links  ErrorDetail  FieldErrorsDetail  BadRequest  Unauthorized  NotFound  DeviceNotFound  MethodNotAllowed  InternalServerError  ValidationErrors  ValueTooLarge  UnprocessableEntity  ConflictAliasInUse  ConflictGroupAlreadyExists  ConflictDeviceAlreadyInGroup

→ Read internal/appengine/http.go [limit=240, offset=150]

$ cd /root/astrate-mule && rg -n '403|Forbidden' .mule/todo.md | rg -i 'appengine|ae-' 
136:- [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
198:- [x] channels-rejoin-joinref-tagging-test: pin in internal/appengine/channels/ws_test.go that a second `phx_join` on the same topic updates the session's `joinRef` (ws.go:284-291) and subsequent `new_event` pushes carry the fresh ref (read per event at ws.go:403) — no test drives a rejoin today (all `TestWireSession_*` join once), so the comment's "server pushes must be tagged with the current one" rots unasserted. Wire: join with ref "1", rejoin with "2", push an event, assert element 0 is `"2"`.
202:- [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

→ Read internal/appengine/http.go [limit=180, offset=389]
→ Read docs/api/astarte_appengine_api.yaml [limit=300, offset=786]
→ Read internal/appengine/http.go [limit=130, offset=568]
→ Read docs/api/astarte_appengine_api.yaml [limit=130, offset=655]

$ cd /root/astrate-mule && rg -n 'func \(s \*Service\) (UnsetProperty|PublishData|GetData|ListGroupDevices|CreateGroup|AddGroupDevice|PatchDevice|PatchGroupDevice|PatchDeviceByAlias|GetGroup|RemoveGroupDevice|GetGroupDevice)\b' internal/appengine/*.go
internal/appengine/data.go:109:func (s *Service) GetData(ctx context.Context, realm, deviceID, ifaceName, path string, opts QueryOpts) (any, error) {
internal/appengine/data.go:466:func (s *Service) PublishData(ctx context.Context, realm, deviceID, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
internal/appengine/data.go:479:func (s *Service) UnsetProperty(ctx context.Context, realm, deviceID, ifaceName, path string) error {
internal/appengine/service.go:291:func (s *Service) PatchDevice(ctx context.Context, realm, deviceID string, p DevicePatch) (*DeviceStatus, error) {
internal/appengine/service.go:309:func (s *Service) PatchDeviceByAlias(ctx context.Context, realm, alias string, p DevicePatch) (*DeviceStatus, error) {
internal/appengine/service.go:384:func (s *Service) PatchGroupDevice(ctx context.Context, realm, groupName, deviceID string, p DevicePatch) (*DeviceStatus, error) {
internal/appengine/service.go:394:func (s *Service) GetGroupDevice(ctx context.Context, realm, groupName, deviceID string) (*DeviceStatus, error) {
internal/appengine/service.go:462:func (s *Service) CreateGroup(ctx context.Context, realm, name string, devices []string) error {
internal/appengine/service.go:508:func (s *Service) GetGroup(ctx context.Context, realm, name string) error {
internal/appengine/service.go:541:func (s *Service) ListGroupDevices(ctx context.Context, realm, name string, details bool, fromToken string, limit int) (*DevicePage, error) {
internal/appengine/service.go:626:func (s *Service) AddGroupDevice(ctx context.Context, realm, name, deviceID string) error {
internal/appengine/service.go:645:func (s *Service) RemoveGroupDevice(ctx context.Context, realm, name, deviceID string) error {

$ cd /root/astrate-mule && python3 - <<'PY'
import re
src=open('docs/api/astrate_appengine_api.yaml').read()
# count operations
ops=re.findall(r'^      operationId: (\w+)', src, re.M)
print(len(ops))
# for each operation, gather documented status codes
blocks=re.split(r'\n(?=  /appengine)', src)
tot=0; no403=[]
for b in blocks:
    m=re.match(r'  (/appengine\S*):', b)
    if not m: continue
    path=m.group(1)
    # split by operation
    for om in re.finditer(r'\n    (get|post|put|patch|delete):\n(.*?)(?=\n    (get|post|put|patch|delete):\n|\n  #|\Z)', b, re.S):
        op=om.group(1); body=om.group(2)
        oi=re.search(r'operationId: (\w+)', body)
        codes=re.findall(r'^        "(\d{3})":', body, re.M)
        tot+=1
        if '403' not in codes:
            no403.append((path,op,oi.group(1) if oi else '?',codes))
print('total ops', tot, 'without 403', len(no403))
for x in no403: print(x)
PY
Traceback (most recent call last):
  File "<stdin>", line 2, in <module>
FileNotFoundError: [Errno 2] No such file or directory: 'docs/api/astrate_appengine_api.yaml'

$ cd /root/astrate-mule && python3 - <<'PY'
import re,time
for _ in range(5):
    try:
        src=open('docs/api/astarte_appengine_api.yaml').read(); break
    except Exception as e:
        time.sleep(0.2)
else:
    raise SystemExit('unreadable')
ops=re.findall(r'^      operationId: (\w+)', src, re.M)
print('operationIds', len(ops))
blocks=re.split(r'\n(?=  /appengine)', src)
tot=0; no403=[]
for b in blocks:
    m=re.match(r'  (/appengine\S*):', b)
    if not m: continue
    path=m.group(1)
    for om in re.finditer(r'\n    (get|post|put|patch|delete):\n(.*?)(?=\n    (get|post|put|patch|delete):\n|\n  #|\Z)', b, re.S):
        op=om.group(1); body=om.group(2)
        oi=re.search(r'operationId: (\w+)', body)
        codes=re.findall(r'^        "(\d{3})":', body, re.M)
        tot+=1
        if '403' not in codes:
            no403.append((path,op,oi.group(1) if oi else '?',codes))
print('total ops', tot, 'without 403', len(no403))
for x in no403: print(x)
PY
operationIds 32
total ops 32 without 403 32
('/appengine/v1/{realm}/devices', 'get', 'listDevices', ['200', '401', '500'])
('/appengine/v1/{realm}/stats/devices', 'get', 'getDeviceStats', ['200', '401', '500'])
('/appengine/v1/{realm}/devices/{device}', 'get', 'getDevice', ['200', '401', '404', '500'])
('/appengine/v1/{realm}/devices/{device}', 'patch', 'patchDevice', ['200', '400', '401', '404', '409', '422', '500'])
('/appengine/v1/{realm}/devices-by-alias/{alias}', 'get', 'getDeviceByAlias', ['200', '401', '404', '500'])
('/appengine/v1/{realm}/devices-by-alias/{alias}', 'patch', 'patchDeviceByAlias', ['200', '400', '401', '404', '422', '500'])
('/appengine/v1/{realm}/devices-by-alias/{alias}/interfaces', 'get', 'listInterfacesByAlias', ['200', '401', '404', '500'])
('/appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}', 'get', 'getDataByAlias', ['200', '401', '404', '422', '500'])
('/appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}', 'get', 'getDataAtPathByAlias', ['200', '401', '404', '422', '500'])
('/appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}', 'put', 'putDataByAlias', ['200', '400', '401', '404', '405', '422', '500'])
('/appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}', 'post', 'publishDataByAlias', ['200', '400', '401', '404', '405', '422', '500'])
('/appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}', 'delete', 'deleteDataByAlias', ['204', '401', '404', '405', '500'])
('/appengine/v1/{realm}/devices/{device}/interfaces', 'get', 'listDeviceInterfaces', ['200', '401', '404', '500'])
('/appengine/v1/{realm}/devices/{device}/interfaces/{interface}', 'get', 'getData', ['200', '401', '404', '422', '500'])
('/appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}', 'get', 'getDataAtPath', ['200', '401', '404', '422', '500'])
('/appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}', 'put', 'putData', ['200', '400', '401', '404', '405', '422', '500'])
('/appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}', 'post', 'publishData', ['200', '400', '401', '404', '405', '422', '500'])
('/appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}', 'delete', 'deleteData', ['204', '401', '404', '405', '500'])
('/appengine/v1/{realm}/groups', 'get', 'listGroups', ['200', '401', '500'])
('/appengine/v1/{realm}/groups', 'post', 'createGroup', ['201', '400', '401', '409', '422', '500'])
('/appengine/v1/{realm}/groups/{group}', 'get', 'getGroup', ['200', '401', '404', '500'])
('/appengine/v1/{realm}/groups/{group}/devices', 'get', 'listGroupDevices', ['200', '401', '404', '422', '500'])
('/appengine/v1/{realm}/groups/{group}/devices', 'post', 'addGroupDevice', ['201', '400', '401', '404', '409', '500'])
('/appengine/v1/{realm}/groups/{group}/devices/{device}', 'get', 'getGroupDevice', ['200', '401', '404', '500'])
('/appengine/v1/{realm}/groups/{group}/devices/{device}', 'patch', 'patchGroupDevice', ['200', '400', '401', '404', '409', '422', '500'])
('/appengine/v1/{realm}/groups/{group}/devices/{device}', 'delete', 'removeGroupDevice', ['204', '401', '404', '500'])
('/appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces', 'get', 'listInterfacesInGroup', ['200', '401', '404', '500'])
('/appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}', 'get', 'getDataInGroup', ['200', '401', '404', '422', '500'])
('/appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}', 'get', 'getDataAtPathInGroup', ['200', '401', '404', '422', '500'])
('/appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}', 'put', 'putDataInGroup', ['200', '400', '401', '404', '405', '422', '500'])
('/appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}', 'post', 'publishDataInGroup', ['200', '400', '401', '404', '405', '422', '500'])
('/appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}', 'delete', 'deleteDataInGroup', ['204', '401', '404', '405', '500'])

→ Read internal/appengine/service.go [limit=140, offset=280]
→ Read docs/api/astarte_appengine_api.yaml [limit=120, offset=236]

$ cd /root/astrate-mule && python3 - <<'PY'
import re,time
for _ in range(5):
    try:
        src=open('docs/api/astarte_appengine_api.yaml').read(); break
    except Exception: time.sleep(0.2)
i=src.index('  responses:')
print(src[i:i+6000])
PY
responses:
        "200":
          description: Paginated list of devices.
          content:
            application/json; charset=utf-8:
              schema:
                oneOf:
                  - $ref: "#/components/schemas/DataEnvelopeDeviceIDList"
                  - $ref: "#/components/schemas/DataEnvelopeDeviceStatusList"
              examples:
                ids:
                  summary: Simple device ID list
                  value:
                    data:
                      - dT6hS2W9TT6LEnP25ks_lg
                    links:
                      self: "/v1/test/devices?details=false&limit=1"
                details:
                  summary: Device status list with pagination
                  value:
                    data:
                      - id: dT6hS2W9TT6LEnP25ks_lg
                        aliases: {}
                        attributes: {}
                        connected: true
                        last_connection: "2024-05-30 13:49:57.045Z"
                        first_registration: "2024-01-15 10:00:00.000Z"
                        last_disconnection: null
                        credentials_inhibited: false
                        groups: []
                        total_received_msgs: 12
                        total_received_bytes: 4096
                        introspection:
                          org.example.Sensors:
                            major: 0
                            minor: 1
                    links:
                      self: "/v1/test/devices?details=true&limit=1"
                      next: "/v1/test/devices?details=true&from_token=dT6hS2W9TT6LEnP25ks_lg&limit=1"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "500":
          $ref: "#/components/responses/InternalServerError"

  /appengine/v1/{realm}/stats/devices:
    get:
      operationId: getDeviceStats
      summary: Get device statistics
      description: Returns total and connected device counts for the realm.
      tags: [AppEngine - Devices]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
      responses:
        "200":
          description: Device statistics.
          content:
            application/json; charset=utf-8:
              schema:
                type: object
                required: [data]
                properties:
                  data:
                    type: object
                    required: [connected_devices, total_devices]
                    properties:
                      connected_devices:
                        type: integer
                        format: int64
                      total_devices:
                        type: integer
                        format: int64
              example:
                data:
                  connected_devices: 42
                  total_devices: 150
        "401":
          $ref: "#/components/responses/Unauthorized"
        "500":
          $ref: "#/components/responses/InternalServerError"

  /appengine/v1/{realm}/devices/{device}:
    get:
      operationId: getDevice
      summary: Get device details
      description: Returns the full device status object.
      tags: [AppEngine - Devices]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - $ref: "#/components/parameters/DeviceID"
      responses:
        "200":
          description: Device status.
          content:
            application/json; charset=utf-8:
              schema:
                $ref: "#/components/schemas/DataEnvelopeDeviceStatus"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "404":
          $ref: "#/components/responses/DeviceNotFound"
        "500":
          $ref: "#/components/responses/InternalServerError"

    patch:
      operationId: patchDevice
      summary: Update device attributes
      description: |
        Partially updates a device's aliases, attributes, or credentials
        inhibited status. Returns the updated device status.
      tags: [AppEngine - Devices]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - $ref: "#/components/parameters/DeviceID"
      requestBody:
        required: true
        content:
          application/merge-patch+json:
            schema:
              type: object
              required: [data]
              properties:
                data:
                  $ref: "#/components/schemas/DevicePatch"
            example:
              data:
                aliases:
                  serial_number: "ABC123"
                attributes:
                  firmware_version: "2.1.0"
                credentials_inhibited: false
      responses:
        "200":
          description: Updated device status.
          content:
            application/json; charset=utf-8:
              schema:
                $ref: "#/components/schemas/DataEnvelopeDeviceStatus"
        "400":
          $ref: "#/components/responses/BadRequest"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "404":
          $ref: "#/components/responses/DeviceNotFound"
        "409":
          $ref: "#/components/responses/ConflictAliasInUse"
        "422":
          $ref: "#/components/responses/UnprocessableEntity"
        "500":
          $ref: "#/components/responses/InternalServerError"

  /appengine/v1/{realm}/devices-by-alias/{alias}:
    get:
      operationId: getDeviceByAlias
      summary: Get device by alias
      description: Returns the device status for a device matching the given alias.
      tags: [AppEngine - Devices]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - name: alias
          in: path
          required: true
          description: The alias value to look up.
          schema:
            type: string
      responses:
        "20

$ cd /root/astrate-mule && python3 - <<'PY'
import re,time
for _ in range(5):
    try:
        src=open('docs/api/astarte_appengine_api.yaml').read(); break
    except Exception: time.sleep(0.2)
i=src.index('\n  responses:\n')
print(src[i:])
PY
responses:
    BadRequest:
      description: Bad Request — malformed or missing data envelope.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Bad Request

    Unauthorized:
      description: Unauthorized — missing or invalid JWT.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Unauthorized

    NotFound:
      description: Not Found — resource does not exist.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Not Found

    DeviceNotFound:
      description: Device not found — unknown device ID.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Device not found

    MethodNotAllowed:
      description: Method Not Allowed — cannot write to a device-owned or read-only resource.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Cannot write to read-only resource

    InternalServerError:
      description: Internal Server Error — unexpected failure.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Internal Server Error

    ValidationErrors:
      description: Validation Error — one or more fields failed validation.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/FieldErrorsDetail"
          example:
            errors:
              group_name:
                - can't be blank

    ValueTooLarge:
      description: Unprocessable Entity — the value is larger than the per-value size limit (64 KiB).
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Value size exceeds size limits

    UnprocessableEntity:
      description: Unprocessable Entity — attribute key not found on device.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Attribute key not found

    ConflictAliasInUse:
      description: Conflict — the alias value is already owned by another device.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Alias already in use

    ConflictGroupAlreadyExists:
      description: Conflict — a group with the given name already exists.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Group already exists

    ConflictDeviceAlreadyInGroup:
      description: Conflict — the device is already a member of the given group.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Device already in group

→ Read internal/appengine/data.go [limit=80, offset=460]

$ cd /root/astrate-mule && rg -n 'func .*UnsetServerProperty|func .*PublishServerValue' internal/engine/*.go internal/*/*.go | head
internal/engine/serverdata_test.go:154:func TestPublishServerValueErrors(t *testing.T) {
internal/engine/serverdata_test.go:213:func TestUnsetServerProperty(t *testing.T) {
internal/engine/serverdata.go:50:func (e *Engine) PublishServerValue(ctx context.Context, realm string, id deviceid.ID, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
internal/engine/serverdata.go:222:func (e *Engine) UnsetServerProperty(ctx context.Context, realm string, id deviceid.ID, ifaceName, path string) error {
internal/engine/serverdata_test.go:154:func TestPublishServerValueErrors(t *testing.T) {
internal/engine/serverdata_test.go:213:func TestUnsetServerProperty(t *testing.T) {
internal/appengine/http_test.go:46:func (f *fakeServerData) PublishServerValue(_ context.Context, realm string, id deviceid.ID, iface, path string, value json.RawMessage, _ *time.Time) error {
internal/appengine/http_test.go:53:func (f *fakeServerData) UnsetServerProperty(_ context.Context, realm string, id deviceid.ID, iface, path string) error {
internal/engine/serverdata.go:50:func (e *Engine) PublishServerValue(ctx context.Context, realm string, id deviceid.ID, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
internal/engine/serverdata.go:222:func (e *Engine) UnsetServerProperty(ctx context.Context, realm string, id deviceid.ID, ifaceName, path string) error {

→ Read internal/engine/serverdata.go [limit=120, offset=40]
→ Read internal/engine/serverdata.go [limit=80, offset=215]
→ Read internal/appengine/service.go [limit=120, offset=500]
→ Read internal/appengine/service.go [limit=130, offset=150]

$ cd /root/astrate-mule && rg -n 'type DeviceStatus struct' -A 30 internal/appengine/service.go
105:type DeviceStatus struct {
106-	ID                   string                        `json:"id"`
107-	Connected            bool                          `json:"connected"`
108-	Introspection        map[string]introspectionEntry `json:"introspection"`
109-	Aliases              map[string]string             `json:"aliases"`
110-	Attributes           map[string]string             `json:"attributes"`
111-	Groups               []string                      `json:"groups"`
112-	CredentialsInhibited bool                          `json:"credentials_inhibited"`
113-	TotalReceivedMsgs    int64                         `json:"total_received_msgs"`
114-	TotalReceivedBytes   int64                         `json:"total_received_bytes"`
115-	FirstRegistration    *time.Time                    `json:"first_registration"`
116-	FirstCredentialsReq  *time.Time                    `json:"first_credentials_request"`
117-	LastConnection       *time.Time                    `json:"last_connection"`
118-	LastDisconnection    *time.Time                    `json:"last_disconnection"`
119-	LastSeenIP           string                        `json:"last_seen_ip,omitempty"`
120-	PreviousInterfaces   map[string]introspectionEntry `json:"previous_interfaces,omitempty"`
121-}
122-
123-// introspectionEntry renders one introspection pair (upstream shape).
124-type introspectionEntry struct {
125-	Major int `json:"major"`
126-	Minor int `json:"minor"`
127-}
128-
129-// DevicePage is one device-list page: IDs, or full statuses when details was
130-// requested. Next is the next-page cursor ("" at the end).
131-type DevicePage struct {
132-	IDs      []string
133-	Statuses []*DeviceStatus
134-	Next     string
135-}

$ cd /root/astrate-mule && python3 - <<'PY'
import re,time
for _ in range(5):
    try:
        src=open('docs/api/astarte_appengine_api.yaml').read(); break
    except Exception: time.sleep(0.2)
i=src.index('    DeviceStatus:')
print(src[i:i+2600])
PY
DeviceStatus:
      type: object
      required: [id, connected]
      properties:
        id:
          type: string
        aliases:
          type: object
          additionalProperties:
            type: string
          nullable: true
        attributes:
          type: object
          additionalProperties:
            type: string
          nullable: true
        connected:
          type: boolean
        last_connection:
          type: string
          nullable: true
          description: Timestamp of last connection (Elixir DateTime format).
        first_registration:
          type: string
          nullable: true
          description: Timestamp of first registration.
        last_disconnection:
          type: string
          nullable: true
          description: Timestamp of last disconnection.
        credentials_inhibited:
          type: boolean
        groups:
          type: array
          items:
            type: string
          description: Group names the device belongs to.
        total_received_msgs:
          type: integer
          format: int64
        total_received_bytes:
          type: integer
          format: int64
        first_credentials_request:
          type: string
          nullable: true
          description: Timestamp of first credentials request.
        last_seen_ip:
          type: string
          description: Last seen IP address. Omitted when never seen.
        previous_interfaces:
          $ref: "#/components/schemas/Introspection"
          nullable: true
          description: Introspection before the last update, when present.
        introspection:
          $ref: "#/components/schemas/Introspection"
          nullable: true
          description: Map of interface name to its major and minor version.

    Introspection:
      type: object
      additionalProperties:
        type: object
        required: [major, minor]
        properties:
          major:
            type: integer
          minor:
            type: integer

    DevicePatch:
      type: object
      properties:
        aliases:
          type: object
          additionalProperties:
            type: string
            nullable: true
          description: Aliases to set. Null values remove the alias.
        attributes:
          type: object
          additionalProperties:
            type: string
            nullable: true
          description: Attributes to set. Null values remove the attribute.
        credentials_inhibited:
          type: boolean
          nullable: true
          description: Set to true to inhibit

→ Read internal/appengine/service.go [limit=50, offset=136]
→ Read internal/appengine/service.go [limit=40, offset=620]

$ cd /root/astrate-mule && rg -n 'ErrPathNotFound|ErrInvalidAlias|ErrValidation\b' internal/appengine/*.go | rg -v '_test'
internal/appengine/service.go:27:// ErrValidation wraps a well-formed request that violates a rule (maps to 422).
internal/appengine/service.go:28:var ErrValidation = errors.New("appengine: validation failed")
internal/appengine/service.go:34:	ErrInvalidAlias         = errors.New("Invalid alias")           //nolint:staticcheck // ST1005: detail strings are upstream astarte_appengine_api 1.2.2 wire text
internal/appengine/service.go:153:			return nil, fmt.Errorf("%w: invalid cursor", ErrValidation)
internal/appengine/service.go:328:			return nil, ErrInvalidAlias
internal/appengine/service.go:717:			return nil, fmt.Errorf("%w: downsample_to is not supported on properties interfaces", ErrValidation)
internal/appengine/data.go:29:// ErrPathNotFound marks an object-aggregated series whose concrete path holds
internal/appengine/data.go:32:var ErrPathNotFound = errors.New("appengine: path not found")
internal/appengine/data.go:119:			return nil, fmt.Errorf("%w: downsample_to is not supported on properties interfaces", ErrValidation)
internal/appengine/data.go:158:			return nil, fmt.Errorf("%w: downsample_to is not supported on object-aggregated interfaces", ErrValidation)
internal/appengine/data.go:166:			return nil, fmt.Errorf("%w: downsample_to requires an endpoint path", ErrValidation)
internal/appengine/data.go:218:			return nil, fmt.Errorf("%w: %s", ErrPathNotFound, path)
internal/appengine/http.go:648:	case errors.Is(err, ErrInvalidAlias):
internal/appengine/http.go:658:	case errors.Is(err, ErrValidation):
internal/appengine/http.go:668:	case errors.Is(err, ErrPathNotFound):
internal/appengine/http.go:679:	case errors.Is(err, engine.ErrPathNotFound):
internal/appengine/http.go:688:// validationDetail strips the ErrValidation prefix for the response detail.

→ Read internal/appengine/data.go [limit=130, offset=100]

$ cd /root/astrate-mule && python3 - <<'PY'
import re,time
for _ in range(5):
    try:
        src=open('docs/api/astarte_appengine_api.yaml').read(); break
    except Exception: time.sleep(0.2)
for m in re.finditer(r'operationId: (\w+)(.*?)(?=\n    (?:get|post|put|patch|delete):|\n  /|\Z)', src, re.S):
    oi=m.group(1); body=m.group(2)
    refs=re.findall(r'\$ref: "#/components/responses/(\w+)"', body)
    if 'NotFound' in refs:
        print(oi, refs)
PY
getGroup ['Unauthorized', 'NotFound', 'InternalServerError']
listGroupDevices ['Unauthorized', 'NotFound', 'ValidationErrors', 'InternalServerError']
addGroupDevice ['BadRequest', 'Unauthorized', 'NotFound', 'ConflictDeviceAlreadyInGroup', 'InternalServerError']
getGroupDevice ['Unauthorized', 'NotFound', 'InternalServerError']
patchGroupDevice ['BadRequest', 'Unauthorized', 'NotFound', 'ConflictAliasInUse', 'UnprocessableEntity', 'InternalServerError']
removeGroupDevice ['Unauthorized', 'NotFound', 'InternalServerError']
listInterfacesInGroup ['Unauthorized', 'NotFound', 'InternalServerError']
getDataInGroup ['Unauthorized', 'NotFound', 'ValidationErrors', 'InternalServerError']
getDataAtPathInGroup ['Unauthorized', 'NotFound', 'ValidationErrors', 'InternalServerError']
putDataInGroup ['BadRequest', 'Unauthorized', 'NotFound', 'MethodNotAllowed', 'ValueTooLarge', 'InternalServerError']
publishDataInGroup ['BadRequest', 'Unauthorized', 'NotFound', 'MethodNotAllowed', 'ValueTooLarge', 'InternalServerError']
deleteDataInGroup ['Unauthorized', 'NotFound', 'MethodNotAllowed', 'InternalServerError']

$ cd /root/astrate-mule && rg -n 'func Write(DeviceNotFound|Error|BadRequest|Unauthorized|Forbidden|FieldErrors|InternalServerError)|Detail[A-Za-z]* *=' pkg/astarteapi/*.go | rg -v _test
pkg/astarteapi/envelope.go:37:	DetailBadRequest = "Bad request"
pkg/astarteapi/envelope.go:39:	DetailUnauthorized = "Unauthorized"
pkg/astarteapi/envelope.go:41:	DetailForbidden = "Forbidden"
pkg/astarteapi/envelope.go:43:	DetailNotFound = "Not Found"
pkg/astarteapi/envelope.go:46:	DetailDeviceNotFound = "Device not found"
pkg/astarteapi/envelope.go:51:	DetailRouteNotFound = "Not found"
pkg/astarteapi/envelope.go:54:	DetailPageNotFound = "Page not found"
pkg/astarteapi/envelope.go:58:	DetailInternalServerError = "Internal server error"
pkg/astarteapi/envelope.go:163:func WriteError(w http.ResponseWriter, status int, detail string) error {
pkg/astarteapi/envelope.go:172:func WriteFieldErrors(w http.ResponseWriter, status int, fields map[string][]string) error {
pkg/astarteapi/envelope.go:191:func WriteBadRequest(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:196:func WriteUnauthorized(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:201:func WriteForbidden(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:211:func WriteDeviceNotFound(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:216:func WriteInternalServerError(w http.ResponseWriter) error {

$ cd /root/astrate-mule && for f in docs/api/*.yaml; do echo "== $f"; python3 - "$f" <<'PY'
import sys,re
src=open(sys.argv[1]).read()
for m in re.finditer(r'^    (\w+):\n      description: (.*)\n(?:.*\n)*?', src, re.M):
    pass
for m in re.finditer(r'\n    ([A-Z]\w+):\n      description: [^\n]*\n(?:.*\n)*?          example:\n            errors:\n              detail: (.*)\n', src):
    print(f'  {m.group(1):28} -> {m.group(2)!r}')
PY
done
== docs/api/astarte_appengine_api.yaml
  BadRequest                   -> 'Bad Request'
  Unauthorized                 -> 'Unauthorized'
  NotFound                     -> 'Not Found'
  DeviceNotFound               -> 'Device not found'
  MethodNotAllowed             -> 'Cannot write to read-only resource'
  InternalServerError          -> 'Internal Server Error'
  ValidationErrors             -> 'Value size exceeds size limits'
  UnprocessableEntity          -> 'Attribute key not found'
  ConflictAliasInUse           -> 'Alias already in use'
  ConflictGroupAlreadyExists   -> 'Group already exists'
  ConflictDeviceAlreadyInGroup -> 'Device already in group'
== docs/api/astarte_housekeeping_api.yaml
  BadRequest                   -> 'Bad Request'
  Unauthorized                 -> 'Unauthorized'
  Forbidden                    -> 'Forbidden'
  NotFound                     -> 'Not Found'
  MethodNotAllowed             -> 'Realm deletion disabled'
  ValidationErrors             -> 'Realm already exists'
  ValidationError              -> '"realm_name can\'t be blank"'
  InternalServerError          -> 'Internal Server Error'
== docs/api/astarte_pairing_api.yaml
  BadRequest                   -> 'Bad Request'
  Unauthorized                 -> 'Unauthorized'
  Forbidden                    -> 'Forbidden'
  DeviceNotFound               -> 'Device not found'
  ValidationError              -> 'Too Many Requests'
  InternalServerError          -> 'Internal Server Error'
== docs/api/astarte_realm_management_api.yaml
  BadRequest                   -> 'Bad Request'
  Unauthorized                 -> 'Unauthorized'
  Forbidden                    -> 'Forbidden'
  NotFound                     -> 'Not Found'
  DeviceNotFound               -> 'Device not found'
  Conflict                     -> '"invalid interface: missing \\"interface_name\\""'
  InternalServerError          -> 'Internal Server Error'
== docs/api/astrate_native_api.yaml
  Unauthorized                 -> 'Unauthorized'
  Forbidden                    -> 'Forbidden'
  InternalServerError          -> 'Internal Server Error'

$ cd /root/astrate-mule && rg -n 'maxValueBytes|maxBodyBytes|DefaultDeviceLimit *=' internal/appengine/*.go | rg -v _test
internal/appengine/http.go:18:// maxBodyBytes caps AppEngine request bodies.
internal/appengine/http.go:19:const maxBodyBytes int64 = 1 << 20
internal/appengine/http.go:21:// maxValueBytes is upstream's per-value limit (~64 KiB): a larger REST
internal/appengine/http.go:24:const maxValueBytes = 65536
internal/appengine/http.go:214:	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
internal/appengine/http.go:235:	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
internal/appengine/http.go:294:	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &value); err != nil {
internal/appengine/http.go:298:	if len(value) > maxValueBytes {
internal/appengine/http.go:323:	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &value); err != nil {
internal/appengine/http.go:327:	if len(value) > maxValueBytes {
internal/appengine/http.go:342:	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &value); err != nil {
internal/appengine/http.go:346:	if len(value) > maxValueBytes {
internal/appengine/http.go:400:	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
internal/appengine/http.go:488:	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
internal/appengine/http.go:526:	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
internal/appengine/service.go:25:const DefaultDeviceLimit = 100

→ Read pkg/astarteapi/envelope.go [limit=100, offset=30]
→ Read internal/appengine/data.go [limit=100, offset=1]

$ cd /root/astrate-mule && rg -n 'DetailNotFound|WriteNotFound' --glob '!*_test.go' . | head -20
./pkg/astarteapi/envelope.go:42:	// DetailNotFound is the canonical generic 404 detail.
./pkg/astarteapi/envelope.go:43:	DetailNotFound = "Not Found"
./pkg/astarteapi/envelope.go:49:	// DetailNotFound, which a handler emits for a route that matched and found
./pkg/astarteapi/envelope.go:205:// WriteNotFound writes the canonical generic 404 envelope.
./pkg/astarteapi/envelope.go:206:func WriteNotFound(w http.ResponseWriter) error {
./pkg/astarteapi/envelope.go:207:	return WriteError(w, http.StatusNotFound, DetailNotFound)
./internal/flowapi/http.go:282:		_ = astarteapi.WriteNotFound(w)
./internal/flowapi/http.go:288:		_ = astarteapi.WriteNotFound(w)
./internal/realm/http.go:341:		_ = astarteapi.WriteNotFound(w)
./internal/realm/http.go:421:		_ = astarteapi.WriteNotFound(w)
./internal/pairing/http.go:290:		_ = astarteapi.WriteNotFound(w)
./internal/housekeeping/http.go:245:		_ = astarteapi.WriteNotFound(w)

$ cd /root/astrate-mule && rg -n 'missingGroup' internal/appengine/service.go | head
64:// missingGroup wraps a failed group resolution so it satisfies BOTH
66:func missingGroup(name string) error {
412:		return 0, deviceid.ID{}, nil, missingGroup(groupName)
514:		return missingGroup(name)
561:		return nil, missingGroup(name)

→ Read internal/appengine/service.go [limit=80, offset=20]

$ cd /root/astrate-mule && rg -n 'docs-sync-(appengine|ae)-' .mule/todo.md | rg -o 'docs-sync-[a-z0-9-]+' | sort -u
docs-sync-ae-patch-merge-patch-media-type
docs-sync-ae-post-group-devices-409
docs-sync-ae-post-groups-409
docs-sync-ae-read-query-params
docs-sync-ae-write-405
docs-sync-ae-write-value-422
docs-sync-appengine-by-alias-endpoints
docs-sync-appengine-data-422-interface-level
docs-sync-appengine-data-output-params
docs-sync-appengine-data-set-422
docs-sync-appengine-device-status-schema
docs-sync-appengine-downsample-min
docs-sync-appengine-get-group-device
docs-sync-appengine-group-endpoints
docs-sync-appengine-group-patch-status
docs-sync-appengine-query-params-status

$ cd /root/astrate-mule && rg -n 'docs-sync-(appengine|ae)-[a-z0-9-]+' .mule/todo.md -o | head -40; echo "==="; rg -n '^\S.*docs-sync-appengine-(group-endpoints|group-patch-status|get-group-device|by-alias-endpoints|device-status-schema|downsample-min|query-params-status|data-set-422|data-422-interface-level)' .mule/todo.md
127:docs-sync-appengine-by-alias-endpoints
128:docs-sync-appengine-group-endpoints
129:docs-sync-appengine-get-group-device
130:docs-sync-appengine-query-params-status
132:docs-sync-appengine-group-patch-status
133:docs-sync-appengine-data-422-interface-level
136:docs-sync-appengine-data-output-params
136:docs-sync-appengine-data-output-params
143:docs-sync-appengine-device-status-schema
144:docs-sync-appengine-data-set-422
145:docs-sync-appengine-downsample-min
167:docs-sync-ae-write-405
168:docs-sync-ae-write-value-422
169:docs-sync-ae-post-groups-409
172:docs-sync-ae-post-group-devices-409
173:docs-sync-ae-patch-merge-patch-media-type
208:docs-sync-ae-read-query-params
===
127:- [x] docs-sync-appengine-by-alias-endpoints [auto]: add the 7 undocumented by-alias mirror routes to docs/api/astarte_appengine_api.yaml — PATCH /devices-by-alias/{alias}, GET /devices-by-alias/{alias}/interfaces, GET/PUT/POST/DELETE /devices-by-alias/{alias}/interfaces/{interface}/{path} (plus the interface-scoped GET) — all registered in internal/appengine/http.go:59-65 and backed by handlers patchDeviceByAlias, listInterfacesByAlias, getDataByAlias, putDataByAlias, deleteDataByAlias (http.go:226-367). Match the existing spec style (same schemas, same response codes as the device-scoped counterparts). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
128:- [x] docs-sync-appengine-group-endpoints [auto]: add the 8 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (returns {data: {group_name: "<name>"}}, http.go:428-434), PATCH /groups/{group}/devices/{device} (DevicePatch body → updated status, http.go:518-537), and the 6 group interface-data routes: the interface-level GET and the path-level GET/PUT/POST/DELETE on /groups/{group}/devices/{device}/interfaces (http.go:75-80, handlers listInterfacesInGroup, getDataInGroup, putDataInGroup, deleteDataInGroup at http.go:187-377). Also add GET /devices/{device}/interfaces (listDeviceInterfaces, http.go:48) which is missing from the spec — that is the 9th route. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
129:- [x] docs-sync-appengine-get-group-device [auto]: add the undocumented `GET /groups/{group}/devices/{device}` route to docs/api/astarte_appengine_api.yaml — it exists in code (internal/appengine/http.go:72, handler getGroupDevice at http.go:509-516, returns the device status projection behind the membership gate) but the spec path only documents delete and patch. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
130:- [x] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups/{group}/devices is missing `from_token` and `limit` pagination params that code reads at http.go:439-441 (and returns via links.next), (2) POST /groups is missing 422 validation-error response (FieldErrors for blank group_name/empty devices, http.go:404-416), (3) PATCH /devices/{device} and PATCH /devices-by-alias/{alias} are missing 422 ("Attribute key not found" via ErrAttributeKeyNotFound, service.go:352-357 → http.go:656-657), and the device-scoped PATCH is also missing 409 ("Alias already in use", http.go:651) and its documented 500 covers the merge-patch content-type mismatch which code answers with an unmapped 500 too, (4) GET /devices/{device}/interfaces/{interface}/{path} is missing 422 (invalid since/since_after/to/limit/downsample_to/format params, http.go:554-639). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
132:- [x] docs-sync-appengine-group-patch-status [auto]: PATCH /groups/{group}/devices/{device} (internal/appengine/http.go:518-537, patchGroupDevice) shares applyPatch with the device-scoped PATCH and also answers 422 ("Attribute key not found", ErrAttributeKeyNotFound, service.go:352-357) and 409 ("Alias already in use", http.go:651), but docs/api/astarte_appengine_api.yaml lists only 200/400/401/404/500 for it — add the two missing responses.
133:- [x] docs-sync-appengine-data-422-interface-level [auto]: the interface-level data GETs — GET /devices/{device}/interfaces/{interface}, GET /devices-by-alias/{alias}/interfaces/{interface}, GET /groups/{group}/devices/{device}/interfaces/{interface} — share parseQueryOpts (internal/appengine/http.go:554-639) and answer 422 FieldErrors for invalid since/since_after/to/limit/downsample_to/format params, but the spec documents only the path-level GET with 422 — add 422 to the interface-level GETs.
143:- [x] docs-sync-appengine-device-status-schema [auto]: fix the `DeviceStatus` schema in docs/api/astrate_appengine_api.yaml — the wire emits `id` (internal/appengine/service.go:106 `json:"id"`, decoded as `ds.ID` in mirror_test.go:156 and dashboard_compat_test.go:100/149) but the spec names the field `device_id` (yaml:1409) in both the schema and the list/detail examples (yaml:67, 75); `introspection` is documented as `map<string,integer>` (yaml:1438-1441) but each value is a `{major, minor}` object (introspectionEntry, service.go:123-127, introspectionView service.go:833-836); and the schema omits the six fields `groups`, `total_received_msgs`, `total_received_bytes`, `first_credentials_request`, `last_seen_ip` (omitempty), `previous_interfaces` (omitempty) that deviceStatusView emits (service.go:254-279). Rename `device_id`→`id`, correct `introspection` to `map<string,{major:int,minor:int}>`, add the missing fields, and fix the examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
144:- [!] docs-sync-appengine-data-set-422 [auto]: add the missing `422` response to the six PUT/POST data-set operations in docs/api/astrate_appengine_api.yaml — `putData`, `putDataByAlias`, `putDataInGroup` (internal/appengine/http.go:292-309, 321-338, 340-357) return 422 `{"errors":{"detail":"Value size exceeds size limits"}}` when the body value exceeds 65536 bytes (maxValueBytes, http.go:24), both for PUT and its POST alias (routes registered at http.go:52-53, 63-64, 78-79), but the spec lists only 200/400/401/404/500 for all six operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
145:- [x] docs-sync-appengine-downsample-min [auto]: fix the `downsample_to` constraint in docs/api/astrate_appengine_api.yaml — the DataDownsample description says "(must be >= 2)" (yaml:1357) but parseQueryOpts rejects `downsample_to <= 2` with "must be greater than 2" (internal/appengine/http.go:592-593), so the documented minimum is off by one: readings of `2` (allowed per the doc) actually 422. Change the description to "(must be > 2)" (the 422 error shape is already documented). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cd /root/astrate-mule && rg -n 'invalid cursor|from_token|device_id.*is not a valid|Endpoint not found|Group not found' internal/appengine/*_test.go | head -30
internal/appengine/mirror_test.go:220:	// Unknown group name maps to 404 "Group not found" on every mirrored surface.
internal/appengine/mirror_test.go:233:			if got.Code != http.StatusNotFound || !strings.Contains(got.Body.String(), "Group not found") {
internal/appengine/mirror_test.go:234:				t.Errorf("%s unknown group: got %d (%s), want 404 Group not found", x.name, got.Code, got.Body)
internal/appengine/mirror_test.go:240:	if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "Group not found") {
internal/appengine/mirror_test.go:241:		t.Errorf("patch unknown group: got %d (%s), want 404 Group not found", rec.Code, rec.Body)
internal/appengine/groups_parity_test.go:61:	if rec.Code != http.StatusNotFound || !strings.Contains(rec.Body.String(), "Group not found") {
internal/appengine/groups_parity_test.go:62:		t.Errorf("unknown group: got %d (%s), want 404 Group not found", rec.Code, rec.Body)
internal/appengine/groups_parity_test.go:240:			if got := u.Query().Get("from_token"); got == "" {
internal/appengine/groups_parity_test.go:241:				t.Errorf("links.next %q carries no from_token", nextURL)
internal/appengine/groups_parity_test.go:243:				t.Errorf("links.next from_token %q is not a UUID-v1-format token", got)
internal/appengine/groups_parity_test.go:258:	// from_token must be a v1-format uuid: junk and v4 are both rejected.
internal/appengine/groups_parity_test.go:259:	rec := r.req(t, http.MethodGet, "/groups/"+pg+"/devices?from_token=zzz", "", r.token)
internal/appengine/groups_parity_test.go:261:		t.Errorf("from_token=zzz: got %d (%s), want 422 is invalid", rec.Code, rec.Body)
internal/appengine/groups_parity_test.go:267:	rec = r.req(t, http.MethodGet, "/groups/"+pg+"/devices?from_token="+v4.UUID(), "", r.token)
internal/appengine/groups_parity_test.go:269:		t.Errorf("from_token=v4 uuid: got %d (%s), want 422 is invalid", rec.Code, rec.Body)
internal/appengine/groups_parity_test.go:271:	rec = r.req(t, http.MethodGet, "/groups/"+pg+"/devices?from_token="+groupTokenFor(0), "", r.token)
internal/appengine/groups_parity_test.go:273:		t.Errorf("from_token=groupTokenFor(0): got %d (%s), want 200", rec.Code, rec.Body)
internal/appengine/writeerror_test.go:29:			http.StatusBadRequest, "Endpoint not found"},
internal/appengine/dashboard_compat_test.go:25:// dashboardCursor extracts from_token the way the dashboard does:
internal/appengine/dashboard_compat_test.go:34:		if k, v, _ := strings.Cut(pair, "="); k == "from_token" {
internal/appengine/dashboard_compat_test.go:85:				path += "&from_token=" + url.QueryEscape(cursor)
internal/appengine/dashboard_compat_test.go:114:			if strings.Contains(body.Links.Next, "?from_token=") {
internal/appengine/dashboard_compat_test.go:115:				t.Fatalf("from_token first in %q", body.Links.Next)
internal/appengine/links_test.go:12:// the FIRST key/value pair ("/v1/…?details" becomes the first key). from_token
internal/appengine/links_test.go:24:	if strings.Contains(links.Next, "?from_token=") {
internal/appengine/links_test.go:25:		t.Fatalf("from_token is the first query parameter in %q — the dashboard's URLSearchParams parse would corrupt it", links.Next)
internal/appengine/links_test.go:36:		if k == "from_token" {
internal/appengine/links_test.go:41:		t.Errorf("dashboard-style parse recovered from_token %q, want CURSOR", got)
internal/appengine/links_test.go:55:	// "details" always precedes "from_token" even when the client sent
internal/appengine/links_test.go:58:	if links.Next != "/v1/test/devices?details=false&from_token=N&limit=5" {

$ cd /root/astrate-mule && rg -n 'NotFoundHandler|NotFoundHandler|appengine.*Mount|\.Mount\(' cmd/astrate/main.go | head -20
362:	}).Mount(mux)
371:	housekeeping.NewAPI(hkSvc, mw, hkKeys).Mount(mux)
375:	realm.NewAPI(realmSvc, mw).Mount(mux)
376:	appengine.NewAPI(appengine.NewService(st, e, log), mw).Mount(mux)
377:	apstream.NewAPI(e.Bus(), mw).Mount(mux)
380:	channels.NewAPI(e.Bus(), st, st).Mount(mux)
382:	flowapi.NewAPI(flowSvc, mw).Mount(mux)
414:	health.Mount(mux)
417:	swagger.Mount(mux)

$ cd /root/astrate-mule && rg -n 'invalid cursor|is not a valid device id' internal/ --glob '!*_test.go'
internal/engine/triggers/match.go:497:				fe.add("device_id", "is not a valid device id")
internal/appengine/service.go:153:			return nil, fmt.Errorf("%w: invalid cursor", ErrValidation)
internal/appengine/service.go:633:		return FieldErrors{"device_id": {"is not a valid device id"}}

→ Read internal/appengine/groups_parity_test.go [limit=120, offset=1]
Now let me check the config-key half of the recipe (cheap) before writing.

$ cd /root/astrate-mule && rg -o '\bASTRATE_[A-Z_0-9]+' -N --no-filename internal/ cmd/ | sort -u > /tmp/opencode/code_keys.txt; wc -l /tmp/opencode/code_keys.txt; rg -o '\bASTRATE_[A-Z_0-9]+' -N --no-filename docs/site/configuration-reference.md | sort -u > /tmp/opencode/doc_keys.txt; wc -l /tmp/opencode/doc_keys.txt; echo "=== in docs, not in code ==="; comm -23 /tmp/opencode/doc_keys.txt /tmp/opencode/code_keys.txt; echo "=== in code, not in docs ==="; comm -13 /tmp/opencode/doc_keys.txt /tmp/opencode/code_keys.txt
24 /tmp/opencode/code_keys.txt
20 /tmp/opencode/doc_keys.txt
=== in docs, not in code ===
=== in code, not in docs ===
ASTRATE_FLOW_CONFIG
ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
ASTRATE_TEST_DSN

$ cd /root/astrate-mule && rg -n 'ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED|ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION|ASTRATE_FLOW_CONFIG' --glob '!*_test.go' internal/ cmd/ | head
internal/config/config.go:279:		"ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
internal/config/config.go:297:	if v, ok := os.LookupEnv("ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"); ok {
internal/flow/blocks/info.go:116:		Config:  "image (required), config (object→ASTRATE_FLOW_CONFIG), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",
internal/flow/blocks/schema.go:19:	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as ASTRATE_FLOW_CONFIG"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
internal/flow/blocks/container/block.go:45:	Nested        map[string]any // opaque JSON for ASTRATE_FLOW_CONFIG
internal/flow/blocks/container/docker.go:19:	// FlowConfigJSON is passed as ASTRATE_FLOW_CONFIG.
internal/flow/blocks/container/docker.go:99:		args = append(args, "-e", "ASTRATE_FLOW_CONFIG="+spec.FlowConfigJSON)
internal/flow/blocks/container/docker.go:201:// encodeFlowConfigJSON marshals nested config for ASTRATE_FLOW_CONFIG.

$ cd /root/astrate-mule && rg -n 'HOUSEKEEPING_REALM_DELETION|DEFAULT_DATASTREAM' docs/site/*.md docs/COMPATIBILITY.md | head
(no output)

$ cd /root/astrate-mule && rg -n '^#{2,3} ' docs/site/configuration-reference.md
7:## Required values
16:## Precedence rule
26:## `[database]` — PostgreSQL/TimescaleDB
32:## `[http]` — REST API listener
41:## `[mqtt]` — Embedded MQTT broker
56:## `[engine]` — Ingestion pipeline
66:## `[pairing]` — Credential issuance & rate limits
78:## `[housekeeping]` — Instance-admin keys
85:## `[storage]` — Retention
91:## `[security]` — Master key
100:## `[realm]` — Auto-provision realm
109:## `[log]` — Logging
118:## Example: minimal production config
147:## Example: dev-mode docker compose

→ Read internal/config/config.go [limit=50, offset=265]
→ Read docs/site/configuration-reference.md [limit=30, offset=76]

$ cd /root/astrate-mule && rg -n 'Housekeeping struct' -A 20 internal/config/config.go
(no output)

$ cd /root/astrate-mule && rg -n 'Housekeeping' internal/config/config.go | head -20
49:	Housekeeping HousekeepingConfig `toml:"housekeeping"`
109:// HousekeepingConfig carries the instance-admin JWT public keys (a_ha): either
111:type HousekeepingConfig struct {
291:		cfg.Housekeeping.DefaultDatastreamMaximumStorageRetention = &n
300:			cfg.Housekeeping.RealmDeletionDisabled = true
302:			cfg.Housekeeping.RealmDeletionDisabled = false
392:// HousekeepingKeys resolves the instance-admin JWT public keys, reading any
394:func (c *Config) HousekeepingKeys() ([]string, error) {
395:	keys := append([]string(nil), c.Housekeeping.JWTPublicKeys...)
396:	for _, f := range c.Housekeeping.JWTPublicKeyFiles {

→ Read internal/config/config.go [limit=30, offset=105]

$ cd /root/astrate-mule && rg -n 'ASTRATE_TEST_DSN' --glob '!*_test.go' . | head -5; echo "=== toml keys in config ==="; rg -o 'toml:"([a-z_]+)"' -r '$1' internal/config/config.go | sort -u | tr '\n' ' '
./docs/ROADMAP.md:112:| 0.5 | `internal/testutil/pg.go` | `StartTimescale(t) *pgxpool.Pool`: container (or `ASTRATE_TEST_DSN` reuse for speed), waits ready, returns pool + teardown | 90 |
./docs/site/contributing.md:76:ASTRATE_TEST_DSN="postgres://user:pass@host:5432/dbname?sslmode=disable" make test-integration
./Makefile:40:## test-integration: T2 — requires a Docker daemon (or ASTRATE_TEST_DSN to reuse a database).
./docker-compose.yml:21:#   ASTRATE_TEST_DSN="postgres://astrate:astrate@127.0.0.1:5432/astrate?sslmode=disable" make test-integration
./test/conformance/instance/instance.go:6:// service (`make up`) or ASTRATE_TEST_DSN; New skips the test when neither is
=== toml keys in config ===
addr advertised_url batch_max_rows batch_max_wait bcrypt_cost cert_ttl cors_allowed_origins credentials_burst credentials_rate database default_datastream_maximum_storage_retention dev_addr device_registration_limit dsn enforce_latest_cert engine format forward housekeeping http insecure_dev_mode jwt_public_key jwt_public_key_file jwt_public_key_files jwt_public_keys kind level log master_key_file max_packet_bytes max_payload_bytes method mqtt name pairing realm realm_deletion_disabled register_burst register_rate retention security session_store_path shard_queue shards static_headers storage subject tls_cert_file tls_key_file triggers url

$ cd /root/astrate-mule && rg -n 'func \(s \*Service\) (ListInterfacesInGroup|GetDataInGroup|PutDataInGroup|DeleteDataInGroup|ListInterfacesByAlias|ListInterfaces)\b' -A 22 internal/appengine/data.go | head -80
(no output)

$ cd /root/astrate-mule && rg -n 'ListInterfacesInGroup|ListInterfacesByAlias|func \(s \*Service\) ListInterfaces' internal/appengine/*.go | rg -v _test
internal/appengine/service.go:793:func (s *Service) ListInterfaces(ctx context.Context, realm, deviceID string) ([]string, error) {
internal/appengine/service.go:809:// ListInterfacesByAlias is ListInterfaces addressed by alias.
internal/appengine/service.go:810:func (s *Service) ListInterfacesByAlias(ctx context.Context, realm, alias string) ([]string, error) {
internal/appengine/service.go:822:// ListInterfacesInGroup is ListInterfaces inside a group (membership gate first).
internal/appengine/service.go:823:func (s *Service) ListInterfacesInGroup(ctx context.Context, realm, groupName, deviceID string) ([]string, error) {
internal/appengine/http.go:179:	names, err := a.svc.ListInterfacesByAlias(r.Context(), r.PathValue("realm"), r.PathValue("alias"))
internal/appengine/http.go:188:	names, err := a.svc.ListInterfacesInGroup(r.Context(), r.PathValue("realm"),

$ cd /root/astrate-mule && rg -n 'WriteNotFound|DetailNotFound' internal/appengine/ ; echo "exit=$?"; echo "--- all 404 writes reachable from appengine writeError ---"; rg -n 'StatusNotFound' internal/appengine/*.go | rg -v _test
exit=1
--- all 404 writes reachable from appengine writeError ---
internal/appengine/http.go:667:		_ = astarteapi.WriteError(w, http.StatusNotFound, "Group not found")
internal/appengine/http.go:669:		_ = astarteapi.WriteError(w, http.StatusNotFound, "Path not found")
internal/appengine/http.go:678:		_ = astarteapi.WriteError(w, http.StatusNotFound, "Interface not found in device introspection")

$ cd /root/astrate-mule && rg -n 'REALM_DELETION_DISABLED|realm_deletion_disabled|default_datastream_maximum_storage_retention|DEFAULT_DATASTREAM' .mule/for-giulio.md | head -20; echo "=== for-giulio tail ==="; tail -25 .mule/for-giulio.md; echo; echo "=== lines ==="; wc -l .mule/for-giulio.md
52:  `default_datastream_maximum_storage_retention` (`*int64`, #73) and `realm_deletion_disabled`
54:  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (config.go:277-292, which
55:  also accepts the bare upstream name `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`
56:  and lets it win when both are set) and `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`
59:  (`realm_deletion_disabled` makes `DELETE /housekeeping/v1/realms/{realm}` answer 405) and the
121:- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`ASTRATE_FLOW_CONFIG`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
191:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
275:  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` and
276:  `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` are consumed in `internal/config/config.go`
=== for-giulio tail ===
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

=== lines ===
472 .mule/for-giulio.md

$ cd /root/astrate-mule && sed -n '1,60p' internal/appengine/writeerror_test.go
package appengine

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astrate-platform/astrate/internal/engine"
)

func TestWriteErrorTaxonomy(t *testing.T) {
	a := &API{}
	cases := []struct {
		name   string
		err    error
		status int
		detail string
	}{
		{"device-owned write", fmt.Errorf("%w: context", engine.ErrNotServerOwned),
			http.StatusMethodNotAllowed, "Cannot write to device owned resource"},
		{"unset on server-owned datastream", fmt.Errorf("%w: context", engine.ErrNotAProperty),
			http.StatusMethodNotAllowed, "Cannot write to read-only resource"},
		{"unknown interface", fmt.Errorf("%w: context", engine.ErrInterfaceNotFound),
			http.StatusNotFound, "Interface not found in device introspection"},
		{"unknown endpoint", fmt.Errorf("%w: context", engine.ErrPathNotFound),
			http.StatusBadRequest, "Endpoint not found"},
		{"read-path path-not-found collision guard", ErrPathNotFound,
			http.StatusNotFound, "Path not found"},
		{"unknown cause stays 500", errors.New("boom"),
			http.StatusInternalServerError, "Internal server error"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			a.writeError(rec, tc.err)
			if rec.Code != tc.status {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, tc.status, rec.Body.String())
			}
			var body struct {
				Errors struct {
					Detail string `json:"detail"`
				} `json:"errors"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("body %q: %v", rec.Body.String(), err)
			}
			if body.Errors.Detail != tc.detail {
				t.Errorf("detail = %q, want %q", body.Errors.Detail, tc.detail)
			}
		})
	}
}

$ cd /root/astrate-mule && rg -n '403|Forbidden' internal/appengine/*_test.go | head -20
internal/appengine/http_test.go:184:		if rec := r.req(t, http.MethodGet, "/devices", "", r.wrong); rec.Code != http.StatusForbidden {
internal/appengine/http_test.go:185:			t.Errorf("wrong claim: got %d, want 403", rec.Code)

$ cd /root/astrate-mule && rg -n 'Internal Server Error|Bad Request' internal/ pkg/ --glob '!*_test.go' | head
(no output)

$ cd /root/astrate-mule && rg -n 'func Test' internal/swagger/swagger_test.go | head -20; echo "=== docs Makefile ==="; cat docs/Makefile
18:func TestMount(t *testing.T) {
77:func TestMountSubPanicsOnBrokenFS(t *testing.T) {
90:func TestSpecs(t *testing.T) {
123:func TestRealmManagement403(t *testing.T) {
179:func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
222:func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
297:func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {
=== docs Makefile ===
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

→ Read internal/appengine/service.go [limit=50, offset=455]

$ cd /root/astrate-mule && rg -n 'func \(s \*Service\) (GetDataInGroup|PutDataInGroup|DeleteDataInGroup|GetDataByAlias|PutDataByAlias|DeleteDataByAlias|ListInterfaces|ListInterfacesByAlias|ListInterfacesInGroup)' -A 18 internal/appengine/data.go internal/appengine/service.go | rg -v _test | head -70
internal/appengine/service.go:726:func (s *Service) GetDataByAlias(ctx context.Context, realm, alias, ifaceName, path string, opts QueryOpts) (any, error) {
internal/appengine/service.go-727-	r, err := s.resolveByAlias(ctx, realm, alias, ifaceName)
internal/appengine/service.go-728-	if err != nil {
internal/appengine/service.go-729-		return nil, err
internal/appengine/service.go-730-	}
internal/appengine/service.go-731-	return s.readResolved(ctx, r, path, opts)
internal/appengine/service.go-732-}
internal/appengine/service.go-733-
internal/appengine/service.go-734-// GetDataInGroup reads an interface endpoint on a group member (upstream
internal/appengine/service.go-735-// GET /groups/{group}/devices/{device}/interfaces/{iface}[/{path}]).
internal/appengine/service.go:736:func (s *Service) GetDataInGroup(ctx context.Context, realm, groupName, deviceID, ifaceName, path string, opts QueryOpts) (any, error) {
internal/appengine/service.go-737-	r, err := s.resolveInGroup(ctx, realm, groupName, deviceID, ifaceName)
internal/appengine/service.go-738-	if err != nil {
internal/appengine/service.go-739-		return nil, err
internal/appengine/service.go-740-	}
internal/appengine/service.go-741-	return s.readResolved(ctx, r, path, opts)
internal/appengine/service.go-742-}
internal/appengine/service.go-743-
internal/appengine/service.go-744-// PublishDataByAlias publishes through an alias address; the resolution gates
internal/appengine/service.go-745-// access and carries the real device ID to the engine write.
internal/appengine/service.go-746-func (s *Service) PublishDataByAlias(ctx context.Context, realm, alias, ifaceName, path string, value json.RawMessage, ts *time.Time) error {
internal/appengine/service.go-747-	r, err := s.resolveByAlias(ctx, realm, alias, ifaceName)
internal/appengine/service.go-748-	if err != nil {
internal/appengine/service.go-749-		return err
internal/appengine/service.go-750-	}
internal/appengine/service.go-751-	return s.PublishData(ctx, realm, r.id.String(), ifaceName, path, value, ts)
internal/appengine/service.go-752-}
internal/appengine/service.go-753-
internal/appengine/service.go-754-// PublishDataInGroup publishes through a group address.
--
internal/appengine/service.go:793:func (s *Service) ListInterfaces(ctx context.Context, realm, deviceID string) ([]string, error) {
internal/appengine/service.go-794-	rid, err := s.realmID(ctx, realm)
internal/appengine/service.go-795-	if err != nil {
internal/appengine/service.go-796-		return nil, err
internal/appengine/service.go-797-	}
internal/appengine/service.go-798-	id, err := deviceid.Parse(deviceID)
internal/appengine/service.go-799-	if err != nil {
internal/appengine/service.go-800-		return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, deviceID)
internal/appengine/service.go-801-	}
internal/appengine/service.go-802-	d, err := s.st.GetDevice(ctx, rid, id)
internal/appengine/service.go-803-	if err != nil {
internal/appengine/service.go-804-		return nil, err
internal/appengine/service.go-805-	}
internal/appengine/service.go-806-	return listInterfacesNames(d), nil
internal/appengine/service.go-807-}
internal/appengine/service.go-808-
internal/appengine/service.go-809-// ListInterfacesByAlias is ListInterfaces addressed by alias.
internal/appengine/service.go:810:func (s *Service) ListInterfacesByAlias(ctx context.Context, realm, alias string) ([]string, error) {
internal/appengine/service.go-811-	rid, err := s.realmID(ctx, realm)
internal/appengine/service.go-812-	if err != nil {
internal/appengine/service.go-813-		return nil, err
internal/appengine/service.go-814-	}
internal/appengine/service.go-815-	d, err := s.st.GetDeviceByAlias(ctx, rid, alias)
internal/appengine/service.go-816-	if err != nil {
internal/appengine/service.go-817-		return nil, err
internal/appengine/service.go-818-	}
internal/appengine/service.go-819-	return listInterfacesNames(d), nil
internal/appengine/service.go-820-}
internal/appengine/service.go-821-
internal/appengine/service.go-822-// ListInterfacesInGroup is ListInterfaces inside a group (membership gate first).
internal/appengine/service.go:823:func (s *Service) ListInterfacesInGroup(ctx context.Context, realm, groupName, deviceID string) ([]string, error) {
internal/appengine/service.go-824-	_, _, d, err := s.groupMember(ctx, realm, groupName, deviceID)
internal/appengine/service.go-825-	if err != nil {
internal/appengine/service.go-826-		return nil, err
internal/appengine/service.go-827-	}
internal/appengine/service.go-828-	return listInterfacesNames(d), nil
internal/appengine/service.go-829-}
internal/appengine/service.go-830-
internal/appengine/service.go-831-// --- helpers ----------------------------------------------------------------
internal/appengine/service.go-832-

$ cd /root/astrate-mule && rg -n 'func \(s \*Service\) (resolveInGroup|resolveByAlias)' -A 30 internal/appengine/service.go | head -70
687:func (s *Service) resolveByAlias(ctx context.Context, realm, alias, ifaceName string) (*resolved, error) {
688-	rid, err := s.realmID(ctx, realm)
689-	if err != nil {
690-		return nil, err
691-	}
692-	d, err := s.st.GetDeviceByAlias(ctx, rid, alias)
693-	if err != nil {
694-		return nil, err
695-	}
696-	return s.resolveOnDevice(ctx, rid, d, ifaceName)
697-}
698-
699-// resolveInGroup maps (realm, group, device id, interface name) to a resolved
700-// interface. Membership gates before the interface lookup: upstream answers
701-// "Device not found" for a non-member even when the interface would also be
702-// missing.
703:func (s *Service) resolveInGroup(ctx context.Context, realm, groupName, deviceID, ifaceName string) (*resolved, error) {
704-	rid, _, d, err := s.groupMember(ctx, realm, groupName, deviceID)
705-	if err != nil {
706-		return nil, err
707-	}
708-	return s.resolveOnDevice(ctx, rid, d, ifaceName)
709-}
710-
711-// readResolved dispatches an already-resolved interface read on its type —
712-// the shared tail of GetData and its mirror variants.
713-func (s *Service) readResolved(ctx context.Context, r *resolved, path string, opts QueryOpts) (any, error) {
714-	if r.iface.Type == interfaceschema.Properties {
715-		// See GetData: reducing a snapshot is refused, not silently ignored.
716-		if opts.DownsamplePoints > 0 {
717-			return nil, fmt.Errorf("%w: downsample_to is not supported on properties interfaces", ErrValidation)
718-		}
719-		return s.propertiesData(ctx, r, path)
720-	}
721-	return s.datastreamData(ctx, r, path, opts)
722-}
723-
724-// GetDataByAlias reads an interface endpoint on an alias-addressed device
725-// (upstream GET /devices-by-alias/{alias}/interfaces/{iface}[/{path}]).
726-func (s *Service) GetDataByAlias(ctx context.Context, realm, alias, ifaceName, path string, opts QueryOpts) (any, error) {
727-	r, err := s.resolveByAlias(ctx, realm, alias, ifaceName)
728-	if err != nil {
729-		return nil, err
730-	}
731-	return s.readResolved(ctx, r, path, opts)
732-}
733-

→ Read internal/swagger/swagger_test.go [limit=100, offset=85]

$ cd /root/astrate-mule && rg -n 'appengine' .mule/todo.md | rg -i 'notfound|not found|group-notfound|403' 
130:- [x] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups/{group}/devices is missing `from_token` and `limit` pagination params that code reads at http.go:439-441 (and returns via links.next), (2) POST /groups is missing 422 validation-error response (FieldErrors for blank group_name/empty devices, http.go:404-416), (3) PATCH /devices/{device} and PATCH /devices-by-alias/{alias} are missing 422 ("Attribute key not found" via ErrAttributeKeyNotFound, service.go:352-357 → http.go:656-657), and the device-scoped PATCH is also missing 409 ("Alias already in use", http.go:651) and its documented 500 covers the merge-patch content-type mismatch which code answers with an unmapped 500 too, (4) GET /devices/{device}/interfaces/{interface}/{path} is missing 422 (invalid since/since_after/to/limit/downsample_to/format params, http.go:554-639). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
132:- [x] docs-sync-appengine-group-patch-status [auto]: PATCH /groups/{group}/devices/{device} (internal/appengine/http.go:518-537, patchGroupDevice) shares applyPatch with the device-scoped PATCH and also answers 422 ("Attribute key not found", ErrAttributeKeyNotFound, service.go:352-357) and 409 ("Alias already in use", http.go:651), but docs/api/astarte_appengine_api.yaml lists only 200/400/401/404/500 for it — add the two missing responses.
136:- [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
198:- [x] channels-rejoin-joinref-tagging-test: pin in internal/appengine/channels/ws_test.go that a second `phx_join` on the same topic updates the session's `joinRef` (ws.go:284-291) and subsequent `new_event` pushes carry the fresh ref (read per event at ws.go:403) — no test drives a rejoin today (all `TestWireSession_*` join once), so the comment's "server pushes must be tagged with the current one" rots unasserted. Wire: join with ref "1", rejoin with "2", push an event, assert element 0 is `"2"`.
202:- [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cd /root/astrate-mule && rg -n 'listDevices|invalid cursor' .mule/todo.md | head
(no output)

$ cd /root/astrate-mule && sed -n '172p;169p;167p;168p' .mule/todo.md
- [x] docs-sync-ae-write-405 [auto]: add the missing `405` response to the nine data-write ops in docs/api/astarte_appengine_api.yaml — PUT/POST/DELETE `/devices/{device}/interfaces/{interface}/{path}`, the three `/devices-by-alias/{alias}` twins, and the three `/groups/{group}/devices/{device}` twins. `writeError` emits `405 Method Not Allowed` with "Cannot write to device owned resource" (`ErrNotServerOwned`) and "Cannot write to read-only resource" (`ErrNotAProperty`) at internal/appengine/http.go:673-676, reachable from putData/putDataByAlias/putDataInGroup (http.go:292-357) and deleteData/deleteDataByAlias/deleteDataInGroup (http.go:311-377); `"405"` currently appears nowhere in the spec (those ops document only 200/204, 400, 401, 404, 500). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-write-value-422 [auto]: add the missing `422` response to the six PUT/POST data-write ops in docs/api/astarte_appengine_api.yaml (device-scoped, by-alias, and in-group at `/interfaces/{interface}/{path}`). A value larger than `maxValueBytes` (64 KiB, internal/appengine/http.go:24) is answered `422` "Value size exceeds size limits" by putData/putDataByAlias/putDataInGroup (http.go:298-301, 327-329, 346-348), but the specs for PUT/POST (yaml:686-695, 726-736, 453-463, 499-509, 1176-1186, 1218-1228) list only 200/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-post-groups-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups` in docs/api/astarte_appengine_api.yaml — `Service.CreateGroup` returns `ErrGroupAlreadyExists` (internal/appengine/service.go:495, mapped at http.go:660-661) which the handler answers as `409` "Group already exists"; the spec (yaml:821-839) documents only 201/400/401/422/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-post-group-devices-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml — re-adding a member the group already holds answers `409` "Device already in group" (`Service.AddGroupDevice` returns `ErrDeviceAlreadyInGroup`, internal/appengine/service.go:636-638, mapped at http.go:662-663), but the spec (yaml:967-977) lists only 201/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

## Objective
- Carry out `.mule/recipes/docs-sync.md` as a proposal-only AppEngine sync: append new task lines to `.mule/todo.md`, without implementing them.
- Do not modify code, documentation, Git state, or existing GitHub state.

## Important Details
- `docs/site/` is never edited; prose/config contradictions go to `.mule/for-giulio.md`.
- `docs/api/*.yaml` drift is proposed as implementation tasks; include `make -C docs build` and Swagger YAML load verification in each API task.
- Recipe: inspect one surface, check roughly 3–4 endpoints, maximum five proposals.
- No evidence file is required by the recipe. It does not call for GitHub issue creation.
- Existing AppEngine tasks already cover by-alias/group routes, group statuses, query parameters, data 405/422 work, and merge-patch media types; do not duplicate them.
- The missing housekeeping config keys are already escalated in `.mule/for-giulio.md`; do not repeat them.
- All output must be English.

## Work State
### Completed
- Selected the `appengine` surface; all 32 registered operations have corresponding YAML operations, so there is no path-level drift.
- Confirmed all 32 AppEngine operations omit `403`; `NewAPI` uses `mw.RequireRealm(auth.ClaimAppEngine)`, all routes use that guard, and authorization failures call `astarteapi.WriteForbidden` in `internal/auth/middleware.go`. The YAML has no `Forbidden` response component.
- Confirmed `patchDeviceByAlias` omits reachable `409 ConflictAliasInUse`: it reaches `applyPatch`, which returns `ErrAliasAlreadyInUse`; the other two PATCH operations already document `409`.
- Confirmed `addGroupDevice` omits reachable `422`: an invalid `device_id` returns `FieldErrors{"device_id": {"is not a valid device id"}}`.
- Confirmed all three data DELETE operations omit reachable `400 "Endpoint not found"`: `UnsetServerProperty` returns `engine.ErrPathNotFound`, mapped by `writeError`; `writeerror_test.go` pins the mapping.
- Confirmed `listDevices` omits reachable `422 "invalid cursor"` for an invalid `from_token`; this is the flat `ErrorDetail` shape, not `FieldErrorsDetail`.
- Additional verified candidates:
  - AppEngine's generic `NotFound` example (`detail: Not Found`) is never emitted by AppEngine handlers; reachable details include `Group not found`, `Device not found`, and `Path not found`.
  - Data GET operations model 404 only as `DeviceNotFound`, although `Path not found` is also reachable.
  - YAML examples `Bad Request` and `Internal Server Error` differ from canonical wire strings `Bad request` and `Internal server error`.
- No files have been modified yet.

### Active
- Finalize at most five atomic task lines. The strongest set is:
  1. `docs-sync-ae-forbidden-403`
  2. `docs-sync-ae-patch-by-alias-409`
  3. `docs-sync-ae-add-group-device-422`
  4. `docs-sync-ae-delete-data-400`
  5. `docs-sync-ae-list-devices-422`
- Append unchecked `[ ] ... [auto]` lines to `.mule/todo.md`; do not implement them.

### Blocked
- No hard blocker.
- Access to `docs/api/astarte_appengine_api.yaml` is intermittently failing with `No such file or directory`/`IO error` from `rg`, `stat`, `cp`, and sometimes `read`; the `grep` tool and Python reads with retries work.

## Next Move
1. Use `grep` or Python-with-retry to confirm current line numbers and exact operation blocks in `docs/api/astarte_appengine_api.yaml`, then re-check `.mule/todo.md` for the five proposed slugs.
2. Append the five detailed proposal lines to `.mule/todo.md`, each naming the handler/source evidence and requiring `make -C docs build` plus Swagger YAML load verification.

## Relevant Files
- `.mule/MULE.md`: project safety and standing rules.
- `.mule/recipes/docs-sync.md`: proposal workflow and evidence requirements.
- `.mule/todo.md`: destination for new unchecked task lines.
- `.mule/for-giulio.md`: already contains the duplicate housekeeping config-key escalation.
- `docs/api/astarte_appengine_api.yaml`: target surface; 32 operations with the verified omissions.
- `internal/appengine/http.go`: route guards and `writeError` status/detail mappings.
- `internal/appengine/service.go`: cursor, patch-conflict, and group-device validation behavior.
- `internal/appengine/data.go`: data path-not-found behavior.
- `internal/auth/middleware.go`: verified-token authorization returns `403`.
- `internal/engine/serverdata.go`: DELETE path mismatch returns `engine.ErrPathNotFound`.
- `internal/appengine/writeerror_test.go`: pins `engine.ErrPathNotFound` → `400 "Endpoint not found"`.
- `pkg/astarteapi/envelope.go`: canonical error-detail strings and envelope writers.
- `docs/site/configuration-reference.md`: never-edit prose; missing keys already escalated.
Continue if you have next steps, or stop and ask for clarification if you are unsure how to proceed.
$ cd /root/astrate-mule && python3 - <<'PY'
import re
p='docs/api/astarte_appengine_api.yaml'
L=open(p).read().split('\n')
# find operation blocks
starts=[]
for i,l in enumerate(L):
    m=re.match(r'^    (get|post|put|delete):$',l)
    if m: starts.append((i,m.group(1)))
    m2=re.match(r'^      operationId: (\w+)$',l)
    if m2: pass
# walk backwards from each status key to nearest operationId
def opid(idx):
    for j in range(idx,-1,-1):
        m=re.match(r'^      operationId: (\w+)$',L[j])
        if m: return m.group(1),j
    return None,None
targets={'patchDeviceByAlias','addGroupDevice','listDevices','deleteDataByAlias','deleteData','deleteDataInGroup','getDataByAlias'}
for key in ['"403"','"409"','"422"','"400"']:
    hits={}
    for i,l in enumerate(L):
        if l.strip().startswith(key+':'):
            o,ol=opid(i)
            hits.setdefault(o,[]).append(i+1)
    print(key, {k:v for k,v in hits.items() if k in targets})
print('--- spec lines of interest ---')
for o in sorted(targets):
    pass
PY
"403" {}
"409" {'addGroupDevice': [976]}
"422" {'patchDeviceByAlias': [287], 'getDataByAlias': [364]}
"400" {'patchDeviceByAlias': [281], 'addGroupDevice': [970]}
--- spec lines of interest ---

$ cd /root/astrate-mule && sed -n '240,295p' docs/api/astarte_appengine_api.yaml
patch:
      operationId: patchDeviceByAlias
      summary: Update device attributes by alias
      description: |
        Partially updates a device's aliases, attributes, or credentials
        inhibited status, looked up by alias. Returns the updated device status.
      tags: [AppEngine - Devices]
      security:
        - a_aea: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - name: alias
          in: path
          required: true
          description: The alias value to look up.
          schema:
            type: string
      requestBody:
        required: true
        content:
          application/merge-patch+json:
            schema:
              type: object
              required: [data]
              properties:
                data:
                  $ref: "#/components/schemas/DevicePatch"
            example:
              data:
                aliases:
                  serial_number: "ABC123"
                attributes:
                  firmware_version: "2.1.0"
                credentials_inhibited: false
      responses:
        "200":
          description: Updated device status.
          content:
            application/json; charset=utf-8:
              schema:
                $ref: "#/components/schemas/DataEnvelopeDeviceStatus"
        "400":
          $ref: "#/components/responses/BadRequest"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "404":
          $ref: "#/components/responses/DeviceNotFound"
        "422":
          $ref: "#/components/responses/UnprocessableEntity"
        "500":
          $ref: "#/components/responses/InternalServerError"

  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces:
    get:
      operationId: listInterfacesByAlias
      summary: List interfaces by alias

$ cd /root/astrate-mule && python3 - <<'PY'
import re
p='docs/api/astarte_appengine_api.yaml'
L=open(p).read().split('\n')
def blocks():
    out=[]
    for i,l in enumerate(L):
        m=re.match(r'^    (get|post|put|delete):$',l)
        if m: out.append((i,m.group(1)))
    comp=None
    for i,l in enumerate(L):
        if l=='components:': comp=i; break
    res=[]
    for n,(i,verb) in enumerate(out):
        end=out[n+1][0] if n+1<len(out) else comp
        o=None
        for j in range(i,end):
            m2=re.match(r'^      operationId: (\w+)$',L[j])
            if m2: o=m2.group(1); break
        codes=set()
        for j in range(i,end):
            m3=re.match(r'^        "(\d{3})":$',L[j])
            if m3: codes.add(m3.group(1))
        res.append((o,verb,i+1,sorted(codes)))
    return res
B=blocks()
print("total ops:",len(B))
for o,v,ln,c in B:
    if o in ('patchDeviceByAlias','addGroupDevice','listDevices','deleteData','deleteDataByAlias','deleteDataInGroup','getGroup','getGroupDevice','getDataByAlias'):
        print(f"  {o:28s} {v:6s} line {ln:5d} {c}")
print("ops with 403:",[o for o,v,l,c in B if '403' in c])
print("any Forbidden component:", any(l.strip()=='Forbidden:' for l in L))
print("any 403 anywhere:", any(l.strip()=='"403":' for l in L))
PY
total ops: 29
  listDevices                  get    line    27 ['200', '401', '500']
  getDataByAlias               get    line   325 ['200', '401', '404', '422', '500']
  deleteDataByAlias            delete line   519 ['204', '401', '404', '405', '500']
  deleteData                   delete line   756 ['204', '401', '404', '405', '500']
  getGroup                     get    line   864 ['200', '401', '404', '500']
  addGroupDevice               post   line   939 ['201', '400', '401', '404', '409', '500']
  getGroupDevice               get    line   982 ['200', '400', '401', '404', '409', '422', '500']
  deleteDataInGroup            delete line  1285 ['204', '401', '404', '405', '500']
ops with 403: []
any Forbidden component: False
any 403 anywhere: False

$ cd /root/astrate-mule && rg -n 'ClaimAppEngine|WriteForbidden|DetailForbidden' internal/appengine/http.go internal/auth/middleware.go pkg/astarteapi/envelope.go internal/appengine/http_test.go | head -20; echo "=== alias conflict ==="; rg -n 'ErrAliasAlreadyInUse' internal/appengine/service.go internal/appengine/http.go; echo "=== fielderrors/cursor/pathnotfound ==="; rg -n 'WriteFieldErrors|invalid cursor|ErrPathNotFound' internal/appengine/http.go internal/appengine/service.go internal/appengine/data.go internal/engine/serverdata.go | head -20
internal/auth/middleware.go:90:			_ = astarteapi.WriteForbidden(w)
internal/auth/middleware.go:115:		_ = astarteapi.WriteForbidden(w)
internal/appengine/http.go:35:	return &API{svc: svc, require: mw.RequireRealm(auth.ClaimAppEngine)}
pkg/astarteapi/envelope.go:40:	// DetailForbidden is the canonical 403 detail.
pkg/astarteapi/envelope.go:41:	DetailForbidden = "Forbidden"
pkg/astarteapi/envelope.go:200:// WriteForbidden writes the canonical 403 envelope.
pkg/astarteapi/envelope.go:201:func WriteForbidden(w http.ResponseWriter) error {
pkg/astarteapi/envelope.go:202:	return WriteError(w, http.StatusForbidden, DetailForbidden)
=== alias conflict ===
internal/appengine/service.go:35:	ErrAliasAlreadyInUse    = errors.New("Alias already in use")    //nolint:staticcheck // ST1005: upstream wire text
internal/appengine/service.go:337:			return nil, ErrAliasAlreadyInUse
internal/appengine/http.go:650:	case errors.Is(err, ErrAliasAlreadyInUse):
=== fielderrors/cursor/pathnotfound ===
internal/appengine/http.go:406:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
internal/appengine/http.go:410:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
internal/appengine/http.go:414:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
internal/appengine/http.go:647:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity, fe)
internal/appengine/http.go:668:	case errors.Is(err, ErrPathNotFound):
internal/appengine/http.go:679:	case errors.Is(err, engine.ErrPathNotFound):
internal/appengine/service.go:153:			return nil, fmt.Errorf("%w: invalid cursor", ErrValidation)
internal/appengine/data.go:29:// ErrPathNotFound marks an object-aggregated series whose concrete path holds
internal/appengine/data.go:32:var ErrPathNotFound = errors.New("appengine: path not found")
internal/appengine/data.go:218:			return nil, fmt.Errorf("%w: %s", ErrPathNotFound, path)
internal/engine/serverdata.go:30:	// ErrPathNotFound: the path resolves no endpoint mapping.
internal/engine/serverdata.go:31:	ErrPathNotFound = errors.New("engine: path matches no endpoint")
internal/engine/serverdata.go:137:			return nil, fmt.Errorf("%w: %q is not an aggregation prefix of %s", ErrPathNotFound, path, ifaceName)
internal/engine/serverdata.go:145:			return nil, fmt.Errorf("%w: %q matches no endpoint of %s", ErrPathNotFound, path, ifaceName)
internal/engine/serverdata.go:243:		return fmt.Errorf("%w: %q matches no endpoint of %s", ErrPathNotFound, path, ifaceName)

$ cd /root/astrate-mule && rg -n 'Forbidden:' -A 8 docs/api/astarte_realm_management_api.yaml | head -20; echo "=== 403 on a rm op ==="; rg -n '"403"' -A 2 docs/api/astarte_realm_management_api.yaml | head -8; echo "=== appengine 403 test ==="; rg -n '403|Forbidden' internal/appengine/http_test.go | head -8
1277:    Forbidden:
1278-      description: Forbidden — valid JWT but insufficient permissions.
1279-      content:
1280-        application/json; charset=utf-8:
1281-          schema:
1282-            $ref: "#/components/schemas/ErrorDetail"
1283-          example:
1284-            errors:
1285-              detail: Forbidden
=== 403 on a rm op ===
90:        "403":
91-          $ref: "#/components/responses/Forbidden"
92-        "500":
--
134:        "403":
135-          $ref: "#/components/responses/Forbidden"
136-        "409":
--
=== appengine 403 test ===
184:		if rec := r.req(t, http.MethodGet, "/devices", "", r.wrong); rec.Code != http.StatusForbidden {
185:			t.Errorf("wrong claim: got %d, want 403", rec.Code)

$ cd /root/astrate-mule && rg -n '^## |^# ' .mule/todo.md | head -20; echo "=== last lines ==="; tail -c 600 .mule/todo.md; echo; echo "=== total lines ==="; wc -l .mule/todo.md
5:# Mule queue
37:## Where tasks come from
70:## Nothing merges on its own
78:## Queue
=== last lines ===
ength`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.

=== total lines ===
235 .mule/todo.md

→ Read .mule/todo.md [offset=228, limit=8]

$ cd /root/astrate-mule && sed -n '27,96p' docs/api/astarte_appengine_api.yaml | rg -n 'name:|description:|maximum|minimum|"20[0-9]"|"\d\d\d"' | head -30; echo "=== listGroupDevices codes ==="; python3 -c "
import re
L=open('docs/api/astarte_appengine_api.yaml').read().split('\n')
for i,l in enumerate(L):
    if 'operationId: listGroupDevices' in l:
        print('line',i+1)
        for j in range(i,i+60):
            if re.match(r'^        \"\d{3}\":',L[j]): print('  ',L[j].strip(),j+1)
            if 'from_token' in L[j] or 'ValidationErrors' in L[j]: print('  *',L[j].strip(),j+1)
"
4:      description: |
12:        - name: details
17:          description: When true, return full device status objects instead of just IDs.
18:        - name: from_token
22:          description: Pagination cursor — start after this device token.
23:        - name: limit
27:          description: Maximum number of devices to return.
29:        "200":
30:          description: Paginated list of devices.
67:        "401":
69:        "500":
=== listGroupDevices codes ===
line 893
  * - name: from_token 910
   "200": 922
   "401": 930
   "404": 932
   "422": 934
  * $ref: "#/components/responses/ValidationErrors" 935
   "500": 936

$ cd /root/astrate-mule && rg -n 'func \(a \*API\) addGroupDevice' -A 20 internal/appengine/http.go; echo "=== ValidationErrors response ==="; rg -n '^    ValidationErrors:' -A 12 docs/api/astarte_appengine_api.yaml
486:func (a *API) addGroupDevice(w http.ResponseWriter, r *http.Request) {
487-	var body groupDeviceBody
488-	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
489-		_ = astarteapi.WriteBadRequest(w)
490-		return
491-	}
492-	if err := a.svc.AddGroupDevice(r.Context(), r.PathValue("realm"), r.PathValue("group"), body.DeviceID); err != nil {
493-		a.writeError(w, err)
494-		return
495-	}
496-	w.WriteHeader(http.StatusCreated)
497-}
498-
499-func (a *API) removeGroupDevice(w http.ResponseWriter, r *http.Request) {
500-	if err := a.svc.RemoveGroupDevice(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device")); err != nil {
501-		a.writeError(w, err)
502-		return
503-	}
504-	w.WriteHeader(http.StatusNoContent)
505-}
506-
=== ValidationErrors response ===
1707:    ValidationErrors:
1708-      description: Validation Error — one or more fields failed validation.
1709-      content:
1710-        application/json; charset=utf-8:
1711-          schema:
1712-            $ref: "#/components/schemas/FieldErrorsDetail"
1713-          example:
1714-            errors:
1715-              group_name:
1716-                - can't be blank
1717-
1718-    ValueTooLarge:
1719-      description: Unprocessable Entity — the value is larger than the per-value size limit (64 KiB).

All five findings verified. Appending the proposal lines to the queue.
$ cd /root/astrate-mule && python3 - <<'PY'
lines = [
"- [ ] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml — the string `\"403\"` appears nowhere in the file and `components` defines no `Forbidden` response, yet every route is wrapped by `a.require = mw.RequireRealm(auth.ClaimAppEngine)` (internal/appengine/http.go:35), which answers `403 Forbidden` for a verified a_aea JWT whose grants do not authorize the method+path (`astarteapi.WriteForbidden`, internal/auth/middleware.go:90; `DetailForbidden`, pkg/astarteapi/envelope.go:41; pinned today by internal/appengine/http_test.go:184-185). Add a `Forbidden` response component copied from the realm-management one (docs/api/astarte_realm_management_api.yaml:1277-1285 — `ErrorDetail` schema, example `errors.detail: Forbidden`) and `$ref` it from all 32 operations, which is exactly the parity docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go:123-170) pins; add the same docs-half test for the AppEngine spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.",
"- [ ] docs-sync-ae-patch-by-alias-409 [auto]: add the missing `409` response to `patchDeviceByAlias` in docs/api/astarte_appengine_api.yaml (yaml:241-291, currently 200/400/401/404/422/500) — the handler resolves the alias then calls the same `applyPatch` as the device- and group-scoped PATCH, and a rename to an alias already in use returns `ErrAliasAlreadyInUse` (internal/appengine/service.go:337) which `writeError` maps to `409 ConflictAliasInUse` (internal/appengine/http.go:650-651); the other two PATCH operations already document that response thanks to docs-sync-appengine-query-params-status and docs-sync-appengine-group-patch-status, so this is the last PATCH of the three with a hole. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.",
"- [ ] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml (yaml:939-980, currently 201/400/401/404/409/500) — a well-formed body whose `device_id` does not parse is rejected as a field error (`FieldErrors{\"device_id\": {\"is not a valid device id\"}}`, internal/appengine/service.go:631-634) and answered 422 with the `FieldErrorsDetail` body (`astarteapi.WriteFieldErrors`, internal/appengine/http.go:647), so `$ref` the existing `ValidationErrors` response (yaml:1707-1716) rather than `BadRequest` — the documented 400 only covers a body `DecodeData` cannot parse (internal/appengine/http.go:488-491). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.",
"- [ ] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data DELETE operations in docs/api/astarte_appengine_api.yaml (`deleteDataByAlias` yaml:519-547, `deleteData` yaml:756-784, `deleteDataInGroup` yaml:1285-1313 — each currently 204/401/404/405/500) — unsetting a path that matches no endpoint mapping returns `engine.ErrPathNotFound` (`UnsetServerProperty`, internal/engine/serverdata.go:241-244) which `writeError` answers as `400 \"Endpoint not found\"` (internal/appengine/http.go:679-680), the same 400 the PUT/POST twins already document via `BadRequest`; internal/appengine/writeerror_test.go pins the mapping. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.",
"- [ ] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 200/401/500) — a `from_token` that does not parse as a device id is refused with `ErrValidation: invalid cursor` (internal/appengine/service.go:150-155) which `writeError` maps to 422 with the flat `{\"errors\":{\"detail\":\"invalid cursor\"}}` body, i.e. the `ErrorDetail` shape and not `FieldErrorsDetail`, so `$ref` the same response the group twin uses for the identical parameter (`ValidationErrors`, yaml:934-935) only after fixing its example, or add a dedicated one. Extend the `from_token` parameter description (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.",
]
with open('.mule/todo.md','a') as f:
    f.write('\n'.join(lines) + '\n')
print('appended', len(lines))
PY
tail -c 300 .mule/todo.md; echo; wc -l .mule/todo.md; git status --short
appended 5
scription (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

240 .mule/todo.md
 M .mule/todo.md

$ cd /root/astrate-mule && rg -n '^- \[ \] docs-sync-ae-' .mule/todo.md | cut -c1-95; echo "=== changed files ==="; git status --short; echo "=== no code/docs touched ==="; git diff --stat
236:- [ ] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the
237:- [ ] docs-sync-ae-patch-by-alias-409 [auto]: add the missing `409` response to `patchDevic
238:- [ ] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /ap
239:- [ ] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data
240:- [ ] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengi
=== changed files ===
 M .mule/todo.md
=== no code/docs touched ===
 .mule/todo.md | 5 +++++
 1 file changed, 5 insertions(+)

Done. Five verified AppEngine doc-drift proposals appended to `.mule/todo.md:236-240`, nothing else touched (no YAML, code, site prose, or git state).

**Surface:** `docs/api/astarte_appengine_api.yaml` — all 32 operations verified present and path-matched against `internal/appengine/http.go`; no path-level drift, so every finding is a status/example gap.

| Task | Gap | Code evidence |
|---|---|---|
| `docs-sync-ae-forbidden-403` | No `"403"` anywhere, no `Forbidden` component, on any of 32 ops | guard `RequireRealm(ClaimAppEngine)` (http.go:35) answers 403 via middleware.go:90; pinned by http_test.go:184-185 |
| `docs-sync-ae-patch-by-alias-409` | `patchDeviceByAlias` documents 200/400/401/404/422/500 | `applyPatch` → `ErrAliasAlreadyInUse` (service.go:337) → 409 (http.go:650-651); the other two PATCHes already have it |
| `docs-sync-ae-add-group-device-422` | Documents 201/400/401/404/409/500 | unparseable `device_id` → `FieldErrors` (service.go:631-634) → 422 (http.go:647); the 400 only covers `DecodeData` failure (http.go:488-491) |
| `docs-sync-ae-delete-data-400` | All three data DELETEs lack 400 | `engine.ErrPathNotFound` (serverdata.go:241-244) → `400 "Endpoint not found"` (http.go:679-680), pinned by writeerror_test.go; PUT/POST already document it |
| `docs-sync-ae-list-devices-422` | Documents 200/401/500 | bad `from_token` → `ErrValidation: invalid cursor` (service.go:150-155) → 422 flat `ErrorDetail` body, not `FieldErrorsDetail`; `limit` deliberately excluded (Atoi error discarded, http.go:87) |

Two notes on the lines: the 403 one names a `TestRealmManagement403`-style docs-half test in `internal/swagger/swagger_test.go` for symmetry, and the `listDevices` one flags that `ValidationErrors`' example is group-specific so reusing it verbatim would be misleading. I did not add a `.mule/for-giulio.md` entry — the one prose finding (missing housekeeping config keys) is already escalated three times there.
