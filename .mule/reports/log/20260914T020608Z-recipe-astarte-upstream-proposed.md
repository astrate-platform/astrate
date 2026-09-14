slug: recipe-astarte-upstream
verdict: proposed
at:  d29e382
ran: 2026-09-14T02:06:08Z on DietPi in 100s


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

$ ls -la /root/astrate-mule/.mule/ && ls /root/astrate-mule/docs/
total 3580
drwxr-xr-x  7 root root    4096 Sep 13 22:23 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      14 Sep 13 22:42 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 13 21:19 .heartbeat
-rw-r--r--  1 root root       0 Sep 14 04:04 .last-output
-rw-r--r--  1 root root       0 Sep 13 22:42 .lock
-rw-r--r--  1 root root      17 Sep 13 22:42 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 3477746 Sep 13 22:44 cron.log
drwxr-xr-x  2 root root    4096 Sep 11 20:29 failed
-rw-r--r--  1 root root   33056 Sep 13 22:31 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    6002 Sep 13 21:19 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 13 19:26 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   50378 Sep 13 21:19 todo.md
COMPATIBILITY.md
DESIGN.md
JSON-PAYLOAD-PROFILE.md
Makefile
OPERATIONS.md
ROADMAP.md
UPSTREAM-EXPERIMENTAL.md
api
embed.go
handoff
mkdocs.yml
requirements.txt
site
swagger-ui

→ Read .mule/todo.md
The newest release is still `v1.4.0-rc.5` (2026-08-20). Let me check what the previous run of this recipe concluded, so I know what "last time" saw.
✱ Grep "upstream movement|astarte-upstream|v1\.4\.0-rc|UPSTREAM-EXPERIMENTAL" in .mule · 100 matches
Upstream has not moved. The newest tag is still `v1.4.0-rc.5` (2026-08-20), exactly as the last several recipe runs recorded. Let me look at how a prior no-op run wrote its evidence so I follow the same convention.
→ Read .mule/reports/log/20260910T202936Z-recipe-astarte-upstream-noop.md

$ git rev-parse --short HEAD && date -u +%Y-%m-%dT%H:%M:%SZ
d29e382
2026-09-14T02:05:52Z

← Write .mule/reports/log/20260914T020552Z-recipe-astarte-upstream-noop.md
Wrote file successfully.

Upstream is unchanged since the last run of this recipe. Writing the no-op evidence file in the established convention, and appending nothing to the queue.
Done: upstream watch — no upstream movement since v1.4.0-rc.5 (2026-08-20); newest tag unchanged, appended nothing to `.mule/todo.md`.

Per the recipe's "empty result is a good result" rule: the last `astarte-upstream` run already analyzed every change since Astrate's 1.2.2 target, and the standing compat-note/probe lines (`compat-note-v1.4-rc`, `probe-required-mapping-flag`, etc.) are already queued, so re-proposing would duplicate. Evidence written to `.mule/reports/log/20260914T020552Z-recipe-astarte-upstream-noop.md`.

```
Done: upstream watch — no upstream movement since v1.4.0-rc.5; appended nothing to .mule/todo.md
Files: .mule/reports/log/20260914T020552Z-recipe-astarte-upstream-noop.md
Verified: gh api repos/astarte-platform/astarte/releases -> v1.4.0-rc.5 newest, unchanged since 2026-08-20 -> pass
Unsure: nothing
Follow-ups: none
```
