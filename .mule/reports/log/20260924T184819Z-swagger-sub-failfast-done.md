slug: swagger-sub-failfast
verdict: done
at:  1d321ea
ran: 2026-09-24T18:48:19Z on DietPi in 371s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
✱ Glob "internal/swagger/*" 2 matches
→ Read internal/swagger/swagger.go
✱ Glob "docs/*.go" 1 match
→ Read internal/swagger/swagger_test.go
→ Read docs/embed.go
✱ Grep "swagger.Mount" 1 match
Task: make `Mount` panic when the `docs.SwaggerUI`/`docs.APIYAML` `fs.Sub` fails, via an `fs.FS`-injecting `MountWithFS` variant plus a broken-FS test. Let me check the caller and environment, then implement.
→ Read cmd/astrate/main.go [offset=408, limit=10]

$ go version && git status --short
go version go1.26.7 linux/arm64

Now the implementation:
← Edit internal/swagger/swagger.go
Index: /root/astrate-mule/internal/swagger/swagger.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger.go
+++ /root/astrate-mule/internal/swagger/swagger.go
@@ -2,8 +2,9 @@
 // /swagger/ and /api/ respectively.
 package swagger
 
 import (
+	"fmt"
 	"io/fs"
 	"net/http"
 	"strings"
 

← Edit internal/swagger/swagger.go
Index: /root/astrate-mule/internal/swagger/swagger.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger.go
+++ /root/astrate-mule/internal/swagger/swagger.go
@@ -14,11 +14,20 @@
 // Mount registers the /swagger and /api routes on the given mux.
 // /swagger redirects to /swagger/index.html; /swagger/ serves the static UI;
 // /api/ serves the OpenAPI YAML specs.
 func Mount(mux *http.ServeMux) {
-	uiRoot, _ := fs.Sub(docs.SwaggerUI, "swagger-ui")
-	apiRoot, _ := fs.Sub(docs.APIYAML, "api")
+	uiRoot := mustSub(docs.SwaggerUI, "swagger-ui")
+	apiRoot := mustSub(docs.APIYAML, "api")
 
+	MountWithFS(mux, uiRoot, apiRoot)
+}
+
+// MountWithFS registers the /swagger and /api routes on the given mux from
+// already-subtreed file systems. Mount computes the sub-trees from the
+// embedded docs and delegates here; injecting fs.FS values makes the served
+// roots testable and lets a broken embed fail fast instead of silently
+// serving an empty tree.
+func MountWithFS(mux *http.ServeMux, uiRoot, apiRoot fs.FS) {
 	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
 		http.Redirect(w, r, "/swagger/index.html", http.StatusFound)
 	})
 	mux.Handle("GET /swagger/", http.StripPrefix("/swagger/", http.FileServer(http.FS(uiRoot))))
@@ -27,8 +36,18 @@
 	// resolves correctly when the page is loaded from /swagger/index.html.
 	mux.Handle("GET /api/", http.StripPrefix("/api/", http.FileServer(http.FS(apiRoot))))
 }
 
+// mustSub returns fs.Sub(fsys, name), panicking if the sub-tree is absent so
+// a broken docs embed fails fast instead of silently serving an empty tree.
+func mustSub(fsys fs.FS, name string) fs.FS {
+	sub, err := fs.Sub(fsys, name)
+	if err != nil {
+		panic(fmt.Sprintf("swagger: %q missing from embedded fs: %v", name, err))
+	}
+	return sub
+}
+
 // Specs returns the list of available YAML spec filenames (without path prefix).
 func Specs() []string {
 	var names []string
 	_ = fs.WalkDir(docs.APIYAML, "api", func(path string, d fs.DirEntry, err error) error {

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -9,8 +9,9 @@
 	"regexp"
 	"sort"
 	"strings"
 	"testing"
+	"testing/fstest"
 
 	docs "github.com/astrate-platform/astrate/docs"
 )
 

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -65,9 +65,8 @@
 	})
 }
 
 func TestSpecs(t *testing.T) {
-	got := Specs()
 
 	if len(got) == 0 {
 		t.Fatal("Specs() returned no filenames")
 	}

That edit accidentally removed a line. Let me fix it and insert the new test properly.
→ Read internal/swagger/swagger_test.go [offset=63, limit=18]

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -64,9 +64,24 @@
 		}
 	})
 }
 
