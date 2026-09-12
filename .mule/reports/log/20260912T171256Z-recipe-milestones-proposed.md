slug: recipe-milestones
verdict: proposed
at:  db98e58
ran: 2026-09-12T17:12:56Z on DietPi in 214s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md

$ gh issue list --label milestone-v3.0 --state all --limit 50 2>&1; echo '---'; gh issue list --label upstream-parity --state all --limit 50 2>&1
---
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

✱ Glob ".mule/research/**" 0 matches

$ gh label list --limit 100 2>&1 | head -60
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

✱ Glob "**/upstream-parity*.md" 0 matches

$ ls -la .mule/ .mule/research/ 2>&1; echo '---'; ls .mule/recipes/ 2>&1
ls: cannot access '.mule/research/': No such file or directory
.mule/:
total 3048
drwxr-xr-x  7 root root    4096 Sep 12 19:08 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      13 Sep 12 19:09 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 12 18:45 .heartbeat
-rw-r--r--  1 root root    8192 Sep 12 19:10 .last-output
-rw-r--r--  1 root root       0 Sep 12 19:09 .lock
-rw-r--r--  1 root root      11 Sep 12 19:09 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 2941065 Sep 12 19:10 cron.log
drwxr-xr-x  2 root root    4096 Sep 11 20:29 failed
-rw-r--r--  1 root root   28479 Sep 12 13:04 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    5266 Sep 12 18:45 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 11 21:40 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   43263 Sep 12 18:45 todo.md
---
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

$ ls .mule/tasks/ .mule/reports/ 2>&1; echo '---'; ls docs/handoff/ 2>&1 | head -40
.mule/reports/:
hygiene.md
log

.mule/tasks/:
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md
---
README.md
flow-design-a-named-durable-flows-2026-07-29.md
flow-design-b-container-block-2026-07-29.md
flow-lifecycle-followups-design-2026-08-22.md
flow-parity-audit-2026-07-29.md
flow-v2-decisions-2026-07-29.md
phase-1-handoff.md
phase-1-memory.md
phase-10-handoff.md
phase-10-memory.md
phase-11-handoff.md
phase-2-handoff.md
phase-2-memory.md
phase-3-handoff.md
phase-3-memory.md
phase-4-handoff.md
phase-4-memory.md
phase-5-handoff.md
phase-5-memory.md
phase-6-handoff.md
phase-6-memory.md
phase-7-memory.md
phase-8-handoff.md
phase-8-memory.md
phase-9-handoff.md
phase-9-memory.md
pokemon-agent-handoff.md
pokemon-agent-memory.md
session-2026-07-29-flow-factory-handoff.md
session-2026-07-29-handoff.md

$ find /root -iname '*upstream*parity*' -not -path '*/node_modules/*' 2>/dev/null; echo '---'; find /root -iname '*2026-08-22*' -not -path '*/node_modules/*' 2>/dev/null | head -20
---
/root/.local/share/logs/RustDesk/service/rustdesk_r2026-08-22_00-49-34.log
/root/astrate-mule/docs/handoff/flow-lifecycle-followups-design-2026-08-22.md
/root/g-mind/.agent/plans/2026-08-22 - locale slim Legion Go.md
/root/astrate-survey/.mule/research/survey-2026-08-22.md
/root/astrate-survey/.mule/research/issues-2026-08-22.md

✱ Grep "retention ceiling → alias/group surfaces|2026-08-22" 100 matches

$ git log --oneline -20 -- .mule/research/ 2>&1; echo '---'; find /root/astrate-mule -name '*.md' -path '*research*' 2>/dev/null; echo '---'; git log --all --oneline --diff-filter=D -- '*.md' 2>/dev/null | head
---
---
71e2cdd Move mule and trickle skills to system-wide agent-skills repo
656815c remove Clea Portal references
57338db mule: drop the report written in the old format

→ Read .mule/reports/log/20260911T185303Z-recipe-milestones-proposed.md
→ Read .mule/for-giulio.md
→ Read docs/UPSTREAM-EXPERIMENTAL.md
The milestone is v3.0. The recipe needs a state re-check (did upstream ship stable v1.4.0?), the standing UPSTREAM-EXPERIMENTAL reconcile, and a gap scan. Let me run the checks.

$ gh api repos/astarte-platform/astarte/releases --jq '.[].tag_name' 2>&1 | head -8; echo '=== open upstream-parity ==='; gh issue list --label upstream-parity --state open --limit 30 --json number,title,labels --jq '.[] | "#\(.number) [\([.labels[].name]|join(","))] \(.title)"'; echo '=== milestone-3.0 all ==='; gh issue list --label milestone-3.0 --state all --limit 50 --json number,title,state --jq '.[] | "#\(.number) \(.state) \(.title)"'; echo '=== mule-alarm recent ==='; gh issue list --label mule-alarm --state open --limit 10 --json number,createdAt --jq '.[] | "#\(.number) \(.createdAt)"'
v1.4.0-rc.5
v1.3.3
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
v1.3.2
v1.3.1
=== open upstream-parity ===
#93 [mule-review,upstream-parity] In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
#92 [enhancement,upstream-parity,upstream-experimental] keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds
#78 [enhancement,milestone-4.0,upstream-parity] FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
=== milestone-3.0 all ===
=== mule-alarm recent ===
#101 2026-09-12T10:58:09Z
#100 2026-09-11T10:55:09Z
#99 2026-09-10T11:08:07Z
#98 2026-09-09T11:05:08Z
#97 2026-09-08T11:19:07Z
#96 2026-09-07T11:15:07Z
#95 2026-09-06T11:17:07Z
#94 2026-09-05T11:00:06Z

