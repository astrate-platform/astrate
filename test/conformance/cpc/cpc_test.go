// Package cpc is conformance checkpoint CP-C (docs/ROADMAP.md §0.3, §8.3): the
// wire-frozen proof that Astrate's operator-facing M7 surfaces — Housekeeping,
// Realm Management, and AppEngine — serve the *unmodified* pinned astartectl
// release binary. It composes the real stack (store + engine + broker +
// pairing + the three M7 HTTP APIs + the live-stream socket) in-process
// against a live TimescaleDB and drives it end to end with astartectl:
//
//	housekeeping realms create  → broker realm reload →
//	realm-management interfaces install (genericsensors.Values) →
//	pairing agent register → a real device publishes a sample →
//	appengine devices data-snapshot / get-samples return the value →
//	realm-management triggers install/list/show/delete round-trip.
//
// The device handshake (CSR → credentials) runs in-process — that path is
// already locked by CP-A/CP-B; here the device exists only to put data behind
// the AppEngine queries. The database is the compose service (`make up`) or
// ASTRATE_TEST_DSN; the suite skips with instructions when neither is reachable
// or when astartectl is unavailable.
package cpc

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

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

// composeDSN is the docker-compose database (`make up` at the repo root).
const composeDSN = "postgres://astrate:astrate@127.0.0.1:5432/astrate?sslmode=disable"

// valuesName is the standard interface CP-C installs and queries: an
// individual, parametric, device-owned datastream (one double per sensor).
const valuesName = "org.astarte-platform.genericsensors.Values"

const valuesInterface = `{
	"interface_name": "org.astarte-platform.genericsensors.Values",
	"version_major": 1,
	"version_minor": 0,
	"type": "datastream",
	"ownership": "device",
	"description": "Generic sensors sampled data.",
	"mappings": [
		{
			"endpoint": "/%{sensor_id}/value",
			"type": "double",
			"explicit_timestamp": true
		}
	]
}`

// triggerJSON is a data trigger over the installed interface — installed,
// listed, shown, and deleted to exercise the trigger CRUD wire shapes.
const triggerJSON = `{
	"name": "cpc_on_value",
	"action": {"http_url": "https://example.com/hook", "http_method": "post"},
	"simple_triggers": [
		{
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "org.astarte-platform.genericsensors.Values",
			"interface_major": 1,
			"match_path": "/*",
			"value_match_operator": "*"
		}
	]
}`

// fixture is the composed Astrate instance under test.
type fixture struct {
	st        *store.Store
	e         *engine.Engine
	broker    *broker.Broker
	pairer    *pairing.Service
	server    *httptest.Server
	realmName string
	roots     *x509.CertPool
	sslURL    string
	realmKey  string // path to the realm JWT private key (PEM)
	realmPub  string // path to the realm JWT public key (PEM)
	instKey   string // path to the instance-admin (housekeeping) private key (PEM)
}

// dialDSN finds a reachable database: ASTRATE_TEST_DSN first, then compose.
func dialDSN(t *testing.T) string {
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
	t.Skip("CP-C needs a database: run `make up` at the repo root or set ASTRATE_TEST_DSN")
	return ""
}

// writePEM writes a PEM file under dir and returns its path.
func writePEM(t *testing.T, dir, name string, block *pem.Block) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, pem.EncodeToMemory(block), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

// rsaKeyFiles generates an RSA keypair and writes the PKCS#1 private key and
// PKIX public key as PEM files, returning (privPath, pubPath, pubPEM).
func rsaKeyFiles(t *testing.T, dir, stem string) (string, string, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	priv := writePEM(t, dir, stem+"_private.pem", &pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	pubBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}
	pub := writePEM(t, dir, stem+"_public.pem", pubBlock)
	return priv, pub, string(pem.EncodeToMemory(pubBlock))
}

