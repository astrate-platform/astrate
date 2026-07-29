"""
Executes ControlCommand payloads by injecting button events into pyboy.

Button → WindowEvent mapping covers all 8 Game Boy buttons.
sequenceId monotonicity check prevents MQTT redelivery from firing duplicate inputs.
"""

from __future__ import annotations

import logging
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
    "SELECT": ("PRESS_BUTTON_SELECT","RELEASE_BUTTON_SELECT"),
}

MAX_HOLD_FRAMES = 120  # safety cap: ~2 seconds at 60 fps


class InputExecutor:
    """Translates ControlCommand dicts into pyboy input events.

    Args:
        pyboy: A PyBoy instance (or any mock implementing send_input/tick).
               Pass None in stub mode to skip actual input.
    """

    def __init__(self, pyboy) -> None:
        self._pyboy = pyboy
        self._last_sequence_id: int = -1

    def execute(self, command: dict) -> bool:
        """Execute a ControlCommand.

        Returns:
            True  — command was executed (or was a valid NONE no-op).
            False — command was deduplicated (sequenceId already seen).
        """
        seq_id = int(command.get("sequenceId", 0))
        if seq_id <= self._last_sequence_id:
            log.debug("Deduplicated command sequenceId=%d (last=%d)",
                      seq_id, self._last_sequence_id)
            return False
        self._last_sequence_id = seq_id

        button = str(command.get("button", "NONE")).upper()
        hold_frames = max(1, min(int(command.get("holdFrames", 1)), MAX_HOLD_FRAMES))

        if button == "NONE" or button not in _BUTTON_EVENTS:
            if button != "NONE":
                log.warning("Unknown button %r — treating as NONE", button)
            log.debug("NONE command (seq=%d)", seq_id)
            return True

        log.info("Input: %s ×%d frames (seq=%d)", button, hold_frames, seq_id)

        if self._pyboy is None:
            # Stub mode: log but don't touch pyboy
            return True

        press_name, release_name = _BUTTON_EVENTS[button]
        self._send(press_name)
        for _ in range(hold_frames):
            self._pyboy.tick()
        self._send(release_name)
        return True

    def _send(self, event_name: str) -> None:
        """Resolve a WindowEvent by name and dispatch it."""
        from pyboy.utils import WindowEvent  # local import avoids hard dep at module level
        event = getattr(WindowEvent, event_name)
        self._pyboy.send_input(event)
