# System Architecture & Technical Design Document
## Autonomous Pokémon Red Agent via Astrate IoT Platform

**Status:** Revised after architectural review — see §2 for rationale.
**Version:** 0.2 (2026-07-29)

---

## 1. Executive Summary & Vision

This document outlines the architecture for an autonomous agent system capable of
playing *Pokémon Red* (Game Boy) using a Large Language Model (LLM) as the cognitive
driver.

The system bridges low-level game emulator state with high-level AI reasoning through
**two Python services** with **Astrate** as the central IoT message bus:

| Component | Role |
|---|---|
| **Emulator Agent** | Runs pyboy in-process; reads WRAM every tick; publishes telemetry to Astrate via MQTT; receives and executes controller commands |
| **Astrate Platform** | Wire-compatible Astarte IoT hub (this repo). Handles MQTT device pairing, interface schema enforcement, bidirectional data routing, and App API |
| **LLM Orchestrator** | Subscribes to Astrate App API (WebSocket/SSE); assembles game context; calls an LLM (OpenAI-compatible HTTP **or** `opencode run` / Big Pickle); dispatches controller commands back through Astrate |

```
┌──────────────────────────────────────────────────────────────────────┐
│                         EMULATOR AGENT (Python)                      │
│                                                                      │
│  ┌────────────────┐   ┌─────────────────┐   ┌───────────────────┐  │
│  │  pyboy (Game   │──▶│  WRAM Reader /  │──▶│  Astrate MQTT     │  │
│  │  Boy emulator) │   │  State Decoder  │   │  Client (mTLS)    │  │
│  └────────────────┘   └─────────────────┘   └────────┬──────────┘  │
│         ▲                                             │              │
│         │ joypad inject                               │ GameState    │
│  ┌──────┴──────────┐                                 │ PartyStatus  │
│  │ Input Executor  │◀── ControlCommand ──────────────┘              │
│  └─────────────────┘                                                 │
└──────────────────────────────────────────────────────────────────────┘
                              │ MQTT (mTLS)  ▲
                              ▼              │
┌──────────────────────────────────────────────────────────────────────┐
│                         ASTRATE PLATFORM                              │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────────┐ │
│  │  Realm Services                                                  │ │
│  │  · Data Updater Plant (DUP) — validates & persists telemetry    │ │
│  │  · Interfaces Engine — org.pokemon.emulator.*                   │ │
│  │  · App API & Triggers — WebSocket/SSE live stream               │ │
│  └──────────────────────────┬──────────────────────────────────────┘ │
└─────────────────────────────┼────────────────────────────────────────┘
                              │ WebSocket / HTTP App API
                              ▼
┌──────────────────────────────────────────────────────────────────────┐
│                       LLM ORCHESTRATOR (Python)                      │
│                                                                      │
│  ┌─────────────────┐  ┌──────────────┐  ┌────────────────────────┐ │
│  │ Context Builder │─▶│  LLM Engine  │─▶│  Action Translator &   │ │
│  │ (map/party/text)│  │  (inference) │  │  Command Dispatch      │ │
│  └─────────────────┘  └──────────────┘  └────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 2. Architectural Review & Decisions

### 2.1 Original design vs. revised design

The initial design included four hops:

```
TASEmulator → Python TCP Bridge → AtomVM (Erlang/BEAM) → Astrate → LLM (Python)
```

After review, both the Python TCP bridge and AtomVM were dropped.

### 2.2 Why AtomVM was removed

AtomVM's value proposition is running compiled BEAM bytecode on constrained
microcontrollers (ESP32, STM32, bare-metal). On a UNIX host running a software
emulator, its benefits do not apply:

- AtomVM on UNIX is an unusual deployment; standard Erlang/OTP would be used instead
- Most OTP hex packages (emqtt, etc.) do not compile for AtomVM's limited stdlib
- The BEAM supervisor/fault-tolerance model is fully replicated by Python asyncio
  task supervision for this use case
- It added a separate build toolchain (rebar3 + AtomVM pack) with no operational payoff
- The "IoT edge device" framing does not fit a desktop software emulator

**If the project is ever deployed on real hardware** (e.g., a Raspberry Pi Zero running
the emulator locally), AtomVM becomes a compelling option again. The Astrate interface
schemas and MQTT wire protocol are identical either way.

### 2.3 Why the Python TCP bridge was removed

The TCP bridge was only necessary as a translation layer between pyboy's Python API
and AtomVM's Erlang socket calls. With AtomVM gone, the Emulator Agent imports
pyboy directly in-process — no IPC, no intermediate protocol, no extra port to manage.

### 2.4 Why Astrate is kept

Astrate continues to provide real value as the message bus:

- **Decoupling:** The Emulator Agent and LLM Orchestrator have zero direct knowledge
  of each other. Astrate is the single source of truth for all game state.
- **Interface schemas:** `org.pokemon.emulator.*` interfaces enforce payload structure
  at the platform level — bugs in either service cannot corrupt the other's view.
- **Persistence:** Astrate's TimescaleDB backend stores every telemetry point, enabling
  offline replay and analysis of gameplay sessions.
- **Bidirectional routing:** Server-owned `ControlCommand` interfaces route commands
  from orchestrator to emulator agent without any direct TCP connection between them.
- **Observability:** The App API and Trigger engine provide out-of-the-box hooks for
  dashboards and alerting.

---

## 3. Subsystem Detailed Specifications

### 3.1 Emulator Agent

**Runtime:** Python 3.11+, single process, asyncio event loop.

**Library:** [pyboy](https://github.com/Binjitsu/PyBoy) — a Python Game Boy emulator
that exposes a direct API for reading/writing memory and injecting button presses.

**Key WRAM Addresses Monitored:**

| Address  | Contents |
|----------|----------|
| `$D35E`  | Current Map ID |
| `$D362`  | Player X coordinate |
| `$D361`  | Player Y coordinate |
| `$D163`  | Number of Pokémon in party (0–6) |
| `$D164`–`$D169` | Party species IDs (one byte per slot; list at `$D16A` = `$FF`) |
| `$D16B` + `n×44` | Party mon struct (`PARTYMON_STRUCT_LENGTH`); HP@+1, level@+33, maxHP@+34 |
| `$D057`  | Battle type flag (`0` = overworld, `1` = wild, `2` = trainer) |
| `$CF4B`  | Dialog / string buffer (`wStringBuffer`, 20 bytes; not `$CC2A`) |

**Core modules:**

| Module | Responsibility |
|---|---|
| `emulator_agent/wram.py` | WRAM address constants and named read helpers |
| `emulator_agent/state_decoder.py` | Converts raw bytes → structured `GameState` / `PartyStatus` dataclasses |
| `emulator_agent/astrate_client.py` | MQTT client (paho-mqtt over mTLS); publishes telemetry; subscribes to `ControlCommand` |
| `emulator_agent/input_executor.py` | Queues `ControlCommand` from MQTT; main loop drains + injects joypad; `sequenceId` dedup |
| `emulator_agent/main.py` | asyncio main loop; owns every `pyboy.tick()`; intro auto-press; stasis; shutdown |

**Tick rate:** 60 fps pyboy tick (CLI `--fps`, default 60; `0` = uncapped); WRAM
snapshot + publish on every *changed* state (position, battle flag, or dialog text
changes), plus a heartbeat every 5 seconds.

**Stasis detection:** If `playerX` and `playerY` remain unchanged for **15 seconds**
while `inBattle = false` and no dialog is active, the agent publishes a `GameState`
with `stasis=true` and logs a warning (time-based so it stays correct at 60 fps).

### 3.2 Astrate Platform

Astrate runs as a single binary or via `docker compose`. The emulator agent registers
as a device using Astrate's Pairing API (identical to Astarte's), receives mTLS
credentials, and connects to the embedded MQTT broker.

**Interface installation** (one-time, via `astartectl`):

```sh
astartectl utils install-interface \
  --astarte-url http://localhost:8080 \
  --realm-key <realm-private-key.pem> \
  --realm pokemon-dev \
  astrate-interfaces/org.pokemon.emulator.GameState.json

