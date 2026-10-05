slug: recipe-hygiene
verdict: timeout
at:  d9f6801
ran: 2026-10-05T19:23:12Z on DietPi in 1441s


> build · big-pickle

I'll start by reading the required context files.
→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md

# Todos
[•] Survey repo layout and .mule state (todo.md, for-giulio.md)
[ ] Run dep sweep `go list -m -u all` + check milestone-boundary standing rule
[ ] Run govulncheck and triage reachable advisories
[ ] Run golangci-lint, group findings not already excluded
[ ] Scan TODO/FIXME/XXX/HACK in internal/, pkg/, cmd/
[ ] Find packages with no test file; find skipped tests
[ ] Append proposed task lines to .mule/todo.md (<=8)
[ ] Write any evidence file the recipe requires

$ ls -la && echo "--- .mule ---" && ls -la .mule/ && echo "--- recipes ---" && ls .mule/recipes/
total 38320
drwxr-xr-x 17 root root     4096 Sep 25 19:41 .
drwx------ 29 root root     4096 Sep 22 12:57 ..
drwxr-xr-x  3 root root     4096 Sep  4 20:46 .agents
-rw-r--r--  1 root root      207 Sep  4 20:46 .dockerignore
drwxr-xr-x  8 root root     4096 Oct  5 20:14 .git
drwxr-xr-x  3 root root     4096 Jul 27 05:10 .github
-rw-r--r--  1 root root     1133 Sep  4 21:26 .gitignore
-rw-r--r--  1 root root     1067 Jul 27 05:10 .golangci.yml
drwxr-xr-x  7 root root     4096 Oct  5 20:14 .mule
-rw-r--r--  1 root root     1578 Sep  4 20:46 AGENTS.md
-rw-r--r--  1 root root     1315 Jul 27 05:10 Dockerfile
-rw-r--r--  1 root root    14394 Jul 27 05:10 LICENSE
-rw-r--r--  1 root root     2545 Jul 27 05:10 Makefile
-rw-r--r--  1 root root     3482 Sep  4 20:46 README.md
-rwxr-xr-x  1 root root 30903300 Sep 25 19:41 astrate
drwxr-xr-x  5 root root     4096 Jul 27 05:10 bench
drwxr-xr-x  3 root root     4096 Jul 27 05:10 cmd
drwxr-xr-x  4 root root     4096 Jul 27 05:10 deploy
-rw-r--r--  1 root root     4317 Jul 27 22:18 docker-compose.yml
drwxr-xr-x  6 root root     4096 Sep  4 20:46 docs
drwxr-xr-x  4 root root     4096 Sep  4 20:46 examples
-rwxr-xr-x  1 root root  8175558 Sep 23 19:04 flow-container-echo
-rw-r--r--  1 root root     5092 Sep  4 20:46 go.mod
-rw-r--r--  1 root root    20833 Sep  4 20:46 go.sum
drwxr-xr-x 17 root root     4096 Sep  4 20:46 internal
drwxr-xr-x  2 root root     4096 Sep  4 20:46 migrations
-rw-r--r--  1 root root      102 Sep  4 20:46 opencode.json
drwxr-xr-x  6 root root     4096 Jul 27 05:10 pkg
drwxr-xr-x 30 root root     4096 Oct  4 13:37 site-dist
drwxr-xr-x  4 root root     4096 Jul 27 05:10 test
drwxr-xr-x  3 root root     4096 Sep  4 21:26 tools
--- .mule ---
total 12544
drwxr-xr-x  7 root root     4096 Oct  5 20:14 .
drwxr-xr-x 17 root root     4096 Sep 25 19:41 ..
-rw-r--r--  1 root root       13 Oct  5 20:59 .budget
-rw-r--r--  1 root root      364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root       11 Oct  5 20:14 .heartbeat
-rw-r--r--  1 root root        0 Oct  5 20:59 .last-output
-rw-r--r--  1 root root        0 Oct  5 20:59 .lock
-rw-r--r--  1 root root        8 Oct  5 20:59 .rotation
-rw-r--r--  1 root root     7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root     2762 Aug 31 01:45 config
-rw-r--r--  1 root root 12456217 Oct  5 20:59 cron.log
drwxr-xr-x  2 root root     4096 Oct  4 22:35 failed
-rw-r--r--  1 root root   127766 Oct  4 20:31 for-giulio.md
-rw-r--r--  1 root root     1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root     1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    14600 Oct  5 20:14 log.md
-rw-r--r--  1 root root     8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root     4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root     4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root     4096 Oct  4 21:08 reviews
drwxr-xr-x  2 root root     4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   149912 Oct  5 20:14 todo.md
--- recipes ---
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

$ git status --short && echo "--- branch ---" && git branch --show-current && echo "--- head ---" && git log --oneline -8
--- branch ---
mule/queue
--- head ---
d9f6801 mule: log docs-native-realm-name-pattern
2ddbf7f mule: docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
b13297f mule: log docs-native-phoenix-newevent-payload
a16b5fe mule: docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
30bb22b mule: log docs-native-metrics-content-negotiation
3962921 mule: docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
336d37e mule: log docs-native-hk-version-description
cce1376 mule: docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.

