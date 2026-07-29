-- 000008: Flow pipeline storage (issue #24, astarte_flow parity). A pipeline
-- is a JSON-serialisable DAG of named blocks (internal/flow.Pipeline shape);
-- graph validity (acyclicity, block-reference resolution) is enforced by
-- internal/store before the row is written, not by a DB constraint.

CREATE TABLE pipelines (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    realm_id   smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    name       text NOT NULL,
    definition jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (realm_id, name)
);
