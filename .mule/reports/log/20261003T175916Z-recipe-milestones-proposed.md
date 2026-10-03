slug: recipe-milestones
verdict: proposed
at:  1ac1967
ran: 2026-10-03T17:59:16Z on DietPi in 363s


> build · big-pickle

I'll start by reading the required context files.
→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md

$ ls -la /root/astrate-mule/.mule/ && echo "--- tasks ---" && ls -la /root/astrate-mule/.mule/tasks/ 2>/dev/null | head -50
total 10776
drwxr-xr-x  7 root root     4096 Oct  3 19:05 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       13 Oct  3 19:53 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  3 19:05 .heartbeat
-rw-r--r--  1 root root        0 Oct  3 19:53 .last-output
-rw-r--r--  1 root root        0 Oct  3 19:53 .lock
-rw-r--r--  1 root root       11 Oct  3 19:53 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 10701458 Oct  3 19:53 cron.log
drwxr-xr-x  2 root root     4096 Oct  3 13:30 failed
-rw-r--r--  1 root root    97814 Oct  2 13:44 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    13046 Oct  3 19:05 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  2 13:43 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   128124 Oct  3 19:05 todo.md
--- tasks ---
total 60
drwxr-xr-x 2 root root 4096 Sep  7 21:11 .
drwxr-xr-x 7 root root 4096 Oct  3 19:05 ..
-rw-r--r-- 1 root root 3293 Sep  7 21:11 docs-sync-appengine-data-output-params.md
-rw-r--r-- 1 root root 5539 Jul 27 22:18 issue-16.md
-rw-r--r-- 1 root root 7510 Sep  4 20:46 issue-23.md
-rw-r--r-- 1 root root 6434 Sep  4 20:46 issue-24.md
-rw-r--r-- 1 root root 5665 Sep  4 20:46 issue-25.md
-rw-r--r-- 1 root root 6009 Sep  4 20:46 issue-26.md
-rw-r--r-- 1 root root 5683 Sep  4 20:46 issue-27.md

$ cd /root/astrate-mule && wc -l .mule/research/upstream-parity-2026-08-22.md docs/UPSTREAM-EXPERIMENTAL.md 2>&1; echo "--- labels ---"; gh issue list --label milestone-3.0 --state all --limit 50 2>&1 | head -60
wc: .mule/research/upstream-parity-2026-08-22.md: No such file or directory
  25 docs/UPSTREAM-EXPERIMENTAL.md
  25 total
--- labels ---

$ cd /root/astrate-mule && gh issue list --label upstream-parity --state all --limit 60 --json number,title,state,labels --template '{{range .}}{{.number}}\t{{.state}}\t{{range .labels}}{{.name}},{{end}}\t{{.title}}{{"\n"}}{{end}}' 2>&1
93\tOPEN\tmule-review,upstream-parity,\tIn internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
92\tOPEN\tenhancement,upstream-parity,upstream-experimental,\tkeyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds
91\tCLOSED\tmule-blocked,upstream-parity,\tPairing health endpoint divergence: upstream serves GET /pairing/health (no realm), Astrate serves GET /pairing/v1/{realm}/health
89\tCLOSED\tenhancement,upstream-parity,\tDashboard flow-block schema mismatch (split_map/virtual pools hardcoded; null_sink/log_sink unknown)
88\tCLOSED\tenhancement,upstream-parity,\tFlow auth: support a_f JWT claim
87\tCLOSED\tenhancement,upstream-parity,\tFlow block: lua_map — needs embedded Lua runtime (parked)
86\tCLOSED\tenhancement,upstream-parity,\tFlow: pipeline source DSL — keep DAG-JSON as documented deviation?
85\tCLOSED\tenhancement,upstream-parity,\tFlow API: user-defined composite blocks
84\tCLOSED\tenhancement,upstream-parity,\tFlow blocks: virtual_device_pool / dynamic_virtual_device_pool
83\tCLOSED\tenhancement,upstream-parity,\tFlow blocks: mqtt_source/mqtt_sink (+modbus_tcp_source?) demand-driven
82\tCLOSED\tenhancement,upstream-parity,\tFlow blocks: http_source/http_sink (demand-driven)
81\tCLOSED\tenhancement,upstream-parity,\tFlow block: json_path_map
80\tCLOSED\tenhancement,upstream-parity,\tFlow blocks: pure-transform set (to_json, update_metadata, split_map, random_source, sort)
79\tCLOSED\tenhancement,upstream-parity,\tVerify registration-limit-reached HTTP status vs upstream
78\tOPEN\tenhancement,milestone-4.0,upstream-parity,\tFDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
77\tCLOSED\tenhancement,upstream-parity,\tVerify per-service version endpoints (GET /version, GET /v1/{realm}/version) served everywhere
76\tCLOSED\tenhancement,upstream-parity,\tHousekeeping: GET /v1/realm-defaults/replication — decide reject/deviate (Cassandra-shaped)
75\tCLOSED\tenhancement,upstream-parity,\tHousekeeping: decide realm-deletion gating/preconditions vs always-sync deviation
74\tCLOSED\tenhancement,upstream-parity,\tHousekeeping: PATCH /v1/realms/{realm} (jwt key, registration limit, retention; null=unset)
73\tCLOSED\tenhancement,upstream-parity,\tHousekeeping: default datastream retention injection env var (upstream 1.4)
72\tCLOSED\tenhancement,upstream-parity,\tRealms: datastream_maximum_storage_retention ceiling (create/patch/enforce)
71\tCLOSED\tenhancement,upstream-parity,\tPairing: realm-scoped health check GET /v1/{realm}/health (upstream 1.3)
70\tCLOSED\tenhancement,upstream-parity,\tTriggers: audit wildcard semantics (interface_name '*', match_path '/*' forcing rules)
69\tCLOSED\tenhancement,upstream-parity,\tVerify unknown-realm HTTP status on RM endpoints against upstream
68\tCLOSED\tenhancement,mule-blocked,upstream-parity,upstream-experimental,\tDecide async_operation=false params vs documented always-sync deviation
67\tCLOSED\tenhancement,upstream-parity,upstream-experimental,\tInterfaces: decide handling of required and encrypted mapping fields (upstream 1.4)
66\tCLOSED\tenhancement,upstream-parity,\tRealm Management: detailed=true interface listing with full mappings (upstream 1.4)
65\tCLOSED\tenhancement,upstream-parity,\tPolicies: handler-overlap rejection + retry_times coupling + prefetch_count
64\tCLOSED\tenhancement,upstream-parity,\tTriggers: decide AMQP action behavior (validate-reject vs NATS-forward deviation)
63\tCLOSED\tenhancement,upstream-parity,\tTriggers: HTTP action validation limits (URL/method/header blocklist/template size)
62\tCLOSED\tenhancement,upstream-parity,\tRealm Management: audit install/update/delete error codes and statuses
61\tCLOSED\tenhancement,upstream-parity,\tRealm Management: audit interface/mapping validation matrix against astarte_core
60\tCLOSED\tenhancement,upstream-parity,\tRealm Management: GET config/datastream_maximum_storage_retention (since upstream 1.2.0)
59\tCLOSED\tenhancement,upstream-parity,\tAppEngine: group create-body validation + UUID-v1 from_token for group device listing
58\tCLOSED\tenhancement,upstream-parity,\tAppEngine: PATCH requires Content-Type application/merge-patch+json
57\tCLOSED\tenhancement,upstream-parity,\tAppEngine: audit server-write error taxonomy against upstream
56\tCLOSED\tenhancement,upstream-parity,\tAppEngine: object-aggregate response formats table and disjoint_tables
55\tCLOSED\tenhancement,upstream-parity,\tAppEngine: complete data-query parameter set
54\tCLOSED\tenhancement,upstream-parity,\tAppEngine: GET /groups/{g} show-group endpoint
53\tCLOSED\tenhancement,upstream-parity,\tAppEngine: group-scoped device PATCH + interface data access under groups/{g}/devices/{d}
52\tCLOSED\tenhancement,upstream-parity,\tAppEngine: devices-by-alias full mirror (PATCH via alias + interface data access)
51\tCLOSED\tenhancement,upstream-parity,\tkeyAgreement control topic + encrypted mappings (upstream 1.4 experimental) — investigate/park
50\tCLOSED\tenhancement,upstream-parity,\tBroker/engine hardening: discard malformed introspection/purge/capabilities instead of crashing
49\tCLOSED\tenhancement,upstream-parity,\tMQTT v1: binaryblob properties re-send as BSON subtype 0 in consumer/properties
48\tCLOSED\tenhancement,upstream-parity,\tMQTT v1: tolerate devices connecting with an empty introspection
47\tCLOSED\tenhancement,upstream-parity,\tMQTT v1: accept device capabilities incl. purge_properties_compression_format (plaintext purge)

