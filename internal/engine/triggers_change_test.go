package engine

// T1 coverage for the change-derived data triggers (value_change,
// value_change_applied, path_created, path_removed, value_stored): the full
// pipeline — accept-time previous-value capture, post-commit emission, and
// delivery through the Forwarder seam.

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/pkg/payload"
)

// capturedEvent is one forwarded trigger envelope, parsed for assertions.
type capturedEvent struct {
	TriggerName string `json:"trigger_name"`
	Event       struct {
		Type     string `json:"type"`
		OldValue any    `json:"old_value"`
		NewValue any    `json:"new_value"`
		Value    any    `json:"value"`
	} `json:"event"`
}

// fakeForwarder records every custom-action delivery.
type fakeForwarder struct {
	mu     sync.Mutex
	events []capturedEvent
}

func (fw *fakeForwarder) Forward(_ context.Context, _, _ string, _ json.RawMessage, event []byte) error {
	var ce capturedEvent
	if err := json.Unmarshal(event, &ce); err != nil {
		return err
	}
	fw.mu.Lock()
	fw.events = append(fw.events, ce)
	fw.mu.Unlock()
	return nil
}

func (fw *fakeForwarder) collected() []capturedEvent {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	return append([]capturedEvent(nil), fw.events...)
}

// ofTrigger filters the collected events by trigger name.
func (fw *fakeForwarder) ofTrigger(name string) []capturedEvent {
	var out []capturedEvent
	for _, ev := range fw.collected() {
		if ev.TriggerName == name {
			out = append(out, ev)
		}
	}
	return out
}

// waitForCount blocks until the trigger delivered n events.
func (fw *fakeForwarder) waitForCount(t *testing.T, name string, n int) {
	t.Helper()
	waitFor(t, 5*time.Second, "deliveries on "+name, func() bool {
		return len(fw.ofTrigger(name)) >= n
	})
}

// changeTriggerDef builds a single-condition trigger with a custom action so
// deliveries land on the fakeForwarder instead of a webhook. Non-amqp shape:
// amqp_exchange actions are rejected at compile time since #64.
func changeTriggerDef(name, on, iface string, major int, matchPath string) string {
	return `{
		"name": "` + name + `",
		"action": {"nats_subject": "astarte_events"},
		"simple_triggers": [{
			"type": "data_trigger", "on": "` + on + `",
			"interface_name": "` + iface + `", "interface_major": ` + strconv.Itoa(major) + `,
			"match_path": "` + matchPath + `", "value_match_operator": "*"
		}]
	}`
}

const (
	scalarsIface   = "com.astrate.test.AllScalarTypes"
	propsIface     = "com.astrate.test.PropertyArrays"
	doublePath     = "/double"
	thresholdsPath = "/config/thresholds"
)

// installChangeTriggers registers one trigger per change-derived condition.
func installChangeTriggers(t *testing.T, fs *fakeStore) {
	t.Helper()
	fs.addTrigger(realmAlphaID, "t_vc", changeTriggerDef("t_vc", "value_change", scalarsIface, 1, doublePath))
	fs.addTrigger(realmAlphaID, "t_vca", changeTriggerDef("t_vca", "value_change_applied", scalarsIface, 1, doublePath))
	fs.addTrigger(realmAlphaID, "t_pc", changeTriggerDef("t_pc", "path_created", scalarsIface, 1, doublePath))
	fs.addTrigger(realmAlphaID, "t_vs", changeTriggerDef("t_vs", "value_stored", scalarsIface, 1, doublePath))
	fs.addTrigger(realmAlphaID, "t_pr", changeTriggerDef("t_pr", "path_removed", propsIface, 2, thresholdsPath))
	fs.addTrigger(realmAlphaID, "t_in", changeTriggerDef("t_in", "incoming_data", scalarsIface, 1, doublePath))
}

