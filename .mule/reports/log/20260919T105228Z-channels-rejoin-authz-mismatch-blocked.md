slug: channels-rejoin-authz-mismatch
verdict: blocked
at:  ebc3854
ran: 2026-09-19T10:52:28Z on DietPi in 179s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
✗ Read ws.go failed [offset=240, limit=80]
Error: File not found: /root/astrate-mule/ws.go
✱ Glob "**/ws.go" 2 matches
✱ Grep "handleJoin" 5 matches
→ Read internal/appengine/channels/ws.go [offset=255, limit=60]
→ Read .mule/reviews/channels-2026-09-18.md [offset=40, limit=45]
The task's premise says the rejoin branch bypasses authz. Let me check the current working-tree state vs HEAD to be sure I'm reading the live code.

$ git diff -- internal/appengine/channels/ws.go; git log --oneline -5 -- internal/appengine/channels/ws.go
aeaeb47 mule: channels-group-watch-membership-scope: group watches authorize against `WATCH::groups/<name>/…` (ws.go:178-205, conformance channels.json) but `Room.dispatch` ever only filters by `watchEntry.deviceID` (empty for a group watch — room.go:36-40, 196-200) and the trigger compiles with no group/device scope (triggers/match.go:374-396, because upstream carries group_name at the payload top level), so a group-scoped token receives every realm event matching interface/path/on, members or not. Plumb a group-membership resolver into internal/appengine/channels (it knows only `Bus` today), resolve membership at Watch time and store it in `watchEntry`, filter deliveries by it. Add a container-free test with a fake resolver: watch group `probe`, publish an event from a member and a non-member, assert only the member's arrives. Verify the member-filter rule against upstream on the Legion if disputed; the boundary leak stands regardless.
7f3a600 conformance: measure group WATCH, server-owned and size limits against upstream
7e39ddf Measure what a WATCH claim is matched against, and agree with it
85f5157 Read a_ch by upstream's rule, and enforce JOIN on a join
699cf11 Initial commit under updated license terms

No working-tree diff for this file, so HEAD is the live code. Let me confirm when `handleJoin`'s layout was last changed, to see whether the task premise predates it.
$ git show 85f5157 -- internal/appengine/channels/ws.go | sed -n '1,80p'; git log -1 --format='%h %ad %s' --date=short 85f5157
commit 85f5157d397ee558bf1d1563675a89f1a6a54db4
Author: Giulio Malitesta <atsetilam@gmail.com>
Date:   Mon Jul 27 00:16:04 2026 +0200

    Read a_ch by upstream's rule, and enforce JOIN on a join
    
    Upstream reads the Channels claim by a different grammar than the REST
    claims, and the difference is observable. The REST plug compiles the verb
    field into a regex and matches it against the HTTP method, so "GET|POST"
    selects two methods and ".*::.*" grants everything. The Channels socket
    instead partitions a_ch by an exact match on the verb field, keeping only
    entries whose first field is literally JOIN or WATCH and discarding the
    rest, then matches the path regex within the chosen bucket.
    
    So the astartectl-style blanket ".*::.*" authorizes nothing on Channels:
    ".*" is not the string "JOIN". M11 assumed the opposite and read a_ch with
    the REST rule, which is why the whole wire-session test file was built on a
    claim upstream refuses.
    
    Astrate also never checked JOIN at all - any token bearing a_ch could join
    any room in its realm. That check now runs on every join, including a
    rejoin, against the room name with rooms:<realm>: stripped off.
    
    Authorizes keeps the verb-regex behaviour, because that is what upstream's
    REST plug does; only the Channels surface departs from it.
    
    Both rules are mutation-checked. Two divergences the recording surfaced but
    could not settle - the watch path omitting the trigger's match path, and a
    device trigger's device_id having to sit inside simple_trigger - are
    documented as unreconciled rather than changed on a reading.

diff --git a/internal/appengine/channels/ws.go b/internal/appengine/channels/ws.go
index ffd3ec8..19d14a1 100644
--- a/internal/appengine/channels/ws.go
+++ b/internal/appengine/channels/ws.go
@@ -205,6 +205,12 @@ func (s *session) handleJoin(f Frame) {
 		s.sendErr(f, "unauthorized")
 		return
 	}
+	// Upstream authorizes every join, including a rejoin, against the room
+	// name with the realm stripped off.
+	if !s.tok.AuthorizesChannel(auth.VerbJoin, name) {
+		s.sendErr(f, "unauthorized")
+		return
+	}
 
 	s.roomsMu.Lock()
 	if j, ok := s.rooms[f.Topic]; ok {
@@ -263,7 +269,7 @@ func (s *session) handleWatch(f Frame) {
 		s.sendErr(f, "invalid simple_trigger")
 		return
 	}
-	if !s.tok.Authorizes(auth.ClaimChannels, "WATCH", path) {
+	if !s.tok.AuthorizesChannel(auth.VerbWatch, path) {
 		s.sendErr(f, "unauthorized")
 		return
 	}
85f5157 2026-07-27 Read a_ch by upstream's rule, and enforce JOIN on a join

Confirmed: commit `85f5157` (2026-07-27) already placed `AuthorizesChannel(auth.VerbJoin, name)` **before** the rejoin branch, with exactly the comment the task quotes. The check governs both first-join and rejoin paths; the comment is accurate. The review that spawned the task misread the rejoin branch shape. Verifying the tree as-is:
$ gofmt -l internal/appengine/channels/ && go test ./internal/appengine/channels/ 2>&1 | tail -5
ok  	github.com/astrate-platform/astrate/internal/appengine/channels	(cached)

`gofmt -l` empty; package tests pass.

The premise is stale — there is nothing to fix. Commit `85f5157` (2026-07-27) already inserted the `AuthorizesChannel(auth.VerbJoin, name)` check at ws.go:278, *before* both the rejoin branch (284-291) and the first-join branch, with precisely the comment the task quotes. The check governs every join, the comment matches the code, and behaviour already is what the task's option 1 would produce. "Fixing the comment" would make it false; adding a redundant call inside the rejoin branch would be pure noise.

```
Done: channels-rejoin-authz-mismatch premised on stale code — the authz check already covers the rejoin branch (ws.go:278 precedes lines 284-291, since 85f5157), so the comment is accurate and no change is warranted
Files: (none changed)
Verified: gofmt -l internal/appengine/channels/ -> clean (no output); go test ./internal/appengine/channels/ -> pass
Unsure: the review at .mule/reviews/channels-2026-09-18.md#2 misread the branch shape; its line numbers match a code state that predates 85f5157
Follow-ups: none
```
