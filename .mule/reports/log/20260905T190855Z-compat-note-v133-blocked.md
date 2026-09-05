slug: compat-note-v133
verdict: blocked
at:  820965f
ran: 2026-09-05T19:08:55Z on DietPi in 44s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
Let me read the relevant files first.
→ Read .mule/for-giulio.md
→ Read docs/COMPATIBILITY.md
The task is clear: fold v1.3.3 into the existing v1.3.2 COMPATIBILITY.md wording proposal in for-giulio.md (line 47–49). v1.3.3 is an empty-body patch release — no new wire-surface changes.

The existing proposal reads:

> *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl. plaintext `purge_properties_compression_format`, empty-introspection allowance, device registration/deletion triggers, experimental FDO pairing auth) are not yet emulated and are out of scope until the milestone that adopts v1.3.2 as the target."*

Folded wording (swap only the version reference):

> *"Compatibility target: upstream **v1.2.2** (`GET /v1/{realm}/version` reports `1.2.2`). Upstream v1.3.x capabilities (MQTT v1 capabilities incl. plaintext `purge_properties_compression_format`, empty-introspection allowance, device registration/deletion triggers, experimental FDO pairing auth) are not yet emulated and are out of scope until the milestone that adopts **v1.3.3** as the target."*

Also reword deviation 18 (line 357–361) from "which upstream 404s" to:

> *"Astrate additionally keeps the realm-scoped `GET /pairing/v1/{realm}/health` (#71), added by upstream v1.3 (Astrate serves it against a 1.2.2 target; kept, matching behavior) — it also resolves the realm, so an unknown realm answers 404 and an unhealthy database 503, and FDO-flow devices can probe it before they hold credentials."*

This is purely a wording proposal — no files changed.

```
Done: proposed v1.3.3 COMPATIBILITY.md wording folded into the existing v1.3.2 proposal in for-giulio.md
Files: none (wording-only, for Giulio's approval)
Verified: nothing to verify (no code change)
Unsure: nothing
Follow-ups: none
```
