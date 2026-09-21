slug: broker-acl-coldstart-fallback-flood
verdict: done
at:  0e5078f
ran: 2026-09-21T19:40:20Z on DietPi in 218s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/broker/authhook.go
→ Read internal/broker/aclhook.go
✱ Grep "TestBrokerACLColdStartIntrospectionMiss" in . · 15 matches
→ Read internal/broker/broker_test.go [offset=400, limit=180]
✱ Grep "func newFakeEnv|func newRecorderIntake|func devKey|type fakeStore|struct.*Store" in internal/broker · 4 matches
→ Read internal/broker/broker_test.go [limit=300]
✱ Grep "GetDevice" in internal/broker · 18 matches
→ Read .mule/reviews/broker-2026-09-21.md
→ Read .mule/todo.md [offset=195, limit=30]
Now let me check the test utilities used by the existing test.
✱ Grep "func MQTTConnect|func MQTTTryConnect|func DiscardLogger|discardLogger" in internal/testutil · 2 matches
→ Read internal/testutil/mqttclient.go [offset=95, limit=80]
Now I'll implement the fix. First, the gating in `authhook.go`:

← Edit internal/broker/authhook.go
Index: /root/astrate-mule/internal/broker/authhook.go
===================================================================
--- /root/astrate-mule/internal/broker/authhook.go
+++ /root/astrate-mule/internal/broker/authhook.go
@@ -149,11 +149,12 @@
 	realmID  int16
 	client   *mqtt.Client // pointer identity guards against takeover races
 	remote   netip.Addr
 
-	mu            sync.Mutex
-	ownership     map[string]interfaceschema.Ownership
-	lastIntroLoad time.Time
+	mu               sync.Mutex
+	ownership        map[string]interfaceschema.Ownership
+	lastIntroLoad    time.Time
+	lastFallbackLoad time.Time
 }
 
 // ownershipOf reports the ownership of an introspected interface name.
 func (s *deviceSession) ownershipOf(iface string) (interfaceschema.Ownership, bool) {

← Edit internal/broker/authhook.go
Index: /root/astrate-mule/internal/broker/authhook.go
===================================================================
--- /root/astrate-mule/internal/broker/authhook.go
+++ /root/astrate-mule/internal/broker/authhook.go
@@ -200,18 +200,33 @@
 	}
 	return true
 }
 
