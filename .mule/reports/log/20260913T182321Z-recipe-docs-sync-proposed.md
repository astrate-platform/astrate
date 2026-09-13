slug: recipe-docs-sync
verdict: proposed
at:  5fecfd5
ran: 2026-09-13T18:23:21Z on DietPi in 184s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
✱ Glob "docs/api/astrate_*_api.yaml" in . · 1 match
✱ Glob "internal/**/*" in . · 100 matches
Let me start by understanding the repo structure and picking a surface to check.
→ Read .mule/todo.md

$ ls -la docs/api/
total 152
drwxr-xr-x 2 root root  4096 Sep  4 18:55 .
drwxr-xr-x 6 root root  4096 Sep  4 20:46 ..
-rw-r--r-- 1 root root 56583 Sep 12 22:44 astarte_appengine_api.yaml
-rw-r--r-- 1 root root 13227 Sep  9 19:14 astarte_housekeeping_api.yaml
-rw-r--r-- 1 root root 20064 Sep 10 21:10 astarte_pairing_api.yaml
-rw-r--r-- 1 root root 36818 Sep 11 19:05 astarte_realm_management_api.yaml
-rw-r--r-- 1 root root 14175 Sep  9 22:27 astrate_native_api.yaml

$ ls docs/ docs/site/ 2>/dev/null; ls .mule/
docs/:
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

docs/site/:
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

$ git -C /root/astrate-mule log --oneline -5; git -C /root/astrate-mule branch --show-current; git -C /root/astrate-mule status --short | head -20
5fecfd5 mule: log config-key-resolvers-test
9e68a01 mule: config-key-resolvers-test [auto]: add container-free unit tests for `HousekeepingKeys()` and `RealmJWTPublicKey()` (internal/config/config.go:383-411) — both file-reading helpers have zero test coverage. Cover: happy path (inline PEM + file PEM → concatenated/trimmed), file-not-found error (assert error message includes the path), inline-over-file priority (`RealmJWTPublicKey` returns the inline block when both are set), and TrimSpace behaviour (leading/trailing newlines stripped). `t.TempDir()` + `os.WriteFile` for fixtures; no Docker.
05a7b69 mule: log config-fail-loud-dev-mode
27d8165 mule: config-fail-loud-dev-mode [auto]: `applyEnv` silently discards a malformed `ASTRATE_MQTT_INSECURE_DEV_MODE` value (internal/config/config.go:255-257) — `strconv.ParseBool` error is discarded, the field becomes false, and TLS requirements are silently re-enabled (or a misleading "tls_cert_file required" error surfaces via validate). Return the parse error on bad input; add a failing test that asserts `ASTRATE_MQTT_INSECURE_DEV_MODE=yes` causes Load to error.
a101be2 mule: log config-fail-loud-engine-shards
mule/queue

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/housekeeping/ --glob '*.go'
(no output)

$ rg -n '^\s{2}/' docs/api/astarte_housekeeping_api.yaml
26:  /housekeeping/v1/realms:
106:  /housekeeping/v1/realms/{realm}:

I'm on `mule/queue` as expected. Now let me run the docs-sync comparison. I'll pick the **housekeeping** surface — smallest spec, least recently covered.
→ Read internal/housekeeping/http.go
The housekeeping surface has 5 routes (GET/POST realms, GET/PATCH/DELETE realm), both paths documented. Now let me read the spec to compare details.
→ Read docs/api/astarte_housekeeping_api.yaml
Both paths and all 5 routes are documented. Now let me verify the error-mapping paths against the Service, especially DELETE (which the spec lists without 405/422).
→ Read internal/housekeeping/service.go
Now I have the complete housekeeping picture. Comparing code vs spec:

**All 5 routes/paths documented** — no path drift. But two real response-level gaps:

1. **`DELETE /realms/{realm}` is missing its two deletion-gating responses** — the handler's `writeError` (http.go:233-249) answers 405 "Realm deletion disabled" (`ErrDeletionDisabled`) and 422 `{"errors":{"error_name":["connected_devices_present"]}}` (`ErrConnectedDevicesPresent`), both reachable from `Service.DeleteRealm` (service.go:248-266). The spec lists only 204/401/403/404/500 — and its *only* 405 reference sits on PATCH, which can never return it (UpdateRealm has no deletion-gating path, service.go:196-216).

2. **The `422` ValidationError example is wrong** — it shows `detail: "housekeeping: validation failed: realm name must not be empty"`, but `validationDetail` (http.go:251-258) strips the prefix on the wire and the real blank-name message is `"realm_name can't be blank"` (service.go:138).

Let me verify the exact line numbers before writing the task lines.

