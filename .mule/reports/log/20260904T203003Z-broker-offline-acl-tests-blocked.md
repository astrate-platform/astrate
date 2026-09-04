slug: broker-offline-acl-tests
verdict: blocked
at:  b8b48aa
ran: 2026-09-04T20:30:04Z on DietPi in 208s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/broker/aclhook.go
→ Read internal/broker/aclhook_test.go
Now let me understand the codebase — the Store interface, session registry, and existing test patterns.

$ rg -n "T0|T1|T2|tier" internal/broker/*_test.go | head -50
internal/broker/broker_test.go:106:// fakeStore is an in-memory broker.Store for T1 tests.
internal/broker/broker_test.go:183:// --- T1 tests ---

$ rg -n "type Store interface" internal/broker/ && rg -n "func ParseCN" internal/broker/ && rg -n "type realmPools" internal/broker/ && rg -n "hookDBTimeout" internal/broker/
internal/broker/authhook.go:24:type Store interface {
internal/broker/identity.go:52:func ParseCN(cn string) (Identity, error) {
internal/broker/authhook.go:60:type realmPools struct {
internal/broker/lifecycle.go:46:	ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
internal/broker/lifecycle.go:76:	ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
internal/broker/authhook.go:38:	// hookDBTimeout bounds the database work performed inside broker hooks
internal/broker/authhook.go:40:	hookDBTimeout = 5 * time.Second
internal/broker/authhook.go:364:	ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
internal/broker/aclhook.go:154:	ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)
internal/broker/aclhook.go:196:			ctx, cancel := context.WithTimeout(context.Background(), hookDBTimeout)

→ Read internal/broker/identity.go
→ Read internal/broker/authhook.go
→ Read internal/broker/broker_test.go
Now let me find where `offlineACL` and `aclHook` are constructed, and check the store.Device type and loadOwnership behavior.

$ rg -n "newOfflineACL|offlineACL|aclHook" internal/broker/*.go | rg -v "_test.go"
internal/broker/aclhook.go:20:	// offlineACLCacheTTL bounds how long a disconnected device's
internal/broker/aclhook.go:23:	offlineACLCacheTTL = 10 * time.Second
internal/broker/aclhook.go:95:// offlineACL answers read-side ACL checks for devices that hold a persistent
internal/broker/aclhook.go:100:type offlineACL struct {
internal/broker/aclhook.go:114:func newOfflineACL(st Store, pools *realmPools, log *slog.Logger) *offlineACL {
internal/broker/aclhook.go:115:	return &offlineACL{st: st, pools: pools, log: log, entries: map[string]*offlineEntry{}}
internal/broker/aclhook.go:121:func (o *offlineACL) ownershipOf(cn, iface string) (interfaceschema.Ownership, bool) {
internal/broker/aclhook.go:124:	if e != nil && time.Since(e.loadedAt) < offlineACLCacheTTL {
internal/broker/aclhook.go:145:func (o *offlineACL) load(cn string) map[string]interfaceschema.Ownership {
internal/broker/aclhook.go:164:// aclHook enforces the §3.2 matrix on every SUBSCRIBE filter, device
internal/broker/aclhook.go:167:type aclHook struct {
internal/broker/aclhook.go:171:	offline  *offlineACL
internal/broker/aclhook.go:176:func (h *aclHook) ID() string { return "astrate-acl" }
internal/broker/aclhook.go:179:func (h *aclHook) Provides(b byte) bool { return b == mqtt.OnACLCheck }
internal/broker/aclhook.go:182:func (h *aclHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
internal/broker/broker.go:148:		&aclHook{
internal/broker/broker.go:150:			offline: newOfflineACL(st, b.pools, log),

$ rg -n "type Device struct|type StoredInterface struct|type InterfaceVersion struct|type Realm struct" internal/store/*.go
internal/store/devices.go:32:type InterfaceVersion struct {
internal/store/devices.go:38:type Device struct {
internal/store/realms.go:16:type Realm struct {
internal/store/interfaces.go:18:type StoredInterface struct {

→ Read internal/store/devices.go [offset=28, limit=30]
→ Read internal/store/interfaces.go [offset=14, limit=12]
→ Read internal/broker/broker.go [offset=140, limit=20]

$ rg -n "func discardLogger|func newTestLogger|discardLogger" internal/broker/ | rg -v "_test.go:" ; rg -rn "func discardLogger" internal/
internal/housekeeping/http_test.go:n() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }
internal/engine/fakestore_test.go:n() *slog.Logger {
internal/broker/sessionstore_test.go:n() *slog.Logger {
internal/realm/http_test.go:n() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }
internal/appengine/http_test.go:n() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

Let me check available test helpers (`discardLogger`, `testutil`) and how sessions work with mochi clients for the T1 hook check.
$ rg -n "func ServerTLSCert|func DeviceTLSConfig|func MQTTConnect" internal/testutil/
internal/testutil/mqttclient.go:26:func ServerTLSCert(t testing.TB) (tls.Certificate, *x509.CertPool) {
internal/testutil/mqttclient.go:80:func DeviceTLSConfig(t testing.TB, clientCertPEM string, key *ecdsa.PrivateKey, roots *x509.CertPool) *tls.Config {
internal/testutil/mqttclient.go:142:func MQTTConnect(t testing.TB, brokerURL, clientID string, cleanSession bool, tlsCfg *tls.Config, tweaks ...func(*paho.ClientOptions)) (paho.Client, bool) {

$ rg -n "discardLogger" internal/broker/
internal/broker/broker_test.go:251:		Logger:           discardLogger(),
internal/broker/broker_test.go:385:		Logger:           discardLogger(),
internal/broker/broker_test.go:452:		Logger:           discardLogger(),
internal/broker/broker_test.go:509:		Logger:           discardLogger(),
internal/broker/sessionstore_test.go:15:func discardLogger() *slog.Logger {
internal/broker/sessionstore_test.go:23:	ss := newSessionStore(filepath.Join(t.TempDir(), "sessions.db"), discardLogger())
internal/broker/sessionstore_test.go:28:	srv := mqtt.New(&mqtt.Options{InlineClient: false, Logger: discardLogger()})
internal/broker/sessionstore_test.go:46:	ss := newSessionStore(filepath.Join(t.TempDir(), "s.db"), discardLogger())
internal/broker/broker_e2e_test.go:205:		Logger:            discardLogger(),

Let me verify the mochi-mqtt `Client` struct fields for the T1 hook test.
$ rg -n "type Client struct" $(go env GOMODCACHE)/github.com/mochi-mqtt/server*/v2/*.go && rg -n "type ClientConnectionDetails struct|Inline" $(go env GOMODCACHE)/github.com/mochi-mqtt/server*/v2/*.go | head -5
rg: /root/go/pkg/mod/github.com/mochi-mqtt/server*/v2/*.go: IO error for operation on /root/go/pkg/mod/github.com/mochi-mqtt/server*/v2/*.go: No such file or directory (os error 2)

$ go env GOMODCACHE && rg -n "mochi-mqtt" go.mod && ls $(go env GOMODCACHE)/github.com/mochi-mqtt/ 2>/dev/null
/root/go/pkg/mod
15:	github.com/mochi-mqtt/server/v2 v2.7.9 // embedded MQTT broker (§1.1)
server

$ rg -n "type Client struct" /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/*.go
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/clients.go:105:type Client struct {

! permission requested: external_directory (/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/*); auto-rejecting
✗ Read /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/clients.go failed [limit=40, offset=105]
Error: The user rejected permission to use this specific tool call.
