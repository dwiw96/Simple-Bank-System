package services

import (
	"simple-bank-system/db/pkg"
	"simple-bank-system/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) pkg.Account {
	hashedPass, err := util.HashingPassword(util.RandomString(6))
	if err != nil {
		t.Error("Failed to hashing password, err: ", err)
	}

	dobInput := util.RandomDate()
	date, err := util.GetDOB(dobInput)
	require.NoError(t, err)

	prov, city, zip := util.RandomAdress()
	arg := CreateAccountParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPass,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		DateOfBirth:    date,
		Address: pkg.Addresses{
			Provinces: prov,
			City:      city,
			ZIP:       zip,
			Street:    "Jalan Raya Labuan, Km.19, Rt/Rw 01/02",
		},
	}

	account, err := testQueries.CreateAccount(ctx, arg)

	require.NoError(t, err)
	require.NotNil(t, account)

	assert.Equal(t, arg.Username, account.Username)
	assert.Equal(t, arg.FullName, account.FullName)
	assert.Equal(t, arg.Email, account.Email)
	assert.Equal(t, arg.DateOfBirth, account.DateOfBirth)
	assert.Equal(t, arg.Address.Provinces, account.Address.Provinces)
	assert.Equal(t, arg.Address.City, account.Address.City)
	assert.Equal(t, arg.Address.ZIP, account.Address.ZIP)
	assert.Equal(t, arg.Address.Street, account.Address.Street)
	assert.True(t, account.PasswordChangeAt.IsZero(), "PasswordChangeAt isn't automatically generate")
	assert.NotZero(t, account.CreatedAt)

	return *account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestCreateAccountFail(t *testing.T) {
	account1 := createRandomAccount(t)

	hashedPass, err := util.HashingPassword(util.RandomString(6))
	if err != nil {
		t.Error("Failed to hashing password, err: ", err)
	}

	dobInput := util.RandomDate()
	date, err := util.GetDOB(dobInput)
	require.NoError(t, err)

	prov, city, zip := util.RandomAdress()

	arg := CreateAccountParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPass,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		DateOfBirth:    date,
		Address: pkg.Addresses{
			Provinces: prov,
			City:      city,
			ZIP:       zip,
			Street:    "Jalan Raya Labuan, Km.19, Rt/Rw 01/02",
		},
	}

	tests := []struct {
		name     string
		expected error
	}{
		{
			name:     "existsUsername",
			expected: util.ErrUsernameExists,
		}, {
			name:     "emptyUsername",
			expected: util.ErrUsernameEmpty,
		}, {
			name:     "emptyPassword",
			expected: util.ErrPasswordEmpty,
		}, {
			name:     "emptyFullname",
			expected: util.ErrFullnameEmpty,
		}, {
			name:     "emptyDOB",
			expected: util.ErrDOBEmpty,
		}, {
			name:     "emptyAddress",
			expected: util.ErrAddressEmpty,
		}, {
			name:     "existsEmail",
			expected: util.ErrEmailExists,
		}, {
			name:     "emptyEmail",
			expected: util.ErrEmailEmpty,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			switch test.name {
			case "existsUsername":
				arg.Username = account1.Username
			case "emptyUsername":
				arg.Username = ""
			case "emptyPassword":
				arg.Username = util.RandomOwner()
				arg.HashedPassword = ""
			case "emptyFullname":
				arg.HashedPassword = hashedPass
				arg.FullName = ""
			case "emptyDOB":
				arg.FullName = util.RandomOwner()
				arg.DateOfBirth = time.Time{}
			case "emptyAddress":
				arg.DateOfBirth, _ = util.GetDOB(dobInput)
				arg.Address.City = ""
			case "existsEmail":
				arg.Username = util.RandomOwner()
				arg.Address.City = city
				arg.Email = account1.Email
			case "emptyEmail":
				arg.Email = ""
			}

			account, err := testQueries.CreateAccount(ctx, arg)

			require.Error(t, err)
			require.Empty(t, account)
			require.Equal(t, test.expected, err)
		})
	}
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	t.Run("GetAccByUsername", func(t *testing.T) {
		account2, err := testQueries.GetAccount(ctx, account1.Username)

		require.Nil(t, err)
		require.NotNil(t, account1)
		require.NotNil(t, account2, "Can't read account from database")

		assert.Equal(t, account1.Username, account2.Username)
		//assert.Equal(t, account1.HashedPassword, account2.HashedPassword)
		assert.Equal(t, account1.FullName, account2.FullName)
		assert.Equal(t, account1.Email, account2.Email)
		assert.Equal(t, account1.DateOfBirth, account2.DateOfBirth)
		assert.Equal(t, account1.Address.Provinces, account2.Address.Provinces)
		assert.Equal(t, account1.Address.City, account2.Address.City)
		assert.Equal(t, account1.Address.ZIP, account2.Address.ZIP)
		assert.Equal(t, account1.Address.Street, account2.Address.Street)
		assert.True(t, account2.PasswordChangeAt.IsZero(), "PasswordChangeAt isn't automatically generate")
		assert.NotZero(t, account2.CreatedAt)
		assert.Empty(t, account2.DeletedAt)
	})

	t.Run("GetAccByAccNumber", func(t *testing.T) {
		account2, err := testQueries.GetAccountByNumber(ctx, account1.AccountNumber)

		require.Nil(t, err)
		require.NotNil(t, account1)
		require.NotNil(t, account2, "Can't read account from database")

		assert.Equal(t, account1.Username, account2.Username)
		//assert.Equal(t, account1.HashedPassword, account2.HashedPassword)
		assert.Equal(t, account1.FullName, account2.FullName)
		assert.Equal(t, account1.Email, account2.Email)
		assert.Equal(t, account1.DateOfBirth, account2.DateOfBirth)
		assert.Equal(t, account1.Address.Provinces, account2.Address.Provinces)
		assert.Equal(t, account1.Address.City, account2.Address.City)
		assert.Equal(t, account1.Address.ZIP, account2.Address.ZIP)
		assert.Equal(t, account1.Address.Street, account2.Address.Street)
		assert.True(t, account2.PasswordChangeAt.IsZero(), "PasswordChangeAt isn't automatically generate")
		assert.NotZero(t, account2.CreatedAt)
		assert.Empty(t, account2.DeletedAt)
	})
}
