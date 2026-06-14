// Package cpb is conformance checkpoint CP-B (docs/ROADMAP.md §0.3, §7.3):
// the wire-frozen proof that Astrate runs the full device loop for the
// *unmodified* official astarte-device-sdk-go. It composes the real stack —
// store + pairing API + embedded MQTT broker + ingestion engine — in-process
// against a live TimescaleDB and drives it with the official SDK device:
//
//	register (astarte-go agent) → mTLS connect → introspection →
//	device datastreams (individual + object, BSON) →
//	device properties (set + unset) →
//	server-owned datastream + property delivery →
//	emptyCache resync of a staged server property.
//
// Each step is cross-checked against the database rows Astrate persisted and
// the callbacks the SDK fired. The database is the compose service
// (`make up`) or ASTRATE_TEST_DSN; the suite skips with instructions when
// neither is reachable.
package cpb

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	sdkdevice "github.com/astarte-platform/astarte-device-sdk-go/device"
	astarteclient "github.com/astarte-platform/astarte-go/client"
	"github.com/astarte-platform/astarte-go/interfaces"
	"github.com/astarte-platform/astarte-go/misc"
	"github.com/jackc/pgx/v5"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/engine"
	"github.com/astrate-platform/astrate/internal/pairing"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// composeDSN is the docker-compose database (`make up` at the repo root).
const composeDSN = "postgres://astrate:astrate@127.0.0.1:5432/astrate?sslmode=disable"

// Interfaces installed in Astrate and registered with the SDK device.
const (
	ifSensor     = "org.astrate.cpb.SensorValues" // device datastream, individual
	ifCoords     = "org.astrate.cpb.Coordinates"  // device datastream, object
	ifDeviceConf = "org.astrate.cpb.DeviceConf"   // device properties
	ifServerData = "org.astrate.cpb.ServerData"   // server datastream, individual
	ifServerConf = "org.astrate.cpb.ServerConf"   // server properties
)

var cpbInterfaceDefs = map[string]string{
	ifSensor: `{
		"interface_name": "` + ifSensor + `", "version_major": 1, "version_minor": 0,
		"type": "datastream", "ownership": "device",
		"mappings": [{"endpoint": "/value", "type": "double"}]
	}`,
	ifCoords: `{
		"interface_name": "` + ifCoords + `", "version_major": 1, "version_minor": 0,
		"type": "datastream", "ownership": "device", "aggregation": "object",
		"mappings": [{"endpoint": "/coords/latitude", "type": "double"}, {"endpoint": "/coords/longitude", "type": "double"}]
	}`,
	ifDeviceConf: `{
		"interface_name": "` + ifDeviceConf + `", "version_major": 1, "version_minor": 0,
		"type": "properties", "ownership": "device",
		"mappings": [{"endpoint": "/%{key}", "type": "string", "allow_unset": true}]
	}`,
	ifServerData: `{
		"interface_name": "` + ifServerData + `", "version_major": 1, "version_minor": 0,
		"type": "datastream", "ownership": "server",
		"mappings": [{"endpoint": "/value", "type": "double"}]
	}`,
	ifServerConf: `{
		"interface_name": "` + ifServerConf + `", "version_major": 1, "version_minor": 0,
		"type": "properties", "ownership": "server",
		"mappings": [{"endpoint": "/%{key}", "type": "integer", "allow_unset": true}]
	}`,
}

// fixture is the composed Astrate instance under test.
type fixture struct {
	st         *store.Store
	e          *engine.Engine
	realm      *store.Realm
	realmName  string
	pairingURL string
	jwtKeyPEM  []byte
	ifaces     map[string]*store.StoredInterface
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
	t.Skip("CP-B needs a database: run `make up` at the repo root or set ASTRATE_TEST_DSN")
	return ""
}

