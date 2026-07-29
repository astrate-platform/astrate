import asyncio
import time
import argparse
import dataclasses
from .astrate_client import AstrateClient
from .input_executor import InputExecutor
from .state_decoder import decode_state, states_differ, party_differs, GameState

@dataclasses.dataclass
class AgentConfig:
    rom_path: str
    astrate_url: str
    astrate_realm: str
    astrate_device_id: str
    cert: str
    key: str
    ca: str
    stub: bool

async def run(config: AgentConfig) -> None:
    pass

if __name__ == "__main__":
    pass
