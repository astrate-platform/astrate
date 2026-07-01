package broker

import (
	"context"
	"fmt"
	"net/netip"
	"sync"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/astrate-platform/astrate/internal/pairing/ca"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
)

// Compile-time guard: the real *store.Store satisfies the broker's Store
// port, so M8 wiring is a plain assignment.
var _ Store = (*store.Store)(nil)

// --- shared test doubles (also used by the e2e suite) ---

// recorderIntake records submitted messages; with autoAck it acknowledges
// immediately, otherwise the test drives Ack by hand via the channel.
type recorderIntake struct {
	autoAck bool
	ch      chan InboundMessage

	mu   sync.Mutex
	msgs []InboundMessage
}

func newRecorderIntake(autoAck bool) *recorderIntake {
	return &recorderIntake{autoAck: autoAck, ch: make(chan InboundMessage, 128)}
}

func (r *recorderIntake) Submit(m InboundMessage) {
	r.mu.Lock()
	r.msgs = append(r.msgs, m)
	r.mu.Unlock()
	if r.autoAck {
		m.Ack()
	}
	select {
	case r.ch <- m:
	default:
	}
}

func (r *recorderIntake) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.msgs)
}

// next waits for the next submitted message.
func (r *recorderIntake) next(t testing.TB, timeout time.Duration) InboundMessage {
	t.Helper()
	select {
	case m := <-r.ch:
		return m
	case <-time.After(timeout):
		t.Fatalf("no message reached the intake within %s", timeout)
		return InboundMessage{}
	}
}

// recorderSink records lifecycle events.
type recorderSink struct {
	mu     sync.Mutex
	events []LifecycleEvent
}

func (r *recorderSink) OnLifecycleEvent(e LifecycleEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, e)
}

func (r *recorderSink) snapshot() []LifecycleEvent {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]LifecycleEvent, len(r.events))
	copy(out, r.events)
	return out
}

// waitFor polls cond until it holds or the timeout elapses.
func waitFor(t testing.TB, timeout time.Duration, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		if cond() {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("condition %q not reached within %s", what, timeout)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// fakeStore is an in-memory broker.Store for T1 tests.
type fakeStore struct {
	mu         sync.Mutex
	realms     []store.Realm
	devices    map[string]*store.Device
	interfaces map[string]*store.StoredInterface
	connects   []netip.Addr
	disconns   int
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		devices:    map[string]*store.Device{},
		interfaces: map[string]*store.StoredInterface{},
	}
}

func devKey(realmID int16, id deviceid.ID) string { return fmt.Sprintf("%d/%s", realmID, id) }

func (f *fakeStore) ListRealms(context.Context) ([]store.Realm, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]store.Realm(nil), f.realms...), nil
}

func (f *fakeStore) GetDevice(_ context.Context, realmID int16, id deviceid.ID) (*store.Device, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	d, ok := f.devices[devKey(realmID, id)]
	if !ok {
		return nil, store.ErrNotFound
	}
	cp := *d
	return &cp, nil
}

func (f *fakeStore) GetInterface(_ context.Context, realmID int16, name string, major int) (*store.StoredInterface, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	si, ok := f.interfaces[fmt.Sprintf("%d/%s/%d", realmID, name, major)]
	if !ok {
		return nil, store.ErrNotFound
	}
	return si, nil
}

func (f *fakeStore) SetDeviceConnected(_ context.Context, realmID int16, id deviceid.ID, _ time.Time, ip netip.Addr) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	d, ok := f.devices[devKey(realmID, id)]
	if !ok {
		return store.ErrNotFound
	}
	d.Connected = true
	f.connects = append(f.connects, ip)
	return nil
}

func (f *fakeStore) SetDeviceDisconnected(_ context.Context, realmID int16, id deviceid.ID, _ time.Time) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	d, ok := f.devices[devKey(realmID, id)]
	if !ok {
		return store.ErrNotFound
	}
	d.Connected = false
	f.disconns++
	return nil
}

func (f *fakeStore) connected(realmID int16, id deviceid.ID) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	d, ok := f.devices[devKey(realmID, id)]
	return ok && d.Connected
}

// --- T1 tests ---

