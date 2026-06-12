package broker

import (
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/storage"
	"github.com/mochi-mqtt/server/v2/packets"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// newTestSessionStore opens a session store on a temp file, plus a hookless
// mochi server for fabricating clients.
func newTestSessionStore(t *testing.T) (*sessionStore, *mqtt.Server) {
	t.Helper()
	ss := newSessionStore(filepath.Join(t.TempDir(), "sessions.db"), discardLogger())
	if err := ss.Init(nil); err != nil {
		t.Fatalf("Init: %v", err)
	}
	t.Cleanup(func() { _ = ss.Stop() })
	srv := mqtt.New(&mqtt.Options{InlineClient: false, Logger: discardLogger()})
	ss.attach(srv)
	return ss, srv
}

func qosPacket(topic string, packetID uint16, created, expiry int64, pv byte) packets.Packet {
	return packets.Packet{
		FixedHeader:     packets.FixedHeader{Type: packets.Publish, Qos: 1},
		TopicName:       topic,
		Payload:         []byte("payload-" + topic),
		PacketID:        packetID,
		Created:         created,
		Expiry:          expiry,
		ProtocolVersion: pv,
	}
}

func TestSessionStoreInitRejectsConfig(t *testing.T) {
	ss := newSessionStore(filepath.Join(t.TempDir(), "s.db"), discardLogger())
	if err := ss.Init(struct{}{}); err == nil {
		t.Fatal("Init with foreign config: expected error")
	}
}

func TestSessionStoreClientAndSubscriptionRoundTrip(t *testing.T) {
	ss, srv := newTestSessionStore(t)
	cl := srv.NewClient(nil, listenerMTLS, "test/h4-Dx_RYTU-RbpDOTabhRg", false)
	cl.Properties.ProtocolVersion = 4

	ss.OnSessionEstablished(cl, packets.Packet{})
	ss.OnSubscribed(cl, packets.Packet{Filters: packets.Subscriptions{
		{Filter: "test/h4-Dx_RYTU-RbpDOTabhRg/com.ex.Server/#", Qos: 1},
		{Filter: "test/h4-Dx_RYTU-RbpDOTabhRg/denied/#", Qos: 1},
	}}, []byte{1, packets.ErrNotAuthorized.Code})

	clients, err := ss.StoredClients()
	if err != nil {
		t.Fatalf("StoredClients: %v", err)
	}
	if len(clients) != 1 || clients[0].ID != cl.ID || clients[0].ProtocolVersion != 4 {
		t.Fatalf("StoredClients = %+v, want one record for %s", clients, cl.ID)
	}

	subs, err := ss.StoredSubscriptions()
	if err != nil {
		t.Fatalf("StoredSubscriptions: %v", err)
	}
	if len(subs) != 1 || subs[0].Filter != "test/h4-Dx_RYTU-RbpDOTabhRg/com.ex.Server/#" || subs[0].Qos != 1 {
		t.Fatalf("StoredSubscriptions = %+v, want only the granted filter", subs)
	}

	// Unsubscribe removes the record.
	ss.OnUnsubscribed(cl, packets.Packet{Filters: packets.Subscriptions{
		{Filter: "test/h4-Dx_RYTU-RbpDOTabhRg/com.ex.Server/#"},
	}})
	if subs, _ = ss.StoredSubscriptions(); len(subs) != 0 {
		t.Fatalf("after unsubscribe: %+v, want empty", subs)
	}

	// A persistent-session disconnect keeps the client; an expiring one
	// drops it and its dependents.
	ss.OnDisconnect(cl, nil, false)
	if clients, _ = ss.StoredClients(); len(clients) != 1 {
		t.Fatal("non-expiring disconnect must keep the session")
	}
	ss.OnDisconnect(cl, nil, true)
	if clients, _ = ss.StoredClients(); len(clients) != 0 {
		t.Fatal("expiring disconnect must drop the session")
	}
}

func TestSessionStoreTakeoverKeepsSession(t *testing.T) {
	ss, srv := newTestSessionStore(t)
	cl := srv.NewClient(nil, listenerMTLS, "test/h4-Dx_RYTU-RbpDOTabhRg", false)
	ss.OnSessionEstablished(cl, packets.Packet{})

	cl.Stop(packets.ErrSessionTakenOver)
	ss.OnDisconnect(cl, packets.ErrSessionTakenOver, true)
	clients, err := ss.StoredClients()
	if err != nil {
		t.Fatalf("StoredClients: %v", err)
	}
	if len(clients) != 1 {
		t.Fatal("takeover disconnect must not drop the session record")
	}
}

func TestSessionStoreInflightOrderingAndExpiry(t *testing.T) {
	ss, srv := newTestSessionStore(t)
	cl := srv.NewClient(nil, listenerMTLS, "test/h4-Dx_RYTU-RbpDOTabhRg", false)
	ss.OnSessionEstablished(cl, packets.Packet{})

	now := time.Now().Unix()
	// Three messages queued within the same wall-clock second (mochi's
	// replay sort key collides), in this order; plus one already expired.
	ss.OnQosPublish(cl, qosPacket("t/a", 11, now, 0, 5), now, 0)
	ss.OnQosPublish(cl, qosPacket("t/b", 7, now, 0, 5), now, 0)
	ss.OnQosPublish(cl, qosPacket("t/c", 23, now, 0, 5), now, 0)
	ss.OnQosPublish(cl, qosPacket("t/dead", 5, now, now-10, 5), now, 0)

	msgs, err := ss.StoredInflightMessages()
	if err != nil {
		t.Fatalf("StoredInflightMessages: %v", err)
	}
	if len(msgs) != 3 {
		t.Fatalf("got %d messages, want 3 (expired one dropped)", len(msgs))
	}
	wantTopics := []string{"t/a", "t/b", "t/c"}
	var prev int64
	for i, m := range msgs {
		if m.TopicName != wantTopics[i] {
			t.Errorf("position %d: topic %q, want %q (queueing order)", i, m.TopicName, wantTopics[i])
		}
		if i > 0 && m.Created <= prev {
			t.Errorf("position %d: Created %d not strictly increasing (prev %d)", i, m.Created, prev)
		}
		prev = m.Created
	}

	// The expired record is also gone from disk.
	msgs, err = ss.StoredInflightMessages()
	if err != nil {
		t.Fatalf("StoredInflightMessages (second read): %v", err)
	}
	if len(msgs) != 3 {
		t.Fatalf("second read: got %d messages, want 3", len(msgs))
	}

	// A resend re-fires OnQosPublish: the original Seq must survive, so
	// ordering stays stable.
	ss.OnQosPublish(cl, qosPacket("t/a", 11, now, 0, 5), now, 1)
	msgs, _ = ss.StoredInflightMessages()
	if msgs[0].TopicName != "t/a" {
		t.Errorf("after resend: first message %q, want t/a", msgs[0].TopicName)
	}

	// Completion and drop both remove records.
	ss.OnQosComplete(cl, packets.Packet{PacketID: 11})
	ss.OnQosDropped(cl, packets.Packet{PacketID: 7})
	if msgs, _ = ss.StoredInflightMessages(); len(msgs) != 1 || msgs[0].TopicName != "t/c" {
		t.Fatalf("after complete+drop: %+v, want only t/c", msgs)
	}
}

func TestSessionStoreRetainedRoundTrip(t *testing.T) {
	ss, srv := newTestSessionStore(t)
	cl := srv.NewClient(nil, listenerMTLS, "test/h4-Dx_RYTU-RbpDOTabhRg", false)
	now := time.Now().Unix()

	ss.OnRetainMessage(cl, qosPacket("r/keep", 0, now, 0, 5), 1)
	ss.OnRetainMessage(cl, qosPacket("r/dead", 0, now, now-1, 5), 1)
	ss.OnRetainMessage(cl, qosPacket("r/gone", 0, now, 0, 5), 1)
	ss.OnRetainMessage(cl, qosPacket("r/gone", 0, now, 0, 5), -1) // cleared

	msgs, err := ss.StoredRetainedMessages()
	if err != nil {
		t.Fatalf("StoredRetainedMessages: %v", err)
	}
	if len(msgs) != 1 || msgs[0].TopicName != "r/keep" {
		t.Fatalf("StoredRetainedMessages = %+v, want only r/keep", msgs)
	}

	ss.OnRetainedExpired("r/keep")
	if msgs, _ = ss.StoredRetainedMessages(); len(msgs) != 0 {
		t.Fatalf("after OnRetainedExpired: %+v, want empty", msgs)
	}
}

func TestNormalizeInflightOrderPerClient(t *testing.T) {
	mk := func(client string, seq uint64, created int64, topic string) storedMessage {
		return storedMessage{
			Message: storage.Message{Client: client, Created: created, TopicName: topic},
			Seq:     seq,
		}
	}
	in := []storedMessage{
		mk("a", 2, 100, "a/second"),
		mk("a", 1, 100, "a/first"),
		mk("a", 3, 100, "a/third"),
		mk("b", 9, 100, "b/only"),
	}
	out := normalizeInflightOrder(in)
	if len(out) != 4 {
		t.Fatalf("got %d messages, want 4", len(out))
	}
	want := []string{"a/first", "a/second", "a/third", "b/only"}
	for i, m := range out {
		if m.TopicName != want[i] {
			t.Errorf("position %d: %q, want %q", i, m.TopicName, want[i])
		}
	}
	// Client a's Created values must be strictly increasing despite the
	// collision; client b's is untouched.
	wantCreated := []int64{100, 101, 102, 100}
	for i, m := range out {
		if m.Created != wantCreated[i] {
			t.Errorf("position %d: Created = %d, want %d", i, m.Created, wantCreated[i])
		}
	}
}
