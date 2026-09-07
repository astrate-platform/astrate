# docs-sync-appengine-data-output-params — document the output/format query params on the AppEngine data GETs

Proposed by the docs-sync recipe run on 2026-09-07 (surface: appengine).

## Drift

All six data GET operations in `docs/api/astarte_appengine_api.yaml` document only
`since`, `since_after`, `to`, `limit`, `downsample_to`, `sort`
(parameter refs `DataSince`/`DataSinceAfter`/`DataTo`/`DataLimit`/`DataDownsample`/`DataSort`).
The code accepts five more query parameters via `parseQueryOpts` (`internal/appengine/http.go:554-641`),
none of which the spec mentions. The six operations and their handlers:

- `GET /appengine/v1/{realm}/devices/{device}/interfaces/{interface}` — `getData` (http.go:50)
- `GET /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}` — `getData` (http.go:51)
- `GET /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}` — `getDataByAlias` (http.go:61)
- `GET /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}` — `getDataByAlias` (http.go:62)
- `GET /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}` — `getDataInGroup` (http.go:76)
- `GET /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}` — `getDataInGroup` (http.go:77)

All three handlers go through `parseQueryOpts` (http.go:254).

## What the code actually does (verified 2026-09-07)

Three of the five missing params genuinely change the response and must be documented:

| param | parsed | consumed |
|-------|--------|----------|
| `format` (`structured`/`table`/`disjoint_tables`, default `structured`) | http.go:601-607 | data.go:237, 244 — `renderObject`/`renderIndividual` |
| `allow_bigintegers` (`true`/`false`) | http.go:612 | data.go:230, 403 — biginteger handling |
| `allow_safe_bigintegers` (`true`/`false`) | http.go:613 | data.go:230, 405 — biginteger handling |

Two are parsed but have **no consumer** anywhere in `internal/` (grep for `RetrieveMetadata`/
`DownsampleKey` finds only the struct fields in data.go:53-54 and the parse sites):

- `retrieve_metadata` (http.go:624-627 → `opts.RetrieveMetadata`) — dead, silently ignored.
- `downsample_key` (http.go:633 → `opts.DownsampleKey`) — dead, silently ignored.

Do not document these two as if they had an effect. Options for the executor, in order of
likelihood: (a) document only the three live params, (b) wire the two dead ones up if there is
an issue that wants them, (c) confirm upstream parity on the Legion Go before deciding. At
minimum the spec must stop claiming the GET is described by its current six-param-only list.

## Scope note

This does **not** overlap the blocked `appengine-snapshot-ignores-query-params` line, which is
about the interface-root snapshot branch ignoring the *window* params (`since`/`since_after`/`to`/
`limit`/`sort`). This task is about *output* params never documented at all. `format=table` also
changes the response `data` shape (a table envelope rather than an object), so consider
documenting the response-shape consequence, not just the parameter list.

## Verification

Run `make -C docs build` and confirm every YAML the Swagger UI references still loads
(the docs-sync recipe's build gate). No code changes; `docs/api/*.yaml` only.