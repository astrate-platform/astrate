slug: recipe-milestones
verdict: proposed
at:  4fff71d
ran: 2026-09-23T17:43:41Z on DietPi in 324s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md
→ Read .mule/todo.md
✱ Glob ".mule/**/*.md" 0 matches

$ git -C /root/astrate-mule log --oneline -5 2>&1; git -C /root/astrate-mule status --short --branch 2>&1 | head -20
4fff71d mule: log examples-echo-container-contract-test
f35af21 mule: examples-echo-container-contract-test [auto]: add a container-free httptest suite for examples/flow-container-echo — the only no-test package in the repo that contains logic (`handleMessage`, main.go:49-75; `sanitizeLog`, main.go:40-47) — pinning (1) POST `/v1/message` echoes a valid JSON-object body verbatim with 200, answers 400 on invalid JSON and 204 on an empty body (the 1 MiB read cap), and (2) metadata `echo_drop` == "1" forces the 204 filter/drop path; the example is the canonical reference container contract (README.md, docs/handoff/flow-design-b-container-block-2026-07-29.md), so a regression silently misleads every container author who copies it.
82ae4da mule: refill the queue
5862ce6 mule: recipe hygiene ran (619s)
b8fbc71 mule: recipe docs-sync ran (238s)
## mule/queue

$ ls -la /root/astrate-mule/.mule/ 2>&1; echo ---; ls -la /root/astrate-mule/.mule/research/ 2>&1; echo ---; ls -la /root/astrate-mule/.mule/recipes/ 2>&1
total 6328
drwxr-xr-x  7 root root    4096 Sep 23 19:07 .
drwxr-xr-x 17 root root    4096 Sep 23 19:04 ..
-rw-r--r--  1 root root      13 Sep 23 19:38 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 23 19:07 .heartbeat
-rw-r--r--  1 root root       0 Sep 23 19:38 .last-output
-rw-r--r--  1 root root       0 Sep 23 19:38 .lock
-rw-r--r--  1 root root      11 Sep 23 19:38 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 6244715 Sep 23 19:38 cron.log
drwxr-xr-x  2 root root    4096 Sep 22 12:58 failed
-rw-r--r--  1 root root   54136 Sep 21 22:34 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    8538 Sep 23 19:07 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 21 21:26 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   75394 Sep 23 19:07 todo.md
---
ls: cannot access '/root/astrate-mule/.mule/research/': No such file or directory
---
total 56
drwxr-xr-x 2 root root 4096 Sep  4 20:46 .
drwxr-xr-x 7 root root 4096 Sep 23 19:07 ..
-rw-r--r-- 1 root root 2772 Jul 27 05:27 astarte-upstream.md
-rw-r--r-- 1 root root 3184 Jul 27 05:27 benchmarks.md
-rw-r--r-- 1 root root 2541 Jul 27 05:27 code-review.md
-rw-r--r-- 1 root root 5525 Jul 27 22:18 cross-project-survey.md
-rw-r--r-- 1 root root 2040 Jul 27 05:27 docs-sync.md
-rw-r--r-- 1 root root 2116 Jul 27 05:27 github-issues.md
-rw-r--r-- 1 root root 2103 Jul 27 05:27 hygiene.md
-rw-r--r-- 1 root root 4497 Jul 27 10:28 legion-go.md
-rw-r--r-- 1 root root 5741 Sep  4 20:46 milestones.md

$ git -C /root/astrate-mule branch -a 2>&1; echo ---; git -C /root/astrate-mule log --all --oneline -- .mule/research/ 2>&1 | head; echo ---; git -C /root/astrate-mule ls-tree -r --name-only main .mule/research/ 2>&1 | head; echo ---; git -C /root/astrate-mule ls-tree main .mule/ 2>&1 | head -30
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
916796b survey: triage the 2026-07-28..2026-09-04 backlog
6d8cae6 mule: daily cross-project survey, 2026-09-04
580faf1 mule: daily cross-project survey, 2026-09-02 (no material change)
8d44f59 mule: daily cross-project survey, 2026-09-01 (no material change)
1c618eb mule: daily cross-project survey, 2026-08-31
ce759e1 survey: recover 16 daily reports the broken commit check left staged (2026-07-28..2026-08-26), and take the tooling fix from main
f92d42a mule/research: commit the calibration survey, and mark it triaged
c6e0301 mule/research: seed the branch the daily survey timer will commit to
---
---
100644 blob 867ac3dc1028c247f86875e0cdf6bd810508afe0	.mule/.gitignore
100644 blob b1eafa1e4c18106a6da0cd3d3365d0d508abc64f	.mule/MULE.md
100644 blob 19aaef88642040428dddd9eb7ae59ba716aa991a	.mule/config
100644 blob 67bd57a9a90a34909a46050ad4619c3f61c9f39b	.mule/for-giulio.md
100644 blob ece84824fb45d6e4de612f8fbc0bcfda57a9423f	.mule/hosts.example
100644 blob c2486ae5c1739229c51e241323e0a725e1c8ec36	.mule/log.md
040000 tree ba5c0a71a58e88bdcefaedccbf0caccc3348a3f1	.mule/recipes
100644 blob 521ba621139cdd91ecabca3ee350d2efd66d4c36	.mule/todo.md

