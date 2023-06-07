CREATE TABLE users (
    username VARCHAR PRIMARY KEY,
    hashed_password VARCHAR NOT NULL,
    full_name VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password_change_at TIMESTAMPTZ DEFAULT '0001-01-01 00:00:00Z' NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

ALTER TABLE accounts ADD FOREIGN KEY (owner) REFERENCES users (username);

/*
 * 1 user can have multiple accounts, but those accounts should have different currencies. So, 
 * 2 accounts under the same user, both can't have 'USD' currencies.
 * unique index is set to 'accounts' table use to solve that problem.
 */
CREATE UNIQUE INDEX ON accounts (owner, currency);

/*
 * Another way to prevent multiple accounts under the same user have same currrencies is to add
 * a unique constraint for the pair of 'owner' and 'currency' on the account table.
 */

--ALTER TABLE accounts ADD CONSTRAINT owner_currency_key UNIQUE (owner, currency);