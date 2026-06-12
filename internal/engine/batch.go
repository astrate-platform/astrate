package engine

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/astrate-platform/astrate/internal/store"
	"github.com/astrate-platform/astrate/pkg/interfaceschema"
	"github.com/astrate-platform/astrate/pkg/payload"
)

// batcher is one shard's micro-batch (docs/DESIGN.md §1.4): ops accumulate
// until BatchMaxRows or BatchMaxWait (the shard loop owns the timer), then
// flush applies property ops in arrival order and commits the datastream
// rows through store.AppendDatastreams in a single transaction. Every Ack is
// released only after that commit (docs/DESIGN.md §5.3).
//
// Property ops are individually idempotent (last-value-wins upserts and
// deletes), so applying them as single statements before the datastream
// transaction keeps the at-least-once contract: any failure leaves the whole
// batch unacknowledged and the retry re-applies them harmlessly.
type batcher struct {
	eng      *Engine
	shardIdx int
	ops      []PersistOp
}

// newBatcher builds the batcher of one shard.
func newBatcher(e *Engine, shardIdx int) *batcher {
	return &batcher{eng: e, shardIdx: shardIdx}
}

// size returns the number of pending ops.
func (b *batcher) size() int { return len(b.ops) }

// add appends one validated op; the shard loop checks size afterwards to
// decide on flushing.
func (b *batcher) add(op PersistOp) { b.ops = append(b.ops, op) }

// flush commits the pending batch, retrying transient failures with
// exponential backoff (DB-outage parking, docs/DESIGN.md §5.3). While parked
// the shard does not consume its channel, so backpressure builds exactly as
// designed. Poisoned batches (integrity-class database rejections) split
// into per-op commits so one bad row cannot wedge the shard. On engine
// shutdown (quit) or context cancellation the batch is abandoned
// unacknowledged: QoS >= 1 senders re-deliver after reconnecting.
func (b *batcher) flush(ctx context.Context) {
	if len(b.ops) == 0 {
		return
	}
	// b.ops stays assigned until flushOps returns: if a panic unwinds
	// through here, the restarted shard loop retries the whole batch
	// (duplicate datastream rows from a partially-finalized batch are
	// tolerated by design — at-least-once, docs/DESIGN.md §5.3).
	b.flushOps(ctx, b.ops)
	b.ops = nil
}

// flushOps drives the retry loop for one slice of ops.
func (b *batcher) flushOps(ctx context.Context, ops []PersistOp) {
	e := b.eng
	start := time.Now()
	backoff := parkBackoffStart
	for {
		err := b.commit(ctx, ops)
		if err == nil {
			e.met.flushSeconds.Observe(time.Since(start).Seconds())
			b.finalize(ops)
			return
		}
		if isPermanentCommitError(err) {
			if len(ops) > 1 {
				e.log.Error("batch rejected by the database; committing ops individually",
					"shard", b.shardIdx, "ops", len(ops), "err", err)
				for i := range ops {
					b.flushOps(ctx, ops[i:i+1])
				}
				return
			}
			// A single op the database permanently refuses: consume it so
			// the shard (and the device) is not wedged forever.
			e.met.internalErrors.Inc()
			e.log.Error("op permanently rejected by the database; dropping",
				"shard", b.shardIdx, "kind", ops[0].Kind.String(), "path", ops[0].Path,
				"device", ops[0].DeviceID.String(), "err", err)
			ops[0].ack()
			return
		}

		e.met.flushRetries.Inc()
		e.log.Warn("batch flush failed; parking shard",
			"shard", b.shardIdx, "ops", len(ops), "backoff", backoff, "err", err)
		select {
		case <-e.quit:
			e.log.Warn("abandoning unflushed batch at shutdown",
				"shard", b.shardIdx, "ops", len(ops))
			return
		case <-ctx.Done():
			e.log.Warn("abandoning unflushed batch on context cancellation",
				"shard", b.shardIdx, "ops", len(ops))
			return
		case <-time.After(backoff):
		}
		backoff = min(backoff*2, parkBackoffCap)
	}
}

