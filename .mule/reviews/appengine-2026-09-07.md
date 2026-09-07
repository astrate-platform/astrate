# Code review: internal/appengine — 2026-09-07

Area chosen because recent reviews were `broker` (2026-09-04), `flow`
(2026-09-05), `engine` (2026-09-06); `appengine` is the largest package with no
review file in `.mule/reviews/` (15 files, ~7200 lines). I read the four
non-test source files largest-first plus the unit-tested seams.

## What I read

- `internal/appengine/service.go` (846) — device list/status/patch, groups,
  by-alias/by-group mirrors, applyPatch validation order, group token
  pagination
- `internal/appengine/http.go` (696) — route table, parseQueryOpts, writeError
  taxonomy, links builders
- `internal/appengine/data.go` (514) — resolve/GetData, datastreamData
  (snapshot / series / downsample / object), bigint reshaping, properties
  snapshot
- `internal/appengine/downsample.go` (29) — bucketFor
- tests: `query_opts_test.go`, `downsample_test.go`, `links_test.go`,
  `groups_parity_test.go`, `http_test.go` (grep)

## What I found (worth proposing)

### 1. The interface-root snapshot silently drops documented query params

`internal/appengine/data.go:138-148`: the individual-datastream snapshot branch
(`path == ""`, no downsample) calls `IndividualSnapshot` and ignores
`Since`/`SinceAfter`/`To`/`Limit`/`Descending` entirely. The endpoint is
documented as taking all of them (docs/api/astarte_appengine_api.yaml:577-582,
with no path-dependency stated). This is the exact "silent lie" the same
function refuses to commit elsewhere: downsample-on-properties is a 422
(data.go:118-120), downsample-at-root needs a path (data.go:165-167), an empty
object root renders `{}` instead of dropping (data.go:214-218). The snapshot
branch drops the params with no comment, unlike the downsample branches.

It is plausible upstream also just returns the latest-per-endpoint snapshot here
(so window/limit genuinely don't apply) — but then Astrate should refuse them
explicitly (422, per its own precedent) rather than accept-and-ignore, and the
docs should not advertise them on the root GET. Needs a live upstream probe to
decide honour-vs-refuse; propose the verification-plus-fix task, not a silent
change.

### 2. `parseGroupToken` / `groupTokenFor` have no direct unit test

The cursor round-trip is pure (service.go:602-623) but only exercised through
the pg-backed `groups_parity_test.go` (integration). The interesting rules —
round-trip through `groupTokenFor(offset)` → `parseGroupToken`, and that a
non-v1-nibble UUID is rejected while a well-formed v1 with zeroed fields is
accepted (parseGroupToken service.go:617-623) — are machine-checkable without a
store and currently rot unasserted. Category-2.

## Decided NOT to propose, and why (so the next run does not re-derive them)

- **Snapshot refusing vs honouring is a design call I cannot make blind.** I
  proposed it as a task (above) because the silent drop is a real asymmetry
  with the function's own precedent; but whether the right fix is refuse-422 or
  honour-the-window must be confirmed against upstream on the Legion first. If
  a later reviewer decides it is genuinely intended, it belongs in
  for-giulio.md instead — do not re-open it as a bug without the probe.
- **`applyBigIntegerOpts` asymmetry** (data.go:386-412): with
  `allow_safe_bigintegers=false` (or `allow_bigintegers` nil + safe false) the
  longinteger leaf stays a json.Number. Whether upstream degrades to a string
  in that combination is a semantics call needing a live probe; the current
  behaviour is documented and the parity test pins something. Flagging so the
  next review does not re-derive it, not proposing it.
- **`creatorGroup`'s 201 echo of the original body** (http.go:422-423) is
  deliberate (duplicates included) and matches upstream. Not a bug.
- **`deviceListLinks` self echoing valid-but-minor query params** and the
  from_token-first-encode dance are intentional and have pins in
  `links_test.go`. Fine.
- **`renderObject` indexing `docs[0]`** (data.go:294, 323) is safe: the empty
  object-series case returns before the call (data.go:214-219).
- **`parseGroupToken` returning `int` from a uint32** (service.go:622) could
  overflow on a 32-bit runtime, but the platform is 64-bit and offsets stay
  well under 2^31. Not worth a task.
- **`PublisherData`/`UnsetProperty` nil-`sd` guard** duplicates the disabled
  message on every call; harmless.

The package is exceptionally well-documented and the data/http/service logic is
heavily tested (query-opts, links, parity, downsample). Finding only the two
items above in ~2000 lines is the expected result for code at this level.
