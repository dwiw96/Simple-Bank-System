package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	//"simple-bank-system/db/pkg"
	"simple-bank-system/util"
	"testing"
)

var (
	username string
	password string
)

func createRandomUser(t *testing.T) []byte {
	password := util.RandomPassword()
	if password == "" {
		t.Errorf("Password is nil, \npass = %s", password)
	}
	if len(password) != 6 {
		t.Errorf("Password is less than 6, \npass = %s", password)
	}

	arg := createUserRequest{
		Username: util.RandomOwner(),
		Password: password,
		FullName: util.RandomOwner(),
		Email:    util.RandomEmail(),
	}
	username = arg.Username

	// Marshall data to json (like json_encode)
	argMarshaled, err := json.Marshal(arg)
	fmt.Printf("arg = \n%s\n", arg)
	if err != nil {
		t.Errorf("Failed to Marshall input body, \nerr = %s", err)
	}

	return argMarshaled
}

func TestCreateUser(t *testing.T) {
	arg := createRandomUser(t)

	/*
	 * We build a new request with http.NewRequest, first param HTTP verb, second the URL, and third the body.
	 * Note that the request’s body should be of type io.Reader.
	 * No problem! To create an io.Reader from a slice of bytes (marshalled), we use bytes.NewReader(marshalled)
	 */
	req, err := http.NewRequest("POST", "http://localhost:8080/users", bytes.NewReader(arg))
	if err != nil {
		t.Fatalf("CreateUser QueryRow error: %v", err)
	}

	// Add some header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer secret")

	// Now create http client and sent request
	// create http client
	// do not forget to set timeout; otherwise, no timeout!
	client := http.Client{Timeout: 10 * time.Second}
	// send the request
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("impossible to send request: %s", err)
	}
	if res.StatusCode != 200 {
		t.Errorf("status code is wrong, status: %d", res.StatusCode)
	}

	// Then to read the body of the response, we need an extra step:

	// we do not forget to close the body to free resources
	// defer will execute that at the end of the current function
	defer res.Body.Close()

	// read body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("impossible to read all body of response: %s", err)
	}
	response := string(resBody)
	if response == "" {
		t.Errorf("Response body is empty")
	}
	fmt.Println("//----------------------//")
	fmt.Printf("response :\n%s\n", response)

	/*
	 * We use ioutil.ReadAll to read the body (res.Body). Why ? because res.Body is of type io.ReadCloser.
	 * Internally the response body “will be streamed on demand as the Body field is read. If the network connection fails or the server terminates the response, Body.Read calls return an error.”
	 *
	 * Note that we call Close on res.Body to ensure that the stream will be closed. the defer statement will ensure that this is called when the surrounding function returns.
	 */
}

func TestLoginUser(t *testing.T) {
	arg := loginRequest{
		Username: "tyfyek",
		Password: "grpbr2",
	}

	// Marshall data to json (like json_encode)
	argMarshaled, err := json.Marshal(arg)
	fmt.Printf("arg = \n%s\n", argMarshaled)
	if err != nil {
		t.Errorf("Failed to Marshall input body, \nerr = %s", err)
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/users/login", bytes.NewReader(argMarshaled))
	if err != nil {
		t.Fatalf("CreateUser QueryRow error: %v", err)
	}

	// Add some header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer secret")

	// Now create http client and sent request
	// create http client
	// do not forget to set timeout; otherwise, no timeout!
	client := http.Client{Timeout: 10 * time.Second}
	// send the request
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("impossible to send request: %s", err)
	}
	if res.StatusCode != 200 {
		t.Errorf("status code is wrong, status: %d", res.StatusCode)
	}

	// Then to read the body of the response, we need an extra step:

	// we do not forget to close the body to free resources
	// defer will execute that at the end of the current function
	defer res.Body.Close()

	// read body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("impossible to read all body of response: %s", err)
	}
	response := string(resBody)
	if response == "" {
		t.Errorf("response is empty")
	}
	fmt.Println("//----------------------//")
	fmt.Printf("response :\n%s\n", response)
}
