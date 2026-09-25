slug: recipe-code-review
verdict: proposed
at:  2493a99
ran: 2026-09-25T17:47:27Z on DietPi in 550s


> build · big-pickle

I'll start by reading the required context files.
→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
→ Read .mule/todo.md

$ ls /root/astrate-mule/.mule/ /root/astrate-mule/.mule/reviews/ 2>&1; echo "---INTERNAL---"; ls /root/astrate-mule/internal/ 2>&1; echo "---PKG---"; ls /root/astrate-mule/pkg/ 2>&1
/root/astrate-mule/.mule/:
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
recipes
reports
reviews
tasks
todo.md

/root/astrate-mule/.mule/reviews/:
appengine-2026-09-07.md
astarteapi-2026-09-15.md
auth-2026-09-14.md
broker-2026-09-04.md
broker-2026-09-21.md
channels-2026-09-18.md
config-2026-09-13.md
deviceid-2026-09-18.md
engine-2026-09-06.md
engine-triggers-2026-09-20.md
flow-2026-09-05.md
flow-2026-09-23.md
flowapi-2026-09-10.md
housekeeping-2026-09-15.md
httpx-2026-09-17.md
interfaceschema-2026-09-11.md
observability-2026-09-16.md
pairing-2026-09-09.md
payload-2026-09-12.md
realm-2026-09-09.md
store-2026-09-08.md
store-2026-09-21.md
swagger-2026-09-24.md
---INTERNAL---
appengine
auth
broker
config
engine
flow
flowapi
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil
---PKG---
astarteapi
deviceid
interfaceschema
payload

$ cd /root/astrate-mule && rg -c '' --glob 'internal/engine/*.go' | sort -t: -k2 -rn; echo "--- reviews sizes ---"; wc -l .mule/reviews/engine-2026-09-06.md .mule/reviews/appengine-2026-09-07.md; echo "--- log tail ---"; tail -30 .mule/log.md
internal/engine/fakestore_test.go:818
internal/engine/data.go:552
internal/engine/cache.go:502
internal/engine/router.go:490
internal/engine/data_test.go:461
internal/engine/batch.go:459
internal/engine/engine.go:456
internal/engine/batch_test.go:446
internal/engine/engine_e2e_test.go:440
internal/engine/engine_integration_test.go:435
internal/engine/serverdata_test.go:390
internal/engine/cache_test.go:363
internal/engine/router_test.go:362
internal/engine/serverdata.go:349
internal/engine/control_test.go:330
internal/engine/control.go:316
internal/engine/triggers_change_test.go:293
internal/engine/capabilities_test.go:255
internal/engine/devicedata_test.go:223
internal/engine/introspection_test.go:189
internal/engine/introspection.go:168
internal/engine/mqttv1_parity_test.go:151
internal/engine/triggers_change_integration_test.go:141
internal/engine/bus_events_test.go:140
internal/engine/topics_test.go:109
internal/engine/capabilities.go:102
internal/engine/topics.go:97
--- reviews sizes ---
  29 .mule/reviews/engine-2026-09-06.md
  81 .mule/reviews/appengine-2026-09-07.md
 110 total
--- log tail ---
| 2026-09-16 | observability-health-content-type-test | done | 250s | f6c27f5 |
| 2026-09-17 | docs-sync-pairing-version-value | done | 1059s | c69a164 |
| 2026-09-17 | docs-sync-native-version-value | done | 503s | 08048b3 |
| 2026-09-17 | httpx-cors-vary-origin-passthrough | done | 143s | 986c100 |
| 2026-09-17 | docs-sync-rm-policies-delete-422 | done | 117s | c0521bc |
| 2026-09-17 | docs-sync-rm-triggers-422-nested-envelope | done | 161s | 8bfedd6 |
| 2026-09-18 | compat-note-v134 | done | 154s | bcd4383 |
| 2026-09-18 | channels-group-watch-membership-scope | done | 677s | aeaeb47 |
| 2026-09-18 | channels-rejoin-joinref-tagging-test | done | 208s | d48e17f |
| 2026-09-19 | channels-rejoin-authz-mismatch | blocked | 179s | wrote nothing |
| 2026-09-20 | errorname-missing-required-test | done | 135s | 739df67 |
| 2026-09-20 | triggers-custom-action-policy-nodecide | done | 459s | 3c8b89a |
| 2026-09-20 | compat-note-custom-action-policy-boundary | done | 97s | 8cd17a8 |
| 2026-09-20 | docs-sync-native-metrics-example-fake-series | done | 93s | 43c4a8b |
| 2026-09-20 | docs-sync-native-socket-missing-403-500 | done | 116s | f9d8b27 |
| 2026-09-21 | store-pipelines-empty-name-zero-blocks-test | done | 91s | 82fc0a4 |
| 2026-09-21 | broker-acl-coldstart-fallback-flood | done | 218s | 0e5078f |
| 2026-09-21 | broker-offlineacl-entry-eviction | done | 406s | 72b2686 |
| 2026-09-22 | docs-sync-ae-read-query-params | blocked | 159s | wrote nothing |
| 2026-09-23 | examples-echo-container-contract-test | done | 251s | f35af21 |
| 2026-09-23 | flow-mqtt-source-reconnect-recovery | blocked | 84s | wrote nothing |
| 2026-09-23 | flow-msg-json-integer-precision | done | 612s | 4df3f6b |
| 2026-09-23 | flow-sort-bounded-buffer | done | 634s | 9bf0ab7 |
| 2026-09-23 | flow-randomsource-span-overflow | blocked | 288s | lint failed: internal/flow/blocks/randomsource.go:113:40: G115: integer overflow conversion int64 -> uint64 (gosec) |
| 2026-09-23 | flow-filter-key-contains-test | done | 107s | 3875ba6 |
| 2026-09-23 | docs-sync-rm-put-interface-409 | done | 356s | 7256af9 |
| 2026-09-24 | docs-sync-rm-interface-422-shapes | done | 490s | 93e784e |
| 2026-09-24 | docs-sync-rm-validation-example-prefix | done | 238s | 02e269b |
| 2026-09-24 | docs-sync-rm-auth-403 | done | 668s | 9c05370 |
| 2026-09-24 | swagger-sub-failfast | done | 371s | 1d321ea |

→ Read .mule/todo.md [offset=174]
→ Read .mule/reviews/engine-2026-09-06.md
→ Read .mule/reviews/flow-2026-09-23.md

