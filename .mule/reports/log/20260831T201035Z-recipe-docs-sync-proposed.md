slug: recipe-docs-sync
verdict: proposed
at:  19ad0fe
ran: 2026-08-31T20:10:35Z on DietPi in 203s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
→ Read .mule/todo.md
✱ Glob "docs/api/*.yaml" 5 matches
→ Read .mule/for-giulio.md
Let me pick the **pairing** surface and compare documented vs actual routes.

$ rg -n '^\s{2}/' /root/astrate-mule/docs/api/astarte_pairing_api.yaml
27:  /pairing/v1/{realm}/agent/devices:
77:  /pairing/v1/{realm}/agent/devices/{deviceID}:
100:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials:
154:  /pairing/v1/{realm}/devices/{deviceID}:
192:  /pairing/v1/{realm}/devices/{deviceID}/protocols/astarte_mqtt_v1/credentials/verify:

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' /root/astrate-mule/internal/ --glob '*pairing*'
(no output)

$ rg -n 'HandleFunc|\.Methods\(|\.Handle\(' /root/astrate-mule/internal/ --glob '*pairing*'
(no output)

$ rg -rn 'pairing' /root/astrate-mule/internal/ --glob '*route*' --glob '*router*'
(no output)

$ ls /root/astrate-mule/internal/pairing/
ca
http.go
http_test.go
ratelimit.go
ratelimit_test.go
service.go
service_test.go
testdata

✱ Grep "pairing" in internal · 0 matches
→ Read internal/pairing/http.go
Now let me read the documented status codes and response fields from the YAML:
→ Read docs/api/astarte_pairing_api.yaml

$ rg -n 'Err(AfterUnregister|AlreadyRegistered|RegistrationLimit|Unauthorized|Inhibited|InvalidHWID|InvalidCSR|InvalidPayload)' /root/astrate-mule/internal/pairing/service.go
41:	// ErrInvalidHWID reports a hw_id that is not a 22-character unpadded
43:	ErrInvalidHWID = errors.New("pairing: invalid hw_id")
44:	// ErrInvalidPayloadFormat reports an initial_payload_format outside
46:	ErrInvalidPayloadFormat = errors.New("pairing: invalid initial_payload_format")
47:	// ErrAlreadyRegistered reports re-registration of a device that has
50:	ErrAlreadyRegistered = errors.New("pairing: device already registered")
51:	// ErrRegistrationLimitReached reports the realm's
54:	ErrRegistrationLimitReached = errors.New("pairing: device registration limit reached")
55:	// ErrUnauthorized is the uniform device-authentication failure: unknown
58:	ErrUnauthorized = errors.New("pairing: unauthorized")
59:	// ErrInhibited reports a device blocked by credentials_inhibited (403).
60:	ErrInhibited = errors.New("pairing: credentials request inhibited")
61:	// ErrInvalidCSR reports an unusable certificate signing request (422).
62:	ErrInvalidCSR = errors.New("pairing: invalid CSR")
169:// secret; afterwards it fails with ErrAlreadyRegistered. initialFormat is
174:		return "", fmt.Errorf("%w: %v", ErrInvalidHWID, err)
177:		return "", fmt.Errorf("%w: %q", ErrInvalidPayloadFormat, initialFormat)
193:			return "", ErrRegistrationLimitReached
208:			return "", fmt.Errorf("%w: %s", ErrAlreadyRegistered, hwID)
240:// device by credentials secret (uniform ErrUnauthorized on any mismatch),
250:		return "", ErrInhibited
259:		if errors.Is(err, ca.ErrInvalidCSR) {
260:			return "", fmt.Errorf("%w: %v", ErrInvalidCSR, err)
327:		return nil, ErrInhibited
365:// device, wrong secret) return the same ErrUnauthorized after a comparable
371:		return nil, nil, ErrUnauthorized
376:		return nil, nil, ErrUnauthorized
384:		return nil, nil, ErrUnauthorized
392:		return nil, nil, ErrUnauthorized
395:		return nil, nil, ErrUnauthorized

Let me check what the service actually returns for register/unregister errors:
$ rg -n 'var Err' /root/astrate-mule/internal/pairing/service.go
(no output)

$ rg -n 'func.*VerifyCredentials' /root/astrate-mule/internal/pairing/service.go
321:func (s *Service) VerifyCredentials(ctx context.Context, realmName, deviceIDStr, secret, clientCrtPEM string) (*VerifyResult, error) {

→ Read internal/pairing/service.go [offset=321, limit=80]
The pairing surface routes match perfectly. Now let me check the `GET /devices/{deviceID}` (info) endpoint more closely for status code drift — I noticed the code has rate limiting but want to verify the docs match:

$ rg -n 'credLimiter|regLimiter' /root/astrate-mule/internal/pairing/http.go
47:	regLimiter   *Limiter
48:	credLimiter  *Limiter
69:		regLimiter:   NewLimiter(cfg.RegisterRate, cfg.RegisterBurst),
70:		credLimiter:  NewLimiter(cfg.CredentialsRate, cfg.CredentialsBurst),
106:	if !a.regLimiter.Allow("ip|" + remoteIP(r).String()) {
154:	if !a.credLimiter.Allow("ip|"+ip.String()) || !a.credLimiter.Allow("dev|"+realm+"/"+deviceID) {

I see that `handleInfo` and `handleVerify` have no rate limiter, and `handleInfo` has no `writeServiceError` call. Let me confirm what happens when `Info()` fails:
→ Read internal/pairing/http.go [offset=197, limit=20]
Everything on the pairing surface lines up well. The specs match the code paths and statuses. Let me now check the `credentials/verify` endpoint — I want to double check the `403 Forbidden` documented response. `handleVerify` calls `writeServiceError` which returns 403 for `ErrInhibited`. That's consistent.

Let me look at one more thing — the `PairingInfo.status` enum. The YAML declares `confirmed, pending, denied, expired`. Let me check what the code actually emits:

$ rg -n 'DeviceStatus' /root/astrate-mule/internal/store/*.go | rg -i 'const|Status' | head -30
/root/astrate-mule/internal/store/devices.go:19:	// DeviceStatusRegistered marks a device that has a credentials secret
/root/astrate-mule/internal/store/devices.go:21:	DeviceStatusRegistered = "registered"
/root/astrate-mule/internal/store/devices.go:22:	// DeviceStatusConfirmed marks a device that has requested credentials
/root/astrate-mule/internal/store/devices.go:24:	DeviceStatusConfirmed = "confirmed"
/root/astrate-mule/internal/store/devices.go:25:	// DeviceStatusInhibited marks a device blocked from new credentials and
/root/astrate-mule/internal/store/devices.go:27:	DeviceStatusInhibited = "inhibited"
/root/astrate-mule/internal/store/devices_test.go:32:		if d.Status != DeviceStatusRegistered || d.Connected || d.PayloadFormatHint != "bson" {
/root/astrate-mule/internal/store/devices_test.go:63:		if d.Status != DeviceStatusConfirmed || d.FirstCredentialsRequest == nil {
/root/astrate-mule/internal/store/devices_test.go:82:		if d, _ = s.GetDevice(ctx, realm.ID, id); d.Status != DeviceStatusInhibited {
/root/astrate-mule/internal/store/devices_test.go:88:		if d, _ = s.GetDevice(ctx, realm.ID, id); d.Status != DeviceStatusConfirmed {
/root/astrate-mule/internal/store/devices_test.go:143:		if d.Status != DeviceStatusRegistered || d.FirstCredentialsRequest != nil || d.CertSerial != nil {

The device status enum in the code is `registered, confirmed, inhibited`, but the YAML documents `confirmed, pending, denied, expired`. That's a genuine drift — the YAML's `PairingInfo.status` enum doesn't match the actual Astrate statuses the handler can return.

Let me check where `info.Status` comes from in the service:
$ rg -n 'func.*Info\b' /root/astrate-mule/internal/pairing/service.go
283:func (s *Service) Info(ctx context.Context, realmName, deviceIDStr, secret string) (*Info, error) {

→ Read internal/pairing/service.go [limit=40, offset=283]
Now I have a clear picture. There are two aspects for the `PairingInfo.status` enum in the YAML:

**Actual statuses the code emits** (`service.go:289-297`): `inhibited`, `confirmed`, or `pending`.
**Documented enum** (`astarte_pairing_api.yaml:367`): `confirmed, pending, denied, expired`.

So the code can emit `inhibited` which is **not** in the documented enum, and the YAML documents `denied` and `expired` which the code never emits. That's real drift.

This is a `docs/api/` fix (generated artefact describing an interface), so it's a normal fix task I can propose. Let me also verify the connections — does the status `inhibited` actually get readable? The `Info` handler calls `authenticateDevice` which returns `ErrInhibited` if `dev.Status == DeviceStatusInhibited`... wait no, let me re-check. `Info` calls `authenticateDevice`, and `authenticateDevice` doesn't check inhibited. Let me verify:

Looking at `authenticateDevice` (service.go:367-398), it does NOT check for inhibited status. So `Info` at line 289-297 does check `dev.Status == store.DeviceStatusInhibited` → status = "inhibited". Good, so `inhibited` is genuinely reachable.

Wait, but actually is `inhibited` reachable through `Info`? `authenticateDevice` returns the device without rejecting inhibited devices. So yes, an inhibited device can query its own info and get `status: inhibited`. So the enum should include `inhibited` and drop `denied`/`expired`.

Let me verify the `Info` in the service — actually I need to double-check whether `authenticateDevice` returns inhibited. Looking again at lines 389-396 it does not check status. Good.

Now let me also verify the docs build still works. But since this is a proposal job, I do NOT edit `docs/api/`. I just propose the fix task and append to `todo.md`.

Let me also check one more endpoint for completeness — the register 403. Actually let me reconsider what I've found and keep to the recipe's guidance (3-4 endpoints checked).

Let me also verify the `verify` endpoint's documented responses against code. The YAML documents `403 Forbidden` for verify, and code returns 403 for `ErrInhibited`. Consistent. But the recipe says the code uses 401 uniformly for auth and 403 only for inhibited — matches.

Now let me quickly check the `status` enum issue is the strongest finding. Let me also confirm whether the `docs/site/configuration-reference.md` check is warranted — that's only if a config key drifted. That's a separate probe; I'll keep the run scoped.

Let me now write the proposal to todo.md. The finding: `docs/api/astarte_pairing_api.yaml` `PairingInfo.status` enum lists `confirmed, pending, denied, expired` but the handler (`internal/pairing/service.go:289-297`, emitted via `internal/pairing/http.go:211`) only ever emits `confirmed`, `pending`, or `inhibited` — `denied`/`expired` are dead and `inhibited` is undocumented.

The recipe says "Five proposals maximum" and this is a proposal job. I have one solid fix to propose plus a status code check. Let me be concise and propose this one task.

Let me check the todo.md current state — it's been modified; the queue is authoritative on `mule/queue` branch per the header note. But the user asked me to "carry it out" which is appending task lines. Let me append.

Actually wait — let me reconsider. Let me check the for-giulio angle. Is this a docs/site contradiction (escalate) or a docs/api fix (propose task)? The status enum is in `docs/api/` (generated), so it's a normal fix task. 

Let me append the task to todo.md.
→ Read .mule/todo.md [offset=98, limit=10]
Now let me make sure I haven't missed other drift. Let me double check that the unregister handler only returns 204/404 etc. as documented. The code at http.go:130-137 returns 204, and `writeServiceError` handles the rest. Consistent with the YAML (204, 401, 403, 404, 500). But note: the agent endpoints require a `a_pa` JWT. `handleRegister`/`handleUnregister` are wrapped in `requireAgent`. The YAML documents 401 for these — matches `requireAgent` returning 401 (via auth middleware). Good.

One more consideration — the `handleCredentials` and `handleVerify` and `handleInfo` rate limit. Only register and credentials have limiters. The YAML documents 429 for register and credentials. Info and verify paths in the YAML don't list 429. Consistent.

So the primary finding is the `PairingInfo.status` enum drift. Let me append the task to todo.md.

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -102,4 +102,5 @@
 - [!] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
 - [!] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto] — BLOCKED: gates failed
 - [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
 - [!] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto] — BLOCKED: gates failed
+- [ ] docs-sync-pairing-status-enum: in docs/api/astarte_pairing_api.yaml, fix the PairingInfo.status enum (line 367): it lists `confirmed, pending, denied, expired` but the handler only ever emits `confirmed`, `pending`, or `inhibited` (internal/pairing/service.go:289-297, via internal/pairing/http.go:211) — drop the dead `denied`/`expired` and add the undocumented `inhibited`. [auto]

This recipe is a proposal job — deliverable is the task line appended to todo.md. No code changes, no git. The status enum is in `docs/api/` so it's a normal fix task, not a `for-giulio.md` escalation (which is only for `docs/site/` prose).

Done: Pairing surface docs-sync — found the PairingInfo.status enum in docs/api/astarte_pairing_api.yaml lists `denied`/`expired` which the code never emits and omits `inhibited` which it does (service.go:289-297, http.go:211); proposed it as a task.
Files: .mule/todo.md
Verified: rg route/path comparison pairing http.go Mount vs astarte_pairing_api.yaml -> all 5 routes match; status codes match -> pass
Unsure: nothing
Follow-ups: docs-sync-pairing-status-enum: fix the PairingInfo.status enum in docs/api/astarte_pairing_api.yaml to {confirmed, pending, inhibited} [auto]
