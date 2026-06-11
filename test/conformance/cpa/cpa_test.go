// Package cpa is conformance checkpoint CP-A (docs/ROADMAP.md §0.3, §5):
// the wire-frozen proof that Astrate's Pairing API serves *unmodified*
// official Astarte clients. It boots a real store + auth + pairing stack
// against a live TimescaleDB and drives it with:
//
//   - the official astarte-device-sdk-go pairing client (registration via
//     astarte-go, then the SDK device's own CSR/credentials handshake),
//     asserting the SDK obtains a certificate it considers valid (it
//     persists it and does not re-request on the next start);
//   - the pinned astartectl release binary
//     (`astartectl pairing agent register`).
//
// The database is the compose service (`make up`) or ASTRATE_TEST_DSN; the
// suite skips with instructions when neither is reachable.
package cpa

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	sdkdevice "github.com/astarte-platform/astarte-device-sdk-go/device"
	astarteclient "github.com/astarte-platform/astarte-go/client"
	"github.com/astarte-platform/astarte-go/interfaces"
	"github.com/astarte-platform/astarte-go/misc"
	"github.com/jackc/pgx/v5"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/pairing"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// composeDSN is the docker-compose database (`make up` at the repo root).
const composeDSN = "postgres://astrate:astrate@127.0.0.1:5432/astrate?sslmode=disable"

// testInterface is a minimal device-owned datastream the SDK requires
// before connecting (Connect refuses an interface-less device).
const testInterface = `{
	"interface_name": "org.astrate.conformance.Test",
	"version_major": 0,
	"version_minor": 1,
	"type": "datastream",
	"ownership": "device",
	"mappings": [{"endpoint": "/value", "type": "double"}]
}`

// recordedCall is one HTTP exchange seen by the harness server.
type recordedCall struct {
	Method string
	Path   string
	Status int
}

// recorder captures every request/response pair crossing the harness server.
type recorder struct {
	mu    sync.Mutex
	calls []recordedCall
}

func (r *recorder) record(c recordedCall) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, c)
}

// count returns how many recorded calls match method, path suffix and status.
func (r *recorder) count(method, pathSuffix string, status int) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for _, c := range r.calls {
		if c.Method == method && strings.HasSuffix(c.Path, pathSuffix) && c.Status == status {
			n++
		}
	}
	return n
}

// statusWriter remembers the response status code.
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (s *statusWriter) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// fixture is the composed Astrate pairing instance under test.
type fixture struct {
	st         *store.Store
	realm      *store.Realm
	realmName  string
	server     *httptest.Server
	pairingURL string
	jwtKeyPEM  []byte
	rec        *recorder
}

// dialDSN finds a reachable database: ASTRATE_TEST_DSN first, then the
// compose default.
func dialDSN(t *testing.T) string {
	t.Helper()
	candidates := []string{os.Getenv("ASTRATE_TEST_DSN"), composeDSN}
	for _, dsn := range candidates {
		if dsn == "" {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		conn, err := pgx.Connect(ctx, dsn)
		if err == nil {
			_ = conn.Close(ctx)
			cancel()
			return dsn
		}
		cancel()
	}
	t.Skip("CP-A needs a database: run `make up` at the repo root or set ASTRATE_TEST_DSN")
	return ""
}

// newFixture provisions a fresh realm and serves the pairing API on a real
// TCP listener.
func newFixture(t *testing.T) *fixture {
	t.Helper()
	ctx := context.Background()

	st, err := store.New(ctx, dialDSN(t))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(st.Close)

	masterKey := make([]byte, store.MasterKeySize)
	if _, err := rand.Read(masterKey); err != nil {
		t.Fatal(err)
	}
	sealer, err := store.NewKeySealer(masterKey)
	if err != nil {
		t.Fatal(err)
	}

	jwtKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	jwtKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(jwtKey),
	})
	pubDER, err := x509.MarshalPKIXPublicKey(&jwtKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pubPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))

	realmName := "cpa" + strconv.FormatInt(time.Now().UnixNano(), 36)
	caPEM, sealedKey, err := pairing.ProvisionCA(realmName, 0, sealer)
	if err != nil {
		t.Fatalf("ProvisionCA: %v", err)
	}
	realm, err := st.CreateRealm(ctx, store.NewRealm{
		Name:               realmName,
		JWTPublicKeysPEM:   []string{pubPEM},
		CACertificatePEM:   caPEM,
		CAPrivateKeySealed: sealedKey,
	})
	if err != nil {
		t.Fatalf("CreateRealm: %v", err)
	}

	svc := pairing.New(st, sealer, pairing.Config{
		// A closed local port: pairing must fully succeed, the MQTT
		// connection (out of CP-A scope, broker lands in M5) must fail.
		BrokerURL: "mqtts://127.0.0.1:1",
		Version:   "1.2.0",
	})
	api := pairing.NewAPI(svc, auth.NewMiddleware(st), pairing.APIConfig{})
	mux := http.NewServeMux()
	api.Mount(mux)

	rec := &recorder{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		mux.ServeHTTP(sw, r)
		rec.record(recordedCall{Method: r.Method, Path: r.URL.Path, Status: sw.status})
	}))
	t.Cleanup(server.Close)

	return &fixture{
		st: st, realm: realm, realmName: realmName,
		server: server, pairingURL: server.URL + "/pairing",
		jwtKeyPEM: jwtKeyPEM, rec: rec,
	}
}

