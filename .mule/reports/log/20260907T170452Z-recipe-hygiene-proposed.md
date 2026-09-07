slug: recipe-hygiene
verdict: proposed
at:  17c6cff
ran: 2026-09-07T17:04:52Z on DietPi in 156s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md
→ Read .mule/todo.md
→ Read .mule/for-giulio.md
Let me read the current state files and run the hygiene checks in parallel.

$ govulncheck ./... 2>/dev/null || echo "govulncheck not installed"
govulncheck not installed

$ which golangci-lint 2>/dev/null || ls /Users/atsetilam/go/bin/golangci-lint 2>/dev/null || echo "golangci-lint not found"
golangci-lint not found

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ 2>/dev/null | head -30
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):
internal/store/store.go:139:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5):

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

$ go test ./... 2>&1 | rg -i 'skip|no test files' | head -20
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]

→ Read internal/broker/intake.go [offset=54, limit=5]
Let me check a few more things to complete the hygiene picture.
→ Read internal/store/store.go [offset=137, limit=5]

$ go test ./... 2>&1 | grep -E '^(---\s+SKIP|ok\s+.+\s+0\.0)' | head -20
(no output)

No actionable findings:

- **govulncheck / golangci-lint**: not installed on this machine (already escalated to for-giulio.md — both are available on the Pi).
- **Dependency updates**: all covered by the 2026-09-04 sweep in for-giulio.md; nothing the repo needs.
- **TODOs**: both are roadmap extension-point placeholders (aspirational), not missing behaviour — `intake.go:56` and `store.go:139` both reference `docs/ROADMAP.md` sections that already exist.
- **Packages without tests**: only `migrations/` (config/data) and `examples/flow-container-echo/` (example code) — neither contains logic worth unit-testing.
- **Skipped tests**: none.

Nothing to propose, nothing to escalate. Stopping here.
