package engine

import (
	"context"
	"errors"
	"time"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/engine/triggers"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

// BrokerPort is the engine's broker-side port (hexagonal-lite,
// docs/DESIGN.md §1.3): server→device publishing and the ACL-relevant
// introspection refresh. Defined on the consumer side so tests substitute
// fakes; AdaptBroker wraps the real broker.
type BrokerPort interface {
	// Publish sends a server-side message (docs/ROADMAP.md §6 file 5.7):
	// retain for properties, per-message expiry for datastreams.
	Publish(topic string, payload []byte, qos byte, retain bool, expiry time.Duration) error
	// RefreshIntrospection reloads a connected device's introspection-derived
	// ACL state after the engine persists a new introspection
	// (docs/ROADMAP.md §7.2 file 6.7).
	RefreshIntrospection(ctx context.Context, realm string, id deviceid.ID) error
}

// brokerFacade adapts *broker.Broker to BrokerPort.
type brokerFacade struct {
	b *broker.Broker
}

// Publish implements BrokerPort.
func (f brokerFacade) Publish(topic string, payload []byte, qos byte, retain bool, expiry time.Duration) error {
	return f.b.Publisher().Publish(topic, payload, qos, retain, expiry)
}

// RefreshIntrospection implements BrokerPort.
func (f brokerFacade) RefreshIntrospection(ctx context.Context, realm string, id deviceid.ID) error {
	return f.b.RefreshIntrospection(ctx, realm, id)
}

// AdaptBroker wraps the embedded broker as the engine's BrokerPort.
func AdaptBroker(b *broker.Broker) BrokerPort {
	return brokerFacade{b: b}
}

// consumerPropertiesSendTimeout bounds the asynchronous purge-message send
// that follows a device connection (docs/DESIGN.md §3.4).
const consumerPropertiesSendTimeout = 30 * time.Second

// New assembles the full engine (docs/ROADMAP.md §7.2 file 6.14): the M6a
// pipeline plus the M6b control-channel handlers, trigger evaluation, and
// the live fan-out bus. bp may be nil at construction time — the broker
// needs the engine as its intake, so M8 wires them in two steps — but must
// be attached (AttachBroker) before Start.
func New(st Store, bp BrokerPort, cfg Config) (*Engine, error) {
	e, err := newEngine(cfg, st)
	if err != nil {
		return nil, err
	}
	e.broker = bp
	e.exec = triggers.NewExecutor(triggers.ExecutorConfig{
		Registerer: cfg.Registerer,
		Logger:     cfg.Logger,
	})
	e.bus = stream.New(cfg.Registerer)

	e.onIntrospection = e.handleIntrospection
	e.onControl = e.handleControl
	e.afterCommit = e.fireCommitted
	e.onLifecycle = e.handleLifecycle
	e.onDeviceError = e.fireDeviceError
	return e, nil
}

// AttachBroker binds the broker port; it must run before Start (M8 wiring:
// engine → broker.New(intake=engine) → AttachBroker → Start).
func (e *Engine) AttachBroker(bp BrokerPort) {
	e.broker = bp
}

// Start loads the schema snapshot, subscribes to interface-change
// notifications, and launches the shard goroutines. ctx bounds the engine's
// background work and must outlive Drain.
func (e *Engine) Start(ctx context.Context) error {
	if e.broker == nil {
		return errors.New("engine: no broker attached (AttachBroker before Start)")
	}
	return e.start(ctx)
}

// Drain gracefully stops the engine (docs/DESIGN.md §5.3): intake is
// refused, shards drain and flush their final batches, asynchronous control
// sends finish, the trigger executor drains its queue, and the live bus
// closes — all bounded by ctx. Call broker.Close first so no submitter
// blocks indefinitely.
func (e *Engine) Drain(ctx context.Context) error {
	err := e.drain(ctx)

	bgDone := make(chan struct{})
	go func() {
		e.bg.Wait()
		close(bgDone)
	}()
	select {
	case <-bgDone:
	case <-ctx.Done():
		err = errors.Join(err, ctx.Err())
	}

	if e.exec != nil {
		err = errors.Join(err, e.exec.Close(ctx))
	}
	if e.bus != nil {
		e.bus.Close()
	}
	return err
}

// Bus exposes the live fan-out bus (the M7b stream socket subscribes
// through it).
func (e *Engine) Bus() *stream.Bus {
	return e.bus
}

// RefreshTriggers rebuilds a realm's compiled snapshot — triggers ride in
// the same snapshot as interfaces, so this is the in-process invalidation
// callback for M7 trigger CRUD (docs/DESIGN.md §2.6).
func (e *Engine) RefreshTriggers(ctx context.Context, realmID int16) error {
	return e.schemas.reloadRealm(ctx, realmID)
}

// fireCommitted is the afterCommit observer (docs/ROADMAP.md §7.1 file
// 6.5): it runs in-shard right after a batch commits, preserving per-device
// order, and feeds data triggers plus the live bus.
func (e *Engine) fireCommitted(ops []PersistOp) {
	for i := range ops {
		op := &ops[i]
		rs := e.schemas.realm(op.Realm)
		if rs == nil {
			continue // realm dropped mid-flight; nothing to evaluate against
		}
		e.fireData(rs, op)
	}
}

// fireData evaluates data triggers for one committed op and publishes it on
// the live bus.
func (e *Engine) fireData(rs *realmSchema, op *PersistOp) {
	value := e.eventValue(op.Value)
	deviceID := op.DeviceID.String()

	if len(rs.triggers) > 0 {
		match := triggers.DataEvent{
			DeviceID:  deviceID,
			Interface: op.Interface.Name,
			Major:     op.Interface.Major,
			Path:      op.Path,
			Value:     op.Value,
		}
		body := triggers.NewIncomingDataEvent(op.Interface.Name, op.Path, value)
		for _, tr := range rs.triggers {
			if tr.MatchesData(match) {
				e.exec.Enqueue(triggers.Delivery{
					Realm:   rs.name,
					Trigger: tr,
					Event: triggers.SimpleEvent{
						Timestamp:   op.TS,
						DeviceID:    deviceID,
						TriggerName: tr.Name,
						Event:       body,
					},
				})
			}
		}
	}

	e.bus.Publish(stream.Event{
		Kind:      stream.KindIncomingData,
		Realm:     rs.name,
		DeviceID:  deviceID,
		Interface: op.Interface.Name,
		Path:      op.Path,
		Value:     value,
		Timestamp: op.TS,
	})
}

// fireDevice evaluates device triggers for one device-scoped event.
func (e *Engine) fireDevice(rs *realmSchema, id deviceid.ID, at time.Time, match triggers.DeviceEvent, body any) {
	for _, tr := range rs.triggers {
		if tr.MatchesDevice(match) {
			e.exec.Enqueue(triggers.Delivery{
				Realm:   rs.name,
				Trigger: tr,
				Event: triggers.SimpleEvent{
					Timestamp:   at,
					DeviceID:    id.String(),
					TriggerName: tr.Name,
					Event:       body,
				},
			})
		}
	}
}

// fireDeviceError is the onDeviceError observer: every pipeline rejection
// becomes a device_error trigger event (docs/DESIGN.md §2.6 "failures are
// never silent").
func (e *Engine) fireDeviceError(m broker.InboundMessage, reason, detail string) {
	rs := e.schemas.realm(m.Realm)
	if rs == nil || len(rs.triggers) == 0 {
		return
	}
	e.fireDevice(rs, m.DeviceID, m.ReceivedAt,
		triggers.DeviceEvent{DeviceID: m.DeviceID.String(), On: triggers.OnDeviceError},
		triggers.NewDeviceErrorEvent(reason, map[string]string{"detail": detail}))
}

// handleLifecycle is the onLifecycle observer (after the built-in cache
// eviction): device_connected / device_disconnected triggers, live bus
// events, and — on connect — the asynchronous `consumer/properties` purge
// send. A (re)connecting device always receives the full current truth, a
// safe superset of the session_present=0 requirement of docs/DESIGN.md §3.4
// (the purge list is idempotent). The send runs off the broker's connection
// goroutine: LifecycleSink implementations must not block.
func (e *Engine) handleLifecycle(ev broker.LifecycleEvent) {
	rs := e.schemas.realm(ev.Realm)
	if rs == nil {
		return
	}
	deviceID := ev.DeviceID.String()

	switch ev.Type {
	case broker.EventDeviceConnected:
		ip := ""
		if ev.RemoteIP.IsValid() {
			ip = ev.RemoteIP.String()
		}
		e.fireDevice(rs, ev.DeviceID, ev.At,
			triggers.DeviceEvent{DeviceID: deviceID, On: triggers.OnDeviceConnected},
			triggers.NewDeviceConnectedEvent(ip))
		e.bus.Publish(stream.Event{
			Kind: stream.KindDeviceConnected, Realm: rs.name, DeviceID: deviceID, Timestamp: ev.At,
		})

		e.bg.Add(1)
		go func() {
			defer e.bg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), consumerPropertiesSendTimeout)
			defer cancel()
			if err := e.sendConsumerProperties(ctx, rs, ev.DeviceID); err != nil {
				e.log.Warn("consumer/properties send after connect failed",
					"realm", rs.name, "device", deviceID, "err", err)
			}
		}()
	case broker.EventDeviceDisconnected:
		e.fireDevice(rs, ev.DeviceID, ev.At,
			triggers.DeviceEvent{DeviceID: deviceID, On: triggers.OnDeviceDisconnected},
			triggers.NewDeviceDisconnectedEvent())
		e.bus.Publish(stream.Event{
			Kind: stream.KindDeviceDisconnected, Realm: rs.name, DeviceID: deviceID, Timestamp: ev.At,
		})
	}
}

// eventValue renders a decoded payload value into its canonical
// JSON-friendly form for trigger events and the live bus; nil stays nil
// (property unset renders as JSON null, upstream parity).
func (e *Engine) eventValue(v any) any {
	if v == nil {
		return nil
	}
	jv, err := jsonable(v)
	if err != nil {
		// Impossible for pipeline-validated values; be loud, deliver null.
		e.met.internalErrors.Inc()
		e.log.Error("persisted value does not render as JSON", "err", err)
		return nil
	}
	return jv
}
