"""Unit tests for MQTT URL / insecure-dev connection helpers."""

from emulator_agent.astrate_client import _parse_broker_url


def test_parse_broker_url_defaults_to_mtls_port():
    assert _parse_broker_url("http://localhost:8080") == ("localhost", 8883)
    assert _parse_broker_url("https://astrate.example") == ("astrate.example", 8883)


def test_parse_broker_url_insecure_defaults_to_1883():
    assert _parse_broker_url("http://localhost:8080", insecure=True) == ("localhost", 1883)


def test_parse_broker_url_mqtt_port_override_wins():
    assert _parse_broker_url("http://localhost:8080", mqtt_port=1883) == ("localhost", 1883)
    assert _parse_broker_url(
        "http://localhost:8080", mqtt_port=9999, insecure=True
    ) == ("localhost", 9999)
