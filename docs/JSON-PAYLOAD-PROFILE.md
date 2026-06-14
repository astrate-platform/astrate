# Astrate JSON Payload Profile

This is the **normative** specification of Astrate's plain-JSON payload profile
for constrained device authors (in particular [AtomVM](https://www.atomvm.net/)
targets on ESP32/RP2040-class MCUs) who want to be first-class Astarte devices
without shipping a BSON codec. It corresponds to `docs/DESIGN.md` §3.5 and the
reference implementation in [`pkg/payload`](../pkg/payload).

A JSON-profile device is an **ordinary Astarte MQTT v1 device** in every
respect — mTLS identity, topic structure, introspection, control channel,
pairing — *except* that data documents are encoded as JSON instead of BSON.
Everything else is wire-identical, so the profile is a strict superset that
upstream Astarte SDKs are unaffected by.

## 1. Transport and topics (unchanged)

Identical to Astarte MQTT v1:

- mTLS to the broker; client certificate `Subject CN = <realm>/<device_id>`;
  MQTT client-ID equal to the CN.
- Base topic `"<realm>/<device_id>"`. Data is published to
  `"<base>/<interface>/<path>"`; control to `"<base>/control/..."`;
  introspection to `"<base>"` itself.
- Introspection is the unchanged `;`-separated list of
  `"<interface>:<major>:<minor>"` triples (plain ASCII; not JSON).

## 2. Format detection (sniffing)

The server detects the encoding per message, structurally and unambiguously
(`pkg/payload/sniff.go`):

```
if   len(p) == 0                                  → empty (property unset / control)
elif len(p) >= 5
     and int32LE(p[0:4]) == len(p)
     and p[len(p)-1] == 0x00                       → BSON
elif first non-whitespace byte == '{'              → JSON
else                                               → rejected
```

A device may therefore publish JSON on exactly the same topics a BSON device
uses; no per-device negotiation is required. Valid JSON text can never collide
with the BSON branch (JSON contains no NUL bytes).

## 3. The JSON data document

```json
{ "v": <value>, "t": "2026-06-10T12:34:56.789Z" }
```

- `v` (**required**) — the value, mapped by the *declared interface type* (see
  §4). The mandatory `{"v": …}` envelope is required: a bare value (`22.5`) is
  **rejected**, preserving symmetry with BSON and keeping `t` unambiguous.
- `t` (**optional**) — the explicit timestamp, either an RFC 3339 string
  (UTC recommended, e.g. `"2026-06-10T12:34:56.789Z"`) **or** an integer number
  of milliseconds since the Unix epoch. Required when (and only when) the
  mapping declares `explicit_timestamp: true`.
- Maximum document size (both encodings): **64 KiB** (configurable via
  `engine.max_payload_bytes`).

### Object-aggregated interfaces

`v` is a JSON object whose keys are the **last path segment** of each mapping in
the aggregation:

```json
{ "v": { "latitude": 45.07, "longitude": 7.69 }, "t": "2026-06-10T12:34:56.789Z" }
```

published to the aggregation path (e.g. `"<base>/com.example.Coords/sensor1"`).

## 4. Type mapping

JSON has one number type, so the interface mapping's declared `ValueType`
disambiguates it. The accepted JSON encodings per Astarte type:

| Astarte type   | JSON encoding of `v`                                              |
|----------------|------------------------------------------------------------------|
| `double`       | JSON number                                                      |
| `integer`      | JSON number (32-bit range)                                       |
| `longinteger`  | JSON number **or** a decimal **string** (use the string for \|x\| > 2^53) |
| `boolean`      | JSON `true` / `false`                                            |
| `string`       | JSON string (UTF-8, ≤ 64 KiB)                                    |
| `binaryblob`   | base64 string, **standard alphabet, padded**                    |
| `datetime`     | RFC 3339 string **or** integer epoch-milliseconds               |
| `*array`       | JSON array of the corresponding scalar encoding                 |

Arrays carry ≤ 1024 elements and must be homogeneous. The same coercion rules
as BSON apply (`pkg/payload/value.go`): `double` rejects NaN/±Inf; an `integer`
JSON number must fit int32; `longinteger` as a JSON number must be an exact
integer.

## 5. Property unset

Publishing an **empty** payload (zero bytes) to a settable property path unsets
that property (`allow_unset` mappings), exactly as for BSON devices.

## 6. Control channel (unchanged framing)

Control payloads keep the Astarte zlib + size-prefix framing for **both** BSON
and JSON devices — zlib inflate is available in the AtomVM standard library, and
changing the control framing would fork the protocol. The deviation is confined
to data documents only.

- **`<base>/control/emptyCache`** — published with an empty payload on connect;
  the server replays server-owned properties and a purge list.
- **`<base>/control/producer/properties`** — the device's set of currently-set
  device-owned property paths, framed as: a **4-byte big-endian** uncompressed
  length, followed by the **zlib-deflated** bytes of the `;`-separated path list
  (`"/p1;/p2;…"`). The server purges device-owned properties not in the list.
- **`<base>/control/consumer/properties`** — the same framing, sent
  server→device, listing currently-set server-owned properties.

## 7. Outbound format (server → device)

The server must send data documents the device can decode. Astrate tracks
`devices.payload_format_hint`:

- Default **`bson`** (the official-SDK assumption).
- Flipped to **`json`** the first time the device publishes a JSON data payload
  (sticky; reset on `emptyCache` only if the device's next data payload is
  BSON).
- Settable explicitly at registration with the additive extension field
  `initial_payload_format`, so a JSON-only device receives JSON from its very
  first server-owned message:

```http
POST /pairing/v1/<realm>/agent/devices
Authorization: Bearer <realm JWT with a_pa>
{ "data": { "hw_id": "<22-char base64url id>", "initial_payload_format": "json" } }
```

The field is ignored by upstream-shaped clients (additive), so a BSON device
simply omits it.

## 8. Minimal device checklist

1. Register (optionally with `initial_payload_format: "json"`), obtain the
   credentials secret.
2. Generate a key + CSR, call the pairing credentials endpoint, receive the
   client certificate; fetch broker URL + realm CA from the info endpoint.
3. mTLS-connect (client-ID = CN), publish introspection, subscribe to
   `"<base>/#"`, publish `emptyCache`.
4. Publish data as `{"v": …, "t": …}` JSON documents on the interface topics;
   send `producer/properties` after `emptyCache`.
5. Decode server-owned messages as JSON (the format you opted into).
