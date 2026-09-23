slug: recipe-astarte-upstream
verdict: proposed
at:  4aa1525
ran: 2026-09-23T10:52:00Z on DietPi in 156s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md
The recipe names `docs/COMPATIBILITY.md` which I've read — it references upstream 1.2.x and 1.4 experimental. Now let me check upstream releases and the current todo queue.
→ Read .mule/todo.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z

Upstream is at v1.3.4 (2026-09-18), one step past the last-handled v1.3.3. Let me read the release notes and check what's already proposed in for-giulio.md to avoid duplicating.
→ Read .mule/for-giulio.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.4") | "TAG: \(.tag_name)\nDATE: \(.published_at)\nNAME: \(.name)\nBODY:\n\(.body)"'
TAG: v1.3.4
DATE: 2026-09-18T09:42:57Z
NAME: v1.3.4
BODY:
### Fixed

- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
  uncollected binaries over time by forcing more frequent full sweep garbage collections
  on them.
- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.

$ ls .mule/ .mule/recipes/ 2>/dev/null
.mule/:
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
recipes
reports
reviews
tasks
todo.md

.mule/recipes/:
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

Upstream's newest stable is now **v1.3.4** (2026-09-18) — past the last handled v1.3.3. Its release body has two fixes, both in `astarte_data_updater_plant`. Let me confirm they're both Elixir-shaped/internal and check whether the v1.3.4 compat-note is already covered.
✱ Grep "compat-note|v1\.3\.4" in .mule/todo.md · 100 matches
Upstream has not moved since the last successful run of this recipe.

- Newest tag is still **v1.3.4** (2026-09-18), confirmed just now via the release sweep; v1.4.0 is still rc.5-only and its deltas are already tracked (#137/#138).
- The last astarte-upstream run (2026-09-21, `20260921T105709Z-recipe-astarte-upstream-proposed.md`) reached the same conclusion; the run before it (2026-09-18, `bcd4383`) already proposed and landed the `compat-note-v1.3.4` task (today `[x]` at todo.md:177), and its wording sits in `.mule/for-giulio.md` folding v1.3.3+v1.3.4 into the v1.3.2 proposal.
- Reading v1.3.4's body again: both entries are inside `astarte_data_updater_plant` (fullsweep-GC tuning on AMQPDataConsumer processes; RPC-server availability) — Elixir-shaped internals on a service whose AMQP surface Astrate deliberately does not implement (deviation 13). No wire/API/validation surface change, so nothing qualifies for a proposal line.

A fresh `compat-note-v1.3.4` line would duplicate the already-proposed wording, so I append nothing — exactly the "an empty result is a good result" branch.

```
Done: no upstream movement since v1.3.4 — nothing proposed, no task lines appended
Files: none (did not touch todo.md, for-giulio.md, or git)
Verified: gh api releases -> newest stable v1.3.4 (2026-09-18), no release since prior run bcd4383 -> pass
Unsure: nothing
Follow-ups: none
```
