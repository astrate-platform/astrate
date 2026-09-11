slug: recipe-astarte-upstream
verdict: proposed
at:  514a685
ran: 2026-09-11T19:15:41Z on DietPi in 204s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[entry] | select(.tag_name=="v1.4.0-rc.5") | .body' 2>/dev/null; gh api "repos/astarte-platform/astarte/releases/tags/v1.4.0-rc.5" --jq '.body'
Sixth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.4

### Fixed

- [astarte_data_updater_plant] Make object values only validate type on selected interface
- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated

$ ls -la .mule/ 2>/dev/null; ls .mule/tasks/ 2>/dev/null
total 2776
drwxr-xr-x  7 root root    4096 Sep 11 20:29 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      14 Sep 11 21:12 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 11 19:07 .heartbeat
-rw-r--r--  1 root root       0 Sep 11 21:12 .last-output
-rw-r--r--  1 root root       0 Sep 11 21:12 .lock
-rw-r--r--  1 root root      17 Sep 11 21:12 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 2676084 Sep 11 21:12 cron.log
drwxr-xr-x  2 root root    4096 Sep 11 20:29 failed
-rw-r--r--  1 root root   27682 Sep 11 21:04 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    4647 Sep 11 20:29 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 10 19:00 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   38735 Sep 11 20:29 todo.md
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md

→ Read .mule/todo.md

$ rg -n "astarte-upstream|upstream movement|upstream Astarte" .mule/log.md .mule/for-giulio.md | head -40
(no output)

$ ls .mule/reports/ 2>/dev/null; ls .mule/failed/ 2>/dev/null
hygiene.md
log
appengine-group-token-roundtrip-test.diff
appengine-group-token-roundtrip-test.gate.log
appengine-snapshot-ignores-query-params.diff
broker-disconnect-device-zombie-session.diff
broker-offline-acl-tests.diff
broker-onconnect-doc-comment.diff
broker-onconnect-doc-comment.gate.log
compat-note-v133.diff
compat-note-v14-rc.diff
docs-sync-appengine-data-set-422.diff
docs-sync-rm-version-example.diff
empty-introspection-verification.diff
flow-validate-source-sink.diff
flow-validate-source-sink.gate.log
gate.log
probe-property-resend-encoding.diff
probe-props-resend-error-triggers.diff
probe-trigger-install-notification-delay.diff

