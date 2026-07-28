# Design: Clea Portal — appliance/device management (issue #29)

Status: proposal, not yet approved. Milestone v3.0 (Clea Portal).

## Scope boundary (read first)

Clea's four components split as: Portal (web front-end + billing, v3.0), Astarte (data hub,
already Astrate), Edgehog (fleet/device *management* — provisioning, OTA, remote actions,
v4.0), OS (v5.0+). Several things an operator would naturally ask for here —
**physical warehouse/stock state, box-QR device claim at unboxing, OTA, remote
diagnostics/actions, a mobile app** — are Edgehog territory per `.mule/milestones.md`, not
Portal. This doc keeps #29 to what Portal actually owns per Clea's own docs: **read/write
views over device metadata that already exists in Astrate** (device list, status, detail,
tags, customer assignment) plus a **registration workflow that is UI/API glue**, not device
provisioning itself (Astrate's `internal/pairing` already does provisioning; a device is
"claimed" by *associating* an already-registered Device ID with a customer, not by inventing
a second provisioning path).

The QR-on-the-box / warehouse-inventory / mobile-dashboard ideas below are real and worth
having, but I'm flagging them to `.mule/for-giulio.md` as a v3.0-vs-v4.0 boundary question
rather than baking them into #29 or silently editing #30-33 (they don't fit those either —
#30 is client/supplier relationships, not stock; #31-33 are roles/extensions/billing).

## What Portal's appliance view needs, concretely

Reference behaviour (Clea Portal docs): list of appliances with realtime status, name,
serial number, assigned customer, last-update timestamp, tags; per-appliance detail with tag
and customer-assignment editing; registration workflow capturing name, external serial
number, and the Astarte-style Device ID.

Astrate already has, in `internal/store/devices.go`: `Device{ID, RealmID, ...}` with a
`status` column (`registered`/`confirmed`/`inhibited`) and `internal/pairing` for the actual
credential/registration flow. It has **no** concept of: a friendly device name, an external
(operator-assigned) serial number, tags, or customer assignment. Realtime "status" in the
Portal sense (online/offline, last seen) needs to be derived from connection state the
broker already tracks (`internal/broker`) — needs checking whether that's already surfaced
anywhere in `internal/appengine`, or whether this issue needs to add it.

## Proposed data model

New table, one row per device *as Portal sees it* (keeps this additive — no changes to
`store.Device` or the wire-compatible Astarte device model):

```sql
CREATE TABLE portal_appliances (
    device_id       device_id NOT NULL,   -- FK to devices, matches deviceid.ID encoding
    realm_id        smallint NOT NULL,
    display_name    text NOT NULL,
    external_serial text,                  -- operator's own serial, not the Astarte Device ID
    customer_id     bigint,                 -- FK, nullable until claimed; references #30's
                                             -- client/supplier model once that lands
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (realm_id, device_id)
);

CREATE TABLE portal_appliance_tags (
    realm_id   smallint NOT NULL,
    device_id  device_id NOT NULL,
    tag        text NOT NULL,
    PRIMARY KEY (realm_id, device_id, tag),
    FOREIGN KEY (realm_id, device_id) REFERENCES portal_appliances (realm_id, device_id)
);
```

`customer_id` is left dangling (no FK enforcement) until #30 defines the client/supplier
model — noted as a caveat below, not resolved here.

## Proposed API surface

New package `internal/portal` (mirrors `internal/appengine`'s shape: HTTP handlers + a thin
store layer), realm-scoped like everything else:

- `GET /v3/portal/{realm}/appliances` — list, with `status`/`tag`/`customer_id` filters,
  backed by a join of `devices` (for connection status) + `portal_appliances` +
  `portal_appliance_tags`.
- `GET /v3/portal/{realm}/appliances/{device_id}` — detail.
- `PATCH /v3/portal/{realm}/appliances/{device_id}` — edit `display_name`, `tags`,
  `customer_id`.
- `POST /v3/portal/{realm}/appliances` — registration workflow: body is
  `{display_name, external_serial, device_id}`. This does **not** call `internal/pairing`
  (the device must already exist — either self-registered via Astarte's own pairing flow, or
  pre-registered by an operator through Housekeeping/Realm Management as today); this
  endpoint only creates the `portal_appliances` row. If `device_id` doesn't already exist in
  `devices`, return 404 rather than provisioning — provisioning is out of scope here.

Realtime status: reuse whatever the broker already exposes for connection state (needs the
investigation task below) rather than inventing a second live-status channel; if nothing
exists yet, the minimal version is a `last_connection`/`last_disconnection` timestamp pair
already likely present on `devices`, with "online" derived, not polled from the broker
directly.

## To investigate before implementation starts

1. Does `internal/appengine` (or `internal/broker`) already expose per-device
   online/offline + last-seen anywhere? If yes, reuse it; if no, that's its own small piece
   of work, possibly shared with #30/#31's dashboards.
2. Confirm `device_id` column type/encoding to match FK correctly (`deviceid.ID` — check its
   Postgres representation in `store.Device` before writing the migration).
3. Auth: these endpoints need *some* actor identity to authorize against, but #31 (users/
   roles) doesn't exist yet. Caveat below.

## Caveats / open decisions (not resolved by this doc)

- **Auth ordering**: this issue's endpoints are meaningless without an authenticated
  "operator" actor, which #31 defines. Either #29 ships behind realm-admin JWT auth only
  (today's model) as an interim, or #31 needs to land first. Recommend: interim realm-admin
  auth, revisit once #31 lands — but that's a call for Giulio, not assumed here.
- **`customer_id` FK**: dangling until #30 defines the client/supplier schema. Migration
  above intentionally omits the FK constraint so #29 isn't blocked on #30's design.
- **QR-claim / warehouse / mobile dashboard**: real operator needs, but Edgehog (v4.0)
  scope per milestones.md, not Portal. Flagged to `.mule/for-giulio.md`, not built here.
- **"Realtime" status**: depends on investigation item 1 above; may need a small addition to
  `internal/broker`/`internal/appengine` that's technically outside #29's stated scope but a
  hard dependency of it.
