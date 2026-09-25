slug: recipe-github-issues
verdict: proposed
at:  31f4cc5
ran: 2026-09-25T16:54:08Z on DietPi in 170s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/github-issues.md

$ gh repo view --json nameWithOwner,description 2>&1; echo "---"; ls .mule/ 2>&1
{"description":"","nameWithOwner":"astrate-platform/astrate"}
---
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

$ gh issue list --state open --limit 40 --json number,title,labels,updatedAt,comments --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}' 2>&1
#111 mule: nothing has landed in 16h [mule-alarm ] []c
#110 mule: nothing has landed in 14h [mule-alarm ] []c
#109 mule: nothing has landed in 14h [mule-alarm ] []c
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

→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ gh issue list --state open --limit 40 --json number,title,createdAt --template '{{range .}}#{{.number}} {{.createdAt}}{{"\n"}}{{end}}' 2>&1
#111 2026-09-25T11:01:07Z
#110 2026-09-24T10:54:08Z
#109 2026-09-22T10:54:47Z
#108 2026-09-21T10:54:46Z
#107 2026-09-19T10:49:12Z
#106 2026-09-18T10:54:09Z
#105 2026-09-17T10:51:07Z
#104 2026-09-16T11:12:14Z
#103 2026-09-14T11:27:08Z
#102 2026-09-13T11:09:08Z
#101 2026-09-12T10:58:09Z
#100 2026-09-11T10:55:09Z
#99 2026-09-10T11:08:07Z
#98 2026-09-09T11:05:08Z
#97 2026-09-08T11:19:07Z
#96 2026-09-07T11:15:07Z
#95 2026-09-06T11:17:07Z
#94 2026-09-05T11:00:06Z
#93 2026-09-04T19:44:12Z
#92 2026-09-04T19:43:54Z
#78 2026-08-22T02:24:44Z
#1 2026-03-13T21:54:29Z

$ gh issue view 111 --json number,title,createdAt,labels,body --template '{{.number}} | {{.createdAt}} | [{{range .labels}}{{.name}} {{end}}]{{"\n"}}BODY:{{.body}}' 2>&1; echo "=== open count ==="; gh issue list --state open --limit 100 --json number --jq 'length' 2>&1; echo "=== mule-labelled open ==="; gh issue list --state open --label mule --limit 50 --json number,title 2>&1
111 | 2026-09-25T11:01:07Z | [mule-alarm ]
BODY:The mule has not landed any work in 16 hours (threshold 8h).

It is probably still ticking and exiting 0 — that is how this failure looks. Worth checking,
on the Pi:

    systemctl list-timers mule.timer
    journalctl -u mule.service --since '-1 day' | tail -50
    cd /root/astrate-mule && bash tools/mule.sh status

Common causes, all of which have happened: the provider is refusing every run, the queue has
nothing runnable in it, the working tree is dirty so every tick aborts, or push has been
failing and the work is only on the Pi.

Filed automatically. Close it — the mule reopens a new one if the silence continues.=== open count ===
22
=== mule-labelled open ===
[]

$ wc -l .mule/todo.md .mule/for-giulio.md
220 .mule/todo.md
   360 .mule/for-giulio.md
   580 total

→ Read .mule/todo.md [offset=200]

$ rg -n '^- \[ \]' .mule/todo.md | wc -l; echo "--- open lines ---"; rg -n '^- \[ \]' .mule/todo.md -o -r '' 2>/dev/null | head -30
11
--- open lines ---
86:
88:
89:
101:
102:
103:
104:
105:
200:
203:
205:

Survey is complete. Let me check one more thing — the alarm mechanism itself, since it's the only thing that changed.
$ rg -n 'mule-alarm|nothing has landed|threshold' tools/ .mule/config/ 2>/dev/null | head -20
tools/mule.sh:57:MULE_ALARM_LABEL="${MULE_ALARM_LABEL:-mule-alarm}"
tools/mule.sh:751:  local body="The mule has not landed any work in ${age} hours (threshold ${max}h).
tools/mule.sh:766:    gh issue create --title "mule: nothing has landed in ${age}h" --body "$body" \

→ Read tools/mule.sh [offset=735, limit=40]

