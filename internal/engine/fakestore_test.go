package engine

// Shared T1 test plumbing: an in-memory Store fake with programmable
// failures, delays, and gates, plus fixture helpers reusing the M1 interface
// definitions (docs/ROADMAP.md §7.3 "uses M1 fixtures"). Plain (untagged)
// file so the integration suite reuses the helpers, mirroring the broker
// recorder convention.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// fixtureDir is the M1 interface-definition fixture corpus.
const fixtureDir = "../../pkg/interfaceschema/testdata/valid"

// discardLogger silences engine logs in tests.
func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// waitFor polls cond until it holds or the timeout elapses.
func waitFor(t testing.TB, timeout time.Duration, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("timed out after %s waiting for %s", timeout, what)
}

// fixtureDefinition reads one M1 interface fixture by file name.
func fixtureDefinition(t testing.TB, name string) []byte {
	t.Helper()
	def, err := os.ReadFile(filepath.Join(fixtureDir, name))
	if err != nil {
		t.Fatalf("reading fixture %s: %v", name, err)
	}
	return def
}

// storedInterface builds a store.StoredInterface from a raw definition,
// assigning synthetic endpoint IDs (ifaceID*1000 + ordinal).
func storedInterface(t testing.TB, realmID int16, ifaceID int64, def []byte) *store.StoredInterface {
	t.Helper()
	iface, err := interfaceschema.ParseInterface(def)
	if err != nil {
		t.Fatalf("parsing fixture definition: %v", err)
	}
	endpoints := make(map[string]int64, len(iface.Mappings))
	for i := range iface.Mappings {
		endpoints[iface.Mappings[i].Endpoint] = ifaceID*1000 + int64(i)
	}
	return &store.StoredInterface{
		ID:          ifaceID,
		RealmID:     realmID,
		Name:        iface.Name,
		Major:       iface.Major,
		Minor:       iface.Minor,
		Type:        iface.Type,
		Ownership:   iface.Ownership,
		Aggregation: iface.Aggregation,
		Definition:  def,
		Endpoints:   endpoints,
	}
}

// fixtureStored loads and wraps one M1 fixture in a single call.
func fixtureStored(t testing.TB, realmID int16, ifaceID int64, name string) *store.StoredInterface {
	t.Helper()
	return storedInterface(t, realmID, ifaceID, fixtureDefinition(t, name))
}

// unsetCall records one UnsetProperty invocation.
type unsetCall struct {
	realmID     int16
	deviceID    deviceid.ID
	interfaceID int64
	path        string
}

// fakeStore is the in-memory Store implementation for T1 tests.
type fakeStore struct {
	mu sync.Mutex

	realms     []store.Realm
	interfaces map[int16][]*store.StoredInterface
	devices    map[deviceKey]*store.Device

	// listRealmsCalls / getDeviceCalls count store round-trips.
	listRealmsCalls int
	getDeviceCalls  int

	// getDeviceErrs holds errors returned by the next GetDevice calls
	// (consumed front to back).
	getDeviceErrs []error

	// appendGate, when non-nil, blocks AppendDatastreams until the channel
	// is closed (backpressure tests).
	appendGate chan struct{}
	// appendDelay sleeps inside AppendDatastreams (synthetic slow store).
	appendDelay time.Duration
	// appendErrs holds errors returned by the next AppendDatastreams calls.
	appendErrs []error
	// appendPanics makes the next AppendDatastreams calls panic.
	appendPanics int

	// inFlight / maxInFlight observe AppendDatastreams concurrency.
	inFlight    int
	maxInFlight int

	batches []store.DatastreamBatch
	upserts []store.Property
	unsets  []unsetCall
	hints   map[deviceKey]string

	notifyCh chan store.Notification
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		interfaces: make(map[int16][]*store.StoredInterface),
		devices:    make(map[deviceKey]*store.Device),
		hints:      make(map[deviceKey]string),
		notifyCh:   make(chan store.Notification, 16),
	}
}

// addRealm registers a realm.
func (f *fakeStore) addRealm(id int16, name string) {
	f.mu.Lock()
	f.realms = append(f.realms, store.Realm{ID: id, Name: name})
	f.mu.Unlock()
}

// addInterface installs an interface into a realm.
func (f *fakeStore) addInterface(realmID int16, si *store.StoredInterface) {
	f.mu.Lock()
	f.interfaces[realmID] = append(f.interfaces[realmID], si)
	f.mu.Unlock()
}

// addDevice registers a device row.
func (f *fakeStore) addDevice(realm string, realmID int16, id deviceid.ID, intro map[string]store.InterfaceVersion, hint string) {
	f.mu.Lock()
	f.devices[deviceKey{realm: realm, id: id}] = &store.Device{
		ID: id, RealmID: realmID, Introspection: intro, PayloadFormatHint: hint,
	}
	f.mu.Unlock()
}

func (f *fakeStore) ListRealms(_ context.Context) ([]store.Realm, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.listRealmsCalls++
	out := make([]store.Realm, len(f.realms))
	copy(out, f.realms)
	return out, nil
}

