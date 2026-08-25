// Command astrate is the single Astrate binary (docs/ROADMAP.md §9 file 8.4):
// it wires the store, ingestion engine, embedded MQTT broker, pairing, the M7
// REST surfaces, and the observability endpoints into one process driven by
// one TOML config, with signal-driven graceful shutdown.
package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/astrate-platform/astrate/internal/appengine"
	"github.com/astrate-platform/astrate/internal/appengine/channels"
	apstream "github.com/astrate-platform/astrate/internal/appengine/stream"
	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/config"
	"github.com/astrate-platform/astrate/internal/engine"
	"github.com/astrate-platform/astrate/internal/engine/forward"
	"github.com/astrate-platform/astrate/internal/engine/triggers"
	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
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

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	log := newLogger(cfg.Log)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cfg, log); err != nil {
		log.Error("astrate exited with error", "error", err)
		os.Exit(1)
	}
}

// run assembles and serves the whole stack until ctx is cancelled, then drains
// it in the §5.3 order: HTTP → broker → flow manager → engine → store. It is
// the in-process entry point the boot tests drive.
func run(ctx context.Context, cfg config.Config, log *slog.Logger) error {
	st, err := store.New(ctx, cfg.Database.DSN)
	if err != nil {
		return fmt.Errorf("store: %w", err)
	}
	defer st.Close()

	if d := cfg.Storage.Retention.Std(); d > 0 {
		if err := st.ApplyGlobalRetention(ctx, d); err != nil {
			return fmt.Errorf("applying retention: %w", err)
		}
	}

	sealer, err := loadSealer(cfg)
	if err != nil {
		return err
	}

	metrics := observability.NewMetrics()

	fwd, err := newForwarder(cfg, log)
	if err != nil {
		return fmt.Errorf("triggers.forward: %w", err)
	}

	e, err := engine.New(st, nil, engine.Config{
		Forwarder:       fwd,
		Shards:          cfg.Engine.Shards,
		ShardQueue:      cfg.Engine.ShardQueue,
		BatchMaxRows:    cfg.Engine.BatchMaxRows,
		BatchMaxWait:    cfg.Engine.BatchMaxWait.Std(),
		MaxPayloadBytes: cfg.Engine.MaxPayloadBytes,
		Registerer:      metrics.Registerer(),
		Logger:          log,
	})
	if err != nil {
		return fmt.Errorf("engine: %w", err)
	}

	b, err := newBroker(ctx, cfg, st, e, log)
	if err != nil {
		return fmt.Errorf("broker: %w", err)
	}
	e.AttachBroker(engine.AdaptBroker(b))

	// Pairing service is built here (not in mountAPIs) so the flow runtime
	// can auto-register first-seen virtual devices through it (#84).
	advertised := cfg.MQTT.AdvertisedURL
	if advertised == "" {
		advertised = "mqtts://" + b.TLSAddr()
	}
	pairer := pairing.New(st, sealer, pairing.Config{
		BrokerURL:         advertised,
		CertTTL:           cfg.Pairing.CertTTL.Std(),
		EnforceLatestCert: cfg.Pairing.EnforceLatestCert,
		Version:           version,
		BcryptCost:        cfg.Pairing.BcryptCost,
	})
	pairer.OnRegistered = e.HandleDeviceRegistered

	// Flow runtime shares the engine's live bus so AstarteSource blocks see
	// the same device events as the stream socket (v2.0 process wiring).
	// Virtual-device pools (#84) write through the engine's device-owned
	// ingest path: storage rows without MQTT. With auto_register they also
	// register first-seen ids through the pairing door; an id already taken
	// maps to flow.ErrVirtualDeviceRegistered so the block drops and logs.
	flowMgr := flow.NewManager()
	flowSvc := flowapi.NewService(st, flowMgr, blocks.DefaultRegistry(), e.Bus(),
		func(ctx context.Context, realm, deviceID, ifaceName, path string, payload json.RawMessage, ts *time.Time) error {
			id, err := deviceid.Parse(deviceID)
			if err != nil {
				return fmt.Errorf("virtual device id %q: %w", deviceID, err)
			}
			return e.PublishDeviceValue(ctx, realm, id, ifaceName, path, payload, ts)
		},
		func(ctx context.Context, realm, deviceID string) error {
			_, err := pairer.Register(ctx, realm, deviceID, "")
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
	// right now, before the first tick — capping even no_ttl interfaces
	// (upstream clamps at write time; Astrate stores no per-row TTL).
	go func() {
		t := time.NewTicker(time.Hour)
		defer t.Stop()
		for {
			if err := st.EnforceRealmRetentionCeilings(ctx); err != nil {
				log.Warn("retention ceiling sweep failed", "error", err)
			}
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
	<-serveErr // ListenAndServe returns http.ErrServerClosed after Shutdown
	return nil
}

// newForwarder builds the Forwarder that receives trigger actions which are
// not HTTP webhooks (docs/DESIGN.md §1.1). The empty kind is the default and
// disables forwarding: those actions are logged and counted "skipped", which
// is the designed behaviour rather than a gap.
//
// The nil return is deliberately a nil *interface*, not a typed nil pointer —
// the executor tests its Forwarder against nil, and a typed nil would satisfy
// that test while panicking on the first custom action.
func newForwarder(cfg config.Config, log *slog.Logger) (triggers.Forwarder, error) {
	f := cfg.Triggers.Forward
	switch f.Kind {
	case "":
		return nil, nil
	case "http":
		h, err := forward.New(forward.Config{
			URL:           f.URL,
			Method:        f.Method,
			StaticHeaders: f.StaticHeaders,
		})
		if err != nil {
			return nil, err
		}
		log.Info("trigger forwarding enabled", "kind", f.Kind, "url", f.URL)
		return h, nil
	case "nats":
		h, err := newNATSForwarder(f)
		if err != nil {
			return nil, err
		}
		log.Info("trigger forwarding enabled", "kind", f.Kind, "url", f.URL, "subject", f.Subject)
		return h, nil
	default:
		// config.validate rejects every other kind, so this is unreachable
		// unless the two lists drift apart.
		return nil, fmt.Errorf("unsupported kind %q", f.Kind)
	}
}

// newBroker builds the broker, loading the server TLS identity from config
// unless dev mode runs without it.
func newBroker(ctx context.Context, cfg config.Config, st *store.Store, e *engine.Engine, log *slog.Logger) (*broker.Broker, error) {
	bcfg := broker.Config{
		TLSAddr:           cfg.MQTT.Addr,
		SessionStorePath:  cfg.MQTT.SessionStorePath,
		InsecureDevMode:   cfg.MQTT.InsecureDevMode,
		DevAddr:           cfg.MQTT.DevAddr,
		EnforceLatestCert: cfg.Pairing.EnforceLatestCert,
		MaxPacketBytes:    cfg.MQTT.MaxPacketBytes,
		Logger:            log,
	}
	switch {
	case cfg.MQTT.TLSCertFile != "":
		cert, err := tls.LoadX509KeyPair(cfg.MQTT.TLSCertFile, cfg.MQTT.TLSKeyFile)
		if err != nil {
			return nil, fmt.Errorf("loading broker TLS keypair: %w", err)
		}
		bcfg.ServerTLSCert = cert
	case cfg.MQTT.InsecureDevMode:
		// The broker always binds an mTLS listener and needs a server cert;
		// dev mode without one gets an ephemeral self-signed identity so the
		// binary runs zero-config (the plaintext dev listener is the one
		// devices actually use here).
		cert, err := selfSignedDevCert()
		if err != nil {
			return nil, fmt.Errorf("generating dev broker certificate: %w", err)
		}
		bcfg.ServerTLSCert = cert
		log.Warn("insecure_dev_mode: using an ephemeral self-signed broker certificate")
	}
	return broker.New(ctx, bcfg, st, e, e)
}

// selfSignedDevCert mints a throwaway ECDSA server certificate for the dev-mode
// broker TLS listener (never used in production: config requires real cert
// files outside insecure_dev_mode).
func selfSignedDevCert() (tls.Certificate, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, err
	}
	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: "astrate-dev"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}, nil
}

// mountAPIs builds the HTTP handler carrying every REST surface plus the
// observability endpoints (wrapped in CORS when configured), and returns the
// housekeeping service for auto-provisioning.
func mountAPIs(cfg config.Config, st *store.Store, e *engine.Engine, b *broker.Broker, sealer *store.KeySealer, metrics *observability.Metrics, flowSvc *flowapi.Service, pairer *pairing.Service, log *slog.Logger) (http.Handler, *housekeeping.Service, error) {
	mw := auth.NewMiddleware(st)
	mux := http.NewServeMux()

	pairing.NewAPI(pairer, mw, pairing.APIConfig{
		RegisterRate:     cfg.Pairing.RegisterRate,
		RegisterBurst:    cfg.Pairing.RegisterBurst,
		CredentialsRate:  cfg.Pairing.CredentialsRate,
		CredentialsBurst: cfg.Pairing.CredentialsBurst,
	}).Mount(mux)

	hkKeys, err := cfg.HousekeepingKeys()
	if err != nil {
		return nil, nil, err
	}
	hkSvc := housekeeping.NewService(st, sealer, b, log).
		WithDefaultDatastreamMaximumStorageRetention(cfg.Housekeeping.DefaultDatastreamMaximumStorageRetention).
		WithRealmDeletionDisabled(cfg.Housekeeping.RealmDeletionDisabled)
	housekeeping.NewAPI(hkSvc, mw, hkKeys).Mount(mux)
	realmSvc := realm.NewService(st, e, log).WithDisconnecter(b)
	realmSvc.OnDeletionStart = e.HandleDeviceDeletionStarted
	realmSvc.OnDeletionFinish = e.HandleDeviceDeletionFinished
	realm.NewAPI(realmSvc, mw).Mount(mux)
	appengine.NewAPI(appengine.NewService(st, e, log), mw).Mount(mux)
	apstream.NewAPI(e.Bus(), mw).Mount(mux)
	// Phoenix Channels socket (phoenix.js V2), alongside the Astrate-native
	// socket above: two protocols, one bus.
	channels.NewAPI(e.Bus(), st).Mount(mux)
	// Flow operator API: pipelines CRUD + start/stop/status (v2.0).
	flowapi.NewAPI(flowSvc, mw).Mount(mux)

	// Upstream-parity per-service health endpoints (the dashboard's API
	// status indicators poll them).
	for _, svc := range []string{"appengine", "realmmanagement", "pairing"} {
		observability.MountServiceCompat(mux, svc)
	}

	// Upstream-parity per-service version endpoints (issue #77): every
	// service answers an unauthenticated GET /{service}/version with
	// {"data": version}; AppEngine's realm-scoped variant requires auth,
	// Pairing's is public — both measured on upstream 1.2.0
	// (test/conformance/upstream/verify-versions.json). Realm Management's
	// realm-scoped route is served by its own API.
	for _, svc := range []string{"appengine", "realmmanagement", "pairing", "housekeeping"} {
		observability.MountVersionCompat(mux, svc, version)
	}
	mux.Handle("GET /appengine/v1/{realm}/version",
		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(version)))
	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(version))

	metrics.RegisterBrokerSessions(func() float64 { return float64(b.SessionCount()) })
	metrics.RegisterDBPool(func() observability.DBPoolStats {
		s := st.Stat()
		return observability.DBPoolStats{
			AcquiredConns: s.AcquiredConns(),
			IdleConns:     s.IdleConns(),
			TotalConns:    s.TotalConns(),
			MaxConns:      s.MaxConns(),
		}
	})
	health := observability.NewHealth(metrics.Handler())
	health.AddReadiness("database", st.Health)
	health.AddReadiness("broker", brokerReadiness(b))
	health.Mount(mux)

	// Swagger UI + OpenAPI YAML specs embedded in the binary.
	swagger.Mount(mux)

	// Unmatched routes under a service prefix answer that service's JSON error
	// envelope rather than Go's plain-text "404 page not found", which any
	// client parsing the envelope would choke on.
	handler := httpx.NotFound(mux)
	if len(cfg.HTTP.CORSAllowedOrigins) > 0 {
		handler = httpx.CORS(cfg.HTTP.CORSAllowedOrigins)(handler)
	}
	return handler, hkSvc, nil
}

