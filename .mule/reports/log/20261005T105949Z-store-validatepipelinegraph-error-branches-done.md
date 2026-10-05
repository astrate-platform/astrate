slug: store-validatepipelinegraph-error-branches
verdict: done
at:  73e7387
ran: 2026-10-05T10:59:49Z on DietPi in 316s


> build · big-pickle

→ Read .mule/MULE.md
I'll start by reading the two files.
→ Read .mule/task.md
→ Read internal/store/pipelines_validate_test.go

$ rg -n "validatePipelineGraph" --glob '*.go' -C2 /root/astrate-mule | head -80
/root/astrate-mule/internal/store/pipelines.go-42-}
/root/astrate-mule/internal/store/pipelines.go-43-
/root/astrate-mule/internal/store/pipelines.go:44:// validatePipelineGraph checks that every connection references a declared
/root/astrate-mule/internal/store/pipelines.go-45-// block and that the block graph is acyclic (Kahn's algorithm).
/root/astrate-mule/internal/store/pipelines.go:46:func validatePipelineGraph(definition []byte) error {
/root/astrate-mule/internal/store/pipelines.go-47-	var g pipelineGraph
/root/astrate-mule/internal/store/pipelines.go-48-	if err := json.Unmarshal(definition, &g); err != nil {
--
/root/astrate-mule/internal/store/pipelines.go-102-// error for an unresolved block reference.
/root/astrate-mule/internal/store/pipelines.go-103-func (s *Store) CreatePipeline(ctx context.Context, realmID int16, name string, definition []byte) (*Pipeline, error) {
/root/astrate-mule/internal/store/pipelines.go:104:	if err := validatePipelineGraph(definition); err != nil {
/root/astrate-mule/internal/store/pipelines.go-105-		return nil, err
/root/astrate-mule/internal/store/pipelines.go-106-	}
--
/root/astrate-mule/internal/store/pipelines.go-137-// UpdatePipeline validates and replaces a pipeline's definition.
/root/astrate-mule/internal/store/pipelines.go-138-func (s *Store) UpdatePipeline(ctx context.Context, realmID int16, name string, definition []byte) (*Pipeline, error) {
/root/astrate-mule/internal/store/pipelines.go:139:	if err := validatePipelineGraph(definition); err != nil {
/root/astrate-mule/internal/store/pipelines.go-140-		return nil, err
/root/astrate-mule/internal/store/pipelines.go-141-	}
--
/root/astrate-mule/internal/store/pipelines_validate_test.go-12-		"connections": [{"from": "a", "to": "b"}, {"from": "b", "to": "c"}]
/root/astrate-mule/internal/store/pipelines_validate_test.go-13-	}`)
/root/astrate-mule/internal/store/pipelines_validate_test.go:14:	if err := validatePipelineGraph(def); err != nil {
/root/astrate-mule/internal/store/pipelines_validate_test.go:15:		t.Errorf("validatePipelineGraph(acyclic) = %v, want nil", err)
/root/astrate-mule/internal/store/pipelines_validate_test.go-16-	}
/root/astrate-mule/internal/store/pipelines_validate_test.go-17-}
--
/root/astrate-mule/internal/store/pipelines_validate_test.go-22-		"connections": [{"from": "a", "to": "b"}, {"from": "b", "to": "c"}, {"from": "c", "to": "a"}]
/root/astrate-mule/internal/store/pipelines_validate_test.go-23-	}`)
/root/astrate-mule/internal/store/pipelines_validate_test.go:24:	err := validatePipelineGraph(def)
/root/astrate-mule/internal/store/pipelines_validate_test.go-25-	if !errors.Is(err, ErrPipelineCyclic) {
/root/astrate-mule/internal/store/pipelines_validate_test.go:26:		t.Errorf("validatePipelineGraph(cyclic) = %v, want ErrPipelineCyclic", err)
/root/astrate-mule/internal/store/pipelines_validate_test.go-27-	}
/root/astrate-mule/internal/store/pipelines_validate_test.go-28-}
--
/root/astrate-mule/internal/store/pipelines_validate_test.go-33-		"connections": [{"from": "a", "to": "ghost"}]
/root/astrate-mule/internal/store/pipelines_validate_test.go-34-	}`)
/root/astrate-mule/internal/store/pipelines_validate_test.go:35:	err := validatePipelineGraph(def)
/root/astrate-mule/internal/store/pipelines_validate_test.go-36-	if err == nil || errors.Is(err, ErrPipelineCyclic) {
/root/astrate-mule/internal/store/pipelines_validate_test.go:37:		t.Errorf("validatePipelineGraph(invalid ref) = %v, want a non-cyclic error", err)
/root/astrate-mule/internal/store/pipelines_validate_test.go-38-	}
/root/astrate-mule/internal/store/pipelines_validate_test.go-39-}
--
/root/astrate-mule/internal/store/pipelines_validate_test.go-44-		"connections": []
/root/astrate-mule/internal/store/pipelines_validate_test.go-45-	}`)
/root/astrate-mule/internal/store/pipelines_validate_test.go:46:	err := validatePipelineGraph(def)
/root/astrate-mule/internal/store/pipelines_validate_test.go-47-	if err == nil || errors.Is(err, ErrPipelineCyclic) {
/root/astrate-mule/internal/store/pipelines_validate_test.go:48:		t.Errorf("validatePipelineGraph(duplicate name) = %v, want a non-cyclic error", err)
/root/astrate-mule/internal/store/pipelines_validate_test.go-49-	}
/root/astrate-mule/internal/store/pipelines_validate_test.go-50-}
--
/root/astrate-mule/internal/store/pipelines_validate_test.go-55-		"connections": []
/root/astrate-mule/internal/store/pipelines_validate_test.go-56-	}`)
/root/astrate-mule/internal/store/pipelines_validate_test.go:57:	err := validatePipelineGraph(def)
/root/astrate-mule/internal/store/pipelines_validate_test.go-58-	if err == nil || errors.Is(err, ErrPipelineCyclic) {
/root/astrate-mule/internal/store/pipelines_validate_test.go:59:		t.Errorf("validatePipelineGraph(no blocks) = %v, want a non-cyclic error", err)
/root/astrate-mule/internal/store/pipelines_validate_test.go-60-	}
/root/astrate-mule/internal/store/pipelines_validate_test.go-61-	if !strings.Contains(err.Error(), "pipeline has no blocks") {
/root/astrate-mule/internal/store/pipelines_validate_test.go:62:		t.Errorf("validatePipelineGraph(no blocks) = %v, want the no-blocks error", err)
/root/astrate-mule/internal/store/pipelines_validate_test.go-63-	}
/root/astrate-mule/internal/store/pipelines_validate_test.go-64-}
--
/root/astrate-mule/internal/store/pipelines_validate_test.go-69-		"connections": []
/root/astrate-mule/internal/store/pipelines_validate_test.go-70-	}`)
/root/astrate-mule/internal/store/pipelines_validate_test.go:71:	err := validatePipelineGraph(def)
/root/astrate-mule/internal/store/pipelines_validate_test.go-72-	if err == nil || errors.Is(err, ErrPipelineCyclic) {
/root/astrate-mule/internal/store/pipelines_validate_test.go:73:		t.Errorf("validatePipelineGraph(empty name) = %v, want a non-cyclic error", err)
/root/astrate-mule/internal/store/pipelines_validate_test.go-74-	}
/root/astrate-mule/internal/store/pipelines_validate_test.go-75-	if !strings.Contains(err.Error(), "block with empty name") {
/root/astrate-mule/internal/store/pipelines_validate_test.go:76:		t.Errorf("validatePipelineGraph(empty name) = %v, want the empty-name error", err)
/root/astrate-mule/internal/store/pipelines_validate_test.go-77-	}
/root/astrate-mule/internal/store/pipelines_validate_test.go-78-}

