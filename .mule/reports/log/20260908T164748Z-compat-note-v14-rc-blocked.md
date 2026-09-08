slug: compat-note-v14-rc
verdict: blocked
at:  ca9922b
ran: 2026-09-08T16:47:48Z on DietPi in 326s


> build · big-pickle

→ Read .mule/task.md
→ Read .mule/MULE.md
The task: propose `docs/COMPATIBILITY.md` wording for v1.4.0-rc (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library), folded into the existing v1.3.x proposal in `.mule/for-giulio.md` — without editing that file. It's a `[readonly]` probe; the deliverable is the report.
→ Read .mule/for-giulio.md
→ Read docs/COMPATIBILITY.md
Read the v1.3.x proposal at for-giulio.md:88-90. Now let me research the actual v1.4.0 items against the codebase — first the COMPATIBILITY.md structure and current emulation state.

$ ls /root/astrate-mule
AGENTS.md
Dockerfile
LICENSE
Makefile
README.md
bench
cmd
deploy
docker-compose.yml
docs
examples
go.mod
go.sum
internal
migrations
opencode.json
pkg
site-dist
test
tools

$ rg -i -l "encrypted|FDO|vault|required" internal/ --glob '*.go' | head -40
internal/flow/substitute.go
internal/config/config.go
internal/flow/blocks/mqtt.go
internal/flow/factory.go
internal/flow/blocks/transform.go
internal/config/config_test.go
internal/flow/blocks/virtualdevicepool/virtualdevicepool_test.go
internal/flow/blocks/virtualdevicepool/virtualdevicepool.go
internal/flow/blocks/info.go
internal/flow/blocks/schema.go
internal/flow/blocks/httpblocks.go
internal/flow/blocks/jsonpathmap_test.go
internal/flow/blocks/schema_test.go
internal/flow/blocks/updatemetadata.go
internal/flow/blocks/mqtt_test.go
internal/flow/blocks/updatemetadata_test.go
internal/flow/blocks/container/block.go
internal/flow/blocks/container/block_test.go
internal/flow/blocks/container/docker.go
internal/engine/router.go
internal/engine/triggers/match_test.go
internal/engine/triggers/errorname.go
internal/engine/forward/nats.go
internal/engine/forward/http.go
internal/engine/triggers/match.go
internal/engine/triggers/policy_test.go
internal/flow/blocks/randomsource.go
internal/flow/blocks/httpblocks_test.go
internal/flow/blocks/jsonpathmap.go
internal/flow/blocks/astartesource/source.go
internal/flow/blocks/catalog_test.go
internal/flow/blocks/catalog.go
internal/appengine/channels/room.go
internal/auth/jwt_test.go
internal/auth/jwt.go
internal/pairing/http.go
internal/pairing/http_test.go
internal/broker/broker.go
internal/store/crypto.go
internal/pairing/service.go

Now let me research the four v1.4.0 items against the codebase and against upstream release notes.
$ rg -n "required|encrypted|fdo|vault" internal/interface/ internal/introspection/ docs/UPSTREAM-EXPERIMENTAL.md 2>/dev/null | head -40
docs/UPSTREAM-EXPERIMENTAL.md:24:| Required + encrypted interface mapping fields | #67 | 1.4 experimental | 2026-08-23 | upstream 1.4 final: promoted or removed |

→ Read docs/UPSTREAM-EXPERIMENTAL.md

$ rg -n "required|encrypted|FDO|fdo" internal/pairing/http.go internal/pairing/service.go internal/engine/triggers/match.go internal/interface/ internal/introspection/ 2>/dev/null | head -30
internal/pairing/http.go:79:	// FDO-flow devices probe it before they hold credentials. No rate
internal/pairing/service.go:368:// design — FDO-flow devices probe it before they hold credentials).
internal/engine/triggers/match.go:403:// or concrete) it is required even for name-less added/removed triggers.
internal/engine/triggers/match.go:461:		// Name-less and concrete matchers carry the required major; a "*"

