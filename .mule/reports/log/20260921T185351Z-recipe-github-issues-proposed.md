slug: recipe-github-issues
verdict: proposed
at:  65d7071
ran: 2026-09-21T18:53:51Z on DietPi in 154s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md

$ ls -la /root/astrate-mule/.mule/
total 5916
drwxr-xr-x  7 root root    4096 Sep 21 20:20 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      13 Sep 21 20:51 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 21 19:35 .heartbeat
-rw-r--r--  1 root root       0 Sep 21 20:51 .last-output
-rw-r--r--  1 root root       0 Sep 21 20:51 .lock
-rw-r--r--  1 root root      14 Sep 21 20:51 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 5830305 Sep 21 20:51 cron.log
drwxr-xr-x  2 root root    4096 Sep 19 12:52 failed
-rw-r--r--  1 root root   51848 Sep 21 20:37 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    8227 Sep 21 19:35 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 21 13:32 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   71911 Sep 21 19:35 todo.md

→ Read .mule/todo.md

$ gh issue list --state open --limit 40 \
  --json number,title,labels,updatedAt,comments \
  --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
#108 mule: nothing has landed in 15h [mule-alarm ] []c
#107 mule: nothing has landed in 14h [mule-alarm ] []c
#106 mule: nothing has landed in 14h [mule-alarm ] []c
#105 mule: nothing has landed in 14h [mule-alarm ] []c
#104 mule: nothing has landed in 14h [mule-alarm ] []c
#103 mule: nothing has landed in 16h [mule-alarm ] []c
#102 mule: nothing has landed in 14h [mule-alarm ] []c
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

Let me check `for-giulio.md` for format and prior mentions of these issues.
✱ Grep "mule-alarm|nothing has landed|#95|#96|#108" in .mule · 100 matches
→ Read .mule/for-giulio.md [offset=1, limit=60]

$ date -u '+%Y-%m-%d %H:%MZ'
2026-09-21 18:52Z

The issue set is unchanged except the mule-alarm pile grew. Let me get the current time for the report stamp.
$ gh issue list --state open --limit 40 --json number,title,createdAt --jq '.[] | "#\(.number) \(.createdAt) \(.title)"'
#108 2026-09-21T10:54:46Z mule: nothing has landed in 15h
#107 2026-09-19T10:49:12Z mule: nothing has landed in 14h
#106 2026-09-18T10:54:09Z mule: nothing has landed in 14h
#105 2026-09-17T10:51:07Z mule: nothing has landed in 14h
#104 2026-09-16T11:12:14Z mule: nothing has landed in 14h
#103 2026-09-14T11:27:08Z mule: nothing has landed in 16h
#102 2026-09-13T11:09:08Z mule: nothing has landed in 14h
#101 2026-09-12T10:58:09Z mule: nothing has landed in 14h
#100 2026-09-11T10:55:09Z mule: nothing has landed in 15h
#99 2026-09-10T11:08:07Z mule: nothing has landed in 14h
#98 2026-09-09T11:05:08Z mule: nothing has landed in 16h
#97 2026-09-08T11:19:07Z mule: nothing has landed in 15h
#96 2026-09-07T11:15:07Z mule: nothing has landed in 14h
#95 2026-09-06T11:17:07Z mule: nothing has landed in 15h
#94 2026-09-05T11:00:06Z mule: nothing has landed in 15h
#93 2026-09-04T19:44:12Z In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
#92 2026-09-04T19:43:54Z keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds
#78 2026-08-22T02:24:44Z FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)
#1 2026-03-13T21:54:29Z Provide an Open Source IoT Platform unironically, unironically

The set is unchanged except alarm **#108** is new today (superseding #107). No machine-checkable code-fix candidates — same verdict as every prior run. The recipe's deliverable is the escalation line, so I append one to `for-giulio.md`, newest at top.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **github-issues triage run, 2026-09-21: still nothing proposable — the mule-alarm pile is now 16 days (#94–#108).** Eighteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-20 evening run: **#108**, today's alarm (created 09-21 10:54Z, "nothing has landed in 15h"); **#107** has expired into the pile (superseded by #108). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-20 one: close #94–#107** — each was a one-day event, superseded by next day's alarm, never actionable; leave #108 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
+
+---
+
 - **Milestone recipe run, 2026-09-21: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-20 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
 
 ---
 

→ Read .mule/for-giulio.md [limit=24]
Done. Surveyed the 18 open issues with the recipe command; nothing changed except the alarm pile grew.

Findings:
- **#95–#108** (14 issues): daily `mule-alarm` telemetry — new **#108** today (superseded #107), never actionable code.
- **#93** aclhook comment: `mule-review`, `8c61268` already pushed, its own recipe path.
- **#92** keyAgreement: parked on stable v1.4.0 per the waiting-on row; escalation already wired.
- **#78** FDO: milestone-4.0 design/investigation, already escalated.
- **#1**: untouched per standing instruction.

Zero machine-checkable code-fix candidates, so no `.mule/todo.md` task lines. Appended one escalation line to `.mule/for-giulio.md` (newest at top) extending the closure proposal to **close #94–#107**, leave #108 to self-expire. No git touched, no issues created/closed/commented.

```
Done: Triage of open GitHub issues; nothing proposable — extended the closure proposal to #94–#107
Files: .mule/for-giulio.md
Verified: gh issue list (state open, limit 40) -> pass; edit applied and verified -> pass
Unsure: nothing
Follow-ups: none
```
