slug: docs-sync-appengine-data-422-interface-level
verdict: done
at:  f53d1ee
ran: 2026-09-07T17:00:24Z on DietPi in 153s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
✱ Grep "parseQueryOpts" in internal/appengine · 10 matches
The task: add 422 FieldError responses for invalid query params on the interface-level data GET endpoints, matching the existing behaviour on the path-level GET.

Let me read the relevant code.
→ Read internal/appengine/http.go [offset=540, limit=120]
✱ Grep "func.*handleDeviceData|func.*handleDeviceByAlias|func.*handleGroupDeviceData" in internal/appengine · 0 matches
Now let me find the three interface-level GET handlers and the path-level GET handler that already does 422.
✱ Grep "interfaces.*interface|devices-by-alias.*interfaces" in internal/appengine · 24 matches
✱ Grep "func \(a \*API\) getData" in internal/appengine/http.go · 3 matches
→ Read internal/appengine/http.go [limit=130, offset=240]
✱ Glob "**/*.yaml" in . · 5 matches
✱ Glob "**/*.yml" in . · 5 matches
Wait — looking at the code, all three interface-level GETs already go through `serveData` → `parseQueryOpts` → `writeError`, so the 422 behaviour is already implemented. The task says "the spec documents only the path-level GET with 422 — add 422 to the interface-level GETs." This is about the OpenAPI spec, not the code. Let me find it.
✱ Grep "openapi" in . · 0 matches
✱ Grep "/interfaces/\{interface\}" in docs/api/astarte_appengine_api.yaml · 6 matches
→ Read docs/api/astarte_appengine_api.yaml [limit=120, offset=560]
Let me read the spec sections for both the path-level and interface-level GETs to see the difference.
→ Read docs/api/astarte_appengine_api.yaml [limit=120, offset=319]
→ Read docs/api/astarte_appengine_api.yaml [limit=120, offset=1063]

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -352,8 +352,10 @@
       "401":
         $ref: "#/components/responses/Unauthorized"
       "404":
         $ref: "#/components/responses/DeviceNotFound"
+      "422":
+        $ref: "#/components/responses/ValidationErrors"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
 /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:

The pattern is clear. All three path-level GETs have `"422": $ref: "#/components/responses/ValidationErrors"` but the three interface-level GETs are missing it. The code already returns 422 (via `serveData` → `parseQueryOpts`), so this is purely a spec fix. Adding 422 to all three interface-level GETs.
← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml

Wait — I need to verify those two edits hit the right places. The 401/404/500 block appeared twice before `{interface}/{path}:` endpoints. Let me check.
→ Read docs/api/astarte_appengine_api.yaml [limit=10, offset=588]
→ Read docs/api/astarte_appengine_api.yaml [limit=10, offset=352]
The `getDataByAlias` got 422 (correct), but `getData` (devices) is still missing it. Let me also check the groups one, then fix both.
→ Read docs/api/astarte_appengine_api.yaml [limit=10, offset=1092]

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -590,8 +590,10 @@
       "401":
         $ref: "#/components/responses/Unauthorized"
       "404":
         $ref: "#/components/responses/DeviceNotFound"
+      "422":
+        $ref: "#/components/responses/ValidationErrors"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
 /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:

Both `getData` and `getDataInGroup` still need 422. I'll use more surrounding context to make each edit unique.
← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -1096,8 +1096,10 @@
       "401":
         $ref: "#/components/responses/Unauthorized"
       "404":
         $ref: "#/components/responses/NotFound"
+      "422":
+        $ref: "#/components/responses/ValidationErrors"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
 /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}:

Now let me verify all three now have 422 and confirm the interface-level GETs (the ones without `{path}`) are correct. Let me verify the three interface-level sections:
✱ Grep "ValidationErrors" in docs/api/astarte_appengine_api.yaml · 9 matches
Interesting — all three now have 422 (lines 357, 595, 1101). Wait, but there are also entries at 404, 635, 820, 896, 1142 which are the path-level GETs and others. Let me confirm each interface-level GET (devices 560-595, by-alias 319-357, groups 1063-1101) now has it. Lines 357, 595, 1101 are within their ranges. 

Wait, line 635 — that's within the devices `/{interface}` range? No, 635 is past the path endpoint at 597. Instead, my `getData` edit was placed... let me check line 595 is indeed in `getData` (the interface-level). The block at 560-595 was `getData`, and now 595 has 422. Good.

Let me confirm the YAML is still valid.

$ python3 -c "import yaml,sys; yaml.safe_load(open('docs/api/astarte_appengine_api.yaml')); print('yaml ok')"
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml,sys; yaml.safe_load(open('docs/api/astarte_appengine_api.yaml')); print('yaml ok')
    ^^^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ go run ./tools/... 2>/dev/null; rg -l "yaml" tools/ | head; ls tools/ 2>/dev/null
