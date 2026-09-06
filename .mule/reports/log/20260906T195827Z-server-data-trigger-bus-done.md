slug: server-data-trigger-bus
verdict: done
at:  e46ce11
ran: 2026-09-06T19:58:27Z on DietPi in 827s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/engine/serverdata.go
→ Read internal/engine/batch.go
✱ Grep "afterCommit|fireCommitted|fireData|finalize" in internal/engine · 27 matches
✱ Glob "internal/engine/*.go" 27 matches
→ Read internal/engine/engine.go [offset=60, limit=200]
→ Read internal/engine/data.go
→ Read internal/engine/serverdata_test.go
→ Read internal/engine/devicedata_test.go
→ Read internal/engine/bus_events_test.go
→ Read internal/engine/fakestore_test.go
→ Read internal/engine/triggers_change_test.go [limit=120]
→ Read internal/engine/engine.go [offset=260, limit=200]
✱ Grep "unreadable row|lookupFailed|internalErrors|latestLookups|propLookups|Tracked|PrevExists|pendingPrev" in internal/engine · 11 matches
→ Read internal/engine/triggers_change_test.go [offset=120]
→ Read internal/engine/engine.go [limit=60]
✱ Grep "func New|func newEngine|newPipelineRig|flushShard|func.*pipelineRig" in internal/engine · 55 matches
→ Read internal/engine/stream/bus.go
→ Read internal/engine/triggers/match.go
✱ Grep "PublishServerValue|PublishDeviceValue|UnsetServerProperty|publishAsOwner" 81 matches
→ Read .mule/reviews/engine-2026-09-06.md
→ Read internal/engine/triggers/actions.go [offset=200, limit=120]

# Todos
[•] Extract lookupPrev + add capturePrevious in data.go
[ ] Route publishAsOwner through capturePrevious + fireData in serverdata.go
[ ] Route UnsetServerProperty through capturePrevious + fireData
[ ] Add failing bus+trigger tests in serverdata_test.go and devicedata_test.go
[ ] Run gate: go vet, go test, gofmt

Now I have full context. Let me set up a todo list and start implementing.
← Edit internal/engine/data.go
Index: /root/astrate-mule/internal/engine/data.go
===================================================================
--- /root/astrate-mule/internal/engine/data.go
+++ /root/astrate-mule/internal/engine/data.go
@@ -395,37 +395,12 @@
 )
 if snap, ok := sh.batch.pendingPrev[key]; ok {
 	prev, prevFound = snap.json, snap.exists
 } else {
-	switch op.Kind {
-	case OpPropertySet, OpPropertyUnset:
-		p, err := e.st.GetProperty(ctx, op.RealmID, op.DeviceID, op.Interface.ID, op.Path)
-		switch {
-		case err == nil:
-			prev, prevFound = p.Value, true
-		case errors.Is(err, store.ErrNotFound):
-			// First-ever set (or re-unset): nothing there before.
-		default:
-			e.lookupFailed(op, err)
-			return
-		}
-	case OpIndividual:
-		r, err := e.st.LatestIndividual(ctx, op.RealmID, op.DeviceID, op.Interface.ID, op.Path)
-		switch {
-		case err == nil:
-			prev, prevFound, err = prevIndividualJSON(r)
-			if err != nil {
-				e.met.internalErrors.Inc()
-				e.log.Error("previous-value snapshot skipped: unreadable row",
-					"realm", op.Realm, "device", op.DeviceID.String(), "path", op.Path, "err", err)
-				return
-			}
-		case errors.Is(err, store.ErrNotFound):
-			// First sample on this path.
-		default:
-			e.lookupFailed(op, err)
-			return
-		}
+	prev, prevFound, err = e.lookupPrev(ctx, op)
+	if err != nil {
+		e.lookupFailed(op, err)
+		return
 	}
 }
 
 op.Tracked = true

→ Read internal/engine/data.go [limit=60, offset=388]
Now add `lookupPrev` and `capturePrevious` after `trackPrevious`:

← Edit internal/engine/data.go
Index: /root/astrate-mule/internal/engine/data.go
===================================================================
--- /root/astrate-mule/internal/engine/data.go
+++ /root/astrate-mule/internal/engine/data.go
@@ -412,8 +412,76 @@
 	}
 	sh.batch.pendingPrev[key] = prevSnap{exists: op.Kind != OpPropertyUnset, json: newJSON}
 }
 
