package broker

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// Listener IDs.
const (
	listenerMTLS = "mtls"
	listenerDev  = "dev"
)

// Config defaults.
const (
	// DefaultTLSAddr is the standard Astarte broker port.
	DefaultTLSAddr = ":8883"
	// DefaultDevAddr is the plaintext development listener address, bound
	// only when InsecureDevMode is set.
	DefaultDevAddr = ":1883"
	// DefaultMaxPacketBytes caps inbound MQTT packets: the 64 KiB payload
	// bound (docs/DESIGN.md §3.5.3, §4.5) plus generous topic/header room.
	DefaultMaxPacketBytes = 128 * 1024
)

// Config carries the broker's operational knobs (TOML wiring lands in M8).
type Config struct {
	// TLSAddr is the mTLS listener address (default DefaultTLSAddr).
	TLSAddr string
	// ServerTLSCert is the broker's server-side TLS identity, required for
	// the TLS listener. Deployments issue it from the realm CA or any CA
	// the device fleet trusts (docs/DESIGN.md §4.4 flow C delivers ca_crt).
	ServerTLSCert tls.Certificate
	// InsecureDevMode additionally binds a plaintext listener that
	// authenticates by claimed client ID alone — local development only
	// (docs/DESIGN.md §3.1).
	InsecureDevMode bool
	// DevAddr is the plaintext listener address (default DefaultDevAddr;
	// ignored unless InsecureDevMode).
	DevAddr string
	// SessionStorePath is the bbolt file persisting sessions across
	// restarts (docs/DESIGN.md §3.1). Required.
	SessionStorePath string
	// EnforceLatestCert rejects connections presenting a certificate older
	// than the device's latest issuance (pairing.enforce_latest_cert,
	// docs/DESIGN.md §4.3).
	EnforceLatestCert bool
	// MaxPacketBytes caps inbound MQTT packet size (default
	// DefaultMaxPacketBytes).
	MaxPacketBytes uint32
	// Logger receives broker and hook logs (default slog.Default()).
	Logger *slog.Logger
}

// Broker is the embedded MQTT broker (docs/ROADMAP.md §6 file 5.8): mochi
// server, Astarte hooks, persistent session store, and inline publisher.
type Broker struct {
	cfg      Config
	srv      *mqtt.Server
	st       Store
	pools    *realmPools
	registry *sessionRegistry
	pub      *Publisher
	log      *slog.Logger

	tlsListener *listeners.TCP
	devListener *listeners.TCP

	closing   chan struct{}
	closeOnce sync.Once
}

