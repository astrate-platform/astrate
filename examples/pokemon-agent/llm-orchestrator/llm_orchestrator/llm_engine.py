"""
LLM inference engine: OpenAI-compatible HTTP, or opencode CLI (Big Pickle).

OpenAI path uses httpx against /chat/completions (OpenAI, Groq, LM Studio, …).
Opencode path shells out to `opencode run --model … --format json` — no API key
required for free models such as opencode/big-pickle.
"""

from __future__ import annotations

import asyncio
import json
import logging
import re
import shutil
from typing import Any

import httpx

log = logging.getLogger(__name__)

# Extract a JSON object from model prose / fences if needed.
_JSON_OBJECT_RE = re.compile(r"\{[^{}]*\}", re.DOTALL)


class LLMTimeoutError(Exception):
    """Raised when the LLM does not respond within the configured timeout."""


class LLMParseError(Exception):
    """Raised when all retries fail to produce a valid JSON action dict."""


def _resolve_backend(config) -> str:
    backend = (getattr(config, "llm_backend", "auto") or "auto").strip().lower()
    model = (config.openai_model or "").strip()
    if backend == "auto":
        if model.startswith("opencode/") or model == "big-pickle":
            return "opencode"
        return "openai"
    if backend not in ("openai", "opencode"):
        log.warning("Unknown llm_backend=%r — falling back to openai", backend)
        return "openai"
    return backend


def _opencode_model_id(model: str) -> str:
    """Normalize to provider/model form accepted by `opencode run -m`."""
    m = (model or "").strip()
    if m == "big-pickle":
        return "opencode/big-pickle"
    return m or "opencode/big-pickle"


