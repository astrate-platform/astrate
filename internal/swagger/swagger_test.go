package swagger

import (
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"reflect"
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
