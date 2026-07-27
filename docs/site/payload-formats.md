# Payload Formats

Astrate accepts both BSON and JSON data documents on the same MQTT topics with the same semantics. This enables constrained devices (AtomVM on ESP32/RP2040) to be first-class Astarte devices without a BSON codec.

## Format detection (sniffing)

The server detects the encoding per message, structurally and unambiguously:

```
if   len(p) == 0                                  -> empty (property unset / control)
elif len(p) >= 5
     and int32LE(p[0:4]) == len(p)
     and p[len(p)-1] == 0x00                       -> BSON
elif first non-whitespace byte == '{'              -> JSON
else                                               -> rejected
```

Valid JSON text can never collide with the BSON branch (JSON contains no NUL bytes). No per-device negotiation is required.

## The data document

Both formats use the same envelope:

```json
{ "v": <value>, "t": "2026-06-10T12:34:56.789Z" }
```

- **`v`** (required) -- the value, mapped by the declared interface type.
- **`t`** (optional) -- explicit timestamp. RFC 3339 string or integer epoch-milliseconds. Required when `explicit_timestamp: true`.
- Maximum document size: **64 KiB** (configurable).
- A bare value (`22.5` instead of `{"v": 22.5}`) is **rejected** -- the envelope is mandatory.

## BSON specifics

- Decoded via `bson.Raw` lookups of `v`/`t` -- no maps, no reflection, near-zero per-message allocations.
- `t` is BSON UTC datetime.
- Uses `go.mongodb.org/mongo-driver/v2/bson` raw-document API.

## JSON profile

Strict, documented profile for constrained device authors:

| Astarte type | JSON encoding of `v` |
|---|---|
| `double` | JSON number |
| `integer` | JSON number (32-bit range) |
| `longinteger` | JSON number **or** decimal string (for values > 2^53) |
| `boolean` | JSON `true` / `false` |
| `string` | JSON string (UTF-8, <= 64 KiB) |
| `binaryblob` | base64 string, standard alphabet, padded |
| `datetime` | RFC 3339 string **or** integer epoch-milliseconds |
| `*array` | JSON array of the corresponding scalar encoding |

### Object-aggregated interfaces

`v` is a JSON object whose keys are the last path segment of each mapping:

```json
{ "v": { "latitude": 45.07, "longitude": 7.69 }, "t": "2026-06-10T12:34:56.789Z" }
```

### Coercion rules

- `double`: int32/int64 widen losslessly; JSON: any number. NaN/Inf rejected.
- `integer`: int32 (or int64/double that fits exactly).
- `longinteger`: int64/int32; JSON number **or** decimal string for JS 2^53 safety.
- `string`: valid UTF-8, <= 64 KiB.
- `binaryblob`: BSON binary / JSON base64 (standard alphabet, padded).
- `datetime`: BSON UTC datetime / JSON RFC 3339 or epoch-ms.
- Arrays: homogeneous, <= 1024 elements.

## Outbound format (server to device)

The server must send data documents the device can decode:

- **Default `bson`** (official SDK assumption).
- **Flipped to `json`** the first time a device publishes a JSON data payload (sticky; reset on `emptyCache` only if next payload is BSON).
- **Settable at registration** with the additive extension field:

```http
POST /pairing/v1/<realm>/agent/devices
{ "data": { "hw_id": "...", "initial_payload_format": "json" } }
```

Control payloads (`consumer/properties`) keep the zlib+size framing for both device types.

## Property unset

Publishing an **empty** payload (zero bytes) to a settable property path unsets that property (`allow_unset` mappings).
