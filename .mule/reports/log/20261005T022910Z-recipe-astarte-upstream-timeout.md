slug: recipe-astarte-upstream
verdict: timeout
at:  ca68353
ran: 2026-10-05T02:29:10Z on DietPi in 1201s


> build · big-pickle

I'll read .mule/MULE.md and then the astarte-upstream recipe and follow the instructions precisely.
→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"' 2>&1 | head -5
v1.4.0-rc.6	2026-09-30T07:28:27Z
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z

