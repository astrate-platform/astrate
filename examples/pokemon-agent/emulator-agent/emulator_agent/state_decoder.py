import dataclasses
import time
from dataclasses import dataclass
from . import wram

@dataclass
class PartyMember:
    slot_index: int
    name: str
    species_id: int
    current_hp: int
    max_hp: int
    level: int

@dataclass  
class GameState:
    map_id: int
    map_name: str
    player_x: int
    player_y: int
    in_battle: bool
    battle_type: int
    dialog_text: str
    stasis: bool
    party: list[PartyMember]
    timestamp: float

def decode_state(pyboy) -> GameState:
    map_id = wram.read_byte(pyboy, wram.MAP_ID)
    return GameState(
        map_id=map_id,
        map_name=wram.MAP_NAMES.get(map_id, "Unknown"),
        player_x=wram.read_byte(pyboy, wram.PLAYER_X),
        player_y=wram.read_byte(pyboy, wram.PLAYER_Y),
        in_battle=wram.read_byte(pyboy, wram.BATTLE_TYPE) > 0,
        battle_type=wram.read_byte(pyboy, wram.BATTLE_TYPE),
        dialog_text=wram.read_text(pyboy, wram.DIALOG_BUFFER, 20),
        stasis=False,
        party=[],
        timestamp=time.time()
    )

def states_differ(a: GameState, b: GameState) -> bool:
    return a.player_x != b.player_x or a.player_y != b.player_y or a.in_battle != b.in_battle or a.dialog_text != b.dialog_text

def party_differs(a: GameState, b: GameState) -> list[int]:
    diff = []
    for i, (m_a, m_b) in enumerate(zip(a.party, b.party)):
        if m_a.current_hp != m_b.current_hp:
            diff.append(i)
    return diff
