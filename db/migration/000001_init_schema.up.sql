/*
 * docker cp simple_bank.sql postgres_sql:/simple_bank.sql
 * docker exec -t postgres_sql psql -U db db -f simple_bank.sql
 */
CREATE TABLE accounts (
    id INTEGER GENERATED ALWAYS AS IDENTITY
        CONSTRAINT pk_accounts_id PRIMARY KEY,
    account_number BIGINT NOT NULL CONSTRAINT uq_accounts_accountNumber UNIQUE, 
        CONSTRAINT ck_accounts_accountNumber_range CHECK (account_number >= 1010000000 AND account_number <= 9999999999),
    username VARCHAR(20) NOT NULL CONSTRAINT uq_accounts_username  UNIQUE, 
        CONSTRAINT ck_accounts_username_length CHECK (LENGTH(username) > 3),
    hashed_password VARCHAR NOT NULL CONSTRAINT ck_accounts_password_empty CHECK (hashed_password <> ''),
    full_name VARCHAR NOT NULL CONSTRAINT ck_accounts_fullname_empty CHECK (full_name <> ''),
    date_of_birth DATE NOT NULL CONSTRAINT ck_accounts_dob_empty CHECK (date_of_birth > '1900-01-01'),
    address smallint NOT NULL CONSTRAINT ck_accounts_address_zero CHECK (address > 0),
    email VARCHAR CONSTRAINT uq_accounts_email UNIQUE NOT NULL CONSTRAINT ck_accounts_email_empty CHECK (email <> ''),
    password_change_at TIMESTAMPTZ DEFAULT '0001-01-01 00:00:00Z' NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TYPE valid_currency AS ENUM ('IDR', 'USD', 'EUR', 'YEN');
CREATE TABLE wallets (
    id INT GENERATED ALWAYS AS IDENTITY
        CONSTRAINT pk_wallets_id PRIMARY KEY,
    account_id INT REFERENCES accounts(id) NOT NULL, CONSTRAINT ck_wallets_accountID_zero CHECK(account_id > 0),
    wallet_number BIGINT CONSTRAINT uq_wallets_walletNumber UNIQUE NOT NULL,
        CONSTRAINT ck_wallets_walletNumber_range CHECK (wallet_number >= 1010000000 AND wallet_number <= 9999999999),
    name VARCHAR NOT NULL CONSTRAINT ck_wallets_name_empty CHECK (name <> ''),
    balance BIGINT NOT NULL, CONSTRAINT ck_wallets_balance_minus CHECK (balance >= 0),
    currency valid_currency NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

/* 
 * This table will record all change to the wallets balance.
 * This table also represent 1-to-many relationship between
 * wallets and entries
 */
CREATE TABLE entries (
    id BIGSERIAL CONSTRAINT pk_entries_id PRIMARY KEY,
    account_id INT NOT NULL,
        CONSTRAINT fk_entries_accountId FOREIGN KEY(account_id) REFERENCES accounts(id), 
            CONSTRAINT ck_entries_accountId_zero CHECK (account_id > 0),
    wallet_id INT NOT NULL,
        CONSTRAINT fk_entries_walletId FOREIGN KEY(wallet_id) REFERENCES wallets(id),
            CONSTRAINT ck_entries_walletId_zero CHECK (wallet_id > 0),
    wallet_number BIGINT NOT NULL,
        CONSTRAINT fk_entries_walletNumber FOREIGN KEY(wallet_number) REFERENCES wallets(wallet_number), 
            CONSTRAINT entries_walletNumber_zero CHECK (wallet_number > 0),
    amount BIGINT NOT NULL, -- added money to the wallets balance in thi entries
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE transfers (
    id BIGSERIAL CONSTRAINT pk_transfers_id PRIMARY KEY,
    account_id INT NOT NULL,
        CONSTRAINT fk_transfers_accountId FOREIGN KEY (account_id) REFERENCES accounts(id),
            CONSTRAINT ck_transfers_accountId_zero CHECK (account_id > 0),
    wallet_id BIGINT NOT NULL,
        CONSTRAINT fk_transfers_walletId FOREIGN KEY (wallet_id) REFERENCES wallets(id),
            CONSTRAINT ck_transfers_walletId_zero CHECK (wallet_id > 0),
    from_wallet_number BIGINT NOT NULL,
        CONSTRAINT fk_transfers_fromWalletNumber FOREIGN KEY (from_wallet_number) REFERENCES wallets(wallet_number),
            CONSTRAINT ck_transfers_fromWallet_zero CHECK (from_wallet_number > 0),
    to_wallet_number BIGINT NOT NULL,
        CONSTRAINT fk_transfers_toWalletNumber FOREIGN KEY (to_wallet_number) REFERENCES wallets(wallet_number),
            CONSTRAINT ck_transfers_toWallet_zero CHECK (to_wallet_number > 0),
    amount BIGINT NOT NULL, check (amount >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX ix_wallets_nam ON wallets (name);
CREATE INDEX ix_entries_walletsId ON entries (wallet_id);
CREATE INDEX ix_transfers_fromWalletNumber ON transfers (from_wallet_number);
CREATE INDEX ix_transfers_toWalletNumber ON transfers (to_wallet_number);
CREATE INDEX ix_transfers_fromWallet_toWallet ON transfers (from_wallet_number, to_wallet_number);
CREATE INDEX ix_accounts_accountNumber ON accounts(account_number);
CREATE INDEX ix_wallets_walletNumber ON wallets(wallet_number);

COMMENT ON COLUMN entries.amount IS 'can be negative or positive';
COMMENT ON COLUMN transfers.amount IS 'must be postive';