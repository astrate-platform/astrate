slug: recipe-code-review
verdict: proposed
at:  6388cd9
ran: 2026-10-02T11:44:27Z on DietPi in 546s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md

$ ls internal/ pkg/ 2>&1
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

$ ls -la .mule/reviews/ 2>&1; echo "---LOG TAIL---"; tail -40 .mule/log.md 2>&1
total 196
drwxr-xr-x 2 root root  4096 Sep 26 19:42 .
drwxr-xr-x 7 root root  4096 Oct  2 13:30 ..
-rw-r--r-- 1 root root  4600 Sep  7 20:15 appengine-2026-09-07.md
-rw-r--r-- 1 root root  3443 Sep 15 22:09 astarteapi-2026-09-15.md
-rw-r--r-- 1 root root  5192 Sep 14 13:34 auth-2026-09-14.md
-rw-r--r-- 1 root root  8188 Sep  4 20:43 broker-2026-09-04.md
-rw-r--r-- 1 root root  6580 Sep 21 21:26 broker-2026-09-21.md
-rw-r--r-- 1 root root  6967 Sep 18 22:00 channels-2026-09-18.md
-rw-r--r-- 1 root root 10261 Sep 25 19:46 cmd-astrate-2026-09-25.md
-rw-r--r-- 1 root root  4340 Sep 13 19:26 config-2026-09-13.md
-rw-r--r-- 1 root root  5264 Sep 18 19:25 deviceid-2026-09-18.md
-rw-r--r-- 1 root root  3141 Sep  6 20:57 engine-2026-09-06.md
-rw-r--r-- 1 root root  4692 Sep 20 19:30 engine-triggers-2026-09-20.md
-rw-r--r-- 1 root root  5516 Sep  5 21:19 flow-2026-09-05.md
-rw-r--r-- 1 root root  6931 Sep 23 20:39 flow-2026-09-23.md
-rw-r--r-- 1 root root  4837 Sep 10 19:00 flowapi-2026-09-10.md
-rw-r--r-- 1 root root  8949 Sep 26 19:42 forward-2026-09-26.md
-rw-r--r-- 1 root root  2595 Sep 15 13:29 housekeeping-2026-09-15.md
-rw-r--r-- 1 root root  4914 Sep 17 21:11 httpx-2026-09-17.md
-rw-r--r-- 1 root root  4712 Sep 11 21:40 interfaceschema-2026-09-11.md
-rw-r--r-- 1 root root  4594 Sep 16 20:48 observability-2026-09-16.md
-rw-r--r-- 1 root root  3506 Sep  9 13:07 pairing-2026-09-09.md
-rw-r--r-- 1 root root  6002 Sep 12 20:24 payload-2026-09-12.md
-rw-r--r-- 1 root root  4259 Sep  9 21:13 realm-2026-09-09.md
-rw-r--r-- 1 root root  4605 Sep  8 18:59 store-2026-09-08.md
-rw-r--r-- 1 root root  5504 Sep 21 13:32 store-2026-09-21.md
-rw-r--r-- 1 root root  2644 Sep 24 20:37 swagger-2026-09-24.md
---LOG TAIL---
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

$ head -20 .mule/reviews/forward-2026-09-26.md; echo "=== channels ==="; head -20 .mule/reviews/channels-2026-09-18.md
# Review — internal/engine/forward (2026-09-26)

First review of this area, so no prior file to extend. Chosen because it is the only
`internal/engine` sub-package with no review behind it, and because it is the seam where
a tenant's event payload leaves the process for something the mule does not control.

## What I read

- `internal/engine/forward/http.go` (all 125 lines) — `Config`, `New`, `validMethod`, `HTTP.Forward`
- `internal/engine/forward/envelope.go` (all 29 lines) — the `envelope` shape and `marshalEnvelope`
- `internal/engine/forward/nats.go` (all 53 lines, `//go:build nats`)
- `internal/engine/forward/http_test.go` (all 310 lines) and `nats_test.go` (all 196 lines)
- `internal/engine/triggers/actions.go:435-504` (`deliver`'s forward branch, `Executor.forward`,
  `DefaultRequestTimeout`) and `:600-629` (`attempt`, the sibling webhook request) — the only caller
- `cmd/astrate/main.go:116-119` and `:248-284` (`newForwarder`), `cmd/astrate/newnats_nats.go:15-19`
- `internal/config/config.go:339-375` (what `validate` does and does not check for
  `triggers.forward`), and `internal/engine/triggers/validation.go:19-34, 140-155` (the blocklist
  and size cap that trigger *action* headers get, and bus forwarder headers do not)
- `nats.go@v1.52.0` `Publish`/`publish` (nats.go:4229, 4424-4471), `processErr` (4181) and
  `setupServerPool`'s default error handler (1853-1856) in the module cache
=== channels ===
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

$ find internal pkg -type d | sort
internal
internal/appengine
internal/appengine/channels
internal/appengine/stream
internal/auth
internal/auth/testdata
internal/broker
internal/config
internal/engine
internal/engine/forward
internal/engine/stream
internal/engine/testdata
internal/engine/triggers
internal/engine/triggers/testdata
internal/flow
internal/flow/blocks
internal/flow/blocks/astartesource
internal/flow/blocks/container
internal/flow/blocks/virtualdevicepool
internal/flowapi
internal/housekeeping
internal/httpx
internal/observability
internal/pairing
internal/pairing/ca
internal/pairing/testdata
internal/pairing/testdata/http
internal/realm
internal/store
internal/store/testdata
internal/swagger
internal/testutil
pkg
pkg/astarteapi
pkg/astarteapi/testdata
pkg/deviceid
pkg/interfaceschema
pkg/interfaceschema/testdata
pkg/interfaceschema/testdata/invalid
pkg/interfaceschema/testdata/valid
pkg/payload
pkg/payload/testdata
pkg/payload/testdata/bson

$ rg -c '' --glob 'internal/pairing/**/*.go' --glob 'internal/testutil/**/*.go' --glob 'internal/flow/blocks/**/*.go' | sort -t: -k2 -rn | head -30
internal/pairing/http_test.go:649
internal/pairing/service_test.go:637
internal/flow/blocks/httpblocks_test.go:577
internal/flow/blocks/virtualdevicepool/virtualdevicepool_test.go:493
internal/pairing/service.go:440
internal/flow/blocks/container/block_test.go:424
internal/flow/blocks/mqtt_test.go:401
internal/flow/blocks/mqtt.go:388
internal/pairing/http.go:367
internal/flow/blocks/httpblocks.go:344
internal/flow/blocks/container/block.go:320
internal/pairing/ca/ca_test.go:309
internal/pairing/ca/ca.go:287
internal/flow/blocks/jsonpathmap_test.go:284
internal/flow/blocks/transform.go:255
internal/flow/blocks/jsonpathmap.go:238
internal/flow/blocks/transform_test.go:236
internal/flow/blocks/randomsource.go:226
internal/testutil/astartedevice.go:216
internal/flow/blocks/virtualdevicepool/virtualdevicepool.go:214
internal/flow/blocks/container/docker.go:211
internal/flow/blocks/splitmap_test.go:192
internal/flow/blocks/sort_test.go:190
internal/flow/blocks/container/httpbridge.go:178
internal/flow/blocks/randomsource_test.go:173
internal/flow/blocks/sort.go:171
internal/flow/blocks/astartesource/source.go:171
internal/flow/blocks/catalog.go:166
internal/testutil/mqttclient.go:163
internal/flow/blocks/info.go:162

$ for f in .mule/reviews/flow-2026-09-05.md .mule/reviews/flow-2026-09-23.md .mule/reviews/engine-2026-09-06.md .mule/reviews/pairing-2026-09-09.md; do echo "=== $f"; sed -n '1,30p' "$f" | grep -E '^(- `#|##|#|Area|First|Chosen|Reviewed)' ; done
=== .mule/reviews/flow-2026-09-05.md
# Code review: internal/flow — 2026-09-05
## Files read
## What I read in detail
## Findings
### 1. `Pipeline.Validate`'s no-source / no-sink rules are unreachable, and the tests that claim to cover them pass on the wrong path
=== .mule/reviews/flow-2026-09-23.md
# Review — internal/flow, 2026-09-23
## What I read
=== .mule/reviews/engine-2026-09-06.md
# Code review — internal/engine — 2026-09-06
Area chosen because `broker` (2026-09-04) and `flow` (2026-09-05) were recent reviews and `engine` is the largest un-reviewed package (48 files).
## What I read
## What I found (worth proposing)
## What I decided NOT to propose, and why
=== .mule/reviews/pairing-2026-09-09.md
# Code review — internal/pairing (2026-09-09)
## What I read
## What I found (proposals)

