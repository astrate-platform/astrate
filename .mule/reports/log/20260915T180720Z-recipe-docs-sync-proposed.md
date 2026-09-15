slug: recipe-docs-sync
verdict: proposed
at:  69e44b2
ran: 2026-09-15T18:07:20Z on DietPi in 660s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
I'll pick one surface and compare documented vs actual endpoints. Let me start by finding the available surfaces.
✱ Glob "docs/api/astrate_*_api.yaml" 1 match

$ ls internal/
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

→ Read docs/api/astrate_native_api.yaml
Only one API YAML exists: `astrate_native_api.yaml`. Let me find the actual handlers.

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*native*' 2>/dev/null || rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*realm*' 2>/dev/null
(no output)

$ rg -n 'HandleFunc|\.Get\(|\.Post\(|\.Put\(|\.Delete\(|\.Patch\(' internal/ --glob '*.go' | grep -i 'health\|readiness\|metrics\|version\|socket\|websocket\|stream' | head -60
internal/pairing/http.go:81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/observability/compat.go:20:	mux.HandleFunc("GET /"+service+"/health", func(w http.ResponseWriter, r *http.Request) {
internal/observability/compat.go:51:	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
internal/observability/health.go:43:	mux.HandleFunc("GET /astrate/v1/health", h.handleHealth)
internal/observability/health.go:44:	mux.HandleFunc("GET /astrate/v1/readiness", h.handleReadiness)
internal/observability/compat_version_test.go:19:		resp, err := http.Get(srv.URL + "/" + svc + "/version")
internal/observability/compat_version_test.go:31:		if ct := resp.Header.Get("Content-Type"); ct == "" {
internal/observability/compat_version_test.go:39:	mux.HandleFunc("GET /v1/{realm}/version", VersionHandler("x"))
internal/observability/compat_version_test.go:43:	resp, err := http.Get(srv.URL + "/v1/somerealm/version")
internal/appengine/stream/ws_test.go:81:		if ct := rec.Header().Get("Content-Type"); ct != "text/event-stream" {
internal/appengine/stream/ws.go:72:		DeviceID:  r.URL.Query().Get("device_id"),
internal/appengine/stream/ws.go:73:		Interface: r.URL.Query().Get("interface"),
internal/appengine/stream/ws.go:152:	return r.Header.Get("Accept") == "text/event-stream" || r.URL.Query().Get("transport") == "sse"

Let me find where the native API routes are actually registered.
$ rg -rn 'HandleFunc|\.Get\(|\.Post\(' internal/ --glob '*.go' | grep -v '_test.go' | head -80
internal/flow/blocks/httpblocks.go:		Subtype:   mediaType(resp.Headern"Content-Type")),
internal/flow/blocks/httpblocks.go:	if req.Headern"Content-Type") == "" {
internal/swagger/swagger.go:	mux.n("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
internal/httpx/cors.go:			origin := r.Headern"Origin")
internal/httpx/cors.go:			if r.Method == http.MethodOptions && r.Headern"Access-Control-Request-Method") != "" {
internal/appengine/stream/ws.go:		DeviceID:  r.URL.Query()n"device_id"),
internal/appengine/stream/ws.go:		Interface: r.URL.Query()n"interface"),
internal/appengine/stream/ws.go:	return r.Headern"Accept") == "text/event-stream" || r.URL.Query()n"transport") == "sse"
internal/appengine/http.go:	limit, _ := strconv.Atoi(qn"limit"))
internal/appengine/http.go:	details, _ := strconv.ParseBool(qn"details"))
internal/appengine/http.go:	page, err := a.svc.ListDevices(r.Context(), r.PathValue("realm"), qn"from_token"), limit, details)
internal/appengine/http.go:		details := qn"details")
internal/appengine/http.go:		limit := qn"limit")
internal/appengine/http.go:	if r.Headern"Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:	if r.Headern"Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:	details, _ := strconv.ParseBool(qn"details"))
internal/appengine/http.go:	limit, _ := strconv.Atoi(qn"limit"))
internal/appengine/http.go:		details, qn"from_token"), limit)
internal/appengine/http.go:	limit := qn"limit")
internal/appengine/http.go:	if d := qn"details"); d != "" {
internal/appengine/http.go:	if r.Headern"Content-Type") != "application/merge-patch+json" {
internal/appengine/http.go:		if v := qnp.name); v != "" {
internal/appengine/http.go:	if qn"since") != "" && qn"since_after") != "" {
internal/appengine/http.go:	if v := qn"limit"); v != "" {
internal/appengine/http.go:	if v := qn"downsample_to"); v != "" {
internal/appengine/http.go:	switch v := qn"format"); v {
internal/appengine/http.go:		switch v := qnp.name); v {
internal/appengine/http.go:	if v := qn"retrieve_metadata"); v != "" {
internal/appengine/http.go:	opts.DownsampleKey = qn"downsample_key")
internal/appengine/http.go:	opts.Descending = qn"sort") != "ascending"
internal/pairing/http.go:	mux.n("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/pairing/http.go:	mux.n("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
internal/pairing/http.go:	mux.n("GET /pairing/v1/{realm}/devices/{deviceID}",
internal/pairing/http.go:	mux.n("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
internal/pairing/http.go:	header := r.Headern"Authorization")
internal/appengine/channels/ws.go:	realm := r.URL.Query()n"realm")
internal/appengine/channels/ws.go:	token := r.URL.Query()n"token")
internal/pairing/ratelimit.go:	b, ok := l.bucketsnkey)
internal/auth/middleware.go:	header := r.Headern"Authorization")
internal/broker/sessionstore.go:		if prev := bnkey); prev != nil {
internal/broker/sessionstore.go:	old, ok := s.srv.Clientsncl.ID)
internal/auth/cache.go:	if tok, ok := c.lrunkey); ok {
internal/realm/http.go:	if r.URL.Query()n"detailed") == "true" {
internal/observability/health.go:	mux.n("GET /astrate/v1/health", h.handleHealth)
internal/observability/health.go:	mux.n("GET /astrate/v1/readiness", h.handleReadiness)
internal/observability/compat.go:	mux.n("GET /"+service+"/health", func(w http.ResponseWriter, r *http.Request) {
internal/observability/compat.go:	mux.n("GET /"+service+"/version", VersionHandler(version))

$ rg -n 'Handle\((mux|m)|HandleFunc|/socket/websocket|astrate/v1/\{realm\}/socket|/appengine/socket' internal/ --glob '*.go' | grep -v _test | head -40
internal/swagger/swagger.go:20:	mux.HandleFunc("GET /swagger", func(w http.ResponseWriter, r *http.Request) {
internal/observability/compat.go:20:	mux.HandleFunc("GET /"+service+"/health", func(w http.ResponseWriter, r *http.Request) {
internal/observability/compat.go:51:	mux.HandleFunc("GET /"+service+"/version", VersionHandler(version))
internal/observability/health.go:43:	mux.HandleFunc("GET /astrate/v1/health", h.handleHealth)
internal/observability/health.go:44:	mux.HandleFunc("GET /astrate/v1/readiness", h.handleReadiness)
internal/pairing/http.go:81:	mux.HandleFunc("GET /pairing/v1/{realm}/health", a.handleHealth)
internal/pairing/http.go:86:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials",
internal/pairing/http.go:88:	mux.HandleFunc("GET /pairing/v1/{realm}/devices/{deviceID}",
internal/pairing/http.go:90:	mux.HandleFunc("POST /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify",
internal/appengine/stream/ws.go:3:// at /astrate/v1/{realm}/socket fed by the engine's in-process fan-out bus.
internal/appengine/stream/ws.go:46:	mux.Handle("GET /astrate/v1/{realm}/socket", a.require(http.HandlerFunc(a.handle)))
internal/appengine/channels/ws.go:48:	mux.Handle("GET /appengine/v1/socket/websocket", http.HandlerFunc(a.handle))
internal/appengine/channels/protocol.go:2:// served at /appengine/v1/socket/websocket. It is the compatibility counterpart

Let me find the actual route registrations more precisely.
$ rg -n 'socket' internal/ --glob '*.go' -l | grep -v _test
internal/engine/engine.go
internal/engine/triggers/validation.go
internal/appengine/stream/ws.go
internal/engine/triggers/events.go
internal/appengine/channels/room.go
internal/appengine/channels/protocol.go
internal/appengine/channels/ws.go
internal/pairing/http.go
internal/appengine/service.go
internal/auth/claims.go

$ rg -n 'metrics|promhttp|/astrate/v1/metrics' internal/ cmd/ --glob '*.go' | grep -v _test | head
cmd/astrate/main.go:111:	metrics := observability.NewMetrics()
cmd/astrate/main.go:125:		Registerer:      metrics.Registerer(),
cmd/astrate/main.go:184:	handler, hkSvc, err := mountAPIs(cfg, st, e, b, sealer, metrics, flowSvc, pairer, log)
cmd/astrate/main.go:350:func mountAPIs(cfg config.Config, st *store.Store, e *engine.Engine, b *broker.Broker, sealer *store.KeySealer, metrics *observability.Metrics, flowSvc *flowapi.Service, pairer *pairing.Service, log *slog.Logger) (http.Handler, *housekeeping.Service, error) {
cmd/astrate/main.go:400:	metrics.RegisterBrokerSessions(func() float64 { return float64(b.SessionCount()) })
cmd/astrate/main.go:401:	metrics.RegisterDBPool(func() observability.DBPoolStats {
cmd/astrate/main.go:410:	health := observability.NewHealth(metrics.Handler())
internal/flow/block.go:21:	// Name returns a human-readable label for metrics and logging.
internal/flow/factory.go:62:// node name and should be returned by Block.Name() for metrics and logging.
internal/engine/data.go:86:// String returns the stable snake_case metrics label.

Let me find the metrics registration and check the socket handlers.
→ Read internal/observability/compat.go
→ Read cmd/astrate/main.go [offset=350, limit=120]
→ Read internal/observability/health.go
All 12 documented paths are present. Now let me verify 3-4 endpoints' status codes against handlers, starting with the two WebSocket endpoints (most detail documented).
→ Read internal/appengine/stream/ws.go
→ Read internal/appengine/channels/ws.go
→ Read internal/auth/middleware.go
Let me check the auth middleware (extraction method + 401 shape) since the docs claim header-or-`?token=` auth for the native socket.

$ rg -n 'ClaimChannels|a_ch|let claim' internal/auth/ --glob '*.go' | grep -v _test
internal/auth/claims.go:22:	// ClaimChannels authorizes Astarte Channels; Astrate honours it on the
internal/auth/claims.go:23:	// live stream socket (a_ch).
internal/auth/claims.go:24:	ClaimChannels Claim = "a_ch"
internal/auth/claims.go:71:	Channels        stringList `json:"a_ch,omitempty"`
internal/auth/claims.go:135:		ClaimChannels:        c.Channels,
internal/auth/claims.go:174:// silently discards an a_ch entry carrying anything else.
internal/auth/claims.go:180:// AuthorizesChannel reports whether the a_ch claim grants verb (VerbJoin or
internal/auth/claims.go:184:// It deliberately does not reuse Authorizes, because upstream reads a_ch by a
internal/auth/claims.go:188:// Channels socket instead *partitions* the a_ch list by an exact match on the
internal/auth/claims.go:198:	for _, g := range t.grants[ClaimChannels] {

$ rg -n 'token' internal/appengine/stream/ internal/auth/ --glob '*.go' | grep -v _test
internal/auth/middleware.go:19:// tokenContextKey carries the verified *Token in the request context.
internal/auth/middleware.go:20:type tokenContextKey struct{}
internal/auth/middleware.go:22:// TokenFromContext returns the verified token stored by the middleware, if
internal/auth/middleware.go:25:	tok, ok := ctx.Value(tokenContextKey{}).(*Token)
internal/auth/middleware.go:31:// unverifiable token → 401, verified token whose claims do not authorize the
internal/auth/middleware.go:39:// DefaultCacheSize token cache.
internal/auth/middleware.go:46:// bearer token, and matches the claim's authorization strings against the
internal/auth/middleware.go:131:	tokenString, ok := bearerToken(r)
internal/auth/middleware.go:137:	tok, err := m.cache.Verify(tokenString, keysPEM)
internal/auth/middleware.go:153:	ctx := context.WithValue(r.Context(), tokenContextKey{}, tok)
internal/auth/middleware.go:177:// "Bearer", optionally followed by a colon, then whitespace and the token.
internal/auth/middleware.go:191:	token := strings.TrimSpace(rest)
internal/auth/middleware.go:192:	if token == "" {
internal/auth/middleware.go:195:	return token, true
internal/auth/jwt.go:3:// reproduces upstream Astarte's token semantics exactly — asymmetric keys
internal/auth/jwt.go:7:// relative to the realm base — so tokens minted by astartectl and existing
internal/auth/jwt.go:32:// Sentinel errors returned by token verification. All of them map to 401 at
internal/auth/jwt.go:36:	// configured, so no token can possibly verify.
internal/auth/jwt.go:38:	// ErrNoKeyMatched reports that the token signature verified against
internal/auth/jwt.go:39:	// none of the realm's keys (wrong key, tampered token, or disallowed
internal/auth/jwt.go:41:	ErrNoKeyMatched = errors.New("auth: token matches none of the realm keys")
internal/auth/jwt.go:52:// for every token holder.
internal/auth/jwt.go:103:// Verify parses tokenString, verifies its signature against the key set, and
internal/auth/jwt.go:106:// The token verifies if *any* key in the set matches, which is what makes
internal/auth/jwt.go:108:func Verify(tokenString string, keys []crypto.PublicKey, now func() time.Time) (*Token, error) {
internal/auth/jwt.go:123:		_, err := parser.ParseWithClaims(tokenString, claims, func(*jwt.Token) (any, error) {
internal/auth/jwt.go:130:		// through the set. Any other failure (malformed token, disallowed
internal/auth/jwt.go:136:		return nil, fmt.Errorf("auth: invalid token: %w", err)
internal/auth/cache.go:13:// DefaultCacheSize is the verified-token LRU capacity (docs/DESIGN.md §4.2).
internal/auth/cache.go:16:// Cache memoizes verified tokens so the signature check, claim validation,
internal/auth/cache.go:17:// and regex compilation run once per (token, key set) instead of once per
internal/auth/cache.go:18:// request. Entries are keyed by SHA-256 of the token *and* of the key set,
internal/auth/cache.go:19:// so rotating a realm's keys naturally invalidates its cached tokens.
internal/auth/cache.go:41:// cacheKey fingerprints the token together with the PEM key set. Domain
internal/auth/cache.go:42:// separation between token and keys (and between keys) uses a 0x00 byte,
internal/auth/cache.go:43:// which cannot appear in base64url tokens or PEM text.
internal/auth/cache.go:44:func cacheKey(tokenString string, keysPEM []string) [sha256.Size]byte {
internal/auth/cache.go:46:	h.Write([]byte(tokenString))
internal/auth/cache.go:56:// Verify returns the verified Token for tokenString against the realm's PEM
internal/auth/cache.go:57:// key set, from cache when possible. Cache hits still honour `exp`: a token
internal/auth/cache.go:60:func (c *Cache) Verify(tokenString string, keysPEM []string) (*Token, error) {
internal/auth/cache.go:61:	key := cacheKey(tokenString, keysPEM)
internal/auth/cache.go:66:			return nil, fmt.Errorf("auth: cached token: %w", jwt.ErrTokenExpired)
internal/auth/cache.go:75:	tok, err := c.verify(tokenString, keys, c.now)
internal/auth/claims.go:39:// breaks no valid token).
internal/auth/claims.go:99://     simply never matches; it does not invalidate the token.
internal/auth/claims.go:150:// ExpiresAt returns the token expiry and whether one is set (`exp` is
internal/auth/claims.go:159:// Authorizes reports whether the token grants `claim` for the given verb
internal/auth/claims.go:162:// Multiple authorization strings within a claim are OR-ed; a token without

The native socket middleware reads only the `Authorization` header (`bearerToken`), but the YAML claims auth also works via `?token=` query parameter. Let me confirm no `?token=` handling exists for that route.
$ rg -n 'To(v|u)\b|socket|websocket|transport' internal/appengine/stream/ --glob '*.go' | grep -v _test
internal/appengine/stream/ws.go:1:// Package stream is the AppEngine live event socket (docs/DESIGN.md §3.7,
internal/appengine/stream/ws.go:3:// at /astrate/v1/{realm}/socket fed by the engine's in-process fan-out bus.
internal/appengine/stream/ws.go:18:	"github.com/coder/websocket"
internal/appengine/stream/ws.go:24:// subscribeBuffer is the per-connection bus buffer; a slow socket drops
internal/appengine/stream/ws.go:28:// Bus is the subset of the engine fan-out bus the socket consumes.
internal/appengine/stream/ws.go:33:// API is the live-socket HTTP surface, guarded by an a_ch realm JWT.
internal/appengine/stream/ws.go:39:// NewAPI wires the socket to the engine bus and the a_ch middleware.
internal/appengine/stream/ws.go:44:// Mount registers the socket route on mux.
internal/appengine/stream/ws.go:46:	mux.Handle("GET /astrate/v1/{realm}/socket", a.require(http.HandlerFunc(a.handle)))
internal/appengine/stream/ws.go:87:	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
internal/appengine/stream/ws.go:93:	// the client disconnects (the socket is server-push only).
internal/appengine/stream/ws.go:102:				_ = conn.Close(websocket.StatusGoingAway, "bus closed")
internal/appengine/stream/ws.go:109:			if err := conn.Write(ctx, websocket.MessageText, payload); err != nil {
internal/appengine/stream/ws.go:152:	return r.Header.Get("Accept") == "text/event-stream" || r.URL.Query().Get("transport") == "sse"

$ rg -n 'Query\(\)\.Get\("token"\)|query.*token|token=' internal/ --glob '*.go' | grep -v _test
internal/appengine/channels/ws.go:54:	token := r.URL.Query().Get("token")

One drift candidate found. Let me make sure no global middleware elsewhere handles `?token=` for that route, and check for `?token=` anywhere in the request path.
$ rg -n 'astrate_devices_total|RegisterGauge|prometheus\.NewGauge|Metric' internal/observability/ --glob '*.go' | grep -v _test | head -20
internal/observability/metrics.go:20:// Metrics owns the process-wide Prometheus registry.
internal/observability/metrics.go:21:type Metrics struct {
internal/observability/metrics.go:25:// NewMetrics builds a registry preloaded with the Go runtime and process
internal/observability/metrics.go:27:func NewMetrics() *Metrics {
internal/observability/metrics.go:33:	return &Metrics{reg: reg}
internal/observability/metrics.go:38:func (m *Metrics) Registerer() prometheus.Registerer { return m.reg }
internal/observability/metrics.go:41:func (m *Metrics) Handler() http.Handler {
internal/observability/metrics.go:56:func (m *Metrics) RegisterBrokerSessions(fn func() float64) {
internal/observability/metrics.go:57:	m.reg.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
internal/observability/metrics.go:67:func (m *Metrics) RegisterDBPool(fn func() DBPoolStats) {
internal/observability/metrics.go:69:		m.reg.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{

Confirmed: `?token=` is only handled in the Phoenix socket (`channels/ws.go:54`), not in the native socket. The YAML documents `?token=` for the native socket — that's drift. Now let me read the current todo and for-giulio files before appending.
→ Read .mule/todo.md
→ Read .mule/for-giulio.md
→ Read docs/api/astrate_native_api.yaml [offset=290, limit=6]
Only one drift on the astrate-native surface. Let me confirm the YAML's auth claim line numbers precisely and check the end of the todo file.

$ tail -5 .mule/todo.md; echo "---lines---"; wc -l .mule/todo.md
- [x] auth-iat-not-required-test [auto]: pin in internal/auth/jwt_test.go that a token with a *future* `iat` still verifies — the "`iat` is not required — upstream parity" rule (jwt.go:104-105, no `jwt.WithIssuedAt` in the parser, jwt.go:115-118) has zero coverage, so a regression that starts validating `iat` or `WithIssuedAt` would silently reject tokens from a skewed-clock issuer and pass every existing suite. Sign RS256 with `iat` = now+1h and no `exp`/`nbf`, assert `Verify` succeeds.
- [x] auth-empty-claims-403-test [auto]: pin in internal/auth/middleware_test.go that a validly signed token carrying NO Astarte claims is `403` on a guarded route, not `401` — the authenticated-but-grant-less distinction (middleware.go:149-151, `Authorizes` on a nil grants map → false → `WriteForbidden`) is currently untested: every 403 fixture has a claim on the wrong surface, so a regression mapping grant-less tokens to 401 (or skipping the forbidden path) passes the suite. Sign RS256 with empty `jwt.MapClaims`, hit the pairing route, expect 403 + golden `envelope_403.json`.
- [x] auth-cache-default-size-test [auto]: pin in internal/auth/cache_test.go that `NewCache(0)` and `NewCache(-1)` fall back to `DefaultCacheSize` (cache.go:29-31) — the `size < 1` branch exists purely so `lru.New` cannot panic and has zero coverage; a regression dropping the fallback (panic on `lru.New` with a non-tick size, or a silently tiny cache) would pass every test. Construct with both sizes, assert no panic and a working `Verify` with `lru.Len() == 0`.
- [x] hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto] — already resolved and re-queued by the refill: commit 015a706 chose the upstream-parity behaviour (0 ≡ unset, not invalid, not a literal 0) and pinned it — create folds 0 → nil with the container-free TestCreateRealmZeroRetentionUnsets (internal/housekeeping/service_test.go:186) and PATCH 0 → ClearRetention in the integration suite (internal/housekeeping/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).
---lines---
176 .mule/todo.md

✗ Edit .mule/todo.md failed
Error: No changes to apply: oldString and newString are identical.
The only drift on the astrate-native surface: the YAML's native-socket auth description claims `?token=` works, but the guard reads only the Authorization header. Let me append the proposal line to the queue.
Let me append the proposed task line at the end of the file.
→ Read .mule/todo.md [offset=173]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -172,5 +172,6 @@
 
 - [x] auth-iat-not-required-test [auto]: pin in internal/auth/jwt_test.go that a token with a *future* `iat` still verifies — the "`iat` is not required — upstream parity" rule (jwt.go:104-105, no `jwt.WithIssuedAt` in the parser, jwt.go:115-118) has zero coverage, so a regression that starts validating `iat` or `WithIssuedAt` would silently reject tokens from a skewed-clock issuer and pass every existing suite. Sign RS256 with `iat` = now+1h and no `exp`/`nbf`, assert `Verify` succeeds.
 - [x] auth-empty-claims-403-test [auto]: pin in internal/auth/middleware_test.go that a validly signed token carrying NO Astarte claims is `403` on a guarded route, not `401` — the authenticated-but-grant-less distinction (middleware.go:149-151, `Authorizes` on a nil grants map → false → `WriteForbidden`) is currently untested: every 403 fixture has a claim on the wrong surface, so a regression mapping grant-less tokens to 401 (or skipping the forbidden path) passes the suite. Sign RS256 with empty `jwt.MapClaims`, hit the pairing route, expect 403 + golden `envelope_403.json`.
 - [x] auth-cache-default-size-test [auto]: pin in internal/auth/cache_test.go that `NewCache(0)` and `NewCache(-1)` fall back to `DefaultCacheSize` (cache.go:29-31) — the `size < 1` branch exists purely so `lru.New` cannot panic and has zero coverage; a regression dropping the fallback (panic on `lru.New` with a non-tick size, or a silently tiny cache) would pass every test. Construct with both sizes, assert no panic and a working `Verify` with `lru.Len() == 0`.
-- [x] hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto] — already resolved and re-queued by the refill: commit 015a706 chose the upstream-parity behaviour (0 ≡ unset, not invalid, not a literal 0) and pinned it — create folds 0 → nil with the container-free TestCreateRealmZeroRetentionUnsets (internal/housekeeping/service_test.go:186) and PATCH 0 → ClearRetention in the integration suite (internal/housekeeping/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).
+- [x] hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null`
+- [ ] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication is "via Authorization header or `?token=` query parameter" (yaml:293-294), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto] — already resolved and re-queued by the refill: commit 015a706 chose the upstream-parity behaviour (0 ≡ unset, not invalid, not a literal 0) and pinned it — create folds 0 → nil with the container-free TestCreateRealmZeroRetentionUnsets (internal/housekeeping/service_test.go:186) and PATCH 0 → ClearRetention in the integration suite (internal/housekeeping/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).

Found one drift on the astrate-native surface. Append the proposal line to the queue.
→ Read .mule/todo.md [offset=176, limit=4]
My edit mangled line 176 — the oldString matched only part of it and the remainder got appended to my new line. Let me repair it properly.
→ Read .mule/todo.md [offset=174]
My edit spliced the old line's tail onto the new line. Let me read the full content to repair it.

$ sed -n '176,177p' .mule/todo.md | wc -c; sed -n '177p' .mule/todo.md | tail -c 300
2284
earRetention in the integration suite (internal/housekeeping/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).

Let me repair both lines: restore line 176 to its full original content and trim the tail off line 177.
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -172,6 +172,6 @@
 
 - [x] auth-iat-not-required-test [auto]: pin in internal/auth/jwt_test.go that a token with a *future* `iat` still verifies — the "`iat` is not required — upstream parity" rule (jwt.go:104-105, no `jwt.WithIssuedAt` in the parser, jwt.go:115-118) has zero coverage, so a regression that starts validating `iat` or `WithIssuedAt` would silently reject tokens from a skewed-clock issuer and pass every existing suite. Sign RS256 with `iat` = now+1h and no `exp`/`nbf`, assert `Verify` succeeds.
 - [x] auth-empty-claims-403-test [auto]: pin in internal/auth/middleware_test.go that a validly signed token carrying NO Astarte claims is `403` on a guarded route, not `401` — the authenticated-but-grant-less distinction (middleware.go:149-151, `Authorizes` on a nil grants map → false → `WriteForbidden`) is currently untested: every 403 fixture has a claim on the wrong surface, so a regression mapping grant-less tokens to 401 (or skipping the forbidden path) passes the suite. Sign RS256 with empty `jwt.MapClaims`, hit the pairing route, expect 403 + golden `envelope_403.json`.
 - [x] auth-cache-default-size-test [auto]: pin in internal/auth/cache_test.go that `NewCache(0)` and `NewCache(-1)` fall back to `DefaultCacheSize` (cache.go:29-31) — the `size < 1` branch exists purely so `lru.New` cannot panic and has zero coverage; a regression dropping the fallback (panic on `lru.New` with a non-tick size, or a silently tiny cache) would pass every test. Construct with both sizes, assert no panic and a working `Verify` with `lru.Len() == 0`.
-- [x] hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null`
-- [ ] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication is "via Authorization header or `?token=` query parameter" (yaml:293-294), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto] — already resolved and re-queued by the refill: commit 015a706 chose the upstream-parity behaviour (0 ≡ unset, not invalid, not a literal 0) and pinned it — create folds 0 → nil with the container-free TestCreateRealmZeroRetentionUnsets (internal/housekeeping/service_test.go:186) and PATCH 0 → ClearRetention in the integration suite (internal/housekeeping/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).
+- [x] hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto] — already resolved and re-queued by the refill: commit 015a706 chose the upstream-parity behaviour (0 ≡ unset, not invalid, not a literal 0) and pinned it — create folds 0 → nil with the container-free TestCreateRealmZeroRetentionUnsets (internal/housekeeping/service_test.go:186) and PATCH 0 → ClearRetention in the integration suite (internal/housekeeping/http_test.go:272) — enforced at every layer (http.go:217, service.go:151-159, store/realms.go:170-175). Both options the line offers contradict measured upstream v1.2.0 parity, so neither was applied (see the 20260915T170704Z done log).
+- [ ] docs-sync-native-socket-query-token-auth [auto]: the `/astrate/v1/{realm}/socket` description in docs/api/astrate_native_api.yaml says authentication is "via Authorization header or `?token=` query parameter" (yaml:293-294), but the route's guard `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ wc -l .mule/todo.md; tail -c 600 .mule/todo.md
177 .mule/todo.md
ws.go:46) extracts the credential only from the Authorization header (`bearerToken`, internal/auth/middleware.go:131,178-196) — no `?token=` fallback exists anywhere on this route, so the spec advertises an auth mode the code cannot serve. Align the spec with the code (drop the `?token=` claim; clients like browser EventSource that cannot set headers would need a token param, so if that pathway is intended, implement it instead — the Phoenix socket does read `?token=` at internal/appengine/channels/ws.go:53-54). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

Fixed. Now the todo file has the restored done-line plus my new proposal line. Let me verify the doc-sync proposal doesn't require a for-giulio entry (it's a docs/api fix task, so no — that's a todo line). The recipe says proposals max five; I found one drift worth proposing. Done.

Done: docs-sync scan of the astrate-native surface — found one drift and proposed one task line
Files: .mule/todo.md
Verified: rg/code comparison of docs/api/astrate_native_api.yaml vs internal route registrations -> one drift found; no test needed (proposal job)
Unsure: none
Follow-ups: `- [ ] docs-sync-native-socket-query-token-auth [auto]: ...` (one line appended to .mule/todo.md; not started, no git touched)