// newFixture wires a running engine + broker + pairing API + the three M7 HTTP
// surfaces + the live-stream socket on one httptest server, exactly as
// cmd/astrate will (M8). The realm is *not* pre-created: astartectl creates it.
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

	keyDir := t.TempDir()
	realmKey, realmPub, _ := rsaKeyFiles(t, keyDir, "realm")
	instKey, _, instPubPEM := rsaKeyFiles(t, keyDir, "instance")

	e, err := engine.New(st, nil, engine.Config{
		Shards: 4, BatchMaxRows: 8, BatchMaxWait: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("engine.New: %v", err)
	}
	serverCert, roots := testutil.ServerTLSCert(t)
	b, err := broker.New(ctx, broker.Config{
		TLSAddr:          "127.0.0.1:0",
		ServerTLSCert:    serverCert,
		SessionStorePath: filepath.Join(t.TempDir(), "sessions.db"),
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
		BrokerURL: "mqtts://" + b.TLSAddr(),
		Version:   "1.2.0",
	})

	// All five M7/M4 surfaces on one mux, as M8's cmd/astrate mounts them.
	mw := auth.NewMiddleware(st)
	mux := http.NewServeMux()
	pairing.NewAPI(pairer, mw, pairing.APIConfig{}).Mount(mux)
	housekeeping.NewAPI(housekeeping.NewService(st, sealer, b, nil), mw, []string{instPubPEM}).Mount(mux)
	realm.NewAPI(realm.NewService(st, e, nil), mw).Mount(mux)
	appengine.NewAPI(appengine.NewService(st, e, nil), mw).Mount(mux)
	apstream.NewAPI(e.Bus(), mw).Mount(mux)

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	var suffix [4]byte
	_, _ = rand.Read(suffix[:])
	realmName := "cpc" + hexLower(suffix[:])

	return &fixture{
		st: st, e: e, broker: b, pairer: pairer, server: server, realmName: realmName,
		roots: roots, sslURL: "ssl://" + b.TLSAddr(),
		realmKey: realmKey, realmPub: realmPub, instKey: instKey,
	}
}

// hexLower renders bytes as lowercase hex (a realm name must be lowercase and
// start with a letter — the "cpc" prefix guarantees the latter).
func hexLower(b []byte) string {
	const hexdigits = "0123456789abcdef"
	out := make([]byte, 0, len(b)*2)
	for _, x := range b {
		out = append(out, hexdigits[x>>4], hexdigits[x&0x0f])
	}
	return string(out)
}

