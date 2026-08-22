// Package container implements the Flow "container" block (Design B / #43 PoC).
//
// Operator model: a pipeline stage whose implementation is a Docker image.
// Transport is HTTP request/response (not upstream AMQP). Lifecycle is tied
// to the flow: the container is started at Instantiate (flow start) and
// removed on Stopper.Stop (flow stop).
package container

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
)

// Type is the pipeline block_type string.
const Type = "container"

const defaultContainerPort = 8080

// defaultRunner is used by the catalog constructor. Tests inject via New.
var defaultRunner Runner = &CLIRunner{}

// SetDefaultRunner replaces the Docker runner used by the catalog constructor.
// Prefer New in unit tests. Returns a restore function.
func SetDefaultRunner(r Runner) (restore func()) {
	prev := defaultRunner
	if r == nil {
		defaultRunner = &CLIRunner{}
	} else {
		defaultRunner = r
	}
	return func() { defaultRunner = prev }
}

// Config is the parsed block configuration.
type Config struct {
	Image         string
	Nested        map[string]any // opaque JSON for ASTRATE_FLOW_CONFIG
	ContainerPort int
	Timeout       time.Duration
	ReadyTimeout  time.Duration
}

// Block is a transform that round-trips each Message through a container.
type Block struct {
	name   string
	cfg    Config
	log    *slog.Logger
	bridge *Bridge
	inst   Instance

	mu      sync.Mutex
	stopped bool
}

// Constructor is the catalog entry point (uses defaultRunner).
func Constructor(name string, config map[string]any, deps flow.Deps) (flow.Block, error) {
	return New(name, config, deps, defaultRunner)
}

// New builds and starts a container block with an explicit Runner (tests + PoC).
func New(name string, config map[string]any, deps flow.Deps, runner Runner) (flow.Block, error) {
	if runner == nil {
		return nil, fmt.Errorf("container: docker runner is nil")
	}
	cfg, err := parseConfig(config)
	if err != nil {
		return nil, err
	}

	flowCfgJSON, err := encodeFlowConfigJSON(cfg.Nested)
	if err != nil {
		return nil, err
	}

	labels := map[string]string{
		"astrate.flow":  "1",
		"astrate.block": name,
	}
	if deps.Realm != "" {
		labels["astrate.realm"] = deps.Realm
	}
	if deps.FlowName != "" {
		labels["astrate.flow_name"] = deps.FlowName
	}

	// Unique-ish name: astrate-<realm>-<block>-<nanos>
	cname := dockerName(deps.Realm, name)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	inst, err := runner.Start(ctx, Spec{
		Image:          cfg.Image,
		ContainerPort:  cfg.ContainerPort,
		FlowConfigJSON: flowCfgJSON,
		Labels:         labels,
		Name:           cname,
	})
	if err != nil {
		return nil, err
	}

	bridge := &Bridge{
		BaseURL: inst.BaseURL(),
		Timeout: cfg.Timeout,
		Client:  &http.Client{Timeout: cfg.Timeout},
	}

	readyCtx, readyCancel := context.WithTimeout(context.Background(), cfg.ReadyTimeout)
	defer readyCancel()
	if err := bridge.WaitReady(readyCtx); err != nil {
		_ = inst.Stop(context.Background())
		return nil, fmt.Errorf("container %q: %w", name, err)
	}

	b := &Block{
		name:   name,
		cfg:    cfg,
		log:    slog.Default(),
		bridge: bridge,
		inst:   inst,
	}
	return b, nil
}

func parseConfig(config map[string]any) (Config, error) {
	cfg := Config{
		ContainerPort: defaultContainerPort,
		Timeout:       defaultTimeout,
		ReadyTimeout:  defaultReadyWait,
	}
	if config == nil {
		return cfg, fmt.Errorf("container: image is required")
	}
	img, _ := config["image"].(string)
	img = strings.TrimSpace(img)
	if img == "" {
		return cfg, fmt.Errorf("container: image is required")
	}
	cfg.Image = img

	if raw, ok := config["config"]; ok && raw != nil {
		switch m := raw.(type) {
		case map[string]any:
			cfg.Nested = m
		default:
			return cfg, fmt.Errorf("container: config must be a JSON object")
		}
	}

	if v, ok := config["port"]; ok && v != nil {
		n, err := asInt(v)
		if err != nil || n <= 0 || n > 65535 {
			return cfg, fmt.Errorf("container: port must be an integer 1–65535")
		}
		cfg.ContainerPort = n
	}

	if v, ok := config["timeout_ms"]; ok && v != nil {
		n, err := asInt(v)
		if err != nil || n <= 0 {
			return cfg, fmt.Errorf("container: timeout_ms must be a positive integer")
		}
		cfg.Timeout = time.Duration(n) * time.Millisecond
	}

	if v, ok := config["ready_timeout_ms"]; ok && v != nil {
		n, err := asInt(v)
		if err != nil || n <= 0 {
			return cfg, fmt.Errorf("container: ready_timeout_ms must be a positive integer")
		}
		cfg.ReadyTimeout = time.Duration(n) * time.Millisecond
	}

	return cfg, nil
}

func asInt(v any) (int, error) {
	switch n := v.(type) {
	case int:
		return n, nil
	case int64:
		return int(n), nil
	case float64:
		return int(n), nil
	case json.Number:
		i, err := n.Int64()
		return int(i), err
	default:
		return 0, fmt.Errorf("not an integer")
	}
}

func dockerName(realm, block string) string {
	san := func(s string) string {
		s = strings.ToLower(s)
		var b strings.Builder
		for _, r := range s {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
				b.WriteRune(r)
			} else {
				b.WriteByte('-')
			}
		}
		out := b.String()
		if len(out) > 20 {
			out = out[:20]
		}
		if out == "" {
			out = "x"
		}
		return out
	}
	// Docker name max 255; keep short + unique suffix.
	return fmt.Sprintf("astrate-%s-%s-%d", san(realm), san(block), time.Now().UnixNano()%1_000_000_000)
}

// Name implements flow.Block.
func (b *Block) Name() string { return b.name }

// Process implements flow.Block. PoC error policy: return error so the router
// logs + increments blockErrors and drops the message (no outs).
func (b *Block) Process(msg *flow.Message) ([]*flow.Message, error) {
	b.mu.Lock()
	stopped := b.stopped
	bridge := b.bridge
	b.mu.Unlock()
	if stopped || bridge == nil {
		return nil, fmt.Errorf("container stopped")
	}
	outs, err := bridge.RoundTrip(context.Background(), msg)
	if err != nil {
		// Loud enough for operators; router also counts blockErrors.
		if b.log != nil {
			b.log.Error("container process error", "block", b.name, "err", err)
		}
		return nil, err
	}
	return outs, nil
}

// Stop implements flow.Stopper: stop and remove the container.
func (b *Block) Stop() {
	b.mu.Lock()
	if b.stopped {
		b.mu.Unlock()
		return
	}
	b.stopped = true
	inst := b.inst
	b.inst = nil
	b.mu.Unlock()

	if inst == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := inst.Stop(ctx); err != nil && b.log != nil {
		b.log.Error("container stop failed", "block", b.name, "err", err)
	}
}

var (
	_ flow.Block   = (*Block)(nil)
	_ flow.Stopper = (*Block)(nil)
)
