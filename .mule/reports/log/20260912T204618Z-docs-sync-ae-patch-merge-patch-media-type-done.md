slug: docs-sync-ae-patch-merge-patch-media-type
verdict: done
at:  64361d9
ran: 2026-09-12T20:46:18Z on DietPi in 137s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_appengine_api.yaml [offset=150, limit=140]
→ Read internal/appengine/http.go [offset=200, limit=40]
→ Read docs/api/astarte_appengine_api.yaml [limit=40, offset=995]
→ Read internal/appengine/http.go [limit=20, offset=512]

$ rg -n "required: \[data\]" docs/api/astarte_appengine_api.yaml
115:                required: [data]
179:              required: [data]
263:              required: [data]
449:              required: [data]
499:              required: [data]
691:              required: [data]
736:              required: [data]
832:              required: [data]
848:                required: [data]
955:              required: [data]
1028:              required: [data]
1219:              required: [data]
1265:              required: [data]
1559:      required: [data]
1566:      required: [data]
1595:      required: [data]
1603:      required: [data]

$ rg -n "merge-patch" docs/api/astarte_appengine_api.yaml
(no output)

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -172,9 +172,9 @@
   - $ref: "#/components/parameters/DeviceID"
 requestBody:
   required: true
   content:
-    application/json; charset=utf-8:
+    application/merge-patch+json:
       schema:
         type: object
         required: [data]
         properties:

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -256,9 +256,9 @@
       type: string
 requestBody:
   required: true
   content:
-    application/json; charset=utf-8:
+    application/merge-patch+json:
       schema:
         type: object
         required: [data]
         properties:

All three PATCH request bodies declared at 10-space indent match, and all three handlers enforce `application/merge-patch+json`. Responses stay as-is. Editing each request body individually.
← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -1021,9 +1021,9 @@
   - $ref: "#/components/parameters/DeviceID"
 requestBody:
   required: true
   content:
-    application/json; charset=utf-8:
+    application/merge-patch+json:
       schema:
         type: object
         required: [data]
         properties:

$ rg -n "merge-patch" docs/api/astarte_appengine_api.yaml
176:          application/merge-patch+json:
260:          application/merge-patch+json:
1025:          application/merge-patch+json:

$ make -C docs build
make: Entering directory '/root/astrate-mule/docs'
cp DESIGN.md site/DESIGN.md
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: Leaving directory '/root/astrate-mule/docs'
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127

→ Read docs/Makefile

$ python3 -c "import yaml,sys; yaml.safe_load(open('docs/api/astarte_appengine_api.yaml')); print('yaml ok')"
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml,sys; yaml.safe_load(open('docs/api/astarte_appengine_api.yaml')); print('yaml ok')
    ^^^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

`make -C docs build` fails only at the final `mkdocs` step because `mkdocs` isn't installed here. The copies succeeded. Let me verify the YAML is still parseable and check what references it.
$ rg -n "astarte_appengine_api.yaml" docs --glob '!api/**' --glob '!site/**'
docs/swagger-ui/index.html:47:      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
docs/site/swagger-ui/index.html:47:      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
docs/handoff/phase-1-memory.md:19:| `astarte_appengine_api.yaml` | 15 | 892 | Device list/detail/patch, data CRUD, groups CRUD, `a_aea` JWT, pagination links |
docs/handoff/phase-1-memory.md:55:├── astarte_appengine_api.yaml       (15 op, 892 righe)

$ command -v yq ruby perl node python3 || true
(no output)

$ which python3; python3 --version; ls /usr/bin | rg -i "yaml|json|python" || true
/usr/bin/python3
Python 3.13.5
aarch64-linux-gnu-python3-config
aarch64-linux-gnu-python3.13-config
json_pp
python3
python3-config
python3.13
python3.13-config
spa-json-dump

