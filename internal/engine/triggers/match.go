// Package triggers compiles stored Astarte trigger definitions into fast
// matchers, renders the upstream-parity SimpleEvent JSON payloads, and
// executes HTTP webhook actions with retry (docs/ROADMAP.md §7.2 files
// 6.10–6.12, docs/DESIGN.md §1.1).
//
// The accepted definition shape is upstream Realm Management's trigger JSON
// (astarte_core SimpleTriggerConfig, v1.2): a "simple_triggers" array of
// data_trigger / device_trigger conditions plus one "action". Conditions
// that are valid upstream but outside Astrate's v1 evaluation scope —
// value_change*, path_created/removed, value_stored,
// device_empty_cache_received, interface_minor_updated, and group-scoped
// triggers — compile successfully (so installs round-trip) but never match;
// they are reported in Trigger.Unsupported so callers can log them.
package triggers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// Trigger condition names (upstream SimpleTriggerConfig "on" values).
const (
	// OnIncomingData fires on every accepted device data publish.
	OnIncomingData = "incoming_data"
	// OnValueChange fires when a value differs from the previous one
	// (accepted, not evaluated in v1).
	OnValueChange = "value_change"
	// OnValueChangeApplied is value_change after persistence (accepted, not
	// evaluated in v1).
	OnValueChangeApplied = "value_change_applied"
	// OnPathCreated fires on the first value of a path (accepted, not
	// evaluated in v1).
	OnPathCreated = "path_created"
	// OnPathRemoved fires on property unset of a path (accepted, not
	// evaluated in v1).
	OnPathRemoved = "path_removed"
	// OnValueStored fires after a datastream insert (accepted, not evaluated
	// in v1).
	OnValueStored = "value_stored"
	// OnDeviceConnected fires when a device session is established.
	OnDeviceConnected = "device_connected"
	// OnDeviceDisconnected fires when a device connection ends.
	OnDeviceDisconnected = "device_disconnected"
	// OnDeviceEmptyCacheReceived fires on control/emptyCache (accepted, not
	// evaluated in v1: upstream defines no SimpleEvent variant for it).
	OnDeviceEmptyCacheReceived = "device_empty_cache_received"
	// OnDeviceError fires when the validation pipeline rejects a message.
	OnDeviceError = "device_error"
	// OnIncomingIntrospection fires on every introspection publish.
	OnIncomingIntrospection = "incoming_introspection"
	// OnInterfaceAdded fires when an introspection declares a new
	// name:major pair.
	OnInterfaceAdded = "interface_added"
	// OnInterfaceRemoved fires when an introspection drops a name:major pair.
	OnInterfaceRemoved = "interface_removed"
	// OnInterfaceMinorUpdated fires on a minor bump (accepted, not evaluated
	// in v1).
	OnInterfaceMinorUpdated = "interface_minor_updated"
)

// anyToken is the wildcard accepted for device IDs, interface names, and
// value operators.
const anyToken = "*"

// anyPath is the wildcard match_path.
const anyPath = "/*"

// dataOns enumerates the data_trigger conditions; the value records whether
// this version evaluates them.
var dataOns = map[string]bool{
	OnIncomingData:       true,
	OnValueChange:        false,
	OnValueChangeApplied: false,
	OnPathCreated:        false,
	OnPathRemoved:        false,
	OnValueStored:        false,
}

// deviceOns enumerates the device_trigger conditions; the value records
// whether this version evaluates them.
var deviceOns = map[string]bool{
	OnDeviceConnected:          true,
	OnDeviceDisconnected:       true,
	OnDeviceEmptyCacheReceived: false,
	OnDeviceError:              true,
	OnIncomingIntrospection:    true,
	OnInterfaceAdded:           true,
	OnInterfaceRemoved:         true,
	OnInterfaceMinorUpdated:    false,
}

// introspectionOns are the device_trigger conditions that accept interface
// filters (upstream validate_introspection_triggers_match_conditions).
var introspectionOns = map[string]bool{
	OnInterfaceAdded:        true,
	OnInterfaceRemoved:      true,
	OnInterfaceMinorUpdated: true,
}

// valueOperators enumerates the data_trigger value_match_operator values.
var valueOperators = map[string]bool{
	anyToken: true, "==": true, "!=": true, ">": true, ">=": true,
	"<": true, "<=": true, "contains": true, "not_contains": true,
}

// Trigger is one compiled trigger: matchers plus the parsed action.
type Trigger struct {
	// Name is the trigger's installed name.
	Name string
	// Action is the parsed delivery action.
	Action *Action
	// Unsupported lists upstream-valid features this version accepts but
	// does not evaluate (logged by the engine's trigger cache).
	Unsupported []string

	data   []dataMatcher
	device []deviceMatcher
}

