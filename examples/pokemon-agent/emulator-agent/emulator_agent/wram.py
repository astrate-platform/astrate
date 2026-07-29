"""
WRAM address constants and memory helpers for Pokémon Red.

Address values verified against the pret/pokered disassembly:
https://github.com/pret/pokered/blob/master/wram.asm
"""

from __future__ import annotations

# ---------------------------------------------------------------------------
# WRAM addresses
# ---------------------------------------------------------------------------

MAP_ID            = 0xD35E  # Current map index
PLAYER_X          = 0xD362  # Player tile X (column)
PLAYER_Y          = 0xD361  # Player tile Y (row)
PARTY_COUNT       = 0xD163  # Number of Pokémon in party (0–6)

# Party species IDs: slot n → PARTY_SPECIES_BASE + n  (n = 0..5)
# pret: wPartySpecies (list ends with $FF at D16A)
PARTY_SPECIES_BASE = 0xD164

# Party mon structs: pret party_struct / PARTYMON_STRUCT_LENGTH = $2C (44).
# Slot 0 base = wPartyMon1 = $D16B. Per-slot layout (macros/ram.asm):
#   +0  Species (1)
#   +1  HP current (2, big-endian)     → $D16C
#   +3  BoxLevel (1)  — NOT max HP
#   ...
#   +33 Level (1)                      → $D18C
#   +34 MaxHP (2, big-endian)          → $D18D
PARTY_DATA_BASE    = 0xD16B
PARTY_SLOT_SIZE    = 44          # 0x2C = PARTYMON_STRUCT_LENGTH

# Offsets within each party slot (relative to PARTY_DATA_BASE + slot*PARTY_SLOT_SIZE)
SLOT_CURRENT_HP_OFFSET = 1    # 2 bytes, big-endian
SLOT_MAX_HP_OFFSET     = 34   # 2 bytes, big-endian (0x22) — was wrongly 3 (BoxLevel)
SLOT_LEVEL_OFFSET      = 33   # 1 byte (0x21)

BATTLE_TYPE        = 0xD057  # pret wIsInBattle: 0 overworld, 1 wild, 2 trainer
# pret wStringBuffer (historical label wcf4b). NOT $CC2A (previous menu item id).
DIALOG_BUFFER      = 0xCF4B
DIALOG_LENGTH      = 20       # NAME_BUFFER_LENGTH; string terminator 0x50 '@'


# ---------------------------------------------------------------------------
# Map name table (pret/pokered indices)
# ---------------------------------------------------------------------------

MAP_NAMES: dict[int, str] = {
    0x00: "Pallet Town",
    0x01: "Viridian City",
    0x02: "Pewter City",
    0x03: "Cerulean City",
    0x04: "Lavender Town",
    0x05: "Vermilion City",
    0x06: "Celadon City",
    0x07: "Fuchsia City",
    0x08: "Cinnabar Island",
    0x09: "Indigo Plateau",
    0x0A: "Saffron City",
    0x0C: "Route 1",
    0x0D: "Route 2",
    0x0E: "Route 3",
    0x0F: "Route 4",
    0x10: "Route 5",
    0x11: "Route 6",
    0x12: "Route 7",
    0x13: "Route 8",
    0x14: "Route 9",
    0x15: "Route 10",
    0x16: "Route 11",
    0x17: "Route 12",
    0x18: "Route 13",
    0x19: "Route 14",
    0x1A: "Route 15",
    0x1B: "Route 16",
    0x1C: "Route 17",
    0x1D: "Route 18",
    0x1E: "Route 19",
    0x1F: "Route 20",
    0x20: "Route 21",
    0x21: "Route 22",
    0x22: "Route 23",
    0x23: "Route 24",
    0x24: "Route 25",
    0x25: "Red's House 1F",
    0x26: "Red's House 2F",
    0x27: "Blue's House",
    0x28: "Oak's Lab",
    0x33: "Viridian Forest",
    0x36: "Mt. Moon B1F",
    0x37: "Mt. Moon B2F",
    0x50: "S.S. Anne Deck",
    0x51: "S.S. Anne Kitchen",
    0x52: "S.S. Anne Captain's Room",
    0x58: "Pokémon Tower 1F",
    0x59: "Pokémon Tower 2F",
    0x5A: "Pokémon Tower 3F",
    0x5B: "Pokémon Tower 4F",
    0x5C: "Pokémon Tower 5F",
    0x5D: "Pokémon Tower 6F",
    0x5E: "Pokémon Tower 7F",
    0x6C: "Silph Co. 1F",
    0x91: "Pokémon Mansion 1F",
    0xF5: "Pokémon League",
}

