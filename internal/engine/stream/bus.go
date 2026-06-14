// Package stream is the in-process live fan-out bus (docs/ROADMAP.md §7.2
// file 6.13, docs/DESIGN.md §1.1): the engine publishes every committed
// data operation and device lifecycle event, and consumers — the M7b
// WebSocket/SSE endpoint, tests — subscribe per realm with optional
// filters. Publishing never blocks: a slow consumer's full channel drops
// the event for that consumer with a metric (§1.4 philosophy — live viewers
// must never backpressure ingestion).
package stream

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Event kinds mirror the Astarte trigger event names.
const (
	// KindIncomingData is a committed device data operation (set or unset).
	KindIncomingData = "incoming_data"
	// KindDeviceConnected is a device connection.
	KindDeviceConnected = "device_connected"
	// KindDeviceDisconnected is a device disconnection.
	KindDeviceDisconnected = "device_disconnected"
)

// Event is one live event.
type Event struct {
	// Kind discriminates the event (Kind* constants).
	Kind string
	// Realm is the tenant.
	Realm string
	// DeviceID is the encoded device ID.
	DeviceID string
	// Interface and Path locate data events; empty for lifecycle events.
	Interface string
	Path      string
	// Value is the JSON-friendly rendering of a data event's value (nil for
	// property unset and lifecycle events).
	Value any
	// Timestamp is the event instant (the effective sample timestamp for
	// data events).
	Timestamp time.Time
}

// Filter narrows a subscription; zero-value fields match everything.
type Filter struct {
	// DeviceID keeps only one device's events when set.
	DeviceID string
	// Interface keeps only one interface's data events when set (lifecycle
	// events carry no interface and are kept only by an empty filter).
	Interface string
}

// matches applies the filter.
func (f Filter) matches(ev *Event) bool {
	if f.DeviceID != "" && f.DeviceID != ev.DeviceID {
		return false
	}
	if f.Interface != "" && f.Interface != ev.Interface {
		return false
	}
	return true
}

// DefaultSubscriberBuffer is the per-subscriber channel capacity used when
// Subscribe is called with a non-positive buffer.
const DefaultSubscriberBuffer = 64

// Bus is the fan-out hub. The zero value is not usable; construct with New.
type Bus struct {
	mu     sync.Mutex
	subs   map[int]*subscriber
	nextID int
	closed bool

	dropped prometheus.Counter
}

// subscriber is one registered consumer.
type subscriber struct {
	realm  string
	filter Filter
	ch     chan Event
	once   sync.Once // guards the channel close (cancel vs Close race)
}

// close closes the subscriber channel exactly once.
func (s *subscriber) close() {
	s.once.Do(func() { close(s.ch) })
}

// New builds a bus; a non-nil reg receives its collectors.
func New(reg prometheus.Registerer) *Bus {
	b := &Bus{
		subs: make(map[int]*subscriber),
		dropped: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_engine_stream_dropped_total",
			Help: "Live events dropped because a subscriber's channel was full (docs/DESIGN.md §1.4).",
		}),
	}
	if reg != nil {
		reg.MustRegister(b.dropped)
	}
	return b
}

// Subscribe registers a consumer for one realm's events. The returned
// cancel function unregisters it and closes the channel; the channel also
// closes when the bus shuts down. buffer <= 0 selects
// DefaultSubscriberBuffer. Subscribing to a closed bus returns an
// already-closed channel.
func (b *Bus) Subscribe(realm string, f Filter, buffer int) (<-chan Event, func()) {
	if buffer <= 0 {
		buffer = DefaultSubscriberBuffer
	}
	s := &subscriber{realm: realm, filter: f, ch: make(chan Event, buffer)}

	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		s.close()
		return s.ch, func() {}
	}
	id := b.nextID
	b.nextID++
	b.subs[id] = s
	b.mu.Unlock()

	cancel := func() {
		b.mu.Lock()
		_, present := b.subs[id]
		delete(b.subs, id)
		b.mu.Unlock()
		if present {
			s.close()
		}
	}
	return s.ch, cancel
}

// Publish fans an event out to every matching subscriber without blocking;
// full subscriber channels drop the event with a metric.
func (b *Bus) Publish(ev Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	for _, s := range b.subs {
		if s.realm != ev.Realm || !s.filter.matches(&ev) {
			continue
		}
		select {
		case s.ch <- ev:
		default:
			b.dropped.Inc()
		}
	}
}

// Subscribers reports the number of registered consumers (metrics, tests).
func (b *Bus) Subscribers() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.subs)
}

// Close shuts the bus down: every subscriber channel closes and later
// publishes are no-ops.
func (b *Bus) Close() {
	b.mu.Lock()
	if b.closed {
		b.mu.Unlock()
		return
	}
	b.closed = true
	subs := make([]*subscriber, 0, len(b.subs))
	for _, s := range b.subs {
		subs = append(subs, s)
	}
	b.subs = map[int]*subscriber{}
	b.mu.Unlock()

	for _, s := range subs {
		s.close()
	}
}
