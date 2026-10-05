slug: recipe-milestones
verdict: proposed
at:  c03cfe4
ran: 2026-10-05T19:34:52Z on DietPi in 210s


> build · big-pickle

→ Read .mule/MULE.md
I'll start by reading the required files.
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md

# Todos
[•] Read .mule/research/upstream-parity-2026-08-22.md source doc for v3.0
[ ] List existing milestone-v3.0 issues (state all) to avoid duplicates
[ ] Check what the prior run already filed/queued (todo.md, for-giulio.md)
[ ] Reconcile docs/UPSTREAM-EXPERIMENTAL.md rows against targeted upstream level
[ ] Bucket gaps: machine-checkable -> gh issue create (max 5), design -> for-giulio.md

✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md

$ gh issue list --label milestone-v3.0 --state all --limit 50 2>&1 | head -60
(no output)

$ ls -la .mule/ .mule/research/ 2>&1; echo "---LABELS---"; gh label list --limit 60 2>&1
ls: cannot access '.mule/research/': No such file or directory
.mule/:
total 12640
drwxr-xr-x  7 root root     4096 Oct  5 21:23 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       13 Oct  5 21:31 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  5 20:14 .heartbeat
-rw-r--r--  1 root root        0 Oct  5 21:31 .last-output
-rw-r--r--  1 root root        0 Oct  5 21:32 .lock
-rw-r--r--  1 root root       11 Oct  5 21:31 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 12556790 Oct  5 21:32 cron.log
drwxr-xr-x  2 root root     4096 Oct  4 22:35 failed
-rw-r--r--  1 root root   127766 Oct  4 20:31 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    14600 Oct  5 20:14 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  4 21:08 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   149912 Oct  5 20:14 todo.md
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

