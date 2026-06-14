package testutil

import (
	"bytes"
	"compress/zlib"
	"crypto/tls"
	"encoding/binary"
	"io"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/astrate-platform/astrate/pkg/deviceid"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// AstarteDevice is a test client that speaks the Astarte MQTT v1 protocol
// (docs/ROADMAP.md §7.2 file 6.15): a paho MQTT client over the broker's mTLS
// listener whose client ID is the certificate CN "<realm>/<device_id>", with
// helpers for introspection, BSON / JSON-profile data publishes, control-channel
// payloads (emptyCache, producer/properties), and capture of server-owned
// messages pushed back to the device. It is reused by the engine T3 suite and
// the M7/M9 conformance harnesses, so it deliberately does not import
// internal/broker (which imports this package in its tests).
type AstarteDevice struct {
	// Realm is the device's realm name.
	Realm string
	// ID is the device identifier.
	ID deviceid.ID
	// Client is the underlying paho client (for advanced cases — most tests
	// use the helpers below).
	Client paho.Client

	base string

	mu       sync.Mutex
	received []ServerMessage
}

// ServerMessage is one message the broker delivered to the device.
type ServerMessage struct {
	// Topic is the full publish topic.
	Topic string
	// Payload is the message body (empty for a property-unset retained clear).
	Payload []byte
}

// ConnectAstarteDevice connects a test device over the broker's TLS listener.
// brokerURL is an "ssl://host:port" address; tlsCfg must carry the device's
// realm-CA-issued client certificate (testutil.DeviceTLSConfig). The device
// subscribes to "<realm>/<id>/#" — the tolerated superset filter (docs/DESIGN.md
// §3.2) — so every server-owned and control message is captured.
func ConnectAstarteDevice(t testing.TB, brokerURL, realm string, id deviceid.ID, tlsCfg *tls.Config, cleanSession bool) *AstarteDevice {
	t.Helper()
	d := &AstarteDevice{Realm: realm, ID: id, base: realm + "/" + id.String()}
	collect := func(opts *paho.ClientOptions) {
		opts.SetDefaultPublishHandler(func(_ paho.Client, m paho.Message) {
			d.mu.Lock()
			d.received = append(d.received, ServerMessage{Topic: m.Topic(), Payload: append([]byte(nil), m.Payload()...)})
			d.mu.Unlock()
		})
	}
	client, _ := MQTTConnect(t, brokerURL, d.base, cleanSession, tlsCfg, collect)
	d.Client = client
	WaitToken(t, client.Subscribe(d.base+"/#", 2, nil), 5*time.Second)
	return d
}

// Base returns the device's base topic "<realm>/<id>".
func (d *AstarteDevice) Base() string { return d.base }

// DataTopic builds the publish topic for an interface path
// ("<realm>/<id>/<iface><path>"); path includes its leading slash, or is empty
// for an object-aggregated interface whose root is the interface itself.
func (d *AstarteDevice) DataTopic(iface, path string) string {
	return d.base + "/" + iface + path
}

// PublishIntrospection publishes the introspection string on the bare base
// topic (docs/DESIGN.md §3.3) at QoS 2 and waits for the PUBACK/PUBCOMP.
func (d *AstarteDevice) PublishIntrospection(t testing.TB, introspection string) {
	t.Helper()
	WaitToken(t, d.Client.Publish(d.base, 2, false, []byte(introspection)), 5*time.Second)
}

// Introspection renders a "name:major:minor;..." string with deterministic
// (sorted) ordering.
func Introspection(entries map[string][2]int) string {
	names := make([]string, 0, len(entries))
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)
	parts := make([]string, 0, len(names))
	for _, name := range names {
		v := entries[name]
		parts = append(parts, name+":"+strconv.Itoa(v[0])+":"+strconv.Itoa(v[1]))
	}
	return strings.Join(parts, ";")
}

// PublishValue encodes a value in the given wire format and publishes it on the
// interface path, waiting for the acknowledgment.
func (d *AstarteDevice) PublishValue(t testing.TB, iface, path string, v payload.Value, ts *time.Time, format payload.Format, qos byte) {
	t.Helper()
	body, err := payload.Encode(v, ts, format)
	if err != nil {
		t.Fatalf("encoding %s%s: %v", iface, path, err)
	}
	d.PublishRaw(t, iface, path, body, qos)
}

