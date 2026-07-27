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

Status: not started. No `milestone-2.0` issues filed yet.

---

## v3.0 — CLEA architecture, piece 1

Reference: https://docs.clea.ai/ — the CLEA architecture is bigger than one milestone, so
each remaining tag (3.0, 4.0, …) takes one piece of it until the whole thing is covered.

**Which piece is v3.0 is not yet decided.** The recipe's first and only job on this
milestone, until it stops being TBD, is to read the CLEA docs, propose a short list of
candidate pieces with a one-line scope each, and write that list to
`.mule/for-giulio.md` as a decision to make. Do not file implementation issues for a piece
that has not been chosen.

Status: not started, scope undecided.

---

## v4.0+ — remaining CLEA pieces

Placeholder. Once v3.0's piece is chosen and scoped, add a `## vX.0 — CLEA, piece N` section
here for the next one, following the same shape. The recipe should propose the next
section's draft (name + one-line scope) as a `for-giulio.md` entry once v3.0 is `DONE`,
rather than leaving this placeholder to rot.

Status: not started, scope undecided.
