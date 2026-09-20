# Code review — internal/engine (triggers + forward) — 2026-09-20

Second pass over `internal/engine`. The 09-06 review covered the core data
paths; this run targets the files the rotation had not read since: the
`triggers` subpackage's validation/transient/error-name machinery and the
`forward/` extension package, plus the delivery path in `actions.go`.

## What I read

- `internal/engine/triggers/actions.go` (610 lines) — parseAction, Executor,
  Enqueue/deliver/forward/webhook, policy retry loop
- `internal/engine/triggers/validation.go` (163 lines) — URL/header/template
  limits, FieldErrors/TriggerErrors renderers
- `internal/engine/triggers/transient.go` (22 lines) — CompileCondition
- `internal/engine/triggers/errorname.go` (79 lines) — upstream error_name
  closed set + astrateToUpstream mapping
- `internal/engine/triggers/match.go` (definition/compile skeletons, lines
  230-345), `policy.go` (Decide contract)
- `internal/engine/forward/{http,envelope,nats}.go` + http_test.go
- Tests: `errorname_test.go`, `transient_test.go`, `forward/http_test.go`,
  `actions_test.go` (policy suite + forwarder seam tests)
- Callers: `internal/engine/cache.go:322` (AttachPolicy), `cmd/astrate/main.go`
  `newForwarder`, `internal/config/config.go` forward validation
- `pkg/payload/payload.go:148` (ReasonMissingRequired emission)

## What I found (worth proposing)

### 1. `missing_required` is the only `astrateToUpstream` key with no test row

`internal/engine/triggers/errorname.go:64` maps the upstream-1.4 reject reason
`missing_required` → `unexpected_object_key`. `TestUpstreamErrorNameMapping`
(errorname_test.go:12-31) covers all 19 other map keys verbatim; a regression
that drops or remaps this entry (the only 1.4-era one — `errorname.go:62-63`
notes it is outside dashboard 1.2.2's closed set, so it is the most likely to
be "reviewed into" the wrong value later) passes every suite. Add the one
table row, plus an invariant test that every `astrateToUpstream` value is a
member of `UpstreamErrorNames()` — that invariant also catches a
typo-then-copied mapping value, which the fixed-row table alone cannot.

### 2. The custom-action (forward) delivery path never consults an attached policy

`webhook()` implements the full policy contract (actions.go:505-579):
`policy.Decide`, `RetryTimes+1` cap, backoff. `forward()` (actions.go:469-485)
does a single attempt under a `RequestTimeout` context and skips the policy
entirely — no `Decide`, no retry, no discard-by-policy. But `Enqueue`'s
`maximum_capacity` gate (actions.go:335-346) and `deliver`'s `event_ttl`
(actions.go:406-413) DO apply to custom actions: the trigger's policy is
checked with `d.Trigger.Policy()` before the custom/webhook branch at
actions.go:432-436. So a trigger with a custom action + retry policy gets
capacity and TTL semantics but silently not retry/discard semantics. The
executor's own comment (actions.go:494-499) says a policy "governs every
retry/discard decision" with no carve-out. Intent should be confirmed: either
route the forward path through the same Decide loop (retry on
`StatusTransport`/treat as server error, honour discard), or document and pin
that forwarder deliveries are intentionally single-shot. The minimum tick is
a test that pins today's behaviour for a custom action with a policy attached
so the decision is forced.

## What I decided NOT to propose, and why

- **Unbounded `io.Copy(io.Discard, resp.Body)` drain in forward/http.go:117-120**
  vs the bounded `io.LimitReader(resp.Body, 1<<20)` in the webhook path
  (actions.go:607). Divergence between siblings, but the drain is bounded in
  time by the `RequestTimeout` context the executor passes (actions.go:476-477),
  and the 1 MiB cap on the webhook side is about not reading an irrelevant body
  forever — the forwarder reads at most until ctx expiry. Not proposing.
- **NATS forwarder ignores its ctx** (nats.go:44 `_ context.Context`). NATS
  `Publish` is non-blocking best-effort; nothing to cancel. Not a gap.
- **`validMethod` accepts any RFC-7230 token while config.go:355-361 restricts
  to the seven verbs.** config validate runs first and is the gate; the looser
  forward check is the boot-time backstop. Deliberate layering, no drift.
- **`transient.go` / `CompileCondition`** — fully hand-tested (transient_test.go:
  unknown types, missing fields, invalid ops, non-JSON, and the nil-Action pin).
- **`validation.go` helpers** — exercised through parseAction's modern/legacy
  branches; the limits are probe-frozen upstream values with issue refs.

Finding two items in a swept area is the yield; the rest of the triggers
package reads correct and well-covered.