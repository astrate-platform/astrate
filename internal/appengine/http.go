package appengine

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

// maxBodyBytes caps AppEngine request bodies.
const maxBodyBytes int64 = 1 << 20

// API is the /appengine/v1 HTTP surface (docs/ROADMAP.md §8.2 file 7.8),
// guarded by a realm JWT carrying a_aea.
type API struct {
	svc     *Service
	require func(http.Handler) http.Handler
}

// NewAPI wires the AppEngine service to its HTTP surface.
func NewAPI(svc *Service, mw *auth.Middleware) *API {
	return &API{svc: svc, require: mw.RequireRealm(auth.ClaimAppEngine)}
}

// Mount registers the routes on mux (paths wire-frozen to upstream
// astarte_appengine_api).
func (a *API) Mount(mux *http.ServeMux) {
	h := func(f http.HandlerFunc) http.Handler { return a.require(f) }
	const base = "/appengine/v1/{realm}"
	mux.Handle("GET "+base+"/devices", h(a.listDevices))
	mux.Handle("GET "+base+"/devices/{device}", h(a.getDevice))
	mux.Handle("PATCH "+base+"/devices/{device}", h(a.patchDevice))
	mux.Handle("GET "+base+"/devices-by-alias/{alias}", h(a.getDeviceByAlias))

	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}", h(a.getData))
	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.getData))
	mux.Handle("PUT "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
	mux.Handle("POST "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
	mux.Handle("DELETE "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteData))

	mux.Handle("GET "+base+"/groups", h(a.listGroups))
	mux.Handle("POST "+base+"/groups", h(a.createGroup))
	mux.Handle("GET "+base+"/groups/{group}/devices", h(a.listGroupDevices))
	mux.Handle("POST "+base+"/groups/{group}/devices", h(a.addGroupDevice))
	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}", h(a.removeGroupDevice))
}

// --- devices ----------------------------------------------------------------

func (a *API) listDevices(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	ids, next, err := a.svc.ListDevices(r.Context(), r.PathValue("realm"), r.URL.Query().Get("from_token"), limit)
	if err != nil {
		a.writeError(w, err)
		return
	}
	if next != "" {
		w.Header().Set("X-Astrate-Next-Token", next)
	}
	_ = astarteapi.WriteData(w, http.StatusOK, ids)
}

func (a *API) getDevice(w http.ResponseWriter, r *http.Request) {
	st, err := a.svc.GetDevice(r.Context(), r.PathValue("realm"), r.PathValue("device"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, st)
}

func (a *API) getDeviceByAlias(w http.ResponseWriter, r *http.Request) {
	st, err := a.svc.GetDeviceByAlias(r.Context(), r.PathValue("realm"), r.PathValue("alias"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, st)
}

// devicePatchBody is the PATCH /devices/{id} wire shape.
type devicePatchBody struct {
	Aliases              map[string]*string `json:"aliases"`
	Attributes           map[string]*string `json:"attributes"`
	CredentialsInhibited *bool              `json:"credentials_inhibited"`
}

func (a *API) patchDevice(w http.ResponseWriter, r *http.Request) {
	var body devicePatchBody
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	st, err := a.svc.PatchDevice(r.Context(), r.PathValue("realm"), r.PathValue("device"), DevicePatch(body))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, st)
}

// --- interface data ---------------------------------------------------------

func (a *API) getData(w http.ResponseWriter, r *http.Request) {
	opts, err := parseQueryOpts(r)
	if err != nil {
		_ = astarteapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	data, err := a.svc.GetData(r.Context(), r.PathValue("realm"), r.PathValue("device"),
		r.PathValue("interface"), pathParam(r), opts)
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, data)
}

func (a *API) putData(w http.ResponseWriter, r *http.Request) {
	var value json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &value); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	err := a.svc.PublishData(r.Context(), r.PathValue("realm"), r.PathValue("device"),
		r.PathValue("interface"), pathParam(r), value, nil)
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) deleteData(w http.ResponseWriter, r *http.Request) {
	err := a.svc.UnsetProperty(r.Context(), r.PathValue("realm"), r.PathValue("device"),
		r.PathValue("interface"), pathParam(r))
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- groups -----------------------------------------------------------------

func (a *API) listGroups(w http.ResponseWriter, r *http.Request) {
	names, err := a.svc.ListGroups(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, names)
}

// groupBody is the POST /groups wire shape.
type groupBody struct {
	GroupName string   `json:"group_name"`
	Devices   []string `json:"devices"`
}

func (a *API) createGroup(w http.ResponseWriter, r *http.Request) {
	var body groupBody
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if err := a.svc.CreateGroup(r.Context(), r.PathValue("realm"), body.GroupName, body.Devices); err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusCreated, body)
}

func (a *API) listGroupDevices(w http.ResponseWriter, r *http.Request) {
	ids, err := a.svc.ListGroupDevices(r.Context(), r.PathValue("realm"), r.PathValue("group"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, ids)
}

// groupDeviceBody is the POST /groups/{group}/devices wire shape.
type groupDeviceBody struct {
	DeviceID string `json:"device_id"`
}

func (a *API) addGroupDevice(w http.ResponseWriter, r *http.Request) {
	var body groupDeviceBody
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if err := a.svc.AddGroupDevice(r.Context(), r.PathValue("realm"), r.PathValue("group"), body.DeviceID); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *API) removeGroupDevice(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.RemoveGroupDevice(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device")); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- shared plumbing --------------------------------------------------------

// pathParam reconstructs the Astarte interface path ("/a/b") from the
// trailing {path...} wildcard segment, which arrives without a leading slash.
func pathParam(r *http.Request) string {
	p := r.PathValue("path")
	if p == "" {
		return ""
	}
	return "/" + p
}

// parseQueryOpts reads the datastream query parameters.
func parseQueryOpts(r *http.Request) (QueryOpts, error) {
	q := r.URL.Query()
	var opts QueryOpts
	for _, p := range []struct {
		name string
		dst  **time.Time
	}{
		{"since", &opts.Since}, {"since_after", &opts.SinceAfter}, {"to", &opts.To},
	} {
		if v := q.Get(p.name); v != "" {
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return opts, errors.New(p.name + " is not an RFC 3339 timestamp")
			}
			*p.dst = &t
		}
	}
	if v := q.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			return opts, errors.New("limit is not a non-negative integer")
		}
		opts.Limit = n
	}
	// Upstream default ordering for datastreams is descending (newest first).
	opts.Descending = q.Get("sort") != "ascending"
	return opts, nil
}

// writeError maps service/store errors onto upstream-shaped responses.
func (a *API) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrValidation):
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
	case errors.Is(err, store.ErrAlreadyExists):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Already exists")
	case errors.Is(err, store.ErrNotFound):
		_ = astarteapi.WriteDeviceNotFound(w)
	default:
		_ = astarteapi.WriteInternalServerError(w)
	}
}

// validationDetail strips the ErrValidation prefix for the response detail.
func validationDetail(err error) string {
	const prefix = "appengine: validation failed: "
	msg := err.Error()
	if strings.HasPrefix(msg, prefix) {
		return strings.TrimPrefix(msg, prefix)
	}
	return msg
}
