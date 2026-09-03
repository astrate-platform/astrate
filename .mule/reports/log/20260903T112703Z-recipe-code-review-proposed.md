slug: recipe-code-review
verdict: proposed
at:  2d82565
ran: 2026-09-03T11:27:03Z on DietPi in 230s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
Let me check recent reviews and pick a package to review.

$ ls /root/astrate-mule/.mule/reviews/ 2>/dev/null; echo "---"; ls /root/astrate-mule/.mule/log.md 2>/dev/null
flow-2026-09-02.md
realm-2026-09-01.md
store-2026-07-27.md
---
/root/astrate-mule/.mule/log.md

$ ls /root/astrate-mule/internal/ /root/astrate-mule/pkg/ 2>/dev/null
/root/astrate-mule/internal/:
appengine
auth
broker
config
engine
flow
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil

/root/astrate-mule/pkg/:
astarteapi
deviceid
interfaceschema
payload

$ rg -c '' --glob 'internal/**/*.go' | sort -t: -k2 -rn | head -20
internal/appengine/channels/ws_test.go:916
internal/engine/triggers/actions_test.go:805
internal/store/datastreams_test.go:783
internal/engine/fakestore_test.go:769
internal/pairing/service_test.go:632
internal/broker/broker_e2e_test.go:618
internal/pairing/http_test.go:614
internal/engine/triggers/match.go:595
internal/broker/sessionstore.go:533
internal/appengine/http_test.go:517
internal/engine/triggers/actions.go:514
internal/engine/cache.go:496
internal/realm/service.go:491
internal/engine/router.go:490
internal/broker/broker_test.go:475
internal/store/devices.go:462
internal/engine/data_test.go:461
internal/store/devices_test.go:446
internal/engine/batch_test.go:446
internal/broker/authhook.go:440

`engine` is a big package and hasn't been reviewed. Let me get its file sizes and start reading.
$ rg -c '' --glob 'internal/engine/**/*.go' | sort -t: -k2 -rn
internal/engine/triggers/actions_test.go:805
internal/engine/fakestore_test.go:769
internal/engine/triggers/match.go:595
internal/engine/triggers/actions.go:514
internal/engine/cache.go:496
internal/engine/router.go:490
internal/engine/data_test.go:461
internal/engine/batch_test.go:446
internal/engine/engine_e2e_test.go:439
internal/engine/triggers/match_test.go:436
internal/engine/engine_integration_test.go:434
internal/engine/batch.go:406
internal/engine/data.go:386
internal/engine/router_test.go:362
internal/engine/cache_test.go:362
internal/engine/triggers/policy_test.go:332
internal/engine/control_test.go:329
internal/engine/engine.go:318
internal/engine/control.go:315
internal/engine/serverdata_test.go:312
internal/engine/forward/http_test.go:310
internal/engine/serverdata.go:280
internal/engine/capabilities_test.go:255
internal/engine/triggers/policy.go:252
internal/engine/stream/bus_test.go:234
internal/engine/stream/bus.go:203
internal/engine/forward/nats_test.go:196
internal/engine/introspection_test.go:189
internal/engine/triggers/events.go:188
internal/engine/introspection.go:168
internal/engine/triggers/events_test.go:142
internal/engine/bus_events_test.go:140
internal/engine/triggers/transient_test.go:139
internal/engine/forward/http.go:125
internal/engine/topics_test.go:109
internal/engine/capabilities.go:101
internal/engine/triggers/errorname_test.go:98
internal/engine/topics.go:97
internal/engine/triggers/errorname.go:76
internal/engine/forward/nats.go:53
internal/engine/forward/envelope.go:29
internal/engine/triggers/transient.go:22

→ Read .mule/todo.md
I'll review `internal/engine` — it's the largest unreviewed package. Let me read the biggest source files.
→ Read internal/engine/cache.go
→ Read internal/engine/router.go
→ Read internal/engine/batch.go
→ Read internal/engine/engine.go
→ Read internal/engine/data.go
→ Read internal/engine/control.go
→ Read internal/engine/serverdata.go
→ Read internal/engine/introspection.go
→ Read internal/engine/capabilities.go
→ Read internal/engine/topics.go
→ Read internal/engine/forward/http.go

