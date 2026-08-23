-- 000011: adopt realms.datastream_maximum_storage_retention formally.
-- The column shipped earlier as idempotent startup DDL
-- (store.ensureRealmRetentionColumn) while migrations/ was frozen; this makes
-- fresh and upgraded databases converge on one definition. IF NOT EXISTS keeps
-- either application order safe.
ALTER TABLE realms ADD COLUMN IF NOT EXISTS datastream_maximum_storage_retention bigint;
