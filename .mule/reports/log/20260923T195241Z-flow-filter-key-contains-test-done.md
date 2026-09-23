slug: flow-filter-key-contains-test
verdict: done
at:  3875ba6
ran: 2026-09-23T19:52:41Z on DietPi in 107s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/flow/blocks/transform.go
→ Read internal/flow/blocks/transform_test.go
Task: add `key_contains` test coverage (substring-match, substring-absent, prefix+contains) to the filter tests.

← Edit internal/flow/blocks/transform_test.go
Index: /root/astrate-mule/internal/flow/blocks/transform_test.go
===================================================================
--- /root/astrate-mule/internal/flow/blocks/transform_test.go
+++ /root/astrate-mule/internal/flow/blocks/transform_test.go
@@ -54,8 +54,55 @@
 		t.Fatalf("drop metadata: out=%v err=%v", drop2, err)
 	}
 }
 
+func TestFilter_KeyContains(t *testing.T) {
+	b, err := blocks.Filter("f", map[string]any{"key_contains": "sensor"}, flow.Deps{})
+	if err != nil {
+		t.Fatal(err)
+	}
+
+	// Substring present: message passes through.
+	out, err := b.Process(&flow.Message{Key: "acme/sensor-valves", Type: flow.TypeString, Data: "x"})
+	if err != nil || len(out) != 1 || out[0].Key != "acme/sensor-valves" {
+		t.Fatalf("match: out=%v err=%v", out, err)
+	}
+
+	// Substring absent: message dropped.
+	out, err = b.Process(&flow.Message{Key: "acme/actuator", Type: flow.TypeString, Data: "x"})
+	if err != nil || len(out) != 0 {
+		t.Fatalf("no match: out=%v err=%v", out, err)
+	}
+}
+
+func TestFilter_PrefixAndContains(t *testing.T) {
+	b, err := blocks.Filter("f", map[string]any{
+		"key_prefix":   "acme/",
+		"key_contains": "sensor",
+	}, flow.Deps{})
+	if err != nil {
+		t.Fatal(err)
+	}
+
+	// Both conditions satisfied: passes through.
+	out, err := b.Process(&flow.Message{Key: "acme/dev-sensor-1", Type: flow.TypeString, Data: "x"})
+	if err != nil || len(out) != 1 {
+		t.Fatalf("both: out=%v err=%v", out, err)
+	}
+
+	// Prefix yes, substring no: dropped.
+	out, err = b.Process(&flow.Message{Key: "acme/actuator", Type: flow.TypeString, Data: "x"})
+	if err != nil || len(out) != 0 {
+		t.Fatalf("prefix only: out=%v err=%v", out, err)
+	}
+
+	// Substring yes, prefix no: dropped.
+	out, err = b.Process(&flow.Message{Key: "other-sensor-valve", Type: flow.TypeString, Data: "x"})
+	if err != nil || len(out) != 0 {
+		t.Fatalf("contains only: out=%v err=%v", out, err)
+	}
+}
+
 func TestFilter_Type(t *testing.T) {
 	b, err := blocks.Filter("f", map[string]any{"type": "integer"}, flow.Deps{})
 	if err != nil {
 		t.Fatal(err)

Now run the gate and gofmt:
$ gofmt -l . && go test ./internal/flow/blocks/... && go vet ./internal/flow/blocks/...
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	1.299s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)

$ go test ./... 2>&1 | grep -Ev "^ok|no test files" ; echo "exit: $?"
exit: 1

$ go test ./... 2>&1 | rg -v "^ok|no test files" ; echo "gateexit: $?"
gateexit: 1

Done: Added `key_contains` filter test coverage — substring-match, substring-absent, and prefix+contains combined rows.
Files: internal/flow/blocks/transform_test.go
Verified: gofmt -l . -> pass; go test ./... -> pass; go vet ./internal/flow/blocks/... -> pass (ran without -race, container-free unit tests)
Unsure: nothing
Follow-ups: none
