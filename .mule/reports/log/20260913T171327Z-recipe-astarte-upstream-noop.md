# astarte-upstream watch — noop

Upstream checked 2026-09-13: still no movement.

- Newest tag overall: `v1.4.0-rc.5` (2026-08-20), unchanged since 2026-08-20.
- Newest stable: `v1.3.3` (2026-08-07), unchanged.
- COMPATIBILITY.md compat target remains 1.2.2 (`APICompatVersion`), so upstream
  (1.3.3 stable / 1.4 rc line) is still ahead — but every wire-relevant item from
  v1.3.0..v1.4.0-rc has already been proposed in `.mule/todo.md`
  (categories: probe lines, capability/validation proposals, compat-note-v1.3.x,
  compat-note-v1.4-rc) and the v1.4-rc compat-note line is already queued.

Re-verified the rc.4→rc.5 delta (the newest unreviewed gap): both upstream fixes
touch the DUP value-validation path and neither maps to an Astrate problem:

1. `fix(dup): Make object values only validate type on selected interface` (#2141) —
   upstream was keying object mappings by endpoint leaf name across *all* object
   interfaces, so a value could be validated against another interface's type.
   Astrate validates each object-aggregation document against the *current*
   interface's `CompiledInterface.ObjectLeaves` only (internal/engine/data.go:278,
   internal/engine/serverdata.go:139; pkg/payload/payload.go:102 → decodeBSONObject,
   decodeJSONObject rejects any key absent from the given leaves). The merged-map bug
   class does not exist here. No task proposed.
2. `fix(dup): ensure binaryblob data is correctly validated` — Elixir/Cyanide 2.0
   decode artifact: a raw plain binary was being accepted for a `binaryblob` mapping.
   Astrate's BSON decoder accepts only a BSON binary subtype for binaryblob and
   rejects any other BSON type with `ReasonTypeMismatch`
   (pkg/payload/bson.go:120-125); JSON requires base64 and rejects non-strings.
   Not applicable. No task proposed.

Result: append nothing to `.mule/todo.md`. Empty result is the expected one.

Note: since no task lines were appended, this noop file is the whole deliverable.