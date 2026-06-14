# Operating Astrate

Astrate is one static binary plus PostgreSQL/TimescaleDB. This guide covers
configuration, the deployment profiles, backups, and key management. The
normative architecture is in [DESIGN.md](DESIGN.md); the milestone plan is in
[ROADMAP.md](ROADMAP.md).

## Configuration

Astrate reads one TOML file (`astrate -config path.toml`) with `ASTRATE_*`
environment overrides; precedence is **defaults < TOML < environment**. The
annotated reference is
[`internal/config/config.example.toml`](../internal/config/config.example.toml).

Only two things are mandatory:

- `database.dsn` (`ASTRATE_DATABASE_DSN`) — the PostgreSQL/TimescaleDB DSN.
- The broker TLS identity `mqtt.tls_cert_file` + `mqtt.tls_key_file` — required
  **unless** `mqtt.insecure_dev_mode` is set.

The master encryption key that seals realm CA private keys is supplied out of
band, never in the config body: set `ASTRATE_MASTER_KEY` (64 hex chars or
base64 of 32 bytes), `ASTRATE_MASTER_KEY_FILE`, or `security.master_key_file`.
Losing it means re-issuing realm CAs (devices re-pair automatically at their
next credential rotation, since their credentials secret still works).

Environment overrides exist for the operationally critical fields, named
`ASTRATE_<SECTION>_<FIELD>` — e.g. `ASTRATE_HTTP_ADDR`, `ASTRATE_MQTT_ADDR`,
`ASTRATE_MQTT_INSECURE_DEV_MODE`, `ASTRATE_ENGINE_SHARDS`, `ASTRATE_LOG_LEVEL`,
`ASTRATE_REALM_NAME`.

### Endpoints

- `:8080` — REST API: `/pairing/v1`, `/realmmanagement/v1`,
  `/housekeeping/v1`, `/appengine/v1`, and the Astrate-native
  `/astrate/v1/{health,readiness,metrics,…}`.
- `:8883` — mTLS MQTT (Astarte MQTT v1).
- `:1883` — plaintext MQTT, bound only under `insecure_dev_mode`.

`/astrate/v1/health` is liveness; `/astrate/v1/readiness` pings the database
and the broker listener (use it for load-balancer and orchestrator probes);
`/astrate/v1/metrics` is the Prometheus scrape endpoint.

## Deployment profiles

### Development (docker compose)

```sh
docker compose --profile full up -d --build
```

Runs the broker in `insecure_dev_mode`: devices connect over plaintext `:1883`
and Astrate mints a throwaway self-signed cert for the (unused) mTLS listener.
The compose file uses a fixed throwaway master key and DB password — never
expose this profile.

### Production

Provide real secrets and TLS, and drop dev mode:

- A strong, secret `ASTRATE_MASTER_KEY[_FILE]`.
- Broker TLS: `mqtt.tls_cert_file` / `mqtt.tls_key_file`. Devices trust the
  per-realm CA returned by the pairing info endpoint, so the broker server cert
  must be issued from a CA your fleet trusts. Set `mqtt.advertised_url` to the
  hostname devices dial (`mqtts://host:8883`).
- HTTP TLS: either terminate in-binary (`http.tls_cert_file` /
  `http.tls_key_file`) or run behind a TLS-terminating reverse proxy.
- Housekeeping admin keys: `housekeeping.jwt_public_keys` (or
  `jwt_public_key_files`) — the instance-admin JWT (`a_ha`) public keys.

The container runs as a non-root user on a read-only root filesystem; the only
writable path is the session-store volume at `/var/lib/astrate`
(`mqtt.session_store_path`).

### Bare VPS

```sh
apt install postgresql-16 timescaledb        # enable the timescaledb extension
astrate -config /etc/astrate/astrate.toml    # the single static binary
```

Astrate self-migrates the schema on boot. Target steady-state footprint: the
Astrate process ≤ 150 MB RSS at ~1k devices, PostgreSQL tuned to ≤ 768 MB
(`shared_buffers=256MB`, `max_connections=50`).

## Auto-provisioning a realm

Set the `[realm]` block (`name` + `jwt_public_key`/`_file`) and Astrate creates
the realm — minting and sealing its CA — on first boot, a no-op if it already
exists. Otherwise create realms at runtime via the Housekeeping API; the broker
reloads its CA trust automatically so new realms accept device connections
without a restart.

## Backups

Two pieces of durable state:

- **PostgreSQL** (`pgdata`): all metadata, properties, and datastreams. Back up
  with `pg_dump`/`pg_basebackup` or volume snapshots. This is the source of
  truth — realms, devices, and the (sealed) CA keys live here.
- **Session store** (`/var/lib/astrate`, bbolt): MQTT session/offline-queue
  state. Losing it only forces clean reconnects; it is not a source of truth.

Keep the master key backed up **separately** from the database: together they
decrypt the realm CA private keys.

## CA re-keying runbook

A realm's CA private key is AES-256-GCM-sealed in the database under the master
key. To recover from a suspected master-key or CA compromise:

1. Roll the master key: re-seal is not automatic — re-provision affected realm
   CAs. Deleting and recreating a realm via Housekeeping mints a fresh CA.
2. Devices keep their credentials secret, so they re-pair automatically: on the
   next `verify`/connect the old cert fails to chain, the SDK re-CSRs, and
   pairing issues a cert from the new CA.
3. To force rotation fleet-wide, shorten `pairing.cert_ttl` so existing certs
   expire quickly, and keep `pairing.enforce_latest_cert` on to reject
   superseded serials immediately.

Rotate realm REST/JWT signing keys independently via
`PUT /realmmanagement/v1/<realm>/config/auth` (multiple keys are accepted for
zero-downtime rotation); this is unrelated to the device CA.
