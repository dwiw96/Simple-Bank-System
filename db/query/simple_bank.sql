CREATE TABLE users (
    username VARCHAR PRIMARY KEY REFERENCES accounts(owner),
    hashed_password VARCHAR NOT NULL,
    full_name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password_change_at TIMESTAMPTZ DEFAULT '0001-01-01 00:00:00Z' NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

/*
 * docker cp simple_bank.sql postgres_sql:/simple_bank.sql
 * docker exec -t postgres_sql psql -U db db -f simple_bank.sql
 */

CREATE TYPE valid_currency AS ENUM ('IDR', 'USD', 'EUR');
CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    owner VARCHAR REFERENCES users(username) NOT NULL, check (owner <> ''),
    balance BIGINT NOT NULL, check (balance >= 0),
    currency valid_currency NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
CREATE UNIQUE INDEX ON accounts (owner, currency;);

/* 
 * This table will record all change to the account balance.
 * This table also represent 1-to-many relationship between
 * accounts and entries
 */
CREATE TABLE entries (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT references accounts(id) NOT NULL, check (account_id > 0),
    amount BIGINT NOT NULL, -- added money to the account balance in thi entries
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE transfers (
    id BIGSERIAL PRIMARY KEY,
    from_account_id BIGINT REFERENCES accounts(id) NOT NULL, check (from_account_id > 0),
    to_account_id BIGINT REFERENCES accounts(id) NOT NULL, check (to_account_id > 0),
    amount BIGINT NOT NULL, check (amount >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX accounts_owner_index ON accounts (owner);
CREATE INDEX entries_account_id_index ON entries (account_id);
CREATE INDEX ON transfers (from_account_id);
CREATE INDEX ON transfers (to_account_id);
CREATE INDEX ON transfers (from_account_id, to_account_id);

COMMENT ON COLUMN entries.amount IS 'can be negative or positive';
COMMENT ON COLUMN transfers.amount IS 'must be postive';