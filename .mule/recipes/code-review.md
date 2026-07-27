# Recipe — codebase review

Read one area of the codebase and propose work. **This run writes no code.**

## Pick one area, not the codebase

Your context does not stretch to a whole-repo review, and a review that skims everything
finds nothing. Pick **one** package and go deep. Rotate through them across runs — check
`.mule/log.md` and `.mule/reviews/` for which areas were done recently, and pick one that
was not.

```sh
ls internal/ pkg/                                   # the map
rg -c '' --glob 'internal/<pkg>/**/*.go' | sort -t: -k2 -rn | head   # the size of it
```

Then read that package's files, largest first, and its tests.

## What is worth proposing

In descending order of value:

1. **A behaviour that is wrong or unguarded** — an unchecked error, a nil path, an
   unbounded slice grown from network input, a goroutine with no exit, a lock held across
   I/O, a context ignored. Propose with the file and line.
2. **A missing test on a rule that already exists.** Rules nothing asserts on rot silently.
   Best kind of task there is: fully machine-checkable.
3. **Performance, but only where you can name the workload.** "This allocates per message
   in the ingest hot path" is a proposal. "This could be faster" is not. Propose a
   `go test -bench` measurement task *before* an optimisation task, never the reverse.
4. **Clarity** — a function that does three things, a name that lies, duplicated logic in
   two packages. Propose only when the change is mechanical and the tests already cover the
   behaviour.
5. **New features.** Rarely. A feature is a design decision and it is not yours: write it to
   `.mule/for-giulio.md` as one line saying what and why, not to the queue.

## What not to propose

- Renames and reformatting for their own sake. The linter owns that.
- Anything in `docs/DESIGN.md` or `docs/ROADMAP.md` — that is the frozen spec, and if the
  code contradicts it, **the code is the bug**; say so in `.mule/for-giulio.md`.
- Refactors that touch more than two files. They cannot be verified in one task.
- Anything you found by pattern-matching rather than by reading. Say "I read
  `internal/x/y.go` lines 40-120" or do not propose it.

## Output

Write `.mule/reviews/<pkg>-<date>.md`: what you read, what you found, and what you decided
*not* to propose and why (that half is what stops the next run repeating your work). Then
append at most five task lines to `.mule/todo.md`, each naming its file.

Finding little in a well-written package is a correct result. Say so and stop.
