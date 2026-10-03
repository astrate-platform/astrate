slug: recipe-code-review
verdict: proposed
at:  c7b9a24
ran: 2026-10-03T20:28:17Z on DietPi in 187s


> build · big-pickle

I'll start by reading the required files.
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

$ ls .mule/ && echo "--- reviews ---" && ls .mule/reviews/ 2>/dev/null
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
--- reviews ---
appengine-2026-09-07.md
astarteapi-2026-09-15.md
auth-2026-09-14.md
broker-2026-09-04.md
broker-2026-09-21.md
channels-2026-09-18.md
cmd-astrate-2026-09-25.md
config-2026-09-13.md
container-2026-10-02.md
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

$ ls internal/appengine internal/auth internal/broker internal/config internal/engine internal/flow internal/flowapi internal/housekeeping internal/httpx internal/observability internal/pairing internal/realm internal/store internal/swagger pkg/astarteapi pkg/deviceid pkg/interfaceschema pkg/payload
internal/appengine:
channels
dashboard_compat_test.go
data.go
downsample.go
downsample_test.go
e2e_test.go
formats_parity_test.go
groups_parity_test.go
http.go
http_test.go
links_test.go
mirror_test.go
patch_test.go
query_opts_test.go
service.go
stream
writeerror_test.go

internal/auth:
cache.go
cache_test.go
claims.go
claims_test.go
jwt.go
jwt_test.go
middleware.go
middleware_test.go
testdata

internal/broker:
aclhook.go
aclhook_test.go
authhook.go
broker.go
broker_e2e_test.go
broker_test.go
identity.go
identity_test.go
intake.go
lifecycle.go
publisher.go
sessionstore.go
sessionstore_test.go

internal/config:
config.example.toml
config.go
config_test.go

internal/engine:
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

internal/flow:
block.go
blocks
composite.go
composite_test.go
factory.go
factory_test.go
flow.go
flow_test.go
graph.go
manager_test.go
message.go
message_test.go
pipeline.go
pipeline_test.go
pump_test.go
router.go
router_test.go
substitute.go
substitute_test.go

internal/flowapi:
composite_wire_test.go
http.go
http_test.go
service.go
service_integration_test.go
service_test.go
userblocks.go
userblocks_test.go

internal/housekeeping:
async_operation_test.go
http.go
http_test.go
service.go
service_test.go

internal/httpx:
cors.go
cors_test.go
notfound.go
notfound_test.go

internal/observability:
compat.go
compat_version_test.go
health.go
metrics.go
observability_test.go

internal/pairing:
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

internal/realm:
async_operation_test.go
dashboard_compat_test.go
detailed.go
detailed_required_encrypted_test.go
detailed_test.go
helpers_test.go
http.go
http_test.go
policy.go
policy_test.go
service.go

internal/store:
compression_test.go
crypto.go
crypto_test.go
datastreams.go
datastreams_test.go
devices.go
devices_test.go
flows.go
flows_test.go
groups.go
groups_test.go
interfaces.go
interfaces_test.go
notify.go
notify_test.go
pipelines.go
pipelines_test.go
pipelines_validate_test.go
policies.go
policies_test.go
properties.go
properties_test.go
realms.go
realms_test.go
retention.go
retention_test.go
store.go
store_test.go
testdata
triggers.go
triggers_test.go
userblocks.go
userblocks_test.go

internal/swagger:
swagger.go
swagger_test.go

pkg/astarteapi:
envelope.go
envelope_test.go
testdata

pkg/deviceid:
deviceid.go
deviceid_test.go

pkg/interfaceschema:
compat.go
compat_test.go
compile.go
compile_test.go
parse.go
parse_test.go
required_encrypted_test.go
testdata
trie.go
trie_test.go
types.go
violations.go
violations_test.go

pkg/payload:
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

