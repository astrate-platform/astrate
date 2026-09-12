slug: docs-sync-ae-write-405
verdict: done
at:  9ed08c3
ran: 2026-09-12T11:22:34Z on DietPi in 400s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
✱ Glob "docs/api/astarte_appengine_api.yaml" 1 match

$ rg -n '405' docs/api/astarte_appengine_api.yaml; echo "---"; rg -n '404' docs/api/astarte_appengine_api.yaml | head -40
---
156:        "404":
201:        "404":
235:        "404":
285:        "404":
319:        "404":
362:        "404":
412:        "404":
460:        "404":
506:        "404":
540:        "404":
566:        "404":
606:        "404":
649:        "404":
692:        "404":
733:        "404":
762:        "404":
864:        "404":
910:        "404":
952:        "404":
980:        "404":
1026:        "404":
1051:        "404":
1079:        "404":
1118:        "404":
1162:        "404":
1206:        "404":
1248:        "404":
1278:        "404":

✱ Grep "ErrNotServerOwned|ErrNotAProperty|Cannot write to" in internal/appengine/http.go · 8 matches
→ Read docs/api/astarte_appengine_api.yaml [limit=120]
✱ Grep "devices-by-alias|/groups/|/interfaces/|putData|deleteData" in docs/api/astarte_appengine_api.yaml · 20 matches
→ Read docs/api/astarte_appengine_api.yaml [offset=369, limit=210]
→ Read docs/api/astarte_appengine_api.yaml [limit=230, offset=613]
→ Read docs/api/astarte_appengine_api.yaml [limit=180, offset=1125]
✱ Grep "MethodNotAllowed|NotAllowed|403" in docs/api/astarte_appengine_api.yaml · 13 matches
✱ Grep "^    (Unauthorized|NotFound|DeviceNotFound|BadRequest|InternalServerError|ValidationErrors):" in docs/api/astarte_appengine_api.yaml · 22 matches
→ Read docs/api/astarte_housekeeping_api.yaml [limit=75, offset=178]
→ Read docs/api/astarte_appengine_api.yaml [limit=80, offset=1613]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=15, offset=372]
→ Read internal/appengine/http.go [limit=110, offset=280]
✱ Grep "publishData|MethodPost|MethodPut|MethodDelete|HandleFunc" in internal/appengine · 15 matches
✱ Grep "putData\b|putDataByAlias|putDataInGroup|deleteData\b|MethodPost|MethodPut" in internal/appengine · 22 matches
→ Read docs/api/astarte_appengine_api.yaml [limit=285, offset=841]
→ Read docs/Makefile
All nine write ops confirmed reachable to `writeError`'s 405 branch (http.go:673-676, routed via http.go:52-79). Let me check the docs build gate before editing.

$ rg -n 'astarte_appengine_api|swagger|yaml' mkdocs.yml 2>/dev/null; echo "---"; ls docs/swagger-ui 2>/dev/null | head
---
index.html

$ rg -rn 'astarte_appengine_api' docs --include='*' -l 2>/dev/null; echo "---mkdocs---"; ls docs/
---mkdocs---
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
UPSTREAM-EXPERIMENTAL.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui

✱ Grep "api/astarte|api/" in docs/mkdocs.yml · 57 matches
Now I'll make the edits. Four edit patterns cover all nine ops, plus the component addition.

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -458,8 +458,10 @@
     "401":
       $ref: "#/components/responses/Unauthorized"
     "404":
       $ref: "#/components/responses/DeviceNotFound"
+    "405":
+      $ref: "#/components/responses/MethodNotAllowed"
     "500":
       $ref: "#/components/responses/InternalServerError"
 
 post:
@@ -504,8 +506,10 @@
     "401":
       $ref: "#/components/responses/Unauthorized"
     "404":
       $ref: "#/components/responses/DeviceNotFound"
+    "405":
+      $ref: "#/components/responses/MethodNotAllowed"
     "500":
       $ref: "#/components/responses/InternalServerError"
 
 delete:
@@ -690,8 +694,10 @@
     "401":
       $ref: "#/components/responses/Unauthorized"
     "404":
       $ref: "#/components/responses/DeviceNotFound"
+    "405":
+      $ref: "#/components/responses/MethodNotAllowed"
     "500":
       $ref: "#/components/responses/InternalServerError"
 
 post:
@@ -731,8 +737,10 @@
     "401":
       $ref: "#/components/responses/Unauthorized"
     "404":
       $ref: "#/components/responses/DeviceNotFound"
