slug: recipe-milestones
verdict: proposed
at:  becf3cd
ran: 2026-10-07T11:38:44Z on DietPi in 628s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md

$ wc -l .mule/research/upstream-parity-2026-08-22.md .mule/for-giulio.md .mule/todo.md docs/UPSTREAM-EXPERIMENTAL.md docs/ROADMAP.md docs/DESIGN.md 2>&1; ls .mule
wc: .mule/research/upstream-parity-2026-08-22.md: No such file or directory
   643 .mule/for-giulio.md
   287 .mule/todo.md
    25 docs/UPSTREAM-EXPERIMENTAL.md
   473 docs/ROADMAP.md
   830 docs/DESIGN.md
  2258 total
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

$ gh label list --limit 100 2>&1 | head -50
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

$ ls .mule/research 2>&1; ls .mule/tasks 2>&1 | head; ls docs/handoff 2>&1 | tail -20
ls: cannot access '.mule/research': No such file or directory
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md
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

$ gh issue list --label milestone-3.0 --state all --limit 50 2>&1; echo ---; gh issue list --label upstream-parity --state all --limit 100 2>&1
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

$ ls -la .mule/ ; echo ---; find . -name '*upstream-parity*' -not -path './.git/*' 2>/dev/null; echo ---; find . -name '*parity*' -not -path './.git/*' 2>/dev/null | head -30
total 13916
drwxr-xr-x  7 root root     4096 Oct  7 13:25 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root        0 Oct  7 13:04 .alarmed
-rw-r--r--  1 root root       13 Oct  7 13:28 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  6 22:29 .heartbeat
-rw-r--r--  1 root root     8192 Oct  7 13:30 .last-output
-rw-r--r--  1 root root        0 Oct  7 13:30 .lock
-rw-r--r--  1 root root       11 Oct  7 13:28 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 13824269 Oct  7 13:30 cron.log
drwxr-xr-x  2 root root     4096 Oct  6 22:49 failed
-rw-r--r--  1 root root   150453 Oct  6 20:39 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    15067 Oct  6 22:49 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  4 21:08 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   158037 Oct  6 22:49 todo.md
---
---
./internal/appengine/formats_parity_test.go
./internal/appengine/groups_parity_test.go
./internal/engine/mqttv1_parity_test.go
./docs/handoff/flow-parity-audit-2026-07-29.md

$ gh issue view 90 --json title,state,labels,body 2>&1 | head -60
{"body":"Reminder of a pending decision parked during the upstream-parity investigation (#47–#89). **Do not raise the version now** — this is a milestone-4.0 item.\n\n## Context\n\n`realm.APICompatVersion` (`internal/realm/service.go`) is currently hardcoded to `\"1.2.2\"` — the emulated upstream API level Astrate declares.\n\nEdgehog 0.13.1 gates just-in-time provisioning triggers (`device_registered`, `device_deletion_started`) behind a Realm Management version >= 1.3. Astrate's engine already supports those trigger kinds end-to-end — verified in #28 (Edgehog compatibility investigation) — but Edgehog refuses to use them because of the declared version string.\n\n## What this issue must cover when picked up\n\n1. **Audit first**: what else does declaring \"1.3\" imply for other consumers?\n   - `astarte-dashboard` is currently validated against 1.2.2 behavior (see #13, error-name bundle comparison). Check whether any dashboard code path feature-gates on RM version besides the policies UI gate (>= 1.1.1).\n   - Sweep the parity findings (#47–#89): several are \"1.3/1.4 features\" we deliberately do NOT implement yet — does declaring 1.3 overpromise any of them?\n   - Check device SDKs and `astartectl` for version-gated behavior.\n2. **Then decide one of**:\n   - bump the constant to `\"1.3.0-rc.0\"` (upstream's own pre-release numbering for that line);\n   - or make it configuration (`internal/config`) so operators can declare the level their consumers expect.\n3. Update `docs/COMPATIBILITY.md` deviation #10 wording accordingly (propose, don't edit directly).\n\n## Outcome if done today (why it's parked)\n\nBumping without the audit risks silently promising surfaces we haven't verified; not bumping keeps Edgehog's JIT provisioning disabled even though the engine supports it. Milestone-4.0 scope decision.","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACtvELrg","name":"milestone-4.0","description":"","color":"5319e7"}],"state":"CLOSED","title":"milestone-4.0: decide on bumping APICompatVersion from \"1.2.2\" to \"1.3.0-rc.0\" (or making it configurable)"}

