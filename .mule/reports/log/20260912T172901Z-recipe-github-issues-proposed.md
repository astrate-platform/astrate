slug: recipe-github-issues
verdict: proposed
at:  1c8d6f4
ran: 2026-09-12T17:29:01Z on DietPi in 164s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md

$ ls -la .mule/ && ls -la .mule/recipes/ 2>/dev/null
total 3060
drwxr-xr-x  7 root root    4096 Sep 12 19:08 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      13 Sep 12 19:26 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 12 18:45 .heartbeat
-rw-r--r--  1 root root       0 Sep 12 19:26 .last-output
-rw-r--r--  1 root root       0 Sep 12 19:26 .lock
-rw-r--r--  1 root root      14 Sep 12 19:26 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 2957557 Sep 12 19:26 cron.log
drwxr-xr-x  2 root root    4096 Sep 11 20:29 failed
-rw-r--r--  1 root root   29801 Sep 12 19:12 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    5266 Sep 12 18:45 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 11 21:40 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   43263 Sep 12 18:45 todo.md
total 56
drwxr-xr-x 2 root root 4096 Sep  4 20:46 .
drwxr-xr-x 7 root root 4096 Sep 12 19:08 ..
-rw-r--r-- 1 root root 2772 Jul 27 05:27 astarte-upstream.md
-rw-r--r-- 1 root root 3184 Jul 27 05:27 benchmarks.md
-rw-r--r-- 1 root root 2541 Jul 27 05:27 code-review.md
-rw-r--r-- 1 root root 5525 Jul 27 22:18 cross-project-survey.md
-rw-r--r-- 1 root root 2040 Jul 27 05:27 docs-sync.md
-rw-r--r-- 1 root root 2116 Jul 27 05:27 github-issues.md
-rw-r--r-- 1 root root 2103 Jul 27 05:27 hygiene.md
-rw-r--r-- 1 root root 4497 Jul 27 10:28 legion-go.md
-rw-r--r-- 1 root root 5741 Sep  4 20:46 milestones.md

$ gh issue list --state open --limit 40 --json number,title,labels,updatedAt,comments --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
#101 mule: nothing has landed in 14h [mule-alarm ] []c
#100 mule: nothing has landed in 15h [mule-alarm ] []c
#99 mule: nothing has landed in 14h [mule-alarm ] []c
#98 mule: nothing has landed in 16h [mule-alarm ] []c
#97 mule: nothing has landed in 15h [mule-alarm ] []c
#96 mule: nothing has landed in 14h [mule-alarm ] []c
#95 mule: nothing has landed in 15h [mule-alarm ] []c
#94 mule: nothing has landed in 15h [mule-alarm ] []c
#93 In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged [mule-review upstream-parity ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:The mule pushed `8c61268` to `mule/queue` for this.

Unreviewed: the gates passed — it compiles, the tests pass, a new test was shown to fail without the change, and no frozen file was touched. That is not the same as the change being worth having. Nothing merges to `main` until someone reads it:

    bash tools/mule.sh review

Left open on purpose. Whether this actually resolves the issue is your call, not mine. createdAt:2026-09-04T19:54:58Z id:IC_kwDORmfsJs8AAAABSo2DbA includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/93#issuecomment-5545755500 viewerDidAuthor:true]]c
#92 keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds [enhancement upstream-parity upstream-experimental ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:Parked by Giulio on 2026-09-04: **wait for a stable \`v1.4.0\`.** The spec is complete but ships only in \`v1.4.0-rc.5\`, and building against an rc risks rebuilding it.

