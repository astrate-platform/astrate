slug: recipe-astarte-upstream
verdict: proposed
at:  ff7905b
ran: 2026-09-18T11:39:30Z on linux x86_64


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
→ Read docs/COMPATIBILITY.md
→ Read .mule/todo.md
→ Read .mule/for-giulio.md
→ Read .mule/reports/log/20260917T190703Z-recipe-astarte-upstream-proposed.md (prior run)

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.3.4	2026-09-18T09:42:57Z
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z

→ Release body for v1.3.4 (2026-09-18):
  ### Fixed
  - [astarte_data_updater_plant] Prevent AMQPDataConsumer processes from accumulating
    uncollected binaries over time by forcing more frequent full sweep garbage collections
    on them.
  - [astarte_data_updater_plant] Ensure the RPC server is always available to clients.

$ gh api repos/astarte-platform/astarte/compare/v1.3.3...v1.3.4 --jq '{total_commits, messages: [.commits[].commit.message]}'
- fix(dup): prevent AMQP consumer process to accumulate uncollected binaries (#2194)
- chore!(dup): remove fullsweep_after configuration (#2198)
- fix(dup): ensure RPC server is always available to clients (#2203)
- chore: prepare v1.3.4 (#2208)

Upstream moved: new stable v1.3.4 published 2026-09-18, one maintenance patch over v1.3.3. All four
commits are inside astarte_data_updater_plant and are BEAM/Broker-RPC internals: BEAM fullsweep GC
tuning (and removing the fullsweep_after knob) against binary accumulation in AMQP consumers, and
restart-availability of DUP's Erlang RPC server after a cluster reconnect. Assessed against the
recipe's qualification list:

- Wire-visible behaviour (MQTT topics, BSON payloads, pairing flow, API response shape, trigger
  payload)? No. None of the three patches touches a wire or REST surface.
- New interface-schema field or validation rule? No.
- Conceptual improvement that survives restatement in Go? No. GC-sweep tuning of BEAM processes and
  an Erlang RPC server's availability are Elixir-shaped internals; Astrate's single Go binary has its
  own GC, no AMQP consumer processes, and no service-to-service RPC (the "ensure always available"
  fix has no analog). Nothing to port, not even an idea.

So the only qualified proposal is the recipe-mandated compat-note for the newest tag, which now names
v1.3.4 (it supersedes v1.3.3 as newest stable; v1.4.0 is still rc.5-only). The open v1.3.2 wording
proposal in .mule/for-giulio.md stands and this folds v1.3.3 + v1.3.4 into it, exactly as the blocked
compat-note-v1.3.3 line was meant to. One line appended to .mule/todo.md, nothing else. No gh issue
creation (this recipe's deliverable is the task-line, not an issue).

Done: upstream watch — new stable v1.3.4 (2026-09-18), maintenance-only (DUP GC/RPC internals, no
wire change); proposed compat-note-v1.3.4, nothing else
Files: .mule/reports/log/20260918T113930Z-recipe-astarte-upstream-proposed.md (evidence); .mule/todo.md (1 task line)
Verified: gh api repos/astarte-platform/astarte/releases -> v1.3.4 newest (2026-09-18); compare v1.3.3...v1.3.4 -> 4 DUP-only commits -> pass
Unsure: nothing
Follow-ups: compat-note-v1.3.4