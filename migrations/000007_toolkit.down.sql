-- 000007 down: drop the toolkit if this migration installed it. Safe when the
-- extension was never created, and the lttb path degrades to time_bucket+avg
-- rather than failing.
DROP EXTENSION IF EXISTS timescaledb_toolkit;
