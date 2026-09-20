slug: recipe-code-review
verdict: proposed
at:  80599d7
ran: 2026-09-20T17:31:18Z on DietPi in 542s


> build · big-pickle

I'll read the required files first.
→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
Recipe is a proposal-only job: pick one package, review it, write a review file, append task lines to `.mule/todo.md`. No git, no code changes. Let me check what's been reviewed recently and see the package map.

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

$ ls -la .mule/reviews/ 2>/dev/null; echo "---LOG---"; cat .mule/log.md 2>/dev/null | tail -60; echo "---GITLOG---"; git log --oneline -20 2>/dev/null
total 136
drwxr-xr-x 2 root root 4096 Sep 18 22:00 .
drwxr-xr-x 7 root root 4096 Sep 19 19:08 ..
-rw-r--r-- 1 root root 4600 Sep  7 20:15 appengine-2026-09-07.md
-rw-r--r-- 1 root root 3443 Sep 15 22:09 astarteapi-2026-09-15.md
-rw-r--r-- 1 root root 5192 Sep 14 13:34 auth-2026-09-14.md
-rw-r--r-- 1 root root 8188 Sep  4 20:43 broker-2026-09-04.md
-rw-r--r-- 1 root root 6967 Sep 18 22:00 channels-2026-09-18.md
-rw-r--r-- 1 root root 4340 Sep 13 19:26 config-2026-09-13.md
-rw-r--r-- 1 root root 5264 Sep 18 19:25 deviceid-2026-09-18.md
-rw-r--r-- 1 root root 3141 Sep  6 20:57 engine-2026-09-06.md
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
---LOG---
| 2026-09-08 | docs-sync-appengine-device-status-schema | done | 255s | 00da133 |
| 2026-09-08 | docs-sync-appengine-data-set-422 | blocked | 251s | wrote nothing |
| 2026-09-08 | docs-sync-appengine-downsample-min | done | 86s | db506f9 |
| 2026-09-09 | pairing-bearer-secret-test | done | 90s | 32544c0 |
| 2026-09-09 | pairing-remoteip-fallback-test | done | 469s | 3d6e790 |
| 2026-09-09 | docs-sync-hk-patch-endpoint | done | 403s | 9fc9f33 |
| 2026-09-09 | docs-sync-hk-retention-field | done | 142s | d6d7c7f |
| 2026-09-09 | realm-pure-helper-tests | done | 403s | 39eb70a |
| 2026-09-09 | docs-sync-native-compat-health-503 | done | 202s | 656e21e |
| 2026-09-09 | docs-sync-native-compat-version-endpoints | done | 480s | d148f16 |
| 2026-09-10 | flowapi-autorestart-terminal-failure | done | 251s | 8ce04d2 |
| 2026-09-10 | flowapi-validationdetail-test | done | 121s | 832d8fb |
| 2026-09-10 | flowapi-autorestart-default-test | done | 148s | e272cfa |
| 2026-09-10 | docs-sync-pairing-status-enum | done | 89s | 2f6176c |
| 2026-09-10 | docs-sync-pairing-version-endpoint | done | 95s | c867b2c |
| 2026-09-11 | docs-sync-rm-delete-interface-status | done | 123s | b37d11f |
| 2026-09-11 | docs-sync-rm-mapping-required-encrypted | done | 137s | 5f8cb25 |
| 2026-09-11 | docs-sync-rm-put-auth-422 | done | 145s | 8554720 |
| 2026-09-11 | docs-sync-rm-version-example | blocked | 97s | opencode exited 1 |
| 2026-09-11 | interfaceschema-compat-attrs-flip-test | done | 329s | 4fe1f4f |
| 2026-09-11 | interfaceschema-object-attrs-uniformity | done | 344s | c6c8362 |
| 2026-09-11 | interfaceschema-malformed-placeholder-fixtures | done | 130s | f6a44a2 |
| 2026-09-11 | interfaceschema-enum-roundtrip-coverage | done | 185s | f955958 |
| 2026-09-11 | interfaceschema-trie-interior-literal-param | done | 121s | 705795f |
| 2026-09-12 | docs-sync-ae-write-405 | done | 400s | 9ed08c3 |
| 2026-09-12 | docs-sync-ae-write-value-422 | done | 338s | 17f6f70 |
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
---GITLOG---
80599d7 mule: recipe astarte-upstream ran (115s)
1dba3d5 mule: refill the queue
bba3b36 mule: recipe github-issues ran (99s)
74021c8 mule: refill the queue
f166917 mule: recipe milestones ran (124s)
c0c7959 mule: recipe hygiene timed out
757407a mule: recipe docs-sync ran (114s)
054d305 mule: blocked channels-rejoin-authz-mismatch
ebc3854 mule: log channels-rejoin-joinref-tagging-test
d48e17f mule: channels-rejoin-joinref-tagging-test: pin in internal/appengine/channels/ws_test.go that a second `phx_join` on the same topic updates the session's `joinRef` (ws.go:284-291) and subsequent `new_event` pushes carry the fresh ref (read per event at ws.go:403) — no test drives a rejoin today (all `TestWireSession_*` join once), so the comment's "server pushes must be tagged with the current one" rots unasserted. Wire: join with ref "1", rejoin with "2", push an event, assert element 0 is `"2"`.
ab18c72 mule: log channels-group-watch-membership-scope
aeaeb47 mule: channels-group-watch-membership-scope: group watches authorize against `WATCH::groups/<name>/…` (ws.go:178-205, conformance channels.json) but `Room.dispatch` ever only filters by `watchEntry.deviceID` (empty for a group watch — room.go:36-40, 196-200) and the trigger compiles with no group/device scope (triggers/match.go:374-396, because upstream carries group_name at the payload top level), so a group-scoped token receives every realm event matching interface/path/on, members or not. Plumb a group-membership resolver into internal/appengine/channels (it knows only `Bus` today), resolve membership at Watch time and store it in `watchEntry`, filter deliveries by it. Add a container-free test with a fake resolver: watch group `probe`, publish an event from a member and a non-member, assert only the member's arrives. Verify the member-filter rule against upstream on the Legion if disputed; the boundary leak stands regardless.
7bf6f8e mule: refill the queue
809d541 mule: recipe code-review ran (683s)
454ef60 mule: recipe astarte-upstream ran (209s)
1c7d617 mule: recipe github-issues ran (119s)
79439e6 mule: refill the queue
b3c4bad mule: recipe milestones ran (176s)
eb95cdd mule: recipe hygiene timed out
4c4cf5b mule: recipe docs-sync ran (421s)

