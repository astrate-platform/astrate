slug: broker-acl-coldstart-introspection-miss
verdict: done
at:  d201db4
ran: 2026-09-04T19:50:18Z on DietPi in 629s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
The task: fix the ACL miss path in `internal/broker/aclhook.go` so an interface introspected after connect (still-cold cache due to debounce) falls back to a synchronous store read, plus a T1 test in `broker_test.go`.

Let me read the relevant files.
→ Read internal/broker/aclhook.go
→ Read internal/broker/authhook.go
Now let me look at `server.go` processPublish and the test file.
✗ Read internal/broker/server.go failed [offset=830, limit=80]
Error: File not found: /root/astrate-mule/internal/broker/server.go
✱ Glob "internal/broker/*_test.go" 5 matches
✱ Grep "func \(s \*server\) processPublish|processPublish" in . · 4 matches
✱ Grep "processPublish" 4 matches
`processPublish` lives in mochi (a dependency). Let me read the review that spawned this, plus the test files.
→ Read .mule/reviews/broker-2026-09-04.md
→ Read internal/broker/aclhook_test.go
Now the key file — `broker_test.go` for the T1 harness.
→ Read internal/broker/broker_test.go
`processPublish` is in mochi (dependency) — the drop semantics are as described in the review. Now let me read `broker.go` and `intake.go` to see the internals the test will reach into.
→ Read internal/broker/broker.go
→ Read internal/broker/intake.go
Now let me check the testutil helpers, `discardLogger`, and any other callers of `refreshIfStale`.
✱ Grep "refreshIfStale" 8 matches
✱ Grep "func discardLogger|func MQTTConnect|func ServerTLSCert|func DeviceTLSConfig|func WaitToken|func DeviceCSR|func MQTTTryConnect" 11 matches
→ Read internal/testutil/mqttclient.go
→ Read internal/broker/sessionstore_test.go [limit=30]
I have the full picture. Let me implement.

**The bug:** in `OnACLCheck`'s miss path, `refreshIfStale` is debounced for the first second after connect (admit stamps `lastIntroLoad`), so a publish to an interface introspected after connect gets denied against the still-cold cache; for QoS0 mochi drops it silently.

**The fix:** `refreshIfStale` reports whether it actually reloaded; when the debounce skips it, fall back to a synchronous single-interface store read that warms the map.

← Edit internal/broker/authhook.go
Index: /root/astrate-mule/internal/broker/authhook.go
===================================================================
--- /root/astrate-mule/internal/broker/authhook.go
+++ /root/astrate-mule/internal/broker/authhook.go
@@ -180,23 +180,63 @@
 	return nil
 }
 
 // refreshIfStale runs refresh at most once per introspectionReloadDebounce.
