slug: docs-sync-pairing-body-cap-400
verdict: blocked
at:  d7e14cf
ran: 2026-10-08T19:13:10Z on DietPi in 1200s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

$ grep -rn "BadRequest" docs | head -20
docs/site/api/astarte_realm_management_api.yaml:131:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_realm_management_api.yaml:265:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_realm_management_api.yaml:414:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_realm_management_api.yaml:568:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_realm_management_api.yaml:785:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_realm_management_api.yaml:1309:    BadRequest:
docs/site/api/astarte_appengine_api.yaml:245:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_appengine_api.yaml:333:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_appengine_api.yaml:518:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_appengine_api.yaml:570:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_appengine_api.yaml:613:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_appengine_api.yaml:775:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_appengine_api.yaml:822:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_appengine_api.yaml:860:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_appengine_api.yaml:940:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_appengine_api.yaml:1064:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_appengine_api.yaml:1146:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_appengine_api.yaml:1336:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_appengine_api.yaml:1384:          $ref: "#/components/responses/BadRequest"
docs/site/api/astarte_appengine_api.yaml:1423:          $ref: "#/components/responses/BadRequest"

$ find docs -name "*pairing*" | head -15
docs/site/api/astarte_pairing_api.yaml
docs/site/pairing-and-security.md
docs/api/astarte_pairing_api.yaml

$ grep -rn "BadRequest" docs/site/api/astarte_pairing_api.yaml | head -10
135:          $ref: "#/components/responses/BadRequest"
223:          $ref: "#/components/responses/BadRequest"
330:          $ref: "#/components/responses/BadRequest"
578:    BadRequest:

$ sed -n '560,600p' docs/site/api/astarte_pairing_api.yaml
type: object
          required: [detail]
          properties:
            detail:
              type: string

    ErrorFields:
      type: object
      required: [errors]
      properties:
        errors:
          type: object
          additionalProperties:
            type: array
            items:
              type: string

  responses:
    BadRequest:
      description: Bad Request — malformed or missing data envelope.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
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

$ sed -n '400,440p' docs/site/api/astarte_pairing_api.yaml
type: string
          enum: [bson, json]
          description: |
            Optional Astrate extension specifying the device's initial payload
            format. Any other value is rejected with 422 and the
            `{"errors":{"initial_payload_format":["is invalid"]}}` body.
            Omit the field, or send it empty, to leave the format unset.

    CredentialsRequest:
      type: object
      required: [csr]
      properties:
        csr:
          type: string
          description: PEM-encoded Certificate Signing Request from the device.

    VerifyRequest:
      type: object
      required: [client_crt]
      properties:
        client_crt:
          type: string
          description: PEM-encoded client certificate to verify.

    # ── Response payload types ──

    CredentialsSecret:
      type: object
      required: [credentials_secret]
      properties:
        credentials_secret:
          type: string
          description: One-time device credentials secret (base64-encoded).

    ClientCert:
      type: object
      required: [client_crt]
      properties:
        client_crt:
          type: string
          description: PEM-encoded client certificate issued for the device.

$ sed -n '110,130p' docs/site/api/astarte_pairing_api.yaml | head -30
- $ref: "#/components/parameters/RealmName"
      requestBody:
        required: true
        content:
          application/json; charset=utf-8:
            schema:
              type: object
              required: [data]
              properties:
                data:
                  $ref: "#/components/schemas/RegisterRequest"
            example:
              data:
                hw_id: dT6hS2W9TT6LEnP25ks_lg
      responses:
        "201":
          description: Device registered. Returns the one-time credentials secret.
          content:
            application/json; charset=utf-8:
              schema:
                $ref: "#/components/schemas/DataEnvelopeCredentialsSecret"

$ sed -n '160,180p' docs/site/api/astarte_pairing_api.yaml | head -30
credentials secret or a certificate.

        Returns 204 on success. Because the row is still there, unregistering
        an already-unregistered device is also a 204 no-op; 404 is returned
        only for a device that does not exist in the realm.
      tags: [Pairing]
      security:
        - a_pa: []
      parameters:
        - $ref: "#/components/parameters/RealmName"
        - $ref: "#/components/parameters/DeviceID"
      responses:
        "204":
          description: Device unregistered successfully. No body returned.
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "404":
          $ref: "#/components/responses/DeviceNotFound"
        "500":

$ sed -n '200,230p' docs/site/api/astarte_pairing_api.yaml | head -40
required: true
        content:
          application/json; charset=utf-8:
            schema:
              type: object
              required: [data]
              properties:
                data:
                  $ref: "#/components/schemas/CredentialsRequest"
            example:
              data:
                csr: "-----BEGIN CERTIFICATE REQUEST-----\nMIICY...\n-----END CERTIFICATE REQUEST-----"
      responses:
        "201":
          description: Client certificate issued.
          content:
            application/json; charset=utf-8:
              schema:
                $ref: "#/components/schemas/DataEnvelopeClientCert"
              example:
                data:
                  client_crt: "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----"
        "400":
          $ref: "#/components/responses/BadRequest"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "422":
          $ref: "#/components/responses/ValidationError"
        "429":

