package flowapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/flow"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

// maxBodyBytes caps Flow API request bodies (pipeline DAGs are small JSON).
const maxBodyBytes int64 = 1 << 20

// API is the /flow/v1 HTTP surface. Routes use a realm JWT with a_rma
// (operator control, same claim as Realm Management).
type API struct {
	svc     *Service
	require func(http.Handler) http.Handler
}

// NewAPI wires the service to its HTTP surface.
func NewAPI(svc *Service, mw *auth.Middleware) *API {
	return &API{svc: svc, require: mw.RequireRealm(auth.ClaimRealmManagement)}
}

// Mount registers the routes on mux.
func (a *API) Mount(mux *http.ServeMux) {
	h := func(f http.HandlerFunc) http.Handler { return a.require(f) }
	mux.Handle("GET /flow/v1/{realm}/pipelines", h(a.listPipelines))
	mux.Handle("POST /flow/v1/{realm}/pipelines", h(a.createPipeline))
	mux.Handle("GET /flow/v1/{realm}/pipelines/{name}", h(a.getPipeline))
	mux.Handle("PUT /flow/v1/{realm}/pipelines/{name}", h(a.updatePipeline))
	mux.Handle("DELETE /flow/v1/{realm}/pipelines/{name}", h(a.deletePipeline))

	mux.Handle("GET /flow/v1/{realm}/flows", h(a.listFlows))
	mux.Handle("POST /flow/v1/{realm}/flows", h(a.startFlow))
	mux.Handle("GET /flow/v1/{realm}/flows/{name}", h(a.getFlow))
	mux.Handle("DELETE /flow/v1/{realm}/flows/{name}", h(a.stopFlow))
	mux.Handle("POST /flow/v1/{realm}/flows/{name}/reload", h(a.reloadFlow))
	mux.Handle("PUT /flow/v1/{realm}/flows/{name}/config", h(a.updateFlowConfig))

	mux.Handle("GET /flow/v1/{realm}/blocks", h(a.listBlocks))
	mux.Handle("GET /flow/v1/{realm}/blocks/{type}", h(a.getBlock))
}

type createPipelineBody struct {
	Name       string          `json:"name"`
	Definition json.RawMessage `json:"definition"`
}

// startFlowBody is POST /flows: named multi-instance + optional config.
type startFlowBody struct {
	Name        string          `json:"name"`
	Pipeline    string          `json:"pipeline"`
	Config      json.RawMessage `json:"config"`
	AutoRestart *bool           `json:"auto_restart"`
}

func (a *API) listPipelines(w http.ResponseWriter, r *http.Request) {
	names, err := a.svc.ListPipelines(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, names)
}

func (a *API) createPipeline(w http.ResponseWriter, r *http.Request) {
	var body createPipelineBody
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	view, err := a.svc.CreatePipeline(r.Context(), r.PathValue("realm"), body.Name, body.Definition)
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusCreated, view)
}

func (a *API) getPipeline(w http.ResponseWriter, r *http.Request) {
	view, err := a.svc.GetPipeline(r.Context(), r.PathValue("realm"), r.PathValue("name"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, view)
}

func (a *API) updatePipeline(w http.ResponseWriter, r *http.Request) {
	// data is the pipeline definition object (blocks + connections), same shape
	// as create's "definition" field.
	var def json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &def); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	view, err := a.svc.UpdatePipeline(r.Context(), r.PathValue("realm"), r.PathValue("name"), def)
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, view)
}

func (a *API) deletePipeline(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.DeletePipeline(r.Context(), r.PathValue("realm"), r.PathValue("name")); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) listFlows(w http.ResponseWriter, r *http.Request) {
	list, err := a.svc.ListFlows(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, list)
}

func (a *API) startFlow(w http.ResponseWriter, r *http.Request) {
	var body startFlowBody
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if body.Name == "" || body.Pipeline == "" {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	autoRestart := true
	if body.AutoRestart != nil {
		autoRestart = *body.AutoRestart
	}
	view, err := a.svc.CreateAndStartFlow(r.Context(), r.PathValue("realm"), CreateFlowRequest{
		Name:        body.Name,
		Pipeline:    body.Pipeline,
		Config:      body.Config,
		AutoRestart: autoRestart,
	})
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusCreated, view)
}

func (a *API) getFlow(w http.ResponseWriter, r *http.Request) {
	view, err := a.svc.GetFlow(r.Context(), r.PathValue("realm"), r.PathValue("name"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, view)
}

func (a *API) stopFlow(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.DeleteFlow(r.Context(), r.PathValue("realm"), r.PathValue("name")); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// reloadFlow re-resolves the pipeline + substitutes stored config into a new
// live graph (issue #44).
func (a *API) reloadFlow(w http.ResponseWriter, r *http.Request) {
	view, err := a.svc.ReloadFlow(r.Context(), r.PathValue("realm"), r.PathValue("name"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, view)
}

// updateFlowConfig replaces a flow's config snapshot; a live flow is rebuilt
// with it (issue #46). Body is the config JSON object.
func (a *API) updateFlowConfig(w http.ResponseWriter, r *http.Request) {
	var config json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &config); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	view, err := a.svc.UpdateFlowConfig(r.Context(), r.PathValue("realm"), r.PathValue("name"), config)
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, view)
}

func (a *API) listBlocks(w http.ResponseWriter, r *http.Request) {
	_ = astarteapi.WriteData(w, http.StatusOK, a.svc.ListBlocks(r.PathValue("realm")))
}

func (a *API) getBlock(w http.ResponseWriter, r *http.Request) {
	view, err := a.svc.GetBlock(r.PathValue("realm"), r.PathValue("type"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, view)
}

func (a *API) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrValidation):
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
	case errors.Is(err, flow.ErrFlowExists):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Flow already running")
	case errors.Is(err, flow.ErrFlowNotFound):
		_ = astarteapi.WriteNotFound(w)
	case errors.Is(err, store.ErrAlreadyExists):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Already exists")
	case errors.Is(err, store.ErrPipelineCyclic):
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Pipeline graph contains a cycle")
	case errors.Is(err, store.ErrNotFound):
		_ = astarteapi.WriteNotFound(w)
	default:
		_ = astarteapi.WriteInternalServerError(w)
	}
}

func validationDetail(err error) string {
	msg := err.Error()
	const prefix = "flowapi: validation failed: "
	if len(msg) > len(prefix) && msg[:len(prefix)] == prefix {
		return msg[len(prefix):]
	}
	return msg
}
