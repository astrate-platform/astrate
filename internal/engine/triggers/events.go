package triggers

import (
	"encoding/json"
	"time"
)

// This file renders the trigger event payloads (docs/ROADMAP.md §7.2 file
// 6.11) with upstream parity: the JSON shapes are frozen against Astarte
// v1.2's SimpleEvents JSON encoder (astarte_core
// triggers/simple_events/encoder.ex) and the trigger engine's default
// envelope (astarte_trigger_engine events_consumer.ex):
//
//	{
//	  "timestamp": "2026-06-12T10:00:00.123Z",
//	  "device_id": "<encoded id>",
//	  "event": { "type": "...", ... },
//	  "trigger_name": "<name>"
//	}
//
// The realm rides in the Astarte-Realm HTTP header, not the body, exactly
// as upstream. One documented deviation: longinteger values beyond 2^53
// render as decimal strings (Astrate's uniform JSON re-encoding contract,
// docs/DESIGN.md §2.3), where Elixir's Jason would emit a bignum literal.

// eventTimeLayout renders envelope timestamps: UTC, millisecond precision —
// what upstream's DateTime.from_unix(ms) |> Jason.encode produces.
const eventTimeLayout = "2006-01-02T15:04:05.000Z"

// SimpleEvent is one trigger event envelope, ready for delivery.
type SimpleEvent struct {
	// Timestamp is the event instant.
	Timestamp time.Time
	// DeviceID is the encoded device ID.
	DeviceID string
	// TriggerName is the matched trigger's name (the envelope carries it).
	TriggerName string
	// Event is the typed event body: one of the *Event structs below.
	Event any
}

// MarshalJSON renders the upstream envelope shape.
func (s SimpleEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp   string `json:"timestamp"`
		DeviceID    string `json:"device_id"`
		Event       any    `json:"event"`
		TriggerName string `json:"trigger_name"`
	}{
		Timestamp:   s.Timestamp.UTC().Format(eventTimeLayout),
		DeviceID:    s.DeviceID,
		Event:       s.Event,
		TriggerName: s.TriggerName,
	})
}

// IncomingDataEvent is the incoming_data event body. Value must already be
// JSON-friendly (the engine renders payload values through its canonical
// jsonb encoding before constructing events); nil renders as null (property
// unset, matching upstream's empty-bson-value handling).
type IncomingDataEvent struct {
	// Type is always "incoming_data".
	Type string `json:"type"`
	// Interface is the interface name.
	Interface string `json:"interface"`
	// Path is the concrete data path.
	Path string `json:"path"`
	// Value is the published value (null for property unset).
	Value any `json:"value"`
}

// NewIncomingDataEvent builds an incoming_data event body.
func NewIncomingDataEvent(iface, path string, value any) IncomingDataEvent {
	return IncomingDataEvent{Type: OnIncomingData, Interface: iface, Path: path, Value: value}
}

// DeviceConnectedEvent is the device_connected event body.
type DeviceConnectedEvent struct {
	// Type is always "device_connected".
	Type string `json:"type"`
	// DeviceIPAddress is the peer address of the new connection.
	DeviceIPAddress string `json:"device_ip_address"`
}

// NewDeviceConnectedEvent builds a device_connected event body.
func NewDeviceConnectedEvent(ip string) DeviceConnectedEvent {
	return DeviceConnectedEvent{Type: OnDeviceConnected, DeviceIPAddress: ip}
}

// DeviceDisconnectedEvent is the device_disconnected event body.
type DeviceDisconnectedEvent struct {
	// Type is always "device_disconnected".
	Type string `json:"type"`
}

// NewDeviceDisconnectedEvent builds a device_disconnected event body.
func NewDeviceDisconnectedEvent() DeviceDisconnectedEvent {
	return DeviceDisconnectedEvent{Type: OnDeviceDisconnected}
}

// DeviceErrorEvent is the device_error event body.
type DeviceErrorEvent struct {
	// Type is always "device_error".
	Type string `json:"type"`
	// ErrorName is the rejection reason (Astrate's §2.6 reject-reason
	// labels feed it).
	ErrorName string `json:"error_name"`
	// Metadata carries free-form diagnostic strings.
	Metadata map[string]string `json:"metadata"`
}

// NewDeviceErrorEvent builds a device_error event body. A nil metadata map
// renders as {} (upstream always emits the map).
func NewDeviceErrorEvent(errorName string, metadata map[string]string) DeviceErrorEvent {
	if metadata == nil {
		metadata = map[string]string{}
	}
	return DeviceErrorEvent{Type: OnDeviceError, ErrorName: errorName, Metadata: metadata}
}

// IncomingIntrospectionEvent is the incoming_introspection event body,
// carrying the raw introspection string (upstream's string form).
type IncomingIntrospectionEvent struct {
	// Type is always "incoming_introspection".
	Type string `json:"type"`
	// Introspection is the raw `;`-separated introspection payload.
	Introspection string `json:"introspection"`
}

// NewIncomingIntrospectionEvent builds an incoming_introspection event body.
func NewIncomingIntrospectionEvent(introspection string) IncomingIntrospectionEvent {
	return IncomingIntrospectionEvent{Type: OnIncomingIntrospection, Introspection: introspection}
}

// InterfaceAddedEvent is the interface_added event body.
type InterfaceAddedEvent struct {
	// Type is always "interface_added".
	Type string `json:"type"`
	// Interface is the added interface name.
	Interface string `json:"interface"`
	// MajorVersion and MinorVersion are the declared version.
	MajorVersion int `json:"major_version"`
	MinorVersion int `json:"minor_version"`
}

// NewInterfaceAddedEvent builds an interface_added event body.
func NewInterfaceAddedEvent(iface string, major, minor int) InterfaceAddedEvent {
	return InterfaceAddedEvent{Type: OnInterfaceAdded, Interface: iface, MajorVersion: major, MinorVersion: minor}
}

// InterfaceRemovedEvent is the interface_removed event body.
type InterfaceRemovedEvent struct {
	// Type is always "interface_removed".
	Type string `json:"type"`
	// Interface is the removed interface name.
	Interface string `json:"interface"`
	// MajorVersion is the removed major.
	MajorVersion int `json:"major_version"`
}

// NewInterfaceRemovedEvent builds an interface_removed event body.
func NewInterfaceRemovedEvent(iface string, major int) InterfaceRemovedEvent {
	return InterfaceRemovedEvent{Type: OnInterfaceRemoved, Interface: iface, MajorVersion: major}
}
