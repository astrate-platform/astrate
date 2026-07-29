# Handoff Prompt — Pokémon Agent Next Session (post-P2)

Copy-paste this prompt into a new session to continue work on the Pokémon Red autonomous agent.

````text
I'm working on the `feat/pokemon-agent` branch of ~/astrate — an autonomous Pokémon Red
agent that connects a Game Boy emulator (pyboy) to an LLM via the Astrate IoT platform.

Before doing anything, read:
  - ~/astrate/docs/handoff/pokemon-agent-memory.md   ← session summary, P0–P2 results, risks
  - ~/astrate/examples/pokemon-agent/docs/DESIGN.md  ← architecture (v0.2)
  - ~/astrate/examples/pokemon-agent/docs/DECISIONS.md ← 10 ADRs
  - ~/astrate/examples/pokemon-agent/docs/TESTING.md ← T1–T4 smoke steps

Branch: feat/pokemon-agent

Code lives in:
  examples/pokemon-agent/
    emulator-agent/   — Python/pyboy/paho-mqtt service
    llm-orchestrator/ — Python/aiohttp/httpx LLM service
    astrate-interfaces/ — 3 JSON interface definitions

Architecture (agreed, do NOT re-discuss):
  pyboy (in-process) ↔ Emulator Agent (Python) ←MQTT mTLS→ Astrate ←WS/HTTP→ LLM Orchestrator
  No AtomVM. No TCP bridge. Two services. Astrate is the bus.

P0 DONE (do not re-open unless source changed):
  Live stream: GET /astrate/v1/{realm}/socket?device_id=&interface=
  Publish:     POST /appengine/v1/{realm}/devices/{id}/interfaces/{iface}/{path}
  wireEvent:   event, realm, device_id, interface, path, value, timestamp
  Client fixed in llm-orchestrator/llm_orchestrator/astrate_client.py
  Source of truth: internal/appengine/stream/ws.go, internal/appengine/http.go

P1 DONE: unit tests green (emulator-agent + llm-orchestrator).

P2 DONE: WRAM verified against pret/pokered
  DIALOG_BUFFER = $CF4B (not $CC2A)
  PARTY_SLOT_SIZE=44, CURRENT_HP@+1, LEVEL@+33, MAX_HP@+34 (was wrongly +3)
  Commit 1ce8736; tests in emulator-agent/tests/test_wram.py

IMPORTANT: the branch working tree may contain unrelated WIP (pipelines, broker ACL,
triggers). Only edit/commit files under examples/pokemon-agent/ and docs/handoff/
unless the user explicitly asks otherwise.

Current task (pick the first that's still incomplete):

  P3 — Smoke test (requires running Astrate + astartectl)
       Follow examples/pokemon-agent/docs/TESTING.md T1 + T2.
       Install interfaces, run emulator agent in --stub mode, verify GameState
       via AppEngine GET or live socket (paths above). JWT needs a_ch for socket.

  P4 — ROM integration (requires user to supply Pokémon Red ROM)
       Run full loop: pyboy + emulator agent + orchestrator + LLM.

  P5 — Optional MkDocs nav entry for examples/pokemon-agent/docs/DESIGN.md

Rules:
  - Read source before changing anything.
  - All changes on feat/pokemon-agent branch.
  - Commit after each task with clear messages; do not sweep in
    unrelated pipelines/broker WIP.
  - At the end: update docs/handoff/pokemon-agent-memory.md with what changed,
    then update this handoff file with the remaining tasks, and tell the user
    which file to read next session.
````