$ rg -n '\bbeat\b|\.alarmed|\.heartbeat' tools/mule.sh
322:  # cmd_refill once sliced beat() and check_pulse() out of the file and left their call sites
327:  for fn in beat check_pulse first_open code_id line_slug proof_gate gates issue_tasks \
540:            beat
578:    beat
737:beat() { date +%s > "$MULE/.heartbeat"; rm -f "$MULE/.alarmed"; }
743:  if [ ! -f "$MULE/.heartbeat" ]; then beat; return 0; fi
744:  last="$(cat "$MULE/.heartbeat" 2>/dev/null)"; [ -n "$last" ] || { beat; return 0; }
748:  [ -f "$MULE/.alarmed" ] && return 0
749:  : > "$MULE/.alarmed"

→ Read tools/mule.sh [limit=14, offset=532]
→ Read tools/mule.sh [limit=10, offset=573]

$ rg -n 'cmd_refill' tools/mule.sh | head; echo "---"; sed -n '/^cmd_refill()/,/^}/p' tools/mule.sh | rg -n 'beat|git .*commit|for-giulio' | head -20
322:  # cmd_refill once sliced beat() and check_pulse() out of the file and left their call sites
328:            issue_done refill_from_recipe cmd_next cmd_refill cmd_review cmd_budget \
642:    cmd_refill || note "the refill added nothing — idling until the next tick"
707:cmd_refill() {
976:  refill)    shift; cmd_refill "$@";;
---
14:  git -C "$REPO" diff --cached --quiet || git -C "$REPO" commit -q -m "mule: refill the queue"

