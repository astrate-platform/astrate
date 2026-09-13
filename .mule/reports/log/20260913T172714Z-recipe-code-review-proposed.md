slug: recipe-code-review
verdict: proposed
at:  a609274
ran: 2026-09-13T17:27:14Z on DietPi in 236s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md

$ ls internal/ pkg/ 2>/dev/null
internal/:
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

pkg/:
astarteapi
deviceid
interfaceschema
payload

$ ls .mule/reviews/ 2>/dev/null; echo "---log tail---"; tail -40 .mule/log.md 2>/dev/null
appengine-2026-09-07.md
broker-2026-09-04.md
engine-2026-09-06.md
flow-2026-09-05.md
flowapi-2026-09-10.md
interfaceschema-2026-09-11.md
pairing-2026-09-09.md
payload-2026-09-12.md
realm-2026-09-09.md
store-2026-09-08.md
---log tail---
| 2026-09-07 | appengine-snapshot-ignores-query-params | blocked | 327s | wrote nothing |
| 2026-09-07 | appengine-group-token-roundtrip-test | blocked | 306s | lint failed: internal/appengine/groups_token_test.go:14:5: redefines-builtin-id: redefinition of the built-in function max (revive) |
| 2026-09-07 | docs-sync-appengine-data-output-params | done | 215s | 262fdd5 |
| 2026-09-08 | probe-required-mapping-flag | done | 426s | db9d06e |
| 2026-09-08 | compat-note-v14-rc | blocked | 326s | wrote nothing |
| 2026-09-08 | detailed-listing-required-encrypted-flags | done | 159s | 63b629e |
| 2026-09-08 | store-alias-values-taken-test | done | 85s | a07bd6a |
| 2026-09-08 | store-latest-individual-test | done | 134s | 741794c |
| 2026-09-08 | docs-sync-appengine-device-status-schema | done | 255s | 00da133 |
| 2026-09-08 | docs-sync-appengine-data-set-422 | blocked | 251s | wrote nothing |
| 2026-09-08 | docs-sync-appengine-downsample-min | done | 86s | db506f9 |
| 2026-09-09 | pairing-bearer-secret-test | done | 90s | 32544c0 |
| 2026-09-09 | pairing-remoteip-fallback-test | done | 469s | 3d6e790 |
| 2026-09-09 | docs-sync-hk-patch-endpoint | done | 403s | 9fc9f33 |
| 2026-09-09 | docs-sync-hk-retention-field | done | 142s | d6d7c7f |
| 2026-09-09 | realm-pure-helper-tests | done | 403s | 39eb70a |
| 2026-09-09 | docs-sync-native-compat-health-503 | done | 202s | 656e21e |
| 2026-09-09 | docs-sync-native-compat-version-endpoints | done | 480s | d148f16 |
| 2026-09-10 | flowapi-autorestart-terminal-failure | done | 251s | 8ce04d2 |
| 2026-09-10 | flowapi-validationdetail-test | done | 121s | 832d8fb |
| 2026-09-10 | flowapi-autorestart-default-test | done | 148s | e272cfa |
| 2026-09-10 | docs-sync-pairing-status-enum | done | 89s | 2f6176c |
| 2026-09-10 | docs-sync-pairing-version-endpoint | done | 95s | c867b2c |
| 2026-09-11 | docs-sync-rm-delete-interface-status | done | 123s | b37d11f |
| 2026-09-11 | docs-sync-rm-mapping-required-encrypted | done | 137s | 5f8cb25 |
| 2026-09-11 | docs-sync-rm-put-auth-422 | done | 145s | 8554720 |
| 2026-09-11 | docs-sync-rm-version-example | blocked | 97s | opencode exited 1 |
| 2026-09-11 | interfaceschema-compat-attrs-flip-test | done | 329s | 4fe1f4f |
| 2026-09-11 | interfaceschema-object-attrs-uniformity | done | 344s | c6c8362 |
| 2026-09-11 | interfaceschema-malformed-placeholder-fixtures | done | 130s | f6a44a2 |
| 2026-09-11 | interfaceschema-enum-roundtrip-coverage | done | 185s | f955958 |
| 2026-09-11 | interfaceschema-trie-interior-literal-param | done | 121s | 705795f |
| 2026-09-12 | docs-sync-ae-write-405 | done | 400s | 9ed08c3 |
| 2026-09-12 | docs-sync-ae-write-value-422 | done | 338s | 17f6f70 |
| 2026-09-12 | docs-sync-ae-post-groups-409 | done | 128s | d6865a1 |
| 2026-09-12 | payload-longinteger-fraction-quantize | done | 356s | d7babc6 |
| 2026-09-12 | payload-json-malformed-t-not-tolerated | done | 99s | 5b56557 |
| 2026-09-12 | docs-sync-ae-post-group-devices-409 | done | 179s | cbf81af |
| 2026-09-12 | docs-sync-ae-patch-merge-patch-media-type | done | 137s | 64361d9 |
| 2026-09-13 | store-todo-lttb-stale | done | 66s | 829f408 |

