# MQTT Protocol

Astrate implements the Astarte MQTT v1 protocol exactly. This page covers the wire format that device SDKs speak.

## Broker

- Embedded `mochi-mqtt/server` v2 in-process.
- MQTT 3.1.1 (what all Astarte SDKs speak; 5.0 also accepted as a superset).
- QoS 0/1/2, retained messages, persistent sessions.

## Connection contract

- **Listener:** TLS on `:8883` with `RequireAndVerifyClientCert`. Client CAs = per-realm CA pool.
- **Identity:** client certificate `Subject CN = <realm>/<device_id>`.
- **Client ID:** free-form on the wire (the official Python SDK sends a random paho-generated ID). Rewritten to the CN before session binding, mirroring VerneMQ's subscriber-id remap.
- **Sessions:** `clean_session=false`. Session state persists in bbolt across restarts so `session_present` survives Astrate restarts.
- **Optional plaintext:** `:1883` behind `insecure_dev_mode` flag only.

## ACL model

For identity `<realm>/<device_id>` with base topic `B = <realm>/<device_id>`:

| Action | Allowed topics |
|---|---|
| PUBLISH | `B` (introspection), `B/control/emptyCache`, `B/control/producer/properties`, `B/<interface_name><path>` for device-owned interfaces in introspection |
| SUBSCRIBE | Any filter within `B/...` -- not gated on introspection (SDKs subscribe before sending introspection) |

Everything else is denied and logged. Server-side publishes use mochi's inline client and bypass ACLs.

## Device to Astrate messages

| Topic | Payload | QoS | Handling |
|---|---|---|---|
| `<realm>/<device_id>` | Introspection: `;`-separated `name:major:minor` triples | 2 | Parse, diff vs stored, update introspection, fire triggers |
| `<realm>/<device_id>/control/emptyCache` | `1` | 2 | Re-send server-owned properties + consumer-properties purge |
| `<realm>/<device_id>/control/producer/properties` | 4-byte BE size + zlib deflated path list | 2 | Decompress, parse, purge device-owned properties not in list |
| `<realm>/<device_id>/<interface><path>` | BSON or JSON `{v, t}` document | per-mapping | Full pipeline: validate, persist, triggers, fan-out |

## Astrate to device messages

| Topic | Payload | When |
|---|---|---|
| `<realm>/<device_id>/<interface><path>` | `{v, t}` document (format per device hint) | Server-owned publish; re-send after emptyCache |
| `<realm>/<device_id>/control/consumer/properties` | 4-byte BE size + zlib deflated path list | After emptyCache, after session-present=0, after server-owned unset |

## Topic parsing

The topic is split as `realm / device_id / rest`. The `rest` is matched against the device's introspected interface names by longest-prefix (interface names contain dots, never `/`), and the remainder is the path. An interface match failure produces a rejection metric and optional `device_error` trigger, never a crash.

## Dynamic interface validation

Each message passes through a validation pipeline:

1. **Introspection gate** -- device must declare the interface `name:major` in its introspection.
2. **Trie match** -- path resolved against compiled interface endpoints (segment-wise, O(depth), zero-alloc).
3. **Ownership check** -- device-published data only on `ownership: device` interfaces.
4. **Payload decode** -- BSON or JSON sniffed per message.
5. **Type check** -- value coerced per the endpoint's declared `ValueType`.
6. **Aggregation shape** -- object-aggregated interfaces must arrive as one document of last-level keys.

Failures increment per-reason Prometheus counters, log with device/interface/path, and feed `device_error` triggers.
