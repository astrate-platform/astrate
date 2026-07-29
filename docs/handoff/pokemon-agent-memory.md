# Pokemon Agent — Session Memory (2026-07-29, session 9)

## What happened this session

### Goal
**Finish T4** — LLM orchestrator end-to-end with opencode/Big Pickle against live
Astrate + T3 emulator.

### Results — T4 **PASS**

Stack already healthy from session 8:
- Astrate `/tmp/astrate` + TimescaleDB; device `HcLR8OkjSKat9ZAf8UolJg` connected
- Emulator PID was live (MQTT ESTABLISHED, Red's House 2F after intro skip)
- JWT `/tmp/pokemon-app-token.txt` still valid (`a_aea` + `a_ch` = `.*::.*`)

Orchestrator run:
```sh
POKEMON_OPENAI_MODEL=opencode/big-pickle
POKEMON_LLM_BACKEND=opencode
POKEMON_LLM_TIMEOUT_SECONDS=90
# no API key
python3 -m llm_orchestrator.main
```

Live evidence (2026-07-29 ~05:13–05:14):
- `LLM backend=opencode model=opencode/big-pickle (no API key)`
- WS connected GameState + PartyStatus
- **8 turns**: `LLM → DOWN|LEFT|A|RIGHT|DOWN|UP|DOWN|LEFT` with seq 1–8
- ControlCommand POST → **HTTP 200** each time
- Emulator: non-local `Input: … (seq=N)` for all 8 (not only intro `(local)`)
- AppEngine ControlCommand samples show seq 1–8

Coords stayed at `(3,6)` map 38 with `stasis=true` for the short run (Red's House
2F layout / short holds may block visible movement). **Bus path is closed**;
meaningful navigation is a later polish item, not a T4 gate.

### Code / docs committed this session
Under `examples/pokemon-agent/` (+ handoff + `docs/site/pokemon-agent.md` via sync):
1. `llm-orchestrator` opencode backend (from session 8, now E2E-proven) + unit tests
2. `tests/test_llm_engine.py` — backend auto-select, NDJSON extract, JSON parse
3. `docs/TESTING.md` T4 rewritten: JWT mint recipe + opencode + OpenAI alternative
4. `docs/DESIGN.md` §3.4 / §5.5 backend note; site sync
5. `llm-orchestrator/README.md` real content (was broken single-line)

### Branch
`feat/pokemon-agent`

Unrelated WIP still on the working tree (pipelines, broker ACL, triggers, go.mod).
Do **not** commit those with pokemon-agent work.

### Prior sessions
- Session 1: architecture, module tree, 10 ADRs
- Session 2 (P0): App API stream/publish paths
- Session 3 (P1/P2): unit tests; WRAM dialog `$CF4B`, maxHP@+34
- Session 4 (P3): live smoke T1+T2 stub + `--insecure`
- Session 5 (P4): ROM T3; 60 fps; dialog 0x00; stasis 15s
- Session 6 (P5): MkDocs Examples → Pokémon Agent nav + `make sync`
- Session 7: command queue + intro auto-press
- Session 8: T4 partial (WS+JWT; Ollama failed; opencode code uncommitted)
- **Session 9: T4 closed with Big Pickle**

---

## What's NOT done yet (next session scope)

### Optional polish / next phases
1. **Overworld movement smoke** — leave Red's House / get coords to change under
   LLM control (stairs DOWN from 2F; longer holds; clear stasis).
2. Fill remaining stub unit tests (`action_translator` is a pass; `context_builder`
   still thin; richer `state_decoder` fixtures on emulator side).
3. Mid-dialog `$CF4B` validation on real ROM after intro / with actual text boxes.
4. Endurance / loop-detection demo (DESIGN P4–P5 table is aspirational).
5. Root `.dockerignore` `docs` → compose full image (out of scope unless asked).
6. Do **not** re-introduce Ollama unless user asks.

### T4 is done — do not re-open unless regression.

## Risks and known issues
1. ~~WebSocket path~~ **P0**. ~~Dialog/maxHP WRAM~~ **P2**. ~~Smoke T1/T2~~ **P3**.
   ~~ROM T3~~ **P4**. ~~MkDocs nav~~ **P5**. ~~Command queue~~ **done**.
   ~~Intro auto-press~~ **done**. ~~T4 LLM→ControlCommand~~ **done** (session 9).
2. `paho-mqtt` 2.x `CallbackAPIVersion.VERSION2` required.
3. ~~pyboy headless uncapped CPU~~ mitigated by default `--fps 60`.
4. **App JWT for stream must use REST-style `a_ch` (`.*::.*`), not JOIN/WATCH.**
5. Compose full build: `.dockerignore` has bare `docs` → missing swagger embed path.
6. `astartectl appengine devices data-snapshot` panics on object datastream; use
   `get-samples` or REST `/…/GameState/state`.
7. Unrelated WIP on same branch — leave alone.
8. Cold-boot mapId 0 → “Pallet Town” is WRAM-zero coincidence until intro skip finishes.
9. Site copy `docs/site/pokemon-agent.md` goes stale if DESIGN changes without `make sync`.
10. Intro-skip can false-complete if dialog buffer is non-empty on title; timeout backstop.
11. Emulator can orphan with MQTT CLOSED + high CPU; kill and restart; always log to a file.
12. Default LLM timeout 5 s is too low for `opencode run` cold start — use ≥60.
13. Big Pickle may refuse “play Pokémon” wording; engine frames as IoT fixture (keep that).
14. opencode backend is subprocess-per-turn — slow vs HTTP; fine for smoke, not production.
15. Short T4 run left player stuck at (3,6) stasis on Red's House 2F — path works;
    navigation needs longer/more purposeful play.

## Recipes (local smoke)

### Astrate
```sh
docker compose up -d timescaledb
export ASTRATE_DATABASE_DSN='postgres://astrate:astrate@127.0.0.1:5432/astrate?sslmode=disable'
export ASTRATE_MQTT_INSECURE_DEV_MODE=true
export ASTRATE_MQTT_SESSION_STORE_PATH=/tmp/astrate-pokemon-sessions.db
export ASTRATE_MASTER_KEY='00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff'
export ASTRATE_REALM_NAME=test
export ASTRATE_REALM_JWT_PUBLIC_KEY_FILE=deploy/devrealm/realm_public.pem
go build -o /tmp/astrate ./cmd/astrate && /tmp/astrate
```

### Emulator (T3)
```sh
DEVICE_ID=$(cat /tmp/pokemon-device-id.txt)  # or register a new one
ROM="/Users/atsetilam/Downloads/Pokemon - Red Version (UE)[!]/Pokemon Red.gb"
cd examples/pokemon-agent/emulator-agent
python3 -m emulator_agent.main \
  --rom "$ROM" --insecure --astrate-url http://localhost:8080 \
  --realm test --device-id "$DEVICE_ID" --fps 60 \
  > /tmp/pokemon-emulator.log 2>&1 &
```

### Mint stream+publish JWT
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
```

### Orchestrator T4 (opencode)
```sh
export POKEMON_ASTRATE_URL=http://localhost:8080
export POKEMON_ASTRATE_REALM=test
export POKEMON_ASTRATE_DEVICE_ID=$(cat /tmp/pokemon-device-id.txt)
export POKEMON_ASTRATE_APP_TOKEN="$(cat /tmp/pokemon-app-token.txt)"
export POKEMON_OPENAI_MODEL=opencode/big-pickle
export POKEMON_LLM_BACKEND=opencode
export POKEMON_LLM_TIMEOUT_SECONDS=60
export POKEMON_LLM_MAX_RETRIES=2
cd examples/pokemon-agent/llm-orchestrator
python3 -m llm_orchestrator.main
```

## Architecture (do not re-discuss)
```
pyboy (in-process) ↔ Emulator Agent (Python) ←MQTT mTLS|dev plaintext→ Astrate ←WS/HTTP→ LLM Orchestrator
```
No AtomVM. No TCP bridge. Two services. Astrate is the bus.
Main loop owns all `pyboy.tick()`; MQTT only enqueues ControlCommands.
LLM backend: OpenAI-compatible HTTP **or** `opencode run` (Big Pickle, no key).

## Session 1 file tree (still accurate)
```
examples/pokemon-agent/
├── README.md, Makefile
├── astrate-interfaces/   (3 JSON interfaces)
├── emulator-agent/       (pyboy + MQTT mTLS / --insecure)
├── llm-orchestrator/     (App API + OpenAI HTTP or opencode/Big Pickle)
└── docs/ DESIGN.md DECISIONS.md TESTING.md
```
