slug: recipe-astarte-upstream
verdict: proposed
at:  d37434e
ran: 2026-10-09T19:25:30Z on DietPi in 79s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md
→ Read .mule/todo.md
→ Read .mule/for-giulio.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:8][] | "\(.tag_name)\t\(.published_at)\tprerelease=\(.prerelease)"'
v1.3.5	2026-10-05T14:11:44Z	prerelease=false
v1.4.0-rc.6	2026-09-30T07:28:27Z	prerelease=true
v1.3.4	2026-09-18T09:42:57Z	prerelease=false
v1.4.0-rc.5	2026-08-20T12:28:04Z	prerelease=true
v1.3.3	2026-08-07T14:33:58Z	prerelease=false
v1.4.0-rc.4	2026-08-06T13:40:20Z	prerelease=true
v1.4.0-rc.3	2026-07-31T10:48:39Z	prerelease=true
v1.4.0-rc.2	2026-07-14T10:16:22Z	prerelease=true

$ gh api repos/astarte-platform/astarte/tags --jq '.[0:6][].name'; echo '---newest commits since v1.3.5---'; gh api 'repos/astarte-platform/astarte/commits?since=2026-10-05T14:11:44Z&per_page=20' --jq '.[] | "\(.sha[0:8]) \(.commit.author.date) \(.commit.message | split("\n")[0])"'
v1.4.0-rc.6
v1.4.0-rc.5
v1.4.0-rc.4
v1.4.0-rc.3
v1.4.0-rc.2
v1.4.0-rc.1
---newest commits since v1.3.5---
171c1bd7 2026-10-06T15:57:45Z chore: forward-port astarte_vmq_plugin from release-1.4
1da04832 2026-10-05T14:57:19Z chore: forward port release-1.4
a8a26d0f 2026-10-05T14:40:35Z chore: forward port release-1.3
d85e0ed0 2026-10-06T10:02:32Z docs: sync FDO workflow to the current behavior (#2241)
a6bd2c21 2026-10-06T09:59:16Z fix(fdo): reject owner key names not accepted by OpenBao (#2252)
a51ab2e1 2026-10-06T07:21:18Z refactor: Update Umbrella checks (#2244)
7628d2ed 2026-10-05T15:49:18Z test(generators): validate interfaces with Core changesets
ed1e99f4 2026-10-05T15:49:04Z test(generators): validate mappings with Core changesets
bb17455a 2026-10-05T15:48:50Z test(generators): validate policy handlers with Core changesets
ec020bf8 2026-10-05T15:48:31Z fix(generators): generate valid trigger policy names

$ gh api repos/astarte-platform/astarte/commits/171c1bd7 --jq '.commit.message, "---files---", (.files[] | "\(.status) \(.filename) (+\(.additions)/-\(.deletions))")'
chore: forward-port astarte_vmq_plugin from release-1.4

Signed-off-by: frossq <francesco.asquini@secomind.com>
---files---
modified CHANGELOG.md (+85/-3)
removed apps/astarte_vmq_plugin/CHANGELOG.md (+0/-255)
modified apps/astarte_vmq_plugin/Dockerfile (+3/-1)
modified apps/astarte_vmq_plugin/lib/astarte_vmq_plugin/application.ex (+2/-1)
modified apps/astarte_vmq_plugin/lib/astarte_vmq_plugin/config.ex (+11/-0)
added apps/astarte_vmq_plugin/lib/astarte_vmq_plugin/health_http/handler.ex (+37/-0)
added apps/astarte_vmq_plugin/lib/astarte_vmq_plugin/health_http/route_injection_service.ex (+110/-0)
modified apps/astarte_vmq_plugin/mix.exs (+5/-1)
modified apps/astarte_vmq_plugin/priv/astarte_vmq_plugin.schema (+4/-0)
modified apps/astarte_vmq_plugin/rel/overlays/bin/start_vernemq (+3/-0)
added apps/astarte_vmq_plugin/test/astarte_vmq_plugin_health_http_handler_test.exs (+59/-0)
added apps/astarte_vmq_plugin/test/astarte_vmq_plugin_route_injection_service_test.exs (+86/-0)
modified apps/astarte_vmq_plugin/test/test_helper.exs (+1/-0)
modified mix.lock (+3/-1)

No new release since the last `astarte-upstream` run: newest stable is still **v1.3.5** (2026-10-05), newest overall still the prerelease **v1.4.0-rc.6** (2026-09-30). The only master commits since v1.3.5 are forward-ports, generator tests, FDO docs/fix and CI refactors; the one MQTT-adjacent commit (`171c1bd7`, "forward-port astarte_vmq_plugin") adds a **VerneMQ plugin health HTTP endpoint** — Elixir/VerneMQ infrastructure Astrate does not implement (embedded mochi broker), so wire-inert here. Per the recipe, nothing is appended.

No `gh issue create` (recipe doesn't ask for one); no evidence file named; nothing written to `.mule/todo.md` or git.

```
Done: astarte-upstream probe — no upstream movement since v1.3.5 (newest stable 2026-10-05; rc.6 still newest overall); nothing proposable, queue untouched
Files: none
Verified: gh api .../releases -> newest stable v1.3.5, no stable v1.4.0; gh api .../commits since v1.3.5 -> forward-ports/generators/FDO/CI only, all wire-inert
Unsure: nothing
Follow-ups: none
```