# repeat for PartyStatus and ControlCommand
```

See §3.3 for full interface schemas.

### 3.3 Astrate Interfaces

#### `org.pokemon.emulator.GameState`

- **Type:** Datastream
- **Aggregation:** Object — all fields published atomically as a single snapshot
- **Ownership:** Device (Emulator Agent)
- **Description:** Emitted whenever player position, battle state, or dialog text
  changes, plus a 5-second heartbeat. Object aggregation guarantees the orchestrator
  always sees a consistent snapshot, never a partial one.

```json
{
  "interface_name": "org.pokemon.emulator.GameState",
  "version_major": 1,
  "version_minor": 0,
  "type": "datastream",
  "ownership": "device",
  "aggregation": "object",
  "mappings": [
    { "endpoint": "/state/mapId",      "type": "integer" },
    { "endpoint": "/state/mapName",    "type": "string"  },
    { "endpoint": "/state/playerX",    "type": "integer" },
    { "endpoint": "/state/playerY",    "type": "integer" },
    { "endpoint": "/state/inBattle",   "type": "boolean" },
    { "endpoint": "/state/battleType", "type": "integer" },
    { "endpoint": "/state/dialogText", "type": "string"  },
    { "endpoint": "/state/stasis",     "type": "boolean" }
  ]
}
```

#### `org.pokemon.emulator.PartyStatus`

- **Type:** Datastream
- **Aggregation:** Individual — each slot can be updated independently
- **Ownership:** Device (Emulator Agent)
- **Description:** Emitted when any party member's HP changes. Individual aggregation
  means a single HP update for slot 0 does not require re-publishing all six slots.

```json
{
  "interface_name": "org.pokemon.emulator.PartyStatus",
  "version_major": 1,
  "version_minor": 0,
  "type": "datastream",
  "ownership": "device",
  "aggregation": "individual",
  "mappings": [
    { "endpoint": "/%{slotIndex}/name",      "type": "string"  },
    { "endpoint": "/%{slotIndex}/speciesId", "type": "integer" },
    { "endpoint": "/%{slotIndex}/currentHp", "type": "integer" },
    { "endpoint": "/%{slotIndex}/maxHp",     "type": "integer" },
    { "endpoint": "/%{slotIndex}/level",     "type": "integer" }
  ]
}
```

#### `org.pokemon.emulator.ControlCommand`

- **Type:** Datastream
- **Aggregation:** Object
- **Ownership:** Server (LLM Orchestrator)
- **Description:** Action orders issued by the LLM and routed by Astrate to the
  Emulator Agent. Server ownership means the Orchestrator publishes via the App API
  and Astrate delivers via MQTT to the registered device.

```json
{
  "interface_name": "org.pokemon.emulator.ControlCommand",
  "version_major": 1,
  "version_minor": 0,
  "type": "datastream",
  "ownership": "server",
  "aggregation": "object",
  "mappings": [
    { "endpoint": "/command/button",     "type": "string"  },
    { "endpoint": "/command/holdFrames", "type": "integer" },
    { "endpoint": "/command/sequenceId", "type": "integer" }
  ]
}
```

### 3.4 LLM Orchestrator

**Runtime:** Python 3.11+, asyncio.

**Core modules:**

| Module | Responsibility |
|---|---|
| `llm_orchestrator/astrate_client.py` | Live stream client (`GET /astrate/v1/{realm}/socket?device_id=&interface=`) + AppEngine publish (`POST /appengine/v1/.../interfaces/.../{path}`); reconnects with exponential backoff |
| `llm_orchestrator/context_builder.py` | Formats game state + party + dialog + action history into a structured LLM prompt |
| `llm_orchestrator/llm_engine.py` | LLM inference: OpenAI-compatible HTTP **or** `opencode run` (Big Pickle / free models, no API key); JSON validation; retries; `LLMTimeoutError` |
| `llm_orchestrator/action_translator.py` | Parses LLM JSON output; validates button names; assigns monotonically increasing `sequenceId` |
| `llm_orchestrator/main.py` | asyncio supervisor; event fan-in (GameState + PartyStatus); on `LLMTimeoutError` sends NONE no-op command |

**Execution flow:**

```
1. Astrate WebSocket delivers GameState event
2. context_builder assembles prompt:
     · Current location & coordinates
     · Battle state & active move options
     · Party HP percentages
     · On-screen dialog text
     · Last 5 actions taken (loop prevention)