+func TestMountSubPanicsOnBrokenFS(t *testing.T) {
+	broken := fstest.MapFS{}
+	for _, name := range []string{"swagger-ui", "api"} {
+		func() {
+			defer func() {
+				if r := recover(); r == nil {
+					t.Errorf("mustSub(%q): fs.Sub succeeded on a broken fs, want panic", name)
+				}
+			}()
+			mustSub(broken, name)
+		}()
+	}
+}
+
 func TestSpecs(t *testing.T) {
+	got := Specs()
 
 	if len(got) == 0 {
 		t.Fatal("Specs() returned no filenames")
 	}

Now run the gate:
$ gofmt -l . && go vet ./... && go test ./...
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
ok  	github.com/astrate-platform/astrate/examples/flow-container-echo	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
--- FAIL: TestMountSubPanicsOnBrokenFS (0.00s)
    swagger_test.go:74: mustSub("swagger-ui"): fs.Sub succeeded on a broken fs, want panic
    swagger_test.go:74: mustSub("api"): fs.Sub succeeded on a broken fs, want panic
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.091s
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)
FAIL

The empty `MapFS` didn't error on `fs.Sub` — let me probe why.
$ cat > /tmp/probe.go <<'EOF'
package main

import (
	"fmt"
	"io/fs"
	"testing/fstest"
)

func main() {
	broken := fstest.MapFS{}
	for _, name := range []string{"swagger-ui", "api"} {
		sub, err := fs.Sub(broken, name)
		fmt.Printf("%q: sub=%v err=%v\n", name, sub != nil, err)
		if err == nil {
			_, err2 := sub.Open("index.html")
			fmt.Printf("  open index.html: %v\n", err2)
		}
	}
}
EOF
go run /tmp/probe.go
"swagger-ui": sub=true err=<nil>
  open index.html: open index.html: file does not exist
"api": sub=true err=<nil>
  open index.html: open index.html: file does not exist

$ go doc testing/fstest.MapFS | head -40 && echo "---" && go doc testing/fstest.MapFile
package fstest // import "testing/fstest"

type MapFS map[string]*MapFile
    A MapFS is a simple in-memory file system for use in tests, represented as
    a map from path names (arguments to Open) to information about the files,
    directories, or symbolic links they represent.

    The map need not include parent directories for files contained in the map;
    those will be synthesized if needed. But a directory can still be included
    by setting the MapFile.Mode's fs.ModeDir bit; this may be necessary for
    detailed control over the directory's fs.FileInfo or to create an empty
    directory.

    File system operations read directly from the map, so that the file system
    can be changed by editing the map as needed. An implication is that file
    system operations must not run concurrently with changes to the map, which
    would be a race. Another implication is that opening or reading a directory
    requires iterating over the entire map, so a MapFS should typically be used
    with not more than a few hundred entries or directory reads.

func (fsys MapFS) Glob(pattern string) ([]string, error)
func (fsys MapFS) Lstat(name string) (fs.FileInfo, error)
func (fsys MapFS) Open(name string) (fs.File, error)
func (fsys MapFS) ReadDir(name string) ([]fs.DirEntry, error)
func (fsys MapFS) ReadFile(name string) ([]byte, error)
func (fsys MapFS) ReadLink(name string) (string, error)
func (fsys MapFS) Stat(name string) (fs.FileInfo, error)
func (fsys MapFS) Sub(dir string) (fs.FS, error)
---
package fstest // import "testing/fstest"

type MapFile struct {
	Data    []byte      // file content or symlink destination
	Mode    fs.FileMode // fs.FileInfo.Mode
	ModTime time.Time   // fs.FileInfo.ModTime
	Sys     any         // fs.FileInfo.Sys
}
    A MapFile describes a single file in a MapFS.

`fstest.MapFS` synthesizes sub-FSs instead of erroring. The realistic broker is an `embed.FS`-style `SubFS`. Let me confirm how the real embed behaves, and that a custom broken `SubFS` errors properly.
$ cat > /tmp/probe2.go <<'EOF'
package main

import (
	"fmt"
	"io/fs"

	docs "github.com/astrate-platform/astrate/docs"
)

func main() {
	_, err := fs.Sub(docs.SwaggerUI, "nonexistent")
	fmt.Printf("real embed missing sub: err=%v\n", err)

	_, err = fs.Sub(docs.SwaggerUI, "swagger-ui")
	fmt.Printf("real embed present sub: err=%v\n", err)
}
EOF
go run /tmp/probe2.go
real embed missing sub: err=<nil>
real embed present sub: err=<nil>

