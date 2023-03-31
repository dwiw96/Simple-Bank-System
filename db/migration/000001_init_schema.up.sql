/*
 * docker cp simple_bank.sql postgres_sql:/simple_bank.sql
 * docker exec -t postgres_sql psql -U db db -f simple_bank.sql
 */

CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    owner VARCHAR NOT NULL,
    balance BIGINT NOT NULL,
    currency VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

/* 
 * This table will record all change to the account balance.
 * This table also represent 1-to-many relationship between
 * accounts and entries
 */
CREATE TABLE entries (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT references accounts(id) NOT NULL,
    amount BIGINT NOT NULL, -- added money to the account balance in thi entries
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE transfers (
    id BIGSERIAL PRIMARY KEY,
    from_account_id BIGINT REFERENCES accounts(id) NOT NULL,
    to_account_id BIGINT REFERENCES accounts(id) NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX accounts_owner_index ON accounts (owner);
CREATE INDEX entries_account_id_index ON entries (account_id);
CREATE INDEX ON transfers (from_account_id);
CREATE INDEX ON transfers (to_account_id);
CREATE INDEX ON transfers (from_account_id, to_account_id);

COMMENT ON COLUMN entries.amount IS 'can be negative or positive';
COMMENT ON COLUMN transfers.amount IS 'must be postive';