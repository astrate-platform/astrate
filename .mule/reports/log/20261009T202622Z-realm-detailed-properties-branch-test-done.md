slug: realm-detailed-properties-branch-test
verdict: done
at:  76f80e9
ran: 2026-10-09T20:26:22Z on DietPi in 88s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/realm/detailed_required_encrypted_test.go
→ Read internal/realm/detailed.go
→ Read internal/realm/detailed_test.go

$ rg -n "AllowUnset|allow_unset" /root/astrate-mule/pkg/interfaceschema/ | head -30
/root/astrate-mule/pkg/interfaceschema/compat.go:117:	case nm.AllowUnset != om.AllowUnset:
/root/astrate-mule/pkg/interfaceschema/compat.go:118:		attr = "allow_unset"
/root/astrate-mule/pkg/interfaceschema/parse_test.go:124:	t.Run("SamplingRate allow_unset", func(t *testing.T) {
/root/astrate-mule/pkg/interfaceschema/parse_test.go:133:			if !m.AllowUnset {
/root/astrate-mule/pkg/interfaceschema/parse_test.go:134:				t.Errorf("mapping %s: AllowUnset = false, want true", m.Endpoint)
/root/astrate-mule/pkg/interfaceschema/parse_test.go:181:	if m.Expiry != 0 || m.DatabaseRetentionTTL != 0 || m.ExplicitTimestamp || m.AllowUnset {
/root/astrate-mule/pkg/interfaceschema/trie.go:29:	// AllowUnset permits property unset via empty payload.
/root/astrate-mule/pkg/interfaceschema/trie.go:30:	AllowUnset bool
/root/astrate-mule/pkg/interfaceschema/trie_test.go:243:	AllowUnset:        false,
/root/astrate-mule/pkg/interfaceschema/compile.go:79:			AllowUnset:        m.AllowUnset,
/root/astrate-mule/pkg/interfaceschema/parse.go:88:	// AllowUnset permits unsetting the property (properties only).
/root/astrate-mule/pkg/interfaceschema/parse.go:89:	AllowUnset bool
/root/astrate-mule/pkg/interfaceschema/parse.go:145:	AllowUnset              *bool                    `json:"allow_unset,omitempty"`
/root/astrate-mule/pkg/interfaceschema/parse.go:400:		if raw.AllowUnset != nil {
/root/astrate-mule/pkg/interfaceschema/parse.go:401:			m.AllowUnset = *raw.AllowUnset
/root/astrate-mule/pkg/interfaceschema/parse.go:407:	if raw.AllowUnset != nil {
/root/astrate-mule/pkg/interfaceschema/parse.go:408:		return m, nil, fmt.Errorf(`mapping %q: "allow_unset" is only allowed on properties`, m.Endpoint)
/root/astrate-mule/pkg/interfaceschema/compat_test.go:61:// pinnedPropsV10 is a properties interface carrying allow_unset (a
/root/astrate-mule/pkg/interfaceschema/compat_test.go:70:		{"endpoint": "/%{sensor_id}/enabled", "type": "boolean", "allow_unset": true},
/root/astrate-mule/pkg/interfaceschema/compat_test.go:75:// pinnedPropsV11 is the at-immutability acceptance twin of the allow_unset
/root/astrate-mule/pkg/interfaceschema/compat_test.go:76:// flip case: allow_unset stays true across the minor bump.
/root/astrate-mule/pkg/interfaceschema/compat_test.go:84:		{"endpoint": "/%{sensor_id}/enabled", "type": "boolean", "allow_unset": true},
/root/astrate-mule/pkg/interfaceschema/compat_test.go:376:			name:   "allow_unset at immutability accepted",
/root/astrate-mule/pkg/interfaceschema/compat_test.go:381:			name:   "allow_unset mutated",
/root/astrate-mule/pkg/interfaceschema/compat_test.go:388:					{"endpoint": "/%{sensor_id}/enabled", "type": "boolean", "allow_unset": false},
/root/astrate-mule/pkg/interfaceschema/compat_test.go:392:			wantSub:   "changed allow_unset",
/root/astrate-mule/pkg/interfaceschema/testdata/invalid/manifest.json:22:  "21-allow-unset-datastream.json": "\"allow_unset\" is only allowed on properties",
/root/astrate-mule/pkg/interfaceschema/testdata/invalid/21-allow-unset-datastream.json:2:  "interface_name": "com.astrate.test.AllowUnsetDatastream",
/root/astrate-mule/pkg/interfaceschema/testdata/invalid/21-allow-unset-datastream.json:7:  "mappings": [{"endpoint": "/value", "type": "double", "allow_unset": true}]
/root/astrate-mule/pkg/interfaceschema/testdata/valid/com.astrate.test.ServerProperties.json:16:      "allow_unset": true

← Edit internal/realm/detailed_required_encrypted_test.go
Index: /root/astrate-mule/internal/realm/detailed_required_encrypted_test.go
===================================================================
--- /root/astrate-mule/internal/realm/detailed_required_encrypted_test.go
+++ /root/astrate-mule/internal/realm/detailed_required_encrypted_test.go
@@ -26,4 +26,43 @@
 	if string(got) != want {
 		t.Errorf("rendered =\n%s\nwant\n%s", got, want)
 	}
 }
+
+// TestRenderDetailedInterfacePropertiesAllowUnset pins the properties branch
+// of writeDetailedMapping (detailed.go:61-63) container-free: a properties
+// mapping carries exactly endpoint/type/allow_unset and none of the
+// datastream delivery defaults, for both allow_unset values.
+func TestRenderDetailedInterfacePropertiesAllowUnset(t *testing.T) {
+	tests := []struct {
+		name string
+		flag string
+		want string
+	}{
+		{
+			name: "allow_unset true",
+			flag: "true",
+		},
+		{
+			name: "allow_unset false",
+			flag: "false",
+		},
+	}
+	for _, tt := range tests {
+		t.Run(tt.name, func(t *testing.T) {
+			def := []byte(`{"interface_name":"com.ex.Mule.Props","version_major":0,"version_minor":2,` +
+				`"type":"properties","ownership":"server",` +
+				`"mappings":[{"endpoint":"/state","type":"boolean","allow_unset":` + tt.flag + `}]}`)
+			got, err := renderDetailedInterface(def)
+			if err != nil {
+				t.Fatalf("renderDetailedInterface: %v", err)
+			}
+			want := `{"interface_name":"com.ex.Mule.Props","version_major":0,"version_minor":2,` +
+				`"type":"properties","ownership":"server","aggregation":"individual",` +
+				`"mappings":[{"endpoint":"/state","type":"boolean","allow_unset":` + tt.flag + `}]}`
+			t.Logf("want=%s", want)
+			if string(got) != want {
+				t.Errorf("rendered =\n%s\nwant\n%s", got, want)
+			}
+		})
+	}
+}

← Edit internal/realm/detailed_required_encrypted_test.go
Index: /root/astrate-mule/internal/realm/detailed_required_encrypted_test.go
===================================================================
--- /root/astrate-mule/internal/realm/detailed_required_encrypted_test.go
+++ /root/astrate-mule/internal/realm/detailed_required_encrypted_test.go
@@ -58,9 +58,8 @@
 	}
 	want := `{"interface_name":"com.ex.Mule.Props","version_major":0,"version_minor":2,` +
 		`"type":"properties","ownership":"server","aggregation":"individual",` +
 		`"mappings":[{"endpoint":"/state","type":"boolean","allow_unset":` + tt.flag + `}]}`
-	t.Logf("want=%s", want)
 	if string(got) != want {
 		t.Errorf("rendered =\n%s\nwant\n%s", got, want)
 	}
 })

$ go test ./internal/realm/ 2>&1 | tail -20 && gofmt -l internal/realm/detailed_required_encrypted_test.go
ok  	github.com/astrate-platform/astrate/internal/realm	0.030s

$ go vet ./internal/realm/ && go test -race ./internal/realm/ 2>&1 | tail -5
FATAL: ThreadSanitizer: unsupported VMA range
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/realm	0.024s
FAIL

Done: added a container-free table test pinning the properties branch of `writeDetailedMapping` (exact endpoint/type/allow_unset, no datastream defaults) for `allow_unset` true and false.

Files: `internal/realm/detailed_required_encrypted_test.go`

Verified: `go test ./internal/realm/` -> pass; `gofmt -l internal/realm/detailed_required_encrypted_test.go` -> clean; `go vet ./internal/realm/` -> pass. `-race` unavailable here (ThreadSanitizer "unsupported VMA range", 39-bit VMA).

Unsure: nothing. This is a test-only task, so there is no behaviour change for the runner to revert; the new test cannot "fail without the change" because the renderer is already correct.

Follow-ups: none
