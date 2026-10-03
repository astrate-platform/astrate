package swagger

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	docs "github.com/astrate-platform/astrate/docs"
	"github.com/astrate-platform/astrate/internal/engine/triggers"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
	"github.com/astrate-platform/astrate/pkg/deviceid"
)

func TestMount(t *testing.T) {
	mux := http.NewServeMux()
	Mount(mux)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	t.Run("GET /swagger redirects to /swagger/index.html", func(t *testing.T) {
		client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}}
		resp, err := client.Get(srv.URL + "/swagger")
		if err != nil {
			t.Fatalf("GET /swagger: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusFound {
			t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusFound)
		}
		loc := resp.Header.Get("Location")
		if loc != "/swagger/index.html" {
			t.Errorf("Location = %q, want %q", loc, "/swagger/index.html")
		}
	})

	t.Run("GET /swagger/ serves the embedded UI", func(t *testing.T) {
		body := get(t, srv.URL+"/swagger/index.html")
		want, err := docs.SwaggerUI.ReadFile("swagger-ui/index.html")
		if err != nil {
			t.Fatalf("reading embedded index.html: %v", err)
		}
		if body != string(want) {
			t.Errorf("served index.html does not match embedded copy")
		}
	})

	t.Run("GET /api/ serves every OpenAPI YAML spec", func(t *testing.T) {
		for _, name := range Specs() {
			body := get(t, srv.URL+"/api/"+name)
			want, err := docs.APIYAML.ReadFile("api/" + name)
			if err != nil {
				t.Fatalf("reading embedded %s: %v", name, err)
			}
			if body != string(want) {
				t.Errorf("served /api/%s does not match embedded copy", name)
			}
		}
	})
}

// emptyEmbed simulates the day the docs embed layout drops the swagger-ui/ or
// api/ tree: an embed.FS declared without a //go:embed directive is empty, and
// rising fs.Sub over it silently yields an empty tree, exactly what Mount must
// fail fast on instead of serving 404s at /swagger/ and /api/.
var emptyEmbed embed.FS

// TestMountSubPanicsOnBrokenFS guards the fail-fast contract: Mount must panic
// when a docs embed sub-tree is absent rather than silently serve an empty
// tree that 404s at /swagger/ and /api/.
func TestMountSubPanicsOnBrokenFS(t *testing.T) {
	for _, name := range []string{"swagger-ui", "api"} {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("mustSub(%q): got no panic on a broken embed, want panic", name)
				}
			}()
			mustSub(emptyEmbed, name)
		}()
	}
}

func TestSpecs(t *testing.T) {
	got := Specs()

	if len(got) == 0 {
		t.Fatal("Specs() returned no filenames")
	}

	for _, name := range got {
		if strings.HasPrefix(name, "api/") || strings.Contains(name, "/") {
			t.Errorf("Specs() entry %q should have no path prefix or dirs", name)
		}
		if !strings.HasSuffix(name, ".yaml") {
			t.Errorf("Specs() entry %q does not end in .yaml", name)
		}
		if _, err := docs.APIYAML.ReadFile("api/" + name); err != nil {
			t.Errorf("Specs() entry %q cannot be read from docs.APIYAML: %v", name, err)
		}
	}

	want, err := embeddedYAMLFilenames()
	if err != nil {
		t.Fatalf("enumerating embedded basenames: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Specs() = %v, want %v", got, want)
	}
}

// TestRealmManagement403 guards the contract that every realm-management
// operation guarded by a_rma documents a 403 Forbidden response, not just
// 401: RequireRealm answers 403 for a verified realm JWT whose a_rma grants
// do not authorize the method+path (internal/auth/middleware.go). Matches the
// upstream-parity the housekeeping and pairing specs already document.
func TestRealmManagement403(t *testing.T) {
	specDocuments403(t, "astarte_realm_management_api.yaml")
}

