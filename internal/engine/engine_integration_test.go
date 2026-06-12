//go:build integration

package engine

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/internal/testutil"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// errInjectedCommit is the injected datastream-commit failure.
var errInjectedCommit = errors.New("instrumented store: injected commit failure")

// instrumentedStore wraps the real store: it counts AppendDatastreams calls
// and fails the next failAppends of them before any data reaches the
// database — the docs/ROADMAP.md §7.3 T2 "injected commit failure" knob.
type instrumentedStore struct {
	*store.Store

	mu          sync.Mutex
	appendCalls int
	failAppends int
}

func (s *instrumentedStore) AppendDatastreams(ctx context.Context, batch store.DatastreamBatch) error {
	s.mu.Lock()
	s.appendCalls++
	fail := s.failAppends > 0
	if fail {
		s.failAppends--
	}
	s.mu.Unlock()
	if fail {
		return errInjectedCommit
	}
	return s.Store.AppendDatastreams(ctx, batch)
}

func (s *instrumentedStore) calls() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.appendCalls
}

// engineRig is the shared T2 fixture: migrated store, seeded realm/device,
// running engine.
type engineRig struct {
	st      *instrumentedStore
	raw     *store.Store
	e       *Engine
	realm   string
	realmID int16
	dev     deviceid.ID
	ifaces  map[string]*store.StoredInterface
}

// msg builds an InboundMessage for the rig's device.
func (r *engineRig) msg(iface, path string, qos byte, body []byte, ack *ackCounter) broker.InboundMessage {
	topic := r.realm + "/" + r.dev.String()
	if iface != "" {
		topic += "/" + iface + path
	}
	return broker.InboundMessage{
		Realm:      r.realm,
		DeviceID:   r.dev,
		Topic:      topic,
		Payload:    body,
		QoS:        qos,
		ReceivedAt: time.Now().UTC().Truncate(time.Millisecond),
		Ack:        ack.fn(),
	}
}

// TestEngine is the umbrella T2 suite (docs/ROADMAP.md §7.3): one
// TimescaleDB container (or ASTRATE_TEST_DSN), one migrated store, one
// running engine; sub-suites share the seeded realm. Realm names are
// per-run-unique so reruns against a reused database stay green.
func TestEngine(t *testing.T) {
	pool := testutil.StartTimescale(t)
	ctx := context.Background()

	raw, err := store.New(ctx, pool.Config().ConnString())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	t.Cleanup(raw.Close)

	var suffix [4]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		t.Fatalf("random suffix: %v", err)
	}
	realmName := "eng" + hex.EncodeToString(suffix[:])

	realm, err := raw.CreateRealm(ctx, store.NewRealm{
		Name:               realmName,
		JWTPublicKeysPEM:   []string{"test-key"},
		CACertificatePEM:   "test-ca-cert",
		CAPrivateKeySealed: []byte("test-sealed-key"),
	})
	if err != nil {
		t.Fatalf("CreateRealm: %v", err)
	}

	rig := &engineRig{
		st:      &instrumentedStore{Store: raw},
		raw:     raw,
		realm:   realmName,
		realmID: realm.ID,
		ifaces:  make(map[string]*store.StoredInterface),
	}

	// Pre-installed interfaces; ObjectFlat and DeepParametric arrive
	// mid-run through the two invalidation paths.
	for _, f := range []string{
		"com.astrate.test.AllScalarTypes.json",
		"org.astarte-platform.genericsensors.Geolocation.json",
		"com.astrate.test.PropertyArrays.json",
		"com.astrate.test.Minimal.json",
	} {
		si, err := raw.InstallInterface(ctx, realm.ID, fixtureDefinition(t, f))
		if err != nil {
			t.Fatalf("InstallInterface(%s): %v", f, err)
		}
		rig.ifaces[si.Name] = si
	}

	rig.dev, err = deviceid.Random()
	if err != nil {
		t.Fatalf("deviceid.Random: %v", err)
	}
	if err := raw.RegisterDevice(ctx, realm.ID, rig.dev, "secret-hash"); err != nil {
		t.Fatalf("RegisterDevice: %v", err)
	}
	intro := map[string]store.InterfaceVersion{
		"com.astrate.test.AllScalarTypes":                 {Major: 1, Minor: 0},
		"org.astarte-platform.genericsensors.Geolocation": {Major: 1, Minor: 0},
		"com.astrate.test.PropertyArrays":                 {Major: 2, Minor: 1},
		"com.astrate.test.Minimal":                        {Major: 0, Minor: 1},
		"com.astrate.test.ObjectFlat":                     {Major: 1, Minor: 0},
		"com.astrate.test.DeepParametric":                 {Major: 0, Minor: 1},
	}
	if _, err := raw.UpdateIntrospection(ctx, realm.ID, rig.dev, intro); err != nil {
		t.Fatalf("UpdateIntrospection: %v", err)
	}

	e, err := newEngine(Config{
		Shards: 4, BatchMaxRows: 8, BatchMaxWait: 20 * time.Millisecond,
		Logger: discardLogger(),
	}, rig.st)
	if err != nil {
		t.Fatalf("newEngine: %v", err)
	}
	runCtx, cancel := context.WithCancel(ctx)
	if err := e.start(runCtx); err != nil {
		cancel()
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		dctx, dcancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer dcancel()
		if err := e.drain(dctx); err != nil {
			t.Errorf("drain: %v", err)
		}
		cancel()
	})
	rig.e = e

	t.Run("TypedColumns", func(t *testing.T) { testTypedColumns(t, rig) })
	t.Run("ObjectRows", func(t *testing.T) { testObjectRows(t, rig) })
	t.Run("Properties", func(t *testing.T) { testProperties(t, rig) })
	t.Run("AckAfterCommit", func(t *testing.T) { testAckAfterCommit(t, rig) })
	t.Run("InvalidationListen", func(t *testing.T) { testInvalidationListen(t, rig) })
	t.Run("InvalidationCallback", func(t *testing.T) { testInvalidationCallback(t, rig) })
}

