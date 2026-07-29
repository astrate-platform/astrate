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

## v3.0 — HALTED

Status: halted, not `DONE`. Do not advance past v2.0 into v3.0 territory until this is
lifted.

---

## v4.0 — Edgehog fleet management

Reference: Clea Edgehog (https://docs.edgehog.io/) is built directly on top of Astarte —
it uses standard Astarte interfaces (see `astarte_interfaces.html` in its docs) for
device-cloud communication, so it should in principle work unmodified against any backend
that implements Astarte's wire protocol, Astrate included.

**Not yet verified against actual code — only against Edgehog's own docs.** Before scoping
this milestone, the first task is a `readonly` investigation issue: stand up Edgehog against
a running Astrate instance and confirm the interfaces/APIs it actually depends on (realm
management, pairing, interface introspection, trigger delivery, etc.) work as Edgehog
expects.

- If compatible: v4.0 scope narrows to closing whatever gaps that investigation finds
  (should be small — adoption, not a build).
- If not compatible: escalate the concrete incompatibilities to `.mule/for-giulio.md` and
  evaluate options there (extend Astrate's API surface vs. something else) rather than
  guessing here.
- Either way, once Edgehog compatibility is settled, evaluate **Clea OS** as the natural
  next piece after Edgehog (v5.0) — it's the device-side agent that talks to both Astarte
  and Edgehog, so it only makes sense once Edgehog is in place.

Status: not started, scope conditional on the Edgehog-compatibility investigation above.

---

## v5.0+ — remaining CLEA pieces (Clea OS, ...)

Placeholder — expected next piece is Clea OS, contingent on v4.0's Edgehog-compatibility
finding per above. Once v4.0 is `DONE`, add a `## v5.0 — Clea OS` section here following the
same shape.

Status: not started, scope undecided.
