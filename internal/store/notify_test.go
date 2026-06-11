//go:build integration

package store

import (
	"context"
	"strconv"
	"testing"
	"time"
)

func testNotify(t *testing.T, s *Store) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := s.Listen(ctx, ChannelInterfaces)
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}

	// Emit until received: the listener connects asynchronously, and after a
	// kill it re-dials with backoff, so notifications sent in the gap are
	// lost by design. Distinct realm IDs per phase keep late duplicates from
	// earlier phases from satisfying the wrong wait.
	waitNotify := func(t *testing.T, realmID int16) {
		t.Helper()
		want := strconv.Itoa(int(realmID))
		deadline := time.After(20 * time.Second)
		tick := time.NewTicker(200 * time.Millisecond)
		defer tick.Stop()
		for {
			select {
			case n, ok := <-ch:
				if !ok {
					t.Fatal("notification channel closed early")
				}
				if n.Channel == ChannelInterfaces && n.Payload == want {
					return
				}
			case <-tick.C:
				if err := s.NotifyInterfacesChanged(context.Background(), realmID); err != nil {
					t.Fatalf("NotifyInterfacesChanged: %v", err)
				}
			case <-deadline:
				t.Fatalf("no notification with payload %q within deadline", want)
			}
		}
	}

	t.Run("EmitReceive", func(t *testing.T) { waitNotify(t, 41) })

	t.Run("SurvivesBackendKill", func(t *testing.T) {
		// Terminate the dedicated LISTEN backend; pool connections never run
		// a LISTEN statement, so the filter only matches the listener.
		rows, err := s.pool.Query(context.Background(), `
			SELECT pg_terminate_backend(pid) FROM pg_stat_activity
			WHERE datname = current_database()
			  AND pid <> pg_backend_pid()
			  AND query ILIKE 'listen %'`)
		if err != nil {
			t.Fatalf("pg_terminate_backend: %v", err)
		}
		killed := 0
		for rows.Next() {
			var ok bool
			if err := rows.Scan(&ok); err != nil {
				t.Fatal(err)
			}
			if ok {
				killed++
			}
		}
		rows.Close()
		if killed == 0 {
			t.Fatal("no LISTEN backend found to terminate")
		}
		t.Logf("terminated %d listener backend(s)", killed)

		waitNotify(t, 42)
	})

	t.Run("ClosesOnCancel", func(t *testing.T) {
		cancel()
		deadline := time.After(10 * time.Second)
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					return
				}
			case <-deadline:
				t.Fatal("channel not closed after context cancellation")
			}
		}
	})
}
