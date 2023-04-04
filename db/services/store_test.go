package services

import (
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(dbpool)

	account1, _ := createRandomAccount(t)
	account2, _ := createRandomAccount(t)

	// resAcc1, err := store.GetAccount(ctx, account1.ID)
	// if err != nil {
	// 	t.Fatalf("can't get account1 from DB")
	// }
	// resAcc2, err := store.GetAccount(ctx, account2.ID)
	// if err != nil {
	// 	t.Fatalf("can't get account1 from DB")
	// }

	// fmt.Printf("\naccount1 (%d) balance = %d", resAcc1.ID, resAcc1.Balance)
	// fmt.Printf("\naccount2 (%d) balance = %d\n", resAcc2.ID, resAcc2.Balance)

	var errs = make(chan error)
	var results = make(chan *TransferTXResult)

	amount := int64(10)
	n := 5
	for i := 0; i < n; i++ {
		go func() {
			arg := TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			}
			result, err := store.TransferTx(ctx, arg)

			errs <- err
			results <- result
		}()
	}

	// check result
	for i := 0; i < n; i++ {
		err := <-errs
		if err != nil {
			t.Fatalf("TransferTx() err: %s", err)
		}

		result := <-results
		if result == nil {
			t.Fatalf("result is empty")
		}

		if result.Transfer == nil {
			t.Fatalf("transfer is empty")
		} else if result.Transfer.ID == 0 {
			t.Fatalf("result.Transfer.ID is empty")
		} else if result.Transfer.FromAccountID != account1.ID {
			t.Fatalf("From Account ID in DB %d != %d input", result.Transfer.FromAccountID, account1.ID)
		} else if result.Transfer.ToAccountID != account2.ID {
			t.Fatalf("To Account ID in DB %d != %d input", result.Transfer.ToAccountID, account2.ID)
		} else if result.Transfer.Amount != amount {
			t.Fatalf("Amount in DB %d != %d input", result.Transfer.Amount, 10)
		} else if result.Transfer.CreatedAt.IsZero() {
			t.Fatalf("Created date is zero")
		}

		_, err = store.GetTransfer(ctx, result.Transfer.ID)
		if err != nil {
			t.Fatalf("can't get transfer data from DB")
		}

		if result.FromEntry.ID == 0 {
			t.Fatalf("result From Entry ID is empty")
		} else if result.FromEntry.AccountID != account1.ID {
			t.Fatalf("From Entry ID	%d != %d input", result.FromEntry.ID, account1.ID)
		} else if result.FromEntry.Amount != -amount {
			t.Fatalf("Amount from entry in DB %d != -10 input", result.FromEntry.Amount)
		} else if result.FromEntry.CreatedAt.IsZero() {
			t.Fatalf("from entry date is empty")
		}

		_, err = store.GetEntry(ctx, result.FromEntry.AccountID)
		if err != nil {
			t.Fatalf("can't get From Entry data from DB")
		}

		if result.ToEntry.ID == 0 {
			t.Fatalf("result to Entry ID is empty")
		} else if result.ToEntry.AccountID != account2.ID {
			t.Fatalf("to Entry ID	%d != %d input", result.FromEntry.ID, account2.ID)
		} else if result.ToEntry.Amount != amount {
			t.Fatalf("Amount to entry in DB %d != 10 input", result.ToEntry.Amount)
		} else if result.ToEntry.CreatedAt.IsZero() {
			t.Fatalf("to entry date is empty")
		}

		_, err = store.GetEntry(ctx, result.ToEntry.AccountID)
		if err != nil {
			t.Fatalf("can't get To Entry data from DB")
		}
	}
}
