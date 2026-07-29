---
name: trickle
description: Trickle-down delegation — the strong model plans, a mid model splits and supervises, and a cheap local agent (opencode/big-pickle) does the typing in an isolated git worktree. Use when Giulio asks to "delega", "trickle mode", "fallo fare a Big Pickle", "risparmia token", "spezza il lavoro e delega", "delegate this", or wants a plan executed cheaply instead of typed out by the top model.
---

# Trickle-down delegation

| Layer | Who | Job |
| --- | --- | --- |
| **Architect** | top model (Opus), main session | Decide *what* and *why*. Write `plan.md` and the phase specs. Review escalations. Never type boilerplate. |
| **Supervisor** | a Sonnet subagent, per phase or small group | Drive `trickle.sh`, read every diff, run the gates, decide re-round vs. take-over. Escalate judgment calls. |
| **Delegate** | `opencode --agent build` (big-pickle), throwaway worktree | Type what the spec says. Cannot reach the real branch. Not a collaborator with memory — a function you call. |

Mechanics live in **`~/.claude/scripts/trickle.sh`** (`--help`): worktree lifecycle, running
the delegate harmlessly, enforcing the never-touch list and the phase's own file list, the
formatter pass, the gates, landing unstaged.

**Everything written for other agents — specs, feedback, the log, this file — is in
English**, even when the conversation with Giulio is in Italian. Every model in the chain
follows English more reliably, the delegate most of all.

## How to read the rules

Each rule carries its pedigree: `[mechanical]` follows from the tools and never needs
retesting; `[observed · N]` repeated across N rounds; `[inferred · N]` one or two data points
plus a story; `[untested]` proposed, no evidence. Non-mechanical rules carry a falsifier —
**if a run falsifies one, change it.** Evidence sits in `evidence/`, read only when a rule is
in question; never copy it back here.

## Triage: what to delegate, and to which tier

Two questions. Difficulty is not one of them.

1. **Is the design decided?** Does implementing it require choosing something with consequences?
2. **Is correctness machine-checkable afterwards** — or does it take a human reading the diff
   against the phase's purpose?

| | machine-checkable | needs a human read |
| --- | --- | --- |
| **design decided** | delegate the phase; the gate is the review | delegate the typing, architect reads the diff before landing |
| **design open** | architect decides, then it is the cell on the left | architect does it |

- **Judge on bindability, not complexity.** [inferred · 1] The most intricate phase of the
  one full run landed in a single round; a simpler one needed two, because it existed to
  produce a *log line* and nothing asserted on log output. *Falsified by: rounds tracking
  size across projects while checkability does not.*
- **New file, no callers** → safest shape there is. [observed · 4]
- **Edits control flow in an existing function** → expect structural residue; say so in the
  spec. [observed · 2]
- **For an edits-flow phase, enumerate the branches as a numbered list — in the spec *and*
  again as named checks in the supervisor's brief.** [inferred · 1] The delegate's
  characteristic failure here is getting most branches right and silently dropping one
  (`executor-strategy`). The one edits-flow phase specced this way landed first try with
  every branch present, including the guard flip that was the exact omission-shaped risk:
  both layers were hunting for a *specific* missing branch rather than reading for general
  correctness. *Falsified by: an enumerated edits-flow phase dropping a branch anyway.*
- **Output consumed elsewhere** (log, metric, caller) → bind it with an assertion or budget
  an architect read. No third option. [inferred · 1]
