# Handoff Prompt — Pokémon Agent Next Session

Copy-paste this prompt into a new session to continue work on the Pokémon Red autonomous agent.

````text
I'm working on the `feat/pokemon-agent` branch of ~/astrate — an autonomous Pokémon Red
agent that connects a Game Boy emulator (pyboy) to an LLM via the Astrate IoT platform.

Before doing anything, read:
  - ~/astrate/docs/handoff/pokemon-agent-memory.md   ← full session summary & risks
  - ~/astrate/examples/pokemon-agent/docs/DESIGN.md  ← architecture (v0.2)
  - ~/astrate/examples/pokemon-agent/docs/DECISIONS.md ← 10 ADRs

Branch: feat/pokemon-agent (all code committed, 7 commits ahead of main)

Code lives in:
  examples/pokemon-agent/
    emulator-agent/   — Python/pyboy/paho-mqtt service
    llm-orchestrator/ — Python/aiohttp/httpx LLM service
    astrate-interfaces/ — 3 JSON interface definitions

Architecture (agreed, do NOT re-discuss):
  pyboy (in-process) ←→ Emulator Agent (Python) ←MQTT mTLS→ Astrate ←WebSocket→ LLM Orchestrator (Python)
  No AtomVM. No TCP bridge. Two services. Astrate is the bus.

Current task (pick the first that's still incomplete):

  P0 — Verify Astrate WebSocket stream route
       Read internal/appengine/stream/ in the repo to find the actual WebSocket endpoint
       path pattern. Update llm-orchestrator/llm_orchestrator/astrate_client.py if it
       differs from /v1/{realm}/devices/{device_id}/interfaces/{interface}.

  P1 — Run unit tests and fix failures
       cd examples/pokemon-agent/emulator-agent && pip install pytest pyboy && pytest tests/ -v
       cd examples/pokemon-agent/llm-orchestrator && pip install pytest pydantic pydantic-settings && pytest tests/ -v
       Fix any import errors or assertion failures. Commit fixes.

  P2 — Verify WRAM addresses
       Check pret/pokered wram.asm for:
         - Dialog buffer: $CC2A vs $CF4B (original doc had $CF4B)
         - Party HP offset: PARTY_SLOT_SIZE=44, SLOT_CURRENT_HP_OFFSET=1
       Update wram.py if wrong. Commit.

  P3 — Smoke test (requires running Astrate + astartectl)
       Follow examples/pokemon-agent/docs/TESTING.md T1 + T2.
       Install interfaces, run emulator agent in --stub mode, verify GameState
       events appear in Astrate App API.

  P4 — ROM integration (requires user to supply Pokémon Red ROM)
       Run full loop: pyboy + emulator agent + orchestrator + LLM.

Rules:
  - Read source before changing anything.
  - All changes on feat/pokemon-agent branch.
  - Commit after each task (P0, P1, P2, etc.) with clear messages.
  - At the end: update docs/handoff/pokemon-agent-memory.md with what changed,
    then update this handoff file with the remaining tasks, and tell the user
    which file to read next session.
````
