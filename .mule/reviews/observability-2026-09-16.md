# Review — internal/observability

Date: 2026-09-16
Reviewed files: health.go (1–83), metrics.go (1–80), compat.go (1–52),
observability_test.go (1–128), compat_version_test.go (1–55), and the wiring in
cmd/astrate/main.go:111, 384-413 (NewMetrics, MountServiceCompat,
MountVersionCompat, RegisterBrokerSessions, RegisterDBPool, NewHealth,
AddReadiness ×2, Mount).

Small well-written package: the health/readiness/metrics surface under
/astrate/v1 and the upstream-parity /{service}/health and /{service}/version
endpoints. Most wire contract is already pinned by tests that mirror upstream's
measured envelopes (`{"data":{"status":"ok"}}`, `{"data":"<version>"}`). No
behavioural bug found; three missing-test items below.

## Finding 1: the readiness budget rule has no test

`health.go:12` (`readinessTimeout = 3 * time.Second`) and its comment —
"bounds the whole readiness probe so a wedged dependency can't hang the
endpoint (and the orchestrator that polls it)" — exist, but nothing asserts
them. Both consumers (`handleReadiness`, health.go:56, and `MountServiceCompat`,
compat.go:22) pass the const into `context.WithTimeout`, and no test exercises
the wedged-dependency path at all: `TestReadiness` (observability_test.go:37-63)
runs only checks that return immediately. A regression that dropped the timeout
altogether (hanging every orchestrator poll) would pass the whole suite. The
budget is a const, so a behavioural test would cost 3s of wall time — the
mechanical fix is to make it an injectable `Health` field defaulting to the
const, then assert a ctx-honoring wedged check returns 503 inside the budget.

**Proposed task:** add a `timeout time.Duration` field on `Health` (defaulting
to `readinessTimeout`) and a test proving a ctx-honoring wedged check cannot
hang `/astrate/v1/readiness` past the budget — the handler returns 503 with the
wedged check reported failing, within a short injected budget.

## Finding 2: two of the four db_pool gauges are unasserted

`TestMetricsExposesGauges` (observability_test.go:25-34) asserts
`astrate_db_pool_acquired_conns` and `astrate_db_pool_max_conns`, plus the
broker-sessions gauge and `go_goroutines` — but never `astrate_db_pool_idle_conns`
or `astrate_db_pool_total_conns` (metrics.go:77-78). A regression deleting (or
mis-wiring the pick function of) either gauge would pass the suite.

**Proposed task:** extend `TestMetricsExposesGauges` to assert all four
`astrate_db_pool_*` gauges, matching the values in the supplied
`DBPoolStats{2, 5, 7, 10}` (idle 5, total 7).

## Finding 3: no Content-Type assertion on the /astrate/v1 surface

`writeJSON` (health.go:79-82) sets `Content-Type: application/json`, and the
compat handlers set it by hand (compat.go:26, 31, 44), but only the version test
pins the header (compat_version_test.go:31-33). `/astrate/v1/health`,
`/astrate/v1/readiness`, and the per-service `/health` (both the 200 and the 503
path) are asserted only on status code and body.

**Proposed task:** assert `Content-Type: application/json` on the three
`/astrate/v1` responses in `TestReadiness`/`TestReadinessAllOK` and on both
paths of `TestMountServiceCompat`'s per-service handler.

## Not proposing

- **Per-service `/health` probes only the DB, not the broker** (`main.go:384`
  passes `st.Health` to all three services, while `/astrate/v1/readiness` probes
  DB + broker at main.go:411-412): a broker-only failure would leave the
  Dashboard's per-service indicator green. This is a deliberate-deviation design
  surface (deviation 18, docs/COMPATIBILITY.md) and the fix is a wiring decision
  in main, not a package bug — noted, not proposed.
- **`MustRegister` panic on a second registration** (`metrics.go:57`, `69`):
  startup-only, currently called exactly once each (main.go:400-401), and
  fail-loud-at-init is the Prometheus idiom. Speculative.
- **`AddReadiness` appends unsynchronized** (health.go:38): all four call sites
  are init-time, single-threaded; no live race.
- **`writeJSON` drops the encoder error** (health.go:82): cannot fail for
  `map[string]string`/`map[string]any` bodies of simple strings.
- **Raw dependency error text in the unauthenticated /readiness body**
  (health.go:64): expected errors carry no secrets; guarding would be invented
  behaviour.
- **4× `Store.Stat()` snapshot per scrape** (metrics.go:68-79, each gauge pick
  calls fn independently): the scrape is one Prometheus poll per 15s against a
  pgx pool snapshot; not a named hot workload.

Verdict: no behavioural bug; three cheap, fully machine-checkable missing-test
items.