slug: recipe-hygiene
verdict: timeout
at:  e30b5fd
ran: 2026-10-08T20:06:24Z on DietPi in 1211s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md

$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ 2>/dev/null | head -30; echo "==="; go test ./... 2>&1 | rg -i 'skip|no test files' | head -20
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},
===
?   	github.com/astrate-platform/astrate/docs	[no test files]
?   	github.com/astrate-platform/astrate/migrations	[no test files]

$ go list -m -u all 2>/dev/null | rg '\[' | head -20; echo "---"; which govulncheck || echo "govulncheck not installed"; echo "---"; ls /Users/atsetilam/go/bin/golangci-lint 2>/dev/null || which golangci-lint || echo "golangci-lint not found"
cel.dev/expr v0.24.0 [v0.25.3]
cloud.google.com/go v0.121.6 [v0.123.0]
cloud.google.com/go/auth v0.16.4 [v0.24.0]
cloud.google.com/go/auth/oauth2adapt v0.2.8 [v0.3.0]
cloud.google.com/go/compute/metadata v0.8.0 [v0.10.0]
cloud.google.com/go/iam v1.5.2 [v1.14.0]
cloud.google.com/go/longrunning v0.6.7 [v1.3.0]
cloud.google.com/go/monitoring v1.24.2 [v1.31.0]
cloud.google.com/go/spanner v1.85.0 [v1.96.1]
cloud.google.com/go/storage v1.56.0 [v1.69.0]
github.com/99designs/keyring v1.2.1 [v1.2.2]
github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0 [v1.23.3]
github.com/Azure/azure-sdk-for-go/sdk/internal v1.1.2 [v1.13.0]
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0 [v1.8.2]
github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c [v0.0.0-20261006220739-8c912ac31dd4]
github.com/Azure/go-autorest/autorest/adal v0.9.16 (retracted) [v0.9.24] (deprecated)
github.com/Azure/go-autorest/autorest/date v0.3.0 [v0.3.1] (deprecated)
github.com/Azure/go-autorest/logger v0.2.1 [v0.2.2] (deprecated)
github.com/Azure/go-autorest/tracing v0.6.0 [v0.6.1] (deprecated)
github.com/ClickHouse/clickhouse-go v1.4.3 [v1.5.4]
---
govulncheck not installed
---
golangci-lint not found

$ ls; echo ---; ls cmd pkg 2>/dev/null; echo ---; ls ~/go/bin 2>/dev/null; ls ~/.local/go/bin 2>/dev/null; echo ---; cat .golangci.yml 2>/dev/null | head -60
AGENTS.md
Dockerfile
LICENSE
Makefile
README.md
astrate
bench
cmd
deploy
docker-compose.yml
docs
examples
flow-container-echo
go.mod
go.sum
internal
migrations
opencode.json
pkg
site-dist
test
tools
---
cmd:
astrate

pkg:
astarteapi
deviceid
interfaceschema
payload
---
goimports
golangci-lint
govulncheck
---
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

$ go test ./... 2>&1 | tail -40; echo "===SKIP==="; go test ./... -v 2>&1 | rg -n '--- SKIP' | head -30
ok  	github.com/astrate-platform/astrate/cmd/astrate	(cached)
?   	github.com/astrate-platform/astrate/docs	[no test files]
ok  	github.com/astrate-platform/astrate/examples/flow-container-echo	(cached)
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
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)
===SKIP===
rg: unrecognized flag --- SKIP

$ export PATH=$PATH:$HOME/go/bin; govulncheck ./... 2>&1 | tail -40
govulncheck: fetching vulnerabilities: Get "https://vuln.go.dev/ID/GO-2026-5841.json.gz": dial tcp: lookup vuln.go.dev: i/o timeout

$ export PATH=$PATH:$HOME/go/bin; go test ./... -v 2>&1 | rg -- '--- SKIP' | head -30; echo "=== exit"
=== exit

$ ls .mule/; echo ---; tail -60 .mule/todo.md; echo ---; ls .mule/recipes/
...output truncated...

