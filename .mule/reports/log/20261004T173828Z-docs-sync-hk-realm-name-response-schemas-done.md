slug: docs-sync-hk-realm-name-response-schemas
verdict: done
at:  def0af6
ran: 2026-10-04T17:38:28Z on DietPi in 178s


> build · big-pickle

I'll start by reading the two files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ rg -n "RealmName|realm_name|Realm\b" docs/api/astarte_housekeeping_api.yaml | head -60
40:                $ref: "#/components/schemas/DataEnvelopeRealmNames"
51:      operationId: createRealm
74:                realm_name: test
82:          description: Realm created successfully.
86:                $ref: "#/components/schemas/DataEnvelopeRealm"
89:                  realm_name: test
110:      operationId: getRealm
117:        - $ref: "#/components/parameters/RealmName"
120:          description: Realm details.
124:                $ref: "#/components/schemas/DataEnvelopeRealm"
127:                  realm_name: test
143:      operationId: patchRealm
159:        - $ref: "#/components/parameters/RealmName"
180:          description: Realm updated successfully. No body returned.
195:      operationId: deleteRealm
202:        - $ref: "#/components/parameters/RealmName"
206:          description: Realm deleted successfully. No body returned.
231:    RealmName:
265:    DataEnvelopeRealm:
270:          $ref: "#/components/schemas/Realm"
272:    DataEnvelopeRealmNames:
282:    Realm:
284:      required: [realm_name, jwt_public_key_pem]
286:        realm_name:
305:      required: [realm_name, jwt_public_key_pem]
307:        realm_name:
443:              detail: Realm deletion disabled
475:              detail: Realm already exists
480:        `realm_name` rejections partition the cases and neither overlaps: an
481:        empty name is refused by the handler with `realm_name can't be blank`,
483:        realm row's CHECK constraint with `realm_name is invalid`.
490:              summary: Empty realm_name — refused by the handler
493:                  detail: "realm_name can't be blank"
495:              summary: Non-empty realm_name off the pattern — refused by the database
498:                  detail: realm_name is invalid

