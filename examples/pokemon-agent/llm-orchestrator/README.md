# LLM Orchestrator

Subscribes to Astrate App API streams (GameState + PartyStatus), calls an LLM each
turn, and publishes `ControlCommand` back through AppEngine.

## LLM backends

| Backend | When | Auth |
|---|---|---|
| **opencode** (recommended for local smoke) | `POKEMON_LLM_BACKEND=opencode` or model `opencode/*` / `big-pickle` | none for free models |
| **openai** | `POKEMON_LLM_BACKEND=openai` or non-opencode model under `auto` | `POKEMON_OPENAI_API_KEY` |

## Quick start (opencode / Big Pickle)

```sh
# JWT must use REST-style a_ch (.*::.*), not astartectl JOIN/WATCH — see docs/TESTING.md T4
export POKEMON_ASTRATE_URL=http://localhost:8080
export POKEMON_ASTRATE_REALM=test
export POKEMON_ASTRATE_DEVICE_ID=<device-id>
export POKEMON_ASTRATE_APP_TOKEN="$(cat /tmp/pokemon-app-token.txt)"
export POKEMON_OPENAI_MODEL=opencode/big-pickle
export POKEMON_LLM_BACKEND=opencode
export POKEMON_LLM_TIMEOUT_SECONDS=180

pip install -r requirements.txt
python3 -m llm_orchestrator.main
```

### Light guidance (leave Red's House without a slow LLM)

```sh
export POKEMON_GUIDANCE=light   # deterministic 2F stairs + 1F south door
# POKEMON_GUIDANCE=auto uses light when it has a path, else LLM
python3 -m llm_orchestrator.main
```

Still publishes `ControlCommand` through AppEngine/MQTT. See `../docs/TESTING.md` T4b.

See `../docs/TESTING.md` (T4/T4b) and `../docs/DESIGN.md` §3.4 for architecture and JWT mint.