// TestAppEngine403 is the AppEngine half of the same contract: every route is
// wrapped by RequireRealm(auth.ClaimAppEngine) (internal/appengine/http.go),
// so a verified a_aea JWT whose grants do not authorize the method+path gets
// 403 Forbidden, not 401 (astarteapi.WriteForbidden,
// internal/auth/middleware.go) — the behaviour pinned by
// internal/appengine/http_test.go. Documenting only 401 would teach a
// generated client that 403 cannot happen.
func TestAppEngine403(t *testing.T) {
	specDocuments403(t, "astarte_appengine_api.yaml")
}

// specDocuments403 checks that every operation in the given spec documents a
// 403 response and that components.responses defines the Forbidden response
// they reference.
func specDocuments403(t *testing.T, filename string) {
	t.Helper()
	b, err := docs.APIYAML.ReadFile("api/" + filename)
	if err != nil {
		t.Fatalf("reading %s: %v", filename, err)
	}
	lines := strings.Split(string(b), "\n")

	pathIdx, compIdx := -1, -1
	for i, l := range lines {
		if l == "paths:" {
			pathIdx = i
		}
		if l == "components:" {
			compIdx = i
		}
		if pathIdx >= 0 && compIdx >= 0 {
			break
		}
	}
	if pathIdx < 0 || compIdx < pathIdx {
		t.Fatal("cannot locate the paths and components sections")
	}

	methodRe := regexp.MustCompile(`^    (get|post|put|delete|patch):$`)
	var ops []int
	for i, l := range lines[pathIdx:compIdx] {
		if methodRe.MatchString(l) {
			ops = append(ops, pathIdx+i)
		}
	}
	if len(ops) == 0 {
		t.Fatal("found no operation blocks in the spec")
	}

	for i, start := range ops {
		end := compIdx
		if i+1 < len(ops) {
			end = ops[i+1]
		}
		if !containsLine(lines[start:end], `        "403":`) {
			t.Errorf("operation starting at line %d documents no 403 response", start+1)
		}
	}

	if !containsLine(lines[compIdx:], `    Forbidden:`) {
		t.Error("components.responses defines no Forbidden response")
	}
}

// TestHousekeepingAsyncOperationParamDocumented guards that the housekeeping
// spec tells clients the `?async_operation` parameter exists. Astrate accepts
// and ignores it on realm create and delete (deviation 17,
// docs/COMPATIBILITY.md), so an upstream client keeps working — but a client
// generated from the spec has to learn the parameter from the spec. This is
// the documentation half of the behaviour pinned by
// TestHousekeepingAsyncOperationParam in internal/housekeeping.
func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	const ref = `        - $ref: "#/components/parameters/AsyncOperation"`
	for _, op := range []string{"createRealm", "deleteRealm"} {
		block := operationBlock(t, lines, op)
		if !containsLine(block, ref) {
			t.Errorf("operation %s does not $ref the AsyncOperation parameter", op)
		}
	}

	comp := componentBlock(t, lines, "    AsyncOperation:")
	for _, want := range []string{
		"      name: async_operation",
		"      in: query",
		"      required: false",
		"      schema:",
		"        type: boolean",
		"        default: false",
	} {
		if !containsLine(comp, want) {
			t.Errorf("components parameter AsyncOperation is missing line %q", want)
		}
	}
	if !strings.Contains(strings.Join(comp, "\n"), "Accepted and ignored") {
		t.Error("components parameter AsyncOperation does not say the value is accepted and ignored")
	}
}