$ rg -n "realm" migrations/000002_metadata.up.sql | head -30
5:CREATE TABLE realms (
9:    ca_certificate   text NOT NULL,                  -- realm CA cert, PEM
20:    realm_id      smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
29:    UNIQUE (realm_id, name, major_version)
52:    realm_id            smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
71:    PRIMARY KEY (realm_id, id)
78:    realm_id  smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
80:    UNIQUE (realm_id, name)
84:    realm_id  smallint NOT NULL,
87:    FOREIGN KEY (realm_id, device_id) REFERENCES devices(realm_id, id) ON DELETE CASCADE
93:    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
96:    UNIQUE (realm_id, name)

$ sed -n 1,20p migrations/000002_metadata.up.sql
-- 000002: relational metadata schema, transcribed verbatim from
-- docs/DESIGN.md §2.2 (docs/ROADMAP.md §3.1 file 2.2).

-- Realms (Housekeeping domain)
CREATE TABLE realms (
    id               smallint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name             text NOT NULL UNIQUE CHECK (name ~ '^[a-z][a-z0-9]*$'),
    jwt_public_keys  jsonb NOT NULL DEFAULT '[]',   -- array of PEM strings (RSA/EC)
    ca_certificate   text NOT NULL,                  -- realm CA cert, PEM
    ca_private_key   bytea NOT NULL,                 -- encrypted at rest (AES-256-GCM,
                                                     -- key from config/KMS env var)
    device_registration_limit integer,
    created_at       timestamptz NOT NULL DEFAULT now()
);

-- Interfaces (Realm Management domain). The raw JSON is the source of truth;
-- generated columns lift the routing-critical fields out for indexing.
CREATE TABLE interfaces (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    realm_id      smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,

$ sed -n 225,320p docs/api/astarte_housekeeping_api.yaml
name: Authorization
      description: |
        Instance-level JWT carrying the `a_ha` claim. Used to authenticate
        housekeeping operations (realm CRUD).

  parameters:
    RealmName:
      name: realm
      in: path
      required: true
      description: |
        The realm name: a lowercase ASCII word starting with a letter, so
        `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm`
        are not realm names — no realm can be created under one. This parameter
        is the read path though, and there an off-pattern name is not reported
        as malformed: it is answered with the same `404` and
        `{"errors":{"detail":"Not Found"}}` body as a realm that does not
        exist, because the lookup matches the name verbatim with no format check
        of its own.
      schema:
        type: string
        pattern: '^[a-z][a-z0-9]*$'

    AsyncOperation:
      name: async_operation
      in: query
      required: false
      description: |
        Accepted and ignored on either value. Upstream 1.4 runs this operation in
        the background and lets the caller opt into synchronous execution with
        `async_operation=false`; Astrate always works synchronously and answers
        only once the work is done, so the value changes nothing and the
        response is identical either way (deviation 17 in COMPATIBILITY.md).
      schema:
        type: boolean
        default: false

  schemas:
    # ── Envelope variants ──

    DataEnvelopeRealm:
      type: object
      required: [data]
      properties:
        data:
          $ref: "#/components/schemas/Realm"

    DataEnvelopeRealmNames:
      type: object
      required: [data]
      properties:
        data:
          type: array
          items:
            type: string
          description: List of realm names.

    Realm:
      type: object
      required: [realm_name, jwt_public_key_pem]
      properties:
        realm_name:
          type: string
          description: Unique realm identifier.
        jwt_public_key_pem:
          type: string
          description: PEM-encoded public key used to verify realm JWTs.
        device_registration_limit:
          type: integer
          format: int32
          nullable: true
          description: Maximum number of devices allowed in this realm. Null means unlimited.
        datastream_maximum_storage_retention:
          type: integer
          format: int64
          nullable: true
          description: Maximum datastream storage retention in seconds. Null means unlimited.

    RealmCreate:
      type: object
      required: [realm_name, jwt_public_key_pem]
      properties:
        realm_name:
          type: string
          pattern: '^[a-z][a-z0-9]*$'
          description: |
            Unique realm identifier: a lowercase ASCII word starting with a
            letter, so `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and
            `test realm` are refused. The pattern is the realm row's own CHECK
            constraint, so the database is what refuses an off-pattern name,
            not a check in the handler: the rejection arrives as a 422 only
            after the realm CA has been minted and sealed and is then thrown
            away. Empty is a different failure with a different message — see
            the 422 response.
        jwt_public_key_pem:
          type: string

$ rg -rn "TestHousekeepingRealmNamePatternDocumented" internal/swagger
internal/swagger/swagger_test.go:// n guards that the housekeeping spec
internal/swagger/swagger_test.go:func n(t *testing.T) {

$ sed -n 1,80p internal/swagger/swagger_test.go
package swagger

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	docs "github.com/astrate-platform/astrate/docs"
	"github.com/astrate-platform/astrate/internal/engine/triggers"
	"github.com/astrate-platform/astrate/migrations"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
	"github.com/astrate-platform/astrate/pkg/deviceid"
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

$ grep -rn "TestHousekeepingRealmNamePatternDocumented" internal/swagger/swagger_test.go | head -3
982:// TestHousekeepingRealmNamePatternDocumented guards that the housekeeping spec
1004:func TestHousekeepingRealmNamePatternDocumented(t *testing.T) {

$ sed -n 960,1040p internal/swagger/swagger_test.go
re, err := regexp.Compile(pattern)
	if err != nil {
		t.Fatalf("DeviceID pattern %q does not compile: %v", pattern, err)
	}

	if !re.MatchString("dT6hS2W9TT6LEnP25ks_lg") {
		t.Errorf("DeviceID pattern %q rejects the canonical ID dT6hS2W9TT6LEnP25ks_lg", pattern)
	}
	for _, spelling := range []string{
		"dT6hS2W9TT6LEnP25ks+lg",  // standard-alphabet '+'
		"dT6hS2W9TT6LEnP25ks/lg",  // standard-alphabet '/'
		"dT6hS2W9TT6LEnP25ks_lg=", // padded
	} {
		if _, err := deviceid.Parse(spelling); err == nil {
			t.Errorf("deviceid.Parse accepts %q; this test's premise no longer holds", spelling)
		}
		if re.MatchString(spelling) {
			t.Errorf("DeviceID pattern %q accepts %q, which deviceid.Parse rejects", pattern, spelling)
		}
	}
}

// TestHousekeepingRealmNamePatternDocumented guards that the housekeeping spec
// constrains `realm_name` the way the schema does. The realms table CHECKs its
// name column against a lowercase ASCII word (migrations/000002_metadata.up.sql),
// so `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are
// not realm names: a non-blank name that misses the pattern is rejected by the
// INSERT as a check violation, folded into store.ErrInvalidRealmName
// (internal/store/realms.go), re-reported by Service.CreateRealm as the
// ErrValidation `realm_name is invalid` (internal/housekeeping/service.go) and
// answered 422 with the flat detail body once validationDetail strips the prefix
// (internal/housekeeping/http.go). Both `RealmCreate.realm_name` and the
// `RealmName` path parameter used to be a bare `type: string`, which reads as
// any string and tells a generated client nothing about the six spellings a
// realm cannot have.
//
// The pattern is read out of the migration and compared with the one the spec
// ships, so the documentation cannot drift from the constraint, and the
// migration's own regex is held against the spellings the store suite refuses
// (internal/store/realms_test.go). Two asymmetries are pinned because a client
// reading only the pattern would get both wrong: the two 422 details partition
// the cases rather than being two spellings of one failure, and the name in the
// path is a read, where an off-pattern name is a 404 and not a 422 — the
// read-path hole the DeviceID parameter documents on the realm-management twin.
func TestHousekeepingRealmNamePatternDocumented(t *testing.T) {
	b, err := docs.APIYAML.ReadFile("api/astarte_housekeeping_api.yaml")
	if err != nil {
		t.Fatalf("reading astarte_housekeeping_api.yaml: %v", err)
	}
	lines := strings.Split(string(b), "\n")

	sql, err := migrations.FS.ReadFile("000002_metadata.up.sql")
	if err != nil {
		t.Fatalf("reading 000002_metadata.up.sql: %v", err)
	}
	check := regexp.MustCompile(`CHECK \(name ~ '([^']+)'\)`).FindSubmatch(sql)
	if check == nil {
		t.Fatal("realms.name declares no CHECK (name ~ ...) constraint; this test's premise no longer holds")
	}
	want := string(check[1])

	re, err := regexp.Compile(want)
	if err != nil {
		t.Fatalf("realm name CHECK %q does not compile: %v", want, err)
	}
	for _, name := range []string{"test", "realm2", "a1"} {
		if !re.MatchString(name) {
			t.Errorf("realm name CHECK %q rejects the valid name %q", want, name)
		}
	}
	for _, name := range []string{
		"", "TEST", "test_realm", "test-realm", "1test", "tëst", "test realm",
		"re-alm", "re_alm",
	} {
		if re.MatchString(name) {
			t.Errorf("realm name CHECK %q accepts %q; this test's premise no longer holds", want, name)
		}
	}

	param := componentBlock(t, lines, "    RealmName:")
	if got := schemaPattern(t, param, "        "); got != want {

$ sed -n 1040,1080p internal/swagger/swagger_test.go
if got := schemaPattern(t, param, "        "); got != want {
		t.Errorf("RealmName parameter pattern = %q, want the realm name CHECK %q", got, want)
	}
	for _, wording := range []string{astarteapi.DetailNotFound, "404"} {
		if !strings.Contains(strings.Join(param, "\n"), wording) {
			t.Errorf("RealmName parameter does not record that an off-pattern name reads as %q", wording)
		}
	}

	create := propertyBlock(t, componentBlock(t, lines, "    RealmCreate:"), "realm_name")
	if got := schemaPattern(t, create, "          "); got != want {
		t.Errorf("RealmCreate.realm_name pattern = %q, want the realm name CHECK %q", got, want)
	}

	if !containsLine(operationBlock(t, lines, "createRealm"),
		`          $ref: "#/components/responses/ValidationError"`) {
		t.Error("createRealm 422 does not $ref the ValidationError response")
	}
	unprocessable := strings.Join(componentBlock(t, lines, "    ValidationError:"), "\n")
	for _, wording := range []string{
		"partition", "realm_name can't be blank", "realm_name is invalid",
	} {
		if !strings.Contains(unprocessable, wording) {
			t.Errorf("ValidationError response does not say %q", wording)
		}
	}
}

// schemaPattern returns the `pattern` value carried at the given indentation,
// or "" when the block declares none.
func schemaPattern(t *testing.T, lines []string, indent string) string {
	t.Helper()
	for _, l := range lines {
		if v, ok := strings.CutPrefix(l, indent+"pattern: "); ok {
			return strings.Trim(v, `'"`)
		}
	}
	return ""
}

// TestPairingUnregisterDeviceSemanticsDocumented guards that the pairing spec

$ sed -n 265,300p docs/api/astarte_housekeeping_api.yaml
DataEnvelopeRealm:
      type: object
      required: [data]
      properties:
        data:
          $ref: "#/components/schemas/Realm"

    DataEnvelopeRealmNames:
      type: object
      required: [data]
      properties:
        data:
          type: array
          items:
            type: string
          description: List of realm names.

    Realm:
      type: object
      required: [realm_name, jwt_public_key_pem]
      properties:
        realm_name:
          type: string
          description: Unique realm identifier.
        jwt_public_key_pem:
          type: string
          description: PEM-encoded public key used to verify realm JWTs.
        device_registration_limit:
          type: integer
          format: int32
          nullable: true
          description: Maximum number of devices allowed in this realm. Null means unlimited.
        datastream_maximum_storage_retention:
          type: integer
          format: int64
          nullable: true

$ sed -n 280,290p docs/api/astarte_housekeeping_api.yaml
description: List of realm names.

    Realm:
      type: object
      required: [realm_name, jwt_public_key_pem]
      properties:
        realm_name:
          type: string
          description: Unique realm identifier.
        jwt_public_key_pem:
          type: string

Reading the files, I see the issue. The write side (RealmCreate) and path param (RealmName) both have pattern `^[a-z][a-z0-9]*$`. The read side (Realm.realm_name and DataEnvelopeRealmNames.items) currently just have `type: string` without the pattern. The test wants to extend coverage to the response-side patterns as well.

I need to:
1. Add pattern to `Realm.realm_name` schema
2. Add pattern to `DataEnvelopeRealmNames.items` (which is an array of realm names)
3. Extend the test to check the response-side patterns against the same CHECK

Let me make these changes.
← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -276,8 +276,9 @@
     data:
       type: array
       items:
         type: string
+        pattern: '^[a-z][a-z0-9]*$'
       description: List of realm names.
 
 Realm:
   type: object

← Edit docs/api/astarte_housekeeping_api.yaml
Index: /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
+++ /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
@@ -285,8 +285,9 @@
 required: [realm_name, jwt_public_key_pem]
 properties:
   realm_name:
     type: string
+    pattern: '^[a-z][a-z0-9]*$'
     description: Unique realm identifier.
   jwt_public_key_pem:
     type: string
     description: PEM-encoded public key used to verify realm JWTs.

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -1050,8 +1050,31 @@
 if got := schemaPattern(t, create, "          "); got != want {
 	t.Errorf("RealmCreate.realm_name pattern = %q, want the realm name CHECK %q", got, want)
 }
 
+realm := propertyBlock(t, componentBlock(t, lines, "    Realm:"), "realm_name")
+if got := schemaPattern(t, realm, "          "); got != want {
+	t.Errorf("Realm.realm_name pattern = %q, want the realm name CHECK %q", got, want)
+}
+
+// Also check response array items pattern for realm names list
+dataEnvelopeNames := componentBlock(t, lines, "    DataEnvelopeRealmNames:")
+dataBlock := propertyBlock(t, dataEnvelopeNames, "data")
+// Look for items pattern in dataBlock - items are indented under data
+items := []string{}
+for i, l := range dataBlock {
+	if strings.Contains(l, "items:") {
+		// collect from i onwards in dataBlock, until we hit something at same or lesser indent relative to data properties
+		for j := i; j < len(dataBlock); j++ {
+			items = append(items, dataBlock[j])
+		}
+		break
+	}
+}
+if got := schemaPattern(t, items, "            "); got != want {
+	t.Errorf("DataEnvelopeRealmNames.items pattern = %q, want the realm name CHECK %q", got, want)
+}
+
 if !containsLine(operationBlock(t, lines, "createRealm"),
 	`          $ref: "#/components/responses/ValidationError"`) {
 	t.Error("createRealm 422 does not $ref the ValidationError response")
 }

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

$ mkdocs --version 2>&1 | head -2
/bin/sh: 1: mkdocs: not found

$ python3 -m mkdocs --version 2>&1 | head -1
/usr/bin/python3: No module named mkdocs

$ pip3 show mkdocs 2>&1 | head -1
/bin/sh: 1: pip3: not found

$ cd docs && npx --yes @redocly/cli@latest lint api/astarte_housekeeping_api.yaml 2>&1 | tail -25
No configurations were provided -- using built in recommended configuration by default.

validating api/astarte_housekeeping_api.yaml...
api/astarte_housekeeping_api.yaml: validated in 424ms

Woohoo! Your API description is valid. 🎉

npm notice
npm notice New major version of npm available! 10.9.8 -> 12.2.0
npm notice Changelog: https://github.com/npm/cli/releases/tag/v12.2.0
npm notice To update run: npm install -g npm@12.2.0
npm notice

$ go test -v ./internal/swagger/... 2>&1 | tail -20
=== RUN   TestRealmManagementErrorDetailExamples
--- PASS: TestRealmManagementErrorDetailExamples (0.00s)
=== RUN   TestAppEngineErrorDetailExamples
--- PASS: TestAppEngineErrorDetailExamples (0.00s)
=== RUN   TestNativeErrorDetailExamples
--- PASS: TestNativeErrorDetailExamples (0.00s)
=== RUN   TestHousekeepingErrorDetailExamples
--- PASS: TestHousekeepingErrorDetailExamples (0.00s)
=== RUN   TestRealmManagement422ValidationDetails
--- PASS: TestRealmManagement422ValidationDetails (0.01s)
=== RUN   TestPairingDeviceIDEncodingDocumented
--- PASS: TestPairingDeviceIDEncodingDocumented (0.00s)
=== RUN   TestRealmManagementDeviceIDEncodingDocumented
--- PASS: TestRealmManagementDeviceIDEncodingDocumented (0.00s)
=== RUN   TestHousekeepingRealmNamePatternDocumented
--- PASS: TestHousekeepingRealmNamePatternDocumented (0.00s)
=== RUN   TestPairingUnregisterDeviceSemanticsDocumented
--- PASS: TestPairingUnregisterDeviceSemanticsDocumented (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/swagger	0.127s

Done: docs-sync-hk-realm-name-response-schemas - added realm_name pattern to Realm schema and DataEnvelopeRealmNames.items, extended test coverage.
Files: docs/api/astarte_housekeeping_api.yaml, internal/swagger/swagger_test.go
Verified: `go test ./internal/swagger/...` -> PASS; OpenAPI lint via npx @redocly/cli@latest on the spec -> valid
Unsure: nothing
Follow-ups: none
