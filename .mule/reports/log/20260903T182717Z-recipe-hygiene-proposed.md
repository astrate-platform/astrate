slug: recipe-hygiene
verdict: proposed
at:  6b7ade3
ran: 2026-09-03T18:27:17Z on DietPi in 184s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md

# Todos
[•] Run dependency update check (go list -m -u all)
[ ] Run govulncheck
[ ] Run golangci-lint
[ ] Scan for TODO/FIXME/XXX/HACK comments
[ ] Check for skipped tests / missing test files
[ ] Write findings to .mule/todo.md and/or .mule/for-giulio.md

$ which govulncheck 2>/dev/null && govulncheck ./... 2>&1 || echo "govulncheck not installed"
govulncheck not installed

$ which golangci-lint 2>/dev/null && golangci-lint run ./... 2>&1 | tail -30 || echo "golangci-lint not installed"
golangci-lint not installed

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ 2>/dev/null | head -30
internal/store/store.go:135:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5):
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},

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
cloud.google.com/go/storage v1.56.0 [v1.67.0]
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

# Todos
[✓] Run dependency update check (go list -m -u all)
[✓] Run govulncheck
[✓] Run golangci-lint
[•] Scan for TODO/FIXME/XXX/HACK comments
[ ] Check for skipped tests / missing test files
[ ] Write findings to .mule/todo.md and/or .mule/for-giulio.md

→ Read internal/housekeeping
→ Read internal/swagger
→ Read internal/store/store.go [offset=125, limit=25]
→ Read .mule/todo.md
→ Read .mule/for-giulio.md
→ Read internal/swagger/swagger.go
→ Read internal/housekeeping/http_test.go

# Todos
[✓] Run dependency update check (go list -m -u all)
[✓] Run govulncheck
[✓] Run golangci-lint
[✓] Scan for TODO/FIXME/XXX/HACK comments
[✓] Check for skipped tests / missing test files
[•] Write findings to .mule/todo.md and/or .mule/for-giulio.md

Now let me check the `service.go` in housekeeping to understand the logic there:
→ Read internal/housekeeping/service.go
Here are my findings. Let me write the deliverables.

**Findings summary:**

- **govulncheck**: not installed (already noted in for-giulio.md)
- **golangci-lint**: not installed (already noted in for-giulio.md)
- **TODO at `internal/store/store.go:135`**: real missing behavior — Downsample should use toolkit lttb() when available; probe infrastructure exists but the switch isn't implemented
- **`internal/housekeeping`**: has only integration tests (`//go:build integration`), so unit-testable validation logic in `CreateRealm` (blank name/key, negative reg limit) goes untested on the Pi
- **Dependency updates**: many newer versions available, but all are transitive/cloud SDKs — no fix this repo needs identified; already noted in for-giulio.md
→ Read .mule/todo.md [offset=120]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -122,4 +122,6 @@
 - [!] lint-flow-stutter: in `internal/flow`, rename exported `FlowStatus`/`FlowConfig`/`FlowMessage` to `Status`/`Config`/`Message` to satisfy revive's stutter rule; touches flow.go, message.go, block.go, graph.go, router.go and their tests — a public API rename, so no wire/JSON change (the Type names are unchanged), it is purely Go identifier renaming. [auto] — BLOCKED: gates failed
 - [!] serverdata-aftercommit-test: in `internal/engine/serverdata_test.go`, add a test that wires an `afterCommit` observer on a `newWiredRig` engine, calls `PublishServerValue` for a server-owned property, and asserts whether the observer fires. The batch path fires `afterCommit` for every committed op (batch.go:168), but `PublishServerValue` bypasses the batch entirely (serverdata.go:99-162) — the current tests never check this. The test documents the intended behaviour (exclusion or parity with the data path) and fails if the contract changes. [auto] — BLOCKED: gates failed
 - [!] consumer-properties-plaintext-fallback: in `internal/engine/control.go` `sendConsumerProperties` (line 204), when `devices.peek` returns nil (device not cached after eviction), the code defaults to zlib compression — but a device that declared `purge_properties_compression_format: plaintext` via the capabilities topic cannot decompress zlib. Fix by falling back to a store read of the device's capability state (or persisting the capability), and add a test in `internal/engine/capabilities_test.go` that evicts the cache entry for a plaintext device and asserts the purge payload is uncompressed. [auto] — BLOCKED: gates failed
 - [!] object-path-ok-comment: in `internal/engine/data.go` `objectPathOK` (line 330), the loop over `ci.ObjectLeaves` returns after checking the first entry. Add a comment explaining why one leaf suffices (upstream validation guarantees uniform prefix depth), so future readers do not mistake the early return for a bug. [auto] — BLOCKED: gates failed
+- [ ] housekeeping-unit-tests: in `internal/housekeeping`, add a unit test file (no integration build tag) that exercises `CreateRealm` validation paths — blank name returns ErrValidation, blank JWT key returns ErrValidation, negative registration limit returns ErrValidation — using a mock store and sealer. The package currently has only integration tests gated behind `//go:build integration`, so these paths are never tested on the Pi. [auto]
+- [ ] store-downsample-lttb: in `internal/store`, implement the TODO at store.go:135 — when `s.hasToolkit` is true, switch `Downsample` from `time_bucket+avg` to toolkit `lttb()` downsampling. The probe already records availability; the time_bucket path in datastreams.go is the always-working default. Add a test that exercises both code paths (mocking `hasToolkit`). [auto]

# Todos
[✓] Run dependency update check (go list -m -u all)
[✓] Run govulncheck
[✓] Run golangci-lint
[✓] Scan for TODO/FIXME/XXX/HACK comments
[✓] Check for skipped tests / missing test files
[✓] Write findings to .mule/todo.md and/or .mule/for-giulio.md

→ Read .mule/todo.md [offset=123]
Done: hygiene recipe run — 2 task lines appended to .mule/todo.md
Files: .mule/todo.md (lines 126–127)
Verified: all 5 checks attempted (go list, govulncheck, golangci-lint, rg TODO, go test) → govulncheck and golangci-lint not installed; dep scan, TODO scan, test scan completed
Unsure: whether store.go:135 TODO's Downsample lttb switch is a v3.0 priority or deferred — referenced in docs/ROADMAP §0.1 but not in milestones.md
Follow-ups: 2 lines appended
- `housekeeping-unit-tests`: CreateRealm validation paths have no unit tests (only integration-tagged tests)
- `store-downsample-lttb`: implement the TODO at store.go:135 (timescaledb_toolkit lttb downsampling)
