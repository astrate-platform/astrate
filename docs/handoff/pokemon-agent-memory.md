# Pokemon Agent — Session Memory (2026-07-29, session 10)

## What happened this session

### Goal
**Overworld movement smoke** — get player coords to change under LLM control
after T4 bus path was proven (session 9 left player stuck at (3,6) stasis).

### Results — movement smoke **PASS**

Root causes of “inputs work but coords never change”:

1. **False intro-skip completion.** Pokémon Red preloads Red's House 2F WRAM
   coords `(3,6)` during Oak's speech — long before free overworld control.
   Old `looks_past_cold_boot()` treated that as “in game”, stopped auto-press,
   and left the agent mid-intro. MQTT ControlCommands applied, but the player
   could not walk.
2. **False dialog from `wStringBuffer`.** `$CF4B` retains name-entry residue
   (`ABBAAA`, `AAAAAAA`). That was published as `dialogText`, so the LLM only
   pressed A (“advance dialog”) and stasis stayed false.

Fixes (live-verified 2026-07-29 ~05:28–05:30):

| Fix | Where |
|---|---|
| Intro completes only when position **leaves** first non-zero baseline | `emulator-agent/.../main.py` |
| Direction probes in intro cycle (A/START + RIGHT/UP/LEFT/DOWN ×20f) | same |
| `is_actionable_dialog()` filters short all-caps residue | `state_decoder.py` + orchestrator `context_builder.py` |
| Directional `holdFrames` floored to **16** (Gen 1 tile step) | `action_translator.py` + system prompt |
| No stasis while intro auto-press still running | `main.py` |

Live evidence (device `HcLR8OkjSKat9ZAf8UolJg`, opencode/big-pickle):

- Intro: `spawn preload at (3, 6)` → `player moved (3, 6) → (4, 6)` (~76 s)
- LLM turns: `RIGHT×16, RIGHT×16, UP×16, DOWN×16, RIGHT×16, UP×16, RIGHT×16` (seq 1–7)
- Emulator: matching non-local `Input: … (seq=N)` lines
- GameState path: `(3,6) → (4,6) → (5,6) → (5,5) → (6,5)` map 38, `dialogText=""`
- Stairs target still ~(7,1); did not require leaving the house for this smoke

### Code / docs changed this session
Under `examples/pokemon-agent/` (+ handoff only outside that tree):

1. `emulator-agent/emulator_agent/main.py` — intro baseline + dir probes + stasis gate
2. `emulator-agent/emulator_agent/state_decoder.py` — `is_actionable_dialog`
3. `llm-orchestrator/.../action_translator.py` — floor direction holds to 16
4. `llm-orchestrator/.../context_builder.py` — holdFrames guidance + dialog filter + stairs hint
5. Tests: intro baseline, dialog filter, action_translator, context_builder
6. `docs/DESIGN.md`, `TESTING.md`, `README.md` — intro + dialog notes

### Branch
`feat/pokemon-agent`

Unrelated WIP still on the working tree (pipelines, broker ACL, triggers, go.mod).
Do **not** commit those with pokemon-agent work.

### Prior sessions
- Session 1: architecture, module tree, 10 ADRs
- Session 2–7: App API, tests, ROM T3, MkDocs, command queue, intro mash
- Session 8: T4 partial
- Session 9: T4 closed with Big Pickle (bus path; coords stuck)
- **Session 10: overworld movement under LLM; intro + dialog bugs fixed**

---

## What's NOT done yet (next session scope)

1. **Leave Red's House** — walk to stairs ~(7,1), warp to 1F / Pallet (longer LLM run or scripted guide).
2. Mid-dialog `$CF4B` validation on **real** text boxes (filter may hide edge cases; better textbox flag later).
3. Fill remaining thin unit tests / richer `state_decoder` fixtures.
4. Endurance / loop-detection demo (DESIGN P4–P5 aspirational).
5. Root `.dockerignore` `docs` → compose full image (out of scope unless asked).
6. Do **not** re-introduce Ollama unless user asks.
7. Optional: suppress publishing GameState heartbeats flooding LLM during long holds.

### T4 bus path + movement smoke are done — do not re-open unless regression.

## Risks and known issues
1. ~~WebSocket / JWT / T4 bus~~ **done**. ~~Intro false-complete~~ **fixed session 10**.
   ~~Name-entry dialog residue A-mash~~ **fixed session 10**. ~~holdFrames=8 half-step~~ **floored to 16**.
2. `paho-mqtt` 2.x `CallbackAPIVersion.VERSION2` required.
3. App JWT for stream must use REST-style `a_ch` (`.*::.*`), not JOIN/WATCH.
4. Compose full build: `.dockerignore` has bare `docs` → missing swagger embed path.
5. `astartectl appengine devices data-snapshot` panics on object datastream; use
   `get-samples` or REST `/…/GameState/state`.
6. Unrelated WIP on same branch — leave alone.
7. Cold-boot mapId 0 → “Pallet Town” is WRAM-zero coincidence until intro progresses.
8. Site copy `docs/site/pokemon-agent.md` goes stale if DESIGN changes without `make sync`.
9. Intro-skip timeout 180 s backstop; free play usually ~60–90 s with probes.
10. Emulator can orphan with MQTT CLOSED + high CPU; kill and restart; always log to a file.
11. Default LLM timeout 5 s is too low for `opencode run` cold start — use ≥60.
12. Big Pickle may refuse “play Pokémon” wording; engine frames as IoT fixture (keep that).
13. opencode backend is subprocess-per-turn — slow vs HTTP; fine for smoke, not production.
14. `$CF4B` is **not** a dedicated dialog register — filter is heuristic; real mid-dialog
    validation still open.
15. PyBoy writes `Pokemon Red.gb.ram` next to the ROM; delete for a clean New Game if needed.

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
DEVICE_ID=$(cat /tmp/pokemon-device-id.txt)
ROM="/Users/atsetilam/Downloads/Pokemon - Red Version (UE)[!]/Pokemon Red.gb"
# optional clean battery: rm -f "${ROM}.ram"
cd examples/pokemon-agent/emulator-agent
python3 -m emulator_agent.main \
  --rom "$ROM" --insecure --astrate-url http://localhost:8080 \
  --realm test --device-id "$DEVICE_ID" --fps 60 \
  > /tmp/pokemon-emulator.log 2>&1 &
# Wait for: Intro skip complete — player moved (3, 6) → …
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

### Orchestrator (opencode)
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
