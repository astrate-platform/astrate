slug: docs-sync-ae-data-write-object-errors
verdict: done
at:  e974103
ran: 2026-10-09T17:43:28Z on DietPi in 587s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/appengine/http.go [offset=660, limit=110]
I'll read the relevant files first.
→ Read docs/api/astarte_appengine_api.yaml [offset=480, limit=120]
→ Read docs/api/astarte_appengine_api.yaml [limit=120, offset=740]
→ Read docs/api/astarte_appengine_api.yaml [limit=120, offset=1300]
→ Read docs/api/astarte_appengine_api.yaml [limit=230, offset=1600]
→ Read docs/api/astarte_appengine_api.yaml [limit=100, offset=1829]

$ rg -n "ValueTooLarge|BadRequest|ErrorDetail|Unexpected object key|Missing required|unexpected_keys" docs/api/astarte_appengine_api.yaml
245:          $ref: "#/components/responses/BadRequest"
333:          $ref: "#/components/responses/BadRequest"
518:          $ref: "#/components/responses/BadRequest"
528:          $ref: "#/components/responses/ValueTooLarge"
570:          $ref: "#/components/responses/BadRequest"
580:          $ref: "#/components/responses/ValueTooLarge"
613:          $ref: "#/components/responses/BadRequest"
775:          $ref: "#/components/responses/BadRequest"
785:          $ref: "#/components/responses/ValueTooLarge"
822:          $ref: "#/components/responses/BadRequest"
832:          $ref: "#/components/responses/ValueTooLarge"
860:          $ref: "#/components/responses/BadRequest"
940:          $ref: "#/components/responses/BadRequest"
1064:          $ref: "#/components/responses/BadRequest"
1146:          $ref: "#/components/responses/BadRequest"
1336:          $ref: "#/components/responses/BadRequest"
1346:          $ref: "#/components/responses/ValueTooLarge"
1384:          $ref: "#/components/responses/BadRequest"
1394:          $ref: "#/components/responses/ValueTooLarge"
1423:          $ref: "#/components/responses/BadRequest"
1770:    ErrorDetail:
1794:    BadRequest:
1799:            $ref: "#/components/schemas/ErrorDetail"
1809:            $ref: "#/components/schemas/ErrorDetail"
1819:            $ref: "#/components/schemas/ErrorDetail"
1829:            $ref: "#/components/schemas/ErrorDetail"
1839:            $ref: "#/components/schemas/ErrorDetail"
1849:            $ref: "#/components/schemas/ErrorDetail"
1859:            $ref: "#/components/schemas/ErrorDetail"
1875:    ValueTooLarge:
1880:            $ref: "#/components/schemas/ErrorDetail"
1890:            $ref: "#/components/schemas/ErrorDetail"
1900:            $ref: "#/components/schemas/ErrorDetail"
1910:            $ref: "#/components/schemas/ErrorDetail"
1920:            $ref: "#/components/schemas/ErrorDetail"

$ rg -n "ValueTooLarge" docs/ ; echo "---MAKEFILE---"; ls docs; echo "---"; cat docs/Makefile 2>/dev/null
docs/api/astarte_appengine_api.yaml:528:          $ref: "#/components/responses/ValueTooLarge"
docs/api/astarte_appengine_api.yaml:580:          $ref: "#/components/responses/ValueTooLarge"
docs/api/astarte_appengine_api.yaml:785:          $ref: "#/components/responses/ValueTooLarge"
docs/api/astarte_appengine_api.yaml:832:          $ref: "#/components/responses/ValueTooLarge"
docs/api/astarte_appengine_api.yaml:1346:          $ref: "#/components/responses/ValueTooLarge"
docs/api/astarte_appengine_api.yaml:1394:          $ref: "#/components/responses/ValueTooLarge"
docs/api/astarte_appengine_api.yaml:1875:    ValueTooLarge:
docs/site/api/astarte_appengine_api.yaml:528:          $ref: "#/components/responses/ValueTooLarge"
docs/site/api/astarte_appengine_api.yaml:580:          $ref: "#/components/responses/ValueTooLarge"
docs/site/api/astarte_appengine_api.yaml:785:          $ref: "#/components/responses/ValueTooLarge"
docs/site/api/astarte_appengine_api.yaml:832:          $ref: "#/components/responses/ValueTooLarge"
docs/site/api/astarte_appengine_api.yaml:1346:          $ref: "#/components/responses/ValueTooLarge"
docs/site/api/astarte_appengine_api.yaml:1394:          $ref: "#/components/responses/ValueTooLarge"
docs/site/api/astarte_appengine_api.yaml:1875:    ValueTooLarge:
---MAKEFILE---
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
UPSTREAM-EXPERIMENTAL.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui
---
.PHONY: serve build clean install sync

install:
	pip install -r requirements.txt

sync:
	cp DESIGN.md site/DESIGN.md
	cp ROADMAP.md site/ROADMAP.md
	cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
	rm -rf site/api && cp -r api site/api
	rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui

serve: sync
	mkdocs serve -f mkdocs.yml

build: sync
	mkdocs build -f mkdocs.yml

clean:
	rm -rf ../site-dist

$ cat docs/mkdocs.yml; echo "===SWAGGER-UI==="; ls -R docs/swagger-ui | head -50; echo "===EMBED==="; cat docs/embed.go
site_name: Astrate
site_description: Lean, single-binary, Astarte-wire-compatible IoT platform in Go
repo_url: https://github.com/atsetilam/astrate

docs_dir: site
site_dir: ../site-dist

theme:
  name: material
  palette:
    - scheme: default
      primary: teal
      accent: amber
      toggle:
        icon: material/brightness-7
        name: Switch to dark mode
    - scheme: slate
      primary: teal
      accent: amber
      toggle:
        icon: material/brightness-4
        name: Switch to light mode
  features:
    - navigation.instant
    - navigation.sections
    - navigation.top
    - content.code.copy
    - toc.follow

markdown_extensions:
  - admonition
  - tables
  - toc:
      permalink: true
  - pymdownx.superfences
  - pymdownx.tabbed:
      alternate_style: true
  - pymdownx.highlight:
      anchor_linenums: true
  - pymdownx.inlinehilite

nav:
  - Home: index.md
  - Architecture: architecture.md
  - Data Modeling: data-modeling.md
  - MQTT Protocol: mqtt-protocol.md
  - Pairing & Security: pairing-and-security.md
  - Payload Formats: payload-formats.md
  - JSON Payload Profile: json-payload-profile.md
  - AppEngine API: appengine-api.md
  - Realm Management API: realm-management-api.md
  - Housekeeping API: housekeeping-api.md
  - Interface Schema: interface-schema.md
  - Triggers: triggers.md
  - Deployment: deployment.md
  - Operations: operations.md
  - Observability: observability.md
  - Compatibility: compatibility.md
  - Configuration Reference: configuration-reference.md
  - Troubleshooting: troubleshooting.md
  - Contributing: contributing.md
  - Migration from Astarte: migration-from-astarte.md
  - API Explorer: swagger.md
  - Quickstart: quickstart.md
  - Examples:
    - Pokémon Agent: pokemon-agent.md
===SWAGGER-UI===
docs/swagger-ui:
index.html
===EMBED===
// Package docs embeds the Swagger UI static files and OpenAPI YAML specs so
// they are served directly from the astrate binary without external files.
package docs

