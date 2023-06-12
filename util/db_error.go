package util

import (
	"errors"
)

var (
	ErrUser         = errors.New("username/email is already exists or empty")
	ErrAccUser      = errors.New("owner data is incorrect or user does not exist (SQLSTATE 23503)")
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExist     = errors.New("record does not exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)
