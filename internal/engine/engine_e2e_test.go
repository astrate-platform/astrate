//go:build e2e

package engine

import (
	"bytes"
	"compress/zlib"
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/netip"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/pairing"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// Interface names installed in the e2e realm.
const (
	e2eSensors = "org.astrate.e2e.SensorValues" // device datastream, individual
	e2eGeo     = "org.astrate.e2e.Geolocation"  // device datastream, object
	e2eDevCfg  = "org.astrate.e2e.DeviceConfig" // device properties (parametric)
	e2eSrvData = "org.astrate.e2e.ServerData"   // server datastream, individual
	e2eSrvCfg  = "org.astrate.e2e.ServerConfig" // server properties (parametric)
)

var e2eInterfaceDefs = map[string]string{
	e2eSensors: `{
		"interface_name": "` + e2eSensors + `", "version_major": 1, "version_minor": 0,
		"type": "datastream", "ownership": "device",
		"mappings": [{"endpoint": "/value", "type": "double"}, {"endpoint": "/count", "type": "integer"}]
	}`,
	e2eGeo: `{
		"interface_name": "` + e2eGeo + `", "version_major": 1, "version_minor": 0,
		"type": "datastream", "ownership": "device", "aggregation": "object",
		"mappings": [{"endpoint": "/gps/latitude", "type": "double"}, {"endpoint": "/gps/longitude", "type": "double"}]
	}`,
	e2eDevCfg: `{
		"interface_name": "` + e2eDevCfg + `", "version_major": 1, "version_minor": 0,
		"type": "properties", "ownership": "device",
		"mappings": [{"endpoint": "/%{key}", "type": "string", "allow_unset": true}]
	}`,
	e2eSrvData: `{
		"interface_name": "` + e2eSrvData + `", "version_major": 1, "version_minor": 0,
		"type": "datastream", "ownership": "server",
		"mappings": [{"endpoint": "/value", "type": "double"}]
	}`,
	e2eSrvCfg: `{
		"interface_name": "` + e2eSrvCfg + `", "version_major": 1, "version_minor": 0,
		"type": "properties", "ownership": "server",
		"mappings": [{"endpoint": "/%{param}", "type": "integer", "allow_unset": true}]
	}`,
}

// e2eEnv is the shared T3 environment: TimescaleDB, store, realm with a
// pairing-provisioned CA, the installed interfaces, and a running broker+engine
// pair wired to each other exactly as M8 will wire them.
type e2eEnv struct {
	st     *store.Store
	e      *Engine
	b      *broker.Broker
	svc    *pairing.Service
	realm  *store.Realm
	roots  *x509.CertPool
	sslURL string
	ifaces map[string]*store.StoredInterface
}

func newE2EEnv(t *testing.T) *e2eEnv {
	t.Helper()
	ctx := context.Background()

	pool := testutil.StartTimescale(t)
	st, err := store.New(ctx, pool.Config().ConnString())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(st.Close)

	master := make([]byte, 32)
	if _, err := rand.Read(master); err != nil {
		t.Fatalf("master key: %v", err)
	}
	sealer, err := store.NewKeySealer(master)
	if err != nil {
		t.Fatalf("NewKeySealer: %v", err)
	}

	suffix := make([]byte, 4)
	if _, err := rand.Read(suffix); err != nil {
		t.Fatalf("realm suffix: %v", err)
	}
	realmName := "enge2e" + hex.EncodeToString(suffix)

	certPEM, sealedKey, err := pairing.ProvisionCA(realmName, 0, sealer)
	if err != nil {
		t.Fatalf("ProvisionCA: %v", err)
	}
	realm, err := st.CreateRealm(ctx, store.NewRealm{
		Name: realmName, CACertificatePEM: certPEM, CAPrivateKeySealed: sealedKey,
	})
	if err != nil {
		t.Fatalf("CreateRealm: %v", err)
	}

	env := &e2eEnv{
		st: st, realm: realm,
		svc:    pairing.New(st, sealer, pairing.Config{BrokerURL: "mqtts://localhost:8883"}),
		ifaces: make(map[string]*store.StoredInterface),
	}
	for name, def := range e2eInterfaceDefs {
		si, err := st.InstallInterface(ctx, realm.ID, []byte(def))
		if err != nil {
			t.Fatalf("InstallInterface(%s): %v", name, err)
		}
		env.ifaces[name] = si
	}

	// Wire broker + engine the way cmd/astrate will (docs/ROADMAP.md §7.2
	// file 6.14): the engine is the broker's intake and lifecycle sink; the
	// broker is the engine's publish port.
	e, err := New(st, nil, Config{
		Shards: 4, BatchMaxRows: 8, BatchMaxWait: 20 * time.Millisecond, Logger: discardLogger(),
	})
	if err != nil {
		t.Fatalf("engine.New: %v", err)
	}
	serverCert, roots := testutil.ServerTLSCert(t)
	b, err := broker.New(ctx, broker.Config{
		TLSAddr:          "127.0.0.1:0",
		ServerTLSCert:    serverCert,
		SessionStorePath: filepath.Join(t.TempDir(), "sessions.db"),
		Logger:           discardLogger(),
	}, st, e, e)
	if err != nil {
		t.Fatalf("broker.New: %v", err)
	}
	e.AttachBroker(AdaptBroker(b))

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
		_ = b.Close() // stop accepting first so no submit blocks
		dctx, dcancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer dcancel()
		if err := e.Drain(dctx); err != nil {
			t.Errorf("engine.Drain: %v", err)
		}
		cancel()
	})

	env.e, env.b, env.roots = e, b, roots
	env.sslURL = "ssl://" + b.TLSAddr()
	return env
}