// TestRealmManagementAsyncOperationParamDocumented guards that the realm
// management spec tells clients the `?async_operation` parameter exists on the
// five operations upstream 1.4 runs in the background: interface
// install/update/delete, device deletion and trigger-delivery-policy delete.
// Astrate accepts and ignores it there (deviation 17, docs/COMPATIBILITY.md) —
// none of the five handlers reads the query string, which is the whole of the
// acceptance — so an upstream client keeps working, but a client generated from
// the spec has to learn the parameter from the spec. This is the documentation
// half of the behaviour pinned by TestRealmManagementAsyncOperationParam in
// internal/realm.
func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	const ref = `        - $ref: "#/components/parameters/AsyncOperation"`
	for _, op := range []string{
		"installInterface", "updateInterface", "deleteInterface",
		"deleteDevice", "deletePolicy",
	} {
		block := operationBlock(t, lines, op)
		if !containsLine(block, ref) {
			t.Errorf("operation %s does not $ref the AsyncOperation parameter", op)
		}
	}

	comp := componentBlock(t, lines, "    AsyncOperation:")
	for _, want := range []string{
		"      name: async_operation",
		"      in: query",
		"      required: false",
		"      schema:",
		"        type: boolean",
		"        default: false",
	} {
		if !containsLine(comp, want) {
			t.Errorf("components parameter AsyncOperation is missing line %q", want)
		}
	}
	if !strings.Contains(strings.Join(comp, "\n"), "Accepted and ignored") {
		t.Error("components parameter AsyncOperation does not say the value is accepted and ignored")
	}
}

// propertyBlock returns the lines of the schema property with the given name
// inside a components.schemas entry, up to the next property or the end of the
// entry.
func propertyBlock(t *testing.T, lines []string, name string) []string {
	t.Helper()
	key := "        " + name + ":"
	start := -1
	for i, l := range lines {
		if l == key {
			start = i
			break
		}
	}
	if start < 0 {
		t.Fatalf("schema declares no property %q", name)
	}
	for i := start + 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" {
			continue
		}
		if !strings.HasPrefix(lines[i], "          ") {
			return lines[start+1 : i]
		}
	}
	return lines[start+1:]
}

// TestHousekeepingRetentionZeroFoldDocumented guards that the housekeeping spec
// tells clients that `datastream_maximum_storage_retention: 0` means unset, not
// a literal zero-second retention. The wire folds it on both paths: PATCH maps
// `null || val == 0` to ClearRetention (internal/housekeeping/http.go) and
// create folds 0 to nil before injecting the configured default
// (internal/housekeeping/service.go) — upstream parity measured on v1.2.0. A
// client that follows the spec without this note sends 0 and silently gets
// unlimited instead. `device_registration_limit` has no such fold, so the same
// 0 is stored literally there; the asymmetry is documented on both fields so
// the two do not read as interchangeable. This is the documentation half of the
// behaviour pinned by ZeroRetentionUnsets and TestHousekeepingRetentionDefault
// in internal/housekeeping.
func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	patch := strings.Join(operationBlock(t, lines, "patchRealm"), "\n")
	for _, want := range []string{
		"explicit `0` clears it too",
		"0 is folded to unset",
		"has no such fold",
		"stored as the literal limit `0`",
	} {
		if !strings.Contains(patch, want) {
			t.Errorf("patchRealm description does not say %q", want)
		}
	}

	for _, schema := range []string{"    RealmCreate:", "    RealmPatch:"} {
		block := componentBlock(t, lines, schema)

		retention := strings.Join(propertyBlock(t, block, "datastream_maximum_storage_retention"), "\n")
		if !strings.Contains(retention, "0") || !strings.Contains(retention, "unlimited") {
			t.Errorf("%s retention does not say what 0 means", schema)
		}
		if !strings.Contains(retention, "unset") {
			t.Errorf("%s retention does not say 0 is folded to unset", schema)
		}

		limit := strings.Join(propertyBlock(t, block, "device_registration_limit"), "\n")
		if !strings.Contains(limit, "`0`") {
			t.Errorf("%s device_registration_limit does not record what 0 means", schema)
		}
		if strings.Contains(limit, "folded") {
			t.Errorf("%s device_registration_limit claims a 0 fold that the wire does not have", schema)
		}
		if !strings.Contains(limit, "null") {
			t.Errorf("%s device_registration_limit does not say null clears the limit", schema)
		}
	}

	create := strings.Join(propertyBlock(t,
		componentBlock(t, lines, "    RealmCreate:"), "datastream_maximum_storage_retention"), "\n")
	if !strings.Contains(create, "configured default retention") {
		t.Error("RealmCreate retention does not say 0 skips the configured default retention")
	}
}

