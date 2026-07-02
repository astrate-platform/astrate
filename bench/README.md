# Astarte vs Astrate benchmark harness

A standalone, wire-protocol-only load harness. It deliberately imports
nothing from Astrate: every interaction is the Astarte MQTT v1 protocol plus
the public REST APIs, so **the same binary drives an upstream Astarte
deployment and an Astrate deployment interchangeably** — that wire
compatibility is what makes the comparison apples-to-apples.

```
bench provision   create realm + interfaces + register devices → state file
bench ingest      N devices × R msg/s: e2e latency, PUBACK latency, loss
bench connstorm   mass mTLS connect (cert issuance timed separately)
bench query       AppEngine read mix: latest / 1h window / snapshot
scripts/          stack bring-up + resource/disk sampling
```

## Quick start

Astrate (from `bench/`):

```sh
./scripts/up-astrate.sh
go run . provision -base-url http://127.0.0.1:8080 \
    -housekeeping-key keys/housekeeping.pem -devices 200 -state astrate.json
go run . ingest -state astrate.json -devices 100 -rate 1 -duration 5m
go run . query  -state astrate.json -devices 100
```

Astarte (same host or a twin):

```sh
./scripts/up-astarte.sh        # clones astarte @ v1.2.0, uses its own compose
go run . provision -base-url http://api.astarte.localhost \
    -housekeeping-key .astarte/compose/astarte-keys/housekeeping_private.pem \
    -devices 200 -state astarte.json
go run . ingest -state astarte.json -devices 100 -rate 1 -duration 5m
```

Sample resources in a second terminal during any run:

```sh
./scripts/sample-stats.sh results/astrate-stats.csv 5
```

## What each number means — and the fairness traps

**e2e visibility latency** (the headline latency): probe devices publish a
marker and poll the AppEngine API until it appears; the histogram is
publish→queryable. This is the only latency defined identically on both
platforms. It is quantized by `-probe-poll` (default 25 ms) — report it with
that caveat, or lower the interval and note the extra read load you injected.

**PUBACK latency is NOT comparable across platforms — never headline it.**
Astrate withholds PUBACK until the row is committed (ack-after-commit);
Astarte acks at VerneMQ and persists asynchronously (RabbitMQ → data updater
plant → Cassandra). The same number means "durably stored" on one platform
and "accepted by the broker" on the other. It is reported because it *is*
each platform's device-visible backpressure behavior.

**Loss** counts rows through the AppEngine API against messages sent. On
Astarte, rows lag PUBACKs by design — if loss is nonzero, first increase
`-grace` before concluding anything.

**Closed loop per device:** each simulated device keeps one message in
flight (publish, wait PUBACK, next tick). If acks are slower than the tick,
the device falls behind — `behind_ticks` > 0 means the platform, not the
generator, set the achieved rate. That is a result, not an error.

**Throughput ceiling:** there is no auto-ramp; run `ingest` at increasing
`-devices`/`-rate` and find the knee where e2e p99 degrades or `behind_ticks`
explodes. Change one variable per run.

Other traps the harness cannot absorb for you:

- **Durability knobs decide the winner silently.** PG `synchronous_commit`
  vs Cassandra commitlog sync — run defaults, but document both alongside
  results.
- **Warm-up:** Cassandra behaves differently once compaction starts. Run
  ≥ 15 min steady state; discard the first minutes; report the steady tail.
- **Host sizing is the headline experiment.** Run the matrix on a
  comfortable host (8–16 GB) for a fair performance comparison, *and* on the
  1–2 GB VPS target where "Astarte does not fit" is itself the result.
- **Generator placement:** run `bench` on a separate machine (or CPU-pin it)
  so it does not steal cycles from the system under test. For a remote
  generator against Astarte, map `*.astarte.localhost` to the Docker host in
  `/etc/hosts` and pass `-broker-url` if needed.
- **Reference lines:** Astrate's own budgets (RSS ≤ 150 MB app, ≤ 768 MB PG,
  ingest p99 < 250 ms — DESIGN §5.4) make useful chart annotations.

## Storage efficiency (bytes per datapoint)

After a large ingest, force both platforms to their at-rest representation,
then measure:

```sh
# Astrate: compress the chunks, then measure
docker exec -it astrate-bench-timescaledb-1 psql -U astrate -c \
  "SELECT compress_chunk(c, true) FROM show_chunks('individual_datastreams') c;"
./scripts/disk-usage.sh astrate

# Astarte: flush + major compaction, then measure
docker exec -it <cassandra-container> nodetool flush
docker exec -it <cassandra-container> nodetool compact
./scripts/disk-usage.sh astarte
```

Divide by the row count (`loss_check_found` totals, or query the stores).
Measure a *delta* between two ingests of known size if the baseline schema
overhead should be excluded.

## Mechanics worth knowing

- Both dev stacks use self-signed broker certs; devices default to skipping
  TLS verification (`-insecure-tls`, recorded in the state file). mTLS client
  auth is always on — that part is the real protocol on both platforms.
- Device keys are ECDSA P-256 by default so thousands generate quickly;
  `-rsa` switches to RSA-2048 if an Astarte deployment rejects EC CSRs.
- Provision retries HTTP 429 with backoff (Astrate rate-limits pairing);
  registration parallelism is `-concurrency`.
- The loss check pages with `since_after` above ~9000 rows per device. Row
  timestamps are ms precision on the wire — above ~100 msg/s/device,
  same-millisecond collisions can undercount the paged variant; keep
  per-device rates below that or extend the run instead.
- Interfaces publish with `reliability: guaranteed` (QoS 1) and explicit
  timestamps, so both platforms index identical time series.
- The state file contains the realm private key and device secrets:
  benchmark material, not secrets — but don't commit it (gitignored).
