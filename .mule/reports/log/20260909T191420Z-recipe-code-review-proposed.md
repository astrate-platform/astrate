slug: recipe-code-review
verdict: proposed
at:  a701b05
ran: 2026-09-09T19:14:20Z on DietPi in 364s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
→ Read internal
→ Read pkg
→ Read .mule/reviews
Let me check what packages exist and what was recently reviewed.
→ Read .mule/log.md
Pairing was just reviewed today. Let me check sizes of unreviewed packages.

$ for dir in auth config flowapi housekeeping httpx observability realm swagger testutil; do echo -n "internal/$dir: "; rg -c '' --glob "internal/$dir/**/*.go" 2>/dev/null | awk -F: '{s+=$2} END{print s+0}'; done && echo "---" && for dir in astarteapi deviceid interfaceschema payload; do echo -n "pkg/$dir: "; rg -c '' --glob "pkg/$dir/**/*.go" 2>/dev/null | awk -F: '{s+=$2} END{print s+0}'; done
internal/auth: 1532
internal/config: 818
internal/flowapi: 2230
internal/housekeeping: 1311
internal/httpx: 357
internal/observability: 398
internal/realm: 2573
internal/swagger: 171
internal/testutil: 697
---
pkg/astarteapi: 485
pkg/deviceid: 358
pkg/interfaceschema: 3176
pkg/payload: 2673

I'll review `internal/realm` — 2573 lines, never reviewed. Let me explore it.
✱ Glob "internal/realm/**/*.go" 10 matches

