slug: auth-iat-not-required-test
verdict: done
at:  6524c7f
ran: 2026-09-14T11:44:42Z on DietPi in 257s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/auth/jwt.go
→ Read internal/auth/jwt_test.go
✱ Grep "func newToken|type Token struct|type astarteClaims" in internal/auth · 3 matches
Now I'll add the test. It pins that a future-`iat` token without `exp`/`nbf` verifies successfully — protecting the upstream-parity rule against a regression that starts enforcing `iat`.

← Edit internal/auth/jwt_test.go
Index: /root/astrate-mule/internal/auth/jwt_test.go
===================================================================
--- /root/astrate-mule/internal/auth/jwt_test.go
+++ /root/astrate-mule/internal/auth/jwt_test.go
@@ -199,8 +199,28 @@
 		t.Errorf("Verify without signing key: got %v, want ErrNoKeyMatched", err)
 	}
 }
 
+func TestVerifyFutureIatAccepted(t *testing.T) {
+	tk := keys(t)
+	// iat is not validated (upstream parity, jwt.go:104-105). A token whose
+	// iat is in the future, with no exp/nbf, must still verify — a regression
+	// that adds jwt.WithIssuedAt to the parser would silently reject tokens
+	// from skewed-clock issuers.
+	futureIat := time.Now().Add(time.Hour).Unix()
+	s := signToken(t, jwt.SigningMethodRS256, tk.rsaKey, jwt.MapClaims{
+		"a_pa": []string{".*::.*"},
+		"iat":  futureIat,
+	})
+	tok, err := Verify(s, []crypto.PublicKey{&tk.rsaKey.PublicKey}, time.Now)
+	if err != nil {
+		t.Fatalf("Verify with future iat: %v", err)
+	}
+	if ts, ok := tok.IssuedAt(); !ok || ts.Unix() != futureIat {
+		t.Errorf("IssuedAt: got (%v, %v), want future iat %d", ts, ok, futureIat)
+	}
+}
+
 func TestParsePublicKeysPEM(t *testing.T) {
 	tk := keys(t)
 
 	t.Run("PKIXAndPKCS1", func(t *testing.T) {

Let me verify the `IssuedAt` method exists on `Token`:
→ Read internal/auth/claims.go [offset=121]
`Token` has no `IssuedAt` method. Let me fix the assertion to match the available API:

← Edit internal/auth/jwt_test.go
Index: /root/astrate-mule/internal/auth/jwt_test.go
===================================================================
--- /root/astrate-mule/internal/auth/jwt_test.go
+++ /root/astrate-mule/internal/auth/jwt_test.go
@@ -214,10 +214,10 @@
 	tok, err := Verify(s, []crypto.PublicKey{&tk.rsaKey.PublicKey}, time.Now)
 	if err != nil {
 		t.Fatalf("Verify with future iat: %v", err)
 	}
-	if ts, ok := tok.IssuedAt(); !ok || ts.Unix() != futureIat {
-		t.Errorf("IssuedAt: got (%v, %v), want future iat %d", ts, ok, futureIat)
+	if _, ok := tok.ExpiresAt(); ok {
+		t.Error("ExpiresAt: want none for exp-less token")
 	}
 }
 
 func TestParsePublicKeysPEM(t *testing.T) {

Now let me verify the new test passes and the gate is clean:
$ gofmt -l .
(no output)

$ go test ./internal/auth/... 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/auth	3.591s

Test passes, gofmt clean. Let me verify the test actually pins the behavior — that a `jwt.WithIssuedAt` regression would make it fail — with a probe:
$ mkdir -p /tmp/opencode/iatprobe && cat > /tmp/opencode/iatprobe/main.go <<'EOF'
package main

import (
	"crypto/rsa"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	claims := jwt.MapClaims{"iat": time.Now().Add(time.Hour).Unix()}
	s, _ := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	// parser WITHOUT WithIssuedAt (current code path)
	_, err := jwt.NewParser().ParseWithClaims(s, &jwt.MapClaims{}, func(*jwt.Token) (any, error) { return &key.PublicKey, nil })
	fmt.Printf("without WithIssuedAt: err=%v\n", err)
	// parser WITH WithIssuedAt (hypothetical regression)
	_, err = jwt.NewParser(jwt.WithIssuedAt()).ParseWithClaims(s, &jwt.MapClaims{}, func(*jwt.Token) (any, error) { return &key.PublicKey, nil })
	fmt.Printf("with WithIssuedAt:    err=%v\n", err)
}
EOF
cd /tmp/opencode/iatprobe && go mod init probe >/dev/null 2>&1 && go get github.com/golang-jwt/jwt/v5 >/dev/null 2>&1 && go run .
without WithIssuedAt: err=<nil>
with WithIssuedAt:    err=token has invalid claims: token used before issued

Confirmed: current parser accepts future `iat`; a `jwt.WithIssuedAt()` regression returns "token used before issued", which would fail the new test. Running the full gate:
$ go vet ./... && go test ./... 2>&1 | tail -20
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

All green. Final report:

```
Done: Added TestVerifyFutureIatAccepted — pins that a token with future iat and no exp/nbf verifies, protecting the upstream-parity rule (jwt.go:104-105) against a WithIssuedAt regression.
Files: internal/auth/jwt_test.go
Verified: go vet ./... -> pass; go test ./... -> pass; gofmt -l . -> clean
Unsure: nothing
Follow-ups: none
```
