# Recipe — milestones

Work toward the next release-tag milestone in `.mule/milestones.md`. Unlike the other
recipes, this one's proposed tasks are allowed to **file GitHub issues directly** (via
`gh issue create`) rather than only appending to `.mule/todo.md` — a filed issue is the
deliverable here, not a step toward one. It is still a proposal job in every other sense:
nothing gets implemented by this recipe, and nothing merges without going through the normal
queue/gate.

## Do this

1. Read `.mule/milestones.md`. Find the **first** section not marked `DONE` — that is the
   only milestone you work on this run. If an earlier milestone still has open,
   un-escalated gaps (see step 4), stop and say so instead of moving to the next one.
2. If that milestone's scope says "not yet decided" / "TBD": do not investigate the
   codebase. Read the linked reference docs, propose 3-5 candidate pieces or interpretations
   with a one-line scope each, and write them as a single `.mule/for-giulio.md` entry ("v3.0
   scope: options are — a) ... b) ... c) ...”). Stop there for this run.
3. Otherwise, investigate the gap between the reference and Astrate's current code:
   - Read the reference (upstream repo docs/README, or the linked doc site) for what the
     milestone actually requires, same care as `.mule/recipes/astarte-upstream.md` — a
     README claim is not a fact, check the code/schema it describes when it matters.
   - `gh issue list --label milestone-<tag> --state all --limit 50` — what has already been
     filed, so you never duplicate.
   - Grep Astrate's own `internal/`, `docs/DESIGN.md`, `docs/ROADMAP.md` for existing
     coverage of each capability the reference names.
4. For each gap, decide which bucket it falls in:
   - **Machine-checkable, no design choice needed** (a missing package, an unimplemented
     wire message, a schema field): **file the issue now**, directly:
     ```
     gh issue create --title "<slug>: <one-line outcome>" --label mule,milestone-<tag> \
       --body "<what/why, cite the reference and the file(s) that would need to change,
       ending with 'Acceptance: <a concrete test or check that proves this is done>'>"
     ```
     This is a `gh` call, not a git mutation, and filing it is the deliverable — do not
     also write a `.mule/todo.md` task line proposing to create the same issue. One or the
     other, never both, or the next run duplicates it.
     **Every issue body must end with an `Acceptance:` line naming a specific test or
     check** — "a table-driven test in internal/flow/router_test.go asserting in-order
     delivery within a stream key", not "add tests". If you cannot state one, the gap is
     too vague to be issue-shaped yet — investigate further or escalate it, don't file it.
     This is what turns the milestone recipe's issues into starting points for MULE.md's
     existing rule that no behaviour change lands without a test proving it — the recipe
     names *which* test, the executing task still has to write and prove it.
     If the piece is big enough to need sub-issues, file the parent issue now, then a
     follow-up task line — `- [ ] milestone-<tag>-issue-<slug>-sub1: gh issue create ...
     --body "part of #<parent>, ..."` — for each child, since the parent's issue number is
     only known after it's filed. **Only file one level per run**; let the next run see the
     new parent number in `gh issue list` before filing children.
   - **Needs a design decision** (an API shape, a protocol extension, a choice the reference
     itself doesn't pin down): write it to `.mule/for-giulio.md`, one line, exactly as
     `.mule/recipes/github-issues.md` already does for issues like this. Do not file a
     GitHub issue for it.
   - **Already covered**: say so, propose nothing.
5. If, after step 3, `gh issue list --label milestone-<tag> --state open` is empty and you
   found no new gaps: propose one `.mule/for-giulio.md` line — "milestone <tag> looks
   complete, verify and cut the tag" — and stop. Do not mark `.mule/milestones.md` `DONE`
   yourself; that file is Giulio's.

## Rules

- **File at most five issues per run**, same cap as the other recipes' proposals.
- **Never file the same gap twice.** `gh issue list --label milestone-<tag> --state all` in
  step 3 is what prevents this — trust it, and never both file an issue and queue a task
  line to create that same issue (see step 4).
- **Never run `gh issue edit/close/comment`** or anything else that changes an *existing*
  issue or touches git — creation of new issues is the one exception to "recipes only
  produce text", nothing else is.
- Every filed issue must name the file(s) or package the gap lives in, if the investigation
  found one — that's what makes the follow-up task machine-startable instead of another
  investigation.
- Never port upstream code or its language-specific structure — port the capability,
  restated for Astrate's Go codebase.
- If the reference named in the milestone section (upstream repos, docs) is unreachable,
  say so and stop; do not guess at scope from the milestone's one-paragraph summary alone.
- This recipe never touches `.mule/milestones.md`. Status changes go through
  `.mule/for-giulio.md`, same as every other frozen-file rule in `.mule/MULE.md`.
