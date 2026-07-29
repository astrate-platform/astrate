package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// ErrPipelineCyclic reports that a pipeline definition's block graph
// contains a cycle.
var ErrPipelineCyclic = errors.New("store: pipeline graph contains a cycle")

// Pipeline is one stored Flow pipeline (astarte_flow parity, issue #24): a
// named, realm-scoped DAG of blocks. Definition is the raw pipeline JSON
// (internal/flow.Pipeline shape: "blocks" + "connections"); this package
// does not import internal/flow (store imports none of its callers, see the
// package doc), so the graph shape is decoded locally just far enough to
// validate it.
type Pipeline struct {
	ID         int64
	RealmID    int16
	Name       string
	Definition []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// pipelineGraph is the subset of internal/flow.Pipeline's JSON shape needed
// to validate block references and acyclicity.
type pipelineGraph struct {
	Blocks []struct {
		Name string `json:"name"`
	} `json:"blocks"`
	Connections []struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"connections"`
}

// validatePipelineGraph checks that every connection references a declared
// block and that the block graph is acyclic (Kahn's algorithm).
func validatePipelineGraph(definition []byte) error {
	var g pipelineGraph
	if err := json.Unmarshal(definition, &g); err != nil {
		return fmt.Errorf("store: pipeline definition does not parse: %w", err)
	}
	if len(g.Blocks) == 0 {
		return fmt.Errorf("store: pipeline has no blocks")
	}

	names := make(map[string]bool, len(g.Blocks))
	for _, b := range g.Blocks {
		if b.Name == "" {
			return fmt.Errorf("store: pipeline has a block with empty name")
		}
		if names[b.Name] {
			return fmt.Errorf("store: pipeline has duplicate block name %q", b.Name)
		}
		names[b.Name] = true
	}

	inDeg := make(map[string]int, len(g.Blocks))
	adj := make(map[string][]string, len(g.Blocks))
	for _, c := range g.Connections {
		if !names[c.From] || !names[c.To] {
			return fmt.Errorf("store: pipeline connection references unknown block %q -> %q", c.From, c.To)
		}
		inDeg[c.To]++
		adj[c.From] = append(adj[c.From], c.To)
	}

	var queue []string
	for name := range names {
		if inDeg[name] == 0 {
			queue = append(queue, name)
		}
	}
	visited := 0
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		visited++
		for _, next := range adj[n] {
			inDeg[next]--
			if inDeg[next] == 0 {
				queue = append(queue, next)
			}
		}
	}
	if visited != len(names) {
		return ErrPipelineCyclic
	}
	return nil
}

// CreatePipeline validates and stores a pipeline; a duplicate name yields
// ErrAlreadyExists, an invalid graph yields ErrPipelineCyclic or a plain
// error for an unresolved block reference.
func (s *Store) CreatePipeline(ctx context.Context, realmID int16, name string, definition []byte) (*Pipeline, error) {
	if err := validatePipelineGraph(definition); err != nil {
		return nil, err
	}
	p := Pipeline{RealmID: realmID, Name: name, Definition: definition}
	err := s.pool.QueryRow(ctx,
		`INSERT INTO pipelines (realm_id, name, definition) VALUES ($1, $2, $3)
		 RETURNING id, created_at, updated_at`,
		realmID, name, definition).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if pgErrCode(err) == pgCodeUniqueViolation {
		return nil, fmt.Errorf("%w: pipeline %q", ErrAlreadyExists, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: creating pipeline %q: %w", name, err)
	}
	return &p, nil
}

// GetPipeline fetches one pipeline by name.
func (s *Store) GetPipeline(ctx context.Context, realmID int16, name string) (*Pipeline, error) {
	var p Pipeline
	err := s.pool.QueryRow(ctx,
		`SELECT id, realm_id, name, definition, created_at, updated_at
		 FROM pipelines WHERE realm_id = $1 AND name = $2`,
		realmID, name).Scan(&p.ID, &p.RealmID, &p.Name, &p.Definition, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: pipeline %q", ErrNotFound, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: getting pipeline %q: %w", name, err)
	}
	return &p, nil
}

// UpdatePipeline validates and replaces a pipeline's definition.
func (s *Store) UpdatePipeline(ctx context.Context, realmID int16, name string, definition []byte) (*Pipeline, error) {
	if err := validatePipelineGraph(definition); err != nil {
		return nil, err
	}
	var p Pipeline
	err := s.pool.QueryRow(ctx,
		`UPDATE pipelines SET definition = $3, updated_at = now()
		 WHERE realm_id = $1 AND name = $2
		 RETURNING id, realm_id, name, definition, created_at, updated_at`,
		realmID, name, definition).Scan(&p.ID, &p.RealmID, &p.Name, &p.Definition, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: pipeline %q", ErrNotFound, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: updating pipeline %q: %w", name, err)
	}
	return &p, nil
}

// DeletePipeline removes a pipeline.
func (s *Store) DeletePipeline(ctx context.Context, realmID int16, name string) error {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM pipelines WHERE realm_id = $1 AND name = $2`, realmID, name)
	if err != nil {
		return fmt.Errorf("store: deleting pipeline %q: %w", name, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: pipeline %q", ErrNotFound, name)
	}
	return nil
}

// ListPipelines returns every pipeline of a realm ordered by name.
func (s *Store) ListPipelines(ctx context.Context, realmID int16) ([]Pipeline, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, realm_id, name, definition, created_at, updated_at
		 FROM pipelines WHERE realm_id = $1 ORDER BY name`,
		realmID)
	if err != nil {
		return nil, fmt.Errorf("store: listing pipelines: %w", err)
	}
	defer rows.Close()

	var out []Pipeline
	for rows.Next() {
		var p Pipeline
		if err := rows.Scan(&p.ID, &p.RealmID, &p.Name, &p.Definition, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("store: scanning pipeline: %w", err)
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: listing pipelines: %w", err)
	}
	return out, nil
}
