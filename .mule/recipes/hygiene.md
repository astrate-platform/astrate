# Recipe — dependency, CI and lint upkeep

The unglamorous maintenance nobody schedules. Run this occasionally; expect it to be short.

## Checks, in order of value

```sh
go list -m -u all 2>/dev/null | rg '\[' | head -20        # deps with newer versions
govulncheck ./... 2>/dev/null || echo "govulncheck not installed"
/Users/atsetilam/go/bin/golangci-lint run ./... 2>&1 | tail -30
rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ | head -30
go test ./... 2>&1 | rg -i 'skip|no test files' | head -20
```

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
