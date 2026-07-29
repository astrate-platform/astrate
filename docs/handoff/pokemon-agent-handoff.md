# Handoff Prompt — Pokémon Agent Next Session

Copy-paste this prompt into a new session to continue work on the Pokémon Red autonomous agent.

````text
I'm working on the `feat/pokemon-agent` branch of ~/astrate — an autonomous Pokémon Red
agent that connects a Game Boy emulator (pyboy) to an LLM via the Astrate IoT platform.

Before doing anything, read:
  - ~/astrate/docs/handoff/pokemon-agent-memory.md   ← session 10; movement smoke PASS
  - ~/astrate/examples/pokemon-agent/docs/DESIGN.md  ← architecture (v0.2)
  - ~/astrate/examples/pokemon-agent/docs/DECISIONS.md ← 10 ADRs
  - ~/astrate/examples/pokemon-agent/docs/TESTING.md ← T0–T4 + intro/movement notes

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
  Intro skip: A/START + direction probes; completes when player MOVES off WRAM
  spawn baseline (Oak preloads (3,6) — non-zero WRAM alone is NOT free play).
  dialogText: filter wStringBuffer name-entry residue (is_actionable_dialog).
  Direction holds floored to 16 frames (Gen 1 tile step).
  LLM: prefer opencode/big-pickle (no API key). Do NOT use Ollama unless asked.

P0–P5 scaffolding DONE. Command queue + intro auto-press DONE.
T4 DONE (session 9): WS + JWT + Big Pickle → ControlCommand → emulator Input (seq=N).
Overworld movement smoke DONE (session 10): LLM RIGHT/UP/… ×16 moves coords
  e.g. (3,6)→(4,6)→(5,6)→(5,5)→(6,5) on Red's House 2F.

IMPORTANT: the branch working tree may contain unrelated WIP (pipelines, broker ACL,
triggers). Only edit/commit files under examples/pokemon-agent/ and docs/handoff/
(plus docs/mkdocs.yml / docs/Makefile / docs/site/ when touching site docs)
unless the user explicitly asks otherwise.

Current task (pick one; primary suggestions):

  1. Leave Red's House — longer LLM run (or light guidance) to stairs ~(7,1),
     warp to 1F / Pallet Town (map change under LLM control).

  2. Mid-dialog validation on real ROM text boxes (current filter is heuristic
     on wStringBuffer; pret textbox flags would be more precise).

  3. Fill thin unit tests / richer state_decoder fixtures.

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
  Optional clean New Game: rm companion .ram next to the ROM.

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
