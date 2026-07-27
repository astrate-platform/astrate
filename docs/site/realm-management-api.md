# Realm Management API

The Realm Management API manages interface (schema) and trigger definitions per realm. It is the operator's interface to the data model.

**Base path:** `/realmmanagement/v1/<realm>/`
**Authentication:** JWT with `a_rma` claim.

## Interfaces

### Install interface

```
POST /realmmanagement/v1/<realm>/interfaces
{ "data": <interface JSON definition> }
```

Parses and validates the interface definition (see [Interface Schema](interface-schema.md)), stores it with normalized endpoint rows, and notifies the engine to rebuild its compiled cache.

### List interfaces

```
GET /realmmanagement/v1/<realm>/interfaces
```

### Get interface

```
GET /realmmanagement/v1/<realm>/interfaces/<name>/<major_version>
```

### Update interface (minor bump)

```
PUT /realmmanagement/v1/<realm>/interfaces/<name>/<major_version>
{ "data": <updated interface JSON> }
```

Only minor-version bumps are allowed. The update must be **additive** -- new mappings only, no mutation of existing mapping attributes, same type/ownership/aggregation. Enforced by `CheckMinorUpgrade` (parity with upstream).

### Delete interface

```
DELETE /realmmanagement/v1/<realm>/interfaces/<name>/<major_version>
```

Only allowed when:
- Major version is 0, **or**
- The interface is not in any device's introspection.

### Get interface version

```
GET /realmmanagement/v1/<realm>/version
```

Reports the emulated upstream API compatibility level (currently `1.2.2`), not Astrate's own version. Used by the Dashboard to feature-gate the Policies UI.

## Triggers

### Install trigger

```
POST /realmmanagement/v1/<realm>/triggers
{ "data": <trigger JSON definition> }
```

Validates the trigger definition against the Astarte trigger schema (simple triggers + action).

### List triggers

```
GET /realmmanagement/v1/<realm>/triggers
```

### Get trigger

```
GET /realmmanagement/v1/<realm>/triggers/<name>
```

### Update trigger

```
PUT /realmmanagement/v1/<realm>/triggers/<name>
{ "data": <updated trigger JSON> }
```

### Delete trigger

```
DELETE /realmmanagement/v1/<realm>/triggers/<name>
```

## Auth configuration

### Rotate JWT keys

```
PUT /realmmanagement/v1/<realm>/config/auth
{ "data": { "jwt_public_key_pem": "<PEM string>" } }
```

Replaces the realm's JWT public key. Multiple keys are accepted for zero-downtime rotation.

## Notifications

Every mutation (interface install/update/delete, trigger CRUD) issues a Postgres `NOTIFY astrate_interfaces`. The engine listens and rebuilds the affected realm's compiled cache. In a single process this is a belt-and-suspenders measure; it exists so an optional hot-standby instance stays coherent.