// commit applies one attempt: property ops as individual idempotent
// statements in arrival order, then all datastream rows in one transaction.
// Ops whose row conversion fails — impossible for pipeline-validated values,
// so a programming error — are marked broken, consumed, and skipped.
func (b *batcher) commit(ctx context.Context, ops []PersistOp) error {
	e := b.eng
	var ds store.DatastreamBatch
	for i := range ops {
		op := &ops[i]
		if op.broken {
			continue
		}
		switch op.Kind {
		case OpPropertySet:
			p, err := propertyRow(op)
			if err != nil {
				b.markBroken(op, err)
				continue
			}
			if err := e.st.UpsertProperty(ctx, *p); err != nil {
				return err
			}
		case OpPropertyUnset:
			if _, err := e.st.UnsetProperty(ctx, op.RealmID, op.DeviceID, op.Interface.ID, op.Path); err != nil {
				return err
			}
		case OpIndividual:
			row, err := individualRow(op)
			if err != nil {
				b.markBroken(op, err)
				continue
			}
			ds.Individual = append(ds.Individual, *row)
		case OpObject:
			row, err := objectRow(op)
			if err != nil {
				b.markBroken(op, err)
				continue
			}
			ds.Objects = append(ds.Objects, *row)
		}
	}
	if len(ds.Individual) > 0 || len(ds.Objects) > 0 {
		return e.st.AppendDatastreams(ctx, ds)
	}
	return nil
}

// finalize releases the acknowledgments of a committed slice — this is the
// ack-after-commit point (docs/DESIGN.md §5.3) — counts the ops, and feeds
// the afterCommit observers (M6b triggers and live fan-out) in-shard, so
// per-device event order matches persistence order.
func (b *batcher) finalize(ops []PersistOp) {
	e := b.eng
	live := ops[:0]
	for i := range ops {
		if ops[i].broken {
			continue
		}
		ops[i].ack()
		e.met.persistOps.WithLabelValues(ops[i].Kind.String()).Inc()
		live = append(live, ops[i])
	}
	if e.afterCommit != nil && len(live) > 0 {
		e.afterCommit(live)
	}
}

// markBroken consumes an op whose row conversion failed: acknowledged (the
// device must not stall on a server-side bug), counted, never persisted.
func (b *batcher) markBroken(op *PersistOp, err error) {
	op.broken = true
	b.eng.met.internalErrors.Inc()
	b.eng.log.Error("persist op dropped: row conversion failed",
		"shard", b.shardIdx, "kind", op.Kind.String(), "interface", op.Interface.Name,
		"path", op.Path, "device", op.DeviceID.String(), "err", err)
	op.ack()
}

// isPermanentCommitError reports database rejections retrying cannot fix:
// integrity violations (SQLSTATE class 23, e.g. an interface row deleted
// mid-flight) and data exceptions (class 22). Everything else — connection
// loss, timeouts, server restarts — is treated as transient.
func isPermanentCommitError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return strings.HasPrefix(pgErr.Code, "23") || strings.HasPrefix(pgErr.Code, "22")
	}
	return false
}

// ---------------------------------------------------------------------------
// PersistOp → store row conversion (docs/DESIGN.md §2.4 typed columns, §2.3
// property jsonb rendering).
// ---------------------------------------------------------------------------

// individualRow converts an OpIndividual into its typed-column row: scalars
// land in their dedicated column, arrays in the pre-encoded JSON value_array.
func individualRow(op *PersistOp) (*store.IndividualRow, error) {
	row := &store.IndividualRow{
		RealmID:     op.RealmID,
		DeviceID:    op.DeviceID,
		InterfaceID: op.Interface.ID,
		EndpointID:  op.Mapping.EndpointID,
		Path:        op.Path,
		TS:          op.TS,
		ReceptionTS: op.ReceptionTS,
	}
	vt := op.Mapping.ValueType
	if vt.IsArray() {
		js, err := encodeValueJSON(op.Value)
		if err != nil {
			return nil, err
		}
		row.ValueArray = js
		return row, nil
	}
	// The column is chosen by the DECLARED type, with a checked assertion:
	// a mismatched value (a pipeline bug) must fail, never land in the
	// wrong column.
	mismatch := func() error {
		return fmt.Errorf("engine: %T value on a %s mapping", op.Value, vt)
	}
	switch vt {
	case interfaceschema.Double:
		v, ok := op.Value.(float64)
		if !ok {
			return nil, mismatch()
		}
		row.ValueDouble = &v
	case interfaceschema.Integer:
		v, ok := op.Value.(int32)
		if !ok {
			return nil, mismatch()
		}
		row.ValueInteger = &v
	case interfaceschema.LongInteger:
		v, ok := op.Value.(int64)
		if !ok {
			return nil, mismatch()
		}
		row.ValueLonginteger = &v
	case interfaceschema.Boolean:
		v, ok := op.Value.(bool)
		if !ok {
			return nil, mismatch()
		}
		row.ValueBoolean = &v
	case interfaceschema.String:
		v, ok := op.Value.(string)
		if !ok {
			return nil, mismatch()
		}
		row.ValueString = &v
	case interfaceschema.BinaryBlob:
		v, ok := op.Value.([]byte)
		if !ok {
			return nil, mismatch()
		}
		if v == nil {
			v = []byte{}
		}
		row.ValueBinaryblob = v
	case interfaceschema.DateTime:
		v, ok := op.Value.(time.Time)
		if !ok {
			return nil, mismatch()
		}
		row.ValueDatetime = &v
	default:
		return nil, fmt.Errorf("engine: unhandled value type %s", vt)
	}
	return row, nil
}