$ tail -40 .mule/todo.md && echo "=== TOTAL LINES ===" && wc -l .mule/todo.md
- [x] docs-sync-pairing-info-version-example [auto]: fix the `getDeviceInfo` response example in docs/api/astarte_pairing_api.yaml, which shows `version: "0.1.0-test"` (yaml:252) — a value no build can produce. The field is filled from `pairing.Config.Version` (internal/pairing/service.go:302), wired from `main.version` (cmd/astrate/main.go:151), whose default is `"0.1.0-dev"` (main.go:53) and which is otherwise set only via `-ldflags "-X main.version=..."`. Use `"0.1.0-dev"`, matching the four compat-version examples docs-sync-native-version-value already corrected to that (docs/api/astrate_native_api.yaml:203, 223, 243, 263). Do NOT document the service's own fallback `DefaultVersion = "0.1.0-astrate"` (service.go:86-88): it is unreachable from `main`, which always passes a non-empty value. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-deviceid-base64url [auto]: the pairing spec calls the device identifier "base64-encoded 128-bit" in two places — the `DeviceID` path parameter (docs/api/astarte_pairing_api.yaml:360) and `RegisterRequest.hw_id` (yaml:373) — but `deviceid.Parse` accepts only the 22-character **unpadded base64url** form, in strict mode, and rejects the standard alphabet's `+` and `/` and any `=` padding with an error wrapping `ErrInvalid` (pkg/deviceid/deviceid.go:26, 37, 39-60; upstream parity cited at 36). Measured: with the same `encoding` the parser uses, `dT6hS2W9TT6LEnP25ks+lg` and `.../lg` fail at byte 19 and the padded 23-character form fails at byte 21, so a client generated from the current spec can send a well-formed standard-base64 ID and get 422 `{"errors":{"hw_id":["is not a valid base64 encoded 128 bits id"]}}` (internal/pairing/http.go:307-309). State the exact encoding and the 22-character length in both descriptions and add a `pattern`/`minLength`/`maxLength` to `hw_id`. The existing `hw_id: dT6hS2W9TT6LEnP25ks_lg` example (yaml:123) is valid and stays. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-initial-payload-format-enum [auto]: constrain `RegisterRequest.initial_payload_format` in docs/api/astarte_pairing_api.yaml (yaml:374-376), an unconstrained string described only as "Optional Astrate extension specifying the device's initial payload format" — so a generated client can send any string, while `Service.Register` accepts only `"bson"` or `"json"` and rejects everything else with 422 `{"errors":{"initial_payload_format":["is invalid"]}}` (internal/pairing/service.go:179-181 -> internal/pairing/http.go:310-312; the allowed set is stated at service.go:173). Add `enum: [bson, json]` plus a note that absent/empty means unset, and add the three missing `ValidationError` examples while you are in the component: `initial_payload_format: ["is invalid"]` (http.go:310-312), `csr: ["can't be blank"]` (http.go:173-177) and `client_crt: ["can't be blank"]` (http.go:253-257) — the component currently shows only the four hw_id/error_name rows (yaml:594-618) and the `csr`/`client_crt` blank-field cases are what a device actually hits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-pairing-unregister-description [auto]: fix the `unregisterDevice` description in docs/api/astarte_pairing_api.yaml (yaml:153), "Removes a device from the realm. Returns 204 on success.", which reads as a device-and-data deletion. The row and every byte of its stored data survive: `store.UnregisterDevice` only clears the credential trail (`credentials_secret_hash = ''`, `first_credentials_request = NULL`, `cert_serial = NULL`, `cert_aki = NULL`) and flips `status` back to `'registered'` (internal/store/devices.go:93-110), and the service docstring says so outright ("its data is kept", internal/pairing/service.go:227-229). State the real semantics: the device becomes registrable again, its interfaces, datastream data and group memberships are retained, and a second DELETE is a 204 no-op (the row is still there, so RowsAffected is 1, devices.go:106-108) — not the documented 404. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [ ] deviceid-strict-trailing-bits [legion] [auto]: decide whether `deviceid.Parse` rejecting non-canonical trailing bits is a deviation worth recording — `encoding` is `base64.RawURLEncoding.Strict()` (pkg/deviceid/deviceid.go:37), so of the 64 strings that share the 22-character unpadded base64url alphabet and length, only the 16 whose last character is one of `A`, `Q`, `g`, `w` parse: measured, `dT6hS2W9TT6LEnP25ks_lx` is refused with `ErrInvalid` while `..._lA` is accepted, and the pairing spec's `hw_id` pattern `^[A-Za-z0-9_-]{22}$` (added by docs-sync-pairing-deviceid-base64url) deliberately does not encode the rule because a pattern is a client-side contract and upstream's leniency is unmeasured. The docstring justifies strictness as byte-identical round-tripping, which is a reason, not a parity claim, and docs/COMPATIBILITY.md records no entry for it. Measure upstream (Elixir `Base.url_decode64!(padding: false)`, which does not check those bits) on the Legion Go, then either record the divergence in docs/COMPATIBILITY.md or relax `encoding` to `base64.RawURLEncoding` with the canonical-form check moved to `String`.
- [x] fdo-rc6-scope-delta-for-giulio [auto]: upstream v1.4.0-rc.6 (2026-09-30, 27 commits, 154 files, almost all FDO) grows the onboarding surface issue #78 is scoped against, and none of it is in #78's body (written against upstream 1.3 / rc.5): new `POST /pairing/v1/{realm}/fdo/ownership_vouchers/{guid}/to0` re-runs TO0 and returns the refreshed rendezvous expiry (`{"data":{"expiry":…}}`, ownership_voucher.ex:63-73; 409 `device_already_onboarded` via fallback_controller.ex:120-125); voucher deletion, which revokes the rendezvous registration first and *keeps* the voucher if that fails (`:rendezvous_revocation_failed` → 500) and refuses another realm's voucher through `ensure_voucher_in_realm` (ownership_voucher.ex:40-60 — now expressible because the voucher row carries its realm); `device_id` accepted on voucher upload; the device registered immediately on upload rather than mid-TO2 (`add_unconfirmed_credentials`, device.ex:63-85); single-entry x5chain certs accepted; invalid vouchers no longer stored; vouchers and to2-sessions moved from the per-realm keyspace to the global `astarte` one with `fdo_guid` added to the device row (housekeeping/realms/queries.ex:784, 1129-1133); and RM device deletion now cascading to the voucher (device_removal/core.ex). Write the delta into `.mule/for-giulio.md` as a scope update for #78 and say which of these a client could actually observe — do not edit or comment on the issue (GitHub state is read-only for the mule) and do not implement anything.
- [x] compat-note-v1.4-rc.6 [auto]: record in `.mule/for-giulio.md` that v1.4.0-rc.6 is wire-inert for Astrate, as a delta on the open v1.3.2/v1.3.4 wording proposal rather than a re-derivation. Measured, the whole non-FDO delta is inert: DUP RPC-server reliability is OTP supervisor/GC tuning; the Dashboard's new `unknown_status` is computed client-side from `last_connection && !last_disconnection && !isConnected` (models/Device/index.ts:213-215) and both timestamps are already on Astrate's device shape (internal/appengine/service.go:118), so a device whose disconnection was never recorded reads as "Unknown status" rather than "Disconnected"; the `on_device_registered` dispatch moved from `Agent.register_device` to `Engine.register_device` (pairing/engine.ex:157-165, 186-193) so it also covers the FDO path, but the Agent path is byte-for-byte the same rule and Astrate already fires it on every successful registration (internal/pairing/service.go:221-223); the `hw_id` validation moved into `astarte_fdo`'s `LoadRequest` with the error text unchanged ("is not a valid base64 encoded 128 bits id"); and `NonNegativeIntegerOrUnsetType` changed only by a coveralls pragma. So the proposal's "v1.4.0 is still rc.5-only" reference becomes rc.6, and its "FDO authentication (pairing, disabled by default)" line understates the 1.4 line, which ships a full ownership-voucher surface Astrate does not emulate (#78, parked). Do not edit docs/COMPATIBILITY.md and do not restate the standing proposal.

- [x] container-stop-deadline [auto]: make `cliInstance.Stop` (internal/flow/blocks/container/docker.go:167-175) derive its 15s bound whenever the passed context carries **no deadline**, not only when it is nil — both cleanup callers pass `context.Background()` (docker.go:122, the "Best-effort cleanup so we do not leave orphans on mapping failure" path, and block.go:121, the not-ready cleanup in `New`), so a wedged docker daemon leaves an unbounded `docker rm -f` inside `exec.CommandContext` holding flow instantiation, while `Block.Stop` (block.go:310-311) does the same job correctly with an explicit 15s timeout. Measured with the package's own `Run` hook: the `rm` invocation arrives with `ctx.Deadline()` unset on both paths. Add a case asserting `inst.Stop(context.Background())` reaches the injected `Run` with a deadline set (it does not today) and that the hostPort-failure path still issues `rm -f <id>` (docker.go:120-124, never executed by any test today). Container-free.
- [!] container-timeout-bounds [auto]: bound `timeout_ms` and `ready_timeout_ms` in `parseConfig` (internal/flow/blocks/container/block.go:205-219) — both are only checked `> 0` and then multiplied unchecked (`time.Duration(n) * time.Millisecond`), which wraps: measured `ready_timeout_ms: 18446744073709` gives `-551.616µs`, so `New` fails with "container: not ready after wait" without ever issuing a probe (the wrapped negative deadline fails the loop guard at httpbridge.go:64), while `ready_timeout_ms: 9223372036854` gives 2562047h47m16s and `New` never returns (still polling `/healthz` after 3s in the probe). The two consumers of a non-positive duration also disagree — `Bridge.timeout()` (httpbridge.go:39-44) treats it as "use the 5s default" while net/http treats a non-positive `Client.Timeout` as no deadline at all (measured: a client with `Timeout: -551µs` returned err=nil after a 3s handler). Reject a value that would overflow or exceed a documented maximum, naming the key and the value in the error, and add table cases asserting `New` refuses both (it accepts them today). Nothing bounds these on the wire either: internal/flow/blocks/schema.go:19 advertises both as bare `"type":"integer"`. Container-free. — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [x] container-parseconfig-rules-test [auto]: pin the four `parseConfig` rules that no test touches (internal/flow/blocks/container/block.go:188-219 — all `count=0` in `go test -coverprofile`, which is 71.2% for the package): `config` must be a JSON object (192-193), `port` an integer in 1–65535 including the 0 / 65536 / non-integer rejections (197-202), and `timeout_ms` and `ready_timeout_ms` positive (205-210, 215-217) — beside the existing `TestConstructor_RequiresImage` (block_test.go:99), which is the only config case the suite has. The port rule is the one with teeth: it is what keeps `127.0.0.1::99999` out of the `docker run` argv (docker.go:93). Container-free, pure.
- [!] container-response-cap-test [auto]: pin `Bridge.RoundTrip`'s 1 MiB response cap (internal/flow/blocks/container/httpbridge.go:133-140 — `io.LimitReader(resp.Body, maxBody+1)` plus the explicit rejection; lines 136, 138.16 and 138.40 are all uncovered), the only bound between a tenant-supplied container and process memory on the message path: assert a handler answering exactly `maxBody+1` bytes is rejected with the "exceeds" error and that exactly `maxBody` bytes still round-trip, driving a `Bridge{MaxBodyBytes: N}` so the fixture stays small. Also worth one line in the same change: the cap is absent from the operator-facing catalog entry (internal/flow/blocks/info.go:116) and from the schema, so an author whose container emits 2 MiB learns about it only from a runtime error. Container-free. — BLOCKED: lint failed: internal/flow/blocks/container/block_test.go:500:2: redefines-builtin-id: redefinition of the built-in function cap (rev
- [x] docs-sync-rm-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astrate_realm_management_api.yaml — the `BadRequest` component shows `detail: Bad Request` (yaml:1265) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1348), but the wire emits the lower-r/lower-s forms the constants freeze: `"Bad request"` (pkg/astarteapi/envelope.go:37) and `"Internal server error"` (envelope.go:58), both pinned by the golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json. This is the realm-management half of the pair docs-sync-pairing-error-example-capitalisation fixed there and explicitly deferred with "leave those for a run of their surface". The other five detail examples on this file are already correct — change only those two. The identical wrong pair still sits in the appengine, housekeeping and native specs; leave those for their surfaces. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs is not installed on the Pi — `pip install -r docs/requirements.txt` first, or at minimum confirm the file still parses as YAML).
- [!] docs-sync-rm-legacy-alias-fields [auto]: add the three accepted upstream legacy alias fields to `InterfaceDefinition` and `InterfaceMapping` in docs/api/astarte_realm_management_api.yaml — both declare `additionalProperties: false` (yaml:935, 972) yet omit interface-level `quality` (alias of `ownership`) and `aggregate` (alias of `aggregation`) and mapping-level `path` (alias of `endpoint`), all three of which the parser accepts and canonicalises: decoded into dedicated fields (pkg/interfaceschema/parse.go:116-122, 136-138), normalised before any rule runs (parse.go:212-236), and re-encoded canonically by `ParseInterfaceCanonical` so the store never sees them (parse.go:165-170, 307-311). A spec-conformant client — or any generated client with this schema baked in — therefore rejects a valid upstream-shaped install/update with a spurious 422 while the service would have accepted and stored it as `ownership`/`aggregation`/`endpoint`. Add the three as documented deprecated aliases (enums mirroring ownership/aggregation, `path` a plain string) or relax `additionalProperties`, and record the three rules the code enforces: `ownership` + `quality` is a violation (parse.go:218) and `aggregation` + `aggregate` is a violation (parse.go:227), while `endpoint` + `path` silently keeps `endpoint` (parse.go:307-311). None of the three may appear in a `required` list. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-rm-deviceid-param [auto]: constrain the `DeviceID` path parameter in docs/api/astrate_realm_management_api.yaml (yaml:887-893), an unconstrained string described only as "The device hardware ID." — `Service.DeleteDevice` runs the path value through `deviceid.Parse` (internal/realm/service.go:118-121), which accepts exactly 22 characters of unpadded base64url in strict mode and rejects padding, the standard alphabet's `+` and `/`, and non-canonical trailing bits (pkg/deviceid/deviceid.go:24-26, 33-37, 42-60); every rejection is folded into `store.ErrNotFound` (service.go:119-121) and answered `404 {"errors":{"detail":"Device not found"}}` (internal/realm/http.go:107-109), so a client that sends a well-formed standard-base64 or 23-character padded id is told the device does not exist rather than that its id is malformed. State the encoding and the 22-character length and add `pattern`/`minLength`/`maxLength`, exactly as docs-sync-pairing-deviceid-base64url did for the pairing twin. Do NOT encode the strict trailing-bits rule in the pattern: whether it is a deviation worth recording is still open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_realm_management_api.yaml — `data` is declared as a bare `type: object` with no properties (yaml:250-261) although the handler parses it with the same `interfaceschema.ParseInterfaceCanonical` as `installInterface` (internal/realm/service.go:196 vs 161) and then checks the body's `interface_name`/`version_major` against the URL's `{name}`/`{major}` (service.go:203-208), so the update body — the one carrying every interface rule — is the only interface body in the spec with nothing to validate against. `$ref` `#/components/schemas/InterfaceDefinition` as `installInterface` does (yaml:115-117), and say in the operation description that the body must name the same interface and major as the URL (disagreement is 409, already documented) and that `version_minor` must increase (409). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [ ] container-orphan-boot-cleanup [legion] [auto]: implement the deferred B5 sweep (docs/handoff/flow-design-b-container-block-2026-07-29.md:199, 239) — `CLIRunner.Start` labels every container `astrate.flow=1` / `astrate.realm` / `astrate.flow_name` (internal/flow/blocks/container/block.go:84-93) precisely so they are findable, passes no `--rm` and no restart policy (both explicitly out of PoC scope, handoff §6), and nothing in the tree ever runs `docker ps --filter label=astrate.flow=1` — `rg 'docker (ps|rm|run)'` hits only this package, the example README and the handoff doc — so every astrate crash leaves the previous run's containers up holding loopback-published ports and RAM, forever. Add a best-effort `docker ps` → `docker rm -f` sweep at boot (cmd/astrate/main.go, beside the existing shutdown stages) and prove it on the Legion Go by leaving a labelled container behind, killing astrate mid-flow, restarting, and asserting the orphan is gone while the freshly started one survives. Needs Docker — the Pi has none. If you would rather keep B5 deferred, drop this line; the handoff marks it Optional.
- [x] docs-sync-ae-error-example-capitalisation [auto]: fix the two wrong error-detail examples in docs/api/astarte_appengine_api.yaml — `BadRequest` shows `detail: Bad Request` (yaml:1732) and `InternalServerError` shows `detail: Internal Server Error` (yaml:1792), while the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58 via `WriteBadRequest`/`WriteInternalServerError`, golden fixtures pkg/astarteapi/testdata/error_bad_request.json and error_internal.json); every 400/500 on that surface goes through the same astarteapi writers. The other detail examples on the file are already correct — change only those two — and add the docs-half assertion `TestAppEngineErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go, which pins the pairing and realm-management halves. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-error-example-capitalisation [auto]: fix the three wrong error-detail examples in docs/api/astarte_housekeeping_api.yaml — the `ErrorDetail` schema-level example (yaml:361-362) and the `BadRequest` response (yaml:385) show `detail: Bad Request`, and `InternalServerError` (yaml:477) shows `detail: Internal Server Error`, but the wire emits the frozen lower-r/lower-s constants `"Bad request"` and `"Internal server error"` (pkg/astarteapi/envelope.go:37, 58, via WriteBadRequest/WriteInternalServerError, pinned by pkg/astarteapi/testdata/error_bad_request.json and error_internal.json). While you are in `ErrorDetail`: its `example` sits at the schema level and shows a bare `{detail: ...}` where the envelope the wire sends is `{"errors":{"detail":...}}` — move it under `errors` (or drop it, as the realm-management and pairing `ErrorDetail` schemas carry no example) so a generated client is not shown the wrong shape. Add `TestHousekeepingErrorDetailExamples` beside the pairing and realm-management halves in internal/swagger/swagger_test.go. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-native-error-example-capitalisation [auto]: fix the one wrong error-detail example in docs/api/astrate_native_api.yaml — `InternalServerError` shows `detail: Internal Server Error` (yaml:466) while the wire emits the frozen `"Internal server error"` (pkg/astarteapi/envelope.go:58 via WriteInternalServerError, golden fixture pkg/astarteapi/testdata/error_internal.json); `Unauthorized` (yaml:446) and `Forbidden` (yaml:456) on the same spec are already correct, and this spec has no BadRequest component, so it is a one-line change. Add `TestNativeErrorDetailExamples` beside TestPairingErrorDetailExamples / TestRealmManagementErrorDetailExamples in internal/swagger/swagger_test.go. With that, no spec carries a reconstructed capitalisation any more. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-error-detail-examples-split [auto]: first tick-sized half of the blocked docs-sync-hk-error-example-capitalisation line, and the only spec still carrying a reconstructed error-detail capitalisation now that the native one is fixed (measured: `detail: Bad Request` at docs/api/astarte_housekeeping_api.yaml:362 and :385, `detail: Internal Server Error` at :477; `detail: Not Found` at :415 is the canonical `astarteapi.DetailNotFound` and is already correct). Fix only the two responses — `BadRequest` to `detail: Bad request` and `InternalServerError` to `detail: Internal server error` (pkg/astarteapi/envelope.go:37, 58) — and add `TestHousekeepingErrorDetailExamples` beside `TestNativeErrorDetailExamples` in internal/swagger/swagger_test.go, pinning those two plus the already-correct `Unauthorized`/`Forbidden`/`NotFound`/`DeviceNotFound` responses. Leave the `ErrorDetail` schema-level example (yaml:361-362, a bare `{detail: ...}` sitting under `errors:` where the wire sends `{"errors":{"detail":...}}`) to a second line: that shape half is what made the original line too big. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] device-empty-cache-received-trigger [auto]: upstream master 958bb9fe (#2232, 2026-10-01, **unreleased** — in no tag; newest is still v1.4.0-rc.6) implements `device_empty_cache_received` end to end — `SimpleEvent` field 23, an empty `DeviceEmptyCacheReceivedEvent` whose encoder emits exactly `{"type":"device_empty_cache_received"}` (astarte_core encoder.ex:54-62), fired by DUP at `control_handler.ex:161` *inside* the `with` that runs `send_control_consumer_properties` → `resend_all_properties` → `set_pending_empty_cache`, so it fires only when all three succeed and never on failure, and the commit deletes the `# TODO: implement on_empty_cache_received` from `astarte_events/lib/astarte_events/triggers/core.ex`. Astrate already accepts the trigger but never evaluates it: `OnDeviceEmptyCacheReceived` sits in `deviceOns` with value `false` (internal/engine/triggers/match.go:123, the `evaluated` flag read at match.go:406) and the comment's stated reason ("upstream defines no SimpleEvent variant for it") is now false. Add the event body constructor to internal/engine/triggers/events.go (the payload is `{"type":"device_empty_cache_received"}` and nothing else — no device id, no metadata), flip the `deviceOns` value to true, and fire it from `handleEmptyCache` (internal/engine/control.go:58) with the same idiom `handleLifecycle` uses for `device_connected` (`e.fireDevice` + `stream.Publish`, internal/engine/engine.go:346-353), at the point where `m.Ack()` runs — i.e. after `resendServerProperties` and `sendConsumerProperties` have both succeeded inside `retryStore`. Add a T1 test that a trigger on it fires exactly once per accepted emptyCache and does not fire when the resend fails. — BLOCKED: wrote nothing
- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
- [x] docs-sync-hk-realm-name-pattern [auto]: constrain `realm_name` in docs/api/astarte_housekeeping_api.yaml — the realm name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7), and a non-blank name that fails it is rejected by the INSERT as a check violation → `store.ErrInvalidRealmName` (internal/store/realms.go:67-70) → `ErrValidation: realm_name is invalid` (internal/housekeeping/service.go:178-181), answered 422 with the flat `{"errors":{"detail":"realm_name is invalid"}}` body once `validationDetail` strips the prefix (http.go:240-241, 252-258); but `RealmCreate.realm_name` (yaml:298-300) and the `RealmName` path parameter (yaml:231-237) are bare `type: string` with no `pattern`, and the POST 422 `ValidationError` component shows only `realm_name can't be blank` (yaml:459-467) — the message for the *other* rejection, the blank one (service.go:137-139). Measured against the same regex the CHECK uses: `TEST`, `test_realm`, `test-realm`, `1test`, `tëst` and `test realm` are all refused. Add `pattern: '^[a-z][a-z0-9]*$'` to both, and say in the POST 422 story that the two messages partition the cases (empty → can't be blank, non-empty but off-pattern → is invalid). Record the read-path asymmetry in the `RealmName` description too: on GET/PATCH/DELETE an off-pattern name is a plain 404, because `GetRealmByName` is a bare `WHERE name = $1` with no format check (internal/store/realms.go:91-101) → `store.ErrNotFound` → `astarteapi.WriteNotFound` — the same "told the resource does not exist rather than that its name is malformed" hole docs-sync-rm-deviceid-param closed on the realm-management twin. Do NOT describe the name as validated before the CA is minted: `ca.Generate` runs first (service.go:161) and x509 does not validate a CN, so the 422 comes out of the database, after an ECDSA keygen and a seal that are then thrown away. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-sync-hk-errordetail-schema-example [auto]: fix the schema-level example inside `ErrorDetail` in docs/api/astarte_housekeeping_api.yaml (yaml:361-362) — the second half docs-sync-hk-error-detail-examples-split explicitly deferred ("its shape question belongs to its own line"). It reads `detail: Bad Request`, the reconstructed title-case form that `TestHousekeepingErrorDetailExamples` (internal/swagger/swagger_test.go:687) now pins against everywhere else on this spec, while the wire emits the frozen `astarteapi.DetailBadRequest` = `"Bad request"` (pkg/astarteapi/envelope.go:37, via `WriteBadRequest` on every 400 — internal/housekeeping/http.go:65, 166, 171, 183). Because the example sits on the generic `ErrorDetail` schema, which all eight error responses `$ref` (yaml:377-477), it is the example a generator falls back to for any response whose own example it cannot resolve — so it teaches `"Bad Request"` for a 404 and a 500 too. Change it to `Bad request`, or drop it as the realm-management and pairing `ErrorDetail` schemas do (neither carries one). Then extend `TestHousekeepingErrorDetailExamples`, whose comment currently says the schema-level example is "deliberately not pinned" (swagger_test.go:684-686), to pin it like the responses. One thing not to "fix": the example is nested one level in — it is the `errors` sub-schema's example, not a bare envelope — so do not re-wrap it in another `errors:`. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-sync-hk-patch-422-field-error-examples [auto]: add the three missing PATCH 422 bodies to the `ValidationErrors` response component in docs/api/astarte_housekeeping_api.yaml (yaml:427-436, referenced once at yaml:190) — it carries a single example, `error_name: [invalid_update_parameters]` (internal/housekeeping/http.go:176-177), but `patchRealm` emits three more field-error envelopes: `jwt_public_key_pem: ["can't be blank"]` (http.go:188-189), `device_registration_limit: ["is invalid"]` (http.go:193-194) and `datastream_maximum_storage_retention: ["is invalid"]` (http.go:198-199). Add them as further `examples` on the media type (OpenAPI 3.0.3, yaml:1) or split the component per case. Also record that the *messages* differ from POST for the same conditions: create rejects a negative limit/retention with the flat `ErrValidation` detail `device_registration_limit must be non-negative` / `datastream_maximum_storage_retention must be non-negative` (internal/housekeeping/service.go:143-148, 422 `ErrorDetail` shape), while PATCH answers `is invalid` in the FieldErrors shape — and that the service's own three `ErrValidation` branches (service.go:205-213) are unreachable from REST, because `patchRealm` pre-checks the identical conditions and answers first. Say in the `RealmPatch` field descriptions that a negative value is a 422 rather than silently ignored. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [!] docs-sync-hk-wrong-type-field-400 [auto]: widen the `BadRequest` component description in docs/api/astarte_housekeeping_api.yaml (yaml:378, "malformed or missing data envelope") — measured today, a perfectly well-formed `data` envelope carrying a wrong-typed field answers **400**, not the documented 422. `astarteapi.DecodeData` into `map[string]json.RawMessage` (internal/housekeeping/http.go:170) cannot fail on a type mismatch, so the failure surfaces on the second decode into the typed body (http.go:182-184) and is answered `WriteBadRequest`. Probe output: `{"data":{"device_registration_limit":"100"}}` → `json: cannot unmarshal string into Go struct field patchBody.device_registration_limit of type int32`; `{"data":{"jwt_public_key_pem":123}}` → the string twin; `{"data":{"device_registration_limit":1.5}}` → `cannot unmarshal number 1.5`; `{"data":[]}` → `cannot unmarshal array into Go value of type patchBody`. Same on POST, which decodes straight into `realmBody` (http.go:63-67). This is unreachable from a spec-conformant client — `RealmCreate`/`RealmPatch` already type both fields — so document it as a one-line extension of the 400 description on both operations rather than inventing a 422, and do NOT re-shape the `BadRequest` example, whose `detail: Bad request` is the frozen constant `TestHousekeepingErrorDetailExamples` pins. One distinction worth stating in the same description: a wrong-typed field is *not* the `invalid_update_parameters` case — unknown keys are caught earlier by the `patchAllowedFields` sweep (http.go:174-180, 422 FieldErrors) on PATCH only, and silently dropped on POST. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-sync-hk-realm-name-response-schemas [auto]: the write side of the realm name is now constrained in docs/api/astarte_housekeeping_api.yaml — `components.parameters.RealmName` and `RealmCreate.realm_name` both carry `pattern: '^[a-z][a-z0-9]*$'` and `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go) holds both against the realm name CHECK extracted from migrations/000002_metadata.up.sql — but the read side is still bare `type: string`: `Realm.realm_name` (components.schemas.Realm) and `DataEnvelopeRealmNames.items`. Decide whether a response field earns the same pattern (a generated client learns nothing new from it, but a reader diffing create against get sees one constrained and one not) and, if it does, add it to both and extend that test to hold the response-side pattern against the same CHECK. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] store-register-inhibit-preserve: RegisterDevice must not overwrite status to 'registered' when a device was inhibited before first credentials request (preserve 'inhibited' state). File: internal/store/devices.go — BLOCKED: gates failed
- [x] store-aliasvalues-self-exclusion-test: Add test for AliasValuesTaken to assert self-exclusion (same device's own alias values do not count as taken). File: internal/store/devices_test.go
- [!] store-latestindividual-empty-errnotfound: Add store-level test asserting LatestIndividual returns ErrNotFound for empty series and returns newest on multiple samples. File: internal/store/datastreams_test.go — BLOCKED: TIMEOUT after 1201s — task too big, split it
- [x] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.go
- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
=== TOTAL LINES ===
281 .mule/todo.md

$ echo "=== reports dir ===" && ls -la .mule/reports/ && echo "=== for-giulio tail ===" && tail -30 .mule/for-giulio.md
=== reports dir ===
total 56
drwxr-xr-x 3 root root  4096 Sep  6 19:40 .
drwxr-xr-x 7 root root  4096 Oct  5 20:14 ..
-rw-r--r-- 1 root root  2279 Sep  6 19:40 hygiene.md
drwxr-xr-x 2 root root 40960 Oct  5 20:14 log
=== for-giulio tail ===

**3. Two `[housekeeping]` config keys exist; `configuration-reference.md` does not list them.**
The table at `docs/site/configuration-reference.md:79-84` has exactly two rows and the section
heading calls the block "Instance-admin keys" — but `HousekeepingConfig`
(`internal/config/config.go:109-125`) has four fields:

- `default_datastream_maximum_storage_retention` (`*int64`, env
  `ASTRATE_HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`, validated as a
  non-negative integer at `config.go:286-291`) — injected into a realm at creation when the
  caller omits the field (`internal/housekeeping/service.go:149-150`). An operator who wants it
  has no way to learn it exists.
- `realm_deletion_disabled` (bool, env `ASTRATE_HOUSEKEEPING_REALM_DELETION_DISABLED`, a
  fail-loud boolean accepting only `1|true|TRUE|True|0|false|""` and refusing to load on anything
  else, `config.go:296-306`) — this is the switch behind the spec's documented
  `405 "Realm deletion disabled"` on `DELETE /housekeeping/v1/realms/{realm}`, so it changes a
  status code an API client sees.

Both belong in that table next to the two key lists, and the heading wants widening.

Checked and *not* escalated, for the record. The config-key sweep
(`rg -o '\bASTRATE_[A-Z_]+' -N internal/ | sort -u`) found nothing else stale: the two
`jwt_public_key*` keys the page quotes at lines 69-71 are real (`config.go:112-113`), and the
`master_key_file` row's `ASTRATE_MASTER_KEY`/`ASTRATE_MASTER_KEY_FILE` fallbacks are real too
(`internal/store/crypto.go:24-27`). And the one thing that looked like a live bug is not one:
`view()` returns only `JWTPublicKeysPEM[0]` (`internal/housekeeping/service.go:277-281`) while
the realm-management surface joins *all* keys (`internal/realm/service.go:627`) — but every
writer stores a single-element array (create at `internal/store/realms.go:54-61`, PATCH at
`internal/store/realms.go:165`, `putAuthConfig` at `internal/realm/service.go:636`), so a
two-key realm is unreachable and there is nothing to document. Worth noting because it reads
like a bug on first inspection.

→ Read .mule/reports/hygiene.md

$ rg -n 'hygiene' .mule/log.md | tail -20 && echo "=== open hygiene lines in todo ===" && rg -n '^\- \[ \].*hygiene|^\- \[ \].*lint|^\- \[ \].*vuln' .mule/todo.md
=== open hygiene lines in todo ===
88:- [ ] hygiene-govulncheck [legion]: on the Legion Go, run govulncheck ./... on a fresh ~/astrate clone and file a task line for any REACHABLE advisory it reports (name the CVE and the call path, per the hygiene recipe); this box cannot build it — go install golang.org/x/vuln/cmd/govulncheck@latest is OOM-killed here (3.7GB, 4 cores), so the recipe's highest-priority check never runs unless done there.

$ which go golangci-lint govulncheck 2>&1; echo "--- go version ---"; go version; echo "--- gopath bin ---"; ls $(go env GOPATH)/bin 2>/dev/null; echo "--- golangci config ---"; cat .golangci.yml
/usr/bin/go
--- go version ---
go version go1.26.7 linux/arm64
--- gopath bin ---
goimports
golangci-lint
govulncheck
--- golangci config ---
# golangci-lint v2 configuration (pinned v2.12.2 — see Makefile GOLANGCI_LINT_VERSION).
# Linter set per docs/ROADMAP.md §1.2 (file 0.4): govet, staticcheck, errcheck, gosec, revive.
version: "2"

run:
  build-tags:
    # Lint everything, including T2/T3 suites that are tag-gated out of `go test ./...`.
    - integration
    - e2e

linters:
  default: none
  enable:
    - errcheck
    - gosec
    - govet
    - revive
    - staticcheck
  settings:
    gosec:
      excludes:
        # G104 duplicates errcheck.
        - G104
  exclusions:
    generated: lax
    rules:
      # Tests: hardening lints and strict error checking add noise, not safety.
      - path: _test\.go
        linters:
          - gosec
      - path: _test\.go
        linters:
          - errcheck
      # Test utilities are test-scoped code shipped outside _test.go files.
      - path: internal/testutil/
        linters:
          - gosec

formatters:
  enable:
    - gofmt
    - goimports
  settings:
    goimports:
      local-prefixes:
        - github.com/astrate-platform/astrate

$ ssh -o ConnectTimeout=5 -o BatchMode=yes legion 'echo ok' 2>&1 | tail -3; echo "exit=$?"
ok
exit=0

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ 2>/dev/null | head -40
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):

$ go list -m -u all 2>/dev/null | rg '\[' | head -30; echo "=== direct requires ==="; rg -n '^\s+\S+ v' go.mod | head -40
cel.dev/expr v0.24.0 [v0.25.3]
cloud.google.com/go v0.121.6 [v0.123.0]
cloud.google.com/go/auth v0.16.4 [v0.24.0]
cloud.google.com/go/auth/oauth2adapt v0.2.8 [v0.3.0]
cloud.google.com/go/compute/metadata v0.8.0 [v0.10.0]
cloud.google.com/go/iam v1.5.2 [v1.14.0]
cloud.google.com/go/longrunning v0.6.7 [v1.3.0]
cloud.google.com/go/monitoring v1.24.2 [v1.31.0]
cloud.google.com/go/spanner v1.85.0 [v1.95.1]
cloud.google.com/go/storage v1.56.0 [v1.69.0]
github.com/99designs/keyring v1.2.1 [v1.2.2]
github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0 [v1.23.2]
github.com/Azure/azure-sdk-for-go/sdk/internal v1.1.2 [v1.13.0]
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0 [v1.8.2]
github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c [v0.0.0-20260917205352-e937bb47801a]
github.com/Azure/go-autorest/autorest/adal v0.9.16 (retracted) [v0.9.24] (deprecated)
github.com/Azure/go-autorest/autorest/date v0.3.0 [v0.3.1] (deprecated)
github.com/Azure/go-autorest/logger v0.2.1 [v0.2.2] (deprecated)
github.com/Azure/go-autorest/tracing v0.6.0 [v0.6.1] (deprecated)
github.com/ClickHouse/clickhouse-go v1.4.3 [v1.5.4]
github.com/DataDog/zstd v1.4.5 [v1.5.7]
github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp v1.5.3 [v1.6.0]
github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp v1.27.0 [v1.38.0] (deprecated)
github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric v0.53.0 [v0.62.0] (deprecated)
github.com/GoogleCloudPlatform/opentelemetry-operations-go/internal/resourcemapping v0.53.0 [v0.62.0]
github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 [v0.0.0-20240927000941-0f3dac36c52b]
github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a [v0.0.0-20230218143504-906a9b012302]
github.com/alicebob/miniredis/v2 v2.23.0 [v2.39.0]
github.com/andybalholm/brotli v1.0.4 [v1.2.6]
github.com/apache/thrift v0.16.0 [v0.25.0]
=== direct requires ===
8:	github.com/BurntSushi/toml v1.6.0 // TOML config (§5.1)
9:	github.com/coder/websocket v1.8.14 // live stream socket (M7)
10:	github.com/eclipse/paho.mqtt.golang v1.5.1 // test-only MQTT client (same family as official Go SDK)
11:	github.com/golang-jwt/jwt/v5 v5.3.1 // JWT verification, alg allowlist (§4.2)
12:	github.com/golang-migrate/migrate/v4 v4.19.1 // migrations via source/iofs + go:embed (§1.3)
13:	github.com/hashicorp/golang-lru/v2 v2.0.7 // verified-token LRU cache (§4.2)
14:	github.com/jackc/pgx/v5 v5.10.0 // Postgres/TimescaleDB driver + pgxpool (§1.3)
15:	github.com/mochi-mqtt/server/v2 v2.7.9 // embedded MQTT broker (§1.1)
16:	github.com/prometheus/client_golang v1.23.2 // metrics (§5.2)
17:	github.com/testcontainers/testcontainers-go v0.43.0 // T2 integration harness (§5.4)
18:	github.com/testcontainers/testcontainers-go/modules/postgres v0.42.0 // Timescale container module
19:	go.etcd.io/bbolt v1.4.3 // broker session store (§3.1; pinned over pebble)
20:	go.mongodb.org/mongo-driver/v2 v2.6.0 // BSON raw-document API only (§3.5.5)
21:	golang.org/x/crypto v0.53.0 // bcrypt cost 10 for credentials secrets (§4.1)
25:	github.com/cbroglie/mustache v1.4.0
26:	github.com/nats-io/nats.go v1.52.0
27:	github.com/santhosh-tekuri/jsonschema/v6 v6.0.3
28:	github.com/testcontainers/testcontainers-go/modules/nats v0.43.0
32:	dario.cat/mergo v1.0.2 // indirect
33:	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
34:	github.com/Microsoft/go-winio v0.6.2 // indirect
35:	github.com/beorn7/perks v1.0.1 // indirect
36:	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
37:	github.com/cespare/xxhash/v2 v2.3.0 // indirect
38:	github.com/containerd/errdefs v1.0.0 // indirect
39:	github.com/containerd/errdefs/pkg v0.3.0 // indirect
40:	github.com/containerd/log v0.1.0 // indirect
41:	github.com/containerd/platforms v0.2.1 // indirect
42:	github.com/cpuguy83/dockercfg v0.3.2 // indirect
43:	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
44:	github.com/distribution/reference v0.6.0 // indirect
45:	github.com/docker/go-connections v0.6.0 // indirect
46:	github.com/docker/go-units v0.5.0 // indirect
47:	github.com/ebitengine/purego v0.10.0 // indirect
48:	github.com/felixge/httpsnoop v1.0.4 // indirect
49:	github.com/go-logr/logr v1.4.3 // indirect
50:	github.com/go-logr/stdr v1.2.2 // indirect
51:	github.com/go-ole/go-ole v1.2.6 // indirect
52:	github.com/google/uuid v1.6.0 // indirect
53:	github.com/gorilla/websocket v1.5.3 // indirect

$ rg -n 'APICompatVersion' internal/observability/*.go | head; echo "=== recent milestone tags ==="; git tag --sort=-creatordate | head -10; echo "=== milestones.md head ==="; head -40 .mule/milestones.md
=== recent milestone tags ===
v0.2.0
mule-queue-archive-2026-09-04
v0.1.0
=== milestones.md head ===
# Milestones

Release-tag-gated goals. Each section names a body of work that must be true by the time
that tag is cut. `.mule/recipes/milestones.md` reads this file to find the **current
target** — the first section below not marked `DONE` — and works toward it: investigating,
filing GitHub issues (`milestone-<tag>` label, alongside `mule`) with sub-issues where the
work splits, and escalating anything that needs a design decision to `.mule/for-giulio.md`.

**This file is Giulio's.** The recipe reads it and may propose edits via
`.mule/for-giulio.md`, but never edits it directly — same rule as `docs/COMPATIBILITY.md`.
Mark a milestone `DONE` yourself once the tag is actually cut.

Order matters: milestones are worked **in order**, lowest tag first. Do not start v3.0
investigation while v2.0 has open, un-escalated gaps — say so and stop instead.

---

## v2.0 — astarte-flow feature parity

Reference: astarte_flow (https://github.com/astarte-platform/astarte_flow) — upstream's
Elixir "Flow" component. It lets ingested data get piped through a graph of processing
blocks (native blocks and containerised ones) before it lands in storage, with pipelines
described as reusable, parametrised graphs.

Scope for this milestone: whatever set of Flow's *capabilities* Astrate needs to expose the
same wire-visible behaviour and operator-facing concepts (pipelines, blocks, native vs.
containerised blocks) — not a port of the Elixir implementation. See
`.mule/recipes/astarte-upstream.md`'s rule: port the idea, restated in Go, never the code.

Status: **DONE** (2026-07-29), tagged `v0.2.0` on 2026-09-04 — the milestone names
(`v2.0`, `v3.0`) are milestone names, not release versions; the project is pre-1.0 and the
version number keeps its own line. Runtime, factory, catalog (incl. filter/map), process
wiring, and `/flow/v1` API are on `main`. **Parity audit** + **product decisions**
recorded 2026-07-29 (`docs/handoff/flow-parity-audit-2026-07-29.md`,
`docs/handoff/flow-v2-decisions-2026-07-29.md`).

**Design A + B landed on `main`** (commit `89145e6`, 2026-07-29): durable flows +
auto_restart (**#41 closed**); named multi-instance + config (**#40 closed**) —
migration `000009`, store, `${config.*}`, API, boot rehydrate. Container block
PoC→MVP (**#43 closed**) — registered in the catalog, usable inside stored

$ go list -m -u all 2>/dev/null | rg '\[' | rg -v 'cloud.google.com|google.golang.org/api|Azure|azure|grpc-ecosystem|opentelemetry' | head -40
cel.dev/expr v0.24.0 [v0.25.3]
github.com/99designs/keyring v1.2.1 [v1.2.2]
github.com/ClickHouse/clickhouse-go v1.4.3 [v1.5.4]
github.com/DataDog/zstd v1.4.5 [v1.5.7]
github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp v1.5.3 [v1.6.0]
github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 [v0.0.0-20240927000941-0f3dac36c52b]
github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a [v0.0.0-20230218143504-906a9b012302]
github.com/alicebob/miniredis/v2 v2.23.0 [v2.39.0]
github.com/andybalholm/brotli v1.0.4 [v1.2.6]
github.com/apache/thrift v0.16.0 [v0.25.0]
github.com/aws/aws-sdk-go v1.49.6 [v1.55.8] (deprecated)
github.com/aws/aws-sdk-go-v2 v1.16.16 [v1.47.1]
github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.8 [v1.7.20]
github.com/aws/aws-sdk-go-v2/credentials v1.12.20 [v1.20.6]
github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.33 [v1.23.11] (deprecated)
github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.23 [v1.5.4]
github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.17 [v2.8.4]
github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.14 [v1.5.4]
github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.9 [v1.13.19]
github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.18 [v1.11.5]
github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.17 [v1.14.4]
github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.17 [v1.20.4]
github.com/aws/aws-sdk-go-v2/service/s3 v1.27.11 [v1.114.0]
github.com/aws/smithy-go v1.13.3 [v1.28.2]
github.com/cbroglie/mustache v1.4.0 [v1.4.2]
github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58 [v0.0.0-20240916140612-caecf3c00c06]
github.com/cncf/xds/go v0.0.0-20250501225837-2ac532fd4443 [v0.0.0-20260202195803-dba9d589def2]
github.com/cockroachdb/cockroach-go/v2 v2.1.1 [v2.4.3]
github.com/cockroachdb/errors v1.11.1 [v1.14.0]
github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b [v0.0.0-20241215232642-bb51bb14a506]
github.com/cockroachdb/pebble v1.1.0 [v1.1.5]
github.com/cockroachdb/redact v1.1.5 [v1.1.8]
github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 [v0.0.0-20250429170803-42689b6311bb]
github.com/coder/websocket v1.8.14 [v1.8.15]
github.com/containerd/log v0.1.0 [v0.2.0]
github.com/containerd/typeurl/v2 v2.2.0 [v2.3.0]
github.com/cznic/mathutil v0.0.0-20180504122225-ca4c9f2c1369 [v0.0.0-20181122101859-297441e03548]
github.com/danieljoos/wincred v1.1.2 [v1.2.3]
github.com/dgraph-io/badger/v4 v4.2.0 [v4.9.6]
github.com/dgraph-io/ristretto v0.1.1 [v0.2.0]

$ rg -rn 'APICompatVersion' --include='*.go' . 2>/dev/null | head; rg -n 'APICompatVersion' -g '*.go' . | head; echo "=== for-giulio dep entries ==="; rg -n 'go list -m -u|dependency|dep bump|go.mod' .mule/for-giulio.md | head -20
./cmd/astrate/main.go:437:		mw.RequireRealmAny(auth.ClaimAppEngine)(observability.VersionHandler(realm.APICompatVersion)))
./cmd/astrate/main.go:438:	mux.Handle("GET /pairing/v1/{realm}/version", observability.VersionHandler(realm.APICompatVersion))
./cmd/astrate/version_test.go:25:	want := `{"data":"` + realm.APICompatVersion + `"}`
./internal/realm/http.go:119:// feature UI on it — see APICompatVersion).
./internal/realm/http.go:121:	_ = astarteapi.WriteData(w, http.StatusOK, APICompatVersion)
./internal/realm/dashboard_compat_test.go:35:		if v != APICompatVersion {
./internal/realm/dashboard_compat_test.go:36:			t.Errorf("version = %q, want %q", v, APICompatVersion)
./internal/realm/service.go:580:// APICompatVersion is the upstream Realm Management API level Astrate
./internal/realm/service.go:588:const APICompatVersion = "1.2.2"
=== for-giulio dep entries ===
4:**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
384:- **govulncheck GO-2026-5970: reachable DoS in golang.org/x/text (infinite loop on invalid input, fixed in v0.39.0, available v0.41.0).** Astrate pins `x/text` indirect at v0.38.0 (go.mod:97) and pgx pulls it into production: `internal/store/notify.go:59` `store.Listen` → `pgx.ConnectConfig` → `unicode/norm.*`. This is the only govulncheck symbol finding that is not test-harness-only: GO-2026-6355/6354 (x/crypto/ssh deadlocked-channel DoS) and GO-2026-6253 (moby/go-archive tar path traversal) are reachable only through testcontainers in `internal/testutil/pg.go`, i.e. never in the deployed binary. `x/text` keeps API compatibility minor-to-minor and the modules Astrate exercises (`unicode/norm` via pgx, `text/language` via jsonschema) are unchanged, so this is a fix Astrate actually needs — the hygiene recipe's highest-priority category. Not a mule task (go.mod never-touch): your decision to bump ≥v0.39.0 now or fold into the next milestone-boundary sweep. Raw: https://pkg.go.dev/vuln/GO-2026-5970. (The 2026-09-04 dep sweep did not list x/text.)
408:- **Dependency sweep corrected: direct (pinned) deps DO have newer versions** — the 2026-09-02 note said the `go list -m -u` sweep showed "only version-skew on transitive deps", but that run hit the recipe's `head -20` cutoff (all cloud/azure/transitive) and never reached the directly-required modules. Full sweep, 2026-09-04. None of these is a fix this repo *needs*, so no bump is proposed — recorded for the decision. Per module (current → available; breaking change; repo use):
421:  only what carries a fix Astrate actually needs. `go.mod` stays on the never-touch list.
478:  SDK/wire compatibility, not minimum dependency count — Astrate is allowed to be a

$ timeout 1500 $(go env GOPATH)/bin/govulncheck ./... 2>&1 | tail -60; echo "EXIT=${PIPESTATUS[0]}"
=== Symbol Results ===

Vulnerability #1: GO-2026-6355
    Prevent DoS on deadlocked established channel in golang.org/x/crypto/ssh
  More info: https://pkg.go.dev/vuln/GO-2026-6355
  Module: golang.org/x/crypto
    Found in: golang.org/x/crypto@v0.53.0
    Fixed in: golang.org/x/crypto@v0.56.0
    Example traces found:
      #1: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls ssh.NewClientConn

Vulnerability #2: GO-2026-6354
    Prevent DoS on deadlocked undecided channel in golang.org/x/crypto/ssh
  More info: https://pkg.go.dev/vuln/GO-2026-6354
  Module: golang.org/x/crypto
    Found in: golang.org/x/crypto@v0.53.0
    Fixed in: golang.org/x/crypto@v0.56.0
    Example traces found:
      #1: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls ssh.NewClientConn

Vulnerability #3: GO-2026-6253
    moby/go-archive: Crafted tar archive can write outside the extraction
    directory in github.com/moby/go-archive
  More info: https://pkg.go.dev/vuln/GO-2026-6253
  Module: github.com/moby/go-archive
    Found in: github.com/moby/go-archive@v0.2.0
    Fixed in: github.com/moby/go-archive@v0.3.0
    Example traces found:
      #1: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls archive.TarWithOptions
      #2: internal/flow/pipeline.go:154:21: flow.Pipeline.MarshalJSON calls json.Marshal, which eventually calls archive.init
      #3: internal/testutil/pg.go:13:2: testutil.init calls testcontainers.init, which calls archive.init
      #4: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls compression.CompressStream
      #5: internal/flow/pipeline.go:154:21: flow.Pipeline.MarshalJSON calls json.Marshal, which eventually calls compression.init
      #6: internal/testutil/pg.go:13:2: testutil.init calls testcontainers.init, which eventually calls compression.init
      #7: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls compression.nopWriteCloser.Close
      #8: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls tarheader.FileInfoHeaderNoLookups
      #9: internal/testutil/pg.go:13:2: testutil.init calls testcontainers.init, which eventually calls tarheader.init
      #10: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls tarheader.nosysFileInfo.Gname
      #11: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls tarheader.nosysFileInfo.Sys
      #12: internal/testutil/pg.go:58:28: testutil.StartTimescale calls postgres.Run, which eventually calls tarheader.nosysFileInfo.Uname

Vulnerability #4: GO-2026-5970
    Infinite loop on invalid input in golang.org/x/text
  More info: https://pkg.go.dev/vuln/GO-2026-5970
  Module: golang.org/x/text
    Found in: golang.org/x/text@v0.38.0
    Fixed in: golang.org/x/text@v0.39.0
    Example traces found:
      #1: internal/store/notify.go:59:34: store.Listen calls pgx.ConnectConfig, which eventually calls norm.Form.Properties
      #2: internal/store/notify.go:59:34: store.Listen calls pgx.ConnectConfig, which eventually calls norm.Form.Span
      #3: internal/store/notify.go:59:34: store.Listen calls pgx.ConnectConfig, which eventually calls norm.Form.Transform

Your code is affected by 4 vulnerabilities from 3 modules.
This scan also found 2 vulnerabilities in packages you import and 3
vulnerabilities in modules you require, but your code doesn't appear to call
these vulnerabilities.
Use '-show verbose' for more details.
/bin/sh: 1: Bad substitution

$ sed -n '375,425p' .mule/for-giulio.md
stabilized and the release has not. Implementing it is the largest surface in the parity
  backlog (CBOR codec, X25519/P-256, HKDF, AES-256-GCM, a 5-state handshake machine,
  shared-secret persistence, five new error names). Your call: build against the rc now,
  re-park until v1.4.0 is a stable tag, or take only the narrow #93 fix. (Escalated again
  2026-09-05 — a prior escalation from the 2026-09-04 milestone run was lost in the queue
  rebuild.)

---

- **govulncheck GO-2026-5970: reachable DoS in golang.org/x/text (infinite loop on invalid input, fixed in v0.39.0, available v0.41.0).** Astrate pins `x/text` indirect at v0.38.0 (go.mod:97) and pgx pulls it into production: `internal/store/notify.go:59` `store.Listen` → `pgx.ConnectConfig` → `unicode/norm.*`. This is the only govulncheck symbol finding that is not test-harness-only: GO-2026-6355/6354 (x/crypto/ssh deadlocked-channel DoS) and GO-2026-6253 (moby/go-archive tar path traversal) are reachable only through testcontainers in `internal/testutil/pg.go`, i.e. never in the deployed binary. `x/text` keeps API compatibility minor-to-minor and the modules Astrate exercises (`unicode/norm` via pgx, `text/language` via jsonschema) are unchanged, so this is a fix Astrate actually needs — the hygiene recipe's highest-priority category. Not a mule task (go.mod never-touch): your decision to bump ≥v0.39.0 now or fold into the next milestone-boundary sweep. Raw: https://pkg.go.dev/vuln/GO-2026-5970. (The 2026-09-04 dep sweep did not list x/text.)

---

- **COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable, 2026-07-14; v1.4.0 is still rc-only).** Astrate's doc and `APICompatVersion` still target upstream **1.2.2** (`internal/realm/service.go:588`); v1.3.0 (2026-05-06) introduced wire-surface changes Astrate does not yet emulate, so this is a decision — adopt v1.3.2 as the compatibility target (then update the doc + bump `APICompatVersion` together, per the bump rule) or keep 1.2.2 and add a "not yet emulated" note. Wire-relevant v1.3.0 deltas (release notes): **MQTT v1 capabilities** incl. `purge_properties_compression_format` (plaintext vs zlib purge — touches the `emptyCache`/`producer,properties` contract COMPATIBILITY.md deviation 1 documents); **empty introspection now allowed**; **device registration triggers** (pairing) and **device deletion started/completed triggers** (RM — the latter two already exist as Astrate deviation 9 emits both around the synchronous delete); **FDO authentication** (pairing, disabled by default); **realm-scoped health** — upstream v1.3 added `GET /pairing/v1/{realm}/health`, which Astrate already serves (`internal/pairing/http.go:78-81`, comment already says "upstream 1.3+"), so deviation 18's wording ("which upstream 404s") is now false against 1.3 and the note should be reworded either way. Proposed doc wording (for your approval, edit to taste): in §Infrastructure differences add a sentence — *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl.

plaintext `purge_properties_compression_format`, empty-introspection allowance, device registration/deletion triggers, experimental FDO pairing auth) are not yet emulated and are out of scope until the milestone that adopts v1.3.2 as the target."* — and reword deviation 18's realm-health note from "which upstream 404s" to "added by upstream v1.3 (Astrate serves it against a 1.2.2 target; kept, matching behavior)". Raw upstream changes: [v1.3.0](https://github.com/astarte-platform/astarte/releases/tag/v1.3.0), [v1.3.2](https://github.com/astarte-platform/astarte/releases/tag/v1.3.2).

---
- **The mule's 42 straight failures were ten lint errors in its own base, not a stale branch.**
  `mule/queue` stopped taking `main` on 2026-07-27 and drifted 120 commits behind, which is
  real and is why it is being rebuilt — but it is not what blocked the work. The lint gate runs
  `golangci-lint run ./...` over the whole repo, and the branch's own `internal/flow/*` code
  carried 10 findings (1 goimports, 1 gosec, 8 revive). So every task failed the gate no matter
  how good the change was, including the four tasks queued to fix those very findings: each
  fixed one and died on the other nine. Base tests and `go vet` pass on that checkout, and
  `golangci-lint` and `govulncheck` have been installed on the Pi the whole time — the two
  things that looked broken were not. `main` is lint-clean, so rebuilding from it clears the
  deadlock. `tools/mule.sh preflight` checks the lint baseline and would have said so on day
  one; nobody ran it between 2026-08-31 and 2026-09-04. (Diagnosed 2026-09-04.)

---

- **Dependency sweep corrected: direct (pinned) deps DO have newer versions** — the 2026-09-02 note said the `go list -m -u` sweep showed "only version-skew on transitive deps", but that run hit the recipe's `head -20` cutoff (all cloud/azure/transitive) and never reached the directly-required modules. Full sweep, 2026-09-04. None of these is a fix this repo *needs*, so no bump is proposed — recorded for the decision. Per module (current → available; breaking change; repo use):
  - `github.com/coder/websocket` v1.8.14 → v1.8.15 — no breaking (patch); used in `internal/appengine/stream/ws.go`, `channels/ws.go`; worth it only for the "transmit in single frame when compression enabled" fix + read-path alloc reduction.
  - `go.etcd.io/bbolt` v1.4.3 → v1.5.0 — bbolt's semver promises no API change between patch/minor, so additive-only; used in `internal/broker/sessionstore.go`; v1.5 adds a data-file size limit and panic-recovery hardening, nothing Astrate needs.
  - `go.mongodb.org/mongo-driver/v2` v2.6.0 → v2.8.2 — the 2.8.0 breaking changes are confined to Queryable Encryption string-query options (`options.Text()`→`String()`); Astrate uses only the raw BSON API (`pkg/payload/bson.go`, `internal/engine/capabilities.go`, `bench/`) and is unaffected.
  - `github.com/nats-io/nats.go` v1.52.0 → v1.53.1 — no breaking; the headline fixes (JetStream `resetOrderedConsumer` race, KV dot-rejection) are paths Astrate does not use — `internal/engine/forward/nats.go` is core NATS publish only.
  - `github.com/prometheus/client_golang` v1.23.2 → v1.24.1 — requires Go ≥1.25 (fine, repo is 1.26.1); the breaking `LabelNames`/remote-api renames don't touch repo usage (`prometheus`/`collectors`/`promhttp` in `internal/observability/metrics.go`, flow/engine metrics); would buy `Gather()` panic-recovery and opt-in `CoalesceGather` scrape-pile-up protection.
  - `github.com/testcontainers/testcontainers-go` v0.43.0 → v0.44.0 (modules/postgres v0.42.0, modules/nats v0.43.0) — breaking in `wait.ForSQL` (callback now takes `network.Port`) and `ImageProvider` (new `PullImageWithPlatform`); Astrate's `internal/testutil/pg.go` looks unaffected but it is test-only anyway.
  - `golang.org/x/crypto` v0.53.0 → v0.56.0 — x/crypto keeps API compatibility; used only for bcrypt in `internal/auth`.
  Note (corrected 2026-09-04): `govulncheck` and `golangci-lint` **are** installed on the Pi (`/root/go/bin`, since 2026-07-28 and 2026-09-01) and `.mule/config` finds them there, so both checks were available; the sweep that produced this list simply ran without invoking govulncheck.

  **Decided 2026-09-04: no bumps.** None of the seven fixes anything this repo has, and each
  one costs a full test run to land. The standing rule instead: re-run this sweep at every
  milestone boundary (the point where `APICompatVersion` or a milestone tag moves), and bump
  only what carries a fix Astrate actually needs. `go.mod` stays on the never-touch list.

---

- **milestone 2.0 looks complete, verify and cut the tag** — all 11 `milestone-2.0` issues

$ timeout 1500 $(go env GOPATH)/bin/golangci-lint run ./... 2>&1 | tail -40
0 issues.

$ rg -n 'moby/go-archive' go.mod go.sum | head -5; echo "=== who imports it ==="; timeout 300 go mod why github.com/moby/go-archive 2>&1 | head -20; echo "=== deps of internal/flow ==="; timeout 300 go list -deps ./internal/flow/... 2>/dev/null | rg -i 'archive|testcontainer|docker' | head -20
go.mod:63:	github.com/moby/go-archive v0.2.0 // indirect
go.sum:106:github.com/moby/go-archive v0.2.0 h1:zg5QDUM2mi0JIM9fdQZWC7U8+2ZfixfTYoHL7rWUcP8=
go.sum:107:github.com/moby/go-archive v0.2.0/go.mod h1:mNeivT14o8xU+5q1YnNrkQVpK+dnNe/K6fHqnTg4qPU=
=== who imports it ===
# github.com/moby/go-archive
github.com/astrate-platform/astrate/internal/testutil
github.com/testcontainers/testcontainers-go
github.com/moby/go-archive
=== deps of internal/flow ===

$ for d in $(go list ./... 2>/dev/null); do if ! ls $d/*_test.go >/dev/null 2>&1; then echo "NO TEST: $d"; fi; done
NO TEST: github.com/astrate-platform/astrate/cmd/astrate
NO TEST: github.com/astrate-platform/astrate/docs
NO TEST: github.com/astrate-platform/astrate/examples/flow-container-echo
NO TEST: github.com/astrate-platform/astrate/internal/appengine
NO TEST: github.com/astrate-platform/astrate/internal/appengine/channels
NO TEST: github.com/astrate-platform/astrate/internal/appengine/stream
NO TEST: github.com/astrate-platform/astrate/internal/auth
NO TEST: github.com/astrate-platform/astrate/internal/broker
NO TEST: github.com/astrate-platform/astrate/internal/config
NO TEST: github.com/astrate-platform/astrate/internal/engine
NO TEST: github.com/astrate-platform/astrate/internal/engine/forward
NO TEST: github.com/astrate-platform/astrate/internal/engine/stream
NO TEST: github.com/astrate-platform/astrate/internal/engine/triggers
NO TEST: github.com/astrate-platform/astrate/internal/flow
NO TEST: github.com/astrate-platform/astrate/internal/flow/blocks
NO TEST: github.com/astrate-platform/astrate/internal/flow/blocks/astartesource
NO TEST: github.com/astrate-platform/astrate/internal/flow/blocks/container
NO TEST: github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool
NO TEST: github.com/astrate-platform/astrate/internal/flowapi
NO TEST: github.com/astrate-platform/astrate/internal/housekeeping
NO TEST: github.com/astrate-platform/astrate/internal/httpx
NO TEST: github.com/astrate-platform/astrate/internal/observability
NO TEST: github.com/astrate-platform/astrate/internal/pairing
NO TEST: github.com/astrate-platform/astrate/internal/pairing/ca
NO TEST: github.com/astrate-platform/astrate/internal/realm
NO TEST: github.com/astrate-platform/astrate/internal/store
NO TEST: github.com/astrate-platform/astrate/internal/swagger
NO TEST: github.com/astrate-platform/astrate/internal/testutil
NO TEST: github.com/astrate-platform/astrate/migrations
NO TEST: github.com/astrate-platform/astrate/pkg/astarteapi
NO TEST: github.com/astrate-platform/astrate/pkg/deviceid
NO TEST: github.com/astrate-platform/astrate/pkg/interfaceschema
NO TEST: github.com/astrate-platform/astrate/pkg/payload

$ sed -n '1,30p;145,160p' internal/flow/pipeline.go; echo "=== deps containing archive for ./internal/flow ==="; timeout 300 go list -deps ./internal/flow 2>&1 | rg -i 'archive|tar|compress' | head
package flow

import (
	"encoding/json"
	"fmt"
)

// Pipeline is an acyclic graph (DAG) of named blocks with typed connections.
// It is a serialisable description; calling Manager.StartFlow instantiates it
// into a running Flow.
type Pipeline struct {
	// ID is a unique identifier for this pipeline.
	ID string `json:"id"`
	// Name is a human-readable label.
	Name string `json:"name"`
	// Blocks is the set of nodes in the graph.
	Blocks []PipelineNode `json:"blocks"`
	// Connections is the set of edges linking block output ports to input ports.
	Connections []Connection `json:"connections"`
}

// PipelineNode describes one block within a pipeline.
type PipelineNode struct {
	// Name is a unique identifier for this node within the pipeline.
	Name string `json:"name"`
	// BlockType identifies which block implementation to use.
	BlockType string `json:"block_type"`
	// Config holds block-specific parameters.
	Config map[string]any `json:"config,omitempty"`
}
	return nil
}

// MarshalJSON serialises the Pipeline. It runs validation before encoding.
func (p *Pipeline) MarshalJSON() ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	type Alias Pipeline
	return json.Marshal((*Alias)(p))
}
=== deps containing archive for ./internal/flow ===
compress/flate
compress/gzip

$ for d in $(go list -f '{{.Dir}}' ./... 2>/dev/null); do if ! ls $d/*_test.go >/dev/null 2>&1; then echo "NO TEST: $d"; fi; done
NO TEST: /root/astrate-mule/docs
NO TEST: /root/astrate-mule/migrations

$ timeout 300 go list -deps ./cmd/astrate 2>/dev/null | rg -c 'moby/go-archive' || echo "go-archive NOT in cmd/astrate deps"; echo "=== binary check ==="; timeout 600 go list -deps -f '{{if eq .ImportPath "github.com/moby/go-archive"}}LINKED{{end}}' ./cmd/astrate 2>/dev/null | rg -v '^$' || echo "NOT LINKED in cmd/astrate"; echo "=== x/text pin ==="; rg -n 'golang.org/x/text|golang.org/x/crypto' go.mod
go-archive NOT in cmd/astrate deps
=== binary check ===
NOT LINKED in cmd/astrate
=== x/text pin ===
21:	golang.org/x/crypto v0.53.0 // bcrypt cost 10 for credentials secrets (§4.1)
97:	golang.org/x/text v0.38.0 // indirect

$ gofmt -l . 2>/dev/null | rg -v '^site-dist' | head; echo "=== gofmt clean above ==="; timeout 900 go test ./... 2>&1 | rg -i 'skip|no test files' | head -25
=== gofmt clean above ===
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]

$ rg -n 't\.Skip|Skipf|SkipNow' --glob '*_test.go' internal/ pkg/ cmd/ test/ 2>/dev/null | head -30
cmd/astrate/main_test.go:47:	t.Skip("the boot suite needs a database: run `make up` or set ASTRATE_TEST_DSN")
cmd/astrate/forward_test.go:68:		t.Skip("built with -tags nats; see forward_nats_test.go")
test/conformance/cpb/cpb_test.go:117:	t.Skip("CP-B needs a database: run `make up` at the repo root or set ASTRATE_TEST_DSN")
test/conformance/cpc/astartectl_test.go:47:		t.Skipf("astartectl: %v (set ASTARTECTL_BIN to run this checkpoint)", err)
test/conformance/cpc/cpc_test.go:122:	t.Skip("CP-C needs a database: run `make up` at the repo root or set ASTRATE_TEST_DSN")
test/conformance/cpa/astartectl_test.go:46:		t.Skipf("astartectl: %v (set ASTARTECTL_BIN to run this checkpoint)", err)
test/conformance/pysdk/pysdk_test.go:47:		t.Skipf("pysdk: %q not found (set ASTRATE_PYSDK_PYTHON)", py)
test/conformance/pysdk/pysdk_test.go:50:		t.Skipf("pysdk: astarte-device-sdk not importable by %s — install test/conformance/pysdk/requirements.txt (%s)", py, out)
test/conformance/cpa/cpa_test.go:132:	t.Skip("CP-A needs a database: run `make up` at the repo root or set ASTRATE_TEST_DSN")
internal/store/datastreams_test.go:420:			t.Skip("timescaledb_toolkit not installed; lttb path unavailable")
internal/pairing/http_test.go:629:		t.Skip("openssl not found in PATH; chain already verified with crypto/x509")
internal/pairing/ca/ca_test.go:293:		t.Skip("10k issuance draw skipped in -short mode")

$ timeout 900 go test -v ./internal/pairing/... ./internal/store/... ./cmd/... 2>&1 | rg -i '^\s*--- SKIP|^--- SKIP|SKIP:' | head -30
(no output)

$ cp go.mod /tmp/opencode/go.mod.bak && cp go.sum /tmp/opencode/go.sum.bak && timeout 600 go mod tidy -diff 2>&1 | head -30; echo "EXIT=$?"; echo "=== restore check ==="; git diff --stat go.mod go.sum
EXIT=0
=== restore check ===

$ timeout 300 go list -m -versions golang.org/x/text 2>&1 | tr ' ' '\n' | tail -5; echo "=== x/crypto available ==="; timeout 300 go list -m -versions golang.org/x/crypto 2>&1 | tr ' ' '\n' | tail -3; echo "=== go-archive available ==="; timeout 300 go list -m -versions github.com/moby/go-archive 2>&1 | tr ' ' '\n' | tail -4
v0.38.0
v0.39.0
v0.40.0
v0.41.0
v0.42.0
=== x/crypto available ===
v0.55.0
v0.56.0
v0.57.0
=== go-archive available ===
v0.3.0
v0.3.1
v0.3.2
v0.3.3

$ ls .github/workflows/ && echo "=== conformance refs in CI ===" && rg -n 'conformance|ASTARTECTL_BIN|ASTRATE_PYSDK' .github/workflows/ Makefile 2>/dev/null | head -20
ci.yml
docs.yml
=== conformance refs in CI ===
Makefile:8:#   make test-conformance    T4 official-SDK conformance harness (lands in M9 under test/)
Makefile:19:.PHONY: build lint test test-integration test-e2e test-conformance tools up down clean
Makefile:48:## test-conformance: T4 — official-SDK harness; lives in test/ with its own go.mod (M9).
Makefile:49:test-conformance:
Makefile:50:	@if [ -f test/conformance/go.mod ]; then \
Makefile:51:		cd test/conformance && $(GO) test -count=1 ./...; \
Makefile:53:		echo "conformance harness lands in M9 (docs/ROADMAP.md §10); nothing to run yet"; \
.github/workflows/ci.yml:2:# (T2) on every push/PR; e2e (T3) + official-SDK conformance (T4) on PRs to main;
.github/workflows/ci.yml:65:  # T4 official-SDK conformance (CP-A/B/C/D): the Go + Python SDKs, the AtomVM
.github/workflows/ci.yml:67:  conformance:
.github/workflows/ci.yml:94:        run: pip install -r test/conformance/pysdk/requirements.txt
.github/workflows/ci.yml:98:          ASTRATE_PYSDK_PYTHON: python3
.github/workflows/ci.yml:100:        run: make test-conformance
.github/workflows/ci.yml:135:        run: cd test/conformance && go test -run TestLoadSmoke -timeout 60m ./load/

$ sed -n '30,60p' test/conformance/cpa/astartectl_test.go; echo "=== pysdk ==="; sed -n '30,60p' test/conformance/pysdk/pysdk_test.go
cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Fatalf("resolving user cache dir: %v", err)
	}
	dir := filepath.Join(cacheDir, "astrate-conformance")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(dir, "astartectl-v"+astartectlVersion)
	if _, err := os.Stat(bin); err == nil {
		return bin
	}

	url, err := astartectlURL()
	if err != nil {
		t.Skipf("astartectl: %v (set ASTARTECTL_BIN to run this checkpoint)", err)
	}
	t.Logf("downloading pinned astartectl %s from %s", astartectlVersion, url)
	if err := downloadAstartectl(url, bin); err != nil {
		t.Fatalf("downloading astartectl: %v", err)
	}
	return bin
}

// astartectlURL maps GOOS/GOARCH onto the release artifact name.
func astartectlURL() (string, error) {
	osName, ok := map[string]string{"darwin": "macOS", "linux": "Linux"}[runtime.GOOS]
	if !ok {
		return "", fmt.Errorf("no pinned astartectl artifact for GOOS=%s", runtime.GOOS)
	}
=== pysdk ===
	ifConf   = "org.astrate.pysdk.Conf"   // device properties
)

var defs = map[string]string{
	ifSensor: `{"interface_name":"` + ifSensor + `","version_major":1,"version_minor":0,"type":"datastream","ownership":"device","mappings":[{"endpoint":"/value","type":"double"}]}`,
	ifConf:   `{"interface_name":"` + ifConf + `","version_major":1,"version_minor":0,"type":"properties","ownership":"device","mappings":[{"endpoint":"/%{k}","type":"string","allow_unset":true}]}`,
}

// pythonWithSDK resolves an interpreter that can import the Astarte SDK, or
// skips the test.
func pythonWithSDK(t *testing.T) string {
	t.Helper()
	py := os.Getenv("ASTRATE_PYSDK_PYTHON")
	if py == "" {
		py = "python3"
	}
	if _, err := exec.LookPath(py); err != nil {
		t.Skipf("pysdk: %q not found (set ASTRATE_PYSDK_PYTHON)", py)
	}
	if out, err := exec.Command(py, "-c", "import astarte.device").CombinedOutput(); err != nil {
		t.Skipf("pysdk: astarte-device-sdk not importable by %s — install test/conformance/pysdk/requirements.txt (%s)", py, out)
	}
	return py
}

func TestPythonSDK(t *testing.T) {
	py := pythonWithSDK(t)
	in := instance.New(t, instance.Config{Interfaces: defs})
	ctx := context.Background()

	id, err := deviceid.Random()

$ $(go env GOPATH)/bin/golangci-lint --version; echo "=== Makefile pin ==="; rg -n 'GOLANGCI_LINT_VERSION|govulncheck|GO_VERSION' Makefile | head; echo "=== govulncheck version ==="; $(go env GOPATH)/bin/govulncheck -version 2>&1 | head -3
golangci-lint has version 2.12.2 built with go1.26.5 from (unknown, modified: ?, mod sum: "h1:7+d1uY0bq1MU2UV3R5pW5Q7QWdcoq4naMRXM+gsJKrs=") on (unknown)
=== Makefile pin ===
12:GOLANGCI_LINT_VERSION  ?= v2.12.2
33:		echo "golangci-lint not found; run 'make tools' (installs $(GOLANGCI_LINT_VERSION))"; exit 1; }
58:	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
=== govulncheck version ===
Go: go1.26.7
Scanner: govulncheck@v1.8.0
DB: https://vuln.go.dev

$ echo "=== does the built binary contain go-archive symbols? ==="; timeout 900 go build -o /tmp/opencode/astrate-probe ./cmd/astrate 2>&1 | tail -5 && go tool nm /tmp/opencode/astrate-probe 2>/dev/null | rg -c 'go-archive' || echo "NO go-archive symbols in binary"
=== does the built binary contain go-archive symbols? ===
NO go-archive symbols in binary

→ Read .github/workflows/ci.yml
→ Read test/conformance/cpa/astartectl_test.go

$ ls test/conformance/ && echo "=== go.mod present? ==="; ls test/conformance/go.mod 2>&1; echo "=== load dir ==="; ls test/conformance/load/ 2>&1 | head
atomvm
cpa
cpb
cpc
go.mod
go.sum
gosdk
instance
load
pysdk
upstream
=== go.mod present? ===
test/conformance/go.mod
=== load dir ===
load_test.go

$ echo "=== any archive pkg in internal/flow deps? ==="; timeout 300 go list -deps ./internal/flow 2>/dev/null | rg 'archive' || echo "NO archive* package in internal/flow deps"; echo; echo "=== go.mod go/toolchain ==="; head -6 go.mod
=== any archive pkg in internal/flow deps? ===
NO archive* package in internal/flow deps

=== go.mod go/toolchain ===
module github.com/astrate-platform/astrate

go 1.26.1

// Pinned dependency decisions — docs/ROADMAP.md §1.1 (frozen in docs/DESIGN.md).
// HTTP routing is stdlib net/http (Go >= 1.22 ServeMux patterns): deliberately no framework.

$ rg -n 'nats' Makefile .github/workflows/ci.yml .golangci.yml 2>/dev/null | head -20; echo "=== forward pkg files ==="; ls internal/engine/forward/
=== forward pkg files ===
envelope.go
envelope_test.go
http.go
http_test.go
nats.go
nats_test.go

