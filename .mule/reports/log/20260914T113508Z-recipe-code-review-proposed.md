slug: recipe-code-review
verdict: proposed
at:  6bca5cc
ran: 2026-09-14T11:35:08Z on DietPi in 469s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
→ Read .mule/log.md
✱ Glob ".mule/reviews/*" in . · 0 matches

$ ls /root/astrate-mule/internal/ /root/astrate-mule/pkg/ 2>/dev/null
/root/astrate-mule/internal/:
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

/root/astrate-mule/pkg/:
astarteapi
deviceid
interfaceschema
payload

$ rg -c '' --glob 'internal/auth/**/*.go' /root/astrate-mule | sort -t: -k2 -rn
/root/astrate-mule/internal/auth/middleware_test.go:271
/root/astrate-mule/internal/auth/claims_test.go:249
/root/astrate-mule/internal/auth/jwt_test.go:248
/root/astrate-mule/internal/auth/claims.go:204
/root/astrate-mule/internal/auth/middleware.go:196
/root/astrate-mule/internal/auth/cache_test.go:144
/root/astrate-mule/internal/auth/jwt.go:139
/root/astrate-mule/internal/auth/cache.go:81

→ Read internal/auth/claims.go
→ Read internal/auth/middleware.go
`auth` hasn't appeared in the log at all — a good pick. Reading it, largest first.
→ Read internal/auth/jwt.go
→ Read internal/auth/cache.go
→ Read internal/auth/claims_test.go
Now the tests.
→ Read internal/auth/jwt_test.go
→ Read internal/auth/cache_test.go
→ Read internal/auth/middleware_test.go
✱ Grep "RequireRealmAny|RequireRealm\(|RequireStatic\(|NewMiddleware\(" in . · 33 matches
A tightly-written package with strong tests. Let me check callers and the key source to look for usage-level gaps, and cross-check what the conformance tests assert.
✱ Grep "GetRealmByName" in internal/store · 19 matches
→ Read cmd/astrate/main.go [offset=330, limit=80]
✱ Grep "AuthorizesChannel|Authorizes\(" in . · 26 matches
→ Read internal/appengine/stream/ws.go
✱ Grep "RequireRealm|Mount|AuthorizesChannel|VerbJoin|VerbWatch" in internal/appengine/channels · 20 matches
✱ Grep "a_ch|ClaimChannels|\[\]\string\{|join|JOIN" in internal/appengine/stream · 1 match
→ Read internal/appengine/channels/ws.go [limit=120, offset=240]
→ Read internal/appengine/channels/ws.go [limit=240]
✱ Grep "a_ch|socket|Channels|channels" in docs · 78 matches
Let me check the frozen spec's socket auth wording and the current todo/for-giulio files.
→ Read .mule/todo.md
The todo has strong context. Let me verify the reviews dir, the for-giulio file, and the current branch (read-only) before writing output.

$ ls /root/astrate-mule/.mule/reviews/ 2>&1; echo "---"; ls /root/astrate-mule/.mule/ ; echo "---"; git -C /root/astrate-mule branch --show-current; git -C /root/astrate-mule status --porcelain | head
appengine-2026-09-07.md
broker-2026-09-04.md
config-2026-09-13.md
engine-2026-09-06.md
flow-2026-09-05.md
flowapi-2026-09-10.md
interfaceschema-2026-09-11.md
pairing-2026-09-09.md
payload-2026-09-12.md
realm-2026-09-09.md
store-2026-09-08.md
---
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
---
mule/queue

→ Read .mule/reviews/config-2026-09-13.md
Style matched. Writing the review file and appending the task lines.

← Write .mule/reviews/auth-2026-09-14.md
Wrote file successfully.

← Edit .mule/reviews/auth-2026-09-14.md
Index: /root/astrate-mule/.mule/reviews/auth-2026-09-14.md
===================================================================
--- /root/astrate-mule/.mule/reviews/auth-2026-09-14.md
+++ /root/astrate-mule/.mule/reviews/auth-2026-09-14.md
@@ -54,9 +54,9 @@
 drops the fallback and panics on `NewCache(0)` would pass every test.
 
 ## What I decided not to propose
 
