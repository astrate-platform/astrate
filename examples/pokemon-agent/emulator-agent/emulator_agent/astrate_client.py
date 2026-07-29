import paho.mqtt.client as mqtt
import json

class AstrateClient:
    def __init__(self, config):
        self.config = config
        self.client = mqtt.Client(client_id=config.astrate_device_id, protocol=mqtt.MQTTv311)
        self.client.on_message = self._on_message
        self._command_cb = None

    def _on_message(self, client, userdata, msg):
        if self._command_cb:
            try:
                payload = json.loads(msg.payload)
                self._command_cb(payload.get("v", {}))
            except Exception:
                pass

    async def connect(self) -> None:
        pass

    def publish_game_state(self, state) -> None:
        pass

    def publish_party_member(self, member) -> None:
        pass

    def set_command_callback(self, cb) -> None:
        self._command_cb = cb

    async def disconnect(self) -> None:
        pass
