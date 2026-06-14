//go:build integration && e2e

// Package security is the T5 hardening suite (docs/ROADMAP.md §10 file 9.8,
// DESIGN §4.5): it composes the real stack — store + engine + broker + pairing
// + the M7 REST surfaces — and probes the platform's security invariants:
// every guarded route refuses anonymous callers, pairing rate-limits brute
// force, the broker enforces mTLS and a TLS floor, oversize bodies are
// rejected, and a producer/properties zip-bomb cannot be processed.
package security

import (
	"bytes"
	"compress/zlib"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/astrate-platform/astrate/internal/appengine"
	apstream "github.com/astrate-platform/astrate/internal/appengine/stream"
	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/engine"
	"github.com/astrate-platform/astrate/internal/housekeeping"
	"github.com/astrate-platform/astrate/internal/pairing"
	"github.com/astrate-platform/astrate/internal/realm"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/payload"
)

const confDef = `{"interface_name":"org.astrate.sec.Conf","version_major":1,"version_minor":0,"type":"properties","ownership":"device","mappings":[{"endpoint":"/%{k}","type":"string","allow_unset":true}]}`

type secEnv struct {
	st        *store.Store
	broker    *broker.Broker
	pairer    *pairing.Service
	server    *httptest.Server
	realm     *store.Realm
	realmName string
	roots     *x509.CertPool
	sslURL    string
	realmKey  *rsa.PrivateKey
	confIface *store.StoredInterface
}

func discard() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

func newSecEnv(t *testing.T) *secEnv {
	t.Helper()
	ctx := context.Background()

	pool := testutil.StartTimescale(t)
	st, err := store.New(ctx, pool.Config().ConnString())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(st.Close)

	master := make([]byte, store.MasterKeySize)
	_, _ = rand.Read(master)
	sealer, err := store.NewKeySealer(master)
	if err != nil {
		t.Fatal(err)
	}

	realmKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	pubDER, _ := x509.MarshalPKIXPublicKey(&realmKey.PublicKey)
	realmPub := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
	instKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	instDER, _ := x509.MarshalPKIXPublicKey(&instKey.PublicKey)
	instPub := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: instDER}))

	var suffix [4]byte
	_, _ = rand.Read(suffix[:])
	realmName := "sec" + hex.EncodeToString(suffix[:])
	caPEM, sealedKey, err := pairing.ProvisionCA(realmName, 0, sealer)
	if err != nil {
		t.Fatalf("ProvisionCA: %v", err)
	}
	r, err := st.CreateRealm(ctx, store.NewRealm{
		Name: realmName, JWTPublicKeysPEM: []string{realmPub},
		CACertificatePEM: caPEM, CAPrivateKeySealed: sealedKey,
	})
	if err != nil {
		t.Fatalf("CreateRealm: %v", err)
	}
	conf, err := st.InstallInterface(ctx, r.ID, []byte(confDef))
	if err != nil {
		t.Fatalf("InstallInterface: %v", err)
	}

	e, err := engine.New(st, nil, engine.Config{Shards: 2, BatchMaxRows: 4, BatchMaxWait: 20 * time.Millisecond, Logger: discard()})
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

	pairer := pairing.New(st, sealer, pairing.Config{BrokerURL: "mqtts://" + b.TLSAddr(), Version: "1.2.0"})
	mw := auth.NewMiddleware(st)
	mux := http.NewServeMux()
	// A deliberately tiny register burst exercises the rate limiter.
	pairing.NewAPI(pairer, mw, pairing.APIConfig{RegisterRate: 0.01, RegisterBurst: 2}).Mount(mux)
	housekeeping.NewAPI(housekeeping.NewService(st, sealer, b, discard()), mw, []string{instPub}).Mount(mux)
	realm.NewAPI(realm.NewService(st, e, discard()), mw).Mount(mux)
	appengine.NewAPI(appengine.NewService(st, e, discard()), mw).Mount(mux)
	apstream.NewAPI(e.Bus(), mw).Mount(mux)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return &secEnv{
		st: st, broker: b, pairer: pairer, server: server, realm: r, realmName: realmName,
		roots: roots, sslURL: "ssl://" + b.TLSAddr(), realmKey: realmKey, confIface: conf,
	}
}

