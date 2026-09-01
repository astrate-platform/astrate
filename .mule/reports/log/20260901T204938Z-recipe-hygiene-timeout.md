slug: recipe-hygiene
verdict: timeout
at:  e68d9e3
ran: 2026-09-01T20:49:39Z on DietPi in 1404s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md

$ govulncheck ./... 2>/dev/null || echo "govulncheck not installed"
govulncheck not installed

Recipe loaded. Running the hygiene checks now.
$ which golangci-lint; ls ~/go/bin/golangci-lint /Users/atsetilam/go/bin/golangci-lint 2>/dev/null; golangci-lint --version 2>/dev/null
/root/go/bin/golangci-lint

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ | head -30
internal/store/store.go:135:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5):
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},

→ Read .mule/todo.md

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

$ golangci-lint run ./... 2>&1 | tail -40
/bin/sh: 1: golangci-lint: not found

$ go list -m -f '{{if not .Indirect}}{{.Path}} {{.Version}}{{end}}' all
github.com/astrate-platform/astrate 
github.com/BurntSushi/toml v1.6.0
github.com/coder/websocket v1.8.14
github.com/eclipse/paho.mqtt.golang v1.5.1
github.com/golang-jwt/jwt/v5 v5.3.1
github.com/golang-migrate/migrate/v4 v4.19.1
github.com/hashicorp/golang-lru/v2 v2.0.7
github.com/jackc/pgx/v5 v5.10.0
github.com/mochi-mqtt/server/v2 v2.7.9
github.com/nats-io/nats.go v1.52.0
github.com/prometheus/client_golang v1.23.2
github.com/testcontainers/testcontainers-go v0.43.0
github.com/testcontainers/testcontainers-go/modules/nats v0.43.0
github.com/testcontainers/testcontainers-go/modules/postgres v0.42.0
go.etcd.io/bbolt v1.4.3
go.mongodb.org/mongo-driver/v2 v2.6.0
golang.org/x/crypto v0.53.0

$ go test ./... 2>&1 | rg -i 'no test files'
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/internal/housekeeping	[no test files]
?   	github.com/astrate-platform/astrate/internal/swagger	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]

