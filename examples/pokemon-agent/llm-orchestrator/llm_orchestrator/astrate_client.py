"""
Astrate App API client for the LLM Orchestrator.

Uses aiohttp for WebSocket streaming and httpx for REST publishing.
Reconnects automatically with exponential backoff on disconnect.

App API endpoint conventions (Astrate = wire-compatible with Astarte):
  Stream: GET  /v1/{realm}/devices/{device_id}/interfaces/{interface}  (WebSocket upgrade)
  Publish: POST /v1/{realm}/devices/{device_id}/interfaces/{interface}{path}

References:
  https://docs.astarte-platform.org/astarte/latest/050-query_astarte.html
"""

from __future__ import annotations

import asyncio
import json
import logging
import random
from collections.abc import AsyncIterator
from typing import Optional

import aiohttp
import httpx

log = logging.getLogger(__name__)

_BACKOFF_BASE  = 1.0   # seconds
_BACKOFF_MAX   = 60.0  # seconds
_BACKOFF_JITTER = 0.1  # ±10%


def _backoff(attempt: int) -> float:
    raw = min(_BACKOFF_BASE * (2 ** attempt), _BACKOFF_MAX)
    jitter = raw * _BACKOFF_JITTER * (random.random() * 2 - 1)
    return raw + jitter


class AstrateAppClient:
    """Async client for Astrate App API.

    Args:
        config: OrchestratorConfig with astrate_* fields populated.
    """

    def __init__(self, config) -> None:
        self._base    = config.astrate_url.rstrip("/")
        self._realm   = config.astrate_realm
        self._device  = config.astrate_device_id
        self._headers = {
            "Authorization": f"Bearer {config.astrate_app_token}",
            "Content-Type":  "application/json",
        }
        self._http = httpx.AsyncClient(headers=self._headers, timeout=10.0)

    # ------------------------------------------------------------------
    # Streaming (WebSocket)
    # ------------------------------------------------------------------

    async def stream_device_events(
        self,
        interface: str,
    ) -> AsyncIterator[dict]:
        """Yield parsed event dicts from Astrate's App API WebSocket stream.

        Astrate streams events as JSON lines over a WebSocket connection.
        Each message has the shape:
          {"interface": "org.pokemon.emulator.GameState",
           "path": "/state",
           "value": {...},
           "timestamp": "..."}

        Reconnects automatically with exponential backoff on disconnect.
        """
        url = (
            f"{self._base}/v1/{self._realm}"
            f"/devices/{self._device}/interfaces/{interface}"
        ).replace("http://", "ws://").replace("https://", "wss://")

        attempt = 0
        while True:
            try:
                async with aiohttp.ClientSession(headers=self._headers) as session:
                    async with session.ws_connect(url) as ws:
                        log.info("WebSocket connected: %s", url)
                        attempt = 0  # reset backoff on successful connect
                        async for msg in ws:
                            if msg.type == aiohttp.WSMsgType.TEXT:
                                try:
                                    yield json.loads(msg.data)
                                except json.JSONDecodeError as exc:
                                    log.warning("Unparseable WS message: %s (%s)", msg.data, exc)
                            elif msg.type in (aiohttp.WSMsgType.CLOSED,
                                              aiohttp.WSMsgType.ERROR):
                                log.warning("WebSocket closed/error: %s", msg)
                                break
            except (aiohttp.ClientError, OSError) as exc:
                delay = _backoff(attempt)
                log.warning("WebSocket error (%s) — reconnecting in %.1fs", exc, delay)
                attempt += 1
                await asyncio.sleep(delay)

    # ------------------------------------------------------------------
    # Publishing (HTTP POST)
    # ------------------------------------------------------------------

    async def publish_command(
        self,
        path: str,
        payload: dict,
    ) -> None:
        """POST a server-owned ControlCommand value via Astrate App API.

        Args:
            path:    Endpoint path, e.g. "/command"
            payload: The `v` value object, e.g. {"button": "UP", "holdFrames": 8, "sequenceId": 1}
        """
        url = (
            f"{self._base}/v1/{self._realm}"
            f"/devices/{self._device}"
            f"/interfaces/org.pokemon.emulator.ControlCommand{path}"
        )
        body = {"data": payload}
        try:
            resp = await self._http.post(url, json=body)
            resp.raise_for_status()
            log.debug("→ ControlCommand %s → %s (HTTP %d)", path, payload, resp.status_code)
        except httpx.HTTPStatusError as exc:
            log.error("Failed to publish ControlCommand: HTTP %d — %s",
                      exc.response.status_code, exc.response.text)
        except httpx.RequestError as exc:
            log.error("Network error publishing ControlCommand: %s", exc)

    async def aclose(self) -> None:
        await self._http.aclose()
