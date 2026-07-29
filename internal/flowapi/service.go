// Package flowapi is the operator-facing Flow surface: realm-scoped pipeline
// CRUD (store) and named durable flow lifecycle (store + flow.Manager).
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
	"github.com/astrate-platform/astrate/internal/flow/blocks"
	"github.com/astrate-platform/astrate/internal/store"
)

// ErrValidation marks a well-formed request that fails pipeline/block rules.
var ErrValidation = errors.New("flowapi: validation failed")

// Service implements pipeline CRUD and durable flow lifecycle over store + Manager.
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

// FlowView is the operator JSON shape for one durable flow (merged with live
// status when the instance is running in the Manager).
type FlowView struct {
	Name         string          `json:"name"`
	Pipeline     string          `json:"pipeline"`
	Realm        string          `json:"realm"`
	Config       json.RawMessage `json:"config"`
	AutoRestart  bool            `json:"auto_restart"`
	Status       string          `json:"status"`
	ErrorMessage *string         `json:"error_message"`
	RuntimeID    string          `json:"runtime_id,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	StartedAt    *time.Time      `json:"started_at"`
	StoppedAt    *time.Time      `json:"stopped_at"`
}

// CreateFlowRequest is the POST /flows body after defaults are applied.
type CreateFlowRequest struct {
	Name        string
	Pipeline    string
	Config      json.RawMessage
	AutoRestart bool
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

// CreateAndStartFlow inserts a durable flow row and starts it (single start path).
func (s *Service) CreateAndStartFlow(ctx context.Context, realm string, req CreateFlowRequest) (*FlowView, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("%w: flow name is required", ErrValidation)
	}
	if req.Pipeline == "" {
		return nil, fmt.Errorf("%w: pipeline is required", ErrValidation)
	}
	config := req.Config
	if len(config) == 0 {
		config = json.RawMessage(`{}`)
	}
	// Validate config is a JSON object before insert.
	var cfgMap map[string]any
	if err := json.Unmarshal(config, &cfgMap); err != nil {
		return nil, fmt.Errorf("%w: config must be a JSON object: %v", ErrValidation, err)
	}
	if cfgMap == nil {
		cfgMap = map[string]any{}
	}

	id, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	row, err := s.st.CreateFlow(ctx, id, req.Name, req.Pipeline, []byte(config), req.AutoRestart)
	if err != nil {
		return nil, err
	}
	return s.startFlowInstance(ctx, realm, id, row.Name, row.PipelineName, cfgMap)
}

// RehydrateAutoRestart starts every durable flow with auto_restart=true.
// Individual failures are logged and recorded as status=failed; boot continues.
// Returns an error only if listing rows fails.
func (s *Service) RehydrateAutoRestart(ctx context.Context) error {
	rows, err := s.st.ListAutoRestartFlows(ctx)
	if err != nil {
		return err
	}
	for _, row := range rows {
		var cfgMap map[string]any
		if err := json.Unmarshal(row.Config, &cfgMap); err != nil {
			msg := fmt.Sprintf("invalid stored config: %v", err)
			s.log.Error("flow rehydrate failed", "realm", row.RealmName, "name", row.Name, "error", msg)
			_ = s.st.UpdateFlowRuntime(ctx, row.RealmID, row.Name, "failed", &msg, nil, nil)
			continue
		}
		if cfgMap == nil {
			cfgMap = map[string]any{}
		}
		if _, err := s.startFlowInstance(ctx, row.RealmName, row.RealmID, row.Name, row.PipelineName, cfgMap); err != nil {
			s.log.Error("flow rehydrate failed", "realm", row.RealmName, "name", row.Name, "error", err)
			// startFlowInstance already marked failed when row exists
			continue
		}
		s.log.Info("flow rehydrated", "realm", row.RealmName, "name", row.Name, "pipeline", row.PipelineName)
	}
	return nil
}

// startFlowInstance is the single start path for POST create and boot rehydrate.
// The durable row must already exist (CreateAndStartFlow inserts first).
func (s *Service) startFlowInstance(ctx context.Context, realm string, realmID int16, name, pipelineName string, config map[string]any) (*FlowView, error) {
	markFailed := func(err error) {
		msg := err.Error()
		_ = s.st.UpdateFlowRuntime(ctx, realmID, name, "failed", &msg, nil, nil)
	}

	stored, err := s.st.GetPipeline(ctx, realmID, pipelineName)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			err = fmt.Errorf("%w: pipeline %q not found", ErrValidation, pipelineName)
		}
		markFailed(err)
		return nil, err
	}

	def, err := flow.SubstituteConfig(stored.Definition, config)
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrValidation, err)
		markFailed(err)
		return nil, err
	}
	if err := s.checkBlockTypes(def); err != nil {
		markFailed(err)
		return nil, err
	}
	p, err := flow.ParseDefinition(pipelineName, pipelineName, def)
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrValidation, err)
		markFailed(err)
		return nil, err
	}
	blks, err := s.reg.Instantiate(p, flow.Deps{Bus: s.bus, Realm: realm, FlowName: name})
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrValidation, err)
		markFailed(err)
		return nil, err
	}

	instanceID := flow.FlowInstanceID(realm, name)
	// Creating status while building live graph (optional observation).
	_ = s.st.UpdateFlowRuntime(ctx, realmID, name, "creating", nil, nil, nil)

	f, err := s.mgr.StartFlow(ctx, flow.FlowConfig{
		PipelineID: instanceID,
		Blocks:     blks,
	})
	if err != nil {
		// Already live (e.g. double rehydrate): do not overwrite status with failed.
		if !errors.Is(err, flow.ErrFlowExists) {
			markFailed(err)
		}
		return nil, err
	}

	now := time.Now().UTC()
	if err := s.st.UpdateFlowRuntime(ctx, realmID, name, "running", nil, &now, nil); err != nil {
		s.log.Error("flow started but failed to persist runtime", "realm", realm, "name", name, "error", err)
	}
	s.log.Info("flow started", "realm", realm, "name", name, "pipeline", pipelineName, "runtime_id", f.ID())

	row, err := s.st.GetFlow(ctx, realmID, name)
	if err != nil {
		// Live is up; synthesize view from memory + known fields.
		return mergeFlowView(realm, &store.Flow{
			Name: name, PipelineName: pipelineName, Config: mustJSON(config),
			AutoRestart: true, Status: "running", StartedAt: &now,
			CreatedAt: now, UpdatedAt: now,
		}, f), nil
	}
	return mergeFlowView(realm, row, f), nil
}

// DeleteFlow stops a live instance if present, unregisters it, and deletes
// the durable row.
func (s *Service) DeleteFlow(ctx context.Context, realm, name string) error {
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return err
	}
	instanceID := flow.FlowInstanceID(realm, name)
	if err := s.mgr.StopFlow(ctx, instanceID); err != nil && !errors.Is(err, flow.ErrFlowNotFound) {
		return err
	}
	s.mgr.UnregisterFlow(instanceID)

	if err := s.st.DeleteFlow(ctx, id, name); err != nil {
		return err
	}
	s.log.Info("flow deleted", "realm", realm, "name", name)
	return nil
}

// GetFlow returns one durable flow by name, merged with live status.
func (s *Service) GetFlow(ctx context.Context, realm, name string) (*FlowView, error) {
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	row, err := s.st.GetFlow(ctx, id, name)
	if err != nil {
		return nil, err
	}
	return mergeFlowView(realm, row, s.liveFlow(realm, name)), nil
}

// ListFlows returns durable flows for a realm (including stopped/failed),
// merged with live Manager status when running.
func (s *Service) ListFlows(ctx context.Context, realm string) ([]FlowView, error) {
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	rows, err := s.st.ListFlows(ctx, id)
	if err != nil {
		return nil, err
	}
	out := make([]FlowView, 0, len(rows))
	for i := range rows {
		out = append(out, *mergeFlowView(realm, &rows[i], s.liveFlow(realm, rows[i].Name)))
	}
	return out, nil
}

// MarkRunningFlowsStopped best-effort sets durable status to stopped for
// instances that were running (clean process shutdown).
func (s *Service) MarkRunningFlowsStopped(ctx context.Context) {
	now := time.Now().UTC()
	for _, f := range s.mgr.ListFlows() {
		if f.Status() != flow.FlowStatusRunning && f.Status() != flow.FlowStatusStopped {
			continue
		}
		realm, name, ok := splitInstanceID(f.PipelineID())
		if !ok {
			continue
		}
		id, err := s.realmID(ctx, realm)
		if err != nil {
			continue
		}
		_ = s.st.UpdateFlowRuntime(ctx, id, name, "stopped", nil, nil, &now)
	}
}

func (s *Service) liveFlow(realm, name string) *flow.Flow {
	instanceID := flow.FlowInstanceID(realm, name)
	for _, f := range s.mgr.ListFlows() {
		if f.PipelineID() == instanceID {
			return f
		}
	}
	return nil
}

// ListBlocks returns operator docs for every registered block type.
// Realm is accepted for path symmetry with other Flow routes; the catalog is process-global.
func (s *Service) ListBlocks(_ string) []blocks.Info {
	return blocks.InfoForTypes(s.reg.Types())
}

// GetBlock returns operator docs for one registered type, or ErrNotFound.
func (s *Service) GetBlock(_, blockType string) (*blocks.Info, error) {
	if !s.reg.Has(blockType) {
		return nil, fmt.Errorf("%w: block type %q", store.ErrNotFound, blockType)
	}
	if info, ok := blocks.LookupInfo(blockType); ok {
		return &info, nil
	}
	stub := blocks.Info{Type: blockType, Role: blocks.RoleTransform, Summary: "registered block (no built-in docs)"}
	return &stub, nil
}

func toPipelineView(p *store.Pipeline) *PipelineView {
	return &PipelineView{
		Name:       p.Name,
		Definition: json.RawMessage(p.Definition),
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

func mergeFlowView(realm string, row *store.Flow, live *flow.Flow) *FlowView {
	cfg := json.RawMessage(row.Config)
	if len(cfg) == 0 {
		cfg = json.RawMessage(`{}`)
	}
	v := &FlowView{
		Name:         row.Name,
		Pipeline:     row.PipelineName,
		Realm:        realm,
		Config:       cfg,
		AutoRestart:  row.AutoRestart,
		Status:       row.Status,
		ErrorMessage: row.ErrorMessage,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
		StartedAt:    row.StartedAt,
		StoppedAt:    row.StoppedAt,
	}
	if live != nil {
		v.RuntimeID = live.ID()
		v.Status = live.Status().String()
		if t := live.StartedAt(); !t.IsZero() {
			st := t.UTC()
			v.StartedAt = &st
		}
		if t := live.StoppedAt(); !t.IsZero() {
			st := t.UTC()
			v.StoppedAt = &st
		}
	}
	return v
}

func mustJSON(m map[string]any) []byte {
	b, err := json.Marshal(m)
	if err != nil {
		return []byte("{}")
	}
	return b
}

func splitInstanceID(id string) (realm, name string, ok bool) {
	for i := 0; i < len(id); i++ {
		if id[i] == '/' {
			if i == 0 || i == len(id)-1 {
				return "", "", false
			}
			return id[:i], id[i+1:], true
		}
	}
	return "", "", false
}
