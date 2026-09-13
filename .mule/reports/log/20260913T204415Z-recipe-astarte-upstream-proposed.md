slug: recipe-astarte-upstream
verdict: proposed
at:  661fa76
ran: 2026-09-13T20:44:15Z on DietPi in 119s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

→ Read docs/UPSTREAM-EXPERIMENTAL.md
Newest upstream release is `v1.4.0-rc.5` (2026-08-20). Let me check what the compat doc names as the target and any prior tracking.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | "\(.tag_name)\t\(.published_at)"' | head -40
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z
v1.4.0-rc.1	2026-07-14T08:53:27Z
v1.3.2	2026-07-14T10:07:26Z
v1.3.1	2026-07-13T11:50:05Z
v1.3.0	2026-05-06T15:30:50Z
v1.2.2	2026-04-27T16:39:28Z
v1.4.0-rc.0	2026-04-08T12:32:33Z
v1.3.0-rc.2	2026-04-08T16:34:59Z
v1.2.2-rc.0	2026-04-08T16:31:03Z
v1.2.1	2026-03-12T16:11:20Z
v1.2.1-rc.1	2026-02-13T10:36:20Z
v1.3.0-rc.1	2026-01-26T15:58:33Z
v1.3.0-rc.0	2025-11-21T14:58:46Z
v1.2.1-rc.0	2025-08-27T07:03:28Z
v1.2.1-alpha.0	2025-04-10T15:31:04Z
v1.2.0	2024-07-02T14:09:07Z
v1.2.0-rc.0	2024-06-11T13:02:33Z
v1.0.6	2024-04-23T21:15:57Z
v1.1.1	2023-11-15T15:29:57Z
v1.0.5	2023-09-26T10:27:57Z
v1.1.0	2023-06-20T18:59:09Z
v1.1.0-rc.0	2023-06-09T15:49:44Z
v1.1.0-alpha.0	2022-11-24T16:38:37Z
v1.0.4	2022-10-25T16:50:21Z
v1.0.3	2022-07-05T13:57:17Z
v1.0.2	2022-04-01T17:18:15Z

→ Read .mule/todo.md
→ Read .mule/for-giulio.md
**No upstream movement since v1.4.0-rc.5** (newest tag, published 2026-08-20; newest stable still v1.3.3, 2026-08-07).

The recipe's first gate — compare newest tag to the compatibility target — resolves to "already triaged, nothing new":

- `COMPATIBILITY.md` (never-touch) already documents the 1.4 rc surface: deviation 17 (async_operation), deviation 18 (health endpoints), with `UPSTREAM-EXPERIMENTAL.md` rows #67/#68 (required/encrypted flags, async_operation) waiting on **v1.4 final** for the reconcile-vs-promote call — still rc-only.
- The release-sweep fact was re-verified today by the milestone recipe run (`for-giulio.md:13`): no stable v1.4.0, nothing shipped since rc.5.
- The rc-line content (encrypted endpoints, `required` flag, FDO, Vault/config) is already covered by landed queue items (`compat-note-v1.3.2` done; `compat-note-v1.3.3`/`compat-note-v1.4-rc` already queued, blocked only because the wording is a Giulio-doc decision in `for-giulio.md`), so re-proposing any compat-note line would duplicate queued work.

Appending **nothing** to `.mule/todo.md` — empty result is the correct result. No git touched.

Done: upstream unchanged since v1.4.0-rc.5 (2026-08-20) / v1.3.3 stable; nothing to propose
Files: none
Verified: `gh api repos/astarte-platform/astarte/releases` (full list, newest = v1.4.0-rc.5 2026-08-20) vs docs/COMPATIBILITY.md + UPSTREAM-EXPERIMENTAL.md -> matched, no gap
Unsure: nothing
Follow-ups: none (a compat-note-v1.4-rc line already sits blocked in the queue, tied to the pending wordings in for-giulio.md)
