package realm

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
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
	mux.Handle("GET /realmmanagement/v1/{realm}/config/device_registration_limit", h(a.getRegistrationLimit))
	mux.Handle("GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention", h(a.getDatastreamMaximumStorageRetention))
	mux.Handle("GET /realmmanagement/v1/{realm}/version", h(a.getVersion))
	mux.Handle("DELETE /realmmanagement/v1/{realm}/devices/{device}", h(a.deleteDevice))
	mux.Handle("GET /realmmanagement/v1/{realm}/policies", h(a.listPolicies))
	mux.Handle("POST /realmmanagement/v1/{realm}/policies", h(a.createPolicy))
	mux.Handle("GET /realmmanagement/v1/{realm}/policies/{name}", h(a.getPolicy))
	mux.Handle("DELETE /realmmanagement/v1/{realm}/policies/{name}", h(a.deletePolicy))
}

// --- trigger delivery policies ------------------------------------------------

func (a *API) listPolicies(w http.ResponseWriter, r *http.Request) {
	names, err := a.svc.ListPolicies(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, names)
}

func (a *API) createPolicy(w http.ResponseWriter, r *http.Request) {
	var def json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &def); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	p, err := a.svc.CreatePolicy(r.Context(), r.PathValue("realm"), def)
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusCreated, json.RawMessage(p.Definition))
}