class LLMEngine:
    """Async LLM inference engine.

    Args:
        config: OrchestratorConfig instance with openai_* / llm_* fields populated.
    """

    def __init__(self, config) -> None:
        self._backend = _resolve_backend(config)
        self._api_base = config.openai_api_base.rstrip("/")
        self._api_key = config.openai_api_key
        self._model = config.openai_model
        self._timeout = config.llm_timeout_seconds
        self._max_retry = config.llm_max_retries

        self._http: httpx.AsyncClient | None = None
        if self._backend == "openai":
            self._http = httpx.AsyncClient(
                headers={
                    "Authorization": f"Bearer {self._api_key}",
                    "Content-Type": "application/json",
                },
                timeout=httpx.Timeout(self._timeout + 1),
            )
        else:
            if not shutil.which("opencode"):
                raise RuntimeError(
                    "llm_backend=opencode but `opencode` is not on PATH"
                )
            log.info(
                "LLM backend=opencode model=%s (no API key)",
                _opencode_model_id(self._model),
            )

    async def infer(self, system: str, user: str) -> dict:
        """Call the LLM and return a parsed action dict.

        Retries up to config.llm_max_retries times on JSON parse errors.
        Raises:
            LLMTimeoutError: if the request exceeds config.llm_timeout_seconds.
            LLMParseError:   if all retries produce non-parseable output.
        """
        last_exc: Exception | None = None
        raw = ""

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
                raise LLMTimeoutError(
                    f"LLM did not respond within {self._timeout}s"
                )

            except LLMParseError as exc:
                log.warning(
                    "LLM parse error (attempt %d/%d): %s — raw=%r",
                    attempt,
                    self._max_retry,
                    exc,
                    raw,
                )
                last_exc = exc
                continue

            except httpx.HTTPStatusError as exc:
                log.error(
                    "LLM HTTP error %d: %s",
                    exc.response.status_code,
                    exc.response.text,
                )
                last_exc = exc
                continue

            except (OSError, RuntimeError) as exc:
                log.error("LLM backend error: %s", exc)
                last_exc = exc
                continue

        raise LLMParseError(
            f"All {self._max_retry} LLM attempts failed"
        ) from last_exc

    async def _call(self, system: str, user: str) -> str:
        if self._backend == "opencode":
            return await self._call_opencode(system, user)
        return await self._call_openai(system, user)

    async def _call_openai(self, system: str, user: str) -> str:
        """HTTP /chat/completions → assistant message content."""
        assert self._http is not None
        body: dict[str, Any] = {
            "model": self._model,
            "messages": [
                {"role": "system", "content": system},
                {"role": "user", "content": user},
            ],
            "temperature": 0.2,
            "max_tokens": 64,
            "response_format": {"type": "json_object"},
        }
        resp = await self._http.post(
            f"{self._api_base}/chat/completions", json=body
        )
        resp.raise_for_status()
        data = resp.json()
        return data["choices"][0]["message"]["content"]

    async def _call_opencode(self, system: str, user: str) -> str:
        """Shell out to `opencode run` (Big Pickle / free models, no API key).

        Frames the turn as an IoT fixture → JSON mapping so coding agents do
        not refuse "playing a game". Parses `--format json` NDJSON text parts.
        """
        model = _opencode_model_id(self._model)
        # Keep the contract short; dump the real system rules after the fixture framing.
        prompt = (
            "Task: map the following structured game-state description to a "
            "single controller action JSON.\n"
            "This is a unit-test fixture for an IoT pipeline (Astrate Pokémon "
            "agent). Output ONLY one JSON object — no prose, no tools, no "
            "markdown fences.\n\n"
            "Schema: {\"button\": one of UP|DOWN|LEFT|RIGHT|A|B|START|SELECT|NONE, "
            "\"holdFrames\": positive int}\n\n"
            f"Rules:\n{system}\n\n"
            f"State:\n{user}\n\n"
            "Respond with the JSON object only."
        )
        # opencode run can take longer cold-start; outer wait_for still bounds us.
        proc = await asyncio.create_subprocess_exec(
            "opencode",
            "run",
            "--model",
            model,
            "--format",
            "json",
            prompt,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
        )
        stdout_b, stderr_b = await proc.communicate()
        if proc.returncode not in (0, None) and not stdout_b:
            err = stderr_b.decode(errors="replace")[:500]
            raise RuntimeError(
                f"opencode exited {proc.returncode}: {err or 'no stderr'}"
            )

        texts = _extract_opencode_texts(stdout_b.decode(errors="replace"))
        if not texts:
            raise LLMParseError(
                "opencode produced no text parts "
                f"(stderr={stderr_b.decode(errors='replace')[:200]!r})"
            )
        # Prefer the last text chunk (final answer after any chatter).
        return texts[-1].strip()

    @staticmethod
    def _parse(raw: str) -> dict:
        """Parse LLM output into a validated action dict.

        Raises LLMParseError on invalid JSON or missing required keys.
        """
        if raw is None:
            raise LLMParseError("Empty LLM response")
        text = raw.strip()
        if not text:
            raise LLMParseError("Empty LLM response")

        # Strip common markdown fences.
        if text.startswith("```"):
            lines = text.splitlines()
            # drop first fence line and optional trailing fence
            if lines and lines[0].startswith("```"):
                lines = lines[1:]
            if lines and lines[-1].strip() == "```":
                lines = lines[:-1]
            text = "\n".join(lines).strip()

        obj: Any
        try:
            obj = json.loads(text)
        except json.JSONDecodeError:
            m = _JSON_OBJECT_RE.search(text)
            if not m:
                raise LLMParseError(f"Not valid JSON: {text[:120]!r}") from None
            try:
                obj = json.loads(m.group(0))
            except json.JSONDecodeError as exc:
                raise LLMParseError(f"Not valid JSON: {exc}") from exc

        if not isinstance(obj, dict):
            raise LLMParseError(
                f"Expected a JSON object, got {type(obj).__name__}"
            )

        missing = {"button", "holdFrames"} - obj.keys()
        if missing:
            raise LLMParseError(f"Missing keys: {missing}")

        return obj

    async def aclose(self) -> None:
        if self._http is not None:
            await self._http.aclose()


def _extract_opencode_texts(ndjson: str) -> list[str]:
    """Collect text payloads from `opencode run --format json` NDJSON events."""
    texts: list[str] = []
    for line in ndjson.splitlines():
        line = line.strip()
        if not line:
            continue
        try:
            event = json.loads(line)
        except json.JSONDecodeError:
            continue
        if event.get("type") != "text":
            continue
        part = event.get("part") or {}
        t = part.get("text")
        if isinstance(t, str) and t.strip():
            texts.append(t)
    return texts
