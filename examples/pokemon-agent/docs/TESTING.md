# Testing Guide

**T0 — Python unit tests (no dependencies)**
```sh
cd examples/pokemon-agent/emulator-agent && pip install pytest && pytest tests/
cd examples/pokemon-agent/llm-orchestrator && pip install pytest pydantic pydantic-settings && pytest tests/
```
Expected: all tests pass. No ROM, no Astrate, no LLM needed.

**T1 — Interface schema validation + install**

Prerequisites: a running Astrate with a realm and realm private key.

Local dev (recommended for smoke):

```sh
# DB only (docker compose full profile currently fails docker build if
# .dockerignore excludes docs/; run the binary on the host instead)
docker compose up -d timescaledb

export ASTRATE_DATABASE_DSN='postgres://astrate:astrate@127.0.0.1:5432/astrate?sslmode=disable'
export ASTRATE_MQTT_INSECURE_DEV_MODE=true
export ASTRATE_MQTT_SESSION_STORE_PATH=/tmp/astrate-pokemon-sessions.db
export ASTRATE_MASTER_KEY='00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff'
export ASTRATE_REALM_NAME=test
export ASTRATE_REALM_JWT_PUBLIC_KEY_FILE=deploy/devrealm/realm_public.pem
go build -o /tmp/astrate ./cmd/astrate && /tmp/astrate
# readiness: curl -s localhost:8080/astrate/v1/readiness
```

```sh
# Validate JSON syntax
for f in examples/pokemon-agent/astrate-interfaces/*.json; do
  python3 -m json.tool "$f" > /dev/null && echo "OK: $f"
done

# Install into the auto-provisioned "test" realm (deploy/devrealm keys)
ASTRATE_URL=http://localhost:8080
REALM=test
REALM_KEY=deploy/devrealm/realm_private.pem

for f in examples/pokemon-agent/astrate-interfaces/*.json; do
  astartectl realm-management interfaces install "$f" \
    -u "$ASTRATE_URL" -k "$REALM_KEY" -r "$REALM"
done
astartectl realm-management interfaces list -u "$ASTRATE_URL" -k "$REALM_KEY" -r "$REALM"
```
Expected: all three `org.pokemon.emulator.*` interfaces listed.

**T2 — Emulator Agent stub mode**

With `mqtt.insecure_dev_mode` (compose/local default), devices use **plaintext MQTT on :1883**.
The agent flag `--insecure` skips mTLS and sets the client ID to `<realm>/<device_id>`
(required by Astrate's plaintext auth). Device still must be **registered**.

```sh
ASTRATE_URL=http://localhost:8080
REALM=test
REALM_KEY=deploy/devrealm/realm_private.pem
DEVICE_ID=$(astartectl utils device-id generate-random)
astartectl pairing agent register "$DEVICE_ID" \
  -u "$ASTRATE_URL" -k "$REALM_KEY" -r "$REALM"

cd examples/pokemon-agent/emulator-agent
pip install -r requirements.txt
python -m emulator_agent.main \
  --stub \
  --insecure \
  --astrate-url "$ASTRATE_URL" \
  --realm "$REALM" \
  --device-id "$DEVICE_ID"
```

Verify GameState (paths corrected in P0):

```sh
# Object datastream samples
astartectl appengine devices get-samples "$DEVICE_ID" \
  org.pokemon.emulator.GameState /state \
  -u "$ASTRATE_URL" -k "$REALM_KEY" -r "$REALM" --count 3

# REST
TOKEN=$(astartectl utils gen-jwt appengine -k "$REALM_KEY")
curl -sS -H "Authorization: Bearer $TOKEN" \
  "$ASTRATE_URL/appengine/v1/$REALM/devices/$DEVICE_ID/interfaces/org.pokemon.emulator.GameState/state?limit=2"

# Live stream (JWT needs channels claim a_ch)
# GET /astrate/v1/$REALM/socket?device_id=$DEVICE_ID&interface=org.pokemon.emulator.GameState
```

Expected:
- Device shows introspection for all three pokemon interfaces
- Samples with `mapName: "Pallet Town"`, moving `playerX`
- PartyStatus `/0/name` = `Pikachu`, `/0/currentHp` and `/0/maxHp` = `20`

Production mTLS (no `--insecure`): pass `--cert` / `--key` / `--ca` from Pairing
credentials handshake; MQTT defaults to port 8883.

**T3 — Full emulator loop**
```sh
python -m emulator_agent.main \
  --rom /path/to/pokemon_red.gb \
  --insecure \   # drop when using real mTLS
  --astrate-url http://localhost:8080 \
  --realm test \
  --device-id <device-id>
```
Expected: telemetry flowing, visible via AppEngine as in T2.

**T4 — LLM Orchestrator with stub**
```sh
cd examples/pokemon-agent/llm-orchestrator
pip install -r requirements.txt
export POKEMON_ASTRATE_URL=http://localhost:8080
export POKEMON_ASTRATE_REALM=test
export POKEMON_ASTRATE_DEVICE_ID=<device-id>
export POKEMON_ASTRATE_APP_TOKEN=<jwt with a_ch + appengine>
export POKEMON_OPENAI_API_KEY=<key>
python -m llm_orchestrator.main
```
Expected: ControlCommand events visible in Astrate; if emulator agent is also running, character moves.