$ for p in auth config housekeeping httpx observability swagger testutil; do echo "$p: $(rg -c '' --glob 'internal/'$p'/**/*.go' internal/$p 2>/dev/null | awk -F: '{s+=$2} END {print s+0}')"; done; echo "--- pkg ---"; for p in astarteapi deviceid; do echo "$p: $(rg -c '' --glob 'pkg/'$p'/**/*.go' pkg/$p 2>/dev/null | awk -F: '{s+=$2} END {print s+0}')"; done
auth: 1532
config: 818
housekeeping: 1311
httpx: 357
observability: 398
swagger: 171
testutil: 697
--- pkg ---
astarteapi: 485
deviceid: 358

$ for p in httpx config observability deviceid swagger astarteapi; do echo "=== $p ==="; wc -l internal/$p/*.go pkg/$p/*.go 2>/dev/null | sort -rn | head -20; done
=== httpx ===
  357 total
  139 internal/httpx/notfound_test.go
  113 internal/httpx/cors_test.go
   56 internal/httpx/notfound.go
   49 internal/httpx/cors.go
=== config ===
  818 total
  411 internal/config/config.go
  407 internal/config/config_test.go
=== observability ===
  398 total
  128 internal/observability/observability_test.go
   83 internal/observability/health.go
   80 internal/observability/metrics.go
   55 internal/observability/compat_version_test.go
   52 internal/observability/compat.go
=== deviceid ===
  358 total
  226 pkg/deviceid/deviceid_test.go
  132 pkg/deviceid/deviceid.go
=== swagger ===
 171 total
 128 internal/swagger/swagger_test.go
  43 internal/swagger/swagger.go
=== astarteapi ===
  485 total
  249 pkg/astarteapi/envelope.go
  236 pkg/astarteapi/envelope_test.go

Picking `internal/config` (moderate size, not reviewed recently). Reading the package now.
→ Read internal/config/config.go
→ Read internal/config/config_test.go

$ rg -n "RegisterRate|RegisterBurst|CredentialsRate|CredentialsBurst" internal/ --glob '*.go' | grep -v _test.go
internal/config/config.go:102:	RegisterRate      float64  `toml:"register_rate"`
internal/config/config.go:103:	RegisterBurst     int      `toml:"register_burst"`
internal/config/config.go:104:	CredentialsRate   float64  `toml:"credentials_rate"`
internal/config/config.go:105:	CredentialsBurst  int      `toml:"credentials_burst"`
internal/config/config.go:194:			RegisterRate: 5, RegisterBurst: 10,
internal/config/config.go:195:			CredentialsRate: 5, CredentialsBurst: 10,
internal/pairing/http.go:23:	DefaultRegisterRate     = 5.0
internal/pairing/http.go:24:	DefaultRegisterBurst    = 20
internal/pairing/http.go:25:	DefaultCredentialsRate  = 1.0
internal/pairing/http.go:26:	DefaultCredentialsBurst = 5
internal/pairing/http.go:35:	RegisterRate     float64
internal/pairing/http.go:36:	RegisterBurst    int
internal/pairing/http.go:37:	CredentialsRate  float64
internal/pairing/http.go:38:	CredentialsBurst int
internal/pairing/http.go:54:	if cfg.RegisterRate == 0 {
internal/pairing/http.go:55:		cfg.RegisterRate = DefaultRegisterRate
internal/pairing/http.go:57:	if cfg.RegisterBurst == 0 {
internal/pairing/http.go:58:		cfg.RegisterBurst = DefaultRegisterBurst
internal/pairing/http.go:60:	if cfg.CredentialsRate == 0 {
internal/pairing/http.go:61:		cfg.CredentialsRate = DefaultCredentialsRate
internal/pairing/http.go:63:	if cfg.CredentialsBurst == 0 {
internal/pairing/http.go:64:		cfg.CredentialsBurst = DefaultCredentialsBurst
internal/pairing/http.go:69:		regLimiter:   NewLimiter(cfg.RegisterRate, cfg.RegisterBurst),
internal/pairing/http.go:70:		credLimiter:  NewLimiter(cfg.CredentialsRate, cfg.CredentialsBurst),

Now let me trace how the unvalidated numeric fields are consumed before deciding what to propose.
$ rg -n "CertTTL|BatchMaxRows|MaxPayloadBytes|ShardQueue|Storage.Retention|BatchMaxWait" internal/ --glob '*.go' | grep -v _test.go | head -40
internal/config/config.go:92:	ShardQueue      int      `toml:"shard_queue"`
internal/config/config.go:93:	BatchMaxRows    int      `toml:"batch_max_rows"`
internal/config/config.go:94:	BatchMaxWait    Duration `toml:"batch_max_wait"`
internal/config/config.go:95:	MaxPayloadBytes int      `toml:"max_payload_bytes"`
internal/config/config.go:100:	CertTTL           Duration `toml:"cert_ttl"`
internal/config/config.go:187:			ShardQueue:      4096,
internal/config/config.go:188:			BatchMaxRows:    64,
internal/config/config.go:189:			BatchMaxWait:    Duration(50 * time.Millisecond),
internal/config/config.go:190:			MaxPayloadBytes: 64 * 1024,
internal/config/config.go:193:			CertTTL:      Duration(30 * 24 * time.Hour),
internal/pairing/service.go:107:	// CertTTL is the client certificate validity; zero selects
internal/pairing/service.go:108:	// ca.DefaultCertTTL (30 days).
internal/pairing/service.go:109:	CertTTL time.Duration
internal/pairing/service.go:147:	if cfg.CertTTL <= 0 {
internal/pairing/service.go:148:		cfg.CertTTL = ca.DefaultCertTTL
internal/pairing/service.go:260:	certPEM, serial, aki, err := realmCA.SignCSR(csrPEM, realmName, dev.ID.String(), s.cfg.CertTTL)
internal/pairing/ca/ca.go:28:	// DefaultCertTTL is the client certificate validity
internal/pairing/ca/ca.go:30:	DefaultCertTTL = 30 * 24 * time.Hour
internal/pairing/ca/ca.go:160://     NotAfter; a non-positive ttl selects DefaultCertTTL;
internal/pairing/ca/ca.go:180:		ttl = DefaultCertTTL
internal/engine/batch.go:22:// until BatchMaxRows or BatchMaxWait (the shard loop owns the timer), then
internal/engine/router.go:26:	// DefaultShardQueue is the default per-shard channel capacity.
internal/engine/router.go:27:	DefaultShardQueue = 4096
internal/engine/router.go:28:	// DefaultBatchMaxRows flushes a shard batch when it reaches this size.
internal/engine/router.go:29:	DefaultBatchMaxRows = 64
internal/engine/router.go:30:	// DefaultBatchMaxWait flushes a non-empty shard batch after this delay.
internal/engine/router.go:31:	DefaultBatchMaxWait = 50 * time.Millisecond
internal/engine/router.go:39:	// ShardQueue is the per-shard channel capacity (default
internal/engine/router.go:40:	// DefaultShardQueue). A full shard blocks QoS >= 1 submits (deferred-ack
internal/engine/router.go:42:	ShardQueue int
internal/engine/router.go:43:	// BatchMaxRows is the micro-batch row cap (default DefaultBatchMaxRows).
internal/engine/router.go:44:	BatchMaxRows int
internal/engine/router.go:45:	// BatchMaxWait is the micro-batch time cap (default DefaultBatchMaxWait).
internal/engine/router.go:46:	BatchMaxWait time.Duration
internal/engine/router.go:47:	// MaxPayloadBytes caps accepted data payload size for both formats
internal/engine/router.go:49:	MaxPayloadBytes int
internal/engine/router.go:137:	// timer drives the BatchMaxWait flush; armed tracks whether timer.C is
internal/engine/router.go:152:	if cfg.ShardQueue == 0 {
internal/engine/router.go:153:		cfg.ShardQueue = DefaultShardQueue
internal/engine/router.go:155:	if cfg.BatchMaxRows == 0 {

$ rg -n "func NewLimiter" internal/pairing/ -A 15
internal/pairing/ratelimit.go:35:func NewLimiter(rate float64, burst int) *Limiter {
internal/pairing/ratelimit.go-36-	if rate <= 0 {
internal/pairing/ratelimit.go-37-		rate = 1.0 / 60
internal/pairing/ratelimit.go-38-	}
internal/pairing/ratelimit.go-39-	if burst < 1 {
internal/pairing/ratelimit.go-40-		burst = 1
internal/pairing/ratelimit.go-41-	}
internal/pairing/ratelimit.go-42-	// lru.New only fails for sizes < 1; limiterMaxKeys is a positive const.
internal/pairing/ratelimit.go-43-	cache, _ := lru.New[string, *bucket](limiterMaxKeys)
internal/pairing/ratelimit.go-44-	return &Limiter{
internal/pairing/ratelimit.go-45-		buckets: cache,
internal/pairing/ratelimit.go-46-		rate:    rate,
internal/pairing/ratelimit.go-47-		burst:   float64(burst),
internal/pairing/ratelimit.go-48-		now:     time.Now,
internal/pairing/ratelimit.go-49-	}
internal/pairing/ratelimit.go-50-}

$ rg -n "negative|must be positive|>\s*=\s*0" internal/engine/router.go internal/config/config.go | head
internal/engine/router.go:162:		return nil, fmt.Errorf("engine: negative configuration value: %+v", cfg)
internal/config/config.go:280:			return fmt.Errorf("config: housekeeping.default_datastream_maximum_storage_retention %q must be a non-negative integer", retention)
internal/config/config.go:328:		return fmt.Errorf("config: engine.shards must be positive")

→ Read .mule/todo.md
→ Read .mule/todo.md [offset=162]

← Write .mule/reviews/config-2026-09-13.md
Wrote file successfully.

Now I have all the context. Writing the review file and appending task lines.
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -161,4 +161,7 @@
 - [x] payload-longinteger-fraction-quantize [auto]: `jsonInt64` (pkg/payload/json.go:255-271) accepts fractional longinteger tokens that the ParseFloat fallback silently rounds onto an integer inside its own documented-safety window |x| < 2^53 (json.go:250-254 claims values there never round, "silently corrupting the value") — `{"v":4503599627370496.5}` decodes to 4503599627370496, `{"v":6000000000000000.9}` to 6000000000000001, and `{"v":5e-400}` underflows to 0, all without error, while the contract says fraction forms are "accepted only when integral". Verify the token is an exact decimal integer (`math/big.Rat` SetString + IsInt) on the float fallback — keeps `5.0`, `1e3`, `0.5e1` accepted, rejects the rounding cases as ReasonTypeMismatch; plain-integer ParseInt fast path and the > 2^53 window check unchanged (BSON unaffected: longinteger accepts int64/int32 only). Add failing rows to TestJSONProfile.
 - [x] payload-json-malformed-t-not-tolerated [auto]: pin in pkg/payload/payload_test.go that a malformed `t` is rejected with ReasonBadTimestamp even on a non-explicit mapping — `decodeJSONEnvelope` (json.go:50-66) parses `t` before the facade learns the explicit-timestamp policy (payload.go:52-54), so only a *valid* `t` is "tolerated and ignored" while `{"v":1,"t":"yesterday"}` with explicit=false is still a rejection; only the valid-t-ignored row and the BSON "t wrong type" twin are asserted today, so a regression implementing the literal leniency reading (ignore a bad t on non-explicit mappings) passes every suite. Add one table row to TestJSONProfile.
 - [x] docs-sync-ae-post-group-devices-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml — re-adding a member the group already holds answers `409` "Device already in group" (`Service.AddGroupDevice` returns `ErrDeviceAlreadyInGroup`, internal/appengine/service.go:636-638, mapped at http.go:662-663), but the spec (yaml:967-977) lists only 201/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-ae-patch-merge-patch-media-type [auto]: fix the three PATCH request bodies in docs/api/astarte_appengine_api.yaml — PATCH /devices/{device} (yaml:176), PATCH /devices-by-alias/{alias} (yaml:260), PATCH /groups/{group}/devices/{device} (yaml:1023) declare `application/json; charset=utf-8`, but all three handlers enforce `Content-Type: application/merge-patch+json` exactly and reject every other value — 500 on the device and group paths (internal/appengine/http.go:209-211, 521-523), 400 on the alias path (http.go:230-232). Change the requestBody media type to `application/merge-patch+json` on all three. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
+- [ ] config-fail-loud-dev-mode [auto]: `applyEnv` silently discards a malformed `ASTRATE_MQTT_INSECURE_DEV_MODE` value (internal/config/config.go:255-257) — `strconv.ParseBool` error is discarded, the field becomes false, and TLS requirements are silently re-enabled (or a misleading "tls_cert_file required" error surfaces via validate). Return the parse error on bad input; add a failing test that asserts `ASTRATE_MQTT_INSECURE_DEV_MODE=yes` causes Load to error.
+- [ ] config-key-resolvers-test [auto]: add container-free unit tests for `HousekeepingKeys()` and `RealmJWTPublicKey()` (internal/config/config.go:383-411) — both file-reading helpers have zero test coverage. Cover: happy path (inline PEM + file PEM → concatenated/trimmed), file-not-found error (assert error message includes the path), inline-over-file priority (`RealmJWTPublicKey` returns the inline block when both are set), and TrimSpace behaviour (leading/trailing newlines stripped). `t.TempDir()` + `os.WriteFile` for fixtures; no Docker.

$ tail -5 .mule/todo.md; echo "---"; git status --short 2>&1 | head
- [x] docs-sync-ae-post-group-devices-409 [auto]: add the missing `409` response to `POST /appengine/v1/{realm}/groups/{group}/devices` in docs/api/astarte_appengine_api.yaml — re-adding a member the group already holds answers `409` "Device already in group" (`Service.AddGroupDevice` returns `ErrDeviceAlreadyInGroup`, internal/appengine/service.go:636-638, mapped at http.go:662-663), but the spec (yaml:967-977) lists only 201/400/401/404/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-ae-patch-merge-patch-media-type [auto]: fix the three PATCH request bodies in docs/api/astarte_appengine_api.yaml — PATCH /devices/{device} (yaml:176), PATCH /devices-by-alias/{alias} (yaml:260), PATCH /groups/{group}/devices/{device} (yaml:1023) declare `application/json; charset=utf-8`, but all three handlers enforce `Content-Type: application/merge-patch+json` exactly and reject every other value — 500 on the device and group paths (internal/appengine/http.go:209-211, 521-523), 400 on the alias path (http.go:230-232). Change the requestBody media type to `application/merge-patch+json` on all three. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] config-fail-loud-engine-shards [auto]: `applyEnv` silently discards a malformed `ASTRATE_ENGINE_SHARDS` value (internal/config/config.go:258-262) — `strconv.Atoi` error is discarded and the TOML/default shard count stays in force, unlike the adjacent fail-loud `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED` (lines 288-297). Return the parse error on bad input; add a failing test that asserts `ASTRATE_ENGINE_SHARDS=banana` causes Load to error.
- [ ] config-fail-loud-dev-mode [auto]: `applyEnv` silently discards a malformed `ASTRATE_MQTT_INSECURE_DEV_MODE` value (internal/config/config.go:255-257) — `strconv.ParseBool` error is discarded, the field becomes false, and TLS requirements are silently re-enabled (or a misleading "tls_cert_file required" error surfaces via validate). Return the parse error on bad input; add a failing test that asserts `ASTRATE_MQTT_INSECURE_DEV_MODE=yes` causes Load to error.
- [ ] config-key-resolvers-test [auto]: add container-free unit tests for `HousekeepingKeys()` and `RealmJWTPublicKey()` (internal/config/config.go:383-411) — both file-reading helpers have zero test coverage. Cover: happy path (inline PEM + file PEM → concatenated/trimmed), file-not-found error (assert error message includes the path), inline-over-file priority (`RealmJWTPublicKey` returns the inline block when both are set), and TrimSpace behaviour (leading/trailing newlines stripped). `t.TempDir()` + `os.WriteFile` for fixtures; no Docker.
---
 M .mule/todo.md
?? .mule/reviews/config-2026-09-13.md

Done: reviewed internal/config and proposed work.

Files:
- `.mule/reviews/config-2026-09-13.md` — what I read, 3 findings, what I chose not to propose and why
- `.mule/todo.md` — 3 task lines appended (capped at 5 per recipe)

Findings (all in `internal/config/config.go`):
1. `ASTRATE_ENGINE_SHARDS` parse error swallowed (lines 258-262) — garbage value invisible, shard count silently keeps TOML/default, unlike the adjacent fail-loud realm-deletion var.
2. `ASTRATE_MQTT_INSECURE_DEV_MODE` parse error swallowed (lines 255-257) — garbage value silently flips the field to false, re-enabling TLS requirements.
3. `HousekeepingKeys()`/`RealmJWTPublicKey()` file-read helpers (lines 383-411) — zero test coverage for happy, error, and trim paths.

Other candidate areas (negative pairing rates, negative engine fields, CORS ftp:// scheme) dismissed after tracing callers — all already guarded downstream.

```
Done: code review of internal/config; proposed 3 tasks
Files: .mule/reviews/config-2026-09-13.md, .mule/todo.md
Verified: none needed (proposal job, no code changed)
Unsure: nothing
Follow-ups: config-fail-loud-engine-shards, config-fail-loud-dev-mode, config-key-resolvers-test (appended)
```
