# Triggers

Triggers allow operators to react to device events by firing HTTP webhooks. Astrate implements Astarte's simple trigger model.

## Trigger types

### Data triggers

Fire when a device publishes data matching specific criteria:

- **Interface + path pattern** -- match on a specific interface and optional path.
- **Value conditions** -- filter by value comparisons (e.g., temperature > 30).

### Device triggers

Fire on device lifecycle events:

- **`device_connected`** -- device established an MQTT connection.
- **`device_disconnected`** -- device disconnected.
- **`device_error`** -- device message was rejected (with error name and metadata).
- **`incoming_introspection`** -- device published an introspection update.

## Trigger actions

### HTTP webhook

The primary action: POST the trigger event as JSON to a configured URL.

```json
{
  "realm": "<realm>",
  "device_id": "<device_id>",
  "event": {
    "type": "data",
    "interface": "<interface_name>",
    "path": "<endpoint_path>",
    "value": <decoded_value>,
    "timestamp": "<ts>"
  }
}
```

- Retry with exponential backoff on 5xx / transport failures.
- Configurable retry cap, event TTL, and maximum in-flight deliveries per trigger.
- Delivery outcomes tracked as Prometheus metrics.

## Trigger definition

Trigger JSON stored via Realm Management:

```json
{
  "name": "my-trigger",
  "action": {
    "type": "http",
    "http_url": "https://example.com/webhook",
    "http_method": "POST",
    "http_headers": { "Authorization": "Bearer token" }
  },
  "simple_triggers": [
    {
      "type": "data_trigger",
      "interface_name": "org.example.Sensors",
      "interface_major": 1,
      "match_path": "/temperature",
      "value_match": {
        "type": "match",
        "operator": "gt",
        "value": 30.0
      }
    }
  ]
}
```

## Delivery policies

A trigger may reference a named delivery policy for finer-grained retry control:

- **`error_handlers`** -- per-status-code or per-error-class (e.g., `server_error`, `any_error`) retry/discard verdicts.
- **`retry_times`** -- maximum retry attempts.
- **`event_ttl`** -- time-to-live for queued events.
- **`maximum_capacity`** -- per-policy in-flight delivery bound.

A trigger without a policy uses Astrate's fixed default: 4xx is permanent failure, 5xx and transport failures retried with backoff.

## Execution model

- Triggers are evaluated **post-persist**, in the shard goroutine that handled the message.
- Matching is fast (compiled matchers, no per-message JSON parsing).
- Webhook delivery is asynchronous -- a slow webhook does not block ingestion.
- A slow viewer drops frames rather than backpressuring ingestion.

## Events bus

The in-process fan-out bus (`internal/engine/stream`) provides a `Subscribe(realm, filter)` interface for live consumers (WebSocket/SSE). Non-blocking sends ensure slow consumers do not affect ingestion throughput.
