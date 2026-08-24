package appengine

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/engine"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

// maxBodyBytes caps AppEngine request bodies.
const maxBodyBytes int64 = 1 << 20

// maxValueBytes is upstream's per-value limit (~64 KiB): a larger REST
// write is answered 422 "Value size exceeds size limits" (measured
// 2026-08-24 against 1.2.0, verify-server-writes.json).
const maxValueBytes = 65536

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
	mux.Handle("GET "+base+"/stats/devices", h(a.devicesStats))
	mux.Handle("GET "+base+"/devices/{device}", h(a.getDevice))
	mux.Handle("PATCH "+base+"/devices/{device}", h(a.patchDevice))
	mux.Handle("GET "+base+"/devices-by-alias/{alias}", h(a.getDeviceByAlias))
	mux.Handle("GET "+base+"/devices/{device}/interfaces", h(a.listDeviceInterfaces))

	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}", h(a.getData))
	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.getData))
	mux.Handle("PUT "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
	mux.Handle("POST "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
	mux.Handle("DELETE "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteData))

	// Mirror surfaces (upstream registers these alongside the device-scoped
	// ones; Go 1.22 precedence resolves the overlap with the group routes
	// above/below without re-registering them).
	mux.Handle("PATCH "+base+"/devices-by-alias/{alias}", h(a.patchDeviceByAlias))
	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces", h(a.listInterfacesByAlias))
	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces/{interface}", h(a.getDataByAlias))
	mux.Handle("GET "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.getDataByAlias))
	mux.Handle("PUT "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.putDataByAlias))
	mux.Handle("POST "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.putDataByAlias))
	mux.Handle("DELETE "+base+"/devices-by-alias/{alias}/interfaces/{interface}/{path...}", h(a.deleteDataByAlias))

	mux.Handle("GET "+base+"/groups", h(a.listGroups))
	mux.Handle("POST "+base+"/groups", h(a.createGroup))
	mux.Handle("GET "+base+"/groups/{group}", h(a.getGroup))
	mux.Handle("GET "+base+"/groups/{group}/devices", h(a.listGroupDevices))
	mux.Handle("POST "+base+"/groups/{group}/devices", h(a.addGroupDevice))
	mux.Handle("GET "+base+"/groups/{group}/devices/{device}", h(a.getGroupDevice))
	mux.Handle("PATCH "+base+"/groups/{group}/devices/{device}", h(a.patchGroupDevice))
	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}", h(a.removeGroupDevice))
	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces", h(a.listInterfacesInGroup))
	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}", h(a.getDataInGroup))
	mux.Handle("GET "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.getDataInGroup))
	mux.Handle("PUT "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.putDataInGroup))
	mux.Handle("POST "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.putDataInGroup))
	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteDataInGroup))
}

// --- devices ----------------------------------------------------------------

func (a *API) listDevices(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	details, _ := strconv.ParseBool(q.Get("details"))
	page, err := a.svc.ListDevices(r.Context(), r.PathValue("realm"), q.Get("from_token"), limit, details)
	if err != nil {
		a.writeError(w, err)
		return
	}
	var body any = page.IDs
	if details {
		body = page.Statuses
	}
	_ = astarteapi.WriteDataWithLinks(w, http.StatusOK, body,
		deviceListLinks(r.PathValue("realm"), q, page.Next))
}

// deviceListLinks builds the upstream pagination links object. links.next
// always carries details, from_token, and limit: url.Values.Encode sorts
// keys, so "details" always precedes "from_token" — the dashboard parses
// links.next with URLSearchParams over the raw path+query string, which
// corrupts the FIRST key/value pair, so from_token must never be first.
// Absent request values are filled with their effective defaults.
func deviceListLinks(realm string, q url.Values, next string) astarteapi.Links {
	base := "/v1/" + realm + "/devices"
	self := base
	if enc := q.Encode(); enc != "" {
		self += "?" + enc
	}
	links := astarteapi.Links{Self: self}
	if next != "" {
		details := q.Get("details")
		if details == "" {
			details = "false"
		}
		limit := q.Get("limit")
		if limit == "" {
			limit = strconv.Itoa(DefaultDeviceLimit)
		}
		nq := url.Values{"details": {details}, "from_token": {next}, "limit": {limit}}
		links.Next = base + "?" + nq.Encode()
	}
	return links
}

