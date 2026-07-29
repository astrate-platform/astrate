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

# Same map position for this long (overworld, no dialog) → stasis alert.
# Time-based so it stays correct at 60 fps (frame counting of 15 was for
# the pre-ROM stub loop and fired in 0.25 s at real tick rate).
STASIS_SECONDS = 15.0
HEARTBEAT_INTERVAL = 5.0  # seconds between forced GameState publishes
DEFAULT_FPS = 60          # DESIGN §3.1; 0 = uncapped (high CPU)

# Auto-press through title / Oak intro until WRAM looks "in game".
# Prefer this over loading a save-state (session preference).
SKIP_INTRO_MAX_SECONDS = 180.0
SKIP_INTRO_INTERVAL_FRAMES = 20   # press when idle, every N frames
SKIP_INTRO_HOLD_FRAMES = 4
# Cycle: mostly A (advance text / confirm), occasional START (title menu)
_SKIP_INTRO_BUTTONS = ("A", "A", "A", "START")


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
    # Skip mTLS and use plaintext :1883 (Astrate mqtt.insecure_dev_mode only)
    insecure: bool = False
    # Target emulator ticks per second (0 = uncapped). Default 60.
    fps: int = DEFAULT_FPS
    # Mash A/START until past cold-boot zero WRAM (ROM mode; ignored with --stub).
    skip_intro: bool = True


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
        # Empty dialog: first byte is the 0x50 string terminator ('@')
        self.memory[wram.DIALOG_BUFFER] = 0x50
        # Party slot 0: Pikachu lv5, HP 20/20 (pret party_struct offsets)
        self.memory[wram.PARTY_SPECIES_BASE] = 25  # Pikachu
        base = wram.PARTY_DATA_BASE
        self.memory[base + wram.SLOT_CURRENT_HP_OFFSET] = 0      # current HP high
        self.memory[base + wram.SLOT_CURRENT_HP_OFFSET + 1] = 20  # current HP low
        self.memory[base + wram.SLOT_LEVEL_OFFSET] = 5
        self.memory[base + wram.SLOT_MAX_HP_OFFSET] = 0          # max HP high
        self.memory[base + wram.SLOT_MAX_HP_OFFSET + 1] = 20     # max HP low

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


