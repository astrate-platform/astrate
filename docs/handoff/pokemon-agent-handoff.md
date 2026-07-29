# Handoff Prompt — Pokémon Agent Next Session (post queue + intro skip)

Copy-paste this prompt into a new session to continue work on the Pokémon Red autonomous agent.

````text
I'm working on the `feat/pokemon-agent` branch of ~/astrate — an autonomous Pokémon Red
agent that connects a Game Boy emulator (pyboy) to an LLM via the Astrate IoT platform.

Before doing anything, read:
  - ~/astrate/docs/handoff/pokemon-agent-memory.md   ← session summary + remaining work
  - ~/astrate/examples/pokemon-agent/docs/DESIGN.md  ← architecture (v0.2; §5.2–5.4 queue/intro)
  - ~/astrate/examples/pokemon-agent/docs/DECISIONS.md ← 10 ADRs
  - ~/astrate/examples/pokemon-agent/docs/TESTING.md ← T0–T4 (T4 notes: Big Pickle / opencode)

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
  Main loop owns every pyboy.tick(); MQTT only InputExecutor.enqueue().
  Default --skip-intro mashes A/START past title (not save-state).

P0–P5 DONE (see memory). Also DONE this branch: command queue + intro auto-press.

IMPORTANT: the branch working tree may contain unrelated WIP (pipelines, broker ACL,
triggers). Only edit/commit files under examples/pokemon-agent/ and docs/handoff/
(plus docs/mkdocs.yml / docs/Makefile / docs/site/ when touching site docs)
unless the user explicitly asks otherwise.

Current task (pick what the user wants; numbered phases are complete):

  Preferred next:
  - T4 LLM orchestrator end-to-end — run via **Big Pickle / opencode** (mule),
    not a long interactive strong-model watch. Put POKEMON_OPENAI_API_KEY only
    in that process env; never commit keys. App JWT needs channels claim a_ch.
    Emulator: T3 ROM + --insecure + default --skip-intro first.

  Also optional:
  - Fill remaining pass unit tests (orchestrator action_translator / context_builder,
    richer state_decoder fixtures)
  - Live-validate intro skip + mid-dialog $CF4B after intro on real ROM
  - After DESIGN edits: cd docs && make sync  (site/pokemon-agent.md)

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
