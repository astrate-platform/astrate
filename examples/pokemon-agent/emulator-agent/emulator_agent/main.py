"""
Emulator Agent entry point.

Runs the pyboy Game Boy emulator in-process, reads WRAM every frame, publishes
telemetry to Astrate, and forwards ControlCommand payloads to pyboy joypad input.

Usage:
    python -m emulator_agent.main [options]
    python -m emulator_agent.main --stub ...   # no ROM required

Run `python -m emulator_agent.main --help` for full option list.
"""

from __future__ import annotations

import argparse
import asyncio
import dataclasses
import logging
import signal
import time
from dataclasses import dataclass
from typing import Optional

from .astrate_client import AstrateClient
from .input_executor import InputExecutor
from .state_decoder import GameState, decode_state, party_differs, states_differ

log = logging.getLogger(__name__)

STASIS_THRESHOLD  = 15    # same position for N state-change events → stasis alert
HEARTBEAT_INTERVAL = 5.0  # seconds between forced GameState publishes


# ---------------------------------------------------------------------------
# Configuration dataclass (populated from CLI args)
# ---------------------------------------------------------------------------

@dataclass
class AgentConfig:
    rom_path: str
    astrate_url: str
    astrate_realm: str
    astrate_device_id: str
    cert: str
    key: str
    ca: str
    stub: bool = False
    log_level: str = "INFO"

    # MQTT broker port override (defaults: 8883 mTLS, 1883 plain)
    mqtt_port: Optional[int] = None


# ---------------------------------------------------------------------------
# Stub pyboy (used in --stub mode so no ROM is needed)
# ---------------------------------------------------------------------------

class _StubPyboy:
    """Minimal pyboy-compatible stub that simulates a player walking east."""

    def __init__(self) -> None:
        self._tick = 0
        # Fake WRAM as a bytearray
        self.memory = bytearray(0x10000)
        self._init_memory()

    def _init_memory(self) -> None:
        from . import wram
        self.memory[wram.MAP_ID]      = 0x00   # Pallet Town
        self.memory[wram.PLAYER_X]    = 5
        self.memory[wram.PLAYER_Y]    = 5
        self.memory[wram.PARTY_COUNT] = 1
        self.memory[wram.BATTLE_TYPE] = 0
        # Party slot 0: Pikachu lv5, HP 20/20
        self.memory[wram.PARTY_SPECIES_BASE] = 25  # Pikachu
        base = wram.PARTY_DATA_BASE
        self.memory[base + 1] = 0   # current HP high byte
        self.memory[base + 2] = 20  # current HP low byte
        self.memory[base + 3] = 0   # max HP high byte
        self.memory[base + 4] = 20  # max HP low byte
        self.memory[base + 33] = 5  # level

    def tick(self) -> None:
        from . import wram
        self._tick += 1
        # Simulate player moving east every 60 frames
        if self._tick % 60 == 0:
            self.memory[wram.PLAYER_X] = (self.memory[wram.PLAYER_X] + 1) % 20

    def send_input(self, event) -> None:
        pass  # no-op in stub

    def stop(self) -> None:
        pass


# ---------------------------------------------------------------------------
# Main async run loop
# ---------------------------------------------------------------------------

