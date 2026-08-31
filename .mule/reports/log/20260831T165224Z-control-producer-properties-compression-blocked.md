slug: control-producer-properties-compression
verdict: blocked
at:  d044364
ran: 2026-08-31T16:52:25Z on DietPi in 332s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
I understand the task: make `inflateProperties` in `internal/engine/control.go` accept plaintext payloads (and the 4-zero-byte empty frame) when the device declared `purge_properties_compression_format: plaintext`, instead of rejecting all non-zlib input. Let me read the relevant code.
→ Read internal/engine/control.go
✱ Grep "inflateProperties" in . · 20 matches
Now let me find the callers and understand the compression capability system.
✱ Grep "purgeCompressionFor|compressionPlaintext|compressionZlib|purge_properties_compression" in . · 74 matches
Let me read the upstream `control_handler.ex` reference and the existing test file.
→ Read .mule/reports/log/20260803T021148Z-recipe-astarte-upstream-proposed.md [offset=605, limit=15]
→ Read internal/engine/control_test.go
→ Read internal/engine/cache.go [offset=405, limit=15]
Let me check how the device's compression format is accessed in the handler context.
→ Read internal/engine/capabilities_test.go [offset=185, limit=55]
✱ Grep "func \(e \*Engine\) deviceState" in internal/engine · 1 match
→ Read .mule/reports/log/20260803T021148Z-recipe-astarte-upstream-proposed.md [offset=595, limit=25]
✱ Grep "decode_payload.*plaintext|:plaintext|\"plaintext\".*decode" in .mule/reports/log · 3 matches
✱ Grep "decode_payload|:zlib|:plaintext" in .mule/reports/log/20260810 · 14 matches
→ Read .mule/reports/log/20260803T021148Z-recipe-astarte-upstream-proposed.md [offset=698, limit=20]
→ Read .mule/reports/log/20260810T021151Z-recipe-astarte-upstream-proposed.md [offset=586, limit=20]
→ Read internal/engine/data.go [offset=209, limit=30]

# Todos
[•] Modify inflateProperties to accept plaintext flag
[ ] Update handleProducerProperties to look up device compression format first
[ ] Update all existing callers of inflateProperties
[ ] Write test proving plaintext works and fails without the change
[ ] Run gate: go test ./... && gofmt -l .

Now I have a clear picture. The upstream `decode_payload` does:

- **zlib**: strip 4-byte header, inflate the rest
- **plaintext**: return the raw payload as-is (no header, no decompression)

I'll:
1. Add a `plaintext bool` param to `inflateProperties`
2. Move `deviceState` before `inflateProperties` in `handleProducerProperties` to look up the format
3. Handle plaintext: split raw `;`-joined bytes (empty → nil)
4. Write a test that proves plaintext works and fails without the change
5. Update all existing callers
← Edit internal/engine/control.go
Index: /root/astrate-mule/internal/engine/control.go
===================================================================
--- /root/astrate-mule/internal/engine/control.go
+++ /root/astrate-mule/internal/engine/control.go
@@ -280,15 +280,36 @@
 func plaintextProperties(entries []string) []byte {
 	return []byte(strings.Join(entries, ";"))
 }
 
