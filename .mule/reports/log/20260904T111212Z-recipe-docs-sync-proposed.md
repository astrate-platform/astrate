slug: recipe-docs-sync
verdict: proposed
at:  b20aa3b
ran: 2026-09-04T11:12:12Z on DietPi in 359s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
Let me start by understanding what surfaces exist and picking one.
✱ Glob "docs/api/astrate_*_api.yaml" 1 match

$ ls /root/astrate-mule/docs/api/
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*realm*management*'
(no output)

Good. I'll pick **realm-management** as the surface for this run. Let me gather the code routes and the documented paths.
$ rg -n '^\s{2}/' /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
27:  /realmmanagement/v1/{realm}/interfaces:
98:  /realmmanagement/v1/{realm}/interfaces/{name}:
127:  /realmmanagement/v1/{realm}/interfaces/{name}/{major}:
239:  /realmmanagement/v1/{realm}/triggers:
309:  /realmmanagement/v1/{realm}/triggers/{name}:
361:  /realmmanagement/v1/{realm}/config/auth:
428:  /realmmanagement/v1/{realm}/config/device_registration_limit:
461:  /realmmanagement/v1/{realm}/version:
494:  /realmmanagement/v1/{realm}/devices/{device}:
520:  /realmmanagement/v1/{realm}/policies:
590:  /realmmanagement/v1/{realm}/policies/{name}:

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*realm*' 2>/dev/null | head -40
(no output)

$ rg -n 'realmmanagement|realm.management|realm_management' internal/ --glob '*.go' -l
internal/engine/router.go
internal/realm/http.go
internal/realm/dashboard_compat_test.go
internal/realm/service.go
internal/realm/http_test.go
internal/httpx/notfound_test.go
internal/observability/compat.go
internal/httpx/notfound.go
internal/observability/observability_test.go
internal/auth/middleware_test.go

$ rg -n 'HandleFunc|Handle\(|\.Get\(|\.Post\(|\.Put\(|\.Delete\(|\.Patch\(' internal/ --glob '*realm*' -l
(no output)

→ Read internal/realm/http.go
→ Read internal/engine/router.go
Now let me compare the code routes vs documented paths and read the YAML for detailed checks.

**Code routes** (19): listInterfaces, installInterface, listInterfaceMajors, getInterface, updateInterface, deleteInterface, listTriggers, createTrigger, getTrigger, deleteTrigger, getAuth, putAuth, getRegistrationLimit, getVersion, deleteDevice, listPolicies, createPolicy, getPolicy, deletePolicy.

**Documented paths** (11): interfaces (GET, POST), interfaces/{name} (GET), interfaces/{name}/{major} (GET, PUT, DELETE), triggers (GET, POST), triggers/{name} (GET, DELETE), config/auth (GET, PUT), config/device_registration_limit (GET), version (GET), devices/{device} (DELETE), policies (GET, POST), policies/{name} (GET, DELETE).

Interesting — code has no `PUT /config/device_registration_limit`. Let me read the YAML for detailed status code/response checks.
→ Read docs/api/astarte_realm_management_api.yaml
→ Read docs/api/astarte_realm_management_api.yaml [offset=90]

$ rg -n 'func WriteNotFound|func WriteDeviceNotFound|func WriteData|func WriteError|func WriteBadRequest|func WriteInternalServerError' internal/
(no output)

Let me check the astarteapi helper functions for response shapes and the device_registration_limit PUT handler.
→ Read internal/realm/service.go [limit=60]

