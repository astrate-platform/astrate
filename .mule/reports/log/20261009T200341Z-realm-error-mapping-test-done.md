slug: realm-error-mapping-test
verdict: done
at:  aef0914
ran: 2026-10-09T20:03:41Z on DietPi in 175s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/realm/http.go
→ Read internal/realm/http_test.go
✱ Grep "ErrMaximumDatabaseRetentionExceeded|ErrNameCollision|ErrNameMismatch|ErrMajorMismatch|ErrValidation|errMajorNotFound" in internal/realm · 44 matches
✱ Grep "ErrAlreadyExists|ErrInterfaceMajorNotZero|ErrInterfaceInUse|ErrNotFound" in internal/store · 100 matches
✱ Grep "type ViolationsError|Violations |MappingCount|type Violation |ErrMinorNotIncreased|ErrDowngradeNotAllowed|ErrMissingEndpoints|ErrIncompatibleEndpointChange" in pkg/interfaceschema · 45 matches
→ Read pkg/interfaceschema/violations.go
→ Read internal/realm/service.go [limit=60]
✱ Grep "func WriteError|func WriteFieldErrors|func WriteNotFound|func WriteInternalServerError|WriteRawErrors|ContentType =" in pkg/astarteapi · 9 matches
→ Read pkg/astarteapi/envelope.go [offset=150, limit=90]
✱ Grep "DetailNotFound|DetailInternalServerError|DetailBadRequest|DetailForbidden" in pkg/astarteapi · 14 matches
✱ Grep "\[auto\]|mutation|implementation back" in .mule · 100 matches

$ ls -la /root/astrate-mule/.mule && ls -la /root/astrate-mule/.mule/tasks 2>/dev/null
total 16588
drwxr-xr-x  7 root root     4096 Oct  9 22:00 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       14 Oct  9 22:00 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  9 19:51 .heartbeat
-rw-r--r--  1 root root        0 Oct  9 22:00 .last-output
-rw-r--r--  1 root root        0 Oct  9 22:00 .lock
-rw-r--r--  1 root root       12 Oct  9 21:50 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 16535081 Oct  9 22:00 cron.log
drwxr-xr-x  2 root root     4096 Oct  9 20:07 failed
-rw-r--r--  1 root root   163554 Oct  9 21:05 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    16266 Oct  9 20:07 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  9 21:54 reviews
-rw-r--r--  1 root root     1586 Oct  9 22:00 task.md
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   177588 Oct  9 21:55 todo.md
total 60
drwxr-xr-x 2 root root 4096 Sep  7 21:11 .
drwxr-xr-x 7 root root 4096 Oct  9 22:00 ..
-rw-r--r-- 1 root root 3293 Sep  7 21:11 docs-sync-appengine-data-output-params.md
-rw-r--r-- 1 root root 5539 Jul 27 22:18 issue-16.md
-rw-r--r-- 1 root root 7510 Sep  4 20:46 issue-23.md
-rw-r--r-- 1 root root 6434 Sep  4 20:46 issue-24.md
-rw-r--r-- 1 root root 5665 Sep  4 20:46 issue-25.md
-rw-r--r-- 1 root root 6009 Sep  4 20:46 issue-26.md
-rw-r--r-- 1 root root 5683 Sep  4 20:46 issue-27.md

$ ls /root/astrate-mule/internal/realm/
async_operation_test.go
dashboard_compat_test.go
detailed.go
detailed_required_encrypted_test.go
detailed_test.go
helpers_test.go
http.go
http_test.go
policy.go
policy_test.go
service.go

$ cat /root/astrate-mule/internal/flowapi/http.go | sed -n '290,310p'
_ = astarteapi.WriteInternalServerError(w)
	}
}

// resolveAutoRestart returns the effective auto_restart value from a request
// body. Missing (nil) defaults to true; a present value is honored as-is.
func resolveAutoRestart(p *bool) bool {
	if p == nil {
		return true
	}
	return *p
}

func validationDetail(err error) string {
	msg := err.Error()
	const prefix = "flowapi: validation failed: "
	if len(msg) > len(prefix) && msg[:len(prefix)] == prefix {
		return msg[len(prefix):]
	}
	return msg
}

→ Read internal/realm/helpers_test.go

$ cd /root/astrate-mule && rg -n "func Test|func callWriter|renderViolationsBody|validationDetail|type rmWriter" internal/realm/*_test.go | rg -v "http_test.go"
internal/realm/helpers_test.go:12:func TestNormaliseIfaceName(t *testing.T) {
internal/realm/helpers_test.go:34:func TestJoinPEM(t *testing.T) {
internal/realm/helpers_test.go:53:func TestTriggerErrorBody(t *testing.T) {
internal/realm/helpers_test.go:107:func TestMajorParam(t *testing.T) {
internal/realm/async_operation_test.go:20:func TestRealmManagementAsyncOperationParam(t *testing.T) {
internal/realm/dashboard_compat_test.go:28:func TestDashboardCompat(t *testing.T) {
internal/realm/dashboard_compat_test.go:166:func TestPolicyReferentialIntegrity(t *testing.T) {
internal/realm/policy_test.go:8:func TestValidatePolicy(t *testing.T) {
internal/realm/detailed_test.go:15:func TestRealmManagementDetailedListing(t *testing.T) {
internal/realm/detailed_test.go:49:func TestRenderDetailedInterfaceProperties(t *testing.T) {
internal/realm/detailed_required_encrypted_test.go:12:func TestRenderDetailedInterfaceRequiredEncrypted(t *testing.T) {

← Write internal/realm/errors_test.go
Wrote file successfully.

$ cd /root/astrate-mule && gofmt -l internal/realm/errors_test.go && go vet ./internal/realm/ && go test ./internal/realm/
ok  	github.com/astrate-platform/astrate/internal/realm	0.034s

$ cd /root/astrate-mule && gofmt -l . && go vet ./... && go test ./...
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

Done: added `internal/realm/errors_test.go`, a container-free suite pinning `writeInterfaceError`, `writeError`, `renderViolationsBody` (byte-for-byte) and `validationDetail`.
Files: internal/realm/errors_test.go
Verified: `gofmt -l .` -> clean; `go vet ./...` -> pass; `go test ./...` -> all packages ok (ran without `-race`; no goroutines touched)
Unsure: none. Test-only diff; no `a.svc` deref, so `*API{}` suffices as the task stated.
Follow-ups: none
