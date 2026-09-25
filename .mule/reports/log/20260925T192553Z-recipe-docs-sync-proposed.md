slug: recipe-docs-sync
verdict: proposed
at:  3b06c46
ran: 2026-09-25T19:25:53Z on DietPi in 335s


> build · big-pickle

I'll start by reading the required files.
→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

$ ls /root/astrate-mule/.mule/ /root/astrate-mule/.mule/recipes/ /root/astrate-mule/docs/ /root/astrate-mule/docs/api/ 2>&1; echo "---SITE---"; ls /root/astrate-mule/docs/site/ 2>&1
/root/astrate-mule/.mule/:
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
recipes
reports
reviews
tasks
todo.md

/root/astrate-mule/.mule/recipes/:
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

/root/astrate-mule/docs/:
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

/root/astrate-mule/docs/api/:
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
---SITE---
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
housekeeping-api.md
index.md
interface-schema.md
json-payload-profile.md
migration-from-astarte.md
mqtt-protocol.md
observability.md
operations.md
pairing-and-security.md
payload-formats.md
pokemon-agent.md
quickstart.md
realm-management-api.md
swagger-ui
swagger.md
triggers.md
troubleshooting.md

$ wc -l /root/astrate-mule/.mule/todo.md /root/astrate-mule/.mule/for-giulio.md 2>&1
225 /root/astrate-mule/.mule/todo.md
   384 /root/astrate-mule/.mule/for-giulio.md
   609 total

→ Read .mule/todo.md
→ Read .mule/todo.md [offset=174]
Now let me pick a surface. Housekeeping has the least docs-sync coverage so far (4 tasks vs 12-16 for others), so I'll use that one.

$ cd /root/astrate-mule && echo "=== CODE ROUTES (housekeeping) ===" && rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/housekeeping/ && echo && echo "=== SPEC PATHS ===" && rg -n '^\s{2}/' docs/api/astarte_housekeeping_api.yaml
=== CODE ROUTES (housekeeping) ===

$ cd /root/astrate-mule && ls internal/ && echo "---" && ls internal/housekeeping/ 2>&1
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
async_operation_test.go
http.go
http_test.go
service.go
service_test.go

$ cd /root/astrate-mule && rg -n --no-heading 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc|mux\.' internal/housekeeping/http.go; echo "EXIT=$?"
37:	mux.Handle("GET /housekeeping/v1/realms", h(a.listRealms))
38:	mux.Handle("POST /housekeeping/v1/realms", h(a.createRealm))
39:	mux.Handle("GET /housekeeping/v1/realms/{realm}", h(a.getRealm))
40:	mux.Handle("PATCH /housekeeping/v1/realms/{realm}", h(a.patchRealm))
41:	mux.Handle("DELETE /housekeeping/v1/realms/{realm}", h(a.deleteRealm))
EXIT=0

$ cd /root/astrate-mule && echo "=== SPEC PATHS ===" && rg -n --no-heading '^  /' docs/api/astarte_housekeeping_api.yaml; echo; echo "=== SPEC OPS ===" && rg -n --no-heading '^    (get|post|put|patch|delete):' docs/api/astarte_housekeeping_api.yaml
=== SPEC PATHS ===
26:  /housekeeping/v1/realms:
106:  /housekeeping/v1/realms/{realm}:

=== SPEC OPS ===
27:    get:
50:    post:
107:    get:
140:    patch:
187:    delete:

All 5 routes are documented. Now checking status codes and response fields against the handler.
→ Read internal/housekeeping/http.go
→ Read docs/api/astarte_housekeeping_api.yaml

$ cd /root/astrate-mule && rg -n --no-heading 'func (m \*Middleware) RequireStatic' -A 40 internal/auth/middleware.go
(no output)

$ cd /root/astrate-mule && rg -n --no-heading 'func (s \*Service) (ListRealms|GetRealm|CreateRealm|UpdateRealm|DeleteRealm)' -A 30 internal/housekeeping/service.go
(no output)