// TestAppEngineAddGroupDevice422Documented guards that adding a device to a
// group documents its 422, which is a different failure from the 400: a
// well-formed body whose device_id does not parse is rejected as a field error
// (FieldErrors{"device_id": ...}, internal/appengine/service.go) and answered
// 422 with the changeset-shaped FieldErrorsDetail body
// (astarteapi.WriteFieldErrors, internal/appengine/http.go), while the
// documented 400 only covers a body DecodeData cannot parse. So the 422 must
// $ref the ValidationErrors response — pointing it at BadRequest would teach a
// generated client the wrong body shape for the common case. This is the
// documentation half of the behaviour pinned by TestAddGroupDeviceErrors in
// internal/appengine.
func TestAppEngineAddGroupDevice422Documented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_appengine_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_appengine_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	const (
		key = `        "422":`
		ref = `          $ref: "#/components/responses/ValidationErrors"`
	)

	block := operationBlock(t, lines, "addGroupDevice")
	idx := -1
	for i, l := range block {
		if l == key {
			idx = i
			break
		}
	}
	if idx < 0 {
		t.Fatal("addGroupDevice documents no 422 response")
	}
	if idx+1 >= len(block) || block[idx+1] != ref {
		t.Errorf("addGroupDevice 422 is followed by %q, want the ValidationErrors $ref", block[idx+1])
	}

	comp := strings.Join(componentBlock(t, lines, "    ValidationErrors:"), "\n")
	if !strings.Contains(comp, `#/components/schemas/FieldErrorsDetail`) {
		t.Error("components.responses.ValidationErrors does not carry the FieldErrorsDetail body")
	}
}

// TestAppEngineDataDelete400Documented guards that the three data DELETE
// operations document their 400: unsetting a path that matches no endpoint
// mapping makes the trie lookup fail with engine.ErrPathNotFound
// (UnsetServerProperty, internal/engine/serverdata.go), which writeError
// answers as 400 "Endpoint not found" (internal/appengine/http.go) — the same
// 400 the PUT/POST twins already document, because it is the same sentinel on
// the same write path. Documenting only 404 would teach a generated client
// that an unmatched path is a missing resource; the read path does use 404 for
// that (appengine.ErrPathNotFound), and the two are deliberately distinct. This
// is the documentation half of the behaviour pinned by
// TestWriteErrorTaxonomy in internal/appengine.
func TestAppEngineDataDelete400Documented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_appengine_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_appengine_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	const (
		key = `        "400":`
		ref = `          $ref: "#/components/responses/BadRequest"`
	)

	for _, op := range []string{"deleteDataByAlias", "deleteData", "deleteDataInGroup"} {
		block := operationBlock(t, lines, op)
		idx := -1
		for i, l := range block {
			if l == key {
				idx = i
				break
			}
		}
		if idx < 0 {
			t.Errorf("operation %s documents no 400 response", op)
			continue
		}
		if idx+1 >= len(block) || block[idx+1] != ref {
			t.Errorf("operation %s 400 is followed by %q, want the BadRequest $ref", op, block[idx+1])
		}
	}

	if !containsLine(lines, "    BadRequest:") {
		t.Error("components.responses defines no BadRequest response")
	}
}

// TestPairingErrorDetailExamples pins the error-detail examples in the pairing
// spec's components.responses to the frozen canonical strings the wire emits
// (astarteapi/envelope.go), so a generated client that copies an example does
// not learn a detail that never goes over the wire. Phoenix renders "Bad
// request" and "Internal server error" — the capital-R/capital-S spellings the
// spec used to carry were reconstructed, not measured on upstream 1.2.0, and a
// startectl/SDK error path matching on the canonical string would miss them.
// The 404 detail is deliberately absent here: pairing answers an unmatched
// path with DetailPageNotFound while DetailNotFound is the generic form, and
// the spec's own DeviceNotFound example is the one that carries a constant.
func TestPairingErrorDetailExamples(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	const detailPrefix = "              detail: "

	for _, tc := range []struct {
		response string
		want     string
	}{
		{"    BadRequest:", astarteapi.DetailBadRequest},
		{"    Unauthorized:", astarteapi.DetailUnauthorized},
		{"    Forbidden:", astarteapi.DetailForbidden},
		{"    DeviceNotFound:", astarteapi.DetailDeviceNotFound},
		{"    InternalServerError:", astarteapi.DetailInternalServerError},
	} {
		block := componentBlock(t, lines, tc.response)

		got := ""
		for _, l := range block {
			if strings.HasPrefix(l, detailPrefix) {
				got = strings.TrimPrefix(l, detailPrefix)
				break
			}
		}
		if got == "" {
			t.Errorf("response %q carries no %q example", tc.response, "detail")
			continue
		}
		if got != tc.want {
			t.Errorf("response %q example detail = %q, want the canonical %q",
				strings.TrimSpace(tc.response), got, tc.want)
		}
	}
}

