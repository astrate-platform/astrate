slug: pairing-burn-bcrypt-cost
verdict: blocked
at:  4ab5f93
ran: 2026-10-07T18:02:02Z on DietPi in 359s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/pairing/service.go [offset=60, limit=160]
→ Read internal/pairing/service.go [offset=380]
✱ Grep "dummySecretHash|burnBcrypt" in . · 20 matches
✱ Glob "internal/pairing/*_test.go" 5 matches
✱ Grep "^func Test|New\(" in internal/pairing/service_test.go · 49 matches
→ Read internal/pairing/service_test.go [limit=240]
✱ Grep "dummy|equaliz|uniform" in docs · 15 matches
✱ Grep "BcryptCost" in . · 9 matches
Task means: make the pairing failure-path bcrypt burn use `s.cfg.BcryptCost` instead of a fixed cost-10 package hash, plus a test that fails today.

← Edit internal/pairing/service.go
Index: /root/astrate-mule/internal/pairing/service.go
===================================================================
--- /root/astrate-mule/internal/pairing/service.go
+++ /root/astrate-mule/internal/pairing/service.go
@@ -87,18 +87,12 @@
 	// is empty.
 	DefaultVersion = "0.1.0-astrate"
 )
 
-// dummySecretHash is a bcrypt hash of an unguessable value, compared against
-// when the device row (or its hash) is missing so authentication failures
-// burn comparable time regardless of cause.
-var dummySecretHash = func() string {
-	h, err := bcrypt.GenerateFromPassword([]byte("astrate-dummy-secret-equalizer"), bcrypt.DefaultCost)
-	if err != nil {
-		panic(fmt.Sprintf("pairing: generating dummy bcrypt hash: %v", err))
-	}
-	return string(h)
-}()
+// dummySecretPlaintext is the value the equalizer hash is generated from. It
+// never leaves the process: it only exists so a failure path can spend one
+// bcrypt comparison of the configured cost.
+const dummySecretPlaintext = "astrate-dummy-secret-equalizer"
 
 // Config carries the service's operational knobs.
 type Config struct {
 	// BrokerURL is handed to devices by the info endpoint

← Edit internal/pairing/service.go
Index: /root/astrate-mule/internal/pairing/service.go
===================================================================
--- /root/astrate-mule/internal/pairing/service.go
+++ /root/astrate-mule/internal/pairing/service.go
@@ -123,8 +123,13 @@
 	st     Store
 	sealer *store.KeySealer
 	cfg    Config
 	now    func() time.Time
+	// dummyOnce guards lazy generation of the equalizer hash at
+	// cfg.BcryptCost, so failure paths burn the same work as a real
+	// comparison (docs/DESIGN.md §4.4).
+	dummyOnce sync.Once
+	dummyHash string
 	// OnRegistered, when non-nil, is called after a successful registration.
 	OnRegistered OnRegisteredFunc
 }
 

← Edit internal/pairing/service.go
Index: /root/astrate-mule/internal/pairing/service.go
===================================================================
--- /root/astrate-mule/internal/pairing/service.go
+++ /root/astrate-mule/internal/pairing/service.go
@@ -421,12 +421,26 @@
 	}
 	return realmCA, nil
 }
 
+// burnHash returns the Service's equalizer hash, generated on first use at
+// the configured bcrypt cost so exactly one hash per cost ever exists.
+func (s *Service) burnHash() string {
+	s.dummyOnce.Do(func() {
+		h, err := bcrypt.GenerateFromPassword([]byte(dummySecretPlaintext), s.cfg.BcryptCost)
+		if err != nil {
+			panic(fmt.Sprintf("pairing: generating dummy bcrypt hash: %v", err))
+		}
+		s.dummyHash = string(h)
+	})
+	return s.dummyHash
+}
+
 // burnBcrypt spends one bcrypt comparison against the dummy hash, equalizing
