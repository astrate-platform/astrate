-- 000007: enable timescaledb_toolkit where the server actually ships it.
--
-- Deliberately conditional rather than a plain CREATE EXTENSION. The toolkit
-- is an *optional* capability: Store.probeToolkit reads pg_extension at
-- startup and Downsample falls back to time_bucket+avg when it is absent
-- (docs/DESIGN.md §2.5), and that fallback is a load-bearing invariant — a
-- deployment on a server without the toolkit must still migrate cleanly.
-- Creating it unconditionally would turn an optional path into a hard
-- requirement for every operator.
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_available_extensions WHERE name = 'timescaledb_toolkit') THEN
        CREATE EXTENSION IF NOT EXISTS timescaledb_toolkit;
    END IF;
END
$$;
