slug: interfaceschema-enum-roundtrip-coverage
verdict: done
at:  f955958
ran: 2026-09-11T20:16:04Z on DietPi in 185s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read pkg/interfaceschema/parse_test.go
✱ Glob "pkg/interfaceschema/*.go" 12 matches
✱ Grep "type (Ownership|Aggregation|Retention|DatabaseRetentionPolicy)|func \(.*(Ownership|Aggregation|Retention|DatabaseRetentionPolicy)\) (String|MarshalJSON|UnmarshalJSON)" in pkg/interfaceschema · 16 matches
→ Read pkg/interfaceschema/types.go
✱ Glob ".mule/**/*.md" 0 matches
✱ Glob ".mule/**/*.{sh,py}" 0 matches

$ ls -la .mule/ && ls -la .mule/todo.md 2>/dev/null
total 2908
drwxr-xr-x  7 root root    4096 Sep 11 22:12 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      14 Sep 11 22:12 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 11 22:04 .heartbeat
-rw-r--r--  1 root root       0 Sep 11 22:12 .last-output
-rw-r--r--  1 root root       0 Sep 11 22:12 .lock
-rw-r--r--  1 root root      12 Sep 11 21:33 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 2803695 Sep 11 22:12 cron.log
drwxr-xr-x  2 root root    4096 Sep 11 20:29 failed
-rw-r--r--  1 root root   27682 Sep 11 21:04 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    4896 Sep 11 22:04 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 11 21:40 reviews
-rw-r--r--  1 root root     450 Sep 11 22:12 task.md
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   41326 Sep 11 22:04 todo.md
-rw-r--r-- 1 root root 41326 Sep 11 22:04 .mule/todo.md

✱ Grep "interfaceschema-enum-roundtrip-coverage" in .mule/todo.md · 3 matches

$ go version && go test ./pkg/interfaceschema/ -coverprofile=/tmp/isan_cov.out -covermode=count >/dev/null 2>&1 && go tool cover -func=/tmp/isan_cov.out | rg 'Retention|DatabaseRetentionPolicy|Ownership|Aggregation'
go version go1.26.7 linux/arm64
github.com/astrate-platform/astrate/pkg/interfaceschema/parse.go:599:		checkObjectAggregation		100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:73:		ParseOwnership			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:118:		ParseAggregation		100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:218:		ParseRetention			60.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:267:		ParseDatabaseRetentionPolicy	75.0%

$ go tool cover -func=/tmp/isan_cov.out | rg 'types.go:(85|97|100|130|142|145|232|246|249|279|291|296)'
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:85:		String				75.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:97:		UnmarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:100:		MarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:130:		String				75.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:142:		UnmarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:145:		MarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:232:		String				0.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:246:		UnmarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:249:		MarshalJSON			0.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:279:		String				0.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:291:		UnmarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:296:		MarshalJSON			0.0%

Confirmed the coverage gaps. Now extending the test:
← Edit pkg/interfaceschema/parse_test.go
Index: /root/astrate-mule/pkg/interfaceschema/parse_test.go
===================================================================
--- /root/astrate-mule/pkg/interfaceschema/parse_test.go
+++ /root/astrate-mule/pkg/interfaceschema/parse_test.go
@@ -272,8 +272,44 @@
 		err := json.Unmarshal(b, &out)
 		return out.String(), err
 	})
 }
+for _, v := range []interfaceschema.Ownership{
+	interfaceschema.OwnershipDevice, interfaceschema.OwnershipServer,
+} {
+	check(t, v.String(), v, func(b []byte) (string, error) {
+		var out interfaceschema.Ownership
+		err := json.Unmarshal(b, &out)
+		return out.String(), err
+	})
+}
+for _, v := range []interfaceschema.Aggregation{
+	interfaceschema.AggregationIndividual, interfaceschema.AggregationObject,
+} {
+	check(t, v.String(), v, func(b []byte) (string, error) {
+		var out interfaceschema.Aggregation
+		err := json.Unmarshal(b, &out)
+		return out.String(), err
+	})
+}
+for _, v := range []interfaceschema.Retention{
+	interfaceschema.RetentionDiscard, interfaceschema.RetentionVolatile, interfaceschema.RetentionStored,
+} {
+	check(t, v.String(), v, func(b []byte) (string, error) {
+		var out interfaceschema.Retention
+		err := json.Unmarshal(b, &out)
+		return out.String(), err
+	})
+}
+for _, v := range []interfaceschema.DatabaseRetentionPolicy{
+	interfaceschema.NoTTL, interfaceschema.UseTTL,
+} {
+	check(t, v.String(), v, func(b []byte) (string, error) {
+		var out interfaceschema.DatabaseRetentionPolicy
+		err := json.Unmarshal(b, &out)
+		return out.String(), err
+	})
+}
 for vt := interfaceschema.Double; vt <= interfaceschema.DateTimeArray; vt++ {
 	check(t, vt.String(), vt, func(b []byte) (string, error) {
 		var out interfaceschema.ValueType
 		err := json.Unmarshal(b, &out)

$ gofmt -l pkg/interfaceschema/
(no output)

$ go test ./pkg/interfaceschema/ -run 'TestEnumRoundTrips' -v 2>&1 | rg 'PASS|FAIL|ok' | head -20
--- PASS: TestEnumRoundTrips (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	0.008s

$ go test ./pkg/interfaceschema/ -coverprofile=/tmp/isan_cov2.out -covermode=count >/dev/null 2>&1 && go tool cover -func=/tmp/isan_cov2.out | rg 'types.go:(85|100|130|145|232|249|279|296)'
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:85:		String				75.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:100:		MarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:130:		String				75.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:145:		MarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:232:		String				80.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:249:		MarshalJSON			100.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:279:		String				75.0%
github.com/astrate-platform/astrate/pkg/interfaceschema/types.go:296:		MarshalJSON			100.0%

$ go vet ./... && go test ./... 2>&1 | rg -v '^ok|no test files' | head -20; echo "gate exit: $?"
gate exit: 0

Gate green, coverage confirmed on all four enums.

**Report:**

Done: Extended TestEnumRoundTrips with round-trip coverage for Ownership, Aggregation, Retention, DatabaseRetentionPolicy — the four omitted enums (Retention/String+MarshalJSON and DatabaseRetentionPolicy/String+MarshalJSON were at 0% stmt coverage).
Files: pkg/interfaceschema/parse_test.go
Verified: `go vet ./... && go test ./...` (no -race, non-race gate per mule rule — this task touches no goroutines) -> pass; `gofmt -l pkg/interfaceschema/` -> clean; per-func cover -> Retention.String 0→80%, Retention.MarshalJSON 0→100%, DatabaseRetentionPolicy.String 0→75%, MarshalJSON 0→100%, Ownership/Aggregation String now hit (75%, is the unmatched default-error branch that stays as-is)
Unsure: this is a test-only change, so the runner's "tests must break when implementation is reverted" check has no implementation diff to revert; the guard's value is that any future typo in String/Parse/MarshalJSON (e.g. "descard", "obilver") makes marshal→unmarshal disagree and the suite fails
Follow-ups: none
