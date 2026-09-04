slug: purge-properties-compression-capabilityauto
verdict: done
at:  6f7a3d6
ran: 2026-09-04T19:09:11Z on DietPi in 86s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
Let me understand the task. It's about checking whether Astrate's capabilities handling needs to recognize the upstream `purge_properties_compression_format` device capability. First I need to look at internal/broker and the capabilities topic (issue #16).

$ rg -l "capabil" /root/astrate-mule --ignore-case -g '!vendor' -g '!*.sum' | head -50
/root/astrate-mule/docs/ROADMAP.md
/root/astrate-mule/docs/handoff/flow-v2-decisions-2026-07-29.md
/root/astrate-mule/docs/handoff/flow-parity-audit-2026-07-29.md
/root/astrate-mule/migrations/000007_toolkit.up.sql
/root/astrate-mule/docs/DESIGN.md
/root/astrate-mule/internal/engine/router.go
/root/astrate-mule/internal/store/store.go
/root/astrate-mule/internal/engine/data.go
/root/astrate-mule/internal/engine/engine.go
/root/astrate-mule/internal/engine/control.go
/root/astrate-mule/internal/engine/data_test.go
/root/astrate-mule/internal/engine/capabilities.go
/root/astrate-mule/internal/engine/topics_test.go
/root/astrate-mule/internal/engine/topics.go
/root/astrate-mule/internal/engine/capabilities_test.go
/root/astrate-mule/internal/engine/cache.go
/root/astrate-mule/internal/broker/broker.go
/root/astrate-mule/internal/broker/aclhook.go
/root/astrate-mule/internal/broker/aclhook_test.go

$ gh issue view 16
title:	capabilities-topic: handle the <realm>/<device_id>/capabilities BSON topic
state:	CLOSED
author:	atsetilam (Giulio Malitesta)
labels:	
comments:	2
assignees:	
projects:	
milestone:	
issue-type:	
parent:	
sub-issues:	
sub-issues-completed:	
blocked-by:	
blocking:	
number:	16
--
Upstream's MQTT v1 protocol (`080-mqtt-v1-protocol.md:119-136`) defines a
`<realm>/<device_id>/capabilities` topic with a BSON payload that sets device capabilities;
today the only capability is `purge_properties_compression_format` (`"zlib"` or `"plaintext"`,
default `"zlib"`). Astrate's broker ACL (`internal/broker/aclhook.go`) does not allow
publishing to this topic and the engine does not parse it — confirmed absent by grep. This is
a live, exercised upstream feature (CHANGELOG v1.3.1 fixed a "capabilities message stalling
queues" bug), not a theoretical one.

This is bigger than a one-line task — see the companion spec this issue's mule task will pick
up automatically (`.mule/tasks/issue-16.md`, added in the same commit as this issue).

Found by a cross-project survey comparing Astrate against upstream astarte-platform
(`.mule/research/survey-2026-07-27.md`, source 1/4, source 4 table row 1).

There's already a `.mule/tasks/issue-16.md` companion spec and an `internal/engine/capabilities.go`. Let me read them.
$ ls /root/astrate-mule/.mule/tasks/ 2>/dev/null
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md

→ Read internal/engine/capabilities.go
→ Read internal/broker/aclhook.go
Interesting. The capabilities feature has actually already been implemented — there's `internal/engine/capabilities.go` with `handleCapabilities`, and the ACL already allows `capabilities` publishes (`aclhook.go:62`). Let me check the current task slug: `purge-properties-compression-capability`.

The task asks whether Astrate's capabilities handling needs to recognize/honour the `purge_properties_compression_format` capability, or whether zlib-only is already the deliberate default.

Wait — this appears to already be implemented. Let me check the current state. Issue #16 is CLOSED. So the capabilities topic handling is done. Let me check whether the handler actually honors the capability (uses zlib vs plaintext) in consumer/properties delivery.

Let me look at how `purgeCompressionFor` is used and whether zlib vs plaintext is actually honored in the consumer/properties payload.

$ cat /root/astrate-mule/.mule/tasks/issue-16.md
# issue-16 — handle the `<realm>/<device_id>/capabilities` BSON topic

<!-- trickle-allow: internal/engine/topics.go internal/engine/topics_test.go internal/engine/control.go internal/engine/control_test.go internal/broker/aclhook.go internal/broker/aclhook_test.go -->

## Context

Astrate classifies every inbound device MQTT publish in `internal/engine/topics.go`'s
`classify()` (lines 61-72): the bare device root is introspection, `control/...` is the
control channel, everything else falls through to `kindData` and gets matched against the
device's introspected interfaces. Upstream Astarte's MQTT v1 protocol also defines a
`<realm>/<device_id>/capabilities` topic (sibling of the root, NOT under `control/`) carrying
a BSON payload that sets device capabilities. Today that topic falls into the `kindData`
default case and Astrate tries (and fails) to match it against an interface — it is not
denied by the broker, but nothing meaningful happens with it, and the ACL doesn't even admit
the *publish* in the first place (see below).

Only one capability exists today: `purge_properties_compression_format`, value `"zlib"` or
`"plaintext"` (default `"zlib"`). It controls the format of the `consumer/properties` purge
message Astrate already sends from `sendConsumerProperties` (`internal/engine/control.go:182`,
read the surrounding ~20 lines) and from `handleEmptyCache`'s resync path
(`internal/engine/control.go:57-74`, which calls `formatForHint`/`resendServerProperties`
already — read `formatForHint` to see the existing hint mechanism this capability extends).

The broker ACL's write-side check is `checkACL` in `internal/broker/aclhook.go` (lines 50-71):
`rest == "control/emptyCache" || rest == "control/producer/properties"` is the exact pattern to
extend with a `capabilities` case (note: capabilities is NOT under `control/`, so this is a new
top-level `rest == "capabilities"` branch, not an addition to the control list).

## What to do

1. **`internal/engine/topics.go`**: add a `kindCapabilities` value to the `topicKind` enum
   (alongside `kindIntrospection`/`kindControl`/`kindData`) and a case in `classify()` for
   `rest == "capabilities"`.
2. **`internal/broker/aclhook.go`**: in `checkACL`'s write branch, allow `rest == "capabilities"`
   the same way `control/emptyCache` is allowed (publish-only; no read/subscribe/delivery side
   — capabilities is a device-to-server topic).
3. **`internal/engine/control.go`** (or wherever the router dispatches on `topicKind` — check
   `internal/engine/router.go:104` for how `kindControl` is currently wired from
   `handleControl`'s caller, and mirror that for the new `kindCapabilities` case): add
   `handleCapabilities(ctx, m, realm)` that:
   - Decodes the BSON payload into a struct with (at least)
     `purge_properties_compression_format` as an optional string.
   - Validates the value is `"zlib"`, `"plaintext"`, or absent (defaulting to `"zlib"`) —
     reject/log anything else the same way other malformed-payload paths in this file do
     (see `handleProducerProperties` for the existing malformed-payload rejection pattern,
     `internal/engine/control.go:76` onward).
   - Stores the capability per-device (find the existing per-device state mechanism —
     `dev.hint()`/`dev.armHintReset()` used by `handleEmptyCache` is the pattern for
     per-device sticky state; capabilities likely belongs alongside it, not as a new store
     table, since it's a live-connection property upstream keeps in memory too — confirm this
     against how `dev` state is defined before adding a new persistence path).
   - On a malformed message, disconnects the device (upstream's documented behaviour) rather
     than just rejecting the message — check how other hard-reject paths in this engine close
     a device session, and reuse that, don't invent a new disconnect mechanism.

## Constraints

- **Do not implement anything beyond `purge_properties_compression_format`.** It is the only
  capability upstream defines today; a forward-compatible unknown-field-tolerant BSON decode
  is enough headroom.
- **Do not change `control/emptyCache` or `control/producer/properties` behaviour.** Existing
  tests for those must stay green unchanged — this is an additive topic, not a rework of the
  control channel.
- If the "per-device state" question above turns out to be ambiguous or to need a schema
  change, say so in your summary instead of guessing a persistence design — that would make
  this a design decision, not a mule task.

## Acceptance criteria

Run from the repo root, inside the mule's worktree:

```sh
go build ./...
go test ./internal/engine/... ./internal/broker/... -run Capabilit
go test ./internal/engine/... ./internal/broker/...   # nothing else regresses
```

Tests to add:

- `internal/engine/topics_test.go`: `classify("capabilities")` returns `kindCapabilities`.
- `internal/broker/aclhook_test.go`: a write ACL check for `base + "/capabilities"` is allowed;
  a read/subscribe check for the same topic is still denied (capabilities is publish-only).
- `internal/engine/control_test.go` (or wherever `handleEmptyCache`/`handleProducerProperties`
  are tested): a valid capabilities BSON message sets the compression format, observable via
  whatever `resendServerProperties`/`sendConsumerProperties` already exposes for testing; a
  malformed capabilities message results in the device being disconnected (assert on whatever
  signal the existing malformed-payload tests in this file already assert on).

$ rg -n "purgeCompressionFor|purgeCompressionFormat|compressionZlib|compressionPlaintext|setPurgeCompression" /root/astrate-mule --ignore-case
/root/astrate-mule/internal/engine/cache.go:404:// setPurgeCompression stores the device's preferred consumer/properties
/root/astrate-mule/internal/engine/cache.go:407:func (d *deviceState) setPurgeCompression(format string) {
/root/astrate-mule/internal/engine/cache.go:413:// purgeCompressionFormat returns the device's preferred consumer/properties
/root/astrate-mule/internal/engine/cache.go:415:func (d *deviceState) purgeCompressionFormat() string {
/root/astrate-mule/internal/engine/capabilities.go:19:	compressionZlib      = "zlib"
/root/astrate-mule/internal/engine/capabilities.go:20:	compressionPlaintext = "plaintext"
/root/astrate-mule/internal/engine/capabilities.go:46:		dev.setPurgeCompression(v)
/root/astrate-mule/internal/engine/capabilities.go:77:		case compressionZlib, compressionPlaintext:
/root/astrate-mule/internal/engine/capabilities.go:81:				value, key, compressionZlib, compressionPlaintext)
/root/astrate-mule/internal/engine/capabilities.go:89:// purgeCompressionFor returns the consumer/properties compression format for
/root/astrate-mule/internal/engine/capabilities.go:93:func purgeCompressionFor(dev *deviceState) string {
/root/astrate-mule/internal/engine/capabilities.go:95:		return compressionZlib
/root/astrate-mule/internal/engine/capabilities.go:97:	f := dev.purgeCompressionFormat()
/root/astrate-mule/internal/engine/capabilities.go:99:		return compressionZlib
/root/astrate-mule/internal/engine/capabilities_test.go:102:	if f := dev.purgeCompressionFormat(); f != "zlib" {
/root/astrate-mule/internal/engine/capabilities_test.go:114:	if f := dev.purgeCompressionFormat(); f != "plaintext" {
/root/astrate-mule/internal/engine/capabilities_test.go:232:// TestPurgeCompressionFor tests the fallback logic.
/root/astrate-mule/internal/engine/capabilities_test.go:233:func TestPurgeCompressionFor(t *testing.T) {
/root/astrate-mule/internal/engine/capabilities_test.go:234:	if got := purgeCompressionFor(nil); got != compressionZlib {
/root/astrate-mule/internal/engine/capabilities_test.go:235:		t.Errorf("purgeCompressionFor(nil) = %q, want zlib", got)
/root/astrate-mule/internal/engine/capabilities_test.go:238:	if got := purgeCompressionFor(dev); got != compressionZlib {
/root/astrate-mule/internal/engine/capabilities_test.go:239:		t.Errorf("purgeCompressionFor(empty) = %q, want zlib", got)
/root/astrate-mule/internal/engine/capabilities_test.go:241:	dev.setPurgeCompression("plaintext")
/root/astrate-mule/internal/engine/capabilities_test.go:242:	if got := purgeCompressionFor(dev); got != "plaintext" {
/root/astrate-mule/internal/engine/capabilities_test.go:243:		t.Errorf("purgeCompressionFor(plaintext) = %q, want plaintext", got)
/root/astrate-mule/internal/engine/control.go:204:	compression := compressionZlib
/root/astrate-mule/internal/engine/control.go:206:		compression = purgeCompressionFor(dev)
/root/astrate-mule/internal/engine/control.go:210:	case compressionPlaintext:

The issue-16 task appears to already be implemented. Let me verify the state: `capabilities.go` exists, `handleCapabilities` exists, `aclhook.go` allows `capabilities` publish, `classify` presumably has a `kindCapabilities`. Let me check the current task slug again.

The task is `purge-properties-compression-capability`. It says: check whether Astrate's capabilities handling needs to recognize/honour `purge_properties_compression_format`, or whether zlib-only is already the deliberate default; propose the change or note why not needed.

But issue #16 is CLOSED and the code already exists and honours both zlib and plaintext. This includes actual compression in `control.go:200-210`. So the feature is fully implemented.

Let me check git history to understand the state. Maybe this task is asking me to verify that the implementation is complete/honoured. Let me look at `control.go` around the compression usage and check whether the feature is actually being honoured in outgoing control messages.

Let me look at the control.go compression code and whether there's anything remaining. Also let me check the git log to see if capabilities.go was already committed.
→ Read internal/engine/control.go [offset=180, limit=60]

$ cd /root/astrate-mule && git log --oneline -15 -- internal/engine/capabilities.go internal/engine/control.go internal/engine/topics.go internal/broker/aclhook.go 2>/dev/null
06cb17a engine: tolerate unknown control subpaths silently (issue #51)
4c881d0 fix(lint): resolve golangci-lint v2.12.2 failures across flow container work
aba2d18 fix: land blocked mule issues #14 #21 #22 #24 #27
4e6db62 mule: issue-16: capabilities-topic: handle the <realm>/<device_id>/capabilities BSON topic
699cf11 Initial commit under updated license terms

The task (`purge-properties-compression-capability`) is asking whether Astrate's capabilities handling needs to recognize/honour `purge_properties_compression_format`. This has already been fully implemented in commit `4e6db62` (issue-16). The capability is recognized (`capabilities.go`), the ACL allows the topic, and the engine honors both `zlib` and `plaintext` compression in `sendConsumerProperties` and the emptyCache resync path.

Let me verify the git status to confirm the working tree is clean (whether this is already fully landed) and check whether this task is a stale/duplicate in the queue.
$ cd /root/astrate-mule && git status --short && echo "---LOG---" && git log --oneline -3
?? .mule/cron.log
---LOG---
ad3be03 mule: keep the broker code review the last tick produced
c0f28a7 mule: split the race-check task five ways, as mule/queue had
6bc95dc mule: a failing gate now says what broke, and carry the branch's notes to main

$ cd /root/astrate-mule && rg -n "purge-properties-compression-capability|purge_properties_compression|compression" .mule/todo.md 2>/dev/null
91:- [ ] purge-properties-compression-capability [auto]: upstream v1.3.0 adds a `purge_properties_compression_format` device capability (`zlib`|`plaintext`, default `zlib`) — a wire-visible capability value. Check whether Astrate's capabilities handling (internal/broker, the `<realm>/<device_id>/capabilities` topic, issue #16) needs to recognize/honour it, or whether zlib-only is already the deliberate default; propose the change or note why not needed.

The feature is already fully implemented and landed. The task is essentially asking me to verify whether the capability handling needs work — and it does not, because it's already done.

Let me confirm the capability is wire-honoured (not just parsed) — it is: `sendConsumerProperties` (control.go:204-217) uses `purgeCompressionFor(dev)` to pick zlib vs plaintext, and the emptyCache resync path also uses it. Let me quickly check the emptyCache path and the topics.go classify to be complete.
$ cd /root/astrate-mule && rg -n "kindCapabilities|capabilities" internal/engine/topics.go internal/engine/router.go
internal/engine/topics.go:23:	// kindCapabilities is "<realm>/<device_id>/capabilities".
internal/engine/topics.go:24:	kindCapabilities
internal/engine/topics.go:62:// the control handler rejects. The capabilities topic is a standalone kind
internal/engine/topics.go:63:// because it carries a BSON document of device capabilities, not a control
internal/engine/topics.go:69:	case rest == "capabilities":
internal/engine/topics.go:70:		return kindCapabilities, ""
internal/engine/router.go:118:	// onCapabilities handles `<realm>/<device_id>/capabilities` BSON

→ Read internal/engine/control.go [offset=50, limit=40]
The implementation is complete and fully landed. Both `sendConsumerProperties` (the purge message) honors the capability via `purgeCompressionFor(dev)`, and `handleEmptyCache`'s resync path calls `sendConsumerProperties` which honors it too.

This task is fully implemented already. The capability:
- **Recognized**: `capabilities.go` parses and validates `purge_properties_compression_format` (zlib|plaintext)
- **Topic handled**: `topics.go:69` `kindCapabilities`, `router.go:118` dispatch
- **ACL allows**: `aclhook.go:62` allows `capabilities` publish
- **Honoured (not just parsed)**: `control.go:204-217` `sendConsumerProperties` picks zlib vs plaintext frame based on the device capability; the emptyCache resync path goes through `sendConsumerProperties`
- **Default zlib**: `purgeCompressionFor` falls back to `zlib`

So the answer to the task's question ("check whether Astrate needs to recognize/honour it, or whether zlib-only is the deliberate default; propose the change or note why not needed") is: **it is already fully recognized AND honoured (both zlib and plaintext), no change is needed.** This is a `[readonly]`-style conclusion — no code change required.

Let me verify the tests exist to confirm the coverage, and run the gate to give a clean report.

$ cd /root/astrate-mule && go test ./internal/engine/... ./internal/broker/... 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/engine	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/forward	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/stream	(cached)
ok  	github.com/astrate-platform/astrate/internal/engine/triggers	(cached)
ok  	github.com/astrate-platform/astrate/internal/broker	(cached)

The implementation is confirmed complete, tests pass. No code change is needed — the feature is already fully implemented and landed (commit `4e6db62`, "mule: issue-16"). This is a "note why not needed" outcome for the task.

My report:

Done: `purge_properties_compression_format` is already fully recognized and honoured (zlib and plaintext); no change needed
Files: (none changed)
Verified: `go test ./internal/engine/... ./internal/broker/...` -> pass
Unsure: nothing
Follow-ups: none
