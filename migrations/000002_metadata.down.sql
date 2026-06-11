-- 000002 down: drop the metadata tables in reverse dependency order.
DROP TABLE IF EXISTS triggers;
DROP TABLE IF EXISTS group_devices;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS devices;
DROP TABLE IF EXISTS endpoints;
DROP TABLE IF EXISTS interfaces;
DROP TABLE IF EXISTS realms;
