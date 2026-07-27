code: 9c9e6b89dfee6fbb71cd4d08a0cee94888d74ce5
at:  13cd7f1
ran: 2026-07-27T10:34:09Z on DietPi in 59s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
This is a `[legion] [readonly]` task: run `go test -race ./...` on the Legion Go and report failures. Let me try to reach Legion.

$ ssh -o ConnectTimeout=5 -o BatchMode=yes legion "echo reachable" 2>&1
reachable

Legion is reachable. Running the race test now.
$ ssh legion "cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./... 2>&1; echo EXIT_CODE:\$?"
fish: $? is not the exit status. In fish, please use $status.
cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./... 2>&1; echo EXIT_CODE:$?
                                                                                                                ^

Legion uses fish shell. Let me adjust.
$ ssh legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./... 2>&1; echo EXIT_CODE:$status'
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
EXIT_CODE:0

All packages passed, no races detected. Clean run.

```
Done: Race check on Legion Go — all 21 packages passed, no races found
Files: none (readonly task, no changes made)
Verified: ssh legion "cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./..." -> pass (exit 0)
Unsure: nothing
Follow-ups: none
```
