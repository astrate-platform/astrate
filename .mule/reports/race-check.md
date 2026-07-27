code: f555e2f2bc98ae5c7cd7deba5030ab4c54693b63
at:  7099f80
ran: 2026-07-27T09:16:38Z on DietPi in 82s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read .mule/recipes/legion-go.md
The task is a `[legion] [readonly]` race-check: SSH into Legion Go, fetch/merge, run `go test -race ./...`, and report any failures. Let me check if it's reachable first.

$ ssh legion 'echo reachable' 2>&1
reachable

Good, Legion Go is up. Running the race check now.
$ ssh legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./... 2>&1'
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

All packages passed with `-race` — no data races detected (all cached, clean merge to origin/main).

**Done:** race-check on Legion Go — all tests pass, no races detected
**Files:** none changed (readonly task)
**Verified:** `ssh legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...'` → all ok, cached
**Unsure:** nothing
**Follow-ups:** none