// New assembles the broker: hooks registered, listeners bound (so the
// addresses are known), realm CA pools loaded. Call Start to begin serving.
// intake must be non-nil; sink may be nil.
func New(ctx context.Context, cfg Config, st Store, intake Intake, sink LifecycleSink) (*Broker, error) {
	if st == nil || intake == nil {
		return nil, errors.New("broker: store and intake are required")
	}
	if cfg.SessionStorePath == "" {
		return nil, errors.New("broker: SessionStorePath is required")
	}
	if cfg.TLSAddr == "" {
		cfg.TLSAddr = DefaultTLSAddr
	}
	if cfg.DevAddr == "" {
		cfg.DevAddr = DefaultDevAddr
	}
	if cfg.MaxPacketBytes == 0 {
		cfg.MaxPacketBytes = DefaultMaxPacketBytes
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	if len(cfg.ServerTLSCert.Certificate) == 0 {
		return nil, errors.New("broker: ServerTLSCert is required for the TLS listener")
	}
	log := cfg.Logger.With("component", "broker")

	caps := mqtt.NewDefaultServerCapabilities()
	caps.MaximumPacketSize = cfg.MaxPacketBytes
	srv := mqtt.New(&mqtt.Options{
		InlineClient: true,
		Capabilities: caps,
		Logger:       log,
	})

	b := &Broker{
		cfg:      cfg,
		srv:      srv,
		st:       st,
		registry: newSessionRegistry(),
		log:      log,
		closing:  make(chan struct{}),
	}
	b.pools = newRealmPools(st, cfg.ServerTLSCert)
	if err := b.pools.Reload(ctx); err != nil {
		return nil, fmt.Errorf("broker: loading realm CA pools: %w", err)
	}

	sessions := newSessionStore(cfg.SessionStorePath, log)
	sessions.attach(srv)

	devListenerID := ""
	if cfg.InsecureDevMode {
		devListenerID = listenerDev
	}
	hooks := []mqtt.Hook{
		sessions,
		&authHook{
			st: st, pools: b.pools, registry: b.registry,
			enforceLatestCert: cfg.EnforceLatestCert,
			devListenerID:     devListenerID,
			log:               log,
		},
		&aclHook{
			st: st, registry: b.registry,
			offline: newOfflineACL(st, b.pools, log),
			log:     log,
		},
		&intakeHook{registry: b.registry, intake: intake, closing: b.closing, log: log},
		&lifecycleHook{st: st, sink: sink, registry: b.registry, log: log, now: time.Now},
	}
	for _, h := range hooks {
		if err := srv.AddHook(h, nil); err != nil {
			return nil, fmt.Errorf("broker: registering hook %s: %w", h.ID(), err)
		}
	}

	b.tlsListener = listeners.NewTCP(listeners.Config{
		ID:        listenerMTLS,
		Address:   cfg.TLSAddr,
		TLSConfig: b.pools.handshakeConfig(),
	})
	if err := srv.AddListener(b.tlsListener); err != nil {
		return nil, fmt.Errorf("broker: binding TLS listener on %s: %w", cfg.TLSAddr, err)
	}
	if cfg.InsecureDevMode {
		b.devListener = listeners.NewTCP(listeners.Config{
			ID:      listenerDev,
			Address: cfg.DevAddr,
		})
		if err := srv.AddListener(b.devListener); err != nil {
			return nil, fmt.Errorf("broker: binding dev listener on %s: %w", cfg.DevAddr, err)
		}
		log.Warn("insecure_dev_mode plaintext listener enabled", "addr", b.devListener.Address())
	}

	b.pub = newPublisher(srv)
	return b, nil
}

// Start restores persisted session state and begins serving the bound
// listeners. It does not block.
func (b *Broker) Start() error {
	return b.srv.Serve()
}

// Close gracefully stops the broker: blocked publish acknowledgments are
// released (their messages stay unacknowledged on the devices, which re-send
// after reconnecting — at-least-once, docs/DESIGN.md §5.3), clients are
// disconnected, and the session store is flushed and closed.
func (b *Broker) Close() error {
	var err error
	b.closeOnce.Do(func() {
		close(b.closing)
		err = b.srv.Close()
	})
	return err
}

// TLSAddr returns the bound mTLS listener address (useful with ":0").
func (b *Broker) TLSAddr() string { return b.tlsListener.Address() }

// DevAddr returns the bound plaintext listener address, or "" when
// insecure_dev_mode is off.
func (b *Broker) DevAddr() string {
	if b.devListener == nil {
		return ""
	}
	return b.devListener.Address()
}

// Publisher returns the inline publishing facade for the engine/AppEngine.
func (b *Broker) Publisher() *Publisher { return b.pub }

// ReloadRealms rebuilds the per-realm client-CA pools from the store. Realm
// CRUD (M7 housekeeping) calls it after creating, deleting, or re-keying a
// realm; new TLS handshakes pick the change up immediately.
func (b *Broker) ReloadRealms(ctx context.Context) error {
	return b.pools.Reload(ctx)
}

// RefreshIntrospection reloads a connected device's introspection-derived
// ACL state. The engine calls it after persisting a new introspection
// (docs/ROADMAP.md §7.2 file 6.7). Unknown or disconnected devices are a
// no-op: their state loads fresh on the next connect or delivery.
func (b *Broker) RefreshIntrospection(ctx context.Context, realm string, id deviceid.ID) error {
	sess := b.registry.get(Identity{Realm: realm, DeviceID: id}.CN())
	if sess == nil {
		return nil
	}
	return sess.refresh(ctx, b.st, b.log)
}

// intakeHook bridges device publishes into the engine intake with
// deferred-ack wiring (docs/ROADMAP.md §6 file 5.8): for QoS >= 1 the
// PUBACK/PUBREC is withheld until the intake calls Ack, propagating
// persistence-commit ordering and shard backpressure to the device
// (docs/DESIGN.md §1.4, §5.3). Device publishes are consumed here — they
// reach subscribers only through what the engine re-publishes (§3.2).
type intakeHook struct {
	mqtt.HookBase
	registry *sessionRegistry
	intake   Intake
	closing  <-chan struct{}
	log      *slog.Logger
}

// ID implements mqtt.Hook.
func (h *intakeHook) ID() string { return "astrate-intake" }

// Provides implements mqtt.Hook.
func (h *intakeHook) Provides(b byte) bool { return b == mqtt.OnPublish }

// OnPublish implements mqtt.Hook.
func (h *intakeHook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	if cl.Net.Inline {
		return pk, nil // server-side publishes fan out normally
	}
	sess := h.registry.get(cl.ID)
	if sess == nil || sess.client != cl {
		// No authenticated session (cannot normally happen): reject without
		// acknowledging, so nothing is silently lost.
		h.log.Error("publish without session", "client", cl.ID, "topic", pk.TopicName)
		return pk, packets.ErrRejectPacket
	}

	msg := InboundMessage{
		Realm:      sess.identity.Realm,
		DeviceID:   sess.identity.DeviceID,
		Topic:      pk.TopicName,
		Payload:    bytes.Clone(pk.Payload),
		QoS:        pk.FixedHeader.Qos,
		ReceivedAt: time.Now(),
	}

	if pk.FixedHeader.Qos == 0 {
		msg.Ack = func() {}
		h.intake.Submit(msg)
		return pk, packets.CodeSuccessIgnore
	}

	acked := make(chan struct{})
	var once sync.Once
	msg.Ack = func() { once.Do(func() { close(acked) }) }
	h.intake.Submit(msg)

	select {
	case <-acked:
		// CodeSuccessIgnore: mochi sends the PUBACK/PUBREC now but does not
		// retain or fan the device's message out to subscribers.
		return pk, packets.CodeSuccessIgnore
	case <-h.closing:
		// Shutting down before the engine acknowledged: reject the packet
		// so no acknowledgment is sent; the device re-sends on reconnect.
		return pk, packets.ErrRejectPacket
	}
}
