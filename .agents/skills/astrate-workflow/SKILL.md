---
name: astrate-workflow
description: Discoverability map and live-status dashboard for everything automated (or semi-automated) across the astrate dev workflow — Mac, the Raspberry Pi (mule + daily survey timers), and the Legion Go (upstream Astarte + benchmarks). Use when Giulio asks "cosa posso fare", "cosa gira", "cosa è automatizzato", "dashboard", "stato del workflow", "cosa c'è per me", "come lancio X", "help astrate", or seems unsure which command/skill covers something he wants to do or check.
---

# Astrate workflow map

This is a **router, not a worker**: it tells you (and Giulio) what exists, where it runs, how
to trigger it, and whose turn it is — then hands off to the real skill/command, or runs the
cheap read-only status checks itself. It does not replace `mule`, `trickle`, or `mule-triage`;
it exists because there are now enough of them that "how do I..." needs one answer, not a
grep through three skill files.

**When invoked, do two things, in order:**

1. Run the live-status checks below (they're all cheap and read-only) and show a compact
   summary — don't dump raw command output, digest it.
2. If Giulio asked a specific "how do I do X" question, answer it from the table and actually
   do it (invoke the skill, run the command) rather than just describing it. If he asked
   generally ("cosa posso fare"), present the table and ask what he wants.

## The three machines

Hostnames/IPs live in `.mule/hosts` (gitignored — `.mule/hosts.example` is the committed
template) as `MULE_PI_SSH` and `MULE_LEGION_SSH`, sourced by `tools/mule.sh`. `source .mule/hosts`
before running any of the ssh snippets below.

| | role | reach it | what runs there |
|---|---|---|---|
| **Mac** | you work here | — | trickle (architect work), one-off skill runs, `gh` |
| **Pi** `$MULE_PI_SSH` | unattended, always on | `ssh $MULE_PI_SSH` (bash -s, its shell is fish) | `mule.timer` (30 min), `mule-survey.timer` (daily ~03:00) |
| **Legion Go** `$MULE_LEGION_SSH` (alias, works from Mac and Pi) | muscle, sleeps/travels | `ssh $MULE_LEGION_SSH` | upstream Astarte in Docker, `bench/`, the race-detector gate |

This skill itself is installed in two places — Claude Code on the Mac, and opencode on the Pi
(`~/.claude-skills`, see `tools/mule-setup-pi.sh`'s skills step) — so it may be *you* running
it from either machine. That's exactly why the self-detection below exists.

## Live status — run these, digest, don't paste raw

**First, work out which machine you're actually running on.** This skill is installed on the
Mac (Claude Code) *and* on the Pi (opencode, `/root/.claude-skills` — see the setup commit),
so "check the Pi" must never mean "ssh to myself": run the same checks locally instead. Detect
it once, locally, no network needed:

```sh
WHERE=mac
if [ "$(uname -s)" = Linux ]; then
  if [ "$(id -un)" = root ] && [ -d /root/astrate-mule ]; then WHERE=pi
  elif [ -d /home/atsetilam/astrate/.git ]; then WHERE=legion
  fi
fi
echo "running on: $WHERE"
```

Then run the Pi and Legion checks accordingly — locally when you're already there, over ssh
otherwise:

```sh
source .mule/hosts
# Pi: are the timers alive, when do they next fire
if [ "$WHERE" = pi ]; then
  systemctl list-timers mule.timer mule-survey.timer --no-pager
else
  ssh -o BatchMode=yes -o ConnectTimeout=6 "$MULE_PI_SSH" \
    'systemctl list-timers mule.timer mule-survey.timer --no-pager' \
    || echo "Pi unreachable from here — that's a real finding only when \$WHERE != pi; if you ARE on the Pi and this branch still ran, the WHERE detection above is wrong, fix it, don't just report 'unreachable'"
fi

# The queue itself: open mule-labelled issues (what Big Pickle will pick up next),
# mule-review (landed, waiting on you), mule-blocked, mule-alarm (dead-man's switch)
gh issue list --label mule --state open
gh issue list --label mule-review,mule-blocked,mule-alarm --state open

# Anything waiting on a design decision from Giulio
cat .mule/for-giulio.md   # if only the header/--- remains, it's empty — good

# Has the daily survey found anything not yet triaged
git fetch origin mule/research -q
git show origin/mule/research:.mule/research/triaged.md 2>/dev/null
git ls-tree -r --name-only origin/mule/research -- .mule/research | grep '^\.mule/research/survey-'

# Legion Go reachable right now? (optional muscle — absence is not an error)
if [ "$WHERE" = legion ]; then
  docker info >/dev/null 2>&1 && echo "legion: up (this IS the legion)" || echo "legion: docker not responding locally — odd, look into it"
else
  ssh -o BatchMode=yes -o ConnectTimeout=6 legion 'docker info >/dev/null 2>&1 && echo up' 2>/dev/null \
    || echo "legion: asleep/unreachable (fine, tasks tagged [legion] just skip)"
fi

# Trickle: is there an open architect-led plan waiting to resume (Mac-only — .trickle/ is
# excluded via .git/info/exclude and only exists on the Mac's working copy)
[ "$WHERE" = mac ] && [ -f .trickle/plan.md ] && tail -30 .trickle/plan.md
```

**A remote check that legitimately fails (ssh times out, `gh`/`git` has no network) is a real
finding — report it as "couldn't reach X from here", never silently as "X is down".** Those
are different facts and the first one is the far more common cause.

## The menu — "I want to..." → what actually happens

| Giulio says (examples, IT/EN) | What runs | Where | Attended? |
|---|---|---|---|
| "fai lavorare Big Pickle", "svuota la coda", "cos'ha combinato" | `mule` skill → `tools/mule.sh status`/`review` | Pi's clone, read from Mac | you read the diff before merge |
| "fai qualcosa di utile", "cerca cose da fare", "top up the queue" | `mule` skill → `tools/mule.sh recipe <name>` (rotates: github-issues, astarte-upstream, code-review, docs-sync, hygiene, milestones) — **proposes into `.mule/todo.md`, executes nothing** | Mac (recipe runs here, proposals land in git) | **you approve the proposal before it can run** |
| "dai un task a Big Pickle: ..." | `gh issue create --label mule --title "..."` (title = the whole task the mule reads — see `.mule/recipes/github-issues.md` for the shape a runnable one needs) | anywhere with `gh` | you write it, mule picks it up next tick |
| "merge quello che ha fatto il mulo" | review the diff yourself, then `git checkout main && git merge --no-ff mule/queue` | Mac | you decide, never automatic |
| "ferma/riattiva il mulo" | `ssh $MULE_PI_SSH 'systemctl stop\|start mule.timer'` | Pi | manual, instant |
| "cosa ha trovato la survey", "controlla la survey" | `mule-triage` skill | Mac reads `mule/research` branch | mostly autonomous — Claude verifies + files issues, reports to you |
| "rilancia la survey adesso" (fuori orario) | `ssh $MULE_PI_SSH 'systemctl start mule-survey.service'` (one-shot, doesn't wait for 03:00) — or run `tools/mule-survey.sh run` directly from the Mac if you want it faster/local | Pi (or Mac) | unattended once started; triage it after with `mule-triage` |
| "guarda cosa c'è per te" / review escalations | read `.mule/for-giulio.md` together, resolve each line in conversation, delete it once handled (that file is a queue, not a log) | Mac | **this is the one that's always your turn** — nothing files or clears these automatically |
| "delega questo", "trickle mode", "pianifica una fase" | `trickle` skill | Mac (worktrees), delegate types via opencode | architect (you+Claude) design, supervisor reviews every round |
| "riprendi il trickle" | `trickle` skill, reads `.trickle/plan.md`'s Status section for the resume point | Mac | same as above |
| "controlla il Legion", "è sveglio?" | `tools/mule.sh legion check`, or the probe in "Live status" above | Mac→Legion via ssh | read-only |
| "lancia un benchmark" | `.mule/recipes/benchmarks.md` — `bench/` CLI (`provision`/`ingest`/`connstorm`/`query`) against a named tier in `bench/scripts/tiers/`, run on/against the Legion Go | Legion Go, driven from wherever | you kick it off, results are yours to read |
| a `[legion]`-tagged mule task is stuck "blocked" | it's not stuck — it's *skipped* whenever the Legion Go is asleep, and retried next tick automatically | Pi | none needed |
| "cos'ha trovato l'audit upstream" (the routine one, not the daily survey) | that's `.mule/recipes/astarte-upstream.md`, run the same way as "cerca cose da fare" above but by name: `tools/mule.sh recipe astarte-upstream` | Mac | you approve proposals |
| "a che punto siamo con la milestone X", "controlla i milestone", "cosa serve per astrate 2.0/3.0" | `.mule/milestones.md` is the source of truth (release-tag-gated goals: v2.0 astarte-flow parity, v3.0+ CLEA pieces); `.mule/recipes/milestones.md` works the current one, and — unlike other recipes — its approved tasks **file GitHub issues directly** (`milestone-<tag>` label) instead of only proposing todo lines. Undecided scope (e.g. which CLEA piece is v3.0) always lands in `.mule/for-giulio.md`, never guessed | Mac (recipe runs here) / Pi (rotation, filed issues) | you approve proposals before they run; scope decisions are yours by design |
| "riconcilia il Pi/Legion", "salva lo stato di ...", a personal clone (not `astrate-mule`/`astrate-survey`, which timers already keep clean) might be dirty | `tools/reconcile.sh <path> [home-branch]` — commits everything onto a `wip/<host>-<timestamp>` branch with an explicit "not reviewed" note, pushes it, then fast-forwards `<path>` back to a clean `main`. Never merges the rescue branch itself. If it hits a signing failure (unattended clone, no gpg-agent), that's a real question — ask before setting `commit.gpgsign false` locally, don't assume | wherever the dirty clone is (run it there, or `ssh <host> bash -s < tools/reconcile.sh <path>`) | **note the rescue branch in `.mule/for-giulio.md`** so it doesn't get forgotten — the script itself doesn't write there |

## What's fully unattended right now (nothing to do unless something's wrong)

- `mule.timer` — every 30 min on the Pi, works the queue, commits to `mule/queue`, pushes.
  One task per tick, daily cap 16, does nothing when the queue is empty (correct, not broken).
- `mule-survey.timer` — daily ~03:00 on the Pi, own clone/branch (`mule/research`), writes
  nothing on days with no material change upstream or in our own code.
- Both post to `.mule/for-giulio.md` or open GitHub issues rather than ever touching `main`
  or closing anything themselves — the two places genuinely nothing-for-you-to-do can hide
  are those files/labels, which is why "live status" above checks them first.

## What always needs you (or Claude acting as your proxy in conversation)

- Anything in `.mule/for-giulio.md`.
- Merging `mule/queue` or a trickle PR into `main`.
- Approving a recipe's proposals before they become runnable queue tasks.
- Triaging the daily survey's candidates into real issues (delegable to `mule-triage`, but
  that skill still reports to you rather than filing silently forever unreviewed — read what
  it filed occasionally, same as `mule.sh review` for code).

## Keeping this map honest

If you add a new recipe, timer, or skill and this file doesn't mention it, that's a bug in
this file — say so and it should get a row.
