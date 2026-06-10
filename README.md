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

## Development

Requires Go (see `go.mod`) and Docker (for integration tests).

```sh
make tools             # install pinned dev tooling (golangci-lint)
make build             # static compile
make lint test         # T1: lint + unit tests
make test-integration  # T2: tests against a TimescaleDB container
make up                # local TimescaleDB via docker compose
```

A full quickstart (docker compose with the `astrate` service, bare-VPS notes)
arrives with milestone M8.
