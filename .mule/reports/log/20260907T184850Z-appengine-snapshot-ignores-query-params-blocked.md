slug: appengine-snapshot-ignores-query-params
verdict: blocked
at:  7cbb4b1
ran: 2026-09-07T18:48:50Z on DietPi in 327s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

# Todos
[•] Read data.go around lines 118-168 to understand the snapshot branch and param checks
[ ] Probe upstream Astarte to decide honour-vs-refuse
[ ] Implement the chosen behaviour
[ ] Write a T1 test pinning the behaviour
[ ] Run gate: go vet + go test + gofmt

→ Read internal/appengine/data.go [offset=90, limit=110]
→ Read docs/api/astarte_appengine_api.yaml [limit=50, offset=555]
✱ Grep "type QueryOpts struct" in internal/appengine · 1 match
✱ Glob "internal/appengine/*_test.go" 11 matches
→ Read internal/appengine/data.go [limit=90, offset=1]
→ Read internal/appengine/query_opts_test.go
✗ Read internal/appengine/data_test.go failed
Error: File not found: /root/astrate-mule/internal/appengine/data_test.go
✱ Grep "func Test.*GetData|GetData\(" in internal/appengine · 2 matches
✱ Grep "IndividualSnapshot" in internal/appengine · 1 match
✱ Grep "func .*IndividualSnapshot" in internal · 1 match
→ Read internal/appengine/http.go [limit=50, offset=250]
→ Read internal/store/datastreams.go [limit=70, offset=240]
✱ Glob ".mule/recipes/*.md" 0 matches

$ ssh -o ConnectTimeout=5 legion 'echo ok'
ssh: connect to host legion.tail4db093.ts.net port 22: Connection timed out

$ ls -la .mule/ 2>/dev/null; ls .mule/tasks/ 2>/dev/null; ls bench tools 2>/dev/null
total 1416
drwxr-xr-x  7 root root    4096 Sep  7 20:43 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      13 Sep  7 20:43 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep  7 19:00 .heartbeat
-rw-r--r--  1 root root       0 Sep  7 20:43 .last-output
-rw-r--r--  1 root root       0 Sep  7 20:43 .lock
-rw-r--r--  1 root root      12 Sep  7 20:13 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 1311751 Sep  7 20:43 cron.log
drwxr-xr-x  2 root root    4096 Sep  6 20:39 failed
-rw-r--r--  1 root root   20080 Sep  7 19:48 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    2493 Sep  7 19:00 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep  7 20:15 reviews
-rw-r--r--  1 root root     690 Sep  7 20:43 task.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 tasks
-rw-r--r--  1 root root   22903 Sep  7 20:15 todo.md
issue-16.md
issue-23.md
issue-24.md
issue-25.md
issue-26.md
issue-27.md
bench:
README.md
astrate
client.go
client_test.go
connstorm.go
device.go
go.mod
go.sum
ingest.go
interfaces
main.go
provision.go
query.go
scripts
state.go
stats.go

tools:
bsoncapture
mule-plan-day.sh
mule-setup-pi.sh
mule-survey.sh
mule.sh
reconcile.sh

The Legion Go probe timed out — it's off. Let me check whether the repo already holds an answer to the honour-vs-refuse question (prior research/notes), so I know whether this task can proceed at all offline.
✱ Grep "snapshot.*(since|limit|descending)|(since|limit|descending).*snapshot" in . · 1 match

$ ls .mule/recipes/ .mule/reports/ 2>/dev/null | head -60
.mule/recipes/:
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

.mule/reports/:
hygiene.md
log

✱ Grep "snapshot" in .mule · 100 matches
→ Read .mule/reviews/appengine-2026-09-07.md
→ Read .mule/recipes/legion-go.md
→ Read .mule/recipes/astarte-upstream.md
→ Read .mule/config
→ Read .mule/todo.md [limit=30, offset=108]
The Legion Go is off (connection timed out). This task's decision — honour vs refuse — hinges entirely on a live upstream probe that only the Legion Go can run (the review file `.mule/reviews/appengine-2026-09-07.md:36-41,54-59` explicitly says the fix must not be chosen blind, and AGENTS.md + `legion-go.md` require me to stop, not work around). Let me record the blocker for Giulio in the standing channel.
→ Read .mule/for-giulio.md [limit=40]
✱ Grep "legion" in tools · 41 matches
→ Read tools/mule.sh [limit=45, offset=362]

$ rg -n '\[!\]([\[]|-)|"- \[!\]"|"\[!\]"' tools/mule.sh | head -20
(no output)

✱ Grep "\[!\]|BLOCKED" in tools · 3 matches
→ Read tools/mule.sh [limit=120, offset=490]
Let me check .mule git tracking to understand what survives the wrapper's revert, and confirm how the queue treats my working tree.

