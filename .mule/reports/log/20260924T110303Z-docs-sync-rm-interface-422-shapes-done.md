slug: docs-sync-rm-interface-422-shapes
verdict: done
at:  93e784e
ran: 2026-09-24T11:03:03Z on DietPi in 490s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_realm_management_api.yaml [offset=100, limit=290]

$ ls /root/astrate-mule && rg -l "writeInterfaceError" /root/astrate-mule --type go
AGENTS.md
Dockerfile
LICENSE
Makefile
README.md
bench
cmd
deploy
docker-compose.yml
docs
examples
flow-container-echo
go.mod
go.sum
internal
migrations
opencode.json
pkg
site-dist
test
tools
/root/astrate-mule/internal/realm/http.go

$ rg -n "ValidationError|TriggerErrorsDetail|ErrorDetail" /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
134:          $ref: "#/components/responses/ValidationError"
234:          $ref: "#/components/responses/ValidationError"
264:                $ref: "#/components/schemas/ErrorDetail"
352:                  - $ref: "#/components/schemas/ErrorDetail"
353:                  - $ref: "#/components/schemas/TriggerErrorsDetail"
491:          $ref: "#/components/responses/ValidationError"
686:          $ref: "#/components/responses/ValidationError"
738:          $ref: "#/components/responses/ValidationError"
1068:    ErrorDetail:
1079:    TriggerErrorsDetail:
1110:            $ref: "#/components/schemas/ErrorDetail"
1120:            $ref: "#/components/schemas/ErrorDetail"
1130:            $ref: "#/components/schemas/ErrorDetail"
1140:            $ref: "#/components/schemas/ErrorDetail"
1150:            $ref: "#/components/schemas/ErrorDetail"
1168:    ValidationError:
1173:            $ref: "#/components/schemas/ErrorDetail"
1183:            $ref: "#/components/schemas/ErrorDetail"

$ rg -n "writeInterfaceError|writeViolations|validationDetail|ErrValidation|ErrMaximumDatabaseRetentionExceeded" /root/astrate-mule/internal/realm/http.go
153:			a.writeInterfaceError(w, err)
161:		a.writeInterfaceError(w, err)
175:		a.writeInterfaceError(w, err)
184:		a.writeInterfaceError(w, err)
197:		a.writeInterfaceError(w, err)
214:		a.writeInterfaceError(w, err)
226:		a.writeInterfaceError(w, err)
347:// writeInterfaceError maps interface-service/store errors onto the upstream
352:func (a *API) writeInterfaceError(w http.ResponseWriter, err error) {
355:	case errors.Is(err, ErrMaximumDatabaseRetentionExceeded):
359:		writeViolations(w, ve)
392:	case errors.Is(err, ErrValidation):
393:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
405:	case errors.Is(err, ErrMaximumDatabaseRetentionExceeded):
409:		writeViolations(w, ve)
410:	case errors.Is(err, ErrValidation):
411:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
427:// writeViolations renders a *ViolationsError as upstream's Phoenix-changeset
434:func writeViolations(w http.ResponseWriter, ve *interfaceschema.ViolationsError) {
560:// validationDetail strips the ErrValidation sentinel prefix so the response
562:func validationDetail(err error) string {

$ rg -n "maximum_database_retention_exceeded|should be at most|mappings" /root/astrate-mule/internal/realm/http_test.go
34:	ifaceV1  = `{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`
35:	ifaceV1b = `{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":1,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"},{"endpoint":"/count","type":"integer"}]}`
36:	ifaceV1x = `{"interface_name":"com.ex.M7a.Sensors","version_major":1,"version_minor":2,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"integer"}]}`
37:	ifaceV2  = `{"interface_name":"com.ex.M7a.Sensors","version_major":2,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"string"}]}`
38:	draftV0  = `{"interface_name":"com.ex.M7a.Draft","version_major":0,"version_minor":1,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/x","type":"double"}]}`
279:			`"mappings":[{"endpoint":"/value","type":"double","database_retention_policy":"use_ttl","database_retention_ttl":%d}]}`,
282:	const overEnvelope = `{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`
353:	sensors := func(minor int, mappings string) string {
355:			`"type":"datastream","ownership":"device","mappings":[%s]}`, minor, mappings)
359:			`"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`, name)
421:				`"ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`, rmIface),
441:				`"ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`, rmIface),
476:// aligned mappings array.
483:			`"mappings":[{"endpoint":"/value","type":"double"}]}`, strings.Repeat("x", 1001))
484:		want := `{"errors":{"description":["should be at most 1000 character(s)"]}}`
493:			`"type":"datastream","ownership":"device","mappings":[`+
496:		want := `{"errors":{"mappings":[{},{"description":["should be at most 1000 character(s)"]}]}}`
512:		`"mappings":[{"path":"/value","type":"string"}]}`
523:		} `json:"mappings"`
541:			`"type":"datastream","ownership":"device","mappings":[`+
547:	want := `{"errors":{"mappings":[{"database_retention_ttl":["must be greater than or equal to 60"]}]}}`

→ Read internal/realm/http.go [limit=240, offset=340]
→ Read docs/api/astarte_realm_management_api.yaml [limit=130, offset=1060]
→ Read internal/realm/http_test.go [limit=70, offset=270]
→ Read internal/realm/http_test.go [limit=80, offset=470]

$ rg -n "WriteFieldErrors|func WriteError" /root/astrate-mule/internal -g "*.go"
/root/astrate-mule/internal/pairing/http.go:121:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/pairing/http.go:174:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/pairing/http.go:254:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/pairing/http.go:308:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/pairing/http.go:311:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/pairing/http.go:314:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/pairing/http.go:317:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/pairing/http.go:320:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/realm/http.go:356:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/realm/http.go:406:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/appengine/http.go:406:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/appengine/http.go:410:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/appengine/http.go:414:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/appengine/http.go:647:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity, fe)
/root/astrate-mule/internal/housekeeping/http.go:176:			_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/housekeeping/http.go:188:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/housekeeping/http.go:193:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/housekeeping/http.go:198:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
/root/astrate-mule/internal/housekeeping/http.go:238:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,