+// lookupPrev reads the pre-write previous-value state of one endpoint from
+// the store: the property row for property ops, the newest sample for
+// individual datastreams. A missing row is not an error — found is false.
+func (e *Engine) lookupPrev(ctx context.Context, op *PersistOp) (prev []byte, found bool, err error) {
+	switch op.Kind {
+	case OpPropertySet, OpPropertyUnset:
+		p, err := e.st.GetProperty(ctx, op.RealmID, op.DeviceID, op.Interface.ID, op.Path)
+		switch {
+		case err == nil:
+			return p.Value, true, nil
+		case errors.Is(err, store.ErrNotFound):
+			return nil, false, nil
+		default:
+			return nil, false, err
+		}
+	case OpIndividual:
+		r, err := e.st.LatestIndividual(ctx, op.RealmID, op.DeviceID, op.Interface.ID, op.Path)
+		switch {
+		case err == nil:
+			return prevIndividualJSON(r)
+		case errors.Is(err, store.ErrNotFound):
+			return nil, false, nil
+		default:
+			return nil, false, err
+		}
+	default:
+		return nil, false, nil
+	}
+}
+
+// capturePrevious fills the previous-value snapshot of an op that persists
+// immediately — the server-/device-owned publish paths (PublishServerValue,
+// PublishDeviceValue, UnsetServerProperty) — with the same gate and store
+// lookups as trackPrevious, minus the shard micro-batch chain: each such
+// publish commits synchronously, so the store alone is always the correct
+// "previous". Change-derived conditions (value_change*, path_created/removed)
+// fire only when a trigger actually watches the endpoint, exactly like the
+// broker ingress path.
+func (e *Engine) capturePrevious(ctx context.Context, rs *realmSchema, op *PersistOp) {
+	if op.Kind == OpObject {
+		return // object-aggregated publishes stay outside v1's change scope
+	}
+	gate := triggers.DataEvent{
+		DeviceID:  op.DeviceID.String(),
+		Interface: op.Interface.Name,
+		Major:     op.Interface.Major,
+		Path:      op.Path,
+	}
+	watched := false
+	for _, tr := range rs.triggers {
+		if tr.TracksChanges(gate) {
+			watched = true
+			break
+		}
+	}
+	if !watched {
+		return
+	}
+	prev, prevFound, err := e.lookupPrev(ctx, op)
+	if err != nil {
+		e.lookupFailed(op, err)
+		return
+	}
+	op.Tracked = true
+	op.PrevExists = prevFound
+	op.Prev = prev
+}
+
 // lookupFailed records a degraded previous-value lookup.
 func (e *Engine) lookupFailed(op *PersistOp, err error) {
 	e.met.internalErrors.Inc()
 	e.log.Warn("previous-value lookup failed; change events skipped",

← Edit internal/engine/serverdata.go
Index: /root/astrate-mule/internal/engine/serverdata.go
===================================================================
--- /root/astrate-mule/internal/engine/serverdata.go
+++ /root/astrate-mule/internal/engine/serverdata.go
@@ -165,32 +165,38 @@
 	switch {
 	case ci.Type == interfaceschema.Properties:
 		op.Kind = OpPropertySet
 		out.retain = true // properties stay retained (§3.4)
+	case ci.Aggregation == interfaceschema.AggregationObject:
+		op.Kind = OpObject
+		out.wireTS = &effTS
+		out.qos = mapping.Reliability
+		out.expiry = mapping.Expiry
+	default:
+		op.Kind = OpIndividual
+		out.wireTS = &effTS
+		out.qos = mapping.Reliability
+		out.expiry = mapping.Expiry
+	}
+	e.capturePrevious(ctx, rs, &op)
+	switch {
+	case ci.Type == interfaceschema.Properties:
 		row, err := propertyRow(&op)
 		if err != nil {
 			return nil, err
 		}
 		if err := e.st.UpsertProperty(ctx, *row); err != nil {
 			return nil, err
 		}
 	case ci.Aggregation == interfaceschema.AggregationObject:
-		op.Kind = OpObject
-		out.wireTS = &effTS
-		out.qos = mapping.Reliability
-		out.expiry = mapping.Expiry
 		row, err := objectRow(&op)
 		if err != nil {
 			return nil, err
 		}
 		if err := e.st.AppendDatastreams(ctx, store.DatastreamBatch{Objects: []store.ObjectRow{*row}}); err != nil {
 			return nil, err
 		}
 	default:
-		op.Kind = OpIndividual
-		out.wireTS = &effTS
-		out.qos = mapping.Reliability
-		out.expiry = mapping.Expiry
 		row, err := individualRow(&op)
 		if err != nil {
 			return nil, err
 		}
@@ -198,8 +204,9 @@
 			return nil, err
 		}
 	}
 	e.met.persistOps.WithLabelValues(op.Kind.String()).Inc()
+	e.fireData(rs, &op)
 	return out, nil
 }
 
 // UnsetServerProperty deletes a server-owned property (AppEngine DELETE,

