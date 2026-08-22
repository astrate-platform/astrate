-- 000010: record which block killed a flow (issue #45 phase 1).
-- Set when a block signals a fatal runtime failure (e.g. its container died
-- mid-run); cleared on any successful start/restart.
ALTER TABLE flows ADD COLUMN failed_block text;
