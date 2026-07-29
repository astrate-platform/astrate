"""Unit tests for llm_engine backend selection, opencode NDJSON parse, JSON parse."""

from __future__ import annotations

import types

import pytest

from llm_orchestrator.llm_engine import (
    LLMParseError,
    LLMEngine,
    _extract_opencode_texts,
    _opencode_model_id,
    _resolve_backend,
)


def _cfg(**kwargs):
    defaults = {
        "llm_backend": "auto",
        "openai_model": "gpt-4o",
        "openai_api_base": "https://api.openai.com/v1",
        "openai_api_key": "sk-test",
        "llm_timeout_seconds": 5.0,
        "llm_max_retries": 3,
    }
    defaults.update(kwargs)
    return types.SimpleNamespace(**defaults)


# ---------------------------------------------------------------------------
# Backend resolution
# ---------------------------------------------------------------------------


def test_resolve_backend_auto_opencode_prefix():
    assert _resolve_backend(_cfg(openai_model="opencode/big-pickle")) == "opencode"


def test_resolve_backend_auto_big_pickle_alias():
    assert _resolve_backend(_cfg(openai_model="big-pickle")) == "opencode"


def test_resolve_backend_auto_openai_default():
    assert _resolve_backend(_cfg(openai_model="gpt-4o")) == "openai"


def test_resolve_backend_explicit_opencode():
    assert (
        _resolve_backend(_cfg(llm_backend="opencode", openai_model="gpt-4o"))
        == "opencode"
    )


def test_resolve_backend_explicit_openai_overrides_prefix():
    assert (
        _resolve_backend(
            _cfg(llm_backend="openai", openai_model="opencode/big-pickle")
        )
        == "openai"
    )


def test_resolve_backend_unknown_falls_back_openai():
    assert _resolve_backend(_cfg(llm_backend="ollama")) == "openai"


def test_opencode_model_id_normalizes_alias():
    assert _opencode_model_id("big-pickle") == "opencode/big-pickle"
    assert _opencode_model_id("opencode/big-pickle") == "opencode/big-pickle"
    assert _opencode_model_id("") == "opencode/big-pickle"


# ---------------------------------------------------------------------------
# opencode NDJSON text extraction
# ---------------------------------------------------------------------------


def test_extract_opencode_texts_collects_text_parts():
    ndjson = "\n".join(
        [
            '{"type":"step_start","part":{"type":"step-start"}}',
            '{"type":"text","part":{"type":"text","text":"thinking…"}}',
            '{"type":"text","part":{"type":"text","text":"{\\"button\\":\\"DOWN\\",\\"holdFrames\\":8}"}}',
            '{"type":"step_finish","part":{"type":"step-finish"}}',
            "not-json",
            "",
        ]
    )
    texts = _extract_opencode_texts(ndjson)
    assert texts == [
        "thinking…",
        '{"button":"DOWN","holdFrames":8}',
    ]


def test_extract_opencode_texts_skips_empty_and_non_text():
    ndjson = "\n".join(
        [
            '{"type":"text","part":{"type":"text","text":"   "}}',
            '{"type":"tool","part":{"type":"tool","text":"nope"}}',
            '{"type":"text","part":{}}',
        ]
    )
    assert _extract_opencode_texts(ndjson) == []


# ---------------------------------------------------------------------------
# JSON action parse (shared by both backends)
# ---------------------------------------------------------------------------


def test_parse_plain_json():
    assert LLMEngine._parse('{"button":"A","holdFrames":1}') == {
        "button": "A",
        "holdFrames": 1,
    }


def test_parse_markdown_fence():
    raw = '```json\n{"button":"UP","holdFrames":8}\n```'
    assert LLMEngine._parse(raw)["button"] == "UP"


def test_parse_prose_with_embedded_object():
    raw = 'Sure, here you go: {"button":"LEFT","holdFrames":8} hope that helps'
    assert LLMEngine._parse(raw)["button"] == "LEFT"


def test_parse_empty_raises():
    with pytest.raises(LLMParseError):
        LLMEngine._parse("")
    with pytest.raises(LLMParseError):
        LLMEngine._parse(None)  # type: ignore[arg-type]


def test_parse_missing_keys_raises():
    with pytest.raises(LLMParseError, match="Missing keys"):
        LLMEngine._parse('{"button":"A"}')
