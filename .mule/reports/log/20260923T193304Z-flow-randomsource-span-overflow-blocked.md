slug: flow-randomsource-span-overflow
verdict: blocked
at:  2e1c216
ran: 2026-09-23T19:33:04Z on DietPi in 288s


> build · big-pickle

I'll start by reading the two files as instructed.
→ Read .mule/MULE.md
→ Read .mule/task.md
Task means: reject a randomSource span wider than `math.MaxInt64` at construction, and add a test. Let me read the named file.
→ Read internal/flow/blocks/randomsource.go
Now let me look at the existing tests for this block to match style.
✱ Glob "internal/flow/blocks/*random*" 2 matches
✱ Grep "RandomSource|random_source|randomsource" in internal/flow/blocks · 19 matches
→ Read internal/flow/blocks/randomsource_test.go
Task: reject a `randomSource` integer span wider than `math.MaxInt64` at config parsing (right where `min <= max` is checked), plus a construct test for that config. The overflow is real: `s.maxInt-s.minInt+1` at randomsource.go:171 wraps non-positive → `rand.Int64N` panics in `Emit`.

Edit 1 — the construction check:

← Edit internal/flow/blocks/randomsource.go
Index: /root/astrate-mule/internal/flow/blocks/randomsource.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/randomsource.go
+++ /root/astrate-mule/internal/flow/blocks/randomsource.go
@@ -105,8 +105,15 @@
 	}
 	if (haveMin || haveMax) && fmin > fmax {
 		return cfg, fmt.Errorf("min must be <= max")
 	}
+	// The inclusive integer span s.maxInt-s.minInt+1 (randomsource.next) must
+	// fit in int64: a wider span overflows to a non-positive n, which makes
+	// rand.Int64N panic. Only the integer path computes that value, so reals
+	// keep arbitrary-width ranges.
+	if cfg.kind == randomInteger && uint64(cfg.maxInt)-uint64(cfg.minInt) >= math.MaxInt64 {
+		return cfg, fmt.Errorf("min and max span wider than the int64 range")
+	}
 	cfg.key = stringConfig(config, "key", "random")
 	return cfg, nil
 }
 

← Edit internal/flow/blocks/randomsource.go
Index: /root/astrate-mule/internal/flow/blocks/randomsource.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/randomsource.go
+++ /root/astrate-mule/internal/flow/blocks/randomsource.go
@@ -109,8 +109,12 @@
 // The inclusive integer span s.maxInt-s.minInt+1 (randomsource.next) must
 // fit in int64: a wider span overflows to a non-positive n, which makes
 // rand.Int64N panic. Only the integer path computes that value, so reals
 // keep arbitrary-width ranges.
+// The inclusive integer span s.maxInt-s.minInt+1 (randomsource.next) must
+// fit in int64: a wider span overflows to a non-positive n, which makes
+// rand.Int64N panic. Only the integer path computes that value, so reals
+// keep arbitrary-width ranges.
 if cfg.kind == randomInteger && uint64(cfg.maxInt)-uint64(cfg.minInt) >= math.MaxInt64 {
 	return cfg, fmt.Errorf("min and max span wider than the int64 range")
 }
 cfg.key = stringConfig(config, "key", "random")

Let me verify the current state of the file:
→ Read internal/flow/blocks/randomsource.go [offset=100, limit=22]
Duplicate comment block — fixing:

← Edit internal/flow/blocks/randomsource.go
Index: /root/astrate-mule/internal/flow/blocks/randomsource.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/randomsource.go
+++ /root/astrate-mule/internal/flow/blocks/randomsource.go
@@ -109,12 +109,8 @@
 // The inclusive integer span s.maxInt-s.minInt+1 (randomsource.next) must
 // fit in int64: a wider span overflows to a non-positive n, which makes
 // rand.Int64N panic. Only the integer path computes that value, so reals
 // keep arbitrary-width ranges.
-// The inclusive integer span s.maxInt-s.minInt+1 (randomsource.next) must
-// fit in int64: a wider span overflows to a non-positive n, which makes
-// rand.Int64N panic. Only the integer path computes that value, so reals
-// keep arbitrary-width ranges.
 if cfg.kind == randomInteger && uint64(cfg.maxInt)-uint64(cfg.minInt) >= math.MaxInt64 {
 	return cfg, fmt.Errorf("min and max span wider than the int64 range")
 }
 cfg.key = stringConfig(config, "key", "random")

