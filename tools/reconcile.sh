#!/usr/bin/env bash
# reconcile.sh — make a possibly-dirty, unattended clone usable again without losing anything.
#
# For the clones nothing else is watching. mule.sh and mule-survey.sh already refuse to run
# over a dirty tree rather than auto-committing — that discipline stays exactly as it is, and
# this script does NOT run from either timer. It exists for the other kind of clone: a
# personal/ad-hoc checkout on the Pi or the Legion Go (root's /root/astrate on the Pi is the
# known example, flagged in .mule/for-giulio.md) that can quietly accumulate uncommitted work
# between sessions, with nobody reading its `git status` until something breaks.
#
# What it does, in order:
#   1. If the tree is dirty, commit EVERYTHING (including untracked) onto a fresh branch
#      named wip/<host>-<UTC timestamp>, with a commit message that says in its own words
#      this is a rescue, not reviewed work, and may not even build. Push it if a remote push
#      is possible.
#   2. Switch to <home-branch> (default: main) and fast-forward it from origin.
#   3. Report plainly what happened — a rescue branch name (and whether it pushed), or "clean".
#
# It never force-pushes, never discards anything, and never merges the rescue branch anywhere
# — that read is a human's job, same as reviewing any other unreviewed diff in this project.
#
#   reconcile.sh <path> [<home-branch>]
set -euo pipefail

REPO="${1:?usage: reconcile.sh <path> [<home-branch>]}"
HOME_BRANCH="${2:-main}"

note() { printf '\033[36m>>\033[0m %s\n' "$*" >&2; }
ok()   { printf '\033[32mok\033[0m %s\n' "$*" >&2; }
bad()  { printf '\033[31m!!\033[0m %s\n' "$*" >&2; }

[ -d "$REPO/.git" ] || { bad "$REPO is not a git repo"; exit 1; }
cd "$REPO"

HOST="$(hostname -s 2>/dev/null || hostname 2>/dev/null || echo unknown-host)"
STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
STARTED_ON="$(git rev-parse --abbrev-ref HEAD)"

if [ -n "$(git status --porcelain)" ]; then
  RESCUE="wip/${HOST}-${STAMP}"
  note "dirty tree on '$STARTED_ON' — rescuing onto $RESCUE"
  git checkout -q -b "$RESCUE"
  git add -A
  git commit -q -m "WIP: state rescued by reconcile.sh, NOT reviewed

Uncommitted work found on ${HOST}:${REPO} (was on branch '${STARTED_ON}') at ${STAMP}.
Committed as-is to save it from being lost or silently overwritten by the next pull —
this is not reviewed, may not build, and may mix unrelated changes. Read the diff
before merging anything from this branch anywhere."
  if git push -q -u origin "$RESCUE" 2>/dev/null; then
    ok "rescued and pushed: $RESCUE"
    echo "RESCUED $RESCUE pushed"
  else
    note "rescued locally, push failed (offline? no write access?) — branch stays local: $RESCUE"
    echo "RESCUED $RESCUE local-only"
  fi
else
  ok "tree was already clean on '$STARTED_ON'"
  echo "CLEAN"
fi

git checkout -q "$HOME_BRANCH" 2>/dev/null || git checkout -q -b "$HOME_BRANCH" "origin/$HOME_BRANCH"
git fetch -q origin "$HOME_BRANCH"
if git merge -q --ff-only "origin/$HOME_BRANCH" 2>/dev/null; then
  ok "$HOME_BRANCH fast-forwarded to origin/$HOME_BRANCH ($(git rev-parse --short HEAD))"
else
  bad "$HOME_BRANCH would not fast-forward against origin/$HOME_BRANCH — needs a human, left as-is"
  exit 2
fi
