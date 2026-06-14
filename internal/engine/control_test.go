package engine

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"

	promtest "github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// goldenFrame loads the frozen producer/properties blob: framed with an
// independent zlib implementation (CPython), exactly like an official SDK
// would produce it (docs/ROADMAP.md §7.3 "zlib golden producer/properties
// blob").
func goldenFrame(t *testing.T) ([]byte, []string) {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("testdata", "producer_properties.hex"))
	if err != nil {
		t.Fatalf("reading golden frame: %v", err)
	}
	frame, err := hex.DecodeString(strings.TrimSpace(string(raw)))
	if err != nil {
		t.Fatalf("golden frame is not hex: %v", err)
	}
	return frame, []string{
		"com.astrate.test.PropertyArrays/config/thresholds",
		"com.astrate.test.PropertyArrays/config/labels",
		"com.example.Other/a/b",
	}
}

// TestControlFrameGolden: the foreign-built golden frame inflates to the
// expected entries, and our own framing round-trips through the parser.
func TestControlFrameGolden(t *testing.T) {
	frame, want := goldenFrame(t)
	got, err := inflateProperties(frame)
	if err != nil {
		t.Fatalf("inflateProperties(golden): %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("golden entries: %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("entry %d: %q, want %q", i, got[i], want[i])
		}
	}

	// Round trip through our builder.
	ours, err := deflateProperties(want)
	if err != nil {
		t.Fatalf("deflateProperties: %v", err)
	}
	back, err := inflateProperties(ours)
	if err != nil {
		t.Fatalf("inflateProperties(ours): %v", err)
	}
	if strings.Join(back, ";") != strings.Join(want, ";") {
		t.Errorf("round trip: %v, want %v", back, want)
	}

	// The frame header declares the uncompressed size, big-endian.
	if declared := binary.BigEndian.Uint32(ours[:4]); int(declared) != len(strings.Join(want, ";")) {
		t.Errorf("declared size %d, want %d", declared, len(strings.Join(want, ";")))
	}

	// The empty list round-trips too (the SDK deliberately sends it).
	empty, err := deflateProperties(nil)
	if err != nil {
		t.Fatalf("deflateProperties(nil): %v", err)
	}
	if got, err := inflateProperties(empty); err != nil || len(got) != 0 {
		t.Errorf("empty list round trip: %v, %v", got, err)
	}
}

// TestInflatePropertiesRejects covers the §4.5 bounds: truncated frames,
// declared sizes above the absolute ceiling (the 1 GiB zip bomb of
// docs/ROADMAP.md §7.3), streams that lie about their size, and garbage.
func TestInflatePropertiesRejects(t *testing.T) {
	deflate := func(s string) []byte {
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		if _, err := zw.Write([]byte(s)); err != nil {
			t.Fatal(err)
		}
		if err := zw.Close(); err != nil {
			t.Fatal(err)
		}
		return buf.Bytes()
	}
	frame := func(declared uint32, z []byte) []byte {
		return append(binary.BigEndian.AppendUint32(nil, declared), z...)
	}

	cases := []struct {
		name string
		in   []byte
	}{
		{name: "empty payload", in: nil},
		{name: "truncated header", in: []byte{0, 0, 1}},
		{name: "zip bomb declared 1 GiB", in: frame(1<<30, deflate(strings.Repeat("a", 64)))},
		{name: "declared above ceiling by one", in: frame(maxControlInflated+1, deflate("x"))},
		{name: "lying header", in: frame(4, deflate(strings.Repeat("b", 100)))},
		{name: "not zlib", in: frame(8, []byte("garbage!"))},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := inflateProperties(tc.in); err == nil {
				t.Fatalf("inflateProperties accepted %q: %v", tc.name, got)
			}
		})
	}

	// Boundary sanity: a payload exactly at its declared size is accepted.
	s := strings.Repeat("c", 512)
	if got, err := inflateProperties(frame(512, deflate(s))); err != nil || len(got) != 1 || got[0] != s {
		t.Errorf("exact-size payload rejected: %v, %v", got, err)
	}
}

