slug: recipe-code-review
verdict: proposed
at:  104d379
ran: 2026-09-04T18:43:36Z on DietPi in 443s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
I'll carry out the code review. First, let me check which areas were recently reviewed and understand the codebase map.

$ ls internal/ pkg/ 2>/dev/null; echo "---REVIEWS---"; ls -la .mule/reviews/ 2>/dev/null; echo "---LOG---"; tail -30 .mule/log.md 2>/dev/null
internal/:
appengine
auth
broker
config
engine
flow
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil

pkg/:
astarteapi
deviceid
interfaceschema
payload
---REVIEWS---
total 48
drwxr-xr-x 2 root root 4096 Sep  3 21:58 .
drwxr-xr-x 8 root root 4096 Sep  4 18:55 ..
-rw-r--r-- 1 root root 7300 Sep  3 13:26 engine-2026-09-03.md
-rw-r--r-- 1 root root 5793 Sep  2 19:01 flow-2026-09-02.md
-rw-r--r-- 1 root root 5106 Sep  3 21:58 pairing-2026-09-03.md
-rw-r--r-- 1 root root 4291 Sep  1 20:46 realm-2026-09-01.md
-rw-r--r-- 1 root root 4406 Jul 27 14:57 store-2026-07-27.md
---LOG---
| 2026-08-31 | compat-note-v140-rc3 | blocked | 222s | gates failed |
| 2026-08-31 | probe-emptycache-resend-device-error | blocked | 74s | gates failed |
| 2026-08-31 | compat-note-v133 | blocked | 199s | gates failed |
| 2026-08-31 | probe-mqtt-capabilities-declaration | blocked | 459s | gates failed |
| 2026-08-31 | probe-binaryblob-validation | blocked | 185s | gates failed |
| 2026-08-31 | probe-properties-on-connect-encoding | blocked | 315s | gates failed |
| 2026-08-31 | docs-sync-pairing-status-enum | blocked | 59s | gates failed |
| 2026-09-01 | issue-91-pairing-health-serve-root | blocked | 930s | gates failed |
| 2026-09-01 | issue-68-async-operation-accepted | blocked | 520s | gates failed |
| 2026-09-01 | probe-object-validation-selected-interface | blocked | 201s | gates failed |
| 2026-09-01 | compat-note-v140-rc5 | blocked | 282s | gates failed |
| 2026-09-01 | realm-policy-list-sorted | blocked | 182s | gates failed |
| 2026-09-01 | realm-interface-lookup-404 | blocked | 83s | gates failed |
| 2026-09-02 | flow-setstatus-race | blocked | 743s | gates failed |
| 2026-09-02 | flow-validate-deadcode | blocked | 223s | gates failed |
| 2026-09-02 | flow-unmarshal-error-tests | blocked | 206s | gates failed |
| 2026-09-02 | flow-datawirescalar-fallthrough | blocked | 190s | gates failed |
| 2026-09-02 | docs-sync-appengine-422 | blocked | 244s | gates failed |
| 2026-09-02 | docs-sync-appengine-socket | blocked | 139s | gates failed |
| 2026-09-02 | lint-engine-goimports | blocked | 421s | gates failed |
| 2026-09-02 | lint-flow-findings | blocked | 260s | gates failed |
| 2026-09-02 | lint-pairing-test-param | blocked | 186s | gates failed |
| 2026-09-02 | lint-flow-stutter | blocked | 98s | gates failed |
| 2026-09-03 | serverdata-aftercommit-test | blocked | 114s | gates failed |
| 2026-09-03 | consumer-properties-plaintext-fallback | blocked | 1182s | gates failed |
| 2026-09-03 | object-path-ok-comment | blocked | 113s | gates failed |
| 2026-09-03 | housekeeping-unit-tests | blocked | 182s | gates failed |
| 2026-09-03 | store-downsample-lttb | blocked | 614s | gates failed |
| 2026-09-04 | pairing-info-verify-ratelimit | blocked | 186s | gates failed |
| 2026-09-04 | docs-sync-realm-policies-422 | blocked | 133s | gates failed |

→ Read .mule/todo.md
Recently reviewed: engine (09-03), flow (09-02), pairing (09-03), realm (09-01), store (07-27). Let me pick an unreviewed package and explore the codebase map to decide.

