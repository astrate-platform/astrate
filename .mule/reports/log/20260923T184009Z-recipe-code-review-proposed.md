slug: recipe-code-review
verdict: proposed
at:  454929a
ran: 2026-09-23T18:40:09Z on DietPi in 1132s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
Recipe says: pick one package, review deep, write `.mule/reviews/<pkg>-<date>.md`, append ≤5 task lines, no code, no git. Let me look at the map and what's been covered.

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

$ ls -la .mule/ && echo "---reviews---" && ls -la .mule/reviews/ 2>/dev/null && echo "---log---" && ls -la .mule/log.md 2>/dev/null
total 6516
drwxr-xr-x  7 root root    4096 Sep 23 19:07 .
drwxr-xr-x 17 root root    4096 Sep 23 19:04 ..
-rw-r--r--  1 root root      13 Sep 23 20:21 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 23 19:07 .heartbeat
-rw-r--r--  1 root root       0 Sep 23 20:21 .last-output
-rw-r--r--  1 root root       0 Sep 23 20:21 .lock
-rw-r--r--  1 root root      12 Sep 23 20:21 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 6435349 Sep 23 20:21 cron.log
drwxr-xr-x  2 root root    4096 Sep 22 12:58 failed
-rw-r--r--  1 root root   56560 Sep 23 19:47 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    8538 Sep 23 19:07 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 21 21:26 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   75394 Sep 23 19:07 todo.md
---reviews---
total 160
drwxr-xr-x 2 root root 4096 Sep 21 21:26 .
drwxr-xr-x 7 root root 4096 Sep 23 19:07 ..
-rw-r--r-- 1 root root 4600 Sep  7 20:15 appengine-2026-09-07.md
-rw-r--r-- 1 root root 3443 Sep 15 22:09 astarteapi-2026-09-15.md
-rw-r--r-- 1 root root 5192 Sep 14 13:34 auth-2026-09-14.md
-rw-r--r-- 1 root root 8188 Sep  4 20:43 broker-2026-09-04.md
-rw-r--r-- 1 root root 6580 Sep 21 21:26 broker-2026-09-21.md
-rw-r--r-- 1 root root 6967 Sep 18 22:00 channels-2026-09-18.md
-rw-r--r-- 1 root root 4340 Sep 13 19:26 config-2026-09-13.md
-rw-r--r-- 1 root root 5264 Sep 18 19:25 deviceid-2026-09-18.md
-rw-r--r-- 1 root root 3141 Sep  6 20:57 engine-2026-09-06.md
-rw-r--r-- 1 root root 4692 Sep 20 19:30 engine-triggers-2026-09-20.md
-rw-r--r-- 1 root root 5516 Sep  5 21:19 flow-2026-09-05.md
-rw-r--r-- 1 root root 4837 Sep 10 19:00 flowapi-2026-09-10.md
-rw-r--r-- 1 root root 2595 Sep 15 13:29 housekeeping-2026-09-15.md
-rw-r--r-- 1 root root 4914 Sep 17 21:11 httpx-2026-09-17.md
-rw-r--r-- 1 root root 4712 Sep 11 21:40 interfaceschema-2026-09-11.md
-rw-r--r-- 1 root root 4594 Sep 16 20:48 observability-2026-09-16.md
-rw-r--r-- 1 root root 3506 Sep  9 13:07 pairing-2026-09-09.md
-rw-r--r-- 1 root root 6002 Sep 12 20:24 payload-2026-09-12.md
-rw-r--r-- 1 root root 4259 Sep  9 21:13 realm-2026-09-09.md
-rw-r--r-- 1 root root 4605 Sep  8 18:59 store-2026-09-08.md
-rw-r--r-- 1 root root 5504 Sep 21 13:32 store-2026-09-21.md
---log---
-rw-r--r-- 1 root root 8538 Sep 23 19:07 .mule/log.md

→ Read .mule/log.md
→ Read .mule/reviews/store-2026-09-21.md

$ for p in flow engine appengine realm pairing flowapi auth config housekeeping httpx observability store broker swagger testutil; do n=$(rg -c '' --glob "internal/$p/**/*.go" 2>/dev/null | awk -F: '{s+=$2} END {print s+0}'); echo "$p: $n"; done
flow: 11157
engine: 15177
appengine: 7435
realm: 2709
pairing: 2942
flowapi: 2335
auth: 1588
config: 1004
housekeeping: 1372
httpx: 369
observability: 451
store: 6812
broker: 3757
swagger: 171
testutil: 697

