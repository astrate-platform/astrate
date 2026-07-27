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
