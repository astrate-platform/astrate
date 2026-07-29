#!/usr/bin/env bash
# mule-setup-pi.sh — install the mule on the Raspberry Pi, run from your Mac.
#
# Idempotent: run it again after changing anything and it converges. It never touches
# /root/astrate (your own working copy on the Pi, which currently has uncommitted work) —
# the mule gets its own clone at /root/astrate-mule.
#
#   bash tools/mule-setup-pi.sh          # do it
#   bash tools/mule-setup-pi.sh check    # report only, change nothing
#
set -euo pipefail

PI="${MULE_PI_SSH:-root@192.168.0.96}"
LEGION="${MULE_LEGION_SSH:-atsetilam@192.168.0.212}"
LEGION_HOST="${LEGION#*@}"
REPO_URL="git@github.com:astrate-platform/astrate.git"
PI_REPO="/root/astrate-mule"
SSH_OPTS="-o BatchMode=yes -o ConnectTimeout=6"

note() { printf '\033[36m>>\033[0m %s\n' "$*"; }
ok()   { printf '\033[32mok\033[0m %s\n' "$*"; }
bad()  { printf '\033[31m!!\033[0m %s\n' "$*"; }
# shellcheck disable=SC2086
pish() { ssh $SSH_OPTS "$PI" 'bash -s' ; }   # the Pi's root shell is fish — always bash -s

MODE="${1:-install}"

# --- 1. reachability --------------------------------------------------------
# shellcheck disable=SC2086
ssh $SSH_OPTS "$PI" 'true' 2>/dev/null || { bad "cannot reach the Pi at $PI"; exit 1; }
ok "Pi reachable at $PI"

# --- 2. the Pi -> Legion Go link -------------------------------------------
# A dedicated key, used for nothing else, so revoking it is one line in one file on the
# Legion Go and nothing else breaks.
if ssh $SSH_OPTS "$PI" 'ssh -o BatchMode=yes -o ConnectTimeout=6 legion true' 2>/dev/null; then
  ok "Pi can already reach the Legion Go"
else
  note "the Pi cannot reach the Legion Go yet — setting that up"
  [ "$MODE" = check ] || {
    PUB="$(pish <<'EOS'
[ -f ~/.ssh/id_ed25519_mule_legion ] || \
  ssh-keygen -q -t ed25519 -f ~/.ssh/id_ed25519_mule_legion -N "" -C "mule@dietpi -> legion"
cat ~/.ssh/id_ed25519_mule_legion.pub
EOS
)"
    note "authorizing the Pi's key on the Legion Go"
    # shellcheck disable=SC2086
    ssh $SSH_OPTS "$LEGION" "mkdir -p ~/.ssh && chmod 700 ~/.ssh && \
      touch ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys && \
      grep -qxF '$PUB' ~/.ssh/authorized_keys || printf '%s\n' '$PUB' >> ~/.ssh/authorized_keys"

    note "trusting the Legion Go's host key on the Pi, and naming it 'legion'"
    pish <<EOS
ssh-keyscan -H $LEGION_HOST >> ~/.ssh/known_hosts 2>/dev/null
sort -u ~/.ssh/known_hosts -o ~/.ssh/known_hosts
grep -q '^Host legion\$' ~/.ssh/config 2>/dev/null || printf '\nHost legion\n  HostName $LEGION_HOST\n  User ${LEGION%%@*}\n  IdentityFile ~/.ssh/id_ed25519_mule_legion\n  StrictHostKeyChecking accept-new\n' >> ~/.ssh/config
EOS
    ssh $SSH_OPTS "$PI" 'ssh -o BatchMode=yes legion "echo reached \$(hostname)"' \
      && ok "Pi -> Legion Go works" || bad "Pi -> Legion Go still failing"
  }
fi

# --- 3. the mule's own clone -----------------------------------------------
[ "$MODE" = check ] || pish <<EOS
set -e
if [ ! -d $PI_REPO/.git ]; then git clone -q $REPO_URL $PI_REPO; fi
cd $PI_REPO
git fetch -q origin
# Fast-forward main, but only if that is where we are — never disturb mule/queue, which may
# hold landed work that has not been merged yet.
if [ "\$(git rev-parse --abbrev-ref HEAD)" = main ]; then
  git merge --ff-only -q origin/main 2>/dev/null || echo "WARN: main would not fast-forward"
fi
mkdir -p $PI_REPO/.mule
echo "clone at \$(git log -1 --format='%h %s') on \$(git rev-parse --abbrev-ref HEAD)"
EOS
ok "mule clone at $PI_REPO (your /root/astrate is left alone)"

# --- 3b. the hosts file -----------------------------------------------------
# .mule/hosts is deliberately NOT in git: the repo is public and there is no reason to
# publish a LAN layout and two usernames. So it is copied, and re-copied whenever it changes.
if [ -f "$(dirname "$0")/../.mule/hosts" ]; then
  [ "$MODE" = check ] || {
    # shellcheck disable=SC2086
    scp $SSH_OPTS -q "$(dirname "$0")/../.mule/hosts" "$PI:$PI_REPO/.mule/hosts"
    ok "copied .mule/hosts to the Pi (not committed — the repo is public)"
  }
