"""WRAM constants verified against pret/pokered (ram/wram.asm, macros/ram.asm)."""

from emulator_agent import wram


def test_core_addresses():
    assert wram.MAP_ID == 0xD35E          # wCurMap
    assert wram.PLAYER_Y == 0xD361        # wYCoord
    assert wram.PLAYER_X == 0xD362        # wXCoord
    assert wram.PARTY_COUNT == 0xD163     # wPartyCount
    assert wram.PARTY_SPECIES_BASE == 0xD164
    assert wram.PARTY_DATA_BASE == 0xD16B  # wPartyMon1
    assert wram.BATTLE_TYPE == 0xD057     # wIsInBattle


def test_party_struct_layout():
    """party_struct: species@0, HP@1, level@33, maxHP@34; size $2C."""
    assert wram.PARTY_SLOT_SIZE == 44
    assert wram.SLOT_CURRENT_HP_OFFSET == 1
    assert wram.SLOT_LEVEL_OFFSET == 33
    assert wram.SLOT_MAX_HP_OFFSET == 34
    # Absolute slot-0 addresses used by community RAM maps / PWhiddy baselines
    assert wram.PARTY_DATA_BASE + wram.SLOT_CURRENT_HP_OFFSET == 0xD16C
    assert wram.PARTY_DATA_BASE + wram.SLOT_LEVEL_OFFSET == 0xD18C
    assert wram.PARTY_DATA_BASE + wram.SLOT_MAX_HP_OFFSET == 0xD18D


def test_dialog_buffer_is_cf4b_not_cc2a():
    """$CC2A is previous menu item; string buffer is pret wStringBuffer / wcf4b."""
    assert wram.DIALOG_BUFFER == 0xCF4B
    assert wram.DIALOG_BUFFER != 0xCC2A
    assert wram.DIALOG_LENGTH == 20


class _Mem:
    """Minimal pyboy.memory stand-in for pure WRAM helper tests."""

    def __init__(self, size: int = 0x10000) -> None:
        self.memory = bytearray(size)


def test_read_text_zeroed_buffer_is_empty():
    """Title-screen / cold boot WRAM is zeros — must not become '????…'."""
    m = _Mem()
    assert wram.read_text(m, wram.DIALOG_BUFFER, wram.DIALOG_LENGTH) == ""


def test_read_text_stops_at_terminator_and_decodes_charset():
    m = _Mem()
    # Gen I: H e l l o @
    for i, b in enumerate((0x87, 0xA4, 0xAB, 0xAB, 0xAE, 0x50)):
        m.memory[wram.DIALOG_BUFFER + i] = b
    assert wram.read_text(m, wram.DIALOG_BUFFER, wram.DIALOG_LENGTH) == "Hello"
