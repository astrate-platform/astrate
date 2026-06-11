-- 000005: per-endpoint TTL enforcement (docs/DESIGN.md §2.5, docs/ROADMAP.md
-- §3.1 file 2.5). A TimescaleDB user-defined action deletes datastream rows
-- older than their endpoint's database_retention_ttl. Deletes target one chunk
-- at a time and commit between chunks to bound transaction size. The optional
-- global hard-cap retention (retention.max_days → add_retention_policy) is
-- applied at runtime from config, not here.
--
-- The (job_id, config) parameters match the signature TimescaleDB requires of
-- job procedures; both default to NULL so `CALL astrate_apply_endpoint_ttl()`
-- also works for tests and operators.
CREATE PROCEDURE astrate_apply_endpoint_ttl(job_id integer DEFAULT NULL, config jsonb DEFAULT NULL)
LANGUAGE plpgsql
AS $proc$
DECLARE
    ep     record;
    chunk  regclass;
    cutoff timestamptz;
BEGIN
    -- Individual datastreams: TTL is declared per endpoint.
    FOR ep IN
        SELECT e.interface_id, e.id AS endpoint_id, e.database_retention_ttl AS ttl
        FROM endpoints e
        WHERE e.database_retention_policy = 'use_ttl'
          AND e.database_retention_ttl IS NOT NULL
    LOOP
        cutoff := now() - make_interval(secs => ep.ttl);
        FOR chunk IN
            SELECT format('%I.%I', chunk_schema, chunk_name)::regclass
            FROM timescaledb_information.chunks
            WHERE hypertable_schema = 'public'
              AND hypertable_name = 'individual_datastreams'
              AND range_start < cutoff
        LOOP
            EXECUTE format(
                'DELETE FROM %s WHERE interface_id = $1 AND endpoint_id = $2 AND ts < $3',
                chunk)
            USING ep.interface_id, ep.endpoint_id, cutoff;
            COMMIT;
        END LOOP;
    END LOOP;

    -- Object datastreams: rows carry no endpoint_id. All mappings of an
    -- object-aggregated interface share the same retention attributes
    -- (pkg/interfaceschema rejects mixed attributes), so the interface-level
    -- minimum is the exact shared TTL.
    FOR ep IN
        SELECT e.interface_id, min(e.database_retention_ttl) AS ttl
        FROM endpoints e
        JOIN interfaces i ON i.id = e.interface_id
        WHERE i.aggregation = 'object'
          AND e.database_retention_policy = 'use_ttl'
          AND e.database_retention_ttl IS NOT NULL
        GROUP BY e.interface_id
    LOOP
        cutoff := now() - make_interval(secs => ep.ttl);
        FOR chunk IN
            SELECT format('%I.%I', chunk_schema, chunk_name)::regclass
            FROM timescaledb_information.chunks
            WHERE hypertable_schema = 'public'
              AND hypertable_name = 'object_datastreams'
              AND range_start < cutoff
        LOOP
            EXECUTE format(
                'DELETE FROM %s WHERE interface_id = $1 AND ts < $2',
                chunk)
            USING ep.interface_id, cutoff;
            COMMIT;
        END LOOP;
    END LOOP;
END;
$proc$;

SELECT add_job('astrate_apply_endpoint_ttl', '1 hour');
