package blocks

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"

	"github.com/astrate-platform/astrate/internal/flow"
)

// MQTT block types (issue #83, astarte_flow parity): mqtt_source subscribes
// to broker topics and emits one binary message per delivery, mqtt_sink
// publishes message payloads to a topic. Both talk plain MQTT over TCP via
// paho.mqtt.golang, like the internal/broker client side.
const (
	TypeMQTTSource = "mqtt_source"
	TypeMQTTSink   = "mqtt_sink"
)

const (
	mqttConnectTimeout    = 10 * time.Second
	mqttSubscribeTimeout  = 10 * time.Second
	mqttPublishTimeout    = 5 * time.Second
	mqttSourceQueueCap    = 256
	mqttDisconnectQuiesce = 250
)

// newMQTTClient builds a connected paho client for url. The returned client
// auto-reconnects in the background; callers own Disconnect.
func newMQTTClient(url, clientID, username, password string, onLost func(error)) (paho.Client, error) {
	if clientID == "" {
		var b [4]byte
		if _, err := rand.Read(b[:]); err != nil {
			return nil, fmt.Errorf("generate client id: %w", err)
		}
		clientID = "astrate-flow-" + hex.EncodeToString(b[:])
	}
	opts := paho.NewClientOptions()
	opts.AddBroker(url)
	opts.SetClientID(clientID)
	if username != "" {
		opts.SetUsername(username)
	}
	if password != "" {
		opts.SetPassword(password)
	}
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(true)
	if onLost != nil {
		opts.SetConnectionLostHandler(func(_ paho.Client, err error) { onLost(err) })
	}

	client := paho.NewClient(opts)
	token := client.Connect()
	if !token.WaitTimeout(mqttConnectTimeout) {
		client.Disconnect(mqttDisconnectQuiesce) // stop background reconnect attempts
		if err := token.Error(); err != nil {
			return nil, fmt.Errorf("connect %s: %w", url, err)
		}
		return nil, fmt.Errorf("connect %s: timed out", url)
	}
	if err := token.Error(); err != nil {
		client.Disconnect(mqttDisconnectQuiesce)
		return nil, fmt.Errorf("connect %s: %w", url, err)
	}
	return client, nil
}