func (f *fakeStore) LoadRealmInterfaces(_ context.Context, realmID int16) ([]*store.StoredInterface, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]*store.StoredInterface, len(f.interfaces[realmID]))
	copy(out, f.interfaces[realmID])
	return out, nil
}

func (f *fakeStore) Listen(ctx context.Context, channel string) (<-chan store.Notification, error) {
	if channel != store.ChannelInterfaces {
		return nil, fmt.Errorf("fakeStore: unexpected channel %q", channel)
	}
	// Mirror store.Listen: the returned channel closes when ctx is
	// cancelled, so the engine's invalidation goroutine winds down.
	out := make(chan store.Notification, 16)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case n, ok := <-f.notifyCh:
				if !ok {
					return
				}
				select {
				case out <- n:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out, nil
}

func (f *fakeStore) GetDevice(_ context.Context, realmID int16, id deviceid.ID) (*store.Device, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.getDeviceCalls++
	if len(f.getDeviceErrs) > 0 {
		err := f.getDeviceErrs[0]
		f.getDeviceErrs = f.getDeviceErrs[1:]
		return nil, err
	}
	for key, dev := range f.devices {
		if key.id == id && dev.RealmID == realmID {
			cp := *dev
			return &cp, nil
		}
	}
	return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, id)
}

func (f *fakeStore) SetPayloadFormatHint(_ context.Context, realmID int16, id deviceid.ID, hint string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for key, dev := range f.devices {
		if key.id == id && dev.RealmID == realmID {
			dev.PayloadFormatHint = hint
			f.hints[key] = hint
			return nil
		}
	}
	return fmt.Errorf("%w: device %s", store.ErrNotFound, id)
}

func (f *fakeStore) AppendDatastreams(_ context.Context, batch store.DatastreamBatch) error {
	f.mu.Lock()
	gate := f.appendGate
	delay := f.appendDelay
	panics := f.appendPanics > 0
	if panics {
		f.appendPanics--
	}
	var err error
	if !panics && len(f.appendErrs) > 0 {
		err = f.appendErrs[0]
		f.appendErrs = f.appendErrs[1:]
	}
	f.inFlight++
	if f.inFlight > f.maxInFlight {
		f.maxInFlight = f.inFlight
	}
	f.mu.Unlock()

	defer func() {
		f.mu.Lock()
		f.inFlight--
		f.mu.Unlock()
	}()

	if gate != nil {
		<-gate
	}
	if delay > 0 {
		time.Sleep(delay)
	}
	if panics {
		panic("fakeStore: injected AppendDatastreams panic")
	}
	if err != nil {
		return err
	}

	f.mu.Lock()
	f.batches = append(f.batches, batch)
	f.mu.Unlock()
	return nil
}

func (f *fakeStore) UpsertProperty(_ context.Context, p store.Property) error {
	f.mu.Lock()
	f.upserts = append(f.upserts, p)
	f.mu.Unlock()
	return nil
}

func (f *fakeStore) UnsetProperty(_ context.Context, realmID int16, deviceID deviceid.ID, interfaceID int64, path string) (bool, error) {
	f.mu.Lock()
	f.unsets = append(f.unsets, unsetCall{realmID, deviceID, interfaceID, path})
	f.mu.Unlock()
	return true, nil
}

// individualRows flattens every committed individual row in commit order.
func (f *fakeStore) individualRows() []store.IndividualRow {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []store.IndividualRow
	for _, b := range f.batches {
		out = append(out, b.Individual...)
	}
	return out
}

// objectRows flattens every committed object row in commit order.
func (f *fakeStore) objectRows() []store.ObjectRow {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []store.ObjectRow
	for _, b := range f.batches {
		out = append(out, b.Objects...)
	}
	return out
}

// batchCount returns the number of committed batches.
func (f *fakeStore) batchCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.batches)
}

// inFlightNow returns the number of AppendDatastreams calls in progress.
func (f *fakeStore) inFlightNow() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.inFlight
}

// maxInFlightNow returns the high-water concurrent AppendDatastreams count.
func (f *fakeStore) maxInFlightNow() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.maxInFlight
}

// hintFor returns the last persisted payload-format hint for the device.
func (f *fakeStore) hintFor(realm string, id deviceid.ID) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.hints[deviceKey{realm: realm, id: id}]
}

// errTransient is the injected transient store failure.
var errTransient = errors.New("fakeStore: injected transient failure")

// ---------------------------------------------------------------------------
// Shared pipeline test rig: a seeded realm with the M1 fixture interfaces,
// message builders, and engine construction helpers.
// ---------------------------------------------------------------------------

// Seeded realm constants.
const (
	realmAlpha   = "alpha"
	realmAlphaID = int16(1)
)

// devAlpha is the seeded test device.
var devAlpha = deviceid.ID{0xa1, 0x70, 4, 8, 15, 16, 23, 42, 1, 2, 3, 4, 5, 6, 7, 8}

// Fixture interface IDs assigned by seedAlpha.
const (
	ifaceAllScalars  = int64(10)
	ifaceObjectFlat  = int64(11)
	ifaceGeolocation = int64(12)
	ifaceServerProps = int64(13)
	ifacePropArrays  = int64(14)
	ifaceMinimal     = int64(15)
)