// TestRealmManagementErrorDetailExamples pins the error-detail examples in the
// realm management spec's components.responses to the frozen canonical strings
// the wire emits (astarteapi/envelope.go), the same guard the pairing spec gets
// from TestPairingErrorDetailExamples. The BadRequest and InternalServerError
// examples carried the capital-R/capital-S spellings, which upstream Phoenix
// never renders and no handler emits: every 400 on this surface goes through
// astarteapi.WriteBadRequest and every 500 through WriteInternalServerError
// (internal/realm/http.go), and the SDK/astartectl error paths match on the
// canonical strings. The Conflict example is an endpoint-specific detail with no
// constant behind it, so it stays free-form; the three 422s that used to share
// one ValidationError component are pinned by
// TestRealmManagement422ValidationDetails.
func TestRealmManagementErrorDetailExamples(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	const detailPrefix = "              detail: "

	for _, tc := range []struct {
		response string
		want     string
	}{
		{"    BadRequest:", astarteapi.DetailBadRequest},
		{"    Unauthorized:", astarteapi.DetailUnauthorized},
		{"    Forbidden:", astarteapi.DetailForbidden},
		{"    NotFound:", astarteapi.DetailNotFound},
		{"    DeviceNotFound:", astarteapi.DetailDeviceNotFound},
		{"    InternalServerError:", astarteapi.DetailInternalServerError},
	} {
		block := componentBlock(t, lines, tc.response)

		got := ""
		for _, l := range block {
			if strings.HasPrefix(l, detailPrefix) {
				got = strings.TrimPrefix(l, detailPrefix)
				break
			}
		}
		if got == "" {
			t.Errorf("response %q carries no %q example", tc.response, "detail")
			continue
		}
		if got != tc.want {
			t.Errorf("response %q example detail = %q, want the canonical %q",
				strings.TrimSpace(tc.response), got, tc.want)
		}
	}
}

// TestAppEngineErrorDetailExamples pins the error-detail examples in the app
// engine spec's components.responses to the frozen canonical strings the wire
// emits (astarteapi/envelope.go), the same guard the pairing and realm
// management specs get from TestPairingErrorDetailExamples and
// TestRealmManagementErrorDetailExamples. The BadRequest and InternalServerError
// examples carried the capital-R/capital-S spellings, which upstream Phoenix
// never renders and no handler emits: every 400 on this surface goes through
// astarteapi.WriteBadRequest and every 500 through WriteInternalServerError
// (internal/appengine/http.go), and the SDK/astartectl error paths match on the
// canonical strings. The responses whose example is an endpoint-specific detail
// with no constant behind it (MethodNotAllowed, ValueTooLarge,
// UnprocessableEntity, the three Conflicts) stay free-form and are not pinned.
func TestAppEngineErrorDetailExamples(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_appengine_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_appengine_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	const detailPrefix = "              detail: "

	for _, tc := range []struct {
		response string
		want     string
	}{
		{"    BadRequest:", astarteapi.DetailBadRequest},
		{"    Unauthorized:", astarteapi.DetailUnauthorized},
		{"    Forbidden:", astarteapi.DetailForbidden},
		{"    NotFound:", astarteapi.DetailNotFound},
		{"    DeviceNotFound:", astarteapi.DetailDeviceNotFound},
		{"    InternalServerError:", astarteapi.DetailInternalServerError},
	} {
		block := componentBlock(t, lines, tc.response)

		got := ""
		for _, l := range block {
			if strings.HasPrefix(l, detailPrefix) {
				got = strings.TrimPrefix(l, detailPrefix)
				break
			}
		}
		if got == "" {
			t.Errorf("response %q carries no %q example", tc.response, "detail")
			continue
		}
		if got != tc.want {
			t.Errorf("response %q example detail = %q, want the canonical %q",
				strings.TrimSpace(tc.response), got, tc.want)
		}
	}
}

