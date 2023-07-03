package token

import (
	//"fmt"
	"fmt"
	"simple-bank-system/util"
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"
)

func TestPasetoMaker(t *testing.T) {
	//temp := paseto.NewV4SymmetricKey()
	//fmt.Println("NewV4SymmetricKey= ", temp)
	key, err := util.RandomByte(32)
	fmt.Println("key: ", key)
	if err != nil {
		t.Fatal("Failed to generate random 32 byte, \nerr: ", err)
	}
	maker, err := NewPasetoMaker(key)
	if err != nil {
		t.Fatal("Failed to make symmetric key, \nerr: ", err)
	}

	username := util.RandomOwner()
	duration := 1 * time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)

	if err != nil {
		t.Fatal("Failed to create encrypted paseto token, \nerr: ", err)
	}
	if token == "" {
		t.Error("encypted token is empty")
	}

	payload, err := maker.VerifyToken(token)
	if err != nil {
		t.Fatal("Failed to verify encrypted paseto token, \nerr: ", err)
	}
	if payload == nil {
		t.Error("payload is empty")
	}

	if payload.Username != username {
		t.Errorf("input username %s != %s payload.Username", username, payload.Username)
	}
	if payload.IssuedAt.After(issuedAt.Add(1*time.Second)) && payload.IssuedAt.Before(issuedAt.Add(-1*time.Second)) {
		t.Errorf("input issuedAt %v != %v payload.IssuedAt", issuedAt, payload.IssuedAt)
	}
	if payload.ExpiredAt.After(expiredAt.Add(1*time.Second)) && payload.ExpiredAt.Before(expiredAt.Add(-1*time.Second)) {
		t.Errorf("input expiredAt %v != %v payload.ExpiredAT", expiredAt, payload.ExpiredAt)
	}
}

func TestExpiredPasetoToken(t *testing.T) {
	key, err := util.RandomByte(32)
	if err != nil {
		t.Fatal("Failed to generate key, \nerr: ", err)
	}
	maker, err := NewPasetoMaker(key)
	if err != nil {
		t.Fatal("Failed to create paseto maker, \nerr: ", err)
	}

	username := util.RandomOwner()
	duration := -1 * time.Minute

	token, err := maker.CreateToken(username, duration)
	if err != nil {
		t.Fatal("Failed to create encrypted paseto token, \nerr: ", err)
	}
	if token == "" {
		t.Error("encypted token is empty")
	}

	payload, err := maker.VerifyToken(token)
	if err == nil {
		t.Fatal("Failed to detect expired token, \nerr: ", err)
	}
	if payload != nil {
		t.Error("payload isn't empty")
	}
}

func TestInvalidPasetoToken(t *testing.T) {
	key, err := util.RandomByte(32)
	if err != nil {
		t.Fatal("Failed to generate random 32 byte, \nerr: ", err)
	}
	maker, err := NewPasetoMaker(key)
	if err != nil {
		t.Fatal("Failed to make symmetric key, \nerr: ", err)
	}

	username := util.RandomOwner()
	duration := 1 * time.Minute

	token, err := maker.CreateToken(username, duration)

	if err != nil {
		t.Fatal("Failed to create encrypted paseto token, \nerr: ", err)
	}
	if token == "" {
		t.Error("encypted token is empty")
	}

	newToken := paseto.NewToken()
	newKey := paseto.NewV4SymmetricKey() // don't share this!!
	encrypted := newToken.V4Encrypt(newKey, nil)
	payload, err := maker.VerifyToken(encrypted)
	if err == nil {
		t.Fatal("Failed to detect invalid paseto token")
	}
	if payload != nil {
		t.Error("payload isn't empty")
	}

}