// TestEmptyCache: the device receives every server-owned property on its
// data topic (QoS 2, hint format) followed by the consumer/properties purge
// (docs/DESIGN.md §3.3), and the §3.5.4 hint reset is armed.
func TestEmptyCache(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})
	const serverIface = "com.astrate.test.ServerProperties"

	fs.setProperty(store.Property{
		RealmID: realmAlphaID, DeviceID: devAlpha, InterfaceID: ifaceServerProps,
		Path: "/limits/maxConnections", Value: []byte("42"), ValueType: interfaceschema.Integer,
	})
	fs.setProperty(store.Property{
		RealmID: realmAlphaID, DeviceID: devAlpha, InterfaceID: ifaceServerProps,
		Path: "/identity/displayName", Value: []byte(`"panel-7"`), ValueType: interfaceschema.String,
	})
	// A device-owned property must NOT be re-sent.
	fs.setProperty(store.Property{
		RealmID: realmAlphaID, DeviceID: devAlpha, InterfaceID: ifacePropArrays,
		Path: "/config/thresholds", Value: []byte("[1.5]"), ValueType: interfaceschema.DoubleArray,
	})

	ack := &ackCounter{}
	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
	if !ack.acked() {
		t.Fatal("emptyCache not acknowledged")
	}

	base := realmAlpha + "/" + devAlpha.String()
	intPubs := port.publishedTo(base + "/" + serverIface + "/limits/maxConnections")
	if len(intPubs) != 1 {
		t.Fatalf("integer property resends: %d, want 1 (all: %+v)", len(intPubs), port.published())
	}
	wantInt, err := payload.Encode(int32(42), nil, payload.FormatBSON)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(intPubs[0].payload, wantInt) {
		t.Errorf("integer resend payload: %x, want %x", intPubs[0].payload, wantInt)
	}
	if intPubs[0].qos != 2 || intPubs[0].retain {
		t.Errorf("integer resend qos=%d retain=%t, want qos=2 retain=false", intPubs[0].qos, intPubs[0].retain)
	}
	if got := port.publishedTo(base + "/" + serverIface + "/identity/displayName"); len(got) != 1 {
		t.Errorf("string property resends: %d, want 1", len(got))
	}
	for _, p := range port.published() {
		if strings.Contains(p.topic, "PropertyArrays") {
			t.Errorf("device-owned property re-sent on %s", p.topic)
		}
	}

	purges := port.publishedTo(base + "/control/consumer/properties")
	if len(purges) != 1 {
		t.Fatalf("consumer/properties messages: %d, want 1", len(purges))
	}
	entries, err := inflateProperties(purges[0].payload)
	if err != nil {
		t.Fatalf("purge payload does not parse: %v", err)
	}
	want := map[string]bool{
		serverIface + "/identity/displayName":  true,
		serverIface + "/limits/maxConnections": true,
	}
	if len(entries) != len(want) {
		t.Fatalf("purge entries: %v", entries)
	}
	for _, e := range entries {
		if !want[e] {
			t.Errorf("unexpected purge entry %q", e)
		}
	}
	if purges[0].qos != 2 {
		t.Errorf("purge qos = %d, want 2", purges[0].qos)
	}

	// The §3.5.4 reset is armed.
	dev := rig.e.devices.peek(realmAlpha, devAlpha)
	dev.mu.Lock()
	armed := dev.resetHintOnBSON
	dev.mu.Unlock()
	if !armed {
		t.Error("emptyCache did not arm the hint reset")
	}
}

// TestEmptyCacheJSONHint: a device flipped to the JSON profile receives its
// emptyCache resends as JSON documents; the control frame stays zlib
// (docs/DESIGN.md §3.5.4).
func TestEmptyCacheJSONHint(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})
	fs.setProperty(store.Property{
		RealmID: realmAlphaID, DeviceID: devAlpha, InterfaceID: ifaceServerProps,
		Path: "/limits/maxConnections", Value: []byte("42"), ValueType: interfaceschema.Integer,
	})

	// First data publish in JSON flips the sticky hint.
	ack := &ackCounter{}
	rig.handle(deviceMsg("com.astrate.test.Minimal", "/value", 0, []byte(`{"v":2.5}`), ack))
	rig.e.flushShard(t.Context(), rig.sh)

	ack = &ackCounter{}
	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
	if !ack.acked() {
		t.Fatal("emptyCache not acknowledged")
	}

	base := realmAlpha + "/" + devAlpha.String()
	pubs := port.publishedTo(base + "/com.astrate.test.ServerProperties/limits/maxConnections")
	if len(pubs) != 1 {
		t.Fatalf("resends: %d, want 1", len(pubs))
	}
	wantJSON, err := payload.Encode(int32(42), nil, payload.FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pubs[0].payload, wantJSON) {
		t.Errorf("JSON-hinted resend payload: %s, want %s", pubs[0].payload, wantJSON)
	}

	purges := port.publishedTo(base + "/control/consumer/properties")
	if len(purges) != 1 {
		t.Fatalf("consumer/properties messages: %d, want 1", len(purges))
	}
	if _, err := inflateProperties(purges[0].payload); err != nil {
		t.Errorf("control frame must stay zlib for JSON devices: %v", err)
	}
}