// PublishRaw publishes a pre-encoded body on the interface path. An empty body
// is a property unset (docs/DESIGN.md §3.3).
func (d *AstarteDevice) PublishRaw(t testing.TB, iface, path string, body []byte, qos byte) {
	t.Helper()
	WaitToken(t, d.Client.Publish(d.DataTopic(iface, path), qos, false, body), 5*time.Second)
}

// EmptyCache publishes the control/emptyCache signal (docs/DESIGN.md §3.3),
// telling Astrate to re-send every server-owned property.
func (d *AstarteDevice) EmptyCache(t testing.TB) {
	t.Helper()
	WaitToken(t, d.Client.Publish(d.base+"/control/emptyCache", 2, false, []byte{}), 5*time.Second)
}

// SendProducerProperties publishes the control/producer/properties payload: the
// zlib-framed exhaustive list of "<iface>/<path>" entries the device still holds
// (docs/DESIGN.md §3.3). Astrate purges every device-owned property not listed.
func (d *AstarteDevice) SendProducerProperties(t testing.TB, entries []string) {
	t.Helper()
	WaitToken(t, d.Client.Publish(d.base+"/control/producer/properties", 2, false, DeflateControlList(entries)), 5*time.Second)
}

// Messages returns a snapshot of every server message captured so far.
func (d *AstarteDevice) Messages() []ServerMessage {
	d.mu.Lock()
	defer d.mu.Unlock()
	return append([]ServerMessage(nil), d.received...)
}

// WaitForMessage polls until a captured message satisfies pred, returning it.
// It fails the test on timeout.
func (d *AstarteDevice) WaitForMessage(t testing.TB, timeout time.Duration, what string, pred func(ServerMessage) bool) ServerMessage {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		for _, m := range d.Messages() {
			if pred(m) {
				return m
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out after %s waiting for %s", timeout, what)
	return ServerMessage{}
}

// WaitForTopic waits for a (non-empty) message delivered on an exact topic.
func (d *AstarteDevice) WaitForTopic(t testing.TB, timeout time.Duration, topic string) ServerMessage {
	t.Helper()
	return d.WaitForMessage(t, timeout, "message on "+topic, func(m ServerMessage) bool {
		return m.Topic == topic && len(m.Payload) > 0
	})
}

// Disconnect closes the MQTT connection.
func (d *AstarteDevice) Disconnect() {
	d.Client.Disconnect(100)
}

// DeflateControlList builds an Astarte control payload: a 4-byte big-endian
// uncompressed-size prefix followed by the zlib-deflated ";"-joined list
// (docs/DESIGN.md §3.3–3.4). It mirrors the engine's framing so test devices
// and the engine agree on the wire shape.
func DeflateControlList(entries []string) []byte {
	plain := strings.Join(entries, ";")
	var buf bytes.Buffer
	_, _ = buf.Write(binary.BigEndian.AppendUint32(nil, uint32(len(plain))))
	zw := zlib.NewWriter(&buf)
	_, _ = io.WriteString(zw, plain)
	_ = zw.Close()
	return buf.Bytes()
}

// InflateControlList parses an Astarte control payload (the inverse of
// DeflateControlList), returning the ";"-separated entry list. An empty list
// yields nil. It is used by tests to assert consumer/properties payloads.
func InflateControlList(t testing.TB, frame []byte) []string {
	t.Helper()
	if len(frame) < 4 {
		t.Fatalf("control frame is %d bytes, below the 4-byte size prefix", len(frame))
	}
	declared := binary.BigEndian.Uint32(frame[:4])
	zr, err := zlib.NewReader(bytes.NewReader(frame[4:]))
	if err != nil {
		t.Fatalf("control frame is not zlib: %v", err)
	}
	defer func() { _ = zr.Close() }()
	plain, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("inflating control frame: %v", err)
	}
	if uint32(len(plain)) != declared {
		t.Fatalf("control frame inflated to %d bytes, declared %d", len(plain), declared)
	}
	if len(plain) == 0 {
		return nil
	}
	return strings.Split(string(plain), ";")
}
