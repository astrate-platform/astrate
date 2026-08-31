slug: probe-binaryblob-validation
verdict: blocked
at:  25560b4
ran: 2026-08-31T19:35:46Z on DietPi in 185s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md

# Todos
[•] Read pkg/payload/value.go for binaryblob validation logic
[ ] Read internal/engine/serverdata.go for binaryblob handling
[ ] Find upstream Astarte binaryblob validator (v1.4.0-rc.5 corrected)
[ ] Compare boundaries and report

→ Read pkg/payload/value.go
→ Read internal/engine/serverdata.go

# Todos
[✓] Read pkg/payload/value.go for binaryblob validation logic
[✓] Read internal/engine/serverdata.go for binaryblob handling
[•] Find how binaryblob is validated in payload decoders
[ ] Find upstream Astarte binaryblob validator (v1.4.0-rc.5 corrected)
[ ] Compare boundaries and report

✱ Grep "binaryblob|BinaryBlob|binary_blob" in pkg/payload · 19 matches
✱ Glob "pkg/payload/*.go" in . · 9 matches
→ Read pkg/payload/json.go [offset=140, limit=30]
→ Read pkg/payload/bson.go [offset=110, limit=40]
→ Read pkg/payload/payload_test.go [offset=310, limit=40]
Now I need to understand the upstream Astarte binaryblob validation boundaries. The task references upstream v1.4.0-rc.5's "corrected validator" ("Ensure binaryblob data is correctly validated"). Let me look at the existing tests around binaryblob size and understand what boundaries Astrate applies, then compare to upstream.

Let me look at the relevant test around line 561 (the "bigBlob" too large test) and understand the size limit for binaryblob.
→ Read pkg/payload/payload_test.go [limit=60, offset=535]
Now I need to understand the upstream Astarte binaryblob validation. The task references upstream v1.4.0-rc.5's "corrected validator" — the PR "Ensure binaryblob data is correctly validated". Let me search for this upstream.

I should investigate what upstream Astarte does for binaryblob validation. The key question: does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream's corrected validator?

Let me search the web for the upstream PR/commit.
◈ Exa Web Search "Astarte "Ensure binaryblob data is correctly validated" binaryblob validator AstarteBinaryBlob"
I found the upstream `validate_value` in `astarte_core/mapping/value_type.ex`. Let me capture the key binaryblob validation boundary from upstream:

```elixir
@blob_size 65536  # 64 KiB
...
{ v, :binaryblob} when is_binary(v) ->
  if byte_size(v) > @blob_size do
    {:error, :value_size_exceeded}
  else
    :ok
  end
```

So upstream's corrected validator enforces `byte_size(blob) > 65536` → error, i.e. maximum blob size is 65536 bytes (64 KiB). 

Now, the commit referenced (6c19770894f85b13eb3f950da7e5718010a581d2) "introduced the fix for values in individual interfaces". Let me check what that corrected.

Wait — but the task is about `v1.4.0-rc.5`. This version of astrate. Let me focus on the actual boundaries Astrate applies vs upstream.

Astrate's binaryblob validation:
- BSON path (`bson.go:120-125`): Any binary subtype accepted as raw bytes, no size limit on the blob itself (only the overall payload cap of 64KiB).
- JSON path (`json.go:153-162`): base64 standard alphabet, padded, strict. No size limit on the decoded blob itself — but the payload cap `DefaultMaxSize = 64KiB` applies to the whole JSON payload.

So Astrate's binaryblob is bounded only by the overall payload size cap (64KiB), NOT by an explicit per-blob 64KiB limit. Upstream's corrected validator enforces `byte_size(blob) > 65536` → reject (value_size_exceeded).

The key question: does the overall 64KiB payload cap in Astrate effectively bound the binaryblob to 64KiB too? In BSON, the payload includes the binary bytes + envelope. A 64KiB blob plus BSON overhead would exceed the 64KiB cap. So a BSON binaryblob of exactly 64KiB might be rejected by Astrate's payload cap even though upstream would accept it (byte_size <= 65536). Meanwhile a blob of 65536 bytes in BSON would have payload > 64KiB so rejected.

But the more important boundary difference: Upstream allows blob up to exactly 65536 bytes. Astrate, via overall payload cap, allows a blob only up to ~64KiB minus envelope overhead. There is a boundary mismatch. Let me verify what the actual payload cap in Astrate's ingestion path is.

Let me look at how payload.Decoder cap is set in the ingestion path (not serverdata which uses len(envelope)).

The task says "does Astrate accept/reject binaryblob mapping values at ingestion with the same boundaries as upstream v1.4.0-rc.5's corrected validator". Let me focus the probe.

Let me check what the upstream commit actually was about — the "Ensure binaryblob data is correctly validated". Let me check the actual commit and whether it's in astarte vs astarte_core. Actually, the value_type.ex `validate_value` with `@blob_size 65536` is the corrected validator. This is in astarte_core.

Let me look for how Astrate's ingestion sets the MaxSize, and where binaryblob validation lives. Let me also check the value_test.go and bson size cap.

