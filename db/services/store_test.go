package services

import (
	"context"
	//"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransferTxSuccess(t *testing.T) {
	store := NewStore(dbpool)

	// Create accounts and wallet to perform transfer between those two accounts
	account1 := createRandomAccount(t)
	wallet1, _ := createRandomWallet(t, account1)
	account2 := createRandomAccount(t)
	wallet2, _ := createRandomWallet(t, account2)

	//fmt.Printf("\n>> before: %d  -  %d\n", account1.Balance, account2.Balance)

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
				AccountID:        account1.ID,
				WalletID:         wallet1.ID,
				FromWalletNumber: wallet1.WalletNumber,
				ToWalletNumber:   wallet2.WalletNumber,
				Amount:           amount,
			}
			result, err := store.TransferTx(context.Background(), arg)

			errs <- err
			results <- result
		}()
	}

	// check result
	k_Check := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err, "TransferTx() Fail")
		result := <-results
		require.NotNil(t, result, "result is empty")

		// Check for transfer record
		assert.NotNil(t, result.Transfer)
		assert.NotZero(t, result.Transfer, "result.Transfer.ID is empty")
		assert.NotZero(t, result.Transfer.ID)
		assert.Equal(t, account1.ID, result.Transfer.AccountID)
		assert.Equal(t, wallet1.ID, result.Transfer.WalletID)
		assert.Equal(t, wallet1.WalletNumber, result.Transfer.FromWalletNumber)
		assert.Equal(t, wallet2.WalletNumber, result.Transfer.ToWalletNumber)
		assert.Equal(t, amount, result.Transfer.Amount)
		assert.NotZero(t, result.Transfer.CreatedAt)

		_, err = store.GetTransfer(ctx, result.Transfer.ID, "ID")
		require.NoError(t, err)

		// Check for `From Entry` Record
		assert.NotZero(t, result.FromEntry.ID)
		assert.Equal(t, account1.ID, result.FromEntry.AccountID)
		assert.Equal(t, wallet1.ID, result.FromEntry.WalletID)
		assert.Equal(t, wallet1.WalletNumber, result.FromEntry.WalletNumber)
		assert.Equal(t, -amount, result.FromEntry.Amount)
		assert.NotZero(t, result.FromEntry.CreatedAt)

		_, err = store.GetEntry(ctx, result.FromEntry.ID, "ID")
		require.NoError(t, err)

		// Check for `To Entry` Record
		assert.NotZero(t, result.ToEntry.ID)
		assert.Equal(t, account2.ID, result.ToEntry.AccountID)
		assert.Equal(t, wallet2.ID, result.ToEntry.WalletID)
		assert.Equal(t, wallet2.WalletNumber, result.ToEntry.WalletNumber)
		assert.Equal(t, amount, result.ToEntry.Amount)
		assert.NotZero(t, result.ToEntry.CreatedAt)

		_, err = store.GetEntry(ctx, result.ToEntry.ID, "ID")
		require.NoError(t, err, "can't get ToEntry data from DB, ToEntryID: %d", result.ToEntry.ID)

		// Check for `Wallet` Record
		require.NotEmpty(t, result.FromWallet)
		assert.Equal(t, wallet1.ID, result.FromWallet.ID)
		assert.Equal(t, wallet1.WalletNumber, result.FromWallet.WalletNumber)
		require.NotEmpty(t, result.ToWallet)
		assert.Equal(t, wallet2.ID, result.ToWallet.ID)
		assert.Equal(t, wallet2.WalletNumber, result.ToWallet.WalletNumber)

		//Check balance difference before and after transfer the money
		//fmt.Printf("\n>> tx: %d  -  %d\n", result.FromWallet.Balance, result.ToWallet.Balance)
		accBal1 := wallet1.Balance - result.FromWallet.Balance
		accBal2 := result.ToWallet.Balance - wallet2.Balance

		assert.Equal(t, accBal1, accBal2)
		assert.True(t, accBal1 > 0)
		assert.True(t, accBal2 > 0)
		assert.True(t, accBal1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount
		assert.True(t, accBal2%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		// var k must be unique, 1th transfer = 1, 2th transfer = 2, ...
		k := int((accBal1 / amount))
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, k_Check, k) // _, ok := k_Check[k]; ok {t.Errorf("k isn't unique: %d", k)
		k_Check[k] = true
	}

	//Check the final balance of 2 accounts after do transfer
	updateWallet1, err := store.GetWalletByNumber(ctx, wallet1.WalletNumber)
	require.NoError(t, err)
	updateWallet2, err := store.GetWalletByNumber(ctx, wallet2.WalletNumber)
	require.NoError(t, err)

	//fmt.Printf("\n>> after: %d  -  %d\n", updateWallet1.Balance, updateWallet2.Balance)
	acc1Bal := wallet1.Balance - (int64(n) * amount)
	assert.Equal(t, acc1Bal, updateWallet1.Balance, "Account1 final balance %d != %d", &updateWallet1.Balance, acc1Bal)

	acc2Bal := wallet2.Balance + (int64(n) * amount)
	assert.Equal(t, acc2Bal, updateWallet2.Balance, "Account2 final balance %d != %d", &updateWallet2.Balance, acc2Bal)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(dbpool)

	// Create 2 accounts to perform transfer between those two accounts
	account1 := createRandomAccount(t)
	wallet1, _ := createRandomWallet(t, account1)
	account2 := createRandomAccount(t)
	wallet2, _ := createRandomWallet(t, account2)
	//fmt.Printf("\n>> before: %d  -  %d\n", wallet1.Balance, wallet2.Balance)

	var errs = make(chan error)

	amount := int64(10)
	n := 10 // n time of concurency

	// do transfer n time using go concurency
	for i := 0; i < n; i++ {
		//txName := fmt.Sprintf("tx %d", i+1)
		arg := TransferTxParams{
			AccountID:        account1.ID,
			WalletID:         wallet1.ID,
			FromWalletNumber: wallet1.WalletNumber,
			ToWalletNumber:   wallet2.WalletNumber,
			Amount:           amount,
		}

		if i%2 == 1 {
			arg = TransferTxParams{
				AccountID:        account2.ID,
				WalletID:         wallet2.ID,
				FromWalletNumber: wallet2.WalletNumber,
				ToWalletNumber:   wallet1.WalletNumber,
				Amount:           amount,
			}
		}
		go func() {
			//ctx := context.WithValue(context.Background(), txKey, txName)
			_, err := store.TransferTx(context.Background(), arg)

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	//Check the final balance of 2 accounts after do transfer
	updateWallet1, err := store.GetWalletByNumber(ctx, wallet1.WalletNumber)
	require.NoError(t, err)
	updateWallet2, err := store.GetWalletByNumber(ctx, wallet2.WalletNumber)
	require.NoError(t, err)
	//fmt.Printf("\n>> after: %d  -  %d\n", updateWallet1.Balance, updateWallet2.Balance)

	assert.Equal(t, wallet1.Balance, updateWallet1.Balance)
	assert.Equal(t, wallet2.Balance, updateWallet2.Balance)
}