-// inflateProperties parses a zlib control payload with the docs/DESIGN.md
-// §4.5 bounds: the declared size is capped by maxControlInflated, and the
-// stream may not inflate beyond what it declared (a lying header is a
-// zip-bomb attempt, not a tolerable client quirk). It returns the entry
-// list; an empty payload yields no entries (the SDK deliberately sends an
-// empty list when the device holds no properties).
-func inflateProperties(p []byte) ([]string, error) {
+// inflateProperties parses a device→server control payload (producer/properties
+// or consumer/properties purge list) and returns the entry list.
+//
+// When plaintext is false the payload is zlib-compressed with the docs/DESIGN.md
+// §4.5 bounds: a 4-byte big-endian declared-size header followed by a zlib
+// stream. The declared size is capped by maxControlInflated and the stream may
+// not inflate beyond what it declared.
+//
+// When plaintext is true the payload is the raw `;`-joined entry list with no
+// framing — the format used when the device declared
+// purge_properties_compression_format: plaintext. The 4-zero-byte payload
+// (\x00\x00\x00\x00) is accepted as an empty list (the device's representation
+// of "no properties" in plaintext mode).
+//
+// An empty payload in either mode yields no entries (the SDK deliberately sends
+// an empty list when the device holds no properties).
+func inflateProperties(p []byte, plaintext bool) ([]string, error) {
+	if plaintext {
+		// The 4-zero-byte empty frame is a known device encoding for an
+		// empty property list in plaintext mode.
+		if len(p) == controlFrameHeader && binary.BigEndian.Uint32(p) == 0 {
+			return nil, nil
+		}
+		if len(p) == 0 {
+			return nil, nil
+		}
+		return strings.Split(string(p), ";"), nil
+	}
 	if len(p) < controlFrameHeader {
 		return nil, fmt.Errorf("control payload is %d bytes, below the %d byte size prefix", len(p), controlFrameHeader)
 	}
 	declared := binary.BigEndian.Uint32(p[:controlFrameHeader])

# Todos
[✓] Modify inflateProperties to accept plaintext flag
[•] Update handleProducerProperties to look up device compression format first
[ ] Update all existing callers of inflateProperties
[ ] Write test proving plaintext works and fails without the change
[ ] Run gate: go test ./... && gofmt -l .

← Edit internal/engine/control.go
Index: /root/astrate-mule/internal/engine/control.go
===================================================================
--- /root/astrate-mule/internal/engine/control.go
+++ /root/astrate-mule/internal/engine/control.go
@@ -79,17 +79,17 @@
 // row not in it is purged — the device is the source of truth for its own
 // properties. Malformed payloads (bad frame, zip bomb) are rejected and
 // consumed.
 func (e *Engine) handleProducerProperties(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
-	entries, err := inflateProperties(m.Payload)
-	if err != nil {
-		e.reject(m, reasonControlInvalid, err.Error())
-		return
-	}
 	dev, ok := e.deviceState(ctx, m, realm)
 	if !ok {
 		return
 	}
+	entries, err := inflateProperties(m.Payload, purgeCompressionFor(dev) == compressionPlaintext)
+	if err != nil {
+		e.reject(m, reasonControlInvalid, err.Error())
+		return
+	}
 
 	keep := e.resolvePropertyRefs(realm, dev, entries)
 	if !e.retryStore(ctx, m, "producer/properties purge", func() error {
 		purged, err := e.st.PurgeDeviceOwnedExcept(ctx, realm.id, m.DeviceID, keep)

← Edit internal/engine/control_test.go
Index: /root/astrate-mule/internal/engine/control_test.go
===================================================================
--- /root/astrate-mule/internal/engine/control_test.go
+++ /root/astrate-mule/internal/engine/control_test.go
@@ -41,9 +41,9 @@
 // TestControlFrameGolden: the foreign-built golden frame inflates to the
 // expected entries, and our own framing round-trips through the parser.
 func TestControlFrameGolden(t *testing.T) {
 	frame, want := goldenFrame(t)
-	got, err := inflateProperties(frame)
+	got, err := inflateProperties(frame, false)
 	if err != nil {
 		t.Fatalf("inflateProperties(golden): %v", err)
 	}
 	if len(got) != len(want) {

← Edit internal/engine/control_test.go
Index: /root/astrate-mule/internal/engine/control_test.go
===================================================================
--- /root/astrate-mule/internal/engine/control_test.go
+++ /root/astrate-mule/internal/engine/control_test.go
@@ -59,9 +59,9 @@
 ours, err := deflateProperties(want)
 if err != nil {
 	t.Fatalf("deflateProperties: %v", err)
 }
-back, err := inflateProperties(ours)
+back, err := inflateProperties(ours, false)
 if err != nil {
 	t.Fatalf("inflateProperties(ours): %v", err)
 }
 if strings.Join(back, ";") != strings.Join(want, ";") {

← Edit internal/engine/control_test.go
Index: /root/astrate-mule/internal/engine/control_test.go
===================================================================
--- /root/astrate-mule/internal/engine/control_test.go
+++ /root/astrate-mule/internal/engine/control_test.go
@@ -77,9 +77,9 @@
 	empty, err := deflateProperties(nil)
 	if err != nil {
 		t.Fatalf("deflateProperties(nil): %v", err)
 	}
-	if got, err := inflateProperties(empty); err != nil || len(got) != 0 {
+	if got, err := inflateProperties(empty, false); err != nil || len(got) != 0 {
 		t.Errorf("empty list round trip: %v, %v", got, err)
 	}
 }
 

← Edit internal/engine/control_test.go
Index: /root/astrate-mule/internal/engine/control_test.go
===================================================================
--- /root/astrate-mule/internal/engine/control_test.go
+++ /root/astrate-mule/internal/engine/control_test.go
@@ -114,9 +114,9 @@
 	{name: "not zlib", in: frame(8, []byte("garbage!"))},
 }
 for _, tc := range cases {
 	t.Run(tc.name, func(t *testing.T) {
-		if got, err := inflateProperties(tc.in); err == nil {
+		if got, err := inflateProperties(tc.in, false); err == nil {
 			t.Fatalf("inflateProperties accepted %q: %v", tc.name, got)
 		}
 	})
 }

← Edit internal/engine/control_test.go
Index: /root/astrate-mule/internal/engine/control_test.go
===================================================================
--- /root/astrate-mule/internal/engine/control_test.go
+++ /root/astrate-mule/internal/engine/control_test.go
@@ -122,9 +122,9 @@
 	}
 
 	// Boundary sanity: a payload exactly at its declared size is accepted.
 	s := strings.Repeat("c", 512)
-	if got, err := inflateProperties(frame(512, deflate(s))); err != nil || len(got) != 1 || got[0] != s {
+	if got, err := inflateProperties(frame(512, deflate(s)), false); err != nil || len(got) != 1 || got[0] != s {
 		t.Errorf("exact-size payload rejected: %v, %v", got, err)
 	}
 }
 

← Edit internal/engine/control_test.go
Index: /root/astrate-mule/internal/engine/control_test.go
===================================================================
--- /root/astrate-mule/internal/engine/control_test.go
+++ /root/astrate-mule/internal/engine/control_test.go
@@ -182,9 +182,9 @@
 purges := port.publishedTo(base + "/control/consumer/properties")
 if len(purges) != 1 {
 	t.Fatalf("consumer/properties messages: %d, want 1", len(purges))
 }
-entries, err := inflateProperties(purges[0].payload)
+entries, err := inflateProperties(purges[0].payload, false)
 if err != nil {
 	t.Fatalf("purge payload does not parse: %v", err)
 }
 want := map[string]bool{

Now update all test callers of `inflateProperties`. They all use zlib, so I pass `false`.
← Edit internal/engine/control_test.go
Index: /root/astrate-mule/internal/engine/control_test.go
===================================================================
--- /root/astrate-mule/internal/engine/control_test.go
+++ /root/astrate-mule/internal/engine/control_test.go
@@ -250,9 +250,9 @@
 	purges := port.publishedTo(base + "/control/consumer/properties")
 	if len(purges) != 1 {
 		t.Fatalf("consumer/properties messages: %d, want 1", len(purges))
 	}
-	if _, err := inflateProperties(purges[0].payload); err != nil {
+	if _, err := inflateProperties(purges[0].payload, false); err != nil {
 		t.Errorf("control frame must stay zlib for JSON devices: %v", err)
 	}
 }
 

→ Read internal/engine/serverdata_test.go [offset=245, limit=10]
Now update callers in the other test files.
→ Read internal/engine/capabilities_test.go [offset=185, limit=45]

← Edit internal/engine/serverdata_test.go
Index: /root/astrate-mule/internal/engine/serverdata_test.go
===================================================================
--- /root/astrate-mule/internal/engine/serverdata_test.go
+++ /root/astrate-mule/internal/engine/serverdata_test.go
@@ -245,9 +245,9 @@
 purges := port.publishedTo(base + "/control/consumer/properties")
 if len(purges) != 1 {
 	t.Fatalf("purge messages after unset: %d, want 1", len(purges))
 }
-entries, err := inflateProperties(purges[0].payload)
+entries, err := inflateProperties(purges[0].payload, false)
 if err != nil {
 	t.Fatalf("purge payload: %v", err)
 }
 if len(entries) != 1 || entries[0] != iface+"/limits/maxPayload" {

← Edit internal/engine/capabilities_test.go
Index: /root/astrate-mule/internal/engine/capabilities_test.go
===================================================================
--- /root/astrate-mule/internal/engine/capabilities_test.go
+++ /root/astrate-mule/internal/engine/capabilities_test.go
@@ -190,9 +190,9 @@
 		t.Errorf("plaintext purge = %q, want %q", payload, want)
 	}
 
 	// Verify it is NOT a valid zlib frame (the 4-byte header + zlib stream).
-	if _, err := inflateProperties(purges[0].payload); err == nil {
+	if _, err := inflateProperties(purges[0].payload, false); err == nil {
 		t.Error("plaintext purge is a valid zlib frame; should be raw text")
 	}
 }
 

← Edit internal/engine/capabilities_test.go
Index: /root/astrate-mule/internal/engine/capabilities_test.go
===================================================================
--- /root/astrate-mule/internal/engine/capabilities_test.go
+++ /root/astrate-mule/internal/engine/capabilities_test.go
@@ -218,9 +218,9 @@
 	t.Fatalf("consumer/properties messages: %d, want 1", len(purges))
 }
 
 // Must be a valid zlib frame.
-entries, err := inflateProperties(purges[0].payload)
+entries, err := inflateProperties(purges[0].payload, false)
 if err != nil {
 	t.Fatalf("default purge is not zlib: %v", err)
 }
 want := serverIface + "/limits/maxConnections"

# Todos
[✓] Modify inflateProperties to accept plaintext flag
[✓] Update handleProducerProperties to look up device compression format first
[✓] Update all existing callers of inflateProperties
[•] Write test proving plaintext works and fails without the change
[ ] Run gate: go test ./... && gofmt -l .

Now I need to write a test. Let me first look at how the rig test setup works for capabilities to understand how to set `purge_properties_compression_format: plaintext` on a device.
→ Read internal/engine/capabilities_test.go [limit=100, offset=80]
I understand the test pattern. Now let me verify the gate compiles first before writing the new test, so I can confirm the refactor is correct.

$ cd /root/astrate-mule && go build ./internal/engine/ 2>&1 | head -30
(no output)

The build passes. Now I'll write the new test. It should prove that a plaintext-device's `producer/properties` payload (raw `;`-joined list and the 4-zero-byte empty frame) is accepted end-to-end, and that a zlib-device's raw plaintext is still rejected.

Let me add a test modeled on `TestProducerPropertiesPurge` in `control_test.go`:
← Edit internal/engine/control_test.go
Index: /root/astrate-mule/internal/engine/control_test.go
===================================================================
--- /root/astrate-mule/internal/engine/control_test.go
+++ /root/astrate-mule/internal/engine/control_test.go
@@ -257,9 +257,8 @@
 }
 
 // TestProducerPropertiesPurge: rows absent from the device's list are
 // purged, listed rows and server-owned rows survive (docs/DESIGN.md §3.3).
-func TestProducerPropertiesPurge(t *testing.T) {
 	rig, fs, _ := newWiredRig(t, Config{})
 
 	kept := store.PropertyRef{InterfaceID: ifacePropArrays, Path: "/config/thresholds"}
 	stale := store.PropertyRef{InterfaceID: ifacePropArrays, Path: "/config/labels"}

Oops, that edit removed the function signature line. Let me undo that by restoring it.
→ Read internal/engine/control_test.go [limit=10, offset=259]

← Edit internal/engine/control_test.go
Index: /root/astrate-mule/internal/engine/control_test.go
===================================================================
--- /root/astrate-mule/internal/engine/control_test.go
+++ /root/astrate-mule/internal/engine/control_test.go
@@ -257,8 +257,9 @@
 }
 
 // TestProducerPropertiesPurge: rows absent from the device's list are
 // purged, listed rows and server-owned rows survive (docs/DESIGN.md §3.3).
+func TestProducerPropertiesPurge(t *testing.T) {
 	rig, fs, _ := newWiredRig(t, Config{})
 
 	kept := store.PropertyRef{InterfaceID: ifacePropArrays, Path: "/config/thresholds"}
 	stale := store.PropertyRef{InterfaceID: ifacePropArrays, Path: "/config/labels"}

Restored. Now let me add the new plaintext test after the existing `TestProducerPropertiesPurge`. Let me view the end of that function to insert cleanly.
→ Read internal/engine/control_test.go [limit=30, offset=296]
Let me check how `mustBSON` is imported in control_test.go (it's in capabilities_test.go). Let me check the bson helper availability.
✱ Grep "func mustBSON|go.mongodb.org" in internal/engine · 3 matches
`bson` and `mustBSON` are package-level in the engine test package, so I can use them in control_test.go by importing bson. But adding the bson import to control_test.go. Actually, better to write the plaintext test in `capabilities_test.go` where bson is already imported and the plaintext purge pattern already exists. That's the more cohesive home. Let me add an end-to-end test there.

Let me write the test in `capabilities_test.go` after `TestPlaintextConsumerProperties`:
→ Read internal/engine/capabilities_test.go [limit=30, offset=225]
I'll add my end-to-end tests here in `capabilities_test.go`, after `TestZlibConsumerPropertiesDefault` (line 230) and before `TestPurgeCompressionFor`. This test proves the whole `producer/properties` path accepts plaintext for a plaintext device and the 4-zero-byte empty frame.

← Edit internal/engine/capabilities_test.go
Index: /root/astrate-mule/internal/engine/capabilities_test.go
===================================================================
--- /root/astrate-mule/internal/engine/capabilities_test.go
+++ /root/astrate-mule/internal/engine/capabilities_test.go
@@ -228,8 +228,102 @@
 		t.Errorf("purge entries = %v, want [%s]", entries, want)
 	}
 }
 
+// TestPlaintextProducerProperties: a device that declares
+// purge_properties_compression_format: plaintext sends its producer/properties
+// purge list as raw `;`-joined text (no zlib frame), and the 4-zero-byte
+// all-zero empty frame for an empty list. Both must be decoded per-capability
+// (upstream control_handler.ex), not rejected as non-zlib.
+func TestPlaintextProducerProperties(t *testing.T) {
+	rig, fs, _ := newWiredRig(t, Config{})
+
+	kept := store.PropertyRef{InterfaceID: ifacePropArrays, Path: "/config/thresholds"}
+	stale := store.PropertyRef{InterfaceID: ifacePropArrays, Path: "/config/labels"}
+	server := store.PropertyRef{InterfaceID: ifaceServerProps, Path: "/limits/maxConnections"}
+	fs.setProperty(store.Property{RealmID: realmAlphaID, DeviceID: devAlpha,
+		InterfaceID: kept.InterfaceID, Path: kept.Path, Value: []byte("[1.5]"), ValueType: interfaceschema.DoubleArray})
+	fs.setProperty(store.Property{RealmID: realmAlphaID, DeviceID: devAlpha,
+		InterfaceID: stale.InterfaceID, Path: stale.Path, Value: []byte(`["a"]`), ValueType: interfaceschema.StringArray})
+	fs.setProperty(store.Property{RealmID: realmAlphaID, DeviceID: devAlpha,
+		InterfaceID: server.InterfaceID, Path: server.Path, Value: []byte("42"), ValueType: interfaceschema.Integer})
+
+	// Declare plaintext compression.
+	ack := &ackCounter{}
+	rig.handle(deviceMsg("capabilities", "", 2, mustBSON(t, bson.M{
+		"purge_properties_compression_format": "plaintext",
+	}), ack))
+	if !ack.acked() {
+		t.Fatal("capabilities not acknowledged")
+	}
+
+	// Send the purge list as raw plaintext.
+	raw := []byte("com.astrate.test.PropertyArrays/config/thresholds")
+	ack = &ackCounter{}
+	rig.handle(deviceMsg("control", "/producer/properties", 2, raw, ack))
+	if !ack.acked() {
+		t.Fatal("plaintext producer/properties not acknowledged")
+	}
+
+	refs := fs.propertyRefs(realmAlphaID, devAlpha)
+	if len(refs) != 2 {
+		t.Fatalf("surviving rows: %+v, want kept + server-owned", refs)
+	}
+	if refs[0] != server && refs[1] != server {
+		t.Error("server-owned property was purged")
+	}
+	if refs[0] != kept && refs[1] != kept {
+		t.Error("listed device-owned property was purged")
+	}
+}
+
+// TestPlaintextProducerPropertiesEmptyFrame: the 4-zero-byte all-zero frame is
+// the plaintext device's encoding of "I hold no properties", and must decode
+// to an empty keep list rather than being rejected as not-zlib.
+func TestPlaintextProducerPropertiesEmptyFrame(t *testing.T) {
+	rig, fs, _ := newWiredRig(t, Config{})
+
+	// A device-owned row the purge should eliminate (not in the empty list).
+	fs.setProperty(store.Property{RealmID: realmAlphaID, DeviceID: devAlpha,
+		InterfaceID: ifacePropArrays, Path: "/config/thresholds", Value: []byte("[1.5]"), ValueType: interfaceschema.DoubleArray})
+
+	ack := &ackCounter{}
+	rig.handle(deviceMsg("capabilities", "", 2, mustBSON(t, bson.M{
+		"purge_properties_compression_format": "plaintext",
+	}), ack))
+	if !ack.acked() {
+		t.Fatal("capabilities not acknowledged")
+	}
+
+	// All-zero 4-byte empty frame.
+	ack = &ackCounter{}
+	rig.handle(deviceMsg("control", "/producer/properties", 2, []byte{0, 0, 0, 0}, ack))
+	if !ack.acked() {
+		t.Fatal("empty-frame producer/properties not acknowledged")
+	}
+
+	refs := fs.propertyRefs(realmAlphaID, devAlpha)
+	if len(refs) != 0 {
+		t.Fatalf("surviving rows after empty purge: %+v, want none", refs)
+	}
+}
+
+// TestZlibDeviceRejectsPlaintext producer list: a device that has NOT declared
+// plaintext (default zlib) must still be rejected when it sends an unpacked
+// list — the per-capability decode must not leak into zlib devices.
+func TestZlibDeviceRejectsPlaintext(t *testing.T) {
+	rig, _, _ := newWiredRig(t, Config{})
+
+	ack := &ackCounter{}
+	rig.handle(deviceMsg("control", "/producer/properties", 2, []byte("com.astrate.test.PropertyArrays/config/thresholds"), ack))
+	if !ack.acked() {
+		t.Error("invalid producer/properties payload not consumed")
+	}
+	if got := promtest.ToFloat64(rig.e.met.rejects.WithLabelValues(reasonControlInvalid)); got != 1 {
+		t.Errorf("rejects[%s] = %v, want 1", reasonControlInvalid, got)
+	}
+}
+
 // TestPurgeCompressionFor tests the fallback logic.
 func TestPurgeCompressionFor(t *testing.T) {
 	if got := purgeCompressionFor(nil); got != compressionZlib {
 		t.Errorf("purgeCompressionFor(nil) = %q, want zlib", got)

Now let me verify the gate. Let me run the tests for the engine package first.
$ cd /root/astrate-mule && go test ./internal/engine/ 2>&1 | tail -30
# github.com/astrate-platform/astrate/internal/engine [github.com/astrate-platform/astrate/internal/engine.test]
internal/engine/capabilities_test.go:244:86: undefined: interfaceschema
internal/engine/capabilities_test.go:246:88: undefined: interfaceschema
internal/engine/capabilities_test.go:248:87: undefined: interfaceschema
internal/engine/capabilities_test.go:287:96: undefined: interfaceschema
FAIL	github.com/astrate-platform/astrate/internal/engine [build failed]
FAIL

I need to import `interfaceschema` in capabilities_test.go. Let me check its imports.
→ Read internal/engine/capabilities_test.go [limit=15]

← Edit internal/engine/capabilities_test.go
Index: /root/astrate-mule/internal/engine/capabilities_test.go
===================================================================
--- /root/astrate-mule/internal/engine/capabilities_test.go
+++ /root/astrate-mule/internal/engine/capabilities_test.go
@@ -6,8 +6,9 @@
 	promtest "github.com/prometheus/client_golang/prometheus/testutil"
 	"go.mongodb.org/mongo-driver/v2/bson"
 
 	"github.com/astrate-platform/astrate/internal/store"
+	"github.com/astrate-platform/astrate/pkg/interfaceschema"
 )
 
 func TestDecodeCapabilities(t *testing.T) {
 	cases := []struct {

$ cd /root/astrate-mule && go test ./internal/engine/ 2>&1 | tail -30
ok  	github.com/astrate-platform/astrate/internal/engine	1.595s

Tests pass. Now let me verify my new tests actually fail without the implementation change (per MULE.md requirements). Let me temporarily revert the implementation and run the new tests.

First, let me confirm the new tests cover the change. Let me run the full gate (`go test ./...` and `gofmt -l .`) as required.
$ cd /root/astrate-mule && gofmt -l . && echo "---gofmt done---" && go test ./... 2>&1 | tail -30
---gofmt done---
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
?   	github.com/astrate-platform/astrate/internal/housekeeping	[no test files]
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
?   	github.com/astrate-platform/astrate/internal/swagger	[no test files]
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

Gate passes and gofmt is clean. Now, per MULE.md requirement, I must verify my new tests fail without the implementation change. Let me temporarily revert the implementation to prove the tests break.

Let me save the current control.go, revert the implementation change, run the new tests specifically, then restore.

The new tests are `TestPlaintextProducerProperties`, `TestPlaintextProducerPropertiesEmptyFrame`. The implementation change is: (1) the `plaintext` branch in `inflateProperties`, and (2) `handleProducerProperties` passing `purgeCompressionFor(dev) == compressionPlaintext`.

Without these, the plaintext device's raw payload (no 4-byte header, 33 bytes) would be interpreted as zlib frame → rejected. Let me verify by stashing the implementation.
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cd /root/astrate-mule && cp internal/engine/control.go /tmp/control.go.bak && cp internal/engine/control_test.go /tmp/control_test.go.bak && cp internal/engine/capabilities_test.go /tmp/capabilities_test.go.bak && echo saved failed
Error: The user rejected permission to use this specific tool call.
