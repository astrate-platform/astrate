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

**T3 — Full emulator loop (real ROM)**

ROM is **not** in the repo. Pass an absolute path you already own, e.g.:

```sh
ROM="/path/to/Pokemon Red.gb"   # do not commit this file
DEVICE_ID=<registered-device-id>

cd examples/pokemon-agent/emulator-agent
python -m emulator_agent.main \
  --rom "$ROM" \
  --insecure \
  --astrate-url http://localhost:8080 \
  --realm test \
  --device-id "$DEVICE_ID"
# optional: --fps 60 (default) | --fps 0 for uncapped (high CPU)
# default: --skip-intro (mash A/START past title); --no-skip-intro to disable
```

Or via Makefile from `examples/pokemon-agent/`:

```sh
make run-emulator-rom DEVICE_ID="$DEVICE_ID" ROM="$ROM"
```

Expected (verified 2026-07-29 on pyboy + insecure Astrate; intro skip fixed session 10):
- Agent loads ROM, connects MQTT `:1883`, publishes introspection
- With default `--skip-intro`, logs show A/START + direction probes, then:
  - `Intro: WRAM spawn preload at (3, 6) …` (Oak preloads Red's House — **not free yet**)
  - `Intro skip complete — player moved (3, 6) → (x, y)` when a real tile step lands
  - Or timeout ~180 s if movement never happens
- Completing on non-zero WRAM alone is wrong (coords appear mid-Oak intro)
- After intro skip, directional ControlCommands with `holdFrames≥16` change
  `playerX`/`playerY` (Red's House 2F stairs ≈ increase X, decrease Y → (7,1))
- Cold boot has empty party → no `PartyStatus` until in-game
- Process CPU stays modest at default `--fps 60` (uncapped was ~full core)
- ControlCommand from AppEngine is queued on MQTT thread and applied on the
  main loop (no `pyboy.tick` from the paho callback)

**T4 — LLM Orchestrator end-to-end (opencode / Big Pickle preferred)**

Prerequisites: Astrate + T3 emulator running (device registered, intro skip done so
GameState is non-zero). `opencode` on PATH for the free-model path.

### Mint App JWT (stream + publish)

Astrate’s live stream (`GET /astrate/v1/{realm}/socket`) authorizes with
`ClaimChannels` using **REST** grammar (`METHOD::path`), **not** Phoenix
`JOIN`/`WATCH`. `astartectl utils gen-jwt appengine channels` mints
`a_ch: ["JOIN::.*","WATCH::.*"]` → **403** on the socket.

Mint with REST-style claims (PyJWT + realm private key):

```sh
python3 - <<'PY'
import jwt, time
from pathlib import Path
key = Path("deploy/devrealm/realm_private.pem").read_text()
tok = jwt.encode(
    {"a_aea": [".*::.*"], "a_ch": [".*::.*"], "exp": int(time.time()) + 8*3600},
    key, algorithm="RS256",
)
Path("/tmp/pokemon-app-token.txt").write_text(tok)
print("ok", len(tok))
PY
# Auth probe without WS upgrade: 426 = JWT accepted (Upgrade Required)
curl -s -o /dev/null -w '%{http_code}\n' \
  -H "Authorization: Bearer $(cat /tmp/pokemon-app-token.txt)" \
  http://localhost:8080/astrate/v1/test/socket
```

### Run orchestrator — opencode / Big Pickle (no API key)

```sh
cd examples/pokemon-agent/llm-orchestrator
pip install -r requirements.txt

export POKEMON_ASTRATE_URL=http://localhost:8080
export POKEMON_ASTRATE_REALM=test
export POKEMON_ASTRATE_DEVICE_ID=<device-id>   # e.g. cat /tmp/pokemon-device-id.txt
export POKEMON_ASTRATE_APP_TOKEN="$(cat /tmp/pokemon-app-token.txt)"

export POKEMON_OPENAI_MODEL=opencode/big-pickle
export POKEMON_LLM_BACKEND=opencode   # or auto (auto selects opencode for opencode/* models)
export POKEMON_LLM_TIMEOUT_SECONDS=180 # 60s often times out on cold/slow big-pickle turns
export POKEMON_LLM_MAX_RETRIES=2
# no POKEMON_OPENAI_API_KEY needed for free opencode models

python3 -m llm_orchestrator.main
```

Expected (verified 2026-07-29 with `opencode/big-pickle` + live T3):
- Log: `LLM backend=opencode model=opencode/big-pickle (no API key)`
- WebSocket connect for GameState + PartyStatus
- Log: `LLM → <BUTTON> ×N (seq=…)` on each GameState turn
- ControlCommand POST → HTTP 200; samples on AppEngine
- Emulator log: non-`(local)` `Input: … (seq=N)` lines (intro inputs are local only)
- After a correct intro skip: GameState `playerX`/`playerY` change under LLM
  control (overworld holdFrames floored to 16 for directions)

### T4b — Leave Red's House (light guidance)

Same stack as T4, but use the deterministic early-game guide (still via Astrate
ControlCommand — no direct pyboy control from the orchestrator):

```sh
export POKEMON_GUIDANCE=light          # or auto (light when it has a path, else LLM)
export POKEMON_TURN_COOLDOWN_SECONDS=0.4
# other POKEMON_ASTRATE_* same as T4
python3 -m llm_orchestrator.main
```

Expected (verified 2026-07-29):
- After intro: path on Red's House 2F toward stairs (7,1)
- Map warp: `mapId=38` → `37` (Red's House 1F)
- Exit south door → `mapId=0` Pallet Town with non-zero coords
  (e.g. `(5,6)`)
- Orchestrator log: `LIGHT → RIGHT/UP/DOWN/LEFT ×16 (seq=N) @ map=…`
- Do **not** inject a high `sequenceId` via curl mid-run: the emulator rejects
  later commands with `sequenceId <= last` (MQTT redelivery dedup)

`POKEMON_GUIDANCE=llm` remains the default. Prefer `light` when opencode is too
slow for multi-step navigation smoke.

### OpenAI-compatible HTTP alternative

```sh
export POKEMON_LLM_BACKEND=openai   # or auto with a non-opencode model
export POKEMON_OPENAI_API_BASE=https://api.openai.com/v1
export POKEMON_OPENAI_API_KEY=<key>
export POKEMON_OPENAI_MODEL=gpt-4o
export POKEMON_LLM_TIMEOUT_SECONDS=15
python3 -m llm_orchestrator.main
```

Do **not** commit API keys. Prefer opencode free models for local smoke unless you
need a paid HTTP provider.
