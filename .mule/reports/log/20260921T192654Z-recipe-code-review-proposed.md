slug: recipe-code-review
verdict: proposed
at:  9de089e
ran: 2026-09-21T19:26:54Z on DietPi in 997s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md

$ ls internal/ pkg/ 2>/dev/null; echo "---REVIEWS---"; ls -la .mule/reviews/ 2>/dev/null; echo "---LOG TAIL---"; tail -40 .mule/log.md 2>/dev/null
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
---REVIEWS---
total 152
drwxr-xr-x 2 root root 4096 Sep 21 13:32 .
drwxr-xr-x 7 root root 4096 Sep 21 20:20 ..
-rw-r--r-- 1 root root 4600 Sep  7 20:15 appengine-2026-09-07.md
-rw-r--r-- 1 root root 3443 Sep 15 22:09 astarteapi-2026-09-15.md
-rw-r--r-- 1 root root 5192 Sep 14 13:34 auth-2026-09-14.md
-rw-r--r-- 1 root root 8188 Sep  4 20:43 broker-2026-09-04.md
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
---LOG TAIL---
| 2026-09-12 | docs-sync-ae-post-groups-409 | done | 128s | d6865a1 |
| 2026-09-12 | payload-longinteger-fraction-quantize | done | 356s | d7babc6 |
| 2026-09-12 | payload-json-malformed-t-not-tolerated | done | 99s | 5b56557 |
| 2026-09-12 | docs-sync-ae-post-group-devices-409 | done | 179s | cbf81af |
| 2026-09-12 | docs-sync-ae-patch-merge-patch-media-type | done | 137s | 64361d9 |
| 2026-09-13 | store-todo-lttb-stale | done | 66s | 829f408 |
| 2026-09-13 | config-fail-loud-engine-shards | done | 150s | 83482e4 |
| 2026-09-13 | config-fail-loud-dev-mode | done | 227s | 27d8165 |
| 2026-09-13 | config-key-resolvers-test | done | 143s | 9e68a01 |
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
| 2026-09-21 | store-pipelines-empty-name-zero-blocks-test | done | 91s | 82fc0a4 |

→ Read .mule/todo.md