func (a *API) devicesStats(w http.ResponseWriter, r *http.Request) {
	total, connected, err := a.svc.DevicesStats(r.Context(), r.PathValue("realm"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, map[string]int64{
		"connected_devices": connected,
		"total_devices":     total,
	})
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

// writeInterfaces renders an interface-name list; names must be non-nil so an
// empty introspection renders as [] rather than null.
func (a *API) writeInterfaces(w http.ResponseWriter, names []string) {
	if names == nil {
		names = []string{}
	}
	_ = astarteapi.WriteData(w, http.StatusOK, names)
}

func (a *API) listDeviceInterfaces(w http.ResponseWriter, r *http.Request) {
	names, err := a.svc.ListInterfaces(r.Context(), r.PathValue("realm"), r.PathValue("device"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	a.writeInterfaces(w, names)
}

func (a *API) listInterfacesByAlias(w http.ResponseWriter, r *http.Request) {
	names, err := a.svc.ListInterfacesByAlias(r.Context(), r.PathValue("realm"), r.PathValue("alias"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	a.writeInterfaces(w, names)
}

func (a *API) listInterfacesInGroup(w http.ResponseWriter, r *http.Request) {
	names, err := a.svc.ListInterfacesInGroup(r.Context(), r.PathValue("realm"),
		r.PathValue("group"), r.PathValue("device"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	a.writeInterfaces(w, names)
}

// devicePatchBody is the PATCH /devices/{id} wire shape.
type devicePatchBody struct {
	Aliases              map[string]*string `json:"aliases"`
	Attributes           map[string]*string `json:"attributes"`
	CredentialsInhibited *bool              `json:"credentials_inhibited"`
}

func (a *API) patchDevice(w http.ResponseWriter, r *http.Request) {
	// Upstream honors PATCH only with exactly this content type (a header
	// comparison, so no charset parameters). Anything else fails its merge
	// with :patch_mimetype_not_supported — which no FallbackController clause
	// maps, so Phoenix surfaces an unmapped 500. Replicated as-is.
	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
		_ = astarteapi.WriteInternalServerError(w)
		return
	}
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

func (a *API) patchDeviceByAlias(w http.ResponseWriter, r *http.Request) {
	// Same content-type discipline as patchDevice, but this controller's
	// fallback maps the mismatch to 400 "Bad request" upstream (probed live),
	// unlike the device path's unmapped 500.
	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Bad request")
		return
	}
	var body devicePatchBody
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	st, err := a.svc.PatchDeviceByAlias(r.Context(), r.PathValue("realm"), r.PathValue("alias"), DevicePatch(body))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, st)
}

// --- interface data ---------------------------------------------------------

// serveData runs the shared tail of every interface-data GET: query parsing,
// the service read (via the caller's resolution), and the envelope write.
// Parse failures render as 422 changeset field errors, matching upstream
// param casting; a *Tabular result renders {data, metadata}.
func (a *API) serveData(w http.ResponseWriter, r *http.Request, read func(QueryOpts) (any, error)) {
	opts, err := parseQueryOpts(r)
	if err != nil {
		a.writeError(w, err)
		return
	}
	data, err := read(opts)
	if err != nil {
		a.writeError(w, err)
		return
	}
	if t, ok := data.(*Tabular); ok {
		_ = astarteapi.WriteDataWithMetadata(w, http.StatusOK, t.Data, t.Metadata)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, data)
}

func (a *API) getData(w http.ResponseWriter, r *http.Request) {
	a.serveData(w, r, func(opts QueryOpts) (any, error) {
		return a.svc.GetData(r.Context(), r.PathValue("realm"), r.PathValue("device"),
			r.PathValue("interface"), pathParam(r), opts)
	})
}

func (a *API) getDataByAlias(w http.ResponseWriter, r *http.Request) {
	a.serveData(w, r, func(opts QueryOpts) (any, error) {
		return a.svc.GetDataByAlias(r.Context(), r.PathValue("realm"), r.PathValue("alias"),
			r.PathValue("interface"), pathParam(r), opts)
	})
}

func (a *API) getDataInGroup(w http.ResponseWriter, r *http.Request) {
	a.serveData(w, r, func(opts QueryOpts) (any, error) {
		return a.svc.GetDataInGroup(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device"),
			r.PathValue("interface"), pathParam(r), opts)
	})
}

func (a *API) putData(w http.ResponseWriter, r *http.Request) {
	var value json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &value); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if len(value) > maxValueBytes {
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
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

func (a *API) putDataByAlias(w http.ResponseWriter, r *http.Request) {
	var value json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &value); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if len(value) > maxValueBytes {
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
		return
	}
	err := a.svc.PublishDataByAlias(r.Context(), r.PathValue("realm"), r.PathValue("alias"),
		r.PathValue("interface"), pathParam(r), value, nil)
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) putDataInGroup(w http.ResponseWriter, r *http.Request) {
	var value json.RawMessage
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &value); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if len(value) > maxValueBytes {
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
		return
	}
	err := a.svc.PublishDataInGroup(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device"),
		r.PathValue("interface"), pathParam(r), value, nil)
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) deleteDataByAlias(w http.ResponseWriter, r *http.Request) {
	err := a.svc.UnsetPropertyByAlias(r.Context(), r.PathValue("realm"), r.PathValue("alias"),
		r.PathValue("interface"), pathParam(r))
	if err != nil {
		a.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *API) deleteDataInGroup(w http.ResponseWriter, r *http.Request) {
	err := a.svc.UnsetPropertyInGroup(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device"),
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

// groupBody is the POST /groups wire shape. Devices is a pointer so a MISSING
// key is distinguishable from an empty array (upstream rejects both, with
// different messages).
type groupBody struct {
	GroupName string    `json:"group_name"`
	Devices   *[]string `json:"devices"`
}

func (a *API) createGroup(w http.ResponseWriter, r *http.Request) {
	var body groupBody
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	switch {
	case body.GroupName == "":
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			FieldErrors{"group_name": {"can't be blank"}})
		return
	case body.Devices == nil:
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			FieldErrors{"devices": {"can't be blank"}})
		return
	case len(*body.Devices) == 0:
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			FieldErrors{"devices": {"should have at least 1 item(s)"}})
		return
	}
	if err := a.svc.CreateGroup(r.Context(), r.PathValue("realm"), body.GroupName, *body.Devices); err != nil {
		a.writeError(w, err)
		return
	}
	// The 201 echo carries the ORIGINAL body list, duplicates included.
	_ = astarteapi.WriteData(w, http.StatusCreated, body)
}

// getGroup serves GET /groups/{group}: the group identity body, distinct from
// listing its members.
func (a *API) getGroup(w http.ResponseWriter, r *http.Request) {
	if err := a.svc.GetGroup(r.Context(), r.PathValue("realm"), r.PathValue("group")); err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, map[string]string{"group_name": r.PathValue("group")})
}

func (a *API) listGroupDevices(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	details, _ := strconv.ParseBool(q.Get("details"))
	limit, _ := strconv.Atoi(q.Get("limit"))
	page, err := a.svc.ListGroupDevices(r.Context(), r.PathValue("realm"), r.PathValue("group"),
		details, q.Get("from_token"), limit)
	if err != nil {
		a.writeError(w, err)
		return
	}
	var body any = page.IDs
	if details {
		body = page.Statuses
	}
	_ = astarteapi.WriteDataWithLinks(w, http.StatusOK, body,
		groupListLinks(r.PathValue("realm"), r.PathValue("group"), q, page.Next))
}

// groupListLinks builds the upstream pagination links for a group-device
// listing. Self echoes the request query only when one was sent; next carries
// from_token plus the effective limit (and details when the request had it),
// matching the probed upstream bodies — url.Values.Encode sorts the keys, so
// details precedes from_token precedes limit.
func groupListLinks(realm, group string, q url.Values, next string) astarteapi.Links {
	base := "/v1/" + realm + "/groups/" + group + "/devices"
	self := base
	if enc := q.Encode(); enc != "" {
		self += "?" + enc
	}
	links := astarteapi.Links{Self: self}
	if next == "" {
		return links
	}
	limit := q.Get("limit")
	if limit == "" {
		limit = strconv.Itoa(DefaultDeviceLimit)
	}
	nq := url.Values{"from_token": {next}, "limit": {limit}}
	if d := q.Get("details"); d != "" {
		nq.Set("details", d)
	}
	links.Next = base + "?" + nq.Encode()
	return links
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

// getGroupDevice serves GET /groups/{group}/devices/{device}: the normal
// device status projection behind the membership gate.
func (a *API) getGroupDevice(w http.ResponseWriter, r *http.Request) {
	st, err := a.svc.GetGroupDevice(r.Context(), r.PathValue("realm"), r.PathValue("group"), r.PathValue("device"))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, st)
}

func (a *API) patchGroupDevice(w http.ResponseWriter, r *http.Request) {
	// Same header comparison as patchDevice; upstream's group-scoped fallback
	// has no clause for the mismatch either, so it also Phoenix-500s.
	if r.Header.Get("Content-Type") != "application/merge-patch+json" {
		_ = astarteapi.WriteInternalServerError(w)
		return
	}
	var body devicePatchBody
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &body); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	st, err := a.svc.PatchGroupDevice(r.Context(), r.PathValue("realm"), r.PathValue("group"),
		r.PathValue("device"), DevicePatch(body))
	if err != nil {
		a.writeError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusOK, st)
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

// parseQueryOpts reads the datastream query parameters. Malformed values come
// back as FieldErrors so writeError renders upstream's changeset-shaped 422
// ({"errors":{"<param>":["is invalid"]}}), not a generic 400.
func parseQueryOpts(r *http.Request) (QueryOpts, error) {
	q := r.URL.Query()
	var opts QueryOpts
	fe := FieldErrors{}
	for _, p := range []struct {
		name string
		dst  **time.Time
	}{
		{"since", &opts.Since}, {"since_after", &opts.SinceAfter}, {"to", &opts.To},
	} {
		if v := q.Get(p.name); v != "" {
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				fe.addf(p.name, "is invalid")
				continue
			}
			*p.dst = &t
		}
	}
	if q.Get("since") != "" && q.Get("since_after") != "" {
		fe.addf("since_after", "conflicts already set parameter")
	}
	if v := q.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		switch {
		case err != nil:
			fe.addf("limit", "is invalid")
		case n < 0:
			fe.addf("limit", "must be greater than or equal to 0")
		default:
			opts.Limit = n
		}
	}
	if v := q.Get("downsample_to"); v != "" {
		n, err := strconv.Atoi(v)
		switch {
		case err != nil:
			fe.addf("downsample_to", "is invalid")
		case n <= 2:
			fe.addf("downsample_to", "must be greater than 2")
		default:
			opts.DownsamplePoints = n
		}
	}
	// The default is cast before validation so a rejected value still leaves
	// a usable format in opts.
	opts.Format = "structured"
	switch v := q.Get("format"); v {
	case "":
	case "structured", "table", "disjoint_tables":
		opts.Format = v
	default:
		fe.addf("format", "is invalid")
	}
	for _, p := range []struct {
		name string
		dst  **bool
	}{
		{"allow_bigintegers", &opts.AllowBigIntegers},
		{"allow_safe_bigintegers", &opts.AllowSafeBigIntegers},
	} {
		switch v := q.Get(p.name); v {
		case "":
		case "true", "false":
			b := v == "true"
			*p.dst = &b
		default:
			fe.addf(p.name, "is invalid")
		}
	}
	if v := q.Get("retrieve_metadata"); v != "" {
		switch v {
		case "true":
			opts.RetrieveMetadata = true
		case "false":
		default:
			fe.addf("retrieve_metadata", "is invalid")
		}
	}
	opts.DownsampleKey = q.Get("downsample_key")
	// Upstream default ordering for datastreams is descending (newest first).
	opts.Descending = q.Get("sort") != "ascending"
	if len(fe) > 0 {
		return opts, fe
	}
	return opts, nil
}

// writeError maps service/store errors onto upstream-shaped responses.
func (a *API) writeError(w http.ResponseWriter, err error) {
	var fe FieldErrors
	switch {
	case errors.As(err, &fe):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity, fe)
	case errors.Is(err, ErrInvalidAlias):
		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Invalid alias")
	case errors.Is(err, ErrAliasAlreadyInUse):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Alias already in use")
	case errors.Is(err, ErrAliasTagNotFound):
		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Alias tag not found")
	case errors.Is(err, ErrInvalidAttributes):
		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Invalid attributes")
	case errors.Is(err, ErrAttributeKeyNotFound):
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Attribute key not found")
	case errors.Is(err, ErrValidation):
		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
	case errors.Is(err, ErrGroupAlreadyExists):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Group already exists")
	case errors.Is(err, ErrDeviceAlreadyInGroup):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Device already in group")
	case errors.Is(err, store.ErrAlreadyExists):
		_ = astarteapi.WriteError(w, http.StatusConflict, "Already exists")
	case errors.Is(err, ErrGroupNotFound):
		_ = astarteapi.WriteError(w, http.StatusNotFound, "Group not found")
	case errors.Is(err, ErrPathNotFound):
		_ = astarteapi.WriteError(w, http.StatusNotFound, "Path not found")
	// Engine write sentinels must precede store.ErrNotFound: wrapped
	// ErrInterfaceNotFound errors are distinct from store's, and Go's switch
	// takes the first match (issue #57, measured taxonomy).
	case errors.Is(err, engine.ErrNotServerOwned):
		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Cannot write to device owned resource")
	case errors.Is(err, engine.ErrNotAProperty):
		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Cannot write to read-only resource")
	case errors.Is(err, engine.ErrInterfaceNotFound):
		_ = astarteapi.WriteError(w, http.StatusNotFound, "Interface not found in device introspection")
	case errors.Is(err, engine.ErrPathNotFound):
		_ = astarteapi.WriteError(w, http.StatusBadRequest, "Endpoint not found")
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
