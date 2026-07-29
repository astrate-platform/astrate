# Architecture Decision Record

**ADR-001:** Project lives in astrate repo as examples/pokemon-agent/
- Context: astrate is a Go reimplementation of Astarte; examples should live close to the platform code for discoverability and to act as integration tests
- Decision: examples/ subdirectory in the same repo
- Consequences: ROM file never committed (gitignored); example has its own Python virtualenvs

**ADR-002:** Two Python services, not a monolith
- Why not one process: decoupling allows the orchestrator to be remote, swapped, or scaled; Astrate is the interface contract
- Why not Go for the emulator agent: pyboy is Python-only; no Go Game Boy library with comparable memory API

**ADR-003:** AtomVM removed
- See §2.2 of DESIGN.md; referenced here for ADR completeness

**ADR-004:** pyboy as emulator
- Alternatives considered: mGBA (C, no Python API), BGB (Windows-only), custom emulator
- pyboy has a documented Python memory API, is actively maintained, supports headless mode, and can run at uncapped speed

**ADR-005:** paho-mqtt for emulator agent (not aio-pika, not aiomqtt)
- paho-mqtt 2.x has a thread-safe callback model; emulator agent's inner loop is CPU-bound (pyboy ticking), not I/O-bound, so threaded MQTT callbacks are fine and simpler than full async MQTT

**ADR-006:** Object aggregation for GameState
- Astrate delivers the entire object atomically; orchestrator never sees half a snapshot
- Alternative (individual): cheaper bandwidth but risks race conditions between mapId and playerX

**ADR-007:** Individual aggregation for PartyStatus
- A single battle exchange changes only one Pokémon's HP; publishing all 6 slots every frame wastes bandwidth
- The orchestrator maintains its own party cache and merges individual updates

**ADR-008:** sequenceId deduplication
- MQTT QoS 2 provides at-most-once delivery guarantees but reconnection edges still allow redelivery
- A monotonic integer in the command is the simplest replay defence; no external state store needed

**ADR-009:** Stasis detection in Emulator Agent, not Orchestrator
- The agent has frame-level position history; the orchestrator only sees published events
- Detecting stasis at the edge keeps the orchestrator stateless w.r.t. position history

**ADR-010:** LLM timeout → NONE command (not skip)
- Skipping leaves the last sequenceId gap, making dedup logic ambiguous
- A NONE command keeps the sequence contiguous and signals to any observer that the orchestrator is alive but slow
