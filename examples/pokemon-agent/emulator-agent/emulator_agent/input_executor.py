"""
Executes ControlCommand payloads by injecting button events into pyboy.

MQTT callbacks only enqueue commands. The main loop drains the queue and owns
every pyboy.tick() — InputExecutor never ticks the emulator (avoids races).

sequenceId monotonicity check prevents MQTT redelivery from firing duplicate inputs.
Local intro-skip presses use _local=True and skip sequenceId accounting.
"""

from __future__ import annotations

import logging
import queue
from typing import Optional

log = logging.getLogger(__name__)

# Map button name → (press_event_name, release_event_name)
# pyboy.WindowEvent members are used as strings so we don't hard-import pyboy
# at module level (allows unit-testing without pyboy installed).
_BUTTON_EVENTS: dict[str, tuple[str, str]] = {
    "UP":     ("PRESS_ARROW_UP",     "RELEASE_ARROW_UP"),
    "DOWN":   ("PRESS_ARROW_DOWN",   "RELEASE_ARROW_DOWN"),
    "LEFT":   ("PRESS_ARROW_LEFT",   "RELEASE_ARROW_LEFT"),
    "RIGHT":  ("PRESS_ARROW_RIGHT",  "RELEASE_ARROW_RIGHT"),
    "A":      ("PRESS_BUTTON_A",     "RELEASE_BUTTON_A"),
    "B":      ("PRESS_BUTTON_B",     "RELEASE_BUTTON_B"),
    "START":  ("PRESS_BUTTON_START", "RELEASE_BUTTON_START"),
    "SELECT": ("PRESS_BUTTON_SELECT", "RELEASE_BUTTON_SELECT"),
}

MAX_HOLD_FRAMES = 120  # safety cap: ~2 seconds at 60 fps


class InputExecutor:
    """Translates ControlCommand dicts into pyboy input events.

    Threading model:
      - MQTT thread → enqueue() only (no pyboy calls)
      - Main loop → before_tick() / after_tick() around pyboy.tick()

    Args:
        pyboy: A PyBoy instance (or any mock implementing send_input).
               Pass None in stub mode to skip actual input.
    """

    def __init__(self, pyboy) -> None:
        self._pyboy = pyboy
        self._last_sequence_id: int = -1
        self._queue: queue.Queue[dict] = queue.Queue()
        # (release_event_name, frames_remaining) while a button is held
        self._hold: Optional[tuple[str, int]] = None

    # ------------------------------------------------------------------
    # MQTT / producer side
    # ------------------------------------------------------------------

    def enqueue(self, command: dict) -> None:
        """Queue a ControlCommand from any thread (MQTT callback safe)."""
        self._queue.put(command)

    def enqueue_local(self, button: str, hold_frames: int = 4) -> bool:
        """Queue a main-loop-only press (intro skip). Skips sequenceId dedup.

        Returns False if a hold is active or MQTT/local work is already queued
        (caller should try again on a later frame).
        """
        if self.busy:
            return False
        self._queue.put({
            "button": button,
            "holdFrames": hold_frames,
            "_local": True,
        })
        return True

    @property
    def busy(self) -> bool:
        """True while holding a button or with at least one queued command."""
        return self._hold is not None or not self._queue.empty()

    @property
    def last_sequence_id(self) -> int:
        return self._last_sequence_id

    # ------------------------------------------------------------------
    # Main-loop side (single-threaded with pyboy)
    # ------------------------------------------------------------------

    def before_tick(self) -> None:
        """Apply press for a newly started command. Call before pyboy.tick()."""
        if self._hold is not None:
            return
        while self._hold is None:
            try:
                command = self._queue.get_nowait()
            except queue.Empty:
                return
            self._begin(command)

    def after_tick(self) -> None:
        """Count down hold frames and release. Call after pyboy.tick()."""
        if self._hold is None:
            return
        release_name, left = self._hold
        left -= 1
        if left <= 0:
            self._send(release_name)
            self._hold = None
        else:
            self._hold = (release_name, left)

    # ------------------------------------------------------------------
    # Internals
    # ------------------------------------------------------------------

    def _begin(self, command: dict) -> None:
        """Start one command: optional sequenceId check, press, arm hold."""
        is_local = bool(command.get("_local"))
        if not is_local:
            seq_id = int(command.get("sequenceId", 0))
            if seq_id <= self._last_sequence_id:
                log.debug(
                    "Deduplicated command sequenceId=%d (last=%d)",
                    seq_id, self._last_sequence_id,
                )
                return
            self._last_sequence_id = seq_id
        else:
            seq_id = -1

        button = str(command.get("button", "NONE")).upper()
        hold_frames = max(1, min(int(command.get("holdFrames", 1)), MAX_HOLD_FRAMES))

        if button == "NONE" or button not in _BUTTON_EVENTS:
            if button != "NONE":
                log.warning("Unknown button %r — treating as NONE", button)
            log.debug("NONE command (seq=%s local=%s)", seq_id, is_local)
            return

        src = "local" if is_local else f"seq={seq_id}"
        log.info("Input: %s ×%d frames (%s)", button, hold_frames, src)

        if self._pyboy is None:
            # Stub mode: accept the command without touching hardware
            return

        press_name, release_name = _BUTTON_EVENTS[button]
        self._send(press_name)
        self._hold = (release_name, hold_frames)

    def _send(self, event_name: str) -> None:
        """Resolve a WindowEvent by name and dispatch it."""
        if self._pyboy is None:
            return
        from pyboy.utils import WindowEvent  # local import avoids hard dep at module level
        event = getattr(WindowEvent, event_name)
        self._pyboy.send_input(event)
