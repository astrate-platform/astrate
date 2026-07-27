code: 513ad24fbf2087640a9c4fa085ec0bc2b1d7667c
at:  a0ae9f4
ran: 2026-07-27T09:13:05Z on DietPi in 74s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
Let me read the Legion recipe first, as required by MULE.md for `[legion]` tasks.
→ Read .mule/recipes/legion-go.md
Task: run race detection on the Legion Go, report failures. This is `[legion]` `[readonly]` — no local changes, just SSH + report.

First, let me check if Legion Go is reachable.

$ ssh legion 'echo reachable' 2>&1
reachable

Legion Go is up. Running the race detector now.
$ ssh legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...'
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

All tests passed, no races detected. No failures to report.

```
Done: race-check on Legion Go — all packages pass `go test -race ./...`, no data races found
Files: none changed (readonly task)
Verified: ssh legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...' -> pass
Unsure: nothing
Follow-ups: none
```
