package services

import (
	"context"

	"simple-bank-system/db/pkg"
)

type CreateTransferParam struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

func (c *DB) CreateTransfer(ctx context.Context, arg CreateTransferParam) (*pkg.Transfers, error) {
	query := `INSERT INTO transfers(from_account_id, to_account_id, amount
	) VALUES (
		$1, $2, $3
	) RETURNING id, from_account_id, to_account_id, amount, created_at;`

	var res pkg.Transfers
	err := c.db.QueryRow(ctx, query, arg.FromAccountID, arg.ToAccountID, arg.Amount).Scan(&res.ID, &res.FromAccountID, &res.ToAccountID, &res.Amount, &res.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *DB) GetTransfer(ctx context.Context, id int64) (*pkg.Transfers, error) {
	query := `SELECT * FROM transfers WHERE id=$1;`
	var res pkg.Transfers
	err := c.db.QueryRow(ctx, query, id).Scan(&res.ID, &res.FromAccountID, &res.ToAccountID, &res.Amount, &res.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

type ListTransfersParam struct {
	limit  int
	offset int
}

func (c *DB) ListTransfers(ctx context.Context, arg ListTransfersParam) ([]pkg.Transfers, error) {
	query := `SELECT * FROM transfers ORDER BY id LIMIT $1 OFFSET $2;`
	res, err := c.db.Query(ctx, query, arg.limit, arg.offset)
	if err != nil {
		return nil, err
	}

	var list []pkg.Transfers
	for res.Next() {
		var temp pkg.Transfers
		if err = res.Scan(&temp.ID, &temp.FromAccountID, &temp.ToAccountID, &temp.Amount, &temp.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, temp)
	}

	return list, nil
}
