//go:build e2e

package broker

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"net/netip"
	"path/filepath"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/astrate-platform/astrate/internal/pairing"
	"github.com/astrate-platform/astrate/internal/pairing/ca"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// Interface fixtures installed in the test realm.
const (
	e2eDeviceIface = "org.astrate.e2e.DeviceData"
	e2eServerIface = "org.astrate.e2e.ServerData"
)

var e2eInterfaceDefs = []string{
	`{
		"interface_name": "` + e2eDeviceIface + `",
		"version_major": 1, "version_minor": 0,
		"type": "datastream", "ownership": "device",
		"mappings": [{"endpoint": "/value", "type": "double"}]
	}`,
	`{
		"interface_name": "` + e2eServerIface + `",
		"version_major": 1, "version_minor": 0,
		"type": "datastream", "ownership": "server",
		"mappings": [{"endpoint": "/value", "type": "double"}]
	}`,
}

// e2eEnv is the shared T3 environment: one TimescaleDB, one migrated store,
// one realm with a pairing-provisioned CA, and the broker's server TLS
// identity.
type e2eEnv struct {
	st         *store.Store
	sealer     *store.KeySealer
	svc        *pairing.Service
	realm      *store.Realm
	otherRealm *store.Realm
	otherCA    *ca.CA
	serverCert tls.Certificate
	roots      *x509.CertPool
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
	realmName := "brk" + hex.EncodeToString(suffix)

	certPEM, sealedKey, err := pairing.ProvisionCA(realmName, 0, sealer)
	if err != nil {
		t.Fatalf("ProvisionCA: %v", err)
	}
	realm, err := st.CreateRealm(ctx, store.NewRealm{
		Name:               realmName,
		CACertificatePEM:   certPEM,
		CAPrivateKeySealed: sealedKey,
	})
	if err != nil {
		t.Fatalf("CreateRealm: %v", err)
	}
	for _, def := range e2eInterfaceDefs {
		if _, err := st.InstallInterface(ctx, realm.ID, []byte(def)); err != nil {
			t.Fatalf("InstallInterface: %v", err)
		}
	}

	// A second realm whose CA signs the "wrong-realm" adversarial certs.
	otherCA, err := ca.Generate(realmName+"x", 0)
	if err != nil {
		t.Fatalf("generating other realm CA: %v", err)
	}
	otherSealed, err := sealer.Seal(otherCA.PrivateKeyDER())
	if err != nil {
		t.Fatalf("sealing other realm CA key: %v", err)
	}
	otherRealm, err := st.CreateRealm(ctx, store.NewRealm{
		Name:               realmName + "x",
		CACertificatePEM:   otherCA.CertificatePEM(),
		CAPrivateKeySealed: otherSealed,
	})
	if err != nil {
		t.Fatalf("CreateRealm (other): %v", err)
	}

	svc := pairing.New(st, sealer, pairing.Config{BrokerURL: "mqtts://localhost:8883"})
	serverCert, roots := testutil.ServerTLSCert(t)
	return &e2eEnv{
		st: st, sealer: sealer, svc: svc,
		realm: realm, otherRealm: otherRealm, otherCA: otherCA,
		serverCert: serverCert, roots: roots,
	}
}

// e2eDevice is a registered, credentialed device with seeded introspection.
type e2eDevice struct {
	identity Identity
	secret   string
	tlsCfg   *tls.Config
}

func (e *e2eEnv) newDevice(t *testing.T) *e2eDevice {
	t.Helper()
	ctx := context.Background()

	devID, err := deviceid.Random()
	if err != nil {
		t.Fatalf("deviceid.Random: %v", err)
	}
	secret, err := e.svc.Register(ctx, e.realm.Name, devID.String(), "")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	key, csrPEM := testutil.DeviceCSR(t)
	certPEM, err := e.svc.Credentials(ctx, e.realm.Name, devID.String(), secret, csrPEM, netip.MustParseAddr("127.0.0.1"))
	if err != nil {
		t.Fatalf("Credentials: %v", err)
	}
	if _, err := e.st.UpdateIntrospection(ctx, e.realm.ID, devID, map[string]store.InterfaceVersion{
		e2eDeviceIface: {Major: 1, Minor: 0},
		e2eServerIface: {Major: 1, Minor: 0},
	}); err != nil {
		t.Fatalf("UpdateIntrospection: %v", err)
	}
	return &e2eDevice{
		identity: Identity{Realm: e.realm.Name, DeviceID: devID},
		secret:   secret,
		tlsCfg:   testutil.DeviceTLSConfig(t, certPEM, key, e.roots),
	}
}

