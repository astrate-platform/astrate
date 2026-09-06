slug: recipe-docs-sync
verdict: proposed
at:  f68f157
ran: 2026-09-06T11:20:06Z on DietPi in 168s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md

# Todos
[•] Pick one surface and find routes in code
[ ] Compare documented paths vs code paths
[ ] Check 3-4 endpoints for status codes / fields
[ ] Check config keys drift
[ ] Propose tasks in .mule/todo.md (max 5)

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' /root/astrate-mule/internal/ --glob '*realm*management*' 2>/dev/null | head -80
(no output)

I'll pick **realm-management** as the surface this run. Let me find the routes in code and the documented paths.
$ ls /root/astrate-mule/docs/api/*realm* 2>/dev/null
/root/astrate-mule/docs/api/astarte_realm_management_api.yaml

$ ls /root/astrate-mule/internal/
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

The glob for realm management code files didn't match. Let me explore the internal directory structure.
$ rg -l 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' /root/astrate-mule/internal/ | head -40
/root/astrate-mule/internal/flow/blocks/httpblocks.go
/root/astrate-mule/internal/flow/blocks/httpblocks_test.go
/root/astrate-mule/internal/engine/forward/http_test.go
/root/astrate-mule/internal/httpx/notfound_test.go
/root/astrate-mule/internal/httpx/cors.go
/root/astrate-mule/internal/appengine/stream/ws.go
/root/astrate-mule/internal/engine/triggers/actions_test.go
/root/astrate-mule/internal/appengine/http.go
/root/astrate-mule/internal/auth/middleware.go
/root/astrate-mule/internal/swagger/swagger.go
/root/astrate-mule/internal/swagger/swagger_test.go
/root/astrate-mule/internal/observability/compat.go
/root/astrate-mule/internal/observability/health.go
/root/astrate-mule/internal/observability/compat_version_test.go
/root/astrate-mule/internal/pairing/http.go

$ ls /root/astrate-mule/internal/realm/
async_operation_test.go
dashboard_compat_test.go
detailed.go
detailed_test.go
http.go
http_test.go
policy.go
policy_test.go
service.go

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' /root/astrate-mule/internal/realm/
(no output)

$ rg -n 'Handle|\.Get\(|\.Post\(|\.Put\(|\.Delete\(|\.Patch\(' /root/astrate-mule/internal/realm/http.go | head -60
25:	require func(http.Handler) http.Handler
37:	h := func(f http.HandlerFunc) http.Handler { return a.require(f) }
38:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces", h(a.listInterfaces))
39:	mux.Handle("POST /realmmanagement/v1/{realm}/interfaces", h(a.installInterface))
40:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}", h(a.listInterfaceMajors))
41:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.getInterface))
42:	mux.Handle("PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.updateInterface))
43:	mux.Handle("DELETE /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.deleteInterface))
44:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers", h(a.listTriggers))
45:	mux.Handle("POST /realmmanagement/v1/{realm}/triggers", h(a.createTrigger))
46:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers/{name}", h(a.getTrigger))
47:	mux.Handle("DELETE /realmmanagement/v1/{realm}/triggers/{name}", h(a.deleteTrigger))
48:	mux.Handle("GET /realmmanagement/v1/{realm}/config/auth", h(a.getAuth))
49:	mux.Handle("PUT /realmmanagement/v1/{realm}/config/auth", h(a.putAuth))
50:	mux.Handle("GET /realmmanagement/v1/{realm}/config/device_registration_limit", h(a.getRegistrationLimit))
51:	mux.Handle("GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention", h(a.getDatastreamMaximumStorageRetention))
52:	mux.Handle("GET /realmmanagement/v1/{realm}/version", h(a.getVersion))
53:	mux.Handle("DELETE /realmmanagement/v1/{realm}/devices/{device}", h(a.deleteDevice))
54:	mux.Handle("GET /realmmanagement/v1/{realm}/policies", h(a.listPolicies))
55:	mux.Handle("POST /realmmanagement/v1/{realm}/policies", h(a.createPolicy))
56:	mux.Handle("GET /realmmanagement/v1/{realm}/policies/{name}", h(a.getPolicy))
57:	mux.Handle("DELETE /realmmanagement/v1/{realm}/policies/{name}", h(a.deletePolicy))
150:	if r.URL.Query().Get("detailed") == "true" {

The realm package uses `detailed.go`/`http.go` but not the `r.(Get|Post...)` pattern. The recipe's grep pattern is generic. The realm-management surface maps to `internal/realm/`. Let me look at how routes are declared there and read the handler files.
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

Now let me compare the code routes against documented paths.

**Code routes:**
1. `GET/POST /interfaces` ✓
2. `GET /interfaces/{name}` ✓
3. `GET/PUT/DELETE /interfaces/{name}/{major}` ✓
4. `GET/POST /triggers` ✓
5. `GET/DELETE /triggers/{name}` ✓
6. `GET/PUT /config/auth` ✓
7. `GET /config/device_registration_limit` ✓
8. `GET /config/datastream_maximum_storage_retention` ← **NOT in doc!**
9. `GET /version` ✓
10. `DELETE /devices/{device}` ✓
11. `GET/POST /policies` ✓
12. `GET/DELETE /policies/{name}` ✓

So the **`/config/datastream_maximum_storage_retention`** endpoint is present in code but undocumented. That's a proposal candidate.

Let me establish the documented paths more precisely (they have GET/POST etc. listed in the yaml). Let me read the yaml for the config section to confirm what's documented, and check the detailed.go policies and status codes. Let me also look at the `detailed` handling in http.go which uses `/realmmanagement/v1/{realm}/interfaces?detailed=true` - is that documented?

Let me read the yaml sections to compare status codes/fields for 3-4 endpoints.
→ Read docs/api/astarte_realm_management_api.yaml [offset=361, limit=160]
Confirmed: the doc documents `config/auth` and `config/device_registration_limit`, but **not** `config/datastream_maximum_storage_retention` which is in code at http.go:51. Also note `config/auth` doc says PUT returns "204 on success" — let me verify what the code actually returns. Let me check the config/auth put handler and the datastream handler.

Let me look at http.go more thoroughly.
→ Read internal/realm/http.go [offset=120, limit=130]

$ rg -n 'func \(a \*API\) (putAuth|deleteDevice|getAuth)' /root/astrate-mule/internal/realm/http.go
105:func (a *API) deleteDevice(w http.ResponseWriter, r *http.Request) {
312:func (a *API) getAuth(w http.ResponseWriter, r *http.Request) {
321:func (a *API) putAuth(w http.ResponseWriter, r *http.Request) {

Let me check the putAuth handler and the devices delete handler for status codes. The doc says putAuth returns 204 and deleteDevice returns 204.
→ Read docs/api/astarte_realm_management_api.yaml [offset=520, limit=150]
Let me read the putAuth, deleteDevice, getAuth handlers to confirm status codes and response shapes.
→ Read internal/realm/http.go [limit=30, offset=105]
→ Read internal/realm/http.go [limit=40, offset=312]
Now I have enough. Let me check the config keys drift as the recipe's second drift category, and check the `detailed=true` param status in the yaml.

Config keys:

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
/root/astrate-mule/internal/config/config.go:ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
/root/astrate-mule/internal/config/config.go:ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
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
/root/astrate-mule/internal/config/config_test.go:ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION
/root/astrate-mule/internal/config/config_test.go:ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED
/root/astrate-mule/internal/config/config_test.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
/root/astrate-mule/internal/engine/engine_integration_test.go:ASTRATE_TEST_DSN
/root/astrate-mule/internal/flow/blocks/container/block.go:ASTRATE_FLOW_CONFIG
/root/astrate-mule/internal/flow/blocks/container/docker.go:ASTRATE_FLOW_CONFIG
/root/astrate-mule/internal/flow/blocks/info.go:ASTRATE_FLOW_CONFIG
/root/astrate-mule/internal/flow/blocks/schema.go:ASTRATE_FLOW_CONFIG
/root/astrate-mule/internal/store/crypto.go:ASTRATE_MASTER_KEY
/root/astrate-mule/internal/store/crypto.go:ASTRATE_MASTER_KEY_FILE
/root/astrate-mule/internal/store/store_test.go:ASTRATE_TEST_DSN
/root/astrate-mule/internal/testutil/pg.go:ASTRATE_TEST_DSN

Let me compare against the config keys documented in the site's configuration-reference.
$ rg -o '\bASTRATE_[A-Z_]+' -N /root/astrate-mule/docs/site/ 2>/dev/null | sort -u
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_DATABASE_DSN
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_ENGINE_SHARDS
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_HTTP_ADDR
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_HTTP_TLS_CERT_FILE
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_HTTP_TLS_KEY_FILE
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_LOG_FORMAT
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_LOG_LEVEL
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_MASTER_KEY
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_MASTER_KEY_FILE
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_MQTT_ADDR
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_MQTT_ADVERTISED_URL
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_MQTT_INSECURE_DEV_MODE
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_MQTT_SESSION_STORE_PATH
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_MQTT_TLS_CERT_FILE
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_MQTT_TLS_KEY_FILE
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_REALM_JWT_PUBLIC_KEY
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_REALM_NAME
/root/astrate-mule/docs/site/configuration-reference.md:ASTRATE_SECURITY_MASTER_KEY_FILE
/root/astrate-mule/docs/site/contributing.md:ASTRATE_TEST_DSN
/root/astrate-mule/docs/site/deployment.md:ASTRATE_DATABASE_DSN
/root/astrate-mule/docs/site/deployment.md:ASTRATE_ENGINE_SHARDS
/root/astrate-mule/docs/site/deployment.md:ASTRATE_HTTP_ADDR
/root/astrate-mule/docs/site/deployment.md:ASTRATE_LOG_LEVEL
/root/astrate-mule/docs/site/deployment.md:ASTRATE_MASTER_KEY
/root/astrate-mule/docs/site/deployment.md:ASTRATE_MASTER_KEY_FILE
/root/astrate-mule/docs/site/deployment.md:ASTRATE_MQTT_ADDR
/root/astrate-mule/docs/site/observability.md:ASTRATE_LOG_LEVEL
/root/astrate-mule/docs/site/operations.md:ASTRATE_DATABASE_DSN
/root/astrate-mule/docs/site/operations.md:ASTRATE_ENGINE_SHARDS
/root/astrate-mule/docs/site/operations.md:ASTRATE_HTTP_ADDR
/root/astrate-mule/docs/site/operations.md:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
/root/astrate-mule/docs/site/operations.md:ASTRATE_LOG_LEVEL
/root/astrate-mule/docs/site/operations.md:ASTRATE_MASTER_KEY
/root/astrate-mule/docs/site/operations.md:ASTRATE_MASTER_KEY_FILE
/root/astrate-mule/docs/site/operations.md:ASTRATE_MQTT_ADDR
/root/astrate-mule/docs/site/operations.md:ASTRATE_MQTT_INSECURE_DEV_MODE
/root/astrate-mule/docs/site/operations.md:ASTRATE_REALM_NAME
/root/astrate-mule/docs/site/troubleshooting.md:ASTRATE_DATABASE_DSN
/root/astrate-mule/docs/site/troubleshooting.md:ASTRATE_LOG_LEVEL

Comparing config.go keys (the source of truth) against the site's configuration-reference:

**In code but NOT in configuration-reference.md:**
- `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` (config.go)
- `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (config.go)

**In configuration-reference.md but NOT in config.go (from internal/ only):** All config.go keys look present. Let me check config.example.toml keys — `ASTRATE_ENGINE_SHARDS`, `ASTRATE_HTTP_ADDR`, `ASTRATE_LOG_LEVEL`, `ASTRATE_MASTER_KEY`, `ASTRATE_MASTER_KEY_FILE`, `ASTRATE_MQTT_INSECURE_DEV_MODE`, `ASTRATE_DATABASE_DSN` are all also in config.go, so consistent.

Wait — config.example.toml lacks several config.go keys, but that's a different file (example config), not the site's configuration-reference. The recipe specifically targets configuration-reference.md. Let me confirm the housekeeping config keys are real in config.go and referenced.

The two housekeeping keys (`ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` and `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`) are in code but not in `configuration-reference.md`. That's a `for-giulio.md` escalation (site prose is never-touch).

Now let me confirm the detailed listInterfaces `?detailed=true` variant — is it documented in the yaml? Let me check the interfaces GET section.
→ Read docs/api/astarte_realm_management_api.yaml [limit=100, offset=27]
The GET `/interfaces` is documented as "names of all interfaces" only, but code now also supports `?detailed=true` returning detailed interface listing (1.4-style, referenced to #66). The doc does not mention the `detailed` query parameter. That's a documented-interface drift.

Now let me verify the status codes I've claimed against code for a few endpoints:

1. `POST /interfaces` (installInterface) → 201 ✓ documented
2. `PUT /config/auth` (putAuth) → 204 ✓ documented (204)
3. `DELETE /devices/{device}` → 204 ✓ documented
4. `GET /config/datastream_maximum_storage_retention` → 200, **not documented at all**

Let me spot-check two more status-code claims by reading what `writeInterfaceError` returns vs the doc for DELETE /interfaces/{name}/{major}. The doc section at 127-239 covers interfaces/{name}/{major}. Let me quickly scan that for status codes.

The key finding is confirmed. Let me check the interface/{name}/{major} doc status codes quickly and verify update/delete 204.
→ Read docs/api/astarte_realm_management_api.yaml [limit=112, offset=127]
The interface/{name}/{major} status codes match the code (204 for update/delete, 404 for non-numeric major via `majorParam`, 422 for in-use). Documented status codes are accurate there.

I now have enough findings for the surface (realm-management). Let me consolidate the drift I found:

**Docs/api fix proposals (generated artefacts — can propose as tasks):**
1. `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` — present in code (http.go:51, getDatastreamMaximumStorageRetention at http.go:135) but **absent from** `astarte_realm_management_api.yaml` (after config/device_registration_limit). Fix: add the endpoint.
2. `GET /realmmanagement/v1/{realm}/interfaces` — the code supports `?detailed=true` returning a detailed listing (http.go:150, referenced to #66), but the yaml only documents the names-only response. Fix: add the `detailed` query param + detailed response schema.

**Site prose contradictions (for-giulio.md, never-touch):**
3. Config keys `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` and `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` exist in config.go but are undocumented in `docs/site/configuration-reference.md`.

The other endpoints' status codes/fields I checked matched the code. No wrong status codes found among the checked ones.

Let me check whether there's a `for-giulio.md` file and look at todo.md format before appending.

$ ls -la /root/astrate-mule/.mule/ | head; echo "---todo---"; tail -20 /root/astrate-mule/.mule/todo.md
total 740
drwxr-xr-x  7 root root   4096 Sep  6 13:17 .
drwxr-xr-x 16 root root   4096 Sep  4 21:26 ..
-rw-r--r--  1 root root      0 Sep  6 13:17 .alarmed
-rw-r--r--  1 root root     13 Sep  6 13:17 .budget
-rw-r--r--  1 root root    364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root     11 Sep  5 22:15 .heartbeat
-rw-r--r--  1 root root  16384 Sep  6 13:18 .last-output
-rw-r--r--  1 root root      0 Sep  6 13:17 .lock
-rw-r--r--  1 root root     10 Sep  6 13:17 .rotation
---todo---
- [ ] race-check-flow: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/flow/... ./internal/realm/... ./internal/pairing/... ./internal/auth/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
- [ ] race-check-appengine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/appengine/... ./internal/observability/... ./internal/httpx/... ./internal/config/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
- [ ] race-check-pkg: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./pkg/... ./cmd/... ./internal/testutil/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]

- [x] broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [approved 2026-09-04]
- [!] broker-disconnect-device-zombie-session: in `internal/broker/broker.go` `DisconnectDevice` (lines 260-266), force-closing the live connection leaves the persisted bbolt session (subscriptions + offline queue, keyed by CN) in place — mochi derives `expire` from the Clean flag, not the reason code (server.go:485-491), so `dropClient` (sessionstore.go:500) is never reached and a later device with the same CN resurrects the deleted device's session. Wire the session store handle into `Broker` and call `dropClient` from `DisconnectDevice`, and add a T1 test in `broker_test.go`: connect clean=false + subscribe, call `DisconnectDevice`, reconnect — assert `session_present` is false and the stale subscription is gone. [approved 2026-09-04] — BLOCKED: wrote nothing
- [!] broker-offline-acl-tests: in `internal/broker/aclhook_test.go`, unit-test the offline-delivery ACL — `offlineACL.ownershipOf` cache hit within TTL, TTL expiry triggering a reload, and a load failure caching as empty and denying (aclhook.go:116-138) — plus a T1 `OnACLCheck` check that the offline branch (delivery to a session-present but disconnected device, aclhook.go:196-202) consults the store. The rule currently has zero direct tests; `TestCheckACLMatrix` exercises only the pure `checkACL`. [approved 2026-09-04] — BLOCKED: wrote nothing
- [!] broker-onconnect-doc-comment: in `internal/broker/authhook.go:314`, restore the missing first line of the `OnConnect` doc comment — it opens mid-sentence ("of the Astarte MQTT v1 protocol, and mochi publishes them without a publish-side ACL check…") and the Will-clearing security rationale (a retained LWT would escape the §3.2 matrix) is only half-documented. [approved 2026-09-04] — BLOCKED: tests failed: --- FAIL: TestMQTTSink_Retained (0.02s)

- [x] purge-properties-compression-capability [auto]: upstream v1.3.0 adds a `purge_properties_compression_format` device capability (`zlib`|`plaintext`, default `zlib`) — a wire-visible capability value. Check whether Astrate's capabilities handling (internal/broker, the `<realm>/<device_id>/capabilities` topic, issue #16) needs to recognize/honour it, or whether zlib-only is already the deliberate default; propose the change or note why not needed.
- [!] empty-introspection-verification [auto]: upstream v1.3.0 changed "allow devices with empty introspection" — verify whether Astrate's device connection/introspection handling currently rejects an empty introspection string where upstream now accepts it, and propose a fix if so. — BLOCKED: wrote nothing
- [!] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages" — investigate only: does Astrate have an analogous delay between trigger install and first delivery? Report, do not patch. — BLOCKED: wrote nothing
- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
- [x] docs-sync-pairing-health-path [auto]: add the undocumented `GET /pairing/v1/{realm}/health` route to docs/api/astarte_pairing_api.yaml — it exists in code (internal/pairing/http.go:81) and returns 200 `{"data":{"status":"ok"}}` for an existing realm, 404 for an unknown realm, 503 when unhealthy, unauthenticated and unrate-limited. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-register-404 [auto]: document the missing `404` (DeviceNotFound) on `POST /pairing/v1/{realm}/agent/devices` in docs/api/astarte_pairing_api.yaml — the handler returns it when `GetRealmByName` fails for an unknown realm (internal/pairing/service.go:183-186 → http.go:326-327), but the spec lists only 201/400/401/403/422/429/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] swagger-httptest-coverage [auto]: add a container-free httptest suite for internal/swagger, which currently has no test file — cover two rules: (1) `Mount` wires `GET /swagger` to a 302→/swagger/index.html, serves the embedded UI under `/swagger/`, and serves the OpenAPI YAML specs under `/api/`; (2) `Specs` returns exactly the `.yaml` filenames from `docs.APIYAML` with no path prefix or dirs. httptest only, no Docker.
- [!] probe-property-resend-encoding [auto]: upstream v1.3.0 fixed outbound server-property values sent to a device on connect/emptyCache (commit 522ccf4f — a raw stored binaryblob is now wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding, so it ships as a BSON binary, not a string). Investigate only: does Astrate's resendServerProperties/rehydrate path (internal/engine/control.go:144, pkg/payload) emit binaryblob → BSON subtype-0 binary and datetime → UTC datetime exactly like upstream 1.3? Report, do not patch. — BLOCKED: wrote nothing
- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
- [!] flow-validate-source-sink: `Pipeline.Validate`'s no-source/no-sink error branches (internal/flow/pipeline.go:148-153) are unreachable — the preceding Kahn cycle check (pipeline.go:100-115) guarantees any graph that gets past it is acyclic and therefore always has a source and a sink. The two test cases "no source (all blocks have incoming edges)" and "no sink (all blocks have outgoing edges)" (pipeline_test.go:134-163) both use a 2-cycle, so they pass via the cycle error, never the branch they name. Make the rule real (a cycle-independent formulation, e.g. reject a graph where a node has no path to a sink) or delete the vacuous check, and fix the two test cases to assert the actual error message so a regression to wrong-path coverage fails. — BLOCKED: tests failed: --- FAIL: TestMQTTSink_Retained (0.02s)
- [x] flow-validate-dead-source-sink-recompute: in `Pipeline.Validate`, the first `hasSource`/`hasSink` computation (pipeline.go:117-126) is overwritten by the second (pipeline.go:128-147) and its result discarded — dead code (the comment at pipeline.go:127 documents it). Delete the first block; behaviour is unchanged. Tests: existing `TestPipeline_Validate` passes.

$ cat /root/astrate-mule/.mule/for-giulio.md 2>/dev/null | tail -30
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

I have what I need. The recipe says: proposals go as task lines in `.mule/todo.md` (for docs/api fixes), and site-prose contradictions go to `.mule/for-giulio.md`.

Let me identify the config keys in code missing from configuration-reference to nail down the for-giulio entry. From the earlier diff: `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` and `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` are in config.go but absent from configuration-reference.md. Let me confirm they're actual config fields (not just string literals) and check the for-giulio header format.
$ rg -n 'HOUSEKEEPING_DEFAULT_DATASTREAM|HOUSEKEEPING_REALM_DELETION' /root/astrate-mule/internal/config/config.go
270:		"ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
271:		"HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION",
288:	if v, ok := os.LookupEnv("ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED"); ok {

→ Read .mule/for-giulio.md [limit=15]
Confirmed both housekeeping config keys exist in config.go. These are undocumented in `configuration-reference.md`. Let me verify they're not in the config.example.toml either (just to fully characterize, though the recipe targets configuration-reference.md specifically).

Both findable. Now let me write the deliverables:

1. **Append task lines to `.mule/todo.md`** (docs/api fixes - the generated-artefact proposals), max 5. I have 2 solid api/api proposals: the datastream retention endpoint and the `detailed=true` param. That's within 5.

2. **Append to `.mule/for-giulio.md`** for the site configuration-reference prose drift (2 keys undocumented).

Let me draft the todo task lines. Note the existing completed docs-sync tasks follow a format like: `docs-sync-pairing-health-path [auto]: ...`. I'll follow the same style. Let me verify the datastream retention endpoint returns shape (data envelope) to describe the task accurately, and confirm the swagger spec loading consideration.

The recipe says: "If any task in this family touches docs/api/, the executing task must run the docs build." So proposed tasks should mention running `make -C docs build`. Good — the existing pairing docs-sync tasks did this.

Let me write the todo lines.
→ Read .mule/for-giulio.md [limit=6, offset=12]

$ tail -5 /root/astrate-mule/.mule/todo.md
- [x] swagger-httptest-coverage [auto]: add a container-free httptest suite for internal/swagger, which currently has no test file — cover two rules: (1) `Mount` wires `GET /swagger` to a 302→/swagger/index.html, serves the embedded UI under `/swagger/`, and serves the OpenAPI YAML specs under `/api/`; (2) `Specs` returns exactly the `.yaml` filenames from `docs.APIYAML` with no path prefix or dirs. httptest only, no Docker.
- [!] probe-property-resend-encoding [auto]: upstream v1.3.0 fixed outbound server-property values sent to a device on connect/emptyCache (commit 522ccf4f — a raw stored binaryblob is now wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding, so it ships as a BSON binary, not a string). Investigate only: does Astrate's resendServerProperties/rehydrate path (internal/engine/control.go:144, pkg/payload) emit binaryblob → BSON subtype-0 binary and datetime → UTC datetime exactly like upstream 1.3? Report, do not patch. — BLOCKED: wrote nothing
- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
- [!] flow-validate-source-sink: `Pipeline.Validate`'s no-source/no-sink error branches (internal/flow/pipeline.go:148-153) are unreachable — the preceding Kahn cycle check (pipeline.go:100-115) guarantees any graph that gets past it is acyclic and therefore always has a source and a sink. The two test cases "no source (all blocks have incoming edges)" and "no sink (all blocks have outgoing edges)" (pipeline_test.go:134-163) both use a 2-cycle, so they pass via the cycle error, never the branch they name. Make the rule real (a cycle-independent formulation, e.g. reject a graph where a node has no path to a sink) or delete the vacuous check, and fix the two test cases to assert the actual error message so a regression to wrong-path coverage fails. — BLOCKED: tests failed: --- FAIL: TestMQTTSink_Retained (0.02s)
- [x] flow-validate-dead-source-sink-recompute: in `Pipeline.Validate`, the first `hasSource`/`hasSink` computation (pipeline.go:117-126) is overwritten by the second (pipeline.go:128-147) and its result discarded — dead code (the comment at pipeline.go:127 documents it). Delete the first block; behaviour is unchanged. Tests: existing `TestPipeline_Validate` passes.

Newest at top. In for-giulio, I prepend after line 11 (`---`). Let me read the exact current first entry to place the insertion correctly, and check todo.md first line to prepend there too.

The todo.md entries — I'll append to the end of the file. Let me check the very end of todo.md to append cleanly.
$ cat >> /root/astrate-mule/.mule/todo.md <<'EOF'
- [ ] docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` route to docs/api/astarte_realm_management_api.yaml — it exists in code (internal/realm/http.go:51, handler getDatastreamMaximumStorageRetention at http.go:135, returns 200 with a data-envelope retention value via Service.GetDatastreamMaximumStorageRetention) but the spec jumps straight from config/device_registration_limit (yaml:428) to /version. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
EOF
echo done
done

Now let me append the two docs/api proposal task lines to `todo.md`, and prepend the two site-prose contradictions to `for-giulio.md`.
← Edit .mule/for-giulio.md
Index: /root/astrate-mule/.mule/for-giulio.md
===================================================================
--- /root/astrate-mule/.mule/for-giulio.md
+++ /root/astrate-mule/.mule/for-giulio.md
@@ -9,8 +9,15 @@
 line once you have dealt with it — this file is a queue, not a log.
 
 ---
 
+- **docs-sync realm-management: two ASTRATE_ config keys exist in code but are missing from
+  `docs/site/configuration-reference.md`** (site prose is yours, never edited by the mule).
+  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` and
+  `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` are consumed in `internal/config/config.go`
+  (lines 270-271 and 288); neither appears in `docs/site/configuration-reference.md`.
+  Add them or confirm they are meant to be internal-only.
+
 - **github-issues triage run, 2026-09-05: nothing proposable, nothing stale.** All five open
   issues were already covered by another mechanism. #94 is a mule-alarm (nothing landing in
   15h), not a code issue — the idle window is the 2026-09-04 queue rebuild, expected. #93 has
   a pushed commit (`8c61268`) awaiting review, its own recipe path. #92 is parked on a stable

Now the for-giulio entry for the site prose contradiction (undocumented config keys). Newest at top, after the `---` on line 11.
# Todos
[✓] Pick one surface and find routes in code
[✓] Compare documented paths vs code paths
[✓] Check 3-4 endpoints for status codes / fields
[✓] Check config keys drift
[✓] Propose tasks in .mule/todo.md (max 5)

Now let me update my todos and verify the results.
Let me verify both edits landed cleanly.
$ tail -3 /root/astrate-mule/.mule/todo.md; echo "=== for-giulio ==="; sed -n '11,22p' /root/astrate-mule/.mule/for-giulio.md
- [x] flow-validate-dead-source-sink-recompute: in `Pipeline.Validate`, the first `hasSource`/`hasSink` computation (pipeline.go:117-126) is overwritten by the second (pipeline.go:128-147) and its result discarded — dead code (the comment at pipeline.go:127 documents it). Delete the first block; behaviour is unchanged. Tests: existing `TestPipeline_Validate` passes.
- [ ] docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` route to docs/api/astarte_realm_management_api.yaml — it exists in code (internal/realm/http.go:51, handler getDatastreamMaximumStorageRetention at http.go:135, returns 200 with a data-envelope retention value via Service.GetDatastreamMaximumStorageRetention) but the spec jumps straight from config/device_registration_limit (yaml:428) to /version. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
=== for-giulio ===
---

- **docs-sync realm-management: two ASTRATE_ config keys exist in code but are missing from
  `docs/site/configuration-reference.md`** (site prose is yours, never edited by the mule).
  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` and
  `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` are consumed in `internal/config/config.go`
  (lines 270-271 and 288); neither appears in `docs/site/configuration-reference.md`.
  Add them or confirm they are meant to be internal-only.

- **github-issues triage run, 2026-09-05: nothing proposable, nothing stale.** All five open
  issues were already covered by another mechanism. #94 is a mule-alarm (nothing landing in
  15h), not a code issue — the idle window is the 2026-09-04 queue rebuild, expected. #93 has

I've completed the docs-sync proposal run, exploring the realm-management surface.

```
Done: docs-sync realm-management — found 2 api-spec fixes (todo.md) + 1 site-prose contradiction (for-giulio.md)
Files: .mule/todo.md, .mule/for-giulio.md
Verified: rg route vs spec path comparison (realm-management surface) -> pass; status-code spot-check (putAuth 204, deleteDevice 204, interface put/delete 204, majorParam 404) -> all match spec
Unsure: none
Follow-ups: - docs-sync-rm-datastream-retention-endpoint: add undocumented GET .../config/datastream_maximum_storage_retention to astarte_realm_management_api.yaml
  - docs-sync-rm-interfaces-detailed-param: add ?detailed=true param + detailed response to GET /interfaces in astarte_realm_management_api.yaml
  (No implementation done; recipes say proposal only.)
```

**Findings:**
- **Doc/api fixes (proposed in todo.md, 2 of max 5):** `GET /config/datastream_maximum_storage_retention` is in code (http.go:51/135) but absent from the yaml; and `GET /interfaces` `?detailed=true` (1.4-style, #66) is undocumented in the spec.
- **Site prose (for-giulio.md, never-touch):** `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION` and `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` exist in config.go but are missing from `configuration-reference.md`.
- Status codes/fields I spot-checked (config/auth put, devices delete, interface update/delete, majorParam 404) all matched the spec — no drift there. Stopped at 2 proposals per "an accurate spec is a good result; say so and stop."
