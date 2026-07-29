package container

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Spec describes a container to start for a flow block.
type Spec struct {
	Image      string
	// ContainerPort is the HTTP port inside the container (default 8080).
	ContainerPort int
	// FlowConfigJSON is passed as ASTRATE_FLOW_CONFIG.
	FlowConfigJSON string
	// Labels are applied to the container (astrate.*).
	Labels map[string]string
	// Name is an optional docker container name (must be unique).
	Name string
}

// Instance is a running container that serves the HTTP bridge.
type Instance interface {
	// BaseURL is the host URL to POST messages to (e.g. http://127.0.0.1:49152).
	BaseURL() string
	// ID is the docker container id (or a test fake id).
	ID() string
	// Stop stops and removes the container (best-effort).
	Stop(ctx context.Context) error
}

// Runner starts containers. PoC default is CLIRunner (docker CLI).
type Runner interface {
	Start(ctx context.Context, spec Spec) (Instance, error)
}

// CLIRunner shells out to the docker CLI. Suitable for PoC; MVP may switch
// to the Engine API client for cancel/cleanup.
type CLIRunner struct {
	// Docker is the docker binary name or path (default "docker").
	Docker string
	// Run is optional; when set, tests can intercept command execution.
	// argv[0] is the docker binary. Returns stdout, stderr, error.
	Run func(ctx context.Context, name string, args ...string) (stdout, stderr string, err error)
}

func (r *CLIRunner) bin() string {
	if r != nil && r.Docker != "" {
		return r.Docker
	}
	return "docker"
}

func (r *CLIRunner) run(ctx context.Context, args ...string) (string, string, error) {
	if r != nil && r.Run != nil {
		return r.Run(ctx, r.bin(), args...)
	}
	cmd := exec.CommandContext(ctx, r.bin(), args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// Start runs `docker run -d` with a published localhost port and Astrate labels.
func (r *CLIRunner) Start(ctx context.Context, spec Spec) (Instance, error) {
	if spec.Image == "" {
		return nil, fmt.Errorf("container: image is required")
	}
	port := spec.ContainerPort
	if port <= 0 {
		port = defaultContainerPort
	}

	args := []string{
		"run", "-d",
		// Publish container port to an ephemeral port on loopback only.
		"-p", fmt.Sprintf("127.0.0.1::%d", port),
	}
	if spec.Name != "" {
		args = append(args, "--name", spec.Name)
	}
	if spec.FlowConfigJSON != "" {
		args = append(args, "-e", "ASTRATE_FLOW_CONFIG="+spec.FlowConfigJSON)
	}
	for k, v := range spec.Labels {
		args = append(args, "--label", k+"="+v)
	}
	args = append(args, spec.Image)

	stdout, stderr, err := r.run(ctx, args...)
	if err != nil {
		msg := strings.TrimSpace(stderr)
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("container: docker run failed: %s", msg)
	}
	id := strings.TrimSpace(stdout)
	if id == "" {
		return nil, fmt.Errorf("container: docker run returned empty container id")
	}

	hostPort, err := r.hostPort(ctx, id, port)
	if err != nil {
		// Best-effort cleanup so we do not leave orphans on mapping failure.
		_, _, _ = r.run(context.Background(), "rm", "-f", id)
		return nil, err
	}

	return &cliInstance{
		runner: r,
		id:     id,
		base:   fmt.Sprintf("http://127.0.0.1:%s", hostPort),
	}, nil
}

func (r *CLIRunner) hostPort(ctx context.Context, id string, containerPort int) (string, error) {
	// docker port <id> 8080/tcp → "127.0.0.1:54321"
	stdout, stderr, err := r.run(ctx, "port", id, fmt.Sprintf("%d/tcp", containerPort))
	if err != nil {
		msg := strings.TrimSpace(stderr)
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("container: docker port failed: %s", msg)
	}
	line := strings.TrimSpace(stdout)
	// May be multi-line; take first non-empty.
	for _, l := range strings.Split(line, "\n") {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		// "0.0.0.0:12345" or "127.0.0.1:12345"
		if i := strings.LastIndex(l, ":"); i >= 0 && i+1 < len(l) {
			return l[i+1:], nil
		}
	}
	return "", fmt.Errorf("container: could not parse host port from %q", line)
}

type cliInstance struct {
	runner *CLIRunner
	id     string
	base   string
}

func (c *cliInstance) BaseURL() string { return c.base }
func (c *cliInstance) ID() string      { return c.id }

func (c *cliInstance) Stop(ctx context.Context) error {
	if c == nil || c.id == "" {
		return nil
	}
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
	}
	_, stderr, err := c.runner.run(ctx, "rm", "-f", c.id)
	if err != nil {
		msg := strings.TrimSpace(stderr)
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("container: docker rm -f %s: %s", c.id, msg)
	}
	return nil
}

// encodeFlowConfigJSON marshals nested config for ASTRATE_FLOW_CONFIG.
func encodeFlowConfigJSON(cfg map[string]any) (string, error) {
	if len(cfg) == 0 {
		return "{}", nil
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("container: marshal config: %w", err)
	}
	return string(b), nil
}
