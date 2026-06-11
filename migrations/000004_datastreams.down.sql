-- 000004 down: drop the datastream hypertables. Dropping a hypertable also
-- removes its compression settings and policies (TimescaleDB jobs).
DROP TABLE IF EXISTS object_datastreams;
DROP TABLE IF EXISTS individual_datastreams;
