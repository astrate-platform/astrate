# Golden BSON payload vectors — provenance

Every `*.hex` file in this directory is a hex-encoded BSON `{v, t}` data
payload produced by the **official Astarte Go device SDK encoder**, captured
verbatim. They freeze the exact wire bytes an unmodified SDK publishes, so
`pkg/payload` is tested against the real protocol, not against itself.

## Pinned sources

| Component | Version | Role |
|---|---|---|
| `github.com/astarte-platform/astarte-device-sdk-go` | **v0.90.2** (latest tagged release) | Encoder under capture: `device/protocol_mqtt_v1.go`, `enqueueMqttV1Message` (lines 417–428) |
| `github.com/astarte-platform/astarte-go` | v0.90.4 | `interfaces.NormalizePayload`, exactly as the SDK calls it |
| `go.mongodb.org/mongo-driver` | v1.8.0 | `bson.Marshal`, exactly as the SDK pins it |

The SDK's encode path is:

```go
payload := map[string]interface{}{"v": interfaces.NormalizePayload(values, false)}
if !timestamp.IsZero() {
	payload["t"] = timestamp.UTC()
}
doc, err := bson.Marshal(payload)
```

`enqueueMqttV1Message` is unexported and reachable only through a connected
`Device`, so the capture harness (`tools/bsoncapture`, a standalone Go module
kept out of the main Astrate module) invokes that exact expression against
the exact dependency versions the SDK declares in its `go.mod`.

## Regenerating

```sh
cd tools/bsoncapture
go run . ../../pkg/payload/testdata/bson
```

Note: regeneration changes bytes without changing semantics — see the field
order note below — so vectors are only regenerated when the pinned SDK
version is deliberately upgraded.

## File naming and contents

- `<type>.hex` — value only (no explicit timestamp).
- `<type>_t.hex` — same value plus `t` = `2026-06-10T12:34:56.789Z`.
- One vector per Astarte mapping type (`double` … `datetimearray`), plus
  `object` (object aggregation: `lat`/`lon` doubles, `samples` integer,
  `ok` boolean) and `doublearray_empty` (zero-length array).

Values are chosen to pin edge behaviour: `longinteger` is 2^53 + 1 (not
float64-representable), `longintegerarray` spans ±2^53±1 and `MaxInt64`,
`binaryblobarray` contains an empty blob, `string` contains multi-byte UTF-8.

## Field order caveat (documented upstream behaviour)

The SDK marshals a Go **map**, and mongo-driver v1 encodes map keys in Go's
randomized map iteration order — so across `*_t.hex` vectors the `t` element
sometimes precedes `v` (e.g. `double_t.hex` captured t-first, `boolean_t.hex`
v-first). This is real on-the-wire SDK behaviour and is exactly why the
Astrate decoder locates `v`/`t` by raw-document lookup instead of position.
Astrate's own encoder always emits `v` then `t`; conformance tests therefore
compare re-encoded documents element-wise/semantically, not byte-for-byte
against these captures.
