slug: docs-sync-hk-realm-name-pattern
verdict: done
at:  e728217
ran: 2026-10-04T11:40:11Z on DietPi in 525s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ ls docs && ls docs/api 2>/dev/null | head -30
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
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

$ rg -n "realm_name" docs/api/astarte_housekeeping_api.yaml
74:                realm_name: test
89:                  realm_name: test
127:                  realm_name: test
275:      required: [realm_name, jwt_public_key_pem]
277:        realm_name:
296:      required: [realm_name, jwt_public_key_pem]
298:        realm_name:
467:              detail: "realm_name can't be blank"

→ Read docs/api/astarte_housekeeping_api.yaml [offset=180, limit=160]
→ Read docs/api/astarte_housekeeping_api.yaml [offset=340]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=180]

$ rg -n "name" migrations/000002_metadata.up.sql | head -20
7:    name             text NOT NULL UNIQUE CHECK (name ~ '^[a-z][a-z0-9]*$'),
22:    name          text     GENERATED ALWAYS AS (definition->>'interface_name') STORED,
29:    UNIQUE (realm_id, name, major_version)
79:    name      text NOT NULL,
80:    UNIQUE (realm_id, name)
94:    name       text NOT NULL,
96:    UNIQUE (realm_id, name)

→ Read internal/store/realms.go [limit=110, offset=40]