// parseMQTTTopics coerces the topics config entry into []string, accepting a
// []string or a JSON-decoded []any of strings.
func parseMQTTTopics(v any) ([]string, error) {
	switch topics := v.(type) {
	case nil:
		return nil, nil
	case []string:
		out := make([]string, len(topics))
		copy(out, topics)
		for _, topic := range out {
			if topic == "" {
				return nil, fmt.Errorf("topics must contain only non-empty strings")
			}
		}
		return out, nil
	case []any:
		out := make([]string, 0, len(topics))
		for _, item := range topics {
			s, ok := item.(string)
			if !ok || s == "" {
				return nil, fmt.Errorf("topics must contain only non-empty strings")
			}
			out = append(out, s)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("topics must contain only non-empty strings")
	}
}

// parseMQTTQoS reads key from config as an MQTT QoS level (0, 1 or 2),
// applying def when the key is absent or null.
func parseMQTTQoS(config map[string]any) (byte, error) {
	v, ok := config["qos"]
	if !ok || v == nil {
		return 0, nil
	}
	if f, isFloat := v.(float64); isFloat && f != math.Trunc(f) {
		return 0, fmt.Errorf("qos must be 0, 1 or 2")
	}
	n, err := numAsInt64(v)
	if err != nil || n < 0 || n > 2 {
		return 0, fmt.Errorf("qos must be 0, 1 or 2")
	}
	return byte(n), nil
}

func boolConfig(config map[string]any, key string, def bool) bool {
	v, ok := config[key]
	if !ok || v == nil {
		return def
	}
	b, ok := v.(bool)
	if !ok {
		return def
	}
	return b
}

type mqttCommonConfig struct {
	url      string
	qos      byte
	clientID string
	username string
	password string
}

func parseMQTTCommon(config map[string]any) (mqttCommonConfig, error) {
	cfg := mqttCommonConfig{
		url:      stringConfig(config, "url", ""),
		clientID: stringConfig(config, "client_id", ""),
		username: stringConfig(config, "username", ""),
		password: stringConfig(config, "password", ""),
	}
	var err error
	if cfg.qos, err = parseMQTTQoS(config); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// MQTTSource constructs a Source that subscribes to MQTT topics and emits one
// binary message per received publish (issue #83).
//
// Config keys:
//   - url (required string): broker URL, e.g. tcp://127.0.0.1:1883
//   - topics (required string array): at least one topic to subscribe to;
//     empty-string entries are rejected
//   - qos (int, default 0): subscription QoS, must be 0, 1 or 2
//   - client_id, username, password (optional strings)
//
// Construction connects eagerly and subscribes before returning, so config
// and connectivity errors fail fast.
func MQTTSource(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
	common, err := parseMQTTCommon(config)
	if err != nil {
		return nil, fmt.Errorf("mqtt_source: %w", err)
	}
	topics, err := parseMQTTTopics(config["topics"])
	if err != nil {
		return nil, fmt.Errorf("mqtt_source: %w", err)
	}
	if len(topics) == 0 {
		return nil, fmt.Errorf("mqtt_source: at least one topic is required")
	}
	if common.url == "" {
		return nil, fmt.Errorf("mqtt_source: url is required")
	}

	src := &mqttSource{name: name, ch: make(chan paho.Message, mqttSourceQueueCap)}
	client, err := newMQTTClient(common.url, common.clientID, common.username, common.password,
		func(error) { src.lost.Store(true) })
	if err != nil {
		return nil, fmt.Errorf("mqtt_source: %w", err)
	}
	src.client = client

	handler := func(_ paho.Client, m paho.Message) {
		select {
		case src.ch <- m:
		default:
			src.dropped.Add(1) // never block the client callback
		}
	}
	for _, topic := range topics {
		token := client.Subscribe(topic, common.qos, handler)
		if !token.WaitTimeout(mqttSubscribeTimeout) {
			client.Disconnect(mqttDisconnectQuiesce)
			return nil, fmt.Errorf("mqtt_source: subscribe %s: timed out", topic)
		}
		if err := token.Error(); err != nil {
			client.Disconnect(mqttDisconnectQuiesce)
			return nil, fmt.Errorf("mqtt_source: subscribe %s: %w", topic, err)
		}
	}
	return src, nil
}

type mqttSource struct {
	name     string
	client   paho.Client
	ch       chan paho.Message
	dropped  atomic.Int64
	lost     atomic.Bool
	stopOnce sync.Once
}

var (
	_ flow.Block   = (*mqttSource)(nil)
	_ flow.Source  = (*mqttSource)(nil)
	_ flow.Stopper = (*mqttSource)(nil)
)

func (s *mqttSource) Name() string { return s.name }

// Process implements flow.Block for the non-pump path: it drains currently
// buffered deliveries without blocking.
func (s *mqttSource) Process(_ *flow.Message) ([]*flow.Message, error) {
	if s.lost.Load() {
		return nil, fmt.Errorf("mqtt_source: connection lost")
	}
	var out []*flow.Message
	for {
		select {
		case m := <-s.ch:
			out = append(out, mqttSourceMessage(m))
		default:
			return out, nil
		}
	}
}

// Emit implements flow.Source: it blocks until one delivery is buffered or
// ctx is cancelled. If the connection was lost, it fails instead of waiting
// forever for messages that may never arrive.
func (s *mqttSource) Emit(ctx context.Context) ([]*flow.Message, error) {
	if s.lost.Load() {
		return nil, fmt.Errorf("mqtt_source: connection lost")
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case m := <-s.ch:
		return []*flow.Message{mqttSourceMessage(m)}, nil
	}
}

// Stop implements flow.Stopper: Disconnect drops the session and lets paho's
// router goroutines exit. Idempotent; safe if never connected.
func (s *mqttSource) Stop() {
	s.stopOnce.Do(func() {
		if s.client != nil {
			s.client.Disconnect(mqttDisconnectQuiesce)
		}
	})
}

func mqttSourceMessage(m paho.Message) *flow.Message {
	return &flow.Message{
		Key:       m.Topic(),
		Type:      flow.TypeBinary,
		Subtype:   "",
		Data:      m.Payload(),
		Timestamp: time.Now().UnixMicro(),
	}
}

// MQTTSink constructs a Sink that publishes each message payload to an MQTT
// topic (issue #83). Payload mapping follows http_sink's rule: binary sends
// its bytes as-is, strings send their text bytes, everything else is
// JSON-encoded.
//
// Config keys:
//   - url (required string): broker URL
//   - topic (required string): target of every publish
//   - qos (int, default 0), retained (bool, default false)
//   - client_id, username, password (optional strings)
//
// Construction connects eagerly; Stop disconnects on teardown.
func MQTTSink(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
	common, err := parseMQTTCommon(config)
	if err != nil {
		return nil, fmt.Errorf("mqtt_sink: %w", err)
	}
	topic := stringConfig(config, "topic", "")
	retained := boolConfig(config, "retained", false)
	if common.url == "" {
		return nil, fmt.Errorf("mqtt_sink: url is required")
	}
	if topic == "" {
		return nil, fmt.Errorf("mqtt_sink: topic is required")
	}

	client, err := newMQTTClient(common.url, common.clientID, common.username, common.password, nil)
	if err != nil {
		return nil, fmt.Errorf("mqtt_sink: %w", err)
	}
	return &mqttSink{
		name:     name,
		topic:    topic,
		qos:      common.qos,
		retained: retained,
		client:   client,
	}, nil
}

type mqttSink struct {
	name     string
	topic    string
	qos      byte
	retained bool
	client   paho.Client
	stopOnce sync.Once
}

var (
	_ flow.Block   = (*mqttSink)(nil)
	_ flow.Stopper = (*mqttSink)(nil)
)

func (s *mqttSink) Name() string { return s.name }

// Process implements flow.Block: it maps msg to payload bytes and publishes,
// waiting up to 5s for broker acknowledgement.
func (s *mqttSink) Process(msg *flow.Message) ([]*flow.Message, error) {
	if msg == nil {
		return nil, nil
	}
	payload, err := mqttPayload(msg)
	if err != nil {
		return nil, err
	}
	token := s.client.Publish(s.topic, s.qos, s.retained, payload)
	if !token.WaitTimeout(mqttPublishTimeout) {
		if err := token.Error(); err != nil {
			return nil, fmt.Errorf("mqtt_sink: publish: %v", err)
		}
		return nil, fmt.Errorf("mqtt_sink: publish: timed out")
	}
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("mqtt_sink: publish: %v", err)
	}
	return nil, nil
}

// Stop implements flow.Stopper. Idempotent.
func (s *mqttSink) Stop() {
	s.stopOnce.Do(func() {
		if s.client != nil {
			s.client.Disconnect(mqttDisconnectQuiesce)
		}
	})
}

// mqttPayload maps a message to publish body bytes, mirroring http_sink's
// sinkPayload rule without the Content-Type derivation.
func mqttPayload(msg *flow.Message) ([]byte, error) {
	switch msg.Type {
	case flow.TypeBinary:
		if bs, ok := msg.Data.([]byte); ok {
			return bs, nil
		}
	case flow.TypeString:
		if str, ok := msg.Data.(string); ok {
			return []byte(str), nil
		}
	}
	body, err := toJSONBytes(msg)
	if err != nil {
		return nil, fmt.Errorf("mqtt_sink: encode payload: %w", err)
	}
	return body, nil
}
