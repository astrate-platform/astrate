slug: pairing-remoteip-fallback-test
verdict: done
at:  3d6e790
ran: 2026-09-09T16:54:02Z on DietPi in 469s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/pairing/http.go
Task: add a container-free unit test for `remoteIP` in internal/pairing/http.go:359 covering three cases (addr:port → Unmap, bare address → mapped, unparseable → IPv4Unspecified).
✱ Glob "internal/pairing/*_test.go" 4 matches
Now let me check existing test style in the pairing package.
→ Read internal/pairing/bearersecret_test.go
→ Read internal/pairing/http_test.go
Now I'll write the unit test in a new non-integration file, following the `bearersecret_test.go` table-driven style.

← Write internal/pairing/remoteip_test.go
Wrote file successfully.

$ gofmt -l . && go test ./internal/pairing/ 2>&1 | tail -20
internal/pairing/remoteip_test.go
ok  	github.com/astrate-platform/astrate/internal/pairing	13.384s

Test passes but gofmt flags the file. Fixing.
$ gofmt -w internal/pairing/remoteip_test.go && gofmt -l . && go test -count=1 -run TestRemoteIP ./internal/pairing/
ok  	github.com/astrate-platform/astrate/internal/pairing	0.383s

→ Read internal/pairing/remoteip_test.go
Formatting clean, test passes. Running the full gate.

$ go vet ./... && go test ./... 2>&1 | tail -30
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

```
Done: added container-free table-driven unit test for remoteIP covering addr:port (IPv4 + IPv4-mapped Unmap), bare-address, and unparseable → IPv4Unspecified fallback.
Files: internal/pairing/remoteip_test.go
Verified: gofmt -l . (clean); go test ./... (pass, incl. TestRemoteIP); go vet ./... (pass) -> pass
Unsure: nothing
Follow-ups: none
```