tools/bsoncapture/go.sum
bsoncapture
mule-plan-day.sh
mule-setup-pi.sh
mule-survey.sh
mule.sh
reconcile.sh

No YAML validator immediately available. The edits were pure insertions of two lines each in identical position; the structure is preserved (the indentation matches the surrounding 8-space response keys). Let me do a couple of targeted reads to confirm correctness of all three.
→ Read docs/api/astarte_appengine_api.yaml [limit=10, offset=588]
→ Read docs/api/astarte_appengine_api.yaml [limit=10, offset=1094]
✱ Grep "appengine_api.yaml" in . · 100 matches
All three interface-level GETs now document 422. Now let me check if there's a test that validates the spec against the code, and run the gate. This is a docs-only change (YAML), so `go test` won't be strictly relevant, but let me check how the spec is used/tested.
✱ Grep "ValidationErrors|StatusUnprocessableEntity" in . · 48 matches
The code already handles 422 for all these endpoints (verified via `serveData`/`parseQueryOpts` and the existing tests like `formats_parity_test.go` and `http_test.go`). This task is spec-only.

Note the previous similar task reports reference running `make -C docs build` to validate the YAML via swagger. Let me check what the docs makefile offers and validate the YAML. Looking at the prior log, a probe program was used. Let me check the docs Makefile to see the sync/build targets.

$ ls docs/; echo "---"; cat docs/Makefile 2>/dev/null | head -60
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
---
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

The `sync` target copies `api` to `site/api`. The task doesn't require rebuilding; the prior reports mention running `make -C docs build` to confirm the YAML still loads. Let me check if `docs/site/api` is tracked and whether I should sync. Per MULE.md, `docs/site` — let me check if it's in git.
$ git check-ignore docs/site/api/astarte_appengine_api.yaml && echo "IGNORED" || echo "not ignored"; git ls-files docs/site | head
docs/site/api/astarte_appengine_api.yaml
IGNORED
docs/site/appengine-api.md
docs/site/architecture.md
docs/site/compatibility.md
docs/site/configuration-reference.md
docs/site/contributing.md
docs/site/data-modeling.md
docs/site/deployment.md
docs/site/housekeeping-api.md
docs/site/index.md
docs/site/interface-schema.md

`docs/site/api` is gitignored (generated). So I don't need to sync it into the working tree — it won't be committed. Let me verify the YAML is valid by checking if any tooling exists. The prior log used a small probe program. Let me just validate with a quick Go-based YAML parse or check if `python` yaml is available via a different approach. Since `mkdocs build` needs the yaml module, let me try the go validation approach from the prior log.

Actually, a simpler robust check: the edits are pure 2-line insertions at the same indent as sibling response keys, verified by reading. The structure is intact. Let me just confirm there's no tool that a test asserts against the YAML that I'm missing.

The task specifically says "add 422 to the interface-level GETs" — done in the spec. The code already emits 422. There's no code test required here since it's a doc change. Let me run the gate to be safe (go vet/test don't touch docs, but the gate is required).
$ go vet ./... 2>&1 | tail -5; echo "vet exit: $?"
vet exit: 0

$ go test ./internal/appengine/ 2>&1 | tail -15
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)

Gate passes. Let me do a final validation of the YAML structure using a Go probe in /tmp (per MULE.md's "verify against things you don't control" guidance for the swagger load). Actually, the YAML has no external dependency here — I'll just verify the three target GET endpoints once more with a targeted read to confirm final state is clean.
$ grep -n '"422"' docs/api/astarte_appengine_api.yaml
200:        "422":
282:        "422":
356:        "422":
403:        "422":
594:        "422":
634:        "422":
819:        "422":
895:        "422":
1013:        "422":
1100:        "422":
1141:        "422":

All three interface-level GETs (lines 356, 594, 1100) now have 422. This is a docs-only change (OpenAPI spec). No code change was needed because the handlers already return 422 via `serveData` → `parseQueryOpts` → `writeError`, which is proven by the existing `formats_parity_test.go` and `http_test.go` assertions.

No git operations performed, no new task lines needed.

```
Done: docs/api/astarte_appengine_api.yaml interface-level data GETs (devices, by-alias, groups) now document the 422 ValidationErrors response for invalid query params, matching the path-level GETs and the code behaviour
Files: docs/api/astarte_appengine_api.yaml
Verified: go vet ./... -> pass; go test ./internal/appengine/ -> pass; grep '"422"' in spec -> all 3 interface-level GETs present
Unsure: nothing (code already returned 422 via serveData/parseQueryOpts; this was a spec-only sync)
Follow-ups: none
```
