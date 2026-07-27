# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **`device_deletion_started`/`device_deletion_finished` trigger events are not emitted.**
  Astrate performs synchronous device deletion (`COMPATIBILITY.md` deviation 9); upstream
  fires both events around an async deletion (`060-triggers.md:109-110`). Decision needed:
  emit both back-to-back for compatibility, emit only `_finished`, or document the divergence
  as deliberate and do nothing. (Cross-project survey, 2026-07-27,
  `.mule/research/survey-2026-07-27.md` source 4.)
- **Mustache trigger-action templates are accepted but not rendered.** Upstream
  (`060-triggers.md:516-543`) renders HTTP action bodies from a `template_type: "mustache"`
  template with access to realm/device/trigger/event fields; Astrate always sends the default
  JSON envelope instead (`internal/engine/triggers/actions.go:46,69`, already noted as
  "unsupported"). Implementing it means picking and pinning a Go Mustache library — a
  dependency decision `docs/ROADMAP.md` §0.1 freezes deliberately. Decision needed: add the
  dependency and implement it, or keep documenting the divergence. (Same survey, source 4.)
- **`value_change`/`value_change_applied`/`path_created`/`path_removed`/`value_stored` trigger
  types compile but never fire** (`internal/engine/triggers/match.go:30-42`, already documented
  as a deliberate v1 scope boundary). Upstream evaluates all of them. Implementing any of them
  needs the engine to read previous state before matching (a DB read per candidate message,
  or a cache) — a real design and performance decision, not a mechanical gap-fill. Decision
  needed: implement some/all in a future milestone, or keep as documented v1 scope. (Same
  survey, source 4.)
- **Group-scoped triggers (`group_name` on device/data triggers) compile but never match**
  (`internal/engine/triggers/match.go:11-12`). Implementing this needs device group-membership
  lookups at trigger-evaluation time — related to the group-WATCH-paths gap (issue #17), since
  both need groups to actually exist in a test realm to verify against. Decision needed: scope
  this in, or keep documenting the divergence. (Same survey, source 4.)

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
