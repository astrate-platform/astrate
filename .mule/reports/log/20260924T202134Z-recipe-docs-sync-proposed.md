slug: recipe-docs-sync
verdict: proposed
at:  6b148bf
ran: 2026-09-24T20:21:34Z on DietPi in 322s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls docs/api/ && printf -- '---\n' && ls internal/ && printf -- '---\n' && rg -l 'appengine|housekeeping|pairing|realm.management' internal/ -g '*.go' | head -50
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
---
appengine
auth
broker
config
engine
flow
flowapi
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil
---
internal/realm/http.go
internal/realm/dashboard_compat_test.go
internal/realm/http_test.go
internal/flow/blocks/virtualdevicepool/virtualdevicepool_test.go
internal/realm/service.go
internal/flow/blocks/virtualdevicepool/virtualdevicepool.go
internal/observability/compat.go
internal/observability/observability_test.go
internal/observability/compat_version_test.go
internal/flow/factory.go
internal/config/config.go
internal/engine/engine.go
internal/engine/router.go
internal/httpx/cors_test.go
internal/httpx/notfound_test.go
internal/httpx/notfound.go
internal/appengine/data.go
internal/appengine/groups_parity_test.go
internal/appengine/http.go
internal/appengine/query_opts_test.go
internal/appengine/formats_parity_test.go
internal/appengine/dashboard_compat_test.go
internal/appengine/channels/ws_test.go
internal/appengine/channels/protocol.go
internal/appengine/service.go
internal/auth/middleware.go
internal/auth/middleware_test.go
internal/auth/claims_test.go
internal/testutil/mqttclient.go
internal/store/devices.go
internal/swagger/swagger_test.go
internal/housekeeping/http.go
internal/housekeeping/service_test.go
internal/housekeeping/async_operation_test.go
internal/housekeeping/http_test.go
internal/housekeeping/service.go
internal/appengine/downsample.go
internal/appengine/http_test.go
internal/appengine/patch_test.go
internal/appengine/channels/ws.go
internal/appengine/e2e_test.go
internal/appengine/writeerror_test.go
internal/appengine/mirror_test.go
internal/appengine/downsample_test.go
internal/engine/engine_e2e_test.go
internal/pairing/remoteip_test.go
internal/pairing/http.go
internal/store/realms.go
internal/pairing/bearersecret_test.go
internal/pairing/service_test.go

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/pairing/*.go
internal/pairing/http.go:81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/pairing/http.go:86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
internal/pairing/http.go:88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
internal/pairing/http.go:90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
internal/pairing/http.go:344:	header := r.Header.Get("Authorization")

$ rg -n '^\s{2}/' docs/api/astarte_pairing_api.yaml
27:  /pairing/v1/{realm}/health:
70:  /pairing/v1/{realm}/version:
97:  /pairing/v1/{realm}/agent/devices:
149:  /pairing/v1/{realm}/agent/devices/{deviceID}:
172:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
224:  /pairing/v1/{realm}/devices/{deviceID}:
258:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:

$ rg -n 'agent/devices|pairing/v1.*version|HandleFunc.*version' -g '*.go' internal/ | rg -v '_test'
internal/store/devices.go:94:// (DELETE /agent/devices parity, docs/DESIGN.md §4.4): the credentials
internal/pairing/http.go:82:	mux.Handle("POST /pairing/v1/{realm}/agent/devices",
internal/pairing/http.go:84:	mux.Handle("DELETE /pairing/v1/{realm}/agent/devices/{deviceID}",
internal/observability/compat.go:51:	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
internal/auth/middleware.go:163:// Example: RelativePath("/pairing/v1/test/agent/devices", "test") returns
internal/auth/middleware.go:164:// ("agent/devices", true).
internal/auth/claims.go:161:// (the request path relative to the realm base, e.g. "agent/devices").

The documented-but-absent candidates: `/version`, `/agent/devices`, `/agent/devices/{deviceID}`. Let me verify they genuinely don't exist in the code (not just in the pairing package) and check one of the documented endpoints' status codes.
$ rg -n -A 60 '/pairing/v1/\{realm\}/devices/\{deviceID\}/protocols/astarte_mqtt_v1/credentials' docs/api/astarte_pairing_api.yaml
172:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
173-    post:
174-      operationId: requestCredentials
175-      summary: Request device credentials (flow B)
176-      description: |
177-        Issues a client certificate for the device's MQTT connection. The device
178-        authenticates with its credentials secret as a Bearer token and submits
179-        a CSR.
180-
181-        Rate-limited per IP and per device.
182-      tags: [Pairing]
183-      security:
184-        - bearerSecret: []
185-      parameters:
186-        - $ref: "#/components/parameters/RealmName"
187-        - $ref: "#/components/parameters/DeviceID"
188-      requestBody:
189-        required: true
190-        content:
191-          application/json; charset=utf-8:
192-            schema:
193-              type: object
194-              required: [data]
195-              properties:
196-                data:
197-                  $ref: "#/components/schemas/CredentialsRequest"
198-            example:
199-              data:
200-                csr: "-----BEGIN CERTIFICATE REQUEST-----\nMIICY...\n-----END CERTIFICATE REQUEST-----"
201-      responses:
202-        "201":
203-          description: Client certificate issued.
204-          content:
205-            application/json; charset=utf-8:
206-              schema:
207-                $ref: "#/components/schemas/DataEnvelopeClientCert"
208-              example:
209-                data:
210-                  client_crt: "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----"
211-        "400":
212-          $ref: "#/components/responses/BadRequest"
213-        "401":
214-          $ref: "#/components/responses/Unauthorized"
215-        "403":
216-          $ref: "#/components/responses/Forbidden"
217-        "422":
218-          $ref: "#/components/responses/ValidationError"
219-        "429":
220-          $ref: "#/components/responses/TooManyRequests"
221-        "500":
222-          $ref: "#/components/responses/InternalServerError"
223-
224-  /pairing/v1/{realm}/devices/{deviceID}:
225-    get:
226-      operationId: getDeviceInfo
227-      summary: Get device pairing info (flow C)
228-      description: |
229-        Returns the device's pairing status, protocol details (broker URL,
230-        CA certificate), and version. The device authenticates with its
231-        credentials secret as a Bearer token.
232-      tags: [Pairing]
--
258:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:
259-    post:
260-      operationId: verifyCredentials
261-      summary: Verify device credentials
262-      description: |
263-        Verifies a device's client certificate and returns its validity status
264-        (valid, expired, revoked, or invalid). The device authenticates with
265-        its credentials secret as a Bearer token.
266-      tags: [Pairing]
267-      security:
268-        - bearerSecret: []
269-      parameters:
270-        - $ref: "#/components/parameters/RealmName"
271-        - $ref: "#/components/parameters/DeviceID"
272-      requestBody:
273-        required: true
274-        content:
275-          application/json; charset=utf-8:
276-            schema:
277-              type: object
278-              required: [data]
279-              properties:
280-                data:
281-                  $ref: "#/components/schemas/VerifyRequest"
282-            example:
283-              data:
284-                client_crt: "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----"
285-      responses:
286-        "200":
287-          description: Credential verification result.
288-          content:
289-            application/json; charset=utf-8:
290-              schema:
291-                oneOf:
292-                  - $ref: "#/components/schemas/DataEnvelopeVerifyValid"
293-                  - $ref: "#/components/schemas/DataEnvelopeVerifyInvalid"
294-              examples:
295-                valid:
296-                  summary: Valid credentials
297-                  value:
298-                    data:
299-                      timestamp: "2024-05-30 13:49:57.045Z"
300-                      until: "2025-05-30 13:49:57.045Z"
301-                      valid: true
302-                expired:
303-                  summary: Expired credentials
304-                  value:
305-                    data:
306-                      cause: EXPIRED
307-                      details: null
308-                      timestamp: "2024-05-30 13:49:57.045Z"
309-                      valid: false
310-                revoked:
311-                  summary: Revoked credentials
312-                  value:
313-                    data:
314-                      cause: REVOKED
315-                      details: null
316-                      timestamp: "2024-05-30 13:49:57.045Z"
317-                      valid: false
318-        "400":

The agent endpoints exist — they use `mux.Handle` (not `HandleFunc`), which my and the recipe's grep pattern missed. Let me verify the `/version` path in observability and the health handler.
→ Read internal/observability/compat.go
→ Read internal/pairing/http.go [offset=75, limit=25]

$ rg -n 'MountVersionCompat|MountServiceCompat' -g '*.go' internal/ | rg -v '_test'
internal/observability/compat.go:9:// MountServiceCompat registers the upstream-parity unauthenticated
internal/observability/compat.go:19:func MountServiceCompat(mux *http.ServeMux, service string, check Check) {
internal/observability/compat.go:49:// MountVersionCompat registers GET /{service}/version for one service.
internal/observability/compat.go:50:func MountVersionCompat(mux *http.ServeMux, service, version string) {

Version endpoint isn't in pairing's `Mount`. Docs document `/pairing/v1/{realm}/version` — let me check where `MountVersionCompat` is actually mounted for pairing.
$ rg -n -A 25 '/pairing/v1/\{realm\}/version' docs/api/astarte_pairing_api.yaml && rg -n -A 30 '/pairing/v1/\{realm\}/health' docs/api/astarte_pairing_api.yaml
70:  /pairing/v1/{realm}/version:
71-    get:
72-      operationId: getVersion
73-      summary: Get emulated API version
74-      description: |
75-        Returns the emulated upstream API compatibility version. Unauthenticated
76-        — mirrors upstream's public pairing version endpoint.
77-
78-        The Astarte Dashboard gates feature UI on this value.
79-      tags: [Pairing]
80-      parameters:
81-        - $ref: "#/components/parameters/RealmName"
82-      responses:
83-        "200":
84-          description: API compatibility version.
85-          content:
86-            application/json; charset=utf-8:
87-              schema:
88-                type: object
89-                required: [data]
90-                properties:
91-                  data:
92-                    type: string
93-                    description: Version string.
94-              example:
95-                data: "1.2.2"
27:  /pairing/v1/{realm}/health:
28-    get:
29-      operationId: getHealth
30-      summary: Realm health check
31-      description: |
32-        Returns the health status of a realm's backing store. Unauthenticated
33-        and unrate-limited — FDO-flow devices probe it before they hold
34-        credentials.
35-
36-        Returns 200 when the realm exists and the database answers, 404 for an
37-        unknown realm, and 503 when the database is unhealthy.
38-      tags: [Pairing]
39-      parameters:
40-        - $ref: "#/components/parameters/RealmName"
41-      responses:
42-        "200":
43-          description: Realm is healthy.
44-          content:
45-            application/json; charset=utf-8:
46-              schema:
47-                $ref: "#/components/schemas/DataEnvelopeHealthStatus"
48-              example:
49-                data:
50-                  status: ok
51-        "404":
52-          description: Unknown realm.
53-          content:
54-            application/json; charset=utf-8:
55-              schema:
56-                $ref: "#/components/schemas/ErrorDetail"
57-              example:

Callers of `MountVersionCompat`/`MountServiceCompat` aren't in `internal/`. Let me search the whole repo.
$ rg -n 'MountVersionCompat|MountServiceCompat|checkVersion|/version' -g '*.go' -g '!**/_test.go' --glob '!internal/**' . | head -30
./cmd/astrate/main.go:384:		observability.MountServiceCompat(mux, svc, st.Health)
./cmd/astrate/main.go:388:	// service answers an unauthenticated GET /{service}/version with
./cmd/astrate/main.go:394:		observability.MountVersionCompat(mux, svc, version)
./cmd/astrate/main.go:433:	mux.Handle("GET /appengine/v1/{realm}/version",
./cmd/astrate/main.go:435:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
./cmd/astrate/version_test.go:21:	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/pairing/v1/somerealm/version", nil))

Pairing paths are fully in sync (the missing ones exist via `mux.Handle` in `cmd/astrate/main.go:435` and `internal/pairing/http.go:82,84`). Now checking 3-4 documented status codes/fields against handlers.
$ rg -n 'func \(a \*API\) handle' internal/pairing/http.go
109:func (a *API) handleRegister(w http.ResponseWriter, r *http.Request) {
134:func (a *API) handleUnregister(w http.ResponseWriter, r *http.Request) {
155:func (a *API) handleCredentials(w http.ResponseWriter, r *http.Request) {
201:func (a *API) handleInfo(w http.ResponseWriter, r *http.Request) {
242:func (a *API) handleVerify(w http.ResponseWriter, r *http.Request) {
284:func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {

→ Read cmd/astrate/main.go [offset=378, limit=62]

$ rg -n 'APICompatVersion' internal/realm/*.go | rg -v _test
internal/realm/http.go:119:// feature UI on it — see APICompatVersion).
internal/realm/http.go:121:	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
internal/realm/service.go:580:// APICompatVersion is the upstream Realm Management API level Astrate
internal/realm/service.go:588:const APICompatVersion = "1.2.2"

Let me read the documented `/version` behavior in the spec's example (`data: "1.2.2"`) and check what `realm.APICompatVersion` actually is, plus spot-check a few handler status codes.
→ Read internal/pairing/http.go [limit=100, offset=284]
→ Read docs/api/astarte_pairing_api.yaml [limit=70, offset=224]
All documented paths exist and `/version` returns exactly `"1.2.2"` matching the docs. Let me spot-check the remaining handlers' status codes.
→ Read internal/pairing/http.go [limit=145, offset=109]
→ Read docs/api/astarte_pairing_api.yaml [limit=102, offset=70]
Now checking the service-level errors `Info`/`Credentials` can return, and `WriteDeviceNotFound`, to confirm the 404/403 reachability.

$ rg -n 'func WriteDeviceNotFound|func WriteNotFound|func WriteForbidden|func WriteUnauthorized' internal/astarteapi/ internal/httpx/ -g '*.go' 2>/dev/null; rg -rn 'func WriteDeviceNotFound' -g '*.go' internal/
(no output)

→ Read docs/api/astarte_pairing_api.yaml [limit=40, offset=318]

$ rg -rn 'WriteDeviceNotFound' -g '*.go' . | rg -v _test
./pkg/astarteapi/envelope.go:// n writes the canonical 404 "Device not found" envelope.
./pkg/astarteapi/envelope.go:func n(w http.ResponseWriter) error {
./internal/pairing/http.go:		_ = astarteapi.n(w)
./internal/appengine/http.go:		_ = astarteapi.n(w)
./internal/realm/http.go:		_ = astarteapi.n(w)

$ rg -n 'func WriteDeviceNotFound' -g '*.go' internal/ && rg -n 'func \(s \*Service\) Info|ErrNotFound|not found' internal/pairing/service.go
(no output)

$ ls pkg/astarteapi/ 2>/dev/null || echo "no pkg/astarteapi"; rg -n 'WriteDeviceNotFound' pkg/astarteapi/envelope.go
envelope.go
envelope_test.go
testdata
210:// WriteDeviceNotFound writes the canonical 404 "Device not found" envelope.
211:func WriteDeviceNotFound(w http.ResponseWriter) error {

$ rg -n 'func.*Info|ErrNotFound|ErrUnauthorized|ErrInhibited' internal/pairing/service.go
58:	// ErrUnauthorized is the uniform device-authentication failure: unknown
61:	ErrUnauthorized = errors.New("pairing: unauthorized")
62:	// ErrInhibited reports a device blocked by credentials_inhibited (403).
63:	ErrInhibited = errors.New("pairing: credentials request inhibited")
229:// credential trail). store.ErrNotFound is returned for unknown devices.
233:		return fmt.Errorf("%w: device %q", store.ErrNotFound, deviceIDStr)
243:// device by credentials secret (uniform ErrUnauthorized on any mismatch),
253:		return "", ErrInhibited
286:func (s *Service) Info(ctx context.Context, realmName, deviceIDStr, secret string) (*Info, error) {
330:		return nil, ErrInhibited
371:		return err // store.ErrNotFound → 404
378:// device, wrong secret) return the same ErrUnauthorized after a comparable
384:		return nil, nil, ErrUnauthorized
387:	if errors.Is(err, store.ErrNotFound) {
389:		return nil, nil, ErrUnauthorized
395:	if errors.Is(err, store.ErrNotFound) {
397:		return nil, nil, ErrUnauthorized
405:		return nil, nil, ErrUnauthorized
408:		return nil, nil, ErrUnauthorized

→ Read internal/pairing/http.go [limit=31, offset=253]
→ Read internal/pairing/service.go [limit=60, offset=220]
→ Read internal/pairing/service.go [limit=55, offset=280]
→ Read internal/pairing/service.go [limit=70, offset=340]

$ rg -n 'func \(a \*API\) requireAgent' -A 20 internal/pairing/http.go
(no output)

$ rg -n 'requireAgent' internal/pairing/*.go; rg -rn 'requireAgent' -g '*.go' internal/auth/ | head
internal/pairing/http.go:46:	requireAgent func(http.Handler) http.Handler
internal/pairing/http.go:68:		requireAgent: mw.RequireRealm(auth.ClaimPairing),
internal/pairing/http.go:83:		a.requireAgent(http.HandlerFunc(a.handleRegister)))
internal/pairing/http.go:85:		a.requireAgent(http.HandlerFunc(a.handleUnregister)))

Device-facing endpoints collapse unknown-device into 401 (service.go:380-408), so those docs are correct. Let me check one more thing: the service-level `/pairing/health` and `/pairing/version` routes mounted in main.go:384-394 — are they documented anywhere?
$ for f in docs/api/*.yaml; do echo "== $f"; rg -n '^\s{2}/(pairing|appengine|housekeeping|realmmanagement|astrate)/' "$f"; done
== docs/api/astarte_appengine_api.yaml
26:  /appengine/v1/{realm}/devices:
98:  /appengine/v1/{realm}/stats/devices:
136:  /appengine/v1/{realm}/devices/{device}:
210:  /appengine/v1/{realm}/devices-by-alias/{alias}:
292:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces:
324:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}:
369:  /appengine/v1/{realm}/devices-by-alias/{alias}/interfaces/{interface}/{path}:
555:  /appengine/v1/{realm}/devices/{device}/interfaces:
583:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
623:  /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
789:  /appengine/v1/{realm}/groups:
863:  /appengine/v1/{realm}/groups/{group}:
891:  /appengine/v1/{realm}/groups/{group}/devices:
981:  /appengine/v1/{realm}/groups/{group}/devices/{device}:
1080:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces:
1108:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}:
1149:  /appengine/v1/{realm}/groups/{group}/devices/{device}/interfaces/{interface}/{path}:
== docs/api/astarte_housekeeping_api.yaml
26:  /housekeeping/v1/realms:
106:  /housekeeping/v1/realms/{realm}:
== docs/api/astarte_pairing_api.yaml
27:  /pairing/v1/{realm}/health:
70:  /pairing/v1/{realm}/version:
97:  /pairing/v1/{realm}/agent/devices:
149:  /pairing/v1/{realm}/agent/devices/{deviceID}:
172:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
224:  /pairing/v1/{realm}/devices/{deviceID}:
258:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:
== docs/api/astarte_realm_management_api.yaml
27:  /realmmanagement/v1/{realm}/interfaces:
171:  /realmmanagement/v1/{realm}/interfaces/{name}:
202:  /realmmanagement/v1/{realm}/interfaces/{name}/{major}:
349:  /realmmanagement/v1/{realm}/triggers:
446:  /realmmanagement/v1/{realm}/triggers/{name}:
502:  /realmmanagement/v1/{realm}/config/auth:
575:  /realmmanagement/v1/{realm}/config/device_registration_limit:
608:  /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention:
645:  /realmmanagement/v1/{realm}/version:
680:  /realmmanagement/v1/{realm}/devices/{device}:
708:  /realmmanagement/v1/{realm}/policies:
782:  /realmmanagement/v1/{realm}/policies/{name}:
== docs/api/astrate_native_api.yaml
26:  /astrate/v1/health:
44:  /astrate/v1/readiness:
77:  /astrate/v1/metrics:
99:  /appengine/health:
127:  /realmmanagement/health:
155:  /pairing/health:
185:  /appengine/version:
205:  /realmmanagement/version:
225:  /pairing/version:
245:  /housekeeping/version:
267:  /astrate/v1/{realm}/socket:
332:  /appengine/v1/socket/websocket:

Service-level endpoints live in `astrate_native_api.yaml` (intentional split), so no drift there either. Let me run the config-key sweep from the recipe as a final check.
$ ls docs/site/ | rg -i 'config|reference' ; rg -o '\bASTRATE_[A-Z_]+' -N internal/ | sort -u
configuration-reference.md
internal/config/config.example.toml:ASTRATE_DATABASE_DSN
internal/config/config.example.toml:ASTRATE_ENGINE_SHARDS
internal/config/config.example.toml:ASTRATE_HTTP_ADDR
internal/config/config.example.toml:ASTRATE_LOG_LEVEL
internal/config/config.example.toml:ASTRATE_MASTER_KEY
internal/config/config.example.toml:ASTRATE_MASTER_KEY_FILE
internal/config/config.example.toml:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/config/config.go:ASTRATE_DATABASE_DSN
internal/config/config.go:ASTRATE_ENGINE_SHARDS
internal/config/config.go:ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
internal/config/config.go:ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
internal/config/config.go:ASTRATE_HTTP_ADDR
internal/config/config.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
internal/config/config.go:ASTRATE_HTTP_TLS_CERT_FILE
internal/config/config.go:ASTRATE_HTTP_TLS_KEY_FILE
internal/config/config.go:ASTRATE_LOG_FORMAT
internal/config/config.go:ASTRATE_LOG_LEVEL
internal/config/config.go:ASTRATE_MASTER_KEY
internal/config/config.go:ASTRATE_MASTER_KEY_FILE
internal/config/config.go:ASTRATE_MQTT_ADDR
internal/config/config.go:ASTRATE_MQTT_ADVERTISED_URL
internal/config/config.go:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/config/config.go:ASTRATE_MQTT_SESSION_STORE_PATH
internal/config/config.go:ASTRATE_MQTT_TLS_CERT_FILE
internal/config/config.go:ASTRATE_MQTT_TLS_KEY_FILE
internal/config/config.go:ASTRATE_REALM_JWT_PUBLIC_KEY
internal/config/config.go:ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
internal/config/config.go:ASTRATE_REALM_NAME
internal/config/config.go:ASTRATE_SECURITY_MASTER_KEY_FILE
internal/config/config_test.go:ASTRATE_DATABASE_DSN
internal/config/config_test.go:ASTRATE_ENGINE_SHARDS
internal/config/config_test.go:ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
internal/config/config_test.go:ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
internal/config/config_test.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
internal/config/config_test.go:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/engine/engine_integration_test.go:ASTRATE_TEST_DSN
internal/flow/blocks/container/block.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/container/docker.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/info.go:ASTRATE_FLOW_CONFIG
internal/flow/blocks/schema.go:ASTRATE_FLOW_CONFIG
internal/store/crypto.go:ASTRATE_MASTER_KEY
internal/store/crypto.go:ASTRATE_MASTER_KEY_FILE
internal/store/store_test.go:ASTRATE_TEST_DSN
internal/testutil/pg.go:ASTRATE_TEST_DSN

$ rg -o '\bASTRATE_[A-Z_]+' -N docs/site/configuration-reference.md | sort -u
ASTRATE_DATABASE_DSN
ASTRATE_ENGINE_SHARDS
ASTRATE_HTTP_ADDR
ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
ASTRATE_HTTP_TLS_CERT_FILE
ASTRATE_HTTP_TLS_KEY_FILE
ASTRATE_LOG_FORMAT
ASTRATE_LOG_LEVEL
ASTRATE_MASTER_KEY
ASTRATE_MASTER_KEY_FILE
ASTRATE_MQTT_ADDR
ASTRATE_MQTT_ADVERTISED_URL
ASTRATE_MQTT_INSECURE_DEV_MODE
ASTRATE_MQTT_SESSION_STORE_PATH
ASTRATE_MQTT_TLS_CERT_FILE
ASTRATE_MQTT_TLS_KEY_FILE
ASTRATE_REALM_JWT_PUBLIC_KEY
ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
ASTRATE_REALM_NAME
ASTRATE_SECURITY_MASTER_KEY_FILE

Only config-key drift: three keys exist in code but not in `configuration-reference.md`. Let me confirm they're not documented elsewhere in docs/site and verify they're real runtime keys.
$ rg -rn 'ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION|ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED|ASTRATE_FLOW_CONFIG' docs/ .mule/ 2>/dev/null; echo "---"; rg -n 'ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION|ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED' internal/config/config.go
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_0d513a4dd001DVcgjXpul9NP4h

.mule/cron.log:internal/config/config.gon
.mule/cron.log:internal/config/config_test.gon
.mule/cron.log:internal/config/config_test.gon
.mule/cron.log:internal/flow/blocks/container/block.gon
.mule/cron.log:internal/flow/blocks/container/docker.gon
.mule/cron.log:internal/flow/blocks/info.gon
.mule/cron.log:internal/flow/blocks/schema.gon
.mule/cron.log:- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/cron.log:n
.mule/cron.log:n
.mule/cron.log:n
.mule/cron.log:$ rg -n 'n|n|n|ASTRATE_TEST_DSN' internal/ | head -20
.mule/cron.log:internal/flow/blocks/info.go:116:		Config:  "image (required), config (object→n), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",
.mule/cron.log:internal/flow/blocks/schema.go:19:	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as n"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
.mule/cron.log:internal/flow/blocks/container/block.go:45:	Nested        map[string]any // opaque JSON for n
.mule/cron.log:internal/flow/blocks/container/docker.go:19:	// FlowConfigJSON is passed as n.
.mule/cron.log:internal/flow/blocks/container/docker.go:99:		args = append(args, "-e", "n="+spec.FlowConfigJSON)
.mule/cron.log:internal/flow/blocks/container/docker.go:201:// encodeFlowConfigJSON marshals nested config for n.
.mule/cron.log:internal/config/config.go:279:		"n",
.mule/cron.log:internal/config/config.go:297:	if v, ok := os.LookupEnv("n"); ok {
.mule/cron.log:internal/config/config_test.go:392:		prefixed = "n"
.mule/cron.log:internal/config/config_test.go:533:	const env = "n"
.mule/cron.log:$ rg -l 'n|n|n' docs/ || echo "absent from docs/"
.mule/cron.log:- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/cron.log:+- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`n`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`n`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`n`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
.mule/cron.log:- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`n`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`n`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`n`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
.mule/cron.log:- `.mule/for-giulio.md` (newest-first): three `ASTRATE_` keys exist in code but are absent from `configuration-reference.md` (`n`, `n` — real startup keys; `n` — per-theme flow-block env contract, your call where it belongs). Reverse direction clean.
.mule/cron.log:Unsure: whether n belongs in configuration-reference or a flow page (flagged in the escalation)
.mule/cron.log: - **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`n`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`n`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`n`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
.mule/cron.log:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `n` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `n` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
.mule/cron.log:  `n` and
.mule/cron.log:  `n` are consumed in `internal/config/config.go`
.mule/cron.log: - **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`n`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`n`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`n`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
.mule/cron.log:- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`n`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`n`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`n`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
.mule/cron.log:	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as n"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
.mule/cron.log:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `n` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `n` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
.mule/cron.log:  `n` and
.mule/cron.log:  `n` are consumed in `internal/config/config.go`
.mule/cron.log:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `n` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `n` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
.mule/cron.log:internal/config/config.go:n
.mule/cron.log:internal/config/config.go:n
.mule/cron.log:internal/config/config_test.go:n
.mule/cron.log:internal/config/config_test.go:n
.mule/cron.log:internal/flow/blocks/container/block.go:n
.mule/cron.log:internal/flow/blocks/container/docker.go:n
.mule/cron.log:internal/flow/blocks/info.go:n
.mule/cron.log:internal/flow/blocks/schema.go:n
.mule/reports/log/20260908T112032Z-recipe-github-issues-proposed.md:  `n` and
.mule/reports/log/20260908T112032Z-recipe-github-issues-proposed.md:  `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260912T172901Z-recipe-github-issues-proposed.md: - **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `n` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `n` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
.mule/reports/log/20260915T113026Z-recipe-code-review-proposed.md:.mule/todo.md:- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/todo.md:- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/reports/log/20260920T185956Z-recipe-docs-sync-proposed.md:- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/reports/log/20260916T184848Z-recipe-code-review-proposed.md:- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/reports/log/20260919T113513Z-recipe-docs-sync-proposed.md:internal/config/config.gon
.mule/reports/log/20260919T113513Z-recipe-docs-sync-proposed.md:internal/config/config.gon
.mule/reports/log/20260919T113513Z-recipe-docs-sync-proposed.md:internal/config/config_test.gon
.mule/reports/log/20260919T113513Z-recipe-docs-sync-proposed.md:internal/config/config_test.gon
.mule/reports/log/20260919T113513Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/block.gon
.mule/reports/log/20260919T113513Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/docker.gon
.mule/reports/log/20260919T113513Z-recipe-docs-sync-proposed.md:internal/flow/blocks/info.gon
.mule/reports/log/20260919T113513Z-recipe-docs-sync-proposed.md:internal/flow/blocks/schema.gon
.mule/reports/log/20260923T174341Z-recipe-milestones-proposed.md: - **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`n`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`n`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`n`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
.mule/reports/log/20260912T171256Z-recipe-milestones-proposed.md: - **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `n` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `n` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
.mule/reports/log/20260913T203209Z-recipe-milestones-proposed.md:83482e4 mule: config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/reports/log/20260906T180442Z-recipe-milestones-proposed.md:   `n` and
.mule/reports/log/20260906T180442Z-recipe-milestones-proposed.md:   `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260910T170149Z-recipe-code-review-proposed.md:  `n` and
.mule/reports/log/20260910T170149Z-recipe-code-review-proposed.md:  `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/config/config.go:n
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/config/config.go:n
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/config/config_test.go:n
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/config/config_test.go:n
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/container/block.go:n
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/container/docker.go:n
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/info.go:n
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/schema.go:n
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:- `n` (config.go)
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:- `n` (config.go)
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:The two housekeeping keys (`n` and `n`) are in code but not in `configuration-reference.md`. That's a `for-giulio.md` escalation (site prose is never-touch).
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:3. Config keys `n` and `n` exist in config.go but are undocumented in `docs/site/configuration-reference.md`.
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:Let me identify the config keys in code missing from configuration-reference to nail down the for-giulio entry. From the earlier diff: `n` and `n` are in config.go but absent from configuration-reference.md. Let me confirm they're actual config fields (not just string literals) and check the for-giulio header format.
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:270:		"n",
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:288:	if v, ok := os.LookupEnv("n"); ok {
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:+  `n` and
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:+  `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:  `n` and
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:  `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md:- **Site prose (for-giulio.md, never-touch):** `n` and `n` exist in config.go but are missing from `configuration-reference.md`.
.mule/reports/log/20260916T115527Z-recipe-docs-sync-proposed.md:- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/reports/log/20260916T115527Z-recipe-docs-sync-proposed.md:internal/config/config.go:n
.mule/reports/log/20260916T115527Z-recipe-docs-sync-proposed.md:internal/config/config.go:n
.mule/reports/log/20260916T115527Z-recipe-docs-sync-proposed.md:internal/config/config_test.go:n
.mule/reports/log/20260916T115527Z-recipe-docs-sync-proposed.md:internal/config/config_test.go:n
.mule/reports/log/20260916T115527Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/block.go:n
.mule/reports/log/20260916T115527Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/docker.go:n
.mule/reports/log/20260916T115527Z-recipe-docs-sync-proposed.md:internal/flow/blocks/info.go:n
.mule/reports/log/20260916T115527Z-recipe-docs-sync-proposed.md:internal/flow/blocks/schema.go:n
.mule/reports/log/20260916T115527Z-recipe-docs-sync-proposed.md:279:		"n",
.mule/reports/log/20260916T115527Z-recipe-docs-sync-proposed.md:297:	if v, ok := os.LookupEnv("n"); ok {
.mule/reports/log/20260923T174805Z-recipe-github-issues-proposed.md:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `n` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `n` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
.mule/reports/log/20260923T174805Z-recipe-github-issues-proposed.md:  `n` and
.mule/reports/log/20260923T174805Z-recipe-github-issues-proposed.md:  `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260923T174805Z-recipe-github-issues-proposed.md: - **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`n`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`n`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`n`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
.mule/reports/log/20260917T190323Z-recipe-github-issues-proposed.md:- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/reports/log/20260918T114100Z-recipe-astarte-upstream-proposed.md: - [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/reports/log/20260913T203830Z-recipe-github-issues-proposed.md:- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/reports/log/20260916T183922Z-recipe-github-issues-proposed.md:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `n` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `n` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
.mule/reports/log/20260916T183922Z-recipe-github-issues-proposed.md:  `n` and
.mule/reports/log/20260916T183922Z-recipe-github-issues-proposed.md:  `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260907T191157Z-recipe-docs-sync-proposed.md:internal/config/config.go:n
.mule/reports/log/20260907T191157Z-recipe-docs-sync-proposed.md:internal/config/config.go:n
.mule/reports/log/20260907T191157Z-recipe-docs-sync-proposed.md:internal/config/config_test.go:n
.mule/reports/log/20260907T191157Z-recipe-docs-sync-proposed.md:internal/config/config_test.go:n
.mule/reports/log/20260907T191157Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/block.go:n
.mule/reports/log/20260907T191157Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/docker.go:n
.mule/reports/log/20260907T191157Z-recipe-docs-sync-proposed.md:internal/flow/blocks/info.go:n
.mule/reports/log/20260907T191157Z-recipe-docs-sync-proposed.md:internal/flow/blocks/schema.go:n
.mule/reports/log/20260907T191157Z-recipe-docs-sync-proposed.md:270:		"n",
.mule/reports/log/20260907T191157Z-recipe-docs-sync-proposed.md:288:	if v, ok := os.LookupEnv("n"); ok {
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/config/config.go:n
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/config/config.go:n
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/config/config_test.go:n
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/config/config_test.go:n
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/container/block.go:n
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/container/docker.go:n
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/info.go:n
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/schema.go:n
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:$ rg -n 'n|n' /root/astrate-mule/internal/config/config.go
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:270:		"n",
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:288:	if v, ok := os.LookupEnv("n"); ok {
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:15:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `n` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `n` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:99:  `n` and
.mule/reports/log/20260912T192341Z-recipe-docs-sync-proposed.md:100:  `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260905T181053Z-recipe-milestones-proposed.md:Closed by 5111e3c. Decision: upstream wins on the wire, Astrate keeps its synchronous teardown. DELETE /housekeeping/v1/realms/{realm} now answers 405 {"errors":{"detail":"Realm deletion disabled"}} when the new opt-in flag ([housekeeping] realm_deletion_disabled / n) is set — flag check precedes existence, as measured live on the v1.2.0 tunnel where this stack has deletion off — and 422 {"errors":{"error_name":["connected_devices_present"]}} when devices are currently connected. Default (flag unset) keeps today's always-delete behavior for existing deployments; deletion remains synchronous/resumable-in-one-transaction rather than upstream's async job. Covered by TestHousekeepingDeleteGating; both gates mutation-verified.
.mule/reports/log/20260910T200120Z-recipe-github-issues-proposed.md:  `n` and
.mule/reports/log/20260910T200120Z-recipe-github-issues-proposed.md:  `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260924T182956Z-recipe-github-issues-proposed.md:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `n` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `n` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
.mule/reports/log/20260924T182956Z-recipe-github-issues-proposed.md:  `n` and
.mule/reports/log/20260924T182956Z-recipe-github-issues-proposed.md:  `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260915T201017Z-recipe-code-review-proposed.md:- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/config/config.go:n
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/config/config.go:n
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/config/config_test.go:n
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/config/config_test.go:n
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/container/block.go:n
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/container/docker.go:n
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/info.go:n
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/schema.go:n
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:$ rg -n 'n|FLOW_CONFIG' /root/astrate-mule/internal/flow/blocks/container/block.go | head
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:45:	Nested        map[string]any // opaque JSON for n
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/container/docker.go:19:	// FlowConfigJSON is passed as n.
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/container/docker.go:99:		args = append(args, "-e", "n="+spec.FlowConfigJSON)
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/container/docker.go:201:// encodeFlowConfigJSON marshals nested config for n.
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/schema.go:19:	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as n"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:/root/astrate-mule/internal/flow/blocks/info.go:116:		Config:  "image (required), config (object→n), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",
.mule/reports/log/20260909T200459Z-recipe-docs-sync-proposed.md:Noted but not proposed (different surfaces / already escalated): the `ASTRATE_HOUSEKEEPING_*` missing config keys are already in for-giulio.md; `n` is an env var Astrate *sets for spawned containers*, not a config key; realm-scoped appengine/pairing version twins belong to their own surfaces. No evidence file was required, no `gh issue create` asked for, so deliverable is just the task lines.
.mule/reports/log/20260924T193516Z-recipe-github-issues-proposed.md:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `n` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `n` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
.mule/reports/log/20260906T181352Z-recipe-github-issues-proposed.md:   `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md:internal/config/config.go:n
.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md:internal/config/config.go:n
.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md:internal/config/config_test.go:n
.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md:internal/config/config_test.go:n
.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/block.go:n
.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/docker.go:n
.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md:internal/flow/blocks/info.go:n
.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md:internal/flow/blocks/schema.go:n
.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md:270:		"n",
.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md:288:	if v, ok := os.LookupEnv("n"); ok {
.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md:+- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `n` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `n` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
.mule/reports/log/20260912T110449Z-recipe-docs-sync-proposed.md:- **Two config keys read but absent from configuration-reference.md** (`n`, `n`, config.go:268-293).
.mule/reports/log/20260907T174812Z-recipe-github-issues-proposed.md:  `n` and
.mule/reports/log/20260907T174812Z-recipe-github-issues-proposed.md:  `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/config/config.go:n
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/config/config.go:n
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/config/config_test.go:n
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/config/config_test.go:n
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/block.go:n
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/docker.go:n
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/flow/blocks/info.go:n
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/flow/blocks/schema.go:n
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:270:		"n",
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:288:	if v, ok := os.LookupEnv("n"); ok {
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:$ rg -n 'maxValueBytes|n' internal/appengine/*.go | grep -v _test
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:$ rg -n 'n' internal/flow/ | grep -v _test; echo "==="; ls .mule/
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/flow/blocks/info.go:116:		Config:  "image (required), config (object→n), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/flow/blocks/schema.go:19:	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as n"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/block.go:45:	Nested        map[string]any // opaque JSON for n
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/docker.go:19:	// FlowConfigJSON is passed as n.
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/docker.go:99:		args = append(args, "-e", "n="+spec.FlowConfigJSON)
.mule/reports/log/20260908T182108Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/docker.go:201:// encodeFlowConfigJSON marshals nested config for n.
.mule/reports/log/20260911T110803Z-recipe-docs-sync-proposed.md:  `n` and
.mule/reports/log/20260911T110803Z-recipe-docs-sync-proposed.md:  `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260911T110803Z-recipe-docs-sync-proposed.md:internal/config/config.go:n
.mule/reports/log/20260911T110803Z-recipe-docs-sync-proposed.md:internal/config/config.go:n
.mule/reports/log/20260913T180348Z-config-fail-loud-dev-mode-done.md: 	const env = "n"
.mule/reports/log/20260915T185804Z-recipe-github-issues-proposed.md:- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/reports/log/20260915T185804Z-recipe-github-issues-proposed.md:- **docs-sync appengine, 2026-09-12: two `ASTRATE_HOUSEKEEPING_*` env keys read by the code are absent from `docs/site/configuration-reference.md`.** The reference's `ASTRATE_*` inventory (20 keys) names neither `n` (and its bare upstream twin `HOUSEKEEPING_...`, read at internal/config/config.go:268-276, the realm default-retention override, #73; absent when unset — no default behaviour change) nor `n` (`os.LookupEnv`, config.go:288-293, the fail-loud realm-deletion gate, #75; absent/empty keeps the default). Both are live config paths with no documented counterpart. Configuration prose is yours — add both keys, or confirm the reference intentionally documents only a subset.
.mule/reports/log/20260915T185804Z-recipe-github-issues-proposed.md:  `n` and
.mule/reports/log/20260915T185804Z-recipe-github-issues-proposed.md:  `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260913T182321Z-recipe-docs-sync-proposed.md: - [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/reports/log/20260908T193321Z-recipe-github-issues-proposed.md:  `n` and
.mule/reports/log/20260908T193321Z-recipe-github-issues-proposed.md:  `n` are consumed in `internal/config/config.go`
.mule/reports/log/20260923T192627Z-flow-sort-bounded-buffer-done.md:	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as n"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:n
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:n
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:n
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:$ rg -n 'n|n|n|ASTRATE_TEST_DSN' internal/ | head -20
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:internal/flow/blocks/info.go:116:		Config:  "image (required), config (object→n), port (default 8080), timeout_ms (default 5000), ready_timeout_ms (default 15000)",
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:internal/flow/blocks/schema.go:19:	TypeContainer:         `{"type":"object","properties":{"image":{"type":"string","description":"Docker image run as a transform"},"config":{"type":"object","description":"Object passed to the container as n"},"port":{"type":"integer","description":"Container HTTP bridge port (default 8080)"},"timeout_ms":{"type":"integer","description":"Per-message POST timeout in milliseconds (default 5000)"},"ready_timeout_ms":{"type":"integer","description":"Health-check wait in milliseconds (default 15000)"}},"required":["image"]}`,
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/block.go:45:	Nested        map[string]any // opaque JSON for n
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/docker.go:19:	// FlowConfigJSON is passed as n.
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/docker.go:99:		args = append(args, "-e", "n="+spec.FlowConfigJSON)
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:internal/flow/blocks/container/docker.go:201:// encodeFlowConfigJSON marshals nested config for n.
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:internal/config/config.go:279:		"n",
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:internal/config/config.go:297:	if v, ok := os.LookupEnv("n"); ok {
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:internal/config/config_test.go:392:		prefixed = "n"
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:internal/config/config_test.go:533:	const env = "n"
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:$ rg -l 'n|n|n' docs/ || echo "absent from docs/"
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:- [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `n` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:+- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`n`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`n`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`n`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:- **Docs-sync run, 2026-09-21: three `ASTRATE_` config keys exist in the code but are absent from `docs/site/configuration-reference.md` — reverse drift.** `rg -o '\bASTRATE_[A-Z_]+' -N internal/` vs the reference page: **`n`** (read in `internal/config/config.go:277-292`, override injected as `housekeeping.default_datastream_maximum_storage_retention`, upstream-bare alias `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` wins when both set), **`n`** (fail-loud boolean gate, `config.go:294-306`) — these two are server startup keys and belong in `configuration-reference.md` (absent today), and **`n`** (`internal/flow/blocks/container/docker.go:99`, `docker.go:19`) — a per-`flow`-container-block env contract, likely a flow-blocks page concern rather than the config reference; your call where it lives. Reverse direction is clean: no documented key has been dropped from the code. Source is `internal/` (only `ASTRATE_TEST_DSN` is a test-only env var and correctly excluded; keys resolved via `os.LookupEnv` at `config.go:255-311`). Docs pages under `docs/site/` are on the never-edit list, so this is a decision/typing task, not a mule code change.
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:- `.mule/for-giulio.md` (newest-first): three `ASTRATE_` keys exist in code but are absent from `configuration-reference.md` (`n`, `n` — real startup keys; `n` — per-theme flow-block env contract, your call where it belongs). Reverse direction clean.
.mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md:Unsure: whether n belongs in configuration-reference or a flow page (flagged in the escalation)
.mule/reports/log/20260910T110956Z-recipe-github-issues-proposed.md:  `n` and
.mule/reports/log/20260910T110956Z-recipe-github-issues-proposed.md:  `n` are consumed in `internal/config/config.go`
---
279:		"ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
297:	if v, ok := os.LookupEnv("ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"); ok {

Previous runs already found and escalated this exact config-key drift. Let me check the current state of `for-giulio.md` and how a prior docs-sync run recorded its deliverable.
$ ls -t .mule/reports/log/ | rg 'docs-sync' | head -5
20260924T185047Z-recipe-docs-sync-proposed.md
20260924T165605Z-docs-sync-rm-auth-403-done.md
20260924T111049Z-docs-sync-rm-validation-example-prefix-done.md
20260924T110303Z-docs-sync-rm-interface-422-shapes-done.md
20260923T201746Z-docs-sync-rm-put-interface-409-done.md

→ Read .mule/for-giulio.md
Today already has a docs-sync run. Let me read its report to see which surface it covered, and check the prior docs-sync evidence-file format.
→ Read .mule/reports/log/20260924T185047Z-recipe-docs-sync-proposed.md
→ Read .mule/reports/log/20260921T203450Z-recipe-docs-sync-proposed.md [limit=40]
That run failed (API connect error — stub only). Let me look at a fully-executed docs-sync proposed report to see the expected evidence-file format, and check todo.md for pending docs-sync lines.
→ Read .mule/todo.md
Confirmed: pairing was previously swept and its `[auto]` fix tasks are all merged (`docs-sync-pairing-*` lines 116-117, 155-156). Final verification that the earlier pairing spec fixes are in the current yaml:

$ rg -n 'confirmed, pending|status:' docs/api/astarte_pairing_api.yaml | head; echo "---branch---"; git rev-parse --abbrev-ref HEAD
50:                  status: ok
251:                  status: confirmed
429:        status:
431:          enum: [confirmed, pending, inhibited]
497:        status:
---branch---
mule/queue

Pairing is fully in sync — the earlier `[auto]` fix tasks are merged and their fixes are in the current spec (register `404` at yaml:140-141, `status` enum `[confirmed, pending, inhibited]` at yaml:431). The config-key drift is already escalated twice in `for-giulio.md` (2026-09-12, 2026-09-21). No task lines to append.

Done: docs-sync pairing surface — accurate spec, no drift; 7/7 documented paths exist and status codes match handlers; config-key drift is already escalated (not re-raised).
Files: none changed (.mule/todo.md, .mule/for-giulio.md untouched — nothing to add; no git ops)
Verified: rg route/path diff (pairing yaml vs internal/pairing + cmd/astrate main.go:435) and handler status-check for health/version/register/unregister/credentials/info/verify -> all match. Note: the recipe's `HandleFunc` grep pattern misses `mux.Handle`-registered routes (http.go:82-85, main.go:435); a broader rg found them, so no false drift.
Unsure: nothing new — the two housekeeping config keys + ASTRATE_FLOW_CONFIG absent from configuration-reference.md are the known 09-12/09-21 escalation; ASTRATE_FLOW_CONFIG is set *by* Astrate into spawned containers, not a user config key. Unknown-device on device endpoints 404s-collapse-to-401 is deliberate (service.go:380-408), matching the spec.
Follow-ups: none
