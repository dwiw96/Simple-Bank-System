package services

import (
	"testing"

	"simple-bank-system/db/pkg"
	"simple-bank-system/util"
)

func createRandomEntries(t *testing.T, account pkg.Account) (pkg.Entry, CreateEntryParam) {
	arg := CreateEntryParam{
		accountID: account.ID,
		amount:    util.RandomMoney(),
	}
	res, err := testQueries.CreateEntry(ctx, arg)
	if err != nil {
		t.Fatalf("CreateEntry() err = %s", err)
	} else if res == nil {
		t.Fatalf("Entry is empty")
	}

	return *res, arg
}

func TestCreatedEntry(t *testing.T) {
	account, _ := createRandomAccount(t)
	entry, input := createRandomEntries(t, account)

	if entry.AccountID != input.accountID {
		t.Fatalf("ID account in DB %d != %d ID account input", entry.AccountID, input.accountID)
	} else if entry.Amount != input.amount {
		t.Fatalf("Amount in DB \"%d\" != \"%d\" Amount input", entry.Amount, input.amount)
	}
}

func TestGetEntry(t *testing.T) {
	account, _ := createRandomAccount(t)
	entry1, _ := createRandomEntries(t, account)
	entry2, err := testQueries.GetEntry(ctx, entry1.AccountID)
	if err != nil {
		t.Fatalf("GetEntry(%d) err = %v", entry1.AccountID, err)
	}

	if entry2.ID != entry1.ID {
		t.Fatalf("ID entry DB %d != %d ID input", entry2.ID, entry1.ID)
	} else if entry2.AccountID != entry1.AccountID {
		t.Fatalf("Account ID %d != %d input", entry2.AccountID, entry1.AccountID)
	} else if entry2.Amount != entry1.Amount {
		t.Fatalf("Amount DB %d != %d input", entry2.Amount, entry1.Amount)
	} else if entry2.CreatedAt != entry1.CreatedAt {
		t.Fatalf("date DB %v != %v input", entry2.CreatedAt, entry1.CreatedAt)
	}
}

func TestListEntry(t *testing.T) {
	account, _ := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomEntries(t, account)
	}

	arg := listEntryParam{
		limit:  5,
		offset: 4,
	}
	entries, err := testQueries.ListEntry(ctx, arg)
	if err != nil {
		t.Fatalf("ListEntry() err = %v", err)
	} else if len(entries) != 5 {
		t.Fatalf("amount of data %v <= 5 input ", len(entries))
	}

	for _, entry := range entries {
		if entry == (pkg.Entry{}) {
			t.Fatalf("account is empty")
		}
	}
}