$ for p in swagger testutil testutil/*; do echo "=== $p ==="; rg -l '' --glob "$p/**/*.go" 2>/dev/null | while read f; do wc -l "$f"; done; done; echo "=== channels ==="; find internal -type d -name channels; echo "=== store subdirs ==="; ls internal/store/
=== swagger ===
=== testutil ===
=== testutil/* ===
=== channels ===
internal/appengine/channels
=== store subdirs ===
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

→ Read .mule/reviews/store-2026-09-21.md

$ for p in broker pairing realm flow flowapi engine appengine config auth housekeeping httpx observability swagger testutil; do echo "=== $p ==="; rg -c '' --glob "internal/$p/**/*.go" 2>/dev/null | sort -t: -k2 -rn | head -5; done; echo "=== channels review ==="; head -30 .mule/reviews/channels-2026-09-18.md
=== broker ===
internal/broker/broker_e2e_test.go:618
internal/broker/broker_test.go:549
internal/broker/sessionstore.go:533
internal/broker/authhook.go:480
internal/broker/broker.go:343
=== pairing ===
internal/pairing/http_test.go:649
internal/pairing/service_test.go:637
internal/pairing/service.go:440
internal/pairing/http.go:367
internal/pairing/ca/ca_test.go:309
=== realm ===
internal/realm/http_test.go:782
internal/realm/service.go:686
internal/realm/http.go:569
internal/realm/dashboard_compat_test.go:213
internal/realm/helpers_test.go:136
=== flow ===
internal/flow/blocks/httpblocks_test.go:577
internal/flow/router_test.go:565
internal/flow/blocks/virtualdevicepool/virtualdevicepool_test.go:493
internal/flow/blocks/container/block_test.go:424
internal/flow/blocks/mqtt_test.go:401
=== flowapi ===
internal/flowapi/service.go:992
internal/flowapi/service_integration_test.go:334
internal/flowapi/http.go:310
internal/flowapi/userblocks.go:209
internal/flowapi/service_test.go:167
=== engine ===
internal/engine/triggers/actions_test.go:1139
internal/engine/fakestore_test.go:818
internal/engine/triggers/match.go:769
internal/engine/triggers/actions.go:629
internal/engine/triggers/match_test.go:594
=== appengine ===
internal/appengine/channels/ws_test.go:1005
internal/appengine/service.go:846
internal/appengine/http.go:696
internal/appengine/data.go:514
internal/appengine/http_test.go:512
=== config ===
internal/config/config_test.go:584
internal/config/config.go:420
=== auth ===
internal/auth/middleware_test.go:284
internal/auth/jwt_test.go:268
internal/auth/claims_test.go:249
internal/auth/claims.go:204
internal/auth/middleware.go:196
=== housekeeping ===
internal/housekeeping/http_test.go:511
internal/housekeeping/service.go:288
internal/housekeeping/service_test.go:277
internal/housekeeping/http.go:259
internal/housekeeping/async_operation_test.go:37
=== httpx ===
internal/httpx/notfound_test.go:139
internal/httpx/cors_test.go:122
internal/httpx/notfound.go:56
internal/httpx/cors.go:52
=== observability ===
internal/observability/observability_test.go:180
internal/observability/health.go:84
internal/observability/metrics.go:80
internal/observability/compat_version_test.go:55
internal/observability/compat.go:52
=== swagger ===
internal/swagger/swagger_test.go:128
internal/swagger/swagger.go:43
=== testutil ===
internal/testutil/astartedevice.go:216
internal/testutil/mqttclient.go:163
internal/testutil/golden_test.go:117
internal/testutil/pg.go:106
internal/testutil/golden.go:59
=== channels review ===
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

Store was reviewed twice already (08 and 21); the most-stale real package is `broker` (2026-09-04). Let me check that prior review and probe the current state.
→ Read .mule/reviews/broker-2026-09-04.md
Broker is the most stale (09-04); store was just done twice. Let me read the current broker files, focusing on areas the prior review didn't cover deeply (intake, publisher, lifecycle, identity, and the current ACL/session state after the blocked fixes).
→ Read internal/broker/publisher.go
→ Read internal/broker/intake.go
→ Read internal/broker/identity.go
→ Read internal/broker/aclhook.go
→ Read internal/broker/authhook.go
→ Read internal/broker/broker.go
→ Read internal/broker/lifecycle.go
Now let me check what changed in broker since the last review (2026-09-04) and read the sessionstore.
→ Read internal/broker/sessionstore.go

$ git log --oneline --since="2026-09-04" -- internal/broker/ | head -30
24ad5b8 mule: issue-93: In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
d201db4 mule: broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [approved 2026-09-04]

→ Read .mule/reviews/flow-2026-09-05.md
Store was just reviewed today and re-reviewed; the most-stale area is `broker` (09-04) with two new commits since (coldstart miss fix + keyAgreement comment), then `flow` (09-05). Let me check the flow review coverage and the new broker test code before deciding.
✱ Grep "syncOwnershipOf|TestColdStart|coldstart|cold-start" in internal/broker · 6 matches
The most-stale areas are broker (09-04) and engine (09-06). Let me check the engine review to see how deep it went, since broker's known gaps are already queued.
→ Read .mule/reviews/engine-2026-09-06.md

$ git log --oneline --since="2026-09-06" -- internal/engine/ | head -40; echo "==="; ls internal/engine/ | head -50
3c8b89a mule: triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
739df67 mule: errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
e46ce11 mule: server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.
===
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

