slug: recipe-code-review
verdict: proposed
at:  7568e53
ran: 2026-09-21T11:33:18Z on DietPi in 480s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md

$ ls internal/ pkg/ 2>/dev/null; echo "---reviews---"; ls .mule/reviews/ 2>/dev/null; echo "---log tail---"; tail -30 .mule/log.md 2>/dev/null
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
---reviews---
appengine-2026-09-07.md
astarteapi-2026-09-15.md
auth-2026-09-14.md
broker-2026-09-04.md
channels-2026-09-18.md
config-2026-09-13.md
deviceid-2026-09-18.md
engine-2026-09-06.md
engine-triggers-2026-09-20.md
flow-2026-09-05.md
flowapi-2026-09-10.md
housekeeping-2026-09-15.md
httpx-2026-09-17.md
interfaceschema-2026-09-11.md
observability-2026-09-16.md
pairing-2026-09-09.md
payload-2026-09-12.md
realm-2026-09-09.md
store-2026-09-08.md
---log tail---
| 2026-09-13 | docs-sync-hk-delete-gating-responses | done | 331s | d07b6c7 |
| 2026-09-13 | docs-sync-hk-validation-example | done | 68s | d762f00 |
| 2026-09-14 | auth-iat-not-required-test | done | 257s | 6524c7f |
| 2026-09-14 | auth-empty-claims-403-test | done | 348s | b91247c |
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

$ for p in broker flow engine appengine store pairing realm flowapi interfaceschema payload config auth testutil; do echo "$p: $(rg -c '' --glob "$p/**/*.go" . 2>/dev/null | awk -F: '{s+=$3} END {print s+0}')"; done
broker: 0
flow: 0
engine: 0
appengine: 0
store: 0
pairing: 0
realm: 0
flowapi: 0
interfaceschema: 0
payload: 0
config: 0
auth: 0
testutil: 0

$ for p in broker flow engine appengine store pairing realm flowapi interfaceschema payload config auth testutil; do echo "$p: $(rg -c '' internal/$p --glob '*.go' 2>/dev/null | awk -F: '{s+=$2} END {print s+0}')"; done
broker: 3576
flow: 11157
engine: 15177
appengine: 7435
store: 6783
pairing: 2942
realm: 2709
flowapi: 2335
interfaceschema: 0
payload: 0
config: 1004
auth: 1588
testutil: 697

