slug: recipe-code-review
verdict: proposed
at:  454ef60
ran: 2026-09-18T20:00:41Z on DietPi in 683s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
Recipe understood: deep-review one package, write an evidence file, append task lines. No code, no git.

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

$ ls .mule/reviews/ 2>/dev/null; echo "---"; tail -40 .mule/log.md 2>/dev/null
appengine-2026-09-07.md
astarteapi-2026-09-15.md
auth-2026-09-14.md
broker-2026-09-04.md
config-2026-09-13.md
deviceid-2026-09-18.md
engine-2026-09-06.md
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
---
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

All 17 packages have a review within the last fortnight. The oldest are broker, flow, engine, appengine. Let me check sizes and pick one.
$ for d in internal/*/ pkg/*/; do n=$(rg -c '' --glob '*.go' "$d" 2>/dev/null | awk -F: '{s+=$2} END {print s+0}'); echo "$n $d"; done | sort -rn
15079 internal/engine/
11157 internal/flow/
7238 internal/appengine/
6759 internal/store/
3576 internal/broker/
3407 pkg/interfaceschema/
2942 internal/pairing/
2709 internal/realm/
2690 pkg/payload/
2335 internal/flowapi/
1588 internal/auth/
1372 internal/housekeeping/
1004 internal/config/
697 internal/testutil/
525 pkg/astarteapi/
451 internal/observability/
369 internal/httpx/
358 pkg/deviceid/
171 internal/swagger/

Oldest un-rotated area: `broker` (reviewed 09-04, 14 days ago), a manageable 3.6k lines. Let me check its previous review, then the package map.
$ ls internal/broker/
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

→ Read .mule/reviews/broker-2026-09-04.md
Broker's previous review proposed 4 tasks. Let me check if they were actioned, so I don't re-propose stale ones, and check for-giulio.md.

$ grep -n "broker" .mule/todo.md; echo "=== log ==="; grep -n "broker" .mule/log.md | tail -20; echo "=== for-giulio ==="; [ -f .mule/for-giulio.md ] && tail -15 .mule/for-giulio.md || echo "(none)"
82:- [ ] broker-external-bus-intake [legion]: implement the second Intake implementation named in the TODO at internal/broker/intake.go:56-62 — an external-bus-backed (e.g. NATS JetStream) Intake for multi-instance deployment / restart survival, reproducing the embedded broker's per-device ordering, deferred-ack backpressure and QoS 0 drop semantics; an implementation of this interface needs containerised integration tests, so it runs on the Legion Go.
98:- [ ] race-check-engine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/engine/... ./internal/broker/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
103:- [x] broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [approved 2026-09-04]
104:- [!] broker-disconnect-device-zombie-session: in `internal/broker/broker.go` `DisconnectDevice` (lines 260-266), force-closing the live connection leaves the persisted bbolt session (subscriptions + offline queue, keyed by CN) in place — mochi derives `expire` from the Clean flag, not the reason code (server.go:485-491), so `dropClient` (sessionstore.go:500) is never reached and a later device with the same CN resurrects the deleted device's session. Wire the session store handle into `Broker` and call `dropClient` from `DisconnectDevice`, and add a T1 test in `broker_test.go`: connect clean=false + subscribe, call `DisconnectDevice`, reconnect — assert `session_present` is false and the stale subscription is gone. [approved 2026-09-04] — BLOCKED: wrote nothing
105:- [!] broker-offline-acl-tests: in `internal/broker/aclhook_test.go`, unit-test the offline-delivery ACL — `offlineACL.ownershipOf` cache hit within TTL, TTL expiry triggering a reload, and a load failure caching as empty and denying (aclhook.go:116-138) — plus a T1 `OnACLCheck` check that the offline branch (delivery to a session-present but disconnected device, aclhook.go:196-202) consults the store. The rule currently has zero direct tests; `TestCheckACLMatrix` exercises only the pure `checkACL`. [approved 2026-09-04] — BLOCKED: wrote nothing
106:- [!] broker-onconnect-doc-comment: in `internal/broker/authhook.go:314`, restore the missing first line of the `OnConnect` doc comment — it opens mid-sentence ("of the Astarte MQTT v1 protocol, and mochi publishes them without a publish-side ACL check…") and the Will-clearing security rationale (a retained LWT would escape the §3.2 matrix) is only half-documented. [approved 2026-09-04] — BLOCKED: tests failed: --- FAIL: TestMQTTSink_Retained (0.02s)
108:- [x] purge-properties-compression-capability [auto]: upstream v1.3.0 adds a `purge_properties_compression_format` device capability (`zlib`|`plaintext`, default `zlib`) — a wire-visible capability value. Check whether Astrate's capabilities handling (internal/broker, the `<realm>/<device_id>/capabilities` topic, issue #16) needs to recognize/honour it, or whether zlib-only is already the deliberate default; propose the change or note why not needed.
127:- [x] server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.
=== log ===
11:| 2026-09-04 | broker-acl-coldstart-introspection-miss | done | 629s | d201db4 |
13:| 2026-09-04 | broker-disconnect-device-zombie-session | blocked | 117s | wrote nothing |
14:| 2026-09-04 | broker-offline-acl-tests | blocked | 208s | wrote nothing |
15:| 2026-09-04 | broker-onconnect-doc-comment | blocked | 438s | tests failed: --- FAIL: TestMQTTSink_Retained (0.02s) |
=== for-giulio ===
  `/root/astrate-mule` already solved the same way (your call, confirmed 2026-07-27: align it
  with the existing unattended-Pi-clone convention rather than fix gpg-agent or bypass signing
  silently inside the script).

