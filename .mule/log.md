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
| 2026-07-27 | race-check | checked | 54s | 8cc751c |
| 2026-07-27 | race-check | checked | 51s | 886ff24 |
| 2026-07-27 | race-check | checked | 42s | 57338db |
| 2026-07-27 | race-check | checked | 56s | ccb4450 |
| 2026-07-27 | race-check | checked | 74s | a0ae9f4 |
| 2026-07-27 | race-check | checked | 82s | 7099f80 |
| 2026-07-27 | issue-6 | transient | 19s |  > build · big-pickle  Error: No provider available  |
| 2026-07-27 | issue-6 | transient | 20s |  > build · big-pickle  Error: No provider available  |
| 2026-07-27 | race-check | transient | 21s |  > build · big-pickle  Error: No provider available  |
| 2026-07-27 | race-check | transient | 22s |  > build · big-pickle  Error: No provider available  |
| 2026-07-27 | race-check | checked | 59s | 13cd7f1 |
| 2026-07-27 | race-check | checked | 46s | 452b303 |
| 2026-07-27 | race-check | checked | 65s | 5ac1fdb |
| 2026-07-27 | issue-6 | done | 401s | 997cdef |
| 2026-07-27 | race-check | checked | 88s | 70bc4be |
| 2026-07-27 | store-realm-cascade-policies | done | 122s | 0480501 |
| 2026-07-27 | race-check | checked | 108s | 16cc008 |
| 2026-07-27 | store-alias-lowest-id | done | 194s | 4af622f |
| 2026-07-27 | race-check | checked | 48s | 344d013 |
| 2026-07-27 | store-delete-device-objects | done | 321s | deb01ac |
| 2026-07-27 | race-check | checked | 71s | ad335ec |
| 2026-07-27 | issue-16 | done | 1092s | d670932 |
| 2026-07-27 | race-check | checked | 72s | cc13242 |
| 2026-07-27 | issue-20 | checked | 211s | 446b806 |
| 2026-07-27 | issue-22 | blocked | 340s | touched a never-touch path |
| 2026-07-27 | issue-21 | done | 297s | 9d85191 |
| 2026-07-27 | issue-15 | done | 209s | 5ea4c1b |
| 2026-07-28 | race-check | blocked | 1200s | TIMEOUT after 1200s — task too big, split it |
| 2026-07-28 | issue-27 | blocked | 1157s | TIMEOUT after 1157s — task too big, split it |
| 2026-07-28 | issue-26 | done | 731s | e250381 |
| 2026-07-28 | issue-25 | done | 261s | ee59732 |
| 2026-07-28 | issue-24 | blocked | 433s | touched a never-touch path |
| 2026-07-28 | issue-23 | done | 567s | a9481d0 |
| 2026-07-28 | issue-20 | checked | 152s | c59a7eb |
| 2026-07-28 | issue-14 | blocked | 156s | gates failed |
| 2026-07-28 | issue-13 | checked | 136s | 4af2073 |
| 2026-07-28 | issue-12 | done | 607s | 89b1004 |
| 2026-07-28 | issue-20 | checked | 103s | 33d5be6 |
| 2026-07-28 | issue-13 | checked | 128s | beecab8 |
| 2026-08-31 | issue-91 | blocked | 845s | gates failed |
| 2026-08-31 | issue-68 | blocked | 533s | gates failed |
| 2026-08-31 | control-producer-properties-compression | blocked | 332s | gates failed |
| 2026-08-31 | probe-interface-default-values | blocked | 361s | gates failed |
| 2026-08-31 | probe-value-type-validation | blocked | 48s | gates failed |
| 2026-08-31 | compat-note-v140-rc3 | blocked | 222s | gates failed |