3. llm_engine sends prompt → LLM → parses JSON response
4. action_translator validates output → builds ControlCommand payload
5. astrate_client POSTs ControlCommand to Astrate App API
6. Astrate routes command → Emulator Agent via MQTT
7. input_executor injects button press into pyboy
```

---

## 4. Data Flow Sequence

```
[Emulator Agent]          [Astrate Platform]       [LLM Orchestrator]
      │                          │                         │
      │── publish GameState ────▶│                         │
      │   (MQTT device topic)    │── WebSocket event ─────▶│
      │                          │   (App API stream)      │
      │                          │                         │── build prompt
      │                          │                         │── call LLM
      │                          │                         │── parse output
      │                          │◀── POST ControlCommand ─│
      │◀── MQTT server topic ────│   (App API HTTP)        │
      │    (ControlCommand)      │                         │
      │── inject joypad ──▶pyboy │                         │
```

---

## 5. Failure Recovery & Loop Mitigation

### 5.1 Stasis detection

The Emulator Agent tracks `(playerX, playerY)` across state change events. If the
position is unchanged for 15 consecutive events while not in battle and no dialog is
active, it sets `stasis=true` in the next `GameState` publish. The Orchestrator
recognises this flag and adjusts the LLM prompt to explicitly request a different
movement direction.

### 5.2 Command queue (main loop owns ticks)

MQTT callbacks run on paho’s background thread. They only call
`InputExecutor.enqueue()` — never `pyboy.tick()` or `send_input`. Each main-loop
frame does `before_tick()` → `pyboy.tick()` → `after_tick()` so hold-frames and
releases stay single-threaded with the emulator.

Local intro skip uses `enqueue_local()` (`_local=True`) and does **not** advance
the MQTT `sequenceId` space.

### 5.3 Intro auto-press (default for ROM mode)

`--skip-intro` (default on; `--no-skip-intro` to disable) mashes A/START while
WRAM still looks like cold boot (all-zero coords, empty party, no dialog). Stops
when `looks_past_cold_boot()` is true or after a timeout (~180 s). Preferred over
save-state loading for local smoke.

### 5.4 `sequenceId` deduplication

Every `ControlCommand` carries a monotonically increasing `sequenceId`. The
`input_executor` maintains the last-executed ID in memory; commands with
`sequenceId ≤ last_executed` are silently dropped. This prevents duplicate button
presses caused by MQTT redelivery or network latency spikes.

### 5.5 LLM timeout → no-op

If the LLM does not respond within `LLM_TIMEOUT_SECONDS` (default: 5; use ≥60 for
`opencode run` cold starts), the Orchestrator publishes a `ControlCommand` with
`button="NONE"` and `holdFrames=0`. This keeps the Astrate command sequence
contiguous and prevents the Emulator Agent from acting on a stale command.

**Backends** (`POKEMON_LLM_BACKEND`: `auto` | `openai` | `opencode`):
- `openai` — HTTP `/chat/completions` (`POKEMON_OPENAI_*`).
- `opencode` — subprocess `opencode run --model … --format json` (no API key for free models).
- `auto` — selects `opencode` when `POKEMON_OPENAI_MODEL` starts with `opencode/` or is `big-pickle`.

### 5.4 Reconnection

Both the Emulator Agent's MQTT client and the Orchestrator's WebSocket client
implement exponential backoff reconnection (initial 1 s, max 60 s, jitter ±10%).

---

## 6. Information Sources & References

1. **Astrate IoT Platform**
   - Repository & Architecture: <https://github.com/astrate-platform/astrate>
   - Interface Definition Model (Datastreams, Object Aggregation, Server/Device Ownership):
     <https://docs.astarte-platform.org/astarte/latest/040-interface_schema.html>
   - App API & Triggers:
     <https://docs.astarte-platform.org/astarte/latest/050-query_astarte.html>

   > **Note:** Documentation links point to `docs.astarte-platform.org` because
   > Astrate is 100% wire-compatible with Astarte. The interface schemas, MQTT
   > protocol, and App API documented there apply directly to Astrate.

2. **pyboy Game Boy Emulator**
   - Repository: <https://github.com/Binjitsu/PyBoy>
   - Memory API: <https://docs.pyboy.dk/index.html>

3. **Pokémon Red Game Boy Memory Map**
   - pret/pokered disassembly (WRAM addresses & data structures):
     <https://github.com/pret/pokered>
   - Game Boy Memory Map Specifications (`$C000–$DFFF` Work RAM, `$FF00` Joypad I/O):
     <https://gbdev.io/pandocs/Memory_Map.html>

---

## 7. Verification & Next Steps

| Phase | Scope | Success Criterion |
|---|---|---|
| **P0** | Interface install | `astartectl` installs all 3 interfaces without schema errors |
| **P1** | Emulator Agent stub | Agent publishes mock `GameState` every tick; visible in Astrate App API |
| **P2** | Full emulator loop | pyboy running Pokémon Red; telemetry flows to Astrate; manual command injects move |
| **P3** | Orchestrator online | LLM receives state, issues commands; character moves in Pallet Town |
| **P4** | Loop detection | Force stasis; verify `stasis=true` published; verify Orchestrator prompt changes |
| **P5** | Endurance | 30-minute unattended run on Route 1; no crash, no infinite loop |
