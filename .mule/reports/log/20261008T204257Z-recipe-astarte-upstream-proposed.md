slug: recipe-astarte-upstream
verdict: proposed
at:  1ca603c
ran: 2026-10-08T20:42:57Z on DietPi in 108s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.3.5	2026-10-05T14:11:44Z
v1.4.0-rc.6	2026-09-30T07:28:27Z
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z

→ Read .mule/todo.md

$ ls -la /root/astrate-mule/.mule/ && ls -la /root/astrate-mule/.mule/tasks/ 2>/dev/null; ls -la /root/astrate-mule/.mule/recipes/ 2>/dev/null
total 15632
drwxr-xr-x  7 root root     4096 Oct  8 22:06 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       14 Oct  8 22:41 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  8 21:36 .heartbeat
-rw-r--r--  1 root root        0 Oct  8 22:41 .last-output
-rw-r--r--  1 root root        0 Oct  8 22:41 .lock
-rw-r--r--  1 root root       17 Oct  8 22:41 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 15566956 Oct  8 22:41 cron.log
drwxr-xr-x  2 root root     4096 Oct  8 21:13 failed
-rw-r--r--  1 root root   159992 Oct  8 22:38 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    15944 Oct  8 21:36 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  8 19:28 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   170694 Oct  8 21:36 todo.md
total 60
drwxr-xr-x 2 root root 4096 Sep  7 21:11 .
drwxr-xr-x 7 root root 4096 Oct  8 22:06 ..
-rw-r--r-- 1 root root 3293 Sep  7 21:11 docs-sync-appengine-data-output-params.md
-rw-r--r-- 1 root root 5539 Jul 27 22:18 issue-16.md
-rw-r--r-- 1 root root 7510 Sep  4 20:46 issue-23.md
-rw-r--r-- 1 root root 6434 Sep  4 20:46 issue-24.md
-rw-r--r-- 1 root root 5665 Sep  4 20:46 issue-25.md
-rw-r--r-- 1 root root 6009 Sep  4 20:46 issue-26.md
-rw-r--r-- 1 root root 5683 Sep  4 20:46 issue-27.md
total 56
drwxr-xr-x 2 root root 4096 Sep  4 20:46 .
drwxr-xr-x 7 root root 4096 Oct  8 22:06 ..
-rw-r--r-- 1 root root 2772 Jul 27 05:27 astarte-upstream.md
-rw-r--r-- 1 root root 3184 Jul 27 05:27 benchmarks.md
-rw-r--r-- 1 root root 2541 Jul 27 05:27 code-review.md
-rw-r--r-- 1 root root 5525 Jul 27 22:18 cross-project-survey.md
-rw-r--r-- 1 root root 2040 Jul 27 05:27 docs-sync.md
-rw-r--r-- 1 root root 2116 Jul 27 05:27 github-issues.md
-rw-r--r-- 1 root root 2103 Jul 27 05:27 hygiene.md
-rw-r--r-- 1 root root 4497 Jul 27 10:28 legion-go.md
-rw-r--r-- 1 root root 5741 Sep  4 20:46 milestones.md

$ ls /root/astrate-mule/docs/ && grep -rn "v1\.3\|1\.2\.2" /root/astrate-mule/docs/UPSTREAM-EXPERIMENTAL.md 2>/dev/null | head -30
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

$ grep -n "v1\.3\|v1\.4\|1\.2\.2\|COMPATIBILITY" /root/astrate-mule/.mule/for-giulio.md | tail -60
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_11d4037e3001InowER1JCXvkdJ