Full output saved to: /root/.local/share/opencode/tool-output/tool_11d163ba3001mP8NeHmLcApGsL

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
- [!] pairing-burn-bcrypt-cost [auto]: make the dummy-secret timing burn use the configured BcryptCost — burnBcrypt compares against a package-level hash generated once at bcrypt.DefaultCost (internal/pairing/service.go:94-100, 428-430) while real auth hashes and compares at cfg.BcryptCost (service.go:204, 407), and bcrypt_cost is a documented, plumbed lever (internal/config/config.go:106, cmd/astrate/main.go:152, docs/site/configuration-reference.md:76, default 10), so with bcrypt_cost < 10 the unknown-device/empty-hash/malformed-ID failure paths each burn cost-10 work against cost-4 real compares, making failure responses ~an order of magnitude slower and a latency oracle for device existence — the inverse of the DESIGN §4.4 uniform-auth-time property. Keep one dummy hash per configured cost (generate lazily per Service at s.cfg.BcryptCost) and burn against it; add a deterministic test asserting bcrypt.Cost of the burn hash equals cfg.BcryptCost (e.g. a Service built with BcryptCost:4 must burn at cost 4), which fails today. — BLOCKED: lint failed: internal/pairing/service.go:95:7: G101: Potential hardcoded credentials (gosec)
- [x] pairing-onregistered-clock: route the registration event timestamp through the injectable clock — Register calls s.OnRegistered(realmName, hwID, time.Now()) with literal wall-clock time (internal/pairing/service.go:222) while every other timestamp in the package funnels through the s.now seam (VerifyCredentials, service.go:333), so a time-travel test cannot pin the DeviceRegisteredEvent timestamp fired via e.HandleDeviceRegistered (cmd/astrate/main.go:154). Change to s.now(), keep TestRegisterEmitsEvent.
- [x] pairing-unregister-bad-id-test: cover the malformed-device-IDs branch of Unregister in internal/pairing/service_test.go — Unregister wraps deviceid.Parse failure as store.ErrNotFound (internal/pairing/service.go:231-234) but TestUnregister (service_test.go:340) only exercises a well-formed unknown ID; assert both bad-ID shapes and the well-formed unknown-ID 404 in one table.
- [x] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.go
- [x] docs-native-ae-realm-version-endpoint [auto]: add the undocumented `GET /appengine/v1/{realm}/version` to docs/api/astarte_appengine_api.yaml — it is the only one of the three realm-scoped version routes no spec covers (pairing's at docs/api/astarte_pairing_api.yaml:70-96 and RM's at docs/api/astarte_realm_management_api.yaml:659-691 are both documented), and the native spec actively points readers at it while leaving it undocumented: docs/api/astrate_native_api.yaml:193 says "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints". Code: registered by mountRealmVersion at cmd/astrate/main.go:436-437 behind `mw.RequireRealmAny(auth.ClaimAppEngine)` (internal/auth/middleware.go:57-93), answering 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38-47) with the Astrate build `version` var (main.go:53, default "0.1.0-dev", ldflags-overridable) — the value docs-sync-native-version-value settled on, NOT `APICompatVersion` — plus 401 for a missing token or an unknown realm, 403 for a verified a_aea JWT that does not authorize it, 500 when GetRealmByName fails. Two details worth stating in the operation: the authorization path is exactly `version` (RelativePath drops the realm segment, middleware.go:143-151), so the JWT needs a grant matching GET + `version`, and an unknown realm is 401 not 404 (middleware.go:68-72) — do NOT document a 404, the hole docs-sync-pairing-version-404-unreachable closed on the pairing twin. `$ref` the existing `RealmName` parameter (yaml:1405) and the `Unauthorized`/`Forbidden`/`InternalServerError` responses (yaml:1734, 1744, 1784) that docs-sync-ae-forbidden-403 established, set `security: [a_aea]`, and inline the `{"data": "<string>"}` schema the way the pairing twin does (astarte_pairing_api.yaml:87-93). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-hk-version-description [auto]: fix the `/housekeeping/version` description in docs/api/astrate_native_api.yaml (yaml:249-253), which promises "The emulated upstream API level is served by the realm-scoped /v1/{realm}/version endpoints" for a path that does not exist — `mountRealmVersion` registers only appengine and pairing (cmd/astrate/main.go:435-438), realm-management's is served by its own API (internal/realm/http.go:52), and housekeeping's Mount registers exactly five realm-lifecycle routes with no version among them (internal/housekeeping/http.go:37-41), so `/housekeeping/v1/{realm}/version` falls through to httpx.NotFound (cmd/astrate/main.go:422). The same sentence appears on the other three compat-version ops and is correct on all of them (appengine main.go:436, realmmanagement realm/http.go:52, pairing main.go:438) — so change only the housekeeping one and say plainly that housekeeping is the one service with no realm-scoped twin. Note that the appengine pointer is true but currently dangling (that endpoint is undocumented; see docs-native-ae-realm-version-endpoint), so do not treat it as a model for correctness. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-metrics-content-negotiation [auto]: the `/astrate/v1/metrics` operation in docs/api/astrate_native_api.yaml (yaml:77-92) says "The response content type is `text/plain; version=0.0.4`" and uses that truncated string as the media-type key, but the handler is `promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{})` (internal/observability/metrics.go:42) and sets a longer header. Measured today through that exact handler with the repo's own client_golang v1.23.2 / common v0.66.1: with no Accept, and with `text/plain;version=0.0.4`, the Content-Type is `text/plain; version=0.0.4; charset=utf-8; escaping=underscores` (expfmt.FmtText plus the escaping parameter — prometheus/common expfmt/encode.go:63-91, `model.NameEscapingScheme = UnderscoreEscaping`); a scraper sending `Accept: application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited` gets a binary body (83 bytes for one series) under `application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited; escaping=underscores` that no response documents; and `Accept-Encoding: gzip` is honoured (`Content-Encoding: gzip`, promhttp defaultCompressionFormats at http.go:68-74). Do NOT claim zstd — it is not offered in this build. Correct the media-type key and the prose to the measured text/plain form, and either document the protobuf and gzip paths or state that the documented text form is what a Prometheus scraper gets. Every case measured is 200, so no new status code. The example body is already a real series thanks to docs-sync-native-metrics-example-fake-series — leave it alone. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-phoenix-newevent-payload [auto]: fix the `new_event` frame example in docs/api/astrate_native_api.yaml (yaml:346-349), which shows `[null, null, "rooms:test:all", "new_event", {"event_data": "..."}]` — there is no `event_data` key anywhere in the code and the payload is not a wrapper around one: `pumpEvents` marshals the bus event straight into the frame payload (`json.Marshal(ev)`, internal/appengine/channels/ws.go:406, frame at :410-416), and `stream.Event` carries no JSON tags (internal/engine/stream/bus.go:40-65), so encoding/json emits the Go field names verbatim. Measured on a verbatim copy of that struct: `{"Kind":"incoming_data","Realm":"test","DeviceID":"f0VMRgIBAQAAAAAAAAAAAA","Interface":"org.example.Sensors","Path":"/sensors/temperature/value","InterfaceMajor":1,"IP":"","ErrorName":"","ErrorMetadata":null,"Value":23.5,"Timestamp":"2024-05-30T13:49:57.045Z"}`. Replace the example with that shape and say in the description that this socket's payload keys are capitalised Go field names while the native socket above pushes snake_case — that one goes through `wireEvent` (internal/appengine/stream/ws.go:51-58), which is why the other example at yaml:283-293 is right and this one is not. Worth one sentence on the consequence: a client generated from `event_data` reads nothing at all, because Go's unmarshal leaves unknown fields zero rather than failing. The four-element frame shape itself is correct (protocol.go:27-43 pins five elements with null join_ref/ref) — leave that part. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-realm-name-pattern [auto]: constrain the realm name on both sockets in docs/api/astrate_native_api.yaml — `components.parameters.RealmName` (yaml:374-380, used by `/astrate/v1/{realm}/socket`) and the Phoenix `realm` query parameter (yaml:352-357) — both are bare `type: string` while the name is CHECK-constrained to `^[a-z][a-z0-9]*$` (migrations/000002_metadata.up.sql:7) and the housekeeping spec already carries that pattern on its twin parameters after docs-sync-hk-realm-name-pattern. Neither socket answers a bad name with a 404: both resolve the realm first and turn a miss into 401 — `GetRealmByName` → `store.ErrNotFound` → `WriteUnauthorized` at internal/appengine/stream/ws.go:46 → internal/auth/middleware.go:66-72 for the native one, and internal/appengine/channels/ws.go:56-69 for the Phoenix one — so a client sending `TEST` or `test_realm` is told its token is bad rather than that its realm name is malformed. Add `pattern: '^[a-z][a-z0-9]*$'` to both and say in each description that an off-pattern name is a 401, not a 404. Read migrations/ for the CHECK, never edit it (it is on the never-touch list). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
- [x] compat-note-v1.3.5 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.5 (newest stable, 2026-10-05, maintenance-only: mix.exs/mix.lock version bumps plus a single CHANGELOG "astarte_data_updater_plant RPC server reliability" entry, no wire/API surface change; v1.4.0 is still rc.6-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2/v1.3.4 wording proposal there already stands and its version-reference delta was recorded on 2026-10-05 as reading "newest stable v1.3.5", so fold this in as that one-word delta rather than re-deriving the proposal.
- [!] docs-native-socket-handshake-errors [auto]: add the two WebSocket-handshake failures to both socket operations in docs/api/astrate_native_api.yaml — `/astrate/v1/{realm}/socket` responses (yaml:335-345) and `/appengine/v1/socket/websocket` (yaml:398-404) document only 101/401/403/500 and 101/401/500 — but measured today with the repo's own dependency (probe in /tmp, coder/websocket accept.go `verifyClientRequest`): a plain GET with no `Connection: Upgrade` answers **426 Upgrade Required**, `Content-Type: text/plain; charset=utf-8`, body `WebSocket protocol violation: Connection header "" does not contain Upgrade\n`, and `Sec-WebSocket-Version: 14` answers **400** `unsupported WebSocket protocol version (only 13 is supported): "14"\n`. Add both to each operation and say two things explicitly: the bodies are plain text from `http.Error`, NOT the `{"errors":{"detail":...}}` envelope every other response on this spec uses (a client that JSON-parses them breaks), and they are only reached after authentication — native: the route's guard `mw.RequireRealm(auth.ClaimChannels)` wraps `handle` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112), Phoenix: the realm/token checks at internal/appengine/channels/ws.go:51-68 run before `websocket.Accept` at :71 — and only on the WebSocket branch (`wantsSSE` is tested first at ws.go:78), so 401/403 and the 200 SSE path win. Do NOT add 405: ServeMux's method-mismatch applies equally to every GET route in every spec. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: wrote nothing
- [x] docs-native-socket-sse-exact-accept [auto]: fix the SSE negotiation prose on `/astrate/v1/{realm}/socket` in docs/api/astrate_native_api.yaml — the description (yaml:289-291), the `transport` parameter (yaml:328-334) and the `200` response (yaml:338-339) read as if any client asking for `text/event-stream` gets SSE, but the predicate `wantsSSE` is exact string equality (internal/appengine/stream/ws.go:151-153): measured, `Accept: text/event-stream, text/plain` and `Accept: text/event-stream;q=1.0` both return false and fall through to `websocket.Accept`, where a non-upgrade client gets the 426 above; `transport` selects SSE only for the literal `sse`, and `transport=websocket` / `transport=` (empty) both mean WebSocket. Say "exactly" for the header, say any other `Accept` falls through to the upgrade attempt, and do NOT claim the server rejects out-of-enum `transport` values — the enum is a client-side contract the handler never enforces. While there, give the `200` its real media type: the handler sets `Content-Type: text/event-stream`, `Cache-Control: no-cache`, `Connection: keep-alive` and writes one `data: {...}` frame per event followed by a blank line (ws.go:123-126, 156-162). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [x] docs-native-socket-security-scheme [auto]: `/astrate/v1/{realm}/socket` is documented as if public in docs/api/astrate_native_api.yaml — the spec has no `components.securitySchemes` at all and root `security: []` (yaml:21), and the operation declares no `security`, so a client generated from this spec sends no `Authorization` header and is answered 401 `{"errors":{"detail":"Unauthorized"}}` by `mw.RequireRealm(auth.ClaimChannels)` (internal/appengine/stream/ws.go:46 -> internal/auth/middleware.go:48-112, bearer token read only from the Authorization header at middleware.go:131, 178-196). Add a scheme mirroring `a_aea` (docs/api/astarte_appengine_api.yaml:1436-1443: apiKey in header `Authorization`, described as a realm JWT carrying the `a_ch` claim) and set `security` on this operation only — root `security: []` stays so health/readiness/metrics/version remain documented as unauthenticated. If you model the Phoenix twin too, its credential is `?token=` in the query (internal/appengine/channels/ws.go:53-54), i.e. apiKey `in: query`, not header. Description notes worth carrying: the scheme is case-insensitive `Bearer` with an optional colon, and an unknown realm is 401 not 404 (middleware.go:68-72). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
- [!] docs-native-socket-event-schema [auto]: document the native socket's event payload in docs/api/astrate_native_api.yaml — currently only a fenced `incoming_data` example inside the operation description (yaml:300-310), with no `components.schemas` entry anywhere (yaml:425-478), and that example is not the general shape: `wireEvent` marks `interface`, `path` and `value` `omitempty` (internal/appengine/stream/ws.go:51-58), so measured on a verbatim copy a lifecycle event marshals to `{"event":"device_connected","realm":"test","device_id":"dev1","timestamp":"2024-05-30T13:49:57.045Z"}` — no interface/path/value at all (and `value: false`/`value: 0` ARE emitted, because omitempty on an `any` drops only nil). Add a `NativeEvent` schema (event, realm, device_id, optional interface/path/value, timestamp) referenced from the operation, state that `event` is one of the internal/engine/stream/bus.go Kind* constants (bus.go:20-26), and record that this socket drops `InterfaceMajor`, `IP`, `ErrorName` and `ErrorMetadata`, which `stream.Event` carries (bus.go:39-65) — so a `device_error` arrives here with no reason, unlike the Phoenix twin, which marshals the whole struct (internal/appengine/channels/ws.go:406-416). That last sentence documents a code gap, not a spec one: do not add the fields; if they are wanted on the wire, that is a separate code task. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: opencode exited 1
- [x] testutil-wait-since-cursor [auto]: give AstarteDevice's message waits a cursor so a second wait for the same topic cannot re-match an older message — `WaitForMessage`/`WaitForTopic` (internal/testutil/astartedevice.go:149-170) scan the capture buffer from index 0 on every poll, so a repeated wait returns the first match again, and in-tree that makes testE2EEmptyCache's resync wait vacuous (see the companion line). Move the paho capture handler out of the `ConnectAstarteDevice` closure (astartedevice.go:61-67) into an injectable method, add `Mark() int` and a `From(start int, …)` wait variant, and add a container-free internal/testutil/astartedevice_test.go that drives two captures on one topic and asserts a wait started after the second returns the second message — it cannot even compile today, which is the failure proof. Keep every new read of `received` under the existing `d.mu`, and say in the report that the change is unverified for races (the Pi gate runs no `-race`). The consumer of this API is engine-e2e-emptycache-resync-assert.
- [ ] engine-e2e-emptycache-resync-assert [legion] [auto]: make the second wait in `testE2EEmptyCache` (internal/engine/engine_e2e_test.go:367) wait for the resync re-publish instead of re-matching the retained set captured at :362 — mark the cursor before `dev.EmptyCache(t)` (:364) and take only messages after the mark, so the "resent property value = int32(42)" assertion at :368-369 fails when the resend never happens; today it cannot fail, because `WaitForTopic` returns the pre-EmptyCache message, which carries the same value 42. Needs `-tags e2e` plus Docker, so this half is Legion-only — the container-free helper and its test land in testutil-wait-since-cursor.
- [!] testutil-control-frame-test [auto]: pin the producer/properties control frame that `SendProducerProperties` puts on the wire (internal/testutil/astartedevice.go:181-216) — `DeflateControlList`/`InflateControlList` have zero tests and every caller of the helpers sits behind an `e2e`/`integration`/conformance build tag, so nothing on the default gate asserts the 4-byte big-endian uncompressed-length prefix plus zlib framing, while the parser side is pinned in internal/engine/control_test.go:89-130 against a CPython golden frame. Add a container-free internal/testutil/astartedevice_test.go with exact-bytes, round-trip, empty-list→nil, sub-4-byte header and declared≠inflated cases; bound the inflate with `io.LimitReader(zr, int64(declared)+1)` the way internal/engine/control.go:301-315 does (astartedevice.go:205 reads with unbounded `io.ReadAll` today, and a lying header from the code under test makes the harness allocate it before failing) and assert the truncated byte count in the mismatch message so the bound itself is pinned. One extra line in internal/engine/control_test.go feeding `testutil.DeflateControlList` into `inflateProperties` is the cross-check that closes the loop. — BLOCKED: wrote nothing
- [x] engine-introspection-producer-roundtrip [auto]: pin the producer half of the docs/DESIGN.md §3.3 introspection rule on the default gate — `testutil.Introspection` (internal/testutil/astartedevice.go:93-105) is called from six places, all behind `e2e`, `integration && e2e` or test/conformance build tags, while the consumer `parseIntrospection` (internal/engine/introspection.go:99-127) is tested only against hand-written literals in internal/engine/introspection_test.go:17-71, so nothing container-free checks that the string the harness renders is the one the engine accepts or that the documented "deterministic (sorted) ordering" holds. Add a case to internal/engine/introspection_test.go that builds the string from an unsorted multi-entry map via `testutil.Introspection`, asserts `parseIntrospection` returns the same map, and pins the exact output (sorted, `name:major:minor;…`, no trailing `;`, empty map → empty string).
- [x] docs-sync-pairing-register-404-unreachable [auto]: drop the `404` `DeviceNotFound` response from `registerDevice` in docs/api/astarte_pairing_api.yaml (yaml:140-141) — it is documented but absent: the route is wrapped by `requireAgent` (internal/pairing/http.go:82-83), i.e. `mw.RequireRealm(auth.ClaimPairing)`, which resolves the realm *before* the handler and answers `401 {"errors":{"detail":"Unauthorized"}}` for an unknown one (internal/auth/middleware.go:66-72, "no existence oracle on auth failures"), so `Service.Register`'s own `GetRealmByName` (internal/pairing/service.go:183-186) can return `store.ErrNotFound` only if the realm is deleted between the middleware's lookup and the handler's — a race, not a contract. Nothing else in Register returns `store.ErrNotFound`: `RegisterDevice`'s `ErrDeviceAlreadyConfirmed` becomes 422 (service.go:209-213) and the CountDevices/SetPayloadFormatHint failures are DB errors → 500. No test covers it (internal/pairing/http_test.go asserts 404 only for the agent DELETE at :527 and for health at :595-598), and the `DeviceNotFound` component's own description — "unknown device ID" — does not match what that path would be even if it fired. The premise came from docs-sync-pairing-register-404, which read service.go without the middleware in front of it. Keep the `404` on `unregisterDevice` (reachable: unknown device, and a malformed `deviceID` via service.go:231-234) and on `getHealth`. State the reachability argument in the report; if you find a reachable path instead, keep the response and fix its description rather than deleting it. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
- [!] docs-sync-pairing-body-cap-400 [auto]: document the 64 KiB request-body cap on the three body-taking pairing ops — `registerDevice`, `requestCredentials` and `verifyCredentials` all decode through `astarteapi.DecodeData(r.Body, maxBodyBytes, …)` with `maxBodyBytes = 64 << 10` (internal/pairing/http.go:17, 116, 169, 249), and any body over 65536 bytes fails with an error wrapping `ErrBodyTooLarge` (pkg/astarteapi/envelope.go:233-235) that each handler turns into a plain 400 `{"errors":{"detail":"Bad request"}}` (http.go:117-118, 170-171, 250-251). The shared `BadRequest` component (yaml:578-586) says only "malformed or missing data envelope" and no request-body schema carries a size hint, so a client is told nothing about the cap. Widen that one shared description to name the 64 KiB cap and the failure shape (a single edit covers all three operations) and say explicitly that an over-size body is 400, not 413 — do not invent a 413. One clause on `csr` (yaml:414) and `client_crt` (yaml:422) too: the cap is on the whole envelope, not on the PEM. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML). — BLOCKED: TIMEOUT after 1200s — task too big, split it
- [x] docs-sync-pairing-realm-name-pattern [auto]: constrain `components.parameters.RealmName` in docs/api/astarte_pairing_api.yaml (yaml:359-365, a bare `type: string` described only as "The realm name.") with `pattern: '^[a-z][a-z0-9]*$'` — the realm name is CHECK-constrained to exactly that (migrations/000002_metadata.up.sql:7, read-only, never edit) and the two surfaces whose tasks already ran carry it on their twins (docs/api/astarte_housekeeping_api.yaml:246, docs/api/astrate_native_api.yaml:420 and 467). Because this one parameter is shared by all seven operations, the description must record that an off-pattern name gets three different answers, not one: the agent and device routes answer **401** (the realm is resolved by the middleware, which treats a miss as unauthenticated rather than 404 — internal/auth/middleware.go:66-72); `getHealth` answers **404** (service.Health returns store.ErrNotFound → WriteNotFound, internal/pairing/service.go:369-371, http.go:289-290); and `getVersion` answers **200** `{"data":"1.2.2"}` for any realm string whatsoever, because `GET /pairing/v1/{realm}/version` is registered with no middleware and no realm lookup (cmd/astrate/main.go:438 → observability.VersionHandler, internal/observability/compat.go:38-47) — say that in the version operation so nobody reads the shared parameter as a guarantee the handler enforces. Mirror `TestHousekeepingRealmNamePatternDocumented` (internal/swagger/swagger_test.go:1247) with a pairing twin only if it stays a few lines; the YAML edit is the deliverable. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
- [x] docs-sync-pairing-deviceid-path-param [auto]: give `components.parameters.DeviceID` the constraints `RegisterRequest.hw_id` already carries in docs/api/astarte_pairing_api.yaml — the parameter (yaml:367-377) is prose-only, while `hw_id` (yaml:386-398) got `pattern: '^[A-Za-z0-9_-]{22}$'` plus `minLength`/`maxLength` from docs-sync-pairing-deviceid-base64url, and `TestPairingDeviceIDEncodingDocumented` (internal/swagger/swagger_test.go:1086) pins only the two descriptions and hw_id's bounds, never the parameter's schema; the realm-management twin was constrained later by docs-sync-rm-deviceid-param and its test (swagger_test.go:1168) does hold the parameter's pattern against `deviceid.Parse`. Add the same three schema lines to the parameter, and record in its description the asymmetry the prose omits: a malformed `deviceID` is **401** on all three device-authenticated operations — `deviceid.Parse` failing inside `authenticateDevice` burns bcrypt and returns `ErrUnauthorized` (internal/pairing/service.go:381-385) → `WriteUnauthorized` (http.go:322-323), exactly as an unknown device or a wrong secret is (service.go:395-397), so the client is told its credentials are bad rather than its ID malformed — but **404** on the agent `unregisterDevice` (service.go:231-234 → http.go:326-327). Extend TestPairingDeviceIDEncodingDocumented to compile the new pattern and hold it against `deviceid.Parse` the way the realm-management test does. Do not encode the strict trailing-bits rule in the pattern; that stays open in deviceid-strict-trailing-bits. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads (mkdocs may be missing — `pip install -r docs/requirements.txt` first, or at minimum confirm the file parses as YAML).
---
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