← Edit internal/engine/serverdata.go
Index: /root/astrate-mule/internal/engine/serverdata.go
===================================================================
--- /root/astrate-mule/internal/engine/serverdata.go
+++ /root/astrate-mule/internal/engine/serverdata.go
@@ -88,10 +88,14 @@
 // publishAsOwner is the shared validate-and-persist core of the publish
 // paths (docs/ROADMAP.md §7.2 file 6.9): realm, device, and interface
 // resolution with the ownership gate parameterised, §2.6 payload
 // validation, and the persist switch. ownership selects which side may
-// write and which sentinel an ownership mismatch carries; delivery stays
-// with the caller.
+// write and which sentinel an ownership mismatch carries. On success the
+// op persists immediately and then runs the same post-commit observers as
+// the broker ingress path — previous-value capture for the change-derived
+// triggers, then fireData (data triggers + live bus), so a server- or
+// device-owned write is observed exactly like a real device's data
+// (devicedata_test.go:3-4). Wire delivery stays with the caller.
 func (e *Engine) publishAsOwner(ctx context.Context, realm string, id deviceid.ID,
 	ifaceName, path string, value json.RawMessage, ts *time.Time,
 	ownership interfaceschema.Ownership) (*ownerDelivery, error) {
 	notOwned := ErrNotDeviceOwned

Now `UnsetServerProperty`:
→ Read internal/engine/serverdata.go [limit=40, offset=245]

← Edit internal/engine/serverdata.go
Index: /root/astrate-mule/internal/engine/serverdata.go
===================================================================
--- /root/astrate-mule/internal/engine/serverdata.go
+++ /root/astrate-mule/internal/engine/serverdata.go
@@ -245,15 +245,32 @@
 if !mapping.AllowUnset {
 	return fmt.Errorf("%w: %s%s", ErrUnsetNotAllowed, ifaceName, path)
 }
 
+// Previous-value snapshot taken before the row is gone, so the
+// change-derived conditions (path_removed, value_change) evaluate against
+// the pre-unset value, then the op fires the same observers as any other
+// committed unset.
+op := PersistOp{
+	Realm:     realm,
+	RealmID:   rs.id,
+	DeviceID:  id,
+	Interface: ci,
+	Mapping:   mapping,
+	Path:      path,
+	Kind:      OpPropertyUnset,
+	TS:        time.Now().UTC(),
+}
+e.capturePrevious(ctx, rs, &op)
+
 existed, err := e.st.UnsetProperty(ctx, rs.id, id, ci.ID, path)
 if err != nil {
 	return err
 }
 if existed {
 	e.met.persistOps.WithLabelValues(OpPropertyUnset.String()).Inc()
 }
+e.fireData(rs, &op)
 
 topic := deviceTopic(realm, id, ifaceName+path)
 if err := e.broker.Publish(topic, []byte{}, 2, true, 0); err != nil {
 	return fmt.Errorf("engine: publishing %s: %w", topic, err)

# Todos
[✓] Extract lookupPrev + add capturePrevious in data.go
[✓] Route publishAsOwner through capturePrevious + fireData in serverdata.go
[✓] Route UnsetServerProperty through capturePrevious + fireData
[•] Add failing bus+trigger tests in serverdata_test.go and devicedata_test.go
[ ] Run gate: go vet, go test, gofmt

Now mark progress and add the tests. First the server test:
← Edit internal/engine/serverdata_test.go
Index: /root/astrate-mule/internal/engine/serverdata_test.go
===================================================================
--- /root/astrate-mule/internal/engine/serverdata_test.go
+++ /root/astrate-mule/internal/engine/serverdata_test.go
@@ -7,8 +7,9 @@
 	"errors"
 	"testing"
 	"time"
 
+	"github.com/astrate-platform/astrate/internal/engine/stream"
 	"github.com/astrate-platform/astrate/internal/store"
 	"github.com/astrate-platform/astrate/pkg/deviceid"
 	"github.com/astrate-platform/astrate/pkg/payload"
 )

→ Read internal/engine/serverdata_test.go [limit=15, offset=300]

← Edit internal/engine/serverdata_test.go
Index: /root/astrate-mule/internal/engine/serverdata_test.go
===================================================================
--- /root/astrate-mule/internal/engine/serverdata_test.go
+++ /root/astrate-mule/internal/engine/serverdata_test.go
@@ -310,4 +310,81 @@
 	if !errors.Is(err, ErrPathNotFound) {
 		t.Errorf("bad prefix err = %v, want ErrPathNotFound", err)
 	}
 }