$ wc -l /root/astrate-mule/internal/realm/*.go | sort -rn
2573 total
   782 /root/astrate-mule/internal/realm/http_test.go
   686 /root/astrate-mule/internal/realm/service.go
   569 /root/astrate-mule/internal/realm/http.go
   213 /root/astrate-mule/internal/realm/dashboard_compat_test.go
    76 /root/astrate-mule/internal/realm/detailed.go
    70 /root/astrate-mule/internal/realm/policy_test.go
    68 /root/astrate-mule/internal/realm/async_operation_test.go
    63 /root/astrate-mule/internal/realm/detailed_test.go
    29 /root/astrate-mule/internal/realm/detailed_required_encrypted_test.go
    17 /root/astrate-mule/internal/realm/policy.go

→ Read internal/realm/service.go
→ Read internal/realm/http.go
→ Read internal/realm/policy.go
→ Read internal/realm/detailed.go
→ Read internal/realm/http_test.go
→ Read internal/realm/dashboard_compat_test.go
→ Read internal/realm/async_operation_test.go
Now let me read the tests.
→ Read internal/realm/detailed_test.go
→ Read internal/realm/detailed_required_encrypted_test.go
→ Read internal/realm/policy_test.go
✱ Grep "realm\.(Service|API|NewService|NewAPI)" in internal · 0 matches
Now let me check what store functions realm calls, and what callers there are of the realm service.
✱ Grep "func \(s \*Service\)" in internal/realm · 26 matches
✱ Grep ""internal/realm"" in . · 0 matches
✱ Grep "astrate/internal/realm" in . · 3 matches

$ grep -n "realm" /root/astrate-mule/test/conformance/cpc/cpc_test.go | head -30
8://	housekeeping realms create  → broker realm reload →
9://	realm-management interfaces install (genericsensors.Values) →
12://	realm-management triggers install/list/show/delete round-trip.
45:	"github.com/astrate-platform/astrate/internal/realm"
99:	realmName string
102:	realmKey  string // path to the realm JWT private key (PEM)
103:	realmPub  string // path to the realm JWT public key (PEM)
158:// cmd/astrate will (M8). The realm is *not* pre-created: astartectl creates it.
179:	realmKey, realmPub, _ := rsaKeyFiles(t, keyDir, "realm")
226:	realm.NewAPI(realm.NewService(st, e, nil), mw).Mount(mux)
235:	realmName := "cpc" + hexLower(suffix[:])
238:		st: st, e: e, broker: b, pairer: pairer, server: server, realmName: realmName,
240:		realmKey: realmKey, realmPub: realmPub, instKey: instKey,
244:// hexLower renders bytes as lowercase hex (a realm name must be lowercase and
262:	// 1. Housekeeping: create the realm with astartectl (instance-admin JWT).
264:		"housekeeping", "realms", "create", f.realmName,
266:		"--realm-public-key", f.realmPub,
269:	realmRow, err := f.st.GetRealmByName(ctx, f.realmName)
271:		t.Fatalf("realm not created by astartectl: %v", err)
273:	// The broker only learns the new realm's CA on reload (M8 wires this to the
279:	// 2. Realm Management: install the interface (realm JWT).
285:		"realm-management", "interfaces", "install", ifaceFile,
286:		"--realm-key", f.realmKey, "-r", f.realmName, "-u", url,
288:	si, err := f.st.GetInterface(ctx, realmRow.ID, valuesName, 1)
293:	// 3. Pairing: register the device with astartectl (realm JWT → a_pa).
301:		"--realm-name", f.realmName,
302:		"--realm-key", f.realmKey,
311:	crt, err := f.pairer.Credentials(ctx, f.realmName, id.String(), secret, csr, netip.MustParseAddr("127.0.0.1"))
316:	dev := testutil.ConnectAstarteDevice(t, f.sslURL, f.realmName, id, tlsCfg, true)
320:		d, err := f.st.GetDevice(ctx, realmRow.ID, id)

Let me look at one more thing — the conformance test to understand if there are gaps.
✱ Grep "TODO|FIXME|HACK|XXX|BUG" in internal/realm · 0 matches
Now let me check the existing todo.md so I know where to append.
→ Read .mule/todo.md
Now I have a thorough understanding of the package. Let me write the review file and propose tasks.

← Write .mule/reviews/realm-2026-09-09.md
Wrote file successfully.

Now let me append the task lines to todo.md.
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -131,4 +131,6 @@
 - [x] pairing-bearer-secret-test [auto]: add a container-free unit test for `bearerSecret` (internal/pairing/http.go:343) locking the upstream-verbatim permissive match `~r/bearer\:?\s+(.*)$/i` — the colon form `Bearer: <secret>`, the case variants (`bearer `, `BEARER `), and the reject paths (no whitespace before value, non-bearer scheme, empty value). Only the plain `Bearer <secret>` form is exercised today (integration http_test.go); a regression to case-sensitive or drop-colon matching would pass every existing test. Pure helper in a non-integration file; no Docker.
 - [x] pairing-remoteip-fallback-test [auto]: add a container-free unit test for `remoteIP` (internal/pairing/http.go:359) — a valid addr:port yields `ap.Addr().Unmap()`, a bare address yields the mapped address, and an unparseable `RemoteAddr` (unix-socket/proxy case) falls back to `netip.IPv4Unspecified()`. That fallback branch is currently uncovered (only reachable through the integration HTTP stack). Pure helper in a non-integration file; no Docker.
 - [x] docs-sync-hk-patch-endpoint [auto]: add the undocumented `PATCH /housekeeping/v1/realms/{realm}` route to docs/api/astarte_housekeeping_api.yaml — it exists in code (internal/housekeeping/http.go:40, handler patchRealm at http.go:163-229, returns 200 on success) but the spec jumps straight from GET to DELETE. Document the request body (`RealmPatch` — same three optional fields as `RealmCreate`: `jwt_public_key_pem`, `device_registration_limit`, `datastream_maximum_storage_retention`, all nullable), the success response, and the error responses: 400 (bad envelope), 401, 403, 404, 405 (deletion disabled), 422 (invalid update parameters for unknown fields, blank jwt_public_key_pem, negative limit/retention, connected devices present), 500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
 - [x] docs-sync-hk-retention-field [auto]: add the `datastream_maximum_storage_retention` field to the `Realm` and `RealmCreate` schemas in docs/api/astarte_housekeeping_api.yaml — the field is on the wire (realmBody struct, internal/housekeeping/http.go:50, returned by viewBody at http.go:99, accepted in create at http.go:68 and patch at http.go:216) but absent from both schemas (yaml:200-231). Type: integer, format int64, nullable, description matching the device_registration_limit style. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
+
+- [ ] realm-pure-helper-tests [auto]: add container-free unit tests for the four untested pure helpers in internal/realm. Create `internal/realm/helpers_test.go` (no build tag) covering: (1) `normaliseIfaceName` — lowercase, hyphen strip, mixed case+hyphens, empty string, already-normalized passthrough; (2) `joinPEM` — empty slice, single key, two keys separated by `\n`; (3) `triggerErrorBody` — empty action+simple_triggers, action-only, simple_triggers-only with index-aligned `{}` substitution, both present; (4) `majorParam` — valid positive, zero, negative (→404), non-numeric (→404). All tests use httptest.ResponseRecorder for majorParam; the rest are direct function calls. No Docker.

Let me verify the review and the appended task line.
→ Read .mule/todo.md [offset=130]
→ Read .mule/reviews/realm-2026-09-09.md
Done: reviewed `internal/realm` (1343 LOC source + 1225 LOC tests), proposed one task line for container-free unit tests covering four untested pure helpers
Files: `.mule/reviews/realm-2026-09-09.md`, `.mule/todo.md`
Verified: read all 10 `.go` files in the package — no behavioral bugs found, code is defensively written
Unsure: whether the store's `DeleteTriggerPolicy` has a FK constraint that makes the service-level referential check in `DeletePolicy` redundant (it provides a better error message either way, so no action needed)
Follow-ups: `realm-pure-helper-tests` (1 task line appended)
