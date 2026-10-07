slug: pairing-unregister-bad-id-test
verdict: done
at:  322242a
ran: 2026-10-07T18:47:17Z on DietPi in 225s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

$ rg -n "func (s \*service|) Unregister|ErrNotFound" internal/pairing/service.go | head -50; echo ---; rg -n "TestUnregister" internal/pairing/service_test.go
229:// credential trail). store.ErrNotFound is returned for unknown devices.
233:		return fmt.Errorf("%w: device %q", store.ErrNotFound, deviceIDStr)
371:		return err // store.ErrNotFound → 404
387:	if errors.Is(err, store.ErrNotFound) {
395:	if errors.Is(err, store.ErrNotFound) {
---
359:func TestUnregister(t *testing.T) {

→ Read internal/pairing/service.go [offset=200, limit=60]
→ Read internal/pairing/service_test.go [offset=300, limit=120]

$ rg -n "func Parse" -A 40 internal/deviceid/*.go | head -80
rg: internal/deviceid/*.go: IO error for operation on internal/deviceid/*.go: No such file or directory (os error 2)

$ rg -rn --glob '!*_test.go' "func Parse" | rg -i deviceid; echo ---; rg -n "deviceid\"" internal/pairing/service.go
pkg/deviceid/deviceid.go:n(s string) (ID, error) {
---
22:	"github.com/astrate-platform/astrate/pkg/deviceid"

→ Read pkg/deviceid/deviceid.go
→ Read internal/pairing/service_test.go [limit=90]

$ rg -n "func randomDeviceID|func deviceCSR" internal/pairing/service_test.go
158:func deviceCSR(t *testing.T) string {
173:func randomDeviceID(t *testing.T) string {

← Edit internal/pairing/service_test.go
Index: /root/astrate-mule/internal/pairing/service_test.go
===================================================================
--- /root/astrate-mule/internal/pairing/service_test.go
+++ /root/astrate-mule/internal/pairing/service_test.go
@@ -360,10 +360,22 @@
 ctx := context.Background()
 svc, fs, _ := newServiceFixture(t, Config{})
 hwID := randomDeviceID(t)
 
-if err := svc.Unregister(ctx, "test", hwID); !errors.Is(err, store.ErrNotFound) {
-	t.Errorf("unknown device: got %v, want store.ErrNotFound", err)
+// Unregister turns every lookup failure into store.ErrNotFound: both
+// deviceid.Parse rejections (wrong length, wrong alphabet) before the
+// store is touched, and a well-formed ID the store does not know.
+for _, tc := range []struct {
+	name string
+	id   string
+}{
+	{"malformed id, wrong length", "not-a-device-id"},
+	{"malformed id, wrong alphabet", "h4+Dx/RYTU+RbpDOTabhRg"},
+	{"well-formed unknown id", randomDeviceID(t)},
+} {
+	if err := svc.Unregister(ctx, "test", tc.id); !errors.Is(err, store.ErrNotFound) {
+		t.Errorf("Unregister(%s = %q): got %v, want store.ErrNotFound", tc.name, tc.id, err)
+	}
 }
 
 secret, err := svc.Register(ctx, "test", hwID, "")
 if err != nil {

$ gofmt -l . && echo "gofmt-clean" && go vet ./... && go test -race ./internal/pairing/ 2>&1 | tail -20
gofmt-clean
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/pairing	0.023s
FAIL

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
ok  	github.com/astrate-platform/astrate/internal/pairing	6.489s
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

$ go test -run TestUnregister -coverprofile=/tmp/cov.out ./internal/pairing/ >/dev/null 2>&1; go tool cover -func=/tmp/cov.out | rg "Unregister"
github.com/astrate-platform/astrate/internal/pairing/http.go:134:	handleUnregister	0.0%
github.com/astrate-platform/astrate/internal/pairing/service.go:230:	Unregister		85.7%

$ awk -F: '/service.go/ {print}' /tmp/cov.out | awk '$3>0 {print "hit:", $1":"$2} $3==0 {print "miss:", $1":"$2}' | rg "23[0-9]"
hit: github.com/astrate-platform/astrate/internal/pairing/service.go:230.88,232.16:2
hit: github.com/astrate-platform/astrate/internal/pairing/service.go:232.16,234.3:1
hit: github.com/astrate-platform/astrate/internal/pairing/service.go:235.2,236.16:2
miss: github.com/astrate-platform/astrate/internal/pairing/service.go:236.16,238.3:1
hit: github.com/astrate-platform/astrate/internal/pairing/service.go:239.2,239.49:1

Done: TestUnregister now asserts the malformed-ID branches (wrong length, wrong alphabet) and the well-formed unknown-ID 404 in one table.
Files: internal/pairing/service_test.go (TestUnregister, :359)
Verified: gofmt -l . -> clean; go vet ./... -> pass; go test ./... -> pass (`go test -race` unusable: ThreadSanitizer "unsupported VMA range / 39-bit kernel", so non-race gate run per MULE.md); coverage probe: service.go:232-234 hit by the new rows (only GetRealmByName's error return still uncovered)
Unsure: no production code changed, so no test that fails without the change exists — this is a coverage-only task; the table passes against unmodified service.go by construction. "Both bad-ID shapes" read as the two deviceid.Parse rejections (length, alphabet).
Follow-ups: none
