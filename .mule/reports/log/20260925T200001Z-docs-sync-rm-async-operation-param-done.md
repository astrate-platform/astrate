slug: docs-sync-rm-async-operation-param
verdict: done
at:  3f0eca8
ran: 2026-09-25T20:00:01Z on DietPi in 425s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ cd /root/astrate-mule && rg -n "TestHousekeepingAsyncOperationParamDocumented" -A 60 internal/swagger/swagger_test.go
172:// TestHousekeepingAsyncOperationParamDocumented guards that the housekeeping
173-// spec tells clients the `?async_operation` parameter exists. Astrate accepts
174-// and ignores it on realm create and delete (deviation 17,
175-// docs/COMPATIBILITY.md), so an upstream client keeps working — but a client
176-// generated from the spec has to learn the parameter from the spec. This is
177-// the documentation half of the behaviour pinned by
178-// TestHousekeepingAsyncOperationParam in internal/housekeeping.
179:func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
180-	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
181-	if err != nil {
182-		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
183-	}
184-	lines := strings.Split(string(b), "\n")
185-
186-	const ref = `        - $ref: "#/components/parameters/AsyncOperation"`
187-	for _, op := range []string{"createRealm", "deleteRealm"} {
188-		block := operationBlock(t, lines, op)
189-		if !containsLine(block, ref) {
190-			t.Errorf("operation %s does not $ref the AsyncOperation parameter", op)
191-		}
192-	}
193-
194-	comp := componentBlock(t, lines, "    AsyncOperation:")
195-	for _, want := range []string{
196-		"      name: async_operation",
197-		"      in: query",
198-		"      required: false",
199-		"      schema:",
200-		"        type: boolean",
201-		"        default: false",
202-	} {
203-		if !containsLine(comp, want) {
204-			t.Errorf("components parameter AsyncOperation is missing line %q", want)
205-		}
206-	}
207-	if !strings.Contains(strings.Join(comp, "\n"), "Accepted and ignored") {
208-		t.Error("components parameter AsyncOperation does not say the value is accepted and ignored")
209-	}
210-}
211-
212-// propertyBlock returns the lines of the schema property with the given name
213-// inside a components.schemas entry, up to the next property or the end of the
214-// entry.
215-func propertyBlock(t *testing.T, lines []string, name string) []string {
216-	t.Helper()
217-	key := "        " + name + ":"
218-	start := -1
219-	for i, l := range lines {
220-		if l == key {
221-			start = i
222-			break
223-		}
224-	}
225-	if start < 0 {
226-		t.Fatalf("schema declares no property %q", name)
227-	}
228-	for i := start + 1; i < len(lines); i++ {
229-		if strings.TrimSpace(lines[i]) == "" {
230-			continue
231-		}
232-		if !strings.HasPrefix(lines[i], "          ") {
233-			return lines[start+1 : i]
234-		}
235-	}
236-	return lines[start+1:]
237-}
238-
239-// TestHousekeepingRetentionZeroFoldDocumented guards that the housekeeping spec

$ cd /root/astrate-mule && ls && rg -n "async_operation" --type go -l | head -30
AGENTS.md
Dockerfile
LICENSE
Makefile
README.md
astrate
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
internal/swagger/swagger_test.go
internal/housekeeping/async_operation_test.go
internal/realm/async_operation_test.go

