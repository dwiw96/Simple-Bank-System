package services

import (
	"context"
	"errors"
	"log"

	"simple-bank-system/db/pkg"
	"simple-bank-system/util"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type CreateWalletParams struct {
	WalletNumber int64
	Name         string
	AccountID    int64
	Balance      int64
	Currency     string
}

func (r *DB) CreateWallet(ctx context.Context, wallet CreateWalletParams) (*pkg.Wallet, error) {
	var res pkg.Wallet

	query := `INSERT INTO wallets(wallet_number, name, account_id, balance, currency
		) VALUES(
			$1+CAST(1000 + floor(random() * 9000) AS bigint), $2, $3, $4, $5
		) RETURNING id, name, account_id, wallet_number, balance, currency, created_at;`
	err := r.db.QueryRow(ctx, query, wallet.WalletNumber, wallet.Name, wallet.AccountID, wallet.Balance, wallet.Currency).Scan(&res.ID, &res.Name, &res.AccountID, &res.WalletNumber, &res.Balance, &res.Currency, &res.CreatedAt)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.ConstraintName == "account_name_key" {
				return nil, errors.New("wallet name already exists")
			}
			if pgxError.ConstraintName == "wallets_wallet_number_key" {
				for pgxError.ConstraintName == "wallets_wallet_number_key" {
					err = r.db.QueryRow(ctx, query, wallet.WalletNumber, wallet.Name, wallet.AccountID, wallet.Balance, wallet.Currency).Scan(&res.ID, &res.Name, &res.AccountID, &res.WalletNumber, &res.Balance, &res.Currency, &res.CreatedAt)
					errors.As(err, &pgxError)
				}
			}
			// 23503 (foreign_key_violation) -> Create wallet for unexists user
			if pgxError.Code == "23503" {
				return nil, util.ErrAccUser
			}
			// 23505 (unique_violation) ->  Duplicate wallet with the same currency
			if pgxError.Code == "23505" {
				return nil, util.ErrDuplicate
			}
		}
		return nil, err
	}

	return &res, nil
}

func (r *DB) CreatePrimaryWallet(ctx context.Context, wallet CreateWalletParams) (*pkg.Wallet, error) {
	var res pkg.Wallet

	query := `INSERT INTO wallets(wallet_number, name, account_id, balance, currency
		) VALUES(
			$1, $2, $3, $4, $5
		) RETURNING id, name, account_id, wallet_number, balance, currency, created_at;`
	err := r.db.QueryRow(ctx, query, wallet.WalletNumber, wallet.Name, wallet.AccountID, wallet.Balance, wallet.Currency).Scan(&res.ID, &res.Name, &res.AccountID, &res.WalletNumber, &res.Balance, &res.Currency, &res.CreatedAt)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.ConstraintName == "account_name_key" {
				return nil, errors.New("wallet name already exists")
			}
			if pgxError.ConstraintName == "wallets_wallet_number_key" {
				for pgxError.ConstraintName == "wallets_wallet_number_key" {
					err = r.db.QueryRow(ctx, query, wallet.WalletNumber, wallet.Name, wallet.AccountID, wallet.Balance, wallet.Currency).Scan(&res.ID, &res.Name, &res.AccountID, &res.WalletNumber, &res.Balance, &res.Currency, &res.CreatedAt)
					errors.As(err, &pgxError)
				}
			}
			// 23503 (foreign_key_violation) -> Create wallet for unexists user
			if pgxError.Code == "23503" {
				return nil, util.ErrAccUser
			}
			// 23505 (unique_violation) ->  Duplicate wallet with the same currency
			if pgxError.Code == "23505" {
				return nil, pgxError
			}
		}
		return nil, err
	}

	return &res, nil
}

