# Design A — Named multi-instance + durable flows (#40 + #41)

**Status:** **implemented** (A1–A5 landed 2026-07-29; open points §12 accepted as defaults).  
**Date:** 2026-07-29  
**Issues:** #40 (named multi-instance + config), #41 (durable + `auto_restart` rehydrate)  
**Decisions:** `docs/handoff/flow-v2-decisions-2026-07-29.md`  
**Follow-ups:** #42 (edge cases after MVP)

This is a short design note: enough for an implementer to land schema, API, and boot
rehydrate without reopening product choices.

---

## 1. Problem (today)

| Today | Wanted |
|---|---|
| At most one running instance per `realm/pipeline` | Many **named** flows per realm; many can share one pipeline |
| Manager key = `realm + "/" + pipeline` | Manager key = `realm + "/" + flow name` |
| `POST /flows` body: `{ "pipeline" }` only | `{ "name", "pipeline", "config"?, "auto_restart"? }` |
| Flow state lives only in `flow.Manager` (RAM) | Durable row in Postgres; survives process restart |
| No config bag / placeholders | Optional config; string placeholders in block config JSON |
| No boot rehydrate | On process start, start every row with `auto_restart=true` |

Pipelines stay as they are: realm-scoped recipes (`pipelines` table, `000008`).

---

## 2. Concepts

```
Pipeline  — reusable DAG blueprint (stored). Name unique per realm.
Flow      — named running (or desired) instance of a pipeline + config snapshot.
```

- **Identity:** `(realm, flow_name)` is unique. Not pipeline name.
- **Pipeline ref:** flow stores the **pipeline name** (text), resolved at each start.
  Do not FK to `pipelines.id` as a hard ON DELETE CASCADE for MVP — deleting a
  pipeline while flows reference it is a #42 concern; first cut: start fails loudly
  if the pipeline is missing.
