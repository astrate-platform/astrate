# Review — pkg/astarteapi

Date: 2026-09-15
Reviewed files: envelope.go (lines 1–249), envelope_test.go (lines 1–236), httpx/notfound.go
(lines 1–56) as the consumer of the two 404 detail constants, appengine/http.go:88-127 &
440-479 as the consumers of `WriteDataWithLinks`.

This is the wire-frozen envelope package: every constructor's exact bytes are pinned by
golden fixtures, so most of the contract is already protected. I picked it because it is
the largest package with no entry in `.mule/reviews/` yet.

## Finding 1: `WriteDataWithMetadata` has no direct coverage

`envelope.go:154-158` — the constructor's wire rule (nil metadata omits the `metadata`
key, `envelope.go:151` `omitempty`) is asserted nowhere. The non-nil branch is only
exercised indirectly, and only under the `//go:build integration` suite
(`internal/appengine/formats_parity_test.go:103-146`), which runs on the Legion Go, not on
the container-free gate. The nil-omission branch is covered nowhere at all — production
never passes nil (`internal/appengine/data.go:262-268` always sets it), so a regression
dropping or relocating the `omitempty` (e.g. serializing `"metadata": null`) would pass
every suite on this box.

**Proposed task:** add a container-free golden test for `WriteDataWithMetadata` asserting
(1) `nil` metadata renders `{"data": v}` with no `metadata` key, (2) a populated metadata
map renders `{"data": v, "metadata": {...}}`, both with 200 and the package `Content-Type`,
matching the existing `TestGoldenEnvelopes` style.

## Finding 2: the "keys emitted in sorted order" claim is unasserted

`envelope.go:170` and `envelope.go:184` document that `WriteFieldErrors` and
`WriteRawErrors` emit map keys in sorted order, but every golden fixture uses a
single-key map (`envelope_test.go:101-121`). Ordering today is guaranteed only by
`encoding/json`'s map handling; a change to the marshalling path (or a switch of the
encoder) that stopped sorting would pass the suite. The comment at `envelope_test.go:112`
("keys must come out sorted") states the rule but does not pin it.

**Proposed task:** extend `TestGoldenEnvelopes` with multi-key `WriteFieldErrors` and
`WriteRawErrors` maps asserting the sorted key order on the wire.

## Not proposing

- **`DecodeData` size accounting counts envelope overhead (`envelope.go:228-235`):** the
  cap is applied to the whole `{"data": ...}` bytes, not the payload. Documented ("reads
  at most maxBytes bytes from r") and unlikely to diverge from upstream's plug-level body
  limit; changing it would be invented behaviour without an upstream probe.

- **Duplicate `data` keys take the last occurrence (Go json semantics):** matches the
  documented sibling-key tolerance in spirit; no evidence of an upstream difference.

- **`Links.Self` without `omitempty` (`envelope.go:134-136`):** both consumers
  (`deviceListLinks`, `groupListLinks`) always set Self; no live caller can emit
  `"self":""`. A defensive `omitempty` would be speculative.

- **A `Content-Type: application/json` check in `DecodeData`:** invented behaviour,
  upstream-compat unverified; skip.

- **R5 marshal-failure fallback (`envelope.go:106-115`)** is already covered by
  `TestWriteDataMarshalFailure`; a write-failure-on-the-wire test is unbuildable with
  `httptest.ResponseRecorder` and low value.

Verdict: a small, well-frozen package; the two gaps above are both missing-test items, no
behavioural bug found.