slug: errorname-missing-required-test
verdict: done
at:  739df67
ran: 2026-09-20T18:11:28Z on DietPi in 135s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/engine/triggers/errorname_test.go
→ Read internal/engine/triggers/errorname.go
Task: add the missing `missing_required` table row and an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()`. Both are test-only additions in `errorname_test.go`.

← Edit internal/engine/triggers/errorname_test.go
Index: /root/astrate-mule/internal/engine/triggers/errorname_test.go
===================================================================
--- /root/astrate-mule/internal/engine/triggers/errorname_test.go
+++ /root/astrate-mule/internal/engine/triggers/errorname_test.go
@@ -26,8 +26,9 @@
 	{"bad_timestamp", "unexpected_value_type"},
 	{"type_mismatch", "unexpected_value_type"},
 	{"value_too_large", "value_size_exceeded"},
 	{"bad_object", "unexpected_object_key"},
+	{"missing_required", "unexpected_object_key"},
 	{"unset_not_allowed", "unexpected_value_type"},
 }
 for _, tt := range tests {
 	t.Run(tt.input, func(t *testing.T) {

← Edit internal/engine/triggers/errorname_test.go
Index: /root/astrate-mule/internal/engine/triggers/errorname_test.go
===================================================================
--- /root/astrate-mule/internal/engine/triggers/errorname_test.go
+++ /root/astrate-mule/internal/engine/triggers/errorname_test.go
@@ -70,8 +70,21 @@
 		})
 	}
 }
 
+func TestAstrateToUpstreamValuesInClosedSet(t *testing.T) {
+	for key, value := range astrateToUpstream {
+		t.Run(key, func(t *testing.T) {
+			for _, n := range UpstreamErrorNames() {
+				if n == value {
+					return
+				}
+			}
+			t.Errorf("astrateToUpstream[%q] = %q is not in UpstreamErrorNames()", key, value)
+		})
+	}
+}
+
 func TestNewDeviceErrorEventWiring(t *testing.T) {
 	t.Run("translated", func(t *testing.T) {
 		evt := NewDeviceErrorEvent("unknown_format", map[string]string{"detail": "x"})
 		if evt.ErrorName != "undecodable_bson_payload" {

Now run the gate for this package:
$ go test -race ./internal/engine/triggers/ && gofmt -l internal/engine/triggers/errorname_test.go
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/engine/triggers	0.066s
FAIL

Race detector unavailable (39-bit VMA, as MULE.md warns). Running without `-race`:
$ go test ./internal/engine/triggers/ && gofmt -l internal/engine/triggers/errorname_test.go
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	1.988s

$ go vet ./internal/engine/triggers/
(no output)

Done. Both additions verified.

```
Done: Added missing_required→unexpected_object_key row to TestUpstreamErrorNameMapping and TestAstrateToUpstreamValuesInClosedSet invariant (every astrateToUpstream value must be in UpstreamErrorNames())
Files: internal/engine/triggers/errorname_test.go
Verified: go test -race ./internal/engine/triggers/ -> failed (TSan unsupported VMA 39, Pi env); go test ./internal/engine/triggers/ -> pass; go vet ./internal/engine/triggers/ -> pass; gofmt -l -> clean
Unsure: Tests-only task — both new assertions hold against the current map (the mapping and values are already correct), so no component of this change fails without the other; it is coverage/invariant hardening, not behavior-change-proof per MULE.md rule 4
Follow-ups: none
```
