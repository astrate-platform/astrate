slug: auth-empty-claims-403-test
verdict: done
at:  b91247c
ran: 2026-09-14T16:54:13Z on DietPi in 348s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
The task: write a test that a validly signed token with no Astarte claims gets 403 (not 401) on a guarded route. Let me understand the existing test patterns first.
→ Read internal/auth/middleware_test.go
→ Read internal/auth/middleware.go
Let me look at the test helpers (`keys`, `publicPEM`, `signToken`) and the golden files.
✱ Grep "func (keys|signToken|publicPEM)" in internal/auth · 3 matches
✱ Glob "internal/auth/testdata/*" 2 matches
→ Read internal/auth/jwt_test.go [offset=1, limit=80]
→ Read internal/auth/testdata/envelope_403.json
Let me check `Authorizes` to confirm the nil-grants-map path.
✱ Grep "func \(t \*Token\) Authorizes" in internal/auth · 2 matches
→ Read internal/auth/claims.go [offset=140, limit=80]
→ Read internal/auth/claims.go [limit=140]
Let me check the `Verify` cache path to make sure an empty-claims token actually verifies (no required claim).
✱ Grep "func \(c \*Cache\) Verify" in internal/auth · 1 match
→ Read internal/auth/cache.go [limit=60, offset=60]
✱ Grep "func \(c \*Cache\) verify" in internal/auth/cache.go · 0 matches
✱ Grep "c.verify|func .*\.verify\(" in internal/auth · 3 matches
✱ Grep "func verify|verify\(tokenString" in internal/auth · 1 match
✱ Grep "verify = |c.verify = func" in internal/auth · 2 matches
→ Read internal/auth/cache.go [limit=59]
✱ Grep "^func Verify" in internal/auth · 1 match
→ Read internal/auth/jwt.go [limit=90, offset=90]
Current behavior is confirmed 403 for a grant-less token (`Authorizes` hits a nil grant list → `authForbidden`). Now I'll add the pinning subtest.

← Edit internal/auth/middleware_test.go
Index: /root/astrate-mule/internal/auth/middleware_test.go
===================================================================
--- /root/astrate-mule/internal/auth/middleware_test.go
+++ /root/astrate-mule/internal/auth/middleware_test.go
@@ -129,8 +129,21 @@
 	}
 	testutil.Golden(t, "envelope_403.json", w.Body.Bytes())
 })
 
+t.Run("NoAstarteClaims403", func(t *testing.T) {
+	// A validly signed token carrying no Astarte claims authenticates
+	// (signature ok) but authorizes nothing: 403, not 401. Every other
+	// 403 fixture has some claim on the wrong surface; this pins the
+	// grant-less path (empty grants map -> Authorizes false).
+	grantlessToken := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{})
+	w := do(mux, "POST", "/pairing/v1/test/agent/devices", grantlessToken)
+	if w.Code != http.StatusForbidden {
+		t.Fatalf("status: got %d, want 403", w.Code)
+	}
+	testutil.Golden(t, "envelope_403.json", w.Body.Bytes())
+})
+
 t.Run("Authorized200", func(t *testing.T) {
 	w := do(mux, "POST", "/pairing/v1/test/agent/devices", pairingToken)
 	if w.Code != http.StatusOK {
 		t.Fatalf("status: got %d, want 200 (body %s)", w.Code, w.Body.String())

Now the gate: race-enabled tests plus gofmt.
$ gofmt -l internal/auth/
(no output)

$ go test -race ./internal/auth/ 2>&1 | tail -20
(no output)

<shell_metadata>
shell tool terminated command after exceeding timeout 120000 ms. If this command is expected to take longer and is not waiting for interactive input, retry with a larger timeout value in milliseconds.
</shell_metadata>

Race build timed out (cold compile + RSA keygen). Retrying with a larger timeout.
$ go test -race ./internal/auth/ 2>&1 | tail -20
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/auth	0.064s
FAIL

Race detector unsupported here (Pi, 39-bit VMA) — that's the expected constraint from MULE.md. Running the container-free gate instead.
$ go vet ./... && go test ./internal/auth/ 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/auth	2.551s

$ go test ./... 2>&1 | tail -30
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

Done: pinned in internal/auth/middleware_test.go that an RS256 token with empty MapClaims is 403 (Forbidden envelope) on the pairing route, not 401.
Files: internal/auth/middleware_test.go
Verified: `go vet ./...` -> pass; `go test ./...` -> pass; `gofmt -l internal/auth/` -> clean; `go test -race` -> unavailable (39-bit VMA, as MULE.md warns) — change is test-only, no concurrency touched
Unsure: nothing
Follow-ups: none
