#!/usr/bin/env bash
# Brings up upstream Astarte for benchmarking using the project's own
# standalone docker-compose, pinned to the version Astrate targets. Clones
# into bench/.astarte (gitignored), generates the compose keys/certificates
# with Astarte's script, and waits for the housekeeping API.
#
# The gateway answers on http://api.astarte.localhost and the broker on
# mqtts://broker.astarte.localhost:8883 — *.localhost resolves to loopback on
# macOS and systemd-resolved Linux. If the load generator runs on another
# machine, point those names at the Docker host in /etc/hosts.
set -euo pipefail

ASTARTE_VERSION="${ASTARTE_VERSION:-v1.2.0}"
cd "$(dirname "$0")/.."

if [[ ! -d .astarte ]]; then
    git clone --depth 1 --branch "$ASTARTE_VERSION" \
        https://github.com/astarte-platform/astarte .astarte
fi

cd .astarte
if [[ ! -x generate-compose-files.sh ]]; then
    echo "generate-compose-files.sh not found in the $ASTARTE_VERSION checkout —" >&2
    echo "check the Astarte release layout and adjust this script." >&2
    exit 1
fi
# Generates the housekeeping keypair + device broker certificates (idempotent
# enough: skip when the key already exists so provisioned realms stay valid).
if [[ ! -f compose/astarte-keys/housekeeping_private.pem ]]; then
    ./generate-compose-files.sh
fi

docker compose up -d

echo -n "waiting for the housekeeping API"
for _ in $(seq 1 120); do
    for probe in \
        "http://api.astarte.localhost/housekeeping/health" \
        "http://api.astarte.localhost/housekeeping/v1/health" \
        "http://api.astarte.localhost/housekeeping/v1/version"; do
        if curl -fsS "$probe" >/dev/null 2>&1; then
            echo " — up"
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
