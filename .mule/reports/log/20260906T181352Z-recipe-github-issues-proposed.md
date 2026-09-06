slug: recipe-github-issues
verdict: proposed
at:  8206be7
ran: 2026-09-06T18:13:52Z on DietPi in 97s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md
→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ gh issue list --state open --limit 40 \
  --json number,title,labels,updatedAt,comments \
  --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
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

All five open issues are already covered by existing mechanisms:

- **#95, #94** — mule-alarm (idle queue), not code issues
- **#93** — pushed commit `8c61268` awaiting review, its own recipe path
- **#92** — parked on stable v1.4.0, already escalated
- **#78** — milestone-4.0 FDO investigation, already escalated
- **#1** — standing instruction: never touch

No new task lines proposed. I'll append the for-giulio entry noting this run.
→ Read .mule/for-giulio.md [limit=10, offset=43]
The existing for-giulio entry from 2026-09-05 is nearly identical to today's findings — only #95 is new (another mule-alarm, same category as #94). I'll update the date and note on the new alarm rather than add a duplicate.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -39,15 +39,15 @@
   `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` are consumed in `internal/config/config.go`
   (lines 270-271 and 288); neither appears in `docs/site/configuration-reference.md`.
   Add them or confirm they are meant to be internal-only.
 
-- **github-issues triage run, 2026-09-05: nothing proposable, nothing stale.** All five open
-  issues were already covered by another mechanism. #94 is a mule-alarm (nothing landing in
-  15h), not a code issue — the idle window is the 2026-09-04 queue rebuild, expected. #93 has
-  a pushed commit (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable
-  v1.4.0 per the waiting-on row. #78 is the milestone-4.0 FDO design/investigation, already
-  escalated below. #1 untouched per standing instruction. No task lines proposed, no
-  for-giulio close suggestions.
+- **github-issues triage run, 2026-09-06: nothing proposable, nothing stale.** All six open
+  issues were already covered by another mechanism. #95 and #94 are mule-alarms (idle queue),
+  not code issues — expected during low-activity windows. #93 has a pushed commit (`8c61268`)
+  awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the waiting-on
+  row. #78 is the milestone-4.0 FDO design/investigation, already escalated below. #1
+  untouched per standing instruction. No task lines proposed, no for-giulio close suggestions.
+  (Identical to 2026-09-05 run — issue list unchanged except new alarm #95.)
 
 ---
 
 - **keyAgreement: the parking condition from #51 has fired — implement now, or wait for a

Done: github-issues triage run — nothing proposable, nothing stale. All six open issues already covered by existing mechanisms; no new task lines, no issues filed.
Files: `.mule/for-giulio.md` (updated existing entry)
Verified: `gh issue list` returned 6 open issues; all cross-checked against todo.md, for-giulio.md, and standing instructions.
Unsure: nothing
Follow-ups: none