+
+// nextBusEvent returns the next live bus event, failing the test if none
+// arrives within the timeout.
+func nextBusEvent(t *testing.T, events <-chan stream.Event) stream.Event {
+	t.Helper()
+	select {
+	case ev, ok := <-events:
+		if !ok {
+			t.Fatal("bus closed before an event arrived")
+		}
+		return ev
+	case <-time.After(2 * time.Second):
+		t.Fatal("no bus event within 2s")
+		return stream.Event{}
+	}
+}
+
+// TestServerPublishFiresBusAndTriggers: a server-originated write fans out on
+// the live bus and through the data triggers exactly like a broker-ingressed
+// value would — incoming_data on every set, value_change against the previous
+// value, path_removed on unset. The virtual-device contract
+// (devicedata_test.go:3-4) extends to server-owned writes.
+func TestServerPublishFiresBusAndTriggers(t *testing.T) {
+	ctx := context.Background()
+	fw := &fakeForwarder{}
+	rig, fs, _ := newWiredRig(t, Config{Forwarder: fw})
+	const iface = "com.astrate.test.ServerProperties"
+	const path = "/limits/maxConnections"
+	fs.addTrigger(realmAlphaID, "t_in", changeTriggerDef("t_in", "incoming_data", iface, 1, path))
+	fs.addTrigger(realmAlphaID, "t_vc", changeTriggerDef("t_vc", "value_change", iface, 1, path))
+	fs.addTrigger(realmAlphaID, "t_pr", changeTriggerDef("t_pr", "path_removed", iface, 1, path))
+	if err := rig.e.RefreshTriggers(ctx, realmAlphaID); err != nil {
+		t.Fatalf("RefreshTriggers: %v", err)
+	}
+	events, cancel := rig.e.bus.Subscribe(realmAlpha, stream.Filter{}, 16)
+	defer cancel()
+
+	// Set #1: a bus card, incoming_data, and value_change(null → 42) —
+	// the first value creates the path.
+	if err := rig.e.PublishServerValue(ctx, realmAlpha, devAlpha, iface, path,
+		json.RawMessage("42"), nil); err != nil {
+		t.Fatalf("PublishServerValue: %v", err)
+	}
+	if ev := nextBusEvent(t, events); ev.Kind != stream.KindIncomingData ||
+		ev.Interface != iface || ev.Path != path || ev.Value != int32(42) {
+		t.Errorf("set #1 bus event = %+v", ev)
+	}
+	fw.waitForCount(t, "t_in", 1)
+	fw.waitForCount(t, "t_vc", 1)
+	assertEvent(t, fw.ofTrigger("t_vc")[0], "value_change", nil, 42)
+
+	// Set #2, different value: value_change carries the previous one.
+	if err := rig.e.PublishServerValue(ctx, realmAlpha, devAlpha, iface, path,
+		json.RawMessage("43"), nil); err != nil {
+		t.Fatalf("PublishServerValue #2: %v", err)
+	}
+	if ev := nextBusEvent(t, events); ev.Value != int32(43) {
+		t.Errorf("set #2 bus event = %+v", ev)
+	}
+	fw.waitForCount(t, "t_in", 2)
+	fw.waitForCount(t, "t_vc", 2)
+	assertEvent(t, lastOf(fw.ofTrigger("t_vc")), "value_change", 42, 43)
+
+	// Unset: another bus card (null value) plus path_removed.
+	if err := rig.e.UnsetServerProperty(ctx, realmAlpha, devAlpha, iface, path); err != nil {
+		t.Fatalf("UnsetServerProperty: %v", err)
+	}
+	if ev := nextBusEvent(t, events); ev.Kind != stream.KindIncomingData ||
+		ev.Interface != iface || ev.Path != path {
+		t.Errorf("unset bus event = %+v", ev)
+	}
+	fw.waitForCount(t, "t_in", 3)
+	fw.waitForCount(t, "t_pr", 1)
+	if got := fw.ofTrigger("t_pr")[0].Event.Type; got != "path_removed" {
+		t.Errorf("unset trigger type = %q, want path_removed", got)
+	}
+}