-// syncOwnershipOf resolves a single interface's ownership synchronously from
-// the store, bypassing the reload debounce — the cold-start fallback for an
-// interface introspected after connect. When the debounce skips refreshIfStale
-// the cache is still cold, and denying the packet against it would drop a
-// legitimate QoS0 publish (mochi's processPublish discards denied QoS0
-// silently). The resolved ownership is written back to the cache so the
-// per-interface store read happens once per introspection refresh cycle; an
-// interface the store does not know caches as denied, the same safe posture
-// as loadOwnership skipping unresolvable interfaces.
+// syncOwnershipOf resolves an unknown interface's ownership synchronously from
+// the store — the cold-start fallback for an interface introspected after
+// connect, reached when the refreshIfStale debounce skips its reload and the
+// cache is still cold (denying the packet would drop a legitimate QoS0
+// publish: mochi's processPublish discards denied QoS0 silently). It is gated
+// to one synchronous store read per session per introspectionReloadDebounce
+// window — the same budget refreshIfStale's full reload pays — so an
+// adversarial topic flood of distinct unknown interface names cannot turn each
+// name into a full-device store read; the extra names within the window are
+// denied until the debounced reload repopulates the cache (≤1s later). The
+// claim is stamped before the slow read, matching the claim-before-reload
+// pattern in refreshIfStale, so concurrent misses beat on the counter, not the
+// database. The resolved ownership is written back to the cache so repeat
+// publishes to the same name stay cheap; an interface the store does not know
+// caches as denied, the same safe posture as loadOwnership skipping
+// unresolvable interfaces.
 func (s *deviceSession) syncOwnershipOf(ctx context.Context, st Store, log *slog.Logger, iface string) (interfaceschema.Ownership, bool) {
+	s.mu.Lock()
+	if time.Since(s.lastFallbackLoad) < introspectionReloadDebounce {
+		s.mu.Unlock()
+		return 0, false
+	}
+	s.lastFallbackLoad = time.Now() // claim the slot before the slow read
+	s.mu.Unlock()
+
 	dev, err := st.GetDevice(ctx, s.realmID, s.identity.DeviceID)
 	if err != nil {
 		log.Debug("ACL cold-start interface resolution failed",
 			"client", s.identity.CN(), "interface", iface, "error", err)

Now the GetDevice-counting fake:
← Edit internal/broker/broker_test.go
Index: /root/astrate-mule/internal/broker/broker_test.go
===================================================================
--- /root/astrate-mule/internal/broker/broker_test.go
+++ /root/astrate-mule/internal/broker/broker_test.go
@@ -104,14 +104,15 @@
 }
 
 // fakeStore is an in-memory broker.Store for T1 tests.
 type fakeStore struct {
-	mu         sync.Mutex
-	realms     []store.Realm
-	devices    map[string]*store.Device
-	interfaces map[string]*store.StoredInterface
-	connects   []netip.Addr
-	disconns   int
+	mu             sync.Mutex
+	realms         []store.Realm
+	devices        map[string]*store.Device
+	interfaces     map[string]*store.StoredInterface
+	connects       []netip.Addr
+	disconns       int
+	getDeviceCalls int
 }
 
 func newFakeStore() *fakeStore {
 	return &fakeStore{

← Edit internal/broker/broker_test.go
Index: /root/astrate-mule/internal/broker/broker_test.go
===================================================================
--- /root/astrate-mule/internal/broker/broker_test.go
+++ /root/astrate-mule/internal/broker/broker_test.go
@@ -131,16 +131,24 @@
 
 func (f *fakeStore) GetDevice(_ context.Context, realmID int16, id deviceid.ID) (*store.Device, error) {
 	f.mu.Lock()
 	defer f.mu.Unlock()
+	f.getDeviceCalls++
 	d, ok := f.devices[devKey(realmID, id)]
 	if !ok {
 		return nil, store.ErrNotFound
 	}
 	cp := *d
 	return &cp, nil
 }
 
+// getDeviceCount returns how many GetDevice calls the fake has served.
+func (f *fakeStore) getDeviceCount() int {
+	f.mu.Lock()
+	defer f.mu.Unlock()
+	return f.getDeviceCalls
+}
+
 func (f *fakeStore) GetInterface(_ context.Context, realmID int16, name string, major int) (*store.StoredInterface, error) {
 	f.mu.Lock()
 	defer f.mu.Unlock()
 	si, ok := f.interfaces[fmt.Sprintf("%d/%s/%d", realmID, name, major)]

Now add the new T1 test after `TestBrokerACLColdStartIntrospectionMiss`:
← Edit internal/broker/broker_test.go
Index: /root/astrate-mule/internal/broker/broker_test.go
===================================================================
--- /root/astrate-mule/internal/broker/broker_test.go
+++ /root/astrate-mule/internal/broker/broker_test.go
@@ -504,8 +504,90 @@
 		t.Fatalf("unexpected intake message: %+v", msg)
 	}
 }
 
+// TestBrokerACLColdStartFallbackFloodGated covers the anti-flood gate on the
+// ACL cold-start fallback: within one introspectionReloadDebounce window a
+// device publishing to K>=2 distinct names the store does not know must
+// trigger at most one synchronous fallback store read (the first name's), not
+// one per name — otherwise an adversarial QoS0 topic flood turns each distinct
+// name into a full-device GetDevice every second. lastIntroLoad is re-stamped
+// just before the publishes so the debounce is deterministically pending, and
+// the fake store counts GetDevice calls so the fallback budget is measured,
+// not assumed.
+func TestBrokerACLColdStartFallbackFloodGated(t *testing.T) {
+	ctx := context.Background()
+	st, realmCA, identity, _ := newFakeEnv(t)
+
+	// Connect with an empty-introspection device so every published name is
+	// genuinely unknown to the session cache.
+	st.mu.Lock()
+	dev := st.devices[devKey(1, identity.DeviceID)]
+	dev.Introspection = map[string]store.InterfaceVersion{}
+	st.mu.Unlock()
+
+	intake := newRecorderIntake(true)
+	serverCert, roots := testutil.ServerTLSCert(t)
+
+	b, err := New(ctx, Config{
+		TLSAddr:          "127.0.0.1:0",
+		ServerTLSCert:    serverCert,
+		SessionStorePath: t.TempDir() + "/sessions.db",
+		Logger:           discardLogger(),
+	}, st, intake, nil)
+	if err != nil {
+		t.Fatalf("New: %v", err)
+	}
+	if err := b.Start(); err != nil {
+		t.Fatalf("Start: %v", err)
+	}
+	t.Cleanup(func() { _ = b.Close() })
+
+	devKeyPriv, csrPEM := testutil.DeviceCSR(t)
+	certPEM, _, _, err := realmCA.SignCSR(csrPEM, identity.Realm, identity.DeviceID.String(), time.Hour)
+	if err != nil {
+		t.Fatalf("issuing device certificate: %v", err)
+	}
+	tlsCfg := testutil.DeviceTLSConfig(t, certPEM, devKeyPriv, roots)
+	client, _ := testutil.MQTTConnect(t, "ssl://"+b.TLSAddr(), identity.CN(), true, tlsCfg)
+
+	// The auth handshake accounts for one GetDevice; every read after this
+	// point is the fallback under test.
+	authReads := st.getDeviceCount()
+	if authReads != 1 {
+		t.Fatalf("auth GetDevice reads = %d, want 1", authReads)
+	}
+
+	// Re-stamp the reload mark so every miss deterministically hits the
+	// debounce and the gated synchronous fallback under test.
+	sess := b.registry.get(identity.CN())
+	if sess == nil {
+		t.Fatalf("no live session for %s", identity.CN())
+	}
+	sess.mu.Lock()
+	sess.lastIntroLoad = time.Now()
+	sess.mu.Unlock()
+
+	base := identity.BaseTopic()
+	const distinctNames = 4 // K >= 2
+	for i := 0; i < distinctNames; i++ {
+		token := client.Publish(base+fmt.Sprintf("/com.ex.Bogus%d/value", i), 0, false, []byte("v"))
+		testutil.WaitToken(t, token, 5*time.Second)
+	}
+
+	// The first distinct name must take its one fallback read...
+	waitFor(t, 5*time.Second, "first fallback GetDevice", func() bool {
+		return st.getDeviceCount() >= authReads+1
+	})
+	// ...and the remaining names in the window must not buy more: give their
+	// publishes time to land, then assert the cap (still inside the 1s debounce,
+	// so an un-gated implementation would already have served K reads).
+	time.Sleep(400 * time.Millisecond)
+	if got := st.getDeviceCount(); got > authReads+1 {
+		t.Fatalf("fallback GetDevice reads = %d, want at most %d (one per introspectionReloadDebounce window); got %d distinct names", got, authReads+1, distinctNames)
+	}
+}
+
 func TestReloadRealmsPicksUpNewRealm(t *testing.T) {
 	ctx := context.Background()
 	st, _, _, _ := newFakeEnv(t)
 	intake := newRecorderIntake(true)

Now run the gate — gofmt, vet, and the broker tests:
$ gofmt -l . && go vet ./internal/broker/ && go test ./internal/broker/ -run 'TestBrokerACLColdStart' -count=1 -v 2>&1 | tail -30
=== RUN   TestBrokerACLColdStartIntrospectionMiss
--- PASS: TestBrokerACLColdStartIntrospectionMiss (0.29s)
=== RUN   TestBrokerACLColdStartFallbackFloodGated
--- PASS: TestBrokerACLColdStartFallbackFloodGated (0.68s)
PASS
ok  	github.com/astrate-platform/astrate/internal/broker	1.094s

Both pass, gofmt clean. Now verify the new test actually fails without the fix (temporarily removing the gate):
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cp internal/broker/authhook.go /tmp/authhook.go.bak failed
Error: The user rejected permission to use this specific tool call.
