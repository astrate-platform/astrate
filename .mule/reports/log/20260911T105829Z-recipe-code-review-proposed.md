slug: recipe-code-review
verdict: proposed
at:  321fc79
ran: 2026-09-11T10:58:29Z on DietPi in 190s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/code-review.md
→ Read .mule/reviews
→ Read .mule/log.md

$ ls /root/astrate-mule/internal/ /root/astrate-mule/pkg/ 2>/dev/null
/root/astrate-mule/internal/:
appengine
auth
broker
config
engine
flow
flowapi
housekeeping
httpx
observability
pairing
realm
store
swagger
testutil

/root/astrate-mule/pkg/:
astarteapi
deviceid
interfaceschema
payload

$ rg -c '' --glob 'internal/housekeeping/**/*.go' | sort -t: -k2 -rn | head -20
internal/housekeeping/http_test.go:511
internal/housekeeping/service.go:280
internal/housekeeping/http.go:259
internal/housekeeping/service_test.go:224
internal/housekeeping/async_operation_test.go:37

$ rg -c '' --glob 'internal/swagger/**/*.go' | sort -t: -k2 -rn | head -20
internal/swagger/swagger_test.go:128
internal/swagger/swagger.go:43

$ rg -c '' --glob 'internal/auth/**/*.go' | sort -t: -k2 -rn | head -20
internal/auth/middleware_test.go:271
internal/auth/claims_test.go:249
internal/auth/jwt_test.go:248
internal/auth/claims.go:204
internal/auth/middleware.go:196
internal/auth/cache_test.go:144
internal/auth/jwt.go:139
internal/auth/cache.go:81

$ rg -c '' --glob 'pkg/interfaceschema/**/*.go' | sort -t: -k2 -rn | head -20
pkg/interfaceschema/parse.go:664
pkg/interfaceschema/types.go:406
pkg/interfaceschema/parse_test.go:352
pkg/interfaceschema/violations_test.go:309
pkg/interfaceschema/compat_test.go:304
pkg/interfaceschema/compile_test.go:235
pkg/interfaceschema/trie_test.go:216
pkg/interfaceschema/required_encrypted_test.go:194
pkg/interfaceschema/trie.go:161
pkg/interfaceschema/compat.go:131
pkg/interfaceschema/compile.go:122
pkg/interfaceschema/violations.go:82

→ Read pkg/interfaceschema/parse.go
→ Read pkg/interfaceschema/types.go
→ Read pkg/interfaceschema/trie.go
→ Read pkg/interfaceschema/compat.go
→ Read pkg/interfaceschema/compile.go
I'll pick `pkg/interfaceschema` — large package, never reviewed, and foundational to the platform. Let me read the source files.
→ Read pkg/interfaceschema/violations.go
Now let me read the test files to see coverage gaps.
→ Read pkg/interfaceschema/parse_test.go
→ Read pkg/interfaceschema/compile_test.go
→ Read pkg/interfaceschema/trie_test.go
→ Read pkg/interfaceschema/compat_test.go
→ Read pkg/interfaceschema/violations_test.go
→ Read pkg/interfaceschema/required_encrypted_test.go

$ head -5 /root/astrate-mule/go.mod
module github.com/astrate-platform/astrate

go 1.26.1

// Pinned dependency decisions — docs/ROADMAP.md §1.1 (frozen in docs/DESIGN.md).

Now let me look at the invalid fixtures manifest and check the module path so I can write a probe.
→ Read pkg/interfaceschema/testdata/invalid/manifest.json
Now let me write a probe in `/tmp` to verify the two findings: the trie placeholder name conflict across depths, and the expiry overflow.
! permission requested: external_directory (/tmp/*); auto-rejecting
✗ mkdir -p /tmp/probe-interfaceschema && cat > /tmp/probe-interfaceschema/go.mod << 'EOF'
module probe

go 1.26.1

replace github.com/astrate-platform/astrate => /root/astrate-mule
EOF
 failed
Error: The user rejected permission to use this specific tool call.
