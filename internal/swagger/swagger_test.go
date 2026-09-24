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
	b, err := docs.APIYAML.ReadFile("api/astarte_realm_management_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_realm_management_api.yaml: %v", err)
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

	methodRe := regexp.MustCompile(`^    (get|post|put|delete):$`)
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
