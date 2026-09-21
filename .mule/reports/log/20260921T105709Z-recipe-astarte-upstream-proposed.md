slug: recipe-astarte-upstream
verdict: proposed
at:  7aa7f65
ran: 2026-09-21T10:57:09Z on DietPi in 110s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md
→ Read .mule/todo.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z

→ Read .mule/for-giulio.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.4") | .body'
### Fixed

- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
  uncollected binaries over time by forcing more frequent full sweep garbage collections
  on them.
- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body'

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.2") | .body'
Astarte v1.3.2 release.

## Changes since v1.3.1

## Fixed

- Avoid crash on invalid properties message

$ ls -la .mule/; echo "--- tasks ---"; ls .mule/tasks/ 2>/dev/null
total 5792
drwxr-xr-x  7 root root    4096 Sep 21 12:55 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root       0 Sep 21 12:55 .alarmed
-rw-r--r--  1 root root      13 Sep 21 12:55 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 20 21:27 .heartbeat
-rw-r--r--  1 root root       0 Sep 21 12:55 .last-output
-rw-r--r--  1 root root       0 Sep 21 12:55 .lock
-rw-r--r--  1 root root      17 Sep 21 12:55 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 5705925 Sep 21 12:55 cron.log
drwxr-xr-x  2 root root    4096 Sep 19 12:52 failed
-rw-r--r--  1 root root   50415 Sep 20 22:42 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    8143 Sep 20 21:27 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 20 19:30 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   70177 Sep 20 21:27 todo.md
--- tasks ---
docs-sync-appengine-data-output-params.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md

$ git log --oneline -15 -- .mule/todo.md; echo "---"; git log --oneline -8 -- .mule/for-giulio.md; echo "--- reports ---"; ls -la .mule/reports/ 2>/dev/null | head -40
395d8bb mule: log docs-sync-native-socket-missing-403-500
ba8afaf mule: log docs-sync-native-metrics-example-fake-series
d2a8c51 mule: refill the queue
ea4f1b9 mule: log compat-note-custom-action-policy-boundary
d68bc52 mule: log triggers-custom-action-policy-nodecide
3c8b89a mule: triggers-custom-action-policy-nodecide [auto]: the custom-action (forward) delivery path never consults an attached policy — `forward()` (internal/engine/triggers/actions.go:469-485) does a single attempt with no `policy.Decide`, no retry, no discard, while `webhook()` (actions.go:505-579) honours the full contract; yet `Enqueue`'s maximum_capacity (actions.go:335-346) and `deliver`'s event_ttl (actions.go:406-413) both apply to custom actions, so a custom-action trigger with a retry policy gets capacity/TTL but silently not retry/discard. Decide intent (route forward through the Decide loop, or pin single-shot as designed) and add a test in actions_test.go that forces the choice.
7e2ceeb mule: log errorname-missing-required-test
e645b75 mule: refill the queue
054d305 mule: blocked channels-rejoin-authz-mismatch
ebc3854 mule: log channels-rejoin-joinref-tagging-test
ab18c72 mule: log channels-group-watch-membership-scope
7bf6f8e mule: refill the queue
f6dc460 mule: refill the queue
216e5a3 mule: log compat-note-v134
194ac6a mule: refill the queue
---
7c9809b mule: refill the queue
2e34867 mule: refill the queue
8cd17a8 mule: compat-note-custom-action-policy-boundary [auto]: propose a docs/COMPATIBILITY.md §8 sentence stating that trigger-policy retry/discard handlers apply only to HTTP webhook actions — custom (forwarder) actions are single-shot by design, receiving only maximum_capacity/event_ttl — via .mule/for-giulio.md (do not edit the file).
1dba3d5 mule: refill the queue
74021c8 mule: refill the queue
79439e6 mule: refill the queue
f6dc460 mule: refill the queue
bcd4383 mule: compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.
--- reports ---
total 40
drwxr-xr-x 3 root root  4096 Sep  6 19:40 .
drwxr-xr-x 7 root root  4096 Sep 21 12:55 ..
-rw-r--r-- 1 root root  2279 Sep  6 19:40 hygiene.md
drwxr-xr-x 2 root root 24576 Sep 21 04:06 log

