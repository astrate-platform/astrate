package flow

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// FlowStatus enumerates the lifecycle states of a flow.
type FlowStatus uint8

const (
	// FlowStatusCreating indicates the flow is being initialised.
	FlowStatusCreating FlowStatus = iota
	// FlowStatusRunning indicates the flow is accepting and processing messages.
	FlowStatusRunning
	// FlowStatusStopped indicates the flow has been gracefully shut down.
	FlowStatusStopped
	// FlowStatusFailed indicates the flow failed during initialisation.
	FlowStatusFailed
)

// String returns a human-readable label for s.
func (s FlowStatus) String() string {
	switch s {
	case FlowStatusCreating:
		return "creating"
	case FlowStatusRunning:
		return "running"
	case FlowStatusStopped:
		return "stopped"
	case FlowStatusFailed:
		return "failed"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}

var (
	// ErrFlowExists is returned by StartFlow when a flow with the same
	// instance key is already registered.
	ErrFlowExists = errors.New("flow: already running")
	// ErrFlowNotFound is returned when a flow ID does not match any
	// registered flow.
	ErrFlowNotFound = errors.New("flow: not found")
)

// FlowConfig holds the parameters needed to instantiate a running flow.
type FlowConfig struct {
	// PipelineID is the Manager map key for this instance (realm/flowName).
	// Historical field name; value is FlowInstanceID, not the pipeline recipe name.
	PipelineID string
	// Blocks is the ordered list of blocks forming the processing graph.
	Blocks []Block
	// RouterCfg is applied to the underlying Router.
	RouterCfg RouterConfig
	// Registerer receives Prometheus collectors; nil leaves them
	// unregistered.
	Registerer prometheus.Registerer
}

// Flow is a running instance of a pipeline. It owns a Router that processes
// messages through a BlockGraph, a source pump that feeds Source blocks into
// the router, and exposes status information.
type Flow struct {
	id         string
	pipelineID string
	status     FlowStatus
	router     *Router
	graph      *BlockGraph

	// cancelPump stops the source pump; pumpWG waits for it to exit.
	cancelPump context.CancelFunc
	pumpWG     sync.WaitGroup

	mu      sync.RWMutex
	created time.Time
	started time.Time
	stopped time.Time
}

// ID returns the flow's unique identifier.
func (f *Flow) ID() string { return f.id }

// PipelineID returns the pipeline this flow was created from.
func (f *Flow) PipelineID() string { return f.pipelineID }

// Status returns the current lifecycle status.
func (f *Flow) Status() FlowStatus {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.status
}

// Router returns the underlying router for message submission.
func (f *Flow) Router() *Router { return f.router }

// CreatedAt returns when the flow was created.
func (f *Flow) CreatedAt() time.Time {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.created
}

// StartedAt returns when the flow started processing. Zero value if not yet
// started.
func (f *Flow) StartedAt() time.Time {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.started
}

// StoppedAt returns when the flow was stopped. Zero value if still running.
func (f *Flow) StoppedAt() time.Time {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.stopped
}

// setStatus is called under lock by the manager.
func (f *Flow) setStatus(s FlowStatus) { f.status = s }

// Manager manages the lifecycle of running flows. It is safe for concurrent
// use.
type Manager struct {
	mu    sync.RWMutex
	flows map[string]*Flow
	seq   int64
}

// NewManager returns an empty flow manager.
func NewManager() *Manager {
	return &Manager{flows: make(map[string]*Flow)}
}

// StartFlow instantiates a pipeline into a running Flow. The block graph is
// built before the flow is registered so construction failures leave no map
// entry (durable layer records status=failed separately). On success the
// status transitions to running.
func (m *Manager) StartFlow(ctx context.Context, cfg FlowConfig) (*Flow, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.flows[cfg.PipelineID]; exists {
		return nil, fmt.Errorf("%w: %s", ErrFlowExists, cfg.PipelineID)
	}

	// Build graph before insert so a bad config never occupies the instance key.
	graph, err := NewBlockGraph(cfg.Blocks...)
	if err != nil {
		return nil, err
	}

	m.seq++
	f := &Flow{
		id:         fmt.Sprintf("flow-%d", m.seq),
		pipelineID: cfg.PipelineID,
		status:     FlowStatusCreating,
		created:    time.Now(),
		graph:      graph,
		router:     NewRouter(graph, cfg.RouterCfg, cfg.Registerer),
	}
	m.flows[cfg.PipelineID] = f

	f.router.Run(ctx)

	// Pump Source blocks independently of the StartFlow caller's context so
	// a short-lived ctx does not tear the pump down while the flow is running.
	pumpCtx, cancel := context.WithCancel(context.Background())
	f.cancelPump = cancel
	f.startSourcePump(pumpCtx)

	f.mu.Lock()
	f.setStatus(FlowStatusRunning)
	f.started = time.Now()
	f.mu.Unlock()

	return f, nil
}

// UnregisterFlow removes an instance key from the manager map. Call after
// StopFlow when deleting a durable row so the name can be reused.
func (m *Manager) UnregisterFlow(instanceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.flows, instanceID)
}

