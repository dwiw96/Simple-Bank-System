package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTransfer(t *testing.T) {
	accRes := loginAccount(t)
	require.NotEmpty(t, accRes)
	wallet1 := createWalletCurrency(t, accRes.AccessToken, "IDR")
	require.NotEmpty(t, wallet1)

	arg := transferRequest{
		FromWalletNumber: accRes.Account.AccountNumber,
		ToWalletNumber:   wallet1.WalletNumber,
		Amount:           250000,
		Currency:         "IDR",
	}
	argMarshaled, err := json.Marshal(arg)
	require.NoError(t, err)
	require.NotEmpty(t, argMarshaled)

	req, err := http.NewRequest("POST", "http://localhost:8080/transfer", bytes.NewReader(argMarshaled))
	require.NoError(t, err)
	require.NotEmpty(t, req)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accRes.AccessToken)

	client := http.Client{Timeout: 10 * time.Second}
	require.NotEmpty(t, client)
	res, err := client.Do(req)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	defer res.Body.Close()

	var response transferTxResponse

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.NotEmpty(t, resBody)
	err = json.Unmarshal(resBody, &response)
	require.NoError(t, err)

	// Check for transfer record
	assert.NotNil(t, response.Transfer)
	assert.Equal(t, accRes.Account.AccountNumber, response.Transfer.FromWalletNumber)
	assert.Equal(t, wallet1.WalletNumber, response.Transfer.ToWalletNumber)
	assert.Equal(t, arg.Amount, response.Transfer.Amount)
	assert.NotZero(t, response.Transfer.CreatedAt)
	//_, err = server.store.GetTransfer(ctx, response.Transfer.FromWalletNumber, "Last")
	//require.NoError(t, err)

	// Check for `From Entry` Record
	assert.Equal(t, accRes.Account.AccountNumber, response.FromEntry.WalletNumber)
	assert.Equal(t, -arg.Amount, response.FromEntry.Amount)
	assert.NotZero(t, response.FromEntry.CreatedAt)
	//_, err = server.store.GetEntry(ctx, response.FromEntry.WalletNumber, "Last")
	//require.NoError(t, err, "can't get ToEntry data from DB, ToEntryNumber: %d", response.ToEntry.WalletNumber)

	// Check for `To Entry` Record
	assert.Equal(t, wallet1.WalletNumber, response.ToEntry.WalletNumber)
	assert.Equal(t, arg.Amount, response.ToEntry.Amount)
	assert.NotZero(t, response.ToEntry.CreatedAt)
	//_, err = server.store.GetEntry(ctx, response.ToEntry.WalletNumber, "Last")
	//require.NoError(t, err, "can't get ToEntry data from DB, ToEntryNumber: %d", response.ToEntry.WalletNumber)

	// Check for `From Wallet` Record
	require.NotEmpty(t, response.FromWallet)
	assert.Equal(t, accRes.Account.AccountNumber, response.FromWallet.WalletNumber)
	/*fromWallet, err := server.store.GetWalletByNumber(ctx, accRes.Account.AccountNumber)
	require.NoError(t, err)
	require.NotEmpty(t, fromWallet)
	assert.Equal(t, fromWallet.Balance, response.FromWallet.Balance)
	assert.Equal(t, fromWallet.Currency, response.FromWallet.Currency)
	assert.Equal(t, fromWallet.CreatedAt, response.FromWallet.CreatedAt)*/

	// Check for `From Wallet` Record
	require.NotEmpty(t, response.ToWallet)
	assert.Equal(t, wallet1.WalletNumber, response.ToWallet.WalletNumber)
	/*toWallet, err := server.store.GetWalletByNumber(ctx, wallet1.WalletNumber)
	require.NoError(t, err)
	require.NotEmpty(t, toWallet)
	assert.Equal(t, toWallet.Balance, response.ToWallet.Balance)
	assert.Equal(t, wallet1.Currency, response.ToWallet.Currency)
	assert.Equal(t, wallet1.CreatedAt, response.ToWallet.CreatedAt)*/

	//Check balance difference before and after transfer the money
	//fmt.Printf("\n>> tx: %d  -  %d\n", response.FromWallet.Balance, response.ToWallet.Balance)
	/*accBal1 := wallet1.Balance - response.FromWallet.Balance
	accBal2 := response.ToWallet.Balance - wallet2.Balance

	assert.Equal(t, accBal1, accBal2)
	assert.True(t, accBal1 > 0)
	assert.True(t, accBal2 > 0)
	assert.True(t, accBal1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount
	assert.True(t, accBal2%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

	// var k must be unique, 1th transfer = 1, 2th transfer = 2, ...
	k := int((accBal1 / amount))
	require.True(t, k >= 1 && k <= n)
	require.NotContains(t, k_Check, k) // _, ok := k_Check[k]; ok {t.Errorf("k isn't unique: %d", k)
	k_Check[k] = true*/
}
