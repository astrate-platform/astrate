VALID_BUTTONS = frozenset(["UP", "DOWN", "LEFT", "RIGHT", "A", "B", "START", "SELECT", "NONE"])
# Gen 1 overworld: one tile step is ~16 frames; short holds often fail to move.
OVERWORLD_DEFAULT_HOLD = 16
DIRECTION_BUTTONS = frozenset(["UP", "DOWN", "LEFT", "RIGHT"])


class ActionTranslator:
    def __init__(self):
        self._sequence_id = 0

    def translate(self, llm_output: dict) -> dict:
        button = str(llm_output.get("button", "NONE")).upper()
        if button not in VALID_BUTTONS:
            raise ValueError(f"Invalid button: {button!r}")
        if "holdFrames" in llm_output and llm_output["holdFrames"] is not None:
            hold_frames = int(llm_output["holdFrames"])
        else:
            hold_frames = OVERWORLD_DEFAULT_HOLD if button in DIRECTION_BUTTONS else 1
        # Bump weak directional holds — LLM often emits 8 which is half a step.
        if button in DIRECTION_BUTTONS and hold_frames < OVERWORLD_DEFAULT_HOLD:
            hold_frames = OVERWORLD_DEFAULT_HOLD
        hold_frames = max(1, min(hold_frames, 120))
        self._sequence_id += 1
        return {
            "button": button,
            "holdFrames": hold_frames,
            "sequenceId": self._sequence_id,
        }
