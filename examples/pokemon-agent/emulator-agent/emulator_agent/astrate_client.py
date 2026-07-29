"""
Astrate MQTT client for the Emulator Agent.

Implements the Astarte MQTT v1 protocol wire format used by Astrate:
  - Introspection: published to  <realm>/<device_id>
  - Device data:   published to  <realm>/<device_id>/<interface>/<path>
  - Server data:   subscribed at <realm>/<device_id>/<interface>/<path>

Payload format (plain-JSON profile accepted by Astrate alongside BSON):
  {"t": "<ISO8601>", "v": <value_or_object>}

References:
  https://docs.astarte-platform.org/astarte/latest/080-mqtt-v1-protocol.html
  https://docs.astarte-platform.org/astarte/latest/040-interface_schema.html
"""

from __future__ import annotations

import json
import logging
import ssl
import threading
import time
from datetime import datetime, timezone
from typing import Callable, Optional

import paho.mqtt.client as mqtt

from .state_decoder import GameState, PartyMember

log = logging.getLogger(__name__)

# Interface names — must match the installed Astrate interface definitions
IFACE_GAME_STATE      = "org.pokemon.emulator.GameState"
IFACE_PARTY_STATUS    = "org.pokemon.emulator.PartyStatus"
IFACE_CONTROL_COMMAND = "org.pokemon.emulator.ControlCommand"

# Introspection string template: <interface>:<major>:<minor>;...
_INTROSPECTION = (
    f"{IFACE_GAME_STATE}:1:0;"
    f"{IFACE_PARTY_STATUS}:1:0;"
    f"{IFACE_CONTROL_COMMAND}:1:0"
)


def _now_iso() -> str:
    return datetime.now(timezone.utc).isoformat(timespec="milliseconds")


def _wrap(value) -> bytes:
    """Wrap a value in the Astrate plain-JSON envelope."""
    return json.dumps({"t": _now_iso(), "v": value}).encode()


