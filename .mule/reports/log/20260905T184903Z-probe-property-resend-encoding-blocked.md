slug: probe-property-resend-encoding
verdict: blocked
at:  86cecdc
ran: 2026-09-05T18:49:03Z on DietPi in 100s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/engine/control.go [offset=130, limit=80]
✱ Grep "resendServerProperties|rehydrate" in internal/engine · 7 matches
✱ Grep "binaryblob|BinaryBlob|binary_blob" in pkg/payload · 19 matches
This is a `[readonly]` probe: check whether Astrate's `resendServerProperties`/rehydrate path emits `binaryblob` → BSON subtype-0 binary and `datetime` → UTC datetime correctly, matching upstream v1.3.0. No code changes allowed.

Let me read the relevant source files.
✱ Grep "datetime|DateTime" in pkg/payload · 78 matches
→ Read internal/engine/control.go [offset=225, limit=90]
Now let me read the `rehydrateValue` function and the BSON encoding path for binaryblob/datetime.
→ Read pkg/payload/bson.go
→ Read pkg/payload/value.go
✱ Grep "TypeBinaryGeneric" in pkg/payload · 1 match
→ Read pkg/payload/testdata/bson/README.md
Let me check the BSON test data and the binary subtype used.
✱ Grep "3q0=|3q2\+|binaryblob.*bson|bson.*binaryblob" in pkg/payload · 6 matches
Let me check the BSON encode test to confirm binary subtype-0 is asserted, and confirm the datetime encoding is UTC.
✱ Grep "BinaryOK|binaryblob|BinaryBlob|encodeBSON|Encode\(|nonce|subtype" in pkg/payload · 62 matches
I have enough to answer this probe. Let me analyze the path thoroughly.

**The resend path** (`internal/engine/control.go:144-181`):
1. `e.st.ListServerOwnedProperties` → each property's stored `Value` is a raw byte slice.
2. `e.rehydrateValue(p.Value, mapping)` (control.go:230-244): wraps the stored JSON in `{"v":...}` and runs the JSON-profile decoder against the mapping. For a `binaryblob`, the stored value is a base64 string → decoded to `[]byte`. For `datetime`, stored RFC 3339 string → decoded to `time.Time` (UTC, truncated to ms via `checkDateTime`).
3. `payload.Encode(val, nil, format)` → for BSON, `encodeBSON`.

**BSON encoding** (`bson.go:184-293`):
- `[]byte` → `bson.TypeBinary` with subtype `bson.TypeBinaryGeneric` (bson.go:263-266). `TypeBinaryGeneric` = 0x00, which is subtype 0. ✓
- `time.Time` → `bson.TypeDateTime` with `t.UnixMilli()` (bson.go:267-273). The value is normalized to UTC+truncated-to-ms via `checkDateTime` (bson.go:268). ✓

Upstream's fix: a raw stored binaryblob is wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding so it ships as a BSON binary, not a string. Astrate's path rehydrates the blob to `[]byte` (the Go type) before encoding to BSON, so it goes through the `[]byte` case emitting generic binary. That matches subtype 0.

datetime: rehydrated to `time.Time` UTC, encoded as UTC datetime. Matches.

Let me verify `bson.TypeBinaryGeneric` is indeed subtype 0 and confirm the round-trip test at payload_test.go:672 covers binaryblob. Let me also double-check the rehydrate path for datetime key detail — the stored value is milliseconds-truncated UTC.

I have enough to report. This is an investigate-only `[readonly]` task; no code changes, no test to add.

Let me verify the `bson.TypeBinaryGeneric` constant value to be precise.

$ grep -rn "TypeBinaryGeneric" $(go env GOMODCACHE 2>/dev/null)/go.mongodb.org/mongo-driver*/bson/*.go 2>/dev/null | head; grep -rn "TypeBinaryGeneric" /root/go/pkg/mod/go.mongodb.org/mongo-driver*/bson/ 2>/dev/null | head
(no output)

→ Read pkg/payload/payload_test.go [limit=40, offset=645]
I have enough to report. Let me verify the `bson.TypeBinaryGeneric` value in the vendored mongo driver to state subtype 0 precisely.

