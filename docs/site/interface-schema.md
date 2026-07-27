# Interface Schema

Astarte interfaces define the data contracts between devices and the platform. Each interface specifies endpoints, value types, reliability, retention, and ownership.

## Interface types

| Type | Description |
|---|---|
| **datastream** | Time-series data. Values are stored in TimescaleDB hypertables with timestamps. |
| **properties** | Last-value-wins key/value state. Stored in a relational table with upsert semantics. |

## Ownership

| Value | Description |
|---|---|
| **device** | Device publishes data to the server. |
| **server** | Server publishes data to the device. |

## Aggregation

| Value | Description |
|---|---|
| **individual** | Each endpoint maps to a single value. Published as `{"v": <value>, "t": <ts>}`. |
| **object** | Multiple endpoints aggregated into one publish. `v` is an object of last-level keys. |

## Endpoint syntax

- Rooted paths: `/%{sensor_id}/value` or `/temperature`.
- `%{param}` segments are whole-segment placeholders (no partial matches).
- No duplicate or conflicting endpoints within an interface.
- Object-aggregated interfaces: all mappings must be at the same depth, with uniform `explicit_timestamp`, `reliability`, `retention`, and `expiry`. Last-level names must be distinct.

## Value types

### Scalars

| Type | Description |
|---|---|
| `double` | 64-bit floating point |
| `integer` | 32-bit signed integer |
| `longinteger` | 64-bit signed integer |
| `boolean` | True/false |
| `string` | UTF-8 text (<= 64 KiB) |
| `binaryblob` | Raw bytes |
| `datetime` | Timestamp |

### Arrays

Each scalar type has an array variant: `doublearray`, `integerarray`, `longintegerarray`, `booleanarray`, `stringarray`, `binaryblobarray`, `datetimearray`.

Arrays are homogeneous and limited to 1024 elements.

## Mapping attributes

| Attribute | Values | Description |
|---|---|---|
| `value_type` | See above | Declared type of the value |
| `reliability` | `unreliable` / `available` / `guaranteed` | Maps to MQTT QoS 0 / 1 / 2 |
| `retention` | `discard` / `stored` / `volatile` | Whether to retain for offline delivery |
| `expiry` | seconds | Message validity for offline devices |
| `database_retention_policy` | `no_ttl` / `use_ttl` | Whether to apply time-based cleanup |
| `database_retention_ttl` | seconds | TTL when `use_ttl` is set |
| `explicit_timestamp` | boolean | Whether the device provides its own `t` |
| `allow_unset` | boolean | Whether empty payload = property unset |

## Example interface

```json
{
  "interface_name": "org.astarte-platform.example.Values",
  "version_major": 1,
  "version_minor": 0,
  "type": "datastream",
  "ownership": "device",
  "description": "Example datastream interface",
  "mappings": [
    {
      "endpoint": "/%{sensor_id}/value",
      "type": "double",
      "explicit_timestamp": true,
      "reliability": "guaranteed",
      "retention": "stored"
    }
  ]
}
```

## Versioning

- New major versions coexist (a device pins which `name:major` its messages validate against).
- Minor bumps must be additive only -- new mappings, no mutation of existing mapping attributes, same type/ownership/aggregation.
- The `minor` advertised by a device may be <= the installed minor.

## Compilation

On load, each interface is compiled into a `CompiledInterface` with an endpoint trie for O(depth) matching. The compiled cache is held behind an `atomic.Pointer` snapshot (copy-on-write; readers never lock). Invalidation fires on Realm Management CRUD via Postgres `NOTIFY`.
