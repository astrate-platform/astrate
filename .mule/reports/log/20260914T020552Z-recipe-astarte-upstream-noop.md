slug: recipe-astarte-upstream
verdict: noop
at:  d29e382
ran: 2026-09-14T02:05:52Z on linux x86_64


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

The release list is unchanged since the 2026-08-20 v1.4.0-rc.5 tag: no new stable, no new
rc, no movement in this recipe's scope. The last `astarte-upstream` recipe run
(2026-09-10T20:29Z, at eb44d02, report
.mule/reports/log/20260910T202936Z-recipe-astarte-upstream-noop.md) already analyzed every
change shipped since Astrate's 1.2.2 target — rc.4→rc.5 (three data_updater_plant fixes
absorbed by Astrate's `UpstreamErrorName` fallback, BSON TypeBinary-only acceptance, and
the #2141 object-mapping filter that cannot exist in Astrate's decode path), rc.4/rc.3
(config/Vault/FDO/GC/AMQP, no wire surface Astrate implements), and v1.3.3 (empty release
body). Astrate's own target is unchanged: `APICompatVersion` stays 1.2.2 (#90's frozen
decision), and the 1.4 experience is tracked as experimental.

No new proposal is qualified: nothing to add, and re-proposing the compat-note or probe
lines already in the queue (`compat-note-v1.3.3`, `compat-note-v1.4-rc`,
`probe-required-mapping-flag`, and the v1.3.0 probes) would duplicate. Recipe rule
honoured: "If nothing changed since last time, append nothing to the queue and say 'no
upstream movement since <tag>'."

No upstream movement since v1.4.0-rc.5.

Done: upstream watch — no upstream movement since v1.4.0-rc.5; appended nothing to .mule/todo.md
Files: .mule/reports/log/20260914T020552Z-recipe-astarte-upstream-noop.md (evidence; todo.md untouched)
Verified: gh api repos/astarte-platform/astarte/releases -> v1.4.0-rc.5 newest, unchanged since 2026-08-20
Unsure: nothing
Follow-ups: none