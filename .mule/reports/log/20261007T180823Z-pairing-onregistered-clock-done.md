slug: pairing-onregistered-clock
verdict: done
at:  4e71533
ran: 2026-10-07T18:08:23Z on DietPi in 143s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

$ rg -n "OnRegistered|time\.Now\(\)|s\.now\(\)" internal/pairing/ | head -50
internal/pairing/service_test.go:101:	now := time.Now()
internal/pairing/service_test.go:196:	svc.OnRegistered = func(realmName, deviceID string, _ time.Time) {
internal/pairing/service_test.go:209:		t.Fatalf("OnRegistered calls: got %d, want 1", len(calls))
internal/pairing/service_test.go:227:		t.Errorf("OnRegistered after rejected re-register: got %d calls, want 1", len(calls))
internal/pairing/service_test.go:567:		svcExpired.now = func() time.Time { return time.Now().Add(31 * 24 * time.Hour) }
internal/pairing/ratelimit_test.go:11:	clock := time.Now()
internal/pairing/ratelimit_test.go:35:	clock := time.Now()
internal/pairing/ratelimit_test.go:54:	clock := time.Now()
internal/pairing/service.go:122:// OnRegisteredFunc is called after a successful device registration so the
internal/pairing/service.go:125:type OnRegisteredFunc func(realmName string, deviceID string, at time.Time)
internal/pairing/service.go:133:	// OnRegistered, when non-nil, is called after a successful registration.
internal/pairing/service.go:134:	OnRegistered OnRegisteredFunc
internal/pairing/service.go:221:	if s.OnRegistered != nil {
internal/pairing/service.go:222:		s.OnRegistered(realmName, hwID, time.Now())
internal/pairing/service.go:333:	res := &VerifyResult{Timestamp: s.now()}
internal/pairing/http_test.go:78:	realmName := "p" + strconv.FormatInt(time.Now().UnixNano(), 36)
internal/pairing/http_test.go:121:		"exp":         time.Now().Add(time.Hour).Unix(),
internal/pairing/http_test.go:330:		startWrong := time.Now()
internal/pairing/http_test.go:337:		startUnknown := time.Now()
internal/pairing/http_test.go:597:			"/pairing/v1/nope"+strconv.FormatInt(time.Now().UnixNano(), 36)+"/health", "", nil)
internal/pairing/ca/ca.go:84:	now := time.Now()
internal/pairing/ca/ca.go:175:	now := time.Now()
internal/pairing/ca/ca_test.go:75:	wantAfter := time.Now().Add(DefaultCALifetime - 24*time.Hour)
internal/pairing/ca/ca_test.go:96:	if _, err := loaded.Verify(certPEM, time.Now()); err != nil {
internal/pairing/ca/ca_test.go:255:		until, err := c.Verify(certPEM, time.Now())
internal/pairing/ca/ca_test.go:267:		_, err := c.Verify(certPEM, time.Now().Add(2*time.Hour))
internal/pairing/ca/ca_test.go:279:		if _, err := c.Verify(foreignCert, time.Now()); !errors.Is(err, ErrCertificateInvalid) {
internal/pairing/ca/ca_test.go:285:		if _, err := c.Verify("not a certificate", time.Now()); !errors.Is(err, ErrCertificateInvalid) {

→ Read internal/pairing/service.go [offset=190, limit=80]

$ rg -n "now\b" internal/pairing/service.go | head -30; echo ---; rg -n "TestRegisterEmitsEvent" -A 45 internal/pairing/service_test.go
132:	now    func() time.Time
150:	return &Service{st: st, sealer: sealer, cfg: cfg, now: time.Now}
333:	res := &VerifyResult{Timestamp: s.now()}
---
186:func TestRegisterEmitsEvent(t *testing.T) {
187-	ctx := context.Background()
188-	svc, fs, _ := newServiceFixture(t, Config{})
189-	hwID := randomDeviceID(t)
190-
191-	type registeredCall struct {
192-		realm string
193-		id    string
194-	}
195-	var calls []registeredCall
196-	svc.OnRegistered = func(realmName, deviceID string, _ time.Time) {
197-		calls = append(calls, registeredCall{realm: realmName, id: deviceID})
198-	}
199-
200-	secret, err := svc.Register(ctx, "test", hwID, "")
201-	if err != nil {
202-		t.Fatalf("Register: %v", err)
203-	}
204-	if len(secret) != 44 {
205-		t.Errorf("secret length: got %d, want 44", len(secret))
206-	}
207-
208-	if len(calls) != 1 {
209-		t.Fatalf("OnRegistered calls: got %d, want 1", len(calls))
210-	}
211-	if calls[0].realm != "test" {
212-		t.Errorf("realm: got %q, want test", calls[0].realm)
213-	}
214-	if calls[0].id != hwID {
215-		t.Errorf("device_id: got %q, want %s", calls[0].id, hwID)
216-	}
217-
218-	// After credentials request, re-registration is rejected — callback must not fire.
219-	id, _ := deviceid.Parse(hwID)
220-	if err := fs.SetDeviceCredentials(ctx, 1, id, "1", "aa", testIP); err != nil {
221-		t.Fatal(err)
222-	}
223-	if _, err := svc.Register(ctx, "test", hwID, ""); !errors.Is(err, ErrAlreadyRegistered) {
224-		t.Errorf("post-credentials re-register: got %v, want ErrAlreadyRegistered", err)
225-	}
226-	if len(calls) != 1 {
227-		t.Errorf("OnRegistered after rejected re-register: got %d calls, want 1", len(calls))
228-	}
229-}
230-
231-func TestRegister(t *testing.T) {

$ rg -n "func newServiceFixture" -A 30 internal/pairing/service_test.go
140:func newServiceFixture(t *testing.T, cfg Config) (*Service, *fakeStore, *store.Realm) {
141-	t.Helper()
142-	sealer := newSealer(t)
143-	certPEM, sealed, err := ProvisionCA("test", 0, sealer)
144-	if err != nil {
145-		t.Fatalf("ProvisionCA: %v", err)
146-	}
147-	fs := newFakeStore()
148-	realm := fs.addRealm(&store.Realm{
149-		ID:                 1,
150-		Name:               "test",
151-		CACertificatePEM:   certPEM,
152-		CAPrivateKeySealed: sealed,
153-	})
154-	return New(fs, sealer, cfg), fs, realm
155-}
156-
157-// deviceCSR builds a fresh EC key + CSR for a device.
158-func deviceCSR(t *testing.T) string {
159-	t.Helper()
160-	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
161-	if err != nil {
162-		t.Fatal(err)
163-	}
164-	der, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
165-		Subject: pkix.Name{CommonName: "ignored"},
166-	}, key)
167-	if err != nil {
168-		t.Fatal(err)
169-	}
170-	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: der}))

← Edit internal/pairing/service.go
Index: /root/astrate-mule/internal/pairing/service.go
===================================================================
--- /root/astrate-mule/internal/pairing/service.go
+++ /root/astrate-mule/internal/pairing/service.go
@@ -218,9 +218,9 @@
 			return "", err
 		}
 	}
 	if s.OnRegistered != nil {
-		s.OnRegistered(realmName, hwID, time.Now())
+		s.OnRegistered(realmName, hwID, s.now())
 	}
 	return secret, nil
 }
 