$ cd /root/astrate-mule && rg -c '' --glob 'cmd/**/*.go' --glob 'internal/testutil/*.go' | sort -t: -k2 -rn; echo "--- log head/recipe list ---"; ls .mule/recipes/; echo "--- reviews mentioning cmd/ ---"; rg -l 'cmd/astrate' .mule/reviews/ | head
cmd/astrate/main.go:581
cmd/astrate/main_test.go:265
internal/testutil/astartedevice.go:216
internal/testutil/mqttclient.go:163
internal/testutil/golden_test.go:117
internal/testutil/pg.go:106
cmd/astrate/forward_test.go:79
internal/testutil/golden.go:59
cmd/astrate/forward_nats_test.go:59
internal/testutil/pg_test.go:36
cmd/astrate/version_test.go:29
cmd/astrate/newnats_default.go:21
cmd/astrate/newnats_nats.go:19
--- log head/recipe list ---
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md
--- reviews mentioning cmd/ ---
.mule/reviews/flowapi-2026-09-10.md
.mule/reviews/engine-triggers-2026-09-20.md
.mule/reviews/httpx-2026-09-17.md
.mule/reviews/auth-2026-09-14.md
.mule/reviews/swagger-2026-09-24.md
.mule/reviews/observability-2026-09-16.md
.mule/reviews/deviceid-2026-09-18.md

→ Read cmd/astrate/main.go

$ cd /root/astrate-mule && rg -c '' --glob 'internal/engine/forward/*.go'; echo "=== newnats_default.go ==="; cat cmd/astrate/newnats_default.go; echo "=== newnats_nats.go ==="; cat cmd/astrate/newnats_nats.go; echo "=== go.mod head ==="; head -5 go.mod; echo "=== Makefile targets ==="; rg -n 'nats|tags' Makefile 2>/dev/null | head -20
internal/engine/forward/envelope.go:29
internal/engine/forward/nats.go:53
internal/engine/forward/http.go:125
internal/engine/forward/http_test.go:310
internal/engine/forward/nats_test.go:196
=== newnats_default.go ===
//go:build !nats

package main

import (
	"fmt"

	"github.com/astrate-platform/astrate/internal/config"
	"github.com/astrate-platform/astrate/internal/engine/triggers"
)

// natsBuildTagEnabled lets untagged test files (forward_test.go) tell which
// half of the build-tag split they are running against.
const natsBuildTagEnabled = false

// newNATSForwarder reports that this binary was not built with -tags nats, so
// triggers.forward.kind = "nats" fails at boot rather than silently
// forwarding nothing. See newnats_nats.go for the tagged build.
func newNATSForwarder(_ config.ForwardConfig) (triggers.Forwarder, error) {
	return nil, fmt.Errorf("triggers.forward.kind \"nats\" requires building astrate with -tags nats")
}
=== newnats_nats.go ===
//go:build nats

package main

import (
	"github.com/astrate-platform/astrate/internal/config"
	"github.com/astrate-platform/astrate/internal/engine/forward"
	"github.com/astrate-platform/astrate/internal/engine/triggers"
)

// natsBuildTagEnabled lets untagged test files (forward_test.go) tell which
// half of the build-tag split they are running against.
const natsBuildTagEnabled = true

// newNATSForwarder builds the NATS bus forwarder. Only compiled into a binary
// built with -tags nats; see newnats_default.go for the untagged build.
func newNATSForwarder(f config.ForwardConfig) (triggers.Forwarder, error) {
	return forward.NewNATS(forward.NATSConfig{URL: f.URL, Subject: f.Subject})
}
=== go.mod head ===
module github.com/astrate-platform/astrate

go 1.26.1

// Pinned dependency decisions — docs/ROADMAP.md §1.1 (frozen in docs/DESIGN.md).
=== Makefile targets ===
42:	$(GO) test -race -count=1 -tags integration ./...
46:	$(GO) test -race -count=1 -tags "integration e2e" ./...

$ cd /root/astrate-mule && time go build -tags nats ./cmd/... 2>&1 | tail -20; echo "exit=$?"; echo "=== vet tagged ==="; go vet -tags nats ./cmd/... 2>&1 | tail -20; echo "exit=$?"
/bin/sh: 1: time: not found
exit=0
=== vet tagged ===
exit=0

$ cd /root/astrate-mule && go build -tags nats ./cmd/... ; echo "build exit=$?"; go vet -tags nats ./cmd/... ; echo "vet exit=$?"; go test -tags nats -count=1 ./cmd/... 2>&1 | tail -5
build exit=0
vet exit=0
--- FAIL: TestForwarderNATSBuilt (0.00s)
    forward_nats_test.go:22: failed to start NATS container: run nats: generic container: get provider: rootless Docker not found, failed to create Docker provider
FAIL
FAIL	github.com/astrate-platform/astrate/cmd/astrate	0.467s
FAIL

→ Read cmd/astrate/main_test.go

