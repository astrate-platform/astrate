slug: recipe-milestones
verdict: proposed
at:  cbdd3e8
ran: 2026-10-06T14:11:29Z on DietPi in 327s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md

$ wc -l .mule/research/upstream-parity-2026-08-22.md .mule/todo.md .mule/for-giulio.md docs/UPSTREAM-EXPERIMENTAL.md 2>/dev/null; ls .mule/
281 .mule/todo.md
   624 .mule/for-giulio.md
    25 docs/UPSTREAM-EXPERIMENTAL.md
   930 total
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

$ ls -la .mule/research 2>&1; find . -name '*upstream-parity*' -not -path './.git/*' 2>/dev/null
ls: cannot access '.mule/research': No such file or directory

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

$ find . -path ./.git -prune -o -iname '*upstream*' -print 2>/dev/null; echo ---; ls .mule/reports .mule/tasks 2>/dev/null; echo ---; grep -rl "upstream-parity-2026-08-22" --include='*.md' . 2>/dev/null | head -20
./test/conformance/upstream
./.mule/recipes/astarte-upstream.md
./.mule/reports/log/20260918T113930Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260923T181743Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260913T171327Z-recipe-astarte-upstream-noop.md
./.mule/reports/log/20260922T194051Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260921T105709Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md
./.mule/reports/log/20260915T111354Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260927T203020Z-recipe-astarte-upstream-timeout.md
./.mule/reports/log/20260929T194353Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260924T194457Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20261001T111852Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260924T183158Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260929T180551Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260907T020359Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20261005T022910Z-recipe-astarte-upstream-timeout.md
./.mule/reports/log/20260915T193626Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260910T202949Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260918T114100Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260923T105200Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260920T113611Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260917T190703Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260906T183102Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20261003T183445Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260922T170456Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260921T020624Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260914T020552Z-recipe-astarte-upstream-noop.md
./.mule/reports/log/20261001T184353Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260930T191053Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260921T190635Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260930T171256Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260908T113422Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260918T194746Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260910T112515Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260913T204415Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260914T020608Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260909T184315Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260925T170557Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260910T202936Z-recipe-astarte-upstream-noop.md
./.mule/reports/log/20260916T184218Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260916T184159Z-recipe-astarte-upstream-noop.md
./.mule/reports/log/20260917T190650Z-recipe-astarte-upstream-noop.md
./.mule/reports/log/20260910T112545Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20261004T184052Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260905T184625Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20261001T203619Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260912T180735Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260929T165352Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260928T203150Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20261005T201617Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260926T173418Z-recipe-astarte-upstream-timeout.md
./.mule/reports/log/20260911T191541Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260928T020420Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260913T171351Z-recipe-astarte-upstream-proposed.md
./.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md
./docs/UPSTREAM-EXPERIMENTAL.md
---
.mule/reports:
hygiene.md
log

.mule/tasks:
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md
---
./.mule/for-giulio.md
./.mule/milestones.md
./.mule/reports/log/20261003T180922Z-recipe-github-issues-proposed.md
./.mule/reports/log/20261004T181112Z-recipe-hygiene-timeout.md
./.mule/reports/log/20261005T193452Z-recipe-milestones-proposed.md
./.mule/reports/log/20260918T193641Z-recipe-milestones-proposed.md
./.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md
./.mule/reports/log/20260906T180442Z-recipe-milestones-proposed.md
./.mule/reports/log/20260912T171256Z-recipe-milestones-proposed.md
./.mule/reports/log/20260923T174341Z-recipe-milestones-proposed.md
./.mule/reports/log/20260910T170149Z-recipe-code-review-proposed.md
./.mule/reports/log/20260916T175801Z-recipe-milestones-proposed.md
./.mule/reports/log/20260926T174259Z-recipe-code-review-proposed.md
./.mule/reports/log/20260908T193127Z-recipe-milestones-proposed.md
./.mule/reports/log/20260914T234550Z-recipe-milestones-proposed.md
./.mule/reports/log/20260923T174805Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260916T183922Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260924T192153Z-recipe-milestones-proposed.md
./.mule/reports/log/20260924T182956Z-recipe-github-issues-proposed.md
./.mule/reports/log/20260918T110039Z-recipe-milestones-proposed.md

$ ls -t .mule/reports/log/ | head -20; echo ---; git status --short | head; git branch --show-current; git log --oneline -5
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
20261004T190849Z-recipe-code-review-proposed.md
20261004T184052Z-recipe-astarte-upstream-proposed.md
---
mule/queue
cbdd3e8 mule: recipe hygiene timed out
230dd11 mule: recipe docs-sync timed out
02003c1 mule: recipe code-review ran (125s)
1e8d87f mule: refill the queue
28ffb8d mule: recipe astarte-upstream ran (187s)

