slug: docs-sync-native-socket-query-token-auth
verdict: done
at:  8148aae
ran: 2026-09-15T18:17:48Z on DietPi in 98s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astrate_native_api.yaml [offset=270, limit=50]
→ Read internal/appengine/stream/ws.go [offset=30, limit=40]
→ Read internal/auth/middleware.go [offset=120, limit=90]
→ Read internal/appengine/channels/ws.go [limit=40, offset=30]
✱ Grep "require\(http.Handler\)|RequireRealm" in internal/appengine/stream · 1 match

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -289,10 +289,9 @@
   ```
 
   **SSE protocol:** each event is a `data:` frame followed by a blank line.
 
-  Authentication: `a_ch` realm JWT (via Authorization header or
-  `?token=` query parameter).
+  Authentication: `a_ch` realm JWT via Authorization header.
 tags: [WebSocket]
 parameters:
   - $ref: "#/components/parameters/RealmName"
   - name: device_id

The Phoenix socket (`channels/ws.go:53-54`) reads `?token=`, but the stream socket route `/astrate/v1/{realm}/socket` only goes through `auth.Middleware.RequireRealm`, which uses `bearerToken` (header-only). The spec over-claims. Aligning the spec to the code: drop the `?token=` claim.
Now verify the Swagger-referenced YAML still loads with a docs build.
$ make -C docs build
make: Entering directory '/root/astrate-mule/docs'
cp DESIGN.md site/DESIGN.md
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: Leaving directory '/root/astrate-mule/docs'
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127

`mkdocs` isn't installed. Let me check the docs Makefile and validate the YAML loads another way.
✱ Glob "docs/swagger-ui/**" 1 match
→ Read docs/Makefile

$ python3 -c "import yaml,sys; d=yaml.safe_load(open('docs/api/astrate_native_api.yaml')); print('YAML OK, top-level paths:', list(d['paths'].keys())[:5])"
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml,sys; d=yaml.safe_load(open('docs/api/astrate_native_api.yaml')); print('YAML OK, top-level paths:', list(d['paths'].keys())[:5])
    ^^^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

✱ Grep "astrate_native_api.yaml|spec-url|url:" in docs/swagger-ui/index.html · 2 matches
The `sync` steps (including `cp -r api site/api`) succeeded; only `mkdocs` is missing, so I'll validate the YAML loads directly.
✱ Grep "astrate_native_api" in docs · 0 matches
→ Read docs/mkdocs.yml

$ which ruby node perl 2>/dev/null; echo "---"; ruby -e "require 'yaml'; YAML.load_file('docs/api/astrate_native_api.yaml'); puts 'YAML OK'" 2>&1 | tail -3
/usr/local/bin/node
/usr/bin/perl
---
/bin/sh: 1: ruby: not found

$ perl -MYAML::Tiny -e 'YAML::Tiny->read("docs/api/astrate_native_api.yaml")->[0]; print "PERL YAML OK\n"' 2>&1 | tail -2; echo "---node---"; node -e "try{require.resolve('js-yaml');console.log('js-yaml present')}catch(e){console.log('js-yaml absent')}"
Can't locate YAML/Tiny.pm in @INC (you may need to install the YAML::Tiny module) (@INC entries checked: /etc/perl /usr/local/lib/aarch64-linux-gnu/perl/5.40.1 /usr/local/share/perl/5.40.1 /usr/lib/aarch64-linux-gnu/perl5/5.40 /usr/share/perl5 /usr/lib/aarch64-linux-gnu/perl-base /usr/lib/aarch64-linux-gnu/perl/5.40 /usr/share/perl/5.40 /usr/local/lib/site_perl).
BEGIN failed--compilation aborted.
---node---
js-yaml absent

✱ Grep "yaml" in go.mod · 100 matches
! permission requested: external_directory (/tmp/yamlcheck/*); auto-rejecting
No YAML tooling is installed anywhere on this box, but the repo's module cache has `gopkg.in/yaml.v3`. I'll validate the edited spec with a throwaway probe in `/tmp`.
✗ Write /tmp/yamlcheck/main.go failed
Error: The user rejected permission to use this specific tool call.
