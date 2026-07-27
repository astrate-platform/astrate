# Troubleshooting

Common failure modes and how to fix them.

## Certificate errors

### Wrong CA / certificate not trusted

**Symptom:** Device fails to connect to the broker. Logs show TLS handshake errors.

**Cause:** The broker's TLS cert is not issued from a CA the device trusts, or the per-realm CA is wrong.

**Fix:**
1. The device trusts the per-realm CA returned by `GET /pairing/v1/<realm>/devices/<id>` (the `ca_crt` field). Verify this CA matches what you expect.
2. The broker's own TLS cert (`mqtt.tls_cert_file`) must be issued from a CA your fleet trusts. It does **not** need to be the realm CA — it just needs to chain to a trusted root.
3. If you suspect CA compromise, re-key the realm CA via the Housekeeping API (delete and recreate the realm).

### Expired client certificate

**Symptom:** Device was connecting fine, now gets rejected. Logs show "certificate has expired".

**Cause:** Client certs expire after `pairing.cert_ttl` (default 30 days).

**Fix:**
1. The SDK should call `POST /pairing/v1/<realm>/devices/<id>/protocols/astarte_mqtt_v1/credentials/verify` on boot. If it returns `"valid": false, "cause": "EXPIRED"`, the SDK re-CSRs automatically.
2. If the SDK doesn't re-CSR, check that `pairing.cert_ttl` is not set too short for your device reboot frequency.
3. To force fleet-wide rotation, shorten `cert_ttl` and ensure `enforce_latest_cert` is on.

### Serial mismatch (`enforce_latest_cert`)

**Symptom:** Device is rejected even though its cert is not expired. Logs show serial mismatch.

**Cause:** `pairing.enforce_latest_cert = true` and the device presented an older certificate.

**Fix:**
1. The device must re-CSR to get a cert with the latest serial. The SDK's `verify` call should trigger this.
2. If devices are offline and cannot re-CSR, disable `enforce_latest_cert` temporarily.
3. This is stricter than upstream Astarte's CRL-less default — it's a security feature, not a bug.

## Database connection refused

**Symptom:** Astrate fails to start or reports "connection refused" to PostgreSQL.

**Fix:**
1. Verify PostgreSQL is running: `pg_isready -h 127.0.0.1 -p 5432`
2. Check the DSN: `database.dsn` or `ASTRATE_DATABASE_DSN` must include `sslmode=disable` for local dev, or `sslmode=require`/`sslmode=verify-full` for production.
3. In Docker Compose, ensure `timescaledb` is healthy before `astrate` starts (the `depends_on` condition handles this).
4. For bare VPS: `systemctl status postgresql` and check `pg_hba.conf` allows the connection.

## Migration failures

**Symptom:** Astrate exits on boot with a migration error.

**Cause:** Schema migrations (embedded via `golang-migrate`) failed against the database.

**Fix:**
1. Check that the TimescaleDB extension is enabled: `SELECT * FROM pg_extension WHERE extname = 'timescaledb';`
2. Ensure the database user has CREATE/DROP privileges on the schema.
3. If a previous migration partially ran, check `schema_migrations` table and manually fix the state if needed.
4. Astrate runs migrations automatically on boot — do not run them manually unless you know what you're doing.

## MQTT connection refused

**Symptom:** Device cannot connect to port `:8883`. Connection refused.

**Fix:**
1. Check `mqtt.insecure_dev_mode`: if `true`, the broker only binds `dev_addr` (`:1883`), not `:8883`. Connect to `:1883` instead.
2. In production (`insecure_dev_mode = false`), verify `mqtt.tls_cert_file` and `mqtt.tls_key_file` are set and the files exist.
3. Check firewall rules: `:8883` (or `:1883` in dev) must be open.
4. In Docker Compose, verify the port mapping: `127.0.0.1:8883:8883`.

## Master key lost

**Symptom:** You cannot decrypt realm CA private keys. New device registrations work but existing devices cannot get new certs.

