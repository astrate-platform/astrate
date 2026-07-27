code: cc20099af379d72b0d71208ee0259d6a404ddc90
at:  ad335ec
ran: 2026-07-27T16:15:28Z on DietPi in 71s
toolchain: go version go1.26.5 linux/amd64


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read .mule/recipes/legion-go.md
This is a `[legion]` `[readonly]` task: run the race detector on the Legion Go and report failures to `.mule/for-giulio.md`. Let me first check if Legion is reachable.

$ ssh legion 'echo reachable' 2>&1
reachable

Legion is reachable. Running the race check:
$ ssh legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...' 2>&1
go: downloading github.com/testcontainers/testcontainers-go v0.43.0
go: downloading github.com/moby/moby/api v1.54.2
go: downloading github.com/shirou/gopsutil/v4 v4.26.5
ok  	github.com/astrate-platform/astrate/cmd/astrate	1.489s
?   	github.com/astrate-platform/astrate/docs	[no test files]
ok  	github.com/astrate-platform/astrate/internal/appengine	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)
ok  	github.com/astrate-platform/astrate/internal/appengine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/auth	1.155s
ok  	github.com/astrate-platform/astrate/internal/broker	2.909s
ok  	github.com/astrate-platform/astrate/internal/config	1.016s
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	1.017s
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	2.359s
?   	github.com/astrate-platform/astrate/internal/housekeeping	[no test files]
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
?   	github.com/astrate-platform/astrate/internal/swagger	[no test files]
ok  	github.com/astrate-platform/astrate/internal/testutil	1.063s
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	1.064s
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

All tests passed with no race conditions detected. Since this is a `[readonly]` task and there were no failures, nothing needs to be written to `.mule/for-giulio.md` (the task specifies reporting only failures).

Done: Race check passed cleanly on Legion Go — all packages green, zero data races.
Files: none (readonly task)
Verified: `ssh legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...'` -> pass
Unsure: nothing
Follow-ups: none
