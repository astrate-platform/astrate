package engine

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/astrate-platform/astrate/internal/broker"
	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/engine/triggers"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// Configuration defaults (docs/DESIGN.md §1.4).
const (
	// DefaultShards is the default shard count.
	DefaultShards = 16
	// DefaultShardQueue is the default per-shard channel capacity.
	DefaultShardQueue = 4096
	// DefaultBatchMaxRows flushes a shard batch when it reaches this size.
	DefaultBatchMaxRows = 64
	// DefaultBatchMaxWait flushes a non-empty shard batch after this delay.
	DefaultBatchMaxWait = 50 * time.Millisecond
)

// Config carries the engine's operational knobs (TOML wiring lands in M8).
type Config struct {
	// Shards is the number of ordered pipeline shards (default
	// DefaultShards). Messages of one device always land on the same shard.
	Shards int
	// ShardQueue is the per-shard channel capacity (default
	// DefaultShardQueue). A full shard blocks QoS >= 1 submits (deferred-ack
	// backpressure) and drops QoS 0 messages with a metric (§1.4).
	ShardQueue int
	// BatchMaxRows is the micro-batch row cap (default DefaultBatchMaxRows).
	BatchMaxRows int
	// BatchMaxWait is the micro-batch time cap (default DefaultBatchMaxWait).
	BatchMaxWait time.Duration
	// MaxPayloadBytes caps accepted data payload size for both formats
	// (default payload.DefaultMaxSize).
	MaxPayloadBytes int
	// Registerer receives the engine's Prometheus collectors; nil leaves
	// them unregistered (the collectors still work, which tests rely on).
	Registerer prometheus.Registerer
	// Logger receives engine logs (default slog.Default()).
	Logger *slog.Logger
}

// Engine is the sharded ingestion pipeline (docs/ROADMAP.md §7.1 file 6.2).
// It implements broker.Intake and broker.LifecycleSink. M6a builds the data
// path; the M6b handlers (introspection, control, triggers, stream) attach
// to the seam fields below without touching the pipeline.
type Engine struct {
	cfg     Config
	st      Store
	log     *slog.Logger
	met     *metrics
	schemas *schemaCache
	devices *deviceCache
	dec     payload.Decoder

	shards  []*shard
	shardWG sync.WaitGroup

	// mu guards closed against in-flight Submits: drain takes the write
	// lock only after every blocked submit has completed.
	mu     sync.RWMutex
	closed bool

	// quit is closed when drain begins: parked flush retries give up
	// (leaving their messages unacknowledged — at-least-once, docs/DESIGN.md
	// §5.3) so blocked submitters and shard loops can wind down.
	quit     chan struct{}
	quitOnce sync.Once

	// broker is the server→device publish + introspection-refresh port,
	// wired by New/AttachBroker (M6b file 6.14).
	broker BrokerPort
	// exec executes matched trigger actions (M6b file 6.12); set by New.
	exec *triggers.Executor
	// bus is the live fan-out hub (M6b file 6.13); set by New.
	bus *stream.Bus
	// bg tracks asynchronous control sends (consumer/properties after a
	// connect) so Drain can wait them out.
	bg sync.WaitGroup

	// onIntrospection handles device introspection publishes (M6b file 6.7).
	// While nil, such messages are acknowledged, counted as unhandled, and
	// dropped — the safe default until the handler is wired.
	onIntrospection func(ctx context.Context, m broker.InboundMessage, realm *realmSchema)
	// onControl handles `/control/...` publishes (M6b file 6.8); subpath is
	// the topic remainder after "control/" (e.g. "emptyCache"). Nil behaves
	// like onIntrospection.
	onControl func(ctx context.Context, m broker.InboundMessage, realm *realmSchema, subpath string)
	// afterCommit observes every persisted op right after its batch commits,
	// in-shard, preserving per-device order. M6b wires trigger evaluation
	// and live fan-out here (docs/ROADMAP.md §7.2). Nil means no observers.
	afterCommit func(ops []PersistOp)
	// onLifecycle observes broker lifecycle events after the built-in cache
	// eviction (M6b wires device connect/disconnect triggers here).
	onLifecycle func(ev broker.LifecycleEvent)
	// onDeviceError observes every rejection (M6b wires device_error trigger
	// events here, docs/DESIGN.md §2.6). Nil means metrics and logs only.
	onDeviceError func(m broker.InboundMessage, reason, detail string)
}