$ for d in internal/* pkg/*; do if [ -d "$d" ]; then n=$(rg -c '' --glob '*.go' "$d" 2>/dev/null | awk -F: '{s+=$2} END {print s}'); echo "$n $d"; fi; done | sort -rn
15079 internal/engine
11157 internal/flow
7435 internal/appengine
6783 internal/store
3576 internal/broker
3407 pkg/interfaceschema
2942 internal/pairing
2709 internal/realm
2690 pkg/payload
2335 internal/flowapi
1588 internal/auth
1372 internal/housekeeping
1004 internal/config
697 internal/testutil
525 pkg/astarteapi
451 internal/observability
369 internal/httpx
358 pkg/deviceid
171 internal/swagger

Now let me check which areas were reviewed recently and sizes, to pick an un-reviewed one. No review file exists for `internal/engine` since 2026-09-06, `broker` since 2026-09-04, `swagger`, or `testutil`.
$ rg -c '' --glob '*.go' internal/broker | sort -t: -k2 -rn | head; echo "---ENGINE---"; rg -c '' --glob '*.go' internal/engine | sort -t: -k2 -rn | head
internal/broker/broker_e2e_test.go:618
internal/broker/broker_test.go:549
internal/broker/sessionstore.go:533
internal/broker/authhook.go:480
internal/broker/broker.go:343
internal/broker/sessionstore_test.go:227
internal/broker/aclhook.go:221
internal/broker/aclhook_test.go:146
internal/broker/identity.go:104
internal/broker/identity_test.go:103
---ENGINE---
internal/engine/triggers/actions_test.go:1074
internal/engine/fakestore_test.go:818
internal/engine/triggers/match.go:769
internal/engine/triggers/actions.go:610
internal/engine/triggers/match_test.go:594
internal/engine/data.go:552
internal/engine/cache.go:502
internal/engine/router.go:490
internal/engine/data_test.go:461
internal/engine/batch.go:459

`internal/broker` is the oldest review (2026-09-04), hasn't been touched since. Let me read the previous broker review to avoid repeating its work, then go deep.
$ ls internal/broker/; echo "---"; rg -c '' --glob '*.go' internal/broker | sort -t: -k2 -rn
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
---
internal/broker/broker_e2e_test.go:618
internal/broker/broker_test.go:549
internal/broker/sessionstore.go:533
internal/broker/authhook.go:480
internal/broker/broker.go:343
internal/broker/sessionstore_test.go:227
internal/broker/aclhook.go:221
internal/broker/aclhook_test.go:146
internal/broker/identity.go:104
internal/broker/identity_test.go:103
internal/broker/intake.go:100
internal/broker/lifecycle.go:90
internal/broker/publisher.go:62

→ Read .mule/reviews/broker-2026-09-04.md
All main packages have been reviewed; the rotation has gone full circle. The oldest/thinnest reviews are broker (09-04), flow (09-05), and engine (09-06). Let me check those to pick the best one.
→ Read .mule/reviews/engine-2026-09-06.md
→ Read .mule/reviews/flow-2026-09-05.md
→ Read .mule/todo.md
Broker (09-04) is the least-recently reviewed but its review was deep and thorough. Engine (09-06) is the largest package (15k lines) and its review was the thinnest relative to size. Let me check the current todo.md and what happened with previously-proposed tasks before deciding.

$ git log --oneline -60 | rg -i 'broker|engine|review|code-review' ; echo "---"; git log --oneline --all | rg -i 'broker-acl|broker-disconnect|broker-offline|broker-onconnect' | head
d48e17f mule: channels-rejoin-joinref-tagging-test: pin in internal/appengine/channels/ws_test.go that a second `phx_join` on the same topic updates the session's `joinRef` (ws.go:284-291) and subsequent `new_event` pushes carry the fresh ref (read per event at ws.go:403) — no test drives a rejoin today (all `TestWireSession_*` join once), so the comment's "server pushes must be tagged with the current one" rots unasserted. Wire: join with ref "1", rejoin with "2", push an event, assert element 0 is `"2"`.
aeaeb47 mule: channels-group-watch-membership-scope: group watches authorize against `WATCH::groups/<name>/…` (ws.go:178-205, conformance channels.json) but `Room.dispatch` ever only filters by `watchEntry.deviceID` (empty for a group watch — room.go:36-40, 196-200) and the trigger compiles with no group/device scope (triggers/match.go:374-396, because upstream carries group_name at the payload top level), so a group-scoped token receives every realm event matching interface/path/on, members or not. Plumb a group-membership resolver into internal/appengine/channels (it knows only `Bus` today), resolve membership at Watch time and store it in `watchEntry`, filter deliveries by it. Add a container-free test with a fake resolver: watch group `probe`, publish an event from a member and a non-member, assert only the member's arrives. Verify the member-filter rule against upstream on the Legion if disputed; the boundary leak stands regardless.
809d541 mule: recipe code-review ran (683s)
42acc40 mule: recipe code-review ran (320s)
1514f97 mule: recipe code-review ran (246s)
08048b3 mule: docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
8875c1d mule: recipe code-review ran (270s)
---
8f6dee4 mule: blocked broker-onconnect-doc-comment
9c49516 mule: blocked broker-offline-acl-tests
b8b48aa mule: blocked broker-disconnect-device-zombie-session
7a205dc mule: log broker-acl-coldstart-introspection-miss
d201db4 mule: broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [approved 2026-09-04]

The broker review was thorough; engine is the largest package (15k lines) yet its review was the thinnest (29 lines, one finding). Let me see what engine files the 09-06 review did NOT cover.
$ rg -c '' --glob '*.go' internal/engine internal/engine/triggers internal/engine/stream | sort -t: -k2 -rn
internal/engine/triggers/actions_test.go:1074
internal/engine/triggers/actions_test.go:1074
internal/engine/fakestore_test.go:818
internal/engine/triggers/match.go:769
internal/engine/triggers/match.go:769
internal/engine/triggers/actions.go:610
internal/engine/triggers/actions.go:610
internal/engine/triggers/match_test.go:594
internal/engine/triggers/match_test.go:594
internal/engine/data.go:552
internal/engine/cache.go:502
internal/engine/router.go:490
internal/engine/data_test.go:461
internal/engine/batch.go:459
internal/engine/engine.go:456
internal/engine/batch_test.go:446
internal/engine/engine_e2e_test.go:440
internal/engine/engine_integration_test.go:435
internal/engine/triggers/policy_test.go:410
internal/engine/triggers/policy_test.go:410
internal/engine/serverdata_test.go:390
internal/engine/cache_test.go:363
internal/engine/router_test.go:362
internal/engine/serverdata.go:349
internal/engine/control_test.go:330
internal/engine/control.go:316
internal/engine/triggers/events.go:312
internal/engine/triggers/events.go:312
internal/engine/forward/http_test.go:310
internal/engine/triggers/policy.go:297
internal/engine/triggers/policy.go:297
internal/engine/triggers_change_test.go:293
internal/engine/capabilities_test.go:255
internal/engine/stream/bus_test.go:234
internal/engine/stream/bus_test.go:234
internal/engine/devicedata_test.go:223
internal/engine/stream/bus.go:209
internal/engine/stream/bus.go:209
internal/engine/forward/nats_test.go:196
internal/engine/introspection_test.go:189
internal/engine/introspection.go:168
internal/engine/triggers/validation.go:163
internal/engine/triggers/validation.go:163
internal/engine/triggers/change_test.go:161
internal/engine/triggers/change_test.go:161
internal/engine/triggers/events_test.go:155
internal/engine/triggers/events_test.go:155
internal/engine/mqttv1_parity_test.go:151
internal/engine/triggers/transient_test.go:142
internal/engine/triggers/transient_test.go:142
internal/engine/triggers_change_integration_test.go:141
internal/engine/bus_events_test.go:140
internal/engine/forward/http.go:125
internal/engine/topics_test.go:109
internal/engine/capabilities.go:102
internal/engine/triggers/errorname_test.go:98
internal/engine/triggers/errorname_test.go:98
internal/engine/topics.go:97
internal/engine/triggers/errorname.go:79
internal/engine/triggers/errorname.go:79
internal/engine/forward/nats.go:53
internal/engine/forward/envelope.go:29
internal/engine/triggers/transient.go:22
internal/engine/triggers/transient.go:22

$ ls internal/engine/ internal/engine/triggers/ internal/engine/stream/ 2>/dev/null
internal/engine/:
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

The 09-06 engine review thoroughly covered the core files but skipped several: `forward/` (http, nats, envelope), `triggers/validation.go`, `triggers/transient.go`, and `triggers/errorname.go`. Since engine is the largest package with the thinnest review, I'll re-review it focusing on the uncovered files.
→ Read internal/engine/forward/http.go
→ Read internal/engine/forward/nats.go
→ Read internal/engine/forward/envelope.go
→ Read internal/engine/triggers/validation.go
Now the uncovered trigger files: validation.go, transient.go, errorname.go.
→ Read internal/engine/triggers/transient.go
→ Read internal/engine/triggers/errorname.go
→ Read internal/engine/forward/http_test.go
→ Read internal/engine/forward/nats_test.go
Now let me look at the tests for these files to check coverage, and check git history for what's changed recently.
✱ Grep "Forward\(|Forwarder|retry|backoff|maxRetries" in internal/engine/triggers/actions.go · 100 matches
The forward package is well-covered. Let me check how the Forwarder is actually used in engine — error handling and retry semantics.
→ Read internal/engine/triggers/actions.go [offset=1, limit=120]
Let me read the rest of actions.go, focusing on delivery, retry, and the HTTP client construction.
→ Read internal/engine/triggers/actions.go [offset=120, limit=490]
✱ Grep "forward\.New|Forwarder|forward\.Config|triggers\.forward|forward:" 100 matches
Let me check how the forwarder is wired in config/cmd to verify timeouts and default client behaviour.
→ Read internal/engine/triggers/actions.go [offset=609]
→ Read internal/config/config.go [offset=338, limit=40]
Let me look at the config validation for the forward method to see if config and forward.New agree, and check main.go's newForwarder.
→ Read cmd/astrate/main.go [offset=105, limit=175]
Let me check the errorname and transient tests for coverage gaps.
→ Read internal/engine/triggers/errorname_test.go
→ Read internal/engine/triggers/transient_test.go
Let me verify a key suspicion: whether delivery policies apply to custom (forwarder) actions, or only to HTTP webhooks.

$ rg -n "Policy\(\)|policy" internal/engine/triggers/actions.go; echo "===TEST POLICY vs CUSTOM==="; rg -n -i "policy" internal/engine/triggers/actions_test.go | rg -i "forward|custom" ; echo "===glob for policy usage==="; rg -rn "\.Policy\(\)" internal/engine/triggers/ | rg -v _test
247:	policyMu       sync.Mutex
248:	policyInflight map[string]int
300:		policyInflight: make(map[string]int),
335:	if pol := d.Trigger.Policy(); pol != nil {
336:		x.policyMu.Lock()
337:		if x.policyInflight[pol.Name] >= pol.MaximumCapacity {
338:			x.policyMu.Unlock()
340:			x.log.Warn("trigger delivery policy capacity full; event dropped",
341:				"realm", d.Realm, "trigger", d.Trigger.Name, "policy", pol.Name)
344:		x.policyInflight[pol.Name]++
345:		x.policyMu.Unlock()
351:		if pol := d.Trigger.Policy(); pol != nil {
352:			x.policyMu.Lock()
353:			x.policyInflight[pol.Name]--
354:			x.policyMu.Unlock()
397:	if pol := d.Trigger.Policy(); pol != nil {
398:		// Balances the in-flight count Enqueue took for this policy, on
401:			x.policyMu.Lock()
402:			x.policyInflight[pol.Name]--
403:			x.policyMu.Unlock()
409:				"realm", d.Realm, "trigger", d.Trigger.Name, "policy", pol.Name,
490:// When the trigger has no attached policy (Policy() == nil), behaviour is the
494:// When a policy is attached, it governs every retry/discard decision: after
496:// passed to policy.Decide. StrategyRetry retries with the same exponential
498:// policy.RetryTimes+1 (RetryTimes counts retries, not attempts). Success
499:// remains status in 200..399, checked before consulting the policy.
502:// treated as a server error: upstream's policy schema speaks only in status
506:	policy := d.Trigger.Policy()
509:	if policy != nil {
510:		maxAttempts = policy.RetryTimes + 1
518:		if policy != nil {
523:			dec := policy.Decide(effectiveStatus)
526:				x.log.Warn("webhook delivery discarded by policy",
528:					"status", status, "policy", policy.Name, "reason", dec.Reason)
536:					"policy", policy.Name, "reason", dec.Reason)
541:				"status", status, "policy", policy.Name, "reason", dec.Reason)
548:					"policy", policy.Name, "reason", dec.Reason)
===TEST POLICY vs CUSTOM===
===glob for policy usage===
internal/engine/triggers/actions.go:	if pol := d.Triggern; pol != nil {
internal/engine/triggers/actions.go:		if pol := d.Triggern; pol != nil {
internal/engine/triggers/actions.go:	if pol := d.Triggern; pol != nil {
internal/engine/triggers/actions.go:	policy := d.Triggern

✱ Grep "CompileCondition|transient|Forwarder" in internal/appengine/channels · 3 matches
Let me verify the `missing_required` mapping coverage and check the DESIGN note on the forwarder/intent.

$ rg -n "missing_required" internal/ pkg/ docs/ 2>/dev/null | rg -v ".\\.git"; echo "===DESIGN 1.1==="; rg -n "forward|Forwarder|AMQP|trigger" docs/DESIGN.md | head -30
internal/engine/triggers/errorname.go:62:	// Upstream 1.4 emits missing_required_mapping, outside dashboard
internal/engine/triggers/errorname.go:64:	"missing_required":  "unexpected_object_key",
pkg/payload/value.go:158:	"missing_required",
pkg/payload/value_test.go:199:		ReasonMissingRequired: "missing_required",
===DESIGN 1.1===
22:| Message bus | RabbitMQ (inter-service AMQP) | In-process Go channels (sharded worker pools) |
57:| **VerneMQ** + `astarte_vmq_plugin` | MQTT broker, mTLS termination, ACLs, bridging publishes onto AMQP | `internal/broker` | Embedded `mochi-mqtt/server` v2 with auth/ACL hooks; publishes flow into the engine through a Go channel, not AMQP. |
58:| **Data Updater Plant (DUP)** | Consumes device messages from AMQP, validates against interfaces, writes to Cassandra, detects trigger conditions | `internal/engine` | The heart of Astrate: sharded per-device ordered pipeline → validation → persistence → trigger evaluation → live fan-out. |
60:| **Realm Management API** | Interface (schema) and trigger CRUD per realm | `internal/realm` | Interfaces stored as JSONB + compiled in-memory (§2.6). Trigger CRUD included; trigger *execution* lives in `internal/engine/triggers`. |
62:| **Trigger Engine** | Executes trigger actions (HTTP webhooks, AMQP messages) | `internal/engine/triggers` | HTTP webhook actions with retry/backoff; AMQP action replaced by optional NATS/HTTP forwarding (extension point). |
99:│   │   ├── triggers/         # trigger matching + action execution (HTTP webhooks)
104:│   ├── realm/                # /realmmanagement/v1 API: interfaces, triggers CRUD
136:  to the device exactly as a slow AMQP consumer would in upstream Astarte. QoS 0 messages are
274:CREATE TABLE triggers (
278:    definition jsonb NOT NULL,        -- Astarte trigger JSON (simple_triggers + action)
428:   `device_error`/`unexpected_value` trigger event + metric).
447:with device/interface/path, and feeds `device_error` triggers (parity with Astarte's DUP).
494:  `device_connected`/`device_disconnected` triggers — the work `astarte_vmq_plugin` does
513:| `<realm>/<device_id>` | Introspection string: `;`-separated `interface_name:major:minor` triples (UTF-8 plain text, e.g. `com.ex.Sensors:1:0;com.ex.Geo:0:1`) | 2 | Parse; diff against stored introspection; update `devices.introspection` (+`old_introspection` for removed pairs); fire `incoming_introspection` / interface added/removed triggers; recompute the device's server-owned subscription expectations. |
516:| `<realm>/<device_id>/<interface_name><path>` (path always starts with `/`) | BSON **or JSON** document `{ "v": <value>, "t": <timestamp, optional> }`; empty payload = property unset | per-mapping reliability | Full pipeline §2.6: validate → persist (upsert property / insert datastream row) → triggers → live fan-out. |
521:`device_error` trigger, never a crash.
555:else                                → reject (metrics + device_error trigger)
625:- `/realmmanagement/v1/<realm>/interfaces...`, `/triggers...`.
787:  batch flush latency, broker sessions, DB pool stats, trigger delivery outcomes.

✱ Grep "strategy|discard|retry|policy" in internal/engine/triggers/match.go · 100 matches
Let me check whether a custom (forwarder) action can carry a policy, and how that interacts with `forward()`.
✱ Grep "AttachPolicy|PolicyName|policy" in internal/engine/triggers · 9 matches
→ Read internal/engine/triggers/match.go [limit=120, offset=1]
✱ Grep "AttachPolicy" 15 matches
Let me confirm policies can be attached to custom-action triggers and see how they resolve.
→ Read internal/engine/triggers/match.go [limit=120, offset=225]
Let me read the forwarder test section in actions_test.go to see what's pinned.
→ Read internal/engine/triggers/actions_test.go [limit=70, offset=555]
Let me verify `missing_required` is actually emitted in the engine (not a dead mapping), which determines whether the missing test row is meaningful.
✱ Grep "func TestPolicy|policy" in internal/engine/triggers/actions_test.go · 100 matches

$ rg -n "ReasonMissingRequired|missing_required" internal/engine pkg/payload --glob '*.go' | rg -v errorname.go; echo "==="; rg -n "missing_required" pkg/payload/value.go pkg/payload/value_test.go
pkg/payload/value_test.go:199:		ReasonMissingRequired: "missing_required",
pkg/payload/value.go:149:	// ReasonMissingRequired: an object-aggregation document omits a key
pkg/payload/value.go:151:	ReasonMissingRequired
pkg/payload/value.go:158:	"missing_required",
pkg/payload/value.go:163:	if r >= ReasonTooLarge && r <= ReasonMissingRequired {
pkg/payload/payload.go:148:			return DecodedPayload{}, rejectf(ReasonMissingRequired,
pkg/payload/required_test.go:13:// declares required is rejected with ReasonMissingRequired in both wire
pkg/payload/required_test.go:38:	if _, err := DecodeObject([]byte(`{"v":{"lon":9.0}}`), leaves); ReasonOf(err) != ReasonMissingRequired {
pkg/payload/required_test.go:39:		t.Fatalf("JSON rejection reason = %v (err %v); want %v", ReasonOf(err), err, ReasonMissingRequired)
pkg/payload/required_test.go:47:	if _, err := DecodeObject(bsonBad, leaves); ReasonOf(err) != ReasonMissingRequired {
pkg/payload/required_test.go:48:		t.Fatalf("BSON rejection reason = %v (err %v); want %v", ReasonOf(err), err, ReasonMissingRequired)
===
pkg/payload/value.go:158:	"missing_required",
pkg/payload/value_test.go:199:		ReasonMissingRequired: "missing_required",

I have my findings. Let me confirm the errorname test omission precisely (map has 20 entries, test table has 19 rows).
$ echo "map entries:"; rg -o '"[a-z_]+"\s*:\s*"[a-z_]+"' internal/engine/triggers/errorname.go | wc -l; echo "test rows:"; rg -o '"[a-z_]+", "[a-z_]+"' internal/engine/triggers/errorname_test.go | wc -l; echo "--- map keys not in test:"; rg -o '^\s*"[a-z_]+":' internal/engine/triggers/errorname.go | sed 's/[": ]//g' | sort > /tmp/mapkeys; rg -o '\{"[a-z_]+",' internal/engine/triggers/errorname_test.go | sed 's/{"//;s/,//' | sort > /tmp/testkeys; comm -23 /tmp/mapkeys /tmp/testkeys
map entries:
20
test rows:
19
--- map keys not in test:
	bad_object
	bad_timestamp
	control_payload_invalid
	control_unknown
	device_unknown
	interface_not_in_introspection
	interface_not_installed
	introspection_invalid
	malformed
	malformed_topic
	missing_required
	no_value
	ownership_violation
	realm_unknown
	too_large
	type_mismatch
	unexpected_path
	unknown_format
	unset_not_allowed
	value_too_large

The `comm` shows `missing_required` among map keys not matched in the test list format. Let me confirm directly.
$ rg -n 'missing_required' internal/engine/triggers/errorname_test.go; echo "exit=$?"; echo "=== map lines ==="; rg -n '^\s*"[a-z_]+":\s*"[a-z_]+",$' internal/engine/triggers/errorname.go
exit=1
=== map lines ===
44:	"realm_unknown":                  "device_session_not_found",
45:	"device_unknown":                 "device_session_not_found",
46:	"malformed_topic":                "invalid_path",
47:	"interface_not_in_introspection": "interface_loading_failed",
48:	"interface_not_installed":        "interface_loading_failed",
49:	"ownership_violation":            "write_on_server_owned_interface",
50:	"unexpected_path":                "mapping_not_found",
51:	"introspection_invalid":          "invalid_introspection",
52:	"control_unknown":                "unexpected_control_message",
53:	"control_payload_invalid":        "unexpected_control_message",
54:	"too_large":                      "value_size_exceeded",
55:	"unknown_format":                 "undecodable_bson_payload",
56:	"malformed":                      "undecodable_bson_payload",
57:	"no_value":                       "unexpected_value_type",
58:	"bad_timestamp":                  "unexpected_value_type",
59:	"type_mismatch":                  "unexpected_value_type",
60:	"value_too_large":                "value_size_exceeded",
61:	"bad_object":                     "unexpected_object_key",
64:	"missing_required":  "unexpected_object_key",
65:	"unset_not_allowed": "unexpected_value_type",

Confirmed both findings. `missing_required` is the only `astrateToUpstream` key with no test row (bare `rg` on the test file finds nothing), and the policy/retry machinery in `forward()` never consults an attached policy. Let me verify git state and write the review.
$ git status --short; git log --oneline -3
80599d7 mule: recipe astarte-upstream ran (115s)
1dba3d5 mule: refill the queue
bba3b36 mule: recipe github-issues ran (99s)

Both findings confirmed and the tree is clean. Writing the review file and appending the task lines.
← Write .mule/reviews/engine-triggers-2026-09-20.md
Wrote file successfully.

Now append the two task lines to `.mule/todo.md`.
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -1,4 +1,7 @@
+- [ ] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
+- [ ] triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
+
 # Mule queue
 
 One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).
 The mule runs the topmost open line. Reorder freely; the order is the priority.

→ Read .mule/todo.md [limit=6]
Reviewed `internal/engine` second pass (triggers/forward subpackages — the areas the 09-06 engine review and rotation hadn't covered).

Done: reviewed internal/engine triggers+forward; proposed 2 task lines in .mule/todo.md, wrote .mule/reviews/engine-triggers-2026-09-20.md
Files: .mule/reviews/engine-triggers-2026-09-20.md, .mule/todo.md
Verified: rg/regex checks (map keys 20 vs test rows 19, `missing_required` absent from test), git status clean -> pass
Unsure: nothing
Follow-ups: errorname-missing-required-test [auto], triggers-custom-action-policy-nodecide [auto] (both appended; no git operations performed)
