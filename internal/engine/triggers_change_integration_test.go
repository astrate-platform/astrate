//go:build integration

package engine

// T2 coverage for the change-derived data triggers (docs/ROADMAP.md §7.3):
// the previous-value lookups run against a real TimescaleDB through the
// store (individual_datastreams ORDER BY ts DESC LIMIT 1 and the properties
// table), one subtest per condition family.

import (
	"context"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// newIntegrationDevice registers a fresh device declaring the given
// interfaces, so change-trigger scenarios start from an empty history.
func (r *engineRig) newIntegrationDevice(t *testing.T, ifaces ...string) (dev deviceid.ID, err error) {
	t.Helper()
	ctx := context.Background()
	dev, err = deviceid.Random()
	if err != nil {
		return dev, err
	}
	if err := r.raw.RegisterDevice(ctx, r.realmID, dev, "secret-hash"); err != nil {
		return dev, err
	}
	intro := make(map[string]store.InterfaceVersion, len(ifaces))
	for _, name := range ifaces {
		si := r.ifaces[name]
		intro[name] = store.InterfaceVersion{Major: si.Major, Minor: si.Minor}
	}
	_, err = r.raw.UpdateIntrospection(ctx, r.realmID, dev, intro)
	return dev, err
}

// testEngineChangeTriggers drives the five change-derived conditions over
// the real store: creation, change, no-change, stored, and property removal.
// Wired as a TestEngine subtest (T2).
func testEngineChangeTriggers(t *testing.T, rig *engineRig) {
	ctx := context.Background()

	fw := &fakeForwarder{}
	e, err := New(rig.st, nil, Config{
		Shards: 1, BatchMaxRows: 8, BatchMaxWait: 10 * time.Millisecond,
		Logger: discardLogger(), Forwarder: fw,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := e.schemas.loadAll(ctx); err != nil {
		t.Fatalf("loadAll: %v", err)
	}

	// A fresh device so earlier subtests' publishes on /double don't count
	// as this device's previous values.
	dev2, err := rig.newIntegrationDevice(t, "com.astrate.test.AllScalarTypes", "com.astrate.test.PropertyArrays")
	if err != nil {
		t.Fatalf("newIntegrationDevice: %v", err)
	}

	for name, def := range map[string]string{
		"t_vc":  changeTriggerDef("t_vc", "value_change", scalarsIface, 1, doublePath),
		"t_vca": changeTriggerDef("t_vca", "value_change_applied", scalarsIface, 1, doublePath),
		"t_pc":  changeTriggerDef("t_pc", "path_created", scalarsIface, 1, doublePath),
		"t_vs":  changeTriggerDef("t_vs", "value_stored", scalarsIface, 1, doublePath),
	} {
		if _, err := rig.raw.CreateTrigger(ctx, rig.realmID, name, []byte(def)); err != nil {
			t.Fatalf("CreateTrigger(%s): %v", name, err)
		}
	}
	prDef := `{
		"name": "t_pr",
		"action": {"amqp_exchange": "astarte_events", "amqp_routing_key": "k"},
		"simple_triggers": [{
			"type": "data_trigger", "on": "path_removed",
			"interface_name": "` + propsIface + `", "interface_major": 2,
			"match_path": "` + thresholdsPath + `", "value_match_operator": "*"
		}]
	}`
	if _, err := rig.raw.CreateTrigger(ctx, rig.realmID, "t_pr", []byte(prDef)); err != nil {
		t.Fatalf("CreateTrigger(t_pr): %v", err)
	}
	if err := e.RefreshTriggers(ctx, rig.realmID); err != nil {
		t.Fatalf("RefreshTriggers: %v", err)
	}

	dev2Msg := func(iface, path string, body []byte) broker.InboundMessage {
		topic := rig.realm + "/" + dev2.String()
		if iface != "" {
			topic += "/" + iface + path
		}
		return broker.InboundMessage{
			Realm: rig.realm, DeviceID: dev2, Topic: topic, Payload: body,
			QoS: 1, ReceivedAt: time.Now().UTC().Truncate(time.Millisecond),
			Ack: (&ackCounter{}).fn(),
		}
	}
	pub := func(iface, path string, body []byte) {
		t.Helper()
		e.handle(ctx, e.shards[0], dev2Msg(iface, path, body))
		e.flushShard(ctx, e.shards[0])
	}

	// Creation: value_change(null→22.5), applied, path_created, stored.
	pub(scalarsIface, doublePath, enc(t, 22.5, nil, payload.FormatBSON))
	fw.waitForCount(t, "t_vc", 1)
	fw.waitForCount(t, "t_pc", 1)
	assertEvent(t, fw.ofTrigger("t_vc")[0], "value_change", nil, 22.5)
	assertValue(t, fw.ofTrigger("t_pc")[0], "path_created", 22.5)
	assertEvent(t, fw.ofTrigger("t_vca")[0], "value_change_applied", nil, 22.5)
	assertStored(t, fw.ofTrigger("t_vs")[0], 22.5)

	// Change: old carried through, no new creation.
	pub(scalarsIface, doublePath, enc(t, 23.5, nil, payload.FormatBSON))
	fw.waitForCount(t, "t_vc", 2)
	assertNoNew(t, fw.ofTrigger("t_pc"), 1, "path_created")
	assertHasEvent(t, fw.ofTrigger("t_vc"), "value_change", 22.5, 23.5)

	// Same value: stored fires, change events don't.
	pub(scalarsIface, doublePath, enc(t, 23.5, nil, payload.FormatBSON))
	fw.waitForCount(t, "t_vs", 3)
	assertNoNew(t, fw.ofTrigger("t_vc"), 2, "value_change")

	// Property set then unset on the real properties table: one removal.
	pub(propsIface, thresholdsPath, []byte(`{"v":[1.5,2.5]}`))
	waitFor(t, 10*time.Second, "thresholds property committed", func() bool {
		_, err := rig.raw.GetProperty(ctx, rig.realmID, dev2, rig.ifaces[propsIface].ID, thresholdsPath)
		return err == nil
	})
	pub(propsIface, thresholdsPath, nil)
	fw.waitForCount(t, "t_pr", 1)
	if ev := fw.ofTrigger("t_pr")[0]; ev.Event.Type != "path_removed" {
		t.Errorf("event type = %q, want path_removed", ev.Event.Type)
	}
}
