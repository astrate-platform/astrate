slug: flowapi-autorestart-terminal-failure
verdict: done
at:  8ce04d2
ran: 2026-09-10T17:16:29Z on DietPi in 251s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

# Todos
[•] Read internal/flowapi/service.go around restartWithBackoff (lines 360-621)
[ ] Understand error types: ErrValidation, store/manager errors
[ ] Write pure helper to classify terminal vs transient errors
[ ] Modify restartWithBackoff to stop on terminal failures
[ ] Write container-free unit test for the helper
[ ] Run gate: go test ./... and gofmt -l .

→ Read internal/flowapi/service.go [offset=340, limit=300]
✱ Grep "ErrValidation" in . · 97 matches

# Todos
[✓] Read internal/flowapi/service.go around restartWithBackoff (lines 360-621)
[•] Understand error types: ErrValidation, store/manager errors
[ ] Write pure helper to classify terminal vs transient errors
[ ] Modify restartWithBackoff to stop on terminal failures
[ ] Write container-free unit test for the helper
[ ] Run gate: go test ./... and gofmt -l .

✱ Grep "var Err(NotFound|AlreadyExists|FlowExists)" in . · 0 matches
✱ Glob "internal/flowapi/*_test.go" 5 matches
✱ Grep "var ErrNotFound|var ErrFlowExists|var ErrAlreadyExists" in . · 0 matches
→ Read internal/flowapi/service_test.go
✱ Grep "var Err" in internal/store · 2 matches
✱ Grep "var Err" in internal/flow · 2 matches
✱ Grep "ErrNotFound|ErrAlreadyExists" in internal/store · 100 matches
→ Read internal/flowapi/service.go [limit=130]
→ Read internal/store/store.go [limit=70]

# Todos
[✓] Read internal/flowapi/service.go around restartWithBackoff (lines 360-621)
[✓] Understand error types: ErrValidation, store/manager errors
[•] Write pure helper to classify terminal vs transient errors
[ ] Modify restartWithBackoff to stop on terminal failures
[ ] Write container-free unit test for the helper
[ ] Run gate: go test ./... and gofmt -l .

The task: make `restartWithBackoff` stop retrying when a rebuild failure is permanent (`ErrValidation`), using a pure helper, plus a container-free unit test on it.
← Edit internal/flowapi/service.go
Index: /root/astrate-mule/internal/flowapi/service.go
===================================================================
--- /root/astrate-mule/internal/flowapi/service.go
+++ /root/astrate-mule/internal/flowapi/service.go
@@ -581,10 +581,21 @@
 	defer s.restartMu.Unlock()
 	delete(s.restarting, instanceID)
 }
 
+// isTerminalRestartFailure reports whether a RestartFlowInstance rebuild
+// error is permanent: ErrValidation means the flow's pipeline/config no
+// longer resolves (deleted pipeline, invalid stored definition, unknown
+// block), so no number of retries can succeed. Store/manager errors are
+// transient and keep the backoff loop alive.
+func isTerminalRestartFailure(err error) bool {
+	return errors.Is(err, ErrValidation)
+}
+
 // restartWithBackoff keeps re-running RestartFlowInstance until it succeeds,
-// the durable row disappears (flow deleted), or auto_restart is turned off.
+// the durable row disappears (flow deleted), auto_restart is turned off, or
+// the rebuild fails with a terminal error (ErrValidation — the flow marks
+// itself failed and there is nothing left to retry).
 func (s *Service) restartWithBackoff(realm string, realmID int16, name, instanceID string) {
 	defer s.endRestart(instanceID)
 	bg := context.Background()
 	delay := restartBackoffBase

← Edit internal/flowapi/service.go
Index: /root/astrate-mule/internal/flowapi/service.go
===================================================================
--- /root/astrate-mule/internal/flowapi/service.go
+++ /root/astrate-mule/internal/flowapi/service.go
@@ -625,8 +625,13 @@
 			s.log.Info("flow auto-restarted after block death",
 				"realm", realm, "flow", name, "attempt", attempt+1)
 			return
 		}
+		if isTerminalRestartFailure(err) {
+			s.log.Error("flow auto-restart stopped: rebuild failure is permanent",
+				"realm", realm, "flow", name, "attempt", attempt+1, "error", err)
+			return
+		}
 		s.log.Warn("flow auto-restart attempt failed",
 			"realm", realm, "flow", name, "attempt", attempt+1, "next_in", delay.String(), "error", err)
 	}
 }

Now the container-free unit test for the helper:
← Edit internal/flowapi/service_test.go
Index: /root/astrate-mule/internal/flowapi/service_test.go
===================================================================
--- /root/astrate-mule/internal/flowapi/service_test.go
+++ /root/astrate-mule/internal/flowapi/service_test.go
@@ -2,8 +2,10 @@
 
 import (
 	"context"
 	"encoding/json"
+	"errors"
+	"fmt"
 	"testing"
 	"time"
 
 	"github.com/astrate-platform/astrate/internal/flow"

← Edit internal/flowapi/service_test.go
Index: /root/astrate-mule/internal/flowapi/service_test.go
===================================================================
--- /root/astrate-mule/internal/flowapi/service_test.go
+++ /root/astrate-mule/internal/flowapi/service_test.go
@@ -119,8 +119,36 @@
 		t.Fatal("PipelineID should alias InstanceID")
 	}
 }
 
