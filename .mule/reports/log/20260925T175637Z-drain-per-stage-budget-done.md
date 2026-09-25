slug: drain-per-stage-budget
verdict: done
at:  84a050e
ran: 2026-09-25T17:56:37Z on DietPi in 265s


> build · big-pickle

I'll start by reading the two files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ sed -n '40,70p' cmd/astrate/main.go
"github.com/astrate-platform/astrate/internal/flowapi"
	"github.com/astrate-platform/astrate/internal/housekeeping"
	"github.com/astrate-platform/astrate/internal/httpx"
	"github.com/astrate-platform/astrate/internal/observability"
	"github.com/astrate-platform/astrate/internal/pairing"
	"github.com/astrate-platform/astrate/internal/realm"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/swagger"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// version is the reported build version (override with
// -ldflags "-X main.version=vX.Y.Z").
var version = "0.1.0-dev"

// shutdownTimeout bounds the whole graceful drain (docs/DESIGN.md §5.3).
const shutdownTimeout = 30 * time.Second

func main() {
	configPath := flag.String("config", "", "path to the TOML config file (env-only when empty)")
	showVersion := flag.Bool("version", false, "print the version and exit")
	healthcheck := flag.Bool("healthcheck", false, "probe the local readiness endpoint and exit (for container HEALTHCHECK)")
	flag.Parse()

	if *showVersion {
		fmt.Println("astrate", version)
		return
	}
	if *healthcheck {
		os.Exit(runHealthcheck())
	}

$ sed -n '460,530p' cmd/astrate/main.go
return nil
	} else if !errors.Is(err, store.ErrNotFound) {
		return fmt.Errorf("checking auto-provision realm: %w", err)
	}
	key, err := cfg.RealmJWTPublicKey()
	if err != nil {
		return err
	}
	if _, err := hk.CreateRealm(ctx, cfg.Realm.Name, key, cfg.Realm.DeviceRegistrationLimit, nil); err != nil {
		return fmt.Errorf("auto-provisioning realm %q: %w", cfg.Realm.Name, err)
	}
	log.Info("auto-provisioned realm", "realm", cfg.Realm.Name)
	return nil
}

// shutdown drains the stack in the §5.3 order. srv may be nil on a startup
// error before the HTTP server began serving. flowSvc may be nil when Flow
// was never wired (should not happen in run).
func shutdown(srv *http.Server, b *broker.Broker, flowSvc *flowapi.Service, e *engine.Engine, log *slog.Logger) {
	sctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if srv != nil {
		log.Info("draining http listener")
		if err := srv.Shutdown(sctx); err != nil {
			log.Warn("http shutdown", "error", err)
		}
	}
	log.Info("stopping broker")
	if err := b.Close(); err != nil {
		log.Warn("broker close", "error", err)
	}
	if flowSvc != nil {
		log.Info("stopping flows")
		if err := flowSvc.Manager().Shutdown(sctx); err != nil {
			log.Warn("flow shutdown", "error", err)
		}
		flowSvc.MarkRunningFlowsStopped(sctx)
	}
	drainEngine(e, log)
	log.Info("shutdown complete")
}

func drainEngine(e *engine.Engine, log *slog.Logger) {
	dctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	log.Info("draining engine")
	if err := e.Drain(dctx); err != nil {
		log.Warn("engine drain", "error", err)
	}
}

// loadSealer builds the CA-key sealer: the configured master-key file feeds the
// store's env-based loader (ASTRATE_MASTER_KEY[_FILE]); decoding rules are
// shared with internal/store.
func loadSealer(cfg config.Config) (*store.KeySealer, error) {
	if cfg.Security.MasterKeyFile != "" {
		if err := os.Setenv(store.EnvMasterKeyFile, cfg.Security.MasterKeyFile); err != nil {
			return nil, err
		}
	}
	sealer, err := store.NewKeySealerFromEnv()
	if err != nil {
		return nil, fmt.Errorf("master key: %w (set ASTRATE_MASTER_KEY, ASTRATE_MASTER_KEY_FILE, or security.master_key_file)", err)
	}
	return sealer, nil
}

