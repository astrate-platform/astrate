# Handoff — Flow after product decisions (2026-07-29)

Use this prompt to continue general Astrate work (not a docs-only phase).

````text
I'm working on the Astrate project in ~/astrate on branch main.

Before changing files, read:
- ~/astrate/docs/handoff/session-2026-07-29-flow-factory-handoff.md  (this file)
- ~/astrate/docs/handoff/flow-design-a-named-durable-flows-2026-07-29.md  (Design A — accept then implement)
- ~/astrate/docs/handoff/flow-design-b-container-block-2026-07-29.md  (Design B — draft; code after A)
- ~/astrate/docs/handoff/flow-v2-decisions-2026-07-29.md  (settled product choices)
- ~/astrate/docs/handoff/flow-parity-audit-2026-07-29.md
- ~/astrate/docs/handoff/README.md
- ~/astrate/.mule/milestones.md
- ~/astrate/.mule/for-giulio.md  (open decisions only)

## Product decisions (settled — do not re-open without Giulio)

Recorded in `docs/handoff/flow-v2-decisions-2026-07-29.md`:

1. **Durable flows** — DB records; `auto_restart` **default true**; optional never;
   process rehydrates on boot (same start path as API); fail loudly. Issues **#41**, edge **#42**.
2. **Named multi-instance + config** — flow name + pipeline + config bag. Issue **#40**.
3. **Containers in v2.0** — phased PoC → MVP → later. Issue **#43**.
   Native Lua/MQTT blocks are **not** a v2.0 gate.

**Hard rule:** do **not** implement containers before the flow model (names, config,
durable auto-restart). Container code waits on Implement A.

## Design status (2026-07-29 session)

| Design | File | Status |
|---|---|---|
| **A** — durable named multi-instance (#40+#41) | `flow-design-a-named-durable-flows-2026-07-29.md` | **Implemented** (migration `000009`, store, substitute, API, boot rehydrate) |
| **B** — container block (#43) | `flow-design-b-container-block-2026-07-29.md` | **Draft written — await accept; implement after A** |

Design A accepted as defaults (DELETE=stop+remove; `${config.key}` flat; failed
start leaves no Manager entry; rehydrate before HTTP listen; list durable rows).

## Split status: committed vs working tree

**On `main` (may be ahead of origin by 2+):** factory, catalog filter/map, `/flow/v1`,
source pump (#37 closed), pokemon #38.

**Likely uncommitted / local:**
- #39 blocks discovery (`internal/flow/blocks/info.go` + flowapi)
- Decision + audit + handoff + **Design A/B** docs under `docs/handoff/`
- `.mule/for-giulio.md`, `.mule/milestones.md` updates

Push/commit only with confirmation.

## Practical menu

| Priority | Action |
|---|---|
| **0** | Commit/push Design A implementation + prior Flow work when Giulio wants |
| **1** | Close/comment GitHub **#40** + **#41** after smoke/integration |
| **2** | **Giulio: accept Design B** (section 10) |
| **3** | After B accepted: container PoC then MVP (#43) |
| **4** | #42 edge cases after durable MVP |
| **5** | Optional: MkDocs Flow operator page; DESIGN.md Flow scope already fixed |
| **6** | Legion #20 / mule #13 / #17–19 when wanted |

Short answer: **Design A is implemented** (named durable flows + boot rehydrate).
Next product gate is **Design B** (#43 containers). Do not implement containers
until B is accepted.

### filter / map config (operator reference)

`filter`: key_prefix, key_contains, type, metadata (AND; drop = zero outputs).
`map`: key template `{key}` / `{metadata.*}`, set_metadata, delete_metadata.

Rules:
- Prefer small, test-backed changes when coding starts.
- Do not commit unless asked; do not push without confirmation.
- End session: update this handoff (or successor) and name the next file to read.
````

## Menu only (quick scan)

| Priority | Action |
|---|---|
| **0** | Commit/push local work (confirm) |
| **1** | Accept Design A |
| **2** | Implement #40+#41 |
| **3** | Accept Design B |
| **4** | Container PoC/MVP after A |
| **5** | #42 later |
| **6** | Optional MkDocs Flow page |
| **7** | Other open issues |

## Landed / decided

| Piece | Notes |
|---|---|
| Runtime + factory + catalog + API | on main |
| Parity audit | `flow-parity-audit-2026-07-29.md` |
| Product decisions | `flow-v2-decisions-2026-07-29.md` |
| **Design A** | `flow-design-a-named-durable-flows-2026-07-29.md` — **await accept** |
| **Design B** | `flow-design-b-container-block-2026-07-29.md` — **await accept** |
| #39 discovery | local code; land/close when pushed |
| #40 multi-instance | design draft ready |
| #41 durable + auto_restart | design draft ready (same package as #40) |
| #42 edge cases | tracking only |
| #43 containers | design draft ready; code after A |

## Also open (context)

| Issue | Notes |
|---|---|
| #28 | Edgehog (readonly, mule) |
| #20 | previous-value bench (Legion) |
| #13 | dashboard error_names |
| #17–19 | enhancements |
| #10 | design question |
| #1 | wontfix — ignore |

## Next file to read

Same file, or after accept: implement from
`docs/handoff/flow-design-a-named-durable-flows-2026-07-29.md`.
