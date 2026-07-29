# Handoff Prompt — Pokémon Agent Next Session

Copy-paste this prompt into a new session to continue work on the Pokémon Red autonomous agent.

````text
I'm working on the `feat/pokemon-agent` branch of ~/astrate — an autonomous Pokémon Red
agent that connects a Game Boy emulator (pyboy) to an LLM via the Astrate IoT platform.

Before doing anything, read:
  - ~/astrate/docs/handoff/pokemon-agent-memory.md   ← session 9; T4 PASS
  - ~/astrate/examples/pokemon-agent/docs/DESIGN.md  ← architecture (v0.2)
  - ~/astrate/examples/pokemon-agent/docs/DECISIONS.md ← 10 ADRs
  - ~/astrate/examples/pokemon-agent/docs/TESTING.md ← T0–T4 (JWT + opencode)

Branch: feat/pokemon-agent

Code lives in:
  examples/pokemon-agent/
    emulator-agent/   — Python/pyboy/paho-mqtt service
    llm-orchestrator/ — App API + LLM (OpenAI HTTP or opencode/Big Pickle)
    astrate-interfaces/ — 3 JSON interface definitions

Architecture (agreed, do NOT re-discuss):
  pyboy (in-process) ↔ Emulator Agent (Python) ←MQTT→ Astrate ←WS/HTTP→ LLM Orchestrator
  No AtomVM. No TCP bridge. Two services. Astrate is the bus.
  Local smoke uses --insecure (plaintext :1883); production uses mTLS :8883.
  Main loop owns every pyboy.tick(); MQTT only InputExecutor.enqueue().
  Default --skip-intro mashes A/START past title (live-verified → Red's House 2F).
  LLM: prefer opencode/big-pickle (no API key). Do NOT use Ollama unless asked.

P0–P5 scaffolding DONE. Command queue + intro auto-press DONE.
T4 DONE (session 9): WS + JWT + Big Pickle → ControlCommand → emulator Input (seq=N).

IMPORTANT: the branch working tree may contain unrelated WIP (pipelines, broker ACL,
triggers). Only edit/commit files under examples/pokemon-agent/ and docs/handoff/
(plus docs/mkdocs.yml / docs/Makefile / docs/site/ when touching site docs)
unless the user explicitly asks otherwise.

Current task (pick one; primary suggestions):

  1. Overworld movement smoke — get player coords to change under LLM control
     (leave Red's House 2F via stairs; longer holds; clear stasis). Stack recipes
     and JWT mint are in pokemon-agent-memory.md / TESTING.md T4.

  2. Fill thin unit tests: action_translator (already has a pass stub — expand),
     context_builder, emulator state_decoder fixtures.

  3. Mid-dialog $CF4B validation on real ROM once text boxes appear.

  4. Longer unattended run / stasis-loop demo (DESIGN verification P4–P5).

JWT gotcha (do not forget):
  Mint with a_ch REST grant — NOT astartectl channels JOIN/WATCH:
    {"a_aea":[".*::.*"],"a_ch":[".*::.*"]}  signed with deploy/devrealm/realm_private.pem
  (Astrate stream uses Authorizes REST grammar on ClaimChannels; JOIN/WATCH → 403.)

Orchestrator (T4 recipe):
  POKEMON_OPENAI_MODEL=opencode/big-pickle
  POKEMON_LLM_BACKEND=opencode
  POKEMON_LLM_TIMEOUT_SECONDS=60
  POKEMON_ASTRATE_* as in memory
  No POKEMON_OPENAI_API_KEY for free opencode models.

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
