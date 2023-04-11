package services

import (
	"context"
	"fmt"
	"simple-bank-system/db/pkg"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

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
		// if err != nil {
		// 	return err
		// }

		//fmt.Println(txName, "update account 1")
		err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		// check updated account
		result.FromAccount, err = q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		//To Account
		//fmt.Println(txName, "update account 2")
		err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		return nil
	})

	return &result, err
}