// runHealthcheck probes the local readiness endpoint and returns a process
// exit code, so a distroless container (no shell or curl) can self-check via
// `astrate -healthcheck`. The HTTP address comes from ASTRATE_HTTP_ADDR.

$ sed -n '200,240p' cmd/astrate/main.go
select {
			case <-ctx.Done():
				return
			case <-t.C:
			}
		}
	}()

	if cfg.Realm.Name != "" {
		if err := autoProvisionRealm(ctx, st, hkSvc, cfg, log); err != nil {
			shutdown(nil, b, flowSvc, e, log)
			return err
		}
	}

	// Rehydrate durable auto_restart flows before accepting traffic (Design A / #41).
	if err := flowSvc.RehydrateAutoRestart(ctx); err != nil {
		shutdown(nil, b, flowSvc, e, log)
		return fmt.Errorf("flow rehydrate: %w", err)
	}

	srv := &http.Server{Addr: cfg.HTTP.Addr, Handler: handler, ReadHeaderTimeout: 10 * time.Second}
	serveErr := make(chan error, 1)
	go func() {
		log.Info("astrate listening", "http", cfg.HTTP.Addr, "mqtt", b.TLSAddr(), "version", version)
		if cfg.HTTP.TLSCertFile != "" {
			serveErr <- srv.ListenAndServeTLS(cfg.HTTP.TLSCertFile, cfg.HTTP.TLSKeyFile)
		} else {
			serveErr <- srv.ListenAndServe()
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("shutdown signal received")
	case err := <-serveErr:
		shutdown(nil, b, flowSvc, e, log)
		return fmt.Errorf("http server: %w", err)
	}

	shutdown(srv, b, flowSvc, e, log)

$ ls cmd/astrate/ && rg -n "shutdown|drainEngine" --glob '*.go' -l
forward_nats_test.go
forward_test.go
main.go
main_test.go
newnats_default.go
newnats_nats.go
version_test.go
cmd/astrate/main.go
cmd/astrate/main_test.go
internal/flow/flow.go
internal/flow/manager_test.go
internal/flow/blocks/container/block.go
internal/engine/router_test.go
internal/engine/data.go
internal/engine/triggers/actions.go
internal/engine/router.go
internal/engine/triggers/actions_test.go
internal/engine/batch.go
internal/engine/introspection.go
internal/flowapi/service.go
internal/broker/broker.go
internal/appengine/channels/room.go

$ rg -n "shutdown|drainEngine|func Test" cmd/astrate/main_test.go | head -50
90:func TestBoot(t *testing.T) {
167:	// Graceful shutdown: run returns cleanly within the drain budget.
172:			t.Fatalf("run returned error on shutdown: %v", err)
174:	case <-time.After(shutdownTimeout + 10*time.Second):
193:		case <-time.After(shutdownTimeout + 10*time.Second):

$ rg -n "func (.*Manager) (Shutdown|Drain)|func .*MarkRunningFlowsStopped" internal/flow internal/flowapi internal/engine -g '*.go'
internal/flowapi/service.go:885:func (s *Service) MarkRunningFlowsStopped(ctx context.Context) {

$ rg -n "func (s \*Service) Manager|type Service struct" -A 12 internal/flowapi/service.go | head -40; rg -n "Manager interface|manager " internal/flowapi/service.go | head -20
29:type Service struct {
30-	st       *store.Store
31-	mgr      *flow.Manager
32-	reg      *flow.Registry
33-	bus      *stream.Bus
34-	ingest   flow.IngestFunc
35-	register flow.RegisterFunc
36-	log      *slog.Logger
37-
38-	// restartMu/restarting dedupe auto-restart loops when several blocks of
39-	// the same flow die (or fatal races a manual reload).
40-	restartMu  sync.Mutex
41-	restarting map[string]bool
56:// Manager exposes the flow manager (process shutdown).
588:// block), so no number of retries can succeed. Store/manager errors are

$ rg -n "func New|func (\(b \*Broker\)) Close" internal/broker/broker.go | head; rg -n "func New|func (e \*Engine) Drain" internal/engine/engine.go | head -20
88:func New(ctx context.Context, cfg Config, st Store, intake Intake, sink LifecycleSink) (*Broker, error) {
211:func (b *Broker) Close() error {
58:func New(st Store, bp BrokerPort, cfg Config) (*Engine, error) {

$ sed -n '54,62p' internal/flowapi/service.go; sed -n '880,905p' internal/flowapi/service.go; rg -n "func (m \*Manager) Shutdown" -A 3 internal/flow/manager.go; rg -n "func (e \*Engine) Drain" -A 3 internal/engine/*.go | head -20
}

// Manager exposes the flow manager (process shutdown).
func (s *Service) Manager() *flow.Manager { return s.mgr }

// PipelineView is the operator JSON shape for one stored pipeline.
type PipelineView struct {
	Name       string          `json:"name"`
	Definition json.RawMessage `json:"definition"`
	return out, nil
}

// MarkRunningFlowsStopped best-effort sets durable status to stopped for
// instances that were running (clean process shutdown).
func (s *Service) MarkRunningFlowsStopped(ctx context.Context) {
	now := time.Now().UTC()
	for _, f := range s.mgr.ListFlows() {
		if f.Status() != flow.FlowStatusRunning && f.Status() != flow.FlowStatusStopped {
			continue
		}
		realm, name, ok := splitInstanceID(f.PipelineID())
		if !ok {
			continue
		}
		id, err := s.realmID(ctx, realm)
		if err != nil {
			continue
		}
		_ = s.st.UpdateFlowRuntime(ctx, id, name, "stopped", nil, nil, nil, &now)
	}
}

func (s *Service) liveFlow(realm, name string) *flow.Flow {
	instanceID := flow.InstanceID(realm, name)
	for _, f := range s.mgr.ListFlows() {
rg: internal/flow/manager.go: IO error for operation on internal/flow/manager.go: No such file or directory (os error 2)

$ sed -n '1,40p' cmd/astrate/main_test.go; sed -n '160,200p' cmd/astrate/main_test.go
//go:build integration

package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"

	"github.com/astrate-platform/astrate/internal/config"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

const composeDSN = "postgres://astrate:astrate@127.0.0.1:5432/astrate?sslmode=disable"

// dialDSN finds a reachable database: ASTRATE_TEST_DSN first, then compose.
func dialDSN(t *testing.T) string {
	t.Helper()
	for _, dsn := range []string{envOr("ASTRATE_TEST_DSN", ""), composeDSN} {
		if dsn == "" {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		conn, err := pgx.Connect(ctx, dsn)
		token := mintAgentToken(t, realmKey)
		secret := registerDevice(t, base, realmName, id.String(), token)
		if len(secret) != 44 {
			t.Errorf("credentials secret length = %d, want 44 (%q)", len(secret), secret)
		}
	})

	// Graceful shutdown: run returns cleanly within the drain budget.
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("run returned error on shutdown: %v", err)
		}
	case <-time.After(shutdownTimeout + 10*time.Second):
		t.Fatal("run did not shut down within the drain budget")
	}

	// Restart against the same DB + session store: the DB survives (the realm
	// auto-provision is a no-op) and the bbolt session store reopens.
	t.Run("RestartReusesState", func(t *testing.T) {
		cfg2 := cfg
		cfg2.HTTP.Addr = freeAddr(t)
		ctx2, cancel2 := context.WithCancel(context.Background())
		done2 := make(chan error, 1)
		go func() { done2 <- run(ctx2, cfg2, log) }()
		waitReady(t, "http://"+cfg2.HTTP.Addr)
		cancel2()
		select {
		case err := <-done2:
			if err != nil {
				t.Fatalf("restart run returned error: %v", err)
			}
		case <-time.After(shutdownTimeout + 10*time.Second):
			t.Fatal("restart did not shut down in time")
		}
	})
}

// waitReady polls /readiness until it returns 200.
func waitReady(t *testing.T, base string) {

$ head -20 cmd/astrate/version_test.go; echo ---; rg -n "func (m \*Manager) Shutdown" -A 12 internal/flow/*.go; echo ---; rg -n "func (e \*Engine) Drain" -A 4 internal/engine/engine.go internal/engine/*.go | head -20
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
---
---

$ rg -n "func \(m \*Manager\) Shutdown" -A 15 internal/flow/
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

$ rg -n "func \(e \*Engine\) Drain" -A 6 internal/engine/
internal/engine/engine.go:101:func (e *Engine) Drain(ctx context.Context) error {
internal/engine/engine.go-102-	err := e.drain(ctx)
internal/engine/engine.go-103-
internal/engine/engine.go-104-	bgDone := make(chan struct{})
internal/engine/engine.go-105-	go func() {
internal/engine/engine.go-106-		e.bg.Wait()
internal/engine/engine.go-107-		close(bgDone)

$ rg -n "drainEngine" . 2>/dev/null
./cmd/astrate/main.go:180:		drainEngine(e, log)
./cmd/astrate/main.go:499:	drainEngine(e, log)
./cmd/astrate/main.go:503:func drainEngine(e *engine.Engine, log *slog.Logger) {

$ sed -n '170,190p' cmd/astrate/main.go
if errors.Is(err, pairing.ErrAlreadyRegistered) || errors.Is(err, store.ErrDeviceAlreadyConfirmed) {
				return fmt.Errorf("%w: %s", flow.ErrVirtualDeviceRegistered, deviceID)
			}
			return err
		}, log)

	if err := e.Start(ctx); err != nil {
		return fmt.Errorf("starting engine: %w", err)
	}
	if err := b.Start(); err != nil {
		drainEngine(e, log)
		return fmt.Errorf("starting broker: %w", err)
	}

	handler, hkSvc, err := mountAPIs(cfg, st, e, b, sealer, metrics, flowSvc, pairer, log)
	if err != nil {
		shutdown(nil, b, flowSvc, e, log)
		return err
	}

	// Realm datastream retention ceilings (#72): sweep hourly — and once

$ rg -n "5.3" -A 12 docs/DESIGN.md | rg -n "shutdown|drain|30|order" | head -20
31:793-- Graceful shutdown: broker stops accepting, shards drain (bounded by timeout), batches flush.

$ rg -n "### 5.3" -A 25 docs/DESIGN.md | head -45
791:### 5.3 Lifecycle & resilience
792-
793-- Graceful shutdown: broker stops accepting, shards drain (bounded by timeout), batches flush.
794-- Crash safety: QoS ≥ 1 messages are PUBACK'd **only after** the persistence batch commits —
795-  at-least-once into Postgres; datastream inserts are idempotence-tolerant (duplicate
796-  (series, ts) rows are acceptable per Astarte semantics; properties are upserts).
797-- DB outage: shards park with exponential backoff; broker applies backpressure (§1.4); QoS 0
798-  data degrades first, by design.
799-
800-### 5.4 Deployment
801-
802-`docker-compose.yml`: `timescale/timescaledb:latest-pg16` (tuned: `shared_buffers=256MB`,
803-`max_connections=50`) + `astrate` (distroless, ~20 MB image). Volumes: pgdata + session store.
804-Or: one binary + one `apt install postgresql-16 timescaledb` on a VPS. Target steady-state RSS:
805-Astrate ≤ 150 MB at 1k devices; Postgres tuned to ≤ 768 MB.
806-
807----
808-
809-## 6. Key Risks & Mitigations
810-
811-| Risk | Mitigation |
812-|---|---|
813-| Subtle protocol drift vs. real SDKs | Phase 2 defines CI conformance suites driving the **official** Astarte Go + Python SDKs (and an AtomVM JSON simulator) against Astrate end-to-end. |
814-| Per-interface TTL via DELETE is heavier than `drop_chunks` | Batched chunk-aware deletes off-peak; global chunk-drop cap as the backstop; revisit table-per-retention-class if profiling demands. |
815-| mochi-mqtt session persistence semantics under restart | Dedicated storage-hook tests for `session_present`, offline queue replay order, and QoS 2 exactly-once handshake. |
816-| Embedded CA key compromise | Encrypted at rest, short cert TTLs, latest-serial enforcement, documented re-keying runbook. |

← Edit cmd/astrate/main.go
Index: /root/astrate-mule/cmd/astrate/main.go
===================================================================
--- /root/astrate-mule/cmd/astrate/main.go
+++ /root/astrate-mule/cmd/astrate/main.go
@@ -51,10 +51,13 @@
 // version is the reported build version (override with
 // -ldflags "-X main.version=vX.Y.Z").
 var version = "0.1.0-dev"
 
-// shutdownTimeout bounds the whole graceful drain (docs/DESIGN.md §5.3).
-const shutdownTimeout = 30 * time.Second
+// shutdownTimeout bounds each stage of the graceful drain
+// (docs/DESIGN.md §5.3). It is a var so tests can shrink it; a stage that
+// blows its budget (a slow in-flight HTTP request, say) must not take the
+// stages after it down with an already-expired context.
+var shutdownTimeout = 30 * time.Second
 
 func main() {
 	configPath := flag.String("config", "", "path to the TOML config file (env-only when empty)")
 	showVersion := flag.Bool("version", false, "print the version and exit")

← Edit cmd/astrate/main.go
Index: /root/astrate-mule/cmd/astrate/main.go
===================================================================
--- /root/astrate-mule/cmd/astrate/main.go
+++ /root/astrate-mule/cmd/astrate/main.go
@@ -476,43 +476,48 @@
 }
 
 // shutdown drains the stack in the §5.3 order. srv may be nil on a startup
 // error before the HTTP server began serving. flowSvc may be nil when Flow
-// was never wired (should not happen in run).
+// was never wired (should not happen in run). Every stage gets its own
+// shutdownTimeout budget, so one slow stage only delays itself.
 func shutdown(srv *http.Server, b *broker.Broker, flowSvc *flowapi.Service, e *engine.Engine, log *slog.Logger) {
-	sctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
-	defer cancel()
-
 	if srv != nil {
 		log.Info("draining http listener")
-		if err := srv.Shutdown(sctx); err != nil {
-			log.Warn("http shutdown", "error", err)
-		}
+		stage("http shutdown", log, srv.Shutdown)
 	}
 	log.Info("stopping broker")
 	if err := b.Close(); err != nil {
 		log.Warn("broker close", "error", err)
 	}
 	if flowSvc != nil {
 		log.Info("stopping flows")
-		if err := flowSvc.Manager().Shutdown(sctx); err != nil {
-			log.Warn("flow shutdown", "error", err)
-		}
-		flowSvc.MarkRunningFlowsStopped(sctx)
+		mgr := flowSvc.Manager()
+		stage("flow shutdown", log, mgr.Shutdown)
+		stage("flows marked stopped", log, func(ctx context.Context) error {
+			flowSvc.MarkRunningFlowsStopped(ctx)
+			return nil
+		})
 	}
 	drainEngine(e, log)
 	log.Info("shutdown complete")
 }
 
-func drainEngine(e *engine.Engine, log *slog.Logger) {
-	dctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
+// stage runs one drain stage under a fresh shutdownTimeout context and logs
+// whatever it returns, so an exhausted stage never leaves the next one with a
+// dead context.
+func stage(label string, log *slog.Logger, fn func(context.Context) error) {
+	sctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
 	defer cancel()
-	log.Info("draining engine")
-	if err := e.Drain(dctx); err != nil {
-		log.Warn("engine drain", "error", err)
+	if err := fn(sctx); err != nil {
+		log.Warn(label, "error", err)
 	}
 }
 
+func drainEngine(e *engine.Engine, log *slog.Logger) {
+	log.Info("draining engine")
+	stage("engine drain", log, e.Drain)
+}
+
 // loadSealer builds the CA-key sealer: the configured master-key file feeds the
 // store's env-based loader (ASTRATE_MASTER_KEY[_FILE]); decoding rules are
 // shared with internal/store.
 func loadSealer(cfg config.Config) (*store.KeySealer, error) {

Now the container-free test:
← Write cmd/astrate/shutdown_test.go
Wrote file successfully.

$ gofmt -l . && go vet ./cmd/... && go test ./cmd/astrate/ -run 'TestDrainStage' -v 2>&1 | tail -20
=== RUN   TestDrainStageBudgetIsPerStage
--- PASS: TestDrainStageBudgetIsPerStage (0.12s)
=== RUN   TestDrainStageBoundsItsOwnStage
--- PASS: TestDrainStageBoundsItsOwnStage (0.02s)
PASS
ok  	github.com/astrate-platform/astrate/cmd/astrate	0.342s

$ cp cmd/astrate/main.go /tmp/opencode/main.go.new && python3 - <<'EOF'
import re
p='cmd/astrate/main.go'
s=open(p).read()
old='''func stage(label string, log *slog.Logger, fn func(context.Context) error) {
	sctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := fn(sctx); err != nil {
		log.Warn(label, "error", err)
	}
}'''
new='''var sharedStageCtx context.Context

func stage(label string, log *slog.Logger, fn func(context.Context) error) {
	if sharedStageCtx == nil {
		sharedStageCtx, _ = context.WithTimeout(context.Background(), shutdownTimeout)
	}
	if err := fn(sharedStageCtx); err != nil {
		log.Warn(label, "error", err)
	}
}'''
assert old in s
open(p,'w').write(s.replace(old,new))
EOF
go test ./cmd/astrate/ -run 'TestDrainStage' 2>&1 | tail -12; cp /tmp/opencode/main.go.new cmd/astrate/main.go; gofmt -l cmd/astrate/; go test ./cmd/astrate/ -run 'TestDrainStage' 2>&1 | tail -3
--- FAIL: TestDrainStageBudgetIsPerStage (0.12s)
    shutdown_test.go:42: stage 0 after an exhausted stage got context deadline exceeded, want a live context
    shutdown_test.go:42: stage 1 after an exhausted stage got context deadline exceeded, want a live context
    shutdown_test.go:42: stage 2 after an exhausted stage got context deadline exceeded, want a live context
--- FAIL: TestDrainStageBoundsItsOwnStage (0.00s)
    shutdown_test.go:63: stage deadline 2026-09-25 19:53:57.437273229 +0200 CEST m=+0.436946303 is not within the last 20ms of the budget
FAIL
FAIL	github.com/astrate-platform/astrate/cmd/astrate	0.561s
FAIL
ok  	github.com/astrate-platform/astrate/cmd/astrate	0.589s

$ go vet ./... 2>&1 | tail -5; go test ./... 2>&1 | rg -v "^ok|no test files" | head -20
(no output)

$ go test -race ./cmd/... 2>&1 | rg -v "no test files" | head -10; git status --short
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/cmd/astrate	0.026s
FAIL
 M cmd/astrate/main.go
?? cmd/astrate/shutdown_test.go

Done: per-stage drain budgets — `shutdown` now runs each of the four stages through an extracted `stage` helper that mints a fresh `shutdownTimeout` context, and `shutdownTimeout` became an injectable var; the two container-free tests pin the contract (a stage after an exhausted one gets a live context; each stage's own budget still cuts off a hung stage).

Files: `cmd/astrate/main.go` (var + comment at ~53, `shutdown`/`stage`/`drainEngine` at ~478-512), `cmd/astrate/shutdown_test.go` (new)

Verified: `gofmt -l .` -> clean; `go vet ./...` -> pass; `go test ./...` -> pass (all packages ok); `go test -race ./cmd/astrate/` -> could not run, `ThreadSanitizer: unsupported VMA range (39-bit VMA kernel)`, same limit the Pi has, so the gate was the non-race form. Negative check: with `stage` temporarily rewritten to share one context, both new tests fail (`got context deadline exceeded, want a live context`), then pass again once restored — the file was restored byte-for-byte and re-gofmt'd.

Unsure: the tests pin the `stage` helper's contract, not `shutdown`'s wiring, because a real end-to-end `shutdown` test needs a DB-backed `flowapi.Service`/`engine.Engine` to observe the context their constructors would receive. A future refactor that reintroduced one shared context *inside* `shutdown` while keeping `stage` would not be caught. I also did not touch the per-request deadline the task mentions at main.go:221 (only `ReadHeaderTimeout` exists) — adding one is a behaviour change beyond this task's scope.

Follow-ups: none.