→ Read tools/mule.sh [limit=40, offset=707]
Confirmed the mechanism. Now I have one genuinely new finding. Writing the escalation entry.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,12 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **github-issues triage run, 2026-09-25: still nothing proposable — the mule-alarm pile is 18 issues (#94–#111), and two corrections to what earlier runs inferred about it.** Twenty-two open issues, same set as every run since 2026-09-05: **#111–#94** are the daily `mule-alarm` noise (zero comments, telemetry not code), **#93** aclhook comment (mule-review, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable upstream v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-24 run: **#111** (today, 09-25 11:01:07Z, "nothing has landed in 16h"). Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines added**. **Proposal, extending the 2026-09-24 one: close #94–#110**; leave #111 (live today) to self-expire. Two corrections, both from reading `tools/mule.sh` rather than the issue titles: **(1) the "missed days" in earlier entries (09-15, 09-20, 09-23) are evidence of activity, not its absence** — I previously called them timer gaps; the alarm is gated on a `.alarmed` sentinel that `beat` deletes on every land (`mule.sh:737,748-749`), so a day with no new issue is a day something landed, not a day the timer slept. **(2) the alarm body's closing line is wrong**: "Close it — the mule reopens a new one if the silence continues" (`mule.sh:764`) implies issue state drives the alarm; it does not. Nothing reads or reconciles the 18 open issues, so the pile can only grow and closing the old ones is safe but inert — and `cmd_refill` commits `mule: refill the queue` (mule.sh:719-720) **without calling `beat`**, so a doc-only refill never counts as landing work. The titles' ages also imply a beat roughly every day around 19:00–21:00Z (#111 filed 09-25 11:01Z at 16h → last beat ≈ 09-24 19:01Z; #110 filed 09-24 10:54Z at 15h → ≈ 09-23 19:54Z) while the ~11:00Z daily check has never once found under the 8h threshold — **what calls `beat` at that time is unverified**, and one command on the Pi settles it (`grep -E 'landed|checked:' /root/astrate-mule/.mule/log | tail -20`). The queue has still landed nothing since ~2026-09-04/05; if that idle streak is not intentional, the wedge is on the Pi and a queue review is the fix.
+
+---
+
 - **Milestone recipe run, 2026-09-24: v3.0 verified unchanged — still gated on upstream, no action.** Release sweep (`gh api repos/astarte-platform/astarte/releases`) still shows **no stable v1.4.0**: newest stable remains `v1.3.4` (2026-09-18, maintenance-only, no wire/API surface change), newest overall remains `v1.4.0-rc.5` (2026-08-20) — nothing shipped since the 2026-09-23 run. `milestone-3.0` label still empty (this milestone's ledger lives under `upstream-parity`): open items are exactly the already-escalated **#92 keyAgreement** (1.4 experimental decision, parked on stable v1.4.0 per Giulio 2026-09-04), **#93** aclhook comment rewrite (mule-review, pushed `8c61268`, its own recipe path), **#78 FDO** (milestone-4.0, already escalated); the #47–#89 rest are all closed. Both UPSTREAM-EXPERIMENTAL rows (#67 required/encrypted flags, #68 async_operation) still wait on **v1.4 final** for their reconcile-vs-promote call; the 1.3 surface itself is fully delivered. `APICompatVersion` stays 1.2.2 per #90's frozen decision — the bump is v3.0's final-phase item, due once v1.4.0 goes stable and the experimental rows are reconciled. No new gaps, no issues filed, no task lines proposed; deliberately NOT proposing "complete, cut the tag" — the milestone targets 1.3/**1.4** and the 1.4 gate is unresolved. Your call when a stable v1.4.0 ships: answer #92, run the final-phase bump, then cut the v3.0 tag.
 
 ---
 

→ Read .mule/for-giulio.md [limit=16]

$ git status --short; echo "=== diff stat ==="; git diff --stat; echo "=== todo.md unchanged? ==="; git diff --quiet .mule/todo.md && echo "todo.md: no changes"
M .mule/for-giulio.md
=== diff stat ===
 .mule/for-giulio.md | 4 ++++
 1 file changed, 4 insertions(+)
=== todo.md unchanged? ===
todo.md: no changes

Recipe carried out. No code touched, no git commands run, nothing created on GitHub.

**Survey result: zero proposable issues.** All 22 open issues are accounted for:

| Issue | Why not proposed |
|---|---|
| #94–#111 (18) | `mule-alarm` telemetry, zero comments, no code — close proposal instead |
| #93 | mule-review, `8c61268` already pushed, awaiting human review (its own recipe path) |
| #92 | parked by Giulio on a stable upstream v1.4.0, wired to `waiting-on.md` |
| #78 | milestone-4.0 FDO, design decision, already escalated |
| #1 | never touched per standing instruction |

**Deliverable:** one entry at the top of `.mule/for-giulio.md` (4 lines added, nothing else changed). No `.mule/todo.md` lines, because the recipe only permits issue-derived lines and no issue qualified. The recipe did *not* ask for `gh issue create`, so I ran no `gh` write commands at all.

The one new thing I found was worth the run: I read `tools/mule.sh` instead of trusting the alarm titles, and **two beliefs held by the last ~20 triage entries are wrong**.

1. The "missed days" (09-15, 09-20, 09-23) were read as timer gaps. They're the opposite — `check_pulse` is gated on a `.alarmed` sentinel that `beat` deletes on every land (`tools/mule.sh:737,748-749`), so a day with *no* new alarm is a day something landed. Earlier runs drew the inference backwards.
2. The alarm body tells you to close it because "the mule reopens a new one" — but no code path reads issue state, so the pile can only ever grow. Closing the old 17 is safe but inert, and `cmd_refill` commits its refill without calling `beat` (mule.sh:719-720), so doc-only refills never count as work.

I also flagged, explicitly as unverified, that the alarm titles imply a `beat` call around 19:00–21:00Z daily while the 11:00Z check has never once found under threshold — I did not guess the cause, and left the one Pi command that settles it in the entry.

One thing worth your attention, outside the recipe's scope: the queue has landed nothing since ~2026-09-04/05, and 9 of its 11 open lines are `[legion]`, which the Pi skips. If the Legion is meant to be reachable, that idle streak is the thing to look at.

```
Done: Triaged all 22 open issues for astrate-platform/astrate; none proposable, escalated the alarm pile with two corrections to prior runs' reasoning
Files: .mule/for-giulio.md
Verified: gh issue list --state open --limit 40 (22 open, 0 labelled `mule`); gh issue view 111 (alarm body); read tools/mule.sh:707-774 to confirm the beat/check_pulse/.alarmed mechanism; git status/diff --stat -> 1 file, +4 lines, todo.md untouched
Unsure: what calls `beat` at ~19:00-21:00Z daily — the 14-16h ages in the alarm titles imply it, but I did not verify it and left the Pi command that would
Follow-ups: none appended to .mule/todo.md
```
