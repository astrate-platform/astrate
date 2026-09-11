slug: docs-sync-rm-put-auth-422
verdict: done
at:  8554720
ran: 2026-09-11T17:07:16Z on DietPi in 145s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_realm_management_api.yaml [offset=400, limit=90]
Task: docs-sync-rm-put-auth-422 — add a missing `422` (ValidationError) response to `PUT /realmmanagement/v1/{realm}/config/auth` in the realm management API spec, since `SetAuthKey` returns `ErrValidation` for a blank `jwt_public_key_pem`, which `writeError` maps to 422.

Let me read the relevant spec section and the code.
→ Read internal/realm/service.go [offset=610, limit=60]
✱ Grep "ErrValidation|ValidationError|422" 100 matches
The code confirms the claim. Let me verify how `ErrValidation` maps to 422 and check whether a ValidationError response component exists in the spec.
✱ Grep "ValidationError|"422"" in docs/api/astarte_realm_management_api.yaml · 42 matches
→ Read internal/realm/http.go [limit=50, offset=385]
The spec already defines a `ValidationError` component referenced at several endpoints. Let me confirm the realm `writeError` mapping and look at an existing usage for style.
→ Read docs/api/astarte_realm_management_api.yaml [limit=40, offset=120]
Confirmed: realm `writeError` maps `ErrValidation` → 422 (internal/realm/http.go:410-411), and `ValidationError` is an existing component in this spec (yaml:1101). Adding the missing 422 to the `put` block.

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -461,8 +461,10 @@
       "400":
         $ref: "#/components/responses/BadRequest"
       "401":
         $ref: "#/components/responses/Unauthorized"
+      "422":
+        $ref: "#/components/responses/ValidationError"
       "500":
         $ref: "#/components/responses/InternalServerError"
 
 /realmmanagement/v1/{realm}/config/device_registration_limit:

Now run the docs build as the task requires, plus the gates.
$ ls docs/ && ls docs/Makefile 2>/dev/null || echo "no docs Makefile"
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
docs/Makefile

$ make -C docs build 2>&1 | tail -30
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

→ Read docs/Makefile
`mkdocs` isn't installed. Let me check the Makefile and whether the YAML itself is valid with a tool that's available.

$ command -v python3; command -v yamllint; command -v ruby; python3 -c "import yaml" 2>&1 | head -2
/usr/bin/python3
Traceback (most recent call last):
  File "<string>", line 1, in <module>

$ rg -n "yaml" go.mod
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

Make's `sync` step ran (copied my change to `site/api`); `mkdocs build` failed only because `mkdocs` isn't installed on this machine. Let me validate the YAML still loads properly with a Go-based check, since the repo has a YAML dependency.
✱ Grep "astarte_realm_management_api|\.yaml" in docs · 1 match
→ Read docs/embed.go

