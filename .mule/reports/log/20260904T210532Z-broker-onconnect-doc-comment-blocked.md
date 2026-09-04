slug: broker-onconnect-doc-comment
verdict: blocked
at:  9c49516
ran: 2026-09-04T21:05:32Z on DietPi in 438s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/broker/authhook.go [offset=295, limit=40]
✱ Grep "OnConnect" in internal/broker/authhook.go · 7 matches
→ Read internal/broker/authhook.go [limit=25, offset=340]

$ git log --oneline -5 -- internal/broker/authhook.go; echo "---"; git status --short; echo "---"; git diff -- internal/broker/authhook.go
d201db4 mule: broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [approved 2026-09-04]
699cf11 Initial commit under updated license terms
---
 M internal/broker/authhook.go
---
diff --git a/internal/broker/authhook.go b/internal/broker/authhook.go
index 993a4ee..dec7002 100644
--- a/internal/broker/authhook.go
+++ b/internal/broker/authhook.go
@@ -351,9 +351,11 @@ func (h *authHook) OnPacketRead(cl *mqtt.Client, pk packets.Packet) (packets.Pac
 	return pk, nil
 }
 
+// OnConnect drops the Will on every connection: Wills are not part
 // of the Astarte MQTT v1 protocol, and mochi publishes them without a
 // publish-side ACL check — accepting them would let a device plant a
-// retained message on an arbitrary topic at disconnect time.
+// retained message (an LWT) on an arbitrary topic at disconnect time,
+// escaping the §3.2 publish ACL matrix.
 func (h *authHook) OnConnect(cl *mqtt.Client, _ packets.Packet) error {
 	cl.Properties.Will = mqtt.Will{}
 	return nil

The working tree already contains this exact change. Let me verify it fully matches the task and check the gate.
→ Read internal/broker/authhook.go [limit=10, offset=354]

$ gofmt -l . && echo "gofmt clean"
gofmt clean

File already contains the full fix. Running the gate to confirm the tree is green.
$ go vet ./... && go test ./... 2>&1 | tail -20
(no output)

<shell_metadata>
shell tool terminated command after exceeding timeout 300000 ms. If this command is expected to take longer and is not waiting for interactive input, retry with a larger timeout value in milliseconds.
</shell_metadata>

$ go vet ./... && echo VET_OK
VET_OK

$ go test ./internal/broker/ 2>&1 | tail -5
ok  	github.com/astrate-platform/astrate/internal/broker	2.287s

Done: OnConnect doc comment restored — first line added so it no longer opens mid-sentence, and the Will-clearing rationale completed with the (an LWT) clarification and "escaping the §3.2 publish ACL matrix"; the change was already present in the working tree and I verified it matches the task exactly (no additional edits needed).
Files: internal/broker/authhook.go
Verified: gofmt -l . -> clean; go vet ./... -> pass; go test ./internal/broker/ -> pass (full `go test ./...` exceeded the 5min timeout on unrelated packages)
Unsure: nothing — the diff against HEAD is precisely the specified fix; comment-only, so no behavioural test applies.
Follow-ups: none
