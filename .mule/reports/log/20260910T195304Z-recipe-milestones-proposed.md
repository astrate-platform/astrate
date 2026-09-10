slug: recipe-milestones
verdict: proposed
at:  8e868f5
ran: 2026-09-10T19:53:04Z on DietPi in 167s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
v3.0 is the first non-DONE milestone. Scope is decided, so I proceed to investigation. Let me gather the source materials in parallel.
→ Read docs/UPSTREAM-EXPERIMENTAL.md
→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ gh issue list --label milestone-v3.0 --state all --limit 50 2>&1
(no output)

The v3.0 milestone is the target. Note the `milestone-v3.0` probe returned empty — prior runs used the `milestone-3.0` label convention. Let me check the actual issue state properly.
$ gh issue list --label upstream-parity --state all --limit 100 2>&1
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
91	CLOSED	Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health	mule-blocked, upstream-parity	2026-09-04T18:34:38Z
89	CLOSED	Dashboard flow-block schema mismatch (split_map/virtual pools hardcoded; null_sink/log_sink unknown)	enhancement, upstream-parity	2026-08-23T15:08:58Z
88	CLOSED	Flow auth: support a_f JWT claim	enhancement, upstream-parity	2026-08-23T15:08:59Z
87	CLOSED	Flow block: lua_map — needs embedded Lua runtime (parked)	enhancement, upstream-parity	2026-09-04T19:20:52Z
86	CLOSED	Flow: pipeline source DSL — keep DAG-JSON as documented deviation?	enhancement, upstream-parity	2026-08-23T11:57:02Z
85	CLOSED	Flow API: user-defined composite blocks	enhancement, upstream-parity	2026-08-23T21:47:32Z
84	CLOSED	Flow blocks: virtual_device_pool / dynamic_virtual_device_pool	enhancement, upstream-parity	2026-08-25T14:56:53Z
83	CLOSED	Flow blocks: mqtt_source/mqtt_sink (+modbus_tcp_source?) demand-driven	enhancement, upstream-parity	2026-08-23T15:09:03Z
82	CLOSED	Flow blocks: http_source/http_sink (demand-driven)	enhancement, upstream-parity	2026-08-23T15:09:04Z
81	CLOSED	Flow block: json_path_map	enhancement, upstream-parity	2026-08-23T15:09:06Z
80	CLOSED	Flow blocks: pure-transform set (to_json, update_metadata, split_map, random_source, sort)	enhancement, upstream-parity	2026-08-23T15:09:07Z
79	CLOSED	Verify registration-limit-reached HTTP status vs upstream	enhancement, upstream-parity	2026-08-24T13:21:57Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
77	CLOSED	Verify per-service version endpoints (GET /version, GET /v1/{realm}/version) served everywhere	enhancement, upstream-parity	2026-08-24T13:23:44Z
76	CLOSED	Housekeeping: GET /v1/realm-defaults/replication — decide reject/deviate (Cassandra-shaped)	enhancement, upstream-parity	2026-08-22T19:12:29Z
75	CLOSED	Housekeeping: decide realm-deletion gating/preconditions vs always-sync deviation	enhancement, upstream-parity	2026-08-23T15:09:10Z
74	CLOSED	Housekeeping: PATCH /v1/realms/{realm} (jwt key, registration limit, retention; null=unset)	enhancement, upstream-parity	2026-08-23T15:09:12Z
73	CLOSED	Housekeeping: default datastream retention injection env var (upstream 1.4)	enhancement, upstream-parity	2026-08-23T15:09:13Z
72	CLOSED	Realms: datastream_maximum_storage_retention ceiling (create/patch/enforce)	enhancement, upstream-parity	2026-08-23T15:09:15Z
71	CLOSED	Pairing: realm-scoped health check GET /v1/{realm}/health (upstream 1.3)	enhancement, upstream-parity	2026-08-23T15:09:16Z
70	CLOSED	Triggers: audit wildcard semantics (interface_name '*', match_path '/*' forcing rules)	enhancement, upstream-parity	2026-08-23T15:09:18Z
69	CLOSED	Verify unknown-realm HTTP status on RM endpoints against upstream	enhancement, upstream-parity	2026-08-24T13:23:42Z
68	CLOSED	Decide async_operation=false params vs documented always-sync deviation	enhancement, mule-blocked, upstream-parity, upstream-experimental	2026-09-04T18:34:35Z
67	CLOSED	Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4)	enhancement, upstream-parity, upstream-experimental	2026-09-04T18:34:19Z
66	CLOSED	Realm Management: detailed=true interface listing with full mappings (upstream 1.4)	enhancement, upstream-parity	2026-08-23T15:09:20Z
65	CLOSED	Policies: handler-overlap rejection + retry_times coupling + prefetch_count	enhancement, upstream-parity	2026-08-23T15:09:22Z
64	CLOSED	Triggers: decide AMQP action behavior (validate-reject vs NATS-forward deviation)	enhancement, upstream-parity	2026-08-23T15:09:23Z
63	CLOSED	Triggers: HTTP action validation limits (URL/method/header blocklist/template size)	enhancement, upstream-parity	2026-08-23T15:09:25Z
62	CLOSED	Realm Management: audit install/update/delete error codes and statuses	enhancement, upstream-parity	2026-08-23T15:09:26Z
61	CLOSED	Realm Management: audit interface/mapping validation matrix against astarte_core	enhancement, upstream-parity	2026-08-23T15:09:28Z
60	CLOSED	Realm Management: GET config/datastream_maximum_storage_retention (since upstream 1.2.0)	enhancement, upstream-parity	2026-08-22T03:23:47Z
59	CLOSED	AppEngine: group create-body validation + UUID-v1 from_token for group device listing	enhancement, upstream-parity	2026-08-22T10:35:18Z
58	CLOSED	AppEngine: PATCH requires Content-Type application/merge-patch+json	enhancement, upstream-parity	2026-08-22T03:15:59Z
57	CLOSED	AppEngine: audit server-write error taxonomy against upstream	enhancement, upstream-parity	2026-08-24T13:24:06Z
56	CLOSED	AppEngine: object-aggregate response formats table and disjoint_tables	enhancement, upstream-parity	2026-08-23T15:09:31Z
55	CLOSED	AppEngine: complete data-query parameter set	enhancement, upstream-parity	2026-08-23T15:09:33Z
54	CLOSED	AppEngine: GET /groups/{g} show-group endpoint	enhancement, upstream-parity	2026-08-22T10:35:16Z
53	CLOSED	AppEngine: group-scoped device PATCH + interface data access under groups/{g}/devices/{d}	enhancement, upstream-parity	2026-08-22T09:18:23Z
52	CLOSED	AppEngine: devices-by-alias full mirror (PATCH via alias + interface data access)	enhancement, upstream-parity	2026-08-22T09:18:20Z
51	CLOSED	keyAgreement control topic + encrypted mappings (upstream 1.4 experimental) — investigate/park	enhancement, upstream-parity	2026-08-22T08:48:28Z
50	CLOSED	Broker/engine hardening: discard malformed introspection/purge/capabilities instead of crashing	enhancement, upstream-parity	2026-08-22T04:04:24Z
49	CLOSED	MQTT v1: binaryblob properties re-send as BSON subtype 0 in consumer/properties	enhancement, upstream-parity	2026-08-22T04:04:23Z
48	CLOSED	MQTT v1: tolerate devices connecting with an empty introspection	enhancement, upstream-parity	2026-08-22T04:04:11Z
47	CLOSED	MQTT v1: accept device capabilities incl. purge_properties_compression_format (plaintext purge)	enhancement, upstream-parity	2026-08-22T04:04:09Z

