package blocks_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"

	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/flow/blocks"
)

// startMQTTBrokerServer embeds a mochi-mqtt broker on a loopback TCP listener
// and returns a stop function (safe to call more than once, also from tests
// that need to kill the broker early) plus its tcp:// URL.
func startMQTTBrokerServer(t *testing.T) (stop func(), url string) {
	t.Helper()
	srv := mqtt.New(&mqtt.Options{})
	if err := srv.AddHook(new(auth.AllowHook), nil); err != nil {
		t.Fatalf("add allow-all hook: %v", err)
	}
	ln := listeners.NewTCP(listeners.Config{ID: "test-tcp", Address: "127.0.0.1:0"})
	if err := srv.AddListener(ln); err != nil {
		t.Fatalf("add listener: %v", err)
	}
	go func() { _ = srv.Serve() }()
	var closeOnce sync.Once
	stop = func() { closeOnce.Do(func() { _ = srv.Close() }) }
	t.Cleanup(stop)
	return stop, "tcp://" + ln.Address()
}

// startMQTTBroker embeds a broker for the whole test and returns its URL.
func startMQTTBroker(t *testing.T) string {
	t.Helper()
	_, url := startMQTTBrokerServer(t)
	return url
}

// connectPaho connects a helper paho client to url, disconnecting on cleanup.
func connectPaho(t *testing.T, url, clientID string, tweak func(*paho.ClientOptions)) paho.Client {
	t.Helper()
	opts := paho.NewClientOptions().
		AddBroker(url).
		SetClientID(clientID).
		SetCleanSession(true)
	if tweak != nil {
		tweak(opts)
	}
	c := paho.NewClient(opts)
	token := c.Connect()
	if !token.WaitTimeout(5*time.Second) || token.Error() != nil {
		t.Fatalf("paho connect %s: %v", url, token.Error())
	}
	t.Cleanup(func() { c.Disconnect(250) })
	return c
}

// publishPaho publishes one payload and waits for the ack.
func publishPaho(t *testing.T, c paho.Client, topic string, qos byte, retained bool, payload any) {
	t.Helper()
	token := c.Publish(topic, qos, retained, payload)
	if !token.WaitTimeout(5*time.Second) || token.Error() != nil {
		t.Fatalf("publish %s: %v", topic, token.Error())
	}
}

// pollEmit calls src.Emit with short contexts until at least one message
// arrives or the deadline passes.
func pollEmit(t *testing.T, src flow.Source, within time.Duration) []*flow.Message {
	t.Helper()
	deadline := time.Now().Add(within)
	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		msgs, err := src.Emit(ctx)
		cancel()
		if err != nil && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
			t.Fatalf("Emit: %v", err)
		}
		if len(msgs) > 0 {
			return msgs
		}
	}
	t.Fatal("timed out waiting for emitted messages")
	return nil
}

// waitForMessage drains ch until one message arrives or the deadline passes.
func waitForMessage(t *testing.T, ch chan paho.Message, within time.Duration) paho.Message {
	t.Helper()
	select {
	case m := <-ch:
		return m
	case <-time.After(within):
		t.Fatal("timed out waiting for published message")
		return nil
	}
}

func mqttSourceConfig(url string, extra map[string]any) map[string]any {
	cfg := map[string]any{"url": url}
	for k, v := range extra {
		cfg[k] = v
	}
	return cfg
}

const mqttTestTimeout = 10 * time.Second

func TestMQTTSource_HappyPath(t *testing.T) {
	url := startMQTTBroker(t)

	block, err := blocks.MQTTSource("mqttsrc", mqttSourceConfig(url, map[string]any{
		"topics": []string{"t/happy"},
	}), flow.Deps{})
	if err != nil {
		t.Fatalf("MQTTSource: %v", err)
	}
	src, ok := block.(flow.Source)
	if !ok {
		t.Fatal("mqtt_source block is not a flow.Source")
	}
	t.Cleanup(func() {
		if s, ok := block.(flow.Stopper); ok {
			s.Stop()
		}
	})

	publisher := connectPaho(t, url, "publisher-happy", nil)
	publishPaho(t, publisher, "t/happy", 0, false, "hello")

	msg := pollEmit(t, src, mqttTestTimeout)[0]
	if msg.Key != "t/happy" {
		t.Errorf("Key = %q, want t/happy", msg.Key)
	}
	if msg.Type != flow.TypeBinary || msg.Subtype != "" {
		t.Errorf("Type/Subtype = %v/%q, want binary/\"\"", msg.Type, msg.Subtype)
	}
	if data, ok := msg.Data.([]byte); !ok || string(data) != "hello" {
		t.Errorf("Data = %#v, want []byte(\"hello\")", msg.Data)
	}
	if msg.Timestamp <= 0 {
		t.Errorf("Timestamp = %d, want positive", msg.Timestamp)
	}
}

