"""
Light guidance for early-game navigation without an LLM call.

Used when POKEMON_GUIDANCE=light (or as a documented alternative to slow
opencode turns). Still publishes ControlCommands through Astrate — only the
action chooser is local heuristics.

Red's House 2F stairs are pret tile (7,1). Red's House 1F exit is south.
"""

from __future__ import annotations

from .context_builder import MAP_REDS_HOUSE_1F, MAP_REDS_HOUSE_2F, _is_reds_house_1f, _is_reds_house_2f

OVERWORLD_HOLD = 16
STAIRS_2F = (7, 1)


def suggest_action(state: dict, action_history: list[str] | None = None) -> dict | None:
    """Return {"button", "holdFrames"} or None if no guidance for this map.

    Prefers RIGHT then UP on 2F toward (7,1). On stairs tile, keep stepping UP
    to trigger warp. On 1F, walk DOWN toward the south door.
    """
    map_name = state.get("mapName", "")
    map_id = state.get("mapId")
    x, y = state.get("playerX"), state.get("playerY")
    if not isinstance(x, int) or not isinstance(y, int):
        return None

    history = action_history or []
    stasis = bool(state.get("stasis", False))

    if _is_reds_house_2f(map_name, map_id):
        return _guide_2f(x, y, stasis, history)

    if _is_reds_house_1f(map_name, map_id):
        return _guide_1f(x, y, stasis, history)

    return None


def _guide_2f(x: int, y: int, stasis: bool, history: list[str]) -> dict:
    goal_x, goal_y = STAIRS_2F
    if x < goal_x:
        btn = "RIGHT"
    elif x > goal_x:
        btn = "LEFT"
    elif y > goal_y:
        btn = "UP"
    elif y < goal_y:
        btn = "DOWN"
    else:
        # On stairs tile — keep walking into the warp (usually UP or RIGHT).
        btn = "UP"

    if stasis and history and history[-1] == btn:
        # Bump off obstacle: try alternate axis toward goal, else reverse.
        alt = _alternate_toward_goal(btn, x, y, goal_x, goal_y)
        btn = alt

    return {"button": btn, "holdFrames": OVERWORLD_HOLD}


def _guide_1f(x: int, y: int, stasis: bool, history: list[str]) -> dict:
    # South door: walk DOWN. If stasis on same DOWN, try LEFT/RIGHT then DOWN.
    btn = "DOWN"
    if stasis and history and history[-1] == "DOWN":
        # Nudge horizontally if door not aligned (1F door is roughly mid-south).
        btn = "LEFT" if (not history or history[-1] != "LEFT") else "RIGHT"
    return {"button": btn, "holdFrames": OVERWORLD_HOLD}


def _alternate_toward_goal(
    blocked: str, x: int, y: int, goal_x: int, goal_y: int
) -> str:
    if blocked in ("LEFT", "RIGHT"):
        if y > goal_y:
            return "UP"
        if y < goal_y:
            return "DOWN"
        return "UP" if blocked == "RIGHT" else "DOWN"
    # blocked vertical
    if x < goal_x:
        return "RIGHT"
    if x > goal_x:
        return "LEFT"
    return "RIGHT" if blocked == "UP" else "LEFT"


# Re-export map constants for tests (map id sanity).
__all__ = ["suggest_action", "STAIRS_2F", "OVERWORLD_HOLD", "MAP_REDS_HOUSE_1F", "MAP_REDS_HOUSE_2F"]
