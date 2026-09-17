slug: recipe-hygiene
verdict: timeout
at:  bb8bd4e
ran: 2026-09-17T21:01:16Z on DietPi in 1376s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/recipes/hygiene.md
Read the recipe. It's a proposal job: run the hygiene checks, propose findings as task lines in `.mule/todo.md` and version notes in `.mule/for-giulio.md`. No `gh issue create` needed (recipe doesn't ask), no implementation, no git.

Let me check the environment and current state first.

$ ls -la /root/astrate-mule/.mule/ 2>&1
total 4980
drwxr-xr-x  7 root root    4096 Sep 17 22:21 .
drwxr-xr-x 17 root root    4096 Sep  6 13:30 ..
-rw-r--r--  1 root root      14 Sep 17 22:38 .budget
-rw-r--r--  1 root root     364 Jul 27 22:18 .gitignore
-rw-r--r--  1 root root      11 Sep 17 22:21 .heartbeat
-rw-r--r--  1 root root       0 Sep 17 22:38 .last-output
-rw-r--r--  1 root root       0 Sep 17 22:38 .lock
-rw-r--r--  1 root root       8 Sep 17 22:38 .rotation
-rw-r--r--  1 root root    7605 Sep  4 21:23 MULE.md
-rw-r--r--  1 root root    2762 Aug 31 01:45 config
-rw-r--r--  1 root root 4895054 Sep 17 22:38 cron.log
drwxr-xr-x  2 root root    4096 Sep 15 19:13 failed
-rw-r--r--  1 root root   37948 Sep 17 21:49 for-giulio.md
-rw-r--r--  1 root root    1438 Jul 27 19:16 hosts
-rw-r--r--  1 root root    1433 Jul 27 05:29 hosts.example
-rw-r--r--  1 root root    7446 Sep 17 22:21 log.md
-rw-r--r--  1 root root    8456 Sep  4 21:23 milestones.md
drwxr-xr-x  2 root root    4096 Sep  4 20:46 recipes
drwxr-xr-x  3 root root    4096 Sep  6 19:40 reports
drwxr-xr-x  2 root root    4096 Sep 17 21:11 reviews
drwxr-xr-x  2 root root    4096 Sep  7 21:11 tasks
-rw-r--r--  1 root root   63437 Sep 17 22:21 todo.md

$ go version; which go govulncheck golangci-lint 2>&1; uname -a
go version go1.26.7 linux/arm64
/usr/bin/go
Linux DietPi 6.18.39+rpt-rpi-v8 #1 SMP PREEMPT Debian 1:6.18.39-1+rpt1 (2026-07-29) aarch64 GNU/Linux

On the Pi (linux/arm64). `govulncheck` and `golangci-lint` are not installed on PATH. Let me run the recipe checks and look at current proposal files.
$ rg -n 'TODO|FIXME|XXX|HACK' internal/ pkg/ cmd/ 2>/dev/null | head -30; echo "---exit $?---"
internal/flow/blocks/httpblocks.go:147:	msg, err := s.fetch(context.TODO())
internal/flow/blocks/astartesource/source.go:56:	return s.drain(context.TODO(), false)
internal/auth/claims_test.go:57:		{"alternation anchoring quirk left", []string{".*::devices|other"}, "GET", "devicesXXX", true},
internal/auth/claims_test.go:58:		{"alternation anchoring quirk right", []string{".*::devices|other"}, "GET", "XXXother", true},
internal/broker/intake.go:56:// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §1.4):
---exit 0---

$ ls ~/go/bin/ /usr/local/bin/ 2>&1 | rg -i 'golangci|govuln' ; command -v golangci-lint govulncheck; echo "done"
golangci-lint
govulncheck
done

$ go list -m -u all 2>/dev/null | rg '\[' | head -20; echo "---exit $?---"
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
github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c [v0.0.0-20260917193836-60816514402b]
github.com/Azure/go-autorest/autorest/adal v0.9.16 (retracted) [v0.9.24] (deprecated)
github.com/Azure/go-autorest/autorest/date v0.3.0 [v0.3.1] (deprecated)
github.com/Azure/go-autorest/logger v0.2.1 [v0.2.2] (deprecated)
github.com/Azure/go-autorest/tracing v0.6.0 [v0.6.1] (deprecated)
github.com/ClickHouse/clickhouse-go v1.4.3 [v1.5.4]
github.com/DataDog/zstd v1.4.5 [v1.5.7]
---exit 0---

$ ls -la ~/go/bin/ 2>&1; echo "==="; ls -la /usr/local/bin/ 2>&1 | rg -i 'golangci|govuln'
total 85260
drwxr-xr-x 2 root root     4096 Sep 15 20:20 .
drwxr-xr-x 4 root root     4096 Jul 28 23:52 ..
-rwxr-xr-x 1 root root  7668042 Sep  2 21:44 goimports
-rwxr-xr-x 1 root root 65622584 Jul 28 23:52 golangci-lint
-rwxr-xr-x 1 root root 13997377 Sep 15 20:20 govulncheck
===

Both tools exist under `/root/go/bin/`. Running the checks — govulncheck (highest priority in the recipe), then lint, then the test scan.