$ git log --oneline -5 2>&1; echo ---; git log --all --oneline --diff-filter=AD -- '.mule/research/*' 2>&1 | head -20; echo ---; git status --short | head -20
becf3cd mule: recipe hygiene timed out
7a24500 mule: blocked docs-native-socket-event-schema
bc3dc9a mule: log docs-native-socket-security-scheme
87c18cf mule: docs-native-socket-security-scheme [auto]: `/astrate/v1/{realm}/socket` is documented as if public in docs/api/astrate_native_api.yaml — the spec has no `components.securitySchemes` at all and root `security: []` (yaml:21), and the operation declares no `security`, so a client generated from this spec sends no `Authorization` header and is answered 401 `{"errors":{"detail":"Unauthorized"}}` by `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112, bearer token read only from the Authorization header at middleware.go:131, 178-196). Add a scheme mirroring `a_aea` (docs/api/astarte_appengine_api.yaml:1436-1443: apiKey in header `Authorization`, described as a realm JWT carrying the `a_ch` claim) and set `security` on this operation only — root `security: []` stays so health/readiness/metrics/version remain documented as unauthenticated. If you model the Phoenix twin too, its credential is `?token=` in the query (internal/appengine/channels/ws.go:53-54), i.e. apiKey `in: query`, not header. Description notes worth carrying: the scheme is case-insensitive `Bearer` with an optional colon, and an unknown realm is 401 not 404 (middleware.go:68-72). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
3676ad9 mule: log docs-native-socket-sse-exact-accept
---
6d8cae6 mule: daily cross-project survey, 2026-09-04
1c618eb mule: daily cross-project survey, 2026-08-31
ce759e1 survey: recover 16 daily reports the broken commit check left staged (2026-07-28..2026-08-26), and take the tooling fix from main
f92d42a mule/research: commit the calibration survey, and mark it triaged
c6e0301 mule/research: seed the branch the daily survey timer will commit to
---

$ grep -rn "upstream-parity-2026-08-22" --include='*.md' . 2>/dev/null | grep -v '^./.git' | head -20
./.mule/for-giulio.md:89:  **New, and the reason this run is not a copy of the twenty before it: `.mule/research/` is listed in `.mule/.gitignore`, so the source document this milestone is built on — `.mule/research/upstream-parity-2026-08-22.md`, named at `.mule/milestones.md:114` — cannot exist in any clone, and the Pi's copy is gone.** Evidence: the last line of `.mule/.gitignore` is `research/`, a bare directory pattern that matches `.mule/research/` at any depth; the directory does not exist here; `git ls-files .mule` has no `research` entry; and `grep -rl 'upstream-parity-2026-08-22' --include='*.md'` over the tree matches **exactly one file, `.mule/milestones.md:114` itself** — there is no tracked equivalent (the only parity document under `docs/handoff/` is the v2.0 `flow-parity-audit-2026-07-29.md`). This is the sharper form of the 2026-09-06 note ("the section's source document does not exist in the repo"), and it is also why that note never got resolved: dropping the file back in place as-is re-dirties the tree on the next tick, because `.mule/.gitignore`'s own header states that any of its unlisted siblings "makes the tree dirty, which makes `tick` refuse to run and so turns the timer into a permanent silent no-op". **Yours, three ways:** (a) commit it once behind a negation (`!.mule/research/upstream-parity-2026-08-22.md`), so the milestone keeps a written reference; (b) move the source document to a tracked path (`docs/handoff/upstream-parity-2026-08-22.md`) and fix the `:114` pointer to match; (c) declare it lost and let the #47–#89 issue bodies plus the per-release CHANGELOG diffs stand as the reference. I did not pick one — every option edits a file that carries your decisions, and (a) versus (b) differ on a question only you can answer: is that document mutable working state or a frozen record?
./.mule/for-giulio.md:145:  **The recipe's standing item is satisfied, and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still has exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental* with reconcile trigger "upstream 1.4 final". v1.4 is not final, so both rows stay as written: nothing to promote, nothing to deprecate, and the register is accurate. `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision; the milestone-bump issue the recipe asks for **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": the milestone's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is done, and the 1.4 half is unresolved (#92 open, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut. Also still true and already escalated on 2026-09-06, so not repeated: the section's source document `.mule/research/upstream-parity-2026-08-22.md` does not exist in the repo.
./.mule/for-giulio.md:381:  upstream-parity-2026-08-22.md` is absent from `main` and from `origin/mule/research`.** The
./.mule/milestones.md:114:- Source document: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 /
./.mule/reports/log/20261003T180922Z-recipe-github-issues-proposed.md:233:  **New, and the reason this run is not a copy of the twenty before it: `.mule/research/` is listed in `.mule/.gitignore`, so the source document this milestone is built on — `.mule/research/upstream-parity-2026-08-22.md`, named at `.mule/milestones.md:114` — cannot exist in any clone, and the Pi's copy is gone.** Evidence: the last line of `.mule/.gitignore` is `research/`, a bare directory pattern that matches `.mule/research/` at any depth; the directory does not exist here; `git ls-files .mule` has no `research` entry; and `grep -rl 'upstream-parity-2026-08-22' --include='*.md'` over the tree matches **exactly one file, `.mule/milestones.md:114` itself** — there is no tracked equivalent (the only parity document under `docs/handoff/` is the v2.0 `flow-parity-audit-2026-07-29.md`). This is the sharper form of the 2026-09-06 note ("the section's source document does not exist in the repo"), and it is also why that note never got resolved: dropping the file back in place as-is re-dirties the tree on the next tick, because `.mule/.gitignore`'s own header states that any of its unlisted siblings "makes the tree dirty, which makes `tick` refuse to run and so turns the timer into a permanent silent no-op". **Yours, three ways:** (a) commit it once behind a negation (`!.mule/research/upstream-parity-2026-08-22.md`), so the milestone keeps a written reference; (b) move the source document to a tracked path (`docs/handoff/upstream-parity-2026-08-22.md`) and fix the `:114` pointer to match; (c) declare it lost and let the #47–#89 issue bodies plus the per-release CHANGELOG diffs stand as the reference. I did not pick one — every option edits a file that carries your decisions, and (a) versus (b) differ on a question only you can answer: is that document mutable working state or a frozen record?
./.mule/reports/log/20261003T180922Z-recipe-github-issues-proposed.md:430:   **New, and the reason this run is not a copy of the twenty before it: `.mule/research/` is listed in `.mule/.gitignore`, so the source document this milestone is built on — `.mule/research/upstream-parity-2026-08-22.md`, named at `.mule/milestones.md:114` — cannot exist in any clone, and the Pi's copy is gone.** Evidence: the last line of `.mule/.gitignore` is `research/`, a bare directory pattern that matches `.mule/research/` at any depth; the directory does not exist here; `git ls-files .mule` has no `research` entry; and `grep -rl 'upstream-parity-2026-08-22' --include='*.md'` over the tree matches **exactly one file, `.mule/milestones.md:114` itself** — there is no tracked equivalent (the only parity document under `docs/handoff/` is the v2.0 `flow-parity-audit-2026-07-29.md`). This is the sharper form of the 2026-09-06 note ("the section's source document does not exist in the repo"), and it is also why that note never got resolved: dropping the file back in place as-is re-dirties the tree on the next tick, because `.mule/.gitignore`'s own header states that any of its unlisted siblings "makes the tree dirty, which makes `tick` refuse to run and so turns the timer into a permanent silent no-op". **Yours, three ways:** (a) commit it once behind a negation (`!.mule/research/upstream-parity-2026-08-22.md`), so the milestone keeps a written reference; (b) move the source document to a tracked path (`docs/handoff/upstream-parity-2026-08-22.md`) and fix the `:114` pointer to match; (c) declare it lost and let the #47–#89 issue bodies plus the per-release CHANGELOG diffs stand as the reference. I did not pick one — every option edits a file that carries your decisions, and (a) versus (b) differ on a question only you can answer: is that document mutable working state or a frozen record?
./.mule/reports/log/20261004T181112Z-recipe-hygiene-timeout.md:72:- Source document: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 /
./.mule/reports/log/20261005T193452Z-recipe-milestones-proposed.md:15:[•] Read .mule/research/upstream-parity-2026-08-22.md source doc for v3.0
./.mule/reports/log/20261005T193452Z-recipe-milestones-proposed.md:21:✗ Read .mule/research/upstream-parity-2026-08-22.md failed
./.mule/reports/log/20261005T193452Z-recipe-milestones-proposed.md:22:Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
./.mule/reports/log/20260918T193641Z-recipe-milestones-proposed.md:12:✗ Read .mule/research/upstream-parity-2026-08-22.md failed
./.mule/reports/log/20260918T193641Z-recipe-milestones-proposed.md:13:Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
./.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:14:$ wc -l .mule/research/upstream-parity-2026-08-22.md docs/UPSTREAM-EXPERIMENTAL.md .mule/todo.md .mule/for-giulio.md 2>/dev/null; ls .mule/tasks/ 2>/dev/null
./.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:78:✗ Read .mule/research/upstream-parity-2026-08-22.md failed
./.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:79:Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
./.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:259:- Source document: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 /
./.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:291:$ rg -n -A6 'upstream-parity-2026-08-22\.md.*names a source doc|names a source doc' .mule/for-giulio.md | head -30
./.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:293:102-  upstream-parity-2026-08-22.md` is absent from `main` and from `origin/mule/research`.** The
./.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:343:$ wc -l .mule/research/upstream-parity-2026-08-22.md docs/UPSTREAM-EXPERIMENTAL.md .mule/recipes/astarte-upstream.md .mule/todo.md .mule/for-giulio.md 2>/dev/null; ls .mule/research/
./.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:897:✗ Read .mule/research/upstream-parity-2026-08-22.md failed

$ ls -t .mule/reports/log/ | head -30; echo ---; ls -t .mule/reports/
20261007T112533Z-recipe-hygiene-timeout.md
20261006T204912Z-docs-native-socket-event-schema-blocked.md
20261006T202927Z-docs-native-socket-security-scheme-done.md
20261006T201112Z-docs-native-socket-sse-exact-accept-done.md
20261006T185740Z-docs-native-socket-handshake-errors-blocked.md
20261006T184046Z-recipe-docs-sync-proposed.md
20261006T182510Z-recipe-code-review-timeout.md
20261006T175547Z-compat-note-v135-done.md
20261006T174931Z-probe-encrypted-mapping-enforcement-blocked.md
20261006T174353Z-recipe-astarte-upstream-proposed.md
20261006T165612Z-recipe-github-issues-proposed.md
20261006T141129Z-recipe-milestones-proposed.md
20261006T140535Z-recipe-hygiene-timeout.md
20261005T210409Z-recipe-docs-sync-timeout.md
20261005T201914Z-recipe-code-review-proposed.md
20261005T201617Z-recipe-astarte-upstream-proposed.md
20261005T195440Z-recipe-github-issues-proposed.md
20261005T193452Z-recipe-milestones-proposed.md
20261005T192311Z-recipe-hygiene-timeout.md
20261005T181428Z-docs-native-realm-name-pattern-done.md
20261005T180602Z-docs-native-phoenix-newevent-payload-done.md
20261005T175822Z-docs-native-metrics-content-negotiation-done.md
20261005T171138Z-docs-native-hk-version-description-done.md
20261005T113206Z-docs-native-ae-realm-version-endpoint-done.md
20261005T111650Z-recipe-docs-sync-proposed.md
20261005T105949Z-store-validatepipelinegraph-error-branches-done.md
20261005T022910Z-recipe-astarte-upstream-timeout.md
20261004T203510Z-store-latestindividual-empty-errnotfound-blocked.md
20261004T200628Z-store-aliasvalues-self-exclusion-test-done.md
20261004T193935Z-store-register-inhibit-preserve-blocked.md
---
log
hygiene.md

