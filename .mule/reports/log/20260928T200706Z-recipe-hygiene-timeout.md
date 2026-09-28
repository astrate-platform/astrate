slug: recipe-hygiene
verdict: timeout
at:  484ec91
ran: 2026-09-28T20:07:07Z on DietPi in 4725s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md

# Todos
[•] Read .mule/todo.md, for-giulio.md, milestones.md for context/duplication
[ ] Run hygiene checks (deps, govulncheck, golangci-lint, TODO sweep, skipped tests)
[ ] Verify each candidate finding against the code before proposing
[ ] Append task lines to .mule/todo.md (and for-giulio.md if warranted)
[ ] Write evidence file if the recipe names one; final report

→ Read .mule/todo.md
→ Read .mule/for-giulio.md

# Todos
[✓] Read .mule/todo.md, for-giulio.md, milestones.md for context/duplication
[•] Run hygiene checks (deps, govulncheck, golangci-lint, TODO sweep, skipped tests)
[ ] Verify each candidate finding against the code before proposing
[ ] Append task lines to .mule/todo.md (and for-giulio.md if warranted)
[ ] Write evidence file if the recipe names one; final report

$ which govulncheck golangci-lint 2>&1; ls /Users/atsetilam/go/bin/ 2>&1 | head; ls ~/go/bin 2>&1 | head
ls: cannot access '/Users/atsetilam/go/bin/': No such file or directory
goimports
golangci-lint
govulncheck

$ go list -m -u all 2>/dev/null | rg '\[' | head -30
cel.dev/expr v0.24.0 [v0.25.3]
cloud.google.com/go v0.121.6 [v0.123.0]
cloud.google.com/go/auth v0.16.4 [v0.24.0]
cloud.google.com/go/auth/oauth2adapt v0.2.8 [v0.3.0]
cloud.google.com/go/compute/metadata v0.8.0 [v0.10.0]
cloud.google.com/go/iam v1.5.2 [v1.14.0]
cloud.google.com/go/longrunning v0.6.7 [v1.3.0]
cloud.google.com/go/monitoring v1.24.2 [v1.31.0]
cloud.google.com/go/spanner v1.85.0 [v1.95.1]
cloud.google.com/go/storage v1.56.0 [v1.68.0]
github.com/99designs/keyring v1.2.1 [v1.2.2]
github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0 [v1.23.1]
github.com/Azure/azure-sdk-for-go/sdk/internal v1.1.2 [v1.12.0]
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0 [v1.8.1]
github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c [v0.0.0-20260917205352-e937bb47801a]
github.com/Azure/go-autorest/autorest/adal v0.9.16 (retracted) [v0.9.24] (deprecated)
github.com/Azure/go-autorest/autorest/date v0.3.0 [v0.3.1] (deprecated)
github.com/Azure/go-autorest/logger v0.2.1 [v0.2.2] (deprecated)
github.com/Azure/go-autorest/tracing v0.6.0 [v0.6.1] (deprecated)
github.com/ClickHouse/clickhouse-go v1.4.3 [v1.5.4]
github.com/DataDog/zstd v1.4.5 [v1.5.7]
github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp v1.5.3 [v1.6.0]
github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.27.0 [v1.38.0] (deprecated)
github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.53.0 [v0.62.0] (deprecated)
github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.53.0 [v0.62.0]
github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 [v0.0.0-20240927000941-0f3dac36c52b]
github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a [v0.0.0-20230218143504-906a9b012302]
github.com/alicebob/miniredis/v2 v2.23.0 [v2.39.0]
github.com/andybalholm/brotli v1.0.4 [v1.2.5]
github.com/apache/thrift v0.16.0 [v0.24.0]

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ 2>/dev/null | head -40
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},