$ export PATH=$PATH:$HOME/go/bin; golangci-lint run ./... 2>&1 | tail -40; echo "EXIT=$?"
0 issues.
EXIT=0

$ sed -n '40,80p' internal/broker/intake.go; echo ===; sed -n '1,40p' .mule/todo.md
// datastream timestamp when the payload carries none).
	ReceivedAt time.Time
	// Ack releases the broker-held acknowledgment for this message. For
	// QoS >= 1 the device's PUBACK (or PUBREC) is not sent until Ack is
	// called — the engine calls it after the persistence batch commits
	// (docs/DESIGN.md §5.3). Ack is never nil, is safe to call multiple
	// times, and must eventually be called for every QoS >= 1 message the
	// intake accepts, or the publishing device stalls (that stall is the
	// designed backpressure path, §1.4).
	Ack func()
}

// Intake consumes accepted device publishes. internal/engine implements it
// (M6); tests use recorders. Submit may block — a full engine shard blocking
// the broker's per-client read loop is the §1.4 backpressure contract.
//
// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):
// external-bus intake — a second implementation of this interface backed by a
// durable external bus (e.g. NATS JetStream) for multi-instance deployment or
// restart survival. Not new device data sources: devices keep publishing over
// the embedded MQTT broker. The frozen decision (2026-08-22) is in §1.4; the
// consumer must reproduce per-device ordering, deferred-ack backpressure, and
// QoS 0 drop semantics.
type Intake interface {
	Submit(InboundMessage)
}