$ ls internal/swagger/ internal/testutil/ 2>/dev/null; echo ---; rg -c '' internal/flow/*.go 2>/dev/null | sort -t: -k2 -rn | head -20
internal/swagger/:
swagger.go
swagger_test.go

internal/testutil/:
astartedevice.go
golden.go
golden_test.go
mqttclient.go
pg.go
pg_test.go
---
internal/flow/router_test.go:565
internal/flow/composite_test.go:391
internal/flow/message.go:364
internal/flow/message_test.go:348
internal/flow/flow.go:332
internal/flow/manager_test.go:318
internal/flow/router.go:284
internal/flow/pipeline_test.go:280
internal/flow/factory_test.go:274
internal/flow/factory.go:249
internal/flow/flow_test.go:198
internal/flow/pump_test.go:194
internal/flow/composite.go:156
internal/flow/pipeline.go:155
internal/flow/substitute.go:126
internal/flow/block.go:113
internal/flow/substitute_test.go:85
internal/flow/graph.go:64

Flow was last reviewed 2026-09-05 — the oldest. Let me read its prior review and check what "channels" refers to.
$ head -30 .mule/reviews/channels-2026-09-18.md
# Code review: internal/appengine/channels (+ internal/appengine/stream) — 2026-09-18

Area chosen: the `channels/` subpackage (Phoenix V2 wire protocol at
`/appengine/v1/socket/websocket`) had no review file of its own — the
2026-09-07 appengine review read only `service.go`, `http.go`, `data.go`,
`downsample.go`, so this websocket/room concurrency code was never reviewed.
It is the kind of place (goroutines, mutexes, per-connection loops) a review
earns its keep, and it is container-free, so everything here is T1-checkable.

## Files read

- `internal/appengine/channels/ws.go` (442) — session loop, join/leave/watch/
  unwatch handlers, `watchAuthPath`, `pumpEvents`, `writeFrame`
- `internal/appengine/channels/room.go` (287) — Registry/room lifecycle,
  dispatch fan-out, Member mailbox/drop, retire
- `internal/appengine/channels/protocol.go` (128) — Frame codec, Reply/OK/Err
- `internal/appengine/stream/ws.go` (162) — Astrate-native socket +
  SSE fallback (the documented deviation from the Phoenix protocol)
- `internal/engine/stream/bus.go` (209) — subscription cancel/close semantics
- `internal/engine/triggers/match.go` (330-504) — CompileCondition /
  deviceFilterErrs / group handling
- Tests: `room_test.go` (332), `ws_test.go` (933, all rows), `protocol_test.go`
- Conformance: `test/conformance/upstream/channels.json` (group-watch rows)

## What I found (worth proposing)

### 1. Group watches authorize against the group but deliver realm-wide — the membership is never consulted

**Where:** `internal/appengine/channels/room.go:36-40` (`watchEntry` has no
group field), `room.go:196-200` (the only delivery filter is

→ Read .mule/reviews/flow-2026-09-05.md

$ ls -la internal/flow/
total 176
drwxr-xr-x  3 root root  4096 Sep  5 21:29 .
drwxr-xr-x 17 root root  4096 Sep  4 20:46 ..
-rw-r--r--  1 root root  3664 Sep  4 20:46 block.go
drwxr-xr-x  5 root root  4096 Sep  4 20:46 blocks
-rw-r--r--  1 root root  5892 Sep  5 21:29 composite.go
-rw-r--r--  1 root root 12656 Sep  4 20:46 composite_test.go
-rw-r--r--  1 root root  8323 Sep  4 20:46 factory.go
-rw-r--r--  1 root root  8189 Sep  4 20:46 factory_test.go
-rw-r--r--  1 root root  8821 Sep  4 20:46 flow.go
-rw-r--r--  1 root root  4461 Sep  4 20:46 flow_test.go
-rw-r--r--  1 root root  1945 Sep  4 20:46 graph.go
-rw-r--r--  1 root root  8243 Sep  4 20:46 manager_test.go
-rw-r--r--  1 root root  9445 Sep  4 20:46 message.go
-rw-r--r--  1 root root  8223 Sep  4 20:46 message_test.go
-rw-r--r--  1 root root  4389 Sep  5 22:12 pipeline.go
-rw-r--r--  1 root root  6131 Sep  5 21:29 pipeline_test.go
-rw-r--r--  1 root root  5553 Sep  4 20:46 pump_test.go
-rw-r--r--  1 root root  8033 Sep  4 21:23 router.go
-rw-r--r--  1 root root 15537 Sep  4 21:23 router_test.go
-rw-r--r--  1 root root  3320 Sep  4 20:46 substitute.go
-rw-r--r--  1 root root  2433 Sep  4 20:46 substitute_test.go

Flow was last reviewed 09-05 — the oldest of the main areas. I'll review `internal/flow`. Let me look at the current state and the todo format.
$ rg -n "flow" .mule/todo.md | head -30
85:- [x] flowapi-autorestart-terminal-failure: `restartWithBackoff` (internal/flowapi/service.go:585-621) retries forever on a permanent rebuild failure — a flow whose pipeline was deleted (service.go:360-363 wraps it as ErrValidation) or whose config/definition no longer resolves keeps marking itself failed and rebuilding every ≤30s with no exit short of flow deletion. Classify terminal vs transient (ErrValidation → permanent; store/manager errors → transient) via a pure helper and stop the loop on terminal failures; add a container-free unit test on the helper. [auto]
89:- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.go:560) with no stop signal, so a block-death racing process shutdown (cmd/astrate/main.go:482-488 drains the Manager but never cancels the goroutine) can rebuild the flow via `mgr.StartFlow` on a background pump (internal/flow/flow.go:173-175) after `Manager.Shutdown` and flip the durable status back to running. Thread a stop channel/context through the Service, checked in the loop's sleep, cancelled on shutdown; verification needs a `[legion]` integration or timing-based test. [auto]
90:- [x] flowapi-validationdetail-test: add a container-free unit test for `validationDetail` (internal/flowapi/http.go:297-304) — the 422-detail prefix-strip rule currently has zero coverage. Assert (1) `"flowapi: validation failed: bad block"` strips to `"bad block"`, (2) a non-prefix message passes through unchanged, (3) a message exactly equal to the prefix is NOT stripped (the guard is `len(msg) > len(prefix)`). Pure helper; no store. [auto]
91:- [x] flowapi-autorestart-default-test: pin the wire rule "POST /flows without `auto_restart` defaults to true" (internal/flowapi/http.go:137-144) — currently untested, so a regression to default-false would pass every existing suite. Extract the default into a pure helper (missing → true, present → honored) and assert all three cases container-free. [auto]
103:- [ ] race-check-flow: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/flow/... ./internal/realm/... ./internal/pairing/... ./internal/auth/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
121:- [!] flow-validate-source-sink: `Pipeline.Validate`'s no-source/no-sink error branches (internal/flow/pipeline.go:148-153) are unreachable — the preceding Kahn cycle check (pipeline.go:100-115) guarantees any graph that gets past it is acyclic and therefore always has a source and a sink. The two test cases "no source (all blocks have incoming edges)" and "no sink (all blocks have outgoing edges)" (pipeline_test.go:134-163) both use a 2-cycle, so they pass via the cycle error, never the branch they name. Make the rule real (a cycle-independent formulation, e.g. reject a graph where a node has no path to a sink) or delete the vacuous check, and fix the two test cases to assert the actual error message so a regression to wrong-path coverage fails. — BLOCKED: tests failed: --- FAIL: TestMQTTSink_Retained (0.02s)
122:- [x] flow-validate-dead-source-sink-recompute: in `Pipeline.Validate`, the first `hasSource`/`hasSink` computation (pipeline.go:117-126) is overwritten by the second (pipeline.go:128-147) and its result discarded — dead code (the comment at pipeline.go:127 documents it). Delete the first block; behaviour is unchanged. Tests: existing `TestPipeline_Validate` passes.
170:- [x] payload-longinteger-fraction-quantize [auto]: `jsonInt64` (pkg/payload/json.go:255-271) accepts fractional longinteger tokens that the ParseFloat fallback silently rounds onto an integer inside its own documented-safety window |x| < 2^53 (json.go:250-254 claims values there never round, "silently corrupting the value") — `{"v":4503599627370496.5}` decodes to 4503599627370496, `{"v":6000000000000000.9}` to 6000000000000001, and `{"v":5e-400}` underflows to 0, all without error, while the contract says fraction forms are "accepted only when integral". Verify the token is an exact decimal integer (`math/big.Rat` SetString + IsInt) on the float fallback — keeps `5.0`, `1e3`, `0.5e1` accepted, rejects the rounding cases as ReasonTypeMismatch; plain-integer ParseInt fast path and the > 2^53 window check unchanged (BSON unaffected: longinteger accepts int64/int32 only). Add failing rows to TestJSONProfile.
201:- [x] docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate_devices_total{realm="test"} 42`, but no such series is ever registered — the registry carries `astrate_broker_sessions` (internal/observability/metrics.go:54-62), `astrate_db_pool_{acquired,idle,total,max}_conns` (metrics.go:65-79), plus the engine/flow/trigger families (`astrate_engine_*` engine/router.go:411-444, `astrate_engine_stream_dropped_total` engine/stream/bus.go:119, `astrate_engine_trigger_*` engine/triggers/actions.go:296-300, `astrate_flow_router_*` flow/router.go:260-276) and the standard `go_*`/`process_*` runtime collectors — so the example advertises a scrape result the wire can never return. Swap it for a real series (e.g. `astrate_broker_sessions 0`), keeping the HELP/TYPE comment style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
209:- [x] examples-echo-container-contract-test [auto]: add a container-free httptest suite for examples/flow-container-echo — the only no-test package in the repo that contains logic (`handleMessage`, main.go:49-75; `sanitizeLog`, main.go:40-47) — pinning (1) POST `/v1/message` echoes a valid JSON-object body verbatim with 200, answers 400 on invalid JSON and 204 on an empty body (the 1 MiB read cap), and (2) metadata `echo_drop` == "1" forces the 204 filter/drop path; the example is the canonical reference container contract (README.md, docs/handoff/flow-design-b-container-block-2026-07-29.md), so a regression silently misleads every container author who copies it.

→ Read internal/flow/pipeline.go
Flow (09-05) is the oldest-reviewed main area — I'll review `internal/flow`. Reading files largest first.
→ Read internal/flow/flow.go
The dead recompute block was already deleted (todo line 122 done, matches pipeline.go now). Continuing with the rest of the package.
→ Read internal/flow/router.go
→ Read internal/flow/message.go
→ Read internal/flow/factory.go
→ Read internal/flow/composite.go
→ Read internal/flow/substitute.go
→ Read internal/flow/block.go
→ Read internal/flow/graph.go
→ Read internal/flow/pipeline_test.go
Now the tests.
→ Read internal/flow/router_test.go
Now the `blocks/` subpackage and remaining tests.

$ ls -R internal/flow/blocks/ | head -50; echo ---; for f in internal/flow/blocks/*.go; do echo "$f: $(rg -c '' $f)"; done 2>/dev/null | sort -t: -k2 -rn
internal/flow/blocks/:
astartesource
catalog.go
catalog_test.go
container
httpblocks.go
httpblocks_test.go
info.go
jsonpathmap.go
jsonpathmap_test.go
mqtt.go
mqtt_test.go
randomsource.go
randomsource_test.go
schema.go
schema_test.go
sort.go
sort_test.go
splitmap.go
splitmap_test.go
tojson.go
tojson_test.go
transform.go
transform_test.go
updatemetadata.go
updatemetadata_test.go
virtualdevicepool

internal/flow/blocks/astartesource:
source.go
source_test.go

internal/flow/blocks/container:
block.go
block_test.go
docker.go
httpbridge.go

internal/flow/blocks/virtualdevicepool:
virtualdevicepool.go
virtualdevicepool_test.go
---
internal/flow/blocks/httpblocks_test.go: 577
internal/flow/blocks/mqtt_test.go: 401
internal/flow/blocks/mqtt.go: 388
internal/flow/blocks/httpblocks.go: 344
internal/flow/blocks/jsonpathmap_test.go: 284
internal/flow/blocks/transform.go: 255
internal/flow/blocks/jsonpathmap.go: 238
internal/flow/blocks/randomsource.go: 226
internal/flow/blocks/splitmap_test.go: 192
internal/flow/blocks/transform_test.go: 189
internal/flow/blocks/randomsource_test.go: 173
internal/flow/blocks/catalog.go: 166
internal/flow/blocks/info.go: 162
internal/flow/blocks/sort_test.go: 155
internal/flow/blocks/tojson_test.go: 120
internal/flow/blocks/schema_test.go: 110
internal/flow/blocks/updatemetadata_test.go: 101
internal/flow/blocks/catalog_test.go: 98
internal/flow/blocks/sort.go: 93
internal/flow/blocks/tojson.go: 77
internal/flow/blocks/splitmap.go: 68
internal/flow/blocks/updatemetadata.go: 57
internal/flow/blocks/schema.go: 23

→ Read internal/flow/manager_test.go
Flow-base was reviewed 09-05; the `internal/flow/blocks/*` subpackage was not covered by that review. Let me read the block implementations, largest first.
→ Read internal/flow/blocks/mqtt.go
→ Read internal/flow/blocks/httpblocks.go
The blocks subpackage was not covered in the 09-05 review. Let me read the rest of the block implementations.
→ Read internal/flow/blocks/transform.go
→ Read internal/flow/blocks/jsonpathmap.go
→ Read internal/flow/blocks/randomsource.go
→ Read internal/flow/blocks/catalog.go
→ Read internal/flow/blocks/info.go
→ Read internal/flow/blocks/sort.go
→ Read internal/flow/blocks/tojson.go
→ Read internal/flow/blocks/splitmap.go
→ Read internal/flow/blocks/updatemetadata.go
→ Read internal/flow/blocks/schema.go
→ Read internal/flow/blocks/astartesource/source.go
→ Read internal/flow/blocks/container/block.go
→ Read internal/flow/blocks/virtualdevicepool/virtualdevicepool.go
→ Read internal/flow/blocks/container/httpbridge.go
Now let me sample the block tests to see coverage of the rules I'm finding.
→ Read internal/flow/blocks/httpblocks_test.go
Let me check the remaining test files for the key candidates I've spotted.
→ Read internal/flow/blocks/transform_test.go
→ Read internal/flow/blocks/mqtt_test.go
Let me check sort, jsonpathmap, astartesource, and virtualdevicepool tests.
→ Read internal/flow/blocks/sort_test.go
→ Read internal/flow/blocks/jsonpathmap_test.go
→ Read internal/flow/message_test.go
Let me verify test coverage on the integer-decode path and check the astartesource/virtualdevicepool/container tests.
✱ Grep "int64\(v\)|Int64\(\)|Float64\(\)|UseNumber" in internal/flow · 7 matches
Let me check the astartesource and virtualdevicepool tests for coverage of the rules I'm about to propose.

$ rg -n "func Test|Path|prefix" internal/flow/blocks/virtualdevicepool/virtualdevicepool_test.go | head -40
99:func TestVirtualDevicePool_AcceptsAndAddressesTarget(t *testing.T) {
142:func TestVirtualDevicePool_DeepPath(t *testing.T) {
164:func TestVirtualDevicePool_ZeroTimestampMeansNilTS(t *testing.T) {
185:func TestVirtualDevicePool_DropsBadMessages(t *testing.T) {
216:func TestVirtualDevicePool_NilMessage(t *testing.T) {
233:func TestVirtualDevicePool_IngestErrorIsWrapped(t *testing.T) {
253:		t.Errorf("err = %q, want prefix %q", got, wantPrefix)
257:func TestVirtualDevicePool_ConfigValidation(t *testing.T) {
327:func TestVirtualDevicePool_StaticInvariantNoRegisterWithoutAutoRegister(t *testing.T) {
348:func TestVirtualDevicePool_RegisterDepRequiredWithAutoRegister(t *testing.T) {
359:func TestVirtualDevicePool_FirstSeenRegisters(t *testing.T) {
403:func TestVirtualDevicePool_AlreadyRegisteredDropsWithoutCaching(t *testing.T) {
431:func TestVirtualDevicePool_InfraRegisterErrorSurfaces(t *testing.T) {
453:func TestVirtualDevicePool_ConcurrentFirstSeenSingleRegister(t *testing.T) {

→ Read internal/flow/blocks/astartesource/source_test.go
Let me check the randomsource and catalog tests to nail down coverage gaps.
→ Read internal/flow/blocks/randomsource_test.go
→ Read internal/flow/factory_test.go
Let me verify a few specifics before writing the review: `key_contains` coverage, and confirm the msg integer truncation and `lost` flag claim.

$ rg -n "key_contains" internal/flow/ ; echo "--- lost cleanup ---"; rg -n "lost" internal/flow/blocks/mqtt.go; echo "--- OnConnect ---"; rg -n "OnConnect|ConnectHandler|IsConnected" internal/flow/blocks/mqtt.go
internal/flow/blocks/info.go:44:		Config:  "key_prefix, key_contains, type (integer|real|boolean|datetime|binary|string|map), metadata (object string→string); at least one required",
internal/flow/blocks/schema.go:7:	TypeFilter:            `{"type":"object","properties":{"key_prefix":{"type":"string","description":"Drop unless key starts with this prefix"},"key_contains":{"type":"string","description":"Drop unless key contains this substring"},"type":{"type":"string","description":"Payload type filter; one of integer|real|boolean|datetime|binary|string|map"},"metadata":{"type":"object","description":"Exact-match metadata string→string"}},"required":[]}`,
internal/flow/blocks/transform.go:20://   - key_contains (string): message Key must contain this substring
internal/flow/blocks/transform.go:49:		keyContains: stringConfig(config, "key_contains", ""),
internal/flow/blocks/transform.go:60:		return cfg, fmt.Errorf("at least one of key_prefix, key_contains, type, metadata is required")
--- lost cleanup ---
187:		func(error) { src.lost.Store(true) })
219:	lost     atomic.Bool
234:	if s.lost.Load() {
235:		return nil, fmt.Errorf("mqtt_source: connection lost")
249:// ctx is cancelled. If the connection was lost, it fails instead of waiting
252:	if s.lost.Load() {
253:		return nil, fmt.Errorf("mqtt_source: connection lost")
--- OnConnect ---

$ rg -n "int64\(v\)|float64" internal/flow/message_test.go internal/flow/blocks/ | head
internal/flow/message_test.go:148:				if got.Data.(float64) != tt.msg.Data.(float64) {
internal/flow/blocks/mqtt.go:113:	if f, isFloat := v.(float64); isFloat && f != math.Trunc(f) {
internal/flow/blocks/splitmap.go:20:// (bool→boolean, int64→integer, float64→real, anything else → string via
internal/flow/blocks/splitmap.go:51:				case float64:
internal/flow/blocks/container/block.go:230:	case float64:
internal/flow/blocks/jsonpathmap_test.go:68:	if n := data["n"]; n != float64(3) {
internal/flow/blocks/jsonpathmap_test.go:69:		t.Errorf("n = %v (%T), want float64(3)", n, n)
internal/flow/blocks/jsonpathmap_test.go:109:	if v := data["v"]; v != float64(7) {
internal/flow/blocks/jsonpathmap_test.go:110:		t.Errorf("v = %v (%T), want float64(7)", v, v)
internal/flow/blocks/jsonpathmap_test.go:178:	if v := data["v"]; v != float64(5) {

## Objective
- Carry out the `code-review.md` recipe: review one Go package area of the Astrate codebase, write a review evidence file under `.mule/reviews/`, and append up to 5 proposed task lines to `.mule/todo.md`.
- This is a proposal-only job: write no code, do not start implementing anything, do not touch git (no commit/branch/checkout/stash).

## Important Details
- **Area selected:** `internal/flow`, last reviewed 2026-09-05 (`flow-2026-09-05.md`) — the oldest-reviewed main area. Its `internal/flow/blocks/` subpackage was **never** covered by any prior review, so this review focuses there (mqtt, http, transform, jsonpathmap, randomsource, sort, tojson, splitmap, updatemetadata, catalog, schema, astartesource, container, virtualdevicepool) plus the core package files whose semantics changed since 09-05.
- **Deliverables:** (1) evidence file `.mule/reviews/flow-2026-09-23.md` (today appears to be 2026-09-23 per `.mule` file mtimes); (2) append **at most 5** new task lines to `.mule/todo.md`. Then stop — no implementation, no git.
- **Todo line format** (from existing todo.md): `- [ ] <slug>: <description with file:line refs> [auto]`, with `[legion]` tag where an integration/timing test is needed.
- **Review file format** follows `store-2026-09-21.md` / `flow-2026-09-05.md`: `# ... Review ... <pkg>, <date>` → `## What I read` → findings with file:line.
- **Proposal priority order** (from recipe): (1) wrong/unguarded behavior, (2) missing test on an existing rule (best: machine-checkable), (3) perf only with named workload — bench task before optimisation, (4) mechanical clarity, (5) new features → write one line to `.mule/for-giulio.md`, never to the queue.
- **No `gh issue` involvement:** this recipe's "new features" channel is `.mule/for-giulio.md`, not GitHub issue creation; never run `gh issue edit/close/comment`.
- **MULE.md rules in force:** English only; ~30k-token context budget (use `rg`/narrow reads); split big work into todo.md lines; standing instructions from Giulio outrank recipes (issue #1 never raised/proposed).
- The 09-05 flow review's top finding (unreachable no-source/no-sink Validate branches) is already tracked in todo.md (`flow-validate-source-sink`, line ~121, currently blocked) — do not re-propose it. The companion dead-recompute task (`flow-validate-dead-source-sink-recompute`) was completed; `pipeline.go` no longer has the duplicated block.

## Work State
### Completed
- Read `.mule/MULE.md` (visible top portion incl. 30k-token rule, standing instructions) and `.mule/recipes/code-review.md`.
- Mapped packages (`ls internal/ pkg/`, line counts): flow 11,157; engine 15,177; appengine 7,435; store 6,812; broker 3,757; etc. `swagger` (171) and `testutil` (697) never reviewed but chosen as not worth deep review.
- Checked `.mule/reviews/` inventory and `.mule/log.md`; confirmed rotation: flow last reviewed 09-05 = oldest main area; store/broker reviewed 09-21, channels 09-18, engine-triggers 09-20.
- Read prior reviews `flow-2026-09-05.md`, `store-2026-09-21.md` (format references).
- Read core files: `pipeline.go`, `flow.go`, `router.go`, `message.go`, `factory.go`, `composite.go`, `substitute.go`, `block.go`, `graph.go` (all full, but `flow.go` and `message.go` were read via truncated tool results in places).
- Read all `internal/flow/blocks/*.go` non-test files (mqtt.go, httpblocks.go, transform.go, jsonpathmap.go, randomsource.go, catalog.go, info.go, sort.go, tojson.go, splitmap.go, updatemetadata.go, schema.go) plus `astartesource/source.go`, `container/block.go` (partial), `container/httpbridge.go` (partial), `virtualdevicepool/virtualdevicepool.go`.
- Sampled tests: `pipeline_test.go`, `router_test.go`, `manager_test.go`, `message_test.go` (truncated at end), `httpblocks_test.go`, `transform_test.go`, `mqtt_test.go`, `sort_test.go`, `jsonpathmap_test.go`, `astartesource/source_test.go`, `virtualdevicepool_test.go` (via grep of test names).
- Grepped for number decoding: `message.go:267` has `m.Data = int64(v)` float64 path; `randomsource.go` uses `rand.Int64N` and `Int64/Float64` helpers; container `block.go:233` uses `n.Int64()`.

### Active
- Candidate findings under evaluation (line numbers / severity not yet final):
  1. **`message.go:267` silent integer truncation:** TypeInteger JSON decoded via `float64` does `int64(v)` — `3.7` silently becomes `3`, out-of-int64-range values wrap; the `json.Number` path (`v.Int64()`) would reject, so behavior is inconsistent. Parallels the already-done payload rule (`payload-longinteger-fraction-quantize` in todo.md). High-value "unguarded behavior" candidate.
  2. **`randomsource.go` overflow panic:** `rand.Int64N(s.maxInt-s.minInt+1)` panics when configured min/max span overflows int64 (e.g., min=MinInt64/max=MaxInt64). Need to confirm this reaches a pump goroutine with no recover (flow.go source pump) — process-crash risk from config input.
  3. **`sort.go` unbounded buffer:** messages buffered until a later arrival pushes the head behind `newest-windowUs`; if a stream's timestamps cluster within `window_ms`, `buf` grows without bound (unbounded slice from network input). Also `windowMs * 1000` can overflow for huge config values.
  4. **`mqtt.go` `lost` flag:** `onLost` handler sets `lost`; need to confirm whether paho auto-reconnect ever clears it — if not, `mqtt_source` stays permanently dead after any transient disconnect despite `SetAutoReconnect(true)`.
- Not yet read (may be needed for verification): tail portions of `mqtt.go` (after line ~63), `message.go` `setDataFromWire`/round-trip internals (after line ~60), `httpblocks.go` after line ~57, `randomsource.go` full config-parse + Emit, `sort.go` full flush loop with exact line numbers, `container/block.go` remainder + `docker.go`/`container` tests, `virtualdevicepool_test.go` bodies.

### Blocked
- No hard blockers. Uncertainties: exact line numbers for final findings; whether mqtt `lost` is cleared on reconnect (would kill finding #4); whether the source pump has a recover for the randomsource panic; reading truncated sections may still be needed.
- Do not implement any proposed tasks, and do not run any git command this session.

## Next Move
1. Verify the four candidate findings with narrow `rg`/reads: (a) message.go decode path + round-trip test coverage for fractional/overflowing integers; (b) randomsource min/max overflow reachability and whether flow.go source pump recovers panics; (c) sort.go flush loop exact lines + any buffer cap; (d) mqtt.go `lost` flag set/clear sites.
2. Confirm no proposed finding duplicates existing todo.md lines (grep todo.md for flow/sort/mqtt/message slugs).
3. Write `.mule/reviews/flow-2026-09-23.md` in the store-2026-09-21 format, with `## What I read` (files, largest first, tests) and findings with file:line.
4. Append ≤5 task lines to `.mule/todo.md` (one per confirmed finding, `- [ ] slug: … [auto]`), preferring the container-free/machine-checkable ones.
5. Stop after the todo.md append — do not implement, do not touch git.

## Relevant Files
- `.mule/MULE.md`: safety/context rules and Giulio's standing instructions governing the run.
- `.mule/recipes/code-review.md`: the recipe being executed (proposal-only, ≤5 todo lines, review evidence file).
- `.mule/reviews/flow-2026-09-05.md`: prior flow review (core package; blocks/ not covered).
- `.mule/reviews/store-2026-09-21.md`: format template for the new review file.
- `.mule/todo.md`: target for the ≤5 new task lines; check for duplicate slugs.
- `internal/flow/message.go` (line ~267 `int64(v)` truncation), `internal/flow/blocks/randomsource.go` (Int64N overflow), `internal/flow/blocks/sort.go` (unbounded window buffer), `internal/flow/blocks/mqtt.go` (`lost` flag, auto-reconnect): candidate findings.
- `internal/flow/flow.go`: source-pump lifecycle / panic containment relevant to randomsource finding.
- `internal/flow/blocks/` + `astartesource/`, `container/`, `virtualdevicepool/`: the never-reviewed subpackage that is the review's focus.
- `internal/flow/*_test.go` and `internal/flow/blocks/*_test.go`: coverage evidence for each proposed task.
Continue if you have next steps, or stop and ask for clarification if you are unsure how to proceed.
→ Read internal/flow/blocks/sort.go
→ Read internal/flow/blocks/mqtt.go [offset=175, limit=110]
→ Read internal/flow/blocks/mqtt.go [limit=60, offset=116]
→ Read internal/flow/blocks/mqtt.go [limit=104, offset=285]
✱ Grep "func newMQTTClient" in internal/flow/blocks · 1 match
→ Read internal/flow/blocks/mqtt.go [limit=85, offset=30]
Let me re-verify the remaining candidates precisely before writing the review.
→ Read internal/flow/message.go [limit=75, offset=230]
→ Read internal/flow/blocks/randomsource.go [limit=60, offset=150]
Let me confirm whether `Message.UnmarshalJSON` uses `UseNumber` (important for the float64-truncation finding).
→ Read internal/flow/message.go [limit=85, offset=120]
Let me check the mqtt tests, sort tests, and current todo.md flow entries to confirm no overlap and to characterize coverage.

$ rg -n "func Test" internal/flow/blocks/mqtt_test.go internal/flow/blocks/sort_test.go; echo "=== todo flow/sort/mqtt/msg ==="; rg -n "flow|sort|mqtt|message|msg|random|unbounded|buffer|reconnect|lost" .mule/todo.md
internal/flow/blocks/mqtt_test.go:117:func TestMQTTSource_HappyPath(t *testing.T) {
internal/flow/blocks/mqtt_test.go:154:func TestMQTTSource_MultiTopic(t *testing.T) {
internal/flow/blocks/mqtt_test.go:190:func TestMQTTSource_ConnectionLost(t *testing.T) {
internal/flow/blocks/mqtt_test.go:219:func TestMQTTSink_HappyPath(t *testing.T) {
internal/flow/blocks/mqtt_test.go:267:func TestMQTTSink_Retained(t *testing.T) {
internal/flow/blocks/mqtt_test.go:301:func TestMQTTSource_ConfigValidation(t *testing.T) {
internal/flow/blocks/mqtt_test.go:339:func TestMQTTSink_ConfigValidation(t *testing.T) {
internal/flow/blocks/mqtt_test.go:373:func TestMQTTRegistration_InfoAndSchema(t *testing.T) {
internal/flow/blocks/sort_test.go:23:func TestSort_FlushesAscendingOnArrival(t *testing.T) {
internal/flow/blocks/sort_test.go:66:func TestSort_Dedup(t *testing.T) {
internal/flow/blocks/sort_test.go:102:func TestSort_NilPassthrough(t *testing.T) {
internal/flow/blocks/sort_test.go:110:func TestSort_ConstructRejections(t *testing.T) {
internal/flow/blocks/sort_test.go:141:func TestSort_EqualTimestampsKeepArrivalOrder(t *testing.T) {
=== todo flow/sort/mqtt/msg ===
85:- [x] flowapi-autorestart-terminal-failure: `restartWithBackoff` (internal/flowapi/service.go:585-621) retries forever on a permanent rebuild failure — a flow whose pipeline was deleted (service.go:360-363 wraps it as ErrValidation) or whose config/definition no longer resolves keeps marking itself failed and rebuilding every ≤30s with no exit short of flow deletion. Classify terminal vs transient (ErrValidation → permanent; store/manager errors → transient) via a pure helper and stop the loop on terminal failures; add a container-free unit test on the helper. [auto]
89:- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.go:560) with no stop signal, so a block-death racing process shutdown (cmd/astrate/main.go:482-488 drains the Manager but never cancels the goroutine) can rebuild the flow via `mgr.StartFlow` on a background pump (internal/flow/flow.go:173-175) after `Manager.Shutdown` and flip the durable status back to running. Thread a stop channel/context through the Service, checked in the loop's sleep, cancelled on shutdown; verification needs a `[legion]` integration or timing-based test. [auto]
90:- [x] flowapi-validationdetail-test: add a container-free unit test for `validationDetail` (internal/flowapi/http.go:297-304) — the 422-detail prefix-strip rule currently has zero coverage. Assert (1) `"flowapi: validation failed: bad block"` strips to `"bad block"`, (2) a non-prefix message passes through unchanged, (3) a message exactly equal to the prefix is NOT stripped (the guard is `len(msg) > len(prefix)`). Pure helper; no store. [auto]
91:- [x] flowapi-autorestart-default-test: pin the wire rule "POST /flows without `auto_restart` defaults to true" (internal/flowapi/http.go:137-144) — currently untested, so a regression to default-false would pass every existing suite. Extract the default into a pure helper (missing → true, present → honored) and assert all three cases container-free. [auto]
103:- [ ] race-check-flow: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/flow/... ./internal/realm/... ./internal/pairing/... ./internal/auth/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
108:- [!] broker-disconnect-device-zombie-session: in `internal/broker/broker.go` `DisconnectDevice` (lines 260-266), force-closing the live connection leaves the persisted bbolt session (subscriptions + offline queue, keyed by CN) in place — mochi derives `expire` from the Clean flag, not the reason code (server.go:485-491), so `dropClient` (sessionstore.go:500) is never reached and a later device with the same CN resurrects the deleted device's session. Wire the session store handle into `Broker` and call `dropClient` from `DisconnectDevice`, and add a T1 test in `broker_test.go`: connect clean=false + subscribe, call `DisconnectDevice`, reconnect — assert `session_present` is false and the stale subscription is gone. [approved 2026-09-04] — BLOCKED: wrote nothing
114:- [!] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages" — investigate only: does Astrate have an analogous delay between trigger install and first delivery? Report, do not patch. — BLOCKED: wrote nothing
121:- [!] flow-validate-source-sink: `Pipeline.Validate`'s no-source/no-sink error branches (internal/flow/pipeline.go:148-153) are unreachable — the preceding Kahn cycle check (pipeline.go:100-115) guarantees any graph that gets past it is acyclic and therefore always has a source and a sink. The two test cases "no source (all blocks have incoming edges)" and "no sink (all blocks have outgoing edges)" (pipeline_test.go:134-163) both use a 2-cycle, so they pass via the cycle error, never the branch they name. Make the rule real (a cycle-independent formulation, e.g. reject a graph where a node has no path to a sink) or delete the vacuous check, and fix the two test cases to assert the actual error message so a regression to wrong-path coverage fails. — BLOCKED: tests failed: --- FAIL: TestMQTTSink_Retained (0.02s)
122:- [x] flow-validate-dead-source-sink-recompute: in `Pipeline.Validate`, the first `hasSource`/`hasSink` computation (pipeline.go:117-126) is overwritten by the second (pipeline.go:128-147) and its result discarded — dead code (the comment at pipeline.go:127 documents it). Delete the first block; behaviour is unchanged. Tests: existing `TestPipeline_Validate` passes.
136:- [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
143:- [x] docs-sync-appengine-device-status-schema [auto]: fix the `DeviceStatus` schema in docs/api/astrate_appengine_api.yaml — the wire emits `id` (internal/appengine/service.go:106 `json:"id"`, decoded as `ds.ID` in mirror_test.go:156 and dashboard_compat_test.go:100/149) but the spec names the field `device_id` (yaml:1409) in both the schema and the list/detail examples (yaml:67, 75); `introspection` is documented as `map<string,integer>` (yaml:1438-1441) but each value is a `{major, minor}` object (introspectionEntry, service.go:123-127, introspectionView service.go:833-836); and the schema omits the six fields `groups`, `total_received_msgs`, `total_received_bytes`, `first_credentials_request`, `last_seen_ip` (omitempty), `previous_interfaces` (omitempty) that deviceStatusView emits (service.go:254-279). Rename `device_id`→`id`, correct `introspection` to `map<string,{major:int,minor:int}>`, add the missing fields, and fix the examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
158:- [x] docs-sync-rm-delete-interface-status [auto]: fix the `deleteInterface` response in docs/api/astarte_realm_management_api.yaml — the spec documents `422` with examples "Interface major version is not 0, can't be deleted" / "Cannot delete an interface that is used by a device introspection" (yaml:255-272), but the handler returns `403` with "Interface can't be deleted" / "Interface can't be deleted since it's currently used" (internal/realm/http.go:387-391, `writeInterfaceError`; both are 403 `StatusForbidden`). The spec's detail strings were lifted from the `writeError` branch (http.go:414-419) which DELETE /interfaces does not use. Fix status to `403` and update the two examples to match the actual wire messages. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
164:- [x] interfaceschema-malformed-placeholder-fixtures [auto]: add invalid fixtures + manifest entries for parseSegment's four untested param branches (pkg/interfaceschema/parse.go:509-517) — "/%{", "/%ax", "/%{x", "/%{1}" — each rejected with the malformed-placeholder / invalid-placeholder-name message; the only placeholder-syntax fixture today is the mid-segment 12-endpoint-partial-placeholder.json (ContainsAny branch), so a regression accepting a dangling "/%{" passes every test.
170:- [x] payload-longinteger-fraction-quantize [auto]: `jsonInt64` (pkg/payload/json.go:255-271) accepts fractional longinteger tokens that the ParseFloat fallback silently rounds onto an integer inside its own documented-safety window |x| < 2^53 (json.go:250-254 claims values there never round, "silently corrupting the value") — `{"v":4503599627370496.5}` decodes to 4503599627370496, `{"v":6000000000000000.9}` to 6000000000000001, and `{"v":5e-400}` underflows to 0, all without error, while the contract says fraction forms are "accepted only when integral". Verify the token is an exact decimal integer (`math/big.Rat` SetString + IsInt) on the float fallback — keeps `5.0`, `1e3`, `0.5e1` accepted, rejects the rounding cases as ReasonTypeMismatch; plain-integer ParseInt fast path and the > 2^53 window check unchanged (BSON unaffected: longinteger accepts int64/int32 only). Add failing rows to TestJSONProfile.
176:- [x] config-key-resolvers-test [auto]: add container-free unit tests for `HousekeepingKeys()` and `RealmJWTPublicKey()` (internal/config/config.go:383-411) — both file-reading helpers have zero test coverage. Cover: happy path (inline PEM + file PEM → concatenated/trimmed), file-not-found error (assert error message includes the path), inline-over-file priority (`RealmJWTPublicKey` returns the inline block when both are set), and TrimSpace behaviour (leading/trailing newlines stripped). `t.TempDir()` + `os.WriteFile` for fixtures; no Docker.
179:- [x] docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_housekeeping_api.yaml (yaml:417) — it shows `detail: "housekeeping: validation failed: realm name must not be empty"`, but `validationDetail` strips the `housekeeping: validation failed: ` prefix on the wire (internal/housekeeping/http.go:251-258) and the real blank-name message is `"realm_name can't be blank"` (service.go:138). Change the example to `detail: "realm_name can't be blank"`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
187:- [x] astarteapi-multikey-sorted-golden [auto]: extend TestGoldenEnvelopes (pkg/astarteapi/envelope_test.go:101-121) with multi-key `WriteFieldErrors` and `WriteRawErrors` maps asserting the sorted key order on the wire — the docs claim "keys are emitted in sorted order" (envelope.go:170, 184) but every fixture uses a single key, so a change to the marshalling path that stopped sorting (e.g. switching encoders) passes the suite.
188:- [x] docs-sync-pairing-deviceendpoints-dead-404-403 [auto]: remove the unreachable `404` DeviceNotFound response from the three device-facing ops in docs/api/astarte_pairing_api.yaml — `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials` (yaml:226-227), `GET /pairing/v1/{realm}/devices/{deviceID}` (yaml:268-269), `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify` (yaml:339-340) — and the `403` Forbidden from the info GET (yaml:266-267). `Service.authenticateDevice` folds every failure (malformed/unknown deviceID, unknown realm, unregistered device, wrong secret) into uniform `ErrUnauthorized` (internal/pairing/service.go:380-411), so `store.ErrNotFound` never reaches `writeServiceError` (http.go:326-327) from those handlers; and `Service.Info` never returns `ErrInhibited` (inhibited devices still read info with 200 + status "inhibited", service.go:286-307, comment at 282-285). Reachable on these three are only 200/201/400/401/422/429/500 plus 403 on credentials and verify (both genuinely raise `ErrInhibited`, service.go:252-254 and 329-331). Keep 404 on the agent endpoints (register yaml:149-150, unregister yaml:176-177) where `GetRealmByName`/`deviceid.Parse` really do return `store.ErrNotFound`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
201:- [x] docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate_devices_total{realm="test"} 42`, but no such series is ever registered — the registry carries `astrate_broker_sessions` (internal/observability/metrics.go:54-62), `astrate_db_pool_{acquired,idle,total,max}_conns` (metrics.go:65-79), plus the engine/flow/trigger families (`astrate_engine_*` engine/router.go:411-444, `astrate_engine_stream_dropped_total` engine/stream/bus.go:119, `astrate_engine_trigger_*` engine/triggers/actions.go:296-300, `astrate_flow_router_*` flow/router.go:260-276) and the standard `go_*`/`process_*` runtime collectors — so the example advertises a scrape result the wire can never return. Swap it for a real series (e.g. `astrate_broker_sessions 0`), keeping the HELP/TYPE comment style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
207:- [x] broker-offlineacl-entry-eviction [auto]: `offlineACL.entries` (internal/broker/aclhook.go:100-162) is append-only — `ownershipOf` upserts one entry per CN ever ACL-checked offline and nothing ever deletes, so a long-lived broker grows the map unboundedly across device churn even after the 10s `offlineACLCacheTTL` passes (the entry is kept and re-stamped, not reaped). Evict stale entries lazily on access (drop an entry whose `loadedAt` predates some multiple of the TTL before inserting another, or cap+LRU), add an injectable `now func() time.Time` in the same style as `lifecycleHook.now`, and cover in a container-free unit test: seed N entries, advance the clock, access one, assert the map stays bounded.
208:- [!] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine_api.yaml — `retrieve_metadata` and `downsample_key` are parsed for every data read by `parseQueryOpts` (internal/appengine/http.go:624 `retrieve_metadata`, 633 `downsample_key`) yet appear on none of the read endpoints: device-scoped GET `{interface}` (yaml:583-606) and `{interface}/{path}` (yaml:623-649), by-alias (yaml:324-368, 369-417), and in-group (yaml:1108-1132, 1149-1176) — the parameter block on each lists only since/since_after/to/limit/downsample_to/sort/format/allow_bigintegers/allow_safe_bigintegers (components DataSince..DataAllowSafeBigIntegers, yaml:1360-1439). Add the two `$ref`s (new DataRetrieveMetadata and DataDownsampleKey components, or inline — mirror the existing style) to all six, and update the DataSort-independent `sort` default note only if verified. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
209:- [x] examples-echo-container-contract-test [auto]: add a container-free httptest suite for examples/flow-container-echo — the only no-test package in the repo that contains logic (`handleMessage`, main.go:49-75; `sanitizeLog`, main.go:40-47) — pinning (1) POST `/v1/message` echoes a valid JSON-object body verbatim with 200, answers 400 on invalid JSON and 204 on an empty body (the 1 MiB read cap), and (2) metadata `echo_drop` == "1" forces the 204 filter/drop path; the example is the canonical reference container contract (README.md, docs/handoff/flow-design-b-container-block-2026-07-29.md), so a regression silently misleads every container author who copies it.

→ Read internal/flow/blocks/transform.go [limit=45, offset=40]
Let me verify the pump panic handling, the filter logic, and the mqtt test broker setup.

$ rg -n "recover|defer func" internal/flow/flow.go internal/flow/blocks/ | head; echo "=== mqtt broker type ==="; rg -n "mochi|broker|listener|1883|Listen" internal/flow/blocks/mqtt_test.go | head
internal/flow/blocks/container/httpbridge.go:131:	defer func() { _ = resp.Body.Close() }()
=== mqtt broker type ===
12:	mqtt "github.com/mochi-mqtt/server/v2"
13:	"github.com/mochi-mqtt/server/v2/hooks/auth"
14:	"github.com/mochi-mqtt/server/v2/listeners"
20:// startMQTTBrokerServer embeds a mochi-mqtt broker on a loopback TCP listener
22:// that need to kill the broker early) plus its tcp:// URL.
29:	ln := listeners.NewTCP(listeners.Config{ID: "test-tcp", Address: "127.0.0.1:0"})
30:	if err := srv.AddListener(ln); err != nil {
31:		t.Fatalf("add listener: %v", err)
40:// startMQTTBroker embeds a broker for the whole test and returns its URL.

Let me read the source pump in flow.go and flow's startSourcePump to understand panic/error handling and check the message_test int coverage.
$ rg -n "TypeInteger|integer|Unmarshal|fraction|3\.7|trunc" internal/flow/message_test.go
16:			name: "integer",
19:				Type:      TypeInteger,
113:			if err := json.Unmarshal(b, &got); err != nil {
114:				t.Fatalf("Unmarshal: %v", err)
143:			case TypeInteger:
212:	if err := json.Unmarshal(b, &raw); err != nil {
217:	if err := json.Unmarshal(raw["schema"], &schema); err != nil {
225:func TestFlowMessage_UnmarshalRejectsUnknownSchema(t *testing.T) {
228:	if err := json.Unmarshal([]byte(raw), &msg); err == nil {
248:	if err := json.Unmarshal(b, &raw); err != nil {
268:	if err := json.Unmarshal(b, &raw); err != nil {
290:	if err := json.Unmarshal(b, &raw); err != nil {
322:	if err := json.Unmarshal(b, &raw); err != nil {
336:	types := []DataType{TypeInteger, TypeReal, TypeBoolean, TypeDatetime, TypeBinary, TypeString}

→ Read internal/flow/flow.go [limit=90, offset=150]
Let me look at the sort equal-timestamp test and the message round-trip test details.
→ Read internal/flow/blocks/sort_test.go

$ rg -n "key_contains|KeyContains|Contains" internal/flow/blocks/transform_test.go; echo "=== msg roundtrip integer test bodies ==="; sed -n '10,40p;95,160p' internal/flow/message_test.go
14:	if err == nil || !strings.Contains(err.Error(), "at least one") {
75:	if err == nil || !strings.Contains(err.Error(), "unknown type") {
82:	if err == nil || !strings.Contains(err.Error(), "at least one") {
=== msg roundtrip integer test bodies ===
	tests := []struct {
		name    string
		msg     Message
		wantErr bool
	}{
		{
			name: "integer",
			msg: Message{
				Key:       "sensor/temp",
				Type:      TypeInteger,
				Data:      int64(42),
				Timestamp: 1551884045074181,
				Metadata:  map[string]string{"source": "device-1"},
			},
		},
		{
			name: "real",
			msg: Message{
				Key:       "sensor/pressure",
				Type:      TypeReal,
				Data:      3.14159,
				Timestamp: 1551884045074181,
			},
		},
		{
			name: "boolean",
			msg: Message{
				Key:       "sensor/active",
				Type:      TypeBoolean,
				Data:      true,
				Timestamp: 1551884045074181,
				Type:      TypeString,
				Data:      "test",
				Timestamp: 1000000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(&tt.msg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal: err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			var got Message
			if err := json.Unmarshal(b, &got); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}

			if got.Key != tt.msg.Key {
				t.Errorf("Key = %q, want %q", got.Key, tt.msg.Key)
			}
			if got.Type != tt.msg.Type {
				t.Errorf("Type = %v, want %v", got.Type, tt.msg.Type)
			}
			if got.Timestamp != tt.msg.Timestamp {
				t.Errorf("Timestamp = %d, want %d", got.Timestamp, tt.msg.Timestamp)
			}
			if got.Subtype != tt.msg.Subtype {
				t.Errorf("Subtype = %q, want %q", got.Subtype, tt.msg.Subtype)
			}

			// Compare metadata.
			if len(got.Metadata) != len(tt.msg.Metadata) {
				t.Errorf("Metadata len = %d, want %d", len(got.Metadata), len(tt.msg.Metadata))
			} else {
				for k, v := range tt.msg.Metadata {
					if got.Metadata[k] != v {
						t.Errorf("Metadata[%q] = %q, want %q", k, got.Metadata[k], v)
					}
				}
			}

			// Compare data based on type.
			switch tt.msg.Type {
			case TypeInteger:
				if got.Data.(int64) != tt.msg.Data.(int64) {
					t.Errorf("Data = %v, want %v", got.Data, tt.msg.Data)
				}
			case TypeReal:
				if got.Data.(float64) != tt.msg.Data.(float64) {
					t.Errorf("Data = %v, want %v", got.Data, tt.msg.Data)
				}
			case TypeBoolean:
				if got.Data.(bool) != tt.msg.Data.(bool) {
					t.Errorf("Data = %v, want %v", got.Data, tt.msg.Data)
				}
			case TypeString:
				if got.Data.(string) != tt.msg.Data.(string) {
					t.Errorf("Data = %v, want %v", got.Data, tt.msg.Data)
				}
			case TypeBinary:
				gotBs := got.Data.([]byte)

→ Read .mule/reviews/store-2026-09-21.md [limit=40]
Let me quickly confirm the review file header format and the exact randomsource construct code before writing.
→ Read internal/flow/blocks/randomsource.go [limit=35, offset=80]
→ Read internal/flow/message.go [limit=62, offset=305]

$ rg -n "UseNumber|json\.Unmarshal|Unmarshal\(" internal/flow/ internal/flow/blocks/container/ | rg -i "number|message|raw|body|decoder" | head -20; echo "=== container bridge decode ==="; rg -n "decodeMessage|Unmarshal|json.NewDecoder|Decode\(" internal/flow/blocks/container/*.go
internal/flow/pipeline_test.go:272:	if err := json.Unmarshal(b, &raw); err != nil {
internal/flow/message_test.go:113:			if err := json.Unmarshal(b, &got); err != nil {
internal/flow/message_test.go:212:	if err := json.Unmarshal(b, &raw); err != nil {
internal/flow/message_test.go:217:	if err := json.Unmarshal(raw["schema"], &schema); err != nil {
internal/flow/message_test.go:228:	if err := json.Unmarshal([]byte(raw), &msg); err == nil {
internal/flow/message_test.go:248:	if err := json.Unmarshal(b, &raw); err != nil {
internal/flow/message_test.go:268:	if err := json.Unmarshal(b, &raw); err != nil {
internal/flow/message_test.go:290:	if err := json.Unmarshal(b, &raw); err != nil {
internal/flow/message_test.go:322:	if err := json.Unmarshal(b, &raw); err != nil {
internal/flow/message.go:143:	if err := json.Unmarshal(b, &w); err != nil {
internal/flow/blocks/container/httpbridge.go:161:		if err := json.Unmarshal(respBody, &list); err != nil {
internal/flow/blocks/container/httpbridge.go:174:	if err := json.Unmarshal(respBody, &one); err != nil {
internal/flow/blocks/httpblocks_test.go:393:	if err := json.Unmarshal(got[3].body, &decoded); err != nil {
internal/flow/blocks/container/httpbridge.go:161:		if err := json.Unmarshal(respBody, &list); err != nil {
internal/flow/blocks/container/httpbridge.go:174:	if err := json.Unmarshal(respBody, &one); err != nil {
internal/flow/blocks/container/block_test.go:71:			if err := json.Unmarshal(body, &probe); err != nil {
internal/flow/blocks/container/block_test.go:71:			if err := json.Unmarshal(body, &probe); err != nil {
=== container bridge decode ===
internal/flow/blocks/container/block_test.go:71:			if err := json.Unmarshal(body, &probe); err != nil {
internal/flow/blocks/container/httpbridge.go:161:		if err := json.Unmarshal(respBody, &list); err != nil {
internal/flow/blocks/container/httpbridge.go:174:	if err := json.Unmarshal(respBody, &one); err != nil {

$ rg -n "func TestFilter" internal/flow/blocks/transform_test.go; echo "=== tail todo ==="; tail -5 .mule/todo.md; echo "=== wc ==="; wc -l .mule/todo.md
12:func TestFilter_RequiresCondition(t *testing.T) {
19:func TestFilter_KeyPrefixAndMetadata(t *testing.T) {
58:func TestFilter_Type(t *testing.T) {
73:func TestFilter_UnknownType(t *testing.T) {
154:func TestFilterMap_PipelineInstantiate(t *testing.T) {
=== tail todo ===
- [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (internal/store/devices.go:130-131, "if several devices share an alias the lowest device ID wins", `ORDER BY id LIMIT 1`) — no test drives it today; add a case in internal/store/devices_test.go that gives two devices the same alias and asserts the lookup resolves the lower ID. Needs the integration DB.
- [x] broker-acl-coldstart-fallback-flood [auto]: `syncOwnershipOf` (internal/broker/authhook.go:212-237, reached from the ACL miss path aclhook.go:196-206) does a full synchronous `GetDevice`+`GetInterface` store read for every distinct interface name a device publishes to, with no rate limit — `introspectionReloadDebounce` (authhook.go:42-45, "so an adversarial topic flood cannot hammer the database") throttles only `refreshIfStale`'s full reload, and the fallback's per-name cache write means N distinct bogus names in one window cause N synchronous full-device reads per second, an open anti-flood hole the pre-fix denied-miss path did not have. Gate the fallback to one sync read per session per debounce window (or fold it into a full reload that stamps `lastIntroLoad`), keep `TestBrokerACLColdStartIntrospectionMiss` (one name) green, and add a T1 test in broker_test.go with a GetDevice-counting fake: connect cold, publish to K>=2 distinct unknown-interface names, assert no more than one fallback read was made.
- [x] broker-offlineacl-entry-eviction [auto]: `offlineACL.entries` (internal/broker/aclhook.go:100-162) is append-only — `ownershipOf` upserts one entry per CN ever ACL-checked offline and nothing ever deletes, so a long-lived broker grows the map unboundedly across device churn even after the 10s `offlineACLCacheTTL` passes (the entry is kept and re-stamped, not reaped). Evict stale entries lazily on access (drop an entry whose `loadedAt` predates some multiple of the TTL before inserting another, or cap+LRU), add an injectable `now func() time.Time` in the same style as `lifecycleHook.now`, and cover in a container-free unit test: seed N entries, advance the clock, access one, assert the map stays bounded.
- [!] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine_api.yaml — `retrieve_metadata` and `downsample_key` are parsed for every data read by `parseQueryOpts` (internal/appengine/http.go:624 `retrieve_metadata`, 633 `downsample_key`) yet appear on none of the read endpoints: device-scoped GET `{interface}` (yaml:583-606) and `{interface}/{path}` (yaml:623-649), by-alias (yaml:324-368, 369-417), and in-group (yaml:1108-1132, 1149-1176) — the parameter block on each lists only since/since_after/to/limit/downsample_to/sort/format/allow_bigintegers/allow_safe_bigintegers (components DataSince..DataAllowSafeBigIntegers, yaml:1360-1439). Add the two `$ref`s (new DataRetrieveMetadata and DataDownsampleKey components, or inline — mirror the existing style) to all six, and update the DataSort-independent `sort` default note only if verified. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] examples-echo-container-contract-test [auto]: add a container-free httptest suite for examples/flow-container-echo — the only no-test package in the repo that contains logic (`handleMessage`, main.go:49-75; `sanitizeLog`, main.go:40-47) — pinning (1) POST `/v1/message` echoes a valid JSON-object body verbatim with 200, answers 400 on invalid JSON and 204 on an empty body (the 1 MiB read cap), and (2) metadata `echo_drop` == "1" forces the 204 filter/drop path; the example is the canonical reference container contract (README.md, docs/handoff/flow-design-b-container-block-2026-07-29.md), so a regression silently misleads every container author who copies it.
=== wc ===
209 .mule/todo.md

Let me confirm the transform tests coverage and check the end of todo.md for the append point.
Now I have enough. Writing the review evidence file, then appending the task lines.
← Write .mule/reviews/flow-2026-09-23.md
Wrote file successfully.

Now appending the five task lines to `.mule/todo.md`.
→ Read .mule/todo.md [offset=205, limit=5]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -206,4 +206,9 @@
 - [x] broker-acl-coldstart-fallback-flood [auto]: `syncOwnershipOf` (internal/broker/authhook.go:212-237, reached from the ACL miss path aclhook.go:196-206) does a full synchronous `GetDevice`+`GetInterface` store read for every distinct interface name a device publishes to, with no rate limit — `introspectionReloadDebounce` (authhook.go:42-45, "so an adversarial topic flood cannot hammer the database") throttles only `refreshIfStale`'s full reload, and the fallback's per-name cache write means N distinct bogus names in one window cause N synchronous full-device reads per second, an open anti-flood hole the pre-fix denied-miss path did not have. Gate the fallback to one sync read per session per debounce window (or fold it into a full reload that stamps `lastIntroLoad`), keep `TestBrokerACLColdStartIntrospectionMiss` (one name) green, and add a T1 test in broker_test.go with a GetDevice-counting fake: connect cold, publish to K>=2 distinct unknown-interface names, assert no more than one fallback read was made.
 - [x] broker-offlineacl-entry-eviction [auto]: `offlineACL.entries` (internal/broker/aclhook.go:100-162) is append-only — `ownershipOf` upserts one entry per CN ever ACL-checked offline and nothing ever deletes, so a long-lived broker grows the map unboundedly across device churn even after the 10s `offlineACLCacheTTL` passes (the entry is kept and re-stamped, not reaped). Evict stale entries lazily on access (drop an entry whose `loadedAt` predates some multiple of the TTL before inserting another, or cap+LRU), add an injectable `now func() time.Time` in the same style as `lifecycleHook.now`, and cover in a container-free unit test: seed N entries, advance the clock, access one, assert the map stays bounded.
 - [!] docs-sync-ae-read-query-params: add the two accepted-but-undocumented query parameters to the six interface-data GET ops in docs/api/astarte_appengine_api.yaml — `retrieve_metadata` and `downsample_key` are parsed for every data read by `parseQueryOpts` (internal/appengine/http.go:624 `retrieve_metadata`, 633 `downsample_key`) yet appear on none of the read endpoints: device-scoped GET `{interface}` (yaml:583-606) and `{interface}/{path}` (yaml:623-649), by-alias (yaml:324-368, 369-417), and in-group (yaml:1108-1132, 1149-1176) — the parameter block on each lists only since/since_after/to/limit/downsample_to/sort/format/allow_bigintegers/allow_safe_bigintegers (components DataSince..DataAllowSafeBigIntegers, yaml:1360-1439). Add the two `$ref`s (new DataRetrieveMetadata and DataDownsampleKey components, or inline — mirror the existing style) to all six, and update the DataSort-independent `sort` default note only if verified. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
 - [x] examples-echo-container-contract-test [auto]: add a container-free httptest suite for examples/flow-container-echo — the only no-test package in the repo that contains logic (`handleMessage`, main.go:49-75; `sanitizeLog`, main.go:40-47) — pinning (1) POST `/v1/message` echoes a valid JSON-object body verbatim with 200, answers 400 on invalid JSON and 204 on an empty body (the 1 MiB read cap), and (2) metadata `echo_drop` == "1" forces the 204 filter/drop path; the example is the canonical reference container contract (README.md, docs/handoff/flow-design-b-container-block-2026-07-29.md), so a regression silently misleads every container author who copies it.
+- [ ] flow-mqtt-source-reconnect-recovery [auto]: `mqttSource` latches `lost` on paho's ConnectionLost handler (internal/flow/blocks/mqtt.go:185-187, gates in Emit/Process at 234-236 and 252-254) but nothing ever clears it, even though `newMQTTClient` enables `SetAutoReconnect(true)` (mqtt.go:54) and the docstring promises background auto-reconnect — so any transient broker blip permanently kills the source for the flow's lifetime ("connection lost" forever) while the flow stays `running`; `TestMQTTSource_ConnectionLost` (mqtt_test.go:190) only asserts the error surfaces, never recovery. Clear the latch on reconnect (a paho OnConnect handler, or gate on `client.IsConnected()` instead of a sticky flag), and add a container-free test that stops and restarts the embedded mochi broker and asserts Emit resumes.
+- [ ] flow-msg-json-integer-precision [auto]: `setDataFromWire`'s TypeInteger float64 branch does `int64(v)` on a value that plain `json.Unmarshal` (message.go:141-145) already rounded to float64, so integer wire values > 2^53 silently lose precision and wrap (`float64(MaxInt64)` → 2^63 → MinInt64 on decode), a fractional `3.7` silently truncates to 3, and `1e300` becomes garbage — the `json.Number` branch (message.go:268-273) that would reject these is unreachable in this path. Mirror the accepted `payload-longinteger-fraction-quantize` rule: decode with `json.Decoder`+`UseNumber` (parse via `json.Number.Int64`) or reject non-integral/out-of-int64 float64; add round-trip rows for `math.MaxInt64`, `1e300`, `3.7` to message_test.go. Container-free.
+- [ ] flow-sort-bounded-buffer [auto]: `Sort` (internal/flow/blocks/sort.go:32-65) appends every message to an uncapped `buf` and only flushes while `len(buf) > 1 && buf[0].ts <= newest-windowUs` (sort.go:60) — a stream whose timestamps stay within `window_ms` of the newest (clustered bursts, or messages sharing one timestamp) never satisfies the condition, so the block emits nothing while the slice grows without bound, and the documented "newest is never emitted" mean the final message is held indefinitely; secondarily `window_ms * 1000` (sort.go:31) overflows for large configs and silently inverts the flush edge. Add a configurable max buffered count with a documented overflow policy (drop-oldest / error / force-flush) and a cap test feeding N same-timestamp messages; guard the multiply. Container-free.
+- [ ] flow-randomsource-span-overflow [auto]: `randomSource.next()` computes `rand.Int64N(s.maxInt-s.minInt+1)` (internal/flow/blocks/randomsource.go:171) whose argument overflows int64 to a non-positive value when min/max span more than the int64 range — config like `min:-9e18, max:9e18` passes the single `min <= max` check (randomsource.go:106-108) but makes `rand.Int64N` panic (it panics on n <= 0) inside the source pump goroutine (flow.go:205-239, no recover), crashing the whole process. Reject a span wider than `math.MaxInt64` at construction alongside the existing check, and add a construct/Emit test with the wide span. Container-free.
+- [ ] flow-filter-key-contains-test [auto]: `Filter`'s `key_contains` rule (`strings.Contains`, internal/flow/blocks/transform.go:69) has zero test coverage — `transform_test.go` only exercises `key_prefix`, `type`, unknown-type and the at-least-one construct rule, so a regression silently dropping the Contains check would pass the whole suite. Add substring-match, substring-absent, and prefix+contains combined rows to the filter tests. Container-free.

Done. Review evidence written to `.mule/reviews/flow-2026-09-23.md`, then 5 new task lines appended to `.mule/todo.md`.

Findings (value-ranked):
1. **`mqtt_source` lost-latch never cleared** — sticky `lost` flag + `SetAutoReconnect(true)` means one broker blip kills the source forever (mqtt.go:187, 234, 252).
2. **Integer wire values lose precision** — plain `json.Unmarshal` → float64 → `int64(v)` silently corrupts >2^53 / `3.7` / `1e300`; `json.Number` branch unreachable (message.go:141-145, 267). Same class as the already-fixed `payload-longinteger-fraction-quantize`.
3. **`Sort` unbounded buffer** — clustered/equal timestamps never pass the flush guard (sort.go:60); slice grows without limit.
4. **`random_source` span overflow** — `rand.Int64N(max-min+1)` overflows on wide configs and panics inside the pump goroutine → process crash (randomsource.go:171).
5. **`key_contains` untested** — existing rule with zero coverage (transform.go:69).

No code written, no git touched.
