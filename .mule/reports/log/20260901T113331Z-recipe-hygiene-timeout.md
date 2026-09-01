slug: recipe-hygiene
verdict: timeout
at:  5006f6e
ran: 2026-09-01T11:33:32Z on DietPi in 1637s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md
→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ which govulncheck 2>/dev/null && govulncheck ./... 2>/dev/null || echo "govulncheck not installed"
govulncheck not installed

$ which golangci-lint 2>/dev/null && golangci-lint run ./... 2>&1 | tail -30 || echo "golangci-lint not installed"
golangci-lint not installed

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ 2>/dev/null | head -30
internal/store/store.go:135:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5):
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},

Now running the hygiene checks.
$ go test ./... 2>&1 | rg -i 'skip|no test files' | head -20
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/internal/housekeeping	[no test files]
?   	github.com/astrate-platform/astrate/internal/swagger	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]

$ go list -m -u all 2>/dev/null | rg '\[' | head -20
cel.dev/expr v0.24.0 [v0.25.3]
cloud.google.com/go v0.121.6 [v0.123.0]
cloud.google.com/go/auth v0.16.4 [v0.23.2]
cloud.google.com/go/compute/metadata v0.8.0 [v0.9.0]
cloud.google.com/go/iam v1.5.2 [v1.13.0]
cloud.google.com/go/longrunning v0.6.7 [v1.2.0]
cloud.google.com/go/monitoring v1.24.2 [v1.30.0]
cloud.google.com/go/spanner v1.85.0 [v1.95.0]
cloud.google.com/go/storage v1.56.0 [v1.66.0]
github.com/99designs/keyring v1.2.1 [v1.2.2]
github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0 [v1.23.1]
github.com/Azure/azure-sdk-for-go/sdk/internal v1.1.2 [v1.12.0]
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0 [v1.8.0]
github.com/Azure/go-autorest/autorest/adal v0.9.16 (retracted) [v0.9.24] (deprecated)
github.com/Azure/go-autorest/autorest/date v0.3.0 [v0.3.1] (deprecated)
github.com/Azure/go-autorest/logger v0.2.1 [v0.2.2] (deprecated)
github.com/Azure/go-autorest/tracing v0.6.0 [v0.6.1] (deprecated)
github.com/ClickHouse/clickhouse-go v1.4.3 [v1.5.4]
github.com/DataDog/zstd v1.4.5 [v1.5.7]
github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp v1.5.3 [v1.6.0]

Let me investigate the findings more deeply. govulncheck and golangci-lint aren't installed, so I can't run those directly. Let me look at the dependency updates, the TODO comment, and the packages with no test files.

Let me check the go version and whether govulncheck can be run. Also examine the packages with no test files to see if they contain logic.
$ go version
go version go1.26.5 linux/arm64

