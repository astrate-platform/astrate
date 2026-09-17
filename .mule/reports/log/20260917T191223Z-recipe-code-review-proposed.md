slug: recipe-code-review
verdict: proposed
at:  66f4932
ran: 2026-09-17T19:12:23Z on DietPi in 246s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
→ Read .mule/reviews
→ Read internal
→ Read pkg
→ Read .mule/log.md

$ wc -l internal/httpx/*.go internal/swagger/*.go internal/testutil/*.go pkg/deviceid/*.go 2>/dev/null | sort -rn | head -40
1583 total
  226 pkg/deviceid/deviceid_test.go
  216 internal/testutil/astartedevice.go
  163 internal/testutil/mqttclient.go
  139 internal/httpx/notfound_test.go
  132 pkg/deviceid/deviceid.go
  128 internal/swagger/swagger_test.go
  117 internal/testutil/golden_test.go
  113 internal/httpx/cors_test.go
  106 internal/testutil/pg.go
   59 internal/testutil/golden.go
   56 internal/httpx/notfound.go
   49 internal/httpx/cors.go
   43 internal/swagger/swagger.go
   36 internal/testutil/pg_test.go

Let me look at an existing review file to match the format, then examine the unreviewed packages.
→ Read .mule/reviews/housekeeping-2026-09-15.md
→ Read .mule/reviews/observability-2026-09-16.md
All 15 `internal/**` main areas have reviews; `httpx`, `swagger`, `testutil`, and `deviceid` do not. I'll go deep on `internal/httpx` (real production middleware, small, unreviewed).
→ Read internal/httpx/cors.go
→ Read internal/httpx/notfound.go
→ Read internal/httpx/notfound_test.go
→ Read internal/httpx/cors_test.go
Package is small and well-tested. Let me check how these are wired in `main` and whether any config drives CORS origins, then verify gate behaviour.
✱ Grep "Vary|cors|CORS|AllowedOrigins|allowed_origins" 43 matches
✱ Grep "httpx\.(CORS|NotFound)" 5 matches
→ Read cmd/astrate/main.go [offset=340, limit=90]
✱ Grep "HandleFunc\(|\.Handle|/v1|Mount\(" 0 matches
Let me confirm which prefixes the other APIs mount under, to write accurate "not proposing" notes.
✱ Grep "api-register|api-new|/flow|Mount\(mux\)|Handle" 37 matches
→ Read .mule/for-giulio.md
Now let me check the current `todo.md` to append the task line in style.
→ Read .mule/todo.md
→ Read .mule/todo.md [offset=173]
Now I'll write the review file and append the task line.

← Write .mule/reviews/httpx-2026-09-17.md
Wrote file successfully.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -183,5 +183,6 @@
 - [x] observability-dbpool-gauge-coverage [auto]: extend `TestMetricsExposesGauges` (internal/observability/observability_test.go:25-34) to assert all four `astrate_db_pool_*` gauges — it currently pins only `acquired_conns` and `max_conns`, so a regression deleting or mis-wiring `idle_conns` or `total_conns` (metrics.go:77-78) passes the suite; expect idle 5, total 7 from the supplied `DBPoolStats{2,5,7,10}`.
 - [x] observability-health-content-type-test [auto]: assert `Content-Type: application/json` on `/astrate/v1/health` and `/astrate/v1/readiness` (set by `writeJSON`, health.go:81) and on both the 200 and 503 paths of the per-service `/health` (set by hand in compat.go:26,31) — only the version endpoint pins the header today (compat_version_test.go:31-33), so a regression dropping the header on the /astrate/v1 surface passes the suite.
 
 - [x] docs-sync-pairing-version-value [auto]: fix `GET /pairing/v1/{realm}/version` in docs/api/astarte_pairing_api.yaml — the description claims it "Returns the emulated upstream API compatibility version" that "the Astarte Dashboard gates feature UI on" (yaml:74-78) and the example is `data: "1.1.0"` (yaml:95), but the route is served by `observability.VersionHandler(version)` (cmd/astrate/main.go:398) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — a value no build serves, and not an API-compat level (the only such constant is RM's, internal/realm/service.go:588 = "1.2.2"). Either make the spec truthful (description + example "0.1.0-dev") or wire the handler to the emulated level like RM's realm-scoped op — say which; the Dashboard-gating sentence makes the code-side fix the likely right choice, but that is Giulio's call. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+- [ ] httpx-cors-vary-origin-passthrough [auto]: in internal/httpx/cors.go non-wildcard mode, the disallowed-origin pass-through (cors.go:26-29) returns without `Vary: Origin`, so responses whose CORS outcome depends on the Origin header reach shared caches without the RFC 9110 §12.5.5 variation marker (a cache could serve the CORS-approving variant to a disallowed origin, or the headerless variant to an allowed one); the stamp lives only in the allowed branch (cors.go:35-36) and only the preflight pins it (cors_test.go:42). Add `h.Add("Vary", "Origin")` on the pass-through branch too, and assert `Vary: Origin` in TestCORSActualRequestStamped (actual non-preflight, pinned only for Allow-Origin today) and TestCORSDisallowedOriginPassesThroughUnstamped. Additive header change; container-free.
 - [x] docs-sync-native-version-value [auto]: fix the four compat-version ops in docs/api/astrate_native_api.yaml — `/appengine/version` (yaml:185-202), `/realmmanagement/version` (yaml:204-221), `/pairing/version` (yaml:223-240), `/housekeeping/version` (yaml:242-259) — plus `DataEnvelopeVersion.description` (yaml:416): all four descriptions say "Returns the API compatibility version of the <Service> service" and every example is `data: "1.2.2"`, but `observability.MountVersionCompat` (cmd/astrate/main.go:393-395) serves the Astrate build `version` var (main.go:53, default "0.1.0-dev"), not the emulated API level; the earlier docs-sync-native-compat-version-endpoints task codified the 1.2.2 premise without checking the wired value, so the spec now advertises a value the wire never serves. Document the build version (description + example "0.1.0-dev") or, to match the RM realm-scoped endpoint (internal/realm/service.go:588) and the compat.md:75 claim, serve `APICompatVersion` instead — decide. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

→ Read .mule/todo.md [offset=186]
Done: reviewed internal/httpx (unreviewed since repo start), wrote the review, appended one task line.
Files: .mule/reviews/httpx-2026-09-17.md, .mule/todo.md
Verified: read cors.go, notfound.go, both test files, main.go wiring, config validation (review only, no code ran)
Unsure: nothing
Follow-ups: `httpx-cors-vary-origin-passthrough` — non-wildcard CORS omits `Vary: Origin` on the disallowed-origin pass-through (cors.go:26-29), so origin-dependent responses reach shared caches unmarked; also `Vary` is unpinned on actual non-preflight responses. One additive header change + tests, container-free.