$ rg -n "func WriteFieldErrors" /root/astrate-mule -g "*.go" -l
/root/astrate-mule/pkg/astarteapi/envelope.go

→ Read pkg/astarteapi/envelope.go
→ Read docs/api/astarte_realm_management_api.yaml [limit=80, offset=470]
→ Read docs/api/astarte_realm_management_api.yaml [limit=100, offset=660]

$ rg -n "^build:|docs" /root/astrate-mule/Makefile | head -30
1:# Astrate build & verification entry points (docs/ROADMAP.md §1.2, verification tiers §0.2).
22:build:
53:		echo "conformance harness lands in M9 (docs/ROADMAP.md §10); nothing to run yet"; \

Now let me check the docs build setup and how the YAML is validated.
→ Read docs/api/astarte_realm_management_api.yaml [limit=45, offset=1]

$ rg -n "build|swagger|yaml|redoc|validate" /root/astrate-mule/Makefile | head -40
1:# Astrate build & verification entry points (docs/ROADMAP.md §1.2, verification tiers §0.2).
3:#   make build               static binary / compile check (CGO disabled)
17:.DEFAULT_GOAL := build
19:.PHONY: build lint test test-integration test-e2e test-conformance tools up down clean
21:## build: compile every package statically; emits $(DIST)/astrate once cmd/astrate lands (M8).
22:build:
23:	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' ./...
26:		CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(DIST)/astrate ./cmd/astrate; \
44:## test-e2e: T3 — component tests behind the e2e build tag (suites land from M5 onward).

$ ls /root/astrate-mule/docs /root/astrate-mule/docs/api
/root/astrate-mule/docs:
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
UPSTREAM-EXPERIMENTAL.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui

/root/astrate-mule/docs/api:
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

→ Read docs/Makefile
→ Read docs/mkdocs.yml

$ ls /root/astrate-mule/docs/swagger-ui; rg -n "astarte_realm_management|YAML|yml" /root/astrate-mule/docs/swagger-ui | head -20
index.html
/root/astrate-mule/docs/swagger-ui/index.html:46:      <option value="../api/astarte_realm_management_api.yaml">Realm Management API — Interfaces, triggers, policies, a_rma auth</option>