func (env *secEnv) token(t *testing.T, claim string) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		claim: []string{".*::.*"}, "exp": time.Now().Add(time.Hour).Unix(),
	})
	s, err := tok.SignedString(env.realmKey)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestSecurity(t *testing.T) {
	env := newSecEnv(t)
	t.Run("AuthzSweep", env.testAuthzSweep)
	t.Run("PairingRateLimit", env.testPairingRateLimit)
	t.Run("TLSHardening", env.testTLSHardening)
	t.Run("OversizeBody", env.testOversizeBody)
	t.Run("ProducerPropertiesZipBomb", env.testZipBomb)
}

// testAuthzSweep asserts every JWT-guarded route refuses an anonymous caller
// (401) and a caller bearing the wrong claim (403). A route table makes a new,
// unguarded route a visible failure.
func (env *secEnv) testAuthzSweep(t *testing.T) {
	id, _ := deviceid.Random()
	r := env.realmName
	dev := id.String()
	routes := []struct{ method, path string }{
		{"POST", "/pairing/v1/" + r + "/agent/devices"},
		{"DELETE", "/pairing/v1/" + r + "/agent/devices/" + dev},
		{"GET", "/realmmanagement/v1/" + r + "/interfaces"},
		{"POST", "/realmmanagement/v1/" + r + "/interfaces"},
		{"GET", "/realmmanagement/v1/" + r + "/interfaces/org.astrate.sec.Conf/1"},
		{"DELETE", "/realmmanagement/v1/" + r + "/interfaces/org.astrate.sec.Conf/1"},
		{"GET", "/realmmanagement/v1/" + r + "/triggers"},
		{"POST", "/realmmanagement/v1/" + r + "/triggers"},
		{"GET", "/realmmanagement/v1/" + r + "/config/auth"},
		{"PUT", "/realmmanagement/v1/" + r + "/config/auth"},
		{"GET", "/housekeeping/v1/realms"},
		{"POST", "/housekeeping/v1/realms"},
		{"GET", "/housekeeping/v1/realms/" + r},
		{"DELETE", "/housekeeping/v1/realms/" + r},
		{"GET", "/appengine/v1/" + r + "/devices"},
		{"GET", "/appengine/v1/" + r + "/devices/" + dev},
		{"PATCH", "/appengine/v1/" + r + "/devices/" + dev},
		{"GET", "/appengine/v1/" + r + "/devices/" + dev + "/interfaces/org.astrate.sec.Conf"},
		{"GET", "/appengine/v1/" + r + "/groups"},
		{"POST", "/appengine/v1/" + r + "/groups"},
		{"GET", "/astrate/v1/" + r + "/socket"},
	}
	for _, rt := range routes {
		code := env.do(t, rt.method, rt.path, "")
		if code != http.StatusUnauthorized {
			t.Errorf("%s %s without a token: got %d, want 401", rt.method, rt.path, code)
		}
	}
	// A valid token with the wrong claim is forbidden, not unauthorized.
	wrong := env.token(t, string(auth.ClaimHousekeeping))
	if code := env.do(t, "GET", "/appengine/v1/"+r+"/devices", wrong); code != http.StatusForbidden {
		t.Errorf("appengine with a_ha token: got %d, want 403", code)
	}
}

// testPairingRateLimit hammers registration past the (tiny) burst and expects a
// 429 from the per-IP token bucket.
func (env *secEnv) testPairingRateLimit(t *testing.T) {
	token := env.token(t, string(auth.ClaimPairing))
	got429 := false
	for range 6 {
		id, _ := deviceid.Random()
		body := fmt.Sprintf(`{"data":{"hw_id":%q}}`, id.String())
		code := env.doBody(t, "POST", "/pairing/v1/"+env.realmName+"/agent/devices", token, body)
		if code == http.StatusTooManyRequests {
			got429 = true
			break
		}
	}
	if !got429 {
		t.Error("registration never rate-limited past the burst")
	}
}

// testTLSHardening proves the broker requires a client certificate (no
// client-auth bypass) and refuses pre-1.2 TLS.
func (env *secEnv) testTLSHardening(t *testing.T) {
	addr := strings.TrimPrefix(env.sslURL, "ssl://")

	// No client certificate, pinned to TLS 1.2 so RequireAndVerifyClientCert
	// surfaces during the handshake (in TLS 1.3 the client finishes before the
	// server validates the client cert, so Dial would return before the alert).
	c1, err := tls.Dial("tcp", addr, &tls.Config{ //nolint:gosec // probing the server's client-auth requirement
		InsecureSkipVerify: true, MinVersion: tls.VersionTLS12, MaxVersion: tls.VersionTLS12,
	})
	if err == nil {
		_ = c1.Close()
		t.Error("broker completed a TLS 1.2 handshake with no client certificate")
	}

	// TLS 1.1 ceiling: below the broker's 1.2 floor, the handshake must fail.
	c2, err := tls.Dial("tcp", addr, &tls.Config{InsecureSkipVerify: true, MaxVersion: tls.VersionTLS11}) //nolint:gosec // probing the TLS floor
	if err == nil {
		_ = c2.Close()
		t.Error("broker accepted a TLS 1.1 handshake")
	}
}