-// the work done on authentication paths that lack a real hash.
-func burnBcrypt(secret string) {
-	_ = bcrypt.CompareHashAndPassword([]byte(dummySecretHash), []byte(secret))
+// the work done on authentication paths that lack a real hash. The comparison
+// runs at cfg.BcryptCost, matching the cost of a real comparison.
+func (s *Service) burnBcrypt(secret string) {
+	_ = bcrypt.CompareHashAndPassword([]byte(s.burnHash()), []byte(secret))
 }
 
 // generateSecret draws the 32-byte credentials secret and encodes it as the
 // upstream-parity 44-character standard base64 string.

← Edit internal/pairing/service.go
Index: /root/astrate-mule/internal/pairing/service.go
===================================================================
--- /root/astrate-mule/internal/pairing/service.go
+++ /root/astrate-mule/internal/pairing/service.go
@@ -378,30 +378,30 @@
 // amount of bcrypt work, so response timing does not leak device existence.
 func (s *Service) authenticateDevice(ctx context.Context, realmName, deviceIDStr, secret string) (*store.Realm, *store.Device, error) {
 	id, err := deviceid.Parse(deviceIDStr)
 	if err != nil {
-		burnBcrypt(secret)
+		s.burnBcrypt(secret)
 		return nil, nil, ErrUnauthorized
 	}
 	realm, err := s.st.GetRealmByName(ctx, realmName)
 	if errors.Is(err, store.ErrNotFound) {
-		burnBcrypt(secret)
+		s.burnBcrypt(secret)
 		return nil, nil, ErrUnauthorized
 	}
 	if err != nil {
 		return nil, nil, err
 	}
 	dev, err := s.st.GetDevice(ctx, realm.ID, id)
 	if errors.Is(err, store.ErrNotFound) {
-		burnBcrypt(secret)
+		s.burnBcrypt(secret)
 		return nil, nil, ErrUnauthorized
 	}
 	if err != nil {
 		return nil, nil, err
 	}
 	if dev.CredentialsSecretHash == "" {
 		// Unregistered (or never-registered) device: no valid hash exists.
-		burnBcrypt(secret)
+		s.burnBcrypt(secret)
 		return nil, nil, ErrUnauthorized
 	}
 	if bcrypt.CompareHashAndPassword([]byte(dev.CredentialsSecretHash), []byte(secret)) != nil {
 		return nil, nil, ErrUnauthorized

→ Read internal/pairing/service.go [limit=40]

← Edit internal/pairing/service.go
Index: /root/astrate-mule/internal/pairing/service.go
===================================================================
--- /root/astrate-mule/internal/pairing/service.go
+++ /root/astrate-mule/internal/pairing/service.go
@@ -12,8 +12,9 @@
 "encoding/base64"
 "errors"
 "fmt"
 "net/netip"
+"sync"
 "time"
 
 "golang.org/x/crypto/bcrypt"
 

→ Read internal/pairing/service_test.go [limit=70, offset=372]
→ Read internal/pairing/service_test.go [limit=35, offset=442]
✗ Edit internal/pairing/service_test.go failed
Error: No changes to apply: oldString and newString are identical.

← Edit internal/pairing/service_test.go
Index: /root/astrate-mule/internal/pairing/service_test.go
===================================================================
--- /root/astrate-mule/internal/pairing/service_test.go
+++ /root/astrate-mule/internal/pairing/service_test.go
@@ -452,8 +452,30 @@
 		}
 	}
 }
 
+// TestBurnHashUsesConfiguredCost pins the equalizer hash to the configured
+// bcrypt cost: with bcrypt_cost below the default, the failure paths that
+// lack a real hash must burn cost-N work like a real comparison, or they
+// answer slower than a wrong-secret attempt and become a latency oracle for
+// device existence (docs/DESIGN.md §4.4 uniform auth time).
+func TestBurnHashUsesConfiguredCost(t *testing.T) {
+	for _, cost := range []int{4, bcrypt.DefaultCost} {
+		svc, _, _ := newServiceFixture(t, Config{BcryptCost: cost})
+		first := svc.burnHash()
+		got, err := bcrypt.Cost([]byte(first))
+		if err != nil {
+			t.Fatalf("BcryptCost %d: burn hash does not parse: %v", cost, err)
+		}
+		if got != cost {
+			t.Errorf("BcryptCost %d: burn hash cost: got %d, want %d", cost, got, cost)
+		}
+		if again := svc.burnHash(); again != first {
+			t.Errorf("BcryptCost %d: burn hash regenerated: got %q, want %q", cost, again, first)
+		}
+	}
+}
+
 func TestCredentialsInhibited(t *testing.T) {
 	ctx := context.Background()
 	svc, fs, _ := newServiceFixture(t, Config{})
 	hwID := randomDeviceID(t)

$ gofmt -l . && go vet ./... 2>&1 | tail -20
(no output)

$ go test -race ./internal/pairing/ 2>&1 | tail -20
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/pairing	0.022s
FAIL

$ go test ./... 2>&1 | tail -40
ok  	github.com/astrate-platform/astrate/cmd/astrate	0.192s
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
ok  	github.com/astrate-platform/astrate/internal/pairing	6.414s
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

$ du -sh /root/astrate-mule 2>/dev/null; du -sh /root/astrate-mule/.git 2>/dev/null
96M	/root/astrate-mule
20M	/root/astrate-mule/.git

! permission requested: external_directory (/tmp/*); auto-rejecting
✗ rm -rf /tmp/mule-burncheck && mkdir -p /tmp/mule-burncheck && cp -a /root/astrate-mule/. /tmp/mule-burncheck/ && ls /tmp/mule-burncheck | head failed
Error: The user rejected permission to use this specific tool call.