- **Ask what the gate can *reach*, not only whether one exists.** [observed · 2] The delegate has
  no network, so a phase whose only meaningful check is "run it against a live external system" has
  an unreachable gate and is architect work — however machine-checkable it looks on paper. M12's
  plan had a recording phase down as "delegate the runner"; the runner could not have been run, and
  every finding it produced (an asynchronous realm creation surfacing as a bare 403, an
  existing-realm status nobody had guessed) came from the live run rather than from reading code.
  Recording, ops and smoke phases fail this test together. **Confirmed again in M12-07, where it did
  the useful work of *splitting* a phase rather than rejecting it:** the plan sized 07 as one phase
  over three files, two on the never-touch list and one whose only real gate is a live database. Cut
  into an architect half (image, volume, migration) and a delegate half (the two Go files, with the
  probed SQL facts written into the spec), the delegate half landed in one 115s round with no
  defects. So the question is not only "is this delegable?" but "which *part* of it has a gate the
  delegate can reach?" — the answer is often a clean seam rather than a refusal.
  **Corollary, and the run's sharpest finding: when the fix is "make X reachable by the tests",
  enumerate every route the tests take to X.** M12-07 existed because a downsampling branch could not
  execute on the pinned database image, and it changed the compose file and the two CI `image:` lines
  — then almost shipped with the branch still unreachable in CI, because the CI job that runs T2 has
  no service container at all and obtains its database through testcontainers, pinned separately in
  the test helper. Two of three routes fixed reads exactly like three, and the symptom of the third
  is a *skipped* test, which is green. Found by asking how CI actually gets a database rather than by
  any gate. *Falsified by: a reachability fix that covers every route on the first, obvious pass.*
  *Falsified by: a delegated phase verified entirely by a gate that never reaches the system under
  test.*
- **Against anything you do not control — someone else's running system, or a dependency's
  semantics — a scratch probe is the cheap instrument.** [observed · 3] Two of M12-04a's three most useful facts — the 2.5s realm propagation delay, the
  422 existing-realm status — came from twenty-line programs written to answer one question each,
  not from reading either codebase. Both were invisible from our source and both changed the
  implementation. M12-04b repeated it at larger scale: five throwaway programs established the
  socket path, the claim grammar, the trigger payload shape and two failure modes *before* a line
  of the real recorder was written, and every one of those answers contradicted a plausible reading
  of the code. **M12-07 extends the rule past remote systems to dependencies:** six psql one-liners
  established `lttb()`'s real contract, and two of the six contradict the obvious reading — it
  refuses a resolution below 3 (so a legal `downsample_to=2` becomes a 500 without an explicit
  floor) and it sorts its input internally (so `LIMIT` inside the aggregate changes the resolution
  rather than trimming the output). Neither is visible from our source; handing them to the delegate
  as *measurements* rather than leaving them as its guesses is a large part of why that phase landed
  first try. *Falsified by: probe programs that repeatedly confirm what a careful read of the
  code already said.*
- **A recorded observation is not a fact until it has been observed twice.** [observed · 1]
  M12-04b's recorder produced a clean, stable-looking table on its first pass and a *different* one
  on its second; either would have been committed as a fixture, and a fixture is precisely the
  artefact later phases stop questioning. Two independent causes, and only one was the recorder's:
  ambient events from the system under test were being attributed to the provocation, and the
  system itself dropped events intermittently. Both are invisible in a single pass by construction.
  So a recording phase repeats each row, records the *count* alongside the value, and states absence
  as "never seen in N tries" rather than "the system emits nothing" — the weaker claim is the only
  one the instrument supports. *Falsified by: a single-pass recording that survives repetition
  unchanged, repeatedly.*
- **In a recorder, hunt for the silent failure that mimics a result.** [observed · 1] The three
  worst bugs in M12-04b all produced coherent, readable, wrong output rather than an error: an idle
  websocket the server closed mid-run (every later row recorded "no event", indistinguishable from
  a finding), a malformed protocol handshake that made *every* case return the same plausible error,
  and first-event-wins attribution that silently credited ambient noise. A recorder's characteristic
  defect is not crashing — it is confidently writing down silence. Make the instrument's own
  failures loud (print the socket closure, assert the handshake, correlate each event to its cause)
  before trusting a row that says nothing happened. *Falsified by: a recorder whose failures
  reliably surface as errors rather than as plausible data.*
- **When the consumer is someone else's code, only running that code is verification.**
  [inferred · 1] M11's tests asserted the exact frames Astrate emits and were green for two
  phases while the Dashboard silently discarded every `device_error`, because it validates one
  field against a closed enum ours was never in. A test written from our side can only assert
  what we *send*; it cannot notice that nobody will accept it. Any plan whose goal is
  compatibility with a shipped client needs a phase that runs the client — and a checkable one,
  since the failure surfaced as a console exception and a missing row, not an error.
  *Falsified by: a client-acceptance defect caught by our own test suite first.*

