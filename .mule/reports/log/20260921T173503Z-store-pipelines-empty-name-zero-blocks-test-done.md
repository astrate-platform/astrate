slug: store-pipelines-empty-name-zero-blocks-test
verdict: done
at:  82fc0a4
ran: 2026-09-21T17:35:03Z on DietPi in 91s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/store/pipelines.go [offset=1, limit=100]
→ Read internal/store/pipelines_validate_test.go

← Edit internal/store/pipelines_validate_test.go
Index: /root/astrate-mule/internal/store/pipelines_validate_test.go
===================================================================
--- /root/astrate-mule/internal/store/pipelines_validate_test.go
+++ /root/astrate-mule/internal/store/pipelines_validate_test.go
@@ -1,8 +1,9 @@
 package store
 
 import (
 	"errors"
+	"strings"
 	"testing"
 )
 
 func TestValidatePipelineGraph_Acyclic(t *testing.T) {

Read the task, the helper (internal/store/pipelines.go:46-98), and the existing suite. The two untested branches are the zero-blocks check (lines 51-53) and the empty-name check (lines 57-59).
← Edit internal/store/pipelines_validate_test.go
Index: /root/astrate-mule/internal/store/pipelines_validate_test.go
===================================================================
--- /root/astrate-mule/internal/store/pipelines_validate_test.go
+++ /root/astrate-mule/internal/store/pipelines_validate_test.go
@@ -47,4 +47,32 @@
 	if err == nil || errors.Is(err, ErrPipelineCyclic) {
 		t.Errorf("validatePipelineGraph(duplicate name) = %v, want a non-cyclic error", err)
 	}
 }
+
+func TestValidatePipelineGraph_NoBlocks(t *testing.T) {
+	def := []byte(`{
+		"blocks": [],
+		"connections": []
+	}`)
+	err := validatePipelineGraph(def)
+	if err == nil || errors.Is(err, ErrPipelineCyclic) {
+		t.Errorf("validatePipelineGraph(no blocks) = %v, want a non-cyclic error", err)
+	}
+	if !strings.Contains(err.Error(), "pipeline has no blocks") {
+		t.Errorf("validatePipelineGraph(no blocks) = %v, want the no-blocks error", err)
+	}
+}
+
+func TestValidatePipelineGraph_EmptyBlockName(t *testing.T) {
+	def := []byte(`{
+		"blocks": [{"name": ""}],
+		"connections": []
+	}`)
+	err := validatePipelineGraph(def)
+	if err == nil || errors.Is(err, ErrPipelineCyclic) {
+		t.Errorf("validatePipelineGraph(empty name) = %v, want a non-cyclic error", err)
+	}
+	if !strings.Contains(err.Error(), "block with empty name") {
+		t.Errorf("validatePipelineGraph(empty name) = %v, want the empty-name error", err)
+	}
+}

$ go test ./internal/store/ && gofmt -l .
ok  	github.com/astrate-platform/astrate/internal/store	0.022s

$ go vet ./internal/store/
(no output)

```
Done: added NoBlocks and EmptyBlockName tests for validatePipelineGraph's two untested error branches in internal/store/pipelines_validate_test.go
Files: internal/store/pipelines_validate_test.go
Verified: go test ./internal/store/ (pass), go vet ./internal/store/ (pass), gofmt -l . (clean) -> pass
Unsure: nothing
Follow-ups: none
```

Both tests assert the specific error branch (message substring), so they only pass if the zero-blocks/empty-name rejection actually fires.
