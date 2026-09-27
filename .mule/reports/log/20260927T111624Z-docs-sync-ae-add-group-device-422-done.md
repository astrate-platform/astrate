slug: docs-sync-ae-add-group-device-422
verdict: done
at:  dff3be8
ran: 2026-09-27T11:16:24Z on DietPi in 274s


> build · big-pickle

I'll start by reading the required context files.
→ Read .mule/MULE.md
→ Read .mule/task.md

# Todos
[•] Read yaml:939-980 and the ValidationErrors response at yaml:1707-1716
[ ] Verify the 422 behaviour claim in Go source (service.go:631-634, http.go:647)
[ ] Add the 422 response entry to the POST group-devices path
[ ] Run `make -C docs build` and confirm the spec still loads
[ ] Run gate: go test ./... and gofmt -l .

→ Read docs/api/astarte_appengine_api.yaml [offset=925, limit=75]
→ Read docs/api/astarte_appengine_api.yaml [offset=1690, limit=45]
✱ Grep "^    ValidationErrors:|^    Conflict:|^    NotFound:|^    InternalServerError:" in docs/api/astarte_appengine_api.yaml · 12 matches
→ Read docs/api/astarte_appengine_api.yaml [limit=55, offset=999]
✗ Read docs/api/astrate_appengine_api.yaml failed [limit=30, offset=1780]
Error: File not found: /root/astrate-mule/docs/api/astrate_appengine_api.yaml

$ rg -n "FieldErrors|WriteFieldErrors|StatusUnprocessableEntity" internal/appengine/service.go internal/appengine/http.go | head -40
internal/appengine/service.go:53:// FieldErrors carries per-field validation failures rendered as the
internal/appengine/service.go:55:type FieldErrors map[string][]string
internal/appengine/service.go:57:func (fe FieldErrors) Error() string { return "appengine: invalid request payload" }
internal/appengine/service.go:60:func (fe FieldErrors) addf(field, format string, args ...any) {
internal/appengine/service.go:467:	fe := FieldErrors{}
internal/appengine/service.go:550:			return nil, FieldErrors{"from_token": {"is invalid"}}
internal/appengine/service.go:554:		return nil, FieldErrors{"limit": {"must be greater than or equal to 0"}}
internal/appengine/service.go:633:		return FieldErrors{"device_id": {"is not a valid device id"}}
internal/appengine/http.go:299:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
internal/appengine/http.go:328:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
internal/appengine/http.go:347:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Value size exceeds size limits")
internal/appengine/http.go:406:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
internal/appengine/http.go:407:			FieldErrors{"group_name": {"can't be blank"}})
internal/appengine/http.go:410:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
internal/appengine/http.go:411:			FieldErrors{"devices": {"can't be blank"}})
internal/appengine/http.go:414:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
internal/appengine/http.go:415:			FieldErrors{"devices": {"should have at least 1 item(s)"}})
internal/appengine/http.go:552:// back as FieldErrors so writeError renders upstream's changeset-shaped 422
internal/appengine/http.go:557:	fe := FieldErrors{}
internal/appengine/http.go:644:	var fe FieldErrors
internal/appengine/http.go:647:		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity, fe)
internal/appengine/http.go:657:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, "Attribute key not found")
internal/appengine/http.go:659:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))

$ sed -n '1770,1800p' docs/api/astrate_appengine_api.yaml
sed: can't read docs/api/astrate_appengine_api.yaml: No such file or directory

$ sed -n '620,640p' internal/appengine/service.go; echo '--- http.go 480-500'; sed -n '480,500p' internal/appengine/http.go; echo '--- http.go 630,665'; sed -n '630,665p' internal/appengine/http.go
return 0, false
	}
	return int(binary.BigEndian.Uint32(u[0:4])), true
}

