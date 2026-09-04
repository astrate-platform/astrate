# Recipe — watch upstream Astarte

Astrate is wire-compatible with Astarte. Upstream is the organisation
`github.com/astarte-platform`, whose main repos are `astarte` (the umbrella, and where the
release tags live), `astarte_e2e`, `astarte-device-sdk-*` and `astarte-dashboard`.

Compatibility state is recorded in `docs/COMPATIBILITY.md` — **read that first**, it tells
you which upstream version Astrate currently targets. It is also on the never-touch list:
propose changes to it, do not make them.

## Do this

```sh
gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
```

**First, before anything else, check the parked decisions.** Read `.mule/waiting-on.md`. For
every row, ask whether the tag it is waiting for now exists upstream (the command above lists
recent releases; use `gh api repos/astarte-platform/astarte/tags` if you need to look further
back). Read each row's "waiting for" cell literally — a row waiting on a **stable** tag is not
satisfied by an `-rc.N` of the same version.

If a row's tag has landed, that is the single most valuable thing this run can report. Write
it to `.mule/for-giulio.md` as its own escalation naming the issue, the tag, and the row's
"what to do when it lands" text — **not** as a `- [ ]` line in `.mule/todo.md`, because a
parked decision is Giulio's to make and is never a queue task. Do this even if you find
nothing else all run, and do it before spending budget on the rest of the recipe.

This step exists because it failed once: issue #51 was parked "until the upstream 1.4
experimental spec stabilizes", the spec was published 2026-08-31, and this job saw the new
tags weekly without ever connecting them to the parked issue. It was caught by hand four days
later.

Then compare the newest tag to the version named in `docs/COMPATIBILITY.md`. If they match, say
so and stop — that is a complete, correct, cheap result, and it is the expected one most of
the time. **Do not go looking for work when there is none.** (Reporting an un-parked row above
is a complete result too, even when the version reference itself is current.)

If upstream is ahead:

```sh
gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="<newtag>") | .body'
```

Read the release notes. Then, and only for the entries that plausibly touch something
Astrate implements, look at the actual upstream change:

```sh
gh search code --owner astarte-platform '<symbol or key>'   # find where it lives
gh api repos/astarte-platform/<repo>/commits?path=<path>    # what changed
```

## What to propose

One line per change, and be strict about which changes qualify:

- **Wire-visible behaviour** — MQTT topics, payload encoding (BSON), the pairing flow, an
  API response shape, a trigger payload. These are compatibility obligations. Propose them.
- **A new interface-schema field or validation rule.** Propose it.
- **Conceptual improvements** — a smarter reconnection policy, a better session-lookup
  structure, a limit that upstream learned it needed. Propose these as *investigation* tasks:
  `- [ ] probe-<slug>: does Astrate have the problem upstream's <change> fixes? report, do not patch`
- **Anything Elixir-shaped** (their supervision trees, their release tooling) — ignore. Do
  not port an implementation, port an idea, and only when the idea survives being restated
  in Go.

Always propose, as the last line, updating the version reference:

```
- [ ] compat-note-<tag>: propose the docs/COMPATIBILITY.md wording for <tag> in .mule/for-giulio.md (do not edit the file)
```

## Rules

- **A release note is a claim, not a fact.** If the change matters, look at the diff.
- Propose at most five items per run.
- Never copy upstream code into this repo. Different licence, different language, different
  architecture. Read it, understand the rule, implement the rule.
- If nothing changed since last time, append nothing to the queue and say "no upstream
  movement since \<tag\>". An empty result is a good result.
