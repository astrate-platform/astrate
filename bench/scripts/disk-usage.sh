#!/usr/bin/env bash
# Reports the on-disk size of each platform's data directories, for the
# bytes-per-datapoint comparison. Run after ingest AND after forcing
# compaction/compression (see README "Storage efficiency") so both platforms
# are measured at rest, not mid-write.
#
#   ./scripts/disk-usage.sh astrate   # timescaledb container
#   ./scripts/disk-usage.sh astarte   # cassandra/scylla container
set -euo pipefail

target="${1:?usage: disk-usage.sh astrate|astarte}"

du_in() { # container path
    docker exec "$1" du -sb "$2" 2>/dev/null || docker exec "$1" du -sk "$2" | awk '{print $1*1024 "\t" $2}'
}

find_container() { # grep pattern
    docker ps --format '{{.Names}}' | grep -m1 "$1" || {
        echo "no running container matching '$1'" >&2
        exit 1
    }
}

case "$target" in
astrate)
    c="$(find_container timescaledb)"
    echo "container: $c"
    du_in "$c" /var/lib/postgresql/data
    ;;
astarte)
    # Upstream compose names the DB service cassandra (or scylla in some
    # revisions); take whichever is running.
    c="$(find_container 'cassandra\|scylla')"
    echo "container: $c"
    du_in "$c" /var/lib/cassandra 2>/dev/null || du_in "$c" /var/lib/scylla
    ;;
*)
    echo "usage: disk-usage.sh astrate|astarte" >&2
    exit 2
    ;;
esac
