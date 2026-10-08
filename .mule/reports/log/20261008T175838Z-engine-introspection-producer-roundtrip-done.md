slug: engine-introspection-producer-roundtrip
verdict: done
at:  d070345
ran: 2026-10-08T17:58:38Z on DietPi in 234s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/testutil/astartedevice.go [offset=70, limit=50]
→ Read internal/engine/introspection.go
→ Read internal/engine/introspection_test.go
✱ Grep "astrate/internal/engine" in internal/testutil · 0 matches
→ Read internal/testutil/astartedevice.go [limit=40]
No import cycle; `testutil` is container-free. Adding the roundtrip pin test.

← Edit internal/engine/introspection_test.go
Index: /root/astrate-mule/internal/engine/introspection_test.go
===================================================================
--- /root/astrate-mule/internal/engine/introspection_test.go
+++ /root/astrate-mule/internal/engine/introspection_test.go
@@ -7,8 +7,9 @@
 
 	promtest "github.com/prometheus/client_golang/prometheus/testutil"
 
 	"github.com/astrate-platform/astrate/internal/store"
+	"github.com/astrate-platform/astrate/internal/testutil"
 	"github.com/astrate-platform/astrate/pkg/payload"
 )
 
 // TestParseIntrospection covers the docs/DESIGN.md §3.3 introspection

← Edit internal/engine/introspection_test.go
Index: /root/astrate-mule/internal/engine/introspection_test.go
===================================================================
--- /root/astrate-mule/internal/engine/introspection_test.go
+++ /root/astrate-mule/internal/engine/introspection_test.go
@@ -78,8 +78,50 @@
 		})
 	}
 }
 
+// TestIntrospectionProducerRoundtrip pins the producer half of the
+// docs/DESIGN.md §3.3 rule on the default gate: testutil.Introspection (the
+// harness renderer, normally exercised only behind e2e/integration build
+// tags) emits a deterministic sorted `name:major:minor;…` string with no
+// trailing `;`, and parseIntrospection accepts exactly that string back.
+func TestIntrospectionProducerRoundtrip(t *testing.T) {
+	entries := map[string][2]int{
+		"com.ex.Zeta":  {1, 0},
+		"com.ex.Alpha": {0, 3},
+		"com.ex.Mid":   {2, 7},
+	}
+	s := testutil.Introspection(entries)
+	const want = "com.ex.Alpha:0:3;com.ex.Mid:2:7;com.ex.Zeta:1:0"
+	if s != want {
+		t.Fatalf("Introspection() = %q, want %q", s, want)
+	}
+	got, err := parseIntrospection(s)
+	if err != nil {
+		t.Fatalf("parseIntrospection(%q): %v", s, err)
+	}
+	wantParsed := map[string]store.InterfaceVersion{
+		"com.ex.Alpha": {Major: 0, Minor: 3},
+		"com.ex.Mid":   {Major: 2, Minor: 7},
+		"com.ex.Zeta":  {Major: 1, Minor: 0},
+	}
+	if len(got) != len(wantParsed) {
+		t.Fatalf("parseIntrospection(%q) = %v, want %v", s, got, wantParsed)
+	}
+	for name, v := range wantParsed {
+		if got[name] != v {
+			t.Errorf("entry %q = %+v, want %+v", name, got[name], v)
+		}
+	}
+
+	if e := testutil.Introspection(map[string][2]int{}); e != "" {
+		t.Errorf("Introspection(empty) = %q, want empty string", e)
+	}
+	if e := testutil.Introspection(nil); e != "" {
+		t.Errorf("Introspection(nil) = %q, want empty string", e)
+	}
+}
+
 // TestIntrospectionHandler drives the wired handler end to end: persistence,
 // device-cache update (the §2.6 step-1 gate follows immediately), and the
 // broker ACL refresh.
 func TestIntrospectionHandler(t *testing.T) {

$ gofmt -l . && go vet ./internal/engine/ ./internal/testutil/
(no output)

$ go test -race ./internal/engine/ 2>&1 | tail -5
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/engine	0.023s
FAIL

Race detector unsupported here (39-bit VMA, as MULE.md predicts). Falling back to `go test ./...`.
$ go test ./... 2>&1 | tail -30
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	1.534s
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
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
ok  	github.com/astrate-platform/astrate/internal/testutil	0.070s
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

$ go vet ./... && gofmt -l . && go test -run 'TestIntrospectionProducerRoundtrip|TestParseIntrospection' -count=3 ./internal/engine/
ok  	github.com/astrate-platform/astrate/internal/engine	0.063s

Done: pinned the introspection producer/consumer roundtrip — `testutil.Introspection` on an unsorted multi-entry map renders the exact sorted `name:major:minor;…` string (no trailing `;`), `parseIntrospection` parses it back to the same map, and empty/nil map renders `""`.
Files: internal/engine/introspection_test.go
Verified: `gofmt -l .` -> clean; `go vet ./...` -> pass; `go test ./...` -> pass; `go test -race` unsupported on this machine (39-bit VMA, fell back to non-race as MULE.md allows); new test run `-count=3` -> pass.
Unsure: nothing.
Follow-ups: none.
