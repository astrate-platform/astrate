# Deployment

Astrate is one static binary plus PostgreSQL/TimescaleDB. Two deployment profiles: Docker Compose (development) and bare VPS (production).

## Docker Compose (development)

```sh
docker compose --profile full up -d --build
```

Runs:
- `astrate` in `insecure_dev_mode` (plaintext MQTT on `:1883`).
- `timescaledb` (tuned: `shared_buffers=256MB`, `max_connections=50`).
- `dashboard` (upstream Astarte Dashboard at `http://localhost:4040`).

The `full` profile auto-provisions a `test` realm with the committed dev-only keypair. Generate a login token:

```sh
astartectl utils gen-jwt all-realm-apis -k deploy/devrealm/realm_private.pem
```

Volumes: `pgdata` (PostgreSQL), session store (`/var/lib/astrate`, bbolt).

**Never expose this profile to untrusted networks.**

## Bare VPS (production)

```sh
apt install postgresql-16 timescaledb
astrate -config /etc/astrate/astrate.toml
```

Astrate self-migrates the schema on boot.

### Target footprint

| Component | Steady-state RSS |
|---|---|
| Astrate | <= 150 MB at ~1k devices |
| PostgreSQL | <= 768 MB (tuned) |

## Configuration

One TOML file + `ASTRATE_*` env overrides. Precedence: **defaults < TOML < environment**.

### Required settings

- `database.dsn` -- PostgreSQL/TimescaleDB DSN (or `ASTRATE_DATABASE_DSN`).
- `mqtt.tls_cert_file` + `mqtt.tls_key_file` -- broker TLS identity (unless `insecure_dev_mode` is set).

### Master key

The master encryption key that seals realm CA private keys is supplied out of band:

- `ASTRATE_MASTER_KEY` (64 hex chars or base64 of 32 bytes)
- `ASTRATE_MASTER_KEY_FILE`
- `security.master_key_file`

Losing it means re-issuing realm CAs (devices re-pair automatically at next credential rotation).

### Environment overrides

Named `ASTRATE_<SECTION>_<FIELD>` -- e.g. `ASTRATE_HTTP_ADDR`, `ASTRATE_MQTT_ADDR`, `ASTRATE_ENGINE_SHARDS`, `ASTRATE_LOG_LEVEL`.

## Endpoints

| Port | Protocol | Purpose |
|---|---|---|
| `:8080` | HTTP/HTTPS | REST API: `/pairing/v1`, `/realmmanagement/v1`, `/housekeeping/v1`, `/appengine/v1`, `/astrate/v1/{health,readiness,metrics}` |
| `:8883` | mTLS MQTT | Astarte MQTT v1 device connections |
| `:1883` | plaintext MQTT | Development only (`insecure_dev_mode`) |

## TLS options

- **In-binary TLS:** set `http.tls_cert_file` / `http.tls_key_file`.
- **Reverse proxy:** terminate TLS at nginx/HAProxy, proxy to `:8080`. Documented compose profiles for both.
- **HSTS:** enabled on the API by default.

## Auto-provisioning a realm

Set the `[realm]` block in config:

```toml
[realm]
name = "home"
jwt_public_key_file = "/etc/astrate/realm_public.pem"
```

Astrate creates the realm on first boot (no-op if it exists). The broker reloads its CA trust automatically.

## Container details

- Distroless static base, non-root user, read-only root filesystem.
- Only writable path: session-store volume at `/var/lib/astrate`.
- Image size budget: <= 30 MB.
- `CGO_ENABLED=0`, `-trimpath`, `-ldflags=-s -w`.

## Backups

Two pieces of durable state:

- **PostgreSQL** (`pgdata`): all metadata, properties, datastreams. Back up with `pg_dump`/`pg_basebackup` or volume snapshots. This is the source of truth.
- **Session store** (`/var/lib/astrate`, bbolt): MQTT session/offline-queue state. Losing it only forces clean reconnects.

Keep the master key backed up **separately** from the database.
