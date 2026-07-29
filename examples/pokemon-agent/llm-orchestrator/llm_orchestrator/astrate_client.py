from typing import AsyncIterator

class AstrateAppClient:
    BASE_WS  = "{url}/v1/{realm}/devices/{device_id}/interfaces/{interface}"
    BASE_HTTP = "{url}/v1/{realm}/devices/{device_id}/interfaces/{interface}{path}"
    
    def __init__(self, config):
        self.config = config
    
    async def stream_device_events(self, interface: str) -> AsyncIterator[dict]:
        yield {}
    
    async def publish_command(self, path: str, payload: dict) -> None:
        pass
