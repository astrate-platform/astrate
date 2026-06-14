package triggers

import (
	"strings"
	"testing"
)

// fixture trigger definitions, shaped exactly like upstream Realm
// Management's trigger JSON (docs/ROADMAP.md §7.3 "fixture upstream trigger
// JSONs").
const (
	fixtureValuesAbove = `{
		"name": "values_above",
		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
		"simple_triggers": [{
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "org.astarte-platform.genericsensors.Values",
			"interface_major": 1,
			"match_path": "/streamTest/value",
			"value_match_operator": ">",
			"known_value": 0.4
		}]
	}`

	fixtureAnyData = `{
		"name": "any_data",
		"action": {"http_post_url": "https://example.com/legacy"},
		"simple_triggers": [{
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "*",
			"match_path": "/*",
			"value_match_operator": "*"
		}]
	}`

	fixtureParametricPath = `{
		"name": "parametric",
		"action": {"http_url": "https://example.com/hook", "http_method": "put"},
		"simple_triggers": [{
			"type": "data_trigger",
			"on": "incoming_data",
			"interface_name": "org.astarte-platform.genericsensors.Values",
			"interface_major": 1,
			"match_path": "/%{sensor_id}/value",
			"value_match_operator": "*"
		}]
	}`

	fixtureDeviceConnected = `{
		"name": "connected",
		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
	}`

	fixtureOneDevice = `{
		"name": "one_device",
		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
		"simple_triggers": [{
			"type": "device_trigger",
			"on": "device_disconnected",
			"device_id": "f0VMRgIBAQAAAAAAAAAAAA"
		}]
	}`

	fixtureInterfaceAdded = `{
		"name": "iface_added",
		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
		"simple_triggers": [{
			"type": "device_trigger",
			"on": "interface_added",
			"interface_name": "com.example.Sensors",
			"interface_major": 2
		}]
	}`

	fixtureValueChange = `{
		"name": "value_change_accepted",
		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
		"simple_triggers": [{
			"type": "data_trigger",
			"on": "value_change",
			"interface_name": "com.example.Sensors",
			"interface_major": 1,
			"match_path": "/a",
			"value_match_operator": "*"
		}]
	}`
)

// compile is the test shorthand.
func compile(t *testing.T, def string) *Trigger {
	t.Helper()
	tr, err := Compile("t", []byte(def))
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}
	return tr
}

// dataTrigger builds a single-condition data trigger definition.
func dataTrigger(t *testing.T, condition string) *Trigger {
	t.Helper()
	return compile(t, `{
		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
		"simple_triggers": [`+condition+`]
	}`)
}

// valuesEvent is the baseline matching event for fixtureValuesAbove.
func valuesEvent() DataEvent {
	return DataEvent{
		DeviceID:  "f0VMRgIBAQAAAAAAAAAAAA",
		Interface: "org.astarte-platform.genericsensors.Values",
		Major:     1,
		Path:      "/streamTest/value",
		Value:     0.5,
	}
}