// connectDevice registers, credentials, and connects a fresh device, then
// publishes an introspection declaring every interface and waits for it to be
// persisted (so subsequent data publishes pass the broker ACL and the engine's
// introspection gate).
func (env *e2eEnv) connectDevice(t *testing.T) *testutil.AstarteDevice {
	t.Helper()
	ctx := context.Background()

	id, err := deviceid.Random()
	if err != nil {
		t.Fatalf("deviceid.Random: %v", err)
	}
	secret, err := env.svc.Register(ctx, env.realm.Name, id.String(), "")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	key, csrPEM := testutil.DeviceCSR(t)
	certPEM, err := env.svc.Credentials(ctx, env.realm.Name, id.String(), secret, csrPEM, netip.MustParseAddr("127.0.0.1"))
	if err != nil {
		t.Fatalf("Credentials: %v", err)
	}
	tlsCfg := testutil.DeviceTLSConfig(t, certPEM, key, env.roots)

	dev := testutil.ConnectAstarteDevice(t, env.sslURL, env.realm.Name, id, tlsCfg, true)
	t.Cleanup(dev.Disconnect)

	dev.PublishIntrospection(t, testutil.Introspection(map[string][2]int{
		e2eSensors: {1, 0}, e2eGeo: {1, 0}, e2eDevCfg: {1, 0}, e2eSrvData: {1, 0}, e2eSrvCfg: {1, 0},
	}))
	waitFor(t, 5*time.Second, "introspection persisted", func() bool {
		d, err := env.st.GetDevice(ctx, env.realm.ID, id)
		return err == nil && len(d.Introspection) == 5
	})
	return dev
}

// decodeIndividual decodes a server→device individual payload against the
// engine's own compiled mapping (the test runs in-package, so it can reach the
// live schema snapshot).
func (env *e2eEnv) decodeIndividual(t *testing.T, iface string, major int, path string, body []byte) payload.DecodedPayload {
	t.Helper()
	rs := env.e.schemas.realm(env.realm.Name)
	if rs == nil {
		t.Fatalf("realm %s not in schema snapshot", env.realm.Name)
	}
	ci := rs.iface(iface, major)
	if ci == nil {
		t.Fatalf("interface %s:%d not compiled", iface, major)
	}
	m, ok := ci.Trie.Match(path)
	if !ok {
		t.Fatalf("path %s matches no endpoint of %s", path, iface)
	}
	dp, err := payload.Decoder{MaxSize: 1 << 16}.Individual(body, m)
	if err != nil {
		t.Fatalf("decoding %s%s: %v", iface, path, err)
	}
	return dp
}

func (env *e2eEnv) series(t *testing.T, dev *testutil.AstarteDevice, iface, path string) []store.IndividualRow {
	t.Helper()
	rows, err := env.st.Series(context.Background(), store.SeriesQuery{
		RealmID: env.realm.ID, DeviceID: dev.ID, InterfaceID: env.ifaces[iface].ID, Path: path,
	})
	if err != nil {
		t.Fatalf("Series(%s%s): %v", iface, path, err)
	}
	return rows
}

