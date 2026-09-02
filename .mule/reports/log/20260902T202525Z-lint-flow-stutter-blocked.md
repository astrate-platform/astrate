slug: lint-flow-stutter
verdict: blocked
at:  2cbc7d3
ran: 2026-09-02T20:25:25Z on DietPi in 98s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
Got it. The task: rename `FlowStatus` → `Status`, `FlowConfig` → `Config`, `FlowMessage` → `Message` in `internal/flow` to fix revive stutter lint. Let me find all occurrences first.

$ rg -n 'FlowStatus|FlowConfig|FlowMessage' /root/astrate-mule/internal/flow/
/root/astrate-mule/internal/flow/router_test.go:19:	msgs []*FlowMessage
/root/astrate-mule/internal/flow/router_test.go:22:func (c *collectBlock) Process(msg *FlowMessage) ([]*FlowMessage, error) {
/root/astrate-mule/internal/flow/router_test.go:38:func (c *collectBlock) get(i int) *FlowMessage {
/root/astrate-mule/internal/flow/router_test.go:51:func (f *failBlock) Process(msg *FlowMessage) ([]*FlowMessage, error) {
/root/astrate-mule/internal/flow/router_test.go:56:	return []*FlowMessage{msg}, nil
/root/astrate-mule/internal/flow/router_test.go:64:func (p *passthroughBlock) Process(msg *FlowMessage) ([]*FlowMessage, error) {
/root/astrate-mule/internal/flow/router_test.go:65:	return []*FlowMessage{msg}, nil
/root/astrate-mule/internal/flow/router_test.go:77:func (p *panicBlock) Process(msg *FlowMessage) ([]*FlowMessage, error) {
/root/astrate-mule/internal/flow/router_test.go:82:	return []*FlowMessage{msg}, nil
/root/astrate-mule/internal/flow/router_test.go:89:func makeMsg(key string, seq int) *FlowMessage {
/root/astrate-mule/internal/flow/router_test.go:90:	return &FlowMessage{
/root/astrate-mule/internal/flow/router_test.go:384:	msgs []*FlowMessage
/root/astrate-mule/internal/flow/router_test.go:387:func (b *blockingSink) Process(msg *FlowMessage) ([]*FlowMessage, error) {
/root/astrate-mule/internal/flow/message.go:74:// FlowMessage is one unit of data flowing through a block graph. Every message
/root/astrate-mule/internal/flow/message.go:77:type FlowMessage struct {
/root/astrate-mule/internal/flow/message.go:111:// MarshalJSON serialises the FlowMessage to the upstream JSON wire format.
/root/astrate-mule/internal/flow/message.go:112:func (m *FlowMessage) MarshalJSON() ([]byte, error) {
/root/astrate-mule/internal/flow/message.go:138:// UnmarshalJSON deserialises a FlowMessage from the upstream JSON wire format.
/root/astrate-mule/internal/flow/message.go:139:func (m *FlowMessage) UnmarshalJSON(b []byte) error {
/root/astrate-mule/internal/flow/message.go:198:func (m *FlowMessage) dataWireScalar() any {
/root/astrate-mule/internal/flow/message.go:213:func (m *FlowMessage) dataWireMap() any {
/root/astrate-mule/internal/flow/message.go:239:func (m *FlowMessage) fieldTypesWire() map[string]string {
/root/astrate-mule/internal/flow/message.go:248:func (m *FlowMessage) fieldSubtypesWire() map[string]string {
/root/astrate-mule/internal/flow/message.go:260:func (m *FlowMessage) setDataFromWire(dt DataType, raw any) error {
/root/astrate-mule/internal/flow/message.go:327:func (m *FlowMessage) setDataFromWireMap(raw any) error {
/root/astrate-mule/internal/flow/flow.go:13:// FlowStatus enumerates the lifecycle states of a flow.
/root/astrate-mule/internal/flow/flow.go:14:type FlowStatus uint8
/root/astrate-mule/internal/flow/flow.go:17:	// FlowStatusCreating indicates the flow is being initialised.
/root/astrate-mule/internal/flow/flow.go:18:	FlowStatusCreating FlowStatus = iota
/root/astrate-mule/internal/flow/flow.go:19:	// FlowStatusRunning indicates the flow is accepting and processing messages.
/root/astrate-mule/internal/flow/flow.go:20:	FlowStatusRunning
/root/astrate-mule/internal/flow/flow.go:21:	// FlowStatusStopped indicates the flow has been gracefully shut down.
/root/astrate-mule/internal/flow/flow.go:22:	FlowStatusStopped
/root/astrate-mule/internal/flow/flow.go:23:	// FlowStatusFailed indicates the flow failed during initialisation.
/root/astrate-mule/internal/flow/flow.go:24:	FlowStatusFailed
/root/astrate-mule/internal/flow/flow.go:28:func (s FlowStatus) String() string {
/root/astrate-mule/internal/flow/flow.go:30:	case FlowStatusCreating:
/root/astrate-mule/internal/flow/flow.go:32:	case FlowStatusRunning:
/root/astrate-mule/internal/flow/flow.go:34:	case FlowStatusStopped:
/root/astrate-mule/internal/flow/flow.go:36:	case FlowStatusFailed:
/root/astrate-mule/internal/flow/flow.go:52:// FlowConfig holds the parameters needed to instantiate a running flow.
/root/astrate-mule/internal/flow/flow.go:53:type FlowConfig struct {
/root/astrate-mule/internal/flow/flow.go:70:	status     FlowStatus
/root/astrate-mule/internal/flow/flow.go:87:func (f *Flow) Status() FlowStatus {
/root/astrate-mule/internal/flow/flow.go:119:func (f *Flow) setStatus(s FlowStatus) { f.status = s }
/root/astrate-mule/internal/flow/flow.go:138:func (m *Manager) StartFlow(ctx context.Context, cfg FlowConfig) (*Flow, error) {
/root/astrate-mule/internal/flow/flow.go:150:		status:     FlowStatusCreating,
/root/astrate-mule/internal/flow/flow.go:157:		f.setStatus(FlowStatusFailed)
/root/astrate-mule/internal/flow/flow.go:166:	f.setStatus(FlowStatusRunning)
/root/astrate-mule/internal/flow/flow.go:190:	f.setStatus(FlowStatusStopped)
/root/astrate-mule/internal/flow/flow.go:197:// GetFlowStatus returns the current status of the flow identified by
/root/astrate-mule/internal/flow/flow.go:199:func (m *Manager) GetFlowStatus(pipelineID string) (FlowStatus, error) {
/root/astrate-mule/internal/flow/flow.go:238:		if s == FlowStatusRunning {
/root/astrate-mule/internal/flow/manager_test.go:18:	f, err := mgr.StartFlow(ctx, FlowConfig{
/root/astrate-mule/internal/flow/manager_test.go:33:	if got := f.Status(); got != FlowStatusRunning {
/root/astrate-mule/internal/flow/manager_test.go:43:	got, err := mgr.GetFlowStatus("pipe-1")
/root/astrate-mule/internal/flow/manager_test.go:45:		t.Fatalf("GetFlowStatus: %v", err)
/root/astrate-mule/internal/flow/manager_test.go:47:	if got != FlowStatusRunning {
/root/astrate-mule/internal/flow/manager_test.go:48:		t.Fatalf("GetFlowStatus = %v, want running", got)
/root/astrate-mule/internal/flow/manager_test.go:59:	f, err := mgr.StartFlow(ctx, FlowConfig{
/root/astrate-mule/internal/flow/manager_test.go:89:	if got := f.Status(); got != FlowStatusStopped {
/root/astrate-mule/internal/flow/manager_test.go:110:	f, err := mgr.StartFlow(context.Background(), FlowConfig{
/root/astrate-mule/internal/flow/manager_test.go:120:	if got := f.Status(); got != FlowStatusFailed {
/root/astrate-mule/internal/flow/manager_test.go:124:	got, err := mgr.GetFlowStatus("pipe-3")
/root/astrate-mule/internal/flow/manager_test.go:126:		t.Fatalf("GetFlowStatus: %v", err)
/root/astrate-mule/internal/flow/manager_test.go:128:	if got != FlowStatusFailed {
/root/astrate-mule/internal/flow/manager_test.go:129:		t.Fatalf("GetFlowStatus = %v, want failed", got)
/root/astrate-mule/internal/flow/manager_test.go:140:		_, err := mgr.StartFlow(ctx, FlowConfig{
/root/astrate-mule/internal/flow/manager_test.go:165:	_, err := mgr.StartFlow(ctx, FlowConfig{
/root/astrate-mule/internal/flow/manager_test.go:174:	_, err = mgr.StartFlow(ctx, FlowConfig{
/root/astrate-mule/internal/flow/manager_test.go:187:	_, err := mgr.GetFlowStatus("nope")
/root/astrate-mule/internal/flow/manager_test.go:189:		t.Fatalf("GetFlowStatus error = %v, want ErrFlowNotFound", err)
/root/astrate-mule/internal/flow/manager_test.go:213:	f, err := mgr.StartFlow(ctx, FlowConfig{
/root/astrate-mule/internal/flow/manager_test.go:261:func (o *orderedCollectBlock) Process(msg *FlowMessage) ([]*FlowMessage, error) {
/root/astrate-mule/internal/flow/manager_test.go:281:		_, err := mgr.StartFlow(ctx, FlowConfig{
/root/astrate-mule/internal/flow/block.go:20:	Process(msg *FlowMessage) ([]*FlowMessage, error)
/root/astrate-mule/internal/flow/block.go:27:type SourceFunc func() ([]*FlowMessage, error)
/root/astrate-mule/internal/flow/block.go:31:type TransformFunc func(msg *FlowMessage) ([]*FlowMessage, error)
/root/astrate-mule/internal/flow/block.go:35:type SinkFunc func(msg *FlowMessage) error
/root/astrate-mule/internal/flow/block.go:43:func (s *sourceBlock) Process(_ *FlowMessage) ([]*FlowMessage, error) {
/root/astrate-mule/internal/flow/block.go:60:func (t *transformBlock) Process(msg *FlowMessage) ([]*FlowMessage, error) {
/root/astrate-mule/internal/flow/block.go:77:func (s *sinkBlock) Process(msg *FlowMessage) ([]*FlowMessage, error) {
/root/astrate-mule/internal/flow/graph.go:29:func (g *BlockGraph) Run(msg *FlowMessage) ([]*FlowMessage, error) {
/root/astrate-mule/internal/flow/graph.go:30:	cur := []*FlowMessage{msg}
/root/astrate-mule/internal/flow/graph.go:32:		var next []*FlowMessage
/root/astrate-mule/internal/flow/message_test.go:9:func TestFlowMessage_RoundTrip(t *testing.T) {
/root/astrate-mule/internal/flow/message_test.go:12:		msg     FlowMessage
/root/astrate-mule/internal/flow/message_test.go:17:			msg: FlowMessage{
/root/astrate-mule/internal/flow/message_test.go:27:			msg: FlowMessage{
/root/astrate-mule/internal/flow/message_test.go:36:			msg: FlowMessage{
/root/astrate-mule/internal/flow/message_test.go:45:			msg: FlowMessage{
/root/astrate-mule/internal/flow/message_test.go:54:			msg: FlowMessage{
/root/astrate-mule/internal/flow/message_test.go:64:			msg: FlowMessage{
/root/astrate-mule/internal/flow/message_test.go:73:			msg: FlowMessage{
/root/astrate-mule/internal/flow/message_test.go:93:			msg: FlowMessage{
/root/astrate-mule/internal/flow/message_test.go:112:			var got FlowMessage
/root/astrate-mule/internal/flow/message_test.go:199:func TestFlowMessage_WireFormatSchema(t *testing.T) {
/root/astrate-mule/internal/flow/message_test.go:200:	msg := FlowMessage{
/root/astrate-mule/internal/flow/message_test.go:225:func TestFlowMessage_UnmarshalRejectsUnknownSchema(t *testing.T) {
/root/astrate-mule/internal/flow/message_test.go:227:	var msg FlowMessage
/root/astrate-mule/internal/flow/message_test.go:233:func TestFlowMessage_MapSubtypesOmittedWhenEmpty(t *testing.T) {
/root/astrate-mule/internal/flow/message_test.go:234:	msg := FlowMessage{
/root/astrate-mule/internal/flow/message_test.go:256:func TestFlowMessage_ScalarSubtypeOmittedWhenEmpty(t *testing.T) {
/root/astrate-mule/internal/flow/message_test.go:257:	msg := FlowMessage{
/root/astrate-mule/internal/flow/message_test.go:276:func TestFlowMessage_BinaryBase64Encoding(t *testing.T) {
/root/astrate-mule/internal/flow/message_test.go:277:	msg := FlowMessage{
/root/astrate-mule/internal/flow/message_test.go:302:func TestFlowMessage_MapBinaryBase64Encoding(t *testing.T) {
/root/astrate-mule/internal/flow/message_test.go:303:	msg := FlowMessage{
/root/astrate-mule/internal/flow/flow_test.go:7:func TestFlowStatus_String(t *testing.T) {
/root/astrate-mule/internal/flow/flow_test.go:9:		status FlowStatus
/root/astrate-mule/internal/flow/flow_test.go:12:		{FlowStatusCreating, "creating"},
/root/astrate-mule/internal/flow/flow_test.go:13:		{FlowStatusRunning, "running"},
/root/astrate-mule/internal/flow/flow_test.go:14:		{FlowStatusStopped, "stopped"},
/root/astrate-mule/internal/flow/flow_test.go:15:		{FlowStatusFailed, "failed"},
/root/astrate-mule/internal/flow/flow_test.go:16:		{FlowStatus(99), "unknown(99)"},
/root/astrate-mule/internal/flow/flow_test.go:28:func TestFlowStatus_Transitions(t *testing.T) {
/root/astrate-mule/internal/flow/flow_test.go:32:		wantFinal  FlowStatus
/root/astrate-mule/internal/flow/flow_test.go:39:				f, err := mgr.StartFlow(t.Context(), FlowConfig{
/root/astrate-mule/internal/flow/flow_test.go:49:			wantFinal: FlowStatusRunning,
/root/astrate-mule/internal/flow/flow_test.go:56:				f, err := mgr.StartFlow(t.Context(), FlowConfig{
/root/astrate-mule/internal/flow/flow_test.go:68:			wantFinal: FlowStatusStopped,
/root/astrate-mule/internal/flow/flow_test.go:74:				f, err := mgr.StartFlow(t.Context(), FlowConfig{
/root/astrate-mule/internal/flow/flow_test.go:83:			wantFinal: FlowStatusFailed,
/root/astrate-mule/internal/flow/flow_test.go:103:	f, err := mgr.StartFlow(t.Context(), FlowConfig{
/root/astrate-mule/internal/flow/flow_test.go:126:	f, err := mgr.StartFlow(t.Context(), FlowConfig{
/root/astrate-mule/internal/flow/flow_test.go:134:	if got := f.Status(); got != FlowStatusRunning {
/root/astrate-mule/internal/flow/flow_test.go:142:	if got := f.Status(); got != FlowStatusStopped {
/root/astrate-mule/internal/flow/flow_test.go:152:	f, _ := mgr.StartFlow(t.Context(), FlowConfig{
/root/astrate-mule/internal/flow/flow_test.go:167:	f, err := mgr.StartFlow(t.Context(), FlowConfig{
/root/astrate-mule/internal/flow/flow_test.go:184:	f, err := mgr.StartFlow(t.Context(), FlowConfig{
/root/astrate-mule/internal/flow/router.go:46:// Router is the stream-based message router. It accepts FlowMessages, hashes
/root/astrate-mule/internal/flow/router.go:72:	msg    *FlowMessage
/root/astrate-mule/internal/flow/router.go:113:func (r *Router) Submit(msg *FlowMessage, qos byte) {

$ rg -rn 'FlowStatus|FlowConfig|FlowMessage' /root/astrate-mule/ --include '*.go' | grep -v 'internal/flow/'
rg: unrecognized flag --include

similar flags that are available: --include-zero

$ rg -n 'FlowStatus|FlowConfig|FlowMessage' /root/astrate-mule/ --glob '*.go' | grep -v 'internal/flow/'
(no output)

No usage outside `internal/flow/`. Three renames needed, including enum constants (`FlowStatusCreating` → `StatusCreating`, etc.) and `GetFlowStatus` → `GetStatus`. Let me do the renames via sed since there are many occurrences:
$ cd /root/astrate-mule/internal/flow && sed -i 's/FlowStatusCreating/StatusCreating/g; s/FlowStatusRunning/StatusRunning/g; s/FlowStatusStopped/StatusStopped/g; s/FlowStatusFailed/StatusFailed/g; s/FlowStatus/Status/g; s/FlowConfig/Config/g; s/FlowMessage/Message/g; s/GetFlowStatus/GetStatus/g' flow.go message.go block.go graph.go router.go flow_test.go message_test.go manager_test.go router_test.go
(no output)

$ rg -n 'FlowStatus|FlowConfig|FlowMessage|GetFlowStatus' /root/astrate-mule/internal/flow/
(no output)

All occurrences replaced. Now run the gates:
$ cd /root/astrate-mule && go vet ./...
(no output)

$ cd /root/astrate-mule && go test ./...
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
ok  	github.com/astrate-platform/astrate/internal/flow	0.246s
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

$ cd /root/astrate-mule && gofmt -l .
(no output)

Done: Renamed `FlowStatus`/`FlowConfig`/`FlowMessage` to `Status`/`Config`/`Message` (and enum constants `FlowStatusCreating` etc. to `StatusCreating`, method `GetFlowStatus` to `GetStatus`) in `internal/flow`.
Files: `internal/flow/flow.go`, `internal/flow/message.go`, `internal/flow/block.go`, `internal/flow/graph.go`, `internal/flow/router.go`, `internal/flow/flow_test.go`, `internal/flow/message_test.go`, `internal/flow/manager_test.go`, `internal/flow/router_test.go`
Verified: `go vet ./...` → pass, `go test ./...` → pass, `gofmt -l .` → clean
Unsure: nothing
Follow-ups: none
