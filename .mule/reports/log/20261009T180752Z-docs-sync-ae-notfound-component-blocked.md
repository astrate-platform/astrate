slug: docs-sync-ae-notfound-component
verdict: blocked
at:  d68d111
ran: 2026-10-09T18:07:52Z on DietPi in 220s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read pkg/astarteapi/envelope.go
→ Read docs/api/astarte_appengine_api.yaml [offset=950, limit=520]
→ Read docs/api/astarte_appengine_api.yaml [limit=280, offset=1590]
→ Read internal/appengine/http.go [limit=80, offset=640]
→ Read internal/appengine/service.go [limit=30, offset=30]
→ Read docs/api/astarte_appengine_api.yaml [offset=1866]
→ Read internal/appengine/service.go [offset=495, limit=40]
✱ Grep "missingGroup|ErrGroupNotFound|WriteDeviceNotFound|GetGroupDevice|ListGroupDevices|AddGroupDevice|RemoveGroupDevice|ListInterfacesInGroup|GetDataInGroup|GetDataAtPathInGroup|PutDataInGroup|DeleteDataInGroup" in internal/appengine · 37 matches
→ Read internal/appengine/service.go [limit=100, offset=60]
→ Read internal/appengine/service.go [limit=300, offset=375]
→ Read internal/appengine/service.go [offset=730]
✱ Grep "func \(s \*Service\) resolveInGroup" in internal/appengine/service.go · 1 match
→ Read internal/appengine/service.go [offset=675, limit=60]
→ Read internal/appengine/http.go [offset=1, limit=60]
→ Read internal/appengine/http.go [limit=80, offset=60]
✱ Grep "func \(a \*API\) (getGroup|listGroupDevices|addGroupDevice|getGroupDevice|patchGroupDevice|removeGroupDevice|listInterfacesInGroup|getDataInGroup|putDataInGroup|deleteDataInGroup)\(" in internal/appengine/http.go · 10 matches
✱ Grep "func \(s \*Store\) (AddGroupDevice|RemoveGroupDevice|AddGroupDeviceTx|GetGroupByName)" in internal/store · 3 matches
✱ Grep "responses/NotFound|responses/DeviceNotFound|responses/GroupNotFound" in docs · 82 matches
→ Read internal/appengine/http.go [limit=120, offset=180]
→ Read docs

$ ls -la /root/astrate-mule/docs && echo "---MAKEFILE---" && cat /root/astrate-mule/docs/Makefile 2>/dev/null
total 184
drwxrwxr-x  6 root root  4096 Oct  9 19:43 .
drwxr-xr-x 17 root root  4096 Sep 25 19:41 ..
-rw-r--r--  1 root root 23369 Sep 25 22:05 COMPATIBILITY.md
-rw-r--r--  1 root root 50906 Sep  4 20:46 DESIGN.md
-rw-r--r--  1 root root  6699 Jul 27 05:10 JSON-PAYLOAD-PROFILE.md
-rw-r--r--  1 root root   426 Sep  4 20:46 Makefile
-rw-r--r--  1 root root  6196 Jul 27 05:40 OPERATIONS.md
-rw-r--r--  1 root root 45928 Sep  4 20:46 ROADMAP.md
-rw-r--r--  1 root root  1336 Sep  4 20:46 UPSTREAM-EXPERIMENTAL.md
drwxr-xr-x  2 root root  4096 Oct  8 21:26 api
-rw-r--r--  1 root root   435 Jul 27 22:18 embed.go
drwxr-xr-x  2 root root  4096 Sep  4 20:46 handoff
-rw-r--r--  1 root root  1763 Sep  4 20:46 mkdocs.yml
-rw-r--r--  1 root root    40 Jul 27 05:10 requirements.txt
drwxrwxr-x  4 root root  4096 Oct  9 19:48 site
drwxr-xr-x  2 root root  4096 Jul 27 05:40 swagger-ui
---MAKEFILE---
.PHONY: serve build clean install sync

install:
	pip install -r requirements.txt

sync:
	cp DESIGN.md site/DESIGN.md
	cp ROADMAP.md site/ROADMAP.md
	cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
	rm -rf site/api && cp -r api site/api
	rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui

serve: sync
	mkdocs serve -f mkdocs.yml

build: sync
	mkdocs build -f mkdocs.yml

clean:
	rm -rf ../site-dist

→ Read internal/appengine/http.go [limit=240, offset=300]
→ Read internal/store/groups.go [limit=150, offset=35]

$ cd /root/astrate-mule && rg -n "NotFound" docs/api/astarte_appengine_api.yaml | rg -v "DeviceNotFound:" 
204:          $ref: "#/components/responses/DeviceNotFound"
251:          $ref: "#/components/responses/DeviceNotFound"
287:          $ref: "#/components/responses/DeviceNotFound"
339:          $ref: "#/components/responses/DeviceNotFound"
377:          $ref: "#/components/responses/DeviceNotFound"
422:          $ref: "#/components/responses/DeviceNotFound"
474:          $ref: "#/components/responses/DeviceNotFound"
524:          $ref: "#/components/responses/DeviceNotFound"
576:          $ref: "#/components/responses/DeviceNotFound"
619:          $ref: "#/components/responses/DeviceNotFound"
649:          $ref: "#/components/responses/DeviceNotFound"
691:          $ref: "#/components/responses/DeviceNotFound"
736:          $ref: "#/components/responses/DeviceNotFound"
781:          $ref: "#/components/responses/DeviceNotFound"
828:          $ref: "#/components/responses/DeviceNotFound"
866:          $ref: "#/components/responses/DeviceNotFound"
978:          $ref: "#/components/responses/NotFound"
1026:          $ref: "#/components/responses/NotFound"
1070:          $ref: "#/components/responses/NotFound"
1104:          $ref: "#/components/responses/NotFound"
1152:          $ref: "#/components/responses/NotFound"
1179:          $ref: "#/components/responses/NotFound"
1209:          $ref: "#/components/responses/NotFound"
1250:          $ref: "#/components/responses/NotFound"
1296:          $ref: "#/components/responses/NotFound"
1342:          $ref: "#/components/responses/NotFound"
1390:          $ref: "#/components/responses/NotFound"
1429:          $ref: "#/components/responses/NotFound"
1856:    NotFound:

$ cd /root/astrate-mule && rg -l "astarte_appengine_api" --glob '*.go' --glob '!docs/site/**' && echo "---" && ls docs/api && echo "---git status---" && git status --short && git log --oneline -3
internal/swagger/swagger_test.go
internal/appengine/http.go
internal/appengine/service.go
---
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
---git status---
d68d111 mule: log docs-sync-ae-group-create-devices-required
9a2cf24 mule: docs-sync-ae-group-create-devices-required [auto]: in docs/api/astarte_appengine_api.yaml, `GroupCreate` (yaml:1652-1663) marks only `group_name` required, but `createGroup` treats `devices` as required — `groupBody.Devices *[]string` (internal/appengine/http.go:391-397) makes a MISSING key a 422 `{"errors":{"devices":["can't be blank"]}}` and an EMPTY array a 422 `{"errors":{"devices":["should have at least 1 item(s)"]}}` (http.go:405-418), and only the MISSING-key case is unreachable from the schema as written. The shared `ValidationErrors` example (yaml:1841-1844) that POST /groups refs (yaml:948) shows `group_name: can't be blank` alone. Add `devices` to `GroupCreate`'s `required` list (the 201 echo always carries it, http.go:424) and add the two device messages to the POST 422 story. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
067a8d0 mule: log docs-sync-ae-data-write-object-errors

