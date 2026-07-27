//go:build nats

package forward

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

// NATSConfig is the NATS bus forwarder's configuration.
type NATSConfig struct {
	// URL is the NATS server address (e.g. "nats://host:4222"). Required.
	URL string
	// Subject is the subject every custom action is published to. Required.
	Subject string
}

// NATS publishes every custom trigger action to a single NATS subject.
type NATS struct {
	conn    *nats.Conn
	subject string
}

// NewNATS builds a NATS bus forwarder. It returns an error when URL or
// Subject is empty, or when the server cannot be reached.
func NewNATS(cfg NATSConfig) (*NATS, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("forward: URL must not be empty")
	}
	if cfg.Subject == "" {
		return nil, fmt.Errorf("forward: Subject must not be empty")
	}
	conn, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("forward: %w", err)
	}
	return &NATS{conn: conn, subject: cfg.Subject}, nil
}

// Forward implements the triggers.Forwarder contract.
func (n *NATS) Forward(_ context.Context, realm, trigger string, action json.RawMessage, event []byte) error {
	payload, err := marshalEnvelope(realm, trigger, action, event)
	if err != nil {
		return fmt.Errorf("forward: marshal: %w", err)
	}
	if err := n.conn.Publish(n.subject, payload); err != nil {
		return fmt.Errorf("forward: %w", err)
	}
	return nil
}