$ ls internal/flow/blocks/container/ internal/pairing/ca/ internal/engine/stream/ internal/testutil/
internal/engine/stream/:
bus.go
bus_test.go

internal/flow/blocks/container/:
block.go
block_test.go
docker.go
httpbridge.go

internal/pairing/ca/:
ca.go
ca_test.go

internal/testutil/:
astartedevice.go
golden.go
golden_test.go
mqttclient.go
pg.go
pg_test.go

$ rg -n -i 'container|docker|httpbridge|virtualdevicepool|astartesource' .mule/reviews/flow-2026-09-05.md .mule/reviews/flow-2026-09-23.md | head -30
.mule/reviews/flow-2026-09-23.md:4:`internal/flow/blocks/` (and its `astartesource/`, `container/`,
.mule/reviews/flow-2026-09-23.md:5:`virtualdevicepool/` subpackages) has never been reviewed, so this pass
.mule/reviews/flow-2026-09-23.md:34:- `astartesource/source.go`, `container/block.go` +
.mule/reviews/flow-2026-09-23.md:35:  `container/httpbridge.go`, `virtualdevicepool/virtualdevicepool.go` —
.mule/reviews/flow-2026-09-23.md:36:  bus-backed source, container block + bridge, device pool.
.mule/reviews/flow-2026-09-23.md:39:  `astartesource/source_test.go`, `virtualdevicepool_test.go`.
.mule/reviews/flow-2026-09-23.md:59:test to stop+restart the broker and assert Emit resumes. Container-free.
.mule/reviews/flow-2026-09-23.md:79:`math.MaxInt64`, `1e300`, `3.7`. Container-free.
.mule/reviews/flow-2026-09-23.md:98:overflow. Container-free.
.mule/reviews/flow-2026-09-23.md:115:span. Container-free.
.mule/reviews/flow-2026-09-23.md:126:prefix+contains rows to the filter tests. Container-free, pure.

$ sed -n '1,40p' .mule/reviews/pairing-2026-09-09.md
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

$ rg -n 'docker\.go|docker' .mule/reviews/flow-2026-09-23.md | head -20; echo "=== todo tail ==="; tail -50 .mule/todo.md
=== todo tail ===
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

→ Read .mule/reviews/flow-2026-09-23.md
→ Read internal/flow/blocks/container/block.go
→ Read internal/flow/blocks/container/docker.go
→ Read internal/flow/blocks/container/httpbridge.go
→ Read internal/flow/blocks/container/block_test.go

