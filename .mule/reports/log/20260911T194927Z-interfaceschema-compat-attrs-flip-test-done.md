slug: interfaceschema-compat-attrs-flip-test
verdict: done
at:  4fe1f4f
ran: 2026-09-11T19:49:28Z on DietPi in 329s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read pkg/interfaceschema/compat_test.go
→ Read pkg/interfaceschema/compat.go
✱ Glob "pkg/interfaceschema/*.go" 12 matches
This is a test-extension task: pin the four untested immutable attributes (retention, expiry, database_retention_policy/ttl, allow_unset) in the upgrade table. Let me check the Mapping struct and parse validation to write valid docs.
✱ Grep "database_retention|DatabaseRetention|allow_unset|AllowUnset|retention|Retention|expiry|Expiry" in pkg/interfaceschema · 100 matches
→ Read pkg/interfaceschema/parse.go

← Edit pkg/interfaceschema/compat_test.go
Index: /root/astrate-mule/pkg/interfaceschema/compat_test.go
===================================================================
--- /root/astrate-mule/pkg/interfaceschema/compat_test.go
+++ /root/astrate-mule/pkg/interfaceschema/compat_test.go
@@ -22,8 +22,71 @@
 		{"endpoint": "/%{sensor_id}/status", "type": "string"}
 	]
 }`
 
+// pinnedV10 pins the three datastream-only immutable attribute pairs at
+// non-zero values, so the flip cases below exercise a real mutation of one
+// attribute instead of a default-to-value transition.
+const pinnedV10 = `{
+	"interface_name": "com.astrate.test.Upgrade",
+	"version_major": 1,
+	"version_minor": 0,
+	"type": "datastream",
+	"ownership": "device",
+	"mappings": [
+		{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
+		 "retention": "stored", "expiry": 3600,
+		 "database_retention_policy": "use_ttl", "database_retention_ttl": 3600},
+		{"endpoint": "/%{sensor_id}/status", "type": "string"}
+	]
+}`
+
+// pinnedV11 is pinnedV10 with the minor bumped: every pinned immutable
+// attribute stays put, which an upgrade must accept. It is the
+// at-immutability acceptance twin of the retention/expiry/database-retention
+// flip cases below.
+const pinnedV11 = `{
+	"interface_name": "com.astrate.test.Upgrade",
+	"version_major": 1,
+	"version_minor": 1,
+	"type": "datastream",
+	"ownership": "device",
+	"mappings": [
+		{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
+		 "retention": "stored", "expiry": 3600,
+		 "database_retention_policy": "use_ttl", "database_retention_ttl": 3600},
+		{"endpoint": "/%{sensor_id}/status", "type": "string"}
+	]
+}`
+
+// pinnedPropsV10 is a properties interface carrying allow_unset (a
+// properties-only immutable attribute) at its non-zero value.
+const pinnedPropsV10 = `{
+	"interface_name": "com.astrate.test.Upgrade",
+	"version_major": 1,
+	"version_minor": 0,
+	"type": "properties",
+	"ownership": "device",
+	"mappings": [
+		{"endpoint": "/%{sensor_id}/enabled", "type": "boolean", "allow_unset": true},
+		{"endpoint": "/%{sensor_id}/mode", "type": "string"}
+	]
+}`
+
+// pinnedPropsV11 is the at-immutability acceptance twin of the allow_unset
+// flip case: allow_unset stays true across the minor bump.
+const pinnedPropsV11 = `{
+	"interface_name": "com.astrate.test.Upgrade",
+	"version_major": 1,
+	"version_minor": 1,
+	"type": "properties",
+	"ownership": "device",
+	"mappings": [
+		{"endpoint": "/%{sensor_id}/enabled", "type": "boolean", "allow_unset": true},
+		{"endpoint": "/%{sensor_id}/mode", "type": "string"}
+	]
+}`
+
 func mustParse(t *testing.T, doc string) *interfaceschema.Interface {
 	t.Helper()
 	iface, err := interfaceschema.ParseInterface([]byte(doc))
 	if err != nil {

Task understood: extend the upgrade table to pin the four untested immutable attributes — datastream retention/expiry and the database_retention_policy/ttl pair, plus properties allow_unset — each with a flip-rejection case and an at-immutability acceptance twin. The rejection logic already exists in `sameMappingAttributes` (compat.go:109-118), so this is a test-only regression pin.

Let me add the pinned fixture docs first:
← Edit pkg/interfaceschema/compat_test.go
Index: /root/astrate-mule/pkg/interfaceschema/compat_test.go
===================================================================
--- /root/astrate-mule/pkg/interfaceschema/compat_test.go
+++ /root/astrate-mule/pkg/interfaceschema/compat_test.go
@@ -288,8 +288,111 @@
 		}`,
 		wantSub:   "changed explicit_timestamp",
 		wantClass: interfaceschema.ErrIncompatibleEndpointChange,
 	},
