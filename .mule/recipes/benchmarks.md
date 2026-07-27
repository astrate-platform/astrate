# Recipe — benchmarks

The harness already exists: `bench/` is a standalone wire-protocol load driver that runs
against **either** Astrate or upstream Astarte (`bench/README.md`, ~40 lines — read it).
Subcommands: `provision`, `ingest`, `connstorm`, `query`. Bring-up scripts are in
`bench/scripts/`.

What is missing is *tiers*: named, repeatable deployment sizes so two runs are comparable.

## First run only — build the tiers

If `bench/scripts/tiers/` does not exist, this is the whole job. Create:

- `bench/scripts/tiers/<tier>.env` for each of the four tiers below — plain `KEY=value`,
  no logic.
- `bench/scripts/run-tier.sh` — sources a tier env file, runs `provision`, `ingest`,
  `connstorm` and `query` in that order against a `-base-url` given on the command line,
  and writes every result under `bench/results/<tier>-<target>-<UTC timestamp>/`
  (one file per subcommand, plus the tier env file copied in verbatim, plus a
  `host.txt` recording `uname -a`, CPU count and total RAM).
- `bench/scripts/tiers/README.md` — the table below, plus how to run one.

The tiers, chosen to bracket real Astarte deployments:

| tier | devices | msg/s per device | aggregate msg/s | ingest | connstorm | where |
| --- | ---: | ---: | ---: | --- | ---: | --- |
| small | 100 | 0.2 | 20 | 5m | 100 | laptop |
| medium | 1000 | 0.2 | 200 | 10m | 1000 | laptop |
| big | 10000 | 0.1 | 1000 | 15m | 5000 | Legion Go |
| giant | 50000 | 0.1 | 5000 | 20m | 20000 | Legion Go |

Rules for the runner, all of them load-bearing:

- **Never overwrite a results directory.** A benchmark result is evidence; if the target
  directory exists, fail.
- **Record the host.** A number without the machine that produced it is worthless.
- **Provision is separate from measurement.** Time certificate issuance on its own — that is
  what `connstorm` exists for — and never fold it into the ingest latency.
- Sample resources during the run with `bench/scripts/sample-stats.sh` and keep its output
  in the results directory.
- The runner must be resumable: if `provision` already produced the state file for this
  tier, reuse it rather than re-registering 50k devices.

Do **not** run a benchmark in this task. Building the harness and running it are separate
tasks with separate failure modes.

## Later runs — propose measurement tasks

Append task lines, one per tier and target, e.g.:

```
- [ ] bench-small-astrate: run bench/scripts/run-tier.sh small against a local Astrate stack, commit the results directory
```

Rules:

- **Two runs of the same tier and target before believing a number.** A single pass has been
  wrong here before. Report the spread, not just the value.
- `big` and `giant` do not belong on the laptop — propose them for the Legion Go and read
  `.mule/recipes/legion-go.md` for how that handoff works.
- Comparing against upstream Astarte means the *same tier on the same host*, and
  `bench/scripts/up-astarte.sh` clones a pinned Astarte version — record which.
- If a run produces a surprising number, the follow-up task is a **probe** that explains it,
  not an optimisation. Do not propose a performance fix for a number nobody has explained.