// DataEvent is the match input for an accepted data publish: the typed
// decoded value rides along for value conditions (nil for property unset).
type DataEvent struct {
	// DeviceID is the encoded publishing device ID.
	DeviceID string
	// Interface and Major identify the interface the publish validated
	// against.
	Interface string
	Major     int
	// Path is the concrete data path.
	Path string
	// Value is the decoded payload value (payload.Value closed set; nil for
	// unset).
	Value any
}

// DeviceEvent is the match input for a device-scoped event.
type DeviceEvent struct {
	// DeviceID is the encoded device ID.
	DeviceID string
	// On is the condition name (OnDeviceConnected, ...).
	On string
	// Interface and Major carry the interface filter input for
	// interface_added / interface_removed events; empty otherwise.
	Interface string
	Major     int
}

// dataMatcher is one compiled data_trigger condition.
type dataMatcher struct {
	deviceID  string // "" or "*" = any
	iface     string // "*" = any
	major     int    // valid when iface != "*"
	anyPath   bool
	pathSegs  []string // segment-wise; "%{...}" segments match anything
	operator  string
	known     any
	evaluated bool
}

// deviceMatcher is one compiled device_trigger condition.
type deviceMatcher struct {
	on        string
	deviceID  string // "" or "*" = any
	iface     string // interface filter for introspection conditions; "" or "*" = any
	major     *int
	evaluated bool
}

// definition is the stored trigger JSON shape.
type definition struct {
	Name           json.RawMessage       `json:"name"`
	Action         json.RawMessage       `json:"action"`
	SimpleTriggers []simpleTriggerConfig `json:"simple_triggers"`
}

// simpleTriggerConfig mirrors upstream SimpleTriggerConfig's JSON fields.
type simpleTriggerConfig struct {
	Type               string          `json:"type"`
	On                 string          `json:"on"`
	GroupName          string          `json:"group_name"`
	DeviceID           string          `json:"device_id"`
	InterfaceName      *string         `json:"interface_name"`
	InterfaceMajor     *int            `json:"interface_major"`
	MatchPath          string          `json:"match_path"`
	ValueMatchOperator string          `json:"value_match_operator"`
	KnownValue         json.RawMessage `json:"known_value"`
}

// Compile parses and validates one stored trigger definition. Validation
// follows upstream SimpleTriggerConfig: M7's Realm Management reuses it at
// install time, so an error here maps to an upstream-shaped 422.
func Compile(name string, def []byte) (*Trigger, error) {
	var d definition
	if err := json.Unmarshal(def, &d); err != nil {
		return nil, fmt.Errorf("triggers: %q does not parse: %w", name, err)
	}
	if len(d.SimpleTriggers) == 0 {
		return nil, fmt.Errorf("triggers: %q declares no simple_triggers", name)
	}
	t := &Trigger{Name: name}

	action, unsupported, err := parseAction(d.Action)
	if err != nil {
		return nil, fmt.Errorf("triggers: %q action: %w", name, err)
	}
	t.Action = action
	t.Unsupported = append(t.Unsupported, unsupported...)

	for i := range d.SimpleTriggers {
		if err := t.compileSimple(&d.SimpleTriggers[i]); err != nil {
			return nil, fmt.Errorf("triggers: %q simple_triggers[%d]: %w", name, i, err)
		}
	}
	return t, nil
}

// compileSimple validates one condition and appends its matcher.
func (t *Trigger) compileSimple(c *simpleTriggerConfig) error {
	switch c.Type {
	case "data_trigger":
		return t.compileData(c)
	case "device_trigger":
		return t.compileDevice(c)
	default:
		return fmt.Errorf("unknown trigger type %q", c.Type)
	}
}