## 2026-08-23 — FDO promoted to milestone-4.0 candidate (Giulio's decision, recorded)

#78 is no longer parked: zero-touch onboarding is strategic for commercial
viability. Scope frozen on the issue: owner-side TO1/TO2 in our Pairing
service only (last mile, like upstream), reuse fdo-rs for
manufacturing/rendezvous, acceptance = official `astarte-device-fdo-rust`
SDK completes onboarding against Astrate, docs as a first-class deliverable.
When v3.0 is marked DONE, the v4.0 section of `.mule/milestones.md` should be
drafted with this investigation as its first item (issue #78 has the full
verified context).
- **The mule has been idle 17h.** Filed by the dead-man's switch; see journalctl on the Pi.

→ Read .mule/reviews/appengine-2026-09-07.md
All packages have been reviewed; the oldest rotation is broker (09-04) but it has pending blocked tasks already queued. The next-oldest un-rotated areas: flow (09-05), engine (09-06), appengine (09-07). Given my context budget, appengine (7.2k lines) is the best deep-read candidate. Let me check its prior review and structure.

$ rg -c '' --glob '*.go' internal/appengine/ | sort -t: -k2 -rn
internal/appengine/channels/ws_test.go:933
internal/appengine/service.go:846
internal/appengine/http.go:696
internal/appengine/data.go:514
internal/appengine/http_test.go:512
internal/appengine/formats_parity_test.go:467
internal/appengine/channels/ws.go:442
internal/appengine/e2e_test.go:351
internal/appengine/channels/room_test.go:332
internal/appengine/groups_parity_test.go:309
internal/appengine/channels/room.go:287
internal/appengine/mirror_test.go:243
internal/appengine/channels/protocol_test.go:214
internal/appengine/query_opts_test.go:177
internal/appengine/stream/ws.go:162
internal/appengine/dashboard_compat_test.go:160
internal/appengine/channels/protocol.go:128
internal/appengine/patch_test.go:113
internal/appengine/stream/ws_test.go:110
internal/appengine/downsample_test.go:93
internal/appengine/links_test.go:65
internal/appengine/writeerror_test.go:55
internal/appengine/downsample.go:29

The 09-07 appengine review covered only the main package files (service/http/data/downsample). The `channels/` subpackage (ws.go, room.go, protocol.go) and `stream/ws.go` have **never** been reviewed — websocket/room concurrency code is exactly where a code review earns its keep. I'll go deep there.

Let me read the files, largest first.
→ Read internal/appengine/channels/ws.go
→ Read internal/appengine/channels/room.go
→ Read internal/appengine/channels/protocol.go
Now let me read the tests and the paired `stream/ws.go`, plus the Bus → stream semantics (buffer 0 / close behaviour).
→ Read internal/appengine/channels/room_test.go
→ Read internal/appengine/stream/ws.go
Let me check what the ws tests cover and how the engine stream bus cancel actually behaves (whether cancelling a subscription closes the returned channel).

