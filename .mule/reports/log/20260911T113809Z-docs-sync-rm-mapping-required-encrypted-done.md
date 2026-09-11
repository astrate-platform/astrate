slug: docs-sync-rm-mapping-required-encrypted
verdict: done
at:  5f8cb25
ran: 2026-09-11T11:38:09Z on DietPi in 137s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_realm_management_api.yaml [offset=55, limit=140]
→ Read docs/api/astarte_realm_management_api.yaml [offset=800, limit=80]
✗ Read internal/interfaceschema/parse.go failed [limit=45, offset=385]
Error: File not found: /root/astrate-mule/internal/interfaceschema/parse.go
✗ Read internal/interfaceschema/compat.go failed [limit=30, offset=105]
Error: File not found: /root/astrate-mule/internal/interfaceschema/compat.go
Now let me verify the Go-side claims.
✗ Read internal/interfaceschema/detailed.go failed [limit=40, offset=50]
Error: File not found: /root/astrate-mule/internal/interfaceschema/detailed.go
✱ Glob "**/interfaceschema/**/*.go" 12 matches
→ Read pkg/interfaceschema/parse.go [limit=45, offset=380]
→ Read pkg/interfaceschema/compat.go [limit=30, offset=105]
✱ Grep "Required|Encrypted" 1 match
→ Read internal/realm/detailed.go [limit=45, offset=40]

← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -855,8 +855,16 @@
   default: false
 explicit_timestamp:
   type: boolean
   default: false
+required:
+  type: boolean
+  description: Datastream-only. Whether the mapping value is mandatory.
+  default: false
+encrypted:
+  type: boolean
+  description: Datastream-only. Whether the mapping value is stored encrypted.
+  default: false
 description:
   type: string
 doc:
   type: string

All claims verified. Making the edits.
← Edit docs/api/astarte_realm_management_api.yaml
Index: /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
+++ /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
@@ -81,8 +81,10 @@
                     reliability: unreliable
                     retention: discard
                     expiry: 0
                     explicit_timestamp: false
+                    required: false
+                    encrypted: false
                     database_retention_policy: no_ttl
 "401":
   $ref: "#/components/responses/Unauthorized"
 "500":

→ Read docs/Makefile
→ Read docs/api/astarte_realm_management_api.yaml [limit=30, offset=818]
→ Read docs/api/astarte_realm_management_api.yaml [limit=35, offset=847]
Now run the docs build to confirm the YAML still loads.

$ make -C docs build 2>&1; echo "EXIT=$?"
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
EXIT=2

→ Read docs/requirements.txt
mkdocs is missing. Let me check the requirements and whether site/ is tracked.

$ git check-ignore docs/site docs/site/index.html api 2>&1; echo "---"; git -C /root/astrate-mule status --porcelain | head -20
---
 M docs/api/astarte_realm_management_api.yaml

$ python3 --version; pip3 --version 2>&1 | head -1; ls /root/astrate-mule/docs | head; git -C /root/astrate-mule check-ignore docs/site; echo "ignore-exit=$?"; git -C /root/astrate-mule ls-files docs/site | head -5
Python 3.13.5
/bin/sh: 1: pip3: not found
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
ignore-exit=1
docs/site/appengine-api.md
docs/site/architecture.md
docs/site/compatibility.md
docs/site/configuration-reference.md
docs/site/contributing.md

$ python3 -m pip --version 2>&1 | head -1; which uv poetry pip 2>&1; python3 -c "import mkdocs" 2>&1 | head -1
/usr/bin/python3: No module named pip
Traceback (most recent call last):

$ ls ~/go/pkg/mod/gopkg.in/ 2>/dev/null; ls ~/go/pkg/mod/cache/download/gopkg.in/yaml.v3/@v/ 2>/dev/null | head; cat /root/astrate-mule/go.mod | head -30
check.v1@v1.0.0-20201130134442-10cb98267c6c
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1
list
v3.0.0-20200313102051-9f266ea9e77c.mod
v3.0.1.info
v3.0.1.lock
v3.0.1.mod
v3.0.1.zip
v3.0.1.ziphash
module github.com/astrate-platform/astrate

go 1.26.1

// Pinned dependency decisions — docs/ROADMAP.md §1.1 (frozen in docs/DESIGN.md).
// HTTP routing is stdlib net/http (Go >= 1.22 ServeMux patterns): deliberately no framework.
require (
	github.com/BurntSushi/toml v1.6.0 // TOML config (§5.1)
	github.com/coder/websocket v1.8.14 // live stream socket (M7)
	github.com/eclipse/paho.mqtt.golang v1.5.1 // test-only MQTT client (same family as official Go SDK)
	github.com/golang-jwt/jwt/v5 v5.3.1 // JWT verification, alg allowlist (§4.2)
	github.com/golang-migrate/migrate/v4 v4.19.1 // migrations via source/iofs + go:embed (§1.3)
	github.com/hashicorp/golang-lru/v2 v2.0.7 // verified-token LRU cache (§4.2)
	github.com/jackc/pgx/v5 v5.10.0 // Postgres/TimescaleDB driver + pgxpool (§1.3)
	github.com/mochi-mqtt/server/v2 v2.7.9 // embedded MQTT broker (§1.1)
	github.com/prometheus/client_golang v1.23.2 // metrics (§5.2)
	github.com/testcontainers/testcontainers-go v0.43.0 // T2 integration harness (§5.4)
	github.com/testcontainers/testcontainers-go/modules/postgres v0.42.0 // Timescale container module
	go.etcd.io/bbolt v1.4.3 // broker session store (§3.1; pinned over pebble)
	go.mongodb.org/mongo-driver/v2 v2.6.0 // BSON raw-document API only (§3.5.5)
	golang.org/x/crypto v0.53.0 // bcrypt cost 10 for credentials secrets (§4.1)
)

require (
	github.com/cbroglie/mustache v1.4.0
	github.com/nats-io/nats.go v1.52.0
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.3
	github.com/testcontainers/testcontainers-go/modules/nats v0.43.0
)

! permission requested: external_directory (/tmp/yamlprobe/*); auto-rejecting
✗ Write /tmp/yamlprobe/go.mod failed
Error: The user rejected permission to use this specific tool call.
! permission requested: external_directory (/tmp/yamlprobe/*); auto-rejecting
✗ Write /tmp/yamlprobe/main.go failed
Error: The user rejected permission to use this specific tool call.
yaml.v3 is in the module cache. Writing a throwaway probe in /tmp to confirm the Swagger-referenced YAML loads.
