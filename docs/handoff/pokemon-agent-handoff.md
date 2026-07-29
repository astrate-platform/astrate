# Handoff Prompt — Pokémon Agent Next Session (post-P5)

Copy-paste this prompt into a new session to continue work on the Pokémon Red autonomous agent.

````text
I'm working on the `feat/pokemon-agent` branch of ~/astrate — an autonomous Pokémon Red
agent that connects a Game Boy emulator (pyboy) to an LLM via the Astrate IoT platform.

Before doing anything, read:
  - ~/astrate/docs/handoff/pokemon-agent-memory.md   ← session summary, P0–P5 results, risks
  - ~/astrate/examples/pokemon-agent/docs/DESIGN.md  ← architecture (v0.2)
  - ~/astrate/examples/pokemon-agent/docs/DECISIONS.md ← 10 ADRs
  - ~/astrate/examples/pokemon-agent/docs/TESTING.md ← T0–T4 smoke steps (updated for --insecure + T3 ROM)

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
P4 DONE: ROM integration — pyboy + real ROM + insecure Astrate; 60 fps pacing; dialog 0x00;
         stasis 15s; ControlCommand START delivered via AppEngine POST.
P5 DONE: MkDocs Examples → Pokémon Agent nav; docs/Makefile sync copies DESIGN →
         docs/site/pokemon-agent.md; index.md link. Run `cd docs && make sync` after DESIGN edits.

IMPORTANT: the branch working tree may contain unrelated WIP (pipelines, broker ACL,
triggers). Only edit/commit files under examples/pokemon-agent/ and docs/handoff/
(plus docs/mkdocs.yml / docs/Makefile / docs/site/ when touching site docs)
unless the user explicitly asks otherwise.

Current task (pick the first the user wants; all numbered P0–P5 are complete):

  Optional:
  - T4 LLM orchestrator when API key is available
  - Save-state / skip title so map/party telemetry is real
  - Queue ControlCommand ticks on the main loop (avoid MQTT-thread pyboy.tick race)
  - Fill remaining pass unit tests

ROM path (do NOT commit the ROM):
  /Users/atsetilam/Downloads/Pokemon - Red Version (UE)[!]/Pokemon Red.gb

Local Astrate smoke:
  docker compose up -d timescaledb
  ASTRATE_DATABASE_DSN=… ASTRATE_MQTT_INSECURE_DEV_MODE=true \
    ASTRATE_REALM_NAME=test ASTRATE_REALM_JWT_PUBLIC_KEY_FILE=deploy/devrealm/realm_public.pem \
    go build -o /tmp/astrate ./cmd/astrate && /tmp/astrate
  (docker compose --profile full build is broken: .dockerignore excludes docs/)

Rules:
  - Read source before changing anything.
  - All changes on feat/pokemon-agent branch.
  - Commit after each task with clear messages; do not sweep in
    unrelated pipelines/broker WIP.
  - At the end: update docs/handoff/pokemon-agent-memory.md with what changed,
    then update this handoff file with the remaining tasks, and tell the user
    which file to read next session.
````