// LifecycleEventType discriminates LifecycleEvent values. The constants
// double as the Astarte trigger event names fed to the M6 trigger engine.
type LifecycleEventType string

const (
	// EventDeviceConnected fires after a device session is established.
	EventDeviceConnected LifecycleEventType = "device_connected"
	// EventDeviceDisconnected fires after a device connection ends (not on
	// session takeover: the device is still connected, on a new channel).
	EventDeviceDisconnected LifecycleEventType = "device_disconnected"
)

// LifecycleEvent is one device connect/disconnect observation
// (docs/DESIGN.md §3.1 — the work astarte_vmq_plugin does upstream).
===
- [x] errorname-missing-required-test [auto]: add the missing `missing_required` → `unexpected_object_key` row to `TestUpstreamErrorNameMapping` (internal/engine/triggers/errorname_test.go) — it is the only `astrateToUpstream` key (errorname.go:64) with no table row, so a regression remapping or dropping this upstream-1.4 entry passes every suite; also add an invariant that every `astrateToUpstream` value is a member of `UpstreamErrorNames()` to catch typo-then-copied mapping values.
- [x] triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
- [x] compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).

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

$ export PATH=$PATH:$HOME/go/bin; for d in $(go list ./... ); do dir=${d#github.com/astrate-platform/astrate/}; [ -z "$dir" ] && continue; if ! ls $dir/*_test.go >/dev/null 2>&1; then echo "NOTEST $dir"; fi; done
NOTEST docs
NOTEST migrations

$ cat .mule/for-giulio.md | tail -80; echo "=== LINES: $(wc -l < .mule/for-giulio.md)"
verified context).
- **The mule has been idle 17h.** Filed by the dead-man's switch; see journalctl on the Pi.

## 2026-09-25 — the `-tags nats` half of the tree is compiled by no gate

`cmd/astrate/newnats_nats.go`, `internal/engine/forward/nats.go` and
`internal/engine/forward/nats_test.go` are all behind `//go:build nats`, and nothing
builds or tests that tag: not `make build`/`make test`/`make test-integration`/
`make test-e2e` (Makefile:23, 38, 42, 46), not the mule gate, not the lint target.
So the whole NATS trigger-forwarding path is invisible to CI — while
`config.validate` accepts `triggers.forward.kind = "nats"` (internal/config/config.go:339-343)
and the untagged binary then fails at boot with the build-tag message
(cmd/astrate/newnats_default.go). I measured it today: `go build -tags nats ./cmd/...`
and `go vet -tags nats ./cmd/...` both exit 0, so there is nothing broken — just nothing
keeping it unbroken. Worth a `make build-nats` (and the NATS test tier, which needs a
container, so it would be a `[legion]` line) in the gate; not a queue line because the
only test I can imagine for it shells out to `go build`, and the mule's own check strips
the implementation and the guard together, so it would never fail.
Not queued for the same reason: `IdleTimeout` is missing on the HTTP server
(cmd/astrate/main.go:221) and `-healthcheck` ignores the config file
(cmd/astrate/main.go:531) — both real, both yours to call, neither has a test worth
writing. See `.mule/reviews/cmd-astrate-2026-09-25.md`.
- **The mule has been idle 16h.** Filed by the dead-man's switch; see journalctl on the Pi.

## 2026-10-04 — three housekeeping drifts in `docs/site/` the recipe is not allowed to touch

Ran the docs-sync recipe over the housekeeping surface (code routes vs
`docs/api/astarte_housekeeping_api.yaml`). The spec half is queued as four lines in
`.mule/todo.md`; the prose half needs you, because `docs/site/` is on the never-touch list.

**1. `docs/site/housekeeping-api.md:30` still talks about Cassandra.**

> Cassandra-specific fields (`replication_class`, `replication_factor`, etc.) from upstream are accepted but ignored.

