# Mule log

One row per task attempt, written by `tools/mule.sh`. This is the record of what the cheap
layer can actually do — read it before deciding whether a kind of task is worth delegating.

`secs` is the honest signal: a task that used most of its 900s budget was too big.

| date | task | outcome | secs | note |
| --- | --- | --- | --- | --- |
| 2026-09-04 | purge-properties-compression-capabilityauto | done | 86s | 6f7a3d6 |
| 2026-09-04 | broker-acl-coldstart-introspection-miss | done | 629s | d201db4 |
| 2026-09-04 | issue-93 | done | 108s | 24ad5b8 |
| 2026-09-04 | broker-disconnect-device-zombie-session | blocked | 117s | wrote nothing |
| 2026-09-04 | broker-offline-acl-tests | blocked | 208s | wrote nothing |
| 2026-09-04 | broker-onconnect-doc-comment | blocked | 438s | tests failed: --- FAIL: TestMQTTSink_Retained (0.02s) |
| 2026-09-05 | empty-introspection-verification | blocked | 106s | wrote nothing |
| 2026-09-05 | probe-trigger-install-notification-delay | blocked | 308s | wrote nothing |