else
  bad "no .mule/hosts here — copy .mule/hosts.example to .mule/hosts and fill it in"
fi

# --- 4. the schedule ---------------------------------------------------------
# mule.service stays around for manual/off-schedule runs (systemctl start mule.service), but
# ticks are no longer driven by a flat 30-minute systemd timer. A daily planner
# (mule-plan-day.sh) picks a random tick count and random times inside two windows and writes
# them as one-shot cron entries for the day; a small systemd timer runs the planner itself
# once a day, early, before either window opens.
[ "$MODE" = check ] || pish <<'EOS'
set -e
cat > /etc/systemd/system/mule.service <<'UNIT'
[Unit]
Description=Astrate mule — one queued task, if there is one
After=network-online.target
Wants=network-online.target

[Service]
Type=oneshot
WorkingDirectory=/root/astrate-mule
Environment=MULE_PUSH=1
Environment=HOME=/root
Environment=PATH=/root/.opencode/bin:/usr/local/bin:/usr/bin:/bin
ExecStart=/bin/bash /root/astrate-mule/tools/mule.sh tick
TimeoutStartSec=1800
Nice=10
UNIT

# Old flat 30-minute timer, if present from a previous install — cron now drives ticks.
systemctl disable --now mule.timer 2>/dev/null || true
rm -f /etc/systemd/system/mule.timer

cat > /etc/systemd/system/mule-planner.service <<'UNIT'
[Unit]
Description=Plan the Astrate mule's tick times for today
After=network-online.target

[Service]
Type=oneshot
WorkingDirectory=/root/astrate-mule
ExecStart=/bin/bash /root/astrate-mule/tools/mule-plan-day.sh
UNIT

cat > /etc/systemd/system/mule-planner.timer <<'UNIT'
[Unit]
Description=Pick the Astrate mule's tick times for today, once a day

[Timer]
OnCalendar=*-*-* 06:00:00
RandomizedDelaySec=1800
AccuracySec=5min
Persistent=true

[Install]
WantedBy=timers.target
UNIT

systemctl daemon-reload
systemctl enable --now mule-planner.timer
systemctl start mule-planner.service
systemctl list-timers mule-planner.timer --no-pager | head -3
crontab -l | grep -A30 'BEGIN mule-daily-schedule' || true
EOS
ok "daily schedule installed (mule-planner.timer + cron)"

# --- 5. the daily survey: its own clone, its own timer ----------------------
# A second, separate clone (not /root/astrate-mule) because the survey works its own branch
# (mule/research) and a `git checkout` there while a queued task has mule/queue checked out
# in the same working tree would collide. Two small clones on an SD card is cheaper than that
# coordination problem.
SURVEY_REPO="/root/astrate-survey"
[ "$MODE" = check ] || pish <<EOS
set -e
if [ ! -d $SURVEY_REPO/.git ]; then git clone -q $REPO_URL $SURVEY_REPO; fi
cd $SURVEY_REPO
git fetch -q origin
if [ "\$(git rev-parse --abbrev-ref HEAD)" = main ]; then
  git merge --ff-only -q origin/main 2>/dev/null || echo "WARN: main would not fast-forward"
fi
echo "survey clone at \$(git log -1 --format='%h %s') on \$(git rev-parse --abbrev-ref HEAD)"
EOS
ok "survey clone at $SURVEY_REPO"

# systemd, not cron, same reasoning as the mule timer. OnCalendar picks an off-peak hour so
# it never contends with a mule tick for the Pi's 4 cores / 3.7GB; RandomizedDelaySec spreads
# load if this pattern is ever reused across more than one host. Persistent=true means a Pi
# that was off at 03:00 still runs the survey once it's back, instead of skipping a whole day.
[ "$MODE" = check ] || pish <<EOS
set -e
cat > /etc/systemd/system/mule-survey.service <<UNIT
[Unit]
Description=Astrate daily cross-project survey (astarte-platform, AtomVM)
After=network-online.target
Wants=network-online.target

[Service]
Type=oneshot
WorkingDirectory=$SURVEY_REPO
Environment=MULE_PUSH=1
Environment=HOME=/root
Environment=PATH=/root/.opencode/bin:/usr/local/bin:/usr/bin:/bin
ExecStart=/bin/bash $SURVEY_REPO/tools/mule-survey.sh run
TimeoutStartSec=4000
Nice=15
UNIT

cat > /etc/systemd/system/mule-survey.timer <<'UNIT'
[Unit]
Description=Run the Astrate daily cross-project survey once a day

[Timer]
OnCalendar=*-*-* 03:00:00
RandomizedDelaySec=600
AccuracySec=5min
Persistent=true

[Install]
WantedBy=timers.target
UNIT

systemctl daemon-reload
systemctl enable --now mule-survey.timer
systemctl list-timers mule-survey.timer --no-pager | head -3
EOS
ok "daily survey timer installed and enabled (~03:00, +/- 10min)"

