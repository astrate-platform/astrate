class InputExecutor:
    def __init__(self, pyboy):
        self._pyboy = pyboy
        self._last_sequence_id = -1
    
    def execute(self, command: dict) -> bool:
        seq_id = command.get("sequenceId", 0)
        if seq_id <= self._last_sequence_id:
            return False
        self._last_sequence_id = seq_id
        return True
