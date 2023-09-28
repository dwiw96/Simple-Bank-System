package util

import (
	"errors"
)

var (
	//username already exists
	ErrUsernameExists      = errors.New("username already exists")
	ErrUsernameEmpty       = errors.New("username is empty")
	ErrAccountNumberExists = errors.New("accountNumber already exists")
	ErrAccountNumberWrong  = errors.New("accountNumber is wrong")
	ErrPasswordEmpty       = errors.New("password is empty")
	ErrFullnameEmpty       = errors.New("fullname is empty")
	ErrDOBEmpty            = errors.New("date of birth is empty")
	ErrAddressEmpty        = errors.New("address is empty")
	ErrEmailExists         = errors.New("email already exists")
	ErrEmailEmpty          = errors.New("email is empty")

	ErrAccUser      = errors.New("owner data is incorrect or user does not exist (SQLSTATE 23503)")
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("record does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
	ErrValid        = errors.New("wallet isn't valid")
)

var ErrReturn = []error{ErrUsernameExists, ErrUsernameEmpty, ErrAccountNumberExists, ErrAccountNumberWrong, ErrPasswordEmpty, ErrFullnameEmpty, ErrDOBEmpty, ErrAddressEmpty, ErrEmailExists, ErrEmailEmpty}

var (
	ErrValidateUsername = "Key: 'createAccountRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag"
	ErrValidatePassword = "Key: 'createAccountRequest.Username' Error:Field validation for 'Password' failed on the 'required' tag"
	ErrValidateFullname = "Key: 'createAccountRequest.Username' Error:Field validation for 'Fullname' failed on the 'required' tag"
	ErrValidateDOB      = "Key: 'createAccountRequest.Username' Error:Field validation for 'DateOfBirth' failed on the 'required' tag"
	ErrValidateAddress  = "Key: 'createAccountRequest.Username' Error:Field validation for 'Address' failed on the 'required' tag"
	ErrValidateEmail    = "Key: 'createAccountRequest.Username' Error:Field validation for 'Email' failed on the 'required' tag"
	ErrValidateProvince = "Key: 'createAccountRequest.Username' Error:Field validation for 'Province' failed on the 'required' tag"
	ErrValidateCity     = "Key: 'createAccountRequest.Username' Error:Field validation for 'City' failed on the 'required' tag"
	ErrValidateZIP      = "Key: 'createAccountRequest.Username' Error:Field validation for 'ZIP' failed on the 'required' tag"
	ErrValidateStreet   = "Key: 'createAccountRequest.Username' Error:Field validation for 'Street' failed on the 'required' tag"
)
