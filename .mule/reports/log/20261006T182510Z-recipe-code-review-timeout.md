slug: recipe-code-review
verdict: timeout
at:  d567596
ran: 2026-10-06T18:25:10Z on DietPi in 1200s


> build · big-pickle

→ Read .mule/MULE.md [offset=1, limit=2000]
→ Read .mule/recipes/code-review.md [offset=1, limit=4000]

$ ls internal/ pkg/
internal/:
appengine
auth
broker
config
engine
flow
flowapi
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil

pkg/:
astarteapi
deviceid
interfaceschema
payload

$ rg -c '' --glob 'internal/engine/**/*.go' | sort -t: -k2 -rn | head -10
internal/engine/triggers/actions_test.go:1199
internal/engine/fakestore_test.go:818
internal/engine/triggers/match.go:769
internal/engine/triggers/actions.go:629
internal/engine/triggers/match_test.go:594
internal/engine/data.go:552
internal/engine/cache.go:502
internal/engine/router.go:490
internal/engine/data_test.go:461
internal/engine/batch.go:459

Let me read the triggers implementation and tests to understand the codebase.
$ sed -n '1,120p' internal/engine/triggers/match.go
// Package triggers compiles stored Astarte trigger definitions into fast
// matchers, renders the upstream-parity SimpleEvent JSON payloads, and
// executes HTTP webhook actions with retry (docs/ROADMAP.md §7.2 files
// 6.10–6.12, docs/DESIGN.md §1.1).
//
// The accepted definition shape is upstream Realm Management's trigger JSON
// (astarte_core SimpleTriggerConfig, v1.2): a "simple_triggers" array of
// data_trigger / device_trigger conditions plus one "action". Conditions
// that are valid upstream but outside Astrate's v1 evaluation scope —
// device_empty_cache_received, interface_minor_updated, and group-scoped
// triggers — compile successfully (so installs round-trip) but never match;
// they are reported in Trigger.Unsupported so callers can log them.
//
// The change-derived data conditions (value_change, value_change_applied,
// path_created, path_removed, value_stored) follow upstream data_handler.ex
// semantics; the engine captures the previous-value snapshot pre-write and
// emits the events post-commit (internal/engine trackPrevious/fireData).
package triggers

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// Trigger condition names (upstream SimpleTriggerConfig "on" values).
const (
	// OnIncomingData fires on every accepted device data publish.
	OnIncomingData = "incoming_data"
	// OnValueChange fires when a value differs from the previous one
	// (upstream execute_pre_change_triggers; Astrate evaluates it
	// post-commit against the accept-time previous-value snapshot).
	OnValueChange = "value_change"
	// OnValueChangeApplied is value_change after persistence (upstream
	// execute_post_change_triggers).
	OnValueChangeApplied = "value_change_applied"
	// OnPathCreated fires on the first value of a path (previous missing,
	// new value present — upstream Core.Trigger post-change semantics).
	OnPathCreated = "path_created"
	// OnPathRemoved fires when an existing property is unset (previous
	// present, new value absent).
	OnPathRemoved = "path_removed"
	// OnValueStored fires after an accepted individual-datastream insert.
	OnValueStored = "value_stored"
	// OnDeviceRegistered fires when a device is registered via the Pairing API.
	OnDeviceRegistered = "device_registered"
	// OnDeviceConnected fires when a device session is established.
	OnDeviceConnected = "device_connected"
	// OnDeviceDisconnected fires when a device connection ends.
	OnDeviceDisconnected = "device_disconnected"
	// OnDeviceDeletionStarted fires immediately before a device is deleted
	// (upstream: start of async deletion; Astrate: before the sync delete).
	OnDeviceDeletionStarted = "device_deletion_started"
	// OnDeviceDeletionFinished fires immediately after a device is deleted
	// (upstream: end of async deletion; Astrate: after the sync delete).
	OnDeviceDeletionFinished = "device_deletion_finished"
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

// upstream v1.2.0 formats, probe-frozen (#70): Interface.interface_name_regex
// and Mapping.mapping_regex. Concrete (non-wildcard) interface_name and
// match_path values must match these or the definition is rejected.
var (
	ifaceNameRE = regexp.MustCompile(`^([a-zA-Z][a-zA-Z0-9]*\.([a-zA-Z0-9][a-zA-Z0-9-]*\.)*)?[a-zA-Z][a-zA-Z0-9]*$`)
	matchPathRE = regexp.MustCompile(`^(\/(%{([a-zA-Z_]+[a-zA-Z0-9_]*)}|[a-zA-Z_]+[a-zA-Z0-9_]*)){1,64}$`)
)

// dataOns enumerates the data_trigger conditions; the value records whether
// this version evaluates them.
var dataOns = map[string]bool{
	OnIncomingData:       true,
	OnValueChange:        true,
	OnValueChangeApplied: true,
	OnPathCreated:        true,
	OnPathRemoved:        true,
	OnValueStored:        true,
}

// prevOns are the data conditions decided by the previous-value comparison:
// they need the engine's pre-write snapshot (upstream get_value_change_
// triggers). value_stored is deliberately absent — it fires on every
// accepted insert and needs no lookup.
var prevOns = map[string]bool{
	OnValueChange:        true,
	OnValueChangeApplied: true,
	OnPathCreated:        true,
	OnPathRemoved:        true,
}

// deviceOns enumerates the device_trigger conditions; the value records
// whether this version evaluates them.
var deviceOns = map[string]bool{
	OnDeviceRegistered:         true,
	OnDeviceConnected:          true,
	OnDeviceDisconnected:       true,

$ sed -n '120,240p' internal/engine/triggers/match.go
OnDeviceDisconnected:       true,
	OnDeviceDeletionStarted:    true,
	OnDeviceDeletionFinished:   true,
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
	// PolicyName is the delivery policy named by the trigger definition.
	// An empty value means the implicit @default behaviour.
	PolicyName string

	data   []dataMatcher
	device []deviceMatcher
	policy *Policy
}

// AttachPolicy sets the resolved delivery policy on the trigger.
func (t *Trigger) AttachPolicy(p *Policy) { t.policy = p }

// Policy returns the resolved delivery policy, nil when none is attached.
// Safe on a nil receiver.
func (t *Trigger) Policy() *Policy {
	if t == nil {
		return nil
	}
	return t.policy
}

// DataEvent is the match input for an accepted data publish: the typed
// decoded value rides along for value conditions (nil for property unset).
// On selects which condition set evaluates: one of the On* data constants
// (incoming_data, value_change, value_change_applied, path_created,
// path_removed, value_stored) — a matcher only reacts to its own condition.
type DataEvent struct {
	// DeviceID is the encoded publishing device ID.
	DeviceID string
	// On is the condition name being evaluated.
	On string
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
	on        string // the condition this matcher reacts to
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
	Policy         string                `json:"policy"`
	SimpleTriggers []simpleTriggerConfig `json:"simple_triggers"`
}

// simpleTriggerConfig mirrors upstream SimpleTriggerConfig's JSON fields.
type simpleTriggerConfig struct {
	Type               string          `json:"type"`

$ sed -n '240,360p' internal/engine/triggers/match.go
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
//
// Field-scoped violations accumulate into a *TriggerErrors (action plus an
// index-aligned simple_triggers array, upstream changeset semantics); errors
// that have no field shape (does-not-parse, no simple_triggers, unknown
// trigger type) stay plain wrapped errors.
func Compile(name string, def []byte) (*Trigger, error) {
	var d definition
	if err := json.Unmarshal(def, &d); err != nil {
		return nil, fmt.Errorf("triggers: %q does not parse: %w", name, err)
	}
	if len(d.SimpleTriggers) == 0 {
		return nil, fmt.Errorf("triggers: %q declares no simple_triggers", name)
	}
	t := &Trigger{Name: name, PolicyName: d.Policy}
	agg := &TriggerErrors{}

	action, unsupported, err := parseAction(d.Action)
	if err != nil {
		var fe *FieldErrors
		if errors.As(err, &fe) {
			agg.Action = fe.Fields
		} else {
			return nil, fmt.Errorf("triggers: %q action: %w", name, err)
		}
	} else {
		t.Action = action
		t.Unsupported = append(t.Unsupported, unsupported...)
	}

	for i := range d.SimpleTriggers {
		if err := t.compileSimple(&d.SimpleTriggers[i]); err != nil {
			var fe fieldErrs
			if !errors.As(err, &fe) {
				return nil, fmt.Errorf("triggers: %q simple_triggers[%d]: %w", name, i, err)
			}
			for len(agg.SimpleTriggers) < i {
				agg.SimpleTriggers = append(agg.SimpleTriggers, nil)
			}
			agg.SimpleTriggers = append(agg.SimpleTriggers, fe)
		}
	}
	if agg.Action != nil || len(agg.SimpleTriggers) > 0 {
		return nil, fmt.Errorf("triggers: %q: %w", name, agg)
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

// compileData applies upstream's data_trigger validation rules (v1.2.0
// simple_trigger_config.ex). Violations accumulate into a fieldErrs error so
// Compile renders them as the per-index changeset envelope; every accepted
// definition compiles to exactly the same matcher it did before #70.
func (t *Trigger) compileData(c *simpleTriggerConfig) error {
	fe := fieldErrs{}
	_, known := dataOns[c.On]
	if c.On == "" {
		fe.add("on", "can't be blank")
	} else if !known {
		fe.add("on", "is invalid")
	}
	iface := ""
	if c.InterfaceName != nil {
		iface = *c.InterfaceName
	}
	if iface == "" {
		fe.add("interface_name", "can't be blank")
	}
	if c.MatchPath == "" {
		fe.add("match_path", "can't be blank")
	}
	op := c.ValueMatchOperator
	switch {
	case op == "":
		fe.add("value_match_operator", "can't be blank")
	case !valueOperators[op]:
		fe.add("value_match_operator", "is invalid")
	}
	if iface == anyToken {
		// Upstream's cond reports one forcing violation at a time; a major
		// riding along with "*" is deleted (accepted, ignored).
		switch {
		case c.On != OnIncomingData:
			fe.add("on", "must be incoming_data when interface_name is "+anyToken)
		case c.MatchPath != anyPath:
			fe.add("match_path", "must be /* when interface_name is "+anyToken)
		}
	} else {
		if iface != "" && !ifaceNameRE.MatchString(iface) {
			fe.add("interface_name", "has invalid format")
		}
		if c.InterfaceMajor == nil {
			fe.add("interface_major", "can't be blank")
		}
	}
	deviceFilterErrs(fe, c.DeviceID, c.GroupName)
	if c.MatchPath == anyPath {

$ sed -n '360,480p' internal/engine/triggers/match.go
if c.MatchPath == anyPath {
		if op != anyToken && op != "" {
			fe.add("value_match_operator", "must be * when match_path is "+anyPath)
		}
	} else if c.MatchPath != "" && !matchPathRE.MatchString(c.MatchPath) {
		fe.add("match_path", "has invalid format")
	}
	if op != anyToken && op != "" && len(c.KnownValue) == 0 {
		fe.add("known_value", "can't be blank")
	}
	if len(fe) > 0 {
		return fe
	}

	m := dataMatcher{
		on:        c.On,
		deviceID:  c.DeviceID,
		iface:     iface,
		anyPath:   c.MatchPath == anyPath,
		operator:  c.ValueMatchOperator,
		evaluated: known && c.GroupName == "",
	}
	if c.InterfaceMajor != nil {
		m.major = *c.InterfaceMajor
	}
	if !m.anyPath {
		m.pathSegs = strings.Split(c.MatchPath, "/")[1:]
	}
	if len(c.KnownValue) > 0 {
		if err := json.Unmarshal(c.KnownValue, &m.known); err != nil {
			return fmt.Errorf("known_value does not parse: %w", err)
		}
	}
	t.noteUnsupported(!known, "data_trigger on "+c.On)
	t.noteUnsupported(c.GroupName != "", "group-scoped trigger (group_name)")
	t.data = append(t.data, m)
	return nil
}

// compileDevice applies upstream's device_trigger validation rules (v1.2.0
// simple_trigger_config.ex): interface filters only on introspection
// conditions, minor_updated's name requirements, and the major rule — with
// interface_name "*" the major is accepted and ignored, otherwise (name nil
// or concrete) it is required even for name-less added/removed triggers.
func (t *Trigger) compileDevice(c *simpleTriggerConfig) error {
	fe := fieldErrs{}
	evaluated, known := deviceOns[c.On]
	if c.On == "" {
		fe.add("on", "can't be blank")
	} else if !known {
		fe.add("on", "is invalid")
	}
	deviceFilterErrs(fe, c.DeviceID, c.GroupName)

	name := ""
	hasName := c.InterfaceName != nil && *c.InterfaceName != ""
	if hasName {
		name = *c.InterfaceName
	}
	if hasName && !introspectionOns[c.On] {
		fe.add("interface_name",
			"is allowed only in if 'on' is one of interface_minor_updated, interface_removed, interface_added")
	}
	if introspectionOns[c.On] {
		switch name {
		case "":
			if c.On == OnInterfaceMinorUpdated {
				fe.add("interface_name", "must be set in interface_minor_updated triggers")
			}
			if c.InterfaceMajor == nil {
				fe.add("interface_major", "can't be blank")
			}
		case anyToken:
			if c.On == OnInterfaceMinorUpdated {
				fe.add("interface_name", "must not be '*' in interface_minor_updated triggers")
			}
		default:
			if !ifaceNameRE.MatchString(name) {
				fe.add("interface_name", "has invalid format")
			}
			if c.InterfaceMajor == nil {
				fe.add("interface_major", "can't be blank")
			}
		}
	}
	if c.On == OnIncomingIntrospection && hasName {
		fe.add("interface_name", "must not be set in incoming_introspection triggers")
	}
	if len(fe) > 0 {
		return fe
	}

	m := deviceMatcher{
		on:        c.On,
		deviceID:  c.DeviceID,
		evaluated: evaluated && c.GroupName == "",
	}
	if hasName {
		m.iface = name
	}
	if introspectionOns[c.On] && name != anyToken {
		// Name-less and concrete matchers carry the required major; a "*"
		// filter ignores it entirely.
		m.major = c.InterfaceMajor
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

$ sed -n '480,600p' internal/engine/triggers/match.go
}
	t.Unsupported = append(t.Unsupported, what)
}

// deviceFilterErrs validates the device_id / group_name pair into fe: at most
// one may be set (upstream validate_device_id_xor_group_name) and device_id
// must be "*" or a valid encoded device ID, with upstream's separate message
// for extended (longer than 128-bit) IDs.
func deviceFilterErrs(fe fieldErrs, deviceID, groupName string) {
	if deviceID != "" && groupName != "" {
		fe.add("group_name", "must not be defined if device_id is defined")
	}
	if deviceID != "" && deviceID != anyToken {
		if _, err := deviceid.Parse(deviceID); err != nil {
			if isExtendedDeviceID(deviceID) {
				fe.add("device_id", "is too long, device id must be 128 bits")
			} else {
				fe.add("device_id", "is not a valid device id")
			}
		}
	}
}

// isExtendedDeviceID reports whether s is pure unpadded-base64url text that
// decodes to more than the canonical 128 bits (upstream's extended-id shape:
// every character beyond 22 adds decoded bytes past the 16-byte device ID).
func isExtendedDeviceID(s string) bool {
	if len(s) <= deviceid.EncodedLen {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case 'A' <= c && c <= 'Z', 'a' <= c && c <= 'z', '0' <= c && c <= '9', c == '-', c == '_':
		default:
			return false
		}
	}
	return true
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
		if introspectionOns[m.on] {
			// The name filters only when concrete: a "*" or name-less matcher
			// matches every interface, with the compiled major still applied
			// when present (name-less introspection triggers always carry one).
			if m.iface != "" && m.iface != anyToken && m.iface != ev.Interface {
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

// TracksChanges reports whether any evaluated data condition of this trigger
// reacts to one of the previous-value-derived conditions (value_change*,
// path_created/removed) for this event's device, interface, and path — the
// gate that keeps the engine's pre-write previous-value lookup off messages
// nobody is watching (upstream get_value_change_triggers).
func (t *Trigger) TracksChanges(ev DataEvent) bool {
	for i := range t.data {
		m := &t.data[i]
		if !m.evaluated || !prevOns[m.on] {
			continue
		}
		if !deviceIDMatches(m.deviceID, ev.DeviceID) {
			continue
		}
		if m.iface != anyToken && (m.iface != ev.Interface || m.major != ev.Major) {
			continue
		}
		if !m.anyPath && !pathMatches(m.pathSegs, ev.Path) {
			continue
		}
		return true
	}
	return false
}

// matches evaluates one data condition.
func (m *dataMatcher) matches(ev DataEvent) bool {
	if !m.evaluated || m.on != ev.On {
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
	if ev.On == OnPathRemoved {
		// Upstream evaluates removals without a value: no operator applies.
		return true
	}

$ sed -n '600,770p' internal/engine/triggers/match.go
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

$ sed -n '770,860p' internal/engine/triggers/match.go
(no output)

$ wc -l internal/engine/triggers/match_test.go && sed -n '1,80p' internal/engine/triggers/match_test.go
594 internal/engine/triggers/match_test.go
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

$ grep -n "#70" internal/engine/triggers/match* internal/engine/triggers/actions* internal/engine/* 2>&1 | head -15
internal/engine/triggers/match.go:85:// upstream v1.2.0 formats, probe-frozen (#70): Interface.interface_name_regex
internal/engine/triggers/match.go:316:// definition compiles to exactly the same matcher it did before #70.
grep: internal/engine/forward: Is a directory
grep: internal/engine/stream: Is a directory
grep: internal/engine/testdata: Is a directory
grep: internal/engine/triggers: Is a directory

$ rg "interface_name_regex\|match_path_regex\|#70" docs/DESIGN.md docs/ROADMAP.md 2>&1 | head -5
(no output)

$ grep -rn "#70" .mule 2>&1 | head -5
.mule/.last-output:129:// upstream v1.2.0 formats, probe-frozen (#70): Interface.interface_name_regex
.mule/.last-output:366:// definition compiles to exactly the same matcher it did before #70.
.mule/.last-output:916:$ grep -n "#70" internal/engine/triggers/match* internal/engine/triggers/actions* internal/engine/* 2>&1 | head -15
.mule/.last-output:917:internal/engine/triggers/match.go:85:// upstream v1.2.0 formats, probe-frozen (#70): Interface.interface_name_regex
.mule/.last-output:918:internal/engine/triggers/match.go:316:// definition compiles to exactly the same matcher it did before #70.

$ ls .mule/reviews/ 2>&1 | tail -10
httpx-2026-09-17.md
interfaceschema-2026-09-11.md
observability-2026-09-16.md
pairing-2026-09-09.md
payload-2026-09-12.md
realm-2026-09-09.md
store-2026-09-08.md
store-2026-09-21.md
store-2026-10-04.md
swagger-2026-09-24.md

$ cat .mule/todo.md 2>&1 | head -40
- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
- [x] triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
- [x] compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).

# Mule queue

One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).
The mule runs the topmost open line. Reorder freely; the order is the priority.

If a line needs more detail than fits, write `.mule/tasks/<slug>.md` and keep the line short.

    tools/mule.sh add "<slug>: <outcome>"
    tools/mule.sh loop

A line tagged `[legion]` needs the Legion Go and is skipped automatically while it is asleep;
the queue moves on to the next runnable line rather than stalling behind it.

A line tagged `[readonly]` is a **standing check**, not a piece of work. It verifies something
and is expected to produce no diff at all — so an empty diff is success, not the `wrote
nothing` failure it would be for any other line. It never gets ticked off, because a gate is
never done. Its output lands in `.mule/reports/<slug>.md` with the sha it ran against, and it
is skipped while that sha is still HEAD: re-checking code that has not moved buys nothing and
spends a call on a free provider.

**`mule/queue` is authoritative for this file.** The mule ticks on the Pi and writes its
progress there, so editing the queue on `main` resurrects completed tasks the moment the two
are merged — which has happened once. Add and reorder tasks on `mule/queue`; let them reach
`main` only when that branch is merged.

A line marked `- [~]` is **parked**: real work, but not tick-sized. The mule only ever
picks up `- [ ]`. Benchmark runs live here because a single tier is 5-20 minutes of ingest
alone -- past the per-task budget -- and because they want someone watching. Run one by hand:

    tools/mule.sh legion bench-push
    ssh legion 'cd ~/astrate/bench && ./scripts/run-tier.sh small astrate -base-url ... -housekeeping-key ...'

## Where tasks come from

**This file is not the whole queue, and for real work it is not even the main part of it.**

$ tail -60 .mule/todo.md
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_11266228a001RHKi0D4DeDvwYK

- [!] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 200/401/500) — a `from_token` that does not parse as a device id is refused with `ErrValidation: invalid cursor` (internal/appengine/service.go:150-155) which `writeError` maps to 422 with the flat `{"errors":{"detail":"invalid cursor"}}` body, i.e. the `ErrorDetail` shape and not `FieldErrorsDetail`, so `$ref` the same response the group twin uses for the identical parameter (`ValidationErrors`, yaml:934-935) only after fixing its example, or add a dedicated one. Extend the `from_token` parameter description (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:556) and `InternalServerError` shows `detail: Internal Server Error` (yaml:638), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37, golden fixture pkg/astarteapi/testdata/error_bad_request.json) and `"Internal server error"` (envelope.go:58, testdata/error_internal.json), both measured on upstream 1.2.0 (envelope.go:33-36, 55-57). The other five detail examples on this file (`Not Found`, `Device not found`, `Unauthorized`, `Forbidden`, `Too Many Requests`) are already correct — change only those two. The identical wrong pair is in the other four specs (appengine 1732/1792, housekeeping 362/385/477, realm_management 1265/1348, native 466); leave those for a run of their surface. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
- [x] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pairing_api.yaml, which shows `version: "0.1.0-test"` (yaml:252) — a value no build can produce. The field is filled from `pairing.Config.Version` (internal/pairing/service.go:302), wired from `main.version` (cmd/astrate/main.go:151), whose default is `"0.1.0-dev"` (main.go:53) and which is otherwise set only via `-ldflags "-X main.version=..."`. Use `"0.1.0-dev"`, matching the four compat-version examples docs-sync-native-version-value already corrected to that (docs/api/astrate_native_api.yaml:203, 223, 243, 263). Do NOT document the service's own fallback `DefaultVersion = "0.1.0-astrate"` (service.go:86-88): it is unreachable from `main`, which always passes a non-empty value. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` path parameter (docs/api/astarte_pairing_api.yaml:360) and `RegisterRequest.hw_id` (yaml:373) — but `deviceid.Parse` accepts only the 22-character **unpadded base64url** form, in strict mode, and rejects the standard alphabet's `+` and `/` and any `=` padding with an error wrapping `ErrInvalid` (pkg/deviceid/deviceid.go:26, 37, 39-60; upstream parity cited at 36). Measured: with the same `encoding` the parser uses, `dT6hS2W9TT6LEnP25ks+lg` and `.../lg` fail at byte 19 and the padded 23-character form fails at byte 21, so a client generated from the current spec can send a well-formed standard-base64 ID and get 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}` (internal/pairing/http.go:307-309). State the exact encoding and the 22-character length in both descriptions and add a `pattern`/`minLength`/`maxLength` to `hw_id`. The existing `hw_id: dT6hS2W9TT6LEnP25ks_lg` example (yaml:123) is valid and stays. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in docs/api/astarte_pairing_api.yaml (yaml:374-376), an unconstrained string described only as "Optional Astrate extension specifying the device's initial payload format" — so a generated client can send any string, while `Service.Register` accepts only `"bson"` or `"json"` and rejects everything else with 422 `{"errors":{"initial_payload_format":["is invalid"]}}` (internal/pairing/service.go:179-181 -> internal/pairing/http.go:310-312; the allowed set is stated at service.go:173). Add `enum: [bson, json]` plus a note that absent/empty means unset, and add the three missing `ValidationError` examples while you are in the component: `initial_payload_format: ["is invalid"]` (http.go:310-312), `csr: ["can't be blank"]` (http.go:173-177) and `client_crt: ["can't be blank"]` (http.go:253-257) — the component currently shows only the four hw_id/error_name rows (yaml:594-618) and the `csr`/`client_crt` blank-field cases are what a device actually hits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` description in docs/api/astarte_pairing_api.yaml (yaml:153), "Removes a device from the realm. Returns 204 on success.", which reads as a device-and-data deletion. The row and every byte of its stored data survive: `store.UnregisterDevice` only clears the credential trail (`credentials_secret_hash = ''`, `first_credentials_request = NULL`, `cert_serial = NULL`, `cert_aki = NULL`) and flips `status` back to `'registered'` (internal/store/devices.go:93-110), and the service docstring says so outright ("its data is kept", internal/pairing/service.go:227-229). State the real semantics: the device becomes registrable again, its interfaces, datastream data and group memberships are retained, and a second DELETE is a 204 no-op (the row is still there, so RowsAffected is 1, devices.go:106-108) — not the documented 404. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `base64.RawURLEncoding.Strict()` (pkg/deviceid/deviceid.go:37), so of the 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with `ErrInvalid` while `..._lA` is accepted, and the pairing spec's `hw_id` pattern `^[A-Za-z0-9_-]{22}$` (added by docs-sync-pairing-deviceid-base64url) deliberately does not encode the rule because a pattern is a client-side contract and upstream's leniency is unmeasured. The docstring justifies strictness as byte-identical round-tripping, which is a reason, not a parity claim, and docs/COMPATIBILITY.md records no entry for it. Measure upstream (Elixir `Base.url_decode64!(padding: false)`, which does not check those bits) on the Legion Go, then either record the divergence in docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.
- [x] fdo-rc6-scope-delta-for-giulio [auto]: upstream v1.4.0-rc.6 (2026-09-30, 27 commits, 154 files, almost all FDO) grows the onboarding surface issue #78 is scoped against, and none of it is in #78's body (written against upstream 1.3 / rc.5): new `POST /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}/to0` re-runs TO0 and returns the refreshed rendezvous expiry (`{"data":{"expiry":…}}`, ownership_voucher.ex:63-73; 409 `device_already_onboarded` via fallback_controller.ex:120-125); voucher deletion, which revokes the rendezvous registration first and *keeps* the voucher if that fails (`:rendezvous_revocation_failed` → 500) and refuses another realm's voucher through `ensure_voucher_in_realm` (ownership_voucher.ex:40-60 — now expressible because the voucher row carries its realm); `device_id` accepted on voucher upload; the device registered immediately on upload rather than mid-TO2 (`add_unconfirmed_credentials`, device.ex:63-85); single-entry x5chain certs accepted; invalid vouchers no longer stored; vouchers and to2-sessions moved from the per-realm keyspace to the global `astarte` one with `fdo_guid` added to the device row (housekeeping/realms/queries.ex:784, 1129-1133); and RM device deletion now cascading to the voucher (device_removal/core.ex). Write the delta into `.mule/for-giulio.md` as a scope update for #78 and say which of these a client could actually observe — do not edit or comment on the issue (GitHub state is read-only for the mule) and do not implement anything.
- [x] compat-note-v1.4-rc.6 [auto]: record in `.mule/for-giulio.md` that v1.4.0-rc.6 is wire-inert for Astrate, as a delta on the open v1.3.2/v1.3.4 wording proposal rather than a re-derivation. Measured, the whole non-FDO delta is inert: DUP RPC-server reliability is OTP supervisor/GC tuning; the Dashboard's new `unknown_status` is computed client-side from `last_connection && !last_disconnection && !isConnected` (models/Device/index.ts:213-215) and both timestamps are already on Astrate's device shape (internal/appengine/service.go:118), so a device whose disconnection was never recorded reads as "Unknown status" rather than "Disconnected"; the `on_device_registered` dispatch moved from `Agent.register_device` to `Engine.register_device` (pairing/engine.ex:157-165, 186-193) so it also covers the FDO path, but the Agent path is byte-for-byte the same rule and Astrate already fires it on every successful registration (internal/pairing/service.go:221-223); the `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"); and `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma. So the proposal's "v1.4.0 is still rc.5-only" reference becomes rc.6, and its "FDO authentication (pairing, disabled by default)" line understates the 1.4 line, which ships a full ownership-voucher surface Astrate does not emulate (#78, parked). Do not edit docs/COMPATIBILITY.md and do not restate the standing proposal.

- [x] container-stop-deadline [auto]: make `cliInstance.Stop` (internal/flow/blocks/container/docker.go:167-175) derive its 15s bound whenever the passed context carries **no deadline**, not only when it is nil — both cleanup callers pass `context.Background()` (docker.go:122, the "Best-effort cleanup so we do not leave orphans on mapping failure" path, and block.go:121, the not-ready cleanup in `New`), so a wedged docker daemon leaves an unbounded `docker rm -f` inside `exec.CommandContext` holding flow instantiation, while `Block.Stop` (block.go:310-311) does the same job correctly with an explicit 15s timeout. Measured with the package's own `Run` hook: the `rm` invocation arrives with `ctx.Deadline()` unset on both paths. Add a case asserting `inst.Stop(context.Background())` reaches the injected `Run` with a deadline set (it does not today) and that the hostPort-failure path still issues `rm -f <id>` (docker.go:120-124, never executed by any test today). Container-free.
- [!] container-timeout-bounds [auto]: bound `timeout_ms` and `ready_timeout_ms` in `parseConfig` (internal/flow/blocks/container/block.go:205-219) — both are only checked `> 0` and then multiplied unchecked (`time.Duration(n) * time.Millisecond`), which wraps: measured `ready_timeout_ms: 18446744073709` gives `-551.616µs`, so `New` fails with "container: not ready after wait" without ever issuing a probe (the wrapped negative deadline fails the loop guard at httpbridge.go:64), while `ready_timeout_ms: 9223372036854` gives 2562047h47m16s and `New` never returns (still polling `/healthz` after 3s in the probe). The two consumers of a non-positive duration also disagree — `Bridge.timeout()` (httpbridge.go:39-44) treats it as "use the 5s default" while net/http treats a non-positive `Client.Timeout` as no deadline at all (measured: a client with `Timeout: -551µs` returned err=nil after a 3s handler). Reject a value that would overflow or exceed a documented maximum, naming the key and the value in the error, and add table cases asserting `New` refuses both (it accepts them today). Nothing bounds these on the wire either: internal/flow/blocks/schema.go:19 advertises both as bare `"type":"integer"`. Container-free. — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [x] container-parseconfig-rules-test [auto]: pin the four `parseConfig` rules that no test touches (internal/flow/blocks/container/block.go:188-219 — all `count=0` in `go test -coverprofile`, which is 71.2% for the package): `config` must be a JSON object (192-193), `port` an integer in 1–65535 including the 0 / 65536 / non-integer rejections (197-202), and `timeout_ms` and `ready_timeout_ms` positive (205-210, 215-217) — beside the existing `TestConstructor_RequiresImage` (block_test.go:99), which is the only config case the suite has. The port rule is the one with teeth: it is what keeps `127.0.0.1::99999` out of the `docker run` argv (docker.go:93). Container-free, pure.
- [!] container-response-cap-test [auto]: pin `Bridge.RoundTrip`'s 1 MiB response cap (internal/flow/blocks/container/httpbridge.go:133-140 — `io.LimitReader(resp.Body, maxBody+1)` plus the explicit rejection; lines 136, 138.16 and 138.40 are all uncovered), the only bound between a tenant-supplied container and process memory on the message path: assert a handler answering exactly `maxBody+1` bytes is rejected with the "exceeds" error and that exactly `maxBody` bytes still round-trip, driving a `Bridge{MaxBodyBytes: N}` so the fixture stays small. Also worth one line in the same change: the cap is absent from the operator-facing catalog entry (internal/flow/blocks/info.go:116) and from the schema, so an author whose container emits 2 MiB learns about it only from a runtime error. Container-free. — BLOCKED: lint failed: internal/flow/blocks/container/block_test.go:500:2: redefines-builtin-id: redefinition of the built-in function cap (rev
- [x] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, housekeeping and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
- [!] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy alias fields to `InterfaceDefinition` and `InterfaceMapping` in docs/api/astarte_realm_management_api.yaml — both declare `additionalProperties: false` (yaml:935, 972) yet omit interface-level `quality` (alias of `ownership`) and `aggregate` (alias of `aggregation`) and mapping-level `path` (alias of `endpoint`), all three of which the parser accepts and canonicalises: decoded into dedicated fields (pkg/interfaceschema/parse.go:116-122, 136-138), normalised before any rule runs (parse.go:212-236), and re-encoded canonically by `ParseInterfaceCanonical` so the store never sees them (parse.go:165-170, 307-311). A spec-conformant client — or any generated client with this schema baked in — therefore rejects a valid upstream-shaped install/update with a spurious 422 while the service would have accepted and stored it as `ownership`/`aggregation`/`endpoint`. Add the three as documented deprecated aliases (enums mirroring ownership/aggregation, `path` a plain string) or relax `additionalProperties`, and record the three rules the code enforces: `ownership` + `quality` is a violation (parse.go:218) and `aggregation` + `aggregate` is a violation (parse.go:227), while `endpoint` + `path` silently keeps `endpoint` (parse.go:307-311). None of the three may appear in a `required` list. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstrained string described only as "The device hardware ID." — `Service.DeleteDevice` runs the path value through `deviceid.Parse` (internal/realm/service.go:118-121), which accepts exactly 22 characters of unpadded base64url in strict mode and rejects padding, the standard alphabet's `+` and `/`, and non-canonical trailing bits (pkg/deviceid/deviceid.go:24-26, 33-37, 42-60); every rejection is folded into `store.ErrNotFound` (service.go:119-121) and answered `404 {"errors":{"detail":"Device not found"}}` (internal/realm/http.go:107-109), so a client that sends a well-formed standard-base64 or 23-character padded id is told the device does not exist rather than that its id is malformed. State the encoding and the 22-character length and add `pattern`/`minLength`/`maxLength`, exactly as docs-sync-pairing-deviceid-base64url did for the pairing twin. Do NOT encode the strict trailing-bits rule in the pattern: whether it is a deviation worth recording is still open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_realm_management_api.yaml — `data` is declared as a bare `type: object` with no properties (yaml:250-261) although the handler parses it with the same `interfaceschema.ParseInterfaceCanonical` as `installInterface` (internal/realm/service.go:196 vs 161) and then checks the body's `interface_name`/`version_major` against the URL's `{name}`/`{major}` (service.go:203-208), so the update body — the one carrying every interface rule — is the only interface body in the spec with nothing to validate against. `$ref` `#/components/schemas/InterfaceDefinition` as `installInterface` does (yaml:115-117), and say in the operation description that the body must name the same interface and major as the URL (disagreement is 409, already documented) and that `version_minor` must increase (409). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.
- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` shows `detail: Bad Request` (yaml:1732) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1792), while the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58 via `WriteBadRequest`/`WriteInternalServerError`, golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json); every 400/500 on that surface goes through the same astarteapi writers. The other detail examples on the file are already correct — change only those two — and add the docs-half assertion `TestAppEngineErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go, which pins the pairing and realm-management halves. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDetail` schema-level example (yaml:361-362) and the `BadRequest` response (yaml:385) show `detail: Bad Request`, and `InternalServerError` (yaml:477) shows `detail: Internal Server Error`, but the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58, via WriteBadRequest/WriteInternalServerError, pinned by pkg/astarteapi/testdata/error_bad_request.json and error_internal.json). While you are in `ErrorDetail`: its `example` sits at the schema level and shows a bare `{detail: ...}` where the envelope the wire sends is `{"errors":{"detail":...}}` — move it under `errors` (or drop it, as the realm-management and pairing `ErrorDetail` schemas carry no example) so a generated client is not shown the wrong shape. Add `TestHousekeepingErrorDetailExamples` beside the pairing and realm-management halves in internal/swagger/swagger_test.go. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `"Internal server error"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implements `device_empty_cache_received` end to end — `SimpleEvent` field 23, an empty `DeviceEmptyCacheReceivedEvent` whose encoder emits exactly `{"type":"device_empty_cache_received"}` (astarte_core encoder.ex:54-62), fired by DUP at `control_handler.ex:161` *inside* the `with` that runs `send_control_consumer_properties` → `resend_all_properties` → `set_pending_empty_cache`, so it fires only when all three succeed and never on failure, and the commit deletes the `# TODO: implement on_empty_cache_received` from `astarte_events/lib/astarte_events/triggers/core.ex`. Astrate already accepts the trigger but never evaluates it: `OnDeviceEmptyCacheReceived` sits in `deviceOns` with value `false` (internal/engine/triggers/match.go:123, the `evaluated` flag read at match.go:406) and the comment's stated reason ("upstream defines no SimpleEvent variant for it") is now false. Add the event body constructor to internal/engine/triggers/events.go (the payload is `{"type":"device_empty_cache_received"}` and nothing else — no device id, no metadata), flip the `deviceOns` value to true, and fire it from `handleEmptyCache` (internal/engine/control.go:58) with the same idiom `handleLifecycle` uses for `device_connected` (`e.fireDevice` + `stream.Publish`, internal/engine/engine.go:346-353), at the point where `m.Ack()` runs — i.e. after `resendServerProperties` and `sendConsumerProperties` have both succeeded inside `retryStore`. Add a T1 test that a trigger on it fires exactly once per accepted emptyCache and does not fire when the resend fails. — BLOCKED: wrote nothing
- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362) — the second half docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail: Bad Request`, the reconstructed title-case form that `TestHousekeepingErrorDetailExamples` (internal/swagger/swagger_test.go:687) now pins against everywhere else on this spec, while the wire emits the frozen `astarteapi.DetailBadRequest` = `"Bad request"` (pkg/astarteapi/envelope.go:37, via `WriteBadRequest` on every 400 — internal/housekeeping/http.go:65, 166, 171, 183). Because the example sits on the generic `ErrorDetail` schema, which all eight error responses `$ref` (yaml:377-477), it is the example a generator falls back to for any response whose own example it cannot resolve — so it teaches `"Bad Request"` for a 404 and a 500 too. Change it to `Bad request`, or drop it as the realm-management and pairing `ErrorDetail` schemas do (neither carries one). Then extend `TestHousekeepingErrorDetailExamples`, whose comment currently says the schema-level example is "deliberately not pinned" (swagger_test.go:684-686), to pin it like the responses. One thing not to "fix": the example is nested one level in — it is the `errors` sub-schema's example, not a bare envelope — so do not re-wrap it in another `errors:`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-patch-422-field-error-examples [auto]: add the three missing PATCH 422 bodies to the `ValidationErrors` response component in docs/api/astarte_housekeeping_api.yaml (yaml:427-436, referenced once at yaml:190) — it carries a single example, `error_name: [invalid_update_parameters]` (internal/housekeeping/http.go:176-177), but `patchRealm` emits three more field-error envelopes: `jwt_public_key_pem: ["can't be blank"]` (http.go:188-189), `device_registration_limit: ["is invalid"]` (http.go:193-194) and `datastream_maximum_storage_retention: ["is invalid"]` (http.go:198-199). Add them as further `examples` on the media type (OpenAPI 3.0.3, yaml:1) or split the component per case. Also record that the *messages* differ from POST for the same conditions: create rejects a negative limit/retention with the flat `ErrValidation` detail `device_registration_limit must be non-negative` / `datastream_maximum_storage_retention must be non-negative` (internal/housekeeping/service.go:143-148, 422 `ErrorDetail` shape), while PATCH answers `is invalid` in the FieldErrors shape — and that the service's own three `ErrValidation` branches (service.go:205-213) are unreachable from REST, because `patchRealm` pre-checks the identical conditions and answers first. Say in the `RealmPatch` field descriptions that a negative value is a 422 rather than silently ignored. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [!] docs-sync-hk-wrong-type-field-400 [auto]: widen the `BadRequest` component description in docs/api/astarte_housekeeping_api.yaml (yaml:378, "malformed or missing data envelope") — measured today, a perfectly well-formed `data` envelope carrying a wrong-typed field answers **400**, not the documented 422. `astarteapi.DecodeData` into `map[string]json.RawMessage` (internal/housekeeping/http.go:170) cannot fail on a type mismatch, so the failure surfaces on the second decode into the typed body (http.go:182-184) and is answered `WriteBadRequest`. Probe output: `{"data":{"device_registration_limit":"100"}}` → `json: cannot unmarshal string into Go struct field patchBody.device_registration_limit of type int32`; `{"data":{"jwt_public_key_pem":123}}` → the string twin; `{"data":{"device_registration_limit":1.5}}` → `cannot unmarshal number 1.5`; `{"data":[]}` → `cannot unmarshal array into Go value of type patchBody`. Same on POST, which decodes straight into `realmBody` (http.go:63-67). This is unreachable from a spec-conformant client — `RealmCreate`/`RealmPatch` already type both fields — so document it as a one-line extension of the 400 description on both operations rather than inventing a 422, and do NOT re-shape the `BadRequest` example, whose `detail: Bad request` is the frozen constant `TestHousekeepingErrorDetailExamples` pins. One distinction worth stating in the same description: a wrong-typed field is *not* the `invalid_update_parameters` case — unknown keys are caught earlier by the `patchAllowedFields` sweep (http.go:174-180, 422 FieldErrors) on PATCH only, and silently dropped on POST. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `components.parameters.RealmName` and `RealmCreate.realm_name` both carry `pattern: '^[a-z][a-z0-9]*$'` and `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go) holds both against the realm name CHECK extracted from migrations/000002_metadata.up.sql — but the read side is still bare `type: string`: `Realm.realm_name` (components.schemas.Realm) and `DataEnvelopeRealmNames.items`. Decide whether a response field earns the same pattern (a generated client learns nothing new from it, but a reader diffing create against get sees one constrained and one not) and, if it does, add it to both and extend that test to hold the response-side pattern against the same CHECK. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] store-register-inhibit-preserve: RegisterDevice must not overwrite status to 'registered' when a device was inhibited before first credentials request (preserve 'inhibited' state). File: internal/store/devices.go — BLOCKED: gates failed
- [x] store-aliasvalues-self-exclusion-test: Add test for AliasValuesTaken to assert self-exclusion (same device's own alias values do not count as taken). File: internal/store/devices_test.go
- [!] store-latestindividual-empty-errnotfound: Add store-level test asserting LatestIndividual returns ErrNotFound for empty series and returns newest on multiple samples. File: internal/store/datastreams_test.go — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [x] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.go
- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
- [x] compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (newest stable, 2026-10-05, maintenance-only: mix.exs/mix.lock version bumps plus a single CHANGELOG "astarte_data_updater_plant RPC server reliability" entry, no wire/API surface change; v1.4.0 is still rc.6-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2/v1.3.4 wording proposal there already stands and its version-reference delta was recorded on 2026-10-05 as reading "newest stable v1.3.5", so fold this in as that one-word delta rather than re-deriving the proposal.

$ grep -rn "triggers" .mule/reviews/*.md 2>&1 | head -5
.mule/reviews/channels-2026-09-18.md:20:- `internal/engine/triggers/match.go` (330-504) — CompileCondition /
.mule/reviews/channels-2026-09-18.md:33:`triggers/match.go:374-396` (a group data trigger compiles with no device or
.mule/reviews/channels-2026-09-18.md:47:triggers to group members; Astrate's group watch degrades to a realm watch
.mule/reviews/cmd-astrate-2026-09-25.md:99:  `triggers.forward.kind = "nats"` (config.go:339-343) on a binary that then fails at boot with
.mule/reviews/deviceid-2026-09-18.md:12:  `internal/engine/triggers/match.go:488-519` (`deviceFilterErrs` /

I’ll review the `internal/engine/triggers` package deep and write a review.
$ sed -n '1,150p' internal/engine/triggers/actions.go
package triggers

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cbroglie/mustache"
	"github.com/prometheus/client_golang/prometheus"
)

// Action is a trigger's parsed delivery action (docs/ROADMAP.md §7.2 file
// 6.12): an HTTP webhook (upstream "http_url"+"http_method", or the legacy
// "http_post_url"), or a custom action routed to the Forwarder extension
// point. Upstream AMQP-shaped actions ("amqp_exchange") are rejected at
// parse time (#64); docs/DESIGN.md §1.1 records the original seam.
type Action struct {
	// Method is the upper-cased HTTP method.
	Method string
	// URL is the webhook endpoint.
	URL string
	// StaticHeaders are the action's extra request headers.
	StaticHeaders map[string]string
	// IgnoreSSLErrors disables server-certificate verification for this
	// action's requests (upstream "ignore_ssl_errors").
	IgnoreSSLErrors bool
	// Custom is the raw action object of a non-HTTP action, delivered
	// through the Forwarder extension point; nil for HTTP actions.
	Custom json.RawMessage
	// Template is the Mustache body template (upstream "template"), set
	// only when TemplateType == "mustache".
	Template string
	// TemplateType is the upstream "template_type" field. Only "mustache"
	// is rendered; any other non-empty value is accepted but ignored (the
	// default JSON envelope is sent), same as before this field existed.
	TemplateType string
}

// httpAction is the upstream HTTP action JSON shape.
type httpAction struct {
	HTTPURL           string            `json:"http_url"`
	HTTPMethod        string            `json:"http_method"`
	HTTPPostURL       string            `json:"http_post_url"` // pre-1.1 legacy: implies POST
	HTTPStaticHeaders map[string]string `json:"http_static_headers"`
	IgnoreSSLErrors   bool              `json:"ignore_ssl_errors"`
	Template          string            `json:"template"`
	TemplateType      string            `json:"template_type"`
}

// httpMethods is the upstream-accepted method set (lowercase on the wire).
var httpMethods = map[string]bool{
	"delete": true, "get": true, "head": true, "options": true,
	"patch": true, "post": true, "put": true,
}

// parseAction validates one action object. It returns the parsed action
// plus the list of accepted-but-not-evaluated features (non-mustache
// template_type values: the default JSON envelope is sent instead).
func parseAction(raw json.RawMessage) (*Action, []string, error) {
	if len(raw) == 0 {
		return nil, nil, fmt.Errorf("missing action")
	}
	var h httpAction
	if err := json.Unmarshal(raw, &h); err != nil {
		return nil, nil, fmt.Errorf("action does not parse: %w", err)
	}

	var unsupported []string
	if h.TemplateType != "" && h.TemplateType != "mustache" {
		unsupported = append(unsupported, fmt.Sprintf("template_type %q (default JSON envelope sent)", h.TemplateType))
	}

	switch {
	case h.HTTPPostURL != "":
		// Legacy pre-1.1 action (#63): its presence picks this branch (probe:
		// post_url + http_url answers "must be blank" on http_url). post_url
		// validates like http_url, the modern fields must all be blank, and
		// no header blocklist/size check applies (headers must be blank).
		fe := newFieldErrors()
		if h.HTTPURL != "" {
			fe.add("http_url", "must be blank")
		}
		if h.HTTPMethod != "" {
			fe.add("http_method", "must be blank")
		}
		if len(h.HTTPStaticHeaders) > 0 {
			fe.add("http_static_headers", "must be blank")
		}
		validateURLField(fe, "http_post_url", h.HTTPPostURL)
		validateTemplate(fe, h.Template)
		if len(fe.Fields) > 0 {
			return nil, nil, fe
		}
		return &Action{
			Method: http.MethodPost, URL: h.HTTPPostURL,
			StaticHeaders: h.HTTPStaticHeaders, IgnoreSSLErrors: h.IgnoreSSLErrors,
			Template: h.Template, TemplateType: h.TemplateType,
		}, unsupported, nil
	case h.HTTPURL != "":
		// Modern action (#63): upstream-shaped field-scoped validation.
		fe := newFieldErrors()
		if h.HTTPMethod == "" {
			fe.add("http_method", "can't be blank")
		} else if !httpMethods[h.HTTPMethod] {
			fe.add("http_method", "is invalid")
		}
		// Mirrors upstream's modern-branch validator; unreachable given the
		// legacy branch above owns every non-empty http_post_url.
		if h.HTTPPostURL != "" {
			fe.add("http_post_url", "must be blank")
		}
		validateURLField(fe, "http_url", h.HTTPURL)
		validateStaticHeaders(fe, h.HTTPStaticHeaders)
		validateTemplate(fe, h.Template)
		if len(fe.Fields) > 0 {
			return nil, nil, fe
		}
		return &Action{
			// Methods arrive lowercase on the wire (upstream convention) but
			// HTTP methods are case-sensitive uppercase tokens.
			Method: strings.ToUpper(h.HTTPMethod), URL: h.HTTPURL,
			StaticHeaders: h.HTTPStaticHeaders, IgnoreSSLErrors: h.IgnoreSSLErrors,
			Template: h.Template, TemplateType: h.TemplateType,
		}, unsupported, nil
	default:
		// Upstream AMQP action (#64): rejected at compile time rather than
		// silently turned into a Custom action whose events vanish into a
		// subject nobody may consume (validate-and-reject beats silent
		// misbehavior). A legacy stored amqp trigger therefore fails loudly
		// at engine reload — that is the intended trade. The Forwarder
		// extension point itself stays open for non-amqp custom actions.
		var sniff struct {
			AMQPExchange string `json:"amqp_exchange"`
		}
		if err := json.Unmarshal(raw, &sniff); err == nil && sniff.AMQPExchange != "" {
			return nil, nil, fmt.Errorf("amqp trigger actions are not supported: " +
				"use an http action (upstream AMQP delivery is deliberately not implemented)")
		}
		// Not an HTTP action either (e.g. a NATS-flavoured custom object):
		// keep it verbatim for the Forwarder extension point.
		return &Action{Custom: raw}, unsupported, nil
	}

$ sed -n '150,250p' internal/engine/triggers/actions.go
}
}

// Forwarder is the extension point for non-HTTP trigger actions
// (docs/DESIGN.md §1.1: "AMQP action replaced by optional NATS/HTTP
// forwarding"). The executor hands it every matched event whose action is
// not an HTTP webhook. internal/engine/forward provides the HTTP and NATS
// (build-tag "nats") implementations, wired from [triggers.forward] by
// internal/config and cmd/astrate. The default (nil) logs the event and
// counts it as skipped, which is the designed behaviour, not a gap in this
// code path.
type Forwarder interface {
	// Forward delivers one rendered event for a custom action. It is a
	// single-shot hand-off: Astrate never retries a returned error, and a
	// trigger policy's retry/discard handlers do not apply to this seam
	// (only maximum_capacity and event_ttl do — see forward).
	Forward(ctx context.Context, realm, trigger string, action json.RawMessage, event []byte) error
}

// Delivery is one matched (trigger, event) pair queued for execution.
type Delivery struct {
	// Realm is the tenant (rides in the Astarte-Realm header).
	Realm string
	// Trigger is the matched trigger.
	Trigger *Trigger
	// Event is the rendered envelope.
	Event SimpleEvent

	// enqueuedAt is stamped by Enqueue and used for TTL checks.
	enqueuedAt time.Time
}

// Executor defaults.
const (
	// DefaultWorkers is the default delivery worker count.
	DefaultWorkers = 4
	// DefaultQueueSize is the default delivery queue capacity.
	DefaultQueueSize = 256
	// DefaultMaxAttempts bounds delivery attempts per event.
	DefaultMaxAttempts = 5
	// DefaultBackoffStart is the first retry delay.
	DefaultBackoffStart = 250 * time.Millisecond
	// DefaultBackoffCap bounds the exponential retry delay.
	DefaultBackoffCap = 5 * time.Second
	// DefaultRequestTimeout bounds one webhook request.
	DefaultRequestTimeout = 10 * time.Second
)

// Delivery outcome labels.
const (
	outcomeDelivered = "delivered"
	outcomeFailed    = "failed"
	outcomeDropped   = "dropped"
	outcomeExpired   = "expired"
	outcomeForwarded = "forwarded"
	outcomeSkipped   = "skipped"
)

// ExecutorConfig carries the executor's knobs; zero values select the
// defaults above.
type ExecutorConfig struct {
	// Workers is the number of delivery goroutines.
	Workers int
	// QueueSize is the bounded delivery queue capacity; a full queue drops
	// the event with a metric (triggers must never backpressure ingestion,
	// docs/DESIGN.md §1.4).
	QueueSize int
	// MaxAttempts bounds attempts per delivery (1 = no retries).
	MaxAttempts int
	// BackoffStart and BackoffCap shape the exponential retry delay.
	BackoffStart time.Duration
	BackoffCap   time.Duration
	// RequestTimeout bounds each webhook request.
	RequestTimeout time.Duration
	// Forwarder handles non-HTTP actions; nil logs and skips them.
	Forwarder Forwarder
	// Registerer receives the executor's collectors; nil leaves them
	// unregistered (they still work, which tests rely on).
	Registerer prometheus.Registerer
	// Logger receives delivery logs (default slog.Default()).
	Logger *slog.Logger
}

// Executor runs trigger actions asynchronously: a bounded queue feeding
// worker goroutines, HTTP webhooks with exponential-backoff retry, and
// delivery-outcome metrics (docs/ROADMAP.md §7.2 file 6.12).
type Executor struct {
	cfg      ExecutorConfig
	log      *slog.Logger
	ch       chan Delivery
	wg       sync.WaitGroup
	closing  chan struct{}
	stopOnce sync.Once

	client         *http.Client
	insecureClient *http.Client

	outcomes *prometheus.CounterVec
	retries  prometheus.Counter

	policyMu       sync.Mutex

$ sed -n '250,350p' internal/engine/triggers/actions.go
policyMu       sync.Mutex
	policyInflight map[string]int

	// sendMu serialises senders against the channel close in Close.
	sendMu sync.RWMutex
}

// NewExecutor builds and starts an executor.
func NewExecutor(cfg ExecutorConfig) *Executor {
	if cfg.Workers <= 0 {
		cfg.Workers = DefaultWorkers
	}
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = DefaultQueueSize
	}
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = DefaultMaxAttempts
	}
	if cfg.BackoffStart <= 0 {
		cfg.BackoffStart = DefaultBackoffStart
	}
	if cfg.BackoffCap <= 0 {
		cfg.BackoffCap = DefaultBackoffCap
	}
	if cfg.RequestTimeout <= 0 {
		cfg.RequestTimeout = DefaultRequestTimeout
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	x := &Executor{
		cfg:     cfg,
		log:     cfg.Logger.With("component", "triggers"),
		ch:      make(chan Delivery, cfg.QueueSize),
		closing: make(chan struct{}),
		client:  &http.Client{Timeout: cfg.RequestTimeout},
		insecureClient: &http.Client{
			Timeout: cfg.RequestTimeout,
			Transport: &http.Transport{
				// Explicitly requested per action via ignore_ssl_errors
				// (upstream parity).
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402 -- opt-in upstream-parity action flag
			},
		},
		outcomes: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "astrate_engine_trigger_deliveries_total",
			Help: "Trigger action deliveries by outcome (docs/DESIGN.md §5.2).",
		}, []string{"outcome"}),
		retries: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_engine_trigger_retries_total",
			Help: "Failed webhook attempts that were retried.",
		}),
		policyInflight: make(map[string]int),
	}
	for _, o := range []string{outcomeDelivered, outcomeFailed, outcomeDropped, outcomeExpired, outcomeForwarded, outcomeSkipped} {
		x.outcomes.WithLabelValues(o)
	}
	if cfg.Registerer != nil {
		cfg.Registerer.MustRegister(x.outcomes, x.retries)
	}

	x.wg.Add(cfg.Workers)
	for range cfg.Workers {
		go x.worker()
	}
	return x
}

// Enqueue queues one delivery without blocking; a full queue drops it with
// a metric and a log line.
//
// The read lock spans the closing check and the send: without it a caller
// that passed the check could still be holding an unsent delivery when Close
// closes the channel underneath it, and the send would panic.
func (x *Executor) Enqueue(d Delivery) {
	x.sendMu.RLock()
	defer x.sendMu.RUnlock()

	select {
	case <-x.closing:
		x.outcomes.WithLabelValues(outcomeDropped).Inc()
		return
	default:
	}

	d.enqueuedAt = time.Now()

	if pol := d.Trigger.Policy(); pol != nil {
		x.policyMu.Lock()
		if x.policyInflight[pol.Name] >= pol.MaximumCapacity {
			x.policyMu.Unlock()
			x.outcomes.WithLabelValues(outcomeDropped).Inc()
			x.log.Warn("trigger delivery policy capacity full; event dropped",
				"realm", d.Realm, "trigger", d.Trigger.Name, "policy", pol.Name)
			return
		}
		x.policyInflight[pol.Name]++
		x.policyMu.Unlock()
	}

$ sed -n '350,450p' internal/engine/triggers/actions.go
select {
	case x.ch <- d:
	default:
		if pol := d.Trigger.Policy(); pol != nil {
			x.policyMu.Lock()
			x.policyInflight[pol.Name]--
			x.policyMu.Unlock()
		}
		x.outcomes.WithLabelValues(outcomeDropped).Inc()
		x.log.Warn("trigger delivery queue full; event dropped",
			"realm", d.Realm, "trigger", d.Trigger.Name, "device", d.Event.DeviceID)
	}
}

// Close stops accepting deliveries, lets the workers drain the queue, and
// waits for them bounded by ctx. In-flight retry sleeps abort immediately.
func (x *Executor) Close(ctx context.Context) error {
	x.stopOnce.Do(func() {
		close(x.closing)
		// Take the write lock before closing: it can only be acquired once
		// every Enqueue in flight has finished its send, and every later one
		// sees the closed x.closing and drops instead.
		x.sendMu.Lock()
		close(x.ch)
		x.sendMu.Unlock()
	})
	done := make(chan struct{})
	go func() {
		x.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("triggers: executor close interrupted: %w", ctx.Err())
	}
}

// worker drains the delivery queue.
func (x *Executor) worker() {
	defer x.wg.Done()
	for d := range x.ch {
		x.deliver(d)
	}
}

// deliver executes one delivery end to end.
func (x *Executor) deliver(d Delivery) {
	if pol := d.Trigger.Policy(); pol != nil {
		// Balances the in-flight count Enqueue took for this policy, on
		// every path out of deliver — expiry, marshal failure, delivery.
		defer func() {
			x.policyMu.Lock()
			x.policyInflight[pol.Name]--
			x.policyMu.Unlock()
		}()

		if pol.EventTTL > 0 && time.Since(d.enqueuedAt) > pol.EventTTL {
			x.outcomes.WithLabelValues(outcomeExpired).Inc()
			x.log.Info("trigger delivery expired",
				"realm", d.Realm, "trigger", d.Trigger.Name, "policy", pol.Name,
				"age", time.Since(d.enqueuedAt), "event_ttl", pol.EventTTL)
			return
		}
	}

	body, err := json.Marshal(d.Event)
	if err != nil {
		x.outcomes.WithLabelValues(outcomeFailed).Inc()
		x.log.Error("trigger event does not marshal",
			"realm", d.Realm, "trigger", d.Trigger.Name, "err", err)
		return
	}

	if d.Trigger.Action.TemplateType == "mustache" {
		if rendered, err := renderMustache(d.Trigger.Action.Template, d.Realm, d.Event); err != nil {
			x.log.Warn("mustache action template failed to render; sending default envelope",
				"realm", d.Realm, "trigger", d.Trigger.Name, "err", err)
		} else {
			body = rendered
		}
	}

	if d.Trigger.Action.Custom != nil {
		x.forward(d, body)
		return
	}
	x.webhook(d, body)
}

// renderMustache renders a "template_type": "mustache" action body
// (060-triggers.md:516-543) with the realm/device/trigger/event fields
// upstream documents. event fields (interface, path, value, ...) are
// flattened in alongside the envelope fields rather than nested, matching
// upstream's flat template namespace.
func renderMustache(template, realm string, event SimpleEvent) ([]byte, error) {
	ctx := map[string]any{
		"realm":        realm,
		"device_id":    event.DeviceID,

$ sed -n '450,550p' internal/engine/triggers/actions.go
"device_id":    event.DeviceID,
		"timestamp":    event.Timestamp.UTC().Format(EventTimeLayout),
		"trigger_name": event.TriggerName,
	}
	if raw, err := json.Marshal(event.Event); err == nil {
		var fields map[string]any
		if json.Unmarshal(raw, &fields) == nil {
			for k, v := range fields {
				ctx[k] = v
			}
		}
	}
	// forceRaw: HTML-escaping {{var}} would mangle JSON/text bodies, which is
	// what these templates almost always render.
	rendered, err := mustache.RenderRaw(template, true, ctx)
	if err != nil {
		return nil, err
	}
	return []byte(rendered), nil
}

// forward hands a non-HTTP action to the Forwarder extension point. It is
// single-shot by design: the Forwarder seam returns only an error, never an
// HTTP status, so policy.Decide's retry/discard handlers (which decide on a
// status code, and are what webhook deliveries retry on) cannot apply to
// custom actions. The admission/staleness policy knobs do still apply to the
// seam: Enqueue bounds maximum_capacity in-flight and deliver drops past
// event_ttl, exactly as for webhooks. A forwarder failure counts as one
// failed attempt, and when the trigger's policy carries a retry rule the log
// says the retry was not applied instead of abandoning it silently.
func (x *Executor) forward(d Delivery, body []byte) {
	if x.cfg.Forwarder == nil {
		x.outcomes.WithLabelValues(outcomeSkipped).Inc()
		x.log.Info("no forwarder configured; custom trigger action skipped",
			"realm", d.Realm, "trigger", d.Trigger.Name)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), x.cfg.RequestTimeout)
	defer cancel()
	if err := x.cfg.Forwarder.Forward(ctx, d.Realm, d.Trigger.Name, d.Trigger.Action.Custom, body); err != nil {
		x.outcomes.WithLabelValues(outcomeFailed).Inc()
		if pol := d.Trigger.Policy(); pol != nil && pol.RetryTimes > 0 {
			// RetryTimes is nonzero iff some handler retries (CompilePolicy
			// rejects retry_times without a retry handler): the operator
			// expects a retry this path cannot give.
			x.log.Warn("custom action delivery failed; forward path is single-shot and does not apply policy retry",
				"realm", d.Realm, "trigger", d.Trigger.Name, "policy", pol.Name, "err", err)
			return
		}
		x.log.Warn("custom trigger action failed",
			"realm", d.Realm, "trigger", d.Trigger.Name, "err", err)
		return
	}
	x.outcomes.WithLabelValues(outcomeForwarded).Inc()
}

// webhook POSTs (or whatever the action's method is) the event JSON and
// handles the delivery result.
//
// When the trigger has no attached policy (Policy() == nil), behaviour is the
// hardcoded default: 2xx/3xx are delivered, 4xx are permanently failed, 5xx
// and transport errors retry with exponential backoff up to cfg.MaxAttempts.
//
// When a policy is attached, it governs every retry/discard decision: after
// each attempt the HTTP status (or StatusTransport for transport errors) is
// passed to policy.Decide. StrategyRetry retries with the same exponential
// backoff; StrategyDiscard ends the delivery as failed. The attempt cap is
// policy.RetryTimes+1 (RetryTimes counts retries, not attempts). Success
// remains status in 200..399, checked before consulting the policy.
//
// A request that never produced a response is passed as StatusTransport and
// treated as a server error: upstream's policy schema speaks only in status
// codes, so this is Astrate's reading, recorded in docs/COMPATIBILITY.md.
// Every decision logs the reason that produced it.
func (x *Executor) webhook(d Delivery, body []byte) {
	policy := d.Trigger.Policy()
	backoff := x.cfg.BackoffStart
	maxAttempts := x.cfg.MaxAttempts
	if policy != nil {
		maxAttempts = policy.RetryTimes + 1
	}
	for attempt := 1; ; attempt++ {
		status, err := x.attempt(d, body)
		if err == nil && status < 400 {
			x.outcomes.WithLabelValues(outcomeDelivered).Inc()
			return
		}
		if policy != nil {
			effectiveStatus := status
			if err != nil {
				effectiveStatus = StatusTransport
			}
			dec := policy.Decide(effectiveStatus)
			if dec.Strategy == StrategyDiscard {
				x.outcomes.WithLabelValues(outcomeFailed).Inc()
				x.log.Warn("webhook delivery discarded by policy",
					"realm", d.Realm, "trigger", d.Trigger.Name,
					"status", status, "policy", policy.Name, "reason", dec.Reason)
				return
			}
			if attempt >= maxAttempts {

$ sed -n '550,630p' internal/engine/triggers/actions.go
if attempt >= maxAttempts {
				x.outcomes.WithLabelValues(outcomeFailed).Inc()
				x.log.Warn("webhook delivery failed after final attempt",
					"realm", d.Realm, "trigger", d.Trigger.Name, "attempts", attempt,
					"status", status, "err", err,
					"policy", policy.Name, "reason", dec.Reason)
				return
			}
			x.log.Info("webhook delivery retrying",
				"realm", d.Realm, "trigger", d.Trigger.Name, "attempt", attempt,
				"status", status, "policy", policy.Name, "reason", dec.Reason)
			x.retries.Inc()
			select {
			case <-x.closing:
				x.outcomes.WithLabelValues(outcomeFailed).Inc()
				x.log.Warn("webhook delivery abandoned at shutdown",
					"realm", d.Realm, "trigger", d.Trigger.Name, "attempts", attempt,
					"policy", policy.Name, "reason", dec.Reason)
				return
			case <-time.After(backoff):
			}
			backoff = min(backoff*2, x.cfg.BackoffCap)
			continue
		}
		if err == nil && status < 500 {
			x.outcomes.WithLabelValues(outcomeFailed).Inc()
			x.log.Warn("webhook permanently refused",
				"realm", d.Realm, "trigger", d.Trigger.Name, "status", status)
			return
		}
		if attempt >= maxAttempts {
			x.outcomes.WithLabelValues(outcomeFailed).Inc()
			x.log.Warn("webhook delivery failed after final attempt",
				"realm", d.Realm, "trigger", d.Trigger.Name, "attempts", attempt,
				"status", status, "err", err)
			return
		}
		x.retries.Inc()
		select {
		case <-x.closing:
			x.outcomes.WithLabelValues(outcomeFailed).Inc()
			x.log.Warn("webhook delivery abandoned at shutdown",
				"realm", d.Realm, "trigger", d.Trigger.Name, "attempts", attempt)
			return
		case <-time.After(backoff):
		}
		backoff = min(backoff*2, x.cfg.BackoffCap)
	}
}

// attempt performs one webhook request, returning the response status (0 on
// transport errors).
func (x *Executor) attempt(d Delivery, body []byte) (int, error) {
	a := d.Trigger.Action
	ctx, cancel := context.WithTimeout(context.Background(), x.cfg.RequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, a.Method, a.URL, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Astarte-Realm", d.Realm)
	for k, v := range a.StaticHeaders {
		req.Header.Set(k, v)
	}

	client := x.client
	if a.IgnoreSSLErrors {
		client = x.insecureClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	// Drain so the connection is reusable; the body itself is irrelevant.
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
	_ = resp.Body.Close()
	return resp.StatusCode, nil
}

$ wc -l internal/engine/triggers/actions_test.go && sed -n '1190,1210p' internal/engine/triggers/actions_test.go
1199 internal/engine/triggers/actions_test.go

	// Close while the senders are mid-flight; every Enqueue that loses the
	// race must drop its delivery, not panic.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := x.Close(ctx); err != nil {
		t.Fatalf("close during concurrent enqueue: %v", err)
	}
	wg.Wait()
}

Let me read the forward function fully and understand the issue.
$ sed -n '430,490p' internal/engine/triggers/actions.go
} else {
			body = rendered
		}
	}

	if d.Trigger.Action.Custom != nil {
		x.forward(d, body)
		return
	}
	x.webhook(d, body)
}

// renderMustache renders a "template_type": "mustache" action body
// (060-triggers.md:516-543) with the realm/device/trigger/event fields
// upstream documents. event fields (interface, path, value, ...) are
// flattened in alongside the envelope fields rather than nested, matching
// upstream's flat template namespace.
func renderMustache(template, realm string, event SimpleEvent) ([]byte, error) {
	ctx := map[string]any{
		"realm":        realm,
		"device_id":    event.DeviceID,
		"timestamp":    event.Timestamp.UTC().Format(EventTimeLayout),
		"trigger_name": event.TriggerName,
	}
	if raw, err := json.Marshal(event.Event); err == nil {
		var fields map[string]any
		if json.Unmarshal(raw, &fields) == nil {
			for k, v := range fields {
				ctx[k] = v
			}
		}
	}
	// forceRaw: HTML-escaping {{var}} would mangle JSON/text bodies, which is
	// what these templates almost always render.
	rendered, err := mustache.RenderRaw(template, true, ctx)
	if err != nil {
		return nil, err
	}
	return []byte(rendered), nil
}

// forward hands a non-HTTP action to the Forwarder extension point. It is
// single-shot by design: the Forwarder seam returns only an error, never an
// HTTP status, so policy.Decide's retry/discard handlers (which decide on a
// status code, and are what webhook deliveries retry on) cannot apply to
// custom actions. The admission/staleness policy knobs do still apply to the
// seam: Enqueue bounds maximum_capacity in-flight and deliver drops past
// event_ttl, exactly as for webhooks. A forwarder failure counts as one
// failed attempt, and when the trigger's policy carries a retry rule the log
// says the retry was not applied instead of abandoning it silently.
func (x *Executor) forward(d Delivery, body []byte) {
	if x.cfg.Forwarder == nil {
		x.outcomes.WithLabelValues(outcomeSkipped).Inc()
		x.log.Info("no forwarder configured; custom trigger action skipped",
			"realm", d.Realm, "trigger", d.Trigger.Name)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), x.cfg.RequestTimeout)
	defer cancel()
	if err := x.cfg.Forwarder.Forward(ctx, d.Realm, d.Trigger.Name, d.Trigger.Action.Custom, body); err != nil {
		x.outcomes.WithLabelValues(outcomeFailed).Inc()

$ sed -n '485,500p' internal/engine/triggers/actions.go
return
	}
	ctx, cancel := context.WithTimeout(context.Background(), x.cfg.RequestTimeout)
	defer cancel()
	if err := x.cfg.Forwarder.Forward(ctx, d.Realm, d.Trigger.Name, d.Trigger.Action.Custom, body); err != nil {
		x.outcomes.WithLabelValues(outcomeFailed).Inc()
		if pol := d.Trigger.Policy(); pol != nil && pol.RetryTimes > 0 {
			// RetryTimes is nonzero iff some handler retries (CompilePolicy
			// rejects retry_times without a retry handler): the operator
			// expects a retry this path cannot give.
			x.log.Warn("custom action delivery failed; forward path is single-shot and does not apply policy retry",
				"realm", d.Realm, "trigger", d.Trigger.Name, "policy", pol.Name, "err", err)
			return
		}
		x.log.Warn("custom trigger action failed",
			"realm", d.Realm, "trigger", d.Trigger.Name, "err", err)

$ grep -rn "policy.Decide" internal/engine/triggers/actions.go
473:// HTTP status, so policy.Decide's retry/discard handlers (which decide on a
515:// passed to policy.Decide. StrategyRetry retries with the same exponential
542:			dec := policy.Decide(effectiveStatus)

$ grep -rn "forward.*policy" internal/engine/triggers/actions_test.go | head -3
(no output)

$ grep "forward.*policy\|custom.*policy" .mule/todo.md | head -2
- [x] triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
- [x] compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).

$ cat .mule/reviews/observability-2026-09-16.md | head -80
# Review — internal/observability

Date: 2026-09-16
Reviewed files: health.go (1–83), metrics.go (1–80), compat.go (1–52),
observability_test.go (1–128), compat_version_test.go (1–55), and the wiring in
cmd/astrate/main.go:111, 384-413 (NewMetrics, MountServiceCompat,
MountVersionCompat, RegisterBrokerSessions, RegisterDBPool, NewHealth,
AddReadiness ×2, Mount).

Small well-written package: the health/readiness/metrics surface under
/astrate/v1 and the upstream-parity /{service}/health and /{service}/version
endpoints. Most wire contract is already pinned by tests that mirror upstream's
measured envelopes (`{"data":{"status":"ok"}}`, `{"data":"<version>"}`). No
behavioural bug found; three missing-test items below.

## Finding 1: the readiness budget rule has no test

`health.go:12` (`readinessTimeout = 3 * time.Second`) and its comment —
"bounds the whole readiness probe so a wedged dependency can't hang the
endpoint (and the orchestrator that polls it)" — exist, but nothing asserts
them. Both consumers (`handleReadiness`, health.go:56, and `MountServiceCompat`,
compat.go:22) pass the const into `context.WithTimeout`, and no test exercises
the wedged-dependency path at all: `TestReadiness` (observability_test.go:37-63)
runs only checks that return immediately. A regression that dropped the timeout
altogether (hanging every orchestrator poll) would pass the whole suite. The
budget is a const, so a behavioural test would cost 3s of wall time — the
mechanical fix is to make it an injectable `Health` field defaulting to the
const, then assert a ctx-honoring wedged check returns 503 inside the budget.

**Proposed task:** add a `timeout time.Duration` field on `Health` (defaulting
to `readinessTimeout`) and a test proving a ctx-honoring wedged check cannot
hang `/astrate/v1/readiness` past the budget — the handler returns 503 with the
wedged check reported failing, within a short injected budget.

## Finding 2: two of the four db_pool gauges are unasserted

`TestMetricsExposesGauges` (observability_test.go:25-34) asserts
`astrate_db_pool_acquired_conns` and `astrate_db_pool_max_conns`, plus the
broker-sessions gauge and `go_goroutines` — but never `astrate_db_pool_idle_conns`
or `astrate_db_pool_total_conns` (metrics.go:77-78). A regression deleting (or
mis-wiring the pick function of) either gauge would pass the suite.

**Proposed task:** extend `TestMetricsExposesGauges` to assert all four
`astrate_db_pool_*` gauges, matching the values in the supplied
`DBPoolStats{2, 5, 7, 10}` (idle 5, total 7).

## Finding 3: no Content-Type assertion on the /astrate/v1 surface

`writeJSON` (health.go:79-82) sets `Content-Type: application/json`, and the
compat handlers set it by hand (compat.go:26, 31, 44), but only the version test
pins the header (compat_version_test.go:31-33). `/astrate/v1/health`,
`/astrate/v1/readiness`, and the per-service `/health` (both the 200 and the 503
path) are asserted only on status code and body.

**Proposed task:** assert `Content-Type: application/json` on the three
`/astrate/v1` responses in `TestReadiness`/`TestReadinessAllOK` and on both
paths of `TestMountServiceCompat`'s per-service handler.

## Not proposing

- **Per-service `/health` probes only the DB, not the broker** (`main.go:384`
  passes `st.Health` to all three services, while `/astrate/v1/readiness` probes
  DB + broker at main.go:411-412): a broker-only failure would leave the
  Dashboard's per-service indicator green. This is a deliberate-deviation design
  surface (deviation 18, docs/COMPATIBILITY.md) and the fix is a wiring decision
  in main, not a package bug — noted, not proposed.
- **`MustRegister` panic on a second registration** (`metrics.go:57`, `69`):
  startup-only, currently called exactly once each (main.go:400-401), and
  fail-loud-at-init is the Prometheus idiom. Speculative.
- **`AddReadiness` appends unsynchronized** (health.go:38): all four call sites
  are init-time, single-threaded; no live race.
- **`writeJSON` drops the encoder error** (health.go:82): cannot fail for
  `map[string]string`/`map[string]any` bodies of simple strings.
- **Raw dependency error text in the unauthenticated /readiness body**
  (health.go:64): expected errors carry no secrets; guarding would be invented
  behaviour.
- **4× `Store.Stat()` snapshot per scrape** (metrics.go:68-79, each gauge pick
  calls fn independently): the scrape is one Prometheus poll per 15s against a
  pgx pool snapshot; not a named hot workload.

$ grep -rn "custom.*action\|Forwarder" internal/engine/triggers/actions_test.go | head -10
83:// actions routed to the Forwarder, and the rejection / unsupported paths.
125:	t.Run("non-amqp custom action kept verbatim for the forwarder", func(t *testing.T) {
617:// recordingForwarder captures the last Forward call and returns a fixed error.
618:type recordingForwarder struct {
628:func (f *recordingForwarder) Forward(_ context.Context, realm, trigger string, action json.RawMessage, event []byte) error {
636:// TestForwarderSkippedWhenUnset: a custom action with no Forwarder is skipped
638:func TestForwarderSkippedWhenUnset(t *testing.T) {
647:// TestForwarderForwards: a custom action with a Forwarder is handed the realm,
649:func TestForwarderForwards(t *testing.T) {
650:	fwd := &recordingForwarder{}

$ sed -n '630,680p' internal/engine/triggers/actions_test.go
defer f.mu.Unlock()
	f.called++
	f.realm, f.trigger, f.action, f.event = realm, trigger, action, event
	return f.err
}

// TestForwarderSkippedWhenUnset: a custom action with no Forwarder is skipped
// (the designed v1 default, not a failure).
func TestForwarderSkippedWhenUnset(t *testing.T) {
	x := newTestExecutor(t, ExecutorConfig{Workers: 1})
	x.Enqueue(testDelivery(&Action{Custom: json.RawMessage(`{"amqp_exchange":"e"}`)}))
	eventually(t, 1, func() float64 { return outcome(x, outcomeSkipped) })
	if got := outcome(x, outcomeForwarded); got != 0 {
		t.Errorf("forwarded = %v, want 0", got)
	}
}

// TestForwarderForwards: a custom action with a Forwarder is handed the realm,
// trigger name, raw action, and rendered event, and counted as forwarded.
func TestForwarderForwards(t *testing.T) {
	fwd := &recordingForwarder{}
	x := newTestExecutor(t, ExecutorConfig{Workers: 1, Forwarder: fwd})
	raw := json.RawMessage(`{"amqp_exchange":"events"}`)
	d := testDelivery(&Action{Custom: raw})
	x.Enqueue(d)

	eventually(t, 1, func() float64 { return outcome(x, outcomeForwarded) })
	fwd.mu.Lock()
	defer fwd.mu.Unlock()
	if fwd.realm != "testrealm" || fwd.trigger != "hook" || string(fwd.action) != string(raw) {
		t.Errorf("forward args: realm=%q trigger=%q action=%s", fwd.realm, fwd.trigger, fwd.action)
	}
	wantEvent, _ := json.Marshal(d.Event)
	if string(fwd.event) != string(wantEvent) {
		t.Errorf("forward event = %s, want %s", fwd.event, wantEvent)
	}
}

// TestForwarderError: a Forwarder error counts as failed.
func TestForwarderError(t *testing.T) {
	fwd := &recordingForwarder{err: io.ErrUnexpectedEOF}
	x := newTestExecutor(t, ExecutorConfig{Workers: 1, Forwarder: fwd})
	x.Enqueue(testDelivery(&Action{Custom: json.RawMessage(`{"amqp_exchange":"e"}`)}))
	eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
}

// TestForwarderPolicySingleShot pins the designed boundary between delivery
// policies and the custom-action (forward) path: the Forwarder seam returns
// only an error, never an HTTP status, so policy.Decide's retry/discard
// handlers cannot apply to it the way they do to webhook deliveries. Only the
// admission/staleness knobs (maximum_capacity, event_ttl) reach custom

$ sed -n '680,710p' internal/engine/triggers/actions_test.go
// admission/staleness knobs (maximum_capacity, event_ttl) reach custom
// actions. A retry policy therefore still drives exactly one forward attempt,
// counted as failed; routing forward through the Decide loop would break this
// test by retrying. The failure is not silent: when the policy carries a
// retry rule, the log states that the forward path is single-shot.
func TestForwarderPolicySingleShot(t *testing.T) {
	policy := mustCompilePolicy(t, `{
		"name":"fwd-retry",
		"error_handlers":[{"on":"any_error","strategy":"retry"}],
		"retry_times":3,
		"maximum_capacity":1
	}`)

	fwd := &recordingForwarder{err: io.ErrUnexpectedEOF}
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuf, nil))
	x := NewExecutor(ExecutorConfig{
		Workers:      1,
		BackoffStart: time.Millisecond,
		BackoffCap:   5 * time.Millisecond,
		Registerer:   prometheus.NewRegistry(),
		Forwarder:    fwd,
		Logger:       logger,
	})
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := x.Close(ctx); err != nil {
			t.Errorf("executor close: %v", err)
		}
	})

$ sed -n '710,730p' internal/engine/triggers/actions_test.go
})

	d := testDelivery(&Action{Custom: json.RawMessage(`{"nats_subject":"events"}`)})
	d.Trigger.AttachPolicy(policy)
	x.Enqueue(d)

	eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
	if got := outcome(x, outcomeForwarded); got != 0 {
		t.Errorf("forwarded = %v, want 0", got)
	}
	// A Decide-loop forwarder would retry (any_error=retry, retry_times=3):
	// watch the attempt count across a window so any erroneous retry is seen.
	deadline := time.Now().Add(100 * time.Millisecond)
	for time.Now().Before(deadline) {
		fwd.mu.Lock()
		called := fwd.called
		fwd.mu.Unlock()
		if called != 1 {
			t.Fatalf("forwarder called %d times, want 1 (forward path is single-shot)", called)
		}
		time.Sleep(time.Millisecond)

$ sed -n '730,740p' internal/engine/triggers/actions_test.go
time.Sleep(time.Millisecond)
	}
	if got := promtest.ToFloat64(x.retries); got != 0 {
		t.Errorf("retries = %v, want 0", got)
	}
	if !strings.Contains(logBuf.String(), "does not apply policy retry") {
		t.Errorf("log %q missing the single-shot retry notice", logBuf.String())
	}
}

// TestEnqueueAfterCloseDrops: enqueueing after Close drops with a metric and

$ sed -n '740,760p' internal/engine/triggers/actions_test.go
// TestEnqueueAfterCloseDrops: enqueueing after Close drops with a metric and
// never panics on the closed channel.
func TestEnqueueAfterCloseDrops(t *testing.T) {
	x := newTestExecutor(t, ExecutorConfig{Workers: 1})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := x.Close(ctx); err != nil {
		t.Fatalf("close: %v", err)
	}
	x.Enqueue(testDelivery(&Action{Method: http.MethodPost, URL: "https://unused"}))
	if got := outcome(x, outcomeDropped); got != 1 {
		t.Errorf("dropped = %v, want 1", got)
	}
}

// TestQueueFullDrops: when every worker is busy and the bounded queue is full,
// further enqueues drop with a metric (triggers never backpressure ingestion,
// docs/DESIGN.md §1.4).
func TestQueueFullDrops(t *testing.T) {
	started := make(chan struct{}, 8)
	release := make(chan struct{})

$ sed -n '760,790p' internal/engine/triggers/actions_test.go
release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		started <- struct{}{}
		<-release
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	defer close(release)

	x := newTestExecutor(t, ExecutorConfig{Workers: 1, QueueSize: 1})
	a := &Action{Method: http.MethodPost, URL: srv.URL}

	x.Enqueue(testDelivery(a)) // dequeued by the sole worker, now blocked in the handler
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("worker never started the first delivery")
	}
	x.Enqueue(testDelivery(a)) // fills the 1-slot queue
	x.Enqueue(testDelivery(a)) // queue full -> dropped

	if got := outcome(x, outcomeDropped); got != 1 {
		t.Errorf("dropped = %v, want 1", got)
	}
}

// TestCloseDrainsQueue: Close lets queued deliveries finish before returning.
func TestCloseDrainsQueue(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)

$ grep -rn "TestForwarder" internal/engine/triggers/actions_test.go
636:// TestForwarderSkippedWhenUnset: a custom action with no Forwarder is skipped
638:func TestForwarderSkippedWhenUnset(t *testing.T) {
647:// TestForwarderForwards: a custom action with a Forwarder is handed the realm,
649:func TestForwarderForwards(t *testing.T) {
668:// TestForwarderError: a Forwarder error counts as failed.
669:func TestForwarderError(t *testing.T) {
676:// TestForwarderPolicySingleShot pins the designed boundary between delivery
685:func TestForwarderPolicySingleShot(t *testing.T) {

$ sed -n '60,140p' internal/engine/triggers/actions_test.go
}
}

// outcome reads one outcome counter.
func outcome(x *Executor, label string) float64 {
	return promtest.ToFloat64(x.outcomes.WithLabelValues(label))
}

// eventually polls get until it equals want or the deadline passes.
func eventually(t *testing.T, want float64, get func() float64) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if get() == want {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("value did not reach %v (last %v)", want, get())
}

// TestParseAction validates the action-object parser (docs/ROADMAP.md §7.2
// file 6.12): http_url+method, the legacy http_post_url, custom (non-HTTP)
// actions routed to the Forwarder, and the rejection / unsupported paths.
func TestParseAction(t *testing.T) {
	t.Run("http_url+method lowercases on the wire, uppercases for net/http", func(t *testing.T) {
		a, unsupported, err := parseAction([]byte(`{"http_url":"https://x/h","http_method":"post"}`))
		if err != nil {
			t.Fatal(err)
		}
		if a.Method != http.MethodPost || a.URL != "https://x/h" {
			t.Errorf("got %q %q", a.Method, a.URL)
		}
		if a.Custom != nil || len(unsupported) != 0 {
			t.Errorf("custom=%s unsupported=%v", a.Custom, unsupported)
		}
	})

	t.Run("legacy http_post_url implies POST", func(t *testing.T) {
		a, _, err := parseAction([]byte(`{"http_post_url":"https://y/legacy"}`))
		if err != nil {
			t.Fatal(err)
		}
		if a.Method != http.MethodPost || a.URL != "https://y/legacy" {
			t.Errorf("got %q %q", a.Method, a.URL)
		}
	})

	t.Run("static headers and ignore_ssl_errors carry through", func(t *testing.T) {
		a, _, err := parseAction([]byte(`{"http_url":"https://x","http_method":"put",` +
			`"http_static_headers":{"X-Foo":"bar"},"ignore_ssl_errors":true}`))
		if err != nil {
			t.Fatal(err)
		}
		if a.Method != http.MethodPut || a.StaticHeaders["X-Foo"] != "bar" || !a.IgnoreSSLErrors {
			t.Errorf("got %+v", a)
		}
	})

	t.Run("unknown http_method rejected", func(t *testing.T) {
		if _, _, err := parseAction([]byte(`{"http_url":"https://x","http_method":"frobnicate"}`)); err == nil {
			t.Fatal("want error for unknown method")
		}
	})

	t.Run("non-amqp custom action kept verbatim for the forwarder", func(t *testing.T) {
		raw := `{"nats_subject":"events"}`
		a, _, err := parseAction([]byte(raw))
		if err != nil {
			t.Fatal(err)
		}
		if a.Method != "" || a.URL != "" || string(a.Custom) != raw {
			t.Errorf("got %+v", a)
		}
	})

	// #64: upstream AMQP trigger actions are rejected at compile time instead
	// of silently becoming a Custom action that is skipped or forwarded into
	// an unconsumed subject. A legacy stored amqp trigger now fails loudly at
	// engine reload — the intended trade (clear error beats silent
	// misbehavior).

$ grep -rn "mustache\|template" internal/engine/triggers/actions_test.go | head -10
163:	t.Run("mustache template parsed and not marked unsupported", func(t *testing.T) {
165:			`{"http_url":"https://x","http_method":"post","template":"{{ value }}","template_type":"mustache"}`))
310:		// --- template ---
312:			name: "template over the limit regardless of template_type",
313:			raw: `{"http_url":"https://x/h","http_method":"post","template_type":"jinja2","template":"` +
315:			want: map[string][]string{"template": {"should be at most 1048576 character(s)"}},
318:			name: "twin template at the limit compiles even as jinja2", wantMeth: "POST", wantURL: "https://x/h",
320:			raw: `{"http_url":"https://x/h","http_method":"post","template_type":"jinja2","template":"` +
465:// TestWebhookMustacheTemplateRendered: a mustache action renders the body
481:		TemplateType: "mustache",

$ sed -n '440,500p' internal/engine/triggers/actions_test.go
a, _, err := parseAction([]byte(`{"http_url":"` + srv.URL +
		`","http_method":"post","http_static_headers":{"content-type":"application/x-ndjson","astarte-realm":"spoofed"}}`))
	if err != nil {
		t.Fatalf("parseAction rejected a colliding header: %v", err)
	}
	x := newTestExecutor(t, ExecutorConfig{Workers: 1})
	x.Enqueue(testDelivery(a))

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("webhook never received")
	}

	if gotRealm != wantRealm {
		t.Errorf("Astarte-Realm = %q, want the static value %q", gotRealm, wantRealm)
	}
	if len(gotRealmValues) != 1 {
		t.Errorf("Astarte-Realm sent %d time(s) (%q), want a single Set value", len(gotRealmValues), gotRealmValues)
	}
	if gotCT != "application/x-ndjson" {
		t.Errorf("Content-Type = %q, want the static value", gotCT)
	}
}

// TestWebhookMustacheTemplateRendered: a mustache action renders the body
// from realm/device/trigger/event fields instead of the default envelope.
func TestWebhookMustacheTemplateRendered(t *testing.T) {
	var gotBody []byte
	done := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
		close(done)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1})
	d := testDelivery(&Action{
		Method: http.MethodPost, URL: srv.URL,
		Template:     "{{realm}}/{{device_id}}/{{interface}}{{path}}={{value}}",
		TemplateType: "mustache",
	})
	x.Enqueue(d)

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("webhook never received")
	}
	eventually(t, 1, func() float64 { return outcome(x, outcomeDelivered) })

	const want = "testrealm/f0VMRgIBAQAAAAAAAAAAAA/com.ex.Sensors/v=1"
	if string(gotBody) != want {
		t.Errorf("body = %s, want %s", gotBody, want)
	}
}

// TestWebhookMustacheTemplateMalformedFallsBack: an unclosed tag fails to
// render and the delivery falls back to the default JSON envelope rather
// than crashing the dispatcher.

$ sed -n '500,550p' internal/engine/triggers/actions_test.go
// than crashing the dispatcher.
func TestWebhookMustacheTemplateMalformedFallsBack(t *testing.T) {
	var gotBody []byte
	done := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusNoContent)
		close(done)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1})
	d := testDelivery(&Action{
		Method: http.MethodPost, URL: srv.URL,
		Template:     "{{unclosed",
		TemplateType: "mustache",
	})
	x.Enqueue(d)

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("webhook never received")
	}
	eventually(t, 1, func() float64 { return outcome(x, outcomeDelivered) })

	wantBody, _ := json.Marshal(d.Event)
	if string(gotBody) != string(wantBody) {
		t.Errorf("body = %s, want default envelope %s", gotBody, wantBody)
	}
}

// TestWebhookRetriesThenSucceeds is the docs/ROADMAP.md §7.3 case: a 500
// followed by a 200 retries once and then delivers.
func TestWebhookRetriesThenSucceeds(t *testing.T) {
	var calls atomic.Int32
	reqs := make(chan struct{}, 4)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if calls.Add(1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		reqs <- struct{}{}
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1, MaxAttempts: 5})
	x.Enqueue(testDelivery(&Action{Method: http.MethodPost, URL: srv.URL}))

	for i := 0; i < 2; i++ {

$ sed -n '550,600p' internal/engine/triggers/actions_test.go
for i := 0; i < 2; i++ {
		select {
		case <-reqs:
		case <-time.After(5 * time.Second):
			t.Fatalf("only %d webhook requests arrived", i)
		}
	}
	eventually(t, 1, func() float64 { return outcome(x, outcomeDelivered) })
	if got := promtest.ToFloat64(x.retries); got != 1 {
		t.Errorf("retries = %v, want 1", got)
	}
	if got := outcome(x, outcomeFailed); got != 0 {
		t.Errorf("failed = %v, want 0", got)
	}
}

// TestWebhook4xxIsPermanent: a 4xx response is a permanent refusal — no
// retry, counted as failed.
func TestWebhook4xxIsPermanent(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	x := newTestExecutor(t, ExecutorConfig{Workers: 1, MaxAttempts: 5})
	x.Enqueue(testDelivery(&Action{Method: http.MethodPost, URL: srv.URL}))

	eventually(t, 1, func() float64 { return outcome(x, outcomeFailed) })
	// Give any (erroneous) retry a chance to fire before asserting call count.
	time.Sleep(20 * time.Millisecond)
	if got := calls.Load(); got != 1 {
		t.Errorf("attempts = %d, want 1 (4xx must not retry)", got)
	}
	if got := promtest.ToFloat64(x.retries); got != 0 {
		t.Errorf("retries = %v, want 0", got)
	}
}

// TestWebhookFailsAfterMaxAttempts: persistent 5xx exhausts the attempt
// budget and is counted as failed, with MaxAttempts-1 retries.
func TestWebhookFailsAfterMaxAttempts(t *testing.T) {
	reqs := make(chan struct{}, 8)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		reqs <- struct{}{}
	}))
	defer srv.Close()

	const maxAttempts = 3

$ grep -rn "templates must be at most" internal/engine/triggers/actions.go
(no output)

$ grep -rn "MaxTemplateBytes\|template.*limit" internal/engine/triggers/actions.go | head -5
(no output)

$ grep -rn "1048576\|template.*size" internal/engine/triggers/actions_test.go | head -2
315:			want: map[string][]string{"template": {"should be at most 1048576 character(s)"}},
321:				strings.Repeat("x", 1048576) + `"}`,

$ grep -rn "should be at most" internal/engine/triggers
internal/engine/triggers/actions_test.go:212:			want: map[string][]string{"http_url": {"should be at most 8192 character(s)"}},
internal/engine/triggers/actions_test.go:315:			want: map[string][]string{"template": {"should be at most 1048576 character(s)"}},
internal/engine/triggers/validation.go:124:		e.add(field, fmt.Sprintf("should be at most %d character(s)", maxURLRunes))
internal/engine/triggers/validation.go:161:		e.add("template", fmt.Sprintf("should be at most %d character(s)", maxTemplateRunes))

$ cat internal/engine/triggers/validation.go
package triggers

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"unicode/utf8"
)

// Upstream http_action.ex v1.2.0 limits, probe-frozen (issue #63).
const (
	minURLRunes      = 8
	maxURLRunes      = 8192
	maxHeadersBytes  = 8192 // Σ(byte_size(name)+2+byte_size(value)) must be < this
	maxTemplateRunes = 1048576
)

// blockedHeaderNames is upstream's hop-by-hop/sensitive header blocklist,
// compared against strings.ToLower(strings.TrimSpace(name)).
var blockedHeaderNames = map[string]bool{
	"connection":           true,
	"content-length":       true,
	"date":                 true,
	"host":                 true,
	"te":                   true,
	"upgrade":              true,
	"x-forwarded-for":      true,
	"x-forwarded-host":     true,
	"x-forwarded-proto":    true,
	"sec-websocket-accept": true,
	"proxy-authorization":  true,
	"proxy-authenticate":   true,
}

// FieldErrors carries an upstream-shaped field→messages map for one part of a
// trigger definition. errors.Is/As-friendly so the realm layer can render the
// nested changeset envelope without string parsing.
type FieldErrors struct {
	Part   string // e.g. "action"
	Fields map[string][]string
}

// Error renders a stable human form ("action: field=message, ..." with fields
// sorted) so logs and plain-error consumers stay deterministic.
func (e *FieldErrors) Error() string {
	return e.Part + ": " + formatFieldErrors(e.Fields)
}

// TriggerErrors carries upstream-shaped field errors across the parts of a
// trigger definition; rendered as {"errors":{"action":{...}}} and/or
// {"errors":{"simple_triggers":[...]}} (index-aligned, {} for clean entries).
type TriggerErrors struct {
	Action         map[string][]string
	SimpleTriggers []map[string][]string // nil entries = no errors at that index
}

// Error renders a stable human form: one "action: ..." segment plus one
// "simple_triggers[i]: ..." segment per offending condition (fields sorted
// inside each), joined by "; ".
func (e *TriggerErrors) Error() string {
	var parts []string
	if len(e.Action) > 0 {
		parts = append(parts, "action: "+formatFieldErrors(e.Action))
	}
	for i, fields := range e.SimpleTriggers {
		if len(fields) > 0 {
			parts = append(parts, fmt.Sprintf("simple_triggers[%d]: %s", i, formatFieldErrors(fields)))
		}
	}
	if len(parts) == 0 {
		return "no trigger errors"
	}
	return strings.Join(parts, "; ")
}

// fieldErrs is the per-condition accumulator compileData/compileDevice build;
// it doubles as an error so Compile can collect it into TriggerErrors without
// changing compileSimple's signature.
type fieldErrs map[string][]string

// Error renders the sorted "field=message" form (see formatFieldErrors).
func (f fieldErrs) Error() string { return formatFieldErrors(f) }

// add appends one message for field, preserving accumulation order.
func (f fieldErrs) add(field, msg string) {
	f[field] = append(f[field], msg)
}

// formatFieldErrors renders fields sorted as "field=msg1, msg2" entries
// joined with ", ", so logs and plain-error consumers stay deterministic.
func formatFieldErrors(fields map[string][]string) string {
	names := make([]string, 0, len(fields))
	for f := range fields {
		names = append(names, f)
	}
	sort.Strings(names)
	parts := make([]string, 0, len(names))
	for _, f := range names {
		parts = append(parts, fmt.Sprintf("%s=%s", f, strings.Join(fields[f], ", ")))
	}
	return strings.Join(parts, ", ")
}

// newFieldErrors starts an empty accumulation for the "action" part.
func newFieldErrors() *FieldErrors {
	return &FieldErrors{Part: "action", Fields: map[string][]string{}}
}

// add appends one message for field, preserving accumulation order.
func (e *FieldErrors) add(field, msg string) {
	e.Fields[field] = append(e.Fields[field], msg)
}

// validateURLField applies upstream's URL rules to one URL-typed field
// (http_url, or the legacy http_post_url): length messages first, then the
// format message appended after them on the same field.
func validateURLField(e *FieldErrors, field, raw string) {
	n := utf8.RuneCountInString(raw)
	switch {
	case n < minURLRunes:
		e.add(field, fmt.Sprintf("should be at least %d character(s)", minURLRunes))
	case n > maxURLRunes:
		e.add(field, fmt.Sprintf("should be at most %d character(s)", maxURLRunes))
	}
	if !isHTTPURL(raw) {
		e.add(field, "must be a valid http(s) URL")
	}
}

// isHTTPURL reports whether raw parses as an absolute URL with an http(s)
// scheme and a non-empty host (upstream's scheme+host check).
func isHTTPURL(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

// validateStaticHeaders applies the modern-branch header checks: the
// blocklist first, then the total-size ceiling; both can fail together.
func validateStaticHeaders(e *FieldErrors, headers map[string]string) {
	for name := range headers {
		blocked := blockedHeaderNames[strings.ToLower(strings.TrimSpace(name))]
		if blocked {
			e.add("http_static_headers", "must contain only allowed http headers")
			break
		}
	}
	total := 0
	for name, value := range headers {
		total += len(name) + 2 + len(value)
	}
	if total >= maxHeadersBytes {
		e.add("http_static_headers", fmt.Sprintf("headers total size must be lower than %d", maxHeadersBytes))
	}
}

// validateTemplate enforces the template rune-length cap regardless of
// template_type (upstream never validates template_type itself).
func validateTemplate(e *FieldErrors, template string) {
	if utf8.RuneCountInString(template) > maxTemplateRunes {
		e.add("template", fmt.Sprintf("should be at most %d character(s)", maxTemplateRunes))
	}
}

$ cat .mule/reviews/engine-2026-09-10.md 2>&1 | head -100 | tail -60
cat: .mule/reviews/engine-2026-09-10.md: No such file or directory

$ ls .mule/reviews/ | grep engine
appengine-2026-09-07.md
engine-2026-09-06.md
engine-triggers-2026-09-20.md

$ cat .mule/reviews/engine-triggers-2026-09-20.md | head -140
# Code review — internal/engine (triggers + forward) — 2026-09-20

Second pass over `internal/engine`. The 09-06 review covered the core data
paths; this run targets the files the rotation had not read since: the
`triggers` subpackage's validation/transient/error-name machinery and the
`forward/` extension package, plus the delivery path in `actions.go`.

## What I read

- `internal/engine/triggers/actions.go` (610 lines) — parseAction, Executor,
  Enqueue/deliver/forward/webhook, policy retry loop
- `internal/engine/triggers/validation.go` (163 lines) — URL/header/template
  limits, FieldErrors/TriggerErrors renderers
- `internal/engine/triggers/transient.go` (22 lines) — CompileCondition
- `internal/engine/triggers/errorname.go` (79 lines) — upstream error_name
  closed set + astrateToUpstream mapping
- `internal/engine/triggers/match.go` (definition/compile skeletons, lines
  230-345), `policy.go` (Decide contract)
- `internal/engine/forward/{http,envelope,nats}.go` + http_test.go
- Tests: `errorname_test.go`, `transient_test.go`, `forward/http_test.go`,
  `actions_test.go` (policy suite + forwarder seam tests)
- Callers: `internal/engine/cache.go:322` (AttachPolicy), `cmd/astrate/main.go`
  `newForwarder`, `internal/config/config.go` forward validation
- `pkg/payload/payload.go:148` (ReasonMissingRequired emission)

## What I found (worth proposing)

### 1. `missing_required` is the only `astrateToUpstream` key with no test row

`internal/engine/triggers/errorname.go:64` maps the upstream-1.4 reject reason
`missing_required` → `unexpected_object_key`. `TestUpstreamErrorNameMapping`
(errorname_test.go:12-31) covers all 19 other map keys verbatim; a regression
that drops or remaps this entry (the only 1.4-era one — `errorname.go:62-63`
notes it is outside dashboard 1.2.2's closed set, so it is the most likely to
be "reviewed into" the wrong value later) passes every suite. Add the one
table row, plus an invariant test that every `astrateToUpstream` value is a
member of `UpstreamErrorNames()` — that invariant also catches a
typo-then-copied mapping value, which the fixed-row table alone cannot.

### 2. The custom-action (forward) delivery path never consults an attached policy

`webhook()` implements the full policy contract (actions.go:505-579):
`policy.Decide`, `RetryTimes+1` cap, backoff. `forward()` (actions.go:469-485)
does a single attempt under a `RequestTimeout` context and skips the policy
entirely — no `Decide`, no retry, no discard-by-policy. But `Enqueue`'s
`maximum_capacity` gate (actions.go:335-346) and `deliver`'s `event_ttl`
(actions.go:406-413) DO apply to custom actions: the trigger's policy is
checked with `d.Trigger.Policy()` before the custom/webhook branch at
actions.go:432-436. So a trigger with a custom action + retry policy gets
capacity and TTL semantics but silently not retry/discard semantics. The
executor's own comment (actions.go:494-499) says a policy "governs every
retry/discard decision" with no carve-out. Intent should be confirmed: either
route the forward path through the same Decide loop (retry on
`StatusTransport`/treat as server error, honour discard), or document and pin
that forwarder deliveries are intentionally single-shot. The minimum tick is
a test that pins today's behaviour for a custom action with a policy attached
so the decision is forced.

## What I decided NOT to propose, and why

- **Unbounded `io.Copy(io.Discard, resp.Body)` drain in forward/http.go:117-120**
  vs the bounded `io.LimitReader(resp.Body, 1<<20)` in the webhook path
  (actions.go:607). Divergence between siblings, but the drain is bounded in
  time by the `RequestTimeout` context the executor passes (actions.go:476-477),
  and the 1 MiB cap on the webhook side is about not reading an irrelevant body
  forever — the forwarder reads at most until ctx expiry. Not proposing.
- **NATS forwarder ignores its ctx** (nats.go:44 `_ context.Context`). NATS
  `Publish` is non-blocking best-effort; nothing to cancel. Not a gap.
- **`validMethod` accepts any RFC-7230 token while config.go:355-361 restricts
  to the seven verbs.** config validate runs first and is the gate; the looser
  forward check is the boot-time backstop. Deliberate layering, no drift.
- **`transient.go` / `CompileCondition`** — fully hand-tested (transient_test.go:
  unknown types, missing fields, invalid ops, non-JSON, and the nil-Action pin).
- **`validation.go` helpers** — exercised through parseAction's modern/legacy
  branches; the limits are probe-frozen upstream values with issue refs.

Finding two items in a swept area is the yield; the rest of the triggers
package reads correct and well-covered.

$ tail -15 .mule/reviews/engine-triggers-2026-09-20.md
time by the `RequestTimeout` context the executor passes (actions.go:476-477),
  and the 1 MiB cap on the webhook side is about not reading an irrelevant body
  forever — the forwarder reads at most until ctx expiry. Not proposing.
- **NATS forwarder ignores its ctx** (nats.go:44 `_ context.Context`). NATS
  `Publish` is non-blocking best-effort; nothing to cancel. Not a gap.
- **`validMethod` accepts any RFC-7230 token while config.go:355-361 restricts
  to the seven verbs.** config validate runs first and is the gate; the looser
  forward check is the boot-time backstop. Deliberate layering, no drift.
- **`transient.go` / `CompileCondition`** — fully hand-tested (transient_test.go:
  unknown types, missing fields, invalid ops, non-JSON, and the nil-Action pin).
- **`validation.go` helpers** — exercised through parseAction's modern/legacy
  branches; the limits are probe-frozen upstream values with issue refs.

Finding two items in a swept area is the yield; the rest of the triggers
package reads correct and well-covered.

$ ls .mule/reviews/engine-triggers* 2>&1 | tail -5
.mule/reviews/engine-triggers-2026-09-20.md

The 2026-09-20 review already covered `internal/engine/triggers` (actions, validation, etc.) and found:
1. The `missing_required` → `unexpected_object_key` mapping row was missing from `TestUpstreamErrorNameMapping` (errorname_test.go)
2. The forward/custom-action path doesn't consult policy (single-shot by design)

Looking at the current state - both tasks related to these findings were marked done: `triggers-custom-action-policy-nodecide` [x] and the errorname missing required test [x]. Also the forward policy test exists.

The current code-review recipe asks us to pick one area not done recently. Let’s look for packages/areas not reviewed recently compared to dates we see (recent reviews are Sep 8-24, Oct 4). Maybe `internal/engine/triggers` was reviewed 2026-09-20 - that's recent (today is Oct 6). Let us check other areas.

Looking at the todo list, there are many triggers-related docs sync tasks, but the code review area? The recipe says "pick one package and go deep. Rotate through them across runs — check `.mule/log.md` and `.mule/reviews/` for which areas were done recently". Let us look at .mule/log.md? Let us check.
$ cat .mule/log.md | head -30
# Mule log

One row per task attempt, written by `tools/mule.sh`. This is the record of what the cheap
layer can actually do — read it before deciding whether a kind of task is worth delegating.

`secs` is the honest signal: a task that used most of its 900s budget was too big.

| date | task | outcome | secs | note |
| --- | --- | --- | --- | --- |
| 2026-09-04 | purge-properties-compression-capabilityauto | done | 86s | 6f7a3d6 |
| 2026-09-04 | broker-acl-coldstart-introspection-miss | done | 629s | d201db4 |
| 2026-09-04 | issue-93 | done | 108s | 24ad5b8 |
| 2026-09-04 | broker-disconnect-device-zombie-session | blocked | 117s | wrote nothing |
| 2026-09-04 | broker-offline-acl-tests | blocked | 208s | wrote nothing |
| 2026-09-04 | broker-onconnect-doc-comment | blocked | 438s | tests failed: --- FAIL: TestMQTTSink_Retained (0.02s) |
| 2026-09-05 | empty-introspection-verification | blocked | 106s | wrote nothing |
| 2026-09-05 | probe-trigger-install-notification-delay | blocked | 308s | wrote nothing |
| 2026-09-05 | compat-note-v132 | done | 73s | 4bc3e1a |
| 2026-09-05 | docs-sync-pairing-health-path | done | 186s | d73e225 |
| 2026-09-05 | docs-sync-pairing-register-404 | done | 106s | 44a7cae |
| 2026-09-05 | swagger-httptest-coverage | done | 267s | ba79b34 |
| 2026-09-05 | probe-property-resend-encoding | blocked | 100s | wrote nothing |
| 2026-09-05 | compat-note-v133 | blocked | 44s | wrote nothing |
| 2026-09-05 | flow-validate-source-sink | blocked | 360s | tests failed: --- FAIL: TestMQTTSink_Retained (0.02s) |
| 2026-09-05 | flow-validate-dead-source-sink-recompute | done | 200s | 7445176 |
| 2026-09-06 | docs-sync-rm-datastream-retention-endpoint | done | 232s | 0574fd9 |
| 2026-09-06 | docs-sync-rm-interfaces-detailed-param | done | 352s | 4094145 |
| 2026-09-06 | housekeeping-tests | done | 455s | 7583a7d |
| 2026-09-06 | probe-props-resend-error-triggers | blocked | 331s | wrote nothing |
| 2026-09-06 | server-data-trigger-bus | done | 827s | e46ce11 |

$ ls .mule/reviews/*.md | wc -l && echo "---" && ls .mule/reviews/*.md | tail -5
27
---
.mule/reviews/realm-2026-09-09.md
.mule/reviews/store-2026-09-08.md
.mule/reviews/store-2026-09-21.md
.mule/reviews/store-2026-10-04.md
.mule/reviews/swagger-2026-09-24.md

$ cat .mule/reviews/store-2026-10-04.md | head -40
# Code review: internal/store — 2026-10-04

Area chosen: internal/store. Previous reviews were 2026-09-08 and 2026-09-21.
The package is large (30 source files) and central; I read non-test source files
in descending size order, with particular attention to devices, datastreams,
interfaces, groups/realms/store plumbing, and their corresponding tests.

## What I read

- devices.go (485), devices_test.go (427) — register/unregister/get, inhibit,
  credentials stamp, connect/disconnect/stats, alias values taken, patch merge,
  payload hint, delete
- datastreams.go (434), datastreams_test.go (823) — COPY append, Series,
  LatestIndividual, IndividualSnapshot, ObjectSeries, Downsample, DownsampleLTTB,
  SeriesSpan
- interfaces.go (398), interfaces_test.go (249) — install/update/delete/get/load,
  StoredInterface resolver
- groups.go (265), groups_test.go; realms.go (250), realms_test.go; store.go (204),
  store_test.go; flows.go (207), flows_test.go; pipelines.go (193),
  pipelines_validate_test.go; properties.go (183), properties_test.go;
  crypto.go (136), crypto_test.go; userblocks.go (134), userblocks_test.go;
  notify.go (109), notify_test.go; triggers.go (88); policies.go (88);
  retention.go (54), retention_test.go

## What I found (worth proposing)

### 1. RegisterDevice unconditionally clears inhibit on an unconfirmed device
devices.go:75-91 does `ON CONFLICT (realm_id, id) DO UPDATE SET credentials_secret_hash = EXCLUDED.credentials_secret_hash, status = 'registered' WHERE devices.first_credentials_request IS NULL`. devices.go:251-268 (SetDeviceInhibited) sets status='inhibited' regardless of confirmation state. So an admin can inhibit an unconfirmed device; its next re-registration before first credentials request flips it back to 'registered' silently. The inhibit flag is meant to block credentials/connections (§5.3 DESIGN.md); this reopens an inhibited unconfirmed device.

Untested: devices_test.Lifecycle covers re-register before credentials (rotates secret) but not the inhibit→re-register case. Uninhibit already restores 'registered' if never confirmed (SetDeviceInhibited logic), but the re-registration path overwrites status to 'registered' regardless of prior inhibited state — that’s the bug.

### 2. AliasValuesTaken lacks full branch coverage
AliasValuesTaken (devices.go:343-361) excludes self (`d.id <> $2`) and scans for any other device whose aliases jsonb contains values in the input. It returns true on first match, false on ErrNoRows. The three branches worth pinning are: (a) no other device has any of the values → false; (b) another device has a value → true; (c) the same device has the value (self-match excluded) → must not count as taken (false). devices_test covers basic cases but not (c) explicitly as a store-level assertion isolating the self-exclusion.

### 3. LatestIndividual has no store-level test asserting empty-series yields ErrNotFound
datastreams.go:247-274 implements LatestIndividual (ordered DESC by ts with limit 1). Engine tests use a fake store; the real store path’s ErrNoRows→ErrNotFound conversion and value retrieval (both columns) aren’t directly asserted. Worth pinning: (a) empty series → ErrNotFound; (b) two samples on same series, returns newest; (c) correct value column populated. The ErrNotFound case is particularly valuable (used by trigger previous-value lookups).

### 4. validatePipelineGraph: two error branches untested
pipelines.go (validatePipelineGraph) has branches that reject certain graph states; pipelines_validate_test.go exists but does not assert both error conditions explicitly. Pinning these makes the validation contract robust.

$ sed -n '40,120p' .mule/reviews/store-2026-10-04.md
### 5. Properties purge/ownership semantics edge cases could be tightened
properties_test covers PurgeComplement and server-owned preservation, but a small additional assertion (e.g. preserve server-owned across purge of non-complement set, or mixed update paths) would lock behaviour; not essential. Keep to the strongest four above.

## What I decided NOT to propose, and why

- patchDeviceJSONB no-op when patch empty (devices.go:377-404): harmless, callers don’t send empty patches, no bug.
- Interface notify layering: NotifyInterfacesChanged issued by realm service, not store — deliberate.
- Downsample/LTTB parameter indexing: verified in multi-filter cases; tests already exercise.
- ApplyGlobalRetention/EnforceRealmRetentionCeilings/UpdateFlowRuntime: already tested.
- Cosmetic renames/reformats: linter domain.

$ tail -15 .mule/reviews/store-2026-10-04.md
datastreams.go:247-274 implements LatestIndividual (ordered DESC by ts with limit 1). Engine tests use a fake store; the real store path’s ErrNoRows→ErrNotFound conversion and value retrieval (both columns) aren’t directly asserted. Worth pinning: (a) empty series → ErrNotFound; (b) two samples on same series, returns newest; (c) correct value column populated. The ErrNotFound case is particularly valuable (used by trigger previous-value lookups).

### 4. validatePipelineGraph: two error branches untested
pipelines.go (validatePipelineGraph) has branches that reject certain graph states; pipelines_validate_test.go exists but does not assert both error conditions explicitly. Pinning these makes the validation contract robust.

### 5. Properties purge/ownership semantics edge cases could be tightened
properties_test covers PurgeComplement and server-owned preservation, but a small additional assertion (e.g. preserve server-owned across purge of non-complement set, or mixed update paths) would lock behaviour; not essential. Keep to the strongest four above.

## What I decided NOT to propose, and why

- patchDeviceJSONB no-op when patch empty (devices.go:377-404): harmless, callers don’t send empty patches, no bug.
- Interface notify layering: NotifyInterfacesChanged issued by realm service, not store — deliberate.
- Downsample/LTTB parameter indexing: verified in multi-filter cases; tests already exercise.
- ApplyGlobalRetention/EnforceRealmRetentionCeilings/UpdateFlowRuntime: already tested.
- Cosmetic renames/reformats: linter domain.

$ ls .mule/reviews/deviceid* 2>&1 | head -2
.mule/reviews/deviceid-2026-09-18.md

$ ls .mule/reviews/payload* 2>&1 | head -2
.mule/reviews/payload-2026-09-12.md

$ ls .mule/reviews/interfaceschema* 2>&1 | head -2
.mule/reviews/interfaceschema-2026-09-11.md

$ ls .mule/reviews/broker* 2>&1 | head -2
.mule/reviews/broker-2026-09-04.md
.mule/reviews/broker-2026-09-21.md

$ ls .mule/reviews/flow* 2>&1 | head -5
.mule/reviews/flow-2026-09-05.md
.mule/reviews/flow-2026-09-23.md
.mule/reviews/flowapi-2026-09-10.md

$ ls .mule/reviews/httpx* 2>&1 | head -2
.mule/reviews/httpx-2026-09-17.md

$ ls .mule/reviews/pairing* 2>&1 | head -2
.mule/reviews/pairing-2026-09-09.md

$ ls .mule/reviews/housekeeping* 2>&1 | head -2
.mule/reviews/housekeeping-2026-09-15.md

$ ls .mule/reviews/auth* 2>&1 | head -2
.mule/reviews/auth-2026-09-14.md

$ ls .mule/reviews/config* 2>&1 | head -3
.mule/reviews/config-2026-09-13.md

$ ls .mule/reviews/observability* 2>&1 | head -2
.mule/reviews/observability-2026-09-16.md

$ ls .mule/reviews/realm* 2>&1 | head -2
.mule/reviews/realm-2026-09-09.md

$ ls .mule/reviews/cmd* 2>&1 | head -3
.mule/reviews/cmd-astrate-2026-09-25.md

$ ls .mule/reviews/appengine* 2>&1 | head -2
.mule/reviews/appengine-2026-09-07.md

$ ls .mule/reviews/channels* 2>&1 | head -2
.mule/reviews/channels-2026-09-18.md

$ ls .mule/reviews/testutil* 2>&1 | head -2
ls: cannot access '.mule/reviews/testutil*': No such file or directory

$ ls .mule/reviews/bench* 2>&1 | head -2
ls: cannot access '.mule/reviews/bench*': No such file or directory

$ ls .mule/reviews/pkg* 2>&1 | head -5
ls: cannot access '.mule/reviews/pkg*': No such file or directory

$ ls internal/testutil/ 2>&1 | head -15
astartedevice.go
golden.go
golden_test.go
mqttclient.go
pg.go
pg_test.go

$ ls pkg/ | sort
astarteapi
deviceid
interfaceschema
payload

$ ls internal/ | sort
appengine
auth
broker
config
engine
flow
flowapi
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil

$ ls bench/ 2>&1 | head -5
README.md
astrate
client.go
client_test.go
connstorm.go

$ ls tools/ 2>&1 | head -10
bsoncapture
mule-plan-day.sh
mule-setup-pi.sh
mule-survey.sh
mule.sh
reconcile.sh

$ ls cmd/ 2>&1 | head -5
astrate