# ---------------------------------------------------------------------------
# Species name table (National Pokédex order, Gen 1)
# ---------------------------------------------------------------------------

SPECIES_NAMES: dict[int, str] = {
    1: "Bulbasaur",   2: "Ivysaur",    3: "Venusaur",
    4: "Charmander",  5: "Charmeleon", 6: "Charizard",
    7: "Squirtle",    8: "Wartortle",  9: "Blastoise",
    10: "Caterpie",  11: "Metapod",   12: "Butterfree",
    13: "Weedle",    14: "Kakuna",    15: "Beedrill",
    16: "Pidgey",    17: "Pidgeotto", 18: "Pidgeot",
    19: "Rattata",   20: "Raticate",
    21: "Spearow",   22: "Fearow",
    23: "Ekans",     24: "Arbok",
    25: "Pikachu",   26: "Raichu",
    27: "Sandshrew", 28: "Sandslash",
    29: "Nidoran♀",  30: "Nidorina",  31: "Nidoqueen",
    32: "Nidoran♂",  33: "Nidorino",  34: "Nidoking",
    35: "Clefairy",  36: "Clefable",
    37: "Vulpix",    38: "Ninetales",
    39: "Jigglypuff",40: "Wigglytuff",
    41: "Zubat",     42: "Golbat",
    43: "Oddish",    44: "Gloom",     45: "Vileplume",
    46: "Paras",     47: "Parasect",
    48: "Venonat",   49: "Venomoth",
    50: "Diglett",   51: "Dugtrio",
    52: "Meowth",    53: "Persian",
    54: "Psyduck",   55: "Golduck",
    56: "Mankey",    57: "Primeape",
    58: "Growlithe", 59: "Arcanine",
    60: "Poliwag",   61: "Poliwhirl", 62: "Poliwrath",
    63: "Abra",      64: "Kadabra",   65: "Alakazam",
    66: "Machop",    67: "Machoke",   68: "Machamp",
    69: "Bellsprout",70: "Weepinbell",71: "Victreebel",
    72: "Tentacool", 73: "Tentacruel",
    74: "Geodude",   75: "Graveler",  76: "Golem",
    77: "Ponyta",    78: "Rapidash",
    79: "Slowpoke",  80: "Slowbro",
    81: "Magnemite", 82: "Magneton",
    83: "Farfetch'd",
    84: "Doduo",     85: "Dodrio",
    86: "Seel",      87: "Dewgong",
    88: "Grimer",    89: "Muk",
    90: "Shellder",  91: "Cloyster",
    92: "Gastly",    93: "Haunter",   94: "Gengar",
    95: "Onix",
    96: "Drowzee",   97: "Hypno",
    98: "Krabby",    99: "Kingler",
    100: "Voltorb",  101: "Electrode",
    102: "Exeggcute",103: "Exeggutor",
    104: "Cubone",   105: "Marowak",
    106: "Hitmonlee",107: "Hitmonchan",
    108: "Lickitung",
    109: "Koffing",  110: "Weezing",
    111: "Rhyhorn",  112: "Rhydon",
    113: "Chansey",
    114: "Tangela",
    115: "Kangaskhan",
    116: "Horsea",   117: "Seadra",
    118: "Goldeen",  119: "Seaking",
    120: "Staryu",   121: "Starmie",
    122: "Mr. Mime",
    123: "Scyther",
    124: "Jynx",
    125: "Electabuzz",
    126: "Magmar",
    127: "Pinsir",
    128: "Tauros",
    129: "Magikarp", 130: "Gyarados",
    131: "Lapras",
    132: "Ditto",
    133: "Eevee",    134: "Vaporeon", 135: "Jolteon", 136: "Flareon",
    137: "Porygon",
    138: "Omanyte",  139: "Omastar",
    140: "Kabuto",   141: "Kabutops",
    142: "Aerodactyl",
    143: "Snorlax",
    144: "Articuno", 145: "Zapdos",   146: "Moltres",
    147: "Dratini",  148: "Dragonair",149: "Dragonite",
    150: "Mewtwo",
    151: "Mew",
}

