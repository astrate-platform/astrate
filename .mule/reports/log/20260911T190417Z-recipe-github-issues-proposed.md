slug: recipe-github-issues
verdict: proposed
at:  d2e91d9
ran: 2026-09-11T19:04:17Z on DietPi in 179s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md
Surveying open issues as the recipe directs.
→ Read .mule/todo.md

$ gh issue list --state open --limit 40 --json number,title,labels,updatedAt,comments --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
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

Issues #94–#100 look like auto-generated alarms. Checking their dates before deciding anything.
$ gh issue list --state open --limit 40 --json number,title,labels,createdAt,updatedAt --template '{{range .}}#{{.number}} {{.title}} | created={{.createdAt}} | updated={{.updatedAt}} | [{{range .labels}}{{.name}},{{end}}]\n{{end}}'
#100 mule: nothing has landed in 15h | created=2026-09-11T10:55:09Z | updated=2026-09-11T10:55:09Z | [mule-alarm,]\n#99 mule: nothing has landed in 14h | created=2026-09-10T11:08:07Z | updated=2026-09-10T11:08:07Z | [mule-alarm,]\n#98 mule: nothing has landed in 16h | created=2026-09-09T11:05:08Z | updated=2026-09-09T11:05:08Z | [mule-alarm,]\n#97 mule: nothing has landed in 15h | created=2026-09-08T11:19:07Z | updated=2026-09-08T11:19:07Z | [mule-alarm,]\n#96 mule: nothing has landed in 14h | created=2026-09-07T11:15:07Z | updated=2026-09-07T11:15:07Z | [mule-alarm,]\n#95 mule: nothing has landed in 15h | created=2026-09-06T11:17:07Z | updated=2026-09-06T11:17:07Z | [mule-alarm,]\n#94 mule: nothing has landed in 15h | created=2026-09-05T11:00:06Z | updated=2026-09-05T11:00:06Z | [mule-alarm,]\n#93 In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged | created=2026-09-04T19:44:12Z | updated=2026-09-04T19:54:58Z | [mule-review,upstream-parity,]\n#92 keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds | created=2026-09-04T19:43:54Z | updated=2026-09-04T19:52:21Z | [enhancement,upstream-parity,upstream-experimental,]\n#78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate) | created=2026-08-22T02:24:44Z | updated=2026-08-23T16:16:08Z | [enhancement,milestone-4.0,upstream-parity,]\n#1 Provide an Open Source IoT Platform unironically, unironically | created=2026-03-13T21:54:29Z | updated=2026-03-13T21:55:11Z | [wontfix,]\n

Alarms #94–#100 are daily auto-generated noise (Sep 5–11, same text every day). Checking current for-giulio.md state before appending.
→ Read .mule/for-giulio.md
Reading the state of things: last two runs already triaged these same issues. Let me compare today's list against those notes.

The list is unchanged in kind — but #100 is a **new** alarm created today (10:55Z), and the alarm pile now spans 7 straight days (Sep 5–11), whereas the 2026-09-08/09-10 runs saw only 1–2 alarms and called them "expected during low-activity windows." Seven days without a landing is a different signal. No issue is a proposable code task: #93 is already implemented (`mule-review`, `8c61268`), #92 is parked via the waiting-on mechanism, #78 is milestone-governed and already escalated, #1 is forbidden. So: no todo.md lines, and one for-giulio entry.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,10 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **github-issues triage run, 2026-09-11: nothing new proposable — the mule-alarm pile is now 7 straight days and worth a look.** Twelve open issues. Still no machine-checkable fix candidates, so no task lines: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-10 run: **#100**, today's alarm (created 10:55Z, "nothing has landed in 15h"). #99–#94 are the same `mule-alarm` (one per day 2026-09-05→09-10, ~11:00Z each, zero comments) — telemetry, not code issues, so never proposable. But seven consecutive alarms is past the "low-activity window" wording earlier runs used: the mule has landed no commit since ~2026-09-04/05, and `main`'s queue copy shows the top line as `- [ ]` while the landed work below is all `[!] BLOCKED` (`wrote nothing` / `tests failed`) — though `mule/queue` is the authoritative copy. **Proposal: close #94–#99** — each was a one-day event, superseded by next day's alarm, never actionable; leave #100 (live today) to self-expire. If the idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
+
 - **Milestone recipe run, 2026-09-11: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable `v1.3.3` (2026-08-07), newest overall `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-10 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
 
 - **Milestone recipe run, 2026-09-10: v3.0 verified unchanged — still gated on upstream, no action.** Latest release sweep (`gh api repos/astarte-platform/astarte/releases`) confirms upstream still has **no stable v1.4.0**: newest stable is `v1.3.3` (2026-08-07), newest pre-release is `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-09 run. `milestone-3.0` label remains empty (this milestone's ledger lives under `upstream-parity`): of #47–#89 all are closed except the already-escalated **#92 keyAgreement** (1.4 experimental decision, gated on a stable v1.4.0), **#93** aclhook comment rewrite (mule-review, pushed commit `8c61268`), **#78 FDO** (milestone-4.0). Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on upstream **v1.4 final** for their reconcile-vs-promote call — the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per **#90's frozen decision** (2026-08-23, #90 closed with it recorded): the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the 1.4-experimental rows are reconciled. No new gaps found, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
 

Done: triaged the open issues — nothing proposable as a task, appended one escalation to for-giulio.md about the 7-day mule-alarm pile.

Files: .mule/for-giulio.md
Verified: gh issue list (open, limit 40, with timestamps) -> pass
Unsure: whether the idle streak on the Pi is intentional — the alarm pile says "no landing since ~2026-09-04/05", but `mule/queue` (authoritative) may disagree with `main`'s todo copy; flagged for Giulio.
Follow-ups: none (no todo.md task lines — #93 already implemented, #92/#78 are decision/parked items, #1 forbidden; the 7 `mule-alarm` issues are telemetry, escalated only)