// TestMatchData is the docs/ROADMAP.md §7.3 match/no-match table over the
// upstream fixtures.
func TestMatchData(t *testing.T) {
	above := compile(t, fixtureValuesAbove)
	anyData := compile(t, fixtureAnyData)
	parametric := compile(t, fixtureParametricPath)

	cases := []struct {
		name    string
		trigger *Trigger
		mutate  func(*DataEvent)
		want    bool
	}{
		{name: "value above threshold", trigger: above, mutate: func(*DataEvent) {}, want: true},
		{name: "value below threshold", trigger: above,
			mutate: func(ev *DataEvent) { ev.Value = 0.3 }, want: false},
		{name: "value equal to threshold", trigger: above,
			mutate: func(ev *DataEvent) { ev.Value = 0.4 }, want: false},
		{name: "integer value widens", trigger: above,
			mutate: func(ev *DataEvent) { ev.Value = int32(7) }, want: true},
		{name: "wrong interface", trigger: above,
			mutate: func(ev *DataEvent) { ev.Interface = "com.example.Other" }, want: false},
		{name: "wrong major", trigger: above,
			mutate: func(ev *DataEvent) { ev.Major = 0 }, want: false},
		{name: "wrong path", trigger: above,
			mutate: func(ev *DataEvent) { ev.Path = "/streamTest/other" }, want: false},
		{name: "unset only matches any", trigger: above,
			mutate: func(ev *DataEvent) { ev.Value = nil }, want: false},
		{name: "string value does not order-compare", trigger: above,
			mutate: func(ev *DataEvent) { ev.Value = "high" }, want: false},

		{name: "catch-all matches anything", trigger: anyData,
			mutate: func(ev *DataEvent) { ev.Interface = "com.whatever.Iface"; ev.Path = "/x/y" }, want: true},
		{name: "catch-all matches unset", trigger: anyData,
			mutate: func(ev *DataEvent) { ev.Value = nil }, want: true},

		{name: "placeholder segment matches", trigger: parametric,
			mutate: func(ev *DataEvent) { ev.Path = "/sensor7/value" }, want: true},
		{name: "placeholder needs the suffix", trigger: parametric,
			mutate: func(ev *DataEvent) { ev.Path = "/sensor7/other" }, want: false},
		{name: "placeholder is one segment", trigger: parametric,
			mutate: func(ev *DataEvent) { ev.Path = "/a/b/value" }, want: false},
		{name: "placeholder path too short", trigger: parametric,
			mutate: func(ev *DataEvent) { ev.Path = "/value" }, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ev := valuesEvent()
			tc.mutate(&ev)
			if got := tc.trigger.MatchesData(ev); got != tc.want {
				t.Errorf("MatchesData(%+v) = %t, want %t", ev, got, tc.want)
			}
		})
	}
}

// TestMatchDataOperators covers every value_match_operator across the typed
// value domain.
func TestMatchDataOperators(t *testing.T) {
	condition := func(op, known string) string {
		return `{
			"type": "data_trigger", "on": "incoming_data",
			"interface_name": "com.example.A", "interface_major": 1,
			"match_path": "/v", "value_match_operator": "` + op + `",
			"known_value": ` + known + `}`
	}
	ev := func(v any) DataEvent {
		return DataEvent{DeviceID: "d", Interface: "com.example.A", Major: 1, Path: "/v", Value: v}
	}

	cases := []struct {
		name  string
		op    string
		known string
		value any
		want  bool
	}{
		{name: "eq number", op: "==", known: "5", value: int64(5), want: true},
		{name: "eq number mismatch", op: "==", known: "5", value: int64(6), want: false},
		{name: "eq string", op: "==", known: `"on"`, value: "on", want: true},
		{name: "eq bool", op: "==", known: "true", value: true, want: true},
		{name: "eq type mismatch never matches", op: "==", known: `"5"`, value: int64(5), want: false},
		{name: "neq number", op: "!=", known: "5", value: int64(6), want: true},
		{name: "neq equal value", op: "!=", known: "5", value: int64(5), want: false},
		{name: "gte boundary", op: ">=", known: "5", value: 5.0, want: true},
		{name: "lt", op: "<", known: "5", value: 4.5, want: true},
		{name: "lte above", op: "<=", known: "5", value: 5.5, want: false},
		{name: "contains substring", op: "contains", known: `"err"`, value: "fatal error", want: true},
		{name: "contains missing substring", op: "contains", known: `"err"`, value: "all good", want: false},
		{name: "contains array element", op: "contains", known: "3", value: []int32{1, 2, 3}, want: true},
		{name: "contains missing element", op: "contains", known: "9", value: []int32{1, 2, 3}, want: false},
		{name: "contains string array", op: "contains", known: `"b"`, value: []string{"a", "b"}, want: true},
		{name: "contains on scalar number", op: "contains", known: "1", value: 1.0, want: false},
		{name: "not_contains array", op: "not_contains", known: "9", value: []int32{1, 2}, want: true},
		{name: "not_contains present element", op: "not_contains", known: "1", value: []int32{1, 2}, want: false},
		{name: "not_contains on scalar number", op: "not_contains", known: "1", value: 2.0, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tr := dataTrigger(t, condition(tc.op, tc.known))
			if got := tr.MatchesData(ev(tc.value)); got != tc.want {
				t.Errorf("op %s known %s value %v = %t, want %t", tc.op, tc.known, tc.value, got, tc.want)
			}
		})
	}
}

