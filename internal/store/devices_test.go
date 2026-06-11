//go:build integration

package store

import (
	"context"
	"errors"
	"net/netip"
	"testing"
	"time"

	"github.com/astrate-platform/astrate/pkg/deviceid"
)

func strPtr(s string) *string { return &s }

func testDevices(t *testing.T, s *Store) {
	ctx := context.Background()

	t.Run("Lifecycle", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		id := mustRegisterDevice(t, s, realm.ID)

		d, err := s.GetDevice(ctx, realm.ID, id)
		if err != nil {
			t.Fatalf("GetDevice: %v", err)
		}
		if d.ID != id || d.RealmID != realm.ID {
			t.Errorf("identity mismatch: %+v", d)
		}
		if d.Status != DeviceStatusRegistered || d.Connected || d.PayloadFormatHint != "bson" {
			t.Errorf("defaults: status=%q connected=%v hint=%q", d.Status, d.Connected, d.PayloadFormatHint)
		}
		if len(d.Introspection) != 0 || len(d.Aliases) != 0 || len(d.Attributes) != 0 {
			t.Errorf("fresh device carries state: %+v", d)
		}
		if d.FirstCredentialsRequest != nil || d.CertSerial != nil {
			t.Errorf("fresh device has credential trail: %+v", d)
		}

		// Re-registration before first credentials request rotates the secret.
		if err := s.RegisterDevice(ctx, realm.ID, id, "rotated-hash"); err != nil {
			t.Fatalf("re-register before credentials: %v", err)
		}
		d, err = s.GetDevice(ctx, realm.ID, id)
		if err != nil {
			t.Fatal(err)
		}
		if d.CredentialsSecretHash != "rotated-hash" {
			t.Errorf("secret not rotated: %q", d.CredentialsSecretHash)
		}

		// Credentials issuance stamps the cert and confirms the device.
		reqIP := netip.MustParseAddr("192.0.2.10")
		if err := s.SetDeviceCredentials(ctx, realm.ID, id, "0a1b2c", "aki-1", reqIP); err != nil {
			t.Fatalf("SetDeviceCredentials: %v", err)
		}
		d, err = s.GetDevice(ctx, realm.ID, id)
		if err != nil {
			t.Fatal(err)
		}
		if d.Status != DeviceStatusConfirmed || d.FirstCredentialsRequest == nil {
			t.Errorf("after credentials: status=%q first=%v", d.Status, d.FirstCredentialsRequest)
		}
		if d.CertSerial == nil || *d.CertSerial != "0a1b2c" || d.CertAKI == nil || *d.CertAKI != "aki-1" {
			t.Errorf("cert stamp: serial=%v aki=%v", d.CertSerial, d.CertAKI)
		}
		if d.LastCredentialsRequestIP == nil || *d.LastCredentialsRequestIP != reqIP {
			t.Errorf("request IP: %v", d.LastCredentialsRequestIP)
		}

		// Now re-registration must conflict.
		if err := s.RegisterDevice(ctx, realm.ID, id, "again"); !errors.Is(err, ErrDeviceAlreadyConfirmed) {
			t.Errorf("re-register after credentials: got %v, want ErrDeviceAlreadyConfirmed", err)
		}

		// Inhibit round-trip.
		if err := s.SetDeviceInhibited(ctx, realm.ID, id, true); err != nil {
			t.Fatal(err)
		}
		if d, _ = s.GetDevice(ctx, realm.ID, id); d.Status != DeviceStatusInhibited {
			t.Errorf("after inhibit: %q", d.Status)
		}
		if err := s.SetDeviceInhibited(ctx, realm.ID, id, false); err != nil {
			t.Fatal(err)
		}
		if d, _ = s.GetDevice(ctx, realm.ID, id); d.Status != DeviceStatusConfirmed {
			t.Errorf("after uninhibit: %q, want confirmed", d.Status)
		}

		// Connection lifecycle + stats.
		at := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
		connIP := netip.MustParseAddr("203.0.113.7")
		if err := s.SetDeviceConnected(ctx, realm.ID, id, at, connIP); err != nil {
			t.Fatal(err)
		}
		if err := s.AddDeviceStats(ctx, realm.ID, id, 3, 1024); err != nil {
			t.Fatal(err)
		}
		if err := s.AddDeviceStats(ctx, realm.ID, id, 2, 512); err != nil {
			t.Fatal(err)
		}
		d, err = s.GetDevice(ctx, realm.ID, id)
		if err != nil {
			t.Fatal(err)
		}
		if !d.Connected || d.LastConnection == nil || !d.LastConnection.Equal(at) {
			t.Errorf("after connect: connected=%v last=%v", d.Connected, d.LastConnection)
		}
		if d.LastSeenIP == nil || *d.LastSeenIP != connIP {
			t.Errorf("last seen IP: %v", d.LastSeenIP)
		}
		if d.TotalReceivedMsgs != 5 || d.TotalReceivedBytes != 1536 {
			t.Errorf("stats: msgs=%d bytes=%d", d.TotalReceivedMsgs, d.TotalReceivedBytes)
		}
		if err := s.SetDeviceDisconnected(ctx, realm.ID, id, at.Add(time.Hour)); err != nil {
			t.Fatal(err)
		}
		if d, _ = s.GetDevice(ctx, realm.ID, id); d.Connected || d.LastDisconnection == nil {
			t.Errorf("after disconnect: connected=%v last=%v", d.Connected, d.LastDisconnection)
		}

		// Payload format hint flip (docs/DESIGN.md §3.5.4).
		if err := s.SetPayloadFormatHint(ctx, realm.ID, id, "json"); err != nil {
			t.Fatal(err)
		}
		if d, _ = s.GetDevice(ctx, realm.ID, id); d.PayloadFormatHint != "json" {
			t.Errorf("hint: %q", d.PayloadFormatHint)
		}
		if err := s.SetPayloadFormatHint(ctx, realm.ID, id, "xml"); err == nil {
			t.Error("invalid hint accepted")
		}

		// Unregister keeps the row but resets the credential trail.
		if err := s.UnregisterDevice(ctx, realm.ID, id); err != nil {
			t.Fatal(err)
		}
		d, err = s.GetDevice(ctx, realm.ID, id)
		if err != nil {
			t.Fatalf("device gone after unregister: %v", err)
		}
		if d.Status != DeviceStatusRegistered || d.FirstCredentialsRequest != nil || d.CertSerial != nil {
			t.Errorf("after unregister: %+v", d)
		}
		if err := s.RegisterDevice(ctx, realm.ID, id, "fresh-secret"); err != nil {
			t.Errorf("re-register after unregister: %v", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		ghost, err := deviceid.Random()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := s.GetDevice(ctx, realm.ID, ghost); !errors.Is(err, ErrNotFound) {
			t.Errorf("GetDevice: %v", err)
		}
		if err := s.SetDeviceInhibited(ctx, realm.ID, ghost, true); !errors.Is(err, ErrNotFound) {
			t.Errorf("SetDeviceInhibited: %v", err)
		}
		if err := s.UnregisterDevice(ctx, realm.ID, ghost); !errors.Is(err, ErrNotFound) {
			t.Errorf("UnregisterDevice: %v", err)
		}
	})

	t.Run("IntrospectionDiff", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		id := mustRegisterDevice(t, s, realm.ID)

		removed, err := s.UpdateIntrospection(ctx, realm.ID, id, map[string]InterfaceVersion{
			"com.example.A": {Major: 1, Minor: 0},
			"com.example.B": {Major: 1, Minor: 2},
		})
		if err != nil {
			t.Fatalf("first introspection: %v", err)
		}
		if len(removed) != 0 {
			t.Errorf("first introspection removed %v", removed)
		}

		// Minor bump keeps A; B disappears.
		removed, err = s.UpdateIntrospection(ctx, realm.ID, id, map[string]InterfaceVersion{
			"com.example.A": {Major: 1, Minor: 1},
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(removed) != 1 || removed["com.example.B"] != (InterfaceVersion{Major: 1, Minor: 2}) {
			t.Errorf("removed after B drop: %v", removed)
		}

		// Major change of A moves the old pair too.
		removed, err = s.UpdateIntrospection(ctx, realm.ID, id, map[string]InterfaceVersion{
			"com.example.A": {Major: 2, Minor: 0},
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(removed) != 1 || removed["com.example.A"] != (InterfaceVersion{Major: 1, Minor: 1}) {
			t.Errorf("removed after A major bump: %v", removed)
		}

		d, err := s.GetDevice(ctx, realm.ID, id)
		if err != nil {
			t.Fatal(err)
		}
		if d.Introspection["com.example.A"] != (InterfaceVersion{Major: 2, Minor: 0}) || len(d.Introspection) != 1 {
			t.Errorf("introspection: %v", d.Introspection)
		}
		if d.OldIntrospection["com.example.B"] != (InterfaceVersion{Major: 1, Minor: 2}) ||
			d.OldIntrospection["com.example.A"] != (InterfaceVersion{Major: 1, Minor: 1}) {
			t.Errorf("old_introspection: %v", d.OldIntrospection)
		}

		ghost, _ := deviceid.Random()
		if _, err := s.UpdateIntrospection(ctx, realm.ID, ghost, nil); !errors.Is(err, ErrNotFound) {
			t.Errorf("introspection of unknown device: %v", err)
		}
	})

	t.Run("AliasesAndAttributes", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		id := mustRegisterDevice(t, s, realm.ID)

		if err := s.PatchDeviceAliases(ctx, realm.ID, id, map[string]*string{
			"serial":   strPtr("sn-001"),
			"location": strPtr("lab"),
		}); err != nil {
			t.Fatalf("PatchDeviceAliases: %v", err)
		}

		got, err := s.GetDeviceByAlias(ctx, realm.ID, "sn-001")
		if err != nil {
			t.Fatalf("GetDeviceByAlias: %v", err)
		}
		if got.ID != id {
			t.Errorf("alias lookup returned %s, want %s", got.ID, id)
		}

		// Replace one tag, remove the other.
		if err := s.PatchDeviceAliases(ctx, realm.ID, id, map[string]*string{
			"serial":   strPtr("sn-002"),
			"location": nil,
		}); err != nil {
			t.Fatal(err)
		}
		d, err := s.GetDevice(ctx, realm.ID, id)
		if err != nil {
			t.Fatal(err)
		}
		if len(d.Aliases) != 1 || d.Aliases["serial"] != "sn-002" {
			t.Errorf("aliases after patch: %v", d.Aliases)
		}
		if _, err := s.GetDeviceByAlias(ctx, realm.ID, "lab"); !errors.Is(err, ErrNotFound) {
			t.Errorf("removed alias still resolves: %v", err)
		}

		if err := s.PatchDeviceAttributes(ctx, realm.ID, id, map[string]*string{
			"firmware": strPtr("1.2.3"),
		}); err != nil {
			t.Fatal(err)
		}
		if d, _ = s.GetDevice(ctx, realm.ID, id); d.Attributes["firmware"] != "1.2.3" {
			t.Errorf("attributes: %v", d.Attributes)
		}
	})

	t.Run("PaginationCursorStability", func(t *testing.T) {
		realm := mustCreateRealm(t, s)
		initial := map[deviceid.ID]bool{}
		for range 5 {
			initial[mustRegisterDevice(t, s, realm.ID)] = true
		}

		n, err := s.CountDevices(ctx, realm.ID)
		if err != nil {
			t.Fatal(err)
		}
		if n != 5 {
			t.Fatalf("CountDevices: %d", n)
		}

		seen := map[deviceid.ID]int{}
		var cursor *deviceid.ID
		pages := 0
		for {
			page, err := s.ListDevices(ctx, realm.ID, cursor, 2)
			if err != nil {
				t.Fatalf("ListDevices: %v", err)
			}
			if len(page) == 0 {
				break
			}
			for _, d := range page {
				seen[d.ID]++
			}
			// A device registered mid-iteration must not disturb the cursor
			// walk over the pre-existing set.
			if pages == 0 {
				mustRegisterDevice(t, s, realm.ID)
			}
			last := page[len(page)-1].ID
			cursor = &last
			pages++
			if pages > 10 {
				t.Fatal("pagination did not terminate")
			}
		}

		for id := range initial {
			if seen[id] != 1 {
				t.Errorf("pre-existing device %s seen %d times, want exactly 1", id, seen[id])
			}
		}
		for id, c := range seen {
			if c != 1 {
				t.Errorf("device %s emitted %d times", id, c)
			}
		}
	})
}