The register already fences #67 (Required + encrypted fields) at `1.4 experimental`, and nothing is implemented. Now let me verify the four v1.4.0-rc items against upstream release notes.
$ git -C /root/astrate-mule log --oneline -5 && rg -n "1\.4\.0" /root/astrate-mule/docs/*.md /root/astrate-mule/.mule/milestones.md 2>/dev/null | head -20
ca9922b mule: log probe-required-mapping-flag
db9d06e mule: probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
6c8dbf1 mule: refill the queue
f117601 mule: recipe astarte-upstream ran (305s)
1e56fda mule: refill the queue
/root/astrate-mule/.mule/milestones.md:115:  v1.4.0-rc.5 diffed service-by-service against our emulated 1.2.2 level).

◈ Exa Web Search "Astarte v1.4.0 release notes encrypted endpoints required mapping FDO v1.1 Vault"

$ git -C /root/astrate-mule log --oneline -20 --all --grep="probe-required\|required-mapping\|compat-note" ; echo "---"; git -C /root/astrate-mule show db9d06e --stat | head -20
ca9922b mule: log probe-required-mapping-flag
db9d06e mule: probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
f407d45 mule: blocked compat-note-v133
9a0da50 mule: log compat-note-v132
4bc3e1a mule: compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
d67d8ea mule: blocked compat-note-v140-rc5
21f74d8 mule: blocked compat-note-v133
4c1228c mule: blocked compat-note-v140-rc3
d145c7d mule: propose upstream v1.3.0 wire-visible changes + compat-note-v1.3.2
---
commit db9d06e23829dce94374d1c65641de34d4395dbe
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Tue Sep 8 13:42:56 2026 +0200

    mule: probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.

 .mule/todo.md | 1 +
 1 file changed, 1 insertion(+)

$ rg -n "required|compat-note|v1.4" /root/astrate-mule/.mule/todo.md | head; echo "---"; cat /root/astrate-mule/.mule/todo.md | head -30
99:- [x] compat-note-v1.3.2 [auto]: propose the docs/COMPATIBILITY.md wording update for upstream v1.3.2 (latest stable; v1.4.0 is still rc-only) in .mule/for-giulio.md — do not edit docs/COMPATIBILITY.md directly, it is on the never-touch list.
104:- [!] compat-note-v1.3.3 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.3 (newest stable, 2026-08-07, empty-body patch release; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 in rather than re-deriving it. — BLOCKED: wrote nothing
121:- [x] probe-required-mapping-flag [auto]: upstream v1.4.0-rc.0 added a `required` boolean field on object-aggregated interface mappings — validated at install time in RM, enforced by DUP and AppEngine on data writes (commits #1846/#1847/#1849/#1854). Investigate only: does Astrate's Realm Management accept and persist the `required` field on interface installation (internal/realm/service.go install path), and does the DUP/AppEngine data-validation path reject writes that omit a `required` mapping? Report the gap, do not patch.
122:- [ ] compat-note-v1.4-rc [auto]: propose the docs/COMPATIBILITY.md wording for the v1.4.0-rc line (encrypted endpoints, `required` mapping flag, FDO v1.1, Vault/config library) in .mule/for-giulio.md (do not edit the file) — note that v1.4.0 is still RC-only and the wording should reflect "experimental, not yet emulated" until a stable v1.4.0 tag ships. Fold into the existing v1.3.x wording proposal already in for-giulio.md rather than re-deriving it.
123:- [ ] detailed-listing-required-encrypted-flags: emit the `required`/`encrypted` mapping flags in `writeDetailedMapping` (internal/realm/detailed.go) so `GET /realmmanagement/v1/{realm}/interfaces?detailed=true` round-trips them — the detailed listing currently materializes a fixed field set that drops both flags even when a datastream mapping declares one, while the plain GET preserves them via the stored definition. Datastream branch only (properties reject both at parse); add a T1 assertion in internal/realm/detailed_test.go.
---
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

There's prior related work — a `probe-required-mapping-flag` task. Let me check its result and the todo queue for context on how these v1.4 items are tracked.
$ ls /root/astrate-mule/.mule/reports/ 2>/dev/null; echo "---"; cat /root/astrate-mule/.mule/reports/probe-required-mapping-flag.md 2>/dev/null
hygiene.md
log
---

→ Read .mule/milestones.md [offset=95, limit=40]
Let me get the full v1.4.0-rc release-note picture to describe the four items accurately.
% WebFetch https://github.com/astarte-platform/astarte/releases/tag/v1.4.0-rc.5
rc.5 is thin; rc.0 has the v1.4 additions. Let me get the full rc.0 notes.
% WebFetch https://github.com/astarte-platform/astarte/releases/tag/v1.4.0-rc.0

$ cat /root/astrate-mule/.mule/reports/log | tail -40
cat: /root/astrate-mule/.mule/reports/log: Is a directory

Let me check the probe's logged report and how Astrate's schema parsing treats `required`/`encrypted`.
$ rg -n "required|encrypted" /root/astrate-mule/internal/realm/detailed.go | head; echo "==="; rg -n "required|encrypted" /root/astrate-mule/internal/interface/ /root/astrate-mule/internal/realm/*.go 2>/dev/null | grep -iv "test" | head -20
===

The report log is a directory. Let me look inside it and verify the schema-parse state directly.
$ ls -la /root/astrate-mule/.mule/reports/ /root/astrate-mule/.mule/reports/log/ 2>/dev/null; echo "==="; ls /root/astrate-mule/internal/realm/
/root/astrate-mule/.mule/reports/:
total 16
drwxr-xr-x 3 root root 4096 Sep  6 19:40 .
drwxr-xr-x 7 root root 4096 Sep  8 18:42 ..
-rw-r--r-- 1 root root 2279 Sep  6 19:40 hygiene.md
drwxr-xr-x 2 root root 4096 Sep  8 13:42 log

/root/astrate-mule/.mule/reports/log/:
total 1596
drwxr-xr-x 2 root root   4096 Sep  8 13:42 .
drwxr-xr-x 3 root root   4096 Sep  6 19:40 ..
-rw-r--r-- 1 root root  19332 Sep  4 21:09 20260904T190911Z-purge-properties-compression-capabilityauto-done.md
-rw-r--r-- 1 root root  25676 Sep  4 21:50 20260904T195018Z-broker-acl-coldstart-introspection-miss-done.md
-rw-r--r-- 1 root root   5576 Sep  4 21:55 20260904T195529Z-issue-93-done.md
-rw-r--r-- 1 root root    970 Sep  4 21:58 20260904T195814Z-broker-disconnect-device-zombie-session-blocked.md
-rw-r--r-- 1 root root   7195 Sep  4 22:30 20260904T203003Z-broker-offline-acl-tests-blocked.md
-rw-r--r-- 1 root root   3691 Sep  4 23:05 20260904T210532Z-broker-onconnect-doc-comment-blocked.md
-rw-r--r-- 1 root root  15724 Sep  5 13:02 20260905T110211Z-empty-introspection-verification-blocked.md
-rw-r--r-- 1 root root   4497 Sep  5 13:14 20260905T111428Z-probe-trigger-install-notification-delay-blocked.md
-rw-r--r-- 1 root root  12798 Sep  5 13:25 20260905T112539Z-compat-note-v132-done.md
-rw-r--r-- 1 root root  31514 Sep  5 18:55 20260905T165512Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root  13915 Sep  5 19:17 20260905T171755Z-docs-sync-pairing-health-path-done.md
-rw-r--r-- 1 root root   5035 Sep  5 19:22 20260905T172234Z-docs-sync-pairing-register-404-done.md
-rw-r--r-- 1 root root  32904 Sep  5 19:36 20260905T173657Z-recipe-hygiene-proposed.md
-rw-r--r-- 1 root root  11141 Sep  5 19:45 20260905T174528Z-swagger-httptest-coverage-done.md
-rw-r--r-- 1 root root 120837 Sep  5 20:10 20260905T181053Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  24739 Sep  5 20:12 20260905T181219Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root 152510 Sep  5 20:46 20260905T184625Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root   9342 Sep  5 20:49 20260905T184903Z-probe-property-resend-encoding-blocked.md
-rw-r--r-- 1 root root   2052 Sep  5 21:08 20260905T190855Z-compat-note-v133-blocked.md
-rw-r--r-- 1 root root  38218 Sep  5 21:20 20260905T192010Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root  42265 Sep  5 21:29 20260905T192911Z-flow-validate-source-sink-blocked.md
-rw-r--r-- 1 root root   6007 Sep  5 22:15 20260905T201529Z-flow-validate-dead-source-sink-recompute-done.md
-rw-r--r-- 1 root root  41275 Sep  6 13:20 20260906T112006Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root  12532 Sep  6 13:31 20260906T113143Z-docs-sync-rm-datastream-retention-endpoint-done.md
-rw-r--r-- 1 root root  19509 Sep  6 19:10 20260906T171045Z-docs-sync-rm-interfaces-detailed-param-done.md
-rw-r--r-- 1 root root  12896 Sep  6 19:40 20260906T174031Z-recipe-hygiene-proposed.md
-rw-r--r-- 1 root root  17313 Sep  6 19:57 20260906T175748Z-housekeeping-tests-done.md
-rw-r--r-- 1 root root  63922 Sep  6 20:04 20260906T180442Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root   8576 Sep  6 20:13 20260906T181352Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  57962 Sep  6 20:31 20260906T183102Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  24479 Sep  6 20:39 20260906T183952Z-probe-props-resend-error-triggers-blocked.md
-rw-r--r-- 1 root root  26465 Sep  6 20:57 20260906T185735Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root  27807 Sep  6 21:58 20260906T195827Z-server-data-trigger-bus-done.md
-rw-r--r-- 1 root root  39030 Sep  6 22:04 20260906T200431Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root  15579 Sep  6 22:10 20260906T201002Z-docs-sync-appengine-by-alias-endpoints-done.md
-rw-r--r-- 1 root root  39405 Sep  6 22:27 20260906T202744Z-docs-sync-appengine-group-endpoints-done.md
-rw-r--r-- 1 root root   8072 Sep  6 22:43 20260906T204333Z-docs-sync-appengine-get-group-device-done.md
-rw-r--r-- 1 root root    218 Sep  7 04:03 20260907T020359Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  46239 Sep  7 13:24 20260907T112447Z-docs-sync-appengine-query-params-status-done.md
-rw-r--r-- 1 root root   3915 Sep  7 13:26 20260907T112654Z-docs-sync-appengine-group-patch-status-done.md
-rw-r--r-- 1 root root  11879 Sep  7 19:00 20260907T170024Z-docs-sync-appengine-data-422-interface-level-done.md
-rw-r--r-- 1 root root   3648 Sep  7 19:04 20260907T170452Z-recipe-hygiene-proposed.md
-rw-r--r-- 1 root root   8038 Sep  7 19:19 20260907T171928Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  49686 Sep  7 19:48 20260907T174812Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  82600 Sep  7 19:55 20260907T175524Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  47843 Sep  7 20:15 20260907T181549Z-recipe-code-review-proposed.md
-rw-r--r-- 1 root root   9669 Sep  7 20:48 20260907T184850Z-appengine-snapshot-ignores-query-params-blocked.md
-rw-r--r-- 1 root root   5653 Sep  7 21:03 20260907T190319Z-appengine-group-token-roundtrip-test-blocked.md
-rw-r--r-- 1 root root  34209 Sep  7 21:11 20260907T191157Z-recipe-docs-sync-proposed.md
-rw-r--r-- 1 root root  28430 Sep  7 22:12 20260907T201232Z-docs-sync-appengine-data-output-params-done.md
-rw-r--r-- 1 root root   4535 Sep  7 22:18 20260907T201806Z-recipe-hygiene-proposed.md
-rw-r--r-- 1 root root  11829 Sep  7 22:47 20260907T204751Z-recipe-milestones-proposed.md
-rw-r--r-- 1 root root  49792 Sep  8 13:20 20260908T112032Z-recipe-github-issues-proposed.md
-rw-r--r-- 1 root root  85176 Sep  8 13:34 20260908T113422Z-recipe-astarte-upstream-proposed.md
-rw-r--r-- 1 root root  23027 Sep  8 13:42 20260908T114256Z-probe-required-mapping-flag-done.md
===
async_operation_test.go
dashboard_compat_test.go
detailed.go
detailed_test.go
http.go
http_test.go
policy.go
policy_test.go
service.go

$ rg -rn "required|encrypted" /root/astrate-mule/internal/interface/ 2>/dev/null | head -15; echo "===realm_files==="; rg -ln "writeDetailedMapping|Detailed" /root/astrate-mule/internal/realm/ 2>/dev/null
===realm_files===
/root/astrate-mule/internal/realm/http.go
/root/astrate-mule/internal/realm/detailed.go
/root/astrate-mule/internal/realm/service.go
/root/astrate-mule/internal/realm/detailed_test.go

The probe report exists. Let me read it plus the v1.3.2 compat-note report for the existing proposal pattern.
→ Read .mule/reports/log/20260908T114256Z-probe-required-mapping-flag-done.md
The probe found `required` is **already implemented** (#67, closed 2026-09-04). Let me nail exactly how `encrypted` is handled in the parser so the wording is precise about "encrypted endpoints".

$ rg -n "encrypted|Encrypted" /root/astrate-mule/pkg/interfaceschema/*.go | head -20
/root/astrate-mule/pkg/interfaceschema/parse.go:97:	// Encrypted records the declared encryption flag; parsed and stored,
/root/astrate-mule/pkg/interfaceschema/parse.go:99:	Encrypted bool
/root/astrate-mule/pkg/interfaceschema/parse.go:148:	Encrypted               *bool                    `json:"encrypted,omitempty"`
/root/astrate-mule/pkg/interfaceschema/parse.go:397:		case raw.Encrypted != nil:
/root/astrate-mule/pkg/interfaceschema/parse.go:398:			return m, nil, fmt.Errorf(`mapping %q: "encrypted" is not allowed on properties`, m.Endpoint)
/root/astrate-mule/pkg/interfaceschema/parse.go:413:	if raw.Encrypted != nil {
/root/astrate-mule/pkg/interfaceschema/parse.go:414:		m.Encrypted = *raw.Encrypted
/root/astrate-mule/pkg/interfaceschema/compat.go:123:	case nm.Encrypted != om.Encrypted:
/root/astrate-mule/pkg/interfaceschema/compat.go:124:		attr = "encrypted"
/root/astrate-mule/pkg/interfaceschema/required_encrypted_test.go:11:func TestParseRequiredEncryptedDatastream(t *testing.T) {
/root/astrate-mule/pkg/interfaceschema/required_encrypted_test.go:22:		if iface.Mappings[0].Encrypted {
/root/astrate-mule/pkg/interfaceschema/required_encrypted_test.go:23:			t.Error("Encrypted = true, want default false")
/root/astrate-mule/pkg/interfaceschema/required_encrypted_test.go:27:	t.Run("encrypted true parses", func(t *testing.T) {
/root/astrate-mule/pkg/interfaceschema/required_encrypted_test.go:29:			"interface_name": "com.astrate.test.EncryptedField",
/root/astrate-mule/pkg/interfaceschema/required_encrypted_test.go:32:			"mappings": [{"endpoint": "/v", "type": "integer", "encrypted": true}]
/root/astrate-mule/pkg/interfaceschema/required_encrypted_test.go:34:		if !iface.Mappings[0].Encrypted {
/root/astrate-mule/pkg/interfaceschema/required_encrypted_test.go:35:			t.Error("Encrypted = false, want true")
/root/astrate-mule/pkg/interfaceschema/required_encrypted_test.go:49:		if iface.Mappings[0].Required || iface.Mappings[0].Encrypted {
/root/astrate-mule/pkg/interfaceschema/required_encrypted_test.go:50:			t.Errorf("Required = %v, Encrypted = %v, want both false",
/root/astrate-mule/pkg/interfaceschema/required_encrypted_test.go:51:				iface.Mappings[0].Required, iface.Mappings[0].Encrypted)

$ rg -in "fdo|vault|rendezvous|encrypt" /root/astrate-mule/internal /root/astrate-mule/pkg --glob '*.go' --glob '!**/*_test.go' | rg -v "flow|broker" | head -20; echo "==="; rg -c "Encrypted" /root/astrate-mule/pkg/payload/*.go 2>/dev/null
/root/astrate-mule/pkg/interfaceschema/compat.go:123:	case nm.Encrypted != om.Encrypted:
/root/astrate-mule/pkg/interfaceschema/compat.go:124:		attr = "encrypted"
/root/astrate-mule/pkg/interfaceschema/trie.go:36:	// Encrypted records the declared encryption flag; parsed and stored,
/root/astrate-mule/pkg/interfaceschema/trie.go:38:	Encrypted bool
/root/astrate-mule/pkg/interfaceschema/compile.go:82:			Encrypted:         m.Encrypted,
/root/astrate-mule/internal/config/config.go:5:// identity must be supplied. The master encryption key is referenced here but
/root/astrate-mule/internal/config/config.go:133:// SecurityConfig references the master encryption key that seals realm CA
/root/astrate-mule/pkg/interfaceschema/parse.go:97:	// Encrypted records the declared encryption flag; parsed and stored,
/root/astrate-mule/pkg/interfaceschema/parse.go:99:	Encrypted bool
/root/astrate-mule/pkg/interfaceschema/parse.go:148:	Encrypted               *bool                    `json:"encrypted,omitempty"`
/root/astrate-mule/pkg/interfaceschema/parse.go:397:		case raw.Encrypted != nil:
/root/astrate-mule/pkg/interfaceschema/parse.go:398:			return m, nil, fmt.Errorf(`mapping %q: "encrypted" is not allowed on properties`, m.Endpoint)
/root/astrate-mule/pkg/interfaceschema/parse.go:413:	if raw.Encrypted != nil {
/root/astrate-mule/pkg/interfaceschema/parse.go:414:		m.Encrypted = *raw.Encrypted
/root/astrate-mule/internal/store/crypto.go:15:// CA private keys are encrypted at rest (docs/DESIGN.md §2.2, §4.3): the
/root/astrate-mule/internal/store/crypto.go:116:// Seal encrypts plaintext, returning nonce || ciphertext+tag.
/root/astrate-mule/internal/pairing/ca/ca.go:149:// plaintext handed to store.KeySealer for at-rest encryption.
===

Facts confirmed. Nothing uses `fdo`/`vault`/`rendezvous`; `encrypted` exists only as a schema flag (parsed, stored, no payload crypto). `required` is enforced (`payload.Oject` → ReasonMissingRequired). This is a proposal-only task (no file edits allowed, no code), so the deliverable is the wording in my report.

---

**Proposed wording** (folded into the standing v1.3.x COMPATIBILITY.md proposal at `for-giulio.md:88-90` — replace that sentence's single clause with this, rather than adding a second sentence):

> In §Infrastructure differences: *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream **v1.3.x** capabilities (MQTT v1 capabilities incl. plaintext `purge_properties_compression_format`, empty-introspection allowance, device registration/deletion triggers, experimental FDO pairing auth) are not yet emulated and are out of scope until the milestone that adopts v1.3.2 as the target. Upstream **v1.4.0** features are experimental until a stable v1.4.0 tag ships (currently v1.4.0-rc.5-only) and are likewise not yet emulated: **encrypted endpoints** (the `encrypted` mapping flag is accepted and persisted on datastream mappings but the E2E payload encryption is not — there is no key negotiation or decrypt path), **FDO v1.1 device registration** (astarte_pairing; parked here as milestone-4.0 issue #78), and the **Vault-backed secret storage + shared config library** upstream introduced (no wire impact; Astrate is a single binary with an embedded CA and at-rest-sealed realm secrets, `internal/store/crypto.go`). The one early-adopted slice is the `required`/`encrypted` mapping-flag **fields** (UPSTREAM-EXPERIMENTAL.md #67, closed 2026-09-04): Astrate accepts, persists and — for `required` — enforces them (missing required leaves are rejected on both JSON and BSON ingest) at the 1.2.2 tag level, reconciled again at upstream v1.4 final."*

Reasoning: the `required` mapping flag must NOT be painted "not yet emulated" — it already landed as #67; only the encrypted-endpoint *feature* (actual crypto), FDO v1.1 and the Vault/config infra remain un-emulated, and FDO/Vault are infra/decision items (parked `#78`) rather than wire gaps.

```
Done: Proposed the docs/COMPATIBILITY.md v1.4.0-rc wording (encrypted endpoints, `required` flag, FDO v1.1, Vault/config library) folded into the standing v1.3.x proposal — no file edited, per the task's do-not-edit instruction
Files: none (proposal-only; deliverable is this report)
Verified: rg across internal/, pkg/ and docs/ confirmed no FDO/Vault/rendezvous code exists, `encrypted` is parsed+stored but never cryptographically consumed, `required` is enforced (payload.Object, ReasonMissingRequired); probe report .mule/reports/log/20260908T114256Z-probe-required-mapping-flag-done.md and UPSTREAM-EXPERIMENTAL.md #67 confirm the flag fields already landed (2026-09-04); upstream v1.4.0-rc.0/rc.5 release notes verified the four items -> n/a (no code, so no go vet/test gate run)
Unsure: Exact upstream wording of the "config library" (I kept the task's own "Vault/config library" unit without naming an upstream library); whether upstream rejects `encrypted` on individual-aggregation datastreams at install — carried over from the probe, not re-verified against live upstream
Follow-ups: none
```
