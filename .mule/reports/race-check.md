code: 1f4b64306fc19834a1fa913b655299e71dae02bd
at:  16cc008
ran: 2026-07-27T14:03:26Z on DietPi in 108s
toolchain: go version go1.26.5 linux/amd64


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
This task is tagged `[legion]` and `[readonly]`. I need to check if I can reach Legion Go, and also read the recipe.

$ ssh -o ConnectTimeout=5 -o BatchMode=yes legion "echo reachable" 2>&1
reachable

→ Read .mule/recipes/legion-go.md
Legion is reachable. Running the race check now.

$ ssh legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...' 2>&1
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	2.848s
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
ok  	github.com/astrate-platform/astrate/internal/testutil	1.056s
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

All packages passed `go test -race ./...` with no failures. No data races detected. Since this is `[readonly]` and the gate passed, I change nothing.

Done: Race check on Legion Go — all tests pass cleanly
Files: none (readonly)
Verified: `ssh legion 'cd ~/astrate && go test -race ./...'` -> pass
Unsure: nothing
Follow-ups: none
