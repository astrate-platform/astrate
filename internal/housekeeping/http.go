package housekeeping

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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
	mux.Handle("PATCH /housekeeping/v1/realms/{realm}", h(a.patchRealm))
	mux.Handle("DELETE /housekeeping/v1/realms/{realm}", h(a.deleteRealm))
}

// realmBody is the realm create/get wire shape. Astrate omits the
// Cassandra-specific fields (replication factor/class) upstream carries.
type realmBody struct {
	RealmName                         string `json:"realm_name"`
	JWTPublicKeyPEM                   string `json:"jwt_public_key_pem"`
	DeviceRegistrationLimit           *int32 `json:"device_registration_limit"`
	DatastreamMaximumStorageRetention *int64 `json:"datastream_maximum_storage_retention"`
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
	rv, err := a.svc.CreateRealm(r.Context(), req.RealmName, req.JWTPublicKeyPEM, req.DeviceRegistrationLimit, req.DatastreamMaximumStorageRetention)
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
		RealmName:                         rv.Name,
		JWTPublicKeyPEM:                   rv.JWTPublicKeyPEM,
		DeviceRegistrationLimit:           rv.DeviceRegistrationLimit,
		DatastreamMaximumStorageRetention: rv.DatastreamMaximumStorageRetention,
	}
}

// optString tracks field presence for PATCH decoding; a JSON null decodes to
// the empty value (the blank check rejects it downstream).
type optString struct {
	present bool
	val     string
}

func (o *optString) UnmarshalJSON(b []byte) error {
	o.present = true
	return json.Unmarshal(b, &o.val)
}

// optInt32 tracks presence and nullness of an optional integer field.
type optInt32 struct {
	present bool
	null    bool
	val     int32
}

func (o *optInt32) UnmarshalJSON(b []byte) error {
	o.present = true
	if bytes.Equal(b, []byte("null")) {
		o.null = true
		return nil
	}
	return json.Unmarshal(b, &o.val)
}

// optInt64 tracks presence and nullness of an optional integer field.
type optInt64 struct {
	present bool
	null    bool
	val     int64
}

func (o *optInt64) UnmarshalJSON(b []byte) error {
	o.present = true
	if bytes.Equal(b, []byte("null")) {
		o.null = true
		return nil
	}
	return json.Unmarshal(b, &o.val)
}

// patchBody is the PATCH request's typed view; presence flags distinguish
// absent from null from valued.
type patchBody struct {
	JWTPublicKeyPEM                   optString `json:"jwt_public_key_pem"`
	DeviceRegistrationLimit           optInt32  `json:"device_registration_limit"`
	DatastreamMaximumStorageRetention optInt64  `json:"datastream_maximum_storage_retention"`
}

// patchAllowedFields is the exact set of updatable realm fields (#74);
// anything else is rejected with upstream's invalid_update_parameters shape.
var patchAllowedFields = map[string]struct{}{
	"jwt_public_key_pem":                   {},
	"device_registration_limit":            {},
	"datastream_maximum_storage_retention": {},
}

func (a *API) patchRealm(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes+1))
	if err != nil || int64(len(raw)) > maxBodyBytes {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	var fields map[string]json.RawMessage
	if err := astarteapi.DecodeData(bytes.NewReader(raw), maxBodyBytes, &fields); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	for k := range fields {
		if _, ok := patchAllowedFields[k]; !ok {
			_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
				map[string][]string{"error_name": {"invalid_update_parameters"}})
			return
		}
	}
	var pb patchBody
	if err := astarteapi.DecodeData(bytes.NewReader(raw), maxBodyBytes, &pb); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}

	if pb.JWTPublicKeyPEM.present && pb.JWTPublicKeyPEM.val == "" {
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"jwt_public_key_pem": {"can't be blank"}})
		return
	}
	if pb.DeviceRegistrationLimit.present && !pb.DeviceRegistrationLimit.null && pb.DeviceRegistrationLimit.val < 0 {
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"device_registration_limit": {"is invalid"}})
		return
	}
	if pb.DatastreamMaximumStorageRetention.present && !pb.DatastreamMaximumStorageRetention.null && pb.DatastreamMaximumStorageRetention.val < 0 {
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"datastream_maximum_storage_retention": {"is invalid"}})
		return
	}

	u := RealmUpdate{}
	if pb.JWTPublicKeyPEM.present {
		u.PatchJWTPublicKeyPEM = true
		u.SetJWTPublicKeyPEM = pb.JWTPublicKeyPEM.val
	}
	if pb.DeviceRegistrationLimit.present {
		if pb.DeviceRegistrationLimit.null {
			u.ClearRegistrationLimit = true
		} else {
			u.PatchRegistrationLimit = true
			u.SetRegistrationLimit = pb.DeviceRegistrationLimit.val
		}
	}
	if pb.DatastreamMaximumStorageRetention.present {
		if pb.DatastreamMaximumStorageRetention.null || pb.DatastreamMaximumStorageRetention.val == 0 {
			u.ClearRetention = true
		} else {
			u.PatchRetention = true
			u.SetRetention = pb.DatastreamMaximumStorageRetention.val
		}
	}

	if err := a.svc.UpdateRealm(r.Context(), r.PathValue("realm"), u); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// writeError maps service/store errors onto upstream-shaped responses.
func (a *API) writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrDeletionDisabled):
		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Realm deletion disabled")
	case errors.Is(err, ErrConnectedDevicesPresent):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"error_name": {"connected_devices_present"}})
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
