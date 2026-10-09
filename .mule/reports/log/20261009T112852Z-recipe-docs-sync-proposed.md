slug: recipe-docs-sync
verdict: proposed
at:  a086f26
ran: 2026-10-09T11:28:52Z on DietPi


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

Surface picked: **appengine**. Rotation evidence: rm 2026-10-02, housekeeping
2026-10-04, astrate-native 2026-10-05/06, pairing 2026-10-07/08; appengine is the
least-recently swept of the five. One surface only, five proposals maximum.

Route-set comparison (no proposal)
----------------------------------
`git grep`-style check of `internal/appengine/http.go:41` `Mount` (18 registered
paths) against `docs/api/astarte_appengine_api.yaml` (paths at yaml:26, 67, 141,
181, 259, 347, 381, 428, 625, 655, 697, 874, 952, 982, 1078, 1183, 1213, 1256)
found **no route gap**: every code route is documented and no documented path is
absent. The route-shape drift on this surface was already closed by
docs-sync-appengine-* / docs-sync-ae-* tasks, so this run is about fields and
error bodies.

Proposals appended to .mule/todo.md
-----------------------------------
Four lines (all `[auto]`, `docs/api/` only — no prose touched):

1. `docs-sync-ae-data-envelope-metadata` — `DataEnvelopeObject`
   (yaml:1711-1717) requires only `data` and types it `object`, but every
   `format=table` data GET returns a sibling `metadata` object
   (`astarteapi.WriteDataWithMetadata` / `metadataEnvelope`, envelope.go:149-158;
   `Tabular`, data.go:59-64; populated by `renderIndividual` data.go:256-268 and
   `renderObject` data.go:293-321). Add optional `metadata` (+ correct the `data`
   shape claim) on the six GETs that `$ref` it (yaml:416, 468, 685, 730, 1244,
   1290).

2. `docs-sync-ae-data-write-object-errors` — the six data PUT/POST ops document
   400 `BadRequest` and 422 `ValueTooLarge`, but two bodies are missing: 400
   `{"errors":{"detail":"Unexpected object key","unexpected_keys":[...]}}`
   (`writeBadObjectError`, http.go:722-748, const :711, ReasonBadObject case
   :684-691) and 422 `{"errors":{"detail":"Missing required mapping key"}}`
   (`detailMissingRequiredMapping` http.go:713-720, ReasonMissingRequired case
   :692-702), both implemented by appengine-unexpected-object-key /
   appengine-missing-required-422 with no spec update. `ErrorDetail`
   (yaml:1741-1750) also needs the optional `unexpected_keys` array.

3. `docs-sync-ae-group-create-devices-required` — `GroupCreate` (yaml:1652-1663)
   marks only `group_name` required, but `createGroup` requires `devices`
   (missing → 422 `can't be blank`, empty → 422 `should have at least 1
   item(s)`, http.go:391-418). Add `devices` to `required` and the two messages
   to the POST /groups 422 story (shared `ValidationErrors`, yaml:1841-1844,
   currently `group_name` only).

4. `docs-sync-ae-notfound-component` — the `NotFound` component (yaml:1795-1803,
   example `detail: Not Found`) is `$ref`d by all twelve group-route 404s
   (yaml:978, 1026, 1070, 1104, 1152, 1179, 1209, 1250, 1296, 1342, 1390, 1429),
   yet `astarteapi.WriteNotFound` / `DetailNotFound` ("Not Found") is emitted
   nowhere in internal/appengine. A missing group answers "Group not found"
   (`ErrGroupNotFound`, service.go:41-42, 508-516; writeError http.go:667-668)
   and a missing member answers "Device not found" (http.go:682-683).
   Repoint/rename the component accordingly.

Escalated to .mule/for-giulio.md (docs/site is never touched)
-------------------------------------------------------------
One line: `docs/site/appengine-api.md:60` writes `&downsample_to=<bucket_duration>`
and line 65 says it "maps onto Timescale `time_bucket()`", but the code reads it
as a **point count** (`opts.DownsamplePoints`, must be `> 2`,
internal/appengine/http.go:588-598; the bucket is computed from the count at
data.go:187-196), matching the OpenAPI description at
`docs/api/astarte_appengine_api.yaml:1509-1514`. Proposed wording supplied;
page untouched.

Deliberately not proposed (already queued)
------------------------------------------
- `docs-sync-ae-list-devices-422` (todo.md:240) — blocked; `GET /devices`
  missing 422 for invalid `from_token`. Not re-proposed.
- `docs-sync-appengine-data-set-422` (todo.md:144) — blocked; the value-size
  422 on the six data writes is covered by the done `docs-sync-ae-write-value-422`.
- `docs-sync-ae-read-query-params` (todo.md:208) — blocked; `retrieve_metadata`
  / `downsample_key` undocumented.
- The general `payload.ReasonOf` 500 family (todo.md:266) — left to its own
  blocked line; proposal 2 covers only the two reasons already mapped in code.

Nothing implemented. Nothing under `docs/` changed, so `make -C docs build` was
not run (the proposals name it for their executor). Git untouched: no commit,
branch, checkout, stash; no GitHub mutation.
