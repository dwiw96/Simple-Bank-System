package services

import (
	"context"
	"errors"

	"log"
	"simple-bank-system/db/pkg"
	"simple-bank-system/util"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type CreateAccountParams struct {
	Username       string
	HashedPassword string
	FullName       string
	DateOfBirth    time.Time
	Address        pkg.Addresses
	Email          string
}

func (r *DB) CreateAccount(ctx context.Context, account CreateAccountParams) (*pkg.Account, error) {
	var res pkg.Account

	hashedPass, err := util.HashingPassword(account.HashedPassword)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO addresses(provinces, city, zip, street
		) VALUES(
			$1, $2, $3, $4
		) RETURNING id, provinces, city, zip, street;`
	dbReturn := r.db.QueryRow(ctx, query, account.Address.Provinces, account.Address.City, account.Address.ZIP, account.Address.Street)
	err = dbReturn.Scan(&res.Address.ID, &res.Address.Provinces, &res.Address.City, &res.Address.ZIP, &res.Address.Street)
	if err != nil {
		log.Println("--- (1)database")
		err = addressErrHandling(err)
		return nil, err
	}

	query = `INSERT INTO accounts(account_number, username, hashed_password, full_name, date_of_birth, address, email
		) VALUES(
			1010000000+CAST(1000000 + floor(random() * 9000000) AS bigint), $1, $2, $3, $4, $5, $6
		) RETURNING id, account_number, username, full_name, date_of_birth, email, password_change_at, created_at;`
	dbReturn = r.db.QueryRow(ctx, query, account.Username, hashedPass, account.FullName, account.DateOfBirth, res.Address.ID, account.Email)
	err = dbReturn.Scan(&res.ID, &res.AccountNumber, &res.Username, &res.FullName, &res.DateOfBirth, &res.Email, &res.PasswordChangeAt, &res.CreatedAt)
	if err != nil {
		log.Println("--- (2)database")
		err = accErrHandling(err)
		return nil, err
	}

	return &res, nil
}

func (r *DB) GetAccount(ctx context.Context, username string) (*pkg.Account, error) {
	/*
	 * Inner join using the WHERE clause.
	 * To use the WHERE clause to perform the same join as you perform using the INNER JOIN syntax, enter both the join condition and the additional
	 * selection condition in the WHERE clause.
	 * The tables to be joined are listed in the FROM clause, separated by commas.
	 */
	query := `SELECT * FROM accounts JOIN addresses ON accounts.address = addresses.id WHERE accounts.username=$1 AND deleted_at IS NULL;`
	row := r.db.QueryRow(ctx, query, &username)

	var account pkg.Account
	err := row.Scan(&account.ID, &account.AccountNumber, &account.Username, &account.HashedPassword, &account.FullName, &account.DateOfBirth, &account.Address.ID, &account.Email, &account.PasswordChangeAt, &account.CreatedAt,
		&account.DeletedAt, &account.Address.ID, &account.Address.Provinces, &account.Address.City, &account.Address.ZIP, &account.Address.Street)
	if err == pgx.ErrNoRows {
		return nil, util.ErrNotExist
	}
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// Get account_number using account_id
func (r *DB) GetAccountByID(ctx context.Context, id int64) (accountNumber *int64, err error) {
	/*
	 * Inner join using the WHERE clause.
	 * To use the WHERE clause to perform the same join as you perform using the INNER JOIN syntax, enter both the join condition and the additional
	 * selection condition in the WHERE clause.
	 * The tables to be joined are listed in the FROM clause, separated by commas.
	 */
	query := `SELECT account_number FROM accounts WHERE id=$1 AND deleted_at IS NULL;`
	row := r.db.QueryRow(ctx, query, &id)

	//var accountNumber int64
	err = row.Scan(&accountNumber)
	if err == pgx.ErrNoRows {
		return nil, util.ErrNotExist
	}
	if err != nil {
		return nil, err
	}
	return accountNumber, nil
}

func (r *DB) GetAccountByNumber(ctx context.Context, accountNumber int64) (*pkg.Account, error) {
	query := `SELECT * FROM accounts JOIN addresses ON accounts.address = addresses.id WHERE accounts.account_number=$1 AND deleted_at IS NULL;`
	row := r.db.QueryRow(ctx, query, &accountNumber)

	var account pkg.Account
	err := row.Scan(&account.ID, &account.AccountNumber, &account.Username, &account.HashedPassword, &account.FullName, &account.DateOfBirth, &account.Address.ID, &account.Email, &account.PasswordChangeAt, &account.CreatedAt,
		&account.DeletedAt, &account.Address.ID, &account.Address.Provinces, &account.Address.City, &account.Address.ZIP, &account.Address.Street)
	if err == pgx.ErrNoRows {
		return nil, util.ErrNotExist
	}
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func addressErrHandling(err error) error {
	//fmt.Println("Error Handling")
	var pgxError *pgconn.PgError
	if errors.As(err, &pgxError) {
		addressErrMsg := []string{"ck_addresses_province_empty", "ck_addresses_city_empty", "ck_zip_zero"}
		for i := range addressErrMsg {
			//fmt.Println("Error Handling Address, i =", i)
			if pgxError.ConstraintName == addressErrMsg[i] {
				//fmt.Println("Error Handling Address, constraint =", pgxError.ConstraintName)
				//fmt.Println("Error Handling Address, Msg =", addressErrMsg[i])
				return util.ErrAddressEmpty
			}
		}
	}
	return err
}

func accErrHandling(err error) error {
	//fmt.Println("Error Handling")
	var pgxError *pgconn.PgError
	if errors.As(err, &pgxError) {
		if pgxError.ConstraintName == "uq_accounts_username" {
			return util.ErrUsernameExists
		}
		if pgxError.ConstraintName == "ck_accounts_username_length" {
			return util.ErrUsernameEmpty
		}
		if pgxError.ConstraintName == "uq_accounts_accountNumber" {
			return util.ErrAccountNumberExists
		}
		if pgxError.ConstraintName == "ck_accounts_accountNumber_range" {
			return util.ErrAccountNumberWrong
		}
		if pgxError.ConstraintName == "ck_accounts_password_empty" {
			return util.ErrPasswordEmpty
		}
		if pgxError.ConstraintName == "ck_accounts_fullname_empty" {
			return util.ErrFullnameEmpty
		}
		if pgxError.ConstraintName == "ck_accounts_dob_empty" {
			return util.ErrDOBEmpty
		}
		if pgxError.ConstraintName == "uq_accounts_email" {
			return util.ErrEmailExists
		}
		if pgxError.ConstraintName == "ck_accounts_email_empty" {
			return util.ErrEmailEmpty
		}
		if pgxError.ConstraintName == "ck_accounts_address_zero" {
			return util.ErrAddressEmpty
		}
	}
	return err
}