$ go list -m -u all 2>/dev/null | rg '\['
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
github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.27.0 [v1.37.0] (deprecated)
github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.53.0 [v0.61.0] (deprecated)
github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.53.0 [v0.61.0]
github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 [v0.0.0-20240927000941-0f3dac36c52b]
github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a [v0.0.0-20230218143504-906a9b012302]
github.com/alicebob/miniredis/v2 v2.23.0 [v2.38.0]
github.com/andybalholm/brotli v1.0.4 [v1.2.3]
github.com/apache/thrift v0.16.0 [v0.24.0]
github.com/aws/aws-sdk-go v1.49.6 [v1.55.8] (deprecated)
github.com/aws/aws-sdk-go-v2 v1.16.16 [v1.45.1]
github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.8 [v1.7.20]
github.com/aws/aws-sdk-go-v2/credentials v1.12.20 [v1.20.2]
github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.33 [v1.23.2] (deprecated)
github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.23 [v1.5.1]
github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.17 [v2.8.1]
github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.14 [v1.5.1]
github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.9 [v1.13.19]
github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.18 [v1.11.1]
github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.17 [v1.14.1]
github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.17 [v1.20.1]
github.com/aws/aws-sdk-go-v2/service/s3 v1.27.11 [v1.110.0]
github.com/aws/smithy-go v1.13.3 [v1.28.1]
github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58 [v0.0.0-20240916140612-caecf3c00c06]
github.com/cncf/xds/go v0.0.0-20250501225837-2ac532fd4443 [v0.0.0-20260202195803-dba9d589def2]
github.com/cockroachdb/cockroach-go/v2 v2.1.1 [v2.4.3]
github.com/cockroachdb/errors v1.11.1 [v1.14.0]
github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b [v0.0.0-20241215232642-bb51bb14a506]
github.com/cockroachdb/pebble v1.1.0 [v1.1.5]
github.com/cockroachdb/redact v1.1.5 [v1.1.8]
github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 [v0.0.0-20250429170803-42689b6311bb]
github.com/coder/websocket v1.8.14 [v1.8.15]
github.com/containerd/typeurl/v2 v2.2.0 [v2.3.0]
github.com/cznic/mathutil v0.0.0-20180504122225-ca4c9f2c1369 [v0.0.0-20181122101859-297441e03548]
github.com/danieljoos/wincred v1.1.2 [v1.2.3]
github.com/dgraph-io/badger/v4 v4.2.0 [v4.9.6]
github.com/dgraph-io/ristretto v0.1.1 [v0.2.0]
github.com/docker/docker v28.3.3+incompatible [v28.5.2+incompatible]
github.com/docker/go-connections v0.6.0 [v0.8.1]
github.com/dustin/go-humanize v1.0.0 [v1.0.1]
github.com/dvsekhvalnov/jose2go v1.7.0 [v1.10.0]
github.com/ebitengine/purego v0.10.0 [v0.11.0]
github.com/edsrzf/mmap-go v0.0.0-20170320065105-0bce6a688712 [v1.2.0]
github.com/envoyproxy/go-control-plane/envoy v1.32.4 [v1.39.0]
github.com/envoyproxy/protoc-gen-validate v1.2.1 [v1.3.3]
github.com/felixge/httpsnoop v1.0.4 [v1.1.0]
github.com/fsouza/fake-gcs-server v1.17.0 [v1.56.1]
github.com/gabriel-vasile/mimetype v1.4.1 [v1.4.15]
github.com/getsentry/sentry-go v0.18.0 [v0.49.0]
github.com/go-jose/go-jose/v4 v4.0.5 [v4.1.4]
github.com/go-logr/logr v1.4.3 [v1.4.4]
github.com/go-ole/go-ole v1.2.6 [v1.3.0]
github.com/go-sql-driver/mysql v1.5.0 [v1.10.0]
github.com/go-stack/stack v1.8.0 [v1.8.1]
github.com/gobuffalo/here v0.6.0 [v0.6.7]
github.com/goccy/go-json v0.9.11 [v0.10.6]
github.com/gocql/gocql v0.0.0-20210515062232-b7ef815b4556 [v1.7.0]
github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe [v0.0.0-20220223132316-b832511892a9]
github.com/golang/glog v1.2.4 [v1.2.5]
github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da [v0.0.0-20241129210726-2c02b8208cf8]
github.com/golang/protobuf v1.5.2 [v1.5.4] (deprecated)
github.com/golang/snappy v0.0.4 [v1.0.0]
github.com/google/flatbuffers v2.0.8+incompatible [v25.12.19+incompatible]
github.com/google/go-querystring v1.1.0 [v1.2.0]
github.com/googleapis/enterprise-certificate-proxy v0.3.6 [v0.3.21]
github.com/googleapis/gax-go/v2 v2.15.0 [v2.24.0]
github.com/gorilla/handlers v1.4.2 [v1.5.2]
github.com/gorilla/mux v1.7.4 [v1.8.1]
github.com/jackc/pgerrcode v0.0.0-20220416144525-469b46aa5efa [v0.0.0-20250907135507-afb5586c32a6]
github.com/jackc/pgtype v1.14.0 [v1.14.4]
github.com/jackc/pgx/v4 v4.18.2 [v4.18.3]
github.com/jinzhu/copier v0.3.5 [v0.4.0]
github.com/k0kubun/pp v2.3.0+incompatible [v3.0.1+incompatible]
github.com/klauspost/compress v1.18.5 [v1.19.2]
github.com/klauspost/cpuid/v2 v2.0.9 [v2.4.0]
github.com/ktrysmt/go-bitbucket v0.6.4 [v0.10.0]
github.com/lib/pq v1.10.9 [v1.12.3]
github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 [v0.0.0-20260802145828-341c2f0c90b5]
github.com/magiconair/properties v1.8.10 [v1.18.11]
github.com/markbates/pkger v0.15.1 [v0.17.1]
github.com/mattn/go-colorable v0.1.6 [v0.1.15]
github.com/mattn/go-isatty v0.0.16 [v0.0.24]
github.com/mattn/go-sqlite3 v1.14.22 [v1.14.50]
github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 [v1.0.4]
github.com/microsoft/go-mssqldb v1.0.0 [v1.11.0]
github.com/mitchellh/mapstructure v1.1.2 [v1.5.0]
github.com/moby/go-archive v0.2.0 [v0.3.3]
github.com/moby/moby/api v1.54.2 [v1.55.0]
github.com/moby/moby/client v0.4.0 [v0.5.1]
github.com/moby/sys/mount v0.3.4 [v0.3.5]
github.com/moby/sys/sequential v0.6.0 [v0.7.0]
github.com/moby/sys/user v0.4.0 [v0.4.1]
github.com/moby/sys/userns v0.1.0 [v0.2.0]
github.com/morikuni/aec v1.0.0 [v1.1.0]
github.com/mutecomm/go-sqlcipher/v4 v4.4.0 [v4.4.2]
github.com/nakagami/firebirdsql v0.0.0-20190310045651-3c02a58cfed8 [v0.9.20]
github.com/nats-io/nats.go v1.52.0 [v1.53.1]
github.com/nats-io/nkeys v0.4.15 [v0.4.16]
github.com/neo4j/neo4j-go-driver v1.8.1-0.20200803113522-b626aa943eba [v1.8.3]
github.com/onsi/ginkgo v1.16.4 [v1.16.5]
github.com/onsi/gomega v1.15.0 [v1.43.0]
github.com/pierrec/lz4/v4 v4.1.16 [v4.1.29]
github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 [v0.0.0-20240102092130-5ac0b6a4141c]
github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 [v0.0.0-20260805114148-88456608a4f6]
github.com/prometheus/client_golang v1.23.2 [v1.24.1]
github.com/prometheus/client_model v0.6.2 [v0.6.3]
github.com/prometheus/common v0.66.1 [v0.71.0]
github.com/prometheus/procfs v0.16.1 [v0.22.0]
github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 [v0.0.0-20230129092748-24d4a6f8daec]
github.com/rogpeppe/go-internal v1.14.1 [v1.16.0]
github.com/rqlite/gorqlite v0.0.0-20230708021416-2acd02b70b79 [v0.0.0-20260504155303-50d445fd0ab9]
github.com/rs/xid v1.4.0 [v1.6.0]
github.com/shirou/gopsutil/v4 v4.26.5 [v4.26.8]
github.com/shopspring/decimal v1.2.0 [v1.4.0]
github.com/sirupsen/logrus v1.9.4 [v1.10.2]
github.com/snowflakedb/gosnowflake v1.6.19 [v1.19.1]
github.com/spf13/cobra v1.8.1 [v1.10.2]
github.com/spf13/pflag v1.0.6 [v1.0.10]
github.com/spiffe/go-spiffe/v2 v2.5.0 [v2.8.1]
github.com/stretchr/testify v1.11.1 [v1.12.1]
github.com/testcontainers/testcontainers-go v0.43.0 [v0.44.0]
github.com/testcontainers/testcontainers-go/modules/nats v0.43.0 [v0.44.0]
github.com/testcontainers/testcontainers-go/modules/postgres v0.42.0 [v0.44.0]
github.com/tklauser/go-sysconf v0.3.16 [v0.4.0]
github.com/tklauser/numcpus v0.11.0 [v0.12.0]
github.com/xanzy/go-gitlab v0.15.0 [v0.115.0] (deprecated)
github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9 [v1.1.2]
github.com/zeebo/xxh3 v1.0.2 [v1.1.0]
go.etcd.io/bbolt v1.4.3 [v1.5.0]
go.mongodb.org/mongo-driver v1.7.5 [v1.17.9] (deprecated)
go.mongodb.org/mongo-driver/v2 v2.6.0 [v2.8.2]
go.opentelemetry.io/contrib/detectors/gcp v1.36.0 [v1.46.0]
go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.61.0 [v0.71.0]
go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0 [v0.71.0]
go.opentelemetry.io/otel v1.41.0 [v1.46.0]
go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.29.0 [v1.46.0]
go.opentelemetry.io/otel/metric v1.41.0 [v1.46.0]
go.opentelemetry.io/otel/sdk v1.36.0 [v1.46.0]
go.opentelemetry.io/otel/sdk/metric v1.36.0 [v1.46.0]
go.opentelemetry.io/otel/trace v1.41.0 [v1.46.0]
go.opentelemetry.io/proto/otlp v1.3.1 [v1.11.0]
go.yaml.in/yaml/v2 v2.4.2 [v2.4.4]
golang.org/x/crypto v0.53.0 [v0.55.0]
golang.org/x/exp v0.0.0-20230626212559-97b1e661b5df [v0.0.0-20260824195058-e88cd73687aa]
golang.org/x/mod v0.36.0 [v0.40.0]
golang.org/x/net v0.55.0 [v0.58.0]
golang.org/x/oauth2 v0.30.0 [v0.36.0]
golang.org/x/sync v0.21.0 [v0.22.0]
golang.org/x/sys v0.46.0 [v0.47.0]
golang.org/x/telemetry v0.0.0-20251008203120-078029d740a8 [v0.0.0-20260828145429-86cb5733f5b7]
golang.org/x/term v0.44.0 [v0.45.0]
golang.org/x/text v0.38.0 [v0.41.0]
golang.org/x/time v0.12.0 [v0.15.0]
golang.org/x/tools v0.45.0 [v0.49.0]
golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 [v0.0.0-20240903120638-7835f813f4da]
google.golang.org/api v0.247.0 [v0.296.0]
google.golang.org/genproto v0.0.0-20250603155806-513f23925822 [v0.0.0-20260831171406-18b4a7587f8a]
google.golang.org/genproto/googleapis/api v0.0.0-20250818200422-3122310a409c [v0.0.0-20260831171406-18b4a7587f8a]
google.golang.org/genproto/googleapis/rpc v0.0.0-20250818200422-3122310a409c [v0.0.0-20260831171406-18b4a7587f8a]
google.golang.org/grpc v1.74.2 [v1.83.2]
google.golang.org/protobuf v1.36.8 [v1.36.12]
lukechampine.com/uint128 v1.2.0 [v1.3.0]
modernc.org/b v1.0.0 [v1.1.0]
modernc.org/cc/v3 v3.36.3 [v3.41.0]
modernc.org/ccgo/v3 v3.16.9 [v3.17.0]
modernc.org/db v1.0.0 [v1.4.0]
modernc.org/file v1.0.0 [v1.1.4]
modernc.org/fileutil v1.0.0 [v1.4.0]
modernc.org/golex v1.0.0 [v1.1.0]
modernc.org/internal v1.0.0 [v1.1.12]
modernc.org/libc v1.17.1 [v1.75.6]
modernc.org/lldb v1.0.0 [v1.0.8]
modernc.org/mathutil v1.5.0 [v1.7.1]
modernc.org/memory v1.2.1 [v1.12.1]
modernc.org/opt v0.1.3 [v0.2.0]
modernc.org/ql v1.0.0 [v1.6.0]
modernc.org/sortutil v1.1.0 [v1.2.1]
modernc.org/sqlite v1.18.1 [v1.57.0]
modernc.org/strutil v1.1.3 [v1.2.1]
modernc.org/token v1.0.0 [v1.1.0]
modernc.org/zappy v1.0.0 [v1.1.0]
pgregory.net/rapid v1.2.0 [v1.3.0]

$ /root/go/bin/golangci-lint run ./... 2>&1 | tail -50
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

$ go run golang.org/x/vuln/cmd/govulncheck@latest ./... 2>&1 | tail -40
signal: killed

govulncheck OOM-killed. Retrying with lighter settings.
