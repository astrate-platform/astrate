# Pokemon Agent — Session Memory (2026-07-29)

## What happened this session

### Context
The user provided a system architecture document for an "Autonomous Pokémon Red Agent"
connecting a Game Boy emulator to an LLM via Astrate. The document had a recurring typo
(astarte → astrate) and an over-engineered architecture (4 hops including AtomVM and a
Python TCP bridge).

### Architectural review (key decision)
After discussion, the architecture was simplified from:
```
TASEmulator → Python TCP Bridge → AtomVM (Erlang/BEAM) → Astrate → LLM (Python)
```
to:
```
Emulator Agent (Python/pyboy) ←→ Astrate ←→ LLM Orchestrator (Python)
```
Rationale documented in DESIGN.md §2 and DECISIONS.md (ADRs 003–005):
- AtomVM is valuable on microcontrollers; wasteful on UNIX for a desktop emulator
- Python bridge was only needed as AtomVM→pyboy glue; without AtomVM it's gone
- Astrate kept: decoupling, schema enforcement, TimescaleDB persistence, observability

### Branch
`feat/pokemon-agent` (branched from `main` at `656815c`)

### Commits on feat/pokemon-agent
```
70e69c9  feat(pokemon-agent): fill in real implementation (replace stubs)
30af401  feat(pokemon-agent): add top-level README and Makefile
5bfd535  feat(pokemon-agent): add ADR decision log and testing guide
eb11115  feat(pokemon-agent): add LLM orchestrator (Astrate App API + OpenAI inference)
5b56ff6  feat(pokemon-agent): add emulator agent (pyboy + Astrate MQTT client)
3e17614  feat(pokemon-agent): add Astrate interface JSON definitions
7a9dc06  feat(pokemon-agent): add revised design document (v0.2)
```

### Files created
```
examples/pokemon-agent/
├── README.md                          — overview, prerequisites, quick start, config ref
├── Makefile                           — setup/test/lint/interfaces-install/run targets
├── astrate-interfaces/
│   ├── org.pokemon.emulator.GameState.json       — device-owned, object aggregation
│   ├── org.pokemon.emulator.PartyStatus.json     — device-owned, individual aggregation
│   └── org.pokemon.emulator.ControlCommand.json  — server-owned, object aggregation
├── emulator-agent/
│   ├── requirements.txt               — pyboy>=2.0, paho-mqtt>=2.0
│   ├── emulator_agent/
│   │   ├── wram.py           ✅ REAL  — 151 species, 50+ maps, GB charset, HP big-endian
│   │   ├── state_decoder.py  ✅ REAL  — full party decode, states_differ, party_differs
│   │   ├── astrate_client.py ✅ REAL  — paho mTLS, introspection, JSON envelope, QoS 2
│   │   ├── input_executor.py ✅ REAL  — 8 buttons, press+hold+release, sequenceId dedup
│   │   └── main.py           ✅ REAL  — asyncio loop, stasis, heartbeat, --stub, CLI
│   └── tests/
│       ├── test_wram.py               — map names, species count, read_text, read_word_le
│       ├── test_state_decoder.py      — decode round-trip, differ helpers
│       └── fixtures/sample_snapshot.json
├── llm-orchestrator/
│   ├── requirements.txt               — aiohttp>=3.9, httpx>=0.27, pydantic>=2, pydantic-settings>=2
│   ├── llm_orchestrator/
│   │   ├── config.py         ✅ REAL  — pydantic-settings, POKEMON_* env vars
│   │   ├── context_builder.py ✅ REAL — full system + user prompt, stasis warning
│   │   ├── llm_engine.py     ✅ REAL  — httpx async, wait_for, 3 retries, LLMTimeoutError
│   │   ├── action_translator.py ✅ REAL — VALID_BUTTONS, sequenceId monotonic, holdFrames clamp
│   │   ├── astrate_client.py ✅ REAL  — aiohttp WS stream, backoff, httpx publish
│   │   └── main.py           ✅ REAL  — TaskGroup consumers, party cache, SIGTERM
│   └── tests/
│       ├── test_context_builder.py    — stasis, battle state, history in prompt
│       ├── test_action_translator.py  — valid/invalid buttons, sequenceId, holdFrames clamp
│       └── fixtures/sample_game_state.json
└── docs/
    ├── DESIGN.md      — full revised architecture doc (v0.2), corrected typos
    ├── DECISIONS.md   — 10 ADRs (ADR-001 through ADR-010)
    └── TESTING.md     — T0–T4 testing guide, known limitations

```

### What is REAL vs. still approximate
- All 9 core Python modules: fully implemented, not stubs
- Tests: structure and test cases written; will need `pip install pytest` to run
- WRAM address for dialog buffer: `$CC2A` used (verify against pret/pokered; `$CF4B` was in the original doc but `$CC2A` is what pokered wram.asm shows for `wTileMap`)
- `read_word_le` in wram.py is defined but HP uses `read_word_be` (Pokémon Red stores HP big-endian) — both exist
- The MQTT broker port detection heuristic in `astrate_client.py` assumes Astrate MQTT runs on 8883; verify against actual Astrate config
- The Astrate App API WebSocket path format (`/v1/{realm}/devices/{device_id}/interfaces/{interface}`) should be verified against the running Astrate instance — it may differ from upstream Astarte's Channels implementation (noted in Astrate's DESIGN.md §1.1 as a documented deviation)

## What's NOT done yet (next session scope)

### P0 — Verify Astrate App API WebSocket path
The Astrate `internal/appengine/stream` package uses a simplified WebSocket endpoint.
Read `internal/appengine/stream/` to confirm the actual route pattern and update
`llm-orchestrator/llm_orchestrator/astrate_client.py` if needed.

### P1 — Run unit tests
```sh
cd examples/pokemon-agent/emulator-agent && pip install pytest pyboy && pytest tests/ -v
cd examples/pokemon-agent/llm-orchestrator && pip install pytest pydantic pydantic-settings && pytest tests/ -v
```
Fix any import errors or assertion failures.

### P2 — Verify WRAM addresses against pokered
Cross-check `wram.py` constants against https://github.com/pret/pokered/blob/master/wram.asm:
- Dialog buffer address (`$CC2A` vs `$CF4B` — original doc had `$CF4B` which is `wcd4b` / menu text)
- Party HP offset layout (confirm `PARTY_SLOT_SIZE = 44` and `SLOT_CURRENT_HP_OFFSET = 1`)

### P3 — Smoke test with Astrate running
Follow TESTING.md T1 + T2: install interfaces via astartectl, run emulator agent in stub mode,
confirm GameState events appear in Astrate App API.

### P4 — ROM integration (requires user to supply ROM)
Run with a real Pokémon Red ROM, confirm pyboy ticks without crashing,
telemetry flows, LLM orchestrator can dispatch commands.

### P5 — Add to MkDocs nav (optional)
Add `examples/pokemon-agent/docs/DESIGN.md` to the main docs site nav as an "Examples" section.

## Risks and known issues
1. Astrate WebSocket stream format may differ from Astarte Channels — needs verification before P3.
2. Dialog buffer WRAM address needs confirming; wrong address = empty dialog strings (not a crash).
3. `paho-mqtt` 2.x API changed significantly from 1.x — the `CallbackAPIVersion.VERSION2` usage is correct for 2.x but will fail on older installs.
4. pyboy `set_emulation_speed(0)` behavior in headless mode should be tested — some versions behave differently.
5. The `--stub` mode's `_StubPyboy` uses a plain `bytearray`; `pyboy.memory` in real pyboy is a `memoryview` — tests should mock accordingly.