$ cd /root/astrate-mule && git show 9a2cf24 --stat && git show 9a2cf24 -- docs/api/astarte_appengine_api.yaml | head -120
commit 9a2cf24ccabe110457475c4af6b9db232b05884c
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Fri Oct 9 19:51:14 2026 +0200

    mule: docs-sync-ae-group-create-devices-required [auto]: in docs/api/astarte_appengine_api.yaml, `GroupCreate` (yaml:1652-1663) marks only `group_name` required, but `createGroup` treats `devices` as required — `groupBody.Devices *[]string` (internal/appengine/http.go:391-397) makes a MISSING key a 422 `{"errors":{"devices":["can't be blank"]}}` and an EMPTY array a 422 `{"errors":{"devices":["should have at least 1 item(s)"]}}` (http.go:405-418), and only the MISSING-key case is unreachable from the schema as written. The shared `ValidationErrors` example (yaml:1841-1844) that POST /groups refs (yaml:948) shows `group_name: can't be blank` alone. Add `devices` to `GroupCreate`'s `required` list (the 201 echo always carries it, http.go:424) and add the two device messages to the POST 422 story. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

 docs/api/astarte_appengine_api.yaml | 25 ++++++++++++++++++++-----
 1 file changed, 20 insertions(+), 5 deletions(-)
commit 9a2cf24ccabe110457475c4af6b9db232b05884c
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Fri Oct 9 19:51:14 2026 +0200

    mule: docs-sync-ae-group-create-devices-required [auto]: in docs/api/astarte_appengine_api.yaml, `GroupCreate` (yaml:1652-1663) marks only `group_name` required, but `createGroup` treats `devices` as required — `groupBody.Devices *[]string` (internal/appengine/http.go:391-397) makes a MISSING key a 422 `{"errors":{"devices":["can't be blank"]}}` and an EMPTY array a 422 `{"errors":{"devices":["should have at least 1 item(s)"]}}` (http.go:405-418), and only the MISSING-key case is unreachable from the schema as written. The shared `ValidationErrors` example (yaml:1841-1844) that POST /groups refs (yaml:948) shows `group_name: can't be blank` alone. Add `devices` to `GroupCreate`'s `required` list (the 201 echo always carries it, http.go:424) and add the two device messages to the POST 422 story. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

diff --git a/docs/api/astarte_appengine_api.yaml b/docs/api/astarte_appengine_api.yaml
index 3d04a5f..762ffce 100644
--- a/docs/api/astarte_appengine_api.yaml
+++ b/docs/api/astarte_appengine_api.yaml
@@ -1651,7 +1651,7 @@ components:
 
     GroupCreate:
       type: object
-      required: [group_name]
+      required: [group_name, devices]
       properties:
         group_name:
           type: string
@@ -1899,10 +1899,25 @@ components:
         application/json; charset=utf-8:
           schema:
             $ref: "#/components/schemas/FieldErrorsDetail"
-          example:
-            errors:
-              group_name:
-                - can't be blank
+          examples:
+            group_name_blank:
+              summary: Blank group name
+              value:
+                errors:
+                  group_name:
+                    - can't be blank
+            devices_missing:
+              summary: Devices key omitted
+              value:
+                errors:
+                  devices:
+                    - can't be blank
+            devices_empty:
+              summary: Empty devices array
+              value:
+                errors:
+                  devices:
+                    - should have at least 1 item(s)
 
     ValueTooLarge:
       description: |

→ Read internal/swagger/swagger_test.go

$ cd /root/astrate-mule && sed -n '1303,1600p' internal/swagger/swagger_test.go
"        maxLength: " + wantLen,
	} {
		if !containsLine(param, want) {
			t.Errorf("DeviceID parameter is missing line %q", want)
		}
	}

	pattern := ""
	for _, l := range param {
		if v, ok := strings.CutPrefix(l, "        pattern: "); ok {
			pattern = strings.Trim(v, `'"`)
			break
		}
	}
	if pattern == "" {
		t.Fatal("DeviceID parameter carries no pattern")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		t.Fatalf("DeviceID pattern %q does not compile: %v", pattern, err)
	}

	if !re.MatchString("dT6hS2W9TT6LEnP25ks_lg") {
		t.Errorf("DeviceID pattern %q rejects the canonical ID dT6hS2W9TT6LEnP25ks_lg", pattern)
	}
	for _, spelling := range []string{
		"dT6hS2W9TT6LEnP25ks+lg",  // standard-alphabet '+'
		"dT6hS2W9TT6LEnP25ks/lg",  // standard-alphabet '/'
		"dT6hS2W9TT6LEnP25ks_lg=", // padded
	} {
		if _, err := deviceid.Parse(spelling); err == nil {
			t.Errorf("deviceid.Parse accepts %q; this test's premise no longer holds", spelling)
		}
		if re.MatchString(spelling) {
			t.Errorf("DeviceID pattern %q accepts %q, which deviceid.Parse rejects", pattern, spelling)
		}
	}
}

// TestHousekeepingRealmNamePatternDocumented guards that the housekeeping spec
// constrains `realm_name` the way the schema does. The realms table CHECKs its
// name column against a lowercase ASCII word (migrations/000002_metadata.up.sql),
// so `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are
// not realm names: a non-blank name that misses the pattern is rejected by the
// INSERT as a check violation, folded into store.ErrInvalidRealmName
// (internal/store/realms.go), re-reported by Service.CreateRealm as the
// ErrValidation `realm_name is invalid` (internal/housekeeping/service.go) and
// answered 422 with the flat detail body once validationDetail strips the prefix
// (internal/housekeeping/http.go). Both `RealmCreate.realm_name` and the
// `RealmName` path parameter used to be a bare `type: string`, which reads as
// any string and tells a generated client nothing about the six spellings a
// realm cannot have.
//
// The pattern is read out of the migration and compared with the one the spec
// ships, so the documentation cannot drift from the constraint, and the
// migration's own regex is held against the spellings the store suite refuses
// (internal/store/realms_test.go). Two asymmetries are pinned because a client
// reading only the pattern would get both wrong: the two 422 details partition
// the cases rather than being two spellings of one failure, and the name in the
// path is a read, where an off-pattern name is a 404 and not a 422 — the
// read-path hole the DeviceID parameter documents on the realm-management twin.
func TestHousekeepingRealmNamePatternDocumented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	sql, err := migrations.FS.ReadFile("000002_metadata.up.sql")
	if err != nil {
		t.Fatalf("reading 000002_metadata.up.sql: %v", err)
	}
	check := regexp.MustCompile(`CHECK \(name ~ '([^']+)'\)`).FindSubmatch(sql)
	if check == nil {
		t.Fatal("realms.name declares no CHECK (name ~ ...) constraint; this test's premise no longer holds")
	}
	want := string(check[1])

	re, err := regexp.Compile(want)
	if err != nil {
		t.Fatalf("realm name CHECK %q does not compile: %v", want, err)
	}
	for _, name := range []string{"test", "realm2", "a1"} {
		if !re.MatchString(name) {
			t.Errorf("realm name CHECK %q rejects the valid name %q", want, name)
		}
	}
	for _, name := range []string{
		"", "TEST", "test_realm", "test-realm", "1test", "tëst", "test realm",
		"re-alm", "re_alm",
	} {
		if re.MatchString(name) {
			t.Errorf("realm name CHECK %q accepts %q; this test's premise no longer holds", want, name)
		}
	}

	param := componentBlock(t, lines, "    RealmName:")
	if got := schemaPattern(t, param, "        "); got != want {
		t.Errorf("RealmName parameter pattern = %q, want the realm name CHECK %q", got, want)
	}
	for _, wording := range []string{astarteapi.DetailNotFound, "404"} {
		if !strings.Contains(strings.Join(param, "\n"), wording) {
			t.Errorf("RealmName parameter does not record that an off-pattern name reads as %q", wording)
		}
	}

	create := propertyBlock(t, componentBlock(t, lines, "    RealmCreate:"), "realm_name")
	if got := schemaPattern(t, create, "          "); got != want {
		t.Errorf("RealmCreate.realm_name pattern = %q, want the realm name CHECK %q", got, want)
	}

	realm := propertyBlock(t, componentBlock(t, lines, "    Realm:"), "realm_name")
	if got := schemaPattern(t, realm, "          "); got != want {
		t.Errorf("Realm.realm_name pattern = %q, want the realm name CHECK %q", got, want)
	}

	// Also check response array items pattern for realm names list
	dataEnvelopeNames := componentBlock(t, lines, "    DataEnvelopeRealmNames:")
	dataBlock := propertyBlock(t, dataEnvelopeNames, "data")
	// Look for items pattern in dataBlock - items are indented under data
	items := []string{}
	for i, l := range dataBlock {
		if strings.Contains(l, "items:") {
			// collect from i onwards in dataBlock, until we hit something at same or lesser indent relative to data properties
			for j := i; j < len(dataBlock); j++ {
				items = append(items, dataBlock[j])
			}
			break
		}
	}
	if got := schemaPattern(t, items, "            "); got != want {
		t.Errorf("DataEnvelopeRealmNames.items pattern = %q, want the realm name CHECK %q", got, want)
	}

	if !containsLine(operationBlock(t, lines, "createRealm"),
		`          $ref: "#/components/responses/ValidationError"`) {
		t.Error("createRealm 422 does not $ref the ValidationError response")
	}
	unprocessable := strings.Join(componentBlock(t, lines, "    ValidationError:"), "\n")
	for _, wording := range []string{
		"partition", "realm_name can't be blank", "realm_name is invalid",
	} {
		if !strings.Contains(unprocessable, wording) {
			t.Errorf("ValidationError response does not say %q", wording)
		}
	}
}