$ gh issue list --state all --limit 100 --json number,title,state,labels 2>&1 | head -c 4000
[{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":118,"state":"OPEN","title":"mule: nothing has landed in 14h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":117,"state":"OPEN","title":"mule: nothing has landed in 15h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":116,"state":"OPEN","title":"mule: nothing has landed in 14h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":115,"state":"OPEN","title":"mule: nothing has landed in 14h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":114,"state":"OPEN","title":"mule: nothing has landed in 8h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":113,"state":"OPEN","title":"mule: nothing has landed in 14h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":112,"state":"OPEN","title":"mule: nothing has landed in 14h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":111,"state":"OPEN","title":"mule: nothing has landed in 16h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":110,"state":"OPEN","title":"mule: nothing has landed in 14h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":109,"state":"OPEN","title":"mule: nothing has landed in 14h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":108,"state":"OPEN","title":"mule: nothing has landed in 15h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":107,"state":"OPEN","title":"mule: nothing has landed in 14h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":106,"state":"OPEN","title":"mule: nothing has landed in 14h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":105,"state":"OPEN","title":"mule: nothing has landed in 14h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":104,"state":"OPEN","title":"mule: nothing has landed in 14h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":103,"state":"OPEN","title":"mule: nothing has landed in 16h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"number":102,"state":"OPEN","title":"mule: nothing has landed in 14h"},{"labels":[{"id":"LA_kwDORmfsJs8AAAACte7BRA","name":"mule-alarm","description":"Filed by the mule's dead-man switch: it has gone quiet","color":"b60205"}],"nu

$ ls -d .mule/research 2>&1; git log --oneline --all -- '.mule/research/*' 2>&1 | head -20; echo "---gitignore---"; cat .mule/.gitignore 2>&1
ls: cannot access '.mule/research': No such file or directory
916796b survey: triage the 2026-07-28..2026-09-04 backlog
6d8cae6 mule: daily cross-project survey, 2026-09-04
580faf1 mule: daily cross-project survey, 2026-09-02 (no material change)
8d44f59 mule: daily cross-project survey, 2026-09-01 (no material change)
1c618eb mule: daily cross-project survey, 2026-08-31
ce759e1 survey: recover 16 daily reports the broken commit check left staged (2026-07-28..2026-08-26), and take the tooling fix from main
f92d42a mule/research: commit the calibration survey, and mark it triaged
c6e0301 mule/research: seed the branch the daily survey timer will commit to
---gitignore---
# The mule's own working files. None of these belong in git, and — more sharply — any of
# them that is not listed here makes the tree dirty, which makes `tick` refuse to run and so
# turns the timer into a permanent silent no-op.
failed/
task.md
.timeout
.lock
.budget
.bench-linux-amd64
.last-output
.rotation
.heartbeat
.alarmed
.impl.patch
hosts
research/

$ gh issue list --label milestone-3.0 --state all --limit 50 --json number,title,state 2>&1 | head -c 4000; echo; echo "---count---"; gh issue list --label milestone-3.0 --state all --limit 50 --json number 2>&1 | python3 -c "import json,sys; d=json.load(sys.stdin); print(len(d))"
[]

---count---
0

$ ls docs/ 2>&1; echo "---UPSTREAM-EXPERIMENTAL---"; ls -la docs/UPSTREAM-EXPERIMENTAL.md 2>&1
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
UPSTREAM-EXPERIMENTAL.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui
---UPSTREAM-EXPERIMENTAL---
-rw-r--r-- 1 root root 1336 Sep  4 20:46 docs/UPSTREAM-EXPERIMENTAL.md

$ gh issue list --label upstream-parity --state all --limit 60 --json number,title,state,labels 2>&1 | python3 -c "
import json,sys
d=json.load(sys.stdin)
for i in d:
    labs=[l['name'] for l in i['labels']]
    print(i['number'], i['state'], '|', i['title'], '|', ','.join(labs))
print('total', len(d))
"
93 OPEN | In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged | mule-review,upstream-parity
92 OPEN | keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds | enhancement,upstream-parity,upstream-experimental
91 CLOSED | Pairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health | mule-blocked,upstream-parity
89 CLOSED | Dashboard flow-block schema mismatch (split_map/virtual pools hardcoded; null_sink/log_sink unknown) | enhancement,upstream-parity
88 CLOSED | Flow auth: support a_f JWT claim | enhancement,upstream-parity
87 CLOSED | Flow block: lua_map — needs embedded Lua runtime (parked) | enhancement,upstream-parity
86 CLOSED | Flow: pipeline source DSL — keep DAG-JSON as documented deviation? | enhancement,upstream-parity
85 CLOSED | Flow API: user-defined composite blocks | enhancement,upstream-parity
84 CLOSED | Flow blocks: virtual_device_pool / dynamic_virtual_device_pool | enhancement,upstream-parity
83 CLOSED | Flow blocks: mqtt_source/mqtt_sink (+modbus_tcp_source?) demand-driven | enhancement,upstream-parity
82 CLOSED | Flow blocks: http_source/http_sink (demand-driven) | enhancement,upstream-parity
81 CLOSED | Flow block: json_path_map | enhancement,upstream-parity
80 CLOSED | Flow blocks: pure-transform set (to_json, update_metadata, split_map, random_source, sort) | enhancement,upstream-parity
79 CLOSED | Verify registration-limit-reached HTTP status vs upstream | enhancement,upstream-parity
78 OPEN | FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate) | enhancement,milestone-4.0,upstream-parity
77 CLOSED | Verify per-service version endpoints (GET /version, GET /v1/{realm}/version) served everywhere | enhancement,upstream-parity
76 CLOSED | Housekeeping: GET /v1/realm-defaults/replication — decide reject/deviate (Cassandra-shaped) | enhancement,upstream-parity
75 CLOSED | Housekeeping: decide realm-deletion gating/preconditions vs always-sync deviation | enhancement,upstream-parity
74 CLOSED | Housekeeping: PATCH /v1/realms/{realm} (jwt key, registration limit, retention; null=unset) | enhancement,upstream-parity
73 CLOSED | Housekeeping: default datastream retention injection env var (upstream 1.4) | enhancement,upstream-parity
72 CLOSED | Realms: datastream_maximum_storage_retention ceiling (create/patch/enforce) | enhancement,upstream-parity
71 CLOSED | Pairing: realm-scoped health check GET /v1/{realm}/health (upstream 1.3) | enhancement,upstream-parity
70 CLOSED | Triggers: audit wildcard semantics (interface_name '*', match_path '/*' forcing rules) | enhancement,upstream-parity
69 CLOSED | Verify unknown-realm HTTP status on RM endpoints against upstream | enhancement,upstream-parity
68 CLOSED | Decide async_operation=false params vs documented always-sync deviation | enhancement,mule-blocked,upstream-parity,upstream-experimental
67 CLOSED | Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4) | enhancement,upstream-parity,upstream-experimental
66 CLOSED | Realm Management: detailed=true interface listing with full mappings (upstream 1.4) | enhancement,upstream-parity
65 CLOSED | Policies: handler-overlap rejection + retry_times coupling + prefetch_count | enhancement,upstream-parity
64 CLOSED | Triggers: decide AMQP action behavior (validate-reject vs NATS-forward deviation) | enhancement,upstream-parity
63 CLOSED | Triggers: HTTP action validation limits (URL/method/header blocklist/template size) | enhancement,upstream-parity
62 CLOSED | Realm Management: audit install/update/delete error codes and statuses | enhancement,upstream-parity
61 CLOSED | Realm Management: audit interface/mapping validation matrix against astarte_core | enhancement,upstream-parity
60 CLOSED | Realm Management: GET config/datastream_maximum_storage_retention (since upstream 1.2.0) | enhancement,upstream-parity
59 CLOSED | AppEngine: group create-body validation + UUID-v1 from_token for group device listing | enhancement,upstream-parity
58 CLOSED | AppEngine: PATCH requires Content-Type application/merge-patch+json | enhancement,upstream-parity
57 CLOSED | AppEngine: audit server-write error taxonomy against upstream | enhancement,upstream-parity
56 CLOSED | AppEngine: object-aggregate response formats table and disjoint_tables | enhancement,upstream-parity
55 CLOSED | AppEngine: complete data-query parameter set | enhancement,upstream-parity
54 CLOSED | AppEngine: GET /groups/{g} show-group endpoint | enhancement,upstream-parity
53 CLOSED | AppEngine: group-scoped device PATCH + interface data access under groups/{g}/devices/{d} | enhancement,upstream-parity
52 CLOSED | AppEngine: devices-by-alias full mirror (PATCH via alias + interface data access) | enhancement,upstream-parity
51 CLOSED | keyAgreement control topic + encrypted mappings (upstream 1.4 experimental) — investigate/park | enhancement,upstream-parity
50 CLOSED | Broker/engine hardening: discard malformed introspection/purge/capabilities instead of crashing | enhancement,upstream-parity
49 CLOSED | MQTT v1: binaryblob properties re-send as BSON subtype 0 in consumer/properties | enhancement,upstream-parity
48 CLOSED | MQTT v1: tolerate devices connecting with an empty introspection | enhancement,upstream-parity
47 CLOSED | MQTT v1: accept device capabilities incl. purge_properties_compression_format (plaintext purge) | enhancement,upstream-parity
total 46