$ for p in compatible payload interfaceschema; do echo "$p: $(rg -c '' pkg/$p --glob '*.go' 2>/dev/null | awk -F: '{s+=$2} END {print s+0}')"; done; echo "--- largest files in engine ---"; ls -S internal/engine/*.go | head -15; echo "--- largest files in store ---"; ls -S internal/store/*.go | head -15
compatible: 0
payload: 2690
interfaceschema: 3407
--- largest files in engine ---
internal/engine/fakestore_test.go
internal/engine/data.go
internal/engine/cache.go
internal/engine/router.go
internal/engine/data_test.go
internal/engine/batch_test.go
internal/engine/engine_e2e_test.go
internal/engine/engine_integration_test.go
internal/engine/engine.go
internal/engine/serverdata_test.go
internal/engine/batch.go
internal/engine/serverdata.go
internal/engine/control.go
internal/engine/control_test.go
internal/engine/router_test.go
--- largest files in store ---
internal/store/datastreams_test.go
internal/store/devices.go
internal/store/datastreams.go
internal/store/devices_test.go
internal/store/interfaces.go
internal/store/realms_test.go
internal/store/store_test.go
internal/store/realms.go
internal/store/interfaces_test.go
internal/store/groups.go
internal/store/store.go
internal/store/properties.go
internal/store/flows.go
internal/store/pipelines.go
internal/store/flows_test.go

Package `store` was last reviewed 2026-09-08 — good rotation candidate. Let me map its files.
$ for f in internal/store/*.go; do echo "$(wc -l < $f) $f"; done | sort -rn
823 internal/store/datastreams_test.go
485 internal/store/devices.go
434 internal/store/datastreams.go
427 internal/store/devices_test.go
398 internal/store/interfaces.go
355 internal/store/realms_test.go
274 internal/store/store_test.go
265 internal/store/groups.go
250 internal/store/realms.go
249 internal/store/interfaces_test.go
207 internal/store/flows.go
207 internal/store/crypto_test.go
204 internal/store/store.go
193 internal/store/pipelines.go
183 internal/store/properties.go
177 internal/store/flows_test.go
161 internal/store/groups_test.go
148 internal/store/properties_test.go
136 internal/store/crypto.go
134 internal/store/userblocks.go
134 internal/store/compression_test.go
123 internal/store/userblocks_test.go
120 internal/store/retention_test.go
109 internal/store/notify.go
96 internal/store/notify_test.go
88 internal/store/triggers.go
88 internal/store/policies.go
76 internal/store/pipelines_test.go
70 internal/store/policies_test.go
66 internal/store/triggers_test.go
54 internal/store/retention.go
49 internal/store/pipelines_validate_test.go

`store` it is. Reading the largest implementation files first.
→ Read internal/store/devices.go
→ Read internal/store/datastreams.go
→ Read internal/store/interfaces.go
→ Read internal/store/groups.go
→ Read internal/store/store.go
→ Read internal/store/realms.go
→ Read internal/store/properties.go
→ Read internal/store/notify.go
→ Read internal/store/flows.go
→ Read internal/store/pipelines.go
→ Read internal/store/triggers.go
→ Read internal/store/policies.go
→ Read internal/store/retention.go
→ Read internal/store/crypto.go
→ Read internal/store/userblocks.go
Reading the remaining small files and key tests.
✱ Grep "CreateTrigger\(|CreateTriggerPolicy\(|validatePipelineGraph" 30 matches
→ Read internal/store/pipelines_validate_test.go
→ Read internal/realm/service.go [offset=410, limit=130]
✱ Grep "ParseTrigger|ParseSimpleTrigger|action\.|unmarshal.*[Tt]rigger" 14 matches
✱ Grep "drops existing endpoint|Downgrade|drops.*endpoint|inhibit" 2 matches
✱ Grep "func Test|points|bucket|LTTB" 44 matches
✱ Grep "inhibited|inhibit" in docs/DESIGN.md · 46 matches

$ rg -n "4\.4|flow A|register|RegisterDevice" /root/astrate-mule/docs/DESIGN.md | head -30
29:Java/Android, Elixir) pointed at Astrate's Pairing API URL must register, obtain credentials,
83:  (register/CSR :8080) │  internal/appengine  ◀────────────┘          ▼             │
237:    status              text NOT NULL DEFAULT 'registered',  -- registered|confirmed|inhibited
689:### 4.4 Pairing flows (wire-compatible)
702:- Re-registering an existing device that has **not yet requested credentials** rotates the
705:- Unregister: `DELETE /pairing/v1/<realm>/agent/devices/<device_id>` (makes the device
719:`first_credentials_request` and flips status `registered → confirmed`.

→ Read docs/DESIGN.md [offset=689, limit=45]
✱ Grep "func Test|Endpoint|endpoint" 27 matches
→ Read internal/store/notify_test.go

$ rg -n "^func Test" internal/store/devices_test.go internal/store/properties_test.go internal/store/groups_test.go internal/store/realms_test.go internal/store/flows_test.go internal/store/userblocks_test.go internal/store/pipelines_test.go internal/store/retention_test.go internal/store/compression_test.go
internal/store/retention_test.go:17:func TestEnforceRealmRetentionCeilings(t *testing.T) {

→ Read internal/store/store_test.go

$ rg -n "t\.Run\(" internal/store/devices_test.go internal/store/properties_test.go internal/store/groups_test.go internal/store/realms_test.go internal/store/flows_test.go internal/store/retention_test.go
internal/store/devices_test.go:20:	t.Run("Lifecycle", func(t *testing.T) {
internal/store/devices_test.go:150:	t.Run("NotFound", func(t *testing.T) {
internal/store/devices_test.go:167:	t.Run("IntrospectionDiff", func(t *testing.T) {
internal/store/devices_test.go:222:	t.Run("AliasesAndAttributes", func(t *testing.T) {
internal/store/devices_test.go:269:	t.Run("PaginationCursorStability", func(t *testing.T) {
internal/store/devices_test.go:323:	t.Run("StatsAndDelete", func(t *testing.T) {
internal/store/devices_test.go:383:	t.Run("AliasValuesTaken", func(t *testing.T) {
internal/store/properties_test.go:48:	t.Run("UpsertIdempotent", func(t *testing.T) {
internal/store/properties_test.go:72:	t.Run("Unset", func(t *testing.T) {
internal/store/properties_test.go:93:	t.Run("PurgeComplement", func(t *testing.T) {
internal/store/realms_test.go:17:	t.Run("CreateGetRoundTrip", func(t *testing.T) {
internal/store/realms_test.go:99:	t.Run("DuplicateName", func(t *testing.T) {
internal/store/realms_test.go:111:	t.Run("InvalidName", func(t *testing.T) {
internal/store/realms_test.go:124:	t.Run("Updates", func(t *testing.T) {
internal/store/realms_test.go:149:	t.Run("RetentionRoundTrip", func(t *testing.T) {
internal/store/realms_test.go:191:	t.Run("UpdateRealmPatches", func(t *testing.T) {
internal/store/realms_test.go:266:	t.Run("CascadeDelete", func(t *testing.T) {

→ Read internal/store/devices_test.go [limit=170]
→ Read internal/store/devices_test.go [limit=110, offset=222]
→ Read internal/store/devices_test.go [offset=331]
✱ Grep "SinceAfter|Since|To:|SeriesSpan|IndividualSnapshot|LatestIndividual" 27 matches
The package is in good shape — one genuine behaviour bug found. Let me check the current todo file format before appending.
→ Read .mule/todo.md

$ wc -l .mule/todo.md; tail -5 .mule/todo.md
202 .mule/todo.md
- [x] channels-rejoin-joinref-tagging-test: pin in internal/appengine/channels/ws_test.go that a second `phx_join` on the same topic updates the session's `joinRef` (ws.go:284-291) and subsequent `new_event` pushes carry the fresh ref (read per event at ws.go:403) — no test drives a rejoin today (all `TestWireSession_*` join once), so the comment's "server pushes must be tagged with the current one" rots unasserted. Wire: join with ref "1", rejoin with "2", push an event, assert element 0 is `"2"`.
- [!] channels-rejoin-authz-mismatch: `handleJoin`'s comment says upstream authorizes every join including a rejoin (ws.go:276-277), but the rejoin branch (ws.go:284-291) replies OK without calling `s.tok.AuthorizesChannel(auth.VerbJoin, name)` — behaviourally identical today (a token is immutable for the session) but the comment lies about what the code does. Either move the authz check to cover the rejoin branch too, or fix the comment; mechanical, no behaviour change, choose one. — BLOCKED: wrote nothing
- [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a locally-run `elixir -e` with `Base.url_decode64!`) with the device ID `"AAAAAAAAAAAAAAAAAAAAAB"` — nonzero unused bits in the 22nd base64url char — and report whether upstream accepts it: `pkg/deviceid/deviceid.go:34-37` uses `base64.RawURLEncoding.Strict()` and deviceid_test.go:176 asserts rejection while the comment claims parity with `Elixir Base.url_decode64!(padding: false)`, which decodes and discards those bits rather than erroring. Astrate may therefore be stricter than upstream on the same wire form (a device registered with such an id on upstream would be rejected here). If upstream rejects it identically, the strict encoding is verified parity and the task is done; if upstream accepts it, escalate the relax-vs-strict call to `.mule/for-giulio.md` with the measurement. Probe first, no code change either way.
- [x] docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate_devices_total{realm="test"} 42`, but no such series is ever registered — the registry carries `astrate_broker_sessions` (internal/observability/metrics.go:54-62), `astrate_db_pool_{acquired,idle,total,max}_conns` (metrics.go:65-79), plus the engine/flow/trigger families (`astrate_engine_*` engine/router.go:411-444, `astrate_engine_stream_dropped_total` engine/stream/bus.go:119, `astrate_engine_trigger_*` engine/triggers/actions.go:296-300, `astrate_flow_router_*` flow/router.go:260-276) and the standard `go_*`/`process_*` runtime collectors — so the example advertises a scrape result the wire can never return. Swap it for a real series (e.g. `astrate_broker_sessions 0`), keeping the HELP/TYPE comment style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

I have enough. Writing the review file and appending task lines. Let me get the exact tail text first.
→ Read .mule/todo.md [offset=198]
Now writing the review file and appending the task lines.

← Write .mule/reviews/store-2026-09-21.md
Wrote file successfully.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -199,4 +199,7 @@
 - [!] channels-rejoin-authz-mismatch: `handleJoin`'s comment says upstream authorizes every join including a rejoin (ws.go:276-277), but the rejoin branch (ws.go:284-291) replies OK without calling `s.tok.AuthorizesChannel(auth.VerbJoin, name)` — behaviourally identical today (a token is immutable for the session) but the comment lies about what the code does. Either move the authz check to cover the rejoin branch too, or fix the comment; mechanical, no behaviour change, choose one. — BLOCKED: wrote nothing
 - [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a locally-run `elixir -e` with `Base.url_decode64!`) with the device ID `"AAAAAAAAAAAAAAAAAAAAAB"` — nonzero unused bits in the 22nd base64url char — and report whether upstream accepts it: `pkg/deviceid/deviceid.go:34-37` uses `base64.RawURLEncoding.Strict()` and deviceid_test.go:176 asserts rejection while the comment claims parity with `Elixir Base.url_decode64!(padding: false)`, which decodes and discards those bits rather than erroring. Astrate may therefore be stricter than upstream on the same wire form (a device registered with such an id on upstream would be rejected here). If upstream rejects it identically, the strict encoding is verified parity and the task is done; if upstream accepts it, escalate the relax-vs-strict call to `.mule/for-giulio.md` with the measurement. Probe first, no code change either way.
 - [x] docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate_devices_total{realm="test"} 42`, but no such series is ever registered — the registry carries `astrate_broker_sessions` (internal/observability/metrics.go:54-62), `astrate_db_pool_{acquired,idle,total,max}_conns` (metrics.go:65-79), plus the engine/flow/trigger families (`astrate_engine_*` engine/router.go:411-444, `astrate_engine_stream_dropped_total` engine/stream/bus.go:119, `astrate_engine_trigger_*` engine/triggers/actions.go:296-300, `astrate_flow_router_*` flow/router.go:260-276) and the standard `go_*`/`process_*` runtime collectors — so the example advertises a scrape result the wire can never return. Swap it for a real series (e.g. `astrate_broker_sessions 0`), keeping the HELP/TYPE comment style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] store-devices-inhibit-re-register [legion] [auto]: `RegisterDevice` (internal/store/devices.go:75-91) silently clears the inhibit flag — its `ON CONFLICT ... DO UPDATE SET status = 'registered'` fires for any device with `first_credentials_request IS NULL`, including one the admin inhibited via `SetDeviceInhibited` (devices.go:251-268), which sets `status='inhibited'` on unconfirmed devices too; §5.3 says an inhibited device blocks new credentials and connections, so the re-registration re-opens it. Preserve `'inhibited'` in the SET (e.g. `status = CASE WHEN devices.status = 'inhibited' THEN 'inhibited' ELSE 'registered' END`) and add a Lifecycle case in internal/store/devices_test.go: inhibit an unconfirmed device, re-register, assert status stays inhibited and the secret still rotates. Verify the assert against upstream's register-not-touching-inhibit on the Legion while the integration suite runs.
+- [ ] store-pipelines-empty-name-zero-blocks-test [auto]: extend internal/store/pipelines_validate_test.go with the two untested error branches of `validatePipelineGraph` (internal/store/pipelines.go:51-64) — a definition with zero blocks and one whose block has an empty name must both be rejected. Pure helper, no DB; the existing suite covers only acyclic/cyclic/unknown-ref/duplicate-name.
+- [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (internal/store/devices.go:130-131, "if several devices share an alias the lowest device ID wins", `ORDER BY id LIMIT 1`) — no test drives it today; add a case in internal/store/devices_test.go that gives two devices the same alias and asserts the lookup resolves the lower ID. Needs the integration DB.

Done: reviewed `internal/store`, wrote the review file, appended 3 task lines.

Done: code review of internal/store (2026-09-21) — found one inhibit-clearing bug, two test gaps
Files: .mule/reviews/store-2026-09-21.md, .mule/todo.md
Verified: read-only; no code changed, no gate run (proposal job)
Unsure: upstream parity for inhibit-vs-register confirmed at High/medium — noted in the task line to verify on the Legion
Follow-ups: store-devices-inhibit-re-register [legion], store-pipelines-empty-name-zero-blocks-test, store-devices-alias-lowest-id-test [legion] — appended to .mule/todo.md; nothing implemented, no git touched.
