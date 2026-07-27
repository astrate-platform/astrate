# Migration from Astarte

!!! warning "Preliminary"
    This guide is based on the design documents and architecture analysis. No real-world Astarte-to-Astrate migration has been performed yet. Treat this as a planning reference, not a battle-tested runbook.

## Overview

Astrate is wire-compatible with Astarte: unmodified official device SDKs work against it. This means migration is possible without changing firmware, but requires careful planning around data export, schema differences, and the cutover window.

## Architecture comparison

| Concern | Astarte | Astrate |
|---|---|---|
| Runtime | ~8 Elixir/OTP microservices | One Go binary |
| Orchestration | Kubernetes + astarte-operator | `docker-compose` or bare binary |
| Data store | Cassandra/ScyllaDB | PostgreSQL 16 + TimescaleDB |
| Message bus | RabbitMQ | In-process Go channels |
| MQTT broker | VerneMQ + plugin | Embedded mochi-mqtt |
| Certificate authority | CFSSL sidecar | Embedded (Go crypto/x509) |
| Payloads | BSON only | BSON + JSON |

## Service mapping

| Astarte component | Astrate equivalent |
|---|---|
| Pairing API + CFSSL | `internal/pairing` (embedded CA) |
| VerneMQ + plugin | `internal/broker` (embedded) |
| DUP (Data Updater Plant) | `internal/engine` (sharded pipeline) |
| AppEngine API | `internal/appengine` |
| Realm Management API | `internal/realm` |
| Housekeeping API | `internal/housekeeping` |
| Trigger Engine | `internal/engine/triggers` |
| Channels (Phoenix socket) | WebSocket/SSE endpoint (additive) |

## Pre-migration checklist

- [ ] **Audit your Astarte deployment:** note version, number of realms, device count, installed interfaces, active triggers.
- [ ] **Export realm data:** use `astartectl` to export interface definitions, trigger definitions, and device lists per realm.
- [ ] **Export device data:** use the AppEngine API or direct Cassandra queries to export device properties and recent datastreams. Astrate does not import Cassandra data directly — you'll need to write it to the new PostgreSQL instance.
- [ ] **Document your JWT keys:** note the realm-level JWT public keys and the instance-admin (`a_ha`) key for Housekeeping.
- [ ] **Plan the DNS/window:** devices will need to point at the new Astrate instance (Pairing API URL, MQTT broker URL).
- [ ] **Test in staging:** deploy Astrate alongside your Astarte instance, register test devices, verify the full flow.

## Data export from Astarte

### Interfaces and triggers

```sh
# Export all interfaces for a realm
astartectl realm-management interfaces list -r <realm> -t <jwt> -a <api-url>

# Export all triggers
astartectl realm-management triggers list -r <realm> -t <jwt> -a <api-url>
```

### Device data

Use the AppEngine API:

```sh
# List devices
curl -H "Authorization: Bearer <jwt>" \
  <api-url>/appengine/v1/<realm>/devices?details=true

# Get device properties
curl -H "Authorization: Bearer <jwt>" \
  <api-url>/appengine/v1/<realm>/devices/<id>/interfaces/<interface>

# Get datastream data
curl -H "Authorization: Bearer <jwt>" \
  "<api-url>/appengine/v1/<realm>/devices/<id>/interfaces/<interface>/<path>?since=..."
```

For large-scale export, query Cassandra directly and transform to the Astrate PostgreSQL schema (see [DESIGN.md §2](DESIGN.md) for the target schema).

### Realm configuration

```sh
# List realms
astartectl housekeeping realms list -t <jwt> -a <api-url>

# Get realm details (JWT keys, device limit)
astartectl housekeeping realms show <realm> -t <jwt> -a <api-url>
```

## Schema mapping

Astrate uses PostgreSQL instead of Cassandra, with a shared-table tenancy model (`realm_id` column) instead of per-keyspace isolation.

Key differences:
- **No Cassandra keyspace per realm.** Realms are rows in a `realms` table.
- **No RabbitMQ.** Inter-service communication is in-process Go channels.
- **TimescaleDB hypertables** replace Cassandra tables for datastreams. Compression and retention are handled by TimescaleDB policies.
- **Interface definitions** are stored as JSONB with generated columns for routing-critical fields.

The migration tool will need to:
1. Create realms via the Housekeeping API (Astrate auto-generates CAs).
2. Install interfaces via the Realm Management API.
3. Import device data into PostgreSQL (properties via AppEngine API, datastreams via direct inserts or the API).

## Device behavior during cutover

### What devices need

Devices need:
1. The Pairing API URL (to register and get credentials).
2. The MQTT broker URL (to connect).
3. A valid credentials secret (obtained at registration).
4. A client certificate issued by the realm CA.

### Cutover strategy

1. **Stop writes to Astarte** (optional but recommended for data consistency).
2. **Export final device state** from Astarte.
3. **Deploy Astrate** with the same realm names.
4. **Create realms** via Housekeeping API (mints new CAs).
5. **Install interfaces** via Realm Management API.
6. **Re-register devices** — devices need new credentials from the new realm CA.
   - If devices store their credentials secret, they can re-CSR on next boot.
   - If not, re-register via the agent API and deliver new credentials out-of-band.
7. **Update device configuration** to point at the new Pairing API and broker URLs.
8. **Verify connectivity** — devices should connect, send introspection, and resume data flow.

### During the cutover window

- Devices trying to connect to the old Astarte instance will fail (expected).
- Devices pointing at the new Astrate instance will need to go through the full pairing flow (register → credentials → connect).
- There is no "transparent handoff" — devices must re-pair because the CA is different.

## Rollback plan

If migration fails:

1. **Revert DNS** to point devices back at the original Astarte instance.
2. **Devices will re-pair** with Astarte on their next connection attempt (they lost their Astrate-issued certs, but their original credentials secret still works with Astarte if you re-imported device registrations).
3. **Data written to Astrate** during the migration window is not automatically synced back to Astarte. You'll need to export and re-import if you want to preserve it.
4. **Keep the Astarte deployment running** in parallel during the migration window so rollback is immediate.

## Known limitations

- **No automated migration tool yet.** The steps above are manual. A migration CLI tool is not on the v1 roadmap but could be built by the community.
- **Data not migrated automatically.** Historical data in Cassandra must be exported and re-imported into PostgreSQL separately.
- **Device re-pairing required.** Because the CA changes, all devices must re-pair. There is no CA import mechanism from Astarte to Astrate.
- **Dashboard:** The upstream Astarte Dashboard works against Astrate, but Device Live Events requires the Channels socket (see [Compatibility](compatibility.md) for details).

## See also

- [Configuration Reference](configuration-reference.md) — for setting up the new Astrate instance
- [Deployment](deployment.md) — Docker Compose and bare VPS deployment
- [Deployment](deployment.md) — Docker Compose and bare VPS deployment details
- [Compatibility](compatibility.md) — wire-compatible surfaces and deliberate deviations
