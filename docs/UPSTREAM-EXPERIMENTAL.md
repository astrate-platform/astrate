# Upstream experimental features register

Astarte's brain is the upstream Astarte team. When upstream ships a feature
marked experimental / unstable / rc-only, Astrate adopts it early (so the wire
behaviour is ready) but tracks it here until upstream promotes or deprecates
it. **Review trigger: every time a future parity milestone targets a newer
upstream version, check each row below against what that release did with the
feature** — promoted → keep and drop the row; deprecated → deprecate in the
same milestone, in sync.

## The version-declaration rule

`APICompatVersion` (internal/realm/service.go) must never claim a level higher
than what Astrate actually implements, and must never be operator-configurable.
It is bumped only in the same change that completes the full surface of that
upstream level. When that happens, every open row below tagged with that level
must be reconciled first — the register is the checklist, so nothing has to be
remembered.

## Rows

| Feature | Issue | Upstream level | Adopted | Reconcile when |
|---|---|---|---|---|
| Required + encrypted interface mapping fields | #67 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
| `async_operation=false` on mutating RM endpoints | #68 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |
