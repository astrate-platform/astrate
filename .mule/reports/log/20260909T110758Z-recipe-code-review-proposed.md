slug: recipe-code-review
verdict: proposed
at:  95ef44e
ran: 2026-09-09T11:07:58Z on DietPi in 159s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
→ Read .mule/log.md [limit=50]
→ Read .mule/reviews

$ ls internal/ pkg/ 2>/dev/null
internal/:
appengine
auth
broker
config
engine
flow
flowapi
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

$ rg -c '' --glob 'internal/**/*.go' | sort -t: -k2 -rn | head -20
internal/engine/triggers/actions_test.go:1074
internal/flowapi/service.go:976
internal/appengine/channels/ws_test.go:933
internal/appengine/service.go:846
internal/store/datastreams_test.go:823
internal/engine/fakestore_test.go:818
internal/realm/http_test.go:782
internal/engine/triggers/match.go:769
internal/appengine/http.go:696
internal/realm/service.go:686
internal/pairing/http_test.go:649
internal/pairing/service_test.go:637
internal/broker/broker_e2e_test.go:618
internal/engine/triggers/actions.go:610
internal/engine/triggers/match_test.go:594
internal/flow/blocks/httpblocks_test.go:577
internal/realm/http.go:569
internal/flow/router_test.go:565
internal/engine/data.go:552
internal/broker/broker_test.go:549

$ for p in auth config flowapi housekeeping httpx observability pairing realm swagger testutil; do echo -n "$p: "; rg -c '' --glob "internal/$p/**/*.go" 2>/dev/null | awk -F: '{s+=$2} END{print s+0}'; done
auth: 1532
config: 818
flowapi: 2230
housekeeping: 1311
httpx: 357
observability: 398
pairing: 2849
realm: 2573
swagger: 171
testutil: 697

Let me check which packages have been reviewed recently. Recent reviews: appengine (09-07), broker (09-04), engine (09-06), flow (09-05), store (09-08). The unreviewed ones include: auth, config, flowapi, housekeeping, httpx, observability, pairing, realm, swagger, testutil, and the pkg/ ones.

Let me look at the sizes of the unreviewed packages to pick a reasonable one.
$ for p in astarteapi deviceid interfaceschema payload; do echo -n "$p: "; rg -c '' --glob "pkg/$p/**/*.go" 2>/dev/null | awk -F: '{s+=$2} END{print s+0}'; done
astarteapi: 485
deviceid: 358
interfaceschema: 3176
payload: 2673

None of the pkg/ packages have been reviewed. `interfaceschema` and `payload` are substantial. Let me pick `pairing` (2849 lines, unreviewed) or one of the pkg/ ones. 

Given the principle of going deep on one package, let me pick `internal/housekeeping` (1311 lines, unreviewed) — a moderate-size unreviewed package. Actually, let me consider value. The recipe says rotate through areas. Let me look at what's been done: appengine, broker, engine, flow, store. The `pairing` package is a large unreviewed area (2849 lines). That's a good candidate for deep review.

Let me pick `internal/pairing`.
$ ls -la internal/pairing/
total 100
drwxr-xr-x  4 root root  4096 Sep  4 20:46 .
drwxr-xr-x 17 root root  4096 Sep  4 20:46 ..
drwxr-xr-x  2 root root  4096 Jul 27 05:10 ca
-rw-r--r--  1 root root 12461 Sep  4 20:46 http.go
-rw-r--r--  1 root root 21458 Sep  4 20:46 http_test.go
-rw-r--r--  1 root root  1827 Jul 27 05:10 ratelimit.go
-rw-r--r--  1 root root  1932 Jul 27 05:10 ratelimit_test.go
-rw-r--r--  1 root root 16040 Sep  4 20:46 service.go
-rw-r--r--  1 root root 19291 Sep  4 20:46 service_test.go
drwxr-xr-x  3 root root  4096 Jul 27 05:10 testdata

→ Read internal/pairing/service.go
→ Read internal/pairing/http.go
Now let me read the HTTP layer and the ratelimit file.
→ Read internal/pairing/ratelimit.go
Let me look at the tests to understand what's covered, and check the `ca` subpackage. Let me start with the service tests and check coverage of the behaviors I noticed. Let me first note potential findings.

Let me think about what I noticed in the code:

1. **`handleCredentials` rate limiting** (http.go:158): It calls `a.credLimiter.Allow("ip|"+ip.String())` AND `a.credLimiter.Allow("dev|"+realm+"/"+deviceID)`. Due to `||` short-circuit, if the first (IP) allows, the second (device) is evaluated. But here's the subtle bug: if the IP bucket allows, then the device bucket is checked. But **if the first `Allow` returns true and the second returns false, the request is rejected — yet the first `Allow` already consumed a token.** More importantly, the ordering consumes a token from the IP bucket even in some cases. Actually wait — consider: `!Allow(ip) || !Allow(dev)`. If `Allow(ip)` is false, short-circuit → rejected, and the ip token was NOT consumed (it was already below 1). If `Allow(ip)` is true and `Allow(dev)` is false → rejected, ip token consumed, dev not. If both true → allowed, both tokens consumed. This is mostly fine semantically, though it can consume an IP token on a dev-rejected request. Not a big deal.

Actually a more subtle issue: this is two separate Limiter instances? No — `a.credLimiter` is one Limiter with a shared LRU keyed by "ip|..." and "dev|...". So both use the same bucket cache but different keys. Fine.

2. **`remoteIP` / `bearerSecret`** — fine.

3. **`Register` device registration limit**: Ctrl — the limit gates every registration *including secret rotations for existing devices* per the comment (upstream parity). This is documented behavior.

4. **`VerifyCredentials`** (service.go:349): When `EnforceLatestCert` and `dev.CertSerial != nil` — compares. But if `EnforceLatestCert` is on and `dev.CertSerial` is nil (device has no cert yet), the cert is accepted as valid bypassing the latest-cert check. That's a possible edge but probably fine.

5. **`Unregister` (service.go:230-240)**: If `deviceid.Parse` fails, it returns `store.ErrNotFound` wrapped. But note — the realm lookup (`GetRealmByName`) happens *after* parsing. If deviceID invalid, it returns ErrNotFound without looking up realm. OK.

6. **`Register` re-registration / `ErrAlreadyRegistered`**: Interesting: In Register, it calls `CountDevices` for the limit check even when re-registering an existing device. The limit gates rotations too. Once a device requests credentials (confirmed), re-registering returns ErrAlreadyRegistered. That matches upstream.

Now, let me focus on the strongest findings. Let me look at the tests to know what's covered before proposing missing tests.

Let me read the service_test and http_test to see coverage.
→ Read internal/pairing/service_test.go
Good coverage. Now let me read the http_test.go and the ca subpackage to complete my understanding, plus the ratelimit test.

Let me focus. This is quite well-tested. Let me find genuine findings. Let me examine a few things I flagged:

1. **The rate limiter double-Allow in `handleCredentials`** (http.go:158). Order matters: `!a.credLimiter.Allow("ip|"+ip.String()) || !a.credLimiter.Allow("dev|"+realm+"/"+deviceID)`. If IP bucket over limit → short circuit, dev not checked, rejected. If IP allows but dev over → rejected but IP token consumed. This is a minor token-waste issue. Not a strong proposal.

2. **`bearerSecret` matching**: it does `strings.Cut(header, " ")`, so a header with no space returns not found. That matches upstream's `bearer\:?\s+` — requires whitespace. Fine.

3. **`Register` initial payload format validation** happens before realm lookup. Fine.

Let me look at a potential real issue: **the `Limiter` and `API` share `credLimiter` for both IP and device keys in the same LRU**. Not a bug.

Let me consider the **`burnBcrypt` equalization**. In `authenticateDevice`, when `deviceid.Parse` fails at the very start, it calls `burnBcrypt(secret)` — fine. But note: in `Register`, the `SetPayloadFormatHint` happens after `RegisterDevice`. If `RegisterDevice` succeeds but `SetPayloadFormatHint` fails, the device is registered but the hint wasn't set, and the secret is returned — fine, the registration succeeded.

Actually wait — more interesting: in `Register`, after `RegisterDevice` succeeds, if `SetPayloadFormatHint` fails it returns the error. But the device is already registered. On retry, re-registration would rotate the secret... but if the hint keeps failing, the call keeps failing after having registered. Minor, probably fine since store ops rarely fail.

Let me check the `http_test.go` coverage to see what's untested, especially for the limiter and bearer handling.
→ Read internal/pairing/http_test.go
→ Read internal/pairing/ratelimit_test.go
The HTTP tests are behind the `integration` build tag (needs a live TimescaleDB — can't run on this Mac, needs Legion Go). The `service_test.go` and `ratelimit_test.go` are container-free unit tests.

Let me note an important observation: `service_test.go` and `ratelimit_test.go` are the container-free unit tests that can run on the Mac. `http_test.go` is `//go:build integration`.

Let me look at the `ca` subpackage and ratelimit_test to fully cover the pairing package.

Let me focus now. I have a few candidate findings:

1. **`handleCredentials` double-token consumption** (http.go:158) — minor, likely not worth proposing.

2. **`VerifyCredentials` with EnforceLatestCert and nil CertSerial** (service.go:349): If enforcement is on but the device has `CertSerial == nil` (e.g. some path left it unset), a presented valid cert would be accepted even though its serial wasn't recorded as latest. But this is an edge/thin. Actually more subtle: if `dev.CertSerial == nil`, the cert is accepted valid regardless of the serial. Is there a legit state where CertSerial is nil but a cert exists? Probably not. Weak.

3. **A missing test** candidates: Let me think about rules that exist but aren't asserted:
   - `remoteIP` fallback to IPv4 unspecified — is it tested? Probably not.
   - `bearerSecret` permissive matching (`Bearer:` with colon, case-insensitive) — is that tested? Not obviously.
   - The `Register` validation of `initial_payload_format` when invalid — tested in `TestRegisterInitialPayloadFormat` (the "cbor" case).

Let me look at ratelimit_test.go to see coverage of the limiter.

Let me consider the strongest genuinely-wrong behavior. Let me re-examine the rate limiter more carefully.

In `handleCredentials`, the `||` with two `Allow` calls: Suppose an attacker shares an IP. The IP bucket (burst 5, rate 1/s) - a single device legitimately renewing is rare. The device bucket (burst 5). Actually the double check is intentional: limit both by IP and by device. Fine.

Actually, here's a real subtle issue with the early-return ordering: `!Allow(ip) || !Allow(dev)`. When the IP is over-limit, the device bucket is NOT consulted and no device token is eaten — fine. When IP is fine but device is over-limit, IP token is wasted. This is asymmetric but harmless (both are refilled).

Let me look at whether there's an unguarded behavior that's actually wrong. Let me reconsider `Register`:

In `Register`, the registration limit check `CountDevices` counts devices in the realm. Note comment says the limit gates every registration attempt including rotations. So a rotated secret consumes no new device but still gated. Upstream parity. Fine.

Let me look more carefully at a possible real bug: **`Unregister` parse order**. Actually it's fine.

Let me reconsider the `bearerSecret`. Upstream regex is `~r/bearer\:?\s+(.*)$/i`. Astrate's implementation: `strings.Cut(header, " ")` splits on first space; requires the scheme part to be "bearer" with optional colon, case-insensitive. The regex also allows the value to start after any whitespace (`\s+`); Astrate uses a single space in Cut. That's fine for normal headers. Edge: `Bearer:\tsecret` (tab) — Cut on space finds no space → not found → rejected, while upstream would accept. Very edge, not worth.

Let me settle on what's genuinely worth proposing. Let me check the `ca` package for a potential real issue, since that's where crypto lives. Let me read the ca files.

Actually, let me keep my context budget in check. Let me look at the ca subpackage briefly and the ratelimit test.

$ ls internal/pairing/ca/
ca.go
ca_test.go

The ratelimit tests are thorough. The service tests are thorough. This is a well-covered package.

Let me look at the `ca` package since crypto is the highest-value area and I haven't checked it yet. Let me read it.
→ Read internal/pairing/ca/ca.go
The `ca` package is well-written. Let me analyze for findings.

One thing I noticed in `SignCSR`: the serial is `leaf.SerialNumber.String()` which is a **decimal** string, and `VerifyCredentials` in service.go compares `cert.SerialNumber.String()` (decimal) against `*dev.CertSerial`. Both use `.String()` so they're consistent. Good.

Let me now consolidate my findings. This package is very thoroughly tested (unit tests cover service, ca, ratelimit; integration covers HTTP). Let me identify genuine, defensible proposals.

Candidate findings:

**Finding A — `handleCredentials` rate-limit double-consume ordering** (http.go:158): The IP bucket's `Allow` runs before the device bucket check. When the IP is over-limit, short-circuit means the device bucket is never touched — fine. But there's an asymmetry worth a note, not a bug.

Actually, let me reconsider more carefully whether there's a real bug in the credLimiter `||` logic. Consider the scenario designed to matter here. The two `Allow` calls share one Limiter but distinct keys. The intent: rate-limit credentials per-IP **and** per-device. This works. Not a bug.

**Finding B — `VerifyCredentials` when `EnforceLatestCert` is on but `dev.CertSerial == nil`** (service.go:349): A valid, in-window cert is accepted as valid even though no serial was recorded as "latest". Is there a reachable state where a device has credentials but `CertSerial` is nil? `SetDeviceCredentials` always sets CertSerial. The only nil path is a device that never requested credentials — but then `authenticateDevice` fails because `CredentialsSecretHash` would be set only after Register... actually after Register the hash is set and status is Registered, CertSerial nil. Such a device hasn't gotten a cert yet, so it can't present one it was issued. It could present a *foreign* cert, but that fails the chain check first. So this is not exploitable. Weak — not proposing.

**Finding C — `Unregister` does a `deviceid.Parse` and returns `store.ErrNotFound` on parse failure** (service.go:233): This means an agent DELETE with a malformed device ID returns 404 rather than 422. That's arguably correct (matches upstream 404-shape for unknown device). Not a bug.

Let me think about **missing tests** — highest-value proposals. What rules exist but aren't asserted in the container-free unit tests?

1. `remoteIP` fallback to `netip.IPv4Unspecified()` when RemoteAddr is unparseable — not unit-tested (only reachable in integration tests). A small pure function with a clear rule. Good test candidate.
2. `bearerSecret` permissive matching — the `Bearer:` colon form and case-insensitivity — not unit-tested directly. But it's exercised in integration tests only with the plain `Bearer ` form. The permissive `bearer:secret` and case variants aren't tested anywhere. Good candidate.
3. `Register` with `initial_payload_format` set — there IS a test (`TestRegisterInitialPayloadFormat`). Good.
4. In service.go, `VerifyCredentials` with `EnforceLatestCert` and a cert serial mismatch tested (RevokedByRotation). Good.

Also worth checking: **the CRL/serial recording in `SignCSR`**. Is the serial *uniqueness* tested — i.e., two issued certs have distinct serials? A regression where serials collide would be catastrophic security-wise. There's no test asserting two issuances produce different serials. That's a genuinely valuable missing test on a security-critical rule. Let me verify it's not tested.

Let me grep for serial-related tests in ca_test.go.

Let me also check whether the `Verify` precedence (expired before invalid, and the `x509.Expired` fallback inside Verify) is tested.

Let me read ca_test.go.
→ Read internal/pairing/ca/ca_test.go
The `ca` package is exceptionally well-tested, including the serial-uniqueness test (TestSerialUniqueness10k). So serial uniqueness is covered.

Let me now identify genuine, defensible proposals. This package is among the most thoroughly tested I've seen. Let me settle on the strongest findings.

Let me reconsider the `Verify` precedence in ca.go lines 234-240. When `cert.Verify` fails, it checks `invalid.Reason == x509.Expired` via `errors.As`. But note: the function has a documented precedence "expiry first, so an expired certificate reports EXPIRED even when other problems coexist." However, the *first* check at line 223 (`at.After(cert.NotAfter) || at.Before(cert.NotBefore)`) catches expiry using the provided `at`. Then `cert.Verify` is called with `CurrentTime: at`. For a cert whose validity window is fine at `at` but chain fails for another reason, it goes to the `errors.As` branch. All good. Covered by tests (Expired, ForeignCA, Garbage).

Now, let me consider genuine gaps not covered by unit tests:

1. **`remoteIP` fallback** (http.go:359-366) — the parse-failure fallback to `netip.IPv4Unspecified()`. Only integration tests reach this? Actually `remoteIP` is used in `handleRegister` and `handleCredentials`, which are integration-only. So this pure helper is untested. But is there a *rule* here worth asserting? Yes: unparseable RemoteAddr → IPv4 unspecified. It's a small, purely-verifiable, container-free helper test. Reasonable but minor.

2. **`bearerSecret`** (http.go:343-354) — the permissive upstream matching: `Bearer`, `Bearer:`, case-insensitive, and trimming. Only `Bearer ` (space) is exercised in integration tests. The colon form and `bearer` lowercase aren't tested anywhere. My reading of it: the colon form `Bearer:foo` should be accepted (upstream regex `bearer\:?\s+` allows optional colon). Let me trace the code: `strings.Cut("Bearer:foo", " ")` — there's no space, so `found` is false → returns `("", false)` → **rejected**. But upstream's regex `bearer\:?\s+` requires `\s+` after the optional colon, so `Bearer:foo` (no whitespace after colon) would NOT match upstream either. So both agree that `Bearer:foo` is invalid. What about `Bearer: foo` (colon, space)? `strings.Cut` splits on the first space → scheme="Bearer:", rest="foo". `TrimSuffix(scheme,":")` → "Bearer", EqualFold "bearer" → true. Accepted. Good, matches upstream. And `bearer foo` (lowercase) → EqualFold → accepted. Good. Case-insensitive tested? Not directly, but likely fine.

So `bearerSecret` behaves correctly. Adding a unit test would be valuable but minor (the integration tests cover the happy path, not the colon/lowercase variants).

Let me front-load the genuinely most valuable proposal. The recipe's #1 value is "a behaviour that is wrong or unguarded." I haven't found a clear wrong behaviour. Let me look harder at the rate limiter issue and the `Register`/`SetPayloadFormatHint` non-atomicity.

Actually, reconsider **`Register` + `SetPayloadFormatHint`** (service.go:216-220): after `RegisterDevice` succeeds, if `SetPayloadFormatHint` fails, the function returns an error — but the secret was already generated, the device already registered and has a valid secret stored. The HTTP layer returns 500, but the device is actually registered. Worst case: on the client retrying, re-registration would succeed again (rotating the secret) and possibly set the hint. It's not a serious issue since store ops are within one DB transaction-ish. Weak.

Let me reconsider the strongest candidate: **the `handleCredentials`/`handleRegister` don't consume a rate-limit token on the authenticate-before-limit path?** Actually the limiter is called first, before auth. That's correct (rate limit before expensive bcrypt — good DoS hygiene).

Let me look once more at http.go:158 with fresh eyes for a real ordering issue:

```go
if !a.credLimiter.Allow("ip|"+ip.String()) || !a.credLimiter.Allow("dev|"+realm+"/"+deviceID) {
```

The two `Allow` share the same `Limiter` (same mutex, same LRU). Both keys tracked. This is semantically: reject if either IP bucket or device bucket is empty. Fine.

I think the honest conclusion is: **this package is very well-written and well-tested.** Per the recipe, "Finding little in a well-written package is a correct result. Say so and stop."

But let me identify at least the defensible minor proposals that a careful reviewer would raise, in descending order:

1. **Missing container-free unit test for `bearerSecret`'s permissive matching** (the `Bearer:`-colon and lowercase variants that upstream's regex accepts and that integration tests never exercise). This asserts a real wire-compat rule that currently rests only on the happy-path integration test. `internal/pairing/http.go:343` → but `bearerSecret` is in an integration-tagged file (`http.go`). Wait, is `http.go` itself integration-tagged? No — only `http_test.go` is. `http.go` is a normal file. So a `bearerSecret`-only test can live in a non-integration `_test.go`. But `bearerSecret` calls nothing else special. That works as a container-free test. 

2. **Missing unit test for `remoteIP` fallback** to `netip.IPv4Unspecified()` on unparseable RemoteAddr. Pure helper, container-free, rule currently untested.

3. A note about the `handleCredentials` asymmetric token consumption (not proposing as a task — it's a nit).

Let me reconsider whether there's a more valuable finding I'm missing. Let me reconsider the `Unregister` flow — after unregister, device data is kept and it becomes re-registrable. Re-register calls `RegisterDevice` which, in the real store, likely resets. Fine, tested.

Let me reconsider `Register` registration-limit with the `CountDevices` — a TOCTOU between CountDevices and RegisterDevice, but store-level; not for this review.

I think my two missing-test proposals are the right items, plus I should record in the review file what I decided *not* to propose (the cosmetic/nit items and the nil-CertSerial non-issue).

Actually, let me reconsider one more genuinely-wrong candidate: **`dummySecretHash` equalization**. `burnBcrypt` is called for: invalid deviceID, unknown realm, unknown device, and device-with-no-hash. But in the success path (`Register` → Credentials with correct secret), it does ONE bcrypt compare. In the failure paths it also does ONE burn. Good equalization. But note: in `authenticateDevice`, when the secret compares successfully, that IS the same cost. Fine.

However — consider `Info` and `VerifyCredentials` and `Credentials`: all three call `authenticateDevice`, which burns one bcrypt on failure. On the `Credentials` success path, one compare. Uniform. Good.

One subtle thing: `ErrInhibited` devices — `Info` still allows inhibited devices to read info (by design, documented). `Credentials`/`Verify` block them after auth (one bcrypt burned). Consistent.

OK. Let me also double check: is there an unbounded growth concern in the Limiter? It's LRU-capped at 65536 — addressed. 

Let me write the review file and append task lines.

Wait — let me reconsider the strongest "missing test" more carefully to make sure I propose genuinely valuable, non-duplicated tasks. The integration `http_test.go` does cover the happy paths but is behind a build tag and needs the DB (Legion Go). The unit tests cover service + ca + limiter but NOT the pure HTTP helper functions `bearerSecret` and `remoteIP` (which live in `http.go`, a non-integration file, and so are testable without containers).

Proposal candidates:
- **A**: unit-test `bearerSecret` permissive forms (colon, lowercase, whitespace) — internal/pairing/http.go:343
- **B**: unit-test `remoteIP` unparseable fallback — internal/pairing/http.go:359

These are both fully machine-checkable, container-free, and currently untested rules. Both are good "missing test" proposals (recipe value tier #2).

I should also note what I decided not to propose and why (the `||` token asymmetry, the nil-CertSerial non-issue, `Register`/`SetPayloadFormatHint` non-atomicity) so the next run doesn't repeat.

Let me check the date. Today is 2026-09-09. So file: `.mule/reviews/pairing-2026-09-09.md`.

Let me write the review file.

One more: let me double-check that `bearerSecret` is truly reachable/untested. Confirmed it's in http.go (non-integration), used by integration handlers. Let me write it into the review as a legitimately testable unit.

Let me also verify there aren't existing tests for bearerSecret/remoteIP elsewhere.
✱ Grep "bearerSecret|remoteIP" in internal/pairing · 0 matches
Confirmed: `bearerSecret` and `remoteIP` have no tests anywhere. These are pure, container-free, testable helper functions.

Now let me write the review file and append task lines to todo.md.

$ cat .mule/todo.md
# Mule queue

One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).
The mule runs the topmost open line. Reorder freely; the order is the priority.

If a line needs more detail than fits, write `.mule/tasks/<slug>.md` and keep the line short.

    tools/mule.sh add "<slug>: <outcome>"
    tools/mule.sh loop

A line tagged `[legion]` needs the Legion Go and is skipped automatically while it is asleep;
the queue moves on to the next runnable line rather than stalling behind it.

A line tagged `[readonly]` is a **standing check**, not a piece of work. It verifies something
and is expected to produce no diff at all — so an empty diff is success, not the `wrote
nothing` failure it would be for any other line. It never gets ticked off, because a gate is
never done. Its output lands in `.mule/reports/<slug>.md` with the sha it ran against, and it
is skipped while that sha is still HEAD: re-checking code that has not moved buys nothing and
spends a call on a free provider.

**`mule/queue` is authoritative for this file.** The mule ticks on the Pi and writes its
progress there, so editing the queue on `main` resurrects completed tasks the moment the two
are merged — which has happened once. Add and reorder tasks on `mule/queue`; let them reach
`main` only when that branch is merged.

A line marked `- [~]` is **parked**: real work, but not tick-sized. The mule only ever
picks up `- [ ]`. Benchmark runs live here because a single tier is 5-20 minutes of ingest
alone -- past the per-task budget -- and because they want someone watching. Run one by hand:

    tools/mule.sh legion bench-push
    ssh legion 'cd ~/astrate/bench && ./scripts/run-tier.sh small astrate -base-url ... -housekeeping-key ...'

## Where tasks come from

**This file is not the whole queue, and for real work it is not even the main part of it.**

The queue is: the standing lines below, plus **every open GitHub issue labelled `mule`**.
Issues are read live on each tick and are never copied into this file — a copy would be a
second place the same fact lives, on a branch the mule commits to and you edit on `main`,
and that produced three merge conflicts in one afternoon.

**To give the mule work, file an issue and label it `mule`.** From anywhere, by anyone,
including another model with repo access. No SSH, no editing this file:

    gh issue create --label mule --title "<slug>: <outcome>" --body "<the detail>"

Labels on the issue are the tags: `legion` and `readonly` mean what `[legion]` and
`[readonly]` mean here. State lives on the issue, as labels, because there is exactly one
copy of it there:

| label          | meaning                                                          |
|----------------|------------------------------------------------------------------|
| `mule`         | queued                                                            |
| `mule-review`  | the mule pushed something; **it is not merged and not reviewed**  |
| `mule-blocked` | it tried and could not; re-label `mule` to try again              |

The mule never closes an issue. Whether the work actually resolves it is a judgement about
intent, which is the reviewer's call.

When both sources are empty a tick runs a **proposal recipe** instead, rotating through
`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
stuck re-reviewing the same code. Lines it invents are tagged `[auto]`: nobody approved those.
A refill costs a tick from the daily budget and never runs what it just proposed — the lines
sit for one tick, which is your window to cut a bad one. `MULE_NO_REFILL=1` turns it off.

## Nothing merges on its own

Everything lands on `mule/queue`. The gates prove a change compiles, passes the tests, ships
a test that fails without it, and touches no frozen file — none of which means the change is
worth having. Before any of it reaches `main`:

    bash tools/mule.sh review

## Queue

- [x] bench-tiers: create the tiered benchmark definitions per .mule/recipes/benchmarks.md (first run only — this task builds the harness, it does not run it)
- [~] bench-small-astrate: run `bench/scripts/run-tier.sh small astrate` against local Astrate, commit results (two runs minimum) [legion]
- [~] bench-medium-astrate: run `bench/scripts/run-tier.sh medium astrate` against local Astrate, commit results (two runs minimum) [legion]
- [~] bench-small-astarte: run `bench/scripts/run-tier.sh small astarte` against local Astarte, commit results (two runs minimum) [legion]
- [~] bench-medium-astarte: run `bench/scripts/run-tier.sh medium astarte` against local Astarte, commit results (two runs minimum) [legion]
- [~] bench-big-astrate [legion]: run `bench/scripts/run-tier.sh big astrate` against Legion Go Astrate, commit results (two runs minimum)
- [~] bench-giant-astrate [legion]: run `bench/scripts/run-tier.sh giant astrate` against Legion Go Astrate, commit results (two runs minimum)
- [~] bench-big-astarte [legion]: run `bench/scripts/run-tier.sh big astarte` against Legion Go Astarte, commit results (two runs minimum)
- [~] bench-giant-astarte [legion]: run `bench/scripts/run-tier.sh giant astarte` against Legion Go Astarte, commit results (two runs minimum)
- [ ] race-check-store: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/store/... ./internal/housekeeping/... ./migrations/...`. Report any failure to .mule/for-giulio.md with the full race report. Split out of the former single `race-check` line, which timed out running the whole tree at once. [legion] [readonly]
- [ ] race-check-engine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/engine/... ./internal/broker/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
- [ ] race-check-flow: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/flow/... ./internal/realm/... ./internal/pairing/... ./internal/auth/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
- [ ] race-check-appengine: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./internal/appengine/... ./internal/observability/... ./internal/httpx/... ./internal/config/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
- [ ] race-check-pkg: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./pkg/... ./cmd/... ./internal/testutil/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]

- [x] broker-acl-coldstart-introspection-miss: in `internal/broker/aclhook.go` `OnACLCheck` (lines 183-195), when a device publishes to an interface introspected after connect, `refreshIfStale` is skipped for the first second (admit stamps `lastIntroLoad` at authhook.go:404, debounce is authhook.go:186) and the recheck re-reads the still-cold cache — a denied QoS0 publish is silently dropped by mochi (processPublish server.go:867-873). Fix the miss path to fall back to a synchronous store read for the unknown interface when the debounce skips the reload, and add a T1 test in `broker_test.go` (fake store, no Docker) that connects with an empty-introspection store, adds the interface+introspection to the store after connect (stamp `sess.lastIntroLoad` to de-flake), and asserts a QoS0 publish to that interface reaches the intake. [approved 2026-09-04]
- [!] broker-disconnect-device-zombie-session: in `internal/broker/broker.go` `DisconnectDevice` (lines 260-266), force-closing the live connection leaves the persisted bbolt session (subscriptions + offline queue, keyed by CN) in place — mochi derives `expire` from the Clean flag, not the reason code (server.go:485-491), so `dropClient` (sessionstore.go:500) is never reached and a later device with the same CN resurrects the deleted device's session. Wire the session store handle into `Broker` and call `dropClient` from `DisconnectDevice`, and add a T1 test in `broker_test.go`: connect clean=false + subscribe, call `DisconnectDevice`, reconnect — assert `session_present` is false and the stale subscription is gone. [approved 2026-09-04] — BLOCKED: wrote nothing
- [!] broker-offline-acl-tests: in `internal/broker/aclhook_test.go`, unit-test the offline-delivery ACL — `offlineACL.ownershipOf` cache hit within TTL, TTL expiry triggering a reload, and a load failure caching as empty and denying (aclhook.go:116-138) — plus a T1 `OnACLCheck` check that the offline branch (delivery to a session-present but disconnected device, aclhook.go:196-202) consults the store. The rule currently has zero direct tests; `TestCheckACLMatrix` exercises only the pure `checkACL`. [approved 2026-09-04] — BLOCKED: wrote nothing
- [!] broker-onconnect-doc-comment: in `internal/broker/authhook.go:314`, restore the missing first line of the `OnConnect` doc comment — it opens mid-sentence ("of the Astarte MQTT v1 protocol, and mochi publishes them without a publish-side ACL check…") and the Will-clearing security rationale (a retained LWT would escape the §3.2 matrix) is only half-documented. [approved 2026-09-04] — BLOCKED: tests failed: --- FAIL: TestMQTTSink_Retained (0.02s)

- [x] purge-properties-compression-capability [auto]: upstream v1.3.0 adds a `purge_properties_compression_format` device capability (`zlib`|`plaintext`, default `zlib`) — a wire-visible capability value. Check whether Astrate's capabilities handling (internal/broker, the `<realm>/<device_id>/capabilities` topic, issue #16) needs to recognize/honour it, or whether zlib-only is already the deliberate default; propose the change or note why not needed.
- [!] empty-introspection-verification [auto]: upstream v1.3.0 changed "allow devices with empty introspection" — verify whether Astrate's device connection/introspection handling currently rejects an empty introspection string where upstream now accepts it, and propose a fix if so. — BLOCKED: wrote nothing
- [!] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages" — investigate only: does Astrate have an analogous delay between trigger install and first delivery? Report, do not patch. — BLOCKED: wrote nothing
- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
- [x] docs-sync-pairing-health-path [auto]: add the undocumented `GET /pairing/v1/{realm}/health` route to docs/api/astarte_pairing_api.yaml — it exists in code (internal/pairing/http.go:81) and returns 200 `{"data":{"status":"ok"}}` for an existing realm, 404 for an unknown realm, 503 when unhealthy, unauthenticated and unrate-limited. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-register-404 [auto]: document the missing `404` (DeviceNotFound) on `POST /pairing/v1/{realm}/agent/devices` in docs/api/astarte_pairing_api.yaml — the handler returns it when `GetRealmByName` fails for an unknown realm (internal/pairing/service.go:183-186 → http.go:326-327), but the spec lists only 201/400/401/403/422/429/500. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] swagger-httptest-coverage [auto]: add a container-free httptest suite for internal/swagger, which currently has no test file — cover two rules: (1) `Mount` wires `GET /swagger` to a 302→/swagger/index.html, serves the embedded UI under `/swagger/`, and serves the OpenAPI YAML specs under `/api/`; (2) `Specs` returns exactly the `.yaml` filenames from `docs.APIYAML` with no path prefix or dirs. httptest only, no Docker.
- [!] probe-property-resend-encoding [auto]: upstream v1.3.0 fixed outbound server-property values sent to a device on connect/emptyCache (commit 522ccf4f — a raw stored binaryblob is now wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding, so it ships as a BSON binary, not a string). Investigate only: does Astrate's resendServerProperties/rehydrate path (internal/engine/control.go:144, pkg/payload) emit binaryblob → BSON subtype-0 binary and datetime → UTC datetime exactly like upstream 1.3? Report, do not patch. — BLOCKED: wrote nothing
- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
- [!] flow-validate-source-sink: `Pipeline.Validate`'s no-source/no-sink error branches (internal/flow/pipeline.go:148-153) are unreachable — the preceding Kahn cycle check (pipeline.go:100-115) guarantees any graph that gets past it is acyclic and therefore always has a source and a sink. The two test cases "no source (all blocks have incoming edges)" and "no sink (all blocks have outgoing edges)" (pipeline_test.go:134-163) both use a 2-cycle, so they pass via the cycle error, never the branch they name. Make the rule real (a cycle-independent formulation, e.g. reject a graph where a node has no path to a sink) or delete the vacuous check, and fix the two test cases to assert the actual error message so a regression to wrong-path coverage fails. — BLOCKED: tests failed: --- FAIL: TestMQTTSink_Retained (0.02s)
- [x] flow-validate-dead-source-sink-recompute: in `Pipeline.Validate`, the first `hasSource`/`hasSink` computation (pipeline.go:117-126) is overwritten by the second (pipeline.go:128-147) and its result discarded — dead code (the comment at pipeline.go:127 documents it). Delete the first block; behaviour is unchanged. Tests: existing `TestPipeline_Validate` passes.
- [x] docs-sync-rm-datastream-retention-endpoint [auto]: add the undocumented `GET /realmmanagement/v1/{realm}/config/datastream_maximum_storage_retention` route to docs/api/astarte_realm_management_api.yaml — it exists in code (internal/realm/http.go:51, handler getDatastreamMaximumStorageRetention at http.go:135, returns 200 with a data-envelope retention value via Service.GetDatastreamMaximumStorageRetention) but the spec jumps straight from config/device_registration_limit (yaml:428) to /version. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-interfaces-detailed-param [auto]: add the `?detailed=true` query parameter to `GET /realmmanagement/v1/{realm}/interfaces` in docs/api/astarte_realm_management_api.yaml — code serves a 1.4-style detailed interface listing when the param is `true` (internal/realm/http.go:150, issue #66); the spec documents only the names-only response (yaml:27-51) and neither the param nor the detailed response. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] housekeeping-tests [auto]: add a container-free unit test file for `internal/housekeeping` (service.go, 262 lines, zero test coverage). Cover two rules: (1) `CreateRealm` validation — blank name, blank JWT key, negative regLimit, negative retention, and default retention injection when `WithDefaultDatastreamMaximumStorageRetention` is set and the caller passes nil; (2) `DeleteRealm` gating — `ErrDeletionDisabled` when `WithRealmDeletionDisabled(true)`, `ErrConnectedDevicesPresent` when `DeviceStats` reports connected > 0. Use httptest-style fakes for the store and sealer (no Docker).
- [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
- [x] docs-sync-appengine-by-alias-endpoints [auto]: add the 7 undocumented by-alias mirror routes to docs/api/astarte_appengine_api.yaml — PATCH /devices-by-alias/{alias}, GET /devices-by-alias/{alias}/interfaces, GET/PUT/POST/DELETE /devices-by-alias/{alias}/interfaces/{interface}/{path} (plus the interface-scoped GET) — all registered in internal/appengine/http.go:59-65 and backed by handlers patchDeviceByAlias, listInterfacesByAlias, getDataByAlias, putDataByAlias, deleteDataByAlias (http.go:226-367). Match the existing spec style (same schemas, same response codes as the device-scoped counterparts). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-appengine-group-endpoints [auto]: add the 8 undocumented group data routes to docs/api/astarte_appengine_api.yaml — GET /groups/{group} (returns {data: {group_name: "<name>"}}, http.go:428-434), PATCH /groups/{group}/devices/{device} (DevicePatch body → updated status, http.go:518-537), and the 6 group interface-data routes: the interface-level GET and the path-level GET/PUT/POST/DELETE on /groups/{group}/devices/{device}/interfaces (http.go:75-80, handlers listInterfacesInGroup, getDataInGroup, putDataInGroup, deleteDataInGroup at http.go:187-377). Also add GET /devices/{device}/interfaces (listDeviceInterfaces, http.go:48) which is missing from the spec — that is the 9th route. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-appengine-get-group-device [auto]: add the undocumented `GET /groups/{group}/devices/{device}` route to docs/api/astarte_appengine_api.yaml — it exists in code (internal/appengine/http.go:72, handler getGroupDevice at http.go:509-516, returns the device status projection behind the membership gate) but the spec path only documents delete and patch. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-appengine-query-params-status [auto]: fix missing query parameters and status codes in docs/api/astarte_appengine_api.yaml — (1) GET /groups/{group}/devices is missing `from_token` and `limit` pagination params that code reads at http.go:439-441 (and returns via links.next), (2) POST /groups is missing 422 validation-error response (FieldErrors for blank group_name/empty devices, http.go:404-416), (3) PATCH /devices/{device} and PATCH /devices-by-alias/{alias} are missing 422 ("Attribute key not found" via ErrAttributeKeyNotFound, service.go:352-357 → http.go:656-657), and the device-scoped PATCH is also missing 409 ("Alias already in use", http.go:651) and its documented 500 covers the merge-patch content-type mismatch which code answers with an unmapped 500 too, (4) GET /devices/{device}/interfaces/{interface}/{path} is missing 422 (invalid since/since_after/to/limit/downsample_to/format params, http.go:554-639). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] server-data-trigger-bus [auto]: the server- and device-owned publish paths `PublishServerValue` / `PublishDeviceValue` / `publishAsOwner` (internal/engine/serverdata.go:94-203) persist ops and publish to the broker but never fire data triggers (`incoming_data`/`value_stored`/`value_change`/…) nor fan out to the live stream bus — unlike the broker-ingress path, which reaches `fireCommitted` via `finalize`→`afterCommit` (batch.go:186, engine.go:140). `UnsetServerProperty` (serverdata.go:211) is likewise silent. The virtual-device contract ("lands storage rows exactly like a real device's data would", devicedata_test.go:3-4) implies these should also fire. Route the persisted op through the same trigger/bus emission as `fireData`, and add a failing test in serverdata_test.go/devicedata_test.go that subscribes to the bus and installs a trigger then asserts both fire on a server-owned and a device-owned publish. Confirm with Giulio that server-originated writes are not meant to be deliberately silent before committing.
- [x] docs-sync-appengine-group-patch-status [auto]: PATCH /groups/{group}/devices/{device} (internal/appengine/http.go:518-537, patchGroupDevice) shares applyPatch with the device-scoped PATCH and also answers 422 ("Attribute key not found", ErrAttributeKeyNotFound, service.go:352-357) and 409 ("Alias already in use", http.go:651), but docs/api/astarte_appengine_api.yaml lists only 200/400/401/404/500 for it — add the two missing responses.
- [x] docs-sync-appengine-data-422-interface-level [auto]: the interface-level data GETs — GET /devices/{device}/interfaces/{interface}, GET /devices-by-alias/{alias}/interfaces/{interface}, GET /groups/{group}/devices/{device}/interfaces/{interface} — share parseQueryOpts (internal/appengine/http.go:554-639) and answer 422 FieldErrors for invalid since/since_after/to/limit/downsample_to/format params, but the spec documents only the path-level GET with 422 — add 422 to the interface-level GETs.
- [!] appengine-snapshot-ignores-query-params [auto]: in `internal/appengine/data.go:138-148`, the individual-datastream interface-root snapshot branch ignores `Since`/`SinceAfter`/`To`/`Limit`/`Descending` even though the root GET is documented with all of them (docs/api/astarte_appengine_api.yaml:577-582) — the same function refuses incompatible params elsewhere (downsample 422, data.go:118-120/165-167), so this is a silent drop. Probe upstream on the Legion Go to decide honour-vs-refuse, then either apply the window/limit or answer an explicit 422 and stop advertising the params on the root GET; add a T1 test in data_test.go/query_opts_test.go pinning the chosen behaviour. — BLOCKED: wrote nothing
- [!] appengine-group-token-roundtrip-test [auto]: add a container-free unit test for `parseGroupToken`/`groupTokenFor` (internal/appengine/service.go:602-623) — currently only exercised via the pg-backed groups_parity_test.go. Assert (1) round-trip `groupTokenFor(offset)` → `parseGroupToken` recovers offset, (2) a non-v1-nibble UUID is rejected, (3) a well-formed v1 UUID with zeroed time fields is accepted (parseGroupToken comment, service.go:614-616). — BLOCKED: lint failed: internal/appengine/groups_token_test.go:14:5: redefines-builtin-id: redefinition of the built-in function max (revive)
- [x] docs-sync-appengine-data-output-params [auto]: document the output/format query params `format`, `allow_bigintegers`, `allow_safe_bigintegers` (accepted by parseQueryOpts, internal/appengine/http.go:601-623, consumed in data.go:230-244,403-405) on the six data GET operations of docs/api/astarte_appengine_api.yaml, which document only since/since_after/to/limit/downsample_to/sort; `retrieve_metadata` and `downsample_key` are parsed but unconsumed — detail in .mule/tasks/docs-sync-appengine-data-output-params.md. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
- [x] detailed-listing-required-encrypted-flags: emit the `required`/`encrypted` mapping flags in `writeDetailedMapping` (internal/realm/detailed.go) so `GET /realmmanagement/v1/{realm}/interfaces?detailed=true` round-trips them — the detailed listing currently materializes a fixed field set that drops both flags even when a datastream mapping declares one, while the plain GET preserves them via the stored definition. Datastream branch only (properties reject both at parse); add a T1 assertion in internal/realm/detailed_test.go.

- [x] store-alias-values-taken-test: add a store test for `AliasValuesTaken` (internal/store/devices.go:343, the `alias_already_in_use` ownership check behind PATCH /devices) pinning the three branches that currently have no direct coverage — no other device in the realm carries any of the values → false; another device carries one → true; the same device (excluded via `d.id <> $2`) carries the value → false (must not count itself). Add to internal/store/devices_test.go. [auto]
- [x] store-latest-individual-test: add a store test for `LatestIndividual` (internal/store/datastreams.go:247, the previous-value lookup feeding value_change/path_created triggers) — only exercised through the engine's fakeStore today; assert (a) of two samples on one series, `ORDER BY ts DESC LIMIT 1` returns the newest with the right value column, (b) an empty series yields ErrNotFound. Add to internal/store/datastreams_test.go. [auto]
- [x] docs-sync-appengine-device-status-schema [auto]: fix the `DeviceStatus` schema in docs/api/astrate_appengine_api.yaml — the wire emits `id` (internal/appengine/service.go:106 `json:"id"`, decoded as `ds.ID` in mirror_test.go:156 and dashboard_compat_test.go:100/149) but the spec names the field `device_id` (yaml:1409) in both the schema and the list/detail examples (yaml:67, 75); `introspection` is documented as `map<string,integer>` (yaml:1438-1441) but each value is a `{major, minor}` object (introspectionEntry, service.go:123-127, introspectionView service.go:833-836); and the schema omits the six fields `groups`, `total_received_msgs`, `total_received_bytes`, `first_credentials_request`, `last_seen_ip` (omitempty), `previous_interfaces` (omitempty) that deviceStatusView emits (service.go:254-279). Rename `device_id`→`id`, correct `introspection` to `map<string,{major:int,minor:int}>`, add the missing fields, and fix the examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-appengine-data-set-422 [auto]: add the missing `422` response to the six PUT/POST data-set operations in docs/api/astrate_appengine_api.yaml — `putData`, `putDataByAlias`, `putDataInGroup` (internal/appengine/http.go:292-309, 321-338, 340-357) return 422 `{"errors":{"detail":"Value size exceeds size limits"}}` when the body value exceeds 65536 bytes (maxValueBytes, http.go:24), both for PUT and its POST alias (routes registered at http.go:52-53, 63-64, 78-79), but the spec lists only 200/400/401/404/500 for all six operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-appengine-downsample-min [auto]: fix the `downsample_to` constraint in docs/api/astrate_appengine_api.yaml — the DataDownsample description says "(must be >= 2)" (yaml:1357) but parseQueryOpts rejects `downsample_to <= 2` with "must be greater than 2" (internal/appengine/http.go:592-593), so the documented minimum is off by one: readings of `2` (allowed per the doc) actually 422. Change the description to "(must be > 2)" (the 422 error shape is already documented). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

The todo queue is long with existing blocked tasks. My code-review proposals should be appended. I have two solid missing-test proposals. Let me pick the strongest 2 and keep them tight.

Let me now write the review file. First, let me double check `remoteIP` behavior once more for the review notes, and confirm my proposal about `bearerSecret` - both pure helpers in a non-integration file.

Let me also decide whether to propose the rate-limit asymmetric token-consumption note. I'll note it in the "decided not to propose" section.

Let me write the review file and append the task lines.
$ ls .mule/reviews/
appengine-2026-09-07.md
broker-2026-09-04.md
engine-2026-09-06.md
flow-2026-09-05.md
store-2026-09-08.md

Now let me write the review file.
← Write .mule/reviews/pairing-2026-09-09.md
Wrote file successfully.

Now let me verify my claims about `bearerSecret` once more against upstream — I want to make sure the proposal task line describes correct expected behavior. The upstream Astarte pairing regex is `~r/bearer\:?\s+(.*)$/i`. Let me confirm this is what the code comment claims (it does, at http.go:342). Since the code comment explicitly documents this regex, and the code implements it, my proposal to add tests locking that behavior is sound.

Let me append the two task lines to todo.md. I need to add them at the bottom of the queue. Let me look at how lines are appended - per MULE.md, one line each. I'll add them with clear slugs and file references.
$ tail -5 .mule/todo.md
- [x] store-alias-values-taken-test: add a store test for `AliasValuesTaken` (internal/store/devices.go:343, the `alias_already_in_use` ownership check behind PATCH /devices) pinning the three branches that currently have no direct coverage — no other device in the realm carries any of the values → false; another device carries one → true; the same device (excluded via `d.id <> $2`) carries the value → false (must not count itself). Add to internal/store/devices_test.go. [auto]
- [x] store-latest-individual-test: add a store test for `LatestIndividual` (internal/store/datastreams.go:247, the previous-value lookup feeding value_change/path_created triggers) — only exercised through the engine's fakeStore today; assert (a) of two samples on one series, `ORDER BY ts DESC LIMIT 1` returns the newest with the right value column, (b) an empty series yields ErrNotFound. Add to internal/store/datastreams_test.go. [auto]
- [x] docs-sync-appengine-device-status-schema [auto]: fix the `DeviceStatus` schema in docs/api/astrate_appengine_api.yaml — the wire emits `id` (internal/appengine/service.go:106 `json:"id"`, decoded as `ds.ID` in mirror_test.go:156 and dashboard_compat_test.go:100/149) but the spec names the field `device_id` (yaml:1409) in both the schema and the list/detail examples (yaml:67, 75); `introspection` is documented as `map<string,integer>` (yaml:1438-1441) but each value is a `{major, minor}` object (introspectionEntry, service.go:123-127, introspectionView service.go:833-836); and the schema omits the six fields `groups`, `total_received_msgs`, `total_received_bytes`, `first_credentials_request`, `last_seen_ip` (omitempty), `previous_interfaces` (omitempty) that deviceStatusView emits (service.go:254-279). Rename `device_id`→`id`, correct `introspection` to `map<string,{major:int,minor:int}>`, add the missing fields, and fix the examples. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-appengine-data-set-422 [auto]: add the missing `422` response to the six PUT/POST data-set operations in docs/api/astrate_appengine_api.yaml — `putData`, `putDataByAlias`, `putDataInGroup` (internal/appengine/http.go:292-309, 321-338, 340-357) return 422 `{"errors":{"detail":"Value size exceeds size limits"}}` when the body value exceeds 65536 bytes (maxValueBytes, http.go:24), both for PUT and its POST alias (routes registered at http.go:52-53, 63-64, 78-79), but the spec lists only 200/400/401/404/500 for all six operations. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-appengine-downsample-min [auto]: fix the `downsample_to` constraint in docs/api/astrate_appengine_api.yaml — the DataDownsample description says "(must be >= 2)" (yaml:1357) but parseQueryOpts rejects `downsample_to <= 2` with "must be greater than 2" (internal/appengine/http.go:592-593), so the documented minimum is off by one: readings of `2` (allowed per the doc) actually 422. Change the description to "(must be > 2)" (the 422 error shape is already documented). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ cat >> .mule/todo.md << 'EOF'

- [ ] pairing-bearer-secret-test [auto]: add a container-free unit test for `bearerSecret` (internal/pairing/http.go:343) locking the upstream-verbatim permissive match `~r/bearer\:?\s+(.*)$/i` — the colon form `Bearer: <secret>`, the case variants (`bearer `, `BEARER `), and the reject paths (no whitespace before value, non-bearer scheme, empty value). Only the plain `Bearer <secret>` form is exercised today (integration http_test.go); a regression to case-sensitive or drop-colon matching would pass every existing test. Pure helper in a non-integration file; no Docker.
- [ ] pairing-remoteip-fallback-test [auto]: add a container-free unit test for `remoteIP` (internal/pairing/http.go:359) — a valid addr:port yields `ap.Addr().Unmap()`, a bare address yields the mapped address, and an unparseable `RemoteAddr` (unix-socket/proxy case) falls back to `netip.IPv4Unspecified()`. That fallback branch is currently uncovered (only reachable through the integration HTTP stack). Pure helper in a non-integration file; no Docker.
EOF
echo appended
appended

Now let me verify the todo additions are in place and report.

Let me confirm the review file and todo are consistent.
$ tail -2 .mule/todo.md
- [ ] pairing-bearer-secret-test [auto]: add a container-free unit test for `bearerSecret` (internal/pairing/http.go:343) locking the upstream-verbatim permissive match `~r/bearer\:?\s+(.*)$/i` — the colon form `Bearer: <secret>`, the case variants (`bearer `, `BEARER `), and the reject paths (no whitespace before value, non-bearer scheme, empty value). Only the plain `Bearer <secret>` form is exercised today (integration http_test.go); a regression to case-sensitive or drop-colon matching would pass every existing test. Pure helper in a non-integration file; no Docker.
- [ ] pairing-remoteip-fallback-test [auto]: add a container-free unit test for `remoteIP` (internal/pairing/http.go:359) — a valid addr:port yields `ap.Addr().Unmap()`, a bare address yields the mapped address, and an unparseable `RemoteAddr` (unix-socket/proxy case) falls back to `netip.IPv4Unspecified()`. That fallback branch is currently uncovered (only reachable through the integration HTTP stack). Pure helper in a non-integration file; no Docker.

Let me produce the report per MULE.md format.

```
Done: reviewed internal/pairing (service, http, ca, ratelimit) and proposed 2 missing-test tasks
Files: .mule/reviews/pairing-2026-09-09.md, .mule/todo.md (2 appended lines)
Verified: grep for bearerSecret/remoteIP test coverage (none) + read the package; no commands altered code
Unsure: nothing
Follow-ups: pairing-bearer-secret-test, pairing-remoteip-fallback-test (both tagged [auto])
```

Done: reviewed internal/pairing (service, http, ca, ratelimit) and proposed 2 missing-test tasks
Files: .mule/reviews/pairing-2026-09-09.md, .mule/todo.md (2 appended lines)
Verified: grep for bearerSecret/remoteIP test coverage (none) + read the package; no code changed
Unsure: nothing
Follow-ups: pairing-bearer-secret-test, pairing-remoteip-fallback-test (both [auto])

The package was well-tested (service/CA/limiter unit coverage is deep). The only real gaps are two pure helpers in the non-integration `http.go` — the permissive `bearerSecret` bearer match and the unparseable-address fallback in `remoteIP` — each currently untested and each asserted only by an untestable-on-this-Mac integration suite. I recorded both as the sole proposals, and logged in the review file what I looked at and deliberately did not propose (the limiter token asymmetry, the nil-CertSerial edge, the hint non-atomicity) so the next run doesn't repeat the exercise.