+	{
+		name:    "retention at immutability accepted",
+		oldDoc:  pinnedV10,
+		next:    pinnedV11,
+	},
+	{
+		name:   "retention mutated",
+		oldDoc: pinnedV10,
+		next: `{
+			"interface_name": "com.astrate.test.Upgrade",
+			"version_major": 1, "version_minor": 1,
+			"type": "datastream", "ownership": "device",
+			"mappings": [
+				{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
+				 "retention": "volatile", "expiry": 3600,
+				 "database_retention_policy": "use_ttl", "database_retention_ttl": 3600},
+				{"endpoint": "/%{sensor_id}/status", "type": "string"}
+			]
+		}`,
+		wantSub:   "changed retention",
+		wantClass: interfaceschema.ErrIncompatibleEndpointChange,
+	},
+	{
+		name:    "expiry at immutability accepted",
+		oldDoc:  pinnedV10,
+		next:    pinnedV11,
+	},
+	{
+		name:   "expiry mutated",
+		oldDoc: pinnedV10,
+		next: `{
+			"interface_name": "com.astrate.test.Upgrade",
+			"version_major": 1, "version_minor": 1,
+			"type": "datastream", "ownership": "device",
+			"mappings": [
+				{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
+				 "retention": "stored", "expiry": 7200,
+				 "database_retention_policy": "use_ttl", "database_retention_ttl": 3600},
+				{"endpoint": "/%{sensor_id}/status", "type": "string"}
+			]
+		}`,
+		wantSub:   "changed expiry",
+		wantClass: interfaceschema.ErrIncompatibleEndpointChange,
+	},
+	{
+		name:    "database retention at immutability accepted",
+		oldDoc:  pinnedV10,
+		next:    pinnedV11,
+	},
+	{
+		name:   "database_retention_policy mutated",
+		oldDoc: pinnedV10,
+		next: `{
+			"interface_name": "com.astrate.test.Upgrade",
+			"version_major": 1, "version_minor": 1,
+			"type": "datastream", "ownership": "device",
+			"mappings": [
+				{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
+				 "retention": "stored", "expiry": 3600,
+				 "database_retention_policy": "no_ttl"},
+				{"endpoint": "/%{sensor_id}/status", "type": "string"}
+			]
+		}`,
+		wantSub:   "changed database_retention_policy",
+		wantClass: interfaceschema.ErrIncompatibleEndpointChange,
+	},
+	{
+		name:   "database_retention_ttl mutated",
+		oldDoc: pinnedV10,
+		next: `{
+			"interface_name": "com.astrate.test.Upgrade",
+			"version_major": 1, "version_minor": 1,
+			"type": "datastream", "ownership": "device",
+			"mappings": [
+				{"endpoint": "/%{sensor_id}/value", "type": "double", "reliability": "guaranteed",
+				 "retention": "stored", "expiry": 3600,
+				 "database_retention_policy": "use_ttl", "database_retention_ttl": 7200},
+				{"endpoint": "/%{sensor_id}/status", "type": "string"}
+			]
+		}`,
+		wantSub:   "changed database_retention_ttl",
+		wantClass: interfaceschema.ErrIncompatibleEndpointChange,
+	},
+	{
+		name:    "allow_unset at immutability accepted",
+		oldDoc:  pinnedPropsV10,
+		next:    pinnedPropsV11,
+	},
+	{
+		name:   "allow_unset mutated",
+		oldDoc: pinnedPropsV10,
+		next: `{
+			"interface_name": "com.astrate.test.Upgrade",
+			"version_major": 1, "version_minor": 1,
+			"type": "properties", "ownership": "device",
+			"mappings": [
+				{"endpoint": "/%{sensor_id}/enabled", "type": "boolean", "allow_unset": false},
+				{"endpoint": "/%{sensor_id}/mode", "type": "string"}
+			]
+		}`,
+		wantSub:   "changed allow_unset",
+		wantClass: interfaceschema.ErrIncompatibleEndpointChange,
+	},
 }
 
 for _, tc := range cases {
 	t.Run(tc.name, func(t *testing.T) {

Now the new table cases, inserted after the `explicit_timestamp mutated` case:
Now run the gate — gofmt and the package tests:
$ gofmt -l pkg/interfaceschema/compat_test.go && go test ./pkg/interfaceschema/
pkg/interfaceschema/compat_test.go
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	0.164s

gofmt flags the file (field alignment). Fix and re-run:
$ gofmt -w pkg/interfaceschema/compat_test.go && gofmt -l . && go test ./pkg/interfaceschema/
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)

gofmt is clean. Let me force a non-cached run to confirm the final file content passes:
$ go test ./pkg/interfaceschema/ -count=1 -run 'TestCheckMinorUpgrade' -v 2>&1 | tail -40
=== RUN   TestCheckMinorUpgradeTable/allow_unset_mutated
--- PASS: TestCheckMinorUpgradeTable (0.02s)
    --- PASS: TestCheckMinorUpgradeTable/additive_mapping_with_minor_bump (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/doc_and_description_changes_are_non-semantic (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/placeholder_rename_is_non-semantic (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/minor_not_increased (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/minor_downgraded (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/major_changed (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/name_changed (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/type_changed (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/ownership_changed (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/mapping_removed (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/value_type_mutated (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/reliability_mutated (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/explicit_timestamp_mutated (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/retention_at_immutability_accepted (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/retention_mutated (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/expiry_at_immutability_accepted (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/expiry_mutated (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/database_retention_at_immutability_accepted (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/database_retention_policy_mutated (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/database_retention_ttl_mutated (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/allow_unset_at_immutability_accepted (0.00s)
    --- PASS: TestCheckMinorUpgradeTable/allow_unset_mutated (0.00s)
=== RUN   TestCheckMinorUpgradeAggregationChange
--- PASS: TestCheckMinorUpgradeAggregationChange (0.00s)
=== RUN   TestCheckMinorUpgradeNil
--- PASS: TestCheckMinorUpgradeNil (0.00s)
=== RUN   TestCheckMinorUpgradeRequiredEncrypted
=== RUN   TestCheckMinorUpgradeRequiredEncrypted/required_flip_rejected
=== RUN   TestCheckMinorUpgradeRequiredEncrypted/required_flip_rejected_acceptance_twin
=== RUN   TestCheckMinorUpgradeRequiredEncrypted/encrypted_flip_rejected
=== RUN   TestCheckMinorUpgradeRequiredEncrypted/encrypted_flip_rejected_acceptance_twin
--- PASS: TestCheckMinorUpgradeRequiredEncrypted (0.01s)
    --- PASS: TestCheckMinorUpgradeRequiredEncrypted/required_flip_rejected (0.00s)
    --- PASS: TestCheckMinorUpgradeRequiredEncrypted/required_flip_rejected_acceptance_twin (0.00s)
    --- PASS: TestCheckMinorUpgradeRequiredEncrypted/encrypted_flip_rejected (0.00s)
    --- PASS: TestCheckMinorUpgradeRequiredEncrypted/encrypted_flip_rejected_acceptance_twin (0.00s)
PASS
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	0.045s

All 22 table cases pass, including the 9 new ones (matching the existing flip + acceptance-twin convention in `TestCheckMinorUpgradeRequiredEncrypted`). Full repo test sweep to close the gate:
$ go vet ./... && go test ./... 2>&1 | tail -20
ok  	github.com/astrate-platform/astrate/internal/flow	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/astartesource	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/container	(cached)
ok  	github.com/astrate-platform/astrate/internal/flow/blocks/virtualdevicepool	(cached)
ok  	github.com/astrate-platform/astrate/internal/flowapi	(cached)
ok  	github.com/astrate-platform/astrate/internal/housekeeping	(cached)
ok  	github.com/astrate-platform/astrate/internal/httpx	(cached)
ok  	github.com/astrate-platform/astrate/internal/observability	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing	(cached)
ok  	github.com/astrate-platform/astrate/internal/pairing/ca	(cached)
ok  	github.com/astrate-platform/astrate/internal/realm	(cached)
ok  	github.com/astrate-platform/astrate/internal/store	(cached)
ok  	github.com/astrate-platform/astrate/internal/swagger	(cached)
ok  	github.com/astrate-platform/astrate/internal/testutil	(cached)
?   	github.com/astrate-platform/astrate/migrations	[no test files]
ok  	github.com/astrate-platform/astrate/pkg/astarteapi	(cached)
ok  	github.com/astrate-platform/astrate/pkg/deviceid	(cached)
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	(cached)
ok  	github.com/astrate-platform/astrate/pkg/payload	(cached)

All gates pass. Report:

```
Done: Extended TestCheckMinorUpgradeTable with 9 cases pinning the four untested immutable attributes — retention, expiry, database_retention_policy/ttl, allow_unset — each with a flip case (rejected as ErrIncompatibleEndpointChange, with the "changed <attr>" substring) and an at-immutability acceptance twin, matching the existing TestCheckMinorUpgradeRequiredEncrypted convention.
Files: pkg/interfaceschema/compat_test.go
Verified: gofmt -l . -> clean; go vet ./... -> pass; go test ./pkg/interfaceschema/ -count=1 -run TestCheckMinorUpgrade -v -> 22/22 pass; go test ./... -> pass. Ran without -race (per MULE.md allowance; no goroutine/lock code touched).
Unsure: nothing
Follow-ups: none
```
