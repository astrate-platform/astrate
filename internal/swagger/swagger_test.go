package swagger

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	docs "github.com/astrate-platform/astrate/docs"
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
