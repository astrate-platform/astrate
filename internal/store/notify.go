package store

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

// ChannelInterfaces is the NOTIFY channel signalling interface CRUD
// (docs/DESIGN.md §2.6): the engine LISTENs and rebuilds the affected
// realm's compiled-interface snapshot. The payload is the realm ID in
// decimal. In a single process the in-process callback already fires; the
// channel keeps an optional hot-standby instance coherent.
const ChannelInterfaces = "astrate_interfaces"

// Notification is one received LISTEN/NOTIFY event.
type Notification struct {
	Channel string
	Payload string
}

// listener reconnect backoff bounds.
const (
	listenBackoffStart = 250 * time.Millisecond
	listenBackoffCap   = 5 * time.Second
)

// NotifyInterfacesChanged emits a ChannelInterfaces notification for realmID.
func (s *Store) NotifyInterfacesChanged(ctx context.Context, realmID int16) error {
	_, err := s.pool.Exec(ctx, `SELECT pg_notify($1, $2)`,
		ChannelInterfaces, strconv.Itoa(int(realmID)))
	if err != nil {
		return fmt.Errorf("store: notifying %s: %w", ChannelInterfaces, err)
	}
	return nil
}

// Listen subscribes to a NOTIFY channel on a dedicated connection (LISTEN
// pins a session, so the pool is not used) and streams notifications until
// ctx is cancelled, at which point the returned channel is closed. Lost
// connections are re-dialled with exponential backoff and the LISTEN is
// re-issued; notifications emitted while disconnected are lost, which is
// fine for the cache-invalidation use case (the listener reloads from the
// tables on resubscribe anyway).
func (s *Store) Listen(ctx context.Context, channel string) (<-chan Notification, error) {
	ident := pgx.Identifier{channel}.Sanitize()
	connCfg := s.pool.Config().ConnConfig

	out := make(chan Notification, 16)
	// The lifecycle context is ctx throughout; context.Background appears
	// only for Close calls that must still run after ctx is cancelled.
	go func() { // #nosec G118 -- see comment above
		defer close(out)
		backoff := listenBackoffStart
		for ctx.Err() == nil {
			conn, err := pgx.ConnectConfig(ctx, connCfg.Copy())
			if err != nil {
				if !sleepCtx(ctx, backoff) {
					return
				}
				backoff = min(backoff*2, listenBackoffCap)
				continue
			}
			if _, err := conn.Exec(ctx, "listen "+ident); err != nil {
				_ = conn.Close(context.Background())
				if !sleepCtx(ctx, backoff) {
					return
				}
				backoff = min(backoff*2, listenBackoffCap)
				continue
			}
			backoff = listenBackoffStart

			for {
				n, err := conn.WaitForNotification(ctx)
				if err != nil {
					// Cancelled or connection lost (e.g.
					// pg_terminate_backend): close and either stop or
					// reconnect.
					_ = conn.Close(context.Background())
					break
				}
				select {
				case out <- Notification{Channel: n.Channel, Payload: n.Payload}:
				case <-ctx.Done():
					_ = conn.Close(context.Background())
					return
				}
			}
		}
	}()
	return out, nil
}

// sleepCtx sleeps for d unless ctx is cancelled first; it reports whether
// the full sleep elapsed.
func sleepCtx(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return true
	case <-ctx.Done():
		return false
	}
}
