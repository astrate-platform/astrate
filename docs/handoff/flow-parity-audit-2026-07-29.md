# Flow v2.0 parity audit — 2026-07-29

Walk of **astarte_flow** operator-visible concepts vs Astrate on `main`
(factory + catalog + `/flow/v1` + filter/map). Reference:

- OpenAPI: `https://github.com/astarte-platform/astarte_flow/blob/master/priv/static/astarte_flow_api.yaml`
- Docs: https://docs.astarte-platform.org/flow/latest/
- Default blocks: `priv/blocks/*.json` in that repo

Policy (from `.mule/milestones.md`): port *capabilities* and operator concepts in Go;
match upstream HTTP only if a real client needs it; container blocks later unless needed.

---

## Concept matrix

| Concept | Upstream | Astrate | Verdict |
|---|---|---|---|
| **Pipeline** (reusable graph) | DSL string (`source`); optional JSON Schema for params | JSON DAG (`blocks` + `connections`); validated + stored | **Parity of idea** — different wire shape (intentional) |
| **Pipeline CRUD** | list / create / get / delete | list / create / get / **update** / delete | **OK** (Astrate superset: PUT) |
| **Flow** (running instance) | Named; config specializes pipeline params; N flows per pipeline | One running instance per `realm/pipeline`; no separate flow name; no config bag | **Gap** — multi-instance + config (see residual) |
| **Flow lifecycle** | create (=start) / get / delete (=stop) | POST/GET/DELETE `/flows` | **OK** |
| **Flow durability** | Process-managed (Elixir app) | In-memory `flow.Manager` only | **Open decision** — for-giulio (flows table) |
| **Blocks catalog** | HTTP CRUD; default + user-defined | Code registry only; no `/blocks` API | **Gap** — discovery API; user-defined deferred |
| **User-defined blocks** | Pipeline-DSL fragment as reusable block | None | **Defer v2.0** (no client) |
| **Container blocks** | Docker + AMQP bridge | None | **Out of v2.0** (milestones) |
| **Auth claim** | JWT `a_f` | Realm Management `a_rma` on `/flow/v1` | **OK for now**; wire-compat only if client |
| **HTTP base path** | `{host}/flow/v1/{realm}/…` | `/flow/v1/{realm}/…` on Astrate process | **Close enough**; full path layout differs from multi-service Astarte |
| **FlowMessage** | `astarte_flow/message/v0.1` | Same schema (`internal/flow/message.go`) | **OK** |

---

## Default block catalog

| Upstream block | Role | Astrate | Notes |
|---|---|---|---|
| `astarte_devices_source` | producer | `astarte_source` | Bus-backed; filters realm/interface/path |
| `filter` | p/c (Lua script) | `filter` | Declarative AND conditions, not Lua — idea ported |
| `update_metadata` | p/c | partly `map` (`set_metadata` / `delete_metadata`) | Key rewrite also on `map` |
| `http_sink` | consumer | — | Optional catalog (# menu 7) |
| `http_source` | producer | — | Defer |
| `mqtt_source` / `mqtt_sink` | I/O | — | Defer (external I/O) |
| `modbus_tcp_source` | producer | — | Defer |
| `random_source` | producer | — | Nice for tests; not operator-critical |
| `lua_map` | p/c | — | Defer scripting |
| `json_path_map` | p/c | — | Defer |
| `to_json` | p/c | — | Defer |
| `sort` | p/c | — | Defer (window reorder) |
| `container` | p/c | — | Out of v2.0 |
| — | — | `null_sink`, `log_sink` | Astrate extras for compose/test |

**Operator-complete path today:** `astarte_source` → `filter`/`map` → `null_sink`/`log_sink`.

---

## API shape details (not bugs)

| Detail | Upstream | Astrate |
|---|---|---|
| Create flow body | `{ name, pipeline, config? }` | `{ pipeline }` only; flow id derived |
| List flows | array of **names** | array of **FlowView** objects |
| Pipeline body | `{ name, source, description?, schema? }` | `{ name, definition }` where definition is DAG JSON |
| Blocks API | full CRUD | missing |

These are intentional Astrate-native choices unless a dashboard/SDK client forces wire match.

---

## Residual gaps (actionable)

### A. Product decisions — **settled 2026-07-29**

See `docs/handoff/flow-v2-decisions-2026-07-29.md`.

1. **Durable flows + auto_restart** — **#41** (design then implement).
2. **Named multi-instance + config** — **#40**.
3. **Containers in scope (PoC→MVP)** — **#43**. Edge cases **#42**.

### B. Filed `milestone-2.0` issues

1. **~~Blocks discovery API~~** — **#39** done locally: `GET .../blocks` + `.../blocks/{type}`.
2. **#40 / #41 / #42 / #43** — flows model + containers (design first).
3. (Optional, demand-driven) **http_sink** catalog block when a client needs outbound webhooks.

### C. Explicitly out of v2.0 gate (updated by product decisions 2026-07-29)

**Superseded:** containers are **in** v2.0, phased — see
`docs/handoff/flow-v2-decisions-2026-07-29.md` and **#43**.

Still out of the v2.0 *gate* (demand-driven later):

- Pipeline DSL parser / user-defined blocks from DSL
- Native Lua / JSONPath scripting blocks (containers can host custom logic)
- External protocol sources/sinks (MQTT, Modbus, HTTP poll) unless a concrete client appears
- JWT claim `a_f` and full upstream path compatibility without a client
- Matching upstream filter *script* semantics (declarative filter is the Go port)

### D. Already done (no action)

- #37 source pump — closed on GitHub; on `main`
- Pipeline store `000008`, factory, process wiring, operator pipelines+flows API

---

## Recommendation for “is v2.0 done?” (updated after decisions)

**Capability floor (filter/map path) already works** on main. Product **DONE** for v2.0
now also requires (after design):

1. Named multi-instance flows + config (**#40**).
2. Durable flows + default auto_restart + loud failure (**#41**).
3. Container block at least to **MVP** (**#43**), via PoC first.
4. Land **#39** blocks discovery if not yet on origin.
5. Docs: DESIGN.md Flow scope fixed 2026-07-29; optional MkDocs Flow page still open.

Do **not** block v2.0 on native Lua, native MQTT, or full DSL parity.
