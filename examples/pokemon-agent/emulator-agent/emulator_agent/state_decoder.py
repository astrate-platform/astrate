"""
Converts raw pyboy WRAM reads into structured GameState / PartyMember dataclasses.

All decode functions are pure (no side effects) and take a pyboy instance directly,
so they can be tested by injecting a mock that implements the memory[] interface.
"""

from __future__ import annotations

import dataclasses
import time
from dataclasses import dataclass

from . import wram


# ---------------------------------------------------------------------------
# Data classes
# ---------------------------------------------------------------------------

@dataclass
class PartyMember:
    slot_index: int
    name: str
    species_id: int
    current_hp: int
    max_hp: int
    level: int

    @property
    def hp_percent(self) -> float:
        """Return HP as a fraction [0.0, 1.0]; 1.0 if max_hp is 0 (safety)."""
        if self.max_hp == 0:
            return 1.0
        return self.current_hp / self.max_hp


@dataclass
class GameState:
    map_id: int
    map_name: str
    player_x: int
    player_y: int
    in_battle: bool
    battle_type: int      # 0 = overworld, 1 = wild battle, 2 = trainer battle
    dialog_text: str
    stasis: bool          # True when stuck in same position for STASIS_THRESHOLD events
    party: list[PartyMember]
    timestamp: float      # time.time() at decode


# ---------------------------------------------------------------------------
# Decode helpers
# ---------------------------------------------------------------------------

def _decode_party(pyboy) -> list[PartyMember]:
    """Read all party members from WRAM."""
    count = min(wram.read_byte(pyboy, wram.PARTY_COUNT), 6)
    members: list[PartyMember] = []
    for i in range(count):
        base = wram.PARTY_DATA_BASE + i * wram.PARTY_SLOT_SIZE
        species_id = wram.read_byte(pyboy, wram.PARTY_SPECIES_BASE + i)
        current_hp = wram.read_word_be(pyboy, base + wram.SLOT_CURRENT_HP_OFFSET)
        max_hp     = wram.read_word_be(pyboy, base + wram.SLOT_MAX_HP_OFFSET)
        level      = wram.read_byte(pyboy, base + wram.SLOT_LEVEL_OFFSET)
        name       = wram.SPECIES_NAMES.get(species_id, f"#{species_id}")
        members.append(PartyMember(
            slot_index=i,
            name=name,
            species_id=species_id,
            current_hp=current_hp,
            max_hp=max_hp,
            level=level,
        ))
    return members


def decode_state(pyboy) -> GameState:
    """Read WRAM and return a fully populated GameState.

    The `stasis` field is always False here; it is set by the main loop
    which tracks position history across multiple decode calls.
    """
    map_id      = wram.read_byte(pyboy, wram.MAP_ID)
    battle_type = wram.read_byte(pyboy, wram.BATTLE_TYPE)
    return GameState(
        map_id      = map_id,
        map_name    = wram.MAP_NAMES.get(map_id, f"Map 0x{map_id:02X}"),
        player_x    = wram.read_byte(pyboy, wram.PLAYER_X),
        player_y    = wram.read_byte(pyboy, wram.PLAYER_Y),
        battle_type = battle_type,
        in_battle   = battle_type > 0,
        dialog_text = wram.read_text(pyboy, wram.DIALOG_BUFFER, wram.DIALOG_LENGTH),
        stasis      = False,  # set by caller
        party       = _decode_party(pyboy),
        timestamp   = time.time(),
    )


# ---------------------------------------------------------------------------
# Change detection helpers (used to decide what to publish)
# ---------------------------------------------------------------------------

def states_differ(a: GameState, b: GameState) -> bool:
    """Return True if any player-visible field changed between two states."""
    return (
        a.player_x    != b.player_x or
        a.player_y    != b.player_y or
        a.in_battle   != b.in_battle or
        a.battle_type != b.battle_type or
        a.map_id      != b.map_id or
        a.dialog_text != b.dialog_text or
        a.stasis      != b.stasis
    )


def party_differs(a: GameState, b: GameState) -> list[int]:
    """Return slot indices where HP changed between two states.

    Only compares slots present in both; new/removed party members are ignored
    (those trigger a full party publish in the main loop).
    """
    changed: list[int] = []
    for ma, mb in zip(a.party, b.party):
        if ma.current_hp != mb.current_hp or ma.max_hp != mb.max_hp:
            changed.append(ma.slot_index)
    return changed
