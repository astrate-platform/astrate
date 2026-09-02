# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **`govulncheck` is not installed on the Pi**, so the hygiene recipe's highest-priority check (reachable dependency vulnerabilities) could not run here. `go install golang.org/x/vuln/cmd/govulncheck@latest` would close it — I did not pick a version because the pin is a decision. The `go list -m -u` sweep shows only version-skew on transitive deps (cloud/azure SDKs, cel/expr, clickhouse); none corresponds to an identified fix this repo needs, so I proposed no bumps. (Hygiene recipe run, 2026-09-02.)

---

- **`docs/site/appengine-api.md:87-88` documents `GET` and `DELETE /appengine/v1/<realm>/groups/<name>`**, but no such routes exist in `internal/appengine/http.go` (only `/groups/{group}/devices` and the nested device ops — there is no `GET` or `DELETE` on `/groups/{group}`). The endpoint the page describes is absent from the code. Decide: drop the two lines, mark them not-yet-implemented, or have the code grow them. (Docs-sync recipe run, 2026-09-02.)

---

- **`Router.Submit` TOCTOU on `closed` flag** (`internal/flow/router.go:113-120`): Submit
  reads `r.closed` under the mutex, drops the lock, then sends on the channel. A concurrent
  `Drain` could close the channel between the unlock and the send, causing a send-on-closed
  channel panic in the caller (not recovered by `processOne`'s defer). Fix options: (a) hold
  the lock during channel send (hot-path perf hit), (b) use a select with `r.quit` instead of
  sending on a closed channel, (c) accept the risk and document that Submit must not be called
  after Drain. Needs a design call — not proposing a task. (Code review 2026-09-02.)

---

- **#87 `lua_map` — needs embedded Lua runtime, parked.** Design/implementation decision: embedding a Lua VM in Astrate is not machine-checkable by the mule. Consider closing if Lua flow support is not on any active roadmap.
- **#78 FDO device onboarding — milestone-4.0, investigation phase.** Too large for a single mule task; the investigation work (reading upstream's TO2 handling, inventorying endpoints/schema/keys) is a multi-session project. Parking for now until the v3.0 queue clears and this becomes the next milestone target.
- **#1 stale — "Provide an Open Source IoT Platform unironically" (wontfix).** Has been open since the repo's founding with no activity. Consider closing.

---

- **Milestones recipe has nothing to file: v2.0 is already DONE, and `mule/queue` is 120
  commits behind `origin/main`, so this branch's `.mule/milestones.md` is stale.** My branch
  still shows v2.0 "not started" (and v3.0 as CLEA TBD), so the recipe pointed me at v2.0 —
  but on `main` v2.0 is marked **DONE (2026-07-29)**: every `milestone-2.0` issue is CLOSED
  (#23–#27, #37, #39, #40, #41, #42, #43), and the factory/catalog/container-block/durable
  named flows/`/flow/v1` API all live on `main` (`git show origin/main:.mule/milestones.md`).
  `gh issue list --label milestone-2.0 --state open` is empty and there are no new v2.0 gaps
  to file, so per the recipe's step 5 I filed **nothing** rather than duplicate closed work.
  **Real current target is v3.0** ("upstream Astarte feature parity 1.3/1.4", scoped on main
  2026-08-22): the first recipe job there is to triage `upstream-parity` backlog **#47–#89**
  (source `.mule/research/upstream-parity-2026-08-22.md`) into a plan. But I should not start
  that triage from this stale branch. Action needed: **merge/rebuild `mule/queue` on top of
  `origin/main`**, then the next milestones run targets v3.0 with the correct research doc.
  (Milestones recipe run, 2026-09-01.)

---

- ~~`device_deletion_started`/`device_deletion_finished` trigger events are not emitted~~ —
  **decided 2026-07-27: emit both, back-to-back, around the synchronous delete.** Filed as
  issue #21 (`mule`). (Cross-project survey, 2026-07-27,
  `.mule/research/survey-2026-07-27.md` source 4.)
- ~~Mustache trigger-action templates are accepted but not rendered~~ — **decided
  2026-07-27: implement it.** Guiding principle clarified: Astarte compatibility means
  SDK/wire compatibility, not minimum dependency count — Astrate is allowed to be a
  compatible *superset*. Library picked: `github.com/cbroglie/mustache`. Filed as issue #22
  (`mule`). (Same survey, source 4.)
- **`value_change`/`value_change_applied`/`path_created`/`path_removed`/`value_stored` trigger
  types compile but never fire** (`internal/engine/triggers/match.go:30-42`). Decision
  deferred pending data: issue #20 (`mule`, `readonly`) asks Big Pickle to benchmark the
  ingest-path cost of a previous-value lookup on the Legion Go before this gets decided one
  way or the other. (Same survey, source 4.)
- **Group-scoped triggers (`group_name` on device/data triggers) compile but never match**
  (`internal/engine/triggers/match.go:11-12`). Decision deferred, tied to issue #17
  (group-WATCH-path reconciliation, trickle work, not mule): whatever group-membership
  mechanism comes out of that phase should also report the perf cost for this decision —
  noted in a comment on #17 so it isn't benchmarked twice. (Same survey, source 4.)

---

- ~~The Pi cannot run the race detector~~ — **resolved 2026-07-27** by installing Go 1.26.5
  as a userland toolchain on the Legion Go (`~/.local/go`, no root, `rm -rf` to undo). The
  Pi still cannot run `-race` (39-bit VMA kernel vs the 48 ThreadSanitizer needs), so its
  gate remains `go vet ./... && go test ./...` — but race coverage now exists on the Legion
  Go, where the full suite runs clean in ~40s on 16 cores. The standing `race-check` task is
  the concurrency gate. Concurrency work is queueable again, provided the race-check runs
  after it.
- **golangci-lint is not installed on the Pi**, so the mule's second gate is silently absent
  there — `gofmt` still runs, the linter does not. `go install
  github.com/golangci/golangci-lint/cmd/golangci-lint@<the pinned version>` on the Pi would
  close it; I did not pick a version for you, since the pin is a decision.
- ~~`/root/astrate` on the Pi has uncommitted work~~ — **resolved 2026-07-27** with the new
  `tools/reconcile.sh`: rescued onto `origin/wip/DietPi-20260727T171543Z` (pushed, not
  reviewed — read the diff before merging anything from it) and `/root/astrate` is now clean
  on `main`. Also had to set `commit.gpgsign false` locally on that clone first — it had
  signing on with no working gpg-agent for a non-interactive session, same class of problem
  `/root/astrate-mule` already solved the same way (your call, confirmed 2026-07-27: align it
  with the existing unattended-Pi-clone convention rather than fix gpg-agent or bypass signing
  silently inside the script).
- **The mule has been idle 15h.** Filed by the dead-man's switch; see journalctl on the Pi.
