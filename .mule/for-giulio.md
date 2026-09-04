# For Giulio

The escalation channel. The mule writes here instead of acting whenever something needs a
**decision** rather than typing: a design choice, a dependency bump, a schema or protocol
change, a contradiction between the code and the frozen spec, a docs page that needs your
voice.

One line each, newest at the top, with the evidence (file:line, tag, CVE) inline. Delete a
line once you have dealt with it — this file is a queue, not a log.

---

- **Group-scoped triggers (`group_name` on device/data triggers) compile but never match**
  (`internal/engine/triggers/match.go:11-12`). Decision deferred, tied to issue #17
  (group-WATCH-path reconciliation, trickle work, not mule): whatever group-membership
  mechanism comes out of that phase should also report the perf cost for this decision —
  noted in a comment on #17 so it isn't benchmarked twice. (Cross-project survey,
  2026-07-27, `.mule/research/survey-2026-07-27.md` source 4.)

- **#78 FDO device onboarding — milestone-4.0, investigation phase.** Too large for a single
  mule task; the investigation work (reading upstream's TO2 handling, inventorying endpoints,
  schema and keys) is a multi-session project. Parked until the v3.0 queue clears and this
  becomes the next milestone target. Promoted from parked to a milestone-4.0 candidate by
  your decision on 2026-08-23: zero-touch onboarding is strategic for commercial viability.
  Scope frozen on the issue — owner-side TO1/TO2 in our Pairing service only (last mile, like
  upstream), reuse fdo-rs for manufacturing/rendezvous, acceptance is the official
  `astarte-device-fdo-rust` SDK completing onboarding against Astrate, docs a first-class
  deliverable. When v3.0 is marked DONE, draft the v4.0 section of `.mule/milestones.md` with
  this investigation as its first item (#78 carries the full verified context).

- **#1 is never to be raised again.** Your standing instruction, 2026-09-04: it stays open
  permanently and is not a candidate for closing, triage, or a for-giulio entry. Do not
  propose it again.
