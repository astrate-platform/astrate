slug: recipe-code-review
verdict: proposed
at:  25137f2
ran: 2026-09-03T19:59:14Z on DietPi in 241s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
✱ Glob ".mule/reviews/*.md" 0 matches
→ Read .mule/log.md

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

$ rg -c '' --glob '*.go' /root/astrate-mule/internal/flow/ | sort -t: -k2 -rn | head -20
/root/astrate-mule/internal/flow/router_test.go:402
/root/astrate-mule/internal/flow/message.go:362
/root/astrate-mule/internal/flow/message_test.go:348
/root/astrate-mule/internal/flow/manager_test.go:313
/root/astrate-mule/internal/flow/pipeline_test.go:280
/root/astrate-mule/internal/flow/router.go:257
/root/astrate-mule/internal/flow/flow.go:245
/root/astrate-mule/internal/flow/flow_test.go:196
/root/astrate-mule/internal/flow/pipeline.go:170
/root/astrate-mule/internal/flow/block.go:91
/root/astrate-mule/internal/flow/graph.go:43

→ Read .mule/todo.md

$ for d in /root/astrate-mule/internal/*/; do pkg=$(basename "$d"); count=$(rg -c '' --glob '*.go' "$d" 2>/dev/null | awk -F: '{s+=$2}END{print s}'); echo "$count $pkg"; done | sort -rn
12413 engine
5348 store
4885 appengine
3441 broker
2777 pairing
2707 flow
1404 auth
1386 realm
697 testutil
647 config
500 housekeeping
357 httpx
272 observability
43 swagger

$ ls -la /root/astrate-mule/.mule/; echo "---for-giulio---"; cat /root/astrate-mule/.mule/for-giulio.md 2>/dev/null | head -80
total 1040
drwxr-xr-x  8 root root   4096 Sep  3 20:48 .
drwxr-xr-x 14 root root   4096 Aug 31 21:44 ..
-rw-r--r--  1 root root      0 Jul 28 22:45 .alarmed
-rw-r--r--  1 root root     14 Sep  3 21:55 .budget
-rw-r--r--  1 root root    364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root     11 Jul 28 07:36 .heartbeat
-rw-r--r--  1 root root      0 Sep  3 21:55 .last-output
-rw-r--r--  1 root root      0 Sep  3 21:55 .lock
-rw-r--r--  1 root root     12 Sep  3 21:55 .rotation
-rw-r--r--  1 root root   6429 Jul 27 22:18 MULE.md
-rw-r--r--  1 root root   2762 Aug 31 01:45 config
-rw-r--r--  1 root root 948612 Sep  3 21:55 cron.log
drwxr-xr-x  2 root root   4096 Sep  3 20:48 failed
-rw-r--r--  1 root root   9160 Sep  3 21:47 for-giulio.md
-rw-r--r--  1 root root   1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root   1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root   5622 Sep  3 20:48 log.md
-rw-r--r--  1 root root   2643 Jul 27 22:18 milestones.md
drwxr-xr-x  2 root root   4096 Jul 27 21:20 plans
drwxr-xr-x  2 root root   4096 Jul 27 22:54 recipes
drwxr-xr-x  3 root root   4096 Jul 28 05:50 reports
drwxr-xr-x  2 root root   4096 Sep  3 13:26 reviews
drwxr-xr-x  2 root root   4096 Aug  3 04:11 tasks
-rw-r--r--  1 root root  20126 Sep  3 20:48 todo.md
---for-giulio---
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **milestone 2.0 looks complete, verify and cut the tag** — all 11 `milestone-2.0` issues
  CLOSED (#23–#27, #37, #39–#43), no open issues, no new gaps after re-checking upstream
  astarte_flow block catalog against `internal/flow/` + git log (MQTT/HTTP source/sink,
  json_path_map, pure-transform set, virtual_device_pool, container block MVP, flow API,
  durable named flows all landed). `mule/queue` branch is stale — see entry below.
  (Milestones recipe run, 2026-09-03.)

- **v3.0 decide #68** (`upstream-parity`, `mule-blocked`): `async_operation=false` acceptance is already decided (Astrate stays always-sync; accept and ignore on mutating endpoints: housekeeping realm create/delete, realm-management interface install/update/delete, trigger/policy delete in `internal/realm/http.go`). The mule's code+test task (`issue-68-async-operation-accepted`) failed gates and is parked. Needs re-queue or gate fix.

---

- **v3.0 decide #67** (`upstream-parity`): `required` and `encrypted` mapping fields exist on upstream master but not at 1.2.2/1.3.0 (`required` = missing key at runtime → `missing_required_mapping` on object-aggregated mappings; `encrypted` ties into keyAgreement #51). Decide now: reject (strict 1.2.2), ignore-and-accept, or implement `required` early. Issue is OPEN, no decision on record, not yet escalated here.

---

- **v3.0 decide #91** (`upstream-parity`, `mule-blocked`): upstream serves the unauthenticated health check at `GET /pairing/health` and returns 404 for `GET /pairing/v1/{realm}/health`, but Astrate only implemented the `v1/{realm}` shape (issue #71, commit 5111e3c) — the dashboard polls per-service health at the service root, so the current route may be invisible to it. Decide: move/duplicate the handler at `GET /pairing/health` (keep or drop the v1/{realm} variant), and re-check the original realmcfg-02 measurement (a port override may have skewed it).

---

- **`govulncheck` is not installed on the Pi**, so the hygiene recipe's highest-priority check (reachable dependency vulnerabilities) could not run here. `go install golang.org/x/vuln/cmd/govulncheck@latest` would close it — I did not pick a version because the pin is a decision. The `go list -m -u` sweep shows only version-skew on transitive deps (cloud/azure SDKs, cel/expr, clickhouse); none corresponds to an identified fix this repo needs, so I proposed no bumps. (Hygiene recipe run, 2026-09-02.)

---

- **`docs/site/appengine-api.md:87-88` documents `GET` and `DELETE /appengine/v1/<realm>/groups/<name>`**, but no such routes exist in `internal/appengine/http.go` (only `/groups/{group}/devices` and the nested device ops — there is no `GET` or `DELETE` on `/groups/{group}`). The endpoint the page describes is absent from the code. Decide: drop the two lines, mark them not-yet-implemented, or have the code grow them. (Docs-sync recipe run, 2026-09-02.)

---

- **`Router.Submit` TOCTOU on `closed` flag** (`internal/flow/router.go:113-120`): Submit
  reads `r.closed` under the mutex, drops the lock, then sends on the channel. A concurrent
  `Drain` could close the channel between the unlock and the send, causing a send-on-closed
  channel panic in the caller (not recovered by `processOne`'s defer). Fix options: (a) hold
  the lock during channel send (hot-path perf hit), (b) use a select with `r.quit` instead of
  sending on a closed channel, (c) accept the risk and document that Submit must not be called
  after Drain. Needs a design call — not proposing a task. (Code review 2026-09-02.)

---

- **#87 `lua_map` — needs embedded Lua runtime, parked.** Design/implementation decision: embedding a Lua VM in Astrate is not machine-checkable by the mule. Consider closing if Lua flow support is not on any active roadmap.
- **#78 FDO device onboarding — milestone-4.0, investigation phase.** Too large for a single mule task; the investigation work (reading upstream's TO2 handling, inventorying endpoints/schema/keys) is a multi-session project. Parking for now until the v3.0 queue clears and this becomes the next milestone target.
- **#1 stale — "Provide an Open Source IoT Platform unironically" (wontfix).** Has been open since the repo's founding with no activity. Consider closing.

---

- **Milestones recipe has nothing to file: v2.0 is already DONE, and `mule/queue` is 120
  commits behind `origin/main`, so this branch's `.mule/milestones.md` is stale.** My branch
  still shows v2.0 "not started" (and v3.0 as CLEA TBD), so the recipe pointed me at v2.0 —
  but on `main` v2.0 is marked **DONE (2026-07-29)**: every `milestone-2.0` issue is CLOSED
  (#23–#27, #37, #39, #40, #41, #42, #43), and the factory/catalog/container-block/durable
  named flows/`/flow/v1` API all live on `main` (`git show origin/main:.mule/milestones.md`).
  `gh issue list --label milestone-2.0 --state open` is empty and there are no new v2.0 gaps
  to file, so per the recipe's step 5 I filed **nothing** rather than duplicate closed work.
  **Real current target is v3.0** ("upstream Astarte feature parity 1.3/1.4", scoped on main
  2026-08-22): the first recipe job there is to triage `upstream-parity` backlog **#47–#89**
  (source `.mule/research/upstream-parity-2026-08-22.md`) into a plan. But I should not start
  that triage from this stale branch. Action needed: **merge/rebuild `mule/queue` on top of
  `origin/main`**, then the next milestones run targets v3.0 with the correct research doc.
  (Milestones recipe run, 2026-09-01.)

---

- ~~`device_deletion_started`/`device_deletion_finished` trigger events are not emitted~~ —
  **decided 2026-07-27: emit both, back-to-back, around the synchronous delete.** Filed as
  issue #21 (`mule`). (Cross-project survey, 2026-07-27,
  `.mule/research/survey-2026-07-27.md` source 4.)
- ~~Mustache trigger-action templates are accepted but not rendered~~ — **decided
  2026-07-27: implement it.** Guiding principle clarified: Astarte compatibility means
  SDK/wire compatibility, not minimum dependency count — Astrate is allowed to be a
  compatible *superset*. Library picked: `github.com/cbroglie/mustache`. Filed as issue #22

The reviews directory doesn't exist yet (glob found nothing). Let me check what review areas were done recently and pick an area that wasn't recently reviewed. Let me check git history for review activity and check the `.mule/for-giulio.md` too.

Let me pick an area. The log shows recent work focused on flow, engine, store, housekeeping, pairing, realm. No code-review entries in log. Since reviews/ is empty, this may be the first review. Let me pick a package not heavily reviewed recently.

The recipe says pick one package and go deep. Given the size (engine is largest at 12413 lines — too big), let me pick something moderate. Let me check what was recently worked on. Recent work themes: flow, engine, store, housekeeping. Let me pick `realm` (1386 lines) — moderate and not heavily touched recently in the log, or `auth` or `observability`.

Let me look at the log to see which areas recent tasks touched. Recent tasks: pairing, engine, realm, flow, housekeeping, store. Let me pick a package not in the recent rotation — `auth` (1404 lines) or `observability` (272) or `config` (647) or `broker` (3441).

`observability` is small. But let me consider. The recipe warns a review that skims finds nothing. Let me pick `auth` — moderately sized, and I don't see recent auth-specific work in the log.

Let me check `.mule/for-giulio.md` and check git log for recent review-related commits.
$ cd /root/astrate-mule && rg -l 'review' .mule/ 2>/dev/null; echo "---reviews dir---"; ls -la .mule/reviews/ 2>/dev/null || echo "no reviews dir"
.mule/cron.log
.mule/recipes/code-review.md
.mule/reports/issue-13.md
.mule/reviews/engine-2026-09-03.md
.mule/reviews/flow-2026-09-02.md
.mule/reviews/realm-2026-09-01.md
.mule/reviews/store-2026-07-27.md
.mule/for-giulio.md
.mule/reports/log/20260824T021445Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260903T194732Z-recipe-github-issues-proposed.md
.mule/reports/log/20260902T203149Z-recipe-milestones-proposed.md
.mule/reports/log/20260903T110658Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260902T191445Z-docs-sync-appengine-422-blocked.md
.mule/reports/log/20260902T112257Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260831T184221Z-compat-note-v140-rc3-blocked.md
.mule/reports/log/20260728T035046Z-issue-13-checked.md
.mule/reports/log/20260901T205135Z-recipe-milestones-proposed.md
.mule/reports/log/20260901T114404Z-recipe-milestones-proposed.md
.mule/reports/log/20260831T172139Z-probe-interface-default-values-blocked.md
.mule/reports/log/20260901T185332Z-realm-policy-list-sorted-blocked.md
.mule/reports/log/20260728T053605Z-issue-13-checked.md
.mule/reports/log/20260901T175716Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260901T184704Z-recipe-code-review-proposed.md
.mule/reports/log/20260803T021148Z-recipe-astarte-upstream-proposed.md
.mule/todo.md
.mule/reports/log/20260901T183224Z-compat-note-v140-rc5-blocked.md
.mule/reports/log/20260903T192902Z-recipe-milestones-proposed.md
.mule/reports/log/20260902T111801Z-recipe-github-issues-proposed.md
.mule/reports/log/20260901T200632Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260902T170345Z-recipe-code-review-proposed.md
.mule/reports/log/20260810T021151Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260901T170200Z-recipe-github-issues-proposed.md
.mule/reports/log/20260831T114758Z-issue-68-blocked.md
.mule/reports/log/20260901T173234Z-issue-91-pairing-health-serve-root-blocked.md
.mule/reports/log/20260831T190559Z-compat-note-v133-blocked.md
.mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260903T112703Z-recipe-code-review-proposed.md
---reviews dir---
total 40
drwxr-xr-x 2 root root 4096 Sep  3 13:26 .
drwxr-xr-x 8 root root 4096 Sep  3 20:48 ..
-rw-r--r-- 1 root root 7300 Sep  3 13:26 engine-2026-09-03.md
-rw-r--r-- 1 root root 5793 Sep  2 19:01 flow-2026-09-02.md
-rw-r--r-- 1 root root 4291 Sep  1 20:46 realm-2026-09-01.md
-rw-r--r-- 1 root root 4406 Jul 27 14:57 store-2026-07-27.md

The reviews directory exists. Areas already reviewed: `engine` (09-03), `flow` (09-02), `realm` (09-01), `store` (07-27). I should pick an area NOT recently reviewed. Let me look at the recent code-review proposal to see what was proposed and understand past picks. Let me look at the for-giulio code-review entries and the existing reviews to avoid repeating.

Areas not yet reviewed: `appengine`, `broker`, `auth`, `observability`, `config`, `housekeeping`, `pairing`, `httpx`, and the `pkg/` packages.

Let me pick a moderate package not reviewed. `pairing` (2777) or `broker` (3441) or `auth` (1404). The recent code-review log ran on 09-03 (engine). Let me check the last code-review proposal log to see what areas it considered.

Let me look at the recent code-review proposal logs to understand the rotation and avoid duplication.
$ cd /root/astrate-mule && cat .mule/reports/log/20260903T112703Z-recipe-code-review-proposed.md
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

$ cd /root/astrate-mule && cat .mule/reviews/engine-2026-09-03.md
# Code review: internal/engine — 2026-09-03

## Files read

- `internal/engine/cache.go` (496 lines) — schema/device caches, copy-on-write snapshot
- `internal/engine/router.go` (490 lines) — sharded pipeline, Submit, metrics
- `internal/engine/batch.go` (406 lines) — micro-batcher, flush, row conversion
- `internal/engine/data.go` (386 lines) — §2.6 validation pipeline, handle, reject
- `internal/engine/engine.go` (318 lines) — New, Drain, fireCommitted, lifecycle
- `internal/engine/control.go` (315 lines) — emptyCache, producer/properties, consumer/properties
- `internal/engine/serverdata.go` (280 lines) — PublishServerValue, UnsetServerProperty
- `internal/engine/introspection.go` (168 lines) — handleIntrospection, parseIntrospection
- `internal/engine/capabilities.go` (101 lines) — BSON capabilities parsing
- `internal/engine/topics.go` (97 lines) — splitDeviceTopic, classify, matchInterface
- `internal/engine/forward/http.go` (125 lines) — HTTP trigger forwarder

Plus test files: `data_test.go`, `serverdata_test.go`, `control_test.go`, `cache_test.go`, `router_test.go`, `fakestore_test.go`, `introspection_test.go`, `capabilities_test.go`, `bus_events_test.go`, `batch_test.go`.

## What I read in detail

I focused on the M6b seam wiring — how `PublishServerValue` and the control-channel handlers interact with the batch/trigger/bus subsystem — because the code's comments consistently reference the `afterCommit` hook (engine.go:109, batch.go:168) as the single observation point for committed ops. I verified the full call graph from `PublishServerValue` through persistence to broker publish, and from `finalize` through `afterCommit` to `fireCommitted`/`fireData`.

## Findings

### 1. `PublishServerValue` bypasses `afterCommit` — no trigger or live bus events

**Where:** `internal/engine/serverdata.go:99-162`

The batch path persists ops, then `finalize` (batch.go:168) calls `afterCommit`, which is wired to `fireCommitted` (engine.go:139) and fires data triggers + live bus events.

`PublishServerValue` bypasses the batch entirely: it persists directly (`UpsertProperty` / `AppendDatastreams`), then publishes to the broker. It never calls `afterCommit`. Consequences:

- Server-initiated property/datastream writes produce no `incoming_data` live bus events — a Channels viewer watching the stream will not see server-side writes.
- No data triggers fire for server-owned interface paths (the `triggers.DataEvent` matching loop in `fireData` is never reached).

This may be intentional — server writes are admin-initiated, not device-initiated — but it is undocumented and untested. The batch path's `TestBatchTriggers` proves triggers fire for device-originated ops; no equivalent test exists for server-originated ops. At minimum the live-bus gap should be documented or deliberately excluded.

**Proposed task:** in `internal/engine/serverdata_test.go`, add a test that installs a data trigger matching a server-owned interface, calls `PublishServerValue`, and asserts the trigger was NOT fired (proving the documented exclusion), or DOES fire (if the intended behaviour is parity with the data path). The test should also verify whether a Channels live-bus subscriber sees the event.

### 2. `sendConsumerProperties` zlib fallback when device not cached

**Where:** `internal/engine/control.go:204`

```go
compression := compressionZlib
if dev := e.devices.peek(realm.name, id); dev != nil {
    compression = purgeCompressionFor(dev)
}
```

When `peek` returns nil (device not cached — e.g. after an eviction), compression defaults to zlib. A device that declared `purge_properties_compression_format: plaintext` via the capabilities topic, then disconnected and reconnected, may briefly have no cache entry when `sendConsumerProperties` runs on the connect lifecycle event. The purge message arrives compressed, which a strictly-plaintext device cannot decompress.

The cache is populated lazily on first data message, not on connect — so there is always a window between connect and the first data message where the device is uncached. In practice this window is tiny and most devices support zlib, but the fallback is asymmetric with the rest of the code which respects the declared format.

**Proposed task:** in `internal/engine/capabilities_test.go`, add a test that creates a plaintext-declared device, evicts its cache entry, then calls `sendConsumerProperties` and asserts the purge payload is uncompressed (plaintext). This will currently fail, confirming the bug. A follow-up fix would load the device from the store when `peek` returns nil, or persist the capability in the store.

### 3. `deflateProperties` does not cap entry count — only inflated size

**Where:** `internal/engine/control.go:262-274`

The sender-side `deflateProperties` accepts an unbounded `[]string` and joins + compresses it. The receiver-side `inflateProperties` caps inflated size at `maxControlInflated` (1 MiB). A pathological sender (or a bug in `sendConsumerProperties` accumulating a huge list) could produce a payload that inflates to 1 MiB of semicolons — well-formed but wasteful. Adding an explicit entry-count cap on the sender side (e.g. 100,000) would be a cheap defense-in-depth.

This is low priority — the sender is always internal — but it is a classic "missing guard on a loop that grows from data" pattern.

**Proposed task:** not worth a standalone task — fold a `maxConsumerPropertyEntries` constant check into the existing `sendConsumerProperties` function and add a test for the cap.

### 4. `objectPathOK` checks only the first ObjectLeaves entry

**Where:** `internal/engine/data.go:330-336`

```go
func objectPathOK(ci *interfaceschema.CompiledInterface, prefix string) bool {
    for leaf := range ci.ObjectLeaves {
        _, ok := ci.Trie.Match(prefix + "/" + leaf)
        return ok
    }
    return false
}
```

This iterates `ObjectLeaves` but returns after checking the first entry. For all Astarte interfaces today, all leaves share the same prefix structure, so this is correct. However, the function name (`objectPathOK`) implies full validation. If a future interface had leaves at different depths (e.g. `/setpoints/heating` and `/setpoints/zone/a/b`), only the first leaf would be checked. The iteration is a no-op; it should either check all leaves or document why one suffices.

**Proposed task:** clarify with a comment that one leaf suffices because upstream validation guarantees uniform prefix depth, or change to check all leaves (both approaches are safe). Low priority.

## What I decided NOT to propose

- **Renames/reformatting** — the code is idiomatic Go throughout; no stutter, no naming issues.
- **The `noteFormat` lock discipline** — I verified that `dev.mu` is held only for the read-modify-write of `resetHintOnBSON`/`formatHint` and released before the DB call. Correct.
- **The `deviceCache.get` TOCTOU** — two goroutines loading the same device concurrently is impossible by the per-device shard design. Correct.
- **The `backoff` loop in `retryStore`/`deviceState`** — exponential backoff with cap and quit/channel check is textbook. Correct.
- **Anything in `triggers/` or `stream/`** — they are separate packages and not the focus of this review.

$ cd /root/astrate-mule && for f in internal/pairing/*.go internal/pairing/**/*.go; do [ -f "$f" ] && echo "$(rg -c '' "$f") $f"; done 2>/dev/null | sort -rn
632 internal/pairing/service_test.go
614 internal/pairing/http_test.go
427 internal/pairing/service.go
348 internal/pairing/http.go
309 internal/pairing/ca/ca_test.go
287 internal/pairing/ca/ca.go
81 internal/pairing/ratelimit_test.go
79 internal/pairing/ratelimit.go

