slug: recipe-astarte-upstream
verdict: proposed
at:  34ba9ca
ran: 2026-09-21T19:06:35Z on DietPi in 130s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read .mule/todo.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z

New upstream release: **v1.3.4** (2026-09-18). Previous newest was v1.3.3, already in queue. Upstream is ahead — reading the release notes.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.4") | .body'
### Fixed

- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
  uncollected binaries over time by forcing more frequent full sweep garbage collections
  on them.
- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.

Both fixes are Elixir/Erlang-runtime specific (garbage collection, RPC server availability) — nothing Astrate implements. Let me check the existing compat-note proposals before appending.
→ Read .mule/for-giulio.md
→ Read .mule/todo.md [offset=174]
No git changes, nothing to append. Summary:

- Newest upstream stable: **v1.3.4** (2026-09-18); newest overall still v1.4.0-rc.5. Nothing shipped since the 2026-09-18 astarte-upstream run already covered this tag.
- v1.3.4 release body: two `astarte_data_updater_plant` fixes — fuller fullsweep GC on AMQPDataConsumer processes, and RPC-server availability. Both Elixir/Erlang-runtime concerns; neither touches a surface Astrate implements (no Erlang VM, no AMQP, single-binary no node-RPC). Per the recipe, Elixir-shaped changes are ignored.
- The version-reference update is already delivered: `compat-note-v1.3.4` is ticked `- [x]` in `.mule/todo.md:177`, and the wording proposal for v1.3.2/v1.3.3/v1.3.4 already stands in `.mule/for-giulio.md` (line 49). Re-proposing would duplicate.

Appending nothing to `.mule/todo.md` — an empty result per the recipe's "no upstream movement" rule.

---

Done: no upstream movement since v1.3.4 — still newest stable (2026-09-18, DUP GC/RPC fixes, Elixir-shaped, nothing proposable); v1.4.0 still rc.5-only
Files: none
Verified: `gh api repos/astarte-platform/astarte/releases` -> newest stable v1.3.4, newest overall v1.4.0-rc.5 (pass); read v1.3.4 release body (pass)
Unsure: nothing
Follow-ups: none
