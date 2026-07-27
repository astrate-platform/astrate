# Observability

Astrate exposes Prometheus metrics, structured JSON logging, and health/readiness probes.

## Health endpoints

| Endpoint | Purpose |
|---|---|
| `GET /astrate/v1/health` | Liveness probe. Returns 200 if the process is running. |
| `GET /astrate/v1/readiness` | Readiness probe. Pings the database and checks the broker listener. Returns 503 with per-check status if any dependency is down. Use for load-balancer and orchestrator probes. |
| `GET /astrate/v1/metrics` | Prometheus scrape endpoint. |

All three endpoints are served on the same HTTP listener as the REST API (`:8080` by default).

## Prometheus metrics

Exposed at `/astrate/v1/metrics`. All counters and histograms use the `astrate_` namespace.

### Engine (ingestion pipeline)

| Metric | Type | Labels | Description |
|---|---|---|---|
| `astrate_engine_messages_total` | counter | — | Inbound device messages submitted to the engine. |
| `astrate_engine_rejects_total` | counter | `reason` | Messages rejected by the validation pipeline (e.g. `introspection_mismatch`, `type_mismatch`, `ownership_violation`). |
| `astrate_engine_persist_ops_total` | counter | `kind` | Validated operations committed to storage, by kind. |
| `astrate_engine_shard_depth` | gauge | `shard` | Messages queued per pipeline shard (sampled at scrape time). |
| `astrate_engine_batch_flush_seconds` | histogram | — | Micro-batch flush latency, including retries. |
| `astrate_engine_batch_flush_retries_total` | counter | — | Failed flush attempts that were retried (DB-outage parking). |
| `astrate_engine_qos0_dropped_total` | counter | — | QoS 0 messages dropped because their shard was full. |
| `astrate_engine_dropped_shutdown_total` | counter | — | Messages refused because the engine was draining during shutdown. |
| `astrate_engine_internal_errors_total` | counter | — | Recovered shard panics and other engine-side faults. |

### Broker

| Metric | Type | Description |
|---|---|---|
| `astrate_broker_sessions` | gauge | Live authenticated MQTT device sessions. |

### Database pool

| Metric | Type | Description |
|---|---|---|
| `astrate_db_pool_acquired_conns` | gauge | Connections currently in use. |
| `astrate_db_pool_idle_conns` | gauge | Idle connections in the pool. |
| `astrate_db_pool_total_conns` | gauge | Total connections in the pool. |
| `astrate_db_pool_max_conns` | gauge | Configured maximum pool size. |

### Triggers

| Metric | Type | Labels | Description |
|---|---|---|---|
| `astrate_engine_trigger_deliveries_total` | counter | `outcome` | Trigger action deliveries by outcome: `delivered`, `failed`, `dropped`, `expired`, `forwarded`, `skipped`. |
| `astrate_engine_trigger_retries_total` | counter | — | Failed webhook attempts that were retried. |

### Live stream

| Metric | Type | Description |
|---|---|---|
| `astrate_engine_stream_dropped_total` | counter | Live events dropped because a subscriber's channel was full. |

### Go runtime (built-in)

The Prometheus registry also includes the standard `go_*` and `process_*` collectors for RSS, GC, and goroutine counts.

## Structured logging

- JSON format via `log/slog`.
- Per-domain log levels configurable.
- Every rejection includes: device ID, interface, path, and reason.
- Trigger deliveries log the applied rule (`reason`, `policy`).

## Configuration

Logging is configured via the `[log]` section (see [Configuration Reference](configuration-reference.md)):

```toml
[log]
level = "info"    # debug, info, warn, error
format = "json"   # json, text
```

Environment: `ASTRATE_LOG_LEVEL=debug`.

## See also

- [Deployment](deployment.md) — health probe usage in Docker Compose and orchestrators
- [Configuration Reference](configuration-reference.md) — all TOML keys including `[log]`
- [Troubleshooting](troubleshooting.md) — debugging with health/readiness/metrics endpoints
