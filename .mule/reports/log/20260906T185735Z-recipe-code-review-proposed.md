slug: recipe-code-review
verdict: proposed
at:  0f86ad6
ran: 2026-09-06T18:57:35Z on DietPi in 199s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
→ Read .mule/log.md

$ ls /root/astrate-mule/.mule/reviews/ 2>/dev/null || echo "no reviews dir"
broker-2026-09-04.md
flow-2026-09-05.md

$ ls /root/astrate-mule/internal/ /root/astrate-mule/pkg/ 2>/dev/null
/root/astrate-mule/internal/:
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

/root/astrate-mule/pkg/:
astarteapi
deviceid
interfaceschema
payload

Recent reviews covered `broker` and `flow`. Let me pick a fresh package and check its size.
$ for pkg in /root/astrate-mule/internal/*/; do name=$(basename "$pkg"); count=$(find "$pkg" -name '*.go' | wc -l); echo "$count $name"; done | sort -rn
49 flow
48 engine
32 store
23 appengine
13 broker
9 realm
8 pairing
8 flowapi
8 auth
6 testutil
5 observability
5 housekeeping
4 httpx
2 swagger
2 config

$ find /root/astrate-mule/internal/engine -name '*.go' | xargs wc -l | sort -rn | head -30
14892 total
  1074 /root/astrate-mule/internal/engine/triggers/actions_test.go
   818 /root/astrate-mule/internal/engine/fakestore_test.go
   769 /root/astrate-mule/internal/engine/triggers/match.go
   610 /root/astrate-mule/internal/engine/triggers/actions.go
   594 /root/astrate-mule/internal/engine/triggers/match_test.go
   509 /root/astrate-mule/internal/engine/data.go
   502 /root/astrate-mule/internal/engine/cache.go
   490 /root/astrate-mule/internal/engine/router.go
   461 /root/astrate-mule/internal/engine/data_test.go
   459 /root/astrate-mule/internal/engine/batch.go
   456 /root/astrate-mule/internal/engine/engine.go
   446 /root/astrate-mule/internal/engine/batch_test.go
   440 /root/astrate-mule/internal/engine/engine_e2e_test.go
   435 /root/astrate-mule/internal/engine/engine_integration_test.go
   410 /root/astrate-mule/internal/engine/triggers/policy_test.go
   363 /root/astrate-mule/internal/engine/cache_test.go
   362 /root/astrate-mule/internal/engine/router_test.go
   330 /root/astrate-mule/internal/engine/control_test.go
   321 /root/astrate-mule/internal/engine/serverdata.go
   316 /root/astrate-mule/internal/engine/control.go
   312 /root/astrate-mule/internal/engine/triggers/events.go
   312 /root/astrate-mule/internal/engine/serverdata_test.go
   310 /root/astrate-mule/internal/engine/forward/http_test.go
   297 /root/astrate-mule/internal/engine/triggers/policy.go
   293 /root/astrate-mule/internal/engine/triggers_change_test.go
   255 /root/astrate-mule/internal/engine/capabilities_test.go
   234 /root/astrate-mule/internal/engine/stream/bus_test.go
   209 /root/astrate-mule/internal/engine/stream/bus.go
   196 /root/astrate-mule/internal/engine/forward/nats_test.go

`engine` is 48 files, unrecent, and large. I'll go deep there. Let me survey what's in it.
$ ls /root/astrate-mule/internal/engine/
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

→ Read internal/engine/triggers/match.go
→ Read internal/engine/triggers/actions.go
→ Read internal/engine/triggers/events.go
→ Read internal/engine/triggers/policy.go
Now let me read the core engine files.
→ Read internal/engine/engine.go
→ Read internal/engine/data.go
→ Read internal/engine/router.go
→ Read internal/engine/batch.go
Let me read the remaining core files — batch.go, cache.go, control.go, serverdata.go, capabilities.go, topics.go, introspection.go, and the stream/bus.
→ Read internal/engine/cache.go
→ Read internal/engine/serverdata.go
→ Read internal/engine/control.go
→ Read internal/engine/capabilities.go
→ Read internal/engine/topics.go
→ Read internal/engine/introspection.go
Let me read the remaining files: control.go, capabilities.go, topics.go, introspection.go, and the stream/bus.
→ Read internal/engine/stream/bus.go
✱ Grep "PublishServerValue|fireCommitted|bus|trigger" in internal/engine · 14 matches
I've now read the core engine. Let me analyze for real issues worth proposing, focusing on the recipe's categories. Let me examine a few specific concerns more closely.

**Concern 1: `PublishServerValue`/`publishAsOwner` do not fire triggers or fan out to the live bus.** Looking at `publishAsOwner` (serverdata.go:94-203): it persists the value and (for `PublishServerValue`) publishes to the broker, but never runs `trackPrevious`/`fireData` to fire trigger events or publish to the live bus. Compare with the inbound data path (`processData` → `fireCommitted`). This is an asymmetry. But wait — is this a bug or a design decision? Let me check the tests to understand expected behavior, and whether there's a documented reason.

Let me check the tests for the engine and serverdata.
✱ Grep "PublishServerValue|Trigger|Bus\(\)|Subscribe|fireCommitted" in internal/engine · 100 matches
→ Read internal/engine/serverdata_test.go [offset=1, limit=150]
The server-data publish path (`publishAsOwner`) persists values without firing triggers or the live bus. This is a real asymmetry: an inbound device publish fires all data triggers and routes to the live bus; a server-owned publish (AppEngine PUT/POST) does neither, even for `device-owner` publishes via `PublishDeviceValue`. Let me verify the upstream/design intent — is this documented as intended, or is it a gap?

Let me check the behavior of a server-owned publish on triggering, and whether there's a test asserting it. Let me look at the serverdata test around the publish and check whether the bus/triggers are verified there.
✱ Grep "PublishDeviceValue|device-owned|virtual|direct ingest|no MQTT" in internal · 96 matches
This confirms my read. The `publishAsOwner`/`PublishServerValue`/`PublishDeviceValue` path persists data but **never fires data triggers** and **never fans out to the live bus** (no `e.bus.Publish`, no `fireCommitted`/`fireData` equivalent). The server-owned and device-owned ingest paths bypass both the trigger evaluation and live-stream fan-out that the broker-driven ingress path gets via `fireCommitted`.

But wait — is this a behavior gap actually worth proposing, or is it deliberate and tested as-is? The tests assert persistence and wire delivery but don't check triggers/bus for these paths. Let me verify by reading the devicedata_test to confirm no trigger/bus assertion exists, and check whether these paths are documented. This is the strongest candidate (category 1: "a behaviour that is wrong or unguarded" — a rule that exists on one path but is silently absent on another).

Let me confirm by reading the devicedata_test fully and checking DESIGN.md reasoning if accessible.
→ Read internal/engine/devicedata_test.go
✱ Grep "EncodedLen" in pkg/deviceid · 6 matches
Confirmed: the server-owned and device-owned publish paths (`publishAsOwner`) persist rows but never fire data triggers nor fan out to the live bus. This is a genuine and unguarded behavior gap — the exact "a rule that exists on one path is silently absent on another" kind of find. The tests only assert persistence and wire delivery, never triggers/bus for these paths. This is my top proposal.

Let me look at a couple more potential issues quickly, then check `trigger` compile for the mustache `verifyTemplate` and other edge cases. Let me also reconsider a subtle one I noticed: in `Enqueue` (actions.go), when a policy is set, the in-flight count is incremented, and `deliver` decrements it via `defer`. But notice `Enqueue` increments `policyInflight` **before** the queue send, and if the queue is full it decrements. That's balanced. Good.

Let me check the `match.go` `TracksChanges` vs `matches` for the `path_removed` condition. In `fireDataChanges`, `removed` is only set when `op.PrevExists && op.Value == nil`. And `matches` for `OnPathRemoved` returns true without value check. Consistent.

Let me look at one more possible concern: the `valueMatches` for `not_contains` on a nil value. Actually `valueMatches` returns false early if `value == nil`. Fine.

Let me consider the `matchPathRE` and interface name regexes — established, tested.

Let me examine one more candidate I noticed: `splitDeviceTopic` uses `deviceIDStringLen = 22` hardcoded rather than `deviceid.EncodedLen`. Let me check if `deviceid.EncodedLen` exists and equals 22.
✱ Grep "deviceIDStringLen|deviceid.EncodedLen" in internal/engine · 4 matches
The `deviceIDStringLen = 22` in `topics.go:34` duplicates `deviceid.EncodedLen` (a magic-number duplication — same constant defined in two places). This is a minor clarity/maintenance issue, not behavior (they equal 22). It's category 4 (duplicated logic), but a very small change touching a hot path. The recipe says clarity proposals only when mechanical and tests cover behavior. This is mechanical but trivially small — arguably not worth a task on its own. I'll note it in the review's "decided not to propose" half.

Let me look at one more substantive candidate: in `actions.go` `Enqueue`/`deliver`, the policy in-flight accounting. Notice `Enqueue` increments `policyInflight` and sends to the channel. But what about the `Drop` path at the top (`<-x.closing`)? That returns before incrementing — fine. The full-queue path decrements — fine. `deliver` decrements via defer — fine.

However — one subtle thing: `channelOf` — when `d.Trigger.Policy() != nil` but the trigger is dropped at the closing check, no increment. OK.

Let me check whether `policyInflight` is ever inspected for a leak. It's incremented in `Enqueue` only when the send to channel succeeds. `deliver` decrements via defer on every path. Balanced.

Let me now consider the strongest proposals. I have:

1. **Server/device-owned publish paths bypass triggers and live bus** (serverdata.go) — strongest, category 1 behavior gap.

2. **Missing test for that gap** (category 2) — actually this is the same fix; the task is "fire triggers + bus on publishAsOwner, prove with a test."

Let me verify one more area: the `docs/DESIGN.md` is frozen spec, and the recipe says if code contradicts DESIGN.md, that's the code being buggy. But I can't read DESIGN.md (never-touch, but reading is fine). Actually DESIGN.md is on the never-touch list for *modification*. Reading is fine, but the recipe says don't propose changes to files under DESIGN.md. My proposal is about engine code, so fine.

Let me double check: is it possible server-values intentionally don't trigger? Let me reconsider. `PublishDeviceValue` exists so virtual devices (via flow) write rows "exactly like a real device's data would". A real device's data publish fires value_stored/incoming_data triggers and the bus. So a virtual device's publish producing no trigger events is inconsistent with the stated contract "lands storage rows exactly like a real device's data would." This strengthens the proposal.

But I should be careful — this could be deliberate and documented. The comment says "without delivering it... produces no MQTT traffic." It only says no MQTT; it says nothing about triggers/bus. The discrepancy with "exactly like a real device's data" is a real asymmetry worth flagging. Since this is a proposal job, I'll propose it as a task (with a test), and note in the review it needs confirmation of intent.

Now let me pick a few more well-supported tasks to make the output useful but within the 5-line cap. Let me scan the triggers tests to see if any rule lacks coverage. Let me look at what `match_test.go` covers and whether there's an obvious missing test.

Actually, let me not over-expand. The recipe says "at most five task lines." I have one strong one. Let me find 1-2 more genuinely machine-checkable proposals worth including.

Let me reconsider the `not_contains` gap: `valueMatches` for `not_contains` — `containable(value)` returns applicable only for string/array types. Fine, tests likely cover.

Let me look at whether there's a missing test for `TracksChanges` with a `path_removed` on a property path, or the `%{placeholder}` matching branches. Actually, let me check one concrete thing: the `segMatches` placeholder path and the `pathMatches` last-segment handling. These have tests presumably.

Let me instead verify a concurrency/robustness concern worth proposing: in `bus.go` `Publish`, the bus holds `b.mu` while sending to each subscriber channel via non-blocking select — not blocking, fine. In `Close`, it holds no lock during `s.close()`. Fine.

Let me settle on solid proposals. I'll go with:

1. **Server/device-owned publish paths never fire data triggers nor fan out to live bus** — behavior gap, with a failing test to prove. (serverdata.go, PublishServerValue/PublishDeviceValue/publishAsOwner)

Let me check one more: the `webhook` retry loop `backoff` variable and `time.After` inside a loop, and the delivery `attempt` — the `body` reused across attempts (bytes.NewReader each time) fine.

Let me look at whether the `Executor.enqueuedAt` stamp happens before the policy capacity check — it does (line 333 set right after closing check, before policy). For non-policy deliveries, enqueuedAt is set but unused. Fine.

I think I have a strong primary finding. Let me consolidate into the review file and 2-3 task lines. Let me check whether there's already coverage of the server-data trigger question in an existing test that would contradict (so I pick the right fix boundary). Let me quickly check if any test subscribes to the bus during a server publish.
✱ Grep "PublishServerValue|PublishDeviceValue" in internal/engine · 29 matches
called
No server-publish test asserts trigger or bus output — my finding stands: this path is unguarded for both trigger delivery and live fan-out.

I have enough. Let me confirm the review directory/date naming and the todo format, then check whether engine was reviewed before (it wasn't — only broker and flow).

Let me write the review file.

$ date +%Y-%m-%d
2026-09-06

→ Read .mule/todo.md
Now I understand the queue format. The recipe says to append tasks to `.mule/todo.md`. Note the header says the queue is authoritative on `mule/queue`, but the recipe (which I'm told to carry out) explicitly says "append at most five task lines to `.mule/todo.md`". The task instructions also say the deliverable is new task lines appended to `.mule/todo.md`. So I'll append there per the recipe.

Let me now also double-check a second solid proposal. Let me reconsider the tasks I could propose:

**Candidate A (strong):** Server/device-owned publish paths don't fire triggers or live bus — serverdata.go PublishServerValue/PublishDeviceValue/publishAsOwner. Behavior gap with a testable assertion.

**Candidate B (medium):** `deviceIDStringLen = 22` in topics.go duplicates `deviceid.EncodedLen` (magic-number duplication). Trivial clarity, but mechanical and behavior-covered by existing tests. Per recipe category 4, small enough that it's arguably not worth a ticket. I'll note in "decided not to propose."

Let me look for one more genuinely valuable, machine-checkable find. Let me reconsider the `TracksChanges` gate vs `matches` for the `path_removed` condition with a property on a datastream — no. 

Let me reconsider a concurrency subtlety in the trigger executor worth checking: in `Enqueue`, the policy capacity increment happens, then the send. If `Close` is called between the closing-check and the send, the send is protected by `sendMu.RLock`. Fine.

Let me look at whether `path_removed`'s `matches` handling is tested. In `valueMatches`, `OnPathRemoved` returns true before value check (good). But consider `OnPathCreated` — it goes through `valueMatches(m.operator, ev.Value, m.known)` with `ev.Value` being the new value (non-nil) — correct.

One thing worth checking: in `fireDataChanges`, `emit` is called for `OnValueChange` and `OnValueChangeApplied` with the SAME `body`? No — `emit(OnValueChange, NewValueChangeEvent...)` and later `emit(OnValueChangeApplied, NewValueChangeAppliedEvent...)`. Correct.

Let me look at a potential issue in `match.go` `matches` for `OnPathRemoved` when `anyPath` is false and `TracksChanges` — actually there's a subtle bug candidate: `TracksChanges` and `matches` for path_removed. Both use `pathMatches(m.pathSegs, ev.Path)`. Consistent.

I think Candidate A is the main deliverable. Let me also add a related test-coverage task since the gap is that no test asserts trigger/bus behavior on these paths. Actually, Candidate A IS the test-coverage task merged with the fix. Per the recipe I should propose the fix+test as one task.

Let me also consider proposing a second independent task from the triggers area. Let me look more carefully at the `valueMatches` `not_contains` and `contains` for a nil known value edge — probably covered.

Let me look at whether `Executor.Enqueue`'s policy capacity accounting can leak when the delivery is dropped at `Clos()`. Actually there's a subtle one: `Enqueue` increments `policyInflight` (line 344) only when capacity allows, and sends. But if `Close` is called *after* Enqueue succeeded and the worker reads the delivery, `deliver` decrements. Fine. But what if the queue send succeeds but then... no leak.

I'll finalize with two task lines: the primary server-publish trigger/bus fix, and I'll add a second smaller but real one. Let me reconsider the introspection handler: `parseIntrospection` — it doesn't validate the interface name against any real interface-name grammar, only rejects topic metacharacters and empty. Upstream may require valid names. But that's speculative; the design says "reject metacharacters". Not a clear gap.

Actually, let me reconsider `handleIntrospection` — it does `dev.setIntrospection(intro)` then fires events but never updates `rs` triggers — fine, that's per-device state.

I'll keep it tight with 1 primary task and 1 secondary (the duplication one is marginal; I'll put it in "decided not to propose"). Actually, let me look once more for a real second find in the hot path to make the review genuinely useful.

Let me re-examine `noteFormat`: it's called with `dp.Format` from the decode, but note that in `processData`, `noteFormat` is called (data.go:303) and `op.Format` is set to `dp.Format`. But `noteFormat` flips `resetHintOnBSON` — only on BSON with the armed flag. Looks fine.

Let me reconsider whether there's a missed `path_removed` trigger on the *server-owned UnsetServerProperty* path (serverdata.go:211) — unsetting a server property also doesn't fire triggers/bus. Reinforces Candidate A (the whole server-data path bypasses trigger/bus).

I'm confident. Let me also verify: is `fireData` reachable only from `fireCommitted` (the inbound path)? Yes. `publishAsOwner` never calls `afterCommit`/`fireData`/`bus.Publish`. This confirms the asymmetry spans both server-owned and (virtual) device-owned ingestion.

Let me write the review file and append one task line. Actually let me make the primary task precise. Since both PublishServerValue and PublishDeviceValue share `publishAsOwner`, the fix is to route the op through trigger evaluation + bus fan-out after persist. But note the server path also publishes to the broker (MQTT) — the design question is whether a server-owned value should ALSO fire `incoming_data`/`value_stored` triggers and appear on the live bus. Upstream semantics: a server-owned publish is still a data store, and the Dashboard live card and triggers subscribe to *events*, so yes it should appear. For device-owned (virtual) ingest, the stated contract "exactly like a real device's data would" strongly implies triggers + bus.

I'll write the task line recommending the fix and a test, flagging that intent should be confirmed since it changes observable event behavior.

Now write the files.

← Write .mule/reviews/engine-2026-09-06.md
Wrote file successfully.

I accidentally made a typo in the review file ("reivews" → I wrote "reivews" at line 3). Let me fix that.
← Edit .mule/reviews/engine-2026-09-06.md
Index: /root/astrate-mule/.mule/reviews/engine-2026-09-06.md
===================================================================
--- /root/astrate-mule/.mule/reviews/engine-2026-09-06.md
+++ /root/astrate-mule/.mule/reviews/engine-2026-09-06.md
@@ -1,7 +1,7 @@
 # Code review — internal/engine — 2026-09-06
 
-Area chosen because `broker` (2026-09-04) and `flow` (2026-09-05) were recent reivews and `engine` is the largest un-reviewed package (48 files).
+Area chosen because `broker` (2026-09-04) and `flow` (2026-09-05) were recent reviews and `engine` is the largest un-reviewed package (48 files).
 
 ## What I read
 
 - `internal/engine/router.go` (engine struct, sharding, Submit, Drain, metrics)

Now append one task line to `.mule/todo.md`. The recipe allows up to five; I found one clear, well-scoped finding that is a genuine behavior gap. Let me append it after the existing queue (I'll place it at the end).
→ Read .mule/todo.md [offset=108, limit=5]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -107,4 +107,5 @@
 - [x] docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` route to docs/api/astarte_realm_management_api.yaml — it exists in code (internal/realm/http.go:51, handler getDatastreamMaximumStorageRetention at http.go:135, returns 200 with a data-envelope retention value via Service.GetDatastreamMaximumStorageRetention) but the spec jumps straight from config/device_registration_limit (yaml:428) to /version. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] housekeeping-tests [auto]: add a container-free unit test file for `internal/housekeeping` (service.go, 262 lines, zero test coverage). Cover two rules: (1) `CreateRealm` validation — blank name, blank JWT key, negative regLimit, negative retention, and default retention injection when `WithDefaultDatastreamMaximumStorageRetention` is set and the caller passes nil; (2) `DeleteRealm` gating — `ErrDeletionDisabled` when `WithRealmDeletionDisabled(true)`, `ErrConnectedDevicesPresent` when `DeviceStats` reports connected > 0. Use httptest-style fakes for the store and sealer (no Docker).
 - [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
+- [ ] server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.

Now let me verify I haven't violated git rules — the task said don't touch git. I only wrote files. Let me confirm no git operations occurred and check git status (read-only) to report accurately.
$ git -C /root/astrate-mule status --short
M .mule/todo.md
?? .mule/reviews/engine-2026-09-06.md

I only wrote the two deliverable files and did not touch git. The recipe is a proposal job, so I did not implement the proposed task, and I did not run `gh issue create` (this recipe's deliverable is todo.md lines, not an issue).

Let me report per the MULE.md format.

```
Done: Reviewed internal/engine (triggers, data path, server data path); proposed one task line.
Files: .mule/reviews/engine-2026-09-06.md, .mule/todo.md
Verified: no gate run (no code changed); git status shows only the two deliverable files
Unsure: whether the server/device-owned publish paths are deliberately silent on triggers/bus — flagged in the task line for Giulio to confirm
Follow-ups: server-data-trigger-bus appended to .mule/todo.md
```

Note also: the todo.md header warns the queue is authoritative on `mule/queue`, but I was instructed to append to `.mule/todo.md` and did so, tagging the line `[auto]` per the existing convention.
