package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"simple-bank-system/db/pkg"
	"simple-bank-system/util"

	//"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var password string

// common recieve interface{} as argument so this func can recieve different struct
func common(t *testing.T, arg interface{}, method string, path string) []byte {
	// Marshall data to json (like json_encode)
	argMarshaled, err := json.Marshal(arg)
	//fmt.Printf("arg = \n%v\n", arg)
	require.NoError(t, err, "Failed to Marshall input body")

	/*
	 * http.Request that can be created using httptest.NewRequest exported function
	 ** NewRequest returns a new incoming server Request, suitable for passing to an http.Handler for testing.
	 *
	 * http.ResponseWriter that can be created by using httptest.NewRecorder type which returns a httptest.ResponseRecorder
	 ** ResponseRecorder is an implementation of http.ResponseWriter that records its mutations for later inspection in tests.
	 */
	req := httptest.NewRequest(method, path, bytes.NewReader(argMarshaled))
	// Add some header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer secret")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	w := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	server.router.ServeHTTP(w, req)

	resBody, err := io.ReadAll(w.Body)
	require.NoError(t, err, "can't read all body of response")

	return resBody
}

func mockCreateRandomAccount(t *testing.T) accountResponse {
	username := util.RandomOwner()

	password = util.RandomPassword()
	require.NotEmpty(t, password, "Password is nil")
	require.Len(t, password, 6, "Password is less than 6")

	prov, city, zip := util.RandomAdress()

	arg := createAccountRequest{
		Username:    username,
		Password:    password,
		FullName:    util.RandomOwner(),
		DateOfBirth: util.RandomDate(),
		Address: pkg.Addresses{
			Provinces: prov,
			City:      city,
			ZIP:       zip,
			Street:    "Jalan Raya Labuan, Km.19, Rt/Rw 01/02",
		},
		Email: util.RandomEmailWithUsername(username),
	}

	resBody := common(t, arg, "POST", "/account")

	// read body
	var response accountResponse
	//response := string(resBody) // convert response body to string
	json.Unmarshal(resBody, &response) // convert json directly to struct
	require.NotEmpty(t, response, "Response body is empty")

	assert.NotZero(t, response.AccountNumber)
	assert.Equal(t, arg.Username, response.Username)
	assert.Equal(t, arg.FullName, response.FullName)
	assert.Equal(t, arg.DateOfBirth, response.DateOfBirth)
	assert.NotZero(t, response.Address)
	assert.Equal(t, arg.Address.Provinces, response.Address.Provinces)
	assert.Equal(t, arg.Address.City, response.Address.City)
	assert.Equal(t, arg.Address.ZIP, response.Address.ZIP)
	assert.Equal(t, arg.Address.Street, response.Address.Street)
	assert.Equal(t, arg.Email, response.Email)
	assert.NotZero(t, response.CreatedAt)
	assert.NotZero(t, response.PasswordChangeAt)

	return response
}

func TestMockCreateAccount(t *testing.T) {
	mockCreateRandomAccount(t)
}

func mockLoginAccount(t *testing.T) {
	account := mockCreateRandomAccount(t)

	arg := loginRequest{
		Username: account.Username,
		Password: password,
	}

	resBody := common(t, arg, "POST", "/account/login")

	// read body
	var response loginResponse
	//response := string(resBody) // convert response body to string
	json.Unmarshal(resBody, &response) // convert json directly to struct
	require.NotEmpty(t, response, "Response body is empty")

	assert.Equal(t, account.AccountNumber, response.Account.AccountNumber)
	assert.Equal(t, arg.Username, response.Account.Username)
	assert.Equal(t, account.FullName, response.Account.FullName)
	assert.Equal(t, account.DateOfBirth, response.Account.DateOfBirth)
	assert.Equal(t, account.Address, response.Account.Address)
	assert.Equal(t, account.Address.Provinces, response.Account.Address.Provinces)
	assert.Equal(t, account.Address.City, response.Account.Address.City)
	assert.Equal(t, account.Address.ZIP, response.Account.Address.ZIP)
	assert.Equal(t, account.Address.Street, response.Account.Address.Street)
	assert.Equal(t, account.Email, response.Account.Email)
	assert.Equal(t, account.CreatedAt, response.Account.CreatedAt)
	assert.Equal(t, account.PasswordChangeAt, response.Account.PasswordChangeAt)
}

func TestMockLogin(t *testing.T) {
	mockLoginAccount(t)
}