$ git -C /root/astrate-mule ls-tree -r --name-only remotes/origin/mule/research | grep -i upstream 2>&1 | head; echo ---; git -C /root/astrate-mule log --all --oneline --source -- '*upstream-parity*' 2>&1 | head
.mule/recipes/astarte-upstream.md
test/conformance/upstream/README.md
test/conformance/upstream/channels.json
test/conformance/upstream/channels.transcript.txt
test/conformance/upstream/record/main.go
test/conformance/upstream/recordchannels/main.go
test/conformance/upstream/rest-errors.json
test/conformance/upstream/rest-errors.transcript.txt
---

$ git -C /root/astrate-mule log --all --name-only --format= 2>/dev/null | grep -i '.mule/research' | sort -u 2>&1 | head; echo ---; git -C /root/astrate-mule ls-tree -r --name-only remotes/origin/mule/research .mule/research 2>&1 | head -40
.mule/research/README.md
.mule/research/issues-2026-07-27.md
.mule/research/issues-2026-07-28.md
.mule/research/issues-2026-07-30.md
.mule/research/issues-2026-08-01.md
.mule/research/issues-2026-08-06.md
.mule/research/issues-2026-08-07.md
.mule/research/issues-2026-08-08.md
.mule/research/issues-2026-08-13.md
.mule/research/issues-2026-08-15.md
---
.mule/research/README.md
.mule/research/issues-2026-07-27.md
.mule/research/issues-2026-07-28.md
.mule/research/issues-2026-07-30.md
.mule/research/issues-2026-08-01.md
.mule/research/issues-2026-08-06.md
.mule/research/issues-2026-08-07.md
.mule/research/issues-2026-08-08.md
.mule/research/issues-2026-08-13.md
.mule/research/issues-2026-08-15.md
.mule/research/issues-2026-08-18.md
.mule/research/issues-2026-08-21.md
.mule/research/issues-2026-08-22.md
.mule/research/issues-2026-08-23.md
.mule/research/issues-2026-08-24.md
.mule/research/issues-2026-08-25.md
.mule/research/issues-2026-08-26.md
.mule/research/issues-2026-09-04.md
.mule/research/log.md
.mule/research/survey-2026-07-27.md
.mule/research/survey-2026-07-28.md
.mule/research/survey-2026-07-30.md
.mule/research/survey-2026-08-01.md
.mule/research/survey-2026-08-06.md
.mule/research/survey-2026-08-07.md
.mule/research/survey-2026-08-08.md
.mule/research/survey-2026-08-13.md
.mule/research/survey-2026-08-15.md
.mule/research/survey-2026-08-18.md
.mule/research/survey-2026-08-21.md
.mule/research/survey-2026-08-22.md
.mule/research/survey-2026-08-23.md
.mule/research/survey-2026-08-24.md
.mule/research/survey-2026-08-25.md
.mule/research/survey-2026-08-26.md
.mule/research/survey-2026-08-31.md
.mule/research/survey-2026-09-04.md
.mule/research/triaged.md

$ git -C /root/astrate-mule remote -v 2>&1; echo ---; git -C /root/astrate-mule log --all --format=%H --name-only 2>/dev/null | grep -i 'upstream-parity\|parity-2026' | sort -u | head; echo ---; gh repo view --json nameWithOwner -q .nameWithOwner 2>&1
origin	git@github.com:astrate-platform/astrate.git (fetch)
origin	git@github.com:astrate-platform/astrate.git (push)
---
---
astrate-platform/astrate

