slug: recipe-astarte-upstream
verdict: noop
at:  ed715ec
ran: 2026-09-16T18:41:59Z on linux x86_64


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md
→ Read .mule/todo.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

The release list is unchanged since the last astarte-upstream recipe run (2026-09-08, report
.mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md): no new stable, no new
rc, still `v1.4.0-rc.5` (2026-08-20) as the newest tag, `v1.3.3` (2026-08-07) as the newest
stable. No movement on any repo in this recipe's scope since then.

Astrate targets upstream 1.2.2 (docs/COMPATIBILITY.md, `APICompatVersion` 1.2.2, #90 frozen
decision); everything upstream has shipped since that target is already accounted for in the
queue from prior runs:

- v1.3.x: purge-properties-compression-capability [x], empty-introspection-verification [!],
  probe-trigger-install-notification-delay [!], probe-property-resend-encoding [!],
  probe-props-resend-error-triggers [!], compat-note-v1.3.2 [x], compat-note-v1.3.3 [!]
- v1.4.0-rc line: probe-required-mapping-flag [x], compat-note-v1.4-rc [!]

No new proposal is qualified: nothing to add, and re-proposing the compat-note or a probe
line already in the queue would duplicate. Recipe rule honoured: "If nothing changed since
last time, append nothing to the queue and say 'no upstream movement since <tag>'."

No upstream movement since v1.4.0-rc.5.

Done: upstream watch — no upstream movement since v1.4.0-rc.5; appended nothing to .mule/todo.md
Files: .mule/reports/log/20260916T184159Z-recipe-astarte-upstream-noop.md (evidence; todo.md untouched)
Verified: gh api repos/astarte-platform/astarte/releases -> v1.4.0-rc.5 newest, unchanged since 2026-08-20
Unsure: nothing
Follow-ups: none