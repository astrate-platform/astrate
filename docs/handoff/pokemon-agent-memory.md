# Pokemon Agent — Session Memory (2026-07-29, session 6)

## What happened this session

### Goal
P5 — Optional MkDocs nav for `examples/pokemon-agent/docs/DESIGN.md`.

### P5 result — **DONE**

MkDocs `docs_dir` is `docs/site/` (not the example tree). Wired the example
DESIGN into the site the same way as top-level `DESIGN.md` / `ROADMAP.md`:

1. **`docs/Makefile` `sync`** — copies
   `../examples/pokemon-agent/docs/DESIGN.md` → `site/pokemon-agent.md`
2. **`docs/mkdocs.yml`** — nav section **Examples → Pokémon Agent**
3. **`docs/site/index.md`** — Examples link under home page
4. **`make build`** — clean success; page in `site-dist/pokemon-agent/`

Canonical source remains `examples/pokemon-agent/docs/DESIGN.md`.
Refresh the site copy with `cd docs && make sync` (or `make build` / `make serve`).

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

---

## What's NOT done yet (next session scope)

### All numbered P0–P5 complete

### Optional follow-ups
- T4 when an OpenAI-compatible API key is available (orchestrator + live ControlCommand loop).
- Save-state / auto-boot past title intro so map/party telemetry is meaningful without manual play.
- InputExecutor still calls `pyboy.tick()` from the MQTT thread while the main loop also ticks — race; prefer a command queue drained in the main loop.
- Replace remaining `pass` unit tests (state_decoder fixtures, action_translator, context_builder).
- Dialog via `$CF4B` still may be empty during speech (tilemap/text engine); validate mid-dialog on real ROM after intro.
- Fix root `.dockerignore` excluding `docs/` (breaks compose full image build). Out of pokemon-agent scope.

## Risks and known issues
1. ~~WebSocket path~~ **P0**. ~~Dialog/maxHP WRAM~~ **P2**. ~~Smoke T1/T2~~ **P3**. ~~ROM T3~~ **P4**. ~~MkDocs nav~~ **P5**.
2. `paho-mqtt` 2.x `CallbackAPIVersion.VERSION2` required.
3. ~~pyboy headless uncapped CPU~~ mitigated by default `--fps 60`.
4. App token for orchestrator must include **channels** claim (`a_ch`).
5. Compose full build: `.dockerignore` has bare `docs` → missing swagger embed path.
6. `astartectl appengine devices data-snapshot` panics on object datastream shape; use `get-samples` or REST `/…/GameState/state`.
7. Unrelated WIP on same branch — leave alone.
8. Cold-boot mapId 0 → “Pallet Town” is WRAM-zero coincidence, not real overworld yet.
9. Site copy `docs/site/pokemon-agent.md` goes stale if DESIGN changes without `make sync`.

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
