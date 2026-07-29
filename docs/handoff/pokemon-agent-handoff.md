# Handoff Prompt — Pokémon Agent Next Session

> **Project line complete on platform `main`.** Merged via **PR #38** (2026-07-29,
> merge `bb67a95`). Example lives at `examples/pokemon-agent/`. Optional polish only —
> default Astrate work is Flow (`session-2026-07-29-flow-factory-handoff.md`).

Copy-paste this prompt into a new session to continue **optional** Pokémon Red agent work.

````text
I'm working on the Astrate project in ~/astrate. The Pokémon Red autonomous agent
example is already on main (PR #38). Optional follow-ups only.

Before doing anything, read:
  - ~/astrate/docs/handoff/pokemon-agent-memory.md   ← session 11; leave house PASS; merged
  - ~/astrate/examples/pokemon-agent/docs/DESIGN.md  ← architecture (v0.2)
  - ~/astrate/examples/pokemon-agent/docs/DECISIONS.md ← 10 ADRs
  - ~/astrate/examples/pokemon-agent/docs/TESTING.md ← T0–T4 + T4b leave-house

Branch: main (code under examples/pokemon-agent/; feat/pokemon-agent was merged)

Code lives in:
  examples/pokemon-agent/
    emulator-agent/   — Python/pyboy/paho-mqtt service
    llm-orchestrator/ — App API + LLM (OpenAI HTTP or opencode/Big Pickle) + light guide
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
  Light guide: POKEMON_GUIDANCE=light|auto for Red's House 2F→1F→Pallet (still via MQTT).

P0–P5 scaffolding DONE. Command queue + intro auto-press DONE.
T4 DONE (session 9). Overworld movement DONE (session 10).
Leave Red's House DONE (session 11, light guide): 2F → 1F → Pallet Town (5,6).

IMPORTANT: only edit/commit files under examples/pokemon-agent/ and docs/handoff/
(plus docs/mkdocs.yml / docs/Makefile / docs/site/ when touching site docs)
unless the user explicitly asks otherwise.

Current task (pick one; primary suggestions):

  1. LLM-only leave house / outdoor steps with POKEMON_GUIDANCE=llm and
     POKEMON_LLM_TIMEOUT_SECONDS≥180 (prove pure LLM path to Pallet or further).

  2. Mid-dialog validation on real ROM text boxes (current filter is heuristic
     on wStringBuffer; pret textbox flags would be more precise).

  3. Fill thin unit tests / richer state_decoder fixtures.

  4. Longer unattended run / stasis-loop demo (DESIGN verification P4–P5).

  5. Optional: intro skip selects CONTINUE when a battery .ram exists.

JWT gotcha (do not forget):
  Mint with a_ch REST grant — NOT astartectl channels JOIN/WATCH:
    {"a_aea":[".*::.*"],"a_ch":[".*::.*"]}  signed with deploy/devrealm/realm_private.pem

Orchestrator recipes:
  # Leave house (fast smoke)
  POKEMON_GUIDANCE=light POKEMON_TURN_COOLDOWN_SECONDS=0.4
  # Pure LLM
  POKEMON_GUIDANCE=llm POKEMON_LLM_TIMEOUT_SECONDS=180
  POKEMON_OPENAI_MODEL=opencode/big-pickle POKEMON_LLM_BACKEND=opencode

sequenceId gotcha:
  Do not curl-inject a high sequenceId mid-run; emulator drops later lower IDs
  until process restart.

ROM path (do NOT commit the ROM):
  /Users/atsetilam/Downloads/Pokemon - Red Version (UE)[!]/Pokemon Red.gb
  Optional clean New Game: rm companion .ram next to the ROM.

Local Astrate smoke:
  docker compose up -d timescaledb
  ASTRATE_DATABASE_DSN=… ASTRATE_MQTT_INSECURE_DEV_MODE=true \
    ASTRATE_REALM_NAME=test ASTRATE_REALM_JWT_PUBLIC_KEY_FILE=deploy/devrealm/realm_public.pem \
    go build -o /tmp/astrate ./cmd/astrate && /tmp/astrate
  If migrate errors on “version 8”, drop/recreate DB (this branch has fewer migrations than main).

Rules:
  - Read source before changing anything.
  - All changes on feat/pokemon-agent branch.
  - Commit after each task with clear messages; do not sweep in
    unrelated pipelines/broker WIP.
  - At the end: update docs/handoff/pokemon-agent-memory.md with what changed,
    then update this handoff file with the remaining tasks, and tell the user
    which file to read next session.
````