// TestMatchDevice covers the device_trigger conditions and their filters.
func TestMatchDevice(t *testing.T) {
	connected := compile(t, fixtureDeviceConnected)
	oneDevice := compile(t, fixtureOneDevice)
	added := compile(t, fixtureInterfaceAdded)

	cases := []struct {
		name    string
		trigger *Trigger
		ev      DeviceEvent
		want    bool
	}{
		{name: "connected matches any device", trigger: connected,
			ev: DeviceEvent{DeviceID: "x", On: OnDeviceConnected}, want: true},
		{name: "connected ignores disconnects", trigger: connected,
			ev: DeviceEvent{DeviceID: "x", On: OnDeviceDisconnected}, want: false},
		{name: "device filter matches", trigger: oneDevice,
			ev: DeviceEvent{DeviceID: "f0VMRgIBAQAAAAAAAAAAAA", On: OnDeviceDisconnected}, want: true},
		{name: "device filter rejects others", trigger: oneDevice,
			ev: DeviceEvent{DeviceID: "AAAAAAAAAAAAAAAAAAAAAQ", On: OnDeviceDisconnected}, want: false},
		{name: "interface added with matching filter", trigger: added,
			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.Sensors", Major: 2}, want: true},
		{name: "interface added wrong major", trigger: added,
			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.Sensors", Major: 1}, want: false},
		{name: "interface added wrong interface", trigger: added,
			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.Other", Major: 2}, want: false},
		{name: "interface added ignores removals", trigger: added,
			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceRemoved, Interface: "com.example.Sensors", Major: 2}, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.trigger.MatchesDevice(tc.ev); got != tc.want {
				t.Errorf("MatchesDevice(%+v) = %t, want %t", tc.ev, got, tc.want)
			}
		})
	}
}

// TestCompileAcceptsUnsupported: upstream-valid-but-unevaluated conditions
// compile (installs round-trip), report themselves, and never match.
func TestCompileAcceptsUnsupported(t *testing.T) {
	tr := compile(t, fixtureValueChange)
	if len(tr.Unsupported) == 0 {
		t.Fatal("value_change compiled without an Unsupported report")
	}
	if tr.MatchesData(DataEvent{Interface: "com.example.Sensors", Major: 1, Path: "/a", Value: 1.0}) {
		t.Error("unevaluated condition matched")
	}

	// Group-scoped triggers compile and never match.
	grouped := dataTrigger(t, `{
		"type": "data_trigger", "on": "incoming_data",
		"interface_name": "com.example.A", "interface_major": 1,
		"match_path": "/v", "value_match_operator": "*",
		"group_name": "fleet-a"}`)
	if len(grouped.Unsupported) == 0 {
		t.Error("group-scoped trigger compiled without an Unsupported report")
	}
	if grouped.MatchesData(DataEvent{DeviceID: "d", Interface: "com.example.A", Major: 1, Path: "/v", Value: 1.0}) {
		t.Error("group-scoped trigger matched")
	}
}

// TestCompileRejects is the upstream-parity validation table.
func TestCompileRejects(t *testing.T) {
	action := `"action": {"http_url": "https://example.com", "http_method": "post"}`
	cases := []struct {
		name string
		def  string
		want string
	}{
		{name: "not json", def: `{`, want: "does not parse"},
		{name: "no simple triggers", def: `{` + action + `, "simple_triggers": []}`, want: "no simple_triggers"},
		{name: "missing action", def: `{"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]}`,
			want: "missing action"},
		{name: "bad http method", def: `{"action": {"http_url": "https://e.com", "http_method": "yeet"},
			"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]}`,
			want: "unsupported http_method"},
		{name: "unknown trigger type", def: `{` + action + `, "simple_triggers": [{"type": "mystery"}]}`,
			want: "unknown trigger type"},
		{name: "unknown data condition", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "data_happened", "interface_name": "a.B",
			"interface_major": 0, "match_path": "/*", "value_match_operator": "*"}]}`,
			want: "unknown data_trigger condition"},
		{name: "data without interface", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data",
			"match_path": "/*", "value_match_operator": "*"}]}`,
			want: "requires interface_name"},
		{name: "data without major", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "a.B",
			"match_path": "/v", "value_match_operator": "*"}]}`,
			want: "requires interface_major"},
		{name: "any interface with concrete path", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "*",
			"match_path": "/v", "value_match_operator": "*"}]}`,
			want: `requires match_path "/*"`},
		{name: "any path with operator", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "a.B",
			"interface_major": 0, "match_path": "/*", "value_match_operator": ">",
			"known_value": 1}]}`,
			want: `requires value_match_operator "*"`},
		{name: "operator without known value", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "a.B",
			"interface_major": 0, "match_path": "/v", "value_match_operator": ">"}]}`,
			want: "requires known_value"},
		{name: "unknown operator", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "a.B",
			"interface_major": 0, "match_path": "/v", "value_match_operator": "~="}]}`,
			want: "unknown value_match_operator"},
		{name: "device id and group name", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "device_connected",
			"device_id": "f0VMRgIBAQAAAAAAAAAAAA", "group_name": "g"}]}`,
			want: "mutually exclusive"},
		{name: "invalid device id", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "device_connected", "device_id": "nope!"}]}`,
			want: "not a valid device id"},
		{name: "unknown device condition", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "device_exploded"}]}`,
			want: "unknown device_trigger condition"},
		{name: "interface filter on connect", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "device_connected", "interface_name": "a.B"}]}`,
			want: "only valid on introspection conditions"},
		{name: "interface_added without interface", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "interface_added"}]}`,
			want: "requires interface_name"},
		{name: "interface_added without major", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "interface_added", "interface_name": "a.B"}]}`,
			want: "requires interface_major"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Compile("t", []byte(tc.def))
			if err == nil {
				t.Fatalf("Compile accepted %s", tc.name)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("err = %v, want substring %q", err, tc.want)
			}
		})
	}
}

// TestCompileActions: both upstream HTTP action shapes parse; custom
// actions land on the Forwarder seam.
func TestCompileActions(t *testing.T) {
	modern := compile(t, fixtureValuesAbove)
	if modern.Action.Method != "POST" || modern.Action.URL != "https://example.com/hook" || modern.Action.Custom != nil {
		t.Errorf("modern action: %+v", modern.Action)
	}

	legacy := compile(t, fixtureAnyData)
	if legacy.Action.Method != "POST" || legacy.Action.URL != "https://example.com/legacy" {
		t.Errorf("legacy http_post_url action: %+v", legacy.Action)
	}

	custom := compile(t, `{
		"action": {"amqp_exchange": "astarte_events_test", "amqp_routing_key": "k"},
		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
	}`)
	if custom.Action.Custom == nil {
		t.Error("AMQP-style action did not land on the Forwarder seam")
	}

	templated := compile(t, `{
		"action": {"http_url": "https://example.com", "http_method": "post",
			"template_type": "mustache", "template": "{{ value }}"},
		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
	}`)
	if len(templated.Unsupported) == 0 {
		t.Error("mustache template accepted silently")
	}

	headers := compile(t, `{
		"action": {"http_url": "https://example.com", "http_method": "delete",
			"http_static_headers": {"X-Custom": "v"}, "ignore_ssl_errors": true},
		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
	}`)
	if headers.Action.Method != "DELETE" || headers.Action.StaticHeaders["X-Custom"] != "v" || !headers.Action.IgnoreSSLErrors {
		t.Errorf("action attributes: %+v", headers.Action)
	}
}