// schemaPattern returns the `pattern` value carried at the given indentation,
// or "" when the block declares none.
func schemaPattern(t *testing.T, lines []string, indent string) string {
	t.Helper()
	for _, l := range lines {
		if v, ok := strings.CutPrefix(l, indent+"pattern: "); ok {
			return strings.Trim(v, `'"`)
		}
	}
	return ""
}

// TestPairingRealmNamePatternDocumented guards that the pairing spec constrains
// its `realm` path parameter the way the schema does (the realms table CHECKs
// the name column, migrations/000002_metadata.up.sql) and records the split
// answer an off-pattern name gets across the operations sharing the parameter:
// the agent and device routes 401 through the middleware, getHealth 404 through
// the verbatim store lookup, and getVersion 200 for any string because it is
// mounted with no middleware and no realm lookup.
func TestPairingRealmNamePatternDocumented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	sql, err := migrations.FS.ReadFile("000002_metadata.up.sql")
	if err != nil {
		t.Fatalf("reading 000002_metadata.up.sql: %v", err)
	}
	check := regexp.MustCompile(`CHECK \(name ~ '([^']+)'\)`).FindSubmatch(sql)
	if check == nil {
		t.Fatal("realms.name declares no CHECK (name ~ ...) constraint; this test's premise no longer holds")
	}
	want := string(check[1])

	param := componentBlock(t, lines, "    RealmName:")
	if got := schemaPattern(t, param, "        "); got != want {
		t.Errorf("RealmName parameter pattern = %q, want the realm name CHECK %q", got, want)
	}
	for _, wording := range []string{"401", "404", "200", astarteapi.DetailNotFound} {
		if !strings.Contains(strings.Join(param, "\n"), wording) {
			t.Errorf("RealmName parameter does not record the %q answer an off-pattern name gets on some operation", wording)
		}
	}

	version := operationDescription(t, operationBlock(t, lines, "getVersion"))
	for _, wording := range []string{"200", "any realm string", "no middleware"} {
		if !strings.Contains(version, wording) {
			t.Errorf("getVersion description does not say %q; a reader would take the shared parameter as a guarantee this handler enforces", wording)
		}
	}
}

// TestPairingUnregisterDeviceSemanticsDocumented guards that the pairing spec
// describes `unregisterDevice` for what it does. "Removes a device from the
// realm" reads as a device-and-data deletion; store.UnregisterDevice clears the
// credential trail and flips status back to 'registered' and touches nothing
// else (internal/store/devices.go), so the row, its interfaces, its datastream
// data and its group memberships all survive, and a second DELETE still matches
// the row — a 204 no-op, not the 404 a deletion would have produced.
func TestPairingUnregisterDeviceSemanticsDocumented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	desc := operationDescription(t, operationBlock(t, lines, "unregisterDevice"))
	for _, want := range []string{
		"not a deletion",
		"registrable again",
		"interfaces",
		"datastream data",
		"group memberships",
		"credential",
		"204 no-op",
		"404 is returned only for a device that does not exist",
	} {
		if !strings.Contains(desc, want) {
			t.Errorf("unregisterDevice description does not say %q", want)
		}
	}
	if strings.Contains(desc, "Removes a device from the realm") {
		t.Error("unregisterDevice description still claims the device is removed from the realm")
	}
}

// operationDescription returns the text of the operation block's description,
// block-scalar or single-line, with the wrapping flattened so assertions do not
// depend on it.
func operationDescription(t *testing.T, block []string) string {
	t.Helper()
	var text []string
	for i, l := range block {
		if l == "      description: |" {
			for _, b := range block[i+1:] {
				if strings.TrimSpace(b) == "" {
					continue
				}
				if !strings.HasPrefix(b, "        ") {
					break
				}
				text = append(text, strings.TrimSpace(b))
			}
			return strings.Join(text, " ")
		}
		if v, ok := strings.CutPrefix(l, "      description: "); ok {
			return strings.Trim(strings.TrimSpace(v), `"'`)
		}
	}
	t.Fatal("operation declares no description")
	return ""
}

// operationBlock returns the lines of the operation with the given operationId,
// up to the next operation, path, or the components section.
func operationBlock(t *testing.T, lines []string, operationID string) []string {
	t.Helper()
	marker := "      operationId: " + operationID
	start := -1
	for i, l := range lines {
		if l == marker {
			start = i
			break
		}
	}
	if start < 0 {
		t.Fatalf("spec declares no operation %q", operationID)
	}
	endRe := regexp.MustCompile(`^      operationId: |^  \S|^components:`)
	for i := start + 1; i < len(lines); i++ {
		if endRe.MatchString(lines[i]) {
			return lines[start+1 : i]
		}
	}
	return lines[start+1:]
}

// componentBlock returns the lines of the components entry starting with the
// given key line, up to the next entry at the same or shallower indentation.
func componentBlock(t *testing.T, lines []string, key string) []string {
	t.Helper()
	start := -1
	for i, l := range lines {
		if l == key {
			start = i
			break
		}
	}

→ Read docs/embed.go

$ cd /root/astrate-mule && sed -n '1600,1720p' internal/swagger/swagger_test.go
}
	if start < 0 {
		t.Fatalf("spec declares no components entry %q", key)
	}
	for i := start + 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" {
			continue
		}
		if !strings.HasPrefix(lines[i], "      ") {
			return lines[start+1 : i]
		}
	}
	return lines[start+1:]
}

func containsLine(lines []string, want string) bool {
	for _, l := range lines {
		if l == want {
			return true
		}
	}
	return false
}

func embeddedYAMLFilenames() ([]string, error) {
	var names []string
	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".yaml") {
			names = append(names, strings.TrimPrefix(path, "api/"))
		}
		return nil
	})
	sort.Strings(names)
	return names, err
}

func get(t *testing.T, url string) string {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s: status = %d, want %d", url, resp.StatusCode, http.StatusOK)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading %s: %v", url, err)
	}
	return string(b)
}

$ cd /root/astrate-mule && cat .mule/todo.md 2>/dev/null; echo "=== for-giulio ==="; cat .mule/for-giulio.md 2>/dev/null | tail -60
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_121d913df001mT13v2Fx5QjXG4

- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362) — the second half docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail: Bad Request`, the reconstructed title-case form that `TestHousekeepingErrorDetailExamples` (internal/swagger/swagger_test.go:687) now pins against everywhere else on this spec, while the wire emits the frozen `astarteapi.DetailBadRequest` = `"Bad request"` (pkg/astarteapi/envelope.go:37, via `WriteBadRequest` on every 400 — internal/housekeeping/http.go:65, 166, 171, 183). Because the example sits on the generic `ErrorDetail` schema, which all eight error responses `$ref` (yaml:377-477), it is the example a generator falls back to for any response whose own example it cannot resolve — so it teaches `"Bad Request"` for a 404 and a 500 too. Change it to `Bad request`, or drop it as the realm-management and pairing `ErrorDetail` schemas do (neither carries one). Then extend `TestHousekeepingErrorDetailExamples`, whose comment currently says the schema-level example is "deliberately not pinned" (swagger_test.go:684-686), to pin it like the responses. One thing not to "fix": the example is nested one level in — it is the `errors` sub-schema's example, not a bare envelope — so do not re-wrap it in another `errors:`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-patch-422-field-error-examples [auto]: add the three missing PATCH 422 bodies to the `ValidationErrors` response component in docs/api/astarte_housekeeping_api.yaml (yaml:427-436, referenced once at yaml:190) — it carries a single example, `error_name: [invalid_update_parameters]` (internal/housekeeping/http.go:176-177), but `patchRealm` emits three more field-error envelopes: `jwt_public_key_pem: ["can't be blank"]` (http.go:188-189), `device_registration_limit: ["is invalid"]` (http.go:193-194) and `datastream_maximum_storage_retention: ["is invalid"]` (http.go:198-199). Add them as further `examples` on the media type (OpenAPI 3.0.3, yaml:1) or split the component per case. Also record that the *messages* differ from POST for the same conditions: create rejects a negative limit/retention with the flat `ErrValidation` detail `device_registration_limit must be non-negative` / `datastream_maximum_storage_retention must be non-negative` (internal/housekeeping/service.go:143-148, 422 `ErrorDetail` shape), while PATCH answers `is invalid` in the FieldErrors shape — and that the service's own three `ErrValidation` branches (service.go:205-213) are unreachable from REST, because `patchRealm` pre-checks the identical conditions and answers first. Say in the `RealmPatch` field descriptions that a negative value is a 422 rather than silently ignored. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [!] docs-sync-hk-wrong-type-field-400 [auto]: widen the `BadRequest` component description in docs/api/astarte_housekeeping_api.yaml (yaml:378, "malformed or missing data envelope") — measured today, a perfectly well-formed `data` envelope carrying a wrong-typed field answers **400**, not the documented 422. `astarteapi.DecodeData` into `map[string]json.RawMessage` (internal/housekeeping/http.go:170) cannot fail on a type mismatch, so the failure surfaces on the second decode into the typed body (http.go:182-184) and is answered `WriteBadRequest`. Probe output: `{"data":{"device_registration_limit":"100"}}` → `json: cannot unmarshal string into Go struct field patchBody.device_registration_limit of type int32`; `{"data":{"jwt_public_key_pem":123}}` → the string twin; `{"data":{"device_registration_limit":1.5}}` → `cannot unmarshal number 1.5`; `{"data":[]}` → `cannot unmarshal array into Go value of type patchBody`. Same on POST, which decodes straight into `realmBody` (http.go:63-67). This is unreachable from a spec-conformant client — `RealmCreate`/`RealmPatch` already type both fields — so document it as a one-line extension of the 400 description on both operations rather than inventing a 422, and do NOT re-shape the `BadRequest` example, whose `detail: Bad request` is the frozen constant `TestHousekeepingErrorDetailExamples` pins. One distinction worth stating in the same description: a wrong-typed field is *not* the `invalid_update_parameters` case — unknown keys are caught earlier by the `patchAllowedFields` sweep (http.go:174-180, 422 FieldErrors) on PATCH only, and silently dropped on POST. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `components.parameters.RealmName` and `RealmCreate.realm_name` both carry `pattern: '^[a-z][a-z0-9]*$'` and `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go) holds both against the realm name CHECK extracted from migrations/000002_metadata.up.sql — but the read side is still bare `type: string`: `Realm.realm_name` (components.schemas.Realm) and `DataEnvelopeRealmNames.items`. Decide whether a response field earns the same pattern (a generated client learns nothing new from it, but a reader diffing create against get sees one constrained and one not) and, if it does, add it to both and extend that test to hold the response-side pattern against the same CHECK. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] store-register-inhibit-preserve: RegisterDevice must not overwrite status to 'registered' when a device was inhibited before first credentials request (preserve 'inhibited' state). File: internal/store/devices.go — BLOCKED: gates failed
- [x] store-aliasvalues-self-exclusion-test: Add test for AliasValuesTaken to assert self-exclusion (same device's own alias values do not count as taken). File: internal/store/devices_test.go
- [!] store-latestindividual-empty-errnotfound: Add store-level test asserting LatestIndividual returns ErrNotFound for empty series and returns newest on multiple samples. File: internal/store/datastreams_test.go — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [!] pairing-burn-bcrypt-cost [auto]: make the dummy-secret timing burn use the configured BcryptCost — burnBcrypt compares against a package-level hash generated once at bcrypt.DefaultCost (internal/pairing/service.go:94-100, 428-430) while real auth hashes and compares at cfg.BcryptCost (service.go:204, 407), and bcrypt_cost is a documented, plumbed lever (internal/config/config.go:106, cmd/astrate/main.go:152, docs/site/configuration-reference.md:76, default 10), so with bcrypt_cost < 10 the unknown-device/empty-hash/malformed-ID failure paths each burn cost-10 work against cost-4 real compares, making failure responses ~an order of magnitude slower and a latency oracle for device existence — the inverse of the DESIGN §4.4 uniform-auth-time property. Keep one dummy hash per configured cost (generate lazily per Service at s.cfg.BcryptCost) and burn against it; add a deterministic test asserting bcrypt.Cost of the burn hash equals cfg.BcryptCost (e.g. a Service built with BcryptCost:4 must burn at cost 4), which fails today. — BLOCKED: lint failed: internal/pairing/service.go:95:7: G101: Potential hardcoded credentials (gosec)
- [x] pairing-onregistered-clock: route the registration event timestamp through the injectable clock — Register calls s.OnRegistered(realmName, hwID, time.Now()) with literal wall-clock time (internal/pairing/service.go:222) while every other timestamp in the package funnels through the s.now seam (VerifyCredentials, service.go:333), so a time-travel test cannot pin the DeviceRegisteredEvent timestamp fired via e.HandleDeviceRegistered (cmd/astrate/main.go:154). Change to s.now(), keep TestRegisterEmitsEvent.
- [x] pairing-unregister-bad-id-test: cover the malformed-device-IDs branch of Unregister in internal/pairing/service_test.go — Unregister wraps deviceid.Parse failure as store.ErrNotFound (internal/pairing/service.go:231-234) but TestUnregister (service_test.go:340) only exercises a well-formed unknown ID; assert both bad-ID shapes and the well-formed unknown-ID 404 in one table.
- [x] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.go
- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
- [x] compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (newest stable, 2026-10-05, maintenance-only: mix.exs/mix.lock version bumps plus a single CHANGELOG "astarte_data_updater_plant RPC server reliability" entry, no wire/API surface change; v1.4.0 is still rc.6-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2/v1.3.4 wording proposal there already stands and its version-reference delta was recorded on 2026-10-05 as reading "newest stable v1.3.5", so fold this in as that one-word delta rather than re-deriving the proposal.
- [!] docs-native-socket-handshake-errors [auto]: add the two WebSocket-handshake failures to both socket operations in docs/api/astrate_native_api.yaml — `/astrate/v1/{realm}/socket` responses (yaml:335-345) and `/appengine/v1/socket/websocket` (yaml:398-404) document only 101/401/403/500 and 101/401/500 — but measured today with the repo's own dependency (probe in /tmp, coder/websocket accept.go `verifyClientRequest`): a plain GET with no `Connection: Upgrade` answers **426 Upgrade Required**, `Content-Type: text/plain; charset=utf-8`, body `WebSocket protocol violation: Connection header "" does not contain Upgrade\n`, and `Sec-WebSocket-Version: 14` answers **400** `unsupported WebSocket protocol version (only 13 is supported): "14"\n`. Add both to each operation and say two things explicitly: the bodies are plain text from `http.Error`, NOT the `{"errors":{"detail":...}}` envelope every other response on this spec uses (a client that JSON-parses them breaks), and they are only reached after authentication — native: the route's guard `mw.RequireRealm(auth.ClaimChannels)` wraps `handle` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112), Phoenix: the realm/token checks at internal/appengine/channels/ws.go:51-68 run before `websocket.Accept` at :71 — and only on the WebSocket branch (`wantsSSE` is tested first at ws.go:78), so 401/403 and the 200 SSE path win. Do NOT add 405: ServeMux's method-mismatch applies equally to every GET route in every spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-native-socket-sse-exact-accept [auto]: fix the SSE negotiation prose on `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml — the description (yaml:289-291), the `transport` parameter (yaml:328-334) and the `200` response (yaml:338-339) read as if any client asking for `text/event-stream` gets SSE, but the predicate `wantsSSE` is exact string equality (internal/appengine/stream/ws.go:151-153): measured, `Accept: text/event-stream, text/plain` and `Accept: text/event-stream;q=1.0` both return false and fall through to `websocket.Accept`, where a non-upgrade client gets the 426 above; `transport` selects SSE only for the literal `sse`, and `transport=websocket` / `transport=` (empty) both mean WebSocket. Say "exactly" for the header, say any other `Accept` falls through to the upgrade attempt, and do NOT claim the server rejects out-of-enum `transport` values — the enum is a client-side contract the handler never enforces. While there, give the `200` its real media type: the handler sets `Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive` and writes one `data: {...}` frame per event followed by a blank line (ws.go:123-126, 156-162). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-socket-security-scheme [auto]: `/astrate/v1/{realm}/socket` is documented as if public in docs/api/astrate_native_api.yaml — the spec has no `components.securitySchemes` at all and root `security: []` (yaml:21), and the operation declares no `security`, so a client generated from this spec sends no `Authorization` header and is answered 401 `{"errors":{"detail":"Unauthorized"}}` by `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112, bearer token read only from the Authorization header at middleware.go:131, 178-196). Add a scheme mirroring `a_aea` (docs/api/astarte_appengine_api.yaml:1436-1443: apiKey in header `Authorization`, described as a realm JWT carrying the `a_ch` claim) and set `security` on this operation only — root `security: []` stays so health/readiness/metrics/version remain documented as unauthenticated. If you model the Phoenix twin too, its credential is `?token=` in the query (internal/appengine/channels/ws.go:53-54), i.e. apiKey `in: query`, not header. Description notes worth carrying: the scheme is case-insensitive `Bearer` with an optional colon, and an unknown realm is 401 not 404 (middleware.go:68-72). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-native-socket-event-schema [auto]: document the native socket's event payload in docs/api/astrate_native_api.yaml — currently only a fenced `incoming_data` example inside the operation description (yaml:300-310), with no `components.schemas` entry anywhere (yaml:425-478), and that example is not the general shape: `wireEvent` marks `interface`, `path` and `value` `omitempty` (internal/appengine/stream/ws.go:51-58), so measured on a verbatim copy a lifecycle event marshals to `{"event":"device_connected","realm":"test","device_id":"dev1","timestamp":"2024-05-30T13:49:57.045Z"}` — no interface/path/value at all (and `value: false`/`value: 0` ARE emitted, because omitempty on an `any` drops only nil). Add a `NativeEvent` schema (event, realm, device_id, optional interface/path/value, timestamp) referenced from the operation, state that `event` is one of the internal/engine/stream/bus.go Kind* constants (bus.go:20-26), and record that this socket drops `InterfaceMajor`, `IP`, `ErrorName` and `ErrorMetadata`, which `stream.Event` carries (bus.go:39-65) — so a `device_error` arrives here with no reason, unlike the Phoenix twin, which marshals the whole struct (internal/appengine/channels/ws.go:406-416). That last sentence documents a code gap, not a spec one: do not add the fields; if they are wanted on the wire, that is a separate code task. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1
- [x] testutil-wait-since-cursor [auto]: give AstarteDevice's message waits a cursor so a second wait for the same topic cannot re-match an older message — `WaitForMessage`/`WaitForTopic` (internal/testutil/astartedevice.go:149-170) scan the capture buffer from index 0 on every poll, so a repeated wait returns the first match again, and in-tree that makes testE2EEmptyCache's resync wait vacuous (see the companion line). Move the paho capture handler out of the `ConnectAstarteDevice` closure (astartedevice.go:61-67) into an injectable method, add `Mark() int` and a `From(start int, …)` wait variant, and add a container-free internal/testutil/astartedevice_test.go that drives two captures on one topic and asserts a wait started after the second returns the second message — it cannot even compile today, which is the failure proof. Keep every new read of `received` under the existing `d.mu`, and say in the report that the change is unverified for races (the Pi gate runs no `-race`). The consumer of this API is engine-e2e-emptycache-resync-assert.
- [ ] engine-e2e-emptycache-resync-assert [legion] [auto]: make the second wait in `testE2EEmptyCache` (internal/engine/engine_e2e_test.go:367) wait for the resync re-publish instead of re-matching the retained set captured at :362 — mark the cursor before `dev.EmptyCache(t)` (:364) and take only messages after the mark, so the "resent property value = int32(42)" assertion at :368-369 fails when the resend never happens; today it cannot fail, because `WaitForTopic` returns the pre-EmptyCache message, which carries the same value 42. Needs `-tags e2e` plus Docker, so this half is Legion-only — the container-free helper and its test land in testutil-wait-since-cursor.
- [!] testutil-control-frame-test [auto]: pin the producer/properties control frame that `SendProducerProperties` puts on the wire (internal/testutil/astartedevice.go:181-216) — `DeflateControlList`/`InflateControlList` have zero tests and every caller of the helpers sits behind an `e2e`/`integration`/conformance build tag, so nothing on the default gate asserts the 4-byte big-endian uncompressed-length prefix plus zlib framing, while the parser side is pinned in internal/engine/control_test.go:89-130 against a CPython golden frame. Add a container-free internal/testutil/astartedevice_test.go with exact-bytes, round-trip, empty-list→nil, sub-4-byte header and declared≠inflated cases; bound the inflate with `io.LimitReader(zr, int64(declared)+1)` the way internal/engine/control.go:301-315 does (astartedevice.go:205 reads with unbounded `io.ReadAll` today, and a lying header from the code under test makes the harness allocate it before failing) and assert the truncated byte count in the mismatch message so the bound itself is pinned. One extra line in internal/engine/control_test.go feeding `testutil.DeflateControlList` into `inflateProperties` is the cross-check that closes the loop. — BLOCKED: wrote nothing
- [x] engine-introspection-producer-roundtrip [auto]: pin the producer half of the docs/DESIGN.md §3.3 introspection rule on the default gate — `testutil.Introspection` (internal/testutil/astartedevice.go:93-105) is called from six places, all behind `e2e`, `integration && e2e` or test/conformance build tags, while the consumer `parseIntrospection` (internal/engine/introspection.go:99-127) is tested only against hand-written literals in internal/engine/introspection_test.go:17-71, so nothing container-free checks that the string the harness renders is the one the engine accepts or that the documented "deterministic (sorted) ordering" holds. Add a case to internal/engine/introspection_test.go that builds the string from an unsorted multi-entry map via `testutil.Introspection`, asserts `parseIntrospection` returns the same map, and pins the exact output (sorted, `name:major:minor;…`, no trailing `;`, empty map → empty string).
- [x] docs-sync-pairing-register-404-unreachable [auto]: drop the `404` `DeviceNotFound` response from `registerDevice` in docs/api/astarte_pairing_api.yaml (yaml:140-141) — it is documented but absent: the route is wrapped by `requireAgent` (internal/pairing/http.go:82-83), i.e. `mw.RequireRealm(auth.ClaimPairing)`, which resolves the realm *before* the handler and answers `401 {"errors":{"detail":"Unauthorized"}}` for an unknown one (internal/auth/middleware.go:66-72, "no existence oracle on auth failures"), so `Service.Register`'s own `GetRealmByName` (internal/pairing/service.go:183-186) can return `store.ErrNotFound` only if the realm is deleted between the middleware's lookup and the handler's — a race, not a contract. Nothing else in Register returns `store.ErrNotFound`: `RegisterDevice`'s `ErrDeviceAlreadyConfirmed` becomes 422 (service.go:209-213) and the CountDevices/SetPayloadFormatHint failures are DB errors → 500. No test covers it (internal/pairing/http_test.go asserts 404 only for the agent DELETE at :527 and for health at :595-598), and the `DeviceNotFound` component's own description — "unknown device ID" — does not match what that path would be even if it fired. The premise came from docs-sync-pairing-register-404, which read service.go without the middleware in front of it. Keep the `404` on `unregisterDevice` (reachable: unknown device, and a malformed `deviceID` via service.go:231-234) and on `getHealth`. State the reachability argument in the report; if you find a reachable path instead, keep the response and fix its description rather than deleting it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
- [!] docs-sync-pairing-body-cap-400 [auto]: document the 64 KiB request-body cap on the three body-taking pairing ops — `registerDevice`, `requestCredentials` and `verifyCredentials` all decode through `astarteapi.DecodeData(r.Body, maxBodyBytes, …)` with `maxBodyBytes = 64 << 10` (internal/pairing/http.go:17, 116, 169, 249), and any body over 65536 bytes fails with an error wrapping `ErrBodyTooLarge` (pkg/astarteapi/envelope.go:233-235) that each handler turns into a plain 400 `{"errors":{"detail":"Bad request"}}` (http.go:117-118, 170-171, 250-251). The shared `BadRequest` component (yaml:578-586) says only "malformed or missing data envelope" and no request-body schema carries a size hint, so a client is told nothing about the cap. Widen that one shared description to name the 64 KiB cap and the failure shape (a single edit covers all three operations) and say explicitly that an over-size body is 400, not 413 — do not invent a 413. One clause on `csr` (yaml:414) and `client_crt` (yaml:422) too: the cap is on the whole envelope, not on the PEM. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML). — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-pairing-realm-name-pattern [auto]: constrain `components.parameters.RealmName` in docs/api/astarte_pairing_api.yaml (yaml:359-365, a bare `type: string` described only as "The realm name.") with `pattern: '^[a-z][a-z0-9]*$'` — the realm name is CHECK-constrained to exactly that (migrations/000002_metadata.up.sql:7, read-only, never edit) and the two surfaces whose tasks already ran carry it on their twins (docs/api/astarte_housekeeping_api.yaml:246, docs/api/astrate_native_api.yaml:420 and 467). Because this one parameter is shared by all seven operations, the description must record that an off-pattern name gets three different answers, not one: the agent and device routes answer **401** (the realm is resolved by the middleware, which treats a miss as unauthenticated rather than 404 — internal/auth/middleware.go:66-72); `getHealth` answers **404** (service.Health returns store.ErrNotFound → WriteNotFound, internal/pairing/service.go:369-371, http.go:289-290); and `getVersion` answers **200** `{"data":"1.2.2"}` for any realm string whatsoever, because `GET /pairing/v1/{realm}/version` is registered with no middleware and no realm lookup (cmd/astrate/main.go:438 → observability.VersionHandler, internal/observability/compat.go:38-47) — say that in the version operation so nobody reads the shared parameter as a guarantee the handler enforces. Mirror `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go:1247) with a pairing twin only if it stays a few lines; the YAML edit is the deliverable. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
- [x] docs-sync-pairing-deviceid-path-param [auto]: give `components.parameters.DeviceID` the constraints `RegisterRequest.hw_id` already carries in docs/api/astarte_pairing_api.yaml — the parameter (yaml:367-377) is prose-only, while `hw_id` (yaml:386-398) got `pattern: '^[A-Za-z0-9_-]{22}$'` plus `minLength`/`maxLength` from docs-sync-pairing-deviceid-base64url, and `TestPairingDeviceIDEncodingDocumented` (internal/swagger/swagger_test.go:1086) pins only the two descriptions and hw_id's bounds, never the parameter's schema; the realm-management twin was constrained later by docs-sync-rm-deviceid-param and its test (swagger_test.go:1168) does hold the parameter's pattern against `deviceid.Parse`. Add the same three schema lines to the parameter, and record in its description the asymmetry the prose omits: a malformed `deviceID` is **401** on all three device-authenticated operations — `deviceid.Parse` failing inside `authenticateDevice` burns bcrypt and returns `ErrUnauthorized` (internal/pairing/service.go:381-385) → `WriteUnauthorized` (http.go:322-323), exactly as an unknown device or a wrong secret is (service.go:395-397), so the client is told its credentials are bad rather than its ID malformed — but **404** on the agent `unregisterDevice` (service.go:231-234 → http.go:326-327). Extend TestPairingDeviceIDEncodingDocumented to compile the new pattern and hold it against `deviceid.Parse` the way the realm-management test does. Do not encode the strict trailing-bits rule in the pattern; that stays open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
- [x] docs-sync-ae-data-envelope-metadata [auto]: add the `metadata` property to `DataEnvelopeObject` in docs/api/astarte_appengine_api.yaml (yaml:1711-1717) — the schema requires only `data` and types it `object`, but every `format=table` data response carries a sibling `metadata` object. Code: `astarteapi.WriteDataWithMetadata` renders `{"data": v, "metadata": m}` through `metadataEnvelope` (`Metadata map[string]any json:"metadata,omitempty"`, pkg/astarteapi/envelope.go:149-158), the `Tabular` carrier sets it (internal/appengine/data.go:59-64), and both renderers populate it — `renderIndividual` (data.go:256-268, `columns` keyed on the path's last segment + `table_header`) and `renderObject` (data.go:293-321, `columns` + sorted `table_header`). Add `metadata` as an optional object property (`columns: map<string,int>`, `table_header: array of string`) described as present only when `format=table`, on all six data GET 200s that `$ref` this schema (yaml:416, 468, 685, 730, 1244, 1290). While there, correct `data`'s type/description: it is not always an object — `format=structured` individual is `[{value,timestamp}]` (data.go:275-281) and `format=disjoint_tables` is an object keyed by path (data.go:269-274, 322-331). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-data-write-object-errors [auto]: document the two object-aggregated write rejections on the six data PUT/POST ops (putData/publishData, by-alias, in-group) in docs/api/astarte_appengine_api.yaml — each refs `BadRequest` for 400 (yaml:518, 570, 775, 822, 1336, 1384) and `ValueTooLarge` for 422 (yaml:528, 580, 785, 832, 1346, 1394), but three bodies are missing. (a) 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` from `writeBadObjectError` (internal/appengine/http.go:722-748, const `detailUnexpectedObjectKey` :711), reached by the `payload.ReasonOf(err) == payload.ReasonBadObject` case (http.go:684-691) — so the shared `ErrorDetail` schema (yaml:1741-1750) needs an optional `unexpected_keys` array, or a dedicated response. (b) 422 `{"errors":{"detail":"Missing required mapping key"}}` (const `detailMissingRequiredMapping`, http.go:713-720) from the `ReasonMissingRequired` case (http.go:692-702) — add an example beside the value-size one already on those ops. Both were implemented by appengine-unexpected-object-key/appengine-missing-required-422 (upstream master b6d46ad4, #2237) with no spec update. Do NOT invent other `payload.ReasonOf` statuses — that family is the separate blocked line appengine-payload-reason-status-map. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-group-create-devices-required [auto]: in docs/api/astarte_appengine_api.yaml, `GroupCreate` (yaml:1652-1663) marks only `group_name` required, but `createGroup` treats `devices` as required — `groupBody.Devices *[]string` (internal/appengine/http.go:391-397) makes a MISSING key a 422 `{"errors":{"devices":["can't be blank"]}}` and an EMPTY array a 422 `{"errors":{"devices":["should have at least 1 item(s)"]}}` (http.go:405-418), and only the MISSING-key case is unreachable from the schema as written. The shared `ValidationErrors` example (yaml:1841-1844) that POST /groups refs (yaml:948) shows `group_name: can't be blank` alone. Add `devices` to `GroupCreate`'s `required` list (the 201 echo always carries it, http.go:424) and add the two device messages to the POST 422 story. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-ae-notfound-component [auto]: the `NotFound` response component in docs/api/astarte_appengine_api.yaml (yaml:1795-1803, example `errors.detail: Not Found`) is `$ref`d by all twelve group-route 404s (yaml:978, 1026, 1070, 1104, 1152, 1179, 1209, 1250, 1296, 1342, 1390, 1429), yet `astarteapi.WriteNotFound`/`DetailNotFound` ("Not Found", pkg/astarteapi/envelope.go:42-43, 205-207) is emitted nowhere in internal/appengine (rg finds only `WriteDeviceNotFound`, http.go:683). A missing group answers 404 `{"errors":{"detail":"Group not found"}}` (`ErrGroupNotFound`, internal/appengine/service.go:41-42, 508-516; writeError http.go:667-668) and a missing member device answers `"Device not found"` (writeError store.ErrNotFound → WriteDeviceNotFound, http.go:682-683). Repoint/rename the component so the group routes show `Group not found` (add a `GroupNotFound` component, keep `DeviceNotFound` on the routes that can 404 for a missing member) and state which routes can emit which. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
=== for-giulio ===
(cmd/astrate/main.go:531) — both real, both yours to call, neither has a test worth
writing. See `.mule/reviews/cmd-astrate-2026-09-25.md`.
- **The mule has been idle 16h.** Filed by the dead-man's switch; see journalctl on the Pi.

## 2026-10-04 — three housekeeping drifts in `docs/site/` the recipe is not allowed to touch

Ran the docs-sync recipe over the housekeeping surface (code routes vs
`docs/api/astarte_housekeeping_api.yaml`). The spec half is queued as four lines in
`.mule/todo.md`; the prose half needs you, because `docs/site/` is on the never-touch list.

**1. `docs/site/housekeeping-api.md:30` still talks about Cassandra.**

> Cassandra-specific fields (`replication_class`, `replication_factor`, etc.) from upstream are accepted but ignored.

Astrate stores PostgreSQL — `internal/store` on pgx, schema in `migrations/*.sql`; there is no
Cassandra anywhere in the tree. The wire shape is four fields (`internal/housekeeping/http.go:46-51`),
and the comment right above it says the opposite of "accepted": *"Astrate omits the
Cassandra-specific fields (replication factor/class) upstream carries"*
(`internal/housekeeping/http.go:44-45`). From a client's side the two readings are the same
(the request succeeds, the fields do nothing), so this is vocabulary rather than behaviour —
but it is the kind of sentence that sends an operator hunting for a Cassandra setting that
cannot exist.

**2. `PATCH /housekeeping/v1/realms/{realm}` is missing from the page entirely.**
`docs/site/housekeeping-api.md` goes Create → List → Get → Delete (lines 10-50), so the route
`docs-sync-hk-patch-endpoint` added and documented in the spec
(`docs/api/astarte_housekeeping_api.yaml:142-192`) has no prose at all. The same page also
predates `datastream_maximum_storage_retention`, which the spec has carried on both `Realm` and
`RealmCreate` since `docs-sync-hk-retention-field`, so its create example (lines 12-21) shows a
three-field body that is not the whole body.

**3. Two `[housekeeping]` config keys exist; `configuration-reference.md` does not list them.**
The table at `docs/site/configuration-reference.md:79-84` has exactly two rows and the section
heading calls the block "Instance-admin keys" — but `HousekeepingConfig`
(`internal/config/config.go:109-125`) has four fields:

- `default_datastream_maximum_storage_retention` (`*int64`, env
  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`, validated as a
  non-negative integer at `config.go:286-291`) — injected into a realm at creation when the
  caller omits the field (`internal/housekeeping/service.go:149-150`). An operator who wants it
  has no way to learn it exists.
- `realm_deletion_disabled` (bool, env `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`, a
  fail-loud boolean accepting only `1|true|TRUE|True|0|false|""` and refusing to load on anything
  else, `config.go:296-306`) — this is the switch behind the spec's documented
  `405 "Realm deletion disabled"` on `DELETE /housekeeping/v1/realms/{realm}`, so it changes a
  status code an API client sees.

Both belong in that table next to the two key lists, and the heading wants widening.

Checked and *not* escalated, for the record. The config-key sweep
(`rg -o '\bASTRATE_[A-Z_]+' -N internal/ | sort -u`) found nothing else stale: the two
`jwt_public_key*` keys the page quotes at lines 69-71 are real (`config.go:112-113`), and the
`master_key_file` row's `ASTRATE_MASTER_KEY`/`ASTRATE_MASTER_KEY_FILE` fallbacks are real too
(`internal/store/crypto.go:24-27`). And the one thing that looked like a live bug is not one:
`view()` returns only `JWTPublicKeysPEM[0]` (`internal/housekeeping/service.go:277-281`) while
the realm-management surface joins *all* keys (`internal/realm/service.go:627`) — but every
writer stores a single-element array (create at `internal/store/realms.go:54-61`, PATCH at
`internal/store/realms.go:165`, `putAuthConfig` at `internal/realm/service.go:636`), so a
two-key realm is unreachable and there is nothing to document. Worth noting because it reads
like a bug on first inspection.

$ cd /root/astrate-mule && rg -n "NotFound|Group not found|Device not found" .mule/ 2>/dev/null | head -40
.mule/reviews/store-2026-09-08.md:33:endpoint, every error branch (ErrNotFound / ErrAlreadyExists /
.mule/reviews/store-2026-09-08.md:54:   yields ErrNotFound. The ErrNotFound branch is the particularly valuable
.mule/reviews/flowapi-2026-09-10.md:32:   permanent error — pipeline deleted (`resolveAndBuild` wraps store.ErrNotFound
.mule/reviews/pairing-2026-10-07.md:37:   `deviceid.Parse` failure as `store.ErrNotFound`; `TestUnregister` (service_test.go:340) only
.mule/reviews/store-2026-09-21.md:11:  payload hint, delete cascade. Tests: `devices_test.go` (Lifecycle, NotFound,
.mule/reviews/pairing-2026-09-09.md:53:- **`Unregister` maps a malformed device ID to `store.ErrNotFound`** (404)
.mule/reviews/store-2026-10-04.md:35:### 3. LatestIndividual has no store-level test asserting empty-series yields ErrNotFound
.mule/reviews/store-2026-10-04.md:36:datastreams.go:247-274 implements LatestIndividual (ordered DESC by ts with limit 1). Engine tests use a fake store; the real store path’s ErrNoRows→ErrNotFound conversion and value retrieval (both columns) aren’t directly asserted. Worth pinning: (a) empty series → ErrNotFound; (b) two samples on same series, returns newest; (c) correct value column populated. The ErrNotFound case is particularly valuable (used by trigger previous-value lookups).
.mule/reviews/httpx-2026-09-17.md:6:(NotFound wraps the mux, CORS wraps NotFound only when
.mule/reviews/httpx-2026-09-17.md:11:tested for their size. `NotFound`'s four-segment envelope (appengine /
.mule/reviews/httpx-2026-09-17.md:12:realmmanagement → `DetailRouteNotFound`, pairing / housekeeping →
.mule/reviews/httpx-2026-09-17.md:13:`DetailPageNotFound`, everything else keeps Go's plain-text 404) is pinned
.mule/reviews/httpx-2026-09-17.md:56:`NotFound` replay (which never overwrites an existing `Vary`).
.mule/reviews/httpx-2026-09-17.md:60:- **The `NotFound` envelope covers only four prefixes** (notfound.go:49-54):
.mule/reviews/httpx-2026-09-17.md:85:Verdict: `NotFound` is clean; one CORS variation-marker gap (small, additive,
.mule/tasks/issue-25.md:48:	ErrFlowNotFound       = errors.New("flow not found")
.mule/tasks/issue-25.md:102:		return ErrFlowNotFound
.mule/tasks/issue-25.md:135:		return nil, ErrFlowNotFound
.mule/tasks/issue-25.md:193:	if _, err := mgr.GetFlowStatus("flow-uuid-1"); !errors.Is(err, flow.ErrFlowNotFound) {
.mule/tasks/issue-25.md:194:		t.Fatalf("FAIL: expected ErrFlowNotFound after StopFlow, got %v", err)
.mule/tasks/issue-25.md:200:	if err := mgr.StopFlow(context.Background(), "unknown"); !errors.Is(err, flow.ErrFlowNotFound) {
.mule/tasks/issue-25.md:201:		t.Fatalf("FAIL: expected ErrFlowNotFound stopping non-existent flow, got %v", err)
.mule/tasks/issue-24.md:9:This specification provides the exact SQL migration, Store methods wrapping `store.ErrNotFound` / `store.ErrAlreadyExists`, and a **loud-failing test suite** for graph validation and CRUD.
.mule/reports/log/20261005T210409Z-recipe-docs-sync-timeout.md:67:- [x] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstrained string described only as "The device hardware ID." — `Service.DeleteDevice` runs the path value through `deviceid.Parse` (internal/realm/service.go:118-121), which accepts exactly 22 characters of unpadded base64url in strict mode and rejects padding, the standard alphabet's `+` and `/`, and non-canonical trailing bits (pkg/deviceid/deviceid.go:24-26, 33-37, 42-60); every rejection is folded into `store.ErrNotFound` (service.go:119-121) and answered `404 {"errors":{"detail":"Device not found"}}` (internal/realm/http.go:107-109), so a client that sends a well-formed standard-base64 or 23-character padded id is told the device does not exist rather than that its id is malformed. State the encoding and the 22-character length and add `pattern`/`minLength`/`maxLength`, exactly as docs-sync-pairing-deviceid-base64url did for the pairing twin. Do NOT encode the strict trailing-bits rule in the pattern: whether it is a deviation worth recording is still open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/reports/log/20261005T210409Z-recipe-docs-sync-timeout.md:73:- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/reports/log/20261005T210409Z-recipe-docs-sync-timeout.md:78:- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/reports/log/20261005T210409Z-recipe-docs-sync-timeout.md:85:- [!] store-latestindividual-empty-errnotfound: Add store-level test asserting LatestIndividual returns ErrNotFound for empty series and returns newest on multiple samples. File: internal/store/datastreams_test.go — BLOCKED: TIMEOUT after 1201s — task too big, split it
.mule/reports/log/20261005T210409Z-recipe-docs-sync-timeout.md:88:- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/reports/log/20261005T210409Z-recipe-docs-sync-timeout.md:91:- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/reports/log/20260920T191642Z-docs-sync-native-metrics-example-fake-series-done.md:60:ModuleNotFoundError: No module named 'yaml'
.mule/reports/log/20261003T180922Z-recipe-github-issues-proposed.md:53:- [x] docs-sync-ae-delete-data-400 [auto]: add the missing `400` response to the three data DELETE operations in docs/api/astarte_appengine_api.yaml (`deleteDataByAlias` yaml:519-547, `deleteData` yaml:756-784, `deleteDataInGroup` yaml:1285-1313 — each currently 204/401/404/405/500) — unsetting a path that matches no endpoint mapping returns `engine.ErrPathNotFound` (`UnsetServerProperty`, internal/engine/serverdata.go:241-244) which `writeError` answers as `400 "Endpoint not found"` (internal/appengine/http.go:679-680), the same 400 the PUT/POST twins already document via `BadRequest`; internal/appengine/writeerror_test.go pins the mapping. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/reports/log/20261003T180922Z-recipe-github-issues-proposed.md:55:- [x] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:556) and `InternalServerError` shows `detail: Internal Server Error` (yaml:638), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37, golden fixture pkg/astarteapi/testdata/error_bad_request.json) and `"Internal server error"` (envelope.go:58, testdata/error_internal.json), both measured on upstream 1.2.0 (envelope.go:33-36, 55-57). The other five detail examples on this file (`Not Found`, `Device not found`, `Unauthorized`, `Forbidden`, `Too Many Requests`) are already correct — change only those two. The identical wrong pair is in the other four specs (appengine 1732/1792, housekeeping 362/385/477, realm_management 1265/1348, native 466); leave those for a run of their surface. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
.mule/reports/log/20261003T180922Z-recipe-github-issues-proposed.md:71:- [x] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstrained string described only as "The device hardware ID." — `Service.DeleteDevice` runs the path value through `deviceid.Parse` (internal/realm/service.go:118-121), which accepts exactly 22 characters of unpadded base64url in strict mode and rejects padding, the standard alphabet's `+` and `/`, and non-canonical trailing bits (pkg/deviceid/deviceid.go:24-26, 33-37, 42-60); every rejection is folded into `store.ErrNotFound` (service.go:119-121) and answered `404 {"errors":{"detail":"Device not found"}}` (internal/realm/http.go:107-109), so a client that sends a well-formed standard-base64 or 23-character padded id is told the device does not exist rather than that its id is malformed. State the encoding and the 22-character length and add `pattern`/`minLength`/`maxLength`, exactly as docs-sync-pairing-deviceid-base64url did for the pairing twin. Do NOT encode the strict trailing-bits rule in the pattern: whether it is a deviation worth recording is still open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/reports/log/20261003T180922Z-recipe-github-issues-proposed.md:77:- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/reports/log/20261003T180922Z-recipe-github-issues-proposed.md:387:29:- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/reports/log/20260917T194954Z-recipe-docs-sync-proposed.md:190:f0d6a4d mule: docs-sync-pairing-deviceendpoints-dead-404-403 [auto]: remove the unreachable `404` DeviceNotFound response from the three device-facing ops in docs/api/astarte_pairing_api.yaml — `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials` (yaml:226-227), `GET /pairing/v1/{realm}/devices/{deviceID}` (yaml:268-269), `POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify` (yaml:339-340) — and the `403` Forbidden from the info GET (yaml:266-267). `Service.authenticateDevice` folds every failure (malformed/unknown deviceID, unknown realm, unregistered device, wrong secret) into uniform `ErrUnauthorized` (internal/pairing/service.go:380-411), so `store.ErrNotFound` never reaches `writeServiceError` (http.go:326-327) from those handlers; and `Service.Info` never returns `ErrInhibited` (inhibited devices still read info with 200 + status "inhibited", service.go:286-307, comment at 282-285). Reachable on these three are only 200/201/400/401/422/429/500 plus 403 on credentials and verify (both genuinely raise `ErrInhibited`, service.go:252-254 and 329-331). Keep 404 on the agent endpoints (register yaml:149-150, unregister yaml:176-177) where `GetRealmByName`/`deviceid.Parse` really do return `store.ErrNotFound`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
.mule/reports/log/20260917T194954Z-recipe-docs-sync-proposed.md:244:✱ Grep "func Write(Data|Error|NotFound|BadRequest|InternalServerError|FieldErrors)" in pkg/astarteapi · 8 matches
.mule/reports/log/20260917T194954Z-recipe-docs-sync-proposed.md:245:Both `policies` and `triggers` handlers route through the generic `writeError`. Let me check `WriteError`/`WriteNotFound`/`WriteData` shapes and the spec's response components.
.mule/reports/log/20260917T194954Z-recipe-docs-sync-proposed.md:249:✱ Grep "func WriteError|func WriteNotFound|func WriteInternalServerError|WriteRawErrors" in pkg/astarteapi/envelope.go · 7 matches
.mule/todo.md:117:- [x] docs-sync-pairing-register-404 [auto]: document the missing `404` (DeviceNotFound) on `POST /pairing/v1/{realm}/agent/devices` in docs/api/astarte_pairing_api.yaml — the handler returns it when `GetRealmByName` fails for an unknown realm (internal/pairing/service.go:183-186 → http.go:326-327), but the spec lists only 201/400/401/403/422/429/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

✗ WebFetch https://raw.githubusercontent.com/astarte-platform/astarte/master/apps/astarte_appengine_api/priv/static/astarte_appengine_api.yaml failed
Error: StatusCode: non 2xx status code (404 GET https://raw.githubusercontent.com/astarte-platform/astarte/master/apps/astarte_appengine_api/priv/static/astarte_appengine_api.yaml)
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cd /tmp && for u in \
"https://raw.githubusercontent.com/astarte-platform/astarte/v1.2.0/apps/astarte_appengine_api/priv/static/astarte_appengine_api.yaml" \
"https://raw.githubusercontent.com/astarte-platform/astarte/v1.3.0/apps/astarte_appengine_api/priv/static/astarte_appengine_api.yaml" \
"https://raw.githubusercontent.com/astarte-platform/astarte/master/apps/astarte_appengine_api/priv/static/astarte_appengine_api.yaml" ; do
echo "=== $u"; curl -fsSL --max-time 20 "$u" -o ae.yaml && { echo OK; rg -n "GroupNotFound|DeviceNotFound|NotFound:|Group not found|Device not found" ae.yaml | head -40; break; } || echo "fail"
done failed
Error: The user rejected permission to use this specific tool call.