Surprising — `fs.Sub` over embed.FS doesn't validate the sub-path at all. Let me check the actual `Sub` implementation in the installed Go to be sure.
$ grep -rn "func (f \*FS) Sub\|func (f FS) Sub" "$(go env GOROOT)/src/embed/embed.go"; echo "---"; sed -n '/func (f \*FS) Sub/,/^}/p' "$(go env GOROOT)/src/embed/embed.go"
---

$ cd "$(go env GOROOT)/src/embed" && rg -n "Sub\b" embed.go
(no output)

$ go doc embed.FS && echo "---" && sed -n '/^func Sub(/,/^}/p' "$(go env GOROOT)/src/io/fs/fs.go"
package embed // import "embed"

type FS struct {
	// Has unexported fields.
}
    An FS is a read-only collection of files, usually initialized with a
    //go:embed directive. When declared without a //go:embed directive, an FS is
    an empty file system.

    An FS is a read-only value, so it is safe to use from multiple goroutines
    simultaneously and also safe to assign values of type FS to each other.

    FS implements fs.FS, so it can be used with any package that understands
    file system interfaces, including net/http, text/template, and
    html/template.

    See the package documentation for more details about initializing an FS.

func (f FS) Open(name string) (fs.File, error)
func (f FS) ReadDir(name string) ([]fs.DirEntry, error)
func (f FS) ReadFile(name string) ([]byte, error)
---

`embed.FS` in Go 1.26 implements no `Sub` method, so `fs.Sub` returns a prefix wrapper and *never* errors for a valid path — the real failure mode (lost embed sub-tree) must be caught by validating the sub-root actually opens. Let me probe end-to-end before finalizing.
$ cat > /tmp/probe3.go <<'EOF'
package main

import (
	"embed"
	"fmt"
	"io/fs"

	docs "github.com/astrate-platform/astrate/docs"
)

var empty embed.FS

func main() {
	for _, e := range []struct{ fsys fs.FS; name string }{
		{docs.SwaggerUI, "swagger-ui"},
		{docs.APIYAML, "api"},
		{empty, "swagger-ui"},
		{empty, "api"},
	} {
		sub, err := fs.Sub(e.fsys, e.name)
		if err != nil {
			fmt.Printf("%T %q: Sub err=%v\n", e.fsys, e.name, err)
			continue
		}
		_, err = sub.Open(".")
		fmt.Printf("%T %q: Sub ok, Open(.) err=%v\n", e.fsys, e.name, err)
	}
}
EOF
go run /tmp/probe3.go
embed.FS "swagger-ui": Sub ok, Open(.) err=<nil>
embed.FS "api": Sub ok, Open(.) err=<nil>
embed.FS "swagger-ui": Sub ok, Open(.) err=open .: file does not exist
embed.FS "api": Sub ok, Open(.) err=open .: file does not exist

Confirmed. The faithful fail-fast must validate the sub-root opens — `fs.Sub` alone never errors on embed.FS. Updating `mustSub` and the test:
← Edit internal/swagger/swagger.go
Index: /root/astrate-mule/internal/swagger/swagger.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger.go
+++ /root/astrate-mule/internal/swagger/swagger.go
@@ -38,11 +38,17 @@
 }
 
 // mustSub returns fs.Sub(fsys, name), panicking if the sub-tree is absent so
 // a broken docs embed fails fast instead of silently serving an empty tree.
+// fs.Sub itself only errors on an invalid path (embed.FS has no Sub method),
+// so the sub-root is also opened to catch an embed whose layout dropped the
+// tree entirely.
 func mustSub(fsys fs.FS, name string) fs.FS {
 	sub, err := fs.Sub(fsys, name)
 	if err != nil {
+		panic(fmt.Sprintf("swagger: fs.Sub(%q): %v", name, err))
+	}
+	if _, err := sub.Open("."); err != nil {
 		panic(fmt.Sprintf("swagger: %q missing from embedded fs: %v", name, err))
 	}
 	return sub
 }

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -1,7 +1,8 @@
 package swagger
 
 import (
+	"embed"
 	"io"
 	"io/fs"
 	"net/http"
 	"net/http/httptest"
@@ -9,9 +10,8 @@
 	"regexp"
 	"sort"
 	"strings"
 	"testing"
-	"testing/fstest"
 
 	docs "github.com/astrate-platform/astrate/docs"
 )
 

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -64,18 +64,26 @@
 		}
 	})
 }
 
