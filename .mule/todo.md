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
- [ ] race-check: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...` (~40s). Report any failure to .mule/for-giulio.md with the full race report. This is the only race coverage that exists — the Pi cannot run -race. [legion]
