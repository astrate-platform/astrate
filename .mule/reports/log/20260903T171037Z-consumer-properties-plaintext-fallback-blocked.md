slug: consumer-properties-plaintext-fallback
verdict: blocked
at:  19d552b
ran: 2026-09-03T17:10:37Z on DietPi in 1182s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
The task: When `devices.peek` returns nil in `sendConsumerProperties`, the code defaults to zlib compression, but a plaintext-capable device can't decompress zlib. I need to fall back to a store read for the device's compression capability, and add a test.

Let me read the relevant files first.
✱ Grep "devices\.peek" in internal/engine · 4 matches
✱ Grep "plaintext" in internal/engine · 27 matches
→ Read internal/engine/control.go [offset=180, limit=80]
→ Read internal/engine/cache.go [limit=70, offset=350]
→ Read internal/engine/capabilities.go
✱ Grep "capability|purge.*compress|purgeCompression" in internal/store · 1 match
Now let me look at the store interface for capabilities and the existing test file.
→ Read internal/engine/capabilities_test.go
✱ Grep "interface Store|func.*Store\b" in internal/store · 89 matches
Let me look at the store interface and how the cache is evicted.
✱ Grep "evict|emptyCache|devices\.delete" in internal/engine · 50 matches
Let me look at the device cache and how it reloads from the store.
→ Read internal/engine/cache.go [offset=420]
→ Read internal/engine/cache.go [limit=60, offset=1]
→ Read internal/store/devices.go [limit=75, offset=1]
Let me look at the `handleEmptyCache` function and how it calls `sendConsumerProperties`.
→ Read internal/engine/control.go [limit=130, offset=1]
Let me check all callers of `sendConsumerProperties`.
✱ Grep "sendConsumerProperties" in internal/engine · 5 matches
→ Read internal/engine/serverdata.go [offset=195, limit=25]
→ Read internal/engine/engine.go [offset=255, limit=25]
Let me understand the full picture - how `deviceState` is created and the rig for tests.
✱ Grep "func.*deviceState\b" in internal/engine · 14 matches
✱ Grep "func newWiredRig" in internal/engine · 1 match
→ Read internal/engine/data.go [limit=30, offset=209]
→ Read internal/engine/fakestore_test.go [limit=60, offset=739]
✱ Grep "func.*GetDevice|type fakeStore|PurgeCompression" in internal/engine/fakestore_test.go · 8 matches
→ Read internal/engine/fakestore_test.go [limit=40, offset=245]
→ Read internal/engine/cache.go [limit=15, offset=28]
→ Read internal/store/store.go [limit=50, offset=30]
Let me check how migrations are managed in the store.
✱ Glob "**/migrations/**" in . · 15 matches
→ Read internal/store/store.go [limit=30, offset=90]
✱ Grep "PatchDeviceAttributes|attributes" in internal/engine · 1 match
→ Read internal/engine/fakestore_test.go [limit=100, offset=100]
✱ Grep "func.*Store.*Device" in internal/store/devices.go · 22 matches
→ Read internal/store/devices.go [offset=385, limit=80]
✱ Grep "func deviceMsg|func seedAlpha" in internal/engine · 3 matches
→ Read internal/engine/control_test.go [limit=50, offset=1]
→ Read internal/engine/fakestore_test.go [limit=60, offset=556]
Let me look at the fake store's `handle` method and the `pipelineRig` to understand how tests work.
✱ Grep "type pipelineRig|func.*handle\b" in internal/engine/fakestore_test.go · 3 matches
→ Read internal/engine/data_test.go [limit=50, offset=1]
✱ Grep "func.*rig.*handle|wiredRig|newWiredRig" in internal/engine/fakestore_test.go · 19 matches
Now let me understand the full cache loading path and what `GetDevice` returns.
→ Read internal/engine/cache.go [offset=440, limit=60]
→ Read internal/store/devices.go [offset=113, limit=40]
Let me look at the server data path that calls sendConsumerProperties and how the deviceView works.
→ Read internal/engine/serverdata.go [limit=100, offset=210]
✱ Grep "func.*retryStore" in internal/engine · 1 match
Let me check how `retryStore` works and if there are any other store methods for capabilities.
→ Read internal/engine/control.go [limit=50, offset=259]
→ Read internal/engine/serverdata.go [limit=100, offset=130]
→ Read migrations/000002_metadata.up.sql
Let me check the store migration for devices table and if there's any capability storage already.
✱ Grep "capability|capabilities|purge_compression|compression_format" in internal/store · 2 matches
→ Read migrations/migrations.go
Let me check what methods the Store interface already has and what's needed for the fix.
→ Read internal/engine/cache.go [limit=35, offset=28]