**Fix:**
1. The master key is used to AES-256-GCM-seal realm CA private keys in the database.
2. **Losing it does not brick existing devices.** Devices keep their credentials secret and will re-pair automatically at their next credential rotation.
3. To recover: generate a new master key, delete and recreate affected realms via the Housekeeping API (this mints fresh CAs sealed under the new key).
4. **Prevention:** back up the master key separately from the database. The two together decrypt everything.

## Dashboard 401 / 403

**Symptom:** Astarte Dashboard shows "401 Unauthorized" or "403 Forbidden".

**Fix:**
1. **401:** The JWT token is missing, expired, or signed with an unknown key. Generate a fresh token:
   ```sh
   astartectl utils gen-jwt all-realm-apis -k <realm_private.pem>
   ```
2. **403:** The token is valid but lacks the required `a_aea` claim for AppEngine access.
3. **CORS errors:** Add the Dashboard origin to `http.cors_allowed_origins` (e.g. `"http://localhost:4040"`).
4. **Realm mismatch:** Ensure the token is signed for the same realm you're querying.

## Rate limit 429 (Too Many Requests)

**Symptom:** Pairing or credentials endpoint returns `429`.

**Cause:** Token-bucket rate limit exceeded (`pairing.register_rate` / `pairing.credentials_rate`).

**Fix:**
1. Check `pairing.register_rate` and `pairing.register_burst` (defaults: 5 req/s, burst 10).
2. Check `pairing.credentials_rate` and `pairing.credentials_burst` (same defaults).
3. Rate limits are per-IP and per-device. If a fleet of devices is rebooting simultaneously, they may hit the burst limit.
4. Increase the rate/burst values if your fleet size justifies it.
5. Check logs for the specific rate-limit rejection to identify which endpoint and IP.

## Payload rejected

**Symptom:** Device sends data but Astrate rejects it (logs show `device_error` / `unexpected_value` or similar).

**Common causes:**
1. **Introspection mismatch:** The device's introspection doesn't declare the interface+major it's publishing on. Fix: update the device's introspection.
2. **Type mismatch:** The payload value doesn't match the endpoint's declared type (e.g. sending a string where an integer is expected).
3. **Payload too large:** Exceeds `engine.max_payload_bytes` (default 64 KB).
4. **BSON/JSON detection failed:** Payload is neither valid BSON nor valid JSON starting with `{`.
5. **Ownership violation:** Device publishing on a `server`-owned interface, or AppEngine publishing on a `device`-owned interface.

Check the specific reject reason in the logs — each rejection includes the device, interface, path, and reason.

## Device not connecting

**Symptom:** Device was working, now shows as disconnected. No data arriving.

**Diagnostic steps:**
1. Check device status: `GET /appengine/v1/<realm>/devices/<id>` — is `connected: false`?
2. Check `last_connection` and `last_disconnection` timestamps.
3. Verify the device's client cert hasn't expired (`pairing.cert_ttl`).
4. Check broker logs for connection attempts and rejection reasons.
5. If using `insecure_dev_mode`, ensure the device connects to `:1883`, not `:8883`.
6. For production: verify the broker TLS cert is valid and the device can reach the host/port.

## General debugging

- **Enable debug logging:** set `log.level = "debug"` in the TOML file or `ASTRATE_LOG_LEVEL=debug`.
- **Health check:** `curl localhost:8080/astrate/v1/health` — returns 200 if the process is alive.
- **Readiness check:** `curl localhost:8080/astrate/v1/readiness` — returns 200 if DB + broker are reachable.
- **Metrics:** `curl localhost:8080/astrate/v1/metrics` — Prometheus scrape endpoint with ingest rate, reject counts, shard depth, DB pool stats.

See [Deployment](deployment.md) for deployment profiles and backup.

## See also

- [Configuration Reference](configuration-reference.md) — all TOML keys, env overrides, and defaults
- [Operations](operations.md) — configuration, backups, CA re-keying, and runtime management
- [Observability](observability.md) — health probes, Prometheus metrics, and structured logging
- [Pairing and Security](pairing-and-security.md) — credential planes and certificate lifecycle
