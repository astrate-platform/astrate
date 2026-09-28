slug: docs-sync-pairing-initial-payload-format-enum
verdict: done
at:  1f6a6c2
ran: 2026-09-28T17:26:27Z on DietPi in 95s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/task.md

$ cd /root/astrate-mule && sed -n '350,400p' docs/api/astarte_pairing_api.yaml
in: path
      required: true
      description: The realm name.
      schema:
        type: string

    DeviceID:
      name: deviceID
      in: path
      required: true
      description: |
        The device hardware ID: a 128-bit value in the 22-character unpadded
        base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and `_`,
        with no `=` padding. The standard-alphabet characters `+` and `/`, and
        a padded 23-character spelling of the same value, are not accepted.
      schema:
        type: string

  schemas:
    # ── Request types ──

    RegisterRequest:
      type: object
      required: [hw_id]
      properties:
        hw_id:
          type: string
          minLength: 22
          maxLength: 22
          pattern: '^[A-Za-z0-9_-]{22}$'
          description: |
            Device hardware ID: a 128-bit value in the 22-character unpadded
            base64url wire form — the alphabet is `A-Z`, `a-z`, `0-9`, `-` and
            `_`, with no `=` padding. The standard-alphabet characters `+` and
            `/`, and a padded 23-character spelling of the same value, are
            rejected with 422 and the
            `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}`
            body.
        initial_payload_format:
          type: string
          description: Optional Astrate extension specifying the device's initial payload format.

    CredentialsRequest:
      type: object
      required: [csr]
      properties:
        csr:
          type: string
          description: PEM-encoded Certificate Signing Request from the device.

    VerifyRequest:

$ cd /root/astrate-mule && sed -n '570,640p' docs/api/astarte_pairing_api.yaml
detail: Bad request

    Unauthorized:
      description: Unauthorized — missing or invalid JWT/Bearer token.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Unauthorized

    Forbidden:
      description: Forbidden — valid credentials but insufficient permissions or device inhibited.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Forbidden

    DeviceNotFound:
      description: Device not found — unknown device ID.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Device not found

    ValidationError:
      description: Unprocessable Entity — field-level validation errors.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorFields"
          examples:
            blank_hwid:
              summary: Blank hw_id
              value:
                errors:
                  hw_id:
                    - "can't be blank"
            already_registered:
              summary: Device already registered
              value:
                errors:
                  error_name:
                    - already_registered
            limit_reached:
              summary: Registration limit reached
              value:
                errors:
                  error_name:
                    - device_registration_limit_reached
            invalid_hwid:
              summary: Invalid hw_id format
              value:
                errors:
                  hw_id:
                    - "is not a valid base64 encoded 128 bits id"

    TooManyRequests:
      description: Too Many Requests — rate limit exceeded.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:

$ cd /root/astrate-mule && sed -n '160,200p' internal/pairing/service.go
}
	sealed, err := sealer.Seal(realmCA.PrivateKeyDER())
	if err != nil {
		return "", nil, fmt.Errorf("pairing: sealing CA key: %w", err)
	}
	return realmCA.CertificatePEM(), sealed, nil
}

// Register implements flow A: it validates hw_id, enforces the realm's
// device registration limit, generates a 44-character credentials secret,
// stores its bcrypt hash, and returns the secret (shown exactly once).
// Re-registering a device that has not yet requested credentials rotates the
// secret; afterwards it fails with ErrAlreadyRegistered. initialFormat is
// the Astrate initial_payload_format extension ("", "bson" or "json").
func (s *Service) Register(ctx context.Context, realmName, hwID, initialFormat string) (string, error) {
	id, err := deviceid.Parse(hwID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidHWID, err)
	}
	if initialFormat != "" && initialFormat != "bson" && initialFormat != "json" {
		return "", fmt.Errorf("%w: %q", ErrInvalidPayloadFormat, initialFormat)
	}

	realm, err := s.st.GetRealmByName(ctx, realmName)
	if err != nil {
		return "", err
	}

	// Upstream parity: the limit gates every registration attempt, secret
	// rotations for existing devices included.
	if realm.DeviceRegistrationLimit != nil {
		n, err := s.st.CountDevices(ctx, realm.ID)
		if err != nil {
			return "", err
		}
		if n >= int64(*realm.DeviceRegistrationLimit) {
			return "", ErrRegistrationLimitReached
		}
	}

	secret, err := generateSecret()

