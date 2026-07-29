# Testing Guide

**T0 — Python unit tests (no dependencies)**
```sh
cd examples/pokemon-agent/emulator-agent && pip install pytest && pytest tests/
cd examples/pokemon-agent/llm-orchestrator && pip install pytest pydantic pydantic-settings && pytest tests/
```
Expected: all tests pass. No ROM, no Astrate, no LLM needed.

**T1 — Interface schema validation**
```sh
# Validate JSON syntax
for f in examples/pokemon-agent/astrate-interfaces/*.json; do python3 -m json.tool $f > /dev/null && echo "OK: $f"; done
# Install into running Astrate
astartectl utils install-interface --astarte-url http://localhost:8080 --realm-key realm.key --realm pokemon-dev \
  examples/pokemon-agent/astrate-interfaces/org.pokemon.emulator.GameState.json
```
Expected: 200 OK from Astrate, interface visible in realm management API.

**T2 — Emulator Agent stub mode**
```sh
cd examples/pokemon-agent/emulator-agent
pip install -r requirements.txt
python -m emulator_agent.main \
  --stub \
  --astrate-url http://localhost:8080 \
  --realm pokemon-dev \
  --device-id <device-id> \
  --cert path/to/device.crt \
  --key path/to/device.key \
  --ca path/to/ca.crt
```
Expected: GameState events visible in Astrate App API at /v1/pokemon-dev/devices/<id>/interfaces/org.pokemon.emulator.GameState

**T3 — Full emulator loop**
```sh
python -m emulator_agent.main \
  --rom /path/to/pokemon_red.gb \
  ...
```
Expected: telemetry flowing, visible in Astrate dashboard or via astartectl.

**T4 — LLM Orchestrator with stub**
```sh
cd examples/pokemon-agent/llm-orchestrator
pip install -r requirements.txt
export POKEMON_ASTRATE_URL=http://localhost:8080
export POKEMON_ASTRATE_REALM=pokemon-dev
export POKEMON_ASTRATE_DEVICE_ID=<device-id>
export POKEMON_ASTRATE_APP_TOKEN=<jwt>
export POKEMON_OPENAI_API_KEY=<key>
python -m llm_orchestrator.main
```
Expected: ControlCommand events visible in Astrate; if emulator agent is also running, character moves.
