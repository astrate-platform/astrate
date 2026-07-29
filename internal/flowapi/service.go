// Package flowapi is the operator-facing Flow surface: realm-scoped pipeline
// CRUD (store) and start/stop/status for running flows (flow.Manager).
// Paths are Astrate-native (/flow/v1/...) until a real upstream client forces
// wire-compatible routes (milestone v2.0 gap 3).
package flowapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/astrate-platform/astrate/internal/engine/stream"
	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/store"
)

// ErrValidation marks a well-formed request that fails pipeline/block rules.
var ErrValidation = errors.New("flowapi: validation failed")

// Service implements pipeline CRUD and flow lifecycle over store + Manager.
type Service struct {
	st  *store.Store
	mgr *flow.Manager
	reg *flow.Registry
	bus *stream.Bus
	log *slog.Logger
}

// NewService wires store, manager, block registry, and the live bus.
// mgr and reg must be non-nil; bus may be nil only if no pipeline uses
// bus-backed sources.
func NewService(st *store.Store, mgr *flow.Manager, reg *flow.Registry, bus *stream.Bus, log *slog.Logger) *Service {
	if log == nil {
		log = slog.Default()
	}
	return &Service{st: st, mgr: mgr, reg: reg, bus: bus, log: log}
}

// Manager exposes the flow manager (process shutdown).
func (s *Service) Manager() *flow.Manager { return s.mgr }

// PipelineView is the operator JSON shape for one stored pipeline.
type PipelineView struct {
	Name       string          `json:"name"`
	Definition json.RawMessage `json:"definition"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// FlowView is the operator JSON shape for one running (or registered) flow.
type FlowView struct {
	ID         string `json:"id"`
	Pipeline   string `json:"pipeline"`
	Realm      string `json:"realm"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at,omitempty"`
	StartedAt  string `json:"started_at,omitempty"`
	StoppedAt  string `json:"stopped_at,omitempty"`
}

func (s *Service) realmID(ctx context.Context, realm string) (int16, error) {
	r, err := s.st.GetRealmByName(ctx, realm)
	if err != nil {
		return 0, err
	}
	return r.ID, nil
}

// ListPipelines returns every pipeline name for a realm.
func (s *Service) ListPipelines(ctx context.Context, realm string) ([]string, error) {
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	list, err := s.st.ListPipelines(ctx, id)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(list))
	for i := range list {
		names[i] = list[i].Name
	}
	return names, nil
}

// GetPipeline returns one pipeline by name.
func (s *Service) GetPipeline(ctx context.Context, realm, name string) (*PipelineView, error) {
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	p, err := s.st.GetPipeline(ctx, id, name)
	if err != nil {
		return nil, err
	}
	return toPipelineView(p), nil
}

// CreatePipeline validates and stores a pipeline definition.
func (s *Service) CreatePipeline(ctx context.Context, realm, name string, definition []byte) (*PipelineView, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: pipeline name is required", ErrValidation)
	}
	if _, err := flow.ParseDefinition(name, name, definition); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	// Ensure every block_type is registered before persisting.
	if err := s.checkBlockTypes(definition); err != nil {
		return nil, err
	}
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	p, err := s.st.CreatePipeline(ctx, id, name, definition)
	if err != nil {
		return nil, err
	}
	return toPipelineView(p), nil
}

// UpdatePipeline replaces a pipeline definition.
func (s *Service) UpdatePipeline(ctx context.Context, realm, name string, definition []byte) (*PipelineView, error) {
	if _, err := flow.ParseDefinition(name, name, definition); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	if err := s.checkBlockTypes(definition); err != nil {
		return nil, err
	}
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	p, err := s.st.UpdatePipeline(ctx, id, name, definition)
	if err != nil {
		return nil, err
	}
	return toPipelineView(p), nil
}

// DeletePipeline removes a stored pipeline. Running flows are not stopped.
func (s *Service) DeletePipeline(ctx context.Context, realm, name string) error {
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return err
	}
	return s.st.DeletePipeline(ctx, id, name)
}

