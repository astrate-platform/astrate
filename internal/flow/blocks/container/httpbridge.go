package container

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/astrate-platform/astrate/internal/flow"
)

const (
	messagePath         = "/v1/message"
	healthPath          = "/healthz"
	defaultTimeout      = 5 * time.Second
	defaultMaxBodyBytes = 1 << 20 // 1 MiB
	defaultReadyWait    = 15 * time.Second
)

// Bridge POSTs FlowMessages to a container HTTP endpoint.
type Bridge struct {
	Client       *http.Client
	BaseURL      string
	Timeout      time.Duration
	MaxBodyBytes int64
}

func (b *Bridge) client() *http.Client {
	if b.Client != nil {
		return b.Client
	}
	return http.DefaultClient
}

func (b *Bridge) timeout() time.Duration {
	if b.Timeout > 0 {
		return b.Timeout
	}
	return defaultTimeout
}

func (b *Bridge) maxBody() int64 {
	if b.MaxBodyBytes > 0 {
		return b.MaxBodyBytes
	}
	return defaultMaxBodyBytes
}

// WaitReady polls GET /healthz until success or ctx is done.
func (b *Bridge) WaitReady(ctx context.Context) error {
	url := strings.TrimRight(b.BaseURL, "/") + healthPath
	deadline, ok := ctx.Deadline()
	if !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultReadyWait)
		defer cancel()
		deadline = time.Now().Add(defaultReadyWait)
	}
	var last error
	for time.Now().Before(deadline) {
		if err := ctx.Err(); err != nil {
			if last != nil {
				return fmt.Errorf("container: not ready: %w (last: %v)", err, last)
			}
			return fmt.Errorf("container: not ready: %w", err)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		resp, err := b.client().Do(req)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return nil
			}
			last = fmt.Errorf("status %d", resp.StatusCode)
		} else {
			last = err
		}
		select {
		case <-ctx.Done():
			if last != nil {
				return fmt.Errorf("container: not ready: %w (last: %v)", ctx.Err(), last)
			}
			return fmt.Errorf("container: not ready: %w", ctx.Err())
		case <-time.After(100 * time.Millisecond):
		}
	}
	if last != nil {
		return fmt.Errorf("container: not ready after wait: %v", last)
	}
	return fmt.Errorf("container: not ready after wait")
}

// RoundTrip sends one Message and parses zero or more outputs.
//
// Contract (PoC):
//   - Request: POST /v1/message, body = Message JSON
//   - 204 or empty body → drop (zero outs)
//   - 200 + object → one message
//   - 200 + array → N messages
//   - other status → error
func (b *Bridge) RoundTrip(ctx context.Context, msg *flow.Message) ([]*flow.Message, error) {
	if msg == nil {
		return nil, nil
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("container: marshal message: %w", err)
	}

	url := strings.TrimRight(b.BaseURL, "/") + messagePath
	reqCtx, cancel := context.WithTimeout(ctx, b.timeout())
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := b.client().Do(req)
	if err != nil {
		return nil, fmt.Errorf("container: http post: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	limited := io.LimitReader(resp.Body, b.maxBody()+1)
	respBody, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("container: read response: %w", err)
	}
	if int64(len(respBody)) > b.maxBody() {
		return nil, fmt.Errorf("container: response body exceeds %d bytes", b.maxBody())
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet := strings.TrimSpace(string(respBody))
		if len(snippet) > 200 {
			snippet = snippet[:200] + "…"
		}
		return nil, fmt.Errorf("container: http %d: %s", resp.StatusCode, snippet)
	}

	respBody = bytes.TrimSpace(respBody)
	if len(respBody) == 0 {
		return nil, nil
	}

	// Array of messages?
	if respBody[0] == '[' {
		var list []*flow.Message
		if err := json.Unmarshal(respBody, &list); err != nil {
			return nil, fmt.Errorf("container: decode message array: %w", err)
		}
		out := make([]*flow.Message, 0, len(list))
		for _, m := range list {
			if m != nil {
				out = append(out, m)
			}
		}
		return out, nil
	}

	var one flow.Message
	if err := json.Unmarshal(respBody, &one); err != nil {
		return nil, fmt.Errorf("container: decode message: %w", err)
	}
	return []*flow.Message{&one}, nil
}
