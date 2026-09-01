slug: recipe-milestones
verdict: proposed
at:  23ebc21
ran: 2026-09-01T20:51:35Z on DietPi in 90s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/milestones.md
Let me read the milestones file and check the todo state.
→ Read .mule/milestones.md
→ Read .mule/todo.md
Let me look at the milestones recipe. The first section not marked `DONE` is **v2.0 — astarte-flow feature parity**. Its scope says "not started. No `milestone-2.0` issues filed yet." — it's not marked "TBD" in the sense of "not yet decided" for scope; the v2.0 section has a decided reference and scope. Let me also check the `for-giulio.md` and check v3.0 is TBD but it's a later milestone.

Per the recipe step 1: find the **first** section not marked DONE. That's v2.0. v2.0's scope is not "TBD/not yet decided" — it's decided (astarte_flow feature parity). So I proceed to step 3 (investigate).

Wait — but step 2 says "If that milestone's scope says 'not yet decided' / 'TBD'". v2.0 does not say that. So I proceed to step 3.

Let me check for-giulio.md and set up my plan.
→ Read .mule/for-giulio.md

$ gh issue list --label milestone-2.0 --state all --limit 50 2>&1; echo "---"; gh issue list --label milestone-2.0 --state open 2>&1
43	CLOSED	flow-container-block: PoC → MVP (custom processing in containers)	enhancement, milestone-2.0	2026-07-29T18:54:52Z
42	CLOSED	flow-rehydrate-edge-cases: follow-ups after durable auto_restart MVP	enhancement, milestone-2.0	2026-07-29T20:39:19Z
41	CLOSED	flow-durable: persist flows + auto_restart rehydrate on process boot	enhancement, milestone-2.0	2026-07-29T18:54:49Z
40	CLOSED	flow-multi-instance: named flows + pipeline config (decided — design then implement)	enhancement, question, milestone-2.0	2026-07-29T18:54:47Z
39	CLOSED	flow-blocks-discovery: GET /flow/v1/{realm}/blocks for registered catalog types	enhancement, milestone-2.0	2026-07-29T18:54:44Z
37	CLOSED	flow-source-pump: poll AstarteSource and Stop on flow teardown	milestone-2.0	2026-07-29T14:13:12Z
27	CLOSED	flow-astarte-source: Astarte device events source block	mule-review, milestone-2.0	2026-07-29T07:45:22Z
26	CLOSED	flow-router: Stream-based message routing through the block graph	milestone-2.0	2026-07-29T07:33:30Z
25	CLOSED	flow-lifecycle: Flow start/stop lifecycle and runtime management	milestone-2.0	2026-07-29T07:33:28Z
24	CLOSED	flow-pipeline-store: Pipeline CRUD storage and graph validation	mule-review, milestone-2.0	2026-07-29T07:45:19Z
23	CLOSED	flow-models: FlowMessage wire format and Block/Pipeline/Flow data models	milestone-2.0	2026-07-29T07:33:26Z
---

The `for-giulio.md` file contains a previous milestones run note (2026-09-01, today) explaining that on `main`, v2.0 is already DONE. But I'm on my current branch. Let me check the git branch and state — the note says `mule/queue` is 120 commits behind `origin/main` and the local `.mule/milestones.md` is stale.

Let me understand my current situation: which branch am I on, and does the note's claim matter? The recipe says find the first section not marked DONE in `.mule/milestones.md`. The file I read shows v2.0 as "not started". But that could be stale.

Let me check git state.
$ git show origin/main:.mule/milestones.md 2>&1 | head -60
# Milestones

Release-tag-gated goals. Each section names a body of work that must be true by the time
that tag is cut. `.mule/recipes/milestones.md` reads this file to find the **current
target** — the first section below not marked `DONE` — and works toward it: investigating,
filing GitHub issues (`milestone-<tag>` label, alongside `mule`) with sub-issues where the
work splits, and escalating anything that needs a design decision to `.mule/for-giulio.md`.

**This file is Giulio's.** The recipe reads it and may propose edits via
`.mule/for-giulio.md`, but never edits it directly — same rule as `docs/COMPATIBILITY.md`.
Mark a milestone `DONE` yourself once the tag is actually cut.

Order matters: milestones are worked **in order**, lowest tag first. Do not start v3.0
investigation while v2.0 has open, un-escalated gaps — say so and stop instead.

---

## v2.0 — astarte-flow feature parity

