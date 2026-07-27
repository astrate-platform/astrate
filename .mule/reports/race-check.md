code: 06dd2cf21ab5e334b31791a28fb762596126bc82
at:  344d013
ran: 2026-07-27T15:08:25Z on DietPi in 48s
toolchain: go version go1.26.5 linux/amd64


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
This is a `[legion]` [readonly] task: SSH into the Legion Go, fetch/merge main, run `go test -race ./...`, and report failures. Let me read the Legion recipe first.
→ Read .mule/recipes/legion-go.md
→ Read .mule/for-giulio.md
This is a `[legion]` `[readonly]` task: SSH into the Legion Go, run the race-detector test suite, and report failures to `for-giulio.md`. Executing now.

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

All tests passed with `-race`, no failures, no races detected. Nothing to report to `for-giulio.md`.

Done: Race-check on Legion Go — all packages pass with -race, no data races found.
Files: none changed (readonly task)
Verified: `ssh legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...'` -> pass (26 packages, all ok or cached)
Unsure: nothing
Follow-ups: none
