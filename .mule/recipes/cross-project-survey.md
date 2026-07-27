# Recipe — daily cross-project survey (astarte-platform, AtomVM)

Run by `tools/mule-survey.sh` on a **daily** timer, separate from the 30-minute mule tick.
This is a bigger, slower job than a normal recipe: it clones two upstream repos and reads
docs, not just release notes. It is also the only recipe that runs on its own dedicated
branch (`mule/research`) instead of proposing into `.mule/todo.md` — its output needs a
human or a strong model to read and triage before anything becomes a queued task, exactly
like the one-off run this recipe was distilled from (`.mule/research/survey-2026-07-27.md`,
`issues-2026-07-27.md`, still on `main` history as the worked example — read the *shape* of
those two files before writing new ones, they are the format to match).

## Be incremental — this is the part that keeps a daily job cheap

Most days, nothing worth reporting changed. Check before doing the expensive thing:

1. Find the most recent prior `survey-*.md` in this same directory (or on the `mule/research`
   branch history if this checkout doesn't have it — `git log --all --oneline -- .mule/research/`).
   If none exists, this is the first run: do the full survey (skip to "Do the survey" below).
2. Otherwise, cheaply check whether anything could have changed since that run's date:
   - `git -C .mule/research/upstream/astarte log -1` after a fresh shallow pull (or
     `git ls-remote https://github.com/astarte-platform/astarte HEAD` if you don't want to
     re-clone yet) — has upstream's default branch moved?
   - `git log --oneline <date-of-last-survey>..HEAD -- internal/ docs/` in this repo — has
     Astrate's own code or docs changed since?
3. **If neither moved, write a one-line entry to `.mule/research/log.md`** (create it if
   absent) — `<date>: no material change since <last-date>, checked <astarte-platform HEAD>
   / <astrate main HEAD>` — and stop. Do not write a new survey/issues file. An empty result
   is a good result; say so and stop, exactly as `.mule/recipes/astarte-upstream.md` already
   does for its narrower check.
4. If something moved, do the survey below, but **scope it to what changed** — you do not
   need to re-read pairing, MQTT protocol, and standard interfaces every day if only the
   trigger docs changed upstream. Read the diff (`git log -p <old>..<new> -- doc/` in the
   upstream clone) to know what to focus on, the same way `.mule/recipes/astarte-upstream.md`
   reads release notes before deciding what to look at.

## Do the survey (first run, or when something moved)

Same five sources as the original one-off pass:

1. **Astrate's own codebase and docs** — `docs/DESIGN.md`, `docs/ROADMAP.md`,
   `docs/COMPATIBILITY.md`, package layout, `// TODO` markers.
2. **`mule/queue` branch** — `git log --oneline origin/mule/queue` (not `mule/queue`, which
   may not exist locally) and `.mule/log.md`, so you don't re-propose what's already covered
   or already blocked.
3. **Upstream `astarte-platform`** — shallow clone into `.mule/research/upstream/astarte`
   (delete it at the end of the run, see "Cleanup" below — the Pi has limited disk).
4. **Parity gap analysis** between 1 and 3, wire-visible behaviour first.
5. **Upstream `atomvm/AtomVM`** — shallow clone into `.mule/research/upstream/AtomVM`, same
   cleanup rule. The connection to Astrate was already established once (it motivates the
   JSON payload profile, `docs/DESIGN.md` §3.5) — you do not need to re-derive that unless
   you find something that changes the conclusion; say so if you do.

Write `.mule/research/survey-<today>.md` and `.mule/research/issues-<today>.md`, same format
as the worked example. **Cap candidates at 8**, tighter than the one-off run's 15 — a daily
job that produces a dozen items a day is not being read. Every `mule-line`/`mule-spec`
candidate must be grep-confirmed not already implemented, exactly as before.

## Rules carried over from the one-off run and from the other recipes

- **Never copy upstream code into this repo.** Read the rule or behaviour, propose that
  Astrate implement or verify it — never port an implementation.
- **Never file a GitHub issue and never run `gh issue create`.** That is a deliberate design
  choice for this recipe: filing needs judgment about scope (mule-line vs. mule-spec vs.
  design-decision) and a check for the kind of trap the live-upstream-recording candidates
  hit last time (they looked like normal mule-spec work and were not — see
  `.mule/for-giulio.md`'s design-decision entries and the three issues that were deliberately
  left unlabelled for the mule). That triage is a human/strong-model job. Your output is the
  input to that triage, not the triage itself.
- **Never touch anything outside `.mule/research/`.** No commits to `.mule/todo.md`, no
  edits to `docs/DESIGN.md`/`docs/ROADMAP.md`/`docs/COMPATIBILITY.md` (propose changes in the
  issues file). Do not touch git yourself — `tools/mule-survey.sh` commits and pushes what you
  wrote after you finish, the same separation `tools/mule.sh` keeps for queued tasks.

## Cleanup

Before finishing, `rm -rf .mule/research/upstream` — the clones are disposable and the Pi's
disk is not large. Keep only the dated `.md` files.

## When you're done

Finish with a short status line: full survey or "no material change", how many candidates in
each scope bucket if a survey was written, and anything that blocked you (couldn't clone,
couldn't reach GitHub, etc — say so rather than papering over it).