async def run(config: AgentConfig) -> None:
    logging.basicConfig(
        level=getattr(logging, config.log_level.upper(), logging.INFO),
        format="%(asctime)s %(levelname)-8s %(name)s: %(message)s",
    )

    # ----- Set up pyboy (or stub) -----
    if config.stub:
        log.info("Running in STUB mode — no ROM required.")
        pyboy = _StubPyboy()
    else:
        from pyboy import PyBoy
        log.info("Loading ROM: %s", config.rom_path)
        pyboy = PyBoy(config.rom_path, window="null")
        pyboy.set_emulation_speed(0)  # uncapped speed in headless mode

    # ----- Set up Astrate client and input executor -----
    client   = AstrateClient(config)
    executor = InputExecutor(pyboy if not config.stub else None)
    client.set_command_callback(executor.execute)

    await client.connect()

    # ----- Main loop state -----
    prev_state: Optional[GameState] = None
    stasis_counter  = 0
    last_heartbeat  = time.monotonic()
    running         = True

    def _shutdown(sig, frame):
        nonlocal running
        log.info("Received signal %s — shutting down…", sig)
        running = False

    signal.signal(signal.SIGINT,  _shutdown)
    signal.signal(signal.SIGTERM, _shutdown)

    log.info("Emulator Agent running. Press Ctrl-C to stop.")

    try:
        while running:
            pyboy.tick()

            raw = decode_state(pyboy)
            now = time.monotonic()

            # ---- Stasis detection ----
            if (prev_state is not None
                    and not raw.in_battle
                    and not raw.dialog_text):
                if (raw.player_x == prev_state.player_x
                        and raw.player_y == prev_state.player_y):
                    stasis_counter += 1
                else:
                    stasis_counter = 0
            else:
                stasis_counter = 0

            state = dataclasses.replace(
                raw, stasis=(stasis_counter >= STASIS_THRESHOLD)
            )
            if stasis_counter == STASIS_THRESHOLD:
                log.warning("Stasis detected at (%d,%d) — publishing alert.",
                            state.player_x, state.player_y)

            # ---- Decide what to publish ----
            changed         = prev_state is None or states_differ(state, prev_state)
            heartbeat_due   = (now - last_heartbeat) >= HEARTBEAT_INTERVAL

            if changed or heartbeat_due:
                client.publish_game_state(state)
                if heartbeat_due:
                    last_heartbeat = now

            # ---- Party HP changes ----
            if prev_state is not None:
                changed_slots = party_differs(state, prev_state)
                for slot_idx in changed_slots:
                    client.publish_party_member(state.party[slot_idx])
            elif state.party:
                # First tick: publish initial party snapshot
                for member in state.party:
                    client.publish_party_member(member)

            prev_state = state

            # Yield to the asyncio event loop so paho callbacks can run
            await asyncio.sleep(0)

    finally:
        log.info("Disconnecting…")
        await client.disconnect()
        pyboy.stop()
        log.info("Emulator Agent stopped.")


# ---------------------------------------------------------------------------
# CLI entry point
# ---------------------------------------------------------------------------

def _parse_args() -> AgentConfig:
    p = argparse.ArgumentParser(
        description="Pokémon Red Emulator Agent — publishes game state to Astrate."
    )
    p.add_argument("--rom",       default="pokemon_red.gb",
                   help="Path to Pokémon Red ROM (required unless --stub)")
    p.add_argument("--astrate-url", default="http://localhost:8080",
                   help="Astrate base URL (default: http://localhost:8080)")
    p.add_argument("--realm",     default="pokemon-dev",
                   help="Astrate realm name")
    p.add_argument("--device-id", required=True,
                   help="Astrate device ID (from pairing registration)")
    p.add_argument("--cert",      default="device.crt",
                   help="Path to device TLS certificate (PEM)")
    p.add_argument("--key",       default="device.key",
                   help="Path to device TLS private key (PEM)")
    p.add_argument("--ca",        default="ca.crt",
                   help="Path to Astrate CA certificate (PEM)")
    p.add_argument("--stub",      action="store_true",
                   help="Run without a ROM using a synthetic game state")
    p.add_argument("--log-level", default="INFO",
                   choices=["DEBUG", "INFO", "WARNING", "ERROR"],
                   help="Log verbosity (default: INFO)")
    p.add_argument("--mqtt-port", type=int, default=None,
                   help="Override MQTT broker port (default: 8883)")
    args = p.parse_args()
    return AgentConfig(
        rom_path          = args.rom,
        astrate_url       = args.astrate_url,
        astrate_realm     = args.realm,
        astrate_device_id = args.device_id,
        cert              = args.cert,
        key               = args.key,
        ca                = args.ca,
        stub              = args.stub,
        log_level         = args.log_level,
        mqtt_port         = args.mqtt_port,
    )


if __name__ == "__main__":
    asyncio.run(run(_parse_args()))
