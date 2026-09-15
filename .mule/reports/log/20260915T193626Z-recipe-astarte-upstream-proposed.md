slug: recipe-astarte-upstream
verdict: proposed
at:  7c7ca2b
ran: 2026-09-15T19:36:26Z on DietPi in 370s


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

→ Read docs/UPSTREAM-EXPERIMENTAL.md
The newest upstream tag is `v1.4.0-rc.5` (2026-08-20) — ahead of what COMPATIBILITY.md pins (1.2.x, with 1.4 experimental). Let me check the experimental doc and the queue to see what's already tracked.
→ Read .mule/todo.md
→ Read .mule/for-giulio.md
The queue already tracks several v1.3.x/v1.4-rc items. Let me check `for-giulio.md` for existing proposals and see whether `rc.5` was already reviewed.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
Sixth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.4

### Fixed

- [astarte_data_updater_plant] Make object values only validate type on selected interface
- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated

The picture: `COMPATIBILITY.md` targets upstream 1.2.2, and the queue already tracks v1.3.x and the v1.4-rc line. Newest overall is `v1.4.0-rc.5` (2026-08-20). Let me check the last astarte-upstream recipe run to avoid re-reviewing what's already been handled.
$ ls .mule/reports/ 2>/dev/null | head -50; echo "---"; git log --oneline --all -30 --grep='upstream' -i
hygiene.md
log
---
203af1c mule: hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto]
015a706 mule: hk-zero-retention-asymmetry [auto]: `PATCH /housekeeping/v1/realms/{realm}` treats `datastream_maximum_storage_retention: 0` identically to `null` — both route to `ClearRetention = true` (internal/housekeeping/http.go:217-218) — but `POST /realms` accepts 0 as a valid non-negative value and persists it (`service.go:146` does not treat 0 as special). Pin the chosen upstream-parity behaviour in a container-free unit test: if 0 is invalid for retention (recommended), add a `*retention == 0` → `ErrValidation` branch in `CreateRealm` (service.go:146) and a `val == 0` → 422 rejection in `patchRealm` (http.go:217); if 0 is valid, drop the `|| val == 0` clause in http.go:217 and add a PATCH test asserting the value round-trips via GET. [auto]
f0863f2 mule: recipe astarte-upstream ran (97s)
6524c7f mule: auth-iat-not-required-test [auto]: pin in internal/auth/jwt_test.go that a token with a *future* `iat` still verifies — the "`iat` is not required — upstream parity" rule (jwt.go:104-105, no `jwt.WithIssuedAt` in the parser, jwt.go:115-118) has zero coverage, so a regression that starts validating `iat` or `WithIssuedAt` would silently reject tokens from a skewed-clock issuer and pass every existing suite. Sign RS256 with `iat` = now+1h and no `exp`/`nbf`, assert `Verify` succeeds.
6bca5cc mule: recipe astarte-upstream ran (100s)
d29e382 mule: recipe astarte-upstream ran (119s)
a609274 mule: recipe astarte-upstream ran (334s)
3cbd3d6 mule: recipe astarte-upstream ran (375s)
ee76ba1 mule: recipe astarte-upstream ran (204s)
321fc79 mule: recipe astarte-upstream ran (152s)
c867b2c mule: docs-sync-pairing-version-endpoint [auto]: add the undocumented `GET /pairing/v1/{realm}/version` route to docs/api/astarte_pairing_api.yaml — it exists in code (cmd/astrate/main.go:398, unauthenticated, upstream-parity per main.go:387-391) and answers 200 `{"data":"<version>"}` via observability.VersionHandler (internal/observability/compat.go:38); mirror the RM spec's getVersion op (astarte_realm_management_api.yaml:532) without its a_rma security. Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
2f6176c mule: docs-sync-pairing-status-enum [auto]: fix the `PairingInfo.status` enum in docs/api/astarte_pairing_api.yaml from `[confirmed, pending, denied, expired]` (yaml:412) to `[confirmed, pending, inhibited]` — the handler never emits `denied` or `expired`: `service.Info` maps only `DeviceStatusInhibited`→"inhibited", `DeviceStatusConfirmed`→"confirmed", and everything else→"pending" (internal/pairing/service.go:296-299, statuses defined in store/devices.go:19-27; upstream-parity per the comment at service.go:282-285). Run `make -C docs build` and confirm the Swagger-referenced YAML still loads.
f463039 mule: recipe astarte-upstream ran (268s)
a701b05 mule: recipe astarte-upstream ran (239s)
32544c0 mule: pairing-bearer-secret-test [auto]: add a container-free unit test for `bearerSecret` (internal/pairing/http.go:343) locking the upstream-verbatim permissive match `~r/bearer\:?\s+(.*)$/i` — the colon form `Bearer: <secret>`, the case variants (`bearer `, `BEARER `), and the reject paths (no whitespace before value, non-bearer scheme, empty value). Only the plain `Bearer <secret>` form is exercised today (integration http_test.go); a regression to case-sensitive or drop-colon matching would pass every existing test. Pure helper in a non-integration file; no Docker.
95ef44e mule: recipe astarte-upstream ran (237s)
db9d06e mule: probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
f117601 mule: recipe astarte-upstream ran (305s)
ab49dcc mule: recipe astarte-upstream ran (127s)
c3cad68 mule: recipe astarte-upstream ran (85s)
4637135 mule: recipe astarte-upstream ran (526s)
4071671 mule: recipe astarte-upstream ran (730s)
4bc3e1a mule: compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
24ad5b8 mule: issue-93: In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged
2a033d4 mule: parking a decision on an upstream release now has something that notices
99743ca mule: the keyAgreement parking condition has fired — record the scope decision
f1d0069 mule: the survey stops stranding reports, and the map stops lying about timers
e0dd32a flow: a lane can no longer be retired underneath an in-flight Submit
6f7a3d6 mule: purge-properties-compression-capability [auto]: upstream v1.3.0 adds a `purge_properties_compression_format` device capability (`zlib`|`plaintext`, default `zlib`) — a wire-visible capability value. Check whether Astrate's capabilities handling (internal/broker, the `<realm>/<device_id>/capabilities` topic, issue #16) needs to recognize/honour it, or whether zlib-only is already the deliberate default; propose the change or note why not needed.
b6e1ca4 compat: make the dashboard health indicators real, pin always-sync (#68, #91)

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select((.tag_name|startswith("v1.4.0-rc")) and (.tag_name != "v1.4.0-rc.5")) | "===== \(.tag_name) \(.published_at) =====\n\(.body)\n"' | head -120
===== v1.4.0-rc.4 2026-08-06T13:40:20Z =====
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

===== v1.4.0-rc.3 2026-07-31T10:48:39Z =====
Fourth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.2

### Added

- [secrets] Configurable base namespace for Vault using `ASTARTE_VAULT_BASE_NAMESPACE`: all created namespaces will be children of the base namespace.

### Fixed

- [fdo] Support for Vault and Vault managed
- [fdo] Properly handle chunked service info messages
- [astarte_data_updater_plant] Avoid clogging AMQP queues when receiving invalid device messages
- [astarte_data_updater_plant] Restore correct type validation checks for incoming data

===== v1.4.0-rc.2 2026-07-14T10:16:22Z =====
Third Astarte v1.4.0 release candidate.

===== v1.4.0-rc.1 2026-07-14T08:53:27Z =====
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

===== v1.4.0-rc.0 2026-04-08T12:32:33Z =====
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

$ git log --all --oneline --format='%h %ad %s' --date=short --grep='recipes/astarte-upstream\|recipe astarte-upstream' | head -15
f0863f2 2026-09-15 mule: recipe astarte-upstream ran (97s)
6bca5cc 2026-09-14 mule: recipe astarte-upstream ran (100s)
d29e382 2026-09-13 mule: recipe astarte-upstream ran (119s)
a609274 2026-09-13 mule: recipe astarte-upstream ran (334s)
3cbd3d6 2026-09-12 mule: recipe astarte-upstream ran (375s)
ee76ba1 2026-09-11 mule: recipe astarte-upstream ran (204s)
321fc79 2026-09-10 mule: recipe astarte-upstream ran (152s)
f463039 2026-09-10 mule: recipe astarte-upstream ran (268s)
a701b05 2026-09-09 mule: recipe astarte-upstream ran (239s)
95ef44e 2026-09-08 mule: recipe astarte-upstream ran (237s)
f117601 2026-09-08 mule: recipe astarte-upstream ran (305s)
ab49dcc 2026-09-07 mule: recipe astarte-upstream ran (127s)
c3cad68 2026-09-07 mule: recipe astarte-upstream ran (85s)
4637135 2026-09-06 mule: recipe astarte-upstream ran (526s)
4071671 2026-09-05 mule: recipe astarte-upstream ran (730s)

Let me see what the recent astarte-upstream runs concluded, to know whether rc.4/rc.5 were already analysed.
$ git show f0863f2 --stat; echo ====; git show 6bca5cc --stat
commit f0863f2cee0005d3d152aafe4351df8296ab3d77
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Tue Sep 15 13:13:54 2026 +0200

    mule: recipe astarte-upstream ran (97s)

 ...915T111354Z-recipe-astarte-upstream-proposed.md | 34 ++++++++++++++++++++++
 1 file changed, 34 insertions(+)
====
commit 6bca5cc2abd2c1cc68ba77419216af2df3ca110b
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Sep 14 04:06:08 2026 +0200

    mule: recipe astarte-upstream ran (100s)

 ...0260914T020552Z-recipe-astarte-upstream-noop.md | 43 +++++++++++
 ...914T020608Z-recipe-astarte-upstream-proposed.md | 83 ++++++++++++++++++++++
 2 files changed, 126 insertions(+)

The runs write report files. Let me read the most recent astarte-upstream report (2026-09-15, the current day's earlier run) to see what was proposed.
$ git show f0863f2 --name-only --format='' ; echo ---; git show f0863f2 --format='' | head -60
.mule/reports/log/20260915T111354Z-recipe-astarte-upstream-proposed.md
---
diff --git a/.mule/reports/log/20260915T111354Z-recipe-astarte-upstream-proposed.md b/.mule/reports/log/20260915T111354Z-recipe-astarte-upstream-proposed.md
new file mode 100644
index 0000000..45e94b7
--- /dev/null
+++ b/.mule/reports/log/20260915T111354Z-recipe-astarte-upstream-proposed.md
@@ -0,0 +1,34 @@
+slug: recipe-astarte-upstream
+verdict: proposed
+at:  b36f50a
+ran: 2026-09-15T11:13:54Z on DietPi in 97s
+
+
+> build · big-pickle
+
+→ Read .mule/MULE.md
+→ Read .mule/recipes/astarte-upstream.md
+→ Read docs/COMPATIBILITY.md
+→ Read .mule/todo.md
+Recipe read. It's a proposal job — output goes to `.mule/todo.md` + any evidence files. No issue creation in this recipe, no git. Starting: check upstream releases vs compatibility target.
+
+$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
+v1.4.0-rc.5	2026-08-20T12:28:04Z
+v1.3.3	2026-08-07T14:33:58Z
+v1.4.0-rc.4	2026-08-06T13:40:20Z
+v1.4.0-rc.3	2026-07-31T10:48:39Z
+v1.4.0-rc.2	2026-07-14T10:16:22Z
+
+Newest tag: v1.4.0-rc.5 (2026-08-20); newest stable: v1.3.3 (2026-08-07). The prior run already saw both (todo.md's `compat-note-v1.3.3` line says "v1.4.0 is still rc.5-only"). Let me confirm the compat-note state in for-giulio.md before concluding.
+→ Read .mule/for-giulio.md
+Upstream has not moved. Newest tag is still **v1.4.0-rc.5** (2026-08-20); newest stable is still **v1.3.3** (2026-08-07). Both are exactly what the previous runs already saw: the `compat-note-v1.3.3` line and `.mule/for-giulio.md` (milestone runs 2026-09-09→09-13) record the same pair, and the v1.3.2/v1.3.3/v1.4-rc compat notes and the v1.3.0 probes are already proposed in the queue.
+
+Per the recipe, an empty result is a good result: I append nothing to `.mule/todo.md` and write no evidence file.
+
+```
+Done: no upstream movement since v1.4.0-rc.5 — newest tag unchanged (v1.4.0-rc.5, 2026-08-20; stable still v1.3.3, 2026-08-07); nothing new to propose
+Files: none
+Verified: gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"' -> pass (5 tags, newest v1.4.0-rc.5); compared against docs/COMPATIBILITY.md target (1.2.2, v1.3 surface delivered, v1.4 experimental) and prior queue state
+Unsure: nothing — release set matches the last recorded run exactly
+Follow-ups: none
+```

The recipe's comparison point is `docs/COMPATIBILITY.md` (target 1.2.2), and upstream is ahead at `v1.4.0-rc.5`. Recent runs concluded "no movement since rc.5" — but I need to check whether any of them actually read rc.5's release body (its two DUP validation fixes touch areas Astrate implements). Let me check the report history.
$ ls -la .mule/reports/log/ 2>/dev/null | tail -40
-rw-r--r-- 1 root root  10121 Sep 12 22:34 20260912T203456Z-docs-sync-ae-post-group-devices-409-done.md
-rw-r--r-- 1 root root   8676 Sep 12 22:46 20260912T204618Z-docs-sync-ae-patch-merge-patch-media-type-done.md
-rw-r--r-- 1 root root  27354 Sep 13 13:15 20260913T111551Z-recipe-hygiene-proposed.md
-rw-r--r-- 1 root root   2381 Sep 13 13:22 20260913T112223Z-store-todo-lttb-stale-done.md
-rw-r--r-- 1 root root    212 Sep 13 18:52 20260913T165246Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  68245 Sep 13 19:05 20260913T170550Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   1971 Sep 13 19:13 20260913T171327Z-recipe-astarte-upstream-noop.md
-rw-r--r-- 1 root root  25213 Sep 13 19:13 20260913T171351Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  20763 Sep 13 19:27 20260913T172714Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root   5810 Sep 13 19:57 20260913T175723Z-config-fail-loud-engine-shards-done.md
-rw-r--r-- 1 root root   7368 Sep 13 20:03 20260913T180348Z-config-fail-loud-dev-mode-done.md
-rw-r--r-- 1 root root   8077 Sep 13 20:09 20260913T180935Z-config-key-resolvers-test-done.md
-rw-r--r-- 1 root root  10426 Sep 13 20:23 20260913T182321Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root   7448 Sep 13 20:50 20260913T185042Z-docs-sync-hk-delete-gating-responses-done.md
-rw-r--r-- 1 root root   3469 Sep 13 21:19 20260913T191910Z-docs-sync-hk-validation-example-done.md
-rw-r--r-- 1 root root   8799 Sep 13 22:23 20260913T202313Z-recipe-hygiene-timeout.md
-rw-r--r-- 1 root root  45997 Sep 13 22:32 20260913T203209Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root   2991 Sep 13 22:38 20260913T203800Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  54978 Sep 13 22:38 20260913T203830Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   3313 Sep 13 22:44 20260913T204415Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   2038 Sep 14 04:06 20260914T020552Z-recipe-astarte-upstream-noop.md
-rw-r--r-- 1 root root   3697 Sep 14 04:06 20260914T020608Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   9178 Sep 14 13:35 20260914T113508Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root   6128 Sep 14 13:44 20260914T114442Z-auth-iat-not-required-test-done.md
-rw-r--r-- 1 root root   6482 Sep 14 18:54 20260914T165413Z-auth-empty-claims-403-test-done.md
-rw-r--r-- 1 root root   4029 Sep 14 19:15 20260914T171558Z-auth-cache-default-size-test-done.md
-rw-r--r-- 1 root root  23102 Sep 14 19:45 20260914T174510Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root   7060 Sep 14 20:18 20260914T181827Z-recipe-hygiene-timeout.md
-rw-r--r-- 1 root root  25311 Sep 15 01:45 20260914T234550Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  26800 Sep 15 12:57 20260915T105740Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root   2011 Sep 15 13:13 20260915T111354Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  18283 Sep 15 13:30 20260915T113026Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root 151656 Sep 15 19:07 20260915T170704Z-hk-zero-retention-asymmetry-done.md
-rw-r--r-- 1 root root    256 Sep 15 19:13 20260915T171349Z-hk-zero-reglimit-create-asymmetry-blocked.md
-rw-r--r-- 1 root root  83771 Sep 15 19:37 20260915T173739Z-hk-zero-retention-asymmetry-done.md
-rw-r--r-- 1 root root  36601 Sep 15 20:07 20260915T180720Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root   4185 Sep 15 20:17 20260915T181748Z-docs-sync-native-socket-query-token-auth-done.md
-rw-r--r-- 1 root root   4410 Sep 15 20:41 20260915T184105Z-recipe-hygiene-timeout.md
-rw-r--r-- 1 root root  14548 Sep 15 20:44 20260915T184406Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  78833 Sep 15 20:58 20260915T185804Z-recipe-github-issues-proposed.md

$ rg -l "rc\.5|rc\.4|rc\.3|rc\.2" .mule/reports/log/ | head; echo ===; rg -il "object values|binaryblob|purge_properties|required mapping|object aggregated" .mule/reports/log/ | head
.mule/reports/log/20260913T171327Z-recipe-astarte-upstream-noop.md
.mule/reports/log/20260908T164748Z-compat-note-v14-rc-blocked.md
.mule/reports/log/20260907T112447Z-docs-sync-appengine-query-params-status-done.md
.mule/reports/log/20260908T193127Z-recipe-milestones-proposed.md
.mule/reports/log/20260913T203830Z-recipe-github-issues-proposed.md
.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260910T202949Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260914T020552Z-recipe-astarte-upstream-noop.md
.mule/reports/log/20260913T171351Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260905T181053Z-recipe-milestones-proposed.md
===
.mule/reports/log/20260913T171327Z-recipe-astarte-upstream-noop.md
.mule/reports/log/20260907T112447Z-docs-sync-appengine-query-params-status-done.md
.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260913T171351Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260915T173739Z-hk-zero-retention-asymmetry-done.md
.mule/reports/log/20260904T190911Z-purge-properties-compression-capabilityauto-done.md
.mule/reports/log/20260908T164748Z-compat-note-v14-rc-blocked.md
.mule/reports/log/20260910T170149Z-recipe-code-review-proposed.md
.mule/reports/log/20260909T230411Z-recipe-milestones-proposed.md

The prior run today concluded "no movement since rc.5" after checking tags — but per the recipe, the comparison is against COMPATIBILITY.md's target (1.2.2), and upstream is ahead. I need to check whether the rc.5 release *notes* were ever actually triaged.
$ git log --all --format='%h %ad %s' --date=iso --grep='astarte-upstream' | rg '2026-08-2|2026-08-3|2026-09-0' ; echo ===; ls .mule/reports/log/ | rg 'astarte-upstream'
a701b05 2026-09-09 20:43:15 +0200 mule: recipe astarte-upstream ran (239s)
95ef44e 2026-09-08 21:38:14 +0200 mule: recipe astarte-upstream ran (237s)
f117601 2026-09-08 13:34:22 +0200 mule: recipe astarte-upstream ran (305s)
ab49dcc 2026-09-07 19:55:24 +0200 mule: recipe astarte-upstream ran (127s)
c3cad68 2026-09-07 04:03:59 +0200 mule: recipe astarte-upstream ran (85s)
4637135 2026-09-06 20:31:02 +0200 mule: recipe astarte-upstream ran (526s)
4071671 2026-09-05 20:46:25 +0200 mule: recipe astarte-upstream ran (730s)
2a033d4 2026-09-04 21:52:33 +0200 mule: parking a decision on an upstream release now has something that notices
ca8e88e 2026-09-04 20:13:24 +0200 mule: recipe astarte-upstream ran (72s)
0468c86 2026-09-03 21:54:22 +0200 mule: recipe astarte-upstream ran (129s)
3debf0d 2026-09-03 13:06:58 +0200 mule: recipe astarte-upstream ran (105s)
2e434c1 2026-09-02 13:22:57 +0200 mule: recipe astarte-upstream ran (104s)
709bc6a 2026-09-01 19:57:16 +0200 mule: recipe astarte-upstream ran (478s)
7b03846 2026-08-31 04:12:36 +0200 mule: recipe astarte-upstream ran (224s)
071d0a2 2026-08-24 04:14:45 +0200 mule: recipe astarte-upstream ran (396s)
===
20260905T184625Z-recipe-astarte-upstream-proposed.md
20260906T183102Z-recipe-astarte-upstream-proposed.md
20260907T020359Z-recipe-astarte-upstream-proposed.md
20260907T175524Z-recipe-astarte-upstream-proposed.md
20260908T113422Z-recipe-astarte-upstream-proposed.md
20260908T193614Z-recipe-astarte-upstream-noop.md
20260908T193814Z-recipe-astarte-upstream-proposed.md
20260909T184315Z-recipe-astarte-upstream-proposed.md
20260910T112515Z-recipe-astarte-upstream-proposed.md
20260910T112545Z-recipe-astarte-upstream-proposed.md
20260910T202936Z-recipe-astarte-upstream-noop.md
20260910T202949Z-recipe-astarte-upstream-proposed.md
20260911T191541Z-recipe-astarte-upstream-proposed.md
20260912T180735Z-recipe-astarte-upstream-proposed.md
20260913T171327Z-recipe-astarte-upstream-noop.md
20260913T171351Z-recipe-astarte-upstream-proposed.md
20260913T204415Z-recipe-astarte-upstream-proposed.md
20260914T020552Z-recipe-astarte-upstream-noop.md
20260914T020608Z-recipe-astarte-upstream-proposed.md
20260915T111354Z-recipe-astarte-upstream-proposed.md

Many runs exist. Let me check whether any run actually triaged the rc.5 release-note bullets against Astrate, by reading the reports from around when rc.5 first appeared.
$ rg -l "astarte_data_updater_plant|object values|Ensure binaryblob|selected interface|config library|Mississippi" .mule/reports/log/; echo "==="; git log --format='%h %ad %s' --date=iso | rg 'astarte-upstream' | rg '08-2[0-4]|08-3'
.mule/reports/log/20260908T164748Z-compat-note-v14-rc-blocked.md
.mule/reports/log/20260913T171327Z-recipe-astarte-upstream-noop.md
.mule/reports/log/20260908T193814Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260913T171351Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260907T175524Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260911T185303Z-recipe-milestones-proposed.md
.mule/reports/log/20260910T202949Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260908T113422Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260908T114256Z-probe-required-mapping-flag-done.md
.mule/reports/log/20260911T110803Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260910T184706Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260911T191541Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260913T203830Z-recipe-github-issues-proposed.md
.mule/reports/log/20260905T181053Z-recipe-milestones-proposed.md
.mule/reports/log/20260910T200120Z-recipe-github-issues-proposed.md
.mule/reports/log/20260915T173739Z-hk-zero-retention-asymmetry-done.md
.mule/reports/log/20260906T183102Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260913T170550Z-recipe-github-issues-proposed.md
.mule/reports/log/20260914T234550Z-recipe-milestones-proposed.md
.mule/reports/log/20260911T194026Z-recipe-code-review-proposed.md
.mule/reports/log/20260904T195529Z-issue-93-done.md
.mule/reports/log/20260915T105740Z-recipe-github-issues-proposed.md
.mule/reports/log/20260913T203209Z-recipe-milestones-proposed.md
.mule/reports/log/20260912T180735Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260909T110758Z-recipe-code-review-proposed.md
.mule/reports/log/20260908T170007Z-recipe-code-review-proposed.md
.mule/reports/log/20260909T184315Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260905T184625Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260908T193321Z-recipe-github-issues-proposed.md
.mule/reports/log/20260910T112545Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260906T183952Z-probe-props-resend-error-triggers-blocked.md
===

$ git show 071d0a2 --format='%h %ad %s' --stat | head -20; echo ===; git show 7b03846 --format='%h %ad %s' --stat | head -20
071d0a2 Mon Aug 24 04:14:45 2026 +0200 mule: recipe astarte-upstream ran (396s)

 ...824T021445Z-recipe-astarte-upstream-proposed.md | 301 +++++++++++++++++++++
 1 file changed, 301 insertions(+)
===
7b03846 Mon Aug 31 04:12:36 2026 +0200 mule: recipe astarte-upstream ran (224s)

 ...831T021236Z-recipe-astarte-upstream-proposed.md | 280 +++++++++++++++++++++
 1 file changed, 280 insertions(+)

$ git show 071d0a2 --format='' | rg -n "rc\.4|rc\.5|binaryblob|object values|selected interface|secret|Announce|config|FDO|Vault" | head -40
21:+v1.4.0-rc.5	2026-08-20T12:28:04Z
23:+v1.4.0-rc.4	2026-08-06T13:40:20Z
28:+v1.4.0-rc.5	2026-08-20T12:28:04Z
30:+v1.4.0-rc.4	2026-08-06T13:40:20Z
49:+Upstream is ahead: `docs/COMPATIBILITY.md` targets upstream **1.2.x** (`APICompatVersion` = `1.2.2`, measurements recorded against `v1.2.0`, Dashboard `v1.2.2`), while upstream's newest tags are `v1.4.0-rc.5` (2026-08-20) and stable `v1.3.3` (2026-08-07). Per the recipe, I read the release notes before proposing anything.
52:+$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
55:+## Changes since v1.4.0-rc.4
59:+- [astarte_data_updater_plant] Make object values only validate type on selected interface
60:+- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated
62:+`v1.3.3`'s body is empty; `v1.4.0-rc.5` only lists two data-updater validation fixes since rc.4. Astrate targets `v1.2.2` (2026-04-27), so the real gap is the whole **v1.3 line** (v1.3.0, May 2026) plus the v1.4.0 RC series. Fetching those notes.
75:+- [astarte_pairing] Add realm-scoped health checks, which can be used to test astarte health or connectivity by devices who are only aware of realm-scoped paths, such as during FDO
80:+- Added database events handling configuration across all services:
93:+- [astarte_pairing] FDO authentication (EXPERIMENTAL feature, disabled by default). New environment variables are needed in order to use FDO:
94:+  - `PAIRING_ENABLE_FDO` - whether the FDO feature is enabled or not (default: false)
95:+  - `PAIRING_FDO_RENDEZVOUS_URL` - URL of the rendezvous server (default: "http://rendezvous:8041")
96:+  - `ASTARTE_BASE_URL_DOMAIN` - domain part of the base URL of astarte, used by devices to connect in TO2 phase (required if FDO enabled)
97:+  - `ASTARTE_BASE_URL_PORT` - port of the base URL of astarte (required if FDO enabled)
98:+  - `ASTARTE_BASE_URL_PROTOCOL` - protocol of the base URL of astarte (required if FDO enabled)
99:+  - `PAIRING_FDO_RENDEZVOUS_SSL_ENABLED` - whether SSL is enabled for the connection to the rendezvous server (default: false)
100:+  - `PAIRING_FDO_RENDEZVOUS_SSL_CA_FILE` - path to the CA certificate file for the rendezvous server TLS connection; when not specified, the bundled cURL certificate bundle will be used
101:+  - `PAIRING_FDO_RENDEZVOUS_SSL_DISABLE_SNI` - disable Server Name Indication for the rendezvous server TLS connection (default: false)
102:+  - `PAIRING_FDO_RENDEZVOUS_SSL_CUSTOM_SNI` - custom SNI hostname for the rendezvous server TLS connection; when not specified, the Rendezvous hostname will be used
103:+- [astarte_housekeeping] AMQP management configuration moved to `astarte_events` library. Environment variables changed:
117:+- BREAKING: [astarte_housekeeping] Housekeeping now creates an AMQP vhost for each created realm. Some required configuration was introduced:
126:+- BREAKING: [astarte_data_updater_plant] Added a separate the AMQP producer configuration, which is mandatory:
140:+- BREAKING: [astarte_pairing] AMQP Producer configuration is now mandatory using the `ASTARTE_EVENTS_PRODUCER_AMQP_*` environment variables
141:+- BREAKING: [astarte_realm_management] AMQP Producer configuration is now mandatory using the `ASTARTE_EVENTS_PRODUCER_AMQP_*` environment variables
147:+- [astarte_housekeeping] Changed logger config to make the metadata fields appear in logfmt log output; in particular the msg field
148:+- [astarte_pairing] Changed logger config to make the metadata fields appear in logfmt log output; in particular the msg field
174:+- Enable FDO by default
176:+The gap is exactly **v1.3.0's "Changes since 1.2.2"** plus the v1.4.0 RCs. Filtering strictly per the recipe (wire-visible / schema rules qualify; AMQP vhosts, clustering env vars, Cassandra replication, container merges, FDO plumbing are infra/Elixir-shaped and ignored), the qualifying items are: empty-introspection acceptance, the new MQTT v1 **capabilities** mechanism (`purge_properties_compression_format`), the v1.4.0-rc.5 data-validation fixes, and the v1.3.0 properties-on-connect encoding fix. Before writing task lines, three narrow greps so each line names a real file.
187:+$ rg -n -i "object" internal/engine/*.go --type go | grep -i -E "aggregat|valid" | head -6; rg -n "binaryblob|BinaryBlob|binary_blob" --type go -l | head -5
217:+Remaining uncovered items from the notes: binaryblob validation (rc.5), property encoding on connect (v1.3.0 Fixed), realm-scoped pairing health checks (v1.3.0 Added), device **registration** triggers (deletion ones are queued, registration is not), and the capability-declaration wire syntax itself. One cheap grep so the trigger line names real symbols:
236:+`OnDeviceRegistered` (`device_registered`) already exists at `internal/engine/triggers/match.go:44`, so upstream's new registration triggers are already matched by Astrate — dropped. Checking the remaining un-read RC notes (rc.0, rc.3, rc.4) so I don't miss anything before appending:
237:+$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.3" or .tag_name=="v1.4.0-rc.4") | .body'
244:+- [config] Initialize `config` library. All external services declared using the `config` library helpers can be configured either using the `_URL` environment variable or by using single component variables `_SCHEME`, `_HOST`, `_PORT`, `_PATH`, `_QUERY`, `_FRAGMENT` (`_URL` has priority), together with normal `_SSL_*` variables.
245:+- [secrets] Use `config` for vault configuration. Vault can now be configured using component variables instead of `ASTARTE_VAULT_URL`:
252:+- [fdo] Use `config` for rendezvous configuration. The rendezvous can now be configured using component variables instead of `PAIRING_FDO_RENDEZVOUS_URL`:
253:+  - `PAIRING_FDO_RENDEZVOUS_SCHEME`: scheme part of the URL (binary, default: `"http"`, or `"https"` if `PAIRING_FDO_RENDEZVOUS_SSL_ENABLED` is `true`)
254:+  - `PAIRING_FDO_RENDEZVOUS_HOST`: host part of the URL (binary, default: `"localhost"`)

Let me read the first reports after rc.5 shipped (2026-08-20) to see how the rc.5/rc.4 notes were triaged.
$ git show 7b03846 --format='' | rg -n "rc\.4|rc\.5|binaryblob|object values|selected interface|secret|Announce|FDO|Vault" | head -40
22:+v1.4.0-rc.5	2026-08-20T12:28:04Z
24:+v1.4.0-rc.4	2026-08-06T13:40:20Z
29:+Upstream is ahead: newest tag is `v1.4.0-rc.5` (stable also has `v1.3.3`), while COMPATIBILITY.md records upstream v1.2.0/v1.2.2. Let me check the todo queue for prior runs and pull the release notes.
31:+$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body'
34:+## Changes since v1.4.0-rc.4
38:+- [astarte_data_updater_plant] Make object values only validate type on selected interface
39:+- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated
42:+The queue already covers v1.3.0/v1.3.3/v1.4.0-rc.3 and rc.5's binaryblob fix, but not rc.5's other change: "Make object values only validate type on selected interface." Let me look at that upstream change before proposing anything.
47:+$ gh search commits --owner astarte-platform 'Make object values only validate type on selected interface' --json sha,repository,commit --jq '.[] | "\(.repository.full_name)\t\(.sha[0:9])\t\(.commit.message | split("\n")[0])"'
50:+$ gh search commits --owner astarte-platform 'object values only validate type on selected interface' --json sha,repository,commit 2>/dev/null; echo "---"; gh search commits --owner astarte-platform 'object values' --json sha,repository,commit.url,commit.author --jq '.[] | "\(.repository.full_name)\t\(.sha[0:9])\t\(.commit.message | split("\n")[0])"' 2>&1 | head -20
81:+The data updater code now lives in the astarte monorepo. Let me find the actual commit for the rc.5 object-validation change.
88:+$ gh search commits --owner astarte-platform --repository astarte 'binaryblob data is correctly validated' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo "==="; gh search commits --owner astarte-platform --repository astarte 'object values only validate type' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'
146:+$ gh search commits -R astarte-platform/astarte 'binaryblob data is correctly validated' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo "==="; gh search commits -R astarte-platform/astarte 'object values only validate type on selected interface' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'
147:+aedc41483	fix(dup): ensure binaryblob data is correctly validated
150:+$ gh search commits -R astarte-platform/astarte 'object values only validate type' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"'; echo "===loose==="; gh search commits -R astarte-platform/astarte 'object values' --json sha,commit --jq '.[] | "\(.sha[0:9])\t\(.commit.message|split("\n")[0])"' | head -10
152:+cc89858b3	DUP: do not crash on binaryblob insertion
157:+$ gh api repos/astarte-platform/astarte/compare/v1.4.0-rc.4...v1.4.0-rc.5 --jq '.commits[] | "\(.sha[0:9])\t\(.commit.author.date)\t\(.commit.message | split("\n")[0])"'
160:+aedc41483	2026-08-07T07:21:33Z	fix(dup): ensure binaryblob data is correctly validated
162:+fdc6d29b8	2026-08-19T08:14:18Z	fix(dup): Make object values only validate type on selected interface (#2141)
164:+3484096d2	2026-08-20T12:10:25Z	Prepare release v1.4.0-rc.5 (#2147)
168:+{"files":[{"additions":1,"deletions":1,"filename":"apps/astarte_appengine_api/mix.lock","status":"modified"},{"additions":2,"deletions":1,"filename":"apps/astarte_data_updater_plant/lib/astarte_data_updater_plant/data_updater/core/interface.ex","status":"modified"},{"additions":40,"deletions":0,"filename":"apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/data_handler_test.exs","status":"modified"},{"additions":23,"deletions":0,"filename":"apps/astarte_data_updater_plant/test/astarte_data_updater_plant/data_updater/core/interface_test.exs","status":"modified"},{"additions":58,"deletions":0,"filename":"apps/astarte_data_updater_plant/test/support/cases/device.ex","status":"modified"},{"additions":1,"deletions":1,"filename":"apps/astarte_pairing/mix.lock","status":"modified"}],"message":"fix(dup): Make object values only validate type on selected interface (#2141)\n\nSigned-off-by: Francesco Noacco \u003cfrancesco.noacco@secomind.com\u003e","sha":"fdc6d29b84ec06a12316340b60f1b43491ca3d5c"}
210:++      assert %{"first" => %{value_type: :boolean}, "second" => %{value_type: :binaryblob}} =
219:+This is a distinct wire-visible fix: object datastream values were validated against whatever interface had the same endpoint names, since `extract_mappings` for an `:object` aggregation pulled mappings across all interfaces. Let me see where Astrate validates object values so the probe line targets the right files.
226:+Astrate decodes object values against `ci.ObjectLeaves` — the topic-named interface's own leaves — so the rc.5 #2141 bug may or may not have a mirror; that's exactly what a probe verifies. Let me check rc.4's notes too, since it sits between what prior runs covered (rc.3, rc.5-partial), and confirm nothing else is pending.
228:+$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.4") | .body'
236:+- [secrets] Use `config` for vault configuration. Vault can now be configured using component variables instead of `ASTARTE_VAULT_URL`:
243:+- [fdo] Use `config` for rendezvous configuration. The rendezvous can now be configured using component variables instead of `PAIRING_FDO_RENDEZVOUS_URL`:
244:+  - `PAIRING_FDO_RENDEZVOUS_SCHEME`: scheme part of the URL (binary, default: `"http"`, or `"https"` if `PAIRING_FDO_RENDEZVOUS_SSL_ENABLED` is `true`)
245:+  - `PAIRING_FDO_RENDEZVOUS_HOST`: host part of the URL (binary, default: `"localhost"`)
246:+  - `PAIRING_FDO_RENDEZVOUS_PORT`: port part of the url (integer, default `8041`)
247:+  - `PAIRING_FDO_RENDEZVOUS_PATH`: path part of the url (string)
248:+  - `PAIRING_FDO_RENDEZVOUS_QUERY`: query part of the url (string)
249:+  - `PAIRING_FDO_RENDEZVOUS_FRAGMENT`: fragment part of the url (string)
267:+ - [ ] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto]
269:++- [ ] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto]
270:++- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]
272:+Nothing new in rc.4 (config-library and Elixir-infra changes — ignore). The only uncovered item since the last run is rc.5's #2141. Appending the two task lines to the queue.
275:+- [ ] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto]
277:+- [ ] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto]
278:+- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]

$ rg -n "probe-binaryblob|probe-object-validation|compat-note-v1.4.0-rc.5|OBJECT|binaryblob" .mule/todo.md .mule/for-giulio.md
.mule/todo.md:112:- [!] probe-property-resend-encoding [auto]: upstream v1.3.0 fixed outbound server-property values sent to a device on connect/emptyCache (commit 522ccf4f — a raw stored binaryblob is now wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding, so it ships as a BSON binary, not a string). Investigate only: does Astrate's resendServerProperties/rehydrate path (internal/engine/control.go:144, pkg/payload) emit binaryblob → BSON subtype-0 binary and datetime → UTC datetime exactly like upstream 1.3? Report, do not patch. — BLOCKED: wrote nothing

The 08-31 run already triaged rc.5 fully (binaryblob validation, object-validation-selected-interface, compat-note-v1.4.0-rc.5). Let me confirm their current state in the queue.
The rc.5 probes proposed on 08-31 don't appear in the current queue — let me check whether they were ever added or were lost in the queue rebuild.
$ git log --all --oneline -S 'probe-binaryblob-validation' -- .mule/todo.md; echo ===; git log --all --oneline -S 'probe-object-validation-selected-interface' -- .mule/todo.md; echo ===; git log --all --oneline -S 'compat-note-v1.4.0-rc.5' -- .mule/todo.md
6e30a1c queue: drop three lines already done or duplicated, approve the plaintext-purge line, take the issues-first ordering from main
===
4cd1e88 mule: refill the queue
===
4cd1e88 mule: refill the queue

$ git show 6e30a1c --stat; echo ===; git show 4cd1e88 --stat | head -20
commit 6e30a1c319d2cc5ec9e054052aebb437335c53b9
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 31 02:10:57 2026 +0200

    queue: drop three lines already done or duplicated, approve the plaintext-purge line, take the issues-first ordering from main

 .mule/todo.md | 11 ++++++++---
 tools/mule.sh | 11 +++++++----
 2 files changed, 15 insertions(+), 7 deletions(-)
===
commit 4cd1e88cd992d8f34e721ac3da8fb662476abc3a
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Tue Sep 1 19:57:16 2026 +0200

    mule: refill the queue

 .mule/cron.log | 918 +++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 .mule/todo.md  |   2 +
 2 files changed, 920 insertions(+)

$ git show 6e30a1c -- .mule/todo.md | head -60
commit 6e30a1c319d2cc5ec9e054052aebb437335c53b9
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 31 02:10:57 2026 +0200

    queue: drop three lines already done or duplicated, approve the plaintext-purge line, take the issues-first ordering from main

diff --git a/.mule/todo.md b/.mule/todo.md
index 0b3fe2f..76bd7f2 100644
--- a/.mule/todo.md
+++ b/.mule/todo.md
@@ -39,6 +39,10 @@ Issues are read live on each tick and are never copied into this file — a copy
 second place the same fact lives, on a branch the mule commits to and you edit on `main`,
 and that produced three merge conflicts in one afternoon.
 
+**Open issues labelled `mule` are worked BEFORE the standing lines below** (changed
+2026-08-31). An issue is work a human approved; a standing line is either a gate or something
+the mule proposed to itself, and self-proposed work never runs ahead of approved work.
+
 **To give the mule work, file an issue and label it `mule`.** From anywhere, by anyone,
 including another model with repo access. No SSH, no editing this file:
 
@@ -90,11 +94,12 @@ worth having. Before any of it reaches `main`:
 - [x] store-realm-cascade-policies: in `internal/store/realms_test.go` `CascadeDelete`, add `trigger_policies` to the post-delete verification loop (query `SELECT count(*) FROM trigger_policies WHERE realm_id = $1`). The migration 000006 has ON DELETE CASCADE but nothing asserts on it. [auto]
 - [x] store-alias-lowest-id: in `internal/store/devices_test.go`, add a subtest that registers two devices in the same realm, sets the same alias tag on both, and asserts `GetDeviceByAlias` returns the one with the lower UUID. The SQL uses `ORDER BY id LIMIT 1` but no test proves it. [auto]
 - [x] store-delete-device-objects: in `internal/store/devices_test.go` `StatsAndDelete`, insert object datastream rows for the device before deleting it, and assert they are gone after the delete. `DeleteDevice` explicitly sweeps `object_datastreams` but the test only checks individual rows. [auto]
-- [ ] control-producer-properties-compression: in `internal/engine/control.go`, accept a plaintext (and the 4-zero-byte empty-frame) device→server `producer/properties` purge list when the device declared `purge_properties_compression_format: plaintext` — `inflateProperties` currently rejects every non-zlib payload while upstream's `control_handler.ex` decodes per-capability [auto]
-- [ ] realm-config-datastream-retention: add `GET /realmmanagement/v1/<realm>/config/datastream_maximum_storage_retention` to `internal/realm/http.go` and carry `datastream_maximum_storage_retention` through `internal/housekeeping/http.go` + `internal/store/realms.go` (upstream v1.3.0 per-realm retention, v1.4.0-rc.1 adds an instance-wide default); detail in `.mule/tasks/realm-config-datastream-retention.md` [auto]
+- [ ] control-producer-properties-compression: in `internal/engine/control.go`, accept a plaintext (and the 4-zero-byte empty-frame) device→server `producer/properties` purge list when the device declared `purge_properties_compression_format: plaintext` — `inflateProperties` currently rejects every non-zlib payload while upstream's `control_handler.ex` decodes per-capability **Approved by Giulio 2026-08-31**: yes, devices that cannot compress must be able to talk to us, following upstream.
 - [ ] probe-interface-default-values: does `GET /realmmanagement/v1/<realm>/interfaces/<name>/<major>` return the same defaulted mapping-parameter values upstream reports after 3f0b864 ("Properly show default values for all mapping parameters")? report, do not patch [auto]
 - [ ] probe-value-type-validation: does Astrate's per-mapping value-type validation reject an aggregated object on an individual-value path and accept nil the way upstream's restored v1.4.0-rc.3 `validate_value_type` does? report, do not patch [auto]
 - [ ] compat-note-v1.4.0-rc.3: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.3 in .mule/for-giulio.md (do not edit the file) [auto]
 - [ ] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto]
 - [ ] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto]
-- [ ] trigger-device-deletion-events: add `device_deletion_started` and `device_deletion_finished` trigger conditions to `internal/engine/triggers/match.go` (new constants + entries in `deviceOns`), emit them from the device-deletion path in `internal/engine/engine.go`, and add them to the `on` enum in `docs/api/astarte_realm_management_api.yaml` (upstream v1.3.0-rc.0) [auto]
+- [ ] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto]
+- [ ] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto]
+- [ ] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto]

$ git log --all --oneline -S 'probe-binaryblob-validation' -- .mule/todo.md; echo "=== current state:"; rg -n "probe-binaryblob|probe-properties-on-connect|probe-mqtt-capabilities|probe-value-type|probe-interface-default" .mule/todo.md
6e30a1c queue: drop three lines already done or duplicated, approve the plaintext-purge line, take the issues-first ordering from main
=== current state:

$ git show 6e30a1c -- .mule/todo.md | rg '^-' | head -20
--- a/.mule/todo.md
-- [ ] control-producer-properties-compression: in `internal/engine/control.go`, accept a plaintext (and the 4-zero-byte empty-frame) device→server `producer/properties` purge list when the device declared `purge_properties_compression_format: plaintext` — `inflateProperties` currently rejects every non-zlib payload while upstream's `control_handler.ex` decodes per-capability [auto]
-- [ ] realm-config-datastream-retention: add `GET /realmmanagement/v1/<realm>/config/datastream_maximum_storage_retention` to `internal/realm/http.go` and carry `datastream_maximum_storage_retention` through `internal/housekeeping/http.go` + `internal/store/realms.go` (upstream v1.3.0 per-realm retention, v1.4.0-rc.1 adds an instance-wide default); detail in `.mule/tasks/realm-config-datastream-retention.md` [auto]
-- [ ] trigger-device-deletion-events: add `device_deletion_started` and `device_deletion_finished` trigger conditions to `internal/engine/triggers/match.go` (new constants + entries in `deviceOns`), emit them from the device-deletion path in `internal/engine/engine.go`, and add them to the `on` enum in `docs/api/astarte_realm_management_api.yaml` (upstream v1.3.0-rc.0) [auto]

$ git branch --show-current; echo ===; git log --all --oneline -S 'probe-binaryblob-validation' --format='%h %ad %s' --date=iso; echo ===; rg -n "probe-binaryblob|probe-object-validation" $(git rev-parse --abbrev-ref HEAD 2>/dev/null) 2>/dev/null; git grep -l 'probe-binaryblob-validation' $(git rev-list --all --max-count=40) 2>/dev/null | sort -u | head
mule/queue
===
dde6634 2026-09-04 20:43:37 +0200 mule: refill the queue
c33a929 2026-09-04 20:43:36 +0200 mule: recipe code-review ran (443s)
104d379 2026-09-04 20:13:25 +0200 mule: refill the queue
ca8e88e 2026-09-04 20:13:24 +0200 mule: recipe astarte-upstream ran (72s)
c085169 2026-09-04 20:00:09 +0200 mule: refill the queue
f0ae2a9 2026-09-04 20:00:08 +0200 mule: recipe github-issues ran (56s)
6d8cae6 2026-09-04 03:11:17 +0200 mule: daily cross-project survey, 2026-09-04
25137f2 2026-09-03 21:54:22 +0200 mule: refill the queue
0468c86 2026-09-03 21:54:22 +0200 mule: recipe astarte-upstream ran (129s)
2d82565 2026-09-03 13:06:58 +0200 mule: refill the queue
3debf0d 2026-09-03 13:06:58 +0200 mule: recipe astarte-upstream ran (105s)
ca9fcad 2026-09-02 20:43:30 +0200 mule: refill the queue
7ff4280 2026-09-02 20:43:30 +0200 mule: recipe docs-sync ran (316s)
345e129 2026-09-02 13:22:57 +0200 mule: refill the queue
2e434c1 2026-09-02 13:22:57 +0200 mule: recipe astarte-upstream ran (104s)
e68d9e3 2026-09-01 22:06:32 +0200 mule: refill the queue
48020a8 2026-09-01 22:06:32 +0200 mule: recipe docs-sync ran (259s)
d67d8ea 2026-09-01 20:32:24 +0200 mule: blocked compat-note-v140-rc5
54b406b 2026-09-01 20:13:03 +0200 mule: blocked probe-object-validation-selected-interface
4cd1e88 2026-09-01 19:57:16 +0200 mule: refill the queue
709bc6a 2026-09-01 19:57:16 +0200 mule: recipe astarte-upstream ran (478s)
e0f8823 2026-09-01 19:32:34 +0200 mule: blocked issue-91-pairing-health-serve-root
733f965 2026-09-01 19:02:01 +0200 mule: refill the queue
b69e32e 2026-09-01 19:02:00 +0200 mule: recipe github-issues ran (167s)
b8424d8 2026-09-01 13:44:04 +0200 mule: refill the queue
2777546 2026-09-01 13:44:04 +0200 mule: recipe milestones ran (607s)
49fcd62 2026-08-31 22:10:36 +0200 mule: refill the queue
45171b9 2026-08-31 22:10:36 +0200 mule: recipe docs-sync ran (203s)
8a983a3 2026-08-31 21:35:46 +0200 mule: blocked probe-binaryblob-validation
4c1228c 2026-08-31 20:42:21 +0200 mule: blocked compat-note-v140-rc3
7b03846 2026-08-31 04:12:36 +0200 mule: recipe astarte-upstream ran (224s)
1c618eb 2026-08-31 03:25:38 +0200 mule: daily cross-project survey, 2026-08-31
6e30a1c 2026-08-31 02:10:57 +0200 queue: drop three lines already done or duplicated, approve the plaintext-purge line, take the issues-first ordering from main
071d0a2 2026-08-24 04:14:45 +0200 mule: recipe astarte-upstream ran (396s)
===

The lines were all proposed on 08-31 and later **blocked**. Let me check their fate after blocking.
$ git log --all --format='%h %ad %s' --date=short -S 'probe-binaryblob-validation' | head; echo "=== issues about them:"; gh issue list --label mule --state all --limit 40 --json number,title,labels,state 2>/dev/null | grep -io "binaryblob\|object values\|object-validation\|selected interface\|rc.3\|rc.5" | sort | uniq -c
dde6634 2026-09-04 mule: refill the queue
c33a929 2026-09-04 mule: recipe code-review ran (443s)
104d379 2026-09-04 mule: refill the queue
ca8e88e 2026-09-04 mule: recipe astarte-upstream ran (72s)
c085169 2026-09-04 mule: refill the queue
f0ae2a9 2026-09-04 mule: recipe github-issues ran (56s)
6d8cae6 2026-09-04 mule: daily cross-project survey, 2026-09-04
25137f2 2026-09-03 mule: refill the queue
0468c86 2026-09-03 mule: recipe astarte-upstream ran (129s)
2d82565 2026-09-03 mule: refill the queue
=== issues about them:

$ gh issue list --label mule --state all --limit 60 --json number,title,labels,state --jq '.[] | "\(.number)\t\(.state)\t\(.title)"' 2>/dev/null | head -60
28	CLOSED	milestone-4.0: investigate Edgehog compatibility with Astrate (readonly)
20	CLOSED	bench: measure ingest-path cost of a previous-value lookup (dormant trigger types decision)
13	CLOSED	audit-dashboard-error-names: compare errorname.go's upstreamErrorNames against the current astarte-dashboard bundle

$ git show dde6634 -- .mule/todo.md | rg -i "binaryblob|object-validation|rc\.5|rc\.3"
(no output)

$ gh issue list --label mule --state all --limit 100 --jq '.[] | "\(.number)\t\(.state)\t\(.title)"' 2>/dev/null
(no output)

$ gh issue list --label mule --state open --limit 20 --json number,title,labels --jq '.[] | "\(.number)\t\(.[title])\n labels: \([.labels[].name]|join(","))"' 2>&1 | head -40
function not defined: title/0

$ gh issue list --label mule --state open --limit 30 --json number,title,labels --jq '.[] | "\(.number)\t\(.title)\t\([.labels[].name]|join(","))"' 2>&1 | head -40
(no output)

$ gh issue list --state open --limit 50 --json number,title,labels --jq '.[] | "\(.number)\t\(.title)\t\([.labels[].name]|join(","))"' 2>&1 | head -60
103	mule: nothing has landed in 16h	mule-alarm
102	mule: nothing has landed in 14h	mule-alarm
101	mule: nothing has landed in 14h	mule-alarm
100	mule: nothing has landed in 15h	mule-alarm
99	mule: nothing has landed in 14h	mule-alarm
98	mule: nothing has landed in 16h	mule-alarm
97	mule: nothing has landed in 15h	mule-alarm
96	mule: nothing has landed in 14h	mule-alarm
95	mule: nothing has landed in 15h	mule-alarm
94	mule: nothing has landed in 15h	mule-alarm
93	In internal/broker/aclhook.go the control/keyAgreement deny comment (lines ~66-69) quotes an upstream doc sentence that upstream deleted in d084308 on 2026-08-31 — rewrite the comment so it justifies the deny as Astrate's own deliberate choice (upstream now publishes a full wire spec in 082-key_agreement_protocol.md and states the handshake IS implemented on the Astarte side; scope decision tracked in issue #92), cite 082 instead of the deleted 080 sentence, keep the deny behaviour and aclhook_test.go:53 unchanged	mule-review,upstream-parity
92	keyAgreement: upstream published the full wire spec (082) — decide whether the parking from #51 still holds	enhancement,upstream-parity,upstream-experimental
78	FDO device onboarding: owner-side TO1/TO2 in Pairing (milestone-4.0 candidate)	enhancement,milestone-4.0,upstream-parity
1	Provide an Open Source IoT Platform unironically, unironically	wontfix

The rc.5 probes were proposed on 08-31 but are absent from the current queue. Let me trace what happened to them.
$ git log --all --format='%h %ad %s' --date=iso -S 'probe-binaryblob-validation' -- .mule/todo.md; echo "=== object-validation slug:"; git log --all --format='%h %ad %s' --date=iso -S 'probe-object-validation-selected-interface' -- .mule/todo.md
6e30a1c 2026-08-31 02:10:57 +0200 queue: drop three lines already done or duplicated, approve the plaintext-purge line, take the issues-first ordering from main
=== object-validation slug:
4cd1e88 2026-09-01 19:57:16 +0200 mule: refill the queue

$ git show 8a983a3 -- .mule/todo.md | head -30; echo ===; git show 54b406b -- .mule/todo.md | head -30; echo ===; git show d67d8ea -- .mule/todo.md | head -30
commit 8a983a378f4ff84d5f5cf6268d711f671148529f
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Aug 31 21:35:46 2026 +0200

    mule: blocked probe-binaryblob-validation

diff --git a/.mule/todo.md b/.mule/todo.md
index bc5eec4..44ec66f 100644
--- a/.mule/todo.md
+++ b/.mule/todo.md
@@ -101,5 +101,5 @@ worth having. Before any of it reaches `main`:
 - [!] probe-emptycache-resend-device-error: does Astrate's emptyCache server-property resend (internal/engine/control.go resendServerProperties) need to emit device_error trigger events the way upstream v1.3.3 (#2119) now does — `interface_loading_failed` when a stored property's interface is unloadable, `resend_interface_properties_failed` on a send failure? report, do not patch [auto] — BLOCKED: gates failed
 - [!] compat-note-v1.3.3: propose the docs/COMPATIBILITY.md wording for v1.3.3 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed
 - [!] probe-mqtt-capabilities-declaration: find how v1.3.x device SDKs declare MQTT v1 capabilities on the wire (upstream device-SDK sources, not the release note) and whether Astrate's parseIntrospection (internal/engine/introspection.go) would accept or reject such a payload; report, do not patch [auto] — BLOCKED: gates failed
-- [ ] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto]
+- [!] probe-binaryblob-validation: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator ("Ensure binaryblob data is correctly validated")? report against pkg/payload/value.go and internal/engine/serverdata.go, do not patch [auto] — BLOCKED: gates failed
 - [ ] probe-properties-on-connect-encoding: does Astrate encode every stored server-property value correctly when resending them to a connecting device (internal/engine/control.go resendServerProperties), as upstream v1.3.0's "correctly encode values when sending properties to device on connection" fix requires? report, do not patch [auto]
===
commit 54b406b1347c3d7eb6841d25a00effb9eb8810b5
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Tue Sep 1 20:13:03 2026 +0200

    mule: blocked probe-object-validation-selected-interface

diff --git a/.mule/todo.md b/.mule/todo.md
index 569b8ed..700cb37 100644
--- a/.mule/todo.md
+++ b/.mule/todo.md
@@ -106,5 +106,5 @@ worth having. Before any of it reaches `main`:
 - [!] docs-sync-pairing-status-enum: in docs/api/astarte_pairing_api.yaml, fix the PairingInfo.status enum (line 367): it lists `confirmed, pending, denied, expired` but the handler only ever emits `confirmed`, `pending`, or `inhibited` (internal/pairing/service.go:289-297, via internal/pairing/http.go:211) — drop the dead `denied`/`expired` and add the undocumented `inhibited`. [auto] — BLOCKED: gates failed
 - [!] issue-91-pairing-health-serve-root: in internal/pairing/http.go, add the unauthenticated `GET /pairing/health` route (no realm segment) sharing the handler/payload of the existing `GET /pairing/v1/{realm}/health`, keeping the v1 route as-is; tests for both paths (closes #91) — BLOCKED: gates failed
 - [!] issue-68-async-operation-accepted: on the mutating endpoints upstream surfaces `async_operation` on (housekeeping realm create/delete; realm-management interface install/update/delete, trigger/policy delete in internal/realm/http.go), accept and ignore `?async_operation=false`, with unparseable/`true` values also not changing behaviour (Astrate stays always-sync); tests cover the flag parsing (closes #68) — BLOCKED: gates failed
-- [ ] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto]
+- [!] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto] — BLOCKED: gates failed
 - [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]
===
commit d67d8ea28b642ce7980444a4ed8fde4387ef5868
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Tue Sep 1 20:32:24 2026 +0200

    mule: blocked compat-note-v140-rc5

diff --git a/.mule/todo.md b/.mule/todo.md
index 700cb37..17e574d 100644
--- a/.mule/todo.md
+++ b/.mule/todo.md
@@ -107,4 +107,4 @@ worth having. Before any of it reaches `main`:
 - [!] issue-91-pairing-health-serve-root: in internal/pairing/http.go, add the unauthenticated `GET /pairing/health` route (no realm segment) sharing the handler/payload of the existing `GET /pairing/v1/{realm}/health`, keeping the v1 route as-is; tests for both paths (closes #91) — BLOCKED: gates failed
 - [!] issue-68-async-operation-accepted: on the mutating endpoints upstream surfaces `async_operation` on (housekeeping realm create/delete; realm-management interface install/update/delete, trigger/policy delete in internal/realm/http.go), accept and ignore `?async_operation=false`, with unparseable/`true` values also not changing behaviour (Astrate stays always-sync); tests cover the flag parsing (closes #68) — BLOCKED: gates failed
 - [!] probe-object-validation-selected-interface: does Astrate decode and validate an object datastream against the topic-named interface's own mappings (internal/engine/data.go ci.ObjectLeaves), so identical last-level endpoint names in a second object interface cannot be type-checked against the wrong mapping the way upstream v1.4.0-rc.5 #2141 ("Make object values only validate type on selected interface") fixed? report, do not patch [auto] — BLOCKED: gates failed
-- [ ] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto]
+- [!] compat-note-v1.4.0-rc.5: propose the docs/COMPATIBILITY.md wording for v1.4.0-rc.5 in .mule/for-giulio.md (do not edit the file) [auto] — BLOCKED: gates failed

$ rg -ni "binaryblob|object.validate|selected interface|rc\.5|v1\.4\.0-rc" .mule/todo.md; echo "=== stale probline check done"; rg -n "probe-" .mule/todo.md | head -30
112:- [!] probe-property-resend-encoding [auto]: upstream v1.3.0 fixed outbound server-property values sent to a device on connect/emptyCache (commit 522ccf4f — a raw stored binaryblob is now wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding, so it ships as a BSON binary, not a string). Investigate only: does Astrate's resendServerProperties/rehydrate path (internal/engine/control.go:144, pkg/payload) emit binaryblob → BSON subtype-0 binary and datetime → UTC datetime exactly like upstream 1.3? Report, do not patch. — BLOCKED: wrote nothing
113:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
130:- [x] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
131:- [!] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it. — BLOCKED: wrote nothing
=== stale probline check done
107:- [!] probe-trigger-install-notification-delay [auto]: upstream v1.3.0 says "services now receive trigger installation and deletion notifications, which should reduce the delay between installing the trigger and starting to receive messages" — investigate only: does Astrate have an analogous delay between trigger install and first delivery? Report, do not patch. — BLOCKED: wrote nothing
112:- [!] probe-property-resend-encoding [auto]: upstream v1.3.0 fixed outbound server-property values sent to a device on connect/emptyCache (commit 522ccf4f — a raw stored binaryblob is now wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding, so it ships as a BSON binary, not a string). Investigate only: does Astrate's resendServerProperties/rehydrate path (internal/engine/control.go:144, pkg/payload) emit binaryblob → BSON subtype-0 binary and datetime → UTC datetime exactly like upstream 1.3? Report, do not patch. — BLOCKED: wrote nothing
119:- [!] probe-props-resend-error-triggers [auto]: upstream v1.3.3 (release body is empty; its one commit, fix #2119) previously folded both properties-resend failure modes into a single `resend_interface_properties_failed` device_error and now fires a distinct one — `interface_loading_failed` when the interface fails to load, `resend_interface_properties_failed` when the send to the device fails. Astrate's resend paths `resendServerProperties` (internal/engine/control.go:144) and `sendConsumerProperties` (control.go:191) skip bad rows with logs and, on the connect path, log a Warn at engine.go:359-362 — they never fire a device_error trigger at all. Investigate only: which of Astrate's resend-failure modes map to which 1.3.3 name, and whether Astrate should fire a device_error trigger at all; report, do not patch. — BLOCKED: wrote nothing
130:- [x] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.

Upstream has not moved: the newest tag is still `v1.4.0-rc.5` (2026-08-20), newest stable `v1.3.3` (2026-08-07), identical to what the last run (2026-09-15, `f0863f2`) and every run since 08-31 recorded. COMPATIBILITY.md targets 1.2.2; the whole v1.3 surface plus the v1.4-rc line (rc.5's two DUP validation fixes triaged 08-31 → probes `probe-binaryblob-validation` / `probe-object-validation-selected-interface` and `compat-note-v1.4-rc` already in the queue, blocked) is already tracked. Nothing new to propose.

Per the recipe: "If nothing changed since last time, append nothing to the queue and say 'no upstream movement since \<tag\>'. An empty result is a good result." I append no lines to `.mule/todo.md` and write no evidence file.

```
Done: no upstream movement since v1.4.0-rc.5 — newest tag unchanged (v1.4.0-rc.5, 2026-08-20; stable still v1.3.3, 2026-08-07); nothing new to propose
Files: none
Verified: gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"' -> pass (5 tags, newest v1.4.0-rc.5); compared against docs/COMPATIBILITY.md target (1.2.2, v1.3 surface delivered, v1.4 experimental) and prior queue state
Unsure: nothing — release set and prior proposals match the last recorded run exactly
Follow-ups: none
```
