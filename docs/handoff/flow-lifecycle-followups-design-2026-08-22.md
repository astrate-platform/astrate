# Design note: flow lifecycle follow-ups (#44 hot-reload, #45 partial restart, #46 config update)

Date: 2026-08-22
Status: proposal — decisions flagged as **[DECIDE]** where there is a real choice.
Context: closes out the open questions in #44, #45, #46. Builds on
`flow-v2-decisions-2026-07-29.md` ("no fancy rolling update of running graph" first cut).

## Shared groundwork: one internal "restart instance" primitive

Today the only start path is `Service.startFlowInstance`
(`internal/flowapi/service.go:253`), used by POST-create and boot rehydrate. Stop is only
reached destructively via `DeleteFlow` (`service.go:327`), which also deletes the durable
row. There is no way to stop a live graph while keeping its durable row and starting again.

All three issues need the same primitive:

```
restartFlowInstance(ctx, realm, realmID, name):
  1. mgr.StopFlow(instanceID) + Unregister   // drains pump, calls Block.Stop() on blocks
  2. re-run the resolve half of startFlowInstance:
     GetPipeline -> SubstituteConfig -> checkBlockTypes -> ParseDefinition -> Instantiate
  3. status "creating" -> mgr.StartFlow -> status "running"
     on error: markFailed (old graph is already gone)
```

Extraction: split `startFlowInstance` into `resolveAndBuild(...)` (steps of the resolve
half) plus the existing orchestration; `restartFlowInstance` = `DeleteFlow`'s stop-half +
`resolveAndBuild`. Durable row is never touched except runtime status/config columns.

Semantics of the swap: **stop-then-start (drain and rebuild)**, not blue-green. If the new
definition/config fails validation or instantiation after teardown, the flow lands in
`failed` with the error persisted — the operator re-fixes and hits reload/restart again.
**[DECIDE]** Blue-green (build new graph first, tear down old only on success) is doable
by instantiating under a temp instance key, but doubles peak resource usage (container
blocks run two containers briefly) and complicates source dedup. Recommendation: defer;
ship stop-then-start.

In-flight messages during swap: whatever `Manager.StopFlow` already guarantees
(`internal/flow/flow.go:248` — cancel pump, drain, then `Stopper.Stop()`). Sources
resubscribe fresh on rebuild; delivery stays at-least-once as today. No new guarantee is
introduced.

---

## #44 Hot-reload of a running flow when its pipeline changes

### Decisions

1. **Yes, an operator can pick up a new definition without manual delete+recreate** — but
   *explicitly*, not automatically on pipeline PUT.
   - New endpoint: `POST /flow/v1/{realm}/flows/{name}/reload` → runs
     `restartFlowInstance`, which re-resolves the pipeline by name and substitutes the
     flow's stored config. This matches how boot rehydrate already works, so behavior is
     consistent: reload == what would have happened on next boot.
2. **Visibility fix ships regardless**: `UpdatePipeline` (`service.go:141`) response gains
     a `referencing_flows` field listing durable flows whose `pipeline_name` matches, so a
     silent no-op edit stops being silent. This needs only a store query joining
     `flows.pipeline_name` (the live `Manager` does not retain the recipe name — see
     `factory.go:207`, historical `PipelineID` is the instance ID). No walking of live
     instances required for v1.
3. **No automatic propagation** of a pipeline PUT to referencing flows, and no rolling
   per-block update — reaffirming the v2 decision. Automatic propagation is the dangerous
   version (a bad edit takes down N running flows at once); explicit reload per flow keeps
   blast radius at one.

### What reload does NOT do

- Does not migrate in-flight state between old and new blocks. A block whose config changed
  is a brand-new instance starting empty (e.g. an HTTP poller re-fetches from scratch).
- Does not preserve the resolved definition durably (unchanged from today: only raw
  pipeline + config snapshot are stored; the resolved graph lives in memory).

---

## #45 Partial restart of only failed blocks

### The missing prerequisite: nobody detects block death mid-run

Runtime block errors today are logged and counted, not fatal (`router.go:190`,
"Block processing errors... not fatal"). `failed` status is only ever written at
start-time. So "restart the failed blocks" has nothing to hook into yet — first we must
*detect and record* a dead block:

1. **Death detection**: container block gains a container-exit watcher (it already owns the
   container lifecycle, `blocks/container/block.go`). Unexpected nonzero exit (vs. a clean
   `Stopper.Stop()` shutdown) raises a new `BlockFatal` signal to the Manager.
   **[DECIDE]** what counts as fatal in v1: recommendation = container exited (any code,
   incl. daemon-reported death) while the flow was `running`; transient errors inside
   message processing (HTTP timeouts etc.) stay non-fatal as today.
2. **Record**: flow transitions to `failed` with the failing block's name persisted (new
   nullable column `failed_block` on the flows row, cleared on successful start).

### Restart granularity: full-flow restart with backoff in v1, true per-block later

The issue's own analysis says the pain is source resubscription, not rebuilding the failed
block. Given that:

- **Auto-restart flows**: on `BlockFatal`, automatically re-run `restartFlowInstance` with
  exponential backoff (reuse/extend the auto-restart flag; cap backoff, keep retrying like
  rehydrate does at boot). Whole graph rebuilds — simple, and preserves the factory
  invariant that a partially built graph never leaks (`factory.go:115`).
- **True per-block restart** (keep surviving blocks, replace one): deferred until there is
  evidence the resubscription cost matters. When specced, it should reuse `Instantiate`'s
  cleanup discipline: stop the survivors only if replacement fails.

**[DECIDE]** whether the manual API surface also grows a `POST /flows/{name}/restart`
(restart a failed/stopped flow without delete+create — falls out of `restartFlowInstance`
almost for free). Recommendation: yes, add it alongside; it is the manual escape hatch for
non-auto-restart flows.

---

## #46 Update a flow's config without delete+recreate

### Decision

- New endpoint: `PUT /flow/v1/{realm}/flows/{name}/config` body = the config object.
- Behavior:
  1. Validate by dry-running `SubstituteConfig` against the flow's current pipeline
     (missing keys fail loudly, same as create — `substitute.go:87`).
  2. Persist the new config snapshot on the durable row.
  3. If the flow is currently live, immediately run `restartFlowInstance` so the running
     graph picks it up; if stopped/failed, just persist — next start uses it.
- **No hot-appliable subset in v1.** Even "simple" values reach blocks through
  `${config.*}` substitution at instantiation time, and the container block bakes them into
  `ASTRATE_FLOW_CONFIG` at container start (`docker.go:90`) — unchangeable without a
  container restart anyway. A uniform "config change ⇒ rebuild" rule is honest and cheap.
- Implementation-wise this is #44's reload with a different trigger and a config write in
  front; both land together or #46 lands right after #44.

---

## Suggested landing order

1. Refactor: extract `resolveAndBuild` / add `restartFlowInstance` (no behavior change).
2. #46 + #44 endpoint + pipeline-update visibility (`referencing_flows`) — one PR.
3. #45 phase 1: container death detection + `failed_block` column + auto-restart w/
   backoff (+ manual `POST /flows/{name}/restart`).

Open decisions to confirm before speccing: blue-green vs stop-then-start (recommendation:
defer blue-green), BlockFatal scope (recommendation: any container exit while running),
manual restart endpoint (recommendation: yes).
