package services

import (
	"context"

	"simple-bank-system/db/pkg"

	"github.com/jackc/pgx/v4"
)

type CreateTransferParam struct {
	AccountID        int64
	WalletID         int64
	FromWalletNumber int64
	ToWalletNumber   int64
	Amount           int64
}

func (c *DB) CreateTransfer(ctx context.Context, arg CreateTransferParam) (*pkg.Transfers, error) {
	query := `INSERT INTO transfers(account_id, wallet_id, from_wallet_number, to_wallet_number, amount
	) VALUES (
		$1, $2, $3, $4, $5
	) RETURNING id, account_id, wallet_id, from_wallet_number, to_wallet_number, amount, created_at;`

	var res pkg.Transfers
	err := c.db.QueryRow(ctx, query, arg.AccountID, arg.WalletID, arg.FromWalletNumber, arg.ToWalletNumber, arg.Amount).Scan(&res.ID, &res.AccountID, &res.WalletID, &res.FromWalletNumber, &res.ToWalletNumber, &res.Amount, &res.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *DB) GetTransfer(ctx context.Context, id int64, option string) (*pkg.Transfers, error) {
	queryID := `SELECT * FROM transfers WHERE id=$1 AND deleted_at IS NULL;`
	queryAcc := `SELECT * FROM transfers WHERE account_id=$1 AND deleted_at IS NULL;`
	queryWal := `SELECT * FROM transfers WHERE from_wallet_number=$1 AND deleted_at IS NULL;`
	queryLast := `SELECT * FROM transfers WHERE from_wallet_number=$1 AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1;`

	queries := []struct {
		option string
		query  string
	}{
		{
			option: "ID",
			query:  queryID,
		}, {
			option: "AccountID",
			query:  queryAcc,
		}, {
			option: "WalletNumber",
			query:  queryWal,
		}, {
			option: "Last",
			query:  queryLast,
		},
	}

	var res pkg.Transfers
	var err error

	for _, dbQuery := range queries {
		if dbQuery.option == option {
			err = c.db.QueryRow(ctx, dbQuery.query, id).Scan(&res.ID, &res.AccountID, &res.WalletID, &res.FromWalletNumber, &res.ToWalletNumber, &res.Amount, &res.CreatedAt, &res.DeletedAt)
			break
		}
	}
	/*if option == 1 {
		err = c.db.QueryRow(ctx, queryID, id).Scan(&res.ID, &res.AccountID, &res.WalletID, &res.FromWalletNumber, &res.ToWalletNumber, &res.Amount, &res.CreatedAt, &res.DeletedAt)
	} else if option == 2 {
		err = c.db.QueryRow(ctx, queryAcc, id).Scan(&res.ID, &res.AccountID, &res.WalletID, &res.FromWalletNumber, &res.ToWalletNumber, &res.Amount, &res.CreatedAt, &res.DeletedAt)
	} else if option == 3 {
		err = c.db.QueryRow(ctx, queryWal, id).Scan(&res.ID, &res.AccountID, &res.WalletID, &res.FromWalletNumber, &res.ToWalletNumber, &res.Amount, &res.CreatedAt, &res.DeletedAt)
	}*/

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
	query := `SELECT * FROM transfers WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2;`
	res, err := c.db.Query(ctx, query, arg.limit, arg.offset)
	if err != nil {
		return nil, err
	}

	var list []pkg.Transfers
	for res.Next() {
		var temp pkg.Transfers
		if err = res.Scan(&temp.ID, &temp.AccountID, &temp.WalletID, &temp.FromWalletNumber, &temp.ToWalletNumber, &temp.Amount, &temp.CreatedAt, &temp.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, temp)
	}

	return list, nil
}

type ListTransfersByDateParam struct {
	Start string
	End   string
}

func (c *DB) ListTransfersByDate(ctx context.Context, arg ListTransfersByDateParam, order bool) ([]pkg.Transfers, error) {
	queryDesc := `SELECT * FROM transfers
	WHERE created_at BETWEEN $1 AND $2 AND deleted_at IS NULL
	ORDER BY created_at DESC; -- Order by the most recent timestamp;
	`
	queryAsc := `SELECT * FROM transfers
	WHERE created_at BETWEEN $1 AND $2 AND deleted_at IS NULL
	ORDER BY created_at ASC; -- Order by the most older timestamp;
	`

	var res pgx.Rows
	var err error

	if order == true {
		res, err = c.db.Query(ctx, queryDesc, arg.Start, arg.End)
	} else {
		res, err = c.db.Query(ctx, queryAsc, arg.Start, arg.End)
	}

	if err != nil {
		return nil, err
	}

	var list []pkg.Transfers
	for res.Next() {
		var temp pkg.Transfers
		if err = res.Scan(&temp.ID, &temp.AccountID, &temp.WalletID, &temp.FromWalletNumber, &temp.ToWalletNumber, &temp.Amount, &temp.CreatedAt, &temp.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, temp)
	}

	return list, nil
}