+// TestIsTerminalRestartFailure pins the auto-restart classification:
+// ErrValidation from the rebuild path (deleted pipeline, unresolvable
+// config/definition) is permanent and must stop the backoff loop, while
+// store/manager errors are transient and must keep it retrying.
+func TestIsTerminalRestartFailure(t *testing.T) {
+	cases := []struct {
+		name string
+		err  error
+		want bool
+	}{
+		{"nil", nil, false},
+		{"pipeline deleted", fmt.Errorf("%w: pipeline %q not found", ErrValidation, "device-to-http"), true},
+		{"invalid stored config", fmt.Errorf("%w: config must be a JSON object: %v", ErrValidation, errors.New("unexpected end")), true},
+		{"definition does not parse", fmt.Errorf("%w: %v", ErrValidation, errors.New("syntax error")), true},
+		{"plain error", errors.New("postgres: connection refused"), false},
+		{"store not found", fmt.Errorf("%w: flow %q", store.ErrNotFound, "f1"), false},
+		{"store already exists", fmt.Errorf("%w: flow %q", store.ErrAlreadyExists, "f1"), false},
+		{"wrapped non-validation", fmt.Errorf("resolve: %w", errors.New("boom")), false},
+	}
+	for _, tc := range cases {
+		t.Run(tc.name, func(t *testing.T) {
+			if got := isTerminalRestartFailure(tc.err); got != tc.want {
+				t.Errorf("isTerminalRestartFailure(%v) = %v, want %v", tc.err, got, tc.want)
+			}
+		})
+	}
+}
+
 // TestFlowDepsCarriesIngest pins the #84 plumbing: every built instance's
 // block Deps carry the service-level IngestFunc, so a virtual_device_pool
 // constructor sees it without any per-pipeline configuration.
 func TestFlowDepsCarriesIngest(t *testing.T) {

# Todos
[✓] Read internal/flowapi/service.go around restartWithBackoff (lines 360-621)
[✓] Understand error types: ErrValidation, store/manager errors
[✓] Write pure helper to classify terminal vs transient errors
[✓] Modify restartWithBackoff to stop on terminal failures
[✓] Write container-free unit test for the helper
[•] Run gate: go test ./... and gofmt -l .

Now the gate — vet, tests, and gofmt:
$ gofmt -l .
(no output)

$ go vet ./... && go test ./... 2>&1 | tail -30
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	0.091s
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

$ go test ./internal/flowapi/ -run TestIsTerminalRestartFailure -v 2>&1 | tail -25
=== RUN   TestIsTerminalRestartFailure
=== RUN   TestIsTerminalRestartFailure/nil
=== RUN   TestIsTerminalRestartFailure/pipeline_deleted
=== RUN   TestIsTerminalRestartFailure/invalid_stored_config
=== RUN   TestIsTerminalRestartFailure/definition_does_not_parse
=== RUN   TestIsTerminalRestartFailure/plain_error
=== RUN   TestIsTerminalRestartFailure/store_not_found
=== RUN   TestIsTerminalRestartFailure/store_already_exists
=== RUN   TestIsTerminalRestartFailure/wrapped_non-validation
--- PASS: TestIsTerminalRestartFailure (0.00s)
    --- PASS: TestIsTerminalRestartFailure/nil (0.00s)
    --- PASS: TestIsTerminalRestartFailure/pipeline_deleted (0.00s)
    --- PASS: TestIsTerminalRestartFailure/invalid_stored_config (0.00s)
    --- PASS: TestIsTerminalRestartFailure/definition_does_not_parse (0.00s)
    --- PASS: TestIsTerminalRestartFailure/plain_error (0.00s)
    --- PASS: TestIsTerminalRestartFailure/store_not_found (0.00s)
    --- PASS: TestIsTerminalRestartFailure/store_already_exists (0.00s)
    --- PASS: TestIsTerminalRestartFailure/wrapped_non-validation (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/flowapi	0.083s

# Todos
[✓] Read internal/flowapi/service.go around restartWithBackoff (lines 360-621)
[✓] Understand error types: ErrValidation, store/manager errors
[✓] Write pure helper to classify terminal vs transient errors
[✓] Modify restartWithBackoff to stop on terminal failures
[✓] Write container-free unit test for the helper
[✓] Run gate: go test ./... and gofmt -l .

Done: `restartWithBackoff` now stops on permanent rebuild failures — a pure helper `isTerminalRestartFailure` (errors.Is vs ErrValidation) short-circuits the backoff loop, while store/manager errors stay transient and keep retrying.
Files: internal/flowapi/service.go (helper + loop exit, ~592-606, ~642), internal/flowapi/service_test.go (TestIsTerminalRestartFailure, 8 table cases)
Verified: `go vet ./... && go test ./...` (no -race) -> pass; `go test ./internal/flowapi -run TestIsTerminalRestartFailure -v` -> pass; `gofmt -l .` -> clean
Unsure: nothing — test is compile-bound to the helper, so removing the implementation breaks the gate
Follow-ups: none
