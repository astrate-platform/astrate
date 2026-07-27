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

# --- 4. the timer -----------------------------------------------------------
# systemd, not cron: OnUnitInactiveSec measures from the previous run *finishing*, so a task
# that overruns 30 minutes delays the next tick instead of racing it. flock inside `tick` is
# the belt to that braces.
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

cat > /etc/systemd/system/mule.timer <<'UNIT'
[Unit]
Description=Run the Astrate mule every 30 minutes

[Timer]
# Measured from when the last run FINISHED, so runs can never overlap or pile up.
OnBootSec=10min
OnUnitInactiveSec=30min
AccuracySec=1min
Persistent=false

[Install]
WantedBy=timers.target
UNIT

systemctl daemon-reload
systemctl enable --now mule.timer
systemctl list-timers mule.timer --no-pager | head -3
EOS
ok "timer installed and enabled"

cat <<EOF

  Installed. From here on:

    ssh $PI 'systemctl list-timers mule.timer'      # when it next fires
    ssh $PI 'journalctl -u mule.service -n 50'      # what it did
    ssh $PI 'systemctl stop mule.timer'             # make it stop
    ssh $PI 'cd $PI_REPO && ./tools/mule.sh status'

  It will do nothing at all until there are approved tasks in .mule/todo.md on the
  branch it tracks. That is deliberate: an idle mule is the correct mule.
EOF