// realmCA opens the primary realm's CA for direct (adversarial) issuance.
func (e *e2eEnv) realmCA(t *testing.T) *ca.CA {
	t.Helper()
	der, err := e.sealer.Open(e.realm.CAPrivateKeySealed)
	if err != nil {
		t.Fatalf("opening realm CA key: %v", err)
	}
	realmCA, err := ca.Load(e.realm.CACertificatePEM, der)
	if err != nil {
		t.Fatalf("loading realm CA: %v", err)
	}
	return realmCA
}

type brokerOpts struct {
	path              string
	enforceLatestCert bool
	insecureDevMode   bool
	intake            Intake
	sink              LifecycleSink
}

func (e *e2eEnv) newBroker(t *testing.T, opts brokerOpts) *Broker {
	t.Helper()
	if opts.path == "" {
		opts.path = filepath.Join(t.TempDir(), "sessions.db")
	}
	if opts.intake == nil {
		opts.intake = newRecorderIntake(true)
	}
	b, err := New(context.Background(), Config{
		TLSAddr:           "127.0.0.1:0",
		ServerTLSCert:     e.serverCert,
		InsecureDevMode:   opts.insecureDevMode,
		DevAddr:           "127.0.0.1:0",
		SessionStorePath:  opts.path,
		EnforceLatestCert: opts.enforceLatestCert,
		Logger:            discardLogger(),
	}, e.st, opts.intake, opts.sink)
	if err != nil {
		t.Fatalf("broker.New: %v", err)
	}
	if err := b.Start(); err != nil {
		t.Fatalf("broker.Start: %v", err)
	}
	t.Cleanup(func() { _ = b.Close() })
	return b
}

func sslURL(b *Broker) string { return "ssl://" + b.TLSAddr() }

// TestBrokerE2E is the umbrella T3 suite (docs/ROADMAP.md §6 verification):
// one database container, one realm, one subtest per gate requirement.
func TestBrokerE2E(t *testing.T) {
	env := newE2EEnv(t)

	t.Run("ConnectAndLifecycle", func(t *testing.T) { testConnectAndLifecycle(t, env) })
	t.Run("Rejects", func(t *testing.T) { testRejects(t, env) })
	t.Run("DevModePlaintext", func(t *testing.T) { testDevModePlaintext(t, env) })
	t.Run("DeferredAck", func(t *testing.T) { testDeferredAck(t, env) })
	t.Run("SessionPersistenceAcrossRestart", func(t *testing.T) { testSessionPersistence(t, env) })
	t.Run("QoS2ExactlyOnce", func(t *testing.T) { testQoS2ExactlyOnce(t, env) })
	t.Run("OfflineMessageExpiry", func(t *testing.T) { testOfflineMessageExpiry(t, env) })
	t.Run("SessionTakeover", func(t *testing.T) { testSessionTakeover(t, env) })
}

