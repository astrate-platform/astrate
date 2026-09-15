slug: recipe-astarte-upstream
verdict: proposed
at:  b36f50a
ran: 2026-09-15T11:13:54Z on DietPi in 97s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md
→ Read .mule/todo.md
Recipe read. It's a proposal job — output goes to `.mule/todo.md` + any evidence files. No issue creation in this recipe, no git. Starting: check upstream releases vs compatibility target.

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

Newest tag: v1.4.0-rc.5 (2026-08-20); newest stable: v1.3.3 (2026-08-07). The prior run already saw both (todo.md's `compat-note-v1.3.3` line says "v1.4.0 is still rc.5-only"). Let me confirm the compat-note state in for-giulio.md before concluding.
→ Read .mule/for-giulio.md
Upstream has not moved. Newest tag is still **v1.4.0-rc.5** (2026-08-20); newest stable is still **v1.3.3** (2026-08-07). Both are exactly what the previous runs already saw: the `compat-note-v1.3.3` line and `.mule/for-giulio.md` (milestone runs 2026-09-09→09-13) record the same pair, and the v1.3.2/v1.3.3/v1.4-rc compat notes and the v1.3.0 probes are already proposed in the queue.

Per the recipe, an empty result is a good result: I append nothing to `.mule/todo.md` and write no evidence file.

```
Done: no upstream movement since v1.4.0-rc.5 — newest tag unchanged (v1.4.0-rc.5, 2026-08-20; stable still v1.3.3, 2026-08-07); nothing new to propose
Files: none
Verified: gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"' -> pass (5 tags, newest v1.4.0-rc.5); compared against docs/COMPATIBILITY.md target (1.2.2, v1.3 surface delivered, v1.4 experimental) and prior queue state
Unsure: nothing — release set matches the last recorded run exactly
Follow-ups: none
```