+    "405":
+      $ref: "#/components/responses/MethodNotAllowed"
     "500":
       $ref: "#/components/responses/InternalServerError"
 
 delete:

$ rg -n '"405"' docs/api/astarte_appengine_api.yaml
462:        "405":
510:        "405":
698:        "405":
741:        "405":

→ Read docs/api/astarte_appengine_api.yaml [limit=14, offset=454]
Pattern A applied to 4 ops. Now the delete-style block (2 ops):

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -154,8 +154,10 @@
       "401":
         $ref: "#/components/responses/Unauthorized"
       "404":
         $ref: "#/components/responses/DeviceNotFound"
+      "405":
+        $ref: "#/components/responses/MethodNotAllowed"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
   patch:
@@ -233,8 +235,10 @@
       "401":
         $ref: "#/components/responses/Unauthorized"
       "404":
         $ref: "#/components/responses/DeviceNotFound"
+      "405":
+        $ref: "#/components/responses/MethodNotAllowed"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
   patch:
@@ -317,8 +321,10 @@
       "401":
         $ref: "#/components/responses/Unauthorized"
       "404":
         $ref: "#/components/responses/DeviceNotFound"
+      "405":
+        $ref: "#/components/responses/MethodNotAllowed"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
 /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
@@ -542,8 +548,10 @@
       "401":
         $ref: "#/components/responses/Unauthorized"
       "404":
         $ref: "#/components/responses/DeviceNotFound"
+      "405":
+        $ref: "#/components/responses/MethodNotAllowed"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
 /appengine/v1/{realm}/devices/{device}/interfaces:
@@ -568,8 +576,10 @@
       "401":
         $ref: "#/components/responses/Unauthorized"
       "404":
         $ref: "#/components/responses/DeviceNotFound"
+      "405":
+        $ref: "#/components/responses/MethodNotAllowed"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
 # ── Interface data ──────────────────────────────────────────────────
@@ -768,8 +778,10 @@
       "401":
         $ref: "#/components/responses/Unauthorized"
       "404":
         $ref: "#/components/responses/DeviceNotFound"
+      "405":
+        $ref: "#/components/responses/MethodNotAllowed"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
 # ── Groups ──────────────────────────────────────────────────────────

$ rg -n '"405"' docs/api/astarte_appengine_api.yaml
158:        "405":
239:        "405":
325:        "405":
468:        "405":
516:        "405":
552:        "405":
580:        "405":
708:        "405":
751:        "405":
782:        "405":

Overshoot — pattern B matched GET ops too. Let me inspect every new `405` location.
$ for l in 158 239 325 468 516 552 580 708 751 782; do start=$((l-12)); end=$((l+4)); sed -n "${start},${end}p" docs/api/astarte_appengine_api.yaml | sed "s/^/$l+/" | rg 'operationId|"4|"5' ; echo "---"; done
158+        "401":
158+        "404":
158+        "405":
158+        "500":
---
239+        "401":
239+        "404":
239+        "405":
239+        "500":
---
325+        "401":
325+        "404":
325+        "405":
325+        "500":
---
468+        "400":
468+        "401":
468+        "404":
468+        "405":
468+        "500":
---
516+        "400":
516+        "401":
516+        "404":
516+        "405":
516+        "500":
---
552+        "401":
552+        "404":
552+        "405":
552+        "500":
---
580+        "401":
580+        "404":
580+        "405":
580+        "500":
---
708+        "400":
708+        "401":
708+        "404":
708+        "405":
708+        "500":
---
751+        "400":
751+        "401":
751+        "404":
751+        "405":
751+        "500":
---
782+        "401":
782+        "404":
782+        "405":
782+        "500":
---

$ awk '/operationId:/{op=$0} /"405":/{print op; op=""}' docs/api/astarte_appengine_api.yaml
operationId: getDevice
      operationId: getDeviceByAlias
      operationId: listInterfacesByAlias
      operationId: putDataByAlias
      operationId: publishDataByAlias
      operationId: deleteDataByAlias
      operationId: listDeviceInterfaces
      operationId: putData
      operationId: publishData
      operationId: deleteData

4 wrong GET ops got 405: getDevice, getDeviceByAlias, listInterfacesByAlias, listDeviceInterfaces. Reverting those (their block has `401` immediately before `404`, unlike the intended write ops):
$ rg -c -U '\s{8}"401":\n\s{10}\$ref: "#/components/responses/Unauthorized"\n\s{8}"404":\n\s{10}\$ref: "#/components/responses/DeviceNotFound"\n\s{8}"405":\n\s{10}\$ref: "#/components/responses/MethodNotAllowed"\n\s{8}"500":' docs/api/astarte_appengine_api.yaml
10