$ cd /root/astrate-mule && cat docs/UPSTREAM-EXPERIMENTAL.md
# Upstream experimental features register

Astarte's brain is the upstream Astarte team. When upstream ships a feature
marked experimental / unstable / rc-only, Astrate adopts it early (so the wire
behaviour is ready) but tracks it here until upstream promotes or deprecates
it. **Review trigger: every time a future parity milestone targets a newer
upstream version, check each row below against what that release did with the
feature** — promoted → keep and drop the row; deprecated → deprecate in the
same milestone, in sync.

## The version-declaration rule

`APICompatVersion` (internal/realm/service.go) must never claim a level higher
than what Astrate actually implements, and must never be operator-configurable.
It is bumped only in the same change that completes the full surface of that
upstream level. When that happens, every open row below tagged with that level
must be reconciled first — the register is the checklist, so nothing has to be
remembered.

## Rows

| Feature | Issue | Upstream level | Adopted | Reconcile when |
|---|---|---|---|---|
| Required + encrypted interface mapping fields | #67 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
| `async_operation=false` on mutating RM endpoints | #68 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |

$ cd /root/astrate-mule && ls .mule/research/ 2>/dev/null | tail -30; echo "--- open all ---"; gh issue list --state open --limit 60 --json number,title,labels --template '{{range .}}#{{.number}}\t[{{range .labels}}{{.name}},{{end}}]\t{{.title}}{{"\n"}}{{end}}' 2>&1
--- open all ---
#116\t[mule-alarm,]\tmule: nothing has landed in 14h
#115\t[mule-alarm,]\tmule: nothing has landed in 14h
#114\t[mule-alarm,]\tmule: nothing has landed in 8h
#113\t[mule-alarm,]\tmule: nothing has landed in 14h
#112\t[mule-alarm,]\tmule: nothing has landed in 14h
#111\t[mule-alarm,]\tmule: nothing has landed in 16h
#110\t[mule-alarm,]\tmule: nothing has landed in 14h
#109\t[mule-alarm,]\tmule: nothing has landed in 14h
#108\t[mule-alarm,]\tmule: nothing has landed in 15h
#107\t[mule-alarm,]\tmule: nothing has landed in 14h
#106\t[mule-alarm,]\tmule: nothing has landed in 14h
#105\t[mule-alarm,]\tmule: nothing has landed in 14h
#104\t[mule-alarm,]\tmule: nothing has landed in 14h
#103\t[mule-alarm,]\tmule: nothing has landed in 16h
#102\t[mule-alarm,]\tmule: nothing has landed in 14h
#101\t[mule-alarm,]\tmule: nothing has landed in 14h
#100\t[mule-alarm,]\tmule: nothing has landed in 15h
#99\t[mule-alarm,]\tmule: nothing has landed in 14h
#98\t[mule-alarm,]\tmule: nothing has landed in 16h
#97\t[mule-alarm,]\tmule: nothing has landed in 15h
#96\t[mule-alarm,]\tmule: nothing has landed in 14h
#95\t[mule-alarm,]\tmule: nothing has landed in 15h
#94\t[mule-alarm,]\tmule: nothing has landed in 15h
#93\t[mule-review,upstream-parity,]\tIn internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
#92\t[enhancement,upstream-parity,upstream-experimental,]\tkeyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds
#78\t[enhancement,milestone-4.0,upstream-parity,]\tFDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
#1\t[wontfix,]\tProvide an Open Source IoT Platform unironically, unironically

$ cd /root/astrate-mule && tail -120 .mule/for-giulio.md
json_path_map, pure-transform set, virtual_device_pool, container block MVP, flow API,
  durable named flows all landed).   (Milestones recipe run, 2026-09-03.)

  **Decided 2026-09-04: the tag is `v0.2.0`, cut on the same day.** The `v2.0`/`v3.0` names
  are milestone names, not release versions — the project is still pre-1.0 and the version
  number keeps its own line rather than jumping two majors to match a milestone label.

---
- ~~`docs/site/appengine-api.md` documents `GET` and `DELETE` on `/appengine/v1/<realm>/groups/<name>`~~
  — **resolved 2026-09-04, and the original note was half wrong.** `GET /groups/{group}` does
  exist (`internal/appengine/http.go:69`); only `DELETE` was absent, and it is absent from
  upstream's own spec too (`docs/api/astarte_appengine_api.yaml` has no `/groups/{group}`
  path at all). So the docs line was spurious, not the code. Line dropped.
---
- ~~`Router.Submit` TOCTOU on the `closed` flag~~ — **fixed 2026-09-04**, and none of the three
  options in the original note was taken: (b) was wrong (a `select` does not save a send on a
  closed channel — it panics either way). The fix is that `Drain` no longer closes the lane
  channels at all; it closes `quit` under a write lock while `Submit` holds the read lock
  across its send, so a lane can never be retired underneath an in-flight sender, and the lanes
  exit on `quit` after draining what is buffered. `TestRouter_SubmitParkedWhenDrainRuns`
  reproduces the old panic deterministically and passes on the fix.
---
- ~~**#87 `lua_map` — needs embedded Lua runtime, parked.**~~ — **closed 2026-09-04.**
  Embedding a Lua VM is not machine-checkable by the mule and Lua flow support is on no active
  roadmap. Reopen in a minute if that changes.
- **#78 FDO device onboarding — milestone-4.0, investigation phase.** Too large for a single
  mule task; the investigation work (reading upstream's TO2 handling, inventorying endpoints,
  schema and keys) is a multi-session project. Parked until the v3.0 queue clears and this
  becomes the next milestone target.
- **#1 is never to be raised again.** Giulio's standing instruction, 2026-09-04: it stays open
  permanently and is not a candidate for closing, triage, or a for-giulio entry. Do not
  propose it again.
---

- ~~**Flow v2.0: named multi-instance flows + pipeline config?**~~ — **decided
  2026-07-29: (b) named multi-instance + config.** Design then implement (#40).
  Doc: `docs/handoff/flow-v2-decisions-2026-07-29.md`.
- ~~**Flow v2.0: durable `flows` table vs in-memory only?**~~ — **decided
  2026-07-29: (b) durable records; `auto_restart` default true, optional never;
  process rehydrates on boot; fail loudly on bad pipeline/start.** Filed #41;
  edge-case follow-ups #42. Same design package as #40.
- ~~**Flow v2.0: containers in scope?**~~ — **decided 2026-07-29: yes, phased
  PoC → MVP (#43).** Native Lua/MQTT blocks are *not* a v2.0 gate (containers cover
  custom logic). Doc: `docs/handoff/flow-v2-decisions-2026-07-29.md`.
- ~~`device_deletion_started`/`device_deletion_finished` trigger events are not emitted~~ —
  **decided 2026-07-27: emit both, back-to-back, around the synchronous delete.** Filed as
  issue #21 (`mule`). (Cross-project survey, 2026-07-27,
  `.mule/research/survey-2026-07-27.md` source 4.)
- ~~Mustache trigger-action templates are accepted but not rendered~~ — **decided
  2026-07-27: implement it.** Guiding principle clarified: Astarte compatibility means
  SDK/wire compatibility, not minimum dependency count — Astrate is allowed to be a
  compatible *superset*. Library picked: `github.com/cbroglie/mustache`. Filed as issue #22
  (`mule`). (Same survey, source 4.)
- ~~`value_change`/`value_change_applied`/`path_created`/`path_removed`/`value_stored` trigger
  types compile but never fire~~ — **resolved: implemented in commit 6bd14a7 (2026-08-22)**,
  and the cost question answered by #20's bench on the Legion Go (2026-08-24 closeout):
  ~0.13 ms keyed read per message, −0.7% throughput even paid once per message over REST.
  Performance is not the constraint for these types (nor, a fortiori, for the group-scoped
  line below). (Same survey, source 4.)
- **Group-scoped triggers (`group_name` on device/data triggers) compile but never match**
  (`internal/engine/triggers/match.go:11-12`). Decision deferred, tied to issue #17
  (group-WATCH-path reconciliation, trickle work, not mule): whatever group-membership
  mechanism comes out of that phase should also report the perf cost for this decision —
  noted in a comment on #17 so it isn't benchmarked twice. (Same survey, source 4.)

---

