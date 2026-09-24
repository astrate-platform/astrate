slug: recipe-astarte-upstream
verdict: proposed
at:  7f03766
ran: 2026-09-24T19:44:58Z on DietPi in 105s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z

→ Read .mule/todo.md
→ Read docs/UPSTREAM-EXPERIMENTAL.md
Upstream is ahead: newest tag is **v1.3.4** (2026-09-18). Let me read its release notes.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.4") | .body'
### Fixed

- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
  uncollected binaries over time by forcing more frequent full sweep garbage collections
  on them.
- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.

v1.3.4 is a patch release: two `astarte_data_updater_plant` fixes (Erlang GC sweep forcing on AMQP consumers, RPC server availability). Both are Elixir/OTP-shaped — no wire-visible behaviour, nothing Astrate implements. Nothing to propose from the diff; only the required compat-note line. Let me check the existing wording proposal to fold into.
$ ls .mule/ && ls .mule/for-giulio.md 2>/dev/null; ls -la .mule/for-giulio.md 2>/dev/null
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
recipes
reports
reviews
tasks
todo.md
.mule/for-giulio.md
-rw-r--r-- 1 root root 59011 Sep 24 21:21 .mule/for-giulio.md

✱ Grep "v1\.3\.[0-9]|COMPATIBILITY" in .mule/for-giulio.md · 100 matches
The v1.3.4 check is not new: the 2026-09-18 run already proposed it (`compat-note-v1.3.4`, done in todo.md:177) and the wording was folded into `for-giulio.md:69`. The release sweep confirms the newest tag is still **v1.3.4** (2026-09-18) — unchanged since that run. No v1.4.0 stable (still rc.5 only).

Per the recipe's "if nothing changed since last time" rule: append nothing, no evidence file.

Done: no upstream movement since v1.3.4 — newest tag still v1.3.4 (2026-09-18, maintenance-only DUP GC/RPC fixes, no wire/API surface change); v1.4.0 still rc.5-only; the 1.3.x compat-note wording already sits in .mule/for-giulio.md and the existing compat-note-v1.3.4 line is done, so appending it again would duplicate it
Files: none
Verified: `gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'` -> pass (newest v1.3.4, 2026-09-18); cross-checked todo.md:177 and for-giulio.md against prior runs -> no gap
Unsure: nothing
Follow-ups: none (empty result is the good result)
Git: untouched