// AddGroupDevice adds a device to a group.
func (s *Service) AddGroupDevice(ctx context.Context, realm, name, deviceID string) error {
	rid, g, err := s.group(ctx, realm, name)
	if err != nil {
		return err
	}
	id, err := deviceid.Parse(deviceID)
	if err != nil {
		return FieldErrors{"device_id": {"is not a valid device id"}}
	}
	if err := s.st.AddGroupDevice(ctx, g.ID, rid, id); err != nil {
		if errors.Is(err, store.ErrAlreadyExists) {
			return fmt.Errorf("%w: %w", ErrDeviceAlreadyInGroup, err)
		}
		return err
	}
--- http.go 480-500

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
--- http.go 630,665
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

$ pwd; ls; ls docs 2>&1 | head
/root/astrate-mule
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

$ ls -la docs/api | head -20
total 168
drwxr-xr-x 2 root root  4096 Sep  4 18:55 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 59057 Sep 27 13:09 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 15547 Sep 25 21:48 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 19506 Sep 17 13:03 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 45387 Sep 25 22:07 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 14977 Sep 20 21:25 astrate_native_api.yaml

✗ Read docs/api/astrate_appengine_api.yaml failed [limit=30, offset=1775]
Error: File not found: /root/astrate-mule/docs/api/astrate_appengine_api.yaml

$ awk 'NR>=1775 && NR<=1800 {print NR": "$0}' docs/api/astrate_appengine_api.yaml
awk: cannot open "docs/api/astrate_appengine_api.yaml" (No such file or directory)

$ awk 'NR>=1775 && NR<=1800 {print NR": "$0}' /root/astrate-mule/docs/api/astarte_appengine_api.yaml
1775:       content:
1776:         application/json; charset=utf-8:
1777:           schema:
1778:             $ref: "#/components/schemas/ErrorDetail"
1779:           example:
1780:             errors:
1781:               detail: Internal Server Error
1782: 
1783:     ValidationErrors:
1784:       description: Validation Error — one or more fields failed validation.
1785:       content:
1786:         application/json; charset=utf-8:
1787:           schema:
1788:             $ref: "#/components/schemas/FieldErrorsDetail"
1789:           example:
1790:             errors:
1791:               group_name:
1792:                 - can't be blank
1793: 
1794:     ValueTooLarge:
1795:       description: Unprocessable Entity — the value is larger than the per-value size limit (64 KiB).
1796:       content:
1797:         application/json; charset=utf-8:
1798:           schema:
1799:             $ref: "#/components/schemas/ErrorDetail"
1800:           example:

→ Read docs/Makefile
Now the edit:

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -1022,8 +1022,10 @@
       "404":
         $ref: "#/components/responses/NotFound"
       "409":
         $ref: "#/components/responses/ConflictDeviceAlreadyInGroup"
+      "422":
+        $ref: "#/components/responses/ValidationErrors"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
 /appengine/v1/{realm}/groups/{group}/devices/{device}:

$ awk 'NR>=1013 && NR<=1031 {print NR": "$0}' /root/astrate-mule/docs/api/astarte_appengine_api.yaml
1013:       responses:
1014:         "201":
1015:           description: Device added to group. No body returned.
1016:         "400":
1017:           $ref: "#/components/responses/BadRequest"
1018:         "401":
1019:           $ref: "#/components/responses/Unauthorized"
1020:         "403":
1021:           $ref: "#/components/responses/Forbidden"
1022:         "404":
1023:           $ref: "#/components/responses/NotFound"
1024:         "409":
1025:           $ref: "#/components/responses/ConflictDeviceAlreadyInGroup"
1026:         "422":
1027:           $ref: "#/components/responses/ValidationErrors"
1028:         "500":
1029:           $ref: "#/components/responses/InternalServerError"
1030: 
1031:   /appengine/v1/{realm}/groups/{group}/devices/{device}:

$ command -v mkdocs python3 pip 2>&1; python3 -c "import yaml; print('pyyaml ok')" 2>&1
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; print('pyyaml ok')
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ make -C docs build 2>&1 | tail -20; echo "EXIT=$?"
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
EXIT=0

$ cat docs/embed.go; echo '--- grep yaml libs'; rg -n "yaml" go.mod; echo '--- tests touching docs yaml'; rg -ln "docs/api|appengine_api.yaml" --glob '*.go' | head
// Package docs embeds the Swagger UI static files and OpenAPI YAML specs so
// they are served directly from the astrate binary without external files.
package docs

import "embed"

// SwaggerUI holds the static Swagger UI files (index.html, CSS, JS refs).
//
//go:embed swagger-ui/*
var SwaggerUI embed.FS

// APIYAML holds the OpenAPI 3.0 YAML specifications for all five API surfaces.
//
//go:embed api/*.yaml
var APIYAML embed.FS
--- grep yaml libs
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect
--- tests touching docs yaml
internal/swagger/swagger_test.go

$ git status --porcelain; echo '--- gitignore site'; rg -n "site|site-dist" .gitignore 2>/dev/null | head
M docs/api/astarte_appengine_api.yaml
--- gitignore site
29:/site-dist/
31:# Docs site copies (refreshed by docs/make sync)
32:/docs/site/DESIGN.md
33:/docs/site/ROADMAP.md
34:/docs/site/api/
35:/docs/site/swagger-ui/

→ Read internal/swagger/swagger_test.go

$ rg -n "addGroupDevice|AddGroupDevice" internal/appengine/*_test.go | head -20; echo '--- 422 documented tests'; rg -n '"422"|ValidationErrors' internal/swagger/swagger_test.go | head
internal/appengine/groups_parity_test.go:158:func TestAddGroupDeviceErrors(t *testing.T) {
--- 422 documented tests

$ awk 'NR>=150 && NR<=205 {print NR": "$0}' internal/appengine/groups_parity_test.go
150: 	}
151: 	var members []string
152: 	decodeData(t, r.req(t, http.MethodGet, "/groups/dup/devices", "", r.token), &members)
153: 	if len(members) != 1 || members[0] != r.dev.String() {
154: 		t.Errorf("duplicated id inserted %d times, want once: %v", len(members), members)
155: 	}
156: }
157: 
158: func TestAddGroupDeviceErrors(t *testing.T) {
159: 	r := newRig(t)
160: 	newGroup(t, r, "ag", r.dev)
161: 
162: 	rec := r.req(t, http.MethodPost, "/groups/ag/devices",
163: 		`{"device_id":"not-a-device-id"}`, r.token)
164: 	if rec.Code != http.StatusUnprocessableEntity || !strings.Contains(rec.Body.String(), "is not a valid device id") {
165: 		t.Errorf("malformed device id: got %d (%s), want 422 is not a valid device id", rec.Code, rec.Body)
166: 	}
167: 
168: 	rec = r.req(t, http.MethodPost, "/groups/ag/devices",
169: 		`{"device_id":`+jsonStr(r.dev.String())+`}`, r.token)
170: 	if rec.Code != http.StatusConflict || !strings.Contains(rec.Body.String(), "Device already in group") {
171: 		t.Errorf("re-add member: got %d (%s), want 409 Device already in group", rec.Code, rec.Body)
172: 	}
173: }
174: 
175: func TestGroupListingPagination(t *testing.T) {
176: 	r := newRig(t)
177: 	ctx := context.Background()
178: 	const pg = "pg"
179: 	ids := make([]deviceid.ID, 0, 5)
180: 	for range 5 {
181: 		id, err := deviceid.Random()
182: 		if err != nil {
183: 			t.Fatal(err)
184: 		}
185: 		if err := r.st.RegisterDevice(ctx, r.realmID, id, "h"); err != nil {
186: 			t.Fatal(err)
187: 		}
188: 		ids = append(ids, id)
189: 	}
190: 	newGroup(t, r, pg, ids...)
191: 
192: 	// Without params: the whole set, bare self link, no next.
193: 	env, decode := decodeLinkedData(t, r.req(t, http.MethodGet, "/groups/"+pg+"/devices", "", r.token))
194: 	var all []string
195: 	decode(&all)
196: 	if env.Links.Self != "/v1/"+r.realm+"/groups/"+pg+"/devices" {
197: 		t.Errorf("links.self = %q, want no query string", env.Links.Self)
198: 	}
199: 	if env.Links.Next != "" {
200: 		t.Errorf("links.next = %q on an unpaged listing, want empty", env.Links.Next)
201: 	}
202: 	if !sameStringSet(all, idsToStrings(ids)) {
203: 		t.Errorf("unpaged listing = %v, want exactly the seeded set %v", all, ids)
204: 	}
205:

Now the documentation-half test, mirroring the existing ones in that file:
← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -361,8 +361,52 @@
 		t.Error("RealmCreate retention does not say 0 skips the configured default retention")
 	}
 }
 
+// TestAppEngineAddGroupDevice422Documented guards that adding a device to a
+// group documents its 422, which is a different failure from the 400: a
+// well-formed body whose device_id does not parse is rejected as a field error
+// (FieldErrors{"device_id": ...}, internal/appengine/service.go) and answered
+// 422 with the changeset-shaped FieldErrorsDetail body
+// (astarteapi.WriteFieldErrors, internal/appengine/http.go), while the
+// documented 400 only covers a body DecodeData cannot parse. So the 422 must
+// $ref the ValidationErrors response — pointing it at BadRequest would teach a
+// generated client the wrong body shape for the common case. This is the
+// documentation half of the behaviour pinned by TestAddGroupDeviceErrors in
+// internal/appengine.
+func TestAppEngineAddGroupDevice422Documented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_appengine_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_appengine_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	const (
+		key = `        "422":`
+		ref = `          $ref: "#/components/responses/ValidationErrors"`
+	)
+
+	block := operationBlock(t, lines, "addGroupDevice")
+	idx := -1
+	for i, l := range block {
+		if l == key {
+			idx = i
+			break
+		}
+	}
+	if idx < 0 {
+		t.Fatal("addGroupDevice documents no 422 response")
+	}
+	if idx+1 >= len(block) || block[idx+1] != ref {
+		t.Errorf("addGroupDevice 422 is followed by %q, want the ValidationErrors $ref", block[idx+1])
+	}
+
+	comp := strings.Join(componentBlock(t, lines, "    ValidationErrors:"), "\n")
+	if !strings.Contains(comp, `#/components/schemas/FieldErrorsDetail`) {
+		t.Error("components.responses.ValidationErrors does not carry the FieldErrorsDetail body")
+	}
+}
+
 // operationBlock returns the lines of the operation with the given operationId,
 // up to the next operation, path, or the components section.
 func operationBlock(t *testing.T, lines []string, operationID string) []string {
 	t.Helper()

$ go test ./internal/swagger/ -run 'TestAppEngine' -v 2>&1 | tail -20
=== RUN   TestAppEngine403
--- PASS: TestAppEngine403 (0.00s)
=== RUN   TestAppEngineAddGroupDevice422Documented
--- PASS: TestAppEngineAddGroupDevice422Documented (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.025s

$ cp /root/astrate-mule/docs/api/astarte_appengine_api.yaml /tmp/opencode/ae.yaml.bak && awk 'NR!=1026 && NR!=1027' /tmp/opencode/ae.yaml.bak > /root/astrate-mule/docs/api/astarte_appengine_api.yaml && awk 'NR>=1023 && NR<=1028 {print NR": "$0}' /root/astrate-mule/docs/api/astarte_appengine_api.yaml
1023:           $ref: "#/components/responses/NotFound"
1024:         "409":
1025:           $ref: "#/components/responses/ConflictDeviceAlreadyInGroup"
1026:         "500":
1027:           $ref: "#/components/responses/InternalServerError"
1028:

$ go test ./internal/swagger/ -run 'TestAppEngineAddGroupDevice422Documented' 2>&1 | tail -10; cp /tmp/opencode/ae.yaml.bak /root/astrate-mule/docs/api/astarte_appengine_api.yaml; echo '--- restored'; go test ./internal/swagger/ -run 'TestAppEngineAddGroupDevice422Documented' 2>&1 | tail -3
--- FAIL: TestAppEngineAddGroupDevice422Documented (0.00s)
    swagger_test.go:397: addGroupDevice documents no 422 response
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.010s
FAIL
--- restored
ok  	github.com/astrate-platform/astrate/internal/swagger	0.020s

$ mkdir -p /tmp/opencode/yamlprobe && cd /tmp/opencode/yamlprobe && go mod init probe 2>&1 | tail -2 && grep -E "gopkg.in/yaml.v3|go.yaml.in/yaml" /root/astrate-mule/go.sum > go.sum && go mod edit -require=gopkg.in/yaml.v3@v3.0.1 && echo ok
go: /tmp/opencode/yamlprobe/go.mod already exists
ok

← Write /tmp/opencode/yamlprobe/main.go
Wrote file successfully.

$ cd /tmp/opencode/yamlprobe && GOFLAGS=-mod=mod GOPROXY=off go run . /root/astrate-mule/docs/api/astarte_appengine_api.yaml 2>&1 | tail -20
YAML PARSE: ok
ALL LOCAL $refs: resolve
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x40 pc=0xfc960]

goroutine 1 [running]:
main.main()
	/tmp/opencode/yamlprobe/main.go:75 +0x330
exit status 2

$ cd /tmp/opencode/yamlprobe && sed -i 's/paths := lookup(&root, \[\]string{"paths"})/doc := root.Content[0]\n\tpaths := lookup(doc, []string{"paths"})/' main.go && GOFLAGS=-mod=mod GOPROXY=off go run . /root/astrate-mule/docs/api/astarte_appengine_api.yaml 2>&1 | tail -20
YAML PARSE: ok
ALL LOCAL $refs: resolve
OPERATION: addGroupDevice
RESPONSES: map[201:(inline) 400:#/components/responses/BadRequest 401:#/components/responses/Unauthorized 403:#/components/responses/Forbidden 404:#/components/responses/NotFound 409:#/components/responses/ConflictDeviceAlreadyInGroup 422:#/components/responses/ValidationErrors 500:#/components/responses/InternalServerError]

$ python3 -m pip --version 2>&1 | tail -1; python3 -c "import mkdocs" 2>&1 | tail -1; timeout 6 curl -sSI https://pypi.org/simple/ 2>&1 | head -3
/usr/bin/python3: No module named pip
ModuleNotFoundError: No module named 'mkdocs'
HTTP/2 200 
access-control-allow-origin: *
access-control-allow-methods: GET

$ rg -n "api|swagger" docs/mkdocs.yml | head -20
50:  - AppEngine API: appengine-api.md
51:  - Realm Management API: realm-management-api.md
52:  - Housekeeping API: housekeeping-api.md
63:  - API Explorer: swagger.md

$ rg -n "api|yaml|spec" docs/site/swagger.md docs/site/appengine-api.md 2>/dev/null | head -20
docs/site/appengine-api.md:24:Returns device status: introspection, connected flag, stats (total received msgs/bytes), timestamps, aliases, attributes.

$ rg -n "url|yaml|spec" docs/site/swagger.md 2>/dev/null | head; echo '--- gitignore check'; git status --porcelain
--- gitignore check
 M docs/api/astarte_appengine_api.yaml
 M internal/swagger/swagger_test.go

$ rg -n "yaml|api|url" /root/astrate-mule/docs/site/swagger.md | head -15
(no output)

$ head -40 /root/astrate-mule/docs/site/swagger.md
# API Explorer

Browse and test every Astrate REST endpoint interactively.

Select an API from the dropdown inside the explorer to switch between the five surfaces:

| API | Purpose |
|-----|---------|
| AppEngine | Device data, groups, live streams |
| Realm Management | Interfaces, triggers, configuration |
| Pairing | Device registration and credentials |
| Housekeeping | Realm lifecycle and cross-realm ops |
| Native | Internal/debug endpoints |

<iframe src="../swagger-ui/" style="width:100%;height:80vh;border:1px solid var(--md-default-fg-color--lightest);border-radius:4px;" loading="lazy"></iframe>

$ rg -n "^(test|vet|lint|check|gate):" -A3 Makefile | head -30
31:lint:
32-	@command -v golangci-lint >/dev/null 2>&1 || { \
33-		echo "golangci-lint not found; run 'make tools' (installs $(GOLANGCI_LINT_VERSION))"; exit 1; }
34-	golangci-lint run ./...
--
37:test:
38-	$(GO) test -race ./...
39-
40-## test-integration: T2 — requires a Docker daemon (or ASTRATE_TEST_DSN to reuse a database).

$ gofmt -l . ; echo "GOFMT_DONE"; go vet ./... 2>&1 | tail -5; echo "VET_EXIT=$?"
GOFMT_DONE
VET_EXIT=0

$ go test ./... 2>&1 | rg -v "^ok|no test files" | head -20; echo "TEST_DONE"
TEST_DONE

$ git diff; echo '--- todo tail'; tail -12 .mule/todo.md
diff --git a/docs/api/astarte_appengine_api.yaml b/docs/api/astarte_appengine_api.yaml
index 94adb6f..595da2d 100644
--- a/docs/api/astarte_appengine_api.yaml
+++ b/docs/api/astarte_appengine_api.yaml
@@ -1023,6 +1023,8 @@ paths:
           $ref: "#/components/responses/NotFound"
         "409":
           $ref: "#/components/responses/ConflictDeviceAlreadyInGroup"
+        "422":
+          $ref: "#/components/responses/ValidationErrors"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
diff --git a/internal/swagger/swagger_test.go b/internal/swagger/swagger_test.go
index 06eba6d..f861565 100644
--- a/internal/swagger/swagger_test.go
+++ b/internal/swagger/swagger_test.go
@@ -362,6 +362,50 @@ func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {
 	}
 }
 
+// TestAppEngineAddGroupDevice422Documented guards that adding a device to a
+// group documents its 422, which is a different failure from the 400: a
+// well-formed body whose device_id does not parse is rejected as a field error
+// (FieldErrors{"device_id": ...}, internal/appengine/service.go) and answered
+// 422 with the changeset-shaped FieldErrorsDetail body
+// (astarteapi.WriteFieldErrors, internal/appengine/http.go), while the
+// documented 400 only covers a body DecodeData cannot parse. So the 422 must
+// $ref the ValidationErrors response — pointing it at BadRequest would teach a
+// generated client the wrong body shape for the common case. This is the
+// documentation half of the behaviour pinned by TestAddGroupDeviceErrors in
+// internal/appengine.
+func TestAppEngineAddGroupDevice422Documented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_appengine_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_appengine_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	const (
+		key = `        "422":`
+		ref = `          $ref: "#/components/responses/ValidationErrors"`
+	)
+
+	block := operationBlock(t, lines, "addGroupDevice")
+	idx := -1
+	for i, l := range block {
+		if l == key {
+			idx = i
+			break
+		}
+	}
+	if idx < 0 {
+		t.Fatal("addGroupDevice documents no 422 response")
+	}
+	if idx+1 >= len(block) || block[idx+1] != ref {
+		t.Errorf("addGroupDevice 422 is followed by %q, want the ValidationErrors $ref", block[idx+1])
+	}
+
+	comp := strings.Join(componentBlock(t, lines, "    ValidationErrors:"), "\n")
+	if !strings.Contains(comp, `#/components/schemas/FieldErrorsDetail`) {
+		t.Error("components.responses.ValidationErrors does not carry the FieldErrorsDetail body")
+	}
+}
+
 // operationBlock returns the lines of the operation with the given operationId,
 // up to the next operation, path, or the components section.
 func operationBlock(t *testing.T, lines []string, operationID string) []string {
--- todo tail
- [x] verify-rm-device-delete-async-leg [legion] [auto]: run TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go) on the Legion Go — the device-delete leg added by docs-sync-rm-delete-device-async-operation-param (register a device, `DELETE /devices/{id}?async_operation=…` answers 204 and the row is gone on both values) compiles clean under `go vet -tags integration ./internal/realm/` but has never been executed, because the Pi has no database; until it runs, TestRealmManagementAsyncOperationParamDocumented in internal/swagger/swagger_test.go is the only assertion of that leg that has actually been verified.
- [!] docs-sync-rm-delete-device-async-operation-param [auto]: the `deleteDevice` operation in docs/api/astarte_realm_management_api.yaml still declares no `?async_operation` parameter although its own prose calls the deviation out ("Synchronously deletes a device and all its data. This is a deliberate deviation from upstream's async deletion", yaml:684-687) — `deleteDevice` (internal/realm/http.go:105-116) never reads the query string, so the value is accepted and ignored exactly like the four operations fixed by docs-sync-rm-async-operation-param and the now-existing `AsyncOperation` component parameter can be `$ref`d with no new component. Two things to decide first: deviation 17 in docs/COMPATIBILITY.md:335-347 lists only realm create/delete, interface install/update/delete and policy delete, so the device leg is missing from the prose there too, and TestRealmManagementAsyncOperationParam (internal/realm/async_operation_test.go:17) never sends the parameter on a device delete, so a docs-half assertion alone would be the only guard — add a case for it too (a DELETE .../devices/{device}?async_operation=true answers 204). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] forward-static-header-validation [auto]: `New` (internal/engine/forward/http.go:38-71) validates the URL and the method at boot precisely because of the rule its own comment states at 42-45 — "an unusable endpoint must fail here rather than surface once per delivery" — but `cfg.StaticHeaders` is stored unvalidated (http.go:68) and applied per delivery with `req.Header.Set` (http.go:110-112), and net/http checks header names and values only at write time. Measured on the real net/http path with a throwaway program: a name of `X Bad Name` or `X:Foo` fails per delivery with `net/http: invalid header field name`, and a value containing a newline fails with `net/http: invalid header field value for "X-Foo"` (a plain space in a value is legal). So one typo in `triggers.forward.static_headers` makes the process boot clean, accept triggers, and then count `astrate_engine_trigger_deliveries_total{outcome="failed"}` and log "custom trigger action failed" (internal/engine/triggers/actions.go:499) on every custom action forever. Validate in `New`: the RFC 7230 field-name grammar is exactly the token grammar `validMethod` (http.go:77-95) already implements, so reuse that predicate per name (rename it if you like) and add a CR/LF rejection per value. Add a table case to `TestNewRejectsBadConfig` (internal/engine/forward/http_test.go:271) asserting `New` returns an error for a bad name and a bad value and that no request is attempted; this test fails today because `New` returns nil error for both.
- [x] forward-static-headers-override-test [auto]: pin the rule http.go:22 states and http.go:110-112 implements — static headers are "applied after the fixed ones", so a static header **overwrites** the fixed `Content-Type` / `Astarte-Realm` / `Astrate-Trigger-Name` set at 107-109. Measured: a static `astarte-realm: spoofed` reaches the bus as `Astarte-Realm="spoofed"`. `TestStaticHeaders` (internal/engine/forward/http_test.go:86-118) asserts only the non-colliding case, so moving the loop above the three fixed `Set`s, or filtering reserved names out of `h.static`, leaves the whole suite green. Add a case that collides on `Astarte-Realm` (the realm-routing header a bus would filter on) and assert the server sees the static value. If the answer is instead that the fixed headers must win, flip the assertion and change 110-112 to skip the three reserved names — either way the current behaviour stops being carried by a comment alone.
- [x] forward-status-error-body [auto]: a non-2xx forward returns `forward: status %d` and nothing else (internal/engine/forward/http.go:121-123) after copying the **entire** response body into `io.Discard` with no limit (http.go:117-120) — the one thing an operator needs to debug a 500 from the bus is destroyed, and the sibling webhook request in the same codebase does bound its drain with `io.LimitReader(resp.Body, 1<<20)` (internal/engine/triggers/actions.go:626), so the two near-identical request paths have drifted on exactly that line. Read a bounded prefix (e.g. 512 bytes) of the body into the returned error and bound the discard the way actions.go:626 does. Add a case beside `TestStatusTable` (internal/engine/forward/http_test.go:164) where the handler answers 500 with body `{"error":"boom"}` and assert the error mentions both `500` and `boom`; it fails today because the error carries only the status code.
- [x] forward-envelope-bytes-test [auto]: `marshalEnvelope` (internal/engine/forward/envelope.go:19-29) has no direct test — `rg marshalEnvelope` finds exactly two callers (http.go:99, nats.go:45) and no test file — while envelope.go:6-8 makes the strong claim that "every Forwarder implementation in this package must produce byte-identical envelopes for the same inputs". Both suites check the shape only by unmarshalling into `bodyShape` (http_test.go:24-29), which cannot see key order or the exact bytes and would pass unchanged if a field were renamed on both sides. Add internal/engine/forward/envelope_test.go with a table pinning the exact output for (nil, nil) → `{"realm":"r","trigger":"t","action":null,"event":null}`, for the empty-non-nil pair the rule at 16-18 is about, and for a valued pair; container-free, no HTTP server needed.
- [x] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.
- [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml — the string `"403"` appears nowhere in the file and `components` defines no `Forbidden` response, yet every route is wrapped by `a.require = mw.RequireRealm(auth.ClaimAppEngine)` (internal/appengine/http.go:35), which answers `403 Forbidden` for a verified a_aea JWT whose grants do not authorize the method+path (`astarteapi.WriteForbidden`, internal/auth/middleware.go:90; `DetailForbidden`, pkg/astarteapi/envelope.go:41; pinned today by internal/appengine/http_test.go:184-185). Add a `Forbidden` response component copied from the realm-management one (docs/api/astarte_realm_management_api.yaml:1277-1285 — `ErrorDetail` schema, example `errors.detail: Forbidden`) and `$ref` it from all 32 operations, which is exactly the parity docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go:123-170) pins; add the same docs-half test for the AppEngine spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-patch-by-alias-409 [auto]: add the missing `409` response to `patchDeviceByAlias` in docs/api/astarte_appengine_api.yaml (yaml:241-291, currently 200/400/401/404/422/500) — the handler resolves the alias then calls the same `applyPatch` as the device- and group-scoped PATCH, and a rename to an alias already in use returns `ErrAliasAlreadyInUse` (internal/appengine/service.go:337) which `writeError` maps to `409 ConflictAliasInUse` (internal/appengine/http.go:650-651); the other two PATCH operations already document that response thanks to docs-sync-appengine-query-params-status and docs-sync-appengine-group-patch-status, so this is the last PATCH of the three with a hole. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml (yaml:939-980, currently 201/400/401/404/409/500) — a well-formed body whose `device_id` does not parse is rejected as a field error (`FieldErrors{"device_id": {"is not a valid device id"}}`, internal/appengine/service.go:631-634) and answered 422 with the `FieldErrorsDetail` body (`astarteapi.WriteFieldErrors`, internal/appengine/http.go:647), so `$ref` the existing `ValidationErrors` response (yaml:1707-1716) rather than `BadRequest` — the documented 400 only covers a body `DecodeData` cannot parse (internal/appengine/http.go:488-491). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data DELETE operations in docs/api/astarte_appengine_api.yaml (`deleteDataByAlias` yaml:519-547, `deleteData` yaml:756-784, `deleteDataInGroup` yaml:1285-1313 — each currently 204/401/404/405/500) — unsetting a path that matches no endpoint mapping returns `engine.ErrPathNotFound` (`UnsetServerProperty`, internal/engine/serverdata.go:241-244) which `writeError` answers as `400 "Endpoint not found"` (internal/appengine/http.go:679-680), the same 400 the PUT/POST twins already document via `BadRequest`; internal/appengine/writeerror_test.go pins the mapping. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 200/401/500) — a `from_token` that does not parse as a device id is refused with `ErrValidation: invalid cursor` (internal/appengine/service.go:150-155) which `writeError` maps to 422 with the flat `{"errors":{"detail":"invalid cursor"}}` body, i.e. the `ErrorDetail` shape and not `FieldErrorsDetail`, so `$ref` the same response the group twin uses for the identical parameter (`ValidationErrors`, yaml:934-935) only after fixing its example, or add a dedicated one. Extend the `from_token` parameter description (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -234,7 +234,7 @@
 - [x] forward-envelope-bytes-test [auto]: `marshalEnvelope` (internal/engine/forward/envelope.go:19-29) has no direct test — `rg marshalEnvelope` finds exactly two callers (http.go:99, nats.go:45) and no test file — while envelope.go:6-8 makes the strong claim that "every Forwarder implementation in this package must produce byte-identical envelopes for the same inputs". Both suites check the shape only by unmarshalling into `bodyShape` (http_test.go:24-29), which cannot see key order or the exact bytes and would pass unchanged if a field were renamed on both sides. Add internal/engine/forward/envelope_test.go with a table pinning the exact output for (nil, nil) → `{"realm":"r","trigger":"t","action":null,"event":null}`, for the empty-non-nil pair the rule at 16-18 is about, and for a valued pair; container-free, no HTTP server needed.
 - [x] webhook-static-headers-override-test [auto]: pin the same precedence rule on the webhook path that forward-static-headers-override-test pinned on the bus forwarder — `attempt` (internal/engine/triggers/actions.go:611-615) sets the fixed `Content-Type` and `Astarte-Realm` and only then applies `a.StaticHeaders`, so an action's `http_static_headers` wins on a colliding name, and the upstream blocklist that gates those names (`blockedHeaderNames`, internal/engine/triggers/validation.go:21-34) blocks `content-length`, `host` and the hop-by-hop set but **not** `content-type` or `astarte-realm`, so both collisions are reachable from a trigger definition. `TestWebhookDelivered` (internal/engine/triggers/actions_test.go:355-400) asserts only the non-colliding `X-Foo` case, so hoisting the static loop above the two fixed `Set` calls leaves the suite green. Add a case with a colliding `content-type` (or `astarte-realm`) asserting the static value arrives at the webhook server; decide deliberately whether upstream parity means "static wins" here as it does in the forwarder, and say so in the test comment.
 - [x] docs-sync-ae-forbidden-403 [auto]: add the missing `403` response to every one of the 32 AppEngine operations in docs/api/astarte_appengine_api.yaml — the string `"403"` appears nowhere in the file and `components` defines no `Forbidden` response, yet every route is wrapped by `a.require = mw.RequireRealm(auth.ClaimAppEngine)` (internal/appengine/http.go:35), which answers `403 Forbidden` for a verified a_aea JWT whose grants do not authorize the method+path (`astarteapi.WriteForbidden`, internal/auth/middleware.go:90; `DetailForbidden`, pkg/astarteapi/envelope.go:41; pinned today by internal/appengine/http_test.go:184-185). Add a `Forbidden` response component copied from the realm-management one (docs/api/astarte_realm_management_api.yaml:1277-1285 — `ErrorDetail` schema, example `errors.detail: Forbidden`) and `$ref` it from all 32 operations, which is exactly the parity docs-sync-rm-auth-403 already established for the realm-management spec and that TestRealmManagement403 (internal/swagger/swagger_test.go:123-170) pins; add the same docs-half test for the AppEngine spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-ae-patch-by-alias-409 [auto]: add the missing `409` response to `patchDeviceByAlias` in docs/api/astarte_appengine_api.yaml (yaml:241-291, currently 200/400/401/404/422/500) — the handler resolves the alias then calls the same `applyPatch` as the device- and group-scoped PATCH, and a rename to an alias already in use returns `ErrAliasAlreadyInUse` (internal/appengine/service.go:337) which `writeError` maps to `409 ConflictAliasInUse` (internal/appengine/http.go:650-651); the other two PATCH operations already document that response thanks to docs-sync-appengine-query-params-status and docs-sync-appengine-group-patch-status, so this is the last PATCH of the three with a hole. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
-- [ ] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml (yaml:939-980, currently 201/400/401/404/409/500) — a well-formed body whose `device_id` does not parse is rejected as a field error (`FieldErrors{"device_id": {"is not a valid device id"}}`, internal/appengine/service.go:631-634) and answered 422 with the `FieldErrorsDetail` body (`astarteapi.WriteFieldErrors`, internal/appengine/http.go:647), so `$ref` the existing `ValidationErrors` response (yaml:1707-1716) rather than `BadRequest` — the documented 400 only covers a body `DecodeData` cannot parse (internal/appengine/http.go:488-491). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [x] docs-sync-ae-add-group-device-422 [auto]: add the missing `422` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml (yaml:939-980, currently 201/400/401/404/409/500) — a well-formed body whose `device_id` does not parse is rejected as a field error (`FieldErrors{"device_id": {"is not a valid device id"}}`, internal/appengine/service.go:631-634) and answered 422 with the `FieldErrorsDetail` body (`astarteapi.WriteFieldErrors`, internal/appengine/http.go:647), so `$ref` the existing `ValidationErrors` response (yaml:1707-1716) rather than `BadRequest` — the documented 400 only covers a body `DecodeData` cannot parse (internal/appengine/http.go:488-491). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data DELETE operations in docs/api/astarte_appengine_api.yaml (`deleteDataByAlias` yaml:519-547, `deleteData` yaml:756-784, `deleteDataInGroup` yaml:1285-1313 — each currently 204/401/404/405/500) — unsetting a path that matches no endpoint mapping returns `engine.ErrPathNotFound` (`UnsetServerProperty`, internal/engine/serverdata.go:241-244) which `writeError` answers as `400 "Endpoint not found"` (internal/appengine/http.go:679-680), the same 400 the PUT/POST twins already document via `BadRequest`; internal/appengine/writeerror_test.go pins the mapping. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [ ] docs-sync-ae-list-devices-422 [auto]: add the missing `422` response to `GET /appengine/v1/{realm}/devices` in docs/api/astarte_appengine_api.yaml (yaml:27-96, currently 200/401/500) — a `from_token` that does not parse as a device id is refused with `ErrValidation: invalid cursor` (internal/appengine/service.go:150-155) which `writeError` maps to 422 with the flat `{"errors":{"detail":"invalid cursor"}}` body, i.e. the `ErrorDetail` shape and not `FieldErrorsDetail`, so `$ref` the same response the group twin uses for the identical parameter (`ValidationErrors`, yaml:934-935) only after fixing its example, or add a dedicated one. Extend the `from_token` parameter description (yaml:44-48) to say the value can 422; do NOT add `limit` to the 422 story, a non-numeric `limit` silently falls back to the default because the `strconv.Atoi` error is discarded (internal/appengine/http.go:87). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

# Todos
[✓] Read yaml:939-980 and the ValidationErrors response at yaml:1707-1716
[✓] Verify the 422 behaviour claim in Go source (service.go:631-634, http.go:647)
[✓] Add the 422 response entry to the POST group-devices path
[✓] Run `make -C docs build` and confirm the spec still loads
[✓] Run gate: go test ./... and gofmt -l .

Done: documented the missing 422 on `addGroupDevice`, `$ref`ing `ValidationErrors` (FieldErrorsDetail), with a docs-half test that fails without the change.
Files: docs/api/astarte_appengine_api.yaml (+2 lines at 1026-1027), internal/swagger/swagger_test.go (+`TestAppEngineAddGroupDevice422Documented`), .mule/todo.md (line marked `[x]`)
Verified: `go test ./internal/swagger -run TestAppEngineAddGroupDevice422` -> pass, and FAIL ("documents no 422 response") with the two YAML lines removed; `go vet ./...` -> pass; `go test ./...` (no `-race`) -> pass; `gofmt -l .` -> empty; `make -C docs build` -> `sync` half ok, then **fails: `mkdocs: No such file or directory`** (no mkdocs and no pip on this box). YAML-load check done instead with a throwaway Go probe in /tmp using gopkg.in/yaml.v3: parses ok, every local `$ref` in the spec resolves, and `addGroupDevice` responses read 201/400/401/403/404/409/**422→ValidationErrors**/500. `mkdocs build` itself never parses the spec (docs/site/swagger.md is just an iframe to swagger-ui, which fetches `/api/*.yaml` in the browser), so the probe covers what the build would have caught. Wire behaviour was confirmed by reading, not guessed: service.go:633 returns `FieldErrors{"device_id": ...}` and http.go:645-647 maps `FieldErrors` to 422; the 400 at http.go:488-491 is `DecodeData` only.
Unsure: whether marking the queue line `[x]` in .mule/todo.md is the mule's job or the script's — the file's existing completed lines are all `[x]`, so I followed that; revert if the script does it. Also worth Giulio's attention: every `docs-sync-*` task ends with "Run `make -C docs build`", and that target cannot succeed on the Pi.
Follow-ups: none