// TestEngineE2E is the M6b component suite (docs/ROADMAP.md §7.3 T3): a real
// device drives the embedded broker, which feeds the engine, which persists to
// TimescaleDB and publishes server-owned and control traffic back. One
// container/realm/broker/engine; one fresh device per subtest.
func TestEngineE2E(t *testing.T) {
	env := newE2EEnv(t)

	t.Run("IndividualAndObjectData", func(t *testing.T) { testE2EData(t, env) })
	t.Run("DeviceProperties", func(t *testing.T) { testE2EDeviceProperties(t, env) })
	t.Run("ServerOwnedDelivery", func(t *testing.T) { testE2EServerDelivery(t, env) })
	t.Run("EmptyCacheResync", func(t *testing.T) { testE2EEmptyCache(t, env) })
	t.Run("ProducerPropertiesPurge", func(t *testing.T) { testE2EProducerProperties(t, env) })
}

// testE2EData: BSON individual, JSON individual (flips the format hint), and
// BSON object-aggregated publishes all land in the right columns. QoS >= 1
// publishes are PUBACKed only after the engine commits, so a returned publish
// token means the row is already persisted.
func testE2EData(t *testing.T, env *e2eEnv) {
	ctx := context.Background()
	dev := env.connectDevice(t)

	dev.PublishValue(t, e2eSensors, "/value", 22.5, nil, payload.FormatBSON, 1)
	if rows := env.series(t, dev, e2eSensors, "/value"); len(rows) != 1 || rows[0].ValueDouble == nil || *rows[0].ValueDouble != 22.5 {
		t.Fatalf("/value rows = %+v", rows)
	}

	dev.PublishValue(t, e2eSensors, "/count", int32(7), nil, payload.FormatJSON, 1)
	if rows := env.series(t, dev, e2eSensors, "/count"); len(rows) != 1 || rows[0].ValueInteger == nil || *rows[0].ValueInteger != 7 {
		t.Fatalf("/count rows = %+v", rows)
	}
	// The JSON device publish flipped the sticky payload-format hint (§3.5.4).
	waitFor(t, 5*time.Second, "format hint flipped to json", func() bool {
		d, err := env.st.GetDevice(ctx, env.realm.ID, dev.ID)
		return err == nil && d.PayloadFormatHint == "json"
	})

	dev.PublishValue(t, e2eGeo, "/gps", map[string]payload.Value{"latitude": 45.07, "longitude": 7.69},
		nil, payload.FormatBSON, 2)
	rows, err := env.st.ObjectSeries(ctx, store.SeriesQuery{
		RealmID: env.realm.ID, DeviceID: dev.ID, InterfaceID: env.ifaces[e2eGeo].ID, Path: "/gps",
	})
	if err != nil {
		t.Fatalf("ObjectSeries: %v", err)
	}
	if len(rows) != 1 || !bytes.Contains(rows[0].Value, []byte("45.07")) || !bytes.Contains(rows[0].Value, []byte("7.69")) {
		t.Fatalf("object rows = %+v", rows)
	}
}

// testE2EDeviceProperties: a device property set is an upsert; an empty payload
// is an unset.
func testE2EDeviceProperties(t *testing.T, env *e2eEnv) {
	ctx := context.Background()
	dev := env.connectDevice(t)
	ifaceID := env.ifaces[e2eDevCfg].ID

	dev.PublishValue(t, e2eDevCfg, "/alpha", "hello", nil, payload.FormatBSON, 1)
	p, err := env.st.GetProperty(ctx, env.realm.ID, dev.ID, ifaceID, "/alpha")
	if err != nil {
		t.Fatalf("GetProperty after set: %v", err)
	}
	if string(p.Value) != `"hello"` {
		t.Errorf("property value = %s, want \"hello\"", p.Value)
	}

	dev.PublishRaw(t, e2eDevCfg, "/alpha", nil, 1) // empty payload = unset
	waitFor(t, 5*time.Second, "property unset", func() bool {
		_, err := env.st.GetProperty(ctx, env.realm.ID, dev.ID, ifaceID, "/alpha")
		return errors.Is(err, store.ErrNotFound)
	})
}

// testE2EServerDelivery: a server-owned datastream value published through the
// engine is persisted and delivered to the subscribed device on its data topic.
func testE2EServerDelivery(t *testing.T, env *e2eEnv) {
	ctx := context.Background()
	dev := env.connectDevice(t)

	if err := env.e.PublishServerValue(ctx, env.realm.Name, dev.ID, e2eSrvData, "/value",
		json.RawMessage("3.5"), nil); err != nil {
		t.Fatalf("PublishServerValue: %v", err)
	}

	topic := dev.Base() + "/" + e2eSrvData + "/value"
	msg := dev.WaitForTopic(t, 5*time.Second, topic)
	if dp := env.decodeIndividual(t, e2eSrvData, 1, "/value", msg.Payload); dp.Value != 3.5 {
		t.Errorf("delivered value = %v, want 3.5", dp.Value)
	}

	rows, err := env.st.Series(ctx, store.SeriesQuery{
		RealmID: env.realm.ID, DeviceID: dev.ID, InterfaceID: env.ifaces[e2eSrvData].ID, Path: "/value",
	})
	if err != nil {
		t.Fatalf("Series: %v", err)
	}
	if len(rows) != 1 || rows[0].ValueDouble == nil || *rows[0].ValueDouble != 3.5 {
		t.Fatalf("server datastream rows = %+v", rows)
	}
}

