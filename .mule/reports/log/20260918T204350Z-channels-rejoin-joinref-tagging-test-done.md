slug: channels-rejoin-joinref-tagging-test
verdict: done
at:  d48e17f
ran: 2026-09-18T20:43:50Z on DietPi in 208s


> build · big-pickle

→ Read .mule/task.md
→ Read .mule/MULE.md
Task: add a test in ws_test.go proving a second `phx_join` updates the session's `joinRef` and subsequent `new_event` pushes carry it. Reading the named files.
→ Read internal/appengine/channels/ws.go
→ Read internal/appengine/channels/ws_test.go
Now let me understand the event flow — check how the room/bus deliver `new_event` (does it need a watch?).
✱ Grep "type fakeBus|func \(r \*Room\)|NewRegistry|KindDeviceConnected|watchFilter|func \(.*Room\) Watch|matches" in internal/appengine/channels · 21 matches
→ Read internal/appengine/channels/room.go [offset=63, limit=180]
→ Read internal/appengine/channels/room_test.go [offset=1, limit=45]
→ Read internal/appengine/channels/room.go [offset=242, limit=101]
Confirmed: an event reaches a member only through a matching watch. So the test must register a watch between the rejoin and the push, like `TestWireSession_NewEvent`. Let me check the `Frame` wire order.
✱ Grep "type Frame struct|JoinRef|func \(f Frame\)" in internal/appengine/channels · 17 matches
Wire format confirmed: `[join_ref, ref, topic, event, payload]`, element 0 is `join_ref`. The rejoin path exists (ws.go:284-291) and `pumpEvents` reads the fresh ref per event (ws.go:403). Now writing the test, mirroring `TestWireSession_NewEvent` (required watch) with a rejoin inserted.

← Edit internal/appengine/channels/ws_test.go
Index: /root/astrate-mule/internal/appengine/channels/ws_test.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/ws_test.go
+++ /root/astrate-mule/internal/appengine/channels/ws_test.go
@@ -624,8 +624,80 @@
 		t.Error("event is nil")
 	}
 }
 
+// TestWireSession_RejoinTagsPushWithFreshRef pins that a second phx_join on the
+// same topic updates the session's joinRef (ws.go:284-291) and that subsequent
+// new_event pushes carry the fresh ref: the comment there says "server pushes
+// must be tagged with the current one" but no test drove a rejoin, so that
+// behaviour rotted unasserted. Every other TestWireSession_* joins once.
+func TestWireSession_RejoinTagsPushWithFreshRef(t *testing.T) {
+	api, key, bus := setupTestAPIWithBus(t)
+	token := mintToken(t, key, jwt.MapClaims{"a_ch": []string{"JOIN::.*", "WATCH::.*"}})
+	mux := http.NewServeMux()
+	api.Mount(mux)
+	conn, cancel := dialSession(t, mux, token)
+	defer cancel()
+	defer conn.CloseNow()
+
+	ctx := context.Background()
+	topic := "rooms:" + testRealmName + ":dashboard_1"
+
+	// First join with ref "1".
+	joinTopic(t, conn, topic, "1")
+	if _, _, err := conn.Read(ctx); err != nil {
+		t.Fatalf("read first join reply: %v", err)
+	}
+
+	// Rejoin the same topic with a fresh ref "2".
+	joinTopic(t, conn, topic, "2")
+	if _, _, err := conn.Read(ctx); err != nil {
+		t.Fatalf("read rejoin reply: %v", err)
+	}
+	if api.reg.Rooms() != 1 {
+		t.Fatalf("Rooms() = %d, want 1 (a rejoin must not create a second room)", api.reg.Rooms())
+	}
+
+	// Register a watch so the pushed event reaches the session.
+	watchPayload := `{"name":"connectiontrigger-D","device_id":"` + validDeviceIDA + `","simple_trigger":{"type":"device_trigger","on":"device_connected","device_id":"` + validDeviceIDA + `"}}`
+	watchMsg := `["2","3","` + topic + `","watch",` + watchPayload + `]`
+	if err := conn.Write(ctx, websocket.MessageText, []byte(watchMsg)); err != nil {
+		t.Fatalf("write watch: %v", err)
+	}
+	if _, _, err := conn.Read(ctx); err != nil {
+		t.Fatalf("read watch reply: %v", err)
+	}
+
+	// Push an event into the bus.
+	now := time.Now().UTC().Truncate(time.Millisecond)
+	bus.ch <- stream.Event{Kind: stream.KindDeviceConnected, DeviceID: validDeviceIDA, IP: "1.2.3.4", Timestamp: now}
+
+	_, raw, err := conn.Read(ctx)
+	if err != nil {
+		t.Fatalf("read new_event: %v", err)
+	}
+
+	var arr []json.RawMessage
+	if err := json.Unmarshal(raw, &arr); err != nil {
+		t.Fatalf("unmarshal raw: %v", err)
+	}
+	if len(arr) != 5 {
+		t.Fatalf("got %d elements, want 5", len(arr))
+	}
+
+	// Element 0: join_ref must carry the fresh ref from the rejoin.
+	if string(arr[0]) != `"2"` {
+		t.Errorf("join_ref = %s, want \"2\" (the rejoin's ref)", string(arr[0]))
+	}
+
+	if string(arr[1]) != "null" {
+		t.Errorf("ref = %s, want null", string(arr[1]))
+	}
+	if string(arr[3]) != `"new_event"` {
+		t.Errorf("event = %s, want \"new_event\"", string(arr[3]))
+	}
+}
+
 // TestWireSession_PumpRacesLeave drives bus events at a joined room while the
 // session leaves it, so the pump goroutine and the read loop touch s.rooms and
 // joinRef concurrently. The specced tests all push events with the room
 // quiescent, which leaves the guarding of that map unexercised: without this,