// TestCPC drives astartectl through the full M7 operator flow.
func TestCPC(t *testing.T) {
	bin := ensureAstartectl(t)
	f := newFixture(t)
	ctx := context.Background()
	url := f.server.URL

	// 1. Housekeeping: create the realm with astartectl (instance-admin JWT).
	runAstartectl(t, bin,
		"housekeeping", "realms", "create", f.realmName,
		"--housekeeping-key", f.instKey,
		"--realm-public-key", f.realmPub,
		"-u", url, "-y",
	)
	realmRow, err := f.st.GetRealmByName(ctx, f.realmName)
	if err != nil {
		t.Fatalf("realm not created by astartectl: %v", err)
	}
	// The broker only learns the new realm's CA on reload (M8 wires this to the
	// housekeeping mutation; here the harness drives it explicitly).
	if err := f.broker.ReloadRealms(ctx); err != nil {
		t.Fatalf("broker ReloadRealms: %v", err)
	}

	// 2. Realm Management: install the interface (realm JWT).
	ifaceFile := filepath.Join(t.TempDir(), "Values.json")
	if err := os.WriteFile(ifaceFile, []byte(valuesInterface), 0o600); err != nil {
		t.Fatal(err)
	}
	runAstartectl(t, bin,
		"realm-management", "interfaces", "install", ifaceFile,
		"--realm-key", f.realmKey, "-r", f.realmName, "-u", url,
	)
	si, err := f.st.GetInterface(ctx, realmRow.ID, valuesName, 1)
	if err != nil {
		t.Fatalf("interface not installed: %v", err)
	}

	// 3. Pairing: register the device with astartectl (realm JWT → a_pa).
	id, err := deviceid.Random()
	if err != nil {
		t.Fatal(err)
	}
	secret := runAstartectl(t, bin,
		"pairing", "agent", "register", id.String(),
		"--pairing-url", url+"/pairing",
		"--realm-name", f.realmName,
		"--realm-key", f.realmKey,
		"--compact-output",
	)
	if len(secret) != 44 {
		t.Fatalf("astartectl credentials secret length: got %d, want 44 (%q)", len(secret), secret)
	}

	// 4. A real device connects (in-process credentials) and publishes a sample.
	key, csr := testutil.DeviceCSR(t)
	crt, err := f.pairer.Credentials(ctx, f.realmName, id.String(), secret, csr, netip.MustParseAddr("127.0.0.1"))
	if err != nil {
		t.Fatalf("Credentials: %v", err)
	}
	tlsCfg := testutil.DeviceTLSConfig(t, crt, key, f.roots)
	dev := testutil.ConnectAstarteDevice(t, f.sslURL, f.realmName, id, tlsCfg, true)
	t.Cleanup(dev.Disconnect)
	dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{valuesName: {1, 0}}))
	waitFor(t, 5*time.Second, "introspection persisted", func() bool {
		d, err := f.st.GetDevice(ctx, realmRow.ID, id)
		return err == nil && len(d.Introspection) == 1
	})

	sampleTS := time.Now().UTC().Truncate(time.Millisecond)
	dev.PublishValue(t, valuesName, "/sensor0/value", 42.5, &sampleTS, payload.FormatBSON, 1)
	waitFor(t, 10*time.Second, "sample persisted", func() bool {
		rows, err := f.st.Series(ctx, store.SeriesQuery{
			RealmID: realmRow.ID, DeviceID: id, InterfaceID: si.ID, Path: "/sensor0/value",
		})
		return err == nil && len(rows) == 1 && rows[0].ValueDouble != nil && *rows[0].ValueDouble == 42.5
	})

	// 5. AppEngine: astartectl must read the published value back. data-snapshot
	// exercises the nested interface-root snapshot shape; get-samples the
	// concrete-path array. A non-zero astartectl exit (parse failure) fails the
	// test inside runAstartectl.
	t.Run("DataSnapshot", func(t *testing.T) {
		out := runAstartectl(t, bin,
			"appengine", "devices", "data-snapshot", id.String(), valuesName,
			"--realm-key", f.realmKey, "-r", f.realmName, "-u", url, "-o", "json",
		)
		if !strings.Contains(out, "42.5") || !strings.Contains(out, "sensor0") {
			t.Fatalf("data-snapshot did not return the published sample:\n%s", out)
		}
	})

	t.Run("GetSamples", func(t *testing.T) {
		out := runAstartectl(t, bin,
			"appengine", "devices", "get-samples", id.String(), valuesName, "/sensor0/value",
			"--realm-key", f.realmKey, "-r", f.realmName, "-u", url, "-o", "json", "-c", "10",
		)
		if !strings.Contains(out, "42.5") {
			t.Fatalf("get-samples did not return the published sample:\n%s", out)
		}
	})

	t.Run("DeviceList", func(t *testing.T) {
		out := runAstartectl(t, bin,
			"appengine", "devices", "list",
			"--realm-key", f.realmKey, "-r", f.realmName, "-u", url,
		)
		if !strings.Contains(out, id.String()) {
			t.Fatalf("device list did not include %s:\n%s", id.String(), out)
		}
	})

	// 5b. AppEngine device-metadata and groups (CP-C superset, ROADMAP §9.5):
	// aliases + attributes PATCH round-trips and group membership via astartectl.
	t.Run("DeviceMetadataAndGroups", func(t *testing.T) {
		dev := id.String()
		ae := func(args ...string) string {
			return runAstartectl(t, bin, append(append([]string{"appengine"}, args...),
				"--realm-key", f.realmKey, "-r", f.realmName, "-u", url)...)
		}
		ae("devices", "aliases", "add", dev, "room=lab1")
		if out := ae("devices", "aliases", "list", dev); !strings.Contains(out, "lab1") {
			t.Fatalf("alias not listed after add:\n%s", out)
		}
		ae("devices", "aliases", "remove", dev, "room")

		ae("devices", "attributes", "set", dev, "owner=ops")
		if out := ae("devices", "attributes", "list", dev); !strings.Contains(out, "ops") {
			t.Fatalf("attribute not listed after set:\n%s", out)
		}

		ae("groups", "create", "cpc-group", dev)
		if out := ae("groups", "list"); !strings.Contains(out, "cpc-group") {
			t.Fatalf("group not listed after create:\n%s", out)
		}
	})

	// 6. Realm Management trigger CRUD round-trip.
	t.Run("TriggerCRUD", func(t *testing.T) {
		triggerFile := filepath.Join(t.TempDir(), "trigger.json")
		if err := os.WriteFile(triggerFile, []byte(triggerJSON), 0o600); err != nil {
			t.Fatal(err)
		}
		runAstartectl(t, bin,
			"realm-management", "triggers", "install", triggerFile,
			"--realm-key", f.realmKey, "-r", f.realmName, "-u", url,
		)
		list := runAstartectl(t, bin,
			"realm-management", "triggers", "list",
			"--realm-key", f.realmKey, "-r", f.realmName, "-u", url,
		)
		if !strings.Contains(list, "cpc_on_value") {
			t.Fatalf("trigger not listed:\n%s", list)
		}
		show := runAstartectl(t, bin,
			"realm-management", "triggers", "show", "cpc_on_value",
			"--realm-key", f.realmKey, "-r", f.realmName, "-u", url,
		)
		if !strings.Contains(show, "cpc_on_value") {
			t.Fatalf("trigger show did not return the trigger:\n%s", show)
		}
		runAstartectl(t, bin,
			"realm-management", "triggers", "delete", "cpc_on_value",
			"--realm-key", f.realmKey, "-r", f.realmName, "-u", url,
		)
		afterList := runAstartectl(t, bin,
			"realm-management", "triggers", "list",
			"--realm-key", f.realmKey, "-r", f.realmName, "-u", url,
		)
		if strings.Contains(afterList, "cpc_on_value") {
			t.Fatalf("trigger still listed after delete:\n%s", afterList)
		}
	})
}

// waitFor polls cond until it returns true or the deadline passes.
func waitFor(t *testing.T, timeout time.Duration, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("timed out after %s waiting for %s", timeout, what)
}