So this issue is not waiting on anyone to remember it. It has a row in [`.mule/waiting-on.md`](https://github.com/astrate-platform/astrate/blob/main/.mule/waiting-on.md), and the weekly `mule-upstream-watch` job now reads that file as its first step every run: when a stable `v1.4.0` tag exists upstream, it escalates to `.mule/for-giulio.md` naming this issue. An `-rc.N` does not satisfy the row.

That mechanism was added because it was missing here: #51 was parked on the same condition, the condition fired on 2026-08-31, and nothing connected the new upstream tags to the parked issue — it took a hand triage four days later to catch it. createdAt:2026-09-04T19:52:21Z id:IC_kwDORmfsJs8AAAABSo0f8A includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/92#issuecomment-5545730032 viewerDidAuthor:true]]c
#78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate) [enhancement milestone-4.0 upstream-parity ] [map[author:map[login:atsetilam] authorAssociation:MEMBER body:Reframed 2026-08-23 by Giulio's decision: no longer parked — **strategic feature** (zero-touch factory onboarding is a commercial-viability requirement for an IoT platform), milestone-4.0 candidate.

## Verified facts (2026-08-23, upstream master)

- Upstream integrates FDO **inside the Pairing service**: `fdo_onboarding_controller.ex`, session plugs (`setup_fdo`, `fdo_session`), CBOR codec + voucher queries in `astarte_data_access`, and two dedicated libs `astarte_fdo` / `astarte_fdo_core` (SECO Mind, 2025; bugfixes through 2026-05). Not an optional sidecar.
- Device side is the official SDK `astarte-device-fdo-rust` (crate on docs.rs, pushed 2026-08-21) — actively maintained.
- Upstream does NOT implement manufacturing / rendezvous services either; those come from the FIDO Alliance reference stack (`fdo-rs/fido-device-onboard-rs`, used by Red Hat/Fedora IoT).

## Scope decision

Astrate implements **only the last mile**, mirroring upstream's cut:
- owner-side protocol surface in our pairing service (TO1 redirect + TO2 onboarding, CBOR wire format, ownership-voucher storage, owner key management)
- reuse the existing open-source FDO ecosystem for manufacturing/rendezvous — never reimplement it.

**Acceptance:** a device running the official `astarte-device-fdo-rust` SDK completes onboarding against Astrate end-to-end and lands as a provisioned realm device.

**Docs are a first-class deliverable:** upstream ships almost no public FDO documentation (code-level only); ours must exceed that — operator guide covering the full chain (manufacturing → rendezvous → owner) with Astrate-specific setup.

## Next step

Investigation phase before any implementation: read upstream's TO2 handling end-to-end (`libs/astarte_fdo*", pairing controllers/plugs, `fdo/queries`), inventory the exact endpoints, credential/voucher schema and key material required, and measure what our pairing service lacks. Output feeds the v4.0 scope decision in `.mule/milestones.md`. createdAt:2026-08-23T16:16:08Z id:IC_kwDORmfsJs8AAAABQRcugw includesCreatedEdit:false isMinimized:false minimizedReason: reactionGroups:[] url:https://github.com/astrate-platform/astrate/issues/78#issuecomment-5387005571 viewerDidAuthor:true]]c
#1 Provide an Open Source IoT Platform unironically, unironically [wontfix ] []c

$ rg -n "issue-(93|92|78|94|95|96|97|98|99|100|101)|#93|#92|#78" .mule/todo.md .mule/for-giulio.md .mule/waiting-on.md 2>/dev/null | head -40
.mule/for-giulio.md:13:- **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:17:- **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pile is now 7 straight days and worth a look.** Twelve open issues. Still no machine-checkable fix candidates, so no task lines: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-10 run: **#100**, today's alarm (created 10:55Z, "nothing has landed in 15h"). #99–#94 are the same `mule-alarm` (one per day 2026-09-05→09-10, ~11:00Z each, zero comments) — telemetry, not code issues, so never proposable. But seven consecutive alarms is past the "low-activity window" wording earlier runs used: the mule has landed no commit since ~2026-09-04/05, and `main`'s queue copy shows the top line as `- [ ]` while the landed work below is all `[!] BLOCKED` (`wrote nothing` / `tests failed`) — though `mule/queue` is the authoritative copy. **Proposal: close #94–#99** — each was a one-day event, superseded by next day's alarm, never actionable; leave #100 (live today) to self-expire. If the idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:19:- **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:21:- **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
.mule/for-giulio.md:46:  issues: #99–#94 are mule-alarms (idle queue), not code issues. #93 has a pushed commit
.mule/for-giulio.md:47:  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
.mule/for-giulio.md:48:  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated. #1
.mule/for-giulio.md:53:  issues are the same: #92 keyAgreement (upstream decision, gated on stable v1.4.0), #93
.mule/for-giulio.md:54:  aclhook comment (pushed commit `8c61268`, mule-review), #78 FDO (milestone-4.0). Both
.mule/for-giulio.md:63:  queue), not code issues — expected during low-activity windows. #93 has a pushed commit
.mule/for-giulio.md:64:  (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the
.mule/for-giulio.md:65:  waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already escalated below.
.mule/for-giulio.md:78:  remaining open items are already escalated, not fileable gaps: **#92 keyAgreement** (1.4
.mule/for-giulio.md:79:  experimental decision, escalated 2026-09-05, gated on a stable upstream v1.4.0), **#78 FDO**
.mule/for-giulio.md:80:  (milestone-4.0, escalated), **#93** aclhook comment rewrite (mule-review, pushed commit
.mule/for-giulio.md:85:  call when v1.4.0 goes stable: answer #92 (implement against the rc, re-park, or narrow #93
.mule/for-giulio.md:104:  not code issues — expected during low-activity windows. #93 has a pushed commit (`8c61268`)
.mule/for-giulio.md:105:  awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the waiting-on
.mule/for-giulio.md:106:  row. #78 is the milestone-4.0 FDO design/investigation, already escalated below. #1
.mule/for-giulio.md:113:  stable v1.4.0?** (issue #92, `upstream-parity`/`upstream-experimental`). #51 closed
.mule/for-giulio.md:120:  codes), and deleted the "not yet implemented" sentences the old parking quoted (#93 already
.mule/for-giulio.md:126:  re-park until v1.4.0 is a stable tag, or take only the narrow #93 fix. (Escalated again
.mule/for-giulio.md:201:- **#78 FDO device onboarding — milestone-4.0, investigation phase.** Too large for a single
.mule/for-giulio.md:266:#78 is no longer parked: zero-touch onboarding is strategic for commercial
.mule/for-giulio.md:272:drafted with this investigation as its first item (issue #78 has the full

$ gh issue view 93 --json title,state,stateReason,labels,createdAt,updatedAt,comments --jq '{title:.title, state:.state, stateReason:.stateReason, labels:[.labels[].name], createdAt:.createdAt, lastComment:.comments[-1]}' 2>&1 | head -60
{"createdAt":"2026-09-04T19:44:12Z","labels":["mule-review","upstream-parity"],"lastComment":{"author":{"login":"atsetilam"},"authorAssociation":"MEMBER","body":"The mule pushed `8c61268` to `mule/queue` for this.\n\nUnreviewed: the gates passed — it compiles, the tests pass, a new test was shown to fail without the change, and no frozen file was touched. That is not the same as the change being worth having. Nothing merges to `main` until someone reads it:\n\n    bash tools/mule.sh review\n\nLeft open on purpose. Whether this actually resolves the issue is your call, not mine.","createdAt":"2026-09-04T19:54:58Z","id":"IC_kwDORmfsJs8AAAABSo2DbA","includesCreatedEdit":false,"isMinimized":false,"minimizedReason":"","reactionGroups":[],"url":"https://github.com/astrate-platform/astrate/issues/93#issuecomment-5545755500","viewerDidAuthor":true},"state":"OPEN","stateReason":"","title":"In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged"}

→ Read .mule/for-giulio.md [limit=50]

$ gh issue list --state open --limit 40 --json number,updatedAt --template '{{range .}}#{{.number}} {{.updatedAt}}{{"\n"}}{{end}}'
#101 2026-09-12T10:58:09Z
#100 2026-09-11T10:55:09Z
#99 2026-09-10T11:08:07Z
#98 2026-09-09T11:05:08Z
#97 2026-09-08T11:19:07Z
#96 2026-09-07T11:15:07Z
#95 2026-09-06T11:17:07Z
#94 2026-09-05T11:00:06Z
#93 2026-09-04T19:54:58Z
#92 2026-09-04T19:52:21Z
#78 2026-08-23T16:16:08Z
#1 2026-03-13T21:55:11Z

→ Read .mule/for-giulio.md [offset=15, limit=4]
Nothing new is proposable — same set as prior runs, plus one more alarm day. Appending today's triage note to the escalation channel.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -13,8 +13,10 @@
 - **Milestone recipe run, 2026-09-12: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) again shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing since the 2026-09-11 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked per Giulio 2026-09-04 on stable v1.4.0), **#93** aclhook comment rewrite (mule-review; `8c61268` on `mule/queue`, **not yet on `origin/main`**), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
 
 - **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
 
