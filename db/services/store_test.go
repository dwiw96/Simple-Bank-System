package services

import (
	"context"
	"fmt"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(dbpool)

	// Create 2 accounts to perform transfer between those two accounts
	account1, _ := createRandomAccount(t)
	account2, _ := createRandomAccount(t)
	fmt.Printf("\n>> before: %d  -  %d\n", account1.Balance, account2.Balance)

	var errs = make(chan error)
	var results = make(chan *TransferTXResult)

	amount := int64(10)
	n := 5

	// do transfer n time using go concurency
	for i := 0; i < n; i++ {
		//txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			//ctx := context.WithValue(context.Background(), txKey, txName)
			arg := TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			}
			result, err := store.TransferTx(context.Background(), arg)

			errs <- err
			results <- result
		}()
	}

	// check result
	k_Check := make(map[int64]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		if err != nil {
			t.Errorf("TransferTx() err: %s", err)
		}

		result := <-results
		if result == nil {
			t.Errorf("result is empty")
		}

		if result.Transfer == nil {
			t.Errorf("transfer is empty")
		}
		if result.Transfer.ID == 0 {
			t.Errorf("result.Transfer.ID is empty")
		}
		if result.Transfer.FromAccountID != account1.ID {
			t.Errorf("From Account ID in DB %d != %d input", result.Transfer.FromAccountID, account1.ID)
		}
		if result.Transfer.ToAccountID != account2.ID {
			t.Errorf("To Account ID in DB %d != %d input", result.Transfer.ToAccountID, account2.ID)
		}
		if result.Transfer.Amount != amount {
			t.Errorf("Amount in DB %d != %d input", result.Transfer.Amount, 10)
		}
		if result.Transfer.CreatedAt.IsZero() {
			t.Errorf("Created date is zero")
		}

		_, err = store.GetTransfer(ctx, result.Transfer.ID)
		if err != nil {
			t.Errorf("can't get transfer data from DB, ID: %d", result.Transfer.ID)
		}

		if result.FromEntry.ID == 0 {
			t.Errorf("result From Entry ID is empty")
		}
		if result.FromEntry.AccountID != account1.ID {
			t.Errorf("From Entry ID	%d != %d input", result.FromEntry.ID, account1.ID)
		}
		if result.FromEntry.Amount != -amount {
			t.Errorf("Amount from entry in DB %d != -10 input", result.FromEntry.Amount)
		}
		if result.FromEntry.CreatedAt.IsZero() {
			t.Errorf("from entry date is empty")
		}

		_, err = store.GetEntry(ctx, result.FromEntry.ID)
		if err != nil {
			t.Errorf("can't get FromEntry data from DB, FromEntryID: %d", result.FromEntry.ID)
		}

		if result.ToEntry.ID == 0 {
			t.Errorf("result to Entry ID is empty")
		}
		if result.ToEntry.AccountID != account2.ID {
			t.Errorf("to Entry ID	%d != %d input", result.FromEntry.ID, account2.ID)
		}
		if result.ToEntry.Amount != amount {
			t.Errorf("Amount to entry in DB %d != 10 input", result.ToEntry.Amount)
		}
		if result.ToEntry.CreatedAt.IsZero() {
			t.Errorf("to entry date is empty")
		}

		_, err = store.GetEntry(ctx, result.ToEntry.ID)
		if err != nil {
			t.Errorf("can't get ToEntry data from DB, ToEntryID: %d", result.ToEntry.ID)
		}

		// Account
		if result.FromAccount == nil {
			t.Errorf("FromAccount is empty")
		}
		if result.FromAccount.ID != account1.ID {
			t.Errorf("FromAccount %d != %d input", result.FromAccount.ID, account1.ID)
		}

		if result.ToAccount == nil {
			t.Errorf("ToAccount is empty")
		}
		if result.ToAccount.ID != account2.ID {
			t.Errorf("ToAccount ID %d != %d input", result.ToAccount.ID, account2.ID)
		}

		//Check balance difference before and after transfer the money
		fmt.Printf("\n>> tx: %d  -  %d\n", result.FromAccount.Balance, result.ToAccount.Balance)
		accBal1 := account1.Balance - result.FromAccount.Balance
		accBal2 := result.ToAccount.Balance - account2.Balance

		if accBal1 != accBal2 {
			t.Errorf("Balance different account1 %d != %d account2", accBal1, accBal2)
		}
		if accBal1 < 0 && accBal2 < 0 {
			t.Errorf("Balance different account1 %d & account2 %d <= 0", accBal1, accBal2)
		}
		if accBal1%amount != 0 && accBal2%amount != 0 {
			t.Errorf("Modulus of accBal1 %d & accBal2 %d != 0", accBal1, accBal2)
		}

		// var k must be unique, 1th transfer = 1, 2th transfer = 2, ...
		k := (accBal1 / amount)
		if k < 1 && k > int64(n) {
			t.Errorf("1 > k(%d) > n", k)
		}
		if _, ok := k_Check[k]; ok {
			t.Errorf("k isn't unique: %d", k)
		} else {
			k_Check[k] = true
		}
	}

	//Check the final balance of 2 accounts after do transfer
	updateAcc1, err := store.GetAccount(ctx, account1.ID)
	if err != nil {
		t.Errorf("GetAccount1 err: %v", err)
	}
	updateAcc2, err := store.GetAccount(ctx, account2.ID)
	if err != nil {
		t.Errorf("GetAccoun2 err: %v", err)
	}

	fmt.Printf("\n>> after: %d  -  %d\n", updateAcc1.Balance, updateAcc2.Balance)
	acc1Bal := account1.Balance - (int64(n) * amount)
	if updateAcc1.Balance != acc1Bal {
		t.Errorf("Account1 final balance %d != %d", &updateAcc1.Balance, acc1Bal)
	}
	acc2Bal := account2.Balance + (int64(n) * amount)
	if updateAcc2.Balance != acc2Bal {
		t.Errorf("Account2 final balance %d != %d", &updateAcc2.Balance, acc2Bal)
	}
}
