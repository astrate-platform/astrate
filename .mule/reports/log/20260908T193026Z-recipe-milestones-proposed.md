slug: recipe-milestones
verdict: proposed
at:  HEAD
ran: 2026-09-08T19:30:26Z on DietPi

> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
→ Read .mule/milestones.md
✗ Read .mule/research/upstream-parity-2026-08-22.md failed
Error: File not found: /root/astrate-mule/.mule/research/upstream-parity-2026-08-22.md

The v3.0 milestone is the current target (first section not marked DONE). Prior runs
(2026-09-05, 09-06, 09-07T17:19Z, 09-07T20:47Z) all concluded the milestone is gated on
externals. Verify nothing has moved since the last milestone run (2026-09-07T20:47Z).

→ Read docs/UPSTREAM-EXPERIMENTAL.md
→ Read .mule/for-giulio.md

$ gh issue list --label milestone-3.0 --state all --limit 50 2>&1
(empty)

$ gh issue list --label upstream-parity --state open --limit 50 2>&1
93 OPEN mule-review, upstream-parity (aclhook keyAgreement comment rewrite, commit 8c61268 awaiting review)
92 OPEN enhancement, upstream-parity, upstream-experimental (keyAgreement parking decision, gated on stable v1.4.0)
78 OPEN enhancement, milestone-4.0, upstream-parity (FDO device onboarding, milestone-4.0 candidate)

$ gh issue list --label upstream-parity --state open --limit 60 | wc -l; gh issue list --label upstream-parity --state closed --limit 60 | wc -l
3
43

$ gh issue list --state open --limit 100
97 OPEN mule-alarm "nothing has landed in 15h" (2026-09-08, idle-queue noise)
96/95/94 OPEN mule-alarm (idle-queue noise, already covered by github-issues triage)
93/92/78 OPEN — as above
1 OPEN wontfix — Issue #1, never touched per standing instruction

Unchanged from the 2026-09-07 runs. Open `upstream-parity` issues are exactly the three
already escalated; nothing new opened or closed since.

$ gh api repos/astarte-platform/astarte/releases --paginate -q '.[] | select(.prerelease==false) | .tag_name' | head -5
v1.3.3
v1.3.2
v1.3.1
v1.3.0
v1.2.2

$ gh api repos/astarte-platform/astarte/releases --paginate -q '.[] | select(.prerelease==true) | .tag_name' | head -3
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3

Upstream v1.4.0 is still rc-only (rc.5); the newest stable remains v1.3.3 (2026-08-07).
No new upstream release since the 2026-08-22 investigation or the 09-07 run, so the
milestone's two gates have not moved:

1. **UPSTREAM-EXPERIMENTAL.md rows #67/#68** ("required/encrypted mapping fields",
   "async_operation=false") are 1.4-experimental and reconcile only when upstream 1.4
   final ships — it hasn't. No promoted → keep, no deprecated → deprecate, no reconcile due.
2. **APICompatVersion bump** — per #90's frozen decision (2026-08-23), stays 1.2.2 until
   the full 1.3 surface is complete and every UPSTREAM-EXPERIMENTAL row at that level is
   reconciled; the 1.4-experimental rows can't reconcile until 1.4 final, so the bump is
   still not due. #90 itself is the parked issue covering this; re-filing it would
   duplicate (recipe: never file the same gap twice).

Source doc `.mule/research/upstream-parity-2026-08-22.md` remains absent from `main` and
`origin/mule/research`; the investigation's output (#47–#89) is intact and reachable, so
milestone work is unaffected. Already escalated 2026-09-06; not re-escalated here.

**Proposed: nothing.**
- No machine-checkable gap → no `gh issue create` (step-4 bucket unmet; all implementable
  #47–#89 are closed).
- No `.mule/todo.md` lines appended (no open implementable work on v3.0; todo.md has no
  milestone-3.0/upstream-parity lines).
- No new `.mule/for-giulio.md` entry — the 2026-09-06 escalation already states this exact
  state and is still accurate; re-adding would duplicate a live queue line.
- The recipe's step-5 "milestone looks complete, verify and cut the tag" line is
  deliberately **not** proposed: open items #92/#93 in the milestone's own backlog and the
  upstream v1.4.0 gate are real, live decisions that are Giulio's, so the milestone is not
  complete and must not be cut.

Milestone v3.0 state: 1.3 surface delivered (retention ceiling #72, alias/group #52–#54
#59, query formats #55 #56, wire capabilities #47–#49, validation/error-code fidelity
#57 #61 #62 #79, per-service version #77, housekeeping #73–#76 — all closed); the rest
waits on Giulio to answer #92 when v1.4.0 goes stable, and the #90 final-phase
APICompatVersion audit. This run filed no issues and proposes no task lines.

```
Done: v3.0 milestone recipe — state unchanged since 2026-09-07T20:47Z; nothing to file, propose, or escalate (milestone gated on upstream v1.4.0 final + Giulio's #92 decision)
Files: .mule/reports/log/20260908T193026Z-recipe-milestones-proposed.md
Verified: gh issue list --label milestone-3.0 -> empty; --label upstream-parity -> 3 open/43 closed; gh api astarte releases -> stable v1.3.3, prerelease v1.4.0-rc.5; git ls-tree origin/mule/research parity doc absent -> confirmed
Unsure: nothing
Follow-ups: none
```