Broker is the most-stale area (09-04) and has fresh code since (the coldstart-miss fix, keyAgreement comment). Flow's review is recent relative to it and its finding is already queued. I'll go deep on broker again, focusing on what changed and what the 09-04 review skimmed (identity.go, publisher.go, lifecycle.go, intake, and the new coldstart fallback). Let me read the tests that pin the new code.
→ Read internal/broker/aclhook_test.go
→ Read internal/broker/identity_test.go
✱ Grep "LastIntroLoad|lastIntroLoad|ColdStart|coldstart|IntrospectionMiss|introspection-miss" in internal/broker · 10 matches

$ rg -n "func Test" internal/broker/broker_test.go internal/broker/sessionstore_test.go internal/broker/broker_e2e_test.go | head -40
internal/broker/broker_test.go:185:func TestNewValidation(t *testing.T) {
internal/broker/broker_test.go:240:func TestBrokerInMemorySmoke(t *testing.T) {
internal/broker/broker_test.go:375:func TestBrokerClientIDRemappedToCertCN(t *testing.T) {
internal/broker/broker_test.go:434:func TestBrokerACLColdStartIntrospectionMiss(t *testing.T) {
internal/broker/broker_test.go:499:func TestReloadRealmsPicksUpNewRealm(t *testing.T) {
internal/broker/sessionstore_test.go:45:func TestSessionStoreInitRejectsConfig(t *testing.T) {
internal/broker/sessionstore_test.go:52:func TestSessionStoreClientAndSubscriptionRoundTrip(t *testing.T) {
internal/broker/sessionstore_test.go:99:func TestSessionStoreTakeoverKeepsSession(t *testing.T) {
internal/broker/sessionstore_test.go:115:func TestSessionStoreInflightOrderingAndExpiry(t *testing.T) {
internal/broker/sessionstore_test.go:172:func TestSessionStoreRetainedRoundTrip(t *testing.T) {
internal/broker/sessionstore_test.go:196:func TestNormalizeInflightOrderPerClient(t *testing.T) {
internal/broker/broker_e2e_test.go:221:func TestBrokerE2E(t *testing.T) {

→ Read internal/broker/broker_test.go [offset=240, limit=280]
Now let me examine the interactions the prior review didn't reach — the new coldstart fallback cache semantics and the Ownership enum values.
✱ Grep "Ownership" in pkg/interfaceschema · 43 matches
✗ Read pkg/interfaceschema/ownership.go failed
Error: File not found: /root/astrate-mule/pkg/interfaceschema/ownership.go
✱ Grep "isLatestCert|enforceLatestCert|EnforceLatestCert|TestIsLatestCert" in internal/broker · 11 matches
Let me verify a couple of claims before writing the review — the `isLatestCert` coverage and the exact pre-fix behavior.

$ rg -n "lastIntroLoad|syncOwnershipOf|refreshIfStale" internal/broker/*_test.go; echo "---"; rg -n "TestCheckACLMatrix|offlineACL" internal/broker/*_test.go; echo "--- retained test ---"; rg -n "Retained" internal/broker/*_test.go | head
internal/broker/broker_test.go:428:// stamps lastIntroLoad at connect, so refreshIfStale's debounce skips the
internal/broker/broker_test.go:431:// silently by mochi for QoS0. lastIntroLoad is re-stamped just before the
internal/broker/broker_test.go:487:	sess.lastIntroLoad = time.Now()
---
internal/broker/aclhook_test.go:10:// TestCheckACLMatrix is the exhaustive §3.2 table test (docs/ROADMAP.md §6):
internal/broker/aclhook_test.go:13:func TestCheckACLMatrix(t *testing.T) {
--- retained test ---
internal/broker/sessionstore_test.go:172:func TestSessionStoreRetainedRoundTrip(t *testing.T) {
internal/broker/sessionstore_test.go:182:	msgs, err := ss.StoredRetainedMessages()
internal/broker/sessionstore_test.go:184:		t.Fatalf("StoredRetainedMessages: %v", err)
internal/broker/sessionstore_test.go:187:		t.Fatalf("StoredRetainedMessages = %+v, want only r/keep", msgs)
internal/broker/sessionstore_test.go:190:	ss.OnRetainedExpired("r/keep")
internal/broker/sessionstore_test.go:191:	if msgs, _ = ss.StoredRetainedMessages(); len(msgs) != 0 {
internal/broker/sessionstore_test.go:192:		t.Fatalf("after OnRetainedExpired: %+v, want empty", msgs)

Let me check `SplitTopic` callers and the pre-fix ACL miss code, since the new fallback is the most interesting new seam.
$ rg -n "SplitTopic" internal/ cmd/ pkg/ --glob '*.go' | rg -v "_test"; echo "=== pre-fix miss path ==="; git show d201db4^:internal/broker/aclhook.go | sed -n '175,215p'
internal/broker/identity.go:67:// SplitTopic splits a wire topic into its realm, device, and rest parts
internal/broker/identity.go:72:func SplitTopic(topic string) (realm, device, rest string, err error) {
=== pre-fix miss path ===
func (h *aclHook) ID() string { return "astrate-acl" }

// Provides implements mqtt.Hook.
func (h *aclHook) Provides(b byte) bool { return b == mqtt.OnACLCheck }

// OnACLCheck implements mqtt.Hook.
func (h *aclHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	if cl.Net.Inline {
		return true // engine/AppEngine publishes bypass ACLs (§3.2)
	}

	var allowed bool
	if sess := h.registry.get(cl.ID); sess != nil && sess.client == cl {
		allowed = checkACL(sess.identity.BaseTopic(), topic, write, func(iface string) (interfaceschema.Ownership, bool) {
			own, ok := sess.ownershipOf(iface)
			if !ok {
				// The interface may have been introspected after connect
				// (the engine updates the row); reload, debounced.
				ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
				defer cancel()
				sess.refreshIfStale(ctx, h.st, h.log)
				own, ok = sess.ownershipOf(iface)
			}
			return own, ok
		})
	} else if !write {
		// Delivery to a session-present but disconnected device (offline
		// queue, retained messages): no live session, consult the store.
		allowed = checkACL(cl.ID, topic, write, func(iface string) (interfaceschema.Ownership, bool) {
			return h.offline.ownershipOf(cl.ID, iface)
		})
	}

	if !allowed {
		h.log.Warn("mqtt ACL denied",
			"client", cl.ID, "topic", topic, "write", write, "remote", cl.Net.Remote)
	}
	return allowed
}

→ Read internal/broker/broker_e2e_test.go [offset=270, limit=120]
Let me check the e2e latest-cert test coverage and whether offlineACL has any tests.

$ rg -n "offlineACL|newOfflineACL|offline\." internal/broker/*_test.go; echo "=== retained QoS1 drop? ==="; rg -n "QosDropped|OnQosDropped|OnRetainedExpired" internal/broker/*_test.go
=== retained QoS1 drop? ===
internal/broker/sessionstore_test.go:166:	ss.OnQosDropped(cl, packets.Packet{PacketID: 7})
internal/broker/sessionstore_test.go:190:	ss.OnRetainedExpired("r/keep")
internal/broker/sessionstore_test.go:192:		t.Fatalf("after OnRetainedExpired: %+v, want empty", msgs)

I have enough for a focused second-pass broker review. The area worth examining was the new coldstart ACL fallback (landed after the 09-04 review). Let me verify one load-count claim before writing the review.
$ rg -n "func (f \*fakeStore)|GetDevice|GetInterface|getCalls|callCount|calls " internal/broker/broker_test.go | head -30
131:func (f *fakeStore) GetDevice(_ context.Context, realmID int16, id deviceid.ID) (*store.Device, error) {
142:func (f *fakeStore) GetInterface(_ context.Context, realmID int16, name string, major int) (*store.StoredInterface, error) {

→ Read internal/broker/broker_test.go [limit=140, offset=100]
The post-09-04 commit `d201db4` added a new ACL-miss fallback, `syncOwnershipOf`, which the first review never saw. Let me confirm the flood cost precisely and check the file tail before writing up.

$ wc -l .mule/todo.md; tail -6 .mule/todo.md
205 .mule/todo.md
- [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a locally-run `elixir -e` with `Base.url_decode64!`) with the device ID `"AAAAAAAAAAAAAAAAAAAAAB"` — nonzero unused bits in the 22nd base64url char — and report whether upstream accepts it: `pkg/deviceid/deviceid.go:34-37` uses `base64.RawURLEncoding.Strict()` and deviceid_test.go:176 asserts rejection while the comment claims parity with `Elixir Base.url_decode64!(padding: false)`, which decodes and discards those bits rather than erroring. Astrate may therefore be stricter than upstream on the same wire form (a device registered with such an id on upstream would be rejected here). If upstream rejects it identically, the strict encoding is verified parity and the task is done; if upstream accepts it, escalate the relax-vs-strict call to `.mule/for-giulio.md` with the measurement. Probe first, no code change either way.
- [x] docs-sync-native-metrics-example-fake-series [auto]: the `/astrate/v1/metrics` example in docs/api/astrate_native_api.yaml (yaml:93-95) shows `astrate_devices_total{realm="test"} 42`, but no such series is ever registered — the registry carries `astrate_broker_sessions` (internal/observability/metrics.go:54-62), `astrate_db_pool_{acquired,idle,total,max}_conns` (metrics.go:65-79), plus the engine/flow/trigger families (`astrate_engine_*` engine/router.go:411-444, `astrate_engine_stream_dropped_total` engine/stream/bus.go:119, `astrate_engine_trigger_*` engine/triggers/actions.go:296-300, `astrate_flow_router_*` flow/router.go:260-276) and the standard `go_*`/`process_*` runtime collectors — so the example advertises a scrape result the wire can never return. Swap it for a real series (e.g. `astrate_broker_sessions 0`), keeping the HELP/TYPE comment style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] store-devices-inhibit-re-register [legion] [auto]: `RegisterDevice` (internal/store/devices.go:75-91) silently clears the inhibit flag — its `ON CONFLICT ... DO UPDATE SET status = 'registered'` fires for any device with `first_credentials_request IS NULL`, including one the admin inhibited via `SetDeviceInhibited` (devices.go:251-268), which sets `status='inhibited'` on unconfirmed devices too; §5.3 says an inhibited device blocks new credentials and connections, so the re-registration re-opens it. Preserve `'inhibited'` in the SET (e.g. `status = CASE WHEN devices.status = 'inhibited' THEN 'inhibited' ELSE 'registered' END`) and add a Lifecycle case in internal/store/devices_test.go: inhibit an unconfirmed device, re-register, assert status stays inhibited and the secret still rotates. Verify the assert against upstream's register-not-touching-inhibit on the Legion while the integration suite runs.
- [x] store-pipelines-empty-name-zero-blocks-test [auto]: extend internal/store/pipelines_validate_test.go with the two untested error branches of `validatePipelineGraph` (internal/store/pipelines.go:51-64) — a definition with zero blocks and one whose block has an empty name must both be rejected. Pure helper, no DB; the existing suite covers only acyclic/cyclic/unknown-ref/duplicate-name.
- [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (internal/store/devices.go:130-131, "if several devices share an alias the lowest device ID wins", `ORDER BY id LIMIT 1`) — no test drives it today; add a case in internal/store/devices_test.go that gives two devices the same alias and asserts the lookup resolves the lower ID. Needs the integration DB.

$ git show d201db4 --stat | head -20
commit d201db478ba0d9475b9615f4a8d5c531d0d7ce75
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Fri Sep 4 21:50:18 2026 +0200

    mule: broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [approved 2026-09-04]

 internal/broker/aclhook.go     | 23 ++++++++-----
 internal/broker/authhook.go    | 44 +++++++++++++++++++++++--
 internal/broker/broker_test.go | 74 ++++++++++++++++++++++++++++++++++++++++++
 3 files changed, 131 insertions(+), 10 deletions(-)

I have my findings. The 09-04 review's fix (`d201db4`) added a synchronous per-name store read to the ACL miss path that the debounce was designed to prevent, and I found an unbounded offlineACL entries map plus a nil-guard with no direct test. Let me verify the Ownership enum and `ownershipOf` connect-order claim once more, then write the review.
→ Read pkg/interfaceschema/types.go [offset=58, limit=15]
I have a clear picture. The key new seam since 09-04 is the cold-start synchronization fallback landed in `d201db4`. Let me write the review and append the proposed lines.

← Write .mule/reviews/broker-2026-09-21.md
Wrote file successfully.

← Edit .mule/reviews/broker-2026-09-21.md
Index: /root/astrate-mule/.mule/reviews/broker-2026-09-21.md
===================================================================
--- /root/astrate-mule/.mule/reviews/broker-2026-09-21.md
+++ /root/astrate-mule/.mule/reviews/broker-2026-09-21.md
@@ -47,10 +47,10 @@
   reopens every second via the debounced full reload, so this is a standing
   per-session throughput of ~N reads/sec, attacker-controlled (QoS0 rate is
   unbounded apart from MaxPacketBytes).
 
-Before the fix (per the review's description of the miss path) a denied uint
-kept that cost at ~1/sec; the fix traded a specific cold-start race for
+Before the fix a denied miss cost nothing (the debounce-skipped reload re-checked
+a cold map and denied); the fix traded a specific cold-start race for
 re-opening the exact anti-flood property the debounce comment (authhook.go:
 44-45) exists to guarantee. The fix's own comment ("the per-interface store
 read happens once per introspection refresh cycle") is only true for names
 already seen.

Now append the two task lines to `.mule/todo.md`.
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -202,4 +202,6 @@
 - [x] docs-sync-native-socket-missing-403-500 [auto]: `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml documents only `101`/`200`/`401` (yaml:318-324), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46; internal/auth/middleware.go:57-92) also answers `403` — a valid a_ch JWT that does not grant the socket path (`WriteForbidden`, middleware.go:90) — and `500` when `GetRealmByName` fails (`WriteInternalServerError`, middleware.go:74), exactly the pair the Phoenix endpoint already documents (yaml:360-366). Add a `403` Forbidden response (new response component, same `{"errors":{"detail":...}}` envelope shape as the Unauthorized one but `detail: Forbidden`, from `pkg/astarteapi` DetailForbidden) and `500` (reuse the existing InternalServerError ref) to the socket's responses. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] store-devices-inhibit-re-register [legion] [auto]: `RegisterDevice` (internal/store/devices.go:75-91) silently clears the inhibit flag — its `ON CONFLICT ... DO UPDATE SET status = 'registered'` fires for any device with `first_credentials_request IS NULL`, including one the admin inhibited via `SetDeviceInhibited` (devices.go:251-268), which sets `status='inhibited'` on unconfirmed devices too; §5.3 says an inhibited device blocks new credentials and connections, so the re-registration re-opens it. Preserve `'inhibited'` in the SET (e.g. `status = CASE WHEN devices.status = 'inhibited' THEN 'inhibited' ELSE 'registered' END`) and add a Lifecycle case in internal/store/devices_test.go: inhibit an unconfirmed device, re-register, assert status stays inhibited and the secret still rotates. Verify the assert against upstream's register-not-touching-inhibit on the Legion while the integration suite runs.
 - [x] store-pipelines-empty-name-zero-blocks-test [auto]: extend internal/store/pipelines_validate_test.go with the two untested error branches of `validatePipelineGraph` (internal/store/pipelines.go:51-64) — a definition with zero blocks and one whose block has an empty name must both be rejected. Pure helper, no DB; the existing suite covers only acyclic/cyclic/unknown-ref/duplicate-name.
 - [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (internal/store/devices.go:130-131, "if several devices share an alias the lowest device ID wins", `ORDER BY id LIMIT 1`) — no test drives it today; add a case in internal/store/devices_test.go that gives two devices the same alias and asserts the lookup resolves the lower ID. Needs the integration DB.
+- [ ] broker-acl-coldstart-fallback-flood [auto]: `syncOwnershipOf` (internal/broker/authhook.go:212-237, reached from the ACL miss path aclhook.go:196-206) does a full synchronous `GetDevice`+`GetInterface` store read for every distinct interface name a device publishes to, with no rate limit — `introspectionReloadDebounce` (authhook.go:42-45, "so an adversarial topic flood cannot hammer the database") throttles only `refreshIfStale`'s full reload, and the fallback's per-name cache write means N distinct bogus names in one window cause N synchronous full-device reads per second, an open anti-flood hole the pre-fix denied-miss path did not have. Gate the fallback to one sync read per session per debounce window (or fold it into a full reload that stamps `lastIntroLoad`), keep `TestBrokerACLColdStartIntrospectionMiss` (one name) green, and add a T1 test in broker_test.go with a GetDevice-counting fake: connect cold, publish to K>=2 distinct unknown-interface names, assert no more than one fallback read was made.
+- [ ] broker-offlineacl-entry-eviction [auto]: `offlineACL.entries` (internal/broker/aclhook.go:100-162) is append-only — `ownershipOf` upserts one entry per CN ever ACL-checked offline and nothing ever deletes, so a long-lived broker grows the map unboundedly across device churn even after the 10s `offlineACLCacheTTL` passes (the entry is kept and re-stamped, not reaped). Evict stale entries lazily on access (drop an entry whose `loadedAt` predates some multiple of the TTL before inserting another, or cap+LRU), add an injectable `now func() time.Time` in the same style as `lifecycleHook.now`, and cover in a container-free unit test: seed N entries, advance the clock, access one, assert the map stays bounded.

$ tail -4 .mule/todo.md | cut -c1-80
- [x] store-pipelines-empty-name-zero-blocks-test [auto]: extend internal/store/
- [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie
- [ ] broker-acl-coldstart-fallback-flood [auto]: `syncOwnershipOf` (internal/br
- [ ] broker-offlineacl-entry-eviction [auto]: `offlineACL.entries` (internal/br

Done. This was a proposal-only run: wrote the review and appended task lines, no code, no git.

**Findings** (second-pass on `broker`, the most-stale area — 09-04 reviews covered the package, but two commits landed since, notably `d201db4`):

1. **Cold-start ACL fallback re-opens the anti-flood hole** — `syncOwnershipOf` (authhook.go:212-237) does a synchronous full `GetDevice`+`GetInterface` per distinct interface name on the ACL miss path (aclhook.go:196-206), with no rate limit; `introspectionReloadDebounce` only throttles `refreshIfStale`'s full reload, so N distinct bogus names in a window cost N store reads/sec. The pre-fix denied-miss path never touched the store.
2. **`offlineACL.entries` is append-only** (aclhook.go:100-162) — entries are upserted per CN and never evicted, so the map grows unboundedly over a long process lifetime even past the 10s TTL.

Also re-checked and confirmed unchanged: `intakeHook` deferred-ack, `Publisher` expiry, `lifecycleHook` takeover guards, session store Seq. The zombie-session / offline-ACL-tests / doc-comment items are already queued and were not duplicated.

- Files: `.mule/reviews/broker-2026-09-21.md`, `.mule/todo.md`
- Verified: read-side verification only; no code or tests written, nothing to run
- Unsure: nothing
- Follow-ups: appended `broker-acl-coldstart-fallback-flood` and `broker-offlineacl-entry-eviction`, both `[auto]`, container-free-testable