Now I understand the rotation. Reviewed so far: `store` (07-27), `realm` (09-01), `flow` (09-02), `engine` (09-03). I need to pick an area NOT reviewed. Let me look at the other two code-review proposal logs to see what was covered on 09-01 and 09-02.

Let me check the flow and realm reviews briefly and the 09-02 code-review log to be thorough about what was covered, then pick a new area.

Given the rotation, good candidates not yet reviewed: `appengine` (4885), `broker` (3441), `pairing` (2777), `auth` (1404), `housekeeping` (500), `httpx` (357), `config` (647), `observability` (272), and `pkg/` packages.

I'll pick `pairing` or `broker`. Both are substantial and untouched by reviews. Let me pick one. The broker is core to a message platform. But pairing (2777) is moderate. Let me go with `broker` since it's a larger unreviewed package, or actually let me consider the size/budget — my context budget matters. The code-review recipe warns against skimming.

Let me look at the other review logs to avoid picking something that was partially covered, but given reviews show only flow/realm/engine/store, all other areas are open. Let me pick `pairing` (2777 lines) as a good moderate area not yet reviewed.

Actually, let me reconsider - `broker` at 3441 is the bigger prize and untouched. But given context budget limits, `pairing` at 2777 is more feasible. Let me go with `pairing`.

