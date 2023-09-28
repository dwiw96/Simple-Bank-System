package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"simple-bank-system/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// http://localhost:8080/wallet
var url = "http://localhost:8080/wallet"

func createRandomWallet(t *testing.T, token string) walletResponse {
	wallArg := createWalletRequest{
		Name:     util.RandomOwner(),
		Currency: util.RandomCurrency(),
	}

	argMarshaled, err := json.Marshal(wallArg)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", url, bytes.NewReader(argMarshaled))
	require.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	client := http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	require.NoError(t, err, "can't send request")
	assert.True(t, res.StatusCode == 200, "status code is wrong, status code:", res.StatusCode, res.Status)

	defer res.Body.Close()

	var response walletResponse

	resBody, err := io.ReadAll(res.Body)
	err = json.Unmarshal(resBody, &response)
	require.NoError(t, err)

	assert.Equal(t, wallArg.Name, response.Name)
	assert.NotEmpty(t, response.WalletNumber)
	assert.Zero(t, response.Balance)
	ok := util.IsSupportedCurrency(response.Currency)
	assert.True(t, ok)

	return response
}

func TestCreateWallet(t *testing.T) {
	accRes := loginAccount(t)
	createRandomWallet(t, accRes.AccessToken)
}

func getWalletTest(t *testing.T, token string, wallet walletResponse) walletResponse {
	arg := getWalletRequest{
		WalletNumber: wallet.WalletNumber,
	}

	argMarshaled, err := json.Marshal(arg)
	require.NoError(t, err)
	require.NotEmpty(t, argMarshaled)

	/*var url strings.Builder
	url.WriteString("http://localhost:8080/wallet")
	url.WriteString(string(wallet.WalletNumber))
	newUrl := url.String()*/
	newUrl := url + "/" + strconv.FormatInt(wallet.WalletNumber, 10)
	//log.Println("url:", newUrl)
	req, err := http.NewRequest("GET", newUrl, nil)
	require.NoError(t, err)
	require.NotEmpty(t, req)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	client := http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	require.NoError(t, err, "can't send request")
	assert.True(t, res.StatusCode == 200, "status code is wrong, status code:", res.StatusCode, res.Status, res.Body)

	defer res.Body.Close()

	var response walletResponse

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.NotEmpty(t, resBody)

	err = json.Unmarshal(resBody, &response)
	require.NoError(t, err)

	assert.Equal(t, wallet.Name, response.Name)
	assert.Equal(t, wallet.WalletNumber, response.WalletNumber)
	assert.Equal(t, wallet.Balance, response.Balance)
	assert.Equal(t, wallet.Currency, response.Currency)
	assert.Equal(t, wallet.CreatedAt, response.CreatedAt)

	return response
}

func TestGetWallet(t *testing.T) {
	accRes := loginAccount(t)
	//log.Println("token:", accRes.AccessToken)
	//log.Println("account:", accRes.Account.Username, accRes.Account.AccounNumber)
	walRes := createRandomWallet(t, accRes.AccessToken)
	//log.Println("wallet:", walRes.WalletNumber)
	require.NotEmpty(t, walRes)

	getWalletTest(t, accRes.AccessToken, walRes)
	//assert.Equal(t, walRes.CreatedAt, wallet.CreatedAt)
	//assert.Equal(t, walRes.Balance, wallet.Balance)
}

func TestListWallet(t *testing.T) {
	accRes := loginAccount(t)

	var wallet []walletResponse
	for i := 0; i < 10; i++ {
		//log.Println("---for")
		//log.Println("i:", i)
		wallet = append(wallet, createRandomWallet(t, accRes.AccessToken))
		//log.Println("wallet:", wallet)
	}

	arg := listWalletsRequest{
		PageID:   2,
		PageSize: 5,
	}

	argMarshaled, err := json.Marshal(arg)
	require.NoError(t, err)
	require.NotEmpty(t, argMarshaled)

	newUrl := url + "?page_id=" + strconv.Itoa(arg.PageID) + "&page_size=" + strconv.Itoa(arg.PageSize)
	//log.Println("newUrl:", newUrl)
	req, err := http.NewRequest("GET", newUrl, nil)
	require.NoError(t, err)
	require.NotEmpty(t, req)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accRes.AccessToken)

	client := http.Client{Timeout: 10 * time.Second}
	require.NotEmpty(t, client)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	assert.True(t, res.StatusCode == 200)

	defer res.Body.Close()

	var response []walletResponse
	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.NotEmpty(t, resBody)

	err = json.Unmarshal(resBody, &response)
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		assert.Equal(t, wallet[i+5].Name, response[i].Name)
		assert.Equal(t, wallet[i+5].WalletNumber, response[i].WalletNumber)
		assert.Equal(t, wallet[i+5].Balance, response[i].Balance)
		assert.Equal(t, wallet[i+5].Currency, response[i].Currency)
		assert.Equal(t, wallet[i+5].CreatedAt, response[i].CreatedAt)
	}
}

