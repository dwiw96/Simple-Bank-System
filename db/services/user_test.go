package services

import (
	"simple-bank-system/db/pkg"
	"simple-bank-system/util"
	"testing"
)

func createRandomUser(t *testing.T) pkg.User {
	hashedPass, err := util.HashingPassword(util.RandomString(6))
	if err != nil {
		t.Error("Failed to hashing password, err: ", err)
	}
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPass,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(ctx, arg)
	if err != nil {
		t.Fatalf("CreateUser QueryRow error: %v", err)
	} else if user == nil {
		t.Fatal("user is empty")
	}

	if user.Username != arg.Username {
		t.Errorf("Username name in database \"%s\" isn't same as arg \"%s\"", user.Username, arg.Username)
	}
	if user.HashedPassword != arg.HashedPassword {
		t.Errorf("HashedPassword in database \"%s\" isn't same as arg \"%s\"", user.HashedPassword, arg.HashedPassword)
	}
	if user.FullName != arg.FullName {
		t.Errorf("FullName in database \"%s\" isn't same as arg \"%s\"", user.FullName, arg.FullName)
	}
	if user.Email != arg.Email {
		t.Errorf("Email in database \"%s\" isn't same as arg \"%s\"", user.Email, arg.Email)
	}
	if user.PasswordChangeAt.IsZero() == true {
		t.Error("PasswordChangeAt isn't automatically generate")
	}
	if user.CreatedAt.IsZero() == true {
		t.Error("created_at is nill")
	}

	return *user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(ctx, user1.Username)
	if err != nil {
		t.Fatalf("GetUser QueryRow error: %v", err)
	}

	if user2.Username != user1.Username {
		t.Errorf("Username name in database \"%s\" isn't same as arg \"%s\"", user2.Username, user1.Username)
	}
	if user2.HashedPassword != user1.HashedPassword {
		t.Errorf("HashedPassword in database \"%s\" isn't same as arg \"%s\"", user2.HashedPassword, user1.HashedPassword)
	}
	if user2.FullName != user1.FullName {
		t.Errorf("FullName in database \"%s\" isn't same as arg \"%s\"", user2.FullName, user1.FullName)
	}
	if user2.Email != user1.Email {
		t.Errorf("Email in database \"%s\" isn't same as arg \"%s\"", user2.Email, user1.Email)
	}
	if user2.PasswordChangeAt.IsZero() == true {
		t.Error("PasswordChangeAt isn't automatically generate")
	}
	if user2.CreatedAt.IsZero() == true {
		t.Error("created_at is nill")
	}
}
