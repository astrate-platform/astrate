VALID_BUTTONS = frozenset(["UP", "DOWN", "LEFT", "RIGHT", "A", "B", "START", "SELECT", "NONE"])

class ActionTranslator:
    def __init__(self):
        self._sequence_id = 0
    
    def translate(self, llm_output: dict) -> dict:
        button = str(llm_output.get("button", "NONE")).upper()
        if button not in VALID_BUTTONS:
            raise ValueError(f"Invalid button: {button!r}")
        hold_frames = max(1, min(int(llm_output.get("holdFrames", 1)), 120))
        self._sequence_id += 1
        return {
            "button": button,
            "holdFrames": hold_frames,
            "sequenceId": self._sequence_id,
        }