+// emptyEmbed simulates the day the docs embed layout drops the swagger-ui/ or
+// api/ tree: an embed.FS declared without a //go:embed directive is empty, and
+// rising fs.Sub over it silently yields an empty tree, exactly what Mount must
+// fail fast on instead of serving 404s at /swagger/ and /api/.
+var emptyEmbed embed.FS
+
+// TestMountSubPanicsOnBrokenFS guards the fail-fast contract: Mount must panic
+// when a docs embed sub-tree is absent rather than silently serve an empty
+// tree that 404s at /swagger/ and /api/.
 func TestMountSubPanicsOnBrokenFS(t *testing.T) {
-	broken := fstest.MapFS{}
 	for _, name := range []string{"swagger-ui", "api"} {
 		func() {
 			defer func() {
 				if r := recover(); r == nil {
-					t.Errorf("mustSub(%q): fs.Sub succeeded on a broken fs, want panic", name)
+					t.Errorf("mustSub(%q): got no panic on a broken embed, want panic", name)
 				}
 			}()
-			mustSub(broken, name)
+			mustSub(emptyEmbed, name)
 		}()
 	}
 }
 

$ gofmt -l . && go vet ./... && go test ./...
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
ok  	github.com/astrate-platform/astrate/examples/flow-container-echo	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	0.065s
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

Gate is green. Let me confirm the swagger tests verbose and check race capability for an accurate report.
$ go test -v -run 'TestMount|TestSpecs' ./internal/swagger 2>&1 | tail -25; echo "=== race probe ==="; go test -race -run TestMount ./internal/swagger 2>&1 | tail -5
=== RUN   TestMount
=== RUN   TestMount/GET_/swagger_redirects_to_/swagger/index.html
=== RUN   TestMount/GET_/swagger/_serves_the_embedded_UI
=== RUN   TestMount/GET_/api/_serves_every_OpenAPI_YAML_spec
--- PASS: TestMount (0.05s)
    --- PASS: TestMount/GET_/swagger_redirects_to_/swagger/index.html (0.01s)
    --- PASS: TestMount/GET_/swagger/_serves_the_embedded_UI (0.02s)
    --- PASS: TestMount/GET_/api/_serves_every_OpenAPI_YAML_spec (0.02s)
=== RUN   TestMountSubPanicsOnBrokenFS
--- PASS: TestMountSubPanicsOnBrokenFS (0.00s)
=== RUN   TestSpecs
--- PASS: TestSpecs (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.091s
=== race probe ===
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.048s
FAIL

$ git diff --stat && git status --short
internal/swagger/swagger.go      | 30 ++++++++++++++++++++++++++++--
 internal/swagger/swagger_test.go | 23 +++++++++++++++++++++++
 2 files changed, 51 insertions(+), 2 deletions(-)
 M internal/swagger/swagger.go
 M internal/swagger/swagger_test.go

Done. Fail-fast swagger mount via `MountWithFS`, proven by a broken-embed test.

An important correction surfaced while implementing: in Go 1.26, `embed.FS` has no `Sub` method, so `fs.Sub` over the docs embeds returns a prefix wrapper and *never* errors for a missing-but-valid path — a pure "panic on `fs.Sub` error" would never fire on the real failure mode (docs layout dropping the tree → silent 404s). So `mustSub` validates the sub-root by opening `"."` afterward; the test simulates the layout change with an empty `embed.FS` (declared without `//go:embed`, exactly the reported scenario) and requires a panic. Probing confirmed real embeds still open cleanly, and the public `Mount(mux)` + main.go:414 caller are untouched.

```
Done: Mount now panics on a missing docs sub-tree (fs-injecting MountWithFS; mustSub validates the sub-root) instead of silently serving 404-empty trees
Files: internal/swagger/swagger.go, internal/swagger/swagger_test.go
Verified: gofmt -l . -> clean; go vet ./... -> pass; go test ./... -> pass; go test -race -> FAIL (ThreadSanitizer unsupported VMA, no concurrency touched)
Unsure: the fail-fast fires on root-Open validation, not fs.Sub itself — embed.FS.Sub never errors for a missing dir in Go 1.26 (probed), so a literal Sub-error panic would never trigger; the Open check is the only thing catching the real layout-change case
Follow-ups: none
```