// Compile-time port guards.
var (
	_ broker.Intake        = (*Engine)(nil)
	_ broker.LifecycleSink = (*Engine)(nil)
)

// shard is one ordered pipeline lane: a bounded channel, the lane goroutine,
// and its micro-batcher.
type shard struct {
	idx   int
	ch    chan broker.InboundMessage
	batch *batcher

	// timer drives the BatchMaxWait flush; armed tracks whether timer.C is
	// live (single-goroutine state, owned by the shard loop).
	timer *time.Timer
	armed bool
}

// newEngine validates cfg and assembles the pipeline. The public
// constructor (New, M6b file 6.14) wraps it; M6a tests call it directly.
func newEngine(cfg Config, st Store) (*Engine, error) {
	if st == nil {
		return nil, errors.New("engine: store is required")
	}
	if cfg.Shards == 0 {
		cfg.Shards = DefaultShards
	}
	if cfg.ShardQueue == 0 {
		cfg.ShardQueue = DefaultShardQueue
	}
	if cfg.BatchMaxRows == 0 {
		cfg.BatchMaxRows = DefaultBatchMaxRows
	}
	if cfg.BatchMaxWait == 0 {
		cfg.BatchMaxWait = DefaultBatchMaxWait
	}
	if cfg.Shards < 0 || cfg.ShardQueue < 0 || cfg.BatchMaxRows < 0 || cfg.BatchMaxWait < 0 || cfg.MaxPayloadBytes < 0 {
		return nil, fmt.Errorf("engine: negative configuration value: %+v", cfg)
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	log := cfg.Logger.With("component", "engine")

	e := &Engine{
		cfg:     cfg,
		st:      st,
		log:     log,
		schemas: newSchemaCache(st, log),
		devices: newDeviceCache(st),
		dec:     payload.Decoder{MaxSize: cfg.MaxPayloadBytes},
		quit:    make(chan struct{}),
	}
	e.shards = make([]*shard, cfg.Shards)
	for i := range e.shards {
		sh := &shard{
			idx:   i,
			ch:    make(chan broker.InboundMessage, cfg.ShardQueue),
			timer: time.NewTimer(time.Hour),
		}
		if !sh.timer.Stop() {
			<-sh.timer.C
		}
		sh.batch = newBatcher(e, i)
		e.shards[i] = sh
	}
	e.met = newMetrics(cfg.Registerer, e.shards)
	return e, nil
}

// start loads the schema snapshot, subscribes to interface-change
// notifications, and launches the shard goroutines. ctx bounds the engine's
// background work; cancel it only after drain.
func (e *Engine) start(ctx context.Context) error {
	if err := e.schemas.loadAll(ctx); err != nil {
		return err
	}
	notifications, err := e.st.Listen(ctx, store.ChannelInterfaces)
	if err != nil {
		return fmt.Errorf("engine: subscribing to %s: %w", store.ChannelInterfaces, err)
	}
	go e.runInvalidation(ctx, notifications)
	for _, sh := range e.shards {
		e.shardWG.Add(1)
		go e.runShard(ctx, sh)
	}
	return nil
}

// Submit implements broker.Intake (docs/DESIGN.md §1.4): route to the
// device's shard; block when the shard is full and the message is QoS >= 1
// (the broker withholds the device's PUBACK, which is the backpressure);
// drop QoS 0 messages on a full shard with a metric.
func (e *Engine) Submit(m broker.InboundMessage) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	e.met.messages.Inc()
	if e.closed {
		// Shutting down: drop without acknowledging, so QoS >= 1 messages
		// are re-sent after reconnect (docs/DESIGN.md §5.3).
		e.met.droppedShutdown.Inc()
		return
	}
	sh := e.shards[shardOf(m.DeviceID[:], len(e.shards))]
	if m.QoS == 0 {
		select {
		case sh.ch <- m:
		default:
			e.met.qos0Drops.Inc()
		}
		return
	}
	sh.ch <- m
}

// shardOf maps a device ID to its shard with inline FNV-1a (allocation-free
// on the hot path).
func shardOf(id []byte, n int) int {
	const (
		offset64 = 14695981039346656037
		prime64  = 1099511628211
	)
	h := uint64(offset64)
	for _, b := range id {
		h ^= uint64(b)
		h *= prime64
	}
	return int(h % uint64(n)) // #nosec G115 -- value already reduced mod n
}

