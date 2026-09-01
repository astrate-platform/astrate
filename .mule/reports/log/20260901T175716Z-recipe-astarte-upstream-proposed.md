slug: recipe-astarte-upstream
verdict: proposed
at:  f4ad882
ran: 2026-09-01T17:57:16Z on DietPi in 478s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

The COMPATIBILITY.md doesn't name an explicit upstream version tag. The latest stable release is v1.3.3, and there are release candidates for v1.4.0. Let me check what v1.3.3 changed, and also look at the v1.4.0-rc.5 notes since that's the newest.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
Sixth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.4

### Fixed

- [astarte_data_updater_plant] Make object values only validate type on selected interface
- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated

v1.3.3 has empty release notes. Let me look at the v1.4.0 RCs more carefully to understand what changed since v1.2.x (which is what the conformance measurements reference).
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.4") | .body'
Fifth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.3

### Added

- [config] Initialize `config` library. All external services declared using the `config` library helpers can be configured either using the `_URL` environment variable or by using single component variables `_SCHEME`, `_HOST`, `_PORT`, `_PATH`, `_QUERY`, `_FRAGMENT` (`_URL` has priority), together with normal `_SSL_*` variables.
- [secrets] Use `config` for vault configuration. Vault can now be configured using component variables instead of `ASTARTE_VAULT_URL`:
  - `ASTARTE_VAULT_SCHEME`: scheme part of the URL (binary, default: `"http"`, or `"https"` if `ASTARTE_VAULT_SSL_ENABLED` is `true`)
  - `ASTARTE_VAULT_HOST`: host part of the URL (binary, default: `"localhost"`)
  - `ASTARTE_VAULT_PORT`: port part of the url (integer, default `8200`)
  - `ASTARTE_VAULT_PATH`: path part of the url (string)
  - `ASTARTE_VAULT_QUERY`: query part of the url (string)
  - `ASTARTE_VAULT_FRAGMENT`: fragment part of the url (string)
- [fdo] Use `config` for rendezvous configuration. The rendezvous can now be configured using component variables instead of `PAIRING_FDO_RENDEZVOUS_URL`:
  - `PAIRING_FDO_RENDEZVOUS_SCHEME`: scheme part of the URL (binary, default: `"http"`, or `"https"` if `PAIRING_FDO_RENDEZVOUS_SSL_ENABLED` is `true`)
  - `PAIRING_FDO_RENDEZVOUS_HOST`: host part of the URL (binary, default: `"localhost"`)
  - `PAIRING_FDO_RENDEZVOUS_PORT`: port part of the url (integer, default `8041`)
  - `PAIRING_FDO_RENDEZVOUS_PATH`: path part of the url (string)
  - `PAIRING_FDO_RENDEZVOUS_QUERY`: query part of the url (string)
  - `PAIRING_FDO_RENDEZVOUS_FRAGMENT`: fragment part of the url (string)

### Fixed

- [astarte_data_updater_plant] Ensure RPC server is always available to clients. Resolved the issue where a temporary disconnection and reconnection of data_updater_plant to the cluster would make the RPC server inaccessible.

### Changed

- [astarte_data_updater_plant] Ensure memory is properly garbage collected

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.3") | .body'
Fourth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.2

### Added

- [secrets] Configurable base namespace for Vault using `ASTARTE_VAULT_BASE_NAMESPACE`: all created namespaces will be children of the base namespace.

### Fixed

- [fdo] Support for Vault and Vault managed
- [fdo] Properly handle chunked service info messages
- [astarte_data_updater_plant] Avoid clogging AMQP queues when receiving invalid device messages
- [astarte_data_updater_plant] Restore correct type validation checks for incoming data

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.2") | .body'
Third Astarte v1.4.0 release candidate.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body' ; echo "---EXIT: $?"
---EXIT: 0

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.1") | .body'
Second Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.0

### Added

- [astarte_housekeeping] Added HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION to set default, instance wise, realm default realm data retention, expressed in seconds

### Fixed

- [fdo] Do not crash when using rotated keys as owner keys
- [fdo] Delete devices from unfinished sessions
- [fdo] Delete old sessions from device
- [fdo] Properly respond to service info messages after IsDone=True
- [astarte_realm_management] Properly show default values for all mapping parameters

### Changed

- Enable FDO by default

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.2") | .body'
Third Astarte v1.4.0 release candidate.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(startswith("v1.2")) | "\(.tag_name)\t\(.published_at)"' | head -5
startswith("v1.2") cannot be applied to: object ({"assets":[],"assets_url" ...})

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z
v1.4.0-rc.1	2026-07-14T08:53:27Z
v1.3.2	2026-07-14T10:07:26Z
v1.3.1	2026-07-13T11:50:05Z
v1.3.0	2026-05-06T15:30:50Z
v1.2.2	2026-04-27T16:39:28Z
v1.4.0-rc.0	2026-04-08T12:32:33Z
v1.3.0-rc.2	2026-04-08T16:34:59Z
v1.2.2-rc.0	2026-04-08T16:31:03Z
v1.2.1	2026-03-12T16:11:20Z
v1.2.1-rc.1	2026-02-13T10:36:20Z
v1.3.0-rc.1	2026-01-26T15:58:33Z
v1.3.0-rc.0	2025-11-21T14:58:46Z
v1.2.1-rc.0	2025-08-27T07:03:28Z
v1.2.1-alpha.0	2025-04-10T15:31:04Z
v1.2.0	2024-07-02T14:09:07Z
v1.2.0-rc.0	2024-06-11T13:02:33Z
v1.0.6	2024-04-23T21:15:57Z
v1.1.1	2023-11-15T15:29:57Z
v1.0.5	2023-09-26T10:27:57Z
v1.1.0	2023-06-20T18:59:09Z
v1.1.0-rc.0	2023-06-09T15:49:44Z
v1.1.0-alpha.0	2022-11-24T16:38:37Z
v1.0.4	2022-10-25T16:50:21Z
v1.0.3	2022-07-05T13:57:17Z
v1.0.2	2022-04-01T17:18:15Z

$ gh api 'repos/astarte-platform/astarte/releases?per_page=100' --jq '.[] | select(.tag_name=="v1.4.0-rc.0") | .body'
First Astarte v1.4.0 release candidate.

## Changes since v1.3

### Added

- Add required flag for mappings of object aggregated interfaces
- [astarte_pairing] Support for FIDO Device Onboard v1.1 for device registration. The following environment variables have been added for the Vault connectivity:
  - `ASTARTE_VAULT_URL` - The URL to access Vault.
  - `ASTARTE_VAULT_AUTHENTICATION_MECHANISM` - The mechanism to use for authenticating with Vault (default: `""`, allowed: \["token"])
  - `ASTARTE_VAULT_TOKEN` - Token to authenticate with Vault, used if `ASTARTE_VAULT_AUTHENTICATION_MECHANISM` is `"token"` (default: `""`)
  - `ASTARTE_VAULT_SSL_CA_FILE` - path to the CA certificate file for the vault TLS connection; when not specified, the bundled cURL certificate bundle will be used
  - `ASTARTE_VAULT_SSL_CUSTOM_SNI` - custom SNI hostname for the vault TLS connection; when not specified, the vault hostname will be used
  - `ASTARTE_VAULT_SSL_DISABLE_SNI` - Disable Server Name Indication for Vault. (default: `false`)
  - `ASTARTE_VAULT_SSL_ENABLED` - Whether to enable SSL for the Vault connection. (default `false`)