$ ls /root/astrate-mule/internal/engine/*_test.go /root/astrate-mule/internal/engine/**/*_test.go 2>/dev/null
/root/astrate-mule/internal/engine/batch_test.go
/root/astrate-mule/internal/engine/bus_events_test.go
/root/astrate-mule/internal/engine/cache_test.go
/root/astrate-mule/internal/engine/capabilities_test.go
/root/astrate-mule/internal/engine/control_test.go
/root/astrate-mule/internal/engine/data_test.go
/root/astrate-mule/internal/engine/engine_e2e_test.go
/root/astrate-mule/internal/engine/engine_integration_test.go
/root/astrate-mule/internal/engine/fakestore_test.go
/root/astrate-mule/internal/engine/forward/http_test.go
/root/astrate-mule/internal/engine/forward/nats_test.go
/root/astrate-mule/internal/engine/introspection_test.go
/root/astrate-mule/internal/engine/router_test.go
/root/astrate-mule/internal/engine/serverdata_test.go
/root/astrate-mule/internal/engine/stream/bus_test.go
/root/astrate-mule/internal/engine/topics_test.go
/root/astrate-mule/internal/engine/triggers/actions_test.go
/root/astrate-mule/internal/engine/triggers/errorname_test.go
/root/astrate-mule/internal/engine/triggers/events_test.go
/root/astrate-mule/internal/engine/triggers/match_test.go
/root/astrate-mule/internal/engine/triggers/policy_test.go
/root/astrate-mule/internal/engine/triggers/transient_test.go

