slug: store-alias-values-taken-test
verdict: done
at:  a07bd6a
ran: 2026-09-08T17:30:33Z on DietPi in 85s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
Task: add a store test for `AliasValuesTaken` covering three branches — no other device has the values, another device has one, same device excluded via `d.id <> $2`.
✱ Glob "internal/store/devices_test.go" 1 match
→ Read internal/store/devices.go [offset=320, limit=80]
Now let me look at the existing test file to understand the test patterns used.
→ Read internal/store/devices_test.go
The existing test uses `mustCreateRealm`, `mustRegisterDevice` helpers. Let me check those helpers and understand the test harness setup. I'll look at the store test helpers.
✱ Grep "func mustCreateRealm|func mustRegisterDevice|func testStore|func Test" in internal/store · 18 matches
→ Read internal/store/store_test.go [offset=225, limit=40]
Now I understand the pattern. I'll add a new subtest to the `testDevices` function in devices_test.go that pins the three branches of `AliasValuesTaken`.

Let me verify the three branches logically against the implementation:
1. No other device carries any value → `false`
2. Another device carries one of the values → `true`
3. Same device (excluded via `d.id <> $2`) carries the value → `false`

I'll add the test before the closing of `testDevices` (after the "StatsAndDelete" subtest). Let me add it with three devices to distinguish branches: `self` (carries a value), `other` (carries another value), and a device that shares a value with `self`.

Let me write a subtest:

← Edit internal/store/devices_test.go
Index: /root/astrate-mule/internal/store/devices_test.go
===================================================================
--- /root/astrate-mule/internal/store/devices_test.go
+++ /root/astrate-mule/internal/store/devices_test.go
@@ -378,5 +378,50 @@
 			t.Errorf("post-delete stats = %d/%d (err %v), want 1/0", total, connected, err)
 		}
 		_ = d2
 	})
+
+	t.Run("AliasValuesTaken", func(t *testing.T) {
+		realm := mustCreateRealm(t, s)
+		self := mustRegisterDevice(t, s, realm.ID)
+		other := mustRegisterDevice(t, s, realm.ID)
+
+		if err := s.PatchDeviceAliases(ctx, realm.ID, self, map[string]*string{
+			"serial": strPtr("sn-self"),
+		}); err != nil {
+			t.Fatalf("PatchDeviceAliases(self): %v", err)
+		}
+		if err := s.PatchDeviceAliases(ctx, realm.ID, other, map[string]*string{
+			"serial": strPtr("sn-other"),
+		}); err != nil {
+			t.Fatalf("PatchDeviceAliases(other): %v", err)
+		}
+
+		// No device in the realm carries any of these values.
+		taken, err := s.AliasValuesTaken(ctx, realm.ID, self, []string{"sn-none", "nope"})
+		if err != nil {
+			t.Fatalf("AliasValuesTaken (unused values): %v", err)
+		}
+		if taken {
+			t.Error("unused values reported as taken")
+		}
+
+		// Another device carries one of the values.
+		taken, err = s.AliasValuesTaken(ctx, realm.ID, self, []string{"sn-other", "nope"})
+		if err != nil {
+			t.Fatalf("AliasValuesTaken (other's value): %v", err)
+		}
+		if !taken {
+			t.Error("another device's alias value not reported as taken")
+		}
+
+		// The same device carries the value: it must not count itself
+		// (excluded via d.id <> $2).
+		taken, err = s.AliasValuesTaken(ctx, realm.ID, self, []string{"sn-self"})
+		if err != nil {
+			t.Fatalf("AliasValuesTaken (own value): %v", err)
+		}
+		if taken {
+			t.Error("device's own alias value counted as taken by itself")
+		}
+	})
 }

Now let me run the gate. The integration-tagged tests require a database (PostgreSQL), which is a `[legion]`/container task per AGENTS.md. Let me check what the test requires and whether the gate can run container-free.

Let me check the store_test.go setup for how the test DB is connected.
→ Read internal/store/store_test.go [limit=60]
This test starts a TimescaleDB container via `testutil.StartTimescale`, which is an integration suite (`-tags integration`) that needs Docker/containers. Per AGENTS.md, I must not run these on this Mac and the Legion Go needs to probe. But wait — this is not a `[legion]` tagged task; it's `[auto]`. Let me check whether docker is available at all, and whether the Legion Go is on.

Per AGENTS.md: before any work needing containers, probe ONCE cheaply with `ssh -o ConnectTimeout=5 legion 'echo ok'`. If it does not answer, stop. But this is a container-free-eligible test? No — integration-tagged suites and anything starting a testcontainer are not allowed on the Mac. Docker itself is not allowed on the Mac.

Let me probe the Legion Go first to see if a container test run is possible.

$ ssh -o ConnectTimeout=5 -o BatchMode=yes legion 'echo ok' 2>&1
ssh: connect to host legion.tail4db093.ts.net port 22: Connection timed out

The Legion Go is off. Let me check what environment I'm actually on (whether this is the Pi) and whether Docker is even available here.
! permission requested: external_directory (/etc/*); auto-rejecting
✗ uname -a; echo "---"; which docker 2>&1; echo "---"; cat /etc/os-release 2>/dev/null | head -3; echo "---"; nproc; free -h | head -2 failed
Error: The user rejected permission to use this specific tool call.
