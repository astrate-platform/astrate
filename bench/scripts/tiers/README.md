# Benchmark tiers

Named, repeatable deployment sizes so two runs are comparable. Each tier is a
plain `.env` file sourced by `run-tier.sh` — no logic, just `KEY=value`.

| tier | devices | msg/s per device | aggregate msg/s | ingest duration | connstorm | where |
| --- | ---: | ---: | ---: | --- | ---: | --- |
| small | 100 | 0.2 | 20 | 5m | 100 | laptop |
| medium | 1000 | 0.2 | 200 | 10m | 1000 | laptop |
| big | 10000 | 0.1 | 1000 | 15m | 5000 | Legion Go |
| giant | 50000 | 0.1 | 5000 | 20m | 20000 | Legion Go |

## Running a tier

From `bench/`:

```sh
# Bring up the deployment first (one of):
./scripts/up-astrate.sh          # local Astrate
./scripts/up-astarte.sh          # upstream Astarte

# Then run the tier:
./scripts/run-tier.sh small astrate \
    -base-url http://127.0.0.1:8080 \
    -housekeeping-key keys/housekeeping.pem

# Or the Legion Go for large tiers:
./scripts/run-tier.sh big astarte \
    -base-url http://api.astarte.localhost \
    -housekeeping-key .astarte/compose/astarte-keys/housekeeping_private.pem
```

Results land in `bench/results/<tier>-<target>-<UTC timestamp>/` — one JSON per
subcommand, plus the tier env file, `host.txt` with the machine identity, and
`stats.csv` if `sample-stats.sh` ran alongside.

## Resumability

If the state file for this tier already has enough devices, `provision` is
skipped — no need to re-register 50k devices on a re-run.

## Never overwrite

A benchmark result is evidence. If the results directory already exists, the
runner refuses and exits with an error.