+- **github-issues triage run, 2026-09-12: still nothing proposable — the mule-alarm pile is now 8 straight days (#94–#101).** Twelve open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-11 run: **#101**, today's alarm (created 10:58Z, "nothing has landed in 14h"); **#100** has expired into the pile (superseded by #101). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-11 one: close #94–#100** — each was a one-day event, superseded by next day's alarm, never actionable; leave #101 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
+
 - **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pile is now 7 straight days and worth a look.** Twelve open issues. Still no machine-checkable fix candidates, so no task lines: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-10 run: **#100**, today's alarm (created 10:55Z, "nothing has landed in 15h"). #99–#94 are the same `mule-alarm` (one per day 2026-09-05→09-10, ~11:00Z each, zero comments) — telemetry, not code issues, so never proposable. But seven consecutive alarms is past the "low-activity window" wording earlier runs used: the mule has landed no commit since ~2026-09-04/05, and `main`'s queue copy shows the top line as `- [ ]` while the landed work below is all `[!] BLOCKED` (`wrote nothing` / `tests failed`) — though `mule/queue` is the authoritative copy. **Proposal: close #94–#99** — each was a one-day event, superseded by next day's alarm, never actionable; leave #100 (live today) to self-expire. If the idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
 
 - **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
 

```
Done: Triage of github-issues recipe — nothing proposable; no todo.md task lines; escalated today's alarm pile state to for-giulio.md
Files: .mule/for-giulio.md
Verified: gh issue list --state open --limit 40 -> 12 open; no gh issue create needed (recipe asks only for triage)
Unsure: nothing
Follow-ups: none
```

All 12 open issues are either already-handled (#93 pushed commit `8c61268` awaiting review), gated decisions (#92 parked on stable v1.4.0, #78 milestone-4.0 FDO), the untouchable #1, or `mule-alarm` telemetry (#94–#101). #101 is today's new live alarm; #100 just expired into the pile. No machine-checkable fix candidates exist, so per the recipe nothing was appended to todo.md. I extended the prior close-#94–#99 proposal to close #94–#100, leaving #101 to self-expire. No git, no GitHub state changes.
