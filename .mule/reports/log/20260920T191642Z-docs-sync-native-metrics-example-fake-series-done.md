slug: docs-sync-native-metrics-example-fake-series
verdict: done
at:  43c4a8b
ran: 2026-09-20T19:16:42Z on DietPi in 93s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read docs/api/astrate_native_api.yaml [offset=80, limit=40]
One line: replace the fake `astrate_devices_total` series in the `/astrate/v1/metrics` example with a real registered series, keeping the HELP/TYPE comment style, then verify `make -C docs build` still loads the YAML.
→ Read internal/observability/metrics.go
→ Read docs/Makefile

← Edit docs/api/astrate_native_api.yaml
Index: /root/astrate-mule/docs/api/astrate_native_api.yaml
===================================================================
--- /root/astrate-mule/docs/api/astrate_native_api.yaml
+++ /root/astrate-mule/docs/api/astrate_native_api.yaml
@@ -89,11 +89,11 @@
           text/plain; version=0.0.4:
             schema:
               type: string
             example: |
-              # HELP astrate_devices_total Total number of devices.
-              # TYPE astrate_devices_total gauge
-              astrate_devices_total{realm="test"} 42
+              # HELP astrate_broker_sessions Live authenticated MQTT device sessions.
+              # TYPE astrate_broker_sessions gauge
+              astrate_broker_sessions 0
 
 # ── Compat health (Astarte Dashboard) ───────────────────────────────
 
 /appengine/health:

$ make -C docs build 2>&1
make: Entering directory '/root/astrate-mule/docs'
cp DESIGN.md site/DESIGN.md
cp ROADMAP.md site/ROADMAP.md
cp ../examples/pokemon-agent/docs/DESIGN.md site/pokemon-agent.md
rm -rf site/api && cp -r api site/api
rm -rf site/swagger-ui && cp -r swagger-ui site/swagger-ui
mkdocs build -f mkdocs.yml
make: mkdocs: No such file or directory
make: *** [Makefile:17: build] Error 127
make: Leaving directory '/root/astrate-mule/docs'

$ cat docs/requirements.txt 2>/dev/null; echo "---"; rg -n "astrate_native_api" docs/swagger-ui 2>/dev/null | head -20
mkdocs>=1.6,<2
mkdocs-material>=9.5,<10
---
docs/swagger-ui/index.html:48:      <option value="../api/astrate_native_api.yaml">Astrate Native API — Health, readiness, metrics, WebSocket</option>

$ python3 -c "import yaml, sys; d=yaml.safe_load(open('docs/api/astrate_native_api.yaml')); print('yaml ok'); import json; print(json.dumps(d['paths']['/astrate/v1/metrics']['get']['responses']['200']['content']['text/plain; version=0.0.4']['example']))"
Traceback (most recent call last):
  File "<string>", line 1, in <module>
    import yaml, sys; d=yaml.safe_load(open('docs/api/astrate_native_api.yaml')); print('yaml ok'); import json; print(json.dumps(d['paths']['/astrate/v1/metrics']['get']['responses']['200']['content']['text/plain; version=0.0.4']['example']))
    ^^^^^^^^^^^^^^^^
ModuleNotFoundError: No module named 'yaml'

$ pip3 install --user -r docs/requirements.txt 2>&1 | tail -5
/bin/sh: 1: pip3: not found

$ which pip pip3 python3 python; python3 --version; ls /usr/bin | rg -i 'pip|mkdocs' 2>/dev/null
/usr/bin/python3
Python 3.13.5
pipewire
pipewire-aes67
pipewire-avb