// testTypedColumns: one publish per scalar type (BSON and JSON) must land in
// exactly its typed column with the right value (docs/ROADMAP.md §7.3 T2).
func testTypedColumns(t *testing.T, rig *engineRig) {
	ctx := context.Background()
	si := rig.ifaces["com.astrate.test.AllScalarTypes"]
	explicit := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	blob := []byte{0xde, 0xad, 0xbe, 0xef}

	ack := &ackCounter{}
	pubs := []struct {
		path string
		body []byte
	}{
		{"/double", enc(t, 22.5, nil, payload.FormatBSON)},
		{"/integer", enc(t, int32(-7), nil, payload.FormatBSON)},
		{"/longinteger", enc(t, int64(1)<<60, nil, payload.FormatBSON)},
		{"/boolean", enc(t, true, nil, payload.FormatBSON)},
		{"/string", []byte(`{"v":"from json"}`)}, // JSON profile on the same pipeline
		{"/binaryblob", enc(t, blob, nil, payload.FormatBSON)},
		{"/datetime", enc(t, explicit, &explicit, payload.FormatBSON)},
	}
	for _, p := range pubs {
		rig.e.Submit(rig.msg("com.astrate.test.AllScalarTypes", p.path, 1, p.body, ack))
	}
	waitFor(t, 10*time.Second, "all scalar publishes acked", func() bool {
		return ack.n.Load() == int32(len(pubs))
	})

	series := func(path string) store.IndividualRow {
		t.Helper()
		rows, err := rig.raw.Series(ctx, store.SeriesQuery{
			RealmID: rig.realmID, DeviceID: rig.dev, InterfaceID: si.ID, Path: path,
		})
		if err != nil {
			t.Fatalf("Series(%s): %v", path, err)
		}
		if len(rows) != 1 {
			t.Fatalf("Series(%s) returned %d rows, want 1", path, len(rows))
		}
		if rows[0].EndpointID != si.Endpoints[path] {
			t.Errorf("Series(%s) endpoint %d, want %d", path, rows[0].EndpointID, si.Endpoints[path])
		}
		return rows[0]
	}

	if r := series("/double"); r.ValueDouble == nil || *r.ValueDouble != 22.5 {
		t.Errorf("/double: %+v", r.ValueDouble)
	}
	if r := series("/integer"); r.ValueInteger == nil || *r.ValueInteger != -7 {
		t.Errorf("/integer: %+v", r.ValueInteger)
	}
	if r := series("/longinteger"); r.ValueLonginteger == nil || *r.ValueLonginteger != 1<<60 {
		t.Errorf("/longinteger: %+v", r.ValueLonginteger)
	}
	if r := series("/boolean"); r.ValueBoolean == nil || !*r.ValueBoolean {
		t.Errorf("/boolean: %+v", r.ValueBoolean)
	}
	if r := series("/string"); r.ValueString == nil || *r.ValueString != "from json" {
		t.Errorf("/string: %+v", r.ValueString)
	}
	if r := series("/binaryblob"); string(r.ValueBinaryblob) != string(blob) {
		t.Errorf("/binaryblob: %x", r.ValueBinaryblob)
	}
	r := series("/datetime")
	if r.ValueDatetime == nil || !r.ValueDatetime.Equal(explicit) {
		t.Errorf("/datetime value: %+v", r.ValueDatetime)
	}
	if !r.TS.Equal(explicit) {
		t.Errorf("/datetime ts %s, want explicit %s", r.TS, explicit)
	}

	// The JSON publish flipped the sticky payload-format hint (§3.5.4).
	waitFor(t, 5*time.Second, "format hint flip persisted", func() bool {
		dev, err := rig.raw.GetDevice(ctx, rig.realmID, rig.dev)
		return err == nil && dev.PayloadFormatHint == hintJSON
	})
}