def looks_past_cold_boot(state: GameState) -> bool:
    """Heuristic: WRAM is no longer the all-zero title/cold-boot pattern.

    Pallet Town's mapId is 0, so map alone is useless. Non-zero player coords,
    party, dialog, or battle mean the game has progressed past pure boot zeros.
    """
    return (
        state.player_x != 0
        or state.player_y != 0
        or state.in_battle
        or bool(state.dialog_text)
        or len(state.party) > 0
    )


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
    # MQTT thread only enqueues; main loop drains + owns every pyboy.tick().
    client.set_command_callback(executor.enqueue)

    await client.connect()

    # ----- Main loop state -----
    prev_state: Optional[GameState] = None
    stasis_since: Optional[float] = None
    stasis_warned = False
    last_heartbeat  = time.monotonic()
    running         = True
    frame_dt = (1.0 / config.fps) if config.fps > 0 else 0.0
    next_frame_at = time.monotonic()
    frame_i = 0
    skip_intro = bool(config.skip_intro and not config.stub)
    intro_done = not skip_intro
    intro_started_at = time.monotonic()
    intro_press_i = 0

    def _shutdown(sig, frame):
        nonlocal running
        log.info("Received signal %s — shutting down…", sig)
        running = False

    signal.signal(signal.SIGINT,  _shutdown)
    signal.signal(signal.SIGTERM, _shutdown)

    if config.fps > 0:
        log.info("Emulator Agent running at %d fps. Press Ctrl-C to stop.", config.fps)
    else:
        log.info("Emulator Agent running uncapped. Press Ctrl-C to stop.")
    if skip_intro:
        log.info(
            "Intro skip ON — mashing A/START until past cold boot (max %.0fs). "
            "Use --no-skip-intro to disable.",
            SKIP_INTRO_MAX_SECONDS,
        )

    try:
        while running:
            # ---- Intro auto-press (main thread only; yields to MQTT queue) ----
            if not intro_done and not executor.busy:
                now_intro = time.monotonic()
                if (now_intro - intro_started_at) >= SKIP_INTRO_MAX_SECONDS:
                    intro_done = True
                    log.warning(
                        "Intro skip timed out after %.0fs — stopping auto-press.",
                        SKIP_INTRO_MAX_SECONDS,
                    )
                elif frame_i % SKIP_INTRO_INTERVAL_FRAMES == 0:
                    button = _SKIP_INTRO_BUTTONS[
                        intro_press_i % len(_SKIP_INTRO_BUTTONS)
                    ]
                    if executor.enqueue_local(button, SKIP_INTRO_HOLD_FRAMES):
                        intro_press_i += 1

            # MQTT/local presses applied here; main loop owns all ticks.
            executor.before_tick()
            pyboy.tick()
            executor.after_tick()
            frame_i += 1

            raw = decode_state(pyboy)
            now = time.monotonic()

            if not intro_done and looks_past_cold_boot(raw):
                intro_done = True
                log.info(
                    "Intro skip complete — past cold boot "
                    "(map=%s pos=(%d,%d) party=%d dialog=%r).",
                    raw.map_name, raw.player_x, raw.player_y,
                    len(raw.party), (raw.dialog_text[:20] if raw.dialog_text else ""),
                )

            # ---- Stasis detection (time-based, overworld + no dialog) ----
            stationary = (
                prev_state is not None
                and not raw.in_battle
                and not raw.dialog_text
                and raw.player_x == prev_state.player_x
                and raw.player_y == prev_state.player_y
            )
            if stationary:
                if stasis_since is None:
                    stasis_since = now
            else:
                stasis_since = None
                stasis_warned = False

            in_stasis = (
                stasis_since is not None
                and (now - stasis_since) >= STASIS_SECONDS
            )
            state = dataclasses.replace(raw, stasis=in_stasis)
            if in_stasis and not stasis_warned:
                log.warning("Stasis detected at (%d,%d) — publishing alert.",
                            state.player_x, state.player_y)
                stasis_warned = True

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

            # Pace to target fps (DESIGN §3.1). Always yield so MQTT callbacks run.
            if frame_dt > 0:
                next_frame_at += frame_dt
                delay = next_frame_at - time.monotonic()
                if delay > 0:
                    await asyncio.sleep(delay)
                else:
                    await asyncio.sleep(0)
                    # Avoid spiral-of-death after a long stall (e.g. GC, MQTT).
                    if delay < -0.25:
                        next_frame_at = time.monotonic()
            else:
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
                   help="Path to device TLS certificate (PEM); unused with --insecure")
    p.add_argument("--key",       default="device.key",
                   help="Path to device TLS private key (PEM); unused with --insecure")
    p.add_argument("--ca",        default="ca.crt",
                   help="Path to Astrate CA certificate (PEM); unused with --insecure")
    p.add_argument("--stub",      action="store_true",
                   help="Run without a ROM using a synthetic game state")
    p.add_argument("--insecure",  action="store_true",
                   help="Plaintext MQTT on :1883 (Astrate insecure_dev_mode only)")
    p.add_argument("--log-level", default="INFO",
                   choices=["DEBUG", "INFO", "WARNING", "ERROR"],
                   help="Log verbosity (default: INFO)")
    p.add_argument("--mqtt-port", type=int, default=None,
                   help="Override MQTT broker port (default: 8883 mTLS / 1883 --insecure)")
    p.add_argument("--fps", type=int, default=DEFAULT_FPS,
                   help=f"Emulator tick rate (default: {DEFAULT_FPS}; 0 = uncapped)")
    p.add_argument(
        "--skip-intro",
        action=argparse.BooleanOptionalAction,
        default=True,
        help="Mash A/START past title/intro until WRAM looks in-game "
             "(default: on for ROM; ignored with --stub). Use --no-skip-intro to disable.",
    )
    args = p.parse_args()
    if args.fps < 0:
        p.error("--fps must be >= 0")
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
        insecure          = args.insecure,
        fps               = args.fps,
        skip_intro        = args.skip_intro,
    )


if __name__ == "__main__":
    asyncio.run(run(_parse_args()))
