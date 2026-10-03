slug: device-empty-cache-received-trigger
verdict: blocked
at:  91d0191
ran: 2026-10-03T18:45:31Z on DietPi in 202s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/task.md
→ Read .mule/MULE.md

$ ls /root/astrate-mule && echo "---" && ls /root/astrate-mule/internal/engine
AGENTS.md
Dockerfile
LICENSE
Makefile
README.md
astrate
bench
cmd
deploy
docker-compose.yml
docs
examples
flow-container-echo
go.mod
go.sum
internal
migrations
opencode.json
pkg
site-dist
test
tools
---
batch.go
batch_test.go
bus_events_test.go
cache.go
cache_test.go
capabilities.go
capabilities_test.go
control.go
control_test.go
data.go
data_test.go
devicedata_test.go
engine.go
engine_e2e_test.go
engine_integration_test.go
fakestore_test.go
forward
introspection.go
introspection_test.go
mqttv1_parity_test.go
router.go
router_test.go
serverdata.go
serverdata_test.go
stream
testdata
topics.go
topics_test.go
triggers
triggers_change_integration_test.go
triggers_change_test.go

→ Read internal/engine/control.go

$ wc -l /root/astrate-mule/internal/engine/triggers/events.go /root/astrate-mule/internal/engine/triggers/match.go /root/astrate-mule/internal/engine/engine.go
312 /root/astrate-mule/internal/engine/triggers/events.go
  769 /root/astrate-mule/internal/engine/triggers/match.go
  456 /root/astrate-mule/internal/engine/engine.go
 1537 total