→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ rg -n "milestone-3\.0|milestone-v3\.0" .mule/todo.md .mule/for-giulio.md 2>&1 | tail -40
.mule/for-giulio.md:27:  **Gate re-verified, unchanged.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — nothing shipped upstream since the 2026-10-03 run, so #92's parking condition (your 2026-09-04 decision) is still unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** (this milestone's ledger lives under `upstream-parity`); of #47–#89 **all remain closed**; open non-alarm set unchanged — **#93** (`mule-review`, own recipe path), **#92** (parked), **#78** (`milestone-4.0`), **#1** (untouched). The post-rc.6 master delta is now a closed set of **8 commits** (`gh api .../commits -f since=2026-09-30T07:28:27Z`), of which exactly two carry wire surface and **both already have owners in the queue**: `b6d46ad4` (#2237 unexpected object keys) → done as `appengine-unexpected-object-key`, and `958bb9fe` (#2232 `device_empty_cache_received`) → queued as `device-empty-cache-received-trigger` (currently `[!]` blocked). Neither is from stable 1.3 or 1.4, so neither is a v3.0 gate.
.mule/for-giulio.md:48:  **Verified unchanged, no new gaps.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition, frozen by you on 2026-09-04 ("wait for a stable v1.4.0", an `-rc.N` does not satisfy it), is still unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Of #47–#89 **all are closed**; the open non-alarm set is unchanged — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** (untouched per standing instruction). rc.6 adds no v3.0 gap: it is already audited in the two entries below, and I re-measured its non-FDO half myself (`compare/v1.4.0-rc.5...v1.4.0-rc.6`, 27 commits) — outside FDO the only library changes are a `# coveralls-ignore-stop` line in `apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex:43` and `realms/queries.ex` *removing* the ownership_voucher / to2_sessions / unconfirmed_devices table creation (moved into migrations). No Realm Management, AppEngine or Pairing wire surface moved.
.mule/for-giulio.md:104:- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Status: not started", which is no longer true.** Release sweep (`gh api repos/astarte-platform/astarte/releases`, run this session) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.5` (2026-08-20, prerelease) — byte-identical to what the 2026-09-24 run saw, so nothing shipped upstream in the two days since. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Open non-alarm issues remain exactly **#93** (aclhook comment, `mule-review`, its own recipe path), **#92** (keyAgreement, parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** (FDO, `milestone-4.0`), **#1** (untouched per standing instruction). No new gaps, **no issues filed, no `.mule/todo.md` task lines added** — there is nothing machine-checkable left, and filing anything here would duplicate #92/#93/#78.
.mule/for-giulio.md:179:- **Milestone recipe run, 2026-09-24: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-23 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:187:- **Milestone recipe run, 2026-09-23: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-21 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:203:- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:211:- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:223:- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:227:- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:243:- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) shows a new **stable v1.3.4** (2026-09-18) — maintenance only (data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change), so no new gap for Astrate — and still **no stable v1.4.0**: newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:255:- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-13 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:259:- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-12 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:263:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:271:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:273:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:304:  upstream gates.** `milestone-3.0` label remains empty. The three open `upstream-parity`
.mule/for-giulio.md:326:  `milestone-3.0` label is empty (this milestone's work lives under `upstream-parity`).** All

→ Read .mule/for-giulio.md [limit=60]

$ gh api repos/astarte-platform/astarte/releases --paginate 2>&1 | python3 -c "
import json,sys
d=json.load(sys.stdin)
rs=[(r['tag_name'], r['published_at'], r['prerelease']) for r in d]
rs.sort(key=lambda t: t[1], reverse=True)
for t in rs[:12]: print(t)
"
('v1.3.5', '2026-10-05T14:11:44Z', False)
('v1.4.0-rc.6', '2026-09-30T07:28:27Z', True)
('v1.3.4', '2026-09-18T09:42:57Z', False)
('v1.4.0-rc.5', '2026-08-20T12:28:04Z', True)
('v1.3.3', '2026-08-07T14:33:58Z', False)
('v1.4.0-rc.4', '2026-08-06T13:40:20Z', True)
('v1.4.0-rc.3', '2026-07-31T10:48:39Z', True)
('v1.4.0-rc.2', '2026-07-14T10:16:22Z', True)
('v1.3.2', '2026-07-14T10:07:26Z', False)
('v1.4.0-rc.1', '2026-07-14T08:53:27Z', True)
('v1.3.1', '2026-07-13T11:50:05Z', False)
('v1.3.0', '2026-05-06T15:30:50Z', False)

New upstream release today: **v1.3.5** (2026-10-05). Every prior milestone run saw v1.3.4 as newest. Diffing it now.
$ gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5 --jq '{ahead_by, total_commits, files: (.files|length), commits: [.commits[]|{sha: .sha[0:8], date: .commit.committer.date, msg: (.commit.message|split("\n")[0])}]}' 2>&1
{"ahead_by":4,"commits":[{"date":"2026-09-25T08:54:35Z","msg":"chore: update horde (#2224)","sha":"8d47574f"},{"date":"2026-10-05T09:02:23Z","msg":"chore: update horde","sha":"4c4222cf"},{"date":"2026-10-05T09:55:53Z","msg":"chore(rm): bump xandra to 0.19.4 (#2246)","sha":"c4c63ebe"},{"date":"2026-10-05T14:09:21Z","msg":"chore: prepare v1.3.5 (#2245)","sha":"8f0f5ef1"}],"files":35,"total_commits":4}

$ gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5 --jq '.files[] | "\(.status) +\(.additions)/-\(.deletions) \(.filename)"' 2>&1
modified +7/-7 .github/workflows/astarte-build-workflow.yaml
modified +6/-0 CHANGELOG.md
modified +2/-2 README.md
modified +2/-2 apps/astarte_appengine_api/mix.exs
modified +2/-2 apps/astarte_appengine_api/mix.lock
modified +2/-4 apps/astarte_data_updater_plant/mix.exs
modified +2/-2 apps/astarte_data_updater_plant/mix.lock
modified +2/-2 apps/astarte_housekeeping/mix.exs
modified +1/-1 apps/astarte_housekeeping/mix.lock
modified +2/-2 apps/astarte_pairing/mix.exs
modified +2/-2 apps/astarte_pairing/mix.lock
modified +2/-2 apps/astarte_realm_management/mix.exs
modified +3/-3 apps/astarte_realm_management/mix.lock
modified +2/-2 apps/astarte_trigger_engine/mix.exs
modified +1/-1 apps/astarte_trigger_engine/mix.lock
modified +1/-1 astarte-dashboard/package.json
modified +1/-1 doc/mix.exs
modified +8/-8 docker-compose.yml
modified +1/-1 libs/astarte_adapters/mix.exs
modified +2/-2 libs/astarte_data_access/mix.exs
modified +1/-1 libs/astarte_data_access/mix.lock
modified +2/-2 libs/astarte_events/mix.exs
modified +1/-1 libs/astarte_events/mix.lock
modified +2/-2 libs/astarte_generators/mix.exs
modified +1/-1 libs/astarte_generators/mix.lock
modified +2/-2 libs/astarte_rpc/mix.exs
modified +1/-1 libs/astarte_rpc/mix.lock
modified +1/-1 libs/astarte_test_suite/mix.exs
modified +1/-1 libs/astarte_test_suite/mix.lock
modified +1/-1 tools/astarte_device_fleet_simulator/mix.exs
modified +2/-2 tools/astarte_e2e/mix.exs
modified +1/-1 tools/astarte_e2e/mix.lock
modified +2/-2 tools/astarte_export/mix.exs
modified +1/-1 tools/astarte_export/mix.lock
modified +1/-1 tools/astarte_import/mix.exs

$ gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5 --jq '.files[] | select(.filename=="CHANGELOG.md") | .patch' 2>&1; echo "=== master commits since 2026-10-04 ==="; gh api repos/astarte-platform/astarte/commits -f since=2026-10-04T00:00:00Z --jq '.[] | "\(.sha[0:8]) \(.commit.committer.date) \(.commit.message|split("\n")[0])"' 2>&1
@@ -7,6 +7,12 @@ and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.
 
 ## [Unreleased]
 
+## [1.3.5] - 2026-10-05
+
+### Fixed
+
+- [astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available
+
 ## [1.3.4] - 2026-09-17
 
 ### Fixed
=== master commits since 2026-10-04 ===
{
  "message": "Not Found",
  "documentation_url": "https://docs.github.com/rest",
  "status": "404"
}gh: Not Found (HTTP 404)

$ gh api "repos/astarte-platform/astarte/commits?since=2026-10-04T00:00:00Z&per_page=100" --jq '.[] | "\(.sha[0:8]) \(.commit.committer.date) \(.commit.message|split("\n")[0])"' 2>&1
731bcddc 2026-10-05T09:34:08Z chore: forward port release-1.4

$ gh api repos/astarte-platform/astarte/commits/731bcddc --jq '{files: [.files[]|"\(.status) +\(.additions)/-\(.deletions) \(.filename)"]}' 2>&1
{"files":["modified +2/-1 apps/astarte_data_access/lib/astarte_data_access/database/database.ex","added +31/-0 apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0011_remove_replacement_data.ex","removed +0/-4 apps/astarte_housekeeping/priv/migrations/astarte/0004_create_session_key_type.sql","removed +0/-29 apps/astarte_housekeeping/priv/migrations/astarte/0005_create_to2_sessions_table.sql","removed +0/-15 apps/astarte_housekeeping/priv/migrations/astarte/0006_create_ownership_vouchers_table.sql","removed +0/-4 apps/astarte_housekeeping/priv/migrations/astarte/0007_add_expiry_to_ownership_vouchers.sql","removed +0/-1 apps/astarte_housekeeping/priv/migrations/realm/0020_drop_unconfirmed_devices_table.sql","removed +0/-4 apps/astarte_housekeeping/priv/migrations/realm/0021_add_device_id_to_ownership_voucher.sql","removed +0/-1 apps/astarte_housekeeping/priv/migrations/realm/0022_delete_vouchers.sql","removed +0/-1 apps/astarte_housekeeping/priv/migrations/realm/0023_drop_ownership_voucher_table.sql","removed +0/-1 apps/astarte_housekeeping/priv/migrations/realm/0024_drop_to2_session_table.sql","removed +0/-4 apps/astarte_housekeeping/priv/migrations/realm/0025_add_fdo_guid_to_devices.sql","modified +2/-2 mix.lock"]}

$ rg -n "APICompatVersion\s*=" internal/realm/service.go 2>&1; echo "--- dup/rpc in astrate ---"; rg -il "data_updater_plant|rpc" --glob '!*.md' internal/ cmd/ 2>&1 | head -10
588:const APICompatVersion = "1.2.2"
--- dup/rpc in astrate ---
internal/pairing/http.go

$ gh issue list --state open --limit 100 --json number,title,labels 2>&1 | python3 -c "
import json,sys
d=json.load(sys.stdin)
alarms=[i['number'] for i in d if any(l['name']=='mule-alarm' for l in i['labels'])]
other=[(i['number'],i['title'][:60]) for i in d if not any(l['name']=='mule-alarm' for l in i['labels'])]
print('open total:', len(d))
print('alarms:', len(alarms), alarms)
print('non-alarm:')
for n,t in other: print(' ', n, t)
"
open total: 29
alarms: 25 [118, 117, 116, 115, 114, 113, 112, 111, 110, 109, 108, 107, 106, 105, 104, 103, 102, 101, 100, 99, 98, 97, 96, 95, 94]
non-alarm:
  93 In internal/broker/aclhook.go the control/keyAgreement deny 
  92 keyAgreement: upstream published the full wire spec (082) — 
  78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milest
  1 Provide an Open Source IoT Platform unironically, unironical

$ rg -in "rpc" internal/pairing/http.go | head -5
300:// ["<msg>"]}}), and RPC-originated conflicts surface as the literal

No new gaps, but one genuinely new upstream fact (v1.3.5, shipped today). Recording it in `for-giulio.md`.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,16 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-10-05: v3.0's gate is unchanged, but the release sweep finally has something new to swallow — upstream shipped stable v1.3.5 today, seven hours ago, and it is wire-inert for Astrate. No issue filed, no task line queued.**  **The new release, measured.** `gh api repos/astarte-platform/astarte/releases` now shows a stable **v1.3.5** at `2026-10-05T14:11:44Z` where every prior run of this recipe (2026-09-10 onward, 20+ entries in this file) saw `v1.3.4` (2026-09-18) as newest. `compare/v1.3.4...v1.3.5` is **4 commits, 35 files**: two `chore: update horde` (`8d47574f`, `4c4222cf`), `chore(rm): bump xandra to 0.19.4` (`c4c63ebe`), `chore: prepare v1.3.5` (`8f0f5ef1`). The file list is `mix.exs`/`mix.lock` version bumps across all apps and libs, `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `astarte-dashboard/package.json`, `doc/mix.exs`, `README.md` — plus a **single** CHANGELOG entry, `## [1.3.5] - 2026-10-05` / `### Fixed` / "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available". That is the whole delta and none of it is Astarte-facing: Astrate implements no data_updater_plant and no Erlang RPC server (the only `rpc` hit under `internal/` is a comment at `internal/pairing/http.go:300` about AMQP reply-code conflicts), so no HTTP route, MQTT topic, AMQP control message or schema field moved. Same shape as the v1.3.2/v1.3.3/v1.3.4 maintenance releases already recorded here, so it is **not a v3.0 gap** and there is nothing to file.
+  **#92's parking condition is still unmet, now for the twenty-first run.** Your 2026-09-04 decision ("wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) — newest stable is `v1.3.5` (today), newest overall remains `v1.4.0-rc.6` (2026-09-30, prerelease). The post-rc.6 master delta has grown by exactly one commit: `731bcddc` "chore: forward port release-1.4" (2026-10-05T09:34Z), and it is the housekeeping/FDO migration relocation the 2026-10-03 run already measured inside rc.5→rc.6 (six `astarte_housekeeping/priv/migrations/{astarte,realm}/*.sql` files removed, table creation moved into migrations) plus a DUP-internal `apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0011_remove_replacement_data.ex` (+31) and `database.ex` (+2/-1). No Realm Management, AppEngine, Pairing or Trigger wire surface. **The milestone ledger is unchanged:** `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** (this milestone's whole ledger lives under `upstream-parity`, as the 2026-09-10 entry established), of #47–#89 **all remain closed**, and the authoritative open set is 29 issues — 25 alarms (#94–#118, one new today: **#118**, that recipe's business) and exactly four non-alarm: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** (untouched per standing instruction). `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`, re-checked), per **#90's** frozen decision.
+  **The standing item is unchanged and I am not re-litigating it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries its two rows, both tagged `1.4 experimental` with trigger "upstream 1.4 final: promoted or removed" — the mis-tagging the 2026-10-04 entry measured out of upstream's CHANGELOG (row #68's `async_operation` feature is 1.0.2 and stable for four years; row #67's `required` half is a plain 1.4.0-rc.0 addition, only its `encrypted` half being genuinely experimental and already tracked by #92/#93). v1.3.5 touches neither feature, so it changes nothing about your four pending decisions (a)–(d) in that entry, and I have not edited the register.
+  **Step 5 deliberately not taken.** I am *not* proposing "milestone v3.0 looks complete, verify and cut the tag", for the same reason as the twenty runs before this one: the section's declared scope is **1.2.2 → 1.3/1.4** and only the 1.3 half is delivered (now including v1.3.5, which is inert anyway). The 1.4 half is open — #92 unanswered, both experimental rows unreconciled, `APICompatVersion` not bumped — so the tag must not be cut. Your call when a stable v1.4.0 ships: answer #92, settle the register rows, run the final-phase bump, then cut the tag.
+  No `gh issue create`, nothing commented, closed or edited on GitHub, and **no `.mule/todo.md` lines**: there is nothing machine-checkable left in the 1.3 line, every open item already has an owner, and the one new fact is a wire-inert maintenance release.
+
+---
+
 - **github-issues triage run, 2026-10-04: nothing proposable for the 22nd run, and the reason the alarm pile has never been actioned is finally located — `mule/queue` forked from `main` 27 minutes before your four commits landed there, so this Pi has been running for a month with a dead-man's switch, a dashboard and a parked-decision check that all point at things which do not exist here.** 28 open issues (authoritative `--limit 100` list: `1, 78, 92, 93, 94…117`; note the recipe's own `--limit 40` command printed 27 of these and silently dropped **#94**, which is open with `mule-alarm` — so treat "27 open" in the 2026-10-03 entry as 28). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`; `8c61268` re-verified still an ancestor of `origin/mule/queue` this run), **#92** keyAgreement (parked on a stable v1.4.0 — measured again today, still unmet: newest upstream stable `v1.3.4`, newest overall `v1.4.0-rc.6`, 2026-09-30), **#78** FDO (`milestone-4.0`, already escalated), **#1** untouched per standing instruction.
   **The branch divergence, measured.** `git rev-list --left-right --count origin/main...HEAD` → **4 643**: `main` is 4 commits ahead, this branch is 643 ahead, merge-base `215d409` at **2026-09-04T21:25:31+02:00**. Your four commits are all dated 2026-09-04 21:40–21:52 — `f1d0069`, `ca47b35`, `99743ca`, `2a033d4`, the last titled "parking a decision on an upstream release now has something that notices". They landed on `main` **after** this branch forked and were never reconciled. `git merge-base --is-ancestor 2a033d4 HEAD` → false. So for 30 days this working tree has been missing:
   - **The fixed alarm body.** `tools/mule.sh:755-756` on this host still tells whoever opens an idle alarm to run `systemctl list-timers mule.timer` and `journalctl -u mule.service --since '-1 day'`. Measured here: `systemctl list-timers --all` returns exactly three mule units — `mule-survey.timer`, `mule-upstream-watch.timer`, `mule-planner.timer`. **`mule.timer` does not exist**, and `mule.service` is not how ticks run. Ticks are the one-shot cron entries inside the `BEGIN mule-daily-schedule` block (`crontab -l`; today `31 13`, `36 13`, `17 13` CEST). Your `f1d0069` message says the reason in one line — "its absence reads as a dead mule to anyone following the map, which is exactly what happened" — and the fix has been sitting on `main` since the night it was written. **All 24 open alarm bodies point a reader at two commands that cannot work on this machine.** That is the mechanical reason the pile has sat for a month: the alarm told the only human who would read it to go and look at a timer that is not there, and he found nothing, and it filed again the next morning.
   - **The dashboard fix, same commit.** `.agents/skills/astrate-dashboard/SKILL.md` on this host still says `mule.timer` (30 min) at **:35, :126**, still runs `systemctl list-timers mule.timer …` at **:66, :69**, and still offers `systemctl stop|start mule.timer` at **:111** as the manual recovery for *"ferma/riattiva il mulo"*. On this Pi that command is a no-op. The map Giulio follows when the mule looks stuck points at a timer that does not exist.

$ head -20 .mule/for-giulio.md | cut -c1-160; echo "=== todo.md unchanged? ==="; git status --porcelain
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **Milestone recipe run, 2026-10-05: v3.0's gate is unchanged, but the release sweep finally has something new to swallow — upstream shipped stable v1.3.5 to
  **#92's parking condition is still unmet, now for the twenty-first run.** Your 2026-09-04 decision ("wait for a stable v1.4.0"; an `-rc.N` does not satisfy it
  **The standing item is unchanged and I am not re-litigating it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries its two rows, both tagged `1.4 experimental` w
  **Step 5 deliberately not taken.** I am *not* proposing "milestone v3.0 looks complete, verify and cut the tag", for the same reason as the twenty runs before
  No `gh issue create`, nothing commented, closed or edited on GitHub, and **no `.mule/todo.md` lines**: there is nothing machine-checkable left in the 1.3 line

---

=== todo.md unchanged? ===
 M .mule/for-giulio.md

Recipe carried out. Working tree has one change (the escalation entry), no code, no git mutations, no `gh issue create`.

```
Done: Milestone recipe run — v3.0 audited, no gaps to file; new upstream fact (stable v1.3.5, shipped today) recorded in .mule/for-giulio.md
Files: .mule/for-giulio.md (one new entry, newest at top)
Verified: gh api repos/astarte-platform/astarte/releases -> newest stable v1.3.5 (2026-10-05T14:11:44Z), newest overall v1.4.0-rc.6 (no stable 1.4, so #92's 2026-09-04 parking condition still unmet);
  compare/v1.3.4...v1.3.5 -> 4 commits / 35 files, wire-inert (mix.exs+mix.lock version bumps, one CHANGELOG "Fixed" line about astarte_data_updater_plant RPC availability; Astrate has no DUP/RPC server);
  gh issue list --label milestone-3.0 --state all -> empty (ledger lives under upstream-parity); #47–#89 all closed; 29 open = 25 alarms + #93/#92/#78/#1; APICompatVersion == "1.2.2" (internal/realm/service.go:588)
Unsure: nothing — the one judgement call was step 5 ("looks complete, cut the tag"), deliberately not taken because v3.0's scope is 1.3/1.4 and the 1.4 gate is unresolved
Follow-ups: none (no task lines appended — nothing machine-checkable remains in the 1.3 line; docs/UPSTREAM-EXPERIMENTAL.md's mis-tagged rows are already escalated and unchanged by v1.3.5)
```
