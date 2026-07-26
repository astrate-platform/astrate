-- 000005 down: unregister the TTL job and drop the procedure.
DO $$
DECLARE
    jid integer;
BEGIN
    FOR jid IN
        SELECT job_id FROM timescaledb_information.jobs
        WHERE proc_name = 'astrate_apply_endpoint_ttl'
    LOOP
        PERFORM delete_job(jid);
    END LOOP;
END
$$;
DROP PROCEDURE IF EXISTS astrate_apply_endpoint_ttl(integer, jsonb);
