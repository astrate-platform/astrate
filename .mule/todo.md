# Mule queue

One task per line. `- [ ]` open, `- [x]` done, `- [!]` blocked (the script writes the reason).
The mule runs the topmost open line. Reorder freely; the order is the priority.

If a line needs more detail than fits, write `.mule/tasks/<slug>.md` and keep the line short.

    tools/mule.sh add "<slug>: <outcome>"
    tools/mule.sh loop

A line tagged `[legion]` needs the Legion Go and is skipped automatically while it is asleep;
the queue moves on to the next runnable line rather than stalling behind it.

A line tagged `[readonly]` is a **standing check**, not a piece of work. It verifies something
and is expected to produce no diff at all — so an empty diff is success, not the `wrote
nothing` failure it would be for any other line. It never gets ticked off, because a gate is
never done. Its output lands in `.mule/reports/<slug>.md` with the sha it ran against, and it
is skipped while that sha is still HEAD: re-checking code that has not moved buys nothing and
spends a call on a free provider.

**`mule/queue` is authoritative for this file.** The mule ticks on the Pi and writes its
progress there, so editing the queue on `main` resurrects completed tasks the moment the two
are merged — which has happened once. Add and reorder tasks on `mule/queue`; let them reach
`main` only when that branch is merged.

A line marked `- [~]` is **parked**: real work, but not tick-sized. The mule only ever
picks up `- [ ]`. Benchmark runs live here because a single tier is 5-20 minutes of ingest
alone -- past the per-task budget -- and because they want someone watching. Run one by hand:

    tools/mule.sh legion bench-push
    ssh legion 'cd ~/astrate/bench && ./scripts/run-tier.sh small astrate -base-url ... -housekeeping-key ...'

## Where tasks come from

**This file is not the whole queue, and for real work it is not even the main part of it.**

The queue is: the standing lines below, plus **every open GitHub issue labelled `mule`**.
Issues are read live on each tick and are never copied into this file — a copy would be a
second place the same fact lives, on a branch the mule commits to and you edit on `main`,
and that produced three merge conflicts in one afternoon.

**To give the mule work, file an issue and label it `mule`.** From anywhere, by anyone,
including another model with repo access. No SSH, no editing this file:

    gh issue create --label mule --title "<slug>: <outcome>" --body "<the detail>"

Labels on the issue are the tags: `legion` and `readonly` mean what `[legion]` and
`[readonly]` mean here. State lives on the issue, as labels, because there is exactly one
copy of it there:

| label          | meaning                                                          |
|----------------|------------------------------------------------------------------|
| `mule`         | queued                                                            |
| `mule-review`  | the mule pushed something; **it is not merged and not reviewed**  |
| `mule-blocked` | it tried and could not; re-label `mule` to try again              |

The mule never closes an issue. Whether the work actually resolves it is a judgement about
intent, which is the reviewer's call.

When both sources are empty a tick runs a **proposal recipe** instead, rotating through
`github-issues`, `astarte-upstream`, `code-review`, `docs-sync`, `hygiene` so it cannot get
stuck re-reviewing the same code. Lines it invents are tagged `[auto]`: nobody approved those.
A refill costs a tick from the daily budget and never runs what it just proposed — the lines
sit for one tick, which is your window to cut a bad one. `MULE_NO_REFILL=1` turns it off.

## Nothing merges on its own

Everything lands on `mule/queue`. The gates prove a change compiles, passes the tests, ships
a test that fails without it, and touches no frozen file — none of which means the change is
worth having. Before any of it reaches `main`:

    bash tools/mule.sh review

## Queue

- [x] bench-tiers: create the tiered benchmark definitions per .mule/recipes/benchmarks.md (first run only — this task builds the harness, it does not run it)
- [~] bench-small-astrate: run `bench/scripts/run-tier.sh small astrate` against local Astrate, commit results (two runs minimum) [legion]
- [~] bench-medium-astrate: run `bench/scripts/run-tier.sh medium astrate` against local Astrate, commit results (two runs minimum) [legion]
- [~] bench-small-astarte: run `bench/scripts/run-tier.sh small astarte` against local Astarte, commit results (two runs minimum) [legion]
- [~] bench-medium-astarte: run `bench/scripts/run-tier.sh medium astarte` against local Astarte, commit results (two runs minimum) [legion]
- [~] bench-big-astrate [legion]: run `bench/scripts/run-tier.sh big astrate` against Legion Go Astrate, commit results (two runs minimum)
- [~] bench-giant-astrate [legion]: run `bench/scripts/run-tier.sh giant astrate` against Legion Go Astrate, commit results (two runs minimum)
- [~] bench-big-astarte [legion]: run `bench/scripts/run-tier.sh big astarte` against Legion Go Astarte, commit results (two runs minimum)
- [~] bench-giant-astarte [legion]: run `bench/scripts/run-tier.sh giant astarte` against Legion Go Astarte, commit results (two runs minimum)
- [ ] race-check: on the Legion Go, `cd ~/astrate && git fetch -q && git merge --ff-only -q origin/main && go test -race ./...` (~40s). Report any failure to .mule/for-giulio.md with the full race report. This is the only race coverage that exists — the Pi cannot run -race. [legion] [readonly]
- [x] store-realm-cascade-policies: in `internal/store/realms_test.go` `CascadeDelete`, add `trigger_policies` to the post-delete verification loop (query `SELECT count(*) FROM trigger_policies WHERE realm_id = $1`). The migration 000006 has ON DELETE CASCADE but nothing asserts on it. [auto]
- [x] store-alias-lowest-id: in `internal/store/devices_test.go`, add a subtest that registers two devices in the same realm, sets the same alias tag on both, and asserts `GetDeviceByAlias` returns the one with the lower UUID. The SQL uses `ORDER BY id LIMIT 1` but no test proves it. [auto]
- [x] store-delete-device-objects: in `internal/store/devices_test.go` `StatsAndDelete`, insert object datastream rows for the device before deleting it, and assert they are gone after the delete. `DeleteDevice` explicitly sweeps `object_datastreams` but the test only checks individual rows. [auto]