func testConnectAndLifecycle(t *testing.T, env *e2eEnv) {
	ctx := context.Background()
	dev := env.newDevice(t)
	sink := &recorderSink{}
	b := env.newBroker(t, brokerOpts{sink: sink})

	client, sessionPresent := testutil.MQTTConnect(t, sslURL(b), dev.identity.CN(), false, dev.tlsCfg)
	if sessionPresent {
		t.Error("first connect reported session_present")
	}

	waitFor(t, 5*time.Second, "device row connected", func() bool {
		d, err := env.st.GetDevice(ctx, env.realm.ID, dev.identity.DeviceID)
		return err == nil && d.Connected
	})
	d, err := env.st.GetDevice(ctx, env.realm.ID, dev.identity.DeviceID)
	if err != nil {
		t.Fatalf("GetDevice: %v", err)
	}
	if d.LastConnection == nil {
		t.Error("last_connection not stamped")
	}
	if d.LastSeenIP == nil || !d.LastSeenIP.IsLoopback() {
		t.Errorf("last_seen_ip = %v, want loopback", d.LastSeenIP)
	}

	client.Disconnect(100)
	waitFor(t, 5*time.Second, "device row disconnected", func() bool {
		d, err := env.st.GetDevice(ctx, env.realm.ID, dev.identity.DeviceID)
		return err == nil && !d.Connected && d.LastDisconnection != nil
	})

	events := sink.snapshot()
	if len(events) != 2 {
		t.Fatalf("got %d lifecycle events, want 2: %+v", len(events), events)
	}
	if events[0].Type != EventDeviceConnected || events[0].Realm != env.realm.Name ||
		events[0].DeviceID != dev.identity.DeviceID || !events[0].RemoteIP.IsLoopback() {
		t.Errorf("connect event wrong: %+v", events[0])
	}
	if events[1].Type != EventDeviceDisconnected || events[1].DeviceID != dev.identity.DeviceID {
		t.Errorf("disconnect event wrong: %+v", events[1])
	}
}

func testRejects(t *testing.T, env *e2eEnv) {
	ctx := context.Background()
	dev := env.newDevice(t)
	b := env.newBroker(t, brokerOpts{enforceLatestCert: true})
	url := sslURL(b)

	t.Run("WrongRealmCA", func(t *testing.T) {
		key, csrPEM := testutil.DeviceCSR(t)
		certPEM, _, _, err := env.otherCA.SignCSR(csrPEM, env.realm.Name, dev.identity.DeviceID.String(), time.Hour)
		if err != nil {
			t.Fatalf("SignCSR: %v", err)
		}
		cfg := testutil.DeviceTLSConfig(t, certPEM, key, env.roots)
		if _, _, err := testutil.MQTTTryConnect(t, url, dev.identity.CN(), true, cfg); err == nil {
			t.Error("certificate from another realm's CA accepted")
		}
	})

	t.Run("UnknownDevice", func(t *testing.T) {
		ghostID, _ := deviceid.Random()
		key, csrPEM := testutil.DeviceCSR(t)
		certPEM, _, _, err := env.realmCA(t).SignCSR(csrPEM, env.realm.Name, ghostID.String(), time.Hour)
		if err != nil {
			t.Fatalf("SignCSR: %v", err)
		}
		cfg := testutil.DeviceTLSConfig(t, certPEM, key, env.roots)
		cn := env.realm.Name + "/" + ghostID.String()
		if _, _, err := testutil.MQTTTryConnect(t, url, cn, true, cfg); err == nil {
			t.Error("unregistered device accepted")
		}
	})

	t.Run("ClientIDMismatch", func(t *testing.T) {
		otherID, _ := deviceid.Random()
		cn := env.realm.Name + "/" + otherID.String()
		if _, _, err := testutil.MQTTTryConnect(t, url, cn, true, dev.tlsCfg); err == nil {
			t.Error("client ID different from certificate CN accepted")
		}
	})

	t.Run("Inhibited", func(t *testing.T) {
		if err := env.st.SetDeviceInhibited(ctx, env.realm.ID, dev.identity.DeviceID, true); err != nil {
			t.Fatalf("SetDeviceInhibited: %v", err)
		}
		defer func() {
			if err := env.st.SetDeviceInhibited(ctx, env.realm.ID, dev.identity.DeviceID, false); err != nil {
				t.Fatalf("clearing inhibition: %v", err)
			}
		}()
		if _, _, err := testutil.MQTTTryConnect(t, url, dev.identity.CN(), true, dev.tlsCfg); err == nil {
			t.Error("inhibited device accepted")
		}
	})

	t.Run("StaleSerial", func(t *testing.T) {
		// Issue a newer certificate: the DB now records its serial, so the
		// original (still time-valid) certificate must be rejected while
		// enforcement is on.
		key2, csrPEM2 := testutil.DeviceCSR(t)
		certPEM2, err := env.svc.Credentials(context.Background(), env.realm.Name,
			dev.identity.DeviceID.String(), dev.secret, csrPEM2, netip.MustParseAddr("127.0.0.1"))
		if err != nil {
			t.Fatalf("rotating credentials: %v", err)
		}
		if _, _, err := testutil.MQTTTryConnect(t, url, dev.identity.CN(), true, dev.tlsCfg); err == nil {
			t.Error("stale-serial certificate accepted with enforcement on")
		}
		freshCfg := testutil.DeviceTLSConfig(t, certPEM2, key2, env.roots)
		if _, _, err := testutil.MQTTTryConnect(t, url, dev.identity.CN(), true, freshCfg); err != nil {
			t.Errorf("latest certificate rejected: %v", err)
		}

		// With enforcement off, the superseded certificate still connects.
		lenient := env.newBroker(t, brokerOpts{enforceLatestCert: false})
		if _, _, err := testutil.MQTTTryConnect(t, sslURL(lenient), dev.identity.CN(), true, dev.tlsCfg); err != nil {
			t.Errorf("superseded certificate rejected with enforcement off: %v", err)
		}
	})

	t.Run("PlaintextWithDevModeOff", func(t *testing.T) {
		if b.DevAddr() != "" {
			t.Fatal("dev listener bound although insecure_dev_mode is off")
		}
		if _, _, err := testutil.MQTTTryConnect(t, "tcp://"+b.TLSAddr(), dev.identity.CN(), true, nil); err == nil {
			t.Error("plaintext connection to the TLS listener accepted")
		}
	})
}

