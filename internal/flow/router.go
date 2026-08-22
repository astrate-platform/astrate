package flow

import (
	"context"
	"log/slog"
	"runtime/debug"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// DefaultLaneCount is the default number of processing lanes.
const DefaultLaneCount = 16

// DefaultLaneCapacity is the default per-lane channel capacity.
const DefaultLaneCapacity = 256

// OverflowPolicy controls what happens when a lane's channel is full.
type OverflowPolicy uint8

const (
	// OverflowBlock blocks the caller until space is available (QoS ≥ 1).
	OverflowBlock OverflowPolicy = iota
	// OverflowDrop discards the message without blocking (QoS 0).
	OverflowDrop
)

// RouterConfig configures a Router.
type RouterConfig struct {
	// Lanes is the number of processing lanes (default DefaultLaneCount).
	Lanes int
	// LaneCapacity is the per-lane channel capacity (default
	// DefaultLaneCapacity).
	LaneCapacity int
	// QoS0Overflow is the policy when a QoS 0 message targets a full lane.
	QoS0Overflow OverflowPolicy
	// QoS1Overflow is the policy when a QoS ≥ 1 message targets a full lane.
	QoS1Overflow OverflowPolicy
	// Registerer receives Prometheus collectors; nil leaves them
	// unregistered.
	Registerer prometheus.Registerer
	// Logger receives router logs (default slog.Default()).
	Logger *slog.Logger
}

// Router is the stream-based message router. It accepts FlowMessages, hashes
// their Key to a lane, and processes them through the block graph. Messages
// with the same Key are always processed in submission order; different Keys
// may interleave across lanes.
type Router struct {
	cfg   RouterConfig
	log   *slog.Logger
	graph *BlockGraph
	met   *routerMetrics

	lanes  []*lane
	laneWG sync.WaitGroup

	mu     sync.Mutex
	closed bool
	quit   chan struct{}
}

// lane is one processing lane: a bounded channel and a goroutine that drains
// it through the block graph.
type lane struct {
	ch chan *flowMsg
}

// flowMsg pairs a message with its submission-time QoS and acknowledgement.
type flowMsg struct {
	msg    *Message
	qos    byte
	onDrop func()
}

// NewRouter builds a router that feeds every submitted message through graph.
func NewRouter(graph *BlockGraph, cfg RouterConfig, reg prometheus.Registerer) *Router {
	if cfg.Lanes <= 0 {
		cfg.Lanes = DefaultLaneCount
	}
	if cfg.LaneCapacity <= 0 {
		cfg.LaneCapacity = DefaultLaneCapacity
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	r := &Router{
		cfg:   cfg,
		log:   cfg.Logger.With("component", "flow-router"),
		graph: graph,
		quit:  make(chan struct{}),
	}
	r.lanes = make([]*lane, cfg.Lanes)
	for i := range r.lanes {
		r.lanes[i] = &lane{ch: make(chan *flowMsg, cfg.LaneCapacity)}
	}
	r.met = newRouterMetrics(reg)
	return r
}

// Run starts the lane goroutines. Call it after NewRouter. ctx bounds
// background work; cancel to initiate drain.
func (r *Router) Run(ctx context.Context) {
	for _, l := range r.lanes {
		r.laneWG.Add(1)
		go r.runLane(ctx, l)
	}
}

// Submit routes msg to the lane determined by FNV-1a(msg.Key). Behaviour
// depends on qos and the configured overflow policies.
func (r *Router) Submit(msg *Message, qos byte) {
	r.met.submitted.Inc()
	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return
	}
	r.mu.Unlock()

	l := r.lanes[laneOf(msg.Key, len(r.lanes))]
	fm := &flowMsg{msg: msg, qos: qos}

	if qos == 0 && r.cfg.QoS0Overflow == OverflowDrop {
		select {
		case l.ch <- fm:
		default:
			r.met.droppedQoS0.Inc()
			if fm.onDrop != nil {
				fm.onDrop()
			}
		}
		return
	}

	if qos >= 1 && r.cfg.QoS1Overflow == OverflowDrop {
		select {
		case l.ch <- fm:
		default:
			r.met.droppedQoS1.Inc()
			if fm.onDrop != nil {
				fm.onDrop()
			}
		}
		return
	}

	l.ch <- fm
}

// Drain stops accepting new messages, lets lanes finish, and waits for all
// lane goroutines to exit.
func (r *Router) Drain(ctx context.Context) error {
	r.mu.Lock()
	if !r.closed {
		r.closed = true
		close(r.quit)
		for _, l := range r.lanes {
			close(l.ch)
		}
	}
	r.mu.Unlock()

	done := make(chan struct{})
	go func() {
		r.laneWG.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// runLane drains one lane's channel through the block graph, recovering
// panics per-message so one block's bug cannot crash the router.
func (r *Router) runLane(_ context.Context, l *lane) {
	defer r.laneWG.Done()
	for fm := range l.ch {
		r.processOne(fm)
	}
}

// processOne feeds a single message through the block graph, catching
// panics. Block-level errors are logged and counted but never crash the
// lane.
func (r *Router) processOne(fm *flowMsg) {
	defer func() {
		if rec := recover(); rec != nil {
			r.met.blockErrors.Inc()
			r.log.Error("panic in block graph", "key", fm.msg.Key,
				"panic", rec, "stack", string(debug.Stack()))
		}
	}()
	_, err := r.graph.Run(fm.msg)
	if err != nil {
		r.met.blockErrors.Inc()
		r.log.Error("block error", "key", fm.msg.Key, "err", err)
	}
	r.met.processed.Inc()
}

// laneOf hashes a key to a lane index using FNV-1a (allocation-free on the
// hot path, same algorithm as the engine shardOf).
func laneOf(key string, n int) int {
	const (
		offset64 = 14695981039346656037
		prime64  = 1099511628211
	)
	h := uint64(offset64)
	for i := 0; i < len(key); i++ {
		h ^= uint64(key[i]) // per-byte FNV-1a, same as the engine shardOf
		h *= prime64
	}
	return int(h % uint64(n)) // #nosec G115 -- value already reduced mod n
}

// routerMetrics are the router's Prometheus collectors.
type routerMetrics struct {
	submitted   prometheus.Counter
	processed   prometheus.Counter
	droppedQoS0 prometheus.Counter
	droppedQoS1 prometheus.Counter
	blockErrors prometheus.Counter
}

func newRouterMetrics(reg prometheus.Registerer) *routerMetrics {
	m := &routerMetrics{
		submitted: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_flow_router_submitted_total",
			Help: "Messages submitted to the flow router.",
		}),
		processed: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_flow_router_processed_total",
			Help: "Messages successfully processed through the block graph.",
		}),
		droppedQoS0: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_flow_router_dropped_qos0_total",
			Help: "QoS 0 messages dropped due to lane overflow.",
		}),
		droppedQoS1: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_flow_router_dropped_qos1_total",
			Help: "QoS ≥ 1 messages dropped due to lane overflow.",
		}),
		blockErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "astrate_flow_router_block_errors_total",
			Help: "Block processing errors (logged, counted, not fatal).",
		}),
	}
	if reg != nil {
		reg.MustRegister(m.submitted, m.processed, m.droppedQoS0, m.droppedQoS1, m.blockErrors)
	}
	return m
}