$ rg -n 'detail: "housekeeping' docs/api/astarte_housekeeping_api.yaml; rg -n 'MethodNotAllowed|ValidationErrors' docs/api/astarte_housekeeping_api.yaml
417:              detail: "housekeeping: validation failed: realm name must not be empty"
183:          $ref: "#/components/responses/MethodNotAllowed"
185:          $ref: "#/components/responses/ValidationErrors"
378:    MethodNotAllowed:
388:    ValidationErrors:

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -164,4 +164,6 @@
 - [x] docs-sync-ae-patch-merge-patch-media-type [auto]: fix the three PATCH request bodies in docs/api/astarte_appengine_api.yaml — PATCH /devices/{device} (yaml:176), PATCH /devices-by-alias/{alias} (yaml:260), PATCH /groups/{group}/devices/{device} (yaml:1023) declare `application/json; charset=utf-8`, but all three handlers enforce `Content-Type: application/merge-patch+json` exactly and reject every other value — 500 on the device and group paths (internal/appengine/http.go:209-211, 521-523), 400 on the alias path (http.go:230-232). Change the requestBody media type to `application/merge-patch+json` on all three. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
 - [x] config-fail-loud-dev-mode [auto]: `applyEnv` silently discards a malformed `ASTRATE_MQTT_INSECURE_DEV_MODE` value (internal/config/config.go:255-257) — `strconv.ParseBool` error is discarded, the field becomes false, and TLS requirements are silently re-enabled (or a misleading "tls_cert_file required" error surfaces via validate). Return the parse error on bad input; add a failing test that asserts `ASTRATE_MQTT_INSECURE_DEV_MODE=yes` causes Load to error.
 - [x] config-key-resolvers-test [auto]: add container-free unit tests for `HousekeepingKeys()` and `RealmJWTPublicKey()` (internal/config/config.go:383-411) — both file-reading helpers have zero test coverage. Cover: happy path (inline PEM + file PEM → concatenated/trimmed), file-not-found error (assert error message includes the path), inline-over-file priority (`RealmJWTPublicKey` returns the inline block when both are set), and TrimSpace behaviour (leading/trailing newlines stripped). `t.TempDir()` + `os.WriteFile` for fixtures; no Docker.
+- [ ] docs-sync-hk-delete-gating-responses [auto]: add the missing `405` and `422` responses to `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — `writeError` answers 405 "Realm deletion disabled" for `ErrDeletionDisabled` and 422 `{"errors":{"error_name":["connected_devices_present"]}}` for `ErrConnectedDevicesPresent` (internal/housekeeping/http.go:233-249), both reachable from `Service.DeleteRealm` (service.go:248-266), but the spec lists only 204/401/403/404/500 (yaml:198-208). Its sole `405` reference sits on PATCH (yaml:182-183), which can never return it (`UpdateRealm`, service.go:196-216, has no deletion-gating path) — re-point 405 to DELETE (reusing the MethodNotAllowed component) and add 422 with the FieldErrors shape and a connected_devices_present example. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_housekeeping_api.yaml (yaml:417) — it shows `detail: "housekeeping: validation failed: realm name must not be empty"`, but `validationDetail` strips the `housekeeping: validation failed: ` prefix on the wire (internal/housekeeping/http.go:251-258) and the real blank-name message is `"realm_name can't be blank"` (service.go:138). Change the example to `detail: "realm_name can't be blank"`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

Confirmed: example at yaml:417, 405/422 refs live only on PATCH (yaml:182-185). Both findings are `docs/api/*.yaml` fixes, so they go to `.mule/todo.md` as task lines.
Done: compared the housekeeping surface (chosen: smallest spec, least recently covered) against its handlers and proposed two `docs/api` fixes.

```
Done: docs-sync sweep of the housekeeping surface — all 5 routes/paths match; found 2 response-level drifts, proposed as task lines
Files: .mule/todo.md
Verified: rg route/path comparison + handler/spec read-through (internal/housekeeping/http.go, service.go, astarte_housekeeping_api.yaml) -> pass
Unsure: nothing
Follow-ups: docs-sync-hk-delete-gating-responses, docs-sync-hk-validation-example (appended)
```

Findings: (1) `DELETE /realms/{realm}` is missing 405 + 422 (its real deletion-gating responses), while the phantom 405 sits on PATCH which can never return it; (2) the 422 ValidationError example shows the internal prefix-stripped string `"housekeeping: validation failed: realm name must not be empty"` but the wire emits `"realm_name can't be blank"` (http.go:251-258, service.go:138). No docs/site prose conflicts found; no `gh issue create` since the recipe doesn't call for one. Not touching git.