$ ls internal/housekeeping/ internal/swagger/ 2>&1; echo "---migrations---"; ls migrations/*.go 2>/dev/null
internal/housekeeping/:
http.go
http_test.go
service.go

internal/swagger/:
swagger.go
---migrations---
migrations/migrations.go

housekeeping has no test file at all (http_test.go doesn't exist in the listing — wait, it does show http_test.go). Let me recheck.
$ ls -la internal/housekeeping/; echo "---"; rg -l 'func Test' internal/housekeeping/ 2>/dev/null || echo "no test funcs"
total 28
drwxr-xr-x  2 root root 4096 Aug 31 13:47 .
drwxr-xr-x 16 root root 4096 Jul 28 02:09 ..
-rw-r--r--  1 root root 3939 Aug 31 13:47 http.go
-rw-r--r--  1 root root 6897 Aug 31 13:47 http_test.go
-rw-r--r--  1 root root 5222 Jul 27 05:10 service.go
---
internal/housekeeping/http_test.go

$ rg -c 'func ' internal/housekeeping/*.go internal/swagger/*.go migrations/*.go 2>/dev/null
internal/housekeeping/http.go:9
internal/housekeeping/http_test.go:10
internal/housekeeping/service.go:7
internal/swagger/swagger.go:2

The `go test` output said "no test files" for internal/housekeeping but there IS an http_test.go with Test funcs. That's odd — the grep output may have been truncated/mismatched. Actually the output listed `internal/housekeeping [no test files]` — but there's clearly a test. Wait, maybe the test output line was from a parallel listing. Let me recheck the actual go test result for housekeeping.
$ go test ./internal/housekeeping/ 2>&1; echo "---swagger---"; go test ./internal/swagger/ 2>&1; echo "---migrations---"; go test ./migrations/ 2>&1; echo "---docs---"; go test ./docs/ 2>&1
?   	github.com/astrate-platform/astrate/internal/housekeeping	[no test files]
---swagger---
?   	github.com/astrate-platform/astrate/internal/swagger	[no test files]
---migrations---
?   	github.com/astrate-platform/astrate/migrations	[no test files]
---docs---
?   	github.com/astrate-platform/astrate/docs	[no test files]

Interesting - `go test` says housekeeping has no test files even though http_test.go exists. This means http_test.go must have a build constraint excluding it. Let me check.
→ Read internal/housekeeping/http_test.go
The housekeeping http_test.go is gated behind `//go:build integration`, so it doesn't run in the normal gate. That's a build-tagged integration test, not a "package with no test file". The recipe's concern about "package with no test file at all" is about the opposite — packages with logic but no test. Let me focus: `internal/swagger` has 2 funcs and no test file, and `migrations` has logic and no test. Let me check those.

Actually, let me reconsider the scope. The housekeeping test is an integration test (build tag), so it's not "no test file at all" in the sense of the recipe — it has tests, they just need the integration build tag. That's arguably a different matter.

Let me verify swagger and migrations logic.

$ rg -n '^func ' internal/swagger/swagger.go migrations/migrations.go
internal/swagger/swagger.go:16:func Mount(mux *http.ServeMux) {
internal/swagger/swagger.go:31:func Specs() []string {

→ Read internal/store/store.go [offset=125, limit=20]
The TODO in store.go is an extension point tied to an existing downstream feature (`s.hasToolkit` already exists, `datastreams.go` lttb path). It names a real missing behaviour (switching to lttb downsampling when the toolkit is present). This is a candidate task.

Let me now look at the migrations package - it has logic (migrations.go) and no test file. Let me check what's there.

Let me also check the golangci-lint situation — it's not installed. The recipe's lint check can't run. I have a note already in for-giulio.md about golangci-lint not installed. Let me check the .golangci.yml exists and focus.

Let me look at migrations.go logic and the swagger.go.
→ Read migrations/migrations.go
→ Read internal/swagger/swagger.go
Evaluating findings:

1. **Dependency updates** — `go list -m -u all` shows many newer versions but these are almost all transitive/indirect deps of the Go Cloud SDK. Per the recipe, I should not propose bumps "to stay current" — only for a fix this repo needs, and go.mod/go.sum are on the never-touch list anyway. These go to for-giulio.md only if there's a fix needed. Nothing indicates a needed fix. I'll skip unless there's a compelling reason.

2. **govulncheck / golangci-lint** — not installed. Can't run. The lint-not-installed is already documented in for-giulio.md (line 42-45). govulncheck is the highest-value check per the recipe (reachable vulnerabilities). Since I can't run it, I should note this. Let me check if I can get govulncheck quickly. Actually govulncheck needs network to download DB. Let me try; it might be feasible. The recipe explicitly says the vulnerability check is highest priority.

3. **TODO comment** — store.go:135 has a real missing behaviour (lttb downsampling fallback when timescaledb_toolkit present). This is a genuine task candidate. But wait — is this already covered elsewhere? It says `s.hasToolkit` is already recorded. The missing piece is switching the downsampling path. This is a real task.

4. **Packages with no test file** — `internal/swagger` (2 funcs, serving static files — mostly wiring/types, minimal logic) and `migrations` (just an embedded FS, a single embed var — types, no logic, actually explicitly a type/embed-only package). Per the recipe: "if it contains logic rather than types". `migrations` is pure embed/types — no logic, skip. `internal/swagger` — it serves files and walks a dir; there's a bit of logic in Specs() but it's thin wiring. Marginal; probably not worth a test task.

   `internal/housekeeping` has an integration-tagged test — not "no test file". Not applicable.

Let me try to run govulncheck since it's the highest priority item. Let me check for network/installation options.

$ ls /root/go/bin 2>/dev/null; echo "---"; ls ~/go/bin 2>/dev/null; echo "---PATH---"; echo $PATH; echo "---try install---"; timeout 90 go install golang.org/x/vuln/cmd/govulncheck@latest 2>&1 | tail -5
golangci-lint
---
golangci-lint
---PATH---
/root/.opencode/bin:/root/.opencode/bin:/usr/local/bin:/usr/bin:/bin
---try install---
go: downloading golang.org/x/vuln v1.7.0
go: downloading golang.org/x/telemetry v0.0.0-20260811182544-a038080d80e5
go: downloading golang.org/x/mod v0.39.0
go: downloading golang.org/x/tools v0.49.0
go: downloading golang.org/x/sync v0.22.0

govulncheck installed successfully. Let me run it. Also, golangci-lint IS installed at ~/go/bin (just not in PATH). Let me run both now.
$ /root/go/bin/golangci-lint run ./... 2>&1 | tail -40
internal/engine/capabilities.go:7:1: File is not properly formatted (goimports)
	"github.com/astrate-platform/astrate/internal/broker"
^
internal/flow/router.go:215:14: G115: integer overflow conversion rune -> uint64 (gosec)
		h ^= uint64(b)
		           ^
internal/flow/flow.go:14:6: exported: type name will be used as flow.FlowStatus by other packages, and that stutters; consider calling this Status (revive)
type FlowStatus uint8
     ^
internal/flow/flow.go:53:6: exported: type name will be used as flow.FlowConfig by other packages, and that stutters; consider calling this Config (revive)
type FlowConfig struct {
     ^
internal/flow/flow_test.go:118:28: empty-block: this block is empty, you can remove it (revive)
	if f.StoppedAt().IsZero() {
		// Not stopped yet; StoppedAt should be zero.
	}
internal/flow/message.go:21:2: exported: exported const TypeInteger should have comment (or a comment on this block) or be unexported (revive)
	TypeInteger DataType = iota
	^
internal/flow/message.go:77:6: exported: type name will be used as flow.FlowMessage by other packages, and that stutters; consider calling this Message (revive)
type FlowMessage struct {
     ^
internal/flow/pipeline.go:120:48: empty-block: this block is empty, you can remove it (revive)
		if inDeg[b.Name] == 0 || outDeg[b.Name] == 0 {
			// Re-check using original degrees (topo sort modified inDeg).
		}
internal/flow/router.go:180:26: unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _ (revive)
func (r *Router) runLane(ctx context.Context, l *lane) {
                         ^
internal/pairing/service_test.go:192:54: unused-parameter: parameter 'at' seems to be unused, consider removing or renaming it as _ (revive)
	svc.OnRegistered = func(realmName, deviceID string, at time.Time) {
	                                                    ^
10 issues:
* goimports: 1
* gosec: 1
* revive: 8

$ /root/go/bin/govulncheck ./... 2>&1 | tail -40
#1: internal/store/interfaces.go:366:22: store.loadEndpointIDs calls pgx.baseRows.Scan, which eventually calls xml.Unmarshal

Vulnerability #6: GO-2026-5972
    Enforce maximum recursion depth in encoding/asn1
  More info: https://pkg.go.dev/vuln/GO-2026-5972
  Standard library
    Found in: encoding/asn1@go1.26.5
    Fixed in: encoding/asn1@go1.26.6
    Example traces found:
      #1: internal/testutil/mqttclient.go:67:43: testutil.DeviceCSR calls x509.CreateCertificateRequest, which calls asn1.Unmarshal

Vulnerability #7: GO-2026-5970
    Infinite loop on invalid input in golang.org/x/text
  More info: https://pkg.go.dev/vuln/GO-2026-5970
  Module: golang.org/x/text
    Found in: golang.org/x/text@v0.38.0
    Fixed in: golang.org/x/text@v0.39.0
    Example traces found:
      #1: internal/store/notify.go:59:34: store.Listen calls pgx.ConnectConfig, which eventually calls norm.Form.Properties
      #2: internal/store/notify.go:59:34: store.Listen calls pgx.ConnectConfig, which eventually calls norm.Form.Span
      #3: internal/store/notify.go:59:34: store.Listen calls pgx.ConnectConfig, which eventually calls norm.Form.Transform

Vulnerability #8: GO-2026-5026
    Invoking failure to reject ASCII-only Punycode-encoded labels in
    golang.org/x/net/idna
  More info: https://pkg.go.dev/vuln/GO-2026-5026
  Standard library
    Found in: net/http@go1.26.5
    Fixed in: net/http@go1.26.6
    Example traces found:
      #1: internal/engine/forward/http.go:113:26: forward.HTTP.Forward calls http.Client.Do
      #2: internal/engine/triggers/actions.go:307:15: triggers.Executor.Close calls sync.Once.Do, which eventually calls http.Get
      #3: internal/engine/triggers/actions.go:307:15: triggers.Executor.Close calls sync.Once.Do, which eventually calls http.Transport.CloseIdleConnections
      #4: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls http.Transport.RoundTrip

Your code is affected by 8 vulnerabilities from 2 modules and the Go standard library.
This scan also found 3 vulnerabilities in packages you import and 3
vulnerabilities in modules you require, but your code doesn't appear to call
these vulnerabilities.
Use '-show verbose' for more details.

This is very useful. Now I have concrete findings. Let me check whether these lint findings are already excluded in the config, and check the golangci config. Let me look at what's there.

Also, govulncheck found 8 reachable vulnerabilities. Let me look at each carefully. The recipe's highest priority: **A vulnerability in a dependency `govulncheck` says is actually reachable.**

Key ones to assess:
- GO-2026-5972 (asn1 recursion depth) — stdlib, reachable in testutil (test-only). 
- GO-2026-5970 (x/text infinite loop) — reachable via store.Listen/pgx. **This is a real one** — x/text is a transitively used module, reachable in production code (store.Listen).
- GO-2026-5026 (net/http idna) — reachable via forward.HTTP.Forward and triggers http.Get. **This is a real production path.**

These stdlib ones (asn1, net/http) are fixed in Go 1.26.6 — a Go toolchain bump, not a go.mod change. That's a decision for Giulio.

Let me see the full govulncheck output including the top vulnerabilities (there were 8, I only saw #6-8).

Let me re-run to capture all of them and understand the full picture, especially which are reachable in production vs test-only.
$ /root/go/bin/govulncheck ./... 2>&1 | head -80
(no output)

<shell_metadata>
shell tool terminated command after exceeding timeout 300000 ms. If this command is expected to take longer and is not waiting for interactive input, retry with a larger timeout value in milliseconds.
</shell_metadata>