→ Read docs/api/astarte_appengine_api.yaml [limit=18, offset=146]
replaceAll would hit the intended delete ops too (they also have `401` before `404`). I need op-specific context. Let me read each wrong block's unique predecessor:
→ Read docs/api/astarte_appengine_api.yaml [limit=18, offset=227]
The two `get` blocks are byte-identical below; I'll revert the 4 wrong GET-op insertions surgically by op, via a small script.

$ python3 - <<'EOF'
import re
p = "docs/api/astarte_appengine_api.yaml"
s = open(p).read()
lines = s.split("\n")
out = []
i = 0
removed = []
while i < len(lines):
    if lines[i].strip() == '"405":' and lines[i+1].strip() == '$ref: "#/components/responses/MethodNotAllowed"':
        # walk back to nearest operationId
        j = len(out) - 1
        op = None
        while j >= 0:
            m = re.search(r'operationId:\s*(\w+)', out[j])
            if m:
                op = m.group(1)
                break
            j -= 1
        if op and (op.startswith("get") or op.startswith("list")):
            removed.append(op)
            i += 2
            continue
    out.append(lines[i])
    i += 1
open(p, "w").write("\n".join(out))
print("removed 405 from:", removed)
EOF
removed 405 from: ['getDevice', 'getDeviceByAlias', 'listInterfacesByAlias', 'listDeviceInterfaces']

$ awk '/operationId:/{op=$0} /"405":/{print op; op=""}' docs/api/astarte_appengine_api.yaml
operationId: putDataByAlias
      operationId: publishDataByAlias
      operationId: deleteDataByAlias
      operationId: putData
      operationId: publishData
      operationId: deleteData

6 intended ops in place; GET ops cleanly reverted. Now the 3 group twins use `NotFound` refs:
$ python3 - <<'EOF'
import re
p = "docs/api/astarte_appengine_api.yaml"
lines = open(p).read().split("\n")
targets = {"putDataInGroup", "publishDataInGroup", "deleteDataInGroup"}
out = []
cur_op = None
added = []
insert_block = ['        "405":', '          $ref: "#/components/responses/MethodNotAllowed"']
i = 0
while i < len(lines):
    line = lines[i]
    m = re.search(r'operationId:\s*(\w+)', line)
    if m:
        cur_op = m.group(1)
    if (cur_op in targets and line.strip() == '"404":'
            and lines[i+1].strip() == '$ref: "#/components/responses/NotFound"'
            and lines[i+2].strip() == '"500":'):
        out.extend([line, lines[i+1]] + insert_block)
        added.append(cur_op)
        i += 2
        continue
    out.append(line)
    i += 1
open(p, "w").write("\n".join(out))
print("added 405 to:", added)
EOF
added 405 to: ['putDataInGroup', 'publishDataInGroup', 'deleteDataInGroup']

$ awk '/operationId:/{op=$0} /"405":/{print op; op=""}' docs/api/astarte_appengine_api.yaml
operationId: putDataByAlias
      operationId: publishDataByAlias
      operationId: deleteDataByAlias
      operationId: putData
      operationId: publishData
      operationId: deleteData
      operationId: putDataInGroup
      operationId: publishDataInGroup
      operationId: deleteDataInGroup

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -1667,8 +1667,18 @@
       example:
         errors:
           detail: Device not found
 
+MethodNotAllowed:
+  description: Method Not Allowed — cannot write to a device-owned or read-only resource.
+  content:
+    application/json; charset=utf-8:
+      schema:
+        $ref: "#/components/schemas/ErrorDetail"
+      example:
+        errors:
+          detail: Cannot write to read-only resource
+
 InternalServerError:
   description: Internal Server Error — unexpected failure.
   content:
     application/json; charset=utf-8:

All nine ops done. Now the `MethodNotAllowed` response component:
Now verify YAML validity, `$ref` targets, and diff sanity, then run the docs build:
$ python3 - <<'EOF'
import yaml
d = yaml.safe_load(open("docs/api/astarte_appengine_api.yaml"))
resp = d["components"]["responses"]
refs = []
for path, ops in d["paths"].items():
    for meth, op in ops.items():
        for code, r in (op.get("responses") or {}).items():
            if isinstance(r, dict) and "$ref" in r:
                refs.append((path, meth, code, r["$ref"]))
