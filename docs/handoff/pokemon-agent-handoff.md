# Handoff Prompt — Pokémon Agent Next Session (post-P3)

Copy-paste this prompt into a new session to continue work on the Pokémon Red autonomous agent.

````text
I'm working on the `feat/pokemon-agent` branch of ~/astrate — an autonomous Pokémon Red
agent that connects a Game Boy emulator (pyboy) to an LLM via the Astrate IoT platform.

Before doing anything, read:
  - ~/astrate/docs/handoff/pokemon-agent-memory.md   ← session summary, P0–P3 results, risks
  - ~/astrate/examples/pokemon-agent/docs/DESIGN.md  ← architecture (v0.2)
  - ~/astrate/examples/pokemon-agent/docs/DECISIONS.md ← 10 ADRs
  - ~/astrate/examples/pokemon-agent/docs/TESTING.md ← T0–T4 smoke steps (updated for --insecure)

Branch: feat/pokemon-agent

Code lives in:
  examples/pokemon-agent/
    emulator-agent/   — Python/pyboy/paho-mqtt service
    llm-orchestrator/ — Python/aiohttp/httpx LLM service
    astrate-interfaces/ — 3 JSON interface definitions

Architecture (agreed, do NOT re-discuss):
  pyboy (in-process) ↔ Emulator Agent (Python) ←MQTT→ Astrate ←WS/HTTP→ LLM Orchestrator
  No AtomVM. No TCP bridge. Two services. Astrate is the bus.
  Local smoke uses --insecure (plaintext :1883); production uses mTLS :8883.

P0 DONE: App API stream/publish paths (llm-orchestrator client).
P1 DONE: unit tests green.
P2 DONE: WRAM verified (DIALOG $CF4B, maxHP@+34).
P3 DONE: live smoke T1+T2 — interfaces install, stub agent, GameState+PartyStatus in AppEngine.

IMPORTANT: the branch working tree may contain unrelated WIP (pipelines, broker ACL,
triggers). Only edit/commit files under examples/pokemon-agent/ and docs/handoff/
unless the user explicitly asks otherwise.

Current task (pick the first that's still incomplete):

  P4 — ROM integration (ROM path supplied by user — do NOT commit the ROM)
       ROM path:
         /Users/atsetilam/Downloads/Pokemon - Red Version (UE)[!]/Pokemon Red.gb
       Run full loop: pyboy + emulator agent (+ optional orchestrator + LLM).
       Follow TESTING.md T3 (and T4 if LLM key available).
       Prefer local Astrate: timescaledb compose + host binary insecure_dev_mode
       (docker compose --profile full build is broken: .dockerignore excludes docs/).
       Example:
         python -m emulator_agent.main \
           --rom "/Users/atsetilam/Downloads/Pokemon - Red Version (UE)[!]/Pokemon Red.gb" \
           --insecure --astrate-url http://localhost:8080 --realm test \
           --device-id <device-id>

  P5 — Optional MkDocs nav for examples/pokemon-agent/docs/DESIGN.md
       MkDocs docs_dir is docs/site/; needs site page + mkdocs.yml (ask first).

Rules:
  - Read source before changing anything.
  - All changes on feat/pokemon-agent branch.
  - Commit after each task with clear messages; do not sweep in
    unrelated pipelines/broker WIP.
  - At the end: update docs/handoff/pokemon-agent-memory.md with what changed,
    then update this handoff file with the remaining tasks, and tell the user
    which file to read next session.
````