$ gh issue list --label milestone-3.0 --state all --limit 50 2>&1; echo "---LABELS---"; gh label list --limit 100 2>&1
---LABELS---
bug	Something isn't working	#d73a4a
documentation	Improvements or additions to documentation	#0075ca
duplicate	This issue or pull request already exists	#cfd3d7
enhancement	New feature or request	#a2eeef
good first issue	Good for newcomers	#7057ff
help wanted	Extra attention is needed	#008672
invalid	This doesn't seem right	#e4e669
question	Further information is requested	#d876e3
wontfix	This will not be worked on	#ffffff
mule	Queued for the autonomous mule to pick up	#5319e7
mule-alarm	Filed by the mule's dead-man switch: it has gone quiet	#b60205
mule-review	Mule pushed a change for this; waiting for a human or a strong model to review it	#0e8a16
mule-blocked	The mule tried this and could not do it; see its comment	#d93f0b
legion	Needs the Legion Go (Docker, race detector, benchmarks)	#fbca04
readonly	A standing check: verifies something, changes nothing, never closes	#c5def5
milestone-2.0	astarte-flow feature parity milestone	#0E8A16
milestone-4.0		#5319e7
milestone-3.0		#5319e7
upstream-parity	Feature parity with upstream astarte-platform	#0052cc
upstream-experimental	Adopted from an upstream experimental/unstable feature — review each upstream release	#d4c5f9

