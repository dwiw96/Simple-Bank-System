package services

import (
	"simple-bank-system/db/pkg"
	"simple-bank-system/util"
	"testing"
)

func createRandomUser(t *testing.T) pkg.User {
	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(ctx, arg)
	if err != nil {
		t.Fatalf("CreateUser QueryRow error: %v", err)
	} else if user == nil {
		t.Fatalf("user is empty")
	}

	if user.Username != arg.Username {
		t.Fatalf("Username name in database \"%s\" isn't same as arg \"%s\"", user.Username, arg.Username)
	}
	if user.HashedPassword != arg.HashedPassword {
		t.Fatalf("HashedPassword in database \"%s\" isn't same as arg \"%s\"", user.HashedPassword, arg.HashedPassword)
	}
	if user.FullName != arg.FullName {
		t.Fatalf("FullName in database \"%s\" isn't same as arg \"%s\"", user.FullName, arg.FullName)
	}
	if user.Email != arg.Email {
		t.Fatalf("Email in database \"%s\" isn't same as arg \"%s\"", user.Email, arg.Email)
	}
	if user.PasswordChangeAt.IsZero() == true {
		t.Fatalf("PasswordChangeAt isn't automatically generate")
	}
	if user.CreatedAt.IsZero() == true {
		t.Fatalf("created_at is nill")
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
		t.Fatalf("Username name in database \"%s\" isn't same as arg \"%s\"", user2.Username, user1.Username)
	}
	if user2.HashedPassword != user1.HashedPassword {
		t.Fatalf("HashedPassword in database \"%s\" isn't same as arg \"%s\"", user2.HashedPassword, user1.HashedPassword)
	}
	if user2.FullName != user1.FullName {
		t.Fatalf("FullName in database \"%s\" isn't same as arg \"%s\"", user2.FullName, user1.FullName)
	}
	if user2.Email != user1.Email {
		t.Fatalf("Email in database \"%s\" isn't same as arg \"%s\"", user2.Email, user1.Email)
	}
	if user2.PasswordChangeAt.IsZero() == true {
		t.Fatalf("PasswordChangeAt isn't automatically generate")
	}
	if user2.CreatedAt.IsZero() == true {
		t.Fatalf("created_at is nill")
	}
}