import "embed"

// SwaggerUI holds the static Swagger UI files (index.html, CSS, JS refs).
//
//go:embed swagger-ui/*
var SwaggerUI embed.FS

// APIYAML holds the OpenAPI 3.0 YAML specifications for all five API surfaces.
//
//go:embed api/*.yaml
var APIYAML embed.FS

$ rg -n "astarte_appengine_api|swagger|openapi|OpenAPI|yaml" docs/swagger-ui docs/embed.go docs/mkdocs.yml 2>/dev/null | head -40
docs/embed.go:1:// Package docs embeds the Swagger UI static files and OpenAPI YAML specs so
docs/embed.go:9://go:embed swagger-ui/*
docs/embed.go:12:// APIYAML holds the OpenAPI 3.0 YAML specifications for all five API surfaces.
docs/embed.go:14://go:embed api/*.yaml
docs/mkdocs.yml:63:  - API Explorer: swagger.md
docs/swagger-ui/index.html:7:  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5.32.11/favicon-32x32.png">
docs/swagger-ui/index.html:8:  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui.css">
docs/swagger-ui/index.html:32:    #swagger-ui { max-width: 100%; }
docs/swagger-ui/index.html:39:    <span class="subtitle">OpenAPI 3.0 Reference</span>
docs/swagger-ui/index.html:44:      <option value="../api/astarte_housekeeping_api.yaml">Housekeeping API — Realm management, a_ha auth</option>
docs/swagger-ui/index.html:45:      <option value="../api/astarte_pairing_api.yaml">Pairing API — Device pairing flows A/B/C, a_pa auth</option>
docs/swagger-ui/index.html:46:      <option value="../api/astarte_realm_management_api.yaml">Realm Management API — Interfaces, triggers, policies, a_rma auth</option>
docs/swagger-ui/index.html:47:      <option value="../api/astarte_appengine_api.yaml">AppEngine API — Devices, data, groups, a_aea auth</option>
docs/swagger-ui/index.html:48:      <option value="../api/astrate_native_api.yaml">Astrate Native API — Health, readiness, metrics, WebSocket</option>
docs/swagger-ui/index.html:53:  <div id="swagger-ui"></div>
docs/swagger-ui/index.html:55:  <script src="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui-bundle.js"></script>
docs/swagger-ui/index.html:56:  <script src="https://unpkg.com/swagger-ui-dist@5.32.11/swagger-ui-standalone-preset.js"></script>
docs/swagger-ui/index.html:60:    const uiContainer = document.getElementById('swagger-ui');
docs/swagger-ui/index.html:66:        dom_id: '#swagger-ui',

