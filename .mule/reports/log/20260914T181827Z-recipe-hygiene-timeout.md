slug: recipe-hygiene
verdict: timeout
at:  f5c6010
ran: 2026-09-14T18:18:28Z on DietPi in 1208s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md
I've read both files. Running the hygiene recipe checks now.

$ govulncheck ./... 2>/dev/null || echo "govulncheck not installed"
govulncheck not installed

$ go vet ./... 2>&1 | tail -5; echo "---VET_DONE---"
---VET_DONE---

$ go list -m -u all 2>/dev/null | rg '\[' | head -20 || echo "NO_UPDATES"
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

$ which golangci-lint 2>/dev/null; ls /Users/atsetilam/go/bin/golangci-lint 2>/dev/null; golangci-lint run ./... 2>&1 | tail -30 || echo "golangci-lint unavailable"
/bin/sh: 1: golangci-lint: not found

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ | head -30
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},

$ go test ./... 2>&1 | rg -i 'skip|no test files' | head -20
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]

→ Read internal/broker/intake.go [offset=45, limit=30]
lint unavailable (macOS path), no skipped tests. Let me inspect the one real TODO comment and the no-test-file packages.

$ rg -c 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ 2>/dev/null; echo '---filelisting---'; rg -l 'TODO|FIXME|HACK' internal/ pkg/ cmd/ 2>/dev/null
internal/flow/blocks/httpblocks.go:1
internal/flow/blocks/astartesource/source.go:1
internal/auth/claims_test.go:2
internal/broker/intake.go:1
---filelisting---
internal/flow/blocks/httpblocks.go
internal/flow/blocks/astartesource/source.go
internal/broker/intake.go

$ go test ./... 2>&1 | rg -i 'skip|no test files' ; echo "---"; ls examples/flow-container-echo/ migrations/ 2>/dev/null
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]
---
examples/flow-container-echo/:
Dockerfile
README.md
main.go

migrations/:
000001_extensions.down.sql
000001_extensions.up.sql
000002_metadata.down.sql
000002_metadata.up.sql
000003_properties.down.sql
000003_properties.up.sql
000004_datastreams.down.sql
000004_datastreams.up.sql
000005_jobs.down.sql
000005_jobs.up.sql
000006_policies.down.sql
000006_policies.up.sql
000007_toolkit.down.sql
000007_toolkit.up.sql
000008_pipelines.down.sql
000008_pipelines.up.sql
000009_flows.down.sql
000009_flows.up.sql
000010_flow_failed_block.down.sql
000010_flow_failed_block.up.sql
000011_realm_retention.down.sql
000011_realm_retention.up.sql
000012_user_blocks.down.sql
000012_user_blocks.up.sql
migrations.go

The only real TODO (`internal/broker/intake.go:56`) is an explicit extension-point marker pointing at frozen ROADMAP/DESIGN decisions — aspirational, per the recipe I ignore it. `govulncheck` is the highest-value check, so let me attempt it via `go run @latest` rather than skip it. Also confirming the no-test-file dirs contain no real logic.
→ Read examples/flow-container-echo/main.go

$ go run golang.org/x/vuln/cmd/govulncheck@latest ./... 2>&1 | tail -40
signal: killed

govulncheck was killed (OOM on this Pi). Let me verify the full `go test` output for any skips.
$ go test ./... 2>&1 | rg -i 'SKIP' | head -20; echo "TEST_EXIT=$?"
TEST_EXIT=0

→ Read migrations/migrations.go
Now let me check the no-test-file packages to see if any contain real logic (per the recipe), and the `.golangci.yml` to confirm what's excluded before deciding on lint.

$ ls docs/*.go 2>/dev/null; echo "---"; cat .golangci.yml 2>/dev/null | head -60
docs/embed.go
---
# golangci-lint v2 configuration (pinned v2.12.2 — see Makefile GOLANGCI_LINT_VERSION).
# Linter set per docs/ROADMAP.md §1.2 (file 0.4): govet, staticcheck, errcheck, gosec, revive.
version: "2"

run:
  build-tags:
    # Lint everything, including T2/T3 suites that are tag-gated out of `go test ./...`.
    - integration
    - e2e

linters:
  default: none
  enable:
    - errcheck
    - gosec
    - govet
    - revive
    - staticcheck
  settings:
    gosec:
      excludes:
        # G104 duplicates errcheck.
        - G104
  exclusions:
    generated: lax
    rules:
      # Tests: hardening lints and strict error checking add noise, not safety.
      - path: _test\.go
        linters:
          - gosec
      - path: _test\.go
        linters:
          - errcheck
      # Test utilities are test-scoped code shipped outside _test.go files.
      - path: internal/testutil/
        linters:
          - gosec

formatters:
  enable:
    - gofmt
    - goimports
  settings:
    goimports:
      local-prefixes:
        - github.com/astrate-platform/astrate

$ ls ~/go/bin 2>/dev/null; ls /usr/local/go/bin 2>/dev/null; ls $(go env GOPATH)/bin 2>/dev/null; echo "---"; cat docs/embed.go 2>/dev/null | head -20
goimports
golangci-lint
govulncheck
goimports
golangci-lint
govulncheck
---
// Package docs embeds the Swagger UI static files and OpenAPI YAML specs so
// they are served directly from the astrate binary without external files.
package docs

import "embed"

// SwaggerUI holds the static Swagger UI files (index.html, CSS, JS refs).
//
//go:embed swagger-ui/*
var SwaggerUI embed.FS

// APIYAML holds the OpenAPI 3.0 YAML specifications for all five API surfaces.
//
//go:embed api/*.yaml
var APIYAML embed.FS

golangci-lint and govulncheck are installed at `$(go env GOPATH)/bin`. Running them properly.