// OnLifecycleEvent implements broker.LifecycleSink: device disconnects evict
// the per-device cache entry (docs/ROADMAP.md §7.1 file 6.1); the M6b seam
// then observes the event (device connect/disconnect triggers). It never
// blocks (broker.LifecycleSink contract).
func (e *Engine) OnLifecycleEvent(ev broker.LifecycleEvent) {
	if ev.Type == broker.EventDeviceDisconnected {
		e.devices.evict(ev.Realm, ev.DeviceID)
	}
	if e.onLifecycle != nil {
		e.onLifecycle(ev)
	}
}

// RefreshInterfaces is the in-process cache-invalidation callback
// (docs/DESIGN.md §2.6): realm management (M7) calls it after interface
// CRUD, complementing the LISTEN/NOTIFY path.
func (e *Engine) RefreshInterfaces(ctx context.Context, realmID int16) error {
	return e.schemas.reloadRealm(ctx, realmID)
}

// runInvalidation consumes interface-change notifications until the channel
// closes (store.Listen closes it when ctx is cancelled).
func (e *Engine) runInvalidation(ctx context.Context, ch <-chan store.Notification) {
	for n := range ch {
		id, err := strconv.ParseInt(n.Payload, 10, 16)
		if err != nil {
			e.log.Warn("interface notification with non-realm payload; full reload",
				"payload", n.Payload)
			if err := e.schemas.loadAll(ctx); err != nil {
				e.log.Error("full schema reload failed", "err", err)
			}
			continue
		}
		if err := e.schemas.reloadRealm(ctx, int16(id)); err != nil {
			e.log.Error("schema reload failed", "realm_id", id, "err", err)
		}
	}
}

// runShard is one shard goroutine: it restarts its processing loop after a
// recovered panic (docs/ROADMAP.md §7.1 file 6.2) and exits when the loop
// reports completion.
func (e *Engine) runShard(ctx context.Context, sh *shard) {
	defer e.shardWG.Done()
	for e.shardPass(ctx, sh) {
	}
}

// shardPass runs the shard loop until the channel closes or ctx is
// cancelled (restart=false), or a panic is recovered (restart=true). The
// panicking message is dropped unacknowledged: QoS >= 1 senders re-deliver.
func (e *Engine) shardPass(ctx context.Context, sh *shard) (restart bool) {
	defer func() {
		if r := recover(); r != nil {
			e.met.internalErrors.Inc()
			e.log.Error("shard panic recovered", "shard", sh.idx, "panic", r,
				"stack", string(debug.Stack()))
			restart = true
		}
	}()

	// A pending batch surviving a panic restart must still flush on time.
	if sh.batch.size() > 0 && !sh.armed {
		sh.timer.Reset(e.cfg.BatchMaxWait)
		sh.armed = true
	}

	for {
		select {
		case m, ok := <-sh.ch:
			if !ok {
				e.flushShard(ctx, sh)
				return false
			}
			e.handle(ctx, sh, m)
			switch {
			case sh.batch.size() >= e.cfg.BatchMaxRows:
				e.flushShard(ctx, sh)
			case sh.batch.size() > 0 && !sh.armed:
				sh.timer.Reset(e.cfg.BatchMaxWait)
				sh.armed = true
			}
		case <-sh.timer.C:
			sh.armed = false
			e.flushShard(ctx, sh)
		case <-ctx.Done():
			return false
		}
	}
}

// flushShard disarms the flush timer and flushes the shard's batch.
func (e *Engine) flushShard(ctx context.Context, sh *shard) {
	if sh.armed {
		if !sh.timer.Stop() {
			<-sh.timer.C
		}
		sh.armed = false
	}
	sh.batch.flush(ctx)
}

