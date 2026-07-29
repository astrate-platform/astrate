# Pokemon Agent — Session Memory (2026-07-29, session 2)

## What happened this session

### Goal
P0 from the previous handoff: verify Astrate WebSocket stream route against
`internal/appengine/stream/` and fix the LLM orchestrator client if wrong.

### P0 result — paths were wrong; client updated

Verified against source (not Astarte upstream docs):

| Surface | Correct Astrate path | What the client had before |
|---|---|---|
| Live stream (WS/SSE) | `GET /astrate/v1/{realm}/socket?device_id=&interface=` | `GET /v1/{realm}/devices/{id}/interfaces/{iface}` |
| REST publish | `POST /appengine/v1/{realm}/devices/{id}/interfaces/{iface}/{path}` | `POST /v1/{realm}/devices/{id}/interfaces/{iface}{path}` |
| Auth on socket | Bearer JWT, `a_ch` claim (`auth.ClaimChannels`) | Bearer (unchanged; claim type is ops concern) |
| Stream event JSON | `wireEvent`: `event`, `realm`, `device_id`, `interface`, `path`, `value`, `timestamp` | Assumed `interface`/`path`/`value`/`timestamp` only |

Evidence:
- `internal/appengine/stream/ws.go` — `Mount` registers `GET /astrate/v1/{realm}/socket`; query filters; `wireEvent` shape
- `internal/appengine/stream/ws_test.go` — SSE test hits `/astrate/v1/testrealm/socket?device_id=dev1&interface=com.ex.S`
- `internal/appengine/http.go` — AppEngine base `/appengine/v1/{realm}`; POST/PUT data routes
- `docs/site/appengine-api.md` — documents both surfaces
- `docs/COMPATIBILITY.md` — Astrate-native socket is a deliberate §1.1 deviation from Phoenix Channels

Consumer compatibility: `main.py` only uses `event["path"]` and `event["value"]`, which
match `wireEvent`. No main.py change required for P0.

### Files changed this session
```
examples/pokemon-agent/llm-orchestrator/llm_orchestrator/astrate_client.py
  — stream URL + query filters; publish URL under /appengine/v1; docs/comments
examples/pokemon-agent/docs/TESTING.md
  — T2 expected paths corrected
examples/pokemon-agent/docs/DESIGN.md
  — astrate_client module row documents real endpoints
docs/handoff/pokemon-agent-memory.md   — this file
docs/handoff/pokemon-agent-handoff.md  — next-session prompt
```

### Branch
`feat/pokemon-agent`

Note: working tree may also contain **unrelated** WIP (pipelines migration, broker ACL,
triggers). Do **not** commit those with pokemon-agent work unless intentionally stacked.

### Prior session (session 1) still stands
Architecture decision (two Python services, no AtomVM/TCP bridge), full module
implementations under `examples/pokemon-agent/`, 10 ADRs, TESTING guide T0–T4.
See session-1 content below for file tree and REAL-vs-approximate notes.

---

## What's NOT done yet (next session scope)

### P1 — Run unit tests and fix failures
```sh
cd examples/pokemon-agent/emulator-agent && pip install pytest pyboy && pytest tests/ -v
cd examples/pokemon-agent/llm-orchestrator && pip install pytest pydantic pydantic-settings && pytest tests/ -v
```
Fix import/assertion failures. Commit only under `examples/pokemon-agent/`.

### P2 — Verify WRAM addresses against pret/pokered
Cross-check `emulator_agent/wram.py` vs https://github.com/pret/pokered/blob/master/wram.asm:
- Dialog buffer: `$CC2A` vs `$CF4B`
- Party: `PARTY_SLOT_SIZE=44`, `SLOT_CURRENT_HP_OFFSET=1`

### P3 — Smoke test (needs running Astrate + astartectl)
`examples/pokemon-agent/docs/TESTING.md` T1 + T2.
Use **corrected** AppEngine/stream paths above. Socket JWT must carry `a_ch`.

### P4 — ROM integration (user supplies Pokémon Red ROM)
Full loop: pyboy + emulator agent + orchestrator + LLM.

### P5 — Optional MkDocs nav entry for the example DESIGN.md

## Risks and known issues
1. ~~WebSocket path may differ from Astarte~~ **Resolved in P0** — Astrate-native socket.
2. Dialog buffer WRAM address still unverified (P2).
3. `paho-mqtt` 2.x `CallbackAPIVersion.VERSION2` required.
4. pyboy `set_emulation_speed(0)` headless behavior untested.
5. `--stub` `_StubPyboy` uses `bytearray`; real pyboy memory is `memoryview`.
6. App token for the orchestrator must include **channels** claim (`a_ch`), not only
   AppEngine data claims — otherwise the socket returns 401.
7. Unrelated WIP on the same branch: leave alone unless the user says otherwise.

## Architecture (do not re-discuss)
```
pyboy (in-process) ↔ Emulator Agent (Python) ←MQTT mTLS→ Astrate ←WS/HTTP→ LLM Orchestrator
```
No AtomVM. No TCP bridge. Two services. Astrate is the bus.

## Session 1 file tree (still accurate)
```
examples/pokemon-agent/
├── README.md, Makefile
├── astrate-interfaces/   (3 JSON interfaces)
├── emulator-agent/       (pyboy + MQTT mTLS)
├── llm-orchestrator/     (App API + OpenAI-compatible LLM)
└── docs/ DESIGN.md DECISIONS.md TESTING.md
```
