slug: recipe-docs-sync
verdict: proposed
at:  fa5ed7c
ran: 2026-09-09T16:57:35Z on DietPi in 138s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
✱ Glob "docs/api/*.yaml" in . · 5 matches
→ Read .mule/todo.md
Now I'll pick one surface and do the comparison. Let me check which surfaces have existing docs-sync coverage and find one that's uncovered.

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/astartenative/ 2>/dev/null | head -50
(no output)

$ rg -n '^\s{2}/' docs/api/astrate_native_api.yaml 2>/dev/null | head -50
26:  /astrate/v1/health:
44:  /astrate/v1/readiness:
77:  /astrate/v1/metrics:
99:  /appengine/health:
118:  /realmmanagement/health:
137:  /pairing/health:
158:  /astrate/v1/{realm}/socket:
220:  /appengine/v1/socket/websocket:

✱ Glob "internal/*native*/**/*.go" in . · 0 matches

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*native*' 2>/dev/null | head -50
(no output)

→ Read internal

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/housekeeping/ 2>/dev/null | head -50
(no output)

$ rg -n '^\s{2}/' docs/api/astarte_housekeeping_api.yaml 2>/dev/null
26:  /housekeeping/v1/realms:
106:  /housekeeping/v1/realms/{realm}:

$ rg -n 'HandleFunc|r\.(Get|Post|Put|Delete|Patch)' internal/housekeeping/ 2>/dev/null
(no output)

