package services

import (
	"testing"

	"simple-bank-system/db/pkg"
	"simple-bank-system/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomEntries(t *testing.T, accountID int64, wallet pkg.Wallet) (pkg.Entry, CreateEntryParam) {
	arg := CreateEntryParam{
		accountID:    accountID,
		walletID:     wallet.ID,
		walletNumber: wallet.WalletNumber,
		amount:       util.RandomMoney(),
	}

	res, err := testQueries.CreateEntry(ctx, arg)
	require.Nil(t, err, "CreateEntry() err = %s", err)
	require.NotNil(t, res)

	assert.Equal(t, arg.accountID, res.AccountID)
	assert.Equal(t, arg.walletID, res.WalletID)
	assert.Equal(t, arg.walletNumber, res.WalletNumber)
	assert.Equal(t, arg.amount, res.Amount)

	return *res, arg
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	wallet, _ := createRandomWallet(t, account)
	createRandomEntries(t, account.ID, wallet)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	wallet, _ := createRandomWallet(t, account)
	entry1, _ := createRandomEntries(t, account.ID, wallet)
	entry2, err := testQueries.GetEntry(ctx, entry1.ID, "ID")
	require.Nil(t, err, "GetEntry(%d) err = %v", entry1.ID, err)

	assert.Equal(t, entry1.ID, entry2.ID)
	assert.Equal(t, entry1.AccountID, entry2.AccountID)
	assert.Equal(t, entry1.WalletID, entry2.WalletID)
	assert.Equal(t, entry1.WalletNumber, entry2.WalletNumber)
	assert.Equal(t, entry1.Amount, entry2.Amount)
	assert.Equal(t, entry1.CreatedAt, entry2.CreatedAt)
	assert.Empty(t, entry2.DeletedAt)
}

func TestListEntry(t *testing.T) {
	account := createRandomAccount(t)
	wallet1, _ := createRandomWallet(t, account)
	wallet2, _ := createRandomWallet(t, account)
	for i := 0; i < 5; i++ {
		createRandomEntries(t, account.ID, wallet1)
	}
	for i := 0; i < 5; i++ {
		createRandomEntries(t, account.ID, wallet2)
	}

	t.Run("LimitOffset", func(t *testing.T) {
		arg := listEntryParam{
			limit:  5,
			offset: 4,
		}
		entries, err := testQueries.ListEntry(ctx, arg)

		require.NoError(t, err)
		require.Len(t, entries, 5)

		for _, entry := range entries {
			assert.NotEmpty(t, entry)
		}
	})

	t.Run("ListByAccID", func(t *testing.T) {
		entries, err := testQueries.ListEntryByID(ctx, account.ID, false)

		require.NoError(t, err)
		require.Len(t, entries, 10)

		for _, entry := range entries {
			assert.NotEmpty(t, entry)
		}
	})

	t.Run("ListByWalID", func(t *testing.T) {
		entries, err := testQueries.ListEntryByID(ctx, wallet2.ID, true)

		require.NoError(t, err)
		require.Len(t, entries, 5)

		for _, entry := range entries {
			assert.NotEmpty(t, entry)
		}
	})

	t.Run("ListByDate", func(t *testing.T) {
		entries, err := testQueries.ListEntryByDate(ctx, "2023-01-01", "2023-09-23", true)

		require.NoError(t, err)
		//require.Len(t, entries, 5)

		for _, entry := range entries {
			assert.NotEmpty(t, entry)
		}
	})
}