// seedAlpha installs the fixture corpus into realm alpha and registers
// devAlpha with an introspection declaring all of it (plus one interface
// that is deliberately not installed).
func seedAlpha(t testing.TB, fs *fakeStore) {
	t.Helper()
	fs.addRealm(realmAlphaID, realmAlpha)
	fs.addInterface(realmAlphaID, fixtureStored(t, realmAlphaID, ifaceAllScalars, "com.astrate.test.AllScalarTypes.json"))
	fs.addInterface(realmAlphaID, fixtureStored(t, realmAlphaID, ifaceObjectFlat, "com.astrate.test.ObjectFlat.json"))
	fs.addInterface(realmAlphaID, fixtureStored(t, realmAlphaID, ifaceGeolocation, "org.astarte-platform.genericsensors.Geolocation.json"))
	fs.addInterface(realmAlphaID, fixtureStored(t, realmAlphaID, ifaceServerProps, "com.astrate.test.ServerProperties.json"))
	fs.addInterface(realmAlphaID, fixtureStored(t, realmAlphaID, ifacePropArrays, "com.astrate.test.PropertyArrays.json"))
	fs.addInterface(realmAlphaID, fixtureStored(t, realmAlphaID, ifaceMinimal, "com.astrate.test.Minimal.json"))
	fs.addDevice(realmAlpha, realmAlphaID, devAlpha, map[string]store.InterfaceVersion{
		"com.astrate.test.AllScalarTypes":                 {Major: 1, Minor: 0},
		"com.astrate.test.ObjectFlat":                     {Major: 1, Minor: 0},
		"org.astarte-platform.genericsensors.Geolocation": {Major: 1, Minor: 0},
		"com.astrate.test.ServerProperties":               {Major: 1, Minor: 2},
		"com.astrate.test.PropertyArrays":                 {Major: 2, Minor: 1},
		"com.astrate.test.Minimal":                        {Major: 0, Minor: 1},
		"com.astrate.test.DeclaredButMissing":             {Major: 1, Minor: 0},
	}, "")
}

// ackCounter records acknowledgment calls.
type ackCounter struct{ n atomic.Int32 }

// fn returns the Ack callback.
func (a *ackCounter) fn() func() { return func() { a.n.Add(1) } }

// acked reports whether Ack ran at least once.
func (a *ackCounter) acked() bool { return a.n.Load() > 0 }

// deviceMsg builds an InboundMessage from devAlpha. iface "" addresses the
// introspection topic; iface "control" plus a path addresses the control
// channel.
func deviceMsg(iface, path string, qos byte, body []byte, ack *ackCounter) broker.InboundMessage {
	return deviceMsgFor(devAlpha, iface, path, qos, body, ack)
}

// deviceMsgFor builds an InboundMessage from an arbitrary realm-alpha device.
func deviceMsgFor(dev deviceid.ID, iface, path string, qos byte, body []byte, ack *ackCounter) broker.InboundMessage {
	topic := realmAlpha + "/" + dev.String()
	if iface != "" {
		topic += "/" + iface + path
	}
	return broker.InboundMessage{
		Realm:      realmAlpha,
		DeviceID:   dev,
		Topic:      topic,
		Payload:    body,
		QoS:        qos,
		ReceivedAt: time.Now().UTC().Truncate(time.Millisecond),
		Ack:        ack.fn(),
	}
}

// enc encodes a `{v, t}` payload, failing the test on encoder errors.
func enc(t testing.TB, v payload.Value, ts *time.Time, f payload.Format) []byte {
	t.Helper()
	p, err := payload.Encode(v, ts, f)
	if err != nil {
		t.Fatalf("encoding payload: %v", err)
	}
	return p
}

// newTestEngine builds an unstarted engine over fs with quiet logs and the
// schema snapshot preloaded — data_test and batch_test drive handle/flush
// synchronously on shard 0.
func newTestEngine(t testing.TB, fs *fakeStore, cfg Config) *Engine {
	t.Helper()
	cfg.Logger = discardLogger()
	e, err := newEngine(cfg, fs)
	if err != nil {
		t.Fatalf("newEngine: %v", err)
	}
	if err := e.schemas.loadAll(context.Background()); err != nil {
		t.Fatalf("loadAll: %v", err)
	}
	return e
}

// startTestEngine builds and starts an engine; cleanup drains it and then
// cancels its run context.
func startTestEngine(t testing.TB, fs *fakeStore, cfg Config) *Engine {
	t.Helper()
	cfg.Logger = discardLogger()
	e, err := newEngine(cfg, fs)
	if err != nil {
		t.Fatalf("newEngine: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	if err := e.start(ctx); err != nil {
		cancel()
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		dctx, dcancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer dcancel()
		if err := e.drain(dctx); err != nil {
			t.Errorf("drain: %v", err)
		}
		cancel()
	})
	return e
}

// Compile-time guards: the fake satisfies the port, and so does the real
// store (M8 wiring is a plain assignment).
var (
	_ Store = (*fakeStore)(nil)
	_ Store = (*store.Store)(nil)
)
