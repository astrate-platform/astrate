---
name: mule
description: Solo-mule mode — Big Pickle (opencode) works the Astrate repo by itself off a queue, one task per fresh process, gated by tests and git. Use when Giulio says "mule", "modalità mulo", "fai lavorare Big Pickle da solo", "svuota la coda", "fagli fare le sue cose", asks what the mule has been doing, or wants the queue topped up. Not the same as /trickle — trickle has an architect per phase, this has none.
---

# Solo-mule mode

One cheap model, working a queue, unsupervised. `/trickle` is for work *you* have designed;
this is for work that keeps a project healthy while you are not looking.

| | trickle | mule |
| --- | --- | --- |
| who designs the work | architect, per phase | a recipe, then Giulio approves the queue |
| isolation | throwaway worktree | branch `mule/queue`, per-task commits |
| what catches errors | a model reads every diff | the gates, plus a human reading commits later |
| unit of work | a phase spec | a one-line task |

The consequence: **only put machine-checkable work in the mule queue.** There is no diff read
between the mule and a commit. If correctness needs a human read, it is trickle work or your
work — not this.

## The three machines

Hostnames/IPs live in `.mule/hosts` (gitignored — `.mule/hosts.example` is the committed
template) as `MULE_PI_SSH` and `MULE_LEGION_SSH`, sourced by `tools/mule.sh`. `source .mule/hosts`
before running any of the ssh snippets below.

| | role | facts that matter |
| --- | --- | --- |
| **Mac** | where Giulio and you work | full gate: `go test -race`, pinned golangci-lint |
| **Pi** `$MULE_PI_SSH` | runs the queue on a 30-min timer, unattended | DietPi arm64, 4 cores, 3.7GB, **no Docker**, clone at `/root/astrate-mule`, opencode + big-pickle installed |
| **Legion Go** `ssh $MULE_LEGION_SSH` | the muscle | CachyOS x86_64, 16 cores, 12GB, Docker running the full upstream Astarte stack, **no Go toolchain** |

**The Pi cannot run the race detector** — 39-bit VMA kernel, ThreadSanitizer needs 48. Its
gate is `go vet ./... && go test ./...` (green, ~3m, 20 packages). So **never queue
concurrency work for the mule**: nothing on that machine can catch a data race. `.mule/config`
detects this per host rather than hardcoding, so the same file gives the Mac the full gate.

A task line tagged `[legion]` runs only when the Legion Go answers; otherwise it is skipped
and the queue moves on. The bench binary is cross-compiled from wherever the mule runs and
copied over (`mule.sh legion bench-push`) — do not install Go on that machine.

## The mechanics