# ---------------------------------------------------------------------------
# Game Boy character encoding → Python str
#
# Pokémon Red uses a custom charset (not ASCII).  The table below covers
# printable characters used in dialog boxes.  0x50 is the string terminator.
# Reference: https://bulbapedia.bulbagarden.net/wiki/Character_encoding_(Generation_I)
# ---------------------------------------------------------------------------

_CHARSET: dict[int, str] = {
    0x50: "",   # string terminator (mapped to empty so we can stop on it)
    0x7F: " ",
    0x80: "A",  0x81: "B",  0x82: "C",  0x83: "D",  0x84: "E",
    0x85: "F",  0x86: "G",  0x87: "H",  0x88: "I",  0x89: "J",
    0x8A: "K",  0x8B: "L",  0x8C: "M",  0x8D: "N",  0x8E: "O",
    0x8F: "P",  0x90: "Q",  0x91: "R",  0x92: "S",  0x93: "T",
    0x94: "U",  0x95: "V",  0x96: "W",  0x97: "X",  0x98: "Y",
    0x99: "Z",
    0x9A: "(",  0x9B: ")",  0x9C: ":",  0x9D: ";",  0x9E: "[",  0x9F: "]",
    0xA0: "a",  0xA1: "b",  0xA2: "c",  0xA3: "d",  0xA4: "e",
    0xA5: "f",  0xA6: "g",  0xA7: "h",  0xA8: "i",  0xA9: "j",
    0xAA: "k",  0xAB: "l",  0xAC: "m",  0xAD: "n",  0xAE: "o",
    0xAF: "p",  0xB0: "q",  0xB1: "r",  0xB2: "s",  0xB3: "t",
    0xB4: "u",  0xB5: "v",  0xB6: "w",  0xB7: "x",  0xB8: "y",
    0xB9: "z",
    0xF6: "0",  0xF7: "1",  0xF8: "2",  0xF9: "3",  0xFA: "4",
    0xFB: "5",  0xFC: "6",  0xFD: "7",  0xFE: "8",  0xFF: "9",
    0xE0: "'",  0xE1: "PK", 0xE2: "MN", 0xE3: "-",
    0xE6: "?",  0xE7: "!",  0xE8: ".",
}


# ---------------------------------------------------------------------------
# Memory read helpers
# ---------------------------------------------------------------------------

def read_byte(pyboy, addr: int) -> int:
    """Read a single byte from WRAM."""
    return pyboy.memory[addr]


def read_word_be(pyboy, addr: int) -> int:
    """Read a big-endian 16-bit word (used for HP values in Pokémon Red)."""
    return (pyboy.memory[addr] << 8) | pyboy.memory[addr + 1]


def read_word_le(pyboy, addr: int) -> int:
    """Read a little-endian 16-bit word."""
    return pyboy.memory[addr] | (pyboy.memory[addr + 1] << 8)


def read_bytes(pyboy, addr: int, n: int) -> bytes:
    """Read n bytes starting at addr."""
    return bytes(pyboy.memory[addr + i] for i in range(n))


def read_text(pyboy, addr: int, n: int) -> str:
    """Decode n bytes of Game Boy character encoding into a Python string.

    Stops early at the 0x50 string terminator.  Unknown bytes are rendered
    as '?' to avoid crashes on garbled WRAM during emulator startup.
    """
    chars: list[str] = []
    for i in range(n):
        val = pyboy.memory[addr + i]
        if val == 0x50:
            break
        chars.append(_CHARSET.get(val, "?"))
    return "".join(chars).strip()
