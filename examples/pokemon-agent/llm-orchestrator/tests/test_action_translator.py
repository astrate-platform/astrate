"""ActionTranslator: button validation, hold-frame floors for overworld."""

import pytest

from llm_orchestrator.action_translator import (
    OVERWORLD_DEFAULT_HOLD,
    ActionTranslator,
)


def test_translate_basic_a_tap():
    t = ActionTranslator()
    cmd = t.translate({"button": "A", "holdFrames": 1})
    assert cmd["button"] == "A"
    assert cmd["holdFrames"] == 1
    assert cmd["sequenceId"] == 1


def test_sequence_id_monotonic():
    t = ActionTranslator()
    a = t.translate({"button": "A", "holdFrames": 1})
    b = t.translate({"button": "B", "holdFrames": 1})
    assert b["sequenceId"] == a["sequenceId"] + 1


def test_invalid_button_raises():
    t = ActionTranslator()
    with pytest.raises(ValueError, match="Invalid button"):
        t.translate({"button": "X", "holdFrames": 1})


def test_direction_default_hold_when_omitted():
    t = ActionTranslator()
    cmd = t.translate({"button": "RIGHT"})
    assert cmd["holdFrames"] == OVERWORLD_DEFAULT_HOLD


def test_direction_bumps_weak_hold_to_one_tile():
    """LLM often emits holdFrames=8; Gen 1 needs ~16 for a full step."""
    t = ActionTranslator()
    cmd = t.translate({"button": "DOWN", "holdFrames": 8})
    assert cmd["holdFrames"] == OVERWORLD_DEFAULT_HOLD


def test_direction_preserves_longer_holds():
    t = ActionTranslator()
    cmd = t.translate({"button": "UP", "holdFrames": 32})
    assert cmd["holdFrames"] == 32


def test_hold_frames_capped_at_120():
    t = ActionTranslator()
    cmd = t.translate({"button": "LEFT", "holdFrames": 999})
    assert cmd["holdFrames"] == 120


def test_button_case_normalized():
    t = ActionTranslator()
    cmd = t.translate({"button": "a", "holdFrames": 2})
    assert cmd["button"] == "A"