← Edit internal/engine/devicedata_test.go
Index: /root/astrate-mule/internal/engine/devicedata_test.go
===================================================================
--- /root/astrate-mule/internal/engine/devicedata_test.go
+++ /root/astrate-mule/internal/engine/devicedata_test.go
@@ -10,8 +10,9 @@
 	"errors"
 	"testing"
 	"time"
 
+	"github.com/astrate-platform/astrate/internal/engine/stream"
 	"github.com/astrate-platform/astrate/internal/store"
 	"github.com/astrate-platform/astrate/pkg/deviceid"
 	"github.com/astrate-platform/astrate/pkg/payload"
 )

← Edit internal/engine/devicedata_test.go
Index: /root/astrate-mule/internal/engine/devicedata_test.go
===================================================================
--- /root/astrate-mule/internal/engine/devicedata_test.go
+++ /root/astrate-mule/internal/engine/devicedata_test.go
@@ -183,4 +183,41 @@
 		t.Errorf("failed publishes left traces: %d upserts, %d batches, %d publishes",
 			upserts, fs.batchCount(), len(port.published()))
 	}
 }
+
+// TestPublishDeviceValueFiresBusAndTriggers: a device-owned ingest fans out
+// on the live bus and through the data triggers exactly like a real device's
+// data would, while still producing zero MQTT traffic — the virtual-device
+// contract in both directions.
+func TestPublishDeviceValueFiresBusAndTriggers(t *testing.T) {
+	ctx := context.Background()
+	fw := &fakeForwarder{}
+	rig, fs, port := newWiredRig(t, Config{Forwarder: fw})
+	addDeviceProps(t, rig, fs)
+	const iface = "com.astrate.test.DeviceProperties"
+	const path = "/state"
+	fs.addTrigger(realmAlphaID, "t_dev", changeTriggerDef("t_dev", "incoming_data", iface, 1, path))
+	if err := rig.e.RefreshTriggers(ctx, realmAlphaID); err != nil {
+		t.Fatalf("RefreshTriggers: %v", err)
+	}
+	events, cancel := rig.e.bus.Subscribe(realmAlpha, stream.Filter{}, 8)
+	defer cancel()
+
+	if err := rig.e.PublishDeviceValue(ctx, realmAlpha, devAlpha, iface, path,
+		json.RawMessage(`"on"`), nil); err != nil {
+		t.Fatalf("PublishDeviceValue: %v", err)
+	}
+
+	if ev := nextBusEvent(t, events); ev.Kind != stream.KindIncomingData ||
+		ev.Interface != iface || ev.Path != path || ev.Value != "on" {
+		t.Errorf("bus event = %+v", ev)
+	}
+	fw.waitForCount(t, "t_dev", 1)
+	ev := fw.ofTrigger("t_dev")[0]
+	if ev.Event.Type != "incoming_data" || ev.Event.Value != "on" {
+		t.Errorf("trigger event = %+v", ev.Event)
+	}
+	if pubs := port.published(); len(pubs) != 0 {
+		t.Errorf("device-owned ingest published to the broker: %+v", pubs)
+	}
+}

