"""light_guide path heuristics for Red's House."""

from llm_orchestrator.light_guide import suggest_action, STAIRS_2F


def test_2f_prefers_right_then_up():
    state = {
        "mapName": "Red's House 2F",
        "mapId": 38,
        "playerX": 3,
        "playerY": 6,
        "stasis": False,
    }
    a = suggest_action(state, [])
    assert a == {"button": "RIGHT", "holdFrames": 16}


def test_2f_up_when_x_aligned():
    state = {
        "mapName": "Red's House 2F",
        "mapId": 38,
        "playerX": 7,
        "playerY": 4,
        "stasis": False,
    }
    a = suggest_action(state, [])
    assert a["button"] == "UP"
    assert a["holdFrames"] == 16


def test_2f_on_stairs_keeps_stepping():
    gx, gy = STAIRS_2F
    state = {
        "mapName": "Red's House 2F",
        "mapId": 38,
        "playerX": gx,
        "playerY": gy,
        "stasis": False,
    }
    a = suggest_action(state, [])
    assert a["button"] in ("UP", "RIGHT")


def test_2f_stasis_switches_axis():
    state = {
        "mapName": "Red's House 2F",
        "mapId": 38,
        "playerX": 5,
        "playerY": 6,
        "stasis": True,
    }
    a = suggest_action(state, ["RIGHT", "RIGHT"])
    # Prefer vertical alternate when RIGHT is stuck.
    assert a["button"] in ("UP", "DOWN", "LEFT")


def test_1f_walks_down():
    state = {
        "mapName": "Red's House 1F",
        "mapId": 37,
        "playerX": 7,
        "playerY": 1,
        "stasis": False,
    }
    a = suggest_action(state, [])
    assert a["button"] == "DOWN"


def test_unknown_map_returns_none():
    state = {
        "mapName": "Pallet Town",
        "mapId": 0,
        "playerX": 5,
        "playerY": 6,
        "stasis": False,
    }
    assert suggest_action(state, []) is None