$ git -C /root/astrate-mule rev-list --all --objects 2>/dev/null | grep -i 'parity\|UPSTREAM-EXPERIMENTAL' | head; echo ---; rg -n 'upstream-parity|parity-2026|#47|#89|milestone' /root/astrate-mule/.mule/for-giulio.md | head -40
65b7fcb37e2d6a84750e069ca454290513c2c0dc docs/UPSTREAM-EXPERIMENTAL.md
cedef26d6865f963eef918b47e7a83c3d1587424 docs/handoff/flow-parity-audit-2026-07-29.md
2263b29cfd3f72aec2db6a7b868f1b55abcd298e internal/appengine/formats_parity_test.go
166954515a90dc826a0c8af49154c062e861fcee internal/appengine/groups_parity_test.go
21a4ced692f7a172151144a545632dd425d9b077 internal/engine/mqttv1_parity_test.go
---
17:- **github-issues triage run, 2026-09-21: still nothing proposable — the mule-alarm pile is now 16 days (#94–#108).** Eighteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-20 evening run: **#108**, today's alarm (created 09-21 10:54Z, "nothing has landed in 15h"); **#107** has expired into the pile (superseded by #108). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-20 one: close #94–#107** — each was a one-day event, superseded by next day's alarm, never actionable; leave #108 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
21:- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
25:- **github-issues triage re-run, 2026-09-20 (evening): nothing proposable, no change since this morning's 13:10 run.** Re-surveyed with the recipe command: the issue set is identical to the 2026-09-20 morning triage — **#94–#107** are still the daily `mule-alarm` noise (14 issues, zero comments, telemetry not code), **#93** aclhook comment rewrite is still mule-review with `8c61268` already pushed (its own recipe path), **#92** keyAgreement is still parked on a stable upstream v1.4.0 with its waiting-on.md row and the mule-upstream-watch escalation wired (an -rc does not satisfy it), **#78** FDO is still the milestone-4.0 design/investigation already escalated, and **#1** is untouched per standing instruction. No new issue has been filed since the morning run. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines proposed**. This morning's proposal stands unchanged: **close #94–#106** — each was a one-day event, superseded by the next day's alarm, never actionable — and leave **#107** (still the newest, no #108 yet) to self-expire. The queue idle streak since ~2026-09-04/05 persists; if that is not intentional the wedge is on the Pi (queue review = the fix, per the 13:10 line).
29:- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
37:- **github-issues triage run, 2026-09-20: still nothing proposable — the mule-alarm pile is now 15 days (#94–#107).** Seventeen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-18 run: **#107**, yesterday's alarm (created 09-19 10:49Z, "nothing has landed in 14h"); **#106** has expired into the pile (superseded by #107). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-18 one: close #94–#106** — each was a one-day event, superseded by next day's alarm, never actionable; leave #107 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
41:- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
45:- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
53:- **COMPATIBILITY.md wording update for upstream v1.3.4 (newest stable, 2026-09-18; v1.4.0 is still rc.5-only) — folds v1.3.3 and v1.3.4 into the open v1.3.2 proposal below; both are maintenance-only and wire-inert, so that proposal stands unchanged.** v1.3.3 (2026-08-07) was an empty-body patch and v1.3.4 (2026-09-18, `astarte_data_updater_plant` only — more frequent fullsweep GC on AMQPDataConsumer processes plus an RPC-availability fix; release body via `gh api repos/astarte-platform/astarte/releases` `v1.3.4`) introduce **no wire/API surface change**, so the v1.3.2 wording proposal stands complete: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"). The only delta a v1.3.4-aware doc carries is the version reference: the proposed §Infrastructure-differences sentence's "until the milestone that adopts v1.3.2 as the target" reads "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target". Raw: [v1.3.4](https://github.com/astarte-platform/astarte/releases/tag/v1.3.4).
57:- **github-issues triage run, 2026-09-18: still nothing proposable — the mule-alarm pile is now 13 straight days (#94–#106).** Sixteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-17 run: **#106**, today's alarm (created 09-18 10:54Z, "nothing has landed in 14h"); **#105** has expired into the pile (superseded by #106). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-17 one: close #94–#105** — each was a one-day event, superseded by next day's alarm, never actionable; leave #106 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
61:- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) shows a new **stable v1.3.4** (2026-09-18) — maintenance only (data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change), so no new gap for Astrate — and still **no stable v1.4.0**: newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
69:- **github-issues triage run, 2026-09-17: still nothing proposable — the mule-alarm pile is now 12 straight days (#94–#105).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-16 run: **#105**, today's alarm (created 09-17 10:51Z, "nothing has landed in 14h"); **#104** has expired into the pile (superseded by #105). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-16 one: close #94–#104** — each was a one-day event, superseded by next day's alarm, never actionable; leave #105 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
71:- **github-issues triage run, 2026-09-16: still nothing proposable — the mule-alarm pile is now 11 straight days (#94–#104).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), plus the daily alarm. New since the 2026-09-15 run: **#104**, today's alarm (created 09-16 11:12Z, "nothing has landed in 14h"); **#103** has expired into the pile (superseded by #104). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-15 one: close #94–#103** — each was a one-day event, superseded by next day's alarm, never actionable; leave #104 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
73:- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-13 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
75:- **github-issues triage run, 2026-09-15: still nothing proposable — the mule-alarm pile is now 10 straight days (#94–#103).** Fourteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-13 run: **#103**, the alarm created 2026-09-14 11:27Z ("nothing has landed in 16h"); **#102** has expired into the pile (superseded by #103). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-13 one: close #94–#102** — each was a one-day event, superseded by next day's alarm, never actionable; leave #103 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
77:- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-12 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
79:- **github-issues triage run, 2026-09-13: still nothing proposable — the mule-alarm pile is now 9 straight days (#94–#102).** Thirteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-12 run: **#102**, today's alarm (created 11:09Z, "nothing has landed in 14h"); **#101** has expired into the pile (superseded by #102). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-12 one: close #94–#101** — each was a one-day event, superseded by next day's alarm, never actionable; leave #102 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
81:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
85:- **github-issues triage run, 2026-09-12: still nothing proposable — the mule-alarm pile is now 8 straight days (#94–#101).** Twelve open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-11 run: **#101**, today's alarm (created 10:58Z, "nothing has landed in 14h"); **#100** has expired into the pile (superseded by #101). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-11 one: close #94–#100** — each was a one-day event, superseded by next day's alarm, never actionable; leave #101 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
87:- **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pile is now 7 straight days and worth a look.** Twelve open issues. Still no machine-checkable fix candidates, so no task lines: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-10 run: **#100**, today's alarm (created 10:55Z, "nothing has landed in 15h"). #99–#94 are the same `mule-alarm` (one per day 2026-09-05→09-10, ~11:00Z each, zero comments) — telemetry, not code issues, so never proposable. But seven consecutive alarms is past the "low-activity window" wording earlier runs used: the mule has landed no commit since ~2026-09-04/05, and `main`'s queue copy shows the top line as `- [ ]` while the landed work below is all `[!] BLOCKED` (`wrote nothing` / `tests failed`) — though `mule/queue` is the authoritative copy. **Proposal: close #94–#99** — each was a one-day event, superseded by next day's alarm, never actionable; leave #100 (live today) to self-expire. If the idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
89:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
91:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
93:- **Extension point: external-bus Intake (NATS JetStream), 2026-09-10.** `internal/broker/intake.go:56-62` carries a `TODO(extension point)` for a second `Intake` implementation backed by a durable external bus (NATS JetStream) for multi-instance deployment or restart survival. The frozen design decision is in `docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4`; the consumer must reproduce per-device ordering, deferred-ack backpressure, and QoS 0 drop semantics. Not tracked by any open issue. Your call: file a milestone issue or park until the roadmap calls for it.
95:- **Extension point: timescaledb_toolkit lttb downsampling, 2026-09-10.** `internal/store/store.go:139-143` carries a `TODO(extension point)` for switching the `Downsample` method from `time_bucket+avg` to the toolkit's `lttb()` when `timescaledb_toolkit` is present — the probe at `store.go:144-150` already records availability in `s.hasToolkit`, but `datastreams.go` always uses the default path. Design reference: `docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5`. Not tracked by any open issue. Your call: file a milestone issue or park.
99:- **docs-sync pairing, 2026-09-10: docs/site names the pre-credentials wire status `registered`; the API returns `pending`.** `docs/site/pairing-and-security.md:56` ("flips status `registered -> confirmed`") and `:91` ("**registered** -- device registered, awaiting first credentials request"), plus `docs/site/data-modeling.md:94`, all give `registered` as a device status value. But `service.Info` emits `pending` for any device that is neither confirmed nor inhibited (internal/pairing/service.go:296-299 — upstream-parity per the comment at service.go:282-285; the DB value is `registered`, store/devices.go:19-27). Wire values are `pending`/`confirmed`/`inhibited`. Site prose is yours — reword to `pending`, or confirm the site intentionally describes the DB value.
118:  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated. #1
122:  upstream gates.** `milestone-3.0` label remains empty. The three open `upstream-parity`
124:  aclhook comment (pushed commit `8c61268`, mule-review), #78 FDO (milestone-4.0). Both
135:  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated below.
143:  machine-checkable gaps left to file; the milestone now waits on two upstream gates, and the
144:  `milestone-3.0` label is empty (this milestone's work lives under `upstream-parity`).** All
145:  the 1.3-surface gaps from the #47–#89 backlog are delivered/closed — retention ceiling #72,
146:  alias/group #52–#54 #59, query formats #55 #56, wire capabilities #47–#49, validation /
150:  (milestone-4.0, escalated), **#93** aclhook comment rewrite (mule-review, pushed commit
159:- **`milestones.md` v3.0 names a source doc that does not exist: `.mule/research/
160:  upstream-parity-2026-08-22.md` is absent from `main` and from `origin/mule/research`.** The
161:  investigation's *output* (issues #47–#89) is intact and reachable, so milestone work is
163:  milestone paragraph re-pointed) before v4.0 drafting leans on it.
176:  row. #78 is the milestone-4.0 FDO design/investigation, already escalated below. #1
183:  stable v1.4.0?** (issue #92, `upstream-parity`/`upstream-experimental`). #51 closed
197:  2026-09-05 — a prior escalation from the 2026-09-04 milestone run was lost in the queue

$ gh issue list --repo astrate-platform/astrate --label upstream-parity --state all --limit 60 2>&1
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

$ gh issue list --repo astrate-platform/astrate --label milestone-3.0 --state all --limit 60 2>&1
(no output)

$ gh api 'repos/astarte-platform/astarte/releases?per_page=30' --paginate -q '.[] | [.tag_name, .published_at, .prerelease] | @tsv' 2>&1 | head -30
v1.3.4	2026-09-18T09:42:57Z	false
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
v1.2.0-rc.0	2024-06-11T13:02:33Z	true
v1.0.6	2024-04-23T21:15:57Z	false
v1.1.1	2023-11-15T15:29:57Z	false
v1.0.5	2023-09-26T10:27:57Z	false
v1.1.0	2023-06-20T18:59:09Z	false
v1.1.0-rc.0	2023-06-09T15:49:44Z	true
v1.1.0-alpha.0	2022-11-24T16:38:37Z	true
v1.0.4	2022-10-25T16:50:21Z	false
v1.0.3	2022-07-05T13:57:17Z	false

→ Read .mule/for-giulio.md [limit=20]

$ gh issue view 67 --repo astrate-platform/astrate --json number,title,state,labels,closedAt,stateReason 2>&1; echo ---; gh issue view 68 --repo astrate-platform/astrate --json number,title,state,labels,closedAt,stateReason 2>&1
{"closedAt":"2026-09-04T18:34:19Z","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"},{"id":"LA_kwDORmfsJs8AAAACxx3GJg","name":"upstream-experimental","description":"Adopted from an upstream experimental/unstable feature — review each upstream release","color":"d4c5f9"}],"number":67,"state":"CLOSED","stateReason":"COMPLETED","title":"Interfaces: decide handling of required and encrypted mapping fields (upstream 1.4)"}
---
{"closedAt":"2026-09-04T18:34:35Z","labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACte8SoA","name":"mule-blocked","description":"The mule tried this and could not do it; see its comment","color":"d93f0b"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"},{"id":"LA_kwDORmfsJs8AAAACxx3GJg","name":"upstream-experimental","description":"Adopted from an upstream experimental/unstable feature — review each upstream release","color":"d4c5f9"}],"number":68,"state":"CLOSED","stateReason":"COMPLETED","title":"Decide async_operation=false params vs documented always-sync deviation"}

→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ gh issue view 92 --repo astrate-platform/astrate --json number,title,state,labels,body,comments 2>&1 | head -60
{"body":"Issue #51 was closed on 2026-08-22 with: *\"Full keyAgreement/encrypted-mappings support stays parked until the upstream 1.4 experimental spec stabilizes — reopen or file fresh when it does.\"* Filing fresh: that condition has now partly fired.\n\n**What changed.** Upstream commit `d084308` (\"docs: add encrypted endpoints key agreement (#2067)\", 2026-08-31) added `doc/pages/architecture/082-key_agreement_protocol.md` — a 267-line full wire protocol — and rewrote the Key Agreement section of `080-mqtt-v1-protocol.md`.\n\n**Verified against upstream `master` today, not taken from the survey:**\n\n- `082-key_agreement_protocol.md` exists and is 13,832 bytes.\n- The sentences the old parking rested on — \"reserved and routed correctly, but the handshake protocol is not yet implemented\" and \"accepted and acknowledged by Astarte, but no response is sent until the feature is fully implemented\" — are **gone from both 080 and 082** (grep count 0 in each).\n- 082 now states the handshake \"is implemented on the Astarte side (`astarte_data_updater_plant`)\", with only device-side session-key usage still described as liable to evolve.\n\n**What the spec now pins down**, where before there was only a sketch: five QoS-2 topics `control/keyAgreement/0..4` (InitExchange, ExchangeResp, SecretHash, HashOk, ExchangeFailed — the 080 table moved from a single `control/keyAgreement` at QoS 1); CBOR message bodies with CDDL given; `alg` 0 = `ECDH_P256-HKDF_SHA256-AES_256_GCM`; a CBOR-wrapped COSE_Key pubkey plus a 32-byte HkdfSalt; session-scoped keys never reused after reconnect or clean reset; AES-256-GCM whole-document ciphertext; and on \"no key established yet\" Astarte logs and **discards** rather than disconnecting. `ExchangeFailed` error codes are enumerated (0 internal, 1 invalid argument, 2 hash mismatch, 3 unprocessable entity).\n\n**The open question is scope, not feasibility, so this is not a mule task.** The spec is published and implemented upstream, but still ships only in `v1.4.0-rc.5` — `v1.3.3` remains the newest stable tag, and Astrate targets 1.2.2. So \"stabilized\" is true of the document and false of the release. Three ways to go:\n\n1. Implement now against the rc spec — largest surface in the parity backlog (CBOR codec, X25519/P-256, HKDF, AES-256-GCM, a 5-state handshake machine, shared-secret persistence, five new error names), and it would be built against a spec that may still move.\n2. Keep parking until `v1.4.0` is a stable tag, and re-file then.\n3. Do the narrow piece only — see the separate mule issue for the stale ACL comment, which is worth fixing regardless of which of these is chosen.\n\nRecorded in `.mule/for-giulio.md`. Survey source: `survey-2026-08-31.md` on `mule/research`.","comments":[{"id":"IC_kwDORmfsJs8AAAABSo0f8A","author":{"login":"atsetilam"},"authorAssociation":"MEMBER","body":"Parked by Giulio on 2026-09-04: **wait for a stable \\`v1.4.0\\`.** The spec is complete but ships only in \\`v1.4.0-rc.5\\`, and building against an rc risks rebuilding it.\n\nSo this issue is not waiting on anyone to remember it. It has a row in [`.mule/waiting-on.md`](https://github.com/astrate-platform/astrate/blob/main/.mule/waiting-on.md), and the weekly `mule-upstream-watch` job now reads that file as its first step every run: when a stable `v1.4.0` tag exists upstream, it escalates to `.mule/for-giulio.md` naming this issue. An `-rc.N` does not satisfy the row.\n\nThat mechanism was added because it was missing here: #51 was parked on the same condition, the condition fired on 2026-08-31, and nothing connected the new upstream tags to the parked issue — it took a hand triage four days later to catch it.","createdAt":"2026-09-04T19:52:21Z","includesCreatedEdit":false,"isMinimized":false,"minimizedReason":"","reactionGroups":[],"url":"https://github.com/astrate-platform/astrate/issues/92#issuecomment-5545730032","viewerDidAuthor":true}],"labels":[{"id":"LA_kwDORmfsJs8AAAACbTsJIQ","name":"enhancement","description":"New feature or request","color":"a2eeef"},{"id":"LA_kwDORmfsJs8AAAACxi99jg","name":"upstream-parity","description":"Feature parity with upstream astarte-platform","color":"0052cc"},{"id":"LA_kwDORmfsJs8AAAACxx3GJg","name":"upstream-experimental","description":"Adopted from an upstream experimental/unstable feature — review each upstream release","color":"d4c5f9"}],"number":92,"state":"OPEN","title":"keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds"}

$ head -40 /root/astrate-mule/.mule/log.md 2>&1
# Mule log

One row per task attempt, written by `tools/mule.sh`. This is the record of what the cheap
layer can actually do — read it before deciding whether a kind of task is worth delegating.

`secs` is the honest signal: a task that used most of its 900s budget was too big.

| date | task | outcome | secs | note |
| --- | --- | --- | --- | --- |
| 2026-09-04 | purge-properties-compression-capabilityauto | done | 86s | 6f7a3d6 |
| 2026-09-04 | broker-acl-coldstart-introspection-miss | done | 629s | d201db4 |
| 2026-09-04 | issue-93 | done | 108s | 24ad5b8 |
| 2026-09-04 | broker-disconnect-device-zombie-session | blocked | 117s | wrote nothing |
| 2026-09-04 | broker-offline-acl-tests | blocked | 208s | wrote nothing |
| 2026-09-04 | broker-onconnect-doc-comment | blocked | 438s | tests failed: --- FAIL: TestMQTTSink_Retained (0.02s) |
| 2026-09-05 | empty-introspection-verification | blocked | 106s | wrote nothing |
| 2026-09-05 | probe-trigger-install-notification-delay | blocked | 308s | wrote nothing |
| 2026-09-05 | compat-note-v132 | done | 73s | 4bc3e1a |
| 2026-09-05 | docs-sync-pairing-health-path | done | 186s | d73e225 |
| 2026-09-05 | docs-sync-pairing-register-404 | done | 106s | 44a7cae |
| 2026-09-05 | swagger-httptest-coverage | done | 267s | ba79b34 |
| 2026-09-05 | probe-property-resend-encoding | blocked | 100s | wrote nothing |
| 2026-09-05 | compat-note-v133 | blocked | 44s | wrote nothing |
| 2026-09-05 | flow-validate-source-sink | blocked | 360s | tests failed: --- FAIL: TestMQTTSink_Retained (0.02s) |
| 2026-09-05 | flow-validate-dead-source-sink-recompute | done | 200s | 7445176 |
| 2026-09-06 | docs-sync-rm-datastream-retention-endpoint | done | 232s | 0574fd9 |
| 2026-09-06 | docs-sync-rm-interfaces-detailed-param | done | 352s | 4094145 |
| 2026-09-06 | housekeeping-tests | done | 455s | 7583a7d |
| 2026-09-06 | probe-props-resend-error-triggers | blocked | 331s | wrote nothing |
| 2026-09-06 | server-data-trigger-bus | done | 827s | e46ce11 |
| 2026-09-06 | docs-sync-appengine-by-alias-endpoints | done | 144s | 35359df |
| 2026-09-06 | docs-sync-appengine-group-endpoints | done | 365s | 6ba74b8 |
| 2026-09-06 | docs-sync-appengine-get-group-device | done | 114s | 1aa98d6 |
| 2026-09-07 | docs-sync-appengine-query-params-status | done | 547s | bff7e1b |
| 2026-09-07 | docs-sync-appengine-group-patch-status | done | 75s | 81ced8f |
| 2026-09-07 | docs-sync-appengine-data-422-interface-level | done | 153s | f53d1ee |
| 2026-09-07 | appengine-snapshot-ignores-query-params | blocked | 327s | wrote nothing |
| 2026-09-07 | appengine-group-token-roundtrip-test | blocked | 306s | lint failed: internal/appengine/groups_token_test.go:14:5: redefines-builtin-id: redefinition of the built-in function max (revive) |
| 2026-09-07 | docs-sync-appengine-data-output-params | done | 215s | 262fdd5 |
| 2026-09-08 | probe-required-mapping-flag | done | 426s | db9d06e |

The v3.0 section is the current milestone (v2.0 is DONE, v4.0 is a placeholder). Prior runs have established the v3.0 ledger lives under `upstream-parity`, and the file `.mule/research/upstream-parity-2026-08-22.md` doesn't exist (already escalated 2026-09-04, for-giulio.md lines 159-163). Let me confirm nothing has changed since the last milestone run (2026-09-21).
$ gh api 'repos/astarte-platform/astarte/releases?per_page=5' -q '.[] | [.tag_name, .published_at, .prerelease] | @tsv' 2>&1; echo '---tags---'; gh api 'repos/astarte-platform/astarte/tags?per_page=5' -q '.[].name' 2>&1
v1.3.4	2026-09-18T09:42:57Z	false
v1.4.0-rc.5	2026-08-20T12:28:04Z	true
v1.3.3	2026-08-07T14:33:58Z	false
v1.4.0-rc.4	2026-08-06T13:40:20Z	true
v1.4.0-rc.3	2026-07-31T10:48:39Z	true
---tags---
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1

$ gh issue list --repo astrate-platform/astrate --label upstream-parity --state open --limit 30 --json number,title,labels --jq '.[] | [.number, ([.labels[].name] | join(","))] | @tsv' 2>&1; echo '---newest issues---'; gh issue list --repo astrate-platform/astrate --state open --limit 8 --json number,title,createdAt --jq '.[] | [.number, .createdAt, .title] | @tsv' 2>&1
93	mule-review,upstream-parity
92	enhancement,upstream-parity,upstream-experimental
78	enhancement,milestone-4.0,upstream-parity
---newest issues---
109	2026-09-22T10:54:47Z	mule: nothing has landed in 14h
108	2026-09-21T10:54:46Z	mule: nothing has landed in 15h
107	2026-09-19T10:49:12Z	mule: nothing has landed in 14h
106	2026-09-18T10:54:09Z	mule: nothing has landed in 14h
105	2026-09-17T10:51:07Z	mule: nothing has landed in 14h
104	2026-09-16T11:12:14Z	mule: nothing has landed in 14h
103	2026-09-14T11:27:08Z	mule: nothing has landed in 16h
102	2026-09-13T11:09:08Z	mule: nothing has landed in 14h

State is unchanged since the last milestone run (2026-09-21): no stable v1.4.0 (newest stable `v1.3.4`, 2026-09-18, maintenance-only; newest overall `v1.4.0-rc.5`), `milestone-3.0` label empty, the three open `upstream-parity` items are all already escalated (#92, #93, #78), and both `UPSTREAM-EXPERIMENTAL` rows still wait on v1.4 final. No new gaps → no issues to file, no task lines. The recipe's output here is the for-giulio entry, matching every prior run's pattern.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-09-23: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-21 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
+
+---
+
 - **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`ASTRATE_FLOW_CONFIG`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
 
 ---
 

No new gaps this run. Summary:

**Milestone worked:** v3.0 (upstream Astarte parity 1.2.2 → 1.3/1.4) — the first non-DONE section; v2.0 is DONE, v4.0 is a TBD placeholder.

**Investigation** (recipe step 3):
- Upstream release sweep: no stable v1.4.0 — newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change); newest overall `v1.4.0-rc.5`. Nothing shipped since the 2026-09-21 run.
- `milestone-3.0` label: still empty (this milestone's ledger lives under `upstream-parity`).
- Open `upstream-parity` items: only already-escalated **#92** (keyAgreement, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** (aclhook comment, mule-review on its own recipe path), **#78** (FDO, milestone-4.0). #47–#89 rest all closed.
- Standing item, UPSTREAM-EXPERIMENTAL reconcile: both rows (#67 required/encrypted flags, #68 async_operation) are tagged "reconcile when upstream 1.4 final" — still rc-only, so no reconcile due.

**Buckets:** nothing new is machine-checkable-and-unfiled, nothing needs a new design decision, nothing is newly covered. No issues filed (the open items are all already escalated, and prior versions already exist — filing would duplicate).

**Writeups:** one `.mule/for-giulio.md` entry appended (newest-at-top, matching the established milestone-run pattern), recording that v3.0 is verified unchanged and still gated on upstream v1.4. Deliberately did **not** propose "complete, cut the tag": the milestone targets 1.3/1.4 and the 1.4 gate is unresolved.

No `.mule/todo.md` lines appended (nothing to propose), no git touched, no issues created.

Done: milestone recipe run — v3.0 verified unchanged, gated on stable v1.4.0, no action
Files: .mule/for-giulio.md
Verified: gh api releases sweep, gh issue list --label upstream-parity/milestone-3.0, UPSTREAM-EXPERIMENTAL.md reconcile -> pass (no changes needed)
Unsure: nothing
Follow-ups: none (open items #92/#93/#78 already tracked; 1.4 gate is the blocking condition)