func TestMQTTSource_MultiTopic(t *testing.T) {
	url := startMQTTBroker(t)

	block, err := blocks.MQTTSource("mqttsrc", mqttSourceConfig(url, map[string]any{
		"topics": []string{"t/a", "t/b"},
	}), flow.Deps{})
	if err != nil {
		t.Fatalf("MQTTSource: %v", err)
	}
	src := block.(flow.Source)
	t.Cleanup(func() { block.(flow.Stopper).Stop() })

	publisher := connectPaho(t, url, "publisher-multi", nil)
	publishPaho(t, publisher, "t/b", 0, false, "second")
	publishPaho(t, publisher, "t/a", 0, false, "first")

	got := map[string]string{}
	deadline := time.Now().Add(mqttTestTimeout)
	for len(got) < 2 && time.Now().Before(deadline) {
		for _, msg := range pollEmit(t, src, time.Until(deadline)) {
			data, _ := msg.Data.([]byte)
			got[msg.Key] = string(data)
		}
	}
	if len(got) != 2 {
		t.Fatalf("got %d messages (%v), want 2", len(got), got)
	}
	if got["t/a"] != "first" || got["t/b"] != "second" {
		t.Errorf("messages = %v, want t/a=first t/b=second", got)
	}
}

// The drop branch of the source's delivery channel cannot be unit-tested from
// an external test package: both the channel and the dropped counter are
// unexported internals, so this coverage is document-skipped per phase plan.

func TestMQTTSource_ConnectionLost(t *testing.T) {
	stopBroker, url := startMQTTBrokerServer(t)

	block, err := blocks.MQTTSource("mqttsrc", mqttSourceConfig(url, map[string]any{
		"topics": []string{"t/lost"},
	}), flow.Deps{})
	if err != nil {
		t.Fatalf("MQTTSource: %v", err)
	}
	src := block.(flow.Source)
	t.Cleanup(func() { block.(flow.Stopper).Stop() })

	stopBroker()

	deadline := time.Now().Add(mqttTestTimeout)
	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		_, err := src.Emit(ctx)
		cancel()
		if err != nil && strings.Contains(err.Error(), "mqtt_source: connection lost") {
			return
		}
		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Emit: %v", err)
		}
	}
	t.Fatal("Emit never reported the lost connection")
}

func TestMQTTSink_HappyPath(t *testing.T) {
	url := startMQTTBroker(t)

	received := make(chan paho.Message, 8)
	subscriber := connectPaho(t, url, "sink-sub", nil)
	token := subscriber.Subscribe("t/sink", 0, func(_ paho.Client, m paho.Message) {
		received <- m
	})
	if token.Error() != nil {
		t.Fatalf("subscribe: %v", token.Error())
	}

	sink, err := blocks.MQTTSink("mqttsink", mqttSourceConfig(url, map[string]any{
		"topic": "t/sink",
	}), flow.Deps{})
	if err != nil {
		t.Fatalf("MQTTSink: %v", err)
	}
	t.Cleanup(func() { sink.(flow.Stopper).Stop() })

	if _, err := sink.Process(&flow.Message{
		Key:  "k",
		Type: flow.TypeBinary,
		Data: []byte{0x00, 0xff},
	}); err != nil {
		t.Fatalf("Process binary: %v", err)
	}
	m := waitForMessage(t, received, mqttTestTimeout)
	if m.Topic() != "t/sink" {
		t.Errorf("Topic = %q, want t/sink", m.Topic())
	}
	if string(m.Payload()) != "\x00\xff" {
		t.Errorf("Payload = %#v, want raw binary bytes", m.Payload())
	}

	if _, err := sink.Process(&flow.Message{
		Key:  "k",
		Type: flow.TypeMap,
		Data: map[string]any{"k": "v"},
	}); err != nil {
		t.Fatalf("Process map: %v", err)
	}
	m = waitForMessage(t, received, mqttTestTimeout)
	if got := string(m.Payload()); got != `{"k":"v"}` {
		t.Errorf("Payload = %q, want JSON body", got)
	}
}

func TestMQTTSink_Retained(t *testing.T) {
	url := startMQTTBroker(t)

	sink, err := blocks.MQTTSink("mqttsink", mqttSourceConfig(url, map[string]any{
		"topic":    "t/retained",
		"retained": true,
	}), flow.Deps{})
	if err != nil {
		t.Fatalf("MQTTSink: %v", err)
	}
	t.Cleanup(func() { sink.(flow.Stopper).Stop() })

	if _, err := sink.Process(&flow.Message{
		Key:  "k",
		Type: flow.TypeString,
		Data: "keep-me",
	}); err != nil {
		t.Fatalf("Process: %v", err)
	}

	received := make(chan paho.Message, 8)
	late := connectPaho(t, url, "late-subscriber", nil)
	token := late.Subscribe("t/retained", 0, func(_ paho.Client, m paho.Message) {
		received <- m
	})
	if token.Error() != nil {
		t.Fatalf("late subscribe: %v", token.Error())
	}
	m := waitForMessage(t, received, mqttTestTimeout)
	if !m.Retained() || string(m.Payload()) != "keep-me" {
		t.Errorf("retained message = %q (retained=%v), want keep-me", m.Payload(), m.Retained())
	}
}