$ rg -rn "swagger|openapi|pyyaml|yaml.load" /root/astrate-mule/tools /root/astrate-mule/docs --type-not yaml 2>/dev/null | rg -iv "swagger.ui" | head -20
/root/astrate-mule/docs/embed.go://go:embed n-ui/*
/root/astrate-mule/docs/Makefile:	rm -rf site/n-ui && cp -r n-ui site/n-ui
/root/astrate-mule/docs/site/swagger.md:<iframe src="../n-ui/" style="width:100%;height:80vh;border:1px solid var(--md-default-fg-color--lightest);border-radius:4px;" loading="lazy"></iframe>
/root/astrate-mule/docs/handoff/phase-4-memory.md:- Copied DESIGN.md, ROADMAP.md, api/, n-ui/ into `docs/site/` so MkDocs can serve them.
/root/astrate-mule/docs/handoff/phase-4-memory.md:- Verified clean build: `mkdocs build` succeeds, all 15 pages + DESIGN + ROADMAP + api + n-ui present in output.
/root/astrate-mule/docs/handoff/phase-4-memory.md:- `docs/site/n-ui/` — copied from `docs/n-ui/`
/root/astrate-mule/docs/handoff/phase-4-memory.md:- All nav links, DESIGN, ROADMAP, api/, n-ui/ links resolve in built output
/root/astrate-mule/docs/handoff/phase-4-memory.md:- `site-dist/` contains: index.html, 15 page dirs, DESIGN/, ROADMAP/, api/, n-ui/, assets/, search/
/root/astrate-mule/docs/handoff/phase-6-memory.md:- Ran `make sync && make build` locally — build succeeded with only informational warnings (MkDocs 2.0 notice, unnav'd DESIGN.md/ROADMAP.md, relative links to api/ and n-ui/).
/root/astrate-mule/docs/handoff/phase-6-memory.md:- `docs/site/` — 19 entries including synced DESIGN.md, ROADMAP.md, api/, n-ui/
/root/astrate-mule/docs/handoff/phase-6-memory.md:- `site-dist/` contains: index.html, 15 page dirs, api/ (5 YAML), n-ui/ (index.html)
/root/astrate-mule/docs/handoff/phase-6-memory.md:- Relative links from `index.md` to `api/` and `n-ui/` are not resolved by MkDocs but work as static file paths in the built site.
/root/astrate-mule/docs/handoff/phase-3-handoff.md:- Phase 2.2 (go:embed + /n route) was not done.
/root/astrate-mule/docs/site/index.md:- [API Explorer](n.md) -- interactive API docs
/root/astrate-mule/docs/handoff/phase-4-handoff.md:- DESIGN.md, ROADMAP.md, api/, n-ui/ are copied into docs/site/ for MkDocs. Run make sync to refresh.
/root/astrate-mule/docs/handoff/phase-2-memory.md:#### Step 2.1: `docs/n-ui/index.html`
/root/astrate-mule/docs/handoff/phase-2-memory.md:#### Step 2.2: Embed nel binary + route `/n`
/root/astrate-mule/docs/handoff/phase-2-memory.md:   - `SwaggerUI` (via `//go:embed n-ui/*`)
/root/astrate-mule/docs/handoff/phase-2-memory.md:2. **`internal/n/n.go`** — package `n`, funzione `Mount(mux *http.ServeMux)`:
/root/astrate-mule/docs/handoff/phase-2-memory.md:   - `GET /n` → redirect 302 a `/n/index.html`

$ cd /root/astrate-mule && git log --oneline -8 -- docs/api/astarte_realm_management_api.yaml
7256af9 mule: docs-sync-rm-put-interface-409 [auto]: add the missing `409` Conflict response to `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` in docs/api/astarte_realm_management_api.yaml (yaml:222-234 lists only 204/400/401/404/422/500) — `updateInterface` answers 409 for every url/body disagreement and minor-upgrade incompatibility: `ErrNameMismatch` (service.go:204 → "Interface name doesn't match the one in the interface json", http.go:367-369), `ErrMajorMismatch` (service.go:207 → http.go:370-372), and the `CheckMinorUpgrade` sentinels `ErrMinorNotIncreased`/`ErrDowngradeNotAllowed`/`ErrMissingEndpoints`/`ErrIncompatibleEndpointChange` (service.go:223-226 → http.go:373-384) plus `store.ErrAlreadyExists` (http.go:385-386) — the exact set `TestRealmManagementErrorCodes` pins as StatusConflict (http_test.go:417-437). Reuse the `Conflict` response component (already used on POST, yaml:131-132) and add name-mismatch/major-mismatch examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
8bfedd6 mule: docs-sync-rm-triggers-422-nested-envelope [auto]: document the nested changeset envelope on `POST /realmmanagement/v1/{realm}/triggers` in docs/api/astarte_realm_management_api.yaml — when trigger compile fails, `createTrigger` answers 422 with the upstream nested body `{"errors":{"action":{...},"simple_triggers":[...]}}` (internal/realm/http.go:252-258, triggerErrorBody at http.go:268-285, issues #63/#70), but the spec's createTrigger 422 (yaml:344-345) refs only the flat `ValidationError` component (`{"errors":{"detail":...}}`) — which the handler also emits for plain ErrValidation (http.go:259-260, e.g. "trigger requires a name", service.go:429). Model the response as oneOf across the flat ErrorDetail and the nested envelope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
c0521bc mule: docs-sync-rm-policies-delete-422 [auto]: add the missing `422` response to `DELETE /realmmanagement/v1/{realm}/policies/{name}` in docs/api/astarte_realm_management_api.yaml — `Service.DeletePolicy` refuses to delete a policy still referenced by a trigger (`"policy %q is still used by trigger %q"`, internal/realm/service.go:573-574, wrapped as ErrValidation) and `writeError` answers 422 for it (internal/realm/http.go:410-411, detail stripped by validationDetail), but the spec lists only 204/401/404/500 for the delete op (yaml:695-713). Add the 422 ValidationError response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
8554720 mule: docs-sync-rm-put-auth-422 [auto]: add the missing `422` (ValidationError) response to `PUT /realmmanagement/v1/{realm}/config/auth` in docs/api/astarte_realm_management_api.yaml — `SetAuthKey` returns `ErrValidation` for a blank `jwt_public_key_pem` (internal/realm/service.go:632-640), which `writeError` maps to 422; the spec (yaml:431-464) documents only 204/400/401/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
5f8cb25 mule: docs-sync-rm-mapping-required-encrypted [auto]: add `required` (boolean, datastream-only) and `encrypted` (boolean, datastream-only) to the `InterfaceMapping` schema in docs/api/astarte_realm_management_api.yaml — the parser accepts and persists both on datastream mappings (internal/interfaceschema/parse.go:410-414), rejects them on properties mappings (parse.go:395-398), the detailed listing emits them (`detailed.go:69`), and the compatibility check compares them (interfaceschema/compat.go:121-123); the schema (yaml:820-862) currently has `additionalProperties: false` and neither field, so the OpenAPI would mark a valid datastream interface install body as invalid. Also add both fields to the detailed-listing example (yaml:70-85). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
b37d11f mule: docs-sync-rm-delete-interface-status [auto]: fix the `deleteInterface` response in docs/api/astarte_realm_management_api.yaml — the spec documents `422` with examples "Interface major version is not 0, can't be deleted" / "Cannot delete an interface that is used by a device introspection" (yaml:255-272), but the handler returns `403` with "Interface can't be deleted" / "Interface can't be deleted since it's currently used" (internal/realm/http.go:387-391, `writeInterfaceError`; both are 403 `StatusForbidden`). The spec's detail strings were lifted from the `writeError` branch (http.go:414-419) which DELETE /interfaces does not use. Fix status to `403` and update the two examples to match the actual wire messages. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
4094145 mule: docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
0574fd9 mule: docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` route to docs/api/astarte_realm_management_api.yaml — it exists in code (internal/realm/http.go:51, handler getDatastreamMaximumStorageRetention at http.go:135, returns 200 with a data-envelope retention value via Service.GetDatastreamMaximumStorageRetention) but the spec jumps straight from config/device_registration_limit (yaml:428) to /version. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

The previous commit `8bfedd6` did the createTrigger oneOf. Let me mirror its exact style.
$ cd /root/astrate-mule && git show 8bfedd6 -- docs/api/astarte_realm_management_api.yaml | head -120
commit 8bfedd697d4dfb879fe64257e9c4b2c739cbb5e5
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Thu Sep 17 22:21:37 2026 +0200

    mule: docs-sync-rm-triggers-422-nested-envelope [auto]: document the nested changeset envelope on `POST /realmmanagement/v1/{realm}/triggers` in docs/api/astarte_realm_management_api.yaml — when trigger compile fails, `createTrigger` answers 422 with the upstream nested body `{"errors":{"action":{...},"simple_triggers":[...]}}` (internal/realm/http.go:252-258, triggerErrorBody at http.go:268-285, issues #63/#70), but the spec's createTrigger 422 (yaml:344-345) refs only the flat `ValidationError` component (`{"errors":{"detail":...}}`) — which the handler also emits for plain ErrValidation (http.go:259-260, e.g. "trigger requires a name", service.go:429). Model the response as oneOf across the flat ErrorDetail and the nested envelope. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

diff --git a/docs/api/astarte_realm_management_api.yaml b/docs/api/astarte_realm_management_api.yaml
index b0873e0..57d9981 100644
--- a/docs/api/astarte_realm_management_api.yaml
+++ b/docs/api/astarte_realm_management_api.yaml
@@ -342,7 +342,30 @@ paths:
         "409":
           $ref: "#/components/responses/Conflict"
         "422":
-          $ref: "#/components/responses/ValidationError"
+          description: Unprocessable Entity — trigger definition failed validation.
+          content:
+            application/json; charset=utf-8:
+              schema:
+                oneOf:
+                  - $ref: "#/components/schemas/ErrorDetail"
+                  - $ref: "#/components/schemas/TriggerErrorsDetail"
+              examples:
+                flat_detail:
+                  summary: Plain validation error
+                  value:
+                    errors:
+                      detail: "trigger requires a name"
+                nested_changeset:
+                  summary: Trigger compile failed — upstream nested changeset envelope
+                  value:
+                    errors:
+                      action:
+                        http_url:
+                          - "invalid url"
+                      simple_triggers:
+                        - type:
+                            - "bad type"
+                        - {}
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -1051,6 +1074,31 @@ components:
             detail:
               type: string
 
+    TriggerErrorsDetail:
+      type: object
+      required: [errors]
+      properties:
+        errors:
+          type: object
+          description: Upstream nested changeset envelope for a failed trigger compile.
+          properties:
+            action:
+              type: object
+              description: Per-field error messages on the trigger action (changeset shape).
+              additionalProperties:
+                type: array
+                items:
+                  type: string
+            simple_triggers:
+              type: array
+              description: Index-aligned per-condition errors; empty objects mark entries without errors.
+              items:
+                type: object
+                additionalProperties:
+                  type: array
+                  items:
+                    type: string
+
   responses:
     BadRequest:
       description: Bad Request — malformed or missing data envelope.

$ cd /root/astrate-mule && rg -n "StatusUnprocessableEntity|422|detail" internal/realm/http_test.go | rg -i "detail|422|invalid|retention_band|ttl" | head -30
164:		// answers with its measured 409 detail (#62).
165:		want := `{"errors":{"detail":"Interface update contains incompatible endpoint changes"}}`
186:		want := `{"errors":{"detail":"Interface can't be deleted"}}`
203:		wantInUse := `{"errors":{"detail":"Interface can't be deleted since it's currently used"}}`
225:			t.Errorf("invalid trigger: got %d, want 422", rec.Code)
270:// datastream_maximum_storage_retention answer upstream's named 422 envelope
301:	// Install above the ceiling → upstream's named 422, nothing installed.
304:		t.Errorf("install ttl 7200 under a 3600 ceiling: got %d %s, want 422 %s",
320:		t.Errorf("minor update to ttl 7200: got %d %s, want 422 %s",
332:// compares the exact {"errors":{"detail":...}} body, interleaved with success
338:		bodyDup           = `{"errors":{"detail":"Interface already exists"}}`
339:		bodyCollision     = `{"errors":{"detail":"Interface name collision detected. Make sure that the difference between two interface names is not limited to the casing or the presence of hyphens."}}`
340:		bodyNameMismatch  = `{"errors":{"detail":"Interface name doesn't match the one in the interface json"}}`
341:		bodyMajorMismatch = `{"errors":{"detail":"Interface major version doesn't match the one in the interface json"}}`
342:		bodyMinorSame     = `{"errors":{"detail":"Interface minor version was not increased"}}`
343:		bodyDowngrade     = `{"errors":{"detail":"Interface downgrade not allowed"}}`
344:		bodyMutated       = `{"errors":{"detail":"Interface update contains incompatible endpoint changes"}}`
345:		bodyDropped       = `{"errors":{"detail":"Interface update has missing endpoints"}}`
346:		bodyMajorNotFound = `{"errors":{"detail":"Interface major not found"}}`
347:		bodyNotFound      = `{"errors":{"detail":"Interface not found"}}`
348:		bodyCantDelete    = `{"errors":{"detail":"Interface can't be deleted"}}`
349:		bodyInUse         = `{"errors":{"detail":"Interface can't be deleted since it's currently used"}}`
350:		bodyTriggerDup    = `{"errors":{"detail":"Already exists"}}`
454:		{"duplicate_trigger_keeps_generic_detail", http.MethodPost, "/triggers", triggerJSON,
474:// TestRealmManagementViolationEnvelopes pins the probe-frozen 422 bodies of
487:			t.Errorf("got %d %s, want 422 %s", rec.Code, rec.Body, want)
499:			t.Errorf("got %d %s, want 422 %s", rec.Code, rec.Body, want)
549:		t.Errorf("ttl 59: got %d %s, want 422 %s", rec.Code, rec.Body, want)
557:// an amqp-exchange trigger action is rejected at creation with 422 (upstream
560:// overlap/prefetch violations answer 422 with the probed upstream wording

$ cd /root/astrate-mule && rg -n "ErrValidation" internal/realm/*.go internal/interfaceschema/*.go | head -30
rg: internal/interfaceschema/*.go: No such file or directory (os error 2)
internal/realm/policy.go:10:// triggers.CompilePolicy. Violations wrap ErrValidation (mapped to 422).
internal/realm/policy.go:14:		return "", fmt.Errorf("%w: %v", ErrValidation, err)
internal/realm/policy_test.go:57:				if !errors.Is(err, ErrValidation) {
internal/realm/policy_test.go:58:					t.Errorf("error %v does not wrap ErrValidation", err)
internal/realm/service.go:26:// ErrValidation wraps a request that is well-formed JSON but violates an
internal/realm/service.go:28:// and fmt.Errorf("%w: ...", ErrValidation, ...) to attach the detail.
internal/realm/service.go:29:var ErrValidation = errors.New("realm: validation failed")
internal/realm/service.go:151:// store.ErrAlreadyExists; a schema violation yields ErrValidation — carrying
internal/realm/service.go:163:		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
internal/realm/service.go:198:		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
internal/realm/service.go:225:		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
internal/realm/service.go:426:		return nil, fmt.Errorf("%w: trigger does not parse: %v", ErrValidation, err)
internal/realm/service.go:429:		return nil, fmt.Errorf("%w: trigger requires a name", ErrValidation)
internal/realm/service.go:436:		// generic ErrValidation.
internal/realm/service.go:441:		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
internal/realm/service.go:446:				return nil, fmt.Errorf("%w: policy %q does not exist in realm %q", ErrValidation, ct.PolicyName, realm)
internal/realm/service.go:574:			return fmt.Errorf("%w: policy %q is still used by trigger %q", ErrValidation, name, trs[i].Name)
internal/realm/service.go:634:		return fmt.Errorf("%w: jwt_public_key_pem can't be blank", ErrValidation)
internal/realm/http.go:392:	case errors.Is(err, ErrValidation):
internal/realm/http.go:410:	case errors.Is(err, ErrValidation):
internal/realm/http.go:560:// validationDetail strips the ErrValidation sentinel prefix so the response

→ Read internal/realm/service.go [limit=85, offset=145]

$ cd /root/astrate-mule && rg -n "ViolationsError|ErrValidation|fmt.Errorf|errors.New" internal/interfaceschema/parse.go | head -40
rg: internal/interfaceschema/parse.go: IO error for operation on internal/interfaceschema/parse.go: No such file or directory (os error 2)

$ cd /root/astrate-mule && ls internal/interfaceschema/ && rg -n "func ParseInterfaceCanonical|errSchemaValidation|validation failed|errors.New|fmt.Errorf" internal/interfaceschema/*.go | rg -i "validation|canonical|violation"| head -30
ls: cannot access 'internal/interfaceschema/': No such file or directory

$ cd /root/astrate-mule && rg -n "violations" internal/realm/service.go | head; find . -type d -name "*schema*" 2>/dev/null; rg -l "ParseInterfaceCanonical" internal --type go
./pkg/interfaceschema
./site-dist/interface-schema
internal/realm/service.go

$ cd /root/astrate-mule && rg -n "func ParseInterfaceCanonical|func ParseInterface\b" pkg/interfaceschema/ -g "*.go"
pkg/interfaceschema/parse.go:160:func ParseInterface(data []byte) (*Interface, error) {
pkg/interfaceschema/parse.go:171:func ParseInterfaceCanonical(data []byte) (*Interface, []byte, error) {

→ Read pkg/interfaceschema/parse.go [limit=65, offset=155]
→ Read internal/realm/http.go [limit=120, offset=120]

$ cd /root/astrate-mule && rg -n "ErrInvalid|\"interface definition is invalid\"|\"interface name\"|errors.New\(" pkg/interfaceschema/parse.go | head -40
36:// ErrInvalid is wrapped by every ParseInterface failure, so callers can
38:var ErrInvalid = errors.New("invalid interface")
157:// wrapping ErrInvalid. Rejections for the rules whose upstream wire shape was
176:		return nil, nil, fmt.Errorf("%w: %w", ErrInvalid, err)
179:		return nil, nil, fmt.Errorf("%w: trailing data after interface document", ErrInvalid)
186:			return nil, nil, err // the combined value already wraps ErrInvalid
188:		return nil, nil, fmt.Errorf("%w: %w", ErrInvalid, err)
197:		return nil, nil, fmt.Errorf("%w: %w", ErrInvalid, err)
206:// parse as a combined ErrInvalid + *ViolationsError value once every mapping
240:		return nil, errors.New(`missing "interface_name"`)
242:		return nil, errors.New(`missing "version_major"`)
244:		return nil, errors.New(`missing "version_minor"`)
246:		return nil, errors.New(`missing "type"`)
248:		return nil, errors.New(`missing "ownership"`)
261:		return nil, errors.New("version_major must be >= 0")
264:		return nil, errors.New("version_minor must be >= 0")
267:		return nil, errors.New("version_major and version_minor cannot both be 0")
284:		return nil, errors.New("properties interfaces cannot use object aggregation")
288:		return nil, errors.New("interface must declare at least one mapping")
354:		return m, nil, errors.New(`mapping is missing "endpoint"`)

Now I have full context. I'll make the edits — first the two 422 blocks, then the schemas. There's a local docs `site/` copy too; I'll note that. Let me start with the POST interfaces block.
← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -130,9 +130,38 @@
         $ref: "#/components/responses/Unauthorized"
       "409":
         $ref: "#/components/responses/Conflict"
       "422":
-        $ref: "#/components/responses/ValidationError"
+        description: Unprocessable Entity — interface definition failed validation.
+        content:
+          application/json; charset=utf-8:
+            schema:
+              oneOf:
+                - $ref: "#/components/schemas/ErrorDetail"
+                - $ref: "#/components/schemas/InterfaceViolationsDetail"
+                - $ref: "#/components/schemas/InterfaceFieldErrorsDetail"
+            examples:
+              flat_detail:
+                summary: Plain validation error
+                value:
+                  errors:
+                    detail: 'missing "interface_name"'
+              nested_changeset:
+                summary: Schema violations — upstream changeset envelope with the aligned mappings array
+                value:
+                  errors:
+                    description:
+                      - "should be at most 1000 character(s)"
+                    mappings:
+                      - {}
+                      - description:
+                          - "should be at most 1000 character(s)"
+              retention_exceeded:
+                summary: Maximum database storage retention exceeded — named field error
+                value:
+                  errors:
+                    error_name:
+                      - "maximum_database_retention_exceeded"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
 /realmmanagement/v1/{realm}/interfaces/{name}:

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -259,9 +259,38 @@
       $ref: "#/components/responses/NotFound"
     "409":
       $ref: "#/components/responses/Conflict"
     "422":
-      $ref: "#/components/responses/ValidationError"
+      description: Unprocessable Entity — interface definition failed validation.
+      content:
+        application/json; charset=utf-8:
+          schema:
+            oneOf:
+              - $ref: "#/components/schemas/ErrorDetail"
+              - $ref: "#/components/schemas/InterfaceViolationsDetail"
+              - $ref: "#/components/schemas/InterfaceFieldErrorsDetail"
+          examples:
+            flat_detail:
+              summary: Plain validation error
+              value:
+                errors:
+                  detail: 'missing "interface_name"'
+            nested_changeset:
+              summary: Schema violations — upstream changeset envelope with the aligned mappings array
+              value:
+                errors:
+                  description:
+                    - "should be at most 1000 character(s)"
+                  mappings:
+                    - {}
+                    - description:
+                        - "should be at most 1000 character(s)"
+            retention_exceeded:
+              summary: Maximum database storage retention exceeded — named field error
+              value:
+                errors:
+                  error_name:
+                    - "maximum_database_retention_exceeded"
     "500":
       $ref: "#/components/responses/InternalServerError"
 
 delete:

Now the schemas. Let me add them after `TriggerErrorsDetail`.
← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -1158,8 +1158,43 @@
                 type: array
                 items:
                   type: string
 
+  InterfaceViolationsDetail:
+    type: object
+    required: [errors]
+    properties:
+      errors:
+        type: object
+        description: Interface schema violations in upstream's changeset shape.
+        properties:
+          mappings:
+            type: array
+            description: Index-aligned per-mapping errors; empty objects mark entries without errors.
+            items:
+              type: object
+              additionalProperties:
+                type: array
+                items:
+                  type: string
+        additionalProperties:
+          type: array
+          items:
+            type: string
+
+  InterfaceFieldErrorsDetail:
+    type: object
+    required: [errors]
+    properties:
+      errors:
+        type: object
+        required: [error_name]
+        properties:
+          error_name:
+            type: array
+            items:
+              type: string
+
 responses:
   BadRequest:
     description: Bad Request — malformed or missing data envelope.
     content:

Now let me verify the edits and run the gate.
$ cd /root/astrate-mule && sed -n '125,165p' docs/api/astarte_realm_management_api.yaml
description: Installed interface definition.
                    $ref: "#/components/schemas/InterfaceDefinition"
        "400":
          $ref: "#/components/responses/BadRequest"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "409":
          $ref: "#/components/responses/Conflict"
        "422":
          description: Unprocessable Entity — interface definition failed validation.
          content:
            application/json; charset=utf-8:
              schema:
                oneOf:
                  - $ref: "#/components/schemas/ErrorDetail"
                  - $ref: "#/components/schemas/InterfaceViolationsDetail"
                  - $ref: "#/components/schemas/InterfaceFieldErrorsDetail"
              examples:
                flat_detail:
                  summary: Plain validation error
                  value:
                    errors:
                      detail: 'missing "interface_name"'
                nested_changeset:
                  summary: Schema violations — upstream changeset envelope with the aligned mappings array
                  value:
                    errors:
                      description:
                        - "should be at most 1000 character(s)"
                      mappings:
                        - {}
                        - description:
                            - "should be at most 1000 character(s)"
                retention_exceeded:
                  summary: Maximum database storage retention exceeded — named field error
                  value:
                    errors:
                      error_name:
                        - "maximum_database_retention_exceeded"
        "500":
          $ref: "#/components/responses/InternalServerError"

$ cd /root/astrate-mule && sed -n '228,270p' docs/api/astarte_realm_management_api.yaml
operationId: updateInterface
      summary: Update an interface definition
      description: |
        Updates an existing interface definition (same major version). Returns
        204 on success with no body.
      tags: [Realm Management - Interfaces]
      security:
        - a_rma: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - $ref: "#/components/parameters/InterfaceName"
        - $ref: "#/components/parameters/MajorVersion"
      requestBody:
        required: true
        content:
          application/json; charset=utf-8:
            schema:
              type: object
              required: [data]
              properties:
                data:
                  description: Updated interface definition.
                  type: object
      responses:
        "204":
          description: Interface updated. No body returned.
        "400":
          $ref: "#/components/responses/BadRequest"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "404":
          $ref: "#/components/responses/NotFound"
        "409":
          $ref: "#/components/responses/Conflict"
        "422":
          description: Unprocessable Entity — interface definition failed validation.
          content:
            application/json; charset=utf-8:
              schema:
                oneOf:
                  - $ref: "#/components/schemas/ErrorDetail"
                  - $ref: "#/components/schemas/InterfaceViolationsDetail"
                  - $ref: "#/components/schemas/InterfaceFieldErrorsDetail"

$ cd /root/astrate-mule && sed -n '1150,1215p' docs/api/astarte_realm_management_api.yaml
items:
                  type: string
            simple_triggers:
              type: array
              description: Index-aligned per-condition errors; empty objects mark entries without errors.
              items:
                type: object
                additionalProperties:
                  type: array
                  items:
                    type: string

    InterfaceViolationsDetail:
      type: object
      required: [errors]
      properties:
        errors:
          type: object
          description: Interface schema violations in upstream's changeset shape.
          properties:
            mappings:
              type: array
              description: Index-aligned per-mapping errors; empty objects mark entries without errors.
              items:
                type: object
                additionalProperties:
                  type: array
                  items:
                    type: string
          additionalProperties:
            type: array
            items:
              type: string

    InterfaceFieldErrorsDetail:
      type: object
      required: [errors]
      properties:
        errors:
          type: object
          required: [error_name]
          properties:
            error_name:
              type: array
              items:
                type: string

  responses:
    BadRequest:
      description: Bad Request — malformed or missing data envelope.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Bad Request

    Unauthorized:
      description: Unauthorized — missing or invalid JWT.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:

Now the gate. Let me check available tooling and run the docs build.
$ cd /root/astrate-mule && which mkdocs python3 2>/dev/null; ls docs/site/api 2>/dev/null | head; git status --short
/usr/bin/python3
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
 M docs/api/astarte_realm_management_api.yaml

$ cd /root/astrate-mule && git ls-files docs/site | head -3; echo "---"; git check-ignore docs/site/api/astarte_realm_management_api.yaml; echo "ignore-exit:$?"; ls docs/site/swagger-ui 2>/dev/null | head
docs/site/appengine-api.md
docs/site/architecture.md
docs/site/compatibility.md
---
docs/site/api/astarte_realm_management_api.yaml
ignore-exit:0
index.html

→ Read docs/requirements.txt

$ cd /root/astrate-mule && python3 -c "import yaml; print('pyyaml', yaml.__version__)" 2>&1; python3 -m pip --version 2>&1 | head -1
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; print('pyyaml', yaml.__version__)
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'
/usr/bin/python3: No module named pip

$ cd /root/astrate-mule && which go; go version 2>&1; ls tools 2>/dev/null; ls website 2>/dev/null
/usr/bin/go
go version go1.26.7 linux/arm64
bsoncapture
mule-plan-day.sh
mule-setup-pi.sh
mule-survey.sh
mule.sh
reconcile.sh

$ cd /root/astrate-mule && rg -l "sigs.k8s.io/yaml|ghodss/yaml|swagger" --glob '*.go' | head; cat go.mod | rg "yaml|openapi|swagger"
cmd/astrate/main.go
docs/embed.go
internal/swagger/swagger.go
internal/swagger/swagger_test.go
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect

→ Read internal/swagger/swagger.go
→ Read internal/swagger/swagger_test.go

$ cd /root/astrate-mule && make -C docs build 2>&1 | tail -20
make: Entering directory '/root/astrate-mule/docs'
cp DESIGN.md site/DESIGN.md
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'

`sync` succeeds; `mkdocs` itself is missing on this machine (no pip/python-yaml either). I'll verify with a Go probe instead, per MULE.md's live-probe rule. Let me check what's in the module cache.
$ ls "$(go env GOMODCACHE)/gopkg.in" 2>/dev/null; ls "$(go env GOMODCACHE)/go.yaml.in" 2>/dev/null; ls "$(go env GOMODCACHE)/cache/download/gopkg.in" 2>/dev/null
check.v1@v0.0.0-20161208181325-20d25e280405
check.v1@v1.0.0-20201130134442-10cb98267c6c
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1
yaml
check.v1
inf.v0
ini.v1
yaml.v2
yaml.v3

yaml.v3 is in the module cache — I can build an offline probe. Let me write a throwaway checker in `/tmp` per the rules.
← Write /tmp/opencode/yamlcheck/main.go
Wrote file successfully.

$ cd /tmp/opencode/yamlcheck && printf 'module yamlcheck\n\ngo 1.26\n\nrequire gopkg.in/yaml.v3 v3.0.1\n' > go.mod && GOFLAGS=-mod=mod GOPROXY=off go build -o yamlcheck . && echo BUILD_OK
BUILD_OK

$ /tmp/opencode/yamlcheck/yamlcheck /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
yaml parses ok
all 130 internal $refs resolve ok
/realmmanagement/v1/{realm}/interfaces:post 422 oneOf ok: [#/components/schemas/ErrorDetail #/components/schemas/InterfaceViolationsDetail #/components/schemas/InterfaceFieldErrorsDetail]
/realmmanagement/v1/{realm}/interfaces/{name}/{major}:put 422 oneOf ok: [#/components/schemas/ErrorDetail #/components/schemas/InterfaceViolationsDetail #/components/schemas/InterfaceFieldErrorsDetail]
new schemas ok

Probe passes: YAML loads, all 130 internal `$ref`s resolve, both operations carry the 422 oneOf. Let me check for an alternate mkdocs install and then review the final diff.
$ ls ~/.local/bin/mkdocs /usr/local/bin/mkdocs 2>/dev/null; find / -name "mkdocs" -type f 2>/dev/null | head -3; echo done
/tmp/opencode/mkdocs-venv/bin/mkdocs
done

A mkdocs venv already exists at `/tmp/opencode/mkdocs-venv`. Let me use it to run the actual build.
$ cd /root/astrate-mule && PATH="/tmp/opencode/mkdocs-venv/bin:$PATH" make -C docs build 2>&1 | tail -30
cp DESIGN.md site/DESIGN.md
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml

 │  ⚠  Warning from the Material for MkDocs team
 │
 │  MkDocs 2.0, the underlying framework of Material for MkDocs,
 │  will introduce backward-incompatible changes, including:
 │
 │  × All plugins will stop working – the plugin system has been removed
 │  × All theme overrides will break – the theming system has been rewritten
 │  × No migration path exists – existing projects cannot be upgraded
 │  × Closed contribution model – community members can't report bugs
 │  × Currently unlicensed – unsuitable for production use
 │
 │  Our full analysis:
 │
 │  https://squidfunk.github.io/mkdocs-material/blog/2026/02/18/mkdocs-2.0/

INFO    -  Cleaning site directory
INFO    -  Building documentation to directory: /root/astrate-mule/site-dist
INFO    -  The following pages exist in the docs directory, but are not included in the "nav" configuration:
  - DESIGN.md
  - ROADMAP.md
INFO    -  Doc file 'index.md' contains an unrecognized relative link 'api/', it was left as is.
INFO    -  Documentation built in 8.97 seconds
make: Leaving directory '/root/astrate-mule/docs'

$ cd /root/astrate-mule && git status --short && git diff --stat
M docs/api/astarte_realm_management_api.yaml
 docs/api/astarte_realm_management_api.yaml | 97 +++++++++++++++++++++++++++++-
 1 file changed, 95 insertions(+), 2 deletions(-)

Build succeeds. Final review of the tracked diff.
$ cd /root/astrate-mule && git diff docs/api/astarte_realm_management_api.yaml
diff --git a/docs/api/astarte_realm_management_api.yaml b/docs/api/astarte_realm_management_api.yaml
index ba1a8f3..4a2f22a 100644
--- a/docs/api/astarte_realm_management_api.yaml
+++ b/docs/api/astarte_realm_management_api.yaml
@@ -131,7 +131,36 @@ paths:
         "409":
           $ref: "#/components/responses/Conflict"
         "422":
-          $ref: "#/components/responses/ValidationError"
+          description: Unprocessable Entity — interface definition failed validation.
+          content:
+            application/json; charset=utf-8:
+              schema:
+                oneOf:
+                  - $ref: "#/components/schemas/ErrorDetail"
+                  - $ref: "#/components/schemas/InterfaceViolationsDetail"
+                  - $ref: "#/components/schemas/InterfaceFieldErrorsDetail"
+              examples:
+                flat_detail:
+                  summary: Plain validation error
+                  value:
+                    errors:
+                      detail: 'missing "interface_name"'
+                nested_changeset:
+                  summary: Schema violations — upstream changeset envelope with the aligned mappings array
+                  value:
+                    errors:
+                      description:
+                        - "should be at most 1000 character(s)"
+                      mappings:
+                        - {}
+                        - description:
+                            - "should be at most 1000 character(s)"
+                retention_exceeded:
+                  summary: Maximum database storage retention exceeded — named field error
+                  value:
+                    errors:
+                      error_name:
+                        - "maximum_database_retention_exceeded"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -231,7 +260,36 @@ paths:
         "409":
           $ref: "#/components/responses/Conflict"
         "422":
-          $ref: "#/components/responses/ValidationError"
+          description: Unprocessable Entity — interface definition failed validation.
+          content:
+            application/json; charset=utf-8:
+              schema:
+                oneOf:
+                  - $ref: "#/components/schemas/ErrorDetail"
+                  - $ref: "#/components/schemas/InterfaceViolationsDetail"
+                  - $ref: "#/components/schemas/InterfaceFieldErrorsDetail"
+              examples:
+                flat_detail:
+                  summary: Plain validation error
+                  value:
+                    errors:
+                      detail: 'missing "interface_name"'
+                nested_changeset:
+                  summary: Schema violations — upstream changeset envelope with the aligned mappings array
+                  value:
+                    errors:
+                      description:
+                        - "should be at most 1000 character(s)"
+                      mappings:
+                        - {}
+                        - description:
+                            - "should be at most 1000 character(s)"
+                retention_exceeded:
+                  summary: Maximum database storage retention exceeded — named field error
+                  value:
+                    errors:
+                      error_name:
+                        - "maximum_database_retention_exceeded"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -1101,6 +1159,41 @@ components:
                   items:
                     type: string
 
+    InterfaceViolationsDetail:
+      type: object
+      required: [errors]
+      properties:
+        errors:
+          type: object
+          description: Interface schema violations in upstream's changeset shape.
+          properties:
+            mappings:
+              type: array
+              description: Index-aligned per-mapping errors; empty objects mark entries without errors.
+              items:
+                type: object
+                additionalProperties:
+                  type: array
+                  items:
+                    type: string
+          additionalProperties:
+            type: array
+            items:
+              type: string
+
+    InterfaceFieldErrorsDetail:
+      type: object
+      required: [errors]
+      properties:
+        errors:
+          type: object
+          required: [error_name]
+          properties:
+            error_name:
+              type: array
+              items:
+                type: string
+
   responses:
     BadRequest:
       description: Bad Request — malformed or missing data envelope.

$ cd /root/astrate-mule && gofmt -l . ; echo "gofmt-exit:$?"; sed -n '44,48p' docs/swagger-ui/index.html
gofmt-exit:0
      <option value="../api/astarte_housekeeping_api.yaml">Housekeeping API — Realm management, a_ha auth</option>
      <option value="../api/astarte_pairing_api.yaml">Pairing API — Device pairing flows A/B/C, a_pa auth</option>
      <option value="../api/astarte_realm_management_api.yaml">Realm Management API — Interfaces, triggers, policies, a_rma auth</option>
      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
      <option value="../api/astrate_native_api.yaml">Astrate Native API — Health, readiness, metrics, WebSocket</option>

Task done. Report:

```
Done: Modeled the interface POST/PUT 422 in specs as a oneOf across the flat ErrorDetail, the InterfaceViolationsDetail (nested changeset + aligned mappings array), and the InterfaceFieldErrorsDetail (named retention envelope), mirroring the createTrigger pattern.
Files: docs/api/astarte_realm_management_api.yaml
Verified: make -C docs build (with PATH=/tmp/opencode/mkdocs-venv/bin:$PATH) -> pass; /tmp OpenAPI probe (yaml.v3 parse, all 130 internal $refs resolve, both ops carry the 422 oneOf) -> pass; gofmt -l . -> pass. No Go code changed, so no go test run.
Unsure: The flat example uses the wire-accurate detail 'missing "interface_name"' (parse.go:240 verbatim, prefix stripped by validationDetail); the other two examples are verbatim from http_test.go:282/484/496. Nothing else guessed.
Follow-ups: none
```