Now the test — add the wide-span rejection row to the construct table:
← Edit internal/flow/blocks/randomsource_test.go
Index: /root/astrate-mule/internal/flow/blocks/randomsource_test.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/randomsource_test.go
+++ /root/astrate-mule/internal/flow/blocks/randomsource_test.go
@@ -46,8 +46,14 @@
 		bad:  map[string]any{"type": "integer", "min": 10, "max": 9},
 		good: map[string]any{"type": "integer", "min": 9, "max": 10},
 		want: `random_source: min must be <= max`,
 	},
+	{
+		name: "integer span wider than int64",
+		bad:  map[string]any{"type": "integer", "min": -9e18, "max": 9e18},
+		good: map[string]any{"type": "integer", "min": -9e18, "max": 0},
+		want: `random_source: min and max span wider than the int64 range`,
+	},
 }
 for _, tt := range tests {
 	t.Run(tt.name+"/reject", func(t *testing.T) {
 		_, err := blocks.RandomSource("r", tt.bad, flow.Deps{})

Now run the gate (container-free). Note MULE.md: this machine cannot run `-race`, so `go vet` + `go test` is the gate.
$ go vet ./... && go test ./internal/flow/blocks/ && gofmt -l .
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	1.261s

$ go test ./... 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/flow	0.655s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.575s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	0.105s
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

Quick probe to confirm the overflow math (that the span indeed overflows to a non-positive Int64N argument):
$ cat > /tmp/probe.go <<'EOF'
package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

func main() {
	min := int64(-9e18)
	max := int64(9e18)
	n := max - min + 1
	fmt.Printf("min=%d max=%d n(overflowed)=%d valid=%v\n", min, max, n, n > 0 && n <= math.MaxInt64)
	// boundary: span == MaxInt64 must also be rejected
	min2, max2 := int64(math.MinInt64), int64(-1)
	fmt.Printf("MinInt64..-1 span n=%d\n", max2-min2+1)
}
EOF
go run /tmp/probe.go
# command-line-arguments
/tmp/probe.go:6:2: "math/rand/v2" imported as rand and not used

$ printf '%s\n' 'package main' '' 'import (' '	"fmt"' '	"math"' ')' '' 'func main() {' '	min := int64(-9e18)' '	max := int64(9e18)' '	n := max - min + 1' '	fmt.Printf("min=%d max=%d n(overflowed)=%d valid=%v\\n", min, max, n, n > 0 && n <= math.MaxInt64)' '	min2, max2 := int64(math.MinInt64), int64(-1)' '	fmt.Printf("MinInt64..-1 n=%d\\n", max2-min2+1)' '}' > /tmp/probe.go && go run /tmp/probe.go
min=-9000000000000000000 max=9000000000000000000 n(overflowed)=-446744073709551615 valid=false\nMinInt64..-1 n=-9223372036854775808\n

Confirmed: `-9e18..9e18` overflows to n=-446744073709551615 (non-positive → `rand.Int64N` panics); boundary span `== MaxInt64` also overflows (n=-2^63), so `>= math.MaxInt64` is the right cutoff. My check rejects both; span `MaxInt64-1` (e.g. `MinInt64..-2`) stays valid.

The construct-rejection test row fails without the change (construction of the wide span would succeed there, so the subtest's error match fails), which is exactly what the verifier requires.

Done: randomSource now rejects an integer min/max span of `math.MaxInt64` or wider at config parse, preventing the `rand.Int64N(n<=0)` panic in the source pump
Files: internal/flow/blocks/randomsource.go, internal/flow/blocks/randomsource_test.go
Verified: `go vet ./...`, `go test ./...`, `gofmt -l .` -> pass (whole repo green); overflow math probed with throwaway program in /tmp
Unsure: gated the new check to `type: integer` only, so wide real ranges still parse as before — please confirm that's acceptable rather than making the check type-agnostic
Follow-ups: none