$ rg -n -i "upstream" .mule/for-giulio.md | head -30
15:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
17:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
25:- **docs-sync pairing, 2026-09-10: docs/site names the pre-credentials wire status `registered`; the API returns `pending`.** `docs/site/pairing-and-security.md:56` ("flips status `registered -> confirmed`") and `:91` ("**registered** -- device registered, awaiting first credentials request"), plus `docs/site/data-modeling.md:94`, all give `registered` as a device status value. But `service.Info` emits `pending` for any device that is neither confirmed nor inhibited (internal/pairing/service.go:296-299 — upstream-parity per the comment at service.go:282-285; the DB value is `registered`, store/devices.go:19-27). Wire values are `pending`/`confirmed`/`inhibited`. Site prose is yours — reword to `pending`, or confirm the site intentionally describes the DB value.
48:  upstream gates.** `milestone-3.0` label remains empty. The three open `upstream-parity`
49:  issues are the same: #92 keyAgreement (upstream decision, gated on stable v1.4.0), #93
51:  UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait
52:  on upstream v1.4 final (currently rc.5). `APICompatVersion` stays 1.2.2 per #90's frozen
68:- **Milestone recipe run, 2026-09-06: v3.0 (upstream 1.2.2 → 1.3/1.4) has no open
69:  machine-checkable gaps left to file; the milestone now waits on two upstream gates, and the
70:  `milestone-3.0` label is empty (this milestone's work lives under `upstream-parity`).** All
75:  experimental decision, escalated 2026-09-05, gated on a stable upstream v1.4.0), **#78 FDO**
77:  `8c61268`). `docs/UPSTREAM-EXPERIMENTAL.md` rows #67/#68 are "1.4 experimental" and reconcile
78:  only when upstream ships **v1.4 final** — not yet (still rc). Per **#90's frozen decision**
80:  every UPSTREAM-EXPERIMENTAL row at that level is reconciled, so the bump is not yet due. Your
86:  upstream-parity-2026-08-22.md` is absent from `main` and from `origin/mule/research`.** The
109:  stable v1.4.0?** (issue #92, `upstream-parity`/`upstream-experimental`). #51 closed
110:  2026-08-22 with "parked until the upstream 1.4 experimental spec stabilizes — reopen or
111:  file fresh when it does". The document side has now stabilized: upstream `d084308`
132:- **COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable, 2026-07-14; v1.4.0 is still rc-only).** Astrate's doc and `APICompatVersion` still target upstream **1.2.2** (`internal/realm/service.go:588`); v1.3.0 (2026-05-06) introduced wire-surface changes Astrate does not yet emulate, so this is a decision — adopt v1.3.2 as the compatibility target (then update the doc + bump `APICompatVersion` together, per the bump rule) or keep 1.2.2 and add a "not yet emulated" note. Wire-relevant v1.3.0 deltas (release notes): **MQTT v1 capabilities** incl. `purge_properties_compression_format` (plaintext vs zlib purge — touches the `emptyCache`/`producer,properties` contract COMPATIBILITY.md deviation 1 documents); **empty introspection now allowed**; **device registration triggers** (pairing) and **device deletion started/completed triggers** (RM — the latter two already exist as Astrate deviation 9 emits both around the synchronous delete); **FDO authentication** (pairing, disabled by default); **realm-scoped health** — upstream v1.3 added `GET /pairing/v1/{realm}/health`, which Astrate already serves (`internal/pairing/http.go:78-81`, comment already says "upstream 1.3+"), so deviation 18's wording ("which upstream 404s") is now false against 1.3 and the note should be reworded either way. Proposed doc wording (for your approval, edit to taste): in §Infrastructure differences add a sentence — *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl.
134:plaintext `purge_properties_compression_format`, empty-introspection allowance, device registration/deletion triggers, experimental FDO pairing auth) are not yet emulated and are out of scope until the milestone that adopts v1.3.2 as the target."* — and reword deviation 18's realm-health note from "which upstream 404s" to "added by upstream v1.3 (Astrate serves it against a 1.2.2 target; kept, matching behavior)". Raw upstream changes: [v1.3.0](https://github.com/astarte-platform/astarte/releases/tag/v1.3.0), [v1.3.2](https://github.com/astarte-platform/astarte/releases/tag/v1.3.2).
170:  CLOSED (#23–#27, #37, #39–#43), no open issues, no new gaps after re-checking upstream
183:  upstream's own spec too (`docs/api/astarte_appengine_api.yaml` has no `/groups/{group}`
198:  mule task; the investigation work (reading upstream's TO2 handling, inventorying endpoints,
264:service only (last mile, like upstream), reuse fdo-rs for

→ Read .mule/log.md

$ git ls-remote --tags --refs https://github.com/astarte-platform/astarte.git 2>/dev/null | awk -F/ '{print $NF}' | sort -V | tail -8
v1.3.2
v1.3.3
v1.4.0-rc.0
v1.4.0-rc.1
v1.4.0-rc.2
v1.4.0-rc.3
v1.4.0-rc.4
v1.4.0-rc.5

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)\t\(.prerelease)"' 2>/dev/null; echo "---"; gh api repos/astarte-platform/astarte/releases --jq '(length)' 2>/dev/null
v1.4.0-rc.5	2026-08-20T12:28:04Z	true
v1.3.3	2026-08-07T14:33:58Z	false
v1.4.0-rc.4	2026-08-06T13:40:20Z	true
v1.4.0-rc.3	2026-07-31T10:48:39Z	true
v1.4.0-rc.2	2026-07-14T10:16:22Z	true
---
30

$ rg -n -i "astarte-upstream|recipe|proposal run|proposal job" .mule/for-giulio.md | head -20
13:- **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pile is now 7 straight days and worth a look.** Twelve open issues. Still no machine-checkable fix candidates, so no task lines: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-10 run: **#100**, today's alarm (created 10:55Z, "nothing has landed in 15h"). #99–#94 are the same `mule-alarm` (one per day 2026-09-05→09-10, ~11:00Z each, zero comments) — telemetry, not code issues, so never proposable. But seven consecutive alarms is past the "low-activity window" wording earlier runs used: the mule has landed no commit since ~2026-09-04/05, and `main`'s queue copy shows the top line as `- [ ]` while the landed work below is all `[!] BLOCKED` (`wrote nothing` / `tests failed`) — though `mule/queue` is the authoritative copy. **Proposal: close #94–#99** — each was a one-day event, superseded by next day's alarm, never actionable; leave #100 (live today) to self-expire. If the idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
15:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
17:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
43:  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
47:- **Milestone recipe run, 2026-09-09: v3.0 unchanged since 2026-09-06 — still waiting on
60:  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
68:- **Milestone recipe run, 2026-09-06: v3.0 (upstream 1.2.2 → 1.3/1.4) has no open
101:  awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the waiting-on
128:- **govulncheck GO-2026-5970: reachable DoS in golang.org/x/text (infinite loop on invalid input, fixed in v0.39.0, available v0.41.0).** Astrate pins `x/text` indirect at v0.38.0 (go.mod:97) and pgx pulls it into production: `internal/store/notify.go:59` `store.Listen` → `pgx.ConnectConfig` → `unicode/norm.*`. This is the only govulncheck symbol finding that is not test-harness-only: GO-2026-6355/6354 (x/crypto/ssh deadlocked-channel DoS) and GO-2026-6253 (moby/go-archive tar path traversal) are reachable only through testcontainers in `internal/testutil/pg.go`, i.e. never in the deployed binary. `x/text` keeps API compatibility minor-to-minor and the modules Astrate exercises (`unicode/norm` via pgx, `text/language` via jsonschema) are unchanged, so this is a fix Astrate actually needs — the hygiene recipe's highest-priority category. Not a mule task (go.mod never-touch): your decision to bump ≥v0.39.0 now or fold into the next milestone-boundary sweep. Raw: https://pkg.go.dev/vuln/GO-2026-5970. (The 2026-09-04 dep sweep did not list x/text.)
152:- **Dependency sweep corrected: direct (pinned) deps DO have newer versions** — the 2026-09-02 note said the `go list -m -u` sweep showed "only version-skew on transitive deps", but that run hit the recipe's `head -20` cutoff (all cloud/azure/transitive) and never reached the directly-required modules. Full sweep, 2026-09-04. None of these is a fix this repo *needs*, so no bump is proposed — recorded for the decision. Per module (current → available; breaking change; repo use):
173:  durable named flows all landed).   (Milestones recipe run, 2026-09-03.)

$ git log --oneline -15 -- .mule/todo.md 2>/dev/null
1757433 mule: blocked docs-sync-rm-version-example
206bdb9 mule: log docs-sync-rm-put-auth-422
a0de181 mule: log docs-sync-rm-mapping-required-encrypted
379ae97 mule: log docs-sync-rm-delete-interface-status
c093f9f mule: refill the queue
0ab77f3 mule: log docs-sync-pairing-version-endpoint
d9796eb mule: log docs-sync-pairing-status-enum
0e3d897 mule: refill the queue
407f0c8 mule: log flowapi-autorestart-default-test
565daa9 mule: log flowapi-validationdetail-test
39ed0c5 mule: log flowapi-autorestart-terminal-failure
bc1e53c mule: refill the queue
979872f mule: log docs-sync-native-compat-version-endpoints
7369a0d mule: log docs-sync-native-compat-health-503
d1188bd mule: refill the queue

No upstream movement to report. Here is the check.

**Newest release:** `v1.4.0-rc.5` (2026-08-20); newest stable `v1.3.3` (2026-08-07) — confirmed via both the releases API and `git ls-remote`, which ends at `v1.4.0-rc.5`. The v1.4.0-rc line and the v1.3.x surface were already reviewed and proposed by earlier runs: the queue already carries `[x] probe-required-mapping-flag`, the v1.3.0 probes (`purge-properties-compression-capability`, `empty-introspection-verification`, `probe-trigger-install-notification-delay`), `probe-props-resend-error-triggers` (v1.3.3), and the `compat-note-v1.3.2`/`compat-note-v1.3.3`/`compat-note-v1.4-rc` lines. Nothing shipped since the last run, so nothing new qualifies to propose.

```
Done: no upstream movement since v1.4.0-rc.5 (newest stable v1.3.3, unchanged since 2026-08-07); nothing appended to .mule/todo.md
Files: none
Verified: gh api repos/astarte-platform/astarte/releases; git ls-remote --tags -> pass (newest tag v1.4.0-rc.5)
Unsure: nothing
Follow-ups: none
```
