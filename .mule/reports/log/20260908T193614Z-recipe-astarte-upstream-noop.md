slug: recipe-astarte-upstream
verdict: noop
at:  eca0e08
ran: 2026-09-08T19:36:14Z on linux x86_64


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md
→ Read docs/UPSTREAM-EXPERIMENTAL.md
→ Read .mule/todo.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

The astarte-upstream recipe already ran earlier today (2026-09-08T11:34:22Z, at 1e56fda,
report .mule/reports/log/20260908T113422Z-recipe-astarte-upstream-proposed.md). The release
list is unchanged since the 2026-08-20 v1.4.0-rc.5 tag: no new stable, no new rc, no movement
on any other repo in this recipe's scope.

Everything upstream has shipped since Astrate's 1.2.2 target is already accounted for in the
queue from prior runs, including earlier today:

- v1.3.x: purge-properties-compression-capability [x], empty-introspection-verification [!],
  probe-trigger-install-notification-delay [!], probe-property-resend-encoding [!],
  probe-props-resend-error-triggers [!], compat-note-v1.3.2 [x], compat-note-v1.3.3 [!]
- v1.4.0-rc line: probe-required-mapping-flag [x], compat-note-v1.4-rc [!]

No new proposal is qualified: nothing to add, and re-proposing the compat-note or a probe
line already in the queue would duplicate. Recipe rule honoured: "If nothing changed since
last time, append nothing to the queue and say 'no upstream movement since <tag>'."

No upstream movement since v1.4.0-rc.5.

Done: upstream watch — no upstream movement since v1.4.0-rc.5; appended nothing to .mule/todo.md
Files: .mule/reports/log/20260908T193614Z-recipe-astarte-upstream-noop.md (evidence; todo.md untouched)
Verified: gh api repos/astarte-platform/astarte/releases -> v1.4.0-rc.5 newest, unchanged since 2026-08-20; git status clean
Unsure: nothing
Follow-ups: none