// testOversizeBody asserts the request-body cap rejects a >1 MiB interface
// install with 400 rather than buffering it.
func (env *secEnv) testOversizeBody(t *testing.T) {
	token := env.token(t, string(auth.ClaimRealmManagement))
	huge := `{"data":"` + strings.Repeat("A", 2<<20) + `"}`
	code := env.doBody(t, "POST", "/realmmanagement/v1/"+env.realmName+"/interfaces", token, huge)
	if code != http.StatusBadRequest {
		t.Errorf("oversize interface install: got %d, want 400", code)
	}
}

// testZipBomb connects a device, sets a property, then sends a
// producer/properties control frame declaring a gigabyte of inflated data. The
// engine must reject it (cap, not OOM): the legitimately-set property survives
// because the malicious purge is never processed, and the broker stays up.
func (env *secEnv) testZipBomb(t *testing.T) {
	ctx := context.Background()
	id, _ := deviceid.Random()
	secret, err := env.pairer.Register(ctx, env.realmName, id.String(), "")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	key, csr := testutil.DeviceCSR(t)
	crt, err := env.pairer.Credentials(ctx, env.realmName, id.String(), secret, csr, netip.MustParseAddr("127.0.0.1"))
	if err != nil {
		t.Fatalf("Credentials: %v", err)
	}
	tlsCfg := testutil.DeviceTLSConfig(t, crt, key, env.roots)
	dev := testutil.ConnectAstarteDevice(t, env.sslURL, env.realmName, id, tlsCfg, true)
	t.Cleanup(dev.Disconnect)
	dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{"org.astrate.sec.Conf": {1, 0}}))

	dev.PublishValue(t, "org.astrate.sec.Conf", "/keep", "v", nil, payload.FormatBSON, 1)
	waitFor(t, 5*time.Second, func() bool {
		p, err := env.st.GetProperty(ctx, env.realm.ID, id, env.confIface.ID, "/keep")
		return err == nil && string(p.Value) == `"v"`
	})

	// 4-byte BE declared size of 1 GiB, then a tiny zlib stream of a short list:
	// the declared size is far over the absolute ceiling, so the engine rejects
	// the frame before inflating it.
	var zbuf bytes.Buffer
	zw := zlib.NewWriter(&zbuf)
	_, _ = zw.Write([]byte("/x"))
	_ = zw.Close()
	frame := make([]byte, 4)
	binary.BigEndian.PutUint32(frame, 1<<30)
	frame = append(frame, zbuf.Bytes()...)
	testutil.WaitToken(t, dev.Client.Publish(dev.Base()+"/control/producer/properties", 2, false, frame), 5*time.Second)

	// The malicious frame must not have purged the legitimately-set property,
	// and the stack must still be serving.
	time.Sleep(500 * time.Millisecond)
	if p, err := env.st.GetProperty(ctx, env.realm.ID, id, env.confIface.ID, "/keep"); err != nil || string(p.Value) != `"v"` {
		t.Errorf("zip-bomb frame purged the property (err=%v) — cap not enforced", err)
	}
	if got := env.do(t, "GET", "/astrate/v1/"+env.realmName+"/socket", ""); got != http.StatusUnauthorized {
		t.Errorf("server unresponsive after zip-bomb (got %d)", got)
	}
}

// --- HTTP helpers -----------------------------------------------------------

func (env *secEnv) do(t *testing.T, method, path, token string) int {
	return env.doBody(t, method, path, token, "")
}

func (env *secEnv) doBody(t *testing.T, method, path, token, body string) int {
	t.Helper()
	var rdr io.Reader
	if body != "" {
		rdr = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, env.server.URL+path, rdr)
	if err != nil {
		t.Fatal(err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	_ = resp.Body.Close()
	return resp.StatusCode
}

func waitFor(t *testing.T, timeout time.Duration, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatal("condition not met within timeout")
}