func TestNewValidation(t *testing.T) {
	ctx := context.Background()
	st := newFakeStore()
	intake := newRecorderIntake(true)
	cert, _ := testutil.ServerTLSCert(t)

	if _, err := New(ctx, Config{SessionStorePath: "x", ServerTLSCert: cert}, st, nil, nil); err == nil {
		t.Error("New without intake: expected error")
	}
	if _, err := New(ctx, Config{ServerTLSCert: cert}, st, intake, nil); err == nil {
		t.Error("New without SessionStorePath: expected error")
	}
	if _, err := New(ctx, Config{SessionStorePath: "x"}, st, intake, nil); err == nil {
		t.Error("New without ServerTLSCert: expected error")
	}
}

// newFakeEnv assembles a complete in-memory environment: one realm with a
// live CA, one confirmed device (certificate trail stamped), and two
// introspected interfaces. It returns the store, the realm CA (for minting
// further certificates), the device identity, and the issued cert PEM.
func newFakeEnv(t *testing.T) (*fakeStore, *ca.CA, Identity, string) {
	t.Helper()
	realmCA, err := ca.Generate("test", 0)
	if err != nil {
		t.Fatalf("generating realm CA: %v", err)
	}
	devID, err := deviceid.Random()
	if err != nil {
		t.Fatalf("generating device ID: %v", err)
	}
	identity := Identity{Realm: "test", DeviceID: devID}

	st := newFakeStore()
	st.realms = []store.Realm{{ID: 1, Name: "test", CACertificatePEM: realmCA.CertificatePEM()}}
	st.devices[devKey(1, devID)] = &store.Device{
		ID: devID, RealmID: 1, Status: store.DeviceStatusConfirmed,
		Introspection: map[string]store.InterfaceVersion{
			"com.ex.DeviceData": {Major: 1},
			"com.ex.ServerData": {Major: 1},
		},
	}
	st.interfaces["1/com.ex.DeviceData/1"] = &store.StoredInterface{Name: "com.ex.DeviceData", Major: 1, Ownership: interfaceschema.OwnershipDevice}
	st.interfaces["1/com.ex.ServerData/1"] = &store.StoredInterface{Name: "com.ex.ServerData", Major: 1, Ownership: interfaceschema.OwnershipServer}

	_, csrPEM := testutil.DeviceCSR(t)
	certPEM, serial, aki, err := realmCA.SignCSR(csrPEM, "test", devID.String(), time.Hour)
	if err != nil {
		t.Fatalf("issuing device certificate: %v", err)
	}
	dev := st.devices[devKey(1, devID)]
	dev.CertSerial, dev.CertAKI = &serial, &aki
	return st, realmCA, identity, certPEM
}