// testObjectRows: object-aggregated publishes land as one jsonb document on
// the prefix path.
func testObjectRows(t *testing.T, rig *engineRig) {
	ctx := context.Background()
	si := rig.ifaces["org.astarte-platform.genericsensors.Geolocation"]
	explicit := time.Date(2026, 6, 1, 12, 30, 0, 0, time.UTC)

	ack := &ackCounter{}
	body := enc(t, map[string]payload.Value{"latitude": 45.07, "longitude": 7.69}, &explicit, payload.FormatBSON)
	rig.e.Submit(rig.msg("org.astarte-platform.genericsensors.Geolocation", "/gps", 2, body, ack))
	waitFor(t, 10*time.Second, "object publish acked", func() bool { return ack.acked() })

	rows, err := rig.raw.ObjectSeries(ctx, store.SeriesQuery{
		RealmID: rig.realmID, DeviceID: rig.dev, InterfaceID: si.ID, Path: "/gps",
	})
	if err != nil {
		t.Fatalf("ObjectSeries: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("%d object rows, want 1", len(rows))
	}
	if string(rows[0].Value) != `{"latitude": 45.07, "longitude": 7.69}` &&
		string(rows[0].Value) != `{"latitude":45.07,"longitude":7.69}` {
		t.Errorf("object document: %s", rows[0].Value)
	}
	if !rows[0].TS.Equal(explicit) {
		t.Errorf("object ts %s, want %s", rows[0].TS, explicit)
	}
}

// testProperties: device property set is an upsert, empty payload deletes
// (docs/ROADMAP.md §7.3 "property upsert/unset").
func testProperties(t *testing.T, rig *engineRig) {
	ctx := context.Background()
	si := rig.ifaces["com.astrate.test.PropertyArrays"]
	const path = "/config/thresholds"

	set := func(body []byte) {
		t.Helper()
		ack := &ackCounter{}
		rig.e.Submit(rig.msg("com.astrate.test.PropertyArrays", path, 2, body, ack))
		waitFor(t, 10*time.Second, "property publish acked", func() bool { return ack.acked() })
	}

	set([]byte(`{"v":[1.5,2.5]}`))
	p, err := rig.raw.GetProperty(ctx, rig.realmID, rig.dev, si.ID, path)
	if err != nil {
		t.Fatalf("GetProperty after set: %v", err)
	}
	if string(p.Value) != `[1.5, 2.5]` && string(p.Value) != `[1.5,2.5]` {
		t.Errorf("property value: %s", p.Value)
	}

	set([]byte(`{"v":[9.0]}`)) // upsert, not insert
	p, err = rig.raw.GetProperty(ctx, rig.realmID, rig.dev, si.ID, path)
	if err != nil {
		t.Fatalf("GetProperty after upsert: %v", err)
	}
	if string(p.Value) != `[9]` {
		t.Errorf("property value after upsert: %s", p.Value)
	}

	set(nil) // unset
	if _, err := rig.raw.GetProperty(ctx, rig.realmID, rig.dev, si.ID, path); !errors.Is(err, store.ErrNotFound) {
		t.Errorf("property after unset: err=%v, want ErrNotFound", err)
	}
}

// testAckAfterCommit: with an injected commit failure the acknowledgment is
// withheld; the retry path eventually commits and only then acks
// (docs/DESIGN.md §5.3, docs/ROADMAP.md §7.3 T2).
func testAckAfterCommit(t *testing.T, rig *engineRig) {
	ctx := context.Background()
	rig.st.mu.Lock()
	rig.st.failAppends = 2
	callsBefore := rig.st.appendCalls
	rig.st.mu.Unlock()

	ack := &ackCounter{}
	rig.e.Submit(rig.msg("com.astrate.test.Minimal", "/value", 1, enc(t, 3.25, nil, payload.FormatBSON), ack))

	// Both injected failures must be consumed without an ack.
	waitFor(t, 10*time.Second, "injected commit failures consumed", func() bool {
		return rig.st.calls() >= callsBefore+2
	})
	if ack.acked() {
		t.Fatal("message acknowledged although no commit succeeded")
	}

	// The parked batch retries (100ms, then 200ms backoff) and commits.
	waitFor(t, 10*time.Second, "retry committed and acked", func() bool { return ack.acked() })
	rows, err := rig.raw.Series(ctx, store.SeriesQuery{
		RealmID: rig.realmID, DeviceID: rig.dev,
		InterfaceID: rig.ifaces["com.astrate.test.Minimal"].ID, Path: "/value",
	})
	if err != nil {
		t.Fatalf("Series: %v", err)
	}
	if len(rows) != 1 || *rows[0].ValueDouble != 3.25 {
		t.Fatalf("retry path persisted %d rows (%+v), want exactly 1", len(rows), rows)
	}
	if ack.n.Load() != 1 {
		t.Errorf("ack fired %d times, want 1", ack.n.Load())
	}
}

// testInvalidationListen: an interface installed mid-run becomes publishable
// without a restart via LISTEN/NOTIFY (docs/ROADMAP.md §7.3).
func testInvalidationListen(t *testing.T, rig *engineRig) {
	ctx := context.Background()

	// Declared in the introspection but not yet installed: rejected.
	ack := &ackCounter{}
	body := enc(t, map[string]payload.Value{"latitude": 1.0, "longitude": 2.0},
		ptrTime(time.Now().UTC()), payload.FormatBSON)
	rig.e.Submit(rig.msg("com.astrate.test.ObjectFlat", "", 1, body, ack))
	waitFor(t, 10*time.Second, "pre-install publish consumed", func() bool { return ack.acked() })

	si, err := rig.raw.InstallInterface(ctx, rig.realmID, fixtureDefinition(t, "com.astrate.test.ObjectFlat.json"))
	if err != nil {
		t.Fatalf("InstallInterface: %v", err)
	}
	rig.ifaces[si.Name] = si
	if err := rig.raw.NotifyInterfacesChanged(ctx, rig.realmID); err != nil {
		t.Fatalf("NotifyInterfacesChanged: %v", err)
	}

	// The notification is asynchronous: poll by re-publishing until a row
	// lands (each pre-reload attempt is consumed as a reject).
	waitFor(t, 15*time.Second, "mid-run interface becomes publishable", func() bool {
		a := &ackCounter{}
		rig.e.Submit(rig.msg("com.astrate.test.ObjectFlat", "", 1, body, a))
		rows, err := rig.raw.ObjectSeries(ctx, store.SeriesQuery{
			RealmID: rig.realmID, DeviceID: rig.dev, InterfaceID: si.ID, Path: "",
		})
		return err == nil && len(rows) > 0
	})
}

// testInvalidationCallback: the in-process RefreshInterfaces callback makes
// a new interface publishable synchronously.
func testInvalidationCallback(t *testing.T, rig *engineRig) {
	ctx := context.Background()

	si, err := rig.raw.InstallInterface(ctx, rig.realmID, fixtureDefinition(t, "com.astrate.test.DeepParametric.json"))
	if err != nil {
		t.Fatalf("InstallInterface: %v", err)
	}
	rig.ifaces[si.Name] = si
	if err := rig.e.RefreshInterfaces(ctx, rig.realmID); err != nil {
		t.Fatalf("RefreshInterfaces: %v", err)
	}

	ack := &ackCounter{}
	rig.e.Submit(rig.msg("com.astrate.test.DeepParametric", "/a/x/c/y/e", 1,
		enc(t, int64(42), nil, payload.FormatBSON), ack))
	waitFor(t, 10*time.Second, "publish on callback-refreshed interface", func() bool { return ack.acked() })

	rows, err := rig.raw.Series(ctx, store.SeriesQuery{
		RealmID: rig.realmID, DeviceID: rig.dev, InterfaceID: si.ID, Path: "/a/x/c/y/e",
	})
	if err != nil {
		t.Fatalf("Series: %v", err)
	}
	if len(rows) != 1 || rows[0].ValueLonginteger == nil || *rows[0].ValueLonginteger != 42 {
		t.Fatalf("parametric longinteger row: %+v", rows)
	}
}

// ptrTime returns &t.
func ptrTime(t time.Time) *time.Time { return &t }
