slug: recipe-astarte-upstream
verdict: noop
at:  eb44d02
ran: 2026-09-10T20:29:36Z on linux x86_64


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
rc, no movement on any other repo in this recipe's scope. The last recipe run
(2026-09-10T11:25Z, at 9a28cd1, report .mule/reports/log/20260910T112515Z-recipe-astarte-upstream-proposed.md)
already analyzed every change shipped since Astrate's 1.2.2 target — rc.4→rc.5 (three
data_updater_plant fixes: the #2119 error-name unmapping absorbed by Astrate's
`UpstreamErrorName` fallback, the binaryblob validation fix absorbed by the BSON decoder's
TypeBinary-only acceptance, the #2141 object-mapping interface_id filter, which cannot exist
because Astrate's object decode already receives only the target interface's `ObjectLeaves`),
rc.4/rc.3 (config/Vault/FDO/GC/AMQP, no wire surface Astrate implements), and v1.3.3 (empty
release body). The standing compat-note lines (`compat-note-v1.3.3`, `compat-note-v1.4-rc`)
are already in the queue; the 1.4-experimental register and its reconcile-vs-promote call stay
gated on a stable v1.4.0.

No new proposal is qualified: nothing to add, and re-proposing a compat-note or probe line
already in the queue would duplicate. Recipe rule honoured: "If nothing changed since last
time, append nothing to the queue and say 'no upstream movement since <tag>'."

No upstream movement since v1.4.0-rc.5.

Done: upstream watch — no upstream movement since v1.4.0-rc.5; appended nothing to .mule/todo.md
Files: .mule/reports/log/20260910T202936Z-recipe-astarte-upstream-noop.md (evidence; todo.md untouched)
Verified: gh api repos/astarte-platform/astarte/releases -> v1.4.0-rc.5 newest, unchanged since 2026-08-20
Unsure: nothing
Follow-ups: none