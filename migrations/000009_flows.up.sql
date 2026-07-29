-- 000009: durable Flow instances (issues #40 + #41).
-- A flow is a named specialization of a stored pipeline for a realm.

CREATE TABLE flows (
    id              bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    realm_id        smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    name            text NOT NULL,
    pipeline_name   text NOT NULL,
    config          jsonb NOT NULL DEFAULT '{}'::jsonb,
    auto_restart    boolean NOT NULL DEFAULT true,
    -- desired/runtime observation (single process; last write wins)
    status          text NOT NULL DEFAULT 'stopped',
    error_message   text,
    created_at      timestamptz NOT NULL DEFAULT now(),
    updated_at      timestamptz NOT NULL DEFAULT now(),
    started_at      timestamptz,
    stopped_at      timestamptz,
    UNIQUE (realm_id, name)
);

CREATE INDEX flows_realm_auto_restart_idx
    ON flows (realm_id)
    WHERE auto_restart = true;