$ rg -n "yaml" go.mod go.sum | head -5
go.mod:93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
go.mod:99:	gopkg.in/yaml.v3 v3.0.1 // indirect
go.sum:202:go.yaml.in/yaml/v2 v2.4.2 h1:DzmwEr2rDGHl7lsFgAHxmNz/1NlQ7xLIrlN2h5d1eGI=
go.sum:203:go.yaml.in/yaml/v2 v2.4.2/go.mod h1:081UH+NErpNdqlCXm3TtEran0rJZGxAYx9hb/ELlsPU=
go.sum:225:gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=

PyYAML isn't installed, but the Go module cache has `gopkg.in/yaml.v3`. Using a throwaway probe in /tmp:
$ mkdir -p /tmp/opencode/yamlprobe && cd /tmp/opencode/yamlprobe && printf 'module yamlprobe\n\ngo 1.22\n\nrequire gopkg.in/yaml.v3 v3.0.1\n' > go.mod && printf 'package main\n\nimport (\n\t"fmt"\n\t"os"\n\n\t"gopkg.in/yaml.v3"\n)\n\nfunc main() {\n\tdata, err := os.ReadFile("/root/astrate-mule/docs/api/astarte_appengine_api.yaml")\n\tif err != nil { fmt.Println("read err:", err); os.Exit(1) }\n\tvar doc yaml.Node\n\tif err := yaml.Unmarshal(data, &doc); err != nil { fmt.Println("parse err:", err); os.Exit(1) }\n\tfmt.Println("yaml ok, top level nodes:", len(doc.Content[0].Content))\n}\n' > main.go && go run . 2>&1 | tail -3
main.go:7:2: missing go.sum entry for module providing package gopkg.in/yaml.v3 (imported by yamlprobe); to add:
	go get yamlprobe

$ cd /tmp/opencode/yamlprobe && GOFLAGS=-mod=mod go mod tidy >/dev/null 2>&1; go run . 2>&1 | tail -3
yaml ok, top level nodes: 12

YAML parses cleanly. Verifying only the intended file changed and the swagger reference still resolves:
$ git diff --stat; echo ---; git diff
docs/api/astarte_appengine_api.yaml | 6 +++---
 1 file changed, 3 insertions(+), 3 deletions(-)
---
diff --git a/docs/api/astarte_appengine_api.yaml b/docs/api/astarte_appengine_api.yaml
index f9fe644..833d4e8 100644
--- a/docs/api/astarte_appengine_api.yaml
+++ b/docs/api/astarte_appengine_api.yaml
@@ -173,7 +173,7 @@ paths:
       requestBody:
         required: true
         content:
-          application/json; charset=utf-8:
+          application/merge-patch+json:
             schema:
               type: object
               required: [data]
@@ -257,7 +257,7 @@ paths:
       requestBody:
         required: true
         content:
-          application/json; charset=utf-8:
+          application/merge-patch+json:
             schema:
               type: object
               required: [data]
@@ -1022,7 +1022,7 @@ paths:
       requestBody:
         required: true
         content:
-          application/json; charset=utf-8:
+          application/merge-patch+json:
             schema:
               type: object
               required: [data]

$ ls -la docs/api/astarte_appengine_api.yaml
-rw-r--r-- 1 root root 56583 Sep 12 22:44 docs/api/astarte_appengine_api.yaml

Done: Swagger requestBody media type for the three AppEngine PATCH ops changed to `application/merge-patch+json`, matching the exact content-type headers the handlers enforce (internal/appengine/http.go:209, 230, 521).
Files: docs/api/astarte_appengine_api.yaml
Verified: `make -C docs build` -> fail (sync steps ok, `mkdocs` binary not installed here); Go yaml.v3 probe parsing the file -> pass ("yaml ok", 12 top-level nodes); `git diff` confirms only the 3 expected lines changed. Swagger reference `docs/swagger-ui/index.html -> ../api/astarte_appengine_api.yaml` resolves.
Unsure: mkdocs is unavailable on this runner, so the full mkdocs build couldn't complete; file was validated with a standalone YAML parse instead.
Follow-ups: none