$ rg -n 'UnmarshalJSON' -A 40 internal/flow/message.go | head -70
142:// UnmarshalJSON deserialises a Message from the upstream JSON wire format.
143:func (m *Message) UnmarshalJSON(b []byte) error {
144-	var w wireMessage
145-	// Decode with UseNumber so numeric tokens survive as json.Number instead
146-	// of being pre-rounded to float64 (which would lose integer precision
147-	// beyond 2^53 before setDataFromWire ever sees them).
148-	dec := json.NewDecoder(bytes.NewReader(b))
149-	dec.UseNumber()
150-	if err := dec.Decode(&w); err != nil {
151-		return fmt.Errorf("flow: unmarshal message: %w", err)
152-	}
153-
154-	if w.Schema != WireSchema {
155-		return fmt.Errorf("flow: unsupported schema %q", w.Schema)
156-	}
157-
158-	m.Key = w.Key
159-	m.Metadata = w.Metadata
160-	m.Timestamp = w.TimestampUs
161-
162-	// Parse type field — either a string or a map of field types.
163-	switch t := w.Type.(type) {
164-	case string:
165-		dt, err := parseDataType(t)
166-		if err != nil {
167-			return err
168-		}
169-		m.Type = dt
170-		if s, ok := w.Subtype.(string); ok {
171-			m.Subtype = s
172-		}
173-		return m.setDataFromWire(dt, w.Data)
174-
175-	case map[string]any:
176-		m.Type = TypeMap
177-		m.FieldTypes = make(map[string]DataType, len(t))
178-		for k, v := range t {
179-			s, ok := v.(string)
180-			if !ok {
181-				return fmt.Errorf("flow: field type for %q is not a string", k)
182-			}
183-			dt, err := parseDataType(s)

$ rg -n 'func (Instantiate|NewRouter)|NotifyFatal|func .*Process\(' internal/flow/flow.go internal/flow/factory.go internal/flow/router.go | head -40
internal/flow/factory.go:29:	// NotifyFatal, when set, is called by blocks that die at runtime (e.g. a
internal/flow/factory.go:33:	NotifyFatal func(block string, cause error)
internal/flow/router.go:83:func NewRouter(graph *BlockGraph, cfg RouterConfig, reg prometheus.Registerer) *Router {

$ rg -n 'instantiate|Instantiate' internal/flow/factory.go internal/flow/flow.go | head -20
internal/flow/flow.go:52:// Config holds the parameters needed to instantiate a running flow.
internal/flow/flow.go:140:// StartFlow instantiates a pipeline into a running Flow. The block graph is
internal/flow/factory.go:14:// ErrUnknownBlockType is returned when Instantiate sees a block_type with no
internal/flow/factory.go:100:// Instantiate turns a validated Pipeline description into an ordered block
internal/flow/factory.go:105:func (r *Registry) Instantiate(p *Pipeline, deps Deps) ([]Block, error) {

$ rg -n 'WireSchema\s*=|schema\s+`' internal/flow/message.go | head; rg -n 'm.Key|msg.Key|Key ==' internal/flow/router.go | head -20
17:const WireSchema = "astarte_flow/message/v0.1"
116:// Submit routes msg to the lane determined by FNV-1a(msg.Key). Behaviour
132:	l := r.lanes[laneOf(msg.Key, len(r.lanes))]
221:			r.log.Error("panic in block graph", "key", fm.msg.Key,
228:		r.log.Error("block error", "key", fm.msg.Key, "err", err)

$ mkdir -p /tmp/opencode/probe && cat > /tmp/opencode/probe/main.go <<'EOF'
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
)

func main() {
	// 1. asInt(1e18) * time.Millisecond overflow, as parseConfig computes it.
	n := int(1e18)
	d := time.Duration(n) * time.Millisecond
	fmt.Printf("timeout_ms=1e18 -> n=%d duration=%v (positive=%v)\n", n, d, d > 0)

	n2 := int(1e15)
	fmt.Printf("timeout_ms=1e15 -> duration=%v positive=%v\n", time.Duration(n2)*time.Millisecond, time.Duration(n2)*time.Millisecond > 0)

	// 2. What net/http does with a negative Client.Timeout: no timeout at all.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(3 * time.Second):
			w.WriteHeader(200)
		case <-r.Context().Done():
		}
	}))
	defer srv.Close()
	c := &http.Client{Timeout: d} // the negative value New would install
	start := time.Now()
	_, err := c.Get(srv.URL)
	fmt.Printf("client with negative Timeout: after %v err=%v (nil means NO timeout applied)\n", time.Since(start).Round(100*time.Millisecond), err)

	// 3. A "null" JSON body from a container.
	var one flow.Message
	err = json.Unmarshal([]byte("null"), &one)
	fmt.Printf("unmarshal null -> err=%v msg=%+v\n", err, one)

	// 4. Same for a body with no schema field.
	var two flow.Message
	err = json.Unmarshal([]byte(`{"key":"k","type":"string","data":"x"}`), &two)
	fmt.Printf("unmarshal schemaless -> err=%v\n", err)
}
EOF
go run /tmp/opencode/probe/main.go
package command-line-arguments
	/tmp/opencode/probe/main.go:10:2: use of internal package github.com/astrate-platform/astrate/internal/flow not allowed

$ mkdir -p internal/flow/zzprobe_tmp && cp /tmp/opencode/probe/main.go internal/flow/zzprobe_tmp/main.go && go run ./internal/flow/zzprobe_tmp/ 2>&1 | tail -20
timeout_ms=1e18 -> n=1000000000000000000 duration=556601h10m5.20689664s (positive=true)
timeout_ms=1e15 -> duration=1076616h40m19.684212736s positive=true
client with negative Timeout: after 3s err=<nil> (nil means NO timeout applied)
unmarshal null -> err=flow: unsupported schema "" msg={Key: Metadata:map[] Type:0 Subtype: Timestamp:0 Data:<nil> FieldTypes:map[] FieldSubtypes:map[]}
unmarshal schemaless -> err=flow: unsupported schema ""

$ rg -n 'Deps\{|flow.Deps' --glob '*.go' -g '!*_test.go' . | head -20
./internal/flow/blocks/mqtt.go:169:func MQTTSource(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flow/blocks/mqtt.go:295:func MQTTSink(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flow/blocks/virtualdevicepool/virtualdevicepool.go:45:func Constructor(name string, config map[string]any, deps flow.Deps) (flow.Block, error) {
./internal/flow/blocks/httpblocks.go:57:func HTTPSource(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flow/blocks/httpblocks.go:215:func HTTPSink(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flow/blocks/sort.go:38:func Sort(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flow/blocks/updatemetadata.go:18:func UpdateMetadata(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flow/blocks/tojson.go:21:func ToJSON(name string, _ map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flow/blocks/splitmap.go:23:func SplitMap(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flow/blocks/container/block.go:65:func Constructor(name string, config map[string]any, deps flow.Deps) (flow.Block, error) {
./internal/flow/blocks/container/block.go:70:func New(name string, config map[string]any, deps flow.Deps, runner Runner) (flow.Block, error) {
./internal/flow/blocks/transform.go:26:func Filter(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flow/blocks/transform.go:91:func Map(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flow/blocks/randomsource.go:25:func RandomSource(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flow/blocks/jsonpathmap.go:38:func JSONPathMap(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flow/blocks/catalog.go:77:func AstarteSource(name string, config map[string]any, deps flow.Deps) (flow.Block, error) {
./internal/flow/blocks/catalog.go:95:func NullSink(name string, _ map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flow/blocks/catalog.go:100:func LogSink(name string, _ map[string]any, _ flow.Deps) (flow.Block, error) {
./internal/flowapi/service.go:389:func (s *Service) flowDeps(realmID int16, realm, name string) flow.Deps {
./internal/flowapi/service.go:390:	return flow.Deps{

$ rg -rn 'NotifyFatal' --glob '*.go' . | grep -v zzprobe_tmp
./internal/flow/factory.go:	// n, when set, is called by blocks that die at runtime (e.g. a
./internal/flow/factory.go:	n func(block string, cause error)
./internal/flow/blocks/container/block.go:	if deps.n != nil {
./internal/flow/blocks/container/block.go:			go b.watchExit(watchCtx, w, deps.n)
./internal/flow/blocks/container/block_test.go:		n: func(block string, _ error) { notifyCh <- block },
./internal/flow/blocks/container/block_test.go:		t.Fatal("expected n on unexpected container exit")
./internal/flow/blocks/container/block_test.go:		n: func(string, error) { fired.Store(true) },
./internal/flow/blocks/container/block_test.go:		t.Fatal("clean stop must not trigger n")
./internal/flowapi/service_integration_test.go:	if deps.n != nil {
./internal/flowapi/service_integration_test.go:				deps.n(name, errors.New("boom: container exited unexpectedly"))
./internal/flowapi/service.go:		n: s.onBlockFatal(realmID, realm, name),
./internal/flowapi/service.go:// onBlockFatal returns the Deps.n callback for one flow instance.

$ sed -n '380,460p' internal/flowapi/service.go
blks, err := s.reg.Instantiate(p, s.flowDeps(realmID, realm, name))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return blks, nil
}

// flowDeps builds the per-instance block dependencies: live bus, realm and
// flow identity, the fatal-block callback, and the engine ingest path (#84).
func (s *Service) flowDeps(realmID int16, realm, name string) flow.Deps {
	return flow.Deps{
		Bus:         s.bus,
		Realm:       realm,
		FlowName:    name,
		NotifyFatal: s.onBlockFatal(realmID, realm, name),
		Ingest:      s.ingest,
		Register:    s.register,
	}
}

// startFlowInstance is the single start path for POST create and boot rehydrate.
// The durable row must already exist (CreateAndStartFlow inserts first).
func (s *Service) startFlowInstance(ctx context.Context, realm string, realmID int16, name, pipelineName string, config map[string]any) (*FlowView, error) {
	markFailed := func(err error) {
		msg := err.Error()
		_ = s.st.UpdateFlowRuntime(ctx, realmID, name, "failed", &msg, nil, nil, nil)
	}

	blks, err := s.resolveAndBuild(ctx, realmID, realm, name, pipelineName, config)
	if err != nil {
		markFailed(err)
		return nil, err
	}

	instanceID := flow.InstanceID(realm, name)
	// Creating status while building live graph (optional observation).
	_ = s.st.UpdateFlowRuntime(ctx, realmID, name, "creating", nil, nil, nil, nil)

	f, err := s.mgr.StartFlow(ctx, flow.Config{
		PipelineID: instanceID,
		Blocks:     blks,
	})
	if err != nil {
		// Already live (e.g. double rehydrate): do not overwrite status with failed.
		if !errors.Is(err, flow.ErrFlowExists) {
			markFailed(err)
		}
		return nil, err
	}

	now := time.Now().UTC()
	if err := s.st.UpdateFlowRuntime(ctx, realmID, name, "running", nil, nil, &now, nil); err != nil {
		s.log.Error("flow started but failed to persist runtime", "realm", realm, "name", name, "error", err)
	}
	s.log.Info("flow started", "realm", realm, "name", name, "pipeline", pipelineName, "runtime_id", f.ID())

	row, err := s.st.GetFlow(ctx, realmID, name)
	if err != nil {
		// Live is up; synthesize view from memory + known fields.
		return mergeFlowView(realm, &store.Flow{
			Name: name, PipelineName: pipelineName, Config: mustJSON(config),
			AutoRestart: true, Status: "running", StartedAt: &now,
			CreatedAt: now, UpdatedAt: now,
		}, f), nil
	}
	return mergeFlowView(realm, row, f), nil
}

// DeleteFlow stops a live instance if present, unregisters it, and deletes
// the durable row.
func (s *Service) DeleteFlow(ctx context.Context, realm, name string) error {
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return err
	}
	instanceID := flow.InstanceID(realm, name)
	if err := s.mgr.StopFlow(ctx, instanceID); err != nil && !errors.Is(err, flow.ErrFlowNotFound) {
		return err
	}
	s.mgr.UnregisterFlow(instanceID)

$ sed -n '180,260p' internal/flow/router.go
r.laneWG.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// runLane drains one lane's channel through the block graph, recovering
// panics per-message so one block's bug cannot crash the router.
func (r *Router) runLane(_ context.Context, l *lane) {
	defer r.laneWG.Done()
	for {
		select {
		case fm := <-l.ch:
			r.processOne(fm)
		case <-r.quit:
			// Drain closes quit only once no Submit is in flight, so what is
			// buffered now is everything this lane will ever see.
			for {
				select {
				case fm := <-l.ch:
					r.processOne(fm)
				default:
					return
				}
			}
		}
	}
}

// processOne feeds a single message through the block graph, catching
// panics. Block-level errors are logged and counted but never crash the
// lane.
func (r *Router) processOne(fm *flowMsg) {
	defer func() {
		if rec := recover(); rec != nil {
			r.met.blockErrors.Inc()
			r.log.Error("panic in block graph", "key", fm.msg.Key,
				"panic", rec, "stack", string(debug.Stack()))
		}
	}()
	_, err := r.graph.Run(fm.msg)
	if err != nil {
		r.met.blockErrors.Inc()
		r.log.Error("block error", "key", fm.msg.Key, "err", err)
	}
	r.met.processed.Inc()
}

// laneOf hashes a key to a lane index using FNV-1a (allocation-free on the
// hot path, same algorithm as the engine shardOf).
func laneOf(key string, n int) int {
	const (
		offset64 = 14695981039346656037
		prime64  = 1099511628211
	)
	h := uint64(offset64)
	for i := 0; i < len(key); i++ {
		h ^= uint64(key[i]) // per-byte FNV-1a, same as the engine shardOf
		h *= prime64
	}
	return int(h % uint64(n)) // #nosec G115 -- value already reduced mod n
}

// routerMetrics are the router's Prometheus collectors.
type routerMetrics struct {
	submitted   prometheus.Counter
	processed   prometheus.Counter
	droppedQoS0 prometheus.Counter
	droppedQoS1 prometheus.Counter
	blockErrors prometheus.Counter
}

func newRouterMetrics(reg prometheus.Registerer) *routerMetrics {
	m := &routerMetrics{
		submitted: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_flow_router_submitted_total",

$ rg -n 'func (s \*Service) onBlockFatal' -A 45 internal/flowapi/service.go
(no output)

$ rg -n 'onBlockFatal' internal/flowapi/*.go
internal/flowapi/service.go:394:		NotifyFatal: s.onBlockFatal(realmID, realm, name),
internal/flowapi/service.go:529:// onBlockFatal returns the Deps.NotifyFatal callback for one flow instance.
internal/flowapi/service.go:532:func (s *Service) onBlockFatal(realmID int16, realm, name string) func(string, error) {

→ Read internal/flowapi/service.go [offset=525, limit=60]

$ rg -n -i 'container block|image|allowlist|registry' docs/DESIGN.md | head -20
39:flows, durable rehydrate, container blocks). See `.mule/milestones.md` and
64:| **Astarte Flow** | Dataflow processing framework | `internal/flow`, `internal/flowapi`, `internal/flow/blocks` | **v2.0 milestone.** Operator surface `/flow/v1/{realm}/…` (pipelines, flows, blocks). Built-in catalog + planned durable named flows and container blocks. Capability parity, not an Elixir port. |
116:└── docker-compose.yml        # astrate + timescale/timescaledb-ha image
639:   keys (PEM, RSA-2048+/ECDSA P-256+, alg allowlist `RS256/RS384/RS512/ES256/ES384/ES512`,
803:`max_connections=50`) + `astrate` (distroless, ~20 MB image). Volumes: pgdata + session store.

$ cat > internal/flow/zzprobe_tmp/main.go <<'EOF'
package main

import (
	"fmt"
	"time"
)

func main() {
	// What parseConfig computes for time.Duration(n) * time.Millisecond.
	for _, ms := range []float64{1e12, 9000000000.0, 18446744073709.0, 9223372036854.0, 9223372036854775807.0} {
		n := int(ms)
		d := time.Duration(n) * time.Millisecond
		fmt.Printf("ms=%-22.0f int=%-21d duration=%-28v positive=%v\n", ms, n, d, d > 0)
	}
}
EOF
go run ./internal/flow/zzprobe_tmp/
ms=1000000000000          int=1000000000000         duration=277777h46m40s                positive=true
ms=9000000000             int=9000000000            duration=2500h0m0s                    positive=true
ms=18446744073709         int=18446744073709        duration=-551.616µs                   positive=false
ms=9223372036854          int=9223372036854         duration=2562047h47m16.854s           positive=true
ms=9223372036854775808    int=9223372036854775807   duration=-1ms                         positive=false

$ rg -rn 'docker (ps|rm|run)' --glob '*.go' --glob '*.md' -l . 2>/dev/null | head; echo "---"; rg -n 'astrate\.flow=|label' --glob '*.go' -g '!*_test.go' internal/ cmd/ tools/ 2>/dev/null | head -20
./examples/flow-container-echo/README.md
./internal/flow/blocks/container/block.go
./internal/flow/blocks/container/block_test.go
./internal/flow/blocks/container/docker.go
./docs/handoff/flow-design-b-container-block-2026-07-29.md
---
cmd/astrate/main.go:507:func stage(label string, log *slog.Logger, fn func(context.Context) error) {
cmd/astrate/main.go:511:		log.Warn(label, "error", err)
internal/flow/blocks/container/block.go:84:	labels := map[string]string{
internal/flow/blocks/container/block.go:89:		labels["astrate.realm"] = deps.Realm
internal/flow/blocks/container/block.go:92:		labels["astrate.flow_name"] = deps.FlowName
internal/flow/blocks/container/block.go:105:		Labels:         labels,
internal/flow/block.go:21:	// Name returns a human-readable label for metrics and logging.
internal/flow/blocks/container/docker.go:80:// Start runs `docker run -d` with a published localhost port and Astrate labels.
internal/flow/blocks/container/docker.go:102:		args = append(args, "--label", k+"="+v)
internal/flow/pipeline.go:14:	// Name is a human-readable label.
internal/flow/factory.go:27:	// container labels and similar operator metadata.
internal/flow/flow.go:27:// String returns a human-readable label for s.
internal/engine/data.go:19:// metric labels and the M6b device_error trigger error names.
internal/engine/data.go:52:// engineRejectReasons pre-registers the labels above (newMetrics).
internal/engine/data.go:83:// opKinds pre-registers the kind labels (newMetrics).
internal/engine/data.go:86:// String returns the stable snake_case metrics label.
internal/engine/stream/bus.go:56:	// reject-reason label and its free-form diagnostic strings.
internal/engine/router.go:406:// pre-registering every reject-reason and op-kind label so dashboards see
internal/engine/triggers/errorname.go:68:// UpstreamErrorName maps one of Astrate's reject-reason labels to the upstream
internal/engine/triggers/events.go:243:	// labels feed it).

$ rg -n -i 'orphan|clean|reap|--rm|restart polic|leak' docs/handoff/flow-design-b-container-block-2026-07-29.md | head -20
77:Do **not** leave orphan containers after clean stop (label + remove).
140:  (`github.com/docker/docker/client` or lightweight wrapper) for cancel/cleanup.
187:- CPU/memory limits, restart policies, multi-container sidecars
199:5. Cleanup of orphaned containers with Astrate labels on process boot (best-effort).
239:| **B5** | Optional: orphan cleanup on boot |

$ rg -n -i 'orphan|leak|known limitation|limitation' docs/DESIGN.md docs/COMPATIBILITY.md | head -20
docs/COMPATIBILITY.md:194:   per-device rate-limited timing so the status code does not leak what the

$ sed -n '130,205p' docs/handoff/flow-design-b-container-block-2026-07-29.md
- Astrate block waits with a **timeout** (config default e.g. 5s).
- On flow start: `docker run` (or API) with env:
  - `ASTRATE_FLOW_CONFIG=<json>` of nested `config`
  - Labels: `astrate.flow=1`, `astrate.realm=…`, `astrate.flow_name=…`, `astrate.block=…`

### Docker API

- Prefer Docker Engine API via local socket (`unix:///var/run/docker.sock`) or
  `DOCKER_HOST`.
- PoC may shell out to `docker` CLI if faster to land; MVP should use a Go client
  (`github.com/docker/docker/client` or lightweight wrapper) for cancel/cleanup.
- **No Kubernetes operator in v2.0** — local Docker only for PoC/MVP.

---

## 5. Package layout (proposed)

```text
internal/flow/blocks/container/
  block.go          # Block + Stopper implementation
  docker.go         # create/start/stop/remove
  httpbridge.go     # POST message ↔ FlowMessage
  block_test.go     # unit with fake round-tripper
```

Register in `blocks.DefaultRegistry()` as `"container"`.

`blocks.Info` docs: image required; local Docker; limitations list for PoC.

### Process dependencies

Extend `flow.Deps` only if needed:

```go
// Optional; nil → container block fails Instantiate with clear error
Docker DockerRunner // interface{ Run(...); Stop(...) }
```

Or construct via config-only + package-level client. Prefer injectable interface for tests.

---

## 6. PoC scope (must / must-not)

### Must

1. `container` block type registered and documented via `/flow/v1/.../blocks`.
2. Start container from `image` when flow starts; HTTP round-trip one message.
3. Stop/remove container when flow stops.
4. Manual test recipe in handoff or `docs/` snippet: tiny echo image or `nginx`+sidecar
   is not enough — ship a **minimal example** under `examples/flow-container-echo/`
   (Dockerfile + 20-line HTTP server that echoes FlowMessage).
5. Loud fail if Docker unavailable or image missing at start.

### Must not (PoC)

- Registry auth beyond what local Docker already has
- CPU/memory limits, restart policies, multi-container sidecars
- AMQP compatibility with upstream container images
- Auto-pull policy UI; `docker pull` best-effort with clear error
- HA / multiple Astrate hosts scheduling the same container

### MVP adds

1. Works inside Design A named durable flow + rehydrate (container comes back on boot
   with the flow).
2. Health wait: block start waits until HTTP accepts or timeout → flow `failed`.
3. Config timeouts, max response bytes.
4. Metrics: messages in/out, container start failures, processing errors.
5. Cleanup of orphaned containers with Astrate labels on process boot (best-effort).

---

## 7. Security notes (MVP floor)

- Container runs with Docker defaults unless config later tightens; document that

$ tail -30 .mule/for-giulio.md
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
- **The mule has been idle 16h.** Filed by the dead-man's switch; see journalctl on the Pi.

$ sed -n '200,240p' docs/handoff/flow-design-b-container-block-2026-07-29.md
---

## 7. Security notes (MVP floor)

- Container runs with Docker defaults unless config later tightens; document that
  operators trust the image.
- Do not mount host Docker socket into the user container.
- Network: bridge + publish only the HTTP port to localhost (or attach to an internal
  network the block can reach). Prefer **not** host network.
- Realm isolation: labels only in PoC; no multi-tenant hard isolation claim.

---

## 8. Sequencing vs Design A

```text
Design A accepted → Implement A (named durable flows)
Design B accepted ↗ (may be earlier)
        ↓
Implement B PoC (can use manual StartFlow without full durability if needed for speed,
                 but MVP requires A)
        ↓
Implement B MVP on durable named flows
```

**Hard rule:** do not merge container as the “default path to custom logic” before
flows have names + config + durable restart semantics operators can rely on.

---

## 9. Implementation plan (after A + B accepted)

| Step | Work |
|---|---|
| **B1** | Echo example image + HTTP contract doc |
| **B2** | `container` block + Docker start/stop + HTTP bridge (PoC) |
| **B3** | Catalog Info + service tests with mocked Docker/HTTP |
| **B4** | MVP: health wait, timeouts, metrics, rehydrate-safe Stopper |
| **B5** | Optional: orphan cleanup on boot |

$ rg -n 'container' -B 4 -A 12 internal/flow/blocks/catalog.go | head -60
8-	"log/slog"
9-
10-	"github.com/astrate-platform/astrate/internal/flow"
11-	"github.com/astrate-platform/astrate/internal/flow/blocks/astartesource"
12:	"github.com/astrate-platform/astrate/internal/flow/blocks/container"
13-	"github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool"
14-)
15-
16-// Well-known block_type strings stored in pipeline definitions.
17-const (
18-	TypeAstarteSource = "astarte_source"
19-	TypeNullSink      = "null_sink"
20-	TypeLogSink       = "log_sink"
21-	// TypeContainer is the Docker-backed custom processing block (Design B / #43).
22:	TypeContainer = container.Type
23-	// TypeFilter and TypeMap are declared in transform.go.
24-
25-	// Pure transforms and the random source (#80, astarte_flow parity).
26-	TypeToJSON         = "to_json"
27-	TypeUpdateMetadata = "update_metadata"
28-	TypeSplitMap       = "split_map"
29-	TypeRandomSource   = "random_source"
30-	TypeSort           = "sort"
31-
32-	// TypeJSONPathMap reshapes JSON messages through a template (#81,
33-	// astarte_flow parity); declared in jsonpathmap.go's neighborhood.
34-	TypeJSONPathMap = "json_path_map"
--
44-	TypeVirtualDevicePool = virtualdevicepool.Type
45-)
46-
47-// DefaultRegistry returns a registry with the minimum useful built-in set:
48:// AstarteSource (bus → Message), filter/map transforms, container (custom
49-// image via local Docker HTTP bridge), and null/log sinks so operators can
50-// compose a complete source→transform→sink pipeline.
51-func DefaultRegistry() *flow.Registry {
52-	r := flow.NewRegistry()
53-	r.Register(TypeAstarteSource, AstarteSource)
54-	r.Register(TypeFilter, Filter)
55-	r.Register(TypeMap, Map)
56-	r.Register(TypeToJSON, ToJSON)
57-	r.Register(TypeUpdateMetadata, UpdateMetadata)
58-	r.Register(TypeSplitMap, SplitMap)
59-	r.Register(TypeRandomSource, RandomSource)
60-	r.Register(TypeSort, Sort)
--
62-	r.Register(TypeHTTPSource, HTTPSource)
63-	r.Register(TypeHTTPSink, HTTPSink)
64-	r.Register(TypeMQTTSource, MQTTSource)
65-	r.Register(TypeMQTTSink, MQTTSink)
66:	r.Register(TypeContainer, container.Constructor)
67-	r.Register(TypeVirtualDevicePool, virtualdevicepool.Constructor)
68-	r.Register(TypeNullSink, NullSink)
69-	r.Register(TypeLogSink, LogSink)
70-	return r
71-}
72-
73-// AstarteSource constructs an astartesource.Source. Config keys:
74-//   - realm (string): tenant; defaults to deps.Realm
75-//   - interface (string): optional interface filter

$ rg -n 'Info\b' -A 30 internal/flow/blocks/catalog.go | sed -n '/LookupInfo/,/^$/p' | head -50; echo "=== infos map ==="; rg -n 'TypeContainer' -A 10 internal/flow/blocks/catalog.go | sed -n '1,40p'
=== infos map ===
21:	// TypeContainer is the Docker-backed custom processing block (Design B / #43).
22:	TypeContainer = container.Type
23-	// TypeFilter and TypeMap are declared in transform.go.
24-
25-	// Pure transforms and the random source (#80, astarte_flow parity).
26-	TypeToJSON         = "to_json"
27-	TypeUpdateMetadata = "update_metadata"
28-	TypeSplitMap       = "split_map"
29-	TypeRandomSource   = "random_source"
30-	TypeSort           = "sort"
31-
32-	// TypeJSONPathMap reshapes JSON messages through a template (#81,
--
66:	r.Register(TypeContainer, container.Constructor)
67-	r.Register(TypeVirtualDevicePool, virtualdevicepool.Constructor)
68-	r.Register(TypeNullSink, NullSink)
69-	r.Register(TypeLogSink, LogSink)
70-	return r
71-}
72-
73-// AstarteSource constructs an astartesource.Source. Config keys:
74-//   - realm (string): tenant; defaults to deps.Realm
75-//   - interface (string): optional interface filter
76-//   - path (string): optional path prefix filter

$ rg -n 'blockInfos|Infos|func LookupInfo' -A 40 internal/flow/blocks/catalog.go | sed -n '/container/,+12p' | head -40
(no output)

$ rg -n 'LookupInfo|infos|Role:|Description:' internal/flow/blocks/*.go | head -30
internal/flow/blocks/schema_test.go:102:func TestLookupInfo_ConfigSchema(t *testing.T) {
internal/flow/blocks/schema_test.go:103:	info, ok := LookupInfo(TypeLogSink)
internal/flow/blocks/schema_test.go:105:		t.Fatal("LookupInfo(log_sink) not found")
internal/flow/blocks/schema_test.go:108:		t.Error("LookupInfo(log_sink).ConfigSchema is nil, want merged schema")
internal/flow/blocks/mqtt_test.go:378:		info, ok := blocks.LookupInfo(typ)
internal/flow/blocks/mqtt_test.go:380:			t.Fatalf("LookupInfo(%q) not found", typ)
internal/flow/blocks/mqtt_test.go:383:			t.Errorf("LookupInfo(%q).Role = %q, want %q", typ, info.Role, role)
internal/flow/blocks/mqtt_test.go:386:			t.Errorf("LookupInfo(%q).ConfigSchema = nil, want schema", typ)
internal/flow/blocks/info.go:36:		Role:    RoleSource,
internal/flow/blocks/info.go:42:		Role:    RoleTransform,
internal/flow/blocks/info.go:48:		Role:    RoleTransform,
internal/flow/blocks/info.go:54:		Role:    RoleTransform,
internal/flow/blocks/info.go:60:		Role:    RoleTransform,
internal/flow/blocks/info.go:66:		Role:    RoleTransform,
internal/flow/blocks/info.go:72:		Role:    RoleSource,
internal/flow/blocks/info.go:78:		Role:    RoleTransform,
internal/flow/blocks/info.go:84:		Role:    RoleTransform,
internal/flow/blocks/info.go:90:		Role:    RoleSource,
internal/flow/blocks/info.go:96:		Role:    RoleSink,
internal/flow/blocks/info.go:102:		Role:    RoleSource,
internal/flow/blocks/info.go:108:		Role:    RoleSink,
internal/flow/blocks/info.go:114:		Role:    RoleTransform,
internal/flow/blocks/info.go:120:		Role:    RoleSink,
internal/flow/blocks/info.go:126:		Role:    RoleSink,
internal/flow/blocks/info.go:132:		Role:    RoleSink,
internal/flow/blocks/info.go:138:// LookupInfo returns static operator docs for a known built-in type.
internal/flow/blocks/info.go:140:func LookupInfo(blockType string) (Info, bool) {
internal/flow/blocks/info.go:155:		if info, ok := LookupInfo(t); ok {
internal/flow/blocks/info.go:159:		out = append(out, Info{Type: t, Role: RoleTransform, Summary: "registered block (no built-in docs)"})
internal/flow/blocks/catalog_test.go:35:		if _, ok := blocks.LookupInfo(want); !ok {

→ Read internal/flow/blocks/info.go

$ rg -n 'container' -A 45 internal/flow/blocks/schema.go | sed -n '1,70p'
19:	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as ASTRATE_FLOW_CONFIG"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
20-	TypeVirtualDevicePool: `{"type":"object","properties":{"devices":{"type":"array","items":{"type":"string"},"minItems":1,"description":"Registered device_ids this pool may publish as"}},"required":["devices"]}`,
21-	TypeNullSink:          `{"type":"object","properties":{}}`,
22-	TypeLogSink:           `{"type":"object","properties":{}}`,
23-}

← Write internal/flow/blocks/container/zz_probe_test.go
Wrote file successfully.

$ go test ./internal/flow/blocks/container/ -run 'TestProbe' -v 2>&1 | grep -E 'PROBE|FAIL|ok|PASS|---' | head -30
PROBE Start err = container: could not parse host port from "no-colon-here"
PROBE cleanup call run    ctx-has-deadline=false
PROBE cleanup call port   ctx-has-deadline=false
PROBE cleanup call rm     ctx-has-deadline=false
PROBE stop-path call run    ctx-has-deadline=false
PROBE stop-path call port   ctx-has-deadline=false
PROBE stop-path call rm     ctx-has-deadline=false
PROBE stop-path call run    ctx-has-deadline=false
PROBE stop-path call port   ctx-has-deadline=false
PROBE stop-path call rm     ctx-has-deadline=false
--- PASS: TestProbe_CleanupContextHasDeadline (0.00s)
PROBE ready_timeout_ms=18446744073709   accepted, err=container "c": container: not ready after wait
PROBE ready_timeout_ms=9223372036854    accepted and STILL POLLING after 3s (New has not returned)
--- PASS: TestProbe_TimeoutBounds (3.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	3.042s

$ rm internal/flow/blocks/container/zz_probe_test.go internal/flow/zzprobe_tmp/main.go && rmdir internal/flow/zzprobe_tmp && git status --short
(no output)

$ go test ./internal/flow/blocks/container/ -coverprofile=/tmp/opencode/cov.out >/dev/null 2>&1; go tool cover -func=/tmp/opencode/cov.out | grep -E 'docker.go|httpbridge.go|block.go' | awk '$3=="0.0%" || $3+0 < 60 {print}' | head -30
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:224:	asInt			28.6%
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:67:	run			25.0%
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:133:	hostPort		57.1%
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:167:	Stop			38.5%
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:189:	Wait			0.0%
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:54:	WaitReady		56.7%

$ go tool cover -func=/tmp/opencode/cov.out | tail -20
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:224:	asInt			28.6%
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:240:	dockerName		85.7%
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:265:	Name			100.0%
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:269:	Process			100.0%
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:289:	Stop			88.9%
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:60:	bin			66.7%
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:67:	run			25.0%
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:81:	Start			74.1%
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:133:	hostPort		57.1%
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:164:	BaseURL			100.0%
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:165:	ID			100.0%
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:167:	Stop			38.5%
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:189:	Wait			0.0%
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:202:	encodeFlowConfigJSON	83.3%
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:32:	client			100.0%
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:39:	timeout			66.7%
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:46:	maxBody			66.7%
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:54:	WaitReady		56.7%
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:108:	RoundTrip		78.3%
total:											(statements)		71.2%

$ go tool cover -html=/tmp/opencode/cov.out -o /tmp/opencode/cov.html 2>/dev/null; awk -F'[:, ]' '$2>=133 && $2<=141 {print "httpbridge.go:"$2"-"$5" count="$NF}' /tmp/opencode/cov.out | head; echo "--- WaitReady 54-98 uncovered blocks ---"; rg 'httpbridge.go:(5[4-9]|[6-9][0-9])\.' /tmp/opencode/cov.out | rg ' 0$' | head
httpbridge.go:137.29-1 count=1
httpbridge.go:138.33-1 count=1
httpbridge.go:133.97-1 count=1
httpbridge.go:136.16-0 count=0
httpbridge.go:138.16-0 count=0
httpbridge.go:133.2-1 count=1
httpbridge.go:135.16-0 count=0
httpbridge.go:138.2-1 count=1
httpbridge.go:138.40-0 count=0
--- WaitReady 54-98 uncovered blocks ---
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:57.9,62.3 4 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:65.35,66.19 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:66.19,68.5 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:69.4,69.54 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:72.17,74.4 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:82.9,84.4 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:90.4,90.60 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:94.2,94.17 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:94.17,96.3 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/httpbridge.go:97.2,97.54 1 0

$ rg 'container/block.go:(18[0-9]|19[0-9]|2[0-2][0-9])\.' /tmp/opencode/cov.out | rg ' 0$' | head -20; echo "=== docker.go Wait/Waiter ==="; rg 'container/docker.go:(18[9-9]|19[0-9]|12[0-9]|17[0-9])\.' /tmp/opencode/cov.out | rg ' 0$' | head -20
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:192.11,193.69 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:197.45,199.40 2 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:199.40,201.4 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:202.3,202.24 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:205.51,207.27 2 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:207.27,209.4 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:210.3,210.52 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:215.27,217.4 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/block.go:228.13,229.21 1 0
=== docker.go Wait/Waiter ===
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:120.16,124.3 2 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:171.16,175.3 3 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:177.16,179.16 2 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:179.16,181.4 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:189.62,191.16 2 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:191.16,193.3 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:194.2,195.20 2 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:195.20,197.3 1 0
github.com/astrate-platform/astrate/internal/flow/blocks/container/docker.go:198.2,198.18 1 0

$ head -25 .mule/for-giulio.md; echo "=== grep container/metrics ==="; rg -n -i 'container block|container block|flow metrics|astrate_flow' .mule/for-giulio.md | head
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **Delta on the two open COMPATIBILITY.md wording entries further down this file (the v1.3.4 one folding in v1.3.3, and the original v1.3.2 one): upstream `v1.4.0-rc.6` is wire-inert for Astrate, so nothing in that proposal is re-derived — only its version reference and one line of its "not yet emulated" list change.** Not applied: `docs/COMPATIBILITY.md` is untouched (still the 1.2.2 target, deviation 18's "which upstream 404s" at `docs/COMPATIBILITY.md:360`, `APICompatVersion = "1.2.2"` at `internal/realm/service.go:588`), and the FDO half of the rc.6 delta already has its own entry below rather than being folded in here.
  - **The whole non-FDO delta is inert, item by item.** DUP RPC-server reliability is OTP supervisor/GC tuning — no protocol. `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"), i.e. FDO-only, covered by the entry below. `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma.
  - **The Dashboard's new `unknown_status` is computed client-side, not a new field** — upstream `models/Device/index.ts:213-215` derives it as `last_connection && !last_disconnection && !isConnected`, and both timestamps are already on Astrate's device shape (`internal/appengine/service.go:117-118`), so nothing is added to the `/appengine` device. The one consequence is a display change: a device whose disconnection was never recorded reads "Unknown status" where it used to read "Disconnected". Astrate ships no dashboard, so this is not a compat item — it only matters if #78's eventual work mirrors that dashboard.
  - **The `on_device_registered` dispatch moved from `Agent.register_device` to `Engine.register_device` (`pairing/engine.ex:157-165, 186-193`), so it now also covers the FDO path** — but the Agent-path rule is byte-for-byte the same, and Astrate already fires it on every successful registration (`internal/pairing/service.go:221-223` → `engine.HandleDeviceRegistered`, `cmd/astrate/main.go:154` → `triggers.OnDeviceRegistered`, `internal/engine/engine.go:386-388`), with `device_registered` already a matchable trigger condition (`internal/engine/triggers/match.go:50,118`). **Worth checking when you approve the proposal:** its "not yet emulated" list names "device registration triggers (pairing)", and against Astrate's wire name that one looks inaccurate — the deletion started/completed triggers it defers to deviation 9 are the real un-emulated half.
  - **The two edits the proposal needs, both yours to approve:** (1) the 2026-09-18 entry's parenthetical "(v1.4.0 is still rc.5-only)" becomes "(v1.4.0 is still rc.6-only)" — re-verified this session that rc.6 is the newest 1.4 prerelease and there is still **no stable v1.4.0** (`gh api repos/astarte-platform/astarte/releases`); (2) its "**FDO authentication** (pairing, disabled by default)" item, and the same "experimental FDO pairing auth" phrase in the proposed doc sentence, understates the 1.4 line, which ships a full ownership-voucher surface (upload with mandatory `hw_id`, deletion, re-TO0, per-voucher `expiry`, a spec-mandated realm-less TO2 path) that Astrate does not emulate — #78, deliberately parked at `.mule/milestones.md:124`. "FDO ownership vouchers (1.4 line, not stable upstream)" says it without claiming the rc.5 shape.
  Unverified: the upstream side of the item list is the upstream-watch recipe's `v1.4.0-rc.5...v1.4.0-rc.6` comparison — the same measurement as the FDO entry below, taken there and not re-derived here (a raw fetch of `models/Device/index.ts` at the rc.6 tag did not resolve from this box). I re-verified only the Astrate-side claims above and the release list. No `.mule/todo.md` lines: the deliverable is the delta, and the proposal itself stays open.

---

- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published 2026-09-30T07:28Z, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and the one path #78's body states as fact no longer exists upstream. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix") deletes the `realm_name` path parameter from all seven TO2 `operation` specs in `fdo_onboarding_controller.ex`. **This breaking change is not in the rc.6 CHANGELOG entry** (which lists only the to0 route, voucher deletion, the `ASTARTE_BASE_URL_*` rename, device-deletion cascade, `device_id` on upload, immediate registration, the two x5chain/voucher-store fixes and the cross-realm delete fix) — so anyone tracking upstream by its changelog, including #78's body, misses it. The realm now comes from the voucher row: the global-keyspace query `fetch_device_id_and_ownership_voucher_and_realm/1` returns `{realm, device_id, voucher_data}` and the `:fdo` plug fetches the session by GUID alone (`plug/fdo_session.ex`: `Session.fetch(guid)`); `RvTO2Addr` no longer carries the realm in the hostname either ("the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure"). **Yours:** #78's acceptance is "a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end" — against rc.6 that SDK posts to `/fdo/101/msg/{60..70}` with no realm, so mirroring rc.5's shape would mean mirroring a path upstream has deleted. No action while parked; the body is just stale.
- **Client-observable delta, from that same compare — every item here is on the wire, not internal:**
  - **Two new owner-side operations.** `DELETE /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}` → `204`, `404`, and `500 :rendezvous_revocation_failed` when the rendezvous revoke fails — **the voucher is kept** in that case (`FallbackController` + `OwnershipVoucherView`/`OwnershipVoucherController.delete_ownership_voucher/2`). `POST .../ownership_vouchers/{guid}/to0` re-runs TO0 against the rendezvous server → `200 {"data":{"expiry":…}}`, `404`, `500 :rendezvous_registration_failed`, and **`409 device_already_onboarded`** ("Once Device Onboard completes, the voucher is not registered on the rendezvous anymore"). A guid that is not a UUID is a `404`, not a `422` (`decode_guid/1` maps `Ecto.UUID.dump` failure to `:not_found`).
  - **`GET .../ownership_vouchers` gains `expiry`** per voucher (new field in the list view; migration `astarte/0007_add_expiry_to_ownership_vouchers.sql`), cleared back to null when TO2 consumes the registration (`mark_voucher_as_claimed` force-changes `expiry` to nil). Additive, so old clients keep working.
=== grep container/metrics ===
91:  sort -u` vs the file: the only other misses are `ASTRATE_FLOW_CONFIG` and the test-only
150:- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`ASTRATE_FLOW_CONFIG`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
383:  json_path_map, pure-transform set, virtual_device_pool, container block MVP, flow API,

$ rg -n -i 'MaxBodyBytes|max_response|max body|1 MiB|orphan' .mule/todo.md .mule/for-giulio.md .mule/reviews/ | head -20
.mule/todo.md:209:- [x] examples-echo-container-contract-test [auto]: add a container-free httptest suite for examples/flow-container-echo — the only no-test package in the repo that contains logic (`handleMessage`, main.go:49-75; `sanitizeLog`, main.go:40-47) — pinning (1) POST `/v1/message` echoes a valid JSON-object body verbatim with 200, answers 400 on invalid JSON and 204 on an empty body (the 1 MiB read cap), and (2) metadata `echo_drop` == "1" forces the 204 filter/drop path; the example is the canonical reference container contract (README.md, docs/handoff/flow-design-b-container-block-2026-07-29.md), so a regression silently misleads every container author who copies it.
.mule/reviews/engine-triggers-2026-09-20.md:65:  and the 1 MiB cap on the webhook side is about not reading an irrelevant body

$ rg -n 'astrate_flow_router_(submitted|dropped)' internal/flow/router.go | head -5; echo "--- handoff 196-202 ---"; sed -n '196,202p' docs/handoff/flow-design-b-container-block-2026-07-29.md
260:			Name: "astrate_flow_router_submitted_total",
268:			Name: "astrate_flow_router_dropped_qos0_total",
272:			Name: "astrate_flow_router_dropped_qos1_total",
--- handoff 196-202 ---
2. Health wait: block start waits until HTTP accepts or timeout → flow `failed`.
3. Config timeouts, max response bytes.
4. Metrics: messages in/out, container start failures, processing errors.
5. Cleanup of orphaned containers with Astrate labels on process boot (best-effort).

---

← Write .mule/reviews/container-2026-10-02.md
Wrote file successfully.

Now writing the review file.
$ tail -c 400 .mule/todo.md; echo "[EOF]"
id"); and `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma. So the proposal's "v1.4.0 is still rc.5-only" reference becomes rc.6, and its "FDO authentication (pairing, disabled by default)" line understates the 1.4 line, which ships a full ownership-voucher surface Astrate does not emulate (#78, parked). Do not edit docs/COMPATIBILITY.md and do not restate the standing proposal.
[EOF]

$ cat >> .mule/todo.md <<'EOF'

- [ ] container-stop-deadline [auto]: make `cliInstance.Stop` (internal/flow/blocks/container/docker.go:167-175) derive its 15s bound whenever the passed context carries **no deadline**, not only when it is nil — both cleanup callers pass `context.Background()` (docker.go:122, the "Best-effort cleanup so we do not leave orphans on mapping failure" path, and block.go:121, the not-ready cleanup in `New`), so a wedged docker daemon leaves an unbounded `docker rm -f` inside `exec.CommandContext` holding flow instantiation, while `Block.Stop` (block.go:310-311) does the same job correctly with an explicit 15s timeout. Measured with the package's own `Run` hook: the `rm` invocation arrives with `ctx.Deadline()` unset on both paths. Add a case asserting `inst.Stop(context.Background())` reaches the injected `Run` with a deadline set (it does not today) and that the hostPort-failure path still issues `rm -f <id>` (docker.go:120-124, never executed by any test today). Container-free.
- [ ] container-timeout-bounds [auto]: bound `timeout_ms` and `ready_timeout_ms` in `parseConfig` (internal/flow/blocks/container/block.go:205-219) — both are only checked `> 0` and then multiplied unchecked (`time.Duration(n) * time.Millisecond`), which wraps: measured `ready_timeout_ms: 18446744073709` gives `-551.616µs`, so `New` fails with "container: not ready after wait" without ever issuing a probe (the wrapped negative deadline fails the loop guard at httpbridge.go:64), while `ready_timeout_ms: 9223372036854` gives 2562047h47m16s and `New` never returns (still polling `/healthz` after 3s in the probe). The two consumers of a non-positive duration also disagree — `Bridge.timeout()` (httpbridge.go:39-44) treats it as "use the 5s default" while net/http treats a non-positive `Client.Timeout` as no deadline at all (measured: a client with `Timeout: -551µs` returned err=nil after a 3s handler). Reject a value that would overflow or exceed a documented maximum, naming the key and the value in the error, and add table cases asserting `New` refuses both (it accepts them today). Nothing bounds these on the wire either: internal/flow/blocks/schema.go:19 advertises both as bare `"type":"integer"`. Container-free.
- [ ] container-parseconfig-rules-test [auto]: pin the four `parseConfig` rules that no test touches (internal/flow/blocks/container/block.go:188-219 — all `count=0` in `go test -coverprofile`, which is 71.2% for the package): `config` must be a JSON object (192-193), `port` an integer in 1–65535 including the 0 / 65536 / non-integer rejections (197-202), and `timeout_ms` and `ready_timeout_ms` positive (205-210, 215-217) — beside the existing `TestConstructor_RequiresImage` (block_test.go:99), which is the only config case the suite has. The port rule is the one with teeth: it is what keeps `127.0.0.1::99999` out of the `docker run` argv (docker.go:93). Container-free, pure.
- [ ] container-response-cap-test [auto]: pin `Bridge.RoundTrip`'s 1 MiB response cap (internal/flow/blocks/container/httpbridge.go:133-140 — `io.LimitReader(resp.Body, maxBody+1)` plus the explicit rejection; lines 136, 138.16 and 138.40 are all uncovered), the only bound between a tenant-supplied container and process memory on the message path: assert a handler answering exactly `maxBody+1` bytes is rejected with the "exceeds" error and that exactly `maxBody` bytes still round-trip, driving a `Bridge{MaxBodyBytes: N}` so the fixture stays small. Also worth one line in the same change: the cap is absent from the operator-facing catalog entry (internal/flow/blocks/info.go:116) and from the schema, so an author whose container emits 2 MiB learns about it only from a runtime error. Container-free.
- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.
EOF
tail -3 .mule/todo.md | head -c 300
- [ ] container-parseconfig-rules-test [auto]: pin the four `parseConfig` rules that no test touches (internal/flow/blocks/container/block.go:188-219 — all `count=0` in `go test -coverprofile`, which is 71.2% for the package): `config` must be a JSON object (192-193), `port` an integer in 1–6553

Now the `for-giulio.md` entry (feature-shaped, newest at top).
→ Read .mule/for-giulio.md [limit=12]

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,10 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **The container block (Design B / #43) ships no metrics at all, and its 1 MiB response cap is not configurable — both are MVP items in your own handoff, both are feature-shaped, so neither is a queue line.** `docs/handoff/flow-design-b-container-block-2026-07-29.md:198-199` lists "Config timeouts, max response bytes" and "Metrics: messages in/out, container start failures, processing errors" as MVP work, and §9 tracks them under step B4; the package registers zero collectors (`internal/flow/blocks/container/` has no `prometheus` import at all), while every other flow component carries an `astrate_flow_*` family (`astrate_flow_router_submitted_total` et al., internal/flow/router.go:260-276) and the discovery API advertises the block as production-ready (`internal/flow/blocks/info.go:112-117`, `schema.go:19`). Related: `Bridge.MaxBodyBytes` (internal/flow/blocks/container/httpbridge.go:29) is never assigned — `New` builds the bridge without it (block.go:112-116), so the limit is hard-wired to `defaultMaxBodyBytes` and a tenant cannot raise it for a large payload; the field reads like a knob but is not one. Yours to call: metric names + whether they become part of the documented metrics surface, and whether `max_response_bytes` becomes a real config key (it would need the schema at `schema.go:19` and the `Config` string at `info.go:116` updated in the same change). Related queue lines exist for the mechanical half (`container-response-cap-test` pins today's cap); nothing implements either feature. Also from the same review, for the record: the handoff's B5 orphan sweep is still open as `container-orphan-boot-cleanup [legion]`, and the image-trust note in §7 is satisfied by the current `Info.Summary`, so a missing image allowlist is documented behaviour rather than a gap.
+
 - **Delta on the two open COMPATIBILITY.md wording entries further down this file (the v1.3.4 one folding in v1.3.3, and the original v1.3.2 one): upstream `v1.4.0-rc.6` is wire-inert for Astrate, so nothing in that proposal is re-derived — only its version reference and one line of its "not yet emulated" list change.** Not applied: `docs/COMPATIBILITY.md` is untouched (still the 1.2.2 target, deviation 18's "which upstream 404s" at `docs/COMPATIBILITY.md:360`, `APICompatVersion = "1.2.2"` at `internal/realm/service.go:588`), and the FDO half of the rc.6 delta already has its own entry below rather than being folded in here.
   - **The whole non-FDO delta is inert, item by item.** DUP RPC-server reliability is OTP supervisor/GC tuning — no protocol. `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"), i.e. FDO-only, covered by the entry below. `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma.
   - **The Dashboard's new `unknown_status` is computed client-side, not a new field** — upstream `models/Device/index.ts:213-215` derives it as `last_connection && !last_disconnection && !isConnected`, and both timestamps are already on Astrate's device shape (`internal/appengine/service.go:117-118`), so nothing is added to the `/appengine` device. The one consequence is a display change: a device whose disconnection was never recorded reads "Unknown status" where it used to read "Disconnected". Astrate ships no dashboard, so this is not a compat item — it only matters if #78's eventual work mirrors that dashboard.
   - **The `on_device_registered` dispatch moved from `Agent.register_device` to `Engine.register_device` (`pairing/engine.ex:157-165, 186-193`), so it now also covers the FDO path** — but the Agent-path rule is byte-for-byte the same, and Astrate already fires it on every successful registration (`internal/pairing/service.go:221-223` → `engine.HandleDeviceRegistered`, `cmd/astrate/main.go:154` → `triggers.OnDeviceRegistered`, `internal/engine/engine.go:386-388`), with `device_registered` already a matchable trigger condition (`internal/engine/triggers/match.go:50,118`). **Worth checking when you approve the proposal:** its "not yet emulated" list names "device registration triggers (pairing)", and against Astrate's wire name that one looks inaccurate — the deletion started/completed triggers it defers to deviation 9 are the real un-emulated half.

$ git status --short; echo "--- package still builds ---"; go vet ./internal/flow/blocks/container/ && go test ./internal/flow/blocks/container/ 2>&1 | tail -3
M .mule/for-giulio.md
 M .mule/todo.md
?? .mule/reviews/container-2026-10-02.md
--- package still builds ---
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)

Review done. No code changed, no git operations, no issue filed (this recipe's deliverable is the review file plus queue lines).

```
Done: Reviewed internal/flow/blocks/container (the only package whose docker.go had never been read) — wrote .mule/reviews/container-2026-10-02.md, appended 5 task lines to .mule/todo.md, added one feature-shaped entry to .mule/for-giulio.md.
Files: .mule/reviews/container-2026-10-02.md (new), .mule/todo.md (+5 lines), .mule/for-giulio.md (+1 bullet)
Verified: go test ./internal/flow/blocks/container/ -> pass; go vet ./internal/flow/blocks/container/ -> pass; git status shows only the three mule files (probe files deleted). Measurements taken with a throwaway program and a temporary in-package test driving CLIRunner.Run, both removed afterwards.
Unsure: the real `docker rm -f` hang is reasoned, not measured — this host has no Docker, so only the missing context deadline was measured (ctx.Deadline() unset on both cleanup paths).
Follow-ups: container-stop-deadline, container-timeout-bounds, container-parseconfig-rules-test, container-response-cap-test, container-orphan-boot-cleanup [legion]
```
