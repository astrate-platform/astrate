# Branch `mule/research`

Output of the daily cross-project survey (`tools/mule-survey.sh`, `.mule/recipes/cross-project-survey.md`),
kept separate from `mule/queue` deliberately — see the commit that introduced this branch for why.

Each day this branch gets either:

- a new `.mule/research/survey-<date>.md` + `.mule/research/issues-<date>.md`, or
- nothing (a one-line note in `.mule/research/log.md`), when nothing material changed since
  the last survey.

**Nothing here is a queued task.** `issues-<date>.md` candidates need a human or a strong
model to triage before anything becomes a GitHub issue labelled `mule` — some of them may
look machine-checkable and not be (see `.mule/for-giulio.md` and issues #17-#19 on the main
repo for the shape of that trap: recording against a live upstream needs an architect, not an
unsupervised agent with no diff review).

The worked example this recipe was distilled from — read before triaging a new day's
output — is `survey-2026-07-27.md` / `issues-2026-07-27.md`, right here on this branch, from
the one-off run that produced GitHub issues #12-#19 and four entries in `.mule/for-giulio.md`
on `main`. Triage is done with the `mule-triage` skill; `triaged.md` in this directory tracks
which dates have already been handled.
