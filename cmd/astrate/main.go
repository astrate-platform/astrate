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
	apstream "github.com/astrate-platform/astrate/internal/appengine/stream"
	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/config"
	"github.com/astrate-platform/astrate/internal/engine"
	"github.com/astrate-platform/astrate/internal/housekeeping"
	"github.com/astrate-platform/astrate/internal/observability"
	"github.com/astrate-platform/astrate/internal/pairing"
	"github.com/astrate-platform/astrate/internal/realm"
	"github.com/astrate-platform/astrate/internal/store"
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
// it in the §5.3 order: HTTP → broker → engine → store. It is the in-process
// entry point the boot tests drive.
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

	e, err := engine.New(st, nil, engine.Config{
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

	if err := e.Start(ctx); err != nil {
		return fmt.Errorf("starting engine: %w", err)
	}
	if err := b.Start(); err != nil {
		drainEngine(e, log)
		return fmt.Errorf("starting broker: %w", err)
	}

	mux, hkSvc, err := mountAPIs(cfg, st, e, b, sealer, metrics, log)
	if err != nil {
		shutdown(nil, b, e, log)
		return err
	}

	if cfg.Realm.Name != "" {
		if err := autoProvisionRealm(ctx, st, hkSvc, cfg, log); err != nil {
			shutdown(nil, b, e, log)
			return err
		}
	}

	srv := &http.Server{Addr: cfg.HTTP.Addr, Handler: mux, ReadHeaderTimeout: 10 * time.Second}
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
		shutdown(nil, b, e, log)
		return fmt.Errorf("http server: %w", err)
	}

	shutdown(srv, b, e, log)
	<-serveErr // ListenAndServe returns http.ErrServerClosed after Shutdown
	return nil
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

// mountAPIs builds the HTTP mux carrying every REST surface plus the
// observability endpoints, and returns the housekeeping service for
// auto-provisioning.
func mountAPIs(cfg config.Config, st *store.Store, e *engine.Engine, b *broker.Broker, sealer *store.KeySealer, metrics *observability.Metrics, log *slog.Logger) (*http.ServeMux, *housekeeping.Service, error) {
	mw := auth.NewMiddleware(st)
	mux := http.NewServeMux()

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
	hkSvc := housekeeping.NewService(st, sealer, b, log)
	housekeeping.NewAPI(hkSvc, mw, hkKeys).Mount(mux)
	realm.NewAPI(realm.NewService(st, e, log), mw).Mount(mux)
	appengine.NewAPI(appengine.NewService(st, e, log), mw).Mount(mux)
	apstream.NewAPI(e.Bus(), mw).Mount(mux)

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

	return mux, hkSvc, nil
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
	if _, err := hk.CreateRealm(ctx, cfg.Realm.Name, key, cfg.Realm.DeviceRegistrationLimit); err != nil {
		return fmt.Errorf("auto-provisioning realm %q: %w", cfg.Realm.Name, err)
	}
	log.Info("auto-provisioned realm", "realm", cfg.Realm.Name)
	return nil
}

// shutdown drains the stack in the §5.3 order. srv may be nil on a startup
// error before the HTTP server began serving.
func shutdown(srv *http.Server, b *broker.Broker, e *engine.Engine, log *slog.Logger) {
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