$ gh issue list --label upstream-experimental --state all --limit 50 2>&1; echo "---OPEN MULE---"; gh issue list --label mule --state open --limit 50 2>&1
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
68	CLOSED	Decide async_operation=false params vs documented always-sync deviation	enhancement, mule-blocked, upstream-parity, upstream-experimental	2026-09-04T18:34:35Z
67	CLOSED	Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4)	enhancement, upstream-parity, upstream-experimental	2026-09-04T18:34:19Z
---OPEN MULE---

State is clear: `milestone-3.0` label is empty (tracking lives under `upstream-parity`), all #47–#89 gaps are closed except #92/#93/#78 which are already escalated, and UPSTREAM-EXPERIMENTAL rows #67/#68 wait on upstream v1.4 final. Two things to verify before concluding: issue #90's status and whether upstream has shipped v1.4.0 stable since the last run (2026-09-09).
$ gh issue view 90 --json number,title,state,labels,body 2>&1 | head -60
{"body":"Reminder of a pending decision parked during the upstream-parity investigation (#47–#89). **Do not raise the version now** — this is a milestone-4.0 item.\n\n## Context\n\n`realm.APICompatVersion` (`internal/realm/service.go`) is currently hardcoded to `\"1.2.2\"` — the emulated upstream API level Astrate declares.\n\nEdgehog 0.13.1 gates just-in-time provisioning triggers (`device_registered`, `device_deletion_started`) behind a Realm Management version >= 1.3. Astrate's engine already supports those trigger kinds end-to-end — verified in #28 (Edgehog compatibility investigation) — but Edgehog refuses to use them because of the declared version string.\n\n## What this issue must cover when picked up\n\n1. **Audit first**: what else does declaring \"1.3\" imply for other consumers?\n   - `astarte-dashboard` is currently validated against 1.2.2 behavior (see #13, error-name bundle comparison). Check whether any dashboard code path feature-gates on RM version besides the policies UI gate (>= 1.1.1).\n   - Sweep the parity findings (#47–#89): several are \"1.3/1.4 features\" we deliberately do NOT implement yet — does declaring 1.3 overpromise any of them?\n   - Check device SDKs and `astartectl` for version-gated behavior.\n2. **Then decide one of**:\n   - bump the constant to `\"1.3.0-rc.0\"` (upstream's own pre-release numbering for that line);\n   - or make it configuration (`internal/config`) so operators can declare the level their consumers expect.\n3. Update `docs/COMPATIBILITY.md` deviation #10 wording accordingly (propose, don't edit directly).\n\n## Outcome if done today (why it's parked)\n\nBumping without the audit risks silently promising surfaces we haven't verified; not bumping keeps Edgehog's JIT provisioning disabled even though the engine supports it. Milestone-4.0 scope decision.","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACtvELrg","name":"milestone-4.0","description":"","color":"5319e7"}],"number":90,"state":"CLOSED","title":"milestone-4.0: decide on bumping APICompatVersion from \"1.2.2\" to \"1.3.0-rc.0\" (or making it configurable)"}