106:  **The recipe's standing item is satisfied and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still carries exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental*, reconcile trigger "upstream 1.4 final". 1.4 is not final, so there is nothing to promote and nothing to deprecate, and the register is accurate as written. `APICompatVersion = "1.2.2"` (`internal/realm/service.go:588`, re-checked this run) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision ("bumped only in the same change that completes that surface, after reconciling every row tagged with the level"). The milestone-bump issue the recipe asks me to file **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": v3.0's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is delivered, the 1.4 half is open (#92, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut.
115:- **Delta on the two open COMPATIBILITY.md wording entries further down this file (the v1.3.4 one folding in v1.3.3, and the original v1.3.2 one): upstream `v1.4.0-rc.6` is wire-inert for Astrate, so nothing in that proposal is re-derived — only its version reference and one line of its "not yet emulated" list change.** Not applied: `docs/COMPATIBILITY.md` is untouched (still the 1.2.2 target, deviation 18's "which upstream 404s" at `docs/COMPATIBILITY.md:360`, `APICompatVersion = "1.2.2"` at `internal/realm/service.go:588`), and the FDO half of the rc.6 delta already has its own entry below rather than being folded in here.
119:  - **The two edits the proposal needs, both yours to approve:** (1) the 2026-09-18 entry's parenthetical "(v1.4.0 is still rc.5-only)" becomes "(v1.4.0 is still rc.6-only)" — re-verified this session that rc.6 is the newest 1.4 prerelease and there is still **no stable v1.4.0** (`gh api repos/astarte-platform/astarte/releases`); (2) its "**FDO authentication** (pairing, disabled by default)" item, and the same "experimental FDO pairing auth" phrase in the proposed doc sentence, understates the 1.4 line, which ships a full ownership-voucher surface (upload with mandatory `hw_id`, deletion, re-TO0, per-voucher `expiry`, a spec-mandated realm-less TO2 path) that Astrate does not emulate — #78, deliberately parked at `.mule/milestones.md:124`. "FDO ownership vouchers (1.4 line, not stable upstream)" says it without claiming the rc.5 shape.
120:  Unverified: the upstream side of the item list is the upstream-watch recipe's `v1.4.0-rc.5...v1.4.0-rc.6` comparison — the same measurement as the FDO entry below, taken there and not re-derived here (a raw fetch of `models/Device/index.ts` at the rc.6 tag did not resolve from this box). I re-verified only the Astrate-side claims above and the release list. No `.mule/todo.md` lines: the deliverable is the delta, and the proposal itself stays open.
124:- **Scope update for issue #78 (FDO): upstream `v1.4.0-rc.6` (prerelease, published 2026-09-30T07:28Z, CHANGELOG header says 2026-09-28) grew the surface #78 is scoped against, and the one path #78's body states as fact no longer exists upstream. Nothing implemented, nothing commented or edited on the issue, no `todo.md` lines — #78 is parked on a product call that is yours, and this only changes how big the eventual build is.** Measured this session with `gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.5...v1.4.0-rc.6`: **27 commits, 153 files** (the task line said 154), 10 added migrations under `apps/astarte_housekeeping/priv/migrations/{astarte,realm}/`, and the device-side scope moved out of the realm entirely — `router.ex` now reads `scope "/fdo/101"` with the comment "URL mandated by FDO specification", the `LogRealm` + `VerifyRealmExists` plugs are dropped from the `:fdo` pipeline, and `6c8d4b8` (`refactor(fdo)!`, body: "BREAKING CHANGE: the FDO TO2 URL in Pairing no longer includes the /v1/:realm/ prefix") deletes the `realm_name` path parameter from all seven TO2 `operation` specs in `fdo_onboarding_controller.ex`. **This breaking change is not in the rc.6 CHANGELOG entry** (which lists only the to0 route, voucher deletion, the `ASTARTE_BASE_URL_*` rename, device-deletion cascade, `device_id` on upload, immediate registration, the two x5chain/voucher-store fixes and the cross-realm delete fix) — so anyone tracking upstream by its changelog, including #78's body, misses it. The realm now comes from the voucher row: the global-keyspace query `fetch_device_id_and_ownership_voucher_and_realm/1` returns `{realm, device_id, voucher_data}` and the `:fdo` plug fetches the session by GUID alone (`plug/fdo_session.ex`: `Session.fetch(guid)`); `RvTO2Addr` no longer carries the realm in the hostname either ("the realm name is no longer referenced in the hostname to be put inside RvTO2Addr structure"). **Yours:** #78's acceptance is "a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end" — against rc.6 that SDK posts to `/fdo/101/msg/{60..70}` with no realm, so mirroring rc.5's shape would mean mirroring a path upstream has deleted. No action while parked; the body is just stale.
134:- **Not client-observable, so not part of the decision:** the keyspace move itself — vouchers and TO2 sessions now live in the global `astarte` keyspace (`Realm.astarte_keyspace_name()` throughout `fdo/queries.ex`) with a `realm` column on the voucher row and `ALLOW FILTERING` on the list query, so **realm isolation is now application-level rather than keyspace-level** (which is what makes `ensure_voucher_in_realm`-style cross-realm checks necessary, `b3bf731` "ensure users cannot delete ownership vouchers belonging to other realms"); `fdo_guid` added to the device row (migration `realm/0025`) so device removal can find the voucher; the `unconfirmed_devices` drop (`realm/0020`, schema file deleted, `remove_device_from_unconfirmed_devices!` removed); the dashboard voucher page (`astarte-dashboard/src/FdoVoucherPage.tsx`, `hooks/useFdo.ts`, a new cypress spec); and the `data_updater_plant` RPC/GC fixes plus the `v1.3.4` forward-ports.
139:- **github-issues triage run, 2026-10-01: nothing proposable, and the 2026-09-26 reading of the alarm pile is now superseded — this silence is real, and the dead-man's switch is latched off so it cannot report it.** 25 open issues, the same set plus **#113** (09-27 11:07Z) and **#114** (09-27 20:09Z, "nothing has landed in 8h"). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction.
152:- **github-issues triage run, 2026-09-26: nothing proposable (as every run since 2026-09-05), and the 19-issue alarm pile is a false alarm — it says the mule has landed nothing in 19 days and it has landed something every day.** 23 open issues, the same set plus **#112** (today, 09-26 11:04:46Z, "nothing has landed in 14h"). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction.
157:  **Correction on a fact 19 entries have repeated as settled: there is no `.mule/waiting-on.md`.** #92's own comment (2026-09-04) states the issue "has a row in `.mule/waiting-on.md`" and that the weekly `mule-upstream-watch` job reads that file as its first step. The file does not exist here, and `rg 'waiting-on|upstream-watch'` over the whole tree — including gitignored `.mule/` and `tools/` — matches only prose in this file and copies of it under `.mule/reports/log/`. So the escalation that was supposed to catch a stable upstream `v1.4.0` is not in the repo: either it lives in a Pi-side unit outside it, or it was never wired. One command settles it: `systemctl list-timers 'mule-upstream*' --all; grep -rl 'upstream-watch\|waiting-on' /etc/systemd/system /root 2>/dev/null`. If it was never wired, #92 is parked on a condition nothing is watching, and the 2026-09-04 claim that "the mechanism was added because it was missing here" is false. The gate is being caught by hand in the meantime — today's milestone run re-derived it from `gh api .../releases`.
161:- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Status: not started", which is no longer true.** Release sweep (`gh api repos/astarte-platform/astarte/releases`, run this session) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.5` (2026-08-20, prerelease) — byte-identical to what the 2026-09-24 run saw, so nothing shipped upstream in the two days since. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Open non-alarm issues remain exactly **#93** (aclhook comment, `mule-review`, its own recipe path), **#92** (keyAgreement, parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** (FDO, `milestone-4.0`), **#1** (untouched per standing instruction). No new gaps, **no issues filed, no `.mule/todo.md` task lines added** — there is nothing machine-checkable left, and filing anything here would duplicate #92/#93/#78.
162:  **New, and the reason this run is not a copy of the last thirteen: `.mule/milestones.md:127-128` still reads "Status: not started. First recipe job: triage #47–#89 into an ordered plan (which are audits vs features vs decisions), file sub-issues where work splits, escalate the 'decide' set in one batch."** All of that happened: of #47–#89 every issue is closed except the deliberately-parked #78. The section describes as future work a triage that is finished, and it is the exact text every future run of this recipe reads to pick its target — so the recipe has been re-deriving "the backlog is done, the gate is upstream" from scratch on a timer for three weeks. Suggest rewriting those two lines to state the real position: *1.3 surface delivered (#47–#89 closed); milestone blocked on upstream v1.4.0 going stable; final phase = reconcile the two `docs/UPSTREAM-EXPERIMENTAL.md` rows, answer #92, bump `APICompatVersion`.* That also makes the blocking condition greppable instead of inferred.
163:  **The recipe's standing item is satisfied, and I did not re-file it.** `docs/UPSTREAM-EXPERIMENTAL.md` still has exactly its two rows — #67 required/encrypted mapping fields and #68 `async_operation=false`, both tagged *1.4 experimental* with reconcile trigger "upstream 1.4 final". v1.4 is not final, so both rows stay as written: nothing to promote, nothing to deprecate, and the register is accurate. `APICompatVersion` is still `"1.2.2"` (`internal/realm/service.go:588`) with the bump rule inline on the constant, per #90's frozen 2026-08-23 decision; the milestone-bump issue the recipe asks for **already exists as #90**, so filing another would be a duplicate. Still deliberately **not** proposing "milestone looks complete, verify and cut the tag": the milestone's declared scope is 1.2.2 → 1.3/**1.4**, the 1.3 half is done, and the 1.4 half is unresolved (#92 open, both experimental rows unreconciled, version not bumped), so it is not complete and the tag should not be cut. Also still true and already escalated on 2026-09-06, so not repeated: the section's source document `.mule/research/upstream-parity-2026-08-22.md` does not exist in the repo.
164:  **Standing ask, unchanged, for whenever a stable v1.4.0 ships:** answer #92, reconcile the two experimental rows, run the final-phase `APICompatVersion` bump (and update `docs/COMPATIBILITY.md` deviation #10 to match), then cut the v3.0 tag. One process note while this stays blocked: this is the fifteenth milestone-run entry in this file (fourteen prior, 2026-09-06 → 2026-09-24) and the header asks for a queue rather than a log — worth deleting the superseded duplicates whenever you next prune, or the signal-to-noise here keeps falling.
232:- **github-issues triage run, 2026-09-25: still nothing proposable — the mule-alarm pile is 18 issues (#94–#111), and two corrections to what earlier runs inferred about it.** Twenty-two open issues, same set as every run since 2026-09-05: **#111–#94** are the daily `mule-alarm` noise (zero comments, telemetry not code), **#93** aclhook comment (mule-review, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable upstream v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-24 run: **#111** (today, 09-25 11:01:07Z, "nothing has landed in 16h"). Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines added**. **Proposal, extending the 2026-09-24 one: close #94–#110**; leave #111 (live today) to self-expire. Two corrections, both from reading `tools/mule.sh` rather than the issue titles: **(1) the "missed days" in earlier entries (09-15, 09-20, 09-23) are evidence of activity, not its absence** — I previously called them timer gaps; the alarm is gated on a `.alarmed` sentinel that `beat` deletes on every land (`mule.sh:737,748-749`), so a day with no new issue is a day something landed, not a day the timer slept. **(2) the alarm body's closing line is wrong**: "Close it — the mule reopens a new one if the silence continues" (`mule.sh:764`) implies issue state drives the alarm; it does not. Nothing reads or reconciles the 18 open issues, so the pile can only grow and closing the old ones is safe but inert — and `cmd_refill` commits `mule: refill the queue` (mule.sh:719-720) **without calling `beat`**, so a doc-only refill never counts as landing work. The titles' ages also imply a beat roughly every day around 19:00–21:00Z (#111 filed 09-25 11:01Z at 16h → last beat ≈ 09-24 19:01Z; #110 filed 09-24 10:54Z at 15h → ≈ 09-23 19:54Z) while the ~11:00Z daily check has never once found under the 8h threshold — **what calls `beat` at that time is unverified**, and one command on the Pi settles it (`grep -E 'landed|checked:' /root/astrate-mule/.mule/log | tail -20`). The queue has still landed nothing since ~2026-09-04/05; if that idle streak is not intentional, the wedge is on the Pi and a queue review is the fix.
236:- **Milestone recipe run, 2026-09-24: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-23 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
240:- **github-issues triage run, 2026-09-24: still nothing proposable — the mule-alarm pile is now 17 days (#94–#110).** Twenty-one open issues, same set: **#93** aclhook comment (mule-review, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable upstream v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-20 runs: the alarms **#108** (09-21), **#109** (09-22) and **#110** (today, 09-24 10:54Z, the live one) — a four-day gap on 09-23 not withstanding. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines added**. **Proposal, extending the 2026-09-20 one: close #94–#109** — each was a one-day event, superseded by the next day's alarm, never actionable; leave #110 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the wedge is on the Pi and a queue review is the fix.
244:- **Milestone recipe run, 2026-09-23: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-21 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
248:- **github-issues triage run, 2026-09-23: still nothing proposable — the mule-alarm pile is now 16 alarms (#94–#109).** Twenty open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-21 run: **#109**, yesterday's alarm (created 09-22 10:54Z, "nothing has landed in 14h"); **#108** has expired into the pile (superseded by #109). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-21 one: close #94–#108** — each was a one-day event, superseded by next day's alarm, never actionable; leave #109 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — now a ~17-day idle streak — if that is not intentional, the queue wedge is on the Pi and a queue review is the fix.
256:- **github-issues triage run, 2026-09-21: still nothing proposable — the mule-alarm pile is now 16 days (#94–#108).** Eighteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-20 evening run: **#108**, today's alarm (created 09-21 10:54Z, "nothing has landed in 15h"); **#107** has expired into the pile (superseded by #108). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-20 one: close #94–#107** — each was a one-day event, superseded by next day's alarm, never actionable; leave #108 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
260:- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
264:- **github-issues triage re-run, 2026-09-20 (evening): nothing proposable, no change since this morning's 13:10 run.** Re-surveyed with the recipe command: the issue set is identical to the 2026-09-20 morning triage — **#94–#107** are still the daily `mule-alarm` noise (14 issues, zero comments, telemetry not code), **#93** aclhook comment rewrite is still mule-review with `8c61268` already pushed (its own recipe path), **#92** keyAgreement is still parked on a stable upstream v1.4.0 with its waiting-on.md row and the mule-upstream-watch escalation wired (an -rc does not satisfy it), **#78** FDO is still the milestone-4.0 design/investigation already escalated, and **#1** is untouched per standing instruction. No new issue has been filed since the morning run. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines proposed**. This morning's proposal stands unchanged: **close #94–#106** — each was a one-day event, superseded by the next day's alarm, never actionable — and leave **#107** (still the newest, no #108 yet) to self-expire. The queue idle streak since ~2026-09-04/05 persists; if that is not intentional the wedge is on the Pi (queue review = the fix, per the 13:10 line).
268:- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
272:- **COMPATIBILITY.md §8 wording proposal: trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (Forwarder) actions are single-shot by design, receiving only `maximum_capacity`/`event_ttl`.** Propose adding to the "Trigger delivery policies" deviation: *"Policy `error_handlers` (retry/discard) and the `retry_times` attempt cap govern HTTP webhook actions only: a custom action delivered through the Forwarder seam is single-shot by design, receiving from the policy solely its `maximum_capacity` in-flight bound and its `event_ttl` staleness drop, with a forward failure counted once and never retried."* Verified in `internal/engine/triggers/actions.go`: `Action.Custom` routes to `forward` (actions.go:435-439), which never consults `policy.Decide` — the Forwarder seam returns only an error, no HTTP status (actions.go:161-167, doc comment at actions.go:471-479) — while `Enqueue` bounds `maximum_capacity` on every delivery (actions.go:338-349) and `deliver` drops past `event_ttl` before the custom/webhook branch (actions.go:400-416). The existing warn-on-single-shot line already tells the operator when a policy's `retry_times` cannot apply to the forward path (actions.go:491-497).
276:- **github-issues triage run, 2026-09-20: still nothing proposable — the mule-alarm pile is now 15 days (#94–#107).** Seventeen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-18 run: **#107**, yesterday's alarm (created 09-19 10:49Z, "nothing has landed in 14h"); **#106** has expired into the pile (superseded by #107). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-18 one: close #94–#106** — each was a one-day event, superseded by next day's alarm, never actionable; leave #107 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
280:- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
284:- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
292:- **COMPATIBILITY.md wording update for upstream v1.3.4 (newest stable, 2026-09-18; v1.4.0 is still rc.5-only) — folds v1.3.3 and v1.3.4 into the open v1.3.2 proposal below; both are maintenance-only and wire-inert, so that proposal stands unchanged.** v1.3.3 (2026-08-07) was an empty-body patch and v1.3.4 (2026-09-18, `astarte_data_updater_plant` only — more frequent fullsweep GC on AMQPDataConsumer processes plus an RPC-availability fix; release body via `gh api repos/astarte-platform/astarte/releases` `v1.3.4`) introduce **no wire/API surface change**, so the v1.3.2 wording proposal stands complete: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"). The only delta a v1.3.4-aware doc carries is the version reference: the proposed §Infrastructure-differences sentence's "until the milestone that adopts v1.3.2 as the target" reads "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target". Raw: [v1.3.4](https://github.com/astarte-platform/astarte/releases/tag/v1.3.4).
296:- **github-issues triage run, 2026-09-18: still nothing proposable — the mule-alarm pile is now 13 straight days (#94–#106).** Sixteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-17 run: **#106**, today's alarm (created 09-18 10:54Z, "nothing has landed in 14h"); **#105** has expired into the pile (superseded by #106). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-17 one: close #94–#105** — each was a one-day event, superseded by next day's alarm, never actionable; leave #106 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
300:- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) shows a new **stable v1.3.4** (2026-09-18) — maintenance only (data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change), so no new gap for Astrate — and still **no stable v1.4.0**: newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
308:- **github-issues triage run, 2026-09-17: still nothing proposable — the mule-alarm pile is now 12 straight days (#94–#105).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-16 run: **#105**, today's alarm (created 09-17 10:51Z, "nothing has landed in 14h"); **#104** has expired into the pile (superseded by #105). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-16 one: close #94–#104** — each was a one-day event, superseded by next day's alarm, never actionable; leave #105 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
310:- **github-issues triage run, 2026-09-16: still nothing proposable — the mule-alarm pile is now 11 straight days (#94–#104).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), plus the daily alarm. New since the 2026-09-15 run: **#104**, today's alarm (created 09-16 11:12Z, "nothing has landed in 14h"); **#103** has expired into the pile (superseded by #104). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-15 one: close #94–#103** — each was a one-day event, superseded by next day's alarm, never actionable; leave #104 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
312:- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-13 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
314:- **github-issues triage run, 2026-09-15: still nothing proposable — the mule-alarm pile is now 10 straight days (#94–#103).** Fourteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-13 run: **#103**, the alarm created 2026-09-14 11:27Z ("nothing has landed in 16h"); **#102** has expired into the pile (superseded by #103). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-13 one: close #94–#102** — each was a one-day event, superseded by next day's alarm, never actionable; leave #103 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
316:- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-12 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
318:- **github-issues triage run, 2026-09-13: still nothing proposable — the mule-alarm pile is now 9 straight days (#94–#102).** Thirteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-12 run: **#102**, today's alarm (created 11:09Z, "nothing has landed in 14h"); **#101** has expired into the pile (superseded by #102). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-12 one: close #94–#101** — each was a one-day event, superseded by next day's alarm, never actionable; leave #102 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
320:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
324:- **github-issues triage run, 2026-09-12: still nothing proposable — the mule-alarm pile is now 8 straight days (#94–#101).** Twelve open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-11 run: **#101**, today's alarm (created 10:58Z, "nothing has landed in 14h"); **#100** has expired into the pile (superseded by #101). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-11 one: close #94–#100** — each was a one-day event, superseded by next day's alarm, never actionable; leave #101 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
326:- **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pile is now 7 straight days and worth a look.** Twelve open issues. Still no machine-checkable fix candidates, so no task lines: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-10 run: **#100**, today's alarm (created 10:55Z, "nothing has landed in 15h"). #99–#94 are the same `mule-alarm` (one per day 2026-09-05→09-10, ~11:00Z each, zero comments) — telemetry, not code issues, so never proposable. But seven consecutive alarms is past the "low-activity window" wording earlier runs used: the mule has landed no commit since ~2026-09-04/05, and `main`'s queue copy shows the top line as `- [ ]` while the landed work below is all `[!] BLOCKED` (`wrote nothing` / `tests failed`) — though `mule/queue` is the authoritative copy. **Proposal: close #94–#99** — each was a one-day event, superseded by next day's alarm, never actionable; leave #100 (live today) to self-expire. If the idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
328:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
330:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
356:  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
362:  issues are the same: #92 keyAgreement (upstream decision, gated on stable v1.4.0), #93
365:  on upstream v1.4 final (currently rc.5). `APICompatVersion` stays 1.2.2 per #90's frozen
373:  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
381:- **Milestone recipe run, 2026-09-06: v3.0 (upstream 1.2.2 → 1.3/1.4) has no open
388:  experimental decision, escalated 2026-09-05, gated on a stable upstream v1.4.0), **#78 FDO**
391:  only when upstream ships **v1.4 final** — not yet (still rc). Per **#90's frozen decision**
392:  (2026-08-23), `APICompatVersion` stays 1.2.2 until the FULL 1.3 surface is complete *and*
394:  call when v1.4.0 goes stable: answer #92 (implement against the rc, re-park, or narrow #93
414:  awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the waiting-on
422:  stable v1.4.0?** (issue #92, `upstream-parity`/`upstream-experimental`). #51 closed
430:  fixed the stale ACL comment that cited them). But the spec ships only in `v1.4.0-rc.5` —
431:  `v1.3.3` is still the newest stable tag and Astrate targets 1.2.2 — so the document has
435:  re-park until v1.4.0 is a stable tag, or take only the narrow #93 fix. (Escalated again
445:- **COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable, 2026-07-14; v1.4.0 is still rc-only).** Astrate's doc and `APICompatVersion` still target upstream **1.2.2** (`internal/realm/service.go:588`); v1.3.0 (2026-05-06) introduced wire-surface changes Astrate does not yet emulate, so this is a decision — adopt v1.3.2 as the compatibility target (then update the doc + bump `APICompatVersion` together, per the bump rule) or keep 1.2.2 and add a "not yet emulated" note. Wire-relevant v1.3.0 deltas (release notes): **MQTT v1 capabilities** incl. `purge_properties_compression_format` (plaintext vs zlib purge — touches the `emptyCache`/`producer,properties` contract COMPATIBILITY.md deviation 1 documents); **empty introspection now allowed**; **device registration triggers** (pairing) and **device deletion started/completed triggers** (RM — the latter two already exist as Astrate deviation 9 emits both around the synchronous delete); **FDO authentication** (pairing, disabled by default); **realm-scoped health** — upstream v1.3 added `GET /pairing/v1/{realm}/health`, which Astrate already serves (`internal/pairing/http.go:78-81`, comment already says "upstream 1.3+"), so deviation 18's wording ("which upstream 404s") is now false against 1.3 and the note should be reworded either way. Proposed doc wording (for your approval, edit to taste): in §Infrastructure differences add a sentence — *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl.
447:plaintext `purge_properties_compression_format`, empty-introspection allowance, device registration/deletion triggers, experimental FDO pairing auth) are not yet emulated and are out of scope until the milestone that adopts v1.3.2 as the target."* — and reword deviation 18's realm-health note from "which upstream 404s" to "added by upstream v1.3 (Astrate serves it against a 1.2.2 target; kept, matching behavior)". Raw upstream changes: [v1.3.0](https://github.com/astarte-platform/astarte/releases/tag/v1.3.0), [v1.3.2](https://github.com/astarte-platform/astarte/releases/tag/v1.3.2).
467:  - `go.etcd.io/bbolt` v1.4.3 → v1.5.0 — bbolt's semver promises no API change between patch/minor, so additive-only; used in `internal/broker/sessionstore.go`; v1.5 adds a data-file size limit and panic-recovery hardening, nothing Astrate needs.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.5") | .body'
Astarte v1.3.5 release.

## Changes since v1.3.4

### Fixed

- [astarte_data_updater_plant] Increase RPC server reliability. A new corner case has been fixed which resulted in the RPC server not being available

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.4") | .body'
### Fixed

- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
  uncollected binaries over time by forcing more frequent full sweep garbage collections
  on them.
- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.

$ grep -n "astarte-upstream" /root/astrate-mule/.mule/log.md | tail -20; echo ---; grep -rn "astarte-upstream" /root/astrate-mule/.mule/reports/ 2>/dev/null | tail -20; echo ---; ls -R /root/astrate-mule/.mule/reports/
---
/root/astrate-mule/.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md:171:slug: recipe-astarte-upstream
/root/astrate-mule/.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md:175:--- /root/astrate-mule/.mule/reports/log/20260907T020359Z-recipe-astarte-upstream-proposed.md
/root/astrate-mule/.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md:176:slug: recipe-astarte-upstream
/root/astrate-mule/.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md:180:--- /root/astrate-mule/.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md
/root/astrate-mule/.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md:181:slug: recipe-astarte-upstream
/root/astrate-mule/.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md:185:--- /root/astrate-mule/.mule/reports/log/20260908T113422Z-recipe-astarte-upstream-proposed.md
/root/astrate-mule/.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md:186:slug: recipe-astarte-upstream
/root/astrate-mule/.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md:192:→ Read .mule/reports/log/20260908T113422Z-recipe-astarte-upstream-proposed.md
/root/astrate-mule/.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md:224:$ ls /root/astrate-mule/.mule/reports/log/ | grep -vE "astarte-upstream|docs-sync|hygiene|milestones|github-issues|code-review" | head; echo "===verdicts==="; grep -rh "^verdict:" /root/astrate-mule/.mule/reports/log/*.md | sort | uniq -c
/root/astrate-mule/.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md:245:← Write .mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md
/root/astrate-mule/.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md:254:Files: .mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md
/root/astrate-mule/.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:1:slug: recipe-astarte-upstream
/root/astrate-mule/.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:10:→ Read .mule/recipes/astarte-upstream.md
/root/astrate-mule/.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:40:astarte-upstream.md
/root/astrate-mule/.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:50:$ grep -rn "astarte-upstream\|upstream" .mule/*.md | head -60
/root/astrate-mule/.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:195:- **astarte-upstream recipe run, 2026-10-05 (evening): no upstream movement that touches anything Astrate implements, so nothing was queued — no `.mule/todo.md` line, no `gh issue create`, nothing commented, closed or edited on GitHub.** One thing *is* new since this recipe's last run (`2026-10-04T20:40:52+02:00` ran; the `2026-10-05T04:29:10+02:00` run timed out): stable **v1.3.5**, published `2026-10-05T14:11:44Z`. The milestone recipe already recorded it as inert today; I re-measured it rather than trusting that, because "a release note is a claim, not a fact" is this recipe's own rule. `gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5` is **4 commits, 35 files** — `8d47574f` + `4c4222cf` `chore: update horde`, `c4c63ebe` `chore(rm): bump xandra to 0.19.4`, `8f0f5ef1` `chore: prepare v1.3.5` — and **every app/lib file in it is a `mix.exs`/`mix.lock` version bump**. The rest is `CHANGELOG.md` (+6: one `### Fixed` entry, "[astarte_data_updater_plant] Increase RPC server reliability"), `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `README.md`, `astarte-dashboard/package.json`. Not one source file changed, so no route, MQTT topic, BSON encoding, control message or interface-schema field moved: Astrate implements no `data_updater_plant` and no Erlang RPC server. Wire-inert, same shape as v1.3.2/v1.3.3/v1.3.4. **Master since `v1.4.0-rc.6` is still the single commit `731bcddc`** ("chore: forward port release-1.4"; `gh api .../commits?since=2026-10-05T09:34:00Z` returns it and nothing newer), the migration-relocation commit the milestone run measured today as inert — so both wire-bearing post-rc.6 commits (`b6d46ad4`, `958bb9fe`) still have their owners in the queue and nothing new joins them.
/root/astrate-mule/.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:219:  - **`.mule/waiting-on.md` and the recipe step that reads it** (`f1d0069`→`2a033d4`; `origin/main:.mule/recipes/astarte-upstream.md:17` is the line "*First, before anything else, check the parked decisions*"). The file is in `origin/main` and **absent from `HEAD`**; `rg waiting-on tools/mule.sh .mule/recipes/ .agents/` returns nothing on this host. So #92's own comment — "the weekly `mule-upstream-watch` job now reads that file as its first step every run" — is **not true of this host**. `mule-upstream-watch.timer` does exist: LAST **Mon 2026-09-28 04:02:35 CEST**, NEXT **2026-10-05 04:08:35 CEST**. It has therefore already run once without the parked-decisions step, and it runs again tomorrow morning without it. #92's re-parking has survived a month on hand-checks by the milestone and triage runs, not on the mechanism the issue was told about.
/root/astrate-mule/.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:228:$ grep -n "astarte-upstream\|upstream" .mule/log.md | tail -20; echo "=== log head ==="; head -40 .mule/log.md
/root/astrate-mule/.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:271:$ grep -n "compat-note\|probe-encrypted" .mule/todo.md | head; echo "=== check for-giulio for 10-07 upstream run ==="; grep -n "astarte-upstream recipe run, 2026-10" .mule/for-giulio.md | head
/root/astrate-mule/.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:281:38:- **astarte-upstream recipe run, 2026-10-05 (evening): no upstream movement that touches anything Astrate implements, so nothing was queued — no `.mule/todo.md` line, no `gh issue create`, nothing commented, closed or edited on GitHub.** One thing *is* new since this recipe's last run (`2026-10-04T20:40:52+02:00` ran; the `2026-10-05T04:29:10+02:00` run timed out): stable **v1.3.5**, published `2026-10-05T14:11:44Z`. The milestone recipe already recorded it as inert today; I re-measured it rather than trusting that, because "a release note is a claim, not a fact" is this recipe's own rule. `gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5` is **4 commits, 35 files** — `8d47574f` + `4c4222cf` `chore: update horde`, `c4c63ebe` `chore(rm): bump xandra to 0.19.4`, `8f0f5ef1` `chore: prepare v1.3.5` — and **every app/lib file in it is a `mix.exs`/`mix.lock` version bump**. The rest is `CHANGELOG.md` (+6: one `### Fixed` entry, "[astarte_data_updater_plant] Increase RPC server reliability"), `.github/workflows/astarte-build-workflow.yaml`, `docker-compose.yml`, `README.md`, `astarte-dashboard/package.json`. Not one source file changed, so no route, MQTT topic, BSON encoding, control message or interface-schema field moved: Astrate implements no `data_updater_plant` and no Erlang RPC server. Wire-inert, same shape as v1.3.2/v1.3.3/v1.3.4. **Master since `v1.4.0-rc.6` is still the single commit `731bcddc`** ("chore: forward port release-1.4"; `gh api .../commits?since=2026-10-05T09:34:00Z` returns it and nothing newer), the migration-relocation commit the milestone run measured today as inert — so both wire-bearing post-rc.6 commits (`b6d46ad4`, `958bb9fe`) still have their owners in the queue and nothing new joins them.
---
/root/astrate-mule/.mule/reports/:
hygiene.md
log

/root/astrate-mule/.mule/reports/log:
20260904T190911Z-purge-properties-compression-capabilityauto-done.md
20260904T195018Z-broker-acl-coldstart-introspection-miss-done.md
20260904T195529Z-issue-93-done.md
20260904T195814Z-broker-disconnect-device-zombie-session-blocked.md
20260904T203003Z-broker-offline-acl-tests-blocked.md
20260904T210532Z-broker-onconnect-doc-comment-blocked.md
20260905T110211Z-empty-introspection-verification-blocked.md
20260905T111428Z-probe-trigger-install-notification-delay-blocked.md
20260905T112539Z-compat-note-v132-done.md
20260905T165512Z-recipe-docs-sync-proposed.md
20260905T171755Z-docs-sync-pairing-health-path-done.md
20260905T172234Z-docs-sync-pairing-register-404-done.md
20260905T173657Z-recipe-hygiene-proposed.md
20260905T174528Z-swagger-httptest-coverage-done.md
20260905T181053Z-recipe-milestones-proposed.md
20260905T181219Z-recipe-github-issues-proposed.md
20260905T184625Z-recipe-astarte-upstream-proposed.md
20260905T184903Z-probe-property-resend-encoding-blocked.md
20260905T190855Z-compat-note-v133-blocked.md
20260905T192010Z-recipe-code-review-proposed.md
20260905T192911Z-flow-validate-source-sink-blocked.md
20260905T201529Z-flow-validate-dead-source-sink-recompute-done.md
20260906T112006Z-recipe-docs-sync-proposed.md
20260906T113143Z-docs-sync-rm-datastream-retention-endpoint-done.md
20260906T171045Z-docs-sync-rm-interfaces-detailed-param-done.md
20260906T174031Z-recipe-hygiene-proposed.md
20260906T175748Z-housekeeping-tests-done.md
20260906T180442Z-recipe-milestones-proposed.md
20260906T181352Z-recipe-github-issues-proposed.md
20260906T183102Z-recipe-astarte-upstream-proposed.md
20260906T183952Z-probe-props-resend-error-triggers-blocked.md
20260906T185735Z-recipe-code-review-proposed.md
20260906T195827Z-server-data-trigger-bus-done.md
20260906T200431Z-recipe-docs-sync-proposed.md
20260906T201002Z-docs-sync-appengine-by-alias-endpoints-done.md
20260906T202744Z-docs-sync-appengine-group-endpoints-done.md
20260906T204333Z-docs-sync-appengine-get-group-device-done.md
20260907T020359Z-recipe-astarte-upstream-proposed.md
20260907T112447Z-docs-sync-appengine-query-params-status-done.md
20260907T112654Z-docs-sync-appengine-group-patch-status-done.md
20260907T170024Z-docs-sync-appengine-data-422-interface-level-done.md
20260907T170452Z-recipe-hygiene-proposed.md
20260907T171928Z-recipe-milestones-proposed.md
20260907T174812Z-recipe-github-issues-proposed.md
20260907T175524Z-recipe-astarte-upstream-proposed.md
20260907T181549Z-recipe-code-review-proposed.md
20260907T184850Z-appengine-snapshot-ignores-query-params-blocked.md
20260907T190319Z-appengine-group-token-roundtrip-test-blocked.md
20260907T191157Z-recipe-docs-sync-proposed.md
20260907T201232Z-docs-sync-appengine-data-output-params-done.md
20260907T201806Z-recipe-hygiene-proposed.md
20260907T204751Z-recipe-milestones-proposed.md
20260908T112032Z-recipe-github-issues-proposed.md
20260908T113422Z-recipe-astarte-upstream-proposed.md
20260908T114256Z-probe-required-mapping-flag-done.md
20260908T164748Z-compat-note-v14-rc-blocked.md
20260908T165609Z-detailed-listing-required-encrypted-flags-done.md
20260908T170007Z-recipe-code-review-proposed.md
20260908T173033Z-store-alias-values-taken-test-done.md
20260908T180824Z-store-latest-individual-test-done.md
20260908T182108Z-recipe-docs-sync-proposed.md
20260908T183808Z-docs-sync-appengine-device-status-schema-done.md
20260908T184535Z-docs-sync-appengine-data-set-422-blocked.md
20260908T185820Z-docs-sync-appengine-downsample-min-done.md
20260908T192323Z-recipe-hygiene-timeout.md
20260908T193026Z-recipe-milestones-proposed.md
20260908T193127Z-recipe-milestones-proposed.md
20260908T193321Z-recipe-github-issues-proposed.md
20260908T193614Z-recipe-astarte-upstream-noop.md
20260908T193814Z-recipe-astarte-upstream-proposed.md
20260909T110758Z-recipe-code-review-proposed.md
20260909T111504Z-pairing-bearer-secret-test-done.md
20260909T165402Z-pairing-remoteip-fallback-test-done.md
20260909T165735Z-recipe-docs-sync-proposed.md
20260909T170637Z-docs-sync-hk-patch-endpoint-done.md
20260909T171614Z-docs-sync-hk-retention-field-done.md
20260909T174959Z-recipe-hygiene-timeout.md
20260909T175718Z-recipe-milestones-proposed.md
20260909T181348Z-recipe-github-issues-proposed.md
20260909T184315Z-recipe-astarte-upstream-proposed.md
20260909T191420Z-recipe-code-review-proposed.md
20260909T195057Z-realm-pure-helper-tests-done.md
20260909T200459Z-recipe-docs-sync-proposed.md
20260909T201715Z-docs-sync-native-compat-health-503-done.md
20260909T202852Z-docs-sync-native-compat-version-endpoints-done.md
20260909T230006Z-recipe-hygiene-timeout.md
20260909T230411Z-recipe-milestones-proposed.md
20260910T110956Z-recipe-github-issues-proposed.md
20260910T112515Z-recipe-astarte-upstream-proposed.md
20260910T112545Z-recipe-astarte-upstream-proposed.md
20260910T170149Z-recipe-code-review-proposed.md
20260910T171629Z-flowapi-autorestart-terminal-failure-done.md
20260910T172606Z-flowapi-validationdetail-test-done.md
20260910T182250Z-flowapi-autorestart-default-test-done.md
20260910T184706Z-recipe-docs-sync-proposed.md
20260910T190522Z-docs-sync-pairing-status-enum-done.md
20260910T191129Z-docs-sync-pairing-version-endpoint-done.md
20260910T192802Z-recipe-hygiene-proposed.md
20260910T195304Z-recipe-milestones-proposed.md
20260910T200120Z-recipe-github-issues-proposed.md
20260910T202936Z-recipe-astarte-upstream-noop.md
20260910T202949Z-recipe-astarte-upstream-proposed.md
20260911T105829Z-recipe-code-review-proposed.md
20260911T110803Z-recipe-docs-sync-proposed.md
20260911T112356Z-docs-sync-rm-delete-interface-status-done.md
20260911T113809Z-docs-sync-rm-mapping-required-encrypted-done.md
20260911T170716Z-docs-sync-rm-put-auth-422-done.md
20260911T182955Z-docs-sync-rm-version-example-blocked.md
20260911T183945Z-recipe-hygiene-proposed.md
20260911T185303Z-recipe-milestones-proposed.md
20260911T190417Z-recipe-github-issues-proposed.md
20260911T191541Z-recipe-astarte-upstream-proposed.md
20260911T194026Z-recipe-code-review-proposed.md
20260911T194927Z-interfaceschema-compat-attrs-flip-test-done.md
20260911T195637Z-interfaceschema-object-attrs-uniformity-done.md
20260911T200416Z-interfaceschema-malformed-placeholder-fixtures-done.md
20260911T201604Z-interfaceschema-enum-roundtrip-coverage-done.md
20260911T202002Z-interfaceschema-trie-interior-literal-param-done.md
20260912T110449Z-recipe-docs-sync-proposed.md
20260912T112234Z-docs-sync-ae-write-405-done.md
20260912T113032Z-docs-sync-ae-write-value-422-done.md
20260912T164501Z-docs-sync-ae-post-groups-409-done.md
20260912T170850Z-recipe-hygiene-timeout.md
20260912T171256Z-recipe-milestones-proposed.md
20260912T172901Z-recipe-github-issues-proposed.md
20260912T180735Z-recipe-astarte-upstream-proposed.md
20260912T182454Z-recipe-code-review-proposed.md
20260912T184128Z-payload-longinteger-fraction-quantize-done.md
20260912T190803Z-payload-json-malformed-t-not-tolerated-done.md
20260912T192246Z-recipe-docs-sync-proposed.md
20260912T192341Z-recipe-docs-sync-proposed.md
20260912T203456Z-docs-sync-ae-post-group-devices-409-done.md
20260912T204618Z-docs-sync-ae-patch-merge-patch-media-type-done.md
20260913T111551Z-recipe-hygiene-proposed.md
20260913T112223Z-store-todo-lttb-stale-done.md
20260913T165246Z-recipe-milestones-proposed.md
20260913T170550Z-recipe-github-issues-proposed.md
20260913T171327Z-recipe-astarte-upstream-noop.md
20260913T171351Z-recipe-astarte-upstream-proposed.md
20260913T172714Z-recipe-code-review-proposed.md
20260913T175723Z-config-fail-loud-engine-shards-done.md
20260913T180348Z-config-fail-loud-dev-mode-done.md
20260913T180935Z-config-key-resolvers-test-done.md
20260913T182321Z-recipe-docs-sync-proposed.md
20260913T185042Z-docs-sync-hk-delete-gating-responses-done.md
20260913T191910Z-docs-sync-hk-validation-example-done.md
20260913T202313Z-recipe-hygiene-timeout.md
20260913T203209Z-recipe-milestones-proposed.md
20260913T203800Z-recipe-github-issues-proposed.md
20260913T203830Z-recipe-github-issues-proposed.md
20260913T204415Z-recipe-astarte-upstream-proposed.md
20260914T020552Z-recipe-astarte-upstream-noop.md
20260914T020608Z-recipe-astarte-upstream-proposed.md
20260914T113508Z-recipe-code-review-proposed.md
20260914T114442Z-auth-iat-not-required-test-done.md
20260914T165413Z-auth-empty-claims-403-test-done.md
20260914T171558Z-auth-cache-default-size-test-done.md
20260914T174510Z-recipe-docs-sync-proposed.md
20260914T181827Z-recipe-hygiene-timeout.md
20260914T234550Z-recipe-milestones-proposed.md
20260915T105740Z-recipe-github-issues-proposed.md
20260915T111354Z-recipe-astarte-upstream-proposed.md
20260915T113026Z-recipe-code-review-proposed.md
20260915T170704Z-hk-zero-retention-asymmetry-done.md
20260915T171349Z-hk-zero-reglimit-create-asymmetry-blocked.md
20260915T173739Z-hk-zero-retention-asymmetry-done.md
20260915T180720Z-recipe-docs-sync-proposed.md
20260915T181748Z-docs-sync-native-socket-query-token-auth-done.md
20260915T184105Z-recipe-hygiene-timeout.md
20260915T184406Z-recipe-milestones-proposed.md
20260915T185804Z-recipe-github-issues-proposed.md
20260915T193626Z-recipe-astarte-upstream-proposed.md
20260915T201017Z-recipe-code-review-proposed.md
20260915T202012Z-astarteapi-metadata-envelope-test-done.md
20260916T111701Z-astarteapi-multikey-sorted-golden-done.md
20260916T115527Z-recipe-docs-sync-proposed.md
20260916T165431Z-docs-sync-pairing-deviceendpoints-dead-404-403-done.md
20260916T171527Z-docs-sync-pairing-version-404-unreachable-done.md
20260916T173558Z-recipe-hygiene-proposed.md
20260916T175801Z-recipe-milestones-proposed.md
20260916T183922Z-recipe-github-issues-proposed.md
20260916T184159Z-recipe-astarte-upstream-noop.md
20260916T184218Z-recipe-astarte-upstream-proposed.md
20260916T184848Z-recipe-code-review-proposed.md
20260916T193131Z-observability-readiness-wedge-test-done.md
20260916T195125Z-observability-dbpool-gauge-coverage-done.md
20260916T202828Z-observability-health-content-type-test-done.md
20260916T204810Z-recipe-docs-sync-proposed.md
20260917T111017Z-docs-sync-pairing-version-value-done.md
20260917T114541Z-docs-sync-native-version-value-done.md
20260917T182237Z-recipe-hygiene-timeout.md
20260917T184319Z-recipe-milestones-timeout.md
20260917T190323Z-recipe-github-issues-proposed.md
20260917T190650Z-recipe-astarte-upstream-noop.md
20260917T190703Z-recipe-astarte-upstream-proposed.md
20260917T191223Z-recipe-code-review-proposed.md
20260917T192217Z-httpx-cors-vary-origin-passthrough-done.md
20260917T194954Z-recipe-docs-sync-proposed.md
20260917T200755Z-docs-sync-rm-policies-delete-422-done.md
20260917T202137Z-docs-sync-rm-triggers-422-nested-envelope-done.md
20260917T210114Z-recipe-hygiene-timeout.md
20260918T110039Z-recipe-milestones-proposed.md
20260918T113442Z-recipe-github-issues-proposed.md
20260918T113930Z-recipe-astarte-upstream-proposed.md
20260918T114100Z-recipe-astarte-upstream-proposed.md
20260918T165351Z-compat-note-v134-done.md
20260918T172636Z-recipe-code-review-proposed.md
20260918T180119Z-recipe-docs-sync-proposed.md
20260918T184115Z-recipe-hygiene-timeout.md
20260918T193641Z-recipe-milestones-proposed.md
20260918T194216Z-recipe-github-issues-proposed.md
20260918T194746Z-recipe-astarte-upstream-proposed.md
20260918T200041Z-recipe-code-review-proposed.md
20260918T202346Z-channels-group-watch-membership-scope-done.md
20260918T204350Z-channels-rejoin-joinref-tagging-test-done.md
20260919T105228Z-channels-rejoin-authz-mismatch-blocked.md
20260919T113513Z-recipe-docs-sync-proposed.md
20260919T170808Z-recipe-hygiene-timeout.md
20260919T211455Z-recipe-milestones-proposed.md
20260920T111057Z-recipe-github-issues-proposed.md
20260920T113611Z-recipe-astarte-upstream-proposed.md
20260920T173118Z-recipe-code-review-proposed.md
20260920T181128Z-errorname-missing-required-test-done.md
20260920T184908Z-triggers-custom-action-policy-nodecide-done.md
20260920T185539Z-compat-note-custom-action-policy-boundary-done.md
20260920T185956Z-recipe-docs-sync-proposed.md
20260920T191642Z-docs-sync-native-metrics-example-fake-series-done.md
20260920T192702Z-docs-sync-native-socket-missing-403-500-done.md
20260920T195443Z-recipe-hygiene-timeout.md
20260920T203030Z-recipe-milestones-proposed.md
20260920T204223Z-recipe-github-issues-proposed.md
20260921T020624Z-recipe-astarte-upstream-proposed.md
20260921T105709Z-recipe-astarte-upstream-proposed.md
20260921T113318Z-recipe-code-review-proposed.md
20260921T173503Z-store-pipelines-empty-name-zero-blocks-test-done.md
20260921T174646Z-recipe-docs-sync-proposed.md
20260921T182029Z-recipe-hygiene-timeout.md
20260921T183804Z-recipe-milestones-proposed.md
20260921T185351Z-recipe-github-issues-proposed.md
20260921T190635Z-recipe-astarte-upstream-proposed.md
20260921T192654Z-recipe-code-review-proposed.md
20260921T194020Z-broker-acl-coldstart-fallback-flood-done.md
20260921T201044Z-broker-offlineacl-entry-eviction-done.md
20260921T203450Z-recipe-docs-sync-proposed.md
20260922T105805Z-docs-sync-ae-read-query-params-blocked.md
20260922T121234Z-recipe-hygiene-timeout.md
20260922T123011Z-recipe-milestones-proposed.md
20260922T164457Z-recipe-github-issues-proposed.md
20260922T170456Z-recipe-astarte-upstream-proposed.md
20260922T170955Z-recipe-code-review-proposed.md
20260922T173055Z-recipe-docs-sync-proposed.md
20260922T181149Z-recipe-hygiene-proposed.md
20260922T190954Z-recipe-milestones-proposed.md
20260922T193852Z-recipe-github-issues-proposed.md
20260922T194051Z-recipe-astarte-upstream-proposed.md
20260922T195256Z-recipe-code-review-proposed.md
20260922T200147Z-recipe-docs-sync-proposed.md
20260922T200748Z-recipe-hygiene-proposed.md
20260922T202552Z-recipe-milestones-proposed.md
20260922T204751Z-recipe-github-issues-proposed.md
20260923T105200Z-recipe-astarte-upstream-proposed.md
20260923T112257Z-recipe-code-review-proposed.md
20260923T112816Z-recipe-docs-sync-proposed.md
20260923T165837Z-recipe-hygiene-proposed.md
20260923T170727Z-examples-echo-container-contract-test-done.md
20260923T174341Z-recipe-milestones-proposed.md
20260923T174805Z-recipe-github-issues-proposed.md
20260923T181743Z-recipe-astarte-upstream-proposed.md
20260923T184009Z-recipe-code-review-proposed.md
20260923T185047Z-flow-mqtt-source-reconnect-recovery-blocked.md
20260923T190710Z-flow-msg-json-integer-precision-done.md
20260923T192627Z-flow-sort-bounded-buffer-done.md
20260923T193304Z-flow-randomsource-span-overflow-blocked.md
20260923T195241Z-flow-filter-key-contains-test-done.md
20260923T201015Z-recipe-docs-sync-proposed.md
20260923T201746Z-docs-sync-rm-put-interface-409-done.md
20260924T110303Z-docs-sync-rm-interface-422-shapes-done.md
20260924T111049Z-docs-sync-rm-validation-example-prefix-done.md
20260924T165605Z-docs-sync-rm-auth-403-done.md
20260924T174100Z-recipe-hygiene-proposed.md
20260924T175556Z-recipe-milestones-proposed.md
20260924T182956Z-recipe-github-issues-proposed.md
20260924T183158Z-recipe-astarte-upstream-proposed.md
20260924T183755Z-recipe-code-review-proposed.md
20260924T184819Z-swagger-sub-failfast-done.md
20260924T185047Z-recipe-docs-sync-proposed.md
20260924T190543Z-recipe-hygiene-proposed.md
20260924T192153Z-recipe-milestones-proposed.md
20260924T193516Z-recipe-github-issues-proposed.md
20260924T194457Z-recipe-astarte-upstream-proposed.md
20260924T201012Z-recipe-code-review-timeout.md
20260924T202134Z-recipe-docs-sync-proposed.md
20260925T112307Z-recipe-hygiene-timeout.md
20260925T115624Z-recipe-milestones-timeout.md
20260925T165408Z-recipe-github-issues-proposed.md
20260925T170557Z-recipe-astarte-upstream-proposed.md
20260925T174727Z-recipe-code-review-proposed.md
20260925T175637Z-drain-per-stage-budget-done.md
20260925T180502Z-cmd-devcert-fields-test-done.md
20260925T184021Z-cmd-healthcheck-contract-test-blocked.md
20260925T185049Z-cmd-loadsealer-masterkeyfile-test-done.md
20260925T192553Z-recipe-docs-sync-proposed.md
20260925T193504Z-docs-sync-hk-async-operation-param-done.md
20260925T195121Z-docs-sync-hk-retention-zero-is-unset-done.md
20260925T200001Z-docs-sync-rm-async-operation-param-done.md
20260925T201910Z-docs-sync-rm-delete-device-async-operation-param-done.md
20260925T204521Z-docs-sync-rm-delete-device-async-operation-param-blocked.md
20260926T112747Z-recipe-hygiene-timeout.md
20260926T113138Z-recipe-milestones-proposed.md
20260926T164829Z-recipe-github-issues-proposed.md
20260926T173418Z-recipe-astarte-upstream-timeout.md
20260926T174259Z-recipe-code-review-proposed.md
20260926T180257Z-forward-static-header-validation-done.md
20260926T181916Z-forward-static-headers-override-test-done.md
20260926T182449Z-forward-status-error-body-done.md
20260926T190026Z-forward-envelope-bytes-test-done.md
20260926T192057Z-webhook-static-headers-override-test-done.md
20260926T194649Z-recipe-docs-sync-proposed.md
20260926T203415Z-docs-sync-ae-forbidden-403-done.md
20260927T111051Z-docs-sync-ae-patch-by-alias-409-done.md
20260927T111624Z-docs-sync-ae-add-group-device-422-done.md
20260927T113516Z-docs-sync-ae-delete-data-400-done.md
20260927T164857Z-docs-sync-ae-list-devices-422-transient.md
20260927T171520Z-docs-sync-ae-list-devices-422-transient.md
20260927T182324Z-docs-sync-ae-list-devices-422-blocked.md
20260927T184918Z-recipe-hygiene-timeout.md
20260927T191018Z-recipe-milestones-timeout.md
20260927T195417Z-recipe-github-issues-timeout.md
20260927T203020Z-recipe-astarte-upstream-timeout.md
20260927T205417Z-recipe-code-review-timeout.md
20260928T020420Z-recipe-astarte-upstream-proposed.md
20260928T105915Z-recipe-docs-sync-proposed.md
20260928T111752Z-docs-sync-pairing-error-example-capitalisation-done.md
20260928T113726Z-docs-sync-pairing-info-version-example-done.md
20260928T170059Z-docs-sync-pairing-deviceid-base64url-done.md
20260928T172627Z-docs-sync-pairing-initial-payload-format-enum-done.md
20260928T180120Z-docs-sync-pairing-unregister-description-done.md
20260928T200706Z-recipe-hygiene-timeout.md
20260928T200921Z-recipe-milestones-proposed.md
20260928T202852Z-recipe-github-issues-proposed.md
20260928T203150Z-recipe-astarte-upstream-proposed.md
20260928T204149Z-recipe-code-review-proposed.md
20260928T204651Z-recipe-docs-sync-proposed.md
20260929T105054Z-recipe-hygiene-proposed.md
20260929T110049Z-recipe-milestones-proposed.md
20260929T110251Z-recipe-github-issues-proposed.md
20260929T165352Z-recipe-astarte-upstream-proposed.md
20260929T170753Z-recipe-code-review-proposed.md
20260929T170954Z-recipe-docs-sync-proposed.md
20260929T173251Z-recipe-hygiene-proposed.md
20260929T175252Z-recipe-milestones-proposed.md
20260929T175958Z-recipe-github-issues-proposed.md
20260929T180551Z-recipe-astarte-upstream-proposed.md
20260929T181453Z-recipe-code-review-proposed.md
20260929T183752Z-recipe-docs-sync-proposed.md
20260929T185857Z-recipe-hygiene-proposed.md
20260929T192156Z-recipe-milestones-proposed.md
20260929T193152Z-recipe-github-issues-proposed.md
20260929T194353Z-recipe-astarte-upstream-proposed.md
20260930T113252Z-recipe-code-review-proposed.md
20260930T113453Z-recipe-docs-sync-proposed.md
20260930T164953Z-recipe-hygiene-proposed.md
20260930T165653Z-recipe-milestones-proposed.md
20260930T170751Z-recipe-github-issues-proposed.md
20260930T171256Z-recipe-astarte-upstream-proposed.md
20260930T174155Z-recipe-code-review-proposed.md
20260930T175957Z-recipe-docs-sync-proposed.md
20260930T183258Z-recipe-hygiene-proposed.md
20260930T184654Z-recipe-milestones-proposed.md
20260930T190453Z-recipe-github-issues-proposed.md
20260930T191053Z-recipe-astarte-upstream-proposed.md
20260930T191757Z-recipe-code-review-proposed.md
20260930T192852Z-recipe-docs-sync-proposed.md
20260930T193453Z-recipe-hygiene-proposed.md
20260930T194054Z-recipe-milestones-proposed.md
20261001T105953Z-recipe-github-issues-proposed.md
20261001T111852Z-recipe-astarte-upstream-proposed.md
20261001T113558Z-recipe-code-review-proposed.md
20261001T164758Z-recipe-docs-sync-proposed.md
20261001T174651Z-recipe-hygiene-proposed.md
20261001T175146Z-recipe-milestones-proposed.md
20261001T181452Z-recipe-github-issues-proposed.md
20261001T184353Z-recipe-astarte-upstream-proposed.md
20261001T185552Z-recipe-code-review-proposed.md
20261001T185956Z-recipe-docs-sync-proposed.md
20261001T192350Z-recipe-hygiene-proposed.md
20261001T194155Z-recipe-milestones-proposed.md
20261001T201010Z-recipe-github-issues-proposed.md
20261001T203619Z-recipe-astarte-upstream-proposed.md
20261001T204604Z-fdo-rc6-scope-delta-for-giulio-done.md
20261002T113048Z-compat-note-v14-rc6-done.md
20261002T114427Z-recipe-code-review-proposed.md
20261002T165525Z-container-stop-deadline-done.md
20261002T172128Z-container-timeout-bounds-blocked.md
20261002T173949Z-container-parseconfig-rules-test-done.md
20261002T175003Z-container-response-cap-test-blocked.md
20261002T181712Z-recipe-docs-sync-proposed.md
20261002T183434Z-docs-sync-rm-error-example-capitalisation-done.md
20261002T192357Z-docs-sync-rm-legacy-alias-fields-blocked.md
20261002T195225Z-docs-sync-rm-validationerror-example-done.md
20261002T201307Z-docs-sync-rm-deviceid-param-done.md
20261002T203814Z-docs-sync-rm-update-interface-body-blocked.md
20261003T110748Z-docs-sync-ae-error-example-capitalisation-done.md
20261003T113021Z-docs-sync-hk-error-example-capitalisation-blocked.md
20261003T165209Z-docs-sync-native-error-example-capitalisation-done.md
20261003T170519Z-docs-sync-hk-error-detail-examples-split-done.md
20261003T173401Z-recipe-hygiene-proposed.md
20261003T175916Z-recipe-milestones-proposed.md
20261003T180922Z-recipe-github-issues-proposed.md
20261003T183445Z-recipe-astarte-upstream-proposed.md
20261003T184531Z-device-empty-cache-received-trigger-blocked.md
20261003T190400Z-appengine-unexpected-object-key-done.md
20261003T193009Z-appengine-payload-reason-status-map-blocked.md
20261003T195023Z-appengine-missing-required-422-done.md
20261003T202817Z-recipe-code-review-proposed.md
20261004T112505Z-recipe-docs-sync-proposed.md
20261004T114011Z-docs-sync-hk-realm-name-pattern-done.md
20261004T165229Z-docs-sync-hk-errordetail-schema-example-done.md
20261004T171709Z-docs-sync-hk-patch-422-field-error-examples-blocked.md
20261004T172242Z-docs-sync-hk-wrong-type-field-400-blocked.md
20261004T173828Z-docs-sync-hk-realm-name-response-schemas-done.md
20261004T181112Z-recipe-hygiene-timeout.md
20261004T181934Z-recipe-milestones-proposed.md
20261004T183240Z-recipe-github-issues-proposed.md
20261004T184052Z-recipe-astarte-upstream-proposed.md
20261004T190849Z-recipe-code-review-proposed.md
20261004T193935Z-store-register-inhibit-preserve-blocked.md
20261004T200628Z-store-aliasvalues-self-exclusion-test-done.md
20261004T203510Z-store-latestindividual-empty-errnotfound-blocked.md
20261005T022910Z-recipe-astarte-upstream-timeout.md
20261005T105949Z-store-validatepipelinegraph-error-branches-done.md
20261005T111650Z-recipe-docs-sync-proposed.md
20261005T113206Z-docs-native-ae-realm-version-endpoint-done.md
20261005T171138Z-docs-native-hk-version-description-done.md
20261005T175822Z-docs-native-metrics-content-negotiation-done.md
20261005T180602Z-docs-native-phoenix-newevent-payload-done.md
20261005T181428Z-docs-native-realm-name-pattern-done.md
20261005T192311Z-recipe-hygiene-timeout.md
20261005T193452Z-recipe-milestones-proposed.md
20261005T195440Z-recipe-github-issues-proposed.md
20261005T201617Z-recipe-astarte-upstream-proposed.md
20261005T201914Z-recipe-code-review-proposed.md
20261005T210409Z-recipe-docs-sync-timeout.md
20261006T140535Z-recipe-hygiene-timeout.md
20261006T141129Z-recipe-milestones-proposed.md
20261006T165612Z-recipe-github-issues-proposed.md
20261006T174353Z-recipe-astarte-upstream-proposed.md
20261006T174931Z-probe-encrypted-mapping-enforcement-blocked.md
20261006T175547Z-compat-note-v135-done.md
20261006T182510Z-recipe-code-review-timeout.md
20261006T184046Z-recipe-docs-sync-proposed.md
20261006T185740Z-docs-native-socket-handshake-errors-blocked.md
20261006T201112Z-docs-native-socket-sse-exact-accept-done.md
20261006T202927Z-docs-native-socket-security-scheme-done.md
20261006T204912Z-docs-native-socket-event-schema-blocked.md
20261007T112533Z-recipe-hygiene-timeout.md
20261007T113844Z-recipe-milestones-proposed.md
20261007T165810Z-recipe-github-issues-proposed.md
20261007T171801Z-recipe-astarte-upstream-proposed.md
20261007T174416Z-recipe-code-review-proposed.md
20261007T180202Z-pairing-burn-bcrypt-cost-blocked.md
20261007T180823Z-pairing-onregistered-clock-done.md
20261007T184717Z-pairing-unregister-bad-id-test-done.md
20261007T191848Z-recipe-docs-sync-proposed.md
20261007T223040Z-recipe-hygiene-timeout.md
20261007T224444Z-recipe-milestones-proposed.md
20261008T111054Z-recipe-github-issues-proposed.md
20261008T113945Z-recipe-astarte-upstream-proposed.md
20261008T172851Z-recipe-code-review-proposed.md
20261008T173721Z-testutil-wait-since-cursor-done.md
20261008T174752Z-testutil-control-frame-test-blocked.md
20261008T175838Z-engine-introspection-producer-roundtrip-done.md
20261008T182643Z-recipe-docs-sync-proposed.md
20261008T183153Z-docs-sync-pairing-register-404-unreachable-done.md
20261008T191310Z-docs-sync-pairing-body-cap-400-blocked.md
20261008T192820Z-docs-sync-pairing-realm-name-pattern-done.md
20261008T193604Z-docs-sync-pairing-deviceid-path-param-done.md
20261008T200623Z-recipe-hygiene-timeout.md
20261008T202151Z-recipe-milestones-proposed.md
20261008T203843Z-recipe-github-issues-proposed.md

$ gh api repos/astarte-platform/astarte/compare/v1.3.4...v1.3.5 --jq '{commits: (.commits|length), files: (.files|length), messages: [.commits[].commit.message]}' 2>&1 | head -40
{"commits":4,"files":35,"messages":["chore: update horde (#2224)\n\nSigned-off-by: Francesco Noacco \u003cfrancesco.noacco@secomind.com\u003e","chore: update horde\n\nSigned-off-by: Francesco Noacco \u003cfrancesco.noacco@secomind.com\u003e","chore(rm): bump xandra to 0.19.4 (#2246)\n\nSigned-off-by: Riccardo Nalgi \u003criccardo.nalgi@secomind.com\u003e","chore: prepare v1.3.5 (#2245)\n\nSigned-off-by: Francesco Noacco \u003cfrancesco.noacco@secomind.com\u003e"]}

→ Read .mule/reports/log/20261008T113945Z-recipe-astarte-upstream-proposed.md
→ Read .mule/reports/log/20261006T175547Z-compat-note-v135-done.md
→ Read .mule/reports/log/20261008T113945Z-recipe-astarte-upstream-proposed.md [offset=155]

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:6][] | "\(.tag_name)\t\(.published_at)\t\(.prerelease)"'
v1.3.5	2026-10-05T14:11:44Z	false
v1.4.0-rc.6	2026-09-30T07:28:27Z	true
v1.3.4	2026-09-18T09:42:57Z	false
v1.4.0-rc.5	2026-08-20T12:28:04Z	true
v1.3.3	2026-08-07T14:33:58Z	false
v1.4.0-rc.4	2026-08-06T13:40:20Z	true

$ gh api "repos/astarte-platform/astarte/commits?since=2026-10-08T11:39:45Z&per_page=30" --jq '.[] | "\(.sha[0:8])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"'; echo "rc=$?"
rc=0

No upstream movement since the last run: releases are unchanged (`v1.3.5` stable 2026-10-05, `v1.4.0-rc.6` prerelease) and `gh api .../commits?since=2026-10-08T11:39:45Z` returns nothing. The v1.3.5 compat note is already `[x]` and its wording delta is already in `for-giulio.md`, so re-proposing would duplicate.

```
Done: no upstream movement since v1.3.5 — nothing appended to .mule/todo.md, no evidence file (recipe names none)
Files: none
Verified: gh api repos/astarte-platform/astarte/releases -> v1.3.5 stable 2026-10-05 / v1.4.0-rc.6 prerelease, unchanged since this morning's run; gh api '.../commits?since=2026-10-08T11:39:45Z' -> empty (last master commit 171c1bd7, already measured inert) ; pass
Unsure: nothing
Follow-ups: none
```
