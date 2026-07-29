"""
Assembles LLM prompts from current game state.

The system prompt defines the agent's role and output contract once.
The user prompt is rebuilt every turn and contains all game-specific context.
"""

from __future__ import annotations

AVAILABLE_ACTIONS_OVERWORLD = ["UP", "DOWN", "LEFT", "RIGHT", "A", "B", "START"]
AVAILABLE_ACTIONS_BATTLE    = ["UP", "DOWN", "A", "B"]


def is_actionable_dialog(text: str) -> bool:
    """Filter wStringBuffer residue (name entry) that is not a real textbox.

    Mirrors emulator_agent.state_decoder.is_actionable_dialog so the
    orchestrator stays safe even if an older emulator build still publishes
    raw buffer text.
    """
    t = (text or "").strip()
    if not t:
        return False
    if " " not in t and t.isalpha() and t.isupper() and len(t) <= 10:
        return False
    return True


def build_system_prompt() -> str:
    return """\
You are an autonomous agent playing Pokémon Red on a Game Boy emulator.

Each turn you receive the current game state as a JSON object and must respond
with exactly ONE controller action as a JSON object.

OUTPUT FORMAT (respond with ONLY valid JSON, no explanation, no markdown):
{"button": "<BUTTON>", "holdFrames": <integer>}

BUTTON must be one of: UP, DOWN, LEFT, RIGHT, A, B, START, SELECT, NONE
holdFrames is how many frames to hold the button:
  1   = single tap (menus, confirm dialog)
  16  = one tile step (overworld movement — Gen 1 walking needs ~16 frames)
  32  = two tile steps
  1–4 = battle menu navigation

RULES:
- Progress the story: explore new areas, talk to NPCs, battle trainers, earn badges.
- In overworld ALWAYS use holdFrames=16 (or 32) for UP/DOWN/LEFT/RIGHT — short holds fail.
- In battle: UP/DOWN navigate the move menu; A confirms; B goes back.
- If stasis=true you MUST choose a DIFFERENT direction from the last move or use A/B.
- Avoid repeating the same action more than 3 times consecutively unless in a menu.
- If dialog_text is non-empty, press A to advance it.
- Starting in Red's House 2F: stairs down are toward the top-right (increase X, decrease Y).
- NONE is a valid no-op when you need to wait (e.g., during animations).\
"""


def build_user_prompt(
    state: dict,
    party: list[dict],
    action_history: list[str],
    stasis: bool,
) -> str:
    """Build the per-turn user prompt from current game state.

    Args:
        state:          GameState dict from Astrate (keys match interface mapping names).
        party:          List of PartyStatus dicts, one per slot, latest known values.
        action_history: Last N button names pressed (oldest first).
        stasis:         True when the stasis flag is set in the latest GameState.
    """
    lines: list[str] = []

    # --- Location ---
    map_name = state.get("mapName", "Unknown")
    x, y     = state.get("playerX", "?"), state.get("playerY", "?")
    lines.append(f"LOCATION: {map_name}  (tile X={x}, Y={y})")

    # --- Battle state ---
    in_battle    = state.get("inBattle", False)
    battle_type  = state.get("battleType", 0)
    if in_battle:
        btype_str = {1: "WILD", 2: "TRAINER"}.get(battle_type, "UNKNOWN")
        lines.append(f"BATTLE: {btype_str} battle in progress")
        actions = AVAILABLE_ACTIONS_BATTLE
    else:
        lines.append("BATTLE: None (overworld)")
        actions = AVAILABLE_ACTIONS_OVERWORLD

    # --- Dialog ---
    dialog = state.get("dialogText", "").strip()
    if dialog and is_actionable_dialog(dialog):
        lines.append(f'DIALOG: "{dialog}"')
        lines.append("  → Press A to advance dialog.")
    elif dialog:
        # Residue still present in older GameState samples — do not instruct A-mash.
        lines.append("DIALOG: (none — buffer residue ignored)")

    # --- Party ---
    if party:
        lines.append("PARTY:")
        for m in party:
            name    = m.get("name", "?")
            cur_hp  = m.get("currentHp", "?")
            max_hp  = m.get("maxHp", "?")
            level   = m.get("level", "?")
            pct     = f"{100 * cur_hp // max_hp}%" if isinstance(cur_hp, int) and isinstance(max_hp, int) and max_hp > 0 else "?"
            lines.append(f"  Slot {m.get('slotIndex', '?')}: {name} Lv{level}  HP {cur_hp}/{max_hp} ({pct})")
    else:
        lines.append("PARTY: (no data yet)")

    # --- Action history ---
    if action_history:
        lines.append(f"LAST ACTIONS: {' → '.join(action_history)}")
    else:
        lines.append("LAST ACTIONS: (none)")

    # --- Stasis warning ---
    if stasis:
        last = action_history[-1] if action_history else "?"
        lines.append(
            f"⚠ STASIS: Player has not moved for many frames. "
            f"Do NOT press {last!r} again. Choose a different direction or interact."
        )

    # --- Available actions ---
    lines.append(f"AVAILABLE BUTTONS: {', '.join(actions)}")
    lines.append("")
    lines.append("Respond with exactly: {\"button\": \"<BUTTON>\", \"holdFrames\": <int>}")

    return "\n".join(lines)