Let me look at what limits apply. The upstream commit message "Ensure binaryblob data is correctly validated" — let me verify which commit that is. Let me search the git log of astrate to find if there's a reference to this task, and search upstream.

Actually, the crucial comparison: Astrate applies the binaryblob boundary only via the aggregate payload size cap (64KiB). Upstream applies an explicit per-blob 64KiB cap. The effective difference is at the boundary: a binaryblob payload whose size is between (payloadcap - overhead) and payloadcap. Let me measure precisely.

Let me find the ingestion path that uses payload.Decoder.
✱ Grep "DefaultMaxSize|MaxSize" in . · 18 matches
Interesting — there's a `.mule/cron.log` mention at line 4100 that appears to be from a *previous* run of this same or similar task! Let me check that.
→ Read .mule/cron.log [offset=4080, limit=60]
This is the cron.log of a previous run — it seems a previous mule session was working on this exact same task and got cut off. This is a record of prior runs. The `.mule/cron.log` is not something the task told me to handle, but it shows a prior partial investigation. I should not rely on it as authoritative, but it confirms the direction.

Wait — this is a `.mule/cron.log` file which is the log of past mule runs. It indicates a previous attempt at this task was interrupted mid-investigation. My job is to complete the probe.

Let me now verify the upstream validator boundaries precisely. I have the upstream `validate_value` from `value_type.ex`. The key bound: `byte_size(blob) > 65536` → reject. So upstream accepts blobs up to and including 65536 bytes.

Now let me determine Astrate's effective boundary. Astrate's binaryblob at ingestion follows the payload size cap (default 64KiB = 65536). Let me check the ingestion router path and whether binaryblob gets any per-blob cap.

Let me look at the router.go ingestion path and the decode path.
→ Read internal/engine/router.go [limit=60, offset=160]
The ingestion decoder uses `cfg.MaxPayloadBytes` (default `payload.DefaultMaxSize` = 65536). So at ingestion, the binaryblob is bounded by the whole payload being ≤ 64KiB. There is no separate per-blob cap.

Now the crucial boundary comparison. Let me write a live probe to measure exactly at what size Astrate rejects a binaryblob in the JSON and BSON formats, and compare to upstream's 65536-byte boundary.

But I can't run upstream easily. However, I have the upstream validator code clearly: `byte_size(blob) > 65536` → reject. So upstream accepts a blob of up to and **including** 65536 bytes (66*1024 = 65536... wait 64 KiB is 65536 bytes).

So upstream: accept iff `blob_size <= 65536`.

Astrate at ingestion: the blob is bounded by the whole payload ≤ 65536 bytes. For BSON, payload = BSON document overhead + blob. So a blob of exactly 65536 bytes would make payload > 65536 → rejected by Astrate, but accepted by upstream. A blob just under 65536 minus overhead would be accepted by Astrate.

Let me write a probe to measure the exact boundary for both JSON and BSON paths in Astrate.

Let me first understand the BSON encoding to construct test payloads precisely.
- For a BSON binary document, choosing a specific blob size and measuring the resulting payload length.

Let me write a probe program.

$ ls /root/astrate-mule && head -20 /root/astrate-mule/go.mod
Dockerfile
LICENSE
Makefile
README.md
bench
cmd
deploy
docker-compose.yml
docs
go.mod
go.sum
internal
migrations
pkg
test
tools
module github.com/astrate-platform/astrate

go 1.26.1

