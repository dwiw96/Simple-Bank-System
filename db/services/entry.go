package services

import (
	"context"

	"simple-bank-system/db/pkg"
)

type CreateEntryParam struct {
	accountID int64
	amount    int64
}

func (c *DB) CreateEntry(ctx context.Context, arg CreateEntryParam) (*pkg.Entry, error) {
	query := `INSERT INTO entries (account_id, amount
	) VALUES (
		$1, $2
	) RETURNING id, account_id, amount, created_at;`

	var res pkg.Entry
	err := c.db.QueryRow(ctx, query, arg.accountID, arg.amount).Scan(&res.ID, &res.AccountID, &res.Amount, &res.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *DB) GetEntry(ctx context.Context, id int64) (*pkg.Entry, error) {
	query := `SELECT * FROM entries WHERE id=$1;`
	var res pkg.Entry
	err := c.db.QueryRow(ctx, query, id).Scan(&res.ID, &res.AccountID, &res.Amount, &res.CreatedAt)
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
	query := `SELECT * FROM entries ORDER BY id LIMIT $1 OFFSET $2;`
	row, err := c.db.Query(ctx, query, arg.limit, arg.offset)
	if err != nil {
		return nil, err
	}

	var res []pkg.Entry
	for row.Next() {
		var temp pkg.Entry
		if err = row.Scan(&temp.ID, &temp.AccountID, &temp.Amount, &temp.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, temp)
	}

	return res, nil
}
