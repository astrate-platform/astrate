# Recipe — heavy workloads on the Lenovo Legion Go

The Legion Go is reachable over SSH as **`legion`** (CachyOS, 16 cores, 12GB RAM, x86_64).
It already runs the full upstream Astarte stack in Docker — Scylla, VerneMQ, RabbitMQ and
the Astarte services — which is exactly what makes it the only machine here that can hold an
Astarte and an Astrate deployment side by side. The Pi you are running on has 4 cores, 3.7GB
and no Docker at all, so anything heavier than `go test` belongs there, not here.

```sh
ssh legion 'docker ps --format "{{.Names}}"'     # what is up
tools/mule.sh legion check                       # ssh + docker, one command
```

## It is optional muscle, never a dependency

A task line tagged `[legion]` is **skipped automatically** when the machine is asleep or off
the network, and the queue moves on to the next runnable task. So:

- Tag every task that needs it: `- [ ] bench-big-astrate: ... [legion]`
- **Never** write a task that silently assumes the Legion Go is up. If it needs it, tag it.
- If you are proposing work and `tools/mule.sh legion check` fails, that is not a blocker —
  propose the tagged tasks anyway. They will run whenever it comes back.

## The Go problem, already solved

**The Legion Go has no Go toolchain, and you must not install one** — that is a change to a
machine nobody asked you to change. The bench harness imports nothing from Astrate and needs
no cgo, so it cross-compiles:

```sh
tools/mule.sh legion bench-push     # builds linux/amd64 here, copies the binary over
ssh legion '~/astrate/bench/bench ingest -state ... '
```

Run `bench-push` at the start of any benchmarking task — it is cheap and it guarantees the
binary matches the current source rather than whatever was left there last time.

## Before a benchmark run — check this, every time

**The UMA/VRAM allocation must be set to 3GB in BIOS/EFI (the lowest option).** The default
carve-out steals several GB of system RAM, which is precisely the resource a device-fleet
benchmark measures. You cannot change this over SSH — it is a firmware setting. What you
*can* do is measure it and refuse to proceed:

```sh
ssh legion 'free -g | head -2'      # total RAM as the OS sees it
```

If total RAM reads well under 12GB, the carve-out is large: **stop, do not run the
benchmark, and write a line to `.mule/for-giulio.md`** saying the Legion Go needs its BIOS
UMA set to 3GB before tier runs are meaningful. A benchmark against a machine short several
GB produces a number that looks fine and means nothing, which is worse than no number.

Also worth checking, for the same reason — a throttled run is invisible in the result:

```sh
ssh legion 'cat /sys/class/power_supply/*/status 2>/dev/null; sensors 2>/dev/null | head'
```

Say in the results whether it was on mains power.

## Running a tier

```sh
ssh legion 'cd ~/astrate && git fetch -q && git checkout -q <sha>'
tools/mule.sh legion bench-push
ssh legion 'cd ~/astrate/bench && ./bench provision ... && ./bench ingest ...'
```

Rules, all load-bearing:

- **Record the host with the numbers**: `uname -a`, `nproc`, `free -h`, `docker info | head`,
  and the git sha, written into the results directory before the run starts.
- **Copy results back to the Pi and commit them.** A results directory that only exists on
  the Legion Go will be gone the next time someone reinstalls that machine.
  `scp -r legion:~/astrate/bench/results/<dir> bench/results/`
- **Never edit a results file to tidy it.** It is evidence.
- Two runs of the same tier before believing a number; report the spread.
- Long runs: prefix with `nohup`/`setsid` and poll, or the SSH session dying takes the
  benchmark with it. Do not hold a 20-minute foreground SSH connection from a timer-driven
  task — the task will hit its own timeout first.

## When it is unreachable

Say so in one line and move on. Do not wait, do not retry in a loop, do not propose work to
"make it reachable" — it is a handheld, it is sometimes off, and that is fine.
