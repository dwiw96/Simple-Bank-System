package services

import (
	"context"

	//"time"

	"simple-bank-system/db/pkg"

	"github.com/jackc/pgx/v4"
)

type CreateEntryParam struct {
	accountID    int64
	walletID     int64
	walletNumber int64
	amount       int64
}

func (c *DB) CreateEntry(ctx context.Context, arg CreateEntryParam) (*pkg.Entry, error) {
	query := `INSERT INTO entries (account_id, wallet_id, wallet_number, amount
	) VALUES (
		$1, $2, $3, $4
	) RETURNING id, account_id, wallet_id, wallet_number, amount, created_at;`

	var res pkg.Entry
	err := c.db.QueryRow(ctx, query, arg.accountID, arg.walletID, arg.walletNumber, arg.amount).Scan(&res.ID, &res.AccountID, &res.WalletID, &res.WalletNumber, &res.Amount, &res.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *DB) GetEntry(ctx context.Context, id int64, option string) (*pkg.Entry, error) {
	//query := `SELECT * FROM entries WHERE id=$1 AND deleted_at IS NULL;`
	queries := []struct {
		option string
		query  string
	}{
		{
			option: "ID",
			query:  `SELECT * FROM entries WHERE id=$1 AND deleted_at IS NULL;`,
		}, {
			option: "WalletNumber",
			query:  `SELECT * FROM entries WHERE wallet_number=$1 AND deleted_at IS NULL;`,
		}, {
			option: "Last",
			query:  `SELECT * FROM entries WHERE wallet_number=$1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1;`,
		},
	}
	var res pkg.Entry
	var err error
	for _, dbQuery := range queries {
		if dbQuery.option == option {
			err = c.db.QueryRow(ctx, dbQuery.query, id).Scan(&res.ID, &res.AccountID, &res.WalletID, &res.WalletNumber, &res.Amount, &res.CreatedAt, &res.DeletedAt)
		}
	}

	if err != nil {
		return nil, err
	}

	return &res, nil
}

type listEntryParam struct {
	limit  int
	offset int
}

func (c *DB) ListEntry(ctx context.Context, arg listEntryParam) ([]pkg.Entry, error) {
	query := `SELECT * FROM entries WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2;`
	row, err := c.db.Query(ctx, query, arg.limit, arg.offset)
	if err != nil {
		return nil, err
	}

	var res []pkg.Entry
	for row.Next() {
		var temp pkg.Entry
		if err = row.Scan(&temp.ID, &temp.AccountID, &temp.WalletID, &temp.WalletNumber, &temp.Amount, &temp.CreatedAt, &temp.DeletedAt); err != nil {
			return nil, err
		}
		res = append(res, temp)
	}

	return res, nil
}

func (c *DB) ListEntryByID(ctx context.Context, ID int64, walOk bool) ([]pkg.Entry, error) {
	queryAcc := `SELECT * FROM entries WHERE account_id=$1 AND deleted_at IS NULL ORDER BY created_at;`
	queryWal := `SELECT * FROM entries WHERE wallet_id=$1 AND deleted_at IS NULL ORDER BY created_at;`

	var row pgx.Rows
	var err error

	if walOk == true {
		row, err = c.db.Query(ctx, queryWal, ID)
	} else {
		row, err = c.db.Query(ctx, queryAcc, ID)
	}

	if err != nil {
		return nil, err
	}

	var res []pkg.Entry
	for row.Next() {
		var temp pkg.Entry
		if err = row.Scan(&temp.ID, &temp.AccountID, &temp.WalletID, &temp.WalletNumber, &temp.Amount, &temp.CreatedAt, &temp.DeletedAt); err != nil {
			return nil, err
		}
		res = append(res, temp)
	}

	return res, nil
}

func (c *DB) ListEntryByDate(ctx context.Context, startDate, endDate string, order bool) ([]pkg.Entry, error) {
	/*
	 * The BETWEEN clause will retrieve rows with timestamps within that range, and the ORDER BY clause with DESC
	 * will order them in descending order, so the most recent timestamps appear first in the result set.
	 */
	queryDesc := `SELECT * FROM entries
	WHERE created_at BETWEEN $1 AND $2  AND deleted_at IS NULL
	ORDER BY created_at DESC; -- Order by the most recent timestamp;
	`

	queryAsch := `SELECT * FROM entries
	WHERE created_at BETWEEN $1 AND $2  AND deleted_at IS NULL
	ORDER BY created_at ASC; -- Order by the most recent timestamp;
	`

	var row pgx.Rows
	var err error
	if order == true {
		row, err = c.db.Query(ctx, queryDesc, startDate, endDate)
		if err != nil {
			return nil, err
		}
	} else {
		row, err = c.db.Query(ctx, queryAsch, startDate, endDate)
		if err != nil {
			return nil, err
		}
	}

	var res []pkg.Entry
	for row.Next() {
		var temp pkg.Entry
		if err = row.Scan(&temp.ID, &temp.AccountID, &temp.WalletID, &temp.WalletNumber, &temp.Amount, &temp.CreatedAt, &temp.DeletedAt); err != nil {
			return nil, err
		}
		res = append(res, temp)
	}

	return res, nil
}
