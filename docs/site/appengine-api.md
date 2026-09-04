# AppEngine API

The AppEngine API provides REST endpoints for reading device data, publishing server-owned data, managing groups and aliases, and a live-stream WebSocket/SSE endpoint.

**Base path:** `/appengine/v1/<realm>/`
**Authentication:** JWT with `a_aea` claim.

## Device management

### List devices

```
GET /appengine/v1/<realm>/devices?details=true&limit=N&from_token=<cursor>
```

Returns a paginated list. Pagination uses `body.links.next` cursor (upstream parity).

### Get device

```
GET /appengine/v1/<realm>/devices/<device_id>
```

Returns device status: introspection, connected flag, stats (total received msgs/bytes), timestamps, aliases, attributes.

### Update device

```
PATCH /appengine/v1/<realm>/devices/<device_id>
```

Fields: `aliases`, `attributes`, `credentials_inhibited`.

Setting `credentials_inhibited: true` blocks the device from obtaining new credentials and connecting to the broker.

### Devices by alias

```
GET /appengine/v1/<realm>/devices-by-alias/<alias>
```

## Device data

### Properties

```
GET /appengine/v1/<realm>/devices/<device_id>/interfaces/<interface>[/<path>]
```

Returns the properties snapshot tree for a device's interfaces.

### Datastream queries

```
GET /appengine/v1/<realm>/devices/<device_id>/interfaces/<interface>/<path>
    ?since=<timestamp>
    &since_after=<timestamp>
    &to=<timestamp>
    &limit=N
    &downsample_to=<bucket_duration>
```

- `since` is inclusive; `since_after` is exclusive.
- Default ordering: descending (newest first).
- `downsample_to` maps onto Timescale `time_bucket()`.

### Server-owned publish

```
PUT /appengine/v1/<realm>/devices/<device_id>/interfaces/<interface>/<path>
{ "data": { "v": <value>, "t": "<timestamp>" } }
```

Validates against `ownership: server` interfaces, persists, and publishes to the device via the broker.

### Property unset

```
DELETE /appengine/v1/<realm>/devices/<device_id>/interfaces/<interface>/<path>
```

## Groups

```
GET    /appengine/v1/<realm>/groups
POST   /appengine/v1/<realm>/groups
GET    /appengine/v1/<realm>/groups/<name>

GET    /appengine/v1/<realm>/groups/<name>/devices
POST   /appengine/v1/<realm>/groups/<name>/devices
DELETE /appengine/v1/<realm>/groups/<name>/devices/<device_id>
```

## Live stream

### WebSocket/SSE endpoint

```
GET /astrate/v1/<realm>/socket
```

Feeds real-time device events to consumers. Honours `a_ch` claims as room filters.

### Astarte Channels (Dashboard compatibility)

```
GET /appengine/v1/socket/websocket?realm=<realm>&token=<jwt>
```

Phoenix WebSocket V2 wire format for upstream Astarte Dashboard compatibility. Answers `phx_join`, `watch`, `phx_leave`, and the `phoenix` heartbeat. Events pushed as `new_event`.

## Response envelope

All responses use the Astarte envelope format:

```json
{ "data": <value> }
```

Errors:

```json
{ "errors": { "detail": "<message>" } }
```

Status codes match upstream: 401 (no/bad token), 403 (claim mismatch), 404 (not found), 409/422 (conflict/validation).
