# Pokémon Agent

Bridging a Game Boy emulator to an Astrate IoT platform and LLM Orchestrator.

## Prerequisites
- Python 3.11+
- Astrate instance (local: TimescaleDB + `astrate` binary with `insecure_dev_mode`)
- OpenAI API Key (orchestrator only)
- Pokémon Red ROM (not needed for `--stub` smoke)

## Quick smoke (T1 + T2)

See [docs/TESTING.md](docs/TESTING.md). Short path against the auto-provisioned
`test` realm (`deploy/devrealm` keys):

```sh
# after Astrate is up on :8080 with insecure_dev_mode
make -C examples/pokemon-agent interfaces-install
DEVICE_ID=$(astartectl utils device-id generate-random)
astartectl pairing agent register "$DEVICE_ID" \
  -u http://localhost:8080 -k deploy/devrealm/realm_private.pem -r test
make -C examples/pokemon-agent run-emulator-stub DEVICE_ID=$DEVICE_ID
```

Architecture and ADRs: [docs/DESIGN.md](docs/DESIGN.md), [docs/DECISIONS.md](docs/DECISIONS.md).
