package token

import (
	"simple-bank-system/util"
	"testing"
	"time"

	"aidanwoods.dev/go-paseto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	//temp := paseto.NewV4SymmetricKey()
	//fmt.Println("NewV4SymmetricKey= ", temp)
	key, err := util.RandomByte(32)
	//fmt.Println("key: ", key)
	require.NoError(t, err, "Failed to generate random 32 byte, \nerr: ", err)

	maker, err := NewPasetoMaker(key)
	require.NoError(t, err, "Failed to make symmetric key, \nerr: ", err)

	//username := util.RandomOwner()
	id := util.RandomInt(1, 100)
	duration := 1 * time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	//log.Println("id =", id)

	token, err := maker.CreateToken(id, duration)
	require.NoError(t, err, "Failed to create encrypted paseto token, \nerr: ", err)
	require.NotEmpty(t, token, "encypted token is empty")

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err, "Failed to verify encrypted paseto token")
	require.NotEmpty(t, payload, "payload is empty")

	assert.Equal(t, id, payload.AccountID)
	assert.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	assert.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

	/*if payload.IssuedAt.After(issuedAt.Add(1*time.Second)) && payload.IssuedAt.Before(issuedAt.Add(-1*time.Second)) {
		t.Errorf("input issuedAt %v != %v payload.IssuedAt", issuedAt, payload.IssuedAt)
	}
	if payload.ExpiredAt.After(expiredAt.Add(1*time.Second)) && payload.ExpiredAt.Before(expiredAt.Add(-1*time.Second)) {
		t.Errorf("input expiredAt %v != %v payload.ExpiredAT", expiredAt, payload.ExpiredAt)
	}*/
}

func TestExpiredPasetoToken(t *testing.T) {
	key, err := util.RandomByte(32)
	require.NoError(t, err, "Failed to generate key")

	maker, err := NewPasetoMaker(key)
	require.NoError(t, err, "Failed to create paseto maker")

	//username := util.RandomOwner()
	id := util.RandomInt(1, 100)
	duration := -1 * time.Minute

	token, err := maker.CreateToken(id, duration)
	require.NoError(t, err, "Failed to create encrypted paseto token")
	require.NotEmpty(t, token, "encypted token is empty")

	payload, err := maker.VerifyToken(token)
	assert.Error(t, err, "Failed to  detect expired token")
	assert.Empty(t, payload, "payload isn't empty")
}

func TestInvalidPasetoToken(t *testing.T) {
	key, err := util.RandomByte(32)
	require.NoError(t, err, "Failed to generate random 32 byte")

	maker, err := NewPasetoMaker(key)
	require.NoError(t, err, "Failed to make symmetric key")

	//username := util.RandomOwner()
	id := util.RandomInt(0, 100)
	duration := 1 * time.Minute

	token, err := maker.CreateToken(id, duration)
	require.NoError(t, err, "Failed to create encrypted paseto token")
	require.NotEmpty(t, token, "encypted token is empty")

	newToken := paseto.NewToken()
	newKey := paseto.NewV4SymmetricKey() // don't share this!!
	encrypted := newToken.V4Encrypt(newKey, nil)
	payload, err := maker.VerifyToken(encrypted)

	assert.Error(t, err, "Failed to detect invalid paseto token")
	assert.Empty(t, payload, "payload isn't empty")
}