- ~~The Pi cannot run the race detector~~ — **resolved 2026-07-27** by installing Go 1.26.5
  as a userland toolchain on the Legion Go (`~/.local/go`, no root, `rm -rf` to undo). The
  Pi still cannot run `-race` (39-bit VMA kernel vs the 48 ThreadSanitizer needs), so its
  gate remains `go vet ./... && go test ./...` — but race coverage now exists on the Legion
  Go, where the full suite runs clean in ~40s on 16 cores. The standing `race-check` task is
  the concurrency gate. Concurrency work is queueable again, provided the race-check runs
  after it.
- ~~golangci-lint is not installed on the Pi~~ — **resolved 2026-07-28**: installed v2.12.2
  via `go install` (the prebuilt-binary installer's published sha256 for linux/arm64 did not
  verify, reproducibly, so built from source instead). `mule.service` now sets
  `MULE_LINT_CMD=golangci-lint run ./...` and has `/root/go/bin` on `PATH`; the lint gate
  runs starting with the next tick.
- ~~`/root/astrate` on the Pi has uncommitted work~~ — **resolved 2026-07-27** with the new
  `tools/reconcile.sh`: rescued onto `origin/wip/DietPi-20260727T171543Z` (pushed, not
  reviewed — read the diff before merging anything from it) and `/root/astrate` is now clean
  on `main`. Also had to set `commit.gpgsign false` locally on that clone first — it had
  signing on with no working gpg-agent for a non-interactive session, same class of problem
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

## 2026-09-25 — the `-tags nats` half of the tree is compiled by no gate

`cmd/astrate/newnats_nats.go`, `internal/engine/forward/nats.go` and
`internal/engine/forward/nats_test.go` are all behind `//go:build nats`, and nothing
builds or tests that tag: not `make build`/`make test`/`make test-integration`/
`make test-e2e` (Makefile:23, 38, 42, 46), not the mule gate, not the lint target.
So the whole NATS trigger-forwarding path is invisible to CI — while
`config.validate` accepts `triggers.forward.kind = "nats"` (internal/config/config.go:339-343)
and the untagged binary then fails at boot with the build-tag message
(cmd/astrate/newnats_default.go). I measured it today: `go build -tags nats ./cmd/...`
and `go vet -tags nats ./cmd/...` both exit 0, so there is nothing broken — just nothing
keeping it unbroken. Worth a `make build-nats` (and the NATS test tier, which needs a
container, so it would be a `[legion]` line) in the gate; not a queue line because the
only test I can imagine for it shells out to `go build`, and the mule's own check strips
the implementation and the guard together, so it would never fail.
Not queued for the same reason: `IdleTimeout` is missing on the HTTP server
(cmd/astrate/main.go:221) and `-healthcheck` ignores the config file
(cmd/astrate/main.go:531) — both real, both yours to call, neither has a test worth
writing. See `.mule/reviews/cmd-astrate-2026-09-25.md`.
- **The mule has been idle 16h.** Filed by the dead-man's switch; see journalctl on the Pi.