Now run the gate:
$ go build ./... 2>&1 | head -50
(no output)

$ go vet ./... 2>&1 | head -50 && gofmt -l . 2>&1 | head -50
(no output)

$ go test ./internal/engine/ -run 'TestServerPublishFiresBusAndTriggers|TestPublishDeviceValueFiresBusAndTriggers' -v 2>&1 | tail -40
=== RUN   TestPublishDeviceValueFiresBusAndTriggers
--- PASS: TestPublishDeviceValueFiresBusAndTriggers (0.01s)
=== RUN   TestServerPublishFiresBusAndTriggers
    serverdata_test.go:379: UnsetServerProperty: engine: mapping does not allow unset: com.astrate.test.ServerProperties/limits/maxConnections
--- FAIL: TestServerPublishFiresBusAndTriggers (0.02s)
FAIL
FAIL	github.com/astrate-platform/astrate/internal/engine	0.061s
FAIL

→ Read pkg/interfaceschema/testdata/valid/com.astrate.test.ServerProperties.json
`/limits/maxConnections` has no `allow_unset`. Switching the test to `/identity/displayName` (string, `allow_unset`):

← Edit internal/engine/serverdata_test.go
Index: /root/astrate-mule/internal/engine/serverdata_test.go
===================================================================
--- /root/astrate-mule/internal/engine/serverdata_test.go
+++ /root/astrate-mule/internal/engine/serverdata_test.go
@@ -337,9 +337,9 @@
 ctx := context.Background()
 fw := &fakeForwarder{}
 rig, fs, _ := newWiredRig(t, Config{Forwarder: fw})
 const iface = "com.astrate.test.ServerProperties"
-const path = "/limits/maxConnections"
+const path = "/identity/displayName"
 fs.addTrigger(realmAlphaID, "t_in", changeTriggerDef("t_in", "incoming_data", iface, 1, path))
 fs.addTrigger(realmAlphaID, "t_vc", changeTriggerDef("t_vc", "value_change", iface, 1, path))
 fs.addTrigger(realmAlphaID, "t_pr", changeTriggerDef("t_pr", "path_removed", iface, 1, path))
 if err := rig.e.RefreshTriggers(ctx, realmAlphaID); err != nil {
@@ -347,33 +347,33 @@
 }
 events, cancel := rig.e.bus.Subscribe(realmAlpha, stream.Filter{}, 16)
 defer cancel()
 
-// Set #1: a bus card, incoming_data, and value_change(null → 42) —
+// Set #1: a bus card, incoming_data, and value_change(null → "alpha") —
 // the first value creates the path.
 if err := rig.e.PublishServerValue(ctx, realmAlpha, devAlpha, iface, path,
-	json.RawMessage("42"), nil); err != nil {
+	json.RawMessage(`"alpha"`), nil); err != nil {
 	t.Fatalf("PublishServerValue: %v", err)
 }
 if ev := nextBusEvent(t, events); ev.Kind != stream.KindIncomingData ||
