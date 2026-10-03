slug: docs-sync-hk-error-detail-examples-split
verdict: done
at:  23006d4
ran: 2026-10-03T17:05:19Z on DietPi in 114s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ rg -n "detail: " docs/api/astarte_housekeeping_api.yaml
362:            detail: Bad Request
385:              detail: Bad Request
395:              detail: Unauthorized
405:              detail: Forbidden
415:              detail: Not Found
425:              detail: Realm deletion disabled
457:              detail: Realm already exists
467:              detail: "realm_name can't be blank"
477:              detail: Internal Server Error

$ rg -n "TestNativeErrorDetailExamples" -A 60 internal/swagger/swagger_test.go | head -100
613:// TestNativeErrorDetailExamples pins the error-detail examples in the native
614-// spec's components.responses to the frozen canonical strings the wire emits
615-// (astarteapi/envelope.go), the same guard the pairing, realm management and
616-// app engine specs get from TestPairingErrorDetailExamples,
617-// TestRealmManagementErrorDetailExamples and TestAppEngineErrorDetailExamples.
618-// Only the 500 was wrong: it showed "Internal Server Error", the capital-S
619-// spelling reconstructed rather than measured, while every 500 on the two
620-// websocket operations the component is $ref'd from answers through
621-// astarteapi.WriteInternalServerError (auth.Middleware's upgrade failure,
622-// internal/auth/middleware.go, and the channels socket,
623-// internal/appengine/channels/ws.go) with the canonical "Internal server error".
624-// The 401 and 403 are already the canonical spellings. This spec has no
625-// BadRequest component at all: the native surface is the compat health and
626-// version endpoints plus the two sockets, and none of them answers 400 through
627-// WriteBadRequest, so there is nothing here to pin.
628:func TestNativeErrorDetailExamples(t *testing.T) {
629-	b, err := docs.APIYAML.ReadFile("api/astrate_native_api.yaml")
630-	if err != nil {
631-		t.Fatalf("reading astrate_native_api.yaml: %v", err)
632-	}
633-	lines := strings.Split(string(b), "\n")
634-
635-	const detailPrefix = "              detail: "
636-
637-	for _, tc := range []struct {
638-		response string
639-		want     string
640-	}{
641-		{"    Unauthorized:", astarteapi.DetailUnauthorized},
642-		{"    Forbidden:", astarteapi.DetailForbidden},
643-		{"    InternalServerError:", astarteapi.DetailInternalServerError},
644-	} {
645-		block := componentBlock(t, lines, tc.response)
646-
647-		got := ""
648-		for _, l := range block {
649-			if strings.HasPrefix(l, detailPrefix) {
650-				got = strings.TrimPrefix(l, detailPrefix)
651-				break
652-			}
653-		}
654-		if got == "" {
655-			t.Errorf("response %q carries no %q example", tc.response, "detail")
656-			continue
657-		}
658-		if got != tc.want {
659-			t.Errorf("response %q example detail = %q, want the canonical %q",
660-				strings.TrimSpace(tc.response), got, tc.want)
661-		}
662-	}
663-}
664-
665-// TestRealmManagement422ValidationDetails pins the 422 examples the three
666-// non-interface realm-management operations carry to the details the wire
667-// actually emits, so a generated client that copies one does not learn a
668-// message no handler can produce.
669-//
670-// The three operations used to share components.responses.ValidationError, whose
671-// single example read `invalid interface: missing "interface_name"`. That string
672-// is wrong for all three: it names an interface, and none of them installs or
673-// updates one. Each now carries its own inline 422, the pattern
674-// installInterface and updateInterface already use for their own details.
675-//
676-// The createPolicy examples are derived from the compiler here instead of being
677-// written out, so a change to the message fails here instead of drifting the
678-// spec — the coupling TestPairingDeviceIDEncodingDocumented gets from
679-// deviceid. The putAuthConfig and deletePolicy details have no exported
680-// constant behind them (Service.SetAuthKey and Service.DeletePolicy,
681-// internal/realm/service.go), so they are spelled out and cited here.
682-func TestRealmManagement422ValidationDetails(t *testing.T) {
683-	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
684-	if err != nil {
685-		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
686-	}
687-	lines := strings.Split(string(b), "\n")
688-

→ Read docs/api/astarte_housekeeping_api.yaml [offset=340, limit=150]

$ rg -n "func Test(RealmManagement|AppEngine|Pairing)ErrorDetailExamples" -B 30 internal/swagger/swagger_test.go | head -120
440-		for i, l := range block {
441-			if l == key {
442-				idx = i
443-				break
444-			}
445-		}
446-		if idx < 0 {
447-			t.Errorf("operation %s documents no 400 response", op)
448-			continue
449-		}
450-		if idx+1 >= len(block) || block[idx+1] != ref {
451-			t.Errorf("operation %s 400 is followed by %q, want the BadRequest $ref", op, block[idx+1])
452-		}
453-	}
454-
455-	if !containsLine(lines, "    BadRequest:") {
456-		t.Error("components.responses defines no BadRequest response")
457-	}
458-}
459-
460-// TestPairingErrorDetailExamples pins the error-detail examples in the pairing
461-// spec's components.responses to the frozen canonical strings the wire emits
462-// (astarteapi/envelope.go), so a generated client that copies an example does
463-// not learn a detail that never goes over the wire. Phoenix renders "Bad
464-// request" and "Internal server error" — the capital-R/capital-S spellings the
465-// spec used to carry were reconstructed, not measured on upstream 1.2.0, and a
466-// startectl/SDK error path matching on the canonical string would miss them.
467-// The 404 detail is deliberately absent here: pairing answers an unmatched
468-// path with DetailPageNotFound while DetailNotFound is the generic form, and
469-// the spec's own DeviceNotFound example is the one that carries a constant.
470:func TestPairingErrorDetailExamples(t *testing.T) {
--
491-		got := ""
492-		for _, l := range block {
493-			if strings.HasPrefix(l, detailPrefix) {
494-				got = strings.TrimPrefix(l, detailPrefix)
495-				break
496-			}
497-		}
498-		if got == "" {
499-			t.Errorf("response %q carries no %q example", tc.response, "detail")
500-			continue
501-		}
502-		if got != tc.want {
503-			t.Errorf("response %q example detail = %q, want the canonical %q",
504-				strings.TrimSpace(tc.response), got, tc.want)
505-		}
506-	}
507-}
508-
509-// TestRealmManagementErrorDetailExamples pins the error-detail examples in the
510-// realm management spec's components.responses to the frozen canonical strings
511-// the wire emits (astarteapi/envelope.go), the same guard the pairing spec gets
512-// from TestPairingErrorDetailExamples. The BadRequest and InternalServerError
513-// examples carried the capital-R/capital-S spellings, which upstream Phoenix
514-// never renders and no handler emits: every 400 on this surface goes through
515-// astarteapi.WriteBadRequest and every 500 through WriteInternalServerError
516-// (internal/realm/http.go), and the SDK/astartectl error paths match on the
517-// canonical strings. The Conflict example is an endpoint-specific detail with no
518-// constant behind it, so it stays free-form; the three 422s that used to share
519-// one ValidationError component are pinned by
520-// TestRealmManagement422ValidationDetails.
521:func TestRealmManagementErrorDetailExamples(t *testing.T) {
--
543-		got := ""
544-		for _, l := range block {
545-			if strings.HasPrefix(l, detailPrefix) {
546-				got = strings.TrimPrefix(l, detailPrefix)
547-				break
548-			}
549-		}
550-		if got == "" {
551-			t.Errorf("response %q carries no %q example", tc.response, "detail")
552-			continue
553-		}
554-		if got != tc.want {
555-			t.Errorf("response %q example detail = %q, want the canonical %q",
556-				strings.TrimSpace(tc.response), got, tc.want)
557-		}
558-	}
559-}
560-
561-// TestAppEngineErrorDetailExamples pins the error-detail examples in the app
562-// engine spec's components.responses to the frozen canonical strings the wire
563-// emits (astarteapi/envelope.go), the same guard the pairing and realm
564-// management specs get from TestPairingErrorDetailExamples and
565-// TestRealmManagementErrorDetailExamples. The BadRequest and InternalServerError
566-// examples carried the capital-R/capital-S spellings, which upstream Phoenix
567-// never renders and no handler emits: every 400 on this surface goes through
568-// astarteapi.WriteBadRequest and every 500 through WriteInternalServerError
569-// (internal/appengine/http.go), and the SDK/astartectl error paths match on the
570-// canonical strings. The responses whose example is an endpoint-specific detail
571-// with no constant behind it (MethodNotAllowed, ValueTooLarge,
572-// UnprocessableEntity, the three Conflicts) stay free-form and are not pinned.
573:func TestAppEngineErrorDetailExamples(t *testing.T) {

$ rg -n "DeviceNotFound|DetailDeviceNotFound" --glob '!*_test.go' . | head -30
./pkg/astarteapi/envelope.go:44:	// DetailDeviceNotFound is the canonical 404 detail for unknown devices
./pkg/astarteapi/envelope.go:46:	DetailDeviceNotFound = "Device not found"
./pkg/astarteapi/envelope.go:210:// WriteDeviceNotFound writes the canonical 404 "Device not found" envelope.
./pkg/astarteapi/envelope.go:211:func WriteDeviceNotFound(w http.ResponseWriter) error {
./pkg/astarteapi/envelope.go:212:	return WriteError(w, http.StatusNotFound, DetailDeviceNotFound)
./docs/api/astarte_pairing_api.yaml:141:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_pairing_api.yaml:179:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_pairing_api.yaml:608:    DeviceNotFound:
./docs/api/astarte_appengine_api.yaml:163:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:210:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:246:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:298:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:336:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:381:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:433:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:483:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:535:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:578:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:608:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:650:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:695:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:740:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:787:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:825:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_appengine_api.yaml:1764:    DeviceNotFound:
./docs/api/astarte_realm_management_api.yaml:717:          $ref: "#/components/responses/DeviceNotFound"
./docs/api/astarte_realm_management_api.yaml:1349:    DeviceNotFound:
./internal/appengine/http.go:682:		_ = astarteapi.WriteDeviceNotFound(w)
./internal/realm/http.go:108:		_ = astarteapi.WriteDeviceNotFound(w)
./internal/pairing/http.go:327:		_ = astarteapi.WriteDeviceNotFound(w)

$ rg -n "func componentBlock" -A 25 internal/swagger/swagger_test.go
985:func componentBlock(t *testing.T, lines []string, key string) []string {
986-	t.Helper()
987-	start := -1
988-	for i, l := range lines {
989-		if l == key {
990-			start = i
991-			break
992-		}
993-	}
994-	if start < 0 {
995-		t.Fatalf("spec declares no components entry %q", key)
996-	}
997-	for i := start + 1; i < len(lines); i++ {
998-		if strings.TrimSpace(lines[i]) == "" {
999-			continue
1000-		}
1001-		if !strings.HasPrefix(lines[i], "      ") {
1002-			return lines[start+1 : i]
1003-		}
1004-	}
1005-	return lines[start+1:]
1006-}
1007-
1008-func containsLine(lines []string, want string) bool {
1009-	for _, l := range lines {
1010-		if l == want {

$ rg -n "Device" docs/api/astarte_housekeeping_api.yaml | head -20; echo "---"; rg -n "^[ ]{4}[A-Za-z]+:" docs/api/astarte_housekeeping_api.yaml | sed -n '1,200p' | rg -A100 "responses"
---

→ Read pkg/astarteapi/envelope.go [limit=70]
→ Read internal/swagger/swagger_test.go [offset=561, limit=55]

$ rg -n "housekeeping" internal/swagger/swagger_test.go | head -20; echo "=== docs Makefile ==="; cat docs/Makefile 2>/dev/null | head -40
127:// upstream-parity the housekeeping and pairing specs already document.
196:// TestHousekeepingAsyncOperationParamDocumented guards that the housekeeping
202:// TestHousekeepingAsyncOperationParam in internal/housekeeping.
204:	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
206:		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
309:// TestHousekeepingRetentionZeroFoldDocumented guards that the housekeeping spec
312:// `null || val == 0` to ClearRetention (internal/housekeeping/http.go) and
314:// (internal/housekeeping/service.go) — upstream parity measured on v1.2.0. A
320:// in internal/housekeeping.
322:	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
324:		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
=== docs Makefile ===
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

$ rg -n "WriteBadRequest|WriteInternalServerError|WriteUnauthorized|WriteForbidden|WriteError|WriteNotFound" internal/housekeeping/*.go | head -20
internal/housekeeping/http.go:65:		_ = astarteapi.WriteBadRequest(w)
internal/housekeeping/http.go:166:		_ = astarteapi.WriteBadRequest(w)
internal/housekeeping/http.go:171:		_ = astarteapi.WriteBadRequest(w)
internal/housekeeping/http.go:183:		_ = astarteapi.WriteBadRequest(w)
internal/housekeeping/http.go:236:		_ = astarteapi.WriteError(w, http.StatusMethodNotAllowed, "Realm deletion disabled")
internal/housekeeping/http.go:241:		_ = astarteapi.WriteError(w, http.StatusUnprocessableEntity, validationDetail(err))
internal/housekeeping/http.go:243:		_ = astarteapi.WriteError(w, http.StatusConflict, "Realm already exists")
internal/housekeeping/http.go:245:		_ = astarteapi.WriteNotFound(w)
internal/housekeeping/http.go:247:		_ = astarteapi.WriteInternalServerError(w)

Now the yaml edits:
← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -381,9 +381,9 @@
       schema:
         $ref: "#/components/schemas/ErrorDetail"
       example:
         errors:
-          detail: Bad Request
+          detail: Bad request
 
 Unauthorized:
   description: Unauthorized — missing or invalid JWT.
   content:

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -473,5 +473,5 @@
 schema:
   $ref: "#/components/schemas/ErrorDetail"
 example:
   errors:
-    detail: Internal Server Error
+    detail: Internal server error

Now the test, beside `TestNativeErrorDetailExamples`:
← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -661,8 +661,69 @@
 		}
 	}
 }
 
+// TestHousekeepingErrorDetailExamples pins the error-detail examples in the
+// housekeeping spec's components.responses to the frozen canonical strings the
+// wire emits (astarteapi/envelope.go), the same guard the pairing, realm
+// management, app engine and native specs get from
+// TestPairingErrorDetailExamples, TestRealmManagementErrorDetailExamples,
+// TestAppEngineErrorDetailExamples and TestNativeErrorDetailExamples. This is
+// the last spec still carrying a reconstructed capitalisation: its BadRequest
+// and InternalServerError examples read "Bad Request" and "Internal Server
+// Error", spellings no upstream Phoenix render answers with and no handler
+// emits — every 400 on this surface goes through astarteapi.WriteBadRequest
+// and the 500 through WriteInternalServerError (internal/housekeeping/http.go:
+// the writeError switch), and the SDK/astartectl error paths match on the
+// canonical strings. The 401 and 403 were already canonical. The responses
+// whose example is an endpoint-specific detail with no constant behind it
+// (MethodNotAllowed, Conflict, ValidationError) stay free-form and are not
+// pinned.
+//
+// There is no DeviceNotFound pin here, unlike the app engine and realm
+// management specs: this surface never answers a device 404, so
+// components.responses declares no such component to pin. The schema-level
+// ErrorDetail example is deliberately not pinned either — it is not a
+// components.responses entry and its shape question belongs to its own line.
+func TestHousekeepingErrorDetailExamples(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astrate_housekeeping_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	const detailPrefix = "              detail: "
+
+	for _, tc := range []struct {
+		response string
+		want     string
+	}{
+		{"    BadRequest:", astarteapi.DetailBadRequest},
+		{"    Unauthorized:", astarteapi.DetailUnauthorized},
+		{"    Forbidden:", astarteapi.DetailForbidden},
+		{"    NotFound:", astarteapi.DetailNotFound},
+		{"    InternalServerError:", astarteapi.DetailInternalServerError},
+	} {
+		block := componentBlock(t, lines, tc.response)
+
+		got := ""
+		for _, l := range block {
+			if strings.HasPrefix(l, detailPrefix) {
+				got = strings.TrimPrefix(l, detailPrefix)
+				break
+			}
+		}
+		if got == "" {
+			t.Errorf("response %q carries no %q example", tc.response, "detail")
+			continue
+		}
+		if got != tc.want {
+			t.Errorf("response %q example detail = %q, want the canonical %q",
+				strings.TrimSpace(tc.response), got, tc.want)
+		}
+	}
+}
+
 // TestRealmManagement422ValidationDetails pins the 422 examples the three
 // non-interface realm-management operations carry to the details the wire
 // actually emits, so a generated client that copies one does not learn a
 // message no handler can produce.

$ sed -n '225,250p' internal/housekeeping/http.go
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

$ gofmt -l . && go vet ./internal/swagger/ && go test ./internal/swagger/ -run 'ErrorDetailExamples' -v 2>&1 | tail -25
=== RUN   TestPairingErrorDetailExamples
--- PASS: TestPairingErrorDetailExamples (0.00s)
=== RUN   TestRealmManagementErrorDetailExamples
--- PASS: TestRealmManagementErrorDetailExamples (0.00s)
=== RUN   TestAppEngineErrorDetailExamples
--- PASS: TestAppEngineErrorDetailExamples (0.00s)
=== RUN   TestNativeErrorDetailExamples
--- PASS: TestNativeErrorDetailExamples (0.00s)
=== RUN   TestHousekeepingErrorDetailExamples
--- PASS: TestHousekeepingErrorDetailExamples (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.015s

! permission requested: external_directory (/tmp/*); auto-rejecting
✗ cp docs/api/astarte_housekeeping_api.yaml /tmp/hk.yaml.new && echo backed-up failed
Error: The user rejected permission to use this specific tool call.
