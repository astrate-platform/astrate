# Pairing and Security

Astrate reproduces Astarte's three credential planes exactly, replacing the CFSSL sidecar with an embedded per-realm CA.

## Credential planes

| Plane | Purpose | Lifetime |
|---|---|---|
| **Realm JWTs** | Humans/services calling REST APIs | Per-token (configurable expiry) |
| **Credentials secret** | Per-device bearer for Pairing API calls | Long-lived, shown once at registration |
| **mTLS client certs** | Device-to-broker authentication | Configurable (default 30 days) |

## JWT validation

- Asymmetric keys only: RSA-2048+ or ECDSA P-256+.
- Algorithm allowlist: `RS256/384/ES256/ES384/ES512`. `none` and HMAC hard-rejected.
- Claims: `a_aea` (AppEngine), `a_rma` (Realm Management), `a_pa` (Pairing), `a_ha` (Housekeeping), `a_ch` (Channels/live stream).
- Each claim is a list of `"<verb-regex>::<path-regex>"` authorization strings, matched against method + path relative to realm base. Regexes are implicitly anchored (upstream parity).
- Multiple realm public keys allow zero-downtime rotation.

## Embedded per-realm CA

Replaces CFSSL:

- On realm creation: ECDSA P-256 CA key + self-signed cert (default 10 year lifetime), or import operator-provided pair.
- Private key encrypted at rest with AES-256-GCM under a master key (from env/file).
- Issues client certs: `Subject CN = <realm>/<device_id>`, 128-bit random serial, `KeyUsage = digitalSignature`, `ExtKeyUsage = clientAuth`.
- **All CSR subject/attributes are ignored** -- the CSR is proof of key possession only.
- Revocation: new credentials record new serial; broker rejects certs whose serial differs from latest (when `pairing.enforce_latest_cert` is enabled).

## Pairing flows

### Flow A -- Registration (operator)

```
POST /pairing/v1/<realm>/agent/devices
Authorization: Bearer <realm JWT with a_pa>
{ "data": { "hw_id": "<22-char base64url device ID>" } }
→ 201 { "data": { "credentials_secret": "<44-char base64>" } }
```

- Re-registering an existing device that has not yet requested credentials rotates the secret.
- After first credentials request, re-registration conflicts (422).
- Optional extension: `"initial_payload_format": "json"` for JSON-profile devices.

### Flow B -- Credentials (device)

```
POST /pairing/v1/<realm>/devices/<device_id>/protocols/astarte_mqtt_v1/credentials
Authorization: Bearer <credentials_secret>
{ "data": { "csr": "-----BEGIN CERTIFICATE REQUEST-----..." } }
→ 201 { "data": { "client_crt": "-----BEGIN CERTIFICATE-----..." } }
```

- Bcrypt constant-time comparison against stored hash.
- First successful call stamps `first_credentials_request` and flips status `registered -> confirmed`.

### Flow C -- Broker discovery

```
GET /pairing/v1/<realm>/devices/<device_id>
Authorization: Bearer <credentials_secret>
→ 200 { "data": {
    "status": "confirmed",
    "protocols": { "astarte_mqtt_v1": {
        "broker_url": "mqtts://<host>:8883",
        "ca_crt": "<realm CA PEM>"
    } } } }
```

### Full sequence

```
Agent                Astrate(pairing)          Device                 Astrate(broker)
  | POST agent/devices    |                       |                        |
  |------JWT(a_pa)------->|                       |                        |
  |<--credentials_secret--|   (secret delivered   |                        |
  |                       | out-of-band)          |                        |
  |                       |<--POST credentials----|  keygen + CSR          |
  |                       |---client_crt--------->|                        |
  |                       |<--GET device info-----|                        |
  |                       |---broker_url, ca_crt->|                        |
  |                       |                       |--CONNECT (mTLS, CN)--->|
  |                       |                       |<-CONNACK(session_present)|
  |                       |                       |-- introspection, subs, |
  |                       |                       |   emptyCache, data --->|
```

## Device status lifecycle

1. **registered** -- device registered, awaiting first credentials request.
2. **confirmed** -- first credentials request succeeded.
3. **inhibited** -- blocked from new credentials and connections (via `PATCH /appengine/v1/<realm>/devices/<id>` with `credentials_inhibited: true`).

## Rate limiting

- Pairing endpoints: per-IP and per-device token buckets.
- MQTT CONNECT storm damping: per-IP.
- AppEngine write endpoints: per-token.
