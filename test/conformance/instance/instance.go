// Package instance composes a real in-process Astrate stack — store + engine +
// broker + pairing on a live TimescaleDB — for the M9 conformance runners
// (docs/ROADMAP.md §10). It is the shared harness behind the gosdk, atomvm,
// and pysdk packages: each provisions a realm, installs interfaces, and drives
// the running instance with a different client. The database is the compose
// service (`make up`) or ASTRATE_TEST_DSN; New skips the test when neither is
// reachable.
package instance

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/engine"
	"github.com/astrate-platform/astrate/internal/pairing"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

//nolint:gosec // G101: local-development compose DSN, not a real credential.
const composeDSN = "postgres://astrate:astrate@127.0.0.1:5432/astrate?sslmode=disable"

// Instance is a running Astrate stack under test.
type Instance struct {
	Store      *store.Store
	Engine     *engine.Engine
	Broker     *broker.Broker
	Pairer     *pairing.Service
	Realm      *store.Realm
	RealmName  string
	Roots      *x509.CertPool
	SSLURL     string // ssl://host:port for the MQTT broker
	PairingURL string // http://host:port/pairing
	BaseURL    string // http://host:port (mounts /pairing)
	JWTKeyPEM  []byte // realm JWT private key (PEM) for agent tokens
	Interfaces map[string]*store.StoredInterface
}

// Config tunes the composed instance.
type Config struct {
	// CertTTL overrides the issued client-certificate lifetime (e.g. a short
	// TTL for cert-renewal scenarios). Zero selects the pairing default.
	CertTTL time.Duration
	// Interfaces are installed into the realm and keyed by name in Interfaces.
	Interfaces map[string]string
	// Registerer receives the engine's Prometheus collectors (reject counters,
	// shard depth) so the load runner can assert on them; nil leaves them off.
	Registerer prometheus.Registerer
}

// dialDSN finds a reachable database: ASTRATE_TEST_DSN first, then compose.
func dialDSN(t testing.TB) string {
	t.Helper()
	for _, dsn := range []string{os.Getenv("ASTRATE_TEST_DSN"), composeDSN} {
		if dsn == "" {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		conn, err := pgx.Connect(ctx, dsn)
		cancel()
		if err == nil {
			_ = conn.Close(context.Background())
			return dsn
		}
	}
	t.Skip("conformance needs a database: run `make up` at the repo root or set ASTRATE_TEST_DSN")
	return ""
}

func discard() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

// New composes and starts a fresh instance with the given config.
func New(t testing.TB, cfg Config) *Instance {
	t.Helper()
	ctx := context.Background()

	st, err := store.New(ctx, dialDSN(t))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(st.Close)

	master := make([]byte, store.MasterKeySize)
	if _, err := rand.Read(master); err != nil {
		t.Fatal(err)
	}
	sealer, err := store.NewKeySealer(master)
	if err != nil {
		t.Fatal(err)
	}

	jwtKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	jwtKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(jwtKey)})
	pubDER, err := x509.MarshalPKIXPublicKey(&jwtKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pubPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))

	realmName := "cd" + strconv.FormatInt(time.Now().UnixNano(), 36)
	caPEM, sealedKey, err := pairing.ProvisionCA(realmName, 0, sealer)
	if err != nil {
		t.Fatalf("ProvisionCA: %v", err)
	}
	realm, err := st.CreateRealm(ctx, store.NewRealm{
		Name: realmName, JWTPublicKeysPEM: []string{pubPEM},
		CACertificatePEM: caPEM, CAPrivateKeySealed: sealedKey,
	})
	if err != nil {
		t.Fatalf("CreateRealm: %v", err)
	}

	ifaces := make(map[string]*store.StoredInterface, len(cfg.Interfaces))
	for name, def := range cfg.Interfaces {
		si, err := st.InstallInterface(ctx, realm.ID, []byte(def))
		if err != nil {
			t.Fatalf("InstallInterface(%s): %v", name, err)
		}
		ifaces[name] = si
	}

	e, err := engine.New(st, nil, engine.Config{
		Shards: 4, BatchMaxRows: 8, BatchMaxWait: 20 * time.Millisecond,
		Registerer: cfg.Registerer, Logger: discard(),
	})
	if err != nil {
		t.Fatalf("engine.New: %v", err)
	}
	serverCert, roots := testutil.ServerTLSCert(t)
	b, err := broker.New(ctx, broker.Config{
		TLSAddr: "127.0.0.1:0", ServerTLSCert: serverCert,
		SessionStorePath: filepath.Join(t.TempDir(), "sessions.db"), Logger: discard(),
	}, st, e, e)
	if err != nil {
		t.Fatalf("broker.New: %v", err)
	}
	e.AttachBroker(engine.AdaptBroker(b))

	runCtx, cancel := context.WithCancel(ctx)
	if err := e.Start(runCtx); err != nil {
		cancel()
		t.Fatalf("engine.Start: %v", err)
	}
	if err := b.Start(); err != nil {
		cancel()
		t.Fatalf("broker.Start: %v", err)
	}
	t.Cleanup(func() {
		_ = b.Close()
		dctx, dcancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer dcancel()
		_ = e.Drain(dctx)
		cancel()
	})

	pairer := pairing.New(st, sealer, pairing.Config{
		BrokerURL: "mqtts://" + b.TLSAddr(), Version: "1.2.0", CertTTL: cfg.CertTTL,
	})
	mux := http.NewServeMux()
	pairing.NewAPI(pairer, auth.NewMiddleware(st), pairing.APIConfig{}).Mount(mux)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return &Instance{
		Store: st, Engine: e, Broker: b, Pairer: pairer, Realm: realm, RealmName: realmName,
		Roots: roots, SSLURL: "ssl://" + b.TLSAddr(),
		PairingURL: server.URL + "/pairing", BaseURL: server.URL,
		JWTKeyPEM: jwtKeyPEM, Interfaces: ifaces,
	}
}

// Register obtains a credentials secret for a device through the pairing
// service (initialFormat may be "" or "json").
func (in *Instance) Register(t testing.TB, deviceID, initialFormat string) string {
	t.Helper()
	secret, err := in.Pairer.Register(context.Background(), in.RealmName, deviceID, initialFormat)
	if err != nil {
		t.Fatalf("Register(%s): %v", deviceID, err)
	}
	return secret
}

// IssueTLS runs the CSR/credentials handshake and returns a client TLS config a
// testutil.AstarteDevice can connect with.
func (in *Instance) IssueTLS(t testing.TB, deviceID, secret string) *tls.Config {
	t.Helper()
	key, csr := testutil.DeviceCSR(t)
	crt, err := in.Pairer.Credentials(context.Background(), in.RealmName, deviceID, secret, csr,
		netip.MustParseAddr("127.0.0.1"))
	if err != nil {
		t.Fatalf("Credentials(%s): %v", deviceID, err)
	}
	return testutil.DeviceTLSConfig(t, crt, key, in.Roots)
}

// NewDevice is the common path: register, issue a cert, and connect a
// testutil device with a clean session.
func (in *Instance) NewDevice(t testing.TB, initialFormat string) (*testutil.AstarteDevice, deviceid.ID) {
	t.Helper()
	id, err := deviceid.Random()
	if err != nil {
		t.Fatal(err)
	}
	secret := in.Register(t, id.String(), initialFormat)
	tlsCfg := in.IssueTLS(t, id.String(), secret)
	dev := testutil.ConnectAstarteDevice(t, in.SSLURL, in.RealmName, id, tlsCfg, true)
	t.Cleanup(dev.Disconnect)
	return dev, id
}