// sourceIdleBackoff is how long the pump waits after an empty Emit before
// polling again. Blocking Sources (e.g. AstarteSource) never hit this on the
// hot path; it only protects non-blocking SourceFunc adapters from spinning.
const sourceIdleBackoff = 10 * time.Millisecond

// startSourcePump launches one goroutine per Source block. Each goroutine
// blocks in Emit until messages arrive (or the context is cancelled), then
// submits them into the router for the remaining graph stages.
func (f *Flow) startSourcePump(ctx context.Context) {
	for _, src := range f.graph.Sources() {
		src := src
		f.pumpWG.Add(1)
		go func() {
			defer f.pumpWG.Done()
			for {
				msgs, err := src.Emit(ctx)
				if err != nil {
					// Context cancellation is the normal shutdown path.
					if ctx.Err() != nil {
						return
					}
					// Transient produce errors: log via router and keep polling.
					if f.router != nil && f.router.log != nil {
						f.router.log.Error("source emit error",
							"source", src.Name(), "err", err)
					}
					continue
				}
				if ctx.Err() != nil {
					return
				}
				if len(msgs) == 0 {
					select {
					case <-ctx.Done():
						return
					case <-time.After(sourceIdleBackoff):
					}
					continue
				}
				for _, msg := range msgs {
					if msg == nil {
						continue
					}
					// Live device data: QoS 0 matches the stream bus's
					// never-backpressure-ingestion philosophy (§1.4).
					f.router.Submit(msg, 0)
				}
			}
		}()
	}
}

// StopFlow gracefully shuts down the flow identified by pipelineID. It stops
// the source pump, drains in-flight messages, calls Stop on every Stopper
// block, and transitions the status to stopped.
func (m *Manager) StopFlow(ctx context.Context, pipelineID string) error {
	m.mu.RLock()
	f, ok := m.flows[pipelineID]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("%w: %s", ErrFlowNotFound, pipelineID)
	}

	// 1. Stop producing: cancel pump and wait for Emit loops to exit so no
	//    new messages are submitted during drain.
	if f.cancelPump != nil {
		f.cancelPump()
	}
	f.pumpWG.Wait()

	// 2. Drain in-flight lane work.
	if err := f.router.Drain(ctx); err != nil {
		return fmt.Errorf("flow drain %s: %w", pipelineID, err)
	}

	// 3. Release block resources (e.g. AstarteSource bus subscriptions).
	for _, b := range f.graph.Blocks() {
		if s, ok := b.(Stopper); ok {
			s.Stop()
		}
	}

	f.mu.Lock()
	f.setStatus(FlowStatusStopped)
	f.stopped = time.Now()
	f.mu.Unlock()

	return nil
}

// GetFlowStatus returns the current status of the flow identified by
// pipelineID.
func (m *Manager) GetFlowStatus(pipelineID string) (FlowStatus, error) {
	m.mu.RLock()
	f, ok := m.flows[pipelineID]
	m.mu.RUnlock()

	if !ok {
		return 0, fmt.Errorf("%w: %s", ErrFlowNotFound, pipelineID)
	}

	return f.Status(), nil
}

// ListFlows returns a snapshot of all registered flows.
func (m *Manager) ListFlows() []*Flow {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]*Flow, 0, len(m.flows))
	for _, f := range m.flows {
		out = append(out, f)
	}
	return out
}

// Shutdown drains all running flows. It is intended for orderly process exit;
// each flow gets the same context deadline.
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.RLock()
	flows := make([]*Flow, 0, len(m.flows))
	for _, f := range m.flows {
		flows = append(flows, f)
	}
	m.mu.RUnlock()

	var first error
	for _, f := range flows {
		f.mu.RLock()
		s := f.status
		f.mu.RUnlock()
		if s == FlowStatusRunning {
			if err := m.StopFlow(ctx, f.PipelineID()); err != nil && first == nil {
				first = err
			}
		}
	}
	return first
}