Never delegate regardless: prose in Giulio's voice, schema or architecture decisions, git
history, signatures, credentials. Per-project extras go in `TRICKLE_NEVER`.

## Setup, once per project

```bash
bash ~/.claude/scripts/trickle.sh init      # writes .trickle/config and .trickle/log.md
```

Fill in `TRICKLE_TEST_CMD` (the acceptance gate — without one the mode loses most of its
safety), `TRICKLE_NEVER`, `TRICKLE_FIX_CMD` (formatter + linter autofix).

## Every session, before writing a single spec

```bash
bash ~/.claude/scripts/trickle.sh preflight
```

[observed · 4 distinct failures in one session] Environment failures have cost more than
every model in the chain put together — missing toolchains, expired gpg agents, containers
removed by an update, provider outages. Preflight asks all of it at once and captures the
test/lint baseline, so a later failure is provably the delegate's. *Falsified by: preflight
passing and the environment still eating the run.*

**Preflight clears the *delegate* gate, not the phase's own verification.** [observed · 1] M12-07
passed preflight and then lost a step to a `docker pull` that **reported exit 0 while failing**,
because Docker Desktop was down. `TRICKLE_NEEDS_DOCKER=""` was honest — T1 genuinely needs no
Docker — but the phase was verified at T2, and nothing checks the environment *that* needs. When a
phase's real gate is heavier than `TRICKLE_TEST_CMD`, check its environment yourself before
speccing, and do not trust a tool's exit code over its output. *Falsified by: a T2-verified phase
whose environment preflight already covered.*

**A lint baseline warning may be a ghost of a dropped worktree.** [observed · 3] Build-cache
linters key on absolute paths, so after `drop` they keep reporting issues in files that no
longer exist — golangci-lint warned about two on a tree that was clean, and did it again
after every landing — and again in M12-07, two `gosec` issues against a `cmd/astrate/main.go` inside
the deleted worktree path. Clear the linter's cache before believing a non-empty baseline, and
before blaming a delegate for issues it did not introduce. *Falsified by: a cache-cleared
baseline still reporting phantom files.*

It never prompts and never waits on one: the gpg check runs with `--pinentry-mode error`
under a timeout, so a locked key is *reported* rather than reproduced. An unattended TouchID
dialog nobody is watching blocks forever, which is the exact failure preflight exists to
catch. [mechanical]

## The architect's output

A plan on disk, not in context — that is where the saving comes from.

- `plan.md`, **up front**: frozen design decisions, the phase list with a one-line outcome
  each, ordering and dependencies, the load-bearing invariants.
- `phases/NN-<slug>.md`, **one at a time, written after reading the previous diff.**
  [inferred · 1] Writing them all up front is waterfall: later specs get composed in
  ignorance of what earlier rounds just taught, and you pay for specs you may never reach.
  *Falsified by: up-front specs landing as cleanly as just-in-time ones.*

```markdown
# Phase — <one line: the outcome, not the activity>

<!-- trickle-allow: path/to/the.go path/to/the_test.go -->

## Context
Why this exists and where it fits. Then the files to read first, each with a one-line reason
and **line ranges**. Say what NOT to read.

## What to do
The concrete change. Name files to create or modify. Be specific about shape and naming.

## Constraints
What not to touch. "If you think X is broken, say so in your summary instead of fixing it."
No network.

## Acceptance criteria
Commands the delegate can run itself inside the worktree. Never reference /tmp or any path
outside it — the sandbox denies them. Name pre-existing failures so they are not mistaken
for its own.
```

`trickle-allow` is the same file list the phase names in prose, in a form `run` enforces —
straying becomes an exit code instead of something a reviewer must catch.

**What decides whether a round lands first try:**

