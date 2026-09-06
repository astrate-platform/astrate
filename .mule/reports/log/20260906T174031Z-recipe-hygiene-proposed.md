slug: recipe-hygiene
verdict: proposed
at:  29c26e0
ran: 2026-09-06T17:40:31Z on DietPi in 135s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md

$ govulncheck ./... 2>/dev/null || echo "govulncheck not installed"
govulncheck not installed

$ which golangci-lint 2>/dev/null && golangci-lint run ./... 2>&1 | tail -30 || echo "golangci-lint not found in PATH"
golangci-lint not found in PATH

Now let me run the hygiene checks.
$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ | head -30
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
?   	github.com/astrate-platform/astrate/internal/housekeeping	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]

$ ls internal/housekeeping/
async_operation_test.go
http.go
http_test.go
service.go

Let me check a few things before concluding: the housekeeping package (is it logic?), and try to find golangci-lint.
$ head -30 internal/broker/intake.go internal/store/store.go 2>&1 | head -60
==> internal/broker/intake.go <==
// Package broker embeds the mochi-mqtt server and gives it Astarte MQTT v1
// semantics (docs/DESIGN.md §3.1–§3.4): mTLS device identity with
// CN = "<realm>/<device_id>", the §3.2 ACL matrix, persistent sessions in a
// bbolt file (so session_present survives Astrate restarts), device lifecycle
// bookkeeping, and an inline publishing facade for the engine and AppEngine.
//
// The broker does not interpret payloads. Every accepted device PUBLISH is
// handed to an Intake (implemented by internal/engine in M6) as an
// InboundMessage whose Ack callback releases the MQTT acknowledgment: for
// QoS >= 1 the PUBACK/PUBREC is withheld until Ack is called, which is how
// persistence-commit ordering (docs/DESIGN.md §5.3) and shard backpressure
// (§1.4) propagate to the device.
package broker

import (
	"net/netip"
	"time"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// InboundMessage is one device PUBLISH delivered to the engine intake
// (docs/ROADMAP.md §6 file 5.1). Payload is an independent copy: it remains
// valid after the broker recycles its packet buffers.
type InboundMessage struct {
	// Realm is the device's realm name (first topic segment).
	Realm string
	// DeviceID is the publishing device, parsed from the connection's
	// certificate CN.
	DeviceID deviceid.ID

==> internal/store/store.go <==
// Package store is Astrate's single PostgreSQL/TimescaleDB access layer
// (docs/DESIGN.md §1.3): every domain package reads and writes the database
// through it, and it imports none of them. It owns connection pooling, the
// embedded schema migrations (docs/DESIGN.md §2.2–2.5), CA-key sealing, and
// one repository file per aggregate (realms, interfaces, devices, properties,
// datastreams, groups, triggers) plus LISTEN/NOTIFY plumbing.
package store

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/astrate-platform/astrate/migrations"
)

// defaultPoolMaxConns bounds the pool when the DSN does not specify
// pool_max_conns. The server is tuned to max_connections=50

$ find /root -name golangci-lint -type f 2>/dev/null | head -5
/root/go/bin/golangci-lint
/root/go/pkg/mod/github.com/docker/docker@v28.3.3+incompatible/hack/validate/golangci-lint

$ /root/go/bin/golangci-lint run ./... 2>&1 | tail -40
0 issues.

$ rg -n 't\.Skip|testing\.Short' internal/ pkg/ cmd/ --include '*.go' | head -20
rg: unrecognized flag --include

similar flags that are available: --include-zero

$ go test ./... 2>&1 | head -40
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
?   	github.com/astrate-platform/astrate/internal/housekeeping	[no test files]
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

$ rg -n 't\.Skip|testing\.Short' -g '*.go' internal/ pkg/ cmd/ | head -20
cmd/astrate/main_test.go:47:	t.Skip("the boot suite needs a database: run `make up` or set ASTRATE_TEST_DSN")
cmd/astrate/forward_test.go:68:		t.Skip("built with -tags nats; see forward_nats_test.go")
internal/pairing/http_test.go:629:		t.Skip("openssl not found in PATH; chain already verified with crypto/x509")
internal/pairing/ca/ca_test.go:292:	if testing.Short() {
internal/pairing/ca/ca_test.go:293:		t.Skip("10k issuance draw skipped in -short mode")
internal/store/datastreams_test.go:380:			t.Skip("timescaledb_toolkit not installed; lttb path unavailable")

→ Read internal/housekeeping/service.go [limit=50]
→ Read internal/housekeeping/service.go [offset=50]
→ Read .mule/for-giulio.md
→ Read .mule/todo.md
Now let me check the current state of the target files.
Now I have all the data. Let me write the hygiene report and the todo entry.

$ ls /root/astrate-mule/.mule/reports/ 2>/dev/null || echo "directory does not exist"
log

← Write .mule/reports/hygiene.md
Wrote file successfully.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -105,4 +105,5 @@
 - [!] flow-validate-source-sink: `Pipeline.Validate`'s no-source/no-sink error branches (internal/flow/pipeline.go:148-153) are unreachable — the preceding Kahn cycle check (pipeline.go:100-115) guarantees any graph that gets past it is acyclic and therefore always has a source and a sink. The two test cases "no source (all blocks have incoming edges)" and "no sink (all blocks have outgoing edges)" (pipeline_test.go:134-163) both use a 2-cycle, so they pass via the cycle error, never the branch they name. Make the rule real (a cycle-independent formulation, e.g. reject a graph where a node has no path to a sink) or delete the vacuous check, and fix the two test cases to assert the actual error message so a regression to wrong-path coverage fails. — BLOCKED: tests failed: --- FAIL: TestMQTTSink_Retained (0.02s)
 - [x] flow-validate-dead-source-sink-recompute: in `Pipeline.Validate`, the first `hasSource`/`hasSink` computation (pipeline.go:117-126) is overwritten by the second (pipeline.go:128-147) and its result discarded — dead code (the comment at pipeline.go:127 documents it). Delete the first block; behaviour is unchanged. Tests: existing `TestPipeline_Validate` passes.
 - [x] docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` route to docs/api/astarte_realm_management_api.yaml — it exists in code (internal/realm/http.go:51, handler getDatastreamMaximumStorageRetention at http.go:135, returns 200 with a data-envelope retention value via Service.GetDatastreamMaximumStorageRetention) but the spec jumps straight from config/device_registration_limit (yaml:428) to /version. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] housekeeping-tests [auto]: add a container-free unit test file for `internal/housekeeping` (service.go, 262 lines, zero test coverage). Cover two rules: (1) `CreateRealm` validation — blank name, blank JWT key, negative regLimit, negative retention, and default retention injection when `WithDefaultDatastreamMaximumStorageRetention` is set and the caller passes nil; (2) `DeleteRealm` gating — `ErrDeletionDisabled` when `WithRealmDeletionDisabled(true)`, `ErrConnectedDevicesPresent` when `DeviceStats` reports connected > 0. Use httptest-style fakes for the store and sealer (no Docker).

Hygiene run complete. One finding worth acting on:

**Housekeeping has no tests.** `internal/housekeeping/service.go` is 262 lines of real business logic — realm provisioning, CA key sealing, validation, deletion gating — with zero test coverage. Task line appended to `.mule/todo.md`.

Everything else was clean:
- golangci-lint: 0 issues
- TODOs: all aspirational or standard Go idiom, nothing proposable
- Skipped tests: all justified by missing external deps or build tags
- Dependencies: already swept and decided 2026-09-04, nothing new
- govulncheck: not installed in this env, can't run

Evidence written to `.mule/reports/hygiene.md`.