bad = [r for r in refs if not r[3].startswith("#/components/") or r[3].split("/")[-1] not in resp]
mna = [r for r in refs if r[3] == "#/components/responses/MethodNotAllowed"]
print("total response refs:", len(refs))
print("dangling refs:", bad if bad else "none")
print("405 ops referencing MethodNotAllowed:", len(mna))
for r in mna: print("  ", r[1].upper(), r[0])
EOF
Traceback (most recent call last):
  File "<stdin>", line 1, in <module>
ModuleNotFoundError: No module named 'yaml'

$ python3 -c "import sys; print(sys.version)"; which python3; ls docs/requirements.txt && cat docs/requirements.txt
3.13.5 (main, Jul 15 2026, 20:25:40) [GCC 14.2.0]
/usr/bin/python3
docs/requirements.txt
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10

$ rg -n 'yaml' go.mod
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

$ mkdir -p /tmp/opencode/yamlprobe && cat > /tmp/opencode/yamlprobe/main.go <<'EOF'
package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	b, err := os.ReadFile("/root/astrate-mule/docs/api/astarte_appengine_api.yaml")
	if err != nil {
		panic(err)
	}
	var d struct {
		Components struct {
			Responses map[string]any `yaml:"responses"`
		} `yaml:"components"`
		Paths map[string]map[string]struct {
			Responses map[string]map[string]any `yaml:"responses"`
		} `yaml:"paths"`
	}
	if err := yaml.Unmarshal(b, &d); err != nil {
		fmt.Println("YAML PARSE FAIL:", err)
		os.Exit(1)
	}
	fmt.Println("YAML parse OK; responses:", len(d.Components.Responses))
	n := 0
	for path, ops := range d.Paths {
		for meth, op := range ops {
			for code, r := range op.Responses {
				if ref, ok := r["$ref"].(string); ok {
					name := ref[strings.LastIndex(ref, "/")+1:]
					if _, ok := d.Components.Responses[name]; !ok {
						fmt.Println("DANGLING:", name, "at", meth, path, code)
					}
					if name == "MethodNotAllowed" {
						n++
					}
				}
			}
		}
	}
	fmt.Println("ops referencing MethodNotAllowed:", n)
}
EOF
echo 'import "strings"' > /dev/null
(no output)

← Write /tmp/opencode/yamlprobe/main.go
Wrote file successfully.

$ go run /tmp/opencode/yamlprobe/main.go
YAML parse OK; components.responses: 9
  405 -> delete /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}
  405 -> put /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}
  405 -> post /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}
  405 -> put /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}
  405 -> post /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}
  405 -> delete /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}
  405 -> post /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}
  405 -> delete /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}
  405 -> put /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}
ops referencing MethodNotAllowed: 9

YAML valid, exactly 9 `MethodNotAllowed` refs, zero dangling refs. Now the build gate:
$ make -C docs build 2>&1
make: Entering directory '/root/astrate-mule/docs'
cp DESIGN.md site/DESIGN.md
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'

`mkdocs` binary isn't installed here; the `sync` half (which copies the Swagger-referenced YAML into `site/`) succeeded before failing on the mkdocs step. Let me confirm the copied YAML matches and is loadable:
$ cmp docs/api/astarte_appengine_api.yaml docs/site/api/astarte_appengine_api.yaml && echo "site copy identical"; rg -c '"405":' docs/site/api/astarte_appengine_api.yaml
site copy identical
9

$ sed -n '1210,1218p' docs/api/astarte_appengine_api.yaml; echo "..."; sed -n '/^    MethodNotAllowed:/,/^    InternalServerError:/p' docs/api/astarte_appengine_api.yaml
description: The value to set.
      responses:
        "200":
          description: Data set successfully.
        "400":
          $ref: "#/components/responses/BadRequest"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "404":
...
    MethodNotAllowed:
      description: Method Not Allowed — cannot write to a device-owned or read-only resource.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Cannot write to read-only resource

    InternalServerError:

$ git diff --stat docs/api/astarte_appengine_api.yaml docs/site/api/astarte_appengine_api.yaml
docs/api/astarte_appengine_api.yaml | 28 ++++++++++++++++++++++++++++
 1 file changed, 28 insertions(+)

$ git status --short docs/site 2>/dev/null; git check-ignore docs/site/api/astarte_appengine_api.yaml
docs/site/api/astarte_appengine_api.yaml