-	ev.Interface != iface || ev.Path != path || ev.Value != int32(42) {
+	ev.Interface != iface || ev.Path != path || ev.Value != "alpha" {
 	t.Errorf("set #1 bus event = %+v", ev)
 }
 fw.waitForCount(t, "t_in", 1)
 fw.waitForCount(t, "t_vc", 1)
-assertEvent(t, fw.ofTrigger("t_vc")[0], "value_change", nil, 42)
+assertEvent(t, fw.ofTrigger("t_vc")[0], "value_change", nil, "alpha")
 
 // Set #2, different value: value_change carries the previous one.
 if err := rig.e.PublishServerValue(ctx, realmAlpha, devAlpha, iface, path,
-	json.RawMessage("43"), nil); err != nil {
+	json.RawMessage(`"beta"`), nil); err != nil {
 	t.Fatalf("PublishServerValue #2: %v", err)
 }
-if ev := nextBusEvent(t, events); ev.Value != int32(43) {
+if ev := nextBusEvent(t, events); ev.Value != "beta" {
 	t.Errorf("set #2 bus event = %+v", ev)
 }
 fw.waitForCount(t, "t_in", 2)
 fw.waitForCount(t, "t_vc", 2)
-assertEvent(t, lastOf(fw.ofTrigger("t_vc")), "value_change", 42, 43)
+assertEvent(t, lastOf(fw.ofTrigger("t_vc")), "value_change", "alpha", "beta")
 
 // Unset: another bus card (null value) plus path_removed.
 if err := rig.e.UnsetServerProperty(ctx, realmAlpha, devAlpha, iface, path); err != nil {
 	t.Fatalf("UnsetServerProperty: %v", err)

$ go test ./internal/engine/ -run 'TestServerPublishFiresBusAndTriggers|TestPublishDeviceValueFiresBusAndTriggers' -v 2>&1 | tail -20
=== RUN   TestPublishDeviceValueFiresBusAndTriggers
--- PASS: TestPublishDeviceValueFiresBusAndTriggers (0.01s)
=== RUN   TestServerPublishFiresBusAndTriggers
--- PASS: TestServerPublishFiresBusAndTriggers (0.02s)
PASS
ok  	github.com/astrate-platform/astrate/internal/engine	0.056s

Now run the full engine test suite:
$ go test ./internal/engine/ -count=1 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/engine	1.706s

$ go test ./... -count=1 2>&1 | tail -40
ok  	github.com/astrate-platform/astrate/cmd/astrate	0.626s
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
ok  	github.com/astrate-platform/astrate/internal/appengine	0.115s
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	33.181s
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	3.459s
ok  	github.com/astrate-platform/astrate/internal/auth	4.910s
ok  	github.com/astrate-platform/astrate/internal/broker	2.708s
ok  	github.com/astrate-platform/astrate/internal/config	0.045s
ok  	github.com/astrate-platform/astrate/internal/engine	2.325s
ok  	github.com/astrate-platform/astrate/internal/engine/forward	0.120s
ok  	github.com/astrate-platform/astrate/internal/engine/stream	0.089s
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	2.268s
ok  	github.com/astrate-platform/astrate/internal/flow	0.635s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	1.220s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	0.094s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	0.590s
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	0.058s
ok  	github.com/astrate-platform/astrate/internal/flowapi	0.216s
ok  	github.com/astrate-platform/astrate/internal/housekeeping	0.017s
ok  	github.com/astrate-platform/astrate/internal/httpx	0.017s
ok  	github.com/astrate-platform/astrate/internal/observability	0.115s
ok  	github.com/astrate-platform/astrate/internal/pairing	15.746s
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	26.969s
ok  	github.com/astrate-platform/astrate/internal/realm	0.067s
ok  	github.com/astrate-platform/astrate/internal/store	0.034s
ok  	github.com/astrate-platform/astrate/internal/swagger	0.037s
ok  	github.com/astrate-platform/astrate/internal/testutil	0.174s
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	0.141s
ok  	github.com/astrate-platform/astrate/pkg/deviceid	0.022s
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	0.226s
ok  	github.com/astrate-platform/astrate/pkg/payload	0.140s

All green. Let me do a final gofmt/vet pass and review the edits:
$ gofmt -l . && go vet ./... 2>&1 | head -20; echo "gates-done"
gates-done

→ Read internal/engine/serverdata.go [limit=130, offset=152]

$ cp internal/engine/serverdata.go /tmp/opencode/serverdata.go.bak && cp internal/engine/data.go /tmp/opencode/data.go.bak
(no output)

Let me verify the runner's acceptance criterion — the new tests must fail without the implementation. Temporarily removing the two `fireData` calls:
← Edit internal/engine/serverdata.go
Index: /root/astrate-mule/internal/engine/serverdata.go
===================================================================
--- /root/astrate-mule/internal/engine/serverdata.go
+++ /root/astrate-mule/internal/engine/serverdata.go
@@ -208,9 +208,9 @@
 			return nil, err
 		}
 	}
 	e.met.persistOps.WithLabelValues(op.Kind.String()).Inc()
-	e.fireData(rs, &op)
+	_ = rs // temporary: implementation removed to verify tests fail alone
 	return out, nil
 }
 
 // UnsetServerProperty deletes a server-owned property (AppEngine DELETE,