// newFixture provisions a realm, installs the interfaces, and wires a running
// engine + broker + pairing API exactly as cmd/astrate will (M8).
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

	realmName := "cpb" + strconv.FormatInt(time.Now().UnixNano(), 36)
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

	ifaces := make(map[string]*store.StoredInterface, len(cpbInterfaceDefs))
	for name, def := range cpbInterfaceDefs {
		si, err := st.InstallInterface(ctx, realm.ID, []byte(def))
		if err != nil {
			t.Fatalf("InstallInterface(%s): %v", name, err)
		}
		ifaces[name] = si
	}

	// Engine + broker, wired to each other.
	e, err := engine.New(st, nil, engine.Config{
		Shards: 4, BatchMaxRows: 8, BatchMaxWait: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("engine.New: %v", err)
	}
	serverCert, _ := testutil.ServerTLSCert(t)
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

	// The broker's real TLS address feeds the pairing info endpoint, so the
	// SDK learns where to connect (the SDK rewrites mqtts:// to ssl://).
	svc := pairing.New(st, sealer, pairing.Config{
		BrokerURL: "mqtts://" + b.TLSAddr(),
		Version:   "1.2.0",
	})
	api := pairing.NewAPI(svc, auth.NewMiddleware(st), pairing.APIConfig{})
	mux := http.NewServeMux()
	api.Mount(mux)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return &fixture{
		st: st, e: e, realm: realm, realmName: realmName,
		pairingURL: server.URL + "/pairing", jwtKeyPEM: jwtKeyPEM, ifaces: ifaces,
	}
}

// register obtains a credentials secret through the official agent client.
func (f *fixture) register(t *testing.T, deviceID string) string {
	t.Helper()
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
	return secret
}

// received is one server→device message captured by the SDK callbacks.
type received struct {
	iface string
	path  string
	value interface{}
}

// inbox collects SDK callback deliveries.
type inbox struct {
	mu   sync.Mutex
	msgs []received
}

func (in *inbox) add(r received) {
	in.mu.Lock()
	defer in.mu.Unlock()
	in.msgs = append(in.msgs, r)
}

// waitFor polls captured messages for one satisfying pred.
func (in *inbox) waitFor(t *testing.T, timeout time.Duration, what string, pred func(received) bool) received {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		in.mu.Lock()
		for _, m := range in.msgs {
			if pred(m) {
				in.mu.Unlock()
				return m
			}
		}
		in.mu.Unlock()
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("timed out after %s waiting for %s", timeout, what)
	return received{}
}

// connectSDK builds an SDK device with every interface and the capture
// callbacks, connects it (failing the test on error), and registers cleanup.
func (f *fixture) connectSDK(t *testing.T, deviceID, secret, cryptoDir string, in *inbox) *sdkdevice.Device {
	t.Helper()
	if err := os.MkdirAll(cryptoDir, 0o700); err != nil {
		t.Fatalf("crypto dir: %v", err)
	}
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
	for name := range cpbInterfaceDefs {
		iface, err := interfaces.ParseInterface([]byte(cpbInterfaceDefs[name]))
		if err != nil {
			t.Fatalf("ParseInterface(%s): %v", name, err)
		}
		if err := d.AddInterface(iface); err != nil {
			t.Fatalf("AddInterface(%s): %v", name, err)
		}
	}
	d.OnIndividualMessageReceived = func(_ *sdkdevice.Device, m sdkdevice.IndividualMessage) {
		in.add(received{iface: m.Interface.Name, path: m.Path, value: m.Value})
	}
	d.OnAggregateMessageReceived = func(_ *sdkdevice.Device, m sdkdevice.AggregateMessage) {
		in.add(received{iface: m.Interface.Name, path: m.Path, value: m.Values})
	}

	ch := make(chan error, 1)
	d.Connect(ch)
	select {
	case err := <-ch:
		if err != nil {
			t.Fatalf("SDK Connect: %v", err)
		}
	case <-time.After(90 * time.Second):
		t.Fatal("SDK Connect did not settle within 90s")
	}
	t.Cleanup(func() {
		dc := make(chan error, 1)
		d.Disconnect(dc)
		select {
		case <-dc:
		case <-time.After(10 * time.Second):
		}
	})
	return d
}

