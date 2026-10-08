slug: testutil-wait-since-cursor
verdict: done
at:  3654fda
ran: 2026-10-08T17:37:21Z on DietPi in 319s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/testutil/astartedevice.go
✱ Grep "func MQTTConnect|func WaitToken" in internal/testutil · 2 matches
✱ Glob "internal/testutil/*.go" 6 matches

$ rg -n "type Message interface" -A 15 $(go env GOMODCACHE)/github.com/eclipse/paho.mqtt.golang*/message.go 2>/dev/null || rg -n "paho.mqtt.golang" go.mod
31:type Message interface {
32-	Duplicate() bool
33-	Qos() byte
34-	Retained() bool
35-	Topic() string
36-	MessageID() uint16
37-	Payload() []byte
38-	Ack()
39-}
40-
41-type message struct {
42-	duplicate bool
43-	qos       byte
44-	retained  bool
45-	topic     string
46-	messageID uint16

→ Read internal/testutil/mqttclient.go [offset=130, limit=40]

$ rg -n "WaitForTopic|WaitForMessage|\.Mark\(" --type go | head -40
test/conformance/atomvm/atomvm_test.go:95:		msg := dev.WaitForTopic(t, 15*time.Second, dev.Base()+"/"+ifServerData+"/value")
internal/engine/engine_e2e_test.go:333:	msg := dev.WaitForTopic(t, 5*time.Second, topic)
internal/engine/engine_e2e_test.go:362:	dev.WaitForTopic(t, 5*time.Second, propTopic) // initial retained set
internal/engine/engine_e2e_test.go:367:	msg := dev.WaitForTopic(t, 5*time.Second, propTopic)
internal/engine/engine_e2e_test.go:374:	cp := dev.WaitForMessage(t, 5*time.Second, "consumer/properties listing "+want, func(m testutil.ServerMessage) bool {
internal/testutil/astartedevice.go:147:// WaitForMessage polls until a captured message satisfies pred, returning it.
internal/testutil/astartedevice.go:149:func (d *AstarteDevice) WaitForMessage(t testing.TB, timeout time.Duration, what string, pred func(ServerMessage) bool) ServerMessage {
internal/testutil/astartedevice.go:164:// WaitForTopic waits for a (non-empty) message delivered on an exact topic.
internal/testutil/astartedevice.go:165:func (d *AstarteDevice) WaitForTopic(t testing.TB, timeout time.Duration, topic string) ServerMessage {
internal/testutil/astartedevice.go:167:	return d.WaitForMessage(t, timeout, "message on "+topic, func(m ServerMessage) bool {
internal/appengine/e2e_test.go:196:		dev.WaitForTopic(t, 5*time.Second, dev.Base()+"/"+cdServerData+"/value")

→ Read internal/testutil/golden_test.go [limit=30]

$ rg -n "paho.Message|fakeMessage|testMessage" --type go internal/ test/ | head -20; rg -n "package testutil" internal/testutil/*_test.go
internal/flow/blocks/mqtt.go:185:	src := &mqttSource{name: name, ch: make(chan paho.Message, mqttSourceQueueCap)}
internal/flow/blocks/mqtt.go:193:	handler := func(_ paho.Client, m paho.Message) {
internal/flow/blocks/mqtt.go:217:	ch       chan paho.Message
internal/flow/blocks/mqtt.go:273:func mqttSourceMessage(m paho.Message) *flow.Message {
internal/flow/blocks/mqtt_test.go:96:func waitForMessage(t *testing.T, ch chan paho.Message, within time.Duration) paho.Message {
internal/flow/blocks/mqtt_test.go:222:	received := make(chan paho.Message, 8)
internal/flow/blocks/mqtt_test.go:224:	token := subscriber.Subscribe("t/sink", 0, func(_ paho.Client, m paho.Message) {
internal/flow/blocks/mqtt_test.go:287:	received := make(chan paho.Message, 8)
internal/flow/blocks/mqtt_test.go:289:	token := late.Subscribe("t/retained", 0, func(_ paho.Client, m paho.Message) {
internal/testutil/astartedevice.go:62:		opts.SetDefaultPublishHandler(func(_ paho.Client, m paho.Message) {
internal/broker/broker_test.go:315:		token := client.Subscribe(base+"/#", 1, func(_ paho.Client, m paho.Message) {
internal/broker/broker_e2e_test.go:448:		opts.SetDefaultPublishHandler(func(_ paho.Client, m paho.Message) {
internal/broker/broker_e2e_test.go:482:		opts.SetDefaultPublishHandler(func(_ paho.Client, m paho.Message) {
internal/broker/broker_e2e_test.go:551:		opts.SetDefaultPublishHandler(func(_ paho.Client, m paho.Message) {
internal/testutil/golden_test.go:1:package testutil
internal/testutil/pg_test.go:3:package testutil

Now I'll make the change: extract the capture handler into a method, add `Mark`/`From` wait variants, and write the container-free test.
← Edit internal/testutil/astartedevice.go
Index: /root/astrate-mule/internal/testutil/astartedevice.go
===================================================================
--- /root/astrate-mule/internal/testutil/astartedevice.go
+++ /root/astrate-mule/internal/testutil/astartedevice.go
@@ -57,16 +57,11 @@
 // §3.2) — so every server-owned and control message is captured.
 func ConnectAstarteDevice(t testing.TB, brokerURL, realm string, id deviceid.ID, tlsCfg *tls.Config, cleanSession bool) *AstarteDevice {
 	t.Helper()
 	d := &AstarteDevice{Realm: realm, ID: id, base: realm + "/" + id.String()}
-	collect := func(opts *paho.ClientOptions) {
-		opts.SetDefaultPublishHandler(func(_ paho.Client, m paho.Message) {
-			d.mu.Lock()
-			d.received = append(d.received, ServerMessage{Topic: m.Topic(), Payload: append([]byte(nil), m.Payload()...)})
-			d.mu.Unlock()
-		})
-	}
-	client, _ := MQTTConnect(t, brokerURL, d.base, cleanSession, tlsCfg, collect)
+	client, _ := MQTTConnect(t, brokerURL, d.base, cleanSession, tlsCfg, func(opts *paho.ClientOptions) {
+		opts.SetDefaultPublishHandler(d.capture)
+	})
 	d.Client = client
 	WaitToken(t, client.Subscribe(d.base+"/#", 2, nil), 5*time.Second)
 	return d
 }

← Edit internal/testutil/astartedevice.go
Index: /root/astrate-mule/internal/testutil/astartedevice.go
===================================================================
--- /root/astrate-mule/internal/testutil/astartedevice.go
+++ /root/astrate-mule/internal/testutil/astartedevice.go
@@ -138,17 +138,50 @@
 	defer d.mu.Unlock()
 	return append([]ServerMessage(nil), d.received...)
 }
 
+// capture is the paho publish handler that records server-owned messages.
+// It is a method rather than a closure so tests can drive captures without a
+// live broker.
+func (d *AstarteDevice) capture(_ paho.Client, m paho.Message) {
+	d.mu.Lock()
+	d.received = append(d.received, ServerMessage{Topic: m.Topic(), Payload: append([]byte(nil), m.Payload()...)})
+	d.mu.Unlock()
+}
+
+// Mark returns the current capture cursor: the number of messages captured so
+// far. Pass it as the start index of a From wait so the wait only matches
+// messages that arrive after the mark — a repeated wait for the same topic
+// cannot re-match an older message.
+func (d *AstarteDevice) Mark() int {
+	d.mu.Lock()
+	defer d.mu.Unlock()
+	return len(d.received)
+}
+
 // WaitForMessage polls until a captured message satisfies pred, returning it.
 // It fails the test on timeout.
 func (d *AstarteDevice) WaitForMessage(t testing.TB, timeout time.Duration, what string, pred func(ServerMessage) bool) ServerMessage {
 	t.Helper()
+	return d.WaitForMessageFrom(t, timeout, 0, what, pred)
+}
+
+// WaitForMessageFrom is WaitForMessage restricted to messages captured at or
+// after index start (typically a Mark() cursor), so a repeated wait does not
+// re-match an older message.
+func (d *AstarteDevice) WaitForMessageFrom(t testing.TB, timeout time.Duration, start int, what string, pred func(ServerMessage) bool) ServerMessage {
+	t.Helper()
+	if start < 0 {
+		start = 0
+	}
 	deadline := time.Now().Add(timeout)
 	for time.Now().Before(deadline) {
-		for _, m := range d.Messages() {
-			if pred(m) {
-				return m
+		msgs := d.Messages()
+		if start < len(msgs) {
+			for _, m := range msgs[start:] {
+				if pred(m) {
+					return m
+				}
 			}
 		}
 		time.Sleep(10 * time.Millisecond)
 	}
@@ -158,9 +191,17 @@
 
 // WaitForTopic waits for a (non-empty) message delivered on an exact topic.
 func (d *AstarteDevice) WaitForTopic(t testing.TB, timeout time.Duration, topic string) ServerMessage {
 	t.Helper()
-	return d.WaitForMessage(t, timeout, "message on "+topic, func(m ServerMessage) bool {
+	return d.WaitForTopicFrom(t, timeout, 0, topic)
+}
+
+// WaitForTopicFrom is WaitForTopic restricted to messages captured at or after
+// index start (typically a Mark() cursor), so a repeated wait for the same
+// topic returns the next message instead of re-matching an older one.
+func (d *AstarteDevice) WaitForTopicFrom(t testing.TB, timeout time.Duration, start int, topic string) ServerMessage {
+	t.Helper()
+	return d.WaitForMessageFrom(t, timeout, start, "message on "+topic, func(m ServerMessage) bool {
 		return m.Topic == topic && len(m.Payload) > 0
 	})
 }
 

Now the container-free test:
← Write internal/testutil/astartedevice_test.go
Wrote file successfully.

$ gofmt -l . && go vet ./... && go test -race ./internal/testutil/ ./internal/engine/ 2>&1 | tail -20
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/testutil	0.025s
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/engine	0.022s
FAIL

Race detector unavailable here (39-bit VMA kernel, as MULE.md notes). Running the standard gate:
$ gofmt -l . && go vet ./... && go test ./... 2>&1 | tail -30
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
ok  	github.com/astrate-platform/astrate/internal/testutil	0.065s
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

$ go test -count=1 -run TestWaitForTopicFromAfterMark -v ./internal/testutil/
=== RUN   TestWaitForTopicFromAfterMark
--- PASS: TestWaitForTopicFromAfterMark (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/internal/testutil	0.069s

Done: AstarteDevice capture handler is now an injectable method with `Mark()` + `WaitForMessageFrom`/`WaitForTopicFrom` cursor waits, plus a container-free test proving a post-mark wait returns the second message.
Files: internal/testutil/astartedevice.go, internal/testutil/astartedevice_test.go
Verified: `gofmt -l .` (clean), `go vet ./...`, `go test ./...` → pass; `go test -race` → cannot run here (39-bit VMA kernel, TSan abort), so the change is **unverified for races**.
Unsure: nothing — the test fails to compile without the new API, which is the intended failure proof.
Follow-ups: none (consumer task engine-e2e-emptycache-resync-assert already exists as the companion line).