Astrate stores PostgreSQL — `internal/store` on pgx, schema in `migrations/*.sql`; there is no
Cassandra anywhere in the tree. The wire shape is four fields (`internal/housekeeping/http.go:46-51`),
and the comment right above it says the opposite of "accepted": *"Astrate omits the
Cassandra-specific fields (replication factor/class) upstream carries"*
(`internal/housekeeping/http.go:44-45`). From a client's side the two readings are the same
(the request succeeds, the fields do nothing), so this is vocabulary rather than behaviour —
but it is the kind of sentence that sends an operator hunting for a Cassandra setting that
cannot exist.

**2. `PATCH /housekeeping/v1/realms/{realm}` is missing from the page entirely.**
`docs/site/housekeeping-api.md` goes Create → List → Get → Delete (lines 10-50), so the route
`docs-sync-hk-patch-endpoint` added and documented in the spec
(`docs/api/astarte_housekeeping_api.yaml:142-192`) has no prose at all. The same page also
predates `datastream_maximum_storage_retention`, which the spec has carried on both `Realm` and
`RealmCreate` since `docs-sync-hk-retention-field`, so its create example (lines 12-21) shows a
three-field body that is not the whole body.

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
=== LINES: 657

$ rg -n -i 'hygiene|govuln|golangci|lint|no test files|external-bus|NATS' .mule/todo.md | head -30; echo ===; rg -n 't\.Skip|Skipf|SkipNow' --glob '*_test.go' -g '!**/vendor/**' . | head -20
65:`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
86:- [ ] broker-external-bus-intake [legion]: implement the second Intake implementation named in the TODO at internal/broker/intake.go:56-62 — an external-bus-backed (e.g. NATS JetStream) Intake for multi-instance deployment / restart survival, reproducing the embedded broker's per-device ordering, deferred-ack backpressure and QoS 0 drop semantics; an implementation of this interface needs containerised integration tests, so it runs on the Legion Go.
88:- [ ] hygiene-govulncheck [legion]: on the Legion Go, run govulncheck ./... on a fresh ~/astrate clone and file a task line for any REACHABLE advisory it reports (name the CVE and the call path, per the hygiene recipe); this box cannot build it — go install golang.org/x/vuln/cmd/govulncheck@latest is OOM-killed here (3.7GB, 4 cores), so the recipe's highest-priority check never runs unless done there.
135:- [!] appengine-group-token-roundtrip-test [auto]: add a container-free unit test for `parseGroupToken`/`groupTokenFor` (internal/appengine/service.go:602-623) — currently only exercised via the pg-backed groups_parity_test.go. Assert (1) round-trip `groupTokenFor(offset)` → `parseGroupToken` recovers offset, (2) a non-v1-nibble UUID is rejected, (3) a well-formed v1 UUID with zeroed time fields is accepted (parseGroupToken comment, service.go:614-616). — BLOCKED: lint failed: internal/appengine/groups_token_test.go:14:5: redefines-builtin-id: redefinition of the built-in function max (revive)
213:- [!] flow-randomsource-span-overflow [auto]: `randomSource.next()` computes `rand.Int64N(s.maxInt-s.minInt+1)` (internal/flow/blocks/randomsource.go:171) whose argument overflows int64 to a non-positive value when min/max span more than the int64 range — config like `min:-9e18, max:9e18` passes the single `min <= max` check (randomsource.go:106-108) but makes `rand.Int64N` panic (it panics on n <= 0) inside the source pump goroutine (flow.go:205-239, no recover), crashing the whole process. Reject a span wider than `math.MaxInt64` at construction alongside the existing check, and add a construct/Emit test with the wide span. Container-free. — BLOCKED: lint failed: internal/flow/blocks/randomsource.go:113:40: G115: integer overflow conversion int64 -> uint64 (gosec)
234:- [x] forward-envelope-bytes-test [auto]: `marshalEnvelope` (internal/engine/forward/envelope.go:19-29) has no direct test — `rg marshalEnvelope` finds exactly two callers (http.go:99, nats.go:45) and no test file — while envelope.go:6-8 makes the strong claim that "every Forwarder implementation in this package must produce byte-identical envelopes for the same inputs". Both suites check the shape only by unmarshalling into `bodyShape` (http_test.go:24-29), which cannot see key order or the exact bytes and would pass unchanged if a field were renamed on both sides. Add internal/engine/forward/envelope_test.go with a table pinning the exact output for (nil, nil) → `{"realm":"r","trigger":"t","action":null,"event":null}`, for the empty-non-nil pair the rule at 16-18 is about, and for a valued pair; container-free, no HTTP server needed.
253:- [!] container-response-cap-test [auto]: pin `Bridge.RoundTrip`'s 1 MiB response cap (internal/flow/blocks/container/httpbridge.go:133-140 — `io.LimitReader(resp.Body, maxBody+1)` plus the explicit rejection; lines 136, 138.16 and 138.40 are all uncovered), the only bound between a tenant-supplied container and process memory on the message path: assert a handler answering exactly `maxBody+1` bytes is rejected with the "exceeds" error and that exactly `maxBody` bytes still round-trip, driving a `Bridge{MaxBodyBytes: N}` so the fixture stays small. Also worth one line in the same change: the cap is absent from the operator-facing catalog entry (internal/flow/blocks/info.go:116) and from the schema, so an author whose container emits 2 MiB learns about it only from a runtime error. Container-free. — BLOCKED: lint failed: internal/flow/blocks/container/block_test.go:500:2: redefines-builtin-id: redefinition of the built-in function cap (rev
256:- [x] docs-sync-rm-validationerror-example [auto]: fix or split the shared `ValidationError` response component in docs/api/astarte_realm_management_api.yaml (yaml:1330-1338) — it is `$ref`d by exactly three operations, `putAuthConfig` (yaml:574), `createPolicy` (yaml:782) and `deletePolicy` (yaml:839), none of which is an interface surface, so its single example `detail: "invalid interface: missing \"interface_name\""` is wrong twice: the real details are `jwt_public_key_pem can't be blank` (internal/realm/service.go:632-635), `triggers: policy must have at least one error handler` and `triggers: policy name must be 1-128 characters of [a-zA-Z0-9_.~-]` (internal/engine/triggers/policy.go:112-117, wrapped by validatePolicy at internal/realm/policy.go:11-16) and `policy "X" is still used by trigger "Y"` (service.go:568-576) — all rendered by `validationDetail`, which strips the `realm: validation failed: ` prefix from `ErrValidation` (service.go:29, http.go:562-568); and that exact string is not producible at all, because a missing required field returns the bare `missing "interface_name"` (pkg/interfaceschema/parse.go:238-240, no `ErrInvalid` wrap) while the `invalid interface: ` prefix belongs only to the violations renderer (pkg/interfaceschema/violations.go:65-68). OpenAPI here is 3.0.3, so an example next to a `$ref` is ignored — give each of the three operations its own inline 422 with the real detail (the pattern installInterface and updateInterface already use at yaml:138-168 and 274-304), or add one component per surface. Do not touch the two interface ops, whose messages are correct. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
258:- [!] docs-sync-rm-update-interface-body [auto]: give the `updateInterface` request body a schema in docs/api/astarte_realm_management_api.yaml — `data` is declared as a bare `type: object` with no properties (yaml:250-261) although the handler parses it with the same `interfaceschema.ParseInterfaceCanonical` as `installInterface` (internal/realm/service.go:196 vs 161) and then checks the body's `interface_name`/`version_major` against the URL's `{name}`/`{major}` (service.go:203-208), so the update body — the one carrying every interface rule — is the only interface body in the spec with nothing to validate against. `$ref` `#/components/schemas/InterfaceDefinition` as `installInterface` does (yaml:115-117), and say in the operation description that the body must name the same interface and major as the URL (disagreement is 409, already documented) and that `version_minor` must increase (409). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads. — BLOCKED: TIMEOUT after 1201s — task too big, split it
276:- [!] pairing-burn-bcrypt-cost [auto]: make the dummy-secret timing burn use the configured BcryptCost — burnBcrypt compares against a package-level hash generated once at bcrypt.DefaultCost (internal/pairing/service.go:94-100, 428-430) while real auth hashes and compares at cfg.BcryptCost (service.go:204, 407), and bcrypt_cost is a documented, plumbed lever (internal/config/config.go:106, cmd/astrate/main.go:152, docs/site/configuration-reference.md:76, default 10), so with bcrypt_cost < 10 the unknown-device/empty-hash/malformed-ID failure paths each burn cost-10 work against cost-4 real compares, making failure responses ~an order of magnitude slower and a latency oracle for device existence — the inverse of the DESIGN §4.4 uniform-auth-time property. Keep one dummy hash per configured cost (generate lazily per Service at s.cfg.BcryptCost) and burn against it; add a deterministic test asserting bcrypt.Cost of the burn hash equals cfg.BcryptCost (e.g. a Service built with BcryptCost:4 must burn at cost 4), which fails today. — BLOCKED: lint failed: internal/pairing/service.go:95:7: G101: Potential hardcoded credentials (gosec)
===
./cmd/astrate/main_test.go:47:	t.Skip("the boot suite needs a database: run `make up` or set ASTRATE_TEST_DSN")
./cmd/astrate/forward_test.go:68:		t.Skip("built with -tags nats; see forward_nats_test.go")
./test/conformance/cpb/cpb_test.go:117:	t.Skip("CP-B needs a database: run `make up` at the repo root or set ASTRATE_TEST_DSN")
./test/conformance/pysdk/pysdk_test.go:47:		t.Skipf("pysdk: %q not found (set ASTRATE_PYSDK_PYTHON)", py)
./test/conformance/pysdk/pysdk_test.go:50:		t.Skipf("pysdk: astarte-device-sdk not importable by %s — install test/conformance/pysdk/requirements.txt (%s)", py, out)
./test/conformance/cpa/astartectl_test.go:46:		t.Skipf("astartectl: %v (set ASTARTECTL_BIN to run this checkpoint)", err)
./test/conformance/cpc/astartectl_test.go:47:		t.Skipf("astartectl: %v (set ASTARTECTL_BIN to run this checkpoint)", err)
./test/conformance/cpa/cpa_test.go:132:	t.Skip("CP-A needs a database: run `make up` at the repo root or set ASTRATE_TEST_DSN")
./internal/store/datastreams_test.go:420:			t.Skip("timescaledb_toolkit not installed; lttb path unavailable")
./internal/pairing/http_test.go:629:		t.Skip("openssl not found in PATH; chain already verified with crypto/x509")
./internal/pairing/ca/ca_test.go:293:		t.Skip("10k issuance draw skipped in -short mode")
./test/conformance/cpc/cpc_test.go:122:	t.Skip("CP-C needs a database: run `make up` at the repo root or set ASTRATE_TEST_DSN")