→ Read .mule/reports/log/20261006T141129Z-recipe-milestones-proposed.md

$ gh api "repos/astarte-platform/astarte/commits?since=2026-10-06T14:11:29Z&per_page=30" --jq '.[] | "\(.sha[0:8]) \(.commit.committer.date) \(.commit.message|split("\n")[0])"' 2>&1; echo "=== open non-alarm ==="; gh issue list --state open --limit 100 --json number,title,labels --jq '[.[] | select(([.labels[].name] | index("mule-alarm")) | not)] | "count=\(length)", (.[] | "\(.number) \(.title[0:70])")'
=== open non-alarm ===
count=4
93 In internal/broker/aclhook.go the control/keyAgreement deny comment (l
92 keyAgreement: upstream published the full wire spec (082) — decide whe
78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 ca
1 Provide an Open Source IoT Platform unironically, unironically

$ gh api repos/astarte-platform/astarte/releases --paginate --jq '.[] | "\(.tag_name) \(.published_at) prerelease=\(.prerelease)"' 2>&1 | head -6; echo "=== milestone-3.0 ==="; gh issue list --label milestone-3.0 --state all --limit 50 --json number 2>&1; echo "=== upstream-parity open ==="; gh issue list --label upstream-parity --state open --limit 50 --json number,title,labels --jq '.[] | "\(.number) \(.title[0:80]) | \([.labels[].name]|join(","))"'
v1.3.5 2026-10-05T14:11:44Z prerelease=false
v1.4.0-rc.6 2026-09-30T07:28:27Z prerelease=true
v1.3.4 2026-09-18T09:42:57Z prerelease=false
v1.4.0-rc.5 2026-08-20T12:28:04Z prerelease=true
v1.3.3 2026-08-07T14:33:58Z prerelease=false
v1.4.0-rc.4 2026-08-06T13:40:20Z prerelease=true
=== milestone-3.0 ===
[]
=== upstream-parity open ===
93 In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-6 | mule-review,upstream-parity
92 keyAgreement: upstream published the full wire spec (082) — decide whether the p | enhancement,upstream-parity,upstream-experimental
78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate) | enhancement,milestone-4.0,upstream-parity

→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ gh api "repos/astarte-platform/astarte/commits?since=2026-10-06T11:00:00Z&per_page=10" --jq '.[] | "\(.sha[0:8]) \(.commit.committer.date) \(.commit.message|split("\n")[0])"' 2>&1; echo "=== tags ==="; gh api "repos/astarte-platform/astarte/tags?per_page=15" --jq '.[].name' 2>&1; echo "=== astrate HEAD ==="; git log --oneline -3; git status --short | head
1da04832 2026-10-06T11:01:17Z chore: forward port release-1.4
=== tags ===
v1.4.0-rc.6
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
v1.4.0-rc.0
v1.3.5
v1.3.4
v1.3.3
v1.3.2
v1.3.1
v1.3.0
v1.3.0-rc.2
v1.3.0-rc.1
=== astrate HEAD ===
becf3cd mule: recipe hygiene timed out
7a24500 mule: blocked docs-native-socket-event-schema
bc3dc9a mule: log docs-native-socket-security-scheme

→ Read .mule/for-giulio.md [limit=60]

