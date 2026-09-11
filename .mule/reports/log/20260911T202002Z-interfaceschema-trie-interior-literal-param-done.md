slug: interfaceschema-trie-interior-literal-param
verdict: done
at:  705795f
ran: 2026-09-11T20:20:02Z on DietPi in 121s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read pkg/interfaceschema/trie.go
Now let me read the existing test file to match its style.
→ Read pkg/interfaceschema/trie_test.go

$ go version && ls
go version go1.26.7 linux/arm64
AGENTS.md
Dockerfile
LICENSE
Makefile
README.md
bench
cmd
deploy
docker-compose.yml
docs
examples
go.mod
go.sum
internal
migrations
opencode.json
pkg
site-dist
test
tools

Let me verify the package builds and check the current behavior with a quick probe before writing the test.
$ cat > /tmp/ttprobe.go <<'EOF'
package main

import (
	"fmt"

	"astrate/pkg/interfaceschema"
)

func main() {
	trie := interfaceschema.NewEndpointTrie()
	for i, ep := range []string{"/a/b/c", "/a/%{p}"} {
		if err := trie.Add(ep, &interfaceschema.CompiledMapping{EndpointID: int64(i + 1)}); err != nil {
			panic(err)
		}
	}
	for _, p := range []string{"/a/b", "/a/b/c", "/a/b/c/d"} {
		m, ok := trie.Match(p)
		fmt.Printf("%-10s ok=%v id=%d\n", p, ok, func() int64 {
			if m == nil {
				return -1
			}
			return m.EndpointID
		}())
	}
}
EOF
go run /tmp/ttprobe.go
/tmp/ttprobe.go:6:2: package astrate/pkg/interfaceschema is not in std (/usr/lib/go-1.26/src/astrate/pkg/interfaceschema)

$ head -1 go.mod
module github.com/astrate-platform/astrate

$ sed -i 's#astrate/pkg/interfaceschema#github.com/astrate-platform/astrate/pkg/interfaceschema#' /tmp/ttprobe.go && go run /tmp/ttprobe.go
/a/b       ok=true id=2
/a/b/c     ok=true id=1
/a/b/c/d   ok=false id=-1

← Edit pkg/interfaceschema/trie_test.go
Index: /root/astrate-mule/pkg/interfaceschema/trie_test.go
===================================================================
--- /root/astrate-mule/pkg/interfaceschema/trie_test.go
+++ /root/astrate-mule/pkg/interfaceschema/trie_test.go
@@ -110,8 +110,37 @@
 		t.Errorf("Match(/a/b/deep) = EndpointID %d, want %d", got.EndpointID, want.EndpointID)
 	}
 }
 
+// TestInteriorLiteralDoesNotShadowParam pins the rule that an interior-only
+// literal child — one with no mapping at the matching depth — does not shadow
+// its parametric sibling (trie.go matchNode): with only "/a/b/c" (leaf at
+// depth 3) and "/a/%{p}" declared, Match("/a/b") falls through the
+// mapping-less interior node "b" to the parametric branch, while
+// Match("/a/b/c") still resolves literals first. Statement coverage cannot
+// see this: matchNode reports 100% in the broken variant too, because every
+// statement executes while the value only differs.
+func TestInteriorLiteralDoesNotShadowParam(t *testing.T) {
+	trie := NewEndpointTrie()
+	deep := &CompiledMapping{EndpointID: 1, ValueType: Double}
+	param := &CompiledMapping{EndpointID: 2, ValueType: Double}
+	if err := trie.Add("/a/b/c", deep); err != nil {
+		t.Fatalf("Add(/a/b/c): %v", err)
+	}
+	if err := trie.Add("/a/%{p}", param); err != nil {
+		t.Fatalf("Add(/a/%%{p}): %v", err)
+	}
+
+	if got, ok := trie.Match("/a/b"); !ok || got != param {
+		t.Fatalf("Match(/a/b) = %+v (ok=%v), want the /a/%%{p} mapping (EndpointID %d)",
+			got, ok, param.EndpointID)
+	}
+	if got, ok := trie.Match("/a/b/c"); !ok || got != deep {
+		t.Fatalf("Match(/a/b/c) = %+v (ok=%v), want the /a/b/c mapping (EndpointID %d)",
+			got, ok, deep.EndpointID)
+	}
+}
+
 func TestAddRejections(t *testing.T) {
 	trie := NewEndpointTrie()
 	m := &CompiledMapping{ValueType: Double}
 	if err := trie.Add("/a/%{p}/x", m); err != nil {

Current code behaves as intended. Now I'll pin it with a test.
$ gofmt -l . && go vet ./pkg/interfaceschema/ && go test ./pkg/interfaceschema/ -run 'TestInteriorLiteralDoesNotShadowParam|TestMatch' -v 2>&1 | tail -30
--- PASS: TestMatchTable/parametric_fallback (0.00s)
    --- PASS: TestMatchTable/parametric_numeric_value (0.00s)
    --- PASS: TestMatchTable/deep_under_parametric (0.00s)
    --- PASS: TestMatchTable/root_parametric (0.00s)
    --- PASS: TestMatchTable/mixed_literal_parametric_literal (0.00s)
    --- PASS: TestMatchTable/placeholder_at_max_length (0.00s)
    --- PASS: TestMatchTable/placeholder_value_with_dots_and_colons (0.00s)
    --- PASS: TestMatchTable/miss_unknown_root (0.00s)
    --- PASS: TestMatchTable/miss_too_deep (0.00s)
    --- PASS: TestMatchTable/miss_too_shallow (0.00s)
    --- PASS: TestMatchTable/miss_interior_without_leaf (0.00s)
    --- PASS: TestMatchTable/miss_wrong_leaf (0.00s)
    --- PASS: TestMatchTable/empty_path (0.00s)
    --- PASS: TestMatchTable/bare_slash (0.00s)
    --- PASS: TestMatchTable/not_rooted (0.00s)
    --- PASS: TestMatchTable/empty_middle_segment (0.00s)
    --- PASS: TestMatchTable/trailing_slash (0.00s)
    --- PASS: TestMatchTable/overlong_placeholder_value (0.00s)
    --- PASS: TestMatchTable/plus_in_placeholder_value (0.00s)
    --- PASS: TestMatchTable/hash_in_placeholder_value (0.00s)
    --- PASS: TestMatchTable/wildcard_alone (0.00s)
    --- PASS: TestMatchTable/hash_under_root_param (0.00s)
=== RUN   TestMatchBacktracking
--- PASS: TestMatchBacktracking (0.00s)
=== RUN   TestInteriorLiteralDoesNotShadowParam
--- PASS: TestInteriorLiteralDoesNotShadowParam (0.00s)
=== RUN   TestMatchZeroAllocs
--- PASS: TestMatchZeroAllocs (0.01s)
PASS
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	0.033s

Full gate check on the touched package plus the whole module in one go.
$ go test ./pkg/interfaceschema/ && go vet ./... 2>&1 | tail -5
ok  	github.com/astrate-platform/astrate/pkg/interfaceschema	0.139s

Done: pinned the interior-only-literal-does-not-shadow-param rule in the endpoints trie with a purpose-built "/a/b/c" + "/a/%{p}" fixture asserting Match("/a/b")→parametric and Match("/a/b/c")→literal.
Files: pkg/interfaceschema/trie_test.go
Verified: go vet ./... -> pass; go test ./pkg/interfaceschema/ -> pass; gofmt -l . -> clean (gofmt -l reported nothing)
Unsure: nothing
Follow-ups: none
