slug: recipe-milestones
verdict: proposed
at:  c0c7959
ran: 2026-09-19T21:14:55Z on DietPi in 124s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
v3.0 is the active milestone. The recipe says its status is "First recipe job: triage #47–#89". Let me gather the inputs.

$ wc -l .mule/research/upstream-parity-2026-08-22.md docs/UPSTREAM-EXPERIMENTAL.md
25 docs/UPSTREAM-EXPERIMENTAL.md
  25 total
wc: .mule/research/upstream-parity-2026-08-22.md: No such file or directory

$ gh issue list --label milestone-3.0 --state all --limit 50
(no output)

$ gh issue list --label upstream-parity --state all --limit 100
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

The research doc isn't at that path. Let me find it and survey the current state.
$ ls .mule/research/ 2>/dev/null; ls .mule/
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

→ Read docs/UPSTREAM-EXPERIMENTAL.md
→ Read .mule/todo.md

$ git log --all --oneline -- .mule/research/ | head -20; echo ---; git log --all --diff-filter=D --name-only --oneline | grep -i 'upstream-parity' | head
916796b survey: triage the 2026-07-28..2026-09-04 backlog
6d8cae6 mule: daily cross-project survey, 2026-09-04
580faf1 mule: daily cross-project survey, 2026-09-02 (no material change)
8d44f59 mule: daily cross-project survey, 2026-09-01 (no material change)
1c618eb mule: daily cross-project survey, 2026-08-31
ce759e1 survey: recover 16 daily reports the broken commit check left staged (2026-07-28..2026-08-26), and take the tooling fix from main
f92d42a mule/research: commit the calibration survey, and mark it triaged
c6e0301 mule/research: seed the branch the daily survey timer will commit to
---

→ Read .mule/for-giulio.md

$ git log --all --oneline --name-only | grep -i 'upstream-parity-2026-08-22' | head; echo ---; git branch -a; echo ---; git log --oneline -15
---
  main
* mule/queue
  remotes/origin/HEAD -> origin/main
  remotes/origin/feat/pokemon-agent
  remotes/origin/fix/blocked-mule-issues
  remotes/origin/m12-06b-watch-path-onto-main
  remotes/origin/main
  remotes/origin/mule/queue
  remotes/origin/mule/research
  remotes/origin/wip/DietPi-20260727T171543Z
  remotes/origin/worktree-m12-05-match-upstream
  remotes/origin/worktree-m12-06-ach-join
  remotes/origin/worktree-m12-06b-watch-path
---
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

Now let me verify the current state against yesterday's run: check upstream releases and the open issue set.
$ gh api repos/astarte-platform/astarte/releases --jq '.[].tag_name' | sort -V | tail -8
v1.3.3
v1.3.4
v1.4.0-rc.0
v1.4.0-rc.1
v1.4.0-rc.2
v1.4.0-rc.3
v1.4.0-rc.4
v1.4.0-rc.5

$ gh issue list --label upstream-parity --state open --limit 50; echo ---; gh issue list --label milestone-3.0 --state all --limit 20; echo ---; gh issue list --label milestone-4.0 --state all --limit 20
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
---
---
90	CLOSED	milestone-4.0: decide on bumping APICompatVersion from "1.2.2" to "1.3.0-rc.0" (or making it configurable)	enhancement, milestone-4.0	2026-08-23T21:06:34Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z
28	CLOSED	milestone-4.0: investigate Edgehog compatibility with Astrate (readonly)	mule, readonly, milestone-4.0	2026-08-23T15:48:58Z

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | "\(.tag_name) \(.published_at) prerelease=\(.prerelease)"' | sort -V | tail -6
v1.4.0-rc.0 2026-04-08T12:32:33Z prerelease=true
v1.4.0-rc.1 2026-07-14T08:53:27Z prerelease=true
v1.4.0-rc.2 2026-07-14T10:16:22Z prerelease=true
v1.4.0-rc.3 2026-07-31T10:48:39Z prerelease=true
v1.4.0-rc.4 2026-08-06T13:40:20Z prerelease=true
v1.4.0-rc.5 2026-08-20T12:28:04Z prerelease=true

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.prerelease==false) | "\(.tag_name) \(.published_at)"' | sort -V | tail -4
v1.3.1 2026-07-13T11:50:05Z
v1.3.2 2026-07-14T10:07:26Z
v1.3.3 2026-08-07T14:33:58Z
v1.3.4 2026-09-18T09:42:57Z

Nothing shipped since yesterday's runs (stable line still tops out at v1.3.4, 2026-09-18; v1.4.0 still rc.5-only), and the open-issue set is identical: #92 (parked on stable v1.4.0), #93 (mule-review, own path), #78 (milestone-4.0). No new gaps, no fileable issue, nothing to queue. The deliverable is the verification record in `.mule/for-giulio.md`.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+
+---
+
 - **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
 
 ---
 

→ Read .mule/for-giulio.md [offset=9, limit=10]
Milestone recipe carried out for the active milestone (v3.0). Current target is v3.0 (upstream 1.2.2 → 1.3/1.4), scope decided, so I investigated: release sweep shows no stable v1.4.0 (newest stable v1.3.4, maintenance-only, no wire change; newest overall v1.4.0-rc.5) — nothing changed since yesterday's run. `milestone-3.0` label empty (ledger lives under `upstream-parity`); all #47–#89 closed except the already-escalated #92 (parked on stable v1.4.0), #93 (mule-review), #78 (milestone-4.0). UPSTREAM-EXPERIMENTAL rows #67/#68 still wait on v1.4 final; `APICompatVersion` stays 1.2.2 per #90. No new gaps, no fileable issues, no task lines; deliberately did not propose "complete, cut the tag" since the 1.4 gate is unresolved. Recorded the verification in `.mule/for-giulio.md`.

Done: filed the v3.0 verification record in for-giulio.md; no gaps found, no issues filed, no task lines proposed
Files: .mule/for-giulio.md
Verified: gh api releases sweep + gh issue list (upstream-parity / milestone-3.0 / milestone-4.0) -> pass (state unchanged)
Unsure: nothing
Follow-ups: none
