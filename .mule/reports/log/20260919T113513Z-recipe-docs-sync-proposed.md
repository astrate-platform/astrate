slug: recipe-docs-sync
verdict: proposed
at:  054d305
ran: 2026-09-19T11:35:13Z on DietPi in 114s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ ls docs/api/ && ls internal/ | rg 'appengine|housekeeping|pairing|realm|native'
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
appengine
housekeeping
pairing
realm

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*housekeeping*'
(no output)

$ rg -n '^\s{2}/' docs/api/astarte_housekeeping_api.yaml
26:  /housekeeping/v1/realms:
106:  /housekeeping/v1/realms/{realm}:

$ ls internal/housekeeping/
async_operation_test.go
http.go
http_test.go
service.go
service_test.go

$ rg -n 'housekeeping' internal/housekeeping/http.go | head -40
1:package housekeeping
19:// API is the /housekeeping/v1 HTTP surface (docs/ROADMAP.md §8.1 file 7.4).
26:// NewAPI wires the housekeeping service to its HTTP surface. instanceKeysPEM
34:// astarte_housekeeping).
37:	mux.Handle("GET /housekeeping/v1/realms", h(a.listRealms))
38:	mux.Handle("POST /housekeeping/v1/realms", h(a.createRealm))
39:	mux.Handle("GET /housekeeping/v1/realms/{realm}", h(a.getRealm))
40:	mux.Handle("PATCH /housekeeping/v1/realms/{realm}", h(a.patchRealm))
41:	mux.Handle("DELETE /housekeeping/v1/realms/{realm}", h(a.deleteRealm))
254:	const prefix = "housekeeping: validation failed: "

→ Read docs/api/astarte_housekeeping_api.yaml
Housekeeping routes match the spec 1:1. Now let me verify the response codes against the handlers.
→ Read internal/housekeeping/http.go
All five routes and their response codes match the spec. Let me also run the configuration-key drift check from the recipe.

$ rg -o 'ASTRATE_[A-Z_]+' docs/site/configuration-reference.md | sort -u
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

$ rg -o '\bASTRATE_[A-Z_]+' -N internal/ | sed 's/://' | sort -u
internal/config/config.example.tomlASTRATE_DATABASE_DSN
internal/config/config.example.tomlASTRATE_ENGINE_SHARDS
internal/config/config.example.tomlASTRATE_HTTP_ADDR
internal/config/config.example.tomlASTRATE_LOG_LEVEL
internal/config/config.example.tomlASTRATE_MASTER_KEY
internal/config/config.example.tomlASTRATE_MASTER_KEY_FILE
internal/config/config.example.tomlASTRATE_MQTT_INSECURE_DEV_MODE
internal/config/config.goASTRATE_DATABASE_DSN
internal/config/config.goASTRATE_ENGINE_SHARDS
internal/config/config.goASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
internal/config/config.goASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
internal/config/config.goASTRATE_HTTP_ADDR
internal/config/config.goASTRATE_HTTP_CORS_ALLOWED_ORIGINS
internal/config/config.goASTRATE_HTTP_TLS_CERT_FILE
internal/config/config.goASTRATE_HTTP_TLS_KEY_FILE
internal/config/config.goASTRATE_LOG_FORMAT
internal/config/config.goASTRATE_LOG_LEVEL
internal/config/config.goASTRATE_MASTER_KEY
internal/config/config.goASTRATE_MASTER_KEY_FILE
internal/config/config.goASTRATE_MQTT_ADDR
internal/config/config.goASTRATE_MQTT_ADVERTISED_URL
internal/config/config.goASTRATE_MQTT_INSECURE_DEV_MODE
internal/config/config.goASTRATE_MQTT_SESSION_STORE_PATH
internal/config/config.goASTRATE_MQTT_TLS_CERT_FILE
internal/config/config.goASTRATE_MQTT_TLS_KEY_FILE
internal/config/config.goASTRATE_REALM_JWT_PUBLIC_KEY
internal/config/config.goASTRATE_REALM_JWT_PUBLIC_KEY_FILE
internal/config/config.goASTRATE_REALM_NAME
internal/config/config.goASTRATE_SECURITY_MASTER_KEY_FILE
internal/config/config_test.goASTRATE_DATABASE_DSN
internal/config/config_test.goASTRATE_ENGINE_SHARDS
internal/config/config_test.goASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
internal/config/config_test.goASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
internal/config/config_test.goASTRATE_HTTP_CORS_ALLOWED_ORIGINS
internal/config/config_test.goASTRATE_MQTT_INSECURE_DEV_MODE
internal/engine/engine_integration_test.goASTRATE_TEST_DSN
internal/flow/blocks/container/block.goASTRATE_FLOW_CONFIG
internal/flow/blocks/container/docker.goASTRATE_FLOW_CONFIG
internal/flow/blocks/info.goASTRATE_FLOW_CONFIG
internal/flow/blocks/schema.goASTRATE_FLOW_CONFIG
internal/store/crypto.goASTRATE_MASTER_KEY
internal/store/crypto.goASTRATE_MASTER_KEY_FILE
internal/store/store_test.goASTRATE_TEST_DSN
internal/testutil/pg.goASTRATE_TEST_DSN

Config-key drift is already escalated (twice) in for-giulio.md, so nothing new there. Let me verify the schema field types against the service.
$ rg -n 'DeviceRegistrationLimit|DatastreamMaximumStorageRetention' internal/housekeeping/service.go | head -30
39:	DeviceRegistrationLimit           *int32
40:	DatastreamMaximumStorageRetention *int64
98:// WithDefaultDatastreamMaximumStorageRetention sets the realm-level datastream
103:func (s *Service) WithDefaultDatastreamMaximumStorageRetention(defaultRetention *int64) *Service {
175:		DeviceRegistrationLimit:           regLimit,
176:		DatastreamMaximumStorageRetention: retention,
285:		DeviceRegistrationLimit:           r.DeviceRegistrationLimit,
286:		DatastreamMaximumStorageRetention: r.DatastreamMaximumStorageRetention,

Housekeeping is clean end to end. Verification summary:

- **Paths**: code registers exactly the 5 routes `Mount` wires (http.go:37-41); the spec documents the same 5 (yaml:26, 106). No present-but-undocumented, no documented-but-absent.
- **Status codes** per operation vs. `writeError` (http.go:233-249) and handler bodies: GET list 200/401/403/500 ✓; POST create 201/400/401/403/409/422/500 ✓; GET realm 200/401/403/404/500 ✓; PATCH 200/400/401/403/404/422/500 ✓ (405 only fires on `ErrDeletionDisabled`, which only DELETE can hit); DELETE 204/401/403/404/405/422/500 ✓.
- **Schemas**: `Realm`/`RealmCreate`/`RealmPatch` field names, `int32`/`int64`, and nullability match the wire (`realmBody` http.go:46-51, `RealmView` service.go:39-40). Example omissions of the nullable retention field are cosmetic, not drift.
- **Config keys**: the two missing `ASTRATE_HOUSEKEEPING_*` keys in configuration-reference.md are already escalated to Giulio twice (for-giulio.md lines 51 and 133-137, both still open) — not re-escalating a duplicate.

No drift, so no task lines and no evidence file. Per the recipe: an accurate spec is the good result.

Done: docs-sync housekeeping run — no drift; all five routes, response codes, and schemas match the code, so nothing proposed.
Files: none
Verified: `rg` route/doc-set comparison, handler-vs-spec status-code comparison, field-type check vs service.go -> pass
Unsure: nothing (an example omitting the nullable retention field on GET/POST is cosmetic, not a spec error)
Follow-ups: none
