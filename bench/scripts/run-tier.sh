#!/usr/bin/env bash
# Runs a named tier against a deployment and writes all results to a single
# timestamped directory.  Sourced tier env vars drive the device counts and
# rates; the base URL and housekeeping key are passed on the command line.
#
#   bench/scripts/run-tier.sh <tier> <target> -base-url <url> -housekeeping-key <file>
#
# <tier>   small | medium | big | giant  (must have a matching tiers/<tier>.env)
# <target> label for the system under test, e.g. "astrate" or "astarte"
#
# Resumable: if the state file for this tier already has enough devices,
# provision is skipped rather than re-registering tens of thousands of them.
# Never overwrites an existing results directory — the benchmark result is
# evidence, and a duplicate means something is wrong.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BENCH_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# ── parse args ──────────────────────────────────────────────────────────────
TIER="${1:?usage: run-tier.sh <tier> <target> -base-url <url> -housekeeping-key <file>}"
TARGET="${2:?usage: run-tier.sh <tier> <target> -base-url <url> -housekeeping-key <file>}"
shift 2

BASE_URL=""
HK_KEY=""
while [[ $# -gt 0 ]]; do
    case "$1" in
        -base-url)     BASE_URL="$2"; shift 2 ;;
        -housekeeping-key) HK_KEY="$2"; shift 2 ;;
        *) echo "run-tier.sh: unknown arg $1" >&2; exit 2 ;;
    esac
done
if [[ -z "$BASE_URL" || -z "$HK_KEY" ]]; then
    echo "run-tier.sh: -base-url and -housekeeping-key are required" >&2
    exit 2
fi

TIER_ENV="$SCRIPT_DIR/tiers/${TIER}.env"
if [[ ! -f "$TIER_ENV" ]]; then
    echo "run-tier.sh: tier file not found: $TIER_ENV" >&2
    exit 1
fi
# shellcheck source=tiers/small.env
source "$TIER_ENV"

echo "tier=$TIER target=$TARGET devices=$DEVICES rate=$RATE duration=$INGEST_DURATION"

# ── build bench once ────────────────────────────────────────────────────────
echo "building bench…"
BENCH_BIN="$BENCH_DIR/bench"
go build -o "$BENCH_BIN" "$BENCH_DIR"
trap 'rm -f "$BENCH_BIN"' EXIT

# ── results directory ───────────────────────────────────────────────────────
TS="$(date -u +%Y%m%d-%H%M%S)"
RESULTS_DIR="$BENCH_DIR/results/${TIER}-${TARGET}-${TS}"
if [[ -d "$RESULTS_DIR" ]]; then
    echo "run-tier.sh: results directory already exists — refusing to overwrite:" >&2
    echo "  $RESULTS_DIR" >&2
    exit 1
fi
mkdir -p "$RESULTS_DIR"

# ── host info ───────────────────────────────────────────────────────────────
{
    echo "=== host ==="
    uname -a
    echo ""
    echo "=== cpu ==="
    nproc
    echo ""
    echo "=== memory ==="
    free -h 2>/dev/null || vm_stat 2>/dev/null || echo "(unknown)"
} > "$RESULTS_DIR/host.txt"

cp "$TIER_ENV" "$RESULTS_DIR/"

# ── background stats sampler ────────────────────────────────────────────────
STATS_CSV="$RESULTS_DIR/stats.csv"
if [[ -x "$SCRIPT_DIR/sample-stats.sh" ]]; then
    "$SCRIPT_DIR/sample-stats.sh" "$STATS_CSV" 5 &
    STATS_PID=$!
    trap 'kill "$STATS_PID" 2>/dev/null; rm -f "$BENCH_BIN"' EXIT
else
    echo "note: sample-stats.sh not found or not executable — skipping resource sampling"
fi

# ── provision ───────────────────────────────────────────────────────────────
# Skip if a state file already exists with enough devices (resumable).
STATE="$RESULTS_DIR/state.json"
if [[ -f "$STATE" ]]; then
    existing=$(python3 -c "import json; print(len(json.load(open('$STATE'))['devices']))" 2>/dev/null || echo 0)
    if (( existing >= DEVICES )); then
        echo "provision: state file has $existing devices (need $DEVICES) — reusing"
    else
        echo "provision: state file has $existing devices (need $DEVICES) — re-provisioning"
        "$BENCH_BIN" provision \
            -base-url "$BASE_URL" \
            -housekeeping-key "$HK_KEY" \
            -devices "$DEVICES" \
            -state "$STATE"
    fi
else
    "$BENCH_BIN" provision \
        -base-url "$BASE_URL" \
        -housekeeping-key "$HK_KEY" \
        -devices "$DEVICES" \
        -state "$STATE"
fi

# ── ingest ──────────────────────────────────────────────────────────────────
"$BENCH_BIN" ingest \
    -state "$STATE" \
    -devices "$DEVICES" \
    -rate "$RATE" \
    -duration "$INGEST_DURATION" \
    -out "$RESULTS_DIR"

# ── connstorm ───────────────────────────────────────────────────────────────
"$BENCH_BIN" connstorm \
    -state "$STATE" \
    -devices "$STORM_DEVICES" \
    -out "$RESULTS_DIR"

# ── query ───────────────────────────────────────────────────────────────────
"$BENCH_BIN" query \
    -state "$STATE" \
    -devices "$DEVICES" \
    -out "$RESULTS_DIR"

# ── done ────────────────────────────────────────────────────────────────────
if [[ -n "${STATS_PID:-}" ]]; then
    kill "$STATS_PID" 2>/dev/null || true
    wait "$STATS_PID" 2>/dev/null || true
fi

echo ""
echo "done — results in $RESULTS_DIR"
echo "  host:        $RESULTS_DIR/host.txt"
echo "  tier env:    $RESULTS_DIR/${TIER}.env"
echo "  stats:       $STATS_CSV"
echo "  provision:   $RESULTS_DIR/provision-*.json"
echo "  ingest:      $RESULTS_DIR/ingest-*.json"
echo "  connstorm:   $RESULTS_DIR/connstorm-*.json"
echo "  query:       $RESULTS_DIR/query-*.json"