- [astarte_realm_management] Allow listing interfaces with mappings using the `detailed=true` parameter
- [astarte_import] Added support for data types: `doublearray`, `integerarray`,
  `booleanarray`, `longintegerarray`, `stringarray`, `datetimearray`, `binaryblobarray`.

### Changed

- [astarte_data_updater_plant] Use mississippi consumer for data updater processes

Now let me check the current state of the todo queue to avoid proposing duplicates, and see if there's any version target referenced elsewhere.
→ Read .mule/todo.md

$ rg -n "v1\.[0-9]" .mule/todo.md .mule/*.md 2>/dev/null | head -40
.mule/todo.md:99:- [!] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto] — BLOCKED: gates failed
.mule/todo.md:100:- [!] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
.mule/todo.md:101:- [!] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto] — BLOCKED: gates failed
.mule/todo.md:102:- [!] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
.mule/todo.md:103:- [!] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto] — BLOCKED: gates failed
.mule/todo.md:104:- [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
.mule/todo.md:105:- [!] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto] — BLOCKED: gates failed
.mule/todo.md:99:- [!] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto] — BLOCKED: gates failed
.mule/todo.md:100:- [!] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
.mule/todo.md:101:- [!] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto] — BLOCKED: gates failed
.mule/todo.md:102:- [!] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
.mule/todo.md:103:- [!] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto] — BLOCKED: gates failed
.mule/todo.md:104:- [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
.mule/todo.md:105:- [!] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto] — BLOCKED: gates failed

$ git log --oneline -15
f4ad882 mule: blocked issue-68-async-operation-accepted
e0f8823 mule: blocked issue-91-pairing-health-serve-root
733f965 mule: refill the queue
b69e32e mule: recipe github-issues ran (167s)
b8424d8 mule: refill the queue
2777546 mule: recipe milestones ran (607s)
3a70ddb mule: refill the queue
22875f1 mule: recipe hygiene timed out
5006f6e mule: blocked docs-sync-pairing-status-enum
49fcd62 mule: refill the queue
45171b9 mule: recipe docs-sync ran (203s)
19ad0fe mule: blocked probe-properties-on-connect-encoding
8a983a3 mule: blocked probe-binaryblob-validation
25560b4 mule: blocked probe-mqtt-capabilities-declaration
21f74d8 mule: blocked compat-note-v133

$ ls .mule/reports/ 2>/dev/null | head -40
issue-13.md
issue-20.md
log
race-check.md

The queue already carries proposals from previous runs of this recipe (v1.3.0, v1.3.3, v1.4.0-rc.3, v1.4.0-rc.5). Let me check git history and past reports to see exactly what was already proposed, so I don't duplicate.
$ rg -rl "astarte-upstream|upstream" .mule/reports/ 2>/dev/null | head
.mule/reports/issue-13.md:✱ Grep "lErrorNames" 7 matches
.mule/reports/issue-13.md:title:	audit-dashboard-error-names: compare errorname.go's lErrorNames against the current astarte-dashboard bundle
.mule/reports/issue-13.md:`lErrorNames` in `internal/engine/triggers/errorname.go`. Report the diff in a comment
.mule/reports/issue-13.md:Found by a cross-project survey comparing Astrate against l astarte-platform
.mule/reports/issue-13.md:- **`lErrorNames` in `internal/engine/triggers/errorname.go`**: exactly 15 names.
.mule/reports/issue-13.md:Done: audited astarte-dashboard v1.3.2 `deviceErrorNames` against `lErrorNames` — no diff
.mule/reports/log/20260831T201035Z-recipe-docs-sync-proposed.md: - [!] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (l device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto] — BLOCKED: gates failed
.mule/reports/log/20260831T201035Z-recipe-docs-sync-proposed.md: - [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as l v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
.mule/reports/log/20260831T201035Z-recipe-docs-sync-proposed.md: - [!] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as l v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto] — BLOCKED: gates failed
.mule/reports/log/20260831T165224Z-control-producer-properties-compression-blocked.md:Let me read the l `control_handler.ex` reference and the existing test file.

$ git show 733f965 --stat | head -30
commit 733f9651fd4a8e4f1e877c7a53e3df81914c7e53
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Tue Sep 1 19:02:01 2026 +0200

    mule: refill the queue

 .mule/cron.log | 276 +++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 .mule/todo.md  |   2 +
 2 files changed, 278 insertions(+)

$ git log --oneline --all | rg -i "upstream|refill" | head -20
733f965 mule: refill the queue
b8424d8 mule: refill the queue
3a70ddb mule: refill the queue
49fcd62 mule: refill the queue
7b03846 mule: recipe astarte-upstream ran (224s)
ab31319 docs(compat): measured write taxonomy + unknown-realm deviation; fix Internal server error casing to measured upstream form
ad4eb45 fix(appengine): server-write error taxonomy matches measured upstream (#57)
cda4010 feat(http): per-service version endpoints, upstream-parity (#77)
907b6e7 test(conformance): record fase-1 upstream verify fixtures (#69/#57); fix Bad request casing to measured upstream form
071d0a2 mule: recipe astarte-upstream ran (396s)
dbd0267 docs: upstream-experimental register + APICompatVersion bump rule (#67, #68, #90)
ed0342d feat(triggers): upstream wildcard semantics + introspection interface rules (#70)
edf4e1b feat(triggers): upstream-measured HTTP action validation limits (#63)
0b5a35c feat(realm,interfaceschema): upstream-measured install/update/delete error taxonomy (issue #62)
7f3a600 conformance: measure group WATCH, server-owned and size limits against upstream
081abe4 milestones: redefine v3.0 as upstream Astarte feature parity
83568dd mule: recipe astarte-upstream ran (333s)
3a4b693 mule: recipe astarte-upstream ran (431s)
3a746d1 mule: recipe astarte-upstream ran (485s)
d145c7d mule: propose upstream v1.3.0 wire-visible changes + compat-note-v1.3.2

$ git show 733f965 -- .mule/todo.md | head -40
commit 733f9651fd4a8e4f1e877c7a53e3df81914c7e53
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Tue Sep 1 19:02:01 2026 +0200

    mule: refill the queue

diff --git a/.mule/todo.md b/.mule/todo.md
index bc44d10..f189f2f 100644
--- a/.mule/todo.md
+++ b/.mule/todo.md
@@ -104,3 +104,5 @@ worth having. Before any of it reaches `main`:
 - [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
 - [!] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto] — BLOCKED: gates failed
 - [!] docs-sync-pairing-status-enum: in docs/api/astarte_pairing_api.yaml, fix the PairingInfo.status enum (line 367): it lists `confirmed, pending, denied, expired` but the handler only ever emits `confirmed`, `pending`, or `inhibited` (internal/pairing/service.go:289-297, via internal/pairing/http.go:211) — drop the dead `denied`/`expired` and add the undocumented `inhibited`. [auto] — BLOCKED: gates failed
+- [ ] issue-91-pairing-health-serve-root: in internal/pairing/http.go, add the unauthenticated `GET /pairing/health` route (no realm segment) sharing the handler/payload of the existing `GET /pairing/v1/{realm}/health`, keeping the v1 route as-is; tests for both paths (closes #91)
+- [ ] issue-68-async-operation-accepted: on the mutating endpoints upstream surfaces `async_operation` on (housekeeping realm create/delete; realm-management interface install/update/delete, trigger/policy delete in internal/realm/http.go), accept and ignore `?async_operation=false`, with unparseable/`true` values also not changing behaviour (Astrate stays always-sync); tests cover the flag parsing (closes #68)

$ git show 7b03846 --stat | head -20
commit 7b0384603d4f9515a3e89a457002a4aef3761c8b
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 31 04:12:36 2026 +0200

    mule: recipe astarte-upstream ran (224s)

 ...831T021236Z-recipe-astarte-upstream-proposed.md | 280 +++++++++++++++++++++
 1 file changed, 280 insertions(+)

$ git show 7b03846 -- .mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md | head -120
commit 7b0384603d4f9515a3e89a457002a4aef3761c8b
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 31 04:12:36 2026 +0200

    mule: recipe astarte-upstream ran (224s)

diff --git a/.mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md b/.mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md
new file mode 100644
index 0000000..b3a6d53
--- /dev/null
+++ b/.mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md
@@ -0,0 +1,280 @@
+slug: recipe-astarte-upstream
+verdict: proposed
+at:  6e30a1c
+ran: 2026-08-31T02:12:36Z on DietPi in 224s
+
+
+> build · big-pickle
+
+I'll start by reading the two files specified.
+→ Read .mule/MULE.md
+→ Read .mule/recipes/astarte-upstream.md
+Now let me read the compatibility doc and check upstream releases in parallel.
+→ Read docs/COMPATIBILITY.md
+
+$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
+v1.4.0-rc.5	2026-08-20T12:28:04Z
+v1.3.3	2026-08-07T14:33:58Z
+v1.4.0-rc.4	2026-08-06T13:40:20Z
+v1.4.0-rc.3	2026-07-31T10:48:39Z
+v1.4.0-rc.2	2026-07-14T10:16:22Z
+
+→ Read .mule/todo.md
+Upstream is ahead: newest tag is `v1.4.0-rc.5` (stable also has `v1.3.3`), while COMPATIBILITY.md records upstream v1.2.0/v1.2.2. Let me check the todo queue for prior runs and pull the release notes.
+
+$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
+Sixth Astarte v1.4.0 release candidate.
+
+## Changes since v1.4.0-rc.4
+
+### Fixed
+
+- [astarte_data_updater_plant] Make object values only validate type on selected interface
+- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated
+
+$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'
+The queue already covers v1.3.0/v1.3.3/v1.4.0-rc.3 and rc.5's binaryblob fix, but not rc.5's other change: "Make object values only validate type on selected interface." Let me look at that upstream change before proposing anything.
+
+$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | {tag_name, body, published_at}'
+{"body":"","published_at":"2026-08-07T14:33:58Z","tag_name":"v1.3.3"}
+
+$ gh search commits --owner astarte-platform 'Make object values only validate type on selected interface' --json sha,repository,commit --jq '.[] | "\(.repository.full_name)\t\(.sha[0:9])\t\(.commit.message | split("\n")[0])"'
+(no output)
+
+$ gh search commits --owner astarte-platform 'object values only validate type on selected interface' --json sha,repository,commit 2>/dev/null; echo "---"; gh search commits --owner astarte-platform 'object values' --json sha,repository,commit.url,commit.author --jq '.[] | "\(.repository.full_name)\t\(.sha[0:9])\t\(.commit.message | split("\n")[0])"' 2>&1 | head -20
+[]
+---
+Unknown JSON field: "commit.url"
+Available fields:
+  author
+  commit
+  committer
+  id
+  parents
+  repository
+  sha
+  url
+
+$ gh api repos/astarte-platform/astarte_data_updater_plant/commits --jq '.[0:15][] | "\(.sha[0:9])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"'
+8739d915f	2019-12-20T15:06:00Z	This repository has been archived.
+69612e95f	2019-12-18T16:49:09Z	Forward port changes from v0.10
+28ae84960	2019-12-09T11:32:32Z	Merge branch 'prepare-0.10.2' into 'release-0.10'
+c0087ddce	2019-12-09T11:21:25Z	Prepare 0.10.2 release
+64da6f368	2019-11-28T14:49:05Z	Merge pull request #17 from rbino/scylla-beta-support
+a40acb406	2019-11-28T14:18:26Z	.travis.yml: remove Scylla CI from allowed failures
+9a3198a26	2019-11-28T14:12:40Z	queries: don't use shorthand syntax to update maps
+bc11f0383	2019-11-27T15:16:08Z	Merge branch 'ci-upgrade-services' into 'master'
+59bd736d7	2019-11-27T10:49:58Z	.gitlab-ci.yml: upgrade Cassandra and RabbitMQ services
+a4156ce2b	2019-11-26T14:20:37Z	Forward port bugfixes from v0.10
+140f80393	2019-11-26T11:34:43Z	Merge branch 'fix-13' into 'release-0.10'
+745d9b792	2019-11-26T11:22:01Z	impl: load device volatile triggers when they're installed
+afe9d5b4f	2019-11-26T11:21:20Z	Merge branch 'fix-15' into 'release-0.10'
+466ca7c70	2019-11-26T10:49:12Z	impl: fix volatile triggers deletion bug
+6fcd3f5f0	2019-11-19T16:24:14Z	Merge pull request #12 from rbino/ci-wait-for-cassandra
+
+The data updater code now lives in the astarte monorepo. Let me find the actual commit for the rc.5 object-validation change.
+$ gh api 'repos/astarte-platform/astarte/commits?since=2026-08-15T00:00:00Z&until=2026-08-21T00:00:00Z' --jq '.[] | "\(.sha[0:9])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"'
+739724991	2026-08-19T13:06:29Z	chore: Update exandra
+643d06107	2026-08-19T12:36:15Z	ci: add release-please for automated versioning and changelog
+6e65c8aca	2026-08-17T12:25:10Z	refactor(dup): Move encryption logic into library (#2123)
+88004a2cf	2026-08-14T09:21:28Z	chore(core): use the local astarte_core across the monorepo
+
+$ gh search commits --owner astarte-platform --repository astarte 'binaryblob data is correctly validated' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo "==="; gh search commits --owner astarte-platform --repository astarte 'object values only validate type' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'
+unknown flag: --repository
+
+Usage:  gh search commits [<query>] [flags]
+
+Flags:
+      --author string            Filter by author
+      --author-date date         Filter based on authored date
+      --author-email string      Filter on author email
+      --author-name string       Filter on author name
+      --committer string         Filter by committer
+      --committer-date date      Filter based on committed date
+      --committer-email string   Filter on committer email
+      --committer-name string    Filter on committer name
+      --hash string              Filter by commit hash
+  -q, --jq expression            Filter JSON output using a jq expression
+      --json fields              Output JSON with the specified fields
+  -L, --limit int                Maximum number of commits to fetch (default 30)
+      --merge                    Filter on merge commits
+      --order string             Order of commits returned, ignored unless '--sort' flag is specified: {asc|desc} (default "desc")
+      --owner strings            Filter on repository owner
+      --parent string            Filter by parent hash
+  -R, --repo OWNER/REPO          Filter on repository, in OWNER/REPO format
+      --sort string              Sort fetched commits: {author-date|committer-date} (default "best-match")
+  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
+      --tree string              Filter by tree hash
+      --visibility strings       Filter based on repository visibility: {public|private|internal}

$ git show 7b03846 -- .mule/reports/log/20260831T021236Z-recipe-astarte-upstream-proposed.md | sed -n '120,280p'
+      --visibility strings       Filter based on repository visibility: {public|private|internal}
+  -w, --web                      Open the search query in the web browser
+  
+===
+unknown flag: --repository
+
+Usage:  gh search commits [<query>] [flags]
+
+Flags:
+      --author string            Filter by author
+      --author-date date         Filter based on authored date
+      --author-email string      Filter on author email
+      --author-name string       Filter on author name
+      --committer string         Filter by committer
+      --committer-date date      Filter based on committed date
+      --committer-email string   Filter on committer email
+      --committer-name string    Filter on committer name
+      --hash string              Filter by commit hash
+  -q, --jq expression            Filter JSON output using a jq expression
+      --json fields              Output JSON with the specified fields
+  -L, --limit int                Maximum number of commits to fetch (default 30)
+      --merge                    Filter on merge commits
+      --order string             Order of commits returned, ignored unless '--sort' flag is specified: {asc|desc} (default "desc")
+      --owner strings            Filter on repository owner
+      --parent string            Filter by parent hash
+  -R, --repo OWNER/REPO          Filter on repository, in OWNER/REPO format
+      --sort string              Sort fetched commits: {author-date|committer-date} (default "best-match")
+  -t, --template string          Format JSON output using a Go template; see "gh help formatting"
+      --tree string              Filter by tree hash
+      --visibility strings       Filter based on repository visibility: {public|private|internal}
+  -w, --web                      Open the search query in the web browser
+
+$ gh search commits -R astarte-platform/astarte 'binaryblob data is correctly validated' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo "==="; gh search commits -R astarte-platform/astarte 'object values only validate type on selected interface' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'
+aedc41483	fix(dup): ensure binaryblob data is correctly validated
+===
+
+$ gh search commits -R astarte-platform/astarte 'object values only validate type' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo "===loose==="; gh search commits -R astarte-platform/astarte 'object values' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"' | head -10
+===loose===
+cc89858b3	DUP: do not crash on binaryblob insertion
+219b9c79d	Merge branch 'object-aggregation' into 'master'
+084db3757	Handle aggregate interface in live events
+adf94edf8	Fix data parsing on datastream object interfaces w/ parametric endpoints
+
+$ gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.4...v1.4.0-rc.5 --jq '.commits[] | "\(.sha[0:9])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"'
+d8592091e	2026-08-07T06:55:11Z	fix(dup): stop remapping interface loading failed error (#2119)
+a816eaf0e	2026-08-07T06:57:11Z	chore: forward port release-1.3
+aedc41483	2026-08-07T07:21:33Z	fix(dup): ensure binaryblob data is correctly validated
+3e70d63f8	2026-08-07T14:28:36Z	chore: prepare v1.3.3 (#2128)
+fdc6d29b8	2026-08-19T08:14:18Z	fix(dup): Make object values only validate type on selected interface (#2141)
+4141a2ae4	2026-08-19T13:08:32Z	chore: forward-port release-1.3 into release-1.4
+3484096d2	2026-08-20T12:10:25Z	Prepare release v1.4.0-rc.5 (#2147)
+
+The new item is commit `fdc6d29b8` (#2141). Let me read the diff.
+$ gh api repos/astarte-platform/astarte/commits/fdc6d29b8 --jq '{sha, message: .commit.message, files: [.files[] | {filename, status: .status, additions, deletions}]}'
+{"files":[{"additions":1,"deletions":1,"filename":"apps/astarte_appengine_api/mix.lock","status":"modified"},{"additions":2,"deletions":1,"filename":"apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex","status":"modified"},{"additions":40,"deletions":0,"filename":"apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_test.exs","status":"modified"},{"additions":23,"deletions":0,"filename":"apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/interface_test.exs","status":"modified"},{"additions":58,"deletions":0,"filename":"apps/astarte_data_updater_plant/test/support/cases/device.ex","status":"modified"},{"additions":1,"deletions":1,"filename":"apps/astarte_pairing/mix.lock","status":"modified"}],"message":"fix(dup): Make object values only validate type on selected interface (#2141)\n\nSigned-off-by: Francesco Noacco \u003cfrancesco.noacco@secomind.com\u003e","sha":"fdc6d29b84ec06a12316340b60f1b43491ca3d5c"}
+
+$ gh api repos/astarte-platform/astarte/commits/fdc6d29b8 --jq '.files[] | select(.filename=="apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex") | .patch'
+@@ -316,11 +316,12 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.Interface do
+   end
+ 
+   def extract_mappings(
+-        %InterfaceDescriptor{aggregation: :object},
++        %InterfaceDescriptor{aggregation: :object, interface_id: interface_id},
+         _mapping,
+         mappings
+       ) do
+     mappings
++    |> Enum.filter(fn {_id, m} -> m.interface_id == interface_id end)
+     |> Map.new(fn {_id, m} ->
+       key = m.endpoint |> String.split("/") |> List.last()
+       {key, m}
+
+$ gh api repos/astarte-platform/astarte/commits/fdc6d29b8 --jq '.files[] | select(.filename=="apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/interface_test.exs") | .patch'
+@@ -209,6 +209,29 @@ defmodule Astarte.DataUpdaterPlant.DataUpdater.Core.InterfaceTest do
+       end
+     end
+ 
++    @tag regression: true
++    test "extract_mappings/3 only considers mappings for the current object interface", context do
++      %{
++        fixed_object_datastream_1: fixed_object_datastream_1,
++        fixed_object_datastream_2: fixed_object_datastream_2,
++        state: state,
++        realm_name: realm_name
++      } = context
++
++      # ignored for objects
++      mapping = nil
++
++      descriptor = state.interfaces[fixed_object_datastream_1.name]
++
++      assert %{"first" => %{value_type: :string}, "second" => %{value_type: :integer}} =
++               Core.Interface.extract_mappings(descriptor, mapping, state.mappings)
++
++      descriptor = state.interfaces[fixed_object_datastream_2.name]
++
++      assert %{"first" => %{value_type: :boolean}, "second" => %{value_type: :binaryblob}} =
++               Core.Interface.extract_mappings(descriptor, mapping, state.mappings)
++    end
++
+     property "forget_interfaces/2 removes interfaces from state cache", context do
+       %{
+         interfaces: interfaces,
+
+✱ Grep "validateValueType|ValidateValueType|func.*[Vv]alidate.*[Vv]alue" in internal/engine · 1 match
+This is a distinct wire-visible fix: object datastream values were validated against whatever interface had the same endpoint names, since `extract_mappings` for an `:object` aggregation pulled mappings across all interfaces. Let me see where Astrate validates object values so the probe line targets the right files.
+✱ Grep "extract_mappings|extractMappings" in internal · 0 matches
+✱ Grep "[Vv]alue[Tt]ype" in . · 89 matches
+✱ Grep "Lookup|Trie|resolveMapping|Mapping.*by.*path|matchPath" in internal/engine · 13 matches
+Let me find how Astrate resolves the mapping set when decoding an object value — that determines whether the upstream bug has an Astrate mirror.
+→ Read pkg/payload/payload.go
+→ Read internal/engine/data.go [offset=230, limit=120]
+Astrate decodes object values against `ci.ObjectLeaves` — the topic-named interface's own leaves — so the rc.5 #2141 bug may or may not have a mirror; that's exactly what a probe verifies. Let me check rc.4's notes too, since it sits between what prior runs covered (rc.3, rc.5-partial), and confirm nothing else is pending.
+
+$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.4") | .body'
+Fifth Astarte v1.4.0 release candidate.
+
+## Changes since v1.4.0-rc.3
+
+### Added
+
+- [config] Initialize `config` library. All external services declared using the `config` library helpers can be configured either using the `_URL` environment variable or by using single component variables `_SCHEME`, `_HOST`, `_PORT`, `_PATH`, `_QUERY`, `_FRAGMENT` (`_URL` has priority), together with normal `_SSL_*` variables.
+- [secrets] Use `config` for vault configuration. Vault can now be configured using component variables instead of `ASTARTE_VAULT_URL`:
+  - `ASTARTE_VAULT_SCHEME`: scheme part of the URL (binary, default: `"http"`, or `"https"` if `ASTARTE_VAULT_SSL_ENABLED` is `true`)
+  - `ASTARTE_VAULT_HOST`: host part of the URL (binary, default: `"localhost"`)
+  - `ASTARTE_VAULT_PORT`: port part of the url (integer, default `8200`)
+  - `ASTARTE_VAULT_PATH`: path part of the url (string)
+  - `ASTARTE_VAULT_QUERY`: query part of the url (string)
+  - `ASTARTE_VAULT_FRAGMENT`: fragment part of the url (string)
+- [fdo] Use `config` for rendezvous configuration. The rendezvous can now be configured using component variables instead of `PAIRING_FDO_RENDEZVOUS_URL`:
+  - `PAIRING_FDO_RENDEZVOUS_SCHEME`: scheme part of the URL (binary, default: `"http"`, or `"https"` if `PAIRING_FDO_RENDEZVOUS_SSL_ENABLED` is `true`)
+  - `PAIRING_FDO_RENDEZVOUS_HOST`: host part of the URL (binary, default: `"localhost"`)
+  - `PAIRING_FDO_RENDEZVOUS_PORT`: port part of the url (integer, default `8041`)
+  - `PAIRING_FDO_RENDEZVOUS_PATH`: path part of the url (string)
+  - `PAIRING_FDO_RENDEZVOUS_QUERY`: query part of the url (string)
+  - `PAIRING_FDO_RENDEZVOUS_FRAGMENT`: fragment part of the url (string)
+
+### Fixed
+
+- [astarte_data_updater_plant] Ensure RPC server is always available to clients. Resolved the issue where a temporary disconnection and reconnection of data_updater_plant to the cluster would make the RPC server inaccessible.
+
+### Changed
+
+- [astarte_data_updater_plant] Ensure memory is properly garbage collected
+
+← Edit .mule/todo.md
+Index: /root/astrate-mule/.mule/todo.md
+===================================================================
+--- /root/astrate-mule/.mule/todo.md
++++ /root/astrate-mule/.mule/todo.md
+@@ -102,4 +102,6 @@
+ - [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]
+ - [ ] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto]
+ - [ ] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto]
+ - [ ] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto]
++- [ ] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto]
++- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]
+
+Nothing new in rc.4 (config-library and Elixir-infra changes — ignore). The only uncovered item since the last run is rc.5's #2141. Appending the two task lines to the queue.
+$ tail -n 5 .mule/todo.md
+- [ ] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto]

$ git branch --show-current; git log --all --oneline -1 -- .mule/todo.md; git log --oneline -3 -- .mule/todo.md
mule/queue
f4ad882 mule: blocked issue-68-async-operation-accepted
f4ad882 mule: blocked issue-68-async-operation-accepted
e0f8823 mule: blocked issue-91-pairing-health-serve-root
733f965 mule: refill the queue

The last upstream run (2026-08-31) proposed `probe-object-validation-selected-interface` and `compat-note-v1.4.0-rc.5`, but those aren't in the current todo.md on disk. Let me check how the queue state evolved — the file may have been overwritten by refills.
$ git log --all --oneline -S "probe-object-validation-selected-interface" -- .mule/todo.md
(no output)

$ git log --all --oneline --reverse -- .mule/todo.md
956d2a7 Add solo-mule mode: an unsupervised opencode worker driven by a task queue
0921b60 mule: make the Legion Go the concurrency gate
2e885ab mule: note mule/queue as authoritative for todo.md, make race-check recurring
8b333f2 mule: a standing check that finds nothing is passing, not failing
0b5b23f mule: refill its own queue, with GitHub issues as the front door
6f2a33e mule: refill the queue
8d98ab5 mule: issues are the queue, and a review step before main
95fdef5 mule: merge main
71877c3 mule: drop the copied issue line — issues are read live now
8606336 mule: refill the queue
16cc008 mule: log store-realm-cascade-policies
344d013 mule: log store-alias-lowest-id
ad335ec mule: log store-delete-device-objects
8a18476 mule: blocked race-check
4f6d95c mule: split race-check into 5 per-subsystem lines
d145c7d mule: propose upstream v1.3.0 wire-visible changes + compat-note-v1.3.2
28f3adf mule: fix ticks dying outside the repo (cron cwd), recover unsaved queue state
6e30a1c queue: drop three lines already done or duplicated, approve the plaintext-purge line, take the issues-first ordering from main
eeb69d4 mule: blocked control-producer-properties-compression
db8048e mule: blocked probe-interface-default-values
acfcb26 mule: blocked probe-value-type-validation
4c1228c mule: blocked compat-note-v140-rc3
2f072d3 mule: blocked probe-emptycache-resend-device-error
21f74d8 mule: blocked compat-note-v133
25560b4 mule: blocked probe-mqtt-capabilities-declaration
8a983a3 mule: blocked probe-binaryblob-validation
19ad0fe mule: blocked probe-properties-on-connect-encoding
49fcd62 mule: refill the queue
5006f6e mule: blocked docs-sync-pairing-status-enum
733f965 mule: refill the queue
e0f8823 mule: blocked issue-91-pairing-health-serve-root
f4ad882 mule: blocked issue-68-async-operation-accepted

$ git log --all --oneline -S "compat-note-v1.4.0-rc.5"
7b03846 mule: recipe astarte-upstream ran (224s)

$ git show d145c7d --stat | head -20; echo ===; git show d145c7d -- .mule/todo.md | head -40
commit d145c7d98cf19891485a0913928ff02e1db3585a
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Wed Jul 29 00:45:07 2026 +0200

    mule: propose upstream v1.3.0 wire-visible changes + compat-note-v1.3.2

 .mule/todo.md | 5 +++++
 1 file changed, 5 insertions(+)
===
commit d145c7d98cf19891485a0913928ff02e1db3585a
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Wed Jul 29 00:45:07 2026 +0200

    mule: propose upstream v1.3.0 wire-visible changes + compat-note-v1.3.2

diff --git a/.mule/todo.md b/.mule/todo.md
index af5ed5f..32b8fd6 100644
--- a/.mule/todo.md
+++ b/.mule/todo.md
@@ -83,3 +83,8 @@ worth having. Before any of it reaches `main`:
 - [~] bench-big-astarte [legion]: run `bench/scripts/run-tier.sh big astarte` against Legion Go Astarte, commit results (two runs minimum)
 - [~] bench-giant-astarte [legion]: run `bench/scripts/run-tier.sh giant astarte` against Legion Go Astarte, commit results (two runs minimum)
 - [ ] race-check: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...` (~40s). Report any failure to .mule/for-giulio.md with the full race report. This is the only race coverage that exists — the Pi cannot run -race. [legion] [readonly]
+
+- [ ] purge-properties-compression-capability [auto]: upstream v1.3.0 adds a `purge_properties_compression_format` device capability (`zlib`|`plaintext`, default `zlib`) — a wire-visible capability value. Check whether Astrate's capabilities handling (internal/broker, the `<realm>/<device_id>/capabilities` topic, issue #16) needs to recognize/honour it, or whether zlib-only is already the deliberate default; propose the change or note why not needed.
+- [ ] empty-introspection-verification [auto]: upstream v1.3.0 changed "allow devices with empty introspection" — verify whether Astrate's device connection/introspection handling currently rejects an empty introspection string where upstream now accepts it, and propose a fix if so.
+- [ ] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages" — investigate only: does Astrate have an analogous delay between trigger install and first delivery? Report, do not patch.
+- [ ] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.

$ git log --all --format='%h %ad %s' --date=iso -- .mule/todo.md | head -40
f4ad882 2026-09-01 19:46:55 +0200 mule: blocked issue-68-async-operation-accepted
e0f8823 2026-09-01 19:32:34 +0200 mule: blocked issue-91-pairing-health-serve-root
733f965 2026-09-01 19:02:01 +0200 mule: refill the queue
5006f6e 2026-08-31 22:32:39 +0200 mule: blocked docs-sync-pairing-status-enum
49fcd62 2026-08-31 22:10:36 +0200 mule: refill the queue
19ad0fe 2026-08-31 21:45:56 +0200 mule: blocked probe-properties-on-connect-encoding
8a983a3 2026-08-31 21:35:46 +0200 mule: blocked probe-binaryblob-validation
25560b4 2026-08-31 21:20:19 +0200 mule: blocked probe-mqtt-capabilities-declaration
21f74d8 2026-08-31 21:05:59 +0200 mule: blocked compat-note-v133
2f072d3 2026-08-31 20:47:54 +0200 mule: blocked probe-emptycache-resend-device-error
4c1228c 2026-08-31 20:42:21 +0200 mule: blocked compat-note-v140-rc3
acfcb26 2026-08-31 20:26:28 +0200 mule: blocked probe-value-type-validation
db8048e 2026-08-31 19:21:40 +0200 mule: blocked probe-interface-default-values
eeb69d4 2026-08-31 18:52:25 +0200 mule: blocked control-producer-properties-compression
6e30a1c 2026-08-31 02:10:57 +0200 queue: drop three lines already done or duplicated, approve the plaintext-purge line, take the issues-first ordering from main
28f3adf 2026-08-22 02:26:26 +0200 mule: fix ticks dying outside the repo (cron cwd), recover unsaved queue state
d145c7d 2026-07-29 00:45:07 +0200 mule: propose upstream v1.3.0 wire-visible changes + compat-note-v1.3.2
4f6d95c 2026-07-28 23:25:22 +0200 mule: split race-check into 5 per-subsystem lines
8a18476 2026-07-28 00:46:30 +0200 mule: blocked race-check
ad335ec 2026-07-27 17:44:03 +0200 mule: log store-delete-device-objects
344d013 2026-07-27 16:36:55 +0200 mule: log store-alias-lowest-id
16cc008 2026-07-27 15:31:24 +0200 mule: log store-realm-cascade-policies
8606336 2026-07-27 14:58:49 +0200 mule: refill the queue
71877c3 2026-07-27 11:44:22 +0200 mule: drop the copied issue line — issues are read live now
95fdef5 2026-07-27 11:44:22 +0200 mule: merge main
8d98ab5 2026-07-27 11:43:49 +0200 mule: issues are the queue, and a review step before main
6f2a33e 2026-07-27 11:09:20 +0200 mule: refill the queue
0b5b23f 2026-07-27 11:07:06 +0200 mule: refill its own queue, with GitHub issues as the front door
8b333f2 2026-07-27 10:49:27 +0200 mule: a standing check that finds nothing is passing, not failing
2e885ab 2026-07-27 06:12:54 +0200 mule: note mule/queue as authoritative for todo.md, make race-check recurring
0921b60 2026-07-27 06:07:17 +0200 mule: make the Legion Go the concurrency gate
956d2a7 2026-07-27 05:26:54 +0200 Add solo-mule mode: an unsupervised opencode worker driven by a task queue

Let me check what's in `.mule/for-giulio.md` (where compat-note proposals land) and research the v1.4.0-rc.0 "required flag for object-aggregated mappings" change, which is the one wire-visible schema item not yet covered.
→ Read .mule/for-giulio.md

$ gh search commits -R astarte-platform/astarte 'required flag for mappings of object' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo ===; gh search commits -R astarte-platform/astarte 'required flag' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.author.date)\t\(.commit.message|split("\n")[0])"' | head -10
===
82018e1c9	2026-05-20T10:33:00+02:00	feat: add required flag support to interface mappings (#524)
b9a4418d1	2026-05-20T10:33:00+02:00	feat: add required flag support to interface mappings (#524)

Let me examine the upstream `required` flag change (#524) to understand its wire effect before proposing.
$ gh api repos/astarte-platform/astarte/commits/82018e1c9 --jq '{message: .commit.message, files: [.files[] | {filename, status, additions, deletions}]}'
{"files":[{"additions":20,"deletions":2,"filename":"astarte-dashboard/cypress/e2e/interface_builder.cy.js","status":"modified"},{"additions":2,"deletions":1,"filename":"astarte-dashboard/cypress/fixtures/test.astarte.AggregatedObjectInterface.json","status":"modified"},{"additions":6,"deletions":0,"filename":"astarte-dashboard/src/astarte-client/models/Mapping/index.ts","status":"modified"},{"additions":2,"deletions":0,"filename":"astarte-dashboard/src/astarte-client/transforms/mapping.ts","status":"modified"},{"additions":1,"deletions":0,"filename":"astarte-dashboard/src/astarte-client/types/dto/mapping.d.ts","status":"modified"},{"additions":6,"deletions":0,"filename":"astarte-dashboard/src/components/InterfaceEditor.tsx","status":"modified"},{"additions":21,"deletions":0,"filename":"astarte-dashboard/src/components/MappingEditor.tsx","status":"modified"}],"message":"feat: add required flag support to interface mappings (#524)\n\nSigned-off-by: Riccardo Nalgi \u003criccardo.nalgi@secomind.com\u003e"}

$ gh api repos/astarte-platform/astarte/commits/b9a4418d1 --jq '{message: .commit.message, files: [.files[] | {filename, status, additions, deletions}]}'
{"files":[{"additions":20,"deletions":2,"filename":"astarte-dashboard/cypress/e2e/interface_builder.cy.js","status":"modified"},{"additions":2,"deletions":1,"filename":"astarte-dashboard/cypress/fixtures/test.astarte.AggregatedObjectInterface.json","status":"modified"},{"additions":6,"deletions":0,"filename":"astarte-dashboard/src/astarte-client/models/Mapping/index.ts","status":"modified"},{"additions":2,"deletions":0,"filename":"astarte-dashboard/src/astarte-client/transforms/mapping.ts","status":"modified"},{"additions":1,"deletions":0,"filename":"astarte-dashboard/src/astarte-client/types/dto/mapping.d.ts","status":"modified"},{"additions":6,"deletions":0,"filename":"astarte-dashboard/src/components/InterfaceEditor.tsx","status":"modified"},{"additions":21,"deletions":0,"filename":"astarte-dashboard/src/components/MappingEditor.tsx","status":"modified"}],"message":"feat: add required flag support to interface mappings (#524)\n\nSigned-off-by: Riccardo Nalgi \u003criccardo.nalgi@secomind.com\u003e"}

The `required` flag commit is actually a dashboard-UI change (#524), not a backend schema change. Let me check the older recipe-run logs to confirm what v1.4.0-rc.0 items were already considered and rejected, so I don't re-propose dropped items.
$ git show 83568dd --stat | head; echo ===; git show 3a4b693 --stat | head -5; echo ===; git show 3a746d1 --stat | head -5
commit 83568ddfa0a4ebdac26755099fbf9e38ddb0197e
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 17 04:08:56 2026 +0200

    mule: recipe astarte-upstream ran (333s)

 ...817T020856Z-recipe-astarte-upstream-proposed.md | 405 +++++++++++++++++++++
 1 file changed, 405 insertions(+)
===
commit 3a4b69370d6c405411fa7307c7cac3b0b7565812
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 10 04:11:51 2026 +0200

    mule: recipe astarte-upstream ran (431s)
===
commit 3a746d1c64d902f284321c9ca480f02c4ccfe4cb
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 3 04:11:49 2026 +0200

    mule: recipe astarte-upstream ran (485s)

$ git show 83568dd:.mule/reports/log/20260817T020856Z-recipe-astarte-upstream-proposed.md | rg -n "propose|todo|detailed|required|rc\.|v1\.4|Append|Edit" | head -40
2:verdict: proposed
15:v1.4.0-rc.4	2026-08-06T13:40:20Z
16:v1.4.0-rc.3	2026-07-31T10:48:39Z
17:v1.4.0-rc.2	2026-07-14T10:16:22Z
18:v1.4.0-rc.1	2026-07-14T08:53:27Z
58:v1.3.0-rc.2	v1.3.0-rc.2	Third Astarte v1.3.0 release candidate.
60:## Changes since v1.3.0-rc.1
71:v1.3.0-rc.1	v1.3.0-rc.1	Second Astarte v1.3.0 release candidate.
73:# Changes since v1.3.0-rc.0
80:v1.3.0-rc.0	v1.3.0-rc.0	First Astarte v1.3.0 release candidate.
102:## [1.4.0-rc.4] - 2026-08-06
130:## [1.4.0-rc.3] - 2026-07-31
143:## [1.4.0-rc.2] - 2026-07-14
145:## [1.4.0-rc.1] - 2026-07-13
163:## [1.4.0-rc.0] - 2026-04-08
167:- Add required flag for mappings of object aggregated interfaces
176:- [astarte_realm_management] Allow listing interfaces with mappings using the `detailed=true` parameter
221:## [1.3.0-rc.2] - 2026-04-08
232:  - `ASTARTE_BASE_URL_DOMAIN` - domain part of the base URL of astarte, used by devices to connect in TO2 phase (required if FDO enabled)
233:  - `ASTARTE_BASE_URL_PORT` - port of the base URL of astarte (required if FDO enabled)
234:  - `ASTARTE_BASE_URL_PROTOCOL` - protocol of the base URL of astarte (required if FDO enabled)
250:## [1.3.0-rc.1] - 2026-01-26
269:## [1.3.0-rc.0] - 2025-11-21
276:  - `HOUSEKEEPING_ASTARTE_KEYSPACE_NETWORK_REPLICATION_MAP` - Datacenter replication map when using NetworkTopologyStrategy (no default, required when using network strategy)
291:- BREAKING: [astarte_housekeeping] Housekeeping now creates an AMQP vhost for each created realm. Some required configuration was introduced:
292:  - `HOUSEKEEPING_AMQP_HOST` - AMQP host for housekeeping operations (required in production, defaults to "localhost" in dev/test)
297:## [1.3.0-rc.0] - 2025-11-21
304:  - `HOUSEKEEPING_ASTARTE_KEYSPACE_NETWORK_REPLICATION_MAP` - Datacenter replication map when using NetworkTopologyStrategy (no default, required when using network strategy)
319:- BREAKING: [astarte_housekeeping] Housekeeping now creates an AMQP vhost for each created realm. Some required configuration was introduced:
320:  - `HOUSEKEEPING_AMQP_HOST` - AMQP host for housekeeping operations (required in production, defaults to "localhost" in dev/test)
350:- [astarte_realm_management] Allow to delete long-disconnected devices. For this to work, Realm Management needs `CLUSTERING_STRATEGY`, `CLUSTERING_KUBERNETES_NAMESPACE` and `DATA_UPDATER_PLANT_CLUSTERING_KUBERNETES_SELECTOR` to be set, just like AppEngine and DUP. Refer to 1.2.1-rc.0 for additional information on the variables.
355:→ Read .mule/todo.md
360:$ grep -n 'device.deletion.*trigger\|deletion.*trigger\|probe-emptycache\|compat-note-v1.3.3\|compat-note-v1.3' /root/astrate-mule/.mule/todo.md
362:95:- [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]
365:→ Read .mule/todo.md [offset=90]
367:← Edit .mule/todo.md
368:Index: /root/astrate-mule/.mule/todo.md
370:--- /root/astrate-mule/.mule/todo.md
371:+++ /root/astrate-mule/.mule/todo.md
373: - [ ] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto]

$ git show 83568dd:.mule/reports/log/20260817T020856Z-recipe-astarte-upstream-proposed.md | sed -n '150,215p'
### Fixed

- [fdo] Do not crash when using rotated keys as owner keys
- [fdo] Delete devices from unfinished sessions
- [fdo] Delete old sessions from device
- [fdo] Properly respond to service info messages after IsDone=True
- [astarte_realm_management] Properly show default values for all mapping parameters

### Changed

- Enable FDO by default

## [1.4.0-rc.0] - 2026-04-08

### Added

- Add required flag for mappings of object aggregated interfaces
- [astarte_pairing] Support for FIDO Device Onboard v1.1 for device registration. The following environment variables have been added for the Vault connectivity:
  - `ASTARTE_VAULT_URL` - The URL to access Vault.
  - `ASTARTE_VAULT_AUTHENTICATION_MECHANISM` - The mechanism to use for authenticating with Vault (default: `""`, allowed: \["token"])
  - `ASTARTE_VAULT_TOKEN` - Token to authenticate with Vault, used if `ASTARTE_VAULT_AUTHENTICATION_MECHANISM` is `"token"` (default: `""`)
  - `ASTARTE_VAULT_SSL_CA_FILE` - path to the CA certificate file for the vault TLS connection; when not specified, the bundled cURL certificate bundle will be used
  - `ASTARTE_VAULT_SSL_CUSTOM_SNI` - custom SNI hostname for the vault TLS connection; when not specified, the vault hostname will be used
  - `ASTARTE_VAULT_SSL_DISABLE_SNI` - Disable Server Name Indication for Vault. (default: `false`)
  - `ASTARTE_VAULT_SSL_ENABLED` - Whether to enable SSL for the Vault connection. (default `false`)
- [astarte_realm_management] Allow listing interfaces with mappings using the `detailed=true` parameter
- [astarte_import] Added support for data types: `doublearray`, `integerarray`,
  `booleanarray`, `longintegerarray`, `stringarray`, `datetimearray`, `binaryblobarray`.

### Changed

- [astarte_data_updater_plant] Use mississippi consumer for data updater processes

## [1.3.3] - Unreleased

### Fixed

- [astarte_data_updater_plant] Resending properties to a device now reports distinct device error triggers depending on the failure: `interface_loading_failed` when an interface fails to load and `resend_interface_properties_failed` when sending the properties to the device itself fails. Previously, both cases were incorrectly reported as `resend_interface_properties_failed`.

## [1.3.2] - 2026-07-14

## Fixed

- Avoid crash on invalid properties message

## [1.3.1] - 2026-07-13

### Fixed

- Solve capabilities message stalling queues
- [astarte_housekeeping] Allow resuming partially completed realm deletions
- Use network topology as default realm replication
- Fix simple strategy realm creation without replication factor

### Changed

- Immediately delete devices in device deletion best case scenario

## [1.3.0] - 2026-05-06

### Changed

- [astarte_housekeeping] When `HOUSEKEEPING_ASTARTE_KEYSPACE_REPLICATION_STRATEGY` is not set, the `astarte` keyspace is now created using `NetworkTopologyStrategy` with a replication map derived from the current ScyllaDB network topology (one replica per node in each datacenter), instead of falling back to `SimpleStrategy` with replication factor 1.
- [astarte_housekeeping] When creating a realm without specifying its replication strategy, the strategy used in the `astarte` keyspace is applied as default.

$ git show 83568dd:.mule/reports/log/20260817T020856Z-recipe-astarte-upstream-proposed.md | sed -n '355,405p'
→ Read .mule/todo.md
• Check device triggers Astrate Explore Agent
✓ Check device triggers Astrate Explore Agent
Now I have the full picture. Let me check what's already tracked vs. what's new.

$ grep -n 'device.deletion.*trigger\|deletion.*trigger\|probe-emptycache\|compat-note-v1.3.3\|compat-note-v1.3' /root/astrate-mule/.mule/todo.md
94:- [ ] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto]
95:- [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]

Upstream is at v1.3.3, Astrate targets v1.2.x. Most v1.3.x changes are already tracked. The one wire-visible gap not yet in the queue is device deletion trigger events. Let me append the new task lines.
→ Read .mule/todo.md [offset=90]

← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -92,4 +92,5 @@
 - [ ] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto]
 - [ ] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto]
 - [ ] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto]
 - [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]
+- [ ] trigger-device-deletion-events: add `device_deletion_started` and `device_deletion_finished` trigger conditions to `internal/engine/triggers/match.go` (new constants + entries in `deviceOns`), emit them from the device-deletion path in `internal/engine/engine.go`, and add them to the `on` enum in `docs/api/astarte_realm_management_api.yaml` (upstream v1.3.0-rc.0) [auto]

Here's the summary:

**Upstream: v1.3.3 (released 2026-08-07) vs. Astrate target: v1.2.x**

Most v1.3.x changes are already tracked in the queue:
- `control-producer-properties-compression` — plaintext producer/properties (v1.3.0 capability)
- `realm-config-datastream-retention` — per-realm data retention (v1.3.0)
- `probe-emptycache-resend-device-error` — v1.3.3 error differentiation
- `compat-note-v1.3.3` / `compat-note-v1.4.0-rc.3` — compat notes

**New task added:**
- `trigger-device-deletion-events` — upstream v1.3.0-rc.0 added `device_deletion_started` and `device_deletion_finished` trigger conditions. Astrate has no equivalent. This is wire-visible (Realm Management trigger API `on` enum, device lifecycle events).

**Already implemented (no task needed):**
- Empty introspection accepted — already handled in `introspection.go:95-102`
- Capabilities + `purge_properties_compression_format` — fully implemented in `capabilities.go`
- Device registration triggers — already wired in `engine.go:283-301`
- Transient trigger crash fix — N/A (Astrate's transient triggers are matcher-only, never enter pipeline)

**Infrastructure-only / Elixir-only (ignored per recipe rules):**
- AMQP vhost per realm, merged API containers, env var changes, ScyllaDB replication strategy, capabilities queue stalling fix, volatile trigger encoding fix — none of these are wire-visible or applicable to Astrate's architecture.

Done: appended 1 task line (trigger-device-deletion-events) to .mule/todo.md
Files: .mule/todo.md
Verified: grep confirms no duplicate
Unsure: nothing
Follow-ups: trigger-device-deletion-events

Confirmed the picture now. Summary of findings:

- **Upstream movement:** newest tag is still `v1.4.0-rc.5` (2026-08-20); stable is `v1.3.3` (2026-08-07). No new releases since the last upstream run (2026-08-31).
- **The last run's queue append was lost:** commit `7b03846` logged its proposal (`probe-object-validation-selected-interface` for rc.5 #2141, plus `compat-note-v1.4.0-rc.5`) but only the log file landed — the todo.md edit never made it into the branch. Neither line exists in the current queue; the only rc.5 item tracked is the binaryblob fix.
- **Everything else in the v1.4.0 RCs** (config-library, vault, FDO, required-flag dashboard UI #524, housekeeping retention) was reviewed by prior runs as dashboard/infra-only and correctly not proposed.

So I'll restore the two lost rc.5 lines (not duplicates — they're absent from the queue) and nothing else.
← Edit .mule/todo.md
Index: /root/astrate-mule/.mule/todo.md
===================================================================
--- /root/astrate-mule/.mule/todo.md
+++ /root/astrate-mule/.mule/todo.md
@@ -105,4 +105,6 @@
 - [!] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto] — BLOCKED: gates failed
 - [!] docs-sync-pairing-status-enum: in docs/api/astarte_pairing_api.yaml, fix the PairingInfo.status enum (line 367): it lists `confirmed, pending, denied, expired` but the handler only ever emits `confirmed`, `pending`, or `inhibited` (internal/pairing/service.go:289-297, via internal/pairing/http.go:211) — drop the dead `denied`/`expired` and add the undocumented `inhibited`. [auto] — BLOCKED: gates failed
 - [!] issue-91-pairing-health-serve-root: in internal/pairing/http.go, add the unauthenticated `GET /pairing/health` route (no realm segment) sharing the handler/payload of the existing `GET /pairing/v1/{realm}/health`, keeping the v1 route as-is; tests for both paths (closes #91) — BLOCKED: gates failed
 - [!] issue-68-async-operation-accepted: on the mutating endpoints upstream surfaces `async_operation` on (housekeeping realm create/delete; realm-management interface install/update/delete, trigger/policy delete in internal/realm/http.go), accept and ignore `?async_operation=false`, with unparseable/`true` values also not changing behaviour (Astrate stays always-sync); tests cover the flag parsing (closes #68) — BLOCKED: gates failed
+- [ ] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto]
+- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]

$ grep -n "probe-object-validation-selected-interface\|compat-note-v1.4.0-rc.5" .mule/todo.md
109:- [ ] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto]
110:- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]

Done: upstream is at v1.4.0-rc.5 (stable v1.3.3) — no new releases since the last run; restored the two rc.5 proposal lines the previous run appended but lost from the queue (probe-object-validation-selected-interface, compat-note-v1.4.0-rc.5). Nothing new proposed; remaining v1.4.0-rc.x items are config/vault/FDO infra or the dashboard-frontend `required`-flag UI (#524), all correctly skipped.

Files: .mule/todo.md
Verified: `gh api .../releases` + `git log -S` on todo.md (confirm lost append) -> pass; grep confirms no duplicate of the two lines -> pass
Unsure: whether restoring lines whose prior append never committed counts as duplication — no live duplicates exist, so I restored rather than leave the gap
Follow-ups: probe-object-validation-selected-interface, compat-note-v1.4.0-rc.5 (both `[auto]`)