$ rg -n "EmptyCache|empty_cache" /root/astrate-mule --glob '!site-dist' -g '!*.json' | head -60
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml:1099:          enum: [incoming_data, value_change, value_change_applied, path_created, path_removed, value_stored, device_registered, device_connected, device_disconnected, device_empty_cache_received, device_error, incoming_introspection, interface_added, interface_removed, interface_minor_updated]
/root/astrate-mule/test/conformance/cpb/cpb_test.go:385:	t.Run("EmptyCacheResyncServerProperty", func(t *testing.T) {
/root/astrate-mule/internal/testutil/astartedevice.go:125:// EmptyCache publishes the control/emptyCache signal (docs/DESIGN.md §3.3),
/root/astrate-mule/internal/testutil/astartedevice.go:127:func (d *AstarteDevice) EmptyCache(t testing.TB) {
/root/astrate-mule/internal/engine/control.go:21:	controlEmptyCache         = "emptyCache"
/root/astrate-mule/internal/engine/control.go:42:	case controlEmptyCache:
/root/astrate-mule/internal/engine/control.go:43:		e.handleEmptyCache(ctx, m, realm)
/root/astrate-mule/internal/engine/control.go:52:// handleEmptyCache implements `control/emptyCache` (docs/DESIGN.md §3.3):
/root/astrate-mule/internal/engine/control.go:58:func (e *Engine) handleEmptyCache(ctx context.Context, m broker.InboundMessage, realm *realmSchema) {
/root/astrate-mule/internal/engine/triggers/errorname.go:21:	"empty_cache_error",
/root/astrate-mule/internal/engine/control_test.go:131:// TestEmptyCache: the device receives every server-owned property on its
/root/astrate-mule/internal/engine/control_test.go:134:func TestEmptyCache(t *testing.T) {
/root/astrate-mule/internal/engine/control_test.go:216:// TestEmptyCacheJSONHint: a device flipped to the JSON profile receives its
/root/astrate-mule/internal/engine/control_test.go:219:func TestEmptyCacheJSONHint(t *testing.T) {
/root/astrate-mule/internal/engine/triggers/match.go:10:// device_empty_cache_received, interface_minor_updated, and group-scoped
/root/astrate-mule/internal/engine/triggers/match.go:61:	// OnDeviceEmptyCacheReceived fires on control/emptyCache (accepted, not
/root/astrate-mule/internal/engine/triggers/match.go:63:	OnDeviceEmptyCacheReceived = "device_empty_cache_received"
/root/astrate-mule/internal/engine/triggers/match.go:123:	OnDeviceEmptyCacheReceived: false,
/root/astrate-mule/internal/engine/triggers/events_test.go:45:		{name: "device_error_no_metadata.json", body: NewDeviceErrorEvent("empty_cache_error", nil)},
/root/astrate-mule/internal/engine/triggers/events_test.go:119:	got, err = json.Marshal(envelope(NewDeviceErrorEvent("empty_cache_error", nil)))
/root/astrate-mule/internal/engine/engine_e2e_test.go:258:	t.Run("EmptyCacheResync", func(t *testing.T) { testE2EEmptyCache(t, env) })
/root/astrate-mule/internal/engine/engine_e2e_test.go:349:// testE2EEmptyCache: after a server-owned property is set, control/emptyCache
/root/astrate-mule/internal/engine/engine_e2e_test.go:353:func testE2EEmptyCache(t *testing.T, env *e2eEnv) {
/root/astrate-mule/internal/engine/engine_e2e_test.go:364:	dev.EmptyCache(t)

→ Read internal/engine/engine.go [offset=280, limit=140]
→ Read internal/engine/triggers/events.go

$ rg -n "Kind[A-Za-z]+ " /root/astrate-mule/internal/engine/stream/*.go | head -40
/root/astrate-mule/internal/engine/stream/bus.go:19:	// KindIncomingData is a committed device data operation (set or unset).
/root/astrate-mule/internal/engine/stream/bus.go:20:	KindIncomingData = "incoming_data"
/root/astrate-mule/internal/engine/stream/bus.go:21:	// KindDeviceRegistered is a device registration via the Pairing API.
/root/astrate-mule/internal/engine/stream/bus.go:22:	KindDeviceRegistered = "device_registered"
/root/astrate-mule/internal/engine/stream/bus.go:23:	// KindDeviceConnected is a device connection.
/root/astrate-mule/internal/engine/stream/bus.go:24:	KindDeviceConnected = "device_connected"
/root/astrate-mule/internal/engine/stream/bus.go:25:	// KindDeviceDisconnected is a device disconnection.
/root/astrate-mule/internal/engine/stream/bus.go:26:	KindDeviceDisconnected = "device_disconnected"
/root/astrate-mule/internal/engine/stream/bus.go:27:	// KindDeviceDeletionStarted is emitted immediately before a synchronous
/root/astrate-mule/internal/engine/stream/bus.go:29:	KindDeviceDeletionStarted = "device_deletion_started"
/root/astrate-mule/internal/engine/stream/bus.go:30:	// KindDeviceDeletionFinished is emitted immediately after a synchronous
/root/astrate-mule/internal/engine/stream/bus.go:32:	KindDeviceDeletionFinished = "device_deletion_finished"
/root/astrate-mule/internal/engine/stream/bus.go:33:	// KindDeviceError is a rejected device message (docs/DESIGN.md §2.6
/root/astrate-mule/internal/engine/stream/bus.go:36:	KindDeviceError = "device_error"
/root/astrate-mule/internal/engine/stream/bus_test.go:102:		if ev := recv(t, ch); ev.Kind != KindDeviceConnected {

→ Read internal/engine/triggers/match.go [limit=150]
→ Read internal/engine/stream/bus.go
→ Read internal/engine/triggers/match.go [offset=380, limit=70]
→ Read internal/engine/control_test.go [limit=240]
→ Read internal/engine/triggers/match.go [limit=130, offset=449]

$ rg -n "^func |^// Test" /root/astrate-mule/internal/engine/triggers_change_test.go | head -60; echo "---"; wc -l /root/astrate-mule/internal/engine/triggers_change_test.go
36:func (fw *fakeForwarder) Forward(_ context.Context, _, _ string, _ json.RawMessage, event []byte) error {
47:func (fw *fakeForwarder) collected() []capturedEvent {
54:func (fw *fakeForwarder) ofTrigger(name string) []capturedEvent {
65:func (fw *fakeForwarder) waitForCount(t *testing.T, name string, n int) {
75:func changeTriggerDef(name, on, iface string, major int, matchPath string) string {
95:func installChangeTriggers(t *testing.T, fs *fakeStore) {
105:// TestChangeTriggersLifecycle drives the whole lifecycle: first publish
110:func TestChangeTriggersLifecycle(t *testing.T) {
177:// TestChangeTriggersIntraBatch: two publishes of the same path inside one
179:func TestChangeTriggersIntraBatch(t *testing.T) {
199:func assertEvent(t *testing.T, ev capturedEvent, typ string, oldValue, newValue any) {
213:func assertValue(t *testing.T, ev capturedEvent, typ string, value any) {
225:func assertHasEvent(t *testing.T, events []capturedEvent, typ string, oldValue, newValue any) {
238:func assertStored(t *testing.T, ev capturedEvent, value any) {
250:func assertNoNew(t *testing.T, events []capturedEvent, limit int, what string) {
257:func lastOf(xs []capturedEvent) capturedEvent { return xs[len(xs)-1] }
260:func jsonAnyEqual(a, b any) bool {
278:func testNum(v any) (float64, bool) {
290:func jsonMarshalForTest(v any) string {
---
293 /root/astrate-mule/internal/engine/triggers_change_test.go

$ rg -n "emptyCache|EmptyCache" /root/astrate-mule/internal/engine/*_test.go
/root/astrate-mule/internal/engine/control_test.go:131:// TestEmptyCache: the device receives every server-owned property on its
/root/astrate-mule/internal/engine/control_test.go:134:func TestEmptyCache(t *testing.T) {
/root/astrate-mule/internal/engine/control_test.go:153:	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
/root/astrate-mule/internal/engine/control_test.go:155:		t.Fatal("emptyCache not acknowledged")
/root/astrate-mule/internal/engine/control_test.go:212:		t.Error("emptyCache did not arm the hint reset")
/root/astrate-mule/internal/engine/control_test.go:216:// TestEmptyCacheJSONHint: a device flipped to the JSON profile receives its
/root/astrate-mule/internal/engine/control_test.go:217:// emptyCache resends as JSON documents; the control frame stays zlib
/root/astrate-mule/internal/engine/control_test.go:219:func TestEmptyCacheJSONHint(t *testing.T) {
/root/astrate-mule/internal/engine/control_test.go:232:	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
/root/astrate-mule/internal/engine/control_test.go:234:		t.Fatal("emptyCache not acknowledged")
/root/astrate-mule/internal/engine/data_test.go:328:// emptyCache reset followed by a BSON payload.
/root/astrate-mule/internal/engine/data_test.go:362:		t.Error("hint flipped back to bson without an emptyCache reset")
/root/astrate-mule/internal/engine/data_test.go:404:	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
/root/astrate-mule/internal/engine/capabilities_test.go:173:	// emptyCache triggers consumer/properties send.
/root/astrate-mule/internal/engine/capabilities_test.go:175:	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
/root/astrate-mule/internal/engine/capabilities_test.go:177:		t.Fatal("emptyCache not acknowledged")
/root/astrate-mule/internal/engine/capabilities_test.go:208:	// emptyCache without any capability set → zlib (default).
/root/astrate-mule/internal/engine/capabilities_test.go:210:	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
/root/astrate-mule/internal/engine/capabilities_test.go:212:		t.Fatal("emptyCache not acknowledged")
/root/astrate-mule/internal/engine/engine_e2e_test.go:258:	t.Run("EmptyCacheResync", func(t *testing.T) { testE2EEmptyCache(t, env) })
/root/astrate-mule/internal/engine/engine_e2e_test.go:349:// testE2EEmptyCache: after a server-owned property is set, control/emptyCache
/root/astrate-mule/internal/engine/engine_e2e_test.go:353:func testE2EEmptyCache(t *testing.T, env *e2eEnv) {
/root/astrate-mule/internal/engine/engine_e2e_test.go:364:	dev.EmptyCache(t)
/root/astrate-mule/internal/engine/mqttv1_parity_test.go:93:	rig.handle(deviceMsg("control", "/emptyCache", 2, []byte("1"), ack))
/root/astrate-mule/internal/engine/mqttv1_parity_test.go:95:		t.Fatal("emptyCache not acknowledged")
/root/astrate-mule/internal/engine/topics_test.go:25:		{"control", "alpha/" + dev + "/control/emptyCache", "control/emptyCache", true},
/root/astrate-mule/internal/engine/topics_test.go:54:		{"control/emptyCache", kindControl, "emptyCache"},

→ Read internal/engine/triggers_change_test.go [limit=180]

$ rg -n "device_trigger" -A4 -B8 /root/astrate-mule/internal/engine/bus_events_test.go | head -80
(no output)

$ rg -ln "device_trigger" /root/astrate-mule/internal /root/astrate-mule/test | head -20
/root/astrate-mule/test/conformance/upstream/channels.json
/root/astrate-mule/test/conformance/upstream/channels.transcript.txt
/root/astrate-mule/test/conformance/upstream/README.md
/root/astrate-mule/test/conformance/upstream/recordchannels/main.go
/root/astrate-mule/internal/realm/dashboard_compat_test.go
/root/astrate-mule/internal/appengine/channels/room_test.go
/root/astrate-mule/internal/appengine/channels/ws_test.go
/root/astrate-mule/internal/realm/http_test.go
/root/astrate-mule/internal/appengine/channels/ws.go
/root/astrate-mule/internal/engine/cache_test.go
/root/astrate-mule/internal/engine/triggers/match.go
/root/astrate-mule/internal/engine/triggers/match_test.go
/root/astrate-mule/internal/engine/triggers/transient_test.go

$ rg -n "device_trigger" -B12 -A10 /root/astrate-mule/internal/engine/cache_test.go | head -100
274-
275-	policyDef := []byte(`{
276-		"name": "retry_policy",
277-		"maximum_capacity": 100,
278-		"error_handlers": [{"on": "any_error", "strategy": "retry"}],
279-		"retry_times": 3
280-	}`)
281-	fs.addPolicy(1, "retry_policy", policyDef)
282-
283-	fs.addTrigger(1, "with_policy", `{
284-		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
285-		"policy": "retry_policy",
286:		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
287-	}`)
288-	fs.addTrigger(1, "no_policy", `{
289-		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
290:		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
291-	}`)
292-	fs.addTrigger(1, "unknown_policy", `{
293-		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
294-		"policy": "nonexistent",
295:		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
296-	}`)
297-
298-	c := newSchemaCache(fs, discardLogger())
299-	if err := c.loadAll(ctx); err != nil {
300-		t.Fatalf("loadAll: %v", err)
301-	}
302-
303-	alpha := c.realm("alpha")
304-	if alpha == nil {
305-		t.Fatal("realm alpha not resolved")

$ rg -n "device_connected" /root/astrate-mule/internal/engine/*_test.go /root/astrate-mule/internal/engine/triggers/*_test.go | head -30
/root/astrate-mule/internal/engine/triggers/events_test.go:38:		{name: "device_connected.json", body: NewDeviceConnectedEvent("203.0.113.89")},
/root/astrate-mule/internal/engine/triggers/events_test.go:105:	// device_connected, per the upstream encoder.
/root/astrate-mule/internal/engine/triggers/events_test.go:113:	wantEvent := map[string]any{"type": "device_connected", "device_ip_address": "203.0.113.89"}
/root/astrate-mule/internal/engine/triggers/events_test.go:115:		t.Errorf("device_connected event: %v, want %v", gotMap["event"], wantEvent)
/root/astrate-mule/internal/engine/triggers/match_test.go:55:		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
/root/astrate-mule/internal/engine/triggers/match_test.go:314:		{name: "missing action", def: `{"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]}`,
/root/astrate-mule/internal/engine/triggers/match_test.go:317:			"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]}`,
/root/astrate-mule/internal/engine/triggers/match_test.go:377:			"type": "device_trigger", "on": "device_connected",
/root/astrate-mule/internal/engine/triggers/match_test.go:381:			"type": "device_trigger", "on": "device_connected", "device_id": "nope!"}]}`,
/root/astrate-mule/internal/engine/triggers/match_test.go:384:			"type": "device_trigger", "on": "device_connected", "device_id": "` +
/root/astrate-mule/internal/engine/triggers/match_test.go:391:			"type": "device_trigger", "on": "device_connected", "interface_name": "a.B"}]}`,
/root/astrate-mule/internal/engine/triggers/match_test.go:435:			{"type": "device_trigger", "on": "device_connected"},
/root/astrate-mule/internal/engine/triggers/match_test.go:479:		`{"type": "device_trigger", "on": "device_connected", "interface_major": 9}`))
/root/astrate-mule/internal/engine/triggers/match_test.go:532:		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
/root/astrate-mule/internal/engine/triggers/match_test.go:541:		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
/root/astrate-mule/internal/engine/triggers/match_test.go:553:		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
/root/astrate-mule/internal/engine/triggers/match_test.go:562:		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
/root/astrate-mule/internal/engine/triggers/match_test.go:573:		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
/root/astrate-mule/internal/engine/triggers/match_test.go:584:		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
/root/astrate-mule/internal/engine/cache_test.go:286:		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
/root/astrate-mule/internal/engine/cache_test.go:290:		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
/root/astrate-mule/internal/engine/cache_test.go:295:		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
/root/astrate-mule/internal/engine/triggers/transient_test.go:10:	raw := json.RawMessage(`{"type":"device_trigger","on":"device_connected"}`)
/root/astrate-mule/internal/engine/triggers/transient_test.go:60:	raw := json.RawMessage(`{"type":"device_trigger","on":"device_connected"}`)

$ rg -n "empty_cache|EmptyCache|Unsupported" /root/astrate-mule/internal/engine/triggers/match_test.go | head -40
284:// TestCompileAcceptsUnsupported: upstream-valid-but-unevaluated conditions
289:func TestCompileAcceptsUnsupported(t *testing.T) {
296:	if len(grouped.Unsupported) == 0 {
297:		t.Error("group-scoped trigger compiled without an Unsupported report")
481:	if len(minorUpdated.Unsupported) == 0 {
543:	if len(templated.Unsupported) != 0 {
544:		t.Errorf("mustache template should render, not be marked unsupported: %v", templated.Unsupported)
555:	if len(unknownTemplate.Unsupported) == 0 {

$ rg -rn "device_empty_cache_received|OnDeviceEmptyCacheReceived" /root/astrate-mule --glob '!site-dist/**' --glob '!docs/api/**' | head -30
/root/astrate-mule/internal/engine/triggers/match.go:// n, interface_minor_updated, and group-scoped
/root/astrate-mule/internal/engine/triggers/match.go:	// n fires on control/emptyCache (accepted, not
/root/astrate-mule/internal/engine/triggers/match.go:	n = "n"
/root/astrate-mule/internal/engine/triggers/match.go:	n: false,

$ rg -n "KindDeviceDisconnected|KindDeviceDeletionFinished|Kind[A-Za-z]+" /root/astrate-mule/internal/appengine/channels/*.go /root/astrate-mule/internal/realm/*.go 2>/dev/null | head -40
/root/astrate-mule/internal/appengine/channels/ws_test.go:576:	bus.ch <- stream.Event{Kind: stream.KindDeviceConnected, DeviceID: validDeviceIDA, IP: "1.2.3.4", Timestamp: now}
/root/astrate-mule/internal/appengine/channels/ws_test.go:672:	bus.ch <- stream.Event{Kind: stream.KindDeviceConnected, DeviceID: validDeviceIDA, IP: "1.2.3.4", Timestamp: now}
/root/astrate-mule/internal/appengine/channels/ws_test.go:753:				Kind: stream.KindDeviceConnected, DeviceID: validDeviceIDA,
/root/astrate-mule/internal/appengine/channels/room.go:258:			case stream.KindIncomingData:
/root/astrate-mule/internal/appengine/channels/room.go:267:			case stream.KindDeviceConnected, stream.KindDeviceDisconnected, stream.KindDeviceError:
/root/astrate-mule/internal/appengine/channels/room.go:281:			case stream.KindIncomingData:
/root/astrate-mule/internal/appengine/channels/room.go:283:			case stream.KindDeviceConnected:
/root/astrate-mule/internal/appengine/channels/room.go:285:			case stream.KindDeviceDisconnected:
/root/astrate-mule/internal/appengine/channels/room.go:287:			case stream.KindDeviceError:
/root/astrate-mule/internal/appengine/channels/room_test.go:84:	bus.ch <- stream.Event{Kind: stream.KindDeviceConnected, DeviceID: validDeviceIDA, IP: "1.2.3.4", Timestamp: now}
/root/astrate-mule/internal/appengine/channels/room_test.go:85:	bus.ch <- stream.Event{Kind: stream.KindDeviceDisconnected, DeviceID: validDeviceIDA, Timestamp: now}
/root/astrate-mule/internal/appengine/channels/room_test.go:86:	bus.ch <- stream.Event{Kind: stream.KindDeviceError, DeviceID: validDeviceIDA, ErrorName: "unknown_format", ErrorMetadata: map[string]string{"detail": "slow"}, Timestamp: now}
/root/astrate-mule/internal/appengine/channels/room_test.go:87:	bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: validDeviceIDA, Interface: "org.example.V1", Path: "/sensor/temp", Value: 42.0, InterfaceMajor: 1, Timestamp: now}
/root/astrate-mule/internal/appengine/channels/room_test.go:168:	bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: validDeviceIDB, Interface: "org.example.V1", Path: "/sensor/temp", Value: 42.0, InterfaceMajor: 1, Timestamp: now}
/root/astrate-mule/internal/appengine/channels/room_test.go:207:	bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: validDeviceIDB, Interface: "org.example.V1", Path: "/sensor/temp", Value: 1.0, InterfaceMajor: 1, Timestamp: now}
/root/astrate-mule/internal/appengine/channels/room_test.go:214:	bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: validDeviceIDA, Interface: "org.example.V1", Path: "/sensor/temp", Value: 2.0, InterfaceMajor: 1, Timestamp: now}
/root/astrate-mule/internal/appengine/channels/room_test.go:246:		bus.ch <- stream.Event{Kind: stream.KindIncomingData, DeviceID: "X", Interface: "org.example.V1", Path: "/a", Value: i, InterfaceMajor: 1, Timestamp: now}

$ rg -n "device_empty_cache_received" /root/astrate-mule --glob '!site-dist/**' | head -30
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml:1099:          enum: [incoming_data, value_change, value_change_applied, path_created, path_removed, value_stored, device_registered, device_connected, device_disconnected, device_empty_cache_received, device_error, incoming_introspection, interface_added, interface_removed, interface_minor_updated]
/root/astrate-mule/internal/engine/triggers/match.go:10:// device_empty_cache_received, interface_minor_updated, and group-scoped
/root/astrate-mule/internal/engine/triggers/match.go:63:	OnDeviceEmptyCacheReceived = "device_empty_cache_received"

→ Read internal/appengine/channels/room.go [limit=120, offset=230]

$ rg -n "func (e \*Engine) retryStore" -A40 /root/astrate-mule/internal/engine/*.go
(no output)

$ rg -n "type InboundMessage" -A30 /root/astrate-mule/internal/broker/*.go
/root/astrate-mule/internal/broker/intake.go:25:type InboundMessage struct {
/root/astrate-mule/internal/broker/intake.go-26-	// Realm is the device's realm name (first topic segment).
/root/astrate-mule/internal/broker/intake.go-27-	Realm string
/root/astrate-mule/internal/broker/intake.go-28-	// DeviceID is the publishing device, parsed from the connection's
/root/astrate-mule/internal/broker/intake.go-29-	// certificate CN.
/root/astrate-mule/internal/broker/intake.go-30-	DeviceID deviceid.ID
/root/astrate-mule/internal/broker/intake.go-31-	// Topic is the full topic the device published to,
/root/astrate-mule/internal/broker/intake.go-32-	// "<realm>/<device_id>[/<rest>]".
/root/astrate-mule/internal/broker/intake.go-33-	Topic string
/root/astrate-mule/internal/broker/intake.go-34-	// Payload is the raw message body (BSON, JSON, control bytes, or empty
/root/astrate-mule/internal/broker/intake.go-35-	// for property unset — classification happens in the engine).
/root/astrate-mule/internal/broker/intake.go-36-	Payload []byte
/root/astrate-mule/internal/broker/intake.go-37-	// QoS is the publish quality of service as received (0..2).
/root/astrate-mule/internal/broker/intake.go-38-	QoS byte
/root/astrate-mule/internal/broker/intake.go-39-	// ReceivedAt is the broker reception timestamp (used as the fallback
/root/astrate-mule/internal/broker/intake.go-40-	// datastream timestamp when the payload carries none).
/root/astrate-mule/internal/broker/intake.go-41-	ReceivedAt time.Time
/root/astrate-mule/internal/broker/intake.go-42-	// Ack releases the broker-held acknowledgment for this message. For
/root/astrate-mule/internal/broker/intake.go-43-	// QoS >= 1 the device's PUBACK (or PUBREC) is not sent until Ack is
/root/astrate-mule/internal/broker/intake.go-44-	// called — the engine calls it after the persistence batch commits
/root/astrate-mule/internal/broker/intake.go-45-	// (docs/DESIGN.md §5.3). Ack is never nil, is safe to call multiple
/root/astrate-mule/internal/broker/intake.go-46-	// times, and must eventually be called for every QoS >= 1 message the
/root/astrate-mule/internal/broker/intake.go-47-	// intake accepts, or the publishing device stalls (that stall is the
/root/astrate-mule/internal/broker/intake.go-48-	// designed backpressure path, §1.4).
/root/astrate-mule/internal/broker/intake.go-49-	Ack func()
/root/astrate-mule/internal/broker/intake.go-50-}
/root/astrate-mule/internal/broker/intake.go-51-
/root/astrate-mule/internal/broker/intake.go-52-// Intake consumes accepted device publishes. internal/engine implements it
/root/astrate-mule/internal/broker/intake.go-53-// (M6); tests use recorders. Submit may block — a full engine shard blocking
/root/astrate-mule/internal/broker/intake.go-54-// the broker's per-client read loop is the §1.4 backpressure contract.
/root/astrate-mule/internal/broker/intake.go-55-//

$ rg -n "retryStore" /root/astrate-mule/internal/engine/*.go | head
/root/astrate-mule/internal/engine/introspection.go:55:	if !e.retryStore(ctx, m, "introspection update", func() error {
/root/astrate-mule/internal/engine/introspection.go:129:// retryStore runs fn with shard-parking retries for transient store
/root/astrate-mule/internal/engine/introspection.go:135:func (e *Engine) retryStore(ctx context.Context, m broker.InboundMessage, what string, fn func() error) bool {
/root/astrate-mule/internal/engine/control.go:66:	if !e.retryStore(ctx, m, "emptyCache resync", func() error {
/root/astrate-mule/internal/engine/control.go:95:	if !e.retryStore(ctx, m, "producer/properties purge", func() error {

→ Read internal/engine/introspection.go [limit=70, offset=125]

$ rg -n "ListServerOwnedProperties" -B5 -A25 /root/astrate-mule/internal/engine/fakestore_test.go
430-		purged++
431-	}
432-	return purged, nil
433-}
434-
435:func (f *fakeStore) ListServerOwnedProperties(_ context.Context, realmID int16, deviceID deviceid.ID) ([]store.Property, error) {
436-	f.mu.Lock()
437-	defer f.mu.Unlock()
438-	var out []store.Property
439-	for ref, p := range f.properties[propsKey{realmID: realmID, id: deviceID}] {
440-		if f.ownershipOf(realmID, ref.InterfaceID) == interfaceschema.OwnershipServer {
441-			out = append(out, p)
442-		}
443-	}
444-	sort.Slice(out, func(i, j int) bool {
445-		if out[i].InterfaceID != out[j].InterfaceID {
446-			return out[i].InterfaceID < out[j].InterfaceID
447-		}
448-		return out[i].Path < out[j].Path
449-	})
450-	return out, nil
451-}
452-
453-func (f *fakeStore) ListTriggers(_ context.Context, realmID int16) ([]store.Trigger, error) {
454-	f.mu.Lock()
455-	defer f.mu.Unlock()
456-	out := make([]store.Trigger, len(f.triggersByRealm[realmID]))
457-	copy(out, f.triggersByRealm[realmID])
458-	return out, nil
459-}
460-

$ rg -n "errStore|storeErr|fail|ErrNotFound" /root/astrate-mule/internal/engine/fakestore_test.go | head -30
4:// failures, delays, and gates, plus fixture helpers reusing the M1 interface
276:	return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, id)
289:	return fmt.Errorf("%w: device %s", store.ErrNotFound, id)
373:		return nil, fmt.Errorf("%w: property %s", store.ErrNotFound, path)
386:		return nil, fmt.Errorf("%w: series %s", store.ErrNotFound, path)
412:	return nil, fmt.Errorf("%w: device %s", store.ErrNotFound, id)
574:// errTransient is the injected transient store failure.
575:var errTransient = errors.New("fakeStore: injected transient failure")
660:// enc encodes a `{v, t}` payload, failing the test on encoder errors.
730:// fakePort is an in-memory BrokerPort with programmable publish failures.

→ Read internal/engine/fakestore_test.go [limit=120, offset=560]
→ Read internal/engine/fakestore_test.go [offset=680]

$ rg -n "errTransient|pubErrs|inject" /root/astrate-mule/internal/engine/*_test.go | head -40
/root/astrate-mule/internal/engine/engine_integration_test.go:21:// errInjectedCommit is the injected datastream-commit failure.
/root/astrate-mule/internal/engine/engine_integration_test.go:22:var errInjectedCommit = errors.New("instrumented store: injected commit failure")
/root/astrate-mule/internal/engine/engine_integration_test.go:26:// database — the docs/ROADMAP.md §7.3 T2 "injected commit failure" knob.
/root/astrate-mule/internal/engine/engine_integration_test.go:333:// testAckAfterCommit: with an injected commit failure the acknowledgment is
/root/astrate-mule/internal/engine/engine_integration_test.go:346:	// Both injected failures must be consumed without an ack.
/root/astrate-mule/internal/engine/engine_integration_test.go:347:	waitFor(t, 10*time.Second, "injected commit failures consumed", func() bool {
/root/astrate-mule/internal/engine/batch_test.go:272:	fs.appendErrs = []error{errTransient, errTransient}
/root/astrate-mule/internal/engine/batch_test.go:344:	fs.appendErrs = []error{errTransient}
/root/astrate-mule/internal/engine/data_test.go:446:	fs.getDeviceErrs = []error{errTransient}
/root/astrate-mule/internal/engine/fakestore_test.go:324:		panic("fakeStore: injected AppendDatastreams panic")
/root/astrate-mule/internal/engine/fakestore_test.go:574:// errTransient is the injected transient store failure.
/root/astrate-mule/internal/engine/fakestore_test.go:575:var errTransient = errors.New("fakeStore: injected transient failure")
/root/astrate-mule/internal/engine/fakestore_test.go:735:	// pubErrs holds errors returned by the next Publish calls.
/root/astrate-mule/internal/engine/fakestore_test.go:736:	pubErrs []error
/root/astrate-mule/internal/engine/fakestore_test.go:742:	if len(p.pubErrs) > 0 {
/root/astrate-mule/internal/engine/fakestore_test.go:743:		err := p.pubErrs[0]
/root/astrate-mule/internal/engine/fakestore_test.go:744:		p.pubErrs = p.pubErrs[1:]
/root/astrate-mule/internal/engine/cache_test.go:346:	fs.getDeviceErrs = []error{errTransient}
/root/astrate-mule/internal/engine/introspection_test.go:167:	fs.updateIntroErrs = []error{errTransient}

→ Read internal/engine/fakestore_test.go [limit=120, offset=60]

$ rg -n "func isPermanentCommitError" -A25 /root/astrate-mule/internal/engine/*.go; rg -n "parkBackoffStart|parkBackoffCap" /root/astrate-mule/internal/engine/*.go | head
/root/astrate-mule/internal/engine/batch.go:206:func isPermanentCommitError(err error) bool {
/root/astrate-mule/internal/engine/batch.go-207-	var pgErr *pgconn.PgError
/root/astrate-mule/internal/engine/batch.go-208-	if errors.As(err, &pgErr) {
/root/astrate-mule/internal/engine/batch.go-209-		return strings.HasPrefix(pgErr.Code, "23") || strings.HasPrefix(pgErr.Code, "22")
/root/astrate-mule/internal/engine/batch.go-210-	}
/root/astrate-mule/internal/engine/batch.go-211-	return false
/root/astrate-mule/internal/engine/batch.go-212-}
/root/astrate-mule/internal/engine/batch.go-213-
/root/astrate-mule/internal/engine/batch.go-214-// ---------------------------------------------------------------------------
/root/astrate-mule/internal/engine/batch.go-215-// PersistOp → store row conversion (docs/DESIGN.md §2.4 typed columns, §2.3
/root/astrate-mule/internal/engine/batch.go-216-// property jsonb rendering).
/root/astrate-mule/internal/engine/batch.go-217-// ---------------------------------------------------------------------------
/root/astrate-mule/internal/engine/batch.go-218-
/root/astrate-mule/internal/engine/batch.go-219-// individualRow converts an OpIndividual into its typed-column row: scalars
/root/astrate-mule/internal/engine/batch.go-220-// land in their dedicated column, arrays in the pre-encoded JSON value_array.
/root/astrate-mule/internal/engine/batch.go-221-func individualRow(op *PersistOp) (*store.IndividualRow, error) {
/root/astrate-mule/internal/engine/batch.go-222-	row := &store.IndividualRow{
/root/astrate-mule/internal/engine/batch.go-223-		RealmID:     op.RealmID,
/root/astrate-mule/internal/engine/batch.go-224-		DeviceID:    op.DeviceID,
/root/astrate-mule/internal/engine/batch.go-225-		InterfaceID: op.Interface.ID,
/root/astrate-mule/internal/engine/batch.go-226-		EndpointID:  op.Mapping.EndpointID,
/root/astrate-mule/internal/engine/batch.go-227-		Path:        op.Path,
/root/astrate-mule/internal/engine/batch.go-228-		TS:          op.TS,
/root/astrate-mule/internal/engine/batch.go-229-		ReceptionTS: op.ReceptionTS,
/root/astrate-mule/internal/engine/batch.go-230-	}
/root/astrate-mule/internal/engine/batch.go-231-	vt := op.Mapping.ValueType
/root/astrate-mule/internal/engine/introspection.go:136:	backoff := parkBackoffStart
/root/astrate-mule/internal/engine/introspection.go:166:		backoff = min(backoff*2, parkBackoffCap)
/root/astrate-mule/internal/engine/batch.go:78:	backoff := parkBackoffStart
/root/astrate-mule/internal/engine/batch.go:119:		backoff = min(backoff*2, parkBackoffCap)
/root/astrate-mule/internal/engine/data.go:64:	parkBackoffStart = 100 * time.Millisecond
/root/astrate-mule/internal/engine/data.go:65:	parkBackoffCap   = 5 * time.Second
/root/astrate-mule/internal/engine/data.go:219:	backoff := parkBackoffStart
/root/astrate-mule/internal/engine/data.go:239:		backoff = min(backoff*2, parkBackoffCap)
/root/astrate-mule/internal/engine/introspection_test.go:173:	if elapsed := time.Since(start); elapsed < parkBackoffStart {
/root/astrate-mule/internal/engine/introspection_test.go:174:		t.Errorf("handler returned after %s, want >= %s of parking", elapsed, parkBackoffStart)

$ ls /root/astrate-mule/internal/engine/triggers/ /root/astrate-mule/internal/engine/triggers/testdata 2>/dev/null
/root/astrate-mule/internal/engine/triggers/:
actions.go
actions_test.go
change_test.go
errorname.go
errorname_test.go
events.go
events_test.go
match.go
match_test.go
policy.go
policy_test.go
testdata
transient.go
transient_test.go
validation.go

/root/astrate-mule/internal/engine/triggers/testdata:
device_connected.json
device_deletion_finished.json
device_deletion_started.json
device_disconnected.json
device_error.json
device_error_no_metadata.json
device_registered.json
incoming_data.json
incoming_data_unset.json
incoming_introspection.json
interface_added.json
interface_removed.json
path_created.json
path_removed.json
value_change.json
value_change_applied.json
value_change_created.json
value_stored.json

→ Read internal/engine/triggers/events_test.go

$ rg -n "func Golden" -A25 /root/astrate-mule/internal/testutil/*.go
/root/astrate-mule/internal/testutil/golden.go:20:func Golden(t testing.TB, name string, got []byte) {
/root/astrate-mule/internal/testutil/golden.go-21-	t.Helper()
/root/astrate-mule/internal/testutil/golden.go-22-
/root/astrate-mule/internal/testutil/golden.go-23-	path := filepath.Join("testdata", filepath.FromSlash(name))
/root/astrate-mule/internal/testutil/golden.go-24-	if *updateGolden {
/root/astrate-mule/internal/testutil/golden.go-25-		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
/root/astrate-mule/internal/testutil/golden.go-26-			t.Fatalf("golden: creating %s: %v", filepath.Dir(path), err)
/root/astrate-mule/internal/testutil/golden.go-27-			return
/root/astrate-mule/internal/testutil/golden.go-28-		}
/root/astrate-mule/internal/testutil/golden.go-29-		if err := os.WriteFile(path, got, 0o644); err != nil {
/root/astrate-mule/internal/testutil/golden.go-30-			t.Fatalf("golden: writing %s: %v", path, err)
/root/astrate-mule/internal/testutil/golden.go-31-			return
/root/astrate-mule/internal/testutil/golden.go-32-		}
/root/astrate-mule/internal/testutil/golden.go-33-		t.Logf("golden: wrote %s (%d bytes)", path, len(got))
/root/astrate-mule/internal/testutil/golden.go-34-		return
/root/astrate-mule/internal/testutil/golden.go-35-	}
/root/astrate-mule/internal/testutil/golden.go-36-
/root/astrate-mule/internal/testutil/golden.go-37-	want, err := os.ReadFile(path)
/root/astrate-mule/internal/testutil/golden.go-38-	if err != nil {
/root/astrate-mule/internal/testutil/golden.go-39-		t.Fatalf("golden: reading %s: %v (create it with `go test -update`)", path, err)
/root/astrate-mule/internal/testutil/golden.go-40-		return
/root/astrate-mule/internal/testutil/golden.go-41-	}
/root/astrate-mule/internal/testutil/golden.go-42-	if bytes.Equal(got, want) {
/root/astrate-mule/internal/testutil/golden.go-43-		return
/root/astrate-mule/internal/testutil/golden.go-44-	}
/root/astrate-mule/internal/testutil/golden.go-45-

$ ls ~/ /root 2>/dev/null; ls /root/astrate 2>/dev/null | head
/root:
Desktop
SurvivalManual
astrate
astrate-mule
astrate-survey
g-mind
go
govd
libheif
listello-buildtest
listello-mule
milestones-recipe.log
mule-crontab.bak
mule-crontab.paused
mule-failed-2026-09-04.tar.gz
mule-gate.log
pi.tail4db093.ts.net.crt
pi.tail4db093.ts.net.key

/root/:
Desktop
SurvivalManual
astrate
astrate-mule
astrate-survey
g-mind
go
govd
libheif
listello-buildtest
listello-mule
milestones-recipe.log
mule-crontab.bak
mule-crontab.paused
mule-failed-2026-09-04.tar.gz
mule-gate.log
pi.tail4db093.ts.net.crt
pi.tail4db093.ts.net.key
Dockerfile
LICENSE
Makefile
README.md
bench
cmd
deploy
docker-compose.yml
docs
go.mod

! permission requested: external_directory (/root/astrate/*); auto-rejecting
✗ cd /root/astrate && git log --oneline -3 2>/dev/null; git cat-file -t 958bb9fe 2>&1 | head -2 failed
Error: The user rejected permission to use this specific tool call.
