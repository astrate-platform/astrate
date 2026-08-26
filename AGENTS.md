# Agent context notes

Facts every model working in this repo must know. Read before acting.

- **Issue #1 is a joke and must never be touched.** Do not comment on it, do
  not close it, do not reference it in work, do not "fix" it. It stays open
  forever, by Giulio's explicit decision.

- **Never run tests that need containers, builds that need containers, or
  Docker itself on this Mac.** Docker on macOS makes the machine unusably
  slow. All testing that needs Docker happens on the Lenovo Legion Go — load
  the `/machines` skill to reach it. Container-free unit tests (`go test`,
  httptest-based rigs) are fine on this Mac and are the normal delegated
  gate; integration-tagged suites (`-tags integration`) and anything starting
  a testcontainer are not. This rule applies to every agent and must be
  carried into every handoff/phase file you write — quote it precisely,
  including the container-free exception, or a cautious delegate will refuse
  to run its own acceptance tests (seen 2026-08-24).

- **The Legion Go is OFF unless Giulio says otherwise (since 2026-08-26).**
  Before any work that needs it — containers, integration suites, live
  upstream — probe ONCE and cheaply: `ssh -o ConnectTimeout=5 legion 'echo ok'`.
  If it does not answer, stop immediately and tell Giulio he needs to power it
  on. Do NOT retry, do NOT scan the network, do NOT degrade into workarounds:
  a session that probes in circles burns tokens and finds nothing, because the
  machine is simply off. Say what you need it for so he knows it is worth
  switching on.
