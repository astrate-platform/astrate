slug: recipe-astarte-upstream
verdict: proposed
at:  4ed448a
ran: 2026-09-03T11:06:58Z on DietPi in 105s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/astarte-upstream.md
Now I'll read `docs/COMPATIBILITY.md` and check the current todo queue, then run the upstream check.
→ Read docs/COMPATIBILITY.md
→ Read .mule/todo.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[0:5][] | "\(.tag_name)\t\(.published_at)"'
v1.4.0-rc.5	2026-08-20T12:28:04Z
v1.3.3	2026-08-07T14:33:58Z
v1.4.0-rc.4	2026-08-06T13:40:20Z
v1.4.0-rc.3	2026-07-31T10:48:39Z
v1.4.0-rc.2	2026-07-14T10:16:22Z

✱ Grep "APICompatVersion" 6 matches
I need to understand what "upstream version" the project targets. Let me check how the API version is set and see what's already been proposed for each release.

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.5") | .body' 2>&1 | head -200
Sixth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.4

### Fixed

- [astarte_data_updater_plant] Make object values only validate type on selected interface
- [astarte_data_updater_plant] Ensure binaryblob data is correctly validated

→ Read .mule/for-giulio.md

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.3.3") | .body' 2>&1 | head -200

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.3") | .body' 2>&1 | head -200
Fourth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.2

### Added

- [secrets] Configurable base namespace for Vault using `ASTARTE_VAULT_BASE_NAMESPACE`: all created namespaces will be children of the base namespace.

### Fixed

- [fdo] Support for Vault and Vault managed
- [fdo] Properly handle chunked service info messages
- [astarte_data_updater_plant] Avoid clogging AMQP queues when receiving invalid device messages
- [astarte_data_updater_plant] Restore correct type validation checks for incoming data

Let me check the release notes for the two versions I haven't reviewed yet.
$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.4") | .body' 2>&1
Fifth Astarte v1.4.0 release candidate.

## Changes since v1.4.0-rc.3

### Added

- [config] Initialize `config` library. All external services declared using the `config` library helpers can be configured either using the `_URL` environment variable or by using single component variables `_SCHEME`, `_HOST`, `_PORT`, `_PATH`, `_QUERY`, `_FRAGMENT` (`_URL` has priority), together with normal `_SSL_*` variables.
- [secrets] Use `config` for vault configuration. Vault can now be configured using component variables instead of `ASTARTE_VAULT_URL`:
  - `ASTARTE_VAULT_SCHEME`: scheme part of the URL (binary, default: `"http"`, or `"https"` if `ASTARTE_VAULT_SSL_ENABLED` is `true`)
  - `ASTARTE_VAULT_HOST`: host part of the URL (binary, default: `"localhost"`)
  - `ASTARTE_VAULT_PORT`: port part of the url (integer, default `8200`)
  - `ASTARTE_VAULT_PATH`: path part of the url (string)
  - `ASTARTE_VAULT_QUERY`: query part of the url (string)
  - `ASTARTE_VAULT_FRAGMENT`: fragment part of the url (string)
- [fdo] Use `config` for rendezvous configuration. The rendezvous can now be configured using component variables instead of `PAIRING_FDO_RENDEZVOUS_URL`:
  - `PAIRING_FDO_RENDEZVOUS_SCHEME`: scheme part of the URL (binary, default: `"http"`, or `"https"` if `PAIRING_FDO_RENDEZVOUS_SSL_ENABLED` is `true`)
  - `PAIRING_FDO_RENDEZVOUS_HOST`: host part of the URL (binary, default: `"localhost"`)
  - `PAIRING_FDO_RENDEZVOUS_PORT`: port part of the url (integer, default `8041`)
  - `PAIRING_FDO_RENDEZVOUS_PATH`: path part of the url (string)
  - `PAIRING_FDO_RENDEZVOUS_QUERY`: query part of the url (string)
  - `PAIRING_FDO_RENDEZVOUS_FRAGMENT`: fragment part of the url (string)

### Fixed

- [astarte_data_updater_plant] Ensure RPC server is always available to clients. Resolved the issue where a temporary disconnection and reconnection of data_updater_plant to the cluster would make the RPC server inaccessible.

### Changed

- [astarte_data_updater_plant] Ensure memory is properly garbage collected

$ gh api repos/astarte-platform/astarte/releases --jq '.[] | select(.tag_name=="v1.4.0-rc.2") | .body' 2>&1
Third Astarte v1.4.0 release candidate.

I now have the full picture. Let me verify the existing coverage:

**COMPATIBILITY.md**: emulated API level is `1.2.2`. Latest upstream is **v1.4.0-rc.5** (2026-08-20).

**Existing proposals in todo.md for upstream releases** (all BLOCKED):
- `probe-value-type-validation` — v1.4.0-rc.3
- `probe-interface-default-values` — v1.4.0-rc.3
- `probe-binaryblob-validation` — v1.4.0-rc.5
- `probe-object-validation-selected-interface` — v1.4.0-rc.5
- `probe-emptycache-resend-device-error` — v1.3.3
- `probe-properties-on-connect-encoding` — v1.3.3
- `probe-mqtt-capabilities-declaration` — v1.3.x SDKs
- `compat-note-v1.3.3` — v1.3.3
- `compat-note-v1.4.0-rc.3` — v1.4.0-rc.3
- `compat-note-v1.4.0-rc.5` — v1.4.0-rc.5

Releases not covered by existing tasks:
- **v1.4.0-rc.4**: config library / component variables (infra config only), RPC server reconnect fix (internal Elixir clustering), GC — none wire-visible.
- **v1.4.0-rc.2**: empty body — nothing to propose.

Every wire-visible change across all five releases already has a proposed task line in the queue. There are no new items to append.

---

**Result**: no upstream movement since the last run that isn't already captured in `.mule/todo.md`. All five releases (v1.4.0-rc.2 through v1.4.0-rc.5, plus v1.3.3) have been surveyed; the wire-visible changes in v1.3.3, v1.4.0-rc.3, and v1.4.0-rc.5 are already BLOCKED probe/compat-note tasks. v1.4.0-rc.2 and v1.4.0-rc.4 contain only infrastructure and internal fixes with no wire impact. No new task lines appended, no evidence file written.
