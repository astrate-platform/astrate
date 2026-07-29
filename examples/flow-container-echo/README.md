# flow-container-echo

Minimal container image for the Astrate Flow **container** block PoC (#43).

Implements the HTTP message contract used by `internal/flow/blocks/container`:

| Method | Path | Behaviour |
|--------|------|-----------|
| `GET` | `/healthz` | `204` — readiness |
| `POST` | `/v1/message` | Body = one FlowMessage JSON (`astarte_flow/message/v0.1`); response `200` + same JSON (echo), or `204` to drop |

Env:

- `ASTRATE_FLOW_CONFIG` — opaque JSON from the pipeline node’s nested `config` (logged at start).
- `PORT` — listen port (default `8080`).

## Build

```bash
docker build -t astrate/flow-container-echo:poc .
```

## Manual curl

```bash
docker run --rm -p 18080:8080 astrate/flow-container-echo:poc

curl -sS -X POST http://127.0.0.1:18080/v1/message \
  -H 'Content-Type: application/json' \
  -d '{
    "schema": "astarte_flow/message/v0.1",
    "key": "demo",
    "type": "string",
    "data": "hello",
    "timestamp_us": 0
  }'
```

Drop a message (filter semantics):

```bash
curl -sS -o /dev/null -w '%{http_code}\n' -X POST http://127.0.0.1:18080/v1/message \
  -H 'Content-Type: application/json' \
  -d '{
    "schema": "astarte_flow/message/v0.1",
    "key": "demo",
    "type": "string",
    "data": "x",
    "metadata": {"echo_drop": "1"},
    "timestamp_us": 0
  }'
# → 204
```

## Pipeline snippet

```json
{
  "name": "enrich",
  "block_type": "container",
  "config": {
    "image": "astrate/flow-container-echo:poc",
    "config": { "note": "opaque to the container" }
  }
}
```

Limitations (PoC): local Docker only, HTTP request/response, no resource limits,
no registry auth beyond what your Docker daemon already has. See Design B.
