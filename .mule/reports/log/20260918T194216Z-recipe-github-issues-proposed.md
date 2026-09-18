slug: recipe-github-issues
verdict: proposed
at:  79439e6
ran: 2026-09-18T19:42:16Z on DietPi in 119s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md

$ gh issue list --state open --limit 40 \
  --json number,title,labels,updatedAt,comments \
  --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
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

✱ Grep "nothing has landed" 100 matches

$ gh issue list --state open --limit 40 --json number,title,createdAt,updatedAt,labels --template '{{range .}}#{{.number}} | created {{.createdAt}} | upd {{.updatedAt}} | [{{range .labels}}{{.name}} {{end}}]{{"\n"}}{{end}}'
#106 | created 2026-09-18T10:54:09Z | upd 2026-09-18T10:54:09Z | [mule-alarm ]
#105 | created 2026-09-17T10:51:07Z | upd 2026-09-17T10:51:07Z | [mule-alarm ]
#104 | created 2026-09-16T11:12:14Z | upd 2026-09-16T11:12:14Z | [mule-alarm ]
#103 | created 2026-09-14T11:27:08Z | upd 2026-09-14T11:27:08Z | [mule-alarm ]
#102 | created 2026-09-13T11:09:08Z | upd 2026-09-13T11:09:08Z | [mule-alarm ]
#101 | created 2026-09-12T10:58:09Z | upd 2026-09-12T10:58:09Z | [mule-alarm ]
#100 | created 2026-09-11T10:55:09Z | upd 2026-09-11T10:55:09Z | [mule-alarm ]
#99 | created 2026-09-10T11:08:07Z | upd 2026-09-10T11:08:07Z | [mule-alarm ]
#98 | created 2026-09-09T11:05:08Z | upd 2026-09-09T11:05:08Z | [mule-alarm ]
#97 | created 2026-09-08T11:19:07Z | upd 2026-09-08T11:19:07Z | [mule-alarm ]
#96 | created 2026-09-07T11:15:07Z | upd 2026-09-07T11:15:07Z | [mule-alarm ]
#95 | created 2026-09-06T11:17:07Z | upd 2026-09-06T11:17:07Z | [mule-alarm ]
#94 | created 2026-09-05T11:00:06Z | upd 2026-09-05T11:00:06Z | [mule-alarm ]
#93 | created 2026-09-04T19:44:12Z | upd 2026-09-04T19:54:58Z | [mule-review upstream-parity ]
#92 | created 2026-09-04T19:43:54Z | upd 2026-09-04T19:52:21Z | [enhancement upstream-parity upstream-experimental ]
#78 | created 2026-08-22T02:24:44Z | upd 2026-08-23T16:16:08Z | [enhancement milestone-4.0 upstream-parity ]
#1 | created 2026-03-13T21:54:29Z | upd 2026-03-13T21:55:11Z | [wontfix ]

→ Read .mule/todo.md
→ Read .mule/for-giulio.md
→ Read .mule/reports/log/20260917T190323Z-recipe-github-issues-proposed.md

$ ls -la .mule/reports/log/ | tail -20
-rw-r--r-- 1 root root   3405 Sep 17 20:22 20260917T182237Z-recipe-hygiene-timeout.md
-rw-r--r-- 1 root root  97405 Sep 17 20:43 20260917T184319Z-recipe-milestones-timeout.md
-rw-r--r-- 1 root root  74737 Sep 17 21:03 20260917T190323Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   2026 Sep 17 21:06 20260917T190650Z-recipe-astarte-upstream-noop.md
-rw-r--r-- 1 root root  47648 Sep 17 21:07 20260917T190703Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   6840 Sep 17 21:12 20260917T191223Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root   5235 Sep 17 21:22 20260917T192217Z-httpx-cors-vary-origin-passthrough-done.md
-rw-r--r-- 1 root root  43685 Sep 17 21:49 20260917T194954Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root   9070 Sep 17 22:07 20260917T200755Z-docs-sync-rm-policies-delete-422-done.md
-rw-r--r-- 1 root root  11351 Sep 17 22:21 20260917T202137Z-docs-sync-rm-triggers-422-nested-envelope-done.md
-rw-r--r-- 1 root root   4566 Sep 17 23:01 20260917T210114Z-recipe-hygiene-timeout.md
-rw-r--r-- 1 root root  40531 Sep 18 13:00 20260918T110039Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  11055 Sep 18 13:34 20260918T113442Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   3271 Sep 18 13:40 20260918T113930Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  64748 Sep 18 13:41 20260918T114100Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  10751 Sep 18 18:53 20260918T165351Z-compat-note-v134-done.md
-rw-r--r-- 1 root root  35381 Sep 18 19:26 20260918T172636Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root  27631 Sep 18 20:01 20260918T180119Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root  55277 Sep 18 20:41 20260918T184115Z-recipe-hygiene-timeout.md
-rw-r--r-- 1 root root  14495 Sep 18 21:36 20260918T193641Z-recipe-milestones-proposed.md

$ date -u +"%Y-%m-%dT%H:%M:%SZ" && git log --oneline -5
2026-09-18T19:41:42Z
79439e6 mule: refill the queue
b3c4bad mule: recipe milestones ran (176s)
eb95cdd mule: recipe hygiene timed out
4c4cf5b mule: recipe docs-sync ran (421s)
f6dc460 mule: refill the queue

→ Read .mule/reports/log/20260918T113442Z-recipe-github-issues-proposed.md
The recipe already ran today: **`.mule/reports/log/20260918T113442Z-recipe-github-issues-proposed.md`** shows the identical triage completed at 11:34Z on the Pi, and its output is already at the top of `.mule/for-giulio.md:25` — the "close #94–#105, leave #106 to self-expire" proposal.

I re-ran the survey (current time 19:41Z): nothing has changed since that run — same 16 open issues, no new ones; #106 is still today's live alarm (created 09-18 10:54Z), #105–#94 still the expired pile, #93/#92/#78/#1 unchanged. So this run is a duplicate: writing the same for-giulio line again would file the proposal twice. No todo.md task lines either — zero machine-checkable code-fix candidates (all mule-alarms, plus the already-mechanism'd set), and the recipe asks for no `gh issue create`, so nothing is filed or created this run.

```
Done: github-issues triage re-run — issue set identical to today's 11:34Z run; proposal already in for-giulio.md (close #94–#105); no duplicate appended, no machine-checkable candidates so no todo.md lines
Files: (none changed — deliverable already present)
Verified: gh issue list --state open --limit 40 -> pass (16 open, same set)
Unsure: nothing
Follow-ups: none
```