- **A change to an output shape invalidates every committed fixture of that shape — name them
  in `trickle-allow`, or the round dies at a gate it cannot pass.** [inferred · 1] M11's
  `error-name` changed what `NewDeviceErrorEvent` emits, which by construction made two
  committed goldens stale; the spec named only source files, so the delegate did the correct
  thing — update the goldens — and `run` correctly called that a file-list violation. The
  delegate was trapped: it could pass the gate or obey the list, not both, and neither choice
  was a mistake. Before writing the allow-list, grep for fixtures of whatever shape the phase
  changes. *Falsified by: a shape-changing phase that omits its fixtures and lands clean.*
- **Two files written, at most three read.** [observed · 15] One contrary data point:
  `error-name` had five files in its final allow-list (two new, one edited function, two
  goldens) and landed in a single 168s round — but three of the five were mechanical
  consequences, and the spec carried the edited function's replacement text *verbatim*, so the
  delegate's real design surface was still one file. Count the decisions, not the paths. Was "one file"; the falsifier
  fired — four two-file phases have now landed clean in one round (`ttl-capacity`,
  `parse-downsample`, `series-span`, and the T2-paired `wire-downsample`), three of them in
  under 70s. What actually binds is context, not file count: past ~30–40k tokens big-pickle
  degrades, and an overrunning phase produces *nothing at all* after burning the budget. The
  cheap way to stay under two is for the architect to hand-write any one-line change a third
  file needs — a struct field, a constant — before the round starts. *Falsified by: a
  three-file phase landing clean in one round, repeatedly.*
  **Corollary: split a phase the plan sized too big, before speccing it.** [inferred · 1] M11's
  phase 05 was one ~250-line file plus a full wire-session test; split into 05a (socket, join,
  leave, heartbeat) and 05b (watch, authorization, event pump), *both halves landed clean in one
  round*. Cut along the seam where the judgment sits, so one half stays purely machine-checkable
  and only the other needs the careful read. A plan written before the code exists is allowed to
  be wrong about size; the architect re-sizing it at spec time is cheaper than an exit 4, which
  yields nothing at all. *Falsified by: a split phase whose halves cost more rounds together
  than the whole would have.*
- **Line ranges, not just filenames.** [observed · 14] Every one-round success had them.
  *Falsified by: range-less specs landing at the same rate.*
- **State the invariant, not only the change.** [observed · 8] "The no-policy path must be
  byte-for-byte unchanged" held across four consecutive phases editing one function, and
  "the span must ignore Limit and Descending" / "a request without the parameter behaves
  exactly as today" held across three more. *Falsified by: a stated invariant broken anyway.*
- **If the phase writes tests, spec the harness's *gaps*, not just the test.** [inferred · 1]
  The delegate cannot cheaply discover what a rig leaves unwired, and it will spend the whole
  budget grepping to find out: `bus-events` hit exit 4 having written no production code at
  all, because the spec assumed a rig that builds through the package's public constructor
  when it actually builds through a lower-level one and leaves four fields nil. Name the setup
  lines the test needs before its first assertion, or keep the phase to production code and
  write the test yourself. *Falsified by: a test-writing phase landing clean against a rig the
  spec never characterised.*
- **If the phase asserts that bad input is *rejected*, supply the adversarial inputs yourself.**
  [observed · 2] The delegate picks rejection cases that fail for an incidental reason and
  reads that as proof: `protocol` R1 asserted wrong-length frames were refused while the code
  never checked length, because every case it chose also had a type error further down; the
  same shape sank `bucket-arith` R1, whose property check held whether or not the rule existed.
  Applied properly it works: `room-registry`'s three supplied rows each turned red under an
  anchored mutation of exactly the rule they named.
  Write the rows into the spec, each one *valid in every respect except the rule under test*,
  and require the assertion to name the rule (the length, the reason) rather than merely
  demanding a non-nil error.
  **Pair a rejection row with an acceptance whenever a blanket refusal would also pass it.**
  [inferred · 1] "This watch is refused because the token does not authorize that device" is
  proved by the refusal only if some *other* watch on the same token is accepted — otherwise a
  token that denies everything, or a handler that refuses every watch, passes the row. Specced
  that way in M11's `ws-watch`, the supervisor's mutation of the authorization check turned
  exactly that row red and nothing else. *Falsified by: an unpaired rejection row that still
  isolates its rule under mutation.* *Falsified by: a delegate-chosen rejection table that isolates
  the rule.*
