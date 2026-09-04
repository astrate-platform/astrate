# Recipe — dependency, CI and lint upkeep

The unglamorous maintenance nobody schedules. Run this occasionally; expect it to be short.

## Checks, in order of value

```sh
# Directly-required deps with newer versions. Do NOT pipe this through `head` — the 2026-09-02
# run did, spent its whole budget on transitive cloud/azure modules, and reported "only
# transitive skew" while seven direct deps had updates. Ask go.mod which ones are ours.
go list -m -u -f '{{if and .Update (not .Indirect)}}{{.Path}} {{.Version}} -> {{.Update.Version}}{{end}}' all 2>/dev/null | rg .
govulncheck ./... 2>/dev/null || echo "govulncheck not installed"
golangci-lint run ./... 2>&1 | tail -30   # on PATH: /root/go/bin on the Pi, ~/go/bin on the Mac
rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ | head -30
go test ./... 2>&1 | rg -i 'skip|no test files' | head -20
```

## Standing rule on dependency bumps (Giulio, 2026-09-04)

**Do not propose a bump just because a newer version exists.** The full sweep of 2026-09-04
found seven direct deps with updates and none of them fixed anything this repo has; each
would have cost a full test run to land. Re-run this sweep at every **milestone boundary**
(the point where `APICompatVersion` or a milestone tag moves) and propose a bump only when
it carries a fix Astrate actually needs — name the fix and the code path that hits it.
`go.mod` is otherwise on the never-touch list.

## What to propose

- **A vulnerability in a dependency `govulncheck` says is actually reachable.** Highest
  priority thing in this file. Propose it first, name the CVE and the call path.
  Reachability matters: an unreachable advisory is not urgent.
- **A lint finding the config does not already exclude.** Group them: one task per package,
  not one per finding.
- **A `TODO` that names a real missing behaviour** — turn it into a task line quoting the
  comment and its file:line. Ignore TODOs that are notes-to-self or aspirational.
- **A package with no test file at all**, if it contains logic rather than types. Propose a
  test task naming the two most important rules in it.
- **A skipped test.** A skip is a test that stopped being a test; find out why it skips and
  propose either fixing it or deleting it, but say which and why.

## Dependency bumps: careful

`go.mod` and `go.sum` are on the never-touch list — build pins are a decision. So a bump is
**never** a mule task. Write it to `.mule/for-giulio.md` as one line per module: current
version, available version, whether the changelog mentions a breaking change, and whether
anything in this repo uses the affected API. That research is genuinely useful and costs one
run; the decision is not yours.

Same for `.github/` workflow changes — propose in `.mule/for-giulio.md`, never edit.

## Rules

- Do not propose a bump "to stay current". Propose one because of a fix this repo needs.
- Do not propose adding a linter. The `.golangci.yml` set is tuned.
- Nothing to report is the normal outcome. Say so and stop.
