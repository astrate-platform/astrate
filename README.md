# Astrate

Astrate is a lean, single-binary re-implementation of the
[Astarte](https://github.com/astarte-platform/astarte) IoT platform in Go:
a modular monolith (pairing, MQTT ingestion, engine, REST APIs in one process)
backed by PostgreSQL + TimescaleDB, with an embedded MQTT broker — built to
stay wire-compatible with the official Astarte device SDKs while running
comfortably on a 1–2 GB VPS or at the edge, with no Kubernetes and no
Cassandra. Alongside Astarte's BSON payloads it natively accepts a documented
plain-JSON profile for ultra-constrained clients (e.g. AtomVM devices).

## Status

Pre-release, under active milestone-driven development. The architecture and
plan are frozen in:

- [docs/DESIGN.md](docs/DESIGN.md) — architectural design (service mapping,
  data model, wire compatibility, security/pairing)
- [docs/ROADMAP.md](docs/ROADMAP.md) — milestones M0–M9, verification tiers,
  SDK conformance checkpoints
- [docs/OPERATIONS.md](docs/OPERATIONS.md) — running Astrate: configuration,
  deployment profiles, backups, and the CA re-keying runbook

## Quick start (docker compose)

Brings up TimescaleDB and Astrate (development profile: plaintext MQTT,
self-signed broker cert, a throwaway master key — see
[docs/OPERATIONS.md](docs/OPERATIONS.md) before exposing it):

```sh
docker compose --profile full up -d --build
curl localhost:8080/astrate/v1/readiness        # {"status":"ok",...}
curl localhost:8080/astrate/v1/metrics          # Prometheus metrics
```

Configure everything from one TOML file or `ASTRATE_*` environment variables;
see [`internal/config/config.example.toml`](internal/config/config.example.toml)
for the annotated reference. Only the database DSN and (outside dev mode) the
broker TLS files are required.

```sh
astrate -config /etc/astrate/astrate.toml       # file
ASTRATE_DATABASE_DSN=postgres://… astrate        # env-only
```

On a bare VPS, run `apt install postgresql-16 timescaledb`, then the single
static `astrate` binary — no Kubernetes, no Cassandra. See
[docs/OPERATIONS.md](docs/OPERATIONS.md).

## Development

Requires Go (see `go.mod`) and Docker (for integration tests).

```sh
make tools             # install pinned dev tooling (golangci-lint)
make build             # static compile -> dist/astrate
make lint test         # T1: lint + unit tests
make test-integration  # T2: tests against a TimescaleDB container
make up                # local TimescaleDB only (fast T2 iteration)
```