$ cd /root/astrate-mule && rg -n "async_operation|AsyncOperation" internal/realm/async_operation_test.go | head -40
11:// TestRealmManagementAsyncOperationParam pins deviation 17
14:// `?async_operation=false`. Astrate is always synchronous, so the parameter is
17:func TestRealmManagementAsyncOperationParam(t *testing.T) {
20:	for i, q := range []string{"?async_operation=false", "?async_operation=true"} {

$ cd /root/astrate-mule && rg -n "async_operation|asyncOperation|AsyncOperation" --type go -g '!*_test.go' | head -40
(no output)

$ cd /root/astrate-mule && ls internal/realm/
async_operation_test.go
dashboard_compat_test.go
detailed.go
detailed_required_encrypted_test.go
detailed_test.go
helpers_test.go
http.go
http_test.go
policy.go
policy_test.go
service.go

→ Read internal/realm/async_operation_test.go

$ cd /root/astrate-mule && rg -n "interfaces|policies|HandleFunc|mux\.|route" internal/realm/http.go | head -60
14:	"github.com/astrate-platform/astrate/pkg/interfaceschema"
22:// 7.2). Every route is guarded by a realm JWT carrying a_rma.
34:// Mount registers the routes on mux (paths wire-frozen to upstream
38:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces", h(a.listInterfaces))
39:	mux.Handle("POST /realmmanagement/v1/{realm}/interfaces", h(a.installInterface))
40:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}", h(a.listInterfaceMajors))
41:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.getInterface))
42:	mux.Handle("PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.updateInterface))
43:	mux.Handle("DELETE /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.deleteInterface))
44:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers", h(a.listTriggers))
45:	mux.Handle("POST /realmmanagement/v1/{realm}/triggers", h(a.createTrigger))
46:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers/{name}", h(a.getTrigger))
47:	mux.Handle("DELETE /realmmanagement/v1/{realm}/triggers/{name}", h(a.deleteTrigger))
48:	mux.Handle("GET /realmmanagement/v1/{realm}/config/auth", h(a.getAuth))
49:	mux.Handle("PUT /realmmanagement/v1/{realm}/config/auth", h(a.putAuth))
50:	mux.Handle("GET /realmmanagement/v1/{realm}/config/device_registration_limit", h(a.getRegistrationLimit))
51:	mux.Handle("GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention", h(a.getDatastreamMaximumStorageRetention))
52:	mux.Handle("GET /realmmanagement/v1/{realm}/version", h(a.getVersion))
53:	mux.Handle("DELETE /realmmanagement/v1/{realm}/devices/{device}", h(a.deleteDevice))
54:	mux.Handle("GET /realmmanagement/v1/{realm}/policies", h(a.listPolicies))
55:	mux.Handle("POST /realmmanagement/v1/{realm}/policies", h(a.createPolicy))
56:	mux.Handle("GET /realmmanagement/v1/{realm}/policies/{name}", h(a.getPolicy))
57:	mux.Handle("DELETE /realmmanagement/v1/{realm}/policies/{name}", h(a.deletePolicy))
60:// --- trigger delivery policies ------------------------------------------------
144:// --- interfaces -------------------------------------------------------------
350:// route here on purpose — triggers and policies keep the generic writeError,
353:	var ve *interfaceschema.ViolationsError
373:	case errors.Is(err, interfaceschema.ErrMinorNotIncreased):
376:	case errors.Is(err, interfaceschema.ErrDowngradeNotAllowed):
379:	case errors.Is(err, interfaceschema.ErrMissingEndpoints):
382:	case errors.Is(err, interfaceschema.ErrIncompatibleEndpointChange):
403:	var ve *interfaceschema.ViolationsError
434:func writeViolations(w http.ResponseWriter, ve *interfaceschema.ViolationsError) {
443:func renderViolationsBody(ve *interfaceschema.ViolationsError) []byte {

$ cd /root/astrate-mule && rg -n "URL.Query|RawQuery" internal/realm/*.go | head -20
internal/realm/http.go:150:	if r.URL.Query().Get("detailed") == "true" {

$ cd /root/astrate-mule && rg -n "AsyncOperation" -B 4 -A 22 docs/api/astarte_housekeeping_api.yaml
57-      tags: [Housekeeping]
58-      security:
59-        - a_ha: []
60-      parameters:
61:        - $ref: "#/components/parameters/AsyncOperation"
62-      requestBody:
63-        required: true
64-        content:
65-          application/json; charset=utf-8:
66-            schema:
67-              type: object
68-              required: [data]
69-              properties:
70-                data:
71-                  $ref: "#/components/schemas/RealmCreate"
72-            example:
73-              data:
74-                realm_name: test
75-                jwt_public_key_pem: |
76-                  -----BEGIN PUBLIC KEY-----
77-                  MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
78-                  -----END PUBLIC KEY-----
79-                device_registration_limit: 100
80-      responses:
81-        "201":
82-          description: Realm created successfully.
83-          content:
--
199-      security:
200-        - a_ha: []
201-      parameters:
202-        - $ref: "#/components/parameters/RealmName"
203:        - $ref: "#/components/parameters/AsyncOperation"
204-      responses:
205-        "204":
206-          description: Realm deleted successfully. No body returned.
207-        "401":
208-          $ref: "#/components/responses/Unauthorized"
209-        "403":
210-          $ref: "#/components/responses/Forbidden"
211-        "404":
212-          $ref: "#/components/responses/NotFound"
213-        "405":
214-          $ref: "#/components/responses/MethodNotAllowed"
215-        "422":
216-          $ref: "#/components/responses/UnprocessableDeletion"
217-        "500":
218-          $ref: "#/components/responses/InternalServerError"
219-
220-components:
221-  securitySchemes:
222-    a_ha:
223-      type: apiKey
224-      in: header
225-      name: Authorization
--
235-      description: The realm name.
236-      schema:
237-        type: string
238-
239:    AsyncOperation:
240-      name: async_operation
241-      in: query
242-      required: false
243-      description: |
244-        Accepted and ignored on either value. Upstream 1.4 runs this operation in
245-        the background and lets the caller opt into synchronous execution with
246-        `async_operation=false`; Astrate always works synchronously and answers
247-        only once the work is done, so the value changes nothing and the
248-        response is identical either way (deviation 17 in COMPATIBILITY.md).
249-      schema:
250-        type: boolean
251-        default: false
252-
253-  schemas:
254-    # ── Envelope variants ──
255-
256-    DataEnvelopeRealm:
257-      type: object
258-      required: [data]
259-      properties:
260-        data:
261-          $ref: "#/components/schemas/Realm"

$ cd /root/astrate-mule && rg -n "^  [a-zA-Z]|^components:|^paths:|^  /" docs/api/astrate_realm_management_api.yaml | head -80
rg: docs/api/astrate_realm_management_api.yaml: IO error for operation on docs/api/astrate_realm_management_api.yaml: No such file or directory (os error 2)

$ cd /root/astrate-mule && ls docs/api/ && ls docs/
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
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

$ cd /root/astrate-mule && ls -la docs/api/ | cat
total 160
drwxr-xr-x 2 root root  4096 Sep  4 18:55 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 56583 Sep 12 22:44 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 15547 Sep 25 21:48 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 19506 Sep 17 13:03 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 44459 Sep 24 18:54 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 14977 Sep 20 21:25 astrate_native_api.yaml

$ cd /root/astrate-mule && wc -l docs/api/astarte_realm_management_api.yaml && rg -c "operationId" docs/api/astarte_realm_management_api.yaml
1327 docs/api/astarte_realm_management_api.yaml
20

$ cd /root/astrate-mule && rg -n "operationId" docs/api/astarte_realm_management_api.yaml
29:      operationId: listInterfaces
96:      operationId: installInterface
173:      operationId: listInterfaceMajors
204:      operationId: getInterface
236:      operationId: updateInterface
307:      operationId: deleteInterface
351:      operationId: listTriggers
377:      operationId: createTrigger
448:      operationId: getTrigger
479:      operationId: deleteTrigger
504:      operationId: getAuthConfig
537:      operationId: putAuthConfig
577:      operationId: getDeviceRegistrationLimit
610:      operationId: getDatastreamMaximumStorageRetention
647:      operationId: getVersion
682:      operationId: deleteDevice
710:      operationId: listPolicies
736:      operationId: createPolicy
784:      operationId: getPolicy
815:      operationId: deletePolicy

$ cd /root/astrate-mule && rg -n "^components:|^  parameters:|^  responses:|^  schemas:|^paths:" docs/api/astarte_realm_management_api.yaml
24:paths:
838:components:
848:  parameters:
898:  schemas:
1235:  responses:

→ Read docs/api/astarte_realm_management_api.yaml [offset=90, limit=130]
→ Read docs/api/astarte_realm_management_api.yaml [limit=120, offset=230]
→ Read docs/api/astarte_realm_management_api.yaml [limit=120, offset=670]
→ Read docs/api/astarte_realm_management_api.yaml [limit=70, offset=800]
→ Read docs/api/astarte_realm_management_api.yaml [limit=70, offset=838]

$ cd /root/astrate-mule && rg -n "operationBlock|componentBlock|containsLine" -A 30 internal/swagger/swagger_test.go | head -80
162:		if !containsLine(lines[start:end], `        "403":`) {
163-			t.Errorf("operation starting at line %d documents no 403 response", start+1)
164-		}
165-	}
166-
167:	if !containsLine(lines[compIdx:], `    Forbidden:`) {
168-		t.Error("components.responses defines no Forbidden response")
169-	}
170-}
171-
172-// TestHousekeepingAsyncOperationParamDocumented guards that the housekeeping
173-// spec tells clients the `?async_operation` parameter exists. Astrate accepts
174-// and ignores it on realm create and delete (deviation 17,
175-// docs/COMPATIBILITY.md), so an upstream client keeps working — but a client
176-// generated from the spec has to learn the parameter from the spec. This is
177-// the documentation half of the behaviour pinned by
178-// TestHousekeepingAsyncOperationParam in internal/housekeeping.
179-func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
180-	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
181-	if err != nil {
182-		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
183-	}
184-	lines := strings.Split(string(b), "\n")
185-
186-	const ref = `        - $ref: "#/components/parameters/AsyncOperation"`
187-	for _, op := range []string{"createRealm", "deleteRealm"} {
188:		block := operationBlock(t, lines, op)
189:		if !containsLine(block, ref) {
190-			t.Errorf("operation %s does not $ref the AsyncOperation parameter", op)
191-		}
192-	}
193-
194:	comp := componentBlock(t, lines, "    AsyncOperation:")
195-	for _, want := range []string{
196-		"      name: async_operation",
197-		"      in: query",
198-		"      required: false",
199-		"      schema:",
200-		"        type: boolean",
201-		"        default: false",
202-	} {
203:		if !containsLine(comp, want) {
204-			t.Errorf("components parameter AsyncOperation is missing line %q", want)
205-		}
206-	}
207-	if !strings.Contains(strings.Join(comp, "\n"), "Accepted and ignored") {
208-		t.Error("components parameter AsyncOperation does not say the value is accepted and ignored")
209-	}
210-}
211-
212-// propertyBlock returns the lines of the schema property with the given name
213-// inside a components.schemas entry, up to the next property or the end of the
214-// entry.
215-func propertyBlock(t *testing.T, lines []string, name string) []string {
216-	t.Helper()
217-	key := "        " + name + ":"
218-	start := -1
219-	for i, l := range lines {
220-		if l == key {
221-			start = i
222-			break
223-		}
224-	}
225-	if start < 0 {
226-		t.Fatalf("schema declares no property %q", name)
227-	}
228-	for i := start + 1; i < len(lines); i++ {
229-		if strings.TrimSpace(lines[i]) == "" {
230-			continue
231-		}
232-		if !strings.HasPrefix(lines[i], "          ") {
233-			return lines[start+1 : i]
--
258:	patch := strings.Join(operationBlock(t, lines, "patchRealm"), "\n")
259-	for _, want := range []string{
260-		"explicit `0` clears it too",
261-		"0 is folded to unset",
262-		"has no such fold",
263-		"stored as the literal limit `0`",
264-	} {

→ Read docs/COMPATIBILITY.md [offset=320, limit=40]

$ cd /root/astrate-mule && rg -n "^func operationBlock|^func componentBlock|^func containsLine" -A 30 internal/swagger/swagger_test.go
302:func operationBlock(t *testing.T, lines []string, operationID string) []string {
303-	t.Helper()
304-	marker := "      operationId: " + operationID
305-	start := -1
306-	for i, l := range lines {
307-		if l == marker {
308-			start = i
309-			break
310-		}
311-	}
312-	if start < 0 {
313-		t.Fatalf("spec declares no operation %q", operationID)
314-	}
315-	endRe := regexp.MustCompile(`^      operationId: |^  \S|^components:`)
316-	for i := start + 1; i < len(lines); i++ {
317-		if endRe.MatchString(lines[i]) {
318-			return lines[start+1 : i]
319-		}
320-	}
321-	return lines[start+1:]
322-}
323-
324-// componentBlock returns the lines of the components entry starting with the
325-// given key line, up to the next entry at the same or shallower indentation.
326:func componentBlock(t *testing.T, lines []string, key string) []string {
327-	t.Helper()
328-	start := -1
329-	for i, l := range lines {
330-		if l == key {
331-			start = i
332-			break
333-		}
334-	}
335-	if start < 0 {
336-		t.Fatalf("spec declares no components entry %q", key)
337-	}
338-	for i := start + 1; i < len(lines); i++ {
339-		if strings.TrimSpace(lines[i]) == "" {
340-			continue
341-		}
342-		if !strings.HasPrefix(lines[i], "      ") {
343-			return lines[start+1 : i]
344-		}
345-	}
346-	return lines[start+1:]
347-}
348-
349:func containsLine(lines []string, want string) bool {
350-	for _, l := range lines {
351-		if l == want {
352-			return true
353-		}
354-	}
355-	return false
356-}
357-
358-func embeddedYAMLFilenames() ([]string, error) {
359-	var names []string
360-	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
361-		if err != nil {
362-			return err
363-		}
364-		if d.IsDir() {
365-			return nil
366-		}
367-		if strings.HasSuffix(path, ".yaml") {
368-			names = append(names, strings.TrimPrefix(path, "api/"))
369-		}
370-		return nil
371-	})
372-	sort.Strings(names)
373-	return names, err
374-}
375-
376-func get(t *testing.T, url string) string {
377-	t.Helper()
378-	resp, err := http.Get(url)
379-	if err != nil {

$ cd /root/astrate-mule && rg -n "func \(a \*API\) (installInterface|updateInterface|deleteInterface|deletePolicy)\b" -A 30 internal/realm/http.go
94:func (a *API) deletePolicy(w http.ResponseWriter, r *http.Request) {
95-	if err := a.svc.DeletePolicy(r.Context(), r.PathValue("realm"), r.PathValue("name")); err != nil {
96-		a.writeError(w, err)
97-		return
98-	}
99-	w.WriteHeader(http.StatusNoContent)
100-}
101-
102-// deleteDevice synchronously removes a device and its data (the dashboard's
103-// Delete button; docs/COMPATIBILITY.md notes the deviation from upstream's
104-// async deletion).
105-func (a *API) deleteDevice(w http.ResponseWriter, r *http.Request) {
106-	err := a.svc.DeleteDevice(r.Context(), r.PathValue("realm"), r.PathValue("device"))
107-	if errors.Is(err, store.ErrNotFound) {
108-		_ = astarteapi.WriteDeviceNotFound(w)
109-		return
110-	}
111-	if err != nil {
112-		a.writeError(w, err)
113-		return
114-	}
115-	w.WriteHeader(http.StatusNoContent)
116-}
117-
118-// getVersion reports the emulated upstream API level (the dashboard gates
119-// feature UI on it — see APICompatVersion).
120-func (a *API) getVersion(w http.ResponseWriter, _ *http.Request) {
121-	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
122-}
123-
124-func (a *API) getRegistrationLimit(w http.ResponseWriter, r *http.Request) {
--
167:func (a *API) installInterface(w http.ResponseWriter, r *http.Request) {
168-	var def json.RawMessage
169-	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &def); err != nil {
170-		_ = astarteapi.WriteBadRequest(w)
171-		return
172-	}
173-	si, err := a.svc.InstallInterface(r.Context(), r.PathValue("realm"), def)
174-	if err != nil {
175-		a.writeInterfaceError(w, err)
176-		return
177-	}
178-	_ = astarteapi.WriteData(w, http.StatusCreated, json.RawMessage(si.Definition))
179-}
180-
181-func (a *API) listInterfaceMajors(w http.ResponseWriter, r *http.Request) {
182-	majors, err := a.svc.ListInterfaceMajors(r.Context(), r.PathValue("realm"), r.PathValue("name"))
183-	if err != nil {
184-		a.writeInterfaceError(w, err)
185-		return
186-	}
187-	_ = astarteapi.WriteData(w, http.StatusOK, majors)
188-}
189-
190-func (a *API) getInterface(w http.ResponseWriter, r *http.Request) {
191-	major, ok := majorParam(w, r)
192-	if !ok {
193-		return
194-	}
195-	def, err := a.svc.GetInterface(r.Context(), r.PathValue("realm"), r.PathValue("name"), major)
196-	if err != nil {
197-		a.writeInterfaceError(w, err)
--
203:func (a *API) updateInterface(w http.ResponseWriter, r *http.Request) {
204-	major, ok := majorParam(w, r)
205-	if !ok {
206-		return
207-	}
208-	var def json.RawMessage
209-	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &def); err != nil {
210-		_ = astarteapi.WriteBadRequest(w)
211-		return
212-	}
213-	if _, err := a.svc.UpdateInterface(r.Context(), r.PathValue("realm"), r.PathValue("name"), major, def); err != nil {
214-		a.writeInterfaceError(w, err)
215-		return
216-	}
217-	w.WriteHeader(http.StatusNoContent)
218-}
219-
220:func (a *API) deleteInterface(w http.ResponseWriter, r *http.Request) {
221-	major, ok := majorParam(w, r)
222-	if !ok {
223-		return
224-	}
225-	if err := a.svc.DeleteInterface(r.Context(), r.PathValue("realm"), r.PathValue("name"), major); err != nil {
226-		a.writeInterfaceError(w, err)
227-		return
228-	}
229-	w.WriteHeader(http.StatusNoContent)
230-}
231-
232-// --- triggers ---------------------------------------------------------------
233-
234-func (a *API) listTriggers(w http.ResponseWriter, r *http.Request) {
235-	names, err := a.svc.ListTriggers(r.Context(), r.PathValue("realm"))
236-	if err != nil {
237-		a.writeError(w, err)
238-		return
239-	}
240-	_ = astarteapi.WriteData(w, http.StatusOK, names)
241-}
242-
243-func (a *API) createTrigger(w http.ResponseWriter, r *http.Request) {
244-	var def json.RawMessage
245-	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &def); err != nil {
246-		_ = astarteapi.WriteBadRequest(w)
247-		return
248-	}
249-	tr, err := a.svc.CreateTrigger(r.Context(), r.PathValue("realm"), def)
250-	if err != nil {

Handlers confirmed: none of the four read the query string (only `listTriggers` reads `detailed`). All four accept and ignore it. Now the YAML edits.
← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -102,8 +102,9 @@
 security:
   - a_rma: []
 parameters:
   - $ref: "#/components/parameters/RealmName"
+  - $ref: "#/components/parameters/AsyncOperation"
 requestBody:
   required: true
   content:
     application/json; charset=utf-8:

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -245,8 +245,9 @@
 parameters:
   - $ref: "#/components/parameters/RealmName"
   - $ref: "#/components/parameters/InterfaceName"
   - $ref: "#/components/parameters/MajorVersion"
+  - $ref: "#/components/parameters/AsyncOperation"
 requestBody:
   required: true
   content:
     application/json; charset=utf-8:

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -318,8 +318,9 @@
 parameters:
   - $ref: "#/components/parameters/RealmName"
   - $ref: "#/components/parameters/InterfaceName"
   - $ref: "#/components/parameters/MajorVersion"
+  - $ref: "#/components/parameters/AsyncOperation"
 responses:
   "204":
     description: Interface deleted. No body returned.
   "401":

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -816,15 +816,17 @@
 
 delete:
   operationId: deletePolicy
   summary: Delete a trigger delivery policy
-  description: Deletes a policy by name. Returns 204 on success.
+  description: Deletes a policy by name. Returns 204 on success. The
+    `async_operation` parameter is accepted and ignored.
   tags: [Realm Management - Policies]
   security:
     - a_rma: []
   parameters:
     - $ref: "#/components/parameters/RealmName"
     - $ref: "#/components/parameters/PolicyName"
+    - $ref: "#/components/parameters/AsyncOperation"
   responses:
     "204":
       description: Policy deleted. No body returned.
     "401":

Now the component parameter definition:
← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -899,8 +899,23 @@
     description: The trigger delivery policy name.
     schema:
       type: string
 
+  AsyncOperation:
+    name: async_operation
+    in: query
+    required: false
+    description: |
+      Accepted and ignored on either value. Upstream 1.4 runs interface
+      install/update/delete and trigger-delivery-policy delete in the
+      background and lets the caller opt into synchronous execution with
+      `async_operation=false`; Astrate always works synchronously and answers
+      only once the work is done, so the value changes nothing and the
+      response is identical either way (deviation 17 in COMPATIBILITY.md).
+    schema:
+      type: boolean
+      default: false
+
 schemas:
   # ── Config types ──
 
   AuthConfig:

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -816,10 +816,9 @@
 
 delete:
   operationId: deletePolicy
   summary: Delete a trigger delivery policy
-  description: Deletes a policy by name. Returns 204 on success. The
-    `async_operation` parameter is accepted and ignored.
+  description: Deletes a policy by name. Returns 204 on success.
   tags: [Realm Management - Policies]
   security:
     - a_rma: []
   parameters:

$ cd /root/astrate-mule && rg -n "AsyncOperation" docs/api/astarte_realm_management_api.yaml
106:        - $ref: "#/components/parameters/AsyncOperation"
249:        - $ref: "#/components/parameters/AsyncOperation"
322:        - $ref: "#/components/parameters/AsyncOperation"
827:        - $ref: "#/components/parameters/AsyncOperation"
902:    AsyncOperation:

$ cd /root/astrate-mule && python3 -c "
import yaml
d=yaml.safe_load(open('docs/api/astarte_realm_management_api.yaml'))
for p,ops in d['paths'].items():
    for m,o in ops.items():
        refs=[x.get('\$ref') for x in o.get('parameters',[])]
        if any('AsyncOperation' in (r or '') for r in refs):
            print(m, o['operationId'], p)
print('component:', yaml.dump(d['components']['parameters']['AsyncOperation'])[:200])
"
Traceback (most recent call last):
  File "<string>", line 2, in <module>
    import yaml
ModuleNotFoundError: No module named 'yaml'

$ cd /root/astrate-mule && python3 - <<'EOF'
import re
lines=open('docs/api/astarte_realm_management_api.yaml').read().split('\n')
for i,l in enumerate(lines):
    if 'operationId' in l:
        op=l.split('operationId: ')[1]
        blk='\n'.join(lines[i:i+16])
        print(op, 'AsyncOperation' in blk)
EOF
listInterfaces False
installInterface True
listInterfaceMajors False
getInterface False
updateInterface True
deleteInterface True
listTriggers False
createTrigger False
getTrigger False
deleteTrigger False
getAuthConfig False
putAuthConfig False
getDeviceRegistrationLimit False
getDatastreamMaximumStorageRetention False
getVersion False
deleteDevice False
listPolicies False
createPolicy False
getPolicy False
deletePolicy True

Now the docs-half test:
→ Read internal/swagger/swagger_test.go [limit=40, offset=1]

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -208,8 +208,51 @@
 		t.Error("components parameter AsyncOperation does not say the value is accepted and ignored")
 	}
 }
 
+// TestRealmManagementAsyncOperationParamDocumented guards that the realm
+// management spec tells clients the `?async_operation` parameter exists on the
+// four operations upstream 1.4 runs in the background: interface
+// install/update/delete and trigger-delivery-policy delete. Astrate accepts and
+// ignores it there (deviation 17, docs/COMPATIBILITY.md) — none of the four
+// handlers reads the query string, which is the whole of the acceptance — so an
+// upstream client keeps working, but a client generated from the spec has to
+// learn the parameter from the spec. This is the documentation half of the
+// behaviour pinned by TestRealmManagementAsyncOperationParam in
+// internal/realm.
+func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	const ref = `        - $ref: "#/components/parameters/AsyncOperation"`
+	for _, op := range []string{"installInterface", "updateInterface", "deleteInterface", "deletePolicy"} {
+		block := operationBlock(t, lines, op)
+		if !containsLine(block, ref) {
+			t.Errorf("operation %s does not $ref the AsyncOperation parameter", op)
+		}
+	}
+
+	comp := componentBlock(t, lines, "    AsyncOperation:")
+	for _, want := range []string{
+		"      name: async_operation",
+		"      in: query",
+		"      required: false",
+		"      schema:",
+		"        type: boolean",
+		"        default: false",
+	} {
+		if !containsLine(comp, want) {
+			t.Errorf("components parameter AsyncOperation is missing line %q", want)
+		}
+	}
+	if !strings.Contains(strings.Join(comp, "\n"), "Accepted and ignored") {
+		t.Error("components parameter AsyncOperation does not say the value is accepted and ignored")
+	}
+}
+
 // propertyBlock returns the lines of the schema property with the given name
 // inside a components.schemas entry, up to the next property or the end of the
 // entry.
 func propertyBlock(t *testing.T, lines []string, name string) []string {

$ cd /root/astrate-mule && cp docs/api/astarte_realm_management_api.yaml /tmp/opencode/rm_new.yaml && go test ./internal/swagger/ -run 'AsyncOperationParamDocumented' -v 2>&1 | tail -20
=== RUN   TestHousekeepingAsyncOperationParamDocumented
--- PASS: TestHousekeepingAsyncOperationParamDocumented (0.00s)
=== RUN   TestRealmManagementAsyncOperationParamDocumented
--- PASS: TestRealmManagementAsyncOperationParamDocumented (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.019s

Now verify the test actually fails without the YAML change (using a read-only `git show` of HEAD, not touching git state):
$ cd /root/astrate-mule && git show HEAD:docs/api/astarte_realm_management_api.yaml > /tmp/opencode/rm_orig.yaml && cp /tmp/opencode/rm_orig.yaml docs/api/astarte_realm_management_api.yaml && go test ./internal/swagger/ -run 'TestRealmManagementAsyncOperationParamDocumented' 2>&1 | tail -12; cp /tmp/opencode/rm_new.yaml docs/api/astarte_realm_management_api.yaml && git status --porcelain
--- FAIL: TestRealmManagementAsyncOperationParamDocumented (0.01s)
    swagger_test.go:233: operation installInterface does not $ref the AsyncOperation parameter
    swagger_test.go:233: operation updateInterface does not $ref the AsyncOperation parameter
    swagger_test.go:233: operation deleteInterface does not $ref the AsyncOperation parameter
    swagger_test.go:233: operation deletePolicy does not $ref the AsyncOperation parameter
    swagger_test.go:237: spec declares no components entry "    AsyncOperation:"
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.029s
FAIL
 M docs/api/astarte_realm_management_api.yaml
 M internal/swagger/swagger_test.go

→ Read docs/Makefile

$ cd /root/astrate-mule && which mkdocs; mkdocs --version 2>&1 | head -2
/bin/sh: 1: mkdocs: not found

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

$ cd /root/astrate-mule && git status --porcelain | head -20 && rg -n "site" .gitignore 2>/dev/null | head
M docs/api/astarte_realm_management_api.yaml
 M internal/swagger/swagger_test.go
29:/site-dist/
31:# Docs site copies (refreshed by docs/make sync)
32:/docs/site/DESIGN.md
33:/docs/site/ROADMAP.md
34:/docs/site/api/
35:/docs/site/swagger-ui/

$ cd /root/astrate-mule && cat docs/requirements.txt && python3 -m pip --version 2>&1 | head -2 && python3 -c "import venv; print('venv ok')"
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10
/usr/bin/python3: No module named pip
venv ok

$ cd /tmp/opencode && timeout 90 python3 -m venv .mkvenv 2>&1 | tail -3; ls .mkvenv/bin 2>/dev/null | head
Activate.ps1
activate
activate.csh
activate.fish
pip
pip3
pip3.13
python
python3
python3.13

$ cd /tmp/opencode && timeout 180 ./.mkvenv/bin/pip install -q -r /root/astrate-mule/docs/requirements.txt 2>&1 | tail -5; ./.mkvenv/bin/mkdocs --version 2>&1 | head -2
mkdocs, version 1.6.1 from /tmp/opencode/.mkvenv/lib/python3.13/site-packages/mkdocs (Python 3.13)

$ cd /root/astrate-mule && PATH=/tmp/opencode/.mkvenv/bin:$PATH make -C docs build 2>&1 | tail -25
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
INFO    -  Documentation built in 8.51 seconds
make: Leaving directory '/root/astrate-mule/docs'

$ cd /root/astrate-mule && rg -n "embeddedYAMLFilenames" -B 5 -A 40 internal/swagger/swagger_test.go | head -80
104-		if _, err := docs.APIYAML.ReadFile("api/" + name); err != nil {
105-			t.Errorf("Specs() entry %q cannot be read from docs.APIYAML: %v", name, err)
106-		}
107-	}
108-
109:	want, err := embeddedYAMLFilenames()
110-	if err != nil {
111-		t.Fatalf("enumerating embedded basenames: %v", err)
112-	}
113-	if !reflect.DeepEqual(got, want) {
114-		t.Errorf("Specs() = %v, want %v", got, want)
115-	}
116-}
117-
118-// TestRealmManagement403 guards the contract that every realm-management
119-// operation guarded by a_rma documents a 403 Forbidden response, not just
120-// 401: RequireRealm answers 403 for a verified realm JWT whose a_rma grants
121-// do not authorize the method+path (internal/auth/middleware.go). Matches the
122-// upstream-parity the housekeeping and pairing specs already document.
123-func TestRealmManagement403(t *testing.T) {
124-	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
125-	if err != nil {
126-		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
127-	}
128-	lines := strings.Split(string(b), "\n")
129-
130-	pathIdx, compIdx := -1, -1
131-	for i, l := range lines {
132-		if l == "paths:" {
133-			pathIdx = i
134-		}
135-		if l == "components:" {
136-			compIdx = i
137-		}
138-		if pathIdx >= 0 && compIdx >= 0 {
139-			break
140-		}
141-	}
142-	if pathIdx < 0 || compIdx < pathIdx {
143-		t.Fatal("cannot locate the paths and components sections")
144-	}
145-
146-	methodRe := regexp.MustCompile(`^    (get|post|put|delete):$`)
147-	var ops []int
148-	for i, l := range lines[pathIdx:compIdx] {
149-		if methodRe.MatchString(l) {
--
396-		}
397-	}
398-	return false
399-}
400-
401:func embeddedYAMLFilenames() ([]string, error) {
402-	var names []string
403-	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
404-		if err != nil {
405-			return err
406-		}
407-		if d.IsDir() {
408-			return nil
409-		}
410-		if strings.HasSuffix(path, ".yaml") {
411-			names = append(names, strings.TrimPrefix(path, "api/"))
412-		}
413-		return nil
414-	})
415-	sort.Strings(names)
416-	return names, err
417-}
418-
419-func get(t *testing.T, url string) string {
420-	t.Helper()
421-	resp, err := http.Get(url)
422-	if err != nil {
423-		t.Fatalf("GET %s: %v", url, err)
424-	}
425-	defer resp.Body.Close()
426-
427-	if resp.StatusCode != http.StatusOK {
428-		t.Fatalf("GET %s: status = %d, want %d", url, resp.StatusCode, http.StatusOK)

$ cd /root/astrate-mule && rg -n "yaml" go.mod
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

$ mkdir -p /tmp/opencode/yamlprobe && cd /tmp/opencode/yamlprobe && cat > main.go <<'EOF'
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	b, err := os.ReadFile("/root/astrate-mule/docs/api/astarte_realm_management_api.yaml")
	if err != nil {
		panic(err)
	}
	var doc map[string]any
	if err := yaml.Unmarshal(b, &doc); err != nil {
		fmt.Println("PARSE FAIL:", err)
		os.Exit(1)
	}
	comps := doc["components"].(map[string]any)["parameters"].(map[string]any)
	fmt.Println("parsed ok; async_operation component:", comps["AsyncOperation"] != nil)
	paths := doc["paths"].(map[string]any)
	for p, ops := range paths {
		for m, o := range ops.(map[string]any) {
			op := o.(map[string]any)
			params, _ := op["parameters"].([]any)
			for _, pr := range params {
				if ref, ok := pr.(map[string]any)["$ref"].(string); ok && ref == "#/components/parameters/AsyncOperation" {
					fmt.Printf("  %-6s %-20s %s\n", m, op["operationId"], p)
				}
			}
		}
	}
	_ = http.StatusOK
	_ = httptest.NewRecorder
}
EOF
cat > go.mod <<'EOF'
module yamlprobe

go 1.24

require gopkg.in/yaml.v3 v3.0.1
EOF
GOFLAGS=-mod=mod go run . 2>&1 | tail -20
parsed ok; async_operation component: true
  put    updateInterface      /realmmanagement/v1/{realm}/interfaces/{name}/{major}
  delete deleteInterface      /realmmanagement/v1/{realm}/interfaces/{name}/{major}
  delete deletePolicy         /realmmanagement/v1/{realm}/policies/{name}
  post   installInterface     /realmmanagement/v1/{realm}/interfaces

$ cd /root/astrate-mule && uname -m && go vet ./... 2>&1 | tail -5
aarch64

$ cd /root/astrate-mule && gofmt -l . 2>&1 | head; go test ./... 2>&1 | rg -v "^ok|no test files" | head -20
(no output)

$ cd /root/astrate-mule && timeout 240 go test -race ./internal/swagger/ 2>&1 | tail -5
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.048s
FAIL

$ cd /root/astrate-mule && rg -n "deleteDevice|async" .mule/todo.md | head; echo ---; tail -12 .mule/todo.md
226:- [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
228:- [ ] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
---
- [x] docs-sync-rm-interface-422-shapes [auto]: model the `422` on `POST /realmmanagement/v1/{realm}/interfaces` (yaml:133-134) and `PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}` (yaml:231-232) in docs/api/astarte_realm_management_api.yaml as a oneOf — both handlers go through `writeInterfaceError`, which answers three distinct 422 bodies: the flat ErrorDetail (`ErrValidation` via `validationDetail`, http.go:392-393), the nested violations changeset envelope (`writeViolations`, http.go:358-359 + 434-533, e.g. `{"errors":{"description":["should be at most 1000 character(s)"]}}` and the aligned full-length `mappings` array, http_test.go:480-501), and the named FieldErrors envelope for `ErrMaximumDatabaseRetentionExceeded` (`{"errors":{"error_name":["maximum_database_retention_exceeded"]}}`, http.go:355-357, http_test.go:282-320); today only the flat ValidationError is referenced. Mirror the createTrigger oneOf pattern (yaml:344-368). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-validation-example-prefix [auto]: fix the `ValidationError` example in docs/api/astarte_realm_management_api.yaml (yaml:1160-1161) — it shows `detail: "realm: validation failed: interface definition is invalid"`, but `validationDetail` strips the `realm: validation failed: ` prefix on the wire (http.go:562-568) and the parser's real messages are `invalid interface: ...` (pkg/interfaceschema/parse.go:38, wrapped at service.go:163); the example should show the stripped, real message. Same class as the already-fixed docs-sync-hk-validation-example (housekeeping http.go:251-258). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-auth-403 [auto]: add the missing `403` Forbidden response to the a_rma-guarded ops in docs/api/astarte_realm_management_api.yaml — `RequireRealm` answers 403 for a verified realm JWT whose a_rma grants do not authorize the method+path (internal/auth/middleware.go:90, 115, 149-151, upstream-parity per the comment at middleware.go:29-32); the housekeeping and pairing specs already document Forbidden on every guarded route, but the RM spec lists only 401 on every op except `DELETE /interfaces`, whose 403 at yaml:257 is the business-logic one. Add the Forbidden response (ErrorDetail envelope, `detail: Forbidden`) to the a_rma-guarded operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] swagger-sub-failfast [auto]: make `Mount` fail fast (panic) when `fs.Sub` over the `docs.SwaggerUI`/`docs.APIYAML` embeds errors — `uiRoot, _ := fs.Sub(...)` and `apiRoot, _ := fs.Sub(...)` (internal/swagger/swagger.go:17-18) silently serve an empty tree → 404s at /swagger/ and /api/ with no log the day the docs embed layout changes; add an fs-injecting internal variant (`MountWithFS(mux, uiRoot, apiRoot fs.FS)` called by `Mount` after the subs) so a broken-FS unit test in swagger_test.go can prove the fail-fast, keeping the public `Mount(mux)` signature and its single main.go:414 caller untouched. Container-free.
- [ ] flow-boot-resumes-stopped-durable-flows [legion] [auto]: RehydrateAutoRestart (cmd/astrate/main.go:216) restarts every flow with auto_restart=true whatever its last status — `ListAutoRestartFlows` filters on `f.auto_restart = true` alone (internal/store/flows.go:125) and nothing ever clears the flag, which is written once at create (flowapi/http.go:141-146 -> service.go:318) and never on stop — so a durable flow stopped through the API is running again after the next boot, while the shutdown mark that says otherwise (MarkRunningFlowsStopped, main.go:497 -> internal/flowapi/service.go:885-901, which writes status="stopped") is never read by that query. The two halves disagree: service.go:325 documents "every durable flow", the shutdown call assumes "only what was running". Say which is the intent, then make them agree — either add `AND f.status = 'running'` to flows.go:125 (making the shutdown mark load-bearing) or delete the mark as dead code — and add a case in internal/store/flows_test.go pinning the chosen rule (a stopped auto_restart row must, or must not, come back). Needs the DB.
- [x] drain-per-stage-budget [auto]: `shutdown` spends a single 30s sctx (cmd/astrate/main.go:479) across three sequential stages — srv.Shutdown (484), b.Close (489), flowSvc.Manager().Shutdown (494) and MarkRunningFlowsStopped (497) — while `drainEngine` alone opens a second one (504), so the constant's own comment "bounds the whole graceful drain" (main.go:55-56) is wrong in both directions, and one slow in-flight HTTP request (the server sets no per-request deadline, 221) burns the whole first budget and leaves the flow stage on an expired context: StopFlow then returns at `f.router.Drain(ctx)` before step 3 releases block resources (internal/flow/flow.go:265-273) and MarkRunningFlowsStopped skips every flow (flowapi/service.go:895-897), both surfacing as one log.Warn. Give each stage its own budget through an extracted helper (shutdownTimeout has to become an injectable var for the test) and pin it container-free: a second stage still receives a live context after a first stage exhausts its own.
- [x] cmd-devcert-fields-test [auto]: add a container-free test for `selfSignedDevCert` (cmd/astrate/main.go:320-345), the throwaway mTLS identity `insecure_dev_mode` boots on, which no test touches — pin NotBefore = now-1h (331; a plain now would reject a device whose clock is behind), NotAfter = now+365d (332), DNSNames localhost and IPAddresses 127.0.0.1/::1 (336-337, dropping which breaks hostname verification for a device dialling mqtts://127.0.0.1:8883), ExtKeyUsage serverAuth (334), and that the leaf is self-signed and parses (x509.ParseCertificate on Certificate[0]).
- [!] cmd-healthcheck-contract-test [auto]: add a container-free test for `runHealthcheck` (cmd/astrate/main.go:531-558) — the `HEALTHCHECK CMD ["/astrate", "-healthcheck"]` contract (Dockerfile:35-36) has no coverage. Assert it probes /astrate/v1/readiness and not /health (543), returns 0 only on 200 so a 503 readiness exits 1 (554), and maps a bare ":PORT" in ASTRATE_HTTP_ADDR to 127.0.0.1 (536-538). httptest.NewServer plus t.Setenv("ASTRATE_HTTP_ADDR", srv.Listener.Addr().String()). — BLOCKED: tests failed: --- FAIL: TestRunHealthcheckProbesReadiness (0.01s)
- [x] cmd-loadsealer-masterkeyfile-test [auto]: add a container-free test for `loadSealer` (cmd/astrate/main.go:515-526) — the `security.master_key_file` to env hand-off (`os.Setenv(store.EnvMasterKeyFile, ...)`, 517) is the only wiring that makes a configured key file reach the store loader, and dropping or misspelling that branch leaves the file silently ignored and boot dying with the generic "master key" error (523) even though the config named a valid file. Assert a 32-byte key file under t.TempDir() loads, the inline ASTRATE_MASTER_KEY path still wins, and a missing key returns the error naming all three env vars. t.Setenv both ASTRATE_MASTER_KEY and ASTRATE_MASTER_KEY_FILE first so the os.Setenv is restored at cleanup.
- [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

→ Read .mule/todo.md [offset=228, limit=5]
✗ Edit .mule/todo.md failed
Error: Found multiple matches for oldString. Provide more surrounding context to make the match unique.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -225,4 +225,5 @@
 - [x] cmd-loadsealer-masterkeyfile-test [auto]: add a container-free test for `loadSealer` (cmd/astrate/main.go:515-526) — the `security.master_key_file` to env hand-off (`os.Setenv(store.EnvMasterKeyFile, ...)`, 517) is the only wiring that makes a configured key file reach the store loader, and dropping or misspelling that branch leaves the file silently ignored and boot dying with the generic "master key" error (523) even though the config named a valid file. Assert a 32-byte key file under t.TempDir() loads, the inline ASTRATE_MASTER_KEY path still wins, and a missing key returns the error naming all three env vars. t.Setenv both ASTRATE_MASTER_KEY and ASTRATE_MASTER_KEY_FILE first so the os.Setenv is restored at cleanup.
 - [x] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] docs-sync-rm-async-operation-param [auto]: document the `?async_operation` boolean query parameter on the four realm-management operations upstream 1.4 runs in the background — interface install/update/delete and trigger-delivery-policy delete — in docs/api/astarte_realm_management_api.yaml, the same gap the housekeeping spec just closed (deviation 17, docs/COMPATIBILITY.md:336-347; TestRealmManagementAsyncOperationParam in internal/realm/async_operation_test.go:16 pins that neither value 4xxs). Check first which of the four RM handlers actually read the query string before adding the `AsyncOperation` component parameter, and add a docs-half test beside TestHousekeepingAsyncOperationParamDocumented in internal/swagger/swagger_test.go — the realm-management `deleteDevice` operation already calls out the deviation in prose (yaml:684-687) but declares no such parameter. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cd /root/astrate-mule && git status --porcelain && git diff --stat && go vet ./... && go test ./internal/swagger/ 2>&1 | tail -3 && gofmt -l . | head -3
M .mule/todo.md
 M docs/api/astarte_realm_management_api.yaml
 M internal/swagger/swagger_test.go
 .mule/todo.md                              |  1 +
 docs/api/astarte_realm_management_api.yaml | 19 +++++++++++++
 internal/swagger/swagger_test.go           | 43 ++++++++++++++++++++++++++++++
 3 files changed, 63 insertions(+)
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)

$ cd /root/astrate-mule && git diff docs/api/astarte_realm_management_api.yaml
diff --git a/docs/api/astarte_realm_management_api.yaml b/docs/api/astarte_realm_management_api.yaml
index 769a02f..5a229f6 100644
--- a/docs/api/astarte_realm_management_api.yaml
+++ b/docs/api/astarte_realm_management_api.yaml
@@ -103,6 +103,7 @@ paths:
         - a_rma: []
       parameters:
         - $ref: "#/components/parameters/RealmName"
+        - $ref: "#/components/parameters/AsyncOperation"
       requestBody:
         required: true
         content:
@@ -245,6 +246,7 @@ paths:
         - $ref: "#/components/parameters/RealmName"
         - $ref: "#/components/parameters/InterfaceName"
         - $ref: "#/components/parameters/MajorVersion"
+        - $ref: "#/components/parameters/AsyncOperation"
       requestBody:
         required: true
         content:
@@ -317,6 +319,7 @@ paths:
         - $ref: "#/components/parameters/RealmName"
         - $ref: "#/components/parameters/InterfaceName"
         - $ref: "#/components/parameters/MajorVersion"
+        - $ref: "#/components/parameters/AsyncOperation"
       responses:
         "204":
           description: Interface deleted. No body returned.
@@ -821,6 +824,7 @@ paths:
       parameters:
         - $ref: "#/components/parameters/RealmName"
         - $ref: "#/components/parameters/PolicyName"
+        - $ref: "#/components/parameters/AsyncOperation"
       responses:
         "204":
           description: Policy deleted. No body returned.
@@ -895,6 +899,21 @@ components:
       schema:
         type: string
 
+    AsyncOperation:
+      name: async_operation
+      in: query
+      required: false
+      description: |
+        Accepted and ignored on either value. Upstream 1.4 runs interface
+        install/update/delete and trigger-delivery-policy delete in the
+        background and lets the caller opt into synchronous execution with
+        `async_operation=false`; Astrate always works synchronously and answers
+        only once the work is done, so the value changes nothing and the
+        response is identical either way (deviation 17 in COMPATIBILITY.md).
+      schema:
+        type: boolean
+        default: false
+
   schemas:
     # ── Config types ──

$ cd /root/astrate-mule && ls .mule/ && rg -n "todo.md" .mule/*.sh .mule/*.md 2>/dev/null | head -10
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
recipes
reports
reviews
task.md
tasks
todo.md
.mule/MULE.md:25:  and add the rest to `.mule/todo.md` as new task lines.** Splitting is success, not failure.
.mule/MULE.md:109:Follow-ups: <task lines you appended to .mule/todo.md, or "none">
.mule/MULE.md:132:Append to `.mule/todo.md`, one line each, never more than ~8 at a time:
.mule/for-giulio.md:16:`.mule/todo.md` (`docs-sync-hk-async-operation-param`, `docs-sync-hk-retention-zero-is-unset`);
.mule/for-giulio.md:77:- **github-issues triage run, 2026-09-25: still nothing proposable — the mule-alarm pile is 18 issues (#94–#111), and two corrections to what earlier runs inferred about it.** Twenty-two open issues, same set as every run since 2026-09-05: **#111–#94** are the daily `mule-alarm` noise (zero comments, telemetry not code), **#93** aclhook comment (mule-review, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable upstream v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-24 run: **#111** (today, 09-25 11:01:07Z, "nothing has landed in 16h"). Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines added**. **Proposal, extending the 2026-09-24 one: close #94–#110**; leave #111 (live today) to self-expire. Two corrections, both from reading `tools/mule.sh` rather than the issue titles: **(1) the "missed days" in earlier entries (09-15, 09-20, 09-23) are evidence of activity, not its absence** — I previously called them timer gaps; the alarm is gated on a `.alarmed` sentinel that `beat` deletes on every land (`mule.sh:737,748-749`), so a day with no new issue is a day something landed, not a day the timer slept. **(2) the alarm body's closing line is wrong**: "Close it — the mule reopens a new one if the silence continues" (`mule.sh:764`) implies issue state drives the alarm; it does not. Nothing reads or reconciles the 18 open issues, so the pile can only grow and closing the old ones is safe but inert — and `cmd_refill` commits `mule: refill the queue` (mule.sh:719-720) **without calling `beat`**, so a doc-only refill never counts as landing work. The titles' ages also imply a beat roughly every day around 19:00–21:00Z (#111 filed 09-25 11:01Z at 16h → last beat ≈ 09-24 19:01Z; #110 filed 09-24 10:54Z at 15h → ≈ 09-23 19:54Z) while the ~11:00Z daily check has never once found under the 8h threshold — **what calls `beat` at that time is unverified**, and one command on the Pi settles it (`grep -E 'landed|checked:' /root/astrate-mule/.mule/log | tail -20`). The queue has still landed nothing since ~2026-09-04/05; if that idle streak is not intentional, the wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:85:- **github-issues triage run, 2026-09-24: still nothing proposable — the mule-alarm pile is now 17 days (#94–#110).** Twenty-one open issues, same set: **#93** aclhook comment (mule-review, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable upstream v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-20 runs: the alarms **#108** (09-21), **#109** (09-22) and **#110** (today, 09-24 10:54Z, the live one) — a four-day gap on 09-23 not withstanding. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines added**. **Proposal, extending the 2026-09-20 one: close #94–#109** — each was a one-day event, superseded by the next day's alarm, never actionable; leave #110 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:93:- **github-issues triage run, 2026-09-23: still nothing proposable — the mule-alarm pile is now 16 alarms (#94–#109).** Twenty open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-21 run: **#109**, yesterday's alarm (created 09-22 10:54Z, "nothing has landed in 14h"); **#108** has expired into the pile (superseded by #109). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-21 one: close #94–#108** — each was a one-day event, superseded by next day's alarm, never actionable; leave #109 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — now a ~17-day idle streak — if that is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:101:- **github-issues triage run, 2026-09-21: still nothing proposable — the mule-alarm pile is now 16 days (#94–#108).** Eighteen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-20 evening run: **#108**, today's alarm (created 09-21 10:54Z, "nothing has landed in 15h"); **#107** has expired into the pile (superseded by #108). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-20 one: close #94–#107** — each was a one-day event, superseded by next day's alarm, never actionable; leave #108 (live today) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.
.mule/for-giulio.md:109:- **github-issues triage re-run, 2026-09-20 (evening): nothing proposable, no change since this morning's 13:10 run.** Re-surveyed with the recipe command: the issue set is identical to the 2026-09-20 morning triage — **#94–#107** are still the daily `mule-alarm` noise (14 issues, zero comments, telemetry not code), **#93** aclhook comment rewrite is still mule-review with `8c61268` already pushed (its own recipe path), **#92** keyAgreement is still parked on a stable upstream v1.4.0 with its waiting-on.md row and the mule-upstream-watch escalation wired (an -rc does not satisfy it), **#78** FDO is still the milestone-4.0 design/investigation already escalated, and **#1** is untouched per standing instruction. No new issue has been filed since the morning run. Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines proposed**. This morning's proposal stands unchanged: **close #94–#106** — each was a one-day event, superseded by the next day's alarm, never actionable — and leave **#107** (still the newest, no #108 yet) to self-expire. The queue idle streak since ~2026-09-04/05 persists; if that is not intentional the wedge is on the Pi (queue review = the fix, per the 13:10 line).
.mule/for-giulio.md:121:- **github-issues triage run, 2026-09-20: still nothing proposable — the mule-alarm pile is now 15 days (#94–#107).** Seventeen open issues, same set: **#93** aclhook comment (mule-review, pushed `8c61268`, its own recipe path), **#92** keyAgreement (parked on stable v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-18 run: **#107**, yesterday's alarm (created 09-19 10:49Z, "nothing has landed in 14h"); **#106** has expired into the pile (superseded by #107). Still zero machine-checkable code-fix candidates, so no todo.md task lines. **Proposal, extending the 2026-09-18 one: close #94–#106** — each was a one-day event, superseded by next day's alarm, never actionable; leave #107 (live) to self-expire. The queue has landed nothing since ~2026-09-04/05 — if that idle streak is not intentional, the queue wedge is on the Pi and a queue review is the fix.

$ cd /root/astrate-mule && rg -n "todo.md" tools/mule.sh | head -20
13:#   mule.sh next             run the first open task in .mule/todo.md
31:TODO="$MULE/todo.md"
372:# A task is one `- [ ] slug: title` line in .mule/todo.md. If .mule/tasks/<slug>.md exists
661:# The queue is the set of open issues labelled `mule`, plus the standing lines in todo.md.
662:# Issue-derived tasks deliberately do NOT get copied into todo.md: that file lives in git on
892:   tasks — it appends them to .mule/todo.md for you to approve
934:That recipe is a *proposal* job: its output is new task lines appended to .mule/todo.md,
940:and does not need the todo.md approval step. But then that is the *whole* deliverable for
965:  note "edit .mule/todo.md to approve/cut, then: tools/mule.sh loop"

$ cd /root/astrate-mule && rg -n "cmd_land|\[x\]|done_line|mark_done" tools/mule.sh | head -20
261:# makes sed read the *script* as a filename, so the completed task was never marked [x] and
557:    [ "$lineno" != 0 ] && sed_i "${lineno}s/^- \[ \]/- [x]/" "$TODO"

```
Done: documented `?async_operation` on the four background-capable realm-management operations in docs/api/astarte_realm_management_api.yaml (new `components.parameters.AsyncOperation`, `$ref`ed from installInterface, updateInterface, deleteInterface, deletePolicy) and pinned it with a new TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go.
Files: docs/api/astarte_realm_management_api.yaml, internal/swagger/swagger_test.go, .mule/todo.md
Verified: `go vet ./...` -> clean; `go test ./...` -> pass (no `-race`: this box is aarch64 with the 39-bit VMA kernel the race detector rejects, `ThreadSanitizer: unsupported VMA range`; no concurrency touched anyway); `gofmt -l .` -> empty; `go test ./internal/swagger/ -run AsyncOperationParamDocumented` -> pass, and FAIL on the pre-change YAML (restored from a read-only `git show HEAD:` copy, tree verified clean afterwards) with all four missing refs + missing component; `make -C docs build` -> "Documentation built in 8.51 seconds" (needed a throwaway venv at /tmp/opencode/.mkvenv, mkdocs was not installed; remaining warnings about `nav` and `index.md`'s `api/` link are pre-existing); a /tmp probe parsed the YAML with gopkg.in/yaml.v3 and listed exactly those four operations as the $ref holders.
Unsure: whether `deleteDevice` was in scope — the task names "the four" twice, so I left it alone and queued it; its handler also ignores the query string, so the doc claim would be true. Also assumed the shared component description's phrase "Accepted and ignored" and the six `schema` lines must match the housekeeping twin verbatim (MULE.md wants the test to fail without the change, so I copied the twin's assertions rather than inventing new ones).
Follow-ups: one line appended to .mule/todo.md — `docs-sync-rm-delete-device-async-operation-param` (deleteDevice declares no such parameter although its prose names the deviation; deviation 17 in COMPATIBILITY.md omits the device leg from its list, and TestRealmManagementAsyncOperationParam never sends the parameter on a device delete).
```
