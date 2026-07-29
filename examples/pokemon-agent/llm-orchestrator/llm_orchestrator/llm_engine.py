"""
OpenAI-compatible LLM inference engine with timeout and retry.

Uses httpx.AsyncClient directly against the /chat/completions endpoint so
it works with any OpenAI-compatible provider (OpenAI, Groq, Ollama, LM Studio, etc.)
without importing the openai SDK.
"""

from __future__ import annotations

import asyncio
import json
import logging
from typing import Any

import httpx

log = logging.getLogger(__name__)


class LLMTimeoutError(Exception):
    """Raised when the LLM does not respond within the configured timeout."""


class LLMParseError(Exception):
    """Raised when all retries fail to produce a valid JSON action dict."""


class LLMEngine:
    """Async LLM inference engine.

    Args:
        config: OrchestratorConfig instance with openai_* fields populated.
    """

    def __init__(self, config) -> None:
        self._api_base  = config.openai_api_base.rstrip("/")
        self._api_key   = config.openai_api_key
        self._model     = config.openai_model
        self._timeout   = config.llm_timeout_seconds
        self._max_retry = config.llm_max_retries

        self._http = httpx.AsyncClient(
            headers={
                "Authorization": f"Bearer {self._api_key}",
                "Content-Type":  "application/json",
            },
            timeout=httpx.Timeout(self._timeout + 1),  # slightly looser than our own timeout
        )

    async def infer(self, system: str, user: str) -> dict:
        """Call the LLM and return a parsed action dict.

        Retries up to config.llm_max_retries times on JSON parse errors.
        Raises:
            LLMTimeoutError: if the request exceeds config.llm_timeout_seconds.
            LLMParseError:   if all retries produce non-parseable output.
        """
        last_exc: Exception | None = None

        for attempt in range(1, self._max_retry + 1):
            try:
                raw = await asyncio.wait_for(
                    self._call(system, user),
                    timeout=self._timeout,
                )
                result = self._parse(raw)
                if attempt > 1:
                    log.info("LLM succeeded on attempt %d", attempt)
                return result

            except asyncio.TimeoutError:
                log.warning("LLM timed out after %.1fs", self._timeout)
                raise LLMTimeoutError(f"LLM did not respond within {self._timeout}s")

            except LLMParseError as exc:
                log.warning("LLM parse error (attempt %d/%d): %s — raw=%r",
                             attempt, self._max_retry, exc, raw if 'raw' in dir() else "?")
                last_exc = exc
                continue

            except httpx.HTTPStatusError as exc:
                log.error("LLM HTTP error %d: %s", exc.response.status_code, exc.response.text)
                last_exc = exc
                continue

        raise LLMParseError(
            f"All {self._max_retry} LLM attempts failed"
        ) from last_exc

    async def _call(self, system: str, user: str) -> str:
        """Make the raw HTTP call and return the assistant message content."""
        body: dict[str, Any] = {
            "model": self._model,
            "messages": [
                {"role": "system", "content": system},
                {"role": "user",   "content": user},
            ],
            "temperature": 0.2,   # low temperature for consistent JSON output
            "max_tokens":  64,    # {"button":"RIGHT","holdFrames":8} is ~35 tokens
            "response_format": {"type": "json_object"},  # supported by gpt-4o, gpt-4-turbo
        }
        resp = await self._http.post(f"{self._api_base}/chat/completions", json=body)
        resp.raise_for_status()
        data = resp.json()
        return data["choices"][0]["message"]["content"]

    @staticmethod
    def _parse(raw: str) -> dict:
        """Parse LLM output into a validated action dict.

        Raises LLMParseError on invalid JSON or missing required keys.
        """
        try:
            obj = json.loads(raw)
        except json.JSONDecodeError as exc:
            raise LLMParseError(f"Not valid JSON: {exc}") from exc

        if not isinstance(obj, dict):
            raise LLMParseError(f"Expected a JSON object, got {type(obj).__name__}")

        missing = {"button", "holdFrames"} - obj.keys()
        if missing:
            raise LLMParseError(f"Missing keys: {missing}")

        return obj

    async def aclose(self) -> None:
        await self._http.aclose()