// brokerReadiness reports the broker listener as ready when a TCP connection to
// it is accepted (docs/DESIGN.md §5.2 readiness broker check).
func brokerReadiness(b *broker.Broker) observability.Check {
	return func(ctx context.Context) error {
		addr := b.TLSAddr()
		if addr == "" {
			return errors.New("broker listener not bound")
		}
		var d net.Dialer
		conn, err := d.DialContext(ctx, "tcp", addr)
		if err != nil {
			return err
		}
		return conn.Close()
	}
}

// autoProvisionRealm creates the configured realm on first boot, a no-op when
// it already exists (docs/DESIGN.md §5.1).
func autoProvisionRealm(ctx context.Context, st *store.Store, hk *housekeeping.Service, cfg config.Config, log *slog.Logger) error {
	if _, err := st.GetRealmByName(ctx, cfg.Realm.Name); err == nil {
		log.Info("auto-provision realm already exists", "realm", cfg.Realm.Name)
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
func runHealthcheck() int {
	addr := os.Getenv("ASTRATE_HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	if strings.HasPrefix(addr, ":") {
		addr = "127.0.0.1" + addr
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// The target is the operator-configured local listener address (a container
	// self-probe), not attacker-controlled input — not an SSRF sink.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+addr+"/astrate/v1/readiness", nil) //nolint:gosec // G704: self-probe of the local readiness endpoint
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	resp, err := http.DefaultClient.Do(req) //nolint:gosec // G704: self-probe of the local readiness endpoint
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}

// newLogger builds the slog handler per config (docs/DESIGN.md §5.2).
func newLogger(c config.LogConfig) *slog.Logger {
	var level slog.Level
	switch c.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: level}
	var h slog.Handler
	if c.Format == "text" {
		h = slog.NewTextHandler(os.Stderr, opts)
	} else {
		h = slog.NewJSONHandler(os.Stderr, opts)
	}
	return slog.New(h)
}
