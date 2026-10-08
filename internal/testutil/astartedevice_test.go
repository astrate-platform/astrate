package testutil

import (
	"testing"
	"time"
)

// fakeMessage is a minimal paho.Message for driving the capture handler
// without a live broker.
type fakeMessage struct {
	topic   string
	payload []byte
}

func (m fakeMessage) Duplicate() bool   { return false }
func (m fakeMessage) Qos() byte         { return 0 }
func (m fakeMessage) Retained() bool    { return false }
func (m fakeMessage) Topic() string     { return m.topic }
func (m fakeMessage) MessageID() uint16 { return 0 }
func (m fakeMessage) Payload() []byte   { return m.payload }
func (m fakeMessage) Ack()              {}

// TestWaitForTopicFromAfterMark drives two captures on one topic and asserts
// that a wait marked before the second capture returns the second message,
// not a re-match of the first.
func TestWaitForTopicFromAfterMark(t *testing.T) {
	d := &AstarteDevice{base: "realm/device"}
	topic := d.base + "/iface/path"

	d.capture(nil, fakeMessage{topic: topic, payload: []byte("first")})
	mark := d.Mark()
	d.capture(nil, fakeMessage{topic: topic, payload: []byte("second")})

	got := d.WaitForTopicFrom(t, time.Second, mark, topic)
	if string(got.Payload) != "second" {
		t.Fatalf("wait after mark returned payload %q, want %q", got.Payload, "second")
	}

	// An unmarked wait still sees the first match from the buffer start.
	oldest := d.WaitForTopic(t, time.Second, topic)
	if string(oldest.Payload) != "first" {
		t.Fatalf("unmarked wait returned payload %q, want %q", oldest.Payload, "first")
	}
}