$ cp /tmp/opencode/serverdata.go.bak internal/engine/serverdata.go
(no output)

The edit tool may have corrupted the line spacing. Let me restore and verify properly:
→ Read /tmp/opencode/serverdata.go.bak [limit=5, offset=210]
Good, restore is clean. Let me verify the tests still pass after the full restore, then verify the tests fail without implementation:

$ go test ./internal/engine/ -run 'TestServerPublishFiresBusAndTriggers|TestPublishDeviceValueFiresBusAndTriggers' -v 2>&1 | grep -E 'PASS|FAIL'
--- PASS: TestPublishDeviceValueFiresBusAndTriggers (0.01s)
--- PASS: TestServerPublishFiresBusAndTriggers (0.02s)
PASS

$ cp internal/engine/serverdata.go /tmp/opencode/serverdata.go.bak2 && \
  sed -i 's/e\.fireData(rs, \&op)/_ = "disabled"/g' internal/engine/serverdata.go && \
  go test ./internal/engine/ -run 'TestServerPublishFiresBusAndTriggers|TestPublishDeviceValueFiresBusAndTriggers' -timeout 30s 2>&1; echo "EXIT:$?" && \
  cp /tmp/opencode/serverdata.go.bak2 internal/engine/serverdata.go
--- FAIL: TestPublishDeviceValueFiresBusAndTriggers (2.01s)
    devicedata_test.go:211: no bus event within 2s
--- FAIL: TestServerPublishFiresBusAndTriggers (2.01s)
    serverdata_test.go:357: no bus event within 2s
FAIL
FAIL	github.com/astrate-platform/astrate/internal/engine	4.049s
FAIL
EXIT:1

Both new tests fail without the implementation — runner acceptance criterion met. Final clean check:
$ go test ./internal/engine/ -count=1 2>&1 | tail -3 && go vet ./... 2>&1 | head -5 && gofmt -l . 2>&1 | head -5
ok  	github.com/astrate-platform/astrate/internal/engine	1.593s

# Todos
[✓] Extract lookupPrev + add capturePrevious in data.go
[✓] Route publishAsOwner through capturePrevious + fireData in serverdata.go
[✓] Route UnsetServerProperty through capturePrevious + fireData
[✓] Add failing bus+trigger tests in serverdata_test.go and devicedata_test.go
[✓] Run gate: go vet, go test, gofmt

Done: route server- and device-owned publish paths through fireData (bus + data triggers) with previous-value capture for change-derived conditions
Files: internal/engine/data.go, internal/engine/serverdata.go, internal/engine/serverdata_test.go, internal/engine/devicedata_test.go
Verified: `go test ./...` → all pass; `go vet ./...` → clean; `gofmt -l .` → clean; new tests fail when fireData calls removed (runner acceptance criterion: "no bus event within 2s")
Uns
