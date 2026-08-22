package engine

// MQTT v1 upstream-parity conformance tests for the 1.3-era items #48 and
// #47/#50's siblings tracked in issues: empty introspection tolerance and
// binaryblob re-send as BSON binary subtype 0.

import (
	"bytes"
	"testing"

	promtest "github.com/prometheus/client_golang/prometheus/testutil"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// TestEmptyIntrospectionAccepted covers issue #48 end to end: an EMPTY
// introspection payload is a valid zero-interface registration — acked,
// persisted with no entries, and the broker ACL refreshed. The §2.6 gate
// must follow immediately (a previously declared interface stops being
// publishable). Adversarially: a lone space is near-empty garbage and must
// be rejected, while a well-formed-but-uninstalled interface is accepted at
// this layer (existence is enforced later at the data gate).
func TestEmptyIntrospectionAccepted(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})

	// Empty payload on the bare introspection topic: zero interfaces.
	ack := &ackCounter{}
	rig.handle(deviceMsg("", "", 2, nil, ack))
	if !ack.acked() {
		t.Fatal("empty introspection not acknowledged")
	}
	if got := promtest.ToFloat64(rig.e.met.rejects.WithLabelValues(reasonIntrospectionInvalid)); got != 0 {
		t.Errorf("rejects[%s] = %v, want 0", reasonIntrospectionInvalid, got)
	}

	fs.mu.Lock()
	entries := len(fs.devices[deviceKey{realm: realmAlpha, id: devAlpha}].Introspection)
	fs.mu.Unlock()
	if entries != 0 {
		t.Errorf("persisted introspection: %d entries, want exactly 0", entries)
	}

	if n := port.refreshCount(); n != 1 {
		t.Errorf("broker introspection refreshes: %d, want 1", n)
	}

	// The empty introspection really took effect: a previously declared
	// interface is rejected at the step-1 gate.
	ack = &ackCounter{}
	rig.handle(deviceMsg("com.astrate.test.AllScalarTypes", "/double", 1, enc(t, 1.5, nil, payload.FormatBSON), ack))
	if got := promtest.ToFloat64(rig.e.met.rejects.WithLabelValues(reasonInterfaceNotDeclared)); got != 1 {
		t.Errorf("rejects[%s] = %v, want 1", reasonInterfaceNotDeclared, got)
	}

	// A lone space parses like garbage, not like an empty introspection.
	ack = &ackCounter{}
	rig.handle(deviceMsg("", "", 2, []byte(" "), ack))
	if !ack.acked() {
		t.Error("space-only introspection not consumed")
	}
	if got := promtest.ToFloat64(rig.e.met.rejects.WithLabelValues(reasonIntrospectionInvalid)); got != 1 {
		t.Errorf("rejects[%s] = %v, want 1 after space-only payload", reasonIntrospectionInvalid, got)
	}

	// Well-formed but uninstalled: accepted at the introspection layer.
	ack = &ackCounter{}
	rig.handle(deviceMsg("", "", 2, []byte("com.ex.A:1:0"), ack))
	if !ack.acked() {
		t.Error("well-formed uninstalled introspection not acknowledged")
	}
}

// TestBinaryBlobResendBSONSubtypeZero pins issue #49 on the actual
// properties re-send path: a re-sent server-owned binaryblob property must
// arrive as a BSON document whose v element is binary subtype 0x00
// (generic), never 1/2/other. Asserted twice: through the high-level
// decoder and directly on the raw wire octet after the element's length
// prefix.
func TestBinaryBlobResendBSONSubtypeZero(t *testing.T) {
	rig, fs, port := newWiredRig(t, Config{})

	fs.setProperty(store.Property{
		RealmID: realmAlphaID, DeviceID: devAlpha, InterfaceID: ifaceServerProps,
		Path:      "/identity/blob",
		Value:     []byte(`"AAEC/v8="`),
		ValueType: interfaceschema.BinaryBlob,
	})

	ack := &ackCounter{}
	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
	if !ack.acked() {
		t.Fatal("emptyCache not acknowledged")
	}

	topic := realmAlpha + "/" + devAlpha.String() + "/com.astrate.test.ServerProperties/identity/blob"
	pubs := port.publishedTo(topic)
	if len(pubs) != 1 {
		t.Fatalf("resend publications: %d, want 1", len(pubs))
	}
	pub := pubs[0]
	if pub.qos != 2 {
		t.Errorf("resend qos = %d, want 2", pub.qos)
	}
	if pub.retain {
		t.Error("resend published with retain set")
	}

	// High-level check: the decoded v value is BSON binary subtype 0x00.
	rv := bson.Raw(pub.payload).Lookup("v")
	if rv.Type != bson.TypeBinary {
		t.Fatalf("re-sent blob v type = %v, want %v", rv.Type, bson.TypeBinary)
	}
	subtype, data, ok := rv.BinaryOK()
	if !ok {
		t.Fatal("re-sent blob v is not a BSON binary value")
	}
	if subtype != bson.TypeBinaryGeneric {
		t.Errorf("binary subtype = %#x, want %#x (generic)", subtype, bson.TypeBinaryGeneric)
	}
	if subtype != 0x00 {
		t.Errorf("binary subtype = %#x, want exactly 0x00", subtype)
	}
	want := []byte{0x00, 0x01, 0x02, 0xFE, 0xFF}
	if !bytes.Equal(data, want) {
		t.Errorf("binary data = %#v, want %#v", data, want)
	}

	// Wire-octet check: scan the raw document for the binary element header
	// (element type 0x05 followed by key "v"); the octet immediately after
	// the 4-byte little-endian length prefix is the subtype itself and must
	// be 0x00, so an encoder emitting a valid but different subtype cannot
	// slip past the decoder check above.
	raw := pub.payload
	found := false
	for i := 0; i+8 <= len(raw); i++ {
		if raw[i] != 0x05 || raw[i+1] != 'v' || raw[i+2] != 0x00 {
			continue
		}
		found = true
		if st := raw[i+7]; st != 0x00 {
			t.Errorf("wire subtype octet = %#x, want 0x00", st)
		}
		break
	}
	if !found {
		t.Fatal(`no binary element with key "v" found in the raw resend payload`)
	}
}
