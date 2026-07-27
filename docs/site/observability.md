# Observability

Astrate exposes Prometheus metrics, structured JSON logging, and health/readiness probes.

## Health endpoints

| Endpoint | Purpose |
|---|---|
| `GET /astrate/v1/health` | Liveness probe. Returns 200 if the process is running. |
| `GET /astrate/v1/readiness` | Readiness probe. Pings the database and checks the broker listener. Use for load-balancer and orchestrator probes. |
| `GET /astrate/v1/metrics` | Prometheus scrape endpoint. |

## Prometheus metrics

Exposed at `/astrate/v1/metrics`:

### Ingestion

| Metric | Description |
|---|---|
| `astrate_ingest_total` | Total messages received from devices |
| `astrate_ingest_reject_total` | Rejected messages, labeled by reason (`introspection_mismatch`, `unexpected_path`, `type_mismatch`, `ownership_violation`, etc.) |

### Engine

| Metric | Description |
|---|---|
| `astrate_shard_depth` | Current depth of each shard's input channel |
| `astrate_batch_flush_duration_seconds` | Histogram of batch flush latencies |
| `astrate_batch_rows` | Rows per batch flush |

### Broker

| Metric | Description |
|---|---|
| `astrate_broker_sessions` | Active MQTT sessions |
| `astrate_broker_connect_total` | Total MQTT connections |
| `astrate_broker_disconnect_total` | Total MQTT disconnections |

### Database

| Metric | Description |
|---|---|
| `astrate_db_pool_acquired` | Connection pool acquisitions |
| `astrate_db_pool_idle` | Idle connections in pool |

### Triggers

| Metric | Description |
|---|---|
| `astrate_trigger_delivery_total` | Trigger webhook deliveries, labeled by outcome (`success`, `retry`, `dropped`) |

## Structured logging

- JSON format via `log/slog`.
- Per-domain log levels configurable.
- Every rejection includes: device ID, interface, path, and reason.
- Trigger deliveries log the applied rule (`reason`, `policy`).

## Configuration

```toml
[log]
level = "info"    # debug, info, warn, error

[metrics]
addr = ":9090"    # separate scrape endpoint (optional; defaults to HTTP addr)
```

Environment: `ASTRATE_LOG_LEVEL=debug`.
