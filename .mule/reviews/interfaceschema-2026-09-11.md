# Code review: pkg/interfaceschema — 2026-09-11

Area chosen because `appengine` (09-07), `store` (09-08), `pairing`/`realm`
(09-09), `flowapi` (09-10) have review files, `broker`/`flow`/`engine` were
09-04..09-06, and `interfaceschema` is the largest package (~3200 lines) with
no review file. The last review ran 09-10, so this is the next rotation slot.
I read every non-test file largest first — `parse.go` (664), `types.go` (406),
`compat.go` (131), `compile.go` (122), `trie.go` (161), `violations.go` (82) —
plus all seven test files and the invalid-fixture manifest, and I ran
`go test ./pkg/interfaceschema/ -coverprofile` to ground the coverage claims
(87.3% statements overall).

## What I found (worth proposing)

The package is in very good shape: strict decode, validated defaults, every
rule named and file/line-anchored, probes against upstream 1.2.0 documented
(#61/#62/#67). No behaviour struck me as wrong or unguarded; the trie is
allocation-free and bench'd. The gaps are tests on four rules that already
exist but are not fully asserted — a regression in each would pass the suite
today.

1. **`sameMappingAttributes` (compat.go:102) leaves four immutable attributes
   unasserted.** `TestCheckMinorUpgradeTable` flips type, reliability,
   explicit_timestamp, required, encrypted — but never retention, expiry,
   database_retention_policy/ttl, or allow_unset. A regression that relaxed
   expiry or DB-TTL immutability on minor upgrade would pass. Allow-unset can
   only be exercised with a properties old/next pair (a datastream insert
   rejects it at parse time).

2. **`sameObjectAttributes` (parse.go:626) pins only one of six uniformity
   rules.** Object-aggregation rejects heterogeneous datastream attributes
   ("must all share the same <attr>"), but the only fixture is reliability
   (testdata/invalid/17-object-mixed-attrs.json). Retention, expiry,
   explicit_timestamp, database_retention_policy/ttl have no rejection twin.
   Coverage: 54.5% statements.

3. **`parseSegment`'s four malformed-placeholder branches (parse.go:509-517)
   are untested.** The only placeholder-syntax fixture is
   12-endpoint-partial-placeholder.json, which exercises the mid-segment `%`
   case (`strings.ContainsAny` line 520), not the param branch: segment
   starting with `%` but len < 4, second char not `{`, missing closing `}`, or
   an invalid placeholder name. `/ "%{"` or `/%{1}` would be accepted by a
   regression that passed today's suite.

4. **`TestEnumRoundTrips` (parse_test.go:244) omits four enums.** It round-trips
   InterfaceType, Reliability, ValueType only. Ownership, Aggregation,
   Retention, DatabaseRetentionPolicy wire strings are never marshalled in the
   package's own tests (Retention.String/MarshalJSON at 0% coverage,
   DatabaseRetentionPolicy at 0%), so a typo'd wire name would go uncaught.

5. **The `matchNode` (trie.go:124) "literal interior does not shadow its
   parametric sibling" corner is behaviourally untested.** With a legal
   interface declaring `/a/b/c` and `/a/%{p}`, path `/a/b` resolves to the
   parametric mapping (the exact literal child at that depth is interior-only,
   mapping nil, so matching falls through to `n.param`). Statement coverage
   cannot see this branch — `matchNode` reports 100% — so a "fix" that missed
   on any shallow path through an interior literal node would sail through the
   suite.

## What I decided NOT to propose, and why

- **`Compile`'s defensive object-aggregation branches** (duplicate object leaf,
  param last level, compile.go:99-101 / 117-121). Only reachable through a
  hand-built `Interface` missing parse's guarantees; `Compile(nil)` and the
  resolver-error path are already tested. Out of contract.
- **`interfaceNameRe` shape** (parse.go:43). I reasoned that middle labels may
  begin with a digit under this regex, but I did not probe upstream, and the
  recipe forbids proposing from pattern-matching alone. Would need a live probe
  first; not worth a task without measurement.
- **`required`/`encrypted` accepted on individual datastream mappings.**
  Deliberate, documented as 1.4-experimental semantics per-key (issue #67);
  `TestObjectAggregationPerKeyRequired` pins the per-key intent. No change.
- **Performance.** The only hot path is `Match`, and it is contractually
  allocation-free (`TestMatchZeroAllocs`) with benchmarks present. Nothing to
  measure.
- **docs/DESIGN.md contradiction:** none found; code and spec agree.

Five missing tests in a package this well-covered, and no behavioural defect,
is the correct result — reading closely was what turned those five up, and
each names the exact area the next reviewer should not have to re-examine.