- **If the gate is `-race`, spec the interleaving — the flag alone proves nothing.**
  [inferred · 1] A race detector only reports schedules the tests actually run. In M11's
  `ws-watch` every specced test pushed events with the room quiescent, so deleting the pump's
  mutex passed `-race -count=5` clean; a test that drives events *concurrently with* a leave
  trips two data races on the same mutation. When a phase adds a goroutine, name the concurrent
  pair the test must produce (this call racing that one), not just the flag to run under.
  *Falsified by: an unspecced concurrent test catching a lock removal anyway.*
- **A phase that exists for observability must be bound by an assertion.** [inferred · 1] The
  delegate computes the value the spec asked for and does not wire it to the place the spec
  said it was for, whenever no test forces the connection. Name the log line verbatim;
  require a test on captured log output, not on a counter. *Falsified by: it wiring up an
  unasserted consumer correctly, repeatedly.*

## The supervisor's loop

```bash
trickle.sh start <name>
trickle.sh run   <name> .trickle/phases/NN-x.md
trickle.sh diff  <name>                # READ THIS
trickle.sh test  <name>                # the acceptance gate
trickle.sh prove <name> <impl file>    # if the phase adds a check
trickle.sh land  <name> && trickle.sh drop <name>
```

`prove` reverts the implementation, keeps the tests, expects the gate to go **red**. Use it
when a phase adds a *check* rather than a feature, and always when the gate only proves the
change compiles — a test never seen to fail is worth very little. [mechanical]

**`prove` cannot be used on a new-file phase** — reverting the only implementation file stops
the test file compiling, and a compile error proves nothing about any individual rule.
[mechanical] There, mutate instead: for each rule you doubt, delete exactly that rule with an
anchor-asserted edit and check the *intended* test goes red on its own. `room-registry` ran five
such mutations in one scripted pass; each named a different rule and each turned exactly one
test red, which is a far sharper result than a whole-file revert could give.

**When a mutation *survives*, the rule is usually fine and the test is usually the problem — find
out which, because "redundant, delete it" is the wrong conclusion and the tempting one.**
[observed · 2] Both instances came from one session and neither was a missing test; both were
tests that looked like they covered a rule and covered something else. In `m12-08a` the guard
normalising an empty `json.RawMessage` survived deletion because the test passed **nil**, and the
stdlib already emits `null` for a nil slice on its own — the test's input could not reach the
branch. In `m12-08b` the `url is required` check survived because an absent URL is *also* refused
by the absolute-URL rule one line later, and the assertion only demanded the message mention
`"url"` — two rules catching the same input, indistinguishable to that row. Fixes: an input that
reaches only the rule under test (a non-nil, zero-length value), and an assertion on the wording
that only the rule under test produces. So the diagnostic is: **can this test's input, and this
test's assertion, tell the mutated rule apart from its neighbour?** If not, sharpen the test — do
not delete the code. *Falsified by: a surviving mutation that really did mark dead code.*

**A mutation that does not compile proves nothing, and this is the most common way a mutation pass
lies to you.** [observed · 3] `m12-07` produced invalid SQL, `m12-08a` left a variable unused with
`if false`. Both looked like a red test and were a build failure. Assert the build succeeds before
reading the result, and prefer mutations that change semantics while keeping every binding used —
`if x.Scheme == "definitely-not-a-scheme"` rather than `if false`. *Falsified by: a compile-failing
mutation that still isolates its rule.*

- **Read the diff yourself.** [observed · 15] The summary is usually accurate and occasionally
  confident about work it did not do. *Falsified by: summaries proving reliable.*
- **Never re-run a spec against a worktree where it already succeeded.** [observed · 3] It
  rewrites rather than refines — longer, worse, regressing what it had right. A second round
  is a *new*, narrow file naming only what is wrong now, and that works: `protocol` R2 fixed a
  real defect in 68s against a worktree R1 had already passed the gate on, touching nothing it
  had got right. *Falsified by: a re-run improving on its own output.*
