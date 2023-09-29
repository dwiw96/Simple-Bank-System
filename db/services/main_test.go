package services

import (
	"context"
	"log"
	"os"
	"simple-bank-system/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	testQueries *DB
	ctx         context.Context
	cancel      context.CancelFunc
	dbpool      *pgxpool.Pool
)

func TestMain(m *testing.M) {
	log.Println("--- TestMain()")
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	dbpool, err = pgxpool.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB: ", err)
	}
	defer dbpool.Close()

	testQueries = NewDB(dbpool)

	ctx, cancel = context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	log.Println("--- (1) TestMain()")

	os.Exit(m.Run())
	log.Println("--- (2) TestMain()")
}
