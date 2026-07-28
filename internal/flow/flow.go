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
	// pipeline ID is already registered.
	ErrFlowExists = errors.New("flow: pipeline already running")
	// ErrFlowNotFound is returned when a flow ID does not match any
	// registered flow.
	ErrFlowNotFound = errors.New("flow: not found")
)

// FlowConfig holds the parameters needed to instantiate a running flow.
type FlowConfig struct {
	// PipelineID is the unique identifier for the pipeline this flow runs.
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
// messages through a BlockGraph and exposes status information.
type Flow struct {
	id         string
	pipelineID string
	status     FlowStatus
	router     *Router
	graph      *BlockGraph

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

// StartFlow instantiates a pipeline into a running Flow. The flow is
// created with status creating, the block graph and router are built, and
// on success the status transitions to running. If graph construction fails
// the status is set to failed and both the flow and the error are returned.
func (m *Manager) StartFlow(ctx context.Context, cfg FlowConfig) (*Flow, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.flows[cfg.PipelineID]; exists {
		return nil, fmt.Errorf("%w: %s", ErrFlowExists, cfg.PipelineID)
	}

	m.seq++
	f := &Flow{
		id:         fmt.Sprintf("flow-%d", m.seq),
		pipelineID: cfg.PipelineID,
		status:     FlowStatusCreating,
		created:    time.Now(),
	}
	m.flows[cfg.PipelineID] = f

	graph, err := NewBlockGraph(cfg.Blocks...)
	if err != nil {
		f.setStatus(FlowStatusFailed)
		return f, err
	}

	f.graph = graph
	f.router = NewRouter(graph, cfg.RouterCfg, cfg.Registerer)
	f.router.Run(ctx)

	f.mu.Lock()
	f.setStatus(FlowStatusRunning)
	f.started = time.Now()
	f.mu.Unlock()

	return f, nil
}

// StopFlow gracefully shuts down the flow identified by pipelineID. It
// drains in-flight messages, releases resources, and transitions the status
// to stopped.
func (m *Manager) StopFlow(ctx context.Context, pipelineID string) error {
	m.mu.RLock()
	f, ok := m.flows[pipelineID]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("%w: %s", ErrFlowNotFound, pipelineID)
	}

	if err := f.router.Drain(ctx); err != nil {
		return fmt.Errorf("flow drain %s: %w", pipelineID, err)
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