- **When a mutation's result contradicts the code you just read, suspect the mutation.**
  [observed · 1] A malformed `perl -0pi` edit during `protocol` silently changed different
  semantics than intended and produced a red/green pattern that made no sense against the
  source; redoing it in python gave an unambiguous answer. Anchor-assert before mutating
  (`assert old in s`) so a missed anchor fails loudly instead of quietly proving nothing.
  **Anchor the restore too, and assert uniqueness, not just presence.** [observed · 1] In
  `error-name` a supervisor's restore used a non-unique `replace` and silently altered a second
  identical line elsewhere in the same function, leaving the tree subtly wrong for the *next*
  mutation. A mutation pass that corrupts its own baseline reports nonsense from then on.
- **Exit 4 = TIMEOUT**: killed mid-turn, the worktree holds a fragment. Shrink, restart. [mechanical]
- **Exit 3 = never-touch or phase file list hit.** The exit code is [mechanical]; what to do
  about it is a policy that was mis-tagged as mechanical for four sessions and so never got
  tested. **Diagnose it at the supervisor layer, decide it at the architect layer.** [inferred · 2]
  Confirmed a second time in `error-name`, on the *other* branch of the rule: the flagged paths
  were real tracked files, so the supervisor stopped dead without reading further and reported
  what it saw — which files, the commit that added them, and why the delegate had touched them.
  The architect settled it with one `cat` and widened the allow-list. What made that cheap is
  that the supervisor neither decided nor sent a bare "exit 3, please advise".
  The supervisor has the worktree loaded and can establish what actually happened for a fraction
  of what re-deriving it upstairs costs — an M11 supervisor correctly identified a false exit 3
  (a stray `.trickle-phase.md` from an orphaned round poisoning `round_base`, so trickle's own
  scaffold cleanup read as a deletion outside the allow-list) and the architect confirmed it with
  a single `git status`. What the supervisor must not do is *resolve* it: never land on an exit 3,
  and report the **evidence** rather than the conclusion, so the architect's check stays one
  command rather than a re-investigation.
  **Two hard stops survive, and they are what the rule was protecting:** if the flagged path is a
  real project file, or anything on the never-touch list, stop dead and escalate — do not
  diagnose, do not continue to `diff`. A plausible story exists for both the benign and the
  malign case, and that is the whole difficulty: getting it wrong on a false positive costs one
  round, getting it wrong on a true positive lands a diff that touched something forbidden.
  *Falsified by: a supervisor diagnosis of an exit 3 that the architect's one-command check does
  not settle, or a real never-touch hit reasoned away as a scaffold artefact.*
- **Two rounds, then take over.** [inferred · 1] *Falsified by: third rounds landing clean
  and cheap.*
- **Never run the delegate outside a trickle worktree**, never near `main`. [mechanical]
- **Never end a turn waiting on a background command.** [mechanical] Notifications do not
  reach an agent no longer in its turn: it looks like a report and is a silent stall.
  **Restate this in every supervisor brief; the skill file alone does not carry it.**
  [observed · 1] An M11 supervisor backgrounded `run`, ended its turn, and reported success
  about a worktree that held only the scaffold file — caught solely because the architect
  checked the tree. The next supervisor, given the rule verbatim in its prompt, ran clean.
  Architects should also verify the worktree actually contains the files before believing any
  first report. *Falsified by: a supervisor briefed without it that still runs synchronously.*
- **`Error: No provider available` is an outage** — not a spec problem, not auth
  (`auth list` showing 0 credentials is normal for a zero-auth provider). Wait ~60s, retry
  twice, stop. Log as infra. [observed · 1]
- `land` may be refused for subagents while working from the main session; then the
  supervisor stops after `test` and the architect lands and commits. [observed · 1]

Expect one mechanical nit per round if `TRICKLE_FIX_CMD` is unset — the delegate never
formats or lints its own output. [observed · 8] Nor does it tidy code it has just
restructured: a single-case `switch` left after extracting a branch, an accessor called twice
in adjacent `if`s. Cheap post-landing cleanups, not worth a round. [observed · 2]

