-- 000004: datastream hypertables, transcribed verbatim from docs/DESIGN.md §2.4
-- (docs/ROADMAP.md §3.1 file 2.4). The design's closing comment "same
-- compression policy shape as above" is materialized for object_datastreams.

CREATE TABLE individual_datastreams (
    realm_id      smallint NOT NULL,
    device_id     uuid NOT NULL,
    interface_id  bigint NOT NULL,
    endpoint_id   bigint NOT NULL,
    path          text NOT NULL,
    ts            timestamptz NOT NULL,     -- value timestamp (explicit_timestamp or reception)
    reception_ts  timestamptz NOT NULL,
    -- exactly one of the following is non-NULL, per the endpoint's declared type:
    value_double       double precision,
    value_integer      integer,
    value_longinteger  bigint,
    value_boolean      boolean,
    value_string       text,
    value_binaryblob   bytea,
    value_datetime     timestamptz,
    value_array        jsonb               -- all *array types (doublearray, stringarray, ...)
);
SELECT create_hypertable('individual_datastreams', by_range('ts', INTERVAL '7 days'));
CREATE INDEX ids_series_idx ON individual_datastreams
    (realm_id, device_id, interface_id, path, ts DESC);

ALTER TABLE individual_datastreams SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'realm_id, device_id, interface_id, path',
    timescaledb.compress_orderby   = 'ts DESC'
);
SELECT add_compression_policy('individual_datastreams', INTERVAL '7 days');

CREATE TABLE object_datastreams (
    realm_id      smallint NOT NULL,
    device_id     uuid NOT NULL,
    interface_id  bigint NOT NULL,
    path          text NOT NULL,             -- the parametric prefix, e.g. '/12'
    ts            timestamptz NOT NULL,
    reception_ts  timestamptz NOT NULL,
    value         jsonb NOT NULL             -- {"temp": 22.1, "hum": 41.0}
);
SELECT create_hypertable('object_datastreams', by_range('ts', INTERVAL '7 days'));
CREATE INDEX ods_series_idx ON object_datastreams
    (realm_id, device_id, interface_id, path, ts DESC);
-- same compression policy shape as above
ALTER TABLE object_datastreams SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'realm_id, device_id, interface_id, path',
    timescaledb.compress_orderby   = 'ts DESC'
);
SELECT add_compression_policy('object_datastreams', INTERVAL '7 days');