$ rg -n 'func WriteNotFound|func WriteDeviceNotFound|func WriteData|func WriteError|func WriteBadRequest|func WriteInternalServerError' /root/astrate-mule/pkg/
/root/astrate-mule/pkg/astarteapi/envelope.go:121:func WriteData(w http.ResponseWriter, status int, v any) error {
/root/astrate-mule/pkg/astarteapi/envelope.go:140:func WriteDataWithLinks(w http.ResponseWriter, status int, v any, links Links) error {
/root/astrate-mule/pkg/astarteapi/envelope.go:147:func WriteError(w http.ResponseWriter, status int, detail string) error {
/root/astrate-mule/pkg/astarteapi/envelope.go:161:func WriteBadRequest(w http.ResponseWriter) error {
/root/astrate-mule/pkg/astarteapi/envelope.go:176:func WriteNotFound(w http.ResponseWriter) error {
/root/astrate-mule/pkg/astarteapi/envelope.go:181:func WriteDeviceNotFound(w http.ResponseWriter) error {
/root/astrate-mule/pkg/astarteapi/envelope.go:186:func WriteInternalServerError(w http.ResponseWriter) error {

$ rg -n 'DeviceRegistrationLimit|SetDeviceRegistrationLimit|RegistrationLimit' /root/astrate-mule/internal/realm/
/root/astrate-mule/internal/realm/http.go:47:	mux.Handle("GET /realmmanagement/v1/{realm}/config/device_registration_limit", h(a.getRegistrationLimit))
/root/astrate-mule/internal/realm/http.go:120:func (a *API) getRegistrationLimit(w http.ResponseWriter, r *http.Request) {
/root/astrate-mule/internal/realm/http.go:121:	limit, err := a.svc.GetDeviceRegistrationLimit(r.Context(), r.PathValue("realm"))
/root/astrate-mule/internal/realm/dashboard_compat_test.go:40:	t.Run("DeviceRegistrationLimit", func(t *testing.T) {
/root/astrate-mule/internal/realm/service.go:411:// GetDeviceRegistrationLimit returns the realm's device registration limit
/root/astrate-mule/internal/realm/service.go:413:func (s *Service) GetDeviceRegistrationLimit(ctx context.Context, realm string) (*int32, error) {
/root/astrate-mule/internal/realm/service.go:418:	return r.DeviceRegistrationLimit, nil

→ Read pkg/astarteapi/envelope.go [offset=170, limit=30]

$ rg -n 'func.*ListInterfaceMajors' /root/astrate-mule/internal/realm/service.go
196:func (s *Service) ListInterfaceMajors(ctx context.Context, realm, name string) ([]int, error) {

→ Read internal/realm/service.go [offset=196, limit=30]
→ Read pkg/astarteapi/envelope.go [limit=50, offset=115]

$ rg -n 'func.*DeleteTrigger|func.*GetTrigger|func.*ListTriggers|func.*CreateTrigger' /root/astrate-mule/internal/realm/service.go
251:func (s *Service) CreateTrigger(ctx context.Context, realm string, def []byte) (*store.Trigger, error) {
284:func (s *Service) GetTrigger(ctx context.Context, realm, name string) (json.RawMessage, error) {
297:func (s *Service) DeleteTrigger(ctx context.Context, realm, name string) error {
310:func (s *Service) ListTriggers(ctx context.Context, realm string) ([]string, error) {

Let me check if the service methods for triggers and policies actually return 404 errors (as the YAML documents).
$ rg -n 'func.*DeletePolicy|func.*GetPolicy|func.*ListPolicies|func.*CreatePolicy' /root/astrate-mule/internal/realm/service.go
339:func (s *Service) CreatePolicy(ctx context.Context, realm string, def []byte) (*store.TriggerPolicy, error) {
352:func (s *Service) GetPolicy(ctx context.Context, realm, name string) (json.RawMessage, error) {
365:func (s *Service) ListPolicies(ctx context.Context, realm string) ([]string, error) {
383:func (s *Service) DeletePolicy(ctx context.Context, realm, name string) error {

→ Read internal/realm/service.go [limit=60, offset=284]
→ Read internal/realm/service.go [limit=40, offset=352]

$ rg -n 'func.*GetTrigger\b|func.*DeleteTrigger\b|func.*GetTriggerPolicy\b|func.*DeleteTriggerPolicy\b' /root/astrate-mule/internal/store/
/root/astrate-mule/internal/store/policies.go:39:func (s *Store) GetTriggerPolicy(ctx context.Context, realmID int16, name string) (*TriggerPolicy, error) {
/root/astrate-mule/internal/store/policies.go:54:func (s *Store) DeleteTriggerPolicy(ctx context.Context, realmID int16, name string) error {
/root/astrate-mule/internal/store/triggers.go:38:func (s *Store) GetTrigger(ctx context.Context, realmID int16, name string) (*Trigger, error) {
/root/astrate-mule/internal/store/triggers.go:53:func (s *Store) DeleteTrigger(ctx context.Context, realmID int16, name string) error {

→ Read internal/store/triggers.go [limit=40, offset=38]
→ Read internal/store/policies.go [limit=40, offset=39]
→ Read internal/realm/service.go [limit=30, offset=383]
→ Read internal/realm/service.go [limit=40, offset=251]
→ Read internal/realm/service.go [limit=50, offset=62]

$ rg -o '\bASTRATE_[A-Z_]+' -N /root/astrate-mule/internal/ | sort -u
/root/astrate-mule/internal/config/config.example.toml:ASTRATE_DATABASE_DSN
/root/astrate-mule/internal/config/config.example.toml:ASTRATE_ENGINE_SHARDS
/root/astrate-mule/internal/config/config.example.toml:ASTRATE_HTTP_ADDR
/root/astrate-mule/internal/config/config.example.toml:ASTRATE_LOG_LEVEL
/root/astrate-mule/internal/config/config.example.toml:ASTRATE_MASTER_KEY
/root/astrate-mule/internal/config/config.example.toml:ASTRATE_MASTER_KEY_FILE
/root/astrate-mule/internal/config/config.example.toml:ASTRATE_MQTT_INSECURE_DEV_MODE
/root/astrate-mule/internal/config/config.go:ASTRATE_DATABASE_DSN
/root/astrate-mule/internal/config/config.go:ASTRATE_ENGINE_SHARDS
/root/astrate-mule/internal/config/config.go:ASTRATE_HTTP_ADDR
/root/astrate-mule/internal/config/config.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
/root/astrate-mule/internal/config/config.go:ASTRATE_HTTP_TLS_CERT_FILE
/root/astrate-mule/internal/config/config.go:ASTRATE_HTTP_TLS_KEY_FILE
/root/astrate-mule/internal/config/config.go:ASTRATE_LOG_FORMAT
/root/astrate-mule/internal/config/config.go:ASTRATE_LOG_LEVEL
/root/astrate-mule/internal/config/config.go:ASTRATE_MASTER_KEY
/root/astrate-mule/internal/config/config.go:ASTRATE_MASTER_KEY_FILE
/root/astrate-mule/internal/config/config.go:ASTRATE_MQTT_ADDR
/root/astrate-mule/internal/config/config.go:ASTRATE_MQTT_ADVERTISED_URL
/root/astrate-mule/internal/config/config.go:ASTRATE_MQTT_INSECURE_DEV_MODE
/root/astrate-mule/internal/config/config.go:ASTRATE_MQTT_SESSION_STORE_PATH
/root/astrate-mule/internal/config/config.go:ASTRATE_MQTT_TLS_CERT_FILE
/root/astrate-mule/internal/config/config.go:ASTRATE_MQTT_TLS_KEY_FILE
/root/astrate-mule/internal/config/config.go:ASTRATE_REALM_JWT_PUBLIC_KEY
/root/astrate-mule/internal/config/config.go:ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
/root/astrate-mule/internal/config/config.go:ASTRATE_REALM_NAME
/root/astrate-mule/internal/config/config.go:ASTRATE_SECURITY_MASTER_KEY_FILE
/root/astrate-mule/internal/config/config_test.go:ASTRATE_DATABASE_DSN
/root/astrate-mule/internal/config/config_test.go:ASTRATE_ENGINE_SHARDS
/root/astrate-mule/internal/config/config_test.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
/root/astrate-mule/internal/engine/engine_integration_test.go:ASTRATE_TEST_DSN
/root/astrate-mule/internal/store/crypto.go:ASTRATE_MASTER_KEY
/root/astrate-mule/internal/store/crypto.go:ASTRATE_MASTER_KEY_FILE
/root/astrate-mule/internal/store/store_test.go:ASTRATE_TEST_DSN
/root/astrate-mule/internal/testutil/pg.go:ASTRATE_TEST_DSN

Now let me check the `getRegistrationLimit` response shape more carefully, and look at `docs/site/configuration-reference.md` for config key drift.
$ ls /root/astrate-mule/docs/site/configuration-reference.md 2>/dev/null && echo EXISTS || echo NOT_FOUND
/root/astrate-mule/docs/site/configuration-reference.md
EXISTS

$ rg -o '\bASTRATE_[A-Z_]+' -N /root/astrate-mule/docs/site/configuration-reference.md | sort -u
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

$ rg -n 'INSECURE_DEV_MODE' /root/astrate-mule/internal/config/config.go
242:	if v, ok := os.LookupEnv("ASTRATE_MQTT_INSECURE_DEV_MODE"); ok {

$ rg -n 'MQTT_INSECURE_DEV_MODE' /root/astrate-mule/internal/config/config.example.toml
10:# ASTRATE_MQTT_INSECURE_DEV_MODE, ASTRATE_ENGINE_SHARDS, ASTRATE_LOG_LEVEL).

→ Read .mule/todo.md
Now let me read the existing `.mule/todo.md` and `.mule/for-giulio.md` to avoid duplicating lines.
→ Read .mule/for-giulio.md
I've completed the analysis. Let me compile the findings:

**Route-level drift: None.** All 19 code routes match the YAML's 11 path entries (with multiple methods per path). No "present but undocumented" or "documented but absent" endpoints.

**Status-code drift found:**

1. **DELETE `/policies/{name}` missing 422 response** — `service.go:398` returns `ErrValidation` when a policy is still referenced by a trigger, mapping to 422 Unprocessable Entity. The YAML only lists 204, 401, 404, 500.

**Config keys: All in sync** between `config.go` and `docs/site/configuration-reference.md`.

Now let me read `todo.md` (already done) and append the task line.

$ wc -l /root/astrate-mule/.mule/todo.md
128 /root/astrate-mule/.mule/todo.md

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -125,4 +125,5 @@
 - [!] object-path-ok-comment: in `internal/engine/data.go` `objectPathOK` (line 330), the loop over `ci.ObjectLeaves` returns after checking the first entry. Add a comment explaining why one leaf suffices (upstream validation guarantees uniform prefix depth), so future readers do not mistake the early return for a bug. [auto] — BLOCKED: gates failed
 - [!] housekeeping-unit-tests: in `internal/housekeeping`, add a unit test file (no integration build tag) that exercises `CreateRealm` validation paths — blank name returns ErrValidation, blank JWT key returns ErrValidation, negative registration limit returns ErrValidation — using a mock store and sealer. The package currently has only integration tests gated behind `//go:build integration`, so these paths are never tested on the Pi. [auto] — BLOCKED: gates failed
 - [!] store-downsample-lttb: in `internal/store`, implement the TODO at store.go:135 — when `s.hasToolkit` is true, switch `Downsample` from `time_bucket+avg` to toolkit `lttb()` downsampling. The probe already records availability; the time_bucket path in datastreams.go is the always-working default. Add a test that exercises both code paths (mocking `hasToolkit`). [auto] — BLOCKED: gates failed
 - [!] pairing-info-verify-ratelimit: in `internal/pairing/http.go`, add per-IP (and per-device) token-bucket rate limits to `handleInfo` (line 197) and `handleVerify` (line 238), which currently run a full bcrypt `CompareHashAndPassword` (via `authenticateDevice`, service.go:394) on any bearer string with no throttle — the sibling `handleCredentials` that does identical bcrypt work is limited at burst-5, and DESIGN §4.5 says "pairing endpoints" should be rate-limited. Reuse the credentials limiter defaults (1 req/s, burst 5) and extend `TestPairingHTTPRateLimit` to assert a burst of fake-secret requests to `GET /devices/{id}` and `POST .../verify` returns 429. [auto] — BLOCKED: gates failed
+- [ ] docs-sync-realm-policies-422: in docs/api/astarte_realm_management_api.yaml, add the undocumented 422 Unprocessable Entity response to DELETE /policies/{name} — the code (internal/realm/service.go:398) returns ErrValidation when a policy is still referenced by a trigger, which maps to 422 via writeError; the YAML only lists 204, 401, 404, 500. [auto]

Done. Here's the report:

**Surface examined:** realm-management (all 19 code routes vs all 11 YAML path entries)

**Drift found (1 task proposed):**

1. **docs/api/astarte_realm_management_api.yaml** — DELETE `/policies/{name}` is missing the 422 response. The code (`internal/realm/service.go:398`) returns `ErrValidation` when a policy is still referenced by a trigger, which maps to 422 Unprocessable Entity via `writeError` in `internal/realm/http.go:292-293`. The YAML only documents 204, 401, 404, 500 for this endpoint. Contrast with the analogous DELETE `/interfaces/{name}/{major}`, which correctly documents 422 for its "interface in use" case.

**Checked and clean:**
- All route paths and HTTP methods match between code and YAML.
- Config keys (`ASTRATE_*`) in `config.go` match `docs/site/configuration-reference.md` exactly (20 keys).
- Status codes for interfaces, triggers, config/auth, devices, and version endpoints all match.

Done: appended 1 task line to .mule/todo.md for the missing 422 on DELETE /policies/{name}
Files: .mule/todo.md
Verified: rg + read comparisons across internal/realm/http.go, service.go, docs/api/astarte_realm_management_api.yaml, internal/config/config.go, docs/site/configuration-reference.md
Unsure: nothing
Follow-ups: docs-sync-realm-policies-422
