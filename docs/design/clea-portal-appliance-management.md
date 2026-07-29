# Design: Astrate Panel — appliance/device management (issue #29)

Status: proposal, not yet approved. Milestone v3.0 — reimplementation of the surface covered
by Clea Portal. Our own frontend is called **Astrate Panel**; "Clea Portal" / "Portal" below
refer only to the third-party product being reimplemented, never to our own thing.

## Scope boundary (read first)

Clea's four components split as: Portal (web front-end + billing, v3.0), Astarte (data hub,
already Astrate), Edgehog (fleet/device *management* — provisioning, OTA, remote actions,
v4.0), OS (v5.0+). Several things an operator would naturally ask for here —
**physical warehouse/stock state, box-QR device claim at unboxing, OTA, remote
diagnostics/actions, a mobile app** — are Edgehog territory per `.mule/milestones.md`, not
this issue. This doc keeps #29 to what Clea Portal's appliance surface actually owns per its
own docs, reimplemented as part of Astrate Panel: **read/write views over device metadata
that already exists in Astrate** (device list, status, detail, tags, customer assignment)
plus a **registration workflow that is UI/API glue**, not device provisioning itself
(Astrate's `internal/pairing` already does provisioning; a device is "claimed" by
*associating* an already-registered Device ID with a customer, not by inventing a second
provisioning path).

The QR-on-the-box / warehouse-inventory / mobile-dashboard ideas below are real and worth
having, but I'm flagging them to `.mule/for-giulio.md` as a v3.0-vs-v4.0 boundary question
rather than baking them into #29 or silently editing #30-33 (they don't fit those either —
#30 is client/supplier relationships, not stock; #31-33 are roles/extensions/billing).

## What the appliance view needs, concretely

Reference behaviour (Clea Portal docs): list of appliances with realtime status, name,
serial number, assigned customer, last-update timestamp, tags; per-appliance detail with tag
and customer-assignment editing; registration workflow capturing name, external serial
number, and the Astarte-style Device ID.

Astrate already has, in `internal/store/devices.go`: `Device{ID, RealmID, ...}` with a
`status` column (`registered`/`confirmed`/`inhibited`) and `internal/pairing` for the actual
credential/registration flow. It has **no** concept of: a friendly device name, an external
(operator-assigned) serial number, tags, or customer assignment.

**Investigated — connectivity status already exists.** `devices` has `connected bool`,
`last_connection`, `last_disconnection`, `last_seen_ip`, maintained by
`store.SetDeviceConnected`/`SetDeviceDisconnected` (`internal/store/devices.go:292-315`),
driven by the broker's connect/disconnect hooks. `store.DeviceStats` already aggregates
online counts per realm. So online/offline + last-seen needs **no new work** — the list/
detail endpoints just read these columns.

