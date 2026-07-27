# Recipe — GitHub issues

Turn open issues on `astrate-platform/astrate` into task lines. **You are triaging, not
fixing.** Nothing here writes code.

## Do this

```sh
gh issue list --state open --limit 40 \
  --json number,title,labels,updatedAt,comments \
  --template '{{range .}}#{{.number}} {{.title}} [{{range .labels}}{{.name}} {{end}}] {{.comments}}c{{"\n"}}{{end}}'
```

That one command is your whole survey. Do **not** `gh issue view` all of them — open only
the handful you are actually going to propose, and read the body plus the last comment, not
the full thread.

For each issue you propose, append one line:

```
- [ ] issue-<number>-<slug>: <the outcome, in one line> (closes #<number>)
```

## Triage rules

- **Propose at most five per run**, newest-and-smallest first. A queue nobody approves is
  worse than no queue.
- **Only propose an issue whose fix is machine-checkable** — a wrong behaviour with a
  reproducible input, a missing validation, a crash. If verifying the fix needs a human to
  read the diff against intent, or needs a live external system, propose it as a
  `probe-` task that investigates and reports instead.
- **Never propose an issue that asks for a design decision** (a new API shape, a schema
  change, a protocol extension). Write it to `.mule/for-giulio.md` with one line on what the
  decision is, and move on. That file is the escalation channel; create it if absent.
- If an issue is stale, contradicted by the current code, or already fixed, say so — append
  a line to `.mule/for-giulio.md` proposing it be closed. **Do not close or comment on
  anything on GitHub yourself.**
- If the issue names a file, put the path in the task line. That saves the executing session
  a grep it would otherwise pay for.

## Reproduction first

For any bug issue you propose, the task line should ask for the **failing test first**:

```
- [ ] issue-<n>-repro: add a failing test in <pkg>_test.go reproducing #<n>, then fix it
```

An issue whose reproduction cannot be written as a test is an issue whose fix cannot be
verified — say so rather than proposing it.