$ rg -n "embed|RealmManagementApi|management_api|\.yaml" --glob '*.go' docs/ internal/ cmd/ pkg/ | grep -iv "_test.go" | head -30
cmd/astrate/main.go:2:// it wires the store, ingestion engine, embedded MQTT broker, pairing, the M7
cmd/astrate/main.go:415:	// Swagger UI + OpenAPI YAML specs embedded in the binary.
docs/embed.go:1:// Package docs embeds the Swagger UI static files and OpenAPI YAML specs so
docs/embed.go:5:import "embed"
docs/embed.go:9://go:embed swagger-ui/*
docs/embed.go:10:var SwaggerUI embed.FS
docs/embed.go:14://go:embed api/*.yaml
docs/embed.go:15:var APIYAML embed.FS
internal/flow/factory.go:223:	// when non-empty so store name wins over a stale embedded field.
internal/swagger/swagger.go:1:// Package swagger serves the embedded Swagger UI and OpenAPI YAML specs at
internal/swagger/swagger.go:25:	// Serve YAML files at /api/ so the relative ../api/*.yaml in index.html
internal/swagger/swagger.go:37:		if strings.HasSuffix(path, ".yaml") {
internal/config/config.go:69:// MQTTConfig is the embedded broker (§3.1).
pkg/payload/bson.go:24:// UTF-8 string, generic binary, UTC datetime, arrays, embedded documents).
pkg/payload/bson.go:134:// decodeBSONObject decodes an object-aggregation `v`: an embedded document
pkg/payload/bson.go:314:// appendBSONObject appends an object-aggregation embedded document. Keys are
internal/engine/engine.go:44:// AdaptBroker wraps the embedded broker as the engine's BrokerPort.
internal/housekeeping/service.go:4:// realm mints its embedded CA (docs/DESIGN.md §4.3) and seals the CA private
internal/store/store.go:4:// embedded schema migrations (docs/DESIGN.md §2.2–2.5), CA-key sealing, and
internal/store/store.go:67:// New connects to dsn, applies any pending embedded migrations, probes
internal/store/store.go:100:// runMigrations applies the embedded migrations up to head. It is a no-op on
internal/store/store.go:110:		return fmt.Errorf("store: loading embedded migrations: %w", err)
internal/broker/broker.go:67:// Broker is the embedded MQTT broker (docs/ROADMAP.md §6 file 5.8): mochi
internal/broker/intake.go:1:// Package broker embeds the mochi-mqtt server and gives it Astarte MQTT v1
internal/broker/intake.go:60:// the embedded MQTT broker. The frozen decision (2026-08-22) is in §1.4; the
internal/pairing/ca/ca.go:1:// Package ca implements Astrate's embedded per-realm certificate authority
internal/pairing/service.go:3:// credentials secrets, CSR-based credential issuance through the embedded

No repo code parses the YAML into a schema, so "still loads" means parse-valid. Let me probe-gate it with a throwaway Go program using the module's YAML parser, and run the code gates.
← Write /tmp/opencode/yamlprobe/main.go
Wrote file successfully.

$ go run /tmp/opencode/yamlprobe/main.go
# command-line-arguments
/tmp/opencode/yamlprobe/main.go:37:27: not enough arguments in call to yaml.Unmarshal
	have ([]byte)
	want ([]byte, interface{})

← Edit /tmp/opencode/yamlprobe/main.go
Index: /tmp/opencode/yamlprobe/main.go
===================================================================
--- /tmp/opencode/yamlprobe/main.go
+++ /tmp/opencode/yamlprobe/main.go
@@ -32,12 +32,9 @@
 }
 paths := m["paths"].(map[string]any)
 put := paths["/realmmanagement/v1/{realm}/config/auth"].(map[string]any)["put"].(map[string]any)
 prs := put["responses"].(map[string]any)
-buf := bytes.NewBuffer(nil)
-if err := yaml.Unmarshal([]byte("---\n")); err != nil {
-	_ = buf
-}
+var _ = bytes.NewBuffer
 for code, ref := range prs {
 	fmt.Printf("PUT /config/auth response %s -> %v\n", code, ref)
 }
 if _, ok := prs["422"]; !ok {

Probe has a stray bug; fixing it.
$ go run /tmp/opencode/yamlprobe/main.go
PUT /config/auth response 500 -> map[$ref:#/components/responses/InternalServerError]
PUT /config/auth response 204 -> map[description:Auth configuration updated. No body returned.]
PUT /config/auth response 400 -> map[$ref:#/components/responses/BadRequest]
PUT /config/auth response 401 -> map[$ref:#/components/responses/Unauthorized]
PUT /config/auth response 422 -> map[$ref:#/components/responses/ValidationError]
OK: YAML parses, ValidationError refexists, 422 present on PUT /config/auth

YAML parses and 422 is present, referencing the existing `ValidationError` component. Now the gates.
$ gofmt -l . 2>&1; echo "---vet---"; go vet ./... 2>&1 | tail -5; echo "---test---"; go test ./... 2>&1 | tail -15
---vet---
---test---
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	0.043s
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

Gates all pass, docs-only change done.

```
Done: Added missing 422 ValidationError response to PUT /realmmanagement/v1/{realm}/config/auth in docs/api/astarte_realm_management_api.yaml
Files: docs/api/astarte_realm_management_api.yaml
Verified: go vet ./... -> pass; go test ./... -> pass; gofmt -l . -> clean (no output); make -C docs build -> partial: sync step ran (site/api updated), mkdocs build itself failed because `mkdocs` is not installed on this host; substituted with Go/yaml.v3 probe that parses the spec and confirms the 422 ref exists and ValidationError component is present -> pass
Unsure: mkdocs could not be run (binary absent); YAML validity confirmed via gopkg.in/yaml.v3 parse instead
Follow-ups: none
```
