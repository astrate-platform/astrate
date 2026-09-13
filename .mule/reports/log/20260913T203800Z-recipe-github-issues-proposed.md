slug: recipe-github-issues
verdict: proposed
at:  c8218e5
ran: 2026-09-13T20:38:00Z on DietPi

> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md

$ gh issue list --state open --limit 40 \
  --json number,title,labels,updatedAt,comments \
  --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
#102 mule: nothing has landed in 14h [mule-alarm ] []c
#101 mule: nothing has landed in 14h [mule-alarm ] []c
#100 mule: nothing has landed in 15h [mule-alarm ] []c
#99 mule: nothing has landed in 14h [mule-alarm ] []c
#98 mule: nothing has landed in 16h [mule-alarm ] []c
#97 mule: nothing has landed in 15h [mule-alarm ] []c
#96 mule: nothing has landed in 14h [mule-alarm ] []c
#95 mule: nothing has landed in 15h [mule-alarm ] []c
#94 mule: nothing has landed in 15h [mule-alarm ] []c
#93 In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged [mule-review upstream-parity ] [map]c
#92 keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds [enhancement upstream-parity upstream-experimental ] [map]c
#78 FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate) [enhancement milestone-4.0 upstream-parity ] [map]c
#1 Provide an Open Source IoT Platform unironically, unironically [wontfix ] []c

Identical to the 2026-09-13T17:05:50Z run of this recipe — 13 open issues,
the same set, no new issue, no comments since. Triage is unchanged:

- **#94–#102** are `mule-alarm` telemetry (idle-queue alarms, one per day),
  not code issues — never proposable. #102 (created today 11:09Z) is live;
  the rest expired into the pile. Stale-issue proposal → for-giulio.md.
- **#93** aclhook comment (mule-review) has a pushed commit `8c61268`
  awaiting review — its own recipe path, not a fresh proposal candidate.
- **#92** keyAgreement is a design decision parked on a stable upstream
  v1.4.0 per the waiting-on row — an escalation, not a machine-checkable fix.
- **#78** FDO is the milestone-4.0 investigation, already escalated.
- **#1** untouched per standing instruction.

Still zero machine-checkable code-fix candidates → no new lines appended to
.mule/todo.md. The closest proposal for #94–#101 and the leave-#102-open
call are already the newest entry in .mule/for-giulio.md (2026-09-13 entry,
written by the 17:05Z run); duplicating them would be noise, so nothing was
appended there either. This log is the whole record of this run.