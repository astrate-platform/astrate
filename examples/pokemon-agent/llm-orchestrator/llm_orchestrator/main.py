"""
LLM Orchestrator entry point.

Subscribes to Astrate App API (GameState + PartyStatus streams), assembles
an LLM prompt each turn, calls the configured LLM, and dispatches a
ControlCommand back through Astrate.

Configuration is read from environment variables (see config.py / --help).

Usage:
    python -m llm_orchestrator.main
    python -m llm_orchestrator.main --help
"""

from __future__ import annotations

import asyncio
import collections
import logging
import signal
from typing import Optional

from .action_translator import ActionTranslator
from .astrate_client import AstrateAppClient
from .config import OrchestratorConfig
from .context_builder import build_system_prompt, build_user_prompt
from .llm_engine import LLMEngine, LLMParseError, LLMTimeoutError

log = logging.getLogger(__name__)

IFACE_GAME_STATE   = "org.pokemon.emulator.GameState"
IFACE_PARTY_STATUS = "org.pokemon.emulator.PartyStatus"


async def run(config: OrchestratorConfig) -> None:
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)-8s %(name)s: %(message)s",
    )

    astrate    = AstrateAppClient(config)
    llm        = LLMEngine(config)
    translator = ActionTranslator()

    system_prompt  = build_system_prompt()
    action_history: collections.deque[str] = collections.deque(maxlen=config.action_history_len)

    # Party cache: slot_index → latest PartyStatus dict
    latest_party: dict[int, dict] = {}

    running = True

    def _shutdown(sig, frame):
        nonlocal running
        log.info("Received %s — shutting down…", sig)
        running = False

    signal.signal(signal.SIGINT,  _shutdown)
    signal.signal(signal.SIGTERM, _shutdown)

    # ------------------------------------------------------------------
    # Task: consume PartyStatus stream and update cache
    # ------------------------------------------------------------------
    async def consume_party() -> None:
        async for event in astrate.stream_device_events(IFACE_PARTY_STATUS):
            if not running:
                break
            try:
                path  = event.get("path", "")          # e.g. "/0/currentHp"
                value = event.get("value")
                parts = path.strip("/").split("/")       # ["0", "currentHp"]
                if len(parts) >= 2:
                    slot_idx = int(parts[0])
                    field    = parts[1]
                    slot = latest_party.setdefault(slot_idx, {"slotIndex": slot_idx})
                    slot[field] = value
                    log.debug("Party slot %d updated: %s=%s", slot_idx, field, value)
            except (ValueError, KeyError) as exc:
                log.warning("Malformed PartyStatus event: %s (%s)", event, exc)

    # ------------------------------------------------------------------
    # Task: consume GameState stream and drive LLM
    # ------------------------------------------------------------------
    async def consume_game_state() -> None:
        async for event in astrate.stream_device_events(IFACE_GAME_STATE):
            if not running:
                break
            try:
                state = event.get("value", {})
                if not state:
                    continue

                stasis       = bool(state.get("stasis", False))
                party_list   = list(latest_party.values())
                history_list = list(action_history)

                user_prompt = build_user_prompt(state, party_list, history_list, stasis)
                log.debug("User prompt:\n%s", user_prompt)

                try:
                    llm_output = await llm.infer(system_prompt, user_prompt)
                    command    = translator.translate(llm_output)

                    log.info("LLM → %s ×%d (seq=%d)",
                             command["button"], command["holdFrames"], command["sequenceId"])

                    await astrate.publish_command("/command", command)
                    action_history.append(command["button"])

                except LLMTimeoutError:
                    log.warning("LLM timeout — sending NONE no-op command.")
                    noop = translator.translate({"button": "NONE", "holdFrames": 1})
                    await astrate.publish_command("/command", noop)

                except LLMParseError as exc:
                    log.error("LLM parse error (all retries failed): %s — skipping turn.", exc)

                except ValueError as exc:
                    log.error("ActionTranslator rejected LLM output: %s", exc)

            except Exception as exc:  # noqa: BLE001
                log.exception("Unexpected error in GameState consumer: %s", exc)

    log.info("LLM Orchestrator running. Press Ctrl-C to stop.")

    try:
        async with asyncio.TaskGroup() as tg:
            tg.create_task(consume_party(),      name="party-consumer")
            tg.create_task(consume_game_state(), name="state-consumer")
    finally:
        await astrate.aclose()
        await llm.aclose()
        log.info("LLM Orchestrator stopped.")


def main() -> None:
    config = OrchestratorConfig()  # reads from POKEMON_* env vars
    asyncio.run(run(config))


if __name__ == "__main__":
    main()