// TestCPA runs the two CP-A clients against one composed instance.
func TestCPA(t *testing.T) {
	f := newFixture(t)
	t.Run("GoSDKPairing", f.testGoSDKPairing)
	t.Run("AstartectlRegister", f.testAstartectlRegister)
}

// testGoSDKPairing drives the official Go SDK end to end through pairing:
// register (astarte-go agent client) → SDK keygen + CSR → credentials →
// info → MQTT attempt (expected to fail: no broker yet). The certificate
// the SDK stored must parse, carry the Astarte CN, chain to the realm CA,
// and satisfy the SDK's own validity check — proven by a second SDK start
// that reuses it without requesting new credentials.
func (f *fixture) testGoSDKPairing(t *testing.T) {
	id, err := deviceid.Random()
	if err != nil {
		t.Fatal(err)
	}
	deviceID := id.String()

	// Flow A through the official agent client.
	agent, err := astarteclient.NewClientWithIndividualURLs(
		map[misc.AstarteService]string{misc.Pairing: f.pairingURL}, nil)
	if err != nil {
		t.Fatalf("astarte-go client: %v", err)
	}
	if err := agent.SetTokenFromPrivateKey(f.jwtKeyPEM); err != nil {
		t.Fatalf("SetTokenFromPrivateKey: %v", err)
	}
	secret, err := agent.Pairing.RegisterDevice(f.realmName, deviceID)
	if err != nil {
		t.Fatalf("official RegisterDevice: %v", err)
	}
	if len(secret) != 44 {
		t.Fatalf("credentials secret length: got %d, want 44 (%q)", len(secret), secret)
	}
	if _, err := base64.StdEncoding.DecodeString(secret); err != nil {
		t.Fatalf("credentials secret is not standard base64: %v", err)
	}

	// Flows B+C through the SDK device.
	cryptoDir := filepath.Join(t.TempDir(), "crypto")
	if err := os.MkdirAll(cryptoDir, 0o700); err != nil {
		t.Fatal(err)
	}
	connectOnce := func() error {
		opts := sdkdevice.NewDeviceOptions()
		opts.UseMqttStore = false
		opts.UseDatabase = false
		opts.ConnectRetry = false
		opts.AutoReconnect = false
		opts.IgnoreSSLErrors = true
		opts.CryptoDir = cryptoDir

		d, err := sdkdevice.NewDeviceWithOptions(deviceID, f.realmName, secret, f.pairingURL, opts)
		if err != nil {
			t.Fatalf("NewDeviceWithOptions: %v", err)
		}
		iface, err := interfaces.ParseInterface([]byte(testInterface))
		if err != nil {
			t.Fatalf("ParseInterface: %v", err)
		}
		if err := d.AddInterface(iface); err != nil {
			t.Fatalf("AddInterface: %v", err)
		}

		ch := make(chan error, 1)
		d.Connect(ch)
		select {
		case err := <-ch:
			return err
		case <-time.After(90 * time.Second):
			t.Fatal("SDK Connect did not settle within 90s")
			return nil
		}
	}

	if err := connectOnce(); err == nil {
		t.Fatal("Connect succeeded with no broker running — expected an MQTT-stage failure")
	}

	credentialsPath := "/devices/" + deviceID + "/protocols/astarte_mqtt_v1/credentials"
	if n := f.rec.count("POST", credentialsPath, http.StatusCreated); n != 1 {
		t.Fatalf("credentials issuances after first connect: got %d, want 1 (calls: %+v)", n, f.rec.calls)
	}
	if n := f.rec.count("GET", "/devices/"+deviceID, http.StatusOK); n < 1 {
		t.Fatalf("info endpoint never hit (calls: %+v)", f.rec.calls)
	}

	// The certificate the SDK stored.
	crtPEM, err := os.ReadFile(filepath.Join(cryptoDir, "device.crt"))
	if err != nil {
		t.Fatalf("SDK did not store a certificate: %v", err)
	}
	block, _ := pem.Decode(crtPEM)
	if block == nil {
		t.Fatal("stored certificate is not PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("stored certificate does not parse: %v", err)
	}
	if want := f.realmName + "/" + deviceID; cert.Subject.CommonName != want {
		t.Errorf("CN: got %q, want %q", cert.Subject.CommonName, want)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM([]byte(f.realm.CACertificatePEM)) {
		t.Fatal("adding realm CA to pool")
	}
	if _, err := cert.Verify(x509.VerifyOptions{
		Roots:     roots,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}); err != nil {
		t.Errorf("stored certificate does not chain to the realm CA: %v", err)
	}
	if !time.Now().Before(cert.NotAfter) {
		t.Error("stored certificate already expired")
	}

	// "A certificate it considers valid": a fresh SDK instance over the
	// same crypto dir must reuse it — no second credentials request.
	if err := connectOnce(); err == nil {
		t.Fatal("second Connect succeeded with no broker running")
	}
	if n := f.rec.count("POST", credentialsPath, http.StatusCreated); n != 1 {
		t.Fatalf("credentials issuances after second connect: got %d, want still 1 — the SDK re-requested a certificate it should consider valid", n)
	}

	// The device row reflects the handshake.
	dev, err := f.st.GetDevice(context.Background(), f.realm.ID, id)
	if err != nil {
		t.Fatalf("GetDevice: %v", err)
	}
	if dev.Status != store.DeviceStatusConfirmed {
		t.Errorf("device status: got %q, want confirmed", dev.Status)
	}
	if dev.CertSerial == nil || *dev.CertSerial != cert.SerialNumber.String() {
		t.Errorf("recorded serial: got %v, want %s", dev.CertSerial, cert.SerialNumber.String())
	}
}

// testAstartectlRegister registers a device with the pinned astartectl
// release binary, which also exercises its self-generated realm JWT against
// Astrate's auth layer.
func (f *fixture) testAstartectlRegister(t *testing.T) {
	bin := ensureAstartectl(t)

	keyFile := filepath.Join(t.TempDir(), "realm_private.pem")
	if err := os.WriteFile(keyFile, f.jwtKeyPEM, 0o600); err != nil {
		t.Fatal(err)
	}

	id, err := deviceid.Random()
	if err != nil {
		t.Fatal(err)
	}

	secret := runAstartectl(t, bin,
		"pairing", "agent", "register", id.String(),
		"--pairing-url", f.pairingURL,
		"--realm-name", f.realmName,
		"--realm-key", keyFile,
		"--compact-output",
	)
	if len(secret) != 44 {
		t.Fatalf("astartectl secret length: got %d, want 44 (%q)", len(secret), secret)
	}

	dev, err := f.st.GetDevice(context.Background(), f.realm.ID, id)
	if err != nil {
		t.Fatalf("device not registered in the store: %v", err)
	}
	if dev.Status != store.DeviceStatusRegistered {
		t.Errorf("device status: got %q, want registered", dev.Status)
	}
}
