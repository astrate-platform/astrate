package realm

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

// maxBodyBytes caps Realm Management request bodies (docs/DESIGN.md §4.5):
// an interface definition (≤ 1024 mappings) or a trigger fits in 1 MiB.
const maxBodyBytes int64 = 1 << 20

// API is the /realmmanagement/v1 HTTP surface (docs/ROADMAP.md §8.1 file
// 7.2). Every route is guarded by a realm JWT carrying a_rma.
type API struct {
	svc     *Service
	require func(http.Handler) http.Handler
}

// NewAPI wires the realm service to its HTTP surface. mw provides the
// realm-JWT middleware (M3).
func NewAPI(svc *Service, mw *auth.Middleware) *API {
	return &API{svc: svc, require: mw.RequireRealm(auth.ClaimRealmManagement)}
}

// Mount registers the routes on mux (paths wire-frozen to upstream
// astarte_realm_management).
func (a *API) Mount(mux *http.ServeMux) {
	h := func(f http.HandlerFunc) http.Handler { return a.require(f) }
	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces", h(a.listInterfaces))
	mux.Handle("POST /realmmanagement/v1/{realm}/interfaces", h(a.installInterface))
	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}", h(a.listInterfaceMajors))
	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.getInterface))
	mux.Handle("PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.updateInterface))
	mux.Handle("DELETE /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.deleteInterface))
	mux.Handle("GET /realmmanagement/v1/{realm}/triggers", h(a.listTriggers))
	mux.Handle("POST /realmmanagement/v1/{realm}/triggers", h(a.createTrigger))
	mux.Handle("GET /realmmanagement/v1/{realm}/triggers/{name}", h(a.getTrigger))
	mux.Handle("DELETE /realmmanagement/v1/{realm}/triggers/{name}", h(a.deleteTrigger))
	mux.Handle("GET /realmmanagement/v1/{realm}/config/auth", h(a.getAuth))
	mux.Handle("PUT /realmmanagement/v1/{realm}/config/auth", h(a.putAuth))
}

// --- interfaces -------------------------------------------------------------

func (a *API) listInterfaces(w http.ResponseWriter, r *http.Request) {
	names, err := a.svc.ListInterfaces(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, names)
}

func (a *API) installInterface(w http.ResponseWriter, r *http.Request) {
	var def json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &def); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	si, err := a.svc.InstallInterface(r.Context(), r.PathValue("realm"), def)
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusCreated, json.RawMessage(si.Definition))
}

func (a *API) listInterfaceMajors(w http.ResponseWriter, r *http.Request) {
	majors, err := a.svc.ListInterfaceMajors(r.Context(), r.PathValue("realm"), r.PathValue("name"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, majors)
}

func (a *API) getInterface(w http.ResponseWriter, r *http.Request) {
	major, ok := majorParam(w, r)
	if !ok {
		return
	}
	def, err := a.svc.GetInterface(r.Context(), r.PathValue("realm"), r.PathValue("name"), major)
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, def)
}

func (a *API) updateInterface(w http.ResponseWriter, r *http.Request) {
	if _, ok := majorParam(w, r); !ok {
		return
	}
	var def json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &def); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if _, err := a.svc.UpdateInterface(r.Context(), r.PathValue("realm"), def); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteInterface(w http.ResponseWriter, r *http.Request) {
	major, ok := majorParam(w, r)
	if !ok {
		return
	}
	if err := a.svc.DeleteInterface(r.Context(), r.PathValue("realm"), r.PathValue("name"), major); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- triggers ---------------------------------------------------------------

func (a *API) listTriggers(w http.ResponseWriter, r *http.Request) {
	names, err := a.svc.ListTriggers(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, names)
}

func (a *API) createTrigger(w http.ResponseWriter, r *http.Request) {
	var def json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &def); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	tr, err := a.svc.CreateTrigger(r.Context(), r.PathValue("realm"), def)
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusCreated, json.RawMessage(tr.Definition))
}

func (a *API) getTrigger(w http.ResponseWriter, r *http.Request) {
	def, err := a.svc.GetTrigger(r.Context(), r.PathValue("realm"), r.PathValue("name"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, def)
}

func (a *API) deleteTrigger(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.DeleteTrigger(r.Context(), r.PathValue("realm"), r.PathValue("name")); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- config/auth ------------------------------------------------------------

// authConfig is the GET/PUT /config/auth body shape (upstream
// jwt_public_key_pem field).
type authConfig struct {
	JWTPublicKeyPEM string `json:"jwt_public_key_pem"`
}

func (a *API) getAuth(w http.ResponseWriter, r *http.Request) {
	key, err := a.svc.GetAuthKey(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, authConfig{JWTPublicKeyPEM: key})
}

func (a *API) putAuth(w http.ResponseWriter, r *http.Request) {
	var cfg authConfig
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &cfg); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if err := a.svc.SetAuthKey(r.Context(), r.PathValue("realm"), cfg.JWTPublicKeyPEM); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- shared plumbing --------------------------------------------------------

// majorParam parses the {major} path segment, writing a 404 for a
// non-numeric value (no such resource) and reporting whether it succeeded.
func majorParam(w http.ResponseWriter, r *http.Request) (int, bool) {
	major, err := strconv.Atoi(r.PathValue("major"))
	if err != nil || major < 0 {
		_ = astarteapi.WriteNotFound(w)
		return 0, false
	}
	return major, true
}

// writeError maps service/store errors onto upstream-shaped responses.
func (a *API) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrValidation):
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
	case errors.Is(err, store.ErrAlreadyExists):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Already exists")
	case errors.Is(err, store.ErrInterfaceMajorNotZero):
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity,
			"Interface major version is not 0, can't be deleted")
	case errors.Is(err, store.ErrInterfaceInUse):
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity,
			"Cannot delete an interface that is used by a device introspection")
	case errors.Is(err, store.ErrNotFound):
		_ = astarteapi.WriteNotFound(w)
	default:
		_ = astarteapi.WriteInternalServerError(w)
	}
}

// validationDetail strips the ErrValidation sentinel prefix so the response
// detail is the underlying schema message.
func validationDetail(err error) string {
	msg := err.Error()
	const prefix = "realm: validation failed: "
	if len(msg) > len(prefix) && msg[:len(prefix)] == prefix {
		return msg[len(prefix):]
	}
	return msg
}