$ gh api repos/astarte-platform/astarte/releases --paginate -q '.[] | "\(.tag_name)\t\(.published_at)\t\(.prerelease)"' 2>&1 | head -20
v1.4.0-rc.5	2026-08-20T12:28:04Z	true
v1.3.3	2026-08-07T14:33:58Z	false
v1.4.0-rc.4	2026-08-06T13:40:20Z	true
v1.4.0-rc.3	2026-07-31T10:48:39Z	true
v1.4.0-rc.2	2026-07-14T10:16:22Z	true
v1.4.0-rc.1	2026-07-14T08:53:27Z	true
v1.3.2	2026-07-14T10:07:26Z	false
v1.3.1	2026-07-13T11:50:05Z	false
v1.3.0	2026-05-06T15:30:50Z	false
v1.2.2	2026-04-27T16:39:28Z	false
v1.4.0-rc.0	2026-04-08T12:32:33Z	true
v1.3.0-rc.2	2026-04-08T16:34:59Z	true
v1.2.2-rc.0	2026-04-08T16:31:03Z	true
v1.2.1	2026-03-12T16:11:20Z	false
v1.2.1-rc.1	2026-02-13T10:36:20Z	true
v1.3.0-rc.1	2026-01-26T15:58:33Z	true
v1.3.0-rc.0	2025-11-21T14:58:46Z	true
v1.2.1-rc.0	2025-08-27T07:03:28Z	true
v1.2.1-alpha.0	2025-04-10T15:31:04Z	true
v1.2.0	2024-07-02T14:09:07Z	false

State fully verified. Everything is unchanged since the 2026-09-09 run: upstream still has no stable v1.4.0 (newest stable `v1.3.3` 2026-08-07; newest pre-release `v1.4.0-rc.5` 2026-08-20), `milestone-3.0` label empty, and the only open `upstream-parity` issues (#92, #93, #78) are all already escalated. #90 is closed with the frozen bump decision recorded. No new fileable gaps — the "looks complete, cut the tag" line stays inappropriate because the 1.4-final gate is unresolved. I'll record the run in for-giulio.md, matching prior milestone-run entries.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,10 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+
 - **Extension point: external-bus Intake (NATS JetStream), 2026-09-10.** `internal/broker/intake.go:56-62` carries a `TODO(extension point)` for a second `Intake` implementation backed by a durable external bus (NATS JetStream) for multi-instance deployment or restart survival. The frozen design decision is in `docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4`; the consumer must reproduce per-device ordering, deferred-ack backpressure, and QoS 0 drop semantics. Not tracked by any open issue. Your call: file a milestone issue or park until the roadmap calls for it.
 
 - **Extension point: timescaledb_toolkit lttb downsampling, 2026-09-10.** `internal/store/store.go:139-143` carries a `TODO(extension point)` for switching the `Downsample` method from `time_bucket+avg` to the toolkit's `lttb()` when `timescaledb_toolkit` is present — the probe at `store.go:144-150` already records availability in `s.hasToolkit`, but `datastreams.go` always uses the default path. Design reference: `docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5`. Not tracked by any open issue. Your call: file a milestone issue or park.
 

Done. Milestone v3.0 is unchanged and there is nothing to file — the run's deliverable is the for-giulio.md entry recording that.

Done: Milestone recipe v3.0 — verified all #47–#89 gaps delivered, remaining items (#92/#93/#78) already escalated, both UPSTREAM-EXPERIMENTAL rows gated on upstream v1.4 final (still rc.5); no new gaps, no issues filed, no task lines
Files: .mule/for-giulio.md
Verified: `gh issue list --label milestone-3.0 --state all` + `--label upstream-parity --state all` + `--label upstream-experimental` + `gh issue view 90` + `gh api repos/astarte-platform/astarte/releases` -> pass (state matches 2026-09-09 run)
Unsure: nothing — deliberately did not propose "complete, cut the tag" since v1.4.0 (part of the 1.3/1.4 target) has not shipped stable
Follow-ups: none
