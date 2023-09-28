package services

import (
	"testing"

	"simple-bank-system/db/pkg"
	"simple-bank-system/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	account1 pkg.Account
	input    CreateWalletParams
)

func createRandomWallet(t *testing.T, account pkg.Account) (pkg.Wallet, CreateWalletParams) {
	//AccNumb := account.Account_number
	arg := CreateWalletParams{
		WalletNumber: 1015050000,
		Name:         util.RandomOwner(),
		AccountID:    account.ID,
		Balance:      util.RandomMoney(),
		Currency:     util.RandomCurrency(),
	}
	//log.Println("CreateWalletParams", arg)

	wallet, err := testQueries.CreateWallet(ctx, arg)
	require.Nil(t, err)
	require.NotNil(t, wallet)

	return *wallet, arg
}

func createRandomWalletList(t *testing.T, account pkg.Account, currency string) (pkg.Wallet, CreateWalletParams) {
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

func TestCreateWallet(t *testing.T) {
	account := createRandomAccount(t)
	//log.Println("Account", account)
	wallet, input := createRandomWallet(t, account)

	assert.Equal(t, input.Name, wallet.Name)
	//assert.Equal(t, input.)
	assert.Equal(t, input.Balance, wallet.Balance)
	assert.Equal(t, input.Currency, wallet.Currency)
	assert.NotZero(t, wallet.ID, "ID isn't automatically generate")
	assert.False(t, wallet.CreatedAt.IsZero(), "wallet createdAt = %v", wallet.CreatedAt)
	assert.Empty(t, wallet.DeletedAt)
}

func TestCreateWalletErr(t *testing.T) {
	arg := CreateWalletParams{
		WalletNumber: 1015050000,
		Name:         "",
		AccountID:    1,
		Balance:      util.RandomMoney(),
		Currency:     util.RandomCurrency(),
	}
	//log.Println("CreateWalletParams", arg)

	_, err := testQueries.CreateWallet(ctx, arg)
	require.Error(t, err)
	//require.NotNil(t, wallet)
}

func TestGetWallet(t *testing.T) {
	account := createRandomAccount(t)
	wallet1, _ := createRandomWallet(t, account)

	t.Run("Get Wallet by Wallet ID", func(t *testing.T) {
		wallet2, err := testQueries.GetWallet(ctx, wallet1.ID)

		require.NoError(t, err)
		require.NotNil(t, wallet1)
		require.NotNil(t, wallet2, "Can't read wallet from database")

		assert.Equal(t, wallet1.ID, wallet2.ID)
		assert.Equal(t, wallet1.WalletNumber, wallet2.WalletNumber)
		assert.Equal(t, wallet1.AccountID, wallet2.AccountID)
		assert.Equal(t, wallet1.Name, wallet2.Name)
		assert.Equal(t, wallet1.Balance, wallet2.Balance)
		assert.Equal(t, wallet1.Currency, wallet2.Currency)
		assert.Equal(t, wallet1.CreatedAt, wallet2.CreatedAt)
		assert.Empty(t, wallet2.DeletedAt)
	})

	t.Run("Get Wallet by Wallet Number", func(t *testing.T) {
		wallet2, err := testQueries.GetWalletByNumber(ctx, wallet1.WalletNumber)

		require.Nil(t, err)
		require.NotNil(t, wallet1)
		require.NotNil(t, wallet2, "Can't read wallet from database")

		assert.Equal(t, wallet1.ID, wallet2.ID)
		assert.Equal(t, wallet1.WalletNumber, wallet2.WalletNumber)
		assert.Equal(t, wallet1.AccountID, wallet2.AccountID)
		assert.Equal(t, wallet1.Name, wallet2.Name)
		assert.Equal(t, wallet1.Balance, wallet2.Balance)
		assert.Equal(t, wallet1.Currency, wallet2.Currency)
		assert.Equal(t, wallet1.CreatedAt, wallet2.CreatedAt)
		assert.Empty(t, wallet2.DeletedAt)
	})
}

func TestUpdateWallet(t *testing.T) {
	account := createRandomAccount(t)
	wallet1, _ := createRandomWallet(t, account)
	input := UpdateWalletParams{
		ID:      wallet1.ID,
		Balance: util.RandomMoney(),
	}

	err := testQueries.UpdateWallet(ctx, input)
	require.Nil(t, err)

	wallet2, err := testQueries.GetWallet(ctx, input.ID)
	require.Nil(t, err)

	require.NotNil(t, wallet1)
	require.NotNil(t, wallet2, "Can't read wallet from database")

	assert.Equal(t, wallet1.ID, wallet2.ID)
	assert.Equal(t, wallet1.WalletNumber, wallet2.WalletNumber)
	assert.Equal(t, wallet1.AccountID, wallet2.AccountID)
	assert.Equal(t, wallet1.Name, wallet2.Name)
	assert.Equal(t, input.Balance, wallet2.Balance)
	assert.Equal(t, wallet1.Currency, wallet2.Currency)
	assert.Equal(t, wallet1.CreatedAt, wallet2.CreatedAt)
	assert.Empty(t, wallet2.DeletedAt)
}

func TestUpdateWalletInformation(t *testing.T) {
	account := createRandomAccount(t)
	wallet1, _ := createRandomWallet(t, account)

	newCurrency := util.RandomCurrency()
	if newCurrency == wallet1.Currency {
		newCurrency = util.RandomCurrency()
	}

	input := UpdateWalletInformationParams{
		WalletNumber: wallet1.WalletNumber,
		Name:         util.RandomOwner(),
		Currency:     newCurrency,
	}

	err := testQueries.UpdateWalletInformation(ctx, input)
	require.Nil(t, err)

	wallet2, err := testQueries.GetWallet(ctx, wallet1.ID)
	require.Nil(t, err)

	require.NotNil(t, wallet1)
	require.NotNil(t, wallet2, "Can't read wallet from database")

	assert.Equal(t, wallet1.ID, wallet2.ID)
	assert.Equal(t, wallet1.WalletNumber, wallet2.WalletNumber)
	assert.Equal(t, wallet1.AccountID, wallet2.AccountID)
	assert.Equal(t, input.Name, wallet2.Name)
	assert.Equal(t, wallet1.Balance, wallet2.Balance)
	assert.Equal(t, input.Currency, wallet2.Currency)
	assert.Equal(t, wallet1.CreatedAt, wallet2.CreatedAt)
	assert.Empty(t, wallet2.DeletedAt)
}

func TestDeleteWallet(t *testing.T) {
	account := createRandomAccount(t)
	wallet1, _ := createRandomWallet(t, account)

	err := testQueries.DeleteWallet(ctx, wallet1.ID)
	require.NoError(t, err)

	wallet2, err := testQueries.GetWallet(ctx, wallet1.ID)
	require.Error(t, err)
	//assert.EqualError(t, err, pgx.ErrNoRows.Error())
	assert.Nil(t, wallet2)
}

func TestListWallet(t *testing.T) {
	account := createRandomAccount(t)
	var lastWallet pkg.Wallet
	currency := []string{"IDR", "USD", "EUR", "YEN"}
	for i := 0; i < 4; i++ {
		lastWallet, _ = createRandomWalletList(t, account, currency[i])
	}

	arg := ListWalletParams{
		AccountID: account.ID,
		Limit:     4,
		Offset:    1,
	}
	wallets, err := testQueries.ListWallet(ctx, arg)
	require.Nil(t, err)
	assert.Equal(t, 3, len(wallets), "can't get the right amount of wallet")

	for _, wallet := range wallets {
		assert.NotEmpty(t, wallet)
		assert.Equal(t, lastWallet.AccountID, wallet.AccountID)
	}
}
