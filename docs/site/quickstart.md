# Quickstart

Get Astrate running in 5 minutes with Docker Compose.

## Prerequisites

- Docker and Docker Compose
- Git

## Steps

### 1. Clone and start

```sh
git clone https://github.com/astrate/astrate.git
cd astrate
docker compose --profile full up -d --build
```

This starts:
- **Astrate** (API on `:8080`, MQTT on `:8883` and `:1883`)
- **TimescaleDB** (PostgreSQL + TimescaleDB)
- **Astarte Dashboard** at `http://localhost:4040`

### 2. Generate a login token

```sh
astartectl utils gen-jwt all-realm-apis -k deploy/devrealm/realm_private.pem
```

Paste the token into the Dashboard login at `http://localhost:4040`.

### 3. Verify health

```sh
curl http://localhost:8080/astrate/v1/health
# 200 OK

curl http://localhost:8080/astrate/v1/readiness
# 200 OK (database + broker healthy)
```

### 4. Register a device

```sh
# Generate a JWT for the test realm
TOKEN=$(astartectl utils gen-jwt "a_pa" -k deploy/devrealm/realm_private.pem)

# Register a device (use any 16-byte base64url ID)
curl -X POST http://localhost:8080/pairing/v1/test/agent/devices \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"data": {"hw_id": "dFj8kLmNpQrStUvWxYz0ab=="}}'
```

Response includes a `credentials_secret` -- this is shown only once.

### 5. Connect a device (Go SDK example)

```go
// Use the credentials_secret from step 4
// The Go SDK handles CSR, certificate, and MQTT connection
// See: https://github.com/astarte-platform/astarte-device-sdk-go
```

### 6. Explore the API

- **API Explorer:** [Swagger UI](swagger.md) — interactive API docs (no manual server needed)
- **OpenAPI specs:** `docs/api/` -- five YAML specs covering all REST surfaces
- **Metrics:** `http://localhost:8080/astrate/v1/metrics`

## Bare VPS deployment

```sh
# Install PostgreSQL + TimescaleDB
apt install postgresql-16 timescaledb

# Download the astrate binary
# (from releases or build from source)

# Create config
cat > /etc/astrate/astrate.toml << 'EOF'
database.dsn = "postgres://astrate:secret@localhost/astrate?sslmode=disable"

[realm]
name = "home"
jwt_public_key_file = "/etc/astrate/realm_public.pem"

[security]
master_key_file = "/etc/astrate/master.key"
EOF

# Run
astrate -config /etc/astrate/astrate.toml
```

Astrate self-migrates the schema on boot.

## Next steps

- [Architecture](architecture.md) -- understand the process model
- [Interface Schema](interface-schema.md) -- define your data contracts
- [Deployment](deployment.md) -- full configuration reference
- [Pairing and Security](pairing-and-security.md) -- production security model