$ rg -n "watchAuthPath|AuthorizesChannel|VerbWatch|VerbJoin|writeTimeout|teardown|rejoin|phx_join" internal/appengine/channels/*_test.go | head -60
internal/appengine/channels/protocol_test.go:22:	want := `["1","1","rooms:test:dashboard_abc_42","phx_join",{}]`
internal/appengine/channels/protocol_test.go:114:		{"join", `["1","1","rooms:test:dashboard_abc_42","phx_join",{}]`},
internal/appengine/channels/protocol_test.go:182:		{"four_elements_well_typed", `["1","1","rooms:test:x","phx_join"]`},
internal/appengine/channels/protocol_test.go:183:		{"six_elements_well_typed", `["1","1","rooms:test:x","phx_join",{},"extra"]`},
internal/appengine/channels/protocol_test.go:184:		{"non_string_ref", `["1",42,"rooms:test:x","phx_join",{}]`},
internal/appengine/channels/protocol_test.go:198:	input := `["1","1","rooms:test:x","phx_join"]`
internal/appengine/channels/ws_test.go:201:	if err := conn.Write(ctx, websocket.MessageText, []byte(`["1","1","`+topic+`","phx_join",{}]`)); err != nil {
internal/appengine/channels/ws_test.go:254:			msg := `["1","1","` + c.topic + `","phx_join",{}]`
internal/appengine/channels/ws_test.go:299:	if err := conn.Write(ctx, websocket.MessageText, []byte(`["1","1","`+topic+`","phx_join",{}]`)); err != nil {
internal/appengine/channels/ws_test.go:381:	msg := `["` + ref + `","` + ref + `","` + topic + `","phx_join",{}]`
internal/appengine/channels/ws_test.go:646:	if err := conn.Write(ctx, websocket.MessageText, []byte(`["1","1","`+topic+`","phx_join",{}]`)); err != nil {
internal/appengine/channels/ws_test.go:738:				[]byte(`["1","1","`+topic+`","phx_join",{}]`)); err != nil {
internal/appengine/channels/ws_test.go:859:			got, err := watchAuthPath(req)
internal/appengine/channels/ws_test.go:866:					t.Errorf("watchAuthPath() error = %v, want %v (got path %q)", err, want, got)
internal/appengine/channels/ws_test.go:871:				t.Fatalf("watchAuthPath() unexpected error: %v", err)
internal/appengine/channels/ws_test.go:874:				t.Errorf("watchAuthPath() = %q, want %q", got, tc.want)

$ ls internal/engine/stream/; echo ---; rg -c '' internal/engine/stream/*.go | sort -t: -k2 -rn
bus.go
bus_test.go
---
internal/engine/stream/bus_test.go:234
internal/engine/stream/bus.go:209

→ Read internal/engine/stream/bus.go
→ Read internal/appengine/channels/ws_test.go [offset=800, limit=133]
Let me check the remaining test coverage (rejoin/joinRef, group watch end-to-end) and whether the triggers package knows anything about group membership.

$ rg -n "rejoin|joinRef|join_ref|GroupName|group_name|zombie|deadroom|retired" internal/appengine/channels/ws_test.go internal/appengine/channels/room_test.go | head -40
internal/appengine/channels/ws_test.go:156:	// Check wire format: join_ref (element 0) must be null, not "".
internal/appengine/channels/ws_test.go:165:		t.Errorf("join_ref = %s, want null", string(arr[0]))
internal/appengine/channels/ws_test.go:593:	// Element 0: join_ref must be "1", not null.
internal/appengine/channels/ws_test.go:595:		t.Errorf("join_ref = %s, want \"1\"", string(arr[0]))
internal/appengine/channels/ws_test.go:630:// joinRef concurrently. The specced tests all push events with the room
internal/appengine/channels/ws_test.go:830:			name: "group data trigger, group_name top level (measured 2026-08-22)",
internal/appengine/channels/ws_test.go:831:			req:  `{"group_name":"g1","simple_trigger":{"type":"data_trigger","interface_name":"org.example.V1","match_path":"/value"}}`,
internal/appengine/channels/ws_test.go:836:			req:  `{"group_name":"g1","simple_trigger":{"type":"device_trigger","on":"device_error","device_id":"*"}}`,
internal/appengine/channels/ws_test.go:840:			name:    "nested-only group_name is refused like upstream's changeset",
internal/appengine/channels/ws_test.go:841:			req:     `{"simple_trigger":{"type":"data_trigger","group_name":"g1","interface_name":"org.example.V1","match_path":"/value"}}`,
internal/appengine/channels/ws_test.go:843:			wantOf:  errNestedGroupName,
internal/appengine/channels/ws_test.go:847:			req:     `{"group_name":"g1","simple_trigger":{"type":"device_trigger","on":"device_error","device_id":"` + dev + `"}}`,
internal/appengine/channels/room_test.go:327:			t.Fatalf("room not retired: Rooms() = %d, cancelCalled = %v", reg.Rooms(), bus.cancelCalled())

$ rg -n "GroupName|group" internal/engine/triggers/*.go | grep -vi "test" | head -30
internal/engine/triggers/match.go:10:// device_empty_cache_received, interface_minor_updated, and group-scoped
internal/engine/triggers/match.go:242:	GroupName          string          `json:"group_name"`
internal/engine/triggers/match.go:359:	deviceFilterErrs(fe, c.DeviceID, c.GroupName)
internal/engine/triggers/match.go:380:		evaluated: known && c.GroupName == "",
internal/engine/triggers/match.go:394:	t.noteUnsupported(c.GroupName != "", "group-scoped trigger (group_name)")
internal/engine/triggers/match.go:412:	deviceFilterErrs(fe, c.DeviceID, c.GroupName)
internal/engine/triggers/match.go:455:		evaluated: evaluated && c.GroupName == "",
internal/engine/triggers/match.go:466:	t.noteUnsupported(c.GroupName != "", "group-scoped trigger (group_name)")
internal/engine/triggers/match.go:484:// deviceFilterErrs validates the device_id / group_name pair into fe: at most
internal/engine/triggers/match.go:485:// one may be set (upstream validate_device_id_xor_group_name) and device_id
internal/engine/triggers/match.go:488:func deviceFilterErrs(fe fieldErrs, deviceID, groupName string) {
internal/engine/triggers/match.go:489:	if deviceID != "" && groupName != "" {
internal/engine/triggers/match.go:490:		fe.add("group_name", "must not be defined if device_id is defined")

→ Read internal/engine/triggers/match.go [offset=330, limit=175]
Let me examine the ws_test around rejoin/joinRef handling and any dead-room tests.
→ Read internal/appengine/channels/ws_test.go [limit=130, offset=560]
→ Read internal/appengine/channels/ws_test.go [limit=160, offset=200]

$ rg -n "group" docs/DESIGN.md | head -20; echo "---CHANNELS-NOTE---"; rg -n "channels|websocket|socket" docs/DESIGN.md | head -20
59:| **AppEngine API** | REST API to read device data, publish server-owned data, manage groups/aliases; Channels (WebSockets) | `internal/appengine` | Same `/v1/<realm>/devices/...` resource model; server-owned publishes are handed to `internal/broker` via the engine. |
258:-- Device groups (AppEngine)
259:CREATE TABLE groups (
265:CREATE TABLE group_devices (
266:    group_id  bigint NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
269:    PRIMARY KEY (group_id, device_id),
623:  unset), `devices-by-alias/...`, `groups...`. Envelope: `{"data": ...}`; errors:
---CHANNELS-NOTE---
22:| Message bus | RabbitMQ (inter-service AMQP) | In-process Go channels (sharded worker pools) |
35:criterion), and Astarte Channels' full Phoenix-socket protocol (we provide a simpler
50:through typed in-process interfaces and channels. The decoupling is preserved *at the package
63:| **Astarte Channels** | Phoenix WebSocket rooms for live data | `internal/appengine/stream` | Simplified WebSocket + SSE endpoint fed by the engine's fan-out bus. Not wire-compatible with Phoenix sockets in v1 (documented deviation). |
134:- Backpressure: bounded channels (default 4096/shard). When a shard saturates, the broker hook
603:1. **Astarte Channels** Phoenix-socket protocol is replaced by plain WebSocket/SSE (different
627:- Astrate-native additions under `/astrate/v1/...` (health, metrics, live stream socket) to
652:  (Housekeeping), `a_ch` (Channels → honoured by Astrate's stream socket). Each is a list of
658:- **`a_ch` is the exception and does not follow the rule above.** Upstream's Channels socket
662:  surfaces. Measured, not inferred — see `test/conformance/upstream/channels.json` and

Let me verify a couple of things: whether any test pins group-watch scoping, and whether DESIGN says anything about group scoping in channels.
→ Read internal/appengine/channels/ws_test.go [limit=110, offset=690]

$ rg -n "group" test/conformance/upstream/channels.json | head -30
160:      "name": "watch group data_trigger (group_name top level), claim WATCH::groups/\u003cgroup\u003e/\u003cinterface\u003e\u003cmatch_path\u003e",
161:      "why": "Astrate's predicted shape. If upstream builds the same authorization path for a group trigger this is accepted.",
164:        "WATCH::groups/probe/org.astrate.bench.Individual/value"
166:      "operation": "watch (group trigger)",
167:      "payload": "{\"group_name\":\"probe\",\"name\":\"probe\",\"simple_trigger\":{\"interface_major\":1,\"interface_name\":\"org.astrate.bench.Individual\",\"match_path\":\"/value\",\"on\":\"incoming_data\",\"type\":\"data_trigger\",\"value_match_operator\":\"*\"}}",
172:      "name": "watch group data_trigger, claim WATCH::groups/\u003cgroup\u003e/\u003cinterface\u003e (no match path)",
176:        "WATCH::groups/probe/org.astrate.bench.Individual"
178:      "operation": "watch (group trigger)",
179:      "payload": "{\"group_name\":\"probe\",\"name\":\"probe\",\"simple_trigger\":{\"interface_major\":1,\"interface_name\":\"org.astrate.bench.Individual\",\"match_path\":\"/value\",\"on\":\"incoming_data\",\"type\":\"data_trigger\",\"value_match_operator\":\"*\"}}",
184:      "name": "watch group data_trigger, claim WATCH::\u003cdevice\u003e/\u003cinterface\u003e\u003cmatch_path\u003e (the device shape)",
185:      "why": "The alternative reading: that a group trigger authorizes against the publishing device's path. Exactly one of this row and the first group row can be accepted, so together they identify which string upstream builds.",
190:      "operation": "watch (group trigger)",
191:      "payload": "{\"group_name\":\"probe\",\"name\":\"probe\",\"simple_trigger\":{\"interface_major\":1,\"interface_name\":\"org.astrate.bench.Individual\",\"match_path\":\"/value\",\"on\":\"incoming_data\",\"type\":\"data_trigger\",\"value_match_operator\":\"*\"}}",
196:      "name": "watch group data_trigger (group_name inside simple_trigger), claim WATCH::groups/\u003cgroup\u003e/\u003cinterface\u003e\u003cmatch_path\u003e",
197:      "why": "The other payload shape a client might send: group_name nested where Astrate reads it. If upstream refuses this one with a validation error rather than an authorization refusal, the field belongs at the top level only.",
200:        "WATCH::groups/probe/org.astrate.bench.Individual/value"
202:      "operation": "watch (group trigger)",
203:      "payload": "{\"name\":\"probe\",\"simple_trigger\":{\"group_name\":\"probe\",\"interface_major\":1,\"interface_name\":\"org.astrate.bench.Individual\",\"match_path\":\"/value\",\"on\":\"incoming_data\",\"type\":\"data_trigger\",\"value_match_operator\":\"*\"}}",
204:      "reply": "[\"1\",\"2\",\"rooms:bench:probe\",\"phx_reply\",{\"status\":\"error\",\"response\":{\"errors\":{\"group_name\":[\"must be present if device_id is not set\"],\"device_id\":[\"must be present if group_name is not set\"]}}}]",
208:      "name": "watch group device_trigger (device_id \"*\" inside simple_trigger), claim WATCH::groups/\u003cgroup\u003e",
209:      "why": "Astrate's predicted shape for a group device trigger: the bare group name. The payload needs device_id=\"*\" inside simple_trigger — a concrete id is refused with reason \"device_id must be * for group triggers\", the mirror image of plain device triggers where \"*\" is refused.",
212:        "WATCH::groups/probe"
214:      "operation": "watch (group trigger)",
215:      "payload": "{\"group_name\":\"probe\",\"name\":\"probe\",\"simple_trigger\":{\"device_id\":\"*\",\"on\":\"device_error\",\"type\":\"device_trigger\"}}",
220:      "name": "watch group device_trigger, claim WATCH::groups/\u003cother group\u003e",
221:      "why": "The paired refusal: same verb, wrong group name, so the group name is really being matched.",
224:        "WATCH::groups/other"
226:      "operation": "watch (group trigger)",
227:      "payload": "{\"group_name\":\"probe\",\"name\":\"probe\",\"simple_trigger\":{\"device_id\":\"*\",\"on\":\"device_error\",\"type\":\"device_trigger\"}}",

Let me confirm the group-watch semantics against the conformance recording before proposing anything.
$ rg -n "group" docs/DESIGN.md | sed -n '25,60p'
(no output)

I now have the full picture. Let me write the review evidence file and append task lines.
← Write .mule/reviews/channels-2026-09-18.md
Wrote file successfully.

→ Read .mule/todo.md [limit=15]
→ Read .mule/todo.md [limit=20, offset=174]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -189,5 +189,8 @@
 
 - [x] docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] httpx-cors-vary-origin-passthrough [auto]: in internal/httpx/cors.go non-wildcard mode, the disallowed-origin pass-through (cors.go:26-29) returns without `Vary: Origin`, so responses whose CORS outcome depends on the Origin header reach shared caches without the RFC 9110 §12.5.5 variation marker (a cache could serve the CORS-approving variant to a disallowed origin, or the headerless variant to an allowed one); the stamp lives only in the allowed branch (cors.go:35-36) and only the preflight pins it (cors_test.go:42). Add `h.Add("Vary", "Origin")` on the pass-through branch too, and assert `Vary: Origin` in TestCORSActualRequestStamped (actual non-preflight, pinned only for Allow-Origin today) and TestCORSDisallowedOriginPassesThroughUnstamped. Additive header change; container-free.
 - [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] channels-group-watch-membership-scope: group watches authorize against `WATCH::groups/<name>/…` (ws.go:178-205, conformance channels.json) but `Room.dispatch` ever only filters by `watchEntry.deviceID` (empty for a group watch — room.go:36-40, 196-200) and the trigger compiles with no group/device scope (triggers/match.go:374-396, because upstream carries group_name at the payload top level), so a group-scoped token receives every realm event matching interface/path/on, members or not. Plumb a group-membership resolver into internal/appengine/channels (it knows only `Bus` today), resolve membership at Watch time and store it in `watchEntry`, filter deliveries by it. Add a container-free test with a fake resolver: watch group `probe`, publish an event from a member and a non-member, assert only the member's arrives. Verify the member-filter rule against upstream on the Legion if disputed; the boundary leak stands regardless.
+- [ ] channels-rejoin-joinref-tagging-test: pin in internal/appengine/channels/ws_test.go that a second `phx_join` on the same topic updates the session's `joinRef` (ws.go:284-291) and subsequent `new_event` pushes carry the fresh ref (read per event at ws.go:403) — no test drives a rejoin today (all `TestWireSession_*` join once), so the comment's "server pushes must be tagged with the current one" rots unasserted. Wire: join with ref "1", rejoin with "2", push an event, assert element 0 is `"2"`.
+- [ ] channels-rejoin-authz-mismatch: `handleJoin`'s comment says upstream authorizes every join including a rejoin (ws.go:276-277), but the rejoin branch (ws.go:284-291) replies OK without calling `s.tok.AuthorizesChannel(auth.VerbJoin, name)` — behaviourally identical today (a token is immutable for the session) but the comment lies about what the code does. Either move the authz check to cover the rejoin branch too, or fix the comment; mechanical, no behaviour change, choose one.
 - [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a locally-run `elixir -e` with `Base.url_decode64!`) with the device ID `"AAAAAAAAAAAAAAAAAAAAAB"` — nonzero unused bits in the 22nd base64url char — and report whether upstream accepts it: `pkg/deviceid/deviceid.go:34-37` uses `base64.RawURLEncoding.Strict()` and deviceid_test.go:176 asserts rejection while the comment claims parity with `Elixir Base.url_decode64!(padding: false)`, which decodes and discards those bits rather than erroring. Astrate may therefore be stricter than upstream on the same wire form (a device registered with such an id on upstream would be rejected here). If upstream rejects it identically, the strict encoding is verified parity and the task is done; if upstream accepts it, escalate the relax-vs-strict call to `.mule/for-giulio.md` with the measurement. Probe first, no code change either way.

$ git status --short; git diff --stat
M .mule/todo.md
?? .mule/reviews/channels-2026-09-18.md
 .mule/todo.md | 3 +++
 1 file changed, 3 insertions(+)

Done: reviewed internal/appengine/channels + stream (never before reviewed — the 09-07 appengine run only read service/http/data/downsample), wrote evidence file, appended 3 task lines. No code, no git operations beyond status check.

Files: `.mule/reviews/channels-2026-09-18.md` (new), `.mule/todo.md` (+3 lines)
Verified: `git status --short` -> clean tree except the two proposal files
Unsure: nothing
Follow-ups: `channels-group-watch-membership-scope`, `channels-rejoin-joinref-tagging-test`, `channels-rejoin-authz-mismatch` (appended; the group-task notes the Legion probe if the member-rule is disputed)

Key finding: group watches authorize against `WATCH::groups/<name>/…` but `Room.dispatch` never filters by group membership — a group-scoped token receives realm-wide events for the matching interface.
