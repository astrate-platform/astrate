slug: docs-sync-ae-data-envelope-metadata
verdict: done
at:  2ca634d
ran: 2026-10-09T17:11:47Z on DietPi in 245s


> build · big-pickle

I'll start by reading the two files.
→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astarte_appengine_api.yaml [offset=1680, limit=80]
→ Read pkg/astarteapi/envelope.go [offset=130, limit=50]
→ Read internal/appengine/data.go [offset=40, limit=80]
→ Read internal/appengine/data.go [offset=240, limit=100]
→ Read docs/api/astarte_appengine_api.yaml [offset=400, limit=90]

$ ls /root/astrate-mule/docs && cat /root/astrate-mule/docs/Makefile 2>/dev/null
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

$ rg -n 'DataEnvelopeObject' /root/astrate-mule/docs/api/astarte_appengine_api.yaml
416:                $ref: "#/components/schemas/DataEnvelopeObject"
468:                $ref: "#/components/schemas/DataEnvelopeObject"
685:                $ref: "#/components/schemas/DataEnvelopeObject"
730:                $ref: "#/components/schemas/DataEnvelopeObject"
1244:                $ref: "#/components/schemas/DataEnvelopeObject"
1290:                $ref: "#/components/schemas/DataEnvelopeObject"
1711:    DataEnvelopeObject:

$ ls -la /root/astrate-mule/docs/api/ && rg -n 'metadata|table_header|columns' /root/astrate-mule/docs/api/*.yaml
total 184
drwxr-xr-x 2 root root  4096 Oct  8 21:26 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 61063 Oct  5 13:29 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 17437 Oct  4 19:36 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 23478 Oct  8 21:30 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 47623 Oct  2 22:09 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 20586 Oct  6 22:49 astrate_native_api.yaml
/root/astrate-mule/docs/api/astarte_appengine_api.yaml:1534:        flat objects; `table` returns a `{data, metadata}` envelope with rows
/root/astrate-mule/docs/api/astarte_appengine_api.yaml:1535:        aligned to a sorted `table_header`; `disjoint_tables` pivots to a map of

$ rg -ln 'DataEnvelopeObject|astarte_appengine_api' /root/astrate-mule --glob '!docs/api/*.yaml' --glob '!docs/site/**'
/root/astrate-mule/docs/swagger-ui/index.html
/root/astrate-mule/docs/handoff/phase-1-memory.md
/root/astrate-mule/internal/swagger/swagger_test.go
/root/astrate-mule/internal/appengine/http.go
/root/astrate-mule/internal/appengine/service.go

→ Read internal/swagger/swagger_test.go
→ Read docs/api/astarte_appengine_api.yaml [offset=1500, limit=70]

$ rg -n 'additionalProperties|oneOf|anyOf' /root/astrate-mule/docs/api/astarte_appengine_api.yaml
101:                oneOf:
1018:                oneOf:
1570:          additionalProperties:
1575:          additionalProperties:
1623:      additionalProperties:
1637:          additionalProperties:
1643:          additionalProperties:
1758:          additionalProperties:

→ Read docs/api/astarte_appengine_api.yaml [offset=1225, limit=75]
→ Read docs/api/astarte_appengine_api.yaml [offset=90, limit=25]
→ Read docs/api/astarte_appengine_api.yaml [offset=1005, limit=30]
→ Read docs/api/astarte_appengine_api.yaml [offset=1565, limit=20]

$ cd /root/astrate-mule && git log --oneline -20
ed89bfc mule: refill the queue
70b14b1 mule: recipe docs-sync ran (762s)
a086f26 mule: recipe code-review ran (136s)
705c467 mule: recipe astarte-upstream ran (108s)
1ca603c mule: refill the queue
26c12e4 mule: recipe github-issues ran (155s)
271423d mule: recipe milestones ran (574s)
c1b43cd mule: recipe hygiene timed out
e30b5fd mule: log docs-sync-pairing-deviceid-path-param
1cda0fb mule: docs-sync-pairing-deviceid-path-param [auto]: give `components.parameters.DeviceID` the constraints `RegisterRequest.hw_id` already carries in docs/api/astarte_pairing_api.yaml — the parameter (yaml:367-377) is prose-only, while `hw_id` (yaml:386-398) got `pattern: '^[A-Za-z0-9_-]{22}$'` plus `minLength`/`maxLength` from docs-sync-pairing-deviceid-base64url, and `TestPairingDeviceIDEncodingDocumented` (internal/swagger/swagger_test.go:1086) pins only the two descriptions and hw_id's bounds, never the parameter's schema; the realm-management twin was constrained later by docs-sync-rm-deviceid-param and its test (swagger_test.go:1168) does hold the parameter's pattern against `deviceid.Parse`. Add the same three schema lines to the parameter, and record in its description the asymmetry the prose omits: a malformed `deviceID` is **401** on all three device-authenticated operations — `deviceid.Parse` failing inside `authenticateDevice` burns bcrypt and returns `ErrUnauthorized` (internal/pairing/service.go:381-385) → `WriteUnauthorized` (http.go:322-323), exactly as an unknown device or a wrong secret is (service.go:395-397), so the client is told its credentials are bad rather than its ID malformed — but **404** on the agent `unregisterDevice` (service.go:231-234 → http.go:326-327). Extend TestPairingDeviceIDEncodingDocumented to compile the new pattern and hold it against `deviceid.Parse` the way the realm-management test does. Do not encode the strict trailing-bits rule in the pattern; that stays open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
555b6e8 mule: log docs-sync-pairing-realm-name-pattern
ed301fb mule: docs-sync-pairing-realm-name-pattern [auto]: constrain `components.parameters.RealmName` in docs/api/astarte_pairing_api.yaml (yaml:359-365, a bare `type: string` described only as "The realm name.") with `pattern: '^[a-z][a-z0-9]*$'` — the realm name is CHECK-constrained to exactly that (migrations/000002_metadata.up.sql:7, read-only, never edit) and the two surfaces whose tasks already ran carry it on their twins (docs/api/astarte_housekeeping_api.yaml:246, docs/api/astrate_native_api.yaml:420 and 467). Because this one parameter is shared by all seven operations, the description must record that an off-pattern name gets three different answers, not one: the agent and device routes answer **401** (the realm is resolved by the middleware, which treats a miss as unauthenticated rather than 404 — internal/auth/middleware.go:66-72); `getHealth` answers **404** (service.Health returns store.ErrNotFound → WriteNotFound, internal/pairing/service.go:369-371, http.go:289-290); and `getVersion` answers **200** `{"data":"1.2.2"}` for any realm string whatsoever, because `GET /pairing/v1/{realm}/version` is registered with no middleware and no realm lookup (cmd/astrate/main.go:438 → observability.VersionHandler, internal/observability/compat.go:38-47) — say that in the version operation so nobody reads the shared parameter as a guarantee the handler enforces. Mirror `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go:1247) with a pairing twin only if it stays a few lines; the YAML edit is the deliverable. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
bf5aadb mule: blocked docs-sync-pairing-body-cap-400
d7e14cf mule: log docs-sync-pairing-register-404-unreachable
cfa6e76 mule: docs-sync-pairing-register-404-unreachable [auto]: drop the `404` `DeviceNotFound` response from `registerDevice` in docs/api/astarte_pairing_api.yaml (yaml:140-141) — it is documented but absent: the route is wrapped by `requireAgent` (internal/pairing/http.go:82-83), i.e. `mw.RequireRealm(auth.ClaimPairing)`, which resolves the realm *before* the handler and answers `401 {"errors":{"detail":"Unauthorized"}}` for an unknown one (internal/auth/middleware.go:66-72, "no existence oracle on auth failures"), so `Service.Register`'s own `GetRealmByName` (internal/pairing/service.go:183-186) can return `store.ErrNotFound` only if the realm is deleted between the middleware's lookup and the handler's — a race, not a contract. Nothing else in Register returns `store.ErrNotFound`: `RegisterDevice`'s `ErrDeviceAlreadyConfirmed` becomes 422 (service.go:209-213) and the CountDevices/SetPayloadFormatHint failures are DB errors → 500. No test covers it (internal/pairing/http_test.go asserts 404 only for the agent DELETE at :527 and for health at :595-598), and the `DeviceNotFound` component's own description — "unknown device ID" — does not match what that path would be even if it fired. The premise came from docs-sync-pairing-register-404, which read service.go without the middleware in front of it. Keep the `404` on `unregisterDevice` (reachable: unknown device, and a malformed `deviceID` via service.go:231-234) and on `getHealth`. State the reachability argument in the report; if you find a reachable path instead, keep the response and fix its description rather than deleting it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
e06b1d8 mule: refill the queue
1753d91 mule: recipe docs-sync ran (693s)
59c62e8 mule: log engine-introspection-producer-roundtrip
d070345 mule: engine-introspection-producer-roundtrip [auto]: pin the producer half of the docs/DESIGN.md §3.3 introspection rule on the default gate — `testutil.Introspection` (internal/testutil/astartedevice.go:93-105) is called from six places, all behind `e2e`, `integration && e2e` or test/conformance build tags, while the consumer `parseIntrospection` (internal/engine/introspection.go:99-127) is tested only against hand-written literals in internal/engine/introspection_test.go:17-71, so nothing container-free checks that the string the harness renders is the one the engine accepts or that the documented "deterministic (sorted) ordering" holds. Add a case to internal/engine/introspection_test.go that builds the string from an unsorted multi-entry map via `testutil.Introspection`, asserts `parseIntrospection` returns the same map, and pins the exact output (sorted, `name:major:minor;…`, no trailing `;`, empty map → empty string).
4693262 mule: blocked testutil-control-frame-test

$ cd /root/astrate-mule && git log --oneline -30 --all --grep='docs-sync' 2>/dev/null; echo '---todo---'; cat .mule/todo.md 2>/dev/null | head -60
70b14b1 mule: recipe docs-sync ran (762s)
e30b5fd mule: log docs-sync-pairing-deviceid-path-param
1cda0fb mule: docs-sync-pairing-deviceid-path-param [auto]: give `components.parameters.DeviceID` the constraints `RegisterRequest.hw_id` already carries in docs/api/astarte_pairing_api.yaml — the parameter (yaml:367-377) is prose-only, while `hw_id` (yaml:386-398) got `pattern: '^[A-Za-z0-9_-]{22}$'` plus `minLength`/`maxLength` from docs-sync-pairing-deviceid-base64url, and `TestPairingDeviceIDEncodingDocumented` (internal/swagger/swagger_test.go:1086) pins only the two descriptions and hw_id's bounds, never the parameter's schema; the realm-management twin was constrained later by docs-sync-rm-deviceid-param and its test (swagger_test.go:1168) does hold the parameter's pattern against `deviceid.Parse`. Add the same three schema lines to the parameter, and record in its description the asymmetry the prose omits: a malformed `deviceID` is **401** on all three device-authenticated operations — `deviceid.Parse` failing inside `authenticateDevice` burns bcrypt and returns `ErrUnauthorized` (internal/pairing/service.go:381-385) → `WriteUnauthorized` (http.go:322-323), exactly as an unknown device or a wrong secret is (service.go:395-397), so the client is told its credentials are bad rather than its ID malformed — but **404** on the agent `unregisterDevice` (service.go:231-234 → http.go:326-327). Extend TestPairingDeviceIDEncodingDocumented to compile the new pattern and hold it against `deviceid.Parse` the way the realm-management test does. Do not encode the strict trailing-bits rule in the pattern; that stays open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
555b6e8 mule: log docs-sync-pairing-realm-name-pattern
ed301fb mule: docs-sync-pairing-realm-name-pattern [auto]: constrain `components.parameters.RealmName` in docs/api/astarte_pairing_api.yaml (yaml:359-365, a bare `type: string` described only as "The realm name.") with `pattern: '^[a-z][a-z0-9]*$'` — the realm name is CHECK-constrained to exactly that (migrations/000002_metadata.up.sql:7, read-only, never edit) and the two surfaces whose tasks already ran carry it on their twins (docs/api/astarte_housekeeping_api.yaml:246, docs/api/astrate_native_api.yaml:420 and 467). Because this one parameter is shared by all seven operations, the description must record that an off-pattern name gets three different answers, not one: the agent and device routes answer **401** (the realm is resolved by the middleware, which treats a miss as unauthenticated rather than 404 — internal/auth/middleware.go:66-72); `getHealth` answers **404** (service.Health returns store.ErrNotFound → WriteNotFound, internal/pairing/service.go:369-371, http.go:289-290); and `getVersion` answers **200** `{"data":"1.2.2"}` for any realm string whatsoever, because `GET /pairing/v1/{realm}/version` is registered with no middleware and no realm lookup (cmd/astrate/main.go:438 → observability.VersionHandler, internal/observability/compat.go:38-47) — say that in the version operation so nobody reads the shared parameter as a guarantee the handler enforces. Mirror `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go:1247) with a pairing twin only if it stays a few lines; the YAML edit is the deliverable. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
bf5aadb mule: blocked docs-sync-pairing-body-cap-400
d7e14cf mule: log docs-sync-pairing-register-404-unreachable
cfa6e76 mule: docs-sync-pairing-register-404-unreachable [auto]: drop the `404` `DeviceNotFound` response from `registerDevice` in docs/api/astarte_pairing_api.yaml (yaml:140-141) — it is documented but absent: the route is wrapped by `requireAgent` (internal/pairing/http.go:82-83), i.e. `mw.RequireRealm(auth.ClaimPairing)`, which resolves the realm *before* the handler and answers `401 {"errors":{"detail":"Unauthorized"}}` for an unknown one (internal/auth/middleware.go:66-72, "no existence oracle on auth failures"), so `Service.Register`'s own `GetRealmByName` (internal/pairing/service.go:183-186) can return `store.ErrNotFound` only if the realm is deleted between the middleware's lookup and the handler's — a race, not a contract. Nothing else in Register returns `store.ErrNotFound`: `RegisterDevice`'s `ErrDeviceAlreadyConfirmed` becomes 422 (service.go:209-213) and the CountDevices/SetPayloadFormatHint failures are DB errors → 500. No test covers it (internal/pairing/http_test.go asserts 404 only for the agent DELETE at :527 and for health at :595-598), and the `DeviceNotFound` component's own description — "unknown device ID" — does not match what that path would be even if it fired. The premise came from docs-sync-pairing-register-404, which read service.go without the middleware in front of it. Keep the `404` on `unregisterDevice` (reachable: unknown device, and a malformed `deviceID` via service.go:231-234) and on `getHealth`. State the reachability argument in the report; if you find a reachable path instead, keep the response and fix its description rather than deleting it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
1753d91 mule: recipe docs-sync ran (693s)
10efb94 mule: recipe docs-sync ran (518s)
3fe292e mule: recipe docs-sync ran (816s)
230dd11 mule: recipe docs-sync timed out
2ddbf7f mule: docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
3962921 mule: docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
7331e6c mule: docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
f381392 mule: recipe docs-sync ran (761s)
6468242 mule: log docs-sync-hk-realm-name-response-schemas
def0af6 mule: docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `components.parameters.RealmName` and `RealmCreate.realm_name` both carry `pattern: '^[a-z][a-z0-9]*$'` and `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go) holds both against the realm name CHECK extracted from migrations/000002_metadata.up.sql — but the read side is still bare `type: string`: `Realm.realm_name` (components.schemas.Realm) and `DataEnvelopeRealmNames.items`. Decide whether a response field earns the same pattern (a generated client learns nothing new from it, but a reader diffing create against get sees one constrained and one not) and, if it does, add it to both and extend that test to hold the response-side pattern against the same CHECK. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
7e56b19 mule: blocked docs-sync-hk-wrong-type-field-400
d954574 mule: blocked docs-sync-hk-patch-422-field-error-examples
2c5d007 mule: log docs-sync-hk-errordetail-schema-example
5044143 mule: docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362) — the second half docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail: Bad Request`, the reconstructed title-case form that `TestHousekeepingErrorDetailExamples` (internal/swagger/swagger_test.go:687) now pins against everywhere else on this spec, while the wire emits the frozen `astarteapi.DetailBadRequest` = `"Bad request"` (pkg/astarteapi/envelope.go:37, via `WriteBadRequest` on every 400 — internal/housekeeping/http.go:65, 166, 171, 183). Because the example sits on the generic `ErrorDetail` schema, which all eight error responses `$ref` (yaml:377-477), it is the example a generator falls back to for any response whose own example it cannot resolve — so it teaches `"Bad Request"` for a 404 and a 500 too. Change it to `Bad request`, or drop it as the realm-management and pairing `ErrorDetail` schemas do (neither carries one). Then extend `TestHousekeepingErrorDetailExamples`, whose comment currently says the schema-level example is "deliberately not pinned" (swagger_test.go:684-686), to pin it like the responses. One thing not to "fix": the example is nested one level in — it is the `errors` sub-schema's example, not a bare envelope — so do not re-wrap it in another `errors:`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
0ba8d65 mule: log docs-sync-hk-realm-name-pattern
e728217 mule: docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
38dc881 mule: recipe docs-sync ran (474s)
5fb0e21 mule: log docs-sync-hk-error-detail-examples-split
23006d4 mule: docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
87a86e6 mule: log docs-sync-native-error-example-capitalisation
7cc237f mule: docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `"Internal server error"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
4506fb9 mule: blocked docs-sync-hk-error-example-capitalisation
---todo---
- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
- [x] triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
- [x] compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).

# Mule queue

One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).
The mule runs the topmost open line. Reorder freely; the order is the priority.

If a line needs more detail than fits, write `.mule/tasks/<slug>.md` and keep the line short.

    tools/mule.sh add "<slug>: <outcome>"
    tools/mule.sh loop

A line tagged `[legion]` needs the Legion Go and is skipped automatically while it is asleep;
the queue moves on to the next runnable line rather than stalling behind it.

A line tagged `[readonly]` is a **standing check**, not a piece of work. It verifies something
and is expected to produce no diff at all — so an empty diff is success, not the `wrote
nothing` failure it would be for any other line. It never gets ticked off, because a gate is
never done. Its output lands in `.mule/reports/<slug>.md` with the sha it ran against, and it
is skipped while that sha is still HEAD: re-checking code that has not moved buys nothing and
spends a call on a free provider.

**`mule/queue` is authoritative for this file.** The mule ticks on the Pi and writes its
progress there, so editing the queue on `main` resurrects completed tasks the moment the two
are merged — which has happened once. Add and reorder tasks on `mule/queue`; let them reach
`main` only when that branch is merged.

A line marked `- [~]` is **parked**: real work, but not tick-sized. The mule only ever
picks up `- [ ]`. Benchmark runs live here because a single tier is 5-20 minutes of ingest
alone -- past the per-task budget -- and because they want someone watching. Run one by hand:

    tools/mule.sh legion bench-push
    ssh legion 'cd ~/astrate/bench && ./scripts/run-tier.sh small astrate -base-url ... -housekeeping-key ...'

## Where tasks come from

**This file is not the whole queue, and for real work it is not even the main part of it.**

The queue is: the standing lines below, plus **every open GitHub issue labelled `mule`**.
Issues are read live on each tick and are never copied into this file — a copy would be a
second place the same fact lives, on a branch the mule commits to and you edit on `main`,
and that produced three merge conflicts in one afternoon.

**To give the mule work, file an issue and label it `mule`.** From anywhere, by anyone,
including another model with repo access. No SSH, no editing this file:

    gh issue create --label mule --title "<slug>: <outcome>" --body "<the detail>"

Labels on the issue are the tags: `legion` and `readonly` mean what `[legion]` and
`[readonly]` mean here. State lives on the issue, as labels, because there is exactly one
copy of it there:

| label          | meaning                                                          |
|----------------|------------------------------------------------------------------|
| `mule`         | queued                                                            |
| `mule-review`  | the mule pushed something; **it is not merged and not reviewed**  |
| `mule-blocked` | it tried and could not; re-label `mule` to try again              |

$ cd /root/astrate-mule && ls tools && rg -n 'keep.*test|revert|implementation|test' tools/mule.sh | head -40
bsoncapture
mule-plan-day.sh
mule-setup-pi.sh
mule-survey.sh
mule.sh
reconcile.sh
21:#   mule.sh revert           undo the last mule commit
163:# MULE.md tells the mule that a behaviour change must come with a test that fails without it.
164:# That was etiquette: nothing checked, and a cheap model can always write a test that asserts
165:# something already true. This makes it mechanical — take the implementation back out, keep
166:# the tests, and require them to break. If they still pass, the test proves nothing about the
169:# Only fires when the diff has both test and non-test Go changes. A pure refactor with no new
170:# test, or a test-only commit, is not this gate's business.
172:  local impl tests
173:  # A brand-new test file is untracked, and `git diff` cannot see it — so without this the
174:  # gate silently passed exactly the case it exists for: a new test added alongside a change.
177:  impl="$(git -C "$REPO" diff --name-only -- '*.go' ':!*_test.go')"
178:  tests="$(git -C "$REPO" diff --name-only -- '*_test.go')"
179:  [ -n "$impl" ] && [ -n "$tests" ] || return 0
182:  git -C "$REPO" diff -- '*.go' ':!*_test.go' > "$patch"
185:  note "gate: do the new tests actually fail without the change?"
187:    # Could not isolate the implementation — say so rather than pretend the gate ran.
188:    note "  could not un-apply the implementation; skipping this gate"
193:  # A compile failure here is a pass, not an error: the test referencing something that no
195:  ( cd "$REPO" && eval "${MULE_TEST_CMD:-go test ./...}" ) >/dev/null 2>&1 && still_passes=1
198:    bad "could not restore the implementation after the proof gate — leaving the tree alone"
204:    bad "the new tests pass with the implementation removed — they do not prove the change"
207:  ok "the tests fail without the change, as they should"
221:# errors, test failures and golangci-lint findings all carry a file:line prefix.
241:      || { gate_reason "tests failed"; bad "$GATE_REASON"; rc=1; }
253:    # leaves a phantom deletion that the revert path would act on.
301:    if ( cd "$REPO" && eval "$MULE_TEST_CMD" ) >/dev/null 2>&1; then ok "test baseline is green"
302:    else bad "test baseline is ALREADY RED — fix it first, or every task will look failed"; fail=1; fi
572:Unreviewed: the gates passed — it compiles, the tests pass, a new test was shown to fail without the change, and no frozen file was touched. That is not the same as the change being worth having. Nothing merges to \`main\` until someone reads it:
778:# day's allowance spent on testing the runner rather than on running tasks is not a reason to
794:# never merges: the gates prove a change compiles, passes tests and does not touch the frozen
835:Read it before it reaches main. The gates prove it compiles, passes the tests, and leaves
836:the frozen files alone — not that it is worth having, and not that the tests are worth
863:cmd_revert() {
864:  # Skip the bookkeeping commits — "revert" means undo a task's *code*, and the log rows
870:  git -C "$REPO" revert --no-edit "$sha" && ok "reverted: $subj"
981:  revert)    shift; cmd_revert "$@";;

→ Read internal/appengine/data.go [offset=119, limit=125]
Now I'll make the schema edit.

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -1712,10 +1712,39 @@
   type: object
   required: [data]
   properties:
     data:
-      description: Interface data (shape depends on interface type and path).
+      description: |
+        Interface data. The shape depends on the interface type, the
+        endpoint and the requested `format`, and is not always an object.
+        The default `structured` format is an array: an individual
+        datastream is a list of `{value, timestamp}` samples and an object
+        interface is a list of documents. `table` is an array of rows
+        (see `metadata`). `disjoint_tables` is an object keyed by path:
+        `{"value": [[value, timestamp], ...]}` for an individual
+        datastream, or one such pair list per column name for an object
+        interface. An interface-root individual query returns a nested
+        `{segment: ...: {value, timestamp}}` tree regardless of `format`.
+      oneOf:
+        - type: array
+        - type: object
+    metadata:
+      description: |
+        Column metadata for the rows in `data`. Present only when
+        `format=table`; omitted for every other format.
       type: object
+      required: [columns, table_header]
+      properties:
+        columns:
+          description: Column name to its zero-based index in each row.
+          type: object
+          additionalProperties:
+            type: integer
+        table_header:
+          description: Column names in row order.
+          type: array
+          items:
+            type: string
 
 DataEnvelopeStringArray:
   type: object
   required: [data]

$ cd /root/astrate-mule && python3 -c "import yaml,sys; d=yaml.safe_load(open('docs/api/astarte_appengine_api.yaml')); s=d['components']['schemas']['DataEnvelopeObject']; import json; print(json.dumps(s, indent=2))"
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml,sys; d=yaml.safe_load(open('docs/api/astarte_appengine_api.yaml')); s=d['components']['schemas']['DataEnvelopeObject']; import json; print(json.dumps(s, indent=2))
    ^^^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ cd /root/astrate-mule && make -C docs build 2>&1 | tail -30
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

$ cd /root/astrate-mule && cat docs/requirements.txt && pip install -q -r docs/requirements.txt 2>&1 | tail -5; which mkdocs
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10
/bin/sh: 1: pip: not found

$ cd /root/astrate-mule && python3 -m pip --version 2>&1; pyyaml_check=$(python3 -c "import yaml" 2>&1); echo "yaml: $pyyaml_check"; ls /usr/lib/python3/dist-packages 2>/dev/null | head; command -v pip3 python3
/usr/bin/python3: No module named pip
yaml: Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml
ModuleNotFoundError: No module named 'yaml'
Cryptodome
packaging
packaging-25.0.dist-info
pycryptodomex-3.20.0.egg-info

$ cd /root/astrate-mule && rg -n 'yaml' go.mod go.sum | head
go.sum:202:go.yaml.in/yaml/v2 v2.4.2 h1:DzmwEr2rDGHl7lsFgAHxmNz/1NlQ7xLIrlN2h5d1eGI=
go.sum:203:go.yaml.in/yaml/v2 v2.4.2/go.mod h1:081UH+NErpNdqlCXm3TtEran0rJZGxAYx9hb/ELlsPU=
go.sum:225:gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
go.sum:226:gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
go.sum:227:gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
go.mod:93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
go.mod:99:	gopkg.in/yaml.v3 v3.0.1 // indirect

← Write /tmp/opencode/yamlcheck.go
Wrote file successfully.

$ cd /root/astrate-mule && go run /tmp/opencode/yamlcheck.go docs/api/astarte_appengine_api.yaml 2>&1 | tail -50
parsed OK
properties:
    data:
        description: |
            Interface data. The shape depends on the interface type, the
            endpoint and the requested `format`, and is not always an object.
            The default `structured` format is an array: an individual
            datastream is a list of `{value, timestamp}` samples and an object
            interface is a list of documents. `table` is an array of rows
            (see `metadata`). `disjoint_tables` is an object keyed by path:
            `{"value": [[value, timestamp], ...]}` for an individual
            datastream, or one such pair list per column name for an object
            interface. An interface-root individual query returns a nested
            `{segment: ...: {value, timestamp}}` tree regardless of `format`.
        oneOf:
            - type: array
            - type: object
    metadata:
        description: |
            Column metadata for the rows in `data`. Present only when
            `format=table`; omitted for every other format.
        properties:
            columns:
                additionalProperties:
                    type: integer
                description: Column name to its zero-based index in each row.
                type: object
            table_header:
                description: Column names in row order.
                items:
                    type: string
                type: array
        required:
            - columns
            - table_header
        type: object
required:
    - data
type: object

The `make build` ran its `sync` target before failing. Let me check what that touched in the working tree.
$ cd /root/astrate-mule && git status --short
M docs/api/astarte_appengine_api.yaml

$ cd /root/astrate-mule && gofmt -l . ; echo "gofmt-exit:$?"
gofmt-exit:0

$ cd /root/astrate-mule && go test ./internal/swagger/ 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/swagger	0.114s

$ cd /root/astrate-mule && go vet ./... 2>&1 | tail -20 && echo "VET_OK" && go test ./... 2>&1 | tail -40
VET_OK
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
ok  	github.com/astrate-platform/astrate/examples/flow-container-echo	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)
ok  	github.com/astrate-platform/astrate/internal/config	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

→ Read docs/api/astarte_appengine_api.yaml [offset=1711, limit=50]

$ cd /root/astrate-mule && git diff --stat
docs/api/astarte_appengine_api.yaml | 31 ++++++++++++++++++++++++++++++-
 1 file changed, 30 insertions(+), 1 deletion(-)

Done: Added optional `metadata` (`columns: map<string,int>`, `table_header: []string`, table-only) to `DataEnvelopeObject` and corrected `data` (now `oneOf: [array, object]` with a shape-accurate description) in `docs/api/astarte_appengine_api.yaml`; all six data GET 200s `$ref` this shared schema.

Files: docs/api/astarte_appengine_api.yaml

Verified:
- `make -C docs build` -> could not run: `mkdocs: No such file or directory`, and no `pip`/`pip3` is installed to `pip install -r docs/requirements.txt`. `make` reached its `sync` target (which touched only gitignored `docs/site/**`), then failed at the `mkdocs` invocation.
- YAML parse of the edited spec with a throwaway `gopkg.in/yaml.v3` probe -> pass; `DataEnvelopeObject` parses with `data.oneOf: [array, object]` and the new `metadata` object.
- `go vet ./...` -> pass; `go test ./...` (no `-race`; race build unavailable on this host) -> pass; `gofmt -l .` -> clean.

Unsure:
- I used `oneOf: [array, object]` for `data` rather than dropping `type`; both are valid JSON Schema, and I picked `oneOf` as the clearer signal. I also marked `columns`/`table_header` as `required` within `metadata` since both renderers always emit both. Reverse either choice if the house style differs.
- `make -C docs build` is unverified for the mkdocs step only; the YAML itself is confirmed valid.

Follow-ups: none