→ Read internal/store/pipelines.go [limit=100]

$ go test ./internal/store/ -run 'TestValidatePipelineGraph' -coverprofile=/tmp/opencode/cov.out -covermode=count 2>&1 | tail -5; go tool cover -func=/tmp/opencode/cov.out | rg 'pipelines.go'
ok  	github.com/astrate-platform/astrate/internal/store	0.016s	coverage: 2.9% of statements
github.com/astrate-platform/astrate/internal/store/pipelines.go:46:	validatePipelineGraph		97.1%
github.com/astrate-platform/astrate/internal/store/pipelines.go:103:	CreatePipeline			0.0%
github.com/astrate-platform/astrate/internal/store/pipelines.go:122:	GetPipeline			0.0%
github.com/astrate-platform/astrate/internal/store/pipelines.go:138:	UpdatePipeline			0.0%
github.com/astrate-platform/astrate/internal/store/pipelines.go:158:	DeletePipeline			0.0%
github.com/astrate-platform/astrate/internal/store/pipelines.go:171:	ListPipelines			0.0%

$ rg 'pipelines.go' /tmp/opencode/cov.out | rg ':(4[6-9]|[5-9][0-9])\.'
github.com/astrate-platform/astrate/internal/store/pipelines.go:46.53,48.55 2 6
github.com/astrate-platform/astrate/internal/store/pipelines.go:48.55,50.3 1 0
github.com/astrate-platform/astrate/internal/store/pipelines.go:51.2,51.24 1 6
github.com/astrate-platform/astrate/internal/store/pipelines.go:51.24,53.3 1 1
github.com/astrate-platform/astrate/internal/store/pipelines.go:55.2,56.29 2 5
github.com/astrate-platform/astrate/internal/store/pipelines.go:56.29,57.19 1 11
github.com/astrate-platform/astrate/internal/store/pipelines.go:57.19,59.4 1 1
github.com/astrate-platform/astrate/internal/store/pipelines.go:60.3,60.20 1 10
github.com/astrate-platform/astrate/internal/store/pipelines.go:60.20,62.4 1 1
github.com/astrate-platform/astrate/internal/store/pipelines.go:63.3,63.23 1 9
github.com/astrate-platform/astrate/internal/store/pipelines.go:66.2,68.34 3 3
github.com/astrate-platform/astrate/internal/store/pipelines.go:68.34,69.37 1 6
github.com/astrate-platform/astrate/internal/store/pipelines.go:69.37,71.4 1 1
github.com/astrate-platform/astrate/internal/store/pipelines.go:72.3,73.42 2 5
github.com/astrate-platform/astrate/internal/store/pipelines.go:76.2,77.26 2 2
github.com/astrate-platform/astrate/internal/store/pipelines.go:77.26,78.23 1 6
github.com/astrate-platform/astrate/internal/store/pipelines.go:78.23,80.4 1 1
github.com/astrate-platform/astrate/internal/store/pipelines.go:82.2,83.21 2 2
github.com/astrate-platform/astrate/internal/store/pipelines.go:83.21,87.31 4 3
github.com/astrate-platform/astrate/internal/store/pipelines.go:87.31,89.24 2 2
github.com/astrate-platform/astrate/internal/store/pipelines.go:89.24,91.5 1 2
github.com/astrate-platform/astrate/internal/store/pipelines.go:94.2,94.27 1 2
github.com/astrate-platform/astrate/internal/store/pipelines.go:94.27,96.3 1 1
github.com/astrate-platform/astrate/internal/store/pipelines.go:97.2,97.12 1 1