# Todos
[✓] Understand the bug: sendConsumerProperties defaults to zlib when peek returns nil, ignoring plaintext capability
[•] Add PurgeCompression to store.Device and persist capability in handleCapabilities
[ ] Load PurgeCompression in deviceCache.get so the store fallback works
[ ] Add test: evict cache entry for plaintext device, assert purge payload is uncompressed
[ ] Run gates: go vet ./... && go test ./...
[ ] Report in MULE.md format

✱ Grep "PatchDevice|SetPayload|store\." in internal/engine/cache.go · 100 matches
→ Read internal/engine/cache.go [limit=40, offset=60]
✱ Grep "PatchDeviceAttributes" in internal/engine · 0 matches
→ Read internal/engine/cache.go [limit=40, offset=28]
Let me look at the store attributes mechanism and how `deviceView` handles the offline fallback.
→ Read internal/store/devices.go [limit=40, offset=348]
→ Read internal/engine/control_test.go [limit=50, offset=135]
✱ Grep "PatchDeviceAttributes" in internal/engine/fakestore_test.go · 0 matches
→ Read internal/engine/fakestore_test.go [limit=40, offset=650]

# Todos
[✓] Understand the bug: sendConsumerProperties defaults to zlib when peek returns nil, ignoring plaintext capability
[•] Add PatchDeviceAttributes to engine Store interface + fakeStore
[ ] Persist purge compression in handleCapabilities via PatchDeviceAttributes
[ ] Load purge compression from attributes in deviceCache.get
[ ] Add store-read fallback in sendConsumerProperties when peek returns nil
[ ] Add test: evict cache entry for plaintext device, assert purge payload is uncompressed
[ ] Run gates: go vet ./... && go test ./...
[ ] Report in MULE.md format

Now I have a complete understanding. Let me plan the fix:

**The bug**: `sendConsumerProperties` (control.go:204) uses `devices.peek()` to check the device's purge compression. When `peek` returns nil (device not cached after eviction), it defaults to zlib. A device that declared plaintext can't decompress zlib.

**The fix**:
1. Persist the purge compression capability in the device's `attributes` JSONB column when `handleCapabilities` receives it
2. Load it from attributes in `deviceCache.get` when loading from store
3. In `sendConsumerProperties`, when `peek` returns nil, fall back to a direct store read (like `deviceView` does)

Let me implement this:
Let me check the `store.Device` Attributes type and understand how it's loaded.
→ Read internal/store/devices.go [limit=15, offset=37]
Now I have a complete understanding. Here's my plan:

**The bug**: When `devices.peek()` returns nil (device not cached after eviction), `sendConsumerProperties` defaults to zlib compression. A plaintext-capable device can't decompress zlib.

**The fix**: Persist the purge compression capability in the device's `attributes` JSONB column when received, and load it back from the store when the cache is cold.

1. Add `PatchDeviceAttributes` to the engine `Store` interface
2. Implement it in `fakeStore`  
3. Persist capability in `handleCapabilities`
4. Load from `Attributes` in `deviceCache.get`
5. Add store-read fallback in `sendConsumerProperties` for the `peek`-nil case
6. Write test

Let me start implementing:

← Edit internal/engine/cache.go
Index: /root/astrate-mule/internal/engine/cache.go
===================================================================
--- /root/astrate-mule/internal/engine/cache.go
+++ /root/astrate-mule/internal/engine/cache.go
@@ -64,8 +64,11 @@
 	ListTriggers(ctx context.Context, realmID int16) ([]store.Trigger, error)
 	// ListTriggerPolicies returns every delivery policy of a realm so the
 	// snapshot can attach compiled policies to triggers at load time.
 	ListTriggerPolicies(ctx context.Context, realmID int16) ([]store.TriggerPolicy, error)
