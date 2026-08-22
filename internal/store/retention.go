package store

import (
	"context"
	"fmt"
)

// EnforceRealmRetentionCeilings applies every realm's
// datastream_maximum_storage_retention (#72): datastream rows older than the
// ceiling are deleted regardless of their interface's own retention policy —
// a set ceiling caps even no_ttl interfaces, which is Astrate's answer to
// upstream's write-time TTL clamp (Astrate stores no per-row TTL). Realms
// without a ceiling are skipped. One statement per hypertable per capped
// realm (no chunk looping at Astrate's scale); a failure on one realm does
// not abort the others — the first error is returned after all realms ran.
func (s *Store) EnforceRealmRetentionCeilings(ctx context.Context) error {
	rows, err := s.pool.Query(ctx, `
		SELECT id, datastream_maximum_storage_retention FROM realms
		WHERE datastream_maximum_storage_retention IS NOT NULL
		  AND datastream_maximum_storage_retention > 0`)
	if err != nil {
		return fmt.Errorf("store: listing retention ceilings: %w", err)
	}
	type capped struct {
		id  int16
		max int64
	}
	var realms []capped
	for rows.Next() {
		var c capped
		if err := rows.Scan(&c.id, &c.max); err != nil {
			rows.Close()
			return fmt.Errorf("store: scanning retention ceiling: %w", err)
		}
		realms = append(realms, c)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("store: listing retention ceilings: %w", err)
	}
	rows.Close()

	var firstErr error
	for _, c := range realms {
		for _, table := range []string{"individual_datastreams", "object_datastreams"} {
			if _, err := s.pool.Exec(ctx,
				`DELETE FROM `+table+` WHERE realm_id = $1 AND ts < now() - make_interval(secs => $2)`,
				c.id, c.max); err != nil && firstErr == nil {
				firstErr = fmt.Errorf("store: sweeping %s of realm id %d past its retention ceiling: %w", table, c.id, err)
			}
		}
	}
	return firstErr
}
