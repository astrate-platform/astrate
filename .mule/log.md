# Mule log

One row per task attempt, written by `tools/mule.sh`. This is the record of what the cheap
layer can actually do — read it before deciding whether a kind of task is worth delegating.

`secs` is the honest signal: a task that used most of its 900s budget was too big.

| date | task | outcome | secs | note |
| --- | --- | --- | --- | --- |
| 2026-07-27 | race-check | done | 42s | c8bec10 |
| 2026-07-27 | race-check | done | 42s | ran before the sed fix; marked by hand |
| 2026-07-27 | race-check | blocked | 18s | opencode exited 1 |
| 2026-07-27 | race-check | transient | 20s |  > build · big-pickle  Error: No provider available  |
| 2026-07-27 | race-check | blocked | 60s | wrote nothing |