- **Config snapshot:** the `config` JSON stored on the flow row is what start uses.
  Editing the pipeline definition does **not** mutate a flow’s config. Hot-reload of
  a running graph is out of first cut (#42).
- **`auto_restart`:** boolean, **default `true`** on create. When false, the row may
  still exist after stop (see §6), but boot never starts it.

---

## 3. Schema — `migrations/000009_flows`

Align with existing style (`bigint` identity PKs, `realm_id smallint`, `jsonb`), not the
old UUID sketch in `.mule/tasks/issue-25.md` (that sketch predated real `pipelines` and
named multi-instance).

```sql
-- 000009: durable Flow instances (issues #40 + #41).
-- A flow is a named specialization of a stored pipeline for a realm.

CREATE TABLE flows (
    id              bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    realm_id        smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    name            text NOT NULL,
    pipeline_name   text NOT NULL,
    config          jsonb NOT NULL DEFAULT '{}'::jsonb,
    auto_restart    boolean NOT NULL DEFAULT true,
    -- desired/runtime observation (single process; last write wins)
    status          text NOT NULL DEFAULT 'stopped',
    error_message   text,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    started_at      timestamptz,
    stopped_at      timestamptz,
    UNIQUE (realm_id, name)
);

CREATE INDEX flows_realm_auto_restart_idx
    ON flows (realm_id)
    WHERE auto_restart = true;
```

### Status values (text, match `flow.FlowStatus.String()` where possible)

| Value | Meaning |
|---|---|
| `creating` | Start in progress (optional; may be transient and rarely observed) |
| `running` | In-memory graph is up |
| `stopped` | Gracefully stopped; row retained if we keep durable history |
| `failed` | Start/rehydrate failed; `error_message` set; **loud** |

Boot and API always update `status` / `error_message` / timestamps so operators can
`GET` a failed flow after a bad rehydrate without relying only on logs.

**Note:** Historical issue-25 used `pipeline_id UUID` FK. **Do not** use that: pipelines
use `bigint` identity and are addressed by `(realm_id, name)`. Store `pipeline_name text`.

---

## 4. Store API (`internal/store/flows.go`)

Mirror `pipelines.go` patterns (`ErrNotFound`, `ErrAlreadyExists`, no import of
`internal/flow`).

Suggested methods:

| Method | Role |
|---|---|
| `CreateFlow(ctx, realmID, name, pipelineName, config, autoRestart)` | Insert; default autoRestart=true if caller passes nil-like default at service layer |
| `GetFlow(ctx, realmID, name)` | One row |
| `ListFlows(ctx, realmID)` | All durable rows for realm (not only running) |
| `ListFlowsToRehydrate(ctx)` or `ListAutoRestartFlows(ctx)` | Rows with `auto_restart=true` (all realms) for boot |
| `UpdateFlowRuntime(ctx, realmID, name, status, errMsg, startedAt, stoppedAt)` | Persist status after start/stop/fail |
| `UpdateFlowDesired(ctx, …)` | Optional later: patch config / auto_restart without start — **not required for MVP** |
| `DeleteFlow(ctx, realmID, name)` | Remove durable row (after stop if running) |

Config stored as raw `[]byte` / `json.RawMessage` JSON object (default `{}`).

---

## 5. Runtime manager (`internal/flow`)

### 5.1 Key change

Today:

```go
// FlowPipelineID = realm + "/" + pipelineName
flows map[string]*Flow  // keyed by PipelineID
```

After Design A:

```go
// FlowInstanceID = realm + "/" + flowName
// Manager map key = FlowInstanceID (not pipeline name).
```

Rename for clarity (recommended during implement):

- `FlowConfig.PipelineID` → keep field but **document** it as the **instance key**
  (`realm/flowName`), **or** add `InstanceID` and deprecate `PipelineID` in the same PR
  if churn is small.
- Helper: `flow.FlowInstanceID(realm, flowName string) string` (replace call sites of
  `FlowPipelineID` for lifecycle; pipeline name remains a separate field on the durable
  row / view).

`ErrFlowExists` message should say “flow already running” (instance name), not “pipeline”.

### 5.2 In-memory vs durable

- `Manager` remains the source of truth for **live** graphs (channels, pumps).
- Durable row is the source of truth for **desired** config + whether to rehydrate.
- `ListFlows` HTTP should return **durable** rows merged with live status when present
  (so a `failed` rehydrate is visible even with nothing in the Manager map).

Recommendation for failed start:

1. Do **not** leave a half-running entry in Manager (or leave only if status is
   `failed` and Stop is a no-op — prefer **no** map entry on fail).
2. Always write `status=failed` + `error_message` to DB.

Current `StartFlow` already sets `FlowStatusFailed` and keeps the map entry on graph
build failure. For durable MVP, prefer: on failure after durable create, **remove** the
in-memory entry (or never insert until success) so re-`POST` / rehydrate can retry, while
DB keeps `failed`. Exact choice: **never insert into Manager until graph+router succeed**;
on failure only update DB. Slightly cleaner than today’s keep-failed-in-map behaviour;
acceptable to change in the same PR (update unit tests).

---

## 6. API (`internal/flowapi`)

Routes stay Astrate-native:

| Method | Path | Behaviour |
|---|---|---|
| `GET` | `/flow/v1/{realm}/flows` | List durable flows (+ live status when running) |
| `POST` | `/flow/v1/{realm}/flows` | Create durable row + start (same path as rehydrate) |
| `GET` | `/flow/v1/{realm}/flows/{name}` | Get by **flow name** |
| `DELETE` | `/flow/v1/{realm}/flows/{name}` | Stop if running, then **delete durable row** |

### 6.1 Create body

```json
{
  "name": "prod-webhooks",
  "pipeline": "device-to-http",
  "config": { "webhook_url": "https://example.com/hook" },
  "auto_restart": true
}
```

| Field | Required | Notes |
|---|---|---|
| `name` | yes | Unique per realm; path segment for GET/DELETE |
| `pipeline` | yes | Must exist in `pipelines` for realm at start time |
| `config` | no | Object; default `{}` |
| `auto_restart` | no | **Default `true`** if omitted |

Breaking change vs today: body currently accepts only `{ "pipeline" }` and treats path
name as pipeline name. **No production clients documented** — change is OK; note in
CHANGELOG / handoff. Old behaviour (one instance named like the pipeline) can be
emulated by operators with `"name": "<same as pipeline>"`.

### 6.2 Response `FlowView`

```json
{
  "name": "prod-webhooks",
  "pipeline": "device-to-http",
  "realm": "acme",
  "config": { "webhook_url": "https://example.com/hook" },
  "auto_restart": true,
  "status": "running",
  "error_message": null,
  "created_at": "...",
  "updated_at": "...",
  "started_at": "...",
  "stopped_at": null
}
```

- Drop or demote synthetic runtime id (`flow-N`) from the primary operator view; optional
  `runtime_id` field if useful for logs. Operator identity is **name**.
- List returns objects (already Astrate-native; not upstream’s name-only array).

### 6.3 Errors (loud)

| Situation | HTTP | Notes |
|---|---|---|
| Unknown pipeline | 404 or 422 | Prefer 422 with clear detail if pipeline missing at start |
| Duplicate flow name | 409 | |
| Already running | 409 | |
| Placeholder unresolved / instantiate fail | 422 | Body detail + DB `failed` if row was created |
| Rehydrate fail at boot | n/a | Log error + DB `failed`; **continue** other flows |

### 6.4 DELETE semantics (MVP)

1. If in Manager and running → `StopFlow` (drain + Stopper).
2. Delete durable row.
3. 204.

Optional later (#42): “stop but keep row with `auto_restart=false`” via a dedicated
endpoint or query flag. Not required for MVP; product only needs never-restart via
create-time flag + delete to fully remove.

If we need stop-without-delete for ops, add later:

- `DELETE` with `?keep=1` **or** `POST .../flows/{name}/stop` — track under #42.

---

## 7. Config substitution (#40)

**First cut:** string placeholder rewrite over the pipeline definition **after** load,
**before** `ParseDefinition` / `Instantiate`.

### 7.1 Syntax

In any **string** value inside block `config` (and only strings for MVP):

```
${config.<key>}
```

- `<key>` is a single path segment matching `^[A-Za-z_][A-Za-z0-9_]*$` (no nested paths
  in v1 unless easy: optional `${config.a.b}` as nested object lookup).
- **Recommended v1:** flat keys only: `${config.webhook_url}` → `config["webhook_url"]`
  must be a string (or stringable number/bool); missing key → **loud fail** at start.
- Non-string JSON values (numbers, bools, objects) in the pipeline definition are left
  unchanged.
- No full templating language, no Elixir DSL.

### 7.2 Algorithm

1. Load stored pipeline definition bytes.
2. Load flow config map.
3. Walk JSON; for each string, replace all `${config.KEY}` occurrences.
4. If any placeholder remains unmatched or KEY missing → error (do not start).
5. Parse substituted JSON as pipeline → `Instantiate` → `Manager.StartFlow`.

Implementation sketch: recursive walk on `any` after `json.Unmarshal`, or regexp on
raw string with care for escapes — prefer structured walk so only string leaves are
touched.

### 7.3 Out of first cut

- JSON Schema validation of `config` against pipeline-declared schema.
- Placeholder substitution inside non-config pipeline fields (block names, types).
- Hot update of config on a running flow.

---

## 8. Single start path + boot rehydrate

### 8.1 Internal start (one function)

All starts go through one service method, e.g.:

```text
flowapi.Service.startFlowInstance(ctx, realm, name, pipelineName, config, autoRestart, persist bool)
```

Used by:

1. `POST /flows` — `persist=true` (insert or conflict).
2. Boot rehydrate — row already exists; `persist` only updates runtime columns.

Steps (order matters):

1. Resolve realm → `realm_id`.
2. Load pipeline by name; missing → fail (write `failed` if row exists).
3. Apply config substitution.
4. `ParseDefinition` + `checkBlockTypes` + `Instantiate` with `Deps{Bus, Realm}`.
5. `Manager.StartFlow` with instance key `realm/name`.
6. On success: DB status `running`, clear `error_message`, set `started_at`.
7. On failure: DB status `failed`, set `error_message`, log at Error; return error to API.

### 8.2 Boot (process)

In `cmd/astrate/main.go` (or a small `flowapi`/`flow` helper called from main), **after**:

- store open + migrations,
- engine started (bus live),
- flow Service + Manager constructed,

and **before** or **just after** HTTP listen (prefer **before** accept traffic so
rehydrate races less with operator calls):

```text
for each row in ListAutoRestartFlows():
    err := startFlowInstance(...)  // same path
    if err != nil:
        log.Error("flow rehydrate failed", realm, name, err)
        // row already marked failed inside start path
        continue  // do not abort process boot
```

Single Astrate process assumed. No distributed lock for MVP.

### 8.3 Shutdown

Existing `flowMgr.Shutdown` drains running flows. After drain, optionally set durable
status to `stopped` for rows that were running (best-effort). On next boot,
`auto_restart=true` starts them again. If process crashes without clean shutdown,
rows may still say `running` — **rehydrate still starts them** (treat `auto_restart` as
desired, not “status must be stopped”). On start success, overwrite status to `running`.

---

## 9. DELETE pipeline while flows exist

**MVP:** `DeletePipeline` does **not** cascade-stop flows (today already “running flows
are not stopped”). Starts fail with loud “pipeline not found” / failed status.

**#42:** optional block delete if flows reference pipeline, or cascade stop+delete.

---

## 10. Implementation plan (after acceptance)

Small, test-backed PRs preferred (can be one PR if tight):

| Step | Work | Tests |
|---|---|---|
| **A1** | Migration `000009` + `store` CRUD + list auto-restart | store tests (pg or sqlite pattern used elsewhere) |
| **A2** | Config substitution helper in `internal/flow` (or flowapi) | unit tests: replace, missing key, no-op |
| **A3** | Manager instance key = flow name; Start only after success | existing manager tests updated |
| **A4** | `flowapi` create/list/get/delete + FlowView fields | service + http tests |
| **A5** | Boot rehydrate hook in `main` | integration test if cheap; else unit test of rehydrate loop with fake store |

**Do not** implement container blocks (#43) in these PRs.

### Acceptance criteria (MVP)

1. Two flows with different names can run the same pipeline with different configs.
2. Process restart restarts only `auto_restart=true` flows; `false` stays down.
3. Missing pipeline / bad placeholder / unknown block → status `failed` + error text +
   non-2xx on API; boot continues.
4. `DELETE /flows/{name}` stops live work and removes the durable row.
5. Existing pipeline CRUD and built-in filter/map path still work.
6. `go test ./internal/flow/... ./internal/flowapi/... ./internal/store/...` green.

---

## 11. Explicitly out of this design (→ #42 or later)

- Hot-reload when pipeline JSON changes under a running flow.
- Partial block restart / rolling graph update.
- Multi-process or HA flow managers.
- Stop-without-delete API polish.
- JSON Schema for config.
- Upstream wire clone (`a_f` claim, name-only list, DSL `source` string).
- Containers (#43 / Design B).

---

## 12. Open points for Giulio (accept or tweak)

Only these need a nod if the rest is fine; defaults in brackets are recommended:

1. **DELETE = stop + remove row** for MVP? **[yes]**
2. **Placeholder syntax** `${config.key}` flat only? **[yes]**
3. **Failed start:** no Manager entry, DB `failed`? **[yes]**
4. **Boot timing:** rehydrate before HTTP listen? **[yes]**
5. **ListFlows:** durable rows including stopped/failed, not only live? **[yes]**

If accepted as-is, next session: implement A1–A5 (or mule only after explicit accept).

---

## 13. Related files (current code anchors)

| Area | Path |
|---|---|
| Manager | `internal/flow/flow.go` |
| Instance key helper | `internal/flow/factory.go` (`FlowPipelineID`) |
| Service start/stop | `internal/flowapi/service.go` |
| HTTP body | `internal/flowapi/http.go` (`startFlowBody`) |
| Pipelines table | `migrations/000008_pipelines.up.sql`, `internal/store/pipelines.go` |
| Process wiring | `cmd/astrate/main.go` |
| Product decisions | `docs/handoff/flow-v2-decisions-2026-07-29.md` |