$ cd /root/astrate-mule && rg -n --no-heading 'func \(s \*Service\)' internal/housekeeping/service.go
103:func (s *Service) WithDefaultDatastreamMaximumStorageRetention(defaultRetention *int64) *Service {
111:func (s *Service) WithRealmDeletionDisabled(disabled bool) *Service {
119:func (s *Service) notifyBrokerReload(ctx context.Context, realm string) {
136:func (s *Service) CreateRealm(ctx context.Context, name, jwtPublicKeyPEM string, regLimit *int32, retention *int64) (*RealmView, error) {
204:func (s *Service) UpdateRealm(ctx context.Context, name string, p RealmUpdate) error {
228:func (s *Service) GetRealm(ctx context.Context, name string) (*RealmView, error) {
238:func (s *Service) ListRealms(ctx context.Context) ([]string, error) {
256:func (s *Service) DeleteRealm(ctx context.Context, name string) error {

$ cd /root/astrate-mule && rg -n --no-heading 'RequireStatic' internal/auth/*.go
internal/auth/middleware.go:95:// RequireStatic guards an instance-level route (Housekeeping) with a fixed
internal/auth/middleware.go:98:func (m *Middleware) RequireStatic(claim Claim, keysPEM []string) func(http.Handler) http.Handler {
internal/auth/middleware_test.go:46:		m.RequireStatic(ClaimHousekeeping, hkKeys)(okHandler))

→ Read internal/auth/middleware.go [offset=40, limit=80]
→ Read internal/housekeeping/service.go [offset=130, limit=150]

$ cd /root/astrate-mule && echo "=== IN CODE ===" && rg -o --no-filename '\bASTRATE_[A-Z_]+' -N internal/ cmd/ pkg/ | sort -u > /tmp/opencode/code_keys.txt && cat /tmp/opencode/code_keys.txt && echo && echo "=== IN DOCS ===" && rg -o --no-filename '\bASTRATE_[A-Z_]+' -N docs/site/configuration-reference.md | sort -u > /tmp/opencode/doc_keys.txt && cat /tmp/opencode/doc_keys.txt && echo && echo "=== IN CODE NOT DOCS ===" && comm -23 /tmp/opencode/code_keys.txt /tmp/opencode/doc_keys.txt && echo "=== IN DOCS NOT CODE ===" && comm -13 /tmp/opencode/code_keys.txt /tmp/opencode/doc_keys.txt
=== IN CODE ===
ASTRATE_DATABASE_DSN
ASTRATE_ENGINE_SHARDS
ASTRATE_FLOW_CONFIG
ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
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
ASTRATE_TEST_DSN

=== IN DOCS ===
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

=== IN CODE NOT DOCS ===
ASTRATE_FLOW_CONFIG
ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
ASTRATE_TEST_DSN
=== IN DOCS NOT CODE ===

$ cd /root/astrate-mule && rg -n --no-heading -B4 -A8 'ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION|ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED' internal/config/config.go
275-	// are set. Absent vars leave the field nil: no injection, no behaviour
276-	// change for existing deployments.
277-	retention := ""
278-	for _, env := range []string{
279:		"ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
280-		"HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
281-	} {
282-		if v, ok := os.LookupEnv(env); ok && v != "" {
283-			retention = v
284-		}
285-	}
286-	if retention != "" {
287-		n, err := strconv.ParseInt(retention, 10, 64)
--
293-
294-	// The realm-deletion gate (#75) is a fail-loud boolean: only the exact
295-	// spellings below are accepted, anything else refuses to load rather
296-	// than silently meaning "off". Absent (or empty) leaves the default.
297:	if v, ok := os.LookupEnv("ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"); ok {
298-		switch v {
299-		case "1", "true", "TRUE", "True":
300-			cfg.Housekeeping.RealmDeletionDisabled = true
301-		case "0", "false", "":
302-			cfg.Housekeeping.RealmDeletionDisabled = false
303-		default:
304-			return fmt.Errorf("config: housekeeping.realm_deletion_disabled %q must be a boolean (1|true|TRUE|True|0|false)", v)
305-		}

$ cd /root/astrate-mule && rg -n --no-heading '^#{1,4} |ASTRATE_' docs/site/configuration-reference.md
1:# Configuration Reference
3:Astrate reads one TOML file (`astrate -config path.toml`) with `ASTRATE_*` environment overrides. Precedence: **built-in defaults < TOML file < environment variables**.
7:## Required values
11:- `database.dsn` (or `ASTRATE_DATABASE_DSN`) — PostgreSQL/TimescaleDB connection string.
16:## Precedence rule
19:built-in defaults  <  TOML file  <  ASTRATE_* environment
22:Every value below is the built-in default unless noted. Environment overrides exist for operationally critical fields, named `ASTRATE_<SECTION>_<FIELD>` (e.g. `ASTRATE_DATABASE_DSN`, `ASTRATE_HTTP_ADDR`).
26:## `[database]` — PostgreSQL/TimescaleDB
30:| `dsn` | string | `""` | `ASTRATE_DATABASE_DSN` | **Yes** | PostgreSQL + TimescaleDB connection string. |
32:## `[http]` — REST API listener
36:| `addr` | string | `":8080"` | `ASTRATE_HTTP_ADDR` | No | Bind address for the single REST listener (pairing, realm management, housekeeping, appengine, native endpoints). |
37:| `tls_cert_file` | string | `""` | `ASTRATE_HTTP_TLS_CERT_FILE` | No | Optional in-binary TLS cert. Leave empty to terminate TLS at a reverse proxy. |
38:| `tls_key_file` | string | `""` | `ASTRATE_HTTP_TLS_KEY_FILE` | No | Optional in-binary TLS key. Must be set together with `tls_cert_file`. |
39:| `cors_allowed_origins` | list[string] | `[]` | `ASTRATE_HTTP_CORS_ALLOWED_ORIGINS` (comma-separated) | No | Browser origins allowed for CORS. `"*"` allows any; empty disables CORS. The Astarte Dashboard SPA needs its origin here (e.g. `"http://localhost:4040"`). |
41:## `[mqtt]` — Embedded MQTT broker
45:| `addr` | string | `":8883"` | `ASTRATE_MQTT_ADDR` | No | Bind address for the mTLS MQTT listener. |
46:| `tls_cert_file` | string | `""` | `ASTRATE_MQTT_TLS_CERT_FILE` | **Yes*** | Broker server cert. Required unless `insecure_dev_mode` is true. |
47:| `tls_key_file` | string | `""` | `ASTRATE_MQTT_TLS_KEY_FILE` | **Yes*** | Broker server key. Required unless `insecure_dev_mode` is true. |
48:| `insecure_dev_mode` | bool | `true` | `ASTRATE_MQTT_INSECURE_DEV_MODE` | No | Bind a plaintext listener on `dev_addr` that trusts the claimed client ID. **Never enable in production.** |
50:| `session_store_path` | string | `"sessions.db"` | `ASTRATE_MQTT_SESSION_STORE_PATH` | No | bbolt/pebble file for MQTT session persistence. On a read-only container FS, must point at the writable session volume. |
51:| `advertised_url` | string | `""` | `ASTRATE_MQTT_ADVERTISED_URL` | No | Broker URL handed to devices by the pairing info endpoint. Empty derives `"mqtts://<addr>"`; set when devices reach the broker by another host. |
56:## `[engine]` — Ingestion pipeline
60:| `shards` | int | `16` | `ASTRATE_ENGINE_SHARDS` | No | Number of ordered ingestion shards. Each shard processes messages strictly in order. |
66:## `[pairing]` — Credential issuance & rate limits
78:## `[housekeeping]` — Instance-admin keys
85:## `[storage]` — Retention
91:## `[security]` — Master key
95:| `master_key_file` | string | `""` | `ASTRATE_SECURITY_MASTER_KEY_FILE` | No | File holding the AES-256 master key that seals realm CA private keys. When empty, falls back to `ASTRATE_MASTER_KEY` (64 hex chars) or `ASTRATE_MASTER_KEY_FILE`. |
100:## `[realm]` — Auto-provision realm
104:| `name` | string | `""` | `ASTRATE_REALM_NAME` | No | Realm name. Empty disables auto-provisioning. |
105:| `jwt_public_key` | string | `""` | `ASTRATE_REALM_JWT_PUBLIC_KEY` | Conditional | PEM public key for the auto-provisioned realm. Required when `name` is set. |
106:| `jwt_public_key_file` | string | `""` | `ASTRATE_REALM_JWT_PUBLIC_KEY_FILE` | Conditional | File path to the PEM public key. Alternative to inline `jwt_public_key`. |
109:## `[log]` — Logging
113:| `level` | string | `"info"` | `ASTRATE_LOG_LEVEL` | No | Log level: `debug`, `info`, `warn`, `error`. |
114:| `format` | string | `"json"` | `ASTRATE_LOG_FORMAT` | No | Log format: `json` (structured) or `text`. |
118:## Example: minimal production config
147:## Example: dev-mode docker compose
155:This sets `ASTRATE_MQTT_INSECURE_DEV_MODE=true`, `ASTRATE_MASTER_KEY` (throwaway), and auto-provisions a `test` realm. See [Deployment](deployment.md) for details.

→ Read docs/site/configuration-reference.md [offset=62, limit=55]

$ cd /root/astrate-mule && rg -n --no-heading -B12 'HousekeepingConfig|type Housekeeping' internal/config/config.go | head -50; echo "=== struct ==="; rg -n --no-heading -A14 'housekeeping struct|housekeepingConfig' internal/config/config.go | head -40
37-}
38-
39-// Std returns the underlying time.Duration.
40-func (d Duration) Std() time.Duration { return time.Duration(d) }
41-
42-// Config is the whole Astrate configuration.
43-type Config struct {
44-	HTTP         HTTPConfig         `toml:"http"`
45-	MQTT         MQTTConfig         `toml:"mqtt"`
46-	Database     DatabaseConfig     `toml:"database"`
47-	Engine       EngineConfig       `toml:"engine"`
48-	Pairing      PairingConfig      `toml:"pairing"`
49:	Housekeeping HousekeepingConfig `toml:"housekeeping"`
--
97-
98-// PairingConfig tunes credential issuance and its rate limits (§4.3, §4.5).
99-type PairingConfig struct {
100-	CertTTL           Duration `toml:"cert_ttl"`
101-	EnforceLatestCert bool     `toml:"enforce_latest_cert"`
102-	RegisterRate      float64  `toml:"register_rate"`
103-	RegisterBurst     int      `toml:"register_burst"`
104-	CredentialsRate   float64  `toml:"credentials_rate"`
105-	CredentialsBurst  int      `toml:"credentials_burst"`
106-	BcryptCost        int      `toml:"bcrypt_cost"`
107-}
108-
109:// HousekeepingConfig carries the instance-admin JWT public keys (a_ha): either
110-// inline PEM blocks or file references; both are concatenated.
111:type HousekeepingConfig struct {
=== struct ===

→ Read internal/config/config.go [limit=18, offset=109]
→ Read internal/config/config.go [limit=30, offset=270]

$ cd /root/astrate-mule && wc -l docs/site/housekeeping-api.md && rg -n --no-heading 'deletion_disabled|405|422|201|409|retention|device_registration_limit|Returns|returns' docs/site/housekeeping-api.md | head -60
72 docs/site/housekeeping-api.md
18:    "device_registration_limit": <optional integer>

→ Read docs/site/housekeeping-api.md

$ cd /root/astrate-mule && rg -n --no-heading -A45 'func \(s \*Store\) DeleteRealm|func .*\) DeleteRealm\(' internal/store/realms.go | head -60
211:func (s *Store) DeleteRealm(ctx context.Context, name string) error {
212-	tx, err := s.pool.Begin(ctx)
213-	if err != nil {
214-		return fmt.Errorf("store: beginning realm delete: %w", err)
215-	}
216-	defer func() { _ = tx.Rollback(ctx) }()
217-
218-	var realmID int16
219-	err = tx.QueryRow(ctx, `SELECT id FROM realms WHERE name = $1`, name).Scan(&realmID)
220-	if errors.Is(err, pgx.ErrNoRows) {
221-		return fmt.Errorf("%w: realm %q", ErrNotFound, name)
222-	}
223-	if err != nil {
224-		return fmt.Errorf("store: resolving realm %q: %w", name, err)
225-	}
226-
227-	for _, table := range []string{"individual_datastreams", "object_datastreams"} {
228-		if _, err := tx.Exec(ctx, `DELETE FROM `+table+` WHERE realm_id = $1`, realmID); err != nil {
229-			return fmt.Errorf("store: deleting %s of realm %q: %w", table, name, err)
230-		}
231-	}
232-	if _, err := tx.Exec(ctx, `DELETE FROM realms WHERE id = $1`, realmID); err != nil {
233-		return fmt.Errorf("store: deleting realm %q: %w", name, err)
234-	}
235-	if err := tx.Commit(ctx); err != nil {
236-		return fmt.Errorf("store: committing realm delete: %w", err)
237-	}
238-	return nil
239-}
240-
241-// scanRealm reads one realm row in realmColumns order.
242-func scanRealm(row pgx.Row) (*Realm, error) {
243-	var r Realm
244-	if err := row.Scan(&r.ID, &r.Name, &r.JWTPublicKeysPEM, &r.CACertificatePEM,
245-		&r.CAPrivateKeySealed, &r.DeviceRegistrationLimit,
246-		&r.DatastreamMaximumStorageRetention, &r.CreatedAt); err != nil {
247-		return nil, err
248-	}
249-	return &r, nil
250-}

$ cd /root/astrate-mule && rg -n --no-heading 'func Generate|func Import' pkg/ca/*.go internal/**/ca*.go 2>/dev/null; echo "--- callers of ca. ---"; rg -n --no-heading 'ca\.' internal/housekeeping/ internal/store/ cmd/ 2>/dev/null
--- callers of ca. ---
internal/housekeeping/service.go:161:	realmCA, err := ca.Generate(name, 0)

$ cd /root/astrate-mule && rg -n --no-heading -i 'REFERENCES realms' migrations/ | head -30
migrations/000012_user_blocks.up.sql:9:    realm_id      smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000009_flows.up.sql:6:    realm_id        smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000002_metadata.up.sql:20:    realm_id      smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000002_metadata.up.sql:52:    realm_id            smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000002_metadata.up.sql:78:    realm_id  smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000002_metadata.up.sql:93:    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000006_policies.up.sql:9:    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000008_pipelines.up.sql:8:    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,

$ cd /root/astrate-mule && rg -rn --files -g '*ca*' pkg internal | grep -i 'ca/' | head; echo "---"; rg -n --no-heading 'package ca' -l -g '*.go' . | head
internal/pairing/ca/ca.go
internal/pairing/ca/ca_test.go
---
./internal/pairing/ca/ca.go
./internal/pairing/ca/ca_test.go

$ cd /root/astrate-mule && rg -n --no-heading 'func (Generate|.*Import)' internal/pairing/ca/ca.go; echo "--- exported funcs ---"; rg -n --no-heading '^func [A-Z]' internal/pairing/ca/ca.go
71:func Generate(realm string, lifetime time.Duration) (*CA, error) {
--- exported funcs ---
71:func Generate(realm string, lifetime time.Duration) (*CA, error) {
121:func Load(certPEM string, keyDER []byte) (*CA, error) {
245:func ParseCertificatePEM(certPEM string) (*x509.Certificate, error) {

$ cd /root/astrate-mule && echo "=== tables with realm_id column ===" && rg -o --no-filename 'CREATE TABLE (IF NOT EXISTS )?(\w+)' -r '$2' migrations/*.up.sql | sort -u > /tmp/opencode/tables.txt && cat /tmp/opencode/tables.txt && echo && echo "=== which have FK to realms ===" && rg -n --no-heading -i 'realm_id\s+\w+.*(REFERENCES realms|references realms)' migrations/*.up.sql | head -40
=== tables with realm_id column ===
devices
endpoints
flows
group_devices
groups
individual_datastreams
interfaces
object_datastreams
pipelines
properties
realms
trigger_policies
triggers
user_blocks

=== which have FK to realms ===
migrations/000008_pipelines.up.sql:8:    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000009_flows.up.sql:6:    realm_id        smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000002_metadata.up.sql:20:    realm_id      smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000002_metadata.up.sql:52:    realm_id            smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000002_metadata.up.sql:78:    realm_id  smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000002_metadata.up.sql:93:    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000012_user_blocks.up.sql:9:    realm_id      smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000006_policies.up.sql:9:    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,

$ cd /root/astrate-mule && for t in devices endpoints interfaces properties groups triggers trigger_policies group_devices; do echo "### $t"; rg -n --no-heading -A25 "CREATE TABLE (IF NOT EXISTS )?$t\b" migrations/*.up.sql | rg -i 'realm|CREATE TABLE' | head -6; done
### devices
migrations/000002_metadata.up.sql:50:CREATE TABLE devices (
migrations/000002_metadata.up.sql-52-    realm_id            smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000002_metadata.up.sql-71-    PRIMARY KEY (realm_id, id)
### endpoints
migrations/000002_metadata.up.sql:33:CREATE TABLE endpoints (
migrations/000002_metadata.up.sql-50-CREATE TABLE devices (
migrations/000002_metadata.up.sql-52-    realm_id            smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
### interfaces
migrations/000002_metadata.up.sql:18:CREATE TABLE interfaces (
migrations/000002_metadata.up.sql-20-    realm_id      smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000002_metadata.up.sql-29-    UNIQUE (realm_id, name, major_version)
migrations/000002_metadata.up.sql-33-CREATE TABLE endpoints (
### properties
migrations/000003_properties.up.sql:3:CREATE TABLE properties (
migrations/000003_properties.up.sql-4-    realm_id     smallint NOT NULL,
migrations/000003_properties.up.sql-12-    PRIMARY KEY (realm_id, device_id, interface_id, path),
migrations/000003_properties.up.sql-13-    FOREIGN KEY (realm_id, device_id) REFERENCES devices(realm_id, id) ON DELETE CASCADE
### groups
migrations/000002_metadata.up.sql:76:CREATE TABLE groups (
migrations/000002_metadata.up.sql-78-    realm_id  smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000002_metadata.up.sql-80-    UNIQUE (realm_id, name)
migrations/000002_metadata.up.sql-82-CREATE TABLE group_devices (
migrations/000002_metadata.up.sql-84-    realm_id  smallint NOT NULL,
migrations/000002_metadata.up.sql-87-    FOREIGN KEY (realm_id, device_id) REFERENCES devices(realm_id, id) ON DELETE CASCADE
### triggers
migrations/000002_metadata.up.sql:91:CREATE TABLE triggers (
migrations/000002_metadata.up.sql-93-    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000002_metadata.up.sql-96-    UNIQUE (realm_id, name)
### trigger_policies
migrations/000006_policies.up.sql:7:CREATE TABLE trigger_policies (
migrations/000006_policies.up.sql-9-    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
migrations/000006_policies.up.sql-12-    UNIQUE (realm_id, name)
### group_devices
migrations/000002_metadata.up.sql:82:CREATE TABLE group_devices (
migrations/000002_metadata.up.sql-84-    realm_id  smallint NOT NULL,
migrations/000002_metadata.up.sql-87-    FOREIGN KEY (realm_id, device_id) REFERENCES devices(realm_id, id) ON DELETE CASCADE
migrations/000002_metadata.up.sql-90--- Triggers (Realm Management domain; executed by the engine)
migrations/000002_metadata.up.sql-91-CREATE TABLE triggers (
migrations/000002_metadata.up.sql-93-    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,

→ Read internal/housekeeping/service.go [offset=276]

$ cd /root/astrate-mule && rg -n --no-heading -B3 -A40 'func \(s \*Store\) UpdateRealm' internal/store/realms.go | head -55
154-// UpdateRealm applies the patched fields to one realm in a single UPDATE built
155-// from only the touched columns. An unknown realm yields ErrNotFound wrapped
156-// like the other realm methods.
157:func (s *Store) UpdateRealm(ctx context.Context, name string, p RealmPatch) error {
158-	sets := []string{}
159-	args := []any{name} // $1 is the WHERE name
160-	add := func(col string, v any) {
161-		sets = append(sets, fmt.Sprintf("%s = $%d", col, len(args)+1))
162-		args = append(args, v)
163-	}
164-	if p.PatchJWTPublicKeyPEM {
165-		add("jwt_public_keys", []string{p.SetJWTPublicKeyPEM})
166-	}
167-	if p.ClearRegistrationLimit {
168-		add("device_registration_limit", nil)
169-	} else if p.PatchRegistrationLimit {
170-		add("device_registration_limit", p.SetRegistrationLimit)
171-	}
172-	if p.ClearRetention || (p.PatchRetention && p.SetRetention <= 0) {
173-		add("datastream_maximum_storage_retention", nil)
174-	} else if p.PatchRetention {
175-		add("datastream_maximum_storage_retention", p.SetRetention)
176-	}
177-	if len(sets) == 0 {
178-		return nil
179-	}
180-
181-	tag, err := s.pool.Exec(ctx,
182-		`UPDATE realms SET `+strings.Join(sets, ", ")+` WHERE name = $1`, args...)
183-	if err != nil {
184-		return fmt.Errorf("store: updating realm %q: %w", name, err)
185-	}
186-	if tag.RowsAffected() == 0 {
187-		return fmt.Errorf("%w: realm %q", ErrNotFound, name)
188-	}
189-	return nil
190-}
191-
192-// SetRealmCA replaces the realm's CA certificate and sealed private key
193-// (re-keying flow, docs/DESIGN.md §4.3).
194-func (s *Store) SetRealmCA(ctx context.Context, name, caCertPEM string, caKeySealed []byte) error {
195-	tag, err := s.pool.Exec(ctx,
196-		`UPDATE realms SET ca_certificate = $2, ca_private_key = $3 WHERE name = $1`,
197-		name, caCertPEM, caKeySealed)

$ cd /root/astrate-mule && cat docs/Makefile 2>/dev/null | head -40; echo "=== site/api ==="; ls docs/site/api/ 2>&1
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
=== site/api ===
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

$ cd /root/astrate-mule && rg -n --no-heading 'Content-Type|charset' pkg/astarteapi/envelope.go | head -20
21:// ContentType is the Content-Type header value set on every envelope,
23:const ContentType = "application/json; charset=utf-8"
112:		w.Header().Set("Content-Type", ContentType)
117:	w.Header().Set("Content-Type", ContentType)

$ cd /root/astrate-mule && rg -n --no-heading -A30 'func .*Authorizes' internal/auth/*.go pkg/auth/*.go 2>/dev/null | head -45
internal/auth/claims.go:164:func (t *Token) Authorizes(claim Claim, verb, authPath string) bool {
internal/auth/claims.go-165-	for _, g := range t.grants[claim] {
internal/auth/claims.go-166-		if g.verb.MatchString(verb) && g.path.MatchString(authPath) {
internal/auth/claims.go-167-			return true
internal/auth/claims.go-168-		}
internal/auth/claims.go-169-	}
internal/auth/claims.go-170-	return false
internal/auth/claims.go-171-}
internal/auth/claims.go-172-
internal/auth/claims.go-173-// Channels authorization verbs. Upstream recognises exactly these two and
internal/auth/claims.go-174-// silently discards an a_ch entry carrying anything else.
internal/auth/claims.go-175-const (
internal/auth/claims.go-176-	VerbJoin  = "JOIN"
internal/auth/claims.go-177-	VerbWatch = "WATCH"
internal/auth/claims.go-178-)
internal/auth/claims.go-179-
internal/auth/claims.go-180-// AuthorizesChannel reports whether the a_ch claim grants verb (VerbJoin or
internal/auth/claims.go-181-// VerbWatch) on authPath — the room name for a join, the trigger's target for
internal/auth/claims.go-182-// a watch.
internal/auth/claims.go-183-//
internal/auth/claims.go-184-// It deliberately does not reuse Authorizes, because upstream reads a_ch by a
internal/auth/claims.go-185-// different rule than the REST claims and the difference is observable. The
internal/auth/claims.go-186-// REST plug compiles the verb field into a regex and matches it against the
internal/auth/claims.go-187-// HTTP method, so "GET|POST" works and ".*::.*" grants everything. The
internal/auth/claims.go-188-// Channels socket instead *partitions* the a_ch list by an exact match on the
internal/auth/claims.go-189-// verb field, keeping only entries whose first field is literally "JOIN" or
internal/auth/claims.go-190-// "WATCH" and discarding the rest, then matches the path regex within the
internal/auth/claims.go-191-// chosen bucket.
internal/auth/claims.go-192-//
internal/auth/claims.go-193-// So a blanket ".*::.*" — which authorizes every REST surface — authorizes
internal/auth/claims.go-194-// nothing here: ".*" is not the string "JOIN". This is measured, not inferred;
--
internal/auth/claims.go:197:func (t *Token) AuthorizesChannel(verb, authPath string) bool {
internal/auth/claims.go-198-	for _, g := range t.grants[ClaimChannels] {
internal/auth/claims.go-199-		if g.verbLiteral == verb && g.path.MatchString(authPath) {
internal/auth/claims.go-200-			return true
internal/auth/claims.go-201-		}
internal/auth/claims.go-202-	}
internal/auth/claims.go-203-	return false
internal/auth/claims.go-204-}
--
internal/auth/claims_test.go:159:func TestAuthorizesChannelParity(t *testing.T) {
internal/auth/claims_test.go-160-	const room = "probe"
internal/auth/claims_test.go-161-
internal/auth/claims_test.go-162-	cases := []struct {

$ cd /root/astrate-mule && rg -n --no-heading '^func Test' internal/housekeeping/async_operation_test.go internal/housekeeping/http_test.go | head -20; echo "=== tags ==="; head -3 internal/housekeeping/async_operation_test.go
internal/housekeeping/async_operation_test.go:16:func TestHousekeepingAsyncOperationParam(t *testing.T) {
internal/housekeeping/http_test.go:116:func TestHousekeeping(t *testing.T) {
internal/housekeeping/http_test.go:197:func TestHousekeepingPatchRealm(t *testing.T) {
internal/housekeeping/http_test.go:309:func TestHousekeepingRetentionDefault(t *testing.T) {
internal/housekeeping/http_test.go:360:func TestHousekeepingDeleteGating(t *testing.T) {
=== tags ===
//go:build integration

package housekeeping

$ cd /root/astrate-mule && rg -n --no-heading -i 'async' internal/housekeeping/http.go internal/housekeeping/service.go cmd/astrate/main.go | head -20; echo "=== test ==="; sed -n '1,60p' internal/housekeeping/async_operation_test.go
=== test ===
//go:build integration

package housekeeping

import (
	"net/http"
	"testing"
)

// TestHousekeepingAsyncOperationParam pins deviation 17
// (docs/COMPATIBILITY.md): upstream 1.4 runs realm create and delete
// asynchronously unless the caller opts out with `?async_operation=false`.
// Astrate is always synchronous, so the parameter is accepted and ignored on
// both values — a client that sends it must never see a 4xx for it, and the
// response must be identical to the parameterless call.
func TestHousekeepingAsyncOperationParam(t *testing.T) {
	r := newHKRig(t)

	for _, q := range []string{"?async_operation=false", "?async_operation=true"} {
		realmName := "as" + randSuffix(t)
		createBody := `{"realm_name":` + jsonStr(realmName) + `,"jwt_public_key_pem":` + jsonStr(r.realmKey) + `}`

		if rec := r.req(t, http.MethodPost, "/realms"+q, createBody, r.haToken); rec.Code != http.StatusCreated {
			t.Fatalf("create realm %s: got %d, want 201 (%s)", q, rec.Code, rec.Body)
		}
		// Synchronous means done on return: the realm reads back immediately.
		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusOK {
			t.Errorf("get realm after create %s: got %d, want 200", q, rec.Code)
		}
		if rec := r.req(t, http.MethodDelete, "/realms/"+realmName+q, "", r.haToken); rec.Code != http.StatusNoContent {
			t.Fatalf("delete realm %s: got %d, want 204 (%s)", q, rec.Code, rec.Body)
		}
		if rec := r.req(t, http.MethodGet, "/realms/"+realmName, "", r.haToken); rec.Code != http.StatusNotFound {
			t.Errorf("get realm after delete %s: got %d, want 404", q, rec.Code)
		}
	}
}

$ cd /root/astrate-mule && rg -n --no-heading -B3 -A8 'async_operation' docs/COMPATIBILITY.md | head -40
336-    upstream 1.4 runs realm create/delete (Housekeeping), interface
337-    install/update/delete (Realm Management) and trigger-delivery-policy
338-    delete in the background, answering immediately and letting the caller opt
339:    into synchronous execution with `?async_operation=false`. Astrate performs
340-    every one of these synchronously and answers only once the work is done —
341-    a strictly stronger guarantee, so a client written against upstream never
342-    observes a regression. The parameter is accepted and ignored on either
343-    value, on every one of those endpoints, so upstream clients that send it
344-    keep working unchanged; pinned by
345-    `TestRealmManagementAsyncOperationParam` and
346-    `TestHousekeepingAsyncOperationParam`. Tracked in
347-    `docs/UPSTREAM-EXPERIMENTAL.md` (#68) until upstream 1.4 is final.

$ cd /root/astrate-mule && rg -n --no-heading 'ca\.Load|ca\.ParseCertificatePEM|ca\.Generate' --glob '*.go' . | grep -v '_test.go'
./internal/housekeeping/service.go:161:	realmCA, err := ca.Generate(name, 0)
./internal/pairing/service.go:157:	realmCA, err := ca.Generate(realmName, lifetime)
./internal/pairing/service.go:350:		cert, perr := ca.ParseCertificatePEM(clientCrtPEM)
./internal/pairing/service.go:419:	realmCA, err := ca.Load(realm.CACertificatePEM, keyDER)
./internal/broker/authhook.go:90:		cert, err := ca.ParseCertificatePEM(r.CACertificatePEM)

$ cd /root/astrate-mule && tail -60 .mule/for-giulio.md
- **Group-scoped triggers (`group_name` on device/data triggers) compile but never match**
  (`internal/engine/triggers/match.go:11-12`). Decision deferred, tied to issue #17
  (group-WATCH-path reconciliation, trickle work, not mule): whatever group-membership
  mechanism comes out of that phase should also report the perf cost for this decision —
  noted in a comment on #17 so it isn't benchmarked twice. (Same survey, source 4.)

---

- ~~The Pi cannot run the race detector~~ — **resolved 2026-07-27** by installing Go 1.26.5
  as a userland toolchain on the Legion Go (`~/.local/go`, no root, `rm -rf` to undo). The
  Pi still cannot run `-race` (39-bit VMA kernel vs the 48 ThreadSanitizer needs), so its
  gate remains `go vet ./... && go test ./...` — but race coverage now exists on the Legion
  Go, where the full suite runs clean in ~40s on 16 cores. The standing `race-check` task is
  the concurrency gate. Concurrency work is queueable again, provided the race-check runs
  after it.
- ~~golangci-lint is not installed on the Pi~~ — **resolved 2026-07-28**: installed v2.12.2
  via `go install` (the prebuilt-binary installer's published sha256 for linux/arm64 did not
  verify, reproducibly, so built from source instead). `mule.service` now sets
  `MULE_LINT_CMD=golangci-lint run ./...` and has `/root/go/bin` on `PATH`; the lint gate
  runs starting with the next tick.
- ~~`/root/astrate` on the Pi has uncommitted work~~ — **resolved 2026-07-27** with the new
  `tools/reconcile.sh`: rescued onto `origin/wip/DietPi-20260727T171543Z` (pushed, not
  reviewed — read the diff before merging anything from it) and `/root/astrate` is now clean
  on `main`. Also had to set `commit.gpgsign false` locally on that clone first — it had
  signing on with no working gpg-agent for a non-interactive session, same class of problem
  `/root/astrate-mule` already solved the same way (your call, confirmed 2026-07-27: align it
  with the existing unattended-Pi-clone convention rather than fix gpg-agent or bypass signing
  silently inside the script).

## 2026-08-23 — FDO promoted to milestone-4.0 candidate (Giulio's decision, recorded)

#78 is no longer parked: zero-touch onboarding is strategic for commercial
viability. Scope frozen on the issue: owner-side TO1/TO2 in our Pairing
service only (last mile, like upstream), reuse fdo-rs for
manufacturing/rendezvous, acceptance = official `astarte-device-fdo-rust`
SDK completes onboarding against Astrate, docs as a first-class deliverable.
When v3.0 is marked DONE, the v4.0 section of `.mule/milestones.md` should be
drafted with this investigation as its first item (issue #78 has the full
verified context).
- **The mule has been idle 17h.** Filed by the dead-man's switch; see journalctl on the Pi.

## 2026-09-25 — the `-tags nats` half of the tree is compiled by no gate

`cmd/astrate/newnats_nats.go`, `internal/engine/forward/nats.go` and
`internal/engine/forward/nats_test.go` are all behind `//go:build nats`, and nothing
builds or tests that tag: not `make build`/`make test`/`make test-integration`/
`make test-e2e` (Makefile:23, 38, 42, 46), not the mule gate, not the lint target.
So the whole NATS trigger-forwarding path is invisible to CI — while
`config.validate` accepts `triggers.forward.kind = "nats"` (internal/config/config.go:339-343)
and the untagged binary then fails at boot with the build-tag message
(cmd/astrate/newnats_default.go). I measured it today: `go build -tags nats ./cmd/...`
and `go vet -tags nats ./cmd/...` both exit 0, so there is nothing broken — just nothing
keeping it unbroken. Worth a `make build-nats` (and the NATS test tier, which needs a
container, so it would be a `[legion]` line) in the gate; not a queue line because the
only test I can imagine for it shells out to `go build`, and the mule's own check strips
the implementation and the guard together, so it would never fail.
Not queued for the same reason: `IdleTimeout` is missing on the HTTP server
(cmd/astrate/main.go:221) and `-healthcheck` ignores the config file
(cmd/astrate/main.go:531) — both real, both yours to call, neither has a test worth
writing. See `.mule/reviews/cmd-astrate-2026-09-25.md`.

→ Read .mule/for-giulio.md [limit=30]
Let me verify the docs build currently works, so my proposals name a verifiable command.

$ cd /root/astrate-mule && (command -v mkdocs >/dev/null && timeout 180 make -C docs build 2>&1 | tail -15) || echo "mkdocs NOT INSTALLED on this box"
mkdocs NOT INSTALLED on this box

$ cd /root/astrate-mule && python3 -c "
import yaml,sys
for f in ['astarte_housekeeping_api.yaml']:
    p='docs/api/'+f
    d=yaml.safe_load(open(p))
    print(p,'OK  paths:',list(d['paths'].keys()))
" 2>&1
Traceback (most recent call last):
  File "<string>", line 2, in <module>
    import yaml,sys
ModuleNotFoundError: No module named 'yaml'

$ cd /root/astrate-mule && cat > /tmp/opencode/yamlcheck.go <<'EOF'
package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	for _, p := range os.Args[1:] {
		b, err := os.ReadFile(p)
		if err != nil {
			fmt.Println(p, "READ ERR", err)
			os.Exit(1)
		}
		var d map[string]any
		if err := yaml.Unmarshal(b, &d); err != nil {
			fmt.Println(p, "PARSE ERR", err)
			os.Exit(1)
		}
		paths, _ := d["paths"].(map[string]any)
		fmt.Println(p, "OK,", len(paths), "paths")
	}
}
EOF
mkdir -p /tmp/opencode/yc && cp /tmp/opencode/yamlcheck.go /tmp/opencode/yc/main.go && go run /tmp/opencode/yc/main.go docs/api/astarte_housekeeping_api.yaml 2>&1 | tail -5
docs/api/astarte_housekeeping_api.yaml OK, 2 paths

Verified the spec parses. Now appending the two `docs/api` proposals to the queue.
→ Read .mule/todo.md [offset=222]
✗ Edit .mule/todo.md failed
Error: Could not find oldString in the file. It must match exactly, including whitespace, indentation, and line endings.

$ cd /root/astrate-mule && tail -c 120 .mule/todo.md | od -c | tail -5
0000100   K   E   Y   _   F   I   L   E       f   i   r   s   t       s
0000120   o       t   h   e       o   s   .   S   e   t   e   n   v    
0000140   i   s       r   e   s   t   o   r   e   d       a   t       c
0000160   l   e   a   n   u   p   .  \n
0000170

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -222,4 +222,6 @@
 - [x] drain-per-stage-budget [auto]: `shutdown` spends a single 30s sctx (cmd/astrate/main.go:479) across three sequential stages — srv.Shutdown (484), b.Close (489), flowSvc.Manager().Shutdown (494) and MarkRunningFlowsStopped (497) — while `drainEngine` alone opens a second one (504), so the constant's own comment "bounds the whole graceful drain" (main.go:55-56) is wrong in both directions, and one slow in-flight HTTP request (the server sets no per-request deadline, 221) burns the whole first budget and leaves the flow stage on an expired context: StopFlow then returns at `f.router.Drain(ctx)` before step 3 releases block resources (internal/flow/flow.go:265-273) and MarkRunningFlowsStopped skips every flow (flowapi/service.go:895-897), both surfacing as one log.Warn. Give each stage its own budget through an extracted helper (shutdownTimeout has to become an injectable var for the test) and pin it container-free: a second stage still receives a live context after a first stage exhausts its own.
 - [x] cmd-devcert-fields-test [auto]: add a container-free test for `selfSignedDevCert` (cmd/astrate/main.go:320-345), the throwaway mTLS identity `insecure_dev_mode` boots on, which no test touches — pin NotBefore = now-1h (331; a plain now would reject a device whose clock is behind), NotAfter = now+365d (332), DNSNames localhost and IPAddresses 127.0.0.1/::1 (336-337, dropping which breaks hostname verification for a device dialling mqtts://127.0.0.1:8883), ExtKeyUsage serverAuth (334), and that the leaf is self-signed and parses (x509.ParseCertificate on Certificate[0]).
 - [!] cmd-healthcheck-contract-test [auto]: add a container-free test for `runHealthcheck` (cmd/astrate/main.go:531-558) — the `HEALTHCHECK CMD ["/astrate", "-healthcheck"]` contract (Dockerfile:35-36) has no coverage. Assert it probes /astrate/v1/readiness and not /health (543), returns 0 only on 200 so a 503 readiness exits 1 (554), and maps a bare ":PORT" in ASTRATE_HTTP_ADDR to 127.0.0.1 (536-538). httptest.NewServer plus t.Setenv("ASTRATE_HTTP_ADDR", srv.Listener.Addr().String()). — BLOCKED: tests failed: --- FAIL: TestRunHealthcheckProbesReadiness (0.01s)
 - [x] cmd-loadsealer-masterkeyfile-test [auto]: add a container-free test for `loadSealer` (cmd/astrate/main.go:515-526) — the `security.master_key_file` to env hand-off (`os.Setenv(store.EnvMasterKeyFile, ...)`, 517) is the only wiring that makes a configured key file reach the store loader, and dropping or misspelling that branch leaves the file silently ignored and boot dying with the generic "master key" error (523) even though the config named a valid file. Assert a 32-byte key file under t.TempDir() loads, the inline ASTRATE_MASTER_KEY path still wins, and a missing key returns the error naming all three env vars. t.Setenv both ASTRATE_MASTER_KEY and ASTRATE_MASTER_KEY_FILE first so the os.Setenv is restored at cleanup.
+- [ ] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /housekeeping/v1/realms` and `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — Astrate accepts and ignores it on either value (deviation 17, docs/COMPATIBILITY.md:336-347; neither handler reads the query string, internal/housekeeping/http.go:62-74 and 85-91, and TestHousekeepingAsyncOperationParam in internal/housekeeping/async_operation_test.go:16 pins that neither value 4xxs and that the response is identical to the parameterless call), but the `post` operation declares no `parameters` block at all (yaml:50-104) and `delete` carries only `RealmName` (yaml:195), so a client generated from this spec never learns the parameter exists. Add an `AsyncOperation` component parameter (`in: query`, boolean, default `false`) and `$ref` it from both operations, with a description stating the value is accepted and ignored because Astrate performs create and delete synchronously (deviation 17). The RM twins (interface install/update/delete, policy delete) have the same gap in astarte_realm_management_api.yaml but are out of this spec's scope — leave them. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_housekeeping_api.yaml — the descriptions say null is the only way to clear `datastream_maximum_storage_retention` ("A null `device_registration_limit` or `datastream_maximum_storage_retention` clears that field", yaml:146-147; "Null clears the retention", yaml:309) and give the field no create-time caveat (yaml:291), but the wire folds an explicit 0 to unset on both paths: PATCH maps `null || val == 0` to `ClearRetention` (internal/housekeeping/http.go:216-219) and the store repeats the rule independently (`SetRetention <= 0` → nil, internal/store/realms.go:172-173), while create folds 0 to nil before injecting the configured default (internal/housekeeping/service.go:149-159, upstream parity measured on v1.2.0). A client following the spec today sends `0` expecting a literal zero-second retention and silently gets unlimited. State the fold in the `patchRealm` description, the `RealmPatch` and `RealmCreate` field descriptions, and record the deliberate asymmetry: `device_registration_limit` has no such fold, so `0` there is stored literally (http.go:208-215) and only `null` clears it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

Now the escalations to `.mule/for-giulio.md` (newest at the top, per that file's header):
✗ Edit .mule/for-giulio.md failed
Error: Found multiple matches for oldString. Provide more surrounding context to make the match unique.

← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,72 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+## 2026-09-25 — docs-sync run, surface: **housekeeping** (the least-covered of the five)
+
+Nothing filed, no `gh issue create` in this recipe. Two `docs/api/` proposals appended to
+`.mule/todo.md` (`docs-sync-hk-async-operation-param`, `docs-sync-hk-retention-zero-is-unset`);
+the three items below are escalations, because they are your prose.
+
+**The path list is clean.** All five routes in code (`internal/housekeeping/http.go:37-41`) are
+documented and the documented status codes now match the handler on all five ops — the four
+earlier housekeeping tasks (PATCH route, retention field, DELETE 405/422, the 422 example)
+closed the real gaps. What is left is description-level, not path-level.
+
+- **`docs/site/configuration-reference.md` `[housekeeping]` lists 2 of the 4 keys it has, and
+  shows `—` in the env-override column for both it does list.** The table (lines 78-83) carries
+  only `jwt_public_keys` and `jwt_public_key_files`. `HousekeepingConfig`
+  (`internal/config/config.go:111-124`) also has
+  `default_datastream_maximum_storage_retention` (`*int64`, #73) and `realm_deletion_disabled`
+  (`bool`, #75), and **both have working env overrides**:
+  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (config.go:277-292, which
+  also accepts the bare upstream name `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`
+  and lets it win when both are set) and `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`
+  (config.go:297-305, fail-loud: only `1|true|TRUE|True|0|false|""` load, anything else refuses
+  to boot). Neither appears anywhere in the file. The keys that gate a destructive operation
+  (`realm_deletion_disabled` makes `DELETE /housekeeping/v1/realms/{realm}` answer 405) and the
+  one that sets a fleet-wide retention ceiling are both invisible to an operator reading the
+  reference. Cross-checked with the recipe's sweep — `rg -o '\bASTRATE_[A-Z_]+' -N internal/ |
+  sort -u` vs the file: the only other misses are `ASTRATE_FLOW_CONFIG` and the test-only
+  `ASTRATE_TEST_DSN`; no key is documented that the code does not have.
+- **`docs/site/housekeeping-api.md` has no PATCH section.** The page enumerates Create / List /
+  Get / Delete (lines 10-50) and stops, while the route exists (`internal/housekeeping/http.go:40`,
+  handler `patchRealm` 163-230) and is fully specified in the OpenAPI file
+  (`docs/api/astarte_housekeeping_api.yaml:140-185`). Related, in the same page: the create
+  body example (line 18) lists only `realm_name` / `jwt_public_key_pem` /
+  `device_registration_limit` and omits `datastream_maximum_storage_retention`, which is on the
+  wire (`realmBody`, http.go:50) and in `RealmCreate` (yaml:287-291). The two nullable integer
+  fields are the ones operators most often need to set.
+- **"Generate (or import) the per-realm CA key + certificate" (`housekeeping-api.md:25`) — there
+  is no import path.** `CreateRealm` calls `ca.Generate(name, 0)` unconditionally
+  (`internal/housekeeping/service.go:161`); the `ca` package's only other constructor is
+  `ca.Load` (`internal/pairing/ca/ca.go:121`), whose sole caller is pairing loading the realm
+  CA that housekeeping already stored (`internal/pairing/service.go:419`) — nothing anywhere
+  imports an existing realm CA at create time. If "import" is a roadmap item rather than a
+  description of today, say so; if it is aspirational, the page reads as a capability a client
+  can use. Smaller, same paragraph: "in a single transaction" lists generate-then-seal before
+  the row insert, but the code generates and seals the CA *first* (service.go:161-168) and only
+  then calls `store.CreateRealm` (170) — the DB work is one transaction, the CA work is not
+  part of it.
+
+Non-findings worth recording so a later run does not re-chase them: the realm cascade really is
+complete (the page's "removes all interfaces, devices, properties, datastreams, groups, and
+triggers" holds — every table reaches `realms` by `ON DELETE CASCADE` directly or through
+`devices`, and `individual_datastreams`/`object_datastreams` are deleted explicitly first,
+`internal/store/realms.go:227-234`); `view` returning only `JWTPublicKeysPEM[0]`
+(service.go:277-281) loses nothing, because `store.CreateRealm` stores a one-element slice and
+`store.UpdateRealm` *replaces* the array rather than appending (realms.go:164-166); `403` is
+genuinely reachable on every housekeeping op (`RequireStatic` → `authorize` writes
+`WriteForbidden` when a valid `a_ha` grant does not match, `internal/auth/middleware.go:110-117`);
+and the PATCH op correctly documents **no** 405 or connected-devices 422, since `UpdateRealm`
+has no deletion-gating path. Spec parses clean today
+(`yaml.Unmarshal`, 2 paths) — `mkdocs` is not installed on this box, so `make -C docs build`
+could not be run here; both proposed tasks carry it as their own verification step, as every
+earlier docs-sync line does.
+
+---
+
 - **github-issues triage run, 2026-09-25: still nothing proposable — the mule-alarm pile is 18 issues (#94–#111), and two corrections to what earlier runs inferred about it.** Twenty-two open issues, same set as every run since 2026-09-05: **#111–#94** are the daily `mule-alarm` noise (zero comments, telemetry not code), **#93** aclhook comment (mule-review, `8c61268` already pushed, its own recipe path), **#92** keyAgreement (parked on a stable upstream v1.4.0 per the waiting-on row), **#78** FDO (milestone-4.0, already escalated), **#1** untouched per standing instruction. New since the 2026-09-24 run: **#111** (today, 09-25 11:01:07Z, "nothing has landed in 16h"). Still zero machine-checkable code-fix candidates, so **no .mule/todo.md task lines added**. **Proposal, extending the 2026-09-24 one: close #94–#110**; leave #111 (live today) to self-expire. Two corrections, both from reading `tools/mule.sh` rather than the issue titles: **(1) the "missed days" in earlier entries (09-15, 09-20, 09-23) are evidence of activity, not its absence** — I previously called them timer gaps; the alarm is gated on a `.alarmed` sentinel that `beat` deletes on every land (`mule.sh:737,748-749`), so a day with no new issue is a day something landed, not a day the timer slept. **(2) the alarm body's closing line is wrong**: "Close it — the mule reopens a new one if the silence continues" (`mule.sh:764`) implies issue state drives the alarm; it does not. Nothing reads or reconciles the 18 open issues, so the pile can only grow and closing the old ones is safe but inert — and `cmd_refill` commits `mule: refill the queue` (mule.sh:719-720) **without calling `beat`**, so a doc-only refill never counts as landing work. The titles' ages also imply a beat roughly every day around 19:00–21:00Z (#111 filed 09-25 11:01Z at 16h → last beat ≈ 09-24 19:01Z; #110 filed 09-24 10:54Z at 15h → ≈ 09-23 19:54Z) while the ~11:00Z daily check has never once found under the 8h threshold — **what calls `beat` at that time is unverified**, and one command on the Pi settles it (`grep -E 'landed|checked:' /root/astrate-mule/.mule/log | tail -20`). The queue has still landed nothing since ~2026-09-04/05; if that idle streak is not intentional, the wedge is on the Pi and a queue review is the fix.
 
 ---
 

$ cd /root/astrate-mule && git status --porcelain && echo "=== new queue lines ===" && tail -2 .mule/todo.md | cut -c1-120
M .mule/for-giulio.md
 M .mule/todo.md
=== new queue lines ===
- [ ] docs-sync-hk-async-operation-param [auto]: document the `?async_operation` boolean query parameter on `POST /house
- [ ] docs-sync-hk-retention-zero-is-unset [auto]: document the `0 ≡ unset` retention rule in docs/api/astarte_houseke

Only the two intended files changed, no git operations, nothing implemented.

```
Done: Ran the docs-sync recipe against the housekeeping surface (least-covered of the five).
Files: .mule/todo.md (2 new [auto] task lines), .mule/for-giulio.md (one dated section, 3 prose
  escalations + a non-findings record)
Verified: rg route-vs-path comparison (all 5 routes documented, no path drift); status-code
  cross-check of all 5 ops against internal/housekeeping/http.go + service.go (all now correct);
  config-key sweep `rg -o '\bASTRATE_[A-Z_]+'` code vs docs/site/configuration-reference.md;
  yaml.Unmarshal on the spec (OK, 2 paths). `make -C docs build` could NOT be run — mkdocs is
  not installed on this box; both proposals carry it as their own verification step.
Unsure: whether the `async_operation` parameter belongs in the spec at all while it is still a
  tracked deviation (#68 in docs/UPSTREAM-EXPERIMENTAL.md, pending upstream v1.4 final) — I
  proposed documenting the current accept-and-ignore behaviour rather than omitting the
  parameter, but the reconcile-vs-promote call is #68's and is yours. Also: no containers, no
  Docker, no tests run — nothing in this recipe needed them.
Follow-ups: docs-sync-hk-async-operation-param, docs-sync-hk-retention-zero-is-unset (both
  appended to .mule/todo.md)
```

Two notes on scope decisions: I picked **housekeeping** because it had 4 prior docs-sync tasks against appengine's 16, and the recipe allows one surface per run. And per the recipe, this recipe does **not** ask for `gh issue create` as a deliverable, so nothing was filed — the whole output is the two task lines plus the escalation.
