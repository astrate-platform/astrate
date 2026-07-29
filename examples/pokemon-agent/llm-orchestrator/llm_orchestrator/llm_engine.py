class LLMTimeoutError(Exception): pass
class LLMParseError(Exception): pass

class LLMEngine:
    def __init__(self, config):
        pass
    
    async def infer(self, system: str, user: str) -> dict:
        return {"button": "NONE", "holdFrames": 1}