+	// PatchDeviceAttributes merges key-value pairs into the device's
+	// user-visible attributes JSONB column. A nil value removes the key.
+	PatchDeviceAttributes(ctx context.Context, realmID int16, id deviceid.ID, patch map[string]*string) error
 }
 
 // fullReloadDebounce rate-limits snapshot self-heal reloads triggered by
 // realm-lookup misses, so a flood of messages for an unknown realm cannot

Now add the attribute key constant and modify `handleCapabilities` to persist the capability:
← Edit internal/engine/capabilities.go
Index: /root/astrate-mule/internal/engine/capabilities.go
===================================================================
--- /root/astrate-mule/internal/engine/capabilities.go
+++ /root/astrate-mule/internal/engine/capabilities.go
@@ -18,8 +18,14 @@
 	compressionZlib      = "zlib"
 	compressionPlaintext = "plaintext"
 )
 
+// attrPurgeCompression is the device-attributes key used to persist the
+// session-scoped purge compression capability across cache evictions.
+// Stored in the devices.attributes JSONB column; namespace-prefixed to
+// avoid collisions with user-set attributes.
+const attrPurgeCompression = "__purge_compression"
+
 // handleCapabilities parses the BSON capabilities payload published by the
 // device on `<realm>/<device_id>/capabilities` and stores the recognised
 // capabilities in the device's in-memory state. The only upstream capability
 // today is purge_properties_compression_format, which controls whether the

Now modify `handleCapabilities` to persist the capability to the store:
← Edit internal/engine/capabilities.go
Index: /root/astrate-mule/internal/engine/capabilities.go
===================================================================
--- /root/astrate-mule/internal/engine/capabilities.go
+++ /root/astrate-mule/internal/engine/capabilities.go
@@ -48,8 +48,18 @@
 		return
 	}
 	if v, ok := caps[capPurgePropertiesCompressionFormat]; ok {
 		dev.setPurgeCompression(v)
+		// Persist the capability so it survives cache evictions: when the
+		// device is offline and the server sends consumer/properties, the
+		// store fallback reads it back. The key is namespace-prefixed to
+		// avoid collisions with user-set attributes.
+		vCopy := v
+		patch := map[string]*string{attrPurgeCompression: &vCopy}
+		if err := e.st.PatchDeviceAttributes(ctx, realm.id, m.DeviceID, patch); err != nil {
+			e.log.Warn("persisting purge compression capability",
+				"realm", realm.name, "device", m.DeviceID.String(), "err", err)
+		}
 	}
 	m.Ack()
 }
 

Now modify `deviceCache.get` to load purge compression from the store's attributes:
← Edit internal/engine/cache.go
Index: /root/astrate-mule/internal/engine/cache.go
===================================================================
--- /root/astrate-mule/internal/engine/cache.go
+++ /root/astrate-mule/internal/engine/cache.go
@@ -467,8 +467,11 @@
 	realmID:       realmID,
 	introspection: dev.Introspection,
 	formatHint:    hint,
 }
+if v, ok := dev.Attributes[attrPurgeCompression]; ok {
+	st.purgeCompression = v
+}
 c.mu.Lock()
 c.m[key] = st
 c.mu.Unlock()
 return st, nil

Now modify `sendConsumerProperties` to fall back to a store read when `peek` returns nil:
→ Read internal/engine/control.go [limit=40, offset=190]

← Edit internal/engine/control.go
Index: /root/astrate-mule/internal/engine/control.go
===================================================================
--- /root/astrate-mule/internal/engine/control.go
+++ /root/astrate-mule/internal/engine/control.go
@@ -202,8 +202,17 @@
 }
 compression := compressionZlib
 if dev := e.devices.peek(realm.name, id); dev != nil {
 	compression = purgeCompressionFor(dev)
+} else {
+	// Device not cached (evicted on disconnect). Load from the store
+	// so a previously declared plaintext capability is respected
+	// rather than defaulting to zlib.
+	if d, err := e.st.GetDevice(ctx, realm.id, id); err == nil {
+		if v, ok := d.Attributes[attrPurgeCompression]; ok && v == compressionPlaintext {
+			compression = compressionPlaintext
+		}
+	}
 }
 var frame []byte
 switch compression {
 case compressionPlaintext:

Now implement `PatchDeviceAttributes` in the fakeStore:
→ Read internal/engine/fakestore_test.go [limit=30, offset=263]

← Edit internal/engine/fakestore_test.go
Index: /root/astrate-mule/internal/engine/fakestore_test.go
===================================================================
--- /root/astrate-mule/internal/engine/fakestore_test.go
+++ /root/astrate-mule/internal/engine/fakestore_test.go
@@ -272,8 +272,29 @@
 	}
 	return fmt.Errorf("%w: device %s", store.ErrNotFound, id)
 }
 