→ Read .mule/reports/log/20261005T193452Z-recipe-milestones-proposed.md
→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ gh api repos/astarte-platform/astarte/releases --paginate --jq '.[] | "\(.tag_name) \(.published_at) prerelease=\(.prerelease)"' 2>&1 | head -8; echo "=== milestone-3.0 ==="; gh issue list --label milestone-3.0 --state all --limit 50 --json number 2>&1; echo "=== upstream-parity open ==="; gh issue list --label upstream-parity --state open --limit 50 --json number,title,labels --jq '.[] | "\(.number) \(.title[0:70]) | \([.labels[].name]|join(","))"'
v1.3.5 2026-10-05T14:11:44Z prerelease=false
v1.4.0-rc.6 2026-09-30T07:28:27Z prerelease=true
v1.3.4 2026-09-18T09:42:57Z prerelease=false
v1.4.0-rc.5 2026-08-20T12:28:04Z prerelease=true
v1.3.3 2026-08-07T14:33:58Z prerelease=false
v1.4.0-rc.4 2026-08-06T13:40:20Z prerelease=true
v1.4.0-rc.3 2026-07-31T10:48:39Z prerelease=true
v1.4.0-rc.2 2026-07-14T10:16:22Z prerelease=true
=== milestone-3.0 ===
[]
=== upstream-parity open ===
93 In internal/broker/aclhook.go the control/keyAgreement deny comment (l | mule-review,upstream-parity
92 keyAgreement: upstream published the full wire spec (082) — decide whe | enhancement,upstream-parity,upstream-experimental
78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 ca | enhancement,milestone-4.0,upstream-parity

$ git log --all --oneline -- '.mule/research/upstream-parity*' 2>&1 | head; echo "---"; git cat-file -e $(git rev-list --all -- '.mule/research/upstream-parity-2026-08-22.md' 2>/dev/null | head -1):.mule/research/upstream-parity-2026-08-22.md 2>&1 && echo EXISTS || echo "not in any ref"; echo "--- APICompat ---"; rg -n "APICompatVersion\s*=" internal/realm/service.go
---
fatal: path '.mule/research/upstream-parity-2026-08-22.md' does not exist (neither on disk nor in the index)
not in any ref
--- APICompat ---
588:const APICompatVersion = "1.2.2"

$ gh api "repos/astarte-platform/astarte/commits?since=2026-10-05T14:20:00Z&per_page=30" --jq '.[] | "\(.sha[0:8]) \(.commit.committer.date) \(.commit.message|split("\n")[0])"' 2>&1; echo "=== open non-alarm ==="; gh issue list --state open --limit 100 --json number,title,labels --jq '[.[] | select(([.labels[].name] | index("mule-alarm")) | not)] | "count=\(length)", (.[] | "\(.number) \(.title[0:60])")'
1da04832 2026-10-06T11:01:17Z chore: forward port release-1.4
a8a26d0f 2026-10-06T10:59:25Z chore: forward port release-1.3
d85e0ed0 2026-10-06T10:02:32Z docs: sync FDO workflow to the current behavior (#2241)
a6bd2c21 2026-10-06T09:59:16Z fix(fdo): reject owner key names not accepted by OpenBao (#2252)
a51ab2e1 2026-10-06T07:21:18Z refactor: Update Umbrella checks (#2244)
7628d2ed 2026-10-05T15:54:31Z test(generators): validate interfaces with Core changesets
ed1e99f4 2026-10-05T15:54:31Z test(generators): validate mappings with Core changesets
bb17455a 2026-10-05T15:54:31Z test(generators): validate policy handlers with Core changesets
ec020bf8 2026-10-05T15:54:31Z fix(generators): generate valid trigger policy names
=== open non-alarm ===
count=4
93 In internal/broker/aclhook.go the control/keyAgreement deny 
92 keyAgreement: upstream published the full wire spec (082) — 
78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milest
1 Provide an Open Source IoT Platform unironically, unironical

$ for sha in 1da04832 a8a26d0f a51ab2e1 7628d2ed ec020bf8; do echo "== $sha"; gh api repos/astarte-platform/astarte/commits/$sha --jq '.files[] | "\(.status) +\(.additions)/-\(.deletions) \(.filename)"' 2>&1 | head -25; done
== 1da04832
modified +2/-1 .credo.exs
modified +8/-1 .dockerignore
added +5/-0 .formatter.exs
modified +20/-11 .github/codecov.yml
removed +0/-17 .github/labeler.yml
removed +0/-39 .github/workflows/astarte-appengine-api-workflow.yaml
modified +1/-56 .github/workflows/astarte-apps-build-workflow.yaml
added +191/-0 .github/workflows/astarte-apps-workflow.yaml
modified +7/-7 .github/workflows/astarte-build-workflow.yaml
modified +1/-1 .github/workflows/astarte-dashboard-compatibility.yml
modified +1/-1 .github/workflows/astarte-dashboard-docker-build.yml
modified +1/-1 .github/workflows/astarte-dashboard-tests.yml
removed +0/-39 .github/workflows/astarte-data-updater-plant-workflow.yaml
added +46/-0 .github/workflows/astarte-e2e-lint-workflow.yaml
added +30/-0 .github/workflows/astarte-e2e-next-test-build-workflow.yaml
modified +82/-2 .github/workflows/astarte-end-to-end-test-workflow.yaml
removed +0/-39 .github/workflows/astarte-housekeeping-workflow.yaml
removed +0/-144 .github/workflows/astarte-libs-build-workflow.yaml
removed +0/-79 .github/workflows/astarte-libs-workflow.yml
removed +0/-39 .github/workflows/astarte-pairing-workflow.yaml
removed +0/-39 .github/workflows/astarte-realm-management-workflow.yaml
removed +0/-39 .github/workflows/astarte-trigger-engine-workflow.yaml
modified +4/-1 .github/workflows/docs-check.yml
modified +3/-1 .github/workflows/docs-workflow.yaml
modified +1/-0 .github/workflows/pr-container-build-workflow.yaml
== a8a26d0f
modified +6/-0 CHANGELOG.md
modified +2/-2 README.md
modified +1/-1 libs/astarte_adapters/mix.exs
modified +1/-1 libs/astarte_test_suite/mix.exs
== a51ab2e1
modified +0/-73 .github/workflows/astarte-apps-build-workflow.yaml
modified +72/-0 .github/workflows/astarte-apps-workflow.yaml
modified +1/-1 apps/astarte_adapters/mix.exs
modified +1/-1 apps/astarte_generators/mix.exs
modified +1/-1 apps/astarte_test_suite/mix.exs
modified +2/-1 mix.exs
== 7628d2ed
modified +11/-1 apps/astarte_generators/test/astarte/core/generators/interface_test.exs
== ec020bf8
modified +2/-1 apps/astarte_generators/lib/astarte/core/generators/triggers/policy.ex
modified +5/-2 apps/astarte_generators/test/astarte/core/generators/triggers/policy/policy_test.exs

$ gh api repos/astarte-platform/astarte/commits/1da04832 --jq '.files | length' 2>&1; gh api repos/astarte-platform/astarte/commits/1da04832 --jq '.files[] | select(.filename | test("apps/astarte_(realm_management|pairing|appengine_api|trigger_engine|housekeeping|data_updater_plant)"; "i")) | "\(.status) \(.filename)"' 2>&1 | head -30; echo "=== release-1.3 forward port changelog ==="; gh api repos/astarte-platform/astarte/commits/a8a26d0f --jq '.files[] | select(.filename=="CHANGELOG.md") | .patch' 2>&1
300
removed apps/astarte_appengine_api/config/ci.exs
removed apps/astarte_appengine_api/config/config.exs
removed apps/astarte_appengine_api/config/dev.exs
removed apps/astarte_appengine_api/config/prod.exs
removed apps/astarte_appengine_api/config/runtime.exs
removed apps/astarte_appengine_api/config/test.exs
modified apps/astarte_appengine_api/lib/astarte_appengine_api/application.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api/auth/queries.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api/device/device.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api/device/interface_value.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api/device/queries.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api/rooms/events_dispatcher.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api/rooms/room.ex
removed apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant.ex
removed apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/behaviour.ex
removed apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/delete_volatile_trigger/request_data.ex
removed apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/data_updater_plant/install_volatile_trigger/request_data.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api/rpc/vmq_plugin.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api/stats/queries.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/device_status_by_alias_controller.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/device_status_by_group_controller.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/device_status_controller.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/fallback_controller.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api_web/controllers/interface_values_by_group_controller.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api_web/telemetry.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api_web/telemetry/api_usage.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api_web/telemetry/database_events.ex
modified apps/astarte_appengine_api/lib/astarte_appengine_api_web/views/error_view.ex
modified apps/astarte_appengine_api/mix.exs
removed apps/astarte_appengine_api/mix.lock
=== release-1.3 forward port changelog ===
@@ -123,6 +123,12 @@ Includes changes from v1.3.4
 
 - [astarte_data_updater_plant] Use mississippi consumer for data updater processes
 
+## [1.3.5] - 2026-10-05
+
+### Fixed
+
+- [astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available
+
 ## [1.3.4] - 2026-09-17
 
 ### Fixed

→ Read .mule/for-giulio.md [limit=18]

$ gh api repos/astarte-platform/astarte/commits/1da04832 --jq '.files[] | select(.filename | test("_web/(router|channel|routes)|router.ex")) | .filename' 2>&1; echo "--- fdo #2241/#2252 files ---"; for sha in d85e0ed0 a6bd2c21; do echo "== $sha"; gh api repos/astarte-platform/astarte/commits/$sha --jq '.files[] | "\(.status) \(.filename)"' 2>&1 | head -10; done
apps/astarte_appengine_api/test/astarte_appengine_api_web/channels/rooms_channel_test.exs
--- fdo #2241/#2252 files ---
== d85e0ed0
modified doc/pages/architecture/050-pairing_mechanism.md
modified doc/pages/user/015-astarte_dashboard.md
modified doc/pages/user/035-register_device.md
== a6bd2c21
modified apps/astarte_pairing/test/astarte_pairing_web/controllers/owner_key_controller_test.exs
modified libs/astarte_secrets/lib/astarte_secrets/owner_key_initialization_options.ex
modified libs/astarte_secrets/test/astarte_secrets/core_test.exs
modified libs/astarte_secrets/test/astarte_secrets/owner_key_initialization_test.exs

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,13 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-10-06: v3.0 verified unchanged — still gated on upstream — and the master delta since the 2026-10-05 run is a closed set of 13 commits, every one wire-inert for Astrate. No issue filed, no `.mule/todo.md` line queued.** Release sweep is byte-identical to yesterday's (`gh api repos/astarte-platform/astarte/releases`): **no stable v1.4.0** — newest stable `v1.3.5` (2026-10-05), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition (your 2026-09-04 "wait for a stable v1.4.0"; an `-rc.N` does not satisfy it) stays unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty** — this milestone's whole ledger lives under `upstream-parity`, where the open set is unchanged: **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`); **#47–#89** otherwise all closed. Non-alarm open issues total exactly four (plus the 25 `mule-alarm` #94–#118): the three above and **#1** (untouched per standing instruction). `APICompatVersion` still `"1.2.2"` (`internal/realm/service.go:588`, re-checked).
+  **The new master commits — measured, and none is a gap.** `gh api .../commits?since=2026-10-05T14:20:00Z` returns 9 today (13 with the four landing between yesterday's run and that mark): `1da04832` + `a8a26d0f` "chore: forward port release-1.4 / release-1.3" are release-branch merges — CI/workflow consolidation, CHANGELOG, `mix.exs` bumps, and the appengine RPC refactor that *deletes* the `astarte_data_updater_plant` RPC client (Astrate implements no `data_updater_plant`); a router-level sweep of the 300-file forward-port finds **no new route or channel** — the only `_web/(router)` hit is `rooms_channel_test.exs`. `d85e0ed0` (docs) + `a6bd2c21` `fix(fdo): reject owner key names not accepted by OpenBao` are both FDO, i.e. **#78**'s parked body. `a51ab2e1` is an umbrella-CI refactor; `7628d2ed`/`ed1e99f4`/`bb17455a`/`ec020bf8` are `test(generators)`/`fix(generators)` Core-changeset validation — upstream test tooling, not Astrate surface. No HTTP route, MQTT topic, AMQP control message or interface-schema field moved.
+  **Standing item unchanged.** `docs/UPSTREAM-EXPERIMENTAL.md`'s two rows (#67 `required`/`encrypted` mapping fields, #68 `async_operation=false`) still wait on **upstream 1.4 final** for their promoted-vs-deprecated call; nothing in today's commit set promotes or removes either and none of it is released, so the register is untouched.
+  **Step 5 deliberately not taken**, same reason as the 21 runs before: the section's declared scope is **1.2.2 → 1.3/1.4** and only the 1.3 half is delivered. The 1.4 half stays open — #92 unanswered, both register rows unreconciled, `APICompatVersion` unbumped — so the tag must not be cut. Your call when a stable v1.4.0 ships: answer #92, reconcile the register, run the final-phase bump, then cut the tag.
+
 - **astarte-upstream recipe run, 2026-10-05 (evening): no upstream movement that touches anything Astrate implements, so nothing was queued — no `.mule/todo.md` line, no `gh issue create`, nothing commented, closed or edited on GitHub.** One thing *is* new since this recipe's last run (`2026-10-04T20:40:52+02:00` ran; the `2026-10-05T04:29:10+02:00` run timed out): stable **v1.3.5**, published `2026-10-05T14:11:44Z`. The milestone recipe already recorded it as inert today; I re-measured it rather than trusting that, because "a release note is a claim, not a fact" is this recipe's own rule. `gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5` is **4 commits, 35 files** — `8d47574f` + `4c4222cf` `chore: update horde`, `c4c63ebe` `chore(rm): bump xandra to 0.19.4`, `8f0f5ef1` `chore: prepare v1.3.5` — and **every app/lib file in it is a `mix.exs`/`mix.lock` version bump**. The rest is `CHANGELOG.md` (+6: one `### Fixed` entry, "[astarte_data_updater_plant] Increase RPC server reliability"), `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `README.md`, `astarte-dashboard/package.json`. Not one source file changed, so no route, MQTT topic, BSON encoding, control message or interface-schema field moved: Astrate implements no `data_updater_plant` and no Erlang RPC server. Wire-inert, same shape as v1.3.2/v1.3.3/v1.3.4. **Master since `v1.4.0-rc.6` is still the single commit `731bcddc`** ("chore: forward port release-1.4"; `gh api .../commits?since=2026-10-05T09:34:00Z` returns it and nothing newer), the migration-relocation commit the milestone run measured today as inert — so both wire-bearing post-rc.6 commits (`b6d46ad4`, `958bb9fe`) still have their owners in the queue and nothing new joins them.
   **The parked-decisions row for #92 is still unmet, re-checked against the tag list rather than the release list** (this host has no `.mule/waiting-on.md` — absent from `HEAD`, already reported; I read it out of `origin/main` read-only): `gh api repos/astarte-platform/astarte/tags?per_page=100` gives newest stable **v1.3.5** and newest overall **v1.4.0-rc.6**, with `v1.4.0-rc.0`…`rc.6` and no stable `v1.4.0`. The row asks for a **stable** tag, which an `-rc.N` does not satisfy, so the twenty-second run ends the same way as the twenty-one before it.
   **The one delta to your open COMPATIBILITY.md wording proposal (the v1.3.2 one, folding v1.3.3/v1.3.4, further down this file): its version-reference sentence should read "newest stable v1.3.5" where it currently says "v1.3.4". Nothing else in it changes** — v1.3.5 adds no capability, so the same decision and the same "not yet emulated" list stand. Not applied: `docs/COMPATIBILITY.md` is untouched and still targets **1.2.2**, `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked). I deliberately did **not** queue a fresh `compat-note-v1.3.5` line: the proposal is already open here, and a second line would re-derive it for the next run to find.
 

No new gaps, but the master delta since yesterday's run needs recording for Giulio. The existing entry style carries evidence inline, so I'll append one entry at the top.
$ git status --short; echo ---; sed -n '11,14p' .mule/for-giulio.md | cut -c1-120
M .mule/for-giulio.md
---
---

- **Milestone recipe run, 2026-10-06: v3.0 verified unchanged — still gated on upstream — and the master delta since
  **The new master commits — measured, and none is a gap.** `gh api .../commits?since=2026-10-05T14:20:00Z` returns 9

First non-DONE milestone is **v3.0**. Its scope is decided (1.2.2 → 1.3/1.4), so I ran the investigation: release sweep, `milestone-3.0` ledger, `upstream-parity` open set, UPSTREAM-EXPERIMENTAL reconcile, and diffed every new `master` commit since the last run. No new gap found — everything is either closed, already owned, or parked. The one genuinely new fact (today's 13 upstream master commits, all wire-inert) went into `.mule/for-giulio.md` as the deliverable.

Done: v3.0 re-verified — still gated on upstream v1.4.0; recorded today's wire-inert master delta for Giulio
Files: .mule/for-giulio.md
Verified: gh api releases / issue list milestone-3.0 + upstream-parity / commits sweep + 5 commit file-diffs -> pass
Unsure: nothing (reference doc .mule/research/… is absent from the working tree and every git ref; the ledger under `upstream-parity` + live upstream API served as the reference, per the prior 21 runs)
Follow-ups: none (no todo.md lines, no `gh issue create` — step 5 "cut the tag" deliberately not proposed, 1.4 gate unresolved)