// drain gracefully stops the pipeline (docs/DESIGN.md §5.3): new submits are
// refused (unacknowledged), shards process their queued messages and flush
// their final batches, bounded by ctx. The public Drain (M6b file 6.14)
// wraps it. Call broker.Close first so no submitter blocks indefinitely.
func (e *Engine) drain(ctx context.Context) error {
	// 1. Tell parked flush loops to give up so saturated shards resume
	//    consuming and blocked submitters can finish.
	e.quitOnce.Do(func() { close(e.quit) })

	// 2. Wait out in-flight submits, then close the lanes.
	e.mu.Lock()
	if !e.closed {
		e.closed = true
		for _, sh := range e.shards {
			close(sh.ch)
		}
	}
	e.mu.Unlock()

	// 3. Wait for the shard goroutines, bounded by ctx.
	done := make(chan struct{})
	go func() {
		e.shardWG.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("engine: drain interrupted: %w", ctx.Err())
	}
}

// metrics are the engine's Prometheus collectors (docs/DESIGN.md §5.2).
// They are constructed unregistered and attached to a Registerer when one is
// configured, so tests can read them without a registry.
type metrics struct {
	messages        prometheus.Counter
	persistOps      *prometheus.CounterVec
	rejects         *prometheus.CounterVec
	qos0Drops       prometheus.Counter
	droppedShutdown prometheus.Counter
	internalErrors  prometheus.Counter
	unhandled       *prometheus.CounterVec
	flushSeconds    prometheus.Histogram
	flushRetries    prometheus.Counter
}

// newMetrics builds (and, with a non-nil reg, registers) the collectors,
// pre-registering every reject-reason and op-kind label so dashboards see
// zeroes instead of absent series.
func newMetrics(reg prometheus.Registerer, shards []*shard) *metrics {
	m := &metrics{
		messages: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_engine_messages_total",
			Help: "Inbound device messages submitted to the engine.",
		}),
		persistOps: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "astrate_engine_persist_ops_total",
			Help: "Validated operations committed to storage, by kind.",
		}, []string{"kind"}),
		rejects: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "astrate_engine_rejects_total",
			Help: "Messages rejected by the validation pipeline, by reason (docs/DESIGN.md §2.6).",
		}, []string{"reason"}),
		qos0Drops: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_engine_qos0_dropped_total",
			Help: "QoS 0 messages dropped because their shard was full (§1.4).",
		}),
		droppedShutdown: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_engine_dropped_shutdown_total",
			Help: "Messages refused (unacknowledged) because the engine was draining.",
		}),
		internalErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_engine_internal_errors_total",
			Help: "Recovered shard panics and other engine-side faults.",
		}),
		unhandled: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "astrate_engine_unhandled_total",
			Help: "Messages acknowledged and dropped because no handler is wired (pre-M6b).",
		}, []string{"kind"}),
		flushSeconds: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "astrate_engine_batch_flush_seconds",
			Help:    "Micro-batch flush latency, including retries.",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
		}),
		flushRetries: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_engine_batch_flush_retries_total",
			Help: "Failed flush attempts that were retried (DB-outage parking, §5.3).",
		}),
	}
	for _, r := range payload.RejectReasons() {
		m.rejects.WithLabelValues(r.String())
	}
	for _, r := range engineRejectReasons {
		m.rejects.WithLabelValues(r)
	}
	for _, k := range opKinds {
		m.persistOps.WithLabelValues(k.String())
	}
	m.unhandled.WithLabelValues("introspection")
	m.unhandled.WithLabelValues("control")

	if reg != nil {
		reg.MustRegister(m.messages, m.persistOps, m.rejects, m.qos0Drops,
			m.droppedShutdown, m.internalErrors, m.unhandled, m.flushSeconds,
			m.flushRetries, &shardDepthCollector{shards: shards})
	}
	return m
}

// shardDepthCollector exports per-shard queue depth as a gauge, sampled at
// scrape time.
type shardDepthCollector struct {
	shards []*shard
}

// shardDepthDesc describes the gauge.
var shardDepthDesc = prometheus.NewDesc(
	"astrate_engine_shard_depth",
	"Messages queued per pipeline shard.",
	[]string{"shard"}, nil,
)

// Describe implements prometheus.Collector.
func (c *shardDepthCollector) Describe(ch chan<- *prometheus.Desc) { ch <- shardDepthDesc }

// Collect implements prometheus.Collector.
func (c *shardDepthCollector) Collect(ch chan<- prometheus.Metric) {
	for _, sh := range c.shards {
		ch <- prometheus.MustNewConstMetric(shardDepthDesc, prometheus.GaugeValue,
			float64(len(sh.ch)), strconv.Itoa(sh.idx))
	}
}
