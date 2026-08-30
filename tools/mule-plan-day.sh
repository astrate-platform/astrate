#!/usr/bin/env bash
# mule-plan-day.sh — write today's mule tick times into root's crontab.
#
# The mule used to fire on a flat 30-minute cadence, 24 hours a day, gated only by a daily
# tick count. That reads as a machine, not as someone sitting down to a side project between
# other things. This picks a random daily tick count and random minute-level times inside two
# such windows and installs them as one-shot cron entries for today, then gets out of the way.
# Run once a day, early, before either window opens — see mule-planner.timer.
set -euo pipefail

REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MULE_BIN="$REPO/tools/mule.sh"
LOG="$REPO/.mule/cron.log"
MARK_BEGIN="# BEGIN mule-daily-schedule"
MARK_END="# END mule-daily-schedule"

rand_between() { echo $(( $1 + RANDOM % ($2 - $1 + 1) )); }

DOM="$(date +%-d)"
MON="$(date +%-m)"

# 18-22 ticks a day, averaging 20 -> ~140/week.
total=$(rand_between 18 22)
# Windows are ~52min and ~245min long; split roughly in that proportion.
morning=$(rand_between 2 4)
evening=$(( total - morning ))

gen_times() {
  local start=$(( $1 * 60 + $2 )) end=$(( $3 * 60 + $4 )) span n=$5 i t
  span=$(( end - start ))
  for ((i = 0; i < n; i++)); do
    t=$(( start + RANDOM % span ))
    printf '%d %d\n' $(( t / 60 )) $(( t % 60 ))
  done
}

ENV="MULE_PUSH=1 HOME=/root PATH=/root/.opencode/bin:/usr/local/bin:/usr/bin:/bin"

lines=()
while read -r h m; do
  lines+=("$m $h $DOM $MON * $ENV MULE_DAILY_MAX=$total /bin/bash $MULE_BIN tick >> $LOG 2>&1")
done < <(gen_times 12 48 13 40 "$morning"; gen_times 18 42 22 47 "$evening")

tmp="$(mktemp)"
{ crontab -l 2>/dev/null || true; } | awk -v b="$MARK_BEGIN" -v e="$MARK_END" '
  $0==b {skip=1; next}
  $0==e {skip=0; next}
  skip!=1 {print}
' > "$tmp"

{
  cat "$tmp"
  echo "$MARK_BEGIN"
  printf '%s\n' "${lines[@]}"
  echo "$MARK_END"
} | crontab -
rm -f "$tmp"

echo "planned $total ticks for $(date +%F): $morning midday, $evening evening"
