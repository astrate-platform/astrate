slug: recipe-astarte-upstream
verdict: timeout
at:  f0ac79c
ran: 2026-09-26T17:34:18Z on DietPi in 1200s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"' 2>&1
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z

$ rg -n -i 'target|upstream v?1\.[0-9]' docs/COMPATIBILITY.md | head -40
79:     Recorded against upstream v1.2.0 in
182:     upstream 1.2.0 on 2026-08-24 (`test/conformance/upstream/verify-versions.json`,
270:13. **AMQP trigger actions are rejected at creation** — upstream 1.2.2 accepts
335:17. **Always synchronous where upstream 1.4 defaults to asynchronous** —
336:    upstream 1.4 runs realm create/delete (Housekeeping), interface
348:    `docs/UPSTREAM-EXPERIMENTAL.md` (#68) until upstream 1.4 is final.

$ ls .mule/ && ls .mule/recipes/ && wc -l .mule/todo.md
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
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md
230 .mule/todo.md

$ rg -n 'compat-note|1\.3\.[0-9]|1\.4\.0' .mule/todo.md .mule/for-giulio.md | head -40
.mule/todo.md:3:- [x] compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).
.mule/todo.md:112:- [x] purge-properties-compression-capability [auto]: upstream v1.3.0 adds a `purge_properties_compression_format` device capability (`zlib`|`plaintext`, default `zlib`) — a wire-visible capability value. Check whether Astrate's capabilities handling (internal/broker, the `<realm>/<device_id>/capabilities` topic, issue #16) needs to recognize/honour it, or whether zlib-only is already the deliberate default; propose the change or note why not needed.
.mule/todo.md:113:- [!] empty-introspection-verification [auto]: upstream v1.3.0 changed "allow devices with empty introspection" — verify whether Astrate's device connection/introspection handling currently rejects an empty introspection string where upstream now accepts it, and propose a fix if so. — BLOCKED: wrote nothing
.mule/todo.md:114:- [!] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages" — investigate only: does Astrate have an analogous delay between trigger install and first delivery? Report, do not patch. — BLOCKED: wrote nothing
.mule/todo.md:115:- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
.mule/todo.md:119:- [!] probe-property-resend-encoding [auto]: upstream v1.3.0 fixed outbound server-property values sent to a device on connect/emptyCache (commit 522ccf4f — a raw stored binaryblob is now wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding, so it ships as a BSON binary, not a string). Investigate only: does Astrate's resendServerProperties/rehydrate path (internal/engine/control.go:144, pkg/payload) emit binaryblob → BSON subtype-0 binary and datetime → UTC datetime exactly like upstream 1.3? Report, do not patch. — BLOCKED: wrote nothing
.mule/todo.md:120:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
.mule/todo.md:126:- [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
.mule/todo.md:137:- [x] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
.mule/todo.md:138:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
.mule/todo.md:177:- [x] compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.
.mule/for-giulio.md:13:- **github-issues triage run, 2026-09-26: nothing proposable (as every run since 2026-09-05), and the 19-issue alarm pile is a false alarm — it says the mule has landed nothing in 19 days and it has landed something every day.** 23 open issues, the same set plus **#112** (today, 09-26 11:04:46Z, "nothing has landed in 14h"). Zero machine-checkable candidates, so **no `.mule/todo.md` lines, no `gh issue create`, nothing commented, closed or edited on GitHub**. Dispositions unchanged: **#93** aclhook comment (`mule-review`, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction.
.mule/for-giulio.md:18:  **Correction on a fact 19 entries have repeated as settled: there is no `.mule/waiting-on.md`.** #92's own comment (2026-09-04) states the issue "has a row in `.mule/waiting-on.md`" and that the weekly `mule-upstream-watch` job reads that file as its first step. The file does not exist here, and `rg 'waiting-on|upstream-watch'` over the whole tree — including gitignored `.mule/` and `tools/` — matches only prose in this file and copies of it under `.mule/reports/log/`. So the escalation that was supposed to catch a stable upstream `v1.4.0` is not in the repo: either it lives in a Pi-side unit outside it, or it was never wired. One command settles it: `systemctl list-timers 'mule-upstream*' --all; grep -rl 'upstream-watch\|waiting-on' /etc/systemd/system /root 2>/dev/null`. If it was never wired, #92 is parked on a condition nothing is watching, and the 2026-09-04 claim that "the mechanism was added because it was missing here" is false. The gate is being caught by hand in the meantime — today's milestone run re-derived it from `gh api .../releases`.
.mule/for-giulio.md:22:- **Milestone recipe run, 2026-09-26: v3.0 verified unchanged — still gated on upstream. One new thing for you, and it is a `milestones.md` fix: the v3.0 section still says "Status: not started", which is no longer true.** Release sweep (`gh api repos/astarte-platform/astarte/releases`, run this session) still shows **no stable v1.4.0**: newest stable `v1.3.4` (2026-09-18), newest overall `v1.4.0-rc.5` (2026-08-20, prerelease) — byte-identical to what the 2026-09-24 run saw, so nothing shipped upstream in the two days since. `gh issue list --label milestone-3.0 --state all --limit 50` is still **empty**: this milestone's whole ledger lives under `upstream-parity` instead. Open non-alarm issues remain exactly **#93** (aclhook comment, `mule-review`, its own recipe path), **#92** (keyAgreement, parked on a stable v1.4.0 per your 2026-09-04 decision), **#78** (FDO, `milestone-4.0`), **#1** (untouched per standing instruction). No new gaps, **no issues filed, no `.mule/todo.md` task lines added** — there is nothing machine-checkable left, and filing anything here would duplicate #92/#93/#78.
.mule/for-giulio.md:23:  **New, and the reason this run is not a copy of the last thirteen: `.mule/milestones.md:127-128` still reads "Status: not started. First recipe job: triage #47–#89 into an ordered plan (which are audits vs features vs decisions), file sub-issues where work splits, escalate the 'decide' set in one batch."** All of that happened: of #47–#89 every issue is closed except the deliberately-parked #78. The section describes as future work a triage that is finished, and it is the exact text every future run of this recipe reads to pick its target — so the recipe has been re-deriving "the backlog is done, the gate is upstream" from scratch on a timer for three weeks. Suggest rewriting those two lines to state the real position: *1.3 surface delivered (#47–#89 closed); milestone blocked on upstream v1.4.0 going stable; final phase = reconcile the two `docs/UPSTREAM-EXPERIMENTAL.md` rows, answer #92, bump `APICompatVersion`.* That also makes the blocking condition greppable instead of inferred.
.mule/for-giulio.md:25:  **Standing ask, unchanged, for whenever a stable v1.4.0 ships:** answer #92, reconcile the two experimental rows, run the final-phase `APICompatVersion` bump (and update `docs/COMPATIBILITY.md` deviation #10 to match), then cut the v3.0 tag. One process note while this stays blocked: this is the fifteenth milestone-run entry in this file (fourteen prior, 2026-09-06 → 2026-09-24) and the header asks for a queue rather than a log — worth deleting the superseded duplicates whenever you next prune, or the signal-to-noise here keeps falling.
.mule/for-giulio.md:93:- **github-issues triage run, 2026-09-25: still nothing proposable — the mule-alarm pile is 18 issues (#94–#111), and two corrections to what earlier runs inferred about it.** Twenty-two open issues, same set as every run since 2026-09-05: **#111–#94** are the daily `mule-alarm` noise (zero comments, telemetry not code), **#93** aclhook comment (mule-review, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable upstream v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-24 run: **#111** (today, 09-25 11:01:07Z, "nothing has landed in 16h"). Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines added**. **Proposal, extending the 2026-09-24 one: close #94–#110**; leave #111 (live today) to self-expire. Two corrections, both from reading `tools/mule.sh` rather than the issue titles: **(1) the "missed days" in earlier entries (09-15, 09-20, 09-23) are evidence of activity, not its absence** — I previously called them timer gaps; the alarm is gated on a `.alarmed` sentinel that `beat` deletes on every land (`mule.sh:737,748-749`), so a day with no new issue is a day something landed, not a day the timer slept. **(2) the alarm body's closing line is wrong**: "Close it — the mule reopens a new one if the silence continues" (`mule.sh:764`) implies issue state drives the alarm; it does not. Nothing reads or reconciles the 18 open issues, so the pile can only grow and closing the old ones is safe but inert — and `cmd_refill` commits `mule: refill the queue` (mule.sh:719-720) **without calling `beat`**, so a doc-only refill never counts as landing work. The titles' ages also imply a beat roughly every day around 19:00–21:00Z (#111 filed 09-25 11:01Z at 16h → last beat ≈ 09-24 19:01Z; #110 filed 09-24 10:54Z at 15h → ≈ 09-23 19:54Z) while the ~11:00Z daily check has never once found under the 8h threshold — **what calls `beat` at that time is unverified**, and one command on the Pi settles it (`grep -E 'landed|checked:' /root/astrate-mule/.mule/log | tail -20`). The queue has still landed nothing since ~2026-09-04/05; if that idle streak is not intentional, the wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:97:- **Milestone recipe run, 2026-09-24: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-23 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:101:- **github-issues triage run, 2026-09-24: still nothing proposable — the mule-alarm pile is now 17 days (#94–#110).** Twenty-one open issues, same set: **#93** aclhook comment (mule-review, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable upstream v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-20 runs: the alarms **#108** (09-21), **#109** (09-22) and **#110** (today, 09-24 10:54Z, the live one) — a four-day gap on 09-23 not withstanding. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines added**. **Proposal, extending the 2026-09-20 one: close #94–#109** — each was a one-day event, superseded by the next day's alarm, never actionable; leave #110 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:105:- **Milestone recipe run, 2026-09-23: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-21 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:109:- **github-issues triage run, 2026-09-23: still nothing proposable — the mule-alarm pile is now 16 alarms (#94–#109).** Twenty open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-21 run: **#109**, yesterday's alarm (created 09-22 10:54Z, "nothing has landed in 14h"); **#108** has expired into the pile (superseded by #109). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-21 one: close #94–#108** — each was a one-day event, superseded by next day's alarm, never actionable; leave #109 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — now a ~17-day idle streak — if that is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:117:- **github-issues triage run, 2026-09-21: still nothing proposable — the mule-alarm pile is now 16 days (#94–#108).** Eighteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-20 evening run: **#108**, today's alarm (created 09-21 10:54Z, "nothing has landed in 15h"); **#107** has expired into the pile (superseded by #108). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-20 one: close #94–#107** — each was a one-day event, superseded by next day's alarm, never actionable; leave #108 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:121:- **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:125:- **github-issues triage re-run, 2026-09-20 (evening): nothing proposable, no change since this morning's 13:10 run.** Re-surveyed with the recipe command: the issue set is identical to the 2026-09-20 morning triage — **#94–#107** are still the daily `mule-alarm` noise (14 issues, zero comments, telemetry not code), **#93** aclhook comment rewrite is still mule-review with `8c61268` already pushed (its own recipe path), **#92** keyAgreement is still parked on a stable upstream v1.4.0 with its waiting-on.md row and the mule-upstream-watch escalation wired (an -rc does not satisfy it), **#78** FDO is still the milestone-4.0 design/investigation already escalated, and **#1** is untouched per standing instruction. No new issue has been filed since the morning run. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines proposed**. This morning's proposal stands unchanged: **close #94–#106** — each was a one-day event, superseded by the next day's alarm, never actionable — and leave **#107** (still the newest, no #108 yet) to self-expire. The queue idle streak since ~2026-09-04/05 persists; if that is not intentional the wedge is on the Pi (queue review = the fix, per the 13:10 line).
.mule/for-giulio.md:129:- **Milestone recipe run, 2026-09-20: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-19 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:137:- **github-issues triage run, 2026-09-20: still nothing proposable — the mule-alarm pile is now 15 days (#94–#107).** Seventeen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-18 run: **#107**, yesterday's alarm (created 09-19 10:49Z, "nothing has landed in 14h"); **#106** has expired into the pile (superseded by #107). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-18 one: close #94–#106** — each was a one-day event, superseded by next day's alarm, never actionable; leave #107 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:141:- **Milestone recipe run, 2026-09-19: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-18 runs. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:145:- **Milestone recipe run, 2026-09-18 (second, same day): v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable is today's `v1.3.4` (2026-09-18, maintenance-only — data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change, no new gap), newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:153:- **COMPATIBILITY.md wording update for upstream v1.3.4 (newest stable, 2026-09-18; v1.4.0 is still rc.5-only) — folds v1.3.3 and v1.3.4 into the open v1.3.2 proposal below; both are maintenance-only and wire-inert, so that proposal stands unchanged.** v1.3.3 (2026-08-07) was an empty-body patch and v1.3.4 (2026-09-18, `astarte_data_updater_plant` only — more frequent fullsweep GC on AMQPDataConsumer processes plus an RPC-availability fix; release body via `gh api repos/astarte-platform/astarte/releases` `v1.3.4`) introduce **no wire/API surface change**, so the v1.3.2 wording proposal stands complete: same decision (adopt the 1.3 line, then doc + `APICompatVersion` bump together per the bump rule, or keep 1.2.2 and add the "not yet emulated" note), same "not yet emulated" list, same deviation 18 reword ("added by upstream v1.3"). The only delta a v1.3.4-aware doc carries is the version reference: the proposed §Infrastructure-differences sentence's "until the milestone that adopts v1.3.2 as the target" reads "until the milestone that adopts the v1.3 line (**newest stable v1.3.4**) as the target". Raw: [v1.3.4](https://github.com/astarte-platform/astarte/releases/tag/v1.3.4).
.mule/for-giulio.md:157:- **github-issues triage run, 2026-09-18: still nothing proposable — the mule-alarm pile is now 13 straight days (#94–#106).** Sixteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-17 run: **#106**, today's alarm (created 09-18 10:54Z, "nothing has landed in 14h"); **#105** has expired into the pile (superseded by #106). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-17 one: close #94–#105** — each was a one-day event, superseded by next day's alarm, never actionable; leave #106 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:161:- **Milestone recipe run, 2026-09-18: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) shows a new **stable v1.3.4** (2026-09-18) — maintenance only (data_updater_plant GC-sweep + RPC-availability fixes, no wire/API surface change), so no new gap for Astrate — and still **no stable v1.4.0**: newest overall remains `v1.4.0-rc.5` (2026-08-20). `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:169:- **github-issues triage run, 2026-09-17: still nothing proposable — the mule-alarm pile is now 12 straight days (#94–#105).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-16 run: **#105**, today's alarm (created 09-17 10:51Z, "nothing has landed in 14h"); **#104** has expired into the pile (superseded by #105). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-16 one: close #94–#104** — each was a one-day event, superseded by next day's alarm, never actionable; leave #105 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:171:- **github-issues triage run, 2026-09-16: still nothing proposable — the mule-alarm pile is now 11 straight days (#94–#104).** Fifteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), plus the daily alarm. New since the 2026-09-15 run: **#104**, today's alarm (created 09-16 11:12Z, "nothing has landed in 14h"); **#103** has expired into the pile (superseded by #104). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-15 one: close #94–#103** — each was a one-day event, superseded by next day's alarm, never actionable; leave #104 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:173:- **Milestone recipe run, 2026-09-15: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-13 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:175:- **github-issues triage run, 2026-09-15: still nothing proposable — the mule-alarm pile is now 10 straight days (#94–#103).** Fourteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-13 run: **#103**, the alarm created 2026-09-14 11:27Z ("nothing has landed in 16h"); **#102** has expired into the pile (superseded by #103). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-13 one: close #94–#102** — each was a one-day event, superseded by next day's alarm, never actionable; leave #103 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:177:- **Milestone recipe run, 2026-09-13: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-12 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, its own recipe path), **#78 FDO** (milestone-4.0, already escalated). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:179:- **github-issues triage run, 2026-09-13: still nothing proposable — the mule-alarm pile is now 9 straight days (#94–#102).** Thirteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-12 run: **#102**, today's alarm (created 11:09Z, "nothing has landed in 14h"); **#101** has expired into the pile (superseded by #102). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-12 one: close #94–#101** — each was a one-day event, superseded by next day's alarm, never actionable; leave #102 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:181:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:185:- **github-issues triage run, 2026-09-12: still nothing proposable — the mule-alarm pile is now 8 straight days (#94–#101).** Twelve open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-11 run: **#101**, today's alarm (created 10:58Z, "nothing has landed in 14h"); **#100** has expired into the pile (superseded by #101). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-11 one: close #94–#100** — each was a one-day event, superseded by next day's alarm, never actionable; leave #101 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:187:- **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pile is now 7 straight days and worth a look.** Twelve open issues. Still no machine-checkable fix candidates, so no task lines: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-10 run: **#100**, today's alarm (created 10:55Z, "nothing has landed in 15h"). #99–#94 are the same `mule-alarm` (one per day 2026-09-05→09-10, ~11:00Z each, zero comments) — telemetry, not code issues, so never proposable. But seven consecutive alarms is past the "low-activity window" wording earlier runs used: the mule has landed no commit since ~2026-09-04/05, and `main`'s queue copy shows the top line as `- [ ]` while the landed work below is all `[!] BLOCKED` (`wrote nothing` / `tests failed`) — though `mule/queue` is the authoritative copy. **Proposal: close #94–#99** — each was a one-day event, superseded by next day's alarm, never actionable; leave #100 (live today) to self-expire. If the idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.

$ rg -n '^- \[ \]' .mule/todo.md | head -30; echo ---; git -C /root/astrate-mule log --oneline -5 -- .mule/todo.md
86:- [ ] broker-external-bus-intake [legion]: implement the second Intake implementation named in the TODO at internal/broker/intake.go:56-62 — an external-bus-backed (e.g. NATS JetStream) Intake for multi-instance deployment / restart survival, reproducing the embedded broker's per-device ordering, deferred-ack backpressure and QoS 0 drop semantics; an implementation of this interface needs containerised integration tests, so it runs on the Legion Go.
88:- [ ] hygiene-govulncheck [legion]: on the Legion Go, run govulncheck ./... on a fresh ~/astrate clone and file a task line for any REACHABLE advisory it reports (name the CVE and the call path, per the hygiene recipe); this box cannot build it — go install golang.org/x/vuln/cmd/govulncheck@latest is OOM-killed here (3.7GB, 4 cores), so the recipe's highest-priority check never runs unless done there.
89:- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.go:560) with no stop signal, so a block-death racing process shutdown (cmd/astrate/main.go:482-488 drains the Manager but never cancels the goroutine) can rebuild the flow via `mgr.StartFlow` on a background pump (internal/flow/flow.go:173-175) after `Manager.Shutdown` and flip the durable status back to running. Thread a stop channel/context through the Service, checked in the loop's sleep, cancelled on shutdown; verification needs a `[legion]` integration or timing-based test. [auto]
101:- [ ] race-check-store: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/store/... ./internal/housekeeping/... ./migrations/...`. Report any failure to .mule/for-giulio.md with the full race report. Split out of the former single `race-check` line, which timed out running the whole tree at once. [legion] [readonly]
102:- [ ] race-check-engine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/engine/... ./internal/broker/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
103:- [ ] race-check-flow: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/flow/... ./internal/realm/... ./internal/pairing/... ./internal/auth/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
104:- [ ] race-check-appengine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/appengine/... ./internal/observability/... ./internal/httpx/... ./internal/config/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
105:- [ ] race-check-pkg: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./pkg/... ./cmd/... ./internal/testutil/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
200:- [ ] deviceid-trailing-bits-upstream-probe [legion]: on the Legion Go, probe the running upstream Astarte (or a locally-run `elixir -e` with `Base.url_decode64!`) with the device ID `"AAAAAAAAAAAAAAAAAAAAAB"` — nonzero unused bits in the 22nd base64url char — and report whether upstream accepts it: `pkg/deviceid/deviceid.go:34-37` uses `base64.RawURLEncoding.Strict()` and deviceid_test.go:176 asserts rejection while the comment claims parity with `Elixir Base.url_decode64!(padding: false)`, which decodes and discards those bits rather than erroring. Astrate may therefore be stricter than upstream on the same wire form (a device registered with such an id on upstream would be rejected here). If upstream rejects it identically, the strict encoding is verified parity and the task is done; if upstream accepts it, escalate the relax-vs-strict call to `.mule/for-giulio.md` with the measurement. Probe first, no code change either way.
203:- [ ] store-devices-inhibit-re-register [legion] [auto]: `RegisterDevice` (internal/store/devices.go:75-91) silently clears the inhibit flag — its `ON CONFLICT ... DO UPDATE SET status = 'registered'` fires for any device with `first_credentials_request IS NULL`, including one the admin inhibited via `SetDeviceInhibited` (devices.go:251-268), which sets `status='inhibited'` on unconfirmed devices too; §5.3 says an inhibited device blocks new credentials and connections, so the re-registration re-opens it. Preserve `'inhibited'` in the SET (e.g. `status = CASE WHEN devices.status = 'inhibited' THEN 'inhibited' ELSE 'registered' END`) and add a Lifecycle case in internal/store/devices_test.go: inhibit an unconfirmed device, re-register, assert status stays inhibited and the secret still rotates. Verify the assert against upstream's register-not-touching-inhibit on the Legion while the integration suite runs.
205:- [ ] store-devices-alias-lowest-id-test [legion] [auto]: pin the documented tie-break of `GetDeviceByAlias` (internal/store/devices.go:130-131, "if several devices share an alias the lowest device ID wins", `ORDER BY id LIMIT 1`) — no test drives it today; add a case in internal/store/devices_test.go that gives two devices the same alias and asserts the lookup resolves the lower ID. Needs the integration DB.
221:- [ ] flow-boot-resumes-stopped-durable-flows [legion] [auto]: RehydrateAutoRestart (cmd/astrate/main.go:216) restarts every flow with auto_restart=true whatever its last status — `ListAutoRestartFlows` filters on `f.auto_restart = true` alone (internal/store/flows.go:125) and nothing ever clears the flag, which is written once at create (flowapi/http.go:141-146 -> service.go:318) and never on stop — so a durable flow stopped through the API is running again after the next boot, while the shutdown mark that says otherwise (MarkRunningFlowsStopped, main.go:497 -> internal/flowapi/service.go:885-901, which writes status="stopped") is never read by that query. The two halves disagree: service.go:325 documents "every durable flow", the shutdown call assumes "only what was running". Say which is the intent, then make them agree — either add `AND f.status = 'running'` to flows.go:125 (making the shutdown mark load-bearing) or delete the mark as dead code — and add a case in internal/store/flows_test.go pinning the chosen rule (a stopped auto_restart row must, or must not, come back). Needs the DB.
---
85b5697 mule: blocked docs-sync-rm-delete-device-async-operation-param
9f20120 mule: log docs-sync-rm-delete-device-async-operation-param
f965f39 mule: docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
7d78bd3 mule: log docs-sync-rm-async-operation-param
3f0eca8 mule: docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ rg -n 'astarte-upstream|upstream watch|upstream-watch' .mule/for-giulio.md | head -8; echo "--- reports today:"; ls -t .mule/reports/ 2>/dev/null | head -8
18:  **Correction on a fact 19 entries have repeated as settled: there is no `.mule/waiting-on.md`.** #92's own comment (2026-09-04) states the issue "has a row in `.mule/waiting-on.md`" and that the weekly `mule-upstream-watch` job reads that file as its first step. The file does not exist here, and `rg 'waiting-on|upstream-watch'` over the whole tree — including gitignored `.mule/` and `tools/` — matches only prose in this file and copies of it under `.mule/reports/log/`. So the escalation that was supposed to catch a stable upstream `v1.4.0` is not in the repo: either it lives in a Pi-side unit outside it, or it was never wired. One command settles it: `systemctl list-timers 'mule-upstream*' --all; grep -rl 'upstream-watch\|waiting-on' /etc/systemd/system /root 2>/dev/null`. If it was never wired, #92 is parked on a condition nothing is watching, and the 2026-09-04 claim that "the mechanism was added because it was missing here" is false. The gate is being caught by hand in the meantime — today's milestone run re-derived it from `gh api .../releases`.
125:- **github-issues triage re-run, 2026-09-20 (evening): nothing proposable, no change since this morning's 13:10 run.** Re-surveyed with the recipe command: the issue set is identical to the 2026-09-20 morning triage — **#94–#107** are still the daily `mule-alarm` noise (14 issues, zero comments, telemetry not code), **#93** aclhook comment rewrite is still mule-review with `8c61268` already pushed (its own recipe path), **#92** keyAgreement is still parked on a stable upstream v1.4.0 with its waiting-on.md row and the mule-upstream-watch escalation wired (an -rc does not satisfy it), **#78** FDO is still the milestone-4.0 design/investigation already escalated, and **#1** is untouched per standing instruction. No new issue has been filed since the morning run. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines proposed**. This morning's proposal stands unchanged: **close #94–#106** — each was a one-day event, superseded by the next day's alarm, never actionable — and leave **#107** (still the newest, no #108 yet) to self-expire. The queue idle streak since ~2026-09-04/05 persists; if that is not intentional the wedge is on the Pi (queue review = the fix, per the 13:10 line).
--- reports today:
log
hygiene.md

$ gh api 'repos/astarte-platform/astarte/commits?sha=master&since=2026-09-18T00:00:00Z&per_page=100' --jq 'length as $n | "count=\($n)", (.[0:15][] | "\(.commit.author.date[0:10])\t\(.sha[0:8])\t\(.commit.message | split("\n")[0])")' 2>&1
count=16
2026-09-24	23bb488e	refactor(generators): move Generators in Umbrella (#2221)
2026-09-24	53778f24	refactor(data_access): move Data Access into Umbrella (#2219)
2026-09-23	f78baa94	Merge pull request #2218 from nedimtokic/move-config
2026-09-23	9c07eb83	refactor(config): update Config references
2026-09-23	d7838d2a	refactor(config): move Astarte Config into Umbrella
2026-09-23	5f46fd34	refactor(events): move Events into Umbrella (#2214)
2026-09-23	7af65c80	fix(core): encode false simple trigger known values (#2215)
2026-09-22	2844f9d6	refactor(secrets): move Astarte Secrets into Umbrella (#2212)
2026-09-22	35cb7d36	refactor(rpc): move Astarte RPC into Umbrella (#2213)
2026-09-18	049bdd45	refactor(adapters): move Astarte Adapters into Umbrella (#2207)
2026-09-16	257086cb	feat(generators): add trigger generator
2026-09-16	cd51834e	feat(generators): add AMQP trigger target generator
2026-09-15	2ef81788	feat(generators): add simple trigger config generator
2026-09-15	461c8a64	feat(generators): add device trigger config generator
2026-09-15	4f99a33f	feat(generators): add data trigger config generator

$ rg -n -i '2215|encode false|known values|known_values' .mule/todo.md .mule/for-giulio.md docs/UPSTREAM-EXPERIMENTAL.md 2>/dev/null | head -20
(no output)

$ gh api repos/astarte-platform/astarte/commits/7af65c80 --jq '.commit.message, (.files[] | "\(.filename)\t+\(.additions)/-\(.deletions)")' 2>&1 | head -30
fix(core): encode false simple trigger known values (#2215)

Signed-off-by: Gabriele Ghio <gabriele.ghio@secomind.com>
libs/astarte_core/lib/astarte_core/triggers/simple_trigger_config.ex	+5/-2
libs/astarte_core/test/astarte_core/triggers/simple_trigger_config_test.exs	+26/-0

$ gh api repos/astarte-platform/astarte/commits/7af65c80 --jq '.files[] | select(.filename|test("simple_trigger_config.ex$")) | .patch' 2>&1
@@ -1,7 +1,7 @@
 #
 # This file is part of Astarte.
 #
-# Copyright 2017 - 2025 SECO Mind Srl
+# Copyright 2017 - 2026 SECO Mind Srl
 #
 # Licensed under the Apache License, Version 2.0 (the "License");
 # you may not use this file except in compliance with the License.
@@ -469,7 +469,7 @@ defmodule Astarte.Core.Triggers.SimpleTriggerConfig do
       group_name: group_name,
       interface_name: interface_name,
       interface_major: interface_major,
-      known_value: known_value && Cyanide.encode!(%{v: known_value}),
+      known_value: encode_known_value(known_value),
       match_path: match_path,
       data_trigger_type: trigger_type,
       value_match_operator: value_match_operator
@@ -484,6 +484,9 @@ defmodule Astarte.Core.Triggers.SimpleTriggerConfig do
     }
   end
 
+  defp encode_known_value(nil), do: nil
+  defp encode_known_value(known_value), do: Cyanide.encode!(%{v: known_value})
+
   defp get_data_trigger_object(%SimpleTriggerConfig{} = config) do
     %SimpleTriggerConfig{
       device_id: device_id,

$ rg -n -i 'known_value|KnownValue' --type go -l | head -20; echo "--- hits:"; rg -n -i 'knownvalue' --type go | head -40
internal/realm/http_test.go
internal/engine/triggers/match_test.go
internal/engine/triggers/match.go
internal/engine/triggers/change_test.go
--- hits:
internal/engine/triggers/match.go:248:	KnownValue         json.RawMessage `json:"known_value"`
internal/engine/triggers/match.go:367:	if op != anyToken && op != "" && len(c.KnownValue) == 0 {
internal/engine/triggers/match.go:388:	if len(c.KnownValue) > 0 {
internal/engine/triggers/match.go:389:		if err := json.Unmarshal(c.KnownValue, &m.known); err != nil {

→ Read internal/engine/triggers/match.go [offset=230, limit=180]

$ rg -n 'known' internal/engine/triggers/match.go | sed -n '1,60p'
217:	known     any
248:	KnownValue         json.RawMessage `json:"known_value"`
257:// that have no field shape (does-not-parse, no simple_triggers, unknown
309:		return fmt.Errorf("unknown trigger type %q", c.Type)
319:	_, known := dataOns[c.On]
322:	} else if !known {
368:		fe.add("known_value", "can't be blank")
380:		evaluated: known && c.GroupName == "",
389:		if err := json.Unmarshal(c.KnownValue, &m.known); err != nil {
390:			return fmt.Errorf("known_value does not parse: %w", err)
393:	t.noteUnsupported(!known, "data_trigger on "+c.On)
406:	evaluated, known := deviceOns[c.On]
409:	} else if !known {
601:	return valueMatches(m.operator, ev.Value, m.known)
652:func valueMatches(op string, value, known any) bool {
661:		eq, comparable := valueEquals(value, known)
664:		eq, comparable := valueEquals(value, known)
668:		b, bok := numeric(known)
683:		return containsValue(value, known)
686:		return applicable && !c(known)
692:// valueEquals compares a typed payload value with a JSON known_value; the
694:func valueEquals(value, known any) (eq, comparable bool) {
696:		b, ok := numeric(known)
701:		s, ok := known.(string)
704:		b, ok := known.(bool)
729:func containsValue(value, known any) bool {
731:	return applicable && c(known)
736:func containable(value any) (func(known any) bool, bool) {
739:		return func(known any) bool {
740:			s, ok := known.(string)
760:func sliceContains[T any](xs []T) func(known any) bool {
761:	return func(known any) bool {
763:			if eq, ok := valueEquals(xs[i], known); ok && eq {

→ Read internal/engine/triggers/match.go [offset=645, limit=80]

$ rg -n -i 'known_value.*false|false.*known_value' internal/engine/triggers/*_test.go | head -20
(no output)

$ gh api 'repos/astarte-platform/astarte/commits?sha=master&since=2026-09-18T00:00:00Z&per_page=100' --jq '.[] | "\(.commit.author.date[0:10])\t\(.sha[0:8])\t\(.commit.message | split("\n")[0])"' 2>&1 | tail -5
2026-09-16	cd51834e	feat(generators): add AMQP trigger target generator
2026-09-15	2ef81788	feat(generators): add simple trigger config generator
2026-09-15	461c8a64	feat(generators): add device trigger config generator
2026-09-15	4f99a33f	feat(generators): add data trigger config generator
2026-09-15	d2caff5c	feat(generators): add simple trigger scope generator

$ gh api repos/astarte-platform/astarte/commits/257086cb --jq '.commit.message, (.files[] | .filename)' 2>&1 | head -25
feat(generators): add trigger generator

Signed-off-by: Gabriele Ghio <gabriele.ghio@secomind.com>
libs/astarte_generators/lib/astarte/core/generators/triggers/trigger.ex
libs/astarte_generators/test/astarte/core/generators/triggers/trigger_test.exs

$ rg -n -i 'master|unreleased|not yet released' .mule/todo.md | head -10
225:- [x] cmd-loadsealer-masterkeyfile-test [auto]: add a container-free test for `loadSealer` (cmd/astrate/main.go:515-526) — the `security.master_key_file` to env hand-off (`os.Setenv(store.EnvMasterKeyFile, ...)`, 517) is the only wiring that makes a configured key file reach the store loader, and dropping or misspelling that branch leaves the file silently ignored and boot dying with the generic "master key" error (523) even though the config named a valid file. Assert a 32-byte key file under t.TempDir() loads, the inline ASTRATE_MASTER_KEY path still wins, and a missing key returns the error naming all three env vars. t.Setenv both ASTRATE_MASTER_KEY and ASTRATE_MASTER_KEY_FILE first so the os.Setenv is restored at cleanup.

$ gh api 'repos/astarte-platform/astarte/contents/libs/astarte_core/lib/astarte_core/triggers/simple_trigger_config.ex?ref=master' --jq '.content' 2>/dev/null | base64 -d > /tmp/opencode/stc.ex; wc -l /tmp/opencode/stc.ex; rg -n 'known_value' /tmp/opencode/stc.ex
689 /tmp/opencode/stc.ex
51:    field :known_value, :any, virtual: true
63:      |> put_if("known_value", config.known_value, config.value_match_operator != "*")
95:    :known_value
368:      |> delete_change(:known_value)
371:      |> validate_required([:known_value])
460:      known_value: known_value,
472:      known_value: encode_known_value(known_value),
487:  defp encode_known_value(nil), do: nil
488:  defp encode_known_value(known_value), do: Cyanide.encode!(%{v: known_value})
613:      known_value: known_value
621:    decoded_known_value =
622:      if known_value do
623:        Cyanide.decode!(known_value)
640:      known_value: decoded_known_value

$ rg -n -i 'known_value' internal/engine/triggers/match_test.go internal/engine/triggers/change_test.go internal/realm/http_test.go | head -20
internal/engine/triggers/change_test.go:45:			"known_value": 42
internal/engine/triggers/match_test.go:23:			"known_value": 0.4
internal/engine/triggers/match_test.go:188:			"known_value": ` + known + `}`
internal/engine/triggers/match_test.go:344:			"known_value": 1}]}`,
internal/engine/triggers/match_test.go:349:			want: "known_value=can't be blank"},
internal/engine/triggers/match_test.go:353:			want: "known_value=can't be blank, value_match_operator=is invalid"},
internal/engine/triggers/match_test.go:429:// "/*"+"==" definition demands only the operator fix even though known_value
internal/engine/triggers/match_test.go:438:			 "known_value": 1}
internal/realm/http_test.go:714:			`"interface_major":1,"match_path":"/*","value_match_operator":"==","known_value":1}]}`

$ sed -n '600,645p' /tmp/opencode/stc.ex; echo "=== changeset region ==="; sed -n '355,380p' /tmp/opencode/stc.ex
{any_device_id, SimpleTriggersUtils.object_type_to_int!(:any_device)}
    end
  end

  defp from_data_trigger(%DataTrigger{} = data_trigger) do
    %DataTrigger{
      data_trigger_type: data_trigger_type,
      group_name: group_name,
      device_id: device_id,
      interface_name: interface_name,
      interface_major: interface_major,
      value_match_operator: value_match_operator,
      match_path: match_path,
      known_value: known_value
    } = data_trigger

    condition = Map.fetch!(@data_trigger_condition_to_string, data_trigger_type)

    value_match_operator_string =
      Map.fetch!(@data_trigger_operator_to_string, value_match_operator)

    decoded_known_value =
      if known_value do
        Cyanide.decode!(known_value)
        |> Map.get("v")
      else
        nil
      end

    # TODO: interface_name and interface_major can't be deducted from interface_id,
    # leaving them nil waiting for an API to retrieve them
    %SimpleTriggerConfig{
      type: "data_trigger",
      on: condition,
      device_id: normalize_proto_string_default(device_id),
      group_name: normalize_proto_string_default(group_name),
      interface_name: interface_name,
      interface_major: interface_major,
      value_match_operator: value_match_operator_string,
      match_path: match_path,
      known_value: decoded_known_value
    }
  end

  defp from_device_trigger(%DeviceTrigger{} = device_trigger, object_id, object_type) do
    %DeviceTrigger{
=== changeset region ===
      if get_field(changeset, :value_match_operator) != "*" do
        add_error(changeset, :value_match_operator, "must be * when match_path is /*")
      else
        changeset
      end
    else
      validate_format(changeset, :match_path, Mapping.mapping_regex())
    end
  end

  defp validate_match_parameters(%Ecto.Changeset{} = changeset) do
    if get_field(changeset, :value_match_operator, "*") == @data_trigger_any_match_operator do
      changeset
      |> delete_change(:known_value)
    else
      changeset
      |> validate_required([:known_value])
    end
  end

  defp validate_device_id(%Ecto.Changeset{} = changeset, field) do
    validate_change(changeset, field, fn field, encoded_id ->
      case validate_device_id_or_any(encoded_id) do
        :ok ->
          []

$ mkdir -p /tmp/opencode/probe && cat > /tmp/opencode/probe/main.go <<'EOF'
package main

import (
	"encoding/json"
	"fmt"
)

type c struct {
	KnownValue json.RawMessage `json:"known_value"`
}

func main() {
	for _, body := range []string{
		`{"known_value":false}`,
		`{"known_value":null}`,
		`{}`,
		`{"known_value":0}`,
		`{"known_value":0.0}`,
		`{"known_value":""}`,
	} {
		var v c
		if err := json.Unmarshal([]byte(body), &v); err != nil {
			fmt.Printf("%-26s ERR %v\n", body, err)
			continue
		}
		var any1 any
		if len(v.KnownValue) > 0 {
			if err := json.Unmarshal(v.KnownValue, &any1); err != nil {
				fmt.Printf("%-26s parse-err %v\n", body, err)
				continue
			}
		}
		fmt.Printf("%-26s len=%d blank-check(len==0)=%-5v decoded=%#v\n", body, len(v.KnownValue), len(v.KnownValue) == 0, any1)
	}
}
EOF
cd /tmp/opencode/probe && go mod init probe >/dev/null 2>&1; go run main.go
{"known_value":false}      len=5 blank-check(len==0)=false decoded=false
{"known_value":null}       len=4 blank-check(len==0)=false decoded=<nil>
{}                         len=0 blank-check(len==0)=true  decoded=<nil>
{"known_value":0}          len=1 blank-check(len==0)=false decoded=0
{"known_value":0.0}        len=3 blank-check(len==0)=false decoded=0
{"known_value":""}         len=2 blank-check(len==0)=false decoded=""

$ rg -n 'simple_triggers|SimpleTriggers' internal/realm/*.go internal/triggers/*.go 2>/dev/null | grep -v _test | head -20
internal/realm/http.go:266:// definition: {"action": {...}} and/or "simple_triggers" as the index-aligned
internal/realm/http.go:273:	if len(te.SimpleTriggers) > 0 {
internal/realm/http.go:274:		arr := make([]map[string][]string, len(te.SimpleTriggers))
internal/realm/http.go:275:		for i, fields := range te.SimpleTriggers {
internal/realm/http.go:282:		body["simple_triggers"] = arr
internal/realm/service.go:434:		// across action and simple_triggers (issues #63/#70); the HTTP layer

$ rg -rn --files-with-matches 'trigger' internal/realm/ | grep -v _test | head
internal/realm/service.go
internal/realm/policy.go
internal/realm/http.go

$ rg -n 'known|KnownValue' internal/realm/service.go | head; echo "=== trigger store round-trip:"; rg -n 'definition|Definition' internal/store/triggers.go 2>/dev/null | head -20
136:// realmID resolves a realm name to its id; an unknown realm surfaces
138:// middleware normally rejects unknown realms first).
362:// name (upstream GET /interfaces/{name}), ascending. An unknown name yields
=== trigger store round-trip:
19:	Definition []byte
23:func (s *Store) CreateTrigger(ctx context.Context, realmID int16, name string, definition []byte) (*Trigger, error) {
24:	t := Trigger{RealmID: realmID, Name: name, Definition: definition}
26:		`INSERT INTO triggers (realm_id, name, definition) VALUES ($1, $2, $3) RETURNING id`,
27:		realmID, name, definition).Scan(&t.ID)
41:		`SELECT id, realm_id, name, definition FROM triggers WHERE realm_id = $1 AND name = $2`,
42:		realmID, name).Scan(&t.ID, &t.RealmID, &t.Name, &t.Definition)
69:		`SELECT id, realm_id, name, definition FROM triggers WHERE realm_id = $1 ORDER BY name`,
79:		if err := rows.Scan(&t.ID, &t.RealmID, &t.Name, &t.Definition); err != nil {

$ gh api 'repos/astarte-platform/astarte/contents/libs/astarte_generators?ref=master' --jq '.[].name' 2>&1 | head; echo "=== mix.exs apps:"; gh api 'repos/astarte-platform/astarte/contents/mix.exs?ref=master' --jq '.content' 2>/dev/null | base64 -d | rg -n 'generators' | head
{"message":"Not Found","documentation_url":"https://docs.github.com/rest/repos/contents#get-repository-content","status":"404"}gh: Not Found (HTTP 404)
=== mix.exs apps:

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.4") | .body' 2>&1 | head -40
### Fixed

- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
  uncollected binaries over time by forcing more frequent full sweep garbage collections
  on them.
- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.

$ gh api 'repos/astarte-platform/astarte/git/trees/v1.3.4?recursive=0' --jq '.tree[].path' 2>&1 | head -20; echo "=== generators in v1.3.4?"; gh api 'repos/astarte-platform/astarte/contents/libs?ref=v1.3.4' --jq '.[].name' 2>&1 | tr '\n' ' '
.adr-dir
.credo.exs
.dockerignore
.dprint.json
.env
.github
.github/codecov.yml
.github/dco.yml
.github/labeler.yml
.github/pull_request_template.md
.github/pull_request_template.md.license
.github/workflows
.github/workflows/astarte-appengine-api-workflow.yaml
.github/workflows/astarte-apps-build-workflow.yaml
.github/workflows/astarte-build-workflow.yaml
.github/workflows/astarte-dashboard-compatibility.yml
.github/workflows/astarte-dashboard-docker-build.yml
.github/workflows/astarte-dashboard-tests.yml
.github/workflows/astarte-data-updater-plant-workflow.yaml
.github/workflows/astarte-e2e-build-workflow.yaml
=== generators in v1.3.4?
astarte_adapters astarte_data_access astarte_events astarte_generators astarte_rpc astarte_test_suite

