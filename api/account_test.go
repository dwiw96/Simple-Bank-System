package api

import (
	"bytes"
	"encoding/json"

	//"fmt"
	"io"
	"net/http"
	"time"

	"simple-bank-system/db/pkg"
	"simple-bank-system/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var Password string

func createRandomAccount(t *testing.T) accountResponse {
	username := util.RandomOwner()

	Password = util.RandomPassword()
	require.NotEmpty(t, Password, "Password is nil")
	require.Len(t, Password, 6, "Password is less than 6")

	prov, city, zip := util.RandomAdress()

	arg := createAccountRequest{
		Username:    username,
		Password:    Password,
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

	// Marshall data to json (like json_encode)
	argMarshaled, err := json.Marshal(arg)
	//fmt.Printf("arg = \n%v\n", arg)
	require.NoError(t, err, "Failed to Marshall input body")

	/*
	 * We build a new request with http.NewRequest, first param HTTP verb, second the URL, and third the body.
	 * Note that the request’s body should be of type io.Reader.
	 * No problem! To create an io.Reader from a slice of bytes (marshalled), we use bytes.NewReader(marshalled)
	 */
	req, err := http.NewRequest("POST", "http://localhost:8080/account", bytes.NewReader(argMarshaled))
	require.NoError(t, err, "CreateUser QueryRow error")

	// Add some header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer secret")

	// Now create http client and sent request
	// create http client
	// do not forget to set timeout; otherwise, no timeout!
	client := http.Client{Timeout: 10 * time.Second}
	// send the request
	res, err := client.Do(req)
	require.NoError(t, err, "can't send request")
	assert.True(t, res.StatusCode == 200, "status code is wrong, status code:", res.StatusCode, res.Status)

	// Then to read the body of the response, we need an extra step:

	// we do not forget to close the body to free resources
	// defer will execute that at the end of the current function
	defer res.Body.Close()

	/*
	 * We use ioutil.ReadAll to read the body (res.Body). Why ? because res.Body is of type io.ReadCloser.
	 * Internally the response body “will be streamed on demand as the Body field is read. If the network connection fails or the server terminates the response, Body.Read calls return an error.”
	 *
	 * Note that we call Close on res.Body to ensure that the stream will be closed. the defer statement will ensure that this is called when the surrounding function returns.
	 */
	// read body
	var response accountResponse
	//err = json.NewDecoder(res.Body).Decode(response)

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err, "can't read all body of response")

	//response := string(resBody) // convert response body to string
	json.Unmarshal(resBody, &response) // convert json directly to struct
	require.NotEmpty(t, response, "Response body is empty")

	//fmt.Println("------------//------------//------------")
	//fmt.Printf("response :\n%v\n", response)

	return response
}

func TestCreateAccountHandler(t *testing.T) {
	createRandomAccount(t)
}

/*func TestCreateAccountHandlerFail(t *testing.T) {
	account := createRandomAccount(t)
	username := util.RandomOwner()

	password := util.RandomPassword()
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

	var (
		ErrFormat   = "Format input data is wrong"
		ErrDatabase = "failed to pass data into database"
	)

	tests := []struct {
		name       string
		expected   []string
		statusCode int
	}{
		{
			name:       "existsUsername",
			expected:   []string{ErrDatabase, util.ErrUsernameExists.Error()},
			statusCode: 409,
		}, {
			name:       "emptyUsername",
			expected:   []string{ErrFormat, util.ErrValidateUsername},
			statusCode: 422,
		}, {
			name:       "emptyPassword",
			expected:   []string{ErrFormat, util.ErrValidatePassword},
			statusCode: 422,
		}, {
			name:       "emptyFullname",
			expected:   []string{ErrFormat, util.ErrValidateFullname},
			statusCode: 422,
		}, {
			name:       "emptyDOB",
			expected:   []string{ErrFormat, util.ErrValidateDOB},
			statusCode: 422,
		}, {
			name:       "emptyAddress",
			expected:   []string{ErrFormat, util.ErrValidateAddress},
			statusCode: 422,
		}, {
			name:       "existsEmail",
			expected:   []string{ErrDatabase, util.ErrEmailExists.Error()},
			statusCode: 409,
		}, {
			name:       "emptyEmail",
			expected:   []string{ErrFormat, util.ErrValidateEmail},
			statusCode: 422,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.name {
			case "existsUsername":
				arg.Username = account.Username
			case "emptyUsername":
				arg.Username = ""
			case "emptyPassword":
				arg.Username = util.RandomOwner()
				arg.Password = ""
			case "emptyFullname":
				arg.Password = password
				arg.FullName = ""
			case "emptyDOB":
				arg.FullName = util.RandomOwner()
				arg.DateOfBirth = ""
			case "emptyAddress":
				arg.DateOfBirth = util.RandomDate()
				arg.Address.City = ""
			case "existsEmail":
				arg.Username = util.RandomOwner()
				arg.Address.City = city
				arg.Email = account.Email
			case "emptyEmail":
				arg.Email = ""
			}

			argMarshaled, err := json.Marshal(arg)
			fmt.Printf("arg = \n%v\n", arg)
			require.NoError(t, err, "Failed to Marshall input body")

			req, err := http.NewRequest("POST", "http://localhost:8080/account", bytes.NewReader(argMarshaled))
			require.NoError(t, err, "CreateUser QueryRow error")

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer secret")

			client := http.Client{Timeout: 10 * time.Second}
			// send the request
			res, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, test.statusCode, res.StatusCode, "status code is wrong")
			defer res.Body.Close()

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err, "can't read all body of response")

			response := string(resBody) // convert response body to string
			fmt.Println("--- (err)response:", response)
			fields := strings.Split(response, "\n")
			fmt.Println("--- (err)fields:", fields[0])
			fmt.Println("--- (err)fields:", fields[1])
			assert.Equal(t, test.expected[0], fields[0])
			assert.Equal(t, test.expected[1], fields[1])
		})
	}
}*/

func loginAccount(t *testing.T) loginResponse {
	accResponse := createRandomAccount(t)

	arg := loginRequest{
		Username: accResponse.Username,
		Password: Password,
	}

	// Marshall data to json (like json_encode)
	argMarshaled, err := json.Marshal(arg)
	//fmt.Printf("arg = \n%s\n", argMarshaled)
	require.NoError(t, err, "Failed to Marshall input body")

	req, err := http.NewRequest("POST", "http://localhost:8080/account/login", bytes.NewReader(argMarshaled))
	require.NoError(t, err, "CreateUser QueryRow error")

	// Add some header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer secret")

	// Now create http client and sent request
	// create http client
	// do not forget to set timeout; otherwise, no timeout!
	client := http.Client{Timeout: 10 * time.Second}
	// send the request
	res, err := client.Do(req)
	require.NoError(t, err, "can't send request")
	assert.False(t, res.StatusCode != 200, "status code is wrong")

	// Then to read the body of the response, we need an extra step:

	// we do not forget to close the body to free resources
	// defer will execute that at the end of the current function
	defer res.Body.Close()

	var response loginResponse
	//err = json.NewDecoder(res.Body).Decode(response)

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err, "can't read all body of response")

	//response := string(resBody) // convert response body to string
	json.Unmarshal(resBody, &response) // convert json directly to struct
	require.NotEmpty(t, response, "Response body is empty")

	return response
}

func TestLoginAccount(t *testing.T) {
	loginAccount(t)
}
