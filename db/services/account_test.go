package services

import (
	"database/sql"
	"errors"

	//"log"
	"testing"

	"simple-bank-system/db/pkg"
	"simple-bank-system/util"
)

func createRandomAccount(t *testing.T) (pkg.Account, CreateAccountParams) {
	arg := CreateAccountParams{
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
		t.Errorf("owner name in database \"%s\" isn't same as input \"%s\"", account.Owner, input.Owner)
	} else if account.Balance != input.Balance {
		t.Errorf("Balance in database \"%d\" isn't same as input \"%d\"", account.Balance, input.Balance)
	} else if account.Currency != input.Currency {
		t.Errorf("Currency in database \"%s\" isn't same as input \"%s\"", account.Currency, input.Currency)
	}

	if account.ID == 0 {
		t.Errorf("ID isn't automatically generate")
	} else if account.CreatedAt.IsZero() == true {
		t.Errorf("created_at is nill")
	}
}

func TestGetAccount(t *testing.T) {
	account1, _ := createRandomAccount(t)
	account2, err := testQueries.GetAccount(ctx, account1.ID)

	if err != nil {
		t.Fatal(err)
	}
	if account2 == nil {
		t.Error("Can't read account from database")
	}

	if account1.ID != account2.ID {
		t.Errorf("ID input \"%d\" != ID output \"%d\"", account1.ID, account2.ID)
	} else if account2.Owner != account2.Owner {
		t.Errorf("Owner input \"%s\" != Owner output \"%s\"", account1.Owner, account2.Owner)
	} else if account1.Balance != account2.Balance {
		t.Errorf("Balance input \"%d\" != Balance output \"%d\"", account1.Balance, account2.Balance)
	} else if account1.Currency != account2.Currency {
		t.Errorf("Currency input \"%s\" != Currency output \"%s\"", account1.Currency, account2.Currency)
	}
}

func TestUpdateAccount(t *testing.T) {
	account1, _ := createRandomAccount(t)
	input := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	err := testQueries.UpdateAccount(ctx, input)
	if err != nil {
		t.Fatalf("Exec Error = %s", err)
	}

	account2, err := testQueries.GetAccount(ctx, input.ID)
	if err != nil {
		t.Errorf("GetAccount() function failed = %s", err)
	} else if account2 == nil {
		t.Errorf("failed to get an account from database")
	}

	if account1.ID != account2.ID {
		t.Errorf("ID input \"%d\" != ID output \"%d\"", account1.ID, account2.ID)
	} else if account2.Owner != account2.Owner {
		t.Errorf("Owner input \"%s\" != Owner output \"%s\"", account1.Owner, account2.Owner)
	} else if account1.Balance == account2.Balance {
		t.Errorf("Balance input \"%d\" isn't updated, Balance output \"%d\"", account1.Balance, account2.Balance)
	} else if account1.Currency != account2.Currency {
		t.Errorf("Currency input \"%s\" != Currency output \"%s\"", account1.Currency, account2.Currency)
	}
}

func TestDeleteAccount(t *testing.T) {
	account1, _ := createRandomAccount(t)
	if account1 == (pkg.Account{}) {
		t.Errorf("failed to create random account")
	}
	err := testQueries.DeleteAccount(ctx, account1.ID)
	if err != nil {
		t.Errorf("DeleteAccount() err : %s", err)
	}

	account2, err := testQueries.GetAccount(ctx, account1.ID)
	if errors.Is(err, sql.ErrNoRows) {
		t.Errorf("Account still not deleted, %s", sql.ErrNoRows)
	} else if account2 != nil {
		t.Errorf("Account still not deleted")
	}
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountParams{
		Limit:  5,
		Offset: 2,
	}
	accounts, err := testQueries.ListAccount(ctx, arg)
	if err != nil {
		t.Fatalf("ListAccount() func err : %s", err)
	} else if len(accounts) != 5 {
		t.Fatalf("can't get the right amount of account")
	}

	for _, account := range accounts {
		if account == (pkg.Account{}) {
			t.Errorf("the accounts are empty")
		}
	}
}
