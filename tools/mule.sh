#!/usr/bin/env bash
# mule.sh — drive Big Pickle (opencode) as a solo, self-paced worker on this repo.
#
# The difference from trickle.sh: there is no architect and no supervisor. The mule works
# its own queue, one task per *fresh opencode process*, and the safety comes from the gates
# and from git rather than from a model reading the diff.
#
# The one design constraint everything here follows: big-pickle has ~200k of window and
# ~30k of *good* window. So the loop lives in bash, not in the model's context. Every task
# starts from an empty session that reads two short files and nothing else.
#
#   mule.sh preflight        environment check — run once per sitting
#   mule.sh next             run the first open task in .mule/todo.md
#   mule.sh loop [N]         run up to N tasks (default: until the queue is empty)
#   mule.sh status           queue, branch, last few log rows
#   mule.sh add "<title>"    append a task to the queue
#   mule.sh menu             print the "what next?" options for the user
#   mule.sh recipe <name>    run a recipe: proposes new tasks, executes none
#   mule.sh revert           undo the last mule commit
#
set -uo pipefail

REPO="$(git rev-parse --show-toplevel 2>/dev/null)" || { echo "not in a git repo" >&2; exit 1; }
MULE="$REPO/.mule"
CONFIG="$MULE/config"
TODO="$MULE/todo.md"
LOG="$MULE/log.md"

note() { printf '\033[36m>>\033[0m %s\n' "$*" >&2; }
bad()  { printf '\033[31m!!\033[0m %s\n' "$*" >&2; }
ok()   { printf '\033[32mok\033[0m %s\n' "$*" >&2; }
die()  { bad "$*"; exit 1; }

command -v opencode >/dev/null || die "opencode is not on PATH"
[ -f "$CONFIG" ] || die "no $CONFIG — this repo is not set up for mule mode"
# shellcheck disable=SC1090
. "$CONFIG"
# shellcheck disable=SC1091
[ -f "$MULE/hosts" ] && . "$MULE/hosts"
MULE_BRANCH="${MULE_BRANCH:-mule/queue}"
MULE_TIMEOUT="${MULE_TIMEOUT:-900}"
MULE_AGENT="${MULE_AGENT:-build}"
MULE_MODEL="${MULE_MODEL:-}"

# --- helpers ----------------------------------------------------------------

# macOS has no timeout(1). Same shape as trickle.sh: watchdog output to /dev/null (it
# inherits our stdout otherwise and holds the pipe open for the full budget), and a sentinel
# file, because opencode exits 0 on SIGTERM and a killed run is otherwise indistinguishable
# from a clean one.
#
# Both background jobs close fd 9 (`9>&-`). fd 9 is tick's flock, and a child that inherits
# it holds the lock as long as it lives — the watchdog is a `sleep $secs`, so an orphaned one
# kept the lock for a full 20 minutes after its tick had finished, and every tick in that
# window was skipped with "previous run is still going" while nothing was running. The timer
# looks alive and does nothing, which is the worst way for this to fail.
run_with_timeout() {
  local secs="$1" sentinel="$2"; shift 2
  "$@" 9>&- & local pid=$!
  ( sleep "$secs"; kill -0 "$pid" 2>/dev/null && { : > "$sentinel"; kill -TERM "$pid" 2>/dev/null; } ) \
    9>&- >/dev/null 2>&1 & local watchdog=$!
  wait "$pid"; local rc=$?
  kill -TERM "$watchdog" 2>/dev/null
  return "$rc"
}

ensure_branch() {
  local cur; cur="$(git -C "$REPO" rev-parse --abbrev-ref HEAD)"
  [ "$cur" = "$MULE_BRANCH" ] && return 0
  [ -z "$(git -C "$REPO" status --porcelain)" ] || \
    die "working tree is dirty on '$cur' — commit or stash before starting the mule"
  if git -C "$REPO" show-ref --verify -q "refs/heads/$MULE_BRANCH"; then
    git -C "$REPO" checkout -q "$MULE_BRANCH" || die "cannot switch to $MULE_BRANCH"
  else
    git -C "$REPO" checkout -q -b "$MULE_BRANCH" || die "cannot create $MULE_BRANCH"
    note "created branch $MULE_BRANCH — the mule never commits to main"
  fi
}