$ sed -n '310,340p' docs/site/api/astarte_pairing_api.yaml | head -40
timestamp: "2024-05-30 13:49:57.045Z"
                      until: "2025-05-30 13:49:57.045Z"
                      valid: true
                expired:
                  summary: Expired credentials
                  value:
                    data:
                      cause: EXPIRED
                      details: null
                      timestamp: "2024-05-30 13:49:57.045Z"
                      valid: false
                revoked:
                  summary: Revoked credentials
                  value:
                    data:
                      cause: REVOKED
                      details: null
                      timestamp: "2024-05-30 13:49:57.045Z"
                      valid: false
        "400":
          $ref: "#/components/responses/BadRequest"
        "401":
          $ref: "#/components/responses/Unauthorized"
        "403":
          $ref: "#/components/responses/Forbidden"
        "422":
          $ref: "#/components/responses/ValidationError"
        "500":
          $ref: "#/components/responses/InternalServerError"

components:

$ sed -n '570,590p' docs/api/astarte_pairing_api.yaml | head -30
additionalProperties:
            type: array
            items:
              type: string

  responses:
    BadRequest:
      description: Bad Request — malformed or missing data envelope.
      content:
        application/json; charset=utf-8:
          schema:
            $ref: "#/components/schemas/ErrorDetail"
          example:
            errors:
              detail: Bad request

    Unauthorized:
      description: Unauthorized — missing or invalid JWT/Bearer token.
      content:
        application/json; charset=utf-8:
          schema:

$ sed -n '1,30p' internal/pairing/http.go
package pairing

import (
	"errors"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"github.com/astrate-platform/astrate/internal/auth"
	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/astarteapi"
)

// maxBodyBytes caps pairing request bodies (docs/DESIGN.md §4.5 input
// bounds): CSRs and certificates fit comfortably in 64 KiB.
const maxBodyBytes int64 = 64 << 10

// Default rate-limit parameters (docs/DESIGN.md §4.5). Registration is
// operator-driven (JWT-protected, fleet provisioning bursts); credentials
// requests are device-driven and rare (a renewal per device per cert TTL).
const (
	DefaultRegisterRate     = 5.0
	DefaultRegisterBurst    = 20
	DefaultCredentialsRate  = 1.0
	DefaultCredentialsBurst = 5
)

// detailTooManyRequests is the 429 envelope detail.
const detailTooManyRequests = "Too Many Requests"

$ sed -n '105,125p' internal/pairing/http.go
type registerResponse struct {
	CredentialsSecret string `json:"credentials_secret"`
}

func (a *API) handleRegister(w http.ResponseWriter, r *http.Request) {
	if !a.regLimiter.Allow("ip|" + remoteIP(r).String()) {
		_ = astarteapi.WriteError(w, http.StatusTooManyRequests, detailTooManyRequests)
		return
	}

	var req registerRequest
	if err := astarteapi.DecodeData(r.Body, maxBodyBytes, &req); err != nil {
		_ = astarteapi.WriteBadRequest(w)
		return
	}
	if req.HwID == "" {
		_ = astarteapi.WriteFieldErrors(w, http.StatusUnprocessableEntity,
			map[string][]string{"hw_id": {"can't be blank"}})
		return
	}

$ sed -n '160,175p' internal/pairing/http.go
return
	}

	secret, ok := bearerSecret(r)
	if !ok {
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

$ sed -n '240,255p' internal/pairing/http.go
}

func (a *API) handleVerify(w http.ResponseWriter, r *http.Request) {
	secret, ok := bearerSecret(r)
	if !ok {
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

$ sed -n '220,240p' pkg/astarteapi/envelope.go
// DecodeData reads at most maxBytes bytes from r, unwraps the mandatory
// {"data": ...} request envelope, and unmarshals the "data" value into dst.
//
// It fails with an error wrapping ErrBodyTooLarge when the body exceeds
// maxBytes, with one wrapping ErrMissingData when the "data" key is absent or
// null, and with the underlying JSON error for malformed bodies (including
// trailing garbage after the top-level value). Sibling keys next to "data"
// are ignored, matching upstream parameter handling.
func DecodeData(r io.Reader, maxBytes int64, dst any) error {
	body, err := io.ReadAll(io.LimitReader(r, maxBytes+1))
	if err != nil {
		return fmt.Errorf("astarteapi: reading request body: %w", err)
	}
	if int64(len(body)) > maxBytes {
		return fmt.Errorf("astarteapi: %w: exceeds %d bytes", ErrBodyTooLarge, maxBytes)
	}
	var env struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("astarteapi: decoding request envelope: %w", err)

$ grep -rn "ErrBodyTooLarge" pkg/astarteapi
pkg/astarteapi/envelope_test.go:242:		{"oversized", `{"data": "` + strings.Repeat("a", 64) + `"}`, 32, astarteapi.ErrBodyTooLarge},
pkg/astarteapi/envelope_test.go:243:		{"oversized by one", `{"data": 11}`, 11, astarteapi.ErrBodyTooLarge},
pkg/astarteapi/envelope.go:65:// ErrBodyTooLarge is wrapped by DecodeData when the request body exceeds the
pkg/astarteapi/envelope.go:67:var ErrBodyTooLarge = errors.New("request body too large")
pkg/astarteapi/envelope.go:223:// It fails with an error wrapping ErrBodyTooLarge when the body exceeds
pkg/astarteapi/envelope.go:234:		return fmt.Errorf("astarteapi: %w: exceeds %d bytes", ErrBodyTooLarge, maxBytes)

$ grep -rn "WriteBadRequest" pkg/astarteapi
pkg/astarteapi/envelope_test.go:90:			write: astarteapi.WriteBadRequest,
pkg/astarteapi/envelope.go:190:// WriteBadRequest writes the canonical 400 envelope.
pkg/astarteapi/envelope.go:191:func WriteBadRequest(w http.ResponseWriter) error {

