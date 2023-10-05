package api

import (
	"context"
	"log"
	"os"
	"simple-bank-system/db/services"
	"simple-bank-system/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	server *Server
)

func TestMain(m *testing.M) {
	log.Println("--- Test Main()")

	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	dbpool, err := pgxpool.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Failed connect to db: ", err)
	}
	defer dbpool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Second)
	defer cancel()

	store := services.NewStore(dbpool)
	server, err = NewServer(store, ctx, config)
	if err != nil {
		log.Fatal("Can't create server, \nerr: ", err)
	}

	log.Println("--- (1) Test Main()")
	server.Start(config.ServerAddress)

	os.Exit(m.Run())
}