**But connectivity ≠ operational health, and that distinction matters here.** "Online" only
means the device has an open MQTT session; it says nothing about whether the appliance is
doing its job (a smart bed's sensor stuck, a coffee machine unable to brew). The appliance
status column should show a coarser "operational status," not raw connectivity, and that has
to be device-defined, not something Astrate infers — severity is domain-specific (a sensor
glitch on a bed might be a warning, a coffee machine that can't dispense is critical).
Proposal: ship a **standard Astarte interface** (e.g. `com.astrate.DeviceHealth`,
device-owned, properties-type) with `status` (`ok`/`warning`/`critical`) and `message`
(string), documented as an *optional convention* devices can implement and Astrate Panel
treats specially — same pattern Astarte itself uses for well-known interfaces. The appliance
list shows connectivity (from `devices.connected`) and, where the interface is implemented,
the self-reported health severity as a separate badge; devices that don't implement it just
show connectivity. This needs its own small piece of design (exact interface JSON, whether
severity levels are fixed or extensible) — noted as an open item below rather than speccing
the interface in full here. **Now split out as its own issue** (see below) since it's a
device-facing interface spec, not an Astrate Panel API.

## Proposed data model

New table, one row per device *as Astrate Panel sees it* (keeps this additive — no changes
to `store.Device` or the wire-compatible Astarte device model):

**Investigated — column type.** `devices.id` is plain `uuid` (`migrations/000002_metadata.up.sql`),
FK'd as `(realm_id, id)`; `deviceid.ID` (`pkg/deviceid`) is just the base64url *encoding* of
that UUID used on the wire/API, not a separate storage type. So the new tables below FK
against `uuid`, matching `devices.id` exactly — no custom type needed.

```sql
CREATE TABLE panel_appliances (
    device_id       uuid NOT NULL,          -- = devices.id
    realm_id        smallint NOT NULL,
    display_name    text NOT NULL,
    external_serial text,                  -- operator's own serial, not the Astarte Device ID
    customer_id     bigint,                 -- FK, nullable until claimed; references #30's
                                             -- client/supplier model once that lands
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (realm_id, device_id),
    FOREIGN KEY (realm_id, device_id) REFERENCES devices (realm_id, id) ON DELETE CASCADE
);

CREATE TABLE panel_appliance_tags (
    realm_id   smallint NOT NULL,
    device_id  uuid NOT NULL,
    tag        text NOT NULL,
    PRIMARY KEY (realm_id, device_id, tag),
    FOREIGN KEY (realm_id, device_id) REFERENCES panel_appliances (realm_id, device_id) ON DELETE CASCADE
);
```

`customer_id` is left dangling (no FK enforcement) until #30 defines the client/supplier
model — noted as a caveat below, not resolved here.

## Proposed API surface

New package `internal/panel` (mirrors `internal/appengine`'s shape: HTTP handlers + a thin
store layer), realm-scoped like everything else:

- `GET /v3/panel/{realm}/appliances` — list, with `status`/`tag`/`customer_id` filters,
  backed by a join of `devices` (for connection status) + `panel_appliances` +
  `panel_appliance_tags`.
- `GET /v3/panel/{realm}/appliances/{device_id}` — detail.
- `PATCH /v3/panel/{realm}/appliances/{device_id}` — edit `display_name`, `tags`,
  `customer_id`.
- `POST /v3/panel/{realm}/appliances` — registration workflow: body is
  `{display_name, external_serial, device_id}`. This does **not** call `internal/pairing`
  (the device must already exist — either self-registered via Astarte's own pairing flow, or
  pre-registered by an operator through Housekeeping/Realm Management as today); this
  endpoint only creates the `panel_appliances` row. If `device_id` doesn't already exist in
  `devices`, return 404 rather than provisioning — provisioning is out of scope here.

## Investigated

1. ~~Does connectivity/last-seen already exist?~~ **Yes** — `devices.connected` +
   `last_connection`/`last_disconnection`, already maintained by the broker. No new work.
2. ~~`device_id` column type?~~ **`uuid`**, same as `devices.id`. FKs above use it directly.

## Still open

3. Auth: these endpoints need *some* actor identity to authorize against, but #31 (users/
   roles) doesn't exist yet. Caveat below.
4. `com.astrate.DeviceHealth` interface spec — split into its own issue (see project
   tracker), not resolved here.

## Caveats / open decisions (not resolved by this doc)

- **Auth ordering**: this issue's endpoints are meaningless without an authenticated
  "operator" actor, which #31 defines. Either #29 ships behind realm-admin JWT auth only
  (today's model) as an interim, or #31 needs to land first. Recommend: interim realm-admin
  auth, revisit once #31 lands — but that's a call for Giulio, not assumed here.
- **`customer_id` FK**: dangling until #30 defines the client/supplier schema. Migration
  above intentionally omits the FK constraint so #29 isn't blocked on #30's design.
- **QR-claim / warehouse / mobile dashboard**: real operator needs, but Edgehog (v4.0)
  scope per milestones.md, not this issue. Flagged to `.mule/for-giulio.md`, not built here.
- **"Realtime" status**: resolved — connectivity already exists in `devices`, no new plumbing
  needed for that part.
- **Operational health vs. connectivity**: split into its own issue (item 4 above).
- **Which frontend Astrate Panel actually is** (fork of the Astarte dashboard vs. new build)
  is out of scope for this doc — see the separate device/appliance-lifecycle issue and the
  Astrate Panel proof-of-concept issue for that thread.
