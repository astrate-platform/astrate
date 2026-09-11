slug: interfaceschema-malformed-placeholder-fixtures
verdict: done
at:  f6a44a2
ran: 2026-09-11T20:04:16Z on DietPi in 130s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
✱ Glob "pkg/interfaceschema/**/*.json" 50 matches
→ Read pkg/interfaceschema/parse.go [offset=460, limit=120]
✱ Glob "pkg/interfaceschema/*_test.go" 6 matches
→ Read pkg/interfaceschema/testdata/invalid/manifest.json
→ Read pkg/interfaceschema/testdata/invalid/12-endpoint-partial-placeholder.json
✱ Glob "pkg/interfaceschema/*.go" 12 matches
✱ Grep "manifest" in pkg/interfaceschema · 16 matches
→ Read pkg/interfaceschema/parse_test.go [offset=175, limit=80]
→ Read pkg/interfaceschema/testdata/README.md
→ Read pkg/interfaceschema/testdata/invalid/11-endpoint-not-rooted.json
→ Read pkg/interfaceschema/testdata/invalid/13-endpoint-duplicate.json
→ Read pkg/interfaceschema/testdata/invalid/19-object-param-leaf.json

← Write pkg/interfaceschema/testdata/invalid/33-endpoint-malformed-placeholder-open.json
Wrote file successfully.

← Write pkg/interfaceschema/testdata/invalid/34-endpoint-malformed-placeholder-mid.json
Wrote file successfully.

← Write pkg/interfaceschema/testdata/invalid/35-endpoint-malformed-placeholder-unclosed.json
Wrote file successfully.

← Write pkg/interfaceschema/testdata/invalid/36-endpoint-invalid-placeholder-name.json
Wrote file successfully.

The task: add four invalid JSON fixtures + manifest entries covering `parseSegment`'s four placeholder-syntax rejections (`/%{`, `/%ax`, `/%{x` → "malformed placeholder"; `/%{1}` → "invalid placeholder name"). I'll create them and wire them into the lockstep-driven `TestInvalidFixturesRejected`.
← Edit pkg/interfaceschema/testdata/invalid/manifest.json
Index: /root/astrate-mule/pkg/interfaceschema/testdata/invalid/manifest.json
===================================================================
--- /root/astrate-mule/pkg/interfaceschema/testdata/invalid/manifest.json
+++ /root/astrate-mule/pkg/interfaceschema/testdata/invalid/manifest.json
@@ -29,6 +29,10 @@
   "28-object-mixed-retention.json": "share the same retention",
   "29-object-mixed-expiry.json": "share the same expiry",
   "30-object-mixed-explicit-timestamp.json": "share the same explicit_timestamp",
   "31-object-mixed-database-retention.json": "share the same database_retention_policy",
-  "32-object-mixed-database-retention-ttl.json": "share the same database_retention_ttl"
+  "32-object-mixed-database-retention-ttl.json": "share the same database_retention_ttl",
+  "33-endpoint-malformed-placeholder-open.json": "malformed placeholder",
+  "34-endpoint-malformed-placeholder-mid.json": "malformed placeholder",
+  "35-endpoint-malformed-placeholder-unclosed.json": "malformed placeholder",
+  "36-endpoint-invalid-placeholder-name.json": "invalid placeholder name"
 }

$ go test ./pkg/interfaceschema/ -run 'TestInvalidFixturesRejected' -v 2>&1 | tail -50
=== RUN   TestInvalidFixturesRejected/25-retention-on-properties.json
=== RUN   TestInvalidFixturesRejected/33-endpoint-malformed-placeholder-open.json
=== RUN   TestInvalidFixturesRejected/35-endpoint-malformed-placeholder-unclosed.json
=== RUN   TestInvalidFixturesRejected/20-properties-explicit-timestamp.json
=== RUN   TestInvalidFixturesRejected/28-object-mixed-retention.json
=== RUN   TestInvalidFixturesRejected/34-endpoint-malformed-placeholder-mid.json
=== RUN   TestInvalidFixturesRejected/02-name-too-long.json
=== RUN   TestInvalidFixturesRejected/08-bad-ownership.json
=== RUN   TestInvalidFixturesRejected/09-bad-aggregation.json
=== RUN   TestInvalidFixturesRejected/15-object-on-properties.json
=== RUN   TestInvalidFixturesRejected/31-object-mixed-database-retention.json
--- PASS: TestInvalidFixturesRejected (0.01s)
    --- PASS: TestInvalidFixturesRejected/01-bad-name.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/03-version-zero.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/11-endpoint-not-rooted.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/17-object-mixed-attrs.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/18-object-divergent-prefix.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/22-ttl-without-policy.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/06-missing-type.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/10-bad-value-type.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/21-allow-unset-datastream.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/24-unknown-field.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/29-object-mixed-expiry.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/05-no-mappings.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/07-bad-type.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/14-endpoint-conflict.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/27-trailing-garbage.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/32-object-mixed-database-retention-ttl.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/36-endpoint-invalid-placeholder-name.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/04-negative-major.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/13-endpoint-duplicate.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/16-object-mixed-depth.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/26-negative-expiry.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/30-object-mixed-explicit-timestamp.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/12-endpoint-partial-placeholder.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/19-object-param-leaf.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/23-use-ttl-missing-ttl.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/25-retention-on-properties.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/33-endpoint-malformed-placeholder-open.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/35-endpoint-malformed-placeholder-unclosed.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/20-properties-explicit-timestamp.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/28-object-mixed-retention.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/34-endpoint-malformed-placeholder-mid.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/02-name-too-long.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/08-bad-ownership.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/09-bad-aggregation.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/15-object-on-properties.json (0.00s)
    --- PASS: TestInvalidFixturesRejected/31-object-mixed-database-retention.json (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	0.022s

All four new fixtures pass. Now the full gate:
$ go test ./pkg/interfaceschema/ && go vet ./pkg/interfaceschema/ && gofmt -l pkg/interfaceschema/
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	0.186s

$ go vet ./... && go test ./... 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

Done: added four invalid fixtures + manifest entries for parseSegment's placeholder-syntax rejections (/%{, /%ax, /%{x → "malformed placeholder"; /%{1} → "invalid placeholder name"), closing the regression gap where a dangling "/%{" passed all tests.
Files: pkg/interfaceschema/testdata/invalid/33-endpoint-malformed-placeholder-open.json, 34-endpoint-malformed-placeholder-mid.json, 35-endpoint-malformed-placeholder-unclosed.json, 36-endpoint-invalid-placeholder-name.json, pkg/interfaceschema/testdata/invalid/manifest.json
Verified: go test ./pkg/interfaceschema/ (new fixtures + lockstep all pass), go vet ./... && go test ./... (all ok), gofmt -l pkg/interfaceschema (clean) -> pass
Unsure: nothing — message substrings were traced directly from parse.go:509-517, and the manifest-driven test asserts containment
Follow-ups: none