$ export PATH=$PATH:$HOME/go/bin; go list -m -f '{{.Path}} {{.Version}}{{if .Update}} [{{.Update.Version}}]{{end}}{{if .Indirect}} (indirect){{end}}' all 2>/dev/null | rg -v '\(indirect\)' | rg '\[' ; echo "=== retracted:"; go list -m -f '{{.Path}} {{.Version}}{{if .Retracted}} RETRACTED: {{.Retracted}}{{end}}' all 2>/dev/null | rg RETRACTED
=== retracted:

$ export PATH=$PATH:$HOME/go/bin; go test -count=1 -v ./cmd/astrate/ ./internal/pairing/ ./internal/pairing/ca/ ./internal/store/ 2>&1 | rg -- '--- (SKIP|FAIL)' | head -30; echo "EXIT=$?"
EXIT=0

$ sed -n '35,55p' cmd/astrate/main_test.go; echo "=== forward_test"; sed -n '55,75p' cmd/astrate/forward_test.go; echo "=== datastreams 405-430"; sed -n '405,430p' internal/store/datastreams_test.go; echo "=== pairing http 615-640"; sed -n '615,640p' internal/pairing/http_test.go; echo "=== ca 285-300"; sed -n '285,300p' internal/pairing/ca/ca_test.go
for _, dsn := range []string{envOr("ASTRATE_TEST_DSN", ""), composeDSN} {
		if dsn == "" {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		conn, err := pgx.Connect(ctx, dsn)
		cancel()
		if err == nil {
			_ = conn.Close(context.Background())
			return dsn
		}
	}
	t.Skip("the boot suite needs a database: run `make up` or set ASTRATE_TEST_DSN")
	return ""
}

func envOr(k, def string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return def
=== forward_test
func TestForwarderRejectsBadURL(t *testing.T) {
	var cfg config.Config
	cfg.Triggers.Forward = config.ForwardConfig{Kind: "http", URL: "nope"}
	if _, err := newForwarder(cfg, quietLogger()); err == nil {
		t.Fatal("expected an error for a relative URL")
	}
}

// TestForwarderNATSWithoutBuildTag: this file (and newnats_default.go) build
// without -tags nats, so kind = "nats" must fail at boot with a message
// naming the missing build tag, not silently forward nothing.
func TestForwarderNATSWithoutBuildTag(t *testing.T) {
	if natsBuildTagEnabled {
		t.Skip("built with -tags nats; see forward_nats_test.go")
	}
	var cfg config.Config
	cfg.Triggers.Forward = config.ForwardConfig{Kind: "nats", URL: "nats://bus.example:4222", Subject: "astrate.triggers"}
	_, err := newForwarder(cfg, quietLogger())
	if err == nil {
		t.Fatal("expected an error building a NATS forwarder without -tags nats")
	}
=== datastreams 405-430
		}, 5*time.Minute)
		if err != nil {
			t.Fatalf("Downsample integers: %v", err)
		}
		if len(intBuckets) != 1 || intBuckets[0].Value != 15 { // avg(0,10,20,30)
			t.Errorf("integer downsample: %v", intBuckets)
		}

		if _, err := s.Downsample(ctx, SeriesQuery{}, 0); err == nil {
			t.Error("zero bucket accepted")
		}
	})

	t.Run("DownsampleLTTB", func(t *testing.T) {
		if !s.HasToolkitLTTB() {
			t.Skip("timescaledb_toolkit not installed; lttb path unavailable")
		}
		realm := mustCreateRealm(t, s)
		device := mustRegisterDevice(t, s, realm.ID)
		si := mustInstallInterface(t, s, realm.ID, allTypesDef)

		// A flat series with one sharp spike. This is the whole reason lttb
		// exists: time_bucket+avg would dilute the spike into its bucket's
		// mean and the shape would be lost, whereas lttb selects the real
		// sample. 60 points, spike of 1000 at index 30.
		base := time.Date(2026, 1, 18, 12, 0, 0, 0, time.UTC)
=== pairing http 615-640
	t.Helper()
	id, err := deviceid.Random()
	if err != nil {
		t.Fatal(err)
	}
	return id.String()
}

// opensslVerify shells out to `openssl verify` as an independent
// implementation cross-check of the issued chain.
func opensslVerify(t *testing.T, caPEM, certPEM string) {
	t.Helper()
	openssl, err := exec.LookPath("openssl")
	if err != nil {
		t.Skip("openssl not found in PATH; chain already verified with crypto/x509")
	}

	dir := t.TempDir()
	caPath := filepath.Join(dir, "ca.pem")
	certPath := filepath.Join(dir, "client.pem")
	if err := os.WriteFile(caPath, []byte(caPEM), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(certPath, []byte(certPEM), 0o600); err != nil {
		t.Fatal(err)
	}
=== ca 285-300
		if _, err := c.Verify("not a certificate", time.Now()); !errors.Is(err, ErrCertificateInvalid) {
			t.Errorf("garbage: got %v, want ErrCertificateInvalid", err)
		}
	})
}