func TestBrokerInMemorySmoke(t *testing.T) {
	ctx := context.Background()
	st, realmCA, identity, _ := newFakeEnv(t)
	intake := newRecorderIntake(false)
	sink := &recorderSink{}
	serverCert, roots := testutil.ServerTLSCert(t)

	b, err := New(ctx, Config{
		TLSAddr:          "127.0.0.1:0",
		ServerTLSCert:    serverCert,
		SessionStorePath: t.TempDir() + "/sessions.db",
		Logger:           discardLogger(),
	}, st, intake, sink)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := b.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(func() { _ = b.Close() })

	devKeyPriv, csrPEM := testutil.DeviceCSR(t)
	certPEM, _, _, err := realmCA.SignCSR(csrPEM, identity.Realm, identity.DeviceID.String(), time.Hour)
	if err != nil {
		t.Fatalf("issuing device certificate: %v", err)
	}
	tlsCfg := testutil.DeviceTLSConfig(t, certPEM, devKeyPriv, roots)
	url := "ssl://" + b.TLSAddr()

	client, sessionPresent := testutil.MQTTConnect(t, url, identity.CN(), true, tlsCfg)
	if sessionPresent {
		t.Error("first clean connect reported session_present")
	}
	waitFor(t, 5*time.Second, "device row connected", func() bool {
		return st.connected(1, identity.DeviceID)
	})

	base := identity.BaseTopic()

	t.Run("DeferredAck", func(t *testing.T) {
		token := client.Publish(base+"/com.ex.DeviceData/value", 1, false, []byte("hello"))
		msg := intake.next(t, 5*time.Second)
		if msg.Realm != "test" || msg.DeviceID != identity.DeviceID || msg.QoS != 1 ||
			msg.Topic != base+"/com.ex.DeviceData/value" || string(msg.Payload) != "hello" {
			t.Fatalf("unexpected InboundMessage: %+v", msg)
		}
		if token.WaitTimeout(500 * time.Millisecond) {
			t.Fatal("PUBACK arrived before Ack() — deferred ack broken")
		}
		msg.Ack()
		testutil.WaitToken(t, token, 5*time.Second)
		msg.Ack() // idempotent
	})

	t.Run("QoS0Flows", func(t *testing.T) {
		token := client.Publish(base, 0, false, []byte("com.ex.DeviceData:1:0"))
		testutil.WaitToken(t, token, 5*time.Second)
		msg := intake.next(t, 5*time.Second)
		if msg.Topic != base || msg.QoS != 0 {
			t.Fatalf("unexpected QoS0 message: %+v", msg)
		}
		msg.Ack() // no-op, must not panic
	})

	t.Run("NoEchoToPublisher", func(t *testing.T) {
		echoed := make(chan string, 8)
		token := client.Subscribe(base+"/#", 1, func(_ paho.Client, m paho.Message) {
			echoed <- m.Topic()
		})
		testutil.WaitToken(t, token, 5*time.Second)

		pubToken := client.Publish(base+"/com.ex.DeviceData/value", 1, false, []byte("no-echo"))
		intake.next(t, 5*time.Second).Ack()
		testutil.WaitToken(t, pubToken, 5*time.Second)

		select {
		case topic := <-echoed:
			t.Fatalf("device publish echoed back on %s", topic)
		case <-time.After(500 * time.Millisecond):
		}

		// Server-side publishes do reach the superset subscription.
		if err := b.Publisher().Publish(base+"/com.ex.ServerData/value", []byte("downlink"), 1, false, 0); err != nil {
			t.Fatalf("inline publish: %v", err)
		}
		select {
		case topic := <-echoed:
			if topic != base+"/com.ex.ServerData/value" {
				t.Fatalf("unexpected delivery on %s", topic)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("server-owned publish not delivered")
		}
	})

	t.Run("SubscribeDenied", func(t *testing.T) {
		token := client.Subscribe("other/AAAAAAAAAAAAAAAAAAAAAA/#", 1, nil)
		if !token.WaitTimeout(5 * time.Second) {
			t.Fatal("subscribe token timed out")
		}
		sub, ok := token.(*paho.SubscribeToken)
		if !ok {
			t.Fatalf("unexpected token type %T", token)
		}
		if granted := sub.Result()["other/AAAAAAAAAAAAAAAAAAAAAA/#"]; granted != 0x80 && token.Error() == nil {
			t.Fatalf("forbidden subscription granted qos %#x", granted)
		}
	})

	t.Run("LifecycleEvents", func(t *testing.T) {
		client.Disconnect(100)
		waitFor(t, 5*time.Second, "device row disconnected", func() bool {
			return !st.connected(1, identity.DeviceID)
		})
		events := sink.snapshot()
		if len(events) < 2 {
			t.Fatalf("got %d lifecycle events, want >= 2", len(events))
		}
		first, last := events[0], events[len(events)-1]
		if first.Type != EventDeviceConnected || first.Realm != "test" || first.DeviceID != identity.DeviceID {
			t.Errorf("first event %+v, want device_connected", first)
		}
		if !first.RemoteIP.IsLoopback() {
			t.Errorf("connect event IP = %s, want loopback", first.RemoteIP)
		}
		if last.Type != EventDeviceDisconnected || last.DeviceID != identity.DeviceID {
			t.Errorf("last event %+v, want device_disconnected", last)
		}
	})
}

// TestBrokerClientIDRemappedToCertCN: the device identity is the certificate
// CN alone — the wire client ID is free-form (upstream parity: the official
// Python SDK connects with a random paho-generated ID) and is rewritten to
// the CN before the session binds, so it cannot impersonate another device.
func TestBrokerClientIDRemappedToCertCN(t *testing.T) {
	ctx := context.Background()
	st, realmCA, identity, _ := newFakeEnv(t)
	intake := newRecorderIntake(true)
	serverCert, roots := testutil.ServerTLSCert(t)

	b, err := New(ctx, Config{
		TLSAddr:          "127.0.0.1:0",
		ServerTLSCert:    serverCert,
		SessionStorePath: t.TempDir() + "/sessions.db",
		Logger:           discardLogger(),
	}, st, intake, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := b.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(func() { _ = b.Close() })

	devKeyPriv, csrPEM := testutil.DeviceCSR(t)
	certPEM, _, _, err := realmCA.SignCSR(csrPEM, identity.Realm, identity.DeviceID.String(), time.Hour)
	if err != nil {
		t.Fatalf("issuing device certificate: %v", err)
	}
	tlsCfg := testutil.DeviceTLSConfig(t, certPEM, devKeyPriv, roots)
	url := "ssl://" + b.TLSAddr()

	publishAndCheckAttribution := func(t *testing.T, clientID string) {
		t.Helper()
		client, _ := testutil.MQTTConnect(t, url, clientID, true, tlsCfg)
		defer client.Disconnect(100)
		token := client.Publish(identity.BaseTopic()+"/com.ex.DeviceData/value", 1, false, []byte("v"))
		msg := intake.next(t, 5*time.Second)
		if msg.Realm != identity.Realm || msg.DeviceID != identity.DeviceID {
			t.Fatalf("publish as %q attributed to %s/%s, want the certificate identity %s",
				clientID, msg.Realm, msg.DeviceID, identity.CN())
		}
		testutil.WaitToken(t, token, 5*time.Second)
	}

	// A random (paho-style) client ID connects and acts as the certificate's
	// device.
	publishAndCheckAttribution(t, "paho-random-3f2a")
	// A client ID naming ANOTHER device is remapped all the same: no spoofing.
	publishAndCheckAttribution(t, "test/AAAAAAAAAAAAAAAAAAAAAA")
	// The CN as client ID (the Go SDK convention) still connects.
	publishAndCheckAttribution(t, identity.CN())
}

func TestReloadRealmsPicksUpNewRealm(t *testing.T) {
	ctx := context.Background()
	st, _, _, _ := newFakeEnv(t)
	intake := newRecorderIntake(true)
	serverCert, roots := testutil.ServerTLSCert(t)

	b, err := New(ctx, Config{
		TLSAddr:          "127.0.0.1:0",
		ServerTLSCert:    serverCert,
		SessionStorePath: t.TempDir() + "/sessions.db",
		Logger:           discardLogger(),
	}, st, intake, nil)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := b.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(func() { _ = b.Close() })

	// A realm created after the broker booted.
	newCA, err := ca.Generate("late", 0)
	if err != nil {
		t.Fatalf("generating CA: %v", err)
	}
	devID, _ := deviceid.Random()
	st.mu.Lock()
	st.realms = append(st.realms, store.Realm{ID: 2, Name: "late", CACertificatePEM: newCA.CertificatePEM()})
	st.devices[devKey(2, devID)] = &store.Device{ID: devID, RealmID: 2, Status: store.DeviceStatusConfirmed}
	st.mu.Unlock()

	devKeyPriv, csrPEM := testutil.DeviceCSR(t)
	certPEM, _, _, err := newCA.SignCSR(csrPEM, "late", devID.String(), time.Hour)
	if err != nil {
		t.Fatalf("issuing certificate: %v", err)
	}
	tlsCfg := testutil.DeviceTLSConfig(t, certPEM, devKeyPriv, roots)
	cn := "late/" + devID.String()
	url := "ssl://" + b.TLSAddr()

	// Before the reload the handshake fails (CA not in the union pool).
	if _, _, err := testutil.MQTTTryConnect(t, url, cn, true, tlsCfg); err == nil {
		t.Fatal("connect for unloaded realm unexpectedly succeeded")
	}
	if err := b.ReloadRealms(ctx); err != nil {
		t.Fatalf("ReloadRealms: %v", err)
	}
	if _, _, err := testutil.MQTTTryConnect(t, url, cn, true, tlsCfg); err != nil {
		t.Fatalf("connect after ReloadRealms failed: %v", err)
	}
}