+func (f *fakeStore) PatchDeviceAttributes(_ context.Context, realmID int16, id deviceid.ID, patch map[string]*string) error {
+	f.mu.Lock()
+	defer f.mu.Unlock()
+	for _, dev := range f.devices {
+		if dev.ID == id && dev.RealmID == realmID {
+			if dev.Attributes == nil {
+				dev.Attributes = make(map[string]string)
+			}
+			for k, v := range patch {
+				if v == nil {
+					delete(dev.Attributes, k)
+				} else {
+					dev.Attributes[k] = *v
+				}
+			}
+			return nil
+		}
+	}
+	return fmt.Errorf("%w: device %s", store.ErrNotFound, id)
+}
+
 func (f *fakeStore) AppendDatastreams(_ context.Context, batch store.DatastreamBatch) error {
 	f.mu.Lock()
 	gate := f.appendGate
 	delay := f.appendDelay

Now let me write the test. Let me re-check `devAlpha` and constants used in tests, plus how eviction works.
✱ Grep "devAlpha\s*=|realmAlpha\s*=|realmAlphaID\s*=|ifaceServerProps\s*=" in internal/engine/fakestore_test.go · 4 matches
→ Read internal/engine/fakestore_test.go [limit=30, offset=552]
Now let me write the test. I'll add it after `TestZlibConsumerPropertiesDefault`. The test needs to:
1. Set a plaintext capability (which now persists it to the store)
2. Evict the device from the cache
3. Trigger a `sendConsumerProperties` (e.g. via emptyCache, but that requires the device cached again — let me think about how to trigger when cache is cold).

The `emptyCache` handler calls `deviceState` which reloads the device into cache on a miss. So emptyCache would repopulate the cache and lose the fallback path. I need a path where `sendConsumerProperties` is called while the cache is cold.

The paths that call `sendConsumerProperties` without loading the cache:
- `UnsetServerProperty` (serverdata.go:209) — calls `deviceView` + `sendConsumerProperties` directly, no cache load
- `OnDeviceConnected` (engine.go:268) — goroutine, but a connect event.

The cleanest way to test the store fallback is to call `UnsetServerProperty` (or `sendConsumerProperties` directly) after evicting the cache. Let me use the `UnsetServerProperty` path since it's an existing public method that leads to `sendConsumerProperties`.

Actually, simpler and most direct: call `e.sendConsumerProperties` directly after eviction. The test is in the engine package, so it can access the unexported method.

Let me set up:
1. Set a server property (so there's something in the purge list)
2. Set plaintext capability (which persists to store + cache)
3. Evict the device from cache
4. Call `sendConsumerProperties` directly
5. Assert payload is plaintext

Let me verify the fakePort's `publishedTo` and the publish topic.
→ Read internal/engine/fakestore_test.go [limit=60, offset=705]
Now let me write the test. Let me add it after `TestZlibConsumerPropertiesDefault` in capabilities_test.go:

← Edit internal/engine/capabilities_test.go
Index: /root/astrate-mule/internal/engine/capabilities_test.go
===================================================================
--- /root/astrate-mule/internal/engine/capabilities_test.go
+++ /root/astrate-mule/internal/engine/capabilities_test.go
@@ -230,20 +230,8 @@
 }
 
 // TestPurgeCompressionFor tests the fallback logic.
 func TestPurgeCompressionFor(t *testing.T) {
-	if got := purgeCompressionFor(nil); got != compressionZlib {
-		t.Errorf("purgeCompressionFor(nil) = %q, want zlib", got)
-	}
-	dev := &deviceState{}
-	if got := purgeCompressionFor(dev); got != compressionZlib {
-		t.Errorf("purgeCompressionFor(empty) = %q, want zlib", got)
-	}
-	dev.setPurgeCompression("plaintext")
-	if got := purgeCompressionFor(dev); got != "plaintext" {
-		t.Errorf("purgeCompressionFor(plaintext) = %q, want plaintext", got)
-	}
-}
 
 // mustBSON marshals a bson.M, failing the test on error.
 func mustBSON(t *testing.T, v bson.M) []byte {
 	t.Helper()

That edit was wrong - I removed the body of TestPurgeCompressionFor. Let me revert it.
← Edit internal/engine/capabilities_test.go
Index: /root/astrate-mule/internal/engine/capabilities_test.go
===================================================================
--- /root/astrate-mule/internal/engine/capabilities_test.go
+++ /root/astrate-mule/internal/engine/capabilities_test.go
@@ -230,8 +230,20 @@
 }
 
 // TestPurgeCompressionFor tests the fallback logic.
 func TestPurgeCompressionFor(t *testing.T) {
+	if got := purgeCompressionFor(nil); got != compressionZlib {
+		t.Errorf("purgeCompressionFor(nil) = %q, want zlib", got)
+	}
+	dev := &deviceState{}
+	if got := purgeCompressionFor(dev); got != compressionZlib {
+		t.Errorf("purgeCompressionFor(empty) = %q, want zlib", got)
+	}
+	dev.setPurgeCompression("plaintext")
+	if got := purgeCompressionFor(dev); got != "plaintext" {
+		t.Errorf("purgeCompressionFor(plaintext) = %q, want plaintext", got)
+	}
+}
 
 // mustBSON marshals a bson.M, failing the test on error.
 func mustBSON(t *testing.T, v bson.M) []byte {
 	t.Helper()

Good, restored. Now let me add my new test after `TestZlibConsumerPropertiesDefault`. Let me view the current state of the file to find where to insert:
→ Read internal/engine/capabilities_test.go [limit=30, offset=220]
I'll add the new test before `TestPurgeCompressionFor`. I'll need a context import. Let me add the test first:

← Edit internal/engine/capabilities_test.go
Index: /root/astrate-mule/internal/engine/capabilities_test.go
===================================================================
--- /root/astrate-mule/internal/engine/capabilities_test.go
+++ /root/astrate-mule/internal/engine/capabilities_test.go
@@ -228,8 +228,61 @@
 		t.Errorf("purge entries = %v, want [%s]", entries, want)
 	}
 }
 
