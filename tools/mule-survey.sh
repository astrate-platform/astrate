#!/usr/bin/env bash
# mule-survey.sh — a daily, separate-from-the-queue cross-project survey run.
#
# Different job from mule.sh's tick, deliberately kept in its own script rather than bolted
# onto it: this runs once a day (not every 30 min), it works its own branch (mule/research,
# not mule/queue) so it never collides with a queued task's commit, and its output is a
# report for a human/strong model to triage, not a diff a gate can wave through — it never
# files a GitHub issue itself. See .mule/recipes/cross-project-survey.md for what it asks
# opencode to actually do.
#
#   mule-survey.sh run          do it (git pull, run opencode, commit+push what it wrote)
#   mule-survey.sh preflight    environment check only
set -uo pipefail

REPO="$(git rev-parse --show-toplevel 2>/dev/null)" || { echo "not in a git repo" >&2; exit 1; }
MULE="$REPO/.mule"
RESEARCH="$MULE/research"
CONFIG="$MULE/config"

note() { printf '\033[36m>>\033[0m %s\n' "$*" >&2; }
bad()  { printf '\033[31m!!\033[0m %s\n' "$*" >&2; }
ok()   { printf '\033[32mok\033[0m %s\n' "$*" >&2; }
die()  { bad "$*"; exit 1; }

command -v opencode >/dev/null || die "opencode is not on PATH"
[ -f "$CONFIG" ] || die "no $CONFIG — this repo is not set up for mule mode"
# shellcheck disable=SC1090
. "$CONFIG"
SURVEY_BRANCH="${MULE_SURVEY_BRANCH:-mule/research}"
SURVEY_TIMEOUT="${MULE_SURVEY_TIMEOUT:-3600}"
MULE_AGENT="${MULE_AGENT:-build}"
MULE_MODEL="${MULE_MODEL:-}"

# Same watchdog shape as mule.sh's run_with_timeout: opencode exits 0 on SIGTERM, so a
# sentinel file is the only reliable way to tell "finished" from "killed".
run_with_timeout() {
  local secs="$1" sentinel="$2"; shift 2
  "$@" & local pid=$!
  ( sleep "$secs"; kill -0 "$pid" 2>/dev/null && { : > "$sentinel"; kill -TERM "$pid" 2>/dev/null; } ) \
    >/dev/null 2>&1 & local watchdog=$!
  wait "$pid"; local rc=$?
  kill -TERM "$watchdog" 2>/dev/null
  return "$rc"
}

ensure_branch() {
  local cur; cur="$(git -C "$REPO" rev-parse --abbrev-ref HEAD)"
  [ "$cur" = "$SURVEY_BRANCH" ] && return 0
  [ -z "$(git -C "$REPO" status --porcelain)" ] || \
    die "working tree is dirty on '$cur' — refusing to switch branches over it"
  git -C "$REPO" fetch -q origin "$SURVEY_BRANCH" 2>/dev/null || true
  if git -C "$REPO" show-ref --verify -q "refs/remotes/origin/$SURVEY_BRANCH"; then
    git -C "$REPO" checkout -q -B "$SURVEY_BRANCH" "origin/$SURVEY_BRANCH" \
      || die "cannot switch to $SURVEY_BRANCH"
  elif git -C "$REPO" show-ref --verify -q "refs/heads/$SURVEY_BRANCH"; then
    git -C "$REPO" checkout -q "$SURVEY_BRANCH" || die "cannot switch to $SURVEY_BRANCH"
  else
    git -C "$REPO" checkout -q -b "$SURVEY_BRANCH" || die "cannot create $SURVEY_BRANCH"
    note "created branch $SURVEY_BRANCH — the survey never commits to main or mule/queue"
  fi
}

cmd_preflight() {
  note "checking environment for the daily survey"
  opencode models >/dev/null 2>&1 && ok "opencode resolves models (agent '$MULE_AGENT')" \
    || bad "opencode cannot resolve a model — check auth"
  [ -z "$(git -C "$REPO" status --porcelain)" ] && ok "working tree clean" \
    || bad "working tree is dirty"
  git ls-remote https://github.com/astarte-platform/astarte HEAD >/dev/null 2>&1 \
    && ok "can reach github.com" || bad "cannot reach github.com — the survey needs it"
  df -h "$REPO" | tail -1 | awk '{print "  disk on "$6": "$4" free"}' >&2
}

cmd_run() {
  ensure_branch
  git -C "$REPO" merge -q --ff-only "origin/main" 2>/dev/null || true

  mkdir -p "$RESEARCH"
  local today; today="$(date +%Y-%m-%d)"
  local outlog="$RESEARCH/.last-run.log"
  local sentinel="$RESEARCH/.timeout"
  rm -f "$sentinel"

  local prompt="Read .mule/recipes/cross-project-survey.md in full, then do exactly what it
says. Today's date is $today — use it for any file you write
(.mule/research/survey-$today.md, .mule/research/issues-$today.md), and for checking what
has changed since the last run. This can take a long time (cloning two repos, reading docs);
that is expected, do not rush it. Do not touch git yourself and never run 'gh issue create' —
this script commits and pushes what you wrote after you finish. When done, finish with the
short status line the recipe asks for."

  note "running the survey (budget ${SURVEY_TIMEOUT}s) — this is slow, that's expected"
  local -a run=(opencode run --agent "$MULE_AGENT")
  [ -n "$MULE_MODEL" ] && run+=(--model "$MULE_MODEL")
  run+=("$prompt")

  ( cd "$REPO" && run_with_timeout "$SURVEY_TIMEOUT" "$sentinel" "${run[@]}" 2>&1 ) \
    | sed $'s/\033\\[[0-9;]*[A-Za-z]//g' | tee "$outlog"
  local rc=${PIPESTATUS[0]}

  # Disposable clones — never let these reach git, and don't leave them on a small disk.
  rm -rf "$RESEARCH/upstream"

  if [ -f "$sentinel" ]; then
    bad "survey timed out after ${SURVEY_TIMEOUT}s — nothing committed"
    rm -f "$sentinel"
    exit 4
  fi
  [ "$rc" = 0 ] || { bad "opencode exited $rc — nothing committed"; exit "$rc"; }

  local changed=0
  for f in "$RESEARCH/survey-$today.md" "$RESEARCH/issues-$today.md" "$RESEARCH/log.md"; do
    [ -f "$f" ] || continue
    git -C "$REPO" add -f "$f"
    changed=1
  done
  # `git status --cached` is not a thing — it errors out and prints nothing, so this test
  # used to be true on every run and the survey staged 16 reports between 2026-07-28 and
  # 2026-08-26 without ever committing one of them. Ask git about the index directly.
  if [ "$changed" = 0 ] || git -C "$REPO" diff --cached --quiet; then
    note "nothing to commit — the recipe found no material change today"
    return 0
  fi

  local msg="mule: daily cross-project survey, $today"
  if [ ! -f "$RESEARCH/survey-$today.md" ]; then
    msg="mule: daily cross-project survey, $today (no material change)"
  fi
  git -C "$REPO" commit -q -m "$msg"
  ok "committed to $SURVEY_BRANCH: $msg"

  if [ -n "${MULE_PUSH:-}" ]; then
    git -C "$REPO" push -q origin "$SURVEY_BRANCH" 2>/dev/null \
      && ok "pushed $SURVEY_BRANCH" || note "push failed — the work is safe locally"
  fi
}

case "${1:-run}" in
  preflight) cmd_preflight ;;
  run)       cmd_run ;;
  *)         die "usage: mule-survey.sh [preflight|run]" ;;
esac
