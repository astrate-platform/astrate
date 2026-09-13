slug: store-todo-lttb-stale
verdict: done
at:  829f408
ran: 2026-09-13T11:22:23Z on DietPi in 66s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
→ Read internal/store/store.go [offset=120, limit=60]
→ Read internal/store/datastreams.go [offset=360, limit=60]
→ Read internal/appengine/data.go [offset=165, limit=25]
The task: remove/replace the stale TODO at store.go:139-143 that claims the lttb switch is unbuilt, when it's actually already built as `DownsampleLTTB` + `HasToolkitLTTB`. I'll delete the stale block.

← Edit internal/store/store.go
Index: /root/astrate-mule/internal/store/store.go
===================================================================
--- /root/astrate-mule/internal/store/store.go
+++ /root/astrate-mule/internal/store/store.go
@@ -134,14 +134,11 @@
 	return nil
 }
 
 // probeCapabilities detects optional database features once at startup.
-//
-// TODO(extension point, docs/ROADMAP.md §0.1 rule 3 / docs/DESIGN.md §2.5):
-// when timescaledb_toolkit is present, Downsample should switch from the
-// time_bucket+avg default to toolkit lttb() downsampling. The probe already
-// records availability in s.hasToolkit; the time_bucket path in
-// datastreams.go is the always-working default.
+// The toolkit result is recorded in s.hasToolkit and exposed to callers via
+// HasToolkitLTTB, which selects the lttb() path (DownsampleLTTB) over the
+// time_bucket+avg default (Downsample) — docs/DESIGN.md §2.5.
 func (s *Store) probeCapabilities(ctx context.Context) error {
 	const q = `SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'timescaledb_toolkit')`
 	if err := s.pool.QueryRow(ctx, q).Scan(&s.hasToolkit); err != nil {
 		return fmt.Errorf("store: probing timescaledb_toolkit: %w", err)

Now verify formatting and vet the affected package.
$ gofmt -l internal/store/store.go; go vet ./internal/store/ 2>&1
(no output)

Done: Removed the stale TODO at internal/store/store.go:139 — the lttb() switch it predicted is already built, and the probe doc comment now states the real design (probe records toolkit availability → HasToolkitLTTB → caller picks DownsampleLTTB over Downsample).
Files: internal/store/store.go
Verified: gofmt -l internal/store/store.go -> clean; go vet ./internal/store/ -> pass. Comment-only change: no behavior test possible.
Unsure: nothing
Follow-ups: none
