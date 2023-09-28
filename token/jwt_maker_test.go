package token

import (
	"errors"
	"testing"
	"time"

	"simple-bank-system/util"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	if err != nil {
		t.Fatal("NewJWTMaker error, \nerr: ", err)
	}

	//username := util.RandomOwner()
	accountID := util.RandomInt(1, 100)
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(accountID, duration)
	if err != nil {
		t.Error("CreateToken error, \nerr: ", err)
	}
	if token == "" {
		t.Error("token is empty")
	}

	payload, err := maker.VerifyToken(token)
	if err != nil {
		t.Error("Vailed to verify token, \nerr: ", err)
	}
	if payload == nil {
		t.Error("payload from verify token is empty")
	}

	if payload.AccountID != accountID {
		t.Errorf("input accountID %d != %d payload.AccountID", accountID, payload.AccountID)
	}
	if payload.IssuedAt.After(issuedAt.Add(1*time.Second)) && payload.IssuedAt.Before(issuedAt.Add(-1*time.Second)) {
		t.Errorf("input issuedAt %v != %v payload.IssuedAt", issuedAt, payload.IssuedAt)
	}
	if payload.ExpiredAt.After(expiredAt.Add(1*time.Second)) && payload.ExpiredAt.Before(expiredAt.Add(-1*time.Second)) {
		t.Errorf("input expiredAt %v != %v payload.ExpiredAT", expiredAt, payload.ExpiredAt)
	}
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32))
	if err != nil {
		t.Fatal("Failed to create jwt maker, \nerr: ", err)
	}

	//username := util.RandomOwner()
	accountID := util.RandomInt(1, 100)
	duration := -time.Minute

	token, err := maker.CreateToken(accountID, duration)
	if err != nil {
		t.Error("Failed to create token, \nerr: ", err)
	}

	payload, err := maker.VerifyToken(token)
	if err == nil {
		t.Error("Failed to recognize expired token")
	}
	if payload != nil {
		t.Error("payload isn't empty for expired token, \npayload: ", payload)
	}
}

func TestInvalidJWTTokenAlgNode(t *testing.T) {
	//username := util.RandomOwner()
	accountID := util.RandomInt(1, 100)
	duration := time.Minute

	payload, err := NewPayLoad(accountID, duration)
	if err != nil {
		t.Error("Failed to create new payload, \nerr: ", err)
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	if jwtToken == nil {
		t.Error("jwt token is empty")
	}

	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Error("signed token is empty, \nerr: ", err)
	}

	maker, err := NewJWTMaker(util.RandomString(32))
	if err != nil {
		t.Error("Failed to make new maker, \nerr: ", err)
	}

	payload, err = maker.VerifyToken(token)
	if err == nil {
		t.Error("failed to verify invalid token")
	}
	if payload != nil {
		t.Error("payload isn't empty for invalid token, \npayload: ", payload)
	}
	if !errors.Is(err, ErrInvalidToken) {
		t.Errorf("Error isn't because of invalid tokens, \nerr: %s \nErrInvalidToken: %s", err, ErrInvalidToken)
	}
}
