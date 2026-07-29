"""WRAM constants and helpers."""
MAP_ID = 0xD35E
PLAYER_Y = 0xD361
PLAYER_X = 0xD362
PARTY_COUNT = 0xD163
PARTY_SPECIES_BASE = 0xD164
PARTY_HP_BASE = 0xD16C
BATTLE_TYPE = 0xD057
DIALOG_BUFFER = 0xCC2A

MAP_NAMES = {
    0: "Pallet Town",
    1: "Viridian City",
    2: "Pewter City",
    12: "Route 1",
    # ... more mappings as needed
}

SPECIES_NAMES = {
    1: "Bulbasaur",
    2: "Ivysaur",
    3: "Venusaur",
    4: "Charmander",
    # ... up to 151
}

def read_byte(pyboy, addr: int) -> int:
    return pyboy.memory[addr]

def read_word_le(pyboy, addr: int) -> int:
    return pyboy.memory[addr] | (pyboy.memory[addr + 1] << 8)

def read_bytes(pyboy, addr: int, n: int) -> bytes:
    return bytes([pyboy.memory[addr + i] for i in range(n)])

def read_text(pyboy, addr: int, n: int) -> str:
    chars = []
    for i in range(n):
        val = pyboy.memory[addr + i]
        if val == 0x50:
            break
        chars.append(chr(val))
    return "".join(chars)
