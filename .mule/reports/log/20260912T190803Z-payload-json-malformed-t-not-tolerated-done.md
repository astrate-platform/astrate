slug: payload-json-malformed-t-not-tolerated
verdict: done
at:  5b56557
ran: 2026-09-12T19:08:03Z on DietPi in 99s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read pkg/payload/json.go
→ Read pkg/payload/payload.go
→ Read pkg/payload/payload_test.go
Task understood: add one `TestJSONProfile` row pinning that a malformed `t` is rejected with `ReasonBadTimestamp` even when `explicit` is false (only a *valid* `t` is tolerated-and-ignored).

← Edit pkg/payload/payload_test.go
Index: /root/astrate-mule/pkg/payload/payload_test.go
===================================================================
--- /root/astrate-mule/pkg/payload/payload_test.go
+++ /root/astrate-mule/pkg/payload/payload_test.go
@@ -350,8 +350,9 @@
 // timestamps
 {name: "t rfc3339", in: `{"v":22.5,"t":"2026-06-10T12:34:56.789Z"}`, vt: interfaceschema.Double, explicit: true, want: 22.5, wantTS: &ts},
 {name: "t epoch ms", in: `{"v":22.5,"t":1781094896789}`, vt: interfaceschema.Double, explicit: true, want: 22.5, wantTS: &ts},
 {name: "t ignored without explicit", in: `{"v":22.5,"t":"2026-06-10T12:34:56.789Z"}`, vt: interfaceschema.Double, want: 22.5},
+{name: "t bad string still rejected without explicit", in: `{"v":1,"t":"yesterday"}`, vt: interfaceschema.Double, wantReason: ReasonBadTimestamp},
 {name: "t required but missing", in: `{"v":22.5}`, vt: interfaceschema.Double, explicit: true, wantReason: ReasonBadTimestamp},
 {name: "t bad string", in: `{"v":1,"t":"yesterday"}`, vt: interfaceschema.Double, explicit: true, wantReason: ReasonBadTimestamp},
 {name: "t bool", in: `{"v":1,"t":true}`, vt: interfaceschema.Double, explicit: true, wantReason: ReasonBadTimestamp},
 {name: "t null", in: `{"v":1,"t":null}`, vt: interfaceschema.Double, explicit: true, wantReason: ReasonBadTimestamp},

$ go test ./pkg/payload/ && gofmt -l pkg/payload/
ok  	github.com/astrate-platform/astrate/pkg/payload	0.042s

$ go vet ./... && go test ./...
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/examples/flow-container-echo	[no test files]
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

Done: pinned that a malformed `t` is rejected with ReasonBadTimestamp even on a non-explicit mapping by adding one row to TestJSONProfile.
Files: pkg/payload/payload_test.go
Verified: go vet ./... && go test ./... (no -race) -> pass; gofmt -l pkg/payload/ -> clean
Unsure: nothing
Follow-ups: none
