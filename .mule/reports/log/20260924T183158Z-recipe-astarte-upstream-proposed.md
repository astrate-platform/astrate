slug: recipe-astarte-upstream
verdict: proposed
at:  375e038
ran: 2026-09-24T18:31:58Z on DietPi in 105s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z

→ Read .mule/todo.md
→ Read .mule/for-giulio.md
Upstream's newest tag is **v1.3.4** (2026-09-18). COMPATIBILITY.md targets upstream 1.2.2 (emulated API compat) with 1.4 tracked as experimental. Let me check the current queue state to avoid duplicates and understand the tracking of 1.4.
→ Read docs/UPSTREAM-EXPERIMENTAL.md
Newest tag is **v1.3.4** (2026-09-18) — unchanged since the prior astarte-upstream/milestone runs documented in `.mule/for-giulio.md` (through 2026-09-23, which already assessed v1.3.4 as maintenance-only and folded it into the open compat-note proposal). The recipe demands a release-note check rather than trusting that claim, so let me verify the v1.3.4 body and confirm nothing newer exists.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.4") | .body'
### Fixed

- [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
  uncollected binaries over time by forcing more frequent full sweep garbage collections
  on them.
- [astarte_data_updater_plant] Ensure the RPC server is always available to clients.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | .tag_name'
v1.3.4
v1.4.0-rc.5
v1.3.3
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
v1.3.2
v1.3.1
v1.3.0
v1.2.2
v1.4.0-rc.0
v1.3.0-rc.2
v1.2.2-rc.0
v1.2.1
v1.2.1-rc.1
v1.3.0-rc.1
v1.3.0-rc.0
v1.2.1-rc.0
v1.2.1-alpha.0
v1.2.0
v1.2.0-rc.0
v1.0.6
v1.1.1
v1.0.5
v1.1.0
v1.1.0-rc.0
v1.1.0-alpha.0
v1.0.4
v1.0.3

No upstream movement since **v1.3.4**. Confirmed: newest release is still `v1.3.4` (2026-09-18), whose body is maintenance-only (`astarte_data_updater_plant` GC-sweep + RPC-availability — internal Elixir fixes, nothing on the wire). Newest rc remains `v1.4.0-rc.5` (2026-08-20), also already assessed. Nothing shipped since the prior runs; the v1.3.4 compat-note wording is already folded into the open proposal in `.mule/for-giulio.md`, so re-proposing it would duplicate that. Per the recipe's empty-result rule I append nothing to `.mule/todo.md`.

```
Done: no upstream movement since v1.3.4; appended nothing to the queue
Files: none
Verified: gh api repos/astarte-platform/astarte/releases --jq '.[0:5][]...' -> pass; v1.3.4 body + full tag list -> pass
Unsure: nothing
Follow-ups: none
```