// TestNativeErrorDetailExamples pins the error-detail examples in the native
// spec's components.responses to the frozen canonical strings the wire emits
// (astarteapi/envelope.go), the same guard the pairing, realm management and
// app engine specs get from TestPairingErrorDetailExamples,
// TestRealmManagementErrorDetailExamples and TestAppEngineErrorDetailExamples.
// Only the 500 was wrong: it showed "Internal Server Error", the capital-S
// spelling reconstructed rather than measured, while every 500 on the two
// websocket operations the component is $ref'd from answers through
// astarteapi.WriteInternalServerError (auth.Middleware's upgrade failure,
// internal/auth/middleware.go, and the channels socket,
// internal/appengine/channels/ws.go) with the canonical "Internal server error".
// The 401 and 403 are already the canonical spellings. This spec has no
// BadRequest component at all: the native surface is the compat health and
// version endpoints plus the two sockets, and none of them answers 400 through
// WriteBadRequest, so there is nothing here to pin.
func TestNativeErrorDetailExamples(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astrate_native_api.yaml")
	if err != nil {
		t.Fatalf("reading astrate_native_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	const detailPrefix = "              detail: "

	for _, tc := range []struct {
		response string
		want     string
	}{
		{"    Unauthorized:", astarteapi.DetailUnauthorized},
		{"    Forbidden:", astarteapi.DetailForbidden},
		{"    InternalServerError:", astarteapi.DetailInternalServerError},
	} {
		block := componentBlock(t, lines, tc.response)

		got := ""
		for _, l := range block {
			if strings.HasPrefix(l, detailPrefix) {
				got = strings.TrimPrefix(l, detailPrefix)
				break
			}
		}
		if got == "" {
			t.Errorf("response %q carries no %q example", tc.response, "detail")
			continue
		}
		if got != tc.want {
			t.Errorf("response %q example detail = %q, want the canonical %q",
				strings.TrimSpace(tc.response), got, tc.want)
		}
	}
}

// TestHousekeepingErrorDetailExamples pins the error-detail examples in the
// housekeeping spec's components.responses to the frozen canonical strings the
// wire emits (astarteapi/envelope.go), the same guard the pairing, realm
// management, app engine and native specs get from
// TestPairingErrorDetailExamples, TestRealmManagementErrorDetailExamples,
// TestAppEngineErrorDetailExamples and TestNativeErrorDetailExamples. This is
// the last spec still carrying a reconstructed capitalisation: its BadRequest
// and InternalServerError examples read "Bad Request" and "Internal Server
// Error", spellings no upstream Phoenix render answers with and no handler
// emits — every 400 on this surface goes through astarteapi.WriteBadRequest
// and the 500 through WriteInternalServerError (internal/housekeeping/http.go:
// the writeError switch), and the SDK/astartectl error paths match on the
// canonical strings. The 401 and 403 were already canonical. The responses
// whose example is an endpoint-specific detail with no constant behind it
// (MethodNotAllowed, Conflict, ValidationError) stay free-form and are not
// pinned.
//
// There is no DeviceNotFound pin here, unlike the app engine and realm
// management specs: this surface never answers a device 404, so
// components.responses declares no such component to pin. The schema-level
// ErrorDetail example is deliberately not pinned either — it is not a
// components.responses entry and its shape question belongs to its own line.
func TestHousekeepingErrorDetailExamples(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
	if err != nil {
		t.Fatalf("reading astrate_housekeeping_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	const detailPrefix = "              detail: "

	for _, tc := range []struct {
		response string
		want     string
	}{
		{"    BadRequest:", astarteapi.DetailBadRequest},
		{"    Unauthorized:", astarteapi.DetailUnauthorized},
		{"    Forbidden:", astarteapi.DetailForbidden},
		{"    NotFound:", astarteapi.DetailNotFound},
		{"    InternalServerError:", astarteapi.DetailInternalServerError},
	} {
		block := componentBlock(t, lines, tc.response)

		got := ""
		for _, l := range block {
			if strings.HasPrefix(l, detailPrefix) {
				got = strings.TrimPrefix(l, detailPrefix)
				break
			}
		}
		if got == "" {
			t.Errorf("response %q carries no %q example", tc.response, "detail")
			continue
		}
		if got != tc.want {
			t.Errorf("response %q example detail = %q, want the canonical %q",
				strings.TrimSpace(tc.response), got, tc.want)
		}
	}
}

// TestRealmManagement422ValidationDetails pins the 422 examples the three
// non-interface realm-management operations carry to the details the wire
// actually emits, so a generated client that copies one does not learn a
// message no handler can produce.
//
// The three operations used to share components.responses.ValidationError, whose
// single example read `invalid interface: missing "interface_name"`. That string
// is wrong for all three: it names an interface, and none of them installs or
// updates one. Each now carries its own inline 422, the pattern
// installInterface and updateInterface already use for their own details.
//
// The createPolicy examples are derived from the compiler here instead of being
// written out, so a change to the message fails here instead of drifting the
// spec — the coupling TestPairingDeviceIDEncodingDocumented gets from
// deviceid. The putAuthConfig and deletePolicy details have no exported
// constant behind them (Service.SetAuthKey and Service.DeletePolicy,
// internal/realm/service.go), so they are spelled out and cited here.
func TestRealmManagement422ValidationDetails(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	if strings.Contains(string(b), "#/components/responses/ValidationError") {
		t.Error("spec still $refs a shared ValidationError response; one example cannot be right for every surface that refs it")
	}
	if containsLine(lines, "    ValidationError:") {
		t.Error("spec still defines components.responses.ValidationError")
	}

	_, noHandlers := triggers.CompilePolicy([]byte(`{"name":"p","error_handlers":[]}`))
	if noHandlers == nil {
		t.Fatal("triggers.CompilePolicy accepted a policy with no error handler; this test's premise no longer holds")
	}
	_, badName := triggers.CompilePolicy(
		[]byte(`{"name":"","error_handlers":[{"on":"any_error","strategy":"discard"}]}`))
	if badName == nil {
		t.Fatal("triggers.CompilePolicy accepted an empty policy name; this test's premise no longer holds")
	}

	for _, tc := range []struct {
		op   string
		want []string
	}{
		{"putAuthConfig", []string{`jwt_public_key_pem can't be blank`}},
		{"createPolicy", []string{noHandlers.Error(), badName.Error()}},
		{"deletePolicy", []string{`policy "audit" is still used by trigger "on_audit_failure"`}},
	} {
		block := strings.Join(responseBlock(t, operationBlock(t, lines, tc.op), "422"), "\n")
		for _, want := range tc.want {
			if !strings.Contains(block, want) {
				t.Errorf("%s 422 carries no example with detail %q", tc.op, want)
			}
		}
		if strings.Contains(block, "#/components/responses/") {
			t.Errorf("%s 422 still $refs a response component instead of carrying its own example", tc.op)
		}
	}
}

// responseBlock returns the lines of the operation's response with the given
// status code, up to the next status code or the end of the operation block.
func responseBlock(t *testing.T, block []string, status string) []string {
	t.Helper()
	key := `        "` + status + `":`
	start := -1
	for i, l := range block {
		if l == key {
			start = i
			break
		}
	}
	if start < 0 {
		t.Fatalf("operation block declares no %s response", status)
	}
	for i := start + 1; i < len(block); i++ {
		if strings.HasPrefix(block[i], `        "`) {
			return block[start+1 : i]
		}
	}
	return block[start+1:]
}

// TestPairingDeviceIDEncodingDocumented guards that the pairing spec states the
// exact wire form of a device identifier — the 22-character unpadded base64url
// string deviceid.Parse accepts — instead of the "base64-encoded 128-bit"
// wording, which read as standard base64 and let a generated client send a
// well-formed `+`/`/` identifier or a padded 23-character spelling and collect
// a 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}`
// (internal/pairing/http.go). The shipped pattern is compiled here and held to
// the same spellings the parser is, so the two cannot drift apart: the bounds
// and the description wording are derived from deviceid.EncodedLen rather than
// written out, so a change to the accepted length fails this test too.
func TestPairingDeviceIDEncodingDocumented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_pairing_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_pairing_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	wantLen := strconv.Itoa(deviceid.EncodedLen)
	wantDesc := fmt.Sprintf("%s-character unpadded", wantLen)

	param := strings.Join(componentBlock(t, lines, "    DeviceID:"), "\n")
	for _, want := range []string{wantDesc, "base64url"} {
		if !strings.Contains(param, want) {
			t.Errorf("DeviceID parameter description does not say %q", want)
		}
	}
	if strings.Contains(param, "(base64-encoded 128-bit)") {
		t.Error("DeviceID parameter still describes the ID as plain base64")
	}

	block := propertyBlock(t, lines, "hw_id")
	hwID := strings.Join(block, "\n")
	for _, want := range []string{wantDesc, "base64url"} {
		if !strings.Contains(hwID, want) {
			t.Errorf("RegisterRequest.hw_id description does not say %q", want)
		}
	}
	for _, want := range []string{
		"          minLength: " + wantLen,
		"          maxLength: " + wantLen,
	} {
		if !containsLine(block, want) {
			t.Errorf("RegisterRequest.hw_id is missing line %q", want)
		}
	}

	pattern := ""
	for _, l := range block {
		if v, ok := strings.CutPrefix(l, "          pattern: "); ok {
			pattern = strings.Trim(v, `'"`)
			break
		}
	}
	if pattern == "" {
		t.Fatal("RegisterRequest.hw_id carries no pattern")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		t.Fatalf("hw_id pattern %q does not compile: %v", pattern, err)
	}

	if !re.MatchString("dT6hS2W9TT6LEnP25ks_lg") {
		t.Errorf("hw_id pattern %q rejects the spec's own example dT6hS2W9TT6LEnP25ks_lg", pattern)
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
			t.Errorf("hw_id pattern %q accepts %q, which deviceid.Parse rejects", pattern, spelling)
		}
	}
}

// TestRealmManagementDeviceIDEncodingDocumented is the realm-management twin of
// TestPairingDeviceIDEncodingDocumented: the `DeviceID` path parameter of
// deleteDevice must state the same wire form, and constrain itself to it. The
// parameter used to be a bare `type: string` described only as "The device
// hardware ID", which read as any string and let a client send a well-formed
// standard-base64 or padded 23-character ID. Service.DeleteDevice folds every
// deviceid.Parse failure into store.ErrNotFound (internal/realm/service.go), so
// such a client is answered 404 Device not found rather than told its ID is
// malformed — the description has to name that consequence, and the pattern has
// to reject the spellings the parser rejects. The bounds and the description
// wording are derived from deviceid.EncodedLen, and the shipped pattern is
// compiled and held against deviceid.Parse here, so the spec cannot drift from
// the parser. The strict trailing-bits rule is deliberately not encoded in the
// pattern: whether recording it is worth it is still open.
func TestRealmManagementDeviceIDEncodingDocumented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	wantLen := strconv.Itoa(deviceid.EncodedLen)
	wantDesc := fmt.Sprintf("%s-character unpadded", wantLen)

	param := componentBlock(t, lines, "    DeviceID:")
	for _, want := range []string{wantDesc, "base64url", astarteapi.DetailDeviceNotFound} {
		if !strings.Contains(strings.Join(param, "\n"), want) {
			t.Errorf("DeviceID parameter does not say %q", want)
		}
	}
	for _, want := range []string{
		"        minLength: " + wantLen,
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
