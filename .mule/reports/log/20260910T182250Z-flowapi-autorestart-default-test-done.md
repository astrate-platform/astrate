slug: flowapi-autorestart-default-test
verdict: done
at:  e272cfa
ran: 2026-09-10T18:22:50Z on DietPi in 148s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
I understand the task: extract the auto_restart default logic from `internal/flowapi/http.go:137-144` into a pure helper and write container-free tests for all three cases. Let me read the code first.
→ Read internal/flowapi/http.go
✱ Glob "internal/flowapi/*_test.go" 5 matches
Now let me check for existing test files and the `CreateFlowRequest` type:
✱ Grep "type CreateFlowRequest" in internal/flowapi · 1 match
→ Read internal/flowapi/http_test.go
Let me check the existing test file to understand the testing style:
→ Read internal/flowapi/service.go [offset=85, limit=20]
Good, I have full context. The logic at lines 141-144 is simple: `nil → true`, otherwise honor the value. I'll extract that into a pure function and add table-driven tests.

← Edit internal/flowapi/http.go
Index: /root/astrate-mule/internal/flowapi/http.go
===================================================================
--- /root/astrate-mule/internal/flowapi/http.go
+++ /root/astrate-mule/internal/flowapi/http.go
@@ -137,12 +137,9 @@
 if body.Name == "" || body.Pipeline == "" {
 	_ = astarteapi.WriteBadRequest(w)
 	return
 }
-autoRestart := true
-if body.AutoRestart != nil {
-	autoRestart = *body.AutoRestart
-}
+autoRestart := resolveAutoRestart(body.AutoRestart)
 view, err := a.svc.CreateAndStartFlow(r.Context(), r.PathValue("realm"), CreateFlowRequest{
 	Name:        body.Name,
 	Pipeline:    body.Pipeline,
 	Config:      body.Config,

Now add the helper function. I'll place it near the other unexported helpers at the bottom:
← Edit internal/flowapi/http.go
Index: /root/astrate-mule/internal/flowapi/http.go
===================================================================
--- /root/astrate-mule/internal/flowapi/http.go
+++ /root/astrate-mule/internal/flowapi/http.go
@@ -290,8 +290,17 @@
 		_ = astarteapi.WriteInternalServerError(w)
 	}
 }
 
+// resolveAutoRestart returns the effective auto_restart value from a request
+// body. Missing (nil) defaults to true; a present value is honored as-is.
+func resolveAutoRestart(p *bool) bool {
+	if p == nil {
+		return true
+	}
+	return *p
+}
+
 func validationDetail(err error) string {
 	msg := err.Error()
 	const prefix = "flowapi: validation failed: "
 	if len(msg) > len(prefix) && msg[:len(prefix)] == prefix {

Now add the tests to `http_test.go`:
← Edit internal/flowapi/http_test.go
Index: /root/astrate-mule/internal/flowapi/http_test.go
===================================================================
--- /root/astrate-mule/internal/flowapi/http_test.go
+++ /root/astrate-mule/internal/flowapi/http_test.go
@@ -38,8 +38,29 @@
 		}
 	}
 }
 
+func TestResolveAutoRestart(t *testing.T) {
+	boolPtr := func(v bool) *bool { return &v }
+	cases := []struct {
+		name string
+		p    *bool
+		want bool
+	}{
+		{"nil defaults to true", nil, true},
+		{"explicit true honored", boolPtr(true), true},
+		{"explicit false honored", boolPtr(false), false},
+	}
+	for _, tc := range cases {
+		t.Run(tc.name, func(t *testing.T) {
+			got := resolveAutoRestart(tc.p)
+			if got != tc.want {
+				t.Errorf("resolveAutoRestart(%v) = %v, want %v", tc.p, got, tc.want)
+			}
+		})
+	}
+}
+
 func TestValidationDetail(t *testing.T) {
 	cases := []struct {
 		name string
 		err  error

→ Read internal/flowapi/http.go [offset=139, limit=12]
Let me verify the edits are correct and run the gate:
→ Read internal/flowapi/http.go [offset=296, limit=20]
Now run the gate — build, vet, and the flowapi tests, plus gofmt check:

$ gofmt -l internal/flowapi/ && go vet ./internal/flowapi/ && go test ./internal/flowapi/
ok  	github.com/astrate-platform/astrate/internal/flowapi	0.037s

Gate passes for the package. Let me run the full container-free suite and whole-repo gofmt check as the gate requires:
$ gofmt -l . && go test ./...
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
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

All gates green. Task complete.

Done: Pin the "POST /flows without auto_restart defaults to true" wire rule by extracting the default into a pure `resolveAutoRestart` helper and testing missing→true, explicit true, explicit false.
Files: internal/flowapi/http.go (helper + callsite at startFlow http.go:141), internal/flowapi/http_test.go (TestResolveAutoRestart)
Verified: `gofmt -l .` -> pass (no output); `go vet ./internal/flowapi/` -> pass; `go test ./...` (without -race, container-free Pi gate) -> pass.
Unsure: nothing
Follow-ups: none
