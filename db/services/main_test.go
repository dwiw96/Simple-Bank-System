package services

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	//"Simple-Bank-System/db/services"
)

var (
	testQueries *DB
	ctx         context.Context
	cancel      context.CancelFunc
	dbpool      *pgxpool.Pool
)

func TestMain(m *testing.M) {
	var err error
	dbpool, err = pgxpool.Connect(context.Background(), "postgresql://db:secret@localhost:5432/bank?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to DB: ", err)
	}
	defer dbpool.Close()

	testQueries = NewDB(dbpool)

	ctx, cancel = context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	//NewCtx(ctx)

	os.Exit(m.Run())
}