Let me check what tests exist and what areas might be under-tested:
$ rg -n 'func Test' /root/astrate-mule/internal/engine/ --glob '*_test.go' | sort
/root/astrate-mule/internal/engine/batch_test.go:180:func TestObjectAndPropertyRows(t *testing.T) {
/root/astrate-mule/internal/engine/batch_test.go:223:func TestFlushAckAfterCommit(t *testing.T) {
/root/astrate-mule/internal/engine/batch_test.go:268:func TestFlushRetryTransient(t *testing.T) {
/root/astrate-mule/internal/engine/batch_test.go:292:func TestFlushPoisonedBatch(t *testing.T) {
/root/astrate-mule/internal/engine/batch_test.go:340:func TestFlushAbandonOnQuit(t *testing.T) {
/root/astrate-mule/internal/engine/batch_test.go:34:func TestEncodeValueJSON(t *testing.T) {
/root/astrate-mule/internal/engine/batch_test.go:374:func TestBatchTriggers(t *testing.T) {
/root/astrate-mule/internal/engine/batch_test.go:412:func TestAfterCommitHook(t *testing.T) {
/root/astrate-mule/internal/engine/batch_test.go:81:func TestIndividualRowColumns(t *testing.T) {
/root/astrate-mule/internal/engine/bus_events_test.go:55:func TestBusEventDetail(t *testing.T) {
/root/astrate-mule/internal/engine/cache_test.go:108:func TestSchemaCacheSelfHeal(t *testing.T) {
/root/astrate-mule/internal/engine/cache_test.go:136:func TestSchemaCacheBrokenDefinitions(t *testing.T) {
/root/astrate-mule/internal/engine/cache_test.go:15:func TestSchemaCacheLoadAndLookup(t *testing.T) {
/root/astrate-mule/internal/engine/cache_test.go:166:func TestSchemaCacheSnapshotRace(t *testing.T) {
/root/astrate-mule/internal/engine/cache_test.go:214:func TestDeviceCache(t *testing.T) {
/root/astrate-mule/internal/engine/cache_test.go:270:func TestSchemaCacheTriggerPolicyAttachment(t *testing.T) {
/root/astrate-mule/internal/engine/cache_test.go:339:func TestDeviceCacheLoadError(t *testing.T) {
/root/astrate-mule/internal/engine/cache_test.go:56:func TestSchemaCacheReloadRealm(t *testing.T) {
/root/astrate-mule/internal/engine/cache_test.go:86:func TestSchemaCacheReloadUnknownRealm(t *testing.T) {
/root/astrate-mule/internal/engine/capabilities_test.go:119:func TestHandleCapabilitiesRejects(t *testing.T) {
/root/astrate-mule/internal/engine/capabilities_test.go:12:func TestDecodeCapabilities(t *testing.T) {
/root/astrate-mule/internal/engine/capabilities_test.go:155:func TestPlaintextConsumerProperties(t *testing.T) {
/root/astrate-mule/internal/engine/capabilities_test.go:199:func TestZlibConsumerPropertiesDefault(t *testing.T) {
/root/astrate-mule/internal/engine/capabilities_test.go:233:func TestPurgeCompressionFor(t *testing.T) {
/root/astrate-mule/internal/engine/capabilities_test.go:68:func TestValidateCapability(t *testing.T) {
/root/astrate-mule/internal/engine/capabilities_test.go:90:func TestHandleCapabilities(t *testing.T) {
/root/astrate-mule/internal/engine/control_test.go:134:func TestEmptyCache(t *testing.T) {
/root/astrate-mule/internal/engine/control_test.go:219:func TestEmptyCacheJSONHint(t *testing.T) {
/root/astrate-mule/internal/engine/control_test.go:261:func TestProducerPropertiesPurge(t *testing.T) {
/root/astrate-mule/internal/engine/control_test.go:309:func TestControlRejects(t *testing.T) {
/root/astrate-mule/internal/engine/control_test.go:43:func TestControlFrameGolden(t *testing.T) {
/root/astrate-mule/internal/engine/control_test.go:89:func TestInflatePropertiesRejects(t *testing.T) {
/root/astrate-mule/internal/engine/data_test.go:224:func TestRejectSeam(t *testing.T) {
/root/astrate-mule/internal/engine/data_test.go:239:func TestAcceptedOps(t *testing.T) {
/root/astrate-mule/internal/engine/data_test.go:329:func TestFormatHintFlip(t *testing.T) {
/root/astrate-mule/internal/engine/data_test.go:391:func TestUnhandledDispatch(t *testing.T) {
/root/astrate-mule/internal/engine/data_test.go:443:func TestDeviceLoadParking(t *testing.T) {
/root/astrate-mule/internal/engine/data_test.go:46:func TestRejectReasons(t *testing.T) {
/root/astrate-mule/internal/engine/engine_e2e_test.go:251:func TestEngineE2E(t *testing.T) {
/root/astrate-mule/internal/engine/engine_integration_test.go:88:func TestEngine(t *testing.T) {
/root/astrate-mule/internal/engine/forward/http_test.go:120:func TestMethodHonoured(t *testing.T) {
/root/astrate-mule/internal/engine/forward/http_test.go:138:func TestNilActionAndEvent(t *testing.T) {
/root/astrate-mule/internal/engine/forward/http_test.go:164:func TestStatusTable(t *testing.T) {
/root/astrate-mule/internal/engine/forward/http_test.go:212:func TestEmptyNonNilActionAndEvent(t *testing.T) {
/root/astrate-mule/internal/engine/forward/http_test.go:238:func TestTransportError(t *testing.T) {
/root/astrate-mule/internal/engine/forward/http_test.go:253:func TestCancelledContext(t *testing.T) {
/root/astrate-mule/internal/engine/forward/http_test.go:271:func TestNewRejectsBadConfig(t *testing.T) {
/root/astrate-mule/internal/engine/forward/http_test.go:31:func TestHappyPath(t *testing.T) {
/root/astrate-mule/internal/engine/forward/http_test.go:86:func TestStaticHeaders(t *testing.T) {
/root/astrate-mule/internal/engine/forward/nats_test.go:131:func TestNATSForwardEmptyNonNilActionAndEvent(t *testing.T) {
/root/astrate-mule/internal/engine/forward/nats_test.go:174:func TestNewNATSRejectsBadConfig(t *testing.T) {
/root/astrate-mule/internal/engine/forward/nats_test.go:37:func TestNATSForwardHappyPath(t *testing.T) {
/root/astrate-mule/internal/engine/forward/nats_test.go:88:func TestNATSForwardNilActionAndEvent(t *testing.T) {
/root/astrate-mule/internal/engine/introspection_test.go:133:func TestIntrospectionRejects(t *testing.T) {
/root/astrate-mule/internal/engine/introspection_test.go:164:func TestIntrospectionParking(t *testing.T) {
/root/astrate-mule/internal/engine/introspection_test.go:17:func TestParseIntrospection(t *testing.T) {
/root/astrate-mule/internal/engine/introspection_test.go:84:func TestIntrospectionHandler(t *testing.T) {
/root/astrate-mule/internal/engine/router_test.go:161:func TestBackpressureQoS0Drops(t *testing.T) {
/root/astrate-mule/internal/engine/router_test.go:206:func TestDrain(t *testing.T) {
/root/astrate-mule/internal/engine/router_test.go:21:func TestOrderingSingleDevice(t *testing.T) {
/root/astrate-mule/internal/engine/router_test.go:251:func TestShardPanicRecovery(t *testing.T) {
/root/astrate-mule/internal/engine/router_test.go:283:func TestInvalidation(t *testing.T) {
/root/astrate-mule/internal/engine/router_test.go:318:func TestLifecycleEviction(t *testing.T) {
/root/astrate-mule/internal/engine/router_test.go:348:func TestShardOf(t *testing.T) {
/root/astrate-mule/internal/engine/router_test.go:61:func TestCrossDeviceParallelism(t *testing.T) {
/root/astrate-mule/internal/engine/router_test.go:97:func TestBackpressureQoS1Blocks(t *testing.T) {
/root/astrate-mule/internal/engine/serverdata_test.go:108:func TestPublishServerDatastream(t *testing.T) {
/root/astrate-mule/internal/engine/serverdata_test.go:153:func TestPublishServerValueErrors(t *testing.T) {
/root/astrate-mule/internal/engine/serverdata_test.go:212:func TestUnsetServerProperty(t *testing.T) {
/root/astrate-mule/internal/engine/serverdata_test.go:269:func TestServerObjectAggregate(t *testing.T) {
/root/astrate-mule/internal/engine/serverdata_test.go:40:func TestPublishServerProperty(t *testing.T) {
/root/astrate-mule/internal/engine/serverdata_test.go:80:func TestPublishServerPropertyJSONHint(t *testing.T) {
/root/astrate-mule/internal/engine/stream/bus_test.go:110:func TestSlowConsumerDrops(t *testing.T) {
/root/astrate-mule/internal/engine/stream/bus_test.go:133:func TestCancelUnsubscribes(t *testing.T) {
/root/astrate-mule/internal/engine/stream/bus_test.go:157:func TestCloseClosesAllSubscribers(t *testing.T) {
/root/astrate-mule/internal/engine/stream/bus_test.go:176:func TestSubscribeAfterClose(t *testing.T) {
/root/astrate-mule/internal/engine/stream/bus_test.go:192:func TestConcurrentPublishSubscribe(t *testing.T) {
/root/astrate-mule/internal/engine/stream/bus_test.go:42:func TestPublishRoutesByRealm(t *testing.T) {
/root/astrate-mule/internal/engine/stream/bus_test.go:60:func TestFilter(t *testing.T) {
/root/astrate-mule/internal/engine/topics_test.go:14:func TestSplitDeviceTopic(t *testing.T) {
/root/astrate-mule/internal/engine/topics_test.go:45:func TestClassify(t *testing.T) {
/root/astrate-mule/internal/engine/topics_test.go:70:func TestMatchInterface(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:153:func TestWebhookDelivered(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:204:func TestWebhookRetriesThenSucceeds(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:238:func TestWebhook4xxIsPermanent(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:262:func TestWebhookFailsAfterMaxAttempts(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:308:func TestForwarderSkippedWhenUnset(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:319:func TestForwarderForwards(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:339:func TestForwarderError(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:348:func TestEnqueueAfterCloseDrops(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:364:func TestQueueFullDrops(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:393:func TestCloseDrainsQueue(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:421:func TestIgnoreSSLErrors(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:453:func TestPolicyRetriesClientError(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:493:func TestPolicyDiscardsServerError(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:524:func TestPolicyExplicitStatusRetries503Discards500(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:570:func TestPolicyTransportFailureFollowsServerError(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:597:func TestNoPolicyFollowsDefaultBehaviour(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:624:func TestPolicyTTExpiresDelivery(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:679:func TestPolicyCapacityDrops(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:742:func TestPolicyCapacityResetsAfterDelivery(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:774:func TestEnqueueRacingCloseDoesNotPanic(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/actions_test.go:80:func TestParseAction(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/errorname_test.go:42:func TestUpstreamErrorNamePassthrough(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/errorname_test.go:53:func TestUpstreamErrorNameClosedSetInvariant(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/errorname_test.go:73:func TestNewDeviceErrorEventWiring(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/errorname_test.go:7:func TestUpstreamErrorNameMapping(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/events_test.go:124:func TestEventTimestampPrecision(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/events_test.go:29:func TestEventGoldens(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/events_test.go:63:func TestEventUpstreamShape(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/match_test.go:124:func TestMatchData(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/match_test.go:180:func TestMatchDataOperators(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/match_test.go:230:func TestMatchDevice(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/match_test.go:269:func TestCompileAcceptsUnsupported(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/match_test.go:293:func TestCompileRejects(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/match_test.go:373:func TestCompileActions(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/match_test.go:411:func TestCompilePolicyField(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/policy_test.go:152:func TestDecide(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/policy_test.go:261:func TestDecideNilPolicy(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/policy_test.go:275:func TestDecideTransportReasonDistinct(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/policy_test.go:294:func TestDecideEveryDecisionHasReason(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/policy_test.go:7:func TestCompilePolicy(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/transient_test.go:23:func TestCompileConditionDataTrigger(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/transient_test.go:57:func TestCompileConditionActionIsNil(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/transient_test.go:68:func TestCompileConditionWildcardData(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/transient_test.go:93:func TestCompileConditionInvalidInputs(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/transient_test.go:9:func TestCompileConditionDeviceTrigger(t *testing.T) {

Let me check a couple of specific test files to verify my findings:
→ Read internal/engine/serverdata_test.go
→ Read internal/engine/control_test.go
→ Read internal/engine/data_test.go
→ Read internal/engine/fakestore_test.go
Now I have a thorough understanding. Let me write the review file and propose tasks.

← Write .mule/reviews/engine-2026-09-03.md
Wrote file successfully.

Now let me append the proposed task lines to todo.md:
→ Read .mule/todo.md [offset=120]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -119,4 +119,7 @@
 - [!] lint-engine-goimports: run goimports -w on `internal/engine/capabilities.go` — its import block is not goimports-formatted (the broker import is not grouped under the local local-prefix), which the golangci goimports formatter flags at every run. [auto] — BLOCKED: gates failed
 - [!] lint-flow-findings: clean up the unexcluded revive/gosec findings in `internal/flow` — add `// #nosec G115` to the false-positive rune→uint64 conversion at router.go:215 (a rune can never overflow uint64, mirroring the existing nosec on line 218), delete the empty `if` block at flow_test.go:118, rename the unused `ctx` param of `runLane` to `_` (router.go:180), and add the missing doc comment for the `TypeInteger` constant (message.go:21). [auto] — BLOCKED: gates failed
 - [!] lint-pairing-test-param: in `internal/pairing/service_test.go:192`, rename the unused `at` parameter of the `OnRegistered` func to `_` to silence revive's unused-parameter. [auto] — BLOCKED: gates failed
 - [!] lint-flow-stutter: in `internal/flow`, rename exported `FlowStatus`/`FlowConfig`/`FlowMessage` to `Status`/`Config`/`Message` to satisfy revive's stutter rule; touches flow.go, message.go, block.go, graph.go, router.go and their tests — a public API rename, so no wire/JSON change (the Type names are unchanged), it is purely Go identifier renaming. [auto] — BLOCKED: gates failed
+- [ ] serverdata-aftercommit-test: in `internal/engine/serverdata_test.go`, add a test that wires an `afterCommit` observer on a `newWiredRig` engine, calls `PublishServerValue` for a server-owned property, and asserts whether the observer fires. The batch path fires `afterCommit` for every committed op (batch.go:168), but `PublishServerValue` bypasses the batch entirely (serverdata.go:99-162) — the current tests never check this. The test documents the intended behaviour (exclusion or parity with the data path) and fails if the contract changes. [auto]
+- [ ] consumer-properties-plaintext-fallback: in `internal/engine/control.go` `sendConsumerProperties` (line 204), when `devices.peek` returns nil (device not cached after eviction), the code defaults to zlib compression — but a device that declared `purge_properties_compression_format: plaintext` via the capabilities topic cannot decompress zlib. Fix by falling back to a store read of the device's capability state (or persisting the capability), and add a test in `internal/engine/capabilities_test.go` that evicts the cache entry for a plaintext device and asserts the purge payload is uncompressed. [auto]
+- [ ] object-path-ok-comment: in `internal/engine/data.go` `objectPathOK` (line 330), the loop over `ci.ObjectLeaves` returns after checking the first entry. Add a comment explaining why one leaf suffices (upstream validation guarantees uniform prefix depth), so future readers do not mistake the early return for a bug. [auto]

Done. Here's what I produced:

**Review file:** `.mule/reviews/engine-2026-09-03.md`

**3 task lines appended to `.mule/todo.md`:**

1. **`serverdata-aftercommit-test`** — `PublishServerValue` bypasses the batch and never calls `afterCommit`, so server-initiated property/datastream writes produce no trigger events and no live bus events. No test checks this. The task adds one to document the intended exclusion or prove it's a bug.

2. **`consumer-properties-plaintext-fallback`** — `sendConsumerProperties` defaults to zlib when `devices.peek` returns nil (uncached device). A plaintext-only device that just connected cannot decompress the purge message. Needs a fix and a regression test.

3. **`object-path-ok-comment`** — `objectPathOK` iterates `ObjectLeaves` but returns after the first entry. The early return is correct (upstream validation guarantees uniform prefix depth), but the function name implies full validation and the code looks like a bug to future readers. A clarifying comment suffices.

**What I decided not to propose:** renames/reformatting (code is idiomatic), the `noteFormat` lock discipline (verified correct), the `deviceCache.get` TOCTOU (impossible by per-device shard invariant), the retry/backoff loops (textbook), and anything in `triggers/` or `stream/` (out of scope).