# --- the Legion Go ----------------------------------------------------------
#
# The Legion Go is optional muscle, never a dependency. Every call here fails soft: an
# unreachable Legion means legion-tagged tasks are skipped and the rest of the queue runs.
# A timer that stalls whenever a handheld is asleep is a timer that does nothing.

legion_sh() {
  [ -n "${MULE_LEGION_SSH:-}" ] || return 1
  # shellcheck disable=SC2086
  ssh $MULE_SSH_OPTS "$MULE_LEGION_SSH" "$@"
}

legion_up() {
  [ -n "${MULE_LEGION_SSH:-}" ] || return 1
  legion_sh "${MULE_LEGION_PROBE:-true}" >/dev/null 2>&1
}

cmd_legion() {
  case "${1:-check}" in
    check)
      legion_sh 'echo up' >/dev/null 2>&1 \
        || { bad "Legion Go unreachable at ${MULE_LEGION_SSH:-<unset>}"; return 1; }
      ok "ssh works: $(legion_sh 'hostname; uname -m' 2>/dev/null | tr '\n' ' ')"
      if legion_sh 'docker info >/dev/null 2>&1'; then
        ok "docker is up ($(legion_sh 'docker ps -q | wc -l' 2>/dev/null | tr -d ' ') containers)"
      else bad "docker is not up — heavy tasks cannot run"; return 1; fi
      legion_up && ok "[legion] tasks are runnable" || { bad "probe failed"; return 1; }
      ;;
    sh) shift; legion_sh "$@";;
    # The Legion Go has no Go toolchain, and putting one there is a machine-state decision
    # nobody made. bench/ imports nothing from Astrate and needs no cgo, so cross-compiling
    # a static binary here and copying it over costs nothing and leaves that machine alone.
    bench-push)
      legion_up || die "Legion Go is not available"
      local out="$MULE/.bench-linux-amd64"
      note "cross-compiling bench for linux/amd64"
      ( cd "$REPO/bench" && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$out" . ) \
        || die "cross-compile failed"
      legion_sh "mkdir -p '$MULE_LEGION_REPO/bench'" || die "mkdir failed on the Legion Go"
      # shellcheck disable=SC2086
      scp $MULE_SSH_OPTS "$out" "$MULE_LEGION_SSH:$MULE_LEGION_REPO/bench/bench" \
        || die "scp failed"
      legion_sh "chmod +x '$MULE_LEGION_REPO/bench/bench'"
      ok "bench binary is on the Legion Go at $MULE_LEGION_REPO/bench/bench"
      ;;
    *) die "legion: check | sh <cmd> | bench-push";;
  esac
}

# The mule may not touch these even though it has a whole worktree. Same list as trickle's,
# for the same reason: these files carry decisions or Giulio's voice, not typing.
check_never() {
  [ -n "${MULE_NEVER:-}" ] || return 0
  local changed hit=0
  changed="$(git -C "$REPO" status --porcelain | awk '{print $NF}')"
  local f pat
  for f in $changed; do
    for pat in $MULE_NEVER; do
      # shellcheck disable=SC2254
      case "$f" in $pat) bad "never-touch violation: $f (matches '$pat')"; hit=1;; esac
    done
  done
  return $hit
}

gates() {
  local rc=0
  if [ -n "${MULE_FIX_CMD:-}" ]; then ( cd "$REPO" && eval "$MULE_FIX_CMD" ) >/dev/null 2>&1; fi
  if [ -n "${MULE_TEST_CMD:-}" ]; then
    note "gate: $MULE_TEST_CMD"
    ( cd "$REPO" && eval "$MULE_TEST_CMD" ) >/dev/null 2>&1 || { bad "tests failed"; rc=1; }
  fi
  if [ "$rc" = 0 ] && [ -n "${MULE_LINT_CMD:-}" ]; then
    note "gate: lint"
    ( cd "$REPO" && eval "$MULE_LINT_CMD" ) >/dev/null 2>&1 || { bad "lint failed"; rc=1; }
  fi
  return $rc
}