**Escalate — don't decide** — when the spec is wrong about the codebase, the delegate flags
something that changes the plan, the gate passes but the diff smells wrong, or a phase blows
its budget twice. Escalate with the diff and a specific question.

**But "don't decide" is not "don't think".** [inferred · 1] The split that actually holds is
*investigation down, judgment up*: the supervisor is the cheapest layer that still has the
worktree, the diff and the delegate's summary in context, so gathering evidence — mutating a
rule to see what turns red, establishing whether a flagged path exists, characterising a failure
— belongs there, and the architect's job is to decide on that evidence. A supervisor that
escalates a bare "exit 3, please advise" has moved the whole investigation upstairs at the most
expensive rate in the chain. **Escalate with evidence, and the escalation gets cheap enough to be
worth doing often.** *Falsified by: supervisor-gathered evidence that repeatedly misleads the
architect's decision.*

## The log

One row per phase in `.trickle/log.md`. The note holds the *why*; the columns are what let a
rule be promoted, demoted, or killed without reasoning from memory.

```
| date | repo | lang | phase | checkable? | shape | spec by | rounds | outcome | secs | note |
```

`checkable?` yes/partial/no — the triage answer, recorded so the triage itself can be tested.
`shape` new-file / additive / edits-flow (combine with `+`). `spec by` the tier that wrote
the spec. `outcome` landed / taken-over / abandoned / infra — infra rows carry 0 rounds, so
the rounds column stays honest about the cheap layer.

**Each session also closes with a short summary paragraph**, and it carries two fields that
exist to test the handoff rule above rather than any phase: how many phases this session ran
and whether it ended **voluntarily or forced** (context exhausted, a blocker, the user
stopping); and — written by the *next* architect, at the top of its own summary — **what it
had to re-derive that the plan did not carry.** "Nothing" is the result that keeps the rule;
a recurring list is what kills it. Without that second field the handoff timing can only ever
be argued from memory, which is the thing this file exists to avoid.

## When to end the session, not just how

A plan usually outlives one session. Deciding *when* to stop and hand the rest to a fresh
architect is the architect's own call, and nothing else will prompt it — so make it explicitly,
at every phase boundary, rather than discovering it when the context is already spent.

**The reframe that makes this checkable: handoff-readiness is a test of the plan, not a
scheduling decision.** [inferred · 1] In this mode the architect's context is disposable by
construction — everything durable is already supposed to be in `plan.md`, the specs, `log.md`
and the project's memory. So "should I hand off?" collapses into "what would the next
architect have to re-derive?" If that answer is not empty, the fix is **not to stay — it is to
write the missing thing down**, which the plan owed anyway. An architect who cannot hand off
cheaply has under-written the plan. *Falsified by: a session that hands off with a plan it
believed complete, and a successor that repeatedly reports missing context anyway.*

So, after each `land` + log + plan update, ask one question: *if this session ended now, what
would the next architect have to re-derive?* Non-empty → write it into the plan. Empty → the
handoff is nearly free, and only cost remains.

Weigh that cost with signals you can actually observe, not a token estimate you cannot:

- **The first summarization is the sharp trigger.** [mechanical] Once the context has been
  compressed, staying carries no context advantage left to protect: you are re-reading a lossy
  digest of what a clean session would read fresh from `plan.md`, which is more precise than
  your recollection of it. Hand off at the next phase boundary.
- **Only ever at a phase boundary.** [mechanical] A worktree in flight, a phase specced but
  unrun, a round awaiting review — each is exactly the "leftover state a next session pays for
  twice" this file warns about below. Land or drop first.
- **The shape of what is coming next.** Browser smokes, e2e log wrangling and large-diff
  reviews load context fast; a new-file phase with a two-method fake barely at all. A heavy
  phase is better started fresh than finished tired.