$ rg -n "payload-reason|reason-status|ReasonMissingRequired|ReasonBadObject|unexpected_keys|Missing required mapping|Unexpected object key" .mule/ docs/ internal/ --glob '!docs/site/**' 2>/dev/null | head -50
.mule/reviews/engine-triggers-2026-09-20.md:24:- `pkg/payload/payload.go:148` (ReasonMissingRequired emission)
.mule/todo.md:265:- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
.mule/todo.md:266:- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress. — BLOCKED: TIMEOUT after 1200s — task too big, split it
.mule/todo.md:267:- [x] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
.mule/todo.md:285:- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
.mule/todo.md:300:- [ ] docs-sync-ae-data-write-object-errors [auto]: document the two object-aggregated write rejections on the six data PUT/POST ops (putData/publishData, by-alias, in-group) in docs/api/astarte_appengine_api.yaml — each refs `BadRequest` for 400 (yaml:518, 570, 775, 822, 1336, 1384) and `ValueTooLarge` for 422 (yaml:528, 580, 785, 832, 1346, 1394), but three bodies are missing. (a) 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` from `writeBadObjectError` (internal/appengine/http.go:722-748, const `detailUnexpectedObjectKey` :711), reached by the `payload.ReasonOf(err) == payload.ReasonBadObject` case (http.go:684-691) — so the shared `ErrorDetail` schema (yaml:1741-1750) needs an optional `unexpected_keys` array, or a dedicated response. (b) 422 `{"errors":{"detail":"Missing required mapping key"}}` (const `detailMissingRequiredMapping`, http.go:713-720) from the `ReasonMissingRequired` case (http.go:692-702) — add an example beside the value-size one already on those ops. Both were implemented by appengine-unexpected-object-key/appengine-missing-required-422 (upstream master b6d46ad4, #2237) with no spec update. Do NOT invent other `payload.ReasonOf` statuses — that family is the separate blocked line appengine-payload-reason-status-map. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
internal/engine/data_test.go:186:			reason: payload.ReasonBadObject.String(),
internal/engine/serverdata_test.go:338:	if payload.ReasonOf(err) != payload.ReasonBadObject {
internal/appengine/http.go:684:	case payload.ReasonOf(err) == payload.ReasonBadObject:
internal/appengine/http.go:692:	case payload.ReasonOf(err) == payload.ReasonMissingRequired:
internal/appengine/http.go:711:const detailUnexpectedObjectKey = "Unexpected object key"
internal/appengine/http.go:716:// renders `%{errors: %{detail: "Missing required mapping key"}}`). Frozen
internal/appengine/http.go:720:const detailMissingRequiredMapping = "Missing required mapping key"
internal/appengine/http.go:723:// upstream envelope: {"errors": {"detail": "Unexpected object key",
internal/appengine/http.go:724:// "unexpected_keys": [...]}}. The other object-shape failures — not a
internal/appengine/http.go:736:		Keys   []string `json:"unexpected_keys"`
internal/appengine/writeerror_test.go:35:			&payload.RejectError{Reason: payload.ReasonBadObject,
internal/appengine/writeerror_test.go:38:			http.StatusBadRequest, "Unexpected object key", []string{"ghost"}},
internal/appengine/writeerror_test.go:40:			fmt.Errorf("engine: %w", &payload.RejectError{Reason: payload.ReasonBadObject,
internal/appengine/writeerror_test.go:43:			http.StatusBadRequest, "Unexpected object key", []string{"a", "b"}},
internal/appengine/writeerror_test.go:45:			&payload.RejectError{Reason: payload.ReasonBadObject, Detail: "object-aggregation document is empty"},
internal/appengine/writeerror_test.go:48:			&payload.RejectError{Reason: payload.ReasonMissingRequired,
internal/appengine/writeerror_test.go:50:			http.StatusUnprocessableEntity, "Missing required mapping key", nil},
internal/appengine/writeerror_test.go:64:					Keys   []string `json:"unexpected_keys"`
internal/appengine/writeerror_test.go:74:				t.Fatalf("unexpected_keys = %v, want %v", body.Errors.Keys, tc.keys)
internal/appengine/writeerror_test.go:78:					t.Errorf("unexpected_keys[%d] = %q, want %q", i, body.Errors.Keys[i], k)
internal/appengine/writeerror_test.go:87:// {"errors":{"detail":"Unexpected object key","unexpected_keys":["ghost"]}}.
internal/appengine/writeerror_test.go:92:		Reason:         payload.ReasonBadObject,
internal/appengine/writeerror_test.go:96:	want := `{"errors":{"detail":"Unexpected object key","unexpected_keys":["ghost"]}}`
.mule/log.md:168:| 2026-10-03 | appengine-payload-reason-status-map | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
.mule/for-giulio.md:96:- **The Legion Go answers ssh but its Docker is not running, so every "probe live upstream" line is blocked on you, not on the hardware (measured 2026-10-03).** `ssh -o ConnectTimeout=5 legion 'echo ok'` returns `ok`, but `docker ps` fails with `dial unix /home/atsetilam/.docker/desktop/docker.sock: no such file or directory`, `systemctl is-active docker` answers `inactive`, and `sudo -n true` answers `è necessaria una password` — so the mule cannot start it. It is Docker Desktop's socket, not a rootful `dockerd` the mule could launch as a user, and no rootless extras are installed. Consequence for the queue: `appengine-payload-reason-status-map` (the remaining eight `payload.Reason` values that still answer 500) cannot be measured the way its line demands. Today's `appengine-missing-required-422` was therefore mapped from the frozen upstream source at `b6d46ad4` (`Device.validate_required_mappings/2` → `fallback_controller.ex:215-219` → `ErrorView "422_missing_required_mapping.json"`), and the commit says so plainly. Nothing about that sibling line needs rewriting — only Docker up on the Legion Go.
.mule/reports/log/20261009T105827Z-recipe-code-review-proposed.md:95:| 2026-10-03 | appengine-payload-reason-status-map | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
.mule/reports/log/20261009T105827Z-recipe-code-review-proposed.md:335:- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
.mule/reports/log/20261009T105827Z-recipe-code-review-proposed.md:336:- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress. — BLOCKED: TIMEOUT after 1200s — task too big, split it
.mule/reports/log/20261009T105827Z-recipe-code-review-proposed.md:337:- [x] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
.mule/reports/log/20261009T105827Z-recipe-code-review-proposed.md:355:- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
.mule/reports/log/20261009T112852Z-recipe-docs-sync-proposed.md:41:   `{"errors":{"detail":"Unexpected object key","unexpected_keys":[...]}}`
.mule/reports/log/20261009T112852Z-recipe-docs-sync-proposed.md:42:   (`writeBadObjectError`, http.go:722-748, const :711, ReasonBadObject case
.mule/reports/log/20261009T112852Z-recipe-docs-sync-proposed.md:43:   :684-691) and 422 `{"errors":{"detail":"Missing required mapping key"}}`
.mule/reports/log/20261009T112852Z-recipe-docs-sync-proposed.md:44:   (`detailMissingRequiredMapping` http.go:713-720, ReasonMissingRequired case
.mule/reports/log/20261009T112852Z-recipe-docs-sync-proposed.md:47:   (yaml:1741-1750) also needs the optional `unexpected_keys` array.
.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:61:.mule/for-giulio.md:82:- **The Legion Go answers ssh but its Docker is not running, so every "probe live upstream" line is blocked on you, not on the hardware (measured 2026-10-03).** `ssh -o ConnectTimeout=5 legion 'echo ok'` returns `ok`, but `docker ps` fails with `dial unix /home/atsetilam/.docker/desktop/docker.sock: no such file or directory`, `systemctl is-active docker` answers `inactive`, and `sudo -n true` answers `è necessaria una password` — so the mule cannot start it. It is Docker Desktop's socket, not a rootful `dockerd` the mule could launch as a user, and no rootless extras are installed. Consequence for the queue: `appengine-payload-reason-status-map` (the remaining eight `payload.Reason` values that still answer 500) cannot be measured the way its line demands. Today's `appengine-missing-required-422` was therefore mapped from the frozen upstream source at `b6d46ad4` (`Device.validate_required_mappings/2` → `fallback_controller.ex:215-219` → `ErrorView "422_missing_required_mapping.json"`), and the commit says so plainly. Nothing about that sibling line needs rewriting — only Docker up on the Legion Go.
.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:104:282:- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:124:- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:125:- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress. — BLOCKED: TIMEOUT after 1200s — task too big, split it
.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:126:- [x] appengine-missing-required-422 [auto]: `ReasonMissingRequired` (raised at pkg/payload/payload.go:148 for an object-aggregated document that omits a mapping's required key, verified for both wire formats in pkg/payload/required_test.go) answers 500 where upstream renders a distinct 422 `missing_required_mapping` — the exact status and detail string are unmeasured, so do not invent them: probe the write against upstream master b6d46ad4 on the Legion Go, record the status and the `{"errors":{...}}` body, then map the reason in `writeError` (internal/appengine/http.go) beside the `bad_object` case and add the case to the internal/appengine/writeerror_test.go table. Separate from the general line above because it is the one half appengine-unexpected-object-key was told to report rather than fix.
.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:141:- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
.mule/reports/log/20261007T171801Z-recipe-astarte-upstream-proposed.md:278:282:- [!] probe-encrypted-mapping-enforcement [auto]: upstream master 1da04832 "chore: forward port release-1.4" (2026-10-06, 1232 files, **unreleased** — absent from both v1.4.0-rc.6 and v1.3.5, verified by fetching those tags: rc.6's DUP data_handler and its libs/astarte_secrets contain no "encrypt" at all) turned the `encrypted` mapping flag from inert into enforced: new `apps/astarte_secrets/lib/astarte_secrets/encrypted_messages.ex`, DEK-encrypted values at rest (`encrypted_dek` column), transport-encrypted server→device publishes and decrypt-on-read in both DUP (`data_updater/core/{data_handler,control_handler}.ex`) and AppEngine (`device/device.ex` `maybe_apply_transport_encryption`, `device/queries.ex` `maybe_decrypt_values`), a dashboard MappingEditor toggle (+7), and +15 lines in `082-key_agreement_protocol.md`. Astrate accepts and persists `encrypted: true` on datastream mappings (internal/interfaceschema/parse.go:410-414) but has no encryption code anywhere (no SharedSecret/encrypted_dek under internal/ or pkg/). Investigate only: list exactly where Astrate's behaviour would diverge from this once it ships — interface acceptance, AppEngine server-write payload, read-back shape, MQTT server-owned publish — and whether an unmodified device SDK observes it. Report, do not patch; docs/UPSTREAM-EXPERIMENTAL.md row #67 and issues #92/#93 stay parked. Same commit also brings the `unexpected_keys: [...]` field on the 400 `Unexpected object key` body — already implemented (appengine-unexpected-object-key, internal/appengine/http.go:711) — and the libs→apps monorepo restructure with vendored `apps/astarte_core` interface/mapping/policy JSON specs, which recapitulate the existing field set (no new schema key). No `gh issue create`, nothing commented, closed or edited on GitHub. — BLOCKED: wrote nothing
.mule/reports/log/20261008T111054Z-recipe-github-issues-proposed.md:97:- [x] appengine-unexpected-object-key [auto]: upstream master b6d46ad4 (#2237, 2026-10-02, **unreleased**) answers an object-aggregated write carrying an undeclared key with 400 `{"errors":{"detail":"Unexpected object key","unexpected_keys":["<sorted keys>"]}}` — `InterfaceValue.unexpected_keys/2` (interface_value.ex:151-162) computes the offending keys sorted, the error becomes `{:error, :unexpected_object_key, keys: keys}` in both `Device` and `InterfaceValue`, and ErrorView "422_unexpected_object_key" now renders the list (fallback_controller.ex:201-206 still puts `:bad_request`; before this commit the detail alone was rendered). Astrate already rejects the same write — `rejectf(ReasonBadObject, "key %q matches no declared object leaf")` (pkg/payload/bson.go:157, pkg/payload/json.go:216) raised by the `dec.Object` call in `publishAsOwner` (internal/engine/serverdata.go) — but `writeError` (internal/appengine/http.go:643) maps no `payload.Reason` whatsoever, so a bad object shape falls through to `default:` and answers **500** "Internal server error" on all six data PUT/POST ops (device, by-alias, group — putData/putDataByAlias/putDataInGroup, http.go:292-357) where upstream answers 400, naming no key. Map the object-shape rejection to 400 with `detail: "Unexpected object key"` and the sorted unexpected keys beside it in the same `{"errors":{...}}` envelope, and add a T1 test: an object write carrying one declared and one undeclared key answers 400 with the undeclared key listed. While in the mapping, check whether `ReasonMissingRequired` has the same hole (upstream renders a distinct 422 `missing_required_mapping`) — report that half, do not guess the status.
.mule/reports/log/20261008T111054Z-recipe-github-issues-proposed.md:98:- [!] appengine-payload-reason-status-map [auto]: the whole `payload.RejectReason` family except `bad_object` still answers **500** on the six AppEngine data write ops — measured 2026-10-03 with a throwaway probe over `writeError` (internal/appengine/http.go): `missing_required`, `too_large`, `malformed`, `no_value`, `type_mismatch`, `value_too_large`, `unset_not_allowed`, `bad_timestamp` and `unknown_format` all reach `default:` → `WriteInternalServerError`, i.e. a client that sends a mistyped double is told the server broke. `bad_object` is now mapped (appengine-unexpected-object-key); the other nine need one upstream status+detail each, and none of them may be guessed — probe each against upstream master b6d46ad4 on the Legion Go, then extend the `payload.ReasonOf(err)` case in `writeError` into a full switch with a table-driven case per reason in internal/appengine/writeerror_test.go. `too_large` is the one with a second caller to check: `Decoder.MaxSize` is fed from `len(envelope)` in `publishAsOwner` (internal/engine/serverdata.go:134), so the REST cap may be unreachable and the reason may only fire from the MQTT ingress. — BLOCKED: TIMEOUT after 1200s — task too big, split it

→ Read .mule/reports/log/20261009T112852Z-recipe-docs-sync-proposed.md

$ rg -n "astarte_appengine_api|APIYAML|docs/api|swagger-ui" --glob '*_test.go' --glob '*.go' -l | head -30
docs/embed.go
internal/swagger/swagger.go
internal/swagger/swagger_test.go
internal/appengine/service.go
internal/appengine/http.go

$ rg -n "docs-sync|docs/api" .mule/recipes/ 2>/dev/null | head; echo "=== recipes ==="; ls .mule/recipes/
.mule/recipes/docs-sync.md:3:The docs site (`docs/site/`, MkDocs, plus OpenAPI specs in `docs/api/` and Swagger UI) drifts
.mule/recipes/docs-sync.md:7:The exception, and it is the useful half of this recipe: `docs/api/*.yaml` are *generated
.mule/recipes/docs-sync.md:18:rg -n '^\s{2}/' docs/api/astrate_<surface>_api.yaml   # the documented paths
.mule/recipes/docs-sync.md:26:- **Documented but absent**, or **present but undocumented**: a `docs/api/` fix task, naming
.mule/recipes/docs-sync.md:38:If any task in this family touches `docs/api/` or `docs/mkdocs.yml`, the executing task must
=== recipes ===
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

→ Read .mule/recipes/docs-sync.md
→ Read internal/swagger/swagger_test.go

$ rg -n "^ *examples:|^ *example:" docs/api/astarte_appengine_api.yaml; echo "=== all specs examples: ==="; rg -n "examples:" docs/api/*.yaml
56:              example:
104:              examples:
170:              example:
230:            example:
318:            example:
891:              example:
923:            example:
970:              example:
1057:            example:
1131:            example:
1800:          example:
1810:          example:
1820:          example:
1830:          example:
1840:          example:
1850:          example:
1860:          example:
1870:          example:
1881:          example:
1891:          example:
1901:          example:
1911:          example:
1921:          example:
=== all specs examples: ===
docs/api/astarte_appengine_api.yaml:104:              examples:
docs/api/astarte_housekeeping_api.yaml:490:          examples:
docs/api/astarte_pairing_api.yaml:308:              examples:
docs/api/astarte_pairing_api.yaml:650:          examples:
docs/api/astarte_realm_management_api.yaml:61:              examples:
docs/api/astarte_realm_management_api.yaml:147:              examples:
docs/api/astarte_realm_management_api.yaml:283:              examples:
docs/api/astarte_realm_management_api.yaml:336:              examples:
docs/api/astarte_realm_management_api.yaml:429:              examples:
docs/api/astarte_realm_management_api.yaml:803:              examples:
docs/api/astarte_realm_management_api.yaml:1365:          examples:

$ rg -n "unexpected_keys|object key|Missing required|ValueTooLarge|BadObject|MissingRequired" docs/api/*.yaml
docs/api/astarte_appengine_api.yaml:528:          $ref: "#/components/responses/ValueTooLarge"
docs/api/astarte_appengine_api.yaml:580:          $ref: "#/components/responses/ValueTooLarge"
docs/api/astarte_appengine_api.yaml:785:          $ref: "#/components/responses/ValueTooLarge"
docs/api/astarte_appengine_api.yaml:832:          $ref: "#/components/responses/ValueTooLarge"
docs/api/astarte_appengine_api.yaml:1346:          $ref: "#/components/responses/ValueTooLarge"
docs/api/astarte_appengine_api.yaml:1394:          $ref: "#/components/responses/ValueTooLarge"
docs/api/astarte_appengine_api.yaml:1722:            (see `metadata`). `disjoint_tables` is an object keyed by path:
docs/api/astarte_appengine_api.yaml:1875:    ValueTooLarge:

→ Read docs/api/astarte_appengine_api.yaml [limit=40, offset=90]
→ Read docs/api/astarte_housekeeping_api.yaml [limit=40, offset=475]

$ rg -n "BadRequest|ValueTooLarge|BadObject|Unexpected object key|Missing required|unexpected_keys|data-write|DataWrite" internal/swagger/swagger_test.go
380:// $ref the ValidationErrors response — pointing it at BadRequest would teach a
437:		ref = `          $ref: "#/components/responses/BadRequest"`
454:			t.Errorf("operation %s 400 is followed by %q, want the BadRequest $ref", op, block[idx+1])
458:	if !containsLine(lines, "    BadRequest:") {
459:		t.Error("components.responses defines no BadRequest response")
486:		{"    BadRequest:", astarteapi.DetailBadRequest},
515:// from TestPairingErrorDetailExamples. The BadRequest and InternalServerError
518:// astarteapi.WriteBadRequest and every 500 through WriteInternalServerError
537:		{"    BadRequest:", astarteapi.DetailBadRequest},
568:// TestRealmManagementErrorDetailExamples. The BadRequest and InternalServerError
571:// astarteapi.WriteBadRequest and every 500 through WriteInternalServerError
574:// with no constant behind it (MethodNotAllowed, ValueTooLarge,
589:		{"    BadRequest:", astarteapi.DetailBadRequest},
628:// BadRequest component at all: the native surface is the compat health and
630:// WriteBadRequest, so there is nothing here to pin.
915:// the last spec still carrying a reconstructed capitalisation: its BadRequest
918:// emits — every 400 on this surface goes through astarteapi.WriteBadRequest
931:// string as BadRequest. It sits on the generic schema all eight error
951:		{"    BadRequest:", astarteapi.DetailBadRequest},
986:	if got != astarteapi.DetailBadRequest {
988:			got, astarteapi.DetailBadRequest)

$ rg -rn "ValueTooLarge|BadObject|Unexpected object key|Missing required mapping" --glob '*_test.go' .
./internal/engine/data_test.go:			reason: payload.Reasonn.String(),
./internal/engine/data_test.go:			reason: payload.Reasonn.String(),
./internal/engine/serverdata_test.go:	if payload.ReasonOf(err) != payload.Reasonn {
./pkg/payload/payload_test.go:		`{"v":{"lat":1.0,"nope":2.0}}`: Reasonn,    // undeclared key
./pkg/payload/payload_test.go:		`{"v":{}}`:                     Reasonn,    // empty document
./pkg/payload/payload_test.go:		`{"v":5}`:                      Reasonn,    // not a document
./pkg/payload/payload_test.go:		`{"v":null}`:                   Reasonn,    // null
./pkg/payload/payload_test.go:		`{"v":[1]}`:                    Reasonn,    // array
./pkg/payload/payload_test.go:	if _, err := DecodeObject(dup, objectLeaves(false)); ReasonOf(err) != Reasonn {
./pkg/payload/payload_test.go:			if ReasonOf(err) != Reasonn {
./pkg/payload/payload_test.go:	if _, err := wide.Individual(overString, mapping(interfaceschema.String, false)); ReasonOf(err) != Reasonn {
./pkg/payload/payload_test.go:	if _, err := Decode(overJSON, mapping(interfaceschema.IntegerArray, false)); ReasonOf(err) != Reasonn {
./pkg/payload/payload_test.go:	if _, err := Decode(overBSON, mapping(interfaceschema.IntegerArray, false)); ReasonOf(err) != Reasonn {
./pkg/payload/payload_test.go:		if _, err := Encode(strings.Repeat("a", MaxStringLen+1), nil, f); ReasonOf(err) != Reasonn {
./pkg/payload/payload_test.go:		if _, err := Encode(make([]bool, MaxArrayLen+1), nil, f); ReasonOf(err) != Reasonn {
./pkg/payload/payload_test.go:		if _, err := Encode(map[string]Value{}, nil, f); ReasonOf(err) != Reasonn {
./pkg/payload/payload_test.go:		if _, err := Encode(map[string]Value{"a": map[string]Value{"b": 1.0}}, nil, f); ReasonOf(err) != Reasonn {
./pkg/payload/value_test.go:	if _, err := checkString(atCap + "x"); ReasonOf(err) != Reasonn {
./pkg/payload/value_test.go:	if _, err := collectArray[int32](MaxArrayLen+1, nil); ReasonOf(err) != Reasonn {
./pkg/payload/value_test.go:		Reasonn:   "value_too_large",
./pkg/payload/value_test.go:		Reasonn:       "bad_object",
./internal/swagger/swagger_test.go:// with no constant behind it (MethodNotAllowed, n,
./internal/appengine/writeerror_test.go:			&payload.RejectError{Reason: payload.Reasonn,
./internal/appengine/writeerror_test.go:			http.StatusBadRequest, "n", []string{"ghost"}},
./internal/appengine/writeerror_test.go:			fmt.Errorf("engine: %w", &payload.RejectError{Reason: payload.Reasonn,
./internal/appengine/writeerror_test.go:			http.StatusBadRequest, "n", []string{"a", "b"}},
./internal/appengine/writeerror_test.go:			&payload.RejectError{Reason: payload.Reasonn, Detail: "object-aggregation document is empty"},
./internal/appengine/writeerror_test.go:			http.StatusUnprocessableEntity, "n key", nil},
./internal/appengine/writeerror_test.go:// TestWriteErrornEnvelope freezes the exact bytes upstream emits for
./internal/appengine/writeerror_test.go:// {"errors":{"detail":"n","unexpected_keys":["ghost"]}}.
./internal/appengine/writeerror_test.go:func TestWriteErrornEnvelope(t *testing.T) {
./internal/appengine/writeerror_test.go:		Reason:         payload.Reasonn,
./internal/appengine/writeerror_test.go:	want := `{"errors":{"detail":"n","unexpected_keys":["ghost"]}}`

$ wc -l internal/swagger/swagger_test.go; sed -n '1304,1600p' internal/swagger/swagger_test.go 2>/dev/null | head -0; rg -n "func Test|func componentBlock|func operationBlock|func operationDescription|func containsLine" internal/swagger/swagger_test.go
1579 internal/swagger/swagger_test.go
26:func TestMount(t *testing.T) {
85:func TestMountSubPanicsOnBrokenFS(t *testing.T) {
98:func TestSpecs(t *testing.T) {
131:func TestRealmManagement403(t *testing.T) {
142:func TestAppEngine403(t *testing.T) {
206:func TestHousekeepingAsyncOperationParamDocumented(t *testing.T) {
249:func TestRealmManagementAsyncOperationParamDocumented(t *testing.T) {
324:func TestHousekeepingRetentionZeroFoldDocumented(t *testing.T) {
384:func TestAppEngineAddGroupDevice422Documented(t *testing.T) {
428:func TestAppEngineDataDelete400Documented(t *testing.T) {
473:func TestPairingErrorDetailExamples(t *testing.T) {
524:func TestRealmManagementErrorDetailExamples(t *testing.T) {
576:func TestAppEngineErrorDetailExamples(t *testing.T) {
631:func TestNativeErrorDetailExamples(t *testing.T) {
682:func TestNativeVersionDescriptions(t *testing.T) {
731:func TestNativeSocketSecurityDocumented(t *testing.T) {
802:func TestNativeMetricsContentNegotiationDocumented(t *testing.T) {
938:func TestHousekeepingErrorDetailExamples(t *testing.T) {
1009:func TestRealmManagement422ValidationDetails(t *testing.T) {
1086:func TestPairingDeviceIDEncodingDocumented(t *testing.T) {
1206:func TestRealmManagementDeviceIDEncodingDocumented(t *testing.T) {
1285:func TestHousekeepingRealmNamePatternDocumented(t *testing.T) {
1391:func TestPairingRealmNamePatternDocumented(t *testing.T) {
1433:func TestPairingUnregisterDeviceSemanticsDocumented(t *testing.T) {
1463:func operationDescription(t *testing.T, block []string) string {
1489:func operationBlock(t *testing.T, lines []string, operationID string) []string {
1513:func componentBlock(t *testing.T, lines []string, key string) []string {
1536:func containsLine(lines []string, want string) bool {

→ Read internal/swagger/swagger_test.go [offset=1463, limit=120]
Now I'll implement. First, extend the `ErrorDetail` schema:

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -1776,8 +1776,16 @@
       required: [detail]
       properties:
         detail:
           type: string
+        unexpected_keys:
+          type: array
+          items:
+            type: string
+          description: |
+            Keys present in an object-aggregated write document that match
+            no declared object leaf. Present only on the 400 `Unexpected
+            object key` rejection, sorted.
 
 FieldErrorsDetail:
   type: object
   required: [errors]

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -1808,8 +1808,32 @@
       example:
         errors:
           detail: Bad request
 
+BadObject:
+  description: |
+    Bad Request — the object-aggregated write document is rejected. A body
+    that is not a JSON object, is empty, or repeats a key answers the
+    generic `Bad request`; a document carrying a key that matches no
+    declared object leaf answers `Unexpected object key` and lists the
+    offending keys, sorted.
+  content:
+    application/json; charset=utf-8:
+      schema:
+        $ref: "#/components/schemas/ErrorDetail"
+      examples:
+        malformed:
+          summary: Not an object, empty, or duplicate key
+          value:
+            errors:
+              detail: Bad request
+        unexpected_key:
+          summary: A key matches no declared object leaf
+          value:
+            errors:
+              detail: Unexpected object key
+              unexpected_keys: [ghost]
+
 Unauthorized:
   description: Unauthorized — missing or invalid JWT.
   content:
     application/json; charset=utf-8:

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -1904,16 +1904,28 @@
           group_name:
             - can't be blank
 
 ValueTooLarge:
-  description: Unprocessable Entity — the value is larger than the per-value size limit (64 KiB).
+  description: |
+    Unprocessable Entity — the object-aggregated write document is
+    rejected. A value larger than the per-value size limit (64 KiB)
+    answers `Value size exceeds size limits`; a document that omits a key
+    the mapping declares required answers `Missing required mapping key`.
   content:
     application/json; charset=utf-8:
       schema:
         $ref: "#/components/schemas/ErrorDetail"
-      example:
-        errors:
-          detail: Value size exceeds size limits
+      examples:
+        value_too_large:
+          summary: Value exceeds the 64 KiB per-value size limit
+          value:
+            errors:
+              detail: Value size exceeds size limits
+        missing_required_mapping:
+          summary: Document omits a key the mapping declares required
+          value:
+            errors:
+              detail: Missing required mapping key
 
 UnprocessableEntity:
   description: Unprocessable Entity — attribute key not found on device.
   content:

Now repoint the six data PUT/POST 400 responses to the new `BadObject` component. Starting with the two by-alias ops:
← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -514,9 +514,9 @@
 responses:
   "200":
     description: Data set successfully.
   "400":
-    $ref: "#/components/responses/BadRequest"
+    $ref: "#/components/responses/BadObject"
   "401":
     $ref: "#/components/responses/Unauthorized"
   "403":
     $ref: "#/components/responses/Forbidden"

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -566,9 +566,9 @@
 responses:
   "200":
     description: Data published successfully.
   "400":
-    $ref: "#/components/responses/BadRequest"
+    $ref: "#/components/responses/BadObject"
   "401":
     $ref: "#/components/responses/Unauthorized"
   "403":
     $ref: "#/components/responses/Forbidden"

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -771,9 +771,9 @@
 responses:
   "200":
     description: Data set successfully.
   "400":
-    $ref: "#/components/responses/BadRequest"
+    $ref: "#/components/responses/BadObject"
   "401":
     $ref: "#/components/responses/Unauthorized"
   "403":
     $ref: "#/components/responses/Forbidden"

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -818,9 +818,9 @@
 responses:
   "200":
     description: Data published successfully.
   "400":
-    $ref: "#/components/responses/BadRequest"
+    $ref: "#/components/responses/BadObject"
   "401":
     $ref: "#/components/responses/Unauthorized"
   "403":
     $ref: "#/components/responses/Forbidden"

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -1332,9 +1332,9 @@
 responses:
   "200":
     description: Data set successfully.
   "400":
-    $ref: "#/components/responses/BadRequest"
+    $ref: "#/components/responses/BadObject"
   "401":
     $ref: "#/components/responses/Unauthorized"
   "403":
     $ref: "#/components/responses/Forbidden"

✗ Edit docs/api/astarte_appengine_api.yaml failed
Error: No changes to apply: oldString and newString are identical.

← Edit docs/api/astarte_appengine_api.yaml
Index: /root/astrate-mule/docs/api/astarte_appengine_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astarte_appengine_api.yaml
+++ /root/astrate-mule/docs/api/astarte_appengine_api.yaml
@@ -1380,9 +1380,9 @@
 responses:
   "200":
     description: Data published successfully.
   "400":
-    $ref: "#/components/responses/BadRequest"
+    $ref: "#/components/responses/BadObject"
   "401":
     $ref: "#/components/responses/Unauthorized"
   "403":
     $ref: "#/components/responses/Forbidden"

$ rg -n "BadObject|BadRequest|ValueTooLarge" docs/api/astarte_appengine_api.yaml
245:          $ref: "#/components/responses/BadRequest"
333:          $ref: "#/components/responses/BadRequest"
518:          $ref: "#/components/responses/BadObject"
528:          $ref: "#/components/responses/ValueTooLarge"
570:          $ref: "#/components/responses/BadObject"
580:          $ref: "#/components/responses/ValueTooLarge"
613:          $ref: "#/components/responses/BadRequest"
775:          $ref: "#/components/responses/BadObject"
785:          $ref: "#/components/responses/ValueTooLarge"
822:          $ref: "#/components/responses/BadObject"
832:          $ref: "#/components/responses/ValueTooLarge"
860:          $ref: "#/components/responses/BadRequest"
940:          $ref: "#/components/responses/BadRequest"
1064:          $ref: "#/components/responses/BadRequest"
1146:          $ref: "#/components/responses/BadRequest"
1336:          $ref: "#/components/responses/BadObject"
1346:          $ref: "#/components/responses/ValueTooLarge"
1384:          $ref: "#/components/responses/BadObject"
1394:          $ref: "#/components/responses/ValueTooLarge"
1423:          $ref: "#/components/responses/BadRequest"
1802:    BadRequest:
1812:    BadObject:
1907:    ValueTooLarge:

→ Read docs/api/astarte_appengine_api.yaml [limit=160, offset=1795]
→ Read docs/api/astarte_appengine_api.yaml [limit=40, offset=1768]
Now add a test that pins the documented bodies, following the repo's existing doc-drift test style:

← Edit internal/swagger/swagger_test.go
Index: /root/astrate-mule/internal/swagger/swagger_test.go
===================================================================
--- /root/astrate-mule/internal/swagger/swagger_test.go
+++ /root/astrate-mule/internal/swagger/swagger_test.go
@@ -459,8 +459,87 @@
 		t.Error("components.responses defines no BadRequest response")
 	}
 }
 
+// TestAppEngineObjectWriteErrorsDocumented guards that the six object-aggregated
+// data write operations (putData/publishData, by-alias and in-group) document
+// the two payload rejections their handler emits beyond the generic bodies.
+//
+// internal/appengine.writeError (internal/appengine/http.go) answers a
+// ReasonBadObject rejection — an object-aggregated document carrying a key that
+// matches no declared object leaf — with 400
+// {"errors":{"detail":"Unexpected object key","unexpected_keys":[...]}}
+// (writeBadObjectError, const detailUnexpectedObjectKey), and a
+// ReasonMissingRequired rejection — a document omitting a key the mapping
+// declares required — with 422 {"errors":{"detail":"Missing required mapping
+// key"}} (const detailMissingRequiredMapping). Both were implemented by
+// appengine-unexpected-object-key / appengine-missing-required-422 with no spec
+// update then. The six operations previously $ref'd the shared BadRequest for
+// 400 and ValueTooLarge for 422, neither of which carried these bodies: the 400
+// example never named the offending keys and the 422 example only showed the
+// value-size rejection. So the 400 now $refs the dedicated BadObject response
+// (it keeps the generic body and adds the key list), the shared ErrorDetail
+// schema declares the optional unexpected_keys array, and ValueTooLarge carries
+// both 422 examples. The three data DELETE operations keep the plain BadRequest;
+// this is the PUT/POST half only. This is the documentation half of the
+// behaviour pinned by TestWriteErrorTaxonomy in internal/appengine.
+func TestAppEngineObjectWriteErrorsDocumented(t *testing.T) {
+	b, err := docs.APIYAML.ReadFile("api/astarte_appengine_api.yaml")
+	if err != nil {
+		t.Fatalf("reading astarte_appengine_api.yaml: %v", err)
+	}
+	lines := strings.Split(string(b), "\n")
+
+	for _, op := range []string{
+		"putData", "publishData",
+		"putDataByAlias", "publishDataByAlias",
+		"putDataInGroup", "publishDataInGroup",
+	} {
+		block := operationBlock(t, lines, op)
+		for _, tc := range []struct {
+			status string
+			ref    string
+		}{
+			{"400", `          $ref: "#/components/responses/BadObject"`},
+			{"422", `          $ref: "#/components/responses/ValueTooLarge"`},
+		} {
+			body := strings.Join(responseBlock(t, block, tc.status), "\n")
+			if !strings.Contains(body, tc.ref) {
+				t.Errorf("operation %s %s does not $ref %s", op, tc.status, tc.ref)
+			}
+		}
+	}
+
+	schema := strings.Join(componentBlock(t, lines, "    ErrorDetail:"), "\n")
+	if !strings.Contains(schema, "unexpected_keys:") {
+		t.Error("ErrorDetail schema declares no unexpected_keys array")
+	}
+	if strings.Contains(schema, "required: [detail, unexpected_keys]") {
+		t.Error("ErrorDetail requires unexpected_keys; it belongs only to the unexpected-object-key 400")
+	}
+
+	bad := strings.Join(componentBlock(t, lines, "    BadObject:"), "\n")
+	for _, want := range []string{
+		"detail: Unexpected object key",
+		"unexpected_keys: [ghost]",
+		"detail: Bad request",
+	} {
+		if !strings.Contains(bad, want) {
+			t.Errorf("BadObject response carries no %q", want)
+		}
+	}
+
+	tooLarge := strings.Join(componentBlock(t, lines, "    ValueTooLarge:"), "\n")
+	for _, want := range []string{
+		"detail: Value size exceeds size limits",
+		"detail: Missing required mapping key",
+	} {
+		if !strings.Contains(tooLarge, want) {
+			t.Errorf("ValueTooLarge response carries no %q example", want)
+		}
+	}
+}
+
 // TestPairingErrorDetailExamples pins the error-detail examples in the pairing
 // spec's components.responses to the frozen canonical strings the wire emits
 // (astarteapi/envelope.go), so a generated client that copies an example does
 // not learn a detail that never goes over the wire. Phoenix renders "Bad

Now run the focused test and formatting:
$ gofmt -l internal/swagger/swagger_test.go; echo "gofmt-exit=$?"
gofmt-exit=0

$ go test ./internal/swagger/... 2>&1 | tail -30
ok  	github.com/astrate-platform/astrate/internal/swagger	0.081s

$ ls examples/pokemon-agent/docs/DESIGN.md docs/site/swagger.md 2>&1; command -v mkdocs && mkdocs --version 2>&1 | head -2
docs/site/swagger.md
examples/pokemon-agent/docs/DESIGN.md

$ go vet ./... 2>&1 | tail -20; echo "vet-exit=$?"
vet-exit=0

$ go test ./... 2>&1 | tail -40; echo "test-exit=${PIPESTATUS[0]}"
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
/bin/sh: 1: Bad substitution

$ make -C docs build 2>&1 | tail -25
make: Entering directory '/root/astrate-mule/docs'
cp DESIGN.md site/DESIGN.md
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'

$ command -v mkdocs; command -v python3; python3 -c "import yaml; print('pyyaml ok')" 2>&1; ls docs/site/api 2>&1 | head; ls .venv docs/.venv 2>/dev/null
/usr/bin/python3
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml; print('pyyaml ok')
    ^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'
astarte_appengine_api.yaml
astarte_housekeeping_api.yaml
astarte_pairing_api.yaml
astarte_realm_management_api.yaml
astrate_native_api.yaml

$ git status --porcelain; echo "=== ignored check ==="; git check-ignore docs/site/api/astarte_appengine_api.yaml docs/site/DESIGN.md 2>&1; echo "=== yaml dep ==="; rg -n "yaml" go.mod
M docs/api/astarte_appengine_api.yaml
 M internal/swagger/swagger_test.go
=== ignored check ===
docs/site/api/astarte_appengine_api.yaml
docs/site/DESIGN.md
=== yaml dep ===
93:	go.yaml.in/yaml/v2 v2.4.2 // indirect
99:	gopkg.in/yaml.v3 v3.0.1 // indirect

$ command -v yq ruby node npx pip3 2>&1; ls /root/go/pkg/mod/gopkg.in/ 2>/dev/null; ls /root/go/pkg/mod/go.yaml.in/ 2>/dev/null
check.v1@v0.0.0-20161208181325-20d25e280405
check.v1@v1.0.0-20201130134442-10cb98267c6c
ini.v1@v1.67.0
yaml.v2@v2.2.8
yaml.v3@v3.0.1
yaml

$ mkdir -p /tmp/opencode/yamlcheck && cat > /tmp/opencode/yamlcheck/go.mod <<'EOF'
module yamlcheck

go 1.26

require gopkg.in/yaml.v3 v3.0.1
EOF
cat > /tmp/opencode/yamlcheck/main.go <<'EOF'
package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	failed := false
	for _, path := range os.Args[1:] {
		b, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("READ  %s: %v\n", path, err)
			failed = true
			continue
		}
		var root map[string]any
		if err := yaml.Unmarshal(b, &root); err != nil {
			fmt.Printf("PARSE %s: %v\n", path, err)
			failed = true
			continue
		}
		if _, ok := root["openapi"]; !ok {
			fmt.Printf("NOT-OPENAPI %s\n", path)
			failed = true
			continue
		}
		fmt.Printf("OK    %s\n", path)
	}
	if failed {
		os.Exit(1)
	}
	_ = strings.TrimSpace
}
EOF
cd /tmp/opencode/yamlcheck && GOFLAGS=-mod=mod go mod tidy >/dev/null 2>&1; go run . /root/astrate-mule/docs/api/astarte_appengine_api.yaml /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml /root/astrate-mule/docs/api/astarte_pairing_api.yaml /root/astrate-mule/docs/api/astarte_realm_management_api.yaml /root/astrate-mule/docs/api/astrate_native_api.yaml; echo "exit=$?"
OK    /root/astrate-mule/docs/api/astarte_appengine_api.yaml
OK    /root/astrate-mule/docs/api/astarte_housekeeping_api.yaml
OK    /root/astrate-mule/docs/api/astarte_pairing_api.yaml
OK    /root/astrate-mule/docs/api/astarte_realm_management_api.yaml
OK    /root/astrate-mule/docs/api/astrate_native_api.yaml
exit=0

$ go test -race ./internal/swagger/... 2>&1 | tail -5; echo "exit=${PIPESTATUS[0]}"
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/swagger	0.027s
FAIL
/bin/sh: 1: Bad substitution

$ gofmt -l . 2>&1; echo "exit=$?"
exit=0

$ cat docs/requirements.txt 2>/dev/null; pip3 --version 2>&1 | head -1
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10
/bin/sh: 1: pip3: not found

$ git diff --stat; echo "==="; git diff docs/api/astarte_appengine_api.yaml | head -120
docs/api/astarte_appengine_api.yaml | 64 +++++++++++++++++++++++++-----
 internal/swagger/swagger_test.go    | 79 +++++++++++++++++++++++++++++++++++++
 2 files changed, 133 insertions(+), 10 deletions(-)
===
diff --git a/docs/api/astarte_appengine_api.yaml b/docs/api/astarte_appengine_api.yaml
index 442a3f6..3d04a5f 100644
--- a/docs/api/astarte_appengine_api.yaml
+++ b/docs/api/astarte_appengine_api.yaml
@@ -515,7 +515,7 @@ paths:
         "200":
           description: Data set successfully.
         "400":
-          $ref: "#/components/responses/BadRequest"
+          $ref: "#/components/responses/BadObject"
         "401":
           $ref: "#/components/responses/Unauthorized"
         "403":
@@ -567,7 +567,7 @@ paths:
         "200":
           description: Data published successfully.
         "400":
-          $ref: "#/components/responses/BadRequest"
+          $ref: "#/components/responses/BadObject"
         "401":
           $ref: "#/components/responses/Unauthorized"
         "403":
@@ -772,7 +772,7 @@ paths:
         "200":
           description: Data set successfully.
         "400":
-          $ref: "#/components/responses/BadRequest"
+          $ref: "#/components/responses/BadObject"
         "401":
           $ref: "#/components/responses/Unauthorized"
         "403":
@@ -819,7 +819,7 @@ paths:
         "200":
           description: Data published successfully.
         "400":
-          $ref: "#/components/responses/BadRequest"
+          $ref: "#/components/responses/BadObject"
         "401":
           $ref: "#/components/responses/Unauthorized"
         "403":
@@ -1333,7 +1333,7 @@ paths:
         "200":
           description: Data set successfully.
         "400":
-          $ref: "#/components/responses/BadRequest"
+          $ref: "#/components/responses/BadObject"
         "401":
           $ref: "#/components/responses/Unauthorized"
         "403":
@@ -1381,7 +1381,7 @@ paths:
         "200":
           description: Data published successfully.
         "400":
-          $ref: "#/components/responses/BadRequest"
+          $ref: "#/components/responses/BadObject"
         "401":
           $ref: "#/components/responses/Unauthorized"
         "403":
@@ -1777,6 +1777,14 @@ components:
           properties:
             detail:
               type: string
+            unexpected_keys:
+              type: array
+              items:
+                type: string
+              description: |
+                Keys present in an object-aggregated write document that match
+                no declared object leaf. Present only on the 400 `Unexpected
+                object key` rejection, sorted.
 
     FieldErrorsDetail:
       type: object
@@ -1801,6 +1809,30 @@ components:
             errors:
               detail: Bad request
 
+    BadObject:
+      description: |
+        Bad Request — the object-aggregated write document is rejected. A body
+        that is not a JSON object, is empty, or repeats a key answers the
+        generic `Bad request`; a document carrying a key that matches no
+        declared object leaf answers `Unexpected object key` and lists the
+        offending keys, sorted.
+      content:
+        application/json; charset=utf-8:
+          schema:
+            $ref: "#/components/schemas/ErrorDetail"
+          examples:
+            malformed:
+              summary: Not an object, empty, or duplicate key
+              value:
+                errors:
+                  detail: Bad request
+            unexpected_key:
+              summary: A key matches no declared object leaf
+              value:
+                errors:
+                  detail: Unexpected object key
+                  unexpected_keys: [ghost]
+
     Unauthorized:
       description: Unauthorized — missing or invalid JWT.
       content:
@@ -1873,14 +1905,26 @@ components:
                 - can't be blank
 
     ValueTooLarge:
-      description: Unprocessable Entity — the value is larger than the per-value size limit (64 KiB).
+      description: |
+        Unprocessable Entity — the object-aggregated write document is
+        rejected. A value larger than the per-value size limit (64 KiB)
+        answers `Value size exceeds size limits`; a document that omits a key
+        the mapping declares required answers `Missing required mapping key`.
       content:
         application/json; charset=utf-8:
           schema:
             $ref: "#/components/schemas/ErrorDetail"
-          example:
-            errors:

$ git diff docs/api/astarte_appengine_api.yaml | tail -30
- can't be blank
 
     ValueTooLarge:
-      description: Unprocessable Entity — the value is larger than the per-value size limit (64 KiB).
+      description: |
+        Unprocessable Entity — the object-aggregated write document is
+        rejected. A value larger than the per-value size limit (64 KiB)
+        answers `Value size exceeds size limits`; a document that omits a key
+        the mapping declares required answers `Missing required mapping key`.
       content:
         application/json; charset=utf-8:
           schema:
             $ref: "#/components/schemas/ErrorDetail"
-          example:
-            errors:
-              detail: Value size exceeds size limits
+          examples:
+            value_too_large:
+              summary: Value exceeds the 64 KiB per-value size limit
+              value:
+                errors:
+                  detail: Value size exceeds size limits
+            missing_required_mapping:
+              summary: Document omits a key the mapping declares required
+              value:
+                errors:
+                  detail: Missing required mapping key
 
     UnprocessableEntity:
       description: Unprocessable Entity — attribute key not found on device.

All changes are in place, the YAML parses, and the full non-race gate passes. The task is done.

```
Done: Documented the two object-aggregated write rejections on the six data PUT/POST ops — 400 `Unexpected object key` now on a dedicated BadObject response with the optional `unexpected_keys` array added to the shared ErrorDetail schema, and 422 `Missing required mapping key` added as a second example on ValueTooLarge — with a pinning test.
Files: docs/api/astarte_appengine_api.yaml; internal/swagger/swagger_test.go
Verified: go vet ./... -> pass; go test ./... (non-race, Pi gate) -> pass; go test ./internal/swagger/... -> pass (incl. new TestAppEngineObjectWriteErrorsDocumented); gofmt -l . -> clean; yaml.v3 parse of all five docs/api/*.yaml -> pass; go test -race ./internal/swagger/... -> fail "ThreadSanitizer: unsupported VMA range" (known 39-bit Pi limitation; nothing concurrent touched); make -C docs build -> fail only at `mkdocs build` (mkdocs not installed here; sync half ran and regenerated the gitignored docs/site copies unchanged in content)
Unsure: Took the task's "or a dedicated response" branch for the 400 because the shared BadRequest is pinned by TestAppEngineErrorDetailExamples and referenced by every non-data op, so the unexpected-keys body cannot live there — the six ops' 400 now $ref the new BadObject and keep BadRequest on the delete/other ops. Kept the 422 under the existing ValueTooLarge name (per "beside the value-size one") and broadened its description to name both rejections instead of renaming it.
Follow-ups: none
```
