code: befc43a26a26c602008cf425f41d79fdc737ee99
at:  57338db
ran: 2026-07-27T08:56:39Z on DietPi in 42s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
This is a `[legion] [readonly]` task: run `go test -race ./...` on the Legion Go machine and report any race failures. Let me SSH in and do it.

$ ssh legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...' 2>&1
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
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
?   	github.com/astrate-platform/astrate/internal/housekeeping	[no test files]
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
?   	github.com/astrate-platform/astrate/internal/swagger	[no test files]
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

All tests pass with `-race`. No failures, no race conditions detected. This is a `[readonly]` task — nothing to change. No failures to report to `.mule/for-giulio.md`.

```
Done: Race-check on Legion Go — all tests pass, no races detected
Files: none changed
Verified: ssh legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...' -> pass (25 packages, all cached or ok)
Unsure: nothing
Follow-ups: none
```
