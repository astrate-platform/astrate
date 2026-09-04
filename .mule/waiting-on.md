# Waiting on an upstream release

Decisions deliberately parked until a **specific upstream release** exists. This file is the
missing half of parking something: `.mule/for-giulio.md` records *what* was decided, this
records *what would un-park it* and makes a machine responsible for noticing.

It exists because the un-parking failed once already. Issue #51 was closed on 2026-08-22 with
"parked until the upstream 1.4 experimental spec stabilizes — reopen or file fresh when it
does". The spec was published on 2026-08-31 and nobody noticed until a hand triage on
2026-09-04. The weekly upstream watch had seen the new tags the whole time; nothing connected
them to the parked decision.

**Who reads this:** `.mule/recipes/astarte-upstream.md`, on every run of the weekly
`mule-upstream-watch` job. It compares each `waiting for` tag below against the tags upstream
actually has, and when one exists it writes the row to `.mule/for-giulio.md` — loudly, as its
own escalation, not buried in a list of proposals.

**One row per parked decision.** Delete the row once the decision is actually made — like
`for-giulio.md`, this is a queue, not a log.

| issue | waiting for | why it was parked | what to do when it lands |
|---|---|---|---|
| [#92](https://github.com/astrate-platform/astrate/issues/92) | `v1.4.0` (a **stable** tag — an `-rc.N` does not count) | The key-agreement wire spec (`082-key_agreement_protocol.md`) is complete and implemented upstream, but ships only in `v1.4.0-rc.5`; `v1.3.3` is the newest stable. Document stabilized, release did not. Giulio's decision, 2026-09-04: wait rather than build against a spec that can still move. | Say so in `.mule/for-giulio.md`, naming #92. The scope call (implement the full handshake — CBOR, X25519/P-256, HKDF, AES-256-GCM, 5-state machine, shared-secret persistence, five error names) is Giulio's, not a queue task. |