// objectRow converts an OpObject into its jsonb-document row.
func objectRow(op *PersistOp) (*store.ObjectRow, error) {
	if _, ok := op.Value.(map[string]payload.Value); !ok {
		return nil, fmt.Errorf("engine: %T value on an object-aggregated publish", op.Value)
	}
	js, err := encodeValueJSON(op.Value)
	if err != nil {
		return nil, err
	}
	return &store.ObjectRow{
		RealmID:     op.RealmID,
		DeviceID:    op.DeviceID,
		InterfaceID: op.Interface.ID,
		Path:        op.Path,
		TS:          op.TS,
		ReceptionTS: op.ReceptionTS,
		Value:       js,
	}, nil
}

// propertyRow converts an OpPropertySet into its properties row.
func propertyRow(op *PersistOp) (*store.Property, error) {
	js, err := encodeValueJSON(op.Value)
	if err != nil {
		return nil, err
	}
	return &store.Property{
		RealmID:     op.RealmID,
		DeviceID:    op.DeviceID,
		InterfaceID: op.Interface.ID,
		EndpointID:  op.Mapping.EndpointID,
		Path:        op.Path,
		Value:       js,
		ValueType:   op.Mapping.ValueType,
	}, nil
}

// maxSafeJSONInt is the largest integer JSON consumers can assume exact
// (2^53 - 1); larger longintegers render as decimal strings, mirroring the
// JSON payload profile (docs/DESIGN.md §3.5.3) and the §2.3 re-encoding
// contract.
const maxSafeJSONInt = int64(1)<<53 - 1

// jsonTimeLayout renders datetime values: RFC 3339, UTC, millisecond
// precision (the Astarte datetime resolution).
const jsonTimeLayout = "2006-01-02T15:04:05.000Z07:00"

// encodeValueJSON renders a decoded payload value as its canonical jsonb
// form (docs/DESIGN.md §2.3): binaryblob as padded standard base64,
// datetime as RFC 3339 milli, longinteger beyond 2^53 as a decimal string.
func encodeValueJSON(v payload.Value) ([]byte, error) {
	jv, err := jsonable(v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(jv)
}

// jsonable maps the closed payload.Value set onto encoding/json-friendly
// values under the conventions above.
func jsonable(v payload.Value) (any, error) {
	switch x := v.(type) {
	case float64, int32, bool, string:
		return x, nil
	case int64:
		if x > maxSafeJSONInt || x < -maxSafeJSONInt {
			return strconv.FormatInt(x, 10), nil
		}
		return json.Number(strconv.FormatInt(x, 10)), nil
	case []byte:
		return base64.StdEncoding.EncodeToString(x), nil
	case time.Time:
		return x.UTC().Format(jsonTimeLayout), nil
	case []float64:
		return jsonableSlice(x)
	case []int32:
		return jsonableSlice(x)
	case []int64:
		return jsonableSlice(x)
	case []bool:
		return jsonableSlice(x)
	case []string:
		return jsonableSlice(x)
	case [][]byte:
		return jsonableSlice(x)
	case []time.Time:
		return jsonableSlice(x)
	case map[string]payload.Value:
		out := make(map[string]any, len(x))
		for k, e := range x {
			jv, err := jsonable(e)
			if err != nil {
				return nil, fmt.Errorf("key %q: %w", k, err)
			}
			out[k] = jv
		}
		return out, nil
	default:
		return nil, fmt.Errorf("engine: %T is not a persistable payload value", v)
	}
}

// jsonableSlice maps one homogeneous array.
func jsonableSlice[T any](xs []T) (any, error) {
	out := make([]any, len(xs))
	for i, e := range xs {
		jv, err := jsonable(e)
		if err != nil {
			return nil, fmt.Errorf("element %d: %w", i, err)
		}
		out[i] = jv
	}
	return out, nil
}