$ tail -60 .mule/log.md
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
| 2026-09-25 | drain-per-stage-budget | done | 265s | 84a050e |
| 2026-09-25 | cmd-devcert-fields-test | done | 357s | f1434eb |
| 2026-09-25 | cmd-healthcheck-contract-test | blocked | 199s | tests failed: --- FAIL: TestRunHealthcheckProbesReadiness (0.01s) |
| 2026-09-25 | cmd-loadsealer-masterkeyfile-test | done | 221s | a6711f9 |
| 2026-09-25 | docs-sync-hk-async-operation-param | done | 248s | ed1cd72 |
| 2026-09-25 | docs-sync-hk-retention-zero-is-unset | done | 325s | c2ad442 |
| 2026-09-25 | docs-sync-rm-async-operation-param | done | 425s | 3f0eca8 |
| 2026-09-25 | docs-sync-rm-delete-device-async-operation-param | done | 903s | f965f39 |
| 2026-09-25 | docs-sync-rm-delete-device-async-operation-param | blocked | 177s | wrote nothing |
| 2026-09-26 | forward-static-header-validation | done | 874s | 1cbebd5 |
| 2026-09-26 | forward-static-headers-override-test | done | 263s | 3a0f153 |
| 2026-09-26 | forward-status-error-body | done | 209s | 915c260 |
| 2026-09-26 | forward-envelope-bytes-test | done | 213s | ed869b6 |
| 2026-09-26 | webhook-static-headers-override-test | done | 178s | f9db0a6 |
| 2026-09-26 | docs-sync-ae-forbidden-403 | done | 1037s | 6a29ac7 |
| 2026-09-27 | docs-sync-ae-patch-by-alias-409 | done | 119s | d1af059 |
| 2026-09-27 | docs-sync-ae-add-group-device-422 | done | 274s | dff3be8 |
| 2026-09-27 | docs-sync-ae-delete-data-400 | done | 330s | 66bc366 |
| 2026-09-27 | docs-sync-ae-list-devices-422 | transient | 344s |  > build · big-pickle  Error: Error from provider (Console): Rate limit exceeded. Please try again later.  |
| 2026-09-27 | docs-sync-ae-list-devices-422 | transient | 248s |  > build · big-pickle  Error: Error from provider (Console): Rate limit exceeded. Please try again later.  |
| 2026-09-27 | docs-sync-ae-list-devices-422 | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-09-28 | docs-sync-pairing-error-example-capitalisation | done | 166s | 7ed78cb |
| 2026-09-28 | docs-sync-pairing-info-version-example | done | 333s | d9ec98f |
| 2026-09-28 | docs-sync-pairing-deviceid-base64url | done | 480s | 32693cf |
| 2026-09-28 | docs-sync-pairing-initial-payload-format-enum | done | 95s | 1f6a6c2 |
| 2026-09-28 | docs-sync-pairing-unregister-description | done | 261s | 8ffb39e |
| 2026-10-01 | fdo-rc6-scope-delta-for-giulio | done | 479s | 41af3c4 |
| 2026-10-02 | compat-note-v14-rc6 | done | 398s | 0215f97 |
| 2026-10-02 | container-stop-deadline | done | 410s | f9af73b |
| 2026-10-02 | container-timeout-bounds | blocked | 1201s | TIMEOUT after 1201s — task too big, split it |
| 2026-10-02 | container-parseconfig-rules-test | done | 274s | 32c85af |
| 2026-10-02 | container-response-cap-test | blocked | 271s | lint failed: internal/flow/blocks/container/block_test.go:500:2: redefines-builtin-id: redefinition of the built-in function cap (rev |
| 2026-10-02 | docs-sync-rm-error-example-capitalisation | done | 751s | 1d1e6f4 |
| 2026-10-02 | docs-sync-rm-legacy-alias-fields | blocked | 329s | wrote nothing |
| 2026-10-02 | docs-sync-rm-validationerror-example | done | 732s | 8a09a1a |
| 2026-10-02 | docs-sync-rm-deviceid-param | done | 303s | 2d1a1c1 |
| 2026-10-02 | docs-sync-rm-update-interface-body | blocked | 1201s | TIMEOUT after 1201s — task too big, split it |
| 2026-10-03 | docs-sync-ae-error-example-capitalisation | done | 298s | f3bcc6d |
| 2026-10-03 | docs-sync-hk-error-example-capitalisation | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-10-03 | docs-sync-native-error-example-capitalisation | done | 284s | 7cc237f |
| 2026-10-03 | docs-sync-hk-error-detail-examples-split | done | 114s | 23006d4 |
| 2026-10-03 | device-empty-cache-received-trigger | blocked | 202s | wrote nothing |
| 2026-10-03 | appengine-unexpected-object-key | done | 956s | f7d60af |
| 2026-10-03 | appengine-payload-reason-status-map | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-10-03 | appengine-missing-required-422 | done | 584s | f84bbb5 |

$ head -12 .mule/reviews/channels-2026-09-18.md .mule/reviews/container-2026-10-02.md .mule/reviews/forward-2026-09-26.md .mule/reviews/engine-triggers-2026-09-20.md
==> .mule/reviews/channels-2026-09-18.md <==
# Code review: internal/appengine/channels (+ internal/appengine/stream) — 2026-09-18

Area chosen: the `channels/` subpackage (Phoenix V2 wire protocol at
`/appengine/v1/socket/websocket`) had no review file of its own — the
2026-09-07 appengine review read only `service.go`, `http.go`, `data.go`,
`downsample.go`, so this websocket/room concurrency code was never reviewed.
It is the kind of place (goroutines, mutexes, per-connection loops) a review
earns its keep, and it is container-free, so everything here is T1-checkable.

## Files read

- `internal/appengine/channels/ws.go` (442) — session loop, join/leave/watch/

==> .mule/reviews/container-2026-10-02.md <==
# Review — internal/flow/blocks/container, 2026-10-02

Rotation check first: every `internal/` and `pkg/` package has a review file except
`internal/testutil` (test helpers, no logic worth proposing on). The 2026-09-23 flow
pass read `container/block.go` and `container/httpbridge.go` but **not `container/docker.go`**
(211 lines, the whole Docker-CLI surface: `CLIRunner`, `cliInstance`, `Spec`, `Waiter`,
`encodeFlowConfigJSON`) and drew no findings from this package at all. This pass is that
missing half.

## What I read

- `internal/flow/blocks/container/block.go` (all 320) — `Config`, `New` (Start, bridge

==> .mule/reviews/forward-2026-09-26.md <==
# Review — internal/engine/forward (2026-09-26)

First review of this area, so no prior file to extend. Chosen because it is the only
`internal/engine` sub-package with no review behind it, and because it is the seam where
a tenant's event payload leaves the process for something the mule does not control.

## What I read

- `internal/engine/forward/http.go` (all 125 lines) — `Config`, `New`, `validMethod`, `HTTP.Forward`
- `internal/engine/forward/envelope.go` (all 29 lines) — the `envelope` shape and `marshalEnvelope`
- `internal/engine/forward/nats.go` (all 53 lines, `//go:build nats`)
- `internal/engine/forward/http_test.go` (all 310 lines) and `nats_test.go` (all 196 lines)

==> .mule/reviews/engine-triggers-2026-09-20.md <==
# Code review — internal/engine (triggers + forward) — 2026-09-20

Second pass over `internal/engine`. The 09-06 review covered the core data
paths; this run targets the files the rotation had not read since: the
`triggers` subpackage's validation/transient/error-name machinery and the
`forward/` extension package, plus the delivery path in `actions.go`.

## What I read

- `internal/engine/triggers/actions.go` (610 lines) — parseAction, Executor,
  Enqueue/deliver/forward/webhook, policy retry loop
- `internal/engine/triggers/validation.go` (163 lines) — URL/header/template

$ ls internal/testutil internal/flow/blocks internal/appengine/stream internal/engine/stream internal/pairing/ca tools/ 2>/dev/null; echo "--- root ---"; ls
internal/appengine/stream:
ws.go
ws_test.go

internal/engine/stream:
bus.go
bus_test.go

internal/flow/blocks:
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

internal/pairing/ca:
ca.go
ca_test.go

internal/testutil:
astartedevice.go
golden.go
golden_test.go
mqttclient.go
pg.go
pg_test.go

tools/:
bsoncapture
mule-plan-day.sh
mule-setup-pi.sh
mule-survey.sh
mule.sh
reconcile.sh
--- root ---
AGENTS.md
Dockerfile
LICENSE
Makefile
README.md
astrate
bench
cmd
deploy
docker-compose.yml
docs
examples
flow-container-echo
go.mod
go.sum
internal
migrations
opencode.json
pkg
site-dist
test
tools

$ rg -c '' --glob 'internal/flow/blocks/**/*.go' | sort -t: -k2 -rn
internal/flow/blocks/container/block_test.go:766
internal/flow/blocks/httpblocks_test.go:577
internal/flow/blocks/virtualdevicepool/virtualdevicepool_test.go:493
internal/flow/blocks/mqtt_test.go:401
internal/flow/blocks/mqtt.go:388
internal/flow/blocks/httpblocks.go:344
internal/flow/blocks/container/block.go:320
internal/flow/blocks/jsonpathmap_test.go:284
internal/flow/blocks/transform.go:255
internal/flow/blocks/jsonpathmap.go:238
internal/flow/blocks/transform_test.go:236
internal/flow/blocks/randomsource.go:226
internal/flow/blocks/container/docker.go:223
internal/flow/blocks/virtualdevicepool/virtualdevicepool.go:214
internal/flow/blocks/splitmap_test.go:192
internal/flow/blocks/sort_test.go:190
internal/flow/blocks/container/httpbridge.go:178
internal/flow/blocks/randomsource_test.go:173
internal/flow/blocks/sort.go:171
internal/flow/blocks/astartesource/source.go:171
internal/flow/blocks/catalog.go:166
internal/flow/blocks/info.go:162
internal/flow/blocks/astartesource/source_test.go:153
internal/flow/blocks/tojson_test.go:120
internal/flow/blocks/schema_test.go:110
internal/flow/blocks/updatemetadata_test.go:101
internal/flow/blocks/catalog_test.go:98
internal/flow/blocks/tojson.go:77
internal/flow/blocks/splitmap.go:68
internal/flow/blocks/updatemetadata.go:57
internal/flow/blocks/schema.go:23

$ head -30 .mule/reviews/flow-2026-09-23.md; echo "=== engine 0906 ==="; head -30 .mule/reviews/engine-2026-09-06.md
# Review — internal/flow, 2026-09-23

Last reviewed: 2026-09-05 (`flow`). That pass covered the core package;
`internal/flow/blocks/` (and its `astartesource/`, `container/`,
`virtualdevicepool/` subpackages) has never been reviewed, so this pass
focuses there plus the message wire-format decode path that feeds them.

## What I read

Core (`internal/flow/`):
- `pipeline.go` — Parse/Validate; the 09-05 dead source/sink branch was
  confirmed fixed (the duplicate hasSource/hasSink recompute is now gone,
  matching todo line 122).
- `flow.go`, `router.go`, `message.go`, `factory.go`, `composite.go`,
  `substitute.go`, `block.go`, `graph.go` — lifecycle, lanes, message
  wire format (Marshal/UnmarshalJSON + setDataFromWire), registry,
  composite expansion, block graph. Tests read: `message_test.go`
  (roundtrip + rejection), `factory_test.go` (instantiate linear /
  unknown type / construct-error / stop-on-failure / parse def /
  astarte source → null sink lifecycle).

Blocks (`internal/flow/blocks/`):
- `mqtt.go` — MQTT source/sink, shared `newMQTTClient` (auto-reconnect),
  connection-lost latch.
- `httpblocks.go` — HTTP source/sink (opt, deadline derivation, JSON body
  rules).
- `transform.go` — Filter (key_prefix/key_contains/type/metadata) + Map.
- `jsonpathmap.go` — JSONPath w/ value extraction + arrays.
- `randomsource.go` — random integer/real/boolean source (min/max/interval).
- `sort.go` — windowed ascending-timestamp sort + dedup.
=== engine 0906 ===
# Code review — internal/engine — 2026-09-06

Area chosen because `broker` (2026-09-04) and `flow` (2026-09-05) were recent reviews and `engine` is the largest un-reviewed package (48 files).

## What I read

- `internal/engine/router.go` (engine struct, sharding, Submit, Drain, metrics)
- `internal/engine/engine.go` (fireCommitted / fireData / fireDataChanges / fireDevice / handleLifecycle)
- `internal/engine/data.go` (validation pipeline, PersistOp, trackPrevious)
- `internal/engine/batch.go` (micro-batcher, row conversion, encodeValueJSON/jsons constraints)
- `internal/engine/cache.go` (schema + device caches, buildRealm trigger/policy load)
- `internal/engine/control.go` (emptyCache / producer+consumer properties, zip-bomb guard)
- `internal/engine/introspection.go`, `topics.go`, `capabilities.go`, `serverdata.go`
- `internal/engine/triggers/{match,actions,events,policy}.go`
- `internal/engine/stream/bus.go`
- tests: `serverdata_test.go`, `devicedata_test.go`, `triggers_change_test.go`

## What I found (worth proposing)

1. **Server- and device-owned publish paths never fire data triggers nor fan out to the live bus.** `publishAsOwner` (internal/engine/serverdata.go:94) persists the op (`UpsertProperty` / `AppendDatastreams`) and (for `PublishServerValue`) publishes to the broker, but never calls `fireData`/`fireCommitted` and never does `e.bus.Publish`. The broker-driven ingress path gets both through `finalize` → `afterCommit` (batch.go:186) → `fireCommitted` (engine.go:140). So:
   - a server-owned AppEngine PUT/POST or `UnsetServerProperty` produces no `incoming_data`/`value_stored`/`value_change` trigger event and no live-stream card;
   - a virtual device write via `PublishDeviceValue` (issue #84) — whose stated contract is "lands storage rows exactly like a real device's data would" (devicedata_test.go:3-4) — likewise fires no triggers.
   None of the existing tests (serverdata_test.go, devicedata_test.go) assert trigger or bus output on these paths, so the absence is unguarded. This is category-1 (a behaviour that is wrong/unbalanced with a sibling path). Proposed as a task with a failing test. Intent should be confirmed: it changes observable event behaviour, so the reviewer should decide it is not a deliberate "server-originated writes are silent" choice.

## What I decided NOT to propose, and why

- **`deviceIDStringLen = 22` duplicated in internal/engine/topics.go:34** instead of `deviceid.EncodedLen` (the canonical constant, pkg/deviceid/deviceid.go:26). Same value today, so no behaviour change; a rename to the shared constant is mechanical. Too small to be worth a task on its own, and it touches the hot topic-split path for cosmetic gain. Flagging here so the next review does not re-derive it.

- Nothing else rose to the bar: the trigger matcher (`match.go`), delivery executor (`actions.go`), and stream bus are well-constructed, heavily tested, and had no wrong/unguarded behaviour I could name with a file:line. The zip-bomb guard (§4.5), shard-parking backoff, policy disjointness rule, and per-requisite triggers all read correct and are machine-covered.

$ sed -n '30,200p' .mule/reviews/flow-2026-09-23.md
- `sort.go` — windowed ascending-timestamp sort + dedup.
- `tojson.go`, `splitmap.go`, `updatemetadata.go`, `catalog.go`, `info.go`,
  `schema.go` — catalog wire, field-splitting, metadata updates, registry
  JSON metadata.
- `astartesource/source.go`, `container/block.go` +
  `container/httpbridge.go`, `virtualdevicepool/virtualdevicepool.go` —
  bus-backed source, container block + bridge, device pool.
- Tests sampled: `httpblocks_test.go`, `transform_test.go`,
  `sort_test.go`, `randomsource_test.go`, `mqtt_test.go`, `jsonpathmap_test.go`,
  `astartesource/source_test.go`, `virtualdevicepool_test.go`.

## What I found

Four behavior issues (unguarded/wrong) and one test gap, in value order.

### 1. `mqtt_source` connection-lost latch is never cleared — a transient broker blip kills the source forever

`newMQTTClient` sets `SetAutoReconnect(true)` (mqtt.go:54) and the docstring
promises background auto-reconnect (mqtt.go:35-36). But the source latches
`lost = true` in the ConnectionLost handler (mqtt.go:185-187) and nothing
ever clears it: there is no `SetOnConnectHandler`, and Emit/Process gate on
the sticky flag (mqtt.go:234-236, 252-254). A single transient disconnect
therefore makes `mqtt_source` answer "connection lost" for the rest of the
flow's life, even after paho finishes reconnecting and re-subscribing — the
flow stays `running` while going silently dark. `TestMQTTSource_ConnectionLost`
(mqtt_test.go:190) only asserts the error surfaces; it never proves recovery.

Fix: clear `lost` on a successful reconnect (e.g. `SetOnConnectHandler`, or
check `client.IsConnected()` in Emit/Process) and extend the embedded-mochi
test to stop+restart the broker and assert Emit resumes. Container-free.

### 2. Message integer wire values silently lose precision / wrap (TypeInteger via float64)

`Message.UnmarshalJSON` (message.go:141-145) decodes with plain
`json.Unmarshal`, so every JSON number lands in `w.Data` as `float64`. The
`setDataFromWire` TypeInteger branch then does `m.Data = int64(v)`
(message.go:265-267) *after* the float64 rounding: an integer over 2^53
loses its low bits (`float64(MaxInt64)` rounds to 2^63 → `int64(2^63)`
wraps to MinInt64 on amd64), a fractional `3.7` silently becomes `3`, and
`1e300` becomes implementation-defined garbage. The `json.Number` branch
(message.go:268-273) — which would reject all of these — is effectively
dead in this path, because plain `json.Unmarshal` never produces
`json.Number`. Compare the already-accepted rule in `pkg/payload`
(`payload-longinteger-fraction-quantize`, todo line 170), which fixed the
same class for payload values; the message wire still has it. Roundtrip
tests only use `int64(42)` (message_test.go:19-22), so nothing catches it.

Fix: decode with `json.Decoder`+`UseNumber` (parse via `json.Number.Int64`)
or reject non-integral/out-of-int64 float64; add roundtrip rows for
`math.MaxInt64`, `1e300`, `3.7`. Container-free.

### 3. `Sort` buffers messages without bound when timestamps cluster in the window

`Sort` (sort.go:32-65) appends every message to an uncapped `buf` and
releases entries only while `len(buf) > 1 && buf[0].ts <= newest-windowUs`
(sort.go:60). A source that emits messages whose timestamps stay within
`window_ms` of the newest — a clustered burst, or a stream where every
message shares the same timestamp — never satisfies the flush condition, so
the block emits nothing and the slice grows without any limit: an unbounded
buffer grown from flow input, exactly the "unbounded slice from network
input" risk class. The documented "newest message is never emitted"
(sort.go:21-25) is implied by the same guard. Secondary: `windowUs :=
cfg.windowMs * 1000` (sort.go:31) overflows for a large-but-valid
`window_ms`, silently inverting the flush edge.

Fix: add a configurable maximum buffered count (overflow policy:
documented drop-oldest / error / force-flush), plus an injectable cap test
feeding many same-timestamp messages; guard the `window_ms * 1000`
overflow. Container-free.

### 4. `random_source` integer span overflows → `rand.Int64N` panics inside the pump goroutine

`randomSource.next()` computes
`rand.Int64N(s.maxInt-s.minInt+1)` (randomsource.go:171). `min`/`max` are
parsed from config as float64 then truncated to int64 (randomsource.go:92,
102) with only a `min <= max` check (randomsource.go:106-107) — not a span
check. A config like `min: -9e18, max: 9e18` (both valid int64) makes
`maxInt-minInt+1` overflow int64 to a non-positive value, and
`rand.Int64N` panics ("argument out of range"). The panic happens inside
the source pump goroutine (flow.go:205-239), which has no `recover`, so it
takes down the whole process. `randomsource_test.go` covers only
min==max (deterministic) and real/boolean shapes.

Fix: reject a span wider than `math.MaxInt64` at construction (alongside the
existing `min <= max` check) and add a construct/Emit test with the wide
span. Container-free.

### 5. `Filter.key_contains` is an existing rule with zero test coverage

`filterConfig.match` applies `key_contains` via `strings.Contains`
(transform.go:69), but `transform_test.go` only tests `key_prefix`
(TestFilter_KeyPrefixAndMetadata), `type`, unknown-type rejection and the
"at least one" construct rule (transform_test.go:12-71). A regression that
silently dropped the Contains check would pass the entire suite.

Fix: add substring-match, substring-absent, and combined
prefix+contains rows to the filter tests. Container-free, pure.

$ sed -n '1,40p' .mule/reviews/pairing-2026-09-09.md; echo "=== bench/test ==="; ls bench test cmd deploy examples
# Code review — internal/pairing (2026-09-09)

## What I read

- `internal/pairing/service.go` (440 lines) — flows A/B/C, `authenticateDevice`
  (uniform bcrypt equalization via `burnBcrypt`), `VerifyCredentials`
  (EXPIRED/INVALID/REVOKED classification), registration-limit gating.
- `internal/pairing/http.go` (367 lines) — the `/pairing/v1` route surface,
  `writeServiceError`, `bearerSecret`, `remoteIP`, 429 rate limiting.
- `internal/pairing/ratelimit.go` + `ratelimit_test.go` — LRU-cap-bounded token
  buckets.
- `internal/pairing/service_test.go` (637 lines) — fakeStore-backed unit tests.
- `internal/pairing/ca/ca.go` + `ca_test.go` — embedded per-realm CA: issue,
  load, verify, serial uniqueness.
- `internal/pairing/http_test.go` (integration-tagged, needs TimescaleDB —
  not run here).

## What I found (proposals)

This is a strongly-tested package (service, CA, and rate-limiter all have deep
container-free unit coverage; serial uniqueness even gets a 10k-issuance test).
Two genuinely untested rules remain, both pure helpers in the **non**-integration
`http.go`, so both are testable without containers:

1. **`bearerSecret` permissive matching** (`internal/pairing/http.go:343`).
   Astrate mirrors upstream's `~r/bearer\:?\s+(.*)$/i` — case-insensitive scheme,
   optional `:` separator, whitespace before the value. Only the plain
   `Bearer <secret>` form is exercised (integration `http_test.go`); the
   `Bearer: ` colon form and `bearer `/`BEARER ` case variants that upstream's
   regex accepts are never asserted anywhere, so a regression that tightens the
   match to case-sensitive or drops the colon would pass every existing test.
2. **`remoteIP` unparseable-address fallback** (`internal/pairing/http.go:359`).
   When `RemoteAddr` is neither a valid addr:port nor a bare address, it falls
   back to `netip.IPv4Unspecified()`. This branch (unix-socket/proxy case) is
   unreachable in the container-free tests and untested.

## What I decided *not* to propose

- **`handleCredentials` double token-bucket check** (`http.go:158`):
  `!Allow(ip) || !Allow(dev)` — when the IP bucket allows but the device bucket
=== bench/test ===
bench:
README.md
astrate
client.go
client_test.go
connstorm.go
device.go
go.mod
go.sum
ingest.go
interfaces
main.go
provision.go
query.go
scripts
state.go
stats.go

cmd:
astrate

deploy:
dashboard
devrealm

examples:
flow-container-echo
pokemon-agent

test:
conformance
security

$ rg -l 'bench/|test/conformance|test/security|tools/' .mule/reviews/ | head; echo "=== sizes ==="; rg -c '' --glob 'bench/**/*.go' | sort -t: -k2 -rn; echo "=== test dirs ==="; ls -R test | head -40
.mule/reviews/httpx-2026-09-17.md
.mule/reviews/channels-2026-09-18.md
=== sizes ===
bench/client.go:362
bench/ingest.go:349
bench/device.go:163
bench/provision.go:159
bench/connstorm.go:145
bench/client_test.go:134
bench/query.go:125
bench/stats.go:82
bench/main.go:71
bench/state.go:68
=== test dirs ===
test:
conformance
security

test/conformance:
atomvm
cpa
cpb
cpc
go.mod
go.sum
gosdk
instance
load
pysdk
upstream

test/conformance/atomvm:
atomvm_test.go

test/conformance/cpa:
astartectl_test.go
cpa_test.go

test/conformance/cpb:
cpb_test.go

test/conformance/cpc:
astartectl_test.go
cpc_test.go

test/conformance/gosdk:
gosdk_test.go

test/conformance/instance:
instance.go

test/conformance/load:
load_test.go

$ ls test/security; echo "=== sizes ==="; rg -c '' --glob 'test/security/**' | sort -t: -k2 -rn | head -20; echo "=== conformance load ==="; rg -c '' --glob 'test/conformance/load/**'
security_test.go
=== sizes ===
test/security/security_test.go:363
=== conformance load ===
test/conformance/load/load_test.go:205

→ Read bench/main.go
→ Read bench/client.go
→ Read bench/ingest.go
→ Read bench/device.go
→ Read bench/stats.go
→ Read bench/query.go
→ Read bench/connstorm.go
Error: Error from provider (Console): Rate limit exceeded. Please try again later.
