"""context_builder prompt assembly + dialog residue filter."""

from llm_orchestrator.context_builder import (
    build_system_prompt,
    build_user_prompt,
    is_actionable_dialog,
)


def test_system_prompt_mentions_overworld_hold():
    p = build_system_prompt()
    assert "holdFrames" in p
    assert "16" in p
    assert "Red's House" in p


def test_is_actionable_dialog_filters_residue():
    assert is_actionable_dialog("ABBAAA") is False
    assert is_actionable_dialog("RED") is False
    assert is_actionable_dialog("") is False
    assert is_actionable_dialog("Hello there!") is True
    assert is_actionable_dialog("It's a PC.") is True


def test_user_prompt_ignores_name_entry_residue():
    state = {
        "mapName": "Red's House 2F",
        "playerX": 5,
        "playerY": 6,
        "inBattle": False,
        "battleType": 0,
        "dialogText": "ABBAAA",
        "stasis": False,
    }
    prompt = build_user_prompt(state, party=[], action_history=["A", "A"], stasis=False)
    assert "Press A to advance dialog" not in prompt
    assert "LOCATION: Red's House 2F" in prompt
    assert "AVAILABLE BUTTONS:" in prompt


def test_user_prompt_advances_real_dialog():
    state = {
        "mapName": "Pallet Town",
        "playerX": 5,
        "playerY": 6,
        "inBattle": False,
        "battleType": 0,
        "dialogText": "Welcome to the world of Pokemon!",
        "stasis": False,
    }
    prompt = build_user_prompt(state, party=[], action_history=[], stasis=False)
    assert "Press A to advance dialog" in prompt
    assert "Welcome to the world of Pokemon!" in prompt


def test_user_prompt_stasis_warning():
    state = {
        "mapName": "Red's House 2F",
        "playerX": 5,
        "playerY": 6,
        "inBattle": False,
        "battleType": 0,
        "dialogText": "",
        "stasis": True,
    }
    prompt = build_user_prompt(
        state, party=[], action_history=["RIGHT", "RIGHT"], stasis=True,
    )
    assert "STASIS" in prompt
    assert "RIGHT" in prompt