$ cd /root/astrate-mule && sed -n '165,185p;245,265p;300,320p' internal/pairing/http.go
_ = astarteapi.WriteUnauthorized(w)
		return
	}
	var req credentialsRequest
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &req); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if req.CSR == "" {
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"csr": {"can't be blank"}})
		return
	}

	clientCrt, err := a.svc.Credentials(r.Context(), realm, deviceID, secret, req.CSR, ip)
	if err != nil {
		a.writeServiceError(w, err)
		return
	}
	_ = astarteapi.WriteData(w, http.StatusCreated, credentialsResponse{ClientCrt: clientCrt})
}
		_ = astarteapi.WriteUnauthorized(w)
		return
	}
	var req verifyRequest
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &req); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if req.ClientCrt == "" {
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"client_crt": {"can't be blank"}})
		return
	}

	res, err := a.svc.VerifyCredentials(r.Context(), r.PathValue("realm"), r.PathValue("deviceID"), secret, req.ClientCrt)
	if err != nil {
		a.writeServiceError(w, err)
		return
	}
	if res.Valid {
		_ = astarteapi.WriteData(w, http.StatusOK, verifyValidResponse{
// ["<msg>"]}}), and RPC-originated conflicts surface as the literal
// "error_name" field. One mapping deviates deliberately: upstream answers
// wrong-secret/unknown-device with 403; Astrate uses 401 (frozen in
// docs/ROADMAP.md §5 — semantically authentication, and uniform across
// causes). Inhibited devices stay 403.
func (a *API) writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidHWID):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"hw_id": {"is not a valid base64 encoded 128 bits id"}})
	case errors.Is(err, ErrInvalidPayloadFormat):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"initial_payload_format": {"is invalid"}})
	case errors.Is(err, ErrAlreadyRegistered):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"error_name": {"already_registered"}})
	case errors.Is(err, ErrRegistrationLimitReached):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"error_name": {"device_registration_limit_reached"}})
	case errors.Is(err, ErrInvalidCSR):
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,

$ cd /root/astrate-mule && rg -n "enum:" docs/api/ | head -20; echo "--- makefile ---"; cat docs/Makefile 2>/dev/null | head -40
docs/api/astarte_pairing_api.yaml:445:          enum: [confirmed, pending, inhibited]
docs/api/astarte_pairing_api.yaml:463:          enum: [true]
docs/api/astarte_pairing_api.yaml:471:          enum: [EXPIRED, INVALID, REVOKED]
docs/api/astarte_pairing_api.yaml:482:          enum: [false]
docs/api/astarte_pairing_api.yaml:513:          enum: [ok]
docs/api/astarte_appengine_api.yaml:1480:        enum: [ascending, descending]
docs/api/astarte_appengine_api.yaml:1489:        enum: [structured, table, disjoint_tables]
docs/api/astrate_native_api.yaml:315:            enum: [websocket, sse]
docs/api/astrate_native_api.yaml:389:          enum: [ok]
docs/api/astrate_native_api.yaml:398:          enum: [ok, unavailable]
docs/api/astrate_native_api.yaml:416:              enum: [ok, unhealthy]
docs/api/astarte_realm_management_api.yaml:949:          enum: [datastream, properties]
docs/api/astarte_realm_management_api.yaml:952:          enum: [device, server]
docs/api/astarte_realm_management_api.yaml:955:          enum: [individual, object]
docs/api/astarte_realm_management_api.yaml:979:          enum: [double, integer, boolean, longinteger, string, binaryblob, datetime, doublearray, integerarray, booleanarray, longintegerarray, stringarray, binaryblobarray, datetimearray]
docs/api/astarte_realm_management_api.yaml:982:          enum: [unreliable, guaranteed, unique]
docs/api/astarte_realm_management_api.yaml:986:          enum: [discard, volatile, stored]
docs/api/astarte_realm_management_api.yaml:995:          enum: [no_ttl, use_ttl]
docs/api/astarte_realm_management_api.yaml:1044:          enum: [data_trigger, device_trigger]
docs/api/astarte_realm_management_api.yaml:1047:          enum: [incoming_data, value_change, value_change_applied, path_created, path_removed, value_stored, device_registered, device_connected, device_disconnected, device_empty_cache_received, device_error, incoming_introspection, interface_added, interface_removed, interface_minor_updated]
--- makefile ---
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

