# Contributing

This guide covers prerequisites, development setup, testing, code style, and how to contribute to Astrate.

## Prerequisites

- **Go** ≥ 1.22 (toolchain managed by `go.mod`)
- **Docker** (for TimescaleDB test containers and local dev stack)
- **golangci-lint** v2.12.2 (installed via `make tools`)
- **Git**

## Dev setup

```sh
# Clone the repo
git clone https://github.com/atsetilam/astrate.git
cd astrate

# Install developer toolchain (golangci-lint)
make tools

# Start local TimescaleDB
make up

# Build the project
make build
```

### Full dev stack (with Astrate + Dashboard)

```sh
docker compose --profile full up -d --build
curl localhost:8080/astrate/v1/readiness
```

This runs:
- `astrate` with `insecure_dev_mode` (plaintext MQTT on `:1883`)
- `timescaledb` (PostgreSQL 16 + TimescaleDB)
- `dashboard` (Astarte Dashboard at `http://localhost:4040`)

See [Deployment](deployment.md) for details.

## Test tiers

Astrate uses five verification tiers (see [ROADMAP.md](ROADMAP.md) §0.2):

| Tier | Name | Command | What it covers |
|---|---|---|---|
| **T1** | Unit | `make test` | Pure Go tests, no Docker, no network. Runs every commit. |
| **T2** | Integration | `make test-integration` | `//go:build integration` tests; testcontainers boots TimescaleDB. Runs in CI and on-demand locally. |
| **T3** | Component / E2E | `make test-e2e` | Full wired binary, real broker + DB, test MQTT/HTTP clients. |
| **T4** | Conformance | `make test-conformance` | Official Astarte SDKs (Go, Python), `astartectl`, AtomVM JSON simulator against a composed Astrate. |
| **T5** | Non-functional | (nightly CI) | Load/footprint budgets, security probes. |

### Running tests locally

```sh
# T1: unit tests (fast, no dependencies)
make test

# T2: integration tests (requires Docker)
make test-integration

# T3: E2E (requires full stack or compose)
make test-e2e

# T4: conformance (lands in M9)
make test-conformance
```

### Reusing an external database for T2

If you have a TimescaleDB instance already running:

```sh
ASTRATE_TEST_DSN="postgres://user:pass@host:5432/dbname?sslmode=disable" make test-integration
```

## Code style

- **Linting:** `make lint` runs `golangci-lint` with the pinned config (`.golangci.yml`). Fix all warnings before submitting.
- **Formatting:** standard `gofmt`. Run `gofmt -s -w .` before committing.
- **Dependencies:** all third-party deps are pinned in `go.mod`. Never add a dep without discussing it first — lean-first principle means the bar for new dependencies is high.
- **No `// TODO`** except at the three named extension points in the ROADMAP (external-bus intake, trigger forwarding, Timescale toolkit probing).
- **Interfaces before implementations:** where package A consumes B through an interface, the interface lives in A's code (hexagonal-lite).
- **Tests ride with their subject:** every `foo.go` should have a `foo_test.go` in the same package.

## Project structure

```
astrate/
├── cmd/astrate/              # main: config, wiring, lifecycle
├── internal/
│   ├── broker/               # embedded MQTT broker, auth/ACL hooks
│   ├── engine/               # ingestion pipeline: shards, validation, persistence
│   ├── pairing/              # registration, credentials, CSR signing
│   │   └── ca/               # embedded per-realm CA
│   ├── appengine/            # REST API + WebSocket/SSE live stream
│   ├── realm/                # interface/trigger CRUD
│   ├── housekeeping/         # realm lifecycle
│   ├── auth/                 # JWT validation, authz claims
│   ├── store/                # pgx repositories, migrations
│   └── config/               # TOML/env config
├── pkg/
│   ├── interfaceschema/      # Interface JSON parsing, trie compiler
│   ├── payload/              # BSON/JSON dual codec
│   ├── deviceid/             # 128-bit device ID handling
│   └── astarteapi/           # shared API envelope types
├── migrations/               # SQL migrations (embedded)
├── test/                     # conformance and E2E test suites
└── docker-compose.yml
```

**Dependency rule:** `pkg/*` never imports `internal/*`. `internal/store` is imported by domain packages but never imports them back.

## How to contribute

1. **Fork and create a branch** from `main`.
2. **Write your change** following the code style above.
3. **Run `make test` and `make lint`** — both must pass.
4. **Write tests** for new behavior (T1 minimum; integration tests for store/broker changes).
5. **Open a PR** against `main`. Describe what changed and why.
6. **CI must pass** — T1 + T2 at minimum; T3/T4 if your change touches protocol-relevant code.

### What to work on

Check the [ROADMAP.md](ROADMAP.md) for the current milestone. High-value contributions:

- Bug fixes (check open issues)
- Missing tests for existing code
- Documentation improvements
- Conformance test coverage

### Commit conventions

- Keep commits focused: one logical change per commit.
- Write clear commit messages (no need for conventional-commits format, but be descriptive).
- Never commit secrets, keys, or credentials.