// Pinned dependency decisions — docs/ROADMAP.md §1.1 (frozen in docs/DESIGN.md).
// HTTP routing is stdlib net/http (Go >= 1.22 ServeMux patterns): deliberately no framework.
require (
	github.com/BurntSushi/toml v1.6.0 // TOML config (§5.1)
	github.com/coder/websocket v1.8.14 // live stream socket (M7)
	github.com/eclipse/paho.mqtt.golang v1.5.1 // test-only MQTT client (same family as official Go SDK)
	github.com/golang-jwt/jwt/v5 v5.3.1 // JWT verification, alg allowlist (§4.2)
	github.com/golang-migrate/migrate/v4 v4.19.1 // migrations via source/iofs + go:embed (§1.3)
	github.com/hashicorp/golang-lru/v2 v2.0.7 // verified-token LRU cache (§4.2)
	github.com/jackc/pgx/v5 v5.10.0 // Postgres/TimescaleDB driver + pgxpool (§1.3)
	github.com/mochi-mqtt/server/v2 v2.7.9 // embedded MQTT broker (§1.1)
	github.com/prometheus/client_golang v1.23.2 // metrics (§5.2)
	github.com/testcontainers/testcontainers-go v0.43.0 // T2 integration harness (§5.4)
	github.com/testcontainers/testcontainers-go/modules/postgres v0.42.0 // Timescale container module
	go.etcd.io/bbolt v1.4.3 // broker session store (§3.1; pinned over pebble)
	go.mongodb.org/mongo-driver/v2 v2.6.0 // BSON raw-document API only (§3.5.5)

Let me write a probe program to measure the exact acceptance boundary in Astrate's payload decoder for binaryblob in both JSON and BSON, using a mapping where the payload cap is the default 64KiB. The probe will determine what the largest accepted blob is and how that compares to upstream's 65536 boundary.

I need to look at how payload.Decode works and the BSON encoding for binary. Let me check the payload API.
→ Read pkg/payload/payload.go
Now I understand the mechanism. The boundary for Astrate binaryblob is entirely determined by the `gate` size cap (default 64KiB on the whole payload). Upstream's corrected validator enforces a per-blob `byte_size <= 65536`.

The difference: upstream's limit is on the *decoded blob bytes* (<= 65536). Astrate's limit is on the *whole wire payload* (<= 65536). For a BSON binaryblob, the wire payload = BSON envelope + the binary element (which has extra overhead: BSON binary type has a 5-byte length prefix + subtype byte + field name + document overhead). So the same blob that upstream accepts (exactly 64KiB) would produce a wire payload > 64KiB in BSON and be rejected by Astrate.

For the JSON path, base64 encoding inflates the payload by 4/3, so a 64KiB blob becomes ~87KiB of base64 → rejected by the 64KiB payload cap, whereas upstream would accept the 64KiB blob.

So there is a definite boundary mismatch: Astrate rejects binaryblob values that upstream would accept, at the upper boundary. This applies both to BSON (device ingestion) and JSON (server-owned ingestion), and also the serverdata path passes `MaxSize: len(envelope)` which won't gate.

Let me write a probe to measure exactly the max accepted blob size in Astrate for both paths. This satisfies the "write a twenty-line throwaway program" guidance.

Let me construct mappings. I need to see how `mapping()` helper builds a CompiledMapping in tests, and how to build a BSON binary payload.

Let me check interfaceschema.CompiledMapping construction and the test helper.
✱ Grep "func mapping|CompiledMapping\{|func rawBSON|encodeBSON" in pkg/payload/payload_test.go · 23 matches
I have enough to write a probe. Let me write a probe program in /tmp that measures the largest accepted binaryblob (in bytes) for both JSON and BSON paths under the default 64KiB cap in Astrate, and compare against upstream's 65536 boundary.

Let me write the probe using the payload package and encodeBSON (via Encode with FormatBSON), building a mapping with BinaryBlob type.

$ tail -40 /root/astrate-mule/go.mod
github.com/magiconair/properties v1.8.10 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/go-archive v0.2.0 // indirect
	github.com/moby/moby/api v1.54.2 // indirect
	github.com/moby/moby/client v0.4.0 // indirect
	github.com/moby/patternmatcher v0.6.1 // indirect
	github.com/moby/sys/sequential v0.6.0 // indirect
	github.com/moby/sys/user v0.4.0 // indirect
	github.com/moby/sys/userns v0.1.0 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nats-io/nkeys v0.4.15 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/rs/xid v1.4.0 // indirect
	github.com/shirou/gopsutil/v4 v4.26.5 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/tklauser/go-sysconf v0.3.16 // indirect
	github.com/tklauser/numcpus v0.11.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0 // indirect
	go.opentelemetry.io/otel v1.41.0 // indirect
	go.opentelemetry.io/otel/metric v1.41.0 // indirect
	go.opentelemetry.io/otel/trace v1.41.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

I'll write a probe in /tmp as a standalone Go program within the module (using a temp test-style main). Since the probe needs the module's packages, I'll create it inside /tmp but reference the module via a go.mod replace, or simpler: put it in a temp directory inside the repo? The rules say put probes in /tmp, not in the repo. But to import the module packages, I need to be in the module. 

I can create a temp package inside the repo would violate "put probes in /tmp". Better: write a throwaway test file in the payload package temporarily? But I shouldn't add files to the repo. Actually the guidance says don't leave it in the repo; I can create and then delete. But "Never write to a path outside this repository" — /tmp is outside. Hmm, contradictory with "put probes in /tmp, not in the repo."

The MULE rule says: "Put probes in /tmp, not in the repo." But also "Never write to a path outside this repository." The intent: /tmp is acceptable scratch. The "never write outside repo" likely refers to not modifying external code/deps, while /tmp probes are explicitly sanctioned. I'll use /tmp with a go.mod replace pointing at the repo module.

Let me set that up.
$ cat /root/astrate-mule/go.mod | grep -E "^go |^module|^require \(" | head; echo "---"; grep -E "go.mongodb|astarte" /root/astrate-mule/go.mod
module github.com/astrate-platform/astrate
go 1.26.1
require (
require (
require (
---
	go.mongodb.org/mongo-driver/v2 v2.6.0 // BSON raw-document API only (§3.5.5)

Let me create the probe in /tmp with a go.mod replace pointing to the repo as a local module, ensuring the go version matches.
$ go version; ls /root/astrate-mule/pkg/interfaceschema/
go version go1.26.5 linux/arm64
compat.go
compat_test.go
compile.go
compile_test.go
parse.go
parse_test.go
testdata
trie.go
trie_test.go
types.go

! permission requested: external_directory (/tmp/*); auto-rejecting
✗ mkdir -p /tmp/blobprobe && cat /root/astrate-mule/go.mod > /tmp/blobprobe/astrate.mod.bak && echo ok failed
Error: The user rejected permission to use this specific tool call.
