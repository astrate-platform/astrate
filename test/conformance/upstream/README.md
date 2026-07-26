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

**A blanket grant is not a grant.** Upstream matches the claim against `JOIN::<room name>`, so
`.*::.*` fails for want of a `::` the room name does not contain — the opposite of the permissive
reading M11 assumed. The verb is also matched per room, not merely required to be present.

One more, which is a client-facing trap rather than an authorization rule: a `device_trigger` is
authorized only when `device_id` sits **inside `simple_trigger`**. With `device_id` only at the
payload's top level — where the AppEngine REST API puts it — upstream refuses with
`reason: "unauthorized"`, which sends a client author to look at their token for what is really a
payload-shape problem. A wildcard `device_id: "*"` is refused too, even under `WATCH::.*`.

### Two facts a re-recorder needs and will otherwise lose an evening to

- **The Phoenix socket must be heartbeated** (`[null,"hb","phoenix","heartbeat",{}]` every ~15s) or
  upstream closes it part-way through a run. The failure is silent: the reader ends, every later
  row records nothing, and "no event" looks exactly like a finding. The recorder now prints the
  closure loudly for this reason.
- **The introspection payload is `<interface>:<major>:<minor>`.** Sending `:0:1` instead of `:1:0`
  loads no interface at all, and *every* publish then comes back as `interface_loading_failed` —
  which looks like a coherent result, and is not one.