// compileData applies upstream's data_trigger validation rules.
func (t *Trigger) compileData(c *simpleTriggerConfig) error {
	evaluated, known := dataOns[c.On]
	if !known {
		return fmt.Errorf("unknown data_trigger condition %q", c.On)
	}
	if c.InterfaceName == nil || *c.InterfaceName == "" {
		return fmt.Errorf("data_trigger requires interface_name")
	}
	iface := *c.InterfaceName
	if c.MatchPath == "" {
		return fmt.Errorf("data_trigger requires match_path")
	}
	if !valueOperators[c.ValueMatchOperator] {
		return fmt.Errorf("unknown value_match_operator %q", c.ValueMatchOperator)
	}
	if iface == anyToken {
		if c.On != OnIncomingData {
			return fmt.Errorf("interface_name %q requires on %q", anyToken, OnIncomingData)
		}
		if c.MatchPath != anyPath {
			return fmt.Errorf("interface_name %q requires match_path %q", anyToken, anyPath)
		}
	} else if c.InterfaceMajor == nil {
		return fmt.Errorf("data_trigger requires interface_major")
	}
	if c.MatchPath == anyPath && c.ValueMatchOperator != anyToken {
		return fmt.Errorf("match_path %q requires value_match_operator %q", anyPath, anyToken)
	}
	if c.ValueMatchOperator != anyToken && len(c.KnownValue) == 0 {
		return fmt.Errorf("value_match_operator %q requires known_value", c.ValueMatchOperator)
	}
	if err := validDeviceFilter(c.DeviceID, c.GroupName); err != nil {
		return err
	}

	m := dataMatcher{
		deviceID:  c.DeviceID,
		iface:     iface,
		anyPath:   c.MatchPath == anyPath,
		operator:  c.ValueMatchOperator,
		evaluated: evaluated && c.GroupName == "",
	}
	if c.InterfaceMajor != nil {
		m.major = *c.InterfaceMajor
	}
	if !m.anyPath {
		if !strings.HasPrefix(c.MatchPath, "/") {
			return fmt.Errorf("match_path %q does not start with '/'", c.MatchPath)
		}
		m.pathSegs = strings.Split(c.MatchPath, "/")[1:]
	}
	if len(c.KnownValue) > 0 {
		if err := json.Unmarshal(c.KnownValue, &m.known); err != nil {
			return fmt.Errorf("known_value does not parse: %w", err)
		}
	}
	t.noteUnsupported(!evaluated, "data_trigger on "+c.On)
	t.noteUnsupported(c.GroupName != "", "group-scoped trigger (group_name)")
	t.data = append(t.data, m)
	return nil
}

// compileDevice applies upstream's device_trigger validation rules.
func (t *Trigger) compileDevice(c *simpleTriggerConfig) error {
	evaluated, known := deviceOns[c.On]
	if !known {
		return fmt.Errorf("unknown device_trigger condition %q", c.On)
	}
	if err := validDeviceFilter(c.DeviceID, c.GroupName); err != nil {
		return err
	}
	if c.InterfaceName != nil && !introspectionOns[c.On] {
		return fmt.Errorf("interface_name is only valid on introspection conditions, not %q", c.On)
	}
	if introspectionOns[c.On] {
		if c.InterfaceName == nil {
			return fmt.Errorf("condition %q requires interface_name", c.On)
		}
		if *c.InterfaceName != anyToken && c.InterfaceMajor == nil {
			return fmt.Errorf("condition %q requires interface_major", c.On)
		}
	}

	m := deviceMatcher{
		on:        c.On,
		deviceID:  c.DeviceID,
		major:     c.InterfaceMajor,
		evaluated: evaluated && c.GroupName == "",
	}
	if c.InterfaceName != nil {
		m.iface = *c.InterfaceName
	}
	t.noteUnsupported(!evaluated, "device_trigger on "+c.On)
	t.noteUnsupported(c.GroupName != "", "group-scoped trigger (group_name)")
	t.device = append(t.device, m)
	return nil
}

// noteUnsupported records an accepted-but-not-evaluated feature once.
func (t *Trigger) noteUnsupported(cond bool, what string) {
	if !cond {
		return
	}
	for _, u := range t.Unsupported {
		if u == what {
			return
		}
	}
	t.Unsupported = append(t.Unsupported, what)
}

// validDeviceFilter checks the device_id / group_name pair: at most one may
// be set (upstream validate_device_id_xor_group_name) and device_id must be
// "*" or a valid encoded device ID.
func validDeviceFilter(deviceID, groupName string) error {
	if deviceID != "" && groupName != "" {
		return fmt.Errorf("device_id and group_name are mutually exclusive")
	}
	if deviceID != "" && deviceID != anyToken {
		if _, err := deviceid.Parse(deviceID); err != nil {
			return fmt.Errorf("device_id %q is not a valid device id", deviceID)
		}
	}
	return nil
}

// MatchesData reports whether any data_trigger condition matches the event.
func (t *Trigger) MatchesData(ev DataEvent) bool {
	for i := range t.data {
		if t.data[i].matches(ev) {
			return true
		}
	}
	return false
}

// MatchesDevice reports whether any device_trigger condition matches the
// event.
func (t *Trigger) MatchesDevice(ev DeviceEvent) bool {
	for i := range t.device {
		m := &t.device[i]
		if !m.evaluated || m.on != ev.On {
			continue
		}
		if !deviceIDMatches(m.deviceID, ev.DeviceID) {
			continue
		}
		if introspectionOns[m.on] && m.iface != anyToken {
			if m.iface != ev.Interface {
				continue
			}
			if m.major != nil && *m.major != ev.Major {
				continue
			}
		}
		return true
	}
	return false
}

