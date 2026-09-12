slug: recipe-docs-sync
verdict: proposed
at:  1a2ed87
ran: 2026-09-12T19:22:46Z on this host in ~60s

Surface picked this run: **appengine** (`docs/api/astarte_appengine_api.yaml` vs `internal/appengine/http.go`).

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

## Route surface

Documented paths (17) vs routes registered in `API.Mount` (internal/appengine/http.go:40-81):
1:1, same methods on every path, including the websocket route `GET /appengine/v1/socket/websocket`
(channels/ws.go:48), which is documented in `astrate_native_api.yaml:325` and docs/site/appengine-api.md.
No missing or stray paths.

## Status-code / schema spot-checks (more than the recipe's 3-4, kept read-only)

- GET /devices, GET /stats/devices, GET /devices/{device}, PATCH /devices/{device}: statuses match
  handlers and `writeError` mapping (http.go:643-686). PATCH 500 is documented and legitimately covers
  the merge-patch content-type mismatch.
- Data path GET / PUT / POST / DELETE (device-scoped, by-alias, in-group): 200/204, 400, 404, 405,
  422 (ValueTooLarge + ValidationErrors) all present — the earlier done lines
  `docs-sync-ae-write-405` and `docs-sync-ae-write-value-422` are in place.
- POST /groups 409, PATCH device/alias 422+409, PATCH group-device 422+409, interface-level GET 422,
  downsample "(must be > 2)" — all present (earlier done lines verified not re-proposed).
- DeviceStatus schema: `id` (not `device_id`), `introspection` as `{major, minor}` map, and the
  six added fields all present (yaml:1442-1499) — the `docs-sync-appengine-device-status-schema` line is in place.

## New drift found

1. **`POST /appengine/v1/{realm}/groups/{group}/devices` is missing `409`.** Re-adding a member the
   group already holds: `Service.AddGroupDevice` returns `ErrDeviceAlreadyInGroup` (service.go:636-638),
   `writeError` maps it to `409 "Device already in group"` (http.go:662-663). The spec (yaml:967-977)
   documents only 201/400/401/404/500. Proposed.

2. **All three PATCH request bodies declare `application/json; charset=utf-8`** (device yaml:176,
   by-alias yaml:260, group-device yaml:1023) but every PATCH handler enforces
   `Content-Type: application/merge-patch+json` exactly and rejects everything else — 500 on the
   device and group paths (http.go:209-211, 521-523), 400 on the alias path (http.go:230-232). A
   spec-following client sends `application/json` and is rejected. The declared media type is wrong
   for the enforced wire on all three. Proposed as one task.

## Elsewhere

- ASTRATE_ config-key sweep: the two `ASTRATE_HOUSEKEEPING_*` keys absent from the config reference
  were already escalated to `.mule/for-giulio.md` today (2026-09-12, docs-sync appengine entry) — not
  re-raised.
- No `docs/site/` page contradicted by the code on this surface; nothing appended to for-giulio.
- Recipe does not ask for `gh issue create`; no issue filed.

```
Done: docs-sync appengine — 2 drift items proposed (missing 409 on POST /groups/{group}/devices; PATCH bodies must declare application/merge-patch+json)
Files: .mule/todo.md (2 task lines appended), .mule/reports/log/20260912T192246Z-recipe-docs-sync-proposed.md
Verified: registered-route list vs documented paths -> 17/17 matched; status-code spot-checks vs handlers/writeError -> all but the two proposed gaps correct
Unsure: whether the media-type fix should keep parity with upstream's own published spec (the merge-patch enforcement itself is deliberately upstream-verbatim, per http.go:205-208); the doc change is truthful to the enforced wire either way
Follow-ups: docs-sync-ae-post-group-devices-409, docs-sync-ae-patch-merge-patch-media-type
```