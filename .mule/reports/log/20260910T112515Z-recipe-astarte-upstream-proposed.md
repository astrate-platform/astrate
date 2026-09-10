slug: recipe-astarte-upstream
verdict: proposed
at:  9a28cd1
ran: 2026-09-10T11:25:15Z on DietPi in 37s

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

$ gh api repos/astarte-platform/astarte/tags --jq '.[0:3][] | "\(.name)"'
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3

The newest tag is still `v1.4.0-rc.5` (2026-08-20), unchanged since the last recipe run
(2026-09-09, report .mule/reports/log/20260909T184315Z-recipe-astarte-upstream-proposed.md).
That run already analyzed every change between rc.4 and rc.5 (three data_updater_plant
fixes: the `interface_loading_failed` error-name unmapping #2119, the binaryblob
validation fix, and the object-mapping interface_id filter #2141) and found no
wire-visible gap: the first is absorbed by Astrate's `UpstreamErrorName` fallback, the
second by the BSON decoder's TypeBinary-only acceptance, and the third cannot exist
because Astrate's object decode already receives only the target interface's
`ObjectLeaves`. rc.4/rc.3 are config/Vault/FDO/GC/AMQP changes with no wire surface
Astrate implements, and v1.3.3 has an empty release body.

Nothing new to check, nothing to propose. The standing compat-note lines
(`compat-note-v1.3.3`, `compat-note-v1.4-rc`) are already in the queue.

No upstream movement since v1.4.0-rc.5.