Reference: astarte_flow (https://github.com/astarte-platform/astarte_flow) — upstream's
Elixir "Flow" component. It lets ingested data get piped through a graph of processing
blocks (native blocks and containerised ones) before it lands in storage, with pipelines
described as reusable, parametrised graphs.

Scope for this milestone: whatever set of Flow's *capabilities* Astrate needs to expose the
same wire-visible behaviour and operator-facing concepts (pipelines, blocks, native vs.
containerised blocks) — not a port of the Elixir implementation. See
`.mule/recipes/astarte-upstream.md`'s rule: port the idea, restated in Go, never the code.

Status: **DONE** (2026-07-29). Runtime, factory, catalog (incl. filter/map), process
wiring, and `/flow/v1` API are on `main`. **Parity audit** + **product decisions**
recorded 2026-07-29 (`docs/handoff/flow-parity-audit-2026-07-29.md`,
`docs/handoff/flow-v2-decisions-2026-07-29.md`).

**Design A + B landed on `main`** (commit `89145e6`, 2026-07-29): durable flows +
auto_restart (**#41 closed**); named multi-instance + config (**#40 closed**) —
migration `000009`, store, `${config.*}`, API, boot rehydrate. Container block
PoC→MVP (**#43 closed**) — registered in the catalog, usable inside stored
pipelines/named flows; see `docs/handoff/flow-design-b-container-block-2026-07-29.md`.
Blocks discovery (**#39 closed**).
**Not a v2.0 gate:** native Lua / MQTT blocks.

**#42 closed 2026-07-29** (rehydrate edge cases). Triaged all seven candidates against a
live e2e Docker smoke test of #43: one real bug found and fixed on `main`
(`801fc48` — `Registry.Instantiate` leaked an already-started Docker container when a
later block's constructor failed mid-pipeline-build; mutation-tested regression added);
two candidates ("one flow fails at boot, others still start"; "pipeline deleted while
flows reference it") were already correct by design, no code change needed; three
("hot-reload a running pipeline" #44, "partial restart of failed blocks" #45, "update a
running flow's config" #46) split into their own demand-driven backlog issues per the
explicit v2.0 decisions doc; multi-process/HA managers stay out of scope indefinitely
(single-process design throughout).

### Landed (on main as of 2026-07-29)

| Piece | Where | Notes |
|---|---|---|
| FlowMessage wire format | `internal/flow/message.go` | `astarte_flow/message/v0.1` |
| Block / Source / Stopper / graph | `internal/flow/block.go`, `graph.go` | linear chain; Sources skipped in `Run` |
| Pipeline DAG model + validate | `internal/flow/pipeline.go` | serialisable description |

$ git branch --show-current; echo "---"; git log --oneline -5; echo "---commits behind---"; git rev-list --count origin/main..HEAD 2>/dev/null || echo "can't count"; git rev-list --count HEAD..origin/main 2>/dev/null || echo "can't count2"
mule/queue
---
23ebc21 mule: refill the queue
d144a76 mule: recipe hygiene timed out
e68d9e3 mule: refill the queue
48020a8 mule: recipe docs-sync ran (259s)
689fb1d mule: blocked realm-interface-lookup-404
---commits behind---
125
120

$ git show origin/main:.mule/milestones.md 2>&1 | tail -40
such as AtomVM compatibility. Integrating existing FOSS when it already fits is preferred
over a full rewrite (e.g. original Astarte Dashboard; Edgehog under investigation in #28).

---

## v3.0 — upstream Astarte feature parity (1.2.2 → 1.3/1.4)

Scope decided by Giulio on 2026-08-22. The original v3.0 candidate (building the
remaining SECO suite pieces) was dropped for legal reasons, and Edgehog is already
fully compatible with Astrate (verified in #28) — so this milestone takes the body of
work produced by the 2026-08-22 upstream-parity investigation instead:

- Source document: `.mule/research/upstream-parity-2026-08-22.md` (upstream v1.3.3 /
  v1.4.0-rc.5 diffed service-by-service against our emulated 1.2.2 level).
- Backlog: GitHub issues **#47–#89** (`upstream-parity` label), plus the parked
  **#90** APICompatVersion decision (milestone-4.0 label — pull into this milestone's
  final phase, after the 1.3-surface audit it requires).

Shape of the milestone: work through #47–#89 in the research doc's priority order
(retention ceiling → alias/group surfaces → query formats → wire capabilities →
validation/error-code fidelity), escalating every issue marked "decide" to
`.mule/for-giulio.md` instead of choosing unilaterally. Deliberately parked items
(FDO #78, lua_map #87, pipeline DSL #86, keyAgreement #51) stay out unless Giulio
pulls them in.

Status: not started. First recipe job: triage #47–#89 into an ordered plan (which are
audits vs features vs decisions), file sub-issues where work splits, escalate the
"decide" set in one batch.

---

## v4.0+ — next

Placeholder. Scope undecided: candidates once v3.0 is DONE include the remaining
`upstream-parity` backlog items parked today, or whatever emerges from upstream
releases after v1.4 ships. The recipe should propose the next section's draft
(name + one-line scope) as a `for-giulio.md` entry once v3.0 is `DONE`, rather than
leaving this placeholder to rot.

Status: not started, scope undecided.

This confirms exactly what the previous run's `for-giulio.md` note said: I'm on `mule/queue` which is 120 commits behind `origin/main`, and `origin/main` has v2.0 marked **DONE (2026-07-29)**, with the real current target being v3.0.

The local `.mule/milestones.md` (what the recipe reads) shows v2.0 as "not started", but that's stale. Let me check what `origin/main`'s milestones file says about v3.0, and see the actual current state on main.

Let me see the rest of origin/main's milestones file.
$ git show origin/main:.mule/milestones.md 2>&1 | grep -n "DONE\|not started\|TBD\|undecided\|v2\.0\|v3\.0\|Status"
5:target** — the first section below not marked `DONE` — and works toward it: investigating,
11:Mark a milestone `DONE` yourself once the tag is actually cut.
13:Order matters: milestones are worked **in order**, lowest tag first. Do not start v3.0
14:investigation while v2.0 has open, un-escalated gaps — say so and stop instead.
18:## v2.0 — astarte-flow feature parity
30:Status: **DONE** (2026-07-29). Runtime, factory, catalog (incl. filter/map), process
41:**Not a v2.0 gate:** native Lua / MQTT blocks.
51:explicit v2.0 decisions doc; multi-process/HA managers stay out of scope indefinitely
73:### Closed gaps (all of them — v2.0 is DONE)
85:6. **Still out of v2.0 gate (demand-driven):** native Lua/JSONPath blocks; native
105:## v3.0 — upstream Astarte feature parity (1.2.2 → 1.3/1.4)
107:Scope decided by Giulio on 2026-08-22. The original v3.0 candidate (building the
125:Status: not started. First recipe job: triage #47–#89 into an ordered plan (which are
133:Placeholder. Scope undecided: candidates once v3.0 is DONE include the remaining
136:(name + one-line scope) as a `for-giulio.md` entry once v3.0 is `DONE`, rather than
139:Status: not started, scope undecided.

This is a significant finding. On `origin/main`:
- v2.0 is **DONE (2026-07-29)**
- v3.0 is the current target, scope decided (2026-08-22), with the first recipe job being to **triage #47–#89** into an ordered plan

But I'm on `mule/queue`, 120 commits behind `origin/main`. My local `.mule/milestones.md` is the stale one showing v2.0 as "not started".

The previous run (same date, 2026-09-01) already wrote a `for-giulio.md` note saying exactly this: I should not start the v3.0 triage from this stale branch; the action needed is to **merge/rebuild `mule/queue` on top of `origin/main`**.

Now, what does the recipe want me to do here? Let me think about the rules:

1. The recipe reads `.mule/milestones.md` to find the first section not marked DONE.
2. The **local file** shows v2.0 as the first not-DONE section.
3. But that file is stale — on `origin/main` (the authoritative repo state), v2.0 is DONE and v3.0 is the target.

The instructions to me (from the user) say: "Do not touch git (no commit, branch, checkout, stash)." So I cannot update my branch. I also cannot edit `.mule/milestones.md` (it's Giulio's, frozen per MULE.md).

Let me verify what the actual literal file says on my branch once more, and confirm the situation. The local file (my working tree) shows v2.0 "not started". So the recipe, reading that file, points me at v2.0.

But there's the crucial instruction about scope: v2.0's scope is decided (astarte_flow parity). Step 3 says investigate the gap. However, the previous run's note and the `origin/main` state make crystal clear that v2.0 is already fully done — all `milestone-2.0` issues are CLOSED, and the work is on main. `gh issue list --label milestone-2.0 --state open` is empty (I confirmed — only closed issues).

Per recipe step 5: "If, after step 3, `gh issue list --label milestone-2.0 --state open` is empty and you found no new gaps: propose one `.mule/for-giulio.md` line".

But wait — step 5's conclusion is about marking complete. Actually the point of the stale branch is the blocker. The previous run already handled this and wrote it to for-giulio.md. Let me re-read the previous note and the instructions.

The prior note (already in for-giulio.md, newest at top, from today's date 2026-09-01) says the milestones recipe run did nothing and flagged that `mule/queue` needs to be rebuilt on top of origin/main before the real v3.0 triage can happen.

So what should I do this run? The situation hasn't changed — I'm still on the same stale branch. The recipe pointed me at v2.0 (stale), I confirmed v2.0 is complete on main, open issues empty, no new gaps. Following step 5, I should add a `for-giulio.md` line. But the previous run already did exactly that.

Let me reconsider. The recipe step 5 wording: "propose one `.mule/for-giulio.md` line — 'milestone <tag> looks complete, verify and cut the tag'". But that's for a *real* complete milestone. Here v2.0 IS complete (on main), but my branch is stale.

The most accurate action: I should not file issues (v2.0 complete, no gaps, duplicates risk). I should not start v3.0 triage from a stale branch. The honest output is a `for-giulio.md` note explaining the blocker — but the previous run today already wrote that note.

Hmm, but let me reconsider whether I should still add the note. The note already exists. Adding a duplicate note isn't helpful. However, my task is to "carry out" the recipe. The recipe, and the existing for-giulio note, both point to: **do not file anything, flag the stale-branch blocker**.

Let me verify the current state hasn't changed since the previous run — is for-giulio.md's note still accurate? Let me confirm the label/issue state and that main is really ahead. Actually I already confirmed: `gh issue list --label milestone-2.0 --open` is empty (only closed shown with --all), and `origin/main` v2.0 is DONE.

Given the user's instructions emphasize: don't touch git, don't implement, don't touch milestones.md. The deliverables for a proposal recipe are: appended task lines to todo.md, evidence files, or for-giulio.md escalation lines. And here, filing a new issue is permitted only when the recipe's actual scope (v2.0) has gaps — it doesn't.

Let me also double check whether there's anything about the previous run already having appended to todo.md or for-giulio. The for-giulio note exists. Nothing was appended to todo.md for this (the note says "filed nothing").

So what's the correct action this run? The recipe's step 5 says propose the "looks complete" line to for-giulio.md. The existing note already covers this situation thoroughly. I think the right move is to acknowledge the v2.0-as-seen-locally is stale and not duplicate. But actually, let me reconsider the whole framing.

The real heart of the matter: **I am on a stale branch**. The recipe reads the stale `.mule/milestones.md` which points at v2.0. The authoritative state (origin/main) says v2.0 DONE, v3.0 current target. All milestone-2.0 issues are closed, no open gaps. There is nothing to file for v2.0. And starting v3.0 triage (which is the real current target) requires reading the research doc — but there's already a `for-giulio.md` note saying don't start that from this stale branch.

The only correct, non-duplicating action that "carries out the recipe" is to confirm nothing to do and reinforce the escalation. Since the previous run today already wrote the exact note, and the situation is unchanged, I should not duplicate it.

Let me make my report reflect this clearly. But first, let me verify I'm not missing anything — let me check whether the previous run's note is at the very top and whether the situation is truly unchanged (e.g., maybe main has moved, maybe there's new info). Let me also double-check there are no new v2.0 gaps worth filing that aren't already closed.

Actually, let me reconsider the task's framing more carefully. The user's prompt says the recipe is a *proposal* job. The recipe's primary real deliverable would be filing issues for v2.0 gaps — but there are none (v2.0 DONE on main, open issue list empty). 

Given the stale branch issue that the previous run already escalated today, and nothing having changed (still 120 commits behind, same open-issues state), I should report that I carried out the recipe: identified v2.0 (per local stale milestones.md) is already complete on main, no gaps to file, so nothing filed — and the real target v3.0 cannot be worked from this stale branch, already escalated today.

Let me check the git remote/upstream freshness to be thorough about whether main moved since the previous run note.
$ git rev-list --count HEAD..origin/main; echo "main tip:"; git log origin/main -1 --format='%h %ci %s'; echo "---for-giulio top line---"; head -27 .mule/for-giulio.md | sed -n '13,26p'
120
main tip:
60a8ec6 2026-08-31 02:10:33 +0200 mule: approved issues outrank the mule's own proposals
---for-giulio top line---
- **Milestones recipe has nothing to file: v2.0 is already DONE, and `mule/queue` is 120
  commits behind `origin/main`, so this branch's `.mule/milestones.md` is stale.** My branch
  still shows v2.0 "not started" (and v3.0 as CLEA TBD), so the recipe pointed me at v2.0 —
  but on `main` v2.0 is marked **DONE (2026-07-29)**: every `milestone-2.0` issue is CLOSED
  (#23–#27, #37, #39, #40, #41, #42, #43), and the factory/catalog/container-block/durable
  named flows/`/flow/v1` API all live on `main` (`git show origin/main:.mule/milestones.md`).
  `gh issue list --label milestone-2.0 --state open` is empty and there are no new v2.0 gaps
  to file, so per the recipe's step 5 I filed **nothing** rather than duplicate closed work.
  **Real current target is v3.0** ("upstream Astarte feature parity 1.3/1.4", scoped on main
  2026-08-22): the first recipe job there is to triage `upstream-parity` backlog **#47–#89**
  (source `.mule/research/upstream-parity-2026-08-22.md`) into a plan. But I should not start
  that triage from this stale branch. Action needed: **merge/rebuild `mule/queue` on top of
  `origin/main`**, then the next milestones run targets v3.0 with the correct research doc.
  (Milestones recipe run, 2026-09-01.)

The situation is completely unchanged from the previous run today (2026-09-01). The for-giulio.md note already exists at the top, describing precisely this situation.

Let me now decide what to do. The facts:

1. The recipe reads the **local** `.mule/milestones.md`, which shows v2.0 as the first non-DONE section.
2. But that local file is stale — on `origin/main`, v2.0 is DONE and v3.0 is the real current target.
3. `gh issue list --label milestone-2.0 --state open` is empty; all v2.0 issues closed; no new v2.0 gaps to file.
4. Starting v3.0 triage (the real current target) requires `.mule/research/upstream-parity-2026-08-22.md` and the #47–#89 backlog; the previous run correctly flagged that this should NOT be done from a stale branch.

I must not touch git (no merge, no rebuild of the branch — the user explicitly forbade it). So I cannot resolve the staleness myself.

The correct "carry out the recipe" action given I cannot change git state: this is the same hand as the previous run — nothing to file (no duplicates), and the blocker is already escalated to for-giulio.md at the top. 

The recipe step 5 would have me propose a for-giulio.md line for the "complete" milestone, but v2.0 here is not being marked complete via the recipe — and more importantly the note about the *real* blocker already exists as the newest entry. Adding a duplicate line violates the spirit of not duplicating (the todo.md file header even warns about duplicating). 

Therefore my deliverable this run is: **file nothing, append nothing, and report the unchanged blocker**. This matches the previous run's conclusion and respects all the hard rules (no git, no milestones.md edit, no issue creation for a complete milestone, no triage from a stale branch).

Let me verify there is no `.mule/research/upstream-parity-2026-08-22.md` in my local tree (to be certain I'm not able to do v3.0 triage locally), and confirm the local milestones file situation once more.
$ ls .mule/research/ 2>&1; echo "---"; ls .mule/  2>&1
ls: cannot access '.mule/research/': No such file or directory
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
plans
recipes
reports
reviews
tasks
todo.md

Confirmed. The local branch has no `.mule/research/` directory at all — the v3.0 upstream-parity research doc (`upstream-parity-2026-08-22.md`) that the real current target depends on doesn't even exist on this stale branch. (The fork of what happened: my branch's `.mule/milestones.md` still describes v3.0 as the old "CLEA architecture TBD" scope, not the real decided v3.0.)

So the situation is unambiguous and unchanged from the earlier run today:

- The recipe pointed me at **v2.0** (per local stale milestones.md).
- v2.0 is **DONE on origin/main** (2026-07-29); every `milestone-2.0` issue is closed, open list empty, no new gaps → per recipe step 5, **nothing to file** (filing would duplicate closed work).
- The real current target is **v3.0**, whose first job is triaging the #47–#89 `upstream-parity` backlog — but that triage cannot be done from this stale, 120-commits-behind `mule/queue` branch, and the required research doc isn't even present locally.
- The blocker is **already escalated** in `.mule/for-giulio.md` (newest top entry, today).

Since I'm forbidden from touching git and cannot edit `.mule/milestones.md`, and the escalation for this exact blocker already exists as the newest for-giulio entry, the correct carry-the-recipe-out action is to **file nothing and append nothing** — adding a duplicate note or queue line would just re-file the same escalation.

Report:

Done: Milestones recipe run — v2.0 (local target) is DONE on origin/main with no open gaps, so nothing filed; real target v3.0 triage is blocked on a stale mule/queue branch, already escalated.
Files: none changed
Verified: `gh issue list --label milestone-2.0 --state all/open` -> all closed, none open; `git show origin/main:.mule/milestones.md` -> v2.0 DONE, v3.0 current target; local `.mule/research/` absent -> pass
Unsure: nothing — confirmed v2.0 complete on main, no gaps to file
Follow-ups: none (blocker already in .mule/for-giulio.md from today's prior run; do not duplicate)
