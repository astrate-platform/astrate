slug: recipe-milestones
verdict: proposed
at:  e7a05f6
ran: 2026-09-09T17:57:18Z on DietPi in 115s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
→ Read .mule/todo.md
Now I'll read the milestones file and the todo file to understand current state.
→ Read .mule/for-giulio.md
The first non-DONE milestone is **v3.0**. The for-giulio.md entry from 2026-09-06 says it has no open machine-checkable gaps left. Let me verify by checking existing issues and the experimental status.

$ gh issue list --label upstream-parity --state open --limit 50 2>&1
93	OPEN	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review, upstream-parity	2026-09-04T19:54:58Z
92	OPEN	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement, upstream-parity, upstream-experimental	2026-09-04T19:52:21Z
78	OPEN	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement, milestone-4.0, upstream-parity	2026-08-23T16:16:08Z

$ gh issue list --label milestone-3.0 --state all --limit 50 2>&1
(no output)

$ gh issue list --label milestone-3.0 --state open --limit 50 2>&1
(no output)

→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ gh issue list --label milestone-3.0 --state all --limit 50 --json number,title,state 2>&1
[]

All confirmed: v3.0 has no open machine-checkable gaps. The milestone-3.0 label is empty, the upstream-parity open issues are all escalation items already in for-giulio.md, and the UPSTREAM-EXPERIMENTAL rows (#67/#68) reconcile only when upstream v1.4 goes stable (still rc). This matches step 5 of the recipe exactly.

No task lines or issues to file — the deliverable is the escalation line in for-giulio.md. I won't append anything to todo.md since there is no new work to propose.
→ Read .mule/for-giulio.md [limit=10]
The existing for-giulio.md entry from 2026-09-06 already covers this exact situation. The state hasn't changed (still no machine-checkable gaps, #92/#78/#93 still the same open items, upstream still rc). I'll add a concise update entry.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,18 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Milestone recipe run, 2026-09-09: v3.0 unchanged since 2026-09-06 — still waiting on
+  upstream gates.** `milestone-3.0` label remains empty. The three open `upstream-parity`
+  issues are the same: #92 keyAgreement (upstream decision, gated on stable v1.4.0), #93
+  aclhook comment (pushed commit `8c61268`, mule-review), #78 FDO (milestone-4.0). Both
+  UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait
+  on upstream v1.4 final (currently rc.5). `APICompatVersion` stays 1.2.2 per #90's frozen
+  decision. No new gaps found, no issues filed, no task lines proposed.
+
+---
+
 - **github-issues triage run, 2026-09-08: nothing proposable, nothing stale.** All six open
   issues were already covered by another mechanism. #97/#96/#95/#94 are mule-alarms (idle
   queue), not code issues — expected during low-activity windows. #93 has a pushed commit
   (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable v1.4.0 per the

Done. Let me verify the edit landed correctly and that I didn't touch todo.md:
→ Read .mule/for-giulio.md [limit=25, offset=10]
Entry landed cleanly. todo.md is untouched.

```
Done: Milestone recipe run on v3.0 — no machine-checkable gaps remain, escalated to for-giulio.md
Files: .mule/for-giulio.md
Verified: gh issue list --label milestone-3.0 returns empty; UPSTREAM-EXPERIMENTAL rows #67/#68 still wait on upstream v1.4 final (rc.5)
Unsure: nothing
Follow-ups: none
```