class AstrateClient:
    """Thread-safe paho-mqtt wrapper for the Astrate MQTT v1 protocol.

    Paho's network loop runs in a background thread (loop_start), so MQTT
    callbacks fire on that thread.  The `execute` callback registered via
    set_command_callback is therefore called from the MQTT thread — keep it
    fast and non-blocking.
    """

    def __init__(self, config) -> None:
        self._config = config
        self._realm     = config.astrate_realm
        self._device_id = config.astrate_device_id
        self._prefix    = f"{self._realm}/{self._device_id}"
        self._command_cb: Optional[Callable[[dict], None]] = None
        self._connected = threading.Event()

        self._client = mqtt.Client(
            client_id=self._device_id,
            protocol=mqtt.MQTTv311,
            callback_api_version=mqtt.CallbackAPIVersion.VERSION2,
        )
        self._client.on_connect    = self._on_connect
        self._client.on_disconnect = self._on_disconnect
        self._client.on_message    = self._on_message

    # ------------------------------------------------------------------
    # Public API
    # ------------------------------------------------------------------

    async def connect(self) -> None:
        """Configure TLS, connect, publish introspection, and start MQTT loop."""
        cfg = self._config

        # mTLS: device certificate + key + CA bundle from Astrate Pairing API
        tls_ctx = ssl.create_default_context(ssl.Purpose.SERVER_AUTH, cafile=cfg.ca)
        tls_ctx.load_cert_chain(certfile=cfg.cert, keyfile=cfg.key)
        self._client.tls_set_context(tls_ctx)

        host, port = _parse_broker_url(cfg.astrate_url)
        log.info("Connecting to Astrate broker %s:%d …", host, port)
        self._client.connect(host, port, keepalive=60)
        self._client.loop_start()

        # Wait up to 10 s for the on_connect callback
        if not self._connected.wait(timeout=10):
            raise ConnectionError(f"Timed out connecting to {host}:{port}")

    def publish_game_state(self, state: GameState) -> None:
        """Publish an object-aggregated GameState snapshot."""
        payload = {
            "mapId":      state.map_id,
            "mapName":    state.map_name,
            "playerX":    state.player_x,
            "playerY":    state.player_y,
            "inBattle":   state.in_battle,
            "battleType": state.battle_type,
            "dialogText": state.dialog_text,
            "stasis":     state.stasis,
        }
        topic = f"{self._prefix}/{IFACE_GAME_STATE}/state"
        self._client.publish(topic, _wrap(payload), qos=0)
        log.debug("→ GameState  map=%s pos=(%d,%d) battle=%s stasis=%s",
                  state.map_name, state.player_x, state.player_y,
                  state.in_battle, state.stasis)

    def publish_party_member(self, member: PartyMember) -> None:
        """Publish individual PartyStatus updates for one party slot."""
        slot = str(member.slot_index)
        base = f"{self._prefix}/{IFACE_PARTY_STATUS}/{slot}"
        for endpoint, value in (
            ("name",      member.name),
            ("speciesId", member.species_id),
            ("currentHp", member.current_hp),
            ("maxHp",     member.max_hp),
            ("level",     member.level),
        ):
            self._client.publish(f"{base}/{endpoint}", _wrap(value), qos=0)
        log.debug("→ PartyStatus slot=%s %s HP=%d/%d",
                  slot, member.name, member.current_hp, member.max_hp)

    def set_command_callback(self, cb: Callable[[dict], None]) -> None:
        """Register a callback for incoming ControlCommand messages.

        The callback receives the unwrapped `v` dict:
        {"button": "UP", "holdFrames": 8, "sequenceId": 42}
        """
        self._command_cb = cb

    async def disconnect(self) -> None:
        self._client.loop_stop()
        self._client.disconnect()
        log.info("Disconnected from Astrate broker.")

    # ------------------------------------------------------------------
    # Paho callbacks (run on MQTT background thread)
    # ------------------------------------------------------------------

    def _on_connect(self, client, userdata, flags, reason_code, properties) -> None:
        if reason_code != 0:
            log.error("MQTT connect failed: reason_code=%s", reason_code)
            return

        log.info("Connected to Astrate broker (rc=%s)", reason_code)

        # Astarte MQTT v1: publish introspection to bare device topic
        client.publish(self._prefix, _INTROSPECTION.encode(), qos=2, retain=False)

        # Subscribe to server-owned ControlCommand topic (QoS 2 for dedup guarantee)
        cmd_topic = f"{self._prefix}/{IFACE_CONTROL_COMMAND}/command"
        client.subscribe(cmd_topic, qos=2)
        log.info("Subscribed to %s", cmd_topic)

        self._connected.set()

    def _on_disconnect(self, client, userdata, disconnect_flags, reason_code, properties) -> None:
        log.warning("Disconnected from broker (rc=%s) — paho will reconnect.", reason_code)
        self._connected.clear()

    def _on_message(self, client, userdata, msg) -> None:
        if self._command_cb is None:
            return
        try:
            envelope = json.loads(msg.payload)
            command  = envelope.get("v", {})
            log.debug("← ControlCommand %s", command)
            self._command_cb(command)
        except (json.JSONDecodeError, AttributeError) as exc:
            log.warning("Malformed ControlCommand payload: %s (%s)", msg.payload, exc)


# ------------------------------------------------------------------
# Helper
# ------------------------------------------------------------------

def _parse_broker_url(url: str) -> tuple[str, int]:
    """Extract (host, port) from an Astrate URL.

    Astrate's embedded MQTT broker listens on a separate port from the HTTP API.
    Default: 8883 (mTLS).  Override via ASTRATE_MQTT_PORT env var or by passing
    host:port directly.

    For development without mTLS, set port 1883 and skip TLS config.
    """
    # Strip scheme
    stripped = url.replace("https://", "").replace("http://", "")
    if ":" in stripped:
        host, port_str = stripped.rsplit(":", 1)
        # If the port looks like an HTTP port, assume MQTT is on 8883
        port = int(port_str)
        if port in (80, 8080, 443, 8443):
            port = 8883
    else:
        host = stripped
        port = 8883
    return host, port
