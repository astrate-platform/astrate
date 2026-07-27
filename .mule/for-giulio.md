# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **The Pi cannot run the race detector**, so the unattended gate is weaker than yours.
  ThreadSanitizer needs a 48-bit VMA; the DietPi kernel is built with 39 (`FATAL:
  ThreadSanitizer: unsupported VMA range / Found 39 - Supported 48`). Measured 2026-07-27.
  The gate there is `go vet ./... && go test ./...`, green in ~3m over 20 packages.
  Consequence: **do not queue concurrency work for the mule** — nothing on that machine can
  catch a data race. Two ways out if you want them, both your call: rebuild the Pi kernel
  with 48-bit VA, or install Go on the Legion Go (`pacman -S go`, x86_64, no VMA problem)
  and make the periodic `[legion]` race-check the real concurrency gate.
- **golangci-lint is not installed on the Pi**, so the mule's second gate is silently absent
  there — `gofmt` still runs, the linter does not. `go install
  github.com/golangci/golangci-lint/cmd/golangci-lint@<the pinned version>` on the Pi would
  close it; I did not pick a version for you, since the pin is a decision.
- **`/root/astrate` on the Pi has uncommitted work** (`cmd/astrate/main.go`, `docs/embed.go`,
  `docs/handoff/phase-2-*.md`, `docs/api/astrate_native_api.yaml`) from an earlier session on
  that machine. The mule does not touch it — it uses its own clone at `/root/astrate-mule` —
  but you may want to rescue or discard it.