func testDevModePlaintext(t *testing.T, env *e2eEnv) {
	ctx := context.Background()
	dev := env.newDevice(t)
	b := env.newBroker(t, brokerOpts{insecureDevMode: true})
	if b.DevAddr() == "" {
		t.Fatal("dev listener not bound with insecure_dev_mode on")
	}
	if _, _, err := testutil.MQTTTryConnect(t, "tcp://"+b.DevAddr(), dev.identity.CN(), true, nil); err != nil {
		t.Fatalf("plaintext dev-mode connect failed: %v", err)
	}
	waitFor(t, 5*time.Second, "device row connected via dev listener", func() bool {
		d, err := env.st.GetDevice(ctx, env.realm.ID, dev.identity.DeviceID)
		return err == nil && d.Connected
	})
}

func testDeferredAck(t *testing.T, env *e2eEnv) {
	dev := env.newDevice(t)
	intake := newRecorderIntake(false) // manual ack
	b := env.newBroker(t, brokerOpts{intake: intake})

	client, _ := testutil.MQTTConnect(t, sslURL(b), dev.identity.CN(), true, dev.tlsCfg)
	topic := dev.identity.BaseTopic() + "/" + e2eDeviceIface + "/value"

	token := client.Publish(topic, 1, false, []byte("measurement"))
	msg := intake.next(t, 5*time.Second)
	if msg.Topic != topic || msg.QoS != 1 || msg.Realm != env.realm.Name || msg.DeviceID != dev.identity.DeviceID {
		t.Fatalf("unexpected InboundMessage: %+v", msg)
	}
	if token.WaitTimeout(700 * time.Millisecond) {
		t.Fatal("PUBACK arrived without Ack() — backpressure wiring broken")
	}
	msg.Ack()
	testutil.WaitToken(t, token, 5*time.Second)
}