Let me check the sizes of the unreviewed packages and weigh which to review.

$ for p in appengine auth broker config housekeeping httpx observability payload interfaceschema deviceid astarteapi; do
  d="internal/$p"
  [ -d "$d" ] || d="pkg/$p"
  [ -d "$d" ] && echo "$d: $(rg -c '' --glob '*.go' "$d" | awk -F: '{s+=$2} END{print s}') lines"
done
internal/appengine: 4885 lines
internal/auth: 1404 lines
internal/broker: 3441 lines
internal/config: 647 lines
internal/housekeeping: 500 lines
internal/httpx: 357 lines
internal/observability: 272 lines
pkg/payload: 2578 lines
pkg/interfaceschema: 2345 lines
pkg/deviceid: 358 lines
pkg/astarteapi: 441 lines

None of appengine, broker, housekeeping, httpx, observability, auth, or the pkg packages have reviews yet. Let me look at the review files to see what was covered, and check what a recently-reviewed package like engine referenced (to avoid overlap). The `code-review` recipe says pick one area not recently done. `broker` is a substantial unreviewed package (3441 lines) and is a concurrency-heavy area (goroutines for message routing) — good candidate. Let me verify broker isn't covered by an existing review.
$ head -40 .mule/reviews/engine-2026-09-03.md; echo "=====PAIRING====="; head -40 .mule/reviews/pairing-2026-09-03.md
# Code review: internal/engine — 2026-09-03

## Files read

- `internal/engine/cache.go` (496 lines) — schema/device caches, copy-on-write snapshot
- `internal/engine/router.go` (490 lines) — sharded pipeline, Submit, metrics
- `internal/engine/batch.go` (406 lines) — micro-batcher, flush, row conversion
- `internal/engine/data.go` (386 lines) — §2.6 validation pipeline, handle, reject
- `internal/engine/engine.go` (318 lines) — New, Drain, fireCommitted, lifecycle
- `internal/engine/control.go` (315 lines) — emptyCache, producer/properties, consumer/properties
- `internal/engine/serverdata.go` (280 lines) — PublishServerValue, UnsetServerProperty
- `internal/engine/introspection.go` (168 lines) — handleIntrospection, parseIntrospection
- `internal/engine/capabilities.go` (101 lines) — BSON capabilities parsing
- `internal/engine/topics.go` (97 lines) — splitDeviceTopic, classify, matchInterface
- `internal/engine/forward/http.go` (125 lines) — HTTP trigger forwarder

Plus test files: `data_test.go`, `serverdata_test.go`, `control_test.go`, `cache_test.go`, `router_test.go`, `fakestore_test.go`, `introspection_test.go`, `capabilities_test.go`, `bus_events_test.go`, `batch_test.go`.

## What I read in detail

I focused on the M6b seam wiring — how `PublishServerValue` and the control-channel handlers interact with the batch/trigger/bus subsystem — because the code's comments consistently reference the `afterCommit` hook (engine.go:109, batch.go:168) as the single observation point for committed ops. I verified the full call graph from `PublishServerValue` through persistence to broker publish, and from `finalize` through `afterCommit` to `fireCommitted`/`fireData`.

## Findings

### 1. `PublishServerValue` bypasses `afterCommit` — no trigger or live bus events

**Where:** `internal/engine/serverdata.go:99-162`

The batch path persists ops, then `finalize` (batch.go:168) calls `afterCommit`, which is wired to `fireCommitted` (engine.go:139) and fires data triggers + live bus events.

`PublishServerValue` bypasses the batch entirely: it persists directly (`UpsertProperty` / `AppendDatastreams`), then publishes to the broker. It never calls `afterCommit`. Consequences:

- Server-initiated property/datastream writes produce no `incoming_data` live bus events — a Channels viewer watching the stream will not see server-side writes.
- No data triggers fire for server-owned interface paths (the `triggers.DataEvent` matching loop in `fireData` is never reached).

This may be intentional — server writes are admin-initiated, not device-initiated — but it is undocumented and untested. The batch path's `TestBatchTriggers` proves triggers fire for device-originated ops; no equivalent test exists for server-originated ops. At minimum the live-bus gap should be documented or deliberately excluded.