$ rg -rn "ErrInvalidRealmName" --glob '*.go' | head -20
internal/housekeeping/service.go:		if errors.Is(err, store.n) {
internal/store/store.go:	// n reports a realm name rejected by the schema's
internal/store/store.go:	n = errors.New("store: invalid realm name")
internal/store/realms_test.go:			if !errors.Is(err, n) {
internal/store/realms_test.go:				t.Errorf("name %q: got %v, want n", name, err)
internal/store/realms.go:		return nil, fmt.Errorf("%w: %q", n, nr.Name)

→ Read internal/housekeeping/service.go [limit=80, offset=125]

$ rg -n "validationDetail" -A 22 internal/http/http.go | head -60
rg: internal/http/http.go: IO error for operation on internal/http/http.go: No such file or directory (os error 2)

$ rg -ln "func validationDetail" . | head;
./internal/housekeeping/http.go
./internal/flowapi/http.go
./internal/appengine/http.go
./internal/realm/http.go

→ Read internal/housekeeping/http.go [limit=45, offset=230]

$ rg -n "pattern:" docs/api/*.yaml | head -30
docs/api/astarte_pairing_api.yaml:390:          pattern: '^[A-Za-z0-9_-]{22}$'
docs/api/astarte_realm_management_api.yaml:945:        pattern: '^[A-Za-z0-9_-]{22}$'
docs/api/astarte_realm_management_api.yaml:992:          pattern: '^([a-zA-Z][a-zA-Z0-9]*\.([a-zA-Z0-9][a-zA-Z0-9-]*\.)*)?[a-zA-Z][a-zA-Z0-9]*$'
docs/api/astarte_realm_management_api.yaml:1158:          pattern: '^[a-zA-Z0-9_.~-]{1,128}$'

$ rg -n "device_id" docs/api/astarte_realm_management_api.yaml | head -30
1103:        device_id:

→ Read docs/api/astarte_realm_management_api.yaml [limit=30, offset=940]
→ Read docs/api/astarte_realm_management_api.yaml [limit=80, offset=1085]

$ rg -n "malformed|off-pattern|does not exist rather" docs/api/*.yaml | head -20
docs/api/astarte_appengine_api.yaml:1725:      description: Bad Request — malformed or missing data envelope.
docs/api/astarte_pairing_api.yaml:579:      description: Bad Request — malformed or missing data envelope.
docs/api/astarte_realm_management_api.yaml:938:        and a malformed ID is not reported as malformed: it is answered with
docs/api/astarte_realm_management_api.yaml:1310:      description: Bad Request — malformed or missing data envelope.
docs/api/astarte_housekeeping_api.yaml:378:      description: Bad Request — malformed or missing data envelope.

→ Read docs/api/astarte_realm_management_api.yaml [limit=25, offset=920]

$ rg -n "ErrInvalidRealmName|invalid realm name" -B4 -A12 internal/store/realms_test.go | head -60
114-				Name:               name,
115-				CACertificatePEM:   "x",
116-				CAPrivateKeySealed: []byte("x"),
117-			})
118:			if !errors.Is(err, ErrInvalidRealmName) {
119:				t.Errorf("name %q: got %v, want ErrInvalidRealmName", name, err)
120-			}
121-		}
122-	})
123-
124-	t.Run("Updates", func(t *testing.T) {
125-		realm := mustCreateRealm(t, s)
126-
127-		if err := s.SetRealmJWTPublicKeys(ctx, realm.Name, []string{"new-key"}); err != nil {
128-			t.Fatalf("SetRealmJWTPublicKeys: %v", err)
129-		}
130-		if err := s.SetRealmCA(ctx, realm.Name, "new-ca", []byte("new-sealed")); err != nil {
131-			t.Fatalf("SetRealmCA: %v", err)

$ cat docs/Makefile; echo ---; ls docs/swagger-ui | head; echo ---; rg -ln "docs/api" --glob '*.go' | head
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
---
index.html
---

→ Read internal/store/realms_test.go [limit=125]

$ rg -ln "api/astarte_housekeeping_api|astarte_housekeeping_api.yaml" --glob '!docs/**' . | head -20; echo ---; which mkdocs; mkdocs --version 2>&1 | head -3; python3 -c "import yaml; print('pyyaml ok')" 2>&1 | head -2
./internal/swagger/swagger_test.go
---
/bin/sh: 1: mkdocs: not found
Traceback (most recent call last):
  File "<string>", line 1, in <module>

$ ls internal/swagger; python3 -V; pip3 --version 2>&1|head -2
swagger.go
swagger_test.go
Python 3.13.5
/bin/sh: 1: pip3: not found

→ Read internal/swagger/swagger_test.go

$ cat go.mod; echo ---; ls $(go env GOMODCACHE) 2>/dev/null | head -20
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

require (
	dario.cat/mergo v1.0.2 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/containerd/platforms v0.2.1 // indirect
	github.com/cpuguy83/dockercfg v0.3.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/go-connections v0.6.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/ebitengine/purego v0.10.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jackc/pgerrcode v0.0.0-20220416144525-469b46aa5efa // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/klauspost/compress v1.18.5 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/magiconair/properties v1.8.10 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/go-archive v0.2.0 // indirect
	github.com/moby/moby/api v1.54.2 // indirect
	github.com/moby/moby/client v0.4.0 // indirect
	github.com/moby/patternmatcher v0.6.1 // indirect
	github.com/moby/sys/sequential v0.6.0 // indirect
	github.com/moby/sys/user v0.4.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nats-io/nkeys v0.4.15 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/rs/xid v1.4.0 // indirect
	github.com/shirou/gopsutil/v4 v4.26.5 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/tklauser/go-sysconf v0.3.16 // indirect
	github.com/tklauser/numcpus v0.11.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0 // indirect
	go.opentelemetry.io/otel v1.41.0 // indirect
	go.opentelemetry.io/otel/metric v1.41.0 // indirect
	go.opentelemetry.io/otel/trace v1.41.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
---
4d63.com
cache
charm.land
codeberg.org
dario.cat
dev.gaijin.team
github.com
gitlab.com
go-simpler.org
go.augendre.info
go.etcd.io
go.mongodb.org
go.opentelemetry.io
go.uber.org
go.yaml.in
golang.org
google.golang.org
gopkg.in
gotest.tools
honnef.co

$ python3 -m pip --version 2>&1 | head -3; ls docs/requirements.txt && cat docs/requirements.txt; ls docs/site | head; python3 -c "import mkdocs" 2>&1|tail -1
/usr/bin/python3: No module named pip
docs/requirements.txt
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10
DESIGN.md
ROADMAP.md
api
appengine-api.md
architecture.md
compatibility.md
configuration-reference.md
contributing.md
data-modeling.md
deployment.md
ModuleNotFoundError: No module named 'mkdocs'

$ rg -n "url|spec" docs/swagger-ui/index.html | head -30
42:    <label for="spec-select">Specification:</label>
43:    <select id="spec-select">
51:    <input type="text" id="server-url" value="http://localhost:8080" placeholder="e.g. http://localhost:8080">
58:    const specSelect = document.getElementById('spec-select');
59:    const serverUrl = document.getElementById('server-url');
62:    function initSwagger(specUrl) {
65:        url: specUrl,
83:      const url = serverUrl.value.replace(/\/+$/, '');
89:        nativeSetter.call(input, url);
94:    specSelect.addEventListener('change', function() {
99:    initSwagger(specSelect.value);

$ python3 -m ensurepip --version 2>&1 | tail -2; curl -sS -m 6 -o /dev/null -w 'pypi:%{http_code}\n' https://pypi.org/simple/mkdocs/ 2>&1 | tail -2
pip 25.1.1
pypi:200

$ rg -n "responses/ValidationError" docs/api/astarte_housekeeping_api.yaml; echo ---; rg -n "can't be blank|is invalid" docs/api/astarte_housekeeping_api.yaml
104:          $ref: "#/components/responses/ValidationError"
190:          $ref: "#/components/responses/ValidationErrors"
---
467:              detail: "realm_name can't be blank"

$ rg -n "createPolicy:" -A 60 docs/api/astarte_realm_management_api.yaml | sed -n '30,70p'
(no output)

$ rg -n "examples:" docs/api/astarte_realm_management_api.yaml | head -20
61:              examples:
147:              examples:
283:              examples:
336:              examples:
429:              examples:
803:              examples:
1365:          examples:

→ Read docs/api/astarte_realm_management_api.yaml [limit=50, offset=265]

$ rg -n "CommonName|func Generate" -A 6 internal/ca/*.go | head -50
rg: internal/ca/*.go: IO error for operation on internal/ca/*.go: No such file or directory (os error 2)

$ rg -ln "func Generate" --glob '*.go' pkg internal | head; echo ---; rg -n "CommonName" --glob '*.go' pkg internal | head
internal/pairing/ca/ca.go
---
internal/testutil/mqttclient.go:35:		Subject:               pkix.Name{CommonName: "astrate-test-broker"},
internal/testutil/mqttclient.go:68:		Subject: pkix.Name{CommonName: "ignored"},
internal/pairing/service_test.go:165:		Subject: pkix.Name{CommonName: "ignored"},
internal/pairing/service_test.go:391:	if got, want := cert.Subject.CommonName, "test/"+hwID; got != want {
internal/broker/authhook.go:362:	cn := peers[0].Subject.CommonName
internal/broker/authhook.go:398:		identity, err := ParseCN(leaf.Subject.CommonName)
internal/pairing/http_test.go:298:		if got, want := cert.Subject.CommonName, f.realmName+"/"+hwID; got != want {
internal/pairing/ca/ca.go:88:			CommonName:   "Astrate Realm " + realm + " CA",
internal/pairing/ca/ca.go:193:		Subject:               pkix.Name{CommonName: realm + "/" + deviceID},
internal/pairing/ca/ca_test.go:23:			CommonName:         "evilrealm/evil-device",

→ Read internal/pairing/ca/ca.go [limit=50, offset=60]

$ rg -n "Detail[A-Za-z]+ *=" pkg/astarteapi/envelope.go | head -20
37:	DetailBadRequest = "Bad request"
39:	DetailUnauthorized = "Unauthorized"
41:	DetailForbidden = "Forbidden"
43:	DetailNotFound = "Not Found"
46:	DetailDeviceNotFound = "Device not found"
51:	DetailRouteNotFound = "Not found"
54:	DetailPageNotFound = "Page not found"
58:	DetailInternalServerError = "Internal server error"

$ mkdir -p /tmp/opencode/regexprobe && cat > /tmp/opencode/regexprobe/main.go <<'EOF'
package main

import (
	"fmt"
	"regexp"
)

func main() {
	re := regexp.MustCompile(`^[a-z][a-z0-9]*$`)
	for _, s := range []string{"test", "production", "a", "t3st", "TEST", "test_realm", "test-realm", "1test", "tëst", "test realm", "", "realm2", "re-alm", "re_alm"} {
		fmt.Printf("%-12q -> %v\n", s, re.MatchString(s))
	}
}
EOF
go run /tmp/opencode/regexprobe/main.go
"test"       -> true
"production" -> true
"a"          -> true
"t3st"       -> true
"TEST"       -> false
"test_realm" -> false
"test-realm" -> false
"1test"      -> false
"tëst"       -> false
"test realm" -> false
""           -> false
"realm2"     -> true
"re-alm"     -> false
"re_alm"     -> false

Verified the code chain. Now the edits:
← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -231,11 +231,20 @@
 RealmName:
   name: realm
   in: path
   required: true
-  description: The realm name.
+  description: |
+    The realm name: a lowercase ASCII word starting with a letter, so
+    `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm`
+    are not realm names — no realm can be created under one. This parameter
+    is the read path though, and there an off-pattern name is not reported
+    as malformed: it is answered with the same `404` and
+    `{"errors":{"detail":"Not Found"}}` body as a realm that does not
+    exist, because the lookup matches the name verbatim with no format check
+    of its own.
   schema:
     type: string
+    pattern: '^[a-z][a-z0-9]*$'
 
 AsyncOperation:
   name: async_operation
   in: query

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -305,9 +305,18 @@
 required: [realm_name, jwt_public_key_pem]
 properties:
   realm_name:
     type: string
-    description: Unique realm identifier.
+    pattern: '^[a-z][a-z0-9]*$'
+    description: |
+      Unique realm identifier: a lowercase ASCII word starting with a
+      letter, so `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and
+      `test realm` are refused. The pattern is the realm row's own CHECK
+      constraint, so the database is what refuses an off-pattern name,
+      not a check in the handler: the rejection arrives as a 422 only
+      after the realm CA has been minted and sealed and is then thrown
+      away. Empty is a different failure with a different message — see
+      the 422 response.
   jwt_public_key_pem:
     type: string
     description: PEM-encoded public key used to verify realm JWTs.
   device_registration_limit:

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -474,16 +474,29 @@
         errors:
           detail: Realm already exists
 
 ValidationError:
-  description: Unprocessable Entity — request body failed validation.
+  description: |
+    Unprocessable Entity — request body failed validation. The two
+    `realm_name` rejections partition the cases and neither overlaps: an
+    empty name is refused by the handler with `realm_name can't be blank`,
+    a non-empty name off the `^[a-z][a-z0-9]*$` pattern is refused by the
+    realm row's CHECK constraint with `realm_name is invalid`.
   content:
     application/json; charset=utf-8:
       schema:
         $ref: "#/components/schemas/ErrorDetail"
-      example:
-        errors:
-          detail: "realm_name can't be blank"
+      examples:
+        blank_name:
+          summary: Empty realm_name — refused by the handler
+          value:
+            errors:
+              detail: "realm_name can't be blank"
+        off_pattern_name:
+          summary: Non-empty realm_name off the pattern — refused by the database
+          value:
+            errors:
+              detail: realm_name is invalid
 
 InternalServerError:
   description: Internal Server Error — unexpected failure.
   content:

Now the POST 422 story:
$ rg -n "go:embed" --glob '*.go' internal/store pkg | head; echo ---; rg -n "migrations" internal/store/*.go | head -10
---
internal/store/store.go:4:// embedded schema migrations (docs/DESIGN.md §2.2–2.5), CA-key sealing, and
internal/store/store.go:23:	"github.com/astrate-platform/astrate/migrations"
internal/store/store.go:67:// New connects to dsn, applies any pending embedded migrations, probes
internal/store/store.go:100:// runMigrations applies the embedded migrations up to head. It is a no-op on
internal/store/store.go:107:	src, err := iofs.New(migrations.FS, ".")
internal/store/store.go:110:		return fmt.Errorf("store: loading embedded migrations: %w", err)
internal/store/store.go:126:		return fmt.Errorf("store: applying migrations: %w", upErr)
internal/store/store_test.go:20:	"github.com/astrate-platform/astrate/migrations"
internal/store/store_test.go:77:	if err := s.pool.QueryRow(ctx, `SELECT version, dirty FROM schema_migrations`).Scan(&version, &dirty); err != nil {
internal/store/store_test.go:78:		t.Fatalf("reading schema_migrations: %v", err)

$ cat migrations/embed.go 2>/dev/null || ls migrations
000001_extensions.down.sql
000001_extensions.up.sql
000002_metadata.down.sql
000002_metadata.up.sql
000003_properties.down.sql
000003_properties.up.sql
000004_datastreams.down.sql
000004_datastreams.up.sql
000005_jobs.down.sql
000005_jobs.up.sql
000006_policies.down.sql
000006_policies.up.sql
000007_toolkit.down.sql
000007_toolkit.up.sql
000008_pipelines.down.sql
000008_pipelines.up.sql
000009_flows.down.sql
000009_flows.up.sql
000010_flow_failed_block.down.sql
000010_flow_failed_block.up.sql
000011_realm_retention.down.sql
000011_realm_retention.up.sql
000012_user_blocks.down.sql
000012_user_blocks.up.sql
migrations.go

→ Read migrations/migrations.go
Now the documentation-half test, in the house style of `TestRealmManagementDeviceIDEncodingDocumented`:

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -955,8 +955,106 @@
 		}
 	}
 }
 
+// TestHousekeepingRealmNamePatternDocumented guards that the housekeeping spec
+// constrains `realm_name` the way the schema does. The realms table CHECKs its
+// name column against a lowercase ASCII word (migrations/000002_metadata.up.sql),
+// so `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are
+// not realm names: a non-blank name that misses the pattern is rejected by the
+// INSERT as a check violation, folded into store.ErrInvalidRealmName
+// (internal/store/realms.go), re-reported by Service.CreateRealm as the
+// ErrValidation `realm_name is invalid` (internal/housekeeping/service.go) and
+// answered 422 with the flat detail body once validationDetail strips the prefix
+// (internal/housekeeping/http.go). Both `RealmCreate.realm_name` and the
+// `RealmName` path parameter used to be a bare `type: string`, which reads as
+// any string and tells a generated client nothing about the six spellings a
+// realm cannot have.
+//
+// The pattern is read out of the migration and compared with the one the spec
+// ships, so the documentation cannot drift from the constraint, and the
+// migration's own regex is held against the spellings the store suite refuses
+// (internal/store/realms_test.go). Two asymmetries are pinned because a client
+// reading only the pattern would get both wrong: the two 422 details partition
+// the cases rather than being two spellings of one failure, and the name in the
+// path is a read, where an off-pattern name is a 404 and not a 422 — the
+// read-path hole the DeviceID parameter documents on the realm-management twin.
+func TestHousekeepingRealmNamePatternDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	sql, err := migrations.FS.ReadFile("000002_metadata.up.sql")
+	if err != nil {
+		t.Fatalf("reading 000002_metadata.up.sql: %v", err)
+	}
+	check := regexp.MustCompile(`CHECK \(name ~ '([^']+)'\)`).FindSubmatch(sql)
+	if check == nil {
+		t.Fatal("realms.name declares no CHECK (name ~ ...) constraint; this test's premise no longer holds")
+	}
+	want := string(check[1])
+
+	re, err := regexp.Compile(want)
+	if err != nil {
+		t.Fatalf("realm name CHECK %q does not compile: %v", want, err)
+	}
+	for _, name := range []string{"test", "realm2", "a1"} {
+		if !re.MatchString(name) {
+			t.Errorf("realm name CHECK %q rejects the valid name %q", want, name)
+		}
+	}
+	for _, name := range []string{
+		"", "TEST", "test_realm", "test-realm", "1test", "tëst", "test realm",
+		"re-alm", "re_alm",
+	} {
+		if re.MatchString(name) {
+			t.Errorf("realm name CHECK %q accepts %q; this test's premise no longer holds", want, name)
+		}
+	}
+
+	param := componentBlock(t, lines, "    RealmName:")
+	if got := schemaPattern(t, param, "        "); got != want {
+		t.Errorf("RealmName parameter pattern = %q, want the realm name CHECK %q", got, want)
+	}
+	for _, wording := range []string{astarteapi.DetailNotFound, "404"} {
+		if !strings.Contains(strings.Join(param, "\n"), wording) {
+			t.Errorf("RealmName parameter does not record that an off-pattern name reads as %q", wording)
+		}
+	}
+
+	create := propertyBlock(t, componentBlock(t, lines, "    RealmCreate:"), "realm_name")
+	if got := schemaPattern(t, create, "          "); got != want {
+		t.Errorf("RealmCreate.realm_name pattern = %q, want the realm name CHECK %q", got, want)
+	}
+
+	if !containsLine(operationBlock(t, lines, "createRealm"),
+		`          $ref: "#/components/responses/ValidationError"`) {
+		t.Error("createRealm 422 does not $ref the ValidationError response")
+	}
+	unprocessable := strings.Join(componentBlock(t, lines, "    ValidationError:"), "\n")
+	for _, wording := range []string{
+		"partition", "realm_name can't be blank", "realm_name is invalid",
+	} {
+		if !strings.Contains(unprocessable, wording) {
+			t.Errorf("ValidationError response does not say %q", wording)
+		}
+	}
+}
+
+// schemaPattern returns the `pattern` value carried at the given indentation,
+// or "" when the block declares none.
+func schemaPattern(t *testing.T, lines []string, indent string) string {
+	t.Helper()
+	for _, l := range lines {
+		if v, ok := strings.CutPrefix(l, indent+"pattern: "); ok {
+			return strings.Trim(v, `'"`)
+		}
+	}
+	return ""
+}
+
 // TestPairingUnregisterDeviceSemanticsDocumented guards that the pairing spec
 // describes `unregisterDevice` for what it does. "Removes a device from the
 // realm" reads as a device-and-data deletion; store.UnregisterDevice clears the
 // credential trail and flips status back to 'registered' and touches nothing

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -15,8 +15,9 @@
 	"testing"
 
 	docs "github.com/astrate-platform/astrate/docs"
 	"github.com/astrate-platform/astrate/internal/engine/triggers"
+	"github.com/astrate-platform/astrate/migrations"
 	"github.com/astrate-platform/astrate/pkg/astarteapi"
 	"github.com/astrate-platform/astrate/pkg/deviceid"
 )
 

$ go test ./internal/swagger/ -run 'Housekeeping' -v 2>&1 | tail -20
=== RUN   TestHousekeepingAsyncOperationParamDocumented
--- PASS: TestHousekeepingAsyncOperationParamDocumented (0.00s)
=== RUN   TestHousekeepingRetentionZeroFoldDocumented
--- PASS: TestHousekeepingRetentionZeroFoldDocumented (0.00s)
=== RUN   TestHousekeepingErrorDetailExamples
--- PASS: TestHousekeepingErrorDetailExamples (0.00s)
=== RUN   TestHousekeepingRealmNamePatternDocumented
--- PASS: TestHousekeepingRealmNamePatternDocumented (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.016s

$ cp docs/api/astarte_housekeeping_api.yaml /tmp/opencode/hk.yaml.bak && python3 - <<'EOF'
import re
p='docs/api/astarte_housekeeping_api.yaml'
s=open(p).read()
# strip the two pattern lines
s=s.replace("          pattern: '^[a-z][a-z0-9]*$'\n","",1)
s=s.replace("        pattern: '^[a-z][a-z0-9]*$'\n","",1)
# strip the partition wording + second example
s=s.replace("""      description: |
        Unprocessable Entity — request body failed validation. The two
        `realm_name` rejections partition the cases and neither overlaps: an
        empty name is refused by the handler with `realm_name can't be blank`,
        a non-empty name off the `^[a-z][a-z0-9]*$` pattern is refused by the
        realm row's CHECK constraint with `realm_name is invalid`.
""","""      description: Unprocessable Entity — request body failed validation.
""",1)
s=s.replace("""          examples:
            blank_name:
              summary: Empty realm_name — refused by the handler
              value:
                errors:
                  detail: "realm_name can't be blank"
            off_pattern_name:
              summary: Non-empty realm_name off the pattern — refused by the database
              value:
                errors:
                  detail: realm_name is invalid
""","""          example:
            errors:
              detail: "realm_name can't be blank"
""",1)
open(p,'w').write(s)
EOF
go test ./internal/swagger/ -run TestHousekeepingRealmNamePatternDocumented 2>&1 | tail -20
--- FAIL: TestHousekeepingRealmNamePatternDocumented (0.00s)
    swagger_test.go:1019: RealmName parameter pattern = "", want the realm name CHECK "^[a-z][a-z0-9]*$"
    swagger_test.go:1029: RealmCreate.realm_name pattern = "", want the realm name CHECK "^[a-z][a-z0-9]*$"
    swagger_test.go:1041: ValidationError response does not say "partition"
    swagger_test.go:1041: ValidationError response does not say "realm_name is invalid"
FAIL
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.017s
FAIL

$ cp /tmp/opencode/hk.yaml.bak docs/api/astarte_housekeeping_api.yaml && go test ./internal/swagger/ 2>&1 | tail -5 && rg -n "pattern: '\^\[a-z\]" docs/api/astarte_housekeeping_api.yaml
ok  	github.com/astrate-platform/astrate/internal/swagger	0.050s
246:        pattern: '^[a-z][a-z0-9]*$'
309:          pattern: '^[a-z][a-z0-9]*$'

Now install mkdocs and run the docs build:
$ python3 -m venv /tmp/opencode/mkvenv >/dev/null 2>&1 && /tmp/opencode/mkvenv/bin/python -m pip install -q --disable-pip-version-check -r docs/requirements.txt 2>&1 | tail -3; /tmp/opencode/mkvenv/bin/mkdocs --version
mkdocs, version 1.6.1 from /tmp/opencode/mkvenv/lib/python3.13/site-packages/mkdocs (Python 3.13)

$ PATH=/tmp/opencode/mkvenv/bin:$PATH make -C docs build 2>&1 | tail -25
mkdocs build -f mkdocs.yml

 │  ⚠  Warning from the Material for MkDocs team
 │
 │  MkDocs 2.0, the underlying framework of Material for MkDocs,
 │  will introduce backward-incompatible changes, including:
 │
 │  × All plugins will stop working – the plugin system has been removed
 │  × All theme overrides will break – the theming system has been rewritten
 │  × No migration path exists – existing projects cannot be upgraded
 │  × Closed contribution model – community members can't report bugs
 │  × Currently unlicensed – unsuitable for production use
 │
 │  Our full analysis:
 │
 │  https://squidfunk.github.io/mkdocs-material/blog/2026/02/18/mkdocs-2.0/

INFO    -  Cleaning site directory
INFO    -  Building documentation to directory: /root/astrate-mule/site-dist
INFO    -  The following pages exist in the docs directory, but are not included in the "nav" configuration:
  - DESIGN.md
  - ROADMAP.md
INFO    -  Doc file 'index.md' contains an unrecognized relative link 'api/', it was left as is.
INFO    -  Documentation built in 3.93 seconds
make: Leaving directory '/root/astrate-mule/docs'

$ git status --short | head -30; echo "--- ignored?"; git check-ignore -v site-dist docs/site/api/astarte_housekeeping_api.yaml 2>&1 | head
M docs/api/astarte_housekeeping_api.yaml
 M internal/swagger/swagger_test.go
--- ignored?
.gitignore:29:/site-dist/	site-dist
.gitignore:34:/docs/site/api/	docs/site/api/astarte_housekeeping_api.yaml

$ /tmp/opencode/mkvenv/bin/python - <<'EOF'
import yaml, glob, json, os
for f in sorted(glob.glob('docs/api/*.yaml')):
    doc = yaml.safe_load(open(f))
    refs = []
    def walk(n):
        if isinstance(n, dict):
            for k, v in n.items():
                if k == '$ref' and isinstance(v, str):
                    refs.append(v)
                else:
                    walk(v)
        elif isinstance(n, list):
            for v in n: walk(v)
    walk(doc)
    missing = []
    for r in refs:
        if not r.startswith('#/'): continue
        cur = doc
        for part in r[2:].split('/'):
            part = part.replace('~1','/').replace('~0','~')
            if isinstance(cur, dict) and part in cur:
                cur = cur[part]
            else:
                missing.append(r); break
    print(f"{os.path.basename(f)}: parses OK, openapi={doc.get('openapi')}, refs={len(refs)}, unresolved={missing}")

hk = yaml.safe_load(open('docs/api/astarte_housekeeping_api.yaml'))
print("RealmName param:", json.dumps(hk['components']['parameters']['RealmName'], ensure_ascii=False, indent=1))
print("RealmCreate.realm_name:", json.dumps(hk['components']['schemas']['RealmCreate']['properties']['realm_name'], ensure_ascii=False, indent=1))
print("POST 422 response:", json.dumps(hk['components']['responses']['ValidationError'], ensure_ascii=False, indent=1))
EOF
echo "--- site-dist copy identical?"; cmp docs/api/astarte_housekeeping_api.yaml site-dist/api/astarte_housekeeping_api.yaml && echo identical
astarte_appengine_api.yaml: parses OK, openapi=3.0.3, refs=346, unresolved=[]
astarte_housekeeping_api.yaml: parses OK, openapi=3.0.3, refs=46, unresolved=[]
astarte_pairing_api.yaml: parses OK, openapi=3.0.3, refs=59, unresolved=[]
astarte_realm_management_api.yaml: parses OK, openapi=3.0.3, refs=154, unresolved=[]
astrate_native_api.yaml: parses OK, openapi=3.0.3, refs=22, unresolved=[]
RealmName param: {
 "name": "realm",
 "in": "path",
 "required": true,
 "description": "The realm name: a lowercase ASCII word starting with a letter, so\n`TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm`\nare not realm names — no realm can be created under one. This parameter\nis the read path though, and there an off-pattern name is not reported\nas malformed: it is answered with the same `404` and\n`{\"errors\":{\"detail\":\"Not Found\"}}` body as a realm that does not\nexist, because the lookup matches the name verbatim with no format check\nof its own.\n",
 "schema": {
  "type": "string",
  "pattern": "^[a-z][a-z0-9]*$"
 }
}
RealmCreate.realm_name: {
 "type": "string",
 "pattern": "^[a-z][a-z0-9]*$",
 "description": "Unique realm identifier: a lowercase ASCII word starting with a\nletter, so `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and\n`test realm` are refused. The pattern is the realm row's own CHECK\nconstraint, so the database is what refuses an off-pattern name,\nnot a check in the handler: the rejection arrives as a 422 only\nafter the realm CA has been minted and sealed and is then thrown\naway. Empty is a different failure with a different message — see\nthe 422 response.\n"
}
POST 422 response: {
 "description": "Unprocessable Entity — request body failed validation. The two\n`realm_name` rejections partition the cases and neither overlaps: an\nempty name is refused by the handler with `realm_name can't be blank`,\na non-empty name off the `^[a-z][a-z0-9]*$` pattern is refused by the\nrealm row's CHECK constraint with `realm_name is invalid`.\n",
 "content": {
  "application/json; charset=utf-8": {
   "schema": {
    "$ref": "#/components/schemas/ErrorDetail"
   },
   "examples": {
    "blank_name": {
     "summary": "Empty realm_name — refused by the handler",
     "value": {
      "errors": {
       "detail": "realm_name can't be blank"
      }
     }
    },
    "off_pattern_name": {
     "summary": "Non-empty realm_name off the pattern — refused by the database",
     "value": {
      "errors": {
       "detail": "realm_name is invalid"
      }
     }
    }
   }
  }
 }
}
--- site-dist copy identical?
identical

Now the full gate:
$ go vet ./... 2>&1 | tail -5; echo "vet exit=$?"; gofmt -l . | head; echo "gofmt done"
vet exit=0
gofmt done

$ go test ./... 2>&1 | grep -v "^ok\|no test files" | head -20; echo "test exit=${PIPESTATUS[0]}"
/bin/sh: 1: Bad substitution

$ go test ./... > /tmp/opencode/gotest.log 2>&1; echo "exit=$?"; grep -v "^ok\|no test files" /tmp/opencode/gotest.log | head -20
exit=0

$ go test -race ./internal/swagger/ > /tmp/opencode/race.log 2>&1; echo "exit=$?"; tail -3 /tmp/opencode/race.log
exit=1
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.021s
FAIL

$ git diff --stat; git diff docs/api/astarte_housekeeping_api.yaml
docs/api/astarte_housekeeping_api.yaml | 43 ++++++++++++---
 internal/swagger/swagger_test.go       | 99 ++++++++++++++++++++++++++++++++++
 2 files changed, 136 insertions(+), 6 deletions(-)
diff --git a/docs/api/astarte_housekeeping_api.yaml b/docs/api/astarte_housekeeping_api.yaml
index ac38b17..36b0d2a 100644
--- a/docs/api/astarte_housekeeping_api.yaml
+++ b/docs/api/astarte_housekeeping_api.yaml
@@ -232,9 +232,18 @@ components:
       name: realm
       in: path
       required: true
-      description: The realm name.
+      description: |
+        The realm name: a lowercase ASCII word starting with a letter, so
+        `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm`
+        are not realm names — no realm can be created under one. This parameter
+        is the read path though, and there an off-pattern name is not reported
+        as malformed: it is answered with the same `404` and
+        `{"errors":{"detail":"Not Found"}}` body as a realm that does not
+        exist, because the lookup matches the name verbatim with no format check
+        of its own.
       schema:
         type: string
+        pattern: '^[a-z][a-z0-9]*$'
 
     AsyncOperation:
       name: async_operation
@@ -297,7 +306,16 @@ components:
       properties:
         realm_name:
           type: string
-          description: Unique realm identifier.
+          pattern: '^[a-z][a-z0-9]*$'
+          description: |
+            Unique realm identifier: a lowercase ASCII word starting with a
+            letter, so `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and
+            `test realm` are refused. The pattern is the realm row's own CHECK
+            constraint, so the database is what refuses an off-pattern name,
+            not a check in the handler: the rejection arrives as a 422 only
+            after the realm CA has been minted and sealed and is then thrown
+            away. Empty is a different failure with a different message — see
+            the 422 response.
         jwt_public_key_pem:
           type: string
           description: PEM-encoded public key used to verify realm JWTs.
@@ -457,14 +475,27 @@ components:
               detail: Realm already exists
 
     ValidationError:
-      description: Unprocessable Entity — request body failed validation.
+      description: |
+        Unprocessable Entity — request body failed validation. The two
+        `realm_name` rejections partition the cases and neither overlaps: an
+        empty name is refused by the handler with `realm_name can't be blank`,
+        a non-empty name off the `^[a-z][a-z0-9]*$` pattern is refused by the
+        realm row's CHECK constraint with `realm_name is invalid`.
       content:
         application/json; charset=utf-8:
           schema:
             $ref: "#/components/schemas/ErrorDetail"
-          example:
-            errors:
-              detail: "realm_name can't be blank"
+          examples:
+            blank_name:
+              summary: Empty realm_name — refused by the handler
+              value:
+                errors:
+                  detail: "realm_name can't be blank"
+            off_pattern_name:
+              summary: Non-empty realm_name off the pattern — refused by the database
+              value:
+                errors:
+                  detail: realm_name is invalid
 
     InternalServerError:
       description: Internal Server Error — unexpected failure.

$ tail -8 .mule/todo.md
- [!] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implements `device_empty_cache_received` end to end — `SimpleEvent` field 23, an empty `DeviceEmptyCacheReceivedEvent` whose encoder emits exactly `{"type":"device_empty_cache_received"}` (astarte_core encoder.ex:54-62), fired by DUP at `control_handler.ex:161` *inside* the `with` that runs `send_control_consumer_properties` → `resend_all_properties` → `set_pending_empty_cache`, so it fires only when all three succeed and never on failure, and the commit deletes the `# TODO: implement on_empty_cache_received` from `astarte_events/lib/astarte_events/triggers/core.ex`. Astrate already accepts the trigger but never evaluates it: `OnDeviceEmptyCacheReceived` sits in `deviceOns` with value `false` (internal/engine/triggers/match.go:123, the `evaluated` flag read at match.go:406) and the comment's stated reason ("upstream defines no SimpleEvent variant for it") is now false. Add the event body constructor to internal/engine/triggers/events.go (the payload is `{"type":"device_empty_cache_received"}` and nothing else — no device id, no metadata), flip the `deviceOns` value to true, and fire it from `handleEmptyCache` (internal/engine/control.go:58) with the same idiom `handleLifecycle` uses for `device_connected` (`e.fireDevice` + `stream.Publish`, internal/engine/engine.go:346-353), at the point where `m.Ack()` runs — i.e. after `resendServerProperties` and `sendConsumerProperties` have both succeeded inside `retryStore`. Add a T1 test that a trigger on it fires exactly once per accepted emptyCache and does not fire when the resend fails. — BLOCKED: wrote nothing
- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
- [ ] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362) — the second half docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail: Bad Request`, the reconstructed title-case form that `TestHousekeepingErrorDetailExamples` (internal/swagger/swagger_test.go:687) now pins against everywhere else on this spec, while the wire emits the frozen `astarteapi.DetailBadRequest` = `"Bad request"` (pkg/astarteapi/envelope.go:37, via `WriteBadRequest` on every 400 — internal/housekeeping/http.go:65, 166, 171, 183). Because the example sits on the generic `ErrorDetail` schema, which all eight error responses `$ref` (yaml:377-477), it is the example a generator falls back to for any response whose own example it cannot resolve — so it teaches `"Bad Request"` for a 404 and a 500 too. Change it to `Bad request`, or drop it as the realm-management and pairing `ErrorDetail` schemas do (neither carries one). Then extend `TestHousekeepingErrorDetailExamples`, whose comment currently says the schema-level example is "deliberately not pinned" (swagger_test.go:684-686), to pin it like the responses. One thing not to "fix": the example is nested one level in — it is the `errors` sub-schema's example, not a bare envelope — so do not re-wrap it in another `errors:`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-hk-patch-422-field-error-examples [auto]: add the three missing PATCH 422 bodies to the `ValidationErrors` response component in docs/api/astarte_housekeeping_api.yaml (yaml:427-436, referenced once at yaml:190) — it carries a single example, `error_name: [invalid_update_parameters]` (internal/housekeeping/http.go:176-177), but `patchRealm` emits three more field-error envelopes: `jwt_public_key_pem: ["can't be blank"]` (http.go:188-189), `device_registration_limit: ["is invalid"]` (http.go:193-194) and `datastream_maximum_storage_retention: ["is invalid"]` (http.go:198-199). Add them as further `examples` on the media type (OpenAPI 3.0.3, yaml:1) or split the component per case. Also record that the *messages* differ from POST for the same conditions: create rejects a negative limit/retention with the flat `ErrValidation` detail `device_registration_limit must be non-negative` / `datastream_maximum_storage_retention must be non-negative` (internal/housekeeping/service.go:143-148, 422 `ErrorDetail` shape), while PATCH answers `is invalid` in the FieldErrors shape — and that the service's own three `ErrValidation` branches (service.go:205-213) are unreachable from REST, because `patchRealm` pre-checks the identical conditions and answers first. Say in the `RealmPatch` field descriptions that a negative value is a 422 rather than silently ignored. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-hk-wrong-type-field-400 [auto]: widen the `BadRequest` component description in docs/api/astarte_housekeeping_api.yaml (yaml:378, "malformed or missing data envelope") — measured today, a perfectly well-formed `data` envelope carrying a wrong-typed field answers **400**, not the documented 422. `astarteapi.DecodeData` into `map[string]json.RawMessage` (internal/housekeeping/http.go:170) cannot fail on a type mismatch, so the failure surfaces on the second decode into the typed body (http.go:182-184) and is answered `WriteBadRequest`. Probe output: `{"data":{"device_registration_limit":"100"}}` → `json: cannot unmarshal string into Go struct field patchBody.device_registration_limit of type int32`; `{"data":{"jwt_public_key_pem":123}}` → the string twin; `{"data":{"device_registration_limit":1.5}}` → `cannot unmarshal number 1.5`; `{"data":[]}` → `cannot unmarshal array into Go value of type patchBody`. Same on POST, which decodes straight into `realmBody` (http.go:63-67). This is unreachable from a spec-conformant client — `RealmCreate`/`RealmPatch` already type both fields — so document it as a one-line extension of the 400 description on both operations rather than inventing a 422, and do NOT re-shape the `BadRequest` example, whose `detail: Bad request` is the frozen constant `TestHousekeepingErrorDetailExamples` pins. One distinction worth stating in the same description: a wrong-typed field is *not* the `invalid_update_parameters` case — unknown keys are caught earlier by the `patchAllowedFields` sweep (http.go:174-180, 422 FieldErrors) on PATCH only, and silently dropped on POST. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ tail -c 200 .mule/todo.md | od -c | tail -3
0000260   e   n   c   e   d       Y   A   M   L       s   t   i   l   l
0000300       l   o   a   d   s   .  \n
0000310

$ cat >> .mule/todo.md <<'EOF'
- [ ] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `components.parameters.RealmName` and `RealmCreate.realm_name` both carry `pattern: '^[a-z][a-z0-9]*$'` and `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go) holds both against the realm name CHECK extracted from migrations/000002_metadata.up.sql — but the read side is still bare `type: string`: `Realm.realm_name` (components.schemas.Realm) and `DataEnvelopeRealmNames.items`. Decide whether a response field earns the same pattern (a generated client learns nothing new from it, but a reader diffing create against get sees one constrained and one not) and, if it does, add it to both and extend that test to hold the response-side pattern against the same CHECK. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
EOF
tail -1 .mule/todo.md | head -c 120
- [ ] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/a

$ go vet ./internal/swagger/ && go test ./internal/swagger/ && gofmt -l . ; git status --short
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
 M .mule/todo.md
 M docs/api/astarte_housekeeping_api.yaml
 M internal/swagger/swagger_test.go

Done: constrained `realm_name` with the schema's own `^[a-z][a-z0-9]*$` in the housekeeping spec (both `RealmCreate.realm_name` and the `RealmName` path parameter), recorded the POST 422 partition (blank → `can't be blank` from the handler, non-empty off-pattern → `is invalid` from the CHECK) as two named examples, recorded the read-path 404 asymmetry on `RealmName`, and added `TestHousekeepingRealmNamePatternDocumented`, which derives the pattern from `migrations/000002_metadata.up.sql` so the spec cannot drift from the constraint.

Files: `docs/api/astarte_housekeeping_api.yaml`, `internal/swagger/swagger_test.go`, `.mule/todo.md`

Verified:
- `go vet ./...` → pass; `go test ./...` → pass; `gofmt -l .` → clean. Ran **without** `-race`: the race build is unusable here (`FATAL: Found 39 - Supported 48`), and this change touches no concurrency.
- Test-fails-without-the-change check: stripped the two `pattern:` lines and the partition wording from a copy of the YAML → `TestHousekeepingRealmNamePatternDocumented` FAILs on 4 assertions; restored → passes.
- `make -C docs build` → `Documentation built in 3.93 seconds` (mkdocs 1.6.1 + material installed into a `/tmp/opencode` venv; neither was present on the machine). Then parsed all five `docs/api/*.yaml` with PyYAML: all parse, all 46 local `$ref`s in the housekeeping spec resolve, and the `site-dist/api/astarte_housekeeping_api.yaml` copy is byte-identical to the source.
- Code chain re-read before editing: CHECK at `migrations/000002_metadata.up.sql:7`; `pgCodeCheckViolation` → `ErrInvalidRealmName` (`internal/store/realms.go:68-69`); → `ErrValidation: realm_name is invalid` (`internal/housekeeping/service.go:179-181`); →422 flat detail via `validationDetail` (`internal/housekeeping/http.go:240-241,251-259`); blank branch at `service.go:137-139`; `ca.Generate(name, 0)` at `service.go:161` with the name only interpolated into the CN (`internal/pairing/ca/ca.go:88`), so x509 validates nothing; `GetRealmByName` is a bare `WHERE name = $1` → `ErrNotFound` → `WriteNotFound` (`DetailNotFound = "Not Found"`). Pattern behaviour checked with a throwaway Go probe in `/tmp`: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst`, `test realm`, `""` all refused; `test`, `realm2`, `a1` accepted.

Unsure: two things. (1) The DB-side refusal of those six spellings rests on the task's measurement plus the integration-only `TestInvalidName` (`internal/store/realms_test.go:112`), which I cannot run here (no Docker, no DB) — I verified the regex semantics in Go, not in Postgres. (2) I read "the POST 422 story" as the `ValidationError` component (referenced only by `createRealm`) and converted its single `example:` to two named `examples:`, matching the house style used by the realm-management 422s; if you meant the `createRealm` description instead, that text should be moved.

Follow-ups: one line appended to `.mule/todo.md` — `docs-sync-hk-realm-name-response-schemas`: the read side (`Realm.realm_name`, `DataEnvelopeRealmNames.items`) is still a bare `type: string`, since the task named only the two write-side spots.