func testSessionPersistence(t *testing.T, env *e2eEnv) {
	dev := env.newDevice(t)
	path := filepath.Join(t.TempDir(), "sessions.db")
	b1 := env.newBroker(t, brokerOpts{path: path})

	subTopic := dev.identity.BaseTopic() + "/" + e2eServerIface + "/#"
	dataTopic := dev.identity.BaseTopic() + "/" + e2eServerIface + "/value"

	client, sessionPresent := testutil.MQTTConnect(t, sslURL(b1), dev.identity.CN(), false, dev.tlsCfg)
	if sessionPresent {
		t.Fatal("fresh session reported session_present")
	}
	testutil.WaitToken(t, client.Subscribe(subTopic, 1, nil), 5*time.Second)
	client.Disconnect(100)
	time.Sleep(200 * time.Millisecond) // let the broker observe the disconnect

	// Queue three QoS 1 messages for the offline session — within the same
	// wall-clock second, which is exactly the replay-order trap.
	for _, payload := range []string{"m1", "m2", "m3"} {
		if err := b1.Publisher().Publish(dataTopic, []byte(payload), 1, false, 0); err != nil {
			t.Fatalf("offline publish %s: %v", payload, err)
		}
	}
	time.Sleep(200 * time.Millisecond) // let the session store flush
	if err := b1.Close(); err != nil {
		t.Fatalf("closing first broker: %v", err)
	}

	// Same bbolt file, fresh broker process state.
	b2 := env.newBroker(t, brokerOpts{path: path})

	received := make(chan string, 8)
	collector := func(opts *paho.ClientOptions) {
		opts.SetDefaultPublishHandler(func(_ paho.Client, m paho.Message) {
			received <- string(m.Payload())
		})
	}
	_, sessionPresent = testutil.MQTTConnect(t, sslURL(b2), dev.identity.CN(), false, dev.tlsCfg, collector)
	if !sessionPresent {
		t.Fatal("session_present not set after broker restart")
	}

	var got []string
	timeout := time.After(10 * time.Second)
	for len(got) < 3 {
		select {
		case p := <-received:
			got = append(got, p)
		case <-timeout:
			t.Fatalf("offline queue replay incomplete: got %v", got)
		}
	}
	want := []string{"m1", "m2", "m3"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("replay out of order: got %v, want %v", got, want)
		}
	}
}