// TestCPB drives the official Go SDK through the full device loop against the
// composed Astrate instance.
func TestCPB(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	id, err := deviceid.Random()
	if err != nil {
		t.Fatal(err)
	}
	deviceID := id.String()

	secret := f.register(t, deviceID)
	if len(secret) != 44 {
		t.Fatalf("credentials secret length: got %d, want 44", len(secret))
	}

	// Stage a server-owned property before the device connects: the SDK sends
	// emptyCache on connect, so Astrate must resync it onto the device.
	if err := f.e.PublishServerValue(ctx, f.realmName, id, ifServerConf, "/staged",
		json.RawMessage("7"), nil); err != nil {
		t.Fatalf("staging server property: %v", err)
	}

	in := &inbox{}
	d := f.connectSDK(t, deviceID, secret, filepath.Join(t.TempDir(), "crypto"), in)

	t.Run("Introspection", func(t *testing.T) {
		waitFor(t, 10*time.Second, func() bool {
			dev, err := f.st.GetDevice(ctx, f.realm.ID, id)
			return err == nil && len(dev.Introspection) == len(cpbInterfaceDefs)
		}, "introspection persisted")
		dev, err := f.st.GetDevice(ctx, f.realm.ID, id)
		if err != nil {
			t.Fatalf("GetDevice: %v", err)
		}
		if v, ok := dev.Introspection[ifSensor]; !ok || v.Major != 1 {
			t.Errorf("introspection missing %s:1 (got %+v)", ifSensor, dev.Introspection)
		}
	})

	t.Run("EmptyCacheResyncServerProperty", func(t *testing.T) {
		m := in.waitFor(t, 15*time.Second, "staged server property on the device", func(r received) bool {
			return r.iface == ifServerConf && r.path == "/staged"
		})
		if n, ok := asNumber(m.value); !ok || n != 7 {
			t.Errorf("staged property value = %v (%T), want 7", m.value, m.value)
		}
	})

	t.Run("IndividualDatastream", func(t *testing.T) {
		if err := d.SendIndividualMessage(ifSensor, "/value", 42.5); err != nil {
			t.Fatalf("SendIndividualMessage: %v", err)
		}
		waitFor(t, 10*time.Second, func() bool {
			rows := f.series(t, id, ifSensor, "/value")
			return len(rows) == 1 && rows[0].ValueDouble != nil && *rows[0].ValueDouble == 42.5
		}, "individual datastream row")
	})

	t.Run("ObjectDatastream", func(t *testing.T) {
		if err := d.SendAggregateMessage(ifCoords, "/coords",
			map[string]interface{}{"latitude": 45.07, "longitude": 7.69}); err != nil {
			t.Fatalf("SendAggregateMessage: %v", err)
		}
		waitFor(t, 10*time.Second, func() bool {
			rows, err := f.st.ObjectSeries(ctx, store.SeriesQuery{
				RealmID: f.realm.ID, DeviceID: id, InterfaceID: f.ifaces[ifCoords].ID, Path: "/coords",
			})
			return err == nil && len(rows) == 1 &&
				containsAll(string(rows[0].Value), "45.07", "7.69")
		}, "object datastream row")
	})

	t.Run("DeviceProperties", func(t *testing.T) {
		ifaceID := f.ifaces[ifDeviceConf].ID
		if err := d.SetProperty(ifDeviceConf, "/mode", "eco"); err != nil {
			t.Fatalf("SetProperty: %v", err)
		}
		waitFor(t, 10*time.Second, func() bool {
			p, err := f.st.GetProperty(ctx, f.realm.ID, id, ifaceID, "/mode")
			return err == nil && string(p.Value) == `"eco"`
		}, "device property set")

		if err := d.UnsetProperty(ifDeviceConf, "/mode"); err != nil {
			t.Fatalf("UnsetProperty: %v", err)
		}
		waitFor(t, 10*time.Second, func() bool {
			_, err := f.st.GetProperty(ctx, f.realm.ID, id, ifaceID, "/mode")
			return errors.Is(err, store.ErrNotFound)
		}, "device property unset")
	})

	t.Run("ServerOwnedDatastreamDelivery", func(t *testing.T) {
		if err := f.e.PublishServerValue(ctx, f.realmName, id, ifServerData, "/value",
			json.RawMessage("9.5"), nil); err != nil {
			t.Fatalf("PublishServerValue: %v", err)
		}
		m := in.waitFor(t, 15*time.Second, "server datastream on the device", func(r received) bool {
			return r.iface == ifServerData && r.path == "/value"
		})
		if n, ok := asNumber(m.value); !ok || n != 9.5 {
			t.Errorf("server datastream value = %v (%T), want 9.5", m.value, m.value)
		}
	})
}

func (f *fixture) series(t *testing.T, id deviceid.ID, iface, path string) []store.IndividualRow {
	t.Helper()
	rows, err := f.st.Series(context.Background(), store.SeriesQuery{
		RealmID: f.realm.ID, DeviceID: id, InterfaceID: f.ifaces[iface].ID, Path: path,
	})
	if err != nil {
		t.Fatalf("Series(%s%s): %v", iface, path, err)
	}
	return rows
}

// waitFor polls cond until it returns true or the deadline passes.
func waitFor(t *testing.T, timeout time.Duration, cond func() bool, what string) {
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

// asNumber widens the numeric types the SDK's BSON decoder may produce.
func asNumber(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

// containsAll reports whether s contains every substring.
func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		found := false
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
