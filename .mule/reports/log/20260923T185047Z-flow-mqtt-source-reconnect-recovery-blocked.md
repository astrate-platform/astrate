slug: flow-mqtt-source-reconnect-recovery
verdict: blocked
at:  85db760
ran: 2026-09-23T18:50:47Z on DietPi in 84s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/flow/blocks/mqtt.go
→ Read internal/flow/blocks/mqtt_test.go
→ Read go.mod
✱ Glob "internal/flow/*.go" 18 matches
→ Read internal/flow/block.go

$ ls ~/go/pkg/mod/github.com/eclipse/ 2>/dev/null; go env GOMODCACHE 2>/dev/null
paho.mqtt.golang@v1.5.1
/root/go/pkg/mod

$ ls /root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/
CODE_OF_CONDUCT.md
CONTRIBUTING.md
LICENSE
NOTICE.md
README.md
SECURITY.md
backoff.go
backoff_test.go
client.go
client_test.go
cmd
components.go
connnotf.go
edl-v10
epl-v20
filestore.go
fvt
fvt_client_test.go
fvt_store_test.go
fvt_test.go
go.mod
go.sum
memstore.go
memstore_ordered.go
message.go
messageids.go
net.go
netconn.go
oops.go
options.go
options_reader.go
options_test.go
packets
ping.go
router.go
status.go
store.go
token.go
token_test.go
topic.go
trace.go
unit_client_test.go
unit_message_test.go
unit_messageids_test.go
unit_options_test.go
unit_ping_test.go
unit_router_test.go
unit_status_test.go
unit_store_test.go
unit_topic_test.go
websocket.go

$ rg -n "func \(c \*client\) (reconnect|connectedHandler|internalConnLost|connectionUp|IsConnected)|SetOnConnectHandler|onConnectHandler|resumeSubs|SetResumeSubs|SetAutoReconnect" /root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/client.go /root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/options.go
/root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/client.go:197:func (c *client) IsConnected() bool {
/root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/client.go:308:func (c *client) reconnect(connectionUp connCompletedFn) {
/root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/client.go:459:// `SetAutoReconnect` and/or `SetConnectRetry`options instead of implementing this yourself.
/root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/client.go:530:func (c *client) internalConnLost(whyConnLost error) {
/root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/client.go:853:			// if not connected and resumeSubs not set this sub will be thrown away
/root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/client.go:1165:			// if not connected and resumeSubs not set this unsub will be thrown away
/root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/options.go:186:// SetResumeSubs will enable resuming of stored (un)subscribe messages when connecting
/root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/options.go:188:func (o *ClientOptions) SetResumeSubs(resume bool) *ClientOptions {
/root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/options.go:334:// SetOnConnectHandler sets the function to be called when the client is connected. Both
/root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/options.go:336:func (o *ClientOptions) SetOnConnectHandler(onConn OnConnectHandler) *ClientOptions {
/root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/options.go:394:// SetAutoReconnect sets whether the automatic reconnection logic should be used
/root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/options.go:397:func (o *ClientOptions) SetAutoReconnect(a bool) *ClientOptions {

! permission requested: external_directory (/root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/*); auto-rejecting
✗ Read /root/go/pkg/mod/github.com/eclipse/paho.mqtt.golang@v1.5.1/client.go failed [offset=190, limit=180]
Error: The user rejected permission to use this specific tool call.
