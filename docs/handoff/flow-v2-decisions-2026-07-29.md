# Flow v2.0 product decisions — 2026-07-29

**Status:** decisions recorded; **design drafts written 2026-07-29** (accept then
implement — no mule coding until Giulio accepts each design).

- Design A: `docs/handoff/flow-design-a-named-durable-flows-2026-07-29.md` (#40+#41)
- Design B: `docs/handoff/flow-design-b-container-block-2026-07-29.md` (#43)

**Audience:** power-user product choices, not implementation detail.

Related: parity audit `docs/handoff/flow-parity-audit-2026-07-29.md`.

---

## Goal (what “done” means for you)

Astrate Flow should let you:

1. Save **pipelines** (reusable processing recipes).
2. Run **several named flows** from the same pipeline with different **settings**.
3. Have those flows **come back after Astrate restarts** (unless you opt out).
4. Eventually run **your own logic in a container** on the message path
   (the main reason Flow exists beyond simple filter/map).

Today on main: recipes + one in-memory run per pipeline + built-in filter/map only.

---

## Decision 1 — Remember running flows (durable)

**Chosen: (b) durable records in the database.**

### Intent (Giulio)

- When a flow is first created/started, **auto-restart on process boot defaults to on**.
- Operators can set a flow to **never auto-restart**.
- Restart mechanism can be simple (same path as “start flow”; process-driven is fine).
- If the pipeline is unusable or start fails → **fail loudly** (visible status + error),
  do not hide it.
- Hard corner cases → **separate issues**, not blockers for the first design.

### Refinement (recommended; does not change the decision)

| Topic | Recommendation | Why |
|---|---|---|
| Who restarts? | **Astrate itself on startup** reads durable rows with `auto_restart=true` and starts them | An external “call the API after boot” needs another service; you already own the process |
| Same code path | Rehydrate calls the **same internal start** used by `POST .../flows` | One implementation, API remains for manual start/stop |
| “Pipeline changed” | At start time, resolve pipeline by name; if missing/invalid → status **failed** + error text | Loud failure without silent partial graphs |
| Config drift | Store the flow’s **config snapshot** with the row; do not silently rewrite it | Named multi-instance needs stable config per flow |
| First cut | Persist desired state + last status; rehydrate best-effort; no fancy “rolling update of running graph” | Keeps v1 of durability shippable |

### Out of first cut (track as follow-ups)

- Hot-reload when pipeline definition is edited under a running flow.
- Partial restart of only failed blocks.
- Distributed multi-process managers (single Astrate process assumed).

---

## Decision 2 — Named multi-instance + config

**Chosen: (b) named flows + config bag (upstream mental model).**

### Intent

- A **pipeline** is a reusable blueprint (possibly with placeholders).
- A **flow** has a **unique name per realm**, points at a pipeline, and carries optional
  **config** used to specialize the blueprint.
- Many flows may share one pipeline (e.g. prod vs staging webhook URL).

### Refinement (recommended)

| Topic | Recommendation | Why |
|---|---|---|
| Identity | Manager / store key = `realm` + **flow name** (not pipeline name) | Matches upstream; allows N instances |
| Start body shape | `{ "name", "pipeline", "config"?, "auto_restart"? }` (`auto_restart` default true) | One operator model for decisions 1+2 |
| Parameters | Prefer **string placeholders in block config JSON** (e.g. `${config.webhook_url}`), not a full Elixir-style pipeline DSL | Smaller design surface; same power for most cases |
| Validation | Optional later: JSON Schema on pipeline; first cut = “missing placeholder → loud fail at start” | Fail loud without overbuilding |

### Depends on

Design for decision 1 and 2 should be **one design** (shared schema + API). They are not
independent code paths.

---

## Decision 3 — Containers (and what “done” needs)

### Intent (Giulio)

- Flow should **consume containers** so custom processing is first-class.
- Ship as **PoC → MVP → later features**, not one big bang.
- You asked whether containers (and Lua / MQTT) are required to call the feature done.

### What actually makes Flow “done” for your goal

| Capability | Required for your goal? | Notes |
|---|---|---|
| **Container block** | **Yes** (phased) | “Bring your algorithm” — this is the product heart beyond filter/map |
| **Lua / script blocks** | **No** (not for v2.0 gate) | Useful later; a container can run Lua/Python/anything |
| **Native MQTT / Modbus / HTTP poll sources** | **No** (not for v2.0 gate) | External I/O; containers or other services can own that |
| Named multi-instance + durable rehydrate | **Yes** | Operator model you chose |
| Full upstream DSL / Dashboard parity | **No** | Astrate-native JSON DAG is fine unless a real client forces wire match |

**Push-back (friendly):** treating Lua and MQTT native blocks as “necessary for done”
would inflate scope without unlocking much that a **container PoC** does not already unlock.
Recommend: **containers in scope for v2.0 (phased)**; Lua/MQTT stay demand-driven extras.

### Phasing (recommended)

1. **Design** — how a container block gets messages in/out (upstream uses AMQP; Astrate
   has no AMQP core — design may choose stdio/gRPC/HTTP/NATS/etc. as long as the
   *operator idea* matches: “this step is my image”).
2. **PoC** — one container block type, local Docker, single message path, manual run,
   documented limitations.
3. **MVP** — usable in a stored pipeline + running flow; start/stop with the flow;
   clear failure if image missing / container dies.
4. **Later** — resource limits, pull policies, multi-arch, scaling, richer SDKs, etc.

---

## Recommended work order (meta → design → code)

**Hard sequencing rule:** do **not** implement containers before the **flow model** is
designed (and preferably implemented enough to run a named durable flow). A container is a
**block inside a running flow** — it needs flow name/config/lifecycle and rehydrate semantics
first. Design B may be *drafted* in parallel with Design A, but **Implement B waits on
Implement A** (or at least on Design A accepted + a runnable multi-instance path).

```
[done] Product decisions (this file)
   ↓
Design A — Flows: durable + named multi-instance + config + auto_restart  (#40 + #41)
   ↓
Design B — Container block: transport + lifecycle + PoC scope  (#43)
   ↓   (B design may start after A design exists; B code waits for A)
Implement A (schema, API, rehydrate)
   ↓
Implement B PoC, then MVP
```

Do **not** start implementation for mule until the matching design is accepted
(short design note is enough; not a multi-week whitepaper).

---

## Issue / tracking map

| Work | Tracking |
|---|---|
| Named multi-instance + config | **#40** (decided **b**; design then implement) |
| Durable flows + auto_restart + loud fail | **#41** (pairs with #40 — same design package) |
| Rehydrate / lifecycle edge cases | **#42** (follow-ups after MVP) |
| Container block PoC → MVP | **#43** (design → PoC → MVP) |
| Blocks discovery | **#39** (implemented locally; land/close separately) |
| Lua / MQTT native blocks | **not** v2.0 gate; open only if a client needs them |

---

## Docs notes

- `docs/DESIGN.md` — Flow scope updated 2026-07-29 (v2.0 in scope; was v1 non-goal).
- MkDocs site has no dedicated Flow operator page yet (optional docs phase).
- Parity audit “out of v2.0: containers” is **superseded** by this decision file.
