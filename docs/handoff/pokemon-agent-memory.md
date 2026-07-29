# Pokemon Agent — Session Memory (2026-07-29, session 7)

## What happened this session

### Goal
Implement the two small remaining hardening items: **intro auto-press** and
**ControlCommand queue**; leave T4 notes for Big Pickle / opencode.

### Results

#### Command queue — **DONE**
- MQTT callback only calls `InputExecutor.enqueue()` (no `pyboy.tick` / joypad).
- Main loop: `before_tick()` → `pyboy.tick()` → `after_tick()` owns hold/release.
- Unit tests in `emulator-agent/tests/test_input_executor.py` (dedup, hold, local).

#### Intro auto-press — **DONE** (preferred over save-state)
- Default `--skip-intro` (ROM mode); `--no-skip-intro` to disable; ignored with `--stub`.
- Mashes A/A/A/START while WRAM looks like cold boot; stops on
  `looks_past_cold_boot()` (non-zero coords / party / dialog / battle) or ~180 s.
- Local presses use `enqueue_local()` and do **not** advance MQTT `sequenceId`.

#### T4 — notes only (not run)
- Still no end-to-end orchestrator run in interactive sessions.
- **Next session should run T4 via Big Pickle / opencode** (mule-style), with
  `POKEMON_OPENAI_API_KEY` only in that process env — never commit keys.
- Emulator should be on T3 with default intro skip first; App JWT needs `a_ch`.

#### Docs
- DESIGN §5.2–5.4 queue + intro; TESTING T3/T4 notes; README flag note.
- All emulator unit tests green (15).

### Branch
`feat/pokemon-agent`

Unrelated WIP still on the working tree (pipelines, broker ACL, triggers, go.mod).
Do **not** commit those with pokemon-agent work.

### Prior sessions
- Session 1: architecture (two Python services), module tree, 10 ADRs
- Session 2 (P0): App API stream/publish paths in llm-orchestrator
- Session 3 (P1/P2): unit tests green; WRAM dialog `$CF4B`, maxHP@+34
- Session 4 (P3): live smoke T1+T2 stub + `--insecure`
- Session 5 (P4): ROM integration — pyboy + real ROM; 60 fps; dialog 0x00; stasis 15s
- Session 6 (P5): MkDocs Examples → Pokémon Agent nav + `make sync`

---

## What's NOT done yet (next session scope)

### Numbered P0–P5 complete; intro skip + command queue complete

### Remaining optional
- **T4** — run LLM orchestrator end-to-end. Prefer **Big Pickle via opencode**
  for the long-running smoke; supply OpenAI-compatible key in process env only.
- Fill remaining `pass` unit tests in orchestrator (`action_translator`,
  `context_builder`) and richer `state_decoder` fixtures.
- Dialog via `$CF4B` still may be empty mid-speech; validate after intro on real ROM.
- Fix root `.dockerignore` excluding `docs/` (compose full image). Out of scope.
- Save-state path deliberately **not** pursued; intro auto-press is the chosen approach.

## Risks and known issues
1. ~~WebSocket path~~ **P0**. ~~Dialog/maxHP WRAM~~ **P2**. ~~Smoke T1/T2~~ **P3**.
   ~~ROM T3~~ **P4**. ~~MkDocs nav~~ **P5**. ~~Command queue~~ **done**.
   ~~Intro auto-press~~ **done**.
2. `paho-mqtt` 2.x `CallbackAPIVersion.VERSION2` required.
3. ~~pyboy headless uncapped CPU~~ mitigated by default `--fps 60`.
4. App token for orchestrator must include **channels** claim (`a_ch`).
5. Compose full build: `.dockerignore` has bare `docs` → missing swagger embed path.
6. `astartectl appengine devices data-snapshot` panics on object datastream shape; use `get-samples` or REST `/…/GameState/state`.
7. Unrelated WIP on same branch — leave alone.
8. Cold-boot mapId 0 → “Pallet Town” is WRAM-zero coincidence until intro skip finishes.
9. Site copy `docs/site/pokemon-agent.md` goes stale if DESIGN changes without `make sync`.
10. Intro-skip heuristic can false-complete if dialog buffer happens to be non-empty on title;
    timeout is the backstop. Validate live once after this change.
11. T4 never live-tested; Big Pickle/opencode is the intended runner.

## Architecture (do not re-discuss)
```
pyboy (in-process) ↔ Emulator Agent (Python) ←MQTT mTLS|dev plaintext→ Astrate ←WS/HTTP→ LLM Orchestrator
```
No AtomVM. No TCP bridge. Two services. Astrate is the bus.
Main loop owns all `pyboy.tick()`; MQTT only enqueues ControlCommands.

## Session 1 file tree (still accurate)
```
examples/pokemon-agent/
├── README.md, Makefile
├── astrate-interfaces/   (3 JSON interfaces)
├── emulator-agent/       (pyboy + MQTT mTLS / --insecure)
├── llm-orchestrator/     (App API + OpenAI-compatible LLM)
└── docs/ DESIGN.md DECISIONS.md TESTING.md
```
