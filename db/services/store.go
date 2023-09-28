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
	"log"
	"simple-bank-system/db/pkg"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// "Store" used to embedded DB struct to extend DB functionality
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

// Add execTc() function to the "Store" to execute a generic database transaction.
/*
 * execTx() will start new db transaction, then create new "DB" object.
 * call callback function with the created "DB", and finally commit or rollback
 * the transaction based on error returned by the function.
 */

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
	AccountID        int64
	WalletID         int64
	FromWalletNumber int64
	ToWalletNumber   int64
	Amount           int64
}
type TransferTXResult struct {
	Transfer   *pkg.Transfers
	FromWallet *pkg.Wallet
	ToWallet   *pkg.Wallet
	FromEntry  *pkg.Entry
	ToEntry    *pkg.Entry
}

//var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (*TransferTXResult, error) {
	var result TransferTXResult

	err := store.execTx(ctx, func(q *DB) error {
		var err error

		//txName := ctx.Value(txKey)

		//fmt.Println(txName, "create transfer")
		//fmt.Println("(input) Transfer Tx Params:", arg)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParam{
			AccountID:        arg.AccountID,
			WalletID:         arg.WalletID,
			FromWalletNumber: arg.FromWalletNumber,
			ToWalletNumber:   arg.ToWalletNumber,
			Amount:           arg.Amount,
		})
		if err != nil {
			log.Println("--(err) 1")
			return err
		}

		//fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParam{
			accountID:    arg.AccountID,
			walletID:     arg.WalletID,
			walletNumber: arg.FromWalletNumber,
			amount:       -arg.Amount,
		})
		if err != nil {
			log.Println("--(err) 2")
			return err
		}

		toWallet, err := q.GetWalletByNumber(ctx, arg.ToWalletNumber)
		if err != nil {
			log.Println("--(err) 3")
			return err
		}
		//fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParam{
			accountID:    toWallet.AccountID,
			walletID:     toWallet.ID,
			walletNumber: arg.ToWalletNumber,
			amount:       arg.Amount,
		})
		if err != nil {
			log.Println("--(err) 4")
			return err
		}

		// From Account
		// 1. GetAccount to get the balance, 2. Update the balance
		// 3. GetAccount again to get updated balance
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromWalletID)

		/*
		 * To avoid deadlock, I make the wallet with the smaller wallet_number to update first
		 */

		if arg.FromWalletNumber < arg.ToWalletNumber {
			//fmt.Println(txName, "update account 1")
			result.FromWallet, result.ToWallet, err = AddMoney(ctx, q, arg.FromWalletNumber, -arg.Amount, arg.ToWalletNumber, arg.Amount)
			if err != nil {
				log.Println("--(err) 5:", err)
				return err
			}
		} else {
			//fmt.Println(txName, "update account 2")
			result.ToWallet, result.FromWallet, err = AddMoney(ctx, q, arg.ToWalletNumber, arg.Amount, arg.FromWalletNumber, -arg.Amount)
			if err != nil {
				log.Println("--(err) 6:", err)
				return err
			}
		}

		return nil
	})

	return &result, err
}

func AddMoney(ctx context.Context, q *DB, walletNumb1, amount1, walletNumb2, amount2 int64) (wallet1 *pkg.Wallet, wallet2 *pkg.Wallet, err error) {
	wallet1, err = q.AddWalletBalance(ctx, AddWalletBalanceParams{
		WalletNumber: walletNumb1,
		Amount:       amount1,
	})
	if err != nil {
		return nil, nil, err
	}

	wallet2, err = q.AddWalletBalance(ctx, AddWalletBalanceParams{
		WalletNumber: walletNumb2,
		Amount:       amount2,
	})
	if err != nil {
		return nil, nil, err
	}
	return
}
