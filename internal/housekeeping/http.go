package housekeeping

import (
	"errors"
	"net/http"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

// maxBodyBytes caps Housekeeping request bodies (a realm-create body — name,
// a PEM key, a limit — is small).
const maxBodyBytes int64 = 64 << 10

// API is the /housekeeping/v1 HTTP surface (docs/ROADMAP.md §8.1 file 7.4).
// Every route is guarded by an instance-level JWT carrying a_ha.
type API struct {
	svc     *Service
	require func(http.Handler) http.Handler
}

// NewAPI wires the housekeeping service to its HTTP surface. instanceKeysPEM
// are the instance-admin JWT public keys (M3 RequireStatic; config-provided
// in M8).
func NewAPI(svc *Service, mw *auth.Middleware, instanceKeysPEM []string) *API {
	return &API{svc: svc, require: mw.RequireStatic(auth.ClaimHousekeeping, instanceKeysPEM)}
}

// Mount registers the routes on mux (paths wire-frozen to upstream
// astarte_housekeeping).
func (a *API) Mount(mux *http.ServeMux) {
	h := func(f http.HandlerFunc) http.Handler { return a.require(f) }
	mux.Handle("GET /housekeeping/v1/realms", h(a.listRealms))
	mux.Handle("POST /housekeeping/v1/realms", h(a.createRealm))
	mux.Handle("GET /housekeeping/v1/realms/{realm}", h(a.getRealm))
	mux.Handle("DELETE /housekeeping/v1/realms/{realm}", h(a.deleteRealm))
}

// realmBody is the realm create/get wire shape. Astrate omits the
// Cassandra-specific fields (replication factor/class) upstream carries.
type realmBody struct {
	RealmName               string `json:"realm_name"`
	JWTPublicKeyPEM         string `json:"jwt_public_key_pem"`
	DeviceRegistrationLimit *int32 `json:"device_registration_limit"`
}

func (a *API) listRealms(w http.ResponseWriter, r *http.Request) {
	names, err := a.svc.ListRealms(r.Context())
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, names)
}

func (a *API) createRealm(w http.ResponseWriter, r *http.Request) {
	var req realmBody
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &req); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	rv, err := a.svc.CreateRealm(r.Context(), req.RealmName, req.JWTPublicKeyPEM, req.DeviceRegistrationLimit)
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusCreated, viewBody(rv))
}

func (a *API) getRealm(w http.ResponseWriter, r *http.Request) {
	rv, err := a.svc.GetRealm(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, viewBody(rv))
}

func (a *API) deleteRealm(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.DeleteRealm(r.Context(), r.PathValue("realm")); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// viewBody renders a RealmView as its wire shape.
func viewBody(rv *RealmView) realmBody {
	return realmBody{
		RealmName:               rv.Name,
		JWTPublicKeyPEM:         rv.JWTPublicKeyPEM,
		DeviceRegistrationLimit: rv.DeviceRegistrationLimit,
	}
}

// writeError maps service/store errors onto upstream-shaped responses.
func (a *API) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrValidation):
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
	case errors.Is(err, store.ErrAlreadyExists):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Realm already exists")
	case errors.Is(err, store.ErrNotFound):
		_ = astarteapi.WriteNotFound(w)
	default:
		_ = astarteapi.WriteInternalServerError(w)
	}
}

// validationDetail strips the ErrValidation prefix for the response detail.
func validationDetail(err error) string {
	msg := err.Error()
	const prefix = "housekeeping: validation failed: "
	if len(msg) > len(prefix) && msg[:len(prefix)] == prefix {
		return msg[len(prefix):]
	}
	return msg
}
