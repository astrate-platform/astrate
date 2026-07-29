"""InputExecutor: queue + main-loop hold timing (no pyboy.tick from MQTT path)."""

from emulator_agent.input_executor import InputExecutor


class _FakePyboy:
    def __init__(self) -> None:
        self.events: list = []
        self.ticks = 0

    def send_input(self, event) -> None:
        self.events.append(event)

    def tick(self) -> None:
        self.ticks += 1


def _pump(ex: InputExecutor, n: int = 1) -> None:
    for _ in range(n):
        ex.before_tick()
        ex.after_tick()


def _we(name: str):
    from pyboy.utils import WindowEvent
    return getattr(WindowEvent, name)


def test_enqueue_dedup_by_sequence_id():
    py = _FakePyboy()
    ex = InputExecutor(py)
    ex.enqueue({"button": "A", "holdFrames": 1, "sequenceId": 1})
    ex.enqueue({"button": "A", "holdFrames": 1, "sequenceId": 1})  # redelivery
    ex.enqueue({"button": "B", "holdFrames": 1, "sequenceId": 2})

    _pump(ex, 1)  # start A, hold counts down to 0 → release
    assert ex.last_sequence_id == 1
    _pump(ex, 1)  # start B
    assert ex.last_sequence_id == 2
    # press A, release A, press B, release B — duplicate seq=1 dropped
    assert py.events == [
        _we("PRESS_BUTTON_A"),
        _we("RELEASE_BUTTON_A"),
        _we("PRESS_BUTTON_B"),
        _we("RELEASE_BUTTON_B"),
    ]


def test_hold_frames_release_after_n_ticks():
    py = _FakePyboy()
    ex = InputExecutor(py)
    ex.enqueue({"button": "START", "holdFrames": 3, "sequenceId": 1})

    ex.before_tick()  # press
    assert py.events == [_we("PRESS_BUTTON_START")]
    assert ex.busy
    ex.after_tick()   # remaining 2

    ex.before_tick()
    ex.after_tick()   # remaining 1
    assert ex.busy
    assert len(py.events) == 1

    ex.before_tick()
    ex.after_tick()   # remaining 0 → release
    assert not ex.busy
    assert py.events[-1] == _we("RELEASE_BUTTON_START")


def test_local_intro_press_skips_sequence_id():
    py = _FakePyboy()
    ex = InputExecutor(py)
    assert ex.enqueue_local("A", hold_frames=1)
    _pump(ex, 1)
    assert ex.last_sequence_id == -1  # local does not advance MQTT seq space
    # MQTT command with sequenceId=1 still accepted
    ex.enqueue({"button": "B", "holdFrames": 1, "sequenceId": 1})
    _pump(ex, 1)
    assert ex.last_sequence_id == 1


def test_enqueue_local_false_when_busy():
    py = _FakePyboy()
    ex = InputExecutor(py)
    ex.enqueue({"button": "A", "holdFrames": 5, "sequenceId": 1})
    ex.before_tick()  # start hold
    assert ex.busy
    assert ex.enqueue_local("START") is False


def test_stub_mode_none_pyboy_consumes_queue():
    ex = InputExecutor(None)
    ex.enqueue({"button": "UP", "holdFrames": 2, "sequenceId": 1})
    _pump(ex, 1)
    assert ex.last_sequence_id == 1
    assert not ex.busy


def test_looks_past_cold_boot():
    from emulator_agent.main import looks_past_cold_boot
    from emulator_agent.state_decoder import GameState

    cold = GameState(
        map_id=0, map_name="Pallet Town", player_x=0, player_y=0,
        in_battle=False, battle_type=0, dialog_text="", stasis=False,
        party=[], timestamp=0.0,
    )
    assert looks_past_cold_boot(cold) is False

    moved = GameState(
        map_id=0, map_name="Pallet Town", player_x=5, player_y=6,
        in_battle=False, battle_type=0, dialog_text="", stasis=False,
        party=[], timestamp=0.0,
    )
    assert looks_past_cold_boot(moved) is True
