slug: docs-sync-hk-errordetail-schema-example
verdict: done
at:  5044143
ran: 2026-10-04T16:52:29Z on DietPi in 184s


> build · big-pickle

I'll read the context files first.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_housekeeping_api.yaml [offset=340, limit=150]
→ Read internal/swagger/swagger_test.go [offset=640, limit=120]
→ Read docs/api/astarte_housekeeping_api.yaml [offset=486, limit=25]

$ rg -n -A 16 '^    ErrorDetail:' /root/astrate-mule/docs/api/*.yaml
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml:368:    ErrorDetail:
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-369-      type: object
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-370-      required: [errors]
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-371-      properties:
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-372-        errors:
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-373-          type: object
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-374-          required: [detail]
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-375-          properties:
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-376-            detail:
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-377-              type: string
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-378-              description: Human-readable error detail string.
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-379-          example:
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-380-            detail: Bad Request
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-381-
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-382-    FieldErrorsDetail:
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-383-      type: object
/root/astrate-mule/docs/api/astarte_housekeeping_api.yaml-384-      required: [errors]
--
/root/astrate-mule/docs/api/astarte_pairing_api.yaml:555:    ErrorDetail:
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-556-      type: object
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-557-      required: [errors]
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-558-      properties:
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-559-        errors:
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-560-          type: object
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-561-          required: [detail]
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-562-          properties:
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-563-            detail:
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-564-              type: string
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-565-
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-566-    ErrorFields:
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-567-      type: object
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-568-      required: [errors]
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-569-      properties:
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-570-        errors:
/root/astrate-mule/docs/api/astarte_pairing_api.yaml-571-          type: object
--
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml:1237:    ErrorDetail:
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1238-      type: object
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1239-      required: [errors]
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1240-      properties:
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1241-        errors:
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1242-          type: object
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1243-          required: [detail]
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1244-          properties:
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1245-            detail:
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1246-              type: string
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1247-
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1248-    TriggerErrorsDetail:
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1249-      type: object
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1250-      required: [errors]
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1251-      properties:
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1252-        errors:
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml-1253-          type: object
--
/root/astrate-mule/docs/api/astarte_appengine_api.yaml:1700:    ErrorDetail:
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1701-      type: object
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1702-      required: [errors]
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1703-      properties:
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1704-        errors:
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1705-          type: object
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1706-          required: [detail]
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1707-          properties:
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1708-            detail:
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1709-              type: string
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1710-
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1711-    FieldErrorsDetail:
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1712-      type: object
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1713-      required: [errors]
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1714-      properties:
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1715-        errors:
/root/astrate-mule/docs/api/astarte_appengine_api.yaml-1716-          type: object
--
/root/astrate-mule/docs/api/astrate_native_api.yaml:426:    ErrorDetail:
/root/astrate-mule/docs/api/astrate_native_api.yaml-427-      type: object
/root/astrate-mule/docs/api/astrate_native_api.yaml-428-      required: [errors]
/root/astrate-mule/docs/api/astrate_native_api.yaml-429-      properties:
/root/astrate-mule/docs/api/astrate_native_api.yaml-430-        errors:
/root/astrate-mule/docs/api/astrate_native_api.yaml-431-          type: object
/root/astrate-mule/docs/api/astrate_native_api.yaml-432-          required: [detail]
/root/astrate-mule/docs/api/astrate_native_api.yaml-433-          properties:
/root/astrate-mule/docs/api/astrate_native_api.yaml-434-            detail:
/root/astrate-mule/docs/api/astrate_native_api.yaml-435-              type: string
/root/astrate-mule/docs/api/astrate_native_api.yaml-436-
/root/astrate-mule/docs/api/astrate_native_api.yaml-437-  responses:
/root/astrate-mule/docs/api/astrate_native_api.yaml-438-    Unauthorized:
/root/astrate-mule/docs/api/astrate_native_api.yaml-439-      description: Unauthorized — missing or invalid JWT/token.
/root/astrate-mule/docs/api/astrate_native_api.yaml-440-      content:
/root/astrate-mule/docs/api/astrate_native_api.yaml-441-        application/json; charset=utf-8:
/root/astrate-mule/docs/api/astrate_native_api.yaml-442-          schema:

$ rg -n 'func (componentBlock|schemaBlock|containsLine|responseBlock)' -A 40 /root/astrate-mule/internal/swagger/swagger_test.go | head -100; echo ---; rg -n 'detailPrefix' /root/astrate-mule/internal/swagger/swagger_test.go
790:func responseBlock(t *testing.T, block []string, status string) []string {
791-	t.Helper()
792-	key := `        "` + status + `":`
793-	start := -1
794-	for i, l := range block {
795-		if l == key {
796-			start = i
797-			break
798-		}
799-	}
800-	if start < 0 {
801-		t.Fatalf("operation block declares no %s response", status)
802-	}
803-	for i := start + 1; i < len(block); i++ {
804-		if strings.HasPrefix(block[i], `        "`) {
805-			return block[start+1 : i]
806-		}
807-	}
808-	return block[start+1:]
809-}
810-
811-// TestPairingDeviceIDEncodingDocumented guards that the pairing spec states the
812-// exact wire form of a device identifier — the 22-character unpadded base64url
813-// string deviceid.Parse accepts — instead of the "base64-encoded 128-bit"
814-// wording, which read as standard base64 and let a generated client send a
815-// well-formed `+`/`/` identifier or a padded 23-character spelling and collect
816-// a 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}`
817-// (internal/pairing/http.go). The shipped pattern is compiled here and held to
818-// the same spellings the parser is, so the two cannot drift apart: the bounds
819-// and the description wording are derived from deviceid.EncodedLen rather than
820-// written out, so a change to the accepted length fails this test too.
821-func TestPairingDeviceIDEncodingDocumented(t *testing.T) {
822-	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
823-	if err != nil {
824-		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
825-	}
826-	lines := strings.Split(string(b), "\n")
827-
828-	wantLen := strconv.Itoa(deviceid.EncodedLen)
829-	wantDesc := fmt.Sprintf("%s-character unpadded", wantLen)
830-
--
1145:func componentBlock(t *testing.T, lines []string, key string) []string {
1146-	t.Helper()
1147-	start := -1
1148-	for i, l := range lines {
1149-		if l == key {
1150-			start = i
1151-			break
1152-		}
1153-	}
1154-	if start < 0 {
1155-		t.Fatalf("spec declares no components entry %q", key)
1156-	}
1157-	for i := start + 1; i < len(lines); i++ {
1158-		if strings.TrimSpace(lines[i]) == "" {
1159-			continue
1160-		}
1161-		if !strings.HasPrefix(lines[i], "      ") {
1162-			return lines[start+1 : i]
1163-		}
1164-	}
1165-	return lines[start+1:]
1166-}
1167-
1168:func containsLine(lines []string, want string) bool {
1169-	for _, l := range lines {
1170-		if l == want {
1171-			return true
1172-		}
1173-	}
1174-	return false
1175-}
1176-
1177-func embeddedYAMLFilenames() ([]string, error) {
1178-	var names []string
1179-	err := fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
1180-		if err != nil {
1181-			return err
1182-		}
1183-		if d.IsDir() {
1184-			return nil
1185-		}
1186-		if strings.HasSuffix(path, ".yaml") {
1187-			names = append(names, strings.TrimPrefix(path, "api/"))
1188-		}
1189-		return nil
1190-	})
1191-	sort.Strings(names)
1192-	return names, err
1193-}
1194-
1195-func get(t *testing.T, url string) string {
1196-	t.Helper()
1197-	resp, err := http.Get(url)
1198-	if err != nil {
1199-		t.Fatalf("GET %s: %v", url, err)
1200-	}
1201-	defer resp.Body.Close()
1202-
---
478:	const detailPrefix = "              detail: "
494:			if strings.HasPrefix(l, detailPrefix) {
495:				got = strings.TrimPrefix(l, detailPrefix)
529:	const detailPrefix = "              detail: "
546:			if strings.HasPrefix(l, detailPrefix) {
547:				got = strings.TrimPrefix(l, detailPrefix)
581:	const detailPrefix = "              detail: "
598:			if strings.HasPrefix(l, detailPrefix) {
599:				got = strings.TrimPrefix(l, detailPrefix)
636:	const detailPrefix = "              detail: "
650:			if strings.HasPrefix(l, detailPrefix) {
651:				got = strings.TrimPrefix(l, detailPrefix)
695:	const detailPrefix = "              detail: "
711:			if strings.HasPrefix(l, detailPrefix) {
712:				got = strings.TrimPrefix(l, detailPrefix)

$ ls /root/astrate-mule/.mule/; echo ---; rg -n 'hk-errordetail|error-detail' /root/astrate-mule/.mule/todo.md
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
---
241:- [x] docs-sync-pairing-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_pairing_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:556) and `InternalServerError` shows `detail: Internal Server Error` (yaml:638), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37, golden fixture pkg/astarteapi/testdata/error_bad_request.json) and `"Internal server error"` (envelope.go:58, testdata/error_internal.json), both measured on upstream 1.2.0 (envelope.go:33-36, 55-57). The other five detail examples on this file (`Not Found`, `Device not found`, `Unauthorized`, `Forbidden`, `Too Many Requests`) are already correct — change only those two. The identical wrong pair is in the other four specs (appengine 1732/1792, housekeeping 362/385/477, realm_management 1265/1348, native 466); leave those for a run of their surface. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
254:- [x] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, housekeeping and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
260:- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` shows `detail: Bad Request` (yaml:1732) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1792), while the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58 via `WriteBadRequest`/`WriteInternalServerError`, golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json); every 400/500 on that surface goes through the same astarteapi writers. The other detail examples on the file are already correct — change only those two — and add the docs-half assertion `TestAppEngineErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go, which pins the pairing and realm-management halves. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
261:- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDetail` schema-level example (yaml:361-362) and the `BadRequest` response (yaml:385) show `detail: Bad Request`, and `InternalServerError` (yaml:477) shows `detail: Internal Server Error`, but the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58, via WriteBadRequest/WriteInternalServerError, pinned by pkg/astarteapi/testdata/error_bad_request.json and error_internal.json). While you are in `ErrorDetail`: its `example` sits at the schema level and shows a bare `{detail: ...}` where the envelope the wire sends is `{"errors":{"detail":...}}` — move it under `errors` (or drop it, as the realm-management and pairing `ErrorDetail` schemas carry no example) so a generated client is not shown the wrong shape. Add `TestHousekeepingErrorDetailExamples` beside the pairing and realm-management halves in internal/swagger/swagger_test.go. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
262:- [x] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `"Internal server error"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
263:- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
269:- [ ] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362) — the second half docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail: Bad Request`, the reconstructed title-case form that `TestHousekeepingErrorDetailExamples` (internal/swagger/swagger_test.go:687) now pins against everywhere else on this spec, while the wire emits the frozen `astarteapi.DetailBadRequest` = `"Bad request"` (pkg/astarteapi/envelope.go:37, via `WriteBadRequest` on every 400 — internal/housekeeping/http.go:65, 166, 171, 183). Because the example sits on the generic `ErrorDetail` schema, which all eight error responses `$ref` (yaml:377-477), it is the example a generator falls back to for any response whose own example it cannot resolve — so it teaches `"Bad Request"` for a 404 and a 500 too. Change it to `Bad request`, or drop it as the realm-management and pairing `ErrorDetail` schemas do (neither carries one). Then extend `TestHousekeepingErrorDetailExamples`, whose comment currently says the schema-level example is "deliberately not pinned" (swagger_test.go:684-686), to pin it like the responses. One thing not to "fix": the example is nested one level in — it is the `errors` sub-schema's example, not a bare envelope — so do not re-wrap it in another `errors:`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cd /root/astrate-mule && cat docs/Makefile 2>/dev/null | head -40; echo ---; rg -n 'DetailBadRequest|DetailNotFound\s*=' pkg/astarteapi/envelope.go; echo ---; rg -c '^    ErrorDetail:$' docs/api/astarte_housekeeping_api.yaml
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
---
33:	// DetailBadRequest is the canonical 400 detail. Measured upstream 1.2.0
37:	DetailBadRequest = "Bad request"
43:	DetailNotFound = "Not Found"
192:	return WriteError(w, http.StatusBadRequest, DetailBadRequest)
---
1

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -376,9 +376,9 @@
         detail:
           type: string
           description: Human-readable error detail string.
       example:
-        detail: Bad Request
+        detail: Bad request
 
 FieldErrorsDetail:
   type: object
   required: [errors]

→ Read docs/api/astarte_housekeeping_api.yaml [limit=16, offset=366]
Now the test extension:

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -681,11 +681,18 @@
 // pinned.
 //
 // There is no DeviceNotFound pin here, unlike the app engine and realm
 // management specs: this surface never answers a device 404, so
-// components.responses declares no such component to pin. The schema-level
-// ErrorDetail example is deliberately not pinned either — it is not a
-// components.responses entry and its shape question belongs to its own line.
+// components.responses declares no such component to pin.
+//
+// The schema-level ErrorDetail example is pinned too, to the same canonical
+// string as BadRequest. It sits on the generic schema all eight error
+// responses $ref, so it is the example a generator falls back to for any
+// response whose own it cannot resolve — leaving it on the reconstructed
+// "Bad Request" taught a 404 and a 500 a detail no handler emits. Its shape is
+// already the wire's: the example belongs to the `errors` sub-schema, so it is
+// the value of the `errors` property and not a bare envelope, and it must not
+// be re-wrapped in another `errors:`.
 func TestHousekeepingErrorDetailExamples(t *testing.T) {
 	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
 	if err != nil {
 		t.Fatalf("reading astrate_housekeeping_api.yaml: %v", err)

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -728,8 +728,23 @@
 			t.Errorf("response %q example detail = %q, want the canonical %q",
 				strings.TrimSpace(tc.response), got, tc.want)
 		}
 	}
+
+	schemaBlock := componentBlock(t, lines, "    ErrorDetail:")
+	const schemaDetailPrefix = "            detail: "
+
+	got := ""
+	for _, l := range schemaBlock {
+		if strings.HasPrefix(l, schemaDetailPrefix) {
+			got = strings.TrimPrefix(l, schemaDetailPrefix)
+			break
+		}
+	}
+	if got != astarteapi.DetailBadRequest {
+		t.Errorf("ErrorDetail schema example detail = %q, want the canonical %q",
+			got, astarteapi.DetailBadRequest)
+	}
 }
 
 // TestRealmManagement422ValidationDetails pins the 422 examples the three
 // non-interface realm-management operations carry to the details the wire

$ cd /root/astrate-mule && gofmt -l . && go vet ./internal/swagger/ && go test ./internal/swagger/ -run TestHousekeepingErrorDetailExamples -v 2>&1 | tail -20
=== RUN   TestHousekeepingErrorDetailExamples
--- PASS: TestHousekeepingErrorDetailExamples (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.014s

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -376,9 +376,9 @@
         detail:
           type: string
           description: Human-readable error detail string.
       example:
-        detail: Bad request
+        detail: Bad Request
 
 FieldErrorsDetail:
   type: object
   required: [errors]

$ cd /root/astrate-mule && go test ./internal/swagger/ -run TestHousekeepingErrorDetailExamples 2>&1 | tail -8
--- FAIL: TestHousekeepingErrorDetailExamples (0.00s)
    swagger_test.go:744: ErrorDetail schema example detail = "Bad Request", want the canonical "Bad request"
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.018s
FAIL

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -376,9 +376,9 @@
         detail:
           type: string
           description: Human-readable error detail string.
       example:
-        detail: Bad Request
+        detail: Bad request
 
 FieldErrorsDetail:
   type: object
   required: [errors]

Test breaks without the change. Now the docs build and the gate:
$ cd /root/astrate-mule && (python3 -c "import yaml,sys; d=yaml.safe_load(open('docs/api/astarte_housekeeping_api.yaml')); print('yaml ok, ErrorDetail example:', d['components']['schemas']['ErrorDetail']['properties']['errors']['example'])" ) 2>&1; echo "--- mkdocs:"; command -v mkdocs || echo "mkdocs not installed"
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml,sys; d=yaml.safe_load(open('docs/api/astarte_housekeeping_api.yaml')); print('yaml ok, ErrorDetail example:', d['components']['schemas']['ErrorDetail']['properties']['errors']['example'])
    ^^^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'
--- mkdocs:
mkdocs not installed

$ cd /root/astrate-mule && (pip install -q -r docs/requirements.txt 2>&1 | tail -3); command -v mkdocs || echo "still no mkdocs"
/bin/sh: 1: pip: not found
still no mkdocs

$ cd /root/astrate-mule && rg -n '^func Test' internal/swagger/swagger_test.go
24:func TestMount(t *testing.T) {
83:func TestMountSubPanicsOnBrokenFS(t *testing.T) {
96:func TestSpecs(t *testing.T) {
129:func TestRealmManagement403(t *testing.T) {
140:func TestAppEngine403(t *testing.T) {
204:func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
247:func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
322:func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {
382:func TestAppEngineAddGroupDevice422Documented(t *testing.T) {
426:func TestAppEngineDataDelete400Documented(t *testing.T) {
471:func TestPairingErrorDetailExamples(t *testing.T) {
522:func TestRealmManagementErrorDetailExamples(t *testing.T) {
574:func TestAppEngineErrorDetailExamples(t *testing.T) {
629:func TestNativeErrorDetailExamples(t *testing.T) {
695:func TestHousekeepingErrorDetailExamples(t *testing.T) {
766:func TestRealmManagement422ValidationDetails(t *testing.T) {
843:func TestPairingDeviceIDEncodingDocumented(t *testing.T) {
925:func TestRealmManagementDeviceIDEncodingDocumented(t *testing.T) {
1004:func TestHousekeepingRealmNamePatternDocumented(t *testing.T) {
1087:func TestPairingUnregisterDeviceSemanticsDocumented(t *testing.T) {

→ Read internal/swagger/swagger_test.go [limit=35, offset=96]

$ cd /root/astrate-mule && rg -n 'yaml' go.mod; echo "--- gitignore:"; rg -n 'site' .gitignore docs/.gitignore 2>/dev/null
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect
--- gitignore:
.gitignore:29:/site-dist/
.gitignore:31:# Docs site copies (refreshed by docs/make sync)
.gitignore:32:/docs/site/DESIGN.md
.gitignore:33:/docs/site/ROADMAP.md
.gitignore:34:/docs/site/api/
.gitignore:35:/docs/site/swagger-ui/

$ cat > /tmp/probe_hk_yaml.go <<'EOF'
package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	b, err := os.ReadFile("docs/api/astarte_housekeeping_api.yaml")
	if err != nil {
		fmt.Println("read:", err)
		os.Exit(1)
	}
	var doc map[string]any
	if err := yaml.Unmarshal(b, &doc); err != nil {
		fmt.Println("yaml parse FAILED:", err)
		os.Exit(1)
	}
	schemas := doc["components"].(map[string]any)["schemas"].(map[string]any)
	ed := schemas["ErrorDetail"].(map[string]any)["properties"].(map[string]any)["errors"].(map[string]any)
	ex, ok := ed["example"]
	if !ok {
		fmt.Println("ErrorDetail.errors carries no example")
		os.Exit(1)
	}
	fmt.Printf("yaml ok; ErrorDetail.errors.example = %#v\n", ex)
	// every $ref target still resolves
	resolves := 0
	for _, r := range collectRefs(doc) {
		name := r[len("#/components/schemas/"):]
		if _, ok := schemas[name]; !ok {
			fmt.Println("DANGLING $ref:", r)
			os.Exit(1)
		}
		resolves++
	}
	fmt.Println("yaml ok; all", resolves, "schema $refs resolve")
}

func collectRefs(v any) []string {
	var out []string
	switch t := v.(type) {
	case map[string]any:
		for k, vv := range t {
			if s, ok := vv.(string); ok && k == "$ref" && len(s) > len("#/components/schemas/") {
				out = append(out, s)
			}
			out = append(out, collectRefs(vv)...)
		}
	case []any:
		for _, vv := range t {
			out = append(out, collectRefs(vv)...)
		}
	}
	return out
}
EOF
cd /root/astrate-mule && go run /tmp/probe_hk_yaml.go
yaml ok; ErrorDetail.errors.example = map[string]interface {}{"detail":"Bad request"}
DANGLING $ref: #/components/parameters/RealmName
exit status 1

! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cd /tmp && python3 - <<'EOF'
import re
s = open('/tmp/probe_hk_yaml.go').read()
s = s.replace('''			if s, ok := vv.(string); ok && k == "$ref" && len(s) > len("#/components/schemas/") {
				out = append(out, s)
			}''', '''			if s, ok := vv.(string); ok && k == "$ref" && strings.HasPrefix(s, "#/components/schemas/") {
				out = append(out, s)
			}
			if s, ok := vv.(string); ok && k == "$ref" && strings.HasPrefix(s, "#/components/parameters/") {
				out = append(out, s)
			}''')
s = s.replace('''	resolves := 0
	for _, r := range collectRefs(doc) {
		name := r[len("#/components/schemas/"):]
		if _, ok := schemas[name]; !ok {
			fmt.Println("DANGLING $ref:", r)
			os.Exit(1)
		}
		resolves++
	}''', '''	params := doc["components"].(map[string]any)["parameters"].(map[string]any)
	resolves := 0
	for _, r := range collectRefs(doc) {
		var ok bool
		if strings.HasPrefix(r, "#/components/schemas/") {
			_, ok = schemas[r[len("#/components/schemas/"):]]
		} else {
			_, ok = params[r[len("#/components/parameters/"):]]
		}
		if !ok {
			fmt.Println("DANGLING $ref:", r)
			os.Exit(1)
		}
		resolves++
	}''')
s = s.replace('''	"fmt"
	"os"
''', '''	"fmt"
	"os"
	"strings"
''')
open('/tmp/probe_hk_yaml.go','w').write(s)
EOF
cd /root/astrate-mule && go run /tmp/probe_hk_yaml.go failed
Error: The user rejected permission to use this specific tool call.
