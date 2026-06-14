package triggers

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/internal/testutil"
)

// eventDeviceID and eventTS pin the golden envelopes.
const eventDeviceID = "f0VMRgIBAQAAAAAAAAAAAA"

var eventTS = time.Date(2026, 6, 12, 10, 0, 0, 123_000_000, time.UTC)

// envelope wraps an event body the way the engine does.
func envelope(body any) SimpleEvent {
	return SimpleEvent{
		Timestamp:   eventTS,
		DeviceID:    eventDeviceID,
		TriggerName: "example_trigger",
		Event:       body,
	}
}

// TestEventGoldens freezes every event payload byte-for-byte
// (docs/ROADMAP.md §7.3 "event JSON golden vs upstream SimpleEvent shape").
func TestEventGoldens(t *testing.T) {
	cases := []struct {
		name string
		body any
	}{
		{name: "incoming_data.json", body: NewIncomingDataEvent(
			"org.astarte-platform.genericsensors.Values", "/streamTest/value", 0.3)},
		{name: "incoming_data_unset.json", body: NewIncomingDataEvent(
			"com.example.Props", "/setting", nil)},
		{name: "device_connected.json", body: NewDeviceConnectedEvent("203.0.113.89")},
		{name: "device_disconnected.json", body: NewDeviceDisconnectedEvent()},
		{name: "device_error.json", body: NewDeviceErrorEvent(
			"interface_not_in_introspection", map[string]string{"detail": "no introspected interface matches x"})},
		{name: "device_error_no_metadata.json", body: NewDeviceErrorEvent("unexpected_value", nil)},
		{name: "incoming_introspection.json", body: NewIncomingIntrospectionEvent(
			"com.ex.Sensors:1:0;com.ex.Geo:0:1")},
		{name: "interface_added.json", body: NewInterfaceAddedEvent("com.ex.Sensors", 1, 2)},
		{name: "interface_removed.json", body: NewInterfaceRemovedEvent("com.ex.Sensors", 1)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(envelope(tc.body))
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			testutil.Golden(t, tc.name, got)
		})
	}
}

// TestEventUpstreamShape compares the rendered envelope field-by-field
// against a handcrafted upstream example (Astarte v1.2 trigger engine
// default payload + SimpleEvents JSON encoder), independent of key order.
func TestEventUpstreamShape(t *testing.T) {
	got, err := json.Marshal(envelope(NewIncomingDataEvent(
		"org.astarte-platform.genericsensors.Values", "/streamTest/value", 0.3)))
	if err != nil {
		t.Fatal(err)
	}
	upstream := `{
		"timestamp": "2026-06-12T10:00:00.123Z",
		"device_id": "f0VMRgIBAQAAAAAAAAAAAA",
		"event": {
			"type": "incoming_data",
			"interface": "org.astarte-platform.genericsensors.Values",
			"path": "/streamTest/value",
			"value": 0.3
		},
		"trigger_name": "example_trigger"
	}`

	var gotMap, wantMap map[string]any
	if err := json.Unmarshal(got, &gotMap); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(upstream), &wantMap); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotMap, wantMap) {
		t.Errorf("envelope diverges from upstream shape:\ngot  %v\nwant %v", gotMap, wantMap)
	}

	// device_connected, per the upstream encoder.
	got, err = json.Marshal(envelope(NewDeviceConnectedEvent("203.0.113.89")))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(got, &gotMap); err != nil {
		t.Fatal(err)
	}
	wantEvent := map[string]any{"type": "device_connected", "device_ip_address": "203.0.113.89"}
	if !reflect.DeepEqual(gotMap["event"], wantEvent) {
		t.Errorf("device_connected event: %v, want %v", gotMap["event"], wantEvent)
	}

	// device_error always carries the metadata map.
	got, err = json.Marshal(envelope(NewDeviceErrorEvent("unexpected_value", nil)))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(got, &gotMap); err != nil {
		t.Fatal(err)
	}
	ev, ok := gotMap["event"].(map[string]any)
	if !ok {
		t.Fatalf("event: %v", gotMap["event"])
	}
	if _, ok := ev["metadata"].(map[string]any); !ok {
		t.Errorf("device_error metadata missing or not an object: %v", ev["metadata"])
	}
}

// TestEventTimestampPrecision: envelope timestamps are UTC with millisecond
// precision regardless of the input location/precision.
func TestEventTimestampPrecision(t *testing.T) {
	cet := time.FixedZone("CET", 3600)
	ev := SimpleEvent{
		Timestamp: time.Date(2026, 1, 2, 13, 4, 5, 999_999_999, cet),
		DeviceID:  "d", TriggerName: "t",
		Event: NewDeviceDisconnectedEvent(),
	}
	raw, err := json.Marshal(ev)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if m["timestamp"] != "2026-01-02T12:04:05.999Z" {
		t.Errorf("timestamp = %v, want 2026-01-02T12:04:05.999Z", m["timestamp"])
	}
}