-func (s *deviceSession) refreshIfStale(ctx context.Context, st Store, log *slog.Logger) {
+// It reports whether the reload actually ran: false means the debounce
+// skipped it (the cache is still cold) and the caller must fall back to a
+// synchronous store read for the missed interface.
+func (s *deviceSession) refreshIfStale(ctx context.Context, st Store, log *slog.Logger) bool {
 	s.mu.Lock()
 	stale := time.Since(s.lastIntroLoad) >= introspectionReloadDebounce
 	if stale {
 		s.lastIntroLoad = time.Now() // claim the slot before the slow reload
 	}
 	s.mu.Unlock()
 	if !stale {
-		return
+		return false
 	}
 	if err := s.refresh(ctx, st, log); err != nil {
 		log.Warn("introspection refresh failed", "client", s.identity.CN(), "error", err)
 	}
+	return true
 }
 
+// syncOwnershipOf resolves a single interface's ownership synchronously from
+// the store, bypassing the reload debounce — the cold-start fallback for an
+// interface introspected after connect. When the debounce skips refreshIfStale
+// the cache is still cold, and denying the packet against it would drop a
+// legitimate QoS0 publish (mochi's processPublish discards denied QoS0
+// silently). The resolved ownership is written back to the cache so the
+// per-interface store read happens once per introspection refresh cycle; an
+// interface the store does not know caches as denied, the same safe posture
+// as loadOwnership skipping unresolvable interfaces.
+func (s *deviceSession) syncOwnershipOf(ctx context.Context, st Store, log *slog.Logger, iface string) (interfaceschema.Ownership, bool) {
+	dev, err := st.GetDevice(ctx, s.realmID, s.identity.DeviceID)
+	if err != nil {
+		log.Debug("ACL cold-start interface resolution failed",
+			"client", s.identity.CN(), "interface", iface, "error", err)
+		return 0, false
+	}
+	ver, introspected := dev.Introspection[iface]
+	var own interfaceschema.Ownership
+	if introspected {
+		si, err := st.GetInterface(ctx, s.realmID, iface, ver.Major)
+		if err != nil {
+			log.Debug("ACL cold-start interface resolution failed",
+				"client", s.identity.CN(), "interface", iface, "major", ver.Major, "error", err)
+			return 0, false
+		}
+		own = si.Ownership
+	}
+	s.mu.Lock()
+	if s.ownership == nil {
+		s.ownership = make(map[string]interfaceschema.Ownership, 1)
+	}
+	s.ownership[iface] = own
+	s.mu.Unlock()
+	return own, introspected
+}
+
 // loadOwnership resolves each introspected interface's ownership. Interfaces
 // missing from the realm (not installed, or version mismatch) are skipped:
 // the ACL then denies their topics, which is the correct posture.
 func loadOwnership(ctx context.Context, st Store, realmID int16, intro map[string]store.InterfaceVersion, log *slog.Logger) map[string]interfaceschema.Ownership {

Now the miss path in `aclhook.go`:
← Edit internal/broker/aclhook.go
Index: /root/astrate-mule/internal/broker/aclhook.go
===================================================================
--- /root/astrate-mule/internal/broker/aclhook.go
+++ /root/astrate-mule/internal/broker/aclhook.go
@@ -186,17 +186,24 @@
 var allowed bool
 if sess := h.registry.get(cl.ID); sess != nil && sess.client == cl {
 	allowed = checkACL(sess.identity.BaseTopic(), topic, write, func(iface string) (interfaceschema.Ownership, bool) {
 		own, ok := sess.ownershipOf(iface)
-		if !ok {
-			// The interface may have been introspected after connect
-			// (the engine updates the row); reload, debounced.
-			ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
-			defer cancel()
-			sess.refreshIfStale(ctx, h.st, h.log)
-			own, ok = sess.ownershipOf(iface)
+		if ok {
+			return own, ok
 		}
-		return own, ok
+		// The interface may have been introspected after connect (the
+		// engine updates the row); reload, debounced.
+		ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
+		defer cancel()
+		if sess.refreshIfStale(ctx, h.st, h.log) {
+			return sess.ownershipOf(iface)
+		}
+		// Debounce skipped the reload: still inside the first
+		// introspectionReloadDebounce after connect, so the cache is
+		// cold. Resolve the unknown interface with one synchronous
+		// store read instead of denying the packet — mochi silently
+		// drops a denied QoS0 publish (processPublish).
+		return sess.syncOwnershipOf(ctx, h.st, h.log, iface)
 	})
 } else if !write {
 	// Delivery to a session-present but disconnected device (offline
 	// queue, retained messages): no live session, consult the store.

Now the T1 test in `broker_test.go`:
← Edit internal/broker/broker_test.go
Index: /root/astrate-mule/internal/broker/broker_test.go
===================================================================
--- /root/astrate-mule/internal/broker/broker_test.go
+++ /root/astrate-mule/internal/broker/broker_test.go
@@ -421,8 +421,82 @@
 	// The CN as client ID (the Go SDK convention) still connects.
 	publishAndCheckAttribution(t, identity.CN())
 }
 
+// TestBrokerACLColdStartIntrospectionMiss covers the connect-order race in
+// the ACL miss path: a device introspects (the engine persists the row) and
+// immediately publishes, before Broker.RefreshIntrospection runs. admit
+// stamps lastIntroLoad at connect, so refreshIfStale's debounce skips the
+// reload for the whole first second and the second cache check reads the
+// still-cold map — without the synchronous fallback the deny is dropped
+// silently by mochi for QoS0. lastIntroLoad is re-stamped just before the
+// publish so the debounce is deterministically pending regardless of machine
+// speed.
+func TestBrokerACLColdStartIntrospectionMiss(t *testing.T) {
+	ctx := context.Background()
+	st, realmCA, identity, _ := newFakeEnv(t)
+
+	// Connect with an empty-introspection device so the post-connect
+	// interface is genuinely unknown to the session cache.
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
+	// The engine persists the introspection the device just sent...
+	st.mu.Lock()
+	dev = st.devices[devKey(1, identity.DeviceID)]
+	dev.Introspection["com.ex.DeviceData"] = store.InterfaceVersion{Major: 1}
+	st.interfaces["1/com.ex.DeviceData/1"] = &store.StoredInterface{
+		Name: "com.ex.DeviceData", Major: 1, Ownership: interfaceschema.OwnershipDevice,
+	}
+	st.mu.Unlock()
+
+	// ...but Broker.RefreshIntrospection has not run yet. Re-stamp the
+	// reload mark so the miss deterministically hits the debounce and the
+	// synchronous fallback under test.
+	sess := b.registry.get(identity.CN())
+	if sess == nil {
+		t.Fatalf("no live session for %s", identity.CN())
+	}
+	sess.mu.Lock()
+	sess.lastIntroLoad = time.Now()
+	sess.mu.Unlock()
+
+	base := identity.BaseTopic()
+	token := client.Publish(base+"/com.ex.DeviceData/value", 0, false, []byte("v"))
+	testutil.WaitToken(t, token, 5*time.Second)
+	msg := intake.next(t, 5*time.Second)
+	if msg.Topic != base+"/com.ex.DeviceData/value" || msg.QoS != 0 {
+		t.Fatalf("unexpected intake message: %+v", msg)
+	}
+}
+
 func TestReloadRealmsPicksUpNewRealm(t *testing.T) {
 	ctx := context.Background()
 	st, _, _, _ := newFakeEnv(t)
 	intake := newRecorderIntake(true)

Now run the gate. The Pi can't run `-race` (39-bit VMA), so per MULE.md the gate is `go vet ./... && go test ./...`.
$ gofmt -l internal/broker/
(no output)

$ go vet ./internal/broker/ 2>&1
(no output)

$ go test ./internal/broker/ -run 'TestBrokerACLColdStartIntrospectionMiss|TestCheckACL|TestBrokerClientIDRemappedToCertCN|TestBrokerInMemorySmoke' -count=1 -v 2>&1 | tail -60
--- PASS: TestCheckACLMatrix/pub_other_realm (0.00s)
    --- PASS: TestCheckACLMatrix/pub_realm-prefix_trick (0.00s)
    --- PASS: TestCheckACLMatrix/pub_server-owned_interface (0.00s)
    --- PASS: TestCheckACLMatrix/pub_uninstalled_interface (0.00s)
    --- PASS: TestCheckACLMatrix/pub_interface_without_path (0.00s)
    --- PASS: TestCheckACLMatrix/pub_interface_with_empty_path (0.00s)
    --- PASS: TestCheckACLMatrix/pub_control_consumer_properties (0.00s)
    --- PASS: TestCheckACLMatrix/pub_control_unknown (0.00s)
    --- PASS: TestCheckACLMatrix/pub_control_keyAgreement (0.00s)
    --- PASS: TestCheckACLMatrix/pub_control_prefix_only (0.00s)
    --- PASS: TestCheckACLMatrix/pub_bare_realm (0.00s)
    --- PASS: TestCheckACLMatrix/pub_empty_topic (0.00s)
    --- PASS: TestCheckACLMatrix/pub_sys_topic (0.00s)
    --- PASS: TestCheckACLMatrix/pub_wildcard_hash (0.00s)
    --- PASS: TestCheckACLMatrix/pub_wildcard_plus_interface (0.00s)
    --- PASS: TestCheckACLMatrix/pub_oversized_topic (0.00s)
    --- PASS: TestCheckACLMatrix/sub_control_consumer_properties (0.00s)
    --- PASS: TestCheckACLMatrix/sub_superset_hash (0.00s)
    --- PASS: TestCheckACLMatrix/sub_server-owned_wildcard (0.00s)
    --- PASS: TestCheckACLMatrix/sub_server-owned_inner_wildcard (0.00s)
    --- PASS: TestCheckACLMatrix/sub_device-owned_interface (0.00s)
    --- PASS: TestCheckACLMatrix/sub_uninstalled_interface (0.00s)
    --- PASS: TestCheckACLMatrix/sub_control_hash (0.00s)
    --- PASS: TestCheckACLMatrix/sub_wildcard_interface_segment (0.00s)
    --- PASS: TestCheckACLMatrix/deliver_server-owned_concrete (0.00s)
    --- PASS: TestCheckACLMatrix/deliver_server-owned_nested (0.00s)
    --- PASS: TestCheckACLMatrix/deliver_device-owned_concrete (0.00s)
    --- PASS: TestCheckACLMatrix/deliver_uninstalled_concrete (0.00s)
    --- PASS: TestCheckACLMatrix/deliver_control_producer_properties (0.00s)
    --- PASS: TestCheckACLMatrix/sub_base_topic (0.00s)
    --- PASS: TestCheckACLMatrix/sub_other_device_hash (0.00s)
    --- PASS: TestCheckACLMatrix/sub_other_device_consumer_properties (0.00s)
    --- PASS: TestCheckACLMatrix/sub_other_realm_hash (0.00s)
    --- PASS: TestCheckACLMatrix/sub_global_hash (0.00s)
    --- PASS: TestCheckACLMatrix/sub_global_plus_pair (0.00s)
    --- PASS: TestCheckACLMatrix/sub_global_plus_hash (0.00s)
    --- PASS: TestCheckACLMatrix/sub_realm_hash (0.00s)
    --- PASS: TestCheckACLMatrix/sub_sys_topics (0.00s)
    --- PASS: TestCheckACLMatrix/sub_empty_topic (0.00s)
    --- PASS: TestCheckACLMatrix/sub_oversized_filter (0.00s)
=== RUN   TestCheckACLNoIntrospection
--- PASS: TestCheckACLNoIntrospection (0.00s)
=== RUN   TestBrokerInMemorySmoke
=== RUN   TestBrokerInMemorySmoke/DeferredAck
=== RUN   TestBrokerInMemorySmoke/QoS0Flows
=== RUN   TestBrokerInMemorySmoke/NoEchoToPublisher
=== RUN   TestBrokerInMemorySmoke/SubscribeDenied
=== RUN   TestBrokerInMemorySmoke/LifecycleEvents
--- PASS: TestBrokerInMemorySmoke (1.33s)
    --- PASS: TestBrokerInMemorySmoke/DeferredAck (0.50s)
    --- PASS: TestBrokerInMemorySmoke/QoS0Flows (0.00s)
    --- PASS: TestBrokerInMemorySmoke/NoEchoToPublisher (0.51s)
    --- PASS: TestBrokerInMemorySmoke/SubscribeDenied (0.00s)
    --- PASS: TestBrokerInMemorySmoke/LifecycleEvents (0.02s)
=== RUN   TestBrokerClientIDRemappedToCertCN
--- PASS: TestBrokerClientIDRemappedToCertCN (0.35s)
=== RUN   TestBrokerACLColdStartIntrospectionMiss
--- PASS: TestBrokerACLColdStartIntrospectionMiss (0.30s)
PASS
ok  	github.com/astrate-platform/astrate/internal/broker	2.091s

All pass. Now verify the test truly fails without the fix (the runner's check) — temporary revert of just the miss-path logic, run, restore.
← Edit internal/broker/aclhook.go
Index: /root/astrate-mule/internal/broker/aclhook.go
===================================================================
--- /root/astrate-mule/internal/broker/aclhook.go
+++ /root/astrate-mule/internal/broker/aclhook.go
@@ -196,14 +196,9 @@
 		defer cancel()
 		if sess.refreshIfStale(ctx, h.st, h.log) {
 			return sess.ownershipOf(iface)
 		}
-		// Debounce skipped the reload: still inside the first
-		// introspectionReloadDebounce after connect, so the cache is
-		// cold. Resolve the unknown interface with one synchronous
-		// store read instead of denying the packet — mochi silently
-		// drops a denied QoS0 publish (processPublish).
-		return sess.syncOwnershipOf(ctx, h.st, h.log, iface)
+		return sess.ownershipOf(iface)
 	})
 } else if !write {
 	// Delivery to a session-present but disconnected device (offline
 	// queue, retained messages): no live session, consult the store.

$ go test ./internal/broker/ -run 'TestBrokerACLColdStartIntrospectionMiss' -count=1 2>&1 | tail -15
--- FAIL: TestBrokerACLColdStartIntrospectionMiss (5.27s)
    broker_test.go:493: no message reached the intake within 5s
FAIL
FAIL	github.com/astrate-platform/astrate/internal/broker	5.367s
FAIL

Confirmed — without the fallback the test fails exactly as the bug describes (denied QoS0 dropped). Restoring the fix:
← Edit internal/broker/aclhook.go
Index: /root/astrate-mule/internal/broker/aclhook.go
===================================================================
--- /root/astrate-mule/internal/broker/aclhook.go
+++ /root/astrate-mule/internal/broker/aclhook.go
@@ -196,9 +196,14 @@
 		defer cancel()
 		if sess.refreshIfStale(ctx, h.st, h.log) {
 			return sess.ownershipOf(iface)
 		}
-		return sess.ownershipOf(iface)
+		// Debounce skipped the reload: still inside the first
+		// introspectionReloadDebounce after connect, so the cache is
+		// cold. Resolve the unknown interface with one synchronous
+		// store read instead of denying the packet — mochi silently
+		// drops a denied QoS0 publish (processPublish).
+		return sess.syncOwnershipOf(ctx, h.st, h.log, iface)
 	})
 } else if !write {
 	// Delivery to a session-present but disconnected device (offline
 	// queue, retained messages): no live session, consult the store.

$ gofmt -l internal/broker/ && go vet ./internal/broker/ 2>&1
(no output)

$ go test ./internal/broker/ -count=1 2>&1 | tail -5
ok  	github.com/astrate-platform/astrate/internal/broker	2.439s

Now the full gate across the repo:
$ go vet ./... 2>&1 | tail -5 && go test ./... -count=1 2>&1 | tail -40
ok  	github.com/astrate-platform/astrate/cmd/astrate	0.623s
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
ok  	github.com/astrate-platform/astrate/internal/appengine	0.100s
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	34.118s
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	2.559s
ok  	github.com/astrate-platform/astrate/internal/auth	2.620s
ok  	github.com/astrate-platform/astrate/internal/broker	2.711s
ok  	github.com/astrate-platform/astrate/internal/config	0.031s
ok  	github.com/astrate-platform/astrate/internal/engine	2.453s
ok  	github.com/astrate-platform/astrate/internal/engine/forward	0.074s
ok  	github.com/astrate-platform/astrate/internal/engine/stream	0.083s
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	2.595s
ok  	github.com/astrate-platform/astrate/internal/flow	0.636s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	1.210s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	0.140s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.574s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	0.046s
ok  	github.com/astrate-platform/astrate/internal/flowapi	0.253s
?   	github.com/astrate-platform/astrate/internal/housekeeping	[no test files]
ok  	github.com/astrate-platform/astrate/internal/httpx	0.056s
ok  	github.com/astrate-platform/astrate/internal/observability	0.127s
ok  	github.com/astrate-platform/astrate/internal/pairing	15.815s
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	35.672s
ok  	github.com/astrate-platform/astrate/internal/realm	0.055s
ok  	github.com/astrate-platform/astrate/internal/store	0.038s
?   	github.com/astrate-platform/astrate/internal/swagger	[no test files]
ok  	github.com/astrate-platform/astrate/internal/testutil	0.308s
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	0.164s
ok  	github.com/astrate-platform/astrate/pkg/deviceid	0.026s
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	0.162s
ok  	github.com/astrate-platform/astrate/pkg/payload	0.180s

Full suite green. Let me review the final diff.
$ git diff --stat && git diff internal/broker/authhook.go internal/broker/aclhook.go
internal/broker/aclhook.go     | 23 ++++++++-----
 internal/broker/authhook.go    | 44 +++++++++++++++++++++++--
 internal/broker/broker_test.go | 74 ++++++++++++++++++++++++++++++++++++++++++
 3 files changed, 131 insertions(+), 10 deletions(-)
diff --git a/internal/broker/aclhook.go b/internal/broker/aclhook.go
index 8129d59..bf9a323 100644
--- a/internal/broker/aclhook.go
+++ b/internal/broker/aclhook.go
@@ -187,15 +187,22 @@ func (h *aclHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
 	if sess := h.registry.get(cl.ID); sess != nil && sess.client == cl {
 		allowed = checkACL(sess.identity.BaseTopic(), topic, write, func(iface string) (interfaceschema.Ownership, bool) {
 			own, ok := sess.ownershipOf(iface)
-			if !ok {
-				// The interface may have been introspected after connect
-				// (the engine updates the row); reload, debounced.
-				ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
-				defer cancel()
-				sess.refreshIfStale(ctx, h.st, h.log)
-				own, ok = sess.ownershipOf(iface)
+			if ok {
+				return own, ok
 			}
-			return own, ok
+			// The interface may have been introspected after connect (the
+			// engine updates the row); reload, debounced.
+			ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
+			defer cancel()
+			if sess.refreshIfStale(ctx, h.st, h.log) {
+				return sess.ownershipOf(iface)
+			}
+			// Debounce skipped the reload: still inside the first
+			// introspectionReloadDebounce after connect, so the cache is
+			// cold. Resolve the unknown interface with one synchronous
+			// store read instead of denying the packet — mochi silently
+			// drops a denied QoS0 publish (processPublish).
+			return sess.syncOwnershipOf(ctx, h.st, h.log, iface)
 		})
 	} else if !write {
 		// Delivery to a session-present but disconnected device (offline
diff --git a/internal/broker/authhook.go b/internal/broker/authhook.go
index 650baa2..993a4ee 100644
--- a/internal/broker/authhook.go
+++ b/internal/broker/authhook.go
@@ -181,7 +181,10 @@ func (s *deviceSession) refresh(ctx context.Context, st Store, log *slog.Logger)
 }
 
 // refreshIfStale runs refresh at most once per introspectionReloadDebounce.
-func (s *deviceSession) refreshIfStale(ctx context.Context, st Store, log *slog.Logger) {
+// It reports whether the reload actually ran: false means the debounce
+// skipped it (the cache is still cold) and the caller must fall back to a
+// synchronous store read for the missed interface.
+func (s *deviceSession) refreshIfStale(ctx context.Context, st Store, log *slog.Logger) bool {
 	s.mu.Lock()
 	stale := time.Since(s.lastIntroLoad) >= introspectionReloadDebounce
 	if stale {
@@ -189,11 +192,48 @@ func (s *deviceSession) refreshIfStale(ctx context.Context, st Store, log *slog.
 	}
 	s.mu.Unlock()
 	if !stale {
-		return
+		return false
 	}
 	if err := s.refresh(ctx, st, log); err != nil {
 		log.Warn("introspection refresh failed", "client", s.identity.CN(), "error", err)
 	}
+	return true
+}
+
+// syncOwnershipOf resolves a single interface's ownership synchronously from
+// the store, bypassing the reload debounce — the cold-start fallback for an
+// interface introspected after connect. When the debounce skips refreshIfStale
+// the cache is still cold, and denying the packet against it would drop a
+// legitimate QoS0 publish (mochi's processPublish discards denied QoS0
+// silently). The resolved ownership is written back to the cache so the
+// per-interface store read happens once per introspection refresh cycle; an
+// interface the store does not know caches as denied, the same safe posture
+// as loadOwnership skipping unresolvable interfaces.
+func (s *deviceSession) syncOwnershipOf(ctx context.Context, st Store, log *slog.Logger, iface string) (interfaceschema.Ownership, bool) {
+	dev, err := st.GetDevice(ctx, s.realmID, s.identity.DeviceID)
+	if err != nil {
+		log.Debug("ACL cold-start interface resolution failed",
+			"client", s.identity.CN(), "interface", iface, "error", err)
+		return 0, false
+	}
+	ver, introspected := dev.Introspection[iface]
+	var own interfaceschema.Ownership
+	if introspected {
+		si, err := st.GetInterface(ctx, s.realmID, iface, ver.Major)
+		if err != nil {
+			log.Debug("ACL cold-start interface resolution failed",
+				"client", s.identity.CN(), "interface", iface, "major", ver.Major, "error", err)
+			return 0, false
+		}
+		own = si.Ownership
+	}
+	s.mu.Lock()
+	if s.ownership == nil {
+		s.ownership = make(map[string]interfaceschema.Ownership, 1)
+	}
+	s.ownership[iface] = own
+	s.mu.Unlock()
+	return own, introspected
 }
 
 // loadOwnership resolves each introspected interface's ownership. Interfaces

Let me run `gofmt -l .` repo-wide as MULE.md specifies, to be thorough.
$ gofmt -l . 2>/dev/null; echo "exit: $?"
exit: 0

Task complete. No git operations performed; work is in the working tree.

```
Done: ACL cold-start miss now falls back to a synchronous per-interface store read when the reload debounce skips refreshIfStale, so a QoS0 publish to an interface introspected after connect reaches the intake instead of being silently dropped
Files: internal/broker/authhook.go, internal/broker/aclhook.go, internal/broker/broker_test.go
Verified: go vet ./... -> pass; go test ./... -> pass; gofmt -l . -> clean; new TestBrokerACLColdStartIntrospectionMiss passes with the fix and fails without it (no message reached the intake within 5s)
Unsure: race coverage — the fix adds a mutex-guarded map write in syncOwnershipOf (same pattern as existing refresh/ownershipOf) but the race detector does not run here; also, refreshIfStale returns true even when the reload errors (caller re-checks the cold cache and denies, the safe posture), and unknown-interface negatives are cached as denied to preserve anti-flood
Follow-ups: none
```
