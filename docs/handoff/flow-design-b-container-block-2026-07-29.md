# Design B — Container block transport + PoC → MVP (#43)

**Status:** draft for acceptance (implement **after** Design A is implemented enough
to run a named durable flow). Design may be reviewed in parallel with Design A.  
**Date:** 2026-07-29  
**Issue:** #43  
**Decisions:** `docs/handoff/flow-v2-decisions-2026-07-29.md`  
**Depends on:** Design A (flow name, config, lifecycle, rehydrate)

---

## 1. Goal

Let operators put **their own image** on the message path:

```text
… → filter/map → [container: my-algo] → sink …
```

Same operator idea as upstream Astarte Flow’s container block; **not** a port of the
Elixir/AMQP bridge. Astrate has no RabbitMQ core (DESIGN.md: in-process channels).

Phasing (product decision):

1. **Design** (this file)
2. **PoC** — one block type, local Docker, single message path, manual run, documented limits
3. **MVP** — usable in stored pipeline + named durable flow; start/stop with the flow;
   loud fail if image missing / container dies
4. **Later** — resource limits, pull policies, multi-arch, scaling, richer SDKs

Native Lua / MQTT blocks are **not** a v2.0 gate (containers cover custom logic).

---

## 2. Operator model

### 2.1 Block type

Catalog type: `container` (or `docker` — prefer **`container`** to match upstream naming).

Example pipeline node:

```json
{
  "name": "enrich",
  "block_type": "container",
  "config": {
    "image": "ghcr.io/example/flow-enrich:1.0",
    "config": { "threshold": 0.5 }
  }
}
```

| Config key | Required | Notes |
|---|---|---|
| `image` | yes | Image reference (local or registry) |
| `config` | no | Opaque JSON object passed to the container at start (env or file) |
| `transport` | no | Default from process/global; PoC may hardcode one |

Flow-level Design A config can supply `image` via `${config.image}` after Design A.

### 2.2 Role in the graph

- **Producer/consumer (p/c):** one input stream of `FlowMessage`, zero-or-more outputs
  per input (same contract as filter/map: drop = zero outs; fan-out allowed later).
- Linear graphs only in current `BlockGraph` — container is one stage in the chain.

### 2.3 Lifecycle tied to the flow

| Flow event | Container |
|---|---|
| Flow start (after Instantiate) | Create + start container; fail flow start if image pull/create fails |
| Message processing | Send msg in; wait for response(s) or async depending on transport PoC |
| Flow stop / Shutdown | Stop + remove container (best-effort; log failures) |
| Container dies mid-run | Mark flow **failed** (loud); stop accepting; surface error |

Do **not** leave orphan containers after clean stop (label + remove).

---

## 3. Message contract (container-facing)

Containers should not need to speak Go. Wire format = existing FlowMessage JSON:

- Schema: `astarte_flow/message/v0.1` (`internal/flow/message.go`)
- One JSON object per message (UTF-8)

**Inbound to container:** full `FlowMessage` JSON.  
**Outbound from container:** zero or more `FlowMessage` JSON objects (same schema).  
Dropping a message = zero outputs (filter semantics).

Optional later: binary framing, batching, side-channel metrics.

---

## 4. Transport choice

Upstream uses **AMQP** between Flow and the container. Astrate will **not** add RabbitMQ
just for this block.

### Options considered

| Option | Pros | Cons |
|---|---|---|
| **A. HTTP** (container listens; block POSTs message, response body = outs) | Easy to implement in any language; debuggable | Latency per msg; need port publish + readiness |
| **B. stdio** (NDJSON on stdin/stdout) | No ports; simple local Docker `-i`; good PoC | Harder multi-stream; buffering; no separate health port |
| **C. gRPC** | Streaming, typed | Heavier SDK burden for “bring any image” |
| **D. NATS / other bus** | Fits “external bus” story | New dependency; overkill for PoC |

### Recommendation

| Phase | Transport |
|---|---|
| **PoC** | **HTTP** request/response on a container-private port (e.g. `8080` inside network / published to localhost ephemeral) |
| **MVP** | Same HTTP path, hardened (timeouts, max body, start health check) |
| **Later** | Optional stdio or gRPC as alternate `transport` values |

**Why HTTP for PoC:** language-agnostic; matches “bring your algorithm” (Flask, Go
`net/http`, Node); easy to curl during debug; no AMQP. Stdio remains a documented
alternative if HTTP proves awkward in Docker-on-Mac networking — decide at PoC review,
not before A.

### HTTP sketch (PoC)

- Container **must** expose `POST /v1/message` (or `/message`):
  - Request body: one FlowMessage JSON
  - Response `200`: either one FlowMessage JSON, or a JSON array of messages
  - Response `204` / empty array: drop
  - Non-2xx: block treats as processing error (log + metric; policy: fail flow vs skip — **PoC: log and drop with metric**; **MVP: configurable, default fail-loud after N errors**)
