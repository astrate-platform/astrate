# Hygiene report — 2026-09-06

Run against current HEAD on the working tree (not a specific sha — this was run
locally, not by the timer).

## Dependency sweep

`go list -m -u all` shows many transitive deps with newer versions (cloud/azure
ecosystem mostly). Direct deps were already swept 2026-09-04 and decided: no
bumps. Nothing new since then. The existing for-giulio entry covers this
comprehensively — no new research to add.

## govulncheck

Not installed in this environment. Cannot run.

## golangci-lint

0 issues. Clean.

## TODOs / FIXMEs / XXXs / HACKs

- `context.TODO()` in `internal/flow/blocks/httpblocks.go:147` and
  `internal/flow/blocks/astartesource/source.go:56` — standard Go idiom for
  contexts that will be wired later, not actionable.
- `internal/broker/intake.go:56` and `internal/store/store.go:139` —
  aspirational extension-point notes referencing ROADMAP/DESIGN docs. Ignored
  per recipe (notes-to-self, not missing behaviour).
- `XXX` in `internal/auth/claims_test.go:57-58` — placeholder device IDs in
  test data, not a TODO marker.

**Nothing proposable.**

## Packages with no test files

- `internal/housekeeping` — 262 lines of real business logic (realm provisioning,
  CA key minting/sealing, validation, deletion gating with connected-devices
  guard, default retention injection). No test file at all. **Proposed as a
  task.**
- `docs`, `examples/flow-container-echo`, `migrations` — non-logic packages.
  Ignored.

## Skipped tests

- `cmd/astrate/main_test.go:47` — skips without database (integration test).
  Expected.
- `cmd/astrate/forward_test.go:68` — skips without NATS build tag. Expected.
- `internal/pairing/http_test.go:629` — skips without openssl in PATH.
  Reasonable fallback.
- `internal/pairing/ca/ca_test.go:292` — skips 10k issuance draw in `-short`
  mode. Standard.
- `internal/store/datastreams_test.go:380` — skips when timescaledb_toolkit
  extension not installed. Reasonable for an optional feature path.

**Nothing proposable.** All skips are justified by missing external dependencies
or build tags.

## Summary

One task proposed: housekeeping unit tests (CreateRealm validation, default
retention injection, deletion guards). No for-giulio entries needed. No issues
filed.