$ cd /root/astrate-mule && echo "=== MarkRunningFlowsStopped ==="; rg -n -A18 'func \(s \*Service\) MarkRunningFlowsStopped' internal/flowapi/*.go; echo "=== flow Manager.Shutdown ==="; rg -n -A20 'func \(m \*Manager\) Shutdown' internal/flow/*.go
=== MarkRunningFlowsStopped ===
internal/flowapi/service.go:885:func (s *Service) MarkRunningFlowsStopped(ctx context.Context) {
internal/flowapi/service.go-886-	now := time.Now().UTC()
internal/flowapi/service.go-887-	for _, f := range s.mgr.ListFlows() {
internal/flowapi/service.go-888-		if f.Status() != flow.FlowStatusRunning && f.Status() != flow.FlowStatusStopped {
internal/flowapi/service.go-889-			continue
internal/flowapi/service.go-890-		}
internal/flowapi/service.go-891-		realm, name, ok := splitInstanceID(f.PipelineID())
internal/flowapi/service.go-892-		if !ok {
internal/flowapi/service.go-893-			continue
internal/flowapi/service.go-894-		}
internal/flowapi/service.go-895-		id, err := s.realmID(ctx, realm)
internal/flowapi/service.go-896-		if err != nil {
internal/flowapi/service.go-897-			continue
internal/flowapi/service.go-898-		}
internal/flowapi/service.go-899-		_ = s.st.UpdateFlowRuntime(ctx, id, name, "stopped", nil, nil, nil, &now)
internal/flowapi/service.go-900-	}
internal/flowapi/service.go-901-}
internal/flowapi/service.go-902-
internal/flowapi/service.go-903-func (s *Service) liveFlow(realm, name string) *flow.Flow {
=== flow Manager.Shutdown ===
internal/flow/flow.go:312:func (m *Manager) Shutdown(ctx context.Context) error {
internal/flow/flow.go-313-	m.mu.RLock()
internal/flow/flow.go-314-	flows := make([]*Flow, 0, len(m.flows))
internal/flow/flow.go-315-	for _, f := range m.flows {
internal/flow/flow.go-316-		flows = append(flows, f)
internal/flow/flow.go-317-	}
internal/flow/flow.go-318-	m.mu.RUnlock()
internal/flow/flow.go-319-
internal/flow/flow.go-320-	var first error
internal/flow/flow.go-321-	for _, f := range flows {
internal/flow/flow.go-322-		f.mu.RLock()
internal/flow/flow.go-323-		s := f.status
internal/flow/flow.go-324-		f.mu.RUnlock()
internal/flow/flow.go-325-		if s == FlowStatusRunning {
internal/flow/flow.go-326-			if err := m.StopFlow(ctx, f.PipelineID()); err != nil && first == nil {
internal/flow/flow.go-327-				first = err
internal/flow/flow.go-328-			}
internal/flow/flow.go-329-		}
internal/flow/flow.go-330-	}
internal/flow/flow.go-331-	return first
internal/flow/flow.go-332-}

$ cd /root/astrate-mule && rg -n -B3 -A12 'shutdown|Shutdown|drain' docs/DESIGN.md | rg -n -A12 -B3 '5\.3' | head -60
16-107-│   ├── store/                # pgx repositories, migrations (golang-migrate embedded)
17---
18-790-
19:791-### 5.3 Lifecycle & resilience
20-792-
21-793:- Graceful shutdown: broker stops accepting, shards drain (bounded by timeout), batches flush.
22-794-- Crash safety: QoS ≥ 1 messages are PUBACK'd **only after** the persistence batch commits —
23-795-  at-least-once into Postgres; datastream inserts are idempotence-tolerant (duplicate
24-796-  (series, ts) rows are acceptable per Astarte semantics; properties are upserts).
25-797-- DB outage: shards park with exponential backoff; broker applies backpressure (§1.4); QoS 0
26-798-  data degrades first, by design.
27-799-
28-800-### 5.4 Deployment
29-801-
30-802-`docker-compose.yml`: `timescale/timescaledb:latest-pg16` (tuned: `shared_buffers=256MB`,
31-803-`max_connections=50`) + `astrate` (distroless, ~20 MB image). Volumes: pgdata + session store.

$ cd /root/astrate-mule && echo "=== forward_test.go ==="; cat cmd/astrate/forward_test.go; echo "=== version_test.go ==="; cat cmd/astrate/version_test.go; echo "=== engine.New sig ==="; rg -n 'func New\(' internal/engine/engine.go; echo "=== log level validation ==="; rg -n 'Level|Format' internal/config/config.go | rg -n 'valid|switch|case|Level' | head -20
=== forward_test.go ===
package main

import (
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/astrate-platform/astrate/internal/config"
)

// quietLogger keeps newForwarder's info line out of the test output.
func quietLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// TestForwarderDisabledIsUntypedNil is the one rule in newForwarder that a
// compiler cannot catch and a careless refactor would break: the executor
// decides whether to forward by comparing its Forwarder against nil, so
// returning a typed nil *forward.HTTP would pass that comparison and then
// panic on the first custom trigger action. The interface value itself must
// be nil, which is why the disabled branch returns a bare nil rather than a
// nil pointer.
func TestForwarderDisabledIsUntypedNil(t *testing.T) {
	f, err := newForwarder(config.Config{}, quietLogger())
	if err != nil {
		t.Fatalf("newForwarder: %v", err)
	}
	if f != nil {
		t.Fatalf("forwarder = %#v, want an untyped nil interface", f)
	}
}

// TestForwarderHTTPBuilt checks the configured kind produces a usable
// forwarder, so the disabled case above is not satisfied by a function that
// returns nil for everything.
func TestForwarderHTTPBuilt(t *testing.T) {
	var cfg config.Config
	cfg.Triggers.Forward = config.ForwardConfig{
		Kind: "http",
		URL:  "https://bus.example/trigger",
	}
	f, err := newForwarder(cfg, quietLogger())
	if err != nil {
		t.Fatalf("newForwarder: %v", err)
	}
	if f == nil {
		t.Fatal("forwarder = nil, want an HTTP forwarder")
	}
}

// TestForwarderRejectsBadURL: config.validate normally catches this, but
// newForwarder must not silently ignore a construction failure if a config
// ever reaches it unvalidated.
func TestForwarderRejectsBadURL(t *testing.T) {
	var cfg config.Config
	cfg.Triggers.Forward = config.ForwardConfig{Kind: "http", URL: "nope"}
	if _, err := newForwarder(cfg, quietLogger()); err == nil {
		t.Fatal("expected an error for a relative URL")
	}
}

// TestForwarderNATSWithoutBuildTag: this file (and newnats_default.go) build
// without -tags nats, so kind = "nats" must fail at boot with a message
// naming the missing build tag, not silently forward nothing.
func TestForwarderNATSWithoutBuildTag(t *testing.T) {
	if natsBuildTagEnabled {
		t.Skip("built with -tags nats; see forward_nats_test.go")
	}
	var cfg config.Config
	cfg.Triggers.Forward = config.ForwardConfig{Kind: "nats", URL: "nats://bus.example:4222", Subject: "astrate.triggers"}
	_, err := newForwarder(cfg, quietLogger())
	if err == nil {
		t.Fatal("expected an error building a NATS forwarder without -tags nats")
	}
	if !strings.Contains(err.Error(), "-tags nats") {
		t.Errorf("error = %q, want it to name the missing build tag", err)
	}
}
=== version_test.go ===
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/realm"
)

// TestRealmScopedVersionReportsAPICompatLevel proves the realm-scoped version
// endpoints serve the emulated upstream API level the Dashboard feature-gates
// its UI on (COMPATIBILITY.md deviation 10, issue #77), not Astrate's build
// version. The pairing route is public, so it is reachable without a token.
func TestRealmScopedVersionReportsAPICompatLevel(t *testing.T) {
	mux := http.NewServeMux()
	mountRealmVersion(mux, auth.NewMiddleware(nil))

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/pairing/v1/somerealm/version", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("pairing version status = %d, want 200", rec.Code)
	}
	want := `{"data":"` + realm.APICompatVersion + `"}`
	if got := rec.Body.String(); got != want {
		t.Errorf("pairing version body = %s, want %s", got, want)
	}
}
=== engine.New sig ===
58:func New(st Store, bp BrokerPort, cfg Config) (*Engine, error) {
=== log level validation ===
1:151:	Level  string `toml:"level"`
3:198:		Log: LogConfig{Level: "info", Format: "json"},
4:252:	str("ASTRATE_LOG_LEVEL", &cfg.Log.Level)
6:379:	switch c.Log.Level {
7:382:		return fmt.Errorf("config: log.level %q is not one of debug|info|warn|error", c.Log.Level)
8:384:	switch c.Log.Format {

$ cd /root/astrate-mule && echo "=== store.Close ==="; rg -n -B4 -A12 'func \(s \*Store\) Close' internal/store/*.go; echo "=== Hijack users ==="; rg -n 'Hijack\(' internal/ --glob '!*_test.go'
=== store.Close ===
internal/store/store.go-184-	return nil
internal/store/store.go-185-}
internal/store/store.go-186-
internal/store/store.go-187-// Close releases the connection pool. The Store is unusable afterwards.
internal/store/store.go:188:func (s *Store) Close() { s.pool.Close() }
internal/store/store.go-189-
internal/store/store.go-190-// pgErrCode extracts the PostgreSQL error code from err, or "" if err is not
internal/store/store.go-191-// a server-reported error.
internal/store/store.go-192-func pgErrCode(err error) string {
internal/store/store.go-193-	var pgErr *pgconn.PgError
internal/store/store.go-194-	if errors.As(err, &pgErr) {
internal/store/store.go-195-		return pgErr.Code
internal/store/store.go-196-	}
internal/store/store.go-197-	return ""
internal/store/store.go-198-}
internal/store/store.go-199-
internal/store/store.go-200-const (
=== Hijack users ===

$ cd /root/astrate-mule && rg -n 'websocket' go.mod; echo "=== stream/ws.go accept ==="; rg -n -A25 'func.*ServeWS|Accept\(' internal/appengine/stream/ws.go | head -50
9:	github.com/coder/websocket v1.8.14 // live stream socket (M7)
53:	github.com/gorilla/websocket v1.5.3 // indirect
=== stream/ws.go accept ===
87:	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
88-	if err != nil {
89-		return // Accept already wrote the failure
90-	}
91-	defer func() { _ = conn.CloseNow() }()
92-	// CloseRead drains client frames and gives us a context that cancels when
93-	// the client disconnects (the socket is server-push only).
94-	ctx := conn.CloseRead(r.Context())
95-
96-	for {
97-		select {
98-		case <-ctx.Done():
99-			return
100-		case ev, ok := <-events:
101-			if !ok {
102-				_ = conn.Close(websocket.StatusGoingAway, "bus closed")
103-				return
104-			}
105-			payload, err := json.Marshal(toWire(&ev))
106-			if err != nil {
107-				continue
108-			}
109-			if err := conn.Write(ctx, websocket.MessageText, payload); err != nil {
110-				return
111-			}
112-		}

$ cd /root/astrate-mule && rg -n -A60 'func \(c \*Config\) validate|func \(c Config\) validate' internal/config/config.go | head -90
311:func (c *Config) validate() error {
312-	if c.Database.DSN == "" {
313-		return fmt.Errorf("config: database.dsn is required (or set ASTRATE_DATABASE_DSN)")
314-	}
315-	if c.HTTP.Addr == "" {
316-		return fmt.Errorf("config: http.addr is required")
317-	}
318-	if c.MQTT.SessionStorePath == "" {
319-		return fmt.Errorf("config: mqtt.session_store_path is required")
320-	}
321-	if !c.MQTT.InsecureDevMode && (c.MQTT.TLSCertFile == "" || c.MQTT.TLSKeyFile == "") {
322-		return fmt.Errorf("config: mqtt.tls_cert_file and mqtt.tls_key_file are required unless mqtt.insecure_dev_mode is set")
323-	}
324-	if (c.HTTP.TLSCertFile == "") != (c.HTTP.TLSKeyFile == "") {
325-		return fmt.Errorf("config: http.tls_cert_file and http.tls_key_file must be set together")
326-	}
327-	for _, o := range c.HTTP.CORSAllowedOrigins {
328-		if o == "*" {
329-			continue
330-		}
331-		u, err := url.Parse(o)
332-		if err != nil || u.Scheme == "" || u.Host == "" || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
333-			return fmt.Errorf("config: http.cors_allowed_origins entry %q is not \"*\" or an absolute origin like \"https://host:port\"", o)
334-		}
335-	}
336-	if c.Engine.Shards <= 0 {
337-		return fmt.Errorf("config: engine.shards must be positive")
338-	}
339-	switch c.Triggers.Forward.Kind {
340-	case "", "http", "nats":
341-	default:
342-		return fmt.Errorf("config: triggers.forward.kind %q is not one of \"\"|\"http\"|\"nats\"", c.Triggers.Forward.Kind)
343-	}
344-	// When kind is "" forwarding is disabled; do not validate url, method, or
345-	// static_headers — a stale url left behind after disabling forwarding must
346-	// not prevent the config from loading.
347-	if c.Triggers.Forward.Kind == "http" {
348-		if c.Triggers.Forward.URL == "" {
349-			return fmt.Errorf("config: triggers.forward.url is required when kind is \"http\"")
350-		}
351-		u, err := url.Parse(c.Triggers.Forward.URL)
352-		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
353-			return fmt.Errorf("config: triggers.forward.url %q is not an absolute http or https URL", c.Triggers.Forward.URL)
354-		}
355-		if c.Triggers.Forward.Method != "" {
356-			switch strings.ToUpper(c.Triggers.Forward.Method) {
357-			case "DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT":
358-			default:
359-				return fmt.Errorf("config: triggers.forward.method %q is not a recognised HTTP method", c.Triggers.Forward.Method)
360-			}
361-		}
362-	}
363-	// Deeper validation (whether the URL actually names a reachable NATS
364-	// server) happens at construction time in forward.NewNATS, the same way
365-	// decision 15 has it dial at boot rather than on the first delivery —
366-	// duplicating a URL-shape check here would only drift from what the
367-	// nats.go client actually accepts.
368-	if c.Triggers.Forward.Kind == "nats" {
369-		if c.Triggers.Forward.URL == "" {
370-			return fmt.Errorf("config: triggers.forward.url is required when kind is \"nats\"")
371-		}

$ cd /root/astrate-mule && rg -n -i 'healthcheck|ASTRATE_HTTP_ADDR|http.addr|addr' Dockerfile docker-compose.yml 2>/dev/null | head -30; echo "=== config default http addr ==="; rg -n 'Addr:' internal/config/config.go | head -5
Dockerfile:35:HEALTHCHECK --interval=10s --timeout=3s --start-period=10s --retries=5 \
Dockerfile:36:    CMD ["/astrate", "-healthcheck"]
docker-compose.yml:48:    healthcheck:
=== config default http addr ===
178:		HTTP: HTTPConfig{Addr: ":8080"},
180:			Addr:             ":8883",
181:			DevAddr:          ":1883",

$ cd /root/astrate-mule && sed -n '40,60p' docker-compose.yml
environment:
      POSTGRES_DB: astrate
      POSTGRES_USER: astrate
      POSTGRES_PASSWORD: astrate
    ports:
      - "127.0.0.1:5432:5432"
    volumes:
      - pgdata_ha:/home/postgres/pgdata/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U astrate -d astrate"]
      interval: 5s
      timeout: 3s
      retries: 12

  astrate:
    profiles: ["full"]
    build:
      context: .
    depends_on:
      timescaledb:
        condition: service_healthy

$ cd /root/astrate-mule && rg -n -i 'healthcheck|idletimeout|selfsigned|devcert|drain|shutdownbudget|per-stage' .mule/todo.md .mule/for-giulio.md | head -20; echo "=== for-giulio tail ==="; tail -25 .mule/for-giulio.md; echo "=== open issue check ==="; gh issue list --label mule --state open --limit 20 2>&1 | head -25
.mule/for-giulio.md:282:  closed channel — it panics either way). The fix is that `Drain` no longer closes the lane
.mule/for-giulio.md:285:  exit on `quit` after draining what is buffered. `TestRouter_SubmitParkedWhenDrainRuns`
.mule/todo.md:89:- [ ] flowapi-autorestart-shutdown-cancel: `onBlockFatal` fires `go s.restartWithBackoff` (internal/flowapi/service.go:560) with no stop signal, so a block-death racing process shutdown (cmd/astrate/main.go:482-488 drains the Manager but never cancels the goroutine) can rebuild the flow via `mgr.StartFlow` on a background pump (internal/flow/flow.go:173-175) after `Manager.Shutdown` and flip the durable status back to running. Thread a stop channel/context through the Service, checked in the loop's sleep, cancelled on shutdown; verification needs a `[legion]` integration or timing-based test. [auto]
=== for-giulio tail ===
- ~~golangci-lint is not installed on the Pi~~ — **resolved 2026-07-28**: installed v2.12.2
  via `go install` (the prebuilt-binary installer's published sha256 for linux/arm64 did not
  verify, reproducibly, so built from source instead). `mule.service` now sets
  `MULE_LINT_CMD=golangci-lint run ./...` and has `/root/go/bin` on `PATH`; the lint gate
  runs starting with the next tick.
- ~~`/root/astrate` on the Pi has uncommitted work~~ — **resolved 2026-07-27** with the new
  `tools/reconcile.sh`: rescued onto `origin/wip/DietPi-20260727T171543Z` (pushed, not
  reviewed — read the diff before merging anything from it) and `/root/astrate` is now clean
  on `main`. Also had to set `commit.gpgsign false` locally on that clone first — it had
  signing on with no working gpg-agent for a non-interactive session, same class of problem
  `/root/astrate-mule` already solved the same way (your call, confirmed 2026-07-27: align it
  with the existing unattended-Pi-clone convention rather than fix gpg-agent or bypass signing
  silently inside the script).

## 2026-08-23 — FDO promoted to milestone-4.0 candidate (Giulio's decision, recorded)

#78 is no longer parked: zero-touch onboarding is strategic for commercial
viability. Scope frozen on the issue: owner-side TO1/TO2 in our Pairing
service only (last mile, like upstream), reuse fdo-rs for
manufacturing/rendezvous, acceptance = official `astarte-device-fdo-rust`
SDK completes onboarding against Astrate, docs as a first-class deliverable.
When v3.0 is marked DONE, the v4.0 section of `.mule/milestones.md` should be
drafted with this investigation as its first item (issue #78 has the full
verified context).
- **The mule has been idle 17h.** Filed by the dead-man's switch; see journalctl on the Pi.
=== open issue check ===

$ cd /root/astrate-mule && rg -n -A25 'func \(m \*Manager\) StopFlow' internal/flow/flow.go; echo "=== RehydrateAutoRestart ==="; rg -n -A30 'func \(s \*Service\) RehydrateAutoRestart' internal/flowapi/service.go | head -45
248:func (m *Manager) StopFlow(ctx context.Context, pipelineID string) error {
249-	m.mu.RLock()
250-	f, ok := m.flows[pipelineID]
251-	m.mu.RUnlock()
252-
253-	if !ok {
254-		return fmt.Errorf("%w: %s", ErrFlowNotFound, pipelineID)
255-	}
256-
257-	// 1. Stop producing: cancel pump and wait for Emit loops to exit so no
258-	//    new messages are submitted during drain.
259-	if f.cancelPump != nil {
260-		f.cancelPump()
261-	}
262-	f.pumpWG.Wait()
263-
264-	// 2. Drain in-flight lane work.
265-	if err := f.router.Drain(ctx); err != nil {
266-		return fmt.Errorf("flow drain %s: %w", pipelineID, err)
267-	}
268-
269-	// 3. Release block resources (e.g. AstarteSource bus subscriptions).
270-	for _, b := range f.graph.Blocks() {
271-		if s, ok := b.(Stopper); ok {
272-			s.Stop()
273-		}
=== RehydrateAutoRestart ===
328:func (s *Service) RehydrateAutoRestart(ctx context.Context) error {
329-	rows, err := s.st.ListAutoRestartFlows(ctx)
330-	if err != nil {
331-		return err
332-	}
333-	for _, row := range rows {
334-		var cfgMap map[string]any
335-		if err := json.Unmarshal(row.Config, &cfgMap); err != nil {
336-			msg := fmt.Sprintf("invalid stored config: %v", err)
337-			s.log.Error("flow rehydrate failed", "realm", row.RealmName, "name", row.Name, "error", msg)
338-			_ = s.st.UpdateFlowRuntime(ctx, row.RealmID, row.Name, "failed", &msg, nil, nil, nil)
339-			continue
340-		}
341-		if cfgMap == nil {
342-			cfgMap = map[string]any{}
343-		}
344-		if _, err := s.startFlowInstance(ctx, row.RealmName, row.RealmID, row.Name, row.PipelineName, cfgMap); err != nil {
345-			s.log.Error("flow rehydrate failed", "realm", row.RealmName, "name", row.Name, "error", err)
346-			// startFlowInstance already marked failed when row exists
347-			continue
348-		}
349-		s.log.Info("flow rehydrated", "realm", row.RealmName, "name", row.Name, "pipeline", row.PipelineName)
350-	}
351-	return nil
352-}
353-
354-// resolveAndBuild is the resolve half of every start path: pipeline lookup,
355-// ${config.*} substitution, validation, and block instantiation. Every error
356-// wraps ErrValidation. It has no durable-row side effects; callers own status
357-// transitions.
358-func (s *Service) resolveAndBuild(ctx context.Context, realmID int16, realm, name, pipelineName string, config map[string]any) ([]flow.Block, error) {

$ cd /root/astrate-mule && head -3 internal/engine/forward/nats_test.go internal/engine/forward/nats.go cmd/astrate/forward_nats_test.go; echo "=== exec.Command in tests ==="; rg -ln 'exec.Command' --glob '*_test.go' | head; echo "=== ListAutoRestartFlows where ==="; rg -n -A8 'func \(s \*Store\) ListAutoRestartFlows' internal/store/*.go
==> internal/engine/forward/nats_test.go <==
//go:build nats

package forward

==> internal/engine/forward/nats.go <==
//go:build nats

package forward

==> cmd/astrate/forward_nats_test.go <==
//go:build nats

package main
=== exec.Command in tests ===
test/conformance/cpc/astartectl_test.go
test/conformance/cpa/astartectl_test.go
test/conformance/pysdk/pysdk_test.go
internal/pairing/http_test.go
=== ListAutoRestartFlows where ===
internal/store/flows.go:118:func (s *Store) ListAutoRestartFlows(ctx context.Context) ([]FlowRehydrate, error) {
internal/store/flows.go-119-	rows, err := s.pool.Query(ctx,
internal/store/flows.go-120-		`SELECT f.id, f.realm_id, f.name, f.pipeline_name, f.config, f.auto_restart,
internal/store/flows.go-121-		        f.status, f.error_message, f.failed_block, f.created_at, f.updated_at, f.started_at, f.stopped_at,
internal/store/flows.go-122-		        r.name
internal/store/flows.go-123-		 FROM flows f
internal/store/flows.go-124-		 JOIN realms r ON r.id = f.realm_id
internal/store/flows.go-125-		 WHERE f.auto_restart = true
internal/store/flows.go-126-		 ORDER BY r.name, f.name`)

$ cd /root/astrate-mule && rg -n 'auto_restart' internal/store/flows.go | head -20; echo "=== flowapi stop ==="; rg -n -B4 -A22 'func \(s \*Service\) Stop\b|func \(s \*Service\) Stop\(' internal/flowapi/service.go | head -50
33:// FlowRehydrate is a durable auto_restart row plus its realm name for boot.
39:const flowColumns = `id, realm_id, name, pipeline_name, config, auto_restart,
63:		`INSERT INTO flows (realm_id, name, pipeline_name, config, auto_restart)
116:// ListAutoRestartFlows returns every flow with auto_restart=true across all
120:		`SELECT f.id, f.realm_id, f.name, f.pipeline_name, f.config, f.auto_restart,
125:		 WHERE f.auto_restart = true
=== flowapi stop ===

→ Read Makefile

$ cd /root/astrate-mule && rg -n 'auto_restart|AutoRestart' internal/flowapi/*.go internal/store/flows.go | rg -v '_test.go' | head -30
internal/store/flows.go:21:	AutoRestart  bool
internal/store/flows.go:33:// FlowRehydrate is a durable auto_restart row plus its realm name for boot.
internal/store/flows.go:39:const flowColumns = `id, realm_id, name, pipeline_name, config, auto_restart,
internal/store/flows.go:47:		&f.ID, &f.RealmID, &f.Name, &f.PipelineName, &f.Config, &f.AutoRestart,
internal/store/flows.go:63:		`INSERT INTO flows (realm_id, name, pipeline_name, config, auto_restart)
internal/store/flows.go:116:// ListAutoRestartFlows returns every flow with auto_restart=true across all
internal/store/flows.go:118:func (s *Store) ListAutoRestartFlows(ctx context.Context) ([]FlowRehydrate, error) {
internal/store/flows.go:120:		`SELECT f.id, f.realm_id, f.name, f.pipeline_name, f.config, f.auto_restart,
internal/store/flows.go:125:		 WHERE f.auto_restart = true
internal/store/flows.go:136:			&fr.ID, &fr.RealmID, &fr.Name, &fr.PipelineName, &fr.Config, &fr.AutoRestart,
internal/flowapi/http.go:63:	AutoRestart *bool           `json:"auto_restart"`
internal/flowapi/http.go:141:	autoRestart := resolveAutoRestart(body.AutoRestart)
internal/flowapi/http.go:146:		AutoRestart: autoRestart,
internal/flowapi/http.go:294:// resolveAutoRestart returns the effective auto_restart value from a request
internal/flowapi/http.go:296:func resolveAutoRestart(p *bool) bool {
internal/flowapi/service.go:78:	AutoRestart  bool            `json:"auto_restart"`
internal/flowapi/service.go:94:	AutoRestart bool
internal/flowapi/service.go:318:	row, err := s.st.CreateFlow(ctx, id, req.Name, req.Pipeline, []byte(config), req.AutoRestart)
internal/flowapi/service.go:325:// RehydrateAutoRestart starts every durable flow with auto_restart=true.
internal/flowapi/service.go:328:func (s *Service) RehydrateAutoRestart(ctx context.Context) error {
internal/flowapi/service.go:329:	rows, err := s.st.ListAutoRestartFlows(ctx)
internal/flowapi/service.go:441:			AutoRestart: true, Status: "running", StartedAt: &now,
internal/flowapi/service.go:523:// until the flow starts or its durable row is deleted / auto_restart disabled.
internal/flowapi/service.go:531:// graph down, and — for auto_restart flows — schedules a backoff rebuild.
internal/flowapi/service.go:554:		if !row.AutoRestart {
internal/flowapi/service.go:595:// the durable row disappears (flow deleted), auto_restart is turned off, or
internal/flowapi/service.go:617:		if !row.AutoRestart {
internal/flowapi/service.go:950:		AutoRestart:  row.AutoRestart,

$ cd /root/astrate-mule && rg -n -A25 '^\s*add\)|^add\(\)|add\)' tools/mule.sh | head -60
973:  add)       shift; cmd_add "$@";;
974-  menu)      shift; cmd_menu "$@";;
975-  recipe)    shift; cmd_recipe "$@";;
976-  refill)    shift; cmd_refill "$@";;
977-  review)    shift; cmd_review "$@";;
978-  budget)    shift; cmd_budget "$@";;
979-  legion)    shift; cmd_legion "$@";;
980-  tick)      shift; cmd_tick "$@";;
981-  revert)    shift; cmd_revert "$@";;
982-  reports)   shift; cmd_reports "$@";;
983-  *) sed -n '2,19p' "$0" | sed 's/^# \{0,1\}//'; exit 1;;
984-esac

$ cd /root/astrate-mule && rg -n -A20 'cmd_add\(\)' tools/mule.sh
421:cmd_add() {
422-  [ $# -gt 0 ] || die 'add needs a title: mule.sh add "slug: what to do"'
423-  printf -- '- [ ] %s\n' "$*" >> "$TODO"
424-  ok "queued: $*"
425-}
426-
427-# Full-transcript history written by save_report — filenames sort chronologically, so `ls`
428-# is the ordering. No args: list. A number: print that many full transcripts, newest first.
429-cmd_reports() {
430-  local dir="$MULE/reports/log" n="${1:-}"
431-  [ -d "$dir" ] || { note "no reports yet"; return 0; }
432-  if [ -z "$n" ]; then
433-    ls -1 "$dir" | sort | tac
434-    note "mule.sh reports N   — print the N most recent transcripts in full"
435-    return 0
436-  fi
437-  ls -1 "$dir" | sort | tac | head -n "$n" | while read -r f; do
438-    echo "===== $f ====="
439-    cat "$dir/$f"
440-    echo
441-  done

$ cd /root/astrate-mule && rg -n '^## ' .mule/for-giulio.md | tail -6; echo "---- total lines ----"; wc -l .mule/for-giulio.md
354:## 2026-08-23 — FDO promoted to milestone-4.0 candidate (Giulio's decision, recorded)
---- total lines ----
364 .mule/for-giulio.md

→ Read .mule/for-giulio.md [offset=350]
Now writing the review file.

← Write .mule/reviews/cmd-astrate-2026-09-25.md
Wrote file successfully.

$ cd /root/astrate-mule && rg -n 'MasterKey|EnvMasterKey' internal/config/config.go internal/store/*.go | rg -v '_test' | head -20
internal/config/config.go:134:// private keys (§4.3). When MasterKeyFile is empty, the store falls back to
internal/config/config.go:137:	MasterKeyFile string `toml:"master_key_file"`
internal/config/config.go:248:	str("ASTRATE_SECURITY_MASTER_KEY_FILE", &cfg.Security.MasterKeyFile)
internal/store/crypto.go:22:	// EnvMasterKey names the env var holding the master key itself,
internal/store/crypto.go:24:	EnvMasterKey = "ASTRATE_MASTER_KEY"
internal/store/crypto.go:25:	// EnvMasterKeyFile names the env var holding a path to a file that
internal/store/crypto.go:27:	EnvMasterKeyFile = "ASTRATE_MASTER_KEY_FILE"
internal/store/crypto.go:28:	// MasterKeySize is the required decoded master key length (AES-256).
internal/store/crypto.go:29:	MasterKeySize = 32
internal/store/crypto.go:32:// ErrNoMasterKey reports that neither master-key env reference is set.
internal/store/crypto.go:33:var ErrNoMasterKey = fmt.Errorf("store: neither %s nor %s is set", EnvMasterKey, EnvMasterKeyFile)
internal/store/crypto.go:44:	if len(masterKey) != MasterKeySize {
internal/store/crypto.go:45:		return nil, fmt.Errorf("store: master key must be %d bytes, got %d", MasterKeySize, len(masterKey))
internal/store/crypto.go:59:// (EnvMasterKey first, then EnvMasterKeyFile).
internal/store/crypto.go:61:	key, err := LoadMasterKey()
internal/store/crypto.go:68:// LoadMasterKey resolves the master key from the environment: EnvMasterKey
internal/store/crypto.go:69:// (hex or base64 text) wins over EnvMasterKeyFile (raw 32 bytes, or hex or
internal/store/crypto.go:71:func LoadMasterKey() ([]byte, error) {
internal/store/crypto.go:72:	if v := os.Getenv(EnvMasterKey); v != "" {
internal/store/crypto.go:73:		key, err := decodeMasterKey(v)

← Write /tmp/opencode/queue-lines.txt
Wrote file successfully.

$ cd /root/astrate-mule && while IFS= read -r line; do [ -n "$line" ] && tools/mule.sh add "$line"; done < /tmp/opencode/queue-lines.txt
ok queued: flow-boot-resumes-stopped-durable-flows [legion] [auto]: RehydrateAutoRestart (cmd/astrate/main.go:216) restarts every flow with auto_restart=true whatever its last status — `ListAutoRestartFlows` filters on `f.auto_restart = true` alone (internal/store/flows.go:125) and nothing ever clears the flag, which is written once at create (flowapi/http.go:141-146 -> service.go:318) and never on stop — so a durable flow stopped through the API is running again after the next boot, while the shutdown mark that says otherwise (MarkRunningFlowsStopped, main.go:497 -> internal/flowapi/service.go:885-901, which writes status="stopped") is never read by that query. The two halves disagree: service.go:325 documents "every durable flow", the shutdown call assumes "only what was running". Say which is the intent, then make them agree — either add `AND f.status = 'running'` to flows.go:125 (making the shutdown mark load-bearing) or delete the mark as dead code — and add a case in internal/store/flows_test.go pinning the chosen rule (a stopped auto_restart row must, or must not, come back). Needs the DB.
ok queued: drain-per-stage-budget [auto]: `shutdown` spends a single 30s sctx (cmd/astrate/main.go:479) across three sequential stages — srv.Shutdown (484), b.Close (489), flowSvc.Manager().Shutdown (494) and MarkRunningFlowsStopped (497) — while `drainEngine` alone opens a second one (504), so the constant's own comment "bounds the whole graceful drain" (main.go:55-56) is wrong in both directions, and one slow in-flight HTTP request (the server sets no per-request deadline, 221) burns the whole first budget and leaves the flow stage on an expired context: StopFlow then returns at `f.router.Drain(ctx)` before step 3 releases block resources (internal/flow/flow.go:265-273) and MarkRunningFlowsStopped skips every flow (flowapi/service.go:895-897), both surfacing as one log.Warn. Give each stage its own budget through an extracted helper (shutdownTimeout has to become an injectable var for the test) and pin it container-free: a second stage still receives a live context after a first stage exhausts its own.
ok queued: cmd-devcert-fields-test [auto]: add a container-free test for `selfSignedDevCert` (cmd/astrate/main.go:320-345), the throwaway mTLS identity `insecure_dev_mode` boots on, which no test touches — pin NotBefore = now-1h (331; a plain now would reject a device whose clock is behind), NotAfter = now+365d (332), DNSNames localhost and IPAddresses 127.0.0.1/::1 (336-337, dropping which breaks hostname verification for a device dialling mqtts://127.0.0.1:8883), ExtKeyUsage serverAuth (334), and that the leaf is self-signed and parses (x509.ParseCertificate on Certificate[0]).
ok queued: cmd-healthcheck-contract-test [auto]: add a container-free test for `runHealthcheck` (cmd/astrate/main.go:531-558) — the `HEALTHCHECK CMD ["/astrate", "-healthcheck"]` contract (Dockerfile:35-36) has no coverage. Assert it probes /astrate/v1/readiness and not /health (543), returns 0 only on 200 so a 503 readiness exits 1 (554), and maps a bare ":PORT" in ASTRATE_HTTP_ADDR to 127.0.0.1 (536-538). httptest.NewServer plus t.Setenv("ASTRATE_HTTP_ADDR", srv.Listener.Addr().String()).
ok queued: cmd-loadsealer-masterkeyfile-test [auto]: add a container-free test for `loadSealer` (cmd/astrate/main.go:515-526) — the `security.master_key_file` to env hand-off (`os.Setenv(store.EnvMasterKeyFile, ...)`, 517) is the only wiring that makes a configured key file reach the store loader, and dropping or misspelling that branch leaves the file silently ignored and boot dying with the generic "master key" error (523) even though the config named a valid file. Assert a 32-byte key file under t.TempDir() loads, the inline ASTRATE_MASTER_KEY path still wins, and a missing key returns the error naming all three env vars. t.Setenv both ASTRATE_MASTER_KEY and ASTRATE_MASTER_KEY_FILE first so the os.Setenv is restored at cleanup.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -361,4 +361,24 @@
 When v3.0 is marked DONE, the v4.0 section of `.mule/milestones.md` should be
 drafted with this investigation as its first item (issue #78 has the full
 verified context).
 - **The mule has been idle 17h.** Filed by the dead-man's switch; see journalctl on the Pi.
+
+## 2026-09-25 — the `-tags nats` half of the tree is compiled by no gate
+
+`cmd/astrate/newnats_nats.go`, `internal/engine/forward/nats.go` and
+`internal/engine/forward/nats_test.go` are all behind `//go:build nats`, and nothing
+builds or tests that tag: not `make build`/`make test`/`make test-integration`/
+`make test-e2e` (Makefile:23, 38, 42, 46), not the mule gate, not the lint target.
+So the whole NATS trigger-forwarding path is invisible to CI — while
+`config.validate` accepts `triggers.forward.kind = "nats"` (internal/config/config.go:339-343)
+and the untagged binary then fails at boot with the build-tag message
+(cmd/astrate/newnats_default.go). I measured it today: `go build -tags nats ./cmd/...`
+and `go vet -tags nats ./cmd/...` both exit 0, so there is nothing broken — just nothing
+keeping it unbroken. Worth a `make build-nats` (and the NATS test tier, which needs a
+container, so it would be a `[legion]` line) in the gate; not a queue line because the
+only test I can imagine for it shells out to `go build`, and the mule's own check strips
+the implementation and the guard together, so it would never fail.
+Not queued for the same reason: `IdleTimeout` is missing on the HTTP server
+(cmd/astrate/main.go:221) and `-healthcheck` ignores the config file
+(cmd/astrate/main.go:531) — both real, both yours to call, neither has a test worth
+writing. See `.mule/reviews/cmd-astrate-2026-09-25.md`.

$ cd /root/astrate-mule && git status --short && echo "--- diff stat ---" && git diff --stat && echo "--- todo tail ---" && tail -c 600 .mule/todo.md | tail -3 | cut -c1-160
M .mule/for-giulio.md
 M .mule/todo.md
?? .mule/reviews/cmd-astrate-2026-09-25.md
--- diff stat ---
 .mule/for-giulio.md | 20 ++++++++++++++++++++
 .mule/todo.md       |  5 +++++
 2 files changed, 25 insertions(+)
--- todo tail ---
.master_key_file` to env hand-off (`os.Setenv(store.EnvMasterKeyFile, ...)`, 517) is the only wiring that makes a configured key file reach the store loader, an

Done: reviewed `cmd/astrate` (never the area of a prior pass) and queued 5 task lines — the boot/durable-status mismatch, the shared drain budget, and three unpinned container-free rules.
Files: `.mule/reviews/cmd-astrate-2026-09-25.md`, `.mule/todo.md` (5 lines appended), `.mule/for-giulio.md` (one dated section)
Verified: `go build -tags nats ./cmd/...` and `go vet -tags nats ./cmd/...` -> pass (both exit 0; the NATS half compiles, nothing to fix there)
Unsure: two intents I refused to decide for him — whether boot should resume only flows that were `running` (`internal/store/flows.go:125` vs the shutdown mark at `cmd/astrate/main.go:497`) and which drain stage deserves its own 30s (DESIGN §5.3 gives no number); both are stated as explicit either/or in the task lines
Follow-ups: `flow-boot-resumes-stopped-durable-flows`, `drain-per-stage-budget`, `cmd-devcert-fields-test`, `cmd-healthcheck-contract-test`, `cmd-loadsealer-masterkeyfile-test`