→ Read internal/pairing/service_test.go [offset=95, limit=45]

← Edit internal/pairing/service_test.go
Index: /root/astrate-mule/internal/pairing/service_test.go
===================================================================
--- /root/astrate-mule/internal/pairing/service_test.go
+++ /root/astrate-mule/internal/pairing/service_test.go
@@ -227,8 +227,27 @@
 		t.Errorf("OnRegistered after rejected re-register: got %d calls, want 1", len(calls))
 	}
 }
 
+func TestRegisterTimestampUsesInjectedClock(t *testing.T) {
+	ctx := context.Background()
+	svc, _, _ := newServiceFixture(t, Config{})
+	hwID := randomDeviceID(t)
+
+	pinned := time.Date(2020, time.March, 4, 5, 6, 7, 0, time.UTC)
+	svc.now = func() time.Time { return pinned }
+
+	var got time.Time
+	svc.OnRegistered = func(_, _ string, at time.Time) { got = at }
+
+	if _, err := svc.Register(ctx, "test", hwID, ""); err != nil {
+		t.Fatalf("Register: %v", err)
+	}
+	if !got.Equal(pinned) {
+		t.Errorf("OnRegistered timestamp: got %v, want %v", got, pinned)
+	}
+}
+
 func TestRegister(t *testing.T) {
 	ctx := context.Background()
 	svc, fs, _ := newServiceFixture(t, Config{})
 	hwID := randomDeviceID(t)

$ gofmt -l . 2>&1 | head; go vet ./internal/pairing/ 2>&1 | tail -5; go test -race ./internal/pairing/ 2>&1 | tail -15
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/pairing	0.021s
FAIL

$ go vet ./... 2>&1 | tail -10; echo "vet exit=$?"; go test ./... 2>&1 | tail -40; echo "test exit=$?"
vet exit=0
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
ok  	github.com/astrate-platform/astrate/internal/pairing	6.350s
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
test exit=0

$ cp internal/pairing/service.go /tmp/opencode/service.go.bak && sed -i 's/s.OnRegistered(realmName, hwID, s.now())/s.OnRegistered(realmName, hwID, time.Now())/' internal/pairing/service.go && go test -run TestRegisterTimestampUsesInjectedClock ./internal/pairing/ 2>&1 | tail -6; cp /tmp/opencode/service.go.bak internal/pairing/service.go && go test -run 'TestRegisterTimestampUsesInjectedClock|TestRegisterEmitsEvent' ./internal/pairing/ 2>&1 | tail -4; git status --porcelain
--- FAIL: TestRegisterTimestampUsesInjectedClock (0.15s)
    service_test.go:246: OnRegistered timestamp: got 2026-10-07 20:07:08.42671031 +0200 CEST m=+0.294765625, want 2020-03-04 05:06:07 +0000 UTC
FAIL
FAIL	github.com/astrate-platform/astrate/internal/pairing	0.302s
FAIL
ok  	github.com/astrate-platform/astrate/internal/pairing	0.568s
 M internal/pairing/service.go
 M internal/pairing/service_test.go

$ gofmt -l internal/pairing/ ; echo "gofmt exit=$?"
gofmt exit=0

Done: Register now passes s.now() to OnRegistered instead of time.Now(), so the DeviceRegisteredEvent timestamp is clock-pinnable; added TestRegisterTimestampUsesInjectedClock (fails without the change, proven) and kept TestRegisterEmitsEvent.
Files: internal/pairing/service.go, internal/pairing/service_test.go
Verified: go vet ./... -> pass; go test ./... -> pass; go test -race ./... -> FAIL (ThreadSanitizer: unsupported VMA range, 39-bit kernel — ran without -race instead, and pairing passes as reported above); gofmt -l internal/pairing/ -> clean; reverting the impl makes the new test fail -> confirmed
 Unsure: nothing
Follow-ups: none
