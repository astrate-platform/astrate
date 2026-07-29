---
name: mule-triage
description: Triage the daily cross-project survey Big Pickle writes to branch mule/research in astrate — read what it found, verify the claims against the current code, and turn the good candidates into GitHub issues (mule-labelled) or .mule/for-giulio.md entries. Use when Giulio asks "cosa ha trovato la survey", "controlla la survey", "triage the research", "review what the daily job found", or after being told the mule-survey.timer ran on the Pi.
---

# Triaging the daily cross-project survey

The daily job (`tools/mule-survey.sh`, `.mule/recipes/cross-project-survey.md`, astrate repo)
writes dated reports to branch `mule/research` and **never files an issue itself** — that
judgment call is this skill. It exists because the first run of this survey (2026-07-27)
needed real correction before anything was safe to queue: three of its "mule-spec" candidates
actually needed a live-upstream recording, which this project already learned (M12 phase 04)
is architect work, not something an unsupervised agent with no diff review should touch. That
run is the calibration example — `survey-2026-07-27.md` / `issues-2026-07-27.md`, committed
to `mule/research` itself (not `main` — read them with `git show
origin/mule/research:.mule/research/survey-2026-07-27.md`) — and the resulting issues
**#12–#19** plus four entries in `.mule/for-giulio.md` on `main`.

## Steps

1. `cd` to the astrate repo, then `git fetch origin mule/research`.
2. Find what's untriaged: `git show origin/mule/research:.mule/research/triaged.md 2>/dev/null`
   (a plain list of dates already handled — create the file with the header below if it
   doesn't exist yet) against `git ls-tree -r --name-only origin/mule/research -- .mule/research`
   for `survey-*.md` files. Anything not in `triaged.md` is new.
3. **If nothing is untriaged**, say so and stop — most days this is the correct outcome, the
   survey is incremental and often writes nothing at all (check `.mule/research/log.md` on
   that branch for "no material change" rows, which need no triage).
4. For each untriaged date, oldest first:
   - Read both files: `git show origin/mule/research:.mule/research/survey-<date>.md` and
     the matching `issues-<date>.md`.
   - **Verify, don't trust.** For every candidate, grep/read the *current* `main` (the survey
     may be a day or more stale) to confirm the gap still exists — the calibration run's whole
     issue list was grep-confirmed this way before filing.
   - **Re-check the scope label, don't inherit it.** The recipe asks the survey to bucket each
     candidate as `mule-line` / `mule-spec` / `design-decision`, and it can get this wrong the
     same way the first run did. Downgrade out of the mule queue anything whose verification
     needs a live external system with no automatable pass/fail (the self-hosted upstream
     Astarte on the Legion Go is the recurring case in this repo), or anything that is really a
     design call dressed up as a task. When in doubt, the failure mode to avoid is a false
     `mule-line`/`mule-spec` landing in the unsupervised queue, not an overly cautious
     `design-decision`.
   - Check for duplicates before filing: `gh issue list --search "<slug or key phrase>"` and a
     skim of `.mule/for-giulio.md` — a survey re-run may re-find something already tracked.
   - **File it:**
     - Genuinely safe, machine-checkable, one-or-two-file candidates → `gh issue create --label mule --title "..." --body "..."`, title alone must stand as a runnable one-line task (the mule loop only ever reads the issue *title*, never the body — see `tools/mule.sh`'s `issue_tasks()`). Add `--label readonly` for investigate-only tasks, `--label legion` for anything needing the Legion Go.
     - Bigger but still machine-checkable (a few files) → same as above, then write a full
       spec at `.mule/tasks/issue-<N>.md` (mirrors a trickle phase spec: context, what to do,
       constraints, acceptance criteria) once you have the issue number, commit it to `main`.
     - Needs a live upstream recording, or any other check the mule cannot run unsupervised →
       `gh issue create` **without** the `mule` label (use `enhancement` or whatever fits),
       body explains why it's excluded from the mule queue.
     - Needs a human design decision → append to `.mule/for-giulio.md` on `main`, do not file
       an issue, unless one already exists for it (link instead).
   - Cap yourself the same way the recipe caps the survey: prioritise, don't file everything
     that's technically true. A pile of low-value issues is worse than a short list Giulio
     will actually read.
5. Append the date(s) you triaged to `.mule/research/triaged.md` on `mule/research`
   (`# Dates triaged by mule-triage — one per line, oldest first` as the header if creating
   it), commit, push.
6. Report to Giulio: what was found, what got filed where (with issue numbers/links), what
   went to `for-giulio.md`, and anything you excluded and why.

## What NOT to do

- Don't re-run the survey yourself — that's the daily timer's job. This skill only reads what
  it already wrote.
- Don't skip the verification pass because the survey "looked careful" — the calibration run
  was careful and still got three items wrong in a way that mattered.
- Don't file into `.mule/todo.md` — issues are the only intake for the mule queue in this
  project (three merge conflicts in one afternoon is why, per `.mule/MULE.md`).
