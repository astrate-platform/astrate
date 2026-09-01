slug: recipe-docs-sync
verdict: proposed
at:  689fb1d
ran: 2026-09-01T20:06:32Z on DietPi in 259s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/docs-sync.md
I'll read the recipe and carry it out. Let me start by understanding the current state — checking the surfaces, the todo, and for-giulio files.

$ ls docs/api/ && echo "---TODO---" && tail -20 .mule/todo.md && echo "---FORGIULIO---" && test -f .mule/for-giulio.md && cat .mule/for-giulio.md || echo "no for-giulio"
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml
---TODO---
- [ ] race-check-pkg: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./pkg/... ./cmd/... ./internal/testutil/...`. Report any failure to .mule/for-giulio.md with the full race report. [legion] [readonly]
- [x] store-realm-cascade-policies: in `internal/store/realms_test.go` `CascadeDelete`, add `trigger_policies` to the post-delete verification loop (query `SELECT count(*) FROM trigger_policies WHERE realm_id = $1`). The migration 000006 has ON DELETE CASCADE but nothing asserts on it. [auto]
- [x] store-alias-lowest-id: in `internal/store/devices_test.go`, add a subtest that registers two devices in the same realm, sets the same alias tag on both, and asserts `GetDeviceByAlias` returns the one with the lower UUID. The SQL uses `ORDER BY id LIMIT 1` but no test proves it. [auto]
- [x] store-delete-device-objects: in `internal/store/devices_test.go` `StatsAndDelete`, insert object datastream rows for the device before deleting it, and assert they are gone after the delete. `DeleteDevice` explicitly sweeps `object_datastreams` but the test only checks individual rows. [auto]
- [!] control-producer-properties-compression: in `internal/engine/control.go`, accept a plaintext (and the 4-zero-byte empty-frame) device→server `producer/properties` purge list when the device declared `purge_properties_compression_format: plaintext` — `inflateProperties` currently rejects every non-zlib payload while upstream's `control_handler.ex` decodes per-capability **Approved by Giulio 2026-08-31**: yes, devices that cannot compress must be able to talk to us, following upstream. — BLOCKED: gates failed
- [!] probe-interface-default-values: does `GET /realmmanagement/v1/<realm>/interfaces/<name>/<major>` return the same defaulted mapping-parameter values upstream reports after 3f0b864 ("Properly show default values for all mapping parameters")? report, do not patch [auto] — BLOCKED: gates failed
- [!] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto] — BLOCKED: gates failed
- [!] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
- [!] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto] — BLOCKED: gates failed
- [!] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
- [!] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto] — BLOCKED: gates failed
- [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
- [!] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto] — BLOCKED: gates failed
- [!] docs-sync-pairing-status-enum: in docs/api/astarte_pairing_api.yaml, fix the PairingInfo.status enum (line 367): it lists `confirmed, pending, denied, expired` but the handler only ever emits `confirmed`, `pending`, or `inhibited` (internal/pairing/service.go:289-297, via internal/pairing/http.go:211) — drop the dead `denied`/`expired` and add the undocumented `inhibited`. [auto] — BLOCKED: gates failed
- [!] issue-91-pairing-health-serve-root: in internal/pairing/http.go, add the unauthenticated `GET /pairing/health` route (no realm segment) sharing the handler/payload of the existing `GET /pairing/v1/{realm}/health`, keeping the v1 route as-is; tests for both paths (closes #91) — BLOCKED: gates failed
- [!] issue-68-async-operation-accepted: on the mutating endpoints upstream surfaces `async_operation` on (housekeeping realm create/delete; realm-management interface install/update/delete, trigger/policy delete in internal/realm/http.go), accept and ignore `?async_operation=false`, with unparseable/`true` values also not changing behaviour (Astrate stays always-sync); tests cover the flag parsing (closes #68) — BLOCKED: gates failed
- [!] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto] — BLOCKED: gates failed
- [!] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
- [!] realm-policy-list-sorted: in `internal/realm/service.go`, sort the names returned by `ListPolicies` the way `ListTriggers`/`ListInterfaces`/`ListInterfaceMajors` already do (it is the only list method without `sort.Strings`), and extend `TestDashboardCompat.Policies` (or a new subtest in `http_test.go`) to install several policies in non-alphabetical order and assert `GET /policies` comes back sorted. [auto] — BLOCKED: gates failed
- [!] realm-interface-lookup-404: in `internal/realm/http_test.go` `TestRealmManagement`, assert the 404 paths that nothing currently covers — `GET /interfaces/<nonexistent-name>` and `GET /interfaces/<name>/<no-such-major>` (service.go `ListInterfaceMajors`:211-213 and `GetInterface`, mapped to 404 via writeError). A rule with no test currently. [auto] — BLOCKED: gates failed
---FORGIULIO---
# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **Milestones recipe has nothing to file: v2.0 is already DONE, and `mule/queue` is 120
  commits behind `origin/main`, so this branch's `.mule/milestones.md` is stale.** My branch
  still shows v2.0 "not started" (and v3.0 as CLEA TBD), so the recipe pointed me at v2.0 —
  but on `main` v2.0 is marked **DONE (2026-07-29)**: every `milestone-2.0` issue is CLOSED
  (#23–#27, #37, #39, #40, #41, #42, #43), and the factory/catalog/container-block/durable
  named flows/`/flow/v1` API all live on `main` (`git show origin/main:.mule/milestones.md`).
  `gh issue list --label milestone-2.0 --state open` is empty and there are no new v2.0 gaps
  to file, so per the recipe's step 5 I filed **nothing** rather than duplicate closed work.
  **Real current target is v3.0** ("upstream Astarte feature parity 1.3/1.4", scoped on main
  2026-08-22): the first recipe job there is to triage `upstream-parity` backlog **#47–#89**
  (source `.mule/research/upstream-parity-2026-08-22.md`) into a plan. But I should not start
  that triage from this stale branch. Action needed: **merge/rebuild `mule/queue` on top of
  `origin/main`**, then the next milestones run targets v3.0 with the correct research doc.
  (Milestones recipe run, 2026-09-01.)

---

- ~~`device_deletion_started`/`device_deletion_finished` trigger events are not emitted~~ —
  **decided 2026-07-27: emit both, back-to-back, around the synchronous delete.** Filed as
  issue #21 (`mule`). (Cross-project survey, 2026-07-27,
  `.mule/research/survey-2026-07-27.md` source 4.)
- ~~Mustache trigger-action templates are accepted but not rendered~~ — **decided
  2026-07-27: implement it.** Guiding principle clarified: Astarte compatibility means
  SDK/wire compatibility, not minimum dependency count — Astrate is allowed to be a
  compatible *superset*. Library picked: `github.com/cbroglie/mustache`. Filed as issue #22
  (`mule`). (Same survey, source 4.)
- **`value_change`/`value_change_applied`/`path_created`/`path_removed`/`value_stored` trigger
  types compile but never fire** (`internal/engine/triggers/match.go:30-42`). Decision
  deferred pending data: issue #20 (`mule`, `readonly`) asks Big Pickle to benchmark the
  ingest-path cost of a previous-value lookup on the Legion Go before this gets decided one
  way or the other. (Same survey, source 4.)
- **Group-scoped triggers (`group_name` on device/data triggers) compile but never match**
  (`internal/engine/triggers/match.go:11-12`). Decision deferred, tied to issue #17
  (group-WATCH-path reconciliation, trickle work, not mule): whatever group-membership
  mechanism comes out of that phase should also report the perf cost for this decision —
  noted in a comment on #17 so it isn't benchmarked twice. (Same survey, source 4.)

---

- ~~The Pi cannot run the race detector~~ — **resolved 2026-07-27** by installing Go 1.26.5
  as a userland toolchain on the Legion Go (`~/.local/go`, no root, `rm -rf` to undo). The
  Pi still cannot run `-race` (39-bit VMA kernel vs the 48 ThreadSanitizer needs), so its
  gate remains `go vet ./... && go test ./...` — but race coverage now exists on the Legion
  Go, where the full suite runs clean in ~40s on 16 cores. The standing `race-check` task is
  the concurrency gate. Concurrency work is queueable again, provided the race-check runs
  after it.
- **golangci-lint is not installed on the Pi**, so the mule's second gate is silently absent
  there — `gofmt` still runs, the linter does not. `go install
  github.com/golangci/golangci-lint/cmd/golangci-lint@<the pinned version>` on the Pi would
  close it; I did not pick a version for you, since the pin is a decision.
- ~~`/root/astrate` on the Pi has uncommitted work~~ — **resolved 2026-07-27** with the new
  `tools/reconcile.sh`: rescued onto `origin/wip/DietPi-20260727T171543Z` (pushed, not
  reviewed — read the diff before merging anything from it) and `/root/astrate` is now clean
  on `main`. Also had to set `commit.gpgsign false` locally on that clone first — it had
  signing on with no working gpg-agent for a non-interactive session, same class of problem
  `/root/astrate-mule` already solved the same way (your call, confirmed 2026-07-27: align it
  with the existing unattended-Pi-clone convention rather than fix gpg-agent or bypass signing
  silently inside the script).
- **The mule has been idle 15h.** Filed by the dead-man's switch; see journalctl on the Pi.

$ rg -o 'astrate_[a-z_]+_api' docs/api/*.yaml | sort -u
(no output)

# Todos
[•] Pick a surface and find drift between internal/ handlers and docs/api yaml
[ ] Check documented status codes/response fields for 3-4 endpoints
[ ] Check ASTRATE_ config keys vs configuration-reference.md
[ ] Append proposals to .mule/todo.md and for-giulio.md

Let me check the API directory contents and the surfaces. The YAML files use `astarte_` prefix not `astrate_`. Let me examine each surface one at a time. Let me pick a surface per run — the recipe says pick one surface per run. Let me check the current state and pick one surface.

Let me set up a todo and pick a surface. First let me look at what's already been done (the pairing one is already in todo as blocked). Let me check the individual surfaces.
The pairing surface already has proposals in the todo (blocked ones came from a previous run). Let me pick a fresh surface. The recipe lists five: appengine, housekeeping, pairing, realm-management, astrate-native. Pairing is already well-probed. Let me pick one of the others. Let me check what surfaces exist in the code and pick realm-management or astrate-native.

Let me explore the surfaces in the code.
$ ls internal/ | grep -Ei 'appengine|housekeep|pairing|realm|native' 
appengine
housekeeping
pairing
realm

$ rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/realm/ | head -40
(no output)

Let me look at the routing in internal/realm. The pattern from the recipe may not match the actual code style. Let me explore.
$ rg -n 'Get\(|Post\(|Put\(|Delete\(|Patch\(|HandleFunc|\.Methods\(' internal/realm/ | head -40
(no output)

$ rg -n '^\s{2}/' docs/api/astarte_realm_management_api.yaml | head -60
27:  /realmmanagement/v1/{realm}/interfaces:
98:  /realmmanagement/v1/{realm}/interfaces/{name}:
127:  /realmmanagement/v1/{realm}/interfaces/{name}/{major}:
239:  /realmmanagement/v1/{realm}/triggers:
309:  /realmmanagement/v1/{realm}/triggers/{name}:
361:  /realmmanagement/v1/{realm}/config/auth:
428:  /realmmanagement/v1/{realm}/config/device_registration_limit:
461:  /realmmanagement/v1/{realm}/version:
494:  /realmmanagement/v1/{realm}/devices/{device}:
520:  /realmmanagement/v1/{realm}/policies:
590:  /realmmanagement/v1/{realm}/policies/{name}:

$ ls internal/realm/
dashboard_compat_test.go
http.go
http_test.go
policy.go
policy_test.go
service.go

Let me find the actual routing mechanism in the realm service.
$ rg -n 'mux|Handle|router|ServeHTTP|func.*Serve|\.Path\(' internal/realm/*.go | head -40
internal/realm/dashboard_compat_test.go:70:			ErrorHandlers []any  `json:"error_handlers"`
internal/realm/dashboard_compat_test.go:73:		if got.Name != "retry5xx" || got.RetryTimes != 3 || len(got.ErrorHandlers) != 1 {
internal/realm/http_test.go:46:	mux      *http.ServeMux
internal/realm/http_test.go:88:	mux := http.NewServeMux()
internal/realm/http_test.go:89:	NewAPI(svc, auth.NewMiddleware(st)).Mount(mux)
internal/realm/http_test.go:92:		st: st, svc: svc, mux: mux, realm: realmName, realmID: realm.ID, jwtKey: key,
internal/realm/http_test.go:112:	r.mux.ServeHTTP(rec, httpReq)
internal/realm/http_test.go:252:func discardLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }
internal/realm/http.go:22:	require func(http.Handler) http.Handler
internal/realm/http.go:31:// Mount registers the routes on mux (paths wire-frozen to upstream
internal/realm/http.go:33:func (a *API) Mount(mux *http.ServeMux) {
internal/realm/http.go:34:	h := func(f http.HandlerFunc) http.Handler { return a.require(f) }
internal/realm/http.go:35:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces", h(a.listInterfaces))
internal/realm/http.go:36:	mux.Handle("POST /realmmanagement/v1/{realm}/interfaces", h(a.installInterface))
internal/realm/http.go:37:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}", h(a.listInterfaceMajors))
internal/realm/http.go:38:	mux.Handle("GET /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.getInterface))
internal/realm/http.go:39:	mux.Handle("PUT /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.updateInterface))
internal/realm/http.go:40:	mux.Handle("DELETE /realmmanagement/v1/{realm}/interfaces/{name}/{major}", h(a.deleteInterface))
internal/realm/http.go:41:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers", h(a.listTriggers))
internal/realm/http.go:42:	mux.Handle("POST /realmmanagement/v1/{realm}/triggers", h(a.createTrigger))
internal/realm/http.go:43:	mux.Handle("GET /realmmanagement/v1/{realm}/triggers/{name}", h(a.getTrigger))
internal/realm/http.go:44:	mux.Handle("DELETE /realmmanagement/v1/{realm}/triggers/{name}", h(a.deleteTrigger))
internal/realm/http.go:45:	mux.Handle("GET /realmmanagement/v1/{realm}/config/auth", h(a.getAuth))
internal/realm/http.go:46:	mux.Handle("PUT /realmmanagement/v1/{realm}/config/auth", h(a.putAuth))
internal/realm/http.go:47:	mux.Handle("GET /realmmanagement/v1/{realm}/config/device_registration_limit", h(a.getRegistrationLimit))
internal/realm/http.go:48:	mux.Handle("GET /realmmanagement/v1/{realm}/version", h(a.getVersion))
internal/realm/http.go:49:	mux.Handle("DELETE /realmmanagement/v1/{realm}/devices/{device}", h(a.deleteDevice))
internal/realm/http.go:50:	mux.Handle("GET /realmmanagement/v1/{realm}/policies", h(a.listPolicies))
internal/realm/http.go:51:	mux.Handle("POST /realmmanagement/v1/{realm}/policies", h(a.createPolicy))
internal/realm/http.go:52:	mux.Handle("GET /realmmanagement/v1/{realm}/policies/{name}", h(a.getPolicy))
internal/realm/http.go:53:	mux.Handle("DELETE /realmmanagement/v1/{realm}/policies/{name}", h(a.deletePolicy))

Now I have the routing. Let me compare the documented paths against the mounted routes for the realm-management surface.

Documented paths (from yaml):
- `/realmmanagement/v1/{realm}/interfaces` (GET, POST)
- `/realmmanagement/v1/{realm}/interfaces/{name}` 
- `/realmmanagement/v1/{realm}/interfaces/{name}/{major}`
- `/realmmanagement/v1/{realm}/triggers`
- `/realmmanagement/v1/{realm}/triggers/{name}`
- `/realmmanagement/v1/{realm}/config/auth`
- `/realmmanagement/v1/{realm}/config/device_registration_limit`
- `/realmmanagement/v1/{realm}/version`
- `/realmmanagement/v1/{realm}/devices/{device}`
- `/realmmanagement/v1/{realm}/policies`
- `/realmmanagement/v1/{realm}/policies/{name}`

Mounted routes: all match. Let me check the methods documented in the yaml against the mounted methods. Let me look at the details of the yaml paths for each endpooint methods.

Let me check the documented HTTP methods per path in the yaml.
$ rg -n '^\s{2}/|^    (get|post|put|delete|patch):' docs/api/astarte_realm_management_api.yaml | grep -A1 '^\s*[0-9]*:\s\{2\}/' ; echo "====="; awk '/^  \//{print NR": "$0} /^    (get|post|put|delete|patch):/{print NR": "$0}' docs/api/astarte_realm_management_api.yaml
27:  /realmmanagement/v1/{realm}/interfaces:
28:    get:
--
98:  /realmmanagement/v1/{realm}/interfaces/{name}:
99:    get:
127:  /realmmanagement/v1/{realm}/interfaces/{name}/{major}:
128:    get:
--
239:  /realmmanagement/v1/{realm}/triggers:
240:    get:
--
309:  /realmmanagement/v1/{realm}/triggers/{name}:
310:    get:
--
361:  /realmmanagement/v1/{realm}/config/auth:
362:    get:
--
428:  /realmmanagement/v1/{realm}/config/device_registration_limit:
429:    get:
461:  /realmmanagement/v1/{realm}/version:
462:    get:
494:  /realmmanagement/v1/{realm}/devices/{device}:
495:    delete:
520:  /realmmanagement/v1/{realm}/policies:
521:    get:
--
590:  /realmmanagement/v1/{realm}/policies/{name}:
591:    get:
=====
27:   /realmmanagement/v1/{realm}/interfaces:
28:     get:
53:     post:
98:   /realmmanagement/v1/{realm}/interfaces/{name}:
99:     get:
127:   /realmmanagement/v1/{realm}/interfaces/{name}/{major}:
128:     get:
158:     put:
196:     delete:
239:   /realmmanagement/v1/{realm}/triggers:
240:     get:
264:     post:
309:   /realmmanagement/v1/{realm}/triggers/{name}:
310:     get:
339:     delete:
361:   /realmmanagement/v1/{realm}/config/auth:
362:     get:
393:     put:
428:   /realmmanagement/v1/{realm}/config/device_registration_limit:
429:     get:
461:   /realmmanagement/v1/{realm}/version:
462:     get:
494:   /realmmanagement/v1/{realm}/devices/{device}:
495:     delete:
520:   /realmmanagement/v1/{realm}/policies:
521:     get:
545:     post:
590:   /realmmanagement/v1/{realm}/policies/{name}:
591:     get:
620:     delete:

The methods and paths all match between the yaml and the code for the realm-management surface. The surface is well-synced at the path/method level. Let me pick a different surface — astrate-native or appengine. Let me check which is more likely to have drift.

Let me look at astrate-native which is a new surface, and appengine. Let me first check the documented paths and methods for astrate-native.
$ ls internal/appengine/ 2>/dev/null; echo "===NATIVE==="; ls internal/ 2>/dev/null | grep -i native
channels
dashboard_compat_test.go
data.go
downsample.go
downsample_test.go
e2e_test.go
http.go
http_test.go
links_test.go
query_opts_test.go
service.go
stream
===NATIVE===

$ awk '/^  \//{print NR": "$0} /^    (get|post|put|delete|patch):/{print NR": "$0}' docs/api/astrate_native_api.yaml
26:   /astrate/v1/health:
27:     get:
44:   /astrate/v1/readiness:
45:     get:
77:   /astrate/v1/metrics:
78:     get:
99:   /appengine/health:
100:     get:
118:   /realmmanagement/health:
119:     get:
137:   /pairing/health:
138:     get:
158:   /astrate/v1/{realm}/socket:
159:     get:
220:   /appengine/v1/socket/websocket:
221:     get:

The astrate-native surface has some documented paths. Let me check where these routes are mounted in the code. Let me search for these paths.
$ rg -n 'astrate/v1|appengine/v1/socket|realmmanagement/health|pairing/health|appengine/health|astrate/v1/health|astrate/v1/readiness|astrate/v1/metrics' internal/ --glob '*.go' | grep -i 'Handle\|Mount\|mux' | head -40
internal/observability/observability_test.go:20:	m.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/metrics", nil))
internal/observability/observability_test.go:49:	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/health", nil))
internal/observability/observability_test.go:56:	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/readiness", nil))
internal/observability/observability_test.go:72:	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/astrate/v1/readiness", nil))
internal/observability/health.go:43:	mux.HandleFunc("GET /astrate/v1/health", h.handleHealth)
internal/observability/health.go:44:	mux.HandleFunc("GET /astrate/v1/readiness", h.handleReadiness)
internal/observability/health.go:45:	mux.Handle("GET /astrate/v1/metrics", h.metrics)
internal/appengine/stream/ws.go:46:	mux.Handle("GET /astrate/v1/{realm}/socket", a.require(http.HandlerFunc(a.handle)))
internal/appengine/channels/ws.go:48:	mux.Handle("GET /appengine/v1/socket/websocket", http.HandlerFunc(a.handle))

Let me look at the astrate-native yaml details for the endpoints and compare against the code that serves them. Let me look at the health endpoints and the metrics endpoint in detail. Let me read the yaml section for these endpoints.

Let me check the appengine surface as well. Actually, let me focus. The realm-management surface turned out to be fully in sync. Let me pick a surface and do a focused comparison. Since the recipe says pick one surface per run, let me pick **astrate-native** since it's the newest/most Astrate-specific.

Let me look at the astrate-native yaml in detail — the health endpoints, metrics, and the socket paths.
→ Read docs/api/astrate_native_api.yaml [limit=220]
Now let me check the actual code for these endpoints. The astrate-native surface documents the health, readiness, metrics, compat health, and two websocket endpoints. Let me verify these against the code.

Let me check the observability health handlers and the streams.
→ Read internal/observability/health.go
The health/readiness/metrics endpoints match the docs. Let me check the compat health endpoints (`/appengine/health`, `/realmmanagement/health`, `/pairing/health`) and the websocket endpoints.

Let me find where these compat health endpoints and the two websocket endpoints are mounted.

$ rg -rn 'appengine/health|realmmanagement/health|pairing/health' cmd/ internal/ tools/ 2>/dev/null | grep -i 'Handle\|Mount\|HandleFunc' | head; echo "===COMMAND==="; rg -rl 'Mount' cmd/ | head
===COMMAND===
cmd/astrate/main.go:	}).l(mux)
cmd/astrate/main.go:	housekeeping.NewAPI(hkSvc, mw, hkKeys).l(mux)
cmd/astrate/main.go:	realm.NewAPI(realm.NewService(st, e, log).WithDisconnecter(b), mw).l(mux)
cmd/astrate/main.go:	appengine.NewAPI(appengine.NewService(st, e, log), mw).l(mux)
cmd/astrate/main.go:	apstream.NewAPI(e.Bus(), mw).l(mux)
cmd/astrate/main.go:	channels.NewAPI(e.Bus(), st).l(mux)
cmd/astrate/main.go:		observability.lServiceCompat(mux, svc)
cmd/astrate/main.go:	health.l(mux)
cmd/astrate/main.go:	swagger.l(mux)

$ rg -n 'health|ServiceCompat|Handle' cmd/astrate/main.go | head -30
56:	healthcheck := flag.Bool("healthcheck", false, "probe the local readiness endpoint and exit (for container HEALTHCHECK)")
63:	if *healthcheck {
154:	srv := &http.Server{Addr: cfg.HTTP.Addr, Handler: handler, ReadHeaderTimeout: 10 * time.Second}
283:func mountAPIs(cfg config.Config, st *store.Store, e *engine.Engine, b *broker.Broker, sealer *store.KeySealer, metrics *observability.Metrics, log *slog.Logger) (http.Handler, *housekeeping.Service, error) {
298:	pairer.OnRegistered = e.HandleDeviceRegistered
319:	// Upstream-parity per-service health endpoints (the dashboard's API
322:		observability.MountServiceCompat(mux, svc)
335:	health := observability.NewHealth(metrics.Handler())
336:	health.AddReadiness("database", st.Health)
337:	health.AddReadiness("broker", brokerReadiness(b))
338:	health.Mount(mux)
437:// `astrate -healthcheck`. The HTTP address comes from ASTRATE_HTTP_ADDR.
480:	opts := &slog.HandlerOptions{Level: level}
481:	var h slog.Handler
483:		h = slog.NewTextHandler(os.Stderr, opts)
485:		h = slog.NewJSONHandler(os.Stderr, opts)

$ rg -n 'MountServiceCompat' internal/observability/*.go
internal/observability/compat.go:7:// MountServiceCompat registers the upstream-parity unauthenticated
internal/observability/compat.go:11:func MountServiceCompat(mux *http.ServeMux, service string) {
internal/observability/observability_test.go:78:func TestMountServiceCompat(t *testing.T) {
internal/observability/observability_test.go:81:		MountServiceCompat(mux, svc)

→ Read internal/observability/compat.go
The compat health endpoints match. Now let me verify the websocket endpoints. `internal/appengine/stream/ws.go` handles `/astrate/v1/{realm}/socket` and `internal/appengine/channels/ws.go` handles `/appengine/v1/socket/websocket`. Let me check the response status codes documented (101, 200, 401) and the SSE/websocket behavior against the code.

Let me read the rest of the astrate-native yaml and the two websocket handlers.
→ Read docs/api/astrate_native_api.yaml [offset=220, limit=116]
Now let me check the two websocket handlers to verify the documented behavior. Let me look at the stream ws.go handler for `/astrate/v1/{realm}/socket`.
→ Read internal/appengine/stream/ws.go
The `/astrate/v1/{realm}/socket` endpoint matches the docs well (websocket 101, SSE 200, 401 via middleware). 

Now let me check the Phoenix V2 channels endpoint `/appengine/v1/socket/websocket` — the docs describe a Phoenix V2 protocol with `phx_join`, `heartbeat`, `watch`/`unwatch`, `new_event` frames, and list 500 as a possible response. Let me read the channels ws.go to verify.
→ Read internal/appengine/channels/ws.go
The astrate-native surface's websocket endpoints match the docs well. The status codes documented (101, 200, 401, 500) are consistent with the code.

The astrate-native surface appears well-synced too. Let me look at the appengine surface, which is a larger data API and more likely to have drift. Let me compare its documented paths with the mounted routes.

Let me look at the appengine surface paths in the yaml and the http.go mount.

$ awk '/^  \//{print NR": "$0} /^    (get|post|put|delete|patch):/{print NR": "$0}' docs/api/astarte_appengine_api.yaml
26:   /appengine/v1/{realm}/devices:
27:     get:
93:   /appengine/v1/{realm}/stats/devices:
94:     get:
131:   /appengine/v1/{realm}/devices/{device}:
132:     get:
156:     patch:
201:   /appengine/v1/{realm}/devices-by-alias/{alias}:
202:     get:
233:   /appengine/v1/{realm}/devices/{device}/interfaces/{interface}:
234:     get:
268:   /appengine/v1/{realm}/devices/{device}/interfaces/{interface}/{path}:
269:     get:
306:     put:
347:     post:
388:     delete:
419:   /appengine/v1/{realm}/groups:
420:     get:
445:     post:
489:   /appengine/v1/{realm}/groups/{group}/devices:
490:     get:
524:     post:
564:   /appengine/v1/{realm}/groups/{group}/devices/{device}:
565:     delete:

$ rg -n 'mux.Handle\(' internal/appengine/http.go | head -60
37:	mux.Handle("GET "+base+"/devices", h(a.listDevices))
38:	mux.Handle("GET "+base+"/stats/devices", h(a.devicesStats))
39:	mux.Handle("GET "+base+"/devices/{device}", h(a.getDevice))
40:	mux.Handle("PATCH "+base+"/devices/{device}", h(a.patchDevice))
41:	mux.Handle("GET "+base+"/devices-by-alias/{alias}", h(a.getDeviceByAlias))
43:	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}", h(a.getData))
44:	mux.Handle("GET "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.getData))
45:	mux.Handle("PUT "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
46:	mux.Handle("POST "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.putData))
47:	mux.Handle("DELETE "+base+"/devices/{device}/interfaces/{interface}/{path...}", h(a.deleteData))
49:	mux.Handle("GET "+base+"/groups", h(a.listGroups))
50:	mux.Handle("POST "+base+"/groups", h(a.createGroup))
51:	mux.Handle("GET "+base+"/groups/{group}/devices", h(a.listGroupDevices))
52:	mux.Handle("POST "+base+"/groups/{group}/devices", h(a.addGroupDevice))
53:	mux.Handle("DELETE "+base+"/groups/{group}/devices/{device}", h(a.removeGroupDevice))

The appengine surface paths and methods match too. All surfaces' paths and methods appear in sync at the route level. 

The recipe says pick one surface per run and check 3-4 endpoints' status codes and response fields. I've verified the astrate-native surface thoroughly. Let me note that I did the astrate-native surface. Now I need to check some detailed status codes/response fields for 3-4 endpoints. Let me check the fetch of a specific detail.

Actually, let me reconsider. I've now looked at three surfaces (realm-management, astrate-native, appengine) at the route level and all matched. Let me pick **appengine** as the surface for this run, and do the deeper status-code/response-field check on 3-4 of its endpoints — the data paths (GET/PUT/POST/DELETE on individual/object datastreams).

Let me also check the config keys inventory, which is a separate useful check in the recipe.

Let me check the appengine getData handler and its response behavior against the documented status codes. Let me look at the yaml for the data endpoints in detail (lines 233-418).
→ Read docs/api/astarte_appengine_api.yaml [offset=233, limit=190]
Now let me check the getData handler behavior and its status codes/response envelope against the docs. Let me look at the appengine data.go / http.go handlers for getData, putData (returns 200), deleteData (returns 204).

Let me read the relevant handler code for status codes.

$ rg -n 'func (a \*API) (getData|putData|deleteData)|StatusOK|StatusNoContent|StatusNotFound|writeError|DeviceNotFound|GetDevice\b|WriteData' internal/appengine/http.go | head -50
64:		a.writeError(w, err)
71:	_ = astarteapi.WriteDataWithLinks(w, http.StatusOK, body,
106:		a.writeError(w, err)
109:	_ = astarteapi.WriteData(w, http.StatusOK, map[string]int64{
116:	st, err := a.svc.GetDevice(r.Context(), r.PathValue("realm"), r.PathValue("device"))
118:		a.writeError(w, err)
121:	_ = astarteapi.WriteData(w, http.StatusOK, st)
127:		a.writeError(w, err)
130:	_ = astarteapi.WriteData(w, http.StatusOK, st)
148:		a.writeError(w, err)
151:	_ = astarteapi.WriteData(w, http.StatusOK, st)
165:		a.writeError(w, err)
168:	_ = astarteapi.WriteData(w, http.StatusOK, data)
180:		a.writeError(w, err)
183:	w.WriteHeader(http.StatusOK)
190:		a.writeError(w, err)
193:	w.WriteHeader(http.StatusNoContent)
201:		a.writeError(w, err)
204:	_ = astarteapi.WriteData(w, http.StatusOK, names)
220:		a.writeError(w, err)
223:	_ = astarteapi.WriteData(w, http.StatusCreated, body)
230:		a.writeError(w, err)
234:		_ = astarteapi.WriteData(w, http.StatusOK, page.Statuses)
237:	_ = astarteapi.WriteData(w, http.StatusOK, page.IDs)
252:		a.writeError(w, err)
260:		a.writeError(w, err)
263:	w.WriteHeader(http.StatusNoContent)
315:// writeError maps service/store errors onto upstream-shaped responses.
316:func (a *API) writeError(w http.ResponseWriter, err error) {
323:		_ = astarteapi.WriteDeviceNotFound(w)

The status codes for getData (200), putData (200), deleteData (204) all match the docs. The appengine surface is well-synced.

Now let me do the config keys check, which is part of the recipe (docs/site/configuration-reference.md vs ASTRATE_ keys in code). This is a useful parallel check. Let me run the config key comparison.
$ rg -o '\bASTRATE_[A-Z_]+' -N internal/ | sort -u
internal/config/config.example.toml:ASTRATE_DATABASE_DSN
internal/config/config.example.toml:ASTRATE_ENGINE_SHARDS
internal/config/config.example.toml:ASTRATE_HTTP_ADDR
internal/config/config.example.toml:ASTRATE_LOG_LEVEL
internal/config/config.example.toml:ASTRATE_MASTER_KEY
internal/config/config.example.toml:ASTRATE_MASTER_KEY_FILE
internal/config/config.example.toml:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/config/config.go:ASTRATE_DATABASE_DSN
internal/config/config.go:ASTRATE_ENGINE_SHARDS
internal/config/config.go:ASTRATE_HTTP_ADDR
internal/config/config.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
internal/config/config.go:ASTRATE_HTTP_TLS_CERT_FILE
internal/config/config.go:ASTRATE_HTTP_TLS_KEY_FILE
internal/config/config.go:ASTRATE_LOG_FORMAT
internal/config/config.go:ASTRATE_LOG_LEVEL
internal/config/config.go:ASTRATE_MASTER_KEY
internal/config/config.go:ASTRATE_MASTER_KEY_FILE
internal/config/config.go:ASTRATE_MQTT_ADDR
internal/config/config.go:ASTRATE_MQTT_ADVERTISED_URL
internal/config/config.go:ASTRATE_MQTT_INSECURE_DEV_MODE
internal/config/config.go:ASTRATE_MQTT_SESSION_STORE_PATH
internal/config/config.go:ASTRATE_MQTT_TLS_CERT_FILE
internal/config/config.go:ASTRATE_MQTT_TLS_KEY_FILE
internal/config/config.go:ASTRATE_REALM_JWT_PUBLIC_KEY
internal/config/config.go:ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
internal/config/config.go:ASTRATE_REALM_NAME
internal/config/config.go:ASTRATE_SECURITY_MASTER_KEY_FILE
internal/config/config_test.go:ASTRATE_DATABASE_DSN
internal/config/config_test.go:ASTRATE_ENGINE_SHARDS
internal/config/config_test.go:ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
internal/engine/engine_integration_test.go:ASTRATE_TEST_DSN
internal/store/crypto.go:ASTRATE_MASTER_KEY
internal/store/crypto.go:ASTRATE_MASTER_KEY_FILE
internal/store/store_test.go:ASTRATE_TEST_DSN
internal/testutil/pg.go:ASTRATE_TEST_DSN

$ rg -o '\bASTRATE_[A-Z_]+' -N docs/site/configuration-reference.md | sort -u
ASTRATE_DATABASE_DSN
ASTRATE_ENGINE_SHARDS
ASTRATE_HTTP_ADDR
ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
ASTRATE_HTTP_TLS_CERT_FILE
ASTRATE_HTTP_TLS_KEY_FILE
ASTRATE_LOG_FORMAT
ASTRATE_LOG_LEVEL
ASTRATE_MASTER_KEY
ASTRATE_MASTER_KEY_FILE
ASTRATE_MQTT_ADDR
ASTRATE_MQTT_ADVERTISED_URL
ASTRATE_MQTT_INSECURE_DEV_MODE
ASTRATE_MQTT_SESSION_STORE_PATH
ASTRATE_MQTT_TLS_CERT_FILE
ASTRATE_MQTT_TLS_KEY_FILE
ASTRATE_REALM_JWT_PUBLIC_KEY
ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
ASTRATE_REALM_NAME
ASTRATE_SECURITY_MASTER_KEY_FILE

Let me compare. The config keys in `config.go` (the actual code, excluding `ASTRATE_TEST_DSN` which is test-only) are:

Code keys (from config.go):
- ASTRATE_DATABASE_DSN
- ASTRATE_ENGINE_SHARDS
- ASTRATE_HTTP_ADDR
- ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
- ASTRATE_HTTP_TLS_CERT_FILE
- ASTRATE_HTTP_TLS_KEY_FILE
- ASTRATE_LOG_FORMAT
- ASTRATE_LOG_LEVEL
- ASTRATE_MASTER_KEY
- ASTRATE_MASTER_KEY_FILE
- ASTRATE_MQTT_ADDR
- ASTRATE_MQTT_ADVERTISED_URL
- ASTRATE_MQTT_INSECURE_DEV_MODE
- ASTRATE_MQTT_SESSION_STORE_PATH
- ASTRATE_MQTT_TLS_CERT_FILE
- ASTRATE_MQTT_TLS_KEY_FILE
- ASTRATE_REALM_JWT_PUBLIC_KEY
- ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
- ASTRATE_REALM_NAME
- ASTRATE_SECURITY_MASTER_KEY_FILE

Docs keys:
- ASTRATE_DATABASE_DSN
- ASTRATE_ENGINE_SHARDS
- ASTRATE_HTTP_ADDR
- ASTRATE_HTTP_CORS_ALLOWED_ORIGINS
- ASTRATE_HTTP_TLS_CERT_FILE
- ASTRATE_HTTP_TLS_KEY_FILE
- ASTRATE_LOG_FORMAT
- ASTRATE_LOG_LEVEL
- ASTRATE_MASTER_KEY
- ASTRATE_MASTER_KEY_FILE
- ASTRATE_MQTT_ADDR
- ASTRATE_MQTT_ADVERTISED_URL
- ASTRATE_MQTT_INSECURE_DEV_MODE
- ASTRATE_MQTT_SESSION_STORE_PATH
- ASTRATE_MQTT_TLS_CERT_FILE
- ASTRATE_MQTT_TLS_KEY_FILE
- ASTRATE_REALM_JWT_PUBLIC_KEY
- ASTRATE_REALM_JWT_PUBLIC_KEY_FILE
- ASTRATE_REALM_NAME
- ASTRATE_SECURITY_MASTER_KEY_FILE

These are identical. The config keys are fully in sync.

Let me verify with a diff computation:
$ comm -3 <(rg -o '\bASTRATE_[A-Z_]+' internal/config/config.go | sort -u) <(rg -o '\bASTRATE_[A-Z_]+' docs/site/configuration-reference.md | sort -u); echo "exit: done (empty = identical)"
/bin/sh: 1: Syntax error: "(" unexpected

$ comm -3 <(rg -o '\bASTRATE_[A-Z_]+' internal/config/config.go | sort -u) <(rg -o '\bASTRATE_[A-Z_]+' docs/site/configuration-reference.md | sort -u); echo "exit done"
/bin/sh: 1: Syntax error: "(" unexpected

The `sh` shell doesn't support process substitution. Let me use temporary files.
$ rg -o '\bASTRATE_[A-Z_]+' internal/config/config.go | sort -u > /tmp/opencode/keys_code.txt; rg -o '\bASTRATE_[A-Z_]+' docs/site/configuration-reference.md | sort -u > /tmp/opencode/keys_docs.txt; diff /tmp/opencode/keys_code.txt /tmp/opencode/keys_docs.txt && echo "IDENTICAL"
IDENTICAL

The active config keys are fully in sync (identical sets). Let me verify one response detail on the astrate-native surface: the docs claim the native socket authenticates via a `?token=` query parameter and Authorization header. Let me confirm the middleware accepts that.
$ rg -n 'token|Authorization|Query\(\)' internal/auth/middleware.go | head -30
19:// tokenContextKey carries the verified *Token in the request context.
20:type tokenContextKey struct{}
22:// TokenFromContext returns the verified token stored by the middleware, if
25:	tok, ok := ctx.Value(tokenContextKey{}).(*Token)
31:// unverifiable token → 401, verified token whose claims do not authorize the
39:// DefaultCacheSize token cache.
46:// bearer token, and matches the claim's authorization strings against the
89:	tokenString, ok := bearerToken(r)
95:	tok, err := m.cache.Verify(tokenString, keysPEM)
113:	ctx := context.WithValue(r.Context(), tokenContextKey{}, tok)
134:// bearerToken extracts the credential from the Authorization header.
136:// "Bearer", optionally followed by a colon, then whitespace and the token.
138:	header := r.Header.Get("Authorization")
150:	token := strings.TrimSpace(rest)
151:	if token == "" {
154:	return token, true

The native socket endpoint uses the middleware (`mw.RequireRealm`), which only reads the Authorization header, not a `?token=` query param. Let me confirm the full `RequireRealm` function and check if there's any query-token support.
→ Read internal/auth/middleware.go

$ rg -n 'NewMetrics|Handler\(\)|promhttp|Handle' internal/observability/metrics.go 2>/dev/null | head; echo "---"; rg -rn 'Content-Type|text/plain' internal/observability/metrics.go 2>/dev/null | head
14:	"github.com/prometheus/client_golang/prometheus/promhttp"
25:// NewMetrics builds a registry preloaded with the Go runtime and process
27:func NewMetrics() *Metrics {
40:// Handler serves the registry in the Prometheus text exposition format.
41:func (m *Metrics) Handler() http.Handler {
42:	return promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})
---

Confirmed a real drift. The native socket endpoint is guarded by `mw.RequireRealm` whose `bearerToken` (middleware.go:137-155) only reads the Authorization header, yet docs/api/astrate_native_api.yaml:188-189 claims `?token=` is accepted. Let me verify the metrics content-type and finish the endpoint deep-checks for the astrate-native surface.
$ rg -n 'require|RequireRealm|ClaimChannels' internal/appengine/stream/ws.go internal/appengine/channels/ws.go
internal/appengine/stream/ws.go:36:	require func(http.Handler) http.Handler
internal/appengine/stream/ws.go:41:	return &API{bus: bus, require: mw.RequireRealm(auth.ClaimChannels)}
internal/appengine/stream/ws.go:46:	mux.Handle("GET /astrate/v1/{realm}/socket", a.require(http.HandlerFunc(a.handle)))

The metrics endpoint uses promhttp (text exposition format) — matches. I now have my findings for the astrate-native surface. Let me check the todo.md head to match its format conventions before append
