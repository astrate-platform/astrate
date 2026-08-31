slug: probe-interface-default-values
verdict: blocked
at:  eeb69d4
ran: 2026-08-31T17:21:39Z on DietPi in 361s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
The task is a `[readonly]` probe: I need to check whether Astrate's `GET /realmmanagement/v1/<realm>/interfaces/<name>/<major>` returns the same defaulted mapping-parameter values as upstream Astarte after commit 3f0b864. I must change nothing — just report.

Let me first understand what the upstream commit 3f0b864 does, then examine Astrate's endpoint.
✱ Grep "realmmanagement.*interfaces" 11 matches
% WebFetch https://github.com/astarte-platform/astarte/commit/3f0b864
✱ Grep "func.*getInterface" 1 match
✱ Grep "func.*GetInterface" 3 matches
→ Read internal/realm/http.go [offset=163, limit=50]
→ Read internal/realm/service.go [offset=220, limit=50]
→ Read internal/store/interfaces.go [offset=243, limit=60]
✱ Grep "InstallInterface|InstallRealmInterface" 57 matches
→ Read internal/store/interfaces.go [limit=60, offset=1]
✱ Grep "ParseInterface|parseMapping|Mapping\s*struct" 52 matches
→ Read internal/store/interfaces.go [limit=80, offset=61]
→ Read pkg/interfaceschema/parse.go
✱ Grep "Required.*bool|required.*json" in pkg/interfaceschema · 0 matches
• Check upstream Mapping defaults General Agent
✓ Check upstream Mapping defaults General Agent
✱ Grep "CREATE TABLE.*endpoints" 1 match
✱ Grep "required.*json|encrypted.*json" in pkg/interfaceschema · 0 matches
→ Read migrations/000002_metadata.up.sql [limit=30, offset=33]
Let me check what happens with `required` and `encrypted` fields, and also look at what the endpoint table stores.
✱ Grep "insertEndpoint" 4 matches
→ Read internal/store/interfaces.go [limit=30, offset=329]
Let me check what the zero value of `DatabaseRetentionPolicy` and other enum types are.
✱ Grep "DatabaseRetentionPolicy|type Reliability|type Retention" in pkg/interfaceschema · 27 matches
Now let me check what the GET response actually looks like for a concrete interface, and also verify the upstream response format more carefully.
→ Read pkg/interfaceschema/types.go [limit=50, offset=151]
Now I have enough to write the report. Let me verify one more thing — whether there are any additional mapping parameters in upstream that Astrate is missing entirely.
% WebFetch https://raw.githubusercontent.com/astarte-platform/astarte_core/master/lib/astarte_core/mapping.ex
Now I have the complete picture. Let me summarize my findings.

**Upstream behavior (after3f0b864):**
- The `GET /realmmanagement/v1/<realm>/interfaces/<name>/<major}` endpoint reconstructs each mapping from DB rows via `Mapping.from_db_result!/1`, which applies explicit defaults for `database_retention_policy`→`no_ttl`, `required`→`false`, `encrypted`→`false`.
- The Jason encoder then serializes using `add_key_if_not_default` — any mapping parameter at its default value is **omitted** from the JSON response.

**Astrate behavior:**
- The GET endpoint returns `si.Definition` — the raw JSON blob stored in the `interfaces.definition` column at install time. Whatever the client POSTed is what gets returned, with no normalization.

**Comparison for default-valued parameters:**

| Parameter | Default | Upstream GET | Astrate GET |
|---|---|---|---|
| `reliability` | `unreliable` | omitted | omitted if not in POST |
| `retention` | `discard` | omitted | omitted if not in POST |
| `expiry` | `0` | omitted | omitted if not in POST |
| `database_retention_policy` | `no_ttl` | omitted | omitted if not in POST |
| `database_retention_ttl` | `nil` | omitted | omitted if not in POST |
| `allow_unset` | `false` | omitted | omitted if not in POST |
| `explicit_timestamp` | `false` | omitted | omitted if not in POST |
| `required` | `false` | omitted | N/A (unsupported) |
| `encrypted` | `false` | omitted | N/A (unsupported) |

The result is identical for the common case (client omits default-valued fields in the POST body): both systems return a mapping JSON without those keys. The one divergence is if a client explicitly includes a default-valued key (e.g., `"reliability": "unreliable"`) in the POST — upstream would strip it, Astrate keeps it. But this is an unusual edge case.

Done: Astrate's GET returns the same defaulted mapping-parameter values as upstream after3f0b864 for the normal case (omitted defaults). The one gap is that Astrate lacks `required` and `encrypted` fields entirely (rejected on POST via `DisallowUnknownFields`), but this doesn't affect the GET response shape since upstream also omits them when they are `false`.
Files: none (readonly probe)
Verified: code review only, no running commands — this is a design-level compatibility check
Unsure: Whether upstream clients ever rely on `required`/`encrypted` being present in POST payloads — if they do, Astrate would reject the install
Follow-ups: none
