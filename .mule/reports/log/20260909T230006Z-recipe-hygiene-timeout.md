slug: recipe-hygiene
verdict: timeout
at:  979872f
ran: 2026-09-09T23:00:07Z on DietPi in 8629s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md
→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ which govulncheck 2>/dev/null && govulncheck ./... 2>&1 | tail -40 || echo "govulncheck not installed"
govulncheck not installed

$ which golangci-lint 2>/dev/null && golangci-lint run ./... 2>&1 | tail -30 || echo "golangci-lint not installed"
golangci-lint not installed

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ | head -30
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/store/store.go:139:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5):
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):

$ go list -m -u all 2>/dev/null | rg '\[' | head -20
cel.dev/expr v0.24.0 [v0.25.3]
cloud.google.com/go v0.121.6 [v0.123.0]
cloud.google.com/go/auth v0.16.4 [v0.23.2]
cloud.google.com/go/compute/metadata v0.8.0 [v0.9.0]
cloud.google.com/go/iam v1.5.2 [v1.13.0]
cloud.google.com/go/longrunning v0.6.7 [v1.2.0]
cloud.google.com/go/monitoring v1.24.2 [v1.30.0]
cloud.google.com/go/spanner v1.85.0 [v1.95.0]
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

$ go test ./... 2>&1 | rg -i 'skip|no test files' | head -20
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]

$ ls -la /root/go/bin/ 2>/dev/null; ls ~/go/bin/ 2>/dev/null; echo "---"; go env GOPATH GOBIN
total 85260
drwxr-xr-x 2 root root     4096 Sep  2 21:44 .
drwxr-xr-x 4 root root     4096 Jul 28 23:52 ..
-rwxr-xr-x 1 root root  7668042 Sep  2 21:44 goimports
-rwxr-xr-x 1 root root 65622584 Jul 28 23:52 golangci-lint
-rwxr-xr-x 1 root root 13996581 Sep  1 13:09 govulncheck
goimports
golangci-lint
govulncheck
---
/root/go

$ rg -n 't\.Skip|\.Skipf|markSkipped|Skip\(' internal/ pkg/ cmd/ | head -30
cmd/astrate/main_test.go:47:	t.Skip("the boot suite needs a database: run `make up` or set ASTRATE_TEST_DSN")
cmd/astrate/forward_test.go:68:		t.Skip("built with -tags nats; see forward_nats_test.go")
internal/store/datastreams_test.go:420:			t.Skip("timescaledb_toolkit not installed; lttb path unavailable")
internal/pairing/ca/ca_test.go:293:		t.Skip("10k issuance draw skipped in -short mode")
internal/pairing/http_test.go:629:		t.Skip("openssl not found in PATH; chain already verified with crypto/x509")

$ go list -m -u all 2>/dev/null | rg -v 'cloud.google|azure|aws|gcp|opentelemetry|prometheus|otel|zstd|kafka|clickhouse|data-dog|datadog|grpc|spanner|bigquery|snowflake|microsoft' | rg '\[' | head -40
cel.dev/expr v0.24.0 [v0.25.3]
github.com/99designs/keyring v1.2.1 [v1.2.2]
github.com/Azure/go-autorest/autorest/adal v0.9.16 (retracted) [v0.9.24] (deprecated)
github.com/Azure/go-autorest/autorest/date v0.3.0 [v0.3.1] (deprecated)
github.com/Azure/go-autorest/logger v0.2.1 [v0.2.2] (deprecated)
github.com/Azure/go-autorest/tracing v0.6.0 [v0.6.1] (deprecated)
github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 [v0.0.0-20240927000941-0f3dac36c52b]
github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a [v0.0.0-20230218143504-906a9b012302]
github.com/alicebob/miniredis/v2 v2.23.0 [v2.39.0]
github.com/andybalholm/brotli v1.0.4 [v1.2.3]
github.com/apache/thrift v0.16.0 [v0.24.0]
github.com/cbroglie/mustache v1.4.0 [v1.4.2]
github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58 [v0.0.0-20240916140612-caecf3c00c06]
github.com/cncf/xds/go v0.0.0-20250501225837-2ac532fd4443 [v0.0.0-20260202195803-dba9d589def2]
github.com/cockroachdb/cockroach-go/v2 v2.1.1 [v2.4.3]
github.com/cockroachdb/errors v1.11.1 [v1.14.0]
github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b [v0.0.0-20241215232642-bb51bb14a506]
github.com/cockroachdb/pebble v1.1.0 [v1.1.5]
github.com/cockroachdb/redact v1.1.5 [v1.1.8]
github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 [v0.0.0-20250429170803-42689b6311bb]
github.com/coder/websocket v1.8.14 [v1.8.15]
github.com/containerd/log v0.1.0 [v0.2.0]
github.com/containerd/typeurl/v2 v2.2.0 [v2.3.0]
github.com/cznic/mathutil v0.0.0-20180504122225-ca4c9f2c1369 [v0.0.0-20181122101859-297441e03548]
github.com/danieljoos/wincred v1.1.2 [v1.2.3]
github.com/dgraph-io/badger/v4 v4.2.0 [v4.9.6]
github.com/dgraph-io/ristretto v0.1.1 [v0.2.0]
github.com/dlclark/regexp2 v1.11.0 [v1.12.0]
github.com/docker/docker v28.3.3+incompatible [v28.5.2+incompatible]
github.com/docker/go-connections v0.6.0 [v0.8.1]
github.com/dustin/go-humanize v1.0.0 [v1.0.1]
github.com/dvsekhvalnov/jose2go v1.7.0 [v1.11.0]
github.com/ebitengine/purego v0.10.0 [v0.11.0]
github.com/edsrzf/mmap-go v0.0.0-20170320065105-0bce6a688712 [v1.2.0]
github.com/envoyproxy/go-control-plane/envoy v1.32.4 [v1.39.0]
github.com/envoyproxy/protoc-gen-validate v1.2.1 [v1.3.3]
github.com/felixge/httpsnoop v1.0.4 [v1.1.0]
github.com/fsouza/fake-gcs-server v1.17.0 [v1.56.1]
github.com/gabriel-vasile/mimetype v1.4.1 [v1.4.15]
github.com/getsentry/sentry-go v0.18.0 [v0.49.0]

govulncheck and golangci-lint are at `/root/go/bin` but not on PATH here. Running them directly, plus inspecting the skip sites.
$ /root/go/bin/golangci-lint run ./... 2>&1 | tail -30
0 issues.

