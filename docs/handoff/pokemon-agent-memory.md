# Pokemon Agent — Session Memory (2026-07-29, session 11)

## Merge status (platform)

- **PR #38** merged to `main` 2026-07-29 (`bb67a95`).
- Project line treated as **complete** for Astrate product prioritization; further
  pokemon work is optional polish only.
- Flow work continues separately: committed core on main + uncommitted factory/API
  WIP in the worktree (see `session-2026-07-29-flow-factory-handoff.md`).

## What happened this session

### Goal
**Leave Red's House** — walk to stairs ~(7,1), warp to 1F / Pallet Town under
control path through Astrate (handoff task 1 after session 10 movement smoke).

### Results — leave house **PASS** (light guidance)

| Step | Evidence |
|---|---|
| Intro skip | `player moved (3, 6) → (4, 6)` map Red's House 2F |
| 2F → stairs | light guide RIGHT/UP to (7,1) |
| Warp 1F | `mapId=38` → `37` Red's House 1F `(7,1)` |
| Exit to Pallet | `mapId=0` Pallet Town `(5,6)` (non-zero coords) |
| Bus path | `LIGHT → … (seq=N)` + emulator `Input: … (seq=N)` via MQTT |

Live device: `PaS87idwT5SizSYnH4q1BQ`. Full clean New Game run ~15:07–15:10 local.

### Why not pure LLM this session
- `opencode/big-pickle` consistently **timed out at 60s** on cold/slow turns.
- Raised recommended timeout to **180s**; still used **`POKEMON_GUIDANCE=light`**
  for the multi-step leave-house smoke (handoff explicitly allowed “or light guidance”).
- Per-turn GOAL hints for stairs were still added to the LLM prompt for future `llm` runs.

### Code / docs changed
Under `examples/pokemon-agent/` (+ handoff only outside that tree):

1. `llm-orchestrator/light_guide.py` — 2F stairs (7,1) + 1F south-door heuristics
2. `llm-orchestrator/config.py` — `guidance` (`llm`/`light`/`auto`), `turn_cooldown_seconds`
3. `llm-orchestrator/main.py` — light/auto path, debounce, clearer action logs
4. `llm-orchestrator/context_builder.py` — stairs GOAL + 1F exit hints in prompts
5. Tests: `test_light_guide.py`, context_builder goal tests (35 passed)
6. `docs/TESTING.md` T4b, orchestrator README; handoff memory/handoff

### Gotchas found live
1. Manual curl with `sequenceId=999` poisons the emulator dedup window — later
   orch seq 1…N are dropped until emulator restart.
2. Restarting emulator with `--no-skip-intro` + existing `.ram` sits on title
   screen at WRAM zeros (need CONTINUE / intro skip). Prefer clean New Game or
   intro mash that selects CONTINUE when a save exists (not implemented).
3. Branch `feat/pokemon-agent` has older DB migrations than `main` — reset
   Timescale `astrate` DB if `no migration found for version 8`.
4. Flow-factory WIP was stashed on `main` before checkout (`wip-flow-factory-20260729`).

### Branch
`feat/pokemon-agent`

### Prior sessions
- Session 9: T4 bus path (Big Pickle)
- Session 10: overworld movement; intro + dialog fixes
- **Session 11: leave house 2F→1F→Pallet via light guide + Astrate**

---

## What's NOT done yet (next session scope)

1. **LLM-only leave house** with timeout≥180 (or faster model) — prove pure
   `POKEMON_GUIDANCE=llm` can also reach Pallet without light guide.
2. Mid-dialog `$CF4B` validation on **real** text boxes / pret textbox flags.
3. Fill remaining thin unit tests / richer `state_decoder` fixtures.
4. Endurance / loop-detection demo (DESIGN P4–P5 aspirational).
5. Optional: CONTINUE on battery save when `--skip-intro`.
6. Root `.dockerignore` `docs` → compose full image (out of scope unless asked).
7. Do **not** re-introduce Ollama unless user asks.

### T4 bus + movement smoke + leave-house (light) are done — do not re-open unless regression.

## Risks and known issues
1. ~~WebSocket / JWT / T4~~ **done**. ~~Intro false-complete~~ **fixed s10**.
   ~~Name-entry dialog A-mash~~ **fixed s10**. ~~holdFrames half-step~~ **16**.
   ~~Leave house map change~~ **done s11 (light)**.
2. `paho-mqtt` 2.x `CallbackAPIVersion.VERSION2` required.
3. App JWT for stream must use REST-style `a_ch` (`.*::.*`), not JOIN/WATCH.
4. Compose full build: `.dockerignore` has bare `docs`.
5. High manual `sequenceId` breaks subsequent lower IDs until emulator restart.
6. Default LLM timeout 5s is far too low for opencode; use **≥180** when slow.
7. Big Pickle may refuse “play Pokémon” wording; keep IoT fixture framing.
8. opencode is subprocess-per-turn — fine for smoke, not production control rate.
9. `$CF4B` filter is still heuristic.
10. Light guide only covers Red's House 2F/1F; outdoor navigation is LLM-only.

## Recipes

### Leave house (light) — after T3 emulator is past intro
```sh
export POKEMON_ASTRATE_URL=http://localhost:8080
export POKEMON_ASTRATE_REALM=test
export POKEMON_ASTRATE_DEVICE_ID=$(cat /tmp/pokemon-device-id.txt)
export POKEMON_ASTRATE_APP_TOKEN="$(cat /tmp/pokemon-app-token.txt)"
export POKEMON_GUIDANCE=light
export POKEMON_TURN_COOLDOWN_SECONDS=0.4
export POKEMON_OPENAI_MODEL=opencode/big-pickle
export POKEMON_LLM_BACKEND=opencode
export POKEMON_LLM_TIMEOUT_SECONDS=180
cd examples/pokemon-agent/llm-orchestrator
python3 -m llm_orchestrator.main
```

### Astrate / emulator / JWT
Unchanged from session 10 memory (insecure :1883, REST `a_ch`, ROM path under
Downloads). If migrations fail with “version 8”, drop/recreate the `astrate` DB
when running this branch’s older binary.

## Architecture (do not re-discuss)
```
pyboy (in-process) ↔ Emulator Agent (Python) ←MQTT→ Astrate ←WS/HTTP→ LLM Orchestrator
```
No AtomVM. No TCP bridge. Light guide only chooses buttons; Astrate remains the bus.