func TestMQTTSource_ConfigValidation(t *testing.T) {
	url := startMQTTBroker(t) // acceptance twins need eager connect to succeed

	rows := []struct {
		name    string
		config  map[string]any
		wantErr string
	}{
		{"missing url", map[string]any{"topics": []string{"t/x"}}, "mqtt_source: url is required"},
		{"empty url", map[string]any{"url": "", "topics": []string{"t/x"}}, "mqtt_source: url is required"},
		{"missing topics", map[string]any{"url": url}, "mqtt_source: at least one topic is required"},
		{"empty topics", map[string]any{"url": url, "topics": []string{}}, "mqtt_source: at least one topic is required"},
		{"empty topic entry", map[string]any{"url": url, "topics": []string{"a", ""}}, "mqtt_source: topics must contain only non-empty strings"},
		{"non-string topic", map[string]any{"url": url, "topics": []any{"a", 3}}, "mqtt_source: topics must contain only non-empty strings"},
		{"qos too big", map[string]any{"url": url, "topics": []string{"t/x"}, "qos": 3}, "mqtt_source: qos must be 0, 1 or 2"},
		{"qos negative", map[string]any{"url": url, "topics": []string{"t/x"}, "qos": -1}, "mqtt_source: qos must be 0, 1 or 2"},
		{"qos fractional", map[string]any{"url": url, "topics": []string{"t/x"}, "qos": 1.5}, "mqtt_source: qos must be 0, 1 or 2"},
		{"acceptance twin default qos", map[string]any{"url": url, "topics": []string{"t/x"}}, ""},
		{"acceptance twin json qos", map[string]any{"url": url, "topics": []any{"t/x"}, "qos": float64(2)}, ""},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			block, err := blocks.MQTTSource("mqttsrc", row.config, flow.Deps{})
			if row.wantErr != "" {
				if err == nil || err.Error() != row.wantErr {
					t.Fatalf("err = %v, want %q", err, row.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("MQTTSource: %v", err)
			}
			block.(flow.Stopper).Stop()
			block.(flow.Stopper).Stop() // Stop is idempotent
		})
	}
}

func TestMQTTSink_ConfigValidation(t *testing.T) {
	url := startMQTTBroker(t)

	rows := []struct {
		name    string
		config  map[string]any
		wantErr string
	}{
		{"missing url", map[string]any{"topic": "t/x"}, "mqtt_sink: url is required"},
		{"empty url", map[string]any{"url": "", "topic": "t/x"}, "mqtt_sink: url is required"},
		{"missing topic", map[string]any{"url": url}, "mqtt_sink: topic is required"},
		{"empty topic", map[string]any{"url": url, "topic": ""}, "mqtt_sink: topic is required"},
		{"qos too big", map[string]any{"url": url, "topic": "t/x", "qos": 3}, "mqtt_sink: qos must be 0, 1 or 2"},
		{"qos fractional", map[string]any{"url": url, "topic": "t/x", "qos": 0.5}, "mqtt_sink: qos must be 0, 1 or 2"},
		{"acceptance twin retained", map[string]any{"url": url, "topic": "t/x", "retained": true, "qos": float64(1)}, ""},
	}
	for _, row := range rows {
		t.Run(row.name, func(t *testing.T) {
			block, err := blocks.MQTTSink("mqttsink", row.config, flow.Deps{})
			if row.wantErr != "" {
				if err == nil || err.Error() != row.wantErr {
					t.Fatalf("err = %v, want %q", err, row.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("MQTTSink: %v", err)
			}
			block.(flow.Stopper).Stop()
			block.(flow.Stopper).Stop() // Stop is idempotent
		})
	}
}

func TestMQTTRegistration_InfoAndSchema(t *testing.T) {
	for typ, role := range map[string]blocks.Role{
		blocks.TypeMQTTSource: blocks.RoleSource,
		blocks.TypeMQTTSink:   blocks.RoleSink,
	} {
		info, ok := blocks.LookupInfo(typ)
		if !ok {
			t.Fatalf("LookupInfo(%q) not found", typ)
		}
		if info.Role != role {
			t.Errorf("LookupInfo(%q).Role = %q, want %q", typ, info.Role, role)
		}
		if info.ConfigSchema == nil {
			t.Errorf("LookupInfo(%q).ConfigSchema = nil, want schema", typ)
		}
		reg := func() bool {
			r := blocks.DefaultRegistry()
			for _, registered := range r.Types() {
				if registered == typ {
					return true
				}
			}
			return false
		}
		if !reg() {
			t.Errorf("%q not in DefaultRegistry", typ)
		}
	}
}
