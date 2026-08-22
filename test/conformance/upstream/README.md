# Recorded upstream observations

Fixtures captured from a **real upstream Astarte**, so that claims Astrate makes about parity
stop being reconstructions. Everything here was observed; nothing was inferred.

The distinction this directory exists to enforce: a test that hits the network to decide whether
Astrate is correct fails whenever the upstream is down, and a test that compares Astrate against a
recorded transcript fails only when Astrate changes. **Re-recording is therefore a deliberate,
separate act** — the `record` command — and comparison is an ordinary offline test.

| file | what it holds |
| --- | --- |
| `rest-errors.json` | the fixture the comparison tests read |
| `rest-errors.transcript.txt` | the exchange that produced every entry, verbatim |
| `record/` | the recorder; talks HTTP only and imports nothing from Astrate |

Each observation carries a `why`, because a fixture nobody can motivate is a fixture nobody dares
change.

## Re-recording

Needs a reachable upstream and its housekeeping key. `bench/scripts/up-astarte.sh` brings one up;
if it runs on another machine, tunnel to it (see the script's header — `.localhost` names cannot be
redirected with `/etc/hosts` on macOS).

```sh
cd test/conformance/upstream
ASTARTE_UPSTREAM_URL=http://api.astarte.localhost:8080 \
ASTARTE_UPSTREAM_REALM=bench \
ASTARTE_UPSTREAM_HOUSEKEEPING_KEY=…/compose/astarte-keys/housekeeping_private.pem \
    go run ./record
```