func (a *API) getPolicy(w http.ResponseWriter, r *http.Request) {
	def, err := a.svc.GetPolicy(r.Context(), r.PathValue("realm"), r.PathValue("name"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, def)
}

func (a *API) deletePolicy(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.DeletePolicy(r.Context(), r.PathValue("realm"), r.PathValue("name")); err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// deleteDevice synchronously removes a device and its data (the dashboard's
// Delete button; docs/COMPATIBILITY.md notes the deviation from upstream's
// async deletion).
func (a *API) deleteDevice(w http.ResponseWriter, r *http.Request) {
	err := a.svc.DeleteDevice(r.Context(), r.PathValue("realm"), r.PathValue("device"))
	if errors.Is(err, store.ErrNotFound) {
		_ = astarteapi.WriteDeviceNotFound(w)
		return
	}
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// getVersion reports the emulated upstream API level (the dashboard gates
// feature UI on it — see APICompatVersion).
func (a *API) getVersion(w http.ResponseWriter, _ *http.Request) {
	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
}

func (a *API) getRegistrationLimit(w http.ResponseWriter, r *http.Request) {
	limit, err := a.svc.GetDeviceRegistrationLimit(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, limit)
}

// getDatastreamMaximumStorageRetention serves the realm's retention ceiling in
// seconds (upstream 1.2.0+; see Service.GetDatastreamMaximumStorageRetention).
func (a *API) getDatastreamMaximumStorageRetention(w http.ResponseWriter, r *http.Request) {
	retention, err := a.svc.GetDatastreamMaximumStorageRetention(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, retention)
}

// --- interfaces -------------------------------------------------------------

func (a *API) listInterfaces(w http.ResponseWriter, r *http.Request) {
	// ?detailed=true serves the additive 1.4-style detailed listing (#66).
	// Upstream 1.2.0 ignores the parameter entirely (probed), so every other
	// value — absent included — keeps today's names-only response unchanged.
	if r.URL.Query().Get("detailed") == "true" {
		docs, err := a.svc.ListInterfacesDetailed(r.Context(), r.PathValue("realm"))
		if err != nil {
			a.writeError(w, err)
			return
		}
		_ = astarteapi.WriteData(w, http.StatusOK, docs)
		return
	}
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
	var ve *interfaceschema.ViolationsError
	switch {
	case errors.Is(err, ErrMaximumDatabaseRetentionExceeded):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"error_name": {"maximum_database_retention_exceeded"}})
	case errors.As(err, &ve):
		writeViolations(w, ve)
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

// writeViolations renders a *ViolationsError as upstream's Phoenix-changeset
// 422 body (probe-frozen shapes, issue #61): interface-level violations as
// {"errors":{"<field>":["<msg>",...]}} and mapping-level ones as one entry
// per declared mapping in the "mappings" array — {} for clean mappings,
// {"<field>":[...]} for offending ones. Keys are emitted hand-built in
// collection order: Go map marshalling would scramble key order and make
// bodies non-deterministic.
func writeViolations(w http.ResponseWriter, ve *interfaceschema.ViolationsError) {
	body := renderViolationsBody(ve)
	w.Header().Set("Content-Type", astarteapi.ContentType)
	w.WriteHeader(http.StatusUnprocessableEntity)
	_, _ = w.Write(body)
}

// renderViolationsBody encodes the {"errors": {...}} object. Interface-level
// fields come first in collection order, the mappings array last.
func renderViolationsBody(ve *interfaceschema.ViolationsError) []byte {
	var buf bytes.Buffer
	buf.WriteString(`{"errors":{`)
	first := true
	piece := func(key string, value []byte) {
		if !first {
			buf.WriteByte(',')
		}
		first = false
		buf.Write(quoteJSON(key))
		buf.WriteByte(':')
		buf.Write(value)
	}

	// Interface-level violations, merged per field preserving first-seen order.
	type fieldErrors struct {
		field string
		msgs  []string
	}
	var tops []fieldErrors
	topIndex := make(map[string]int)
	for i := range ve.Violations {
		v := &ve.Violations[i]
		if v.MappingIndex >= 0 {
			continue
		}
		if at, ok := topIndex[v.Field]; ok {
			tops[at].msgs = append(tops[at].msgs, v.Messages...)
			continue
		}
		topIndex[v.Field] = len(tops)
		tops = append(tops, fieldErrors{field: v.Field, msgs: append([]string(nil), v.Messages...)})
	}
	for _, f := range tops {
		piece(f.field, marshalStrings(f.msgs))
	}

	// Mapping-level violations as the full-length aligned array.
	if ve.MappingCount > 0 {
		type entry struct {
			fields []string
			msgs   map[string][]string
		}
		entries := make([]entry, ve.MappingCount)
		hasMappings := false
		for i := range ve.Violations {
			v := &ve.Violations[i]
			if v.MappingIndex < 0 || v.MappingIndex >= len(entries) {
				continue
			}
			hasMappings = true
			e := &entries[v.MappingIndex]
			if e.msgs == nil {
				e.msgs = make(map[string][]string)
			}
			if _, ok := e.msgs[v.Field]; !ok {
				e.fields = append(e.fields, v.Field)
			}
			e.msgs[v.Field] = append(e.msgs[v.Field], v.Messages...)
		}
		if hasMappings {
			var arr bytes.Buffer
			arr.WriteByte('[')
			for i := range entries {
				if i > 0 {
					arr.WriteByte(',')
				}
				e := &entries[i]
				if e.msgs == nil {
					arr.WriteString("{}")
					continue
				}
				arr.WriteByte('{')
				for j, f := range e.fields {
					if j > 0 {
						arr.WriteByte(',')
					}
					arr.Write(quoteJSON(f))
					arr.WriteByte(':')
					arr.Write(marshalStrings(e.msgs[f]))
				}
				arr.WriteByte('}')
			}
			arr.WriteByte(']')
			piece("mappings", arr.Bytes())
		}
	}

	buf.WriteString(`}}`)
	return buf.Bytes()
}

// quoteJSON encodes s as a JSON string without HTML escaping, matching the
// astarteapi envelope style. Cannot fail for a string input.
func quoteJSON(s string) []byte {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s) // always appends '\n'
	out := buf.Bytes()
	return out[:len(out)-1]
}

// marshalStrings renders msgs as a JSON array of strings.
func marshalStrings(msgs []string) []byte {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, msg := range msgs {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.Write(quoteJSON(msg))
	}
	buf.WriteByte(']')
	return buf.Bytes()
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
