--ALTER TABLE wallets ADD FOREIGN KEY (owner) REFERENCES accounts (username);

/*
 * 1 account can have multiple wallets, but those wallets should have different currencies. So, 
 * 2 wallets under the same user, both can't have 'USD' currencies.
 * unique index is set to 'wallets' table use to solve that problem.
 */
--CREATE UNIQUE INDEX ON wallets (owner, currency);
--CREATE UNIQUE INDEX ON wallets (currency);

/*
 * To prevent multiple wallets under the same user have same currrencies is to add
 * a unique constraint for the pair of 'owner' and 'currency' on the account table.
 */
--ALTER TABLE wallets ADD CONSTRAINT account_currency_key UNIQUE (account_id, currency);
ALTER TABLE wallets ADD CONSTRAINT account_name_key UNIQUE (account_id, name);

--ALTER TABLE entries ADD FOREIGN KEY (account_id) REFERENCES accounts(id);
--ALTER TABLE transfers ADD FOREIGN KEY (account_id) REFERENCES accounts(id);

CREATE TABLE addresses (
    id INT GENERATED ALWAYS AS IDENTITY CONSTRAINT pk_addresses_id PRIMARY KEY,
    provinces VARCHAR NOT NULL CONSTRAINT ck_addresses_province_empty CHECK (provinces <>''),
    city VARCHAR NOT NULL CONSTRAINT ck_addresses_city_empty CHECK (city <>''),
    zip INT NOT NULL CONSTRAINT ck_zip_zero CHECK (zip > 0),
    street TEXT NOT NULL CONSTRAINT ck_street_zero CHECK (street <> '')
);

ALTER TABLE accounts ADD CONSTRAINT fk_accounts_address FOREIGN KEY (address) REFERENCES addresses(id);

ALTER TABLE accounts ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE wallets ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE entries ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE transfers ADD COLUMN deleted_at TIMESTAMPTZ;