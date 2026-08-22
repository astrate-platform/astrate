package triggers

import "testing"

// Change-derived condition fixtures (upstream Realm Management JSON shape),
// one per activated "on".
const (
	fixtureValueChangeApplied = `{
		"name": "vc_applied",
		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
		"simple_triggers": [{
			"type": "data_trigger", "on": "value_change_applied",
			"interface_name": "com.example.Sensors", "interface_major": 1,
			"match_path": "/a", "value_match_operator": "*"
		}]
	}`

	fixturePathCreated = `{
		"name": "path_created",
		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
		"simple_triggers": [{
			"type": "data_trigger", "on": "path_created",
			"interface_name": "com.example.Sensors", "interface_major": 1,
			"match_path": "/%{key}", "value_match_operator": "*"
		}]
	}`

	fixturePathRemoved = `{
		"name": "path_removed",
		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
		"simple_triggers": [{
			"type": "data_trigger", "on": "path_removed",
			"interface_name": "com.example.Sensors", "interface_major": 1,
			"match_path": "/setting", "value_match_operator": "*"
		}]
	}`

	fixtureValueStored = `{
		"name": "value_stored",
		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
		"simple_triggers": [{
			"type": "data_trigger", "on": "value_stored",
			"interface_name": "com.example.Sensors", "interface_major": 1,
			"match_path": "/a", "value_match_operator": "==",
			"known_value": 42
		}]
	}`
)

// sensorsEvent is the baseline event for the com.example.Sensors fixtures.
func sensorsEvent(on string) DataEvent {
	return DataEvent{
		DeviceID:  "d",
		On:        on,
		Interface: "com.example.Sensors",
		Major:     1,
		Path:      "/a",
		Value:     1.0,
	}
}

// TestMatchChangeConditions covers the five previously-dormant data
// conditions: each reacts only to its own On, keeps the device/interface/
// path filters, and evaluates its operator the way upstream executes it.
func TestMatchChangeConditions(t *testing.T) {
	change := compile(t, fixtureValueChange)
	applied := compile(t, fixtureValueChangeApplied)
	created := compile(t, fixturePathCreated)
	removed := compile(t, fixturePathRemoved)
	stored := compile(t, fixtureValueStored)

	cases := []struct {
		name    string
		trigger *Trigger
		ev      DataEvent
		want    bool
	}{
		// Each matcher answers only to its own condition.
		{name: "value_change on its own on", trigger: change,
			ev: sensorsEvent(OnValueChange), want: true},
		{name: "value_change ignores incoming_data", trigger: change,
			ev: sensorsEvent(OnIncomingData), want: false},
		{name: "value_change ignores applied", trigger: change,
			ev: sensorsEvent(OnValueChangeApplied), want: false},
		{name: "applied ignores value_change", trigger: applied,
			ev: sensorsEvent(OnValueChange), want: false},
		{name: "created matches first value", trigger: created,
			ev: sensorsEvent(OnPathCreated), want: true},

		// Filters still apply.
		{name: "wrong interface", trigger: change,
			ev: func() DataEvent { ev := sensorsEvent(OnValueChange); ev.Interface = "com.example.Other"; return ev }(), want: false},
		{name: "wrong path", trigger: change,
			ev: func() DataEvent { ev := sensorsEvent(OnValueChange); ev.Path = "/b"; return ev }(), want: false},
		{name: "placeholder created wrong depth", trigger: created,
			ev: func() DataEvent { ev := sensorsEvent(OnPathCreated); ev.Path = "/a/b"; return ev }(), want: false},
		{name: "placeholder created matches segment", trigger: created,
			ev: func() DataEvent { ev := sensorsEvent(OnPathCreated); ev.Path = "/k1"; return ev }(), want: true},

		// Operator evaluation against the (new) value.
		{name: "stored equals known", trigger: stored,
			ev: func() DataEvent { ev := sensorsEvent(OnValueStored); ev.Value = int32(42); return ev }(), want: true},
		{name: "stored differs from known", trigger: stored,
			ev: func() DataEvent { ev := sensorsEvent(OnValueStored); ev.Value = int32(41); return ev }(), want: false},

		// Removals evaluate no value at all (upstream passes none).
		{name: "removed matches unset", trigger: removed,
			ev: func() DataEvent { ev := sensorsEvent(OnPathRemoved); ev.Value = nil; ev.Path = "/setting"; return ev }(), want: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.trigger.MatchesData(tc.ev); got != tc.want {
				t.Errorf("MatchesData(%+v) = %t, want %t", tc.ev, got, tc.want)
			}
		})
	}

	// The evaluated conditions no longer report themselves as unsupported.
	if got := compile(t, fixtureValueChange).Unsupported; len(got) != 0 {
		t.Errorf("value_change reported unsupported: %v", got)
	}
}

// TestTracksChanges gates the engine's previous-value lookup: true only for
// the previous-value conditions with matching filters; incoming_data and
// value_stored never need a snapshot.
func TestTracksChanges(t *testing.T) {
	base := sensorsEvent("")

	cases := []struct {
		name   string
		def    string
		mutate func(*DataEvent)
		want   bool
	}{
		{name: "value_change watches", def: fixtureValueChange, mutate: func(*DataEvent) {}, want: true},
		{name: "applied watches", def: fixtureValueChangeApplied, mutate: func(*DataEvent) {}, want: true},
		{name: "created watches", def: fixturePathCreated, mutate: func(*DataEvent) {}, want: true},
		{name: "removed watches", def: fixturePathRemoved, mutate: func(ev *DataEvent) { ev.Path = "/setting" }, want: true},
		{name: "stored does not watch", def: fixtureValueStored, mutate: func(*DataEvent) {}, want: false},
		{name: "incoming does not watch", def: fixtureValuesAbove, mutate: func(*DataEvent) {}, want: false},
		{name: "filter mismatch", def: fixtureValueChange,
			mutate: func(ev *DataEvent) { ev.Interface = "com.example.Other" }, want: false},
		{name: "grouped never watches", def: `{
			"name": "g", "action": {"http_url": "https://e.com", "http_method": "post"},
			"simple_triggers": [{
				"type": "data_trigger", "on": "value_change",
				"interface_name": "com.example.Sensors", "interface_major": 1,
				"match_path": "/a", "value_match_operator": "*", "group_name": "g"
			}]}`, mutate: func(*DataEvent) {}, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ev := base
			tc.mutate(&ev)
			if got := compile(t, tc.def).TracksChanges(ev); got != tc.want {
				t.Errorf("TracksChanges(%+v) = %t, want %t", ev, got, tc.want)
			}
		})
	}
}
