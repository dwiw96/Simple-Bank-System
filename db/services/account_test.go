package services

import (
	"database/sql"
	"errors"

	//"log"
	"testing"

	"simple-bank-system/db/pkg"
	"simple-bank-system/util"
)

func createRandomAccount(t *testing.T) (pkg.Account, pkg.Account) {
	arg := pkg.Account{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(ctx, arg)
	if err != nil {
		t.Fatalf("QueryRow func error = %v", err)
	} else if account == nil {
		t.Fatalf("account is empty")
	}

	return *account, arg
}
func TestCreateAccount(t *testing.T) {
	account, input := createRandomAccount(t)

	if account.Owner != input.Owner {
		t.Fatalf("owner name in database \"%s\" isn't same as input \"%s\"", account.Owner, input.Owner)
	} else if account.Balance != input.Balance {
		t.Fatalf("Balance in database \"%d\" isn't same as input \"%d\"", account.Balance, input.Balance)
	} else if account.Currency != input.Currency {
		t.Fatalf("Currency in database \"%s\" isn't same as input \"%s\"", account.Currency, input.Currency)
	}

	if account.ID == 0 {
		t.Fatalf("ID isn't automatically generate")
	} else if account.CreatedAt.IsZero() == true {
		t.Fatalf("created_at is nill")
	}
}

func TestGetAccount(t *testing.T) {
	account1, _ := createRandomAccount(t)
	account2, err := testQueries.GetAccount(ctx, account1.ID)

	if err != nil {
		t.Fatal(err)
	}
	if account2 == nil {
		t.Fatalf("Can't read account from database")
	}

	if account1.ID != account2.ID {
		t.Fatalf("ID input \"%d\" != ID output \"%d\"", account1.ID, account2.ID)
	} else if account2.Owner != account2.Owner {
		t.Fatalf("Owner input \"%s\" != Owner output \"%s\"", account1.Owner, account2.Owner)
	} else if account1.Balance != account2.Balance {
		t.Fatalf("Balance input \"%d\" != Balance output \"%d\"", account1.Balance, account2.Balance)
	} else if account1.Currency != account2.Currency {
		t.Fatalf("Currency input \"%s\" != Currency output \"%s\"", account1.Currency, account2.Currency)
	}
}

func TestUpdateAccount(t *testing.T) {
	account1, _ := createRandomAccount(t)
	input := pkg.Account{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	err := testQueries.UpdateAccount(ctx, input.ID, input.Balance)
	if err != nil {
		t.Fatalf("Exec Error = %s", err)
	}

	account2, err := testQueries.GetAccount(ctx, input.ID)
	if err != nil {
		t.Fatalf("GetAccount() function failed = %s", err)
	} else if account2 == nil {
		t.Fatalf("failed to get an account from database")
	}

	if account1.ID != account2.ID {
		t.Fatalf("ID input \"%d\" != ID output \"%d\"", account1.ID, account2.ID)
	} else if account2.Owner != account2.Owner {
		t.Fatalf("Owner input \"%s\" != Owner output \"%s\"", account1.Owner, account2.Owner)
	} else if account1.Balance == account2.Balance {
		t.Fatalf("Balance input \"%d\" isn't updated, Balance output \"%d\"", account1.Balance, account2.Balance)
	} else if account1.Currency != account2.Currency {
		t.Fatalf("Currency input \"%s\" != Currency output \"%s\"", account1.Currency, account2.Currency)
	}
}

func TestDeleteAccount(t *testing.T) {
	account1, _ := createRandomAccount(t)
	if account1 == (pkg.Account{}) {
		t.Fatalf("failed to create random account")
	}
	err := testQueries.DeleteAccount(ctx, account1.ID)
	if err != nil {
		t.Fatalf("DeleteAccount() err : %s", err)
	}

	account2, err := testQueries.GetAccount(ctx, account1.ID)
	if errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("Account still not deleted, %s", sql.ErrNoRows)
	} else if account2 != nil {
		t.Fatalf("Account still not deleted")
	}
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	var (
		limit  int = 5
		offset int = 2
	)
	accounts, err := testQueries.ListAccount(ctx, limit, offset)
	if err != nil {
		t.Fatalf("ListAccount() func err : %s", err)
	} else if len(accounts) != 5 {
		t.Fatalf("can't get the right amount of account")
	}

	for _, account := range accounts {
		if account == (pkg.Account{}) {
			t.Fatalf("the accounts are empty")
		}
	}
}
