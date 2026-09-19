slug: recipe-hygiene
verdict: timeout
at:  757407a
ran: 2026-09-19T17:08:09Z on DietPi in 1251s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md

$ ls -la
total 156
drwxr-xr-x 17 root root  4096 Sep  6 13:30 .
drwx------ 28 root root  4096 Sep  4 22:28 ..
drwxr-xr-x  3 root root  4096 Sep  4 20:46 .agents
-rw-r--r--  1 root root   207 Sep  4 20:46 .dockerignore
drwxr-xr-x  8 root root  4096 Sep 19 13:35 .git
drwxr-xr-x  3 root root  4096 Jul 27 05:10 .github
-rw-r--r--  1 root root  1133 Sep  4 21:26 .gitignore
-rw-r--r--  1 root root  1067 Jul 27 05:10 .golangci.yml
drwxr-xr-x  7 root root  4096 Sep 19 12:52 .mule
-rw-r--r--  1 root root  1578 Sep  4 20:46 AGENTS.md
-rw-r--r--  1 root root  1315 Jul 27 05:10 Dockerfile
-rw-r--r--  1 root root 14394 Jul 27 05:10 LICENSE
-rw-r--r--  1 root root  2545 Jul 27 05:10 Makefile
-rw-r--r--  1 root root  3482 Sep  4 20:46 README.md
drwxr-xr-x  5 root root  4096 Jul 27 05:10 bench
drwxr-xr-x  3 root root  4096 Jul 27 05:10 cmd
drwxr-xr-x  4 root root  4096 Jul 27 05:10 deploy
-rw-r--r--  1 root root  4317 Jul 27 22:18 docker-compose.yml
drwxr-xr-x  6 root root  4096 Sep  4 20:46 docs
drwxr-xr-x  4 root root  4096 Sep  4 20:46 examples
-rw-r--r--  1 root root  5092 Sep  4 20:46 go.mod
-rw-r--r--  1 root root 20833 Sep  4 20:46 go.sum
drwxr-xr-x 17 root root  4096 Sep  4 20:46 internal
drwxr-xr-x  2 root root  4096 Sep  4 20:46 migrations
-rw-r--r--  1 root root   102 Sep  4 20:46 opencode.json
drwxr-xr-x  6 root root  4096 Jul 27 05:10 pkg
drwxr-xr-x 30 root root  4096 Sep 17 13:43 site-dist
drwxr-xr-x  4 root root  4096 Jul 27 05:10 test
drwxr-xr-x  3 root root  4096 Sep  4 21:26 tools

$ ls .mule 2>/dev/null; echo "---recipes---"; ls .mule/recipes 2>/dev/null
MULE.md
config
cron.log
failed
for-giulio.md
hosts
hosts.example
log.md
milestones.md
recipes
reports
reviews
tasks
todo.md
---recipes---
astarte-upstream.md
benchmarks.md
code-review.md
cross-project-survey.md
docs-sync.md
github-issues.md
hygiene.md
legion-go.md
milestones.md

# Todos
[ ] Run deps-vs-newer check (go list -m -u all)
[ ] Run govulncheck
[ ] Run golangci-lint
[ ] Search TODO/FIXME/XXX/HACK
[ ] Run go test for skips/no-test-files
[ ] Write proposals to .mule/todo.md + evidence file

$ go version; which govulncheck; which golangci-lint; which gh; ls ~/go/bin 2>/dev/null
go version go1.26.7 linux/arm64
/usr/bin/gh
goimports
golangci-lint
govulncheck

$ go list -m -u all 2>/dev/null | rg '\[|=> ' | head -20
cel.dev/expr v0.24.0 [v0.25.3]
cloud.google.com/go v0.121.6 [v0.123.0]
cloud.google.com/go/auth v0.16.4 [v0.23.3]
cloud.google.com/go/compute/metadata v0.8.0 [v0.9.1]
cloud.google.com/go/iam v1.5.2 [v1.13.0]
cloud.google.com/go/longrunning v0.6.7 [v1.2.0]
cloud.google.com/go/monitoring v1.24.2 [v1.30.0]
cloud.google.com/go/spanner v1.85.0 [v1.95.1]
cloud.google.com/go/storage v1.56.0 [v1.68.0]
github.com/99designs/keyring v1.2.1 [v1.2.2]
github.com/Azure/azure-sdk-for-go/sdk/azcore v1.4.0 [v1.23.1]
github.com/Azure/azure-sdk-for-go/sdk/internal v1.1.2 [v1.12.0]
github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v1.0.0 [v1.8.1]
github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c [v0.0.0-20260917205352-e937bb47801a]
github.com/Azure/go-autorest/autorest/adal v0.9.16 (retracted) [v0.9.24] (deprecated)
github.com/Azure/go-autorest/autorest/date v0.3.0 [v0.3.1] (deprecated)
github.com/Azure/go-autorest/logger v0.2.1 [v0.2.2] (deprecated)
github.com/Azure/go-autorest/tracing v0.6.0 [v0.6.1] (deprecated)
github.com/ClickHouse/clickhouse-go v1.4.3 [v1.5.4]
github.com/DataDog/zstd v1.4.5 [v1.5.7]