// TestChangeTriggersLifecycle drives the whole lifecycle: first publish
// creates the path and fires value_change with a null old value; a changed
// value carries the previous one; an unchanged value fires only value_stored;
// a property unset removes its path. Object-aggregated publishes stay outside
// the scope.
func TestChangeTriggersLifecycle(t *testing.T) {
	ctx := t.Context()
	fw := &fakeForwarder{}
	rig, fs, _ := newWiredRig(t, Config{Forwarder: fw})
	installChangeTriggers(t, fs)
	if err := rig.e.RefreshTriggers(ctx, realmAlphaID); err != nil {
		t.Fatalf("RefreshTriggers: %v", err)
	}
	lookupsBefore := fs.latestLookups.Load()

	pub := func(iface, path string, body []byte) {
		rig.handle(deviceMsg(iface, path, 1, body, &ackCounter{}))
		rig.e.flushShard(ctx, rig.sh)
	}

	// 1. First value: path_created + value_change(null→22.5) + applied + stored.
	pub(scalarsIface, doublePath, enc(t, 22.5, nil, payload.FormatBSON))
	fw.waitForCount(t, "t_vc", 1)
	fw.waitForCount(t, "t_pc", 1)
	fw.waitForCount(t, "t_vca", 1)
	fw.waitForCount(t, "t_vs", 1)
	assertEvent(t, fw.ofTrigger("t_vc")[0], "value_change", nil, 22.5)
	assertValue(t, fw.ofTrigger("t_pc")[0], "path_created", 22.5)
	assertEvent(t, fw.ofTrigger("t_vca")[0], "value_change_applied", nil, 22.5)
	assertStored(t, fw.ofTrigger("t_vs")[0], 22.5)

	assertValue(t, lastOf(fw.ofTrigger("t_in")), "incoming_data", 22.5)

	// 2. Changed value: no creation this time, old carried through.
	pub(scalarsIface, doublePath, enc(t, 23.5, nil, payload.FormatBSON))
	fw.waitForCount(t, "t_vc", 2)
	fw.waitForCount(t, "t_vs", 2)
	assertNoNew(t, fw.ofTrigger("t_pc"), 1, "path_created")
	assertEvent(t, lastOf(fw.ofTrigger("t_vc")), "value_change", 22.5, 23.5)

	// incoming_data keeps firing for every accepted publish.
	fw.waitForCount(t, "t_in", 2)

	// 3. Same value again: only value_stored, no change events.
	pub(scalarsIface, doublePath, enc(t, 23.5, nil, payload.FormatBSON))
	fw.waitForCount(t, "t_vs", 3)
	assertNoNew(t, fw.ofTrigger("t_vc"), 2, "value_change")
	assertNoNew(t, fw.ofTrigger("t_vca"), 2, "value_change_applied")
	assertNoNew(t, fw.ofTrigger("t_pc"), 1, "path_created")

	// 4. Property set then unset: removal fires once, on the unset only.
	pub(propsIface, thresholdsPath, []byte(`{"v":[1.5,2.5]}`))
	pub(propsIface, thresholdsPath, nil) // empty payload = unset
	fw.waitForCount(t, "t_pr", 1)
	ev := fw.ofTrigger("t_pr")[0]
	if ev.Event.Type != "path_removed" || ev.Event.Value != nil {
		t.Errorf("path_removed event = %+v", ev.Event)
	}

	// 5. An unwatched interface never pays the lookup: Minimal has no
	// change-derived trigger installed.
	watchLookups := fs.latestLookups.Load()
	pub("com.astrate.test.Minimal", "/value", enc(t, 1.25, nil, payload.FormatJSON))
	rig.e.flushShard(ctx, rig.sh)
	if fs.latestLookups.Load() != watchLookups {
		t.Error("previous-value lookup ran for an unwatched endpoint")
	}
	if lookupsBefore != 0 {
		t.Errorf("%d lookups before any watched publish", lookupsBefore)
	}
}