log_row() { printf '| %s | %s | %s | %ss | %s |\n' "$(date +%F)" "$1" "$2" "$3" "${4:-}" >> "$LOG"; }

# --- preflight --------------------------------------------------------------

cmd_preflight() {
  local fail=0
  git -C "$REPO" rev-parse HEAD >/dev/null 2>&1 && ok "git repo at $REPO" || { bad "git"; fail=1; }
  if opencode models >/dev/null 2>&1; then ok "opencode resolves models (agent '$MULE_AGENT')"
  else bad "opencode cannot list models — provider outage or bad install"; fail=1; fi
  if [ -n "${MULE_TEST_CMD:-}" ]; then
    if ( cd "$REPO" && eval "$MULE_TEST_CMD" ) >/dev/null 2>&1; then ok "test baseline is green"
    else bad "test baseline is ALREADY RED — fix it first, or every task will look failed"; fail=1; fi
  else bad "MULE_TEST_CMD is empty — the mule would have no gate at all"; fail=1; fi
  if [ -n "${MULE_LINT_CMD:-}" ]; then
    ( cd "$REPO" && eval "$MULE_LINT_CMD" ) >/dev/null 2>&1 \
      && ok "lint baseline is clean" \
      || bad "lint baseline is not clean (clear the linter cache before believing this)"
  fi
  # Signing is fatal for an unattended mule and fails in the least helpful way: the gates
  # pass, the commit dies on a pinentry that has no tty, the work stays uncommitted, and the
  # dirty tree it leaves behind makes every future tick abort. Cost one run to learn.
  if [ "$(git -C "$REPO" config --get commit.gpgsign)" = "true" ]; then
    if git -C "$REPO" commit --dry-run --allow-empty -q -m probe >/dev/null 2>&1; then
      ok "commit signing is on and works unattended"
    else
      bad "commit.gpgsign is ON and cannot sign without a tty — every task would fail at"
      bad "  the commit and leave a dirty tree. Fix: git -C $REPO config commit.gpgsign false"
      fail=1
    fi
  else ok "commit signing is off (correct for an unattended clone)"; fi
  command -v gh >/dev/null && gh auth status >/dev/null 2>&1 \
    && ok "gh is authenticated (issue tasks will work)" \
    || note "gh not authenticated — the GitHub-issues recipe will not work"
  # A dirty tree is fatal, not cosmetic: ensure_branch refuses to switch branches over it,
  # and a task's diff would be indistinguishable from whatever was already uncommitted.
  if [ -z "$(git -C "$REPO" status --porcelain)" ]; then ok "working tree is clean"
  else bad "working tree is dirty — commit or stash first"; fail=1; fi
  echo
  [ "$fail" = 0 ] && ok "preflight passed" || bad "preflight FAILED — do not start the mule"
  return $fail
}

# --- the queue --------------------------------------------------------------

# A task is one `- [ ] slug: title` line in .mule/todo.md. If .mule/tasks/<slug>.md exists
# it is the spec; otherwise the line itself is the whole spec. Both are legitimate — a
# one-liner is enough for most of what the mule should be doing.
# The first open task the *current environment* can actually run. A line tagged [legion] is
# skipped when the Legion Go is asleep rather than blocking everything behind it — the queue
# is a queue, not a chain. The probe costs one ssh round-trip and only fires when a
# legion-tagged line is genuinely next in line.
first_open() {
  local legion_ok="" l
  while IFS= read -r l; do
    case "$l" in
      *"[legion]"*)
        if [ -z "$legion_ok" ]; then
          if legion_up; then legion_ok=yes
          else legion_ok=no; note "Legion Go unreachable — skipping [legion] tasks this run"; fi
        fi
        [ "$legion_ok" = yes ] || continue
        ;;
    esac
    printf '%s\n' "$l"; return 0
  done < <(grep -n '^- \[ \] ' "$TODO")
  return 1
}

