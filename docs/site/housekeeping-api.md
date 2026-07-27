# Housekeeping API

The Housekeeping API manages realm lifecycle at the instance level. It is the top-level administrative surface.

**Base path:** `/housekeeping/v1/`
**Authentication:** JWT with `a_ha` claim (instance-admin keys, separate from realm keys).

## Realms

### Create realm

```
POST /housekeeping/v1/realms
{
  "data": {
    "realm_name": "<name>",
    "jwt_public_key_pem": "<PEM string>",
    "device_registration_limit": <optional integer>
  }
}
```

Creates the realm in a single transaction:
1. Insert the `realms` row with JWT public keys.
2. Generate (or import) the per-realm CA key + certificate.
3. Seal the CA private key with AES-256-GCM under the master key.

The broker reloads its CA trust automatically -- new realms accept device connections without a restart.

Cassandra-specific fields (`replication_class`, `replication_factor`, etc.) from upstream are accepted but ignored.

### List realms

```
GET /housekeeping/v1/realms
```

### Get realm

```
GET /housekeeping/v1/realms/<realm_name>
```

### Delete realm

```
DELETE /housekeeping/v1/realms/<realm_name>
```

Transactional cascade: removes all interfaces, devices, properties, datastreams, groups, and triggers. The realm's CA key material is deleted.

## Auto-provisioning

For single-realm installs, set the `[realm]` config block (`name` + `jwt_public_key`/`_file`) and Astrate creates the realm on first boot (no-op if it already exists).

```toml
[realm]
name = "home"
jwt_public_key_file = "/etc/astrate/realm_public.pem"
```

## Authentication

Housekeeping uses instance-level admin keys (`a_ha` claim), not realm-level keys. This separation ensures that realm operators cannot modify other realms.

The admin JWT public keys are configured at the instance level:

```toml
housekeeping.jwt_public_keys = ["<PEM string>"]
# or
housekeeping.jwt_public_key_files = ["/etc/astrate/admin_public.pem"]
```
