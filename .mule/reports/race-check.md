code: f12b01a4681ce398918e1c1575ce4cbda7921fb0
at:  cc13242
ran: 2026-07-27T17:35:37Z on DietPi in 72s
toolchain: go version go1.26.5 linux/amd64


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
The task is a `[legion] [readonly]` task: run the race detector on the Legion Go. Let me check if it's reachable.

$ ssh -o ConnectTimeout=5 -o BatchMode=yes legion 'echo ok' 2>&1
ok

Legion Go is up. Running the race check.
$ ssh -o ConnectTimeout=10 legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./... 2>&1' 2>&1
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

All packages pass with `-race`. No data races detected. The full run was cached (nothing changed since last run).

```
Done: race check on Legion Go — all packages pass, no data races detected
Files: none changed
Verified: ssh legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...' -> pass
Unsure: nothing
Follow-ups: none
```