**Proposed task:** in `internal/engine/serverdata_test.go`, add a test that installs a data trigger matching a server-owned interface, calls `PublishServerValue`, and asserts the trigger was NOT fired (proving the documented exclusion), or DOES fire (if the intended behaviour is parity with the data path). The test should also verify whether a Channels live-bus subscriber sees the event.

### 2. `sendConsumerProperties` zlib fallback when device not cached
=====PAIRING=====
# Code review: internal/pairing — 2026-09-03

## Files read

- `internal/pairing/service.go` (427 lines) — flows A–C: Register, Unregister,
  Credentials, Info, VerifyCredentials, authenticateDevice
- `internal/pairing/http.go` (348 lines) — /pairing/v1 routes, rate limiting,
  writeServiceError
- `internal/pairing/ratelimit.go` (79 lines) — keyed token-bucket limiter
- `internal/pairing/ca/ca.go` (287 lines) — embedded per-realm CA, SignCSR, Verify
- Tests: `service_test.go`, `http_test.go`, `ratelimit_test.go`, `ca/ca_test.go`
- Docs: `docs/DESIGN.md` §4.4 (flows A–C), §4.5 (hardening checklist)

This is the first review of `pairing`; the rotation had covered store, realm, flow,
engine.

## What I read in detail

I focused on the rate-limiting seam (§4.5 says "Rate limits: pairing endpoints
(per-IP and per-device)") because the doting design doc and the code diverge there in a
way that matters for CPU burn. I walked every handler in `http.go` and checked which ones
apply a `Limiter` and which ones reach `bcrypt.CompareHashAndPassword` via
`authenticateDevice` (service.go:367-398).

## Findings

### 1. Info + Verify do full bcrypt work with NO rate limit — unauthenticated CPU burn

**Where:** `internal/pairing/http.go` — `handleInfo` (line 197), `handleVerify`
(line 238); bcrypt compare via `authenticateDevice` (service.go:394).

`handleRegister` (line 106) and `handleCredentials` (line 154) both gate on a token-bucket
`Limiter` before doing any work. But `handleInfo` and `handleVerify` — the two device-auth
endpoints — check only `bearerSecret(r)` (any non-empty header passes) and then call
`Info`/`VerifyCredentials`, each of which runs `authenticateDevice` → a full
`bcrypt.CompareHashAndPassword` (cost 10, ~50-100ms of CPU).

Consequences:

- The `verify` endpoint is the SDK boot hot path that upstream expects to run on every

Broker is a good unreviewed target (3441 lines, concurrency-heavy MQTT message routing — exactly the kind of area where a behavior/wrong/unguarded finding lives). Let me look at the broker package structure.
$ for f in internal/broker/*.go; do echo "$f: $(wc -l < $f)"; done
internal/broker/aclhook.go: 209
internal/broker/aclhook_test.go: 145
internal/broker/authhook.go: 440
internal/broker/broker.go: 343
internal/broker/broker_e2e_test.go: 618
internal/broker/broker_test.go: 475
internal/broker/identity.go: 104
internal/broker/identity_test.go: 103
internal/broker/intake.go: 92
internal/broker/lifecycle.go: 90
internal/broker/publisher.go: 62
internal/broker/sessionstore.go: 533
internal/broker/sessionstore_test.go: 227

Let me read the core broker files. I'll read the largest non-test files first. Let me read broker.go, authhook.go, sessionstore.go, and the smaller ones.
→ Read internal/broker/broker.go
→ Read internal/broker/intake.go
→ Read internal/broker/publisher.go
→ Read internal/broker/lifecycle.go
Now let me read the sessionstore.go (largest) and the identity.go and aclhook.go files.
→ Read internal/broker/sessionstore.go
→ Read internal/broker/aclhook.go
→ Read internal/broker/identity.go
Now let me read authhook.go (440 lines), which is the security-critical piece, and check the sessions registry (`sessionRegistry`) used throughout — it's referenced but I haven't seen its definition.

$ rg -l "type sessionRegistry|type realmPools|func \(.*sessionRegistry|type deviceSession" internal/broker/
internal/broker/authhook.go

→ Read internal/broker/authhook.go
The `sessionRegistry` and `realmPools` are defined in authhook.go. I've now read the whole broker package's non-test source. Let me check the tests to understand what's covered, and look at the engine's parser/capabilities handling for potential gaps. Let me review the test files briefly to see what's asserted.

Let me look at the broker tests, particularly around the intake/Ack path and offline ACL, which are the concurrency-sensitive seams.
→ Read internal/broker/broker_test.go
I now have a good understanding of the broker package. Let me read the e2e test and the sessionstore and aclhook tests to see what's covered, focusing on the offline ACL and session takeover paths which are the concurrency-sensitive seams.
→ Read internal/broker/broker_e2e_test.go
→ Read internal/broker/aclhook_test.go
Let me read the aclhook_test and sessionstore_test to see test coverage of the ACL and session restore paths.
→ Read internal/broker/sessionstore_test.go
I've read the whole broker package. Let me verify a couple of mochi-internal behaviors the findings depend on — the ACL-deny semantics for device publishes and what `DisconnectClient` does to a persisted session — before proposing anything.

$ find / -type d -path "*mochi-mqtt/server*" 2>/dev/null | head; echo "---"; rg -l "mochi-mqtt/server" go.mod 2>/dev/null; grep mochi go.mod
/root/go/pkg/mod/github.com/mochi-mqtt/server
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/packets
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks/storage
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks/storage/bolt
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks/storage/badger
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks/storage/redis
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks/storage/pebble
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks/debug
---
go.mod
	github.com/mochi-mqtt/server/v2 v2.7.9 // embedded MQTT broker (§1.1)

$ rg -n "OnACLCheck|ACL" /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go | head -20; echo "=== disconnects ==="; rg -n "func (s \*Server) DisconnectClient|DisconnectClient\(" /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go | head
240:// messages from the embedding application, set the inline flag to true to bypass ACL and
751:// to any topic (including $SYS) and bypass ACL checks. The qos byte is used for limiting the
840:// InlineClients using this method can publish packets to any topic (including $SYS) and bypass ACL checks.
867:	if !cl.Net.Inline && !s.hooks.OnACLCheck(cl, pk.TopicName, true) {
1030:	if !s.hooks.OnACLCheck(cl, pk.TopicName, false) {
1258:		} else if !s.hooks.OnACLCheck(cl, sub.Filter, false) {
=== disconnects ===
526:			_ = s.DisconnectClient(cl, code)
564:		_ = s.DisconnectClient(existing, packets.ErrSessionTakenOver)                                   // [MQTT-3.1.4-3]
864:		return s.DisconnectClient(cl, packets.ErrReceiveMaximum) // ~[MQTT-3.3.4-7] ~[MQTT-3.3.4-8]
873:			return s.DisconnectClient(cl, packets.ErrNotAuthorized)
1414:func (s *Server) DisconnectClient(cl *Client, code packets.Code) error {
1511:		_ = s.DisconnectClient(cl, packets.ErrServerShuttingDown)

$ sed -n '855,880p' /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go
}

// processPublish processes a Publish packet.
func (s *Server) processPublish(cl *Client, pk packets.Packet) error {
	if !cl.Net.Inline && !IsValidFilter(pk.TopicName, true) {
		return nil
	}

	if atomic.LoadInt32(&cl.State.Inflight.receiveQuota) == 0 {
		return s.DisconnectClient(cl, packets.ErrReceiveMaximum) // ~[MQTT-3.3.4-7] ~[MQTT-3.3.4-8]
	}

	if !cl.Net.Inline && !s.hooks.OnACLCheck(cl, pk.TopicName, true) {
		if pk.FixedHeader.Qos == 0 {
			return nil
		}

		if cl.Properties.ProtocolVersion != 5 {
			return s.DisconnectClient(cl, packets.ErrNotAuthorized)
		}

		ackType := packets.Puback
		if pk.FixedHeader.Qos == 2 {
			ackType = packets.Pubrec
		}

$ sed -n '1414,1470p' /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go
func (s *Server) DisconnectClient(cl *Client, code packets.Code) error {
	out := packets.Packet{
		FixedHeader: packets.FixedHeader{
			Type: packets.Disconnect,
		},
		ReasonCode: code.Code,
		Properties: packets.Properties{},
	}

	if code.Code >= packets.ErrUnspecifiedError.Code {
		out.Properties.ReasonString = code.Reason //  // [MQTT-3.14.2-1]
	}

	// We already have a code we are using to disconnect the client, so we are not
	// interested if the write packet fails due to a closed connection (as we are closing it).
	err := cl.WritePacket(out)
	if !s.Options.Capabilities.Compatibilities.PassiveClientDisconnect {
		cl.Stop(code)
		if code.Code >= packets.ErrUnspecifiedError.Code {
			return code
		}
	}

	return err
}

// publishSysTopics publishes the current values to the server $SYS topics.
// Due to the int to string conversions this method is not as cheap as
// some of the others so the publishing interval should be set appropriately.
func (s *Server) publishSysTopics() {
	pk := packets.Packet{
		FixedHeader: packets.FixedHeader{
			Type:   packets.Publish,
			Retain: true,
		},
		Created: time.Now().Unix(),
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	atomic.StoreInt64(&s.Info.MemoryAlloc, int64(m.HeapInuse))
	atomic.StoreInt64(&s.Info.Threads, int64(runtime.NumGoroutine()))
	atomic.StoreInt64(&s.Info.Time, time.Now().Unix())
	atomic.StoreInt64(&s.Info.Uptime, time.Now().Unix()-atomic.LoadInt64(&s.Info.Started))
	atomic.StoreInt64(&s.Info.ClientsTotal, int64(s.Clients.Len()))
	atomic.StoreInt64(&s.Info.ClientsDisconnected, atomic.LoadInt64(&s.Info.ClientsTotal)-atomic.LoadInt64(&s.Info.ClientsConnected))

	info := s.Info.Clone()
	topics := map[string]string{
		SysPrefix + "/broker/version":              s.Info.Version,
		SysPrefix + "/broker/time":                 Int64toa(info.Time),
		SysPrefix + "/broker/uptime":               Int64toa(info.Uptime),
		SysPrefix + "/broker/started":              Int64toa(info.Started),
		SysPrefix + "/broker/load/bytes/received":  Int64toa(info.BytesReceived),
		SysPrefix + "/broker/load/bytes/sent":      Int64toa(info.BytesSent),
		SysPrefix + "/broker/clients/connected":    Int64toa(info.ClientsConnected),
		SysPrefix + "/broker/clients/disconnected": Int64toa(info.ClientsDisconnected),

The ACL-deny path for a denied QoS0 publish silently drops the packet (`return nil`); for QoS≥1 it disconnects the client. Now let me verify what `cl.Stop(code)` does to the persisted session (relevant to `DisconnectDevice` on device deletion), and how `OnDisconnect`'s `expire` flag is set.
$ rg -n "OnDisconnect|expire bool|func (cl \*Client) Stop" /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/clients.go | head -20
(no output)

$ rg -n "OnDisconnect|func.*Stop|expire" /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/*.go | grep -v _test | head -30
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:152:	clientExpiry   *time.Ticker     // interval ticker for cleaning expired clients
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:153:	inflightExpiry *time.Ticker     // interval ticker for cleaning up expired inflight messages
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:485:	expire := (cl.Properties.ProtocolVersion == 5 && cl.Properties.Props.SessionExpiryInterval == 0) || (cl.Properties.ProtocolVersion < 5 && cl.Properties.Clean)
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:486:	s.hooks.OnDisconnect(cl, err, expire)
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:488:	if expire && !cl.IsTakenOver() {
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:1668:		expire := (cl.Properties.ProtocolVersion == 5 && cl.Properties.Props.SessionExpiryInterval == 0) || (cl.Properties.ProtocolVersion < 5 && cl.Properties.Clean)
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:1669:		s.hooks.OnDisconnect(cl, packets.ErrServerShuttingDown, expire)
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:1670:		if expire {
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:1704:		expire := s.Options.Capabilities.MaximumSessionExpiryInterval
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:1706:			expire = client.Properties.Props.SessionExpiryInterval
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:1709:		if disconnected+int64(expire) < dt {
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:1716:// clearExpiredRetainedMessage deletes retained messages from topics if they have expired.
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:1719:		expired := pk.ProtocolVersion == 5 && pk.Expiry > 0 && pk.Expiry < now // [MQTT-3.3.2-5]
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:1726:		if expired || enforced {
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:1733:// clearExpiredInflights deletes any inflight messages which have expired.
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/clients.go:341:// ClearExpiredInflights deletes any inflight messages which have expired.
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/clients.go:345:		expired := tk.ProtocolVersion == 5 && tk.Expiry > 0 && tk.Expiry < now // [MQTT-3.3.2-5]
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/clients.go:351:		if expired || enforced {
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/clients.go:393:func (cl *Client) Stop(err error) {
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/clients.go:413:func (cl *Client) StopCause() error {
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/clients.go:421:func (cl *Client) StopTime() int64 {
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks.go:29:	OnDisconnect
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks.go:88:	OnDisconnect(cl *Client, err error, expire bool)
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks.go:185:func (h *Hooks) Stop() {
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks.go:219:func (h *Hooks) OnStopped() {
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks.go:259:// OnDisconnect is called when a client is disconnected for any reason.
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks.go:260:func (h *Hooks) OnDisconnect(cl *Client, err error, expire bool) {
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks.go:262:		if hook.Provides(OnDisconnect) {
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks.go:263:			hook.OnDisconnect(cl, err, expire)
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/hooks.go:483:// OnQosDropped is called the Qos flow for a message expires. In other words, when

$ sed -n '470,500p' /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go; echo "=== clients.go Stop ==="; sed -n '393,412p' /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/clients.go
return fmt.Errorf("resend inflight: %w", err)
		}
	}

	s.hooks.OnSessionEstablished(cl, pk)

	err = cl.Read(s.receivePacket)
	if err != nil {
		s.sendLWT(cl)
		cl.Stop(err)
	} else {
		cl.Properties.Will = Will{} // [MQTT-3.14.4-3] [MQTT-3.1.2-10]
	}
	s.Log.Debug("client disconnected", "error", err, "client", cl.ID, "remote", cl.Net.Remote, "listener", listener)

	expire := (cl.Properties.ProtocolVersion == 5 && cl.Properties.Props.SessionExpiryInterval == 0) || (cl.Properties.ProtocolVersion < 5 && cl.Properties.Clean)
	s.hooks.OnDisconnect(cl, err, expire)

	if expire && !cl.IsTakenOver() {
		cl.ClearInflights()
		s.UnsubscribeClient(cl)
		s.Clients.Delete(cl.ID) // [MQTT-4.1.0-2] ![MQTT-3.1.2-23]
	}

	return err
}

// readConnectionPacket reads the first incoming header for a connection, and if
// acceptable, returns the valid connection packet.
func (s *Server) readConnectionPacket(cl *Client) (pk packets.Packet, err error) {
	fh := new(packets.FixedHeader)
=== clients.go Stop ===
func (cl *Client) Stop(err error) {
	cl.State.endOnce.Do(func() {

		if cl.Net.Conn != nil {
			_ = cl.Net.Conn.Close() // omit close error
		}

		if err != nil {
			cl.State.stopCause.Store(err)
		}

		if cl.State.cancelOpen != nil {
			cl.State.cancelOpen()
		}

		atomic.StoreInt64(&cl.State.disconnected, time.Now().Unix())
	})
}

// StopCause returns the reason the client connection was stopped, if any.

The `DisconnectDevice` finding is confirmed: `expire` is derived from the client's Clean flag/session expiry, not the reason code — so a persistent-session device deleted via `DisconnectDevice` keeps its bbolt session. Let me quickly verify how mochi derives the per-message `Expiry` (relevant to `Publisher.Publish`'s `MessageExpiryInterval`), then write the review.
$ rg -n "Expiry" /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/topics.go /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/packets/packets.go 2>/dev/null | head -20
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/packets/packets.go:133:	Expiry          int64         // unix timestamp indicating when the packet will expire and should be deleted
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/packets/packets.go:215:		Expiry:         pk.Expiry,

$ rg -n "\.Expiry =|Expiry:" /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/*.go | grep -v _test | head -20
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:183:			clientExpiry:   time.NewTicker(time.Second),
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:184:			inflightExpiry: time.NewTicker(time.Second),
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:185:			retainedExpiry: time.NewTicker(time.Second),
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:890:		pk.Expiry = pk.Created + expiry
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:994:	if pk.Expiry == 0 {
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:997:			pk.Expiry = pk.Created + expiry
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:1090:			out.Expiry = -1
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:1154:		Expiry:     time.Now().Unix() + s.Options.Capabilities.MaximumMessageExpiryInterval,
/root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go:1540:		pk.Expiry = time.Now().Unix() + int64(pk.Connect.WillProperties.WillDelayInterval)

$ sed -n '880,895p' /root/go/pkg/mod/github.com/mochi-mqtt/server/v2@v2.7.9/server.go
ack := s.buildAck(pk.PacketID, ackType, 0, pk.Properties, packets.ErrNotAuthorized)
		return cl.WritePacket(ack)
	}

	pk.Origin = cl.ID
	pk.Created = time.Now().Unix()

	if expiry := minimum(s.Options.Capabilities.MaximumMessageExpiryInterval,
		int64(pk.Properties.MessageExpiryInterval)); expiry > 0 {
		pk.Expiry = pk.Created + expiry
	}

	if !cl.Net.Inline {
		if pki, ok := cl.State.Inflight.Get(pk.PacketID); ok {
			if pki.FixedHeader.Type == packets.Pubrec { // [MQTT-4.3.3-10]

Now I have everything verified. Let me write the review file and append the task lines. First the review:
← Write .mule/reviews/broker-2026-09-04.md
Wrote file successfully.

Now let me append the task lines to `.mule/todo.md`:
→ Read .mule/todo.md [offset=125, limit=5]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -126,4 +126,8 @@
 - [!] housekeeping-unit-tests: in `internal/housekeeping`, add a unit test file (no integration build tag) that exercises `CreateRealm` validation paths — blank name returns ErrValidation, blank JWT key returns ErrValidation, negative registration limit returns ErrValidation — using a mock store and sealer. The package currently has only integration tests gated behind `//go:build integration`, so these paths are never tested on the Pi. [auto] — BLOCKED: gates failed
 - [!] store-downsample-lttb: in `internal/store`, implement the TODO at store.go:135 — when `s.hasToolkit` is true, switch `Downsample` from `time_bucket+avg` to toolkit `lttb()` downsampling. The probe already records availability; the time_bucket path in datastreams.go is the always-working default. Add a test that exercises both code paths (mocking `hasToolkit`). [auto] — BLOCKED: gates failed
 - [!] pairing-info-verify-ratelimit: in `internal/pairing/http.go`, add per-IP (and per-device) token-bucket rate limits to `handleInfo` (line 197) and `handleVerify` (line 238), which currently run a full bcrypt `CompareHashAndPassword` (via `authenticateDevice`, service.go:394) on any bearer string with no throttle — the sibling `handleCredentials` that does identical bcrypt work is limited at burst-5, and DESIGN §4.5 says "pairing endpoints" should be rate-limited. Reuse the credentials limiter defaults (1 req/s, burst 5) and extend `TestPairingHTTPRateLimit` to assert a burst of fake-secret requests to `GET /devices/{id}` and `POST .../verify` returns 429. [auto] — BLOCKED: gates failed
 - [!] docs-sync-realm-policies-422: in docs/api/astarte_realm_management_api.yaml, add the undocumented 422 Unprocessable Entity response to DELETE /policies/{name} — the code (internal/realm/service.go:398) returns ErrValidation when a policy is still referenced by a trigger, which maps to 422 via writeError; the YAML only lists 204, 401, 404, 500. [auto] — BLOCKED: gates failed
+- [ ] broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [auto]
+- [ ] broker-offline-acl-tests: in `internal/broker/aclhook_test.go`, unit-test the offline-delivery ACL — `offlineACL.ownershipOf` cache hit within TTL, TTL expiry triggering a reload, and a load failure caching as empty and denying (aclhook.go:116-138) — plus a T1 `OnACLCheck` check that the offline branch (delivery to a session-present but disconnected device, aclhook.go:196-202) consults the store. The rule currently has zero direct tests; `TestCheckACLMatrix` exercises only the pure `checkACL`. [auto]
+- [ ] broker-disconnect-device-zombie-session: in `internal/broker/broker.go` `DisconnectDevice` (lines 260-266), force-closing the live connection leaves the persisted bbolt session (subscriptions + offline queue, keyed by CN) in place — mochi derives `expire` from the Clean flag, not the reason code (server.go:485-491), so `dropClient` (sessionstore.go:500) is never reached and a later device with the same CN resurrects the deleted device's session. Wire the session store handle into `Broker` and call `dropClient` from `DisconnectDevice`, and add a T1 test in `broker_test.go`: connect clean=false + subscribe, call `DisconnectDevice`, reconnect — assert `session_present` is false and the stale subscription is gone. [auto]
+- [ ] broker-onconnect-doc-comment: in `internal/broker/authhook.go:314`, restore the missing first line of the `OnConnect` doc comment — it opens mid-sentence ("of the Astarte MQTT v1 protocol, and mochi publishes them without a publish-side ACL check…") and the Will-clearing security rationale (a retained LWT would escape the §3.2 matrix) is only half-documented. [auto]

Let me verify the appended lines:
$ tail -5 .mule/todo.md
- [!] docs-sync-realm-policies-422: in docs/api/astarte_realm_management_api.yaml, add the undocumented 422 Unprocessable Entity response to DELETE /policies/{name} — the code (internal/realm/service.go:398) returns ErrValidation when a policy is still referenced by a trigger, which maps to 422 via writeError; the YAML only lists 204, 401, 404, 500. [auto] — BLOCKED: gates failed
- [ ] broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [auto]
- [ ] broker-offline-acl-tests: in `internal/broker/aclhook_test.go`, unit-test the offline-delivery ACL — `offlineACL.ownershipOf` cache hit within TTL, TTL expiry triggering a reload, and a load failure caching as empty and denying (aclhook.go:116-138) — plus a T1 `OnACLCheck` check that the offline branch (delivery to a session-present but disconnected device, aclhook.go:196-202) consults the store. The rule currently has zero direct tests; `TestCheckACLMatrix` exercises only the pure `checkACL`. [auto]
- [ ] broker-disconnect-device-zombie-session: in `internal/broker/broker.go` `DisconnectDevice` (lines 260-266), force-closing the live connection leaves the persisted bbolt session (subscriptions + offline queue, keyed by CN) in place — mochi derives `expire` from the Clean flag, not the reason code (server.go:485-491), so `dropClient` (sessionstore.go:500) is never reached and a later device with the same CN resurrects the deleted device's session. Wire the session store handle into `Broker` and call `dropClient` from `DisconnectDevice`, and add a T1 test in `broker_test.go`: connect clean=false + subscribe, call `DisconnectDevice`, reconnect — assert `session_present` is false and the stale subscription is gone. [auto]
- [ ] broker-onconnect-doc-comment: in `internal/broker/authhook.go:314`, restore the missing first line of the `OnConnect` doc comment — it opens mid-sentence ("of the Astarte MQTT v1 protocol, and mochi publishes them without a publish-side ACL check…") and the Will-clearing security rationale (a retained LWT would escape the §3.2 matrix) is only half-documented. [auto]

Done: reviewed `internal/broker` — first time this package (recent reviews: store, realm, flow, engine, pairing) — and proposed four tasks.
Files: `.mule/reviews/broker-2026-09-04.md` (written), `.mule/todo.md` (4 lines appended)
Verified: read all 8 broker source files + 4 test files; checked mochi v2.7.9 source (module cache) for the ACL-deny semantics, `DisconnectClient`/`expire` derivation, and inline `pk.Expiry` — the findings rest on those reads, not pattern-matching. No code touched, no git operations.
Unsure: nothing major — finding #2's fix direction (drop the persistent bbolt session on device deletion) is a defensible intended behavior stated explicitly in the task line; if upstream intentionally keeps sessions, the test documents that instead.
Follow-ups (appended to `.mule/todo.md`):
- `broker-acl-coldstart-introspection-miss` — QoS0 publish to a freshly-introspected interface silently dropped in the first second after connect (miss-path reload debounced against the admit stamp).
- `broker-offline-acl-tests` — offline-delivery ACL has zero direct tests.
- `broker-disconnect-device-zombie-session` — deleted device leaves a persistent bbolt session that a same-CN re-registration resurrects.
- `broker-onconnect-doc-comment` — truncated security-rationale comment at authhook.go:314.
