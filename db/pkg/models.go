package pkg

import (
	"database/sql"
	"time"
)

type Account struct {
	ID               int64
	AccountNumber    int64
	Username         string
	HashedPassword   string
	FullName         string
	DateOfBirth      time.Time
	Address          Addresses
	Email            string
	PasswordChangeAt time.Time
	CreatedAt        time.Time
	DeletedAt        sql.NullTime
}

type Addresses struct {
	ID        int64
	Provinces string `json:"province" validate:"required"`
	City      string `json:"city" validate:"required"`
	ZIP       int64  `json:"zip" validate:"required"`
	Street    string `json:"street" validate:"required"`
}

type Wallet struct {
	ID           int64
	Name         string
	AccountID    int64
	WalletNumber int64
	Balance      int64
	Currency     string
	CreatedAt    time.Time
	DeletedAt    sql.NullTime
}

type Entry struct {
	ID           int64
	AccountID    int64
	WalletID     int64
	WalletNumber int64
	Amount       int64
	CreatedAt    time.Time
	DeletedAt    sql.NullTime
}

type Transfers struct {
	ID               int64
	AccountID        int64
	WalletID         int64
	FromWalletNumber int64
	ToWalletNumber   int64
	Amount           int64
	CreatedAt        time.Time
	DeletedAt        sql.NullTime
}
