# Mule log

One row per task attempt, written by `tools/mule.sh`. This is the record of what the cheap
layer can actually do — read it before deciding whether a kind of task is worth delegating.

`secs` is the honest signal: a task that used most of its 900s budget was too big.

| date | task | outcome | secs | note |
| --- | --- | --- | --- | --- |
| 2026-09-04 | purge-properties-compression-capabilityauto | done | 86s | 6f7a3d6 |
| 2026-09-04 | broker-acl-coldstart-introspection-miss | done | 629s | d201db4 |
| 2026-09-04 | issue-93 | done | 108s | 24ad5b8 |
