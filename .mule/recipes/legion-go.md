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

## It is the only machine that can run the race detector

Go 1.26.5 is installed there as a **userland toolchain** at `~/.local/go` (no root; removable
with `rm -rf ~/.local/go`), and is on the PATH for non-interactive ssh via `config.fish`.

This matters more than anything else in this file. **The Pi cannot run `-race` at all** — its
kernel has a 39-bit VMA and ThreadSanitizer needs 48 — so the Legion Go is the *only* place
race coverage exists. It takes about 40 seconds there on 16 cores:

```sh
ssh legion 'cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...'
```

Run it after any batch of merged work, and report failures to `.mule/for-giulio.md`. A data
race that the Pi's gate cannot see is exactly the defect this machine exists to catch.

The bench binary can now be built there directly, but `tools/mule.sh legion bench-push` still
cross-compiles and copies one in — cheap, and it guarantees the binary matches the source you
think it does rather than whatever that checkout happens to be on.

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
