# realm-config-datastream-retention

Found by `.mule/recipes/astarte-upstream.md` (run 2026-08-03) against the
v1.2.2 → v1.4.0-rc.3 gap.

Upstream (astarte-platform/astarte) added realm-scoped datastream retention:

- **Realm Management GET** `/realmmanagement/v1/<realm>/config/datastream_maximum_storage_retention`
  (`apps/astarte_realm_management/lib/astarte_realm_management_web/router.ex`,
  present since v1.3.0; read-only — updates go through housekeeping).
  Astrate's `internal/realm/http.go` registers `/config/auth` and
  `/config/device_registration_limit` but not this route → 404.
- **Housekeeping realm create/get body** now carries `datastream_maximum_storage_retention`
  (`apps/astarte_housekeeping/lib/astarte_housekeeping/realms/realm.ex`,
  `apps/astarte_housekeeping/lib/astarte_housekeeping_web/views/realm_view.ex`).
  Astrate's `realmBody` (`internal/housekeeping/http.go:42`) has only
  `realm_name`, `jwt_public_key_pem`, `device_registration_limit`; a client
  setting retention would have it silently dropped (JSON unknown fields ignored).
- **v1.4.0-rc.1** added `HOUSEKEEPING_DEFAULT_DATASTREAM_MAXIMUM_STORAGE_RETENTION`
  (instance-wide default, seconds), injected via
  `apps/astarte_housekeeping/lib/astarte_housekeeping/realms/realms.ex
  maybe_inject_default_retention`.
- Enforcement is per-realm: the DUP reads the realm's
  `datastream_maximum_storage_retention` (stored in the realm `kv_store`
  `realm_config` group) as a datastream TTL
  (`apps/astarte_data_updater_plant/.../data_updater/impl.ex`,
  `time_based_actions.ex reload_datastream_maximum_storage_retention_on_expiry`).
  Astrate has only global retention (`internal/config/config.go` `StorageConfig.Retention`,
  applied by `internal/store/store.go` `ApplyGlobalRetention`).

Suggested split (tick-sized):
1. Wire surface: RM `GET /config/datastream_maximum_storage_retention` + accept/echo the
   field on the housekeeping create/get body.
2. Storage + enforcement: per-realm retention column/table, applied on ingest like
   `ApplyGlobalRetention`.

Evidence recorded against upstream:
- `router.ex@v1.3.0` and `@v1.4.0-rc.3` both list the route; `router.ex@v1.2.2` has no
  `/config/*` block at all (whole block is v1.3.0+).
- Release notes v1.3.0 and v1.4.0-rc.1 (both claims now diff-checked).
