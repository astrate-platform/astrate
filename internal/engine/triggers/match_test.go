package triggers

import (
	"errors"
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
		On:        OnIncomingData,
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
		return DataEvent{DeviceID: "d", On: OnIncomingData, Interface: "com.example.A", Major: 1, Path: "/v", Value: v}
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
		{name: "deletion_started matches its on", trigger: compile(t, `{
			"name": "del", "action": {"http_post_url": "https://x"},
			"simple_triggers": [{"type": "device_trigger", "on": "device_deletion_started"}]
		}`),
			ev: DeviceEvent{DeviceID: "x", On: OnDeviceDeletionStarted}, want: true},
		{name: "deletion_finished matches its on", trigger: compile(t, `{
			"name": "del", "action": {"http_post_url": "https://x"},
			"simple_triggers": [{"type": "device_trigger", "on": "device_deletion_finished"}]
		}`),
			ev: DeviceEvent{DeviceID: "x", On: OnDeviceDeletionFinished}, want: true},
		{name: "deletion_started ignores finished", trigger: compile(t, `{
			"name": "del", "action": {"http_post_url": "https://x"},
			"simple_triggers": [{"type": "device_trigger", "on": "device_deletion_started"}]
		}`),
			ev: DeviceEvent{DeviceID: "x", On: OnDeviceDeletionFinished}, want: false},
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
// compile (installs round-trip), report themselves, and never match. The
// change-derived data conditions are evaluated since the previous-value
// activation; only group-scoping and the two dormant device conditions
// remain unevaluated.
func TestCompileAcceptsUnsupported(t *testing.T) {
	// Group-scoped triggers compile and never match.
	grouped := dataTrigger(t, `{
		"type": "data_trigger", "on": "incoming_data",
		"interface_name": "com.example.A", "interface_major": 1,
		"match_path": "/v", "value_match_operator": "*",
		"group_name": "fleet-a"}`)
	if len(grouped.Unsupported) == 0 {
		t.Error("group-scoped trigger compiled without an Unsupported report")
	}
	if grouped.MatchesData(DataEvent{DeviceID: "d", On: OnIncomingData, Interface: "com.example.A", Major: 1, Path: "/v", Value: 1.0}) {
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
			want: "http_method=is invalid"},
		{name: "unknown trigger type", def: `{` + action + `, "simple_triggers": [{"type": "mystery"}]}`,
			want: "unknown trigger type"},
		{name: "unknown data condition", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "data_happened", "interface_name": "a.B",
			"interface_major": 0, "match_path": "/*", "value_match_operator": "*"}]}`,
			want: "on=is invalid"},
		{name: "data without interface", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data",
			"match_path": "/*", "value_match_operator": "*"}]}`,
			want: "interface_name=can't be blank"},
		{name: "data without major", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "a.B",
			"match_path": "/v", "value_match_operator": "*"}]}`,
			want: "interface_major=can't be blank"},
		{name: "any interface with value_change", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "value_change", "interface_name": "*",
			"match_path": "/*", "value_match_operator": "*"}]}`,
			want: `on=must be incoming_data when interface_name is *`},
		{name: "any interface with concrete path", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "*",
			"match_path": "/v", "value_match_operator": "*"}]}`,
			want: `match_path=must be /* when interface_name is *`},
		{name: "any path with operator", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "a.B",
			"interface_major": 0, "match_path": "/*", "value_match_operator": ">",
			"known_value": 1}]}`,
			want: `value_match_operator=must be * when match_path is /*`},
		{name: "operator without known value", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "a.B",
			"interface_major": 0, "match_path": "/v", "value_match_operator": ">"}]}`,
			want: "known_value=can't be blank"},
		{name: "unknown operator", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "a.B",
			"interface_major": 0, "match_path": "/v", "value_match_operator": "~="}]}`,
			want: "known_value=can't be blank, value_match_operator=is invalid"},
		{name: "bad interface format leading digit", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "1bad.Name",
			"interface_major": 0, "match_path": "/v", "value_match_operator": "*"}]}`,
			want: "interface_name=has invalid format"},
		{name: "bad interface format leading hyphen", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "-bad",
			"interface_major": 0, "match_path": "/v", "value_match_operator": "*"}]}`,
			want: "interface_name=has invalid format"},
		{name: "bad interface format empty label", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "a..b",
			"interface_major": 0, "match_path": "/v", "value_match_operator": "*"}]}`,
			want: "interface_name=has invalid format"},
		{name: "bad path trailing slash", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "a.B",
			"interface_major": 0, "match_path": "/value/", "value_match_operator": "*"}]}`,
			want: "match_path=has invalid format"},
		{name: "bad path empty placeholder", def: `{` + action + `, "simple_triggers": [{
			"type": "data_trigger", "on": "incoming_data", "interface_name": "a.B",
			"interface_major": 0, "match_path": "%{}/x", "value_match_operator": "*"}]}`,
			want: "match_path=has invalid format"},
		// Upstream's mapping regex admits uppercase segments ([a-zA-Z]):
		// "/Value/x" is VALID, unlike a naive lowercase-only reading.
		{name: "device id and group name", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "device_connected",
			"device_id": "f0VMRgIBAQAAAAAAAAAAAA", "group_name": "g"}]}`,
			want: "group_name=must not be defined if device_id is defined"},
		{name: "invalid device id", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "device_connected", "device_id": "nope!"}]}`,
			want: "device_id=is not a valid device id"},
		{name: "extended device id", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "device_connected", "device_id": "` +
			strings.Repeat("A", 43) + `"}]}`,
			want: "device_id=is too long, device id must be 128 bits"},
		{name: "unknown device condition", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "device_exploded"}]}`,
			want: "on=is invalid"},
		{name: "interface filter on connect", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "device_connected", "interface_name": "a.B"}]}`,
			want: "is allowed only in if 'on' is one of interface_minor_updated, interface_removed, interface_added"},
		{name: "interface_added without interface", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "interface_added"}]}`,
			want: "interface_major=can't be blank"},
		{name: "interface_added without major", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "interface_added", "interface_name": "a.B"}]}`,
			want: "interface_major=can't be blank"},
		{name: "interface_removed without major", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "interface_removed", "interface_name": "a.B"}]}`,
			want: "interface_major=can't be blank"},
		{name: "minor_updated without name", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "interface_minor_updated", "interface_major": 1}]}`,
			want: "interface_name=must be set in interface_minor_updated triggers"},
		{name: "minor_updated wildcard name", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "interface_minor_updated", "interface_name": "*",
			"interface_major": 1}]}`,
			want: `interface_name=must not be '*' in interface_minor_updated triggers`},
		{name: "incoming_introspection with name", def: `{` + action + `, "simple_triggers": [{
			"type": "device_trigger", "on": "incoming_introspection",
			"interface_name": "a.B", "interface_major": 1}]}`,
			want: "must not be set in incoming_introspection triggers"},
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

// TestCompileFieldErrorShape pins the aggregate error shape end to end: a
// clean first simple_trigger renders as an empty entry ({} on the wire), a
// "/*"+"==" definition demands only the operator fix even though known_value
// would be required for a concrete path (upstream's probe-frozen combo).
func TestCompileFieldErrorShape(t *testing.T) {
	def := `{
		"action": {"http_url": "https://example.com", "http_method": "post"},
		"simple_triggers": [
			{"type": "device_trigger", "on": "device_connected"},
			{"type": "data_trigger", "on": "incoming_data", "interface_name": "a.B",
			 "interface_major": 0, "match_path": "/*", "value_match_operator": "==",
			 "known_value": 1}
		]
	}`
	_, err := Compile("t", []byte(def))
	if err == nil {
		t.Fatal("Compile accepted wildcard path with ==")
	}
	want := `triggers: "t": simple_triggers[1]: value_match_operator=must be * when match_path is /*`
	if err.Error() != want {
		t.Errorf("err = %q, want %q", err.Error(), want)
	}
	var te *TriggerErrors
	if !errors.As(err, &te) {
		t.Fatalf("err %v is not *TriggerErrors", err)
	}
	if te.Action != nil || len(te.SimpleTriggers) != 2 || te.SimpleTriggers[0] != nil {
		t.Errorf("aggregate = %+v, want clean action, aligned length 2, empty entry 0", te)
	}
}

// TestCompileDeviceIntrospectionMatrix covers rule set 2's accept side:
// name-less added/removed triggers need (and keep) their major, "*" ignores
// it, minor_updated wants a concrete name, stray majors are accepted.
func TestCompileDeviceIntrospectionMatrix(t *testing.T) {
	defDevice := func(condition string) string {
		return `{
			"action": {"http_url": "https://example.com/hook", "http_method": "post"},
			"simple_triggers": [` + condition + `]
		}`
	}

	namelessAdded := compile(t, defDevice(
		`{"type": "device_trigger", "on": "interface_added", "interface_major": 2}`))
	namelessRemoved := compile(t, defDevice(
		`{"type": "device_trigger", "on": "interface_removed", "interface_major": 3}`))
	wildcardAdded := compile(t, defDevice(
		`{"type": "device_trigger", "on": "interface_added", "interface_name": "*", "interface_major": 7}`))
	minorUpdated := compile(t, defDevice(
		`{"type": "device_trigger", "on": "interface_minor_updated",
		  "interface_name": "com.example.Sensors", "interface_major": 1}`))
	connectWithMajor := compile(t, defDevice(
		`{"type": "device_trigger", "on": "device_connected", "interface_major": 9}`))

	if len(minorUpdated.Unsupported) == 0 {
		t.Error("minor_updated should stay unevaluated-but-accepted")
	}

	cases := []struct {
		name    string
		trigger *Trigger
		ev      DeviceEvent
		want    bool
	}{
		{name: "name-less added matches any interface of the major", trigger: namelessAdded,
			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.Other", Major: 2}, want: true},
		{name: "name-less added rejects other majors", trigger: namelessAdded,
			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.Other", Major: 1}, want: false},
		{name: "name-less removed matches any interface of the major", trigger: namelessRemoved,
			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceRemoved, Interface: "org.a.B", Major: 3}, want: true},
		{name: "name-less removed rejects other majors", trigger: namelessRemoved,
			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceRemoved, Interface: "org.a.B", Major: 4}, want: false},
		{name: "wildcard name matches any interface", trigger: wildcardAdded,
			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.X", Major: 9}, want: true},
		{name: "wildcard name ignores the major", trigger: wildcardAdded,
			ev: DeviceEvent{DeviceID: "x", On: OnInterfaceAdded, Interface: "com.example.X", Major: 123}, want: true},
		{name: "connect with major-only stays unfiltered", trigger: connectWithMajor,
			ev: DeviceEvent{DeviceID: "x", On: OnDeviceConnected}, want: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.trigger.MatchesDevice(tc.ev); got != tc.want {
				t.Errorf("MatchesDevice(%+v) = %t, want %t", tc.ev, got, tc.want)
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

	// Non-amqp custom action (#64: amqp_exchange actions are rejected at
	// compile time); it must still land on the Forwarder seam.
	custom := compile(t, `{
		"action": {"nats_subject": "astarte_events_test"},
		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
	}`)
	if custom.Action.Custom == nil {
		t.Error("custom action did not land on the Forwarder seam")
	}

	templated := compile(t, `{
		"action": {"http_url": "https://example.com", "http_method": "post",
			"template_type": "mustache", "template": "{{ value }}"},
		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
	}`)
	if len(templated.Unsupported) != 0 {
		t.Errorf("mustache template should render, not be marked unsupported: %v", templated.Unsupported)
	}
	if templated.Action.Template != "{{ value }}" {
		t.Errorf("template not carried through: %+v", templated.Action)
	}

	unknownTemplate := compile(t, `{
		"action": {"http_url": "https://example.com", "http_method": "post",
			"template_type": "jinja2", "template": "{{ value }}"},
		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
	}`)
	if len(unknownTemplate.Unsupported) == 0 {
		t.Error("non-mustache template_type accepted silently")
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

func TestCompilePolicyField(t *testing.T) {
	withPolicy := compile(t, `{
		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
		"policy": "my_policy",
		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
	}`)
	if withPolicy.PolicyName != "my_policy" {
		t.Errorf("PolicyName = %q, want %q", withPolicy.PolicyName, "my_policy")
	}
	if withPolicy.Policy() != nil {
		t.Error("Policy() should be nil before AttachPolicy")
	}

	withoutPolicy := compile(t, `{
		"action": {"http_url": "https://example.com/hook", "http_method": "post"},
		"simple_triggers": [{"type": "device_trigger", "on": "device_connected"}]
	}`)
	if withoutPolicy.PolicyName != "" {
		t.Errorf("PolicyName = %q, want empty", withoutPolicy.PolicyName)
	}

	withPolicy.AttachPolicy(&Policy{Name: "my_policy"})
	if withPolicy.Policy() == nil || withPolicy.Policy().Name != "my_policy" {
		t.Error("Policy() should return attached policy")
	}
}