+// TestPlaintextConsumerPropertiesAfterEviction: a plaintext-capable device
+// whose cache entry has been evicted (disconnect) must still receive an
+// uncompressed consumer/properties payload. This covers the store fallback
+// in sendConsumerProperties when devices.peek returns nil.
+func TestPlaintextConsumerPropertiesAfterEviction(t *testing.T) {
+	rig, fs, port := newWiredRig(t, Config{})
+	const serverIface = "com.astrate.test.ServerProperties"
+
+	fs.setProperty(store.Property{
+		RealmID: realmAlphaID, DeviceID: devAlpha, InterfaceID: ifaceServerProps,
+		Path: "/limits/maxConnections", Value: []byte("42"),
+	})
+
+	// Set plaintext capability: persisted to the device's attributes in the
+	// store and cached in the in-memory state.
+	ack := &ackCounter{}
+	rig.handle(deviceMsg("capabilities", "", 2, mustBSON(t, bson.M{
+		"purge_properties_compression_format": "plaintext",
+	}), ack))
+	if !ack.acked() {
+		t.Fatal("capabilities not acknowledged")
+	}
+
+	// Evict the device's cache entry, as a disconnect would.
+	rig.e.devices.evict(realmAlpha, devAlpha)
+
+	// sendConsumerProperties with a cold cache must fall back to the store
+	// and honour the persisted plaintext capability.
+	rs := rig.e.schemas.realmOrReload(context.Background(), realmAlpha)
+	if rs == nil {
+		t.Fatal("realm alpha not loaded")
+	}
+	if err := rig.e.sendConsumerProperties(context.Background(), rs, devAlpha); err != nil {
+		t.Fatalf("sendConsumerProperties: %v", err)
+	}
+
+	base := realmAlpha + "/" + devAlpha.String()
+	purges := port.publishedTo(base + "/control/consumer/properties")
+	if len(purges) != 1 {
+		t.Fatalf("consumer/properties messages: %d, want 1", len(purges))
+	}
+
+	// The payload must be raw plaintext, not a zlib frame.
+	payload := string(purges[0].payload)
+	want := serverIface + "/limits/maxConnections"
+	if payload != want {
+		t.Errorf("plaintext purge after eviction = %q, want %q", payload, want)
+	}
+	if _, err := inflateProperties(purges[0].payload); err == nil {
+		t.Error("plaintext purge after eviction is a valid zlib frame; should be raw text")
+	}
+}
+
 // TestPurgeCompressionFor tests the fallback logic.
 func TestPurgeCompressionFor(t *testing.T) {
 	if got := purgeCompressionFor(nil); got != compressionZlib {
 		t.Errorf("purgeCompressionFor(nil) = %q, want zlib", got)

Now I need to add the `context` import in capabilities_test.go:
← Edit internal/engine/capabilities_test.go
Index: /root/astrate-mule/internal/engine/capabilities_test.go
===================================================================
--- /root/astrate-mule/internal/engine/capabilities_test.go
+++ /root/astrate-mule/internal/engine/capabilities_test.go
@@ -1,7 +1,8 @@
 package engine
 
 import (
+	"context"
 	"testing"
 
 	promtest "github.com/prometheus/client_golang/prometheus/testutil"
 	"go.mongodb.org/mongo-driver/v2/bson"

Now let me check that the `realmOrReload` method exists with that signature:
✱ Grep "func.*realmOrReload" in internal/engine · 1 match
→ Read internal/engine/cache.go [limit=20, offset=165]
`realmOrReload` is on `schemaCache`, not the engine. Let me check how it's accessed on the Engine:
✱ Grep "realmOrReload|schemas" in internal/engine/serverdata.go · 29 matches
Good, `rig.e.schemas.realmOrReload` is correct. Now let me build and run the tests:

$ go build ./internal/engine/ 2>&1 | head -30
(no output)

Now let me run the specific test to verify:
$ go test ./internal/engine/ -run 'TestPlaintextConsumerPropertiesAfterEviction|TestPlaintextConsumerProperties|TestZlibConsumerPropertiesDefault|TestPurgeCompressionFor|TestHandleCapabilities' -v 2>&1 | tail -40
=== RUN   TestHandleCapabilities
--- PASS: TestHandleCapabilities (0.01s)
=== RUN   TestHandleCapabilitiesRejects
--- PASS: TestHandleCapabilitiesRejects (0.01s)
=== RUN   TestPlaintextConsumerProperties
--- PASS: TestPlaintextConsumerProperties (0.01s)
=== RUN   TestZlibConsumerPropertiesDefault
--- PASS: TestZlibConsumerPropertiesDefault (0.01s)
=== RUN   TestPlaintextConsumerPropertiesAfterEviction
--- PASS: TestPlaintextConsumerPropertiesAfterEviction (0.01s)
=== RUN   TestPurgeCompressionFor
--- PASS: TestPurgeCompressionFor (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/engine	0.093s

All tests pass. Now let me verify my test actually fails without the fix (per MULE.md requirement). Let me temporarily revert the fix in `sendConsumerProperties` and check:
$ git stash list 2>&1; git diff --stat 2>&1 | tail -20
.mule/cron.log                       | 431 +++++++++++++++++++++++++++++++++++
 internal/engine/cache.go             |   6 +
 internal/engine/capabilities.go      |  16 ++
 internal/engine/capabilities_test.go |  54 +++++
 internal/engine/control.go           |   9 +
 internal/engine/fakestore_test.go    |  21 ++
 6 files changed, 537 insertions(+)

Let me verify the test would fail without the fix by temporarily reverting the fallback logic in control.go:
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cp internal/engine/control.go /tmp/control.go.bak failed
Error: The user rejected permission to use this specific tool call.
