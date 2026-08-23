-- 000012: per-realm user-defined composite blocks (issue #85, astarte_flow
-- parity). block_type carries upstream's enum (producer|consumer|
-- producer_consumer); source is the pipeline-definition body inlined at flow
-- start by internal/flow.ExpandComposites; config_schema is an optional JSON
-- Schema validating the params flows may configure on the block.

CREATE TABLE user_blocks (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    realm_id      smallint NOT NULL REFERENCES realms(id) ON DELETE CASCADE,
    name          text NOT NULL,
    block_type    text NOT NULL CHECK (block_type IN ('producer', 'consumer', 'producer_consumer')),
    source        text NOT NULL,
    config_schema text,
    created_at    timestamptz NOT NULL DEFAULT now(),
    updated_at    timestamptz NOT NULL DEFAULT now(),
    UNIQUE (realm_id, name)
);
