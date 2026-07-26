#!/usr/bin/env bash
# Brings up upstream Astarte for benchmarking using the project's own
# standalone docker-compose, pinned to the version Astrate targets. Clones
# into bench/.astarte (gitignored), generates the compose keys/certificates
# with Astarte's script, and waits for the housekeeping API.
#
# The gateway answers on http://api.astarte.localhost and the broker on
# mqtts://broker.astarte.localhost:8883 — *.localhost resolves to loopback on
# macOS and systemd-resolved Linux. A load generator on another machine cannot
# be pointed at the Docker host with /etc/hosts: macOS collapses every
# *.localhost name to loopback before consulting that file, so the entry is
# silently inert. Tunnel to loopback instead, which also carries the vhost the
# gateway routes on:
#
#   ssh -N -L 8080:127.0.0.1:80 -L 8883:127.0.0.1:8883 user@dockerhost
#   ... -base-url http://api.astarte.localhost:8080
set -euo pipefail

ASTARTE_VERSION="${ASTARTE_VERSION:-v1.2.0}"
scripts="$(cd "$(dirname "$0")" && pwd)"
cd "$scripts/.."

if [[ ! -d .astarte ]]; then
    git clone --depth 1 --branch "$ASTARTE_VERSION" \
        https://github.com/astarte-platform/astarte .astarte
fi

cd .astarte
# Generates the housekeeping keypair + device broker certificates into
# compose/. Astarte used to ship this as generate-compose-files.sh in the
# checkout; as of 1.2 the repo carries no such script and its README calls an
# initializer container instead, which is what this runs. The image tag does
# not track the Astarte version — 1.2.0's README pins 1.1 — hence the separate
# variable.
#
# Skipped when the key already exists, so provisioned realms stay valid across
# re-runs.
ASTARTE_INITIALIZER="${ASTARTE_INITIALIZER:-astarte/docker-compose-initializer:1.1}"
if [[ ! -f compose/astarte-keys/housekeeping_private.pem ]]; then
    docker run --rm -v "$(pwd)/compose:/compose" "$ASTARTE_INITIALIZER"
    if [[ ! -f compose/astarte-keys/housekeeping_private.pem ]]; then
        echo "the initializer ran but produced no housekeeping_private.pem —" >&2
        echo "check $ASTARTE_INITIALIZER against the $ASTARTE_VERSION README." >&2
        exit 1
    fi
fi

# Stock compose plus bench/scripts/astarte-compose-override.yml — environment
# fixes only (see that file's header), never behaviour changes.
compose=(docker compose -f docker-compose.yml -f "$scripts/astarte-compose-override.yml")

"${compose[@]}" pull
"${compose[@]}" up -d

# Readiness means "the gateway routes to housekeeping and housekeeping answers",
# not "the request succeeded". /housekeeping/v1/version is authenticated, so a
# ready stack answers 401 without a token; that is the expected success case.
# /housekeeping/health is unauthenticated and answers 200. Measured on 1.2.0:
# health 200, v1/version 401, and anything traefik has no route for — including
# a stack that is still booting — 404. So 200-or-401 is the up/booting boundary.
echo -n "waiting for the housekeeping API"
for _ in $(seq 1 120); do
    for probe in \
        "http://api.astarte.localhost/housekeeping/health" \
        "http://api.astarte.localhost/housekeeping/v1/version"; do
        code="$(curl -s -o /dev/null -w '%{http_code}' --max-time 5 "$probe" || true)"
        if [[ "$code" == "200" || "$code" == "401" ]]; then
            echo " — up ($probe → $code)"
            key="$(pwd)/compose/astarte-keys/housekeeping_private.pem"
            echo "API:    http://api.astarte.localhost"
            echo "Broker: (advertised by pairing, typically mqtts://broker.astarte.localhost:8883)"
            echo "Next:   go run . provision -base-url http://api.astarte.localhost -housekeeping-key $key"
            exit 0
        fi
    done
    echo -n .
    sleep 5
done
echo " — FAILED (docker compose logs; Cassandra can take minutes on first boot)"
exit 1