$ rg -n 'Route|Mount|Group|Handle|http\.' internal/housekeeping/ 2>/dev/null
internal/housekeeping/http.go:23:	require func(http.Handler) http.Handler
internal/housekeeping/http.go:33:// Mount registers the routes on mux (paths wire-frozen to upstream
internal/housekeeping/http.go:35:func (a *API) Mount(mux *http.ServeMux) {
internal/housekeeping/http.go:36:	h := func(f http.HandlerFunc) http.Handler { return a.require(f) }
internal/housekeeping/http.go:37:	mux.Handle("GET /housekeeping/v1/realms", h(a.listRealms))
internal/housekeeping/http.go:38:	mux.Handle("POST /housekeeping/v1/realms", h(a.createRealm))
internal/housekeeping/http.go:39:	mux.Handle("GET /housekeeping/v1/realms/{realm}", h(a.getRealm))
internal/housekeeping/http.go:40:	mux.Handle("PATCH /housekeeping/v1/realms/{realm}", h(a.patchRealm))
internal/housekeeping/http.go:41:	mux.Handle("DELETE /housekeeping/v1/realms/{realm}", h(a.deleteRealm))
internal/housekeeping/http.go:53:func (a *API) listRealms(w http.ResponseWriter, r *http.Request) {
internal/housekeeping/http.go:59:	_ = astarteapi.WriteData(w, http.StatusOK, names)
internal/housekeeping/http.go:62:func (a *API) createRealm(w http.ResponseWriter, r *http.Request) {
internal/housekeeping/http.go:73:	_ = astarteapi.WriteData(w, http.StatusCreated, viewBody(rv))
internal/housekeeping/http.go:76:func (a *API) getRealm(w http.ResponseWriter, r *http.Request) {
internal/housekeeping/http.go:82:	_ = astarteapi.WriteData(w, http.StatusOK, viewBody(rv))
internal/housekeeping/http.go:85:func (a *API) deleteRealm(w http.ResponseWriter, r *http.Request) {
internal/housekeeping/http.go:90:	w.WriteHeader(http.StatusNoContent)
internal/housekeeping/http.go:163:func (a *API) patchRealm(w http.ResponseWriter, r *http.Request) {
internal/housekeeping/http.go:176:			_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
internal/housekeeping/http.go:188:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
internal/housekeeping/http.go:193:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
internal/housekeeping/http.go:198:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
internal/housekeeping/http.go:229:	w.WriteHeader(http.StatusOK)
internal/housekeeping/http.go:233:func (a *API) writeError(w http.ResponseWriter, err error) {
internal/housekeeping/http.go:236:		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Realm deletion disabled")
internal/housekeeping/http.go:238:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
internal/housekeeping/http.go:241:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
internal/housekeeping/http.go:243:		_ = astarteapi.WriteError(w, http.StatusConflict, "Realm already exists")
internal/housekeeping/async_operation_test.go:23:		if rec := r.req(t, http.MethodPost, "/realms"+q, createBody, r.haToken); rec.Code != http.StatusCreated {
internal/housekeeping/async_operation_test.go:27:		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusOK {
internal/housekeeping/async_operation_test.go:30:		if rec := r.req(t, http.MethodDelete, "/realms/"+realmName+q, "", r.haToken); rec.Code != http.StatusNoContent {
internal/housekeeping/async_operation_test.go:33:		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusNotFound {
internal/housekeeping/http_test.go:33:	mux      *http.ServeMux
internal/housekeeping/http_test.go:87:	mux := http.NewServeMux()
internal/housekeeping/http_test.go:91:		auth.NewMiddleware(st), []string{instancePub}).Mount(mux)
internal/housekeeping/http_test.go:122:		if rec := r.req(t, http.MethodGet, "/realms", "", ""); rec.Code != http.StatusUnauthorized {
internal/housekeeping/http_test.go:127:		if rec := r.req(t, http.MethodGet, "/realms", "", r.wrongTok); rec.Code != http.StatusForbidden {
internal/housekeeping/http_test.go:133:		rec := r.req(t, http.MethodPost, "/realms", createBody, r.haToken)
internal/housekeeping/http_test.go:134:		if rec.Code != http.StatusCreated {
internal/housekeeping/http_test.go:143:		if rec := r.req(t, http.MethodPost, "/realms", createBody, r.haToken); rec.Code != http.StatusConflict {
internal/housekeeping/http_test.go:150:		if rec := r.req(t, http.MethodPost, "/realms", body, r.haToken); rec.Code != http.StatusUnprocessableEntity {
internal/housekeeping/http_test.go:156:		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusOK {
internal/housekeeping/http_test.go:159:		if rec := r.req(t, http.MethodGet, "/realms/nope"+randSuffix(t), "", r.haToken); rec.Code != http.StatusNotFound {
internal/housekeeping/http_test.go:163:		decodeData(t, r.req(t, http.MethodGet, "/realms", "", r.haToken), &names)
internal/housekeeping/http_test.go:185:		if rec := r.req(t, http.MethodDelete, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusNoContent {
internal/housekeeping/http_test.go:188:		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusNotFound {
internal/housekeeping/http_test.go:201:	if rec := r.req(t, http.MethodPost, "/realms", createBody, r.haToken); rec.Code != http.StatusCreated {
internal/housekeeping/http_test.go:207:		return r.req(t, http.MethodPatch, "/realms/"+realmName, rawBody, r.haToken)
internal/housekeeping/http_test.go:212:		decodeData(t, r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken), &body)
internal/housekeeping/http_test.go:218:		if rec.Code != http.StatusUnprocessableEntity {
internal/housekeeping/http_test.go:227:		if rec := r.req(t, http.MethodPatch, "/realms/"+realmName, `{"device_registration_limit":7}`, ""); rec.Code != http.StatusUnauthorized {
internal/housekeeping/http_test.go:233:		if rec := patch(t, `{"datastream_maximum_storage_retention":86400}`); rec.Code != http.StatusOK || rec.Body.Len() != 0 {
internal/housekeeping/http_test.go:242:		if rec := patch(t, `{"device_registration_limit":7}`); rec.Code != http.StatusOK || rec.Body.Len() != 0 {
internal/housekeeping/http_test.go:254:		if rec := patch(t, `{"jwt_public_key_pem":`+string(pemJSON)+`}`); rec.Code != http.StatusOK || rec.Body.Len() != 0 {
internal/housekeeping/http_test.go:263:		if rec := patch(t, `{"device_registration_limit":null,"datastream_maximum_storage_retention":null}`); rec.Code != http.StatusOK || rec.Body.Len() != 0 {
internal/housekeeping/http_test.go:273:		if rec := patch(t, `{"datastream_maximum_storage_retention":60}`); rec.Code != http.StatusOK {
internal/housekeeping/http_test.go:276:		if rec := patch(t, `{"datastream_maximum_storage_retention":0}`); rec.Code != http.StatusOK || rec.Body.Len() != 0 {
internal/housekeeping/http_test.go:300:		rec := r.req(t, http.MethodPatch, "/realms/nope"+randSuffix(t), `{"device_registration_limit":7}`, r.haToken)
internal/housekeeping/http_test.go:301:		if rec.Code != http.StatusNotFound {
internal/housekeeping/http_test.go:316:		if rec := r.req(t, http.MethodPost, "/realms", body, r.haToken); rec.Code != http.StatusCreated {
internal/housekeeping/http_test.go:323:		decodeData(t, r.req(t, http.MethodGet, "/realms/"+name, "", r.haToken), &body)
internal/housekeeping/http_test.go:345:		if rec := r.req(t, http.MethodPost, "/realms", body, r.haToken); rec.Code != http.StatusCreated {
internal/housekeeping/http_test.go:349:		decodeData(t, r.req(t, http.MethodGet, "/realms/"+name, "", r.haToken), &rb)
internal/housekeeping/http_test.go:365:		if rec := r.req(t, http.MethodPost, "/realms", body, r.haToken); rec.Code != http.StatusCreated {
internal/housekeeping/http_test.go:385:		if rec := r.req(t, http.MethodDelete, "/realms/nope"+randSuffix(t), "", r.haToken); rec.Code != http.StatusNotFound {
internal/housekeeping/http_test.go:401:		rec := r.req(t, http.MethodDelete, "/realms/"+realmName, "", r.haToken)
internal/housekeeping/http_test.go:402:		if rec.Code != http.StatusUnprocessableEntity {
internal/housekeeping/http_test.go:409:		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusOK {
internal/housekeeping/http_test.go:417:		if rec := r.req(t, http.MethodDelete, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusNoContent {
internal/housekeeping/http_test.go:420:		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusNotFound {
internal/housekeeping/http_test.go:428:		rec := r.req(t, http.MethodDelete, "/realms/"+realmName, "", r.haToken)
internal/housekeeping/http_test.go:429:		if rec.Code != http.StatusMethodNotAllowed {
internal/housekeeping/http_test.go:436:		if rec := r.req(t, http.MethodDelete, "/realms/nope"+randSuffix(t), "", r.haToken); rec.Code != http.StatusMethodNotAllowed {
internal/housekeeping/http_test.go:444:func discardLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

→ Read docs/api/astarte_housekeeping_api.yaml
Now let me read the housekeeping spec and the full http.go handler code to do a proper comparison.
→ Read internal/housekeeping/http.go

$ rg -n 'DatastreamMaximumStorageRetention' internal/housekeeping/ 2>/dev/null
internal/housekeeping/http.go:50:	DatastreamMaximumStorageRetention *int64 `json:"datastream_maximum_storage_retention"`
internal/housekeeping/http.go:68:	rv, err := a.svc.CreateRealm(r.Context(), req.RealmName, req.JWTPublicKeyPEM, req.DeviceRegistrationLimit, req.DatastreamMaximumStorageRetention)
internal/housekeeping/http.go:99:		DatastreamMaximumStorageRetention: rv.DatastreamMaximumStorageRetention,
internal/housekeeping/http.go:152:	DatastreamMaximumStorageRetention optInt64  `json:"datastream_maximum_storage_retention"`
internal/housekeeping/http.go:197:	if pb.DatastreamMaximumStorageRetention.present && !pb.DatastreamMaximumStorageRetention.null && pb.DatastreamMaximumStorageRetention.val < 0 {
internal/housekeeping/http.go:216:	if pb.DatastreamMaximumStorageRetention.present {
internal/housekeeping/http.go:217:		if pb.DatastreamMaximumStorageRetention.null || pb.DatastreamMaximumStorageRetention.val == 0 {
internal/housekeeping/http.go:221:			u.SetRetention = pb.DatastreamMaximumStorageRetention.val
internal/housekeeping/service_test.go:40:		DatastreamMaximumStorageRetention: nr.DatastreamMaximumStorageRetention,
internal/housekeeping/service_test.go:137:			WithDefaultDatastreamMaximumStorageRetention(&def)
internal/housekeeping/service_test.go:146:		if got := st.created[0].DatastreamMaximumStorageRetention; got == nil || *got != def {
internal/housekeeping/service_test.go:149:		if got := rv.DatastreamMaximumStorageRetention; got == nil || *got != def {
internal/housekeeping/service_test.go:158:			WithDefaultDatastreamMaximumStorageRetention(&def)
internal/housekeeping/service_test.go:163:		if got := st.created[0].DatastreamMaximumStorageRetention; got == nil || *got != explicit {
internal/housekeeping/service_test.go:175:		if st.created[0].DatastreamMaximumStorageRetention != nil {
internal/housekeeping/service_test.go:176:			t.Errorf("store received retention %v, want nil", *st.created[0].DatastreamMaximumStorageRetention)
internal/housekeeping/http_test.go:89:		WithDefaultDatastreamMaximumStorageRetention(opts.defaultRetention).
internal/housekeeping/http_test.go:236:		if body := getBody(t); body.DatastreamMaximumStorageRetention == nil || *body.DatastreamMaximumStorageRetention != 86400 {
internal/housekeeping/http_test.go:237:			t.Errorf("retention after patch = %v, want 86400", body.DatastreamMaximumStorageRetention)
internal/housekeeping/http_test.go:267:		if body.DeviceRegistrationLimit != nil || body.DatastreamMaximumStorageRetention != nil {
internal/housekeeping/http_test.go:268:			t.Errorf("fields after null patch = %v/%v, want nil/nil", body.DeviceRegistrationLimit, body.DatastreamMaximumStorageRetention)
internal/housekeeping/http_test.go:279:		if body := getBody(t); body.DatastreamMaximumStorageRetention != nil {
internal/housekeeping/http_test.go:280:			t.Errorf("retention after zero patch = %v, want nil", *body.DatastreamMaximumStorageRetention)
internal/housekeeping/http_test.go:324:		return body.DatastreamMaximumStorageRetention
internal/housekeeping/http_test.go:350:		if rb.DatastreamMaximumStorageRetention != nil {
internal/housekeeping/http_test.go:351:			t.Errorf("retention without default = %v, want nil", *rb.DatastreamMaximumStorageRetention)
internal/housekeeping/service.go:40:	DatastreamMaximumStorageRetention *int64
internal/housekeeping/service.go:98:// WithDefaultDatastreamMaximumStorageRetention sets the realm-level datastream
internal/housekeeping/service.go:103:func (s *Service) WithDefaultDatastreamMaximumStorageRetention(defaultRetention *int64) *Service {
internal/housekeeping/service.go:168:		DatastreamMaximumStorageRetention: retention,
internal/housekeeping/service.go:278:		DatastreamMaximumStorageRetention: r.DatastreamMaximumStorageRetention,

Let me verify the `datastream_maximum_storage_retention` field in the service layer to confirm it's a real persisted field.
$ rg -n 'datastream_maximum_storage_retention' docs/api/astarte_housekeeping_api.yaml 2>/dev/null
(no output)

Now I have the full picture. Two clear drift items on the housekeeping surface:

1. **PATCH `/housekeeping/v1/realms/{realm}` exists in code but is entirely absent from the YAML spec**
2. **`datastream_maximum_storage_retention` is a real field on the wire (accepted, persisted, returned) but absent from both the `Realm` and `RealmCreate` schemas**

Let me also check if the PATCH handler's error surface (422 with `invalid_update_parameters`, plus the validation errors) would need documenting.
→ Read internal/housekeeping/service.go [offset=1, limit=60]
I have the comparison done. Two clear `docs/api/` drift items:

1. **PATCH `/housekeeping/v1/realms/{realm}`** is in code (http.go:40) but entirely absent from the YAML
2. **`datastream_maximum_storage_retention`** is a real wire field (persisted, returned, accepted in create + patch) but absent from both the `Realm` and `RealmCreate` schemas

Now I'll append the task lines to `.mule/todo.md`.
→ Read .mule/todo.md [offset=130]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -129,4 +129,6 @@
 - [x] docs-sync-appengine-downsample-min [auto]: fix the `downsample_to` constraint in docs/api/astrate_appengine_api.yaml — the DataDownsample description says "(must be >= 2)" (yaml:1357) but parseQueryOpts rejects `downsample_to <= 2` with "must be greater than 2" (internal/appengine/http.go:592-593), so the documented minimum is off by one: readings of `2` (allowed per the doc) actually 422. Change the description to "(must be > 2)" (the 422 error shape is already documented). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 
 - [x] pairing-bearer-secret-test [auto]: add a container-free unit test for `bearerSecret` (internal/pairing/http.go:343) locking the upstream-verbatim permissive match `~r/bearer\:?\s+(.*)$/i` — the colon form `Bearer: <secret>`, the case variants (`bearer `, `BEARER `), and the reject paths (no whitespace before value, non-bearer scheme, empty value). Only the plain `Bearer <secret>` form is exercised today (integration http_test.go); a regression to case-sensitive or drop-colon matching would pass every existing test. Pure helper in a non-integration file; no Docker.
 - [x] pairing-remoteip-fallback-test [auto]: add a container-free unit test for `remoteIP` (internal/pairing/http.go:359) — a valid addr:port yields `ap.Addr().Unmap()`, a bare address yields the mapped address, and an unparseable `RemoteAddr` (unix-socket/proxy case) falls back to `netip.IPv4Unspecified()`. That fallback branch is currently uncovered (only reachable through the integration HTTP stack). Pure helper in a non-integration file; no Docker.
+- [ ] docs-sync-hk-patch-endpoint [auto]: add the undocumented `PATCH /housekeeping/v1/realms/{realm}` route to docs/api/astarte_housekeeping_api.yaml — it exists in code (internal/housekeeping/http.go:40, handler patchRealm at http.go:163-229, returns 200 on success) but the spec jumps straight from GET to DELETE. Document the request body (`RealmPatch` — same three optional fields as `RealmCreate`: `jwt_public_key_pem`, `device_registration_limit`, `datastream_maximum_storage_retention`, all nullable), the success response, and the error responses: 400 (bad envelope), 401, 403, 404, 405 (deletion disabled), 422 (invalid update parameters for unknown fields, blank jwt_public_key_pem, negative limit/retention, connected devices present), 500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-hk-retention-field [auto]: add the `datastream_maximum_storage_retention` field to the `Realm` and `RealmCreate` schemas in docs/api/astarte_housekeeping_api.yaml — the field is on the wire (realmBody struct, internal/housekeeping/http.go:50, returned by viewBody at http.go:99, accepted in create at http.go:68 and patch at http.go:216) but absent from both schemas (yaml:200-231). Type: integer, format int64, nullable, description matching the device_registration_limit style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

→ Read .mule/todo.md [offset=130]
Done. Two task lines appended to `.mule/todo.md`:

- **`docs-sync-hk-patch-endpoint`** (line 133): PATCH `/housekeeping/v1/realms/{realm}` exists in code (`http.go:40`) but is entirely absent from the YAML spec. The handler (`patchRealm`, http.go:163-229) accepts the same three optional fields as create, returns 200, and has a rich error surface (405/422/404/etc.) that needs documenting.

- **`docs-sync-hk-retention-field`** (line 134): `datastream_maximum_storage_retention` is a real wire field (persisted, returned, accepted in both create and patch) but absent from both the `Realm` and `RealmCreate` schemas in the YAML.

No prose drift found on the housekeeping surface — the two existing spec pages (`GET /realms`, `GET/POST/DELETE /realms/{realm}`) accurately reflect the code's status codes and response shapes.

Done: proposed 2 docs/api fix tasks for the housekeeping surface
Files: .mule/todo.md (lines 133-134)
Verified: housekeeping code routes vs YAML paths, status codes, schemas — 2 drift items found
Unsure: nothing
Follow-ups: docs-sync-hk-patch-endpoint, docs-sync-hk-retention-field