- Astrate block waits with a **timeout** (config default e.g. 5s).
- On flow start: `docker run` (or API) with env:
  - `ASTRATE_FLOW_CONFIG=<json>` of nested `config`
  - Labels: `astrate.flow=1`, `astrate.realm=…`, `astrate.flow_name=…`, `astrate.block=…`

### Docker API

- Prefer Docker Engine API via local socket (`unix:///var/run/docker.sock`) or
  `DOCKER_HOST`.
- PoC may shell out to `docker` CLI if faster to land; MVP should use a Go client
  (`github.com/docker/docker/client` or lightweight wrapper) for cancel/cleanup.
- **No Kubernetes operator in v2.0** — local Docker only for PoC/MVP.

---

## 5. Package layout (proposed)

```text
internal/flow/blocks/container/
  block.go          # Block + Stopper implementation
  docker.go         # create/start/stop/remove
  httpbridge.go     # POST message ↔ FlowMessage
  block_test.go     # unit with fake round-tripper
```

Register in `blocks.DefaultRegistry()` as `"container"`.

`blocks.Info` docs: image required; local Docker; limitations list for PoC.

### Process dependencies

Extend `flow.Deps` only if needed:

```go
// Optional; nil → container block fails Instantiate with clear error
Docker DockerRunner // interface{ Run(...); Stop(...) }
```

Or construct via config-only + package-level client. Prefer injectable interface for tests.

---

## 6. PoC scope (must / must-not)

### Must

1. `container` block type registered and documented via `/flow/v1/.../blocks`.
2. Start container from `image` when flow starts; HTTP round-trip one message.
3. Stop/remove container when flow stops.
4. Manual test recipe in handoff or `docs/` snippet: tiny echo image or `nginx`+sidecar
   is not enough — ship a **minimal example** under `examples/flow-container-echo/`
   (Dockerfile + 20-line HTTP server that echoes FlowMessage).
5. Loud fail if Docker unavailable or image missing at start.

### Must not (PoC)

- Registry auth beyond what local Docker already has
- CPU/memory limits, restart policies, multi-container sidecars
- AMQP compatibility with upstream container images
- Auto-pull policy UI; `docker pull` best-effort with clear error
- HA / multiple Astrate hosts scheduling the same container

### MVP adds

1. Works inside Design A named durable flow + rehydrate (container comes back on boot
   with the flow).
2. Health wait: block start waits until HTTP accepts or timeout → flow `failed`.
3. Config timeouts, max response bytes.
4. Metrics: messages in/out, container start failures, processing errors.
5. Cleanup of orphaned containers with Astrate labels on process boot (best-effort).

---

## 7. Security notes (MVP floor)

- Container runs with Docker defaults unless config later tightens; document that
  operators trust the image.
- Do not mount host Docker socket into the user container.
- Network: bridge + publish only the HTTP port to localhost (or attach to an internal
  network the block can reach). Prefer **not** host network.
- Realm isolation: labels only in PoC; no multi-tenant hard isolation claim.

---

## 8. Sequencing vs Design A

```text
Design A accepted → Implement A (named durable flows)
Design B accepted ↗ (may be earlier)
        ↓
Implement B PoC (can use manual StartFlow without full durability if needed for speed,
                 but MVP requires A)
        ↓
Implement B MVP on durable named flows
```

**Hard rule:** do not merge container as the “default path to custom logic” before
flows have names + config + durable restart semantics operators can rely on.

---

## 9. Implementation plan (after A + B accepted)

| Step | Work |
|---|---|
| **B1** | Echo example image + HTTP contract doc |
| **B2** | `container` block + Docker start/stop + HTTP bridge (PoC) |
| **B3** | Catalog Info + service tests with mocked Docker/HTTP |
| **B4** | MVP: health wait, timeouts, metrics, rehydrate-safe Stopper |
| **B5** | Optional: orphan cleanup on boot |

### Acceptance (PoC)

- Pipeline `astarte_source → container(echo) → log_sink` processes a live message.
- `go test` for bridge logic without requiring Docker in unit tests; one integration
  tag test optional when Docker available.

### Acceptance (MVP)

- Same path as durable named flow with `auto_restart`; after process restart, container
  runs again.
- Missing image → flow status `failed` + error_message.

---

## 10. Open points for Giulio

Defaults recommended in brackets:

1. **PoC transport = HTTP** (not stdio)? **[yes]**
2. **Block type name `container`?** **[yes]**
3. **PoC may shell `docker` CLI; MVP uses API client?** **[yes]**
4. **Processing error policy PoC = log+drop+metric?** **[yes]**
5. **Ship `examples/flow-container-echo`?** **[yes]**

---

## 11. Related

| Doc / code | Role |
|---|---|
| Design A | Flow identity, config, start path |
| `internal/flow/message.go` | FlowMessage schema |
| `internal/flow/blocks/catalog.go` | Register block |
| `docs/DESIGN.md` | No AMQP core; single binary |
| Upstream container block | Idea only; AMQP not required |