Everything lives in the repo: `tools/mule.sh`, `.mule/MULE.md` (the mule's rules),
`.mule/todo.md` (the queue), `.mule/recipes/` (six proposal jobs), `.mule/log.md`,
`.mule/for-giulio.md` (escalations), `.mule/config` (gates and the never-touch list).

```bash
tools/mule.sh preflight          # once per sitting — baseline must be green
tools/mule.sh status
tools/mule.sh recipe review      # PROPOSES tasks into todo.md; does not execute
tools/mule.sh loop               # drains the queue, then prints the menu
tools/mule.sh legion check       # ssh + docker on the Legion Go
tools/mule.sh revert             # undo the last mule commit
bash tools/mule-setup-pi.sh      # (re)install the Pi timer; idempotent
```

**`tick`** is the timer's entry point and belongs to the Pi, not to you: `flock -n` so runs
never overlap, **one task per tick and never a loop**, a daily cap (`MULE_DAILY_MAX=16`), and
**an empty queue does nothing at all** — it never runs a recipe to invent work for itself.
That timidity is deliberate; the provider is free and nobody is watching. If Giulio ever asks
why the mule "isn't doing anything", check the queue and the budget file before assuming a
fault — idle is the designed state.

Watching it:

```bash
source .mule/hosts
ssh "$MULE_PI_SSH" 'journalctl -u mule.service -n 50'
ssh "$MULE_PI_SSH" 'systemctl list-timers mule.timer'
ssh "$MULE_PI_SSH" 'systemctl stop mule.timer'      # the off switch
```

`next` runs one task in a **fresh opencode process** and then, if `MULE_TEST_CMD` and
`MULE_LINT_CMD` pass and nothing on the never-touch list moved, commits it. Otherwise the
tree is restored, the fragment is kept in `.mule/failed/`, and the line is marked `- [!]`
with the reason. Two blocked tasks in a row stops the loop.

The fresh process per task is the whole design. Big-pickle has ~30k of good context out of
200k, so the loop lives in bash and no task inherits another's session.

## What Giulio actually asks for, and what to do

**"top up the queue" / "fai qualcosa di utile"** — run a recipe, then *read what it
proposed before he does*. Cut anything that is not machine-checkable, anything touching more
than two files, anything vague. A five-line queue he approves beats a twenty-line one he
ignores.

**"cos'ha combinato?"** — `tools/mule.sh review`. That prints the commits, the diffstat,
unapproved `[auto]` lines and the standing check results, separating code from the mule's own
bookkeeping. Then **actually read the diff it points at**. The gates prove it compiles, the
tests pass, a new test fails without the change, and no frozen file was touched — none of
which means the change was a good idea. Report what landed, what blocked, and anything in
`.mule/for-giulio.md`.

Reviewing is the job Giulio explicitly wants a strong model for, so do it properly rather
than summarising the log: read the code, say whether it should reach `main`, and say plainly
when the honest answer is "this is fine but it is not worth much".

**"merge it"** — review first, then `git checkout main && git merge --no-ff mule/queue`.
Never merge unread. Nothing merges automatically anywhere in this system, deliberately.

**Giving it work** — file a GitHub issue labelled `mule`; do not edit `.mule/todo.md`. Issues
are the queue and are read live on each tick; `todo.md` holds only standing lines, and
copying issues into it caused three merge conflicts in one afternoon. Labels carry the tags
(`legion`, `readonly`) and the state (`mule` → `mule-review` → your call). The mule never
closes an issue.

    gh issue create --label mule --title "<slug>: <outcome>" --body "<detail>"

**Issues labelled `mule-alarm`** are the dead-man's switch: the mule has landed nothing for
8+ hours and said so itself. Do not treat it as a task — diagnose it. It almost always means
the provider is refusing every run, the queue has nothing runnable, the tree is dirty so
every tick aborts, or push is failing and the work is stranded on the Pi.

**a blocked task** — read `.mule/failed/<slug>.diff`. Usually one of: the task was too big
(split it into two lines), the task was ambiguous (rewrite the line naming the file), or the
task needed a decision (move it to `.mule/for-giulio.md`). Rewriting the line is the fix
almost every time; re-running an unchanged blocked line is not.

## Writing a task line

Same failure modes as trickle specs, one line instead of a spec:

- **Name the outcome and the file.** "fix logging" is unrunnable; "trigger-log-level: demote
  the per-message dispatch log in internal/triggers/dispatch.go from Info to Debug" runs.
- **One or two files.** Three files is a task that will time out and yield nothing.
- **If it needs more than a line, write `.mule/tasks/<slug>.md`** — the script picks it up
  automatically and it becomes a real spec, with line ranges and acceptance criteria, like a
  trickle phase.
- **Anything asserted only by a log line, a metric, or an external system will pass its gate
  while being wrong.** Bind it with a test or do not queue it.

## Keeping it honest

`.mule/log.md` is the evidence: `secs` near the 900s budget means the task was too big, and
a recurring blocked reason means a recipe is proposing the wrong kind of work. When a recipe
keeps producing tasks you cut, fix the recipe — that is where the leverage is, not in the
individual task lines.