**The asymmetry decides ties, and it is specific to this mode.** [inferred · 1] Handing off too
early fails visibly and cheaply — the successor reads three files. Staying too long fails
invisibly and gradually: what degrades first is the skeptical diff read. In a mode whose whole
premise is that the review catches what the gates do not, that is the one faculty worth
protecting — and the evidence is in this log, where every defect found so far was found by
reading a diff whose gate was green. **Bias toward handing off; the burden of proof is on
staying.** *Falsified by: sessions that ran long and still caught their defects at the same
rate.*

What is genuinely irreducible is tacit calibration of the delegate's habits. That is an
argument for keeping the evidence-tagged rules above rich enough to carry it, not for staying.

## Closing a run: the handoff command

A run ends when the *next* session can start without re-deriving anything — not when the
last phase lands, and not only when the plan is finished. Do these in order, then stop:

1. commit the code
2. `.trickle/log.md` — one row per phase, notes carrying the why
3. `.trickle/plan.md` — the resume point, or `Status: complete` and what was left open
4. the project's own memory, if it keeps one (a plan finished in `.trickle/` that no other
   agent can see has only half landed)
5. the honesty pass on this file, below

**Then, as the very last thing in the final message, print the command that starts the next
session.** One fenced `bash` block, one line, nothing else in it:

```bash
cd /path/to/repo && claude "riprendi il trickle: .trickle/plan.md, fase 05 — <slug>. preflight per primo."
```

- **`cd` included, copy-pasteable as-is.** The next session starts from wherever the terminal
  happens to be, not from here.
- **Name the next phase and its number.** "Continue where we left off" makes the new session
  read the plan to find out what it is; naming it means the plan is read to find out *how*.
- **Point at `plan.md`, do not restate it.** The plan is on disk exactly so the prompt does
  not have to carry it — restating it in the prompt reintroduces the cost the mode exists to
  avoid, and creates a second version to drift.
- **Say `preflight` in the prompt.** It is the rule a session skips most readily when it
  believes it is merely continuing, and it is the one that has cost the most when skipped.
- **Name leftover state**: a worktree left in place unlanded, a database container left
  running, a phase whose T2 never ran, a decision waiting on Giulio. A next session that has
  to discover those pays for them twice.
- **If the plan is complete**, either the command opens the next piece of work and says so,
  or — when there is no next piece — print no command at all and say the plan is closed.
  A handoff command pointing at nothing is worse than none.

The prompt addresses the next *architect*, so it can be in whatever language the
conversation is in; the plan and specs it points at stay English, per the rule at the top.

## Keeping this file honest

**The layer boundaries here are self-assessments by the layer above.** They record what the
architect kept, not what the architect had to keep, so they are predictions:

> *Phases whose correctness needs a human read have so far needed architect review. Delegate
> one, and if it lands clean the boundary moves.* [untested]

**The boundary moved once already, and in the direction worth watching.** [inferred · 1]
Mutation-based verification — revert a rule, check exactly the intended test goes red — was
written here as the architect's move. In M11's `ws-watch` the *supervisor* ran it unprompted and
found what the gate could not: the pump's mutex could be deleted with every test still green
under `-race`. That is the run's most valuable single finding, and it was produced a layer below
where this file assumed it had to live. **Verification is not the same faculty as judgment**, and
it is the one that delegates: it has a mechanical success criterion (did exactly the intended
test go red?), which is precisely the property that makes anything delegable in this mode. The
skeptical diff read, which has no such criterion, is the part that has not moved and probably
should not. Brief supervisors to mutate the concurrency and authorization rules by name, not just
the rejection tables. *Falsified by: supervisor mutation passes that report green on rules an
architect read then finds unbound.*

> *A supervisor can write the spec for a machine-checkable phase from a one-line intent in
> `plan.md`.* [untested] **Experiment:** next run, take the most machine-checkable phase of
> the batch, have the supervisor spec it, record `spec by = supervisor`. It succeeds only if
> the diff passes the gates **and** survives an architect read — a plausible-but-wrong diff
> that passes is the exact failure nobody catches.

At the end of every run the architect makes one short pass: promote what the run confirmed,
demote or delete what it contradicted, add nothing untagged, move narrative into `evidence/`.
The only thing stopping this file from accreting war stories forever.