func updateWalletTest(t *testing.T, token string, wallet walletResponse, balance int64) {
	arg := updateWalletRequest{
		WalletNumber: wallet.WalletNumber,
		Balance:      balance,
	}

	argMarshaled, err := json.Marshal(arg)
	require.NoError(t, err)
	require.NotEmpty(t, argMarshaled)

	newUrl := url + "/update/" + strconv.FormatInt(wallet.WalletNumber, 10)
	req, err := http.NewRequest("PUT", newUrl, bytes.NewReader(argMarshaled))
	require.NoError(t, err)
	require.NotEmpty(t, req)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	client := http.Client{Timeout: 10 * time.Second}
	require.NotEmpty(t, client)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	assert.True(t, res.StatusCode == 200)

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.NotEmpty(t, resBody)

	resStr := string(resBody)
	assert.Equal(t, "\"Data modified\"\n", resStr)
}

func TestUpdateWallet(t *testing.T) {
	accRes := loginAccount(t)
	walRes := createRandomWallet(t, accRes.AccessToken)
	updateWalletTest(t, accRes.AccessToken, walRes, 100000)
}

func TestUpdateinfoWallet(t *testing.T) {
	accRes := loginAccount(t)
	walRes := createRandomWallet(t, accRes.AccessToken)

	newCurrency := util.RandomCurrency()
	if newCurrency == walRes.Currency {
		newCurrency = util.RandomCurrency()
	}

	arg := updateWalletInfoRequest{
		WalletNumber: walRes.WalletNumber,
		Name:         util.RandomOwner(),
		Currency:     newCurrency,
	}

	argMarshaled, err := json.Marshal(arg)
	require.NoError(t, err)
	require.NotEmpty(t, argMarshaled)

	newUrl := url + "/updateInfo/" + strconv.FormatInt(walRes.WalletNumber, 10)
	req, err := http.NewRequest("PUT", newUrl, bytes.NewReader(argMarshaled))
	require.NoError(t, err)
	require.NotEmpty(t, req)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accRes.AccessToken)

	client := http.Client{Timeout: 10 * time.Second}
	require.NotEmpty(t, client)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	assert.True(t, res.StatusCode == 200)

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.NotEmpty(t, resBody)

	resStr := string(resBody)
	assert.Equal(t, "\"Data modified\"\n", resStr)
}

func createWalletCurrency(t *testing.T, token string, currency string) walletResponse {
	wallArg := createWalletRequest{
		Name:     util.RandomOwner(),
		Currency: currency,
	}

	argMarshaled, err := json.Marshal(wallArg)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", url, bytes.NewReader(argMarshaled))
	require.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	client := http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	require.NoError(t, err, "can't send request")
	assert.True(t, res.StatusCode == 200, "status code is wrong, status code:", res.StatusCode, res.Status)

	defer res.Body.Close()

	var response walletResponse

	resBody, err := io.ReadAll(res.Body)
	err = json.Unmarshal(resBody, &response)
	require.NoError(t, err)

	assert.Equal(t, wallArg.Name, response.Name)
	assert.NotEmpty(t, response.WalletNumber)
	assert.Zero(t, response.Balance)
	ok := util.IsSupportedCurrency(response.Currency)
	assert.True(t, ok)

	return response
}
func TestDeleteWallet(t *testing.T) {
	accRes := loginAccount(t)
	require.NotEmpty(t, accRes)
	//log.Println("token:", accRes.AccessToken)
	//log.Println("account:", accRes.Account.Username, accRes.Account.AccountNumber)

	walRes := createWalletCurrency(t, accRes.AccessToken, "IDR")
	require.NotEmpty(t, walRes)

	updateWalletTest(t, accRes.AccessToken, walRes, 350000)
	//wallet := getWalletTest(t, accRes.AccessToken, walRes)
	//assert.Equal(t, int64(350000), wallet.Balance)

	newUrl := url + "/delete/" + strconv.FormatInt(walRes.WalletNumber, 10)
	require.NotEmpty(t, newUrl)

	req, err := http.NewRequest("DELETE", newUrl, nil)
	require.NoError(t, err)
	require.NotEmpty(t, req)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accRes.AccessToken)

	client := http.Client{Timeout: 10 * time.Second}
	require.NotEmpty(t, client)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.True(t, res.StatusCode == 200, "status code is wrong, status code:", res.StatusCode, res.Status, res.Body)

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.NotEmpty(t, resBody)

	resStr := string(resBody)
	assert.Equal(t, "\"Wallet deleted!\"\n", resStr)

	/*primaryWalletArg := walletResponse{
		Name:         "Primary Wallet",
		WalletNumber: accRes.Account.AccounNumber,
		Balance:      1000000,
		Currency:     "IDR",
	}
	primaryWallet := getWalletTest(t, accRes.AccessToken, primaryWalletArg)
	assert.Equal(t, 1350000, primaryWallet.Balance)
	assert.Equal(t, accRes.Account.AccounNumber, primaryWallet.WalletNumber)
	assert.Equal(t, primaryWalletArg.Name, primaryWallet.Name)*/
}