$ gofmt -l internal/appengine/channels/ws_test.go
(no output)

$ go test ./internal/appengine/channels/ 2>&1
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	25.050s

Package passes and gofmt is clean. Let me empirically verify the test actually pins the rejoin (temporarily removing the update, expecting failure, then restoring).
← Edit internal/appengine/channels/ws.go
Index: /root/astrate-mule/internal/appengine/channels/ws.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/ws.go
+++ /root/astrate-mule/internal/appengine/channels/ws.go
@@ -283,9 +283,9 @@
 s.roomsMu.Lock()
 if j, ok := s.rooms[f.Topic]; ok {
 	// A rejoin on the same socket carries a fresh join_ref; server pushes
 	// must be tagged with the current one.
-	j.joinRef = f.Ref
+	_ = f.Ref
 	s.roomsMu.Unlock()
 	rep, _ := OK(f, nil)
 	s.writeFrame(rep)
 	return

$ go test ./internal/appengine/channels/ -run TestWireSession_RejoinTagsPushWithFreshRef 2>&1
# github.com/astrate-platform/astrate/internal/appengine/channels [github.com/astrate-platform/astrate/internal/appengine/channels.test]
internal/appengine/channels/ws.go:284:5: declared and not used: j
FAIL	github.com/astrate-platform/astrate/internal/appengine/channels [build failed]
FAIL

← Edit internal/appengine/channels/ws.go
Index: /root/astrate-mule/internal/appengine/channels/ws.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/ws.go
+++ /root/astrate-mule/internal/appengine/channels/ws.go
@@ -283,9 +283,9 @@
 s.roomsMu.Lock()
 if j, ok := s.rooms[f.Topic]; ok {
 	// A rejoin on the same socket carries a fresh join_ref; server pushes
 	// must be tagged with the current one.
-	_ = f.Ref
+	j.joinRef = j.joinRef
 	s.roomsMu.Unlock()
 	rep, _ := OK(f, nil)
 	s.writeFrame(rep)
 	return

$ go test ./internal/appengine/channels/ -run TestWireSession_RejoinTagsPushWithFreshRef 2>&1
--- FAIL: TestWireSession_RejoinTagsPushWithFreshRef (1.46s)
    ws_test.go:689: join_ref = "1", want "2" (the rejoin's ref)
FAIL
FAIL	github.com/astrate-platform/astrate/internal/appengine/channels	1.492s
FAIL

The test fails without the rejoin update — the pin is meaningful. Restoring the implementation.
← Edit internal/appengine/channels/ws.go
Index: /root/astrate-mule/internal/appengine/channels/ws.go
===================================================================
--- /root/astrate-mule/internal/appengine/channels/ws.go
+++ /root/astrate-mule/internal/appengine/channels/ws.go
@@ -283,9 +283,9 @@
 s.roomsMu.Lock()
 if j, ok := s.rooms[f.Topic]; ok {
 	// A rejoin on the same socket carries a fresh join_ref; server pushes
 	// must be tagged with the current one.
-	j.joinRef = j.joinRef
+	j.joinRef = f.Ref
 	s.roomsMu.Unlock()
 	rep, _ := OK(f, nil)
 	s.writeFrame(rep)
 	return

$ go vet ./internal/appengine/channels/ && go test ./internal/appengine/channels/ 2>&1 && gofmt -l internal/appengine/channels/ws_test.go
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)

Cached — forcing a fresh run to confirm the restored tree passes.
$ go test -count=1 ./internal/appengine/channels/ 2>&1
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	23.642s

$ git diff --stat && git diff -- internal/appengine/channels/ws.go | head -5
internal/appengine/channels/ws_test.go | 72 ++++++++++++++++++++++++++++++++++
 1 file changed, 72 insertions(+)

Done: Pinned in ws_test.go that a second `phx_join` on the same topic updates the session's `joinRef`, and that a subsequent `new_event` push carries the fresh ref.
Files: internal/appengine/channels/ws_test.go
Verified: `go vet ./internal/appengine/channels/` + `go test -count=1 ./internal/appengine/channels/` -> pass; `gofmt -l` -> clean; temporary back-out of the rejoin-ref update made the test fail (`join_ref = "1", want "2"`), confirming the pin is meaningful; implementation restored unchanged (diff is test-only).
Unsure: nothing — the push needs a registered watch to reach the session (dispatch filters by watch), so the test registers one before pushing, consistent with `TestWireSession_NewEvent`.
Follow-ups: none