$ git show -s --format='%ci %h %s' bcd4383 216e5a3 395d8bb; echo "--- HEAD ---"; git rev-parse --abbrev-ref HEAD; git log --oneline -3
2026-09-18 18:53:51 +0200 bcd4383 mule: compat-note-v1.3.4 [auto]: propose the docs/COMPATIBILITY.md wording for v1.3.4 (newest stable, 2026-09-18, maintenance-only patch — DUP fullsweep-GC tuning + RPC-availability fixes, no wire/API surface change; v1.4.0 is still rc.5-only) in .mule/for-giulio.md (do not edit the file) — the open v1.3.2 wording proposal in for-giulio.md stands; fold v1.3.3 and v1.3.4 in rather than re-deriving it.
2026-09-18 18:53:51 +0200 216e5a3 mule: log compat-note-v134
2026-09-20 21:27:02 +0200 395d8bb mule: log docs-sync-native-socket-missing-403-500
--- HEAD ---
mule/queue
7aa7f65 mule: recipe astarte-upstream ran (99s)
7c9809b mule: refill the queue
b82e0ec mule: recipe github-issues ran (186s)

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:12][] | "\(.tag_name)\t\(.prerelease)\t\(.published_at)"'
v1.3.4	false	2026-09-18T09:42:57Z
v1.4.0-rc.5	true	2026-08-20T12:28:04Z
v1.3.3	false	2026-08-07T14:33:58Z
v1.4.0-rc.4	true	2026-08-06T13:40:20Z
v1.4.0-rc.3	true	2026-07-31T10:48:39Z
v1.4.0-rc.2	true	2026-07-14T10:16:22Z
v1.4.0-rc.1	true	2026-07-14T08:53:27Z
v1.3.2	false	2026-07-14T10:07:26Z
v1.3.1	false	2026-07-13T11:50:05Z
v1.3.0	false	2026-05-06T15:30:50Z
v1.2.2	false	2026-04-27T16:39:28Z
v1.4.0-rc.0	true	2026-04-08T12:32:33Z

$ git show --stat 7aa7f65; echo "==== diff ===="; git show 7aa7f65 | head -80
commit 7aa7f65a0879321c40ee819b961b4c25c1ed8efa
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Sep 21 04:06:24 2026 +0200

    mule: recipe astarte-upstream ran (99s)

 .../log/20260921T020624Z-recipe-astarte-upstream-proposed.md     | 9 +++++++++
 1 file changed, 9 insertions(+)
==== diff ====
commit 7aa7f65a0879321c40ee819b961b4c25c1ed8efa
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Sep 21 04:06:24 2026 +0200

    mule: recipe astarte-upstream ran (99s)

diff --git a/.mule/reports/log/20260921T020624Z-recipe-astarte-upstream-proposed.md b/.mule/reports/log/20260921T020624Z-recipe-astarte-upstream-proposed.md
new file mode 100644
index 0000000..f712b47
--- /dev/null
+++ b/.mule/reports/log/20260921T020624Z-recipe-astarte-upstream-proposed.md
@@ -0,0 +1,9 @@
+slug: recipe-astarte-upstream
+verdict: proposed
+at:  7c9809b
+ran: 2026-09-21T02:06:24Z on DietPi in 99s
+
+
+> build · big-pickle
+
+Error: Cannot connect to API: Unable to connect. Is the computer able to access the url?

Checked upstream releases and the compat state.

Newest upstream tag is `v1.3.4` (2026-09-18), unchanged since the last successful astarte-upstream run (`bcd4383`, 2026-09-18, which already proposed and ran `compat-note-v1.3.4`). The 2026-09-21T02:06Z run of this recipe errored (`Cannot connect to API`) and proposed nothing. v1.3.4 is maintenance-only and Elixir-shaped (DUP fullsweep-GC tuning + RPC-availability), no wire/API surface change; v1.4.0 remains rc.5-only and its deltas are already tracked (#137/#138). The open v1.3.2→v1.3.4 wording proposal already sits in `.mule/for-giulio.md`, so a fresh `compat-note-v1.3.4` line would duplicate it.

No upstream movement since v1.3.4 — appended nothing to the queue, wrote no evidence file, touched no git.

```
Done: no upstream movement since v1.3.4; no task lines proposed
Files: none
Verified: gh api repos/astarte-platform/astarte/releases -> newest stable v1.3.4 (2026-09-18), no release since; prior successful run bcd4383 already covered it -> pass
Unsure: nothing
Follow-ups: none
```