$ git diff docs/api/astarte_appengine_api.yaml
diff --git a/docs/api/astarte_appengine_api.yaml b/docs/api/astarte_appengine_api.yaml
index d0ec324..dfb3662 100644
--- a/docs/api/astarte_appengine_api.yaml
+++ b/docs/api/astarte_appengine_api.yaml
@@ -459,6 +459,8 @@ paths:
           $ref: "#/components/responses/Unauthorized"
         "404":
           $ref: "#/components/responses/DeviceNotFound"
+        "405":
+          $ref: "#/components/responses/MethodNotAllowed"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -505,6 +507,8 @@ paths:
           $ref: "#/components/responses/Unauthorized"
         "404":
           $ref: "#/components/responses/DeviceNotFound"
+        "405":
+          $ref: "#/components/responses/MethodNotAllowed"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -539,6 +543,8 @@ paths:
           $ref: "#/components/responses/Unauthorized"
         "404":
           $ref: "#/components/responses/DeviceNotFound"
+        "405":
+          $ref: "#/components/responses/MethodNotAllowed"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -691,6 +697,8 @@ paths:
           $ref: "#/components/responses/Unauthorized"
         "404":
           $ref: "#/components/responses/DeviceNotFound"
+        "405":
+          $ref: "#/components/responses/MethodNotAllowed"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -732,6 +740,8 @@ paths:
           $ref: "#/components/responses/Unauthorized"
         "404":
           $ref: "#/components/responses/DeviceNotFound"
+        "405":
+          $ref: "#/components/responses/MethodNotAllowed"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -761,6 +771,8 @@ paths:
           $ref: "#/components/responses/Unauthorized"
         "404":
           $ref: "#/components/responses/DeviceNotFound"
+        "405":
+          $ref: "#/components/responses/MethodNotAllowed"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -1205,6 +1217,8 @@ paths:
           $ref: "#/components/responses/Unauthorized"
         "404":
           $ref: "#/components/responses/NotFound"
+        "405":
+          $ref: "#/components/responses/MethodNotAllowed"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -1247,6 +1261,8 @@ paths:
           $ref: "#/components/responses/Unauthorized"
         "404":
           $ref: "#/components/responses/NotFound"
+        "405":
+          $ref: "#/components/responses/MethodNotAllowed"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -1277,6 +1293,8 @@ paths:
           $ref: "#/components/responses/Unauthorized"
         "404":
           $ref: "#/components/responses/NotFound"
+        "405":
+          $ref: "#/components/responses/MethodNotAllowed"
         "500":
           $ref: "#/components/responses/InternalServerError"
 
@@ -1650,6 +1668,16 @@ components:
             errors:
               detail: Device not found
 
+    MethodNotAllowed:
+      description: Method Not Allowed — cannot write to a device-owned or read-only resource.
+      content:
+        application/json; charset=utf-8:
+          schema:
+            $ref: "#/components/schemas/ErrorDetail"
+          example:
+            errors:
+              detail: Cannot write to read-only resource
+
     InternalServerError:
       description: Internal Server Error — unexpected failure.
       content:

Done: added the missing `405 Method Not Allowed` response to all nine data-write ops (PUT/POST/DELETE on `/devices/{device}/interfaces/{interface}/{path}`, the `/devices-by-alias/{alias}` triple, and the `/groups/{group}/devices/{device}` triple) via a new `MethodNotAllowed` component, matching the messages `writeError` emits (http.go:673-676) for `ErrNotServerOwned`/`ErrNotAProperty`.

Files: docs/api/astarte_appengine_api.yaml

Verified:
- yaml.v3 probe (in /tmp): file parses, exactly 9 ops → `MethodNotAllowed`, no dangling `$ref` → pass
- `make -C docs build`: sync step copied the Swagger-referenced YAML into `docs/site/api/`; final `mkdocs build` fails only because `mkdocs` is not installed on this machine (Error 127) — not a YAML problem
- my first draft over-reached and added 405 to 4 GET ops by a broad replace; caught it via op-anchored audit and reverted those 4 before finishing

Unsure: nothing (405 was in code for all nine ops — routes at http.go:52-79; GET ops do not hit the 405 branch, so they were correctly left alone; group ops keep their existing `NotFound` 404 ref, devices/by-alias keep `DeviceNotFound`). No Go behavior changed, so `writeerror_test.go` already covers the 405 cases, and no new test was needed for a docs-only change.

Follow-ups: none