func testQoS2ExactlyOnce(t *testing.T, env *e2eEnv) {
	dev := env.newDevice(t)
	intake := newRecorderIntake(true)
	b := env.newBroker(t, brokerOpts{intake: intake})

	received := make(chan string, 8)
	collector := func(opts *paho.ClientOptions) {
		opts.SetDefaultPublishHandler(func(_ paho.Client, m paho.Message) {
			received <- string(m.Payload())
		})
	}
	client, _ := testutil.MQTTConnect(t, sslURL(b), dev.identity.CN(), true, dev.tlsCfg, collector)

	// Device → server: the full PUBLISH/PUBREC/PUBREL/PUBCOMP handshake
	// completes and the message reaches the intake exactly once.
	topic := dev.identity.BaseTopic() + "/" + e2eDeviceIface + "/value"
	before := intake.count()
	testutil.WaitToken(t, client.Publish(topic, 2, false, []byte("qos2-up")), 10*time.Second)
	time.Sleep(500 * time.Millisecond) // window for any duplicate submission
	if n := intake.count() - before; n != 1 {
		t.Fatalf("QoS2 publish reached the intake %d times, want exactly 1", n)
	}

	// Server → device, QoS 2 subscription: delivered exactly once.
	subTopic := dev.identity.BaseTopic() + "/" + e2eServerIface + "/#"
	dataTopic := dev.identity.BaseTopic() + "/" + e2eServerIface + "/value"
	testutil.WaitToken(t, client.Subscribe(subTopic, 2, nil), 5*time.Second)
	if err := b.Publisher().Publish(dataTopic, []byte("qos2-down"), 2, false, 0); err != nil {
		t.Fatalf("inline QoS2 publish: %v", err)
	}
	select {
	case p := <-received:
		if p != "qos2-down" {
			t.Fatalf("unexpected payload %q", p)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("QoS2 downlink not delivered")
	}
	select {
	case p := <-received:
		t.Fatalf("duplicate QoS2 delivery: %q", p)
	case <-time.After(700 * time.Millisecond):
	}
}

func testOfflineMessageExpiry(t *testing.T, env *e2eEnv) {
	dev := env.newDevice(t)
	path := filepath.Join(t.TempDir(), "sessions.db")
	b1 := env.newBroker(t, brokerOpts{path: path})

	subTopic := dev.identity.BaseTopic() + "/" + e2eServerIface + "/#"
	dataTopic := dev.identity.BaseTopic() + "/" + e2eServerIface + "/value"

	client, _ := testutil.MQTTConnect(t, sslURL(b1), dev.identity.CN(), false, dev.tlsCfg)
	testutil.WaitToken(t, client.Subscribe(subTopic, 1, nil), 5*time.Second)
	client.Disconnect(100)
	time.Sleep(200 * time.Millisecond)

	if err := b1.Publisher().Publish(dataTopic, []byte("dies"), 1, false, time.Second); err != nil {
		t.Fatalf("publishing expiring message: %v", err)
	}
	if err := b1.Publisher().Publish(dataTopic, []byte("keeper"), 1, false, 0); err != nil {
		t.Fatalf("publishing keeper message: %v", err)
	}

	// Wait past the expiry (the broker sweeps expired offline messages
	// every second), then restart on the same session file: the load path
	// must also refuse to resurrect it.
	time.Sleep(2500 * time.Millisecond)
	if err := b1.Close(); err != nil {
		t.Fatalf("closing first broker: %v", err)
	}
	b2 := env.newBroker(t, brokerOpts{path: path})

	received := make(chan string, 8)
	collector := func(opts *paho.ClientOptions) {
		opts.SetDefaultPublishHandler(func(_ paho.Client, m paho.Message) {
			received <- string(m.Payload())
		})
	}
	_, sessionPresent := testutil.MQTTConnect(t, sslURL(b2), dev.identity.CN(), false, dev.tlsCfg, collector)
	if !sessionPresent {
		t.Fatal("session_present not set after restart")
	}

	select {
	case p := <-received:
		if p != "keeper" {
			t.Fatalf("expired message delivered: %q", p)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("keeper message not delivered")
	}
	select {
	case p := <-received:
		t.Fatalf("unexpected second delivery: %q (expired message must not arrive)", p)
	case <-time.After(time.Second):
	}
}

func testSessionTakeover(t *testing.T, env *e2eEnv) {
	ctx := context.Background()
	dev := env.newDevice(t)
	sink := &recorderSink{}
	b := env.newBroker(t, brokerOpts{sink: sink})
	url := sslURL(b)

	c1, _ := testutil.MQTTConnect(t, url, dev.identity.CN(), false, dev.tlsCfg)
	waitFor(t, 5*time.Second, "first connection recorded", func() bool {
		d, err := env.st.GetDevice(ctx, env.realm.ID, dev.identity.DeviceID)
		return err == nil && d.Connected
	})

	c2, _ := testutil.MQTTConnect(t, url, dev.identity.CN(), false, dev.tlsCfg)
	waitFor(t, 5*time.Second, "old connection dropped", func() bool {
		return !c1.IsConnected()
	})

	// The takeover must not mark the device disconnected: the new channel
	// owns the session.
	time.Sleep(300 * time.Millisecond)
	d, err := env.st.GetDevice(ctx, env.realm.ID, dev.identity.DeviceID)
	if err != nil {
		t.Fatalf("GetDevice: %v", err)
	}
	if !d.Connected {
		t.Fatal("device marked disconnected after session takeover")
	}
	for _, e := range sink.snapshot() {
		if e.Type == EventDeviceDisconnected {
			t.Fatalf("spurious device_disconnected during takeover: %+v", e)
		}
	}

	c2.Disconnect(100)
	waitFor(t, 5*time.Second, "device row disconnected", func() bool {
		d, err := env.st.GetDevice(ctx, env.realm.ID, dev.identity.DeviceID)
		return err == nil && !d.Connected
	})
	events := sink.snapshot()
	if len(events) == 0 || events[len(events)-1].Type != EventDeviceDisconnected {
		t.Fatalf("missing trailing device_disconnected: %+v", events)
	}
}