-- **The a_ch gate on the Astrate-native live socket matheres the verb field as a regex
+- **The a_ch gate on the Astrate-native live socket matches the verb field as a regex
   against `GET`** (`internal/appengine/stream/ws.go:41` → `RequireRealm(ClaimChannels)`
   → `Authorizes(ClaimChannels, "GET", "socket")`). Consequence: the canonical upstream
   Channels grants `JOIN::.*`/`WATCH::.*` can never open the socket, while `.*::.*` — which
   the measured a_ch grammar says authorizes nothing on the Phoenix Channels surface —

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -166,4 +166,8 @@
 - [x] config-fail-loud-dev-mode [auto]: `applyEnv` silently discards a malformed `ASTRATE_MQTT_INSECURE_DEV_MODE` value (internal/config/config.go:255-257) — `strconv.ParseBool` error is discarded, the field becomes false, and TLS requirements are silently re-enabled (or a misleading "tls_cert_file required" error surfaces via validate). Return the parse error on bad input; add a failing test that asserts `ASTRATE_MQTT_INSECURE_DEV_MODE=yes` causes Load to error.
 - [x] config-key-resolvers-test [auto]: add container-free unit tests for `HousekeepingKeys()` and `RealmJWTPublicKey()` (internal/config/config.go:383-411) — both file-reading helpers have zero test coverage. Cover: happy path (inline PEM + file PEM → concatenated/trimmed), file-not-found error (assert error message includes the path), inline-over-file priority (`RealmJWTPublicKey` returns the inline block when both are set), and TrimSpace behaviour (leading/trailing newlines stripped). `t.TempDir()` + `os.WriteFile` for fixtures; no Docker.
 - [x] docs-sync-hk-delete-gating-responses [auto]: add the missing `405` and `422` responses to `DELETE /housekeeping/v1/realms/{realm}` in docs/api/astarte_housekeeping_api.yaml — `writeError` answers 405 "Realm deletion disabled" for `ErrDeletionDisabled` and 422 `{"errors":{"error_name":["connected_devices_present"]}}` for `ErrConnectedDevicesPresent` (internal/housekeeping/http.go:233-249), both reachable from `Service.DeleteRealm` (service.go:248-266), but the spec lists only 204/401/403/404/500 (yaml:198-208). Its sole `405` reference sits on PATCH (yaml:182-183), which can never return it (`UpdateRealm`, service.go:196-216, has no deletion-gating path) — re-point 405 to DELETE (reusing the MethodNotAllowed component) and add 422 with the FieldErrors shape and a connected_devices_present example. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-hk-validation-example [auto]: fix the `422` ValidationError example in docs/api/astarte_housekeeping_api.yaml (yaml:417) — it shows `detail: "housekeeping: validation failed: realm name must not be empty"`, but `validationDetail` strips the `housekeeping: validation failed: ` prefix on the wire (internal/housekeeping/http.go:251-258) and the real blank-name message is `"realm_name can't be blank"` (service.go:138). Change the example to `detail: "realm_name can't be blank"`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+
+- [ ] auth-iat-not-required-test [auto]: pin in internal/auth/jwt_test.go that a token with a *future* `iat` still verifies — the "`iat` is not required — upstream parity" rule (jwt.go:104-105, no `jwt.WithIssuedAt` in the parser, jwt.go:115-118) has zero coverage, so a regression that starts validating `iat` or `WithIssuedAt` would silently reject tokens from a skewed-clock issuer and pass every existing suite. Sign RS256 with `iat` = now+1h and no `exp`/`nbf`, assert `Verify` succeeds.
+- [ ] auth-empty-claims-403-test [auto]: pin in internal/auth/middleware_test.go that a validly signed token carrying NO Astarte claims is `403` on a guarded route, not `401` — the authenticated-but-grant-less distinction (middleware.go:149-151, `Authorizes` on a nil grants map → false → `WriteForbidden`) is currently untested: every 403 fixture has a claim on the wrong surface, so a regression mapping grant-less tokens to 401 (or skipping the forbidden path) passes the suite. Sign RS256 with empty `jwt.MapClaims`, hit the pairing route, expect 403 + golden `envelope_403.json`.
+- [ ] auth-cache-default-size-test [auto]: pin in internal/auth/cache_test.go that `NewCache(0)` and `NewCache(-1)` fall back to `DefaultCacheSize` (cache.go:29-31) — the `size < 1` branch exists purely so `lru.New` cannot panic and has zero coverage; a regression dropping the fallback (panic on `lru.New` with a non-tick size, or a silently tiny cache) would pass every test. Construct with both sizes, assert no panic and a working `Verify` with `lru.Len() == 0`.

Now appending the three task lines to `.mule/todo.md`.
Done. Reviewed `internal/auth` (its first code-review pass), a sibling of the well-trodden surfaces; found a tightly-written package with no category-1 issues, three untested rule pins.

Done: code-review of internal/auth — 3 category-2 rule-pin tasks proposed, review file written
Files: .mule/reviews/auth-2026-09-14.md, .mule/todo.md
Verified: read all 8 auth files + 4 test files, traced middleware callers across all surfaces -> no category-1 finding
Unsure: nothing
Follow-ups: auth-iat-not-required-test, auth-empty-claims-403-test, auth-cache-default-size-test (appended to .mule/todo.md)