// testE2EEmptyCache: after a server-owned property is set, control/emptyCache
// makes Astrate re-send that property on the data topic and publish the
// consumer/properties purge list — both decoded and asserted (docs/DESIGN.md
// §3.3–3.4).
func testE2EEmptyCache(t *testing.T, env *e2eEnv) {
	ctx := context.Background()
	dev := env.connectDevice(t)

	if err := env.e.PublishServerValue(ctx, env.realm.Name, dev.ID, e2eSrvCfg, "/brightness",
		json.RawMessage("42"), nil); err != nil {
		t.Fatalf("PublishServerValue: %v", err)
	}
	propTopic := dev.Base() + "/" + e2eSrvCfg + "/brightness"
	dev.WaitForTopic(t, 5*time.Second, propTopic) // initial retained set

	dev.EmptyCache(t)

	// The resync re-publishes the property value on its data topic...
	msg := dev.WaitForTopic(t, 5*time.Second, propTopic)
	if dp := env.decodeIndividual(t, e2eSrvCfg, 1, "/brightness", msg.Payload); dp.Value != int32(42) {
		t.Errorf("resent property value = %v (%T), want int32(42)", dp.Value, dp.Value)
	}

	// ...and publishes the consumer/properties purge list naming it.
	want := e2eSrvCfg + "/brightness"
	cp := dev.WaitForMessage(t, 5*time.Second, "consumer/properties listing "+want, func(m testutil.ServerMessage) bool {
		if m.Topic != dev.Base()+"/control/consumer/properties" {
			return false
		}
		for _, e := range controlEntries(m.Payload) {
			if e == want {
				return true
			}
		}
		return false
	})
	if got := controlEntries(cp.Payload); len(got) != 1 || got[0] != want {
		t.Errorf("consumer/properties entries = %v, want [%s]", got, want)
	}
}

// testE2EProducerProperties: control/producer/properties is the device's
// exhaustive list of held device-owned properties; Astrate purges every
// device-owned property not in it.
func testE2EProducerProperties(t *testing.T, env *e2eEnv) {
	ctx := context.Background()
	dev := env.connectDevice(t)
	ifaceID := env.ifaces[e2eDevCfg].ID

	for _, key := range []string{"/alpha", "/beta", "/gamma"} {
		dev.PublishValue(t, e2eDevCfg, key, "v"+key, nil, payload.FormatBSON, 1)
	}
	for _, key := range []string{"/alpha", "/beta", "/gamma"} {
		if _, err := env.st.GetProperty(ctx, env.realm.ID, dev.ID, ifaceID, key); err != nil {
			t.Fatalf("GetProperty(%s) before purge: %v", key, err)
		}
	}

	// Device reports it only still holds /alpha; /beta and /gamma must be purged.
	dev.SendProducerProperties(t, []string{e2eDevCfg + "/alpha"})

	waitFor(t, 5*time.Second, "non-listed properties purged", func() bool {
		_, errB := env.st.GetProperty(ctx, env.realm.ID, dev.ID, ifaceID, "/beta")
		_, errG := env.st.GetProperty(ctx, env.realm.ID, dev.ID, ifaceID, "/gamma")
		return errors.Is(errB, store.ErrNotFound) && errors.Is(errG, store.ErrNotFound)
	})
	if _, err := env.st.GetProperty(ctx, env.realm.ID, dev.ID, ifaceID, "/alpha"); err != nil {
		t.Errorf("listed property /alpha was purged: %v", err)
	}
}

// controlEntries inflates an Astarte control payload (4-byte size prefix + zlib)
// into its entry list, returning nil on any error — it runs inside poll
// predicates where a malformed frame must not abort the test.
func controlEntries(frame []byte) []string {
	if len(frame) < 4 {
		return nil
	}
	zr, err := zlib.NewReader(bytes.NewReader(frame[4:]))
	if err != nil {
		return nil
	}
	defer func() { _ = zr.Close() }()
	plain, err := io.ReadAll(zr)
	if err != nil || len(plain) == 0 {
		return nil
	}
	if int(binary.BigEndian.Uint32(frame[:4])) != len(plain) {
		return nil
	}
	return strings.Split(string(plain), ";")
}
