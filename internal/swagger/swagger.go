// Package swagger serves the embedded Swagger UI and OpenAPI YAML specs at
// /swagger/ and /api/ respectively.
package swagger

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	docs "github.com/astrate-platform/astrate/docs"
)

// Mount registers the /swagger and /api routes on the given mux.
// /swagger redirects to /swagger/index.html; /swagger/ serves the static UI;
// /api/ serves the OpenAPI YAML specs.
func Mount(mux *http.ServeMux) {
	uiRoot := mustSub(docs.SwaggerUI, "swagger-ui")
	apiRoot := mustSub(docs.APIYAML, "api")

	MountWithFS(mux, uiRoot, apiRoot)
}

// MountWithFS registers the /swagger and /api routes on the given mux from
// already-subtreed file systems. Mount computes the sub-trees from the
// embedded docs and delegates here; injecting fs.FS values makes the served
// roots testable and lets a broken embed fail fast instead of silently
// serving an empty tree.
func MountWithFS(mux *http.ServeMux, uiRoot, apiRoot fs.FS) {
	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
	})
	mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(uiRoot))))

	// Serve YAML files at /api/ so the relative ../api/*.yaml in index.html
	// resolves correctly when the page is loaded from /swagger/index.html.
	mux.Handle("GET /api/", http.StripPrefix("/api/", http.FileServer(http.FS(apiRoot))))
}

// mustSub returns fs.Sub(fsys, name), panicking if the sub-tree is absent so
// a broken docs embed fails fast instead of silently serving an empty tree.
// fs.Sub itself only errors on an invalid path (embed.FS has no Sub method),
// so the sub-root is also opened to catch an embed whose layout dropped the
// tree entirely.
func mustSub(fsys fs.FS, name string) fs.FS {
	sub, err := fs.Sub(fsys, name)
	if err != nil {
		panic(fmt.Sprintf("swagger: fs.Sub(%q): %v", name, err))
	}
	if _, err := sub.Open("."); err != nil {
		panic(fmt.Sprintf("swagger: %q missing from embedded fs: %v", name, err))
	}
	return sub
}

// Specs returns the list of available YAML spec filenames (without path prefix).
func Specs() []string {
	var names []string
	_ = fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".yaml") {
			names = append(names, strings.TrimPrefix(path, "api/"))
		}
		return nil
	})
	return names
}