func (r *DB) GetWallet(ctx context.Context, id int64) (*pkg.Wallet, error) {
	row := r.db.QueryRow(ctx, "SELECT * FROM wallets WHERE id=$1 AND deleted_at IS NULL;", &id)

	var wallet pkg.Wallet
	err := row.Scan(&wallet.ID, &wallet.AccountID, &wallet.WalletNumber, &wallet.Name, &wallet.Balance, &wallet.Currency, &wallet.CreatedAt, &wallet.DeletedAt)
	if err == pgx.ErrNoRows {
		return nil, pgx.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *DB) GetWalletByNumber(ctx context.Context, number int64) (*pkg.Wallet, error) {
	row := r.db.QueryRow(ctx, "SELECT * FROM wallets WHERE wallet_number=$1 AND deleted_at IS NULL;", &number)

	var wallet pkg.Wallet
	err := row.Scan(&wallet.ID, &wallet.AccountID, &wallet.WalletNumber, &wallet.Name, &wallet.Balance, &wallet.Currency, &wallet.CreatedAt, &wallet.DeletedAt)
	if err == pgx.ErrNoRows {
		return nil, util.ErrNotExist
	}
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

type ListWalletParams struct {
	AccountID int64
	Limit     int
	Offset    int
}

func (r *DB) ListWallet(ctx context.Context, arg ListWalletParams) ([]pkg.Wallet, error) {
	//log.Printf("account Id: %d - limit: %d - offset: %d\n", arg.AccountID, arg.Limit, arg.Offset)
	query := `SELECT * FROM wallets WHERE account_id=$1 AND deleted_at IS NULL AND name!='Primary Wallet' ORDER BY id LIMIT $2 OFFSET $3;`
	res, err := r.db.Query(ctx, query, arg.AccountID, arg.Limit, arg.Offset)

	if err != nil {
		log.Println("err 1")
		return nil, err
	}

	var list []pkg.Wallet
	for res.Next() {
		var temp pkg.Wallet
		//log.Println("pass 1")
		if err := res.Scan(&temp.ID, &temp.AccountID, &temp.WalletNumber, &temp.Name, &temp.Balance, &temp.Currency, &temp.CreatedAt, &temp.DeletedAt); err != nil {
			//log.Println("err 2")
			if err == pgx.ErrNoRows {
				return nil, util.ErrNotExist
			}
			return nil, err
		}
		list = append(list, temp)
		//log.Println("pass 2")
	}

	res.Close()

	return list, nil
}

type UpdateWalletParams struct {
	ID      int64
	Balance int64
}

func (r *DB) UpdateWallet(ctx context.Context, arg UpdateWalletParams) error {
	res, err := r.db.Exec(ctx, "UPDATE wallets SET balance=$1 WHERE id=$2 AND deleted_at IS NULL", arg.Balance, arg.ID)
	if err != nil {
		log.Println("Exec error")
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		log.Println("Update Failed")
		return util.ErrUpdateFailed
	}

	return nil
}

type UpdateWalletInformationParams struct {
	WalletNumber int64
	Name         string
	Currency     string
}

func (r *DB) UpdateWalletInformation(ctx context.Context, arg UpdateWalletInformationParams) error {
	res, err := r.db.Exec(ctx, "UPDATE wallets SET name=$1, currency=$2 WHERE wallet_number=$3 AND deleted_at IS NULL", arg.Name, arg.Currency, arg.WalletNumber)
	if err != nil {
		log.Println("Exec error")
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		log.Println("Update Failed")
		return util.ErrUpdateFailed
	}

	return nil
}

/*func (r *DB) DeleteWallet(ctx context.Context, id int64) error {
	res, err := r.db.Exec(ctx, "DELETE FROM wallets WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		log.Println("(DeleteWallet) Exec Error =", err)
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		log.Println("(DeleteWallet) Delete wallet failed")
		return err
	}

	return nil
}*/

func (r *DB) DeleteWallet(ctx context.Context, id int64) error {
	query := `UPDATE wallets SET deleted_at=now() WHERE id=$1`
	res, err := r.db.Exec(ctx, query, id)

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		log.Println("(DeleteWallet) Delete wallet failed")
		return err
	}

	return nil
}

func (r *DB) GetWalletForUpdate(ctx context.Context, id int64) (*pkg.Wallet, error) {
	query := `SELECT * FROM wallets WHERE id=$1 AND deleted_at IS NULL FOR NO KEY UPDATE`
	row := r.db.QueryRow(ctx, query, id)

	var wallet pkg.Wallet
	var err error
	if err = row.Scan(&wallet.ID, &wallet.AccountID, &wallet.WalletNumber, &wallet.Name, &wallet.Balance, &wallet.Currency, &wallet.CreatedAt, &wallet.DeletedAt); err != nil {
		return nil, err
	}

	return &wallet, err
}

type AddWalletBalanceParams struct {
	WalletNumber int64
	Amount       int64
}

func (r *DB) AddWalletBalance(ctx context.Context, arg AddWalletBalanceParams) (*pkg.Wallet, error) {
	var res pkg.Wallet
	query := `UPDATE wallets SET balance=balance+$1 WHERE wallet_number=$2
	RETURNING id, account_id, wallet_number, name, balance, currency, created_at`
	err := r.db.QueryRow(ctx, query, arg.Amount, arg.WalletNumber).Scan(&res.ID, &res.AccountID, &res.WalletNumber, &res.Name, &res.Balance, &res.Currency, &res.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &res, err
}