func (s *Service) checkBlockTypes(definition []byte) error {
	var shape struct {
		Blocks []struct {
			Name      string `json:"name"`
			BlockType string `json:"block_type"`
		} `json:"blocks"`
	}
	if err := json.Unmarshal(definition, &shape); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	for _, b := range shape.Blocks {
		if b.BlockType == "" {
			return fmt.Errorf("%w: block %q has empty block_type", ErrValidation, b.Name)
		}
		if !s.reg.Has(b.BlockType) {
			return fmt.Errorf("%w: unknown block_type %q on block %q (known: %v)",
				ErrValidation, b.BlockType, b.Name, s.reg.Types())
		}
	}
	return nil
}

// StartFlow loads a stored pipeline, instantiates blocks, and starts it.
func (s *Service) StartFlow(ctx context.Context, realm, pipelineName string) (*FlowView, error) {
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	stored, err := s.st.GetPipeline(ctx, id, pipelineName)
	if err != nil {
		return nil, err
	}
	p, err := flow.ParseDefinition(pipelineName, pipelineName, stored.Definition)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	blks, err := s.reg.Instantiate(p, flow.Deps{Bus: s.bus, Realm: realm})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	pipeID := flow.FlowPipelineID(realm, pipelineName)
	f, err := s.mgr.StartFlow(ctx, flow.FlowConfig{
		PipelineID: pipeID,
		Blocks:     blks,
	})
	if err != nil {
		return nil, err
	}
	s.log.Info("flow started", "realm", realm, "pipeline", pipelineName, "flow_id", f.ID())
	return toFlowView(f, realm, pipelineName), nil
}

// StopFlow stops a running flow for the realm/pipeline pair.
func (s *Service) StopFlow(ctx context.Context, realm, pipelineName string) error {
	pipeID := flow.FlowPipelineID(realm, pipelineName)
	if err := s.mgr.StopFlow(ctx, pipeID); err != nil {
		return err
	}
	s.log.Info("flow stopped", "realm", realm, "pipeline", pipelineName)
	return nil
}

// GetFlow returns status for one realm-scoped pipeline instance.
func (s *Service) GetFlow(realm, pipelineName string) (*FlowView, error) {
	pipeID := flow.FlowPipelineID(realm, pipelineName)
	for _, f := range s.mgr.ListFlows() {
		if f.PipelineID() == pipeID {
			return toFlowView(f, realm, pipelineName), nil
		}
	}
	return nil, fmt.Errorf("%w: flow %q", store.ErrNotFound, pipelineName)
}

// ListFlows returns running flows for a realm (pipeline name prefix match).
func (s *Service) ListFlows(realm string) []FlowView {
	prefix := realm + "/"
	out := make([]FlowView, 0)
	for _, f := range s.mgr.ListFlows() {
		pid := f.PipelineID()
		if len(pid) > len(prefix) && pid[:len(prefix)] == prefix {
			name := pid[len(prefix):]
			out = append(out, *toFlowView(f, realm, name))
		}
	}
	return out
}

func toPipelineView(p *store.Pipeline) *PipelineView {
	return &PipelineView{
		Name:       p.Name,
		Definition: json.RawMessage(p.Definition),
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

func toFlowView(f *flow.Flow, realm, pipelineName string) *FlowView {
	v := &FlowView{
		ID:       f.ID(),
		Pipeline: pipelineName,
		Realm:    realm,
		Status:   f.Status().String(),
	}
	if t := f.CreatedAt(); !t.IsZero() {
		v.CreatedAt = t.UTC().Format(time.RFC3339Nano)
	}
	if t := f.StartedAt(); !t.IsZero() {
		v.StartedAt = t.UTC().Format(time.RFC3339Nano)
	}
	if t := f.StoppedAt(); !t.IsZero() {
		v.StoppedAt = t.UTC().Format(time.RFC3339Nano)
	}
	return v
}
