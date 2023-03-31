package services

import (
	"testing"

	"simple-bank-system/db/pkg"
	"simple-bank-system/util"
)

func createRandomTransfer(t *testing.T, account1, account2 pkg.Account) (*pkg.Transfers, CreateTransferParam) {
	arg := CreateTransferParam{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	res, err := testQueries.CreateTransfer(ctx, arg)
	if err != nil {
		t.Fatalf("CreateTransfer() err: %v", err)
	} else if res == nil {
		t.Fatalf("transfer is empty")
	}
	return res, arg
}

func TestCreateTransfer(t *testing.T) {
	account1, _ := createRandomAccount(t)
	account2, _ := createRandomAccount(t)
	transfer, input := createRandomTransfer(t, account1, account2)

	if transfer.FromAccountID != account1.ID {
		t.Fatalf("From ID DB %d != %d ID input", transfer.FromAccountID, account1.ID)
	} else if transfer.ToAccountID != account2.ID {
		t.Fatalf("To ID DB %d != %d ID input", transfer.ToAccountID, account2.ID)
	} else if transfer.Amount != input.Amount {
		t.Fatalf("DB amount %d != %d input", transfer.Amount, input.Amount)
	}
}

func TestGetTransfer(t *testing.T) {
	account1, _ := createRandomAccount(t)
	account2, _ := createRandomAccount(t)
	transfer1, _ := createRandomTransfer(t, account1, account2)
	transfer2, err := testQueries.GetTransfer(ctx, transfer1.ID)

	if err != nil {
		t.Fatalf("GetTransfer() err: %v", err)
	}

	if transfer2 == nil {
		t.Fatalf("can't read transfer")
	} else if transfer2.ID != transfer1.ID {
		t.Fatalf("ID output %d != %d input", transfer2.ID, transfer1.ID)
	} else if transfer2.FromAccountID != transfer1.FromAccountID {
		t.Fatalf("From Account ID output %d != %d input", transfer2.FromAccountID, transfer2.FromAccountID)
	} else if transfer2.ToAccountID != transfer1.ToAccountID {
		t.Fatalf("To Account ID output %d != %d input", transfer2.ToAccountID, transfer1.ToAccountID)
	} else if transfer2.Amount != transfer1.Amount {
		t.Fatalf("Amount output %d != %d input", transfer2.Amount, transfer1.Amount)
	} else if transfer2.CreatedAt != transfer1.CreatedAt {
		t.Fatalf("date output %v != %v input", transfer2.CreatedAt, transfer1.CreatedAt)
	}
}

func TestListTransfers(t *testing.T) {
	account1, _ := createRandomAccount(t)
	account2, _ := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomTransfer(t, account1, account2)
	}

	arg := ListTransfersParam{
		limit:  5,
		offset: 2,
	}
	transfers, err := testQueries.ListTransfers(ctx, arg)
	if err != nil {
		t.Fatalf("ListTransfers() err: %v", err)
	} else if len(transfers) != arg.limit {
		t.Fatalf("amount of transfers list %d != %d input", len(transfers), arg.limit)
	}

	for i, transfer := range transfers {
		if transfer == (pkg.Transfers{}) {
			t.Fatalf("transfers list %d is empty", i)
		}
	}
}
