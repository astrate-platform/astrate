# Configuration Reference

Astrate reads one TOML file (`astrate -config path.toml`) with `ASTRATE_*` environment overrides. Precedence: **built-in defaults < TOML file < environment variables**.

The annotated example is at [`internal/config/config.example.toml`](https://github.com/atsetilam/astrate/blob/main/internal/config/config.example.toml).

## Required values

Only two things are mandatory:

- `database.dsn` (or `ASTRATE_DATABASE_DSN`) — PostgreSQL/TimescaleDB connection string.
- Broker TLS: `mqtt.tls_cert_file` + `mqtt.tls_key_file` — required **unless** `mqtt.insecure_dev_mode` is set.

The master encryption key is supplied out of band, never in the config body.

## Precedence rule

```
built-in defaults  <  TOML file  <  ASTRATE_* environment
```

Every value below is the built-in default unless noted. Environment overrides exist for operationally critical fields, named `ASTRATE_<SECTION>_<FIELD>` (e.g. `ASTRATE_DATABASE_DSN`, `ASTRATE_HTTP_ADDR`).

---

## `[database]` — PostgreSQL/TimescaleDB

| Key | Type | Default | Env override | Required | Description |
|---|---|---|---|---|---|
| `dsn` | string | `""` | `ASTRATE_DATABASE_DSN` | **Yes** | PostgreSQL + TimescaleDB connection string. |

## `[http]` — REST API listener

| Key | Type | Default | Env override | Required | Description |
|---|---|---|---|---|---|
| `addr` | string | `":8080"` | `ASTRATE_HTTP_ADDR` | No | Bind address for the single REST listener (pairing, realm management, housekeeping, appengine, native endpoints). |
| `tls_cert_file` | string | `""` | `ASTRATE_HTTP_TLS_CERT_FILE` | No | Optional in-binary TLS cert. Leave empty to terminate TLS at a reverse proxy. |
| `tls_key_file` | string | `""` | `ASTRATE_HTTP_TLS_KEY_FILE` | No | Optional in-binary TLS key. Must be set together with `tls_cert_file`. |
| `cors_allowed_origins` | list[string] | `[]` | `ASTRATE_HTTP_CORS_ALLOWED_ORIGINS` (comma-separated) | No | Browser origins allowed for CORS. `"*"` allows any; empty disables CORS. The Astarte Dashboard SPA needs its origin here (e.g. `"http://localhost:4040"`). |

## `[mqtt]` — Embedded MQTT broker

| Key | Type | Default | Env override | Required | Description |
|---|---|---|---|---|---|
| `addr` | string | `":8883"` | `ASTRATE_MQTT_ADDR` | No | Bind address for the mTLS MQTT listener. |
| `tls_cert_file` | string | `""` | `ASTRATE_MQTT_TLS_CERT_FILE` | **Yes*** | Broker server cert. Required unless `insecure_dev_mode` is true. |
| `tls_key_file` | string | `""` | `ASTRATE_MQTT_TLS_KEY_FILE` | **Yes*** | Broker server key. Required unless `insecure_dev_mode` is true. |
| `insecure_dev_mode` | bool | `true` | `ASTRATE_MQTT_INSECURE_DEV_MODE` | No | Bind a plaintext listener on `dev_addr` that trusts the claimed client ID. **Never enable in production.** |
| `dev_addr` | string | `":1883"` | — | No | Address for the plaintext dev-mode listener. |
| `session_store_path` | string | `"sessions.db"` | `ASTRATE_MQTT_SESSION_STORE_PATH` | No | bbolt/pebble file for MQTT session persistence. On a read-only container FS, must point at the writable session volume. |
| `advertised_url` | string | `""` | `ASTRATE_MQTT_ADVERTISED_URL` | No | Broker URL handed to devices by the pairing info endpoint. Empty derives `"mqtts://<addr>"`; set when devices reach the broker by another host. |
| `max_packet_bytes` | uint32 | `278528` (272 KB) | — | No | Maximum MQTT packet size. |

*Required unless `insecure_dev_mode` is `true`.

## `[engine]` — Ingestion pipeline

| Key | Type | Default | Env override | Required | Description |
|---|---|---|---|---|---|
| `shards` | int | `16` | `ASTRATE_ENGINE_SHARDS` | No | Number of ordered ingestion shards. Each shard processes messages strictly in order. |
| `shard_queue` | int | `4096` | — | No | Bounded channel capacity per shard. Backpressure hits the broker hook when saturated. |
| `batch_max_rows` | int | `64` | — | No | Micro-batch flush threshold (row count). |
| `batch_max_wait` | duration | `"50ms"` | — | No | Micro-batch flush threshold (time). Batches flush at `batch_max_rows` or `batch_max_wait`, whichever comes first. |
| `max_payload_bytes` | int | `65536` (64 KB) | — | No | Maximum accepted payload size (both BSON and JSON). |

## `[pairing]` — Credential issuance & rate limits

| Key | Type | Default | Env override | Required | Description |
|---|---|---|---|---|---|
| `cert_ttl` | duration | `"720h"` (30 days) | — | No | Client certificate validity. |
| `enforce_latest_cert` | bool | `false` | — | No | Reject connections presenting a certificate older than the device's latest issuance (always-online CRL). Enable for fleets that rotate while devices hold older certs. |
| `register_rate` | float | `5.0` | — | No | Token-bucket rate limit for device registration (per IP, requests/sec). |
| `register_burst` | int | `10` | — | No | Token-bucket burst for device registration. |
| `credentials_rate` | float | `5.0` | — | No | Token-bucket rate limit for credential requests (per IP, requests/sec). |
| `credentials_burst` | int | `10` | — | No | Token-bucket burst for credential requests. |
| `bcrypt_cost` | int | `10` | — | No | bcrypt cost for hashing credentials secrets. |

## `[housekeeping]` — Instance-admin keys

| Key | Type | Default | Env override | Required | Description |
|---|---|---|---|---|---|
| `jwt_public_keys` | list[string] | `[]` | — | No | PEM blocks inline for instance-admin JWT public keys (claim `a_ha`). |
| `jwt_public_key_files` | list[string] | `[]` | — | No | File paths to PEM public keys. Both inline and file references are concatenated. |

## `[storage]` — Retention

| Key | Type | Default | Env override | Required | Description |
|---|---|---|---|---|---|
| `retention` | duration | `""` (disabled) | — | No | Global drop-chunks retention. Empty/0 disables it; per-endpoint TTL still applies. Example: `"8760h"` (one year). |

## `[security]` — Master key

| Key | Type | Default | Env override | Required | Description |
|---|---|---|---|---|---|
| `master_key_file` | string | `""` | `ASTRATE_SECURITY_MASTER_KEY_FILE` | No | File holding the AES-256 master key that seals realm CA private keys. When empty, falls back to `ASTRATE_MASTER_KEY` (64 hex chars) or `ASTRATE_MASTER_KEY_FILE`. |

!!! warning "Master key"
    Losing the master key means re-issuing realm CAs. Devices re-pair automatically at their next credential rotation since their credentials secret still works. Keep a separate backup of the master key.

## `[realm]` — Auto-provision realm

| Key | Type | Default | Env override | Required | Description |
|---|---|---|---|---|---|
| `name` | string | `""` | `ASTRATE_REALM_NAME` | No | Realm name. Empty disables auto-provisioning. |
| `jwt_public_key` | string | `""` | `ASTRATE_REALM_JWT_PUBLIC_KEY` | Conditional | PEM public key for the auto-provisioned realm. Required when `name` is set. |
| `jwt_public_key_file` | string | `""` | `ASTRATE_REALM_JWT_PUBLIC_KEY_FILE` | Conditional | File path to the PEM public key. Alternative to inline `jwt_public_key`. |
| `device_registration_limit` | int32 | `nil` (unlimited) | — | No | Maximum number of devices in this realm. `nil` or 0 = unlimited. |

## `[log]` — Logging

| Key | Type | Default | Env override | Required | Description |
|---|---|---|---|---|---|
| `level` | string | `"info"` | `ASTRATE_LOG_LEVEL` | No | Log level: `debug`, `info`, `warn`, `error`. |
| `format` | string | `"json"` | `ASTRATE_LOG_FORMAT` | No | Log format: `json` (structured) or `text`. |

---

## Example: minimal production config

```toml
[database]
dsn = "postgres://astrate:secret@db.example.com:5432/astrate?sslmode=require"

[mqtt]
addr = ":8883"
tls_cert_file = "/etc/astrate/broker.crt"
tls_key_file  = "/etc/astrate/broker.key"
insecure_dev_mode = false
session_store_path = "/var/lib/astrate/sessions.db"
advertised_url = "mqtts://astrate.example.com:8883"

[http]
addr = ":8080"
cors_allowed_origins = ["https://dashboard.example.com"]

[security]
master_key_file = "/etc/astrate/master.key"

[housekeeping]
jwt_public_key_files = ["/etc/astrate/housekeeping.pub"]

[log]
level = "info"
format = "json"
```

## Example: dev-mode docker compose

The included `docker-compose.yml` uses environment variables:

```sh
docker compose --profile full up -d --build
```

This sets `ASTRATE_MQTT_INSECURE_DEV_MODE=true`, `ASTRATE_MASTER_KEY` (throwaway), and auto-provisions a `test` realm. See [Deployment](deployment.md) for details.