$ cd /root/astrate-mule && sed -n '436,450p' docs/api/astarte_pairing_api.yaml
properties:
                broker_url:
                  type: string
                  description: MQTT broker URL for the device.
                ca_crt:
                  type: string
                  description: PEM-encoded CA certificate for the MQTT TLS connection (Astrate extension).
        status:
          type: string
          enum: [confirmed, pending, inhibited]
          description: Device pairing status.
        version:
          type: string
          description: Astrate version string.

Now the enum change:
← Edit docs/api/astarte_pairing_api.yaml
Index: /root/astrate-mule/docs/api/astarte_pairing_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_pairing_api.yaml
+++ /root/astrate-mule/docs/api/astarte_pairing_api.yaml
@@ -386,9 +386,14 @@
         `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}`
         body.
     initial_payload_format:
       type: string
-      description: Optional Astrate extension specifying the device's initial payload format.
+      enum: [bson, json]
+      description: |
+        Optional Astrate extension specifying the device's initial payload
+        format. Any other value is rejected with 422 and the
+        `{"errors":{"initial_payload_format":["is invalid"]}}` body.
+        Omit the field, or send it empty, to leave the format unset.
 
 CredentialsRequest:
   type: object
   required: [csr]

← Edit docs/api/astarte_pairing_api.yaml
Index: /root/astrate-mule/docs/api/astarte_pairing_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_pairing_api.yaml
+++ /root/astrate-mule/docs/api/astarte_pairing_api.yaml
@@ -634,8 +634,26 @@
           value:
             errors:
               hw_id:
                 - "is not a valid base64 encoded 128 bits id"
+        invalid_payload_format:
+          summary: Unsupported initial_payload_format
+          value:
+            errors:
+              initial_payload_format:
+                - "is invalid"
+        blank_csr:
+          summary: Blank csr
+          value:
+            errors:
+              csr:
+                - "can't be blank"
+        blank_client_crt:
+          summary: Blank client_crt
+          value:
+            errors:
+              client_crt:
+                - "can't be blank"
 
 TooManyRequests:
   description: Too Many Requests — rate limit exceeded.
   content:

Now the build gate:
$ cd /root/astrate-mule && make -C docs build 2>&1 | tail -25
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

$ cd /root/astrate-mule && git status --porcelain | head -20; echo "--- venv? ---"; ls -a docs | head -20; echo "--- which ---"; command -v mkdocs python3 pip3 2>&1
M docs/api/astarte_pairing_api.yaml
--- venv? ---
.
..
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
--- which ---

$ cd /root/astrate-mule && rg -n "yaml" go.mod; echo "--- cache ---"; ls $(go env GOMODCACHE)/gopkg.in 2>/dev/null; echo "--- mkdocs refs ---"; rg -n "api|swagger" docs/mkdocs.yml | head -20
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect
--- cache ---
check.v1@v0.0.0-20161208181325-20d25e280405
check.v1@v1.0.0-20201130134442-10cb98267c6c
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1
--- mkdocs refs ---
50:  - AppEngine API: appengine-api.md
51:  - Realm Management API: realm-management-api.md
52:  - Housekeeping API: housekeeping-api.md
63:  - API Explorer: swagger.md

! permission requested: external_directory (/tmp/*); auto-rejecting
✗ mkdir -p /tmp/yamlcheck && cat > /tmp/yamlcheck/go.mod <<'EOF'
module yamlcheck

go 1.21

require gopkg.in/yaml.v3 v3.0.1
EOF
cat > /tmp/yamlcheck/main.go <<'EOF'
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func main() {
	root := os.Args[1]
	bad := 0
	n := 0
	err := filepath.Walk(root, func(p string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return err
		}
		if filepath.Ext(p) != ".yaml" && filepath.Ext(p) != ".yml" {
			return nil
		}
		n++
		b, err := os.ReadFile(p)
		if err != nil {
			fmt.Println("READ FAIL", p, err)
			bad++
			return nil
		}
		var v any
		if err := yaml.Unmarshal(b, &v); err != nil {
			fmt.Println("PARSE FAIL", p, err)
			bad++
			return nil
		}
		fmt.Println("ok", p)
		return nil
	})
	if err != nil {
		fmt.Println("walk error:", err)
		os.Exit(1)
	}
	fmt.Printf("parsed %d yaml files, %d failures\n", n, bad)
	if bad > 0 {
		os.Exit(1)
	}
}
EOF
cd /tmp/yamlcheck && GOFLAGS=-mod=mod GOPROXY=off go run . /root/astrate-mule/docs 2>&1 | tail -20 failed
Error: The user rejected permission to use this specific tool call.
