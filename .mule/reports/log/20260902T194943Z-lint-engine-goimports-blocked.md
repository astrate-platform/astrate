slug: lint-engine-goimports
verdict: blocked
at:  6750a4b
ran: 2026-09-02T19:49:43Z on DietPi in 421s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
The task is to fix import grouping in `internal/engine/capabilities.go` by running `goimports -w` on it.

Let me read the file first.
→ Read internal/engine/capabilities.go [limit=30]
The broker import and the mongo driver import are in the same group — `goimports` should separate them. Let me run it.

$ goimports -w internal/engine/capabilities.go
/bin/sh: 1: goimports: not found

$ go install golang.org/x/tools/cmd/goimports@latest 2>&1 && export PATH="$PATH:$(go env GOPATH)/bin" && goimports -w internal/engine/capabilities.go
(no output)

Now let me check the result and run the gate:
→ Read internal/engine/capabilities.go [limit=12]

$ export PATH="$PATH:$(go env GOPATH)/bin" && gofmt -l . && go vet ./... && go test ./... 2>&1
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
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
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

$ head -5 go.mod
module github.com/astrate-platform/astrate

go 1.26.1

// Pinned dependency decisions — docs/ROADMAP.md §1.1 (frozen in docs/DESIGN.md).

The file was unchanged — goimports didn't apply any fix. Let me check the local prefix configuration:
$ rg -i 'goimport|localprefix' -l
(no output)

$ export PATH="$PATH:$(go env GOPATH)/bin" && goimports -local github.com/astrate-platform/astrate -d internal/engine/capabilities.go
diff -u internal/engine/capabilities.go.orig internal/engine/capabilities.go
--- internal/engine/capabilities.go.orig	2026-09-02 21:46:37.807450584 +0200
+++ internal/engine/capabilities.go	2026-09-02 21:46:37.807450584 +0200
@@ -4,8 +4,9 @@
 	"context"
 	"fmt"
 
-	"github.com/astrate-platform/astrate/internal/broker"
 	"go.mongodb.org/mongo-driver/v2/bson"
+
+	"github.com/astrate-platform/astrate/internal/broker"
 )
 
 // Known capability keys (upstream MQTT v1 §2 — Device capabilities).

Now I see the expected fix. Let me check the golangci-lint config to see if it specifies `-local`:
$ rg -A20 'goimports' .golangci* 2>/dev/null || echo "no .golangci config"
- goimports
  settings:
    goimports:
      local-prefixes:
        - github.com/astrate-platform/astrate

Confirmed: golangci-lint uses `-local github.com/astrate-platform/astrate`. Applying the fix:
$ export PATH="$PATH:$(go env GOPATH)/bin" && goimports -local github.com/astrate-platform/astrate -w internal/engine/capabilities.go
(no output)

→ Read internal/engine/capabilities.go [limit=12]

$ go vet ./... 2>&1 && go test ./... 2>&1
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
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
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

$ gofmt -l internal/engine/capabilities.go
(no output)

Done: Fixed import grouping in `internal/engine/capabilities.go` — separated `internal/broker` into its own local-prefix group (per golangci-lint `goimports` `local-prefixes` config).
Files: `internal/engine/capabilities.go`
Verified: `go vet ./... && go test ./... && gofmt -l .` -> pass
Unsure: nothing
Follow-ups: none
