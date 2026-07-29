package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Flow is one durable named instance of a pipeline for a realm (issues #40 +
// #41). PipelineName is resolved at start time against the pipelines table;
// Config is the JSON object used for ${config.*} substitution.
type Flow struct {
	ID           int64
	RealmID      int16
	Name         string
	PipelineName string
	Config       []byte
	AutoRestart  bool
	Status       string
	ErrorMessage *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	StartedAt    *time.Time
	StoppedAt    *time.Time
}

// FlowRehydrate is a durable auto_restart row plus its realm name for boot.
type FlowRehydrate struct {
	Flow
	RealmName string
}

const flowColumns = `id, realm_id, name, pipeline_name, config, auto_restart,
	status, error_message, created_at, updated_at, started_at, stopped_at`

func scanFlow(scan interface {
	Scan(dest ...any) error
}) (*Flow, error) {
	var f Flow
	err := scan.Scan(
		&f.ID, &f.RealmID, &f.Name, &f.PipelineName, &f.Config, &f.AutoRestart,
		&f.Status, &f.ErrorMessage, &f.CreatedAt, &f.UpdatedAt, &f.StartedAt, &f.StoppedAt,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// CreateFlow inserts a durable flow row. config defaults to {} when empty.
// autoRestart is stored as given (service layer defaults omitted to true).
func (s *Store) CreateFlow(ctx context.Context, realmID int16, name, pipelineName string, config []byte, autoRestart bool) (*Flow, error) {
	if len(config) == 0 {
		config = []byte("{}")
	}
	row := s.pool.QueryRow(ctx,
		`INSERT INTO flows (realm_id, name, pipeline_name, config, auto_restart)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING `+flowColumns,
		realmID, name, pipelineName, config, autoRestart)
	f, err := scanFlow(row)
	if pgErrCode(err) == pgCodeUniqueViolation {
		return nil, fmt.Errorf("%w: flow %q", ErrAlreadyExists, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: creating flow %q: %w", name, err)
	}
	return f, nil
}

// GetFlow fetches one durable flow by realm and name.
func (s *Store) GetFlow(ctx context.Context, realmID int16, name string) (*Flow, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT `+flowColumns+` FROM flows WHERE realm_id = $1 AND name = $2`,
		realmID, name)
	f, err := scanFlow(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%w: flow %q", ErrNotFound, name)
	}
	if err != nil {
		return nil, fmt.Errorf("store: getting flow %q: %w", name, err)
	}
	return f, nil
}

// ListFlows returns every durable flow for a realm ordered by name.
func (s *Store) ListFlows(ctx context.Context, realmID int16) ([]Flow, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT `+flowColumns+` FROM flows WHERE realm_id = $1 ORDER BY name`,
		realmID)
	if err != nil {
		return nil, fmt.Errorf("store: listing flows: %w", err)
	}
	defer rows.Close()

	var out []Flow
	for rows.Next() {
		f, err := scanFlow(rows)
		if err != nil {
			return nil, fmt.Errorf("store: scanning flow: %w", err)
		}
		out = append(out, *f)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: listing flows: %w", err)
	}
	return out, nil
}

// ListAutoRestartFlows returns every flow with auto_restart=true across all
// realms, for process boot rehydrate.
func (s *Store) ListAutoRestartFlows(ctx context.Context) ([]FlowRehydrate, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT f.id, f.realm_id, f.name, f.pipeline_name, f.config, f.auto_restart,
		        f.status, f.error_message, f.created_at, f.updated_at, f.started_at, f.stopped_at,
		        r.name
		 FROM flows f
		 JOIN realms r ON r.id = f.realm_id
		 WHERE f.auto_restart = true
		 ORDER BY r.name, f.name`)
	if err != nil {
		return nil, fmt.Errorf("store: listing auto-restart flows: %w", err)
	}
	defer rows.Close()

	var out []FlowRehydrate
	for rows.Next() {
		var fr FlowRehydrate
		err := rows.Scan(
			&fr.ID, &fr.RealmID, &fr.Name, &fr.PipelineName, &fr.Config, &fr.AutoRestart,
			&fr.Status, &fr.ErrorMessage, &fr.CreatedAt, &fr.UpdatedAt, &fr.StartedAt, &fr.StoppedAt,
			&fr.RealmName,
		)
		if err != nil {
			return nil, fmt.Errorf("store: scanning auto-restart flow: %w", err)
		}
		out = append(out, fr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: listing auto-restart flows: %w", err)
	}
	return out, nil
}

// UpdateFlowRuntime persists status / error / timestamps after start, stop, or fail.
func (s *Store) UpdateFlowRuntime(ctx context.Context, realmID int16, name, status string, errMsg *string, startedAt, stoppedAt *time.Time) error {
	tag, err := s.pool.Exec(ctx,
		`UPDATE flows SET
		   status = $3,
		   error_message = $4,
		   started_at = COALESCE($5, started_at),
		   stopped_at = $6,
		   updated_at = now()
		 WHERE realm_id = $1 AND name = $2`,
		realmID, name, status, errMsg, startedAt, stoppedAt)
	if err != nil {
		return fmt.Errorf("store: updating flow runtime %q: %w", name, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: flow %q", ErrNotFound, name)
	}
	return nil
}

// DeleteFlow removes a durable flow row.
func (s *Store) DeleteFlow(ctx context.Context, realmID int16, name string) error {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM flows WHERE realm_id = $1 AND name = $2`, realmID, name)
	if err != nil {
		return fmt.Errorf("store: deleting flow %q: %w", name, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("%w: flow %q", ErrNotFound, name)
	}
	return nil
}
