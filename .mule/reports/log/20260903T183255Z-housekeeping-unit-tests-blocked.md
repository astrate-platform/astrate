slug: housekeeping-unit-tests
verdict: blocked
at:  79f9c1e
ran: 2026-09-03T18:32:55Z on DietPi in 182s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

# Todos
[•] Read housekeeping package: CreateRealm, store interface, sealer interface, ErrValidation
[ ] Write unit test file testing CreateRealm validation paths with mocks
[ ] Run gate: go vet ./... && go test ./internal/housekeeping/...

✱ Glob "internal/housekeeping/**/*.go" 3 matches
✱ Grep "CreateRealm" 75 matches
✱ Grep "ErrValidation" in internal/housekeeping · 9 matches
→ Read internal/housekeeping/service.go
→ Read internal/housekeeping/http.go
→ Read internal/housekeeping/http_test.go

# Todos
[✓] Read housekeeping package: CreateRealm, store interface, sealer interface, ErrValidation
[•] Write unit test file testing CreateRealm validation paths with mocks
[ ] Run gate: go vet ./... && go test ./internal/housekeeping/...

Now I understand the code. The three validation checks (blank name, blank JWT key, negative reg limit) return early before touching the sealer or store, so a zero-value `*Service` works fine.
← Write internal/housekeeping/service_test.go
Wrote file successfully.

✱ Grep "func Generate" in internal/pairing/ca · 1 match
→ Read internal/pairing/ca/ca.go [offset=71, limit=20]
The `ca.Generate` will work, but then `s.sealer.Seal(nil)` will panic. Let me focus the test on exactly the three required validation paths and drop the valid-input cases (they'd need real mocks).

← Write internal/housekeeping/service_test.go
Wrote file successfully.

← Edit internal/housekeeping/service_test.go
Index: /root/astrate-mule/internal/housekeeping/service_test.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service_test.go
+++ /root/astrate-mule/internal/housekeeping/service_test.go
@@ -50,11 +50,8 @@
 }
 if !errors.Is(err, ErrValidation) {
 	t.Fatalf("want ErrValidation, got: %v", err)
 }
-if !errors.Is(err, ErrValidation) {
-	t.Fatalf("want ErrValidation, got: %v", err)
-}
 got := err.Error()
 if !contains(got, tt.wantDetail) {
 	t.Errorf("error %q does not contain %q", got, tt.wantDetail)
 }

# Todos
[✓] Read housekeeping package: CreateRealm, store interface, sealer interface, ErrValidation
[✓] Write unit test file testing CreateRealm validation paths with mocks
[•] Run gate: go vet ./... && go test ./internal/housekeeping/...

The task mentions "using a mock store and sealer". The `Service` struct holds concrete `*store.Store` and `*store.KeySealer`, not interfaces, so I can't inject mocks. Instead I use a zero-value `*Service` — which is a valid alternative since the validation paths bail before touching them. Let me verify the gate.
$ go vet ./internal/housekeeping/... && go test ./internal/housekeeping/... && gofmt -l internal/housekeeping/
# github.com/astrate-platform/astrate/internal/housekeeping
# [github.com/astrate-platform/astrate/internal/housekeeping]
vet: internal/housekeeping/service_test.go:55:8: undefined: contains

`contains` is defined in the integration-gated `http_test.go`, so it's not available. I'll use `strings.Contains` instead.
← Edit internal/housekeeping/service_test.go
Index: /root/astrate-mule/internal/housekeeping/service_test.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service_test.go
+++ /root/astrate-mule/internal/housekeeping/service_test.go
@@ -50,11 +50,10 @@
 			}
 			if !errors.Is(err, ErrValidation) {
 				t.Fatalf("want ErrValidation, got: %v", err)
 			}
-			got := err.Error()
-			if !contains(got, tt.wantDetail) {
-				t.Errorf("error %q does not contain %q", got, tt.wantDetail)
+			if !strings.Contains(err.Error(), tt.wantDetail) {
+				t.Errorf("error %q does not contain %q", err.Error(), tt.wantDetail)
 			}
 		})
 	}
 }

← Edit internal/housekeeping/service_test.go
Index: /root/astrate-mule/internal/housekeeping/service_test.go
===================================================================
--- /root/astrate-mule/internal/housekeeping/service_test.go
+++ /root/astrate-mule/internal/housekeeping/service_test.go
@@ -2,8 +2,9 @@
 
 import (
 	"context"
 	"errors"
+	"strings"
 	"testing"
 )
 
 // TestCreateRealm_Validation exercises the three validation paths in

$ go vet ./internal/housekeeping/... && go test ./internal/housekeeping/... && gofmt -l internal/housekeeping/
ok  	github.com/astrate-platform/astrate/internal/housekeeping	0.028s

Gate passes (non-integration), gofmt clean. Let me verify the test actually fails without the change — i.e. it's testing something real. Since the validation checks are intrinsic to the function, removing them should break the test. The test asserts `ErrValidation`, which is returned by those early checks. This is genuine.

Let me double-check the full test file once more.
→ Read internal/housekeeping/service_test.go

# Todos
[✓] Read housekeeping package: CreateRealm, store interface, sealer interface, ErrValidation
[✓] Write unit test file testing CreateRealm validation paths with mocks
[✓] Run gate: go vet ./... && go test ./internal/housekeeping/...

The test implements the three required validation paths. Note about the task wording: "using a mock store and sealer" — the `Service` holds concrete `*store.Store` / `*store.KeySealer` (not interfaces), so mocks can't be injected. Instead I use a zero-value `*Service`, which is equivalent for these paths since they return before dereferencing any field. This is functionally sound and tests the exact same code paths.
The test asserts on `ErrValidation`, which only those validation branches produce — removing any branch would break its corresponding subtest. I did not touch git.

```
Done: added non-integration unit tests in internal/housekeeping for CreateRealm's three ErrValidation paths (blank name, blank JWT key, negative reg limit)
Files: internal/housekeeping/service_test.go
Verified: go vet ./internal/housekeeping/... && go test ./internal/housekeeping/... && gofmt -l internal/housekeeping/ -> pass (passed; non-race build)
Unsure: task said "using a mock store and sealer", but Service holds concrete *store.Store / *store.KeySealer (not interfaces), so mocks can't be injected; used a zero-value *Service instead, which the validation paths hit before dereferencing any field — same code paths, no panic
Follow-ups: none
```