$ git ls-files .mule | sed -n '1,40p'; echo '--- status ---'; git status --porcelain | head; echo '--- mule gitignore ---'; cat .mule/.gitignore
.mule/.gitignore
.mule/MULE.md
.mule/config
.mule/for-giulio.md
.mule/hosts.example
.mule/log.md
.mule/milestones.md
.mule/recipes/astarte-upstream.md
.mule/recipes/benchmarks.md
.mule/recipes/code-review.md
.mule/recipes/cross-project-survey.md
.mule/recipes/docs-sync.md
.mule/recipes/github-issues.md
.mule/recipes/hygiene.md
.mule/recipes/legion-go.md
.mule/recipes/milestones.md
.mule/reports/hygiene.md
.mule/reports/log/20260904T190911Z-purge-properties-compression-capabilityauto-done.md
.mule/reports/log/20260904T195018Z-broker-acl-coldstart-introspection-miss-done.md
.mule/reports/log/20260904T195529Z-issue-93-done.md
.mule/reports/log/20260904T195814Z-broker-disconnect-device-zombie-session-blocked.md
.mule/reports/log/20260904T203003Z-broker-offline-acl-tests-blocked.md
.mule/reports/log/20260904T210532Z-broker-onconnect-doc-comment-blocked.md
.mule/reports/log/20260905T110211Z-empty-introspection-verification-blocked.md
.mule/reports/log/20260905T111428Z-probe-trigger-install-notification-delay-blocked.md
.mule/reports/log/20260905T112539Z-compat-note-v132-done.md
.mule/reports/log/20260905T165512Z-recipe-docs-sync-proposed.md
.mule/reports/log/20260905T171755Z-docs-sync-pairing-health-path-done.md
.mule/reports/log/20260905T172234Z-docs-sync-pairing-register-404-done.md
.mule/reports/log/20260905T173657Z-recipe-hygiene-proposed.md
.mule/reports/log/20260905T174528Z-swagger-httptest-coverage-done.md
.mule/reports/log/20260905T181053Z-recipe-milestones-proposed.md
.mule/reports/log/20260905T181219Z-recipe-github-issues-proposed.md
.mule/reports/log/20260905T184625Z-recipe-astarte-upstream-proposed.md
.mule/reports/log/20260905T184903Z-probe-property-resend-encoding-blocked.md
.mule/reports/log/20260905T190855Z-compat-note-v133-blocked.md
.mule/reports/log/20260905T192010Z-recipe-code-review-proposed.md
.mule/reports/log/20260905T192911Z-flow-validate-source-sink-blocked.md
.mule/reports/log/20260905T201529Z-flow-validate-dead-source-sink-recompute-done.md
.mule/reports/log/20260906T112006Z-recipe-docs-sync-proposed.md
--- status ---
--- mule gitignore ---
# The mule's own working files. None of these belong in git, and — more sharply — any of
# them that is not listed here makes the tree dirty, which makes `tick` refuse to run and so
# turns the timer into a permanent silent no-op.
failed/
task.md
.timeout
.lock
.budget
.bench-linux-amd64
.last-output
.rotation
.heartbeat
.alarmed
.impl.patch
hosts
research/

→ Read tools/mule.sh [limit=45, offset=445]

# Todos
[✓] Read data.go + docs to understand the snapshot branch and param checks
[✓] Probe upstream on the Legion Go (required to decide honour-vs-refuse)
[ ] Implement the chosen behaviour
[ ] Write a T1 test pinning the behaviour
[ ] Run gate

Blocked before any code: the task's decision gate is a live upstream probe that only the Legion Go can run, and it is off.

**What I did:** read the named files (`data.go:138-148` snapshot branch, `docs/api/astarte_appengine_api.yaml:577-582`, the `[!]`/`[legion]` mechanics in `tools/mule.sh`, and `.mule/reviews/appengine-2026-09-07.md` which is the source of this task). Then probed the Legion Go once, cheaply, per AGENTS.md:

`ssh -o ConnectTimeout=5 legion 'echo ok'` -> fail (connection timed out to `legion.tail4db093.ts.net:22`).

**Why I stopped:** the task itself says "Probe upstream on the Legion Go to decide honour-vs-refuse", and the review file confirms it (`.mule/reviews/appengine-2026-09-07.md:36-41,54-59`): "Snapshot refusing vs honouring is a design call I cannot make blind... the right fix is refuse-422 or honour-the-window must be confirmed against upstream on the Legion first." Picking either branch without the probe violates MULE.md's "Never invent a behaviour" and AGENTS.md's "do not degrade into workarounds". There is no independently useful slice I can do offline — both fixes and the docs edit depend on the probe outcome. `gofmt`/`go test` gates are moot since the tree is untouched.

---

Done: blocked — the Legion Go is off and no code should be written for this decision without the live upstream probe it names
Files: (none changed)
Verified: `ssh -o ConnectTimeout=5 legion 'echo ok'` -> fail (timeout); no gate run, no code changed
Unsure: honour-vs-refuse is undecided by design until the probe runs; I did not guess
Follow-ups: none appended to .mule/todo.md (it and for-giulio.md are tracked, and a tracked edit on this run would be mis-recorded as the task "done" by the wrapper). Please power on the Legion Go (needs the Docker upstream Astarte stack) so this task can run; and consider retagging the line `[legion]` in .mule/todo.md:118 — `tools/mule.sh first_open` skips `[legion]` lines while the Go is asleep and runs the next runnable task instead of burning a run on this probe each time.
