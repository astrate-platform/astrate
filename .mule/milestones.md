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

Status: **in progress** (not DONE). Runtime, factory, catalog (incl. filter/map), process
wiring, and `/flow/v1` API are on `main`. Remaining: flows-table decision + parity audit.
Refreshed 2026-07-29.

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
| Block factory + instantiate | `internal/flow/factory.go` | `Registry`, `ParseDefinition`, topo order → `[]Block` |
| Built-in catalog | `internal/flow/blocks/catalog.go` + `transform.go` | `astarte_source`, `filter`, `map`, `null_sink`, `log_sink` |
| Process wiring | `cmd/astrate/main.go` | shared `e.Bus()` + `flow.Manager`; shutdown drains flows |
| Operator HTTP API | `internal/flowapi` | `/flow/v1/{realm}/pipelines` + `/flows` (a_rma) |

### Remaining gaps (file / keep as `milestone-2.0` issues)

Work these roughly in order. Escalate design questions to `.mule/for-giulio.md` rather than
guessing wire shape.

1. **Flow runtime persistence (optional / decide)** — issue-25 sketched a `flows` table
   (`000009`); only `pipelines` landed. Decide whether running-flow records need durable
   status or in-memory manager is enough for v2.0.
2. **Parity audit** — walk astarte_flow's operator-visible concepts (pipelines, blocks,
   flows, container blocks) against what Astrate exposes; file residual gaps or escalate
   "won't do for v2.0" decisions. Match upstream HTTP paths only if a real client needs them
   (current surface is Astrate-native `/flow/v1/...`). Containerised blocks later unless a
   client needs them for parity.

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