$ grep -n "b6d46ad4\|958bb9fe\|rc\.6" .mule/for-giulio.md | head -20; echo "=== todo milestone lines ==="; grep -n "milestone-3.0\|keyAgreement\|aclhook" .mule/todo.md | head -20
13:- **COMPATIBILITY.md wording proposal for upstream v1.3.5 (newest stable, 2026-10-05; v1.4.0 is still rc.6-only) — a one-word delta on the open v1.3.2/v1.3.4 proposal below, nothing re-derived.** v1.3.5 is maintenance-only: `gh api repos/astarte-platform/astarte/releases` shows it published `2026-10-05T14:11:44Z`, not a prerelease, with a body of a single `### Fixed` line — "[astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available" — and `compare/v1.3.4...v1.3.5` is 4 commits / 35 files whose only non-`mix.exs`/`mix.lock` files are `.github/workflows/astarte-build-workflow.yaml`, `CHANGELOG.md`, `README.md`, `astarte-dashboard/package.json` and `docker-compose.yml`. Not one source file changed, so no route, MQTT topic, control message or interface-schema field moved (Astrate implements no `data_updater_plant` and no Erlang RPC server). **The only change a v1.3.5-aware doc carries: the v1.3.4 entry's §Infrastructure-differences sentence — "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target" — reads "…(**newest stable v1.3.5**) as the target".** Everything else stands unchanged: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"); that entry's "v1.4.0 is still rc.5-only" header is already superseded to rc.6 by the open rc.6 delta entry below, not here. This is the same one-word delta recorded on 2026-10-05 in this file ("newest stable v1.3.5"), folded in as asked rather than re-derived. Applied to nothing: `docs/COMPATIBILITY.md` untouched and still 1.2.2-targeted (it contains no `1.3` reference at all), `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). Raw: [v1.3.5](https://github.com/astarte-platform/astarte/releases/tag/v1.3.5).
23:  **One detail to fold into that review rather than queue separately:** the rewritten comment at `internal/broker/aclhook.go` (the `control/keyAgreement` deny) says "…upstream's wire spec (082-key_agreement_protocol.md) only ships with v1.4.0-rc.5, v1.3.3 being the newest stable tag". The 082 citation, the "deliberately denied as Astrate's own choice" framing and the unchanged deny + `aclhook_test.go` case are exactly what #93 asked for, but the version reference has aged — upstream's newest stable is **v1.3.5** (2026-10-05) and newest overall `v1.4.0-rc.6` (2026-09-30). No separate task: one word in the same diff when someone finally reads it.
27:- **Milestone recipe run, 2026-10-06: v3.0 verified unchanged — still gated on upstream — and the master delta since the 2026-10-05 run is a closed set of 13 commits, every one wire-inert for Astrate. No issue filed, no `.mule/todo.md` line queued.** Release sweep is byte-identical to yesterday's (`gh api repos/astarte-platform/astarte/releases`): **no stable v1.4.0** — newest stable `v1.3.5` (2026-10-05), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition (your 2026-09-04 "wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) stays unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** — this milestone's whole ledger lives under `upstream-parity`, where the open set is unchanged: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`); **#47–#89** otherwise all closed. Non-alarm open issues total exactly four (plus the 25 `mule-alarm` #94–#118): the three above and **#1** (untouched per standing instruction). `APICompatVersion` still `"1.2.2"` (`internal/realm/service.go:588`, re-checked).
32:- **astarte-upstream recipe run, 2026-10-05 (evening): no upstream movement that touches anything Astrate implements, so nothing was queued — no `.mule/todo.md` line, no `gh issue create`, nothing commented, closed or edited on GitHub.** One thing *is* new since this recipe's last run (`2026-10-04T20:40:52+02:00` ran; the `2026-10-05T04:29:10+02:00` run timed out): stable **v1.3.5**, published `2026-10-05T14:11:44Z`. The milestone recipe already recorded it as inert today; I re-measured it rather than trusting that, because "a release note is a claim, not a fact" is this recipe's own rule. `gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5` is **4 commits, 35 files** — `8d47574f` + `4c4222cf` `chore: update horde`, `c4c63ebe` `chore(rm): bump xandra to 0.19.4`, `8f0f5ef1` `chore: prepare v1.3.5` — and **every app/lib file in it is a `mix.exs`/`mix.lock` version bump**. The rest is `CHANGELOG.md` (+6: one `### Fixed` entry, "[astarte_data_updater_plant] Increase RPC server reliability"), `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `README.md`, `astarte-dashboard/package.json`. Not one source file changed, so no route, MQTT topic, BSON encoding, control message or interface-schema field moved: Astrate implements no `data_updater_plant` and no Erlang RPC server. Wire-inert, same shape as v1.3.2/v1.3.3/v1.3.4. **Master since `v1.4.0-rc.6` is still the single commit `731bcddc`** ("chore: forward port release-1.4"; `gh api .../commits?since=2026-10-05T09:34:00Z` returns it and nothing newer), the migration-relocation commit the milestone run measured today as inert — so both wire-bearing post-rc.6 commits (`b6d46ad4`, `958bb9fe`) still have their owners in the queue and nothing new joins them.
33:  **The parked-decisions row for #92 is still unmet, re-checked against the tag list rather than the release list** (this host has no `.mule/waiting-on.md` — absent from `HEAD`, already reported; I read it out of `origin/main` read-only): `gh api repos/astarte-platform/astarte/tags?per_page=100` gives newest stable **v1.3.5** and newest overall **v1.4.0-rc.6**, with `v1.4.0-rc.0`…`rc.6` and no stable `v1.4.0`. The row asks for a **stable** tag, which an `-rc.N` does not satisfy, so the twenty-second run ends the same way as the twenty-one before it.
45:  **#92's parking condition is still unmet, now for the twenty-first run.** Your 2026-09-04 decision ("wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) — newest stable is `v1.3.5` (today), newest overall remains `v1.4.0-rc.6` (2026-09-30, prerelease). The post-rc.6 master delta has grown by exactly one commit: `731bcddc` "chore: forward port release-1.4" (2026-10-05T09:34Z), and it is the housekeeping/FDO migration relocation the 2026-10-03 run already measured inside rc.5→rc.6 (six `astarte_housekeeping/priv/migrations/{astarte,realm}/*.sql` files removed, table creation moved into migrations) plus a DUP-internal `apps/astarte_data_access/lib/astarte_data_access/database/migrations/astarte/0011_remove_replacement_data.ex` (+31) and `database.ex` (+2/-1). No Realm Management, AppEngine, Pairing or Trigger wire surface. **The milestone ledger is unchanged:** `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** (this milestone's whole ledger lives under `upstream-parity`, as the 2026-09-10 entry established), of #47–#89 **all remain closed**, and the authoritative open set is 29 issues — 25 alarms (#94–#118, one new today: **#118**, that recipe's business) and exactly four non-alarm: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** (untouched per standing instruction). `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`, re-checked), per **#90's** frozen decision.
52:- **github-issues triage run, 2026-10-04: nothing proposable for the 22nd run, and the reason the alarm pile has never been actioned is finally located — `mule/queue` forked from `main` 27 minutes before your four commits landed there, so this Pi has been running for a month with a dead-man's switch, a dashboard and a parked-decision check that all point at things which do not exist here.** 28 open issues (authoritative `--limit 100` list: `1, 78, 92, 93, 94…117`; note the recipe's own `--limit 40` command printed 27 of these and silently dropped **#94**, which is open with `mule-alarm` — so treat "27 open" in the 2026-10-03 entry as 28). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`; `8c61268` re-verified still an ancestor of `origin/mule/queue` this run), **#92** keyAgreement (parked on a stable v1.4.0 — measured again today, still unmet: newest upstream stable `v1.3.4`, newest overall `v1.4.0-rc.6`, 2026-09-30), **#78** FDO (`milestone-4.0`, already escalated), **#1** untouched per standing instruction.
66:  **Gate re-verified, unchanged.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — nothing shipped upstream since the 2026-10-03 run, so #92's parking condition (your 2026-09-04 decision) is still unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** (this milestone's ledger lives under `upstream-parity`); of #47–#89 **all remain closed**; open non-alarm set unchanged — **#93** (`mule-review`, own recipe path), **#92** (parked), **#78** (`milestone-4.0`), **#1** (untouched). The post-rc.6 master delta is now a closed set of **8 commits** (`gh api .../commits -f since=2026-09-30T07:28:27Z`), of which exactly two carry wire surface and **both already have owners in the queue**: `b6d46ad4` (#2237 unexpected object keys) → done as `appengine-unexpected-object-key`, and `958bb9fe` (#2232 `device_empty_cache_received`) → queued as `device-empty-cache-received-trigger` (currently `[!]` blocked). Neither is from stable 1.3 or 1.4, so neither is a v3.0 gate.
76:- **The Legion Go answers ssh but its Docker is not running, so every "probe live upstream" line is blocked on you, not on the hardware (measured 2026-10-03).** `ssh -o ConnectTimeout=5 legion 'echo ok'` returns `ok`, but `docker ps` fails with `dial unix /home/atsetilam/.docker/desktop/docker.sock: no such file or directory`, `systemctl is-active docker` answers `inactive`, and `sudo -n true` answers `è necessaria una password` — so the mule cannot start it. It is Docker Desktop's socket, not a rootful `dockerd` the mule could launch as a user, and no rootless extras are installed. Consequence for the queue: `appengine-payload-reason-status-map` (the remaining eight `payload.Reason` values that still answer 500) cannot be measured the way its line demands. Today's `appengine-missing-required-422` was therefore mapped from the frozen upstream source at `b6d46ad4` (`Device.validate_required_mappings/2` → `fallback_controller.ex:215-219` → `ErrorView "422_missing_required_mapping.json"`), and the commit says so plainly. Nothing about that sibling line needs rewriting — only Docker up on the Legion Go.
87:  **Verified unchanged, no new gaps.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition, frozen by you on 2026-09-04 ("wait for a stable v1.4.0", an `-rc.N` does not satisfy it), is still unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Of #47–#89 **all are closed**; the open non-alarm set is unchanged — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** (untouched per standing instruction). rc.6 adds no v3.0 gap: it is already audited in the two entries below, and I re-measured its non-FDO half myself (`compare/v1.4.0-rc.5...v1.4.0-rc.6`, 27 commits) — outside FDO the only library changes are a `# coveralls-ignore-stop` line in `apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex:43` and `realms/queries.ex` *removing* the ownership_voucher / to2_sessions / unconfirmed_devices table creation (moved into migrations). No Realm Management, AppEngine or Pairing wire surface moved.
97:- **Delta on the two open COMPATIBILITY.md wording entries further down this file (the v1.3.4 one folding in v1.3.3, and the original v1.3.2 one): upstream `v1.4.0-rc.6` is wire-inert for Astrate, so nothing in that proposal is re-derived — only its version reference and one line of its "not yet emulated" list change.** Not applied: `docs/COMPATIBILITY.md` is untouched (still the 1.2.2 target, deviation 18's "which upstream 404s" at `docs/COMPATIBILITY.md:360`, `APICompatVersion = "1.2.2"` at `internal/realm/service.go:588`), and the FDO half of the rc.6 delta already has its own entry below rather than being folded in here.
101:  - **The two edits the proposal needs, both yours to approve:** (1) the 2026-09-18 entry's parenthetical "(v1.4.0 is still rc.5-only)" becomes "(v1.4.0 is still rc.6-only)" — re-verified this session that rc.6 is the newest 1.4 prerelease and there is still **no stable v1.4.0** (`gh api repos/astarte-platform/astarte/releases`); (2) its "**FDO authentication** (pairing, disabled by default)" item, and the same "experimental FDO pairing auth" phrase in the proposed doc sentence, understates the 1.4 line, which ships a full ownership-voucher surface (upload with mandatory `hw_id`, deletion, re-TO0, per-voucher `expiry`, a spec-mandated realm-less TO2 path) that Astrate does not emulate — #78, deliberately parked at `.mule/milestones.md:124`. "FDO ownership vouchers (1.4 line, not stable upstream)" says it without claiming the rc.5 shape.
102:  Unverified: the upstream side of the item list is the upstream-watch recipe's `v1.4.0-rc.5...v1.4.0-rc.6` comparison — the same measurement as the FDO entry below, taken there and not re-derived here (a raw fetch of `models/Device/index.ts` at the rc.6 tag did not resolve from this box). I re-verified only the Astrate-side claims above and the release list. No `.mule/todo.md` lines: the deliverable is the delta, and the proposal itself stays open.
106:- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published 2026-09-30T07:28Z, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and the one path #78's body states as fact no longer exists upstream. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix") deletes the `realm_name` path parameter from all seven TO2 `operation` specs in `fdo_onboarding_controller.ex`. **This breaking change is not in the rc.6 CHANGELOG entry** (which lists only the to0 route, voucher deletion, the `ASTARTE_BASE_URL_*` rename, device-deletion cascade, `device_id` on upload, immediate registration, the two x5chain/voucher-store fixes and the cross-realm delete fix) — so anyone tracking upstream by its changelog, including #78's body, misses it. The realm now comes from the voucher row: the global-keyspace query `fetch_device_id_and_ownership_voucher_and_realm/1` returns `{realm, device_id, voucher_data}` and the `:fdo` plug fetches the session by GUID alone (`plug/fdo_session.ex`: `Session.fetch(guid)`); `RvTO2Addr` no longer carries the realm in the hostname either ("the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure"). **Yours:** #78's acceptance is "a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end" — against rc.6 that SDK posts to `/fdo/101/msg/{60..70}` with no realm, so mirroring rc.5's shape would mean mirroring a path upstream has deleted. No action while parked; the body is just stale.
117:- **One premise in #78's own comment is now false.** It says "upstream ships almost no public FDO documentation (code-level only); ours must exceed that." rc.6 adds two documented sections: `doc/pages/user/035-register_device.md` gained `## FIDO Device Onboard` (Owner Keys, Ownership Vouchers) and `## Credentials Secret Lifecycle`, and `doc/pages/architecture/050-pairing_mechanism.md` gained `## FIDO Device Onboarding` (TO0, TO2, Owner Keys management) — commits `162062e`, `d1c39fb`. "Docs must exceed upstream" is still the right requirement, but the gap is far smaller than #78 assumes, and those two pages are now the cheapest available scope input.
118:- **Where Astrate stands, for the size estimate:** it has **no FDO surface at all** — the only FDO reference in the tree is the comment on the realm health probe (`internal/pairing/http.go:79`, `internal/pairing/service.go:368`), and `.mule/milestones.md:124` still lists #78 as a deliberately parked item. So none of this is a delta to partially-built code; every item above is net-new scope for a "Size L" issue that is blocked on "who runs a rendezvous server?". The rc.6 delta does not change that call, but it adds four things any eventual implementation must decide: a spec-mandated device-side path with no realm (so any deviation is ours alone), a mandatory `hw_id` in the upload body (so the owner must supply the device id *before* the device has onboarded), an upload that creates a visible realm device as a side effect, and a per-voucher `expiry` the operator has to poll to know when the rendezvous stops serving.
119:  Unverified: read-only — every line above comes from the compare API and the rc.6-tagged sources, and the 409/500 mappings come from `FallbackController` source rather than a live call (no Docker on this box). Deliberately no `.mule/todo.md` lines: the deliverable is the scope update, and #78 stays parked until you answer the rendezvous question.
=== todo milestone lines ===
107:- [x] broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [approved 2026-09-04]
109:- [!] broker-offline-acl-tests: in `internal/broker/aclhook_test.go`, unit-test the offline-delivery ACL — `offlineACL.ownershipOf` cache hit within TTL, TTL expiry triggering a reload, and a load failure caching as empty and denying (aclhook.go:116-138) — plus a T1 `OnACLCheck` check that the offline branch (delivery to a session-present but disconnected device, aclhook.go:196-202) consults the store. The rule currently has zero direct tests; `TestCheckACLMatrix` exercises only the pure `checkACL`. [approved 2026-09-04] — BLOCKED: wrote nothing
206:- [x] broker-acl-coldstart-fallback-flood [auto]: `syncOwnershipOf` (internal/broker/authhook.go:212-237, reached from the ACL miss path aclhook.go:196-206) does a full synchronous `GetDevice`+`GetInterface` store read for every distinct interface name a device publishes to, with no rate limit — `introspectionReloadDebounce` (authhook.go:42-45, "so an adversarial topic flood cannot hammer the database") throttles only `refreshIfStale`'s full reload, and the fallback's per-name cache write means N distinct bogus names in one window cause N synchronous full-device reads per second, an open anti-flood hole the pre-fix denied-miss path did not have. Gate the fallback to one sync read per session per debounce window (or fold it into a full reload that stamps `lastIntroLoad`), keep `TestBrokerACLColdStartIntrospectionMiss` (one name) green, and add a T1 test in broker_test.go with a GetDevice-counting fake: connect cold, publish to K>=2 distinct unknown-interface names, assert no more than one fallback read was made.
207:- [x] broker-offlineacl-entry-eviction [auto]: `offlineACL.entries` (internal/broker/aclhook.go:100-162) is append-only — `ownershipOf` upserts one entry per CN ever ACL-checked offline and nothing ever deletes, so a long-lived broker grows the map unboundedly across device churn even after the 10s `offlineACLCacheTTL` passes (the entry is kept and re-stamped, not reaped). Evict stale entries lazily on access (drop an entry whose `loadedAt` predates some multiple of the TTL before inserting another, or cap+LRU), add an injectable `now func() time.Time` in the same style as `lifecycleHook.now`, and cover in a container-free unit test: seed N entries, advance the clock, access one, assert the map stays bounded.

$ gh issue list --label upstream-parity --state closed --limit 100 --json number,title,stateReason,closedAt --jq '.[] | "\(.number)\t\(.stateReason)\t\(.closedAt[0:10])\t\(.title[0:65])"'
91	COMPLETED	2026-09-04	Pairing health endpoint divergence: upstream serves GET /pairing/
89	COMPLETED	2026-08-23	Dashboard flow-block schema mismatch (split_map/virtual pools har
88	COMPLETED	2026-08-23	Flow auth: support a_f JWT claim
87	COMPLETED	2026-09-04	Flow block: lua_map — needs embedded Lua runtime (parked)
86	COMPLETED	2026-08-23	Flow: pipeline source DSL — keep DAG-JSON as documented deviation
85	COMPLETED	2026-08-23	Flow API: user-defined composite blocks
84	COMPLETED	2026-08-23	Flow blocks: virtual_device_pool / dynamic_virtual_device_pool
83	COMPLETED	2026-08-23	Flow blocks: mqtt_source/mqtt_sink (+modbus_tcp_source?) demand-d
82	COMPLETED	2026-08-23	Flow blocks: http_source/http_sink (demand-driven)
81	COMPLETED	2026-08-23	Flow block: json_path_map
80	COMPLETED	2026-08-23	Flow blocks: pure-transform set (to_json, update_metadata, split_
79	COMPLETED	2026-08-24	Verify registration-limit-reached HTTP status vs upstream
77	COMPLETED	2026-08-24	Verify per-service version endpoints (GET /version, GET /v1/{real
76	COMPLETED	2026-08-22	Housekeeping: GET /v1/realm-defaults/replication — decide reject/
75	COMPLETED	2026-08-22	Housekeeping: decide realm-deletion gating/preconditions vs alway
74	COMPLETED	2026-08-22	Housekeeping: PATCH /v1/realms/{realm} (jwt key, registration lim
73	COMPLETED	2026-08-22	Housekeeping: default datastream retention injection env var (ups
72	COMPLETED	2026-08-22	Realms: datastream_maximum_storage_retention ceiling (create/patc
71	COMPLETED	2026-08-22	Pairing: realm-scoped health check GET /v1/{realm}/health (upstre
70	COMPLETED	2026-08-22	Triggers: audit wildcard semantics (interface_name '*', match_pat
69	COMPLETED	2026-08-24	Verify unknown-realm HTTP status on RM endpoints against upstream
68	COMPLETED	2026-09-04	Decide async_operation=false params vs documented always-sync dev
67	COMPLETED	2026-09-04	Interfaces: decide handling of required and encrypted mapping fie
66	COMPLETED	2026-08-22	Realm Management: detailed=true interface listing with full mappi
65	COMPLETED	2026-08-22	Policies: handler-overlap rejection + retry_times coupling + pref
64	COMPLETED	2026-08-22	Triggers: decide AMQP action behavior (validate-reject vs NATS-fo
63	COMPLETED	2026-08-22	Triggers: HTTP action validation limits (URL/method/header blockl
62	COMPLETED	2026-08-22	Realm Management: audit install/update/delete error codes and sta
61	COMPLETED	2026-08-22	Realm Management: audit interface/mapping validation matrix again
60	COMPLETED	2026-08-22	Realm Management: GET config/datastream_maximum_storage_retention
59	COMPLETED	2026-08-22	AppEngine: group create-body validation + UUID-v1 from_token for 
58	COMPLETED	2026-08-22	AppEngine: PATCH requires Content-Type application/merge-patch+js
57	COMPLETED	2026-08-24	AppEngine: audit server-write error taxonomy against upstream
56	COMPLETED	2026-08-22	AppEngine: object-aggregate response formats table and disjoint_t
55	COMPLETED	2026-08-22	AppEngine: complete data-query parameter set
54	COMPLETED	2026-08-22	AppEngine: GET /groups/{g} show-group endpoint
53	COMPLETED	2026-08-22	AppEngine: group-scoped device PATCH + interface data access unde
52	COMPLETED	2026-08-22	AppEngine: devices-by-alias full mirror (PATCH via alias + interf
51	COMPLETED	2026-08-22	keyAgreement control topic + encrypted mappings (upstream 1.4 exp
50	COMPLETED	2026-08-22	Broker/engine hardening: discard malformed introspection/purge/ca
49	COMPLETED	2026-08-22	MQTT v1: binaryblob properties re-send as BSON subtype 0 in consu
48	COMPLETED	2026-08-22	MQTT v1: tolerate devices connecting with an empty introspection
47	COMPLETED	2026-08-22	MQTT v1: accept device capabilities incl. purge_properties_compre

$ gh issue view 90 --json stateReason,closedAt,comments --jq '.stateReason, .closedAt, (.comments[] | "[\(.createdAt[0:10])] \(.body[0:300])")' 2>&1 | head -20
COMPLETED
2026-08-23T21:06:34Z
[2026-08-23] Decision frozen 2026-08-23: **do not bump now, and never make it configurable.**

- `APICompatVersion` stays `1.2.2` until Astrate implements the FULL 1.3 surface; it is bumped only in the same change that completes that surface, after reconciling every row in docs/UPSTREAM-EXPERIMENTAL.md tagged wi

$ grep -n "1\.3\|1\.4\|upstream\|parity\|APICompat" docs/ROADMAP.md | head -40
20:   external-bus intake (§1.4), NATS/HTTP trigger forwarding (§1.1), and Timescale toolkit `lttb`
23:   interface (hexagonal-lite, §1.3), the interface is defined in A's milestone step, with B
60:| **M3** | `internal/auth` — JWT + Astarte authz claims | T1 green; claim-matching parity table passes |
88:| HTTP routing | **stdlib `net/http`** (Go ≥ 1.22 `ServeMux` method+wildcard patterns) | Zero deps; upstream-compatible paths need no framework |
90:| Postgres | `github.com/jackc/pgx/v5` (+ `pgxpool`) | Frozen in §1.3 |
91:| Migrations | `github.com/golang-migrate/migrate/v4` with `source/iofs` + `go:embed` | Frozen in §1.3 |
100:| Test containers | `github.com/testcontainers/testcontainers-go` (+ postgres module), image `timescale/timescaledb:latest-pg16` | §5.4 parity with production image |
119:### 1.3 Verification / gate
130:**Dependency rule check:** `pkg/*` imports stdlib + bson only; zero `internal/*` imports (§1.3).
136:| 1.1 | `pkg/deviceid/deviceid.go` | `type ID [16]byte`; `Parse(s string) (ID, error)` (exactly 22-char unpadded base64url → 16 bytes), `String()`, `FromUUID`/`UUID()`, `Random()`, `FromNamespace(ns uuid, payload string)` (UUIDv5 deterministic derivation, astartectl parity) | 140 |
149:| 1.3 | `pkg/astarteapi/envelope.go` | `WriteData(w, status, v)` → `{"data": v}`; `WriteError(w, status, detail)` → `{"errors":{"detail":"..."}}`; canonical upstream error constructors (404 `"Device not found"`, 401, 403, 422 shapes); request body `{"data": ...}` unwrapper with size cap | 120 |
150:| 1.4 | `pkg/astarteapi/envelope_test.go` | Golden JSON for every constructor; unwrap rejects missing `data`, oversized bodies | 90 |
155:### 2.3 Step M1.3 — `pkg/interfaceschema`
162:| 1.8 | `pkg/interfaceschema/testdata/valid/*.json` | ≥ 12 fixtures: vendored upstream `astarte-platform/standard-interfaces` (genericsensors, device-info…) + object-aggregated + parametric + properties w/ allow_unset + every value type | — |
167:| 1.13 | `pkg/interfaceschema/compat.go` | `CheckMinorUpgrade(old, new *Interface) error` — additive-mappings-only, no mutation of existing mapping attributes, same type/ownership/aggregation (§2.6 versioning parity) | 100 |
173:### 2.4 Step M1.4 — `pkg/payload`
198:**Goal:** the single Postgres access layer (§1.3: imported by domains, imports none of them).
246:| 3.2 | `internal/auth/claims.go` | `a_aea/a_rma/a_pa/a_ha/a_ch` extraction; `"<verb-regex>::<path-regex>"` matching with **implicit anchoring** (upstream parity), evaluated against method + path relative to realm base | 130 |
258:**M3 gate:** T1 green; parity table reviewed against upstream `astarte_rpc`/dashboard token semantics.
276:- **T2 HTTP flow tests (golden bodies):** Flow A 201 + 44-char secret; second register pre-credentials returns a *different* secret; post-credentials register → 422 upstream-shaped; bad `hw_id` (21 chars, padded, non-url) → 422. Flow B: 201 `client_crt` that (a) parses, (b) has CN `<realm>/<device>`, (c) chains to the realm CA returned by Flow C — verified with `crypto/x509` *and* an `openssl verify` exec smoke; wrong secret → 401 with the same body/timing class as unknown device (uniform error assertion); rate limit → 429. Flow C: info golden body (§4.4 shape exactly); verify on fresh cert → `valid:true` + `until`; on expired (issue with 1 s TTL, sleep) → `EXPIRED`; on cert from a different CA → `INVALID`.
277:- Inhibited device: credentials → 403 parity.
299:| 5.8 | `internal/broker/broker.go` | `New(cfg, store, intake, sink)`: mochi server, TLS listener `:8883` (`RequireAndVerifyClientCert`, client-CA pool), optional `:1883` behind `insecure_dev_mode`, MQTT 3.1.1/5.0 accept, hook registration, `OnPublish`→`Intake.Submit` with **deferred-ack wiring** (QoS ≥ 1 PUBACK held until `Ack()`, §1.4/§5.3), graceful stop | 220 |
312:**Goal:** sharded ordered pipeline (§1.4) + validation (§2.6) + persistence + control channel
321:| 6.2 | `internal/engine/router.go` | `Engine.Submit(InboundMessage)` (implements `broker.Intake`): `shard = FNV1a(deviceID) % N` (default 16), bounded chans (default 4096); QoS ≥ 1 full ⇒ blocking submit (ack deferral = backpressure); QoS 0 full ⇒ drop + metric (§1.4); per-shard goroutine with panic recover-and-log (§6) | 160 |
324:| 6.5 | `internal/engine/batch.go` | Per-shard micro-batch: flush at 64 rows or 50 ms (§1.4) through `store.AppendBatch` / property upserts in one tx; **`Ack()` called only after commit** (§5.3); DB-outage parking with exponential backoff | 200 |
335:| 6.11 | `internal/engine/triggers/events.go` | Astarte trigger event JSON payload shapes (SimpleEvent envelope parity: realm, device_id, event type fields) | 120 |
337:| 6.13 | `internal/engine/stream/bus.go` | In-process fan-out: `Subscribe(realm, filter) (<-chan Event, cancel)`; non-blocking sends, slow-consumer drop + metric (§1.1/§1.4) | 100 |
343:- **T1:** ordering property test — N messages for one device through a multi-shard engine with a synthetic slow store arrive persisted in publish order; cross-device parallelism observed. Backpressure: filled shard blocks QoS 1 submit / drops QoS 0 + metric. Reject-reason table covering every §2.6 failure (uses M1 fixtures). zlib: golden `producer/properties` blob round-trip; zip-bomb (1 GiB declared / absolute ceiling) rejected. Trigger matching: fixture upstream trigger JSONs → match/no-match table; webhook retry on 500-then-200 (httptest); event JSON golden vs upstream SimpleEvent shape.
354:**Goal:** the operator-facing APIs (§3.7), wire-shaped to upstream. Two Phase 3 sessions:
361:| 7.1 | `internal/realm/service.go` | Interfaces: install (parse+validate via `pkg/interfaceschema`), update (CheckMinorUpgrade enforcement), delete (upstream draining rules: only major 0 / not in any introspection — parity), list/get; triggers CRUD (validate action shape); `config/auth` JWT key rotation (§4.2); every mutation → store NOTIFY + in-process engine callback | 220 |
371:| 7.6 | `internal/appengine/service.go` | Device list (paged)/get (upstream status body: introspection, connected, stats…), PATCH (aliases, attributes, `credentials_inhibited`), devices-by-alias; groups CRUD + membership | 240 |
379:- **T2 golden suites per service:** envelopes + status codes byte-compared (404 device, 409/422 interface conflicts, 401/403 split). Interfaces: minor-bump additive accepted, mapping-mutation rejected (CheckMinorUpgrade), major coexistence, delete-while-introspected rejected. Housekeeping: created realm immediately serves pairing (cross-domain test with M4). AppEngine: datastream query boundaries (`since` inclusive vs `since_after` exclusive), `limit`, descending default ordering parity, `downsample_to` bucket-count correctness; longinteger > 2^53 round-trips as string through publish→query; property tree shape golden.
448:| S2 | M1.1–M1.2 + M1.3 (`pkg/deviceid`, `astarteapi`, `interfaceschema`) | ~1,800 |
449:| S3 | M1.4 (`pkg/payload`) | ~1,150 |

$ grep -n "1\.3\|1\.4\|APICompat\|upstream" docs/DESIGN.md | head -30; echo "=== COMPAT deviations count ==="; grep -c "^" docs/COMPATIBILITY.md
16:| Concern | Astarte (upstream) | Astrate |
34:API is compatible, so the upstream dashboard *may* work later, but it is not a v1 acceptance
36:WebSocket/SSE live stream; see §1.4).
54:| Astarte component | Responsibility (upstream) | Astrate package | Notes on the mapping |
91:### 1.3 Package layout (top-level skeleton)
124:### 1.4 Concurrency model (replacing RabbitMQ)
136:  to the device exactly as a slow AMQP consumer would in upstream Astarte. QoS 0 messages are
154:reproduce the in-process contract: per-device ordering (§1.4's shard guarantee), deferred-ack
421:values are *not* semantically interpreted (matching upstream behaviour) but are length- and
427:   `name:major`; otherwise reject (and raise the upstream-equivalent
441:     and ignored, matching upstream leniency.
465:Reference: the *Astarte MQTT v1 Protocol* specification in the upstream docs. Note that the
495:  upstream.
589:  "initial_payload_format": "json"}}` — additive, ignored-by-upstream-shaped).
606:   upstream SDK behaviour is byte-identical.
616:Mounted on one HTTP listener (`:8080`) with upstream-compatible base paths (configurable, so
624:  `{"errors": {"detail": "..."}}` with upstream status codes.
628:  avoid colliding with upstream's namespace.
657:  anchored as upstream does; `exp` honoured if present; `iat` not required (parity).
662:  surfaces. Measured, not inferred — see `test/conformance/upstream/channels.json` and
681:  possession, exactly as upstream treats it.
684:  always-online CRL equivalent — stricter than upstream's CRL-less default, toggleable via
703:  secret (parity); after first credentials request it conflicts (`422`, upstream-shaped error).
797:- DB outage: shards park with exponential backoff; broker applies backpressure (§1.4); QoS 0
823:This document freezes: the module decomposition (§1.3), the storage model (§2), the wire
=== COMPAT deviations count ===
375

$ grep -n "Status: not started\|not started" .mule/for-giulio.md | head -10; echo "=== milestones.md status mentions ==="; grep -rn "milestones.md" .mule/for-giulio.md | grep -i "stale\|not started\|status" | head -10
143:- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Status: not started", which is no longer true.** Release sweep (`gh api repos/astarte-platform/astarte/releases`, run this session) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.5` (2026-08-20, prerelease) — byte-identical to what the 2026-09-24 run saw, so nothing shipped upstream in the two days since. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Open non-alarm issues remain exactly **#93** (aclhook comment, `mule-review`, its own recipe path), **#92** (keyAgreement, parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** (FDO, `milestone-4.0`), **#1** (untouched per standing instruction). No new gaps, **no issues filed, no `.mule/todo.md` task lines added** — there is nothing machine-checkable left, and filing anything here would duplicate #92/#93/#78.
144:  **New, and the reason this run is not a copy of the last thirteen: `.mule/milestones.md:127-128` still reads "Status: not started. First recipe job: triage #47–#89 into an ordered plan (which are audits vs features vs decisions), file sub-issues where work splits, escalate the 'decide' set in one batch."** All of that happened: of #47–#89 every issue is closed except the deliberately-parked #78. The section describes as future work a triage that is finished, and it is the exact text every future run of this recipe reads to pick its target — so the recipe has been re-deriving "the backlog is done, the gate is upstream" from scratch on a timer for three weeks. Suggest rewriting those two lines to state the real position: *1.3 surface delivered (#47–#89 closed); milestone blocked on upstream v1.4.0 going stable; final phase = reconcile the two `docs/UPSTREAM-EXPERIMENTAL.md` rows, answer #92, bump `APICompatVersion`.* That also makes the blocking condition greppable instead of inferred.
=== milestones.md status mentions ===
143:- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Status: not started", which is no longer true.** Release sweep (`gh api repos/astarte-platform/astarte/releases`, run this session) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.5` (2026-08-20, prerelease) — byte-identical to what the 2026-09-24 run saw, so nothing shipped upstream in the two days since. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Open non-alarm issues remain exactly **#93** (aclhook comment, `mule-review`, its own recipe path), **#92** (keyAgreement, parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** (FDO, `milestone-4.0`), **#1** (untouched per standing instruction). No new gaps, **no issues filed, no `.mule/todo.md` task lines added** — there is nothing machine-checkable left, and filing anything here would duplicate #92/#93/#78.
144:  **New, and the reason this run is not a copy of the last thirteen: `.mule/milestones.md:127-128` still reads "Status: not started. First recipe job: triage #47–#89 into an ordered plan (which are audits vs features vs decisions), file sub-issues where work splits, escalate the 'decide' set in one batch."** All of that happened: of #47–#89 every issue is closed except the deliberately-parked #78. The section describes as future work a triage that is finished, and it is the exact text every future run of this recipe reads to pick its target — so the recipe has been re-deriving "the backlog is done, the gate is upstream" from scratch on a timer for three weeks. Suggest rewriting those two lines to state the real position: *1.3 surface delivered (#47–#89 closed); milestone blocked on upstream v1.4.0 going stable; final phase = reconcile the two `docs/UPSTREAM-EXPERIMENTAL.md` rows, answer #92, bump `APICompatVersion`.* That also makes the blocking condition greppable instead of inferred.

$ grep -n "verdict" tools/mule.sh | head -20
282:  local slug="$1" verdict="$2" elapsed="$3" outlog="$4"
285:    echo "verdict: $verdict"
290:  } > "$MULE/reports/log/$(date -u '+%Y%m%dT%H%M%SZ')-$slug-$verdict.md"
499:  local verdict="" reason=""
501:    verdict=blocked; reason="TIMEOUT after ${elapsed}s — task too big, split it"
513:    verdict=blocked; reason="opencode exited $rc"
515:    verdict=blocked; reason="touched a never-touch path"
545:    verdict=blocked; reason="wrote nothing"
547:    verdict=blocked; reason="${GATE_REASON:-gates failed}"
549:    verdict=done
552:  if [ "$verdict" = done ]; then

$ grep -n "proposed\|noop\|recipe" tools/mule.sh | sed -n '1,60p'
18:#   mule.sh recipe <name>    run a recipe: proposes new tasks, executes none
328:            issue_done refill_from_recipe cmd_next cmd_refill cmd_review cmd_budget \
336:    || note "gh not authenticated — the GitHub-issues recipe will not work"
404:  # approved; a standing line is either a gate or something the mule proposed to itself, and
405:  # those must not run ahead of approved work — twelve self-proposed lines sat in front of two
612:#   - **an empty queue does nothing at all.** It never runs a recipe to invent work for
636:  # to the daily budget first, so a recipe that proposes nothing still costs a tick and a
643:    # Whatever it just proposed waits for the next tick. That is deliberate: the proposals
689:# Ask the mule to propose its own work, one recipe per call, rotating. Rotation matters: left
691:refill_from_recipe() {
702:  note "queue is empty — running the '$next' recipe to propose work"
703:  cmd_recipe "$next"
717:  refill_from_recipe
891:   Each option is a recipe in .mule/recipes/. Running one PROPOSES
897:   labelled 'mule' first, then these recipes in rotation. So the
918:# Run a recipe: a fresh mule session that reads the recipe and appends proposed tasks.
919:cmd_recipe() {
920:  local name="${1:-}"; [ -n "$name" ] || die "recipe needs a name (see mule.sh menu)"
930:  local f="$MULE/recipes/$name.md"; [ -f "$f" ] || die "no such recipe: $f"
932:  local prompt="Read .mule/MULE.md, then read .mule/recipes/$name.md and carry it out.
934:That recipe is a *proposal* job: its output is new task lines appended to .mule/todo.md,
938:If the recipe tells you to run 'gh issue create' as the deliverable itself (some do —
939:check the recipe), that is fine to actually run now: filing an issue is not a code change
944:new issues, and only when the recipe explicitly asks for it."
956:    save_report "recipe-$name" "timeout" "$elapsed" "$outlog"
958:      && git -C "$REPO" commit -q -m "mule: recipe $name timed out" >/dev/null 2>&1
959:    bad "recipe timed out"; return 1
961:  save_report "recipe-$name" "proposed" "$elapsed" "$outlog"
963:    && git -C "$REPO" commit -q -m "mule: recipe $name ran (${elapsed}s)" >/dev/null 2>&1
964:  echo; note "proposed queue is now:"; grep -E '^- \[ \] ' "$TODO" | tail -20
975:  recipe)    shift; cmd_recipe "$@";;

$ git status --short; echo "rc=$?"; echo "--- confirm gate facts ---"; rg -n "APICompatVersion\s*=" internal/realm/service.go; grep -c "^" docs/UPSTREAM-EXPERIMENTAL.md; gh issue list --label upstream-experimental --state all --limit 20 --json number,state --jq '.[] | "\(.number) \(.state)"'
rc=0
--- confirm gate facts ---
588:const APICompatVersion = "1.2.2"
25
92 OPEN
68 CLOSED
67 CLOSED

Ran the milestones recipe against `.mule/milestones.md`: first non-`DONE` section is **v3.0** (v2.0 is `DONE`), scope decided (1.2.2 → 1.3/1.4), so step 3's investigation applied.

**What I checked**
- Reference `.mule/research/upstream-parity-2026-08-22.md` is still absent (gitignored dir, not in any ref) — escalated 2026-09-06 and re-escalated 2026-10-03 with three options; still unanswered. Fell back to the ledger + live upstream API, as the prior 22 runs did.
- `gh issue list --label milestone-3.0 --state all` → empty; this milestone's ledger lives under `upstream-parity`, open set unchanged: **#93** (`mule-review`), **#92** (parked on a stable v1.4.0), **#78** (`milestone-4.0`, deliberately parked).
- Closed-set audit: all 43 `upstream-parity` closed issues have `stateReason=COMPLETED` — **none** dropped as `not_planned`, so no gap is hiding behind a closed issue.
- Upstream: no release since `v1.3.5` / `v1.4.0-rc.6`, and **zero** commits since the 2026-10-06 run (`commits?since=2026-10-06T14:11:29Z` → empty).
- Standing item: `docs/UPSTREAM-EXPERIMENTAL.md` still its two rows (#67, #68), both waiting on upstream 1.4 final; `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`). `#90` re-read — closed `COMPLETED` 2026-08-23 with the frozen "do not bump now" decision, and the register is itself declared the checklist, so filing a bump issue would duplicate it.
- Step 5 deliberately **not** taken: scope is 1.2.2 → 1.3/1.4 and only the 1.3 half is delivered — #92 unanswered, both rows unreconciled, version unbumped.

No new gap, so no `gh issue create`, no `.mule/todo.md` lines, no `.mule/for-giulio.md` entry — the pending 10-06 entry, the 10-03 lost-doc escalation and the 09-26 stale-status escalation already carry every decision here, and the file is a queue, not a log. Working tree untouched.

```
Done: v3.0 re-verified — still gated on a stable upstream v1.4.0; no new gaps, nothing filed or queued
Files: none
Verified: gh api releases/tags + commits?since -> pass (no release, 0 new commits); gh issue list milestone-3.0 / upstream-parity --state all + stateReason audit -> pass (43 closed all COMPLETED, 3 open all owned); rg APICompatVersion internal/realm/service.go:588 -> 1.2.2; docs/UPSTREAM-EXPERIMENTAL.md 25 lines / 2 rows unreconciled; git status -> clean
Unsure: scope measured against the ledger + live upstream, not the milestone's source document — .mule/research/upstream-parity-2026-08-22.md is unreachable and its escalation is still open
Follow-ups: none
```
