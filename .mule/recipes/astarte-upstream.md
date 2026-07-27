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

Compare the newest tag to the version named in `docs/COMPATIBILITY.md`. If they match, say
so and stop — that is a complete, correct, cheap result, and it is the expected one most of
the time. **Do not go looking for work when there is none.**

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
