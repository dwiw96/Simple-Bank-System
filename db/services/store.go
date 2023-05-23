/*
 * Implement DB transaction 
 * This code will provide all functions to run database queries individually, as well as
 * their combinations functions within a transaction.
 * DB struct doesn't support transaction because each query only can do 1 operations on 1
 * spesific table, Store struct in this file solved that problem.
*/

package services

import (
	"context"
	"fmt"
	"simple-bank-system/db/pkg"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Store used to embedded DB struct to extend DB functionality
type Store struct {
	*DB
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
		DB: NewDB(db),
	}
}

// function to the Store to execute a generic database transaction.
//func (store *Store) execTx(ctx context.Context, fn func(*DB) error) error {
//}

func (store *Store) execTx(ctx context.Context, fn func(*DB) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := NewDB(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx rr = %s, and rb err = %s", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}
type TransferTXResult struct {
	Transfer    *pkg.Transfers
	FromAccount *pkg.Account
	ToAccount   *pkg.Account
	FromEntry   *pkg.Entry
	ToEntry     *pkg.Entry
}

//var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (*TransferTXResult, error) {
	var result TransferTXResult

	err := store.execTx(ctx, func(q *DB) error {
		var err error

		//txName := ctx.Value(txKey)

		//fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParam{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		//fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParam{
			accountID: arg.FromAccountID,
			amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		//fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParam{
			accountID: arg.ToAccountID,
			amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// From Account
		// 1. GetAccount to get the balance, 2. Update the balance
		// 3. GetAccount again to get updated balance
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)

		if arg.FromAccountID < arg.ToAccountID {
			//fmt.Println(txName, "update account 1")
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			//fmt.Println(txName, "update account 2")
			result.ToAccount, result.FromAccount, err = AddMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return &result, err
}

func AddMoney(ctx context.Context, q *DB, accountID1, amount1, accountID2, amount2 int64) (account1 *pkg.Account, account2 *pkg.Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return nil, nil, err
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return nil, nil, err
	}
	return
}
