# Mule queue

One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).
The mule runs the topmost open line. Reorder freely; the order is the priority.

If a line needs more detail than fits, write `.mule/tasks/<slug>.md` and keep the line short.

    tools/mule.sh add "<slug>: <outcome>"
    tools/mule.sh loop

A line tagged `[legion]` needs the Legion Go and is skipped automatically while it is asleep;
the queue moves on to the next runnable line rather than stalling behind it.

## Queue

- [ ] bench-tiers: create the tiered benchmark definitions per .mule/recipes/benchmarks.md (first run only — this task builds the harness, it does not run it)
- [ ] race-check: run `go test -race ./...` on the Legion Go (`ssh legion`) and report failures to .mule/for-giulio.md — the Pi cannot run it, so this is the only race coverage there is. If Go is not installed there, write that one line to .mule/for-giulio.md and stop; do not install it. [legion]
