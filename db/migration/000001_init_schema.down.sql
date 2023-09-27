--Order position is important
--ALTER TABLE IF EXISTS wallets DROP CONSTRAINT IF EXISTS accounts_owner_currency_idx ;
--ALTER TABLE IF EXISTS wallets DROP CONSTRAINT IF EXISTS accounts_owner_fkey;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS schema_migrations;