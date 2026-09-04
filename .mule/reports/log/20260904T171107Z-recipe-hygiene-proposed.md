slug: recipe-hygiene
verdict: proposed
at:  672f739
ran: 2026-09-04T17:11:07Z on DietPi in 415s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md

$ govulncheck ./... 2>/dev/null || echo "govulncheck not installed"
govulncheck not installed

Let me run the hygiene checks. I'll start with all the commands in parallel.
$ which golangci-lint 2>/dev/null && golangci-lint run ./... 2>&1 | tail -30 || echo "golangci-lint not installed"
golangci-lint not installed

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ 2>&1 | head -30
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},
internal/store/store.go:135:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5):

→ Read .mule/todo.md
→ Read .mule/for-giulio.md

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

Let me check for skipped tests and look at the TODO context.
→ Read internal/store/store.go [offset=130, limit=15]

$ go test ./... 2>&1 | grep -i 'SKIP'
(no output)

Let me check if `internal/housekeeping` has logic worth testing (the only package with no test files and logic).
$ ls /root/astrate-mule/internal/housekeeping/
http.go
http_test.go
service.go

$ rg -n 'func ' /root/astrate-mule/internal/housekeeping/*.go 2>/dev/null | head -20
/root/astrate-mule/internal/housekeeping/http_test.go:38:func newHKRig(t *testing.T) *hkRig {
/root/astrate-mule/internal/housekeeping/http_test.go:76:func (r *hkRig) req(t *testing.T, method, path, rawBody, token string) *httptest.ResponseRecorder {
/root/astrate-mule/internal/housekeeping/http_test.go:91:func TestHousekeeping(t *testing.T) {
/root/astrate-mule/internal/housekeeping/http_test.go:171:func discardLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }
/root/astrate-mule/internal/housekeeping/http_test.go:173:func mintToken(t *testing.T, key *rsa.PrivateKey, claims jwt.MapClaims) string {
/root/astrate-mule/internal/housekeeping/http_test.go:183:func pubPEM(t *testing.T, pub *rsa.PublicKey) string {
/root/astrate-mule/internal/housekeeping/http_test.go:192:func decodeData(t *testing.T, rec *httptest.ResponseRecorder, dst any) {
/root/astrate-mule/internal/housekeeping/http_test.go:208:func contains(xs []string, s string) bool {
/root/astrate-mule/internal/housekeeping/http_test.go:217:func jsonStr(s string) string {
/root/astrate-mule/internal/housekeeping/http_test.go:222:func randSuffix(t *testing.T) string {
/root/astrate-mule/internal/housekeeping/service.go:49:func NewService(st *store.Store, sealer *store.KeySealer, reloader Reloader, log *slog.Logger) *Service {
/root/astrate-mule/internal/housekeeping/service.go:59:func (s *Service) notifyBrokerReload(ctx context.Context, realm string) {
/root/astrate-mule/internal/housekeeping/service.go:73:func (s *Service) CreateRealm(ctx context.Context, name, jwtPublicKeyPEM string, regLimit *int32) (*RealmView, error) {
/root/astrate-mule/internal/housekeeping/service.go:112:func (s *Service) GetRealm(ctx context.Context, name string) (*RealmView, error) {
/root/astrate-mule/internal/housekeeping/service.go:122:func (s *Service) ListRealms(ctx context.Context) ([]string, error) {
/root/astrate-mule/internal/housekeeping/service.go:137:func (s *Service) DeleteRealm(ctx context.Context, name string) error {
/root/astrate-mule/internal/housekeeping/service.go:146:func view(r *store.Realm) *RealmView {
/root/astrate-mule/internal/housekeeping/http.go:26:func NewAPI(svc *Service, mw *auth.Middleware, instanceKeysPEM []string) *API {
/root/astrate-mule/internal/housekeeping/http.go:32:func (a *API) Mount(mux *http.ServeMux) {
/root/astrate-mule/internal/housekeeping/http.go:48:func (a *API) listRealms(w http.ResponseWriter, r *http.Request) {

$ head -5 /root/astrate-mule/internal/housekeeping/http_test.go
//go:build integration

package housekeeping

import (

Let me check if any of the deps with updates are actually direct imports in this repo.
$ grep -E '^module |^\t' /root/astrate-mule/go.mod | head -60
module github.com/astrate-platform/astrate

$ grep -E '^\t[a-z]' /root/astrate-mule/go.mod | head -40
(no output)

→ Read go.mod
Let me get the full dependency-update list to check for direct deps with new versions.

$ go list -m -u all 2>/dev/null | rg '\['
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
github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.27.0 [v1.37.0] (deprecated)
github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.53.0 [v0.61.0] (deprecated)
github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.53.0 [v0.61.0]
github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 [v0.0.0-20240927000941-0f3dac36c52b]
github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a [v0.0.0-20230218143504-906a9b012302]
github.com/alicebob/miniredis/v2 v2.23.0 [v2.39.0]
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
github.com/go-jose/go-jose/v4 v4.0.5 [v4.1.5]
github.com/go-logr/logr v1.4.3 [v1.4.4]
github.com/go-ole/go-ole v1.2.6 [v1.3.0]
github.com/go-sql-driver/mysql v1.5.0 [v1.10.1]
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
github.com/googleapis/gax-go/v2 v2.15.0 [v2.24.1]
github.com/gorilla/handlers v1.4.2 [v1.5.2]
github.com/gorilla/mux v1.7.4 [v1.8.1]
github.com/jackc/pgerrcode v0.0.0-20220416144525-469b46aa5efa [v0.0.0-20250907135507-afb5586c32a6]
github.com/jackc/pgtype v1.14.0 [v1.14.4]
github.com/jackc/pgx/v4 v4.18.2 [v4.18.3]
github.com/jinzhu/copier v0.3.5 [v0.4.0]
github.com/k0kubun/pp v2.3.0+incompatible [v3.0.1+incompatible]
github.com/klauspost/compress v1.18.5 [v1.20.0]
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
github.com/moby/moby/api v1.54.2 [v1.56.0]
github.com/moby/moby/client v0.4.0 [v0.6.0]
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
golang.org/x/crypto v0.53.0 [v0.56.0]
golang.org/x/exp v0.0.0-20230626212559-97b1e661b5df [v0.0.0-20260824195058-e88cd73687aa]
golang.org/x/mod v0.36.0 [v0.40.0]
golang.org/x/net v0.55.0 [v0.58.0]
golang.org/x/oauth2 v0.30.0 [v0.36.0]
golang.org/x/sync v0.21.0 [v0.22.0]
golang.org/x/sys v0.46.0 [v0.47.0]
golang.org/x/telemetry v0.0.0-20251008203120-078029d740a8 [v0.0.0-20260902144106-3ef544be8421]
golang.org/x/term v0.44.0 [v0.45.0]
golang.org/x/text v0.38.0 [v0.41.0]
golang.org/x/time v0.12.0 [v0.15.0]
golang.org/x/tools v0.45.0 [v0.49.0]
golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 [v0.0.0-20240903120638-7835f813f4da]
google.golang.org/api v0.247.0 [v0.297.0]
google.golang.org/genproto v0.0.0-20250603155806-513f23925822 [v0.0.0-20260904163448-b1c236e22ff4]
google.golang.org/genproto/googleapis/api v0.0.0-20250818200422-3122310a409c [v0.0.0-20260904163448-b1c236e22ff4]
google.golang.org/genproto/googleapis/rpc v0.0.0-20250818200422-3122310a409c [v0.0.0-20260904163448-b1c236e22ff4]
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
modernc.org/libc v1.17.1 [v1.75.7]
modernc.org/lldb v1.0.0 [v1.0.8]
modernc.org/mathutil v1.5.0 [v1.7.1]
modernc.org/memory v1.2.1 [v1.12.1]
modernc.org/opt v0.1.3 [v0.2.0]
modernc.org/ql v1.0.0 [v1.6.0]
modernc.org/sortutil v1.1.0 [v1.2.1]
modernc.org/sqlite v1.18.1 [v1.58.0]
modernc.org/strutil v1.1.3 [v1.2.1]
modernc.org/token v1.0.0 [v1.1.0]
modernc.org/zappy v1.0.0 [v1.1.0]
pgregory.net/rapid v1.2.0 [v1.3.0]

✱ Grep "go\.etcd\.io/bbolt" 1 match
✱ Grep "go\.mongodb\.org/mongo-driver/v2" 8 matches
✱ Grep "nats-io/nats\.go" 3 matches
✱ Grep "prometheus/client_golang" 21 matches
✱ Grep "coder/websocket" 5 matches
The full sweep is mostly transitive testcontainers/migrate deps. Let me check what the repo actually uses from the directly-imported modules that have updates.
✱ Grep "testcontainers-go" 4 matches
The prior 2026-09-02 note claimed the sweep showed "only version-skew on transitive deps" — but that scan stopped at `head -20`, missing updates on several **directly pinned** deps. Let me check changelogs for the ones used in production code paths.
◈ Exa Web Search "go.etcd.io bbolt v1.5.0 release breaking changes changelog"
◈ Exa Web Search "coder/websocket v1.8.15 release changelog"
◈ Exa Web Search "go.mongodb.org mongo-driver v2.8.0 release breaking changes BSON"
coder/websocket v1.8.15 is a patch with a compression fix (no breaking), bbolt v1.5.0 is additive-only (bbolt's semver: no API change between patch/minor), mongo-driver 2.8 breaking changes are confined to Queryable Encryption options (repo uses BSON raw only). Let me check the remaining notable direct deps.
◈ Exa Web Search "prometheus client_golang v1.24.0 release notes breaking changes"
◈ Exa Web Search "nats.go v1.53.0 release notes breaking changes nats-io"
I have enough changelog data on most direct deps. Let me quickly check testcontainers v0.44 (the last unverified one) before writing up.
◈ Exa Web Search "testcontainers-go v0.44.0 release notes breaking changes"
testcontainers v0.44 has breaking changes but only in `wait.ForSQL` callbacks and `ImageProvider` (test-only harness here). I now have all the data I need. No new task lines are warranted — the only actionable TODO and the one test-gap package are already queued as `store-downsample-lttb` and `housekeeping-unit-tests`, there are no skipped tests, and govulncheck/golangci-lint are simply not installed. The genuinely new finding is the dependency sweep correction. Let me write it to for-giulio.md.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,20 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **Dependency sweep corrected: direct (pinned) deps DO have newer versions** — the 2026-09-02 note said the `go list -m -u` sweep showed "only version-skew on transitive deps", but that run hit the recipe's `head -20` cutoff (all cloud/azure/transitive) and never reached the directly-required modules. Full sweep, 2026-09-04. None of these is a fix this repo *needs*, so no bump is proposed — recorded for the decision. Per module (current → available; breaking change; repo use):
+  - `github.com/coder/websocket` v1.8.14 → v1.8.15 — no breaking (patch); used in `internal/appengine/stream/ws.go`, `channels/ws.go`; worth it only for the "transmit in single frame when compression enabled" fix + read-path alloc reduction.
+  - `go.etcd.io/bbolt` v1.4.3 → v1.5.0 — bbolt's semver promises no API change between patch/minor, so additive-only; used in `internal/broker/sessionstore.go`; v1.5 adds a data-file size limit and panic-recovery hardening, nothing Astrate needs.
+  - `go.mongodb.org/mongo-driver/v2` v2.6.0 → v2.8.2 — the 2.8.0 breaking changes are confined to Queryable Encryption string-query options (`options.Text()`→`String()`); Astrate uses only the raw BSON API (`pkg/payload/bson.go`, `internal/engine/capabilities.go`, `bench/`) and is unaffected.
+  - `github.com/nats-io/nats.go` v1.52.0 → v1.53.1 — no breaking; the headline fixes (JetStream `resetOrderedConsumer` race, KV dot-rejection) are paths Astrate does not use — `internal/engine/forward/nats.go` is core NATS publish only.
+  - `github.com/prometheus/client_golang` v1.23.2 → v1.24.1 — requires Go ≥1.25 (fine, repo is 1.26.1); the breaking `LabelNames`/remote-api renames don't touch repo usage (`prometheus`/`collectors`/`promhttp` in `internal/observability/metrics.go`, flow/engine metrics); would buy `Gather()` panic-recovery and opt-in `CoalesceGather` scrape-pile-up protection.
+  - `github.com/testcontainers/testcontainers-go` v0.43.0 → v0.44.0 (modules/postgres v0.42.0, modules/nats v0.43.0) — breaking in `wait.ForSQL` (callback now takes `network.Port`) and `ImageProvider` (new `PullImageWithPlatform`); Astrate's `internal/testutil/pg.go` looks unaffected but it is test-only anyway.
+  - `golang.org/x/crypto` v0.53.0 → v0.56.0 — x/crypto keeps API compatibility; used only for bcrypt in `internal/auth`.
+  Note: `govulncheck` and `golangci-lint` are still not installed on the Pi (existing entries above), so the vulnerability and lint checks remain unavailable.
+
+---
+
 - **milestone 2.0 looks complete, verify and cut the tag** — all 11 `milestone-2.0` issues
   CLOSED (#23–#27, #37, #39–#43), no open issues, no new gaps after re-checking upstream
   astarte_flow block catalog against `internal/flow/` + git log (MQTT/HTTP source/sink,
   json_path_map, pure-transform set, virtual_device_pool, container block MVP, flow API,

No new task lines are warranted. Let me double-check for actual `t.Skip` calls (the filtered `go test` output showed no runtime skips, but let me confirm at source level too).
✱ Grep "t\.Skip" 12 matches
Most skips are environmental (DB, external binaries/toolkit, `-short` mode) — legitimate, not "stopped being a test". The one worth a closer look is `forward_test.go:68` with its build-tag conditional.
→ Read cmd/astrate/forward_test.go [limit=80]
This skip is a legitimate both-build-modes pattern, not a dead test. The analysis is complete.

**Checks run and outcome:**

| check | result |
|---|---|
| `go list -m -u` | direct deps **do** have updates (previous run truncated at `head -20` and missed them) → written to `.mule/for-giulio.md` |
| `govulncheck` | not installed (already on record) |
| `golangci-lint` | not installed (already on record) |
| TODO/FIXME scan | only `internal/store/store.go:135` → already queued as `store-downsample-lttb` |
| packages w/o unit tests | only `internal/housekeeping` → already queued as `housekeeping-unit-tests` |
| skipped tests | all environmental (DB/external binaries/toolkit/`-short`) or a both-build-modes conditional — none is a dead test |

No new task lines were appended — nothing new is worth proposing, which the recipe says is the normal outcome. The one genuinely new finding went to `for-giulio.md`. No git touched.

```
Done: Hygiene run: no new task lines; corrected dependency-sweep gap (direct deps have newer versions) recorded in .mule/for-giulio.md
Files: .mule/for-giulio.md
Verified: hygiene checks (deps/govulncheck/lint/TODO/test-gaps/skips) -> govulncheck & golangci-lint fail (not installed, on record); rest pass
Unsure: changelog claims for coder/websocket, bbolt, mongo-driver v2, nats.go, client_golang, testcontainers are from release notes, not from reading each tagged source
Follow-ups: none
```
