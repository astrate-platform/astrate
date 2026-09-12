slug: recipe-hygiene
verdict: timeout
at:  fb9d8bf
ran: 2026-09-12T17:08:52Z on DietPi in 1293s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md

$ which govulncheck 2>/dev/null && govulncheck ./... 2>/dev/null || echo "govulncheck not installed"
govulncheck not installed

$ golangci-lint run ./... 2>&1 | tail -30
/bin/sh: 1: golangci-lint: not found

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ | head -30
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/store/store.go:139:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5):
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},

$ go test ./... 2>&1 | rg -i 'skip|no test files' | head -20
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]

$ go list -m -u all 2>/dev/null | rg '\[' | head -20
cel.dev/expr v0.24.0 [v0.25.3]
cloud.google.com/go v0.121.6 [v0.123.0]
cloud.google.com/go/auth v0.16.4 [v0.23.2]
cloud.google.com/go/compute/metadata v0.8.0 [v0.9.0]
cloud.google.com/go/iam v1.5.2 [v1.13.0]
cloud.google.com/go/longrunning v0.6.7 [v1.2.0]
cloud.google.com/go/monitoring v1.24.2 [v1.30.0]
cloud.google.com/go/spanner v1.85.0 [v1.95.1]
cloud.google.com/go/storage v1.56.0 [v1.67.1]
github.com/99designs/keyring v1.2.1 [v1.2.2]
github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0 [v1.23.1]
github.com/Azure/azure-sdk-for-go/sdk/internal v1.1.2 [v1.12.0]
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0 [v1.8.1]
github.com/Azure/go-autorest/autorest/adal v0.9.16 (retracted) [v0.9.24] (deprecated)
github.com/Azure/go-autorest/autorest/date v0.3.0 [v0.3.1] (deprecated)
github.com/Azure/go-autorest/logger v0.2.1 [v0.2.2] (deprecated)
github.com/Azure/go-autorest/tracing v0.6.0 [v0.6.1] (deprecated)
github.com/ClickHouse/clickhouse-go v1.4.3 [v1.5.4]
github.com/DataDog/zstd v1.4.5 [v1.5.7]
github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp v1.5.3 [v1.6.0]

→ Read internal/flow/blocks/httpblocks.go [offset=140, limit=20]
→ Read internal/store/store.go [offset=135, limit=15]
→ Read internal/broker/intake.go [offset=50, limit=15]
→ Read docs
The TODOs in `internal/flow/blocks/httpblocks.go:147` and `internal/flow/blocks/astartesource/source.go:56` use `context.TODO()` as a placeholder — that's a known Go idiom for "figure out the right context later," not a real TODO for missing behaviour. The two in `store.go:139` and `intake.go:56` are explicit `TODO(extension point, ...)` markers tied to the roadmap — those are planned extension points, not things to act on now.

Let me check the packages with no test files to see if any contain meaningful logic.
→ Read examples/flow-container-echo
→ Read migrations
None of the untested packages contain meaningful Go logic worth testing (`docs/`, `examples/`, `migrations/`).

Now let me check the current todo list.
→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ ls /root/go/bin 2>/dev/null; echo "---"; ls ~/go/bin 2>/dev/null; echo "---PATH---"; echo "$PATH"
goimports
golangci-lint
govulncheck
---
goimports
golangci-lint
govulncheck
---PATH---
/root/.opencode/bin:/root/.opencode/bin:/usr/local/bin:/usr/bin:/bin

$ go version
go version go1.26.7 linux/arm64

Let me check for the toolchain binaries in the expected locations referenced in the repos' own state files.
Tools are in `/root/go/bin`, just not on PATH. Running the three binary-dependent checks with full paths.
$ /root/go/bin/golangci-lint run ./... 2>&1 | tail -30
0 issues.