$ cd /root/astrate-mule && ls -la .mule/research 2>&1 | head; echo "=== git branch/head (read-only) ==="; git log --oneline -5 2>&1; echo "=== #90 ==="; gh issue view 90 --json number,title,state,labels,body --template '{{.number}} {{.state}} [{{range .labels}}{{.name}},{{end}}]{{.title}}{{"\n"}}{{.body}}' 2>&1 | head -40
ls: cannot access '.mule/research': No such file or directory
=== git branch/head (read-only) ===
1ac1967 mule: recipe hygiene ran (172s)
5fb0e21 mule: log docs-sync-hk-error-detail-examples-split
23006d4 mule: docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
87a86e6 mule: log docs-sync-native-error-example-capitalisation
7cc237f mule: docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `"Internal server error"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
=== #90 ===
90 CLOSED [enhancement,milestone-4.0,]milestone-4.0: decide on bumping APICompatVersion from "1.2.2" to "1.3.0-rc.0" (or making it configurable)
Reminder of a pending decision parked during the upstream-parity investigation (#47–#89). **Do not raise the version now** — this is a milestone-4.0 item.

## Context

`realm.APICompatVersion` (`internal/realm/service.go`) is currently hardcoded to `"1.2.2"` — the emulated upstream API level Astrate declares.

Edgehog 0.13.1 gates just-in-time provisioning triggers (`device_registered`, `device_deletion_started`) behind a Realm Management version >= 1.3. Astrate's engine already supports those trigger kinds end-to-end — verified in #28 (Edgehog compatibility investigation) — but Edgehog refuses to use them because of the declared version string.

## What this issue must cover when picked up

1. **Audit first**: what else does declaring "1.3" imply for other consumers?
   - `astarte-dashboard` is currently validated against 1.2.2 behavior (see #13, error-name bundle comparison). Check whether any dashboard code path feature-gates on RM version besides the policies UI gate (>= 1.1.1).
   - Sweep the parity findings (#47–#89): several are "1.3/1.4 features" we deliberately do NOT implement yet — does declaring 1.3 overpromise any of them?
   - Check device SDKs and `astartectl` for version-gated behavior.
2. **Then decide one of**:
   - bump the constant to `"1.3.0-rc.0"` (upstream's own pre-release numbering for that line);
   - or make it configuration (`internal/config`) so operators can declare the level their consumers expect.
3. Update `docs/COMPATIBILITY.md` deviation #10 wording accordingly (propose, don't edit directly).

## Outcome if done today (why it's parked)

Bumping without the audit risks silently promising surfaces we haven't verified; not bumping keeps Edgehog's JIT provisioning disabled even though the engine supports it. Milestone-4.0 scope decision.

$ cd /root/astrate-mule && echo "=== for-giulio mentions of milestone v3.0/cut the tag ==="; grep -n "cut the tag\|milestone v3\|milestone-3.0\|v3.0 looks\|v3.0 is" .mule/for-giulio.md | tail -30; echo "=== todo.md milestone lines ==="; grep -n "^- \[ \] milestone" .mule/todo.md | tail -20; echo "=== todo total lines ==="; wc -l .mule/todo.md; echo "=== open task lines ==="; grep -c "^- \[ \]" .mule/todo.md
=== for-giulio mentions of milestone v3.0/cut the tag ===
61:- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Status: not started", which is no longer true.** Release sweep (`gh api repos/astarte-platform/astarte/releases`, run this session) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.5` (2026-08-20, prerelease) — byte-identical to what the 2026-09-24 run saw, so nothing shipped upstream in the two days since. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Open non-alarm issues remain exactly **#93** (aclhook comment, `mule-review`, its own recipe path), **#92** (keyAgreement, parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** (FDO, `milestone-4.0`), **#1** (untouched per standing instruction). No new gaps, **no issues filed, no `.mule/todo.md` task lines added** — there is nothing machine-checkable left, and filing anything here would duplicate #92/#93/#78.
63:  **The recipe's standing item is satisfied, and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still has exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental* with reconcile trigger "upstream 1.4 final". v1.4 is not final, so both rows stay as written: nothing to promote, nothing to deprecate, and the register is accurate. `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision; the milestone-bump issue the recipe asks for **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": the milestone's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is done, and the 1.4 half is unresolved (#92 open, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut. Also still true and already escalated on 2026-09-06, so not repeated: the section's source document `.mule/research/upstream-parity-2026-08-22.md` does not exist in the repo.
136:- **Milestone recipe run, 2026-09-24: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-23 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
144:- **Milestone recipe run, 2026-09-23: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-21 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
160:- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
168:- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
180:- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
184:- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
200:- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) shows a new **stable v1.3.4** (2026-09-18) — maintenance only (data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change), so no new gap for Astrate — and still **no stable v1.4.0**: newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
212:- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-13 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
216:- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-12 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
220:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
228:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
230:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
261:  upstream gates.** `milestone-3.0` label remains empty. The three open `upstream-parity`
283:  `milestone-3.0` label is empty (this milestone's work lives under `upstream-parity`).** All
296:  proposes **no** "complete, cut the tag" line: the open decision set is real and is yours.
382:- **milestone 2.0 looks complete, verify and cut the tag** — all 11 `milestone-2.0` issues
480:When v3.0 is marked DONE, the v4.0 section of `.mule/milestones.md` should be
=== todo.md milestone lines ===
=== todo total lines ===
263 .mule/todo.md
=== open task lines ===
14

$ cd /root/astrate-mule && echo "=== #90 comments ==="; gh issue view 90 --comments 2>&1 | tail -40; echo; echo "=== #92 ==="; gh issue view 92 --comments 2>&1 | head -80
=== #90 comments ===
author:	atsetilam
association:	member
edited:	false
status:	none
--
Decision frozen 2026-08-23: **do not bump now, and never make it configurable.**

- `APICompatVersion` stays `1.2.2` until Astrate implements the FULL 1.3 surface; it is bumped only in the same change that completes that surface, after reconciling every row in docs/UPSTREAM-EXPERIMENTAL.md tagged with the level.
- Configurability rejected: an operator raising the declared level would make clients believe features exist that don't — the costliest lie available.
- Enforcement is mechanical, not memory-based: the bump rule lives inline on the constant (internal/realm/service.go), as the standing checklist item in .mule/recipes/milestones.md, and in docs/UPSTREAM-EXPERIMENTAL.md's version-declaration rule.
--

=== #92 ===
author:	atsetilam
association:	member
edited:	false
status:	none
--
Parked by Giulio on 2026-09-04: **wait for a stable \`v1.4.0\`.** The spec is complete but ships only in \`v1.4.0-rc.5\`, and building against an rc risks rebuilding it.

So this issue is not waiting on anyone to remember it. It has a row in [`.mule/waiting-on.md`](https://github.com/astrate-platform/astrate/blob/main/.mule/waiting-on.md), and the weekly `mule-upstream-watch` job now reads that file as its first step every run: when a stable `v1.4.0` tag exists upstream, it escalates to `.mule/for-giulio.md` naming this issue. An `-rc.N` does not satisfy the row.

That mechanism was added because it was missing here: #51 was parked on the same condition, the condition fired on 2026-08-31, and nothing connected the new upstream tags to the parked issue — it took a hand triage four days later to catch it.
--

$ cd /root/astrate-mule && cat .mule/waiting-on.md 2>&1 | head -60; echo "=== grep waiting-on in for-giulio tail ==="; grep -n "waiting-on\|upstream-watch" .mule/for-giulio.md | tail -10
cat: .mule/waiting-on.md: No such file or directory
=== grep waiting-on in for-giulio tail ===
196:- **github-issues triage run, 2026-09-18: still nothing proposable — the mule-alarm pile is now 13 straight days (#94–#106).** Sixteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-17 run: **#106**, today's alarm (created 09-18 10:54Z, "nothing has landed in 14h"); **#105** has expired into the pile (superseded by #106). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-17 one: close #94–#105** — each was a one-day event, superseded by next day's alarm, never actionable; leave #106 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
208:- **github-issues triage run, 2026-09-17: still nothing proposable — the mule-alarm pile is now 12 straight days (#94–#105).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-16 run: **#105**, today's alarm (created 09-17 10:51Z, "nothing has landed in 14h"); **#104** has expired into the pile (superseded by #105). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-16 one: close #94–#104** — each was a one-day event, superseded by next day's alarm, never actionable; leave #105 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
210:- **github-issues triage run, 2026-09-16: still nothing proposable — the mule-alarm pile is now 11 straight days (#94–#104).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), plus the daily alarm. New since the 2026-09-15 run: **#104**, today's alarm (created 09-16 11:12Z, "nothing has landed in 14h"); **#103** has expired into the pile (superseded by #104). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-15 one: close #94–#103** — each was a one-day event, superseded by next day's alarm, never actionable; leave #104 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
214:- **github-issues triage run, 2026-09-15: still nothing proposable — the mule-alarm pile is now 10 straight days (#94–#103).** Fourteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-13 run: **#103**, the alarm created 2026-09-14 11:27Z ("nothing has landed in 16h"); **#102** has expired into the pile (superseded by #103). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-13 one: close #94–#102** — each was a one-day event, superseded by next day's alarm, never actionable; leave #103 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
218:- **github-issues triage run, 2026-09-13: still nothing proposable — the mule-alarm pile is now 9 straight days (#94–#102).** Thirteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-12 run: **#102**, today's alarm (created 11:09Z, "nothing has landed in 14h"); **#101** has expired into the pile (superseded by #102). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-12 one: close #94–#101** — each was a one-day event, superseded by next day's alarm, never actionable; leave #102 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
224:- **github-issues triage run, 2026-09-12: still nothing proposable — the mule-alarm pile is now 8 straight days (#94–#101).** Twelve open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-11 run: **#101**, today's alarm (created 10:58Z, "nothing has landed in 14h"); **#100** has expired into the pile (superseded by #101). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-11 one: close #94–#100** — each was a one-day event, superseded by next day's alarm, never actionable; leave #101 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
226:- **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pile is now 7 straight days and worth a look.** Twelve open issues. Still no machine-checkable fix candidates, so no task lines: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-10 run: **#100**, today's alarm (created 10:55Z, "nothing has landed in 15h"). #99–#94 are the same `mule-alarm` (one per day 2026-09-05→09-10, ~11:00Z each, zero comments) — telemetry, not code issues, so never proposable. But seven consecutive alarms is past the "low-activity window" wording earlier runs used: the mule has landed no commit since ~2026-09-04/05, and `main`'s queue copy shows the top line as `- [ ]` while the landed work below is all `[!] BLOCKED` (`wrote nothing` / `tests failed`) — though `mule/queue` is the authoritative copy. **Proposal: close #94–#99** — each was a one-day event, superseded by next day's alarm, never actionable; leave #100 (live today) to self-expire. If the idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
257:  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated. #1
274:  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated below.
314:  awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the waiting-on

$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/releases --jq '.[] | "\(.tag_name)\t\(.published_at)\tprerelease=\(.prerelease)\t\(.name)"' 2>&1 | head -15; echo "=== latest tags ==="; gh api repos/astarte-platform/astarte/tags --jq '.[].name' 2>&1 | head -15
v1.4.0-rc.6	2026-09-30T07:28:27Z	prerelease=true	v1.4.0-rc.6
v1.3.4	2026-09-18T09:42:57Z	prerelease=false	v1.3.4
v1.4.0-rc.5	2026-08-20T12:28:04Z	prerelease=true	v1.4.0-rc.5
v1.3.3	2026-08-07T14:33:58Z	prerelease=false	
v1.4.0-rc.4	2026-08-06T13:40:20Z	prerelease=true	v1.4.0-rc.4
v1.4.0-rc.3	2026-07-31T10:48:39Z	prerelease=true	v1.4.0-rc.3
v1.4.0-rc.2	2026-07-14T10:16:22Z	prerelease=true	v1.4.0-rc.2
v1.4.0-rc.1	2026-07-14T08:53:27Z	prerelease=true	v1.4.0-rc.1
v1.3.2	2026-07-14T10:07:26Z	prerelease=false	v1.3.2
v1.3.1	2026-07-13T11:50:05Z	prerelease=false	v1.3.1
v1.3.0	2026-05-06T15:30:50Z	prerelease=false	v1.3.0
v1.2.2	2026-04-27T16:39:28Z	prerelease=false	v1.2.2 - Bitter Kola Nut
v1.4.0-rc.0	2026-04-08T12:32:33Z	prerelease=true	v1.4.0-rc.0
v1.3.0-rc.2	2026-04-08T16:34:59Z	prerelease=true	v1.3.0-rc.2
v1.2.2-rc.0	2026-04-08T16:31:03Z	prerelease=true	v1.2.2-rc.0
=== latest tags ===
v1.4.0-rc.6
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
v1.4.0-rc.0
v1.3.4
v1.3.3
v1.3.2
v1.3.1
v1.3.0
v1.3.0-rc.2
v1.3.0-rc.1
v1.3.0-rc.0

$ cd /root/astrate-mule && cat .mule/.gitignore; echo "=== tracked waiting-on? ==="; git ls-files .mule | grep -i waiting; echo "(exit $?)"; echo "=== upstream-watch refs ==="; grep -rn "waiting-on" --include="*" .mule/ tools/ 2>/dev/null | grep -v for-giulio | grep -v cron.log | head -20; echo "=== recipes dir ==="; ls .mule/recipes/
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
=== tracked waiting-on? ===
(exit 1)
=== upstream-watch refs ===
.mule/.last-output:366:cat: .mule/waiting-on.md: No such file or directory
.mule/.last-output:368:196:- **github-issues triage run, 2026-09-18: still nothing proposable — the mule-alarm pile is now 13 straight days (#94–#106).** Sixteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-17 run: **#106**, today's alarm (created 09-18 10:54Z, "nothing has landed in 14h"); **#105** has expired into the pile (superseded by #106). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-17 one: close #94–#105** — each was a one-day event, superseded by next day's alarm, never actionable; leave #106 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/.last-output:369:208:- **github-issues triage run, 2026-09-17: still nothing proposable — the mule-alarm pile is now 12 straight days (#94–#105).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-16 run: **#105**, today's alarm (created 09-17 10:51Z, "nothing has landed in 14h"); **#104** has expired into the pile (superseded by #105). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-16 one: close #94–#104** — each was a one-day event, superseded by next day's alarm, never actionable; leave #105 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/.last-output:370:210:- **github-issues triage run, 2026-09-16: still nothing proposable — the mule-alarm pile is now 11 straight days (#94–#104).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), plus the daily alarm. New since the 2026-09-15 run: **#104**, today's alarm (created 09-16 11:12Z, "nothing has landed in 14h"); **#103** has expired into the pile (superseded by #104). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-15 one: close #94–#103** — each was a one-day event, superseded by next day's alarm, never actionable; leave #104 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/.last-output:371:214:- **github-issues triage run, 2026-09-15: still nothing proposable — the mule-alarm pile is now 10 straight days (#94–#103).** Fourteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-13 run: **#103**, the alarm created 2026-09-14 11:27Z ("nothing has landed in 16h"); **#102** has expired into the pile (superseded by #103). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-13 one: close #94–#102** — each was a one-day event, superseded by next day's alarm, never actionable; leave #103 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/.last-output:372:218:- **github-issues triage run, 2026-09-13: still nothing proposable — the mule-alarm pile is now 9 straight days (#94–#102).** Thirteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-12 run: **#102**, today's alarm (created 11:09Z, "nothing has landed in 14h"); **#101** has expired into the pile (superseded by #102). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-12 one: close #94–#101** — each was a one-day event, superseded by next day's alarm, never actionable; leave #102 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/.last-output:373:224:- **github-issues triage run, 2026-09-12: still nothing proposable — the mule-alarm pile is now 8 straight days (#94–#101).** Twelve open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-11 run: **#101**, today's alarm (created 10:58Z, "nothing has landed in 14h"); **#100** has expired into the pile (supe
.mule/reports/log/20260923T181743Z-recipe-astarte-upstream-proposed.md:119:- **github-issues triage run, 2026-09-23: still nothing proposable — the mule-alarm pile is now 16 alarms (#94–#109).** Twenty open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-21 run: **#109**, yesterday's alarm (created 09-22 10:54Z, "nothing has landed in 14h"); **#108** has expired into the pile (superseded by #109). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-21 one: close #94–#108** — each was a one-day event, superseded by next day's alarm, never actionable; leave #109 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — now a ~17-day idle streak — if that is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/reports/log/20260923T181743Z-recipe-astarte-upstream-proposed.md:127:- **github-issues triage run, 2026-09-21: still nothing proposable — the mule-alarm pile is now 16 days (#94–#108).** Eighteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-20 evening run: **#108**, today's alarm (created 09-21 10:54Z, "nothing has landed in 15h"); **#107** has expired into the pile (superseded by #108). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-20 one: close #94–#107** — each was a one-day event, superseded by next day's alarm, never actionable; leave #108 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/reports/log/20260923T181743Z-recipe-astarte-upstream-proposed.md:135:- **github-issues triage re-run, 2026-09-20 (evening): nothing proposable, no change since this morning's 13:10 run.** Re-surveyed with the recipe command: the issue set is identical to the 2026-09-20 morning triage — **#94–#107** are still the daily `mule-alarm` noise (14 issues, zero comments, telemetry not code), **#93** aclhook comment rewrite is still mule-review with `8c61268` already pushed (its own recipe path), **#92** keyAgreement is still parked on a stable upstream v1.4.0 with its waiting-on.md row and the mule-upstream-watch escalation wired (an -rc does not satisfy it), **#78** FDO is still the milestone-4.0 design/investigation already escalated, and **#1** is untouched per standing instruction. No new issue has been filed since the morning run. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines proposed**. This morning's proposal stands unchanged: **close #94–#106** — each was a one-day event, superseded by the next day's alarm, never actionable — and leave **#107** (still the newest, no #108 yet) to self-expire. The queue idle streak since ~2026-09-04/05 persists; if that is not intentional the wedge is on the Pi (queue review = the fix, per the 13:10 line).
.mule/reports/log/20260923T181743Z-recipe-astarte-upstream-proposed.md:147:- **github-issues triage run, 2026-09-20: still nothing proposable — the mule-alarm pile is now 15 days (#94–#107).** Seventeen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-18 run: **#107**, yesterday's alarm (created 09-19 10:49Z, "nothing has landed in 14h"); **#106** has expired into the pile (superseded by #107). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-18 one: close #94–#106** — each was a one-day event, superseded by next day's alarm, never actionable; leave #107 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/reports/log/20260917T194954Z-recipe-docs-sync-proposed.md:421: - **github-issues triage run, 2026-09-17: still nothing proposable — the mule-alarm pile is now 12 straight days (#94–#105).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-16 run: **#105**, today's alarm (created 09-17 10:51Z, "nothing has landed in 14h"); **#104** has expired into the pile (superseded by #105). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-16 one: close #94–#104** — each was a one-day event, superseded by next day's alarm, never actionable; leave #105 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/reports/log/20260917T194954Z-recipe-docs-sync-proposed.md:423: - **github-issues triage run, 2026-09-16: still nothing proposable — the mule-alarm pile is now 11 straight days (#94–#104).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), plus the daily alarm. New since the 2026-09-15 run: **#104**, today's alarm (created 09-16 11:12Z, "nothing has landed in 14h"); **#103** has expired into the pile (superseded by #104). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-15 one: close #94–#103** — each was a one-day event, superseded by next day's alarm, never actionable; leave #104 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/reports/log/20260920T185539Z-compat-note-custom-action-policy-boundary-done.md:106: - **github-issues triage run, 2026-09-20: still nothing proposable — the mule-alarm pile is now 15 days (#94–#107).** Seventeen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-18 run: **#107**, yesterday's alarm (created 09-19 10:49Z, "nothing has landed in 14h"); **#106** has expired into the pile (superseded by #107). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-18 one: close #94–#106** — each was a one-day event, superseded by next day's alarm, never actionable; leave #107 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:97:$ git remote -v; echo "---branch---"; git branch --show-current; echo "---waiting-on---"; ls -la .mule/waiting-on.md 2>/dev/null; echo "---research dir---"; ls -la .mule/research/ 2>/dev/null; echo "---find parity docs---"; rg -l "upstream-parity" --glob '*.md' .mule/ docs/ 2>/dev/null
.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:102:---waiting-on---
.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:210:13:- **github-issues triage run, 2026-09-16: still nothing proposable — the mule-alarm pile is now 11 straight days (#94–#104).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), plus the daily alarm. New since the 2026-09-15 run: **#104**, today's alarm (created 09-16 11:12Z, "nothing has landed in 14h"); **#103** has expired into the pile (superseded by #104). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-15 one: close #94–#103** — each was a one-day event, superseded by next day's alarm, never actionable; leave #104 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:212:17:- **github-issues triage run, 2026-09-15: still nothing proposable — the mule-alarm pile is now 10 straight days (#94–#103).** Fourteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-13 run: **#103**, the alarm created 2026-09-14 11:27Z ("nothing has landed in 16h"); **#102** has expired into the pile (superseded by #103). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-13 one: close #94–#102** — each was a one-day event, superseded by next day's alarm, never actionable; leave #103 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:214:21:- **github-issues triage run, 2026-09-13: still nothing proposable — the mule-alarm pile is now 9 straight days (#94–#102).** Thirteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-12 run: **#102**, today's alarm (created 11:09Z, "nothing has landed in 14h"); **#101** has expired into the pile (superseded by #102). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-12 one: close #94–#101** — each was a one-day event, superseded by next day's alarm, never actionable; leave #102 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/reports/log/20260920T203030Z-recipe-milestones-proposed.md:216:27:- **github-issues triage run, 2026-09-12: still nothing proposable — the mule-alarm pile is now 8 straight days (#94–#101).** Twelve open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-11 run: **#101**, today's alarm (created 10:58Z, "nothing has landed in 14h"); **#100** has expired into the pile (superseded by #101). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-11 one: close #94–#100** — each was a one-day event, superseded by next day's alarm, never actionable; leave #101 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
=== recipes dir ===
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6 --jq '{ahead:.ahead_by, behind:.behind_by, total:.total_commits, files:[.files[]|.filename]}' 2>&1 | head -40
{"ahead":27,"behind":0,"files":[".env",".github/workflows/astarte-build-workflow.yaml",".github/workflows/astarte-end-to-end-test-workflow.yaml",".gitignore",".typos.toml","CHANGELOG.md","apps/astarte_appengine_api/mix.exs","apps/astarte_appengine_api/mix.lock","apps/astarte_appengine_api/test/support/helpers/database.ex","apps/astarte_appengine_api/test/support/helpers/database_v2.ex","apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/rpc/server.ex","apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/rpc/supervisor.ex","apps/astarte_data_updater_plant/mix.exs","apps/astarte_data_updater_plant/mix.lock","apps/astarte_data_updater_plant/test/support/database_test_helper.ex","apps/astarte_data_updater_plant/test/support/helpers/database.ex","apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex","apps/astarte_housekeeping/lib/astarte_housekeeping/realms/queries.ex","apps/astarte_housekeeping/mix.exs","apps/astarte_housekeeping/mix.lock","apps/astarte_housekeeping/priv/migrations/astarte/0004_create_session_key_type.sql","apps/astarte_housekeeping/priv/migrations/astarte/0005_create_to2_sessions_table.sql","apps/astarte_housekeeping/priv/migrations/astarte/0006_create_ownership_vouchers_table.sql","apps/astarte_housekeeping/priv/migrations/astarte/0007_add_expiry_to_ownership_vouchers.sql","apps/astarte_housekeeping/priv/migrations/realm/0020_drop_unconfirmed_devices_table.sql","apps/astarte_housekeeping/priv/migrations/realm/0021_add_device_id_to_ownership_voucher.sql","apps/astarte_housekeeping/priv/migrations/realm/0022_delete_vouchers.sql","apps/astarte_housekeeping/priv/migrations/realm/0023_drop_ownership_voucher_table.sql","apps/astarte_housekeeping/priv/migrations/realm/0024_drop_to2_session_table.sql","apps/astarte_housekeeping/priv/migrations/realm/0025_add_fdo_guid_to_devices.sql","apps/astarte_housekeeping/test/support/helpers/database.ex","apps/astarte_pairing/config/dev.exs","apps/astarte_pairing/config/test.exs","apps/astarte_pairing/lib/astarte_pairing/agent/agent.ex","apps/astarte_pairing/lib/astarte_pairing/agent/device_registration_request.ex","apps/astarte_pairing/lib/astarte_pairing/engine.ex","apps/astarte_pairing/lib/astarte_pairing/queries.ex","apps/astarte_pairing/lib/astarte_pairing_web/api_spec/schemas/agent.ex","apps/astarte_pairing/lib/astarte_pairing_web/api_spec/schemas/ownership_voucher.ex","apps/astarte_pairing/lib/astarte_pairing_web/controllers/fallback_controller.ex","apps/astarte_pairing/lib/astarte_pairing_web/controllers/fdo_onboarding_controller.ex","apps/astarte_pairing/lib/astarte_pairing_web/controllers/ownership_voucher_controller.ex","apps/astarte_pairing/lib/astarte_pairing_web/plug/fdo_session.ex","apps/astarte_pairing/lib/astarte_pairing_web/router.ex","apps/astarte_pairing/lib/astarte_pairing_web/views/error_view.ex","apps/astarte_pairing/lib/astarte_pairing_web/views/ownership_voucher_view.ex","apps/astarte_pairing/mix.exs","apps/astarte_pairing/mix.lock","apps/astarte_pairing/test/astarte_pairing/certs/engine_test.exs","apps/astarte_pairing/test/astarte_pairing_web/controllers/fdo_onboarding_controller_test.exs","apps/astarte_pairing/test/astarte_pairing_web/controllers/ownership_voucher_controller_test.exs","apps/astarte_pairing/test/support/cases/device.ex","apps/astarte_pairing/test/support/cases/fdo_session.ex","apps/astarte_pairing/test/support/helpers/database.ex","apps/astarte_pairing/test/support/helpers/device.ex","apps/astarte_pairing/test/support/helpers/fdo.ex","apps/astarte_pairing/test/test_helper.exs","apps/astarte_realm_management/lib/astarte_realm_management/application.ex","apps/astarte_realm_management/lib/astarte_realm_management/device_removal/core.ex","apps/astarte_realm_management/lib/astarte_realm_management/device_removal/device_remover.ex","apps/astarte_realm_management/lib/astarte_realm_management/device_removal/queries.ex","apps/astarte_realm_management/lib/astarte_realm_management/device_removal/scheduler.ex","apps/astarte_realm_management/mix.exs","apps/astarte_realm_management/mix.lock","apps/astarte_realm_management/test/astarte_realm_management/device_removal/device_remover_test.exs","apps/astarte_realm_management/test/astarte_realm_management/device_removal/scheduler_test.exs","apps/astarte_realm_management/test/support/helpers/database.ex","apps/astarte_trigger_engine/mix.exs","apps/astarte_trigger_engine/mix.lock","apps/astarte_trigger_engine/test/support/helpers/database.ex","astarte-dashboard/cypress/e2e/fdo_vouchers_page.cy.js","astarte-dashboard/package.json","astarte-dashboard/src/DeviceStatusPage/DeviceInfoCard.tsx","astarte-dashboard/src/DevicesPage.tsx","astarte-dashboard/src/FdoVoucherPage.tsx","astarte-dashboard/src/FdoVouchersPage.tsx","astarte-dashboard/src/RegisterDevicePage.tsx","astarte-dashboard/src/astarte-client/client.ts","astarte-dashboard/src/astarte-client/models/Device/index.ts","astarte-dashboard/src/components/Icon.tsx","astarte-dashboard/src/components/IntrospectionTable.tsx","astarte-dashboard/src/hooks/useFdo.ts","compose/traefik/traefik.yaml","doc/mix.exs","doc/pages/architecture/050-pairing_mechanism.md","doc/pages/user/035-register_device.md","docker-compose.yml","libs/astarte_adapters/mix.lock","libs/astarte_config/mix.exs","libs/astarte_config/mix.lock","libs/astarte_data_access/lib/astarte_data_access/device.ex","libs/astarte_data_access/lib/astarte_data_access/device/unconfirmed_device.ex","libs/astarte_data_access/lib/astarte_data_access/devices/device.ex","libs/astarte_data_access/lib/astarte_data_access/fdo/ownership_voucher.ex","libs/astarte_data_access/lib/astarte_data_access/fdo/queries.ex","libs/astarte_data_access/lib/astarte_data_access/fdo/to2_session.ex","libs/astarte_data_access/mix.exs","libs/astarte_data_access/mix.lock","libs/astarte_data_access/test/device_test.exs","libs/astarte_data_access/test/fdo/ownership_voucher/ownership_voucher_test.exs","libs/astarte_data_access/test/fdo/queries_test.exs","libs/astarte_data_access/test/support/database_test_helper.exs","libs/astarte_events/mix.exs","libs/astarte_events/mix.lock","libs/astarte_events/test/support/helpers/database_test_helper.ex","libs/astarte_fdo/config/test.exs","libs/astarte_fdo/lib/config.ex","libs/astarte_fdo/lib/config/base_url_host.ex","libs/astarte_fdo/lib/owner_onboarding.ex","libs/astarte_fdo/lib/owner_onboarding/session.ex","libs/astarte_fdo/lib/ownership_voucher/load_request.ex","libs/astarte_fdo/lib/ownership_voucher/ownership_voucher.ex","libs/astarte_fdo/lib/service_info.ex","libs/astarte_fdo/lib/to0.ex","libs/astarte_fdo/mix.lock","libs/astarte_fdo/test/astarte_fdo/config/base_url_host_test.exs","libs/astarte_fdo/test/astarte_fdo/config_test.exs","libs/astarte_fdo/test/astarte_fdo/onboarding/done_test.exs","libs/astarte_fdo/test/astarte_fdo/onboarding/prove_device_test.exs","libs/astarte_fdo/test/astarte_fdo/owner_onboarding/owner_onboarding_test.exs","libs/astarte_fdo/test/astarte_fdo/owner_onboarding/session_test.exs","libs/astarte_fdo/test/astarte_fdo/ownership_voucher/load_request_test.exs","libs/astarte_fdo/test/astarte_fdo/ownership_voucher_test.exs","libs/astarte_fdo/test/astarte_fdo/rendezvous/core_test.exs","libs/astarte_fdo/test/astarte_fdo/rendezvous/rv_to2_addr_test.exs","libs/astarte_fdo/test/astarte_fdo/service_info_test.exs","libs/astarte_fdo/test/astarte_fdo/to0_test.exs","libs/astarte_fdo/test/support/cases/fdo_session.ex","libs/astarte_fdo/test/support/helpers/database.ex","libs/astarte_fdo/test/support/helpers/fdo.ex","libs/astarte_fdo/test/test_helper.exs","libs/astarte_fdo_core/lib/rendezvous/rv_to2_addr.ex","libs/astarte_fdo_core/mix.exs","libs/astarte_fdo_core/mix.lock","libs/astarte_fdo_core/test/astarte_fdo_core/fdo/rendezvous/rv_to2_addr_test.exs","libs/astarte_generators/mix.exs","libs/astarte_generators/mix.lock","libs/astarte_rpc/mix.exs","libs/astarte_rpc/mix.lock","libs/astarte_rpc/test/support/helpers/database.ex","libs/astarte_secrets/lib/astarte_secrets/astarte_secrets.ex","libs/astarte_secrets/mix.exs","libs/astarte_secrets/mix.lock","libs/astarte_test_suite/lib/astarte/test_suite/helpers/realm.ex","libs/astarte_test_suite/mix.lock","tools/astarte_device_fleet_simulator/mix.exs","tools/astarte_e2e/Dockerfile","tools/astarte_e2e/mix.exs","tools/astarte_e2e/mix.lock","tools/astarte_export/mix.exs","tools/astarte_export/mix.lock","tools/astarte_import/mix.exs","tools/astarte_import/mix.lock"],"total":27}

Now let me verify whether the new upstream delta (rc.5→rc.6, 27 commits) touches any v3.0-relevant wire surface.
$ cd /root/astrate-mule && gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6 --jq '.files[] | select(.filename|test("non_negative_integer_or_unset_type|queries.ex|router.ex|CHANGELOG")) | "=== \(.filename) (\(.status), +\(.additions)/-\(.deletions)) ===\n\(.patch // "(no patch)")"' 2>&1 | head -120
=== CHANGELOG.md (modified, +36/-0) ===
@@ -7,6 +7,33 @@ and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.
 
 ## Unreleased
 
+## [1.4.0-rc.6] - 2026-09-28
+
+Includes changes from v1.3.4
+
+### Added
+
+- [fdo] Add `POST /v1/:realm_name/fdo/ownership_vouchers/:guid/to0` to run TO0 again for an
+  ownership voucher, refreshing how long the rendezvous server keeps serving its registration.
+  The resulting expiry is reported when listing the ownership vouchers of a realm.
+- [fdo] Add ownership voucher deletion
+
+### Changed
+
+- [fdo] Rename `ASTARTE_BASE_URL_DOMAIN` to `ASTARTE_BASE_URL_HOST`. Its value may now be
+  either a domain name or an IP address
+- [fdo] Deleting a device also deletes its ownership voucher.
+- [fdo] Allow specifying device_id on ownership voucher upload
+- [fdo] The device is now immediately registered on ownership voucher upload
+- [astarte_data_updater_plant] Improve RPC server reliability
+
+### Fixed
+
+- [dashboard] Add unknown status when device messages aren't properly consumed
+- [fdo] Accept single entry x5chain certificates on voucher upload
+- [fdo] Invalid vouchers are no longer stored on the database
+- [fdo] Ensure users cannot delete ownership vouchers belonging to other realms
+
 ## [1.4.0-rc.5] - 2026-08-20
 
 ### Fixed
@@ -96,6 +123,15 @@ and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.
 
 - [astarte_data_updater_plant] Use mississippi consumer for data updater processes
 
+## [1.3.4] - 2026-09-17
+
+### Fixed
+
+- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
+  uncollected binaries over time by forcing more frequent full sweep garbage collections
+  on them.
+- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.
+
 ## [1.3.3] - 2026-08-07
 
 ### Fixed
=== apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex (modified, +1/-0) ===
@@ -40,4 +40,5 @@ defmodule Astarte.Housekeeping.Realms.NonNegativeIntegerOrUnsetType do
   def dump(:unset), do: {:ok, :unset}
   def dump(n) when is_integer(n), do: {:ok, n}
   def dump(_), do: :error
+  # coveralls-ignore-stop
 end
=== apps/astarte_housekeeping/lib/astarte_housekeeping/realms/queries.ex (modified, +120/-73) ===
@@ -310,9 +310,6 @@ defmodule Astarte.Housekeeping.Realms.Queries do
          :ok <- create_simple_triggers_table(keyspace_name),
          :ok <- create_grouped_devices_table(keyspace_name),
          :ok <- create_deletion_in_progress_table(keyspace_name),
-         :ok <- create_unconfirmed_devices_table(keyspace_name),
-         :ok <- create_ownership_vouchers_table(keyspace_name),
-         :ok <- create_to2_sessions_table(keyspace_name),
          :ok <- insert_realm_public_key(keyspace_name, public_key_pem),
          :ok <- insert_realm_astarte_schema_version(keyspace_name),
          :ok <- insert_realm(realm_name, device_limit),
@@ -650,42 +647,6 @@ defmodule Astarte.Housekeeping.Realms.Queries do
     end
   end
 
-  defp create_unconfirmed_devices_table(keyspace_name) do
-    query = """
-    CREATE TABLE #{keyspace_name}.unconfirmed_devices (
-      device_id uuid,
-      created_at timestamp,
-      PRIMARY KEY (device_id)
-    );
-    """
-
-    with {:ok, %{rows: nil, num_rows: 1}} <- CSystem.execute_schema_change(query) do
-      :ok
-    end
-  end
-
-  defp create_ownership_vouchers_table(keyspace_name) do
-    query = """
-    CREATE TABLE #{keyspace_name}.ownership_vouchers (
-      guid blob,
-      voucher_data blob,
-      output_voucher blob,
-      replacement_guid blob,
-      replacement_rendezvous_info blob,
-      replacement_public_key blob,
-      key_name varchar,
-      key_algorithm int,
-      user_id blob,
-      status int,
-      PRIMARY KEY (guid)
-    );
-    """
-
-    with {:ok, %{rows: nil, num_rows: 1}} <- CSystem.execute_schema_change(query) do
-      :ok
-    end
-  end
-
   defp create_session_key_type(keyspace_name) do
     query = """
     CREATE TYPE #{keyspace_name}.session_key (
@@ -699,40 +660,6 @@ defmodule Astarte.Housekeeping.Realms.Queries do
     end
   end
 
-  defp create_to2_sessions_table(keyspace_name) do
-    query = """
-    CREATE TABLE #{keyspace_name}.to2_sessions (
-      guid blob,

→ Read .mule/for-giulio.md [limit=60]

$ cd /root/astrate-mule && grep -n "APICompatVersion\s*=" internal/realm/service.go | head -3; echo "--- last-output context ---"; sed -n '360,370p' .mule/.last-output | cut -c1-300; echo "--- gitignore tail ---"; tail -3 .mule/.gitignore
588:const APICompatVersion = "1.2.2"
--- last-output context ---
So this issue is not waiting on anyone to remember it. It has a row in [`.mule/waiting-on.md`](https://github.com/astrate-platform/astrate/blob/main/.mule/waiting-on.md), and the weekly `mule-upstream-watch` job now reads that file as its first step every run: when a stable `v1.4.0` tag exists upstr

That mechanism was added because it was missing here: #51 was parked on the same condition, the condition fired on 2026-08-31, and nothing connected the new upstream tags to the parked issue — it took a hand triage four days later to catch it.
--

$ cd /root/astrate-mule && cat .mule/waiting-on.md 2>&1 | head -60; echo "=== grep waiting-on in for-giulio tail ==="; grep -n "waiting-on\|upstream-watch" .mule/for-giulio.md | tail -10
cat: .mule/waiting-on.md: No such file or directory
=== grep waiting-on in for-giulio tail ===
196:- **github-issues triage run, 2026-09-18: still nothing proposable — the mule-alarm pile is now 13 straight days (#94–#106).** Sixteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the wa
208:- **github-issues triage run, 2026-09-17: still nothing proposable — the mule-alarm pile is now 12 straight days (#94–#105).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the wa
210:- **github-issues triage run, 2026-09-16: still nothing proposable — the mule-alarm pile is now 11 straight days (#94–#104).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the wa
--- gitignore tail ---
.impl.patch
hosts
research/

$ cd /root/astrate-mule && git ls-files | grep -i "parity" ; echo "--- tracked docs mentioning the research doc ---"; grep -rln "upstream-parity-2026-08-22" --include="*.md" . 2>/dev/null | grep -v "^./.mule/reports\|^./.mule/cron.log\|^./.mule/for-giulio" | head
docs/handoff/flow-parity-audit-2026-07-29.md
internal/appengine/formats_parity_test.go
internal/appengine/groups_parity_test.go
internal/engine/mqttv1_parity_test.go
--- tracked docs mentioning the research doc ---
./.mule/milestones.md

Verified enough to write the deliverable. Adding the escalation entry at the top of `for-giulio.md` (newest first, per that file's own convention):
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,17 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-10-03: v3.0's gate is unchanged, and this run's one new finding is not a gap in Astrate's code — it is that this milestone's declared source document can never be committed, so every run re-derives the milestone from the issue bodies.**
+  **Verified unchanged, no new gaps.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.6` (2026-09-30, prerelease) — so #92's parking condition, frozen by you on 2026-09-04 ("wait for a stable v1.4.0", an `-rc.N` does not satisfy it), is still unmet. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Of #47–#89 **all are closed**; the open non-alarm set is unchanged — **#93** aclhook comment (`mule-review`, its own recipe path), **#92** keyAgreement (parked), **#78** FDO (`milestone-4.0`), **#1** (untouched per standing instruction). rc.6 adds no v3.0 gap: it is already audited in the two entries below, and I re-measured its non-FDO half myself (`compare/v1.4.0-rc.5...v1.4.0-rc.6`, 27 commits) — outside FDO the only library changes are a `# coveralls-ignore-stop` line in `apps/astarte_housekeeping/lib/astarte_housekeeping/realms/non_negative_integer_or_unset_type.ex:43` and `realms/queries.ex` *removing* the ownership_voucher / to2_sessions / unconfirmed_devices table creation (moved into migrations). No Realm Management, AppEngine or Pairing wire surface moved.
+  **The recipe's standing item is satisfied and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental*, reconcile trigger "upstream 1.4 final". 1.4 is not final, so there is nothing to promote and nothing to deprecate, and the register is accurate as written. `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked this run) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision ("bumped only in the same change that completes that surface, after reconciling every row tagged with the level"). The milestone-bump issue the recipe asks me to file **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": v3.0's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is delivered, the 1.4 half is open (#92, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut.
+  **New, and the reason this run is not a copy of the twenty before it: `.mule/research/` is listed in `.mule/.gitignore`, so the source document this milestone is built on — `.mule/research/upstream-parity-2026-08-22.md`, named at `.mule/milestones.md:114` — cannot exist in any clone, and the Pi's copy is gone.** Evidence: the last line of `.mule/.gitignore` is `research/`, a bare directory pattern that matches `.mule/research/` at any depth; the directory does not exist here; `git ls-files .mule` has no `research` entry; and `grep -rl 'upstream-parity-2026-08-22' --include='*.md'` over the tree matches **exactly one file, `.mule/milestones.md:114` itself** — there is no tracked equivalent (the only parity document under `docs/handoff/` is the v2.0 `flow-parity-audit-2026-07-29.md`). This is the sharper form of the 2026-09-06 note ("the section's source document does not exist in the repo"), and it is also why that note never got resolved: dropping the file back in place as-is re-dirties the tree on the next tick, because `.mule/.gitignore`'s own header states that any of its unlisted siblings "makes the tree dirty, which makes `tick` refuse to run and so turns the timer into a permanent silent no-op". **Yours, three ways:** (a) commit it once behind a negation (`!.mule/research/upstream-parity-2026-08-22.md`), so the milestone keeps a written reference; (b) move the source document to a tracked path (`docs/handoff/upstream-parity-2026-08-22.md`) and fix the `:114` pointer to match; (c) declare it lost and let the #47–#89 issue bodies plus the per-release CHANGELOG diffs stand as the reference. I did not pick one — every option edits a file that carries your decisions, and (a) versus (b) differ on a question only you can answer: is that document mutable working state or a frozen record?
+  **Read together with the 2026-09-26 `waiting-on.md` entry, not as a repeat of it: the same rule explains why that file still does not exist either.** `.mule/waiting-on.md` is *not* in `.mule/.gitignore`, so under the header's rule it cannot be created here without wedging `tick`, and it is not tracked (`git ls-files .mule` has no `waiting-on`), so it is not here at all — and nothing in the repo references it (`rg waiting-on` over `.mule/` and `tools/` matches only prose in this file and copies under `.mule/reports/log/`). So #92's escalation is currently parked on a condition the repository cannot express. That entry's one-liner — `systemctl list-timers 'mule-upstream*' --all` — is still what settles whether the mechanism lives in a Pi-side unit outside this repo; **I could not run it from this box** (no systemd here), so the 2026-09-26 finding stands unchanged rather than refuted, and this is only the repo-side half of the answer.
+  No `gh issue create`, nothing commented, closed or edited on GitHub, and **no `.mule/todo.md` lines**: every open item already has an owner (#92 the waiting-on mechanism, #78 parked and escalated below, #93 its own recipe path), and the one new fact is a file-layout decision, not a machine-checkable code gap — filing or queueing it would put a choice that is yours in front of the queue as if it were settled.
+
+---
+
 - **The container block (Design B / #43) ships no metrics at all, and its 1 MiB response cap is not configurable — both are MVP items in your own handoff, both are feature-shaped, so neither is a queue line.** `docs/handoff/flow-design-b-container-block-2026-07-29.md:198-199` lists "Config timeouts, max response bytes" and "Metrics: messages in/out, container start failures, processing errors" as MVP work, and §9 tracks them under step B4; the package registers zero collectors (`internal/flow/blocks/container/` has no `prometheus` import at all), while every other flow component carries an `astrate_flow_*` family (`astrate_flow_router_submitted_total` et al., internal/flow/router.go:260-276) and the discovery API advertises the block as production-ready (`internal/flow/blocks/info.go:112-117`, `schema.go:19`). Related: `Bridge.MaxBodyBytes` (internal/flow/blocks/container/httpbridge.go:29) is never assigned — `New` builds the bridge without it (block.go:112-116), so the limit is hard-wired to `defaultMaxBodyBytes` and a tenant cannot raise it for a large payload; the field reads like a knob but is not one. Yours to call: metric names + whether they become part of the documented metrics surface, and whether `max_response_bytes` becomes a real config key (it would need the schema at `schema.go:19` and the `Config` string at `info.go:116` updated in the same change). Related queue lines exist for the mechanical half (`container-response-cap-test` pins today's cap); nothing implements either feature. Also from the same review, for the record: the handoff's B5 orphan sweep is still open as `container-orphan-boot-cleanup [legion]`, and the image-trust note in §7 is satisfied by the current `Info.Summary`, so a missing image allowlist is documented behaviour rather than a gap.
 
 - **Delta on the two open COMPATIBILITY.md wording entries further down this file (the v1.3.4 one folding in v1.3.3, and the original v1.3.2 one): upstream `v1.4.0-rc.6` is wire-inert for Astrate, so nothing in that proposal is re-derived — only its version reference and one line of its "not yet emulated" list change.** Not applied: `docs/COMPATIBILITY.md` is untouched (still the 1.2.2 target, deviation 18's "which upstream 404s" at `docs/COMPATIBILITY.md:360`, `APICompatVersion = "1.2.2"` at `internal/realm/service.go:588`), and the FDO half of the rc.6 delta already has its own entry below rather than being folded in here.
   - **The whole non-FDO delta is inert, item by item.** DUP RPC-server reliability is OTP supervisor/GC tuning — no protocol. `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"), i.e. FDO-only, covered by the entry below. `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma.

$ cd /root/astrate-mule && git status --short; echo "--- exit ---"; sed -n '11,16p' .mule/for-giulio.md | cut -c1-160
M .mule/for-giulio.md
--- exit ---
---

- **Milestone recipe run, 2026-10-03: v3.0's gate is unchanged, and this run's one new finding is not a gap in Astrate's code — it is that this milestone's de
  **Verified unchanged, no new gaps.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.4`
  **The recipe's standing item is satisfied and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries exactly its two rows — #67 required/encr
  **New, and the reason this run is not a copy of the twenty before it: `.mule/research/` is listed in `.mule/.gitignore`, so the source document this milestone

Only `.mule/for-giulio.md` changed; no code, no todo.md lines, no git mutations, no issues filed.

```
Done: Milestone recipe run on v3.0 — gate verified still closed (no stable upstream v1.4.0; rc.6 is wire-inert for v3.0), standing UPSTREAM-EXPERIMENTAL/APICompatVersion item confirmed already satisfied by #90, and one new escalation filed: `.mule/research/` is gitignored, so v3.0's declared source document can never be committed.
Files: .mule/for-giulio.md (one new entry at the top)
Verified: gh issue list --label milestone-3.0 (empty), upstream-parity #47–#89 (all closed), gh api releases + compare v1.4.0-rc.5...v1.4.0-rc.6 (27 commits, non-FDO delta inert), git ls-files / .mule/.gitignore / grep -rl for the missing research doc -> pass (findings confirmed, not fixed)
Unsure: could not run `systemctl list-timers 'mule-upstream*'` from this box (no systemd), so the 2026-09-26 waiting-on.md finding stands unrefuted rather than resolved
Follow-ups: none appended to .mule/todo.md — the two live findings (research doc unrecoverable, waiting-on.md inexpressible) are file-layout decisions recorded in for-giulio.md, and every open issue (#92/#93/#78) already has an owner
```
