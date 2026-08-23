package blocks

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"sync/atomic"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
)

// RandomSource emits one random value per configured interval.
//
// Config keys:
//   - type (required string): integer | real | boolean
//   - interval_ms (int, default 1000): must be positive; every Emit (including
//     the first) waits one interval before returning a message
//   - min / max (numbers, optional): bounds — inclusive [min,max] for integer
//     (defaults 0..100), [min,max] for real (defaults 0.0..1.0); both ignored
//     for boolean
//   - key (string, default "random"): Key of the emitted messages
func RandomSource(name string, config map[string]any, _ flow.Deps) (flow.Block, error) {
	cfg, err := parseRandomSourceConfig(config)
	if err != nil {
		return nil, fmt.Errorf("random_source: %w", err)
	}
	return &randomSource{
		name:     name,
		kind:     cfg.kind,
		interval: cfg.interval,
		minInt:   cfg.minInt,
		maxInt:   cfg.maxInt,
		minReal:  cfg.minReal,
		maxReal:  cfg.maxReal,
		key:      cfg.key,
	}, nil
}

type randomKind uint8

const (
	randomInteger randomKind = iota
	randomReal
	randomBoolean
)

type randomSourceConfig struct {
	kind     randomKind
	interval time.Duration
	// Integer bounds are inclusive.
	minInt, maxInt int64
	// Real bounds.
	minReal, maxReal float64
	key              string
}

func parseRandomSourceConfig(config map[string]any) (randomSourceConfig, error) {
	cfg := randomSourceConfig{interval: time.Second}
	switch s := stringConfig(config, "type", ""); s {
	case "integer":
		cfg.kind = randomInteger
		cfg.minInt, cfg.maxInt = 0, 100
		cfg.minReal, cfg.maxReal = 0, 100
	case "real":
		cfg.kind = randomReal
		cfg.minReal, cfg.maxReal = 0, 1
	case "boolean":
		cfg.kind = randomBoolean
	default:
		return cfg, fmt.Errorf("unknown type %q (want integer|real|boolean)", s)
	}

	if v, ok := config["interval_ms"]; ok && v != nil {
		n, err := numAsInt64(v)
		if err != nil || n <= 0 {
			return cfg, fmt.Errorf("interval_ms must be positive")
		}
		cfg.interval = time.Duration(n) * time.Millisecond
	}

	var haveMin, haveMax bool
	var fmin, fmax float64
	if v, ok := config["min"]; ok && v != nil {
		f, err := numAsFloat64(v)
		if err != nil {
			return cfg, fmt.Errorf("min must be a number")
		}
		fmin = f
		cfg.minInt = int64(f)
		cfg.minReal = f
		haveMin = true
	}
	if v, ok := config["max"]; ok && v != nil {
		f, err := numAsFloat64(v)
		if err != nil {
			return cfg, fmt.Errorf("max must be a number")
		}
		fmax = f
		cfg.maxInt = int64(f)
		cfg.maxReal = f
		haveMax = true
	}
	if (haveMin || haveMax) && fmin > fmax {
		return cfg, fmt.Errorf("min must be <= max")
	}
	cfg.key = stringConfig(config, "key", "random")
	return cfg, nil
}

type randomSource struct {
	name     string
	kind     randomKind
	interval time.Duration
	minInt   int64
	maxInt   int64
	minReal  float64
	maxReal  float64
	key      string
	stopped  atomic.Bool
}

var (
	_ flow.Block   = (*randomSource)(nil)
	_ flow.Source  = (*randomSource)(nil)
	_ flow.Stopper = (*randomSource)(nil)
)

func (s *randomSource) Name() string { return s.name }

// Process implements flow.Block for the non-pump path: it produces one value
// immediately without waiting for the interval (a stopped source yields none).
func (s *randomSource) Process(_ *flow.Message) ([]*flow.Message, error) {
	if s.stopped.Load() {
		return nil, nil
	}
	return []*flow.Message{s.next()}, nil
}

// Emit implements flow.Source: it blocks until one interval elapses or ctx is
// cancelled, then returns exactly one message. A stopped source never returns
// a message.
func (s *randomSource) Emit(ctx context.Context) ([]*flow.Message, error) {
	timer := time.NewTimer(s.interval)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
	}
	if s.stopped.Load() {
		return nil, nil
	}
	return []*flow.Message{s.next()}, nil
}

// Stop implements flow.Stopper: idempotent; after Stop no further messages
// are emitted.
func (s *randomSource) Stop() { s.stopped.Store(true) }

func (s *randomSource) next() *flow.Message {
	msg := &flow.Message{
		Key:       s.key,
		Timestamp: time.Now().UnixMicro(),
	}
	switch s.kind {
	case randomInteger:
		msg.Type = flow.TypeInteger
		msg.Data = s.minInt + rand.Int64N(s.maxInt-s.minInt+1) //nolint:gosec // non-crypto sample values by design (random_source)
	case randomReal:
		msg.Type = flow.TypeReal
		msg.Data = s.minReal + rand.Float64()*(s.maxReal-s.minReal) //nolint:gosec // non-crypto sample values by design (random_source)
	default:
		msg.Type = flow.TypeBoolean
		msg.Data = rand.Int64N(2) == 0 //nolint:gosec // non-crypto sample values by design (random_source)
	}
	return msg
}

// numAsInt64 coerces a config number (int/int64/float64/json.Number) to int64.
func numAsInt64(v any) (int64, error) {
	switch n := v.(type) {
	case int:
		return int64(n), nil
	case int32:
		return int64(n), nil
	case int64:
		return n, nil
	case uint:
		if uint64(n) > math.MaxInt64 {
			return 0, fmt.Errorf("number out of int64 range")
		}
		return int64(n), nil
	case uint64:
		if n > math.MaxInt64 {
			return 0, fmt.Errorf("number out of int64 range")
		}
		return int64(n), nil
	case float64:
		return int64(n), nil
	case json.Number:
		return n.Int64()
	default:
		return 0, fmt.Errorf("not a number")
	}
}

// numAsFloat64 coerces a config number to float64.
func numAsFloat64(v any) (float64, error) {
	switch n := v.(type) {
	case int:
		return float64(n), nil
	case int32:
		return float64(n), nil
	case int64:
		return float64(n), nil
	case float64:
		return n, nil
	case json.Number:
		return n.Float64()
	default:
		return 0, fmt.Errorf("not a number")
	}
}