func TestSerialUniqueness10k(t *testing.T) {
	if testing.Short() {
		t.Skip("10k issuance draw skipped in -short mode")
	}
	c := mustGenerate(t, "test", 0)
	csr := newCSR(t, newECKey(t)) // key reuse is fine; serials must still differ

	seen := make(map[string]struct{}, 10_000)
	for i := 0; i < 10_000; i++ {
		_, serial, _, err := c.SignCSR(csr, "test", "dev", time.Hour)

$ head -20 test/conformance/cpa/cpa_test.go; echo ===; rg -n 'go:build' test/conformance/*/*.go | head; echo ===; rg -n 'hygiene' .mule/todo.md
// Package cpa is conformance checkpoint CP-A (docs/ROADMAP.md §0.3, §5):
// the wire-frozen proof that Astrate's Pairing API serves *unmodified*
// official Astarte clients. It boots a real store + auth + pairing stack
// against a live TimescaleDB and drives it with:
//
//   - the official astarte-device-sdk-go pairing client (registration via
//     astarte-go, then the SDK device's own CSR/credentials handshake),
//     asserting the SDK obtains a certificate it considers valid (it
//     persists it and does not re-request on the next start);
//   - the pinned astartectl release binary
//     (`astartectl pairing agent register`).
//
// The database is the compose service (`make up`) or ASTRATE_TEST_DSN; the
// suite skips with instructions when neither is reachable.
package cpa

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
===
===
65:`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
88:- [ ] hygiene-govulncheck [legion]: on the Legion Go, run govulncheck ./... on a fresh ~/astrate clone and file a task line for any REACHABLE advisory it reports (name the CVE and the call path, per the hygiene recipe); this box cannot build it — go install golang.org/x/vuln/cmd/govulncheck@latest is OOM-killed here (3.7GB, 4 cores), so the recipe's highest-priority check never runs unless done there.

$ export PATH=$PATH:$HOME/go/bin; go test -count=1 -v ./internal/store/ 2>&1 | rg -- 'SKIP|FAIL|^ok|^---' | head -40; echo "---cmd---"; go test -count=1 -v ./cmd/astrate/ 2>&1 | rg -- 'SKIP|FAIL|^ok' | head -20
--- PASS: TestKeySealerRoundTrip (0.00s)
--- PASS: TestKeySealerNonceUniqueness (0.00s)
--- PASS: TestKeySealerWrongKeyFails (0.00s)
--- PASS: TestKeySealerTamperDetected (0.00s)
--- PASS: TestKeySealerShortInput (0.00s)
--- PASS: TestNewKeySealerRejectsBadLength (0.00s)
--- PASS: TestLoadMasterKeyFromEnvHex (0.00s)
--- PASS: TestLoadMasterKeyFromEnvBase64 (0.00s)
--- PASS: TestLoadMasterKeyFromFile (0.00s)
--- PASS: TestLoadMasterKeyErrors (0.00s)
--- PASS: TestValidatePipelineGraph_Acyclic (0.00s)
--- PASS: TestValidatePipelineGraph_Cyclic (0.00s)
--- PASS: TestValidatePipelineGraph_MalformedJSON (0.00s)
--- PASS: TestValidatePipelineGraph_InvalidBlockRef (0.00s)
--- PASS: TestValidatePipelineGraph_DuplicateBlockName (0.00s)
--- PASS: TestValidatePipelineGraph_NoBlocks (0.00s)
--- PASS: TestValidatePipelineGraph_EmptyBlockName (0.00s)
ok  	github.com/astrate-platform/astrate/internal/store	0.012s
---cmd---
ok  	github.com/astrate-platform/astrate/cmd/astrate	0.328s

$ head -5 cmd/astrate/main_test.go; echo ===; rg -n 'go:build' cmd/astrate/*.go internal/store/*.go internal/pairing/*_test.go internal/pairing/ca/*_test.go test/conformance/*/*.go | head -30; echo === ; rg -n 'HasToolkitLTTB|func TestDownsample' internal/store/*_test.go | head
//go:build integration

package main

import (
===
internal/store/notify_test.go:1://go:build integration
internal/store/userblocks_test.go:1://go:build integration
internal/store/groups_test.go:1://go:build integration
internal/store/store_test.go:1://go:build integration
internal/store/devices_test.go:1://go:build integration
internal/store/realms_test.go:1://go:build integration
internal/store/policies_test.go:1://go:build integration
cmd/astrate/main_test.go:1://go:build integration
cmd/astrate/newnats_default.go:1://go:build !nats
cmd/astrate/forward_nats_test.go:1://go:build nats
internal/store/pipelines_test.go:1://go:build integration
internal/store/compression_test.go:1://go:build integration
cmd/astrate/newnats_nats.go:1://go:build nats
internal/store/interfaces_test.go:1://go:build integration
internal/store/flows_test.go:1://go:build integration
internal/store/datastreams_test.go:1://go:build integration
internal/store/retention_test.go:1://go:build integration
internal/pairing/http_test.go:1://go:build integration
internal/store/properties_test.go:1://go:build integration
internal/store/triggers_test.go:1://go:build integration
===
internal/store/store_test.go:45:	t.Logf("timescaledb_toolkit available: %v", s.HasToolkitLTTB())
internal/store/datastreams_test.go:419:		if !s.HasToolkitLTTB() {

$ rg -n 'func Test' internal/store/*_test.go | wc -l; rg -l 'integration' internal/store/*_test.go; echo ===; ls internal/store/
19
internal/store/userblocks_test.go
internal/store/realms_test.go
internal/store/pipelines_test.go
internal/store/flows_test.go
internal/store/compression_test.go
internal/store/groups_test.go
internal/store/retention_test.go
internal/store/interfaces_test.go
internal/store/datastreams_test.go
internal/store/policies_test.go
internal/store/store_test.go
internal/store/notify_test.go
internal/store/devices_test.go
internal/store/properties_test.go
internal/store/triggers_test.go
===
compression_test.go
crypto.go
crypto_test.go
datastreams.go
datastreams_test.go
devices.go
devices_test.go
flows.go
flows_test.go
groups.go
groups_test.go
interfaces.go
interfaces_test.go
notify.go
notify_test.go
pipelines.go
pipelines_test.go
pipelines_validate_test.go
policies.go
policies_test.go
properties.go
properties_test.go
realms.go
realms_test.go
retention.go
retention_test.go
store.go
store_test.go
testdata
triggers.go
triggers_test.go
userblocks.go
userblocks_test.go

$ rg -n 'TODO|FIXME|XXX|HACK' --glob '*.go' --glob '!*_test.go' bench/ tools/ examples/ test/ 2>/dev/null | rg -v 'context\.TODO' | head -20; echo "=== all non-test TODO markers:"; rg -n 'TODO\(|FIXME|XXX|HACK' --glob '*.go' . | rg -v 'context\.TODO|_test\.go' | head -20
=== all non-test TODO markers:
./internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):

$ export PATH=$PATH:$HOME/go/bin; go list ./... | rg 'conformance|bench|tools|examples|test'; echo "=== dirs with .go not in go list:"; for f in $(rg -l --glob '*.go' --glob '!*_test.go' '^package ' . | sed 's|^\./||'); do d=$(dirname $f); echo "$d"; done | sort -u | while read d; do go list ./... 2>/dev/null | rg -q "astrate/$d\$" || echo "NOTLISTED $d"; done
github.com/astrate-platform/astrate/examples/flow-container-echo
github.com/astrate-platform/astrate/internal/testutil
=== dirs with .go not in go list:
NOTLISTED bench
NOTLISTED test/conformance/instance
NOTLISTED test/conformance/upstream/record
NOTLISTED test/conformance/upstream/recordchannels
NOTLISTED tools/bsoncapture

$ for d in bench tools/bsoncapture test/conformance/instance test/conformance/upstream/record test/conformance/upstream/recordchannels; do echo "== $d"; ls $d; done; echo "=== build tags:"; rg -n 'go:build' bench/*.go tools/bsoncapture/*.go test/conformance/instance/*.go test/conformance/upstream/record*/*.go 2>/dev/null | head -30
== bench
README.md
astrate
client.go
client_test.go
connstorm.go
device.go
go.mod
go.sum
ingest.go
interfaces
main.go
provision.go
query.go
scripts
state.go
stats.go
== tools/bsoncapture
go.mod
go.sum
main.go
== test/conformance/instance
instance.go
== test/conformance/upstream/record
main.go
== test/conformance/upstream/recordchannels
main.go
=== build tags:

$ head -20 bench/go.mod; echo ===; rg -n 'bench' Makefile | head -20; echo ===; ls .github/workflows 2>/dev/null; rg -n 'go test|golangci|bench' Makefile | head -30
// Standalone benchmark harness (see README.md). Deliberately a separate
// module that imports NOTHING from astrate: it speaks only the wire protocols
// (Astarte MQTT v1 + the REST APIs), so the same binary drives an upstream
// Astarte deployment and an Astrate deployment interchangeably.
module github.com/astrate-platform/astrate/bench

go 1.26.1

require (
	github.com/eclipse/paho.mqtt.golang v1.5.1 // same pin as the main module
	github.com/golang-jwt/jwt/v5 v5.3.1
	go.mongodb.org/mongo-driver/v2 v2.6.0
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
)
===
===
ci.yml
docs.yml
4:#   make lint                golangci-lint (config: .golangci.yml)
30:## lint: run the pinned golangci-lint (install with `make tools`).
32:	@command -v golangci-lint >/dev/null 2>&1 || { \
33:		echo "golangci-lint not found; run 'make tools' (installs $(GOLANGCI_LINT_VERSION))"; exit 1; }
34:	golangci-lint run ./...
58:	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

$ go list ./... | rg 'conformance|instance|record'; echo "rc=$?"; head -6 test/conformance/instance/instance.go; echo ===; head -6 test/conformance/upstream/record/main.go; echo ===; ls test/conformance/upstream/
rc=1
// Package instance composes a real in-process Astrate stack — store + engine +
// broker + pairing on a live TimescaleDB — for the M9 conformance runners
// (docs/ROADMAP.md §10). It is the shared harness behind the gosdk, atomvm,
// and pysdk packages: each provisions a realm, installs interfaces, and drives
// the running instance with a different client. The database is the compose
// service (`make up`) or ASTRATE_TEST_DSN; New skips the test when neither is
===
// Command record captures upstream Astarte's REST error envelopes and writes
// them to test/conformance/upstream as a fixture plus a verbatim transcript
// (M12 plan, phase 04).
//
// It exists because Astrate's error envelopes were reconstructed rather than
// observed: phase 01 shipped a 404 envelope whose *shape* and status were
===
README.md
channels.json
channels.transcript.txt
record
recordchannels
rest-errors.json
rest-errors.transcript.txt
verify-server-writes.json
verify-server-writes.transcript.txt
verify-versions.json
verify-versions.transcript.txt

$ rg -n 'go:build|\+build' test/conformance/instance/instance.go | head; echo "rc"; go list -e -f '{{.ImportPath}} err={{.Error}}' ./... 2>/dev/null | rg -v '^github' | head; echo "=== all list:"; go list ./... | wc -l; go list ./... | rg 'test/'; echo "=== go.mod locations:"; find . -name go.mod -not -path './.git/*'
rc
=== all list:
33
=== go.mod locations:
./test/conformance/go.mod
./tools/bsoncapture/go.mod
./bench/go.mod
./go.mod

$ cat .github/workflows/ci.yml
# CI pipeline (docs/ROADMAP.md §0.2, §10). Tiers: lint + unit (T1) + integration
# (T2) on every push/PR; e2e (T3) + official-SDK conformance (T4) on PRs to main;
# nightly load/footprint + security (T5) on a schedule.
name: ci

on:
  push:
    branches: [main]
  pull_request:
  schedule:
    - cron: "0 3 * * *" # nightly T5
  workflow_dispatch: # on-demand full sweep incl. the nightly T5 tier

permissions:
  contents: read

env:
  GOLANGCI_LINT_VERSION: v2.12.2 # keep in lockstep with Makefile GOLANGCI_LINT_VERSION
  COMPOSE_DSN: postgres://astrate:astrate@localhost:5432/astrate?sslmode=disable

jobs:
  lint:
    runs-on: ubuntu-latest
    timeout-minutes: 10
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - uses: golangci/golangci-lint-action@v8
        with:
          version: ${{ env.GOLANGCI_LINT_VERSION }}

  unit:
    runs-on: ubuntu-latest
    timeout-minutes: 15
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - run: make test

  integration:
    runs-on: ubuntu-latest # Docker daemon available for testcontainers
    timeout-minutes: 30
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - run: make test-integration

  # T3 component/E2E suites (broker/engine/appengine/cmd e2e + test/security).
  e2e:
    runs-on: ubuntu-latest
    timeout-minutes: 40
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - run: make test-e2e

  # T4 official-SDK conformance (CP-A/B/C/D): the Go + Python SDKs, the AtomVM
  # JSON simulator, astartectl, and the load smoke against a real database.
  conformance:
    runs-on: ubuntu-latest
    timeout-minutes: 45
    services:
      timescaledb:
        # -ha, not the plain image: it is the only one carrying
        # timescaledb_toolkit, and CI must exercise the lttb() downsampling
        # path on the same image developers run locally (docker-compose.yml).
        image: timescale/timescaledb-ha:pg16
        env:
          POSTGRES_DB: astrate
          POSTGRES_USER: astrate
          POSTGRES_PASSWORD: astrate
        ports:
          - 5432:5432
        options: >-
          --health-cmd "pg_isready -U astrate -d astrate"
          --health-interval 5s --health-timeout 3s --health-retries 12
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - uses: actions/setup-python@v5
        with:
          python-version: "3.11"
      - name: Install the pinned Python SDK
        run: pip install -r test/conformance/pysdk/requirements.txt
      - name: Conformance (CP-A..CP-D)
        env:
          ASTRATE_TEST_DSN: ${{ env.COMPOSE_DSN }}
          ASTRATE_PYSDK_PYTHON: python3
          CGO_ENABLED: "1" # the official Go SDK pulls bson/sqlite
        run: make test-conformance

  # T5 nightly: footprint/load budgets and the security hardening sweep.
  nightly:
    if: github.event_name == 'schedule' || github.event_name == 'workflow_dispatch'
    runs-on: ubuntu-latest
    timeout-minutes: 90
    services:
      timescaledb:
        # -ha, not the plain image: it is the only one carrying
        # timescaledb_toolkit, and CI must exercise the lttb() downsampling
        # path on the same image developers run locally (docker-compose.yml).
        image: timescale/timescaledb-ha:pg16
        env:
          POSTGRES_DB: astrate
          POSTGRES_USER: astrate
          POSTGRES_PASSWORD: astrate
        ports:
          - 5432:5432
        options: >-
          --health-cmd "pg_isready -U astrate -d astrate"
          --health-interval 5s --health-timeout 3s --health-retries 12
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - name: Security hardening suite
        run: make test-e2e
      - name: Load + footprint budgets (§5.4)
        env:
          ASTRATE_TEST_DSN: ${{ env.COMPOSE_DSN }}
          LOAD_DEVICES: "200"
          LOAD_MESSAGES: "100"
          LOAD_P99_MS: "250"
        run: cd test/conformance && go test -run TestLoadSmoke -timeout 60m ./load/

$ export PATH=$PATH:$HOME/go/bin; cd test/conformance && go list ./... && echo "=== notest:" && for d in $(go list ./...); do dir=${d#*astrate/test/conformance/}; ls $dir/*_test.go >/dev/null 2>&1 || echo "NOTEST $dir"; done; echo "=== vet:"; go vet ./... 2>&1 | tail -10; echo "vet=$?"
github.com/astrate-platform/astrate/test/conformance/atomvm
github.com/astrate-platform/astrate/test/conformance/cpa
github.com/astrate-platform/astrate/test/conformance/cpb
github.com/astrate-platform/astrate/test/conformance/cpc
github.com/astrate-platform/astrate/test/conformance/gosdk
github.com/astrate-platform/astrate/test/conformance/instance
github.com/astrate-platform/astrate/test/conformance/load
github.com/astrate-platform/astrate/test/conformance/pysdk
github.com/astrate-platform/astrate/test/conformance/upstream/record
github.com/astrate-platform/astrate/test/conformance/upstream/recordchannels
=== notest:
NOTEST instance
NOTEST upstream/record
NOTEST upstream/recordchannels
=== vet:
go: downloading github.com/astarte-platform/astarte-go v0.90.4
go: downloading github.com/cristalhq/jwt/v3 v3.1.0
go: downloading github.com/iancoleman/orderedmap v0.2.0
go: downloading github.com/ispirata/paho.mqtt.golang v1.3.91-0.20220121155423-6e36d43b2ec9
go: downloading gorm.io/gorm v1.22.3
go: downloading gorm.io/driver/sqlite v1.2.6
go: downloading github.com/mattn/go-sqlite3 v1.14.22
go: downloading github.com/go-stack/stack v1.8.0
go: downloading github.com/jinzhu/inflection v1.0.0
go: downloading github.com/jinzhu/now v1.1.3
vet=0

$ timeout 15 curl -sI https://proxy.golang.org 2>&1 | head -3; echo "---"; timeout 15 curl -sI https://vuln.go.dev 2>&1 | head -3; echo "rc=$?"; go env GOPROXY GOFLAGS GONOSUMDB 2>/dev/null
HTTP/2 200 
accept-ranges: bytes
access-control-allow-headers: Range
---
HTTP/2 200 
x-guploader-uploadid: AP6rU80mIWPugYLKVp--7nQBJNMBoLjJhC82i-m1-3iJvZagunN8bC8JY0zvbGHThBVzcph0sDpdz2U
x-goog-generation: 1790281578638935
rc=0
https://proxy.golang.org,direct

