package util

import (
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	HashedPassword1, err := HashingPassword(password)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("HashedPassword1: ", HashedPassword1)
	if HashedPassword1 == "" {
		t.Errorf("hashed password is empty")
	}

	err = VerifyPassword(password, HashedPassword1)
	if err != nil {
		t.Error("Password is wrong, err: ", err)
	}

	wrongPass := RandomString(6)
	err = VerifyPassword(wrongPass, HashedPassword1)
	if err != bcrypt.ErrMismatchedHashAndPassword {
		t.Error("Test for wrong pass isn't equal to bcrypt.ErrMismatchedHashAndPassword, err: ", err)
	}

	HashedPassword2, err := HashingPassword(password)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("HashedPassword2: ", HashedPassword2)
	if HashedPassword2 == "" {
		t.Errorf("hashed password is empty")
	}
	if HashedPassword1 == HashedPassword2 {
		t.Errorf("Hashed password 1 and 2 are equal, %s == %s", HashedPassword1, HashedPassword2)
	}
	err = VerifyPassword(password, HashedPassword2)
	if err != nil {
		t.Error("Password is wrong, err: ", err)
	}
}