State looks identical to the 2026-09-11 run. Let me verify the earlier-milestone check (v2.0 gaps), whether #93's commit landed, and that main hasn't moved.
$ echo '=== milestone-2.0 open ==='; gh issue list --label milestone-2.0 --state open --limit 30 --json number,title --jq '.[] | "#\(.number) \(.title)"'; echo '(end)'; echo '=== recent main commits ==='; git log --oneline -5 2>&1; echo '=== is 8c61268 on main? ==='; git merge-base --is-ancestor 8c61268 HEAD 2>&1 && echo merged || echo 'not on HEAD'; git branch -a --contains 8c61268 2>&1 | head -5
=== milestone-2.0 open ===
(end)
=== recent main commits ===
db98e58 mule: recipe hygiene timed out
fb9d8bf mule: log docs-sync-ae-post-groups-409
d6865a1 mule: docs-sync-ae-post-groups-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups` in docs/api/astarte_appengine_api.yaml — `Service.CreateGroup` returns `ErrGroupAlreadyExists` (internal/appengine/service.go:495, mapped at http.go:660-661) which the handler answers as `409` "Group already exists"; the spec (yaml:821-839) documents only 201/400/401/422/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
b1f5d98 mule: log docs-sync-ae-write-value-422
17f6f70 mule: docs-sync-ae-write-value-422 [auto]: add the missing `422` response to the six PUT/POST data-write ops in docs/api/astarte_appengine_api.yaml (device-scoped, by-alias, and in-group at `/interfaces/{interface}/{path}`). A value larger than `maxValueBytes` (64 KiB, internal/appengine/http.go:24) is answered `422` "Value size exceeds size limits" by putData/putDataByAlias/putDataInGroup (http.go:298-301, 327-329, 346-348), but the specs for PUT/POST (yaml:686-695, 726-736, 453-463, 499-509, 1176-1186, 1218-1228) list only 200/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
=== is 8c61268 on main? ===
merged
* mule/queue
  remotes/origin/mule/queue

$ gh issue list --label milestone-4.0 --state open --limit 20 --json number,title --jq '.[] | "#\(.number) \(.title)"'; echo '=== comments on #93 ==='; gh issue view 93 --json comments --jq '[.comments[].body] | length'
#78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
=== comments on #93 ===
1

$ echo "=== branch ==="; git branch --show-current; echo "=== origin/main contains 8c61268? ==="; git merge-base --is-ancestor 8c61268 origin/main 2>/dev/null && echo 'yes, on origin/main' || echo 'no'; echo "=== origin/main log -5 ==="; git log origin/main --oneline -5 2>&1
=== branch ===
mule/queue
=== origin/main contains 8c61268? ===
no
=== origin/main log -5 ===
2a033d4 mule: parking a decision on an upstream release now has something that notices
99743ca mule: the keyAgreement parking condition has fired — record the scope decision
ca47b35 mule: the hygiene sweep stops lying about deps, and the decision queue is a queue again
f1d0069 mule: the survey stops stranding reports, and the map stops lying about timers
215d409 mule: stop tracking the Pi's tick transcript

State confirmed. Nothing new since the 2026-09-11 run: upstream still at stable `v1.3.3` (newest overall `v1.4.0-rc.5`), `milestone-3.0` label empty (the v3.0 ledger lives under `upstream-parity`), the only open `upstream-parity` items are the already-escalated #92/#93/#78, and no stable v1.4.0 means the two UPSTREAM-EXPERIMENTAL rows (#67/#68) are not yet due for reconcile. No new gaps → no issues to file, and the 1.4 gate is still open so no "cut the tag" proposal. The deliverable is one `for-giulio.md` entry.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,10 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+
 - **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
 
 - **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pile is now 7 straight days and worth a look.** Twelve open issues. Still no machine-checkable fix candidates, so no task lines: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-10 run: **#100**, today's alarm (created 10:55Z, "nothing has landed in 15h"). #99–#94 are the same `mule-alarm` (one per day 2026-09-05→09-10, ~11:00Z each, zero comments) — telemetry, not code issues, so never proposable. But seven consecutive alarms is past the "low-activity window" wording earlier runs used: the mule has landed no commit since ~2026-09-04/05, and `main`'s queue copy shows the top line as `- [ ]` while the landed work below is all `[!] BLOCKED` (`wrote nothing` / `tests failed`) — though `mule/queue` is the authoritative copy. **Proposal: close #94–#99** — each was a one-day event, superseded by next day's alarm, never actionable; leave #100 (live today) to self-expire. If the idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
 

Done: v3.0 milestone re-verified — unchanged, still gated on upstream stable v1.4.0; one `for-giulio.md` line recorded, nothing filed.
Files: .mule/for-giulio.md
Verified: `gh api repos/astarte-platform/astarte/releases` (newest stable v1.3.3, newest release v1.4.0-rc.5) -> no stable v1.4.0; `gh issue list --label milestone-3.0` (0) and `--label upstream-parity --state open` (#92/#93/#78) + `--label milestone-2.0 --state open` (0) -> no change vs 2026-09-11; `git merge-base --is-ancestor 8c61268 origin/main` (no — #93's fix still on `mule/queue`) -> pass
Unsure: nothing (state identical to the 2026-09-11 run; no new gaps, so methodically same conclusion)
Follow-ups: none (no gaps to file; no "cut the tag" line — the 1.4 gate per #92/#67/#68/#90 is still open)
