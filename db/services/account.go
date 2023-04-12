package services

import (
	"context"
	"log"

	"simple-bank-system/db/pkg"
)

func (r *DB) CreateAccount(ctx context.Context, account pkg.Account) (*pkg.Account, error) {
	var res pkg.Account

	query := `INSERT INTO accounts(owner, balance, currency
		) VALUES(
			$1, $2, $3
		) RETURNING id, owner, balance, currency, created_at;`
	err := r.db.QueryRow(ctx, query, account.Owner, account.Balance, account.Currency).Scan(&res.ID, &res.Owner, &res.Balance, &res.Currency, &res.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *DB) GetAccount(ctx context.Context, id int64) (*pkg.Account, error) {
	row := r.db.QueryRow(ctx, "SELECT * FROM accounts WHERE id=$1;", &id)

	var account pkg.Account
	if err := row.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt); err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *DB) ListAccount(ctx context.Context, limit int, offset int) ([]pkg.Account, error) {
	res, err := r.db.Query(ctx, "SELECT * FROM accounts ORDER BY id LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return nil, err
	}

	var list []pkg.Account
	for res.Next() {
		var temp pkg.Account
		if err := res.Scan(&temp.ID, &temp.Owner, &temp.Balance, &temp.Currency, &temp.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, temp)
	}

	return list, nil
}

type UpdateAccountParams struct {
	ID      int64
	Balance int64
}

func (r *DB) UpdateAccount(ctx context.Context, arg UpdateAccountParams) error {
	res, err := r.db.Exec(ctx, "UPDATE accounts SET balance=$1 WHERE id=$2", arg.Balance, arg.ID)
	if err != nil {
		log.Println("Exec error")
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		log.Println("Update Failed")
		return nil
	}

	return nil
}

func (r *DB) DeleteAccount(ctx context.Context, id int64) error {
	res, err := r.db.Exec(ctx, "DELETE FROM accounts WHERE id=$1", id)
	if err != nil {
		log.Println("Exec Error = ")
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		log.Println("Delete account failed")
		return err
	}

	return nil
}

func (r *DB) GetAccountForUpdate(ctx context.Context, id int64) (*pkg.Account, error) {
	query := `SELECT * FROM accounts WHERE id=$1 FOR NO KEY UPDATE`
	row := r.db.QueryRow(ctx, query, id)

	var account pkg.Account
	var err error
	if err = row.Scan(&account.ID, &account.Owner, &account.Balance, &account.Currency, &account.CreatedAt); err != nil {
		return nil, err
	}
	return &account, err
}

type AddAccountBalanceParams struct {
	ID     int64
	Amount int64
}

func (r *DB) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (*pkg.Account, error) {
	var res pkg.Account
	query := `UPDATE accounts SET balance=balance+$1 WHERE id=$2
	RETURNING id, owner, balance, currency, created_at`
	err := r.db.QueryRow(ctx, query, arg.Amount, arg.ID).Scan(&res.ID, &res.Owner, &res.Balance, &res.Currency, &res.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &res, err
}