Let me start reading the pairing package.
$ cd /root/astrate-mule && grep -n 'code-review\|pairing\|broker\|auth\|housekeeping\|observability\|httpx\|config\|appengine' .mule/reviews/*.md .mule/for-giulio.md 2>/dev/null | head -40
.mule/reviews/engine-2026-09-03.md:21:I focused on the M6b seam wiring — how `PublishServerValue` and the control-channel handlers interact with the batch/trigger/bus subsystem — because the code's comments consistently reference the `afterCommit` hook (engine.go:109, batch.go:168) as the single observation point for committed ops. I verified the full call graph from `PublishServerValue` through persistence to broker publish, and from `finalize` through `afterCommit` to `fireCommitted`/`fireData`.
.mule/reviews/engine-2026-09-03.md:31:`PublishServerValue` bypasses the batch entirely: it persists directly (`UpsertProperty` / `AppendDatastreams`), then publishes to the broker. It never calls `afterCommit`. Consequences:
.mule/for-giulio.md:20:- **v3.0 decide #68** (`upstream-parity`, `mule-blocked`): `async_operation=false` acceptance is already decided (Astrate stays always-sync; accept and ignore on mutating endpoints: housekeeping realm create/delete, realm-management interface install/update/delete, trigger/policy delete in `internal/realm/http.go`). The mule's code+test task (`issue-68-async-operation-accepted`) failed gates and is parked. Needs re-queue or gate fix.
.mule/for-giulio.md:28:- **v3.0 decide #91** (`upstream-parity`, `mule-blocked`): upstream serves the unauthenticated health check at `GET /pairing/health` and returns 404 for `GET /pairing/v1/{realm}/health`, but Astrate only implemented the `v1/{realm}` shape (issue #71, commit 5111e3c) — the dashboard polls per-service health at the service root, so the current route may be invisible to it. Decide: move/duplicate the handler at `GET /pairing/health` (keep or drop the v1/{realm} variant), and re-check the original realmcfg-02 measurement (a port override may have skewed it).
.mule/for-giulio.md:36:- **`docs/site/appengine-api.md:87-88` documents `GET` and `DELETE /appengine/v1/<realm>/groups/<name>`**, but no such routes exist in `internal/appengine/http.go` (only `/groups/{group}/devices` and the nested device ops — there is no `GET` or `DELETE` on `/groups/{group}`). The endpoint the page describes is absent from the code. Decide: drop the two lines, mark them not-yet-implemented, or have the code grow them. (Docs-sync recipe run, 2026-09-02.)

`pairing` is unreviewed. It's a moderate package (service.go 427, http.go 348, plus ca and ratelimit). Let me read the main source files. I'll start with the largest source files: service.go and http.go, then ca.go and ratelimit.go.

Let me read them.
→ Read internal/pairing/service.go
→ Read internal/pairing/ratelimit.go
Now let me read http.go and the ca package.
→ Read internal/pairing/http.go
→ Read internal/pairing/ca/ca.go
Now let me look at the tests to understand what's covered. Let me look at the tests. Let me check the ca_test and service_test for coverage of key rules. Let me focus on the most valuable review areas.

Key things I've noticed so far — let me think about candidates:

1. **`Credentials` vs `Info` inhibited behavior asymmetry**: `Credentials` rejects inhibited (line 249-251 → ErrInhibited), `VerifyCredentials` rejects inhibited (line 326-328), but `Info` allows inhibited to read info (line 285-297, comment says that's intended). OK, documented.

2. **`Register` device registration limit**: gates every attempt including secret rotations. Documented as upstream parity. Interesting but there's a potential concern: `CountDevices` counts all devices in realm, not pending/unconfirmed. If limit is hit, even re-registering an existing device with a rotated secret fails. But that's documented upstream parity.

3. **Rate limiter `Allow`** — the `handleCredentials` calls `a.credLimiter.Allow` twice, consuming TWO tokens from the IP bucket AND one from the dev bucket per request. That's intentional (stronger limit).

4. **`ratelimit.go` `Allow` bug**: Look at line 65-72. When `elapsed > 0`, tokens are added and `b.last = now`. But if `elapsed == 0` (two requests in the same instant), tokens aren't refilled and `b.last` stays. That's fine. But there's a subtle issue: the refill logic only adds `elapsed * l.rate` and updates `b.last = now` inside `if elapsed > 0`. This is correct behavior.

Wait — actually let me reconsider a potential real bug in `ratelimit.go`. The comment on `lru.New` says "lru.New only fails for sizes < 1". That's correct. Not a bug.

5. **`remoteIP`** — uses only the peer IP from RemoteAddr. No X-Forwarded-For handling. Behind a proxy, all clients appear as the proxy IP, defeating per-IP rate limiting. But this may be intentional (Astrate listens direct). Could be worth noting as a design decision, but not a clear bug.

Let me look at whether there's a real behavioral gap. Let me examine the test files to find untested rules. Let me look at the service_test and http_test.

Let me focus and look at specific things. The most valuable find would be an untested existing rule. Let me examine the `Info` status logic and the `Unregister` — note that `handleUnregister` calls `Unregister` which returns `store.ErrNotFound` for unknown device → writeServiceError maps to `WriteDeviceNotFound` (404). OK covered.

Let me look at the register limit interaction more carefully and check tests. Let me grep the test files.
→ Read internal/pairing/service_test.go

$ cd /root/astrate-mule && grep -rn "4.5\|rate.limit\|ratelimit\|Limiter\|bcrypt\|DoS\|rate limit" docs/DESIGN.md | head -30
220:    credentials_secret_hash text NOT NULL,      -- bcrypt
551:{ "v": <value>, "t": "2026-06-10T12:34:56.789Z" }
627:   *only* against the Pairing API. Stored bcrypt-hashed (cost 10); shown exactly once.
701:Constant-time-ish secret verification (bcrypt compare against the stored hash; uniform error +
702:per-IP/device token-bucket rate limiting to blunt brute force). First successful call stamps
745:### 4.5 Platform hardening checklist (v1 scope)
749:- Secrets handling: credentials secrets bcrypt-hashed; CA keys AES-GCM-encrypted; JWT public

Good coverage. Now let me consider what's interesting. Let me note a few findings and think about which are worth proposing.

**Finding candidates:**

1. **`Info` does NOT reject inhibited devices, but `Credentials` and `Verify` do.** Documented as intentional (inhibited may read info). Not a bug.

2. **The `Register` device_registration_limit counts ALL devices, not just pending ones.** `CountDevices` returns total device count. When a realm limit is set, registering a second device after the first has already gotten credentials fails. This is documented upstream parity. Not a bug.

3. **`Unregister` on `handleUnregister` — no rate limit.** Registration is rate-limited, credentials is rate-limited, but unregister is not. Unregister is an agent endpoint (JWT-protected), so probably fine.

Now let me look for genuinely wrong/unguarded behaviors or untested rules. Let me reconsider the `ca.Verify` `x509.Expired` handling, the `SignCSR` clamping, etc.

Let me look at one specific thing: **`Credentials` does not clamp when `now` is before CA's NotBefore in `Credentials` path but `Info` reads it fine.** Actually those are internal.

Let me reconsider potential real issues:

**Issue A: `Info` inhibited asymmetry** — already covered by test `TestInfo` line 513-516. Covered.

**Issue B: `Register` — `CountDevices` used for the limit includes devices regardless of status.** Hmm.

Let me look for something concrete and valuable. Let me look at the `http_test.go` to see RPC route coverage, and consider:

**Potential real issue — `handleCredentials` is missing a `maxBodyBytes` on `DecodeData`.** No wait, it has it.

**Rate-limit bypass:** `handleRegister` and `handleCredentials` check rate limits. But `handleInfo` and `handleVerify` are NOT rate-limited. These are device-auth endpoints (bearer secret required) that do bcrypt compare. An attacker with a valid device secret could hammer them — but they have a secret, so it's not an unauthenticated brute-force vector. However, `handleVerify` and `handleInfo` do a bcrypt `CompareHashAndPassword` on every request (via `authenticateDevice`), which is CPU-expensive. But they require a valid bearer secret, so an attacker needs valid credentials. Not a clear unguarded concern.

Actually wait — let me reconsider. `handleVerify` and `handleInfo` require `bearerSecret` but the *wrong*-secret path still burns a bcrypt. But they require a bearer secret header at all. Any client can send a bearer secret header with garbage; the service then does a full bcrypt compare (via `authenticateDevice`) — burning ~50-100ms CPU per request with NO rate limiting. This is a CPU DoS vector: an unauthenticated attacker can send unlimited `GET /devices/{id}` (Info) or `POST .../verify` with a fake bearer secret, forcing a bcrypt comparison each time. Registration and credentials are rate-limited per-IP, but Info and Verify are not.

That's a legitimate finding: **the two device-auth flows do full bcrypt work (authenticateDevice) with no per-IP rate limit, unlike registration and the credentials endpoint.** An unauthenticated attacker can burn server CPU (bcrypt at default cost) at full network speed through `GET /devices/{deviceID}` and `.../credentials/verify` with any fake bearer secret.

Let me verify by re-reading: `handleInfo` (line 197-214) calls `bearerSecret` then `svc.Info` → `authenticateDevice` → bcrypt compare. No limiter call. Yes. And `handleVerify` similarly.

This is a genuine finding — the two flows that perform expensive bcrypt work are the only two device endpoints NOT rate-limited, and they're reachable without valid credentials (any bearer string triggers the compare). This is worth proposing.

Let me double check whether upstream or DESIGN docs say Info/Verify should be rate-limited. Let me check docs/DESIGN.md for rate-limit §4.5.

Let me also examine the `ratelimit.go` and consider the interaction with `handleCredentials` double-token consumption.

Let me look at the http_test.go to see coverage and check the design doc for the rate-limiting intent.
$ cd /root/astrate-mule && grep -n "func Test" internal/pairing/http_test.go; echo "---ratelimit tests---"; grep -n "func Test" internal/pairing/ratelimit_test.go
203:func TestPairingHTTP(t *testing.T) {
557:func TestPairingHTTPRateLimit(t *testing.T) {
---ratelimit tests---
9:func TestLimiterBurstThenDeny(t *testing.T) {
33:func TestLimiterKeysAreIndependent(t *testing.T) {
52:func TestLimiterCapsTokensAtBurst(t *testing.T) {
73:func TestLimiterEvictionBoundsMemory(t *testing.T) {

Let me read the DESIGN doc section §4.4/§4.5 about rate limiting to understand what's the intended protection, and check the http_test RateLimit test to see which endpoints are covered.
→ Read internal/pairing/http_test.go [offset=557, limit=60]
→ Read docs/DESIGN.md [offset=690, limit=75]

$ cd /root/astrate-mule && grep -n "func Test" internal/pairing/ca/ca_test.go
56:func TestGenerateCACertificate(t *testing.T) {
84:func TestLoadRoundTrip(t *testing.T) {
110:func TestSignCSRFieldAssertions(t *testing.T) {
177:func TestSignCSRRSAKey(t *testing.T) {
196:func TestSignCSRTTLClampedToCA(t *testing.T) {
212:func TestSignCSRRejections(t *testing.T) {
247:func TestVerify(t *testing.T) {
291:func TestSerialUniqueness10k(t *testing.T) {

This confirms my finding. DESIGN.md §4.5 says "Rate limits: pairing endpoints (per-IP and per-device)". The credentials endpoint is rate-limited (tested at http_test.go:557), but **Info and Verify are not**. These two do a full bcrypt `CompareHashAndPassword` on every request (via `authenticateDevice` in service.go) and are reachable with any fake bearer secret — an unauthenticated attacker can burn server CPU at full network speed through `GET /devices/{deviceID}` and `POST .../verify`. The second credentials endpoint is protected but Info+Verify, which do the same bcrypt work, are not.

This is a genuine finding: the design doc says "pairing endpoints" should be rate-limited, but the two device-auth endpoints doing the expensive work are skipped.

Let me also verify the `handleInfo` and `handleVerify` don't go through any `requireAgent`/limiter middleware — confirmed from http.go lines 82-87: only register/unregister use `requireAgent`, and only handleRegister/handleCredentials apply limiters.

Let me also check the `ratelimit.go` for a possible second finding. Looking at `Limiter.Allow`, there's the bucket eviction issue but that's fine. Let me reconsider the `handleCredentials` double `Allow`: consumes 2 tokens from IP bucket (one for "ip|" and one for... no wait, `a.credLimiter.Allow("ip|"+ip.String())` and `a.credLimiter.Allow("dev|"+realm+"/"+deviceID)`. The first consumes from IP bucket, second from dev bucket. Two different keys, so one token from each. That's the documented per-IP AND per-device. Fine.

Now let me also verify the `Register` limit concern about `CountDevices` including confirmed devices — that's documented. Not proposing.

Let me now also check whether everything with bcrypt also goes through the same dummy-hash equalization. `authenticateDevice` returns ErrUnauthorized for bad secret after the real compare, and burns full compare. Correct.

Let me also double check one more potential finding: In `handleCredentials`, the rate-limited check happens BEFORE auth, good. In handleInfo/handleVerify, no rate limit at all.

Let me now also look at whether there's an untested rule worth proposing as a missing test. Let me check the `TestPairingHTTP` coverage versus what writeServiceError covers. Actually, let me look at whether `handleUnregister` (agent, DELETE) has rate limiting. It has `requireAgent` (JWT-protected) but no limiter. That's fine — agent endpoints are operator.

Let me also verify: is there coverage that `Verify` returns `Until` only when valid? That's tested.

Let me reconsider the "Info allows inhibited but Credentials/Verify reject" — this is actually an interesting asymmetry worth confirming it's intended. The DESIGN doc and comments say inhibited may read info (only new credentials/connections blocked). It's documented and tested. Not proposing.

Now, one more strong candidate. Let me think about the `remoteIP` no-X-Forwarded-For issue more. DESIGN says rate limit per-IP. If Astrate sits behind a reverse proxy (which DESIGN.md explicitly supports: "HTTP listener TLS-terminated in-binary or behind a reverse proxy"), then `remoteIP` returns the proxy's IP for ALL requests, so per-IP rate limiting becomes per-proxy (ineffective against fleets) — actually it would rate limit ALL devices behind the proxy as if one IP. Since DESIGN explicitly documents a reverse-proxy deployment, and the rate limit keys on `r.RemoteAddr` which would be the proxy, the per-IP limit degrades to a single global bucket in that deployment. This is worth noting but it's a deployment concern; the design doesn't mention proxy IP forwarding. It's a legitimate follow-up question but maybe better as a for-giulio note than a task. Actually upstream uses the real peer. Let me not over-reach.

I have my primary finding (Info/Verify not rate-limited despite doing bcrypt). That's a strong, specific, actionable finding with a clear test.

Let me also identify a second good proposal — a missing test. Let me think: the rate-limit test only covers the *credentials* endpoint's 429. There's no test asserting Info/Verify are NOT rate limited nor whether they should be. Since I'm proposing to ADD rate limiting to Info/Verify, the test would be part of the fix task.

Let me look for another independent missing-test finding to round out. Let me reconsider:

**`Register`'s limiter is per-IP only** (`handleRegister` line 106), while credentials is per-IP AND per-device. Register is agent (JWT), so per-IP on agent is a fleet-provisioning protection. Fine, documented.

Let me settle on findings:

1. **Info + Verify device endpoints do full bcrypt with no rate limit** (http.go handleInfo:197, handleVerify:238 → authenticateDevice → bcrypt, DESIGN §4.5 says pairing endpoints should be rate-limited). Propose adding per-IP (and per-device) limiters, plus a test. HIGH VALUE.

2. Let me reconsider whether there's a behavior that's *wrong* (not just missing). Let me look again at the `Allow` double-consume subtlety... Actually there's a subtle issue I should verify: `handleCredentials` does `a.credLimiter.Allow("ip|"+ip.String()) || a.credLimiter.Allow("dev|"+...)`. Because of short-circuit `||`, if the first (IP) check returns false, the second (dev) is NOT evaluated. That's fine — deny.

Let me look for a second, weaker but valid task: a missing test on an existing rule. Let me check `ca.go` SignCSR clamping — is there a test that TTL clamps to CA NotAfter? And `Load` mismatched key test? Let me check ca_test.go coverage.

Actually, let me look at one more concrete thing: the `formatDateTime` and verify envelope are golden-tested? Let me not spend too much. 

Let me verify the ca_test coverage for the clamping and expired-CA cases, which would determine if there's a missing-test proposal there.
$ cd /root/astrate-mule && grep -n "EnforceLatestCert\|EnforceCert\|Info\|Verify\|limiter\|Limiter\|regLimiter\|credLimiter\|registerRate\|register" internal/pairing/http.go
47:	regLimiter   *Limiter
48:	credLimiter  *Limiter
69:		regLimiter:   NewLimiter(cfg.RegisterRate, cfg.RegisterBurst),
70:		credLimiter:  NewLimiter(cfg.CredentialsRate, cfg.CredentialsBurst),
74:// Mount registers the pairing routes on mux. Paths are wire-frozen
81:		a.requireAgent(http.HandlerFunc(a.handleUnregister)))
85:		a.handleInfo)
87:		a.handleVerify)
92:// registerRequest is the flow A body: hw_id plus the Astrate
95:type registerRequest struct {
100:// registerResponse is the show-once secret envelope payload.
101:type registerResponse struct {
106:	if !a.regLimiter.Allow("ip|" + remoteIP(r).String()) {
111:	var req registerRequest
127:	_ = astarteapi.WriteData(w, http.StatusCreated, registerResponse{CredentialsSecret: secret})
130:func (a *API) handleUnregister(w http.ResponseWriter, r *http.Request) {
131:	err := a.svc.Unregister(r.Context(), r.PathValue("realm"), r.PathValue("deviceID"))
154:	if !a.credLimiter.Allow("ip|"+ip.String()) || !a.credLimiter.Allow("dev|"+realm+"/"+deviceID) {
197:func (a *API) handleInfo(w http.ResponseWriter, r *http.Request) {
203:	info, err := a.svc.Info(r.Context(), r.PathValue("realm"), r.PathValue("deviceID"), secret)
238:func (a *API) handleVerify(w http.ResponseWriter, r *http.Request) {
255:	res, err := a.svc.VerifyCredentials(r.Context(), r.PathValue("realm"), r.PathValue("deviceID"), secret, req.ClientCrt)
296:			map[string][]string{"error_name": {"already_registered"}})

My primary finding is confirmed. The ca_test already covers TTL clamping (`TestSignCSRTTLClampedToCA`) and Verify. So the ca package is well-tested.

Now let me settle on the review. My main finding is solid: **Info + Verify do full bcrypt work with no rate limiting** while the twin credentials endpoint that does the same bcrypt is rate-limited. This is a real, specific, actionable gap with a clear test.

Let me consider a second proposal. The `handleUnregister` agent endpoint is not rate-limited but it's JWT-protected — not a concern.

Let me reconsider the rate-limit protection conceptually. The value: `handleInfo`/`handleVerify` each trigger a bcrypt compare (~50-100ms CPU) on ANY request with a bearer header, regardless of whether the secret is valid. An unauthenticated attacker (fake secret) can pound these at full speed with no per-IP throttle. Meanwhile, the credentials endpoint — which does the *same* bcrypt work — is rate-limited at burst 5/sec. This is clearly inconsistent and contradicts DESIGN §4.5 ("Rate limits: pairing endpoints").

This is worth one strong proposal. For a second, let me pick a smaller, mechanical missing-test/clarity item. Let me reconsider the `Info` status logic — the `default: status = "pending"` catches `DeviceStatusRegistered` (and anything unknown). That's fine and tested (pending before credentials).

Actually, let me look at one more thing: the `Register` rate limit being per-IP applies to the agent endpoint. But agent register is JWT-operator-driven. DESIGN says register rate is about "fleet provisioning bursts". Fine.

Let me settle. I'll propose **one** high-value task (the Info/Verify rate limit) plus maybe a second lower-value one. Given the recipe says "at most five" and values quality over quantity, one strong task plus documenting what I decided not to propose is ideal. Let me also consider the `remoteIP` proxy finding as a for-giulio note (design consideration, not a tick-sized code task).

Let me write the review file. First let me get today's date: 2026-09-03.

Let me create the review file and append the task. Since this is a proposal job, I write the review doc and append task lines to todo.md.

Let me write the review, then append the task line.
$ cd /root/astrate-mule && head -1 .mule/reviews/engine-2026-09-03.md | sed 's/.*— //'; date +%F
2026-09-03
2026-09-03

Now let me write the review file for `pairing`.
← Write .mule/reviews/pairing-2026-09-03.md
Wrote file successfully.

Now let me append the task line to `.mule/todo.md`. Let me add it after the last line (line 127).
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -124,4 +124,5 @@
 - [!] consumer-properties-plaintext-fallback: in `internal/engine/control.go` `sendConsumerProperties` (line 204), when `devices.peek` returns nil (device not cached after eviction), the code defaults to zlib compression — but a device that declared `purge_properties_compression_format: plaintext` via the capabilities topic cannot decompress zlib. Fix by falling back to a store read of the device's capability state (or persisting the capability), and add a test in `internal/engine/capabilities_test.go` that evicts the cache entry for a plaintext device and asserts the purge payload is uncompressed. [auto] — BLOCKED: gates failed
 - [!] object-path-ok-comment: in `internal/engine/data.go` `objectPathOK` (line 330), the loop over `ci.ObjectLeaves` returns after checking the first entry. Add a comment explaining why one leaf suffices (upstream validation guarantees uniform prefix depth), so future readers do not mistake the early return for a bug. [auto] — BLOCKED: gates failed
 - [!] housekeeping-unit-tests: in `internal/housekeeping`, add a unit test file (no integration build tag) that exercises `CreateRealm` validation paths — blank name returns ErrValidation, blank JWT key returns ErrValidation, negative registration limit returns ErrValidation — using a mock store and sealer. The package currently has only integration tests gated behind `//go:build integration`, so these paths are never tested on the Pi. [auto] — BLOCKED: gates failed
 - [!] store-downsample-lttb: in `internal/store`, implement the TODO at store.go:135 — when `s.hasToolkit` is true, switch `Downsample` from `time_bucket+avg` to toolkit `lttb()` downsampling. The probe already records availability; the time_bucket path in datastreams.go is the always-working default. Add a test that exercises both code paths (mocking `hasToolkit`). [auto] — BLOCKED: gates failed
+- [ ] pairing-info-verify-ratelimit: in `internal/pairing/http.go`, add per-IP (and per-device) token-bucket rate limits to `handleInfo` (line 197) and `handleVerify` (line 238), which currently run a full bcrypt `CompareHashAndPassword` (via `authenticateDevice`, service.go:394) on any bearer string with no throttle — the sibling `handleCredentials` that does identical bcrypt work is limited at burst-5, and DESIGN §4.5 says "pairing endpoints" should be rate-limited. Reuse the credentials limiter defaults (1 req/s, burst 5) and extend `TestPairingHTTPRateLimit` to assert a burst of fake-secret requests to `GET /devices/{id}` and `POST .../verify` returns 429. [auto]

I've completed the code-review proposal job. Let me confirm I did not touch git and verify my outputs.

The recipe produced:
- `.mule/reviews/pairing-2026-09-03.md` (evidence file)
- 1 task line appended to `.mule/todo.md`

No code changes, no git operations. Let me verify.
$ cd /root/astrate-mule && git status --short
M .mule/cron.log
 M .mule/todo.md
?? .mule/reviews/pairing-2026-09-03.md

The working tree shows only my two intended file changes (plus a pre-existing `.mule/cron.log` modification from the cron job — not mine). No git operations performed.

```
Done: Reviewed internal/pairing — proposed per-IP/per-device rate limiting on Info+Verify (bcrypt DoS) + recorded a proxy-header note
Files: .mule/reviews/pairing-2026-09-03.md, .mule/todo.md
Verified: no gates run (proposal-only job, no code change) -> n/a
Unsure: nothing
Follow-ups: pairing-info-verify-ratelimit appended to .mule/todo.md (1 line)
```

**Summary of the review** (first review of `pairing`; follows the store→realm→flow→engine rotation):

**Proposed (1 task):** `handleInfo` and `handleVerify` (http.go:197, 238) do a full bcrypt compare via `authenticateDevice` on any bearer string, but — unlike the sibling `handleCredentials` (same bcrypt cost, capped at burst-5) — they have **no rate limit**. An unauthenticated attacker can burn server CPU at full speed, and DESIGN §4.5 explicitly says "pairing endpoints" should be rate-limited.

**Noted for Giulio (not queued):** `remoteIP` ignores proxy headers, so behind a reverse proxy (which DESIGN supports) the per-IP throttles collapse to one shared bucket — a config/design call.

**Deliberately not proposed:** the Info-vs-Credentials inhibited asymmetry (intentional, tested), the `Register` limit counting confirmed devices (upstream parity, tested), the double `Allow` in credentials (correct per-IP+per-device split), and the `ca`/`ratelimit` packages (well-tested, no gaps).

Per your instructions I did not implement the proposed task and did not touch git.
