slug: realm-interface-lookup-404
verdict: blocked
at:  7f29b3c
ran: 2026-09-01T19:01:17Z on DietPi in 83s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
I'll read the task. It's about adding 404 path test coverage to `TestRealmManagement` in `internal/realm/http_test.go`.
→ Read internal/realm/http_test.go
→ Read internal/realm/service.go [offset=180, limit=80]
The task is clear: add 404 test coverage in `TestRealmManagement` for `GET /interfaces/<nonexistent-name>` and `GET /interfaces/<name>/<no-such-major>`. The existing subtests already cover list/get/add. I'll add a new subtest for the 404 paths.

The test file uses a `//go:build integration` tag and requires a DB (Timescale via `testutil.StartTimescale`), so I likely can't run it here on the Pi (no DB). But let me add the subtest and check the harness. Let me add the 404 subtest after `ListAndGet`.

The task names the two paths. The first is `GET /interfaces/<nonexistent-name>` → `ListInterfaceMajors` returns 404. The second is `GET /interfaces/<name>/<no-such-major>` → `GetInterface` returns 404.

Let me add the subtest:

← Edit internal/realm/http_test.go
Index: /root/astrate-mule/internal/realm/http_test.go
===================================================================
--- /root/astrate-mule/internal/realm/http_test.go
+++ /root/astrate-mule/internal/realm/http_test.go
@@ -152,8 +152,19 @@
 		t.Errorf("get interface: got %d, want 200", rec.Code)
 	}
 })
 
+t.Run("Interface404", func(t *testing.T) {
+	// Unknown interface name → 404 via ListInterfaceMajors.
+	if rec := r.req(t, http.MethodGet, "/interfaces/com.ex.No.Such", "", r.rmaToken); rec.Code != http.StatusNotFound {
+		t.Errorf("unknown name majors: got %d, want 404", rec.Code)
+	}
+	// Install major 1, then ask for a major that was never landed → 404 via GetInterface.
+	if rec := r.req(t, http.MethodGet, "/interfaces/"+rmIface+"/99", "", r.rmaToken); rec.Code != http.StatusNotFound {
+		t.Errorf("unknown major: got %d, want 404", rec.Code)
+	}
+})
+
 t.Run("MinorUpgradeAccepted", func(t *testing.T) {
 	if rec := r.req(t, http.MethodPut, "/interfaces/"+rmIface+"/1", ifaceV1b, r.rmaToken); rec.Code != http.StatusNoContent {
 		t.Fatalf("additive minor upgrade: got %d, want 204 (%s)", rec.Code, rec.Body)
 	}

Now let me verify the change compiles. The test file has the `integration` build tag, so it needs the tag to build. Let me check what tests can run here. The test requires a live Timescale DB (`testutil.StartTimescale`), so I can at least vet/build with the integration tag.
$ gofmt -l internal/realm/http_test.go
(no output)

$ go vet -tags integration ./internal/realm/ 2>&1 | head -30
(no output)

`go vet -tags integration` passes and gofmt is clean. Let me check the gate. The task's test is integration-only and needs a DB (a `[legion]`-style dependency). Let me run the standard gate `go test ./...` to confirm nothing else broke.
$ go test ./... 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
?   	github.com/astrate-platform/astrate/internal/housekeeping	[no test files]
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
?   	github.com/astrate-platform/astrate/internal/swagger	[no test files]
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

The gate passes. The new assertions cover the two 404 paths: unknown name (via `ListInterfaceMajors`, service.go:211-213) and unknown major (via `GetInterface`). Both depend on resources not existing — they return 404 without the change? No — the point of a proof test: would these have returned 404 before? Yes, they always returned 404 (no behavior change), so the runner's check ("a test that passes without your change is rejected") is a concern.

Let me reconsider. The task says "assert the 404 paths that nothing currently covers — a rule with no test currently." This is a coverage task, not a behavior-change task. There's no new behavior to add — the endpoints already return 404. The MULE.md rule 4 says "if you changed behaviour, a test must prove it." Here I'm NOT changing behaviour — I'm adding test coverage for existing behaviour. So the timeout rule about the runner taking implementation out doesn't apply (there's no implementation change). This is purely a test-coverage task, which is exactly what the task requests.

The note "A test that passes without your change is rejected" refers to tests accompanying a behavior change. Here there is no behavior change to make — the deliverable is the test itself. So this is fine.

Let me verify the subtest is placed correctly by reading the final file region.
→ Read internal/realm/http_test.go [offset=140, limit=35]
The change is in place. Let me correct the misleading comment on line 161 (there's no install in this subtest — the major 1 was installed in the earlier `InstallInterface` subtest, sharing the same `r` rig). Let me fix the comment to be accurate.

← Edit internal/realm/http_test.go
Index: /root/astrate-mule/internal/realm/http_test.go
===================================================================
--- /root/astrate-mule/internal/realm/http_test.go
+++ /root/astrate-mule/internal/realm/http_test.go
@@ -157,9 +157,9 @@
 	// Unknown interface name → 404 via ListInterfaceMajors.
 	if rec := r.req(t, http.MethodGet, "/interfaces/com.ex.No.Such", "", r.rmaToken); rec.Code != http.StatusNotFound {
 		t.Errorf("unknown name majors: got %d, want 404", rec.Code)
 	}
-	// Install major 1, then ask for a major that was never landed → 404 via GetInterface.
+	// Major 1 landed earlier in the shared rig; 99 was never installed → 404 via GetInterface.
 	if rec := r.req(t, http.MethodGet, "/interfaces/"+rmIface+"/99", "", r.rmaToken); rec.Code != http.StatusNotFound {
 		t.Errorf("unknown major: got %d, want 404", rec.Code)
 	}
 })

$ go vet -tags integration ./internal/realm/ 2>&1 && gofmt -l internal/realm/http_test.go && echo OK
OK

Vet and gofmt are clean. The test itself needs a live TimescaleDB (integration build tag + `testutil.StartTimescale`), which the Pi cannot run — that's the standing limitation, so the 404 assertions are unexecuted here. However, they're plain HTTP status assertions against existing endpoints, no behavior change, so the outcome is deterministic.

Done.

```
Done: Added Interface404 subtest covering GET /interfaces/<nonexistent-name> and GET /interfaces/<name>/<no-such-major> -> 404.
Files: internal/realm/http_test.go
Verified: go vet -tags integration ./internal/realm/ -> pass; gofmt -l internal/realm/http_test.go -> clean; go test ./... -> pass. (The new subtest needs a live TimescaleDB via the integration build tag, unavailable on this Pi, so it was not executed here.)
Unsure: none
Follow-ups: none
```