// TestChangeTriggersIntraBatch: two publishes of the same path inside one
// batch chain their previous values without a committed row in between.
func TestChangeTriggersIntraBatch(t *testing.T) {
	ctx := t.Context()
	fw := &fakeForwarder{}
	rig, fs, _ := newWiredRig(t, Config{Forwarder: fw})
	installChangeTriggers(t, fs)
	if err := rig.e.RefreshTriggers(ctx, realmAlphaID); err != nil {
		t.Fatalf("RefreshTriggers: %v", err)
	}

	ack := &ackCounter{}
	rig.handle(deviceMsg(scalarsIface, doublePath, 1, enc(t, 1.0, nil, payload.FormatBSON), ack))
	rig.handle(deviceMsg(scalarsIface, doublePath, 1, enc(t, 2.0, nil, payload.FormatBSON), ack))
	rig.e.flushShard(ctx, rig.sh)

	fw.waitForCount(t, "t_vc", 2)
	assertHasEvent(t, fw.ofTrigger("t_vc"), "value_change", nil, 1.0)
	assertHasEvent(t, fw.ofTrigger("t_vc"), "value_change", 1.0, 2.0)
}

// assertEvent checks type and old/new payloads of a change event.
func assertEvent(t *testing.T, ev capturedEvent, typ string, oldValue, newValue any) {
	t.Helper()
	if ev.Event.Type != typ {
		t.Fatalf("event type = %q, want %q (%+v)", ev.Event.Type, typ, ev)
	}
	if !jsonAnyEqual(ev.Event.OldValue, oldValue) {
		t.Errorf("%s old_value = %v, want %v", typ, ev.Event.OldValue, oldValue)
	}
	if !jsonAnyEqual(ev.Event.NewValue, newValue) {
		t.Errorf("%s new_value = %v, want %v", typ, ev.Event.NewValue, newValue)
	}
}

// assertValue checks a single-value event body (path_created et al).
func assertValue(t *testing.T, ev capturedEvent, typ string, value any) {
	t.Helper()
	if ev.Event.Type != typ {
		t.Fatalf("event type = %q, want %q (%+v)", ev.Event.Type, typ, ev)
	}
	if !jsonAnyEqual(ev.Event.Value, value) {
		t.Errorf("%s value = %v, want %v", typ, ev.Event.Value, value)
	}
}

// assertHasEvent finds one event with the given old/new payloads regardless
// of delivery order (executor workers may process out of order).
func assertHasEvent(t *testing.T, events []capturedEvent, typ string, oldValue, newValue any) {
	t.Helper()
	for _, ev := range events {
		if ev.Event.Type == typ &&
			jsonAnyEqual(ev.Event.OldValue, oldValue) &&
			jsonAnyEqual(ev.Event.NewValue, newValue) {
			return
		}
	}
	t.Fatalf("no %s event with old=%v new=%v among %+v", typ, oldValue, newValue, events)
}

// assertStored checks a value_stored payload.
func assertStored(t *testing.T, ev capturedEvent, value any) {
	t.Helper()
	if ev.Event.Type != "value_stored" {
		t.Fatalf("event type = %q, want value_stored", ev.Event.Type)
	}
	if !jsonAnyEqual(ev.Event.Value, value) {
		t.Errorf("value_stored value = %v, want %v", ev.Event.Value, value)
	}
}

// assertNoNew fails when more occurrences of a condition appeared than seen
// at the previous checkpoint.
func assertNoNew(t *testing.T, events []capturedEvent, limit int, what string) {
	t.Helper()
	if len(events) > limit {
		t.Errorf("%s fired %d times, want at most %d", what, len(events), limit)
	}
}

func lastOf(xs []capturedEvent) capturedEvent { return xs[len(xs)-1] }

// jsonAnyEqual compares two decoded JSON values numerically.
func jsonAnyEqual(a, b any) bool {
	na, aok := testNum(a)
	nb, bok := testNum(b)
	if aok && bok {
		return na == nb
	}
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	as, aok := a.(string)
	bs, bok := b.(string)
	if aok && bok {
		return as == bs
	}
	return string(jsonMarshalForTest(a)) == string(jsonMarshalForTest(b)) // arrays et al.
}

// testNum widens JSON numbers to float64.
func testNum(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	default:
		return 0, false
	}
}

func jsonMarshalForTest(v any) string {
	js, _ := json.Marshal(v)
	return string(js)
}