The realm must already exist — `go run . provision` in `bench/` creates it. **No credential is
written to these files**: tokens are described in the transcript ("a token signed by a key upstream
has never seen"), never printed, and a realm public key in a response body is elided.

## What the current recording says

Recorded 2026-07-26 against Astarte **v1.2.0**, realm `bench`.

- **Unmatched routes do not answer uniformly, and the split is not where Astrate assumed.**
  AppEngine and RealmManagement answer `{"errors":{"detail":"Not found"}}`; Pairing answers
  `{"errors":{"detail":"Page not found"}}` — and so does **Housekeeping**, which Astrate
  deliberately leaves as Go's plain-text 404 on the grounds that it had never been observed. That
  is a real deviation, found by this recording, and it is phase 05's to settle.
- **Deviation 5 in `docs/COMPATIBILITY.md` is now verified rather than argued.** Against the same
  authenticated route: no token → **401**, a malformed token → **403**, a well-formed token signed
  by an unknown key → **403**. Astrate answers 401 to all three. The valid-token row is recorded
  alongside them on purpose: without an acceptance, a deployment that refused everything would
  satisfy every rejection row.
- **An existing realm is reported as `422` with `error_name: existing_realm`**, where Astrate
  reports `409`. Found when `bench provision` could not be re-run against upstream.

### One caveat, recorded rather than hidden

The `unknown service prefix` row (`/no-such-service/v1/thing` → plain-text `404 page not found`)
is **answered by the gateway, not by Astarte**: traefik routes on `Host` plus `PathPrefix`, and a
prefix belonging to no service matches no router at all, so traefik's own default 404 is what the
client sees. The row is kept because it bounds the envelope — it shows the JSON envelopes above
come from the services and not from something in front of them — but **nothing about Astarte's
behaviour may be concluded from it**, and a deployment that puts Astrate behind a different proxy
would legitimately differ.

## Channels: frames and `device_error` names (recorded 2026-07-26)

`channels.json` and `channels.transcript.txt`, produced by `recordchannels/`. This is the half of
the oracle that needed a device driven over MQTT, which is why it was split from the REST half.

```sh
cd test/conformance
ASTARTE_UPSTREAM_URL=http://api.astarte.localhost:8080 \
ASTARTE_UPSTREAM_STATE=../../bench/upstream.json \
    go run ./upstream/recordchannels
```

`bench/upstream.json` is what `go run . provision` writes in `bench/`; it holds the realm signing
key and the device credentials secrets, is gitignored, and **nothing from it reaches the fixture** —
checked before committing.

### The prize: M11's invented mapping table is 6 of 8 right, and 2 wrong

M11 phase 08 chose `internal/engine/triggers/errorname.go`'s table by reading an enum out of the
Dashboard's JavaScript bundle. Every value it produces is *accepted* by the Dashboard; until now
nobody had checked whether upstream says the same thing for the same rejection. Upstream does, for
six of the eight provocations that can be reached with a device-owned datastream:

| provocation | Astrate's reason → its mapped name | upstream says | |
| --- | --- | --- | --- |
| interface not in introspection | `interface_not_in_introspection` → `invalid_interface` | `interface_loading_failed` | **wrong** |
| unknown path on a known interface | `unexpected_path` → `mapping_not_found` | `mapping_not_found` | ok |
| payload that is not BSON at all | `unknown_format`/`malformed` → `undecodable_bson_payload` | `undecodable_bson_payload` | ok |
| string where the mapping says double | `type_mismatch` → `unexpected_value_type` | `unexpected_value_type` | ok |
| BSON document with no `v` key | `no_value` → `undecodable_bson_payload` | `unexpected_value_type` | **wrong** |
| object with an unexpected key | `bad_object` → `unexpected_object_key` | `unexpected_object_key` | ok |
| unknown control message | `control_unknown` → `unexpected_control_message` | `unexpected_control_message` | ok |
| invalid introspection | `introspection_invalid` → `invalid_introspection` | `invalid_introspection` | ok |

Both errors are the same *kind* of mistake — reasoning from the name of the Astrate reason rather
than from what upstream does. `interface_loading_failed` reads like an internal failure and was
therefore used as Astrate's fallback, when upstream in fact uses it for the ordinary case of a
device publishing to an interface it never declared. And a document with no `v` key is, to
upstream, a *value* problem rather than a decoding one. **Phase 05 settles both.**

A valid publish is recorded alongside as the acceptance row: it produced no `device_error` in any
attempt. Without it, a stack that emitted `device_error` for everything would satisfy every row
above.

### Delivery of these events is unreliable on this stack, so every row is repeated

Each provocation runs **three times**, and each row carries `attempts`/`delivered`. This is not
defensiveness: two consecutive single-pass runs of an earlier recorder *disagreed*, and a fixture
that had been written from either one would have recorded a falsehood.

Two distinct causes were found, and only the first is ours:

- **Session noise was being attributed to the provocation.** Upstream answers the session's own
  `emptyCache` control with `device_error`/`device_session_not_found` until Data Updater Plant has
  registered the session, and taking the first event to arrive credited that to whatever was being
  provoked. The recorder now correlates on the `base64_payload` upstream echoes back in the event
  metadata, and keeps everything it rejected under `uncorrelated_frames` so the correlation can be
  audited rather than trusted.
- **Upstream itself drops some of these events.** Every provocation is *detected* — each one
  appears in `data-updater-plant`'s log with a `tag=` naming exactly the error. But most of these
  error classes also make upstream force a clean session, and on this stack that path times out on
  an RPC to VerneMQ and crashes the handling process (`tag=data_upd_crash_detected`), which is
  when the event is lost. **Whether that is Astarte's behaviour or this deployment's is not
  established**, so `delivered: 2/3` is recorded as what it is and no row claims upstream is
  silent. A row with `delivered: 0` would mean "never observed in three tries", never "upstream
  emits nothing".

### What could not be reached, and is therefore not recorded

The realm holds only device-owned datastream interfaces, so `write_on_server_owned_interface` was
never provoked. `value_size_exceeded` did not appear either — a 100 KB string was accepted rather
than refused, so upstream's limit is higher than that and was not bracketed. Both remain
**unverified guesses in Astrate's table**, exactly as they were before this recording; the
recording's value here is knowing precisely which two rows it did not settle.

### The `a_ch` grammar, which settles M11 decision 5

M11 could not decide how narrowly the Channels claim is read, because a blanket `.*::.*` grant
authorizes every reading of the WATCH path equally. Upstream answers it plainly, and the rows are
paired so a refusal cannot be explained by a server that refuses everything:

| operation | claim | |
| --- | --- | --- |
| join a room | `.*::.*` | refused |
| join a room | `.*` | refused |
| join a room | `JOIN::.*` | accepted |
| join a room | `JOIN::<that room>` | accepted |
| join a room | `JOIN::<another room>` | refused |
| watch a data trigger | `WATCH::.*` | accepted |

**A blanket grant is not a grant**, and the reason is not the one first written here. Upstream
*partitions* the `a_ch` list by an **exact** match on the verb field: an entry is split on `:` into
three parts and kept only if its first field is literally `JOIN` or `WATCH`, everything else being
discarded. The path regex is then matched within the chosen bucket. So `.*::.*` authorizes nothing
because `.*` is not the string `JOIN` — not, as the first version of this note claimed, because the
room name lacks a `::`. That earlier story is refuted by the `.*` row in the same table: under it,
`.*` would have matched the bare room name and been accepted, and it was refused. Confirmed against
upstream's source (`socket_guardian.ex`, `extract_authorization_paths/2`, pinned `^match_prefix`).

### What the `WATCH` claim is matched against

The rows above all used `WATCH::.*`, which accepts whatever string upstream builds and therefore
cannot say what that string *is*. These use narrow claims, so exactly one of a pair can be accepted:

| trigger | claim | |
| --- | --- | --- |
| data trigger on `/value` | `WATCH::<device>/<interface>` | refused |
| data trigger on `/value` | `WATCH::<device>/<interface>/value` | accepted |
| device trigger | `WATCH::<that device>` | accepted |
| device trigger | `WATCH::<another device>` | refused |

So a data trigger's authorization path is **`<device_id>/<interface_name><match_path>`** — the
match path is appended, with no separator because it already begins with a slash — and a device
trigger's is the **bare device id**. Astrate built `<device_id>/<interface_name>` for data triggers
until M12 phase 06b, which is wrong in *both* directions: it refused claims upstream accepts and
accepted claims upstream refuses.

The group shapes (`groups/<name>/<interface><match_path>` and `groups/<name>`) were measured on
2026-08-22, after a group was provisioned in the realm, and both match what Astrate builds — the
last source-read claim in this file is now an observation. The same recording settled two payload
rules a client author needs:

- **`group_name` lives at the watch payload's top level.** A `group_name` nested inside
  `simple_trigger` is refused by upstream's changeset — `{"errors":{"group_name":["must be present
  if device_id is not set"],"device_id":[...]}}` — before authorization is consulted. Astrate used
  to read it only from `simple_trigger`, which meant an upstream-shaped group watch silently
  degraded into a device-shaped path check; fixed in the reconciliation.
- **A group device_trigger requires `device_id: "*"` inside `simple_trigger`.** A concrete id is
  refused with reason `device_id must be * for group triggers` — the exact mirror image of plain
  device triggers, where `"*"` is refused and a concrete id is required.

### Ownership and size (recorded 2026-08-22, issues #18/#19)

Two interfaces were installed by hand to reach what the realm's original device-owned pair could
not provoke: `org.astrate.bench.ServerOwned` (server-owned datastream) and
`org.astrate.bench.Strings` (device-owned datastream with a string endpoint), plus the group above.
The rows live in the fixture's `device_errors` section; the bisection that found the boundaries was
a scratch probe and is summarized here rather than committed as code.

- **`write_on_server_owned_interface` is confirmed verbatim** (2 of 3 attempts delivered — the
  usual flakiness). The interface was in the device's introspection, so the rejection cannot be
  anything else. Astrate's mapping needed no change.
- **`value_size_exceeded` cannot be reached over MQTT on this stack.** The transport enforces its
  own cap first: any publish whose MQTT packet exceeds 65536 bytes is ACKed by VerneMQ and then
  silently discarded — nothing stored, no `device_error`, no session change. Measured by bisection:
  a 65468-byte BSON document (string of 65455 bytes) is stored and queryable; one byte-string longer,
  65473 bytes of BSON, vanishes. At ≥ ~3 MB the broker stops being polite and closes the connection
  outright (accepted ≤ 3,088,710-byte payloads' neighbour; 3,090,000 dropped the TCP connection).
  Upstream's own validator *does* have the name — astarte_core rejects strings over 65536 bytes as
  `value_size_exceeded` — but a publish that large can never reach it over MQTT. Astrate keeps its
  own caps and emits the event upstream does not; that divergence is deliberate and recorded in
  `docs/COMPATIBILITY.md`.

The size rows pin both sides of each boundary with fixed sizes so the fixture stays deterministic;
the exact byte bounds belong to this note, since they describe this deployment's transport as much
as Astarte.

### Re-recording one section at a time

`ASTARTE_UPSTREAM_SECTION=auth|errors|extra|all` (default `all`) re-records one section and carries
the others over from the committed fixture, refusing to write a fixture with a section empty.
Delivery of `device_error` to the room is unreliable on this stack, so settling an authorization
question with a full re-run would churn `attempts`/`delivered` counts that were measured carefully
and are not what is being asked about. `extra` appends the ownership/size rows without touching the
previously measured ones.

### Setup the extra sections assume

`auth` needs a group named `probe` containing the recorder's device; `extra` needs the two extra
interfaces installed. Neither is created by `bench provision` yet — they were made once by hand via
Realm Management/AppEngine APIs (realm token minted from the state file's realm key) and survive in
the realm. If either goes missing, the recorder fails loudly rather than recording a wrong answer.

One more, which is a client-facing trap rather than an authorization rule: a `device_trigger` is
authorized only when `device_id` sits **inside `simple_trigger`** *and equals the request's own
`device_id`*. With `device_id` only at the payload's top level — where the AppEngine REST API puts
it — upstream refuses with `reason: "unauthorized"`, which sends a client author to look at their
token for what is really a payload-shape problem. A wildcard `device_id: "*"` is refused too, even
under `WATCH::.*`.

**A trap that cost this recording a pass:** `interface_major` is required as soon as
`interface_name` is not `"*"`, and omitting it is refused with
`{"errors":{"simple_trigger":{"interface_major":["can't be blank"]}}}` — a *validation* error
returned before authorization is consulted at all. To a recorder checking only for `"status":"ok"`
that is a refusal like any other, so both narrow data-trigger rows first came back "refused" and
would have been read as evidence about the authorization path. The recorder now prints a loud
`!! VALIDATION ERROR, NOT AN AUTHORIZATION REFUSAL` on any reply carrying `"errors"` rather than
`"reason"`.

### Re-recording one section at a time

Covered above — `ASTARTE_UPSTREAM_SECTION` takes `auth`, `errors`, `extra` or `all`.

### Two facts a re-recorder needs and will otherwise lose an evening to

- **The Phoenix socket must be heartbeated** (`[null,"hb","phoenix","heartbeat",{}]` every ~15s) or
  upstream closes it part-way through a run. The failure is silent: the reader ends, every later
  row records nothing, and "no event" looks exactly like a finding. The recorder now prints the
  closure loudly for this reason.
- **The introspection payload is `<interface>:<major>:<minor>`.** Sending `:0:1` instead of `:1:0`
  loads no interface at all, and *every* publish then comes back as `interface_loading_failed` —
  which looks like a coherent result, and is not one.
