slug: recipe-hygiene
verdict: timeout
at:  230dd11
ran: 2026-10-06T14:05:37Z on DietPi in 11179s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md

$ which golangci-lint || ls /Users/atsetilam/go/bin/golangci-lint 2>/dev/null || echo "golangci-lint not found"
golangci-lint not found

$ go list -m -u all 2>/dev/null | rg '\[' | head -20; echo "---GOVULN---"; which govulncheck || echo "govulncheck not installed"
cel.dev/expr v0.24.0 [v0.25.3]
cloud.google.com/go v0.121.6 [v0.123.0]
cloud.google.com/go/auth v0.16.4 [v0.24.0]
cloud.google.com/go/auth/oauth2adapt v0.2.8 [v0.3.0]
cloud.google.com/go/compute/metadata v0.8.0 [v0.10.0]
cloud.google.com/go/iam v1.5.2 [v1.14.0]
cloud.google.com/go/longrunning v0.6.7 [v1.3.0]
cloud.google.com/go/monitoring v1.24.2 [v1.31.0]
cloud.google.com/go/spanner v1.85.0 [v1.95.1]
cloud.google.com/go/storage v1.56.0 [v1.69.0]
github.com/99designs/keyring v1.2.1 [v1.2.2]
github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0 [v1.23.2]
github.com/Azure/azure-sdk-for-go/sdk/internal v1.1.2 [v1.13.0]
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0 [v1.8.2]
github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c [v0.0.0-20260917205352-e937bb47801a]
github.com/Azure/go-autorest/autorest/adal v0.9.16 (retracted) [v0.9.24] (deprecated)
github.com/Azure/go-autorest/autorest/date v0.3.0 [v0.3.1] (deprecated)
github.com/Azure/go-autorest/logger v0.2.1 [v0.2.2] (deprecated)
github.com/Azure/go-autorest/tracing v0.6.0 [v0.6.1] (deprecated)
github.com/ClickHouse/clickhouse-go v1.4.3 [v1.5.4]
---GOVULN---
govulncheck not installed

$ ls "$(go env GOPATH)/bin" 2>/dev/null; echo "---"; go env GOFLAGS GOPATH
goimports
golangci-lint
govulncheck
---

/root/go

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ 2>/dev/null | head -30; echo "===EXIT $?==="
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):
===EXIT 0===

$ PATH="$PATH:$(go env GOPATH)/bin" golangci-lint run ./... 2>&1 | tail -40; echo "===EXIT==="
0 issues.
===EXIT===

