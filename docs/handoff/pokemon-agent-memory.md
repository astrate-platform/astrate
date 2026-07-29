# Pokemon Agent — Session Memory (2026-07-29, session 5)

## What happened this session

### Goal
P4 — ROM integration (T3 full emulator loop with real Pokémon Red + pyboy).

### P4 result — **PASS** (live Astrate + real ROM)

Stack (same as P3):
1. `docker compose up -d timescaledb` (already healthy)
2. Host binary: `/tmp/astrate` with `ASTRATE_MQTT_INSECURE_DEV_MODE=true`,
   realm `test` + `deploy/devrealm/realm_public.pem`
3. Interfaces already installed; registered device `HcLR8OkjSKat9ZAf8UolJg`
4. ROM (outside repo, do not commit):
   `/Users/atsetilam/Downloads/Pokemon - Red Version (UE)[!]/Pokemon Red.gb`
5. `python -m emulator_agent.main --rom … --insecure --realm test --device-id …`

Verified via AppEngine REST `/…/GameState/state`:

| Field | Title-screen (cold boot) value |
|---|---|
| mapId / mapName | `0` / `Pallet Town` (zeroed WRAM until intro finishes) |
| playerX / playerY | `0` / `0` |
| dialogText | `""` (after fix; was `????…` from 0x00 bytes) |
| inBattle | `false` |
| stasis | `true` after 15 s stationary overworld |
| PartyStatus | none (party count 0 until in-game) |
| Heartbeat | ~5 s |
| ControlCommand | AppEngine POST → agent log `Input: START ×N frames` |

CPU: default `--fps 60` ≈ **4%**; uncapped (`set_emulation_speed(0)` + `sleep(0)`) was **~94%**.

T4 (LLM orchestrator) **not run** — no `POKEMON_OPENAI_API_KEY` / OpenAI key in environment.

### Code changes this session (emulator agent + docs)

1. **`--fps` / 60 fps pacing** (`main.py`) — DESIGN §3.1; avoids full-core spin in headless ROM mode. `--fps 0` keeps uncapped.
2. **Stasis is time-based (15 s)** — old counter of 15 *ticks* fired in 0.25 s at 60 fps; now `STASIS_SECONDS = 15.0`.
3. **`read_text` stops on 0x00** — zeroed dialog buffer no longer becomes `????????????????????` (title screen / cold boot).
4. Unit tests: zeroed buffer + charset “Hello” decode (9 emulator tests green).
5. Makefile `run-emulator-rom`, TESTING.md T3, README T3, DESIGN stasis/fps notes.

### Branch
`feat/pokemon-agent`

Unrelated WIP still on the working tree (pipelines, broker ACL, triggers, go.mod).
Do **not** commit those with pokemon-agent work.

### Prior sessions
- Session 1: architecture (two Python services), module tree, 10 ADRs
- Session 2 (P0): App API stream/publish paths in llm-orchestrator
- Session 3 (P1/P2): unit tests green; WRAM dialog `$CF4B`, maxHP@+34
- Session 4 (P3): live smoke T1+T2 stub + `--insecure`

---

## What's NOT done yet (next session scope)

### P5 — Optional MkDocs nav entry
MkDocs `docs_dir` is `docs/site/`; example DESIGN lives under
`examples/pokemon-agent/docs/`. Would need a site page or symlink +
`docs/mkdocs.yml` edit (outside pure example tree — **ask first**).

### Optional follow-ups
- T4 when an OpenAI-compatible API key is available (orchestrator + live ControlCommand loop).
- Save-state / auto-boot past title intro so map/party telemetry is meaningful without manual play.
- InputExecutor still calls `pyboy.tick()` from the MQTT thread while the main loop also ticks — race; prefer a command queue drained in the main loop.
- Replace remaining `pass` unit tests (state_decoder fixtures, action_translator, context_builder).
- Dialog via `$CF4B` still may be empty during speech (tilemap/text engine); validate mid-dialog on real ROM after intro.
- Fix root `.dockerignore` excluding `docs/` (breaks compose full image build). Out of pokemon-agent scope.

## Risks and known issues
1. ~~WebSocket path~~ **P0**. ~~Dialog/maxHP WRAM~~ **P2**. ~~Smoke T1/T2~~ **P3**. ~~ROM T3~~ **P4**.
2. `paho-mqtt` 2.x `CallbackAPIVersion.VERSION2` required.
3. ~~pyboy headless uncapped CPU~~ mitigated by default `--fps 60`.
4. App token for orchestrator must include **channels** claim (`a_ch`).
5. Compose full build: `.dockerignore` has bare `docs` → missing swagger embed path.
6. `astartectl appengine devices data-snapshot` panics on object datastream shape; use `get-samples` or REST `/…/GameState/state`.
7. Unrelated WIP on same branch — leave alone.
8. Cold-boot mapId 0 → “Pallet Town” is WRAM-zero coincidence, not real overworld yet.

## Architecture (do not re-discuss)
```
pyboy (in-process) ↔ Emulator Agent (Python) ←MQTT mTLS|dev plaintext→ Astrate ←WS/HTTP→ LLM Orchestrator
```
No AtomVM. No TCP bridge. Two services. Astrate is the bus.

## Session 1 file tree (still accurate)
```
examples/pokemon-agent/
├── README.md, Makefile
├── astrate-interfaces/   (3 JSON interfaces)
├── emulator-agent/       (pyboy + MQTT mTLS / --insecure)
├── llm-orchestrator/     (App API + OpenAI-compatible LLM)
└── docs/ DESIGN.md DECISIONS.md TESTING.md
```
