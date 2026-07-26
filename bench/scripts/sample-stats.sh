#!/usr/bin/env bash
# Samples docker container CPU/memory into a CSV while a benchmark runs:
#
#   ./scripts/sample-stats.sh results/astrate-stats.csv [interval_s]
#
# Run it in a second terminal (or backgrounded) around the ingest run, and
# stop it with Ctrl-C / kill. Sum the mem_bytes of a stack's containers for
# the footprint comparison; note the per-container breakdown too — the
# database vs application split is part of the story.
set -euo pipefail

out="${1:?usage: sample-stats.sh <out.csv> [interval_s]}"
interval="${2:-5}"

echo "ts,container,cpu_pct,mem_bytes,mem_limit_bytes" > "$out"
echo "sampling docker stats every ${interval}s → $out (Ctrl-C to stop)"

to_bytes() {
    # docker stats prints humanized sizes (e.g. 1.5GiB, 200MiB, 900kB).
    awk 'BEGIN{IGNORECASE=1}
    {
        v=$0
        mult=1
        if (v ~ /KiB/) mult=1024
        else if (v ~ /MiB/) mult=1024*1024
        else if (v ~ /GiB/) mult=1024*1024*1024
        else if (v ~ /kB/)  mult=1000
        else if (v ~ /MB/)  mult=1000*1000
        else if (v ~ /GB/)  mult=1000*1000*1000
        gsub(/[^0-9.]/, "", v)
        printf "%.0f", v*mult
    }'
}

while true; do
    ts="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    docker stats --no-stream --format '{{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}' |
    while IFS=$'\t' read -r name cpu mem; do
        used="${mem%%/*}"
        limit="${mem##*/}"
        printf '%s,%s,%s,%s,%s\n' \
            "$ts" "$name" "${cpu%\%}" \
            "$(printf '%s' "$used" | to_bytes)" \
            "$(printf '%s' "$limit" | to_bytes)" >> "$out"
    done
    sleep "$interval"
done