// matches evaluates one data condition.
func (m *dataMatcher) matches(ev DataEvent) bool {
	if !m.evaluated {
		return false
	}
	if !deviceIDMatches(m.deviceID, ev.DeviceID) {
		return false
	}
	if m.iface != anyToken && (m.iface != ev.Interface || m.major != ev.Major) {
		return false
	}
	if !m.anyPath && !pathMatches(m.pathSegs, ev.Path) {
		return false
	}
	return valueMatches(m.operator, ev.Value, m.known)
}

// deviceIDMatches applies the device_id filter ("" and "*" match any).
func deviceIDMatches(filter, deviceID string) bool {
	return filter == "" || filter == anyToken || filter == deviceID
}

// pathMatches compares a concrete path against compiled match_path segments;
// "%{...}" placeholder segments match any single segment (upstream
// data-trigger paths may use endpoint placeholders).
func pathMatches(segs []string, path string) bool {
	if !strings.HasPrefix(path, "/") {
		return false
	}
	rest := path[1:]
	for i, seg := range segs {
		var cur string
		if i == len(segs)-1 {
			cur = rest
			rest = ""
		} else {
			var ok bool
			cur, rest, ok = strings.Cut(rest, "/")
			if !ok {
				return false
			}
		}
		if strings.Contains(cur, "/") {
			return false
		}
		if !segMatches(seg, cur) {
			return false
		}
	}
	return rest == ""
}

// segMatches compares one segment against its pattern.
func segMatches(pattern, seg string) bool {
	if strings.HasPrefix(pattern, "%{") && strings.HasSuffix(pattern, "}") {
		return seg != ""
	}
	return pattern == seg
}

// valueMatches applies a value_match_operator to the decoded payload value
// (upstream operator set; docs/ROADMAP.md §7.2 file 6.10). Values outside
// an operator's domain — ordering on strings, containment on scalars,
// anything on binary blobs or datetimes — do not match, mirroring upstream's
// guard-clause behaviour.
func valueMatches(op string, value, known any) bool {
	if op == anyToken {
		return true
	}
	if value == nil {
		return false // unset payloads only match "*"
	}
	switch op {
	case "==":
		eq, comparable := valueEquals(value, known)
		return comparable && eq
	case "!=":
		eq, comparable := valueEquals(value, known)
		return comparable && !eq
	case ">", ">=", "<", "<=":
		a, aok := numeric(value)
		b, bok := numeric(known)
		if !aok || !bok {
			return false
		}
		switch op {
		case ">":
			return a > b
		case ">=":
			return a >= b
		case "<":
			return a < b
		default:
			return a <= b
		}
	case "contains":
		return containsValue(value, known)
	case "not_contains":
		c, applicable := containable(value)
		return applicable && !c(known)
	default:
		return false
	}
}

// valueEquals compares a typed payload value with a JSON known_value; the
// second result reports whether the two were comparable at all.
func valueEquals(value, known any) (eq, comparable bool) {
	if a, ok := numeric(value); ok {
		b, ok := numeric(known)
		return ok && a == b, ok
	}
	switch v := value.(type) {
	case string:
		s, ok := known.(string)
		return ok && v == s, ok
	case bool:
		b, ok := known.(bool)
		return ok && v == b, ok
	default:
		return false, false
	}
}

// numeric widens the numeric payload and JSON types to float64.
func numeric(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case int:
		return float64(n), true
	default:
		return 0, false
	}
}

// containsValue implements the "contains" operator: substring on strings,
// element equality on arrays.
func containsValue(value, known any) bool {
	c, applicable := containable(value)
	return applicable && c(known)
}

// containable returns a containment predicate for the value, and whether
// containment applies to its type at all.
func containable(value any) (func(known any) bool, bool) {
	switch v := value.(type) {
	case string:
		return func(known any) bool {
			s, ok := known.(string)
			return ok && strings.Contains(v, s)
		}, true
	case []string:
		return sliceContains(v), true
	case []float64:
		return sliceContains(v), true
	case []int32:
		return sliceContains(v), true
	case []int64:
		return sliceContains(v), true
	case []bool:
		return sliceContains(v), true
	default:
		return nil, false
	}
}

// sliceContains builds an element-equality predicate over one homogeneous
// payload array.
func sliceContains[T any](xs []T) func(known any) bool {
	return func(known any) bool {
		for i := range xs {
			if eq, ok := valueEquals(xs[i], known); ok && eq {
				return true
			}
		}
		return false
	}
}
