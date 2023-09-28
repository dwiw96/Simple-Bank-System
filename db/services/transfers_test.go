package services

import (
	"testing"

	"simple-bank-system/db/pkg"
	"simple-bank-system/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomWalletTransfer(t *testing.T, account pkg.Account, currency string) (pkg.Wallet, CreateWalletParams) {
	//AccNumb := account.Account_number
	arg := CreateWalletParams{
		WalletNumber: 1015050000,
		Name:         util.RandomOwner(),
		AccountID:    account.ID,
		Balance:      util.RandomMoney(),
		Currency:     currency,
	}
	//log.Println("CreateWalletParams", arg)

	wallet, err := testQueries.CreateWallet(ctx, arg)
	require.Nil(t, err)
	require.NotNil(t, wallet)

	return *wallet, arg
}

func createRandomTransfer(t *testing.T, wallet1, wallet2 pkg.Wallet) (*pkg.Transfers, CreateTransferParam) {
	arg := CreateTransferParam{
		AccountID:        wallet1.AccountID,
		WalletID:         wallet1.ID,
		FromWalletNumber: wallet1.WalletNumber,
		ToWalletNumber:   wallet2.WalletNumber,
		Amount:           util.RandomMoney(),
	}

	res, err := testQueries.CreateTransfer(ctx, arg)
	require.NoError(t, err)
	require.NotNil(t, res)

	return res, arg
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	wallet1_1, _ := createRandomWalletTransfer(t, account1, "IDR")
	wallet1_2, _ := createRandomWalletTransfer(t, account1, "USD")
	account2 := createRandomAccount(t)
	wallet2_1, _ := createRandomWallet(t, account2)

	t.Run("Tx To Wallet With Same Account", func(t *testing.T) {
		transfer, input := createRandomTransfer(t, wallet1_1, wallet1_2)
		assert.Equal(t, wallet1_1.AccountID, transfer.AccountID)
		assert.Equal(t, wallet1_1.ID, transfer.WalletID)
		assert.Equal(t, wallet1_1.WalletNumber, transfer.FromWalletNumber)
		assert.Equal(t, wallet1_2.WalletNumber, transfer.ToWalletNumber)
		assert.Equal(t, input.Amount, transfer.Amount)
	})

	t.Run("Tx To Wallet With Different Account", func(t *testing.T) {
		transfer, input := createRandomTransfer(t, wallet2_1, wallet1_1)
		assert.Equal(t, wallet2_1.AccountID, transfer.AccountID)
		assert.Equal(t, wallet2_1.ID, transfer.WalletID)
		assert.Equal(t, wallet2_1.WalletNumber, transfer.FromWalletNumber)
		assert.Equal(t, wallet1_1.WalletNumber, transfer.ToWalletNumber)
		assert.Equal(t, input.Amount, transfer.Amount)
	})
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	wallet1, _ := createRandomWallet(t, account1)
	account2 := createRandomAccount(t)
	wallet2, _ := createRandomWallet(t, account2)
	transfer1, _ := createRandomTransfer(t, wallet1, wallet2)

	t.Run("Get Transfer Record by Account ID", func(t *testing.T) {
		transfer2, err := testQueries.GetTransfer(ctx, account1.ID, "AccountID")
		require.NoError(t, err, "account ID: %d <<>> %d transfer1 ID", account1.ID, transfer1.ID)
		require.NotNil(t, transfer2, "account ID: %d <<>> %d transfer1 ID", account1.ID, transfer1.ID)

		assert.Equal(t, transfer1.ID, transfer2.ID, "account ID: %d <<>> %d transfer1 ID", account1.ID, transfer1.ID)
		assert.Equal(t, account1.ID, transfer2.AccountID, "account ID: %d <<>> %d transfer1 ID", account1.ID, transfer1.ID)
		assert.Equal(t, wallet1.ID, transfer2.WalletID, "account ID: %d <<>> %d transfer1 ID", account1.ID, transfer1.ID)
		assert.Equal(t, transfer1.WalletID, transfer2.WalletID, "account ID: %d <<>> %d transfer1 ID", account1.ID, transfer1.ID)
		assert.Equal(t, wallet1.WalletNumber, transfer2.FromWalletNumber, "account ID: %d <<>> %d transfer1 ID", account1.ID, transfer1.ID)
		assert.Equal(t, wallet2.WalletNumber, transfer2.ToWalletNumber, "account ID: %d <<>> %d transfer1 ID", account1.ID, transfer1.ID)
		assert.Equal(t, transfer1.Amount, transfer2.Amount, "account ID: %d <<>> %d transfer1 ID", account1.ID, transfer1.ID)
		assert.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt, "account ID: %d <<>> %d transfer1 ID", account1.ID, transfer1.ID)
		assert.Empty(t, transfer2.DeletedAt, "account ID: %d <<>> %d transfer1 ID", account1.ID, transfer1.ID)
	})

	t.Run("Get Transfer Record by Wallet Number", func(t *testing.T) {
		transfer2, err := testQueries.GetTransfer(ctx, wallet1.WalletNumber, "WalletNumber")
		require.NoError(t, err, "wallet ID: %d <<>> %d transfer1 ID", wallet1.ID, transfer1.ID)
		require.NotNil(t, transfer2, "wallet ID: %d <<>> %d transfer1 ID", wallet1.ID, transfer1.ID)

		assert.Equal(t, transfer1.ID, transfer2.ID, "wallet ID: %d <<>> %d transfer1 ID", wallet1.ID, transfer1.ID)
		assert.Equal(t, account1.ID, transfer2.AccountID, "wallet ID: %d <<>> %d transfer1 ID", wallet1.ID, transfer1.ID)
		assert.Equal(t, wallet1.ID, transfer2.WalletID, "wallet ID: %d <<>> %d transfer1 ID", wallet1.ID, transfer1.ID)
		assert.Equal(t, wallet1.WalletNumber, transfer2.FromWalletNumber, "wallet ID: %d <<>> %d transfer1 ID", wallet1.ID, transfer1.ID)
		assert.Equal(t, wallet2.WalletNumber, transfer2.ToWalletNumber, "wallet ID: %d <<>> %d transfer1 ID", wallet1.ID, transfer1.ID)
		assert.Equal(t, transfer1.Amount, transfer2.Amount, "wallet ID: %d <<>> %d transfer1 ID", wallet1.ID, transfer1.ID)
		assert.Equal(t, transfer1.CreatedAt, transfer2.CreatedAt, "wallet ID: %d <<>> %d transfer1 ID", wallet1.ID, transfer1.ID)
		assert.Empty(t, transfer2.DeletedAt, "wallet ID: %d <<>> %d transfer1 ID", wallet1.ID, transfer1.ID)
	})
}

func TestListTransfers(t *testing.T) {
	account1 := createRandomAccount(t)
	wallet1, _ := createRandomWallet(t, account1)
	account2 := createRandomAccount(t)
	wallet2, _ := createRandomWallet(t, account2)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, wallet1, wallet2)
	}

	t.Run("List Transfer By Limit and Offset", func(t *testing.T) {
		arg := ListTransfersParam{
			limit:  5,
			offset: 2,
		}
		transfers, err := testQueries.ListTransfers(ctx, arg)

		require.NoError(t, err)
		require.Len(t, transfers, arg.limit, "amount of transfers list %d != %d input", len(transfers), arg.limit)

		for _, transfer := range transfers {
			assert.NotEmpty(t, transfer)
		}
	})
	t.Run("List Transfer By Date", func(t *testing.T) {
		arg := ListTransfersByDateParam{
			Start: "2023-01-01",
			End:   "2023-09-23",
		}
		transfers, err := testQueries.ListTransfersByDate(ctx, arg, true)

		require.NoError(t, err)

		for _, transfer := range transfers {
			assert.NotEmpty(t, transfer)
		}
	})
}