cmd_status() {
  echo "branch:  $(git -C "$REPO" rev-parse --abbrev-ref HEAD)"
  echo "open:    $(grep -c '^- \[ \] ' "$TODO" 2>/dev/null | head -1)"
  echo "blocked: $(grep -c '^- \[!\] ' "$TODO" 2>/dev/null | head -1)"
  echo; echo "-- queue --"; grep -E '^- \[[ !]\] ' "$TODO" 2>/dev/null | head -15
  echo; echo "-- last log rows --"; tail -5 "$LOG" 2>/dev/null
}

cmd_add() {
  [ $# -gt 0 ] || die 'add needs a title: mule.sh add "slug: what to do"'
  printf -- '- [ ] %s\n' "$*" >> "$TODO"
  ok "queued: $*"
}

cmd_next() {
  ensure_branch
  local line; line="$(first_open)"
  [ -n "$line" ] || { note "queue is empty"; cmd_menu; return 9; }
  local lineno="${line%%:*}" text="${line#*:}"
  # Quote the pattern: unquoted, `[ ]` is a glob character class, so the prefix survives
  # and ends up in the slug, the commit message and the task file.
  text="${text#"- [ ] "}"
  local slug="${text%%:*}"; [ "$slug" = "$text" ] && slug="$(echo "$text" | tr -cs '[:alnum:]' '-' | cut -c1-24)"
  slug="$(echo "$slug" | tr -cd '[:alnum:]-_')"

  local spec="$MULE/tasks/$slug.md"
  local taskfile="$MULE/task.md"
  { echo "# Task"; echo; echo "$text"; echo
    if [ -f "$spec" ]; then echo "## Spec"; echo; cat "$spec"; fi
  } > "$taskfile"

  note "task: $text"
  local prompt='Read .mule/MULE.md, then read .mule/task.md, then do exactly that one task.

Do not touch git: no commit, no branch, no checkout, no stash. Leave your work in the
working tree — committing it is the script'"'"'s job, and it only commits if the gates pass.

Stop when the task is done. Do not start the next one. Finish with a short report in the
format MULE.md gives.'

  local -a run=(opencode run --agent "$MULE_AGENT")
  [ -n "$MULE_MODEL" ] && run+=(--model "$MULE_MODEL")
  run+=("$prompt")

  local sentinel="$MULE/.timeout"; rm -f "$sentinel"
  local started; started="$(date +%s)"
  ( cd "$REPO" && run_with_timeout "$MULE_TIMEOUT" "$sentinel" "${run[@]}" 2>&1 ) \
    | sed $'s/\033\\[[0-9;]*[A-Za-z]//g' | cat -s
  local rc=${PIPESTATUS[0]} elapsed=$(( $(date +%s) - started ))
  local timed_out=0; [ -f "$sentinel" ] && timed_out=1
  rm -f "$sentinel" "$taskfile"

  local verdict="" reason=""
  if [ "$timed_out" = 1 ]; then
    verdict=blocked; reason="TIMEOUT after ${elapsed}s — task too big, split it"
  elif [ "$rc" -ne 0 ]; then
    verdict=blocked; reason="opencode exited $rc"
  elif ! check_never; then
    verdict=blocked; reason="touched a never-touch path"
  elif [ -z "$(git -C "$REPO" status --porcelain)" ]; then
    verdict=blocked; reason="wrote nothing"
  elif ! gates; then
    verdict=blocked; reason="gates failed"
  else
    verdict=done
  fi

  if [ "$verdict" = done ]; then
    git -C "$REPO" add -A
    git -C "$REPO" commit -q -m "mule: $text" || { bad "commit failed"; return 1; }
    # mark the line done, in the commit that follows
    sed -i '' "${lineno}s/^- \[ \]/- [x]/" "$TODO"
    log_row "$slug" done "$elapsed" "$(git -C "$REPO" rev-parse --short HEAD)"
    git -C "$REPO" add "$TODO" "$LOG" && git -C "$REPO" commit -q -m "mule: log $slug"
    if [ -n "${MULE_PUSH:-}" ]; then
      git -C "$REPO" push -q origin "$MULE_BRANCH" 2>/dev/null \
        && ok "pushed $MULE_BRANCH" || note "push failed — the work is safe locally"
    fi
    ok "landed: $text (${elapsed}s)"
    return 0
  fi

  # Keep the fragment where it can be looked at, then put the tree back.
  mkdir -p "$MULE/failed"
  git -C "$REPO" diff > "$MULE/failed/$slug.diff" 2>/dev/null
  git -C "$REPO" checkout -- . 2>/dev/null; git -C "$REPO" clean -fdq -e .mule 2>/dev/null
  sed -i '' "${lineno}s/^- \[ \]/- [!]/" "$TODO"
  sed -i '' "${lineno}s|\$| — BLOCKED: $reason|" "$TODO"
  log_row "$slug" "blocked" "$elapsed" "$reason"
  git -C "$REPO" add "$TODO" "$LOG" >/dev/null 2>&1 && git -C "$REPO" commit -q -m "mule: blocked $slug"
  bad "blocked: $reason (fragment kept at .mule/failed/$slug.diff)"
  return 1
}

# The timer's entry point. Everything about it is deliberately timid, because the thing on
# the other end is a free provider and nobody is watching this run:
#
#   - **one task per tick, never a loop.** A tick that drains a queue is a tick that makes
#     twenty provider calls in a burst.
#   - **an empty queue does nothing at all.** It never runs a recipe to invent work for
#     itself. Work comes from Giulio approving a queue; an idle mule is the correct mule.
#   - **a daily cap**, so a runaway queue cannot spend the whole day hammering the provider.
#   - **flock, non-blocking**, so a slow task never overlaps the next tick — the ticks pile
#     up as no-ops instead of as concurrent opencode processes on a 4-core Pi.
cmd_tick() {
  command -v flock >/dev/null || die "tick needs flock(1) — run it on the Pi, not on macOS"
  exec 9>"$MULE/.lock"
  flock -n 9 || { note "previous run is still going — skipping this tick"; exit 0; }

  local today count=0 d c stamp="$MULE/.budget"
  today="$(date +%F)"
  if [ -f "$stamp" ]; then read -r d c < "$stamp"; [ "$d" = "$today" ] && count="$c"; fi
  if [ "$count" -ge "${MULE_DAILY_MAX:-16}" ]; then
    note "daily budget spent ($count/${MULE_DAILY_MAX:-16}) — idling until tomorrow"; exit 0
  fi

  # Start from what is actually on origin, so a task is never written against a stale tree
  # and Giulio's own pushes are picked up without anyone logging into the Pi.
  git -C "$REPO" fetch -q origin 2>/dev/null || note "fetch failed — working offline"

  first_open >/dev/null 2>&1 || { note "queue empty — nothing to do"; exit 0; }
  echo "$today $((count + 1))" > "$stamp"
  note "tick $((count + 1))/${MULE_DAILY_MAX:-16} for $today"
  cmd_next
}

cmd_loop() {
  local budget="${1:-99}" i=0 fails=0
  while [ "$i" -lt "$budget" ]; do
    cmd_next; local rc=$?
    [ "$rc" = 9 ] && return 0                       # queue drained; menu already printed
    [ "$rc" != 0 ] && fails=$((fails+1))
    # Two consecutive blocks means something systemic — stop rather than burn the queue.
    [ "$fails" -ge 2 ] && { bad "two blocked tasks — stopping, a human should look"; return 1; }
    [ "$rc" = 0 ] && fails=0
    i=$((i+1))
    echo
  done
  note "ran $i tasks; $(grep -c '^- \[ \] ' "$TODO" | head -1) still open"
}

cmd_revert() {
  # Skip the bookkeeping commits — "revert" means undo a task's *code*, and the log rows
  # around it are worth keeping as the record that the task was tried.
  local sha; sha="$(git -C "$REPO" log -20 --pretty='%H %s' \
    | grep -v ' mule: log ' | grep -v ' mule: blocked ' | grep -m1 ' mule: ' | cut -d' ' -f1)"
  [ -n "$sha" ] || die "no mule code commit in the last 20 commits"
  local subj; subj="$(git -C "$REPO" log -1 --pretty=%s "$sha")"
  git -C "$REPO" revert --no-edit "$sha" && ok "reverted: $subj"
}

cmd_menu() {
  cat <<'EOF'

  ─────────────────────────────────────────────────────────────
   The queue is empty. What next?

   1  issues     work through the open GitHub issues
   2  upstream   check astarte-platform upstream for new tags and
                 conceptual improvements worth porting
   3  review     read the codebase and propose features, perf work
                 or clarity cleanups
   4  bench      run the tiered benchmarks (small/medium/big/giant)
   5  docs       re-sync the docs site with what the code now does
   6  hygiene    dependency, CI and lint upkeep
   other         tell me what you want instead

   Each option is a recipe in .mule/recipes/. Running one PROPOSES
   tasks — it appends them to .mule/todo.md for you to approve
   before the mule executes them.

     tools/mule.sh add "<slug>: <what to do>"
     tools/mule.sh loop

   Heavy runs (option 4) go to the Lenovo Legion Go — tag those
   task lines [legion] and they are skipped while it is asleep.
   Set its BIOS UMA VRAM to 3GB (the lowest option) first.
  ─────────────────────────────────────────────────────────────
EOF
  if [ -n "${MULE_LEGION_SSH:-}" ]; then
    if legion_up; then echo "   Legion Go: UP — [legion] tasks will run"
    else echo "   Legion Go: down — [legion] tasks will wait"; fi
  fi
}

# Run a recipe: a fresh mule session that reads the recipe and appends proposed tasks.
cmd_recipe() {
  local name="${1:-}"; [ -n "$name" ] || die "recipe needs a name (see mule.sh menu)"
  case "$name" in
    1|issues)   name=github-issues;;
    2|upstream) name=astarte-upstream;;
    3|review)   name=code-review;;
    4|bench)    name=benchmarks;;
    5|docs)     name=docs-sync;;
    6)          name=hygiene;;
  esac
  local f="$MULE/recipes/$name.md"; [ -f "$f" ] || die "no such recipe: $f"
  ensure_branch
  local prompt="Read .mule/MULE.md, then read .mule/recipes/$name.md and carry it out.