$ find / -path "*mongo-driver*/bson" -type d 2>/dev/null | head; echo "---"; grep -rn "TypeBinaryGeneric\|TypeBinary\\b" $(go env GOMODCACHE)/go.mongodb.org 2>/dev/null | grep "bson/" | head
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson
---
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/extjson_parser.go:127:				t = TypeBinary
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/extjson_parser.go:272:	case TypeBinary, TypeRegex, TypeTimestamp, TypeDBPointer:
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/extjson_parser.go:283:		if t == TypeBinary && ejp.s == jpsSawValue {
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/extjson_parser.go:319:					return nil, invalidJSONErrorForType("$uuid and value and then }", TypeBinary)
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/extjson_parser.go:347:				return nil, invalidJSONErrorForType(",", TypeBinary)
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/extjson_parser.go:356:				return nil, invalidJSONErrorForType("$type", TypeBinary)
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/extjson_parser.go:366:				return nil, invalidJSONErrorForType("2 key-value pairs and then }", TypeBinary)
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/type_test.go:23:		{"binary", TypeBinary, "binary"},
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/primitive_codecs_test.go:518:						BSONType: TypeBinary,
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/primitive_codecs_test.go:527:					RawValue{Type: TypeBinary, Value: bsoncore.AppendBinary(nil, 0xFF, []byte{0x01, 0x02, 0x03})},

$ grep -rn "TypeBinaryGeneric\b" /root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/primitive.go /root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/internal/bsoncore/*.go 2>/dev/null | head; echo "---v1---"; find / -path "*mongo-driver/v1*/bson/primitive.go" 2>/dev/null | head -1
---v1---

$ grep -rn "TypeBinaryGeneric" /root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/ 2>/dev/null | head
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/empty_interface_codec.go:106:		if binElem.Subtype == TypeBinaryGeneric || binElem.Subtype == TypeBinaryBinaryOld {
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/types.go:66:	TypeBinaryGeneric     byte = 0x00
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/byte_slice_codec.go:67:		if subtype != TypeBinaryGeneric && subtype != TypeBinaryBinaryOld {
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/unmarshal_value_test.go:66:			bytes:    bsoncore.AppendBinary(nil, TypeBinaryGeneric, []byte(strings.Repeat("t", 4096))),
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/encoder_test.go:237:				AppendBinary("myBytes", TypeBinaryGeneric, []byte{}).
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/decoder_test.go:525:				AppendBinary("myBinary", TypeBinaryGeneric, []byte{}).
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/unmarshal_test.go:55:			bytes:    bsoncore.AppendBinary(nil, TypeBinaryGeneric, []byte("hello world")),
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/slice_codec.go:122:		if subtype != TypeBinaryGeneric && subtype != TypeBinaryBinaryOld {
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/default_value_decoders.go:472:		if subtype != TypeBinaryGeneric && subtype != TypeBinaryBinaryOld {
/root/go/pkg/mod/go.mongodb.org/mongo-driver/v2@v2.6.0/bson/default_value_decoders.go:1172:		if subtype != TypeBinaryGeneric && subtype != TypeBinaryBinaryOld {

Confirmed: `bson.TypeBinaryGeneric = 0x00` = subtype 0. I have everything needed to report.

**Analysis:**

Upstream v1.3.0 fixed the outbound server-property resend so a stored binaryblob is wrapped `%Cyanide.Binary{subtype: :generic}` before BSON encoding → ships as BSON binary (subtype 0), not a string. datetime ships as UTC datetime.

Astrate's resendServerProperties/rehydrate path (`internal/engine/control.go:144-181`):

1. **binaryblob:** stored JSON property value is base64 (per rehydrate docstring, docs/DESIGN §2.3/§3.5.3). `rehydrateValue` (control.go:230-244) wraps it `{"v":...}` and runs the JSON-profile decoder → `[]byte`. Then `payload.Encode` → BSON `encodeBSON` hits the `case []byte:` (bson.go:259-266) which emits `bson.TypeBinary` with subtype `bson.TypeBinaryGeneric` (=0x00, generic). **Matches upstream: BSON binary, subtype 0, not a string.** ✓

2. **datetime:** stored RFC 3339 string rehydrates through `checkDateTime` → `time.Time` in UTC, truncated to ms (value.go:285-291). BSON encode of `time.Time` (bson.go:267-273) emits `bson.TypeDateTime` with `t.UnixMilli()` — UTC datetime. **Matches upstream: UTC datetime.** ✓

The rehydrate route re-derives the typed Go value (`[]byte`/`time.Time`) from the stored JSON rendering before BSON encoding, which is exactly the semantic equivalent of upstream's `%Cyanide.Binary` wrap — both land on BSON generic binary + UTC datetime. The round-trip test at payload_test.go:655-656 / 662-663 already covers BSON binaryblob/datetime encode→decode for both formats.

No divergence found. Investigate-only; no change made.

```
Done: Probe confirms Astrate's resendServerProperties/rehydrate path emits binaryblob → BSON subtype-0 binary and datetime → UTC datetime, matching upstream v1.3.0 — no divergence.
Files: (none — readonly probe)
Verified: go lookup bson.TypeBinaryGeneric=0x00 in mongo-driver v2.6.0; read internal/engine/control.go:144-244, pkg/payload/bson.go, value.go, payload_test.go:645-684 -> pass (no build/race run; no code changed)
Unsure: nothing
Follow-ups: none
```
