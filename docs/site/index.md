# Astrate

A lean, single-binary, Astarte-wire-compatible IoT platform in Go.

## What is Astrate?

Astrate is a spiritual fork of the [Astarte IoT Platform](https://github.com/astarte-platform/astarte) that preserves Astarte's external contracts while replacing its internals entirely. Unmodified official Astarte device SDKs (C/ESP32, Python, Go, Rust, Java/Android, Elixir) work against Astrate without a single line of SDK code changed.

## Key differences from upstream Astarte

| Concern | Astarte (upstream) | Astrate |
|---|---|---|
| Runtime | ~8 Elixir/OTP microservices | One statically linked Go binary |
| Orchestration | Kubernetes + astarte-operator | `docker-compose.yml` or bare binary + systemd |
| Time-series store | Cassandra / ScyllaDB | PostgreSQL 16 + TimescaleDB |
| Message bus | RabbitMQ | In-process Go channels |
| MQTT broker | VerneMQ | Embedded `mochi-mqtt` |
| Certificate authority | CFSSL sidecar | Embedded CA via `crypto/x509` |
| Payloads | BSON only | BSON **and** plain JSON |
| Target footprint | 4-16 GB RAM cluster | <= 1-2 GB RAM single VPS / edge node |

## Documentation pages

- [Architecture](architecture.md) -- process model, package layout, concurrency
- [Data Modeling](data-modeling.md) -- PostgreSQL + TimescaleDB schema
- [MQTT Protocol](mqtt-protocol.md) -- Astarte MQTT v1 wire compatibility
- [Pairing and Security](pairing-and-security.md) -- certificates, JWTs, credentials
- [Payload Formats](payload-formats.md) -- BSON and JSON dual codec
- [AppEngine API](appengine-api.md) -- device data, groups, live stream
- [Realm Management API](realm-management-api.md) -- interfaces and triggers
- [Housekeeping API](housekeeping-api.md) -- realm lifecycle
- [Interface Schema](interface-schema.md) -- Astarte interface JSON specification
- [Triggers](triggers.md) -- event matching and webhook delivery
- [Deployment](deployment.md) -- Docker Compose, bare VPS, configuration
- [Observability](observability.md) -- metrics, health, logging
- [Compatibility](compatibility.md) -- deviations from upstream Astarte
- [Quickstart](quickstart.md) -- get running in 5 minutes

## Design references

- [Design document](DESIGN.md) -- full architectural design (Phase 1)
- [Implementation roadmap](ROADMAP.md) -- milestone plan (Phase 2)
- [OpenAPI specifications](api/) -- REST API specs
- [API Explorer](swagger.md) -- interactive API docs