$ rg -n -i 'validatepipelinegraph|pipelines_validate' /root/astrate-mule/.mule/ 2>/dev/null | head -20
/root/astrate-mule/.mule/todo.md:204:- [x] store-pipelines-empty-name-zero-blocks-test [auto]: extend internal/store/pipelines_validate_test.go with the two untested error branches of `validatePipelineGraph` (internal/store/pipelines.go:51-64) — a definition with zero blocks and one whose block has an empty name must both be rejected. Pure helper, no DB; the existing suite covers only acyclic/cyclic/unknown-ref/duplicate-name.
/root/astrate-mule/.mule/todo.md:276:- [ ] store-validatepipelinegraph-error-branches: Add tests covering both error branches of validatePipelineGraph. File: internal/store/pipelines_validate_test.go
/root/astrate-mule/.mule/reviews/store-2026-09-21.md:28:  (88), `retention.go` (54). Tests: `pipelines_validate_test.go`,
/root/astrate-mule/.mule/reviews/store-2026-09-21.md:60:### 2. `validatePipelineGraph`: two error branches untested
/root/astrate-mule/.mule/reviews/store-2026-09-21.md:63:and (b) a block with an empty name. `pipelines_validate_test.go` covers
/root/astrate-mule/.mule/reviews/store-2026-10-04.md:20:  pipelines_validate_test.go; properties.go (183), properties_test.go;
/root/astrate-mule/.mule/reviews/store-2026-10-04.md:38:### 4. validatePipelineGraph: two error branches untested
/root/astrate-mule/.mule/reviews/store-2026-10-04.md:39:pipelines.go (validatePipelineGraph) has branches that reject certain graph states; pipelines_validate_test.go exists but does not assert both error conditions explicitly. Pinning these makes the validation contract robust.
/root/astrate-mule/.mule/reviews/store-2026-09-08.md:28:  `pipelines_validate_test.go`, `retention_test.go`
/root/astrate-mule/.mule/reviews/store-2026-09-08.md:72:  `UpdateFlowRuntime`, `validatePipelineGraph`.** All already have direct
/root/astrate-mule/.mule/tasks/issue-24.md:66:// ValidatePipelineGraph performs in-memory graph validation without DB access.
/root/astrate-mule/.mule/tasks/issue-24.md:67:func ValidatePipelineGraph(p *flow.Pipeline) error {
/root/astrate-mule/.mule/tasks/issue-24.md:143:func TestValidatePipelineGraph_Acyclic(t *testing.T) {
/root/astrate-mule/.mule/tasks/issue-24.md:155:	if err := store.ValidatePipelineGraph(p); err != nil {
/root/astrate-mule/.mule/tasks/issue-24.md:160:func TestValidatePipelineGraph_DirectCycle(t *testing.T) {
/root/astrate-mule/.mule/tasks/issue-24.md:172:	err := store.ValidatePipelineGraph(p)
/root/astrate-mule/.mule/tasks/issue-24.md:178:func TestValidatePipelineGraph_IndirectCycle(t *testing.T) {
/root/astrate-mule/.mule/tasks/issue-24.md:191:	err := store.ValidatePipelineGraph(p)
/root/astrate-mule/.mule/tasks/issue-24.md:197:func TestValidatePipelineGraph_MissingBlockReference(t *testing.T) {
/root/astrate-mule/.mule/tasks/issue-24.md:206:	err := store.ValidatePipelineGraph(p)

$ sed -n '30,50p' /root/astrate-mule/.mule/reviews/store-2026-10-04.md
Untested: devices_test.Lifecycle covers re-register before credentials (rotates secret) but not the inhibit→re-register case. Uninhibit already restores 'registered' if never confirmed (SetDeviceInhibited logic), but the re-registration path overwrites status to 'registered' regardless of prior inhibited state — that’s the bug.

### 2. AliasValuesTaken lacks full branch coverage
AliasValuesTaken (devices.go:343-361) excludes self (`d.id <> $2`) and scans for any other device whose aliases jsonb contains values in the input. It returns true on first match, false on ErrNoRows. The three branches worth pinning are: (a) no other device has any of the values → false; (b) another device has a value → true; (c) the same device has the value (self-match excluded) → must not count as taken (false). devices_test covers basic cases but not (c) explicitly as a store-level assertion isolating the self-exclusion.

### 3. LatestIndividual has no store-level test asserting empty-series yields ErrNotFound
datastreams.go:247-274 implements LatestIndividual (ordered DESC by ts with limit 1). Engine tests use a fake store; the real store path’s ErrNoRows→ErrNotFound conversion and value retrieval (both columns) aren’t directly asserted. Worth pinning: (a) empty series → ErrNotFound; (b) two samples on same series, returns newest; (c) correct value column populated. The ErrNotFound case is particularly valuable (used by trigger previous-value lookups).

### 4. validatePipelineGraph: two error branches untested
pipelines.go (validatePipelineGraph) has branches that reject certain graph states; pipelines_validate_test.go exists but does not assert both error conditions explicitly. Pinning these makes the validation contract robust.

### 5. Properties purge/ownership semantics edge cases could be tightened
properties_test covers PurgeComplement and server-owned preservation, but a small additional assertion (e.g. preserve server-owned across purge of non-complement set, or mixed update paths) would lock behaviour; not essential. Keep to the strongest four above.

## What I decided NOT to propose, and why

- patchDeviceJSONB no-op when patch empty (devices.go:377-404): harmless, callers don’t send empty patches, no bug.
- Interface notify layering: NotifyInterfacesChanged issued by realm service, not store — deliberate.
- Downsample/LTTB parameter indexing: verified in multi-filter cases; tests already exercise.
- ApplyGlobalRetention/EnforceRealmRetentionCeilings/UpdateFlowRuntime: already tested.
- Cosmetic renames/reformats: linter domain.

$ rg -l 'require\.|assert\.' internal/store/ | head; rg -n 'testify' go.mod; ls internal/store/*_test.go
84:	github.com/stretchr/testify v1.11.1 // indirect
internal/store/compression_test.go
internal/store/crypto_test.go
internal/store/datastreams_test.go
internal/store/devices_test.go
internal/store/flows_test.go
internal/store/groups_test.go
internal/store/interfaces_test.go
internal/store/notify_test.go
internal/store/pipelines_test.go
internal/store/pipelines_validate_test.go
internal/store/policies_test.go
internal/store/properties_test.go
internal/store/realms_test.go
internal/store/retention_test.go
internal/store/store_test.go
internal/store/triggers_test.go
internal/store/userblocks_test.go

← Edit internal/store/pipelines_validate_test.go
Index: /root/astrate-mule/internal/store/pipelines_validate_test.go
===================================================================
--- /root/astrate-mule/internal/store/pipelines_validate_test.go
+++ /root/astrate-mule/internal/store/pipelines_validate_test.go
@@ -26,28 +26,70 @@
 		t.Errorf("validatePipelineGraph(cyclic) = %v, want ErrPipelineCyclic", err)
 	}
 }
 
-func TestValidatePipelineGraph_InvalidBlockRef(t *testing.T) {
-	def := []byte(`{
-		"blocks": [{"name": "a"}, {"name": "b"}],
-		"connections": [{"from": "a", "to": "ghost"}]
-	}`)
+func TestValidatePipelineGraph_MalformedJSON(t *testing.T) {
+	def := []byte(`{"blocks": [{"name": "a"}`)
 	err := validatePipelineGraph(def)
 	if err == nil || errors.Is(err, ErrPipelineCyclic) {
-		t.Errorf("validatePipelineGraph(invalid ref) = %v, want a non-cyclic error", err)
+		t.Fatalf("validatePipelineGraph(malformed) = %v, want a non-cyclic error", err)
 	}
+	if !strings.Contains(err.Error(), "does not parse") {
+		t.Errorf("validatePipelineGraph(malformed) = %v, want the parse error", err)
+	}
+	var syntaxErr *json.SyntaxError
+	if !errors.As(err, &syntaxErr) {
+		t.Errorf("validatePipelineGraph(malformed) = %v, want it to wrap *json.SyntaxError", err)
+	}
+	// The parse must fail before the graph is inspected, so the zero-blocks
+	// check cannot be what rejects this definition.
+	if strings.Contains(err.Error(), "no blocks") {
+		t.Errorf("validatePipelineGraph(malformed) = %v, want the parse error to win over the no-blocks error", err)
+	}
 }
 
+func TestValidatePipelineGraph_InvalidBlockRef(t *testing.T) {
+	tests := []struct {
+		name    string
+		def     string
+		wantMsg string
+	}{
+		{
+			name:    "unknown destination",
+			def:     `{"blocks": [{"name": "a"}, {"name": "b"}], "connections": [{"from": "a", "to": "ghost"}]}`,
+			wantMsg: `unknown block "a" -> "ghost"`,
+		},
+		{
+			name:    "unknown source",
+			def:     `{"blocks": [{"name": "a"}, {"name": "b"}], "connections": [{"from": "ghost", "to": "b"}]}`,
+			wantMsg: `unknown block "ghost" -> "b"`,
+		},
+	}
+	for _, tt := range tests {
+		t.Run(tt.name, func(t *testing.T) {
+			err := validatePipelineGraph([]byte(tt.def))
+			if err == nil || errors.Is(err, ErrPipelineCyclic) {
+				t.Fatalf("validatePipelineGraph(invalid ref) = %v, want a non-cyclic error", err)
+			}
+			if !strings.Contains(err.Error(), tt.wantMsg) {
+				t.Errorf("validatePipelineGraph(invalid ref) = %v, want it to mention %q", err, tt.wantMsg)
+			}
+		})
+	}
+}
+
 func TestValidatePipelineGraph_DuplicateBlockName(t *testing.T) {
 	def := []byte(`{
 		"blocks": [{"name": "a"}, {"name": "a"}],
 		"connections": []
 	}`)
 	err := validatePipelineGraph(def)
 	if err == nil || errors.Is(err, ErrPipelineCyclic) {
-		t.Errorf("validatePipelineGraph(duplicate name) = %v, want a non-cyclic error", err)
+		t.Fatalf("validatePipelineGraph(duplicate name) = %v, want a non-cyclic error", err)
 	}
+	if !strings.Contains(err.Error(), `duplicate block name "a"`) {
+		t.Errorf("validatePipelineGraph(duplicate name) = %v, want the duplicate-name error", err)
+	}
 }
 
 func TestValidatePipelineGraph_NoBlocks(t *testing.T) {
 	def := []byte(`{

Now I'll add the tests. Three gaps: the parse-error branch (genuinely uncovered), and the two branches (unknown-ref, duplicate-name) that only assert "some error".
← Edit internal/store/pipelines_validate_test.go
Index: /root/astrate-mule/internal/store/pipelines_validate_test.go
===================================================================
--- /root/astrate-mule/internal/store/pipelines_validate_test.go
+++ /root/astrate-mule/internal/store/pipelines_validate_test.go
@@ -1,7 +1,8 @@
 package store
 
 import (
+	"encoding/json"
 	"errors"
 	"strings"
 	"testing"
 )

$ gofmt -l . && go vet ./internal/store/ && go test ./internal/store/ -run 'TestValidatePipelineGraph' -v -coverprofile=/tmp/opencode/cov2.out 2>&1 | tail -30
=== RUN   TestValidatePipelineGraph_Acyclic
--- PASS: TestValidatePipelineGraph_Acyclic (0.00s)
=== RUN   TestValidatePipelineGraph_Cyclic
--- PASS: TestValidatePipelineGraph_Cyclic (0.00s)
=== RUN   TestValidatePipelineGraph_MalformedJSON
--- PASS: TestValidatePipelineGraph_MalformedJSON (0.00s)
=== RUN   TestValidatePipelineGraph_InvalidBlockRef
=== RUN   TestValidatePipelineGraph_InvalidBlockRef/unknown_destination
=== RUN   TestValidatePipelineGraph_InvalidBlockRef/unknown_source
--- PASS: TestValidatePipelineGraph_InvalidBlockRef (0.00s)
    --- PASS: TestValidatePipelineGraph_InvalidBlockRef/unknown_destination (0.00s)
    --- PASS: TestValidatePipelineGraph_InvalidBlockRef/unknown_source (0.00s)
=== RUN   TestValidatePipelineGraph_DuplicateBlockName
--- PASS: TestValidatePipelineGraph_DuplicateBlockName (0.00s)
=== RUN   TestValidatePipelineGraph_NoBlocks
--- PASS: TestValidatePipelineGraph_NoBlocks (0.00s)
=== RUN   TestValidatePipelineGraph_EmptyBlockName
--- PASS: TestValidatePipelineGraph_EmptyBlockName (0.00s)
PASS
coverage: 3.0% of statements
ok  	github.com/astrate-platform/astrate/internal/store	0.017s	coverage: 3.0% of statements

$ go tool cover -func=/tmp/opencode/cov2.out | rg 'validatePipelineGraph'
github.com/astrate-platform/astrate/internal/store/pipelines.go:46:	validatePipelineGraph		100.0%

$ cp internal/store/pipelines.go /tmp/opencode/pipelines.go.bak && python3 - <<'EOF'
import subprocess, re, shutil
orig = open('/tmp/opencode/pipelines.go.bak').read()
muts = {
 'drop parse branch': ('\tif err := json.Unmarshal(definition, &g); err != nil {\n\t\treturn fmt.Errorf("store: pipeline definition does not parse: %w", err)\n\t}\n', '\t_ = json.Unmarshal(definition, &g)\n'),
 'drop duplicate-name check': ('\t\tif names[b.Name] {\n\t\t\treturn fmt.Errorf("store: pipeline has duplicate block name %q", b.Name)\n\t\t}\n', ''),
 'drop unknown-ref check': ('\t\tif !names[c.From] || !names[c.To] {\n\t\t\treturn fmt.Errorf("store: pipeline connection references unknown block %q -> %q", c.From, c.To)\n\t\t}\n', ''),
}
for name,(old,new) in muts.items():
    assert old in orig, name
    open('internal/store/pipelines.go','w').write(orig.replace(old,new,1))
    r = subprocess.run(['go','test','./internal/store/','-run','TestValidatePipelineGraph'],capture_output=True,text=True)
    print('=== mutation:',name,'-> exit',r.returncode)
    print('\n'.join(l for l in (r.stdout+r.stderr).splitlines() if 'FAIL' in l or '---' in l or 'cannot' in l)[:900])
shutil.copy('/tmp/opencode/pipelines.go.bak','internal/store/pipelines.go')
EOF
=== mutation: drop parse branch -> exit 1
--- FAIL: TestValidatePipelineGraph_MalformedJSON (0.00s)
FAIL
FAIL	github.com/astrate-platform/astrate/internal/store	0.012s
FAIL
=== mutation: drop duplicate-name check -> exit 1
--- FAIL: TestValidatePipelineGraph_DuplicateBlockName (0.00s)
FAIL
FAIL	github.com/astrate-platform/astrate/internal/store	0.013s
FAIL
=== mutation: drop unknown-ref check -> exit 1
--- FAIL: TestValidatePipelineGraph_InvalidBlockRef (0.00s)
    --- FAIL: TestValidatePipelineGraph_InvalidBlockRef/unknown_destination (0.00s)
    --- FAIL: TestValidatePipelineGraph_InvalidBlockRef/unknown_source (0.00s)
FAIL
FAIL	github.com/astrate-platform/astrate/internal/store	0.011s
FAIL

$ diff /tmp/opencode/pipelines.go.bak internal/store/pipelines.go && echo "pipelines.go restored clean" && git status --porcelain
pipelines.go restored clean
 M internal/store/pipelines_validate_test.go

$ gofmt -l . ; go vet ./... && go test ./... 2>&1 | rg -v '^ok|no test files' | head -20; echo "gate exit: $?"
gate exit: 0

$ go test -race ./internal/store/ -run TestValidatePipelineGraph -count=1 2>&1 | tail -3; rm -f /tmp/opencode/cov.out /tmp/opencode/cov2.out /tmp/opencode/pipelines.go.bak
FATAL: Found 39 - Supported 48
FAIL	github.com/astrate-platform/astrate/internal/store	0.027s
FAIL

$ git diff --stat && go test ./internal/store/ -count=1 2>&1 | tail -2
internal/store/pipelines_validate_test.go | 57 +++++++++++++++++++++++++++----
 1 file changed, 50 insertions(+), 7 deletions(-)
ok  	github.com/astrate-platform/astrate/internal/store	0.011s

Race detector unavailable here (39-bit VMA kernel, as MULE.md predicts), so the gate was run without `-race`.

```
Done: Added the missing validatePipelineGraph parse-error test and pinned the exact messages for the unknown-block-ref and duplicate-block-name branches; the helper is now at 100% statement coverage and each new assertion was mutation-checked.
Files: internal/store/pipelines_validate_test.go
Verified: gofmt -l . (clean) -> pass; go vet ./... && go test ./... -> pass; go test ./internal/store/ -run TestValidatePipelineGraph -coverprofile -> validatePipelineGraph 97.1% -> 100.0%; three mutations of internal/store/pipelines.go (drop parse branch / drop duplicate-name check / drop unknown-ref check) each break exactly the matching new test -> pass; git diff confirms internal/store/pipelines.go was restored byte-identical.
Unsure: The task says "both error branches" without naming them. Only one branch (JSON unmarshal failure, pipelines.go:48) was genuinely uncovered — zero-blocks and empty-name were already added by the earlier store-pipelines-empty-name-zero-blocks-test task (.mule/todo.md:204). I read the review that spawned the line (.mule/reviews/store-2026-10-04.md item 4) as pointing at the two branches the suite exercised only as "some non-nil, non-cyclic error", so I pinned those two messages as well as adding the parse-error test. If "both" meant something else, the surplus is the strengthened InvalidBlockRef/DuplicateBlockName assertions, which are harmless and strictly stronger. Also note: the review claims coverage was missing at pipelines.go:51-64, which is stale — those lines are already tested.
Follow-ups: none
```
