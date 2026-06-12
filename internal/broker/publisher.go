package broker

import (
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

// Publisher is the inline-client facade for server→device messages
// (docs/ROADMAP.md §6 file 5.7). Publishes bypass the ACL (docs/DESIGN.md
// §3.2) and are delivered — or queued on the persistent session, with
// optional per-message expiry — exactly as device-side subscriptions
// dictate. The engine (M6) and AppEngine (M7) consume it.
type Publisher struct {
	srv *mqtt.Server
	cl  *mqtt.Client
}

// newPublisher builds the facade around a dedicated inline client. The
// client's protocol version is pinned to 5 so injected packets carry
// per-message expiry through mochi's offline-queue sweeper (which honours
// expiry only on v5-tagged packets); on the wire each subscriber still
// receives the message encoded for its own protocol version.
func newPublisher(srv *mqtt.Server) *Publisher {
	cl := srv.NewClient(nil, "local", "astrate-publisher", true)
	cl.Properties.ProtocolVersion = 5
	return &Publisher{srv: srv, cl: cl}
}

// Publish sends a server-side message. qos is clamped to 0..2; retain marks
// the message retained on its topic (server-owned properties, docs/DESIGN.md
// §3.4); a positive expiry bounds how long the message may wait in an
// offline device's queue (datastream expiry, §2.5 — rounded up to whole
// seconds, the MQTT expiry granularity). Zero expiry means the broker
// default (mochi's maximum message expiry, 24h).
func (p *Publisher) Publish(topic string, payload []byte, qos byte, retain bool, expiry time.Duration) error {
	if qos > 2 {
		qos = 2
	}
	pk := packets.Packet{
		FixedHeader: packets.FixedHeader{
			Type:   packets.Publish,
			Qos:    qos,
			Retain: retain,
		},
		TopicName: topic,
		Payload:   payload,
		// Inline publishes skip QoS handshakes, but mochi requires a
		// non-zero packet ID for QoS > 0 validity checks (same convention
		// as mqtt.Server.Publish).
		PacketID: uint16(qos),
	}
	if expiry > 0 {
		secs := int64((expiry + time.Second - 1) / time.Second)
		if secs > int64(^uint32(0)) {
			secs = int64(^uint32(0))
		}
		pk.Properties.MessageExpiryInterval = uint32(secs)
	}
	return p.srv.InjectPacket(p.cl, pk)
}