That recipe is a *proposal* job: its output is new task lines appended to .mule/todo.md,
plus any evidence file it tells you to write. Do not start implementing the tasks you
propose. Do not touch git."
  local -a run=(opencode run --agent "$MULE_AGENT")
  [ -n "$MULE_MODEL" ] && run+=(--model "$MULE_MODEL")
  run+=("$prompt")
  local sentinel="$MULE/.timeout"; rm -f "$sentinel"
  ( cd "$REPO" && run_with_timeout "$MULE_TIMEOUT" "$sentinel" "${run[@]}" 2>&1 ) \
    | sed $'s/\033\\[[0-9;]*[A-Za-z]//g' | cat -s
  [ -f "$sentinel" ] && { rm -f "$sentinel"; bad "recipe timed out"; return 1; }
  echo; note "proposed queue is now:"; grep -E '^- \[ \] ' "$TODO" | tail -20
  note "edit .mule/todo.md to approve/cut, then: tools/mule.sh loop"
}

case "${1:-}" in
  preflight) shift; cmd_preflight "$@";;
  next)      shift; cmd_next "$@";;
  loop)      shift; cmd_loop "$@";;
  status)    shift; cmd_status "$@";;
  add)       shift; cmd_add "$@";;
  menu)      shift; cmd_menu "$@";;
  recipe)    shift; cmd_recipe "$@";;
  legion)    shift; cmd_legion "$@";;
  tick)      shift; cmd_tick "$@";;
  revert)    shift; cmd_revert "$@";;
  *) sed -n '2,19p' "$0" | sed 's/^# \{0,1\}//'; exit 1;;
esac