// TestProducerPropertiesPurge: rows absent from the device's list are
// purged, listed rows and server-owned rows survive (docs/DESIGN.md §3.3).
func TestProducerPropertiesPurge(t *testing.T) {
	rig, fs, _ := newWiredRig(t, Config{})

	kept := store.PropertyRef{InterfaceID: ifacePropArrays, Path: "/config/thresholds"}
	stale := store.PropertyRef{InterfaceID: ifacePropArrays, Path: "/config/labels"}
	server := store.PropertyRef{InterfaceID: ifaceServerProps, Path: "/limits/maxConnections"}
	fs.setProperty(store.Property{RealmID: realmAlphaID, DeviceID: devAlpha,
		InterfaceID: kept.InterfaceID, Path: kept.Path, Value: []byte("[1.5]"), ValueType: interfaceschema.DoubleArray})
	fs.setProperty(store.Property{RealmID: realmAlphaID, DeviceID: devAlpha,
		InterfaceID: stale.InterfaceID, Path: stale.Path, Value: []byte(`["a"]`), ValueType: interfaceschema.StringArray})
	fs.setProperty(store.Property{RealmID: realmAlphaID, DeviceID: devAlpha,
		InterfaceID: server.InterfaceID, Path: server.Path, Value: []byte("42"), ValueType: interfaceschema.Integer})

	frame, err := deflateProperties([]string{
		"com.astrate.test.PropertyArrays/config/thresholds",
		"com.astrate.test.UnknownInterface/whatever", // skipped, logged
	})
	if err != nil {
		t.Fatal(err)
	}
	ack := &ackCounter{}
	rig.handle(deviceMsg("control", "/producer/properties", 2, frame, ack))
	if !ack.acked() {
		t.Fatal("producer/properties not acknowledged")
	}

	refs := fs.propertyRefs(realmAlphaID, devAlpha)
	if len(refs) != 2 {
		t.Fatalf("surviving rows: %+v, want kept + server-owned", refs)
	}
	if refs[0] != server && refs[1] != server {
		t.Error("server-owned property was purged")
	}
	if refs[0] != kept && refs[1] != kept {
		t.Error("listed device-owned property was purged")
	}

	fs.mu.Lock()
	calls := len(fs.purgeCalls)
	keepArg := fs.purgeCalls[0].keep
	fs.mu.Unlock()
	if calls != 1 || len(keepArg) != 1 || keepArg[0] != kept {
		t.Errorf("purge call: %d calls, keep=%+v", calls, keepArg)
	}
}

// TestControlRejects: malformed frames and unknown subpaths are rejected
// under their reasons and consumed.
func TestControlRejects(t *testing.T) {
	rig, _, _ := newWiredRig(t, Config{})

	ack := &ackCounter{}
	rig.handle(deviceMsg("control", "/producer/properties", 2, []byte{0, 0}, ack))
	if !ack.acked() {
		t.Error("malformed control payload not consumed")
	}
	if got := promtest.ToFloat64(rig.e.met.rejects.WithLabelValues(reasonControlInvalid)); got != 1 {
		t.Errorf("rejects[%s] = %v, want 1", reasonControlInvalid, got)
	}

	ack = &ackCounter{}
	rig.handle(deviceMsg("control", "/bogus", 2, []byte("1"), ack))
	if !ack.acked() {
		t.Error("unknown control subpath not consumed")
	}
	if got := promtest.ToFloat64(rig.e.met.rejects.WithLabelValues(reasonControlUnknown)); got != 1 {
		t.Errorf("rejects[%s] = %v, want 1", reasonControlUnknown, got)
	}
}
