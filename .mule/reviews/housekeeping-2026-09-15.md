# Review — internal/housekeeping

Date: 2026-09-15
Reviewed files: service.go (lines 1–280), http.go (lines 1–259), service_test.go (lines 1–224),
http_test.go (lines 1–511), async_operation_test.go (lines 1–37).

## Finding: PATCH treats `retention: 0` as "clear"

`http.go:216-223` — when `DatastreamMaximumStorageRetention` is present and non-null, the
handler checks `val == 0` and routes to `ClearRetention = true`:

```go
if pb.DatastreamMaximumStorageRetention.present {
    if pb.DatastreamMaximumStorageRetention.null || pb.DatastreamMaximumStorageRetention.val == 0 {
        u.ClearRetention = true
    } else {
        u.PatchRetention = true
        u.SetRetention = pb.DatastreamMaximumStorageRetention.val
    }
}
```

The `Create` path (`service.go:146-151`) does NOT treat 0 specially — a 0 in the create
body is accepted and persisted (0 ≥ 0 passes the negative check, `retention != nil` means
the default is not injected, and the store receives `*int64(0)`). So the two write paths
disagree on whether 0 is a valid retention value.

This is an upstream-parity deviation worth pinning: if 0 is invalid for retention (which
it semantically should be — "delete immediately" is not a useful datastream policy), the
`Create` path should reject it; if 0 is valid, the PATCH path should not silently clear it.

**Proposed task:** add a container-free unit test pinning the chosen behaviour. If the
decision is "reject 0" (recommended — matches upstream's intent that retention > 0),
both the create and PATCH paths gain a `*retention == 0` / `val == 0` → `ErrValidation`
branch; if "accept 0", the PATCH path drops the `|| val == 0` clause.

## Not proposing

- **`Service.UpdateRealm` missing negative-value guard (service.go:196-216):** the HTTP
  layer catches negative retention/limit values at http.go:192-200 before they reach the
  service, so the store never sees them today. Adding a defensive check at the service
  level is good hygiene but not a behaviour change (the existing tests already pass); it
  stays in the review only.

- **Double-parse of PATCH body (http.go:164-185):** the body is read into `fields
  map[string]json.RawMessage` to check for unknown keys, then re-read into `patchBody`.
  Technically redundant but preserves the clean separation between the unknown-field
  check and the typed decode; not worth the refactor.

- **`notifyBrokerReload` swallows errors (service.go:119-126):** by design — the realm
  row is already committed, and the doc comment says self-heal recovers. This is a
  documented trade-off, not a bug.
