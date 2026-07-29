# Milestones

Release-tag-gated goals. Each section names a body of work that must be true by the time
that tag is cut. `.mule/recipes/milestones.md` reads this file to find the **current
target** — the first section below not marked `DONE` — and works toward it: investigating,
filing GitHub issues (`milestone-<tag>` label, alongside `mule`) with sub-issues where the
work splits, and escalating anything that needs a design decision to `.mule/for-giulio.md`.

**This file is Giulio's.** The recipe reads it and may propose edits via
`.mule/for-giulio.md`, but never edits it directly — same rule as `docs/COMPATIBILITY.md`.
Mark a milestone `DONE` yourself once the tag is actually cut.

Order matters: milestones are worked **in order**, lowest tag first. Do not start v3.0
investigation while v2.0 has open, un-escalated gaps — say so and stop instead.

---

## v2.0 — astarte-flow feature parity

Reference: astarte_flow (https://github.com/astarte-platform/astarte_flow) — upstream's
Elixir "Flow" component. It lets ingested data get piped through a graph of processing
blocks (native blocks and containerised ones) before it lands in storage, with pipelines
described as reusable, parametrised graphs.

Scope for this milestone: whatever set of Flow's *capabilities* Astrate needs to expose the
same wire-visible behaviour and operator-facing concepts (pipelines, blocks, native vs.
containerised blocks) — not a port of the Elixir implementation. See
`.mule/recipes/astarte-upstream.md`'s rule: port the idea, restated in Go, never the code.

Status: **DONE** (2026-07-29). Runtime, factory, catalog (incl. filter/map), process
wiring, and `/flow/v1` API are on `main`. **Parity audit** + **product decisions**
recorded 2026-07-29 (`docs/handoff/flow-parity-audit-2026-07-29.md`,
`docs/handoff/flow-v2-decisions-2026-07-29.md`).

**Design A + B landed on `main`** (commit `89145e6`, 2026-07-29): durable flows +
auto_restart (**#41 closed**); named multi-instance + config (**#40 closed**) —
migration `000009`, store, `${config.*}`, API, boot rehydrate. Container block
PoC→MVP (**#43 closed**) — registered in the catalog, usable inside stored
pipelines/named flows; see `docs/handoff/flow-design-b-container-block-2026-07-29.md`.
Blocks discovery (**#39 closed**).
**Not a v2.0 gate:** native Lua / MQTT blocks.

**#42 closed 2026-07-29** (rehydrate edge cases). Triaged all seven candidates against a
live e2e Docker smoke test of #43: one real bug found and fixed on `main`
(`801fc48` — `Registry.Instantiate` leaked an already-started Docker container when a
later block's constructor failed mid-pipeline-build; mutation-tested regression added);
two candidates ("one flow fails at boot, others still start"; "pipeline deleted while
flows reference it") were already correct by design, no code change needed; three
("hot-reload a running pipeline" #44, "partial restart of failed blocks" #45, "update a
running flow's config" #46) split into their own demand-driven backlog issues per the
explicit v2.0 decisions doc; multi-process/HA managers stay out of scope indefinitely
(single-process design throughout).

### Landed (on main as of 2026-07-29)

| Piece | Where | Notes |
|---|---|---|
| FlowMessage wire format | `internal/flow/message.go` | `astarte_flow/message/v0.1` |
| Block / Source / Stopper / graph | `internal/flow/block.go`, `graph.go` | linear chain; Sources skipped in `Run` |
| Pipeline DAG model + validate | `internal/flow/pipeline.go` | serialisable description |
| Flow lifecycle manager | `internal/flow/flow.go` | Start / Stop / List / Shutdown |
| Stream router (lanes, QoS, metrics) | `internal/flow/router.go` | in-order per key |
| Source pump + Stop on teardown | `internal/flow/flow.go` (#37) | pumps `Source.Emit` → router; `Stopper.Stop` after drain |
| AstarteSource block | `internal/flow/blocks/astartesource` | bus → FlowMessage (#27) |
| Pipeline store + migration | `internal/store/pipelines.go`, `000008_pipelines` | realm-scoped DAG JSON (#24) |
| Durable flows store + migration | `internal/store/flows.go`, `000009_flows` | named multi-instance + auto_restart (#40+#41) |
| Config substitution | `internal/flow/substitute.go` | `${config.key}` in block config strings |
| Block factory + instantiate | `internal/flow/factory.go` | `Registry`, `ParseDefinition`, topo order → `[]Block` |
| Built-in catalog | `internal/flow/blocks/catalog.go` + `transform.go` | `astarte_source`, `filter`, `map`, `null_sink`, `log_sink` |
| Process wiring | `cmd/astrate/main.go` | bus + Manager; **boot rehydrate** before listen; shutdown marks stopped |
| Operator HTTP API | `internal/flowapi` | pipelines + named durable `/flows` (a_rma) |

### Closed gaps (all of them — v2.0 is DONE)

1. **~~Design A: durable + named multi-instance~~** — **implemented, #40+#41 closed.**
2. **~~Design B: container block~~** — **implemented (PoC→MVP), #43 closed.** Doc:
   `docs/handoff/flow-design-b-container-block-2026-07-29.md`. PoC transport: HTTP
   (not AMQP).
3. **~~Blocks discovery API~~** — **implemented, #39 closed.**
4. **~~Parity audit + product decisions~~** — **done** 2026-07-29.
   Docs: `flow-parity-audit-2026-07-29.md`, `flow-v2-decisions-2026-07-29.md`.
5. **~~Rehydrate edge cases~~** — **#42 closed 2026-07-29.** One real bug fixed
   (`Registry.Instantiate` container-leak on partial pipeline build, `801fc48`); rest
   already correct or split into #44/#45/#46 (demand-driven, no milestone).
6. **Still out of v2.0 gate (demand-driven):** native Lua/JSONPath blocks; native
   MQTT/Modbus/HTTP poll I/O; full pipeline DSL; `a_f` path wire-compat without a
   client; `http_sink` unless a client needs it.

### Explicitly out of this milestone (tracked elsewhere)

- Dormant trigger types / previous-value cost → #20 (Legion bench) + for-giulio decision
- Group-scoped triggers → #17 (and related #18–19 wire/authz)
- External-bus intake design → #10
- Edgehog client compatibility → #28 (`readonly` / mule) — may integrate upstream FOSS as-is

### Project scope (standing)

Astrate reimplements **free and open-source** software worth recreating for its purpose: a
wire-compatible Astarte-platform reimplementation in Go (lighter components) plus extras
such as AtomVM compatibility. Integrating existing FOSS when it already fits is preferred
over a full rewrite (e.g. original Astarte Dashboard; Edgehog under investigation in #28).

---

## v3.0 — CLEA architecture, piece 1

Reference: https://docs.clea.ai/ — the CLEA architecture is bigger than one milestone, so
each remaining tag (3.0, 4.0, …) takes one piece of it until the whole thing is covered.

**Which piece is v3.0 is not yet decided.** The recipe's first and only job on this
milestone, until it stops being TBD, is to read the CLEA docs, propose a short list of
candidate pieces with a one-line scope each, and write that list to
`.mule/for-giulio.md` as a decision to make. Do not file implementation issues for a piece
that has not been chosen.

Status: not started, scope undecided. Do not investigate while v2.0 is open.

---

## v4.0+ — remaining CLEA pieces

Placeholder. Once v3.0's piece is chosen and scoped, add a `## vX.0 — CLEA, piece N` section
here for the next one, following the same shape. The recipe should propose the next
section's draft (name + one-line scope) as a `for-giulio.md` entry once v3.0 is `DONE`,
rather than leaving this placeholder to rot.

Status: not started, scope undecided.
