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
	"sync"
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
	st     *store.Store
	mgr    *flow.Manager
	reg    *flow.Registry
	bus    *stream.Bus
	ingest flow.IngestFunc
	log    *slog.Logger

	// restartMu/restarting dedupe auto-restart loops when several blocks of
	// the same flow die (or fatal races a manual reload).
	restartMu  sync.Mutex
	restarting map[string]bool
}

// NewService wires store, manager, block registry, and the live bus.
// mgr and reg must be non-nil; bus may be nil only if no pipeline uses
// bus-backed sources; ingest may be nil only if no pipeline uses
// virtual_device_pool (#84).
func NewService(st *store.Store, mgr *flow.Manager, reg *flow.Registry, bus *stream.Bus, ingest flow.IngestFunc, log *slog.Logger) *Service {
	if log == nil {
		log = slog.Default()
	}
	return &Service{st: st, mgr: mgr, reg: reg, bus: bus, ingest: ingest, log: log}
}

// Manager exposes the flow manager (process shutdown).
func (s *Service) Manager() *flow.Manager { return s.mgr }

// PipelineView is the operator JSON shape for one stored pipeline.
type PipelineView struct {
	Name       string          `json:"name"`
	Definition json.RawMessage `json:"definition"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	// ReferencingFlows is set on update: names of durable flows built from
	// this pipeline. Running instances keep executing the old definition
	// until explicitly reloaded (POST /flows/{name}/reload).
	ReferencingFlows []string `json:"referencing_flows,omitempty"`
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
	FailedBlock  *string         `json:"failed_block,omitempty"`
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

// UpdatePipeline replaces a pipeline definition. Running flows built from it
// keep executing the old definition; the response lists them so the edit is
// never a silent no-op (issue #44).
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
	view := toPipelineView(p)
	if view.ReferencingFlows, err = s.referencingFlows(ctx, id, name); err != nil {
		return nil, err
	}
	return view, nil
}

// referencingFlows lists durable flows in a realm whose pipeline_name matches.
func (s *Service) referencingFlows(ctx context.Context, realmID int16, pipelineName string) ([]string, error) {
	rows, err := s.st.ListFlows(ctx, realmID)
	if err != nil {
		return nil, err
	}
	var out []string
	for i := range rows {
		if rows[i].PipelineName == pipelineName {
			out = append(out, rows[i].Name)
		}
	}
	return out, nil
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
			_ = s.st.UpdateFlowRuntime(ctx, row.RealmID, row.Name, "failed", &msg, nil, nil, nil)
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

// resolveAndBuild is the resolve half of every start path: pipeline lookup,
// ${config.*} substitution, validation, and block instantiation. Every error
// wraps ErrValidation. It has no durable-row side effects; callers own status
// transitions.
func (s *Service) resolveAndBuild(ctx context.Context, realmID int16, realm, name, pipelineName string, config map[string]any) ([]flow.Block, error) {
	stored, err := s.st.GetPipeline(ctx, realmID, pipelineName)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			err = fmt.Errorf("%w: pipeline %q not found", ErrValidation, pipelineName)
		}
		return nil, err
	}
	def, err := flow.SubstituteConfig(stored.Definition, config)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	if err := s.checkBlockTypes(def); err != nil {
		return nil, err
	}
	p, err := flow.ParseDefinition(pipelineName, pipelineName, def)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	blks, err := s.reg.Instantiate(p, s.flowDeps(realmID, realm, name))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}
	return blks, nil
}

// flowDeps builds the per-instance block dependencies: live bus, realm and
// flow identity, the fatal-block callback, and the engine ingest path (#84).
func (s *Service) flowDeps(realmID int16, realm, name string) flow.Deps {
	return flow.Deps{
		Bus:         s.bus,
		Realm:       realm,
		FlowName:    name,
		NotifyFatal: s.onBlockFatal(realmID, realm, name),
		Ingest:      s.ingest,
	}
}

// startFlowInstance is the single start path for POST create and boot rehydrate.
// The durable row must already exist (CreateAndStartFlow inserts first).
func (s *Service) startFlowInstance(ctx context.Context, realm string, realmID int16, name, pipelineName string, config map[string]any) (*FlowView, error) {
	markFailed := func(err error) {
		msg := err.Error()
		_ = s.st.UpdateFlowRuntime(ctx, realmID, name, "failed", &msg, nil, nil, nil)
	}

	blks, err := s.resolveAndBuild(ctx, realmID, realm, name, pipelineName, config)
	if err != nil {
		markFailed(err)
		return nil, err
	}

	instanceID := flow.InstanceID(realm, name)
	// Creating status while building live graph (optional observation).
	_ = s.st.UpdateFlowRuntime(ctx, realmID, name, "creating", nil, nil, nil, nil)

	f, err := s.mgr.StartFlow(ctx, flow.Config{
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
	if err := s.st.UpdateFlowRuntime(ctx, realmID, name, "running", nil, nil, &now, nil); err != nil {
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
	instanceID := flow.InstanceID(realm, name)
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

// RestartFlowInstance tears down the live graph (if any) and rebuilds it from
// the durable row's current pipeline + config. The durable row is never
// deleted; on rebuild failure the flow is marked failed with the error. Works
// for running flows (reload semantics) and stopped/failed ones (manual start).
func (s *Service) RestartFlowInstance(ctx context.Context, realm string, realmID int16, name string) (*FlowView, error) {
	row, err := s.st.GetFlow(ctx, realmID, name)
	if err != nil {
		return nil, err
	}
	instanceID := flow.InstanceID(realm, name)
	if err := s.mgr.StopFlow(ctx, instanceID); err != nil && !errors.Is(err, flow.ErrFlowNotFound) {
		return nil, err
	}
	s.mgr.UnregisterFlow(instanceID)

	var cfgMap map[string]any
	if err := json.Unmarshal(row.Config, &cfgMap); err != nil {
		msg := fmt.Sprintf("invalid stored config: %v", err)
		_ = s.st.UpdateFlowRuntime(ctx, realmID, name, "failed", &msg, nil, nil, nil)
		return nil, fmt.Errorf("%w: %s", ErrValidation, msg)
	}
	if cfgMap == nil {
		cfgMap = map[string]any{}
	}

	blks, err := s.resolveAndBuild(ctx, realmID, realm, name, row.PipelineName, cfgMap)
	if err != nil {
		msg := err.Error()
		_ = s.st.UpdateFlowRuntime(ctx, realmID, name, "failed", &msg, nil, nil, nil)
		return nil, err
	}

	_ = s.st.UpdateFlowRuntime(ctx, realmID, name, "creating", nil, nil, nil, nil)
	f, err := s.mgr.StartFlow(ctx, flow.Config{
		PipelineID: instanceID,
		Blocks:     blks,
	})
	if err != nil {
		if !errors.Is(err, flow.ErrFlowExists) {
			msg := err.Error()
			_ = s.st.UpdateFlowRuntime(ctx, realmID, name, "failed", &msg, nil, nil, nil)
		}
		return nil, err
	}

	now := time.Now().UTC()
	if err := s.st.UpdateFlowRuntime(ctx, realmID, name, "running", nil, nil, &now, nil); err != nil {
		s.log.Error("flow restarted but failed to persist runtime", "realm", realm, "name", name, "error", err)
	}
	s.log.Info("flow restarted", "realm", realm, "name", name, "pipeline", row.PipelineName, "runtime_id", f.ID())
	return s.GetFlow(ctx, realm, name)
}

// restartBackoffBase/max shape the auto-restart delay ladder for flows whose
// blocks die at runtime (issue #45): 1s, 2s, 4s, ... capped at 30s, retrying
// until the flow starts or its durable row is deleted / auto_restart disabled.
const (
	restartBackoffBase = time.Second
	restartBackoffMax  = 30 * time.Second
)

// onBlockFatal returns the Deps.NotifyFatal callback for one flow instance.
// A dead block fails the whole flow (records failed_block), tears the live
// graph down, and — for auto_restart flows — schedules a backoff rebuild.
func (s *Service) onBlockFatal(realmID int16, realm, name string) func(string, error) {
	instanceID := flow.InstanceID(realm, name)
	return func(block string, cause error) {
		bg := context.Background()
		s.log.Error("flow block died", "realm", realm, "flow", name, "block", block, "error", cause)
		msg := fmt.Sprintf("block %q failed: %v", block, cause)
		if err := s.st.UpdateFlowRuntime(bg, realmID, name, "failed", &msg, &block, nil, nil); err != nil {
			s.log.Error("failed to persist block death", "realm", realm, "flow", name, "error", err)
		}

		if err := s.mgr.StopFlow(bg, instanceID); err != nil && !errors.Is(err, flow.ErrFlowNotFound) {
			s.log.Error("failed to stop flow after block death", "realm", realm, "flow", name, "error", err)
		}
		s.mgr.UnregisterFlow(instanceID)

		row, err := s.st.GetFlow(bg, realmID, name)
		if err != nil {
			if !errors.Is(err, store.ErrNotFound) {
				s.log.Error("failed to read flow after block death", "realm", realm, "flow", name, "error", err)
			}
			return
		}
		if !row.AutoRestart {
			return
		}
		if !s.beginRestart(instanceID) {
			return // another fatal/restart loop is already rebuilding
		}
		go s.restartWithBackoff(realm, realmID, name, instanceID)
	}
}

// beginRestart claims the instance for an auto-restart loop; false means one
// is already active. Use endRestart to release.
func (s *Service) beginRestart(instanceID string) bool {
	s.restartMu.Lock()
	defer s.restartMu.Unlock()
	if s.restarting == nil {
		s.restarting = make(map[string]bool)
	}
	if s.restarting[instanceID] {
		return false
	}
	s.restarting[instanceID] = true
	return true
}

func (s *Service) endRestart(instanceID string) {
	s.restartMu.Lock()
	defer s.restartMu.Unlock()
	delete(s.restarting, instanceID)
}

// restartWithBackoff keeps re-running RestartFlowInstance until it succeeds,
// the durable row disappears (flow deleted), or auto_restart is turned off.
func (s *Service) restartWithBackoff(realm string, realmID int16, name, instanceID string) {
	defer s.endRestart(instanceID)
	bg := context.Background()
	delay := restartBackoffBase
	for attempt := 0; ; attempt++ {
		time.Sleep(delay)
		delay *= 2
		if delay > restartBackoffMax {
			delay = restartBackoffMax
		}

		row, err := s.st.GetFlow(bg, realmID, name)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return // deleted while we waited; nothing to bring back
			}
			s.log.Error("auto-restart aborted: cannot read flow row", "realm", realm, "flow", name, "error", err)
			return
		}
		if !row.AutoRestart {
			return
		}

		ctx, cancel := context.WithTimeout(bg, 2*time.Minute)
		_, err = s.RestartFlowInstance(ctx, realm, realmID, name)
		cancel()
		if err == nil {
			s.log.Info("flow auto-restarted after block death",
				"realm", realm, "flow", name, "attempt", attempt+1)
			return
		}
		s.log.Warn("flow auto-restart attempt failed",
			"realm", realm, "flow", name, "attempt", attempt+1, "next_in", delay.String(), "error", err)
	}
}

// ReloadFlow re-resolves a flow's pipeline by name and rebuilds the live
// graph from it (issue #44). Explicit: editing a pipeline never touches
// referencing flows on its own.
func (s *Service) ReloadFlow(ctx context.Context, realm, name string) (*FlowView, error) {
	id, err := s.realmID(ctx, realm)
	if err != nil {
		return nil, err
	}
	return s.RestartFlowInstance(ctx, realm, id, name)
}

// UpdateFlowConfig replaces a flow's config snapshot (issue #46). The new
// config must substitute cleanly against the flow's current pipeline; if the
// flow is live it is rebuilt immediately with the new config.
func (s *Service) UpdateFlowConfig(ctx context.Context, realm, name string, config json.RawMessage) (*FlowView, error) {
	if len(config) == 0 {
		config = json.RawMessage(`{}`)
	}
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
	row, err := s.st.GetFlow(ctx, id, name)
	if err != nil {
		return nil, err
	}
	// Dry run: the new config must satisfy the pipeline's ${config.*} keys.
	if err := s.resolveAndBuildDryRun(ctx, id, realm, name, row.PipelineName, cfgMap); err != nil {
		return nil, err
	}
	if _, err := s.st.UpdateFlowConfig(ctx, id, name, config); err != nil {
		return nil, err
	}
	s.log.Info("flow config updated", "realm", realm, "name", name)

	if s.liveFlow(realm, name) != nil {
		return s.RestartFlowInstance(ctx, realm, id, name)
	}
	row, err = s.st.GetFlow(ctx, id, name)
	if err != nil {
		return nil, err
	}
	return mergeFlowView(realm, row, nil), nil
}

// resolveAndBuildDryRun validates config substitution + definition shape
// without instantiating blocks (no container starts, no bus wiring).
func (s *Service) resolveAndBuildDryRun(ctx context.Context, realmID int16, _, _, pipelineName string, config map[string]any) error {
	stored, err := s.st.GetPipeline(ctx, realmID, pipelineName)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			err = fmt.Errorf("%w: pipeline %q not found", ErrValidation, pipelineName)
		}
		return err
	}
	def, err := flow.SubstituteConfig(stored.Definition, config)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
	if err := s.checkBlockTypes(def); err != nil {
		return err
	}
	_, err = flow.ParseDefinition(pipelineName, pipelineName, def)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}
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
		_ = s.st.UpdateFlowRuntime(ctx, id, name, "stopped", nil, nil, nil, &now)
	}
}

func (s *Service) liveFlow(realm, name string) *flow.Flow {
	instanceID := flow.InstanceID(realm, name)
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
		FailedBlock:  row.FailedBlock,
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