# --- 6. Claude Code's skills, usable from opencode too ----------------------
# opencode has its own skill loader (scans **/SKILL.md under configured directories) that
# reads the exact same file format as Claude Code — same required `name`/`description`
# frontmatter, same "one folder per skill" layout. So the skills written for Claude Code
# (~/.claude/skills on the Mac: astrate-dashboard, mule-triage) are usable from an interactive
# `opencode` session on the Pi too, without copying prose into a second, divergent file.
#
# NEVER commit these into the repo: historically astrate-dashboard's body named the Pi's and
# Legion's LAN addresses directly, the exact information .mule/hosts is kept out of git for.
# That's fixed now (LAN IPs replaced with $MULE_PI_SSH/$MULE_LEGION_SSH sourced from
# .mule/hosts, see RESTRUCTURING.md in agent-skills) — the skill's prose is safe to commit
# today, but this scp-and-never-commit deploy step hasn't been retired to match yet (deferred:
# touching the Pi's live setup needs asking Giulio first). Still travels by scp, like
# .mule/hosts, into a location outside any git working tree.
LOCAL_SKILLS="$HOME/.claude/skills"
SKILL_NAMES="astrate-dashboard mule-triage"

# jq merge helper: add path to config.skills.paths if not already present, preserving
# everything else in the file (the Pi's opencode.jsonc already carries its provider config).
jq_add_skill_path() {
  local file="$1" path="$2" tmp="${1}.tmp.$$"
  jq --arg p "$path" '.skills.paths = ((.skills.paths // []) + [$p] | unique)' "$file" > "$tmp" \
    && mv "$tmp" "$file"
}

if [ -d "$LOCAL_SKILLS/astrate-dashboard" ]; then
  [ "$MODE" = check ] || {
    # 6a. the Mac's own opencode: point straight at ~/.claude/skills, no copy, one source of truth.
    MAC_OC_CONFIG="$HOME/.config/opencode/opencode.jsonc"
    if [ -f "$MAC_OC_CONFIG" ] && command -v jq >/dev/null; then
      jq_add_skill_path "$MAC_OC_CONFIG" "$LOCAL_SKILLS"
      ok "Mac's opencode config points at $LOCAL_SKILLS"
    else
      bad "no $MAC_OC_CONFIG or no local jq — skipped the Mac side, do it by hand"
    fi

    # 6b. the Pi: scp the skill folders (root has no ~/.claude), merge its own config.
    pish <<'EOS'
mkdir -p ~/.claude-skills
EOS
    for name in $SKILL_NAMES; do
      [ -d "$LOCAL_SKILLS/$name" ] || continue
      # shellcheck disable=SC2086
      scp $SSH_OPTS -qr "$LOCAL_SKILLS/$name" "$PI:~/.claude-skills/$name.tmp" \
        && ssh $SSH_OPTS "$PI" "rm -rf ~/.claude-skills/$name && mv ~/.claude-skills/$name.tmp ~/.claude-skills/$name"
    done
    ok "copied $SKILL_NAMES to the Pi's ~/.claude-skills (not committed — see the note above)"

    pish <<'EOS'
set -e
CFG=/root/.config/opencode/opencode.jsonc
command -v jq >/dev/null || { echo "no jq on the Pi — cannot merge $CFG"; exit 1; }
[ -f "$CFG" ] || echo '{"$schema":"https://opencode.ai/config.json"}' > "$CFG"
jq '.skills.paths = ((.skills.paths // []) + ["/root/.claude-skills"] | unique)' "$CFG" > "$CFG.tmp.$$" \
  && mv "$CFG.tmp.$$" "$CFG"
echo "Pi's opencode config now has skills.paths: $(jq -c .skills.paths "$CFG")"
EOS
    ok "Pi's opencode config points at ~/.claude-skills"
  }
else
  bad "no $LOCAL_SKILLS/astrate-dashboard here — skipping the skills step (run this from the Mac that has them)"
fi

cat <<EOF

  Installed. From here on:

    ssh $PI 'crontab -l'                                     # today's planned tick times
    ssh $PI 'journalctl -u mule.service -n 50'              # what the last tick did
    ssh $PI 'systemctl stop mule-planner.timer'             # stop planning new days
    ssh $PI 'cd $PI_REPO && ./tools/mule.sh status'

    ssh $PI 'systemctl list-timers mule-survey.timer'       # when the daily survey next fires
    ssh $PI 'journalctl -u mule-survey.service -n 80'       # what it found
    ssh $PI 'systemctl stop mule-survey.timer'              # make it stop
    git fetch origin mule/research && git log origin/mule/research  # read what it wrote

    ssh $PI                                                  # then just: opencode
                                                               # /astrate-dashboard, /mule-triage work there too now

  The mule does nothing at all until there are approved tasks in .mule/todo.md or issues
  labelled 'mule' — an idle mule is the correct mule. The daily survey runs
  regardless, but is itself incremental: most days it should report "no material change" and
  commit nothing you need to read. Its output on branch mule/research is a report for you or
  a strong model to triage into issues — it never files one itself.
EOF
