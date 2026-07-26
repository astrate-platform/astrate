// Package broker embeds the mochi-mqtt server and gives it Astarte MQTT v1
// semantics (docs/DESIGN.md §3.1–§3.4): mTLS device identity with
// CN = "<realm>/<device_id>", the §3.2 ACL matrix, persistent sessions in a
// bbolt file (so session_present survives Astrate restarts), device lifecycle
// bookkeeping, and an inline publishing facade for the engine and AppEngine.
//
// The broker does not interpret payloads. Every accepted device PUBLISH is
// handed to an Intake (implemented by internal/engine in M6) as an
// InboundMessage whose Ack callback releases the MQTT acknowledgment: for
// QoS >= 1 the PUBACK/PUBREC is withheld until Ack is called, which is how
// persistence-commit ordering (docs/DESIGN.md §5.3) and shard backpressure
// (§1.4) propagate to the device.
package broker

import (
	"net/netip"
	"time"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// InboundMessage is one device PUBLISH delivered to the engine intake
// (docs/ROADMAP.md §6 file 5.1). Payload is an independent copy: it remains
// valid after the broker recycles its packet buffers.
type InboundMessage struct {
	// Realm is the device's realm name (first topic segment).
	Realm string
	// DeviceID is the publishing device, parsed from the connection's
	// certificate CN.
	DeviceID deviceid.ID
	// Topic is the full topic the device published to,
	// "<realm>/<device_id>[/<rest>]".
	Topic string
	// Payload is the raw message body (BSON, JSON, control bytes, or empty
	// for property unset — classification happens in the engine).
	Payload []byte
	// QoS is the publish quality of service as received (0..2).
	QoS byte
	// ReceivedAt is the broker reception timestamp (used as the fallback
	// datastream timestamp when the payload carries none).
	ReceivedAt time.Time
	// Ack releases the broker-held acknowledgment for this message. For
	// QoS >= 1 the device's PUBACK (or PUBREC) is not sent until Ack is
	// called — the engine calls it after the persistence batch commits
	// (docs/DESIGN.md §5.3). Ack is never nil, is safe to call multiple
	// times, and must eventually be called for every QoS >= 1 message the
	// intake accepts, or the publishing device stalls (that stall is the
	// designed backpressure path, §1.4).
	Ack func()
}

// Intake consumes accepted device publishes. internal/engine implements it
// (M6); tests use recorders. Submit may block — a full engine shard blocking
// the broker's per-client read loop is the §1.4 backpressure contract.
type Intake interface {
	Submit(InboundMessage)
}

// LifecycleEventType discriminates LifecycleEvent values. The constants
// double as the Astarte trigger event names fed to the M6 trigger engine.
type LifecycleEventType string

const (
	// EventDeviceConnected fires after a device session is established.
	EventDeviceConnected LifecycleEventType = "device_connected"
	// EventDeviceDisconnected fires after a device connection ends (not on
	// session takeover: the device is still connected, on a new channel).
	EventDeviceDisconnected LifecycleEventType = "device_disconnected"
)

// LifecycleEvent is one device connect/disconnect observation
// (docs/DESIGN.md §3.1 — the work astarte_vmq_plugin does upstream).
type LifecycleEvent struct {
	// Type is the event discriminator.
	Type LifecycleEventType
	// Realm is the device's realm name.
	Realm string
	// DeviceID is the device.
	DeviceID deviceid.ID
	// RemoteIP is the peer address for connects; the zero Addr for
	// disconnects.
	RemoteIP netip.Addr
	// At is the event instant.
	At time.Time
}

// LifecycleSink receives lifecycle events (the device_connected /
// device_disconnected trigger feed). Implementations must not block: they
// are invoked from broker connection goroutines.
type LifecycleSink interface {
	OnLifecycleEvent(LifecycleEvent)
}
